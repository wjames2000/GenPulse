package skills

import (
	"fmt"
	"time"
)

// TriggerType 触发器类型
type TriggerType string

const (
	// 基于工具调用次数的触发器
	TriggerToolUsageCount TriggerType = "tool_usage_count"

	// 基于错误类型的触发器
	TriggerErrorType TriggerType = "error_type"

	// 基于任务复杂度的触发器
	TriggerTaskComplexity TriggerType = "task_complexity"

	// 基于执行时间的触发器
	TriggerExecutionTime TriggerType = "execution_time"

	// 基于成功率的触发器
	TriggerSuccessRate TriggerType = "success_rate"

	// 手动触发器
	TriggerManual TriggerType = "manual"
)

// TriggerCondition 触发条件
type TriggerCondition struct {
	Type      TriggerType `json:"type" yaml:"type"`
	Threshold any         `json:"threshold" yaml:"threshold"` // 阈值，类型根据TriggerType变化
	Operator  string      `json:"operator" yaml:"operator"`   // >, <, >=, <=, ==, !=
}

// TriggerResult 触发结果
type TriggerResult struct {
	Triggered  bool        `json:"triggered" yaml:"triggered"`
	Type       TriggerType `json:"type" yaml:"type"`
	Reason     string      `json:"reason" yaml:"reason"`
	Confidence float64     `json:"confidence" yaml:"confidence"` // 触发置信度 0-1
	Data       any         `json:"data" yaml:"data"`             // 触发数据
	Timestamp  time.Time   `json:"timestamp" yaml:"timestamp"`
}

// TaskExecutionRecord 任务执行记录
type TaskExecutionRecord struct {
	TaskID        string         `json:"task_id" yaml:"task_id"`
	TaskType      string         `json:"task_type" yaml:"task_type"`
	Description   string         `json:"description" yaml:"description"`
	Complexity    string         `json:"complexity" yaml:"complexity"` // simple, medium, complex
	StartTime     time.Time      `json:"start_time" yaml:"start_time"`
	EndTime       time.Time      `json:"end_time" yaml:"end_time"`
	Success       bool           `json:"success" yaml:"success"`
	ErrorType     string         `json:"error_type,omitempty" yaml:"error_type,omitempty"`
	ErrorMessage  string         `json:"error_message,omitempty" yaml:"error_message,omitempty"`
	ToolUsage     map[string]int `json:"tool_usage" yaml:"tool_usage"`         // 工具名 -> 使用次数
	AgentInvolved []string       `json:"agent_involved" yaml:"agent_involved"` // 参与的Agent
	Steps         []TaskStep     `json:"steps" yaml:"steps"`                   // 执行步骤
	Context       map[string]any `json:"context" yaml:"context"`               // 执行上下文
	Output        map[string]any `json:"output" yaml:"output"`                 // 输出结果
	Metadata      map[string]any `json:"metadata" yaml:"metadata"`             // 元数据
}

// TaskStep 任务步骤
type TaskStep struct {
	StepID     string         `json:"step_id" yaml:"step_id"`
	Order      int            `json:"order" yaml:"order"`
	Action     string         `json:"action" yaml:"action"`
	Tool       string         `json:"tool,omitempty" yaml:"tool,omitempty"`
	Parameters map[string]any `json:"parameters" yaml:"parameters"`
	StartTime  time.Time      `json:"start_time" yaml:"start_time"`
	EndTime    time.Time      `json:"end_time" yaml:"end_time"`
	Success    bool           `json:"success" yaml:"success"`
	Error      string         `json:"error,omitempty" yaml:"error,omitempty"`
	Output     map[string]any `json:"output" yaml:"output"`
}

// TriggerManager 触发器管理器
type TriggerManager struct {
	conditions []TriggerCondition
	registry   *Registry
	thresholds map[TriggerType]ThresholdConfig
}

// ThresholdConfig 阈值配置
type ThresholdConfig struct {
	DefaultThreshold any    `json:"default_threshold" yaml:"default_threshold"`
	Description      string `json:"description" yaml:"description"`
	Unit             string `json:"unit,omitempty" yaml:"unit,omitempty"`
}

// NewTriggerManager 创建触发器管理器
func NewTriggerManager(registry *Registry) *TriggerManager {
	tm := &TriggerManager{
		registry: registry,
		conditions: []TriggerCondition{
			// 默认触发器配置
			{
				Type:      TriggerToolUsageCount,
				Threshold: 3,
				Operator:  ">=",
			},
			{
				Type:      TriggerErrorType,
				Threshold: "permission_denied",
				Operator:  "==",
			},
			{
				Type:      TriggerTaskComplexity,
				Threshold: "complex",
				Operator:  "==",
			},
			{
				Type:      TriggerExecutionTime,
				Threshold: time.Minute * 5,
				Operator:  ">",
			},
			{
				Type:      TriggerSuccessRate,
				Threshold: 0.8,
				Operator:  ">",
			},
		},
		thresholds: getDefaultThresholds(),
	}

	return tm
}

// CheckTriggers 检查是否触发技能生成
func (tm *TriggerManager) CheckTriggers(record *TaskExecutionRecord) []TriggerResult {
	var results []TriggerResult

	// 检查每个触发器条件
	for _, condition := range tm.conditions {
		result := tm.checkCondition(condition, record)
		if result.Triggered {
			results = append(results, result)
		}
	}

	return results
}

// checkCondition 检查单个条件
func (tm *TriggerManager) checkCondition(condition TriggerCondition, record *TaskExecutionRecord) TriggerResult {
	result := TriggerResult{
		Type:      condition.Type,
		Timestamp: time.Now(),
	}

	switch condition.Type {
	case TriggerToolUsageCount:
		result = tm.checkToolUsage(condition, record)
	case TriggerErrorType:
		result = tm.checkErrorType(condition, record)
	case TriggerTaskComplexity:
		result = tm.checkTaskComplexity(condition, record)
	case TriggerExecutionTime:
		result = tm.checkExecutionTime(condition, record)
	case TriggerSuccessRate:
		result = tm.checkSuccessRate(condition, record)
	case TriggerManual:
		// 手动触发器总是返回未触发，由外部控制
		result.Triggered = false
		result.Reason = "Manual trigger requires explicit call"
	}

	return result
}

// checkToolUsage 检查工具使用次数
func (tm *TriggerManager) checkToolUsage(condition TriggerCondition, record *TaskExecutionRecord) TriggerResult {
	result := TriggerResult{
		Type:      TriggerToolUsageCount,
		Timestamp: time.Now(),
	}

	threshold, ok := condition.Threshold.(int)
	if !ok {
		result.Triggered = false
		result.Reason = fmt.Sprintf("Invalid threshold type for tool usage: %T", condition.Threshold)
		return result
	}

	// 检查是否有工具使用次数超过阈值
	for toolName, usageCount := range record.ToolUsage {
		if compare(usageCount, threshold, condition.Operator) {
			result.Triggered = true
			result.Reason = fmt.Sprintf("Tool '%s' used %d times (threshold: %d %s)",
				toolName, usageCount, threshold, condition.Operator)
			result.Confidence = calculateConfidence(float64(usageCount), float64(threshold), 10.0) // 最大10次
			result.Data = map[string]any{
				"tool":        toolName,
				"usage_count": usageCount,
				"threshold":   threshold,
			}
			break
		}
	}

	if !result.Triggered {
		result.Reason = fmt.Sprintf("No tool usage exceeds threshold %d", threshold)
	}

	return result
}

// checkErrorType 检查错误类型
func (tm *TriggerManager) checkErrorType(condition TriggerCondition, record *TaskExecutionRecord) TriggerResult {
	result := TriggerResult{
		Type:      TriggerErrorType,
		Timestamp: time.Now(),
	}

	if !record.Success && record.ErrorType != "" {
		threshold, ok := condition.Threshold.(string)
		if !ok {
			result.Triggered = false
			result.Reason = fmt.Sprintf("Invalid threshold type for error type: %T", condition.Threshold)
			return result
		}

		if compare(record.ErrorType, threshold, condition.Operator) {
			result.Triggered = true
			result.Reason = fmt.Sprintf("Error type '%s' matches trigger condition", record.ErrorType)
			result.Confidence = 0.9 // 错误类型匹配置信度高
			result.Data = map[string]any{
				"error_type":    record.ErrorType,
				"error_message": record.ErrorMessage,
			}
		} else {
			result.Reason = fmt.Sprintf("Error type '%s' does not match threshold '%s'",
				record.ErrorType, threshold)
		}
	} else {
		result.Reason = "Task succeeded or no error type specified"
	}

	return result
}

// checkTaskComplexity 检查任务复杂度
func (tm *TriggerManager) checkTaskComplexity(condition TriggerCondition, record *TaskExecutionRecord) TriggerResult {
	result := TriggerResult{
		Type:      TriggerTaskComplexity,
		Timestamp: time.Now(),
	}

	threshold, ok := condition.Threshold.(string)
	if !ok {
		result.Triggered = false
		result.Reason = fmt.Sprintf("Invalid threshold type for task complexity: %T", condition.Threshold)
		return result
	}

	// 定义复杂度等级
	complexityLevels := map[string]int{
		"simple":  1,
		"medium":  2,
		"complex": 3,
	}

	currentLevel, currentOk := complexityLevels[record.Complexity]
	thresholdLevel, thresholdOk := complexityLevels[threshold]

	if currentOk && thresholdOk {
		if compare(currentLevel, thresholdLevel, condition.Operator) {
			result.Triggered = true
			result.Reason = fmt.Sprintf("Task complexity '%s' meets condition", record.Complexity)
			result.Confidence = 0.8
			result.Data = map[string]any{
				"complexity": record.Complexity,
				"level":      currentLevel,
			}
		} else {
			result.Reason = fmt.Sprintf("Task complexity '%s' does not meet condition", record.Complexity)
		}
	} else {
		result.Reason = fmt.Sprintf("Invalid complexity value: current=%s, threshold=%s",
			record.Complexity, threshold)
	}

	return result
}

// checkExecutionTime 检查执行时间
func (tm *TriggerManager) checkExecutionTime(condition TriggerCondition, record *TaskExecutionRecord) TriggerResult {
	result := TriggerResult{
		Type:      TriggerExecutionTime,
		Timestamp: time.Now(),
	}

	threshold, ok := condition.Threshold.(time.Duration)
	if !ok {
		// 尝试从数字转换
		if num, ok := condition.Threshold.(int); ok {
			threshold = time.Duration(num) * time.Second
		} else if num, ok := condition.Threshold.(float64); ok {
			threshold = time.Duration(num) * time.Second
		} else {
			result.Triggered = false
			result.Reason = fmt.Sprintf("Invalid threshold type for execution time: %T", condition.Threshold)
			return result
		}
	}

	executionTime := record.EndTime.Sub(record.StartTime)
	if compare(executionTime, threshold, condition.Operator) {
		result.Triggered = true
		result.Reason = fmt.Sprintf("Execution time %v exceeds threshold %v", executionTime, threshold)
		result.Confidence = calculateConfidence(float64(executionTime), float64(threshold), float64(threshold*2))
		result.Data = map[string]any{
			"execution_time": executionTime,
			"threshold":      threshold,
		}
	} else {
		result.Reason = fmt.Sprintf("Execution time %v does not exceed threshold %v", executionTime, threshold)
	}

	return result
}

// checkSuccessRate 检查成功率
func (tm *TriggerManager) checkSuccessRate(condition TriggerCondition, record *TaskExecutionRecord) TriggerResult {
	result := TriggerResult{
		Type:      TriggerSuccessRate,
		Timestamp: time.Now(),
	}

	// 这里需要从历史记录中计算成功率
	// 简化实现：只检查当前任务是否成功
	threshold, ok := condition.Threshold.(float64)
	if !ok {
		result.Triggered = false
		result.Reason = fmt.Sprintf("Invalid threshold type for success rate: %T", condition.Threshold)
		return result
	}

	// 获取相关技能的历史成功率
	// 这里简化实现，实际应该查询历史记录
	successRate := 0.0
	if record.Success {
		successRate = 1.0
	}

	if compare(successRate, threshold, condition.Operator) {
		result.Triggered = true
		result.Reason = fmt.Sprintf("Success rate %.2f meets condition", successRate)
		result.Confidence = 0.7
		result.Data = map[string]any{
			"success_rate": successRate,
			"threshold":    threshold,
		}
	} else {
		result.Reason = fmt.Sprintf("Success rate %.2f does not meet condition", successRate)
	}

	return result
}

// AddCondition 添加触发条件
func (tm *TriggerManager) AddCondition(condition TriggerCondition) {
	tm.conditions = append(tm.conditions, condition)
}

// RemoveCondition 移除触发条件
func (tm *TriggerManager) RemoveCondition(index int) error {
	if index < 0 || index >= len(tm.conditions) {
		return fmt.Errorf("condition index out of range")
	}

	tm.conditions = append(tm.conditions[:index], tm.conditions[index+1:]...)
	return nil
}

// GetConditions 获取所有触发条件
func (tm *TriggerManager) GetConditions() []TriggerCondition {
	return tm.conditions
}

// SetThreshold 设置阈值配置
func (tm *TriggerManager) SetThreshold(triggerType TriggerType, config ThresholdConfig) {
	tm.thresholds[triggerType] = config
}

// GetThreshold 获取阈值配置
func (tm *TriggerManager) GetThreshold(triggerType TriggerType) (ThresholdConfig, bool) {
	config, exists := tm.thresholds[triggerType]
	return config, exists
}

// ShouldExtractSkill 判断是否应该提取技能
func (tm *TriggerManager) ShouldExtractSkill(record *TaskExecutionRecord) (bool, []TriggerResult) {
	results := tm.CheckTriggers(record)

	// 如果有任何触发器被触发，则应该提取技能
	shouldExtract := len(results) > 0

	// 计算总体置信度
	if shouldExtract {
		// 取最高置信度
		maxConfidence := 0.0
		for _, result := range results {
			if result.Confidence > maxConfidence {
				maxConfidence = result.Confidence
			}
		}

		// 如果最高置信度低于阈值，可能不提取
		if maxConfidence < 0.5 {
			shouldExtract = false
		}
	}

	return shouldExtract, results
}

// compare 通用比较函数
func compare(a, b any, operator string) bool {
	switch aVal := a.(type) {
	case int:
		bVal, ok := b.(int)
		if !ok {
			return false
		}
		switch operator {
		case ">":
			return aVal > bVal
		case "<":
			return aVal < bVal
		case ">=":
			return aVal >= bVal
		case "<=":
			return aVal <= bVal
		case "==":
			return aVal == bVal
		case "!=":
			return aVal != bVal
		}
	case float64:
		bVal, ok := b.(float64)
		if !ok {
			return false
		}
		switch operator {
		case ">":
			return aVal > bVal
		case "<":
			return aVal < bVal
		case ">=":
			return aVal >= bVal
		case "<=":
			return aVal <= bVal
		case "==":
			return aVal == bVal
		case "!=":
			return aVal != bVal
		}
	case string:
		bVal, ok := b.(string)
		if !ok {
			return false
		}
		switch operator {
		case "==":
			return aVal == bVal
		case "!=":
			return aVal != bVal
		}
	case time.Duration:
		bVal, ok := b.(time.Duration)
		if !ok {
			return false
		}
		switch operator {
		case ">":
			return aVal > bVal
		case "<":
			return aVal < bVal
		case ">=":
			return aVal >= bVal
		case "<=":
			return aVal <= bVal
		case "==":
			return aVal == bVal
		case "!=":
			return aVal != bVal
		}
	}

	return false
}

// calculateConfidence 计算置信度
func calculateConfidence(value, threshold, maxValue float64) float64 {
	if maxValue <= 0 {
		return 0.5
	}

	// 归一化到0-1范围
	normalized := (value - threshold) / (maxValue - threshold)

	// 限制在0-1之间
	if normalized < 0 {
		return 0.3
	}
	if normalized > 1 {
		return 1.0
	}

	// 基础置信度 + 归一化增量
	return 0.5 + normalized*0.5
}

// getDefaultThresholds 获取默认阈值配置
func getDefaultThresholds() map[TriggerType]ThresholdConfig {
	return map[TriggerType]ThresholdConfig{
		TriggerToolUsageCount: {
			DefaultThreshold: 3,
			Description:      "工具使用次数阈值",
			Unit:             "次",
		},
		TriggerErrorType: {
			DefaultThreshold: "permission_denied",
			Description:      "触发技能生成的错误类型",
		},
		TriggerTaskComplexity: {
			DefaultThreshold: "complex",
			Description:      "任务复杂度阈值",
		},
		TriggerExecutionTime: {
			DefaultThreshold: time.Minute * 5,
			Description:      "执行时间阈值",
			Unit:             "duration",
		},
		TriggerSuccessRate: {
			DefaultThreshold: 0.8,
			Description:      "成功率阈值",
			Unit:             "ratio",
		},
	}
}

// ManualTrigger 手动触发技能生成
func (tm *TriggerManager) ManualTrigger(record *TaskExecutionRecord, reason string) TriggerResult {
	return TriggerResult{
		Triggered:  true,
		Type:       TriggerManual,
		Reason:     reason,
		Confidence: 1.0,
		Data:       record,
		Timestamp:  time.Now(),
	}
}
