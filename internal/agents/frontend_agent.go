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

// FrontendDeveloperAgent 前端开发工程师Agent
type FrontendDeveloperAgent struct {
	*BaseAgent
}

// NewFrontendDeveloperAgent 创建前端开发工程师Agent
func NewFrontendDeveloperAgent(config AgentConfig, modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine) (*FrontendDeveloperAgent, error) {
	// 设置前端开发工程师Agent的默认配置
	config.Role = RoleFrontendDev
	if config.Description == "" {
		config.Description = "前端开发工程师，负责React/Vue组件开发和前端界面实现"
	}

	// 设置默认能力
	if len(config.Capabilities) == 0 {
		config.Capabilities = []AgentCapability{
			CapabilityCodeGeneration,
			CapabilityFileOperation,
			CapabilityUIUXDesign,
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
		"component_development": `你是一个前端开发工程师。请根据以下需求开发前端组件：

项目类型：{{project_type}}
技术栈：{{tech_stack}}
组件名称：{{component_name}}
组件功能：{{component_functionality}}
设计要求：{{design_requirements}}
API接口：{{api_interfaces}}

请生成以下内容：
1. 组件代码（React/Vue组件）
2. 样式代码（CSS/SCSS/Tailwind）
3. 状态管理（如果需要）
4. 事件处理
5. 数据获取和更新逻辑
6. 错误处理
7. 加载状态
8. 响应式设计

请输出完整的组件代码：`,

		"ui_implementation": `你是一个前端开发工程师。请实现以下UI界面：

页面名称：{{page_name}}
页面功能：{{page_function}}
设计稿：{{design_reference}}
交互要求：{{interaction_requirements}}
响应式要求：{{responsive_requirements}}

请生成以下内容：
1. 页面结构（HTML/JSX）
2. 样式实现（CSS/SCSS/Tailwind）
3. 交互逻辑（JavaScript/TypeScript）
4. 状态管理
5. 路由配置（如果需要）
6. API集成
7. 表单验证
8. 错误处理
9. 加载优化
10. 性能优化

请输出完整的页面实现代码：`,

		"frontend_project_setup": `你是一个前端开发工程师。请设置前端项目：

项目名称：{{project_name}}
项目描述：{{project_description}}
技术栈：{{tech_stack}}
框架选择：{{framework_choice}}
UI库选择：{{ui_library}}
状态管理：{{state_management}}
构建工具：{{build_tool}}

请提供：
1. 项目目录结构
2. 配置文件（package.json, vite.config.js, webpack.config.js等）
3. 基础组件
4. 路由配置
5. 状态管理配置
6. API服务配置
7. 样式配置
8. 开发环境配置
9. 构建和部署配置
10. 测试配置

请输出完整的项目设置方案：`,
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

	agent := &FrontendDeveloperAgent{
		BaseAgent: baseAgent,
	}

	return agent, nil
}

// Execute 执行前端开发工程师Agent任务
func (a *FrontendDeveloperAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
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

	utils.Info("前端开发工程师Agent %s 开始执行任务: %s", a.config.Name, task)

	startTime := time.Now()

	// 根据任务类型选择提示词模板
	templateName := "component_development"
	taskLower := strings.ToLower(task)

	if strings.Contains(taskLower, "页面") || strings.Contains(taskLower, "ui") || strings.Contains(taskLower, "界面") {
		templateName = "ui_implementation"
	} else if strings.Contains(taskLower, "项目") || strings.Contains(taskLower, "设置") || strings.Contains(taskLower, "初始化") {
		templateName = "frontend_project_setup"
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
	if componentName, ok := parameters["component_name"]; ok {
		promptData["component_name"] = componentName
	}

	prompt, err := a.FormatPrompt(templateName, promptData)
	if err != nil {
		prompt = fmt.Sprintf("你是一个前端开发工程师。请根据以下需求开发前端代码：\n\n需求：%s\n\n参数：%v\n\n请输出完整的前端代码。", task, parameters)
	}

	// 调用模型生成代码
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		a.execution.State = StateFailed
		a.execution.Error = fmt.Sprintf("前端开发失败: %v", err)
		a.state = StateFailed
		a.execution.CompletedAt = &time.Time{}
		*a.execution.CompletedAt = time.Now()

		return &AgentResult{
			Success:  false,
			Output:   fmt.Sprintf("前端开发失败: %v", err),
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

	frontendCode := modelResponse.Content

	// 提取代码文件列表
	codeFiles := a.extractCodeFiles(frontendCode)

	result := &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"task_analyzed":  true,
			"frontend_code":  frontendCode,
			"code_files":     codeFiles,
			"model_response": modelResponse,
			"agent_role":     a.config.Role,
			"agent_name":     a.config.Name,
			"template_used":  templateName,
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "frontend_code",
				Name:        "前端代码",
				Content:     frontendCode,
				Description: fmt.Sprintf("%s生成的前端代码", a.config.Name),
			},
			{
				Type:        "code_files",
				Name:        "代码文件列表",
				Content:     strings.Join(codeFiles, "\n"),
				Description: "生成的前端代码文件列表",
			},
		},
		Duration: duration,
		Logs: []string{
			fmt.Sprintf("前端开发完成"),
			fmt.Sprintf("使用模板: %s", templateName),
			fmt.Sprintf("生成代码文件: %d个", len(codeFiles)),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
			fmt.Sprintf("执行时间: %v", duration),
		},
	}

	a.execution.Result = result

	return result, nil
}

// extractCodeFiles 从代码内容中提取文件列表
func (a *FrontendDeveloperAgent) extractCodeFiles(codeContent string) []string {
	var files []string

	lines := strings.Split(codeContent, "\n")
	currentFile := ""
	inCodeBlock := false
	codeBlockLanguage := ""

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 检查代码块开始
		if strings.HasPrefix(trimmedLine, "```") && !inCodeBlock {
			inCodeBlock = true
			codeBlockLanguage = strings.TrimPrefix(trimmedLine, "```")
			if codeBlockLanguage == "" {
				codeBlockLanguage = "unknown"
			}
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
				strings.Contains(line, "<!-- File:") || strings.Contains(line, "<!-- file:") ||
				strings.Contains(line, "/* File:") || strings.Contains(line, "/* file:") {
				// 提取文件名
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					fileName := strings.TrimSpace(parts[1])
					fileName = strings.TrimSuffix(fileName, "*/")
					fileName = strings.TrimSuffix(fileName, "-->")
					fileName = strings.TrimSpace(fileName)

					if currentFile != "" {
						files = append(files, currentFile)
					}
					currentFile = fileName
				}
			}
		}

		// 检查文件扩展名模式
		if !inCodeBlock && (strings.Contains(line, ".js") || strings.Contains(line, ".ts") ||
			strings.Contains(line, ".jsx") || strings.Contains(line, ".tsx") ||
			strings.Contains(line, ".vue") || strings.Contains(line, ".css") ||
			strings.Contains(line, ".scss") || strings.Contains(line, ".html")) {
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

// ValidateTask 验证任务（前端开发工程师特定验证）
func (a *FrontendDeveloperAgent) ValidateTask(task string) error {
	// 调用基础验证
	if err := a.BaseAgent.ValidateTask(task); err != nil {
		return err
	}

	// 前端开发工程师特定验证
	taskLower := strings.ToLower(task)

	// 检查是否包含前端开发相关关键词
	frontendKeywords := []string{
		"前端", "组件", "页面", "界面", "UI", "React", "Vue", "Angular", "JavaScript", "TypeScript",
		"frontend", "component", "page", "ui", "react", "vue", "angular", "javascript", "typescript",
		"样式", "CSS", "HTML", "响应式", "style", "css", "html", "responsive",
	}

	hasFrontendKeyword := false
	for _, keyword := range frontendKeywords {
		if strings.Contains(taskLower, keyword) {
			hasFrontendKeyword = true
			break
		}
	}

	if !hasFrontendKeyword {
		utils.Warn("任务可能不适用于前端开发工程师: %s", task)
	}

	return nil
}

// GetRoleDescription 获取角色描述
func (a *FrontendDeveloperAgent) GetRoleDescription() string {
	return "前端开发工程师：负责React/Vue组件开发和前端界面实现"
}

// DevelopComponent 开发前端组件（专用方法）
func (a *FrontendDeveloperAgent) DevelopComponent(ctx context.Context, componentName, functionality, techStack string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"component_name":          componentName,
		"component_functionality": functionality,
		"tech_stack":              techStack,
	}

	task := fmt.Sprintf("开发前端组件'%s'", componentName)
	return a.Execute(ctx, task, parameters)
}

// ImplementUIPage 实现UI页面（专用方法）
func (a *FrontendDeveloperAgent) ImplementUIPage(ctx context.Context, pageName, pageFunction, designRequirements string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"page_name":           pageName,
		"page_function":       pageFunction,
		"design_requirements": designRequirements,
	}

	task := fmt.Sprintf("实现UI页面'%s'", pageName)
	return a.Execute(ctx, task, parameters)
}

// SetupFrontendProject 设置前端项目（专用方法）
func (a *FrontendDeveloperAgent) SetupFrontendProject(ctx context.Context, projectName, projectDescription, framework string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"project_name":        projectName,
		"project_description": projectDescription,
		"framework_choice":    framework,
	}

	task := fmt.Sprintf("设置前端项目'%s'", projectName)
	return a.Execute(ctx, task, parameters)
}
