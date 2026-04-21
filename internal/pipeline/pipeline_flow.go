package pipeline

import (
	"context"
	"fmt"
	"time"

	"GenPulse/internal/agents"
	"GenPulse/internal/genkit/flows"
	"GenPulse/internal/utils"
)

// PipelineFlow 主流水线Flow
type PipelineFlow struct {
	flowEngine   *flows.FlowEngine
	agentManager *agents.AgentManager
}

// NewPipelineFlow 创建主流水线Flow
func NewPipelineFlow(flowEngine *flows.FlowEngine, agentManager *agents.AgentManager) *PipelineFlow {
	return &PipelineFlow{
		flowEngine:   flowEngine,
		agentManager: agentManager,
	}
}

// DefineMainPipeline 定义主流水线Flow
func (pf *PipelineFlow) DefineMainPipeline() (*flows.FlowDefinition, error) {
	// 主流水线定义
	flowDef := &flows.FlowDefinition{
		ID:          "main_software_development_pipeline",
		Name:        "软件项目开发主流水线",
		Description: "从需求输入到项目完成的完整软件开发流水线，支持多Agent协作",
		Type:        flows.FlowTypeSequential,
		Version:     "1.0.0",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Parameters: map[string]interface{}{
			"project_name": map[string]interface{}{
				"type":        "string",
				"description": "项目名称",
				"required":    true,
			},
			"project_description": map[string]interface{}{
				"type":        "string",
				"description": "项目描述",
				"required":    true,
			},
			"requirements": map[string]interface{}{
				"type":        "string",
				"description": "项目需求",
				"required":    true,
			},
			"tech_stack": map[string]interface{}{
				"type":        "string",
				"description": "技术栈偏好（可选）",
				"required":    false,
			},
		},
		Returns: map[string]interface{}{
			"project_generated": map[string]interface{}{
				"type":        "boolean",
				"description": "项目是否生成成功",
			},
			"project_path": map[string]interface{}{
				"type":        "string",
				"description": "生成的项目路径",
			},
			"execution_summary": map[string]interface{}{
				"type":        "object",
				"description": "执行摘要",
			},
			"artifacts": map[string]interface{}{
				"type":        "array",
				"description": "生成的所有产物",
			},
		},
	}

	// 定义流水线节点
	nodes := []flows.NodeDefinition{
		// 阶段1: 需求分析
		{
			ID:          "node_1_requirements_analysis",
			Name:        "需求分析",
			Type:        flows.NodeTypeAction,
			Description: "产品经理Agent分析需求，生成PRD文档",
			Config: map[string]interface{}{
				"agent_id": "product_manager_001",
				"action":   "analyze_requirements",
			},
			Position: flows.Position{X: 100, Y: 100},
		},
		// 阶段2: 架构设计
		{
			ID:          "node_2_architecture_design",
			Name:        "架构设计",
			Type:        flows.NodeTypeAction,
			Description: "技术架构师Agent设计系统架构",
			Config: map[string]interface{}{
				"agent_id": "architect_001",
				"action":   "design_architecture",
			},
			Position: flows.Position{X: 300, Y: 100},
		},
		// 阶段3: 任务分解
		{
			ID:          "node_3_task_decomposition",
			Name:        "任务分解",
			Type:        flows.NodeTypeAction,
			Description: "编排器Agent分解任务，生成执行计划",
			Config: map[string]interface{}{
				"agent_id": "orchestrator_001",
				"action":   "decompose_tasks",
			},
			Position: flows.Position{X: 500, Y: 100},
		},
		// 阶段4: 并行开发（前端+后端）
		{
			ID:          "node_4_parallel_development",
			Name:        "并行开发",
			Type:        flows.NodeTypeParallel,
			Description: "前端和后端Agent并行开发",
			Config: map[string]interface{}{
				"parallel_nodes": []string{
					"node_4_1_frontend_development",
					"node_4_2_backend_development",
				},
			},
			Position: flows.Position{X: 700, Y: 100},
		},
		// 子节点: 前端开发
		{
			ID:          "node_4_1_frontend_development",
			Name:        "前端开发",
			Type:        flows.NodeTypeAction,
			Description: "前端开发Agent实现用户界面",
			Config: map[string]interface{}{
				"agent_id": "frontend_dev_001",
				"action":   "develop_frontend",
			},
			Position: flows.Position{X: 600, Y: 50},
		},
		// 子节点: 后端开发
		{
			ID:          "node_4_2_backend_development",
			Name:        "后端开发",
			Type:        flows.NodeTypeAction,
			Description: "后端开发Agent实现API和业务逻辑",
			Config: map[string]interface{}{
				"agent_id": "backend_dev_001",
				"action":   "develop_backend",
			},
			Position: flows.Position{X: 600, Y: 150},
		},
		// 阶段5: 测试
		{
			ID:          "node_5_testing",
			Name:        "测试",
			Type:        flows.NodeTypeAction,
			Description: "QA工程师Agent进行测试",
			Config: map[string]interface{}{
				"agent_id": "qa_engineer_001",
				"action":   "execute_tests",
			},
			Position: flows.Position{X: 900, Y: 100},
		},
		// 阶段6: 部署
		{
			ID:          "node_6_deployment",
			Name:        "部署",
			Type:        flows.NodeTypeAction,
			Description: "DevOps工程师Agent进行部署",
			Config: map[string]interface{}{
				"agent_id": "devops_engineer_001",
				"action":   "deploy_project",
			},
			Position: flows.Position{X: 1100, Y: 100},
		},
		// 阶段7: 代码审查
		{
			ID:          "node_7_code_review",
			Name:        "代码审查",
			Type:        flows.NodeTypeAction,
			Description: "代码审查员Agent进行代码审查",
			Config: map[string]interface{}{
				"agent_id": "reviewer_001",
				"action":   "review_code",
			},
			Position: flows.Position{X: 1300, Y: 100},
		},
		// 阶段8: 项目验证
		{
			ID:          "node_8_project_validation",
			Name:        "项目验证",
			Type:        flows.NodeTypeAction,
			Description: "验证生成的项目是否可运行",
			Config: map[string]interface{}{
				"agent_id": "devops_engineer_001",
				"action":   "validate_project",
			},
			Position: flows.Position{X: 1500, Y: 100},
		},
	}

	// 定义边（节点连接）
	edges := []flows.EdgeDefinition{
		// 顺序执行边
		{ID: "edge_1_2", Source: "node_1_requirements_analysis", Target: "node_2_architecture_design"},
		{ID: "edge_2_3", Source: "node_2_architecture_design", Target: "node_3_task_decomposition"},
		{ID: "edge_3_4", Source: "node_3_task_decomposition", Target: "node_4_parallel_development"},
		{ID: "edge_4_5", Source: "node_4_parallel_development", Target: "node_5_testing"},
		{ID: "edge_5_6", Source: "node_5_testing", Target: "node_6_deployment"},
		{ID: "edge_6_7", Source: "node_6_deployment", Target: "node_7_code_review"},
		{ID: "edge_7_8", Source: "node_7_code_review", Target: "node_8_project_validation"},
		// 并行节点到父节点的边
		{ID: "edge_4_1_to_parent", Source: "node_4_1_frontend_development", Target: "node_4_parallel_development"},
		{ID: "edge_4_2_to_parent", Source: "node_4_2_backend_development", Target: "node_4_parallel_development"},
	}

	flowDef.Nodes = nodes
	flowDef.Edges = edges

	return flowDef, nil
}

// ExecutePipeline 执行流水线
func (pf *PipelineFlow) ExecutePipeline(ctx context.Context, parameters map[string]interface{}) (*PipelineResult, error) {
	utils.Info("开始执行软件项目开发流水线")

	// 验证必要参数
	if err := pf.validateParameters(parameters); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}

	// 创建流水线上下文
	pipelineCtx := &PipelineContext{
		Parameters:   parameters,
		Artifacts:    make(map[string]interface{}),
		ExecutionLog: []PipelineLogEntry{},
		StartTime:    time.Now(),
		CurrentStage: "initializing",
	}

	// 记录开始
	pipelineCtx.AddLog("info", "流水线开始执行", nil)

	// 执行各个阶段
	var err error
	var result *PipelineResult

	// 阶段1: 需求分析
	pipelineCtx.CurrentStage = "requirements_analysis"
	if err = pf.executeRequirementsAnalysis(ctx, pipelineCtx); err != nil {
		return pf.handlePipelineError(pipelineCtx, "需求分析失败", err)
	}

	// 阶段2: 架构设计
	pipelineCtx.CurrentStage = "architecture_design"
	if err = pf.executeArchitectureDesign(ctx, pipelineCtx); err != nil {
		return pf.handlePipelineError(pipelineCtx, "架构设计失败", err)
	}

	// 阶段3: 任务分解
	pipelineCtx.CurrentStage = "task_decomposition"
	if err = pf.executeTaskDecomposition(ctx, pipelineCtx); err != nil {
		return pf.handlePipelineError(pipelineCtx, "任务分解失败", err)
	}

	// 阶段4: 并行开发
	pipelineCtx.CurrentStage = "parallel_development"
	if err = pf.executeParallelDevelopment(ctx, pipelineCtx); err != nil {
		return pf.handlePipelineError(pipelineCtx, "并行开发失败", err)
	}

	// 阶段5: 测试
	pipelineCtx.CurrentStage = "testing"
	if err = pf.executeTesting(ctx, pipelineCtx); err != nil {
		return pf.handlePipelineError(pipelineCtx, "测试失败", err)
	}

	// 阶段6: 部署
	pipelineCtx.CurrentStage = "deployment"
	if err = pf.executeDeployment(ctx, pipelineCtx); err != nil {
		return pf.handlePipelineError(pipelineCtx, "部署失败", err)
	}

	// 阶段7: 代码审查
	pipelineCtx.CurrentStage = "code_review"
	if err = pf.executeCodeReview(ctx, pipelineCtx); err != nil {
		return pf.handlePipelineError(pipelineCtx, "代码审查失败", err)
	}

	// 阶段8: 项目验证
	pipelineCtx.CurrentStage = "project_validation"
	if err = pf.executeProjectValidation(ctx, pipelineCtx); err != nil {
		return pf.handlePipelineError(pipelineCtx, "项目验证失败", err)
	}

	// 流水线执行成功
	pipelineCtx.CurrentStage = "completed"
	pipelineCtx.AddLog("success", "流水线执行完成", map[string]interface{}{
		"total_duration": time.Since(pipelineCtx.StartTime).String(),
	})

	result = &PipelineResult{
		Success:       true,
		ProjectPath:   pipelineCtx.GetArtifact("project_path", "").(string),
		ExecutionTime: time.Since(pipelineCtx.StartTime),
		Artifacts:     pipelineCtx.Artifacts,
		Logs:          pipelineCtx.ExecutionLog,
		Summary: map[string]interface{}{
			"total_stages":      8,
			"completed_stages":  8,
			"failed_stages":     0,
			"total_artifacts":   len(pipelineCtx.Artifacts),
			"project_generated": true,
		},
	}

	utils.Info("软件项目开发流水线执行完成，耗时: %v", result.ExecutionTime)
	return result, nil
}

// validateParameters 验证参数
func (pf *PipelineFlow) validateParameters(parameters map[string]interface{}) error {
	requiredParams := []string{"project_name", "project_description", "requirements"}

	for _, param := range requiredParams {
		if value, ok := parameters[param]; !ok || value == "" {
			return fmt.Errorf("缺少必要参数: %s", param)
		}
	}

	return nil
}

// executeRequirementsAnalysis 执行需求分析
func (pf *PipelineFlow) executeRequirementsAnalysis(ctx context.Context, pipelineCtx *PipelineContext) error {
	pipelineCtx.AddLog("info", "开始需求分析阶段", nil)

	// 获取产品经理Agent
	agent, err := pf.agentManager.GetAgent("product_manager_001")
	if err != nil {
		return fmt.Errorf("获取产品经理Agent失败: %w", err)
	}

	// 准备参数
	params := map[string]interface{}{
		"project_name":        pipelineCtx.Parameters["project_name"],
		"project_description": pipelineCtx.Parameters["project_description"],
		"requirements":        pipelineCtx.Parameters["requirements"],
	}

	// 执行需求分析
	task := "分析项目需求并生成PRD文档"
	result, err := agent.Execute(ctx, task, params)
	if err != nil {
		return fmt.Errorf("需求分析执行失败: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("需求分析失败: %v", result.Output)
	}

	// 保存产物
	pipelineCtx.SetArtifact("prd_document", result.Output)
	pipelineCtx.SetArtifact("requirements_analysis", result.Artifacts)

	pipelineCtx.AddLog("success", "需求分析完成", map[string]interface{}{
		"agent":    agent.GetConfig().Name,
		"duration": result.Duration.String(),
	})

	return nil
}

// executeArchitectureDesign 执行架构设计
func (pf *PipelineFlow) executeArchitectureDesign(ctx context.Context, pipelineCtx *PipelineContext) error {
	pipelineCtx.AddLog("info", "开始架构设计阶段", nil)

	// 获取技术架构师Agent
	agent, err := pf.agentManager.GetAgent("architect_001")
	if err != nil {
		return fmt.Errorf("获取技术架构师Agent失败: %w", err)
	}

	// 准备参数
	params := map[string]interface{}{
		"project_name":           pipelineCtx.Parameters["project_name"],
		"project_description":    pipelineCtx.Parameters["project_description"],
		"technical_requirements": pipelineCtx.GetArtifact("prd_document", ""),
	}

	// 执行架构设计
	task := "设计项目技术架构"
	result, err := agent.Execute(ctx, task, params)
	if err != nil {
		return fmt.Errorf("架构设计执行失败: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("架构设计失败: %v", result.Output)
	}

	// 保存产物
	pipelineCtx.SetArtifact("architecture_design", result.Output)
	pipelineCtx.SetArtifact("technical_specification", result.Artifacts)

	pipelineCtx.AddLog("success", "架构设计完成", map[string]interface{}{
		"agent":    agent.GetConfig().Name,
		"duration": result.Duration.String(),
	})

	return nil
}

// executeTaskDecomposition 执行任务分解
func (pf *PipelineFlow) executeTaskDecomposition(ctx context.Context, pipelineCtx *PipelineContext) error {
	pipelineCtx.AddLog("info", "开始任务分解阶段", nil)

	// 获取编排器Agent
	agent, err := pf.agentManager.GetAgent("orchestrator_001")
	if err != nil {
		return fmt.Errorf("获取编排器Agent失败: %w", err)
	}

	// 准备参数
	params := map[string]interface{}{
		"project_name":        pipelineCtx.Parameters["project_name"],
		"project_description": pipelineCtx.Parameters["project_description"],
		"prd_document":        pipelineCtx.GetArtifact("prd_document", ""),
		"architecture_design": pipelineCtx.GetArtifact("architecture_design", ""),
	}

	// 执行任务分解
	task := "分解项目任务并生成执行计划"
	result, err := agent.Execute(ctx, task, params)
	if err != nil {
		return fmt.Errorf("任务分解执行失败: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("任务分解失败: %v", result.Output)
	}

	// 保存产物
	pipelineCtx.SetArtifact("task_plan", result.Output)
	pipelineCtx.SetArtifact("execution_schedule", result.Artifacts)

	pipelineCtx.AddLog("success", "任务分解完成", map[string]interface{}{
		"agent":    agent.GetConfig().Name,
		"duration": result.Duration.String(),
	})

	return nil
}

// handlePipelineError 处理流水线错误
func (pf *PipelineFlow) handlePipelineError(pipelineCtx *PipelineContext, message string, err error) (*PipelineResult, error) {
	pipelineCtx.AddLog("error", message, map[string]interface{}{
		"error": err.Error(),
		"stage": pipelineCtx.CurrentStage,
	})

	result := &PipelineResult{
		Success:       false,
		Error:         fmt.Errorf("%s: %w", message, err),
		ExecutionTime: time.Since(pipelineCtx.StartTime),
		Artifacts:     pipelineCtx.Artifacts,
		Logs:          pipelineCtx.ExecutionLog,
		FailedStage:   pipelineCtx.CurrentStage,
		Summary: map[string]interface{}{
			"total_stages":      8,
			"completed_stages":  getCompletedStagesCount(pipelineCtx.CurrentStage),
			"failed_stages":     1,
			"total_artifacts":   len(pipelineCtx.Artifacts),
			"project_generated": false,
		},
	}

	utils.Error("流水线执行失败: %v", err)
	return result, result.Error
}

// getCompletedStagesCount 获取已完成的阶段数量
func getCompletedStagesCount(currentStage string) int {
	stagesOrder := []string{
		"initializing",
		"requirements_analysis",
		"architecture_design",
		"task_decomposition",
		"parallel_development",
		"testing",
		"deployment",
		"code_review",
		"project_validation",
		"completed",
	}

	for i, stage := range stagesOrder {
		if stage == currentStage {
			return i
		}
	}

	return 0
}
