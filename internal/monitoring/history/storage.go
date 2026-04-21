package history

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type HistoryStorage struct {
	db     *sql.DB
	dbPath string
	mu     sync.RWMutex
}

func NewHistoryStorage(dbPath string) (*HistoryStorage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &HistoryStorage{
		db:     db,
		dbPath: dbPath,
	}

	if err := storage.initTables(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return storage, nil
}

func (hs *HistoryStorage) initTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS execution_records (
			id TEXT PRIMARY KEY,
			trace_id TEXT NOT NULL,
			pipeline_id TEXT,
			name TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME,
			duration INTEGER,
			agent_count INTEGER NOT NULL DEFAULT 0,
			tool_call_count INTEGER NOT NULL DEFAULT 0,
			token_usage INTEGER NOT NULL DEFAULT 0,
			cost_estimate REAL,
			metadata TEXT,
			tags TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_records_trace_id ON execution_records(trace_id)`,
		`CREATE INDEX IF NOT EXISTS idx_records_pipeline_id ON execution_records(pipeline_id)`,
		`CREATE INDEX IF NOT EXISTS idx_records_status ON execution_records(status)`,
		`CREATE INDEX IF NOT EXISTS idx_records_start_time ON execution_records(start_time)`,
		`CREATE INDEX IF NOT EXISTS idx_records_end_time ON execution_records(end_time)`,
		`CREATE INDEX IF NOT EXISTS idx_records_created_at ON execution_records(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_records_name ON execution_records(name)`,
		`CREATE INDEX IF NOT EXISTS idx_records_agent_count ON execution_records(agent_count)`,
		`CREATE INDEX IF NOT EXISTS idx_records_token_usage ON execution_records(token_usage)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS records_fts USING fts5(
			id UNINDEXED,
			name,
			description,
			tags,
			content='execution_records',
			content_rowid='rowid'
		)`,
		`CREATE TRIGGER IF NOT EXISTS records_ai AFTER INSERT ON execution_records BEGIN
			INSERT INTO records_fts(rowid, name, description, tags) 
			VALUES (new.rowid, new.name, new.description, new.tags);
		END`,
		`CREATE TRIGGER IF NOT EXISTS records_ad AFTER DELETE ON execution_records BEGIN
			INSERT INTO records_fts(records_fts, rowid, name, description, tags) 
			VALUES('delete', old.rowid, old.name, old.description, old.tags);
		END`,
		`CREATE TRIGGER IF NOT EXISTS records_au AFTER UPDATE ON execution_records BEGIN
			INSERT INTO records_fts(records_fts, rowid, name, description, tags) 
			VALUES('delete', old.rowid, old.name, old.description, old.tags);
			INSERT INTO records_fts(rowid, name, description, tags) 
			VALUES (new.rowid, new.name, new.description, new.tags);
		END`,
	}

	for _, query := range queries {
		if _, err := hs.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}

	return nil
}

func (hs *HistoryStorage) CreateRecord(ctx context.Context, record *ExecutionRecord) error {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	if record.ID == "" {
		record.ID = uuid.New().String()
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now()
	}
	record.UpdatedAt = time.Now()

	metadataJSON, err := json.Marshal(record.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	tagsJSON, err := json.Marshal(record.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `INSERT INTO execution_records (
		id, trace_id, pipeline_id, name, description, status,
		start_time, end_time, duration, agent_count, tool_call_count,
		token_usage, cost_estimate, metadata, tags, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = hs.db.ExecContext(ctx, query,
		record.ID,
		record.TraceID,
		record.PipelineID,
		record.Name,
		record.Description,
		record.Status,
		record.StartTime,
		record.EndTime,
		record.Duration.Milliseconds(),
		record.AgentCount,
		record.ToolCallCount,
		record.TokenUsage,
		record.CostEstimate,
		string(metadataJSON),
		string(tagsJSON),
		record.CreatedAt,
		record.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	return nil
}

func (hs *HistoryStorage) UpdateRecord(ctx context.Context, id string, updates map[string]interface{}) error {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	// 先获取现有记录
	record, err := hs.getRecordByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get record: %w", err)
	}

	// 应用更新
	if status, ok := updates["status"].(string); ok {
		record.Status = status
	}
	if endTime, ok := updates["end_time"].(time.Time); ok {
		record.EndTime = endTime
		if !record.StartTime.IsZero() {
			record.Duration = endTime.Sub(record.StartTime)
		}
	}
	if agentCount, ok := updates["agent_count"].(int); ok {
		record.AgentCount = agentCount
	}
	if toolCallCount, ok := updates["tool_call_count"].(int); ok {
		record.ToolCallCount = toolCallCount
	}
	if tokenUsage, ok := updates["token_usage"].(int); ok {
		record.TokenUsage = tokenUsage
	}
	if costEstimate, ok := updates["cost_estimate"].(float64); ok {
		record.CostEstimate = costEstimate
	}
	if metadata, ok := updates["metadata"].(map[string]interface{}); ok {
		if record.Metadata == nil {
			record.Metadata = make(map[string]interface{})
		}
		for k, v := range metadata {
			record.Metadata[k] = v
		}
	}
	if tags, ok := updates["tags"].([]string); ok {
		record.Tags = tags
	}

	record.UpdatedAt = time.Now()

	// 更新数据库
	metadataJSON, err := json.Marshal(record.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	tagsJSON, err := json.Marshal(record.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `UPDATE execution_records SET
		status = ?,
		end_time = ?,
		duration = ?,
		agent_count = ?,
		tool_call_count = ?,
		token_usage = ?,
		cost_estimate = ?,
		metadata = ?,
		tags = ?,
		updated_at = ?
	WHERE id = ?`

	_, err = hs.db.ExecContext(ctx, query,
		record.Status,
		record.EndTime,
		record.Duration.Milliseconds(),
		record.AgentCount,
		record.ToolCallCount,
		record.TokenUsage,
		record.CostEstimate,
		string(metadataJSON),
		string(tagsJSON),
		record.UpdatedAt,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	return nil
}

func (hs *HistoryStorage) GetRecord(ctx context.Context, id string) (*ExecutionRecord, error) {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	return hs.getRecordByID(ctx, id)
}

func (hs *HistoryStorage) getRecordByID(ctx context.Context, id string) (*ExecutionRecord, error) {
	query := `SELECT 
		id, trace_id, pipeline_id, name, description, status,
		start_time, end_time, duration, agent_count, tool_call_count,
		token_usage, cost_estimate, metadata, tags, created_at, updated_at
	FROM execution_records WHERE id = ?`

	row := hs.db.QueryRowContext(ctx, query, id)

	var record ExecutionRecord
	var metadataJSON, tagsJSON string
	var durationMs sql.NullInt64
	var endTime, created_at, updated_at sql.NullTime
	var costEstimate sql.NullFloat64

	err := row.Scan(
		&record.ID,
		&record.TraceID,
		&record.PipelineID,
		&record.Name,
		&record.Description,
		&record.Status,
		&record.StartTime,
		&endTime,
		&durationMs,
		&record.AgentCount,
		&record.ToolCallCount,
		&record.TokenUsage,
		&costEstimate,
		&metadataJSON,
		&tagsJSON,
		&created_at,
		&updated_at,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan record: %w", err)
	}

	if endTime.Valid {
		record.EndTime = endTime.Time
	}
	if durationMs.Valid {
		record.Duration = time.Duration(durationMs.Int64) * time.Millisecond
	}
	if costEstimate.Valid {
		record.CostEstimate = costEstimate.Float64
	}
	if created_at.Valid {
		record.CreatedAt = created_at.Time
	}
	if updated_at.Valid {
		record.UpdatedAt = updated_at.Time
	}

	if err := json.Unmarshal([]byte(metadataJSON), &record.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	if err := json.Unmarshal([]byte(tagsJSON), &record.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return &record, nil
}

func (hs *HistoryStorage) QueryRecords(ctx context.Context, query ExecutionQuery) ([]*ExecutionRecord, int, error) {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	baseQuery := `SELECT 
		id, trace_id, pipeline_id, name, description, status,
		start_time, end_time, duration, agent_count, tool_call_count,
		token_usage, cost_estimate, metadata, tags, created_at, updated_at
	FROM execution_records WHERE 1=1`

	countQuery := `SELECT COUNT(*) FROM execution_records WHERE 1=1`

	args := []interface{}{}
	conditions := []string{}

	if len(query.IDs) > 0 {
		placeholders := ""
		for i, id := range query.IDs {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args = append(args, id)
		}
		conditions = append(conditions, "id IN ("+placeholders+")")
	}

	if len(query.Statuses) > 0 {
		placeholders := ""
		for i, status := range query.Statuses {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args = append(args, status)
		}
		conditions = append(conditions, "status IN ("+placeholders+")")
	}

	if query.PipelineID != "" {
		conditions = append(conditions, "pipeline_id = ?")
		args = append(args, query.PipelineID)
	}

	if query.SearchText != "" {
		// 使用全文搜索
		searchQuery := `SELECT rowid FROM records_fts WHERE records_fts MATCH ?`
		searchArgs := []interface{}{query.SearchText + "*"}

		var rowIDs []interface{}
		rows, err := hs.db.QueryContext(ctx, searchQuery, searchArgs...)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var rowID int64
				if err := rows.Scan(&rowID); err == nil {
					rowIDs = append(rowIDs, rowID)
				}
			}
		}

		if len(rowIDs) > 0 {
			placeholders := ""
			for i := range rowIDs {
				if i > 0 {
					placeholders += ","
				}
				placeholders += "?"
			}
			conditions = append(conditions, "rowid IN ("+placeholders+")")
			args = append(args, rowIDs...)
		}
	}

	if len(query.Tags) > 0 {
		// 简单的标签搜索（JSON数组包含）
		for _, tag := range query.Tags {
			conditions = append(conditions, "tags LIKE ?")
			args = append(args, "%\""+tag+"\"%")
		}
	}

	if !query.StartTime.IsZero() {
		conditions = append(conditions, "start_time >= ?")
		args = append(args, query.StartTime)
	}

	if !query.EndTime.IsZero() {
		conditions = append(conditions, "end_time <= ?")
		args = append(args, query.EndTime)
	}

	if len(conditions) > 0 {
		whereClause := " AND " + joinConditions(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// 获取总数
	var total int
	err := hs.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count records: %w", err)
	}

	// 排序
	orderBy := " ORDER BY "
	switch query.SortBy {
	case "start_time":
		orderBy += "start_time"
	case "end_time":
		orderBy += "end_time"
	case "duration":
		orderBy += "duration"
	case "agent_count":
		orderBy += "agent_count"
	case "token_usage":
		orderBy += "token_usage"
	case "cost_estimate":
		orderBy += "cost_estimate"
	case "created_at":
		orderBy += "created_at"
	default:
		orderBy += "start_time"
	}

	if query.SortOrder == "asc" {
		orderBy += " ASC"
	} else {
		orderBy += " DESC"
	}

	baseQuery += orderBy

	// 分页
	if query.Limit > 0 {
		baseQuery += " LIMIT ?"
		args = append(args, query.Limit)
	}

	if query.Offset > 0 {
		baseQuery += " OFFSET ?"
		args = append(args, query.Offset)
	}

	// 执行查询
	rows, err := hs.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query records: %w", err)
	}
	defer rows.Close()

	records := []*ExecutionRecord{}
	for rows.Next() {
		record, err := hs.scanRecord(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan record: %w", err)
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, total, nil
}

func (hs *HistoryStorage) scanRecord(rows *sql.Rows) (*ExecutionRecord, error) {
	var record ExecutionRecord
	var metadataJSON, tagsJSON string
	var durationMs sql.NullInt64
	var endTime, created_at, updated_at sql.NullTime
	var costEstimate sql.NullFloat64

	err := rows.Scan(
		&record.ID,
		&record.TraceID,
		&record.PipelineID,
		&record.Name,
		&record.Description,
		&record.Status,
		&record.StartTime,
		&endTime,
		&durationMs,
		&record.AgentCount,
		&record.ToolCallCount,
		&record.TokenUsage,
		&costEstimate,
		&metadataJSON,
		&tagsJSON,
		&created_at,
		&updated_at,
	)

	if err != nil {
		return nil, err
	}

	if endTime.Valid {
		record.EndTime = endTime.Time
	}
	if durationMs.Valid {
		record.Duration = time.Duration(durationMs.Int64) * time.Millisecond
	}
	if costEstimate.Valid {
		record.CostEstimate = costEstimate.Float64
	}
	if created_at.Valid {
		record.CreatedAt = created_at.Time
	}
	if updated_at.Valid {
		record.UpdatedAt = updated_at.Time
	}

	if err := json.Unmarshal([]byte(metadataJSON), &record.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	if err := json.Unmarshal([]byte(tagsJSON), &record.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return &record, nil
}

func (hs *HistoryStorage) DeleteRecord(ctx context.Context, id string) error {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	query := `DELETE FROM execution_records WHERE id = ?`
	_, err := hs.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	return nil
}

func (hs *HistoryStorage) GetStatistics(ctx context.Context, startTime, endTime time.Time) (map[string]interface{}, error) {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	stats := make(map[string]interface{})

	queries := map[string]string{
		"total_records":     `SELECT COUNT(*) FROM execution_records WHERE start_time >= ? AND start_time <= ?`,
		"completed_records": `SELECT COUNT(*) FROM execution_records WHERE status = 'completed' AND start_time >= ? AND start_time <= ?`,
		"failed_records":    `SELECT COUNT(*) FROM execution_records WHERE status = 'failed' AND start_time >= ? AND start_time <= ?`,
		"avg_duration":      `SELECT AVG(duration) FROM execution_records WHERE start_time >= ? AND start_time <= ? AND duration IS NOT NULL`,
		"total_token_usage": `SELECT SUM(token_usage) FROM execution_records WHERE start_time >= ? AND start_time <= ?`,
		"total_tool_calls":  `SELECT SUM(tool_call_count) FROM execution_records WHERE start_time >= ? AND start_time <= ?`,
		"total_cost":        `SELECT SUM(cost_estimate) FROM execution_records WHERE start_time >= ? AND start_time <= ? AND cost_estimate IS NOT NULL`,
	}

	for name, query := range queries {
		var value interface{}
		err := hs.db.QueryRowContext(ctx, query, startTime, endTime).Scan(&value)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to query %s: %w", name, err)
		}
		stats[name] = value
	}

	// 计算成功率
	total, _ := stats["total_records"].(int64)
	completed, _ := stats["completed_records"].(int64)
	if total > 0 {
		stats["success_rate"] = float64(completed) / float64(total) * 100
	} else {
		stats["success_rate"] = 0.0
	}

	// 平均Token使用量
	if total > 0 {
		totalToken, _ := stats["total_token_usage"].(int64)
		stats["avg_token_usage"] = totalToken / total
	} else {
		stats["avg_token_usage"] = 0
	}

	// 平均工具调用次数
	if total > 0 {
		totalTools, _ := stats["total_tool_calls"].(int64)
		stats["avg_tool_calls"] = totalTools / total
	} else {
		stats["avg_tool_calls"] = 0
	}

	// 平均成本
	if total > 0 {
		totalCost, _ := stats["total_cost"].(float64)
		stats["avg_cost"] = totalCost / float64(total)
	} else {
		stats["avg_cost"] = 0.0
	}

	// Agent统计
	agentStatsQuery := `SELECT 
		json_extract(metadata, '$.primary_agent') as agent_name,
		COUNT(*) as count,
		AVG(duration) as avg_duration,
		SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as success_count
	FROM execution_records 
	WHERE start_time >= ? AND start_time <= ? 
		AND json_extract(metadata, '$.primary_agent') IS NOT NULL
	GROUP BY json_extract(metadata, '$.primary_agent')
	ORDER BY count DESC`

	rows, err := hs.db.QueryContext(ctx, agentStatsQuery, startTime, endTime)
	if err != nil {
		return stats, nil // 不因为agent统计失败而返回错误
	}
	defer rows.Close()

	agentStats := make(map[string]map[string]interface{})
	for rows.Next() {
		var agentName sql.NullString
		var count int64
		var avgDuration sql.NullFloat64
		var successCount int64

		err := rows.Scan(&agentName, &count, &avgDuration, &successCount)
		if err != nil {
			continue
		}

		if agentName.Valid {
			agentData := map[string]interface{}{
				"count":         count,
				"success_count": successCount,
				"avg_duration":  0.0,
				"success_rate":  0.0,
			}

			if avgDuration.Valid {
				agentData["avg_duration"] = time.Duration(avgDuration.Float64) * time.Millisecond
			}
			if count > 0 {
				agentData["success_rate"] = float64(successCount) / float64(count) * 100
			}

			agentStats[agentName.String] = agentData
		}
	}

	stats["agent_statistics"] = agentStats

	return stats, nil
}

func (hs *HistoryStorage) CleanupOldData(ctx context.Context, retentionDays int) error {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	query := `DELETE FROM execution_records WHERE end_time < ?`

	result, err := hs.db.ExecContext(ctx, query, cutoffTime)
	if err != nil {
		return fmt.Errorf("failed to cleanup old data: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Cleaned up %d old history records\n", rowsAffected)

	return nil
}

func (hs *HistoryStorage) Close() error {
	return hs.db.Close()
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
