package memory

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// SearchEngine 记忆检索引擎
type SearchEngine struct {
	workingMemory  *WorkingMemoryManager
	episodicMemory *EpisodicMemory
	semanticMemory *SemanticMemory
}

// MemoryQuery 记忆查询
type MemoryQuery struct {
	Query        string  `json:"query"`
	TaskType     string  `json:"task_type,omitempty"`
	Category     string  `json:"category,omitempty"`
	SessionID    string  `json:"session_id,omitempty"`
	MinRelevance float64 `json:"min_relevance,omitempty"`
	Limit        int     `json:"limit,omitempty"`
	IncludeL1    bool    `json:"include_l1,omitempty"` // 是否包含工作记忆
	IncludeL2    bool    `json:"include_l2,omitempty"` // 是否包含情节记忆
	IncludeL3    bool    `json:"include_l3,omitempty"` // 是否包含语义记忆
}

// MemoryResult 记忆结果
type MemoryResult struct {
	Source     string         `json:"source"` // L1, L2, L3
	Content    string         `json:"content"`
	Relevance  float64        `json:"relevance"`  // 0-1 相关度
	Confidence float64        `json:"confidence"` // 0-1 置信度
	Timestamp  time.Time      `json:"timestamp"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Query      *MemoryQuery    `json:"query"`
	Results    []*MemoryResult `json:"results"`
	Stats      map[string]any  `json:"stats,omitempty"`
	TotalCount int             `json:"total_count"`
}

// NewSearchEngine 创建记忆检索引擎
func NewSearchEngine(wm *WorkingMemoryManager, em *EpisodicMemory, sm *SemanticMemory) *SearchEngine {
	return &SearchEngine{
		workingMemory:  wm,
		episodicMemory: em,
		semanticMemory: sm,
	}
}

// Search 搜索记忆
func (se *SearchEngine) Search(ctx context.Context, query *MemoryQuery) (*SearchResponse, error) {
	// 设置默认值
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.MinRelevance <= 0 {
		query.MinRelevance = 0.3
	}

	// 如果没有指定包含哪些层，默认全部包含
	if !query.IncludeL1 && !query.IncludeL2 && !query.IncludeL3 {
		query.IncludeL1 = true
		query.IncludeL2 = true
		query.IncludeL3 = true
	}

	var allResults []*MemoryResult
	stats := make(map[string]any)

	// L3语义记忆搜索（最高优先级）
	if query.IncludeL3 && se.semanticMemory != nil {
		l3Results, l3Stats, err := se.searchSemanticMemory(query)
		if err != nil {
			return nil, fmt.Errorf("failed to search semantic memory: %w", err)
		}
		allResults = append(allResults, l3Results...)
		stats["l3"] = l3Stats
	}

	// L2情节记忆搜索
	if query.IncludeL2 && se.episodicMemory != nil {
		l2Results, l2Stats, err := se.searchEpisodicMemory(query)
		if err != nil {
			return nil, fmt.Errorf("failed to search episodic memory: %w", err)
		}
		allResults = append(allResults, l2Results...)
		stats["l2"] = l2Stats
	}

	// L1工作记忆搜索
	if query.IncludeL1 && se.workingMemory != nil {
		l1Results, l1Stats, err := se.searchWorkingMemory(query)
		if err != nil {
			return nil, fmt.Errorf("failed to search working memory: %w", err)
		}
		allResults = append(allResults, l1Results...)
		stats["l1"] = l1Stats
	}

	// 按相关度排序
	allResults = se.sortByRelevance(allResults)

	// 过滤低相关度结果
	filteredResults := se.filterByRelevance(allResults, query.MinRelevance)

	// 限制结果数量
	if len(filteredResults) > query.Limit {
		filteredResults = filteredResults[:query.Limit]
	}

	// 计算统计信息
	stats["total_results"] = len(allResults)
	stats["filtered_results"] = len(filteredResults)
	stats["avg_relevance"] = se.calculateAverageRelevance(filteredResults)

	return &SearchResponse{
		Query:      query,
		Results:    filteredResults,
		Stats:      stats,
		TotalCount: len(filteredResults),
	}, nil
}

// searchSemanticMemory 搜索语义记忆（L3）
func (se *SearchEngine) searchSemanticMemory(query *MemoryQuery) ([]*MemoryResult, map[string]any, error) {
	var results []*MemoryResult
	stats := make(map[string]any)

	// 搜索记忆内容
	semanticQuery := &SemanticQuery{
		Query:    query.Query,
		TaskType: query.TaskType,
		Category: query.Category,
		Limit:    query.Limit,
	}

	contentResults, err := se.semanticMemory.Search(semanticQuery)
	if err != nil {
		return nil, nil, err
	}

	// 转换结果
	for _, content := range contentResults {
		relevance := se.calculateRelevance(content, query.Query)
		if relevance >= query.MinRelevance {
			results = append(results, &MemoryResult{
				Source:     "L3",
				Content:    content,
				Relevance:  relevance,
				Confidence: 0.8, // 语义记忆置信度较高
				Timestamp:  time.Now(),
				Metadata: map[string]any{
					"type": "semantic_memory",
				},
			})
		}
	}

	// 获取任务建议
	if query.TaskType != "" {
		advice, err := se.semanticMemory.GetTaskAdvice(query.TaskType)
		if err == nil && advice != "" {
			results = append(results, &MemoryResult{
				Source:     "L3",
				Content:    advice,
				Relevance:  0.9, // 任务建议相关度很高
				Confidence: 0.7,
				Timestamp:  time.Now(),
				Metadata: map[string]any{
					"type":      "task_advice",
					"task_type": query.TaskType,
				},
			})
		}
	}

	stats["count"] = len(results)
	stats["source"] = "semantic_memory"

	return results, stats, nil
}

// searchEpisodicMemory 搜索情节记忆（L2）
func (se *SearchEngine) searchEpisodicMemory(query *MemoryQuery) ([]*MemoryResult, map[string]any, error) {
	var results []*MemoryResult
	stats := make(map[string]any)

	// 构建搜索查询
	searchQuery := &SearchQuery{
		Query:   query.Query,
		Filters: make(map[string]any),
		Limit:   query.Limit,
		SortBy:  "relevance",
	}

	// 添加过滤器
	if query.TaskType != "" {
		searchQuery.Filters["task_type"] = query.TaskType
	}
	if query.Category != "" {
		searchQuery.Filters["category"] = query.Category
	}
	if query.SessionID != "" {
		searchQuery.Filters["session_id"] = query.SessionID
	}

	// 执行搜索
	episodicResults, err := se.episodicMemory.Search(searchQuery)
	if err != nil {
		return nil, nil, err
	}

	// 转换结果
	for _, result := range episodicResults {
		content := fmt.Sprintf("%s: %s", result.Record.Description, result.Record.Content)
		relevance := result.Relevance

		if relevance >= query.MinRelevance {
			results = append(results, &MemoryResult{
				Source:     "L2",
				Content:    content,
				Relevance:  relevance,
				Confidence: 0.6, // 情节记忆置信度中等
				Timestamp:  result.Record.CreatedAt,
				Metadata: map[string]any{
					"type":       "episodic_memory",
					"task_type":  result.Record.TaskType,
					"success":    result.Record.Success,
					"importance": result.Record.Importance,
					"tags":       result.Record.Tags,
				},
			})
		}
	}

	// 如果没有搜索结果，尝试获取相关任务类型的记忆
	if len(results) == 0 && query.TaskType != "" {
		relatedMemories, err := se.episodicMemory.GetByTaskType(query.TaskType, 5)
		if err == nil && len(relatedMemories) > 0 {
			for _, memory := range relatedMemories {
				content := fmt.Sprintf("%s: %s", memory.Description, memory.Content)
				relevance := 0.5 // 默认相关度

				results = append(results, &MemoryResult{
					Source:     "L2",
					Content:    content,
					Relevance:  relevance,
					Confidence: 0.5,
					Timestamp:  memory.CreatedAt,
					Metadata: map[string]any{
						"type":       "related_memory",
						"task_type":  memory.TaskType,
						"success":    memory.Success,
						"importance": memory.Importance,
					},
				})
			}
		}
	}

	stats["count"] = len(results)
	stats["source"] = "episodic_memory"

	return results, stats, nil
}

// searchWorkingMemory 搜索工作记忆（L1）
func (se *SearchEngine) searchWorkingMemory(query *MemoryQuery) ([]*MemoryResult, map[string]any, error) {
	var results []*MemoryResult
	stats := make(map[string]any)

	// 获取当前会话的工作记忆
	sessionID := query.SessionID
	if sessionID == "" {
		// 如果没有指定会话ID，使用默认会话
		sessionID = "default"
	}

	// 获取会话
	session, exists := se.workingMemory.GetSession(sessionID)
	if !exists {
		// 会话不存在，返回空结果
		stats["count"] = 0
		stats["source"] = "working_memory"
		return results, stats, nil
	}

	// 获取会话中的所有数据
	sessionData := session.GetAll()

	// 搜索会话数据
	for key, value := range sessionData {
		// 将值转换为字符串
		valueStr := fmt.Sprintf("%v", value)

		// 计算相关度
		relevance := se.calculateRelevance(valueStr, query.Query)
		if relevance >= query.MinRelevance {
			results = append(results, &MemoryResult{
				Source:     "L1",
				Content:    fmt.Sprintf("%s: %s", key, valueStr),
				Relevance:  relevance,
				Confidence: 0.4, // 工作记忆置信度较低
				Timestamp:  time.Now(),
				Metadata: map[string]any{
					"type":       "working_memory",
					"key":        key,
					"session_id": sessionID,
				},
			})
		}
	}

	stats["count"] = len(results)
	stats["source"] = "working_memory"

	return results, stats, nil
}

// calculateRelevance 计算相关度
func (se *SearchEngine) calculateRelevance(content, query string) float64 {
	if query == "" {
		return 0.5 // 默认相关度
	}

	contentLower := strings.ToLower(content)
	queryLower := strings.ToLower(query)

	// 简单关键词匹配
	queryWords := strings.Fields(queryLower)
	matchedWords := 0

	for _, word := range queryWords {
		if len(word) > 2 && strings.Contains(contentLower, word) {
			matchedWords++
		}
	}

	// 计算相关度
	if len(queryWords) == 0 {
		return 0.5
	}

	relevance := float64(matchedWords) / float64(len(queryWords))

	// 完全匹配加分
	if strings.Contains(contentLower, queryLower) {
		relevance = min(1.0, relevance+0.3)
	}

	return relevance
}

// sortByRelevance 按相关度排序
func (se *SearchEngine) sortByRelevance(results []*MemoryResult) []*MemoryResult {
	// 使用冒泡排序（简单实现）
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Relevance > results[i].Relevance {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	return results
}

// filterByRelevance 按相关度过滤
func (se *SearchEngine) filterByRelevance(results []*MemoryResult, minRelevance float64) []*MemoryResult {
	var filtered []*MemoryResult
	for _, result := range results {
		if result.Relevance >= minRelevance {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

// calculateAverageRelevance 计算平均相关度
func (se *SearchEngine) calculateAverageRelevance(results []*MemoryResult) float64 {
	if len(results) == 0 {
		return 0
	}

	total := 0.0
	for _, result := range results {
		total += result.Relevance
	}

	return total / float64(len(results))
}

// GetContextualMemory 获取上下文相关记忆
func (se *SearchEngine) GetContextualMemory(ctx context.Context, taskType, sessionID string) (*SearchResponse, error) {
	query := &MemoryQuery{
		TaskType:  taskType,
		SessionID: sessionID,
		Limit:     10,
		IncludeL1: true,
		IncludeL2: true,
		IncludeL3: true,
	}

	return se.Search(ctx, query)
}

// GetLearningInsights 获取学习洞察
func (se *SearchEngine) GetLearningInsights(ctx context.Context, taskType string) (*SearchResponse, error) {
	query := &MemoryQuery{
		TaskType:     taskType,
		Query:        "学习 经验 教训 成功 失败",
		Limit:        15,
		IncludeL2:    true,
		IncludeL3:    true,
		MinRelevance: 0.4,
	}

	return se.Search(ctx, query)
}

// GetBestPractices 获取最佳实践
func (se *SearchEngine) GetBestPractices(ctx context.Context, taskType string) (*SearchResponse, error) {
	query := &MemoryQuery{
		TaskType:     taskType,
		Query:        "最佳实践 建议 技巧 方法",
		Limit:        10,
		IncludeL2:    true,
		IncludeL3:    true,
		MinRelevance: 0.5,
	}

	return se.Search(ctx, query)
}

// min 返回最小值
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
