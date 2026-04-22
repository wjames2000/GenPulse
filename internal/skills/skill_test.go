package skills

import (
	"testing"
	"time"
)

// TestSkillBasic 测试技能基础功能
func TestSkillBasic(t *testing.T) {
	// 测试1: 创建技能
	t.Run("创建技能", func(t *testing.T) {
		skill := NewSkill("测试技能", "这是一个测试技能", "测试作者")

		if skill.Name != "测试技能" {
			t.Errorf("技能名称不正确: 期望 %s, 实际 %s", "测试技能", skill.Name)
		}
		if skill.Description != "这是一个测试技能" {
			t.Errorf("技能描述不正确: 期望 %s, 实际 %s", "这是一个测试技能", skill.Description)
		}
		if skill.Author != "测试作者" {
			t.Errorf("技能作者不正确: 期望 %s, 实际 %s", "测试作者", skill.Author)
		}
		if skill.ID == "" {
			t.Error("技能ID应该自动生成")
		}
		if skill.CreatedAt.IsZero() {
			t.Error("创建时间应该被设置")
		}
		if skill.UpdatedAt.IsZero() {
			t.Error("更新时间应该被设置")
		}
		if !skill.Enabled {
			t.Error("新技能应该被启用")
		}
	})

	// 测试2: 添加步骤
	t.Run("添加步骤", func(t *testing.T) {
		skill := NewSkill("测试技能", "这是一个测试技能", "测试作者")

		// 添加步骤
		skill.AddStep("第一步", "test-tool", []Parameter{
			{Name: "param1", Type: "string", Required: true},
		})

		if len(skill.Steps) != 1 {
			t.Errorf("步骤数量不正确: 期望 1, 实际 %d", len(skill.Steps))
		}

		step := skill.Steps[0]
		if step.Action != "第一步" {
			t.Errorf("步骤动作不正确: 期望 %s, 实际 %s", "第一步", step.Action)
		}
		if step.Tool != "test-tool" {
			t.Errorf("步骤工具不正确: 期望 %s, 实际 %s", "test-tool", step.Tool)
		}
		if step.Order != 1 {
			t.Errorf("步骤顺序不正确: 期望 %d, 实际 %d", 1, step.Order)
		}
		if len(step.Parameters) != 1 {
			t.Errorf("步骤参数数量不正确: 期望 1, 实际 %d", len(step.Parameters))
		}
	})

	// 测试3: 添加条件
	t.Run("添加条件", func(t *testing.T) {
		skill := NewSkill("测试技能", "这是一个测试技能", "测试作者")

		// 添加步骤
		skill.AddStep("第一步", "test-tool", []Parameter{
			{Name: "param1", Type: "string", Required: true},
		})

		// 添加条件
		err := skill.AddCondition(0, "equals", "param1", "value1")
		if err != nil {
			t.Errorf("添加条件失败: %v", err)
		}

		if len(skill.Steps[0].Conditions) != 1 {
			t.Errorf("条件数量不正确: 期望 1, 实际 %d", len(skill.Steps[0].Conditions))
		}

		condition := skill.Steps[0].Conditions[0]
		if condition.Type != "equals" {
			t.Errorf("条件类型不正确: 期望 %s, 实际 %s", "equals", condition.Type)
		}
		if condition.Check != "param1" {
			t.Errorf("条件检查字段不正确: 期望 %s, 实际 %s", "param1", condition.Check)
		}
	})

	// 测试4: 使用统计
	t.Run("使用统计", func(t *testing.T) {
		skill := NewSkill("测试技能", "这是一个测试技能", "测试作者")

		// 初始统计
		if skill.UsageCount != 0 {
			t.Errorf("初始使用次数不正确: 期望 0, 实际 %d", skill.UsageCount)
		}
		if skill.SuccessRate != 0 {
			t.Errorf("初始成功率不正确: 期望 0, 实际 %f", skill.SuccessRate)
		}

		// 增加使用次数（成功）
		skill.IncrementUsage(true)
		if skill.UsageCount != 1 {
			t.Errorf("使用次数不正确: 期望 1, 实际 %d", skill.UsageCount)
		}
		if skill.SuccessRate != 1.0 {
			t.Errorf("成功率不正确: 期望 1.0, 实际 %f", skill.SuccessRate)
		}

		// 增加使用次数（失败）
		skill.IncrementUsage(false)
		if skill.UsageCount != 2 {
			t.Errorf("使用次数不正确: 期望 2, 实际 %d", skill.UsageCount)
		}
		if skill.SuccessRate != 0.5 {
			t.Errorf("成功率不正确: 期望 0.5, 实际 %f", skill.SuccessRate)
		}

		// 检查最后使用时间
		if skill.LastUsed.IsZero() {
			t.Error("最后使用时间应该被更新")
		}
	})

	// 测试5: 转换为元数据
	t.Run("转换为元数据", func(t *testing.T) {
		skill := NewSkill("测试技能", "这是一个测试技能", "测试作者")
		skill.Category = "测试"
		skill.Tags = []string{"tag1", "tag2"}
		skill.Complexity = "simple"

		metadata := skill.ToMetadata()

		if metadata.ID != skill.ID {
			t.Errorf("元数据ID不正确: 期望 %s, 实际 %s", skill.ID, metadata.ID)
		}
		if metadata.Name != skill.Name {
			t.Errorf("元数据名称不正确: 期望 %s, 实际 %s", skill.Name, metadata.Name)
		}
		if metadata.Description != skill.Description {
			t.Errorf("元数据描述不正确: 期望 %s, 实际 %s", skill.Description, metadata.Description)
		}
		if metadata.Category != skill.Category {
			t.Errorf("元数据类别不正确: 期望 %s, 实际 %s", skill.Category, metadata.Category)
		}
		if len(metadata.Tags) != len(skill.Tags) {
			t.Errorf("元数据标签数量不正确: 期望 %d, 实际 %d", len(skill.Tags), len(metadata.Tags))
		}
		if metadata.Complexity != skill.Complexity {
			t.Errorf("元数据复杂度不正确: 期望 %s, 实际 %s", skill.Complexity, metadata.Complexity)
		}
	})

	// 测试6: 技能ID生成
	t.Run("技能ID生成", func(t *testing.T) {
		skill := NewSkill("测试技能", "这是一个测试技能", "测试作者")

		// ID应该不为空
		if len(skill.ID) == 0 {
			t.Error("技能ID不应该为空")
		}

		// ID应该包含技能名称（小写）
		if len(skill.ID) < 4 {
			t.Errorf("技能ID长度不合理: %s", skill.ID)
		}
	})

	// 测试7: 技能字段设置
	t.Run("技能字段设置", func(t *testing.T) {
		skill := NewSkill("测试技能", "这是一个测试技能", "测试作者")

		// 设置分类和标签
		skill.Category = "测试分类"
		skill.Tags = []string{"测试", "示例"}
		skill.Complexity = "medium"
		skill.SourceTaskID = "task-123"
		skill.RelatedTools = []string{"tool1", "tool2"}
		skill.AgentTypes = []string{"agent1", "agent2"}
		skill.AvgExecutionTime = 5 * time.Second
		skill.TokenEstimate = 1000

		// 验证字段设置
		if skill.Category != "测试分类" {
			t.Errorf("分类不正确: 期望 %s, 实际 %s", "测试分类", skill.Category)
		}
		if len(skill.Tags) != 2 {
			t.Errorf("标签数量不正确: 期望 2, 实际 %d", len(skill.Tags))
		}
		if skill.Complexity != "medium" {
			t.Errorf("复杂度不正确: 期望 %s, 实际 %s", "medium", skill.Complexity)
		}
		if skill.SourceTaskID != "task-123" {
			t.Errorf("来源任务ID不正确: 期望 %s, 实际 %s", "task-123", skill.SourceTaskID)
		}
		if len(skill.RelatedTools) != 2 {
			t.Errorf("关联工具数量不正确: 期望 2, 实际 %d", len(skill.RelatedTools))
		}
		if len(skill.AgentTypes) != 2 {
			t.Errorf("Agent类型数量不正确: 期望 2, 实际 %d", len(skill.AgentTypes))
		}
		if skill.AvgExecutionTime != 5*time.Second {
			t.Errorf("平均执行时间不正确: 期望 %v, 实际 %v", 5*time.Second, skill.AvgExecutionTime)
		}
		if skill.TokenEstimate != 1000 {
			t.Errorf("Token估计不正确: 期望 %d, 实际 %d", 1000, skill.TokenEstimate)
		}
	})
}

// TestSkillEdgeCases 测试技能边界情况
func TestSkillEdgeCases(t *testing.T) {
	// 测试1: 空技能名称
	t.Run("空技能名称", func(t *testing.T) {
		skill := NewSkill("", "描述", "作者")
		if skill.Name != "" {
			t.Error("空名称应该被允许")
		}
		if skill.ID == "" {
			t.Error("即使名称为空，ID也应该生成")
		}
	})

	// 测试2: 添加条件到不存在的步骤
	t.Run("添加条件到不存在的步骤", func(t *testing.T) {
		skill := NewSkill("测试技能", "这是一个测试技能", "测试作者")

		err := skill.AddCondition(0, "equals", "param1", "value1")
		if err == nil {
			t.Error("添加到不存在的步骤应该失败")
		}
	})

	// 测试3: 多次添加步骤
	t.Run("多次添加步骤", func(t *testing.T) {
		skill := NewSkill("测试技能", "这是一个测试技能", "测试作者")

		// 添加多个步骤
		for i := 0; i < 5; i++ {
			skill.AddStep("步骤", "tool", []Parameter{})
		}

		if len(skill.Steps) != 5 {
			t.Errorf("步骤数量不正确: 期望 5, 实际 %d", len(skill.Steps))
		}

		// 验证步骤顺序
		for i, step := range skill.Steps {
			if step.Order != i+1 {
				t.Errorf("步骤 %d 的顺序不正确: 期望 %d, 实际 %d", i, i+1, step.Order)
			}
		}
	})

	// 测试4: 技能禁用
	t.Run("技能禁用", func(t *testing.T) {
		skill := NewSkill("测试技能", "这是一个测试技能", "测试作者")

		skill.Enabled = false
		if skill.Enabled {
			t.Error("技能应该被禁用")
		}

		skill.Enabled = true
		if !skill.Enabled {
			t.Error("技能应该被启用")
		}
	})

	// 测试5: 技能验证状态
	t.Run("技能验证状态", func(t *testing.T) {
		skill := NewSkill("测试技能", "这是一个测试技能", "测试作者")

		skill.Validated = true
		if !skill.Validated {
			t.Error("技能应该被标记为已验证")
		}

		skill.Validated = false
		if skill.Validated {
			t.Error("技能应该被标记为未验证")
		}
	})
}
