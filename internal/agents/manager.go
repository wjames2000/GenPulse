package agents

import (
	"context"
	"fmt"
	"sync"
	"time"

	"GenPulse/internal/genkit/flows"
	"GenPulse/internal/genkit/models"
	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/utils"
)

// AgentManager Agent管理器
type AgentManager struct {
	agents       map[string]Agent
	modelAdapter *models.UnifiedModelAdapter
	toolRegistry *tools.ToolRegistry
	flowEngine   *flows.FlowEngine
	mutex        sync.RWMutex
	initialized  bool
}

// NewAgentManager 创建Agent管理器
func NewAgentManager(modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine) *AgentManager {
	return &AgentManager{
		agents:       make(map[string]Agent),
		modelAdapter: modelAdapter,
		toolRegistry: toolRegistry,
		flowEngine:   flowEngine,
		initialized:  false,
	}
}

// Initialize 初始化Agent管理器
func (am *AgentManager) Initialize() error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if am.initialized {
		return nil
	}

	utils.Info("初始化Agent管理器...")

	// 检查依赖
	if am.modelAdapter == nil {
		return fmt.Errorf("model adapter is required")
	}

	if am.toolRegistry == nil {
		return fmt.Errorf("tool registry is required")
	}

	// 创建默认Agent
	if err := am.createDefaultAgents(); err != nil {
		utils.Warn("创建默认Agent失败: %v", err)
	}

	// 初始化所有Agent
	var initErrors []string
	for agentID, agent := range am.agents {
		if err := agent.Initialize(); err != nil {
			initErrors = append(initErrors, fmt.Sprintf("%s: %v", agentID, err))
		}
	}

	if len(initErrors) > 0 {
		utils.Warn("部分Agent初始化失败: %v", initErrors)
	}

	am.initialized = true
	utils.Info("Agent管理器初始化完成，注册Agent数量: %d", len(am.agents))

	return nil
}

// createDefaultAgents 创建默认Agent
func (am *AgentManager) createDefaultAgents() error {
	// 创建全栈开发Agent
	fullstackConfig := AgentConfig{
		ID:          "fullstack_dev_001",
		Name:        "全栈开发工程师",
		Role:        RoleFullStackDev,
		Description: "全栈开发工程师，能够处理前后端开发任务",
		ModelConfig: models.ModelConfig{
			Type:     models.ModelTypeGemini,
			Name:     "gemini-1.5-pro",
			Provider: "google",
		},
		Capabilities: []AgentCapability{
			CapabilityCodeGeneration,
			CapabilityFileOperation,
			CapabilityGitOperation,
			CapabilityShellExecution,
			CapabilityProjectSetup,
		},
		Tools: []string{
			"fs_tool",
			"git_tool",
			"shell_tool",
			"project_tool",
		},
		MaxRetries: 3,
		Timeout:    5 * time.Minute,
		Enabled:    true,
	}

	fullstackAgent, err := NewFullstackDeveloperAgent(fullstackConfig, am.modelAdapter, am.toolRegistry, am.flowEngine)
	if err != nil {
		return fmt.Errorf("failed to create fullstack agent: %w", err)
	}

	if err := am.RegisterAgent(fullstackAgent); err != nil {
		return fmt.Errorf("failed to register fullstack agent: %w", err)
	}

	utils.Info("已创建默认全栈开发Agent: %s", fullstackConfig.Name)

	return nil
}

// RegisterAgent 注册Agent
func (am *AgentManager) RegisterAgent(agent Agent) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	agentID := agent.GetConfig().ID
	if agentID == "" {
		return NewAgentError(ErrInvalidConfig, "注册Agent失败", "Agent ID不能为空")
	}

	if _, exists := am.agents[agentID]; exists {
		return NewAgentError(ErrInvalidConfig, "注册Agent失败", fmt.Sprintf("Agent ID %s 已存在", agentID)).
			WithAgentID(agentID)
	}

	am.agents[agentID] = agent
	utils.Info("注册Agent: %s (%s)", agent.GetConfig().Name, agentID)

	return nil
}

// UnregisterAgent 注销Agent
func (am *AgentManager) UnregisterAgent(agentID string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	agent, exists := am.agents[agentID]
	if !exists {
		return fmt.Errorf("agent with ID %s not found", agentID)
	}

	// 关闭Agent
	if err := agent.Shutdown(); err != nil {
		utils.Warn("关闭Agent %s 失败: %v", agentID, err)
	}

	delete(am.agents, agentID)
	utils.Info("注销Agent: %s", agentID)

	return nil
}

// GetAgent 获取Agent
func (am *AgentManager) GetAgent(agentID string) (Agent, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	agent, exists := am.agents[agentID]
	if !exists {
		return nil, NewAgentError(ErrInvalidConfig, "获取Agent失败", fmt.Sprintf("Agent ID %s 未找到", agentID)).
			WithAgentID(agentID)
	}

	return agent, nil
}

// GetAgentByName 通过名称获取Agent
func (am *AgentManager) GetAgentByName(name string) (Agent, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	for _, agent := range am.agents {
		if agent.GetConfig().Name == name {
			return agent, nil
		}
	}

	return nil, fmt.Errorf("agent with name %s not found", name)
}

// ListAgents 列出所有Agent
func (am *AgentManager) ListAgents() []AgentConfig {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	var configs []AgentConfig
	for _, agent := range am.agents {
		configs = append(configs, agent.GetConfig())
	}

	return configs
}

// ListAgentsByRole 按角色列出Agent
func (am *AgentManager) ListAgentsByRole(role AgentRole) []AgentConfig {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	var configs []AgentConfig
	for _, agent := range am.agents {
		if agent.GetConfig().Role == role {
			configs = append(configs, agent.GetConfig())
		}
	}

	return configs
}

// ExecuteAgent 执行Agent任务
func (am *AgentManager) ExecuteAgent(ctx context.Context, agentID string, task string, parameters map[string]interface{}) (*AgentResult, error) {
	agent, err := am.GetAgent(agentID)
	if err != nil {
		return nil, err
	}

	if !agent.IsEnabled() {
		return nil, NewAgentError(ErrAgentDisabled, "执行Agent任务失败", fmt.Sprintf("Agent %s 已禁用", agentID)).
			WithAgentID(agentID).
			WithTask(task)
	}

	utils.Info("执行Agent任务: %s -> %s", agent.GetConfig().Name, task)
	return agent.Execute(ctx, task, parameters)
}

// ExecuteAgentByName 通过名称执行Agent任务
func (am *AgentManager) ExecuteAgentByName(ctx context.Context, agentName string, task string, parameters map[string]interface{}) (*AgentResult, error) {
	agent, err := am.GetAgentByName(agentName)
	if err != nil {
		return nil, err
	}

	return am.ExecuteAgent(ctx, agent.GetConfig().ID, task, parameters)
}

// GetAgentStatus 获取Agent状态
func (am *AgentManager) GetAgentStatus(agentID string) (map[string]interface{}, error) {
	agent, err := am.GetAgent(agentID)
	if err != nil {
		return nil, err
	}

	config := agent.GetConfig()
	execution := agent.GetExecution()

	status := map[string]interface{}{
		"id":           config.ID,
		"name":         config.Name,
		"role":         config.Role,
		"description":  config.Description,
		"state":        agent.GetState(),
		"enabled":      agent.IsEnabled(),
		"capabilities": config.Capabilities,
		"tools":        config.Tools,
		"statistics": map[string]interface{}{
			"execution_count": agent.GetExecutionCount(),
			"success_rate":    agent.GetSuccessRate(),
			"avg_duration":    agent.GetAverageDuration().String(),
		},
	}

	if execution != nil {
		status["current_execution"] = map[string]interface{}{
			"id":         execution.ID,
			"task":       execution.Task,
			"state":      execution.State,
			"started_at": execution.StartedAt,
			"error":      execution.Error,
		}

		if execution.CompletedAt != nil {
			status["current_execution"].(map[string]interface{})["completed_at"] = *execution.CompletedAt
		}
	}

	return status, nil
}

// GetAllAgentsStatus 获取所有Agent状态
func (am *AgentManager) GetAllAgentsStatus() map[string]interface{} {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	status := make(map[string]interface{})
	totalAgents := len(am.agents)
	enabledAgents := 0
	busyAgents := 0

	for agentID, agent := range am.agents {
		agentStatus, err := am.GetAgentStatus(agentID)
		if err == nil {
			status[agentID] = agentStatus

			if agent.IsEnabled() {
				enabledAgents++
			}

			if agent.GetState() != StateIdle {
				busyAgents++
			}
		}
	}

	return map[string]interface{}{
		"total_agents":   totalAgents,
		"enabled_agents": enabledAgents,
		"busy_agents":    busyAgents,
		"idle_agents":    totalAgents - busyAgents,
		"agents":         status,
	}
}

// EnableAgent 启用Agent
func (am *AgentManager) EnableAgent(agentID string) error {
	agent, err := am.GetAgent(agentID)
	if err != nil {
		return err
	}

	agent.SetEnabled(true)
	return nil
}

// DisableAgent 禁用Agent
func (am *AgentManager) DisableAgent(agentID string) error {
	agent, err := am.GetAgent(agentID)
	if err != nil {
		return err
	}

	agent.SetEnabled(false)
	return nil
}

// Shutdown 关闭Agent管理器
func (am *AgentManager) Shutdown() error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if !am.initialized {
		return nil
	}

	utils.Info("关闭Agent管理器...")

	var shutdownErrors []string
	for agentID, agent := range am.agents {
		if err := agent.Shutdown(); err != nil {
			shutdownErrors = append(shutdownErrors, fmt.Sprintf("%s: %v", agentID, err))
		}
	}

	if len(shutdownErrors) > 0 {
		utils.Warn("部分Agent关闭失败: %v", shutdownErrors)
	}

	am.agents = make(map[string]Agent)
	am.initialized = false

	utils.Info("Agent管理器已关闭")
	return nil
}

// IsInitialized 检查是否已初始化
func (am *AgentManager) IsInitialized() bool {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	return am.initialized
}

// GetAgentCount 获取Agent数量
func (am *AgentManager) GetAgentCount() int {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	return len(am.agents)
}

// GetEnabledAgentCount 获取启用的Agent数量
func (am *AgentManager) GetEnabledAgentCount() int {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	count := 0
	for _, agent := range am.agents {
		if agent.IsEnabled() {
			count++
		}
	}

	return count
}

// CreateAgentFromConfig 从配置创建Agent
func (am *AgentManager) CreateAgentFromConfig(config AgentConfig) (Agent, error) {
	switch config.Role {
	case RoleFullStackDev:
		return NewFullstackDeveloperAgent(config, am.modelAdapter, am.toolRegistry, am.flowEngine)
	// 其他角色可以在这里添加
	default:
		return NewBaseAgent(config, am.modelAdapter, am.toolRegistry, am.flowEngine)
	}
}

// LoadAgentsFromConfig 从配置加载多个Agent
func (am *AgentManager) LoadAgentsFromConfig(configs []AgentConfig) ([]string, error) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	var loaded []string
	var errors []string

	for _, config := range configs {
		agent, err := am.CreateAgentFromConfig(config)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", config.Name, err))
			continue
		}

		if err := am.RegisterAgent(agent); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", config.Name, err))
			continue
		}

		loaded = append(loaded, config.ID)
	}

	if len(errors) > 0 {
		return loaded, fmt.Errorf("部分Agent加载失败: %v", errors)
	}

	return loaded, nil
}

// Global agent manager instance
var globalAgentManager *AgentManager

// InitGlobalAgentManager 初始化全局Agent管理器
func InitGlobalAgentManager(modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine) error {
	if globalAgentManager == nil {
		globalAgentManager = NewAgentManager(modelAdapter, toolRegistry, flowEngine)
	}

	return globalAgentManager.Initialize()
}

// GetGlobalAgentManager 获取全局Agent管理器
func GetGlobalAgentManager() *AgentManager {
	return globalAgentManager
}

// GetGlobalAgent 获取全局Agent
func GetGlobalAgent(agentID string) (Agent, error) {
	if globalAgentManager == nil {
		return nil, fmt.Errorf("global agent manager not initialized")
	}

	return globalAgentManager.GetAgent(agentID)
}

// ExecuteGlobalAgent 执行全局Agent任务
func ExecuteGlobalAgent(ctx context.Context, agentID string, task string, parameters map[string]interface{}) (*AgentResult, error) {
	if globalAgentManager == nil {
		return nil, fmt.Errorf("global agent manager not initialized")
	}

	return globalAgentManager.ExecuteAgent(ctx, agentID, task, parameters)
}
