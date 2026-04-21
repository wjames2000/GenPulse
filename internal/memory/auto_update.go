package memory

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// AutoUpdateManager 记忆自动更新管理器
type AutoUpdateManager struct {
	workingMemory  *WorkingMemoryManager
	episodicMemory *EpisodicMemory
	semanticMemory *SemanticMemory
	searchEngine   *SearchEngine
	config         *AutoUpdateConfig
}

// AutoUpdateConfig 自动更新配置
type AutoUpdateConfig struct {
	EnableL2Update      bool    `json:"enable_l2_update"`       // 是否启用L2更新
	EnableL3Update      bool    `json:"enable_l3_update"`       // 是否启用L3更新
	MinImportance       float64 `json:"min_importance"`         // 最小重要性阈值
	SuccessThreshold    float64 `json:"success_threshold"`      // 成功阈值
	LearningThreshold   float64 `json:"learning_threshold"`     // 学习阈值
	UpdateInterval      int     `json:"update_interval"`        // 更新间隔（秒）
	MaxRecordsPerUpdate int     `json:"max_records_per_update"` // 每次更新最大记录数
}

// TaskResult 任务结果
type TaskResult struct {
	TaskID       string         `json:"task_id"`
	SessionID    string         `json:"session_id"`
	TaskType     string         `json:"task_type"`
	Description  string         `json:"description"`
	Content      string         `json:"content"`
	Success      bool           `json:"success"`
	ErrorType    string         `json:"error_type,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	Duration     time.Duration  `json:"duration"`
	Importance   float64        `json:"importance"` // 0-1 重要性
	Tags         []string       `json:"tags,omitempty"`
	Category     string         `json:"category,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	Insights     []string       `json:"insights,omitempty"` // 学习洞察
}

// DefaultAutoUpdateConfig 默认自动更新配置
var DefaultAutoUpdateConfig = &AutoUpdateConfig{
	EnableL2Update:      true,
	EnableL3Update:      true,
	MinImportance:       0.3,
	SuccessThreshold:    0.7,
	LearningThreshold:   0.5,
	UpdateInterval:      300, // 5分钟
	MaxRecordsPerUpdate: 10,
}

// NewAutoUpdateManager 创建自动更新管理器
func NewAutoUpdateManager(wm *WorkingMemoryManager, em *EpisodicMemory, sm *SemanticMemory, se *SearchEngine) *AutoUpdateManager {
	return &AutoUpdateManager{
		workingMemory:  wm,
		episodicMemory: em,
		semanticMemory: sm,
		searchEngine:   se,
		config:         DefaultAutoUpdateConfig,
	}
}

// WithConfig 设置配置
func (aum *AutoUpdateManager) WithConfig(config *AutoUpdateConfig) *AutoUpdateManager {
	aum.config = config
	return aum
}

// RecordTaskResult 记录任务结果并自动更新记忆
func (aum *AutoUpdateManager) RecordTaskResult(ctx context.Context, result *TaskResult) error {
	// 验证重要性阈值
	if result.Importance < aum.config.MinImportance {
		return nil // 重要性太低，不记录
	}

	// 更新L2情节记忆
	if aum.config.EnableL2Update {
		if err := aum.updateEpisodicMemory(result); err != nil {
			return fmt.Errorf("failed to update episodic memory: %w", err)
		}
	}

	// 更新L3语义记忆
	if aum.config.EnableL3Update {
		if err := aum.updateSemanticMemory(result); err != nil {
			return fmt.Errorf("failed to update semantic memory: %w", err)
		}
	}

	// 更新L1工作记忆（会话上下文）
	if err := aum.updateWorkingMemory(result); err != nil {
		return fmt.Errorf("failed to update working memory: %w", err)
	}

	return nil
}

// updateEpisodicMemory 更新情节记忆（L2）
func (aum *AutoUpdateManager) updateEpisodicMemory(result *TaskResult) error {
	if aum.episodicMemory == nil {
		return nil
	}

	// 创建记忆记录
	record := &MemoryRecord{
		ID:           fmt.Sprintf("task_%s_%d", result.TaskID, time.Now().UnixNano()),
		SessionID:    result.SessionID,
		TaskID:       result.TaskID,
		TaskType:     result.TaskType,
		Description:  result.Description,
		Content:      result.Content,
		Metadata:     result.Metadata,
		Tags:         result.Tags,
		Category:     result.Category,
		Importance:   result.Importance,
		Success:      result.Success,
		ErrorType:    result.ErrorType,
		ErrorMessage: result.ErrorMessage,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		AccessedAt:   time.Now(),
		AccessCount:  0,
		RelatedIDs:   []string{},
	}

	// 存储记录
	return aum.episodicMemory.Store(record)
}

// updateSemanticMemory 更新语义记忆（L3）
func (aum *AutoUpdateManager) updateSemanticMemory(result *TaskResult) error {
	if aum.semanticMemory == nil {
		return nil
	}

	// 根据任务结果类型更新不同的语义记忆

	// 1. 更新成功模式
	if result.Success && result.Importance >= aum.config.SuccessThreshold {
		if err := aum.updateSuccessPattern(result); err != nil {
			return fmt.Errorf("failed to update success pattern: %w", err)
		}
	}

	// 2. 更新失败模式
	if !result.Success && result.ErrorType != "" {
		if err := aum.updateFailurePattern(result); err != nil {
			return fmt.Errorf("failed to update failure pattern: %w", err)
		}
	}

	// 3. 更新学习事件
	if len(result.Insights) > 0 && result.Importance >= aum.config.LearningThreshold {
		if err := aum.updateLearningEvents(result); err != nil {
			return fmt.Errorf("failed to update learning events: %w", err)
		}
	}

	// 4. 更新记忆库内容
	if err := aum.updateMemoryContent(result); err != nil {
		return fmt.Errorf("failed to update memory content: %w", err)
	}

	return nil
}

// updateSuccessPattern 更新成功模式
func (aum *AutoUpdateManager) updateSuccessPattern(result *TaskResult) error {
	// 提取关键因素
	keyFactors := aum.extractKeyFactors(result)

	// 创建成功模式
	pattern := &SuccessPattern{
		PatternID:     fmt.Sprintf("success_%s_%d", result.TaskType, time.Now().UnixNano()),
		TaskType:      result.TaskType,
		Description:   result.Description,
		KeyFactors:    keyFactors,
		Conditions:    aum.extractConditions(result),
		Effectiveness: result.Importance,
		UsageCount:    1,
		LastUsed:      time.Now(),
		CreatedAt:     time.Now(),
	}

	return aum.semanticMemory.AddSuccessPattern(pattern)
}

// updateFailurePattern 更新失败模式
func (aum *AutoUpdateManager) updateFailurePattern(result *TaskResult) error {
	// 分析根本原因
	rootCause := aum.analyzeRootCause(result)

	// 创建失败模式
	pattern := &FailurePattern{
		PatternID:    fmt.Sprintf("failure_%s_%d", result.TaskType, time.Now().UnixNano()),
		TaskType:     result.TaskType,
		Description:  result.Description,
		ErrorType:    result.ErrorType,
		RootCause:    rootCause,
		Conditions:   aum.extractConditions(result),
		Frequency:    1,
		LastOccurred: time.Now(),
		CreatedAt:    time.Now(),
	}

	return aum.semanticMemory.AddFailurePattern(pattern)
}

// updateLearningEvents 更新学习事件
func (aum *AutoUpdateManager) updateLearningEvents(result *TaskResult) error {
	for _, insight := range result.Insights {
		// 确定影响类型
		impact := "neutral"
		if result.Success {
			impact = "positive"
		} else if result.ErrorType != "" {
			impact = "negative"
		}

		// 计算置信度
		confidence := result.Importance

		// 创建学习事件
		event := &LearningEvent{
			EventID:      fmt.Sprintf("learn_%s_%d", result.TaskType, time.Now().UnixNano()),
			TaskType:     result.TaskType,
			Description:  result.Description,
			Insight:      insight,
			Impact:       impact,
			Confidence:   confidence,
			Tags:         result.Tags,
			RelatedTasks: []string{result.TaskID},
			CreatedAt:    time.Now(),
		}

		if err := aum.semanticMemory.AddLearningEvent(event); err != nil {
			return err
		}
	}

	return nil
}

// updateMemoryContent 更新记忆库内容
func (aum *AutoUpdateManager) updateMemoryContent(result *TaskResult) error {
	// 构建更新内容
	var content strings.Builder

	if result.Success {
		content.WriteString(fmt.Sprintf("成功完成%s任务: %s\n", result.TaskType, result.Description))
		content.WriteString(fmt.Sprintf("- 重要性: %.0f%%\n", result.Importance*100))
		content.WriteString(fmt.Sprintf("- 耗时: %v\n", result.Duration))

		if len(result.Tags) > 0 {
			content.WriteString(fmt.Sprintf("- 标签: %s\n", strings.Join(result.Tags, ", ")))
		}

		if len(result.Insights) > 0 {
			content.WriteString("- 学习要点:\n")
			for _, insight := range result.Insights {
				content.WriteString(fmt.Sprintf("  - %s\n", insight))
			}
		}
	} else {
		content.WriteString(fmt.Sprintf("执行%s任务失败: %s\n", result.TaskType, result.Description))
		content.WriteString(fmt.Sprintf("- 错误类型: %s\n", result.ErrorType))
		content.WriteString(fmt.Sprintf("- 错误信息: %s\n", result.ErrorMessage))
		content.WriteString(fmt.Sprintf("- 重要性: %.0f%%\n", result.Importance*100))

		if len(result.Insights) > 0 {
			content.WriteString("- 教训总结:\n")
			for _, insight := range result.Insights {
				content.WriteString(fmt.Sprintf("  - %s\n", insight))
			}
		}
	}

	// 追加到记忆库
	section := "经验记录"
	if result.Success {
		section = "成功经验"
	} else {
		section = "失败教训"
	}

	return aum.semanticMemory.AppendMemory(section, content.String())
}

// updateWorkingMemory 更新工作记忆（L1）
func (aum *AutoUpdateManager) updateWorkingMemory(result *TaskResult) error {
	if aum.workingMemory == nil {
		return nil
	}

	// 更新会话上下文
	sessionID := result.SessionID
	if sessionID == "" {
		sessionID = "default"
	}

	// 获取或创建会话
	ctx := context.Background()
	session := aum.workingMemory.GetOrCreateSession(sessionID, ctx)

	// 存储任务结果摘要
	summary := map[string]any{
		"last_task_type":     result.TaskType,
		"last_task_result":   result.Success,
		"last_task_time":     time.Now(),
		"last_task_duration": result.Duration.String(),
		"task_importance":    result.Importance,
	}

	// 如果有错误，存储错误信息
	if !result.Success && result.ErrorType != "" {
		summary["last_error_type"] = result.ErrorType
		summary["last_error_message"] = result.ErrorMessage
	}

	// 更新工作记忆
	for key, value := range summary {
		session.Set(key, value)
	}

	return nil
}

// extractKeyFactors 提取关键因素
func (aum *AutoUpdateManager) extractKeyFactors(result *TaskResult) []string {
	var factors []string

	// 从元数据中提取关键因素
	if metadata, ok := result.Metadata["key_factors"]; ok {
		if factorsList, ok := metadata.([]string); ok {
			factors = append(factors, factorsList...)
		}
	}

	// 从标签中提取
	factors = append(factors, result.Tags...)

	// 从内容中提取关键词
	contentWords := strings.Fields(result.Content)
	for _, word := range contentWords {
		if len(word) > 4 && !aum.isCommonWord(word) {
			factors = append(factors, word)
		}
	}

	// 去重
	return aum.removeDuplicates(factors)
}

// extractConditions 提取条件
func (aum *AutoUpdateManager) extractConditions(result *TaskResult) map[string]any {
	conditions := make(map[string]any)

	// 添加基本条件
	conditions["task_type"] = result.TaskType
	conditions["category"] = result.Category
	conditions["importance"] = result.Importance
	conditions["duration"] = result.Duration.String()

	// 添加标签条件
	if len(result.Tags) > 0 {
		conditions["tags"] = result.Tags
	}

	// 添加元数据中的条件
	for key, value := range result.Metadata {
		if key != "key_factors" { // 避免重复
			conditions[key] = value
		}
	}

	return conditions
}

// analyzeRootCause 分析根本原因
func (aum *AutoUpdateManager) analyzeRootCause(result *TaskResult) string {
	// 根据错误类型分析根本原因
	switch result.ErrorType {
	case "timeout", "deadline_exceeded":
		return "时间管理不足或任务复杂度估计不准确"
	case "resource_exhausted", "memory_limit":
		return "资源规划不足或优化不够"
	case "permission_denied", "access_denied":
		return "权限配置问题或安全限制"
	case "validation_error", "invalid_input":
		return "输入验证不充分或数据格式错误"
	case "network_error", "connection_failed":
		return "网络连接问题或服务不可用"
	case "logic_error", "bug":
		return "代码逻辑错误或边界条件处理不足"
	case "dependency_error", "missing_dependency":
		return "依赖管理问题或版本冲突"
	default:
		// 从错误信息中提取关键词
		errorLower := strings.ToLower(result.ErrorMessage)
		if strings.Contains(errorLower, "time") || strings.Contains(errorLower, "slow") {
			return "性能问题或响应时间过长"
		} else if strings.Contains(errorLower, "memory") || strings.Contains(errorLower, "space") {
			return "内存或存储空间不足"
		} else if strings.Contains(errorLower, "file") || strings.Contains(errorLower, "directory") {
			return "文件系统操作问题"
		} else if strings.Contains(errorLower, "syntax") || strings.Contains(errorLower, "parse") {
			return "语法错误或解析问题"
		} else {
			return "未知原因，需要进一步分析"
		}
	}
}

// isCommonWord 判断是否为常见词
func (aum *AutoUpdateManager) isCommonWord(word string) bool {
	commonWords := map[string]bool{
		"the": true, "and": true, "for": true, "with": true, "this": true,
		"that": true, "have": true, "from": true, "what": true, "when": true,
		"where": true, "which": true, "will": true, "would": true, "could": true,
		"should": true, "about": true, "after": true, "before": true, "between": true,
		"into": true, "through": true, "during": true, "without": true,
	}

	return commonWords[strings.ToLower(word)]
}

// removeDuplicates 去除重复项
func (aum *AutoUpdateManager) removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// BatchUpdate 批量更新记忆
func (aum *AutoUpdateManager) BatchUpdate(ctx context.Context, results []*TaskResult) error {
	count := 0
	for _, result := range results {
		if count >= aum.config.MaxRecordsPerUpdate {
			break
		}

		if err := aum.RecordTaskResult(ctx, result); err != nil {
			return fmt.Errorf("failed to record task result %s: %w", result.TaskID, err)
		}

		count++
	}

	return nil
}

// GetUpdateStats 获取更新统计
func (aum *AutoUpdateManager) GetUpdateStats() map[string]any {
	stats := make(map[string]any)

	// 获取L2统计
	if aum.episodicMemory != nil {
		episodicStats, err := aum.episodicMemory.GetStats()
		if err == nil {
			stats["episodic_memory"] = episodicStats
		}
	}

	// 获取配置
	stats["config"] = aum.config

	return stats
}

// CleanupOldRecords 清理旧记录
func (aum *AutoUpdateManager) CleanupOldRecords(maxAgeDays int, minImportance float64) (int, error) {
	if aum.episodicMemory == nil {
		return 0, nil
	}

	return aum.episodicMemory.Cleanup(maxAgeDays, minImportance)
}
