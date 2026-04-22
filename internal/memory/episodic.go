package memory

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// EpisodicMemory 情节记忆（L2）- 存储任务执行记录
type EpisodicMemory struct {
	db      *sql.DB
	dbPath  string
	initSQL []string

	stmtMu           sync.Mutex
	stmtStore        *sql.Stmt
	stmtGet          *sql.Stmt
	stmtUpdate       *sql.Stmt
	stmtDelete       *sql.Stmt
	stmtUpdateAccess *sql.Stmt
}

// MemoryRecord 记忆记录
type MemoryRecord struct {
	ID           string         `json:"id"`
	SessionID    string         `json:"session_id"`
	TaskID       string         `json:"task_id"`
	TaskType     string         `json:"task_type"`
	Description  string         `json:"description"`
	Content      string         `json:"content"`
	Metadata     map[string]any `json:"metadata"`
	Tags         []string       `json:"tags"`
	Category     string         `json:"category"`
	Importance   float64        `json:"importance"` // 0-1 重要性评分
	Success      bool           `json:"success"`
	ErrorType    string         `json:"error_type,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	AccessedAt   time.Time      `json:"accessed_at"`
	AccessCount  int            `json:"access_count"`
	RelatedIDs   []string       `json:"related_ids,omitempty"`
}

// SearchQuery 搜索查询
type SearchQuery struct {
	Query         string         `json:"query"`
	Filters       map[string]any `json:"filters,omitempty"`
	Limit         int            `json:"limit,omitempty"`
	Offset        int            `json:"offset,omitempty"`
	SortBy        string         `json:"sort_by,omitempty"`    // relevance, created_at, updated_at, accessed_at, importance
	SortOrder     string         `json:"sort_order,omitempty"` // asc, desc
	MinImportance float64        `json:"min_importance,omitempty"`
	MaxImportance float64        `json:"max_importance,omitempty"`
	StartTime     *time.Time     `json:"start_time,omitempty"`
	EndTime       *time.Time     `json:"end_time,omitempty"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Record     *MemoryRecord `json:"record"`
	Relevance  float64       `json:"relevance"` // 0-1 相关度评分
	Highlights []string      `json:"highlights,omitempty"`
}

// NewEpisodicMemory 创建情节记忆
func NewEpisodicMemory(dbPath string) (*EpisodicMemory, error) {
	em := &EpisodicMemory{
		dbPath: dbPath,
		initSQL: []string{
			// 主表
			`CREATE TABLE IF NOT EXISTS memories (
				id TEXT PRIMARY KEY,
				session_id TEXT NOT NULL,
				task_id TEXT NOT NULL,
				task_type TEXT NOT NULL,
				description TEXT NOT NULL,
				content TEXT NOT NULL,
				metadata TEXT NOT NULL,
				tags TEXT NOT NULL,
				category TEXT NOT NULL,
				importance REAL NOT NULL DEFAULT 0.5,
				success INTEGER NOT NULL DEFAULT 1,
				error_type TEXT,
				error_message TEXT,
				created_at DATETIME NOT NULL,
				updated_at DATETIME NOT NULL,
				accessed_at DATETIME NOT NULL,
				access_count INTEGER NOT NULL DEFAULT 0,
				related_ids TEXT NOT NULL DEFAULT '[]'
			)`,

			// 创建索引
			`CREATE INDEX IF NOT EXISTS idx_memories_session_id ON memories(session_id)`,
			`CREATE INDEX IF NOT EXISTS idx_memories_task_id ON memories(task_id)`,
			`CREATE INDEX IF NOT EXISTS idx_memories_task_type ON memories(task_type)`,
			`CREATE INDEX IF NOT EXISTS idx_memories_category ON memories(category)`,
			`CREATE INDEX IF NOT EXISTS idx_memories_importance ON memories(importance)`,
			`CREATE INDEX IF NOT EXISTS idx_memories_success ON memories(success)`,
			`CREATE INDEX IF NOT EXISTS idx_memories_created_at ON memories(created_at)`,
			`CREATE INDEX IF NOT EXISTS idx_memories_updated_at ON memories(updated_at)`,
			`CREATE INDEX IF NOT EXISTS idx_memories_accessed_at ON memories(accessed_at)`,

			// FTS5虚拟表用于全文搜索（使用unicode61分词器支持中文）
			`CREATE VIRTUAL TABLE IF NOT EXISTS memories_fts USING fts5(
				id UNINDEXED,
				description,
				content,
				tags,
				category,
				tokenize='unicode61',
				content='memories',
				content_rowid='rowid'
			)`,

			// 触发器：插入时同步到FTS
			`CREATE TRIGGER IF NOT EXISTS memories_ai AFTER INSERT ON memories BEGIN
				INSERT INTO memories_fts(rowid, description, content, tags, category)
				VALUES (new.rowid, new.description, new.content, new.tags, new.category);
			END`,

			// 触发器：更新时同步到FTS
			`CREATE TRIGGER IF NOT EXISTS memories_au AFTER UPDATE ON memories BEGIN
				INSERT INTO memories_fts(memories_fts, rowid, description, content, tags, category)
				VALUES('delete', old.rowid, old.description, old.content, old.tags, old.category);
				INSERT INTO memories_fts(rowid, description, content, tags, category)
				VALUES (new.rowid, new.description, new.content, new.tags, new.category);
			END`,

			// 触发器：删除时从FTS删除
			`CREATE TRIGGER IF NOT EXISTS memories_ad AFTER DELETE ON memories BEGIN
				INSERT INTO memories_fts(memories_fts, rowid, description, content, tags, category)
				VALUES('delete', old.rowid, old.description, old.content, old.tags, old.category);
			END`,
		},
	}

	// 打开数据库连接
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 设置连接参数 - WAL模式下允许多个读取连接
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)
	db.SetConnMaxLifetime(30 * time.Minute)

	em.db = db

	// 初始化数据库
	if err := em.initialize(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 准备预处理语句
	if err := em.prepareStatements(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to prepare statements: %w", err)
	}

	return em, nil
}

// initialize 初始化数据库
func (em *EpisodicMemory) initialize() error {
	// 确保数据库目录存在
	dir := filepath.Dir(em.dbPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory %s: %w", dir, err)
		}
	}

	// 先执行PRAGMA设置（需要在首次连接时执行）
	pragmas := []string{
		`PRAGMA journal_mode=WAL`,
		`PRAGMA synchronous=NORMAL`,
		`PRAGMA cache_size=-64000`,
		`PRAGMA temp_store=MEMORY`,
		`PRAGMA mmap_size=268435456`,
	}
	for _, pragma := range pragmas {
		if _, err := em.db.Exec(pragma); err != nil {
			// PRAGMA失败不阻断初始化，仅记录
			fmt.Printf("Warning: failed to execute PRAGMA: %s, error: %v\n", pragma, err)
		}
	}

	// 执行初始化SQL（DDL语句）
	for _, sql := range em.initSQL {
		if _, err := em.db.Exec(sql); err != nil {
			return fmt.Errorf("failed to execute SQL: %s, error: %w", sql, err)
		}
	}

	return nil
}

// prepareStatements 预编译频繁执行的SQL语句
func (em *EpisodicMemory) prepareStatements() error {
	em.stmtMu.Lock()
	defer em.stmtMu.Unlock()

	var err error

	em.stmtStore, err = em.db.Prepare(`INSERT INTO memories (
		id, session_id, task_id, task_type, description, content,
		metadata, tags, category, importance, success, error_type,
		error_message, created_at, updated_at, accessed_at, access_count,
		related_ids
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare store statement: %w", err)
	}

	em.stmtGet, err = em.db.Prepare(`SELECT 
		id, session_id, task_id, task_type, description, content,
		metadata, tags, category, importance, success, error_type,
		error_message, created_at, updated_at, accessed_at, access_count,
		related_ids
	FROM memories WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare get statement: %w", err)
	}

	em.stmtUpdate, err = em.db.Prepare(`UPDATE memories SET
		session_id = ?, task_id = ?, task_type = ?, description = ?,
		content = ?, metadata = ?, tags = ?, category = ?,
		importance = ?, success = ?, error_type = ?, error_message = ?,
		updated_at = ?, related_ids = ?
	WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}

	em.stmtDelete, err = em.db.Prepare(`DELETE FROM memories WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}

	em.stmtUpdateAccess, err = em.db.Prepare(`UPDATE memories SET accessed_at = ?, access_count = access_count + 1 WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare updateAccess statement: %w", err)
	}

	return nil
}

// Close 关闭数据库连接
func (em *EpisodicMemory) Close() error {
	em.stmtMu.Lock()
	defer em.stmtMu.Unlock()

	if em.stmtStore != nil {
		em.stmtStore.Close()
	}
	if em.stmtGet != nil {
		em.stmtGet.Close()
	}
	if em.stmtUpdate != nil {
		em.stmtUpdate.Close()
	}
	if em.stmtDelete != nil {
		em.stmtDelete.Close()
	}
	if em.stmtUpdateAccess != nil {
		em.stmtUpdateAccess.Close()
	}

	if em.db != nil {
		return em.db.Close()
	}
	return nil
}

// Store 存储记忆记录
func (em *EpisodicMemory) Store(record *MemoryRecord) error {
	metadataJSON, tagsJSON, relatedIDsJSON, err := marshalRecordFields(record)
	if err != nil {
		return err
	}

	now := time.Now()
	if record.CreatedAt.IsZero() {
		record.CreatedAt = now
	}
	if record.UpdatedAt.IsZero() {
		record.UpdatedAt = now
	}
	if record.AccessedAt.IsZero() {
		record.AccessedAt = now
	}

	em.stmtMu.Lock()
	stmt := em.stmtStore
	em.stmtMu.Unlock()

	if stmt == nil {
		return fmt.Errorf("store statement not prepared")
	}

	_, err = stmt.Exec(
		record.ID,
		record.SessionID,
		record.TaskID,
		record.TaskType,
		record.Description,
		record.Content,
		string(metadataJSON),
		string(tagsJSON),
		record.Category,
		record.Importance,
		record.Success,
		record.ErrorType,
		record.ErrorMessage,
		record.CreatedAt,
		record.UpdatedAt,
		record.AccessedAt,
		record.AccessCount,
		string(relatedIDsJSON),
	)

	if err != nil {
		return fmt.Errorf("failed to store memory record: %w", err)
	}

	return nil
}

// StoreBatch 批量存储记忆记录（在单个事务中）
func (em *EpisodicMemory) StoreBatch(records []*MemoryRecord) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := em.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	em.stmtMu.Lock()
	stmt := em.stmtStore
	em.stmtMu.Unlock()

	if stmt == nil {
		return fmt.Errorf("store statement not prepared")
	}

	txStmt := tx.Stmt(stmt)

	now := time.Now()
	for _, record := range records {
		metadataJSON, tagsJSON, relatedIDsJSON, err := marshalRecordFields(record)
		if err != nil {
			return err
		}

		if record.CreatedAt.IsZero() {
			record.CreatedAt = now
		}
		if record.UpdatedAt.IsZero() {
			record.UpdatedAt = now
		}
		if record.AccessedAt.IsZero() {
			record.AccessedAt = now
		}

		_, err = txStmt.Exec(
			record.ID,
			record.SessionID,
			record.TaskID,
			record.TaskType,
			record.Description,
			record.Content,
			string(metadataJSON),
			string(tagsJSON),
			record.Category,
			record.Importance,
			record.Success,
			record.ErrorType,
			record.ErrorMessage,
			record.CreatedAt,
			record.UpdatedAt,
			record.AccessedAt,
			record.AccessCount,
			string(relatedIDsJSON),
		)
		if err != nil {
			return fmt.Errorf("failed to store memory record %s: %w", record.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch insert: %w", err)
	}

	return nil
}

func marshalRecordFields(record *MemoryRecord) (metadataJSON, tagsJSON, relatedIDsJSON []byte, err error) {
	metadataJSON, err = json.Marshal(record.Metadata)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	tagsJSON, err = json.Marshal(record.Tags)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to marshal tags: %w", err)
	}

	relatedIDsJSON, err = json.Marshal(record.RelatedIDs)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to marshal related_ids: %w", err)
	}

	return metadataJSON, tagsJSON, relatedIDsJSON, nil
}

// Get 获取记忆记录
func (em *EpisodicMemory) Get(id string) (*MemoryRecord, error) {
	em.stmtMu.Lock()
	stmt := em.stmtGet
	em.stmtMu.Unlock()

	if stmt == nil {
		return nil, fmt.Errorf("get statement not prepared")
	}

	row := stmt.QueryRow(id)

	record, err := em.scanRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("memory record not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get memory record: %w", err)
	}

	if err := em.updateAccess(id); err != nil {
		fmt.Printf("Warning: failed to update access time for record %s: %v\n", id, err)
	}

	return record, nil
}

// Update 更新记忆记录
func (em *EpisodicMemory) Update(record *MemoryRecord) error {
	metadataJSON, tagsJSON, relatedIDsJSON, err := marshalRecordFields(record)
	if err != nil {
		return err
	}

	record.UpdatedAt = time.Now()

	em.stmtMu.Lock()
	stmt := em.stmtUpdate
	em.stmtMu.Unlock()

	if stmt == nil {
		return fmt.Errorf("update statement not prepared")
	}

	result, err := stmt.Exec(
		record.SessionID,
		record.TaskID,
		record.TaskType,
		record.Description,
		record.Content,
		string(metadataJSON),
		string(tagsJSON),
		record.Category,
		record.Importance,
		record.Success,
		record.ErrorType,
		record.ErrorMessage,
		record.UpdatedAt,
		string(relatedIDsJSON),
		record.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update memory record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("memory record not found: %s", record.ID)
	}

	return nil
}

// Delete 删除记忆记录
func (em *EpisodicMemory) Delete(id string) error {
	em.stmtMu.Lock()
	stmt := em.stmtDelete
	em.stmtMu.Unlock()

	if stmt == nil {
		return fmt.Errorf("delete statement not prepared")
	}

	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete memory record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("memory record not found: %s", id)
	}

	return nil
}

// scanRow 扫描数据库行到MemoryRecord
func (em *EpisodicMemory) scanRow(row *sql.Row) (*MemoryRecord, error) {
	var record MemoryRecord
	var metadataJSON, tagsJSON, relatedIDsJSON string
	var successInt int

	err := row.Scan(
		&record.ID,
		&record.SessionID,
		&record.TaskID,
		&record.TaskType,
		&record.Description,
		&record.Content,
		&metadataJSON,
		&tagsJSON,
		&record.Category,
		&record.Importance,
		&successInt,
		&record.ErrorType,
		&record.ErrorMessage,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.AccessedAt,
		&record.AccessCount,
		&relatedIDsJSON,
	)

	if err != nil {
		return nil, err
	}

	// 反序列化JSON字段
	if err := json.Unmarshal([]byte(metadataJSON), &record.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	if err := json.Unmarshal([]byte(tagsJSON), &record.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	if err := json.Unmarshal([]byte(relatedIDsJSON), &record.RelatedIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal related_ids: %w", err)
	}

	record.Success = successInt == 1

	return &record, nil
}

// updateAccess 更新访问时间和计数
func (em *EpisodicMemory) updateAccess(id string) error {
	em.stmtMu.Lock()
	stmt := em.stmtUpdateAccess
	em.stmtMu.Unlock()

	if stmt == nil {
		return fmt.Errorf("updateAccess statement not prepared")
	}

	_, err := stmt.Exec(time.Now(), id)
	return err
}

// Search 搜索记忆记录
func (em *EpisodicMemory) Search(query *SearchQuery) ([]*SearchResult, error) {
	var whereClauses []string
	var args []any

	hasFullText := query.Query != ""

	if hasFullText {
		ftsQuery := `SELECT rowid FROM memories_fts WHERE memories_fts MATCH ?`
		whereClauses = append(whereClauses, fmt.Sprintf("rowid IN (%s)", ftsQuery))
		args = append(args, query.Query)
	}

	for key, value := range query.Filters {
		switch key {
		case "session_id", "task_id", "task_type", "category":
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", key))
			args = append(args, value)
		case "success":
			if success, ok := value.(bool); ok {
				whereClauses = append(whereClauses, "success = ?")
				args = append(args, success)
			}
		case "tags":
			if tag, ok := value.(string); ok {
				whereClauses = append(whereClauses, "tags LIKE ?")
				args = append(args, fmt.Sprintf("%%%s%%", tag))
			}
		}
	}

	if query.MinImportance > 0 {
		whereClauses = append(whereClauses, "importance >= ?")
		args = append(args, query.MinImportance)
	}
	if query.MaxImportance > 0 && query.MaxImportance <= 1 {
		whereClauses = append(whereClauses, "importance <= ?")
		args = append(args, query.MaxImportance)
	}

	if query.StartTime != nil {
		whereClauses = append(whereClauses, "created_at >= ?")
		args = append(args, *query.StartTime)
	}
	if query.EndTime != nil {
		whereClauses = append(whereClauses, "created_at <= ?")
		args = append(args, *query.EndTime)
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	var sortSQL string
	if hasFullText && query.SortBy == "relevance" {
		sortSQL = `ORDER BY (
			SELECT bm25(memories_fts, 0.0, 10.0, 5.0, 5.0)
			FROM memories_fts
			WHERE memories_fts.rowid = memories.rowid
		) ASC`
	} else {
		sortSQL = "ORDER BY "
		switch query.SortBy {
		case "created_at":
			sortSQL += "created_at"
		case "updated_at":
			sortSQL += "updated_at"
		case "accessed_at":
			sortSQL += "accessed_at"
		case "importance":
			sortSQL += "importance"
		case "access_count":
			sortSQL += "access_count"
		default:
			sortSQL += "created_at"
		}

		if query.SortOrder == "asc" {
			sortSQL += " ASC"
		} else {
			sortSQL += " DESC"
		}
	}

	limit := 50
	if query.Limit > 0 && query.Limit <= 100 {
		limit = query.Limit
	}

	offset := 0
	if query.Offset > 0 {
		offset = query.Offset
	}

	sqlQuery := fmt.Sprintf(`
		SELECT 
			id, session_id, task_id, task_type, description, content,
			metadata, tags, category, importance, success, error_type,
			error_message, created_at, updated_at, accessed_at, access_count,
			related_ids
		FROM memories %s %s LIMIT ? OFFSET ?
	`, whereSQL, sortSQL)

	args = append(args, limit, offset)

	rows, err := em.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search memory records: %w", err)
	}
	defer rows.Close()

	var results []*SearchResult
	for rows.Next() {
		var record MemoryRecord
		var metadataJSON, tagsJSON, relatedIDsJSON string
		var successInt int

		err := rows.Scan(
			&record.ID,
			&record.SessionID,
			&record.TaskID,
			&record.TaskType,
			&record.Description,
			&record.Content,
			&metadataJSON,
			&tagsJSON,
			&record.Category,
			&record.Importance,
			&successInt,
			&record.ErrorType,
			&record.ErrorMessage,
			&record.CreatedAt,
			&record.UpdatedAt,
			&record.AccessedAt,
			&record.AccessCount,
			&relatedIDsJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if err := json.Unmarshal([]byte(metadataJSON), &record.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		if err := json.Unmarshal([]byte(tagsJSON), &record.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		if err := json.Unmarshal([]byte(relatedIDsJSON), &record.RelatedIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal related_ids: %w", err)
		}

		record.Success = successInt == 1

		relevance := 0.5
		if hasFullText {
			relevance = calculateFTSRelevance(&record, query.Query)
		}

		results = append(results, &SearchResult{
			Record:    &record,
			Relevance: relevance,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

// calculateFTSRelevance 基于BM25算法计算相关度
func calculateFTSRelevance(record *MemoryRecord, query string) float64 {
	if query == "" {
		return 0.5
	}

	queryLower := strings.ToLower(query)
	queryWords := strings.Fields(queryLower)
	if len(queryWords) == 0 {
		return 0.5
	}

	fields := []string{
		strings.ToLower(record.Description),
		strings.ToLower(record.Content),
		strings.ToLower(strings.Join(record.Tags, " ")),
		strings.ToLower(record.Category),
	}

	fieldWeights := []float64{10.0, 5.0, 3.0, 2.0}
	totalScore := 0.0
	maxPossibleScore := 0.0

	for fi, field := range fields {
		fieldLen := len(strings.Fields(field))
		if fieldLen == 0 {
			continue
		}
		avgFieldLen := 50.0
		k1 := 1.2
		b := 0.75

		for _, word := range queryWords {
			if len(word) <= 2 {
				continue
			}
			maxPossibleScore += fieldWeights[fi]

			tf := float64(strings.Count(field, word))
			if tf == 0 {
				continue
			}

			docLenRatio := float64(fieldLen) / avgFieldLen
			score := fieldWeights[fi] * (tf * (k1 + 1)) / (tf + k1*(1-b+b*docLenRatio))
			totalScore += score
		}
	}

	if maxPossibleScore == 0 {
		if strings.Contains(strings.ToLower(record.Description+record.Content), queryLower) {
			return 0.7
		}
		return 0.5
	}

	relevance := totalScore / maxPossibleScore
	if relevance > 1.0 {
		relevance = 1.0
	}

	if strings.Contains(strings.ToLower(record.Description+record.Content), queryLower) {
		relevance = relevance*0.8 + 0.2
	}

	return relevance
}

// GetOptimizationStats 获取PRAGMA优化设置状态
func (em *EpisodicMemory) GetOptimizationStats() (map[string]any, error) {
	stats := make(map[string]any)

	pragmas := []struct {
		name string
		sql  string
	}{
		{"journal_mode", "PRAGMA journal_mode"},
		{"synchronous", "PRAGMA synchronous"},
		{"cache_size", "PRAGMA cache_size"},
		{"temp_store", "PRAGMA temp_store"},
		{"mmap_size", "PRAGMA mmap_size"},
		{"page_count", "PRAGMA page_count"},
		{"page_size", "PRAGMA page_size"},
		{"total_changes", "PRAGMA total_changes"},
	}

	for _, p := range pragmas {
		var val string
		err := em.db.QueryRow(p.sql).Scan(&val)
		if err != nil {
			stats[p.name] = fmt.Sprintf("error: %v", err)
		} else {
			stats[p.name] = val
		}
	}

	var fts5Count int
	err := em.db.QueryRow(`SELECT COUNT(*) FROM memories_fts`).Scan(&fts5Count)
	if err == nil {
		stats["fts5_entry_count"] = fts5Count
	}

	return stats, nil
}

// GetBySession 获取会话的所有记忆记录
func (em *EpisodicMemory) GetBySession(sessionID string) ([]*MemoryRecord, error) {
	query := `SELECT 
		id, session_id, task_id, task_type, description, content,
		metadata, tags, category, importance, success, error_type,
		error_message, created_at, updated_at, accessed_at, access_count,
		related_ids
	FROM memories WHERE session_id = ? ORDER BY created_at DESC`

	rows, err := em.db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session memories: %w", err)
	}
	defer rows.Close()

	var records []*MemoryRecord
	for rows.Next() {
		record, err := em.scanRowFromRows(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

// GetByTaskType 获取特定任务类型的记忆记录
func (em *EpisodicMemory) GetByTaskType(taskType string, limit int) ([]*MemoryRecord, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `SELECT 
		id, session_id, task_id, task_type, description, content,
		metadata, tags, category, importance, success, error_type,
		error_message, created_at, updated_at, accessed_at, access_count,
		related_ids
	FROM memories WHERE task_type = ? ORDER BY importance DESC, accessed_at DESC LIMIT ?`

	rows, err := em.db.Query(query, taskType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get memories by task type: %w", err)
	}
	defer rows.Close()

	var records []*MemoryRecord
	for rows.Next() {
		record, err := em.scanRowFromRows(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

// GetRecent 获取最近的记忆记录
func (em *EpisodicMemory) GetRecent(limit int) ([]*MemoryRecord, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `SELECT 
		id, session_id, task_id, task_type, description, content,
		metadata, tags, category, importance, success, error_type,
		error_message, created_at, updated_at, accessed_at, access_count,
		related_ids
	FROM memories ORDER BY created_at DESC LIMIT ?`

	rows, err := em.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent memories: %w", err)
	}
	defer rows.Close()

	var records []*MemoryRecord
	for rows.Next() {
		record, err := em.scanRowFromRows(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

// GetMostAccessed 获取访问次数最多的记忆记录
func (em *EpisodicMemory) GetMostAccessed(limit int) ([]*MemoryRecord, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `SELECT 
		id, session_id, task_id, task_type, description, content,
		metadata, tags, category, importance, success, error_type,
		error_message, created_at, updated_at, accessed_at, access_count,
		related_ids
	FROM memories ORDER BY access_count DESC, importance DESC LIMIT ?`

	rows, err := em.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get most accessed memories: %w", err)
	}
	defer rows.Close()

	var records []*MemoryRecord
	for rows.Next() {
		record, err := em.scanRowFromRows(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

// GetByImportance 获取重要性高的记忆记录
func (em *EpisodicMemory) GetByImportance(minImportance float64, limit int) ([]*MemoryRecord, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `SELECT 
		id, session_id, task_id, task_type, description, content,
		metadata, tags, category, importance, success, error_type,
		error_message, created_at, updated_at, accessed_at, access_count,
		related_ids
	FROM memories WHERE importance >= ? ORDER BY importance DESC, accessed_at DESC LIMIT ?`

	rows, err := em.db.Query(query, minImportance, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get important memories: %w", err)
	}
	defer rows.Close()

	var records []*MemoryRecord
	for rows.Next() {
		record, err := em.scanRowFromRows(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

// GetSuccessful 获取成功的记忆记录
func (em *EpisodicMemory) GetSuccessful(limit int) ([]*MemoryRecord, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `SELECT 
		id, session_id, task_id, task_type, description, content,
		metadata, tags, category, importance, success, error_type,
		error_message, created_at, updated_at, accessed_at, access_count,
		related_ids
	FROM memories WHERE success = 1 ORDER BY importance DESC, accessed_at DESC LIMIT ?`

	rows, err := em.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get successful memories: %w", err)
	}
	defer rows.Close()

	var records []*MemoryRecord
	for rows.Next() {
		record, err := em.scanRowFromRows(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

// GetFailed 获取失败的记忆记录
func (em *EpisodicMemory) GetFailed(limit int) ([]*MemoryRecord, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `SELECT 
		id, session_id, task_id, task_type, description, content,
		metadata, tags, category, importance, success, error_type,
		error_message, created_at, updated_at, accessed_at, access_count,
		related_ids
	FROM memories WHERE success = 0 ORDER BY importance DESC, accessed_at DESC LIMIT ?`

	rows, err := em.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get failed memories: %w", err)
	}
	defer rows.Close()

	var records []*MemoryRecord
	for rows.Next() {
		record, err := em.scanRowFromRows(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

// GetStats 获取统计信息
func (em *EpisodicMemory) GetStats() (map[string]any, error) {
	stats := make(map[string]any)

	// 总记录数
	var totalCount int
	err := em.db.QueryRow("SELECT COUNT(*) FROM memories").Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}
	stats["total_count"] = totalCount

	// 成功记录数
	var successCount int
	err = em.db.QueryRow("SELECT COUNT(*) FROM memories WHERE success = 1").Scan(&successCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get success count: %w", err)
	}
	stats["success_count"] = successCount

	// 失败记录数
	var failedCount int
	err = em.db.QueryRow("SELECT COUNT(*) FROM memories WHERE success = 0").Scan(&failedCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get failed count: %w", err)
	}
	stats["failed_count"] = failedCount

	// 成功率
	successRate := 0.0
	if totalCount > 0 {
		successRate = float64(successCount) / float64(totalCount)
	}
	stats["success_rate"] = successRate

	// 平均重要性
	var avgImportance float64
	err = em.db.QueryRow("SELECT AVG(importance) FROM memories").Scan(&avgImportance)
	if err != nil {
		return nil, fmt.Errorf("failed to get average importance: %w", err)
	}
	stats["avg_importance"] = avgImportance

	// 总访问次数
	var totalAccessCount int
	err = em.db.QueryRow("SELECT SUM(access_count) FROM memories").Scan(&totalAccessCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total access count: %w", err)
	}
	stats["total_access_count"] = totalAccessCount

	// 按任务类型统计
	query := `SELECT task_type, COUNT(*) as count, AVG(importance) as avg_importance, 
		SUM(CASE WHEN success = 1 THEN 1 ELSE 0 END) as success_count
		FROM memories GROUP BY task_type ORDER BY count DESC`

	rows, err := em.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get task type stats: %w", err)
	}
	defer rows.Close()

	taskTypeStats := make(map[string]map[string]any)
	for rows.Next() {
		var taskType string
		var count, successCount int
		var avgImportance float64

		err := rows.Scan(&taskType, &count, &avgImportance, &successCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task type stats: %w", err)
		}

		taskTypeStats[taskType] = map[string]any{
			"count":          count,
			"avg_importance": avgImportance,
			"success_count":  successCount,
			"success_rate":   float64(successCount) / float64(count),
		}
	}
	stats["task_type_stats"] = taskTypeStats

	// 按类别统计
	categoryQuery := `SELECT category, COUNT(*) as count, AVG(importance) as avg_importance
		FROM memories GROUP BY category ORDER BY count DESC`

	categoryRows, err := em.db.Query(categoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get category stats: %w", err)
	}
	defer categoryRows.Close()

	categoryStats := make(map[string]map[string]any)
	for categoryRows.Next() {
		var category string
		var count int
		var avgImportance float64

		err := categoryRows.Scan(&category, &count, &avgImportance)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category stats: %w", err)
		}

		categoryStats[category] = map[string]any{
			"count":          count,
			"avg_importance": avgImportance,
		}
	}
	stats["category_stats"] = categoryStats

	// 最近活动时间
	var lastActivityStr string
	err = em.db.QueryRow("SELECT MAX(accessed_at) FROM memories").Scan(&lastActivityStr)
	if err == nil && lastActivityStr != "" {
		if lastActivity, err := time.Parse("2006-01-02 15:04:05", lastActivityStr); err == nil {
			stats["last_activity"] = lastActivity
		}
	}

	// 最早记录时间
	var firstRecordStr string
	err = em.db.QueryRow("SELECT MIN(created_at) FROM memories").Scan(&firstRecordStr)
	if err == nil && firstRecordStr != "" {
		if firstRecord, err := time.Parse("2006-01-02 15:04:05", firstRecordStr); err == nil {
			stats["first_record"] = firstRecord
		}
	}

	return stats, nil
}

// Cleanup 清理旧记录
func (em *EpisodicMemory) Cleanup(maxAgeDays int, minImportance float64) (int, error) {
	if maxAgeDays <= 0 {
		maxAgeDays = 30 // 默认保留30天
	}

	cutoffTime := time.Now().AddDate(0, 0, -maxAgeDays)

	query := `DELETE FROM memories WHERE created_at < ? AND importance < ?`

	result, err := em.db.Exec(query, cutoffTime, minImportance)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old records: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

// scanRowFromRows 从sql.Rows扫描记录
func (em *EpisodicMemory) scanRowFromRows(rows *sql.Rows) (*MemoryRecord, error) {
	var record MemoryRecord
	var metadataJSON, tagsJSON, relatedIDsJSON string
	var successInt int

	err := rows.Scan(
		&record.ID,
		&record.SessionID,
		&record.TaskID,
		&record.TaskType,
		&record.Description,
		&record.Content,
		&metadataJSON,
		&tagsJSON,
		&record.Category,
		&record.Importance,
		&successInt,
		&record.ErrorType,
		&record.ErrorMessage,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.AccessedAt,
		&record.AccessCount,
		&relatedIDsJSON,
	)

	if err != nil {
		return nil, err
	}

	// 反序列化JSON字段
	if err := json.Unmarshal([]byte(metadataJSON), &record.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	if err := json.Unmarshal([]byte(tagsJSON), &record.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	if err := json.Unmarshal([]byte(relatedIDsJSON), &record.RelatedIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal related_ids: %w", err)
	}

	record.Success = successInt == 1

	return &record, nil
}
