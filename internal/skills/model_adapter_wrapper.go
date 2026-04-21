package skills

import (
	"context"
	"fmt"

	"GenPulse/internal/genkit/models"
)

// ModelAdapterWrapper 包装models.UnifiedModelAdapter以实现LLMClient接口
type ModelAdapterWrapper struct {
	adapter    *models.UnifiedModelAdapter
	modelID    string
	defaultReq models.ModelRequest
}

// NewModelAdapterWrapper 创建新的模型适配器包装器
func NewModelAdapterWrapper(adapter *models.UnifiedModelAdapter, modelID string) *ModelAdapterWrapper {
	return &ModelAdapterWrapper{
		adapter: adapter,
		modelID: modelID,
		defaultReq: models.ModelRequest{
			Temperature: 0.2,
			MaxTokens:   4096,
		},
	}
}

// Generate 实现LLMClient接口
func (w *ModelAdapterWrapper) Generate(ctx context.Context, prompt string, options map[string]any) (string, error) {
	// 构建模型请求
	req := w.defaultReq
	req.Prompt = prompt

	// 应用选项
	if temp, ok := options["temperature"].(float64); ok {
		req.Temperature = temp
	}
	if maxTokens, ok := options["max_tokens"].(int); ok {
		req.MaxTokens = maxTokens
	}
	if stream, ok := options["stream"].(bool); ok {
		req.Stream = stream
	}
	if metadata, ok := options["metadata"].(map[string]interface{}); ok {
		req.Metadata = metadata
	}

	// 调用模型适配器
	resp, err := w.adapter.Generate(ctx, w.modelID, req)
	if err != nil {
		return "", fmt.Errorf("model generation failed: %w", err)
	}

	return resp.Content, nil
}

// WithTemperature 设置温度
func (w *ModelAdapterWrapper) WithTemperature(temp float64) *ModelAdapterWrapper {
	w.defaultReq.Temperature = temp
	return w
}

// WithMaxTokens 设置最大token数
func (w *ModelAdapterWrapper) WithMaxTokens(maxTokens int) *ModelAdapterWrapper {
	w.defaultReq.MaxTokens = maxTokens
	return w
}

// WithModelID 设置模型ID
func (w *ModelAdapterWrapper) WithModelID(modelID string) *ModelAdapterWrapper {
	w.modelID = modelID
	return w
}
