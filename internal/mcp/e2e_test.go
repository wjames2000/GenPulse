package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"GenPulse/internal/mcp/client"
	"GenPulse/internal/mcp/config"
	"GenPulse/internal/mcp/host"
	"GenPulse/internal/utils"
)

// TestMCPEndToEnd 测试MCP端到端流程
func TestMCPEndToEnd(t *testing.T) {
	utils.Info("开始MCP端到端测试")

	// 创建临时目录
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "mcp_config.json")

	// 1. 测试配置管理器
	t.Run("配置管理器", func(t *testing.T) {
		testConfigManager(t, configPath)
	})

	// 2. 测试MCP主机
	t.Run("MCP主机", func(t *testing.T) {
		testMCPHost(t, configPath)
	})

	// 3. 测试配置同步
	t.Run("配置同步", func(t *testing.T) {
		testConfigSync(t, configPath)
	})

	// 4. 测试错误恢复
	t.Run("错误恢复", func(t *testing.T) {
		testErrorRecovery(t, configPath)
	})

	utils.Info("MCP端到端测试完成")
}

// testConfigManager 测试配置管理器
func testConfigManager(t *testing.T, configPath string) {
	utils.Info("测试配置管理器")

	// 创建配置管理器
	configManager, err := config.NewMCPConfigManager(configPath)
	if err != nil {
		t.Fatalf("创建配置管理器失败: %v", err)
	}

	// 测试初始配置
	initialConfig := configManager.GetConfig()
	if !initialConfig.AutoStart {
		t.Error("期望AutoStart为true")
	}

	// 测试添加服务器
	serverConfig := host.MCPHostServerConfig{
		ID:       "e2e-test-server",
		Name:     "端到端测试服务器",
		Type:     "client",
		Enabled:  true,
		Priority: 50,
		ClientConfig: client.MCPClientConfig{
			ServerType: "stdio",
			Command:    "echo",
			Args:       []string{"e2e test"},
			Namespace:  "e2e",
			Timeout:    10,
		},
	}

	if err := configManager.AddServer(serverConfig); err != nil {
		t.Fatalf("添加服务器失败: %v", err)
	}

	// 验证服务器已添加
	servers := configManager.ListServers()
	if len(servers) != 1 {
		t.Errorf("期望有1个服务器，实际有%d个", len(servers))
	}

	// 测试获取服务器
	retrievedServer, err := configManager.GetServer("e2e-test-server")
	if err != nil {
		t.Fatalf("获取服务器失败: %v", err)
	}

	if retrievedServer.Name != "端到端测试服务器" {
		t.Errorf("期望服务器名称为'端到端测试服务器'，实际为'%s'", retrievedServer.Name)
	}

	// 测试更新服务器
	updatedConfig := *retrievedServer
	updatedConfig.Name = "更新后的测试服务器"
	updatedConfig.Priority = 75

	if err := configManager.UpdateServer("e2e-test-server", updatedConfig); err != nil {
		t.Fatalf("更新服务器失败: %v", err)
	}

	// 验证更新
	updatedServer, err := configManager.GetServer("e2e-test-server")
	if err != nil {
		t.Fatalf("获取更新后的服务器失败: %v", err)
	}

	if updatedServer.Name != "更新后的测试服务器" {
		t.Errorf("期望更新后的名称为'更新后的测试服务器'，实际为'%s'", updatedServer.Name)
	}

	if updatedServer.Priority != 75 {
		t.Errorf("期望更新后的优先级为75，实际为%d", updatedServer.Priority)
	}

	// 测试移除服务器
	if err := configManager.RemoveServer("e2e-test-server"); err != nil {
		t.Fatalf("移除服务器失败: %v", err)
	}

	// 验证移除
	servers = configManager.ListServers()
	if len(servers) != 0 {
		t.Errorf("期望有0个服务器，实际有%d个", len(servers))
	}

	// 测试配置文件存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("配置文件不存在: %s", configPath)
	}

	utils.Info("配置管理器测试通过")
}

// testMCPHost 测试MCP主机
func testMCPHost(t *testing.T, configPath string) {
	utils.Info("测试MCP主机")

	// 创建配置管理器并添加服务器
	configManager, err := config.NewMCPConfigManager(configPath)
	if err != nil {
		t.Fatalf("创建配置管理器失败: %v", err)
	}

	// 添加测试服务器
	serverConfig := host.MCPHostServerConfig{
		ID:       "e2e-host-test",
		Name:     "主机测试服务器",
		Type:     "client",
		Enabled:  true,
		Priority: 60,
		ClientConfig: client.MCPClientConfig{
			ServerType: "stdio",
			Command:    "echo",
			Args:       []string{"host test"},
			Namespace:  "host",
			Timeout:    5,
		},
	}

	if err := configManager.AddServer(serverConfig); err != nil {
		t.Fatalf("添加服务器失败: %v", err)
	}

	// 获取配置并创建MCP主机
	hostConfig := configManager.GetConfig()
	mcpHost := host.NewMCPHost(hostConfig)

	// 测试启动
	ctx := context.Background()
	if err := mcpHost.Start(ctx); err != nil {
		t.Fatalf("启动MCP主机失败: %v", err)
	}

	// 验证运行状态
	if !mcpHost.IsRunning() {
		t.Error("MCP主机未运行")
	}

	// 测试获取服务器状态
	status, err := mcpHost.GetServerStatus("e2e-host-test")
	if err != nil {
		t.Fatalf("获取服务器状态失败: %v", err)
	}

	if status["type"] != "client" {
		t.Errorf("期望服务器类型为'client'，实际为'%v'", status["type"])
	}

	// 测试列出所有工具
	tools := mcpHost.ListAllTools()
	utils.Info("发现 %d 个命名空间的工具", len(tools))

	// 测试停止
	if err := mcpHost.Stop(); err != nil {
		t.Fatalf("停止MCP主机失败: %v", err)
	}

	if mcpHost.IsRunning() {
		t.Error("MCP主机仍在运行")
	}

	// 测试移除服务器
	if err := mcpHost.RemoveServer("e2e-host-test"); err != nil {
		t.Fatalf("移除服务器失败: %v", err)
	}

	utils.Info("MCP主机测试通过")
}

// testConfigSync 测试配置同步
func testConfigSync(t *testing.T, configPath string) {
	utils.Info("测试配置同步")

	// 创建配置管理器
	configManager, err := config.NewMCPConfigManager(configPath)
	if err != nil {
		t.Fatalf("创建配置管理器失败: %v", err)
	}

	// 创建初始配置
	initialConfig := host.MCPHostConfig{
		AutoStart:             false,
		ToolDiscoveryInterval: 30,
		MaxConcurrentCalls:    5,
		Servers: []host.MCPHostServerConfig{
			{
				ID:       "sync-test-1",
				Name:     "同步测试1",
				Type:     "client",
				Enabled:  true,
				Priority: 80,
				ClientConfig: client.MCPClientConfig{
					ServerType: "stdio",
					Command:    "echo",
					Args:       []string{"sync1"},
					Namespace:  "sync",
				},
			},
			{
				ID:       "sync-test-2",
				Name:     "同步测试2",
				Type:     "client",
				Enabled:  false,
				Priority: 40,
				ClientConfig: client.MCPClientConfig{
					ServerType: "stdio",
					Command:    "echo",
					Args:       []string{"sync2"},
					Namespace:  "sync",
				},
			},
		},
	}

	// 更新配置
	if err := configManager.UpdateConfig(initialConfig); err != nil {
		t.Fatalf("更新配置失败: %v", err)
	}

	// 验证配置已保存
	savedConfig := configManager.GetConfig()
	if len(savedConfig.Servers) != 2 {
		t.Errorf("期望有2个服务器，实际有%d个", len(savedConfig.Servers))
	}

	// 创建MCP主机
	mcpHost := host.NewMCPHost(savedConfig)

	// 测试配置更新
	updatedConfig := savedConfig
	updatedConfig.AutoStart = true
	updatedConfig.ToolDiscoveryInterval = 15

	// 使用回调更新配置
	updateErr := mcpHost.UpdateConfigWithCallback(updatedConfig, func(cfg host.MCPHostConfig) error {
		return configManager.UpdateConfig(cfg)
	})

	if updateErr != nil {
		t.Fatalf("更新配置失败: %v", updateErr)
	}

	// 验证MCP主机配置已更新
	hostConfig := mcpHost.GetConfig()
	if !hostConfig.AutoStart {
		t.Error("期望AutoStart为true")
	}

	if hostConfig.ToolDiscoveryInterval != 15 {
		t.Errorf("期望ToolDiscoveryInterval为15，实际为%d", hostConfig.ToolDiscoveryInterval)
	}

	// 验证配置文件已更新
	fileConfig := configManager.GetConfig()
	if !fileConfig.AutoStart {
		t.Error("配置文件中的AutoStart应为true")
	}

	utils.Info("配置同步测试通过")
}

// testErrorRecovery 测试错误恢复
func testErrorRecovery(t *testing.T, configPath string) {
	utils.Info("测试错误恢复")

	// 创建配置管理器
	configManager, err := config.NewMCPConfigManager(configPath)
	if err != nil {
		t.Fatalf("创建配置管理器失败: %v", err)
	}

	// 测试无效服务器配置
	invalidConfig := host.MCPHostServerConfig{
		ID:       "invalid-test",
		Name:     "无效测试服务器",
		Type:     "invalid_type", // 无效类型
		Enabled:  true,
		Priority: 50,
	}

	// 添加无效配置应失败
	if err := configManager.AddServer(invalidConfig); err == nil {
		t.Error("期望添加无效类型服务器时返回错误")
	}

	// 测试重复ID
	server1 := host.MCPHostServerConfig{
		ID:       "duplicate-id",
		Name:     "服务器1",
		Type:     "client",
		Enabled:  true,
		Priority: 50,
		ClientConfig: client.MCPClientConfig{
			ServerType: "stdio",
			Command:    "echo",
			Args:       []string{"test"},
			Namespace:  "test",
		},
	}

	server2 := host.MCPHostServerConfig{
		ID:       "duplicate-id", // 相同ID
		Name:     "服务器2",
		Type:     "client",
		Enabled:  true,
		Priority: 60,
		ClientConfig: client.MCPClientConfig{
			ServerType: "stdio",
			Command:    "echo",
			Args:       []string{"test2"},
			Namespace:  "test2",
		},
	}

	// 第一次添加应成功
	if err := configManager.AddServer(server1); err != nil {
		t.Fatalf("第一次添加服务器失败: %v", err)
	}

	// 第二次添加应失败
	if err := configManager.AddServer(server2); err == nil {
		t.Error("期望添加重复ID服务器时返回错误")
	}

	// 测试不存在的服务器
	if _, err := configManager.GetServer("non-existent"); err == nil {
		t.Error("期望获取不存在的服务器时返回错误")
	}

	if err := configManager.RemoveServer("non-existent"); err == nil {
		t.Error("期望移除不存在的服务器时返回错误")
	}

	if err := configManager.UpdateServer("non-existent", server1); err == nil {
		t.Error("期望更新不存在的服务器时返回错误")
	}

	// 测试MCP主机错误处理
	hostConfig := configManager.GetConfig()
	// 设置不自动启动工具发现，避免测试超时
	hostConfig.ToolDiscoveryInterval = 0
	mcpHost := host.NewMCPHost(hostConfig)

	// 测试启动已运行的主机
	ctx := context.Background()
	if err := mcpHost.Start(ctx); err != nil {
		t.Fatalf("启动MCP主机失败: %v", err)
	}

	// 再次启动应失败
	if err := mcpHost.Start(ctx); err == nil {
		t.Error("期望启动已运行的主机时返回错误")
	}

	// 测试停止未运行的主机
	if err := mcpHost.Stop(); err != nil {
		t.Fatalf("停止MCP主机失败: %v", err)
	}

	// 再次停止应成功（无操作）
	if err := mcpHost.Stop(); err != nil {
		t.Fatalf("再次停止MCP主机失败: %v", err)
	}

	// 测试获取不存在的服务器状态
	if _, err := mcpHost.GetServerStatus("non-existent"); err == nil {
		t.Error("期望获取不存在的服务器状态时返回错误")
	}

	// 测试移除不存在的服务器
	if err := mcpHost.RemoveServer("non-existent"); err == nil {
		t.Error("期望移除不存在的服务器时返回错误")
	}

	utils.Info("错误恢复测试通过")
}

// TestMCPConcurrentAccess 测试并发访问
func TestMCPConcurrentAccess(t *testing.T) {
	utils.Info("测试MCP并发访问")

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "concurrent_config.json")

	// 创建配置管理器
	configManager, err := config.NewMCPConfigManager(configPath)
	if err != nil {
		t.Fatalf("创建配置管理器失败: %v", err)
	}

	// 并发添加服务器
	const numGoroutines = 10
	errCh := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			serverConfig := host.MCPHostServerConfig{
				ID:       fmt.Sprintf("concurrent-%d", id),
				Name:     fmt.Sprintf("并发服务器%d", id),
				Type:     "client",
				Enabled:  true,
				Priority: id * 10,
				ClientConfig: client.MCPClientConfig{
					ServerType: "stdio",
					Command:    "echo",
					Args:       []string{fmt.Sprintf("concurrent%d", id)},
					Namespace:  fmt.Sprintf("concurrent%d", id),
				},
			}

			errCh <- configManager.AddServer(serverConfig)
		}(i)
	}

	// 收集错误
	for i := 0; i < numGoroutines; i++ {
		if err := <-errCh; err != nil {
			t.Errorf("并发添加服务器失败: %v", err)
		}
	}

	// 验证所有服务器都已添加
	servers := configManager.ListServers()
	if len(servers) != numGoroutines {
		t.Errorf("期望有%d个服务器，实际有%d个", numGoroutines, len(servers))
	}

	// 并发读取
	done := make(chan bool, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			serverID := fmt.Sprintf("concurrent-%d", id)
			server, err := configManager.GetServer(serverID)
			if err != nil {
				t.Errorf("获取服务器 %s 失败: %v", serverID, err)
			} else if server.Name != fmt.Sprintf("并发服务器%d", id) {
				t.Errorf("期望服务器名称为'并发服务器%d'，实际为'%s'", id, server.Name)
			}
			done <- true
		}(i)
	}

	// 等待所有读取完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// 并发更新
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			serverID := fmt.Sprintf("concurrent-%d", id)
			server, err := configManager.GetServer(serverID)
			if err != nil {
				t.Errorf("获取服务器 %s 失败: %v", serverID, err)
				return
			}

			updatedConfig := *server
			updatedConfig.Priority = id*10 + 5

			errCh <- configManager.UpdateServer(serverID, updatedConfig)
		}(i)
	}

	// 收集更新错误
	for i := 0; i < numGoroutines; i++ {
		if err := <-errCh; err != nil {
			t.Errorf("并发更新服务器失败: %v", err)
		}
	}

	// 验证更新
	for i := 0; i < numGoroutines; i++ {
		serverID := fmt.Sprintf("concurrent-%d", i)
		server, err := configManager.GetServer(serverID)
		if err != nil {
			t.Errorf("获取更新后的服务器 %s 失败: %v", serverID, err)
		} else if server.Priority != i*10+5 {
			t.Errorf("期望服务器 %s 的优先级为%d，实际为%d", serverID, i*10+5, server.Priority)
		}
	}

	utils.Info("MCP并发访问测试通过")
}

// TestMCPPerformanceBenchmark 测试性能基准
func TestMCPPerformanceBenchmark(t *testing.T) {
	utils.Info("测试MCP性能基准")

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "benchmark_config.json")

	// 创建配置管理器
	configManager, err := config.NewMCPConfigManager(configPath)
	if err != nil {
		t.Fatalf("创建配置管理器失败: %v", err)
	}

	// 基准测试：添加大量服务器
	const numServers = 100
	startTime := time.Now()

	for i := 0; i < numServers; i++ {
		serverConfig := host.MCPHostServerConfig{
			ID:       fmt.Sprintf("benchmark-%d", i),
			Name:     fmt.Sprintf("基准测试服务器%d", i),
			Type:     "client",
			Enabled:  i%2 == 0, // 交替启用/禁用
			Priority: i % 100,
			ClientConfig: client.MCPClientConfig{
				ServerType: "stdio",
				Command:    "echo",
				Args:       []string{fmt.Sprintf("benchmark%d", i)},
				Namespace:  fmt.Sprintf("benchmark%d", i),
			},
		}

		if err := configManager.AddServer(serverConfig); err != nil {
			t.Fatalf("添加服务器 %d 失败: %v", i, err)
		}
	}

	addTime := time.Since(startTime)
	utils.Info("添加 %d 个服务器耗时: %v", numServers, addTime)

	// 基准测试：列出服务器
	startTime = time.Now()
	servers := configManager.ListServers()
	listTime := time.Since(startTime)

	if len(servers) != numServers {
		t.Errorf("期望有%d个服务器，实际有%d个", numServers, len(servers))
	}
	utils.Info("列出 %d 个服务器耗时: %v", numServers, listTime)

	// 基准测试：获取启用的服务器
	startTime = time.Now()
	enabledServers := configManager.GetEnabledServers()
	getEnabledTime := time.Since(startTime)

	expectedEnabled := numServers / 2
	if len(enabledServers) != expectedEnabled {
		t.Errorf("期望有%d个启用的服务器，实际有%d个", expectedEnabled, len(enabledServers))
	}
	utils.Info("获取 %d 个启用的服务器耗时: %v", expectedEnabled, getEnabledTime)

	// 基准测试：随机访问
	const numAccesses = 1000
	startTime = time.Now()

	for i := 0; i < numAccesses; i++ {
		serverID := fmt.Sprintf("benchmark-%d", i%numServers)
		_, err := configManager.GetServer(serverID)
		if err != nil {
			t.Errorf("获取服务器 %s 失败: %v", serverID, err)
		}
	}

	accessTime := time.Since(startTime)
	utils.Info("随机访问 %d 次耗时: %v", numAccesses, accessTime)

	// 性能要求
	if addTime > 2*time.Second {
		t.Errorf("添加服务器耗时过长: %v", addTime)
	}

	if listTime > 100*time.Millisecond {
		t.Errorf("列出服务器耗时过长: %v", listTime)
	}

	if accessTime > 500*time.Millisecond {
		t.Errorf("随机访问耗时过长: %v", accessTime)
	}

	utils.Info("MCP性能基准测试通过")
}
