package main

import (
	"context"
	"fmt"
	"log"

	"GenPulse/internal/agents"
	"GenPulse/internal/genkit/flows"
	"GenPulse/internal/genkit/models"
	"GenPulse/internal/genkit/tools"
)

func main() {
	fmt.Println("=== GenPulse Agent 系统演示 ===")
	
	// 创建模拟的依赖组件
	modelAdapter := &models.UnifiedModelAdapter{}
	toolRegistry := &tools.ToolRegistry{}
	flowEngine := &flows.FlowEngine{}
	
	// 创建Agent管理器
	agentManager := agents.NewAgentManager(modelAdapter, toolRegistry, flowEngine)
	
	// 初始化Agent管理器
	if err := agentManager.Initialize(); err != nil {
		log.Printf("Agent管理器初始化失败: %v", err)
	} else {
		fmt.Println("✓ Agent管理器初始化成功")
	}
	
	// 获取所有Agent状态
	status := agentManager.GetAllAgentsStatus()
	fmt.Printf("\n=== Agent系统状态 ===\n")
	fmt.Printf("总Agent数量: %d\n", status["total_agents"])
	fmt.Printf("启用Agent数量: %d\n", status["enabled_agents"])
	fmt.Printf("忙碌Agent数量: %d\n", status["busy_agents"])
	fmt.Printf("空闲Agent数量: %d\n", status["idle_agents"])
	
	// 显示所有Agent信息
	fmt.Println("\n=== 所有Agent角色 ===")
	agentsMap := status["agents"].(map[string]interface{})
	for _, agentStatus := range agentsMap {
		statusMap := agentStatus.(map[string]interface{})
		fmt.Printf("• %s (%s): %s\n", 
			statusMap["name"].(string),
			statusMap["role"].(string),
			statusMap["description"].(string))
	}
	
	// 演示各个Agent的功能
	fmt.Println("\n=== Agent功能演示 ===")
	
	// 1. 编排器Agent
	fmt.Println("\n1. 编排器Agent (Orchestrator):")
	fmt.Println("   - 任务分解和规划")
	fmt.Println("   - 执行计划生成")
	fmt.Println("   - Agent调度和协调")
	
	// 2. 产品经理Agent
	fmt.Println("\n2. 产品经理Agent (Product Manager):")
	fmt.Println("   - 需求分析和整理")
	fmt.Println("   - PRD文档生成")
	fmt.Println("   - 用户故事定义")
	fmt.Println("   - 产品路线图规划")
	
	// 3. 技术架构师Agent
	fmt.Println("\n3. 技术架构师Agent (Architect):")
	fmt.Println("   - 技术架构设计")
	fmt.Println("   - 系统设计文档")
	fmt.Println("   - 技术方案输出")
	fmt.Println("   - 数据库设计")
	
	// 4. 前端开发Agent
	fmt.Println("\n4. 前端开发Agent (Frontend Developer):")
	fmt.Println("   - React/Vue组件开发")
	fmt.Println("   - UI界面实现")
	fmt.Println("   - 状态管理")
	fmt.Println("   - API集成")
	
	// 5. 后端开发Agent
	fmt.Println("\n5. 后端开发Agent (Backend Developer):")
	fmt.Println("   - API接口开发")
	fmt.Println("   - 数据库操作")
	fmt.Println("   - 业务逻辑实现")
	fmt.Println("   - 系统集成")
	
	// 6. QA工程师Agent
	fmt.Println("\n6. QA工程师Agent (QA Engineer):")
	fmt.Println("   - 测试用例生成")
	fmt.Println("   - 测试执行")
	fmt.Println("   - 质量保证")
	fmt.Println("   - 缺陷分析")
	
	// 7. DevOps工程师Agent
	fmt.Println("\n7. DevOps工程师Agent (DevOps Engineer):")
	fmt.Println("   - 项目构建")
	fmt.Println("   - 部署验证")
	fmt.Println("   - 基础设施设置")
	fmt.Println("   - 监控配置")
	
	// 8. 代码审查员Agent
	fmt.Println("\n8. 代码审查员Agent (Reviewer):")
	fmt.Println("   - 代码审查")
	fmt.Println("   - 安全扫描")
	fmt.Println("   - 质量检查")
	fmt.Println("   - 最佳实践检查")
	
	// 9. 全栈开发Agent
	fmt.Println("\n9. 全栈开发Agent (Fullstack Developer):")
	fmt.Println("   - 前后端全栈开发")
	fmt.Println("   - 项目初始化")
	fmt.Println("   - 完整功能实现")
	
	// 演示工作流程
	fmt.Println("\n=== 典型工作流程 ===")
	fmt.Println("1. 产品经理分析需求并生成PRD")
	fmt.Println("2. 技术架构师设计系统架构")
	fmt.Println("3. 编排器分解任务并分配")
	fmt.Println("4. 前端开发实现用户界面")
	fmt.Println("5. 后端开发实现API和业务逻辑")
	fmt.Println("6. QA工程师进行测试")
	fmt.Println("7. DevOps工程师部署项目")
	fmt.Println("8. 代码审查员进行代码审查")
	fmt.Println("9. 全栈开发处理跨领域任务")
	
	fmt.Println("\n=== 使用示例 ===")
	
	// 获取特定Agent
	orchestrator, err := agentManager.GetAgent("orchestrator_001")
	if err == nil {
		config := orchestrator.GetConfig()
		fmt.Printf("获取编排器Agent: %s (%s)\n", config.Name, config.Role)
		fmt.Printf("描述: %s\n", config.Description)
		fmt.Printf("能力: %v\n", config.Capabilities)
	}
	
	// 演示如何执行任务
	fmt.Println("\n执行任务示例:")
	fmt.Println("1. 创建上下文: ctx := context.Background()")
	fmt.Println("2. 准备参数: params := map[string]interface{}{\"project\": \"myapp\"}")
	fmt.Println("3. 执行任务: result, err := agent.Execute(ctx, \"设计系统架构\", params)")
	fmt.Println("4. 处理结果: if result.Success { ... }")
	
	fmt.Println("\n=== 总结 ===")
	fmt.Println("✓ 已实现8个专业Agent角色")
	fmt.Println("✓ 每个Agent都有专门的职责和能力")
	fmt.Println("✓ 支持通过AgentManager统一管理")
	fmt.Println("✓ 可以注册为Tool供其他Agent调用")
	fmt.Println("✓ 支持配置化创建和初始化")
	
	fmt.Println("\n演示完成!")
}

// 模拟执行任务
func demoTaskExecution(ctx context.Context, agent agents.Agent, task string) {
	fmt.Printf("\n执行任务: %s\n", task)
	
	params := map[string]interface{}{
		"demo": true,
		"task": task,
	}
	
	// 在实际应用中，这里会调用agent.Execute(ctx, task, params)
	fmt.Printf("Agent: %s\n", agent.GetConfig().Name)
	fmt.Printf("角色: %s\n", agent.GetConfig().Role)
	fmt.Printf("参数: %v\n", params)
	fmt.Println("状态: 就绪 (在实际应用中会调用模型生成结果)")
}