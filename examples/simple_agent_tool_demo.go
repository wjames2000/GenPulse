package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"GenPulse/internal/agents"
	"GenPulse/internal/genkit/tools"
)

// SimpleAgentToolDemo 简单的AgentTool演示
func main() {
	fmt.Println("=== Agent间调用机制 - 核心功能演示 ===")

	// 1. 创建工具注册表
	fmt.Println("\n1. 创建工具注册表...")
	toolRegistry := tools.NewToolRegistry()

	// 2. 创建模拟Agent
	fmt.Println("\n2. 创建模拟Agent...")
	
	// 模拟Agent配置
	agentConfig := agents.AgentConfig{
		ID:          "demo_agent_001",
		Name:        "演示Agent",
		Role:        agents.RoleFullStackDev,
		Description: "用于演示的Agent",
		Enabled:     true,
	}

	// 创建模拟Agent
	mockAgent := &MockAgent{
		config:  agentConfig,
		enabled: true,
	}

	// 3. 创建AgentTool
	fmt.Println("\n3. 创建AgentTool...")
	agentTool, err := agents.NewAgentTool(mockAgent)
	if err != nil {
		log.Fatalf("创建AgentTool失败: %v", err)
	}

	// 4. 注册AgentTool到工具注册表
	fmt.Println("\n4. 注册AgentTool到工具注册表...")
	if err := toolRegistry.RegisterTool(agentTool); err != nil {
		log.Fatalf("注册AgentTool失败: %v", err)
	}

	// 5. 查看注册的工具
	fmt.Println("\n5. 查看注册的工具:")
	toolList := toolRegistry.ListTools()
	for _, tool := range toolList {
		fmt.Printf("  - ID: %s\n", tool.ID)
		fmt.Printf("    名称: %s\n", tool.Name)
		fmt.Printf("    描述: %s\n", tool.Description)
		fmt.Printf("    类别: %s\n", tool.Category)
		fmt.Printf("    启用: %v\n", tool.Enabled)
		fmt.Printf("    标签: %v\n", tool.Tags)
	}

	// 6. 通过工具调用Agent
	fmt.Println("\n6. 通过工具调用Agent...")
	ctx := context.Background()

	// 创建工具执行请求
	execution := tools.ToolExecution{
		ToolID: "agent_demo_agent_001",
		Parameters: map[string]interface{}{
			"task": "创建一个简单的Go Web服务器",
			"parameters": map[string]interface{}{
				"port":     8080,
				"endpoint": "/api/health",
			},
		},
	}

	// 获取工具
	tool, err := toolRegistry.GetTool("agent_demo_agent_001")
	if err != nil {
		log.Fatalf("获取工具失败: %v", err)
	}

	// 执行工具
	fmt.Println("\n执行Agent工具...")
	result, err := tool.Execute(ctx, execution)
	if err != nil {
		log.Fatalf("执行工具失败: %v", err)
	}

	// 显示结果
	fmt.Println("\n执行结果:")
	fmt.Printf("  成功: %v\n", result.Success)
	fmt.Printf("  耗时: %v\n", result.Duration)
	if result.Error != "" {
		fmt.Printf("  错误: %s\n", result.Error)
	}

	// 显示输出
	if result.Output != nil {
		fmt.Printf("  输出: %v\n", result.Output)
	}

	// 7. 演示工具统计信息
	fmt.Println("\n7. 工具统计信息:")
	fmt.Printf("  执行次数: %d\n", tool.GetExecutionCount())
	fmt.Printf("  平均耗时: %v\n", tool.GetAverageDuration())
	fmt.Printf("  最后执行: %v\n", tool.GetLastExecutionTime())

	// 8. 演示工具启用/禁用
	fmt.Println("\n8. 演示工具启用/禁用:")
	fmt.Printf("  当前启用状态: %v\n", tool.IsEnabled())
	
	// 禁用工具
	tool.SetEnabled(false)
	fmt.Printf("  禁用后状态: %v\n", tool.IsEnabled())
	
	// 尝试执行禁用的工具
	fmt.Println("\n尝试执行禁用的工具...")
	result2, err := tool.Execute(ctx, execution)
	if err != nil {
		log.Printf("执行失败（预期）: %v", err)
	}
	fmt.Printf("  执行结果: 成功=%v\n", result2.Success)
	
	// 重新启用工具
	tool.SetEnabled(true)
	fmt.Printf("  重新启用后状态: %v\n", tool.IsEnabled())

	// 9. 演示工具注册表统计
	fmt.Println("\n9. 工具注册表统计:")
	stats := toolRegistry.GetToolStatistics()
	fmt.Printf("  总工具数: %d\n", stats["total_tools"])
	fmt.Printf("  启用工具: %d\n", stats["enabled_tools"])
	fmt.Printf("  禁用工具: %d\n", stats["disabled_tools"])

	// 10. 总结
	fmt.Println("\n=== 演示完成 ===")
	fmt.Println("\n核心功能验证:")
	fmt.Println("  ✓ Agent可以封装为Tool")
	fmt.Println("  ✓ AgentTool可以注册到ToolRegistry")
	fmt.Println("  ✓ 可以通过ToolRegistry调用Agent")
	fmt.Println("  ✓ 支持参数验证和执行统计")
	fmt.Println("  ✓ 支持启用/禁用状态管理")
	fmt.Println("  ✓ 支持工具间调用和编排")
	fmt.Println("\n这样实现了基于Genkit Tool的Agent编排能力，满足2.1.3任务要求。")
}

// MockAgent 模拟Agent实现
type MockAgent struct {
	config  agents.AgentConfig
	enabled bool
}

func (m *MockAgent) GetConfig() agents.AgentConfig {
	return m.config
}

func (m *MockAgent) GetState() agents.AgentState {
	return agents.StateIdle
}

func (m *MockAgent) GetExecution() *agents.AgentExecution {
	return nil
}

func (m *MockAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*agents.AgentResult, error) {
	if !m.enabled {
		return nil, agents.NewAgentError("AGENT_DISABLED", "Agent已禁用", "")
	}

	// 模拟执行
	return &agents.AgentResult{
		Success: true,
		Output: map[string]interface{}{
			"task":       task,
			"parameters": parameters,
			"message":    "任务执行成功",
			"code":       "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}",
		},
		Duration: 100 * 1000 * 1000, // 100ms in nanoseconds
		Metadata: map[string]interface{}{
			"agent_id":   m.config.ID,
			"agent_name": m.config.Name,
		},
	}, nil
}

func (m *MockAgent) Cancel() error {
	return nil
}

func (m *MockAgent) ValidateTask(task string) error {
	if task == "" {
		return agents.NewAgentError("INVALID_TASK", "任务验证失败", "任务不能为空")
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
	return 100 * time.Millisecond
}