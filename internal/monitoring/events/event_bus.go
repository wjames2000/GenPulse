package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type EventType string

const (
	EventTypeAgentStarted      EventType = "agent_started"
	EventTypeAgentCompleted    EventType = "agent_completed"
	EventTypeToolInvoked       EventType = "tool_invoked"
	EventTypeToolCompleted     EventType = "tool_completed"
	EventTypePipelineStarted   EventType = "pipeline_started"
	EventTypePipelineCompleted EventType = "pipeline_completed"
	EventTypeMetricUpdated     EventType = "metric_updated"
	EventTypeTraceRecorded     EventType = "trace_recorded"
	EventTypeSkillGenerated    EventType = "skill_generated"
	EventTypeMemoryUpdated     EventType = "memory_updated"
	EventTypeErrorOccurred     EventType = "error_occurred"
	EventTypeStatusChanged     EventType = "status_changed"
)

type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type EventHandler func(ctx context.Context, event Event) error

type EventBus struct {
	mu              sync.RWMutex
	handlers        map[EventType][]EventHandler
	subscribers     map[string]chan Event
	bufferSize      int
	eventBuffer     []Event
	maxBufferSize   int
	metricsRecorder func(event Event)
}

type EventBusConfig struct {
	BufferSize    int
	MaxBufferSize int
}

func DefaultConfig() EventBusConfig {
	return EventBusConfig{
		BufferSize:    100,
		MaxBufferSize: 10000,
	}
}

func NewEventBus(config EventBusConfig) *EventBus {
	return &EventBus{
		handlers:      make(map[EventType][]EventHandler),
		subscribers:   make(map[string]chan Event),
		bufferSize:    config.BufferSize,
		eventBuffer:   make([]Event, 0, config.MaxBufferSize),
		maxBufferSize: config.MaxBufferSize,
	}
}

func (eb *EventBus) SetMetricsRecorder(recorder func(event Event)) {
	eb.metricsRecorder = recorder
}

func (eb *EventBus) Publish(ctx context.Context, eventType EventType, source string, data map[string]interface{}) error {
	event := Event{
		ID:        generateEventID(),
		Type:      eventType,
		Timestamp: time.Now(),
		Source:    source,
		Data:      data,
		Metadata:  make(map[string]interface{}),
	}

	eb.mu.Lock()
	eb.eventBuffer = append(eb.eventBuffer, event)
	if len(eb.eventBuffer) > eb.maxBufferSize {
		eb.eventBuffer = eb.eventBuffer[len(eb.eventBuffer)-eb.maxBufferSize:]
	}
	eb.mu.Unlock()

	if eb.metricsRecorder != nil {
		eb.metricsRecorder(event)
	}

	go eb.dispatchEvent(ctx, event)
	return nil
}

func (eb *EventBus) dispatchEvent(ctx context.Context, event Event) {
	eb.mu.RLock()
	handlers := eb.handlers[event.Type]
	subscribers := make([]chan Event, 0, len(eb.subscribers))
	for _, ch := range eb.subscribers {
		subscribers = append(subscribers, ch)
	}
	eb.mu.RUnlock()

	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h(ctx, event); err != nil {
				fmt.Printf("Error in event handler for %s: %v\n", event.Type, err)
			}
		}(handler)
	}

	for _, ch := range subscribers {
		select {
		case ch <- event:
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("Timeout sending event to subscriber\n")
		}
	}
}

func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

func (eb *EventBus) Unsubscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	handlers := eb.handlers[eventType]
	for i, h := range handlers {
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
			eb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

func (eb *EventBus) CreateSubscription() (string, <-chan Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	id := generateEventID()
	ch := make(chan Event, eb.bufferSize)
	eb.subscribers[id] = ch

	return id, ch
}

func (eb *EventBus) RemoveSubscription(id string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if ch, exists := eb.subscribers[id]; exists {
		close(ch)
		delete(eb.subscribers, id)
	}
}

func (eb *EventBus) GetRecentEvents(limit int) []Event {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if limit <= 0 || limit > len(eb.eventBuffer) {
		limit = len(eb.eventBuffer)
	}

	events := make([]Event, limit)
	copy(events, eb.eventBuffer[len(eb.eventBuffer)-limit:])

	return events
}

func (eb *EventBus) GetEventsByType(eventType EventType, limit int) []Event {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	var filtered []Event
	for i := len(eb.eventBuffer) - 1; i >= 0 && len(filtered) < limit; i-- {
		if eb.eventBuffer[i].Type == eventType {
			filtered = append(filtered, eb.eventBuffer[i])
		}
	}

	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	}

	return filtered
}

func (eb *EventBus) GetEventsBySource(source string, limit int) []Event {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	var filtered []Event
	for i := len(eb.eventBuffer) - 1; i >= 0 && len(filtered) < limit; i-- {
		if eb.eventBuffer[i].Source == source {
			filtered = append(filtered, eb.eventBuffer[i])
		}
	}

	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	}

	return filtered
}

func (eb *EventBus) ClearBuffer() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.eventBuffer = make([]Event, 0, eb.maxBufferSize)
}

func (eb *EventBus) Close() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	for id, ch := range eb.subscribers {
		close(ch)
		delete(eb.subscribers, id)
	}
}

func generateEventID() string {
	return fmt.Sprintf("event_%d", time.Now().UnixNano())
}

func CreateAgentStartedEvent(agentName string, input interface{}) (EventType, map[string]interface{}) {
	data := map[string]interface{}{
		"agent_name": agentName,
		"input":      input,
		"timestamp":  time.Now(),
	}
	return EventTypeAgentStarted, data
}

func CreateAgentCompletedEvent(agentName string, output interface{}, duration time.Duration, status string) (EventType, map[string]interface{}) {
	data := map[string]interface{}{
		"agent_name": agentName,
		"output":     output,
		"duration":   duration.Milliseconds(),
		"status":     status,
		"timestamp":  time.Now(),
	}
	return EventTypeAgentCompleted, data
}

func CreateToolInvokedEvent(toolName string, params interface{}) (EventType, map[string]interface{}) {
	data := map[string]interface{}{
		"tool_name": toolName,
		"params":    params,
		"timestamp": time.Now(),
	}
	return EventTypeToolInvoked, data
}

func CreateToolCompletedEvent(toolName string, result interface{}, duration time.Duration, status string) (EventType, map[string]interface{}) {
	data := map[string]interface{}{
		"tool_name": toolName,
		"result":    result,
		"duration":  duration.Milliseconds(),
		"status":    status,
		"timestamp": time.Now(),
	}
	return EventTypeToolCompleted, data
}

func CreatePipelineStartedEvent(pipelineID string, pipelineName string, input interface{}) (EventType, map[string]interface{}) {
	data := map[string]interface{}{
		"pipeline_id":   pipelineID,
		"pipeline_name": pipelineName,
		"input":         input,
		"timestamp":     time.Now(),
	}
	return EventTypePipelineStarted, data
}

func CreatePipelineCompletedEvent(pipelineID string, pipelineName string, output interface{}, duration time.Duration, status string) (EventType, map[string]interface{}) {
	data := map[string]interface{}{
		"pipeline_id":   pipelineID,
		"pipeline_name": pipelineName,
		"output":        output,
		"duration":      duration.Milliseconds(),
		"status":        status,
		"timestamp":     time.Now(),
	}
	return EventTypePipelineCompleted, data
}

func CreateMetricUpdatedEvent(metricName string, value float64, labels map[string]string) (EventType, map[string]interface{}) {
	data := map[string]interface{}{
		"metric_name": metricName,
		"value":       value,
		"labels":      labels,
		"timestamp":   time.Now(),
	}
	return EventTypeMetricUpdated, data
}

func CreateTraceRecordedEvent(traceID string, spanName string, duration time.Duration) (EventType, map[string]interface{}) {
	data := map[string]interface{}{
		"trace_id":  traceID,
		"span_name": spanName,
		"duration":  duration.Milliseconds(),
		"timestamp": time.Now(),
	}
	return EventTypeTraceRecorded, data
}

func CreateSkillGeneratedEvent(skillName string, skillType string) (EventType, map[string]interface{}) {
	data := map[string]interface{}{
		"skill_name": skillName,
		"skill_type": skillType,
		"timestamp":  time.Now(),
	}
	return EventTypeSkillGenerated, data
}

func CreateErrorOccurredEvent(source string, errorMsg string, details interface{}) (EventType, map[string]interface{}) {
	data := map[string]interface{}{
		"source":    source,
		"error":     errorMsg,
		"details":   details,
		"timestamp": time.Now(),
	}
	return EventTypeErrorOccurred, data
}

func CreateStatusChangedEvent(entityType string, entityID string, oldStatus string, newStatus string) (EventType, map[string]interface{}) {
	data := map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityID,
		"old_status":  oldStatus,
		"new_status":  newStatus,
		"timestamp":   time.Now(),
	}
	return EventTypeStatusChanged, data
}

func EventToJSON(event Event) (string, error) {
	bytes, err := json.Marshal(event)
	if err != nil {
		return "", fmt.Errorf("failed to marshal event to JSON: %w", err)
	}
	return string(bytes), nil
}

func JSONToEvent(jsonStr string) (Event, error) {
	var event Event
	err := json.Unmarshal([]byte(jsonStr), &event)
	if err != nil {
		return Event{}, fmt.Errorf("failed to unmarshal JSON to event: %w", err)
	}
	return event, nil
}
