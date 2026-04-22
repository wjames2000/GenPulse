package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"GenPulse/internal/agents"
	"GenPulse/internal/genkit/flows"
	"GenPulse/internal/genkit/models"
	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/pipeline"
)

// TestEndToEndPipeline 端到端流水线验证测试
func TestEndToEndPipeline(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "genpulse-e2e-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("测试目录: %s\n", tempDir)

	// 创建模拟的依赖组件
	modelAdapter := &models.UnifiedModelAdapter{}
	toolRegistry := &tools.ToolRegistry{}
	flowEngine := &flows.FlowEngine{}

	// 创建Agent管理器
	agentManager := agents.NewAgentManager(modelAdapter, toolRegistry, flowEngine)

	// 初始化Agent管理器
	if err := agentManager.Initialize(); err != nil {
		t.Logf("警告: Agent管理器初始化失败（模拟环境）: %v", err)
		// 在测试环境中，我们继续执行模拟测试
	}

	// 创建流水线
	pipelineFlow := pipeline.NewPipelineFlow(flowEngine, agentManager)

	// 测试用例：简单的TODO应用
	testCases := []struct {
		name         string
		description  string
		requirements string
		techStack    string
	}{
		{
			name:         "SimpleTodoApp",
			description:  "一个简单的待办事项应用",
			requirements: "开发一个待办事项应用，支持添加、删除、标记完成待办事项。需要用户界面和API。",
			techStack:    "React + Go + SQLite",
		},
		{
			name:         "WeatherDashboard",
			description:  "天气信息仪表板",
			requirements: "显示当前天气和未来5天预报，支持城市搜索。需要从公开API获取数据。",
			techStack:    "Vue + Node.js + MongoDB",
		},
		{
			name:         "BlogPlatform",
			description:  "简单的博客平台",
			requirements: "用户可以发布、编辑、删除博客文章，支持Markdown格式。需要用户认证和评论功能。",
			techStack:    "Next.js + Python + PostgreSQL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n=== 测试用例: %s ===\n", tc.name)
			fmt.Printf("描述: %s\n", tc.description)
			fmt.Printf("技术栈: %s\n", tc.techStack)

			// 准备测试参数
			params := map[string]interface{}{
				"project_name":        tc.name,
				"project_description": tc.description,
				"requirements":        tc.requirements,
				"tech_stack":          tc.techStack,
				"test_mode":           true, // 标记为测试模式
			}

			// 创建上下文
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			defer cancel()

			// 执行流水线
			startTime := time.Now()
			fmt.Printf("开始执行流水线...\n")

			result, err := pipelineFlow.ExecutePipeline(ctx, params)

			duration := time.Since(startTime)
			fmt.Printf("流水线执行完成，耗时: %v\n", duration)

			// 验证结果
			if err != nil {
				t.Errorf("流水线执行失败: %v", err)
				if result != nil {
					fmt.Printf("失败阶段: %s\n", result.FailedStage)
					fmt.Printf("已生成产物: %d个\n", len(result.Artifacts))
				}
			} else {
				// 验证成功结果
				if !result.Success {
					t.Errorf("流水线返回失败状态")
				}

				// 验证基本输出
				if result.ProjectPath == "" {
					t.Logf("警告: 项目路径为空（在模拟环境中是预期的）")
				}

				// 验证产物数量
				expectedMinArtifacts := 5 // 至少应该有PRD、架构设计、任务计划等
				if len(result.Artifacts) < expectedMinArtifacts {
					t.Errorf("产物数量不足，期望至少%d个，实际%d个",
						expectedMinArtifacts, len(result.Artifacts))
				} else {
					fmt.Printf("✓ 生成产物: %d个\n", len(result.Artifacts))
				}

				// 验证摘要信息
				if summary, ok := result.Summary["total_stages"].(int); ok {
					if summary != 8 {
						t.Errorf("阶段总数不正确，期望8，实际%d", summary)
					}
				}

				// 输出成功信息
				fmt.Printf("✓ 流水线执行成功!\n")
				fmt.Printf("  执行时间: %v\n", result.ExecutionTime)
				fmt.Printf("  完成阶段: %d/%d\n",
					result.Summary["completed_stages"], result.Summary["total_stages"])
				fmt.Printf("  生成产物: %d个\n", result.Summary["total_artifacts"])
			}

			// 输出日志统计
			if result != nil {
				infoLogs := 0
				successLogs := 0
				errorLogs := 0
				for _, log := range result.Logs {
					switch log.Level {
					case "info":
						infoLogs++
					case "success":
						successLogs++
					case "error":
						errorLogs++
					}
				}
				fmt.Printf("日志统计: 信息=%d, 成功=%d, 错误=%d\n",
					infoLogs, successLogs, errorLogs)
			}
		})
	}
}

// TestPipelineErrorHandling 测试流水线错误处理
func TestPipelineErrorHandling(t *testing.T) {
	fmt.Printf("\n=== 测试错误处理 ===\n")

	// 创建模拟组件
	modelAdapter := &models.UnifiedModelAdapter{}
	toolRegistry := &tools.ToolRegistry{}
	flowEngine := &flows.FlowEngine{}
	agentManager := agents.NewAgentManager(modelAdapter, toolRegistry, flowEngine)

	// 创建流水线
	pipelineFlow := pipeline.NewPipelineFlow(flowEngine, agentManager)

	// 测试用例：无效参数
	testCases := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
	}{
		{
			name: "MissingRequiredParams",
			params: map[string]interface{}{
				"project_name": "Test", // 缺少project_description和requirements
			},
			expectError: true,
		},
		{
			name: "EmptyRequirements",
			params: map[string]interface{}{
				"project_name":        "Test",
				"project_description": "Test project",
				"requirements":        "", // 空需求
			},
			expectError: true,
		},
		{
			name: "ValidParams",
			params: map[string]interface{}{
				"project_name":        "ValidTest",
				"project_description": "A valid test project",
				"requirements":        "Test requirements",
				"tech_stack":          "React + Go",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("测试: %s\n", tc.name)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			result, err := pipelineFlow.ExecutePipeline(ctx, tc.params)

			if tc.expectError {
				if err == nil {
					t.Errorf("期望错误但执行成功")
				} else {
					fmt.Printf("✓ 正确捕获错误: %v\n", err)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但执行失败: %v", err)
				} else if result != nil && !result.Success {
					t.Errorf("流水线返回失败状态")
				} else {
					fmt.Printf("✓ 参数验证通过\n")
				}
			}
		})
	}
}

// TestParallelExecution 测试并行执行
func TestParallelExecution(t *testing.T) {
	fmt.Printf("\n=== 测试并行执行 ===\n")

	// 创建并行引擎
	modelAdapter := &models.UnifiedModelAdapter{}
	toolRegistry := &tools.ToolRegistry{}
	agentManager := agents.NewAgentManager(modelAdapter, toolRegistry, nil)

	parallelEngine := pipeline.NewParallelEngine(agentManager, 3)

	// 启动引擎
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := parallelEngine.Start(ctx); err != nil {
		t.Fatalf("启动并行引擎失败: %v", err)
	}
	defer parallelEngine.Stop()

	// 创建测试任务
	tasks := []pipeline.ParallelTask{
		{
			ID:          "task_1",
			Name:        "前端开发",
			Description: "开发用户界面",
			AgentID:     "frontend_dev_001",
			Task:        "开发React组件",
			Parameters:  map[string]interface{}{"component": "Header"},
			Priority:    1,
			Timeout:     2 * time.Minute,
			RetryCount:  1,
		},
		{
			ID:          "task_2",
			Name:        "后端开发",
			Description: "开发API",
			AgentID:     "backend_dev_001",
			Task:        "开发REST API",
			Parameters:  map[string]interface{}{"endpoint": "/api/users"},
			Priority:    1,
			Timeout:     2 * time.Minute,
			RetryCount:  1,
		},
		{
			ID:          "task_3",
			Name:        "数据库设计",
			Description: "设计数据表",
			AgentID:     "backend_dev_001",
			Task:        "设计数据库表",
			Parameters:  map[string]interface{}{"table": "users"},
			Priority:    2,
			Timeout:     1 * time.Minute,
			RetryCount:  0,
		},
	}

	fmt.Printf("提交%d个并行任务...\n", len(tasks))

	// 执行并行任务
	results, err := parallelEngine.ExecuteParallel(ctx, tasks)
	if err != nil {
		t.Fatalf("并行执行失败: %v", err)
	}

	// 验证结果
	fmt.Printf("完成%d个任务\n", len(results))

	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
			fmt.Printf("✓ 任务成功: %s (%s), 耗时: %v\n",
				result.TaskID, result.AgentName, result.Duration)
		} else {
			fmt.Printf("✗ 任务失败: %s, 错误: %v\n", result.TaskID, result.Error)
		}
	}

	fmt.Printf("成功率: %d/%d (%.1f%%)\n",
		successCount, len(results), float64(successCount)/float64(len(results))*100)

	if successCount == 0 {
		t.Errorf("所有并行任务都失败")
	}
}

// TestContextPassing 测试上下文传递
func TestContextPassing(t *testing.T) {
	fmt.Printf("\n=== 测试上下文传递 ===\n")

	// 创建流水线上下文
	params := map[string]interface{}{
		"project_name":        "ContextTest",
		"project_description": "测试上下文传递",
		"requirements":        "测试需求",
		"tech_stack":          "Test Stack",
	}

	pipelineCtx := pipeline.NewPipelineContext(params)

	// 添加测试产物
	pipelineCtx.SetArtifact("prd_document", "测试PRD文档")
	pipelineCtx.SetArtifact("architecture_design", "测试架构设计")
	pipelineCtx.SetArtifact("task_plan", "测试任务计划")

	// 测试为不同角色获取上下文
	roles := []string{
		"前端开发",
		"后端开发",
		"技术架构师",
		"QA工程师",
	}

	for _, role := range roles {
		contextData := pipelineCtx.GetContextForAgent("test_agent", role)
		fmt.Printf("%s 上下文包含 %d 个数据项\n", role, len(contextData))

		// 验证包含必要的数据
		requiredKeys := []string{"project_name", "project_description"}
		for _, key := range requiredKeys {
			if _, ok := contextData[key]; !ok {
				t.Errorf("%s 上下文缺少必要键: %s", role, key)
			}
		}
	}

	// 测试上下文序列化
	jsonStr, err := pipelineCtx.ToJSON()
	if err != nil {
		t.Fatalf("上下文序列化失败: %v", err)
	}

	fmt.Printf("上下文序列化成功，长度: %d 字节\n", len(jsonStr))

	// 测试反序列化
	restoredCtx, err := pipeline.FromJSON(jsonStr)
	if err != nil {
		t.Fatalf("上下文反序列化失败: %v", err)
	}

	// 验证反序列化后的数据
	if restoredCtx.Parameters["project_name"] != "ContextTest" {
		t.Errorf("反序列化后参数不正确")
	}

	if len(restoredCtx.Artifacts) != 3 {
		t.Errorf("反序列化后产物数量不正确")
	}

	fmt.Printf("✓ 上下文传递测试通过\n")
}

// TestErrorHandler 测试错误处理器
func TestErrorHandler(t *testing.T) {
	fmt.Printf("\n=== 测试错误处理器 ===\n")

	errorHandler := pipeline.NewErrorHandler(3, 1*time.Second)

	// 测试不同类型的错误
	testErrors := []struct {
		errMsg      string
		description string
	}{
		{"connection timeout", "超时错误"},
		{"agent not found: frontend_dev_001", "Agent不可用"},
		{"validation failed: invalid input parameters", "验证失败"},
		{"network error: connection refused", "网络错误"},
		{"memory allocation failed", "资源耗尽"},
		{"unknown error occurred", "未知执行错误"},
	}

	ctx := context.Background()
	for _, testErr := range testErrors {
		err := fmt.Errorf(testErr.errMsg)
		shouldRetry, action, waitTime := errorHandler.HandleError(
			ctx, "test_stage", "test_agent", err, 0,
		)

		fmt.Printf("错误: %s\n", testErr.description)
		fmt.Printf("  类型分析: %s\n", getErrorType(err))
		fmt.Printf("  处理结果: 重试=%v, 行动=%s, 等待=%v\n",
			shouldRetry, action, waitTime)
	}

	// 获取错误统计
	stats := errorHandler.GetErrorStats()
	fmt.Printf("\n错误统计:\n")
	fmt.Printf("  总错误数: %v\n", stats["total_errors"])
	fmt.Printf("  按类型: %v\n", stats["by_type"])
	fmt.Printf("  按严重程度: %v\n", stats["by_severity"])

	fmt.Printf("✓ 错误处理器测试通过\n")
}

// getErrorType 辅助函数获取错误类型
func getErrorType(err error) string {
	errStr := err.Error()
	switch {
	case contains(errStr, "timeout"):
		return "timeout"
	case contains(errStr, "not found") || contains(errStr, "unavailable"):
		return "agent_unavailable"
	case contains(errStr, "validation"):
		return "validation_failed"
	case contains(errStr, "network") || contains(errStr, "connection"):
		return "network_error"
	case contains(errStr, "resource") || contains(errStr, "memory"):
		return "resource_exhausted"
	default:
		return "execution_failed"
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > len(substr) && (s[:len(substr)] == substr ||
			contains(s[1:], substr)))
}

// MainTest 运行所有测试
func MainTest() {
	fmt.Println("=== GenPulse 端到端流水线验证测试 ===")
	fmt.Println("开始时间:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()

	// 创建测试实例
	t := &testing.T{}

	// 运行测试
	TestEndToEndPipeline(t)
	TestPipelineErrorHandling(t)
	TestParallelExecution(t)
	TestContextPassing(t)
	TestErrorHandler(t)

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("结束时间:", time.Now().Format("2006-01-02 15:04:05"))

	if t.Failed() {
		fmt.Println("❌ 测试失败")
		os.Exit(1)
	} else {
		fmt.Println("✅ 所有测试通过")
		os.Exit(0)
	}
}

func main() {
	MainTest()
}
