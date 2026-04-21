package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
)

type MetricDefinition struct {
	Name        string
	Description string
	Type        MetricType
	Unit        string
	Labels      []string
}

type MetricValue struct {
	Name   string
	Value  float64
	Labels map[string]string
	Time   time.Time
}

type MetricsCollector struct {
	meterProvider metric.MeterProvider
	meter         metric.Meter
	mu            sync.RWMutex

	counters   map[string]metric.Int64Counter
	gauges     map[string]metric.Float64ObservableGauge
	histograms map[string]metric.Float64Histogram

	metricDefinitions map[string]MetricDefinition
	metricValues      map[string][]MetricValue
}

func NewMetricsCollector(serviceName string) (*MetricsCollector, error) {
	meterProvider := otel.GetMeterProvider()
	meter := meterProvider.Meter(
		serviceName,
		metric.WithInstrumentationVersion("1.0.0"),
	)

	collector := &MetricsCollector{
		meterProvider:     meterProvider,
		meter:             meter,
		counters:          make(map[string]metric.Int64Counter),
		gauges:            make(map[string]metric.Float64ObservableGauge),
		histograms:        make(map[string]metric.Float64Histogram),
		metricDefinitions: make(map[string]MetricDefinition),
		metricValues:      make(map[string][]MetricValue),
	}

	if err := collector.registerDefaultMetrics(); err != nil {
		return nil, fmt.Errorf("failed to register default metrics: %w", err)
	}

	return collector, nil
}

func (mc *MetricsCollector) registerDefaultMetrics() error {
	defaultMetrics := []MetricDefinition{
		{
			Name:        "agent.execution.count",
			Description: "Total number of agent executions",
			Type:        MetricTypeCounter,
			Unit:        "1",
			Labels:      []string{"agent_name", "status"},
		},
		{
			Name:        "agent.execution.duration",
			Description: "Duration of agent executions in milliseconds",
			Type:        MetricTypeHistogram,
			Unit:        "ms",
			Labels:      []string{"agent_name"},
		},
		{
			Name:        "tool.invocation.count",
			Description: "Total number of tool invocations",
			Type:        MetricTypeCounter,
			Unit:        "1",
			Labels:      []string{"tool_name", "status"},
		},
		{
			Name:        "tool.invocation.duration",
			Description: "Duration of tool invocations in milliseconds",
			Type:        MetricTypeHistogram,
			Unit:        "ms",
			Labels:      []string{"tool_name"},
		},
		{
			Name:        "token.consumption.total",
			Description: "Total tokens consumed",
			Type:        MetricTypeCounter,
			Unit:        "1",
			Labels:      []string{"model", "type"},
		},
		{
			Name:        "pipeline.execution.count",
			Description: "Total number of pipeline executions",
			Type:        MetricTypeCounter,
			Unit:        "1",
			Labels:      []string{"pipeline_name", "status"},
		},
		{
			Name:        "memory.usage.bytes",
			Description: "Memory usage in bytes",
			Type:        MetricTypeGauge,
			Unit:        "By",
			Labels:      []string{"memory_type"},
		},
		{
			Name:        "skill.generation.count",
			Description: "Total number of skills generated",
			Type:        MetricTypeCounter,
			Unit:        "1",
			Labels:      []string{"skill_type"},
		},
		{
			Name:        "skill.reuse.count",
			Description: "Total number of skill reuses",
			Type:        MetricTypeCounter,
			Unit:        "1",
			Labels:      []string{"skill_name"},
		},
	}

	for _, def := range defaultMetrics {
		if err := mc.RegisterMetric(def); err != nil {
			return fmt.Errorf("failed to register metric %s: %w", def.Name, err)
		}
	}

	return nil
}

func (mc *MetricsCollector) RegisterMetric(def MetricDefinition) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	switch def.Type {
	case MetricTypeCounter:
		counter, err := mc.meter.Int64Counter(
			def.Name,
			metric.WithDescription(def.Description),
			metric.WithUnit(def.Unit),
		)
		if err != nil {
			return fmt.Errorf("failed to create counter %s: %w", def.Name, err)
		}
		mc.counters[def.Name] = counter

	case MetricTypeGauge:
		gauge, err := mc.meter.Float64ObservableGauge(
			def.Name,
			metric.WithDescription(def.Description),
			metric.WithUnit(def.Unit),
		)
		if err != nil {
			return fmt.Errorf("failed to create gauge %s: %w", def.Name, err)
		}
		mc.gauges[def.Name] = gauge

	case MetricTypeHistogram:
		histogram, err := mc.meter.Float64Histogram(
			def.Name,
			metric.WithDescription(def.Description),
			metric.WithUnit(def.Unit),
		)
		if err != nil {
			return fmt.Errorf("failed to create histogram %s: %w", def.Name, err)
		}
		mc.histograms[def.Name] = histogram

	default:
		return fmt.Errorf("unknown metric type: %s", def.Type)
	}

	mc.metricDefinitions[def.Name] = def
	mc.metricValues[def.Name] = make([]MetricValue, 0)

	return nil
}

func (mc *MetricsCollector) RecordCounter(name string, value int64, labels map[string]string) error {
	mc.mu.RLock()
	counter, exists := mc.counters[name]
	mc.mu.RUnlock()

	if !exists {
		return fmt.Errorf("counter %s not registered", name)
	}

	attrs := mc.convertLabelsToAttributes(labels)
	counter.Add(context.Background(), value, metric.WithAttributes(attrs...))

	mc.recordMetricValue(name, float64(value), labels)
	return nil
}

func (mc *MetricsCollector) RecordHistogram(name string, value float64, labels map[string]string) error {
	mc.mu.RLock()
	histogram, exists := mc.histograms[name]
	mc.mu.RUnlock()

	if !exists {
		return fmt.Errorf("histogram %s not registered", name)
	}

	attrs := mc.convertLabelsToAttributes(labels)
	histogram.Record(context.Background(), value, metric.WithAttributes(attrs...))

	mc.recordMetricValue(name, value, labels)
	return nil
}

func (mc *MetricsCollector) SetGauge(name string, value float64, labels map[string]string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if _, exists := mc.gauges[name]; !exists {
		return fmt.Errorf("gauge %s not registered", name)
	}

	mc.recordMetricValue(name, value, labels)
	return nil
}

func (mc *MetricsCollector) recordMetricValue(name string, value float64, labels map[string]string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	metricValue := MetricValue{
		Name:   name,
		Value:  value,
		Labels: labels,
		Time:   time.Now(),
	}

	if values, exists := mc.metricValues[name]; exists {
		mc.metricValues[name] = append(values, metricValue)
		if len(mc.metricValues[name]) > 1000 {
			mc.metricValues[name] = mc.metricValues[name][len(mc.metricValues[name])-1000:]
		}
	}
}

func (mc *MetricsCollector) convertLabelsToAttributes(labels map[string]string) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, len(labels))
	for key, value := range labels {
		attrs = append(attrs, attribute.String(key, value))
	}
	return attrs
}

func (mc *MetricsCollector) GetMetricValues(name string, limit int) []MetricValue {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if values, exists := mc.metricValues[name]; exists {
		if limit > 0 && len(values) > limit {
			return values[len(values)-limit:]
		}
		return values
	}
	return nil
}

func (mc *MetricsCollector) GetMetricDefinition(name string) (MetricDefinition, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	def, exists := mc.metricDefinitions[name]
	return def, exists
}

func (mc *MetricsCollector) GetAllMetrics() map[string][]MetricValue {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	result := make(map[string][]MetricValue)
	for name, values := range mc.metricValues {
		result[name] = make([]MetricValue, len(values))
		copy(result[name], values)
	}
	return result
}

func (mc *MetricsCollector) RecordAgentExecution(agentName string, duration time.Duration, status string) {
	labels := map[string]string{
		"agent_name": agentName,
		"status":     status,
	}

	mc.RecordCounter("agent.execution.count", 1, labels)
	mc.RecordHistogram("agent.execution.duration", float64(duration.Milliseconds()), map[string]string{
		"agent_name": agentName,
	})
}

func (mc *MetricsCollector) RecordToolInvocation(toolName string, duration time.Duration, status string) {
	labels := map[string]string{
		"tool_name": toolName,
		"status":    status,
	}

	mc.RecordCounter("tool.invocation.count", 1, labels)
	mc.RecordHistogram("tool.invocation.duration", float64(duration.Milliseconds()), map[string]string{
		"tool_name": toolName,
	})
}

func (mc *MetricsCollector) RecordTokenConsumption(model string, tokenType string, count int64) {
	mc.RecordCounter("token.consumption.total", count, map[string]string{
		"model": model,
		"type":  tokenType,
	})
}

func (mc *MetricsCollector) RecordPipelineExecution(pipelineName string, status string) {
	mc.RecordCounter("pipeline.execution.count", 1, map[string]string{
		"pipeline_name": pipelineName,
		"status":        status,
	})
}

func (mc *MetricsCollector) RecordMemoryUsage(memoryType string, bytes int64) {
	mc.SetGauge("memory.usage.bytes", float64(bytes), map[string]string{
		"memory_type": memoryType,
	})
}

func (mc *MetricsCollector) RecordSkillGeneration(skillType string) {
	mc.RecordCounter("skill.generation.count", 1, map[string]string{
		"skill_type": skillType,
	})
}

func (mc *MetricsCollector) RecordSkillReuse(skillName string) {
	mc.RecordCounter("skill.reuse.count", 1, map[string]string{
		"skill_name": skillName,
	})
}
