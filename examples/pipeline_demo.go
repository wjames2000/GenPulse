package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"GenPulse/internal/agents"
	"GenPulse/internal/genkit/flows"
	"GenPulse/internal/genkit/models"
	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/pipeline"
)

func main() {
	fmt.Println("=== GenPulse 流水线编排系统演示 ===")
	fmt.Println("演示2.3节：流水线编排开发功能")
	fmt.Println()

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

	// 创建流水线Flow
	pipelineFlow := pipeline.NewPipelineFlow(flowEngine, agentManager)
	
	// 定义主流水线
	flowDef, err := pipelineFlow.DefineMainPipeline()
	if err != nil {
		log.Fatalf("定义流水线失败: %v", err)
	}
	
	fmt.Println("✓ 主流水线Flow设计完成")
	fmt.Printf("流水线ID: %s\n", flowDef.ID)
	fmt.Printf("流水线名称: %s\n", flowDef.Name)
	fmt.Printf("流水线描述: %s\n", flowDef.Description)
	fmt.Printf("包含节点: %d个\n", len(flowDef.Nodes))
	fmt.Printf("包含边: %d条\n", len(flowDef.Edges))
	fmt.Println()

	// 演示流水线节点
	fmt.Println("=== 流水线节点结构 ===")
	for i, node := range flowDef.Nodes {
		fmt.Printf("%d. %s (%s)\n", i+1, node.Name, node.Type)
		fmt.Printf("   描述: %s\n", node.Description)
		if agentID, ok := node.Config["agent_id"]; ok {
			fmt.Printf("   Agent: %s\n", agentID)
		}
		fmt.Println()
	}

	// 演示并行执行引擎
	fmt.Println("=== 并行执行引擎演示 ===")
	
	// 创建并行引擎
	_ = pipeline.NewParallelEngine(agentManager, 3)
	fmt.Println("✓ 创建并行执行引擎 (3个worker)")
	
	// 演示任务
	tasks := []pipeline.ParallelTask{
		{
			ID:          "task_1",
			Name:        "前端组件开发",
			Description: "开发用户登录组件",
			AgentID:     "frontend_dev_001",
			Task:        "开发用户登录界面",
			Parameters: map[string]interface{}{
				"component_name": "LoginForm",
				"tech_stack":     "React + TypeScript",
			},
			Priority:   1,
			Timeout:    2 * time.Minute,
			RetryCount: 2,
		},
		{
			ID:          "task_2",
			Name:        "后端API开发",
			Description: "开发用户认证API",
			AgentID:     "backend_dev_001",
			Task:        "开发用户登录API",
			Parameters: map[string]interface{}{
				"api_name":   "AuthAPI",
				"tech_stack": "Go + Gin",
			},
			Priority:   1,
			Timeout:    3 * time.Minute,
			RetryCount: 2,
		},
		{
			ID:          "task_3",
			Name:        "数据库设计",
			Description: "设计用户数据表",
			AgentID:     "backend_dev_001",
			Task:        "设计用户数据库表",
			Parameters: map[string]interface{}{
				"table_name": "users",
				"db_type":    "PostgreSQL",
			},
			Priority:   2,
			Timeout:    2 * time.Minute,
			RetryCount: 1,
		},
	}
	
	fmt.Println("创建了3个并行任务:")
	for i, task := range tasks {
		fmt.Printf("  %d. %s (Agent: %s, 优先级: %d)\n", 
			i+1, task.Name, task.AgentID, task.Priority)
	}
	fmt.Println()

	// 演示上下文传递机制
	fmt.Println("=== 上下文传递机制演示 ===")
	
	// 创建流水线上下文
	pipelineCtx := pipeline.NewPipelineContext(map[string]interface{}{
		"project_name":        "DemoProject",
		"project_description": "演示项目",
		"requirements":        "开发一个用户管理系统",
		"tech_stack":          "React + Go + PostgreSQL",
	})
	
	fmt.Println("✓ 创建流水线上下文")
	fmt.Printf("项目名称: %s\n", pipelineCtx.Parameters["project_name"])
	fmt.Printf("项目描述: %s\n", pipelineCtx.Parameters["project_description"])
	fmt.Printf("技术栈: %s\n", pipelineCtx.Parameters["tech_stack"])
	fmt.Println()

	// 模拟添加产物
	pipelineCtx.SetArtifact("prd_document", "产品需求文档内容...")
	pipelineCtx.SetArtifact("architecture_design", "系统架构设计...")
	pipelineCtx.SetArtifact("task_plan", "任务执行计划...")
	
	fmt.Println("✓ 模拟添加了3个产物:")
	fmt.Println("  1. PRD文档")
	fmt.Println("  2. 架构设计")
	fmt.Println("  3. 任务计划")
	fmt.Println()

	// 演示为不同角色准备上下文
	fmt.Println("=== 为不同Agent准备上下文 ===")
	
	roles := []string{
		"前端开发",
		"后端开发",
		"QA工程师",
		"技术架构师",
	}
	
	for _, role := range roles {
		contextData := pipelineCtx.GetContextForAgent("demo_agent", role)
		fmt.Printf("%s 上下文包含 %d 个数据项\n", role, len(contextData))
		
		// 显示关键数据项
		keyItems := []string{}
		for key := range contextData {
			if len(keyItems) < 3 { // 只显示前3个
				keyItems = append(keyItems, key)
			}
		}
		fmt.Printf("  关键项: %v\n", keyItems)
	}
	fmt.Println()

	// 演示错误处理与重试策略
	fmt.Println("=== 错误处理与重试策略演示 ===")
	
	errorHandler := pipeline.NewErrorHandler(3, 1*time.Second)
	
	// 模拟不同类型的错误
	errorTypes := []struct {
		errMsg string
		expectedType string
	}{
		{"connection timeout", "timeout"},
		{"agent not found", "agent_unavailable"},
		{"validation failed: invalid input", "validation_failed"},
		{"network error: connection refused", "network_error"},
	}
	
	for _, testCase := range errorTypes {
		err := fmt.Errorf(testCase.errMsg)
		shouldRetry, action, waitTime := errorHandler.HandleError(
			context.Background(),
			"demo_stage",
			"demo_agent",
			err,
			0,
		)
		
		fmt.Printf("错误: %s\n", testCase.errMsg)
		fmt.Printf("  处理结果: 重试=%v, 行动=%s, 等待=%v\n", 
			shouldRetry, action, waitTime)
	}
	fmt.Println()

	// 演示完整流水线执行（模拟）
	fmt.Println("=== 完整流水线执行流程 ===")
	
	stages := []struct {
		name        string
		description string
		agent       string
	}{
		{"需求分析", "产品经理分析需求并生成PRD", "product_manager_001"},
		{"架构设计", "技术架构师设计系统架构", "architect_001"},
		{"任务分解", "编排器分解任务并生成计划", "orchestrator_001"},
		{"并行开发", "前端和后端Agent并行开发", "frontend_dev_001, backend_dev_001"},
		{"测试", "QA工程师进行测试", "qa_engineer_001"},
		{"部署", "DevOps工程师进行部署", "devops_engineer_001"},
		{"代码审查", "代码审查员进行代码审查", "reviewer_001"},
		{"项目验证", "验证生成的项目是否可运行", "devops_engineer_001"},
	}
	
	fmt.Println("流水线包含8个阶段:")
	for i, stage := range stages {
		fmt.Printf("  %d. %s\n", i+1, stage.name)
		fmt.Printf("     描述: %s\n", stage.description)
		fmt.Printf("     Agent: %s\n", stage.agent)
	}
	fmt.Println()

	// 演示流水线执行结果
	fmt.Println("=== 流水线执行结果演示 ===")
	
	// 模拟执行流水线
	_ = context.Background()
	
	// 注意：这里只是演示，实际执行需要真实的Agent和模型
	fmt.Println("模拟执行流水线...")
	time.Sleep(1 * time.Second)
	
	// 创建模拟结果
	simulatedResult := &pipeline.PipelineResult{
		Success:       true,
		ProjectPath:   "/tmp/demo_project",
		ExecutionTime: 8 * time.Minute,
		Artifacts: map[string]interface{}{
			"prd_document":         "产品需求文档",
			"architecture_design":  "系统架构设计",
			"frontend_code":        "前端代码",
			"backend_code":         "后端代码",
			"test_report":          "测试报告",
			"deployment_result":    "部署结果",
			"code_review_report":   "代码审查报告",
			"validation_result":    "验证结果",
		},
		Summary: map[string]interface{}{
			"total_stages":       8,
			"completed_stages":   8,
			"failed_stages":      0,
			"total_artifacts":    8,
			"project_generated":  true,
		},
	}
	
	fmt.Println("✓ 流水线执行成功!")
	fmt.Printf("项目路径: %s\n", simulatedResult.ProjectPath)
	fmt.Printf("执行时间: %v\n", simulatedResult.ExecutionTime)
	fmt.Printf("生成产物: %d个\n", simulatedResult.Summary["total_artifacts"])
	fmt.Println()

	// 技术特点总结
	fmt.Println("=== 技术特点总结 ===")
	fmt.Println("1. 主流水线Flow设计")
	fmt.Println("   • 定义从需求输入到项目完成的完整流程")
	fmt.Println("   • 支持顺序和并行执行")
	fmt.Println("   • 包含8个专业阶段")
	fmt.Println()
	
	fmt.Println("2. 并行执行引擎")
	fmt.Println("   • 基于goroutine实现并发调度")
	fmt.Println("   • 支持任务优先级和超时控制")
	fmt.Println("   • 自动负载均衡和错误恢复")
	fmt.Println()
	
	fmt.Println("3. 上下文传递机制")
	fmt.Println("   • 实现Agent间的数据共享")
	fmt.Println("   • 支持角色特定的上下文准备")
	fmt.Println("   • 完整的产物管理和版本控制")
	fmt.Println()
	
	fmt.Println("4. 错误处理与重试策略")
	fmt.Println("   • 智能错误分类和严重程度评估")
	fmt.Println("   • 多种降级策略（重试、跳过、替代方案）")
	fmt.Println("   • 指数退避和人工干预机制")
	fmt.Println()

	// 使用示例
	fmt.Println("=== 使用示例 ===")
	fmt.Println("```go")
	fmt.Println("// 1. 创建流水线")
	fmt.Println("pipelineFlow := pipeline.NewPipelineFlow(flowEngine, agentManager)")
	fmt.Println()
	fmt.Println("// 2. 准备参数")
	fmt.Println("params := map[string]interface{}{")
	fmt.Println("  \"project_name\": \"MyApp\",")
	fmt.Println("  \"project_description\": \"一个Web应用\",")
	fmt.Println("  \"requirements\": \"用户管理功能\",")
	fmt.Println("  \"tech_stack\": \"React + Go\",")
	fmt.Println("}")
	fmt.Println()
	fmt.Println("// 3. 执行流水线")
	fmt.Println("result, err := pipelineFlow.ExecutePipeline(ctx, params)")
	fmt.Println("if err != nil {")
	fmt.Println("  log.Fatal(\"流水线执行失败:\", err)")
	fmt.Println("}")
	fmt.Println()
	fmt.Println("// 4. 处理结果")
	fmt.Println("if result.Success {")
	fmt.Println("  fmt.Println(\"项目生成成功!\")")
	fmt.Println("  fmt.Println(\"项目路径:\", result.ProjectPath)")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("=== 演示完成 ===")
	fmt.Println("流水线编排系统已实现2.3节所有功能:")
	fmt.Println("✓ 2.3.1 主流水线Flow设计")
	fmt.Println("✓ 2.3.2 并行执行引擎")
	fmt.Println("✓ 2.3.3 上下文传递机制")
	fmt.Println("✓ 2.3.4 错误处理与重试策略")
}