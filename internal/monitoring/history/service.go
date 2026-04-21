package history

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type ExecutionRecord struct {
	ID            string                 `json:"id"`
	TraceID       string                 `json:"trace_id"`
	PipelineID    string                 `json:"pipeline_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description,omitempty"`
	Status        string                 `json:"status"` // "running", "completed", "failed", "cancelled"
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time,omitempty"`
	Duration      time.Duration          `json:"duration,omitempty"`
	AgentCount    int                    `json:"agent_count"`
	ToolCallCount int                    `json:"tool_call_count"`
	TokenUsage    int                    `json:"token_usage"`
	CostEstimate  float64                `json:"cost_estimate,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Tags          []string               `json:"tags,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type ExecutionQuery struct {
	IDs        []string
	Statuses   []string
	PipelineID string
	SearchText string
	Tags       []string
	StartTime  time.Time
	EndTime    time.Time
	Limit      int
	Offset     int
	SortBy     string // "start_time", "end_time", "duration", "agent_count", "token_usage"
	SortOrder  string // "asc", "desc"
}

type ReplayState struct {
	RecordID         string                 `json:"record_id"`
	TraceID          string                 `json:"trace_id"`
	Status           string                 `json:"status"` // "initializing", "playing", "paused", "completed", "error"
	CurrentTime      time.Time              `json:"current_time"`
	StartTime        time.Time              `json:"start_time"`
	EndTime          time.Time              `json:"end_time"`
	PlaybackSpeed    float64                `json:"playback_speed"` // 0.5x, 1x, 2x, 4x, etc.
	CurrentSpanIndex int                    `json:"current_span_index"`
	TotalSpans       int                    `json:"total_spans"`
	Progress         float64                `json:"progress"` // 0-100
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

type HistoryService struct {
	storage      *HistoryStorage
	replayStates map[string]*ReplayState
	mu           sync.RWMutex
	traceStore   TraceStore
}

type TraceStore interface {
	GetTraceTree(ctx context.Context, traceID string) ([]interface{}, error)
	QuerySpans(ctx context.Context, query interface{}) ([]interface{}, error)
}

func NewHistoryService(traceStore TraceStore) (*HistoryService, error) {
	// 使用默认数据库路径
	storage, err := NewHistoryStorage("history.db")
	if err != nil {
		return nil, fmt.Errorf("failed to create history storage: %w", err)
	}

	return &HistoryService{
		storage:      storage,
		replayStates: make(map[string]*ReplayState),
		traceStore:   traceStore,
	}, nil
}

func NewHistoryServiceWithStorage(storage *HistoryStorage, traceStore TraceStore) *HistoryService {
	return &HistoryService{
		storage:      storage,
		replayStates: make(map[string]*ReplayState),
		traceStore:   traceStore,
	}
}

func (hs *HistoryService) CreateRecord(ctx context.Context, record *ExecutionRecord) (*ExecutionRecord, error) {
	if record.ID == "" {
		record.ID = uuid.New().String()
	}
	if record.TraceID == "" {
		// 尝试从上下文中获取Trace ID
		if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
			record.TraceID = span.SpanContext().TraceID().String()
		}
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now()
	}

	if err := hs.storage.CreateRecord(ctx, record); err != nil {
		return nil, fmt.Errorf("failed to create record in storage: %w", err)
	}

	return record, nil
}

func (hs *HistoryService) UpdateRecord(ctx context.Context, id string, updates map[string]interface{}) (*ExecutionRecord, error) {
	if err := hs.storage.UpdateRecord(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update record in storage: %w", err)
	}

	// 重新获取更新后的记录
	return hs.storage.GetRecord(ctx, id)
}

func (hs *HistoryService) GetRecord(ctx context.Context, id string) (*ExecutionRecord, error) {
	record, err := hs.storage.GetRecord(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get record from storage: %w", err)
	}
	if record == nil {
		return nil, fmt.Errorf("record not found: %s", id)
	}
	return record, nil
}

func (hs *HistoryService) QueryRecords(ctx context.Context, query ExecutionQuery) ([]*ExecutionRecord, int, error) {
	return hs.storage.QueryRecords(ctx, query)
}

func (hs *HistoryService) DeleteRecord(ctx context.Context, id string) error {
	// 先删除相关的回放状态
	hs.mu.Lock()
	for replayID, state := range hs.replayStates {
		if state.RecordID == id {
			delete(hs.replayStates, replayID)
		}
	}
	hs.mu.Unlock()

	// 然后删除存储中的记录
	return hs.storage.DeleteRecord(ctx, id)
}

func (hs *HistoryService) GetStatistics(ctx context.Context, startTime, endTime time.Time) (map[string]interface{}, error) {
	return hs.storage.GetStatistics(ctx, startTime, endTime)
}

func (hs *HistoryService) StartReplay(ctx context.Context, recordID string, speed float64) (*ReplayState, error) {
	// 获取记录
	record, err := hs.GetRecord(ctx, recordID)
	if err != nil {
		return nil, fmt.Errorf("record not found: %s", recordID)
	}

	if record.TraceID == "" {
		return nil, fmt.Errorf("record has no trace ID")
	}

	if hs.traceStore == nil {
		return nil, fmt.Errorf("trace store not available")
	}

	replayID := uuid.New().String()
	state := &ReplayState{
		RecordID:         recordID,
		TraceID:          record.TraceID,
		Status:           "initializing",
		StartTime:        record.StartTime,
		EndTime:          record.EndTime,
		PlaybackSpeed:    speed,
		CurrentSpanIndex: 0,
		Progress:         0,
		Metadata: map[string]interface{}{
			"record_name":   record.Name,
			"record_status": record.Status,
		},
	}

	hs.mu.Lock()
	hs.replayStates[replayID] = state
	hs.mu.Unlock()

	go hs.initializeReplay(ctx, replayID, state)

	return state, nil
}

func (hs *HistoryService) initializeReplay(ctx context.Context, replayID string, state *ReplayState) {
	hs.mu.Lock()
	state.Status = "playing"
	hs.mu.Unlock()

	spans, err := hs.traceStore.GetTraceTree(ctx, state.TraceID)
	if err != nil {
		hs.mu.Lock()
		state.Status = "error"
		state.Metadata["error"] = err.Error()
		hs.mu.Unlock()
		return
	}

	hs.mu.Lock()
	state.TotalSpans = len(spans)
	state.CurrentTime = state.StartTime
	hs.mu.Unlock()

	hs.broadcastReplayUpdate(replayID, state)
}

func (hs *HistoryService) ControlReplay(ctx context.Context, replayID string, action string, params map[string]interface{}) (*ReplayState, error) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	state, exists := hs.replayStates[replayID]
	if !exists {
		return nil, fmt.Errorf("replay not found: %s", replayID)
	}

	switch action {
	case "pause":
		state.Status = "paused"
	case "resume":
		if state.Status == "paused" {
			state.Status = "playing"
		}
	case "stop":
		state.Status = "completed"
		state.CurrentSpanIndex = state.TotalSpans
		state.Progress = 100
	case "seek":
		if progress, ok := params["progress"].(float64); ok {
			if progress < 0 {
				progress = 0
			}
			if progress > 100 {
				progress = 100
			}
			state.Progress = progress
			if state.TotalSpans > 0 {
				state.CurrentSpanIndex = int(float64(state.TotalSpans) * progress / 100)
			}
			totalDuration := state.EndTime.Sub(state.StartTime)
			state.CurrentTime = state.StartTime.Add(time.Duration(float64(totalDuration) * progress / 100))
		}
	case "speed":
		if speed, ok := params["speed"].(float64); ok && speed > 0 {
			state.PlaybackSpeed = speed
		}
	}

	hs.broadcastReplayUpdate(replayID, state)

	return state, nil
}

func (hs *HistoryService) GetReplayState(ctx context.Context, replayID string) (*ReplayState, error) {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	state, exists := hs.replayStates[replayID]
	if !exists {
		return nil, fmt.Errorf("replay not found: %s", replayID)
	}

	return state, nil
}

func (hs *HistoryService) GetReplayData(ctx context.Context, replayID string, fromIndex, limit int) ([]interface{}, error) {
	hs.mu.RLock()
	state, exists := hs.replayStates[replayID]
	hs.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("replay not found: %s", replayID)
	}

	spans, err := hs.traceStore.GetTraceTree(ctx, state.TraceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trace tree: %w", err)
	}

	if fromIndex < 0 {
		fromIndex = 0
	}
	if fromIndex >= len(spans) {
		return []interface{}{}, nil
	}

	endIndex := fromIndex + limit
	if endIndex > len(spans) {
		endIndex = len(spans)
	}

	return spans[fromIndex:endIndex], nil
}

func (hs *HistoryService) broadcastReplayUpdate(replayID string, state *ReplayState) {
	updateData := map[string]interface{}{
		"replay_id": replayID,
		"state":     state,
		"timestamp": time.Now(),
	}

	jsonData, _ := json.Marshal(updateData)
	fmt.Printf("Replay update: %s\n", string(jsonData))
}

func contains(str, substr string) bool {
	if len(str) < len(substr) {
		return false
	}
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
