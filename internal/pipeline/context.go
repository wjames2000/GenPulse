package pipeline

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"GenPulse/internal/utils"
)

// PipelineContext 流水线上下文
type PipelineContext struct {
	Parameters    map[string]interface{}  `json:"parameters"`
	Artifacts     map[string]interface{}  `json:"artifacts"`
	ExecutionLog  []PipelineLogEntry      `json:"execution_log"`
	StartTime     time.Time               `json:"start_time"`
	CurrentStage  string                  `json:"current_stage"`
	SharedContext map[string]interface{}  `json:"shared_context"`
	AgentContexts map[string]AgentContext `json:"agent_contexts"` // 各Agent的专用上下文
	StageResults  map[string]StageResult  `json:"stage_results"`  // 各阶段的结果
	mu            sync.RWMutex            `json:"-"`
}

// PipelineLogEntry 流水线日志条目
type PipelineLogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"` // info, success, warning, error
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Stage     string                 `json:"stage,omitempty"`
	AgentID   string                 `json:"agent_id,omitempty"`
}

// AgentContext Agent专用上下文
type AgentContext struct {
	AgentID     string                 `json:"agent_id"`
	Role        string                 `json:"role"`
	ContextData map[string]interface{} `json:"context_data"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// StageResult 阶段执行结果
type StageResult struct {
	StageName string        `json:"stage_name"`
	Success   bool          `json:"success"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Output    interface{}   `json:"output,omitempty"`
	Error     string        `json:"error,omitempty"`
	Artifacts []string      `json:"artifacts,omitempty"` // 生成的产物ID列表
	AgentID   string        `json:"agent_id,omitempty"`
}

// PipelineResult 流水线执行结果
type PipelineResult struct {
	Success       bool                   `json:"success"`
	Error         error                  `json:"error,omitempty"`
	ProjectPath   string                 `json:"project_path,omitempty"`
	ExecutionTime time.Duration          `json:"execution_time"`
	Artifacts     map[string]interface{} `json:"artifacts"`
	Logs          []PipelineLogEntry     `json:"logs"`
	FailedStage   string                 `json:"failed_stage,omitempty"`
	Summary       map[string]interface{} `json:"summary"`
}

// NewPipelineContext 创建新的流水线上下文
func NewPipelineContext(parameters map[string]interface{}) *PipelineContext {
	return &PipelineContext{
		Parameters:    parameters,
		Artifacts:     make(map[string]interface{}),
		ExecutionLog:  make([]PipelineLogEntry, 0),
		StartTime:     time.Now(),
		CurrentStage:  "initializing",
		SharedContext: make(map[string]interface{}),
		AgentContexts: make(map[string]AgentContext),
		StageResults:  make(map[string]StageResult),
	}
}

// AddLog 添加日志
func (pc *PipelineContext) AddLog(level, message string, data map[string]interface{}) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	entry := PipelineLogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Data:      data,
		Stage:     pc.CurrentStage,
	}

	pc.ExecutionLog = append(pc.ExecutionLog, entry)

	// 同时输出到系统日志
	logMessage := fmt.Sprintf("[%s] %s", pc.CurrentStage, message)
	switch level {
	case "info":
		utils.Info(logMessage)
	case "success":
		utils.Info("✓ " + logMessage)
	case "warning":
		utils.Warn(logMessage)
	case "error":
		utils.Error(logMessage)
	}
}

// SetArtifact 设置产物
func (pc *PipelineContext) SetArtifact(key string, value interface{}) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.Artifacts[key] = value
	pc.AddLog("info", fmt.Sprintf("设置产物: %s", key), map[string]interface{}{
		"artifact_key":  key,
		"artifact_type": fmt.Sprintf("%T", value),
	})
}

// GetArtifact 获取产物
func (pc *PipelineContext) GetArtifact(key string, defaultValue interface{}) interface{} {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	if value, ok := pc.Artifacts[key]; ok {
		return value
	}
	return defaultValue
}

// HasArtifact 检查是否有产物
func (pc *PipelineContext) HasArtifact(key string) bool {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	_, ok := pc.Artifacts[key]
	return ok
}

// SetSharedContext 设置共享上下文
func (pc *PipelineContext) SetSharedContext(key string, value interface{}) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.SharedContext[key] = value
}

// GetSharedContext 获取共享上下文
func (pc *PipelineContext) GetSharedContext(key string, defaultValue interface{}) interface{} {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	if value, ok := pc.SharedContext[key]; ok {
		return value
	}
	return defaultValue
}

// UpdateAgentContext 更新Agent上下文
func (pc *PipelineContext) UpdateAgentContext(agentID, role string, contextData map[string]interface{}) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	now := time.Now()
	if agentCtx, exists := pc.AgentContexts[agentID]; exists {
		// 更新现有上下文
		for k, v := range contextData {
			agentCtx.ContextData[k] = v
		}
		agentCtx.UpdatedAt = now
		pc.AgentContexts[agentID] = agentCtx
	} else {
		// 创建新上下文
		pc.AgentContexts[agentID] = AgentContext{
			AgentID:     agentID,
			Role:        role,
			ContextData: contextData,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}

	pc.AddLog("info", fmt.Sprintf("更新Agent上下文: %s (%s)", agentID, role), nil)
}

// GetAgentContext 获取Agent上下文
func (pc *PipelineContext) GetAgentContext(agentID string) (AgentContext, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	agentCtx, exists := pc.AgentContexts[agentID]
	return agentCtx, exists
}

// GetContextForAgent 获取为特定Agent准备的上下文数据
func (pc *PipelineContext) GetContextForAgent(agentID, role string) map[string]interface{} {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	contextData := make(map[string]interface{})

	// 添加共享上下文
	for k, v := range pc.SharedContext {
		contextData[k] = v
	}

	// 添加Agent专用上下文
	if agentCtx, exists := pc.AgentContexts[agentID]; exists {
		for k, v := range agentCtx.ContextData {
			contextData[k] = v
		}
	}

	// 添加相关产物
	relevantArtifacts := pc.getRelevantArtifactsForRole(role)
	for artifactKey, artifactValue := range relevantArtifacts {
		contextData[artifactKey] = artifactValue
	}

	// 添加上一阶段的结果
	if pc.CurrentStage != "initializing" {
		contextData["current_stage"] = pc.CurrentStage
		contextData["previous_stage_results"] = pc.StageResults
	}

	return contextData
}

// getRelevantArtifactsForRole 获取与角色相关的产物
func (pc *PipelineContext) getRelevantArtifactsForRole(role string) map[string]interface{} {
	relevantArtifacts := make(map[string]interface{})

	// 根据角色确定需要哪些产物
	switch strings.ToLower(role) {
	case "productmanager", "product_manager":
		// 产品经理需要原始需求
		if req, ok := pc.Parameters["requirements"]; ok {
			relevantArtifacts["raw_requirements"] = req
		}
		if desc, ok := pc.Parameters["project_description"]; ok {
			relevantArtifacts["project_description"] = desc
		}

	case "architect", "技术架构师":
		// 架构师需要PRD和需求
		if prd, ok := pc.Artifacts["prd_document"]; ok {
			relevantArtifacts["prd_document"] = prd
		}
		if req, ok := pc.Parameters["requirements"]; ok {
			relevantArtifacts["requirements"] = req
		}
		if ts, ok := pc.Parameters["tech_stack"]; ok {
			relevantArtifacts["tech_stack_preference"] = ts
		}

	case "frontenddev", "frontend_dev", "前端开发":
		// 前端开发需要PRD、架构设计和任务计划
		if prd, ok := pc.Artifacts["prd_document"]; ok {
			relevantArtifacts["prd_document"] = prd
		}
		if arch, ok := pc.Artifacts["architecture_design"]; ok {
			relevantArtifacts["architecture_design"] = arch
		}
		if plan, ok := pc.Artifacts["task_plan"]; ok {
			relevantArtifacts["task_plan"] = plan
		}
		if ts, ok := pc.Parameters["tech_stack"]; ok {
			relevantArtifacts["tech_stack"] = ts
		}

	case "backenddev", "backend_dev", "后端开发":
		// 后端开发需要PRD、架构设计和任务计划
		if prd, ok := pc.Artifacts["prd_document"]; ok {
			relevantArtifacts["prd_document"] = prd
		}
		if arch, ok := pc.Artifacts["architecture_design"]; ok {
			relevantArtifacts["architecture_design"] = arch
		}
		if plan, ok := pc.Artifacts["task_plan"]; ok {
			relevantArtifacts["task_plan"] = plan
		}
		if ts, ok := pc.Parameters["tech_stack"]; ok {
			relevantArtifacts["tech_stack"] = ts
		}

	case "qaengineer", "qa_engineer", "qa工程师":
		// QA工程师需要PRD、架构设计和代码
		if prd, ok := pc.Artifacts["prd_document"]; ok {
			relevantArtifacts["prd_document"] = prd
		}
		if arch, ok := pc.Artifacts["architecture_design"]; ok {
			relevantArtifacts["architecture_design"] = arch
		}
		if frontendCode, ok := pc.Artifacts["frontend_code"]; ok {
			relevantArtifacts["frontend_code"] = frontendCode
		}
		if backendCode, ok := pc.Artifacts["backend_code"]; ok {
			relevantArtifacts["backend_code"] = backendCode
		}

	case "devops", "devops_engineer", "devops工程师":
		// DevOps工程师需要所有技术相关产物
		if arch, ok := pc.Artifacts["architecture_design"]; ok {
			relevantArtifacts["architecture_design"] = arch
		}
		if frontendCode, ok := pc.Artifacts["frontend_code"]; ok {
			relevantArtifacts["frontend_code"] = frontendCode
		}
		if backendCode, ok := pc.Artifacts["backend_code"]; ok {
			relevantArtifacts["backend_code"] = backendCode
		}
		if ts, ok := pc.Parameters["tech_stack"]; ok {
			relevantArtifacts["tech_stack"] = ts
		}

	case "reviewer", "代码审查员":
		// 代码审查员需要所有代码
		if frontendCode, ok := pc.Artifacts["frontend_code"]; ok {
			relevantArtifacts["frontend_code"] = frontendCode
		}
		if backendCode, ok := pc.Artifacts["backend_code"]; ok {
			relevantArtifacts["backend_code"] = backendCode
		}
		if arch, ok := pc.Artifacts["architecture_design"]; ok {
			relevantArtifacts["architecture_design"] = arch
		}

	case "orchestrator", "编排器":
		// 编排器需要所有信息
		for k, v := range pc.Artifacts {
			relevantArtifacts[k] = v
		}
		for k, v := range pc.Parameters {
			relevantArtifacts[k] = v
		}
	}

	return relevantArtifacts
}

// RecordStageResult 记录阶段执行结果
func (pc *PipelineContext) RecordStageResult(stageName string, success bool, output interface{}, err error, agentID string, artifacts []string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	now := time.Now()
	var errorStr string
	if err != nil {
		errorStr = err.Error()
	}

	// 计算持续时间
	var duration time.Duration
	if stageResult, exists := pc.StageResults[stageName]; exists {
		duration = now.Sub(stageResult.StartTime)
	} else {
		// 如果没有开始时间，使用默认值
		duration = 0
	}

	pc.StageResults[stageName] = StageResult{
		StageName: stageName,
		Success:   success,
		StartTime: now.Add(-duration), // 回推开始时间
		EndTime:   now,
		Duration:  duration,
		Output:    output,
		Error:     errorStr,
		Artifacts: artifacts,
		AgentID:   agentID,
	}

	level := "success"
	if !success {
		level = "error"
	}

	pc.AddLog(level, fmt.Sprintf("阶段完成: %s", stageName), map[string]interface{}{
		"success":   success,
		"duration":  duration.String(),
		"agent_id":  agentID,
		"artifacts": len(artifacts),
	})
}

// GetStageResult 获取阶段结果
func (pc *PipelineContext) GetStageResult(stageName string) (StageResult, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	result, exists := pc.StageResults[stageName]
	return result, exists
}

// GetAllStageResults 获取所有阶段结果
func (pc *PipelineContext) GetAllStageResults() map[string]StageResult {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	// 创建副本
	results := make(map[string]StageResult)
	for k, v := range pc.StageResults {
		results[k] = v
	}
	return results
}

// ToJSON 将上下文转换为JSON
func (pc *PipelineContext) ToJSON() (string, error) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	data, err := json.MarshalIndent(pc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化上下文失败: %w", err)
	}
	return string(data), nil
}

// FromJSON 从JSON恢复上下文
func FromJSON(jsonStr string) (*PipelineContext, error) {
	var ctx PipelineContext
	if err := json.Unmarshal([]byte(jsonStr), &ctx); err != nil {
		return nil, fmt.Errorf("反序列化上下文失败: %w", err)
	}
	return &ctx, nil
}

// GetSummary 获取上下文摘要
func (pc *PipelineContext) GetSummary() map[string]interface{} {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	completedStages := 0
	failedStages := 0
	totalArtifacts := len(pc.Artifacts)
	totalAgents := len(pc.AgentContexts)

	for _, result := range pc.StageResults {
		if result.Success {
			completedStages++
		} else {
			failedStages++
		}
	}

	return map[string]interface{}{
		"current_stage":    pc.CurrentStage,
		"total_stages":     len(pc.StageResults),
		"completed_stages": completedStages,
		"failed_stages":    failedStages,
		"total_artifacts":  totalArtifacts,
		"total_agents":     totalAgents,
		"execution_time":   time.Since(pc.StartTime).String(),
		"log_entries":      len(pc.ExecutionLog),
	}
}
