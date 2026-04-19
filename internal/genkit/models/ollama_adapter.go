package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"GenPulse/internal/utils"
)

// OllamaAdapter Ollama模型适配器
type OllamaAdapter struct {
	config ModelConfig
	client *http.Client
}

// NewOllamaAdapter 创建Ollama适配器
func NewOllamaAdapter(config ModelConfig) (ModelAdapter, error) {
	// 设置默认值
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:11434"
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 4096
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.Timeout == 0 {
		config.Timeout = 60 // Ollama可能需要更长时间
	}

	adapter := &OllamaAdapter{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}

	utils.Info("创建Ollama适配器: %s (URL: %s)", config.Name, config.BaseURL)
	return adapter, nil
}

// Name 获取适配器名称
func (a *OllamaAdapter) Name() string {
	return a.config.Name
}

// Type 获取模型类型
func (a *OllamaAdapter) Type() ModelType {
	return ModelTypeOllama
}

// SupportsStreaming 是否支持流式
func (a *OllamaAdapter) SupportsStreaming() bool {
	return true
}

// SupportsTools 是否支持工具调用
func (a *OllamaAdapter) SupportsTools() bool {
	// Ollama对工具调用的支持取决于模型
	// 这里假设支持，但实际使用时需要检查具体模型
	return true
}

// GetConfig 获取配置
func (a *OllamaAdapter) GetConfig() ModelConfig {
	return a.config
}

// UpdateConfig 更新配置
func (a *OllamaAdapter) UpdateConfig(config ModelConfig) error {
	a.config = config
	// 更新HTTP客户端超时
	a.client.Timeout = time.Duration(config.Timeout) * time.Second
	return nil
}

// Generate 生成文本
func (a *OllamaAdapter) Generate(ctx context.Context, req ModelRequest) (*ModelResponse, error) {
	utils.Debug("Ollama适配器生成文本: %s", a.config.Name)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.config.Timeout)*time.Second)
	defer cancel()

	// 处理消息格式
	var prompt string
	if len(req.Messages) > 0 {
		// Ollama通常使用特定的消息格式
		for _, msg := range req.Messages {
			prompt += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
		}
	} else {
		prompt = req.Prompt
	}

	// TODO: 实际调用Ollama API
	// 这里先返回模拟响应

	response := &ModelResponse{
		Content:      fmt.Sprintf("这是来自本地Ollama模型 %s 的模拟响应。\n\n输入: %s", a.config.Name, prompt),
		FinishReason: "stop",
		Usage: TokenUsage{
			PromptTokens:     len(prompt) / 4,
			CompletionTokens: 200,
			TotalTokens:      len(prompt)/4 + 200,
		},
		Metadata: map[string]interface{}{
			"model":     a.config.Name,
			"provider":  "local",
			"adapter":   "ollama",
			"base_url":  a.config.BaseURL,
			"timestamp": time.Now().Unix(),
		},
	}

	return response, nil
}

// GenerateStream 流式生成文本
func (a *OllamaAdapter) GenerateStream(ctx context.Context, req ModelRequest, callback func(*ModelResponse)) error {
	utils.Debug("Ollama适配器流式生成文本: %s", a.config.Name)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.config.Timeout)*time.Second)
	defer cancel()

	// 处理消息格式
	var prompt string
	if len(req.Messages) > 0 {
		for _, msg := range req.Messages {
			prompt += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
		}
	} else {
		prompt = req.Prompt
	}

	// TODO: 实际实现流式调用
	// 这里模拟流式响应

	// 模拟分块响应
	chunks := []string{
		"这是",
		"来自",
		"本地",
		"Ollama",
		"模型",
		a.config.Name,
		"的",
		"流式",
		"响应。",
		"\n\n",
		"输入内容:",
		prompt[:min(30, len(prompt))],
	}
	if len(prompt) > 30 {
		chunks = append(chunks, "...")
	}

	for i, chunk := range chunks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 模拟延迟（本地模型可能更快或更慢）
			time.Sleep(150 * time.Millisecond)

			response := &ModelResponse{
				Content: chunk,
				Metadata: map[string]interface{}{
					"model":     a.config.Name,
					"chunk":     i + 1,
					"total":     len(chunks),
					"local":     true,
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
			PromptTokens:     len(prompt) / 4,
			CompletionTokens: 100,
			TotalTokens:      len(prompt)/4 + 100,
		},
		Metadata: map[string]interface{}{
			"model":     a.config.Name,
			"completed": true,
			"local":     true,
			"timestamp": time.Now().Unix(),
		},
	})

	return nil
}

// HealthCheck 健康检查
func (a *OllamaAdapter) HealthCheck(ctx context.Context) error {
	utils.Debug("Ollama适配器健康检查: %s", a.config.Name)

	// 检查Ollama服务是否可用
	healthURL := fmt.Sprintf("%s/api/tags", a.config.BaseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("Ollama service not available at %s: %w", a.config.BaseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama service returned status %d", resp.StatusCode)
	}

	utils.Info("Ollama服务健康检查通过: %s", a.config.BaseURL)
	return nil
}
