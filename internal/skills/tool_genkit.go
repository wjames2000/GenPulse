package skills

import (
	"context"
	"fmt"
	"path/filepath"
)

// SkillTool 技能管理工具的Genkit封装
type SkillTool struct {
	manager *SkillManager
}

// NewSkillTool 创建技能管理工具
func NewSkillTool(skillsDir string, llmClient LLMClient) (*SkillTool, error) {
	manager, err := NewSkillManager(skillsDir, llmClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create skill manager: %w", err)
	}

	return &SkillTool{
		manager: manager,
	}, nil
}

// RegisterTools 注册所有技能管理工具
func (st *SkillTool) RegisterTools() error {
	// TODO: 实现Genkit工具注册
	// 注册技能处理工具
	// if err := flow.DefineTool("skill_process_task", "Process task execution and generate skills if triggered", st.ProcessTaskTool); err != nil {
	// 	return fmt.Errorf("failed to register process_task tool: %w", err)
	// }

	// 其他工具注册...

	return nil
}

// ProcessTaskTool 处理任务执行的工具
func (st *SkillTool) ProcessTaskTool(ctx context.Context, input *ProcessTaskInput) (*SkillProcessingResult, error) {
	return st.manager.ProcessTaskExecution(input.Record)
}

// GetSkillTool 获取技能的工具
func (st *SkillTool) GetSkillTool(ctx context.Context, input *GetSkillInput) (*LoadResult, error) {
	return st.manager.GetSkill(input.SkillID, input.Level)
}

// ListSkillsTool 列出技能的工具
func (st *SkillTool) ListSkillsTool(ctx context.Context, input *ListSkillsInput) ([]*SkillMetadata, error) {
	return st.manager.ListSkills(input.Filter)
}

// SearchSkillsTool 搜索技能的工具
func (st *SkillTool) SearchSkillsTool(ctx context.Context, input *SearchSkillsInput) ([]*SkillMetadata, error) {
	return st.manager.SearchSkills(input.Query, input.Filters)
}

// EnableSkillTool 启用技能的工具
func (st *SkillTool) EnableSkillTool(ctx context.Context, input *EnableSkillInput) error {
	return st.manager.EnableSkill(input.SkillID)
}

// DisableSkillTool 禁用技能的工具
func (st *SkillTool) DisableSkillTool(ctx context.Context, input *DisableSkillInput) error {
	return st.manager.DisableSkill(input.SkillID)
}

// ValidateSkillTool 验证技能的工具
func (st *SkillTool) ValidateSkillTool(ctx context.Context, input *ValidateSkillInput) (*ValidationReport, error) {
	return st.manager.ValidateSkill(input.SkillID)
}

// UpdateSkillTool 更新技能的工具
func (st *SkillTool) UpdateSkillTool(ctx context.Context, input *UpdateSkillInput) error {
	return st.manager.UpdateSkill(input.Skill)
}

// DeleteSkillTool 删除技能的工具
func (st *SkillTool) DeleteSkillTool(ctx context.Context, input *DeleteSkillInput) error {
	return st.manager.DeleteSkill(input.SkillID)
}

// GetSkillStatsTool 获取技能统计的工具
func (st *SkillTool) GetSkillStatsTool(ctx context.Context, input *GetSkillStatsInput) (*SkillStats, error) {
	return st.manager.GetSkillStats()
}

// RecordSkillUsageTool 记录技能使用的工具
func (st *SkillTool) RecordSkillUsageTool(ctx context.Context, input *RecordSkillUsageInput) error {
	return st.manager.RecordSkillUsage(input.SkillID, input.Success)
}

// GetRelatedSkillsTool 获取相关技能的工具
func (st *SkillTool) GetRelatedSkillsTool(ctx context.Context, input *GetRelatedSkillsInput) ([]*SkillMetadata, error) {
	return st.manager.GetRelatedSkills(input.SkillID)
}

// ExportSkillTool 导出技能的工具
func (st *SkillTool) ExportSkillTool(ctx context.Context, input *ExportSkillInput) (*ExportSkillOutput, error) {
	data, err := st.manager.ExportSkill(input.SkillID, input.Format)
	if err != nil {
		return nil, err
	}

	return &ExportSkillOutput{
		Data:   data,
		Format: input.Format,
	}, nil
}

// ImportSkillTool 导入技能的工具
func (st *SkillTool) ImportSkillTool(ctx context.Context, input *ImportSkillInput) (*Skill, error) {
	return st.manager.ImportSkill(input.Data, input.Format)
}

// ManualSkillGenerationTool 手动技能生成工具
func (st *SkillTool) ManualSkillGenerationTool(ctx context.Context, input *ManualSkillGenerationInput) (*SkillProcessingResult, error) {
	return st.manager.ManualSkillGeneration(input.Record, input.Reason)
}

// 工具输入输出结构定义

type ProcessTaskInput struct {
	Record *TaskExecutionRecord `json:"record"`
}

type GetSkillInput struct {
	SkillID string    `json:"skill_id"`
	Level   LoadLevel `json:"level,omitempty"`
}

type ListSkillsInput struct {
	Filter map[string]any `json:"filter,omitempty"`
}

type SearchSkillsInput struct {
	Query   string         `json:"query"`
	Filters map[string]any `json:"filters,omitempty"`
}

type EnableSkillInput struct {
	SkillID string `json:"skill_id"`
}

type DisableSkillInput struct {
	SkillID string `json:"skill_id"`
}

type ValidateSkillInput struct {
	SkillID string `json:"skill_id"`
}

type UpdateSkillInput struct {
	Skill *Skill `json:"skill"`
}

type DeleteSkillInput struct {
	SkillID string `json:"skill_id"`
}

type GetSkillStatsInput struct {
	// 空结构，不需要输入参数
}

type RecordSkillUsageInput struct {
	SkillID string `json:"skill_id"`
	Success bool   `json:"success"`
}

type GetRelatedSkillsInput struct {
	SkillID string `json:"skill_id"`
}

type ExportSkillInput struct {
	SkillID string `json:"skill_id"`
	Format  string `json:"format"` // yaml, json, markdown
}

type ExportSkillOutput struct {
	Data   []byte `json:"data"`
	Format string `json:"format"`
}

type ImportSkillInput struct {
	Data   []byte `json:"data"`
	Format string `json:"format"` // yaml, json
}

type ManualSkillGenerationInput struct {
	Record *TaskExecutionRecord `json:"record"`
	Reason string               `json:"reason"`
}

// LLMClientWrapper Genkit LLM客户端包装器
type LLMClientWrapper struct {
	// model *ai.Model
}

// NewLLMClientWrapper 创建LLM客户端包装器
func NewLLMClientWrapper(model interface{}) *LLMClientWrapper {
	return &LLMClientWrapper{
		// model: model,
	}
}

// Generate 实现LLMClient接口
func (w *LLMClientWrapper) Generate(ctx context.Context, prompt string, options map[string]any) (string, error) {
	// TODO: 实现Genkit LLM调用
	// 构建生成请求
	// req := &ai.GenerateRequest{
	// 	Messages: []*ai.Message{
	// 		{
	// 			Role:    ai.RoleUser,
	// 			Content: []*ai.Part{ai.NewTextPart(prompt)},
	// 		},
	// 	},
	// }

	// 调用Genkit生成
	// resp, err := w.model.Generate(ctx, req)
	// if err != nil {
	// 	return "", fmt.Errorf("LLM generation failed: %w", err)
	// }

	// 模拟响应
	return "模拟LLM响应: " + prompt[:minInt(50, len(prompt))] + "...", nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// InitializeSkillTools 初始化技能管理工具
func InitializeSkillTools(genkitConfig *GenkitConfig) (*SkillTool, error) {
	// 确定技能目录
	skillsDir := genkitConfig.SkillsDir
	if skillsDir == "" {
		// 使用默认目录
		skillsDir = filepath.Join(".", "skills")
	}

	// 创建LLM客户端包装器
	llmClient := NewLLMClientWrapper(genkitConfig.Model)

	// 创建技能管理工具
	skillTool, err := NewSkillTool(skillsDir, llmClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create skill tool: %w", err)
	}

	// 注册工具
	if err := skillTool.RegisterTools(); err != nil {
		return nil, fmt.Errorf("failed to register skill tools: %w", err)
	}

	return skillTool, nil
}

// GenkitConfig Genkit配置
type GenkitConfig struct {
	SkillsDir string      `json:"skills_dir" yaml:"skills_dir"`
	Model     interface{} `json:"model" yaml:"model"`
}

// DefaultGenkitConfig 默认Genkit配置
func DefaultGenkitConfig() *GenkitConfig {
	return &GenkitConfig{
		SkillsDir: "./skills",
		Model:     nil, // 需要在外部设置
	}
}
