package main

import (
	"context"
	"fmt"
	"log"

	"GenPulse/internal/agents"
	"GenPulse/internal/genkit/models"
	"GenPulse/internal/genkit/tools"
)

// MockFlowEngine 模拟Flow引擎
type MockFlowEngine struct{}

func (m *MockFlowEngine) ExecuteFlow(ctx context.Context, flowID string, parameters map[string]interface{}) (interface{}, error) {
	return nil, nil
}

func main() {
	fmt.Println("=== Agent间调用机制演示 ===")

	// 1. 初始化依赖组件
	fmt.Println("\n1. 初始化依赖组件...")
	
	// 创建模型适配器（模拟）
	modelAdapter := &models.UnifiedModelAdapter{}
	
	// 创建工具注册表
	toolRegistry := tools.NewToolRegistry()
	
	// 创建Flow引擎（模拟）
	flowEngine := &MockFlowEngine{}
	
	// 2. 创建Agent管理器
	fmt.Println("2. 创建Agent管理器...")
	agentManager := agents.NewAgentManager(modelAdapter, toolRegistry, flowEngine)
	
	// 3. 创建两个测试Agent
	fmt.Println("3. 创建测试Agent...")
	
	// Agent 1: 架构师
	architectConfig := agents.AgentConfig{
		ID:          "architect_001",
		Name:        "架构师",
		Role:        agents.RoleArchitect,
		Description: "负责系统架构设计",
		ModelConfig: models.ModelConfig{
			Name:     "gemini-1.5-pro",
			Provider: "google",
			Type:     models.ModelTypeGemini,
		},
		Enabled: true,
	}
	
	// Agent 2: 全栈开发工程师
	developerConfig := agents.AgentConfig{
		ID:          "developer_001",
		Name:        "全栈开发工程师",
		Role:        agents.RoleFullStackDev,
		Description: "负责前后端代码实现",
		ModelConfig: models.ModelConfig{
			Name:     "gemini-1.5-pro",
			Provider: "google",
			Type:     models.ModelTypeGemini,
		},
		Enabled: true,
	}
	
	// 4. 创建Agent实例
	fmt.Println("4. 创建Agent实例...")
	
	// 创建架构师Agent
	architectAgent, err := agents.NewBaseAgent(architectConfig, modelAdapter, toolRegistry, flowEngine)
	if err != nil {
		log.Fatalf("创建架构师Agent失败: %v", err)
	}
	
	// 创建开发工程师Agent
	developerAgent, err := agents.NewBaseAgent(developerConfig, modelAdapter, toolRegistry, flowEngine)
	if err != nil {
		log.Fatalf("创建开发工程师Agent失败: %v", err)
	}
	
	// 5. 注册Agent到管理器
	fmt.Println("5. 注册Agent到管理器...")
	if err := agentManager.RegisterAgent(architectAgent); err != nil {
		log.Fatalf("注册架构师Agent失败: %v", err)
	}
	
	if err := agentManager.RegisterAgent(developerAgent); err != nil {
		log.Fatalf("注册开发工程师Agent失败: %v", err)
	}
	
	// 6. 初始化Agent管理器（会自动将Agent注册为Tool）
	fmt.Println("6. 初始化Agent管理器...")
	if err := agentManager.Initialize(); err != nil {
		log.Fatalf("初始化Agent管理器失败: %v", err)
	}
	
	// 7. 演示Agent间调用
	fmt.Println("\n7. 演示Agent间调用机制...")
	
	// 7.1 查看注册的工具
	fmt.Println("\n7.1 查看注册的工具:")
	toolList := toolRegistry.ListTools()
	for _, tool := range toolList {
		fmt.Printf("  - %s: %s (类别: %s)\n", tool.ID, tool.Name, tool.Category)
	}
	
	// 7.2 通过工具调用架构师Agent
	fmt.Println("\n7.2 通过工具调用架构师Agent:")
	ctx := context.Background()
	
	// 创建工具执行请求
	architectExecution := tools.ToolExecution{
		ToolID: "agent_architect_001",
		Parameters: map[string]interface{}{
			"task": "设计一个微服务电商系统架构",
			"parameters": map[string]interface{}{
				"requirements": "需要支持用户管理、商品管理、订单管理、支付功能",
				"tech_stack":   "Go + React + PostgreSQL + Redis + Docker",
			},
		},
	}
	
	// 获取架构师工具
	architectTool, err := toolRegistry.GetTool("agent_architect_001")
	if err != nil {
		log.Fatalf("获取架构师工具失败: %v", err)
	}
	
	// 执行工具
	result, err := architectTool.Execute(ctx, architectExecution)
	if err != nil {
		log.Fatalf("执行架构师工具失败: %v", err)
	}
	
	fmt.Printf("执行结果: 成功=%v, 耗时=%v\n", result.Success, result.Duration)
	
	// 7.3 通过工具调用开发工程师Agent
	fmt.Println("\n7.3 通过工具调用开发工程师Agent:")
	
	developerExecution := tools.ToolExecution{
		ToolID: "agent_developer_001",
		Parameters: map[string]interface{}{
			"task": "根据架构设计实现用户管理模块",
			"parameters": map[string]interface{}{
				"architecture": "微服务电商系统",
				"module":       "用户管理",
			},
		},
	}
	
	// 获取开发工程师工具
	developerTool, err := toolRegistry.GetTool("agent_developer_001")
	if err != nil {
		log.Fatalf("获取开发工程师工具失败: %v", err)
	}
	
	// 执行工具
	result2, err := developerTool.Execute(ctx, developerExecution)
	if err != nil {
		log.Fatalf("执行开发工程师工具失败: %v", err)
	}
	
	fmt.Printf("执行结果: 成功=%v, 耗时=%v\n", result2.Success, result2.Duration)
	
	// 8. 演示Agent编排概念
	fmt.Println("\n8. 演示Agent编排概念:")
	fmt.Println("   通过将Agent封装为Tool，可以实现:")
	fmt.Println("   - Agent A 调用 Agent B 的工具")
	fmt.Println("   - Flow引擎顺序/并行执行多个Agent")
	fmt.Println("   - 上下文在Agent间传递")
	fmt.Println("   - 错误处理和重试机制")
	
	// 9. 统计信息
	fmt.Println("\n9. 统计信息:")
	
	// Agent统计
	agentList := agentManager.ListAgents()
	fmt.Printf("注册的Agent数量: %d\n", len(agentList))
	for _, agent := range agentList {
		fmt.Printf("  - %s (%s): %s\n", agent.Name, agent.ID, agent.Role)
	}
	
	// 工具统计
	toolStats := toolRegistry.GetToolStatistics()
	fmt.Printf("\n工具统计:\n")
	fmt.Printf("  总工具数: %d\n", toolStats["total_tools"])
	fmt.Printf("  启用工具: %d\n", toolStats["enabled_tools"])
	fmt.Printf("  禁用工具: %d\n", toolStats["disabled_tools"])
	
	// 10. 清理
	fmt.Println("\n10. 清理资源...")
	
	// 注销Agent
	if err := agentManager.UnregisterAgent("architect_001"); err != nil {
		log.Printf("注销架构师Agent失败: %v", err)
	}
	
	if err := agentManager.UnregisterAgent("developer_001"); err != nil {
		log.Printf("注销开发工程师Agent失败: %v", err)
	}
	
	fmt.Println("\n=== 演示完成 ===")
	fmt.Println("总结: Agent间调用机制已成功实现，Agent可以被封装为Tool，并通过ToolRegistry进行调用和编排。")
}