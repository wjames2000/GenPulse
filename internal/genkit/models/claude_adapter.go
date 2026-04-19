package models

import (
	"context"
	"fmt"
	"time"

	"GenPulse/internal/utils"
)

// ClaudeAdapter Claude模型适配器
type ClaudeAdapter struct {
	config ModelConfig
	// client *anthropic.Client // TODO: 实际Anthropic客户端
}

// NewClaudeAdapter 创建Claude适配器
func NewClaudeAdapter(config ModelConfig) (ModelAdapter, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required for Claude adapter")
	}

	// 设置默认值
	if config.MaxTokens == 0 {
		config.MaxTokens = 4096
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.Timeout == 0 {
		config.Timeout = 30
	}

	adapter := &ClaudeAdapter{
		config: config,
	}

	// TODO: 初始化Anthropic客户端
	// client, err := anthropic.NewClient(config.APIKey)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to create Anthropic client: %w", err)
	// }
	// adapter.client = client

	utils.Info("创建Claude适配器: %s", config.Name)
	return adapter, nil
}

// Name 获取适配器名称
func (a *ClaudeAdapter) Name() string {
	return a.config.Name
}

// Type 获取模型类型
func (a *ClaudeAdapter) Type() ModelType {
	return ModelTypeClaude
}

// SupportsStreaming 是否支持流式
func (a *ClaudeAdapter) SupportsStreaming() bool {
	return true
}

// SupportsTools 是否支持工具调用
func (a *ClaudeAdapter) SupportsTools() bool {
	return true
}

// GetConfig 获取配置
func (a *ClaudeAdapter) GetConfig() ModelConfig {
	return a.config
}

// UpdateConfig 更新配置
func (a *ClaudeAdapter) UpdateConfig(config ModelConfig) error {
	a.config = config
	return nil
}

// Generate 生成文本
func (a *ClaudeAdapter) Generate(ctx context.Context, req ModelRequest) (*ModelResponse, error) {
	utils.Debug("Claude适配器生成文本: %s", a.config.Name)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.config.Timeout)*time.Second)
	defer cancel()

	// Claude通常使用消息格式
	var messages []ChatMessage
	if len(req.Messages) > 0 {
		messages = req.Messages
	} else {
		// 将普通prompt转换为消息格式
		messages = []ChatMessage{
			{Role: "user", Content: req.Prompt},
		}
	}

	// TODO: 实际调用Anthropic API
	// 这里先返回模拟响应

	// 构建响应内容
	responseContent := fmt.Sprintf("这是来自 %s (Claude) 的模拟响应。\n\n", a.config.Name)
	for _, msg := range messages {
		responseContent += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	response := &ModelResponse{
		Content:      responseContent,
		FinishReason: "stop",
		Usage: TokenUsage{
			PromptTokens:     calculateTokenCount(messages),
			CompletionTokens: 150,
			TotalTokens:      calculateTokenCount(messages) + 150,
		},
		Metadata: map[string]interface{}{
			"model":     a.config.Name,
			"provider":  a.config.Provider,
			"adapter":   "claude",
			"timestamp": time.Now().Unix(),
		},
	}

	return response, nil
}

// GenerateStream 流式生成文本
func (a *ClaudeAdapter) GenerateStream(ctx context.Context, req ModelRequest, callback func(*ModelResponse)) error {
	utils.Debug("Claude适配器流式生成文本: %s", a.config.Name)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.config.Timeout)*time.Second)
	defer cancel()

	// 处理消息
	var messages []ChatMessage
	if len(req.Messages) > 0 {
		messages = req.Messages
	} else {
		messages = []ChatMessage{
			{Role: "user", Content: req.Prompt},
		}
	}

	// TODO: 实际实现流式调用
	// 这里模拟流式响应

	// 模拟分块响应
	chunks := []string{
		"这是",
		"来自",
		a.config.Name,
		"(Claude)",
		"的",
		"流式",
		"响应。",
		"\n\n",
	}

	// 添加消息摘要
	for _, msg := range messages {
		chunks = append(chunks, fmt.Sprintf("%s:", msg.Role))
		chunks = append(chunks, msg.Content[:min(20, len(msg.Content))])
		if len(msg.Content) > 20 {
			chunks = append(chunks, "...")
		}
		chunks = append(chunks, "\n")
	}

	for i, chunk := range chunks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 模拟延迟
			time.Sleep(120 * time.Millisecond) // Claude通常响应较慢

			response := &ModelResponse{
				Content: chunk,
				Metadata: map[string]interface{}{
					"model":     a.config.Name,
					"chunk":     i + 1,
					"total":     len(chunks),
					"timestamp": time.Now().Unix(),
				},
			}

			callback(response)
		}
	}

	// 发送完成响应
	callback(&ModelResponse{
		FinishReason: "stop",
		Usage: TokenUsage{
			PromptTokens:     calculateTokenCount(messages),
			CompletionTokens: 80,
			TotalTokens:      calculateTokenCount(messages) + 80,
		},
		Metadata: map[string]interface{}{
			"model":     a.config.Name,
			"completed": true,
			"timestamp": time.Now().Unix(),
		},
	})

	return nil
}

// HealthCheck 健康检查
func (a *ClaudeAdapter) HealthCheck(ctx context.Context) error {
	// TODO: 实际健康检查

	// 检查API密钥
	if a.config.APIKey == "" {
		return fmt.Errorf("API key not configured")
	}

	// 检查模型名称
	if a.config.Name == "" {
		return fmt.Errorf("model name not configured")
	}

	utils.Debug("Claude适配器健康检查: %s", a.config.Name)

	// Claude模型通常以"claude-"开头
	if len(a.config.Name) < 7 || a.config.Name[:7] != "claude-" {
		utils.Warn("Claude模型名称可能不正确: %s", a.config.Name)
	}

	// 这里可以添加实际的API调用测试
	// 暂时返回成功
	return nil
}

// calculateTokenCount 计算消息的token数量（粗略估算）
func calculateTokenCount(messages []ChatMessage) int {
	total := 0
	for _, msg := range messages {
		// 粗略估算：4个字符约等于1个token
		total += len(msg.Role) + len(msg.Content) + len(msg.Name)
	}
	return total / 4
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
