package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"GenPulse/internal/utils"
)

// ShellTool 命令行执行工具
type ShellTool struct {
	*BaseTool
	workspacePath   string
	allowedCommands []string // 命令白名单，空表示允许所有命令
	commandTimeout  time.Duration
	mutex           sync.RWMutex
}

// NewShellTool 创建命令行执行工具
func NewShellTool(workspacePath string) (*ShellTool, error) {
	// 确保工作区目录存在
	if err := os.MkdirAll(workspacePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// 默认允许的命令白名单
	allowedCommands := []string{
		"go", "npm", "yarn", "pnpm", "node", "python", "python3", "pip", "pip3",
		"git", "ls", "cat", "echo", "mkdir", "rm", "cp", "mv", "find", "grep",
		"curl", "wget", "tar", "unzip", "zip", "chmod", "chown",
	}

	definition := ToolDefinition{
		ID:          "shell_tool",
		Name:        "Shell Tool",
		Description: "提供命令行执行功能，支持命令白名单与超时控制",
		Category:    ToolCategoryShell,
		Version:     "1.0.0",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"command": map[string]interface{}{
					"type":        "string",
					"description": "要执行的命令",
				},
				"args": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "string",
					},
					"description": "命令参数",
				},
				"working_dir": map[string]interface{}{
					"type":        "string",
					"description": "工作目录（相对于工作区）",
					"default":     ".",
				},
				"timeout": map[string]interface{}{
					"type":        "integer",
					"description": "超时时间（秒）",
					"default":     30,
				},
				"capture_output": map[string]interface{}{
					"type":        "boolean",
					"description": "是否捕获输出",
					"default":     true,
				},
				"env": map[string]interface{}{
					"type":        "object",
					"description": "环境变量",
				},
			},
			"required": []string{"command"},
		},
		Returns: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"success": map[string]interface{}{
					"type":        "boolean",
					"description": "执行是否成功",
				},
				"exit_code": map[string]interface{}{
					"type":        "integer",
					"description": "退出码",
				},
				"stdout": map[string]interface{}{
					"type":        "string",
					"description": "标准输出",
				},
				"stderr": map[string]interface{}{
					"type":        "string",
					"description": "标准错误",
				},
				"duration": map[string]interface{}{
					"type":        "number",
					"description": "执行时间（秒）",
				},
				"error": map[string]interface{}{
					"type":        "string",
					"description": "错误信息",
				},
			},
		},
		Enabled: true,
		Tags:    []string{"shell", "command", "exec"},
	}

	tool := &ShellTool{
		BaseTool:        NewBaseTool(definition),
		workspacePath:   workspacePath,
		allowedCommands: allowedCommands,
		commandTimeout:  30 * time.Second,
	}

	return tool, nil
}

// Execute 执行命令行操作
func (t *ShellTool) Execute(ctx context.Context, execution ToolExecution) (*ToolResult, error) {
	// 获取参数
	command, _ := execution.Parameters["command"].(string)
	args, _ := execution.Parameters["args"].([]interface{})
	workingDir, _ := execution.Parameters["working_dir"].(string)
	timeout, _ := execution.Parameters["timeout"].(float64)
	captureOutput, _ := execution.Parameters["capture_output"].(bool)
	env, _ := execution.Parameters["env"].(map[string]interface{})

	// 验证命令安全性
	if err := t.validateCommand(command); err != nil {
		return &ToolResult{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// 验证工作目录
	safeWorkingDir, err := t.validateWorkingDir(workingDir)
	if err != nil {
		return &ToolResult{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// 转换参数
	var cmdArgs []string
	for _, arg := range args {
		if str, ok := arg.(string); ok {
			cmdArgs = append(cmdArgs, str)
		}
	}

	// 转换环境变量
	cmdEnv := os.Environ()
	for key, value := range env {
		if str, ok := value.(string); ok {
			cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", key, str))
		}
	}

	// 执行命令
	startTime := time.Now()
	result, execErr := t.executeCommand(command, cmdArgs, safeWorkingDir, cmdEnv, time.Duration(timeout)*time.Second, captureOutput)
	duration := time.Since(startTime)

	// 构建结果
	success := execErr == nil
	if result != nil {
		if exitCode, ok := result["exit_code"].(int); ok {
			success = success && exitCode == 0
		}
	}

	toolResult := &ToolResult{
		Success:   success,
		Output:    result,
		Duration:  duration,
		Timestamp: startTime,
	}

	if execErr != nil {
		toolResult.Error = execErr.Error()
	}

	return toolResult, nil
}

// validateCommand 验证命令安全性
func (t *ShellTool) validateCommand(command string) error {
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// 获取命令名称（第一个单词）
	cmdName := strings.Fields(command)[0]

	// 检查是否在白名单中
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if len(t.allowedCommands) > 0 {
		allowed := false
		for _, allowedCmd := range t.allowedCommands {
			if cmdName == allowedCmd {
				allowed = true
				break
			}
		}

		if !allowed {
			return fmt.Errorf("command '%s' is not in the allowed list", cmdName)
		}
	}

	// 检查危险命令模式
	dangerousPatterns := []string{
		"rm -rf /", "rm -rf /*", "rm -rf .", "rm -rf *",
		":(){ :|:& };:", // fork炸弹
		"mkfs", "dd if=", "shutdown", "halt", "reboot",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(command, pattern) {
			return fmt.Errorf("command contains dangerous pattern: %s", pattern)
		}
	}

	return nil
}

// validateWorkingDir 验证工作目录
func (t *ShellTool) validateWorkingDir(workingDir string) (string, error) {
	if workingDir == "" || workingDir == "." {
		return t.workspacePath, nil
	}

	// 清理路径
	cleanPath := filepath.Clean(workingDir)

	// 检查是否为绝对路径，如果是则转换为相对于工作区的路径
	if filepath.IsAbs(cleanPath) {
		// 确保路径在工作区内
		relPath, err := filepath.Rel(t.workspacePath, cleanPath)
		if err != nil {
			return "", fmt.Errorf("path is outside workspace: %s", cleanPath)
		}

		// 检查是否尝试向上访问
		if strings.HasPrefix(relPath, "..") {
			return "", fmt.Errorf("path traversal not allowed: %s", cleanPath)
		}

		cleanPath = filepath.Join(t.workspacePath, relPath)
	} else {
		// 相对路径，直接连接到工作区
		cleanPath = filepath.Join(t.workspacePath, cleanPath)
	}

	// 确保目录存在
	if err := os.MkdirAll(cleanPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create working directory: %w", err)
	}

	return cleanPath, nil
}

// executeCommand 执行命令
func (t *ShellTool) executeCommand(command string, args []string, workingDir string, env []string, timeout time.Duration, captureOutput bool) (map[string]interface{}, error) {
	// 创建命令上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 创建命令
	var cmd *exec.Cmd
	if len(args) > 0 {
		cmd = exec.CommandContext(ctx, command, args...)
	} else {
		// 如果命令包含空格，使用shell执行
		if strings.Contains(command, " ") {
			cmd = exec.CommandContext(ctx, "sh", "-c", command)
		} else {
			cmd = exec.CommandContext(ctx, command)
		}
	}

	// 设置工作目录和环境变量
	cmd.Dir = workingDir
	cmd.Env = env

	// 捕获输出
	var stdout, stderr []byte
	var execErr error

	if captureOutput {
		stdout, execErr = cmd.Output()
		if execErr != nil {
			if exitErr, ok := execErr.(*exec.ExitError); ok {
				stderr = exitErr.Stderr
			}
		}
	} else {
		// 不捕获输出，直接执行
		execErr = cmd.Run()
		if execErr != nil {
			if exitErr, ok := execErr.(*exec.ExitError); ok {
				stderr = exitErr.Stderr
			}
		}
	}

	// 获取退出码
	exitCode := 0
	if execErr != nil {
		if exitErr, ok := execErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			// 非退出错误（如超时、权限错误等）
			exitCode = -1
		}
	}

	// 检查是否超时
	if ctx.Err() == context.DeadlineExceeded {
		return map[string]interface{}{
			"command":     command,
			"args":        args,
			"working_dir": workingDir,
			"exit_code":   -2,
			"stdout":      string(stdout),
			"stderr":      string(stderr),
			"timeout":     true,
			"error":       "command execution timed out",
		}, nil
	}

	result := map[string]interface{}{
		"command":     command,
		"args":        args,
		"working_dir": workingDir,
		"exit_code":   exitCode,
		"stdout":      string(stdout),
		"stderr":      string(stderr),
		"timeout":     false,
	}

	if execErr != nil {
		result["error"] = execErr.Error()
	}

	return result, nil
}

// SetAllowedCommands 设置允许的命令列表
func (t *ShellTool) SetAllowedCommands(commands []string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.allowedCommands = commands
}

// GetAllowedCommands 获取允许的命令列表
func (t *ShellTool) GetAllowedCommands() []string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.allowedCommands
}

// SetCommandTimeout 设置命令超时时间
func (t *ShellTool) SetCommandTimeout(timeout time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.commandTimeout = timeout
}

// GetCommandTimeout 获取命令超时时间
func (t *ShellTool) GetCommandTimeout() time.Duration {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.commandTimeout
}

// AddAllowedCommand 添加允许的命令
func (t *ShellTool) AddAllowedCommand(command string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// 检查是否已存在
	for _, cmd := range t.allowedCommands {
		if cmd == command {
			return
		}
	}

	t.allowedCommands = append(t.allowedCommands, command)
}

// RemoveAllowedCommand 移除允许的命令
func (t *ShellTool) RemoveAllowedCommand(command string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var newCommands []string
	for _, cmd := range t.allowedCommands {
		if cmd != command {
			newCommands = append(newCommands, cmd)
		}
	}

	t.allowedCommands = newCommands
}

// Initialize 初始化工具
func (t *ShellTool) Initialize() error {
	// 调用父类初始化
	if err := t.BaseTool.Initialize(); err != nil {
		return err
	}

	utils.Info("命令行工具初始化完成，工作区: %s", t.workspacePath)
	utils.Info("允许的命令: %v", t.allowedCommands)
	utils.Info("默认超时: %v", t.commandTimeout)
	return nil
}
