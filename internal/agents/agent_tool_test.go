package agents

import (
	"context"
	"testing"
	"time"

	"GenPulse/internal/genkit/models"
	"GenPulse/internal/genkit/tools"
)

// MockAgent 模拟Agent用于测试
type MockAgent struct {
	config  AgentConfig
	state   AgentState
	enabled bool
}

func NewMockAgent(config AgentConfig) *MockAgent {
	return &MockAgent{
		config:  config,
		state:   StateIdle,
		enabled: true,
	}
}

func (m *MockAgent) GetConfig() AgentConfig {
	return m.config
}

func (m *MockAgent) GetState() AgentState {
	return m.state
}

func (m *MockAgent) GetExecution() *AgentExecution {
	return nil
}

func (m *MockAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	if !m.enabled {
		return nil, NewAgentError("AGENT_DISABLED", "Agent已禁用", "")
	}

	m.state = StateExecuting
	defer func() { m.state = StateIdle }()

	// 模拟执行
	time.Sleep(10 * time.Millisecond)

	return &AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"task":       task,
			"parameters": parameters,
			"message":    "任务执行成功",
			"agent_id":   m.config.ID,
			"agent_name": m.config.Name,
		},
		Duration: 10 * time.Millisecond,
		Metadata: map[string]interface{}{
			"agent_id":   m.config.ID,
			"agent_name": m.config.Name,
			"timestamp":  time.Now(),
		},
	}, nil
}

func (m *MockAgent) Cancel() error {
	m.state = StateFailed
	return nil
}

func (m *MockAgent) ValidateTask(task string) error {
	if task == "" {
		return NewAgentError("INVALID_TASK", "任务验证失败", "任务不能为空")
	}
	return nil
}

func (m *MockAgent) Initialize() error {
	return nil
}

func (m *MockAgent) Shutdown() error {
	return nil
}

func (m *MockAgent) IsEnabled() bool {
	return m.enabled
}

func (m *MockAgent) SetEnabled(enabled bool) {
	m.enabled = enabled
}

func (m *MockAgent) GetExecutionCount() int {
	return 0
}

func (m *MockAgent) GetSuccessRate() float64 {
	return 1.0
}

func (m *MockAgent) GetAverageDuration() time.Duration {
	return 10 * time.Millisecond
}

// TestAgentToolBasic 测试AgentTool基础功能
func TestAgentToolBasic(t *testing.T) {
	// 创建模拟Agent
	agentConfig := AgentConfig{
		ID:          "test-agent-001",
		Name:        "测试Agent",
		Role:        RoleFullStackDev,
		Description: "测试用的Agent",
		ModelConfig: models.ModelConfig{
			Name:     "test-model",
			Provider: "test",
			Type:     "test",
		},
		Enabled: true,
	}

	mockAgent := NewMockAgent(agentConfig)

	// 创建AgentTool
	agentTool, err := NewAgentTool(mockAgent)
	if err != nil {
		t.Fatalf("创建AgentTool失败: %v", err)
	}

	// 测试1: 验证工具定义
	t.Run("工具定义验证", func(t *testing.T) {
		def := agentTool.GetDefinition()
		if def.ID != "agent_test-agent-001" {
			t.Errorf("工具ID不正确: 期望 %s, 实际 %s", "agent_test-agent-001", def.ID)
		}
		if def.Name != "测试Agent Agent" {
			t.Errorf("工具名称不正确: 期望 %s, 实际 %s", "测试Agent Agent", def.Name)
		}
		if def.Category != tools.ToolCategoryCustom {
			t.Errorf("工具类别不正确: 期望 %s, 实际 %s", tools.ToolCategoryCustom, def.Category)
		}
		if !def.Enabled {
			t.Error("工具应该启用")
		}
	})

	// 测试2: 验证参数验证
	t.Run("参数验证", func(t *testing.T) {
		// 有效参数
		validParams := map[string]interface{}{
			"task": "测试任务",
			"parameters": map[string]interface{}{
				"param1": "value1",
			},
		}
		if err := agentTool.ValidateParameters(validParams); err != nil {
			t.Errorf("有效参数验证失败: %v", err)
		}

		// 无效参数：缺少task
		invalidParams1 := map[string]interface{}{
			"parameters": map[string]interface{}{},
		}
		if err := agentTool.ValidateParameters(invalidParams1); err == nil {
			t.Error("缺少task参数应该验证失败")
		}

		// 无效参数：task不是字符串
		invalidParams2 := map[string]interface{}{
			"task": 123,
		}
		if err := agentTool.ValidateParameters(invalidParams2); err == nil {
			t.Error("task参数不是字符串应该验证失败")
		}
	})

	// 测试3: 执行工具
	t.Run("执行工具", func(t *testing.T) {
		ctx := context.Background()
		execution := tools.ToolExecution{
			ToolID: "agent_test-agent-001",
			Parameters: map[string]interface{}{
				"task": "执行测试任务",
				"parameters": map[string]interface{}{
					"test_param": "test_value",
				},
			},
		}

		result, err := agentTool.Execute(ctx, execution)
		if err != nil {
			t.Fatalf("执行工具失败: %v", err)
		}

		if !result.Success {
			t.Errorf("工具执行应该成功，但返回失败: %s", result.Error)
		}

		// 验证结果包含Agent信息
		output, ok := result.Output.(map[string]interface{})
		if !ok {
			t.Fatal("输出结果不是map类型")
		}

		if output["agent_id"] != "test-agent-001" {
			t.Errorf("输出缺少agent_id: %v", output)
		}

		if output["agent_name"] != "测试Agent" {
			t.Errorf("输出缺少agent_name: %v", output)
		}
	})

	// 测试4: 启用/禁用状态
	t.Run("启用禁用状态", func(t *testing.T) {
		if !agentTool.IsEnabled() {
			t.Error("工具初始状态应该启用")
		}

		agentTool.SetEnabled(false)
		if agentTool.IsEnabled() {
			t.Error("工具禁用后应该返回false")
		}

		agentTool.SetEnabled(true)
		if !agentTool.IsEnabled() {
			t.Error("工具启用后应该返回true")
		}
	})

	// 测试5: 统计信息
	t.Run("统计信息", func(t *testing.T) {
		// 执行一次以更新统计
		ctx := context.Background()
		execution := tools.ToolExecution{
			ToolID: "agent_test-agent-001",
			Parameters: map[string]interface{}{
				"task": "更新统计的任务",
			},
		}

		_, err := agentTool.Execute(ctx, execution)
		if err != nil {
			t.Fatalf("执行工具失败: %v", err)
		}

		count := agentTool.GetExecutionCount()
		if count < 1 {
			t.Errorf("执行次数应该至少为1，实际: %d", count)
		}

		avgDuration := agentTool.GetAverageDuration()
		if avgDuration <= 0 {
			t.Errorf("平均执行时间应该大于0，实际: %v", avgDuration)
		}

		lastTime := agentTool.GetLastExecutionTime()
		if lastTime.IsZero() {
			t.Error("最后执行时间不应该为零")
		}
	})
}

// TestAgentToolWithDisabledAgent 测试禁用Agent的工具
func TestAgentToolWithDisabledAgent(t *testing.T) {
	// 创建禁用的模拟Agent
	agentConfig := AgentConfig{
		ID:          "disabled-agent",
		Name:        "禁用Agent",
		Role:        RoleFullStackDev,
		Description: "禁用的Agent",
		Enabled:     false,
	}

	mockAgent := NewMockAgent(agentConfig)
	mockAgent.SetEnabled(false)

	agentTool, err := NewAgentTool(mockAgent)
	if err != nil {
		t.Fatalf("创建AgentTool失败: %v", err)
	}

	if agentTool.IsEnabled() {
		t.Error("禁用Agent创建的工具应该也是禁用的")
	}

	// 尝试执行应该失败
	ctx := context.Background()
	execution := tools.ToolExecution{
		ToolID: "agent_disabled-agent",
		Parameters: map[string]interface{}{
			"task": "测试任务",
		},
	}

	result, err := agentTool.Execute(ctx, execution)
	if err != nil {
		t.Fatalf("执行工具不应该返回错误: %v", err)
	}

	if result.Success {
		t.Error("禁用Agent的工具执行应该失败")
	}
}
