package skills

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Validator 技能验证器
type Validator struct {
	rules        []ValidationRule
	patterns     map[string]*regexp.Regexp
	safeTools    map[string]bool
	dangerousOps []string
}

// ValidationRule 验证规则
type ValidationRule struct {
	Name        string                        `json:"name" yaml:"name"`
	Description string                        `json:"description" yaml:"description"`
	Severity    string                        `json:"severity" yaml:"severity"` // low, medium, high, critical
	Check       func(*Skill) ValidationResult `json:"-" yaml:"-"`
}

// ValidationResult 验证结果
type ValidationResult struct {
	RuleName    string   `json:"rule_name" yaml:"rule_name"`
	Passed      bool     `json:"passed" yaml:"passed"`
	Severity    string   `json:"severity" yaml:"severity"`
	Message     string   `json:"message" yaml:"message"`
	Details     []string `json:"details,omitempty" yaml:"details,omitempty"`
	Suggestions []string `json:"suggestions,omitempty" yaml:"suggestions,omitempty"`
}

// ValidationReport 验证报告
type ValidationReport struct {
	SkillID          string             `json:"skill_id" yaml:"skill_id"`
	SkillName        string             `json:"skill_name" yaml:"skill_name"`
	OverallPass      bool               `json:"overall_pass" yaml:"overall_pass"`
	TotalChecks      int                `json:"total_checks" yaml:"total_checks"`
	PassedChecks     int                `json:"passed_checks" yaml:"passed_checks"`
	FailedChecks     int                `json:"failed_checks" yaml:"failed_checks"`
	Results          []ValidationResult `json:"results" yaml:"results"`
	CriticalFailures []string           `json:"critical_failures,omitempty" yaml:"critical_failures,omitempty"`
	Warnings         []string           `json:"warnings,omitempty" yaml:"warnings,omitempty"`
	Timestamp        string             `json:"timestamp" yaml:"timestamp"`
}

// NewValidator 创建验证器
func NewValidator() *Validator {
	v := &Validator{
		patterns:  make(map[string]*regexp.Regexp),
		safeTools: make(map[string]bool),
		dangerousOps: []string{
			"rm ", "del ", "format ", "shutdown ", "reboot ", "kill ",
			"chmod 777", "chown root", "sudo ", "su ", "passwd ",
		},
	}

	// 编译正则表达式
	v.compilePatterns()

	// 初始化安全工具列表
	v.initSafeTools()

	// 添加验证规则
	v.initRules()

	return v
}

// Validate 验证技能
func (v *Validator) Validate(skill *Skill) *ValidationReport {
	report := &ValidationReport{
		SkillID:     skill.ID,
		SkillName:   skill.Name,
		TotalChecks: len(v.rules),
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
	}

	var results []ValidationResult
	var criticalFailures []string
	var warnings []string

	// 执行所有验证规则
	for _, rule := range v.rules {
		result := rule.Check(skill)
		results = append(results, result)

		if !result.Passed {
			report.FailedChecks++

			// 记录严重失败
			if result.Severity == "critical" {
				criticalFailures = append(criticalFailures, result.Message)
			} else if result.Severity == "high" {
				warnings = append(warnings, result.Message)
			}
		} else {
			report.PassedChecks++
		}
	}

	report.Results = results
	report.CriticalFailures = criticalFailures
	report.Warnings = warnings

	// 如果没有严重失败，则整体通过
	report.OverallPass = len(criticalFailures) == 0

	return report
}

// ValidateWithOptions 使用选项验证技能
func (v *Validator) ValidateWithOptions(skill *Skill, options ValidationOptions) *ValidationReport {
	report := v.Validate(skill)

	// 应用选项
	if options.StrictMode {
		// 严格模式下，任何失败都导致整体不通过
		report.OverallPass = report.FailedChecks == 0
	}

	if options.IgnoreWarnings {
		// 忽略警告，只关注严重失败
		report.OverallPass = len(report.CriticalFailures) == 0
	}

	return report
}

// AddRule 添加验证规则
func (v *Validator) AddRule(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

// RemoveRule 移除验证规则
func (v *Validator) RemoveRule(ruleName string) bool {
	for i, rule := range v.rules {
		if rule.Name == ruleName {
			v.rules = append(v.rules[:i], v.rules[i+1:]...)
			return true
		}
	}
	return false
}

// GetRules 获取所有验证规则
func (v *Validator) GetRules() []ValidationRule {
	return v.rules
}

// compilePatterns 编译正则表达式模式
func (v *Validator) compilePatterns() {
	patterns := map[string]string{
		"dangerous_command": `(?i)(rm\s+-rf|del\s+/q|format\s+|shutdown\s+|reboot\s+|kill\s+-9|chmod\s+777|chown\s+root|sudo\s+|su\s+|passwd\s+)`,
		"sensitive_path":    `(?i)(/etc/passwd|/etc/shadow|/root/|C:\\Windows\\System32|/bin/bash|/sbin/)`,
		"ip_address":        `\b(?:\d{1,3}\.){3}\d{1,3}\b`,
		"url":               `https?://[^\s]+`,
		"email":             `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`,
		"shell_injection":   `[$&|;` + "`" + `]`,
	}

	for name, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Printf("Warning: failed to compile pattern %s: %v\n", name, err)
			continue
		}
		v.patterns[name] = re
	}
}

// initSafeTools 初始化安全工具列表
func (v *Validator) initSafeTools() {
	safeTools := []string{
		"fs_read", "fs_list", "fs_stat",
		"git_status", "git_log", "git_diff",
		"shell_exec", // 需要额外检查
		"project_init", "dep_install",
	}

	for _, tool := range safeTools {
		v.safeTools[tool] = true
	}
}

// initRules 初始化验证规则
func (v *Validator) initRules() {
	v.rules = []ValidationRule{
		{
			Name:        "skill_basic_info",
			Description: "检查技能基本信息完整性",
			Severity:    "medium",
			Check:       v.checkBasicInfo,
		},
		{
			Name:        "skill_steps_integrity",
			Description: "检查技能步骤完整性",
			Severity:    "high",
			Check:       v.checkStepsIntegrity,
		},
		{
			Name:        "dangerous_operations",
			Description: "检查危险操作",
			Severity:    "critical",
			Check:       v.checkDangerousOperations,
		},
		{
			Name:        "tool_safety",
			Description: "检查工具安全性",
			Severity:    "high",
			Check:       v.checkToolSafety,
		},
		{
			Name:        "parameter_validation",
			Description: "检查参数验证",
			Severity:    "medium",
			Check:       v.checkParameterValidation,
		},
		{
			Name:        "sensitive_data",
			Description: "检查敏感数据泄露",
			Severity:    "critical",
			Check:       v.checkSensitiveData,
		},
		{
			Name:        "step_ordering",
			Description: "检查步骤顺序合理性",
			Severity:    "low",
			Check:       v.checkStepOrdering,
		},
		{
			Name:        "condition_validity",
			Description: "检查条件有效性",
			Severity:    "medium",
			Check:       v.checkConditionValidity,
		},
		{
			Name:        "expected_result_clarity",
			Description: "检查预期结果清晰度",
			Severity:    "low",
			Check:       v.checkExpectedResultClarity,
		},
		{
			Name:        "skill_complexity_match",
			Description: "检查技能复杂度匹配",
			Severity:    "low",
			Check:       v.checkSkillComplexityMatch,
		},
	}
}

// checkBasicInfo 检查基本信息
func (v *Validator) checkBasicInfo(skill *Skill) ValidationResult {
	result := ValidationResult{
		RuleName: "skill_basic_info",
		Severity: "medium",
		Passed:   true,
	}

	var issues []string

	// 检查名称
	if skill.Name == "" {
		issues = append(issues, "技能名称为空")
		result.Passed = false
	}

	// 检查描述
	if skill.Description == "" {
		issues = append(issues, "技能描述为空")
		result.Passed = false
	} else if len(skill.Description) < 10 {
		issues = append(issues, "技能描述过短（至少10个字符）")
		result.Passed = false
	}

	// 检查版本
	if skill.Version == "" {
		issues = append(issues, "技能版本为空")
		result.Passed = false
	}

	// 检查分类
	if skill.Category == "" {
		issues = append(issues, "技能分类为空")
		result.Passed = false
	}

	if len(issues) > 0 {
		result.Message = "技能基本信息不完整"
		result.Details = issues
		result.Suggestions = []string{
			"填写完整的技能名称和描述",
			"设置有效的版本号",
			"选择适当的分类",
		}
	} else {
		result.Message = "技能基本信息完整"
	}

	return result
}

// checkStepsIntegrity 检查步骤完整性
func (v *Validator) checkStepsIntegrity(skill *Skill) ValidationResult {
	result := ValidationResult{
		RuleName: "skill_steps_integrity",
		Severity: "high",
		Passed:   true,
	}

	if len(skill.Steps) == 0 {
		result.Passed = false
		result.Message = "技能没有定义任何步骤"
		result.Suggestions = []string{"至少定义一个执行步骤"}
		return result
	}

	var issues []string

	// 检查步骤顺序
	expectedOrder := 1
	for i, step := range skill.Steps {
		if step.Order != expectedOrder {
			issues = append(issues, fmt.Sprintf("步骤 %d 的顺序错误: 期望 %d, 实际 %d", i+1, expectedOrder, step.Order))
			result.Passed = false
		}
		expectedOrder++

		// 检查步骤动作
		if step.Action == "" {
			issues = append(issues, fmt.Sprintf("步骤 %d 的动作描述为空", step.Order))
			result.Passed = false
		}

		// 检查步骤ID
		if step.ID == "" {
			issues = append(issues, fmt.Sprintf("步骤 %d 的ID为空", step.Order))
			result.Passed = false
		}
	}

	if len(issues) > 0 {
		result.Message = "技能步骤完整性检查失败"
		result.Details = issues
		result.Suggestions = []string{
			"确保步骤顺序正确（从1开始连续）",
			"为每个步骤提供清晰的行动描述",
			"为每个步骤设置唯一的ID",
		}
	} else {
		result.Message = "技能步骤完整性检查通过"
	}

	return result
}

// checkDangerousOperations 检查危险操作
func (v *Validator) checkDangerousOperations(skill *Skill) ValidationResult {
	result := ValidationResult{
		RuleName: "dangerous_operations",
		Severity: "critical",
		Passed:   true,
	}

	var dangerousOps []string

	// 检查步骤中的危险操作
	for _, step := range skill.Steps {
		// 检查动作描述中的危险命令
		if v.patterns["dangerous_command"] != nil {
			if matches := v.patterns["dangerous_command"].FindAllString(step.Action, -1); len(matches) > 0 {
				dangerousOps = append(dangerousOps, fmt.Sprintf("步骤 %d: 发现危险命令 - %v", step.Order, matches))
				result.Passed = false
			}
		}

		// 检查参数中的危险内容
		for _, param := range step.Parameters {
			if param.Default != nil {
				defaultStr := fmt.Sprintf("%v", param.Default)
				if v.patterns["dangerous_command"] != nil {
					if matches := v.patterns["dangerous_command"].FindAllString(defaultStr, -1); len(matches) > 0 {
						dangerousOps = append(dangerousOps,
							fmt.Sprintf("步骤 %d 参数 %s: 默认值包含危险命令 - %v", step.Order, param.Name, matches))
						result.Passed = false
					}
				}
			}
		}

		// 检查预期结果中的危险内容
		if v.patterns["dangerous_command"] != nil {
			if matches := v.patterns["dangerous_command"].FindAllString(step.Expected, -1); len(matches) > 0 {
				dangerousOps = append(dangerousOps, fmt.Sprintf("步骤 %d 预期结果: 包含危险命令 - %v", step.Order, matches))
				result.Passed = false
			}
		}
	}

	// 检查示例中的危险操作
	for i, example := range skill.Examples {
		if v.patterns["dangerous_command"] != nil {
			if matches := v.patterns["dangerous_command"].FindAllString(example, -1); len(matches) > 0 {
				dangerousOps = append(dangerousOps, fmt.Sprintf("示例 %d: 包含危险命令 - %v", i+1, matches))
				result.Passed = false
			}
		}
	}

	if len(dangerousOps) > 0 {
		result.Message = "发现危险操作"
		result.Details = dangerousOps
		result.Suggestions = []string{
			"移除或替换危险命令",
			"添加明确的警告说明",
			"考虑使用更安全的替代方案",
		}
	} else {
		result.Message = "未发现危险操作"
	}

	return result
}

// checkToolSafety 检查工具安全性
func (v *Validator) checkToolSafety(skill *Skill) ValidationResult {
	result := ValidationResult{
		RuleName: "tool_safety",
		Severity: "high",
		Passed:   true,
	}

	var unsafeTools []string

	// 检查相关工具
	for _, tool := range skill.RelatedTools {
		if !v.safeTools[tool] {
			unsafeTools = append(unsafeTools, tool)
			result.Passed = false
		}
	}

	// 检查步骤中使用的工具
	for _, step := range skill.Steps {
		if step.Tool != "" && !v.safeTools[step.Tool] {
			unsafeTools = append(unsafeTools, step.Tool)
			result.Passed = false
		}
	}

	if len(unsafeTools) > 0 {
		result.Message = "发现不安全或未知工具"
		result.Details = unsafeTools
		result.Suggestions = []string{
			"只使用经过验证的安全工具",
			"为未知工具添加安全说明",
			"考虑使用工具白名单",
		}
	} else {
		result.Message = "工具安全性检查通过"
	}

	return result
}

// checkParameterValidation 检查参数验证
func (v *Validator) checkParameterValidation(skill *Skill) ValidationResult {
	result := ValidationResult{
		RuleName: "parameter_validation",
		Severity: "medium",
		Passed:   true,
	}

	var unvalidatedParams []string

	for _, step := range skill.Steps {
		for _, param := range step.Parameters {
			// 检查必填参数是否有默认值
			if param.Required && param.Default == nil {
				// 这是可以的，但可以记录
			}

			// 检查是否有约束说明
			if param.Constraints == "" && param.Type != "string" {
				unvalidatedParams = append(unvalidatedParams,
					fmt.Sprintf("步骤 %d 参数 %s: 缺少约束说明", step.Order, param.Name))
				result.Passed = false
			}

			// 检查参数描述
			if param.Description == "" {
				unvalidatedParams = append(unvalidatedParams,
					fmt.Sprintf("步骤 %d 参数 %s: 缺少描述", step.Order, param.Name))
				result.Passed = false
			}
		}
	}

	if len(unvalidatedParams) > 0 {
		result.Message = "参数验证不充分"
		result.Details = unvalidatedParams
		result.Suggestions = []string{
			"为所有参数提供清晰的描述",
			"为复杂类型参数添加约束说明",
			"考虑参数验证逻辑",
		}
	} else {
		result.Message = "参数验证充分"
	}

	return result
}

// checkSensitiveData 检查敏感数据
func (v *Validator) checkSensitiveData(skill *Skill) ValidationResult {
	result := ValidationResult{
		RuleName: "sensitive_data",
		Severity: "critical",
		Passed:   true,
	}

	var sensitiveData []string

	// 检查所有文本字段
	textFields := []string{
		skill.Name,
		skill.Description,
	}

	for _, step := range skill.Steps {
		textFields = append(textFields, step.Action, step.Expected)
		for _, param := range step.Parameters {
			if param.Default != nil {
				textFields = append(textFields, fmt.Sprintf("%v", param.Default))
			}
		}
	}

	for _, example := range skill.Examples {
		textFields = append(textFields, example)
	}

	for _, tip := range skill.Tips {
		textFields = append(textFields, tip)
	}

	for _, warning := range skill.Warnings {
		textFields = append(textFields, warning)
	}

	// 检查敏感路径
	if v.patterns["sensitive_path"] != nil {
		for _, text := range textFields {
			if matches := v.patterns["sensitive_path"].FindAllString(text, -1); len(matches) > 0 {
				sensitiveData = append(sensitiveData, fmt.Sprintf("发现敏感路径: %v", matches))
				result.Passed = false
			}
		}
	}

	// 检查IP地址
	if v.patterns["ip_address"] != nil {
		for _, text := range textFields {
			if matches := v.patterns["ip_address"].FindAllString(text, -1); len(matches) > 0 {
				sensitiveData = append(sensitiveData, fmt.Sprintf("发现IP地址: %v", matches))
				result.Passed = false
			}
		}
	}

	// 检查URL
	if v.patterns["url"] != nil {
		for _, text := range textFields {
			if matches := v.patterns["url"].FindAllString(text, -1); len(matches) > 0 {
				// 检查是否为内部或敏感URL
				for _, url := range matches {
					if strings.Contains(url, "localhost") ||
						strings.Contains(url, "127.0.0.1") ||
						strings.Contains(url, "internal") ||
						strings.Contains(url, "admin") {
						sensitiveData = append(sensitiveData, fmt.Sprintf("发现内部URL: %s", url))
						result.Passed = false
					}
				}
			}
		}
	}

	// 检查邮箱
	if v.patterns["email"] != nil {
		for _, text := range textFields {
			if matches := v.patterns["email"].FindAllString(text, -1); len(matches) > 0 {
				sensitiveData = append(sensitiveData, fmt.Sprintf("发现邮箱地址: %v", matches))
				result.Passed = false
			}
		}
	}

	if len(sensitiveData) > 0 {
		result.Message = "发现敏感数据"
		result.Details = sensitiveData
		result.Suggestions = []string{
			"移除或替换敏感数据",
			"使用占位符代替真实数据",
			"添加数据脱敏说明",
		}
	} else {
		result.Message = "未发现敏感数据"
	}

	return result
}

// checkStepOrdering 检查步骤顺序
func (v *Validator) checkStepOrdering(skill *Skill) ValidationResult {
	result := ValidationResult{
		RuleName: "step_ordering",
		Severity: "low",
		Passed:   true,
	}

	if len(skill.Steps) <= 1 {
		result.Message = "步骤顺序检查跳过（步骤数不足）"
		return result
	}

	// 检查步骤逻辑顺序
	// 这里可以添加更复杂的逻辑检查
	// 例如：文件创建应该在文件写入之前

	result.Message = "步骤顺序基本合理"
	return result
}

// checkConditionValidity 检查条件有效性
func (v *Validator) checkConditionValidity(skill *Skill) ValidationResult {
	result := ValidationResult{
		RuleName: "condition_validity",
		Severity: "medium",
		Passed:   true,
	}

	var invalidConditions []string

	for _, step := range skill.Steps {
		for _, condition := range step.Conditions {
			// 检查条件类型
			if condition.Type == "" {
				invalidConditions = append(invalidConditions,
					fmt.Sprintf("步骤 %d: 条件类型为空", step.Order))
				result.Passed = false
			}

			// 检查检查内容
			if condition.Check == "" {
				invalidConditions = append(invalidConditions,
					fmt.Sprintf("步骤 %d: 条件检查内容为空", step.Order))
				result.Passed = false
			}
		}
	}

	if len(invalidConditions) > 0 {
		result.Message = "发现无效条件"
		result.Details = invalidConditions
		result.Suggestions = []string{
			"为所有条件指定类型",
			"提供明确的检查内容",
			"确保条件逻辑合理",
		}
	} else {
		result.Message = "条件有效性检查通过"
	}

	return result
}

// checkExpectedResultClarity 检查预期结果清晰度
func (v *Validator) checkExpectedResultClarity(skill *Skill) ValidationResult {
	result := ValidationResult{
		RuleName: "expected_result_clarity",
		Severity: "low",
		Passed:   true,
	}

	var unclearResults []string

	for _, step := range skill.Steps {
		if step.Expected == "" {
			unclearResults = append(unclearResults,
				fmt.Sprintf("步骤 %d: 预期结果为空", step.Order))
			result.Passed = false
		} else if len(step.Expected) < 5 {
			unclearResults = append(unclearResults,
				fmt.Sprintf("步骤 %d: 预期结果过短", step.Order))
			result.Passed = false
		}
	}

	if len(unclearResults) > 0 {
		result.Message = "预期结果不够清晰"
		result.Details = unclearResults
		result.Suggestions = []string{
			"为每个步骤提供明确的预期结果",
			"预期结果应该具体、可验证",
			"避免使用模糊的描述",
		}
	} else {
		result.Message = "预期结果清晰度检查通过"
	}

	return result
}

// checkSkillComplexityMatch 检查技能复杂度匹配
func (v *Validator) checkSkillComplexityMatch(skill *Skill) ValidationResult {
	result := ValidationResult{
		RuleName: "skill_complexity_match",
		Severity: "low",
		Passed:   true,
	}

	// 根据步骤数量判断复杂度是否匹配
	stepCount := len(skill.Steps)

	expectedComplexity := "simple"
	if stepCount > 5 {
		expectedComplexity = "complex"
	} else if stepCount > 2 {
		expectedComplexity = "medium"
	}

	if skill.Complexity != expectedComplexity {
		result.Message = fmt.Sprintf("技能复杂度可能不匹配: 当前 %s, 建议 %s",
			skill.Complexity, expectedComplexity)
		result.Suggestions = []string{
			fmt.Sprintf("考虑将复杂度调整为 %s", expectedComplexity),
			"或重新评估步骤数量",
		}
	} else {
		result.Message = "技能复杂度匹配"
	}

	return result
}

// ValidationOptions 验证选项
type ValidationOptions struct {
	StrictMode         bool `json:"strict_mode" yaml:"strict_mode"`
	IgnoreWarnings     bool `json:"ignore_warnings" yaml:"ignore_warnings"`
	CheckSensitiveData bool `json:"check_sensitive_data" yaml:"check_sensitive_data"`
	CheckDangerousOps  bool `json:"check_dangerous_ops" yaml:"check_dangerous_ops"`
}

// DefaultValidationOptions 默认验证选项
func DefaultValidationOptions() ValidationOptions {
	return ValidationOptions{
		StrictMode:         false,
		IgnoreWarnings:     false,
		CheckSensitiveData: true,
		CheckDangerousOps:  true,
	}
}
