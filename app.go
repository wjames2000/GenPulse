package main

import (
	"context"
	"fmt"
	"time"

	"GenPulse/internal/genkit"
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
	// 暂时返回空数组
	return []map[string]interface{}{}
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
