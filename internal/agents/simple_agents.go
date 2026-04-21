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

// SimpleAgent 简化版Agent基类
type SimpleAgent struct {
	*BaseAgent
}

// NewSimpleAgent 创建简化版Agent
func NewSimpleAgent(config AgentConfig, modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine) (*SimpleAgent, error) {
	// 创建基础Agent
	baseAgent, err := NewBaseAgent(config, modelAdapter, toolRegistry, flowEngine)
	if err != nil {
		return nil, err
	}

	// 设置默认提示词模板
	if config.PromptTemplates == nil {
		config.PromptTemplates = make(map[string]string)
	}

	// 根据角色设置不同的默认提示词
	defaultTemplate := getDefaultPromptTemplate(config.Role)
	if _, exists := config.PromptTemplates["default"]; !exists {
		config.PromptTemplates["default"] = defaultTemplate
	}

	agent := &SimpleAgent{
		BaseAgent: baseAgent,
	}

	return agent, nil
}

// getDefaultPromptTemplate 根据角色获取默认提示词模板
func getDefaultPromptTemplate(role AgentRole) string {
	switch role {
	case RoleOrchestrator:
		return `你是一个项目编排器。请分析并处理以下任务：

任务：{{task}}
参数：{{parameters}}

请按照你的专业角色（项目编排器）进行分析和处理，输出详细的分析结果和执行建议。`

	case RoleProductManager:
		return `你是一个产品经理。请分析并处理以下任务：

任务：{{task}}
参数：{{parameters}}

请按照你的专业角色（产品经理）进行分析和处理，重点关注需求分析、用户故事、功能定义和产品规划。`

	case RoleArchitect:
		return `你是一个技术架构师。请分析并处理以下任务：

任务：{{task}}
参数：{{parameters}}

请按照你的专业角色（技术架构师）进行分析和处理，重点关注技术架构、系统设计、数据库设计和API设计。`

	case RoleFrontendDev:
		return `你是一个前端开发工程师。请分析并处理以下任务：

任务：{{task}}
参数：{{parameters}}

请按照你的专业角色（前端开发工程师）进行分析和处理，重点关注React/Vue组件开发、UI实现、状态管理和API集成。`

	case RoleBackendDev:
		return `你是一个后端开发工程师。请分析并处理以下任务：

任务：{{task}}
参数：{{parameters}}

请按照你的专业角色（后端开发工程师）进行分析和处理，重点关注API开发、数据库操作、业务逻辑和系统集成。`

	case RoleQAEngineer:
		return `你是一个QA工程师。请分析并处理以下任务：

任务：{{task}}
参数：{{parameters}}

请按照你的专业角色（QA工程师）进行分析和处理，重点关注测试用例设计、测试执行、质量保证和缺陷管理。`

	case RoleDevOps:
		return `你是一个DevOps工程师。请分析并处理以下任务：

任务：{{task}}
参数：{{parameters}}

请按照你的专业角色（DevOps工程师）进行分析和处理，重点关注项目构建、部署、监控和运维自动化。`

	case RoleReviewer:
		return `你是一个代码审查员。请分析并处理以下任务：

任务：{{task}}
参数：{{parameters}}

请按照你的专业角色（代码审查员）进行分析和处理，重点关注代码质量、安全漏洞、性能问题和最佳实践。`

	default:
		return `请分析并处理以下任务：

任务：{{task}}
参数：{{parameters}}

请提供详细的分析结果和执行建议。`
	}
}

// Execute 执行简化版Agent任务
func (a *SimpleAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
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

	utils.Info("Agent %s (%s) 开始执行任务: %s", a.config.Name, a.config.Role, task)

	startTime := time.Now()

	// 格式化提示词
	promptData := map[string]interface{}{
		"task":       task,
		"parameters": parameters,
	}

	prompt, err := a.FormatPrompt("default", promptData)
	if err != nil {
		prompt = fmt.Sprintf("请分析并处理以下任务：\n任务：%s\n参数：%v", task, parameters)
	}

	// 调用模型分析任务
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		a.execution.State = StateFailed
		a.execution.Error = fmt.Sprintf("任务分析失败: %v", err)
		a.state = StateFailed
		a.execution.CompletedAt = &time.Time{}
		*a.execution.CompletedAt = time.Now()

		return &AgentResult{
			Success:  false,
			Output:   fmt.Sprintf("任务分析失败: %v", err),
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

	analysis := modelResponse.Content

	result := &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"task_analyzed":  true,
			"analysis":       analysis,
			"model_response": modelResponse,
			"agent_role":     a.config.Role,
			"agent_name":     a.config.Name,
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "task_analysis",
				Name:        "任务分析",
				Content:     analysis,
				Description: fmt.Sprintf("%s对任务的分析和建议", a.config.Name),
			},
		},
		Duration: duration,
		Logs: []string{
			fmt.Sprintf("任务分析完成"),
			fmt.Sprintf("Agent角色: %s", a.config.Role),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
			fmt.Sprintf("执行时间: %v", duration),
		},
	}

	a.execution.Result = result

	return result, nil
}

// ValidateTask 验证任务（简化版Agent特定验证）
func (a *SimpleAgent) ValidateTask(task string) error {
	// 调用基础验证
	if err := a.BaseAgent.ValidateTask(task); err != nil {
		return err
	}

	// 根据角色进行特定验证
	taskLower := strings.ToLower(task)

	// 检查是否包含危险内容（所有Agent通用）
	dangerousPatterns := []string{
		"rm -rf", "format c:", "delete all", "destroy", "drop database",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(taskLower, pattern) {
			return fmt.Errorf("task contains dangerous pattern: %s", pattern)
		}
	}

	return nil
}

// GetRoleDescription 获取角色描述
func (a *SimpleAgent) GetRoleDescription() string {
	switch a.config.Role {
	case RoleOrchestrator:
		return "项目编排器：负责任务分解、执行计划生成和Agent调度"
	case RoleProductManager:
		return "产品经理：负责需求分析、PRD文档生成和产品规划"
	case RoleArchitect:
		return "技术架构师：负责技术架构设计、技术方案输出和系统设计"
	case RoleFrontendDev:
		return "前端开发工程师：负责React/Vue组件开发和前端界面实现"
	case RoleBackendDev:
		return "后端开发工程师：负责后端API、数据库代码生成和业务逻辑实现"
	case RoleQAEngineer:
		return "QA工程师：负责测试用例生成与执行、质量保证"
	case RoleDevOps:
		return "DevOps工程师：负责项目构建、启动、部署验证和运维"
	case RoleReviewer:
		return "代码审查员：负责代码审查、安全扫描和质量检查"
	default:
		return "通用Agent：处理各种任务"
	}
}
