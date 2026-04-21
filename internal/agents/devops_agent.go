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

// DevOpsEngineerAgent DevOps工程师Agent
type DevOpsEngineerAgent struct {
	*BaseAgent
}

// NewDevOpsEngineerAgent 创建DevOps工程师Agent
func NewDevOpsEngineerAgent(config AgentConfig, modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine) (*DevOpsEngineerAgent, error) {
	// 设置DevOps工程师Agent的默认配置
	config.Role = RoleDevOps
	if config.Description == "" {
		config.Description = "DevOps工程师，负责项目构建、启动、部署验证和运维"
	}

	// 设置默认能力
	if len(config.Capabilities) == 0 {
		config.Capabilities = []AgentCapability{
			CapabilityShellExecution,
			CapabilityProjectSetup,
			CapabilityDeployment,
			CapabilityMonitoring,
		}
	}

	// 设置默认工具
	if len(config.Tools) == 0 {
		config.Tools = []string{
			"shell_tool",
			"project_tool",
			"fs_tool",
		}
	}

	// 设置默认提示词模板
	if config.PromptTemplates == nil {
		config.PromptTemplates = make(map[string]string)
	}

	// 添加默认提示词模板
	defaultTemplates := map[string]string{
		"project_build": `你是一个DevOps工程师。请构建以下项目：

项目名称：{{project_name}}
项目类型：{{project_type}}
技术栈：{{tech_stack}}
构建工具：{{build_tool}}
构建环境：{{build_environment}}

请执行以下构建步骤：
1. 环境检查
2. 依赖安装
3. 代码编译
4. 测试执行
5. 打包
6. 构建验证
7. 构建报告

请输出详细的构建过程和结果：`,

		"project_deployment": `你是一个DevOps工程师。请部署以下项目：

项目名称：{{project_name}}
部署环境：{{deployment_environment}}
部署目标：{{deployment_target}}
部署策略：{{deployment_strategy}}
配置管理：{{configuration_management}}

请执行以下部署步骤：
1. 环境准备
2. 配置部署
3. 服务部署
4. 数据库部署
5. 网络配置
6. 安全配置
7. 部署验证
8. 回滚方案

请输出详细的部署过程和结果：`,

		"infrastructure_setup": `你是一个DevOps工程师。请设置以下基础设施：

基础设施类型：{{infrastructure_type}}
云提供商：{{cloud_provider}}
资源需求：{{resource_requirements}}
网络需求：{{network_requirements}}
安全需求：{{security_requirements}}

请提供以下基础设施设置：
1. 资源规划
2. 网络架构
3. 安全策略
4. 监控配置
5. 备份策略
6. 高可用配置
7. 成本优化
8. 运维手册

请输出详细的基础设施设置方案：`,

		"monitoring_setup": `你是一个DevOps工程师。请设置以下监控系统：

监控目标：{{monitoring_target}}
监控指标：{{monitoring_metrics}}
告警规则：{{alert_rules}}
可视化需求：{{visualization_requirements}}
存储需求：{{storage_requirements}}

请设置以下监控组件：
1. 指标收集
2. 日志收集
3. 追踪系统
4. 告警系统
5. 仪表板
6. 报表系统
7. 容量规划
8. 性能分析

请输出详细的监控系统设置方案：`,
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

	agent := &DevOpsEngineerAgent{
		BaseAgent: baseAgent,
	}

	return agent, nil
}

// Execute 执行DevOps工程师Agent任务
func (a *DevOpsEngineerAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
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

	utils.Info("DevOps工程师Agent %s 开始执行任务: %s", a.config.Name, task)

	startTime := time.Now()

	// 根据任务类型选择提示词模板
	templateName := "project_build"
	taskLower := strings.ToLower(task)

	if strings.Contains(taskLower, "部署") || strings.Contains(taskLower, "deploy") || strings.Contains(taskLower, "deployment") {
		templateName = "project_deployment"
	} else if strings.Contains(taskLower, "基础设施") || strings.Contains(taskLower, "infrastructure") || strings.Contains(taskLower, "环境") {
		templateName = "infrastructure_setup"
	} else if strings.Contains(taskLower, "监控") || strings.Contains(taskLower, "monitor") || strings.Contains(taskLower, "告警") {
		templateName = "monitoring_setup"
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
	if projectType, ok := parameters["project_type"]; ok {
		promptData["project_type"] = projectType
	}
	if techStack, ok := parameters["tech_stack"]; ok {
		promptData["tech_stack"] = techStack
	}

	prompt, err := a.FormatPrompt(templateName, promptData)
	if err != nil {
		prompt = fmt.Sprintf("你是一个DevOps工程师。请根据以下需求进行项目构建和部署：\n\n需求：%s\n\n参数：%v\n\n请输出详细的构建和部署方案。", task, parameters)
	}

	// 调用模型生成方案
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		a.execution.State = StateFailed
		a.execution.Error = fmt.Sprintf("DevOps任务执行失败: %v", err)
		a.state = StateFailed
		a.execution.CompletedAt = &time.Time{}
		*a.execution.CompletedAt = time.Now()

		return &AgentResult{
			Success:  false,
			Output:   fmt.Sprintf("DevOps任务执行失败: %v", err),
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

	devopsPlan := modelResponse.Content

	// 提取命令列表
	commands := a.extractCommands(devopsPlan)

	result := &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"task_analyzed":  true,
			"devops_plan":    devopsPlan,
			"commands":       commands,
			"model_response": modelResponse,
			"agent_role":     a.config.Role,
			"agent_name":     a.config.Name,
			"template_used":  templateName,
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "devops_plan",
				Name:        "DevOps方案",
				Content:     devopsPlan,
				Description: fmt.Sprintf("%s生成的DevOps方案", a.config.Name),
			},
			{
				Type:        "commands",
				Name:        "命令列表",
				Content:     strings.Join(commands, "\n"),
				Description: "生成的命令列表",
			},
		},
		Duration: duration,
		Logs: []string{
			fmt.Sprintf("DevOps任务完成"),
			fmt.Sprintf("使用模板: %s", templateName),
			fmt.Sprintf("生成命令: %d个", len(commands)),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
			fmt.Sprintf("执行时间: %v", duration),
		},
	}

	a.execution.Result = result

	return result, nil
}

// extractCommands 从DevOps方案中提取命令
func (a *DevOpsEngineerAgent) extractCommands(devopsPlan string) []string {
	var commands []string

	lines := strings.Split(devopsPlan, "\n")
	inCommandBlock := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 检查命令块开始
		if strings.HasPrefix(trimmedLine, "```bash") || strings.HasPrefix(trimmedLine, "```sh") ||
			strings.HasPrefix(trimmedLine, "```shell") {
			inCommandBlock = true
			continue
		}

		// 检查命令块结束
		if strings.HasPrefix(trimmedLine, "```") && inCommandBlock {
			inCommandBlock = false
			continue
		}

		// 在命令块中，提取命令
		if inCommandBlock {
			if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "#") {
				commands = append(commands, trimmedLine)
			}
		}

		// 检查常见的命令模式
		if !inCommandBlock && (strings.HasPrefix(trimmedLine, "$ ") ||
			strings.HasPrefix(trimmedLine, "> ") ||
			strings.HasPrefix(trimmedLine, "# ") ||
			strings.Contains(trimmedLine, "docker ") ||
			strings.Contains(trimmedLine, "kubectl ") ||
			strings.Contains(trimmedLine, "terraform ") ||
			strings.Contains(trimmedLine, "ansible ") ||
			strings.Contains(trimmedLine, "npm ") ||
			strings.Contains(trimmedLine, "yarn ") ||
			strings.Contains(trimmedLine, "go ") ||
			strings.Contains(trimmedLine, "python ") ||
			strings.Contains(trimmedLine, "java ") ||
			strings.Contains(trimmedLine, "mvn ") ||
			strings.Contains(trimmedLine, "gradle ")) {

			// 清理命令
			cleanCommand := strings.TrimPrefix(trimmedLine, "$ ")
			cleanCommand = strings.TrimPrefix(cleanCommand, "> ")
			cleanCommand = strings.TrimPrefix(cleanCommand, "# ")

			if cleanCommand != "" && len(cleanCommand) > 3 {
				commands = append(commands, cleanCommand)
			}
		}
	}

	// 去重
	uniqueCommands := make(map[string]bool)
	var result []string
	for _, cmd := range commands {
		if !uniqueCommands[cmd] {
			uniqueCommands[cmd] = true
			result = append(result, cmd)
		}
	}

	return result
}

// ValidateTask 验证任务（DevOps工程师特定验证）
func (a *DevOpsEngineerAgent) ValidateTask(task string) error {
	// 调用基础验证
	if err := a.BaseAgent.ValidateTask(task); err != nil {
		return err
	}

	// DevOps工程师特定验证
	taskLower := strings.ToLower(task)

	// 检查是否包含DevOps相关关键词
	devopsKeywords := []string{
		"构建", "部署", "运维", "监控", "基础设施", "环境", "配置", "自动化",
		"build", "deploy", "deployment", "operation", "monitor", "infrastructure", "environment", "configuration", "automation",
		"docker", "kubernetes", "k8s", "terraform", "ansible", "jenkins", "gitlab", "ci/cd",
		"云", "服务器", "网络", "安全", "cloud", "server", "network", "security",
	}

	hasDevOpsKeyword := false
	for _, keyword := range devopsKeywords {
		if strings.Contains(taskLower, keyword) {
			hasDevOpsKeyword = true
			break
		}
	}

	if !hasDevOpsKeyword {
		utils.Warn("任务可能不适用于DevOps工程师: %s", task)
	}

	return nil
}

// GetRoleDescription 获取角色描述
func (a *DevOpsEngineerAgent) GetRoleDescription() string {
	return "DevOps工程师：负责项目构建、启动、部署验证和运维"
}

// BuildProject 构建项目（专用方法）
func (a *DevOpsEngineerAgent) BuildProject(ctx context.Context, projectName, projectType, buildTool string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"project_name": projectName,
		"project_type": projectType,
		"build_tool":   buildTool,
	}

	task := fmt.Sprintf("构建项目'%s'", projectName)
	return a.Execute(ctx, task, parameters)
}

// DeployProject 部署项目（专用方法）
func (a *DevOpsEngineerAgent) DeployProject(ctx context.Context, projectName, deploymentEnvironment string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"project_name":           projectName,
		"deployment_environment": deploymentEnvironment,
	}

	task := fmt.Sprintf("部署项目'%s'到环境'%s'", projectName, deploymentEnvironment)
	return a.Execute(ctx, task, parameters)
}

// SetupInfrastructure 设置基础设施（专用方法）
func (a *DevOpsEngineerAgent) SetupInfrastructure(ctx context.Context, infrastructureType, cloudProvider string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"infrastructure_type": infrastructureType,
		"cloud_provider":      cloudProvider,
	}

	task := fmt.Sprintf("设置%s基础设施", infrastructureType)
	return a.Execute(ctx, task, parameters)
}

// SetupMonitoring 设置监控系统（专用方法）
func (a *DevOpsEngineerAgent) SetupMonitoring(ctx context.Context, monitoringTarget string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"monitoring_target": monitoringTarget,
	}

	task := fmt.Sprintf("为'%s'设置监控系统", monitoringTarget)
	return a.Execute(ctx, task, parameters)
}
