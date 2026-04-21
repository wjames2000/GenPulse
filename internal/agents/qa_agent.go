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

// QAEngineerAgent QA工程师Agent
type QAEngineerAgent struct {
	*BaseAgent
}

// NewQAEngineerAgent 创建QA工程师Agent
func NewQAEngineerAgent(config AgentConfig, modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine) (*QAEngineerAgent, error) {
	// 设置QA工程师Agent的默认配置
	config.Role = RoleQAEngineer
	if config.Description == "" {
		config.Description = "QA工程师，负责测试用例生成与执行、质量保证"
	}

	// 设置默认能力
	if len(config.Capabilities) == 0 {
		config.Capabilities = []AgentCapability{
			CapabilityTesting,
			CapabilityQualityAssurance,
			CapabilityBugDetection,
		}
	}

	// 设置默认工具
	if len(config.Tools) == 0 {
		config.Tools = []string{
			"fs_tool",
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
		"test_case_generation": `你是一个QA工程师。请为以下功能生成测试用例：

功能名称：{{feature_name}}
功能描述：{{feature_description}}
技术栈：{{tech_stack}}
测试类型：{{test_type}}
测试范围：{{test_scope}}

请生成以下测试用例：
1. 功能测试用例
2. 边界测试用例
3. 异常测试用例
4. 性能测试用例
5. 安全测试用例
6. 兼容性测试用例
7. 用户体验测试用例

每个测试用例应包括：
- 测试用例ID
- 测试标题
- 前置条件
- 测试步骤
- 预期结果
- 实际结果（执行时填写）
- 测试状态（通过/失败）

请输出详细的测试用例：`,

		"test_execution": `你是一个QA工程师。请执行以下测试：

测试环境：{{test_environment}}
测试工具：{{test_tools}}
测试数据：{{test_data}}
被测系统：{{system_under_test}}
测试计划：{{test_plan}}

请执行以下操作：
1. 环境准备
2. 测试数据准备
3. 测试执行
4. 结果记录
5. 缺陷报告
6. 测试总结

请输出测试执行报告：`,

		"quality_assurance": `你是一个QA工程师。请进行质量保证：

项目名称：{{project_name}}
质量目标：{{quality_goals}}
质量标准：{{quality_standards}}
检查项：{{check_items}}

请进行以下质量检查：
1. 代码质量检查
2. 文档质量检查
3. 性能质量检查
4. 安全质量检查
5. 可用性检查
6. 可靠性检查
7. 可维护性检查
8. 兼容性检查

请输出质量保证报告：`,

		"bug_analysis": `你是一个QA工程师。请分析以下缺陷：

缺陷描述：{{bug_description}}
重现步骤：{{reproduce_steps}}
影响范围：{{impact_scope}}
严重程度：{{severity}}
优先级：{{priority}}

请进行以下分析：
1. 缺陷根本原因分析
2. 影响范围分析
3. 修复建议
4. 预防措施
5. 回归测试方案
6. 风险评估

请输出缺陷分析报告：`,
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

	agent := &QAEngineerAgent{
		BaseAgent: baseAgent,
	}

	return agent, nil
}

// Execute 执行QA工程师Agent任务
func (a *QAEngineerAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
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

	utils.Info("QA工程师Agent %s 开始执行任务: %s", a.config.Name, task)

	startTime := time.Now()

	// 根据任务类型选择提示词模板
	templateName := "test_case_generation"
	taskLower := strings.ToLower(task)

	if strings.Contains(taskLower, "执行") || strings.Contains(taskLower, "run") || strings.Contains(taskLower, "execute") {
		templateName = "test_execution"
	} else if strings.Contains(taskLower, "质量") || strings.Contains(taskLower, "quality") || strings.Contains(taskLower, "保证") {
		templateName = "quality_assurance"
	} else if strings.Contains(taskLower, "缺陷") || strings.Contains(taskLower, "bug") || strings.Contains(taskLower, "问题") {
		templateName = "bug_analysis"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"task":       task,
		"parameters": parameters,
	}

	// 添加特定参数
	if featureName, ok := parameters["feature_name"]; ok {
		promptData["feature_name"] = featureName
	}
	if featureDesc, ok := parameters["feature_description"]; ok {
		promptData["feature_description"] = featureDesc
	}
	if testType, ok := parameters["test_type"]; ok {
		promptData["test_type"] = testType
	}

	prompt, err := a.FormatPrompt(templateName, promptData)
	if err != nil {
		prompt = fmt.Sprintf("你是一个QA工程师。请根据以下需求进行测试和质量保证：\n\n需求：%s\n\n参数：%v\n\n请输出详细的测试和质量保证报告。", task, parameters)
	}

	// 调用模型生成测试用例
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		a.execution.State = StateFailed
		a.execution.Error = fmt.Sprintf("测试生成失败: %v", err)
		a.state = StateFailed
		a.execution.CompletedAt = &time.Time{}
		*a.execution.CompletedAt = time.Now()

		return &AgentResult{
			Success:  false,
			Output:   fmt.Sprintf("测试生成失败: %v", err),
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

	testCases := modelResponse.Content

	// 提取测试用例统计
	testStats := a.extractTestStatistics(testCases)

	result := &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"task_analyzed":   true,
			"test_cases":      testCases,
			"test_statistics": testStats,
			"model_response":  modelResponse,
			"agent_role":      a.config.Role,
			"agent_name":      a.config.Name,
			"template_used":   templateName,
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "test_cases",
				Name:        "测试用例",
				Content:     testCases,
				Description: fmt.Sprintf("%s生成的测试用例", a.config.Name),
			},
			{
				Type:        "test_statistics",
				Name:        "测试统计",
				Content:     fmt.Sprintf("%v", testStats),
				Description: "测试用例统计信息",
			},
		},
		Duration: duration,
		Logs: []string{
			fmt.Sprintf("测试生成完成"),
			fmt.Sprintf("使用模板: %s", templateName),
			fmt.Sprintf("测试用例数量: %d", testStats["total_test_cases"]),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
			fmt.Sprintf("执行时间: %v", duration),
		},
	}

	a.execution.Result = result

	return result, nil
}

// extractTestStatistics 从测试用例内容中提取统计信息
func (a *QAEngineerAgent) extractTestStatistics(testCases string) map[string]interface{} {
	stats := map[string]interface{}{
		"total_test_cases":    0,
		"functional_tests":    0,
		"boundary_tests":      0,
		"exception_tests":     0,
		"performance_tests":   0,
		"security_tests":      0,
		"compatibility_tests": 0,
		"usability_tests":     0,
	}

	lines := strings.Split(testCases, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 统计测试用例数量
		if strings.Contains(trimmedLine, "测试用例") || strings.Contains(trimmedLine, "TestCase") ||
			strings.Contains(trimmedLine, "TC-") || strings.Contains(trimmedLine, "TC_") {
			stats["total_test_cases"] = stats["total_test_cases"].(int) + 1
		}

		// 统计不同类型的测试用例
		if strings.Contains(trimmedLine, "功能测试") || strings.Contains(trimmedLine, "functional") {
			stats["functional_tests"] = stats["functional_tests"].(int) + 1
		} else if strings.Contains(trimmedLine, "边界测试") || strings.Contains(trimmedLine, "boundary") {
			stats["boundary_tests"] = stats["boundary_tests"].(int) + 1
		} else if strings.Contains(trimmedLine, "异常测试") || strings.Contains(trimmedLine, "exception") {
			stats["exception_tests"] = stats["exception_tests"].(int) + 1
		} else if strings.Contains(trimmedLine, "性能测试") || strings.Contains(trimmedLine, "performance") {
			stats["performance_tests"] = stats["performance_tests"].(int) + 1
		} else if strings.Contains(trimmedLine, "安全测试") || strings.Contains(trimmedLine, "security") {
			stats["security_tests"] = stats["security_tests"].(int) + 1
		} else if strings.Contains(trimmedLine, "兼容性测试") || strings.Contains(trimmedLine, "compatibility") {
			stats["compatibility_tests"] = stats["compatibility_tests"].(int) + 1
		} else if strings.Contains(trimmedLine, "用户体验测试") || strings.Contains(trimmedLine, "usability") {
			stats["usability_tests"] = stats["usability_tests"].(int) + 1
		}
	}

	return stats
}

// ValidateTask 验证任务（QA工程师特定验证）
func (a *QAEngineerAgent) ValidateTask(task string) error {
	// 调用基础验证
	if err := a.BaseAgent.ValidateTask(task); err != nil {
		return err
	}

	// QA工程师特定验证
	taskLower := strings.ToLower(task)

	// 检查是否包含测试相关关键词
	testingKeywords := []string{
		"测试", "质量", "缺陷", "bug", "问题", "检查", "验证", "保证",
		"test", "quality", "bug", "issue", "check", "verify", "assurance",
		"用例", "执行", "报告", "分析", "case", "execute", "report", "analysis",
	}

	hasTestingKeyword := false
	for _, keyword := range testingKeywords {
		if strings.Contains(taskLower, keyword) {
			hasTestingKeyword = true
			break
		}
	}

	if !hasTestingKeyword {
		utils.Warn("任务可能不适用于QA工程师: %s", task)
	}

	return nil
}

// GetRoleDescription 获取角色描述
func (a *QAEngineerAgent) GetRoleDescription() string {
	return "QA工程师：负责测试用例生成与执行、质量保证"
}

// GenerateTestCases 生成测试用例（专用方法）
func (a *QAEngineerAgent) GenerateTestCases(ctx context.Context, featureName, featureDescription, testType string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"feature_name":        featureName,
		"feature_description": featureDescription,
		"test_type":           testType,
	}

	task := fmt.Sprintf("为功能'%s'生成测试用例", featureName)
	return a.Execute(ctx, task, parameters)
}

// ExecuteTests 执行测试（专用方法）
func (a *QAEngineerAgent) ExecuteTests(ctx context.Context, testEnvironment, testPlan string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"test_environment": testEnvironment,
		"test_plan":        testPlan,
	}

	task := fmt.Sprintf("在环境'%s'执行测试", testEnvironment)
	return a.Execute(ctx, task, parameters)
}

// PerformQualityAssurance 进行质量保证（专用方法）
func (a *QAEngineerAgent) PerformQualityAssurance(ctx context.Context, projectName, qualityGoals string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"project_name":  projectName,
		"quality_goals": qualityGoals,
	}

	task := fmt.Sprintf("为项目'%s'进行质量保证", projectName)
	return a.Execute(ctx, task, parameters)
}

// AnalyzeBug 分析缺陷（专用方法）
func (a *QAEngineerAgent) AnalyzeBug(ctx context.Context, bugDescription, severity string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"bug_description": bugDescription,
		"severity":        severity,
	}

	task := fmt.Sprintf("分析缺陷'%s'", bugDescription)
	return a.Execute(ctx, task, parameters)
}
