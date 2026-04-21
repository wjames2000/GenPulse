package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WailsEventBridge struct {
	eventBus   *EventBus
	appContext context.Context
	mu         sync.RWMutex
	subscribed map[string]bool
}

func NewWailsEventBridge(eventBus *EventBus) *WailsEventBridge {
	return &WailsEventBridge{
		eventBus:   eventBus,
		subscribed: make(map[string]bool),
	}
}

func (web *WailsEventBridge) SetAppContext(ctx context.Context) {
	web.mu.Lock()
	defer web.mu.Unlock()
	web.appContext = ctx
}

func (web *WailsEventBridge) StartForwarding() error {
	web.mu.Lock()
	defer web.mu.Unlock()

	if web.appContext == nil {
		return fmt.Errorf("app context not set")
	}

	eventTypes := []EventType{
		EventTypeAgentStarted,
		EventTypeAgentCompleted,
		EventTypeToolInvoked,
		EventTypeToolCompleted,
		EventTypePipelineStarted,
		EventTypePipelineCompleted,
		EventTypeMetricUpdated,
		EventTypeTraceRecorded,
		EventTypeSkillGenerated,
		EventTypeErrorOccurred,
		EventTypeStatusChanged,
	}

	for _, eventType := range eventTypes {
		if web.subscribed[string(eventType)] {
			continue
		}

		handler := func(ctx context.Context, event Event) error {
			return web.forwardToFrontend(event)
		}

		web.eventBus.Subscribe(eventType, handler)
		web.subscribed[string(eventType)] = true
	}

	return nil
}

func (web *WailsEventBridge) forwardToFrontend(event Event) error {
	web.mu.RLock()
	ctx := web.appContext
	web.mu.RUnlock()

	if ctx == nil {
		return fmt.Errorf("app context not available")
	}

	eventJSON, err := EventToJSON(event)
	if err != nil {
		return fmt.Errorf("failed to convert event to JSON: %w", err)
	}

	runtime.EventsEmit(ctx, "monitoring_event", eventJSON)
	return nil
}

func (web *WailsEventBridge) StopForwarding() {
	web.mu.Lock()
	defer web.mu.Unlock()

	for eventType := range web.subscribed {
		delete(web.subscribed, eventType)
	}
}

func (web *WailsEventBridge) SendCustomEventToFrontend(eventName string, data interface{}) error {
	web.mu.RLock()
	ctx := web.appContext
	web.mu.RUnlock()

	if ctx == nil {
		return fmt.Errorf("app context not available")
	}

	runtime.EventsEmit(ctx, eventName, data)
	return nil
}

func (web *WailsEventBridge) GetRecentEventsForFrontend(limit int) ([]map[string]interface{}, error) {
	events := web.eventBus.GetRecentEvents(limit)
	result := make([]map[string]interface{}, len(events))

	for i, event := range events {
		eventJSON, err := EventToJSON(event)
		if err != nil {
			return nil, fmt.Errorf("failed to convert event %d to JSON: %w", i, err)
		}

		var eventMap map[string]interface{}
		if err := json.Unmarshal([]byte(eventJSON), &eventMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event JSON: %w", err)
		}

		result[i] = eventMap
	}

	return result, nil
}

func (web *WailsEventBridge) GetEventsByTypeForFrontend(eventTypeStr string, limit int) ([]map[string]interface{}, error) {
	eventType := EventType(eventTypeStr)
	events := web.eventBus.GetEventsByType(eventType, limit)
	result := make([]map[string]interface{}, len(events))

	for i, event := range events {
		eventJSON, err := EventToJSON(event)
		if err != nil {
			return nil, fmt.Errorf("failed to convert event %d to JSON: %w", i, err)
		}

		var eventMap map[string]interface{}
		if err := json.Unmarshal([]byte(eventJSON), &eventMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event JSON: %w", err)
		}

		result[i] = eventMap
	}

	return result, nil
}

func (web *WailsEventBridge) ClearEventBuffer() {
	web.eventBus.ClearBuffer()
}

type FrontendEventReceiver struct {
	eventBus *EventBus
}

func NewFrontendEventReceiver(eventBus *EventBus) *FrontendEventReceiver {
	return &FrontendEventReceiver{
		eventBus: eventBus,
	}
}

func (fer *FrontendEventReceiver) HandleFrontendEvent(eventJSON string) error {
	event, err := JSONToEvent(eventJSON)
	if err != nil {
		return fmt.Errorf("failed to parse frontend event: %w", err)
	}

	ctx := context.Background()
	return fer.eventBus.Publish(ctx, event.Type, "frontend", event.Data)
}

func (fer *FrontendEventReceiver) SubscribeToFrontendEvents(eventType EventType, handler EventHandler) {
	fer.eventBus.Subscribe(eventType, handler)
}
