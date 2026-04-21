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

// OrchestratorAgent 编排器Agent
type OrchestratorAgent struct {
	*BaseAgent
	agentManager *AgentManager
}

// NewOrchestratorAgent 创建编排器Agent
func NewOrchestratorAgent(config AgentConfig, modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine, agentManager *AgentManager) (*OrchestratorAgent, error) {
	// 设置编排器Agent的默认配置
	config.Role = RoleOrchestrator
	if config.Description == "" {
		config.Description = "项目编排器，负责任务分解、执行计划生成和Agent调度"
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
			"project_tool",
		}
	}

	// 设置默认提示词模板
	if config.PromptTemplates == nil {
		config.PromptTemplates = make(map[string]string)
	}

	// 添加默认提示词模板
	defaultTemplates := map[string]string{
		"task_decomposition": `你是一个项目编排器。请分析以下项目需求，并将其分解为具体的开发任务：

项目需求：{{requirement}}
项目类型：{{project_type}}
技术栈：{{tech_stack}}

请按照以下步骤进行任务分解：
1. 需求分析：理解用户的核心需求
2. 任务识别：识别需要完成的具体任务
3. 任务分类：将任务按类型分类（前端、后端、数据库、测试、部署等）
4. 依赖分析：分析任务间的依赖关系
5. 工作量估算：估算每个任务的工作量
6. Agent分配：为每个任务分配合适的Agent角色

请输出一个详细的任务分解计划，包括：
- 任务列表（每个任务包含：ID、名称、描述、类型、依赖任务、预估工作量、负责Agent）
- 执行顺序建议
- 并行执行机会
- 关键路径分析
- 风险评估

请开始你的分析：`,

		"execution_plan": `你是一个项目编排器。请根据以下任务分解生成具体的执行计划：

项目名称：{{project_name}}
任务分解：{{task_decomposition}}
可用Agent：{{available_agents}}

请生成一个详细的执行计划，包括：
1. 阶段划分：将任务划分为不同的执行阶段
2. 时间安排：为每个阶段和任务安排时间
3. Agent调度：指定每个任务的执行Agent
4. 资源分配：考虑Agent的并发限制
5. 检查点：设置关键检查点和里程碑
6. 风险应对：制定风险应对策略

请输出执行计划，格式为：
- 阶段1：需求分析阶段
  - 任务1.1：需求分析（PM Agent，2小时）
  - 任务1.2：PRD文档生成（PM Agent，1小时）
- 阶段2：架构设计阶段
  - 任务2.1：技术架构设计（Architect Agent，3小时）
  - 任务2.2：数据库设计（Architect Agent，2小时）
- ...以此类推

请开始生成执行计划：`,

		"agent_scheduling": `你是一个项目编排器。请为以下任务分配合适的Agent：

任务列表：{{task_list}}
可用Agent状态：{{agent_status}}

请考虑以下因素进行分配：
1. Agent的专业能力匹配
2. Agent的当前负载
3. 任务的优先级
4. 任务间的依赖关系
5. Agent的并发限制

请输出Agent分配方案，格式为：
- 任务1：需求分析 -> PM Agent（理由：PM Agent擅长需求分析和文档生成）
- 任务2：前端界面开发 -> Frontend Agent（理由：Frontend Agent擅长React/Vue开发）
- 任务3：后端API开发 -> Backend Agent（理由：Backend Agent擅长Go API开发）
- ...以此类推

请开始分配：`,

		"progress_monitoring": `你是一个项目编排器。请监控以下项目进度：

项目名称：{{project_name}}
执行计划：{{execution_plan}}
当前状态：{{current_status}}
已完成任务：{{completed_tasks}}
进行中任务：{{in_progress_tasks}}
待完成任务：{{pending_tasks}}

请分析：
1. 进度评估：当前进度是否符合计划
2. 问题识别：识别进度延迟或问题
3. 调整建议：如有需要，提出计划调整建议
4. 资源优化：优化Agent分配和资源使用
5. 风险更新：更新风险评估

请输出进度监控报告：`,

		"coordination": `你是一个项目编排器。请协调以下Agent间的协作：

协作场景：{{scenario}}
涉及Agent：{{involved_agents}}
协作目标：{{goal}}
当前问题：{{current_issue}}

请提供协调方案：
1. 沟通机制：建立有效的沟通渠道
2. 接口定义：明确Agent间的接口和数据格式
3. 同步点：设置必要的同步检查点
4. 冲突解决：制定冲突解决机制
5. 质量保证：确保协作质量

请输出协调方案：`,
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

	agent := &OrchestratorAgent{
		BaseAgent:    baseAgent,
		agentManager: agentManager,
	}

	return agent, nil
}

// Execute 执行编排器任务
func (a *OrchestratorAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
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

	utils.Info("编排器Agent %s 开始执行任务: %s", a.config.Name, task)

	startTime := time.Now()
	var result *AgentResult
	var execErr error

	// 根据任务类型执行不同的处理逻辑
	taskType := a.detectTaskType(task, parameters)

	switch taskType {
	case "task_decomposition":
		result, execErr = a.handleTaskDecomposition(ctx, task, parameters)
	case "execution_plan":
		result, execErr = a.handleExecutionPlan(ctx, task, parameters)
	case "agent_scheduling":
		result, execErr = a.handleAgentScheduling(ctx, task, parameters)
	case "progress_monitoring":
		result, execErr = a.handleProgressMonitoring(ctx, task, parameters)
	case "coordination":
		result, execErr = a.handleCoordination(ctx, task, parameters)
	default:
		result, execErr = a.handleGenericOrchestration(ctx, task, parameters)
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
func (a *OrchestratorAgent) detectTaskType(task string, parameters map[string]interface{}) string {
	taskLower := strings.ToLower(task)

	// 检查关键词
	keywords := map[string]string{
		"task_decomposition":  "分解任务|任务分解|breakdown|decompose|split tasks",
		"execution_plan":      "执行计划|制定计划|plan|schedule|timeline",
		"agent_scheduling":    "分配Agent|调度Agent|assign|schedule agents|allocate",
		"progress_monitoring": "监控进度|进度跟踪|monitor|track|progress",
		"coordination":        "协调|协作|coordinate|collaborate|sync",
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
	if requirement, ok := parameters["requirement"].(string); ok && requirement != "" {
		return "task_decomposition"
	}

	if taskDecomposition, ok := parameters["task_decomposition"].(string); ok && taskDecomposition != "" {
		return "execution_plan"
	}

	if taskList, ok := parameters["task_list"].(string); ok && taskList != "" {
		return "agent_scheduling"
	}

	return "generic_orchestration"
}

// handleTaskDecomposition 处理任务分解
func (a *OrchestratorAgent) handleTaskDecomposition(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理任务分解: %s", task)

	// 获取参数
	requirement, _ := parameters["requirement"].(string)
	projectType, _ := parameters["project_type"].(string)
	techStack, _ := parameters["tech_stack"].(string)

	if requirement == "" {
		requirement = task
	}

	if projectType == "" {
		projectType = "fullstack"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"requirement":  requirement,
		"project_type": projectType,
		"tech_stack":   techStack,
	}

	prompt, err := a.FormatPrompt("task_decomposition", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("请分解以下项目需求为具体任务：\n需求：%s\n项目类型：%s\n技术栈：%s", requirement, projectType, techStack)
	}

	// 调用模型进行任务分解
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("任务分解失败: %v", err),
		}, nil
	}

	// 解析任务分解结果
	decomposition := modelResponse.Content

	// 获取可用Agent信息
	var availableAgentsInfo string
	if a.agentManager != nil {
		agentStatus := a.agentManager.GetAllAgentsStatus()
		if agents, ok := agentStatus["agents"].(map[string]interface{}); ok {
			for agentID, status := range agents {
				if statusMap, ok := status.(map[string]interface{}); ok {
					name, _ := statusMap["name"].(string)
					role, _ := statusMap["role"].(string)
					enabled, _ := statusMap["enabled"].(bool)
					state, _ := statusMap["state"].(string)

					if enabled {
						availableAgentsInfo += fmt.Sprintf("- %s (%s): %s, 状态: %s\n", name, agentID, role, state)
					}
				}
			}
		}
	}

	if availableAgentsInfo == "" {
		availableAgentsInfo = "无可用Agent信息"
	}

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"task_decomposed":      true,
			"requirement":          requirement,
			"project_type":         projectType,
			"decomposition":        decomposition,
			"available_agents":     availableAgentsInfo,
			"model_response":       modelResponse,
			"decomposition_format": "任务列表、执行顺序、并行机会、关键路径、风险评估",
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "task_decomposition",
				Name:        "任务分解报告",
				Content:     decomposition,
				Description: "项目需求的任务分解结果",
			},
		},
		Logs: []string{
			"任务分解完成",
			fmt.Sprintf("需求: %s", requirement),
			fmt.Sprintf("项目类型: %s", projectType),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// handleExecutionPlan 处理执行计划生成
func (a *OrchestratorAgent) handleExecutionPlan(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理执行计划生成: %s", task)

	// 获取参数
	projectName, _ := parameters["project_name"].(string)
	taskDecomposition, _ := parameters["task_decomposition"].(string)
	availableAgents, _ := parameters["available_agents"].(string)

	if projectName == "" {
		projectName = "未命名项目"
	}

	if taskDecomposition == "" {
		// 如果没有提供任务分解，先进行任务分解
		decompositionResult, err := a.handleTaskDecomposition(ctx, task, parameters)
		if err != nil || !decompositionResult.Success {
			return &AgentResult{
				Success: false,
				Output:  "需要先进行任务分解",
			}, nil
		}

		if decomposition, ok := decompositionResult.Output.(map[string]interface{})["decomposition"].(string); ok {
			taskDecomposition = decomposition
		}
	}

	if availableAgents == "" && a.agentManager != nil {
		agentStatus := a.agentManager.GetAllAgentsStatus()
		if agents, ok := agentStatus["agents"].(map[string]interface{}); ok {
			for agentID, status := range agents {
				if statusMap, ok := status.(map[string]interface{}); ok {
					name, _ := statusMap["name"].(string)
					role, _ := statusMap["role"].(string)
					enabled, _ := statusMap["enabled"].(bool)

					if enabled {
						availableAgents += fmt.Sprintf("- %s (%s): %s\n", name, agentID, role)
					}
				}
			}
		}
	}

	if availableAgents == "" {
		availableAgents = "无可用Agent信息"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"project_name":       projectName,
		"task_decomposition": taskDecomposition,
		"available_agents":   availableAgents,
	}

	prompt, err := a.FormatPrompt("execution_plan", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("请为以下项目生成执行计划：\n项目名称：%s\n任务分解：%s\n可用Agent：%s", projectName, taskDecomposition, availableAgents)
	}

	// 调用模型生成执行计划
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("执行计划生成失败: %v", err),
		}, nil
	}

	executionPlan := modelResponse.Content

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"execution_plan_generated": true,
			"project_name":             projectName,
			"execution_plan":           executionPlan,
			"task_decomposition":       taskDecomposition,
			"model_response":           modelResponse,
			"plan_format":              "阶段划分、时间安排、Agent调度、资源分配、检查点、风险应对",
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "execution_plan",
				Name:        "项目执行计划",
				Content:     executionPlan,
				Description: "项目的详细执行计划",
			},
		},
		Logs: []string{
			"执行计划生成完成",
			fmt.Sprintf("项目名称: %s", projectName),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// handleAgentScheduling 处理Agent调度
func (a *OrchestratorAgent) handleAgentScheduling(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理Agent调度: %s", task)

	// 获取参数
	taskList, _ := parameters["task_list"].(string)
	agentStatus, _ := parameters["agent_status"].(string)

	if taskList == "" {
		return &AgentResult{
			Success: false,
			Output:  "task_list is required",
		}, nil
	}

	if agentStatus == "" && a.agentManager != nil {
		// 获取实际Agent状态
		allStatus := a.agentManager.GetAllAgentsStatus()
		if agents, ok := allStatus["agents"].(map[string]interface{}); ok {
			for agentID, status := range agents {
				if statusMap, ok := status.(map[string]interface{}); ok {
					name, _ := statusMap["name"].(string)
					role, _ := statusMap["role"].(string)
					enabled, _ := statusMap["enabled"].(bool)
					state, _ := statusMap["state"].(string)
					executionCount, _ := statusMap["statistics"].(map[string]interface{})["execution_count"].(int)

					agentStatus += fmt.Sprintf("- %s (%s): 角色=%s, 启用=%v, 状态=%s, 执行次数=%d\n",
						name, agentID, role, enabled, state, executionCount)
				}
			}
		}
	}

	if agentStatus == "" {
		agentStatus = "无Agent状态信息"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"task_list":    taskList,
		"agent_status": agentStatus,
	}

	prompt, err := a.FormatPrompt("agent_scheduling", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("请为以下任务分配合适的Agent：\n任务列表：%s\nAgent状态：%s", taskList, agentStatus)
	}

	// 调用模型进行Agent分配
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("Agent调度失败: %v", err),
		}, nil
	}

	schedulingPlan := modelResponse.Content

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"agent_scheduled":     true,
			"task_list":           taskList,
			"scheduling_plan":     schedulingPlan,
			"agent_status":        agentStatus,
			"model_response":      modelResponse,
			"scheduling_criteria": "能力匹配、负载均衡、优先级、依赖关系、并发限制",
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "agent_scheduling",
				Name:        "Agent调度方案",
				Content:     schedulingPlan,
				Description: "任务的Agent分配方案",
			},
		},
		Logs: []string{
			"Agent调度完成",
			fmt.Sprintf("任务数量: %d", len(strings.Split(taskList, "\n"))),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// handleProgressMonitoring 处理进度监控
func (a *OrchestratorAgent) handleProgressMonitoring(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理进度监控: %s", task)

	// 获取参数
	projectName, _ := parameters["project_name"].(string)
	executionPlan, _ := parameters["execution_plan"].(string)
	currentStatus, _ := parameters["current_status"].(string)
	completedTasks, _ := parameters["completed_tasks"].(string)
	inProgressTasks, _ := parameters["in_progress_tasks"].(string)
	pendingTasks, _ := parameters["pending_tasks"].(string)

	if projectName == "" {
		projectName = "当前项目"
	}

	if executionPlan == "" {
		return &AgentResult{
			Success: false,
			Output:  "execution_plan is required for progress monitoring",
		}, nil
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"project_name":      projectName,
		"execution_plan":    executionPlan,
		"current_status":    currentStatus,
		"completed_tasks":   completedTasks,
		"in_progress_tasks": inProgressTasks,
		"pending_tasks":     pendingTasks,
	}

	prompt, err := a.FormatPrompt("progress_monitoring", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("请监控以下项目进度：\n项目名称：%s\n执行计划：%s\n当前状态：%s\n已完成：%s\n进行中：%s\n待完成：%s",
			projectName, executionPlan, currentStatus, completedTasks, inProgressTasks, pendingTasks)
	}

	// 调用模型进行进度分析
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("进度监控失败: %v", err),
		}, nil
	}

	progressReport := modelResponse.Content

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"progress_monitored": true,
			"project_name":       projectName,
			"progress_report":    progressReport,
			"execution_plan":     executionPlan,
			"current_status":     currentStatus,
			"model_response":     modelResponse,
			"monitoring_aspects": "进度评估、问题识别、调整建议、资源优化、风险更新",
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "progress_report",
				Name:        "进度监控报告",
				Content:     progressReport,
				Description: "项目进度监控和分析报告",
			},
		},
		Logs: []string{
			"进度监控完成",
			fmt.Sprintf("项目: %s", projectName),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// handleCoordination 处理协调任务
func (a *OrchestratorAgent) handleCoordination(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理协调任务: %s", task)

	// 获取参数
	scenario, _ := parameters["scenario"].(string)
	involvedAgents, _ := parameters["involved_agents"].(string)
	goal, _ := parameters["goal"].(string)
	currentIssue, _ := parameters["current_issue"].(string)

	if scenario == "" {
		scenario = task
	}

	if involvedAgents == "" {
		return &AgentResult{
			Success: false,
			Output:  "involved_agents is required",
		}, nil
	}

	if goal == "" {
		goal = "完成协作任务"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"scenario":        scenario,
		"involved_agents": involvedAgents,
		"goal":            goal,
		"current_issue":   currentIssue,
	}

	prompt, err := a.FormatPrompt("coordination", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("请协调以下Agent间的协作：\n场景：%s\n涉及Agent：%s\n目标：%s\n当前问题：%s",
			scenario, involvedAgents, goal, currentIssue)
	}

	// 调用模型生成协调方案
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("协调方案生成失败: %v", err),
		}, nil
	}

	coordinationPlan := modelResponse.Content

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"coordination_planned": true,
			"scenario":             scenario,
			"coordination_plan":    coordinationPlan,
			"involved_agents":      involvedAgents,
			"goal":                 goal,
			"model_response":       modelResponse,
			"coordination_aspects": "沟通机制、接口定义、同步点、冲突解决、质量保证",
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "coordination_plan",
				Name:        "协调方案",
				Content:     coordinationPlan,
				Description: "Agent间协作的协调方案",
			},
		},
		Logs: []string{
			"协调方案生成完成",
			fmt.Sprintf("场景: %s", scenario),
			fmt.Sprintf("涉及Agent: %s", involvedAgents),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// handleGenericOrchestration 处理通用编排任务
func (a *OrchestratorAgent) handleGenericOrchestration(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理通用编排任务: %s", task)

	// 使用模型分析任务
	promptData := map[string]interface{}{
		"task":       task,
		"parameters": parameters,
	}

	prompt, err := a.FormatPrompt("task_decomposition", promptData)
	if err != nil {
		prompt = fmt.Sprintf("作为项目编排器，请分析并处理以下任务：%s\n参数：%v", task, parameters)
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
				Type:        "orchestration_analysis",
				Name:        "编排分析",
				Content:     analysis,
				Description: "编排器对任务的分析和建议",
			},
		},
		Logs: []string{
			"通用编排任务完成",
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// ValidateTask 验证任务（编排器Agent特定验证）
func (a *OrchestratorAgent) ValidateTask(task string) error {
	// 调用基础验证
	if err := a.BaseAgent.ValidateTask(task); err != nil {
		return err
	}

	// 编排器Agent特定的验证
	taskLower := strings.ToLower(task)

	// 检查是否包含有效的编排关键词
	orchestrationKeywords := []string{
		"分解", "计划", "调度", "分配", "协调", "监控", "进度",
		"decompose", "plan", "schedule", "assign", "coordinate", "monitor", "progress",
	}

	hasKeyword := false
	for _, keyword := range orchestrationKeywords {
		if strings.Contains(taskLower, keyword) {
			hasKeyword = true
			break
		}
	}

	if !hasKeyword {
		utils.Warn("编排器任务可能不相关: %s", task)
		// 不返回错误，只是警告
	}

	return nil
}

// GetOrchestrationCapabilities 获取编排能力描述
func (a *OrchestratorAgent) GetOrchestrationCapabilities() string {
	capabilities := []string{
		"任务分解：将复杂需求分解为具体可执行任务",
		"执行计划：制定详细的项目执行计划和时间表",
		"Agent调度：根据任务特性和Agent能力进行智能分配",
		"进度监控：实时跟踪项目进度并识别问题",
		"协调协作：协调多个Agent间的协作和接口定义",
		"风险管理：识别和应对项目风险",
		"资源优化：优化Agent和资源的使用效率",
	}

	return strings.Join(capabilities, "\n- ")
}
