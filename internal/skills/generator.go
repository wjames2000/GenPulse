package skills

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Generator 技能生成器
type Generator struct {
	registry *Registry
	llm      LLMClient
	prompts  PromptTemplates
}

// LLMClient LLM客户端接口
type LLMClient interface {
	Generate(ctx context.Context, prompt string, options map[string]any) (string, error)
}

// PromptTemplates 提示词模板
type PromptTemplates struct {
	SkillExtraction string
	StepRefinement  string
	ExampleCreation string
	Validation      string
}

// GenerationOptions 生成选项
type GenerationOptions struct {
	Model             string         `json:"model" yaml:"model"`
	Temperature       float64        `json:"temperature" yaml:"temperature"`
	MaxTokens         int            `json:"max_tokens" yaml:"max_tokens"`
	IncludeExamples   bool           `json:"include_examples" yaml:"include_examples"`
	IncludeTips       bool           `json:"include_tips" yaml:"include_tips"`
	IncludeWarnings   bool           `json:"include_warnings" yaml:"include_warnings"`
	Category          string         `json:"category" yaml:"category"`
	Tags              []string       `json:"tags" yaml:"tags"`
	Complexity        string         `json:"complexity" yaml:"complexity"`
	AdditionalContext map[string]any `json:"additional_context" yaml:"additional_context"`
}

// GenerationResult 生成结果
type GenerationResult struct {
	Skill       *Skill        `json:"skill" yaml:"skill"`
	Success     bool          `json:"success" yaml:"success"`
	Error       string        `json:"error,omitempty" yaml:"error,omitempty"`
	LLMResponse string        `json:"llm_response,omitempty" yaml:"llm_response,omitempty"`
	Duration    time.Duration `json:"duration" yaml:"duration"`
	Trigger     TriggerResult `json:"trigger,omitempty" yaml:"trigger,omitempty"`
}

// NewGenerator 创建技能生成器
func NewGenerator(registry *Registry, llm LLMClient) *Generator {
	return &Generator{
		registry: registry,
		llm:      llm,
		prompts:  getDefaultPrompts(),
	}
}

// GenerateFromExperience 从经验记录生成技能
func (g *Generator) GenerateFromExperience(ctx context.Context, record *TaskExecutionRecord, trigger TriggerResult, options GenerationOptions) (*GenerationResult, error) {
	startTime := time.Now()

	// 构建生成上下文
	context := g.buildGenerationContext(record, trigger, options)

	// 生成技能名称和描述
	skillName, description, err := g.generateSkillInfo(ctx, record, context, options)
	if err != nil {
		return &GenerationResult{
			Success:  false,
			Error:    fmt.Sprintf("Failed to generate skill info: %v", err),
			Duration: time.Since(startTime),
			Trigger:  trigger,
		}, err
	}

	// 创建技能基础结构
	skill := NewSkill(skillName, description, "system")
	skill.SourceTaskID = record.TaskID
	skill.Category = options.Category
	skill.Tags = options.Tags
	skill.Complexity = options.Complexity
	skill.AgentTypes = record.AgentInvolved

	// 提取相关工具
	skill.RelatedTools = g.extractRelatedTools(record)

	// 生成步骤
	steps, err := g.generateSteps(ctx, record, context, options)
	if err != nil {
		return &GenerationResult{
			Success:  false,
			Error:    fmt.Sprintf("Failed to generate steps: %v", err),
			Duration: time.Since(startTime),
			Trigger:  trigger,
		}, err
	}

	skill.Steps = steps

	// 生成示例
	if options.IncludeExamples {
		examples, err := g.generateExamples(ctx, skill, context, options)
		if err != nil {
			// 示例生成失败不影响主要功能
			fmt.Printf("Warning: failed to generate examples: %v\n", err)
		} else {
			skill.Examples = examples
		}
	}

	// 生成技巧
	if options.IncludeTips {
		tips, err := g.generateTips(ctx, skill, context, options)
		if err != nil {
			fmt.Printf("Warning: failed to generate tips: %v\n", err)
		} else {
			skill.Tips = tips
		}
	}

	// 生成警告
	if options.IncludeWarnings {
		warnings, err := g.generateWarnings(ctx, skill, context, options)
		if err != nil {
			fmt.Printf("Warning: failed to generate warnings: %v\n", err)
		} else {
			skill.Warnings = warnings
		}
	}

	// 估算Token和平均执行时间
	skill.TokenEstimate = g.estimateTokens(skill)
	skill.AvgExecutionTime = g.estimateExecutionTime(record)

	// 设置验证状态（需要后续验证）
	skill.Validated = false

	result := &GenerationResult{
		Skill:    skill,
		Success:  true,
		Duration: time.Since(startTime),
		Trigger:  trigger,
	}

	return result, nil
}

// buildGenerationContext 构建生成上下文
func (g *Generator) buildGenerationContext(record *TaskExecutionRecord, trigger TriggerResult, options GenerationOptions) map[string]any {
	context := map[string]any{
		"task_id":            record.TaskID,
		"task_type":          record.TaskType,
		"task_description":   record.Description,
		"task_complexity":    record.Complexity,
		"success":            record.Success,
		"execution_time":     record.EndTime.Sub(record.StartTime).String(),
		"agent_involved":     record.AgentInvolved,
		"tool_usage":         record.ToolUsage,
		"steps_count":        len(record.Steps),
		"trigger_type":       trigger.Type,
		"trigger_reason":     trigger.Reason,
		"trigger_confidence": trigger.Confidence,
		"generation_options": options,
	}

	// 添加步骤摘要
	var stepSummaries []map[string]any
	for _, step := range record.Steps {
		stepSummary := map[string]any{
			"step_id":  step.StepID,
			"order":    step.Order,
			"action":   step.Action,
			"tool":     step.Tool,
			"success":  step.Success,
			"duration": step.EndTime.Sub(step.StartTime).String(),
		}
		stepSummaries = append(stepSummaries, stepSummary)
	}
	context["step_summaries"] = stepSummaries

	// 添加上下文信息
	if record.Context != nil {
		context["task_context"] = record.Context
	}

	// 添加输出信息
	if record.Output != nil {
		context["task_output"] = record.Output
	}

	// 添加额外上下文
	if options.AdditionalContext != nil {
		for key, value := range options.AdditionalContext {
			context[key] = value
		}
	}

	return context
}

// generateSkillInfo 生成技能名称和描述
func (g *Generator) generateSkillInfo(ctx context.Context, record *TaskExecutionRecord, context map[string]any, options GenerationOptions) (string, string, error) {
	prompt := g.prompts.SkillExtraction

	// 替换模板变量
	prompt = replaceTemplateVariables(prompt, map[string]string{
		"{{task_description}}": record.Description,
		"{{task_type}}":        record.TaskType,
		"{{task_complexity}}":  record.Complexity,
		"{{success}}":          fmt.Sprintf("%v", record.Success),
		"{{agent_involved}}":   strings.Join(record.AgentInvolved, ", "),
		"{{trigger_reason}}":   context["trigger_reason"].(string),
	})

	// 调用LLM生成
	llmOptions := map[string]any{
		"temperature": options.Temperature,
		"max_tokens":  options.MaxTokens,
		"model":       options.Model,
	}

	response, err := g.llm.Generate(ctx, prompt, llmOptions)
	if err != nil {
		return "", "", fmt.Errorf("LLM generation failed: %w", err)
	}

	// 解析响应（假设格式：名称|描述）
	parts := strings.SplitN(response, "|", 2)
	if len(parts) != 2 {
		// 尝试其他分隔符
		parts = strings.SplitN(response, "\n", 2)
		if len(parts) != 2 {
			// 使用默认值
			name := fmt.Sprintf("Skill for %s", record.TaskType)
			description := fmt.Sprintf("Automatically generated skill for %s task", record.Description)
			return name, description, nil
		}
	}

	name := strings.TrimSpace(parts[0])
	description := strings.TrimSpace(parts[1])

	// 清理名称（移除特殊字符）
	name = cleanSkillName(name)

	return name, description, nil
}

// generateSteps 生成步骤
func (g *Generator) generateSteps(ctx context.Context, record *TaskExecutionRecord, context map[string]any, options GenerationOptions) ([]Step, error) {
	var steps []Step

	// 从任务记录中提取步骤
	for i, taskStep := range record.Steps {
		// 生成步骤描述
		stepDescription := g.generateStepDescription(taskStep, i+1, len(record.Steps))

		// 提取参数
		parameters := g.extractParameters(taskStep)

		// 提取条件
		conditions := g.extractConditions(taskStep, record)

		// 创建步骤
		step := Step{
			ID:         generateStepID("temp", i+1),
			Order:      i + 1,
			Action:     stepDescription,
			Tool:       taskStep.Tool,
			Parameters: parameters,
			Conditions: conditions,
			Expected:   g.generateExpectedResult(taskStep),
		}

		steps = append(steps, step)
	}

	// 使用LLM优化步骤
	if len(steps) > 0 {
		optimizedSteps, err := g.optimizeStepsWithLLM(ctx, steps, context, options)
		if err != nil {
			fmt.Printf("Warning: step optimization failed: %v\n", err)
			// 使用原始步骤
		} else {
			steps = optimizedSteps
		}
	}

	return steps, nil
}

// generateStepDescription 生成步骤描述
func (g *Generator) generateStepDescription(taskStep TaskStep, stepNumber, totalSteps int) string {
	if taskStep.Action != "" {
		return taskStep.Action
	}

	// 根据工具和参数生成描述
	if taskStep.Tool != "" {
		return fmt.Sprintf("使用 %s 工具执行操作", taskStep.Tool)
	}

	return fmt.Sprintf("执行步骤 %d/%d", stepNumber, totalSteps)
}

// extractParameters 提取参数
func (g *Generator) extractParameters(taskStep TaskStep) []Parameter {
	var parameters []Parameter

	if taskStep.Parameters != nil {
		for name, value := range taskStep.Parameters {
			paramType := inferParameterType(value)
			parameter := Parameter{
				Name:        name,
				Type:        paramType,
				Description: fmt.Sprintf("参数 %s", name),
				Required:    true,
				Default:     value,
			}
			parameters = append(parameters, parameter)
		}
	}

	return parameters
}

// extractConditions 提取条件
func (g *Generator) extractConditions(taskStep TaskStep, record *TaskExecutionRecord) []Condition {
	var conditions []Condition

	// 添加前置步骤成功条件
	if taskStep.Order > 1 {
		conditions = append(conditions, Condition{
			Type:     "previous_step_success",
			Check:    fmt.Sprintf("step_%d_success", taskStep.Order-1),
			Expected: true,
		})
	}

	// 根据错误类型添加条件
	if !taskStep.Success && taskStep.Error != "" {
		// 分析错误类型并添加相应条件
		if strings.Contains(strings.ToLower(taskStep.Error), "permission") {
			conditions = append(conditions, Condition{
				Type:     "permission_check",
				Check:    "has_permission",
				Expected: true,
			})
		}
	}

	return conditions
}

// generateExpectedResult 生成预期结果
func (g *Generator) generateExpectedResult(taskStep TaskStep) string {
	if taskStep.Success && taskStep.Output != nil {
		// 简化输出描述
		return "操作成功完成"
	}

	if !taskStep.Success {
		return fmt.Sprintf("避免错误: %s", taskStep.Error)
	}

	return "成功完成指定操作"
}

// optimizeStepsWithLLM 使用LLM优化步骤
func (g *Generator) optimizeStepsWithLLM(ctx context.Context, steps []Step, context map[string]any, options GenerationOptions) ([]Step, error) {
	// 构建步骤文本
	var stepsText strings.Builder
	for _, step := range steps {
		stepsText.WriteString(fmt.Sprintf("步骤 %d: %s\n", step.Order, step.Action))
		if step.Tool != "" {
			stepsText.WriteString(fmt.Sprintf("  工具: %s\n", step.Tool))
		}
		if len(step.Parameters) > 0 {
			stepsText.WriteString("  参数:\n")
			for _, param := range step.Parameters {
				stepsText.WriteString(fmt.Sprintf("    - %s: %s\n", param.Name, param.Description))
			}
		}
		stepsText.WriteString("\n")
	}

	prompt := g.prompts.StepRefinement
	prompt = replaceTemplateVariables(prompt, map[string]string{
		"{{steps}}":            stepsText.String(),
		"{{task_description}}": context["task_description"].(string),
		"{{task_complexity}}":  context["task_complexity"].(string),
	})

	llmOptions := map[string]any{
		"temperature": options.Temperature,
		"max_tokens":  options.MaxTokens / 2, // 优化步骤使用较少token
		"model":       options.Model,
	}

	_, err := g.llm.Generate(ctx, prompt, llmOptions)
	if err != nil {
		return steps, err
	}

	// 解析优化后的步骤
	// 这里简化实现，实际应该解析LLM返回的结构化数据
	return steps, nil
}

// generateExamples 生成示例
func (g *Generator) generateExamples(ctx context.Context, skill *Skill, context map[string]any, options GenerationOptions) ([]string, error) {
	prompt := g.prompts.ExampleCreation
	prompt = replaceTemplateVariables(prompt, map[string]string{
		"{{skill_name}}":        skill.Name,
		"{{skill_description}}": skill.Description,
		"{{steps_count}}":       fmt.Sprintf("%d", len(skill.Steps)),
		"{{task_type}}":         context["task_type"].(string),
	})

	llmOptions := map[string]any{
		"temperature": options.Temperature,
		"max_tokens":  options.MaxTokens / 3,
		"model":       options.Model,
	}

	response, err := g.llm.Generate(ctx, prompt, llmOptions)
	if err != nil {
		return nil, err
	}

	// 解析示例（假设每行一个示例）
	lines := strings.Split(response, "\n")
	var examples []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "//") {
			examples = append(examples, line)
		}
	}

	// 限制示例数量
	if len(examples) > 3 {
		examples = examples[:3]
	}

	return examples, nil
}

// generateTips 生成技巧
func (g *Generator) generateTips(ctx context.Context, skill *Skill, context map[string]any, options GenerationOptions) ([]string, error) {
	// 基于技能内容生成技巧
	var tips []string

	// 添加通用技巧
	tips = append(tips, "确保所有前置条件满足后再执行")
	tips = append(tips, "仔细检查参数格式和类型")

	// 基于工具添加特定技巧
	for _, tool := range skill.RelatedTools {
		switch tool {
		case "fs_write", "fs_create":
			tips = append(tips, "写入文件前检查目录权限")
		case "shell_exec":
			tips = append(tips, "执行命令前验证命令安全性")
		case "git_commit":
			tips = append(tips, "提交前检查变更内容")
		}
	}

	// 基于复杂度添加技巧
	if skill.Complexity == "complex" {
		tips = append(tips, "复杂任务建议分步执行并验证中间结果")
	}

	return tips, nil
}

// generateWarnings 生成警告
func (g *Generator) generateWarnings(ctx context.Context, skill *Skill, context map[string]any, options GenerationOptions) ([]string, error) {
	// 基于任务记录中的错误生成警告
	var warnings []string

	taskSuccess := context["success"].(bool)
	if !taskSuccess {
		errorMsg := ""
		if context["error_message"] != nil {
			errorMsg = context["error_message"].(string)
		}

		if errorMsg != "" {
			warnings = append(warnings, fmt.Sprintf("注意: 原始任务执行失败，错误: %s", errorMsg))
		}
	}

	// 添加工具相关警告
	for _, tool := range skill.RelatedTools {
		switch tool {
		case "fs_delete", "fs_remove":
			warnings = append(warnings, "删除操作不可逆，请谨慎执行")
		case "shell_exec":
			warnings = append(warnings, "命令执行可能影响系统安全，确保命令来源可信")
		}
	}

	// 添加执行时间警告
	executionTime := context["execution_time"].(string)
	if strings.Contains(executionTime, "m") || strings.Contains(executionTime, "h") {
		warnings = append(warnings, "此技能执行时间较长，请确保有足够时间")
	}

	return warnings, nil
}

// extractRelatedTools 提取相关工具
func (g *Generator) extractRelatedTools(record *TaskExecutionRecord) []string {
	tools := make(map[string]bool)

	// 从工具使用记录中提取
	for toolName := range record.ToolUsage {
		tools[toolName] = true
	}

	// 从步骤中提取
	for _, step := range record.Steps {
		if step.Tool != "" {
			tools[step.Tool] = true
		}
	}

	// 转换为切片
	var toolList []string
	for tool := range tools {
		toolList = append(toolList, tool)
	}

	return toolList
}

// estimateTokens 估算Token数量
func (g *Generator) estimateTokens(skill *Skill) int {
	// 简单估算：基于文本长度
	text := skill.Name + skill.Description
	for _, step := range skill.Steps {
		text += step.Action
		for _, param := range step.Parameters {
			text += param.Name + param.Description
		}
	}

	// 粗略估算：1个中文字符约2个token，1个英文字符约0.25个token
	chineseChars := countChineseChars(text)
	englishChars := len(text) - chineseChars

	tokens := chineseChars*2 + englishChars/4

	// 添加缓冲
	return tokens + 500
}

// estimateExecutionTime 估算执行时间
func (g *Generator) estimateExecutionTime(record *TaskExecutionRecord) time.Duration {
	if record.Success {
		return record.EndTime.Sub(record.StartTime)
	}

	// 如果失败，基于步骤数量估算
	avgStepTime := time.Second * 30
	return time.Duration(len(record.Steps)) * avgStepTime
}

// inferParameterType 推断参数类型
func inferParameterType(value any) string {
	switch value.(type) {
	case string:
		return "string"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "number"
	case float32, float64:
		return "number"
	case bool:
		return "boolean"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return "string"
	}
}

// cleanSkillName 清理技能名称
func cleanSkillName(name string) string {
	// 移除特殊字符
	var cleaned strings.Builder
	for _, ch := range name {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == ' ' || ch == '-' || ch == '_' {
			cleaned.WriteRune(ch)
		}
	}

	result := cleaned.String()

	// 限制长度
	if len(result) > 50 {
		result = result[:50]
	}

	return strings.TrimSpace(result)
}

// countChineseChars 计算中文字符数量
func countChineseChars(text string) int {
	count := 0
	for _, ch := range text {
		// 简单判断：Unicode中文字符范围
		if (ch >= 0x4E00 && ch <= 0x9FFF) ||
			(ch >= 0x3400 && ch <= 0x4DBF) ||
			(ch >= 0x20000 && ch <= 0x2A6DF) {
			count++
		}
	}
	return count
}

// replaceTemplateVariables 替换模板变量
func replaceTemplateVariables(template string, variables map[string]string) string {
	result := template
	for key, value := range variables {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}

// getDefaultPrompts 获取默认提示词模板
func getDefaultPrompts() PromptTemplates {
	return PromptTemplates{
		SkillExtraction: `基于以下任务执行记录，生成一个技能的名称和简要描述。

任务描述: {{task_description}}
任务类型: {{task_type}}
任务复杂度: {{task_complexity}}
执行结果: {{success}}
参与Agent: {{agent_involved}}
触发原因: {{trigger_reason}}

请生成:
1. 技能名称（简洁明了，体现技能核心功能）
2. 技能描述（1-2句话说明技能用途）

格式: 技能名称|技能描述

示例: 文件创建技能|在指定路径创建文件并设置权限`,

		StepRefinement: `优化以下技能步骤，使其更加清晰、可执行。

原始步骤:
{{steps}}

任务背景: {{task_description}}
任务复杂度: {{task_complexity}}

请优化步骤描述，确保:
1. 每个步骤有明确的操作目标
2. 参数说明清晰
3. 条件判断合理
4. 预期结果明确`,

		ExampleCreation: `为以下技能生成使用示例:

技能名称: {{skill_name}}
技能描述: {{skill_description}}
步骤数量: {{steps_count}}
任务类型: {{task_type}}

请生成2-3个具体的使用示例，展示技能在不同场景下的应用。`,

		Validation: `验证以下技能是否安全有效:

技能名称: {{skill_name}}
技能描述: {{skill_description}}
步骤: {{steps}}

请检查:
1. 是否有安全风险（如删除文件、执行命令）
2. 步骤是否逻辑清晰
3. 参数是否完整
4. 条件是否合理`,
	}
}

// SetPrompts 设置提示词模板
func (g *Generator) SetPrompts(prompts PromptTemplates) {
	g.prompts = prompts
}

// DefaultGenerationOptions 获取默认生成选项
func DefaultGenerationOptions() GenerationOptions {
	return GenerationOptions{
		Model:           "gpt-4",
		Temperature:     0.7,
		MaxTokens:       2000,
		IncludeExamples: true,
		IncludeTips:     true,
		IncludeWarnings: true,
		Category:        "general",
		Tags:            []string{"auto-generated"},
		Complexity:      "medium",
	}
}
