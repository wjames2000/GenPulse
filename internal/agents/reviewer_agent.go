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

// ReviewerAgent 代码审查员Agent
type ReviewerAgent struct {
	*BaseAgent
}

// NewReviewerAgent 创建代码审查员Agent
func NewReviewerAgent(config AgentConfig, modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine) (*ReviewerAgent, error) {
	// 设置代码审查员Agent的默认配置
	config.Role = RoleReviewer
	if config.Description == "" {
		config.Description = "代码审查员，负责代码审查、安全扫描和质量检查"
	}

	// 设置默认能力
	if len(config.Capabilities) == 0 {
		config.Capabilities = []AgentCapability{
			CapabilityReview,
			CapabilitySecurityScan,
			CapabilityQualityCheck,
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
		"code_review": `你是一个代码审查员。请审查以下代码：

代码文件：{{code_file}}
代码内容：{{code_content}}
编程语言：{{programming_language}}
代码规范：{{coding_standards}}
审查重点：{{review_focus}}

请进行以下审查：
1. 代码质量审查
2. 代码规范审查
3. 性能审查
4. 安全审查
5. 可维护性审查
6. 可测试性审查
7. 文档审查
8. 最佳实践审查

请输出详细的代码审查报告：`,

		"security_scan": `你是一个代码审查员。请进行安全扫描：

扫描目标：{{scan_target}}
扫描类型：{{scan_type}}
安全标准：{{security_standards}}
风险等级：{{risk_level}}

请检查以下安全漏洞：
1. 注入漏洞（SQL注入、命令注入等）
2. 跨站脚本（XSS）
3. 跨站请求伪造（CSRF）
4. 身份验证和授权问题
5. 敏感数据泄露
6. 安全配置错误
7. 不安全的反序列化
8. 使用已知漏洞的组件
9. 不足的日志记录和监控
10. 其他OWASP Top 10漏洞

请输出详细的安全扫描报告：`,

		"quality_check": `你是一个代码审查员。请进行质量检查：

检查目标：{{check_target}}
质量标准：{{quality_standards}}
检查项：{{check_items}}
验收标准：{{acceptance_criteria}}

请检查以下质量指标：
1. 代码复杂度
2. 代码重复率
3. 测试覆盖率
4. 文档完整性
5. 性能指标
6. 可靠性指标
7. 可用性指标
8. 可维护性指标
9. 可扩展性指标
10. 合规性检查

请输出详细的质量检查报告：`,

		"architecture_review": `你是一个代码审查员。请进行架构审查：

架构目标：{{architecture_target}}
架构原则：{{architecture_principles}}
技术栈：{{tech_stack}}
设计模式：{{design_patterns}}

请审查以下架构方面：
1. 架构一致性
2. 组件设计
3. 接口设计
4. 数据流设计
5. 错误处理设计
6. 安全设计
7. 性能设计
8. 扩展性设计
9. 部署设计
10. 监控设计

请输出详细的架构审查报告：`,
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

	agent := &ReviewerAgent{
		BaseAgent: baseAgent,
	}

	return agent, nil
}

// Execute 执行代码审查员Agent任务
func (a *ReviewerAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
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

	utils.Info("代码审查员Agent %s 开始执行任务: %s", a.config.Name, task)

	startTime := time.Now()

	// 根据任务类型选择提示词模板
	templateName := "code_review"
	taskLower := strings.ToLower(task)

	if strings.Contains(taskLower, "安全") || strings.Contains(taskLower, "security") || strings.Contains(taskLower, "漏洞") {
		templateName = "security_scan"
	} else if strings.Contains(taskLower, "质量") || strings.Contains(taskLower, "quality") || strings.Contains(taskLower, "检查") {
		templateName = "quality_check"
	} else if strings.Contains(taskLower, "架构") || strings.Contains(taskLower, "architecture") || strings.Contains(taskLower, "设计") {
		templateName = "architecture_review"
	}

	// 格式化提示词
	promptData := map[string]interface{}{
		"task":       task,
		"parameters": parameters,
	}

	// 添加特定参数
	if codeFile, ok := parameters["code_file"]; ok {
		promptData["code_file"] = codeFile
	}
	if codeContent, ok := parameters["code_content"]; ok {
		promptData["code_content"] = codeContent
	}
	if programmingLanguage, ok := parameters["programming_language"]; ok {
		promptData["programming_language"] = programmingLanguage
	}

	prompt, err := a.FormatPrompt(templateName, promptData)
	if err != nil {
		prompt = fmt.Sprintf("你是一个代码审查员。请根据以下需求进行代码审查：\n\n需求：%s\n\n参数：%v\n\n请输出详细的审查报告。", task, parameters)
	}

	// 调用模型进行审查
	modelResponse, err := a.GenerateWithModel(ctx, prompt)
	if err != nil {
		a.execution.State = StateFailed
		a.execution.Error = fmt.Sprintf("代码审查失败: %v", err)
		a.state = StateFailed
		a.execution.CompletedAt = &time.Time{}
		*a.execution.CompletedAt = time.Now()

		return &AgentResult{
			Success:  false,
			Output:   fmt.Sprintf("代码审查失败: %v", err),
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

	reviewReport := modelResponse.Content

	// 提取问题和建议
	issues := a.extractIssues(reviewReport)
	recommendations := a.extractRecommendations(reviewReport)

	result := &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"task_analyzed":   true,
			"review_report":   reviewReport,
			"issues_found":    issues,
			"recommendations": recommendations,
			"model_response":  modelResponse,
			"agent_role":      a.config.Role,
			"agent_name":      a.config.Name,
			"template_used":   templateName,
		},
		Artifacts: []AgentArtifact{
			{
				Type:        "review_report",
				Name:        "审查报告",
				Content:     reviewReport,
				Description: fmt.Sprintf("%s生成的审查报告", a.config.Name),
			},
			{
				Type:        "issues",
				Name:        "发现问题",
				Content:     strings.Join(issues, "\n"),
				Description: "审查发现的问题列表",
			},
			{
				Type:        "recommendations",
				Name:        "改进建议",
				Content:     strings.Join(recommendations, "\n"),
				Description: "审查提出的改进建议",
			},
		},
		Duration: duration,
		Logs: []string{
			fmt.Sprintf("代码审查完成"),
			fmt.Sprintf("使用模板: %s", templateName),
			fmt.Sprintf("发现问题: %d个", len(issues)),
			fmt.Sprintf("提出建议: %d个", len(recommendations)),
			fmt.Sprintf("模型Tokens: %d", modelResponse.Usage.TotalTokens),
			fmt.Sprintf("执行时间: %v", duration),
		},
	}

	a.execution.Result = result

	return result, nil
}

// extractIssues 从审查报告中提取问题
func (a *ReviewerAgent) extractIssues(reviewReport string) []string {
	var issues []string

	lines := strings.Split(reviewReport, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 检查问题模式
		if strings.Contains(trimmedLine, "问题") || strings.Contains(trimmedLine, "issue") ||
			strings.Contains(trimmedLine, "缺陷") || strings.Contains(trimmedLine, "bug") ||
			strings.Contains(trimmedLine, "漏洞") || strings.Contains(trimmedLine, "vulnerability") ||
			strings.Contains(trimmedLine, "错误") || strings.Contains(trimmedLine, "error") ||
			strings.Contains(trimmedLine, "风险") || strings.Contains(trimmedLine, "risk") ||
			strings.Contains(trimmedLine, "❌") || strings.Contains(trimmedLine, "⚠️") ||
			strings.Contains(trimmedLine, "[问题]") || strings.Contains(trimmedLine, "[ISSUE]") {

			// 清理问题描述
			cleanIssue := strings.TrimPrefix(trimmedLine, "❌ ")
			cleanIssue = strings.TrimPrefix(cleanIssue, "⚠️ ")
			cleanIssue = strings.TrimPrefix(cleanIssue, "[问题] ")
			cleanIssue = strings.TrimPrefix(cleanIssue, "[ISSUE] ")

			if cleanIssue != "" && len(cleanIssue) > 5 {
				issues = append(issues, cleanIssue)
			}
		}

		// 检查编号列表中的问题
		if strings.HasPrefix(trimmedLine, "1.") || strings.HasPrefix(trimmedLine, "2.") ||
			strings.HasPrefix(trimmedLine, "3.") || strings.HasPrefix(trimmedLine, "4.") ||
			strings.HasPrefix(trimmedLine, "5.") || strings.HasPrefix(trimmedLine, "6.") ||
			strings.HasPrefix(trimmedLine, "7.") || strings.HasPrefix(trimmedLine, "8.") ||
			strings.HasPrefix(trimmedLine, "9.") || strings.HasPrefix(trimmedLine, "10.") {

			// 检查是否包含问题关键词
			if strings.Contains(trimmedLine, "问题") || strings.Contains(trimmedLine, "缺陷") ||
				strings.Contains(trimmedLine, "漏洞") || strings.Contains(trimmedLine, "错误") {
				issues = append(issues, trimmedLine)
			}
		}
	}

	// 去重
	uniqueIssues := make(map[string]bool)
	var result []string
	for _, issue := range issues {
		if !uniqueIssues[issue] {
			uniqueIssues[issue] = true
			result = append(result, issue)
		}
	}

	return result
}

// extractRecommendations 从审查报告中提取建议
func (a *ReviewerAgent) extractRecommendations(reviewReport string) []string {
	var recommendations []string

	lines := strings.Split(reviewReport, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 检查建议模式
		if strings.Contains(trimmedLine, "建议") || strings.Contains(trimmedLine, "recommendation") ||
			strings.Contains(trimmedLine, "改进") || strings.Contains(trimmedLine, "improvement") ||
			strings.Contains(trimmedLine, "优化") || strings.Contains(trimmedLine, "optimization") ||
			strings.Contains(trimmedLine, "修复") || strings.Contains(trimmedLine, "fix") ||
			strings.Contains(trimmedLine, "✅") || strings.Contains(trimmedLine, "💡") ||
			strings.Contains(trimmedLine, "[建议]") || strings.Contains(trimmedLine, "[RECOMMENDATION]") {

			// 清理建议描述
			cleanRec := strings.TrimPrefix(trimmedLine, "✅ ")
			cleanRec = strings.TrimPrefix(cleanRec, "💡 ")
			cleanRec = strings.TrimPrefix(cleanRec, "[建议] ")
			cleanRec = strings.TrimPrefix(cleanRec, "[RECOMMENDATION] ")

			if cleanRec != "" && len(cleanRec) > 5 {
				recommendations = append(recommendations, cleanRec)
			}
		}

		// 检查建议列表
		if strings.HasPrefix(trimmedLine, "- ") || strings.HasPrefix(trimmedLine, "* ") ||
			strings.HasPrefix(trimmedLine, "• ") {

			// 检查是否包含建议关键词
			if strings.Contains(trimmedLine, "建议") || strings.Contains(trimmedLine, "改进") ||
				strings.Contains(trimmedLine, "优化") || strings.Contains(trimmedLine, "修复") {
				recommendations = append(recommendations, trimmedLine)
			}
		}
	}

	// 去重
	uniqueRecs := make(map[string]bool)
	var result []string
	for _, rec := range recommendations {
		if !uniqueRecs[rec] {
			uniqueRecs[rec] = true
			result = append(result, rec)
		}
	}

	return result
}

// ValidateTask 验证任务（代码审查员特定验证）
func (a *ReviewerAgent) ValidateTask(task string) error {
	// 调用基础验证
	if err := a.BaseAgent.ValidateTask(task); err != nil {
		return err
	}

	// 代码审查员特定验证
	taskLower := strings.ToLower(task)

	// 检查是否包含审查相关关键词
	reviewKeywords := []string{
		"审查", "检查", "扫描", "审核", "评估", "分析", "review", "check", "scan", "audit", "evaluate", "analyze",
		"代码", "安全", "质量", "架构", "code", "security", "quality", "architecture",
		"问题", "缺陷", "漏洞", "建议", "改进", "issue", "bug", "vulnerability", "recommendation", "improvement",
	}

	hasReviewKeyword := false
	for _, keyword := range reviewKeywords {
		if strings.Contains(taskLower, keyword) {
			hasReviewKeyword = true
			break
		}
	}

	if !hasReviewKeyword {
		utils.Warn("任务可能不适用于代码审查员: %s", task)
	}

	return nil
}

// GetRoleDescription 获取角色描述
func (a *ReviewerAgent) GetRoleDescription() string {
	return "代码审查员：负责代码审查、安全扫描和质量检查"
}

// ReviewCode 审查代码（专用方法）
func (a *ReviewerAgent) ReviewCode(ctx context.Context, codeFile, codeContent, programmingLanguage string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"code_file":            codeFile,
		"code_content":         codeContent,
		"programming_language": programmingLanguage,
	}

	task := fmt.Sprintf("审查代码文件'%s'", codeFile)
	return a.Execute(ctx, task, parameters)
}

// ScanSecurity 安全扫描（专用方法）
func (a *ReviewerAgent) ScanSecurity(ctx context.Context, scanTarget, scanType string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"scan_target": scanTarget,
		"scan_type":   scanType,
	}

	task := fmt.Sprintf("对'%s'进行安全扫描", scanTarget)
	return a.Execute(ctx, task, parameters)
}

// CheckQuality 质量检查（专用方法）
func (a *ReviewerAgent) CheckQuality(ctx context.Context, checkTarget, qualityStandards string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"check_target":      checkTarget,
		"quality_standards": qualityStandards,
	}

	task := fmt.Sprintf("对'%s'进行质量检查", checkTarget)
	return a.Execute(ctx, task, parameters)
}

// ReviewArchitecture 架构审查（专用方法）
func (a *ReviewerAgent) ReviewArchitecture(ctx context.Context, architectureTarget, techStack string) (*AgentResult, error) {
	parameters := map[string]interface{}{
		"architecture_target": architectureTarget,
		"tech_stack":          techStack,
	}

	task := fmt.Sprintf("审查'%s'架构", architectureTarget)
	return a.Execute(ctx, task, parameters)
}
