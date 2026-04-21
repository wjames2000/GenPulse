package skills

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Skill 表示一个技能，符合 agentskills.io 标准
type Skill struct {
	// 元数据
	ID          string    `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Version     string    `json:"version" yaml:"version"`
	Description string    `json:"description" yaml:"description"`
	Author      string    `json:"author" yaml:"author"`
	CreatedAt   time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" yaml:"updated_at"`

	// 分类和标签
	Category   string   `json:"category" yaml:"category"`
	Tags       []string `json:"tags" yaml:"tags"`
	Complexity string   `json:"complexity" yaml:"complexity"` // simple, medium, complex

	// 使用统计
	UsageCount  int       `json:"usage_count" yaml:"usage_count"`
	LastUsed    time.Time `json:"last_used" yaml:"last_used"`
	SuccessRate float64   `json:"success_rate" yaml:"success_rate"`

	// 技能内容
	Prerequisites []string `json:"prerequisites" yaml:"prerequisites"` // 前置技能
	Steps         []Step   `json:"steps" yaml:"steps"`                 // 执行步骤
	Examples      []string `json:"examples" yaml:"examples"`           // 使用示例
	Tips          []string `json:"tips" yaml:"tips"`                   // 使用技巧
	Warnings      []string `json:"warnings" yaml:"warnings"`           // 注意事项

	// 关联信息
	SourceTaskID string   `json:"source_task_id" yaml:"source_task_id"` // 来源任务ID
	RelatedTools []string `json:"related_tools" yaml:"related_tools"`   // 关联的工具
	AgentTypes   []string `json:"agent_types" yaml:"agent_types"`       // 适用的Agent类型

	// 性能指标
	AvgExecutionTime time.Duration `json:"avg_execution_time" yaml:"avg_execution_time"`
	TokenEstimate    int           `json:"token_estimate" yaml:"token_estimate"`

	// 状态
	Enabled   bool `json:"enabled" yaml:"enabled"`
	Validated bool `json:"validated" yaml:"validated"`
}

// Step 表示技能的一个执行步骤
type Step struct {
	ID         string      `json:"id" yaml:"id"`
	Order      int         `json:"order" yaml:"order"`
	Action     string      `json:"action" yaml:"action"`                         // 动作描述
	Tool       string      `json:"tool" yaml:"tool"`                             // 使用的工具
	Parameters []Parameter `json:"parameters" yaml:"parameters"`                 // 参数
	Conditions []Condition `json:"conditions" yaml:"conditions"`                 // 执行条件
	Expected   string      `json:"expected" yaml:"expected"`                     // 预期结果
	Fallback   *Step       `json:"fallback,omitempty" yaml:"fallback,omitempty"` // 备用步骤
}

// Parameter 表示步骤参数
type Parameter struct {
	Name        string `json:"name" yaml:"name"`
	Type        string `json:"type" yaml:"type"` // string, number, boolean, array, object
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"required" yaml:"required"`
	Default     any    `json:"default,omitempty" yaml:"default,omitempty"`
	Constraints string `json:"constraints,omitempty" yaml:"constraints,omitempty"`
}

// Condition 表示执行条件
type Condition struct {
	Type     string `json:"type" yaml:"type"`         // file_exists, env_set, previous_step_success, etc.
	Check    string `json:"check" yaml:"check"`       // 检查内容
	Expected any    `json:"expected" yaml:"expected"` // 期望值
}

// SkillMetadata 用于渐进式披露加载的元数据
type SkillMetadata struct {
	ID            string   `json:"id" yaml:"id"`
	Name          string   `json:"name" yaml:"name"`
	Version       string   `json:"version" yaml:"version"`
	Description   string   `json:"description" yaml:"description"`
	Category      string   `json:"category" yaml:"category"`
	Tags          []string `json:"tags" yaml:"tags"`
	Complexity    string   `json:"complexity" yaml:"complexity"`
	UsageCount    int      `json:"usage_count" yaml:"usage_count"`
	SuccessRate   float64  `json:"success_rate" yaml:"success_rate"`
	Enabled       bool     `json:"enabled" yaml:"enabled"`
	Validated     bool     `json:"validated" yaml:"validated"`
	TokenEstimate int      `json:"token_estimate" yaml:"token_estimate"`
}

// NewSkill 创建新技能
func NewSkill(name, description, author string) *Skill {
	now := time.Now()
	return &Skill{
		ID:            generateSkillID(name),
		Name:          name,
		Version:       "1.0.0",
		Description:   description,
		Author:        author,
		CreatedAt:     now,
		UpdatedAt:     now,
		Category:      "general",
		Tags:          []string{},
		Complexity:    "medium",
		UsageCount:    0,
		SuccessRate:   0.0,
		Steps:         []Step{},
		Examples:      []string{},
		Tips:          []string{},
		Warnings:      []string{},
		RelatedTools:  []string{},
		AgentTypes:    []string{},
		Enabled:       true,
		Validated:     false,
		TokenEstimate: 1000,
	}
}

// ToMetadata 转换为元数据
func (s *Skill) ToMetadata() *SkillMetadata {
	return &SkillMetadata{
		ID:            s.ID,
		Name:          s.Name,
		Version:       s.Version,
		Description:   s.Description,
		Category:      s.Category,
		Tags:          s.Tags,
		Complexity:    s.Complexity,
		UsageCount:    s.UsageCount,
		SuccessRate:   s.SuccessRate,
		Enabled:       s.Enabled,
		Validated:     s.Validated,
		TokenEstimate: s.TokenEstimate,
	}
}

// AddStep 添加步骤
func (s *Skill) AddStep(action, tool string, parameters []Parameter) {
	step := Step{
		ID:         generateStepID(s.ID, len(s.Steps)+1),
		Order:      len(s.Steps) + 1,
		Action:     action,
		Tool:       tool,
		Parameters: parameters,
		Conditions: []Condition{},
	}
	s.Steps = append(s.Steps, step)
	s.UpdatedAt = time.Now()
}

// AddCondition 为步骤添加条件
func (s *Skill) AddCondition(stepIndex int, conditionType, check string, expected any) error {
	if stepIndex < 0 || stepIndex >= len(s.Steps) {
		return ErrStepNotFound
	}

	condition := Condition{
		Type:     conditionType,
		Check:    check,
		Expected: expected,
	}

	s.Steps[stepIndex].Conditions = append(s.Steps[stepIndex].Conditions, condition)
	s.UpdatedAt = time.Now()
	return nil
}

// IncrementUsage 增加使用计数
func (s *Skill) IncrementUsage(success bool) {
	s.UsageCount++
	s.LastUsed = time.Now()

	// 更新成功率
	if success {
		if s.UsageCount == 1 {
			s.SuccessRate = 1.0
		} else {
			previousSuccess := s.SuccessRate * float64(s.UsageCount-1)
			s.SuccessRate = (previousSuccess + 1.0) / float64(s.UsageCount)
		}
	} else {
		if s.UsageCount == 1 {
			s.SuccessRate = 0.0
		} else {
			previousSuccess := s.SuccessRate * float64(s.UsageCount-1)
			s.SuccessRate = previousSuccess / float64(s.UsageCount)
		}
	}
}

// generateSkillID 生成技能ID
func generateSkillID(name string) string {
	// 简单实现：使用时间戳+名称哈希
	timestamp := time.Now().Unix()
	return formatSkillID(name, timestamp)
}

// generateStepID 生成步骤ID
func generateStepID(skillID string, stepNumber int) string {
	return skillID + "_step_" + string(rune(stepNumber))
}

// formatSkillID 格式化技能ID
func formatSkillID(name string, timestamp int64) string {
	// 移除空格和特殊字符，转换为小写
	formatted := ""
	for _, ch := range name {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			formatted += string(ch)
		} else if ch == ' ' || ch == '-' || ch == '_' {
			formatted += "_"
		}
	}

	// 转换为小写
	formatted = strings.ToLower(formatted)

	// 添加时间戳
	return formatted + "_" + strconv.FormatInt(timestamp, 10)
}

// 错误定义
var (
	ErrStepNotFound = errors.New("step not found")
)
