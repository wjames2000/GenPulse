package history

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WailsBridge struct {
	service *HistoryService
	ctx     context.Context
}

func NewWailsBridge(service *HistoryService) *WailsBridge {
	return &WailsBridge{
		service: service,
	}
}

func (wb *WailsBridge) SetContext(ctx context.Context) {
	wb.ctx = ctx
}

func (wb *WailsBridge) CreateRecord(name, description, pipelineID string, metadata map[string]interface{}, tags []string) (*ExecutionRecord, error) {
	record := &ExecutionRecord{
		Name:        name,
		Description: description,
		PipelineID:  pipelineID,
		Status:      "running",
		StartTime:   time.Now(),
		Metadata:    metadata,
		Tags:        tags,
	}

	return wb.service.CreateRecord(wb.ctx, record)
}

func (wb *WailsBridge) UpdateRecord(id string, updates map[string]interface{}) (*ExecutionRecord, error) {
	return wb.service.UpdateRecord(wb.ctx, id, updates)
}

func (wb *WailsBridge) GetRecord(id string) (*ExecutionRecord, error) {
	return wb.service.GetRecord(wb.ctx, id)
}

func (wb *WailsBridge) QueryRecords(params map[string]interface{}) ([]*ExecutionRecord, int, error) {
	query := ExecutionQuery{
		Limit:     50,
		Offset:    0,
		SortBy:    "start_time",
		SortOrder: "desc",
	}

	if limit, ok := params["limit"].(float64); ok {
		query.Limit = int(limit)
	}
	if offset, ok := params["offset"].(float64); ok {
		query.Offset = int(offset)
	}
	if sortBy, ok := params["sortBy"].(string); ok {
		query.SortBy = sortBy
	}
	if sortOrder, ok := params["sortOrder"].(string); ok {
		query.SortOrder = sortOrder
	}
	if searchText, ok := params["searchText"].(string); ok {
		query.SearchText = searchText
	}
	if pipelineID, ok := params["pipelineID"].(string); ok {
		query.PipelineID = pipelineID
	}

	if statuses, ok := params["statuses"].([]interface{}); ok {
		for _, s := range statuses {
			if status, ok := s.(string); ok {
				query.Statuses = append(query.Statuses, status)
			}
		}
	}

	if tags, ok := params["tags"].([]interface{}); ok {
		for _, t := range tags {
			if tag, ok := t.(string); ok {
				query.Tags = append(query.Tags, tag)
			}
		}
	}

	if startTimeStr, ok := params["startTime"].(string); ok {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			query.StartTime = startTime
		}
	}

	if endTimeStr, ok := params["endTime"].(string); ok {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			query.EndTime = endTime
		}
	}

	return wb.service.QueryRecords(wb.ctx, query)
}

func (wb *WailsBridge) DeleteRecord(id string) error {
	return wb.service.DeleteRecord(wb.ctx, id)
}

func (wb *WailsBridge) GetStatistics(startTimeStr, endTimeStr string) (map[string]interface{}, error) {
	var startTime, endTime time.Time
	var err error

	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			startTime = time.Now().AddDate(0, -1, 0)
		}
	} else {
		startTime = time.Now().AddDate(0, -1, 0)
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			endTime = time.Now()
		}
	} else {
		endTime = time.Now()
	}

	return wb.service.GetStatistics(wb.ctx, startTime, endTime)
}

func (wb *WailsBridge) StartReplay(recordID string, speed float64) (*ReplayState, error) {
	if speed <= 0 {
		speed = 1.0
	}

	state, err := wb.service.StartReplay(wb.ctx, recordID, speed)
	if err != nil {
		return nil, err
	}

	wb.emitEvent("replay:started", map[string]interface{}{
		"replay_id": recordID,
		"state":     state,
	})

	return state, nil
}

func (wb *WailsBridge) ControlReplay(replayID, action string, params map[string]interface{}) (*ReplayState, error) {
	state, err := wb.service.ControlReplay(wb.ctx, replayID, action, params)
	if err != nil {
		return nil, err
	}

	wb.emitEvent("replay:updated", map[string]interface{}{
		"replay_id": replayID,
		"action":    action,
		"state":     state,
	})

	return state, nil
}

func (wb *WailsBridge) GetReplayState(replayID string) (*ReplayState, error) {
	return wb.service.GetReplayState(wb.ctx, replayID)
}

func (wb *WailsBridge) GetReplayData(replayID string, fromIndex, limit int) ([]interface{}, error) {
	return wb.service.GetReplayData(wb.ctx, replayID, fromIndex, limit)
}

func (wb *WailsBridge) emitEvent(eventName string, data interface{}) {
	if wb.ctx == nil {
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Failed to marshal event data: %v\n", err)
		return
	}

	runtime.EventsEmit(wb.ctx, eventName, string(jsonData))
}

func (wb *WailsBridge) SubscribeToReplayEvents(replayID string) error {
	if wb.ctx == nil {
		return fmt.Errorf("context not set")
	}

	go wb.monitorReplay(replayID)

	return nil
}

func (wb *WailsBridge) monitorReplay(replayID string) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			state, err := wb.service.GetReplayState(wb.ctx, replayID)
			if err != nil {
				return
			}

			if state.Status == "completed" || state.Status == "error" {
				wb.emitEvent("replay:ended", map[string]interface{}{
					"replay_id": replayID,
					"state":     state,
				})
				return
			}

			if state.Status == "playing" {
				wb.emitEvent("replay:progress", map[string]interface{}{
					"replay_id": replayID,
					"state":     state,
				})
			}
		}
	}
}
