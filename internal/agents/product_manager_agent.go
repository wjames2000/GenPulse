package agents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"GenPulse/internal/genkit/flows"
	"GenPulse/internal/genkit/models"
	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/utils"
)

// ProductManagerAgent 产品经理Agent
type ProductManagerAgent struct {
	*BaseAgent
}

// NewProductManagerAgent 创建产品经理Agent
func NewProductManagerAgent(config AgentConfig, modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine) (*ProductManagerAgent, error) {
	// 设置产品经理Agent的默认配置
	config.Role = RoleProductManager
	if config.Description == "" {
		config.Description = "产品经理，负责需求分析、PRD文档生成和产品规划"
	}

	// 设置默认能力
	if len(config.Capabilities) == 0 {
		config.Capabilities = []AgentCapability{
			CapabilityPlanning,
		}
	}

	// 设置默认工具
	if len(config.Tools) == 0 {
		config.Tools = []string{
			"fs_tool",
			"project_tool",
		}
	}

	// 设置默认提示词模板
	if config.PromptTemplates == nil {
		config.PromptTemplates = make(map[string]string)
	}

	// 添加默认提示词模板
	defaultTemplates := map[string]string{
		"requirement_analysis": `你是一个产品经理。请分析以下产品需求：

原始需求：{{raw_requirement}}
业务背景：{{business_context}}
目标用户：{{target_users}}
业务目标：{{business_goals}}

请按照以下结构进行需求分析：
1. 需求理解：用你自己的话重新描述需求，确保理解正确
2. 用户故事：从不同用户角度描述使用场景
3. 功能列表：列出所有需要实现的功能点
4. 非功能需求：性能、安全、可用性等要求
5. 优先级评估：根据业务价值和技术可行性评估优先级
6. 验收标准：每个功能的验收条件
7. 假设和约束：项目假设和技术约束

请输出完整的需求分析报告：`,

		"prd_generation": `你是一个产品经理。请根据以下需求分析生成PRD文档：

项目名称：{{project_name}}
需求分析：{{requirement_analysis}}
项目目标：{{project_goals}}
项目范围：{{project_scope}}

请生成完整的PRD文档，包含以下章节：
1. 文档概述（版本、作者、更新时间）
2. 项目背景（业务背景、问题陈述、机会）
3. 项目目标（业务目标、成功指标）
4. 目标用户（用户画像、使用场景）
5. 功能需求（详细功能描述、用户故事）
6. 非功能需求（性能、安全、可用性、兼容性）
7. 项目范围（包含内容、排除内容）
8. 项目假设和约束
9. 验收标准
10. 项目里程碑
11. 风险分析
12. 附录

请使用专业的PRD文档格式：`,

		"user_story_creation": `你是一个产品经理。请为以下功能创建用户故事：

功能描述：{{feature_description}}
目标用户：{{target_users}}
使用场景：{{usage_scenarios}}

请按照以下格式创建用户故事：
作为 [用户角色]
我希望 [实现某个功能]
以便于 [达到某个业务价值]

验收标准：
- [条件1]
- [条件2]
- [条件3]

请为每个主要功能创建至少3个用户故事：`,

		"feature_prioritization": `你是一个产品经理。请对以下功能进行优先级排序：

功能列表：{{feature_list}}
业务目标：{{business_goals}}
技术约束：{{technical_constraints}}
资源限制：{{resource_limitations}}

请使用以下方法进行优先级评估：
1. 业务价值评估（高、中、低）
2. 技术复杂度评估（高、中、低）
3. 用户影响评估（影响用户数量、使用频率）
4. 依赖关系分析（功能间的依赖）
5. 风险评估（技术风险、市场风险）

请输出优先级排序结果，建议使用MoSCoW方法（Must have, Should have, Could have, Won't have）：`,

		"acceptance_criteria": `你是一个产品经理。请为以下功能定义验收标准：

功能名称：{{feature_name}}
功能描述：{{feature_description}}
用户故事：{{user_stories}}

请定义清晰、可测试的验收标准：
1. 功能完整性：功能是否完整实现
2. 用户体验：用户界面和交互是否符合预期
3. 性能要求：响应时间、并发能力等
4. 安全性：数据安全、访问控制等
5. 兼容性：浏览器、设备、操作系统兼容性
6. 错误处理：异常情况下的处理方式
7. 数据验证：输入数据的验证规则

请为每个功能提供至少5条验收标准：`,

		"market_analysis": `你是一个产品经理。请进行市场分析：

产品概念：{{product_concept}}
目标市场：{{target_market}}
竞争对手：{{competitors}}

请分析：
1. 市场规模和增长趋势
2. 目标用户细分
3. 竞争对手分析（优势、劣势、机会、威胁）
4. 市场机会和空白点
5. 市场进入策略
6. 定价策略建议
7. 推广渠道建议

请输出市场分析报告：`,
	}

	for name, template := range defaultTemplates {
		if _, exists := config.PromptTemplates[name]; !exists {
			config.PromptTemplates[name] = template
		}
	}

	// 创建基础Agent
	baseAgent, err := NewBaseAgent(config, modelAdapter, toolRegistry, flowEngine)
	if err != nil {
		return nil, err
	}

	agent := &ProductManagerAgent{
		BaseAgent: baseAgent,
	}

	return agent, nil
}

// Execute 执行产品经理任务
func (a *ProductManagerAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	if !a.enabled {
		return nil, fmt.Errorf("agent is disabled")
	}

	if a.state != StateIdle {
		return nil, fmt.Errorf("agent is busy, current state: %s", a.state)
	}

	// 验证任务
	if err := a.ValidateTask(task); err != nil {
		return nil, fmt.Errorf("task validation failed: %w", err)
	}

	// 创建执行上下文
	executionID := fmt.Sprintf("%s-%d", a.config.ID, time.Now().Unix())
	a.execution = &AgentExecution{
		ID:         executionID,
		AgentID:    a.config.ID,
		Task:       task,
		Parameters: parameters,
		Context:    make(map[string]interface{}),
		StartedAt:  time.Now(),
		State:      StateThinking,
	}

	a.state = StateThinking
	a.executionCount++

	utils.Info("产品经理Agent %s 开始执行任务: %s", a.config.Name, task)

	startTime := time.Now()
	var result *AgentResult
	var execErr error

	// 根据任务类型执行不同的处理逻辑
	taskType := a.detectTaskType(task, parameters)

	switch taskType {
	case "requirement_analysis":
		result, execErr = a.handleRequirementAnalysis(ctx, task, parameters)
	case "prd_generation":
		result, execErr = a.handlePRDGeneration(ctx, task, parameters)
	case "user_story_creation":
		result, execErr = a.handleUserStoryCreation(ctx, task, parameters)
	case "feature_prioritization":
		result, execErr = a.handleFeaturePrioritization(ctx, task, parameters)
	case "acceptance_criteria":
		result, execErr = a.handleAcceptanceCriteria(ctx, task, parameters)
	case "market_analysis":
		result, execErr = a.handleMarketAnalysis(ctx, task, parameters)
	default:
		result, execErr = a.handleGenericProductTask(ctx, task, parameters)
	}

	duration := time.Since(startTime)

	// 更新执行上下文
	a.execution.CompletedAt = &time.Time{}
	*a.execution.CompletedAt = time.Now()

	if execErr != nil {
		a.execution.State = StateFailed
		a.execution.Error = execErr.Error()
		a.state = StateFailed
	} else {
		a.execution.State = StateCompleted
		a.execution.Result = result
		a.state = StateIdle

		if result.Success {
			a.successCount++
		}
		a.totalDuration += duration

		if result.Duration == 0 {
			result.Duration = duration
		}
	}

	return result, execErr
}

// detectTaskType 检测任务类型
func (a *ProductManagerAgent) detectTaskType(task string, parameters map[string]interface{}) string {
	taskLower := strings.ToLower(task)

	// 检查关键词
	keywords := map[string]string{
		"requirement_analysis":   "需求分析|分析需求|analyze requirement|requirement analysis",
		"prd_generation":         "生成PRD|PRD文档|prd|product requirements|需求文档",
		"user_story_creation":    "用户故事|user story|story creation|创建故事",
		"feature_prioritization": "功能优先级|优先级排序|prioritize|priority|排序功能",
		"acceptance_criteria":    "验收标准|acceptance criteria|验收条件|测试标准",
		"market_analysis":        "市场分析|竞品分析|market analysis|competitor analysis",
	}

	for taskType, pattern := range keywords {
		patterns := strings.Split(pattern, "|")
		for _, p := range patterns {
			if strings.Contains(taskLower, strings.ToLower(p)) {
				return taskType
			}
		}
	}

	// 检查参数
	if rawRequirement, ok := parameters["raw_requirement"].(string); ok && rawRequirement != "" {
		return "requirement_analysis"
	}

	if requirementAnalysis, ok := parameters["requirement_analysis"].(string); ok && requirementAnalysis != "" {
		return "prd_generation"
	}

	if featureDescription, ok := parameters["feature_description"].(string); ok && featureDescription != "" {
		return "user_story_creation"
	}

	if featureList, ok := parameters["feature_list"].(string); ok && featureList != "" {
		return "feature_prioritization"
	}

	return "generic_product"
}

// handleRequirementAnalysis 处理需求分析
func (a *ProductManagerAgent) handleRequirementAnalysis(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理需求分析: %s", task)

	// 获取参数
	rawRequirement, _ := parameters["raw_requirement"].(string)
	businessContext, _ := parameters["business_context"].(string)
	targetUsers, _ := parameters["target_users"].(string)
	businessGoals, _ := parameters["business_goals"].(string)

	if rawRequirement == "" {
		rawRequirement = task
	}

	if businessContext == "" {
		businessContext = "未提供业务背景"
	}

	if targetUsers == "" {
		targetUsers = "未指定目标用户"
	}

	if businessGoals == "" {
		businessGoals = "未明确业务目标"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"raw_requirement":  rawRequirement,
		"business_context": businessContext,
		"target_users":     targetUsers,
		"business_goals":   businessGoals,
	}

	prompt, err := a.FormatPrompt("requirement_analysis", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("请分析以下产品需求：\n原始需求：%s\n业务背景：%s\n目标用户：%s\n业务目标：%s",
			rawRequirement, businessContext, targetUsers, businessGoals)
	}

	// 调用模型进行需求分析
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("需求分析失败: %v", err),
		}, nil
	}

	requirementAnalysis := modelResponse.Content

	// 提取关键信息（简化处理）
	analysisSummary := a.extractAnalysisSummary(requirementAnalysis)

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"requirement_analyzed": true,
			"raw_requirement":      rawRequirement,
			"analysis":             requirementAnalysis,
			"analysis_summary":     analysisSummary,
			"model_response":       modelResponse,
			"analysis_structure":   "需求理解、用户故事、功能列表、非功能需求、优先级评估、验收标准、假设约束",
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "requirement_analysis",
				Name:        "需求分析报告",
				Content:     requirementAnalysis,
				Description: "产品需求的分析报告",
			},
			{
				Type:        "analysis_summary",
				Name:        "分析摘要",
				Content:     analysisSummary,
				Description: "需求分析的关键信息摘要",
			},
		},
		Logs: []string{
			"需求分析完成",
			fmt.Sprintf("原始需求: %s", rawRequirement),
			fmt.Sprintf("业务目标: %s", businessGoals),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// extractAnalysisSummary 提取分析摘要
func (a *ProductManagerAgent) extractAnalysisSummary(analysis string) map[string]interface{} {
	summary := make(map[string]interface{})

	// 简单的关键词提取（实际应用中可以使用更复杂的NLP处理）
	lines := strings.Split(analysis, "\n")

	var features []string
	var userStories []string
	var priorities []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 提取功能点
		if strings.Contains(line, "功能") || strings.Contains(line, "feature") || strings.Contains(line, "- ") {
			features = append(features, line)
		}

		// 提取用户故事
		if strings.Contains(line, "用户") || strings.Contains(line, "user") || strings.Contains(line, "作为") {
			userStories = append(userStories, line)
		}

		// 提取优先级
		if strings.Contains(line, "优先级") || strings.Contains(line, "priority") || strings.Contains(line, "重要") {
			priorities = append(priorities, line)
		}
	}

	// 限制数量
	if len(features) > 10 {
		features = features[:10]
	}
	if len(userStories) > 5 {
		userStories = userStories[:5]
	}
	if len(priorities) > 5 {
		priorities = priorities[:5]
	}

	summary["features"] = features
	summary["user_stories"] = userStories
	summary["priorities"] = priorities
	summary["total_features"] = len(features)
	summary["total_user_stories"] = len(userStories)

	return summary
}

// handlePRDGeneration 处理PRD文档生成
func (a *ProductManagerAgent) handlePRDGeneration(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理PRD文档生成: %s", task)

	// 获取参数
	projectName, _ := parameters["project_name"].(string)
	requirementAnalysis, _ := parameters["requirement_analysis"].(string)
	projectGoals, _ := parameters["project_goals"].(string)
	projectScope, _ := parameters["project_scope"].(string)

	if projectName == "" {
		projectName = "未命名项目"
	}

	if requirementAnalysis == "" {
		// 如果没有提供需求分析，先进行分析
		analysisResult, err := a.handleRequirementAnalysis(ctx, task, parameters)
		if err != nil || !analysisResult.Success {
			return &AgentResult{
				Success: false,
				Output:  "需要先进行需求分析",
			}, nil
		}

		if analysis, ok := analysisResult.Output.(map[string]interface{})["analysis"].(string); ok {
			requirementAnalysis = analysis
		}
	}

	if projectGoals == "" {
		projectGoals = "实现产品需求，满足用户期望"
	}

	if projectScope == "" {
		projectScope = "包含所有已识别的功能需求"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"project_name":         projectName,
		"requirement_analysis": requirementAnalysis,
		"project_goals":        projectGoals,
		"project_scope":        projectScope,
	}

	prompt, err := a.FormatPrompt("prd_generation", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("请生成PRD文档：\n项目名称：%s\n需求分析：%s\n项目目标：%s\n项目范围：%s",
			projectName, requirementAnalysis, projectGoals, projectScope)
	}

	// 调用模型生成PRD文档
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("PRD生成失败: %v", err),
		}, nil
	}

	prdDocument := modelResponse.Content

	// 保存PRD文档到文件
	var filePath string
	if projectPath, ok := parameters["project_path"].(string); ok && projectPath != "" {
		fileName := fmt.Sprintf("PRD_%s_%s.md", strings.ReplaceAll(projectName, " ", "_"), time.Now().Format("20060102"))
		filePath = fmt.Sprintf("%s/docs/%s", projectPath, fileName)

		// 使用文件工具保存
		toolResult, err := a.ExecuteTool(ctx, "fs_tool", map[string]interface{}{
			"operation": "write",
			"path":      filePath,
			"content":   prdDocument,
		})

		if err == nil && toolResult.Success {
			utils.Info("PRD文档已保存到: %s", filePath)
		} else {
			utils.Warn("保存PRD文档失败: %v", err)
		}
	}

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"prd_generated":        true,
			"project_name":         projectName,
			"prd_document":         prdDocument,
			"file_path":            filePath,
			"requirement_analysis": requirementAnalysis,
			"model_response":       modelResponse,
			"prd_structure":        "文档概述、项目背景、项目目标、目标用户、功能需求、非功能需求、项目范围、假设约束、验收标准、里程碑、风险分析",
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "prd_document",
				Name:        "PRD文档",
				Content:     prdDocument,
				Path:        filePath,
				Description: "产品需求文档",
			},
		},
		Logs: []string{
			"PRD文档生成完成",
			fmt.Sprintf("项目名称: %s", projectName),
			fmt.Sprintf("文档长度: %d 字符", len(prdDocument)),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// handleUserStoryCreation 处理用户故事创建
func (a *ProductManagerAgent) handleUserStoryCreation(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理用户故事创建: %s", task)

	// 获取参数
	featureDescription, _ := parameters["feature_description"].(string)
	targetUsers, _ := parameters["target_users"].(string)
	usageScenarios, _ := parameters["usage_scenarios"].(string)

	if featureDescription == "" {
		featureDescription = task
	}

	if targetUsers == "" {
		targetUsers = "系统用户"
	}

	if usageScenarios == "" {
		usageScenarios = "日常使用场景"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"feature_description": featureDescription,
		"target_users":        targetUsers,
		"usage_scenarios":     usageScenarios,
	}

	prompt, err := a.FormatPrompt("user_story_creation", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("请为以下功能创建用户故事：\n功能描述：%s\n目标用户：%s\n使用场景：%s",
			featureDescription, targetUsers, usageScenarios)
	}

	// 调用模型创建用户故事
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("用户故事创建失败: %v", err),
		}, nil
	}

	userStories := modelResponse.Content

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"user_stories_created": true,
			"feature_description":  featureDescription,
			"user_stories":         userStories,
			"target_users":         targetUsers,
			"model_response":       modelResponse,
			"story_format":         "作为[用户角色]、我希望[实现功能]、以便于[达到价值]",
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "user_stories",
				Name:        "用户故事",
				Content:     userStories,
				Description: "功能的用户故事描述",
			},
		},
		Logs: []string{
			"用户故事创建完成",
			fmt.Sprintf("功能: %s", featureDescription),
			fmt.Sprintf("目标用户: %s", targetUsers),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// handleFeaturePrioritization 处理功能优先级排序
func (a *ProductManagerAgent) handleFeaturePrioritization(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理功能优先级排序: %s", task)

	// 获取参数
	featureList, _ := parameters["feature_list"].(string)
	businessGoals, _ := parameters["business_goals"].(string)
	technicalConstraints, _ := parameters["technical_constraints"].(string)
	resourceLimitations, _ := parameters["resource_limitations"].(string)

	if featureList == "" {
		return &AgentResult{
			Success: false,
			Output:  "feature_list is required",
		}, nil
	}

	if businessGoals == "" {
		businessGoals = "实现产品价值，满足用户需求"
	}

	if technicalConstraints == "" {
		technicalConstraints = "无特殊技术约束"
	}

	if resourceLimitations == "" {
		resourceLimitations = "标准开发资源"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"feature_list":          featureList,
		"business_goals":        businessGoals,
		"technical_constraints": technicalConstraints,
		"resource_limitations":  resourceLimitations,
	}

	prompt, err := a.FormatPrompt("feature_prioritization", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("请对以下功能进行优先级排序：\n功能列表：%s\n业务目标：%s\n技术约束：%s\n资源限制：%s",
			featureList, businessGoals, technicalConstraints, resourceLimitations)
	}

	// 调用模型进行优先级排序
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("功能优先级排序失败: %v", err),
		}, nil
	}

	prioritizationResult := modelResponse.Content

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"features_prioritized":  true,
			"feature_list":          featureList,
			"prioritization":        prioritizationResult,
			"business_goals":        businessGoals,
			"model_response":        modelResponse,
			"prioritization_method": "MoSCoW方法（Must have, Should have, Could have, Won't have）",
			"evaluation_criteria":   "业务价值、技术复杂度、用户影响、依赖关系、风险评估",
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "feature_prioritization",
				Name:        "功能优先级排序",
				Content:     prioritizationResult,
				Description: "功能优先级排序结果",
			},
		},
		Logs: []string{
			"功能优先级排序完成",
			fmt.Sprintf("功能数量: %d", len(strings.Split(featureList, "\n"))),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// handleAcceptanceCriteria 处理验收标准定义
func (a *ProductManagerAgent) handleAcceptanceCriteria(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理验收标准定义: %s", task)

	// 获取参数
	featureName, _ := parameters["feature_name"].(string)
	featureDescription, _ := parameters["feature_description"].(string)
	userStories, _ := parameters["user_stories"].(string)

	if featureName == "" {
		featureName = "未命名功能"
	}

	if featureDescription == "" {
		featureDescription = task
	}

	if userStories == "" {
		userStories = "未提供用户故事"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"feature_name":        featureName,
		"feature_description": featureDescription,
		"user_stories":        userStories,
	}

	prompt, err := a.FormatPrompt("acceptance_criteria", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("请为以下功能定义验收标准：\n功能名称：%s\n功能描述：%s\n用户故事：%s",
			featureName, featureDescription, userStories)
	}

	// 调用模型定义验收标准
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("验收标准定义失败: %v", err),
		}, nil
	}

	acceptanceCriteria := modelResponse.Content

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"acceptance_criteria_defined": true,
			"feature_name":                featureName,
			"acceptance_criteria":         acceptanceCriteria,
			"feature_description":         featureDescription,
			"model_response":              modelResponse,
			"criteria_types":              "功能完整性、用户体验、性能要求、安全性、兼容性、错误处理、数据验证",
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "acceptance_criteria",
				Name:        "验收标准",
				Content:     acceptanceCriteria,
				Description: "功能的验收标准定义",
			},
		},
		Logs: []string{
			"验收标准定义完成",
			fmt.Sprintf("功能: %s", featureName),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// handleMarketAnalysis 处理市场分析
func (a *ProductManagerAgent) handleMarketAnalysis(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理市场分析: %s", task)

	// 获取参数
	productConcept, _ := parameters["product_concept"].(string)
	targetMarket, _ := parameters["target_market"].(string)
	competitors, _ := parameters["competitors"].(string)

	if productConcept == "" {
		productConcept = task
	}

	if targetMarket == "" {
		targetMarket = "未指定目标市场"
	}

	if competitors == "" {
		competitors = "未提供竞争对手信息"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"product_concept": productConcept,
		"target_market":   targetMarket,
		"competitors":     competitors,
	}

	prompt, err := a.FormatPrompt("market_analysis", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("请进行市场分析：\n产品概念：%s\n目标市场：%s\n竞争对手：%s",
			productConcept, targetMarket, competitors)
	}

	// 调用模型进行市场分析
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("市场分析失败: %v", err),
		}, nil
	}

	marketAnalysis := modelResponse.Content

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"market_analyzed":  true,
			"product_concept":  productConcept,
			"market_analysis":  marketAnalysis,
			"target_market":    targetMarket,
			"competitors":      competitors,
			"model_response":   modelResponse,
			"analysis_aspects": "市场规模、用户细分、竞争对手分析、市场机会、进入策略、定价策略、推广渠道",
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "market_analysis",
				Name:        "市场分析报告",
				Content:     marketAnalysis,
				Description: "产品市场分析报告",
			},
		},
		Logs: []string{
			"市场分析完成",
			fmt.Sprintf("产品概念: %s", productConcept),
			fmt.Sprintf("目标市场: %s", targetMarket),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// handleGenericProductTask 处理通用产品任务
func (a *ProductManagerAgent) handleGenericProductTask(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理通用产品任务: %s", task)

	// 使用模型分析任务
	promptData := map[string]interface{}{
		"task":       task,
		"parameters": parameters,
	}

	prompt, err := a.FormatPrompt("requirement_analysis", promptData)
	if err != nil {
		prompt = fmt.Sprintf("作为产品经理，请分析并处理以下任务：%s\n参数：%v", task, parameters)
	}

	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("任务分析失败: %v", err),
		}, nil
	}

	analysis := modelResponse.Content

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"task_analyzed":  true,
			"analysis":       analysis,
			"model_response": modelResponse,
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "product_analysis",
				Name:        "产品分析",
				Content:     analysis,
				Description: "产品经理对任务的分析和建议",
			},
		},
		Logs: []string{
			"通用产品任务完成",
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// ValidateTask 验证任务（产品经理Agent特定验证）
func (a *ProductManagerAgent) ValidateTask(task string) error {
	// 调用基础验证
	if err := a.BaseAgent.ValidateTask(task); err != nil {
		return err
	}

	// 产品经理Agent特定的验证
	taskLower := strings.ToLower(task)

	// 检查是否包含产品管理关键词
	productKeywords := []string{
		"需求", "产品", "用户", "市场", "功能", "优先级", "故事", "文档",
		"requirement", "product", "user", "market", "feature", "priority", "story", "document",
	}

	hasKeyword := false
	for _, keyword := range productKeywords {
		if strings.Contains(taskLower, keyword) {
			hasKeyword = true
			break
		}
	}

	if !hasKeyword {
		utils.Warn("产品经理任务可能不相关: %s", task)
		// 不返回错误，只是警告
	}

	return nil
}

// GetProductManagementCapabilities 获取产品管理能力描述
func (a *ProductManagerAgent) GetProductManagementCapabilities() string {
	capabilities := []string{
		"需求分析：深入分析用户需求，识别核心功能和价值主张",
		"PRD文档：生成完整的产品需求文档，包含功能描述和验收标准",
		"用户故事：创建详细的用户故事，描述功能使用场景和价值",
		"功能优先级：基于业务价值和技术可行性进行功能优先级排序",
		"验收标准：定义清晰、可测试的功能验收标准",
		"市场分析：分析市场规模、竞争对手和市场机会",
		"产品规划：制定产品路线图和发布计划",
	}

	return strings.Join(capabilities, "\n- ")
}
