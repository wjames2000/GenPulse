package skills

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// SkillManager 技能管理器
type SkillManager struct {
	registry  *Registry
	generator *Generator
	validator *Validator
	loader    *ProgressiveLoader
	trigger   *TriggerManager
}

// NewSkillManager 创建技能管理器
func NewSkillManager(skillsDir string, llm LLMClient) (*SkillManager, error) {
	// 创建注册表
	registry, err := NewRegistry(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry: %w", err)
	}

	// 创建其他组件
	generator := NewGenerator(registry, llm)
	validator := NewValidator()
	loader := NewProgressiveLoader(registry, DefaultLoaderOptions())
	trigger := NewTriggerManager(registry)

	return &SkillManager{
		registry:  registry,
		generator: generator,
		validator: validator,
		loader:    loader,
		trigger:   trigger,
	}, nil
}

// ProcessTaskExecution 处理任务执行记录
func (sm *SkillManager) ProcessTaskExecution(record *TaskExecutionRecord) (*SkillProcessingResult, error) {
	startTime := time.Now()

	// 检查是否触发技能生成
	shouldExtract, triggerResults := sm.trigger.ShouldExtractSkill(record)

	result := &SkillProcessingResult{
		TaskID:         record.TaskID,
		Triggered:      shouldExtract,
		TriggerResults: triggerResults,
		StartTime:      startTime,
	}

	if !shouldExtract {
		result.EndTime = time.Now()
		result.Message = "No skill extraction triggered"
		return result, nil
	}

	// 生成技能
	options := DefaultGenerationOptions()
	options.Category = inferCategory(record)
	options.Tags = inferTags(record)
	options.Complexity = record.Complexity

	// 使用最高置信度的触发器
	var bestTrigger TriggerResult
	if len(triggerResults) > 0 {
		bestTrigger = triggerResults[0]
		for _, trigger := range triggerResults {
			if trigger.Confidence > bestTrigger.Confidence {
				bestTrigger = trigger
			}
		}
	}

	genResult, err := sm.generator.GenerateFromExperience(
		context.Background(),
		record,
		bestTrigger,
		options,
	)

	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("Skill generation failed: %v", err)
		result.EndTime = time.Now()
		return result, err
	}

	// 验证技能
	validationReport := sm.validator.Validate(genResult.Skill)

	// 如果验证通过，注册技能
	if validationReport.OverallPass {
		if err := sm.registry.Register(genResult.Skill); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("Skill registration failed: %v", err)
			result.EndTime = time.Now()
			return result, err
		}

		result.Success = true
		result.GeneratedSkill = genResult.Skill
		result.ValidationReport = validationReport
		result.Message = fmt.Sprintf("Skill '%s' generated and registered successfully", genResult.Skill.Name)
	} else {
		result.Success = false
		result.GeneratedSkill = genResult.Skill
		result.ValidationReport = validationReport
		result.Message = "Skill generated but failed validation"
		result.Error = "Validation failed"
	}

	result.EndTime = time.Now()
	return result, nil
}

// GetSkill 获取技能
func (sm *SkillManager) GetSkill(skillID string, level LoadLevel) (*LoadResult, error) {
	options := LoadOptions{
		Level:    level,
		Strategy: StrategyLazy,
	}

	return sm.loader.Load(skillID, options)
}

// ListSkills 列出技能
func (sm *SkillManager) ListSkills(filter map[string]any) ([]*SkillMetadata, error) {
	// 需要将map[string]any转换为map[string]string
	stringFilter := make(map[string]string)
	for k, v := range filter {
		if s, ok := v.(string); ok {
			stringFilter[k] = s
		}
	}
	return sm.registry.Search("", stringFilter)
}

// EnableSkill 启用技能
func (sm *SkillManager) EnableSkill(skillID string) error {
	return sm.registry.Enable(skillID)
}

// DisableSkill 禁用技能
func (sm *SkillManager) DisableSkill(skillID string) error {
	return sm.registry.Disable(skillID)
}

// ValidateSkill 验证技能
func (sm *SkillManager) ValidateSkill(skillID string) (*ValidationReport, error) {
	skill, err := sm.registry.Get(skillID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}

	report := sm.validator.Validate(skill)

	// 如果验证通过，更新验证状态
	if report.OverallPass {
		if err := sm.registry.Validate(skillID); err != nil {
			return report, fmt.Errorf("validation passed but failed to update status: %w", err)
		}
	}

	return report, nil
}

// UpdateSkill 更新技能
func (sm *SkillManager) UpdateSkill(skill *Skill) error {
	// 重新验证
	report := sm.validator.Validate(skill)
	if !report.OverallPass {
		return fmt.Errorf("skill validation failed: %v", report.CriticalFailures)
	}

	return sm.registry.Update(skill)
}

// DeleteSkill 删除技能
func (sm *SkillManager) DeleteSkill(skillID string) error {
	return sm.registry.Delete(skillID)
}

// SearchSkills 搜索技能
func (sm *SkillManager) SearchSkills(query string, filters map[string]any) ([]*SkillMetadata, error) {
	// 需要将map[string]any转换为map[string]string
	stringFilters := make(map[string]string)
	for k, v := range filters {
		if s, ok := v.(string); ok {
			stringFilters[k] = s
		}
	}
	return sm.registry.Search(query, stringFilters)
}

// GetSkillStats 获取技能统计
func (sm *SkillManager) GetSkillStats() (*SkillStats, error) {
	metadatas, err := sm.registry.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list skills: %w", err)
	}

	stats := &SkillStats{
		TotalSkills:        len(metadatas),
		EnabledSkills:      0,
		ValidatedSkills:    0,
		ByCategory:         make(map[string]int),
		ByComplexity:       make(map[string]int),
		TotalUsage:         0,
		AverageSuccessRate: 0.0,
	}

	var totalSuccessRate float64
	var skillsWithUsage int

	for _, metadata := range metadatas {
		if metadata.Enabled {
			stats.EnabledSkills++
		}
		if metadata.Validated {
			stats.ValidatedSkills++
		}

		stats.ByCategory[metadata.Category]++
		stats.ByComplexity[metadata.Complexity]++

		stats.TotalUsage += metadata.UsageCount

		if metadata.UsageCount > 0 {
			totalSuccessRate += metadata.SuccessRate
			skillsWithUsage++
		}
	}

	if skillsWithUsage > 0 {
		stats.AverageSuccessRate = totalSuccessRate / float64(skillsWithUsage)
	}

	// 获取加载器统计
	loadStats := sm.loader.GetStats()
	stats.LoadStats = loadStats

	return stats, nil
}

// RecordSkillUsage 记录技能使用
func (sm *SkillManager) RecordSkillUsage(skillID string, success bool) error {
	return sm.registry.IncrementUsage(skillID, success)
}

// GetRelatedSkills 获取相关技能
func (sm *SkillManager) GetRelatedSkills(skillID string) ([]*SkillMetadata, error) {
	skill, err := sm.registry.Get(skillID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}

	var related []*SkillMetadata

	// 获取前置技能
	for _, prereqID := range skill.Prerequisites {
		metadata, err := sm.registry.GetMetadata(prereqID)
		if err == nil {
			related = append(related, metadata)
		}
	}

	// 获取相同分类的技能
	sameCategory, err := sm.registry.GetByCategory(skill.Category)
	if err == nil {
		for _, metadata := range sameCategory {
			if metadata.ID != skillID {
				related = append(related, metadata)
			}
		}
	}

	// 去重
	uniqueRelated := make(map[string]*SkillMetadata)
	for _, metadata := range related {
		uniqueRelated[metadata.ID] = metadata
	}

	// 转换为切片
	var result []*SkillMetadata
	for _, metadata := range uniqueRelated {
		result = append(result, metadata)
	}

	return result, nil
}

// ExportSkill 导出技能
func (sm *SkillManager) ExportSkill(skillID string, format string) ([]byte, error) {
	skill, err := sm.registry.Get(skillID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}

	switch format {
	case "yaml":
		return yaml.Marshal(skill)
	case "json":
		return json.MarshalIndent(skill, "", "  ")
	case "markdown":
		md := generateMarkdown(skill)
		return []byte(md), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// ImportSkill 导入技能
func (sm *SkillManager) ImportSkill(data []byte, format string) (*Skill, error) {
	var skill Skill

	switch format {
	case "yaml":
		if err := yaml.Unmarshal(data, &skill); err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
		}
	case "json":
		if err := json.Unmarshal(data, &skill); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	// 验证技能
	report := sm.validator.Validate(&skill)
	if !report.OverallPass {
		return nil, fmt.Errorf("imported skill validation failed: %v", report.CriticalFailures)
	}

	// 注册技能
	if err := sm.registry.Register(&skill); err != nil {
		return nil, fmt.Errorf("failed to register imported skill: %w", err)
	}

	return &skill, nil
}

// ManualSkillGeneration 手动技能生成
func (sm *SkillManager) ManualSkillGeneration(record *TaskExecutionRecord, reason string) (*SkillProcessingResult, error) {
	// 创建手动触发器
	trigger := sm.trigger.ManualTrigger(record, reason)

	// 生成技能
	options := DefaultGenerationOptions()
	options.Category = inferCategory(record)
	options.Tags = inferTags(record)
	options.Complexity = record.Complexity

	genResult, err := sm.generator.GenerateFromExperience(
		context.Background(),
		record,
		trigger,
		options,
	)

	if err != nil {
		return nil, fmt.Errorf("manual skill generation failed: %w", err)
	}

	// 验证并注册
	validationReport := sm.validator.Validate(genResult.Skill)

	result := &SkillProcessingResult{
		TaskID:           record.TaskID,
		Triggered:        true,
		TriggerResults:   []TriggerResult{trigger},
		StartTime:        time.Now(),
		GeneratedSkill:   genResult.Skill,
		ValidationReport: validationReport,
	}

	if validationReport.OverallPass {
		if err := sm.registry.Register(genResult.Skill); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("Skill registration failed: %v", err)
		} else {
			result.Success = true
			result.Message = fmt.Sprintf("Skill '%s' generated and registered successfully", genResult.Skill.Name)
		}
	} else {
		result.Success = false
		result.Message = "Skill generated but failed validation"
		result.Error = "Validation failed"
	}

	result.EndTime = time.Now()
	return result, nil
}

// SkillProcessingResult 技能处理结果
type SkillProcessingResult struct {
	TaskID           string            `json:"task_id" yaml:"task_id"`
	Success          bool              `json:"success" yaml:"success"`
	Triggered        bool              `json:"triggered" yaml:"triggered"`
	TriggerResults   []TriggerResult   `json:"trigger_results,omitempty" yaml:"trigger_results,omitempty"`
	GeneratedSkill   *Skill            `json:"generated_skill,omitempty" yaml:"generated_skill,omitempty"`
	ValidationReport *ValidationReport `json:"validation_report,omitempty" yaml:"validation_report,omitempty"`
	Message          string            `json:"message" yaml:"message"`
	Error            string            `json:"error,omitempty" yaml:"error,omitempty"`
	StartTime        time.Time         `json:"start_time" yaml:"start_time"`
	EndTime          time.Time         `json:"end_time" yaml:"end_time"`
}

// SkillStats 技能统计
type SkillStats struct {
	TotalSkills        int            `json:"total_skills" yaml:"total_skills"`
	EnabledSkills      int            `json:"enabled_skills" yaml:"enabled_skills"`
	ValidatedSkills    int            `json:"validated_skills" yaml:"validated_skills"`
	ByCategory         map[string]int `json:"by_category" yaml:"by_category"`
	ByComplexity       map[string]int `json:"by_complexity" yaml:"by_complexity"`
	TotalUsage         int            `json:"total_usage" yaml:"total_usage"`
	AverageSuccessRate float64        `json:"average_success_rate" yaml:"average_success_rate"`
	LoadStats          *LoadStats     `json:"load_stats,omitempty" yaml:"load_stats,omitempty"`
}

// inferCategory 推断分类
func inferCategory(record *TaskExecutionRecord) string {
	// 基于任务类型推断分类
	switch record.TaskType {
	case "file_operation", "fs_operation":
		return "filesystem"
	case "git_operation":
		return "version_control"
	case "shell_command":
		return "system"
	case "project_setup":
		return "project"
	case "code_generation":
		return "development"
	case "testing":
		return "testing"
	case "deployment":
		return "deployment"
	default:
		return "general"
	}
}

// inferTags 推断标签
func inferTags(record *TaskExecutionRecord) []string {
	var tags []string

	// 添加任务类型标签
	if record.TaskType != "" {
		tags = append(tags, record.TaskType)
	}

	// 添加复杂度标签
	if record.Complexity != "" {
		tags = append(tags, record.Complexity)
	}

	// 添加工具标签
	for tool := range record.ToolUsage {
		tags = append(tags, tool)
	}

	// 添加Agent标签
	tags = append(tags, record.AgentInvolved...)

	// 添加成功/失败标签
	if record.Success {
		tags = append(tags, "success")
	} else {
		tags = append(tags, "failure")
	}

	return tags
}
