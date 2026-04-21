package discovery

import (
	"context"
	"fmt"
	"strings"
	"time"

	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/mcp/host"
	"GenPulse/internal/utils"
)

// MCPToolWrapper MCP工具包装器
type MCPToolWrapper struct {
	tool              DiscoveredTool
	host              *host.MCPHost
	callCount         int
	totalDuration     time.Duration
	lastExecutionTime time.Time
	enabled           bool
}

// NewMCPToolWrapper 创建新的MCP工具包装器
func NewMCPToolWrapper(tool DiscoveredTool, host *host.MCPHost) *MCPToolWrapper {
	return &MCPToolWrapper{
		tool:              tool,
		host:              host,
		callCount:         0,
		totalDuration:     0,
		lastExecutionTime: time.Time{},
		enabled:           tool.Enabled,
	}
}

// Execute 执行工具
func (w *MCPToolWrapper) Execute(ctx context.Context, execution tools.ToolExecution) (*tools.ToolResult, error) {
	startTime := time.Now()

	if !w.enabled {
		return nil, fmt.Errorf("tool is disabled: %s", w.tool.FullName)
	}

	if w.host == nil {
		return nil, fmt.Errorf("MCP host not available")
	}

	utils.Info("执行MCP工具: %s", w.tool.FullName)

	// 解析完整名称获取服务器ID和工具名
	serverID, toolName, err := w.parseFullName(w.tool.FullName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tool name: %w", err)
	}

	// 通过MCP主机调用工具
	result, err := w.host.CallTool(serverID, toolName, execution.Parameters)
	if err != nil {
		// 记录失败调用
		if discovery := GetGlobalToolDiscovery(); discovery != nil {
			discovery.RecordToolCall(w.tool.FullName, false)
		}
		return nil, fmt.Errorf("MCP tool call failed: %w", err)
	}

	// 记录成功调用
	if discovery := GetGlobalToolDiscovery(); discovery != nil {
		discovery.RecordToolCall(w.tool.FullName, true)
	}

	// 更新统计信息
	duration := time.Since(startTime)
	w.callCount++
	w.totalDuration += duration
	w.lastExecutionTime = time.Now()

	utils.Debug("MCP工具执行成功: %s (调用次数: %d, 耗时: %v)", w.tool.FullName, w.callCount, duration)

	// 创建工具结果
	toolResult := &tools.ToolResult{
		Success:   true,
		Output:    result,
		Duration:  duration,
		Timestamp: w.lastExecutionTime,
		Metadata: map[string]interface{}{
			"duration_ms": duration.Milliseconds(),
			"server_id":   serverID,
			"source":      "mcp",
		},
	}

	return toolResult, nil
}

// parseFullName 解析完整工具名称
func (w *MCPToolWrapper) parseFullName(fullName string) (string, string, error) {
	// 格式: namespace.toolName
	parts := strings.SplitN(fullName, ".", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid tool name format: %s, expected 'namespace.toolName'", fullName)
	}

	// 暂时使用命名空间作为服务器ID
	// 在实际实现中，可能需要更复杂的映射
	serverID := parts[0]
	toolName := parts[1]

	return serverID, toolName, nil
}

// GetDefinition 获取工具定义
func (w *MCPToolWrapper) GetDefinition() tools.ToolDefinition {
	return tools.ToolDefinition{
		ID:          w.tool.FullName,
		Name:        w.tool.FullName,
		Description: w.tool.ToolInfo.Description,
		Category:    tools.ToolCategoryCustom,
		Version:     "1.0.0",
		Parameters:  w.tool.ToolInfo.InputSchema,
		Returns: map[string]interface{}{
			"type": "object",
		},
		Enabled: w.tool.Enabled,
		Tags:    []string{"mcp", "external"},
		Metadata: map[string]interface{}{
			"mcp_server_id": w.tool.ServerID,
			"source":        "mcp_discovery",
		},
	}
}

// GetType 获取工具类型
func (w *MCPToolWrapper) GetType() string {
	return "mcp_wrapper"
}

// GetCallCount 获取调用次数
func (w *MCPToolWrapper) GetCallCount() int {
	return w.callCount
}

// GetTool 获取原始工具信息
func (w *MCPToolWrapper) GetTool() DiscoveredTool {
	return w.tool
}

// GetCategory 获取工具类别
func (w *MCPToolWrapper) GetCategory() tools.ToolCategory {
	return tools.ToolCategoryCustom
}

// ValidateParameters 验证参数
func (w *MCPToolWrapper) ValidateParameters(parameters map[string]interface{}) error {
	// 简单验证：检查参数是否为map
	if parameters == nil {
		return fmt.Errorf("parameters cannot be nil")
	}
	return nil
}

// Initialize 初始化工具
func (w *MCPToolWrapper) Initialize() error {
	// MCP工具包装器不需要特殊初始化
	return nil
}

// Shutdown 关闭工具
func (w *MCPToolWrapper) Shutdown() error {
	// MCP工具包装器不需要特殊关闭逻辑
	return nil
}

// IsEnabled 检查是否启用
func (w *MCPToolWrapper) IsEnabled() bool {
	return w.enabled
}

// SetEnabled 设置启用状态
func (w *MCPToolWrapper) SetEnabled(enabled bool) {
	w.enabled = enabled
	w.tool.Enabled = enabled
}

// GetExecutionCount 获取执行次数
func (w *MCPToolWrapper) GetExecutionCount() int {
	return w.callCount
}

// GetAverageDuration 获取平均执行时间
func (w *MCPToolWrapper) GetAverageDuration() time.Duration {
	if w.callCount == 0 {
		return 0
	}
	return w.totalDuration / time.Duration(w.callCount)
}

// GetLastExecutionTime 获取最后执行时间
func (w *MCPToolWrapper) GetLastExecutionTime() time.Time {
	return w.lastExecutionTime
}

// Global tool discovery instance
var globalToolDiscovery *ToolDiscovery

// SetGlobalToolDiscovery 设置全局工具发现实例
func SetGlobalToolDiscovery(discovery *ToolDiscovery) {
	globalToolDiscovery = discovery
}

// GetGlobalToolDiscovery 获取全局工具发现实例
func GetGlobalToolDiscovery() *ToolDiscovery {
	return globalToolDiscovery
}
