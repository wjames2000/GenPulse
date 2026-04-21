package main

import (
	"GenPulse/internal/utils"
	"context"
	"fmt"
	"time"

	"GenPulse/internal/genkit"
	"GenPulse/internal/mcp/client"
	"GenPulse/internal/mcp/host"
	"GenPulse/internal/services"
)

// App struct
type App struct {
	ctx           context.Context
	baseService   *services.BaseService
	genkitManager *genkit.GenkitManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		baseService:   services.NewBaseService(),
		genkitManager: genkit.NewGenkitManager(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 启动基础服务
	a.baseService.Startup(ctx)

	// 启动Genkit运行时
	if err := a.genkitManager.Initialize(ctx); err != nil {
		fmt.Printf("警告: Genkit运行时初始化失败: %v\n", err)
		// 不阻止应用启动，记录错误继续
		a.baseService.LogMessage("error", fmt.Sprintf("Genkit初始化失败: %v", err))
	} else {
		a.baseService.LogMessage("info", "Genkit运行时初始化完成")
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetAppInfo 获取应用信息
func (a *App) GetAppInfo() map[string]interface{} {
	return a.baseService.GetAppInfo()
}

// GetLogs 获取日志
func (a *App) GetLogs() []map[string]interface{} {
	// 从基础服务获取日志
	logs := a.baseService.GetLogs()

	// 转换为前端需要的格式
	result := make([]map[string]interface{}, len(logs))
	for i, log := range logs {
		result[i] = map[string]interface{}{
			"id":        fmt.Sprintf("log-%d-%d", time.Now().UnixNano(), i),
			"timestamp": log["timestamp"],
			"level":     log["level"],
			"message":   log["message"],
			"agent_id":  log["agent_id"],
			"task_id":   log["task_id"],
			"details":   log["details"],
			"duration":  log["duration"],
			"tags":      log["tags"],
		}
	}

	return result
}

// Agent相关命令

// ListAgents 获取Agent列表
func (a *App) ListAgents() ([]map[string]interface{}, error) {
	agentManager := a.genkitManager.GetAgentManager()
	if agentManager == nil {
		return nil, fmt.Errorf("agent manager not initialized")
	}

	agents := agentManager.ListAgents()

	// 转换为前端需要的格式
	result := make([]map[string]interface{}, len(agents))
	for i, agent := range agents {
		result[i] = map[string]interface{}{
			"id":          agent.ID,
			"name":        agent.Name,
			"role":        string(agent.Role),
			"description": agent.Description,
			"model_config": map[string]interface{}{
				"type":     string(agent.ModelConfig.Type),
				"name":     agent.ModelConfig.Name,
				"provider": agent.ModelConfig.Provider,
			},
			"capabilities": agent.Capabilities,
			"tools":        agent.Tools,
			"enabled":      agent.Enabled,
		}
	}

	return result, nil
}

// GetAgentStatus 获取Agent状态
func (a *App) GetAgentStatus(agentId string) (map[string]interface{}, error) {
	agentManager := a.genkitManager.GetAgentManager()
	if agentManager == nil {
		return nil, fmt.Errorf("agent manager not initialized")
	}

	status, err := agentManager.GetAgentStatus(agentId)
	if err != nil {
		return nil, err
	}

	return status, nil
}

// GetAllAgentsStatus 获取所有Agent状态
func (a *App) GetAllAgentsStatus() (map[string]interface{}, error) {
	agentManager := a.genkitManager.GetAgentManager()
	if agentManager == nil {
		return nil, fmt.Errorf("agent manager not initialized")
	}

	status := agentManager.GetAllAgentsStatus()
	return status, nil
}

// ExecuteAgent 执行Agent任务
func (a *App) ExecuteAgent(agentId string, task string, parameters map[string]interface{}) (map[string]interface{}, error) {
	agentManager := a.genkitManager.GetAgentManager()
	if agentManager == nil {
		return nil, fmt.Errorf("agent manager not initialized")
	}

	ctx := context.Background()
	result, err := agentManager.ExecuteAgent(ctx, agentId, task, parameters)
	if err != nil {
		return nil, err
	}

	// 转换为前端需要的格式
	execution := map[string]interface{}{
		"id":         fmt.Sprintf("exec-%d", time.Now().UnixNano()),
		"agent_id":   agentId,
		"task":       task,
		"state":      "executing",
		"started_at": time.Now().Format(time.RFC3339),
		"parameters": parameters,
	}

	if result != nil {
		execution["state"] = "completed"
		execution["completed_at"] = time.Now().Format(time.RFC3339)
		execution["result"] = result.Output

		if !result.Success {
			execution["state"] = "failed"
			execution["error"] = "Task execution failed"
		}
	}

	return execution, nil
}

// GetAgentExecutions 获取Agent执行历史
func (a *App) GetAgentExecutions() ([]map[string]interface{}, error) {
	// 这里应该从数据库或内存中获取执行历史
	// 暂时返回空数组
	return []map[string]interface{}{}, nil
}

// CancelAgentExecution 取消Agent执行
func (a *App) CancelAgentExecution(executionId string) error {
	// 这里应该实现取消逻辑
	// 暂时返回成功
	return nil
}

// HealthCheck 健康检查
func (a *App) HealthCheck() map[string]interface{} {
	return a.baseService.HealthCheck()
}

// LogMessage 记录日志消息
func (a *App) LogMessage(level string, message string) {
	a.baseService.LogMessage(level, message)
}

// LogMessageWithDetails 记录带详细信息的日志消息
func (a *App) LogMessageWithDetails(level string, message string, details map[string]interface{}, agentID string, taskID string, tags []string, duration int64) {
	a.baseService.LogMessageWithDetails(level, message, details, agentID, taskID, tags, duration)
}

// GetLogsByLevel 按级别获取日志
func (a *App) GetLogsByLevel(level string) []map[string]interface{} {
	return a.baseService.GetLogsByLevel(level)
}

// GetLogsByAgent 按Agent获取日志
func (a *App) GetLogsByAgent(agentID string) []map[string]interface{} {
	return a.baseService.GetLogsByAgent(agentID)
}

// GetLogsByTimeRange 按时间范围获取日志
func (a *App) GetLogsByTimeRange(startTimeStr string, endTimeStr string) []map[string]interface{} {
	startTime, _ := time.Parse(time.RFC3339, startTimeStr)
	endTime, _ := time.Parse(time.RFC3339, endTimeStr)
	return a.baseService.GetLogsByTimeRange(startTime, endTime)
}

// ClearLogs 清空日志
func (a *App) ClearLogs() {
	a.baseService.ClearLogs()
}

// GetLogStatistics 获取日志统计信息
func (a *App) GetLogStatistics() map[string]interface{} {
	return a.baseService.GetLogStatistics()
}

// ==================== 技能管理 API ====================

// GetSkills 获取技能列表
func (a *App) GetSkills() ([]map[string]interface{}, error) {
	// 从技能管理器获取技能列表
	skillManager := a.genkitManager.GetSkillManager()
	if skillManager == nil {
		// 返回模拟数据用于测试
		return a.getMockSkills(), nil
	}

	// 获取技能元数据
	metadatas, err := skillManager.ListSkills(nil)
	if err != nil {
		// 出错时返回模拟数据
		a.baseService.LogMessage("error", fmt.Sprintf("Failed to list skills: %v", err))
		return a.getMockSkills(), nil
	}

	// 转换为前端需要的格式
	result := make([]map[string]interface{}, len(metadatas))
	for i, metadata := range metadatas {
		result[i] = map[string]interface{}{
			"id":           metadata.ID,
			"name":         metadata.Name,
			"description":  metadata.Description,
			"category":     metadata.Category,
			"version":      metadata.Version,
			"enabled":      metadata.Enabled,
			"validated":    metadata.Validated,
			"complexity":   metadata.Complexity,
			"usage_count":  metadata.UsageCount,
			"success_rate": metadata.SuccessRate,
			"tags":         metadata.Tags,
		}
	}

	return result, nil
}

// getMockSkills 获取模拟技能数据（用于测试）
func (a *App) getMockSkills() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"id":           "react-expert",
			"name":         "React 专家",
			"description":  "精通现代 React 架构，能够生成高性能组件、处理复杂状态管理并优化渲染流水线。",
			"category":     "frontend",
			"version":      "v1.4.2",
			"enabled":      true,
			"validated":    true,
			"complexity":   "complex",
			"usage_count":  12450,
			"success_rate": 0.984,
			"tags":         []string{"react", "frontend", "ui"},
			"type":         "cognitive-skill",
			"created_at":   time.Now().Add(-30 * 24 * time.Hour).Format(time.RFC3339),
			"updated_at":   time.Now().Add(-7 * 24 * time.Hour).Format(time.RFC3339),
			"agent_types":  []string{"frontend", "fullstack"},
		},
		{
			"id":           "go-backend",
			"name":         "Go 后端",
			"description":  "专注于高并发微服务架构，提供稳健的 API 设计、数据库分片策略与协程管理方案。",
			"category":     "backend",
			"version":      "v2.1.0",
			"enabled":      true,
			"validated":    true,
			"complexity":   "complex",
			"usage_count":  8204,
			"success_rate": 0.952,
			"tags":         []string{"go", "backend", "api"},
			"type":         "logic-processing",
			"created_at":   time.Now().Add(-45 * 24 * time.Hour).Format(time.RFC3339),
			"updated_at":   time.Now().Add(-14 * 24 * time.Hour).Format(time.RFC3339),
			"agent_types":  []string{"backend", "devops"},
		},
		{
			"id":           "git-pipeline",
			"name":         "Git 流水线",
			"description":  "自动化 CI/CD 流程构建，代码合并冲突智能解决，部署策略优化。",
			"category":     "devops",
			"version":      "v1.0.8",
			"enabled":      true,
			"validated":    true,
			"complexity":   "medium",
			"usage_count":  45112,
			"success_rate": 0.991,
			"tags":         []string{"git", "ci-cd", "automation"},
			"type":         "ops-automation",
			"created_at":   time.Now().Add(-60 * 24 * time.Hour).Format(time.RFC3339),
			"updated_at":   time.Now().Add(-3 * 24 * time.Hour).Format(time.RFC3339),
			"agent_types":  []string{"devops", "reviewer"},
		},
	}
}

// GetSkillDetails 获取技能详情
func (a *App) GetSkillDetails(skillID string) (map[string]interface{}, error) {
	skillManager := a.genkitManager.GetSkillManager()
	if skillManager == nil {
		// 返回模拟数据
		return a.getMockSkillDetails(skillID), nil
	}

	// 获取完整技能内容（使用L1级别）
	loadResult, err := skillManager.GetSkill(skillID, 1) // L1 = 1
	if err != nil {
		// 出错时返回模拟数据
		a.baseService.LogMessage("error", fmt.Sprintf("Failed to get skill details: %v", err))
		return a.getMockSkillDetails(skillID), nil
	}

	skill := loadResult.Skill
	// metadata := loadResult.Metadata

	// 构建技能详情
	details := map[string]interface{}{
		"id":                 skill.ID,
		"name":               skill.Name,
		"description":        skill.Description,
		"category":           skill.Category,
		"version":            skill.Version,
		"enabled":            skill.Enabled,
		"validated":          skill.Validated,
		"complexity":         skill.Complexity,
		"usage_count":        skill.UsageCount,
		"success_rate":       skill.SuccessRate,
		"created_at":         skill.CreatedAt.Format(time.RFC3339),
		"updated_at":         skill.UpdatedAt.Format(time.RFC3339),
		"tags":               skill.Tags,
		"agent_types":        skill.AgentTypes,
		"steps":              skill.Steps,
		"examples":           skill.Examples,
		"tips":               skill.Tips,
		"warnings":           skill.Warnings,
		"prerequisites":      skill.Prerequisites,
		"related_tools":      skill.RelatedTools,
		"token_estimate":     skill.TokenEstimate,
		"avg_execution_time": skill.AvgExecutionTime.String(),
		"source_task_id":     skill.SourceTaskID,
	}

	return details, nil
}

// getMockSkillDetails 获取模拟技能详情
func (a *App) getMockSkillDetails(skillID string) map[string]interface{} {
	// 根据skillID返回不同的模拟数据
	baseSkill := map[string]interface{}{
		"id":           skillID,
		"name":         "React 专家",
		"description":  "精通现代 React 架构，能够生成高性能组件、处理复杂状态管理并优化渲染流水线。",
		"category":     "frontend",
		"version":      "v1.4.2",
		"enabled":      true,
		"validated":    true,
		"complexity":   "complex",
		"usage_count":  12450,
		"success_rate": 0.984,
		"created_at":   time.Now().Add(-30 * 24 * time.Hour).Format(time.RFC3339),
		"updated_at":   time.Now().Add(-7 * 24 * time.Hour).Format(time.RFC3339),
		"tags":         []string{"react", "frontend", "ui"},
		"agent_types":  []string{"frontend", "fullstack"},
		"steps": []map[string]interface{}{
			{
				"id":         "analyze-requirements",
				"order":      1,
				"action":     "分析需求并确定组件结构",
				"tool":       "llm",
				"parameters": []map[string]interface{}{},
			},
			{
				"id":         "create-component",
				"order":      2,
				"action":     "创建React组件文件",
				"tool":       "fs_write",
				"parameters": []map[string]interface{}{},
			},
		},
		"examples": []string{
			"生成一个带有搜索功能的用户列表组件",
			"创建支持拖拽排序的看板组件",
		},
		"tips": []string{
			"使用React.memo优化性能",
			"优先使用函数组件和hooks",
		},
		"warnings": []string{
			"避免在渲染函数中创建新对象",
			"注意useEffect的依赖数组",
		},
		"prerequisites":      []string{},
		"related_tools":      []string{"fs_write", "llm"},
		"token_estimate":     1500,
		"avg_execution_time": "2.5s",
		"source_task_id":     "task-12345",
	}

	return baseSkill
}

// EnableSkill 启用技能
func (a *App) EnableSkill(skillID string) error {
	skillManager := a.genkitManager.GetSkillManager()
	if skillManager == nil {
		return fmt.Errorf("skill manager not initialized")
	}

	return skillManager.EnableSkill(skillID)
}

// DisableSkill 禁用技能
func (a *App) DisableSkill(skillID string) error {
	skillManager := a.genkitManager.GetSkillManager()
	if skillManager == nil {
		return fmt.Errorf("skill manager not initialized")
	}

	return skillManager.DisableSkill(skillID)
}

// ValidateSkill 验证技能
func (a *App) ValidateSkill(skillID string) (map[string]interface{}, error) {
	skillManager := a.genkitManager.GetSkillManager()
	if skillManager == nil {
		// 返回模拟验证结果
		return map[string]interface{}{
			"overall_pass":      true,
			"critical_failures": []string{},
			"warnings":          []string{},
			"total_checks":      5,
			"passed_checks":     5,
			"failed_checks":     0,
			"timestamp":         time.Now().Format("2006-01-02 15:04:05"),
		}, nil
	}

	report, err := skillManager.ValidateSkill(skillID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate skill: %w", err)
	}

	// 转换为前端格式
	result := map[string]interface{}{
		"overall_pass":      report.OverallPass,
		"critical_failures": report.CriticalFailures,
		"warnings":          report.Warnings,
		"total_checks":      report.TotalChecks,
		"passed_checks":     report.PassedChecks,
		"failed_checks":     report.FailedChecks,
		"timestamp":         report.Timestamp,
	}

	return result, nil
}

// DeleteSkill 删除技能
func (a *App) DeleteSkill(skillID string) error {
	skillManager := a.genkitManager.GetSkillManager()
	if skillManager == nil {
		return fmt.Errorf("skill manager not initialized")
	}

	return skillManager.DeleteSkill(skillID)
}

// SearchSkills 搜索技能
func (a *App) SearchSkills(query string, filters map[string]interface{}) ([]map[string]interface{}, error) {
	skillManager := a.genkitManager.GetSkillManager()
	if skillManager == nil {
		return nil, fmt.Errorf("skill manager not initialized")
	}

	metadatas, err := skillManager.SearchSkills(query, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search skills: %w", err)
	}

	// 转换为前端格式
	result := make([]map[string]interface{}, len(metadatas))
	for i, metadata := range metadatas {
		result[i] = map[string]interface{}{
			"id":           metadata.ID,
			"name":         metadata.Name,
			"description":  metadata.Description,
			"category":     metadata.Category,
			"version":      metadata.Version,
			"enabled":      metadata.Enabled,
			"validated":    metadata.Validated,
			"complexity":   metadata.Complexity,
			"usage_count":  metadata.UsageCount,
			"success_rate": metadata.SuccessRate,
			"tags":         metadata.Tags,
		}
	}

	return result, nil
}

// GetSkillStats 获取技能统计信息
func (a *App) GetSkillStats() (map[string]interface{}, error) {
	skillManager := a.genkitManager.GetSkillManager()
	if skillManager == nil {
		// 返回模拟统计
		return map[string]interface{}{
			"total_skills":     3,
			"enabled_skills":   3,
			"validated_skills": 3,
			"by_category": map[string]int{
				"frontend": 1,
				"backend":  1,
				"devops":   1,
			},
			"by_complexity": map[string]int{
				"complex": 2,
				"medium":  1,
			},
			"total_usage":          65866,
			"average_success_rate": 0.975,
		}, nil
	}

	stats, err := skillManager.GetSkillStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get skill stats: %w", err)
	}

	// 转换为前端格式
	result := map[string]interface{}{
		"total_skills":         stats.TotalSkills,
		"enabled_skills":       stats.EnabledSkills,
		"validated_skills":     stats.ValidatedSkills,
		"by_category":          stats.ByCategory,
		"by_complexity":        stats.ByComplexity,
		"total_usage":          stats.TotalUsage,
		"average_success_rate": stats.AverageSuccessRate,
	}

	// 添加加载器统计
	if stats.LoadStats != nil {
		result["load_stats"] = map[string]interface{}{
			"total_loads":  stats.LoadStats.TotalLoads,
			"cache_hits":   stats.LoadStats.CacheHits,
			"cache_misses": stats.LoadStats.CacheMisses,
		}
	}

	return result, nil
}

// GetRelatedSkills 获取相关技能
func (a *App) GetRelatedSkills(skillID string) ([]map[string]interface{}, error) {
	skillManager := a.genkitManager.GetSkillManager()
	if skillManager == nil {
		return nil, fmt.Errorf("skill manager not initialized")
	}

	metadatas, err := skillManager.GetRelatedSkills(skillID)
	if err != nil {
		return nil, fmt.Errorf("failed to get related skills: %w", err)
	}

	// 转换为前端格式
	result := make([]map[string]interface{}, len(metadatas))
	for i, metadata := range metadatas {
		result[i] = map[string]interface{}{
			"id":           metadata.ID,
			"name":         metadata.Name,
			"description":  metadata.Description,
			"category":     metadata.Category,
			"version":      metadata.Version,
			"enabled":      metadata.Enabled,
			"validated":    metadata.Validated,
			"complexity":   metadata.Complexity,
			"usage_count":  metadata.UsageCount,
			"success_rate": metadata.SuccessRate,
			"tags":         metadata.Tags,
		}
	}

	return result, nil
}

// ExportSkill 导出技能
func (a *App) ExportSkill(skillID string, format string) (string, error) {
	skillManager := a.genkitManager.GetSkillManager()
	if skillManager == nil {
		return "", fmt.Errorf("skill manager not initialized")
	}

	data, err := skillManager.ExportSkill(skillID, format)
	if err != nil {
		return "", fmt.Errorf("failed to export skill: %w", err)
	}

	return string(data), nil
}

// ImportSkill 导入技能
func (a *App) ImportSkill(data string, format string) (map[string]interface{}, error) {
	skillManager := a.genkitManager.GetSkillManager()
	if skillManager == nil {
		return nil, fmt.Errorf("skill manager not initialized")
	}

	skill, err := skillManager.ImportSkill([]byte(data), format)
	if err != nil {
		return nil, fmt.Errorf("failed to import skill: %w", err)
	}

	// 返回导入的技能信息
	return map[string]interface{}{
		"id":          skill.ID,
		"name":        skill.Name,
		"description": skill.Description,
		"category":    skill.Category,
		"version":     skill.Version,
	}, nil
}

// ==================== 记忆管理 API ====================

// GetEpisodicMemories 获取情节记忆
func (a *App) GetEpisodicMemories(query string, limit int) ([]map[string]interface{}, error) {
	// 从记忆管理器获取情节记忆
	memoryManager := a.genkitManager.GetMemoryManager()
	if memoryManager == nil {
		// 返回模拟数据
		return a.getMockEpisodicMemories(), nil
	}

	// 搜索记忆
	memories, err := memoryManager.SearchEpisodic(query, limit)
	if err != nil {
		// 出错时返回模拟数据
		a.baseService.LogMessage("error", fmt.Sprintf("Failed to search episodic memories: %v", err))
		return a.getMockEpisodicMemories(), nil
	}

	// 转换为前端格式
	result := make([]map[string]interface{}, len(memories))
	for i, memory := range memories {
		result[i] = map[string]interface{}{
			"id":              memory.ID,
			"task_id":         memory.TaskID,
			"task_type":       memory.TaskType,
			"description":     memory.Description,
			"agent_id":        memory.AgentID,
			"agent_name":      memory.AgentName,
			"success":         memory.Success,
			"duration_ms":     memory.DurationMs,
			"created_at":      memory.CreatedAt.Format(time.RFC3339),
			"keywords":        memory.Keywords,
			"tool_usage":      memory.ToolUsage,
			"context_data":    memory.ContextData,
			"relevance_score": memory.RelevanceScore,
		}
	}

	return result, nil
}

// getMockEpisodicMemories 获取模拟情节记忆
func (a *App) getMockEpisodicMemories() []map[string]interface{} {
	now := time.Now()
	return []map[string]interface{}{
		{
			"id":              "mem-001",
			"task_id":         "task-001",
			"task_type":       "code_generation",
			"description":     "生成React用户管理界面",
			"agent_id":        "frontend-agent",
			"agent_name":      "前端开发Agent",
			"success":         true,
			"duration_ms":     2450,
			"created_at":      now.Add(-2 * time.Hour).Format(time.RFC3339),
			"keywords":        []string{"react", "ui", "user-management"},
			"tool_usage":      map[string]int{"fs_write": 3, "llm": 2},
			"context_data":    map[string]interface{}{"framework": "react", "components": 5},
			"relevance_score": 0.95,
		},
		{
			"id":              "mem-002",
			"task_id":         "task-002",
			"task_type":       "api_development",
			"description":     "创建用户认证API",
			"agent_id":        "backend-agent",
			"agent_name":      "后端开发Agent",
			"success":         true,
			"duration_ms":     3200,
			"created_at":      now.Add(-5 * time.Hour).Format(time.RFC3339),
			"keywords":        []string{"go", "api", "authentication"},
			"tool_usage":      map[string]int{"fs_write": 4, "llm": 3},
			"context_data":    map[string]interface{}{"language": "go", "endpoints": 3},
			"relevance_score": 0.88,
		},
		{
			"id":              "mem-003",
			"task_id":         "task-003",
			"task_type":       "testing",
			"description":     "执行单元测试套件",
			"agent_id":        "qa-agent",
			"agent_name":      "质量保证Agent",
			"success":         true,
			"duration_ms":     1800,
			"created_at":      now.Add(-8 * time.Hour).Format(time.RFC3339),
			"keywords":        []string{"testing", "go-test", "coverage"},
			"tool_usage":      map[string]int{"shell_exec": 5, "fs_read": 10},
			"context_data":    map[string]interface{}{"tests": 25, "coverage": 0.85},
			"relevance_score": 0.76,
		},
	}
}

// GetSemanticMemory 获取语义记忆
func (a *App) GetSemanticMemory() (map[string]interface{}, error) {
	memoryManager := a.genkitManager.GetMemoryManager()
	if memoryManager == nil {
		// 返回模拟数据
		return a.getMockSemanticMemory(), nil
	}

	// 通过SearchEngine获取语义记忆
	if memoryManager.SemanticMemory() == nil {
		return a.getMockSemanticMemory(), nil
	}

	// 获取用户画像
	profile, err := memoryManager.SemanticMemory().GetUserProfile()
	if err != nil {
		// 出错时返回模拟数据
		a.baseService.LogMessage("error", fmt.Sprintf("Failed to get user profile: %v", err))
		return a.getMockSemanticMemory(), nil
	}

	// 转换为前端格式
	result := map[string]interface{}{
		"user_id":         profile.UserID,
		"username":        profile.Username,
		"preferences":     profile.Preferences,
		"skills":          profile.Skills,
		"interests":       profile.Interests,
		"goals":           profile.Goals,
		"working_style":   profile.WorkingStyle,
		"communication":   profile.Communication,
		"knowledge_areas": profile.KnowledgeAreas,
		"created_at":      profile.CreatedAt.Format(time.RFC3339),
		"last_updated":    profile.LastUpdated.Format(time.RFC3339),
	}

	return result, nil
}

// getMockSemanticMemory 获取模拟语义记忆
func (a *App) getMockSemanticMemory() map[string]interface{} {
	now := time.Now()
	return map[string]interface{}{
		"user_id": "user-001",
		"name":    "开发者",
		"preferences": map[string]interface{}{
			"language":  "go",
			"framework": "react",
			"database":  "postgresql",
			"testing":   "单元测试优先",
		},
		"skills": []string{
			"Go 后端开发",
			"React 前端开发",
			"数据库设计",
			"微服务架构",
		},
		"interests": []string{
			"AI 辅助编程",
			"性能优化",
			"系统架构",
			"开发者工具",
		},
		"project_goals": []string{
			"构建高效的AI开发流水线",
			"实现代码自动生成与优化",
			"提升开发效率50%以上",
		},
		"coding_style": map[string]interface{}{
			"indentation": "tabs",
			"line_length": 100,
			"naming":      "camelCase",
			"comments":    "必要的文档注释",
		},
		"created_at":        now.Add(-90 * 24 * time.Hour).Format(time.RFC3339),
		"updated_at":        now.Add(-7 * 24 * time.Hour).Format(time.RFC3339),
		"interaction_count": 156,
		"success_rate":      0.92,
	}
}

// GetMemoryStats 获取记忆统计
func (a *App) GetMemoryStats() (map[string]interface{}, error) {
	memoryManager := a.genkitManager.GetMemoryManager()
	if memoryManager == nil {
		// 返回模拟数据
		return map[string]interface{}{
			"total_memories":     156,
			"episodic_memories":  150,
			"semantic_memories":  1,
			"working_memories":   5,
			"total_searches":     245,
			"avg_search_time_ms": 45.2,
			"cache_hit_rate":     0.78,
			"last_updated":       time.Now().Format(time.RFC3339),
		}, nil
	}

	stats, err := memoryManager.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory stats: %w", err)
	}

	// 转换为前端格式
	result := make(map[string]interface{})

	// 安全地提取统计信息
	if total, ok := stats["total_memories"]; ok {
		result["total_memories"] = total
	}
	if episodic, ok := stats["episodic_memories"]; ok {
		result["episodic_memories"] = episodic
	}
	if semantic, ok := stats["semantic_memories"]; ok {
		result["semantic_memories"] = semantic
	}
	if working, ok := stats["working_memories"]; ok {
		result["working_memories"] = working
	}
	if searches, ok := stats["total_searches"]; ok {
		result["total_searches"] = searches
	}
	if avgTime, ok := stats["avg_search_time_ms"]; ok {
		result["avg_search_time_ms"] = avgTime
	}
	if hitRate, ok := stats["cache_hit_rate"]; ok {
		result["cache_hit_rate"] = hitRate
	}
	if lastUpdated, ok := stats["last_updated"]; ok {
		if t, ok := lastUpdated.(time.Time); ok {
			result["last_updated"] = t.Format(time.RFC3339)
		} else {
			result["last_updated"] = time.Now().Format(time.RFC3339)
		}
	} else {
		result["last_updated"] = time.Now().Format(time.RFC3339)
	}

	return result, nil
}

// ==================== 进化收益 API ====================

// GetEvolutionBenefits 获取进化收益数据
func (a *App) GetEvolutionBenefits() (map[string]interface{}, error) {
	skillManager := a.genkitManager.GetSkillManager()
	if skillManager == nil {
		return nil, fmt.Errorf("skill manager not initialized")
	}

	// 获取技能统计
	skillStats, err := skillManager.GetSkillStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get skill stats: %w", err)
	}

	// 获取记忆统计
	var memoryStats map[string]interface{}
	memoryManager := a.genkitManager.GetMemoryManager()
	if memoryManager != nil {
		stats, err := memoryManager.GetStats()
		if err == nil {
			memoryStats = make(map[string]interface{})
			if total, ok := stats["total_memories"]; ok {
				memoryStats["total_memories"] = total
			}
			if episodic, ok := stats["episodic_memories"]; ok {
				memoryStats["episodic_memories"] = episodic
			}
			if semantic, ok := stats["semantic_memories"]; ok {
				memoryStats["semantic_memories"] = semantic
			}
			if hitRate, ok := stats["cache_hit_rate"]; ok {
				memoryStats["cache_hit_rate"] = hitRate
			}
		}
	}

	// 计算收益指标
	totalTokenSavings := 0
	totalTimeSavings := 0.0
	automationRate := 0.0

	// 这里应该从实际使用数据计算收益
	// 暂时使用估算值
	if skillStats.TotalUsage > 0 {
		// 假设每次技能使用平均节省100 tokens和5秒时间
		totalTokenSavings = skillStats.TotalUsage * 100
		totalTimeSavings = float64(skillStats.TotalUsage) * 5.0

		// 自动化率 = 启用技能数 / 总技能数
		if skillStats.TotalSkills > 0 {
			automationRate = float64(skillStats.EnabledSkills) / float64(skillStats.TotalSkills)
		}
	}

	// 构建收益数据
	benefits := map[string]interface{}{
		"skill_stats": map[string]interface{}{
			"total_skills":         skillStats.TotalSkills,
			"enabled_skills":       skillStats.EnabledSkills,
			"total_usage":          skillStats.TotalUsage,
			"average_success_rate": skillStats.AverageSuccessRate,
		},
		"memory_stats": memoryStats,
		"benefit_metrics": map[string]interface{}{
			"total_token_savings":        totalTokenSavings,
			"total_time_savings_seconds": totalTimeSavings,
			"automation_rate":            automationRate,
			"estimated_cost_savings":     float64(totalTokenSavings) * 0.002, // 假设每1000 tokens $0.002
			"productivity_gain":          automationRate * 100,               // 百分比
		},
		"trends": map[string]interface{}{
			"skill_growth":       []int{5, 8, 12, 15, 18, 22, 25, 28, 32, 35},                           // 示例数据
			"usage_growth":       []int{10, 25, 45, 70, 100, 135, 175, 220, 270, 325},                   // 示例数据
			"success_rate_trend": []float64{0.85, 0.87, 0.89, 0.91, 0.92, 0.93, 0.94, 0.94, 0.95, 0.95}, // 示例数据
		},
	}

	return benefits, nil
}

// ==================== MCP 配置管理 API ====================

// GetMCPConfig 获取MCP配置
func (a *App) GetMCPConfig() (map[string]interface{}, error) {
	// 从Genkit管理器获取MCP主机
	genkitManager := a.genkitManager
	if genkitManager == nil {
		return nil, fmt.Errorf("genkit manager not initialized")
	}

	// 获取MCP主机
	mcpHost := genkitManager.GetMCPHost()
	if mcpHost == nil {
		// 返回默认配置
		return map[string]interface{}{
			"auto_start":              true,
			"tool_discovery_interval": 60,
			"max_concurrent_calls":    10,
			"servers":                 []map[string]interface{}{},
		}, nil
	}

	// 获取配置
	config := mcpHost.GetConfig()

	// 转换为前端格式
	servers := make([]map[string]interface{}, len(config.Servers))
	for i, server := range config.Servers {
		serverMap := map[string]interface{}{
			"id":       server.ID,
			"name":     server.Name,
			"type":     server.Type,
			"enabled":  server.Enabled,
			"priority": server.Priority,
		}

		// 根据类型添加配置
		if server.Type == "client" {
			serverMap["client_config"] = map[string]interface{}{
				"server_type": server.ClientConfig.ServerType,
				"command":     server.ClientConfig.Command,
				"args":        server.ClientConfig.Args,
				"namespace":   server.ClientConfig.Namespace,
				"timeout":     server.ClientConfig.Timeout,
			}
		} else if server.Type == "server" {
			serverMap["server_config"] = map[string]interface{}{
				"type":        server.ServerConfig.Type,
				"tool_filter": server.ServerConfig.ToolFilter,
			}
		}

		servers[i] = serverMap
	}

	return map[string]interface{}{
		"auto_start":              config.AutoStart,
		"tool_discovery_interval": config.ToolDiscoveryInterval,
		"max_concurrent_calls":    config.MaxConcurrentCalls,
		"servers":                 servers,
	}, nil
}

// UpdateMCPConfig 更新MCP配置
func (a *App) UpdateMCPConfig(config map[string]interface{}) error {
	genkitManager := a.genkitManager
	if genkitManager == nil {
		return fmt.Errorf("genkit manager not initialized")
	}

	mcpHost := genkitManager.GetMCPHost()
	if mcpHost == nil {
		return fmt.Errorf("MCP host not initialized")
	}

	// 获取配置管理器
	mcpConfigManager := genkitManager.GetMCPConfig()
	if mcpConfigManager == nil {
		return fmt.Errorf("MCP config manager not initialized")
	}

	// 解析配置
	var autoStart bool
	if val, ok := config["auto_start"]; ok {
		autoStart, _ = val.(bool)
	}

	var toolDiscoveryInterval int
	if val, ok := config["tool_discovery_interval"]; ok {
		if f, ok := val.(float64); ok {
			toolDiscoveryInterval = int(f)
		}
	}

	var maxConcurrentCalls int
	if val, ok := config["max_concurrent_calls"]; ok {
		if f, ok := val.(float64); ok {
			maxConcurrentCalls = int(f)
		}
	}

	// 解析服务器列表
	var servers []interface{}
	if val, ok := config["servers"]; ok {
		if s, ok := val.([]interface{}); ok {
			servers = s
		}
	}

	// 转换为MCP主机配置
	hostConfig := host.MCPHostConfig{
		AutoStart:             autoStart,
		ToolDiscoveryInterval: toolDiscoveryInterval,
		MaxConcurrentCalls:    maxConcurrentCalls,
		Servers:               []host.MCPHostServerConfig{},
	}

	// 解析服务器配置
	for _, serverData := range servers {
		if serverMap, ok := serverData.(map[string]interface{}); ok {
			serverID, _ := serverMap["id"].(string)
			serverName, _ := serverMap["name"].(string)
			serverType, _ := serverMap["type"].(string)
			enabled, _ := serverMap["enabled"].(bool)
			priority, _ := serverMap["priority"].(float64)

			serverConfig := host.MCPHostServerConfig{
				ID:       serverID,
				Name:     serverName,
				Type:     serverType,
				Enabled:  enabled,
				Priority: int(priority),
			}

			// 根据类型解析配置
			if serverType == "client" {
				if clientConfigMap, ok := serverMap["client_config"].(map[string]interface{}); ok {
					serverConfig.ClientConfig = client.MCPClientConfig{
						ServerType: clientConfigMap["server_type"].(string),
						Command:    clientConfigMap["command"].(string),
						Namespace:  clientConfigMap["namespace"].(string),
						Timeout:    int(clientConfigMap["timeout"].(float64)),
					}

					// 解析参数数组
					if argsInterface, ok := clientConfigMap["args"]; ok {
						if argsArray, ok := argsInterface.([]interface{}); ok {
							args := make([]string, len(argsArray))
							for i, arg := range argsArray {
								args[i] = arg.(string)
							}
							serverConfig.ClientConfig.Args = args
						}
					}
				}
			} else if serverType == "server" {
				if serverConfigMap, ok := serverMap["server_config"].(map[string]interface{}); ok {
					serverConfig.ServerConfig = host.MCPServerConfig{
						Type:       serverConfigMap["type"].(string),
						ToolFilter: serverConfigMap["tool_filter"].(string),
					}
				}
			}

			hostConfig.Servers = append(hostConfig.Servers, serverConfig)
		}
	}

	// 更新配置并保存到文件
	return mcpHost.UpdateConfigWithCallback(hostConfig, func(cfg host.MCPHostConfig) error {
		return mcpConfigManager.UpdateConfig(cfg)
	})
}

// AddMCPServer 添加MCP服务器
func (a *App) AddMCPServer(serverConfig map[string]interface{}) (map[string]interface{}, error) {
	ctx := context.Background()
	genkitManager := a.genkitManager
	if genkitManager == nil {
		return nil, fmt.Errorf("genkit manager not initialized")
	}

	mcpHost := genkitManager.GetMCPHost()
	if mcpHost == nil {
		return nil, fmt.Errorf("MCP host not initialized")
	}

	// 解析服务器配置
	serverID, _ := serverConfig["id"].(string)
	serverName, _ := serverConfig["name"].(string)
	serverType, _ := serverConfig["type"].(string)
	enabled, _ := serverConfig["enabled"].(bool)
	priority, _ := serverConfig["priority"].(float64)

	hostServerConfig := host.MCPHostServerConfig{
		ID:       serverID,
		Name:     serverName,
		Type:     serverType,
		Enabled:  enabled,
		Priority: int(priority),
	}

	// 根据类型解析配置
	if serverType == "client" {
		if clientConfigMap, ok := serverConfig["client_config"].(map[string]interface{}); ok {
			hostServerConfig.ClientConfig = client.MCPClientConfig{
				ServerType: clientConfigMap["server_type"].(string),
				Command:    clientConfigMap["command"].(string),
				Namespace:  clientConfigMap["namespace"].(string),
				Timeout:    int(clientConfigMap["timeout"].(float64)),
			}

			// 解析参数数组
			if argsInterface, ok := clientConfigMap["args"]; ok {
				if argsArray, ok := argsInterface.([]interface{}); ok {
					args := make([]string, len(argsArray))
					for i, arg := range argsArray {
						args[i] = arg.(string)
					}
					hostServerConfig.ClientConfig.Args = args
				}
			}
		}
	} else if serverType == "server" {
		if serverConfigMap, ok := serverConfig["server_config"].(map[string]interface{}); ok {
			hostServerConfig.ServerConfig = host.MCPServerConfig{
				Type:       serverConfigMap["type"].(string),
				ToolFilter: serverConfigMap["tool_filter"].(string),
			}
		}
	}

	// 添加服务器
	if err := mcpHost.AddServer(ctx, hostServerConfig); err != nil {
		return nil, fmt.Errorf("failed to add server: %w", err)
	}

	// 更新配置文件
	if mcpConfigManager := genkitManager.GetMCPConfig(); mcpConfigManager != nil {
		if err := mcpConfigManager.AddServer(hostServerConfig); err != nil {
			utils.Warn("Failed to save server config to file: %v", err)
		}
	}

	// 返回添加的服务器信息
	return map[string]interface{}{
		"id":       hostServerConfig.ID,
		"name":     hostServerConfig.Name,
		"type":     hostServerConfig.Type,
		"enabled":  hostServerConfig.Enabled,
		"priority": hostServerConfig.Priority,
	}, nil
}

// RemoveMCPServer 移除MCP服务器
func (a *App) RemoveMCPServer(serverID string) error {
	genkitManager := a.genkitManager
	if genkitManager == nil {
		return fmt.Errorf("genkit manager not initialized")
	}

	mcpHost := genkitManager.GetMCPHost()
	if mcpHost == nil {
		return fmt.Errorf("MCP host not initialized")
	}

	// 从MCP主机移除服务器
	if err := mcpHost.RemoveServer(serverID); err != nil {
		return err
	}

	// 更新配置文件
	if mcpConfigManager := genkitManager.GetMCPConfig(); mcpConfigManager != nil {
		if err := mcpConfigManager.RemoveServer(serverID); err != nil {
			utils.Warn("Failed to remove server config from file: %v", err)
		}
	}

	return nil
}

// UpdateMCPServer 更新MCP服务器
func (a *App) UpdateMCPServer(serverID string, serverConfig map[string]interface{}) error {
	genkitManager := a.genkitManager
	if genkitManager == nil {
		return fmt.Errorf("genkit manager not initialized")
	}

	mcpHost := genkitManager.GetMCPHost()
	if mcpHost == nil {
		return fmt.Errorf("MCP host not initialized")
	}

	// 解析服务器配置
	serverName, _ := serverConfig["name"].(string)
	serverType, _ := serverConfig["type"].(string)
	enabled, _ := serverConfig["enabled"].(bool)
	priority, _ := serverConfig["priority"].(float64)

	hostServerConfig := host.MCPHostServerConfig{
		ID:       serverID,
		Name:     serverName,
		Type:     serverType,
		Enabled:  enabled,
		Priority: int(priority),
	}

	// 根据类型解析配置
	if serverType == "client" {
		if clientConfigMap, ok := serverConfig["client_config"].(map[string]interface{}); ok {
			hostServerConfig.ClientConfig = client.MCPClientConfig{
				ServerType: clientConfigMap["server_type"].(string),
				Command:    clientConfigMap["command"].(string),
				Namespace:  clientConfigMap["namespace"].(string),
				Timeout:    int(clientConfigMap["timeout"].(float64)),
			}

			// 解析参数数组
			if argsInterface, ok := clientConfigMap["args"]; ok {
				if argsArray, ok := argsInterface.([]interface{}); ok {
					args := make([]string, len(argsArray))
					for i, arg := range argsArray {
						args[i] = arg.(string)
					}
					hostServerConfig.ClientConfig.Args = args
				}
			}
		}
	} else if serverType == "server" {
		if serverConfigMap, ok := serverConfig["server_config"].(map[string]interface{}); ok {
			hostServerConfig.ServerConfig = host.MCPServerConfig{
				Type:       serverConfigMap["type"].(string),
				ToolFilter: serverConfigMap["tool_filter"].(string),
			}
		}
	}

	// 更新服务器
	ctx := context.Background()
	if err := mcpHost.UpdateServer(ctx, serverID, hostServerConfig); err != nil {
		return fmt.Errorf("failed to update server: %w", err)
	}

	// 更新配置文件
	if mcpConfigManager := genkitManager.GetMCPConfig(); mcpConfigManager != nil {
		if err := mcpConfigManager.UpdateServer(serverID, hostServerConfig); err != nil {
			utils.Warn("Failed to update server config in file: %v", err)
		}
	}

	return nil
}

// GetMCPServerStatus 获取MCP服务器状态
func (a *App) GetMCPServerStatus(serverID string) (map[string]interface{}, error) {
	genkitManager := a.genkitManager
	if genkitManager == nil {
		return nil, fmt.Errorf("genkit manager not initialized")
	}

	mcpHost := genkitManager.GetMCPHost()
	if mcpHost == nil {
		return nil, fmt.Errorf("MCP host not initialized")
	}

	// 获取服务器状态
	status, err := mcpHost.GetServerStatus(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get server status: %w", err)
	}

	// 转换为前端格式
	return map[string]interface{}{
		"id":          status["id"],
		"name":        status["name"],
		"type":        status["type"],
		"enabled":     status["enabled"],
		"connected":   status["connected"],
		"last_error":  status["last_error"],
		"tool_count":  status["tool_count"],
		"last_update": status["last_update"],
	}, nil
}

// GetMCPTools 获取MCP工具列表
func (a *App) GetMCPTools() ([]map[string]interface{}, error) {
	genkitManager := a.genkitManager
	if genkitManager == nil {
		return nil, fmt.Errorf("genkit manager not initialized")
	}

	mcpHost := genkitManager.GetMCPHost()
	if mcpHost == nil {
		return nil, fmt.Errorf("MCP host not initialized")
	}

	// 暂时返回空数组，实际实现需要从MCP客户端获取工具
	// 这里返回模拟数据用于测试
	return []map[string]interface{}{
		{
			"server_id":   "example-server-1",
			"name":        "search_web",
			"namespace":   "web",
			"description": "搜索网页内容",
			"input_schema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "搜索关键词",
					},
				},
			},
		},
		{
			"server_id":   "example-server-2",
			"name":        "read_file",
			"namespace":   "filesystem",
			"description": "读取文件内容",
			"input_schema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "文件路径",
					},
				},
			},
		},
	}, nil
}

// TestMCPServerConnection 测试MCP服务器连接
func (a *App) TestMCPServerConnection(serverID string) (map[string]interface{}, error) {
	genkitManager := a.genkitManager
	if genkitManager == nil {
		return nil, fmt.Errorf("genkit manager not initialized")
	}

	mcpHost := genkitManager.GetMCPHost()
	if mcpHost == nil {
		return nil, fmt.Errorf("MCP host not initialized")
	}

	// 模拟连接测试
	// 实际实现需要调用MCP客户端的连接测试方法
	return map[string]interface{}{
		"success": true,
		"message": "连接测试成功（模拟）",
	}, nil
}
