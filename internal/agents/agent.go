package agents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"GenPulse/internal/genkit/flows"
	"GenPulse/internal/genkit/models"
	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/utils"
)

// AgentRole Agent角色
type AgentRole string

const (
	RoleFullStackDev   AgentRole = "fullstack_developer"
	RoleFrontendDev    AgentRole = "frontend_developer"
	RoleBackendDev     AgentRole = "backend_developer"
	RoleArchitect      AgentRole = "architect"
	RoleProductManager AgentRole = "product_manager"
	RoleQAEngineer     AgentRole = "qa_engineer"
	RoleDevOps         AgentRole = "devops"
	RoleReviewer       AgentRole = "reviewer"
	RoleOrchestrator   AgentRole = "orchestrator"
)

// AgentState Agent状态
type AgentState string

const (
	StateIdle      AgentState = "idle"
	StateThinking  AgentState = "thinking"
	StateExecuting AgentState = "executing"
	StateWaiting   AgentState = "waiting"
	StateCompleted AgentState = "completed"
	StateFailed    AgentState = "failed"
)

// AgentCapability Agent能力
type AgentCapability string

const (
	CapabilityCodeGeneration AgentCapability = "code_generation"
	CapabilityFileOperation  AgentCapability = "file_operation"
	CapabilityGitOperation   AgentCapability = "git_operation"
	CapabilityShellExecution AgentCapability = "shell_execution"
	CapabilityProjectSetup   AgentCapability = "project_setup"
	CapabilityTesting        AgentCapability = "testing"
	CapabilityReview         AgentCapability = "review"
	CapabilityPlanning       AgentCapability = "planning"
)

// AgentConfig Agent配置
type AgentConfig struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Role            AgentRole          `json:"role"`
	Description     string             `json:"description"`
	ModelConfig     models.ModelConfig `json:"model_config"`
	Capabilities    []AgentCapability  `json:"capabilities"`
	Tools           []string           `json:"tools"` // 工具ID列表
	PromptTemplates map[string]string  `json:"prompt_templates"`
	MaxRetries      int                `json:"max_retries"`
	Timeout         time.Duration      `json:"timeout"`
	Enabled         bool               `json:"enabled"`
}

// AgentExecution Agent执行上下文
type AgentExecution struct {
	ID          string                 `json:"id"`
	AgentID     string                 `json:"agent_id"`
	Task        string                 `json:"task"`
	Parameters  map[string]interface{} `json:"parameters"`
	Context     map[string]interface{} `json:"context"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	State       AgentState             `json:"state"`
	Result      *AgentResult           `json:"result,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// AgentResult Agent执行结果
type AgentResult struct {
	Success   bool                   `json:"success"`
	Output    interface{}            `json:"output,omitempty"`
	Artifacts []AgentArtifact        `json:"artifacts,omitempty"`
	Logs      []string               `json:"logs,omitempty"`
	Duration  time.Duration          `json:"duration"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// AgentArtifact Agent生成的产物
type AgentArtifact struct {
	Type        string      `json:"type"`
	Name        string      `json:"name"`
	Content     interface{} `json:"content,omitempty"`
	Path        string      `json:"path,omitempty"`
	Description string      `json:"description,omitempty"`
}

// AgentError Agent错误类型
type AgentError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	AgentID string `json:"agent_id,omitempty"`
	Task    string `json:"task,omitempty"`
}

// Error 实现error接口
func (e *AgentError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
}

// NewAgentError 创建Agent错误
func NewAgentError(code, message, details string) *AgentError {
	return &AgentError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// WithAgentID 设置Agent ID
func (e *AgentError) WithAgentID(agentID string) *AgentError {
	e.AgentID = agentID
	return e
}

// WithTask 设置任务
func (e *AgentError) WithTask(task string) *AgentError {
	e.Task = task
	return e
}

// AgentErrorCode Agent错误码
const (
	ErrAgentDisabled       = "AGENT_DISABLED"
	ErrAgentBusy           = "AGENT_BUSY"
	ErrAgentNotInitialized = "AGENT_NOT_INITIALIZED"
	ErrTaskValidation      = "TASK_VALIDATION"
	ErrTaskExecution       = "TASK_EXECUTION"
	ErrToolNotFound        = "TOOL_NOT_FOUND"
	ErrToolDisabled        = "TOOL_DISABLED"
	ErrModelNotAvailable   = "MODEL_NOT_AVAILABLE"
	ErrTimeout             = "TIMEOUT"
	ErrMaxRetriesExceeded  = "MAX_RETRIES_EXCEEDED"
	ErrInvalidConfig       = "INVALID_CONFIG"
)

// Agent interface Agent接口
type Agent interface {
	// 基本信息
	GetConfig() AgentConfig
	GetState() AgentState
	GetExecution() *AgentExecution

	// 执行功能
	Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error)
	Cancel() error
	ValidateTask(task string) error

	// 生命周期
	Initialize() error
	Shutdown() error

	// 状态管理
	IsEnabled() bool
	SetEnabled(enabled bool)

	// 统计信息
	GetExecutionCount() int
	GetSuccessRate() float64
	GetAverageDuration() time.Duration
}

// BaseAgent 基础Agent实现
type BaseAgent struct {
	config         AgentConfig
	state          AgentState
	execution      *AgentExecution
	modelAdapter   *models.UnifiedModelAdapter
	toolRegistry   *tools.ToolRegistry
	flowEngine     *flows.FlowEngine
	executionCount int
	successCount   int
	totalDuration  time.Duration
	enabled        bool
}

// NewBaseAgent 创建基础Agent
func NewBaseAgent(config AgentConfig, modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry, flowEngine *flows.FlowEngine) (*BaseAgent, error) {
	if config.ID == "" {
		return nil, NewAgentError(ErrInvalidConfig, "Agent配置无效", "Agent ID不能为空")
	}

	if config.Name == "" {
		return nil, NewAgentError(ErrInvalidConfig, "Agent配置无效", "Agent名称不能为空")
	}

	// 设置默认值
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	if config.Timeout == 0 {
		config.Timeout = 5 * time.Minute
	}

	if config.PromptTemplates == nil {
		config.PromptTemplates = make(map[string]string)
	}

	agent := &BaseAgent{
		config:       config,
		state:        StateIdle,
		modelAdapter: modelAdapter,
		toolRegistry: toolRegistry,
		flowEngine:   flowEngine,
		enabled:      config.Enabled,
	}

	return agent, nil
}

// GetConfig 获取Agent配置
func (a *BaseAgent) GetConfig() AgentConfig {
	return a.config
}

// GetState 获取Agent状态
func (a *BaseAgent) GetState() AgentState {
	return a.state
}

// GetExecution 获取当前执行上下文
func (a *BaseAgent) GetExecution() *AgentExecution {
	return a.execution
}

// Execute 执行任务（基础实现，需要子类重写）
func (a *BaseAgent) Execute(ctx context.Context, task string, parameters map[string]interface{}) (*AgentResult, error) {
	if !a.enabled {
		return nil, NewAgentError(ErrAgentDisabled, "Agent已禁用", "").
			WithAgentID(a.config.ID)
	}

	if a.state != StateIdle {
		return nil, NewAgentError(ErrAgentBusy, "Agent正忙", fmt.Sprintf("当前状态: %s", a.state)).
			WithAgentID(a.config.ID)
	}

	// 验证任务
	if err := a.ValidateTask(task); err != nil {
		return nil, NewAgentError(ErrTaskValidation, "任务验证失败", err.Error()).
			WithAgentID(a.config.ID).
			WithTask(task)
	}

	// 创建执行上下文
	executionID := fmt.Sprintf("%s-%d", a.config.ID, time.Now().Unix())
	a.execution = &AgentExecution{
		ID:         executionID,
		AgentID:    a.config.ID,
		Task:       task,
		Parameters: parameters,
		Context:    make(map[string]interface{}),
		StartedAt:  time.Now(),
		State:      StateThinking,
	}

	a.state = StateThinking
	a.executionCount++

	utils.Info("Agent %s 开始执行任务: %s", a.config.Name, task)

	// 基础实现：直接返回成功（子类需要重写）
	result := &AgentResult{
		Success:  true,
		Output:   "Task executed (base implementation)",
		Duration: 0,
	}

	a.execution.CompletedAt = &time.Time{}
	*a.execution.CompletedAt = time.Now()
	a.execution.State = StateCompleted
	a.execution.Result = result
	a.state = StateIdle

	if result.Success {
		a.successCount++
	}
	a.totalDuration += result.Duration

	return result, nil
}

// Cancel 取消当前执行
func (a *BaseAgent) Cancel() error {
	if a.state == StateIdle || a.state == StateCompleted || a.state == StateFailed {
		return fmt.Errorf("agent is not executing")
	}

	a.state = StateFailed
	if a.execution != nil {
		a.execution.State = StateFailed
		a.execution.Error = "execution cancelled by user"
	}

	utils.Info("Agent %s 执行已取消", a.config.Name)
	return nil
}

// ValidateTask 验证任务（基础实现）
func (a *BaseAgent) ValidateTask(task string) error {
	if task == "" {
		return NewAgentError(ErrTaskValidation, "任务验证失败", "任务不能为空")
	}

	// 检查任务长度
	if len(task) < 5 {
		return NewAgentError(ErrTaskValidation, "任务验证失败", "任务太短（至少5个字符）")
	}

	if len(task) > 1000 {
		return NewAgentError(ErrTaskValidation, "任务验证失败", "任务太长（最多1000个字符）")
	}

	return nil
}

// Initialize 初始化Agent
func (a *BaseAgent) Initialize() error {
	if a.state != StateIdle {
		return NewAgentError(ErrAgentBusy, "Agent初始化失败", "Agent不处于空闲状态").
			WithAgentID(a.config.ID)
	}

	// 检查模型适配器
	if a.modelAdapter == nil {
		return NewAgentError(ErrInvalidConfig, "Agent初始化失败", "模型适配器不能为空").
			WithAgentID(a.config.ID)
	}

	// 检查工具注册表
	if a.toolRegistry == nil {
		return NewAgentError(ErrInvalidConfig, "Agent初始化失败", "工具注册表不能为空").
			WithAgentID(a.config.ID)
	}

	// 验证工具可用性
	for _, toolID := range a.config.Tools {
		tool, err := a.toolRegistry.GetTool(toolID)
		if err != nil {
			return NewAgentError(ErrToolNotFound, "Agent初始化失败", fmt.Sprintf("工具 %s 未找到: %v", toolID, err)).
				WithAgentID(a.config.ID)
		}

		if !tool.IsEnabled() {
			utils.Warn("工具 %s 已禁用，Agent %s 可能无法正常工作", toolID, a.config.Name)
		}
	}

	// 初始化Flow引擎（如果可用）
	if a.flowEngine != nil {
		// Flow引擎可能已经有初始化逻辑
		// 暂时不调用初始化方法
		utils.Info("Flow引擎已关联到Agent %s", a.config.Name)
	}

	a.enabled = true
	utils.Info("Agent %s 初始化完成", a.config.Name)

	return nil
}

// Shutdown 关闭Agent
func (a *BaseAgent) Shutdown() error {
	// 如果正在执行，尝试取消
	if a.state == StateThinking || a.state == StateExecuting || a.state == StateWaiting {
		a.Cancel()
	}

	a.enabled = false
	a.state = StateIdle
	a.execution = nil

	utils.Info("Agent %s 已关闭", a.config.Name)
	return nil
}

// IsEnabled 检查Agent是否启用
func (a *BaseAgent) IsEnabled() bool {
	return a.enabled
}

// SetEnabled 设置Agent启用状态
func (a *BaseAgent) SetEnabled(enabled bool) {
	a.enabled = enabled
	if enabled {
		utils.Info("Agent %s 已启用", a.config.Name)
	} else {
		utils.Info("Agent %s 已禁用", a.config.Name)
	}
}

// GetExecutionCount 获取执行次数
func (a *BaseAgent) GetExecutionCount() int {
	return a.executionCount
}

// GetSuccessRate 获取成功率
func (a *BaseAgent) GetSuccessRate() float64 {
	if a.executionCount == 0 {
		return 0.0
	}
	return float64(a.successCount) / float64(a.executionCount) * 100
}

// GetAverageDuration 获取平均执行时间
func (a *BaseAgent) GetAverageDuration() time.Duration {
	if a.executionCount == 0 {
		return 0
	}
	return a.totalDuration / time.Duration(a.executionCount)
}

// GetPromptTemplate 获取提示词模板
func (a *BaseAgent) GetPromptTemplate(templateName string) (string, error) {
	template, exists := a.config.PromptTemplates[templateName]
	if !exists {
		return "", fmt.Errorf("prompt template '%s' not found", templateName)
	}
	return template, nil
}

// SetPromptTemplate 设置提示词模板
func (a *BaseAgent) SetPromptTemplate(templateName, template string) {
	if a.config.PromptTemplates == nil {
		a.config.PromptTemplates = make(map[string]string)
	}
	a.config.PromptTemplates[templateName] = template
}

// HasCapability 检查是否具有特定能力
func (a *BaseAgent) HasCapability(capability AgentCapability) bool {
	for _, cap := range a.config.Capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}

// AddCapability 添加能力
func (a *BaseAgent) AddCapability(capability AgentCapability) {
	if !a.HasCapability(capability) {
		a.config.Capabilities = append(a.config.Capabilities, capability)
	}
}

// RemoveCapability 移除能力
func (a *BaseAgent) RemoveCapability(capability AgentCapability) {
	var newCapabilities []AgentCapability
	for _, cap := range a.config.Capabilities {
		if cap != capability {
			newCapabilities = append(newCapabilities, cap)
		}
	}
	a.config.Capabilities = newCapabilities
}

// ExecuteTool 执行工具
func (a *BaseAgent) ExecuteTool(ctx context.Context, toolID string, parameters map[string]interface{}) (*tools.ToolResult, error) {
	if a.toolRegistry == nil {
		return nil, NewAgentError(ErrInvalidConfig, "执行工具失败", "工具注册表不可用").
			WithAgentID(a.config.ID)
	}

	tool, err := a.toolRegistry.GetTool(toolID)
	if err != nil {
		return nil, NewAgentError(ErrToolNotFound, "执行工具失败", fmt.Sprintf("工具 %s 未找到", toolID)).
			WithAgentID(a.config.ID)
	}

	if !tool.IsEnabled() {
		return nil, NewAgentError(ErrToolDisabled, "执行工具失败", fmt.Sprintf("工具 %s 已禁用", toolID)).
			WithAgentID(a.config.ID)
	}

	execution := tools.ToolExecution{
		ToolID:     toolID,
		Parameters: parameters,
		Context: map[string]interface{}{
			"agent_id":   a.config.ID,
			"agent_name": a.config.Name,
			"task":       a.execution.Task,
		},
	}

	utils.Info("Agent %s 执行工具: %s", a.config.Name, toolID)
	return tool.Execute(ctx, execution)
}

// GenerateWithModel 使用模型生成内容
func (a *BaseAgent) GenerateWithModel(ctx context.Context, prompt string) (*models.ModelResponse, error) {
	if a.modelAdapter == nil {
		return nil, NewAgentError(ErrInvalidConfig, "调用模型失败", "模型适配器不可用").
			WithAgentID(a.config.ID)
	}

	modelName := a.config.ModelConfig.Name
	if modelName == "" {
		return nil, NewAgentError(ErrInvalidConfig, "调用模型失败", "模型名称未配置").
			WithAgentID(a.config.ID)
	}

	request := models.ModelRequest{
		Prompt: prompt,
	}

	utils.Info("Agent %s 调用模型: %s", a.config.Name, modelName)
	return a.modelAdapter.Generate(ctx, modelName, request)
}

// GenerateStreamWithModel 使用模型流式生成内容
func (a *BaseAgent) GenerateStreamWithModel(ctx context.Context, prompt string, callback func(*models.ModelResponse)) error {
	if a.modelAdapter == nil {
		return NewAgentError(ErrInvalidConfig, "流式调用模型失败", "模型适配器不可用").
			WithAgentID(a.config.ID)
	}

	modelName := a.config.ModelConfig.Name
	if modelName == "" {
		return NewAgentError(ErrInvalidConfig, "流式调用模型失败", "模型名称未配置").
			WithAgentID(a.config.ID)
	}

	request := models.ModelRequest{
		Prompt: prompt,
	}

	utils.Info("Agent %s 流式调用模型: %s", a.config.Name, modelName)
	return a.modelAdapter.GenerateStream(ctx, modelName, request, callback)
}

// FormatPrompt 格式化提示词
func (a *BaseAgent) FormatPrompt(templateName string, data map[string]interface{}) (string, error) {
	template, err := a.GetPromptTemplate(templateName)
	if err != nil {
		return "", err
	}

	// 简单的模板替换
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		valueStr := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, valueStr)
	}

	return result, nil
}
