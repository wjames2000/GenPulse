package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"GenPulse/internal/utils"
)

// ToolCategory 工具类别
type ToolCategory string

const (
	ToolCategoryFileSystem ToolCategory = "filesystem"
	ToolCategoryGit        ToolCategory = "git"
	ToolCategoryShell      ToolCategory = "shell"
	ToolCategoryProject    ToolCategory = "project"
	ToolCategoryNetwork    ToolCategory = "network"
	ToolCategoryDatabase   ToolCategory = "database"
	ToolCategoryUtility    ToolCategory = "utility"
	ToolCategoryCustom     ToolCategory = "custom"
)

// ToolDefinition 工具定义
type ToolDefinition struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    ToolCategory           `json:"category"`
	Version     string                 `json:"version"`
	Parameters  map[string]interface{} `json:"parameters"` // JSON Schema格式
	Returns     map[string]interface{} `json:"returns"`    // JSON Schema格式
	Enabled     bool                   `json:"enabled"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ToolExecution 工具执行
type ToolExecution struct {
	ToolID     string                 `json:"tool_id"`
	Parameters map[string]interface{} `json:"parameters"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

// ToolResult 工具执行结果
type ToolResult struct {
	Success   bool                   `json:"success"`
	Output    interface{}            `json:"output,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Duration  time.Duration          `json:"duration"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Tool interface 工具接口
type Tool interface {
	// 基本信息
	GetDefinition() ToolDefinition
	GetCategory() ToolCategory

	// 执行功能
	Execute(ctx context.Context, execution ToolExecution) (*ToolResult, error)
	ValidateParameters(parameters map[string]interface{}) error

	// 生命周期
	Initialize() error
	Shutdown() error

	// 状态
	IsEnabled() bool
	SetEnabled(enabled bool)

	// 统计信息
	GetExecutionCount() int
	GetAverageDuration() time.Duration
	GetLastExecutionTime() time.Time
}

// BaseTool 基础工具实现
type BaseTool struct {
	definition     ToolDefinition
	executionCount int
	totalDuration  time.Duration
	lastExecution  time.Time
	enabled        bool
	mutex          sync.RWMutex
}

// NewBaseTool 创建基础工具
func NewBaseTool(definition ToolDefinition) *BaseTool {
	return &BaseTool{
		definition: definition,
		enabled:    true,
	}
}

// GetDefinition 获取工具定义
func (t *BaseTool) GetDefinition() ToolDefinition {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.definition
}

// GetCategory 获取工具类别
func (t *BaseTool) GetCategory() ToolCategory {
	return t.definition.Category
}

// Execute 执行工具（需要子类实现）
func (t *BaseTool) Execute(ctx context.Context, execution ToolExecution) (*ToolResult, error) {
	return nil, fmt.Errorf("Execute method must be implemented by subclass")
}

// ValidateParameters 验证参数（基础实现）
func (t *BaseTool) ValidateParameters(parameters map[string]interface{}) error {
	// 基础实现：检查必需参数
	// 子类可以重写此方法提供更复杂的验证

	// 这里可以添加参数验证逻辑
	// 暂时返回nil，表示验证通过
	return nil
}

// Initialize 初始化工具
func (t *BaseTool) Initialize() error {
	// 基础实现：记录初始化日志
	utils.Info("初始化工具: %s (%s)", t.definition.Name, t.definition.Category)
	return nil
}

// Shutdown 关闭工具
func (t *BaseTool) Shutdown() error {
	// 基础实现：记录关闭日志
	utils.Info("关闭工具: %s", t.definition.Name)
	return nil
}

// IsEnabled 检查工具是否启用
func (t *BaseTool) IsEnabled() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.enabled
}

// SetEnabled 设置工具启用状态
func (t *BaseTool) SetEnabled(enabled bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.enabled = enabled
	status := "启用"
	if !enabled {
		status = "禁用"
	}
	utils.Info("%s 工具: %s", status, t.definition.Name)
}

// GetExecutionCount 获取执行次数
func (t *BaseTool) GetExecutionCount() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.executionCount
}

// GetAverageDuration 获取平均执行时间
func (t *BaseTool) GetAverageDuration() time.Duration {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if t.executionCount == 0 {
		return 0
	}
	return t.totalDuration / time.Duration(t.executionCount)
}

// GetLastExecutionTime 获取最后执行时间
func (t *BaseTool) GetLastExecutionTime() time.Time {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.lastExecution
}

// recordExecution 记录执行统计
func (t *BaseTool) recordExecution(duration time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.executionCount++
	t.totalDuration += duration
	t.lastExecution = time.Now()
}

// IncrementExecutionCount 增加执行计数（供子类调用）
func (t *BaseTool) IncrementExecutionCount(duration time.Duration) {
	t.recordExecution(duration)
}

// Mutex 获取互斥锁（供子类调用）
func (t *BaseTool) Mutex() *sync.RWMutex {
	return &t.mutex
}

// ToolRegistry 工具注册表
type ToolRegistry struct {
	tools      map[string]Tool
	categories map[ToolCategory][]string
	mutex      sync.RWMutex
}

// NewToolRegistry 创建工具注册表
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools:      make(map[string]Tool),
		categories: make(map[ToolCategory][]string),
	}
}

// RegisterTool 注册工具
func (tr *ToolRegistry) RegisterTool(tool Tool) error {
	definition := tool.GetDefinition()
	toolID := definition.ID

	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	// 检查是否已注册
	if _, exists := tr.tools[toolID]; exists {
		return fmt.Errorf("tool already registered: %s", toolID)
	}

	// 初始化工具
	if err := tool.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize tool %s: %w", toolID, err)
	}

	// 注册工具
	tr.tools[toolID] = tool

	// 更新类别索引
	category := definition.Category
	tr.categories[category] = append(tr.categories[category], toolID)

	utils.Info("注册工具: %s (%s)", definition.Name, category)
	return nil
}

// UnregisterTool 注销工具
func (tr *ToolRegistry) UnregisterTool(toolID string) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	tool, exists := tr.tools[toolID]
	if !exists {
		return fmt.Errorf("tool not found: %s", toolID)
	}

	// 关闭工具
	if err := tool.Shutdown(); err != nil {
		utils.Warn("关闭工具 %s 时出错: %v", toolID, err)
	}

	// 从注册表中移除
	delete(tr.tools, toolID)

	// 从类别索引中移除
	category := tool.GetCategory()
	if toolIDs, ok := tr.categories[category]; ok {
		for i, id := range toolIDs {
			if id == toolID {
				tr.categories[category] = append(toolIDs[:i], toolIDs[i+1:]...)
				break
			}
		}
	}

	utils.Info("注销工具: %s", toolID)
	return nil
}

// GetTool 获取工具
func (tr *ToolRegistry) GetTool(toolID string) (Tool, error) {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	tool, exists := tr.tools[toolID]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", toolID)
	}

	return tool, nil
}

// ExecuteTool 执行工具
func (tr *ToolRegistry) ExecuteTool(ctx context.Context, execution ToolExecution) (*ToolResult, error) {
	startTime := time.Now()

	// 获取工具
	tool, err := tr.GetTool(execution.ToolID)
	if err != nil {
		return &ToolResult{
			Success:   false,
			Error:     err.Error(),
			Duration:  time.Since(startTime),
			Timestamp: startTime,
		}, err
	}

	// 检查工具是否启用
	if !tool.IsEnabled() {
		err := fmt.Errorf("tool is disabled: %s", execution.ToolID)
		return &ToolResult{
			Success:   false,
			Error:     err.Error(),
			Duration:  time.Since(startTime),
			Timestamp: startTime,
		}, err
	}

	// 验证参数
	if err := tool.ValidateParameters(execution.Parameters); err != nil {
		return &ToolResult{
			Success:   false,
			Error:     fmt.Sprintf("parameter validation failed: %v", err),
			Duration:  time.Since(startTime),
			Timestamp: startTime,
		}, err
	}

	// 执行工具
	result, err := tool.Execute(ctx, execution)
	if err != nil {
		// 记录错误但不返回错误，因为错误信息已经在result中
		utils.Error("工具执行失败 %s: %v", execution.ToolID, err)
	}

	// 确保result不为nil
	if result == nil {
		result = &ToolResult{
			Success: false,
			Error:   "tool execution returned nil result",
		}
	}

	// 设置执行时间
	result.Duration = time.Since(startTime)
	result.Timestamp = startTime

	// 记录执行统计（如果工具是BaseTool类型）
	if baseTool, ok := tool.(*BaseTool); ok {
		baseTool.recordExecution(result.Duration)
	}

	// 记录执行日志
	if result.Success {
		utils.Info("工具执行成功: %s (耗时: %v)", execution.ToolID, result.Duration)
	} else {
		utils.Warn("工具执行失败: %s (错误: %s)", execution.ToolID, result.Error)
	}

	return result, nil
}

// ListTools 列出所有工具
func (tr *ToolRegistry) ListTools() []ToolDefinition {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	var definitions []ToolDefinition
	for _, tool := range tr.tools {
		definitions = append(definitions, tool.GetDefinition())
	}

	return definitions
}

// ListToolsByCategory 按类别列出工具
func (tr *ToolRegistry) ListToolsByCategory(category ToolCategory) []ToolDefinition {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	var definitions []ToolDefinition
	if toolIDs, ok := tr.categories[category]; ok {
		for _, toolID := range toolIDs {
			if tool, exists := tr.tools[toolID]; exists {
				definitions = append(definitions, tool.GetDefinition())
			}
		}
	}

	return definitions
}

// GetToolCount 获取工具数量
func (tr *ToolRegistry) GetToolCount() int {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()
	return len(tr.tools)
}

// GetToolCountByCategory 按类别获取工具数量
func (tr *ToolRegistry) GetToolCountByCategory(category ToolCategory) int {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()
	return len(tr.categories[category])
}

// InitializeAllTools 初始化所有已注册的工具
func (tr *ToolRegistry) InitializeAllTools() error {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	var errors []string
	for toolID, tool := range tr.tools {
		if err := tool.Initialize(); err != nil {
			errors = append(errors, fmt.Sprintf("工具 %s 初始化失败: %v", toolID, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("工具初始化错误: %s", strings.Join(errors, "; "))
	}

	return nil
}

// EnableTool 启用工具
func (tr *ToolRegistry) EnableTool(toolID string) error {
	tool, err := tr.GetTool(toolID)
	if err != nil {
		return err
	}

	tool.SetEnabled(true)
	return nil
}

// DisableTool 禁用工具
func (tr *ToolRegistry) DisableTool(toolID string) error {
	tool, err := tr.GetTool(toolID)
	if err != nil {
		return err
	}

	tool.SetEnabled(false)
	return nil
}

// GetToolStatistics 获取工具统计信息
func (tr *ToolRegistry) GetToolStatistics() map[string]interface{} {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	stats := make(map[string]interface{})
	stats["total_tools"] = len(tr.tools)

	categoryStats := make(map[string]int)
	for category, toolIDs := range tr.categories {
		categoryStats[string(category)] = len(toolIDs)
	}
	stats["tools_by_category"] = categoryStats

	// 计算启用/禁用工具数量
	enabledCount := 0
	for _, tool := range tr.tools {
		if tool.IsEnabled() {
			enabledCount++
		}
	}
	stats["enabled_tools"] = enabledCount
	stats["disabled_tools"] = len(tr.tools) - enabledCount

	return stats
}

// ExportTools 导出工具定义
func (tr *ToolRegistry) ExportTools() ([]byte, error) {
	definitions := tr.ListTools()
	return json.MarshalIndent(definitions, "", "  ")
}

// ImportTools 导入工具定义
func (tr *ToolRegistry) ImportTools(data []byte) ([]string, error) {
	var definitions []ToolDefinition
	if err := json.Unmarshal(data, &definitions); err != nil {
		return nil, fmt.Errorf("failed to parse tool definitions: %w", err)
	}

	var imported []string
	for _, definition := range definitions {
		// 这里需要根据定义创建具体的工具实例
		// 暂时只记录
		imported = append(imported, definition.Name)
	}

	return imported, nil
}

// Global tool registry instance
var globalToolRegistry *ToolRegistry

// InitGlobalToolRegistry 初始化全局工具注册表
func InitGlobalToolRegistry() error {
	if globalToolRegistry == nil {
		globalToolRegistry = NewToolRegistry()
		utils.Info("初始化全局工具注册表")
	}
	return nil
}

// GetGlobalToolRegistry 获取全局工具注册表
func GetGlobalToolRegistry() *ToolRegistry {
	if globalToolRegistry == nil {
		// 自动初始化
		InitGlobalToolRegistry()
	}
	return globalToolRegistry
}
