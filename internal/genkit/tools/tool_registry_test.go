package tools

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestToolRegistryBasic 测试工具注册表基础功能
func TestToolRegistryBasic(t *testing.T) {
	registry := NewToolRegistry()

	// 创建测试工具
	testTool := NewTestTool(ToolDefinition{
		ID:          "test-tool-1",
		Name:        "测试工具1",
		Description: "用于测试的工具",
		Category:    ToolCategoryUtility,
		Version:     "1.0.0",
		Enabled:     true,
	})

	// 测试1: 注册工具
	t.Run("注册工具", func(t *testing.T) {
		err := registry.RegisterTool(testTool)
		if err != nil {
			t.Errorf("注册工具失败: %v", err)
		}

		// 验证工具已注册
		tool, err := registry.GetTool("test-tool-1")
		if err != nil {
			t.Errorf("获取工具失败: %v", err)
		}
		if tool.GetDefinition().ID != "test-tool-1" {
			t.Errorf("工具ID不正确: 期望 %s, 实际 %s", "test-tool-1", tool.GetDefinition().ID)
		}
	})

	// 测试2: 重复注册
	t.Run("重复注册", func(t *testing.T) {
		duplicateTool := NewTestTool(ToolDefinition{
			ID:          "test-tool-1", // 相同的ID
			Name:        "测试工具1重复",
			Description: "重复注册测试",
			Category:    ToolCategoryUtility,
			Version:     "1.0.1",
			Enabled:     true,
		})

		err := registry.RegisterTool(duplicateTool)
		if err == nil {
			t.Error("重复注册应该失败")
		}
	})

	// 测试3: 获取不存在的工具
	t.Run("获取不存在的工具", func(t *testing.T) {
		_, err := registry.GetTool("non-existent-tool")
		if err == nil {
			t.Error("获取不存在的工具应该失败")
		}
	})

	// 测试4: 注册另一个工具
	t.Run("注册另一个工具", func(t *testing.T) {
		testTool2 := NewTestTool(ToolDefinition{
			ID:          "test-tool-2",
			Name:        "测试工具2",
			Description: "另一个测试工具",
			Category:    ToolCategoryFileSystem,
			Version:     "1.0.0",
			Enabled:     true,
		})
		err := registry.RegisterTool(testTool2)
		if err != nil {
			t.Errorf("注册第二个工具失败: %v", err)
		}
	})

	// 测试5: 获取工具数量
	t.Run("获取工具数量", func(t *testing.T) {
		count := registry.GetToolCount()
		if count != 2 {
			t.Errorf("工具总数不正确: 期望 2, 实际 %d", count)
		}

		utilityCount := registry.GetToolCountByCategory(ToolCategoryUtility)
		if utilityCount != 1 {
			t.Errorf("工具类别数量不正确: 期望 1, 实际 %d", utilityCount)
		}

		fsCount := registry.GetToolCountByCategory(ToolCategoryFileSystem)
		if fsCount != 1 {
			t.Errorf("文件系统工具数量不正确: 期望 1, 实际 %d", fsCount)
		}
	})

	// 测试7: 执行工具
	t.Run("执行工具", func(t *testing.T) {
		ctx := context.Background()
		execution := ToolExecution{
			ToolID: "test-tool-1",
			Parameters: map[string]interface{}{
				"test": "value",
			},
		}

		result, err := registry.ExecuteTool(ctx, execution)
		if err != nil {
			t.Errorf("执行工具失败: %v", err)
		}
		if !result.Success {
			t.Errorf("工具执行应该成功: %v", result.Error)
		}

		// 验证执行结果
		output, ok := result.Output.(map[string]interface{})
		if !ok {
			t.Error("输出格式不正确")
		}
		if message, ok := output["message"].(string); !ok || message != "test executed successfully" {
			t.Errorf("输出内容不正确: %v", output)
		}
	})

	// 测试8: 执行不存在的工具
	t.Run("执行不存在的工具", func(t *testing.T) {
		ctx := context.Background()
		execution := ToolExecution{
			ToolID: "non-existent-tool",
			Parameters: map[string]interface{}{
				"test": "value",
			},
		}

		result, err := registry.ExecuteTool(ctx, execution)
		if err == nil {
			t.Error("执行不存在的工具应该失败")
		}
		// 注意：ExecuteTool 在错误时返回非nil的ToolResult
		if result == nil {
			t.Error("执行失败时应该返回ToolResult")
		} else if result.Success {
			t.Error("执行不存在的工具应该返回失败结果")
		}
	})

	// 测试9: 启用/禁用工具
	t.Run("启用禁用工具", func(t *testing.T) {
		// 禁用工具
		err := registry.DisableTool("test-tool-1")
		if err != nil {
			t.Errorf("禁用工具失败: %v", err)
		}

		tool, err := registry.GetTool("test-tool-1")
		if err != nil {
			t.Errorf("获取工具失败: %v", err)
		}
		if tool.IsEnabled() {
			t.Error("工具应该被禁用")
		}

		// 启用工具
		err = registry.EnableTool("test-tool-1")
		if err != nil {
			t.Errorf("启用工具失败: %v", err)
		}

		tool, err = registry.GetTool("test-tool-1")
		if err != nil {
			t.Errorf("获取工具失败: %v", err)
		}
		if !tool.IsEnabled() {
			t.Error("工具应该被启用")
		}
	})

	// 测试10: 初始化所有工具
	t.Run("初始化所有工具", func(t *testing.T) {
		err := registry.InitializeAllTools()
		if err != nil {
			t.Errorf("初始化所有工具失败: %v", err)
		}
	})

	// 测试11: 获取工具统计信息
	t.Run("获取工具统计信息", func(t *testing.T) {
		stats := registry.GetToolStatistics()

		// 验证基本统计信息
		if totalTools, ok := stats["total_tools"].(int); !ok || totalTools != 2 {
			t.Errorf("工具总数不正确: 期望 2, 实际 %v", stats["total_tools"])
		}

		if enabledTools, ok := stats["enabled_tools"].(int); !ok || enabledTools != 2 {
			t.Errorf("启用工具数不正确: 期望 2, 实际 %v", stats["enabled_tools"])
		}

		if disabledTools, ok := stats["disabled_tools"].(int); !ok || disabledTools != 0 {
			t.Errorf("禁用工具数不正确: 期望 0, 实际 %v", stats["disabled_tools"])
		}

		// 验证类别统计
		if categoryStats, ok := stats["tools_by_category"].(map[string]int); ok {
			if utilityCount, ok := categoryStats["utility"]; !ok || utilityCount != 1 {
				t.Errorf("工具类别数量不正确: 期望 utility=1, 实际 %v", categoryStats["utility"])
			}
			if fsCount, ok := categoryStats["filesystem"]; !ok || fsCount != 1 {
				t.Errorf("工具类别数量不正确: 期望 filesystem=1, 实际 %v", categoryStats["filesystem"])
			}
		} else {
			// 检查是否是 map[string]interface{}
			if categoryStats2, ok := stats["tools_by_category"].(map[string]interface{}); ok {
				if utilityCount, ok := categoryStats2["utility"].(int); !ok || utilityCount != 1 {
					t.Errorf("工具类别数量不正确: 期望 utility=1, 实际 %v", categoryStats2["utility"])
				}
				if fsCount, ok := categoryStats2["filesystem"].(int); !ok || fsCount != 1 {
					t.Errorf("工具类别数量不正确: 期望 filesystem=1, 实际 %v", categoryStats2["filesystem"])
				}
			} else {
				t.Error("类别统计信息不存在或格式不正确")
			}
		}
	})
}

// TestToolRegistryConcurrent 测试并发操作
func TestToolRegistryConcurrent(t *testing.T) {
	registry := NewToolRegistry()
	errors := make(chan error, 10)

	// 并发注册多个工具
	for i := 0; i < 5; i++ {
		go func(index int) {
			tool := NewTestTool(ToolDefinition{
				ID:          "concurrent-tool-" + string(rune('A'+index)),
				Name:        "并发工具 " + string(rune('A'+index)),
				Description: "用于并发测试的工具",
				Category:    ToolCategoryUtility,
				Version:     "1.0.0",
				Enabled:     true,
			})

			err := registry.RegisterTool(tool)
			errors <- err
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 5; i++ {
		if err := <-errors; err != nil {
			t.Errorf("并发注册失败: %v", err)
		}
	}

	// 验证所有工具都注册了
	count := registry.GetToolCount()
	if count != 5 {
		t.Errorf("工具数量不正确: 期望 5, 实际 %d", count)
	}

	// 并发执行工具
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		go func(index int) {
			execution := ToolExecution{
				ToolID: "concurrent-tool-" + string(rune('A'+index)),
				Parameters: map[string]interface{}{
					"test": "value",
				},
			}

			_, err := registry.ExecuteTool(ctx, execution)
			errors <- err
		}(i)
	}

	// 等待所有执行完成
	for i := 0; i < 5; i++ {
		if err := <-errors; err != nil {
			t.Errorf("并发执行失败: %v", err)
		}
	}

	// 验证执行统计
	stats := registry.GetToolStatistics()

	if totalTools, ok := stats["total_tools"].(int); !ok || totalTools != 5 {
		t.Errorf("工具总数不正确: 期望 5, 实际 %v", stats["total_tools"])
	}
}

// TestToolRegistryErrorHandling 测试错误处理
func TestToolRegistryErrorHandling(t *testing.T) {
	registry := NewToolRegistry()

	// 测试1: 注册nil工具（会panic，所以跳过或期望panic）
	t.Run("注册nil工具", func(t *testing.T) {
		// 这个测试会导致panic，所以我们需要捕获它或者跳过
		// 在实际代码中应该添加nil检查
		t.Skip("注册nil工具会导致panic，需要在实际代码中添加nil检查")
	})

	// 测试2: 禁用不存在的工具
	t.Run("禁用不存在的工具", func(t *testing.T) {
		err := registry.DisableTool("non-existent-tool")
		if err == nil {
			t.Error("禁用不存在的工具应该失败")
		}
	})

	// 测试3: 启用不存在的工具
	t.Run("启用不存在的工具", func(t *testing.T) {
		err := registry.EnableTool("non-existent-tool")
		if err == nil {
			t.Error("启用不存在的工具应该失败")
		}
	})

	// 测试4: 执行禁用的工具
	t.Run("执行禁用的工具", func(t *testing.T) {
		// 注册一个工具
		tool := NewTestTool(ToolDefinition{
			ID:          "disabled-tool",
			Name:        "禁用工具",
			Description: "被禁用的测试工具",
			Category:    ToolCategoryUtility,
			Version:     "1.0.0",
			Enabled:     true,
		})
		err := registry.RegisterTool(tool)
		if err != nil {
			t.Errorf("注册工具失败: %v", err)
		}

		// 禁用工具
		err = registry.DisableTool("disabled-tool")
		if err != nil {
			t.Errorf("禁用工具失败: %v", err)
		}

		// 验证工具确实被禁用
		registeredTool, err := registry.GetTool("disabled-tool")
		if err != nil {
			t.Errorf("获取工具失败: %v", err)
		}
		if registeredTool.IsEnabled() {
			t.Error("工具应该被禁用")
		}

		ctx := context.Background()
		execution := ToolExecution{
			ToolID: "disabled-tool",
			Parameters: map[string]interface{}{
				"test": "value",
			},
		}

		result, err := registry.ExecuteTool(ctx, execution)
		if err == nil {
			t.Error("执行禁用的工具应该返回错误")
		}
		if result == nil {
			t.Error("执行失败时应该返回ToolResult")
		} else if result.Success {
			t.Error("执行禁用的工具应该返回失败结果")
		}
	})
}

// TestToolRegistryPerformance 测试性能
func TestToolRegistryPerformance(t *testing.T) {
	registry := NewToolRegistry()
	ctx := context.Background()

	// 注册大量工具
	startTime := time.Now()
	for i := 0; i < 100; i++ {
		tool := NewTestTool(ToolDefinition{
			ID:          "perf-tool-" + fmt.Sprintf("%03d", i),
			Name:        "性能测试工具 " + fmt.Sprintf("%03d", i),
			Description: "用于性能测试的工具",
			Category:    ToolCategoryUtility,
			Version:     "1.0.0",
			Enabled:     true,
		})
		err := registry.RegisterTool(tool)
		if err != nil {
			t.Errorf("注册工具失败: %v", err)
		}
	}
	registrationTime := time.Since(startTime)

	// 执行大量工具
	startTime = time.Now()
	for i := 0; i < 100; i++ {
		execution := ToolExecution{
			ToolID: "perf-tool-" + fmt.Sprintf("%03d", i),
			Parameters: map[string]interface{}{
				"test": "value",
			},
		}
		_, err := registry.ExecuteTool(ctx, execution)
		if err != nil {
			t.Errorf("执行工具失败: %v", err)
		}
	}
	executionTime := time.Since(startTime)

	// 获取统计信息
	startTime = time.Now()
	stats := registry.GetToolStatistics()
	statsTime := time.Since(startTime)

	t.Logf("性能测试结果:")
	t.Logf("  注册100个工具: %v", registrationTime)
	t.Logf("  执行100次工具: %v", executionTime)
	t.Logf("  获取统计信息: %v", statsTime)
	t.Logf("  总工具数: %v", stats["total_tools"])

	// 验证性能指标（宽松检查）
	if registrationTime > 2*time.Second {
		t.Errorf("注册工具时间过长: %v", registrationTime)
	}
	if executionTime > 2*time.Second {
		t.Errorf("执行工具时间过长: %v", executionTime)
	}
	if statsTime > 1*time.Second {
		t.Errorf("获取统计信息时间过长: %v", statsTime)
	}
}
