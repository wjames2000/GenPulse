package models

import (
	"context"
	"fmt"
	"strings"

	"GenPulse/internal/utils"
)

// ModelType 模型类型
type ModelType string

const (
	ModelTypeGemini ModelType = "gemini"
	ModelTypeGPT    ModelType = "gpt"
	ModelTypeClaude ModelType = "claude"
	ModelTypeOllama ModelType = "ollama"
	ModelTypeCustom ModelType = "custom"
)

// ModelConfig 模型配置
type ModelConfig struct {
	Type        ModelType `json:"type"`
	Name        string    `json:"name"`
	Provider    string    `json:"provider"`
	APIKey      string    `json:"api_key,omitempty"`
	BaseURL     string    `json:"base_url,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Timeout     int       `json:"timeout,omitempty"` // 秒
}

// ModelRequest 模型请求
type ModelRequest struct {
	Prompt      string                 `json:"prompt"`
	Messages    []ChatMessage          `json:"messages,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	Tools       []ToolDefinition       `json:"tools,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ModelResponse 模型响应
type ModelResponse struct {
	Content      string                 `json:"content"`
	FinishReason string                 `json:"finish_reason,omitempty"`
	Usage        TokenUsage             `json:"usage,omitempty"`
	ToolCalls    []ToolCall             `json:"tool_calls,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Error        error                  `json:"error,omitempty"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"` // system, user, assistant, tool
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// ToolDefinition 工具定义
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// ToolCall 工具调用
type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"` // function
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction 工具调用函数
type ToolCallFunction struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// TokenUsage Token使用情况
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ModelAdapter 模型适配器接口
type ModelAdapter interface {
	// 基本信息
	Name() string
	Type() ModelType
	SupportsStreaming() bool
	SupportsTools() bool

	// 核心功能
	Generate(ctx context.Context, req ModelRequest) (*ModelResponse, error)
	GenerateStream(ctx context.Context, req ModelRequest, callback func(*ModelResponse)) error

	// 配置
	GetConfig() ModelConfig
	UpdateConfig(config ModelConfig) error

	// 健康检查
	HealthCheck(ctx context.Context) error
}

// ModelAdapterFactory 模型适配器工厂
type ModelAdapterFactory interface {
	CreateAdapter(config ModelConfig) (ModelAdapter, error)
	SupportedTypes() []ModelType
}

// UnifiedModelAdapter 统一模型适配器
type UnifiedModelAdapter struct {
	adapters     map[string]ModelAdapter
	defaultModel string
	factory      ModelAdapterFactory
}

// NewUnifiedModelAdapter 创建统一模型适配器
func NewUnifiedModelAdapter(factory ModelAdapterFactory) *UnifiedModelAdapter {
	return &UnifiedModelAdapter{
		adapters: make(map[string]ModelAdapter),
		factory:  factory,
	}
}

// RegisterModel 注册模型
func (uma *UnifiedModelAdapter) RegisterModel(config ModelConfig) error {
	// 生成模型ID
	modelID := generateModelID(config)

	// 检查是否已注册
	if _, exists := uma.adapters[modelID]; exists {
		return fmt.Errorf("model already registered: %s", modelID)
	}

	// 创建适配器
	adapter, err := uma.factory.CreateAdapter(config)
	if err != nil {
		return fmt.Errorf("failed to create adapter for %s: %w", modelID, err)
	}

	// 注册适配器
	uma.adapters[modelID] = adapter

	// 如果是第一个模型，设置为默认模型
	if uma.defaultModel == "" {
		uma.defaultModel = modelID
	}

	utils.Info("注册模型: %s (%s)", modelID, config.Type)
	return nil
}

// Generate 生成文本
func (uma *UnifiedModelAdapter) Generate(ctx context.Context, modelID string, req ModelRequest) (*ModelResponse, error) {
	// 如果未指定模型ID，使用默认模型
	if modelID == "" {
		modelID = uma.defaultModel
		if modelID == "" {
			return nil, fmt.Errorf("no model specified and no default model available")
		}
	}

	// 获取适配器
	adapter, exists := uma.adapters[modelID]
	if !exists {
		return nil, fmt.Errorf("model not found: %s", modelID)
	}

	// 执行生成
	utils.Debug("使用模型 %s 生成文本", modelID)
	return adapter.Generate(ctx, req)
}

// GenerateStream 流式生成文本
func (uma *UnifiedModelAdapter) GenerateStream(ctx context.Context, modelID string, req ModelRequest, callback func(*ModelResponse)) error {
	// 如果未指定模型ID，使用默认模型
	if modelID == "" {
		modelID = uma.defaultModel
		if modelID == "" {
			return fmt.Errorf("no model specified and no default model available")
		}
	}

	// 获取适配器
	adapter, exists := uma.adapters[modelID]
	if !exists {
		return fmt.Errorf("model not found: %s", modelID)
	}

	// 检查是否支持流式
	if !adapter.SupportsStreaming() {
		return fmt.Errorf("model does not support streaming: %s", modelID)
	}

	// 执行流式生成
	utils.Debug("使用模型 %s 流式生成文本", modelID)
	return adapter.GenerateStream(ctx, req, callback)
}

// SetDefaultModel 设置默认模型
func (uma *UnifiedModelAdapter) SetDefaultModel(modelID string) error {
	if _, exists := uma.adapters[modelID]; !exists {
		return fmt.Errorf("model not found: %s", modelID)
	}

	uma.defaultModel = modelID
	utils.Info("设置默认模型: %s", modelID)
	return nil
}

// GetDefaultModel 获取默认模型
func (uma *UnifiedModelAdapter) GetDefaultModel() string {
	return uma.defaultModel
}

// ListModels 列出所有已注册的模型
func (uma *UnifiedModelAdapter) ListModels() []string {
	var models []string
	for modelID := range uma.adapters {
		models = append(models, modelID)
	}
	return models
}

// GetModelInfo 获取模型信息
func (uma *UnifiedModelAdapter) GetModelInfo(modelID string) (ModelConfig, error) {
	adapter, exists := uma.adapters[modelID]
	if !exists {
		return ModelConfig{}, fmt.Errorf("model not found: %s", modelID)
	}

	return adapter.GetConfig(), nil
}

// HealthCheck 健康检查
func (uma *UnifiedModelAdapter) HealthCheck(ctx context.Context) map[string]error {
	errors := make(map[string]error)

	for modelID, adapter := range uma.adapters {
		if err := adapter.HealthCheck(ctx); err != nil {
			errors[modelID] = err
			utils.Warn("模型健康检查失败 %s: %v", modelID, err)
		}
	}

	return errors
}

// generateModelID 生成模型ID
func generateModelID(config ModelConfig) string {
	// 格式: provider-type-name
	provider := strings.ToLower(config.Provider)
	modelType := strings.ToLower(string(config.Type))
	name := strings.ToLower(config.Name)

	return fmt.Sprintf("%s-%s-%s", provider, modelType, name)
}

// DefaultModelAdapterFactory 默认模型适配器工厂
type DefaultModelAdapterFactory struct{}

// CreateAdapter 创建适配器
func (f *DefaultModelAdapterFactory) CreateAdapter(config ModelConfig) (ModelAdapter, error) {
	switch config.Type {
	case ModelTypeGemini:
		return NewGeminiAdapter(config)
	case ModelTypeGPT:
		return NewGPTAdapter(config)
	case ModelTypeClaude:
		return NewClaudeAdapter(config)
	case ModelTypeOllama:
		return NewOllamaAdapter(config)
	case ModelTypeCustom:
		return NewCustomAdapter(config)
	default:
		return nil, fmt.Errorf("unsupported model type: %s", config.Type)
	}
}

// SupportedTypes 支持的模型类型
func (f *DefaultModelAdapterFactory) SupportedTypes() []ModelType {
	return []ModelType{
		ModelTypeGemini,
		ModelTypeGPT,
		ModelTypeClaude,
		ModelTypeOllama,
		ModelTypeCustom,
	}
}
