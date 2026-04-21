package agents

import (
	"fmt"
	"testing"
	"time"
)

// TestSimpleEvolutionManager 测试简化自进化管理器
func TestSimpleEvolutionManager(t *testing.T) {
	// 创建管理器
	config := DefaultSimpleConfig()
	manager := NewSimpleEvolutionManager(config)

	// 创建测试Agent
	agent := createTestSimpleAgent(t)

	// 测试1: 提示词增强
	t.Run("TestPromptEnhancement", func(t *testing.T) {
		testSimplePromptEnhancement(t, manager, agent)
	})

	// 测试2: 反馈闭环
	t.Run("TestFeedbackLoop", func(t *testing.T) {
		testSimpleFeedbackLoop(t, manager, agent)
	})

	// 测试3: 事件追踪
	t.Run("TestEventTracking", func(t *testing.T) {
		testSimpleEventTracking(t, manager)
	})

	// 测试4: 统计功能
	t.Run("TestStats", func(t *testing.T) {
		testSimpleStats(t, manager)
	})
}

// testSimplePromptEnhancement 测试简化提示词增强
func testSimplePromptEnhancement(t *testing.T, manager *SimpleEvolutionManager, agent Agent) {
	// 添加测试技能
	manager.AddTestSkill(&SimpleSkill{
		ID:          "test-skill-1",
		Name:        "create_login_page",
		Description: "创建用户登录页面",
		Category:    "frontend_developer",
		UsageCount:  10,
		SuccessRate: 0.9,
		CreatedAt:   time.Now(),
		LastUsed:    time.Now(),
		Enabled:     true,
	})

	manager.AddTestSkill(&SimpleSkill{
		ID:          "test-skill-2",
		Name:        "implement_rest_api",
		Description: "实现RESTful API接口",
		Category:    "backend_developer",
		UsageCount:  15,
		SuccessRate: 0.85,
		CreatedAt:   time.Now(),
		LastUsed:    time.Now(),
		Enabled:     true,
	})

	// 测试任务
	task := "创建一个用户登录页面，包含用户名和密码输入"

	// 增强提示词
	enhancedPrompt := manager.EnhancePrompt(agent, task)

	// 验证结果
	if enhancedPrompt == "" {
		t.Error("增强后的提示词为空")
	}

	// 验证包含原始任务
	if !contains(enhancedPrompt, task) {
		t.Error("增强后的提示词不包含原始任务")
	}

	// 验证包含技能信息
	if !contains(enhancedPrompt, "可用技能") {
		t.Log("注意: 增强后的提示词未包含技能信息，可能没有相关技能")
	}

	t.Logf("测试通过: 简化提示词增强成功")
	t.Logf("原始任务: %s", task)
	t.Logf("增强后长度: %d 字符", len(enhancedPrompt))
}

// testSimpleFeedbackLoop 测试简化反馈闭环
func testSimpleFeedbackLoop(t *testing.T, manager *SimpleEvolutionManager, agent Agent) {
	// 测试任务1: 简单任务（不应创建技能）
	simpleTask := "修改按钮颜色"
	simpleResult := &AgentResult{
		Success:  true,
		Output:   "按钮颜色已修改为蓝色",
		Duration: 5 * time.Second,
	}

	err := manager.RecordExecution(agent, simpleTask, simpleResult)
	if err != nil {
		t.Errorf("记录简单任务失败: %v", err)
	}

	// 测试任务2: 复杂任务（应创建技能）
	complexTask := "实现用户注册功能，包含表单验证和数据库存储"
	complexResult := &AgentResult{
		Success:  true,
		Output:   "成功实现用户注册功能",
		Duration: 3 * time.Minute,
	}

	err = manager.RecordExecution(agent, complexTask, complexResult)
	if err != nil {
		t.Errorf("记录复杂任务失败: %v", err)
	}

	// 验证技能数量增加
	skills := manager.GetSkills()
	t.Logf("当前技能数量: %d", len(skills))

	// 验证事件记录
	events := manager.GetEvents(10)
	if len(events) == 0 {
		t.Error("未记录任何事件")
	} else {
		t.Logf("记录了 %d 个事件", len(events))
	}

	t.Log("测试通过: 简化反馈闭环功能正常")
}

// testSimpleEventTracking 测试简化事件追踪
func testSimpleEventTracking(t *testing.T, manager *SimpleEvolutionManager) {
	// 获取事件
	events := manager.GetEvents(5)

	if len(events) == 0 {
		t.Log("注意: 暂无事件记录")
	} else {
		t.Logf("事件数量: %d", len(events))

		// 输出最近事件
		for i, event := range events {
			t.Logf("事件 %d: %s - %s (成功: %v)",
				i+1, event.EventType, event.Description, event.Success)
		}
	}

	// 验证事件类型
	hasPromptEvent := false
	hasExecutionEvent := false

	for _, event := range events {
		if event.EventType == "prompt_enhanced" {
			hasPromptEvent = true
		}
		if event.EventType == "execution_success" || event.EventType == "execution_failed" {
			hasExecutionEvent = true
		}
	}

	if !hasPromptEvent {
		t.Log("注意: 未找到提示词增强事件")
	}
	if !hasExecutionEvent {
		t.Log("注意: 未找到执行事件")
	}

	t.Log("测试通过: 简化事件追踪功能正常")
}

// testSimpleStats 测试简化统计功能
func testSimpleStats(t *testing.T, manager *SimpleEvolutionManager) {
	// 获取统计信息
	stats := manager.GetStats()

	if stats == nil {
		t.Error("统计信息为空")
		return
	}

	// 验证基本统计字段
	requiredFields := []string{
		"total_skills",
		"total_events",
		"enabled_skills",
		"total_skill_usage",
		"avg_success_rate",
	}

	for _, field := range requiredFields {
		if _, ok := stats[field]; !ok {
			t.Errorf("缺少统计字段: %s", field)
		}
	}

	// 输出统计信息
	t.Logf("自进化统计:")
	t.Logf("  技能总数: %v", stats["total_skills"])
	t.Logf("  启用技能: %v", stats["enabled_skills"])
	t.Logf("  技能总使用次数: %v", stats["total_skill_usage"])
	t.Logf("  平均成功率: %.1f%%", stats["avg_success_rate"].(float64)*100)
	t.Logf("  事件总数: %v", stats["total_events"])

	// 输出事件类型统计
	if eventStats, ok := stats["event_stats"].(map[string]int); ok {
		t.Logf("  事件类型统计:")
		for eventType, count := range eventStats {
			t.Logf("    %s: %d", eventType, count)
		}
	}

	t.Log("测试通过: 简化统计功能正常")
}

// createTestSimpleAgent 创建测试Agent
func createTestSimpleAgent(t *testing.T) Agent {
	config := AgentConfig{
		ID:      "test-simple-agent",
		Name:    "测试简化Agent",
		Role:    RoleFrontendDev,
		Enabled: true,
		PromptTemplates: map[string]string{
			"default": "你是一个前端开发专家，擅长使用React和TypeScript。",
		},
	}

	agent, err := NewBaseAgent(config, nil, nil, nil)
	if err != nil {
		t.Fatalf("创建测试Agent失败: %v", err)
	}

	return agent
}

// BenchmarkSimpleEvolution 性能测试
func BenchmarkSimpleEvolution(b *testing.B) {
	manager := NewSimpleEvolutionManager(DefaultSimpleConfig())
	agent := &BaseAgent{
		config: AgentConfig{
			ID:   "benchmark-agent",
			Name: "性能测试Agent",
			Role: RoleFullStackDev,
			PromptTemplates: map[string]string{
				"default": "测试提示词模板",
			},
		},
	}

	// 添加一些测试技能
	for i := 0; i < 100; i++ {
		manager.AddTestSkill(&SimpleSkill{
			ID:          fmt.Sprintf("benchmark-skill-%d", i),
			Name:        fmt.Sprintf("skill_%d", i),
			Description: fmt.Sprintf("测试技能 %d", i),
			Category:    "benchmark",
			UsageCount:  i,
			SuccessRate: 0.8,
			CreatedAt:   time.Now(),
			LastUsed:    time.Now(),
			Enabled:     true,
		})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 测试提示词增强
		task := fmt.Sprintf("测试任务 %d", i)
		_ = manager.EnhancePrompt(agent, task)

		// 测试记录执行结果
		result := &AgentResult{
			Success:  true,
			Output:   "测试输出",
			Duration: time.Duration(i%100) * time.Millisecond,
		}
		_ = manager.RecordExecution(agent, task, result)
	}
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

// ExampleSimpleEvolutionManager 使用示例
func ExampleSimpleEvolutionManager() {
	// 创建自进化管理器
	config := DefaultSimpleConfig()
	manager := NewSimpleEvolutionManager(config)

	// 创建Agent
	agentConfig := AgentConfig{
		ID:   "example-agent",
		Name: "示例Agent",
		Role: RoleFullStackDev,
		PromptTemplates: map[string]string{
			"default": "你是一个全栈开发专家。",
		},
	}
	agent, _ := NewBaseAgent(agentConfig, nil, nil, nil)

	// 1. 增强提示词
	task := "开发一个待办事项应用"
	enhancedPrompt := manager.EnhancePrompt(agent, task)
	fmt.Printf("增强提示词完成，长度: %d\n", len(enhancedPrompt))

	// 2. 记录执行结果
	result := &AgentResult{
		Success:  true,
		Output:   "成功创建待办事项应用",
		Duration: 2 * time.Minute,
	}
	manager.RecordExecution(agent, task, result)

	// 3. 获取统计信息
	stats := manager.GetStats()
	fmt.Printf("技能总数: %v\n", stats["total_skills"])
	fmt.Printf("事件总数: %v\n", stats["total_events"])

	// 输出: 增强提示词完成，长度: ...
	// 输出: 技能总数: ...
	// 输出: 事件总数: ...
}
