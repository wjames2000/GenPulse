package pipeline

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"GenPulse/internal/agents"
	"GenPulse/internal/utils"
)

// ErrorType 错误类型
type ErrorType string

const (
	ErrorTypeAgentUnavailable  ErrorType = "agent_unavailable"
	ErrorTypeExecutionFailed   ErrorType = "execution_failed"
	ErrorTypeTimeout           ErrorType = "timeout"
	ErrorTypeValidationFailed  ErrorType = "validation_failed"
	ErrorTypeResourceExhausted ErrorType = "resource_exhausted"
	ErrorTypeNetworkError      ErrorType = "network_error"
)

// ErrorSeverity 错误严重程度
type ErrorSeverity string

const (
	ErrorSeverityLow      ErrorSeverity = "low"
	ErrorSeverityMedium   ErrorSeverity = "medium"
	ErrorSeverityHigh     ErrorSeverity = "high"
	ErrorSeverityCritical ErrorSeverity = "critical"
)

// PipelineError 流水线错误
type PipelineError struct {
	Type        ErrorType     `json:"type"`
	Severity    ErrorSeverity `json:"severity"`
	Stage       string        `json:"stage"`
	AgentID     string        `json:"agent_id,omitempty"`
	Message     string        `json:"message"`
	OriginalErr error         `json:"original_error,omitempty"`
	Timestamp   time.Time     `json:"timestamp"`
	RetryCount  int           `json:"retry_count"`
}

// ErrorHandler 错误处理器
type ErrorHandler struct {
	maxRetries         int
	retryDelay         time.Duration
	fallbackStrategies map[ErrorType]FallbackStrategy
	errorHistory       []PipelineError
	mu                 sync.RWMutex
}

// FallbackStrategy 降级策略
type FallbackStrategy struct {
	Enabled     bool                   `json:"enabled"`
	Strategy    string                 `json:"strategy"` // retry, skip, use_alternative, manual_fallback
	Alternative string                 `json:"alternative,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler(maxRetries int, retryDelay time.Duration) *ErrorHandler {
	handler := &ErrorHandler{
		maxRetries:   maxRetries,
		retryDelay:   retryDelay,
		errorHistory: make([]PipelineError, 0),
		fallbackStrategies: map[ErrorType]FallbackStrategy{
			ErrorTypeAgentUnavailable: {
				Enabled:     true,
				Strategy:    "use_alternative",
				Alternative: "simple_agent",
				Config: map[string]interface{}{
					"max_wait_time": "30s",
				},
			},
			ErrorTypeExecutionFailed: {
				Enabled:  true,
				Strategy: "retry",
				Config: map[string]interface{}{
					"max_retries":    3,
					"backoff_factor": 2,
				},
			},
			ErrorTypeTimeout: {
				Enabled:  true,
				Strategy: "retry",
				Config: map[string]interface{}{
					"max_retries":      2,
					"timeout_increase": "50%",
				},
			},
			ErrorTypeValidationFailed: {
				Enabled:  true,
				Strategy: "manual_fallback",
				Config: map[string]interface{}{
					"require_human_review": true,
				},
			},
			ErrorTypeResourceExhausted: {
				Enabled:  true,
				Strategy: "skip",
				Config: map[string]interface{}{
					"skip_if_non_critical": true,
				},
			},
			ErrorTypeNetworkError: {
				Enabled:  true,
				Strategy: "retry",
				Config: map[string]interface{}{
					"max_retries":         5,
					"exponential_backoff": true,
				},
			},
		},
	}

	return handler
}

// HandleError 处理错误
func (eh *ErrorHandler) HandleError(ctx context.Context, stage string, agentID string, err error, retryCount int) (shouldRetry bool, fallbackAction string, waitTime time.Duration) {
	// 分析错误类型和严重程度
	errorType, severity := eh.analyzeError(err)

	pipelineErr := PipelineError{
		Type:        errorType,
		Severity:    severity,
		Stage:       stage,
		AgentID:     agentID,
		Message:     err.Error(),
		OriginalErr: err,
		Timestamp:   time.Now(),
		RetryCount:  retryCount,
	}

	// 记录错误
	eh.recordError(pipelineErr)

	// 获取对应的降级策略
	strategy, exists := eh.fallbackStrategies[errorType]
	if !exists || !strategy.Enabled {
		// 没有可用的降级策略，不重试
		return false, "abort", 0
	}

	// 根据策略决定行动
	switch strategy.Strategy {
	case "retry":
		if retryCount < eh.maxRetries {
			// 计算等待时间（指数退避）
			waitTime = eh.calculateWaitTime(retryCount)
			utils.Warn("错误处理: 阶段 %s 将在 %v 后重试 (尝试 %d/%d)",
				stage, waitTime, retryCount+1, eh.maxRetries)
			return true, "retry", waitTime
		}
		// 重试次数用尽，尝试其他策略
		return eh.handleRetryExhausted(pipelineErr, strategy)

	case "skip":
		if eh.canSkipStage(stage, severity) {
			utils.Warn("错误处理: 跳过阶段 %s (严重程度: %s)", stage, severity)
			return false, "skip", 0
		}
		return false, "abort", 0

	case "use_alternative":
		alternative := strategy.Alternative
		if alternative != "" {
			utils.Warn("错误处理: 使用替代方案 %s 替换 %s", alternative, agentID)
			return false, "use_alternative:" + alternative, 0
		}
		return false, "abort", 0

	case "manual_fallback":
		utils.Error("错误处理: 需要人工干预 - 阶段 %s, 错误: %v", stage, err)
		return false, "manual_intervention_required", 0

	default:
		return false, "abort", 0
	}
}

// analyzeError 分析错误类型和严重程度
func (eh *ErrorHandler) analyzeError(err error) (ErrorType, ErrorSeverity) {
	errStr := err.Error()

	// 检查错误类型
	var errorType ErrorType
	var severity ErrorSeverity

	switch {
	case strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded"):
		errorType = ErrorTypeTimeout
		severity = ErrorSeverityMedium

	case strings.Contains(errStr, "agent not found") || strings.Contains(errStr, "unavailable"):
		errorType = ErrorTypeAgentUnavailable
		severity = ErrorSeverityHigh

	case strings.Contains(errStr, "validation") || strings.Contains(errStr, "invalid"):
		errorType = ErrorTypeValidationFailed
		severity = ErrorSeverityHigh

	case strings.Contains(errStr, "resource") || strings.Contains(errStr, "memory") || strings.Contains(errStr, "disk"):
		errorType = ErrorTypeResourceExhausted
		severity = ErrorSeverityCritical

	case strings.Contains(errStr, "network") || strings.Contains(errStr, "connection") || strings.Contains(errStr, "HTTP"):
		errorType = ErrorTypeNetworkError
		severity = ErrorSeverityMedium

	default:
		errorType = ErrorTypeExecutionFailed
		severity = ErrorSeverityMedium
	}

	return errorType, severity
}

// calculateWaitTime 计算等待时间（指数退避）
func (eh *ErrorHandler) calculateWaitTime(retryCount int) time.Duration {
	// 基础等待时间 * 2^重试次数
	waitTime := eh.retryDelay * time.Duration(1<<uint(retryCount))

	// 设置上限为5分钟
	maxWait := 5 * time.Minute
	if waitTime > maxWait {
		waitTime = maxWait
	}

	return waitTime
}

// handleRetryExhausted 处理重试次数用尽的情况
func (eh *ErrorHandler) handleRetryExhausted(err PipelineError, strategy FallbackStrategy) (bool, string, time.Duration) {
	// 检查是否有备用策略
	if alternative, ok := strategy.Config["fallback_strategy"]; ok {
		fallbackStr := fmt.Sprintf("%v", alternative)
		utils.Warn("重试次数用尽，使用备用策略: %s", fallbackStr)
		return false, "fallback:" + fallbackStr, 0
	}

	// 检查是否可以跳过
	if skipIfNonCritical, ok := strategy.Config["skip_if_non_critical"].(bool); ok && skipIfNonCritical {
		if err.Severity != ErrorSeverityCritical {
			utils.Warn("重试次数用尽，跳过非关键阶段")
			return false, "skip_after_exhausted", 0
		}
	}

	utils.Error("重试次数用尽且无备用策略，中止执行")
	return false, "abort_after_exhausted", 0
}

// canSkipStage 检查是否可以跳过阶段
func (eh *ErrorHandler) canSkipStage(stage string, severity ErrorSeverity) bool {
	// 关键阶段不能跳过
	criticalStages := []string{
		"requirements_analysis",
		"architecture_design",
		"project_validation",
	}

	for _, criticalStage := range criticalStages {
		if stage == criticalStage {
			return false
		}
	}

	// 严重错误不能跳过
	if severity == ErrorSeverityCritical {
		return false
	}

	return true
}

// recordError 记录错误
func (eh *ErrorHandler) recordError(err PipelineError) {
	eh.mu.Lock()
	defer eh.mu.Unlock()

	eh.errorHistory = append(eh.errorHistory, err)

	// 保持错误历史大小
	if len(eh.errorHistory) > 1000 {
		eh.errorHistory = eh.errorHistory[500:] // 保留最近500个错误
	}

	// 记录日志
	utils.Error("流水线错误: 阶段=%s, 类型=%s, 严重程度=%s, 消息=%s",
		err.Stage, err.Type, err.Severity, err.Message)
}

// GetErrorHistory 获取错误历史
func (eh *ErrorHandler) GetErrorHistory() []PipelineError {
	eh.mu.RLock()
	defer eh.mu.RUnlock()

	// 返回副本
	history := make([]PipelineError, len(eh.errorHistory))
	copy(history, eh.errorHistory)
	return history
}

// GetErrorStats 获取错误统计
func (eh *ErrorHandler) GetErrorStats() map[string]interface{} {
	eh.mu.RLock()
	defer eh.mu.RUnlock()

	stats := map[string]interface{}{
		"total_errors": len(eh.errorHistory),
		"by_type":      make(map[string]int),
		"by_severity":  make(map[string]int),
		"by_stage":     make(map[string]int),
	}

	for _, err := range eh.errorHistory {
		// 按类型统计
		typeKey := string(err.Type)
		stats["by_type"].(map[string]int)[typeKey]++

		// 按严重程度统计
		severityKey := string(err.Severity)
		stats["by_severity"].(map[string]int)[severityKey]++

		// 按阶段统计
		stageKey := err.Stage
		stats["by_stage"].(map[string]int)[stageKey]++
	}

	return stats
}

// SetFallbackStrategy 设置降级策略
func (eh *ErrorHandler) SetFallbackStrategy(errorType ErrorType, strategy FallbackStrategy) {
	eh.mu.Lock()
	defer eh.mu.Unlock()

	eh.fallbackStrategies[errorType] = strategy
}

// GetFallbackStrategy 获取降级策略
func (eh *ErrorHandler) GetFallbackStrategy(errorType ErrorType) (FallbackStrategy, bool) {
	eh.mu.RLock()
	defer eh.mu.RUnlock()

	strategy, exists := eh.fallbackStrategies[errorType]
	return strategy, exists
}

// ExecuteWithRetry 带重试的执行
func (eh *ErrorHandler) ExecuteWithRetry(ctx context.Context, stage string, agentID string, executeFunc func() (interface{}, error)) (interface{}, error) {
	var lastErr error

	for retry := 0; retry <= eh.maxRetries; retry++ {
		// 执行函数
		result, err := executeFunc()

		if err == nil {
			// 执行成功
			if retry > 0 {
				utils.Info("重试成功: 阶段 %s (尝试 %d)", stage, retry+1)
			}
			return result, nil
		}

		lastErr = err

		// 如果是最后一次尝试，直接返回错误
		if retry == eh.maxRetries {
			break
		}

		// 处理错误，决定是否重试
		shouldRetry, action, waitTime := eh.HandleError(ctx, stage, agentID, err, retry)

		if !shouldRetry {
			// 不重试，返回错误
			return nil, fmt.Errorf("%s: %w", action, err)
		}

		// 等待重试
		select {
		case <-time.After(waitTime):
			// 继续重试
			continue
		case <-ctx.Done():
			// 上下文取消
			return nil, fmt.Errorf("执行被取消: %w", ctx.Err())
		}
	}

	// 所有重试都失败
	return nil, fmt.Errorf("重试 %d 次后仍然失败: %w", eh.maxRetries, lastErr)
}

// CreateFallbackAgent 创建降级Agent
func (eh *ErrorHandler) CreateFallbackAgent(originalAgentID string, fallbackType string) (agents.Agent, error) {
	// 这里可以根据fallbackType创建不同的降级Agent
	// 例如: "simple_agent" -> 使用简化版Agent
	//       "basic_agent" -> 使用基础功能Agent

	// 这是一个示例实现，实际项目中需要根据具体需求实现
	switch fallbackType {
	case "simple_agent":
		// 返回一个简化版Agent
		// 实际实现中需要创建具体的Agent实例
		return nil, fmt.Errorf("降级Agent创建未实现: %s", fallbackType)
	default:
		return nil, fmt.Errorf("不支持的降级类型: %s", fallbackType)
	}
}
