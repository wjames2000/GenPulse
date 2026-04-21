package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"GenPulse/internal/monitoring/events"
	"GenPulse/internal/monitoring/storage"
)

func main() {
	fmt.Println("=== Simple Monitoring Demo ===")

	ctx := context.Background()

	fmt.Println("1. Creating trace store...")
	traceStore, err := storage.NewTraceStore(":memory:")
	if err != nil {
		log.Fatalf("Failed to create trace store: %v", err)
	}
	defer traceStore.Close()

	fmt.Println("2. Creating event bus...")
	eventBus := events.NewEventBus(events.DefaultConfig())

	fmt.Println("\n=== Recording Demo Data ===")

	demoTraceRecording(ctx, traceStore)
	demoEventPublishing(ctx, eventBus)

	fmt.Println("\n=== Querying Data ===")
	queryTraces(ctx, traceStore)
	queryEvents(eventBus)

	fmt.Println("\n=== Demo Completed ===")
}

func demoTraceRecording(ctx context.Context, ts *storage.TraceStore) {
	fmt.Println("\n--- Recording Traces ---")

	spans := []*storage.TraceSpan{
		{
			ID:        "span_1",
			TraceID:   "trace_1",
			Name:      "agent.execution",
			Kind:      "internal",
			StartTime: time.Now().Add(-5 * time.Minute),
			EndTime:   time.Now().Add(-4*time.Minute - 30*time.Second),
			Duration:  90 * time.Second,
			Attributes: map[string]interface{}{
				"agent.name": "orchestrator",
				"status":     "success",
			},
			Status:    "OK",
			AgentName: "orchestrator",
		},
		{
			ID:        "span_2",
			TraceID:   "trace_1",
			ParentID:  "span_1",
			Name:      "tool.invocation",
			Kind:      "client",
			StartTime: time.Now().Add(-4*time.Minute - 25*time.Second),
			EndTime:   time.Now().Add(-4*time.Minute - 20*time.Second),
			Duration:  5 * time.Second,
			Attributes: map[string]interface{}{
				"tool.name": "fs_write",
				"file":      "/tmp/test.txt",
			},
			Status:   "OK",
			ToolName: "fs_write",
		},
	}

	for i, span := range spans {
		if err := ts.StoreSpan(ctx, span); err != nil {
			log.Printf("Failed to store span %d: %v", i, err)
		} else {
			fmt.Printf("Stored span: %s (%s)\n", span.Name, span.ID)
		}
	}
}

func demoEventPublishing(ctx context.Context, eb *events.EventBus) {
	fmt.Println("\n--- Publishing Events ---")

	eventsToPublish := []struct {
		eventType events.EventType
		source    string
		data      map[string]interface{}
	}{
		{
			eventType: events.EventTypeAgentStarted,
			source:    "demo",
			data: map[string]interface{}{
				"agent_name": "architect",
				"input":      "Design system architecture",
			},
		},
		{
			eventType: events.EventTypeAgentCompleted,
			source:    "demo",
			data: map[string]interface{}{
				"agent_name": "architect",
				"duration":   120000,
				"status":     "success",
			},
		},
		{
			eventType: events.EventTypeToolInvoked,
			source:    "demo",
			data: map[string]interface{}{
				"tool_name": "git_clone",
				"repo":      "https://github.com/example/repo",
			},
		},
		{
			eventType: events.EventTypeErrorOccurred,
			source:    "demo",
			data: map[string]interface{}{
				"source": "pipeline",
				"error":  "Network timeout",
			},
		},
	}

	for i, event := range eventsToPublish {
		if err := eb.Publish(ctx, event.eventType, event.source, event.data); err != nil {
			log.Printf("Failed to publish event %d: %v", i, err)
		} else {
			fmt.Printf("Published event: %s\n", event.eventType)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func queryTraces(ctx context.Context, ts *storage.TraceStore) {
	fmt.Println("\n--- Querying Traces ---")

	recentSpans, err := ts.GetRecentSpans(ctx, 10)
	if err != nil {
		log.Printf("Failed to get recent spans: %v", err)
		return
	}

	fmt.Printf("Found %d spans:\n", len(recentSpans))
	for i, span := range recentSpans {
		idDisplay := span.ID
		if len(span.ID) > 8 {
			idDisplay = span.ID[:8]
		}
		fmt.Printf("%d. %s (ID: %s, Duration: %v)\n",
			i+1, span.Name, idDisplay, span.Duration)
	}

	agentSpans, err := ts.GetAgentSpans(ctx, "orchestrator", 5)
	if err != nil {
		log.Printf("Failed to get agent spans: %v", err)
		return
	}

	fmt.Printf("\nAgent 'orchestrator' spans: %d\n", len(agentSpans))
}

func queryEvents(eb *events.EventBus) {
	fmt.Println("\n--- Querying Events ---")

	recentEvents := eb.GetRecentEvents(5)
	fmt.Printf("Recent events: %d\n", len(recentEvents))
	for i, event := range recentEvents {
		fmt.Printf("%d. [%s] from %s\n", i+1, event.Type, event.Source)
	}

	agentEvents := eb.GetEventsByType(events.EventTypeAgentStarted, 3)
	fmt.Printf("\nAgent started events: %d\n", len(agentEvents))
}
