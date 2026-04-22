package monitoring_test

import (
	"context"
	"testing"
	"time"

	"GenPulse/internal/monitoring"
	"GenPulse/internal/monitoring/events"
	"GenPulse/internal/monitoring/storage"
)

func TestMonitoringService(t *testing.T) {
	ctx := context.Background()

	config := monitoring.DefaultConfig()
	config.TraceStorePath = ":memory:"
	config.EventBufferSize = 100

	service, err := monitoring.NewMonitoringService(config)
	if err != nil {
		t.Fatalf("Failed to create monitoring service: %v", err)
	}
	defer service.Shutdown(ctx)

	if !service.IsEnabled() {
		t.Error("Monitoring service should be enabled")
	}

	stats := service.GetStats()
	if stats["enabled"] != true {
		t.Error("Stats should show service as enabled")
	}
}

func TestEventBus(t *testing.T) {
	ctx := context.Background()
	eventBus := events.NewEventBus(events.DefaultConfig())

	eventReceived := false
	eventBus.Subscribe(events.EventTypeAgentStarted, func(ctx context.Context, event events.Event) error {
		eventReceived = true
		if event.Type != events.EventTypeAgentStarted {
			t.Errorf("Expected event type %s, got %s", events.EventTypeAgentStarted, event.Type)
		}
		return nil
	})

	err := eventBus.Publish(ctx, events.EventTypeAgentStarted, "test", map[string]interface{}{
		"agent_name": "test_agent",
	})
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if !eventReceived {
		t.Error("Event handler was not called")
	}

	recentEvents := eventBus.GetRecentEvents(5)
	if len(recentEvents) != 1 {
		t.Errorf("Expected 1 recent event, got %d", len(recentEvents))
	}
}

func TestTraceStore(t *testing.T) {
	ctx := context.Background()
	traceStore, err := storage.NewTraceStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create trace store: %v", err)
	}
	defer traceStore.Close()

	span := &storage.TraceSpan{
		ID:        "test_span_1",
		TraceID:   "test_trace_1",
		Name:      "test.operation",
		Kind:      "internal",
		StartTime: time.Now().Add(-1 * time.Minute),
		EndTime:   time.Now(),
		Duration:  time.Minute,
		Attributes: map[string]interface{}{
			"test_key": "test_value",
		},
		Status: "OK",
	}

	err = traceStore.StoreSpan(ctx, span)
	if err != nil {
		t.Fatalf("Failed to store span: %v", err)
	}

	retrievedSpan, err := traceStore.GetSpan(ctx, "test_span_1")
	if err != nil {
		t.Fatalf("Failed to get span: %v", err)
	}

	if retrievedSpan == nil {
		t.Fatal("Retrieved span should not be nil")
	}

	if retrievedSpan.Name != "test.operation" {
		t.Errorf("Expected span name 'test.operation', got '%s'", retrievedSpan.Name)
	}

	if retrievedSpan.Status != "OK" {
		t.Errorf("Expected span status 'OK', got '%s'", retrievedSpan.Status)
	}

	recentSpans, err := traceStore.GetRecentSpans(ctx, 10)
	if err != nil {
		t.Fatalf("Failed to get recent spans: %v", err)
	}

	if len(recentSpans) != 1 {
		t.Errorf("Expected 1 recent span, got %d", len(recentSpans))
	}
}

func TestEventCreationHelpers(t *testing.T) {
	eventType, data := events.CreateAgentStartedEvent("test_agent", "test_input")
	if eventType != events.EventTypeAgentStarted {
		t.Errorf("Expected event type %s, got %s", events.EventTypeAgentStarted, eventType)
	}

	if agentName, ok := data["agent_name"].(string); !ok || agentName != "test_agent" {
		t.Errorf("Expected agent_name 'test_agent', got %v", data["agent_name"])
	}

	eventType, data = events.CreateToolCompletedEvent("test_tool", "test_result", 100*time.Millisecond, "success")
	if eventType != events.EventTypeToolCompleted {
		t.Errorf("Expected event type %s, got %s", events.EventTypeToolCompleted, eventType)
	}

	durationVal := data["duration"]
	var duration float64
	switch v := durationVal.(type) {
	case float64:
		duration = v
	case int64:
		duration = float64(v)
	case int:
		duration = float64(v)
	default:
		t.Errorf("Unexpected duration type: %T", durationVal)
	}

	if duration < 99 || duration > 101 {
		t.Errorf("Expected duration around 100, got %v", duration)
	}
}

func TestTraceQuery(t *testing.T) {
	ctx := context.Background()
	traceStore, err := storage.NewTraceStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create trace store: %v", err)
	}
	defer traceStore.Close()

	spans := []*storage.TraceSpan{
		{
			ID:        "span_1",
			TraceID:   "trace_1",
			Name:      "agent.orchestrator",
			Kind:      "internal",
			StartTime: time.Now().Add(-10 * time.Minute),
			EndTime:   time.Now().Add(-9 * time.Minute),
			Duration:  time.Minute,
			Status:    "OK",
			AgentName: "orchestrator",
		},
		{
			ID:        "span_2",
			TraceID:   "trace_1",
			Name:      "agent.architect",
			Kind:      "internal",
			StartTime: time.Now().Add(-9 * time.Minute),
			EndTime:   time.Now().Add(-8 * time.Minute),
			Duration:  time.Minute,
			Status:    "OK",
			AgentName: "architect",
		},
		{
			ID:        "span_3",
			TraceID:   "trace_2",
			Name:      "agent.frontend",
			Kind:      "internal",
			StartTime: time.Now().Add(-5 * time.Minute),
			EndTime:   time.Now().Add(-4 * time.Minute),
			Duration:  time.Minute,
			Status:    "ERROR",
			AgentName: "frontend",
		},
	}

	for _, span := range spans {
		err = traceStore.StoreSpan(ctx, span)
		if err != nil {
			t.Fatalf("Failed to store span: %v", err)
		}
	}

	query := storage.TraceQuery{
		AgentName: "orchestrator",
		Limit:     10,
	}

	orchestratorSpans, err := traceStore.QuerySpans(ctx, query)
	if err != nil {
		t.Fatalf("Failed to query spans: %v", err)
	}

	if len(orchestratorSpans) != 1 {
		t.Errorf("Expected 1 orchestrator span, got %d", len(orchestratorSpans))
	}

	errorSpans, err := traceStore.QuerySpans(ctx, storage.TraceQuery{
		Status: "ERROR",
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("Failed to query error spans: %v", err)
	}

	if len(errorSpans) != 1 {
		t.Errorf("Expected 1 error span, got %d", len(errorSpans))
	}

	traceSpans, err := traceStore.GetTraceTree(ctx, "trace_1")
	if err != nil {
		t.Fatalf("Failed to get trace tree: %v", err)
	}

	if len(traceSpans) != 2 {
		t.Errorf("Expected 2 spans in trace_1, got %d", len(traceSpans))
	}
}

func TestEventSubscription(t *testing.T) {
	ctx := context.Background()
	eventBus := events.NewEventBus(events.DefaultConfig())

	eventCount := 0
	handler := func(ctx context.Context, event events.Event) error {
		eventCount++
		return nil
	}

	eventBus.Subscribe(events.EventTypeAgentStarted, handler)
	eventBus.Subscribe(events.EventTypeAgentCompleted, handler)

	err := eventBus.Publish(ctx, events.EventTypeAgentStarted, "test", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	err = eventBus.Publish(ctx, events.EventTypeAgentCompleted, "test", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	err = eventBus.Publish(ctx, events.EventTypeToolInvoked, "test", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if eventCount != 2 {
		t.Errorf("Expected 2 events handled, got %d", eventCount)
	}

	eventBus.Unsubscribe(events.EventTypeAgentStarted, handler)

	err = eventBus.Publish(ctx, events.EventTypeAgentStarted, "test", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if eventCount != 2 {
		t.Errorf("Expected event count to remain 2 after unsubscribe, got %d", eventCount)
	}
}
