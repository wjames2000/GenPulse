package client

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"GenPulse/internal/utils"
)

// MCPClientConfig MCP客户端配置
type MCPClientConfig struct {
	// ServerType 服务器类型: "stdio" 或 "sse"
	ServerType string `json:"server_type"`

	// Command 当ServerType为"stdio"时，要执行的命令
	Command string `json:"command,omitempty"`

	// Args 命令参数
	Args []string `json:"args,omitempty"`

	// URL 当ServerType为"sse"时，服务器URL
	URL string `json:"url,omitempty"`

	// Namespace 命名空间，用于工具隔离
	Namespace string `json:"namespace,omitempty"`

	// Timeout 连接超时时间（秒）
	Timeout int `json:"timeout,omitempty"`

	// AutoReconnect 是否自动重连
	AutoReconnect bool `json:"auto_reconnect,omitempty"`

	// MaxRetries 最大重试次数
	MaxRetries int `json:"max_retries,omitempty"`
}

// MCPClient MCP客户端接口
type MCPClient interface {
	// Connect 连接到MCP服务器
	Connect(ctx context.Context) error

	// Disconnect 断开连接
	Disconnect() error

	// IsConnected 检查是否已连接
	IsConnected() bool

	// ListTools 列出服务器提供的所有工具
	ListTools() ([]ToolInfo, error)

	// CallTool 调用工具
	CallTool(toolName string, arguments map[string]interface{}) (interface{}, error)

	// GetConfig 获取客户端配置
	GetConfig() MCPClientConfig

	// GetNamespace 获取命名空间
	GetNamespace() string
}

// ToolInfo 工具信息
type ToolInfo struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	InputSchema  map[string]interface{} `json:"inputSchema,omitempty"`
	OutputSchema map[string]interface{} `json:"outputSchema,omitempty"`
}

// StdioMCPClient 基于stdio的MCP客户端
type StdioMCPClient struct {
	config    MCPClientConfig
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	stderr    io.ReadCloser
	connected bool
	tools     []ToolInfo
}

// NewStdioMCPClient 创建新的Stdio MCP客户端
func NewStdioMCPClient(config MCPClientConfig) *StdioMCPClient {
	return &StdioMCPClient{
		config:    config,
		connected: false,
	}
}

// Connect 连接到stdio MCP服务器
func (c *StdioMCPClient) Connect(ctx context.Context) error {
	if c.connected {
		return fmt.Errorf("already connected")
	}

	utils.Info("连接到MCP服务器 (stdio): %s %v", c.config.Command, c.config.Args)

	// 创建命令
	cmd := exec.CommandContext(ctx, c.config.Command, c.config.Args...)

	// 获取标准输入输出管道
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	c.cmd = cmd
	c.stdin = stdin
	c.stdout = stdout
	c.stderr = stderr
	c.connected = true

	// 启动goroutine读取stderr
	go c.readStderr()

	utils.Info("MCP服务器连接成功")

	// 初始化连接并获取工具列表
	if err := c.initializeConnection(ctx); err != nil {
		c.Disconnect()
		return fmt.Errorf("failed to initialize connection: %w", err)
	}

	return nil
}

// initializeConnection 初始化连接并获取工具列表
func (c *StdioMCPClient) initializeConnection(ctx context.Context) error {
	// TODO: 实现MCP协议初始化握手
	// 这里应该发送初始化请求并接收工具列表

	// 临时实现：返回空工具列表
	c.tools = []ToolInfo{}

	utils.Info("MCP连接初始化完成，发现 %d 个工具", len(c.tools))
	return nil
}

// readStderr 读取stderr输出
func (c *StdioMCPClient) readStderr() {
	buf := make([]byte, 1024)
	for {
		n, err := c.stderr.Read(buf)
		if err != nil {
			if err != io.EOF {
				utils.Warn("读取MCP服务器stderr失败: %v", err)
			}
			break
		}
		if n > 0 {
			utils.Debug("MCP服务器stderr: %s", string(buf[:n]))
		}
	}
}

// Disconnect 断开连接
func (c *StdioMCPClient) Disconnect() error {
	if !c.connected {
		return nil
	}

	utils.Info("断开MCP服务器连接")

	// 关闭标准输入
	if c.stdin != nil {
		c.stdin.Close()
	}

	// 等待命令结束
	if c.cmd != nil && c.cmd.Process != nil {
		// 发送中断信号
		c.cmd.Process.Signal(os.Interrupt)

		// 等待进程结束
		done := make(chan error, 1)
		go func() {
			done <- c.cmd.Wait()
		}()

		// 设置超时
		select {
		case <-time.After(5 * time.Second):
			// 超时后强制终止
			c.cmd.Process.Kill()
			utils.Warn("MCP服务器进程强制终止")
		case err := <-done:
			if err != nil {
				utils.Debug("MCP服务器进程退出: %v", err)
			}
		}
	}

	c.connected = false
	c.cmd = nil
	c.stdin = nil
	c.stdout = nil
	c.stderr = nil
	c.tools = nil

	utils.Info("MCP服务器连接已断开")
	return nil
}

// IsConnected 检查是否已连接
func (c *StdioMCPClient) IsConnected() bool {
	return c.connected
}

// ListTools 列出所有工具
func (c *StdioMCPClient) ListTools() ([]ToolInfo, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected")
	}

	return c.tools, nil
}

// CallTool 调用工具
func (c *StdioMCPClient) CallTool(toolName string, arguments map[string]interface{}) (interface{}, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected")
	}

	// TODO: 实现MCP协议工具调用
	// 这里应该发送工具调用请求并接收响应

	utils.Info("调用MCP工具: %s", toolName)

	// 临时实现：返回模拟结果
	return map[string]interface{}{
		"success": true,
		"tool":    toolName,
		"result":  "MCP工具调用成功（模拟）",
	}, nil
}

// GetConfig 获取配置
func (c *StdioMCPClient) GetConfig() MCPClientConfig {
	return c.config
}

// GetNamespace 获取命名空间
func (c *StdioMCPClient) GetNamespace() string {
	return c.config.Namespace
}

// SSEMCPClient 基于SSE的MCP客户端
type SSEMCPClient struct {
	config    MCPClientConfig
	connected bool
	tools     []ToolInfo
}

// NewSSEMCPClient 创建新的SSE MCP客户端
func NewSSEMCPClient(config MCPClientConfig) *SSEMCPClient {
	return &SSEMCPClient{
		config:    config,
		connected: false,
	}
}

// Connect 连接到SSE MCP服务器
func (c *SSEMCPClient) Connect(ctx context.Context) error {
	if c.connected {
		return fmt.Errorf("already connected")
	}

	utils.Info("连接到MCP服务器 (SSE): %s", c.config.URL)

	// TODO: 实现SSE连接
	// 这里应该建立SSE连接并处理事件

	c.connected = true
	utils.Info("MCP服务器连接成功")

	// 初始化连接
	if err := c.initializeConnection(ctx); err != nil {
		c.Disconnect()
		return fmt.Errorf("failed to initialize connection: %w", err)
	}

	return nil
}

// initializeConnection 初始化连接
func (c *SSEMCPClient) initializeConnection(ctx context.Context) error {
	// TODO: 实现SSE连接初始化
	c.tools = []ToolInfo{}
	utils.Info("MCP连接初始化完成，发现 %d 个工具", len(c.tools))
	return nil
}

// Disconnect 断开连接
func (c *SSEMCPClient) Disconnect() error {
	if !c.connected {
		return nil
	}

	utils.Info("断开MCP服务器连接")

	// TODO: 关闭SSE连接

	c.connected = false
	c.tools = nil

	utils.Info("MCP服务器连接已断开")
	return nil
}

// IsConnected 检查是否已连接
func (c *SSEMCPClient) IsConnected() bool {
	return c.connected
}

// ListTools 列出所有工具
func (c *SSEMCPClient) ListTools() ([]ToolInfo, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected")
	}

	return c.tools, nil
}

// CallTool 调用工具
func (c *SSEMCPClient) CallTool(toolName string, arguments map[string]interface{}) (interface{}, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected")
	}

	// TODO: 实现SSE工具调用

	utils.Info("调用MCP工具: %s", toolName)

	return map[string]interface{}{
		"success": true,
		"tool":    toolName,
		"result":  "MCP工具调用成功（模拟）",
	}, nil
}

// GetConfig 获取配置
func (c *SSEMCPClient) GetConfig() MCPClientConfig {
	return c.config
}

// GetNamespace 获取命名空间
func (c *SSEMCPClient) GetNamespace() string {
	return c.config.Namespace
}

// NewMCPClient 根据配置创建MCP客户端
func NewMCPClient(config MCPClientConfig) (MCPClient, error) {
	switch config.ServerType {
	case "stdio":
		return NewStdioMCPClient(config), nil
	case "sse":
		return NewSSEMCPClient(config), nil
	default:
		return nil, fmt.Errorf("unsupported server type: %s", config.ServerType)
	}
}
