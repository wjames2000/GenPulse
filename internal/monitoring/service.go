package monitoring

import (
	"context"
	"fmt"
	"log"
	"time"

	"GenPulse/internal/monitoring/events"
	"GenPulse/internal/monitoring/metrics"
	"GenPulse/internal/monitoring/storage"
	"GenPulse/internal/monitoring/telemetry"
)

type MonitoringService struct {
	telemetry   *telemetry.Telemetry
	metrics     *metrics.MetricsCollector
	traceStore  *storage.TraceStore
	eventBus    *events.EventBus
	wailsBridge *events.WailsEventBridge
	config      Config
	startTime   time.Time
}

type Config struct {
	Enabled         bool
	TelemetryConfig telemetry.Config
	TraceStorePath  string
	EventBufferSize int
	RetentionDays   int
}

func DefaultConfig() Config {
	return Config{
		Enabled:         true,
		TelemetryConfig: telemetry.DefaultConfig(),
		TraceStorePath:  "data/traces.db",
		EventBufferSize: 1000,
		RetentionDays:   30,
	}
}

func NewMonitoringService(config Config) (*MonitoringService, error) {
	if !config.Enabled {
		return &MonitoringService{
			config:    config,
			startTime: time.Now(),
		}, nil
	}

	otel, err := telemetry.NewTelemetry(config.TelemetryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize telemetry: %w", err)
	}

	metricsCollector, err := metrics.NewMetricsCollector("genpulse")
	if err != nil {
		otel.Shutdown(context.Background())
		return nil, fmt.Errorf("failed to create metrics collector: %w", err)
	}

	traceStore, err := storage.NewTraceStore(config.TraceStorePath)
	if err != nil {
		otel.Shutdown(context.Background())
		return nil, fmt.Errorf("failed to create trace store: %w", err)
	}

	eventBusConfig := events.EventBusConfig{
		BufferSize:    config.EventBufferSize,
		MaxBufferSize: config.EventBufferSize * 10,
	}
	eventBus := events.NewEventBus(eventBusConfig)

	wailsBridge := events.NewWailsEventBridge(eventBus)

	service := &MonitoringService{
		telemetry:   otel,
		metrics:     metricsCollector,
		traceStore:  traceStore,
		eventBus:    eventBus,
		wailsBridge: wailsBridge,
		config:      config,
		startTime:   time.Now(),
	}

	service.setupEventHandlers()

	go service.startCleanupWorker(context.Background())

	return service, nil
}

func (ms *MonitoringService) setupEventHandlers() {
	if !ms.config.Enabled {
		return
	}

	ms.eventBus.Subscribe(events.EventTypeAgentStarted, ms.handleAgentStarted)
	ms.eventBus.Subscribe(events.EventTypeAgentCompleted, ms.handleAgentCompleted)
	ms.eventBus.Subscribe(events.EventTypeToolInvoked, ms.handleToolInvoked)
	ms.eventBus.Subscribe(events.EventTypeToolCompleted, ms.handleToolCompleted)
	ms.eventBus.Subscribe(events.EventTypePipelineStarted, ms.handlePipelineStarted)
	ms.eventBus.Subscribe(events.EventTypePipelineCompleted, ms.handlePipelineCompleted)
	ms.eventBus.Subscribe(events.EventTypeSkillGenerated, ms.handleSkillGenerated)
	ms.eventBus.Subscribe(events.EventTypeErrorOccurred, ms.handleErrorOccurred)
}

func (ms *MonitoringService) handleAgentStarted(ctx context.Context, event events.Event) error {
	agentName, _ := event.Data["agent_name"].(string)
	log.Printf("Agent started: %s", agentName)
	return nil
}

func (ms *MonitoringService) handleAgentCompleted(ctx context.Context, event events.Event) error {
	agentName, _ := event.Data["agent_name"].(string)
	duration, _ := event.Data["duration"].(float64)
	status, _ := event.Data["status"].(string)

	ms.metrics.RecordAgentExecution(agentName, time.Duration(duration)*time.Millisecond, status)

	log.Printf("Agent completed: %s, duration: %vms, status: %s",
		agentName, duration, status)
	return nil
}

func (ms *MonitoringService) handleToolInvoked(ctx context.Context, event events.Event) error {
	toolName, _ := event.Data["tool_name"].(string)
	log.Printf("Tool invoked: %s", toolName)
	return nil
}

func (ms *MonitoringService) handleToolCompleted(ctx context.Context, event events.Event) error {
	toolName, _ := event.Data["tool_name"].(string)
	duration, _ := event.Data["duration"].(float64)
	status, _ := event.Data["status"].(string)

	ms.metrics.RecordToolInvocation(toolName, time.Duration(duration)*time.Millisecond, status)

	log.Printf("Tool completed: %s, duration: %vms, status: %s",
		toolName, duration, status)
	return nil
}

func (ms *MonitoringService) handlePipelineStarted(ctx context.Context, event events.Event) error {
	pipelineName, _ := event.Data["pipeline_name"].(string)
	log.Printf("Pipeline started: %s", pipelineName)
	return nil
}

func (ms *MonitoringService) handlePipelineCompleted(ctx context.Context, event events.Event) error {
	pipelineName, _ := event.Data["pipeline_name"].(string)
	status, _ := event.Data["status"].(string)

	ms.metrics.RecordPipelineExecution(pipelineName, status)

	log.Printf("Pipeline completed: %s, status: %s", pipelineName, status)
	return nil
}

func (ms *MonitoringService) handleSkillGenerated(ctx context.Context, event events.Event) error {
	skillType, _ := event.Data["skill_type"].(string)
	ms.metrics.RecordSkillGeneration(skillType)

	skillName, _ := event.Data["skill_name"].(string)
	log.Printf("Skill generated: %s (%s)", skillName, skillType)
	return nil
}

func (ms *MonitoringService) handleErrorOccurred(ctx context.Context, event events.Event) error {
	source, _ := event.Data["source"].(string)
	errorMsg, _ := event.Data["error"].(string)
	log.Printf("Error occurred: %s - %s", source, errorMsg)
	return nil
}

func (ms *MonitoringService) startCleanupWorker(ctx context.Context) {
	if !ms.config.Enabled {
		return
	}

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := ms.traceStore.CleanupOldData(ctx, ms.config.RetentionDays); err != nil {
				log.Printf("Failed to cleanup old trace data: %v", err)
			}
		}
	}
}

func (ms *MonitoringService) GetTelemetry() *telemetry.Telemetry {
	return ms.telemetry
}

func (ms *MonitoringService) GetMetrics() *metrics.MetricsCollector {
	return ms.metrics
}

func (ms *MonitoringService) GetTraceStore() *storage.TraceStore {
	return ms.traceStore
}

func (ms *MonitoringService) GetEventBus() *events.EventBus {
	return ms.eventBus
}

func (ms *MonitoringService) GetWailsBridge() *events.WailsEventBridge {
	return ms.wailsBridge
}

func (ms *MonitoringService) IsEnabled() bool {
	return ms.config.Enabled
}

func (ms *MonitoringService) GetUptime() time.Duration {
	return time.Since(ms.startTime)
}

func (ms *MonitoringService) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	stats["enabled"] = ms.config.Enabled
	stats["uptime"] = ms.GetUptime().String()
	stats["start_time"] = ms.startTime

	if ms.config.Enabled {
		allMetrics := ms.metrics.GetAllMetrics()
		stats["metric_count"] = len(allMetrics)

		recentEvents := ms.eventBus.GetRecentEvents(10)
		stats["recent_event_count"] = len(recentEvents)
	}

	return stats
}

func (ms *MonitoringService) Shutdown(ctx context.Context) error {
	var errors []error

	if ms.config.Enabled {
		if ms.telemetry != nil {
			if err := ms.telemetry.Shutdown(ctx); err != nil {
				errors = append(errors, fmt.Errorf("telemetry shutdown error: %w", err))
			}
		}

		if ms.traceStore != nil {
			if err := ms.traceStore.Close(); err != nil {
				errors = append(errors, fmt.Errorf("trace store close error: %w", err))
			}
		}

		if ms.eventBus != nil {
			ms.eventBus.Close()
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple errors during shutdown: %v", errors)
	}

	return nil
}

func (ms *MonitoringService) RecordAgentTrace(ctx context.Context, span *storage.TraceSpan) error {
	if !ms.config.Enabled {
		return nil
	}

	return ms.traceStore.StoreSpan(ctx, span)
}

func (ms *MonitoringService) RecordToolTrace(ctx context.Context, span *storage.TraceSpan) error {
	if !ms.config.Enabled {
		return nil
	}

	return ms.traceStore.StoreSpan(ctx, span)
}

func (ms *MonitoringService) PublishEvent(ctx context.Context, eventType events.EventType, source string, data map[string]interface{}) error {
	if !ms.config.Enabled {
		return nil
	}

	return ms.eventBus.Publish(ctx, eventType, source, data)
}
