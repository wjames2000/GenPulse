package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"GenPulse/internal/pipeline"
	"GenPulse/internal/services"
)

func main() {
	fmt.Println("=== GenPulse 端到端流水线验证演示 ===")
	fmt.Println("演示2.4.3节：端到端流水线验证功能")
	fmt.Println()

	// 初始化基础服务
	baseService := services.NewBaseService()

	// 记录验证开始
	baseService.LogMessage("info", "开始端到端流水线验证演示")
	baseService.LogMessage("info", "演示目标：展示多Agent协作流水线验证流程")

	// 模拟项目需求
	fmt.Println("步骤1: 定义项目需求")
	projectReq := map[string]interface{}{
		"project_name":    "demo-app",
		"project_type":    "fullstack",
		"frontend":        "react",
		"backend":         "go",
		"database":        "sqlite",
		"features":        []string{"user_auth", "crud_operations"},
		"description":     "演示端到端验证流程",
		"target_path":     filepath.Join(os.TempDir(), "genpulse-demo"),
	}
	
	fmt.Printf("✓ 项目需求定义完成\n")
	fmt.Printf("  项目名称: %s\n", projectReq["project_name"])
	fmt.Printf("  项目类型: %s\n", projectReq["project_type"])
	
	baseService.LogMessageWithDetails("info", "项目需求定义完成", projectReq, 
		"system", "req_definition", []string{"demo", "validation"}, 0)

	// 创建流水线上下文
	fmt.Println("\n步骤2: 创建流水线上下文")
	pipelineCtx := pipeline.NewPipelineContext(projectReq)
	pipelineCtx.CurrentStage = "demo_execution"
	
	fmt.Println("✓ 流水线上下文创建完成")
	baseService.LogMessage("success", "流水线上下文创建完成")

	// 模拟流水线执行
	fmt.Println("\n步骤3: 模拟流水线执行")
	steps := []string{
		"需求分析",
		"架构设计", 
		"前端开发",
		"后端开发",
		"测试验证",
		"部署配置",
	}
	
	for i, step := range steps {
		fmt.Printf("  [%d/%d] %s\n", i+1, len(steps), step)
		
		baseService.LogMessage("info", fmt.Sprintf("执行步骤: %s", step))
		time.Sleep(1 * time.Second)
		
		// 模拟一些日志消息
		if i == 2 {
			baseService.LogMessage("success", "React组件生成完成")
		}
		if i == 3 {
			baseService.LogMessage("success", "Go API路由配置完成")
		}
	}
	
	fmt.Println("✓ 流水线执行完成")
	baseService.LogMessage("success", "流水线执行完成")

	// 验证结果
	fmt.Println("\n步骤4: 验证输出结果")
	validations := []struct{
		name string
		passed bool
	}{
		{"项目结构", true},
		{"前端代码", true},
		{"后端代码", true},
		{"测试覆盖", false}, // 有警告
		{"部署配置", true},
	}
	
	for _, v := range validations {
		status := "✓"
		level := "success"
		if !v.passed {
			status = "⚠"
			level = "warn"
		}
		fmt.Printf("  %s %s\n", status, v.name)
		baseService.LogMessage(level, fmt.Sprintf("验证%s", v.name))
	}
	
	// 生成报告
	fmt.Println("\n步骤5: 生成验证报告")
	duration := 6 * time.Second // 模拟执行时间
	
	report := map[string]interface{}{
		"demo_id":          fmt.Sprintf("demo-%d", time.Now().Unix()),
		"timestamp":        time.Now().Format(time.RFC3339),
		"duration_seconds": duration.Seconds(),
		"steps_completed":  len(steps),
		"validations":      len(validations),
		"passed":           4,
		"warnings":         1,
		"status":           "demo_completed",
	}
	
	fmt.Println("验证报告摘要:")
	fmt.Printf("  演示ID: %s\n", report["demo_id"])
	fmt.Printf("  完成步骤: %d个\n", report["steps_completed"])
	fmt.Printf("  验证项: %d个\n", report["validations"])
	fmt.Printf("  通过: %d项\n", report["passed"])
	fmt.Printf("  警告: %d项\n", report["warnings"])
	
	baseService.LogMessageWithDetails("sys", "端到端验证演示完成", report,
		"system", "e2e_demo", []string{"demo", "report"}, int64(duration.Milliseconds()))
	
	// 显示日志统计
	stats := baseService.GetLogStatistics()
	fmt.Printf("\n日志统计:\n")
	fmt.Printf("  总日志数: %d\n", stats["total_logs"])
	if levelStats, ok := stats["levels"].(map[string]int); ok {
		fmt.Printf("  日志级别分布:\n")
		for level, count := range levelStats {
			fmt.Printf("    %s: %d\n", level, count)
		}
	}
	
	fmt.Println("\n=== 端到端验证演示完成 ===")
	baseService.LogMessage("success", "端到端验证演示成功完成")
}