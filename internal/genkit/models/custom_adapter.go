package models

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"GenPulse/internal/utils"
)

// CustomAdapter 自定义模型适配器
type CustomAdapter struct {
	config ModelConfig
	client *http.Client
}

// NewCustomAdapter 创建自定义适配器
func NewCustomAdapter(config ModelConfig) (ModelAdapter, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required for custom adapter")
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

	adapter := &CustomAdapter{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}

	utils.Info("创建自定义适配器: %s (URL: %s)", config.Name, config.BaseURL)
	return adapter, nil
}

// Name 获取适配器名称
func (a *CustomAdapter) Name() string {
	return a.config.Name
}

// Type 获取模型类型
func (a *CustomAdapter) Type() ModelType {
	return ModelTypeCustom
}

// SupportsStreaming 是否支持流式
func (a *CustomAdapter) SupportsStreaming() bool {
	// 自定义适配器可能支持流式，取决于后端实现
	// 这里假设不支持，除非明确配置
	return false
}

// SupportsTools 是否支持工具调用
func (a *CustomAdapter) SupportsTools() bool {
	// 自定义适配器可能支持工具调用
	// 这里假设不支持，除非明确配置
	return false
}

// GetConfig 获取配置
func (a *CustomAdapter) GetConfig() ModelConfig {
	return a.config
}

// UpdateConfig 更新配置
func (a *CustomAdapter) UpdateConfig(config ModelConfig) error {
	a.config = config
	// 更新HTTP客户端超时
	a.client.Timeout = time.Duration(config.Timeout) * time.Second
	return nil
}

// Generate 生成文本
func (a *CustomAdapter) Generate(ctx context.Context, req ModelRequest) (*ModelResponse, error) {
	utils.Debug("自定义适配器生成文本: %s", a.config.Name)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.config.Timeout)*time.Second)
	defer cancel()

	// 构建请求体（简化版本）
	requestBody := map[string]interface{}{
		"prompt":      req.Prompt,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}

	if len(req.Messages) > 0 {
		requestBody["messages"] = req.Messages
	}

	// TODO: 实际调用自定义API
	// 这里先返回模拟响应

	// 模拟API调用
	time.Sleep(500 * time.Millisecond) // 模拟网络延迟

	response := &ModelResponse{
		Content: fmt.Sprintf("这是来自自定义模型 %s 的响应。\n\nAPI端点: %s\n\n输入: %s",
			a.config.Name, a.config.BaseURL, req.Prompt),
		FinishReason: "stop",
		Usage: TokenUsage{
			PromptTokens:     len(req.Prompt) / 4,
			CompletionTokens: 180,
			TotalTokens:      len(req.Prompt)/4 + 180,
		},
		Metadata: map[string]interface{}{
			"model":     a.config.Name,
			"provider":  a.config.Provider,
			"adapter":   "custom",
			"base_url":  a.config.BaseURL,
			"custom":    true,
			"timestamp": time.Now().Unix(),
		},
	}

	return response, nil
}

// GenerateStream 流式生成文本
func (a *CustomAdapter) GenerateStream(ctx context.Context, req ModelRequest, callback func(*ModelResponse)) error {
	if !a.SupportsStreaming() {
		return fmt.Errorf("streaming not supported by this custom adapter")
	}

	utils.Debug("自定义适配器流式生成文本: %s", a.config.Name)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.config.Timeout)*time.Second)
	defer cancel()

	// TODO: 实际实现流式调用
	// 这里模拟流式响应

	chunks := []string{
		"自定义",
		"模型",
		a.config.Name,
		"的",
		"流式",
		"响应",
		"(实验性功能)",
		"。",
		"\n\n",
		"API:",
		a.config.BaseURL,
		"\n",
		"输入:",
		req.Prompt[:min(20, len(req.Prompt))],
	}

	if len(req.Prompt) > 20 {
		chunks = append(chunks, "...")
	}

	for i, chunk := range chunks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(200 * time.Millisecond)

			response := &ModelResponse{
				Content: chunk,
				Metadata: map[string]interface{}{
					"model":     a.config.Name,
					"chunk":     i + 1,
					"total":     len(chunks),
					"custom":    true,
					"timestamp": time.Now().Unix(),
				},
			}

			callback(response)
		}
	}

	callback(&ModelResponse{
		FinishReason: "stop",
		Usage: TokenUsage{
			PromptTokens:     len(req.Prompt) / 4,
			CompletionTokens: 90,
			TotalTokens:      len(req.Prompt)/4 + 90,
		},
		Metadata: map[string]interface{}{
			"model":     a.config.Name,
			"completed": true,
			"custom":    true,
			"timestamp": time.Now().Unix(),
		},
	})

	return nil
}

// HealthCheck 健康检查
func (a *CustomAdapter) HealthCheck(ctx context.Context) error {
	utils.Debug("自定义适配器健康检查: %s", a.config.Name)

	// 检查自定义API是否可用
	// 尝试发送一个简单的请求到健康检查端点或根端点

	healthURL := a.config.BaseURL
	if !strings.HasSuffix(healthURL, "/") {
		healthURL += "/"
	}

	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	// 添加API密钥头（如果配置了）
	if a.config.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.config.APIKey))
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("custom API not available at %s: %w", a.config.BaseURL, err)
	}
	defer resp.Body.Close()

	// 接受2xx状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("custom API returned status %d", resp.StatusCode)
	}

	utils.Info("自定义API健康检查通过: %s", a.config.BaseURL)
	return nil
}
