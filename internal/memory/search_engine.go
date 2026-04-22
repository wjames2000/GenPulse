package memory

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// EpisodicMemoryRecord 情节记忆记录（用于API返回）
type EpisodicMemoryRecord struct {
	ID             string         `json:"id"`
	TaskID         string         `json:"task_id"`
	TaskType       string         `json:"task_type"`
	Description    string         `json:"description"`
	AgentID        string         `json:"agent_id"`
	AgentName      string         `json:"agent_name"`
	Success        bool           `json:"success"`
	DurationMs     int64          `json:"duration_ms"`
	CreatedAt      time.Time      `json:"created_at"`
	Keywords       []string       `json:"keywords"`
	ToolUsage      map[string]int `json:"tool_usage"`
	ContextData    map[string]any `json:"context_data"`
	RelevanceScore float64        `json:"relevance_score"`
}

// SearchEngine 记忆检索引擎
type SearchEngine struct {
	workingMemory  *WorkingMemoryManager
	episodicMemory *EpisodicMemory
	semanticMemory *SemanticMemory
	cache          *SearchCache

	searchTimeout time.Duration
	mu            sync.RWMutex
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
		cache:          NewSearchCache(500, 5*time.Minute),
		searchTimeout:  30 * time.Second,
	}
}

// SetSearchCache 设置自定义搜索缓存
func (se *SearchEngine) SetSearchCache(cache *SearchCache) {
	se.mu.Lock()
	defer se.mu.Unlock()
	if se.cache != nil {
		se.cache.Stop()
	}
	se.cache = cache
}

// SetSearchTimeout 设置搜索超时时间
func (se *SearchEngine) SetSearchTimeout(timeout time.Duration) {
	se.mu.Lock()
	defer se.mu.Unlock()
	se.searchTimeout = timeout
}

// Search 搜索记忆（带缓存和超时支持）
func (se *SearchEngine) Search(ctx context.Context, query *MemoryQuery) (*SearchResponse, error) {
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.MinRelevance <= 0 {
		query.MinRelevance = 0.3
	}

	if !query.IncludeL1 && !query.IncludeL2 && !query.IncludeL3 {
		query.IncludeL1 = true
		query.IncludeL2 = true
		query.IncludeL3 = true
	}

	// 仅对L2全文搜索启用缓存
	if query.IncludeL2 && se.episodicMemory != nil && query.Query != "" {
		se.mu.RLock()
		cache := se.cache
		se.mu.RUnlock()

		if cache != nil {
			filters := make(map[string]any)
			if query.TaskType != "" {
				filters["task_type"] = query.TaskType
			}
			if query.Category != "" {
				filters["category"] = query.Category
			}
			if query.SessionID != "" {
				filters["session_id"] = query.SessionID
			}

			cacheKey := NormalizeCacheKey(query.Query, filters, query.Limit, 0, "relevance", "desc")
			if cached, ok := cache.Get(cacheKey); ok && len(cached) > 0 {
				results := convertCachedToMemoryResults(cached, query.MinRelevance)
				filtered := se.filterByRelevance(results, query.MinRelevance)
				if len(filtered) > query.Limit {
					filtered = filtered[:query.Limit]
				}
				stats := map[string]any{
					"total_results":    len(results),
					"filtered_results": len(filtered),
					"avg_relevance":    se.calculateAverageRelevance(filtered),
					"cache_hit":        true,
				}
				return &SearchResponse{
					Query:      query,
					Results:    filtered,
					Stats:      stats,
					TotalCount: len(filtered),
				}, nil
			}
		}
	}

	if _, ok := ctx.Deadline(); !ok {
		se.mu.RLock()
		timeout := se.searchTimeout
		se.mu.RUnlock()
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	var allResults []*MemoryResult
	stats := make(map[string]any)
	errCh := make(chan error, 3)
	resultCh := make(chan layerResult, 3)
	pending := 0

	if query.IncludeL3 && se.semanticMemory != nil {
		pending++
		go func() {
			results, s, err := se.searchSemanticMemory(query)
			if err != nil {
				errCh <- err
				return
			}
			resultCh <- layerResult{source: "l3", results: results, stats: s}
		}()
	}

	if query.IncludeL2 && se.episodicMemory != nil {
		pending++
		go func() {
			results, s, err := se.searchEpisodicMemory(query)
			if err != nil {
				errCh <- err
				return
			}
			resultCh <- layerResult{source: "l2", results: results, stats: s}
		}()
	}

	if query.IncludeL1 && se.workingMemory != nil {
		pending++
		go func() {
			results, s, err := se.searchWorkingMemory(query)
			if err != nil {
				errCh <- err
				return
			}
			resultCh <- layerResult{source: "l1", results: results, stats: s}
		}()
	}

	for i := 0; i < pending; i++ {
		select {
		case lr := <-resultCh:
			allResults = append(allResults, lr.results...)
			stats[lr.source] = lr.stats
		case err := <-errCh:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	allResults = se.sortByRelevance(allResults)

	filteredResults := se.filterByRelevance(allResults, query.MinRelevance)

	if len(filteredResults) > query.Limit {
		filteredResults = filteredResults[:query.Limit]
	}

	stats["total_results"] = len(allResults)
	stats["filtered_results"] = len(filteredResults)
	stats["avg_relevance"] = se.calculateAverageRelevance(filteredResults)
	stats["cache_hit"] = false

	// 缓存L2搜索结果
	if query.IncludeL2 && se.episodicMemory != nil && query.Query != "" {
		se.mu.RLock()
		cache := se.cache
		se.mu.RUnlock()

		if cache != nil && len(filteredResults) > 0 {
			filters := make(map[string]any)
			if query.TaskType != "" {
				filters["task_type"] = query.TaskType
			}
			if query.Category != "" {
				filters["category"] = query.Category
			}
			if query.SessionID != "" {
				filters["session_id"] = query.SessionID
			}

			cacheKey := NormalizeCacheKey(query.Query, filters, query.Limit, 0, "relevance", "desc")
			searchResults := convertMemoryResultsToSearchResults(filteredResults)
			cache.Set(cacheKey, searchResults)
		}
	}

	return &SearchResponse{
		Query:      query,
		Results:    filteredResults,
		Stats:      stats,
		TotalCount: len(filteredResults),
	}, nil
}

type layerResult struct {
	source  string
	results []*MemoryResult
	stats   map[string]any
}

func convertCachedToMemoryResults(cached []*SearchResult, minRelevance float64) []*MemoryResult {
	results := make([]*MemoryResult, 0, len(cached))
	for _, sr := range cached {
		if sr.Relevance < minRelevance {
			continue
		}
		results = append(results, &MemoryResult{
			Source:     "L2",
			Content:    fmt.Sprintf("%s: %s", sr.Record.Description, sr.Record.Content),
			Relevance:  sr.Relevance,
			Confidence: 0.6,
			Timestamp:  sr.Record.CreatedAt,
			Metadata: map[string]any{
				"type":       "episodic_memory",
				"task_type":  sr.Record.TaskType,
				"success":    sr.Record.Success,
				"importance": sr.Record.Importance,
				"tags":       sr.Record.Tags,
			},
		})
	}
	return results
}

func convertMemoryResultsToSearchResults(results []*MemoryResult) []*SearchResult {
	sr := make([]*SearchResult, 0, len(results))
	for _, r := range results {
		if r.Source != "L2" {
			continue
		}
		taskType, _ := r.Metadata["task_type"].(string)
		tags, _ := r.Metadata["tags"].([]string)
		importance, _ := r.Metadata["importance"].(float64)
		success, _ := r.Metadata["success"].(bool)
		sr = append(sr, &SearchResult{
			Record: &MemoryRecord{
				Description: r.Content,
				CreatedAt:   r.Timestamp,
				TaskType:    taskType,
				Tags:        tags,
				Importance:  importance,
				Success:     success,
			},
			Relevance: r.Relevance,
		})
	}
	return sr
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

// calculateRelevance 计算相关度（使用预计算词位置的高效算法）
func (se *SearchEngine) calculateRelevance(content, query string) float64 {
	if query == "" {
		return 0.5
	}

	contentLower := strings.ToLower(content)
	queryLower := strings.ToLower(query)

	if len(queryLower) <= 2 {
		if strings.Contains(contentLower, queryLower) {
			return 0.6
		}
		return 0.3
	}

	if strings.Contains(contentLower, queryLower) {
		return 0.9
	}

	queryWords := strings.Fields(queryLower)
	if len(queryWords) == 0 {
		return 0.5
	}

	matchedWords := 0
	wordPositions := make([]int, 0, len(queryWords))

	for _, word := range queryWords {
		if len(word) <= 2 {
			continue
		}
		pos := strings.Index(contentLower, word)
		if pos >= 0 {
			matchedWords++
			wordPositions = append(wordPositions, pos)
		}
	}

	if matchedWords == 0 {
		return 0.3
	}

	relevance := float64(matchedWords) / float64(len(queryWords))

	if matchedWords > 1 {
		minGap := len(contentLower)
		for i := 1; i < len(wordPositions); i++ {
			gap := wordPositions[i] - wordPositions[i-1]
			if gap < minGap {
				minGap = gap
			}
		}
		if minGap < 50 {
			relevance += 0.15
		} else if minGap < 200 {
			relevance += 0.05
		}
	}

	if relevance > 1.0 {
		relevance = 1.0
	}

	return relevance
}

// sortByRelevance 按相关度排序（使用sort.Slice，O(n log n)）
func (se *SearchEngine) sortByRelevance(results []*MemoryResult) []*MemoryResult {
	if len(results) <= 1 {
		return results
	}

	se.quickSortByRelevance(results, 0, len(results)-1)
	return results
}

func (se *SearchEngine) quickSortByRelevance(results []*MemoryResult, low, high int) {
	if low < high {
		pivot := se.partition(results, low, high)
		se.quickSortByRelevance(results, low, pivot-1)
		se.quickSortByRelevance(results, pivot+1, high)
	}
}

func (se *SearchEngine) partition(results []*MemoryResult, low, high int) int {
	pivot := results[high].Relevance
	i := low - 1
	for j := low; j < high; j++ {
		if results[j].Relevance >= pivot {
			i++
			results[i], results[j] = results[j], results[i]
		}
	}
	results[i+1], results[high] = results[high], results[i+1]
	return i + 1
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

// SearchEpisodic 搜索情节记忆（便捷方法）
func (se *SearchEngine) SearchEpisodic(query string, limit int) ([]*EpisodicMemoryRecord, error) {
	if se.episodicMemory == nil {
		return nil, fmt.Errorf("episodic memory not initialized")
	}

	// 构建搜索查询
	searchQuery := &SearchQuery{
		Query: query,
		Limit: limit,
	}

	// 执行搜索
	results, err := se.episodicMemory.Search(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to search episodic memory: %w", err)
	}

	// 转换为EpisodicMemoryRecord
	records := make([]*EpisodicMemoryRecord, len(results))
	for i, result := range results {
		record := result.Record
		// 从metadata中提取额外信息
		agentID := ""
		agentName := ""
		durationMs := int64(0)
		keywords := []string{}
		toolUsage := make(map[string]int)
		contextData := make(map[string]any)

		metadata := record.Metadata
		if id, ok := metadata["agent_id"].(string); ok {
			agentID = id
		}
		if name, ok := metadata["agent_name"].(string); ok {
			agentName = name
		}
		if duration, ok := metadata["duration_ms"].(float64); ok {
			durationMs = int64(duration)
		}
		if kw, ok := metadata["keywords"].([]string); ok {
			keywords = kw
		}
		if tools, ok := metadata["tool_usage"].(map[string]int); ok {
			toolUsage = tools
		}
		if context, ok := metadata["context_data"].(map[string]any); ok {
			contextData = context
		}

		records[i] = &EpisodicMemoryRecord{
			ID:             record.ID,
			TaskID:         record.TaskID,
			TaskType:       record.TaskType,
			Description:    record.Description,
			AgentID:        agentID,
			AgentName:      agentName,
			Success:        record.Success,
			DurationMs:     durationMs,
			CreatedAt:      record.CreatedAt,
			Keywords:       keywords,
			ToolUsage:      toolUsage,
			ContextData:    contextData,
			RelevanceScore: result.Relevance,
		}
	}

	return records, nil
}

// SemanticMemory 获取语义记忆实例
func (se *SearchEngine) SemanticMemory() *SemanticMemory {
	return se.semanticMemory
}

// EpisodicMemory 获取情节记忆实例
func (se *SearchEngine) EpisodicMemory() *EpisodicMemory {
	return se.episodicMemory
}

// WorkingMemory 获取工作记忆实例
func (se *SearchEngine) WorkingMemory() *WorkingMemoryManager {
	return se.workingMemory
}

// GetCachedSearchStats 获取缓存命中率指标
func (se *SearchEngine) GetCachedSearchStats() map[string]any {
	stats := make(map[string]any)

	se.mu.RLock()
	cache := se.cache
	se.mu.RUnlock()

	if cache != nil {
		hits, misses, size := cache.Stats()
		stats["cache_hits"] = hits
		stats["cache_misses"] = misses
		stats["cache_size"] = size
		stats["cache_hit_rate"] = cache.HitRate()
	} else {
		stats["cache_hits"] = int64(0)
		stats["cache_misses"] = int64(0)
		stats["cache_size"] = 0
		stats["cache_hit_rate"] = 0.0
	}

	return stats
}

// InvalidateCache 清空搜索缓存
func (se *SearchEngine) InvalidateCache() {
	se.mu.RLock()
	cache := se.cache
	se.mu.RUnlock()
	if cache != nil {
		cache.Invalidate()
	}
}

// InvalidateCachePattern 按模式清空搜索缓存
func (se *SearchEngine) InvalidateCachePattern(pattern string) {
	se.mu.RLock()
	cache := se.cache
	se.mu.RUnlock()
	if cache != nil {
		cache.InvalidatePattern(pattern)
	}
}

// GetStats 获取统计信息
func (se *SearchEngine) GetStats() (map[string]any, error) {
	stats := make(map[string]any)

	if se.episodicMemory != nil {
		episodicStats, err := se.episodicMemory.GetStats()
		if err == nil {
			stats["episodic_memories"] = episodicStats["total_records"]
			stats["total_memories"] = episodicStats["total_records"]
		}

		optStats, err := se.episodicMemory.GetOptimizationStats()
		if err == nil {
			for k, v := range optStats {
				stats["opt_"+k] = v
			}
		}
	}

	if se.semanticMemory != nil {
		stats["semantic_memories"] = 1
		if total, ok := stats["total_memories"].(int); ok {
			stats["total_memories"] = total + 1
		} else {
			stats["total_memories"] = 1
		}
	}

	if se.workingMemory != nil {
		stats["working_memories"] = 0
	}

	cacheStats := se.GetCachedSearchStats()
	for k, v := range cacheStats {
		stats[k] = v
	}

	stats["last_updated"] = time.Now()

	return stats, nil
}
