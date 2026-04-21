package skills

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

// LoadLevel 加载级别
type LoadLevel int

const (
	// L0 仅加载元数据
	L0 LoadLevel = iota
	// L1 加载完整内容
	L1
	// L2 加载完整内容+相关技能
	L2
)

// LoadStrategy 加载策略
type LoadStrategy string

const (
	// StrategyEager 急切加载：立即加载完整内容
	StrategyEager LoadStrategy = "eager"
	// StrategyLazy 懒加载：按需加载
	StrategyLazy LoadStrategy = "lazy"
	// StrategyPredictive 预测加载：基于使用模式预测加载
	StrategyPredictive LoadStrategy = "predictive"
)

// LoadOptions 加载选项
type LoadOptions struct {
	Level           LoadLevel       `json:"level" yaml:"level"`
	Strategy        LoadStrategy    `json:"strategy" yaml:"strategy"`
	MaxTokens       int             `json:"max_tokens" yaml:"max_tokens"`
	IncludeDisabled bool            `json:"include_disabled" yaml:"include_disabled"`
	IncludeInvalid  bool            `json:"include_invalid" yaml:"include_invalid"`
	Filter          map[string]any  `json:"filter" yaml:"filter"`
	Context         context.Context `json:"-" yaml:"-"`
}

// LoadResult 加载结果
type LoadResult struct {
	Skill      *Skill         `json:"skill,omitempty" yaml:"skill,omitempty"`
	Metadata   *SkillMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Level      LoadLevel      `json:"level" yaml:"level"`
	LoadedAt   time.Time      `json:"loaded_at" yaml:"loaded_at"`
	TokenCount int            `json:"token_count" yaml:"token_count"`
	CacheHit   bool           `json:"cache_hit" yaml:"cache_hit"`
	Error      string         `json:"error,omitempty" yaml:"error,omitempty"`
}

// ProgressiveLoader 渐进式披露加载器
type ProgressiveLoader struct {
	registry     *Registry
	cache        *LoadCache
	stats        *LoadStats
	strategy     LoadStrategy
	defaultLevel LoadLevel
	maxCacheSize int
	mu           sync.RWMutex
}

// LoadCache 加载缓存
type LoadCache struct {
	metadataCache map[string]*SkillMetadata
	fullCache     map[string]*Skill
	accessTime    map[string]time.Time
	size          int
	maxSize       int
}

// LoadStats 加载统计
type LoadStats struct {
	TotalLoads        int64               `json:"total_loads" yaml:"total_loads"`
	CacheHits         int64               `json:"cache_hits" yaml:"cache_hits"`
	CacheMisses       int64               `json:"cache_misses" yaml:"cache_misses"`
	TokenSavings      int64               `json:"token_savings" yaml:"token_savings"`
	AverageLoadTime   time.Duration       `json:"average_load_time" yaml:"average_load_time"`
	LevelDistribution map[LoadLevel]int64 `json:"level_distribution" yaml:"level_distribution"`
}

// NewProgressiveLoader 创建渐进式加载器
func NewProgressiveLoader(registry *Registry, options LoaderOptions) *ProgressiveLoader {
	loader := &ProgressiveLoader{
		registry:     registry,
		cache:        newLoadCache(options.MaxCacheSize),
		stats:        newLoadStats(),
		strategy:     options.DefaultStrategy,
		defaultLevel: options.DefaultLevel,
		maxCacheSize: options.MaxCacheSize,
	}

	return loader
}

// Load 加载技能
func (pl *ProgressiveLoader) Load(skillID string, options LoadOptions) (*LoadResult, error) {
	startTime := time.Now()

	// 设置默认选项
	if options.Level == 0 && options.Strategy == "" {
		options.Level = pl.defaultLevel
		options.Strategy = pl.strategy
	}

	if options.Context == nil {
		options.Context = context.Background()
	}

	// 检查缓存
	cachedResult, cacheHit := pl.checkCache(skillID, options.Level)
	if cacheHit {
		pl.recordCacheHit()
		result := &LoadResult{
			Skill:      cachedResult.Skill,
			Metadata:   cachedResult.Metadata,
			Level:      options.Level,
			LoadedAt:   time.Now(),
			TokenCount: cachedResult.TokenCount,
			CacheHit:   true,
		}

		pl.recordLoad(options.Level, time.Since(startTime), result.TokenCount)
		return result, nil
	}

	pl.recordCacheMiss()

	// 根据策略和级别加载
	var result *LoadResult
	var err error

	switch options.Strategy {
	case StrategyEager:
		result, err = pl.loadEager(skillID, options)
	case StrategyLazy:
		result, err = pl.loadLazy(skillID, options)
	case StrategyPredictive:
		result, err = pl.loadPredictive(skillID, options)
	default:
		result, err = pl.loadLazy(skillID, options)
	}

	if err != nil {
		return nil, err
	}

	// 更新缓存
	pl.updateCache(skillID, result, options.Level)

	// 记录统计
	pl.recordLoad(options.Level, time.Since(startTime), result.TokenCount)

	return result, nil
}

// LoadMultiple 批量加载技能
func (pl *ProgressiveLoader) LoadMultiple(skillIDs []string, options LoadOptions) ([]*LoadResult, error) {
	var results []*LoadResult
	var errors []error

	// 根据策略决定并行或串行加载
	if options.Strategy == StrategyEager && len(skillIDs) > 5 {
		// 急切加载且数量多时使用并行
		var wg sync.WaitGroup
		resultChan := make(chan *LoadResult, len(skillIDs))
		errorChan := make(chan error, len(skillIDs))

		for _, skillID := range skillIDs {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				result, err := pl.Load(id, options)
				if err != nil {
					errorChan <- err
				} else {
					resultChan <- result
				}
			}(skillID)
		}

		wg.Wait()
		close(resultChan)
		close(errorChan)

		for result := range resultChan {
			results = append(results, result)
		}

		for err := range errorChan {
			errors = append(errors, err)
		}
	} else {
		// 串行加载
		for _, skillID := range skillIDs {
			result, err := pl.Load(skillID, options)
			if err != nil {
				errors = append(errors, err)
			} else {
				results = append(results, result)
			}
		}
	}

	if len(errors) > 0 {
		return results, fmt.Errorf("failed to load some skills: %v", errors)
	}

	return results, nil
}

// LoadByFilter 按过滤器加载技能
func (pl *ProgressiveLoader) LoadByFilter(filter map[string]any, options LoadOptions) ([]*LoadResult, error) {
	// 首先获取匹配的元数据
	// 需要将map[string]any转换为map[string]string
	stringFilter := make(map[string]string)
	for k, v := range filter {
		if s, ok := v.(string); ok {
			stringFilter[k] = s
		}
	}
	metadatas, err := pl.registry.Search("", stringFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to search skills: %w", err)
	}

	// 提取技能ID
	var skillIDs []string
	for _, metadata := range metadatas {
		// 应用加载选项过滤
		if !options.IncludeDisabled && !metadata.Enabled {
			continue
		}
		if !options.IncludeInvalid && !metadata.Validated {
			continue
		}
		skillIDs = append(skillIDs, metadata.ID)
	}

	// 加载技能
	return pl.LoadMultiple(skillIDs, options)
}

// Preload 预加载技能
func (pl *ProgressiveLoader) Preload(skillIDs []string, level LoadLevel) {
	// 在后台预加载
	go func() {
		for _, skillID := range skillIDs {
			options := LoadOptions{
				Level:    level,
				Strategy: StrategyLazy,
			}
			pl.Load(skillID, options)
		}
	}()
}

// WarmupCache 预热缓存
func (pl *ProgressiveLoader) WarmupCache(count int, level LoadLevel) error {
	// 获取最常用的技能
	metadatas, err := pl.registry.List()
	if err != nil {
		return fmt.Errorf("failed to list skills: %w", err)
	}

	// 按使用次数排序
	sort.Slice(metadatas, func(i, j int) bool {
		return metadatas[i].UsageCount > metadatas[j].UsageCount
	})

	// 预加载前N个
	preloadCount := min(count, len(metadatas))
	var skillIDs []string
	for i := 0; i < preloadCount; i++ {
		skillIDs = append(skillIDs, metadatas[i].ID)
	}

	pl.Preload(skillIDs, level)
	return nil
}

// GetStats 获取统计信息
func (pl *ProgressiveLoader) GetStats() *LoadStats {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	stats := *pl.stats
	return &stats
}

// ClearCache 清空缓存
func (pl *ProgressiveLoader) ClearCache() {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.cache = newLoadCache(pl.maxCacheSize)
}

// SetStrategy 设置加载策略
func (pl *ProgressiveLoader) SetStrategy(strategy LoadStrategy) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.strategy = strategy
}

// SetDefaultLevel 设置默认加载级别
func (pl *ProgressiveLoader) SetDefaultLevel(level LoadLevel) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.defaultLevel = level
}

// loadEager 急切加载
func (pl *ProgressiveLoader) loadEager(skillID string, options LoadOptions) (*LoadResult, error) {
	// 急切加载总是加载完整内容
	skill, err := pl.registry.Get(skillID)
	if err != nil {
		return nil, fmt.Errorf("failed to load skill: %w", err)
	}

	tokenCount := pl.calculateTokenCount(skill, L1)

	result := &LoadResult{
		Skill:      skill,
		Metadata:   skill.ToMetadata(),
		Level:      L1,
		LoadedAt:   time.Now(),
		TokenCount: tokenCount,
		CacheHit:   false,
	}

	return result, nil
}

// loadLazy 懒加载
func (pl *ProgressiveLoader) loadLazy(skillID string, options LoadOptions) (*LoadResult, error) {
	switch options.Level {
	case L0:
		return pl.loadMetadata(skillID)
	case L1:
		return pl.loadFull(skillID)
	case L2:
		return pl.loadWithRelated(skillID)
	default:
		return pl.loadMetadata(skillID)
	}
}

// loadPredictive 预测加载
func (pl *ProgressiveLoader) loadPredictive(skillID string, options LoadOptions) (*LoadResult, error) {
	// 基于使用模式预测加载级别
	predictedLevel := pl.predictLoadLevel(skillID)

	// 如果预测级别低于请求级别，使用请求级别
	if predictedLevel < options.Level {
		predictedLevel = options.Level
	}

	// 加载
	switch predictedLevel {
	case L0:
		return pl.loadMetadata(skillID)
	case L1:
		return pl.loadFull(skillID)
	case L2:
		return pl.loadWithRelated(skillID)
	default:
		return pl.loadMetadata(skillID)
	}
}

// loadMetadata 加载元数据
func (pl *ProgressiveLoader) loadMetadata(skillID string) (*LoadResult, error) {
	metadata, err := pl.registry.GetMetadata(skillID)
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}

	tokenCount := pl.calculateTokenCountFromMetadata(metadata)

	result := &LoadResult{
		Metadata:   metadata,
		Level:      L0,
		LoadedAt:   time.Now(),
		TokenCount: tokenCount,
		CacheHit:   false,
	}

	return result, nil
}

// loadFull 加载完整内容
func (pl *ProgressiveLoader) loadFull(skillID string) (*LoadResult, error) {
	skill, err := pl.registry.Get(skillID)
	if err != nil {
		return nil, fmt.Errorf("failed to load skill: %w", err)
	}

	tokenCount := pl.calculateTokenCount(skill, L1)

	result := &LoadResult{
		Skill:      skill,
		Metadata:   skill.ToMetadata(),
		Level:      L1,
		LoadedAt:   time.Now(),
		TokenCount: tokenCount,
		CacheHit:   false,
	}

	return result, nil
}

// loadWithRelated 加载完整内容+相关技能
func (pl *ProgressiveLoader) loadWithRelated(skillID string) (*LoadResult, error) {
	// 加载主技能
	skill, err := pl.registry.Get(skillID)
	if err != nil {
		return nil, fmt.Errorf("failed to load skill: %w", err)
	}

	// 加载相关技能（前置技能）
	var relatedSkills []*Skill
	for _, prereqID := range skill.Prerequisites {
		relatedSkill, err := pl.registry.Get(prereqID)
		if err == nil {
			relatedSkills = append(relatedSkills, relatedSkill)
		}
	}

	// 计算总Token数
	tokenCount := pl.calculateTokenCount(skill, L2)
	for _, relatedSkill := range relatedSkills {
		tokenCount += pl.calculateTokenCount(relatedSkill, L0) // 只计算元数据
	}

	result := &LoadResult{
		Skill:      skill,
		Metadata:   skill.ToMetadata(),
		Level:      L2,
		LoadedAt:   time.Now(),
		TokenCount: tokenCount,
		CacheHit:   false,
	}

	return result, nil
}

// checkCache 检查缓存
func (pl *ProgressiveLoader) checkCache(skillID string, level LoadLevel) (*LoadResult, bool) {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	// 检查访问时间（缓存过期）
	if accessTime, exists := pl.cache.accessTime[skillID]; exists {
		if time.Since(accessTime) > time.Hour {
			// 缓存过期
			return nil, false
		}
	}

	// 根据级别检查缓存
	switch level {
	case L0:
		if metadata, exists := pl.cache.metadataCache[skillID]; exists {
			tokenCount := pl.calculateTokenCountFromMetadata(metadata)
			return &LoadResult{
				Metadata:   metadata,
				Level:      L0,
				TokenCount: tokenCount,
			}, true
		}
	case L1, L2:
		if skill, exists := pl.cache.fullCache[skillID]; exists {
			tokenCount := pl.calculateTokenCount(skill, level)
			return &LoadResult{
				Skill:      skill,
				Metadata:   skill.ToMetadata(),
				Level:      level,
				TokenCount: tokenCount,
			}, true
		}
	}

	return nil, false
}

// updateCache 更新缓存
func (pl *ProgressiveLoader) updateCache(skillID string, result *LoadResult, level LoadLevel) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	// 更新访问时间
	pl.cache.accessTime[skillID] = time.Now()

	// 根据级别更新缓存
	switch level {
	case L0:
		if result.Metadata != nil {
			pl.cache.metadataCache[skillID] = result.Metadata
			pl.cache.size++
		}
	case L1, L2:
		if result.Skill != nil {
			pl.cache.fullCache[skillID] = result.Skill
			pl.cache.size++

			// 同时缓存元数据
			pl.cache.metadataCache[skillID] = result.Metadata
		}
	}

	// 如果缓存超过最大大小，清理最久未使用的
	if pl.cache.size > pl.cache.maxSize {
		pl.cleanupCache()
	}
}

// cleanupCache 清理缓存
func (pl *ProgressiveLoader) cleanupCache() {
	// 找到最久未访问的项
	var oldestKey string
	var oldestTime time.Time

	for key, accessTime := range pl.cache.accessTime {
		if oldestKey == "" || accessTime.Before(oldestTime) {
			oldestKey = key
			oldestTime = accessTime
		}
	}

	if oldestKey != "" {
		// 从所有缓存中移除
		delete(pl.cache.metadataCache, oldestKey)
		delete(pl.cache.fullCache, oldestKey)
		delete(pl.cache.accessTime, oldestKey)
		pl.cache.size--
	}
}

// predictLoadLevel 预测加载级别
func (pl *ProgressiveLoader) predictLoadLevel(skillID string) LoadLevel {
	// 基于使用频率预测
	metadata, err := pl.registry.GetMetadata(skillID)
	if err != nil {
		return L0
	}

	// 使用次数越多，加载级别越高
	if metadata.UsageCount > 10 {
		return L1
	} else if metadata.UsageCount > 3 {
		return L0
	}

	return L0
}

// calculateTokenCount 计算Token数量
func (pl *ProgressiveLoader) calculateTokenCount(skill *Skill, level LoadLevel) int {
	tokenCount := 0

	// 总是包含元数据
	tokenCount += pl.calculateTokenCountFromMetadata(skill.ToMetadata())

	if level >= L1 {
		// 包含步骤
		for _, step := range skill.Steps {
			tokenCount += len(step.Action) / 4
			for _, param := range step.Parameters {
				tokenCount += len(param.Name) / 4
				tokenCount += len(param.Description) / 4
			}
			tokenCount += len(step.Expected) / 4
		}

		// 包含示例
		for _, example := range skill.Examples {
			tokenCount += len(example) / 4
		}

		// 包含技巧和警告
		for _, tip := range skill.Tips {
			tokenCount += len(tip) / 4
		}
		for _, warning := range skill.Warnings {
			tokenCount += len(warning) / 4
		}
	}

	return tokenCount
}

// calculateTokenCountFromMetadata 从元数据计算Token数量
func (pl *ProgressiveLoader) calculateTokenCountFromMetadata(metadata *SkillMetadata) int {
	tokenCount := 0

	tokenCount += len(metadata.Name) / 4
	tokenCount += len(metadata.Description) / 4
	tokenCount += len(metadata.Category) / 4
	for _, tag := range metadata.Tags {
		tokenCount += len(tag) / 4
	}

	return tokenCount
}

// recordCacheHit 记录缓存命中
func (pl *ProgressiveLoader) recordCacheHit() {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.stats.CacheHits++
	pl.stats.TotalLoads++
}

// recordCacheMiss 记录缓存未命中
func (pl *ProgressiveLoader) recordCacheMiss() {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.stats.CacheMisses++
	pl.stats.TotalLoads++
}

// recordLoad 记录加载统计
func (pl *ProgressiveLoader) recordLoad(level LoadLevel, duration time.Duration, tokenCount int) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	// 更新平均加载时间
	if pl.stats.TotalLoads > 0 {
		totalTime := pl.stats.AverageLoadTime * time.Duration(pl.stats.TotalLoads-1)
		pl.stats.AverageLoadTime = (totalTime + duration) / time.Duration(pl.stats.TotalLoads)
	} else {
		pl.stats.AverageLoadTime = duration
	}

	// 更新级别分布
	pl.stats.LevelDistribution[level]++

	// 更新Token节省（如果加载了元数据而不是完整内容）
	if level == L0 {
		// 估算完整内容Token数（假设是元数据的3倍）
		fullTokenCount := tokenCount * 3
		pl.stats.TokenSavings += int64(fullTokenCount - tokenCount)
	}
}

// newLoadCache 创建新缓存
func newLoadCache(maxSize int) *LoadCache {
	return &LoadCache{
		metadataCache: make(map[string]*SkillMetadata),
		fullCache:     make(map[string]*Skill),
		accessTime:    make(map[string]time.Time),
		maxSize:       maxSize,
	}
}

// newLoadStats 创建新统计
func newLoadStats() *LoadStats {
	return &LoadStats{
		LevelDistribution: make(map[LoadLevel]int64),
	}
}

// LoaderOptions 加载器选项
type LoaderOptions struct {
	DefaultStrategy LoadStrategy `json:"default_strategy" yaml:"default_strategy"`
	DefaultLevel    LoadLevel    `json:"default_level" yaml:"default_level"`
	MaxCacheSize    int          `json:"max_cache_size" yaml:"max_cache_size"`
}

// DefaultLoaderOptions 默认加载器选项
func DefaultLoaderOptions() LoaderOptions {
	return LoaderOptions{
		DefaultStrategy: StrategyLazy,
		DefaultLevel:    L0,
		MaxCacheSize:    100,
	}
}

// DefaultLoadOptions 默认加载选项
func DefaultLoadOptions() LoadOptions {
	return LoadOptions{
		Level:           L0,
		Strategy:        StrategyLazy,
		MaxTokens:       10000,
		IncludeDisabled: false,
		IncludeInvalid:  false,
		Filter:          nil,
		Context:         context.Background(),
	}
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
