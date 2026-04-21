package main

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

func main() {
	fmt.Println("=== Monitoring Data Collection Demo ===")

	ctx := context.Background()

	fmt.Println("\n1. Initializing OpenTelemetry...")
	telemetryConfig := telemetry.DefaultConfig()
	telemetryConfig.ConsoleExport = true
	telemetryConfig.OTLPEndpoint = ""

	otel, err := telemetry.NewTelemetry(telemetryConfig)
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer otel.Shutdown(ctx)

	fmt.Println("2. Creating metrics collector...")
	metricsCollector, err := metrics.NewMetricsCollector("genpulse-demo")
	if err != nil {
		log.Fatalf("Failed to create metrics collector: %v", err)
	}

	fmt.Println("3. Creating trace store...")
	traceStore, err := storage.NewTraceStore(":memory:")
	if err != nil {
		log.Fatalf("Failed to create trace store: %v", err)
	}
	defer traceStore.Close()

	fmt.Println("4. Creating event bus...")
	eventBus := events.NewEventBus(events.DefaultConfig())

	fmt.Println("\n=== Starting Demo Scenarios ===")

	demoAgentExecution(ctx, metricsCollector, traceStore, eventBus)
	demoToolInvocation(ctx, metricsCollector, traceStore, eventBus)
	demoPipelineExecution(ctx, metricsCollector, traceStore, eventBus)
	demoSkillGeneration(ctx, metricsCollector, eventBus)

	fmt.Println("\n=== Querying and Displaying Results ===")
	displayMetrics(metricsCollector)
	displayTraces(ctx, traceStore)
	displayEvents(eventBus)

	fmt.Println("\n=== Demo Completed ===")
}

func demoAgentExecution(ctx context.Context, mc *metrics.MetricsCollector, ts *storage.TraceStore, eb *events.EventBus) {
	fmt.Println("\n--- Demo: Agent Execution ---")

	agentName := "orchestrator_agent"
	startTime := time.Now()

	eventType, eventData := events.CreateAgentStartedEvent(agentName, "Create a simple web application")
	if err := eb.Publish(ctx, eventType, "demo", eventData); err != nil {
		log.Printf("Failed to publish agent started event: %v", err)
	}

	time.Sleep(150 * time.Millisecond)

	duration := time.Since(startTime)
	mc.RecordAgentExecution(agentName, duration, "success")

	eventType, eventData = events.CreateAgentCompletedEvent(agentName, "Project structure created", duration, "success")
	if err := eb.Publish(ctx, eventType, "demo", eventData); err != nil {
		log.Printf("Failed to publish agent completed event: %v", err)
	}

	span := &storage.TraceSpan{
		ID:        fmt.Sprintf("span_%d", time.Now().UnixNano()),
		TraceID:   "trace_agent_execution",
		Name:      "agent.execution",
		Kind:      "internal",
		StartTime: startTime,
		EndTime:   time.Now(),
		Duration:  duration,
		Attributes: map[string]interface{}{
			"agent.name": agentName,
			"input":      "Create a simple web application",
			"output":     "Project structure created",
		},
		Status:        "OK",
		StatusMessage: "Agent execution completed successfully",
		AgentName:     agentName,
	}

	if err := ts.StoreSpan(ctx, span); err != nil {
		log.Printf("Failed to store span: %v", err)
	}

	fmt.Printf("Recorded agent execution: %s (duration: %v)\n", agentName, duration)
}

func demoToolInvocation(ctx context.Context, mc *metrics.MetricsCollector, ts *storage.TraceStore, eb *events.EventBus) {
	fmt.Println("\n--- Demo: Tool Invocation ---")

	toolName := "fs_write_file"
	startTime := time.Now()

	eventType, eventData := events.CreateToolInvokedEvent(toolName, map[string]interface{}{
		"path":    "/tmp/demo.txt",
		"content": "Hello, World!",
	})
	if err := eb.Publish(ctx, eventType, "demo", eventData); err != nil {
		log.Printf("Failed to publish tool invoked event: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	duration := time.Since(startTime)
	mc.RecordToolInvocation(toolName, duration, "success")

	eventType, eventData = events.CreateToolCompletedEvent(toolName, "File written successfully", duration, "success")
	if err := eb.Publish(ctx, eventType, "demo", eventData); err != nil {
		log.Printf("Failed to publish tool completed event: %v", err)
	}

	span := &storage.TraceSpan{
		ID:        fmt.Sprintf("span_%d", time.Now().UnixNano()),
		TraceID:   "trace_tool_invocation",
		Name:      "tool.invocation",
		Kind:      "client",
		StartTime: startTime,
		EndTime:   time.Now(),
		Duration:  duration,
		Attributes: map[string]interface{}{
			"tool.name": toolName,
			"path":      "/tmp/demo.txt",
			"success":   true,
		},
		Status:        "OK",
		StatusMessage: "Tool invocation completed successfully",
		ToolName:      toolName,
	}

	if err := ts.StoreSpan(ctx, span); err != nil {
		log.Printf("Failed to store span: %v", err)
	}

	fmt.Printf("Recorded tool invocation: %s (duration: %v)\n", toolName, duration)
}

func demoPipelineExecution(ctx context.Context, mc *metrics.MetricsCollector, ts *storage.TraceStore, eb *events.EventBus) {
	fmt.Println("\n--- Demo: Pipeline Execution ---")

	pipelineID := "pipeline_web_app"
	pipelineName := "Web Application Pipeline"
	startTime := time.Now()

	eventType, eventData := events.CreatePipelineStartedEvent(pipelineID, pipelineName, "Create full-stack web app")
	if err := eb.Publish(ctx, eventType, "demo", eventData); err != nil {
		log.Printf("Failed to publish pipeline started event: %v", err)
	}

	agents := []string{"product_manager", "architect", "frontend_dev", "backend_dev", "qa_engineer"}
	for _, agent := range agents {
		agentStart := time.Now()
		time.Sleep(20 * time.Millisecond)
		agentDuration := time.Since(agentStart)

		mc.RecordAgentExecution(agent, agentDuration, "success")

		span := &storage.TraceSpan{
			ID:        fmt.Sprintf("span_%s_%d", agent, time.Now().UnixNano()),
			TraceID:   "trace_pipeline_execution",
			Name:      fmt.Sprintf("agent.%s", agent),
			Kind:      "internal",
			StartTime: agentStart,
			EndTime:   time.Now(),
			Duration:  agentDuration,
			Attributes: map[string]interface{}{
				"agent.name": agent,
				"pipeline":   pipelineName,
			},
			Status:     "OK",
			AgentName:  agent,
			PipelineID: pipelineID,
		}

		if err := ts.StoreSpan(ctx, span); err != nil {
			log.Printf("Failed to store span for agent %s: %v", agent, err)
		}
	}

	time.Sleep(100 * time.Millisecond)

	duration := time.Since(startTime)
	mc.RecordPipelineExecution(pipelineName, "success")

	eventType, eventData = events.CreatePipelineCompletedEvent(pipelineID, pipelineName, "Web application created successfully", duration, "success")
	if err := eb.Publish(ctx, eventType, "demo", eventData); err != nil {
		log.Printf("Failed to publish pipeline completed event: %v", err)
	}

	fmt.Printf("Recorded pipeline execution: %s with %d agents (duration: %v)\n", pipelineName, len(agents), duration)
}

func demoSkillGeneration(ctx context.Context, mc *metrics.MetricsCollector, eb *events.EventBus) {
	fmt.Println("\n--- Demo: Skill Generation ---")

	skillName := "create_react_component"
	skillType := "frontend"

	mc.RecordSkillGeneration(skillType)

	eventType, eventData := events.CreateSkillGeneratedEvent(skillName, skillType)
	if err := eb.Publish(ctx, eventType, "demo", eventData); err != nil {
		log.Printf("Failed to publish skill generated event: %v", err)
	}

	for i := 0; i < 3; i++ {
		mc.RecordSkillReuse(skillName)
		time.Sleep(10 * time.Millisecond)
	}

	mc.RecordTokenConsumption("gpt-4", "completion", 1500)
	mc.RecordTokenConsumption("gpt-4", "prompt", 3200)

	mc.RecordMemoryUsage("working", 1024*1024)
	mc.RecordMemoryUsage("episodic", 5*1024*1024)

	fmt.Printf("Recorded skill generation: %s (%s) with 3 reuses\n", skillName, skillType)
}

func displayMetrics(mc *metrics.MetricsCollector) {
	fmt.Println("\n--- Current Metrics ---")

	metricNames := []string{
		"agent.execution.count",
		"tool.invocation.count",
		"token.consumption.total",
		"pipeline.execution.count",
		"skill.generation.count",
		"skill.reuse.count",
	}

	for _, name := range metricNames {
		values := mc.GetMetricValues(name, 5)
		if len(values) > 0 {
			fmt.Printf("%s: %d recent values\n", name, len(values))
			for _, val := range values {
				fmt.Printf("  - Value: %.2f, Labels: %v, Time: %v\n",
					val.Value, val.Labels, val.Time.Format("15:04:05"))
			}
		}
	}

	allMetrics := mc.GetAllMetrics()
	fmt.Printf("\nTotal metrics collected: %d\n", len(allMetrics))
}

func displayTraces(ctx context.Context, ts *storage.TraceStore) {
	fmt.Println("\n--- Stored Traces ---")

	recentSpans, err := ts.GetRecentSpans(ctx, 5)
	if err != nil {
		log.Printf("Failed to get recent spans: %v", err)
		return
	}

	fmt.Printf("Recent spans: %d\n", len(recentSpans))
	for i, span := range recentSpans {
		fmt.Printf("%d. %s (Agent: %s, Tool: %s) - %v\n",
			i+1, span.Name, span.AgentName, span.ToolName, span.Duration)
	}

	agentSpans, err := ts.GetAgentSpans(ctx, "orchestrator_agent", 3)
	if err != nil {
		log.Printf("Failed to get agent spans: %v", err)
		return
	}

	fmt.Printf("\nAgent 'orchestrator_agent' spans: %d\n", len(agentSpans))

	stats, err := ts.GetStatistics(ctx, time.Now().Add(-1*time.Hour), time.Now())
	if err != nil {
		log.Printf("Failed to get statistics: %v", err)
		return
	}

	fmt.Println("\n--- Statistics (Last Hour) ---")
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func displayEvents(eb *events.EventBus) {
	fmt.Println("\n--- Recent Events ---")

	recentEvents := eb.GetRecentEvents(5)
	fmt.Printf("Recent events: %d\n", len(recentEvents))
	for i, event := range recentEvents {
		fmt.Printf("%d. [%s] %s from %s at %v\n",
			i+1, event.Type, event.ID, event.Source, event.Timestamp.Format("15:04:05"))
	}

	agentEvents := eb.GetEventsByType(events.EventTypeAgentCompleted, 3)
	fmt.Printf("\nAgent completed events: %d\n", len(agentEvents))
	for i, event := range agentEvents {
		if agentName, ok := event.Data["agent_name"].(string); ok {
			fmt.Printf("%d. Agent: %s, Duration: %vms\n",
				i+1, agentName, event.Data["duration"])
		}
	}
}
