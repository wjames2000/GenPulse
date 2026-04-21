package host

import (
	"context"
	"fmt"
	"sync"
	"time"

	"GenPulse/internal/mcp/client"
	"GenPulse/internal/mcp/server"
	"GenPulse/internal/utils"
)

// MCPHostConfig MCP主机配置
type MCPHostConfig struct {
	// Servers 服务器配置列表
	Servers []MCPHostServerConfig `json:"servers"`

	// AutoStart 是否自动启动
	AutoStart bool `json:"auto_start,omitempty"`

	// ToolDiscoveryInterval 工具发现间隔（秒）
	ToolDiscoveryInterval int `json:"tool_discovery_interval,omitempty"`

	// MaxConcurrentCalls 最大并发调用数
	MaxConcurrentCalls int `json:"max_concurrent_calls,omitempty"`
}

// MCPHostServerConfig MCP主机服务器配置
type MCPHostServerConfig struct {
	// ID 服务器唯一标识
	ID string `json:"id"`

	// Name 服务器名称
	Name string `json:"name"`

	// Type 服务器类型: "client" 或 "server"
	Type string `json:"type"`

	// ClientConfig 客户端配置（当Type为"client"时）
	ClientConfig client.MCPClientConfig `json:"client_config,omitempty"`

	// ServerConfig 服务器配置（当Type为"server"时）
	ServerConfig MCPServerConfig `json:"server_config,omitempty"`

	// Enabled 是否启用
	Enabled bool `json:"enabled,omitempty"`

	// Priority 优先级（0-100，越高越优先）
	Priority int `json:"priority,omitempty"`
}

// MCPServerConfig MCP服务器配置
type MCPServerConfig struct {
	// Type 服务器类型: "stdio"
	Type string `json:"type"`

	// ToolFilter 工具过滤器（正则表达式）
	ToolFilter string `json:"tool_filter,omitempty"`
}

// MCPHost MCP主机管理器
type MCPHost struct {
	config    MCPHostConfig
	clients   map[string]client.MCPClient
	servers   map[string]server.MCPServer
	toolCache map[string][]client.ToolInfo // 缓存工具列表：serverID -> tools
	mu        sync.RWMutex
	running   bool
	cancel    context.CancelFunc
}

// NewMCPHost 创建新的MCP主机
func NewMCPHost(config MCPHostConfig) *MCPHost {
	return &MCPHost{
		config:    config,
		clients:   make(map[string]client.MCPClient),
		servers:   make(map[string]server.MCPServer),
		toolCache: make(map[string][]client.ToolInfo),
		running:   false,
	}
}

// Start 启动MCP主机
func (h *MCPHost) Start(ctx context.Context) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.running {
		return fmt.Errorf("MCP host already running")
	}

	utils.Info("启动MCP主机管理器")

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(ctx)
	h.cancel = cancel

	// 启动所有启用的服务器
	for _, serverConfig := range h.config.Servers {
		if !serverConfig.Enabled {
			continue
		}

		if err := h.startServer(ctx, serverConfig); err != nil {
			utils.Warn("启动服务器 %s 失败: %v", serverConfig.Name, err)
			continue
		}
	}

	// 启动工具发现循环
	if h.config.ToolDiscoveryInterval > 0 {
		go h.toolDiscoveryLoop(ctx)
	}

	h.running = true
	utils.Info("MCP主机管理器已启动，管理 %d 个服务器", len(h.clients)+len(h.servers))

	return nil
}

// startServer 启动单个服务器
func (h *MCPHost) startServer(ctx context.Context, config MCPHostServerConfig) error {
	switch config.Type {
	case "client":
		return h.startClient(ctx, config)
	case "server":
		return h.startLocalServer(ctx, config)
	default:
		return fmt.Errorf("unknown server type: %s", config.Type)
	}
}

// startClient 启动MCP客户端
func (h *MCPHost) startClient(ctx context.Context, config MCPHostServerConfig) error {
	utils.Info("启动MCP客户端: %s (%s)", config.Name, config.ID)

	// 创建客户端
	mcpClient, err := client.NewMCPClient(config.ClientConfig)
	if err != nil {
		return fmt.Errorf("failed to create MCP client: %w", err)
	}

	// 连接到服务器
	if err := mcpClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to MCP server: %w", err)
	}

	// 获取工具列表
	tools, err := mcpClient.ListTools()
	if err != nil {
		utils.Warn("获取MCP客户端工具列表失败: %v", err)
		tools = []client.ToolInfo{}
	}

	// 缓存工具列表
	h.toolCache[config.ID] = tools

	// 存储客户端
	h.clients[config.ID] = mcpClient

	utils.Info("MCP客户端启动成功: %s，发现 %d 个工具", config.Name, len(tools))
	return nil
}

// startLocalServer 启动本地MCP服务器
func (h *MCPHost) startLocalServer(ctx context.Context, config MCPHostServerConfig) error {
	utils.Info("启动本地MCP服务器: %s (%s)", config.Name, config.ID)

	// TODO: 需要从全局获取工具注册表
	// 暂时创建空的工具注册表
	// toolRegistry := tools.GetGlobalToolRegistry()

	// 创建服务器
	mcpServer, err := server.NewMCPServer(nil, config.ServerConfig.Type) // 暂时传递nil
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	// 启动服务器
	if err := mcpServer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	// 获取工具列表
	tools, err := mcpServer.ExportTools()
	if err != nil {
		utils.Warn("导出MCP服务器工具失败: %v", err)
	}

	// 转换为客户端工具格式
	clientTools := make([]client.ToolInfo, 0, len(tools))
	for _, tool := range tools {
		clientTools = append(clientTools, client.ToolInfo{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}

	// 缓存工具列表
	h.toolCache[config.ID] = clientTools

	// 存储服务器
	h.servers[config.ID] = mcpServer

	utils.Info("本地MCP服务器启动成功: %s，导出 %d 个工具", config.Name, len(clientTools))
	return nil
}

// toolDiscoveryLoop 工具发现循环
func (h *MCPHost) toolDiscoveryLoop(ctx context.Context) {
	interval := h.config.ToolDiscoveryInterval
	if interval <= 0 {
		interval = 60 // 默认60秒
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	utils.Info("工具发现循环启动，间隔 %d 秒", interval)

	for {
		select {
		case <-ctx.Done():
			utils.Info("工具发现循环停止")
			return
		case <-ticker.C:
			h.discoverTools(ctx)
		}
	}
}

// discoverTools 发现所有服务器的工具
func (h *MCPHost) discoverTools(ctx context.Context) {
	h.mu.RLock()
	// 创建客户端ID和客户端的映射
	clientMap := make(map[string]client.MCPClient)
	for id, client := range h.clients {
		clientMap[id] = client
	}
	h.mu.RUnlock()

	utils.Debug("开始工具发现")

	// 批量收集更新
	updates := make(map[string][]client.ToolInfo)

	for id, client := range clientMap {
		if !client.IsConnected() {
			continue
		}

		// 获取工具列表
		tools, err := client.ListTools()
		if err != nil {
			utils.Debug("获取客户端 %s 工具列表失败: %v", id, err)
			continue
		}

		updates[id] = tools
	}

	// 批量更新缓存
	if len(updates) > 0 {
		h.mu.Lock()
		for id, tools := range updates {
			h.toolCache[id] = tools
			utils.Debug("更新客户端 %s 的工具缓存，发现 %d 个工具", id, len(tools))
		}
		h.mu.Unlock()
	}
}

// Stop 停止MCP主机
func (h *MCPHost) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.running {
		return nil
	}

	utils.Info("停止MCP主机管理器")

	// 取消上下文
	if h.cancel != nil {
		h.cancel()
	}

	// 断开所有客户端连接
	for id, client := range h.clients {
		if err := client.Disconnect(); err != nil {
			utils.Warn("断开客户端 %s 连接失败: %v", id, err)
		}
	}

	// 停止所有服务器
	for id, server := range h.servers {
		if err := server.Stop(); err != nil {
			utils.Warn("停止服务器 %s 失败: %v", id, err)
		}
	}

	// 清理资源
	h.clients = make(map[string]client.MCPClient)
	h.servers = make(map[string]server.MCPServer)
	h.toolCache = make(map[string][]client.ToolInfo)
	h.running = false

	utils.Info("MCP主机管理器已停止")
	return nil
}

// IsRunning 检查是否正在运行
func (h *MCPHost) IsRunning() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.running
}

// ListAllTools 列出所有工具（按命名空间分组）
func (h *MCPHost) ListAllTools() map[string][]client.ToolInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[string][]client.ToolInfo)

	// 收集客户端工具
	for id, client := range h.clients {
		if tools, ok := h.toolCache[id]; ok {
			namespace := client.GetNamespace()
			if namespace == "" {
				namespace = id
			}
			result[namespace] = append(result[namespace], tools...)
		}
	}

	// 收集服务器工具
	for id, tools := range h.toolCache {
		if _, isClient := h.clients[id]; !isClient {
			// 这是服务器工具
			result[id] = tools
		}
	}

	return result
}

// CallTool 调用工具
func (h *MCPHost) CallTool(serverID, toolName string, arguments map[string]interface{}) (interface{}, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 首先在客户端中查找
	if client, ok := h.clients[serverID]; ok {
		if !client.IsConnected() {
			return nil, fmt.Errorf("client %s is not connected", serverID)
		}

		utils.Info("通过客户端 %s 调用工具: %s", serverID, toolName)
		return client.CallTool(toolName, arguments)
	}

	// 然后在服务器中查找
	if server, ok := h.servers[serverID]; ok {
		if !server.IsRunning() {
			return nil, fmt.Errorf("server %s is not running", serverID)
		}

		utils.Info("通过服务器 %s 调用工具: %s", serverID, toolName)
		// 服务器需要不同的调用方式
		return nil, fmt.Errorf("server tool calling not implemented yet")
	}

	return nil, fmt.Errorf("server not found: %s", serverID)
}

// AddClient 添加MCP客户端
func (h *MCPHost) AddClient(ctx context.Context, config MCPHostServerConfig) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[config.ID]; exists {
		return fmt.Errorf("client with ID %s already exists", config.ID)
	}

	if _, exists := h.servers[config.ID]; exists {
		return fmt.Errorf("server with ID %s already exists", config.ID)
	}

	return h.startClient(ctx, config)
}

// AddServer 添加MCP服务器
func (h *MCPHost) AddServer(ctx context.Context, config MCPHostServerConfig) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[config.ID]; exists {
		return fmt.Errorf("client with ID %s already exists", config.ID)
	}

	if _, exists := h.servers[config.ID]; exists {
		return fmt.Errorf("server with ID %s already exists", config.ID)
	}

	return h.startLocalServer(ctx, config)
}

// RemoveServer 移除服务器
func (h *MCPHost) RemoveServer(serverID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 检查是否是客户端
	if client, ok := h.clients[serverID]; ok {
		if err := client.Disconnect(); err != nil {
			return fmt.Errorf("failed to disconnect client: %w", err)
		}
		delete(h.clients, serverID)
		delete(h.toolCache, serverID)
		utils.Info("已移除MCP客户端: %s", serverID)
		return nil
	}

	// 检查是否是服务器
	if server, ok := h.servers[serverID]; ok {
		if err := server.Stop(); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
		delete(h.servers, serverID)
		delete(h.toolCache, serverID)
		utils.Info("已移除MCP服务器: %s", serverID)
		return nil
	}

	return fmt.Errorf("server not found: %s", serverID)
}

// UpdateServer 更新服务器配置
func (h *MCPHost) UpdateServer(ctx context.Context, serverID string, config MCPHostServerConfig) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 确保ID一致
	config.ID = serverID

	// 检查服务器是否存在
	_, isClient := h.clients[serverID]
	_, isServer := h.servers[serverID]

	if !isClient && !isServer {
		return fmt.Errorf("server not found: %s", serverID)
	}

	// 先移除旧的服务器
	if isClient {
		if client, ok := h.clients[serverID]; ok {
			if err := client.Disconnect(); err != nil {
				return fmt.Errorf("failed to disconnect client: %w", err)
			}
			delete(h.clients, serverID)
			delete(h.toolCache, serverID)
		}
	} else if isServer {
		if server, ok := h.servers[serverID]; ok {
			if err := server.Stop(); err != nil {
				return fmt.Errorf("failed to stop server: %w", err)
			}
			delete(h.servers, serverID)
			delete(h.toolCache, serverID)
		}
	}

	// 根据类型启动新的服务器
	if config.Type == "client" {
		return h.startClient(ctx, config)
	} else if config.Type == "server" {
		return h.startLocalServer(ctx, config)
	}

	return fmt.Errorf("unknown server type: %s", config.Type)
}

// GetServerStatus 获取服务器状态
func (h *MCPHost) GetServerStatus(serverID string) (map[string]interface{}, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	status := make(map[string]interface{})

	// 检查客户端
	if client, ok := h.clients[serverID]; ok {
		status["type"] = "client"
		status["connected"] = client.IsConnected()
		status["config"] = client.GetConfig()
		if tools, ok := h.toolCache[serverID]; ok {
			status["tool_count"] = len(tools)
		}
		return status, nil
	}

	// 检查服务器
	if server, ok := h.servers[serverID]; ok {
		status["type"] = "server"
		status["running"] = server.IsRunning()
		if tools, ok := h.toolCache[serverID]; ok {
			status["tool_count"] = len(tools)
		}
		return status, nil
	}

	return nil, fmt.Errorf("server not found: %s", serverID)
}

// GetConfig 获取配置
func (h *MCPHost) GetConfig() MCPHostConfig {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.config
}

// UpdateConfig 更新配置
func (h *MCPHost) UpdateConfig(config MCPHostConfig) error {
	h.mu.Lock()

	// 停止当前运行的主机
	if h.running {
		// 先释放锁，避免死锁
		h.mu.Unlock()
		if err := h.Stop(); err != nil {
			h.mu.Lock()
			return fmt.Errorf("failed to stop current host: %w", err)
		}
		h.mu.Lock()
	}

	// 更新配置
	h.config = config

	// 重新启动
	ctx := context.Background()
	err := h.Start(ctx)
	h.mu.Unlock()

	return err
}

// UpdateConfigWithCallback 更新配置并执行回调
func (h *MCPHost) UpdateConfigWithCallback(config MCPHostConfig, callback func(MCPHostConfig) error) error {
	h.mu.Lock()

	// 停止当前运行的主机
	if h.running {
		// 先释放锁，避免死锁
		h.mu.Unlock()
		if err := h.Stop(); err != nil {
			h.mu.Lock()
			return fmt.Errorf("failed to stop current host: %w", err)
		}
		h.mu.Lock()
	}

	// 更新配置
	h.config = config

	// 执行回调（例如保存到文件）
	if callback != nil {
		if err := callback(config); err != nil {
			h.mu.Unlock()
			return fmt.Errorf("failed to execute callback: %w", err)
		}
	}

	// 重新启动
	ctx := context.Background()
	err := h.Start(ctx)
	h.mu.Unlock()

	return err
}
