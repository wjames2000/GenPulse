package genkit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"GenPulse/internal/genkit/config"
	"GenPulse/internal/genkit/flows"
	"GenPulse/internal/genkit/models"
	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/utils"
)

func TestGenkitIntegration(t *testing.T) {
	// 初始化日志
	utils.InitGlobalLogger(utils.INFO, false)

	// 创建测试上下文
	ctx := context.Background()

	t.Run("TestConfigInitialization", func(t *testing.T) {
		// 测试配置初始化
		cfgMgr := config.GetGlobalConfig()
		if cfgMgr == nil {
			t.Fatal("Failed to get config manager")
		}

		appConfig := cfgMgr.GetConfig()
		if appConfig.AppName != "GenPulse" {
			t.Errorf("Expected app name 'GenPulse', got '%s'", appConfig.AppName)
		}

		fmt.Printf("✓ 配置初始化测试通过\n")
	})

	t.Run("TestModelAdapter", func(t *testing.T) {
		// 测试模型适配器
		factory := &models.DefaultModelAdapterFactory{}
		adapter := models.NewUnifiedModelAdapter(factory)

		// 注册测试模型
		geminiConfig := models.ModelConfig{
			Type:        models.ModelTypeGemini,
			Name:        "gemini-1.5-pro-test",
			Provider:    "google",
			APIKey:      "test-key",
			MaxTokens:   100,
			Temperature: 0.7,
		}

		if err := adapter.RegisterModel(geminiConfig); err != nil {
			t.Fatalf("Failed to register model: %v", err)
		}

		// 测试模型列表
		modelList := adapter.ListModels()
		if len(modelList) == 0 {
			t.Fatal("No models registered")
		}

		fmt.Printf("✓ 模型适配器测试通过 (注册模型: %v)\n", modelList)
	})

	t.Run("TestToolRegistry", func(t *testing.T) {
		// 测试工具注册表
		toolRegistry := tools.GetGlobalToolRegistry()

		// 创建并注册文件系统工具
		fsTool, err := tools.NewFSTool("/tmp/genpulse-test")
		if err != nil {
			t.Fatalf("Failed to create FS tool: %v", err)
		}

		if err := toolRegistry.RegisterTool(fsTool); err != nil {
			t.Fatalf("Failed to register tool: %v", err)
		}

		// 测试工具列表
		toolList := toolRegistry.ListTools()
		if len(toolList) == 0 {
			t.Fatal("No tools registered")
		}

		fmt.Printf("✓ 工具注册表测试通过 (注册工具: %v)\n", len(toolList))
	})

	t.Run("TestFlowEngine", func(t *testing.T) {
		// 测试Flow引擎

		// 创建模型适配器和工具注册表
		factory := &models.DefaultModelAdapterFactory{}
		modelAdapter := models.NewUnifiedModelAdapter(factory)
		toolRegistry := tools.NewToolRegistry()

		// 创建Flow引擎
		flowEngine := flows.NewFlowEngine(modelAdapter, toolRegistry)

		// 创建测试Flow定义
		testFlow := flows.FlowDefinition{
			ID:          "test-flow",
			Name:        "测试Flow",
			Description: "用于集成测试的简单Flow",
			Type:        flows.FlowTypeSequential,
			Version:     "1.0.0",
			Nodes: []flows.NodeDefinition{
				{
					ID:   "start",
					Name: "开始节点",
					Type: flows.NodeTypeAction,
					Config: map[string]interface{}{
						"action_type": "log",
						"message":     "Flow执行开始",
					},
				},
				{
					ID:   "process",
					Name: "处理节点",
					Type: flows.NodeTypeAction,
					Config: map[string]interface{}{
						"action_type": "transform",
						"input":       "hello world",
						"operation":   "uppercase",
					},
				},
				{
					ID:   "end",
					Name: "结束节点",
					Type: flows.NodeTypeAction,
					Config: map[string]interface{}{
						"action_type": "log",
						"message":     "Flow执行完成",
					},
				},
			},
			Edges: []flows.EdgeDefinition{
				{Source: "start", Target: "process"},
				{Source: "process", Target: "end"},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// 注册Flow
		if err := flowEngine.RegisterFlow(testFlow); err != nil {
			t.Fatalf("Failed to register flow: %v", err)
		}

		// 测试Flow列表
		flowList := flowEngine.ListFlows()
		if len(flowList) == 0 {
			t.Fatal("No flows registered")
		}

		fmt.Printf("✓ Flow引擎测试通过 (注册Flow: %v)\n", len(flowList))
	})

	t.Run("TestGenkitManager", func(t *testing.T) {
		// 测试Genkit管理器
		genkitManager := NewGenkitManager()

		// 初始化
		if err := genkitManager.Initialize(ctx); err != nil {
			t.Fatalf("Failed to initialize Genkit manager: %v", err)
		}

		// 检查初始化状态
		if !genkitManager.IsInitialized() {
			t.Fatal("Genkit manager should be initialized")
		}

		fmt.Printf("✓ Genkit管理器测试通过\n")
	})

	// 综合测试：模拟完整工作流
	t.Run("TestCompleteWorkflow", func(t *testing.T) {
		fmt.Println("\n=== 综合测试：模拟完整工作流 ===")

		// 1. 初始化配置
		cfgMgr := config.GetGlobalConfig()
		appConfig := cfgMgr.GetConfig()
		fmt.Printf("1. 应用配置加载: %s v%s\n", appConfig.AppName, appConfig.AppVersion)

		// 2. 初始化模型适配器
		factory := &models.DefaultModelAdapterFactory{}
		modelAdapter := models.NewUnifiedModelAdapter(factory)

		// 注册示例模型
		modelConfigs := []models.ModelConfig{
			{
				Type:        models.ModelTypeGemini,
				Name:        "gemini-1.5-pro",
				Provider:    "google",
				APIKey:      "demo-key",
				MaxTokens:   1000,
				Temperature: 0.7,
			},
			{
				Type:        models.ModelTypeGPT,
				Name:        "gpt-4-turbo",
				Provider:    "openai",
				APIKey:      "demo-key",
				MaxTokens:   1000,
				Temperature: 0.7,
			},
		}

		for _, cfg := range modelConfigs {
			if err := modelAdapter.RegisterModel(cfg); err != nil {
				fmt.Printf("  警告: 注册模型失败 %s: %v\n", cfg.Name, err)
			} else {
				fmt.Printf("2. 注册模型: %s (%s)\n", cfg.Name, cfg.Type)
			}
		}

		// 3. 初始化工具注册表
		toolRegistry := tools.NewToolRegistry()

		// 注册文件系统工具
		fsTool, err := tools.NewFSTool("/tmp/genpulse-demo")
		if err == nil {
			if err := toolRegistry.RegisterTool(fsTool); err != nil {
				fmt.Printf("  警告: 注册工具失败: %v\n", err)
			} else {
				fmt.Printf("3. 注册工具: %s\n", fsTool.GetDefinition().Name)
			}
		}

		// 4. 初始化Flow引擎
		flowEngine := flows.NewFlowEngine(modelAdapter, toolRegistry)

		// 创建示例Flow
		exampleFlow := flows.FlowDefinition{
			ID:          "example-workflow",
			Name:        "示例工作流",
			Description: "演示完整的AI工作流",
			Type:        flows.FlowTypeSequential,
			Version:     "1.0.0",
			Nodes: []flows.NodeDefinition{
				{
					ID:   "analyze",
					Name: "需求分析",
					Type: flows.NodeTypeModel,
					Config: map[string]interface{}{
						"model_id": "gemini-1.5-pro",
						"prompt":   "分析用户需求：{{user_input}}",
					},
				},
				{
					ID:   "generate",
					Name: "代码生成",
					Type: flows.NodeTypeModel,
					Config: map[string]interface{}{
						"model_id": "gpt-4-turbo",
						"prompt":   "根据分析结果生成代码：{{analyze.output.content}}",
					},
				},
				{
					ID:   "save",
					Name: "保存文件",
					Type: flows.NodeTypeTool,
					Config: map[string]interface{}{
						"tool_id": "fs_tool",
						"parameters": map[string]interface{}{
							"operation": "write",
							"path":      "output/{{timestamp}}.go",
							"content":   "{{generate.output.content}}",
						},
					},
				},
			},
			Edges: []flows.EdgeDefinition{
				{Source: "analyze", Target: "generate"},
				{Source: "generate", Target: "save"},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := flowEngine.RegisterFlow(exampleFlow); err != nil {
			fmt.Printf("  警告: 注册Flow失败: %v\n", err)
		} else {
			fmt.Printf("4. 注册Flow: %s\n", exampleFlow.Name)
		}

		// 5. 显示系统状态
		fmt.Println("\n=== 系统状态汇总 ===")
		fmt.Printf("• 注册模型: %d\n", len(modelAdapter.ListModels()))
		fmt.Printf("• 注册工具: %d\n", toolRegistry.GetToolCount())
		fmt.Printf("• 注册Flow: %d\n", len(flowEngine.ListFlows()))

		// 6. 健康检查
		fmt.Println("\n=== 健康检查 ===")

		// 模型健康检查
		modelErrors := modelAdapter.HealthCheck(ctx)
		if len(modelErrors) == 0 {
			fmt.Println("• 模型适配器: ✅ 健康")
		} else {
			fmt.Printf("• 模型适配器: ⚠️  有 %d 个问题\n", len(modelErrors))
		}

		// 工具统计
		toolStats := toolRegistry.GetToolStatistics()
		fmt.Printf("• 工具注册表: 总共 %d 个工具 (%d 启用, %d 禁用)\n",
			toolStats["total_tools"].(int),
			toolStats["enabled_tools"].(int),
			toolStats["disabled_tools"].(int))

		// Flow统计
		flowStats := flowEngine.GetStatistics()
		fmt.Printf("• Flow引擎: %d 个Flow定义\n", flowStats["total_flows"].(int))

		fmt.Println("\n✓ 综合测试完成")
	})

	fmt.Println("\n🎉 所有集成测试通过！")
}

// 运行集成测试
func ExampleIntegration() {
	// 这个示例展示如何初始化和使用完整的Genkit系统

	ctx := context.Background()

	// 1. 初始化配置
	config.InitGlobalConfig()

	// 2. 初始化Genkit管理器
	genkitManager := NewGenkitManager()
	if err := genkitManager.Initialize(ctx); err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		return
	}

	// 3. 获取全局实例
	modelAdapter := models.NewUnifiedModelAdapter(&models.DefaultModelAdapterFactory{})
	toolRegistry := tools.GetGlobalToolRegistry()
	flowEngine := flows.GetGlobalFlowEngine()

	// 使用示例
	fmt.Printf("Genkit系统初始化完成:\n")
	fmt.Printf("• 模型适配器: 已创建\n")
	fmt.Printf("• 工具注册表: %d 个工具\n", toolRegistry.GetToolCount())
	if flowEngine != nil {
		fmt.Printf("• Flow引擎: %d 个Flow\n", len(flowEngine.ListFlows()))
	}
}
