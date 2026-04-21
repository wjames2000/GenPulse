package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/utils"
)

// MCPServer MCP服务器接口
type MCPServer interface {
	// Start 启动MCP服务器
	Start(ctx context.Context) error

	// Stop 停止MCP服务器
	Stop() error

	// IsRunning 检查是否正在运行
	IsRunning() bool

	// GetToolRegistry 获取工具注册表
	GetToolRegistry() *tools.ToolRegistry

	// ExportTools 导出工具到MCP协议
	ExportTools() ([]ToolDefinition, error)
}

// ToolDefinition MCP工具定义
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"inputSchema,omitempty"`
}

// StdioMCPServer 基于stdio的MCP服务器
type StdioMCPServer struct {
	toolRegistry *tools.ToolRegistry
	running      bool
	mu           sync.RWMutex
	stdin        io.Reader
	stdout       io.Writer
	stderr       io.Writer
	cancel       context.CancelFunc
}

// NewStdioMCPServer 创建新的Stdio MCP服务器
func NewStdioMCPServer(toolRegistry *tools.ToolRegistry) *StdioMCPServer {
	return &StdioMCPServer{
		toolRegistry: toolRegistry,
		running:      false,
		stdin:        os.Stdin,
		stdout:       os.Stdout,
		stderr:       os.Stderr,
	}
}

// Start 启动MCP服务器
func (s *StdioMCPServer) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("server already running")
	}

	utils.Info("启动MCP服务器 (stdio)")

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	// 启动消息处理循环
	go s.handleMessages(ctx)

	s.running = true
	utils.Info("MCP服务器已启动，等待客户端连接")

	return nil
}

// handleMessages 处理来自stdin的消息
func (s *StdioMCPServer) handleMessages(ctx context.Context) {
	decoder := json.NewDecoder(s.stdin)

	for {
		select {
		case <-ctx.Done():
			utils.Info("MCP服务器消息处理循环停止")
			return
		default:
			var msg map[string]interface{}
			if err := decoder.Decode(&msg); err != nil {
				if err != io.EOF {
					utils.Warn("读取MCP消息失败: %v", err)
				}
				return
			}

			// 处理消息
			go s.handleMessage(msg)
		}
	}
}

// handleMessage 处理单个消息
func (s *StdioMCPServer) handleMessage(msg map[string]interface{}) {
	msgType, ok := msg["type"].(string)
	if !ok {
		s.sendError("invalid message: missing type field")
		return
	}

	switch msgType {
	case "initialize":
		s.handleInitialize(msg)
	case "tools/list":
		s.handleListTools(msg)
	case "tools/call":
		s.handleCallTool(msg)
	default:
		s.sendError(fmt.Sprintf("unknown message type: %s", msgType))
	}
}

// handleInitialize 处理初始化请求
func (s *StdioMCPServer) handleInitialize(msg map[string]interface{}) {
	response := map[string]interface{}{
		"type": "initialize_result",
		"result": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{
					"listChanged": false,
				},
			},
			"serverInfo": map[string]interface{}{
				"name":    "GenPulse MCP Server",
				"version": "1.0.0",
			},
		},
	}

	s.sendResponse(response)
	utils.Info("MCP客户端初始化完成")
}

// handleListTools 处理列出工具请求
func (s *StdioMCPServer) handleListTools(msg map[string]interface{}) {
	toolDefs, err := s.ExportTools()
	if err != nil {
		s.sendError(fmt.Sprintf("failed to export tools: %v", err))
		return
	}

	// 转换为MCP协议格式
	tools := make([]map[string]interface{}, 0, len(toolDefs))
	for _, tool := range toolDefs {
		tools = append(tools, map[string]interface{}{
			"name":        tool.Name,
			"description": tool.Description,
			"inputSchema": tool.InputSchema,
		})
	}

	response := map[string]interface{}{
		"type": "tools/list_result",
		"result": map[string]interface{}{
			"tools": tools,
		},
	}

	s.sendResponse(response)
	utils.Debug("返回工具列表，共 %d 个工具", len(tools))
}

// handleCallTool 处理调用工具请求
func (s *StdioMCPServer) handleCallTool(msg map[string]interface{}) {
	params, ok := msg["params"].(map[string]interface{})
	if !ok {
		s.sendError("invalid call tool message: missing params")
		return
	}

	name, ok := params["name"].(string)
	if !ok {
		s.sendError("invalid call tool message: missing tool name")
		return
	}

	arguments, _ := params["arguments"].(map[string]interface{})

	utils.Info("MCP工具调用: %s", name)

	// 调用工具
	result, err := s.callTool(name, arguments)
	if err != nil {
		s.sendError(fmt.Sprintf("tool call failed: %v", err))
		return
	}

	response := map[string]interface{}{
		"type": "tools/call_result",
		"result": map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("%v", result),
				},
			},
		},
	}

	s.sendResponse(response)
	utils.Debug("工具调用完成: %s", name)
}

// callTool 调用本地工具
func (s *StdioMCPServer) callTool(name string, arguments map[string]interface{}) (interface{}, error) {
	if s.toolRegistry == nil {
		return nil, fmt.Errorf("tool registry not available")
	}

	// 从工具注册表获取工具
	tool, err := s.toolRegistry.GetTool(name)
	if err != nil {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	// 调用工具
	execution := tools.ToolExecution{
		ToolID:     name,
		Parameters: arguments,
	}

	result, err := tool.Execute(context.Background(), execution)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	return result, nil
}

// sendResponse 发送响应
func (s *StdioMCPServer) sendResponse(response map[string]interface{}) {
	data, err := json.Marshal(response)
	if err != nil {
		utils.Warn("序列化响应失败: %v", err)
		return
	}

	data = append(data, '\n')

	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, err := s.stdout.Write(data); err != nil {
		utils.Warn("发送响应失败: %v", err)
	}
}

// sendError 发送错误响应
func (s *StdioMCPServer) sendError(message string) {
	response := map[string]interface{}{
		"type": "error",
		"error": map[string]interface{}{
			"message": message,
		},
	}

	s.sendResponse(response)
	utils.Warn("发送错误响应: %s", message)
}

// Stop 停止MCP服务器
func (s *StdioMCPServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	utils.Info("停止MCP服务器")

	// 取消上下文
	if s.cancel != nil {
		s.cancel()
	}

	s.running = false
	utils.Info("MCP服务器已停止")

	return nil
}

// IsRunning 检查是否正在运行
func (s *StdioMCPServer) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetToolRegistry 获取工具注册表
func (s *StdioMCPServer) GetToolRegistry() *tools.ToolRegistry {
	return s.toolRegistry
}

// ExportTools 导出工具到MCP协议
func (s *StdioMCPServer) ExportTools() ([]ToolDefinition, error) {
	if s.toolRegistry == nil {
		return nil, fmt.Errorf("tool registry not available")
	}

	// 获取所有工具
	toolList := s.toolRegistry.ListTools()
	toolDefs := make([]ToolDefinition, 0, len(toolList))

	for _, def := range toolList {
		// 创建工具定义
		toolDef := ToolDefinition{
			Name:        def.Name,
			Description: def.Description,
			InputSchema: def.Parameters,
		}

		toolDefs = append(toolDefs, toolDef)
	}

	return toolDefs, nil
}

// NewMCPServer 创建MCP服务器
func NewMCPServer(toolRegistry *tools.ToolRegistry, serverType string) (MCPServer, error) {
	switch serverType {
	case "stdio":
		return NewStdioMCPServer(toolRegistry), nil
	default:
		return nil, fmt.Errorf("unsupported server type: %s", serverType)
	}
}
