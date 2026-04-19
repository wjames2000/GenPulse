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

// FullstackDeveloperAgent 全栈开发Agent
type FullstackDeveloperAgent struct {
	*BaseAgent
}

// NewFullstackDeveloperAgent 创建全栈开发Agent
func NewFullstackDeveloperAgent(config AgentConfig, modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine) (*FullstackDeveloperAgent, error) {
	// 设置全栈开发Agent的默认配置
	config.Role = RoleFullStackDev
	if config.Description == "" {
		config.Description = "全栈开发工程师，能够处理前后端开发任务"
	}

	// 设置默认能力
	if len(config.Capabilities) == 0 {
		config.Capabilities = []AgentCapability{
			CapabilityCodeGeneration,
			CapabilityFileOperation,
			CapabilityGitOperation,
			CapabilityShellExecution,
			CapabilityProjectSetup,
		}
	}

	// 设置默认工具
	if len(config.Tools) == 0 {
		config.Tools = []string{
			"fs_tool",
			"git_tool",
			"shell_tool",
			"project_tool",
		}
	}

	// 设置默认提示词模板
	if config.PromptTemplates == nil {
		config.PromptTemplates = make(map[string]string)
	}

	// 添加默认提示词模板
	defaultTemplates := map[string]string{
		"task_analysis": `你是一个全栈开发工程师。请分析以下任务需求：

任务：{{task}}

请按照以下步骤进行分析：
1. 理解需求：明确用户想要实现什么功能
2. 技术选型：根据需求选择合适的技术栈（前端、后端、数据库等）
3. 项目结构：设计合理的项目目录结构
4. 实现步骤：列出具体的实现步骤
5. 文件清单：需要创建哪些文件，每个文件的内容概要

请开始你的分析：`,

		"code_generation": `你是一个全栈开发工程师。请根据以下需求生成代码：

项目类型：{{project_type}}
项目名称：{{project_name}}
需求描述：{{requirement}}
技术栈：{{tech_stack}}

具体要求：
1. 生成完整可运行的代码
2. 代码要有良好的结构和注释
3. 遵循最佳实践和编码规范
4. 考虑可扩展性和可维护性

请生成以下文件的代码：{{file_list}}

请开始生成代码：`,

		"project_setup": `你是一个全栈开发工程师。请为以下项目设置初始结构：

项目类型：{{project_type}}
项目名称：{{project_name}}
项目描述：{{project_description}}
技术栈：{{tech_stack}}

请提供：
1. 项目目录结构
2. 必要的配置文件（如package.json, go.mod等）
3. 基础代码文件
4. 依赖安装命令
5. 运行和测试命令

请开始设置项目：`,

		"file_creation": `你是一个全栈开发工程师。请创建以下文件：

文件路径：{{file_path}}
文件类型：{{file_type}}
文件用途：{{file_purpose}}
相关代码：{{related_code}}

请生成完整的文件内容，包括必要的注释和文档。`,

		"error_fixing": `你是一个全栈开发工程师。请修复以下错误：

错误信息：{{error_message}}
相关代码：{{related_code}}
错误上下文：{{error_context}}

请分析错误原因并提供修复方案。`,
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

	agent := &FullstackDeveloperAgent{
		BaseAgent: baseAgent,
	}

	return agent, nil
}

// Execute 执行全栈开发任务
func (a *FullstackDeveloperAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
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

	utils.Info("全栈开发Agent %s 开始执行任务: %s", a.config.Name, task)

	startTime := time.Now()
	var result *AgentResult
	var execErr error

	// 根据任务类型执行不同的处理逻辑
	taskType := a.detectTaskType(task, parameters)

	switch taskType {
	case "project_creation":
		result, execErr = a.handleProjectCreation(ctx, task, parameters)
	case "code_generation":
		result, execErr = a.handleCodeGeneration(ctx, task, parameters)
	case "file_operation":
		result, execErr = a.handleFileOperation(ctx, task, parameters)
	case "command_execution":
		result, execErr = a.handleCommandExecution(ctx, task, parameters)
	case "error_fixing":
		result, execErr = a.handleErrorFixing(ctx, task, parameters)
	default:
		result, execErr = a.handleGenericTask(ctx, task, parameters)
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
func (a *FullstackDeveloperAgent) detectTaskType(task string, parameters map[string]interface{}) string {
	taskLower := strings.ToLower(task)

	// 检查关键词
	keywords := map[string]string{
		"project_creation":  "创建项目|新建项目|初始化项目|setup project|init project|create project",
		"code_generation":   "生成代码|编写代码|实现功能|开发|code|implement|generate",
		"file_operation":    "创建文件|修改文件|删除文件|读取文件|file|create file|edit file",
		"command_execution": "运行命令|执行命令|安装依赖|构建项目|run|execute|install|build",
		"error_fixing":      "修复错误|解决问题|调试|fix|debug|error|bug",
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
	if projectType, ok := parameters["project_type"].(string); ok && projectType != "" {
		return "project_creation"
	}

	if filePath, ok := parameters["file_path"].(string); ok && filePath != "" {
		return "file_operation"
	}

	if command, ok := parameters["command"].(string); ok && command != "" {
		return "command_execution"
	}

	return "generic"
}

// handleProjectCreation 处理项目创建任务
func (a *FullstackDeveloperAgent) handleProjectCreation(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理项目创建任务: %s", task)

	// 获取参数
	projectPath, _ := parameters["project_path"].(string)
	projectType, _ := parameters["project_type"].(string)
	projectName, _ := parameters["project_name"].(string)

	if projectPath == "" {
		return &AgentResult{
			Success: false,
			Output:  "project_path is required",
		}, nil
	}

	if projectType == "" {
		// 尝试从任务中推断项目类型
		taskLower := strings.ToLower(task)
		if strings.Contains(taskLower, "go") || strings.Contains(taskLower, "golang") {
			projectType = "go"
		} else if strings.Contains(taskLower, "react") {
			projectType = "react"
		} else if strings.Contains(taskLower, "node") || strings.Contains(taskLower, "javascript") {
			projectType = "nodejs"
		} else if strings.Contains(taskLower, "python") {
			projectType = "python"
		} else if strings.Contains(taskLower, "static") || strings.Contains(taskLower, "html") {
			projectType = "static"
		} else {
			projectType = "go" // 默认使用Go
		}
	}

	if projectName == "" {
		// 从路径中提取项目名称
		parts := strings.Split(projectPath, "/")
		projectName = parts[len(parts)-1]
		if projectName == "" && len(parts) > 1 {
			projectName = parts[len(parts)-2]
		}
		if projectName == "" {
			projectName = "myproject"
		}
	}

	// 使用项目管理工具初始化项目
	toolResult, err := a.ExecuteTool(ctx, "project_tool", map[string]interface{}{
		"operation":    "init",
		"project_path": projectPath,
		"project_type": projectType,
		"project_name": projectName,
	})

	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("failed to initialize project: %v", err),
		}, nil
	}

	if !toolResult.Success {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("project initialization failed: %v", toolResult.Error),
		}, nil
	}

	// 分析任务需求
	promptData := map[string]interface{}{
		"task":         task,
		"project_type": projectType,
		"project_name": projectName,
		"project_path": projectPath,
	}

	prompt, err := a.FormatPrompt("task_analysis", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("分析项目需求：%s\n项目类型：%s\n项目名称：%s", task, projectType, projectName)
	}

	// 调用模型进行分析
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		utils.Warn("模型分析失败: %v", err)
		// 继续执行，即使模型分析失败
	}

	// 收集产物
	artifacts := []AgentArtifact{
		{
			Type:        "project_structure",
			Name:        "项目结构",
			Content:     toolResult.Output,
			Description: "项目初始化结果",
		},
	}

	if modelResponse != nil {
		artifacts = append(artifacts, AgentArtifact{
			Type:        "task_analysis",
			Name:        "任务分析",
			Content:     modelResponse.Content,
			Description: "模型对任务需求的分析",
		})
	}

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"project_initialized": true,
			"project_type":        projectType,
			"project_name":        projectName,
			"project_path":        projectPath,
			"tool_result":         toolResult.Output,
			"analysis":            modelResponse,
		},
		Artifacts: artifacts,
		Logs: []string{
			fmt.Sprintf("项目初始化成功: %s", projectPath),
			fmt.Sprintf("项目类型: %s", projectType),
			fmt.Sprintf("项目名称: %s", projectName),
		},
	}, nil
}

// handleCodeGeneration 处理代码生成任务
func (a *FullstackDeveloperAgent) handleCodeGeneration(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理代码生成任务: %s", task)

	// 获取参数
	projectPath, _ := parameters["project_path"].(string)
	requirement, _ := parameters["requirement"].(string)
	techStack, _ := parameters["tech_stack"].(string)
	fileList, _ := parameters["file_list"].([]interface{})

	if requirement == "" {
		requirement = task
	}

	if projectPath == "" {
		projectPath = "./"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"requirement":  requirement,
		"project_type": "fullstack",
		"project_name": "code_generation",
		"tech_stack":   techStack,
		"file_list":    fileList,
	}

	prompt, err := a.FormatPrompt("code_generation", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("生成代码需求：%s\n技术栈：%s", requirement, techStack)
	}

	// 调用模型生成代码
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("代码生成失败: %v", err),
		}, nil
	}

	// 解析生成的代码（简化处理）
	generatedCode := modelResponse.Content

	// 创建代码文件（如果指定了文件列表）
	var createdFiles []string
	var artifacts []AgentArtifact

	if len(fileList) > 0 {
		for i, fileItem := range fileList {
			if filePath, ok := fileItem.(string); ok && filePath != "" {
				// 简单的文件创建逻辑
				content := fmt.Sprintf("// 文件: %s\n// 生成时间: %s\n\n%s", filePath, time.Now().Format("2006-01-02 15:04:05"), generatedCode)

				toolResult, err := a.ExecuteTool(ctx, "project_tool", map[string]interface{}{
					"operation":    "create_file",
					"project_path": projectPath,
					"file_path":    filePath,
					"file_content": content,
				})

				if err == nil && toolResult.Success {
					createdFiles = append(createdFiles, filePath)

					artifacts = append(artifacts, AgentArtifact{
						Type:        "code_file",
						Name:        filePath,
						Path:        filePath,
						Description: fmt.Sprintf("生成的代码文件 %d", i+1),
					})
				}
			}
		}
	}

	artifacts = append(artifacts, AgentArtifact{
		Type:        "generated_code",
		Name:        "生成的代码",
		Content:     generatedCode,
		Description: "模型生成的代码内容",
	})

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"code_generated": true,
			"generated_code": generatedCode,
			"created_files":  createdFiles,
			"model_response": modelResponse,
			"files_count":    len(createdFiles),
		},
		Artifacts: artifacts,
		Logs: []string{
			fmt.Sprintf("代码生成成功"),
			fmt.Sprintf("生成文件数量: %d", len(createdFiles)),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// handleFileOperation 处理文件操作任务
func (a *FullstackDeveloperAgent) handleFileOperation(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理文件操作任务: %s", task)

	// 获取参数
	operation, _ := parameters["operation"].(string)
	path, _ := parameters["path"].(string)
	content, _ := parameters["content"].(string)

	if operation == "" {
		// 从任务中推断操作类型
		taskLower := strings.ToLower(task)
		if strings.Contains(taskLower, "创建") || strings.Contains(taskLower, "create") || strings.Contains(taskLower, "新建") {
			operation = "write"
		} else if strings.Contains(taskLower, "读取") || strings.Contains(taskLower, "read") || strings.Contains(taskLower, "查看") {
			operation = "read"
		} else if strings.Contains(taskLower, "修改") || strings.Contains(taskLower, "edit") || strings.Contains(taskLower, "更新") {
			operation = "write"
		} else if strings.Contains(taskLower, "删除") || strings.Contains(taskLower, "delete") || strings.Contains(taskLower, "移除") {
			operation = "delete"
		} else {
			operation = "read"
		}
	}

	if path == "" {
		return &AgentResult{
			Success: false,
			Output:  "path is required for file operation",
		}, nil
	}

	// 使用文件系统工具执行操作
	toolResult, err := a.ExecuteTool(ctx, "fs_tool", map[string]interface{}{
		"operation": operation,
		"path":      path,
		"content":   content,
	})

	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("文件操作失败: %v", err),
		}, nil
	}

	if !toolResult.Success {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("文件操作失败: %v", toolResult.Error),
		}, nil
	}

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"operation":   operation,
			"path":        path,
			"tool_result": toolResult.Output,
		},
		Logs: []string{
			fmt.Sprintf("文件操作成功: %s %s", operation, path),
		},
	}, nil
}

// handleCommandExecution 处理命令执行任务
func (a *FullstackDeveloperAgent) handleCommandExecution(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理命令执行任务: %s", task)

	// 获取参数
	command, _ := parameters["command"].(string)
	args, _ := parameters["args"].([]interface{})
	workingDir, _ := parameters["working_dir"].(string)

	if command == "" {
		// 从任务中提取命令
		taskLower := strings.ToLower(task)
		if strings.Contains(taskLower, "安装") || strings.Contains(taskLower, "install") {
			command = "npm"
			args = []interface{}{"install"}
		} else if strings.Contains(taskLower, "运行") || strings.Contains(taskLower, "run") {
			command = "npm"
			args = []interface{}{"start"}
		} else if strings.Contains(taskLower, "构建") || strings.Contains(taskLower, "build") {
			command = "npm"
			args = []interface{}{"run", "build"}
		} else if strings.Contains(taskLower, "测试") || strings.Contains(taskLower, "test") {
			command = "npm"
			args = []interface{}{"test"}
		} else {
			return &AgentResult{
				Success: false,
				Output:  "command is required",
			}, nil
		}
	}

	if workingDir == "" {
		workingDir = "./"
	}

	// 使用Shell工具执行命令
	toolResult, err := a.ExecuteTool(ctx, "shell_tool", map[string]interface{}{
		"command":        command,
		"args":           args,
		"working_dir":    workingDir,
		"timeout":        60,
		"capture_output": true,
	})

	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("命令执行失败: %v", err),
		}, nil
	}

	if !toolResult.Success {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("命令执行失败: %v", toolResult.Error),
		}, nil
	}

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"command":     command,
			"args":        args,
			"working_dir": workingDir,
			"tool_result": toolResult.Output,
		},
		Logs: []string{
			fmt.Sprintf("命令执行成功: %s", command),
			fmt.Sprintf("工作目录: %s", workingDir),
		},
	}, nil
}

// handleErrorFixing 处理错误修复任务
func (a *FullstackDeveloperAgent) handleErrorFixing(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理错误修复任务: %s", task)

	// 获取参数
	errorMessage, _ := parameters["error_message"].(string)
	relatedCode, _ := parameters["related_code"].(string)
	errorContext, _ := parameters["error_context"].(string)

	if errorMessage == "" {
		errorMessage = task
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"error_message": errorMessage,
		"related_code":  relatedCode,
		"error_context": errorContext,
	}

	prompt, err := a.FormatPrompt("error_fixing", promptData)
	if err != nil {
		utils.Warn("Failed to format prompt: %v", err)
		prompt = fmt.Sprintf("修复错误：%s\n相关代码：%s\n错误上下文：%s", errorMessage, relatedCode, errorContext)
	}

	// 调用模型分析错误
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("错误分析失败: %v", err),
		}, nil
	}

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"error_analyzed": true,
			"analysis":       modelResponse.Content,
			"model_response": modelResponse,
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "error_analysis",
				Name:        "错误分析",
				Content:     modelResponse.Content,
				Description: "模型对错误的分析和修复建议",
			},
		},
		Logs: []string{
			"错误分析完成",
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// handleGenericTask 处理通用任务
func (a *FullstackDeveloperAgent) handleGenericTask(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	utils.Info("处理通用任务: %s", task)

	// 使用模型分析任务
	promptData := map[string]interface{}{
		"task":       task,
		"parameters": parameters,
	}

	prompt, err := a.FormatPrompt("task_analysis", promptData)
	if err != nil {
		prompt = fmt.Sprintf("请分析并执行以下任务：%s\n参数：%v", task, parameters)
	}

	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		return &AgentResult{
			Success: false,
			Output:  fmt.Sprintf("任务分析失败: %v", err),
		}, nil
	}

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"task_analyzed":  true,
			"analysis":       modelResponse.Content,
			"model_response": modelResponse,
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "task_analysis",
				Name:        "任务分析",
				Content:     modelResponse.Content,
				Description: "模型对任务的分析和建议",
			},
		},
		Logs: []string{
			"任务分析完成",
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
		},
	}, nil
}

// ValidateTask 验证任务（全栈开发Agent特定验证）
func (a *FullstackDeveloperAgent) ValidateTask(task string) error {
	// 调用基础验证
	if err := a.BaseAgent.ValidateTask(task); err != nil {
		return err
	}

	// 全栈开发Agent特定的验证
	taskLower := strings.ToLower(task)

	// 检查是否包含危险内容
	dangerousPatterns := []string{
		"rm -rf", "format c:", "delete all", "destroy",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(taskLower, pattern) {
			return fmt.Errorf("task contains dangerous pattern: %s", pattern)
		}
	}

	// 检查任务是否过于模糊
	if len(task) < 10 && !strings.Contains(taskLower, "help") && !strings.Contains(taskLower, "test") {
		return fmt.Errorf("task is too vague, please provide more details")
	}

	return nil
}

// GetCapabilitiesDescription 获取能力描述
func (a *FullstackDeveloperAgent) GetCapabilitiesDescription() string {
	descriptions := []string{}
	for _, cap := range a.config.Capabilities {
		switch cap {
		case CapabilityCodeGeneration:
			descriptions = append(descriptions, "代码生成")
		case CapabilityFileOperation:
			descriptions = append(descriptions, "文件操作")
		case CapabilityGitOperation:
			descriptions = append(descriptions, "Git操作")
		case CapabilityShellExecution:
			descriptions = append(descriptions, "命令执行")
		case CapabilityProjectSetup:
			descriptions = append(descriptions, "项目设置")
		case CapabilityTesting:
			descriptions = append(descriptions, "测试")
		case CapabilityReview:
			descriptions = append(descriptions, "代码审查")
		case CapabilityPlanning:
			descriptions = append(descriptions, "项目规划")
		}
	}

	return strings.Join(descriptions, ", ")
}
