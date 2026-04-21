package discovery

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/mcp/client"
	"GenPulse/internal/mcp/host"
	"GenPulse/internal/utils"
)

// ToolDiscovery 工具发现服务
type ToolDiscovery struct {
	host          *host.MCPHost
	toolRegistry  *tools.ToolRegistry
	discovered    map[string]DiscoveredTool // toolFullName -> DiscoveredTool
	mu            sync.RWMutex
	running       bool
	cancel        context.CancelFunc
	discoveryChan chan ToolDiscoveryEvent
}

// DiscoveredTool 发现的工具
type DiscoveredTool struct {
	FullName     string          // 格式: namespace.toolName
	ServerID     string          // 服务器ID
	ToolInfo     client.ToolInfo // 工具信息
	LastSeen     time.Time       // 最后发现时间
	CallCount    int             // 调用次数
	LastCallTime time.Time       // 最后调用时间
	SuccessRate  float64         // 成功率
	Enabled      bool            // 是否启用
}

// ToolDiscoveryEvent 工具发现事件
type ToolDiscoveryEvent struct {
	Type      string         // "added", "removed", "updated"
	Tool      DiscoveredTool // 工具信息
	Timestamp time.Time      // 事件时间
}

// NewToolDiscovery 创建新的工具发现服务
func NewToolDiscovery(mcpHost *host.MCPHost, toolRegistry *tools.ToolRegistry) *ToolDiscovery {
	return &ToolDiscovery{
		host:          mcpHost,
		toolRegistry:  toolRegistry,
		discovered:    make(map[string]DiscoveredTool),
		running:       false,
		discoveryChan: make(chan ToolDiscoveryEvent, 100),
	}
}

// Start 启动工具发现服务
func (td *ToolDiscovery) Start(ctx context.Context) error {
	td.mu.Lock()
	defer td.mu.Unlock()

	if td.running {
		return fmt.Errorf("tool discovery already running")
	}

	utils.Info("启动MCP工具发现服务")

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(ctx)
	td.cancel = cancel

	// 启动发现循环
	go td.discoveryLoop(ctx)

	td.running = true
	utils.Info("MCP工具发现服务已启动")

	return nil
}

// discoveryLoop 发现循环
func (td *ToolDiscovery) discoveryLoop(ctx context.Context) {
	// 初始发现
	td.performDiscovery(ctx)

	// 定期发现
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			utils.Info("工具发现循环停止")
			return
		case <-ticker.C:
			td.performDiscovery(ctx)
		}
	}
}

// performDiscovery 执行工具发现
func (td *ToolDiscovery) performDiscovery(ctx context.Context) {
	utils.Debug("执行MCP工具发现")

	if td.host == nil {
		utils.Warn("MCP主机未初始化，跳过工具发现")
		return
	}

	// 获取所有工具
	allTools := td.host.ListAllTools()

	td.mu.Lock()
	defer td.mu.Unlock()

	// 记录当前发现的所有工具
	currentTools := make(map[string]bool)

	// 处理每个命名空间的工具
	for namespace, toolList := range allTools {
		for _, toolInfo := range toolList {
			fullName := fmt.Sprintf("%s.%s", namespace, toolInfo.Name)
			currentTools[fullName] = true

			// 检查是否是新工具
			if existing, exists := td.discovered[fullName]; exists {
				// 工具已存在，检查是否有更新
				if td.hasToolChanged(existing.ToolInfo, toolInfo) {
					// 工具已更新
					updatedTool := DiscoveredTool{
						FullName:     fullName,
						ServerID:     existing.ServerID, // 需要找到对应的ServerID
						ToolInfo:     toolInfo,
						LastSeen:     time.Now(),
						CallCount:    existing.CallCount,
						LastCallTime: existing.LastCallTime,
						SuccessRate:  existing.SuccessRate,
						Enabled:      existing.Enabled,
					}

					td.discovered[fullName] = updatedTool

					// 发送更新事件
					td.sendDiscoveryEvent(ToolDiscoveryEvent{
						Type:      "updated",
						Tool:      updatedTool,
						Timestamp: time.Now(),
					})

					utils.Debug("工具已更新: %s", fullName)
				} else {
					// 工具未变化，更新最后发现时间
					existing.LastSeen = time.Now()
					td.discovered[fullName] = existing
				}
			} else {
				// 新工具
				// 需要找到对应的ServerID
				serverID := td.findServerIDForNamespace(namespace)

				newTool := DiscoveredTool{
					FullName:    fullName,
					ServerID:    serverID,
					ToolInfo:    toolInfo,
					LastSeen:    time.Now(),
					CallCount:   0,
					SuccessRate: 1.0,
					Enabled:     true,
				}

				td.discovered[fullName] = newTool

				// 发送添加事件
				td.sendDiscoveryEvent(ToolDiscoveryEvent{
					Type:      "added",
					Tool:      newTool,
					Timestamp: time.Now(),
				})

				utils.Info("发现新工具: %s (来自服务器: %s)", fullName, serverID)

				// 自动注册到工具注册表
				if td.toolRegistry != nil {
					td.registerDiscoveredTool(newTool)
				}
			}
		}
	}

	// 检查是否有工具被移除
	for fullName, tool := range td.discovered {
		if !currentTools[fullName] {
			// 工具被移除
			tool.LastSeen = time.Now()
			td.discovered[fullName] = tool

			// 发送移除事件
			td.sendDiscoveryEvent(ToolDiscoveryEvent{
				Type:      "removed",
				Tool:      tool,
				Timestamp: time.Now(),
			})

			utils.Info("工具已移除: %s", fullName)

			// 从工具注册表移除
			if td.toolRegistry != nil {
				td.unregisterDiscoveredTool(tool)
			}
		}
	}

	utils.Debug("工具发现完成，当前发现 %d 个工具", len(td.discovered))
}

// findServerIDForNamespace 根据命名空间查找服务器ID
func (td *ToolDiscovery) findServerIDForNamespace(namespace string) string {
	// 这里需要根据命名空间查找对应的服务器ID
	// 暂时返回命名空间作为服务器ID
	return namespace
}

// hasToolChanged 检查工具是否有变化
func (td *ToolDiscovery) hasToolChanged(oldTool, newTool client.ToolInfo) bool {
	if oldTool.Name != newTool.Name {
		return true
	}

	if oldTool.Description != newTool.Description {
		return true
	}

	// 可以添加更多比较逻辑
	return false
}

// sendDiscoveryEvent 发送发现事件
func (td *ToolDiscovery) sendDiscoveryEvent(event ToolDiscoveryEvent) {
	select {
	case td.discoveryChan <- event:
		// 事件已发送
	default:
		// 通道已满，丢弃事件
		utils.Warn("工具发现事件通道已满，丢弃事件: %s %s", event.Type, event.Tool.FullName)
	}
}

// registerDiscoveredTool 注册发现的工具到工具注册表
func (td *ToolDiscovery) registerDiscoveredTool(tool DiscoveredTool) {
	if td.toolRegistry == nil {
		return
	}

	// 创建MCP工具包装器
	mcpTool := NewMCPToolWrapper(tool, td.host)

	// 注册工具
	if err := td.toolRegistry.RegisterTool(mcpTool); err != nil {
		utils.Warn("注册MCP工具失败 %s: %v", tool.FullName, err)
	} else {
		utils.Debug("已注册MCP工具: %s", tool.FullName)
	}
}

// unregisterDiscoveredTool 从工具注册表移除工具
func (td *ToolDiscovery) unregisterDiscoveredTool(tool DiscoveredTool) {
	if td.toolRegistry == nil {
		return
	}

	// 从工具注册表移除
	td.toolRegistry.UnregisterTool(tool.FullName)
	utils.Debug("已移除MCP工具: %s", tool.FullName)
}

// Stop 停止工具发现服务
func (td *ToolDiscovery) Stop() error {
	td.mu.Lock()
	defer td.mu.Unlock()

	if !td.running {
		return nil
	}

	utils.Info("停止MCP工具发现服务")

	// 取消上下文
	if td.cancel != nil {
		td.cancel()
	}

	// 关闭事件通道
	close(td.discoveryChan)

	// 清理发现的工具
	for _, tool := range td.discovered {
		td.unregisterDiscoveredTool(tool)
	}

	td.discovered = make(map[string]DiscoveredTool)
	td.running = false

	utils.Info("MCP工具发现服务已停止")
	return nil
}

// IsRunning 检查是否正在运行
func (td *ToolDiscovery) IsRunning() bool {
	td.mu.RLock()
	defer td.mu.RUnlock()
	return td.running
}

// GetDiscoveredTools 获取所有发现的工具
func (td *ToolDiscovery) GetDiscoveredTools() []DiscoveredTool {
	td.mu.RLock()
	defer td.mu.RUnlock()

	tools := make([]DiscoveredTool, 0, len(td.discovered))
	for _, tool := range td.discovered {
		tools = append(tools, tool)
	}

	return tools
}

// GetTool 获取特定工具
func (td *ToolDiscovery) GetTool(fullName string) (*DiscoveredTool, error) {
	td.mu.RLock()
	defer td.mu.RUnlock()

	tool, exists := td.discovered[fullName]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", fullName)
	}

	return &tool, nil
}

// EnableTool 启用工具
func (td *ToolDiscovery) EnableTool(fullName string) error {
	td.mu.Lock()
	defer td.mu.Unlock()

	tool, exists := td.discovered[fullName]
	if !exists {
		return fmt.Errorf("tool not found: %s", fullName)
	}

	if tool.Enabled {
		return nil // 已经启用
	}

	tool.Enabled = true
	td.discovered[fullName] = tool

	// 重新注册到工具注册表
	if td.toolRegistry != nil {
		td.registerDiscoveredTool(tool)
	}

	utils.Info("已启用工具: %s", fullName)
	return nil
}

// DisableTool 禁用工具
func (td *ToolDiscovery) DisableTool(fullName string) error {
	td.mu.Lock()
	defer td.mu.Unlock()

	tool, exists := td.discovered[fullName]
	if !exists {
		return fmt.Errorf("tool not found: %s", fullName)
	}

	if !tool.Enabled {
		return nil // 已经禁用
	}

	tool.Enabled = false
	td.discovered[fullName] = tool

	// 从工具注册表移除
	if td.toolRegistry != nil {
		td.unregisterDiscoveredTool(tool)
	}

	utils.Info("已禁用工具: %s", fullName)
	return nil
}

// RecordToolCall 记录工具调用
func (td *ToolDiscovery) RecordToolCall(fullName string, success bool) error {
	td.mu.Lock()
	defer td.mu.Unlock()

	tool, exists := td.discovered[fullName]
	if !exists {
		return fmt.Errorf("tool not found: %s", fullName)
	}

	// 更新调用统计
	tool.CallCount++
	tool.LastCallTime = time.Now()

	// 更新成功率
	totalCalls := float64(tool.CallCount)
	if success {
		// 简化计算：新的成功率 = (旧的成功率 * (总调用数-1) + 1) / 总调用数
		if tool.CallCount == 1 {
			tool.SuccessRate = 1.0
		} else {
			tool.SuccessRate = (tool.SuccessRate*(totalCalls-1) + 1.0) / totalCalls
		}
	} else {
		// 调用失败
		if tool.CallCount == 1 {
			tool.SuccessRate = 0.0
		} else {
			tool.SuccessRate = (tool.SuccessRate * (totalCalls - 1)) / totalCalls
		}
	}

	td.discovered[fullName] = tool
	return nil
}

// GetDiscoveryEvents 获取发现事件通道
func (td *ToolDiscovery) GetDiscoveryEvents() <-chan ToolDiscoveryEvent {
	return td.discoveryChan
}

// SearchTools 搜索工具
func (td *ToolDiscovery) SearchTools(query string) []DiscoveredTool {
	td.mu.RLock()
	defer td.mu.RUnlock()

	// 编译正则表达式（不区分大小写）
	pattern, err := regexp.Compile("(?i)" + regexp.QuoteMeta(query))
	if err != nil {
		// 如果正则表达式编译失败，使用简单字符串匹配
		return td.searchToolsSimple(query)
	}

	results := make([]DiscoveredTool, 0)

	for _, tool := range td.discovered {
		// 检查工具名称或描述是否匹配
		if pattern.MatchString(tool.FullName) || pattern.MatchString(tool.ToolInfo.Description) {
			results = append(results, tool)
		}
	}

	return results
}

// searchToolsSimple 简单字符串搜索
func (td *ToolDiscovery) searchToolsSimple(query string) []DiscoveredTool {
	results := make([]DiscoveredTool, 0)
	queryLower := ""
	if query != "" {
		queryLower = query
		// 简单实现：转换为小写进行比较
	}

	for _, tool := range td.discovered {
		// 简单字符串包含检查
		matches := query == "" ||
			containsIgnoreCase(tool.FullName, queryLower) ||
			containsIgnoreCase(tool.ToolInfo.Description, queryLower)

		if matches {
			results = append(results, tool)
		}
	}

	return results
}

// containsIgnoreCase 简单的不区分大小写包含检查
func containsIgnoreCase(s, substr string) bool {
	// 简单实现：转换为小写
	// 注意：这不能正确处理Unicode，但对于简单搜索足够了
	if len(s) < len(substr) {
		return false
	}

	// 简单实现
	for i := 0; i <= len(s)-len(substr); i++ {
		if strings.EqualFold(s[i:i+len(substr)], substr) {
			return true
		}
	}

	return false
}

// GetToolStatistics 获取工具统计信息
func (td *ToolDiscovery) GetToolStatistics() map[string]interface{} {
	td.mu.RLock()
	defer td.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["total_tools"] = len(td.discovered)

	enabledCount := 0
	totalCalls := 0
	for _, tool := range td.discovered {
		if tool.Enabled {
			enabledCount++
		}
		totalCalls += tool.CallCount
	}

	stats["enabled_tools"] = enabledCount
	stats["disabled_tools"] = len(td.discovered) - enabledCount
	stats["total_calls"] = totalCalls

	// 计算平均成功率
	if len(td.discovered) > 0 {
		totalSuccessRate := 0.0
		for _, tool := range td.discovered {
			totalSuccessRate += tool.SuccessRate
		}
		stats["average_success_rate"] = totalSuccessRate / float64(len(td.discovered))
	} else {
		stats["average_success_rate"] = 0.0
	}

	return stats
}
