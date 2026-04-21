package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type TraceSpan struct {
	ID            string                 `json:"id"`
	TraceID       string                 `json:"trace_id"`
	ParentID      string                 `json:"parent_id,omitempty"`
	Name          string                 `json:"name"`
	Kind          string                 `json:"kind"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	Duration      time.Duration          `json:"duration"`
	Attributes    map[string]interface{} `json:"attributes"`
	Events        []TraceEvent           `json:"events,omitempty"`
	Status        string                 `json:"status"`
	StatusMessage string                 `json:"status_message,omitempty"`
	Resource      map[string]interface{} `json:"resource,omitempty"`
	AgentName     string                 `json:"agent_name,omitempty"`
	ToolName      string                 `json:"tool_name,omitempty"`
	PipelineID    string                 `json:"pipeline_id,omitempty"`
}

type TraceEvent struct {
	Name       string                 `json:"name"`
	Timestamp  time.Time              `json:"timestamp"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

type TraceQuery struct {
	TraceID    string
	SpanName   string
	AgentName  string
	ToolName   string
	PipelineID string
	Status     string
	StartTime  time.Time
	EndTime    time.Time
	Limit      int
	Offset     int
}

type TraceStore struct {
	db     *sql.DB
	dbPath string
	mu     sync.RWMutex
}

func NewTraceStore(dbPath string) (*TraceStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &TraceStore{
		db:     db,
		dbPath: dbPath,
	}

	if err := store.initTables(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return store, nil
}

func (ts *TraceStore) initTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS traces (
			id TEXT PRIMARY KEY,
			trace_id TEXT NOT NULL,
			parent_id TEXT,
			name TEXT NOT NULL,
			kind TEXT NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			duration INTEGER NOT NULL,
			attributes TEXT,
			events TEXT,
			status TEXT NOT NULL,
			status_message TEXT,
			resource TEXT,
			agent_name TEXT,
			tool_name TEXT,
			pipeline_id TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_trace_id ON traces(trace_id)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_parent_id ON traces(parent_id)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_name ON traces(name)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_agent_name ON traces(agent_name)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_tool_name ON traces(tool_name)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_pipeline_id ON traces(pipeline_id)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_status ON traces(status)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_start_time ON traces(start_time)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_end_time ON traces(end_time)`,
		`CREATE INDEX IF NOT EXISTS idx_traces_created_at ON traces(created_at)`,
	}

	for _, query := range queries {
		if _, err := ts.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}

	return nil
}

func (ts *TraceStore) StoreSpan(ctx context.Context, span *TraceSpan) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	attributesJSON, err := json.Marshal(span.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	eventsJSON, err := json.Marshal(span.Events)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	resourceJSON, err := json.Marshal(span.Resource)
	if err != nil {
		return fmt.Errorf("failed to marshal resource: %w", err)
	}

	query := `INSERT INTO traces (
		id, trace_id, parent_id, name, kind, start_time, end_time, duration,
		attributes, events, status, status_message, resource,
		agent_name, tool_name, pipeline_id
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = ts.db.ExecContext(ctx, query,
		span.ID,
		span.TraceID,
		span.ParentID,
		span.Name,
		span.Kind,
		span.StartTime,
		span.EndTime,
		span.Duration.Milliseconds(),
		string(attributesJSON),
		string(eventsJSON),
		span.Status,
		span.StatusMessage,
		string(resourceJSON),
		span.AgentName,
		span.ToolName,
		span.PipelineID,
	)

	if err != nil {
		return fmt.Errorf("failed to insert span: %w", err)
	}

	return nil
}

func (ts *TraceStore) GetSpan(ctx context.Context, id string) (*TraceSpan, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	query := `SELECT 
		id, trace_id, parent_id, name, kind, start_time, end_time, duration,
		attributes, events, status, status_message, resource,
		agent_name, tool_name, pipeline_id
	FROM traces WHERE id = ?`

	row := ts.db.QueryRowContext(ctx, query, id)

	var span TraceSpan
	var attributesJSON, eventsJSON, resourceJSON string
	var durationMs int64

	err := row.Scan(
		&span.ID,
		&span.TraceID,
		&span.ParentID,
		&span.Name,
		&span.Kind,
		&span.StartTime,
		&span.EndTime,
		&durationMs,
		&attributesJSON,
		&eventsJSON,
		&span.Status,
		&span.StatusMessage,
		&resourceJSON,
		&span.AgentName,
		&span.ToolName,
		&span.PipelineID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan span: %w", err)
	}

	span.Duration = time.Duration(durationMs) * time.Millisecond

	if err := json.Unmarshal([]byte(attributesJSON), &span.Attributes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
	}

	if err := json.Unmarshal([]byte(eventsJSON), &span.Events); err != nil {
		return nil, fmt.Errorf("failed to unmarshal events: %w", err)
	}

	if err := json.Unmarshal([]byte(resourceJSON), &span.Resource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	return &span, nil
}

func (ts *TraceStore) QuerySpans(ctx context.Context, queryParams TraceQuery) ([]*TraceSpan, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	baseQuery := `SELECT 
		id, trace_id, parent_id, name, kind, start_time, end_time, duration,
		attributes, events, status, status_message, resource,
		agent_name, tool_name, pipeline_id
	FROM traces WHERE 1=1`

	args := []interface{}{}
	conditions := []string{}

	if queryParams.TraceID != "" {
		conditions = append(conditions, "trace_id = ?")
		args = append(args, queryParams.TraceID)
	}

	if queryParams.SpanName != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+queryParams.SpanName+"%")
	}

	if queryParams.AgentName != "" {
		conditions = append(conditions, "agent_name = ?")
		args = append(args, queryParams.AgentName)
	}

	if queryParams.ToolName != "" {
		conditions = append(conditions, "tool_name = ?")
		args = append(args, queryParams.ToolName)
	}

	if queryParams.PipelineID != "" {
		conditions = append(conditions, "pipeline_id = ?")
		args = append(args, queryParams.PipelineID)
	}

	if queryParams.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, queryParams.Status)
	}

	if !queryParams.StartTime.IsZero() {
		conditions = append(conditions, "start_time >= ?")
		args = append(args, queryParams.StartTime)
	}

	if !queryParams.EndTime.IsZero() {
		conditions = append(conditions, "end_time <= ?")
		args = append(args, queryParams.EndTime)
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + joinConditions(conditions, " AND ")
	}

	baseQuery += " ORDER BY start_time DESC"

	if queryParams.Limit > 0 {
		baseQuery += " LIMIT ?"
		args = append(args, queryParams.Limit)
	}

	if queryParams.Offset > 0 {
		baseQuery += " OFFSET ?"
		args = append(args, queryParams.Offset)
	}

	rows, err := ts.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query spans: %w", err)
	}
	defer rows.Close()

	spans := []*TraceSpan{}
	for rows.Next() {
		var span TraceSpan
		var attributesJSON, eventsJSON, resourceJSON string
		var durationMs int64

		err := rows.Scan(
			&span.ID,
			&span.TraceID,
			&span.ParentID,
			&span.Name,
			&span.Kind,
			&span.StartTime,
			&span.EndTime,
			&durationMs,
			&attributesJSON,
			&eventsJSON,
			&span.Status,
			&span.StatusMessage,
			&resourceJSON,
			&span.AgentName,
			&span.ToolName,
			&span.PipelineID,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan span: %w", err)
		}

		span.Duration = time.Duration(durationMs) * time.Millisecond

		if err := json.Unmarshal([]byte(attributesJSON), &span.Attributes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
		}

		if err := json.Unmarshal([]byte(eventsJSON), &span.Events); err != nil {
			return nil, fmt.Errorf("failed to unmarshal events: %w", err)
		}

		if err := json.Unmarshal([]byte(resourceJSON), &span.Resource); err != nil {
			return nil, fmt.Errorf("failed to unmarshal resource: %w", err)
		}

		spans = append(spans, &span)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return spans, nil
}

func (ts *TraceStore) GetTraceTree(ctx context.Context, traceID string) ([]*TraceSpan, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	query := `SELECT 
		id, trace_id, parent_id, name, kind, start_time, end_time, duration,
		attributes, events, status, status_message, resource,
		agent_name, tool_name, pipeline_id
	FROM traces WHERE trace_id = ? ORDER BY start_time`

	rows, err := ts.db.QueryContext(ctx, query, traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query trace tree: %w", err)
	}
	defer rows.Close()

	spans := []*TraceSpan{}
	for rows.Next() {
		var span TraceSpan
		var attributesJSON, eventsJSON, resourceJSON string
		var durationMs int64

		err := rows.Scan(
			&span.ID,
			&span.TraceID,
			&span.ParentID,
			&span.Name,
			&span.Kind,
			&span.StartTime,
			&span.EndTime,
			&durationMs,
			&attributesJSON,
			&eventsJSON,
			&span.Status,
			&span.StatusMessage,
			&resourceJSON,
			&span.AgentName,
			&span.ToolName,
			&span.PipelineID,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan span: %w", err)
		}

		span.Duration = time.Duration(durationMs) * time.Millisecond

		if err := json.Unmarshal([]byte(attributesJSON), &span.Attributes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
		}

		if err := json.Unmarshal([]byte(eventsJSON), &span.Events); err != nil {
			return nil, fmt.Errorf("failed to unmarshal events: %w", err)
		}

		if err := json.Unmarshal([]byte(resourceJSON), &span.Resource); err != nil {
			return nil, fmt.Errorf("failed to unmarshal resource: %w", err)
		}

		spans = append(spans, &span)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return spans, nil
}

func (ts *TraceStore) GetAgentSpans(ctx context.Context, agentName string, limit int) ([]*TraceSpan, error) {
	return ts.QuerySpans(ctx, TraceQuery{
		AgentName: agentName,
		Limit:     limit,
	})
}

func (ts *TraceStore) GetPipelineSpans(ctx context.Context, pipelineID string) ([]*TraceSpan, error) {
	return ts.QuerySpans(ctx, TraceQuery{
		PipelineID: pipelineID,
	})
}

func (ts *TraceStore) GetRecentSpans(ctx context.Context, limit int) ([]*TraceSpan, error) {
	return ts.QuerySpans(ctx, TraceQuery{
		Limit: limit,
	})
}

func (ts *TraceStore) GetStatistics(ctx context.Context, startTime, endTime time.Time) (map[string]interface{}, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	stats := make(map[string]interface{})

	queries := map[string]string{
		"total_spans":     `SELECT COUNT(*) FROM traces WHERE start_time >= ? AND end_time <= ?`,
		"total_agents":    `SELECT COUNT(DISTINCT agent_name) FROM traces WHERE start_time >= ? AND end_time <= ? AND agent_name IS NOT NULL`,
		"total_tools":     `SELECT COUNT(DISTINCT tool_name) FROM traces WHERE start_time >= ? AND end_time <= ? AND tool_name IS NOT NULL`,
		"total_pipelines": `SELECT COUNT(DISTINCT pipeline_id) FROM traces WHERE start_time >= ? AND end_time <= ? AND pipeline_id IS NOT NULL`,
		"avg_duration":    `SELECT AVG(duration) FROM traces WHERE start_time >= ? AND end_time <= ?`,
		"success_rate": `SELECT 
			(COUNT(CASE WHEN status = 'OK' THEN 1 END) * 100.0 / COUNT(*)) as success_rate 
			FROM traces WHERE start_time >= ? AND end_time <= ?`,
	}

	for name, query := range queries {
		var value interface{}
		err := ts.db.QueryRowContext(ctx, query, startTime, endTime).Scan(&value)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to query %s: %w", name, err)
		}
		stats[name] = value
	}

	agentStatsQuery := `SELECT 
		agent_name, 
		COUNT(*) as count,
		AVG(duration) as avg_duration,
		(COUNT(CASE WHEN status = 'OK' THEN 1 END) * 100.0 / COUNT(*)) as success_rate
		FROM traces 
		WHERE start_time >= ? AND end_time <= ? AND agent_name IS NOT NULL
		GROUP BY agent_name
		ORDER BY count DESC`

	rows, err := ts.db.QueryContext(ctx, agentStatsQuery, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query agent stats: %w", err)
	}
	defer rows.Close()

	agentStats := []map[string]interface{}{}
	for rows.Next() {
		var agentName string
		var count int
		var avgDuration, successRate float64

		err := rows.Scan(&agentName, &count, &avgDuration, &successRate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agent stats: %w", err)
		}

		agentStats = append(agentStats, map[string]interface{}{
			"agent_name":   agentName,
			"count":        count,
			"avg_duration": avgDuration,
			"success_rate": successRate,
		})
	}

	stats["agent_statistics"] = agentStats

	return stats, nil
}

func (ts *TraceStore) CleanupOldData(ctx context.Context, retentionDays int) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	query := `DELETE FROM traces WHERE end_time < ?`

	result, err := ts.db.ExecContext(ctx, query, cutoffTime)
	if err != nil {
		return fmt.Errorf("failed to cleanup old data: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Cleaned up %d old trace records\n", rowsAffected)

	return nil
}

func (ts *TraceStore) Close() error {
	return ts.db.Close()
}

func joinConditions(conditions []string, sep string) string {
	if len(conditions) == 0 {
		return ""
	}
	if len(conditions) == 1 {
		return conditions[0]
	}

	result := conditions[0]
	for i := 1; i < len(conditions); i++ {
		result += sep + conditions[i]
	}
	return result
}
