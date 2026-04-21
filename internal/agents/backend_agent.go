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

// BackendDeveloperAgent 后端开发工程师Agent
type BackendDeveloperAgent struct {
	*BaseAgent
}

// NewBackendDeveloperAgent 创建后端开发工程师Agent
func NewBackendDeveloperAgent(config AgentConfig, modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine) (*BackendDeveloperAgent, error) {
	// 设置后端开发工程师Agent的默认配置
	config.Role = RoleBackendDev
	if config.Description == "" {
		config.Description = "后端开发工程师，负责后端API、数据库代码生成和业务逻辑实现"
	}

	// 设置默认能力
	if len(config.Capabilities) == 0 {
		config.Capabilities = []AgentCapability{
			CapabilityCodeGeneration,
			CapabilityFileOperation,
			CapabilityDatabaseDesign,
			CapabilityAPIDesign,
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
		"api_development": `你是一个后端开发工程师。请根据以下需求开发API：

项目类型：{{project_type}}
技术栈：{{tech_stack}}
API名称：{{api_name}}
API功能：{{api_functionality}}
请求方法：{{http_method}}
请求参数：{{request_parameters}}
响应格式：{{response_format}}
数据库表：{{database_tables}}
业务逻辑：{{business_logic}}

请生成以下内容：
1. API路由定义
2. 控制器代码
3. 服务层代码
4. 数据访问层代码
5. 模型定义
6. 请求验证
7. 响应处理
8. 错误处理
9. 日志记录
10. 单元测试

请输出完整的API代码：`,

		"database_design": `你是一个后端开发工程师。请设计数据库：

项目名称：{{project_name}}
业务需求：{{business_requirements}}
数据模型：{{data_models}}
查询需求：{{query_requirements}}
性能要求：{{performance_requirements}}

请提供以下设计内容：
1. 数据库选型（MySQL/PostgreSQL/MongoDB等）
2. 表结构设计（表名、字段、类型、约束）
3. 索引设计（主键、外键、索引）
4. 关系设计（一对一、一对多、多对多）
5. 视图设计（如果需要）
6. 存储过程设计（如果需要）
7. 触发器设计（如果需要）
8. 数据迁移脚本
9. 数据库初始化脚本
10. 性能优化建议

请输出完整的数据库设计方案：`,

		"backend_project_setup": `你是一个后端开发工程师。请设置后端项目：

项目名称：{{project_name}}
项目描述：{{project_description}}
技术栈：{{tech_stack}}
框架选择：{{framework_choice}}
数据库选择：{{database_choice}}
API风格：{{api_style}}
认证方式：{{authentication_method}}

请提供：
1. 项目目录结构
2. 配置文件（go.mod, package.json, pom.xml等）
3. 依赖配置
4. 数据库配置
5. API路由配置
6. 中间件配置
7. 错误处理配置
8. 日志配置
9. 环境配置
10. 部署配置
11. 测试配置

请输出完整的项目设置方案：`,

		"business_logic": `你是一个后端开发工程师。请实现以下业务逻辑：

业务场景：{{business_scenario}}
输入数据：{{input_data}}
处理规则：{{processing_rules}}
输出结果：{{output_requirements}}
错误情况：{{error_cases}}

请生成以下内容：
1. 业务逻辑代码
2. 数据验证逻辑
3. 数据处理逻辑
4. 计算逻辑
5. 状态管理
6. 事务处理
7. 并发控制
8. 缓存策略
9. 错误处理
10. 日志记录

请输出完整的业务逻辑代码：`,
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

	agent := &BackendDeveloperAgent{
		BaseAgent: baseAgent,
	}

	return agent, nil
}

// Execute 执行后端开发工程师Agent任务
func (a *BackendDeveloperAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
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

	utils.Info("后端开发工程师Agent %s 开始执行任务: %s", a.config.Name, task)

	startTime := time.Now()

	// 根据任务类型选择提示词模板
	templateName := "api_development"
	taskLower := strings.ToLower(task)

	if strings.Contains(taskLower, "数据库") || strings.Contains(taskLower, "表") || strings.Contains(taskLower, "data") {
		templateName = "database_design"
	} else if strings.Contains(taskLower, "项目") || strings.Contains(taskLower, "设置") || strings.Contains(taskLower, "初始化") {
		templateName = "backend_project_setup"
	} else if strings.Contains(taskLower, "业务") || strings.Contains(taskLower, "逻辑") || strings.Contains(taskLower, "business") {
		templateName = "business_logic"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"task":       task,
		"parameters": parameters,
	}

	// 添加特定参数
	if projectType, ok := parameters["project_type"]; ok {
		promptData["project_type"] = projectType
	}
	if techStack, ok := parameters["tech_stack"]; ok {
		promptData["tech_stack"] = techStack
	}
	if apiName, ok := parameters["api_name"]; ok {
		promptData["api_name"] = apiName
	}
	if apiFunctionality, ok := parameters["api_functionality"]; ok {
		promptData["api_functionality"] = apiFunctionality
	}

	prompt, err := a.FormatPrompt(templateName, promptData)
	if err != nil {
		prompt = fmt.Sprintf("你是一个后端开发工程师。请根据以下需求开发后端代码：\n\n需求：%s\n\n参数：%v\n\n请输出完整的后端代码。", task, parameters)
	}

	// 调用模型生成代码
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		a.execution.State = StateFailed
		a.execution.Error = fmt.Sprintf("后端开发失败: %v", err)
		a.state = StateFailed
		a.execution.CompletedAt = &time.Time{}
		*a.execution.CompletedAt = time.Now()

		return &AgentResult{
			Success:  false,
			Output:   fmt.Sprintf("后端开发失败: %v", err),
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

	backendCode := modelResponse.Content

	// 提取代码文件列表
	codeFiles := a.extractCodeFiles(backendCode)

	// 提取API端点列表
	apiEndpoints := a.extractAPIEndpoints(backendCode)

	result := &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"task_analyzed":  true,
			"backend_code":   backendCode,
			"code_files":     codeFiles,
			"api_endpoints":  apiEndpoints,
			"model_response": modelResponse,
			"agent_role":     a.config.Role,
			"agent_name":     a.config.Name,
			"template_used":  templateName,
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "backend_code",
				Name:        "后端代码",
				Content:     backendCode,
				Description: fmt.Sprintf("%s生成的后端代码", a.config.Name),
			},
			{
				Type:        "code_files",
				Name:        "代码文件列表",
				Content:     strings.Join(codeFiles, "\n"),
				Description: "生成的后端代码文件列表",
			},
			{
				Type:        "api_endpoints",
				Name:        "API端点列表",
				Content:     strings.Join(apiEndpoints, "\n"),
				Description: "生成的API端点列表",
			},
		},
		Duration: duration,
		Logs: []string{
			fmt.Sprintf("后端开发完成"),
			fmt.Sprintf("使用模板: %s", templateName),
			fmt.Sprintf("生成代码文件: %d个", len(codeFiles)),
			fmt.Sprintf("生成API端点: %d个", len(apiEndpoints)),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
			fmt.Sprintf("执行时间: %v", duration),
		},
	}

	a.execution.Result = result

	return result, nil
}

// extractCodeFiles 从代码内容中提取文件列表
func (a *BackendDeveloperAgent) extractCodeFiles(codeContent string) []string {
	var files []string

	lines := strings.Split(codeContent, "\n")
	currentFile := ""
	inCodeBlock := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 检查代码块开始
		if strings.HasPrefix(trimmedLine, "```") && !inCodeBlock {
			inCodeBlock = true
			continue
		}

		// 检查代码块结束
		if strings.HasPrefix(trimmedLine, "```") && inCodeBlock {
			inCodeBlock = false
			if currentFile != "" {
				files = append(files, currentFile)
				currentFile = ""
			}
			continue
		}

		// 在代码块中，检查文件名注释
		if inCodeBlock {
			// 检查常见的文件名模式
			if strings.Contains(line, "// File:") || strings.Contains(line, "// file:") ||
				strings.Contains(line, "# File:") || strings.Contains(line, "# file:") ||
				strings.Contains(line, "/* File:") || strings.Contains(line, "/* file:") ||
				strings.Contains(line, "-- File:") || strings.Contains(line, "-- file:") {
				// 提取文件名
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					fileName := strings.TrimSpace(parts[1])
					fileName = strings.TrimSuffix(fileName, "*/")
					fileName = strings.TrimSpace(fileName)

					if currentFile != "" {
						files = append(files, currentFile)
					}
					currentFile = fileName
				}
			}
		}

		// 检查文件扩展名模式
		if !inCodeBlock && (strings.Contains(line, ".go") || strings.Contains(line, ".java") ||
			strings.Contains(line, ".py") || strings.Contains(line, ".js") ||
			strings.Contains(line, ".ts") || strings.Contains(line, ".sql") ||
			strings.Contains(line, ".yaml") || strings.Contains(line, ".yml") ||
			strings.Contains(line, ".json") || strings.Contains(line, ".toml")) {
			// 尝试提取文件名
			words := strings.Fields(line)
			for _, word := range words {
				if strings.Contains(word, ".") && len(word) > 3 {
					// 简单的文件名检测
					files = append(files, word)
				}
			}
		}
	}

	// 添加最后一个文件
	if currentFile != "" {
		files = append(files, currentFile)
	}

	// 去重
	uniqueFiles := make(map[string]bool)
	var result []string
	for _, file := range files {
		if !uniqueFiles[file] {
			uniqueFiles[file] = true
			result = append(result, file)
		}
	}

	return result
}

// extractAPIEndpoints 从代码内容中提取API端点
func (a *BackendDeveloperAgent) extractAPIEndpoints(codeContent string) []string {
	var endpoints []string

	lines := strings.Split(codeContent, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 检查常见的API路由模式
		if strings.Contains(trimmedLine, "@RequestMapping") ||
			strings.Contains(trimmedLine, "@GetMapping") ||
			strings.Contains(trimmedLine, "@PostMapping") ||
			strings.Contains(trimmedLine, "@PutMapping") ||
			strings.Contains(trimmedLine, "@DeleteMapping") ||
			strings.Contains(trimmedLine, "@PatchMapping") ||
			strings.Contains(trimmedLine, "router.get(") ||
			strings.Contains(trimmedLine, "router.post(") ||
			strings.Contains(trimmedLine, "router.put(") ||
			strings.Contains(trimmedLine, "router.delete(") ||
			strings.Contains(trimmedLine, "app.get(") ||
			strings.Contains(trimmedLine, "app.post(") ||
			strings.Contains(trimmedLine, "app.put(") ||
			strings.Contains(trimmedLine, "app.delete(") ||
			strings.Contains(trimmedLine, "http.HandleFunc(") ||
			strings.Contains(trimmedLine, "r.GET(") ||
			strings.Contains(trimmedLine, "r.POST(") ||
			strings.Contains(trimmedLine, "r.PUT(") ||
			strings.Contains(trimmedLine, "r.DELETE(") {

			// 提取路径
			parts := strings.Split(trimmedLine, "\"")
			if len(parts) >= 2 {
				endpoint := parts[1]
				if endpoint != "" && !strings.Contains(endpoint, " ") {
					// 添加HTTP方法
					method := "GET"
					if strings.Contains(trimmedLine, "post") || strings.Contains(trimmedLine, "POST") {
						method = "POST"
					} else if strings.Contains(trimmedLine, "put") || strings.Contains(trimmedLine, "PUT") {
						method = "PUT"
					} else if strings.Contains(trimmedLine, "delete") || strings.Contains(trimmedLine, "DELETE") {
						method = "DELETE"
					} else if strings.Contains(trimmedLine, "patch") || strings.Contains(trimmedLine, "PATCH") {
						method = "PATCH"
					}

					endpoints = append(endpoints, fmt.Sprintf("%s %s", method, endpoint))
				}
			}
		}
	}

	// 去重
	uniqueEndpoints := make(map[string]bool)
	var result []string
	for _, endpoint := range endpoints {
		if !uniqueEndpoints[endpoint] {
			uniqueEndpoints[endpoint] = true
			result = append(result, endpoint)
		}
	}

	return result
}

// ValidateTask 验证任务（后端开发工程师特定验证）
func (a *BackendDeveloperAgent) ValidateTask(task string) error {
	// 调用基础验证
	if err := a.BaseAgent.ValidateTask(task); err != nil {
		return err
	}

	// 后端开发工程师特定验证
	taskLower := strings.ToLower(task)

	// 检查是否包含后端开发相关关键词
	backendKeywords := []string{
		"后端", "API", "接口", "数据库", "服务", "业务逻辑", "服务器", "Go", "Java", "Python", "Node.js",
		"backend", "api", "interface", "database", "service", "business logic", "server", "go", "java", "python", "node",
		"路由", "控制器", "模型", "中间件", "认证", "授权", "route", "controller", "model", "middleware", "auth", "authentication",
	}

	hasBackendKeyword := false
	for _, keyword := range backendKeywords {
		if strings.Contains(taskLower, keyword) {
			hasBackendKeyword = true
			break
		}
	}

	if !hasBackendKeyword {
		utils.Warn("任务可能不适用于后端开发工程师: %s", task)
	}

	return nil
}

// GetRoleDescription 获取角色描述
func (a *BackendDeveloperAgent) GetRoleDescription() string {
	return "后端开发工程师：负责后端API、数据库代码生成和业务逻辑实现"
}

// DevelopAPI 开发API（专用方法）
func (a *BackendDeveloperAgent) DevelopAPI(ctx context.Context, apiName, functionality, httpMethod string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"api_name":          apiName,
		"api_functionality": functionality,
		"http_method":       httpMethod,
	}

	task := fmt.Sprintf("开发API'%s'", apiName)
	return a.Execute(ctx, task, parameters)
}

// DesignDatabase 设计数据库（专用方法）
func (a *BackendDeveloperAgent) DesignDatabase(ctx context.Context, projectName, businessRequirements string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"project_name":          projectName,
		"business_requirements": businessRequirements,
	}

	task := fmt.Sprintf("为项目'%s'设计数据库", projectName)
	return a.Execute(ctx, task, parameters)
}

// SetupBackendProject 设置后端项目（专用方法）
func (a *BackendDeveloperAgent) SetupBackendProject(ctx context.Context, projectName, projectDescription, framework string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"project_name":        projectName,
		"project_description": projectDescription,
		"framework_choice":    framework,
	}

	task := fmt.Sprintf("设置后端项目'%s'", projectName)
	return a.Execute(ctx, task, parameters)
}

// ImplementBusinessLogic 实现业务逻辑（专用方法）
func (a *BackendDeveloperAgent) ImplementBusinessLogic(ctx context.Context, businessScenario, processingRules string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"business_scenario": businessScenario,
		"processing_rules":  processingRules,
	}

	task := fmt.Sprintf("实现业务逻辑'%s'", businessScenario)
	return a.Execute(ctx, task, parameters)
}
