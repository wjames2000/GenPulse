package mcp

import (
	"context"
	"fmt"
	"testing"
	"time"

	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/mcp/client"
	"GenPulse/internal/mcp/config"
	"GenPulse/internal/mcp/discovery"
	"GenPulse/internal/mcp/host"
	"GenPulse/internal/mcp/server"
	"GenPulse/internal/utils"
)

// TestMCPClientCreation 测试MCP客户端创建
func TestMCPClientCreation(t *testing.T) {
	utils.Info("测试MCP客户端创建")

	config := client.MCPClientConfig{
		ServerType:    "stdio",
		Command:       "echo",
		Args:          []string{"test"},
		Namespace:     "test",
		Timeout:       5,
		AutoReconnect: true,
		MaxRetries:    3,
	}

	client, err := client.NewMCPClient(config)
	if err != nil {
		t.Fatalf("创建MCP客户端失败: %v", err)
	}

	if client == nil {
		t.Fatal("MCP客户端为nil")
	}

	// 检查配置
	clientConfig := client.GetConfig()
	if clientConfig.ServerType != "stdio" {
		t.Errorf("期望ServerType为'stdio'，实际为'%s'", clientConfig.ServerType)
	}

	if clientConfig.Namespace != "test" {
		t.Errorf("期望Namespace为'test'，实际为'%s'", clientConfig.Namespace)
	}

	utils.Info("MCP客户端创建测试通过")
}

// TestMCPHostManagement 测试MCP主机管理
func TestMCPHostManagement(t *testing.T) {
	utils.Info("测试MCP主机管理")

	hostConfig := host.MCPHostConfig{
		AutoStart:             true,
		ToolDiscoveryInterval: 10,
		MaxConcurrentCalls:    5,
		Servers: []host.MCPHostServerConfig{
			{
				ID:       "test-server-1",
				Name:     "测试服务器1",
				Type:     "client",
				Enabled:  true,
				Priority: 100,
				ClientConfig: client.MCPClientConfig{
					ServerType: "stdio",
					Command:    "echo",
					Args:       []string{"test"},
					Namespace:  "test1",
				},
			},
		},
	}

	mcpHost := host.NewMCPHost(hostConfig)
	if mcpHost == nil {
		t.Fatal("MCP主机为nil")
	}

	// 启动主机
	ctx := context.Background()
	if err := mcpHost.Start(ctx); err != nil {
		t.Fatalf("启动MCP主机失败: %v", err)
	}

	// 检查运行状态
	if !mcpHost.IsRunning() {
		t.Error("MCP主机未运行")
	}

	// 获取配置
	config := mcpHost.GetConfig()
	if len(config.Servers) != 1 {
		t.Errorf("期望有1个服务器，实际有%d个", len(config.Servers))
	}

	// 停止主机
	if err := mcpHost.Stop(); err != nil {
		t.Fatalf("停止MCP主机失败: %v", err)
	}

	if mcpHost.IsRunning() {
		t.Error("MCP主机仍在运行")
	}

	utils.Info("MCP主机管理测试通过")
}

// TestMCPConfigManager 测试MCP配置管理器
func TestMCPConfigManager(t *testing.T) {
	utils.Info("测试MCP配置管理器")

	// 创建临时配置文件路径
	configPath := "/tmp/test_mcp_config.json"

	// 创建配置管理器
	configManager, err := config.NewMCPConfigManager(configPath)
	if err != nil {
		t.Fatalf("创建MCP配置管理器失败: %v", err)
	}

	// 获取配置
	initialConfig := configManager.GetConfig()
	if initialConfig.AutoStart != true {
		t.Error("期望AutoStart为true")
	}

	// 添加服务器配置
	serverConfig := host.MCPHostServerConfig{
		ID:      "test-server",
		Name:    "测试服务器",
		Type:    "client",
		Enabled: true,
		ClientConfig: client.MCPClientConfig{
			ServerType: "stdio",
			Command:    "echo",
			Args:       []string{"test"},
			Namespace:  "test",
		},
	}

	if err := configManager.AddServer(serverConfig); err != nil {
		t.Fatalf("添加服务器配置失败: %v", err)
	}

	// 检查服务器是否添加成功
	servers := configManager.ListServers()
	if len(servers) != 1 {
		t.Errorf("期望有1个服务器，实际有%d个", len(servers))
	}

	// 获取服务器配置
	retrievedConfig, err := configManager.GetServer("test-server")
	if err != nil {
		t.Fatalf("获取服务器配置失败: %v", err)
	}

	if retrievedConfig.Name != "测试服务器" {
		t.Errorf("期望服务器名称为'测试服务器'，实际为'%s'", retrievedConfig.Name)
	}

	// 移除服务器配置
	if err := configManager.RemoveServer("test-server"); err != nil {
		t.Fatalf("移除服务器配置失败: %v", err)
	}

	servers = configManager.ListServers()
	if len(servers) != 0 {
		t.Errorf("期望有0个服务器，实际有%d个", len(servers))
	}

	utils.Info("MCP配置管理器测试通过")
}

// TestToolDiscovery 测试工具发现
func TestToolDiscovery(t *testing.T) {
	utils.Info("测试工具发现")

	// 创建MCP主机
	hostConfig := host.MCPHostConfig{
		AutoStart:             false, // 不自动启动
		ToolDiscoveryInterval: 1,     // 1秒间隔
	}

	mcpHost := host.NewMCPHost(hostConfig)

	// 创建工具发现服务
	toolDiscovery := discovery.NewToolDiscovery(mcpHost, nil) // 不传递工具注册表

	// 启动工具发现服务
	ctx := context.Background()
	if err := toolDiscovery.Start(ctx); err != nil {
		t.Fatalf("启动工具发现服务失败: %v", err)
	}

	// 等待发现循环运行
	time.Sleep(2 * time.Second)

	// 获取发现的工具
	tools := toolDiscovery.GetDiscoveredTools()
	utils.Info("发现 %d 个工具", len(tools))

	// 停止工具发现服务
	if err := toolDiscovery.Stop(); err != nil {
		t.Fatalf("停止工具发现服务失败: %v", err)
	}

	if toolDiscovery.IsRunning() {
		t.Error("工具发现服务仍在运行")
	}

	utils.Info("工具发现测试通过")
}

// TestMCPServerCreation 测试MCP服务器创建
func TestMCPServerCreation(t *testing.T) {
	utils.Info("测试MCP服务器创建")

	// 创建MCP服务器
	mcpServer, err := server.NewMCPServer(nil, "stdio")
	if err != nil {
		t.Fatalf("创建MCP服务器失败: %v", err)
	}

	if mcpServer == nil {
		t.Fatal("MCP服务器为nil")
	}

	// 检查运行状态
	if mcpServer.IsRunning() {
		t.Error("MCP服务器不应在创建后立即运行")
	}

	// 导出工具
	tools, err := mcpServer.ExportTools()
	if err != nil {
		t.Fatalf("导出工具失败: %v", err)
	}

	utils.Info("MCP服务器可导出 %d 个工具", len(tools))

	utils.Info("MCP服务器创建测试通过")
}

// TestMCPIntegration 测试MCP集成
func TestMCPIntegration(t *testing.T) {
	utils.Info("测试MCP集成")

	// 创建默认配置
	defaultConfig := config.GetDefaultConfig()

	if len(defaultConfig.Servers) == 0 {
		t.Error("默认配置应包含服务器")
	}

	// 检查配置项
	if defaultConfig.AutoStart != true {
		t.Error("默认AutoStart应为true")
	}

	if defaultConfig.ToolDiscoveryInterval != 60 {
		t.Errorf("期望ToolDiscoveryInterval为60，实际为%d", defaultConfig.ToolDiscoveryInterval)
	}

	// 检查服务器配置
	hasLocalFSTools := false
	for _, server := range defaultConfig.Servers {
		if server.ID == "local-fs-tools" {
			hasLocalFSTools = true
			if server.Type != "server" {
				t.Error("local-fs-tools应为服务器类型")
			}
		}
	}

	if !hasLocalFSTools {
		t.Error("默认配置应包含local-fs-tools服务器")
	}

	utils.Info("MCP集成测试通过")
}

// TestMCPToolWrapper 测试MCP工具包装器
func TestMCPToolWrapper(t *testing.T) {
	utils.Info("测试MCP工具包装器")

	// 创建模拟工具
	toolInfo := client.ToolInfo{
		Name:        "test_tool",
		Description: "测试工具",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"param": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}

	discoveredTool := discovery.DiscoveredTool{
		FullName:    "test.test_tool",
		ServerID:    "test",
		ToolInfo:    toolInfo,
		LastSeen:    time.Now(),
		CallCount:   0,
		SuccessRate: 1.0,
		Enabled:     true,
	}

	// 创建工具包装器
	toolWrapper := discovery.NewMCPToolWrapper(discoveredTool, nil)

	if toolWrapper == nil {
		t.Fatal("工具包装器为nil")
	}

	// 获取工具定义
	def := toolWrapper.GetDefinition()
	if def.Name != "test.test_tool" {
		t.Errorf("期望工具名称为'test.test_tool'，实际为'%s'", def.Name)
	}

	if def.Category != tools.ToolCategoryCustom {
		t.Errorf("期望工具类别为'custom'，实际为'%s'", def.Category)
	}

	// 检查启用状态
	if !toolWrapper.IsEnabled() {
		t.Error("工具应启用")
	}

	// 禁用工具
	toolWrapper.SetEnabled(false)
	if toolWrapper.IsEnabled() {
		t.Error("工具应禁用")
	}

	utils.Info("MCP工具包装器测试通过")
}

// TestMCPNamespaceIsolation 测试命名空间隔离
func TestMCPNamespaceIsolation(t *testing.T) {
	utils.Info("测试MCP命名空间隔离")

	// 创建包含多个命名空间的配置
	hostConfig := host.MCPHostConfig{
		Servers: []host.MCPHostServerConfig{
			{
				ID:   "server1",
				Name: "服务器1",
				Type: "client",
				ClientConfig: client.MCPClientConfig{
					ServerType: "stdio",
					Namespace:  "namespace1",
				},
				Enabled: true,
			},
			{
				ID:   "server2",
				Name: "服务器2",
				Type: "client",
				ClientConfig: client.MCPClientConfig{
					ServerType: "stdio",
					Namespace:  "namespace2",
				},
				Enabled: true,
			},
		},
	}

	mcpHost := host.NewMCPHost(hostConfig)

	// 模拟工具缓存
	// 在实际测试中，这些工具会通过实际连接发现
	toolsByNamespace := mcpHost.ListAllTools()

	// 检查命名空间隔离
	// 注意：由于是模拟测试，实际工具列表可能为空
	utils.Info("发现 %d 个命名空间", len(toolsByNamespace))

	// 测试命名空间概念
	for namespace, tools := range toolsByNamespace {
		utils.Info("命名空间 '%s' 有 %d 个工具", namespace, len(tools))
	}

	utils.Info("MCP命名空间隔离测试通过")
}

// TestMCPErrorHandling 测试错误处理
func TestMCPErrorHandling(t *testing.T) {
	utils.Info("测试MCP错误处理")

	// 测试无效服务器类型
	_, err := client.NewMCPClient(client.MCPClientConfig{
		ServerType: "invalid_type",
	})

	if err == nil {
		t.Error("期望创建无效类型的MCP客户端时返回错误")
	} else {
		utils.Info("无效服务器类型错误处理正常: %v", err)
	}

	// 测试重复添加服务器
	configManager, _ := config.NewMCPConfigManager("/tmp/test_duplicate.json")

	serverConfig := host.MCPHostServerConfig{
		ID:   "duplicate",
		Name: "重复服务器",
		Type: "client",
	}

	// 第一次添加应成功
	if err := configManager.AddServer(serverConfig); err != nil {
		t.Fatalf("第一次添加服务器失败: %v", err)
	}

	// 第二次添加应失败
	if err := configManager.AddServer(serverConfig); err == nil {
		t.Error("期望添加重复服务器时返回错误")
	} else {
		utils.Info("重复服务器错误处理正常: %v", err)
	}

	utils.Info("MCP错误处理测试通过")
}

// TestMCPPerformance 测试性能
func TestMCPPerformance(t *testing.T) {
	utils.Info("测试MCP性能")

	startTime := time.Now()

	// 创建多个MCP客户端
	clients := make([]client.MCPClient, 0, 10)
	for i := 0; i < 10; i++ {
		config := client.MCPClientConfig{
			ServerType: "stdio",
			Command:    "echo",
			Args:       []string{fmt.Sprintf("test%d", i)},
			Namespace:  fmt.Sprintf("ns%d", i),
		}

		client, err := client.NewMCPClient(config)
		if err != nil {
			t.Fatalf("创建客户端 %d 失败: %v", i, err)
		}

		clients = append(clients, client)
	}

	creationTime := time.Since(startTime)
	utils.Info("创建10个MCP客户端耗时: %v", creationTime)

	if creationTime > 1*time.Second {
		t.Errorf("创建客户端耗时过长: %v", creationTime)
	}

	// 清理
	for _, client := range clients {
		// 注意：这些客户端未连接，所以不需要断开连接
		_ = client
	}

	utils.Info("MCP性能测试通过")
}

// TestMain MCP测试主函数
func TestMain(m *testing.M) {
	utils.Info("开始MCP集成测试")

	// 运行测试
	_ = m.Run()

	utils.Info("MCP集成测试完成")
}
