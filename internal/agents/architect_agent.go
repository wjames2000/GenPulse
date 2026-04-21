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

// ArchitectAgent 技术架构师Agent
type ArchitectAgent struct {
	*BaseAgent
}

// NewArchitectAgent 创建技术架构师Agent
func NewArchitectAgent(config AgentConfig, modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine) (*ArchitectAgent, error) {
	// 设置技术架构师Agent的默认配置
	config.Role = RoleArchitect
	if config.Description == "" {
		config.Description = "技术架构师，负责技术架构设计、技术方案输出和系统设计"
	}

	// 设置默认能力
	if len(config.Capabilities) == 0 {
		config.Capabilities = []AgentCapability{
			CapabilityPlanning,
			CapabilityDesign,
			CapabilityAnalysis,
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
		"architecture_design": `你是一个技术架构师。请为以下需求设计技术架构：

项目名称：{{project_name}}
项目描述：{{project_description}}
技术需求：{{technical_requirements}}
非功能需求：{{non_functional_requirements}}

请按照以下步骤进行架构设计：
1. 架构目标：明确架构设计的目标和原则
2. 技术选型：选择合适的技术栈（前端、后端、数据库、中间件等）
3. 架构风格：确定架构风格（微服务、单体、分层等）
4. 组件设计：设计系统组件和模块划分
5. 数据设计：设计数据库结构和数据流
6. API设计：设计API接口和通信协议
7. 部署架构：设计部署架构和基础设施
8. 安全设计：设计安全策略和防护措施
9. 性能设计：设计性能优化方案
10. 监控设计：设计监控和运维方案

请输出详细的技术架构设计方案：`,

		"system_design": `你是一个技术架构师。请为以下系统进行详细设计：

系统名称：{{system_name}}
系统描述：{{system_description}}
核心功能：{{core_features}}
用户规模：{{user_scale}}
数据规模：{{data_scale}}

请提供以下设计内容：
1. 系统架构图（用文字描述）
2. 核心模块设计
3. 数据库设计（表结构、索引、关系）
4. API设计（接口定义、参数、返回值）
5. 缓存设计（缓存策略、缓存结构）
6. 消息队列设计（消息类型、队列结构）
7. 安全设计（认证、授权、加密）
8. 性能设计（并发处理、响应时间、吞吐量）
9. 扩展性设计（水平扩展、垂直扩展）
10. 容错设计（故障恢复、数据备份）

请开始系统设计：`,

		"technical_solution": `你是一个技术架构师。请为以下问题提供技术解决方案：

问题描述：{{problem_description}}
业务场景：{{business_scenario}}
技术约束：{{technical_constraints}}
现有系统：{{existing_system}}

请提供以下内容：
1. 问题分析：分析问题的根本原因和影响范围
2. 解决方案：提出具体的技术解决方案
3. 技术选型：推荐合适的技术和工具
4. 实施步骤：详细的实施步骤和时间计划
5. 风险评估：识别潜在风险并提供应对措施
6. 成本评估：评估实施成本（人力、时间、资源）
7. 效果评估：预期效果和收益评估

请输出详细的技术解决方案：`,
	}

	// 合并提示词模板
	for key, template := range defaultTemplates {
		if _, exists := config.PromptTemplates[key]; !exists {
			config.PromptTemplates[key] = template
		}
	}

	// 创建基础Agent
	baseAgent, err := NewBaseAgent(config, modelAdapter, toolRegistry, flowEngine)
	if err != nil {
		return nil, err
	}

	agent := &ArchitectAgent{
		BaseAgent: baseAgent,
	}

	return agent, nil
}

// Execute 执行技术架构师Agent任务
func (a *ArchitectAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
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

	utils.Info("技术架构师Agent %s 开始执行任务: %s", a.config.Name, task)

	startTime := time.Now()

	// 根据任务类型选择提示词模板
	templateName := "architecture_design"
	if strings.Contains(strings.ToLower(task), "系统设计") {
		templateName = "system_design"
	} else if strings.Contains(strings.ToLower(task), "技术方案") || strings.Contains(strings.ToLower(task), "解决方案") {
		templateName = "technical_solution"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"task":       task,
		"parameters": parameters,
	}

	// 添加特定参数
	if projectName, ok := parameters["project_name"]; ok {
		promptData["project_name"] = projectName
	}
	if projectDesc, ok := parameters["project_description"]; ok {
		promptData["project_description"] = projectDesc
	}
	if techReqs, ok := parameters["technical_requirements"]; ok {
		promptData["technical_requirements"] = techReqs
	}

	prompt, err := a.FormatPrompt(templateName, promptData)
	if err != nil {
		prompt = fmt.Sprintf("你是一个技术架构师。请为以下需求设计技术架构：\n\n需求：%s\n\n参数：%v\n\n请输出详细的技术架构设计方案。", task, parameters)
	}

	// 调用模型分析任务
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		a.execution.State = StateFailed
		a.execution.Error = fmt.Sprintf("架构设计失败: %v", err)
		a.state = StateFailed
		a.execution.CompletedAt = &time.Time{}
		*a.execution.CompletedAt = time.Now()

		return &AgentResult{
			Success:  false,
			Output:   fmt.Sprintf("架构设计失败: %v", err),
			Duration: time.Since(startTime),
		}, nil
	}

	duration := time.Since(startTime)

	// 更新执行上下文
	a.execution.CompletedAt = &time.Time{}
	*a.execution.CompletedAt = time.Now()
	a.execution.State = StateCompleted
	a.state = StateIdle
	a.successCount++
	a.totalDuration += duration

	architectureDesign := modelResponse.Content

	// 尝试生成架构图描述
	architectureDiagram := a.generateArchitectureDiagram(architectureDesign)

	result := &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"task_analyzed":        true,
			"architecture_design":  architectureDesign,
			"architecture_diagram": architectureDiagram,
			"model_response":       modelResponse,
			"agent_role":           a.config.Role,
			"agent_name":           a.config.Name,
			"template_used":        templateName,
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "architecture_design",
				Name:        "技术架构设计方案",
				Content:     architectureDesign,
				Description: fmt.Sprintf("%s的技术架构设计方案", a.config.Name),
			},
			{
				Type:        "architecture_diagram",
				Name:        "架构图描述",
				Content:     architectureDiagram,
				Description: "系统架构图文字描述",
			},
		},
		Duration: duration,
		Logs: []string{
			fmt.Sprintf("架构设计完成"),
			fmt.Sprintf("使用模板: %s", templateName),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
			fmt.Sprintf("执行时间: %v", duration),
		},
	}

	a.execution.Result = result

	return result, nil
}

// generateArchitectureDiagram 从架构设计生成架构图描述
func (a *ArchitectAgent) generateArchitectureDiagram(architectureDesign string) string {
	// 简化的架构图生成逻辑
	// 在实际实现中，可以调用专门的图表生成服务或使用图表库

	diagram := "```mermaid\ngraph TD\n"

	// 提取关键组件
	lines := strings.Split(architectureDesign, "\n")
	components := []string{}

	for _, line := range lines {
		if strings.Contains(line, "组件") || strings.Contains(line, "模块") || strings.Contains(line, "服务") {
			// 简化的组件提取逻辑
			cleanLine := strings.TrimSpace(line)
			if len(cleanLine) > 0 && !strings.HasPrefix(cleanLine, "#") {
				components = append(components, cleanLine)
			}
		}
	}

	// 生成简单的架构图
	if len(components) > 0 {
		// 添加用户节点
		diagram += "    User[用户] --> Frontend[前端应用]\n"

		// 添加组件节点
		for i, comp := range components {
			if i < 5 { // 限制组件数量
				compName := fmt.Sprintf("Comp%d[%s]", i+1, comp)
				if i == 0 {
					diagram += fmt.Sprintf("    Frontend --> %s\n", compName)
				} else {
					diagram += fmt.Sprintf("    Comp%d --> %s\n", i, compName)
				}
			}
		}

		// 添加数据库节点
		diagram += "    Comp1 --> Database[(数据库)]\n"
		diagram += "    Comp2 --> Cache[(缓存)]\n"
	}

	diagram += "```"

	return diagram
}

// ValidateTask 验证任务（技术架构师特定验证）
func (a *ArchitectAgent) ValidateTask(task string) error {
	// 调用基础验证
	if err := a.BaseAgent.ValidateTask(task); err != nil {
		return err
	}

	// 技术架构师特定验证
	taskLower := strings.ToLower(task)

	// 检查是否包含架构设计相关关键词
	architectureKeywords := []string{
		"架构", "设计", "系统", "技术", "方案", "规划",
		"architecture", "design", "system", "technical", "solution", "plan",
	}

	hasArchitectureKeyword := false
	for _, keyword := range architectureKeywords {
		if strings.Contains(taskLower, keyword) {
			hasArchitectureKeyword = true
			break
		}
	}

	if !hasArchitectureKeyword {
		utils.Warn("任务可能不适用于技术架构师: %s", task)
	}

	return nil
}

// GetRoleDescription 获取角色描述
func (a *ArchitectAgent) GetRoleDescription() string {
	return "技术架构师：负责技术架构设计、技术方案输出和系统设计"
}

// DesignArchitecture 设计技术架构（专用方法）
func (a *ArchitectAgent) DesignArchitecture(ctx context.Context, projectName, projectDescription, requirements string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"project_name":           projectName,
		"project_description":    projectDescription,
		"technical_requirements": requirements,
	}

	task := fmt.Sprintf("为项目'%s'设计技术架构", projectName)
	return a.Execute(ctx, task, parameters)
}

// DesignSystem 设计系统（专用方法）
func (a *ArchitectAgent) DesignSystem(ctx context.Context, systemName, systemDescription, coreFeatures string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"system_name":        systemName,
		"system_description": systemDescription,
		"core_features":      coreFeatures,
	}

	task := fmt.Sprintf("为系统'%s'进行详细设计", systemName)
	return a.Execute(ctx, task, parameters)
}

// ProvideTechnicalSolution 提供技术解决方案（专用方法）
func (a *ArchitectAgent) ProvideTechnicalSolution(ctx context.Context, problemDescription, businessScenario string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"problem_description": problemDescription,
		"business_scenario":   businessScenario,
	}

	task := fmt.Sprintf("为问题'%s'提供技术解决方案", problemDescription)
	return a.Execute(ctx, task, parameters)
}
