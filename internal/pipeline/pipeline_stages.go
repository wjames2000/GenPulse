package pipeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"GenPulse/internal/agents"
	"GenPulse/internal/utils"
)

// executeParallelDevelopment 执行并行开发
func (pf *PipelineFlow) executeParallelDevelopment(ctx context.Context, pipelineCtx *PipelineContext) error {
	pipelineCtx.AddLog("info", "开始并行开发阶段", nil)

	// 创建并行执行引擎
	parallelEngine := NewParallelEngine(pf.agentManager, 2) // 2个worker用于前端和后端

	// 启动引擎
	if err := parallelEngine.Start(ctx); err != nil {
		return fmt.Errorf("启动并行引擎失败: %w", err)
	}
	defer parallelEngine.Stop()

	// 创建并行任务
	tasks := parallelEngine.CreateFrontendBackendTasks(pipelineCtx)

	// 执行并行任务
	results, err := parallelEngine.ExecuteParallel(ctx, tasks)
	if err != nil {
		return fmt.Errorf("并行执行失败: %w", err)
	}

	// 处理结果
	var wg sync.WaitGroup
	errors := make(chan error, len(results))

	for _, result := range results {
		wg.Add(1)
		go func(result ParallelTaskResult) {
			defer wg.Done()

			if !result.Success {
				errors <- fmt.Errorf("%s 任务失败: %v", result.AgentName, result.Error)
				return
			}

			// 根据Agent角色保存结果
			switch result.AgentRole {
			case string(agents.RoleFrontendDev):
				pipelineCtx.SetArtifact("frontend_code", result.Output)
				pipelineCtx.SetArtifact("frontend_artifacts", result.Artifacts)
				pipelineCtx.AddLog("success", "前端开发完成", map[string]interface{}{
					"agent":    result.AgentName,
					"duration": result.Duration.String(),
				})

			case string(agents.RoleBackendDev):
				pipelineCtx.SetArtifact("backend_code", result.Output)
				pipelineCtx.SetArtifact("backend_artifacts", result.Artifacts)
				pipelineCtx.AddLog("success", "后端开发完成", map[string]interface{}{
					"agent":    result.AgentName,
					"duration": result.Duration.String(),
				})
			}
		}(result)
	}

	wg.Wait()
	close(errors)

	// 检查是否有错误
	var combinedErr error
	for err := range errors {
		if combinedErr == nil {
			combinedErr = err
		} else {
			combinedErr = fmt.Errorf("%v; %v", combinedErr, err)
		}
	}

	if combinedErr != nil {
		return fmt.Errorf("并行开发阶段有任务失败: %w", combinedErr)
	}

	pipelineCtx.AddLog("success", "并行开发阶段完成", map[string]interface{}{
		"completed_tasks":  len(results),
		"successful_tasks": countSuccessfulTasks(results),
	})

	return nil
}

// executeTesting 执行测试
func (pf *PipelineFlow) executeTesting(ctx context.Context, pipelineCtx *PipelineContext) error {
	pipelineCtx.AddLog("info", "开始测试阶段", nil)

	// 获取QA工程师Agent
	agent, err := pf.agentManager.GetAgent("qa_engineer_001")
	if err != nil {
		return fmt.Errorf("获取QA工程师Agent失败: %w", err)
	}

	// 准备参数
	params := map[string]interface{}{
		"project_name":        pipelineCtx.Parameters["project_name"],
		"project_description": pipelineCtx.Parameters["project_description"],
		"prd_document":        pipelineCtx.GetArtifact("prd_document", ""),
		"architecture_design": pipelineCtx.GetArtifact("architecture_design", ""),
		"frontend_code":       pipelineCtx.GetArtifact("frontend_code", ""),
		"backend_code":        pipelineCtx.GetArtifact("backend_code", ""),
		"task_plan":           pipelineCtx.GetArtifact("task_plan", ""),
	}

	// 执行测试
	task := "对生成的项目进行测试"
	result, err := agent.Execute(ctx, task, params)
	if err != nil {
		return fmt.Errorf("测试执行失败: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("测试失败: %v", result.Output)
	}

	// 保存产物
	pipelineCtx.SetArtifact("test_report", result.Output)
	pipelineCtx.SetArtifact("test_cases", result.Artifacts)

	pipelineCtx.AddLog("success", "测试完成", map[string]interface{}{
		"agent":    agent.GetConfig().Name,
		"duration": result.Duration.String(),
	})

	return nil
}

// executeDeployment 执行部署
func (pf *PipelineFlow) executeDeployment(ctx context.Context, pipelineCtx *PipelineContext) error {
	pipelineCtx.AddLog("info", "开始部署阶段", nil)

	// 获取DevOps工程师Agent
	agent, err := pf.agentManager.GetAgent("devops_engineer_001")
	if err != nil {
		return fmt.Errorf("获取DevOps工程师Agent失败: %w", err)
	}

	// 准备参数
	params := map[string]interface{}{
		"project_name":        pipelineCtx.Parameters["project_name"],
		"project_description": pipelineCtx.Parameters["project_description"],
		"architecture_design": pipelineCtx.GetArtifact("architecture_design", ""),
		"frontend_code":       pipelineCtx.GetArtifact("frontend_code", ""),
		"backend_code":        pipelineCtx.GetArtifact("backend_code", ""),
		"test_report":         pipelineCtx.GetArtifact("test_report", ""),
		"tech_stack":          pipelineCtx.Parameters["tech_stack"],
	}

	// 执行部署
	task := "部署生成的项目"
	result, err := agent.Execute(ctx, task, params)
	if err != nil {
		return fmt.Errorf("部署执行失败: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("部署失败: %v", result.Output)
	}

	// 保存产物
	pipelineCtx.SetArtifact("deployment_result", result.Output)
	pipelineCtx.SetArtifact("deployment_artifacts", result.Artifacts)

	// 提取项目路径（如果部署结果中包含）
	if deploymentMap, ok := result.Output.(map[string]interface{}); ok {
		if projectPath, ok := deploymentMap["project_path"]; ok {
			pipelineCtx.SetArtifact("project_path", projectPath)
		}
	}

	pipelineCtx.AddLog("success", "部署完成", map[string]interface{}{
		"agent":    agent.GetConfig().Name,
		"duration": result.Duration.String(),
	})

	return nil
}

// executeCodeReview 执行代码审查
func (pf *PipelineFlow) executeCodeReview(ctx context.Context, pipelineCtx *PipelineContext) error {
	pipelineCtx.AddLog("info", "开始代码审查阶段", nil)

	// 获取代码审查员Agent
	agent, err := pf.agentManager.GetAgent("reviewer_001")
	if err != nil {
		return fmt.Errorf("获取代码审查员Agent失败: %w", err)
	}

	// 准备参数
	params := map[string]interface{}{
		"project_name":        pipelineCtx.Parameters["project_name"],
		"project_description": pipelineCtx.Parameters["project_description"],
		"frontend_code":       pipelineCtx.GetArtifact("frontend_code", ""),
		"backend_code":        pipelineCtx.GetArtifact("backend_code", ""),
		"architecture_design": pipelineCtx.GetArtifact("architecture_design", ""),
		"test_report":         pipelineCtx.GetArtifact("test_report", ""),
	}

	// 执行代码审查
	task := "审查生成的代码"
	result, err := agent.Execute(ctx, task, params)
	if err != nil {
		return fmt.Errorf("代码审查执行失败: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("代码审查失败: %v", result.Output)
	}

	// 保存产物
	pipelineCtx.SetArtifact("code_review_report", result.Output)
	pipelineCtx.SetArtifact("review_issues", result.Artifacts)

	pipelineCtx.AddLog("success", "代码审查完成", map[string]interface{}{
		"agent":    agent.GetConfig().Name,
		"duration": result.Duration.String(),
	})

	return nil
}

// executeProjectValidation 执行项目验证
func (pf *PipelineFlow) executeProjectValidation(ctx context.Context, pipelineCtx *PipelineContext) error {
	pipelineCtx.AddLog("info", "开始项目验证阶段", nil)

	// 获取DevOps工程师Agent进行验证
	agent, err := pf.agentManager.GetAgent("devops_engineer_001")
	if err != nil {
		return fmt.Errorf("获取DevOps工程师Agent失败: %w", err)
	}

	// 准备参数
	params := map[string]interface{}{
		"project_name":        pipelineCtx.Parameters["project_name"],
		"project_description": pipelineCtx.Parameters["project_description"],
		"deployment_result":   pipelineCtx.GetArtifact("deployment_result", ""),
		"code_review_report":  pipelineCtx.GetArtifact("code_review_report", ""),
		"test_report":         pipelineCtx.GetArtifact("test_report", ""),
		"project_path":        pipelineCtx.GetArtifact("project_path", ""),
	}

	// 执行项目验证
	task := "验证生成的项目是否可运行"
	result, err := agent.Execute(ctx, task, params)
	if err != nil {
		return fmt.Errorf("项目验证执行失败: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("项目验证失败: %v", result.Output)
	}

	// 保存产物
	pipelineCtx.SetArtifact("validation_result", result.Output)
	pipelineCtx.SetArtifact("validation_artifacts", result.Artifacts)

	pipelineCtx.AddLog("success", "项目验证完成", map[string]interface{}{
		"agent":             agent.GetConfig().Name,
		"duration":          result.Duration.String(),
		"project_generated": true,
	})

	return nil
}

// countSuccessfulTasks 计算成功任务数量
func countSuccessfulTasks(results []ParallelTaskResult) int {
	count := 0
	for _, result := range results {
		if result.Success {
			count++
		}
	}
	return count
}

// ExecuteStageWithRetry 带重试执行阶段
func (pf *PipelineFlow) ExecuteStageWithRetry(ctx context.Context, stageName string, executeFunc func(context.Context) error, maxRetries int) error {
	errorHandler := NewErrorHandler(maxRetries, 2*time.Second)

	var lastErr error

	for retry := 0; retry <= maxRetries; retry++ {
		err := executeFunc(ctx)

		if err == nil {
			if retry > 0 {
				utils.Info("阶段 %s 重试成功 (尝试 %d)", stageName, retry+1)
			}
			return nil
		}

		lastErr = err

		// 如果是最后一次尝试，直接返回错误
		if retry == maxRetries {
			break
		}

		// 处理错误，决定是否重试
		shouldRetry, action, waitTime := errorHandler.HandleError(ctx, stageName, "", err, retry)

		if !shouldRetry {
			return fmt.Errorf("%s: %w", action, err)
		}

		// 等待重试
		select {
		case <-time.After(waitTime):
			// 继续重试
			utils.Info("阶段 %s 重试 %d/%d，等待 %v", stageName, retry+1, maxRetries, waitTime)
			continue
		case <-ctx.Done():
			// 上下文取消
			return fmt.Errorf("阶段 %s 执行被取消: %w", stageName, ctx.Err())
		}
	}

	return fmt.Errorf("阶段 %s 重试 %d 次后仍然失败: %w", stageName, maxRetries, lastErr)
}
