package agents

import (
	"context"
	"testing"
	"time"

	"GenPulse/internal/genkit/models"
)

// TestAgentBasic 测试Agent基础功能
func TestAgentBasic(t *testing.T) {
	// 创建Agent配置
	config := AgentConfig{
		ID:          "test-agent-1",
		Name:        "测试Agent",
		Role:        RoleFullStackDev,
		Description: "测试用的Agent",
		ModelConfig: models.ModelConfig{
			Name:     "test-model",
			Provider: "test",
			Type:     "test",
		},
		Capabilities: []AgentCapability{
			CapabilityCodeGeneration,
		},
		Tools:      []string{},
		MaxRetries: 3,
		Timeout:    5 * time.Minute,
		Enabled:    true,
	}

	// 创建BaseAgent（使用nil依赖，因为只是测试基础功能）
	agent := &BaseAgent{
		config:  config,
		state:   StateIdle,
		enabled: true,
	}

	// 测试1: 验证Agent基本信息
	t.Run("Agent基本信息", func(t *testing.T) {
		agentConfig := agent.GetConfig()
		if agentConfig.ID != "test-agent-1" {
			t.Errorf("Agent ID不正确: 期望 %s, 实际 %s", "test-agent-1", agentConfig.ID)
		}
		if agentConfig.Name != "测试Agent" {
			t.Errorf("Agent名称不正确: 期望 %s, 实际 %s", "测试Agent", agentConfig.Name)
		}
		if agentConfig.Role != RoleFullStackDev {
			t.Errorf("Agent角色不正确: 期望 %s, 实际 %s", RoleFullStackDev, agentConfig.Role)
		}
	})

	// 测试2: 验证Agent状态
	t.Run("Agent状态管理", func(t *testing.T) {
		initialState := agent.GetState()
		if initialState != StateIdle {
			t.Errorf("初始状态不正确: 期望 %s, 实际 %s", StateIdle, initialState)
		}

		if !agent.IsEnabled() {
			t.Error("Agent应该处于启用状态")
		}

		// 测试禁用/启用
		agent.SetEnabled(false)
		if agent.IsEnabled() {
			t.Error("Agent应该处于禁用状态")
		}
		agent.SetEnabled(true)
		if !agent.IsEnabled() {
			t.Error("Agent应该处于启用状态")
		}
	})

	// 测试3: 验证任务验证
	t.Run("任务验证", func(t *testing.T) {
		// 测试空任务
		if err := agent.ValidateTask(""); err == nil {
			t.Error("空任务应该验证失败")
		}

		// 测试过短任务
		if err := agent.ValidateTask("ab"); err == nil {
			t.Error("过短任务应该验证失败")
		}

		// 测试有效任务
		validTask := "创建一个简单的Go Web服务器"
		if err := agent.ValidateTask(validTask); err != nil {
			t.Errorf("有效任务验证失败: %v", err)
		}
	})

	// 测试4: 验证Agent能力
	t.Run("Agent能力检查", func(t *testing.T) {
		// BaseAgent没有HasCapability方法，跳过这个测试
		// 如果需要测试能力检查，需要在BaseAgent中添加HasCapability方法
	})

	// 测试5: 执行简单任务（模拟）
	t.Run("执行简单任务", func(t *testing.T) {
		ctx := context.Background()
		task := "创建一个简单的Go Web服务器，监听8080端口"

		result, err := agent.Execute(ctx, task, map[string]interface{}{
			"port": 8080,
		})

		if err != nil {
			t.Errorf("执行任务失败: %v", err)
		}

		if result == nil {
			t.Error("执行结果不应该为nil")
		}

		if !result.Success {
			t.Error("任务执行应该成功")
		}

		// 验证执行上下文
		execution := agent.GetExecution()
		if execution == nil {
			t.Error("执行上下文不应该为nil")
		}

		if execution.Task != task {
			t.Errorf("执行任务不匹配: 期望 %s, 实际 %s", task, execution.Task)
		}

		if execution.State != StateCompleted {
			t.Errorf("执行状态不正确: 期望 %s, 实际 %s", StateCompleted, execution.State)
		}
	})

	// 测试6: 统计信息
	t.Run("统计信息", func(t *testing.T) {
		executionCount := agent.GetExecutionCount()
		if executionCount < 1 {
			t.Error("Agent应该至少执行了一次任务")
		}

		successRate := agent.GetSuccessRate()
		if successRate < 0 {
			t.Error("成功率不应该为负数")
		}

		avgDuration := agent.GetAverageDuration()
		if avgDuration < 0 {
			t.Error("平均执行时间不应该为负数")
		}
	})
}

// TestAgentEndToEnd 测试Agent端到端功能
func TestAgentEndToEnd(t *testing.T) {
	t.Run("端到端验证", func(t *testing.T) {
		// 创建Agent配置
		config := AgentConfig{
			ID:          "e2e-test-agent",
			Name:        "端到端测试Agent",
			Role:        RoleFullStackDev,
			Description: "用于端到端测试的Agent",
			ModelConfig: models.ModelConfig{
				Name:     "test-model",
				Provider: "test",
				Type:     "test",
			},
			Capabilities: []AgentCapability{
				CapabilityCodeGeneration,
				CapabilityProjectSetup,
			},
			Enabled: true,
		}

		// 创建BaseAgent
		agent := &BaseAgent{
			config:  config,
			state:   StateIdle,
			enabled: true,
		}

		// 执行端到端任务
		ctx := context.Background()
		task := "创建一个简单的Go项目，包含main.go文件和基本的HTTP服务器"

		result, err := agent.Execute(ctx, task, map[string]interface{}{
			"project_name": "test-project",
			"port":         8080,
		})

		if err != nil {
			t.Errorf("端到端任务执行失败: %v", err)
		}

		if result == nil {
			t.Error("端到端任务执行结果不应该为nil")
		}

		// 验证执行统计
		executionCount := agent.GetExecutionCount()
		if executionCount < 1 {
			t.Error("Agent应该至少执行了一次任务")
		}

		successRate := agent.GetSuccessRate()
		if successRate < 0 {
			t.Error("成功率不应该为负数")
		}

		// 验证任务执行成功
		if !result.Success {
			t.Error("端到端任务应该执行成功")
		}

		// 验证输出
		if result.Output == nil {
			t.Error("执行结果应该有输出")
		}
	})
}
