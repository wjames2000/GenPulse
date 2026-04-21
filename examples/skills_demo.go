package main

import (
	"context"
	"fmt"
	"time"

	"GenPulse/internal/skills"
)

func main() {
	fmt.Println("=== Skills 闭环引擎演示 ===")

	// 创建模拟的LLM客户端
	llmClient := &MockLLMClient{}

	// 创建技能管理器
	manager, err := skills.NewSkillManager("./skills_demo", llmClient)
	if err != nil {
		fmt.Printf("创建技能管理器失败: %v\n", err)
		return
	}

	fmt.Println("1. 创建模拟任务执行记录...")

	// 创建模拟任务执行记录
	taskRecord := &skills.TaskExecutionRecord{
		TaskID:      "task_001",
		TaskType:    "file_operation",
		Description: "创建Go项目结构并初始化git仓库",
		Complexity:  "medium",
		StartTime:   time.Now().Add(-10 * time.Minute),
		EndTime:     time.Now(),
		Success:     true,
		ToolUsage: map[string]int{
			"fs_create": 3,
			"fs_write":  2,
			"git_init":  1,
			"git_add":   1,
		},
		AgentInvolved: []string{"BackendDevAgent", "DevOpsAgent"},
		Steps: []skills.TaskStep{
			{
				StepID:     "step_001",
				Order:      1,
				Action:     "创建项目目录结构",
				Tool:       "fs_create",
				Parameters: map[string]any{"path": "./myproject", "type": "directory"},
				StartTime:  time.Now().Add(-10 * time.Minute),
				EndTime:    time.Now().Add(-9 * time.Minute),
				Success:    true,
				Output:     map[string]any{"created": true, "path": "./myproject"},
			},
			{
				StepID:     "step_002",
				Order:      2,
				Action:     "创建go.mod文件",
				Tool:       "fs_write",
				Parameters: map[string]any{"path": "./myproject/go.mod", "content": "module myproject\n\ngo 1.21\n"},
				StartTime:  time.Now().Add(-9 * time.Minute),
				EndTime:    time.Now().Add(-8 * time.Minute),
				Success:    true,
				Output:     map[string]any{"written": true, "path": "./myproject/go.mod"},
			},
			{
				StepID:     "step_003",
				Order:      3,
				Action:     "初始化git仓库",
				Tool:       "git_init",
				Parameters: map[string]any{"path": "./myproject"},
				StartTime:  time.Now().Add(-8 * time.Minute),
				EndTime:    time.Now().Add(-7 * time.Minute),
				Success:    true,
				Output:     map[string]any{"initialized": true, "path": "./myproject/.git"},
			},
		},
		Context: map[string]any{
			"project_type": "go",
			"requirements": "简单的Go项目骨架",
		},
		Output: map[string]any{
			"project_created": true,
			"project_path":    "./myproject",
			"has_git":         true,
		},
	}

	fmt.Println("2. 处理任务执行记录...")

	// 处理任务执行记录
	result, err := manager.ProcessTaskExecution(taskRecord)
	if err != nil {
		fmt.Printf("处理任务执行记录失败: %v\n", err)
		return
	}

	if result.Triggered {
		fmt.Printf("技能生成触发: %v\n", result.TriggerResults)

		if result.Success {
			fmt.Printf("技能生成成功: %s\n", result.Message)
			fmt.Printf("生成的技能: %s (ID: %s)\n",
				result.GeneratedSkill.Name, result.GeneratedSkill.ID)

			// 显示验证报告
			if result.ValidationReport != nil {
				fmt.Printf("验证结果: %v\n", result.ValidationReport.OverallPass)
				fmt.Printf("通过检查: %d/%d\n",
					result.ValidationReport.PassedChecks,
					result.ValidationReport.TotalChecks)
			}
		} else {
			fmt.Printf("技能生成失败: %s\n", result.Error)
		}
	} else {
		fmt.Println("技能生成未触发")
	}

	fmt.Println("\n3. 列出所有技能...")

	// 列出所有技能
	skillList, err := manager.ListSkills(nil)
	if err != nil {
		fmt.Printf("列出技能失败: %v\n", err)
		return
	}

	fmt.Printf("找到 %d 个技能:\n", len(skillList))
	for i, skill := range skillList {
		fmt.Printf("%d. %s (分类: %s, 使用次数: %d)\n",
			i+1, skill.Name, skill.Category, skill.UsageCount)
	}

	fmt.Println("\n4. 获取技能统计...")

	// 获取技能统计
	stats, err := manager.GetSkillStats()
	if err != nil {
		fmt.Printf("获取技能统计失败: %v\n", err)
		return
	}

	fmt.Printf("总技能数: %d\n", stats.TotalSkills)
	fmt.Printf("已启用: %d\n", stats.EnabledSkills)
	fmt.Printf("已验证: %d\n", stats.ValidatedSkills)
	fmt.Printf("总使用次数: %d\n", stats.TotalUsage)
	fmt.Printf("平均成功率: %.1f%%\n", stats.AverageSuccessRate*100)

	// 按分类显示
	fmt.Println("\n按分类分布:")
	for category, count := range stats.ByCategory {
		fmt.Printf("  %s: %d\n", category, count)
	}

	fmt.Println("\n=== 演示完成 ===")
}

// MockLLMClient 模拟LLM客户端
type MockLLMClient struct{}

func (m *MockLLMClient) Generate(ctx context.Context, prompt string, options map[string]any) (string, error) {
	// 模拟LLM响应
	fmt.Println("[MockLLM] 收到提示词:", prompt[:min(100, len(prompt))], "...")

	// 根据提示词类型返回模拟响应
	if contains(prompt, "技能名称") {
		return "Go项目初始化技能|创建Go项目目录结构并初始化git仓库", nil
	} else if contains(prompt, "优化步骤") {
		return "步骤已优化", nil
	} else if contains(prompt, "使用示例") {
		return "示例1: 初始化一个新的Go Web项目\n示例2: 创建Go库项目结构", nil
	}

	return "模拟响应", nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
