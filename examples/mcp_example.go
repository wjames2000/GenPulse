package main

import (
	"fmt"
	"log"

	"GenPulse/internal/mcp/client"
	"GenPulse/internal/mcp/config"
	"GenPulse/internal/mcp/host"
)

func main() {
	fmt.Println("=== GenPulse MCP 集成示例 ===")

	// 示例1: 创建MCP客户端
	exampleMCPClient()

	// 示例2: 管理MCP配置
	exampleMCPConfig()

	// 示例3: 使用MCP主机
	exampleMCPHost()

	fmt.Println("\n=== MCP 集成示例完成 ===")
}

// exampleMCPClient 示例：创建和使用MCP客户端
func exampleMCPClient() {
	fmt.Println("\n--- 示例1: MCP客户端 ---")

	// 创建MCP客户端配置
	clientConfig := client.MCPClientConfig{
		ServerType:    "stdio",
		Command:       "echo",
		Args:          []string{"Hello from MCP Server"},
		Namespace:     "example",
		Timeout:       10,
		AutoReconnect: true,
		MaxRetries:    3,
	}

	// 创建MCP客户端
	fmt.Printf("创建MCP客户端: %s\n", clientConfig.Command)
	mcpClient, err := client.NewMCPClient(clientConfig)
	if err != nil {
		log.Printf("创建MCP客户端失败: %v\n", err)
		return
	}

	fmt.Printf("MCP客户端创建成功，命名空间: %s\n", mcpClient.GetNamespace())

	// 注意：这个示例客户端不会实际连接，因为echo命令不是真正的MCP服务器
	// 在实际使用中，您需要连接到一个真正的MCP服务器

	config := mcpClient.GetConfig()
	fmt.Printf("客户端配置:\n")
	fmt.Printf("  - 服务器类型: %s\n", config.ServerType)
	fmt.Printf("  - 命令: %s\n", config.Command)
	fmt.Printf("  - 命名空间: %s\n", config.Namespace)
	fmt.Printf("  - 超时: %d秒\n", config.Timeout)
}

// exampleMCPConfig 示例：管理MCP配置
func exampleMCPConfig() {
	fmt.Println("\n--- 示例2: MCP配置管理 ---")

	// 创建临时配置文件路径
	configPath := "/tmp/genpulse_mcp_example.json"

	// 创建配置管理器
	fmt.Printf("创建MCP配置管理器: %s\n", configPath)
	configManager, err := config.NewMCPConfigManager(configPath)
	if err != nil {
		log.Printf("创建配置管理器失败: %v\n", err)
		return
	}

	// 获取默认配置
	defaultConfig := config.GetDefaultConfig()
	fmt.Printf("默认配置包含 %d 个服务器\n", len(defaultConfig.Servers))

	// 列出配置中的服务器
	servers := configManager.ListServers()
	fmt.Printf("当前配置有 %d 个服务器\n", len(servers))

	// 添加一个新的服务器配置
	newServer := host.MCPHostServerConfig{
		ID:       "example-server",
		Name:     "示例服务器",
		Type:     "client",
		Enabled:  true,
		Priority: 50,
		ClientConfig: client.MCPClientConfig{
			ServerType: "stdio",
			Command:    "python3",
			Args:       []string{"-c", "print('MCP Server Ready')"},
			Namespace:  "python",
			Timeout:    30,
		},
	}

	fmt.Printf("添加服务器: %s\n", newServer.Name)
	if err := configManager.AddServer(newServer); err != nil {
		log.Printf("添加服务器失败: %v\n", err)
	} else {
		fmt.Println("服务器添加成功")
	}

	// 验证服务器已添加
	servers = configManager.ListServers()
	fmt.Printf("现在配置有 %d 个服务器\n", len(servers))

	for _, server := range servers {
		fmt.Printf("  - %s (%s): %s\n", server.Name, server.ID, server.Type)
	}

	// 移除服务器
	fmt.Printf("移除服务器: %s\n", newServer.ID)
	if err := configManager.RemoveServer(newServer.ID); err != nil {
		log.Printf("移除服务器失败: %v\n", err)
	} else {
		fmt.Println("服务器移除成功")
	}
}

// exampleMCPHost 示例：使用MCP主机
func exampleMCPHost() {
	fmt.Println("\n--- 示例3: MCP主机管理 ---")

	// 创建MCP主机配置
	hostConfig := host.MCPHostConfig{
		AutoStart:             false, // 不自动启动，仅演示配置
		ToolDiscoveryInterval: 30,
		MaxConcurrentCalls:    5,
		Servers: []host.MCPHostServerConfig{
			{
				ID:       "demo-client-1",
				Name:     "演示客户端1",
				Type:     "client",
				Enabled:  true,
				Priority: 100,
				ClientConfig: client.MCPClientConfig{
					ServerType: "stdio",
					Command:    "echo",
					Args:       []string{"Demo MCP Client 1"},
					Namespace:  "demo1",
				},
			},
			{
				ID:       "demo-client-2",
				Name:     "演示客户端2",
				Type:     "client",
				Enabled:  true,
				Priority: 90,
				ClientConfig: client.MCPClientConfig{
					ServerType: "stdio",
					Command:    "echo",
					Args:       []string{"Demo MCP Client 2"},
					Namespace:  "demo2",
				},
			},
		},
	}

	// 创建MCP主机
	fmt.Println("创建MCP主机...")
	mcpHost := host.NewMCPHost(hostConfig)

	// 获取配置
	config := mcpHost.GetConfig()
	fmt.Printf("MCP主机配置:\n")
	fmt.Printf("  - 自动启动: %v\n", config.AutoStart)
	fmt.Printf("  - 工具发现间隔: %d秒\n", config.ToolDiscoveryInterval)
	fmt.Printf("  - 最大并发调用: %d\n", config.MaxConcurrentCalls)
	fmt.Printf("  - 服务器数量: %d\n", len(config.Servers))

	// 演示命名空间隔离概念
	fmt.Println("\n命名空间隔离演示:")
	fmt.Println("每个MCP客户端可以分配一个唯一的命名空间")
	fmt.Println("这确保了工具名称不会冲突，即使不同服务器提供同名工具")

	for _, server := range config.Servers {
		fmt.Printf("  - %s: 命名空间 = %s\n", server.Name, server.ClientConfig.Namespace)
	}

	// 演示MCP主机的状态管理
	fmt.Println("\nMCP主机状态管理:")
	fmt.Println("MCP主机可以:")
	fmt.Println("  1. 管理多个MCP服务器连接")
	fmt.Println("  2. 自动发现和注册工具")
	fmt.Println("  3. 提供命名空间隔离")
	fmt.Println("  4. 处理工具调用路由")
	fmt.Println("  5. 监控连接状态")

	// 注意：在实际使用中，您需要启动MCP主机
	// ctx := context.Background()
	// if err := mcpHost.Start(ctx); err != nil {
	//     log.Printf("启动MCP主机失败: %v\n", err)
	// }

	fmt.Println("\nMCP主机创建完成（未实际启动）")
}

// 运行MCP示例
func init() {
	// 设置简单的日志输出
	log.SetFlags(0)
}

// 实际运行示例
func runExample() {
	// 在实际应用中，您会这样使用MCP:
	fmt.Println("\n=== 实际MCP使用模式 ===")

	// 1. 加载配置
	configManager, _ := config.NewMCPConfigManager("./data/mcp_config.json")

	// 2. 创建MCP主机
	_ = host.NewMCPHost(configManager.GetConfig())

	// 3. 启动MCP主机（示例中不实际启动）
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := mcpHost.Start(ctx); err != nil {
	//     log.Fatalf("启动MCP主机失败: %v", err)
	// }

	// 4. 使用发现的工具
	// tools := mcpHost.ListAllTools()
	// fmt.Printf("发现 %d 个命名空间的工具\n", len(tools))

	// 5. 调用工具
	// result, err := mcpHost.CallTool("server-id", "tool-name", map[string]interface{}{"param": "value"})

	fmt.Println("MCP集成已准备就绪")
}
