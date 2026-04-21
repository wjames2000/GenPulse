package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/utils"
)

// AgentTool 将Agent封装为Tool
type AgentTool struct {
	agent      Agent
	definition tools.ToolDefinition
	baseTool   *tools.BaseTool
}

// NewAgentTool 创建AgentTool
func NewAgentTool(agent Agent) (*AgentTool, error) {
	if agent == nil {
		return nil, fmt.Errorf("agent不能为空")
	}

	config := agent.GetConfig()

	// 创建工具定义
	definition := tools.ToolDefinition{
		ID:          fmt.Sprintf("agent_%s", config.ID),
		Name:        fmt.Sprintf("%s Agent", config.Name),
		Description: fmt.Sprintf("调用%s Agent执行任务", config.Name),
		Category:    tools.ToolCategoryCustom,
		Version:     "1.0.0",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"task": map[string]interface{}{
					"type":        "string",
					"description": "要执行的任务描述",
				},
				"parameters": map[string]interface{}{
					"type":                 "object",
					"description":          "任务参数",
					"additionalProperties": true,
				},
			},
			"required": []string{"task"},
		},
		Returns: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"success": map[string]interface{}{
					"type":        "boolean",
					"description": "执行是否成功",
				},
				"result": map[string]interface{}{
					"type":                 "object",
					"description":          "执行结果",
					"additionalProperties": true,
				},
				"error": map[string]interface{}{
					"type":        "string",
					"description": "错误信息（如果失败）",
				},
				"agent_id": map[string]interface{}{
					"type":        "string",
					"description": "执行Agent的ID",
				},
				"agent_name": map[string]interface{}{
					"type":        "string",
					"description": "执行Agent的名称",
				},
				"execution_time": map[string]interface{}{
					"type":        "number",
					"description": "执行耗时（秒）",
				},
			},
		},
		Enabled: agent.IsEnabled(),
		Tags: []string{
			"agent",
			string(config.Role),
			"ai",
		},
		Metadata: map[string]interface{}{
			"agent_id":     config.ID,
			"agent_role":   config.Role,
			"capabilities": config.Capabilities,
		},
	}

	baseTool := tools.NewBaseTool(definition)

	return &AgentTool{
		agent:      agent,
		definition: definition,
		baseTool:   baseTool,
	}, nil
}

// GetDefinition 获取工具定义
func (at *AgentTool) GetDefinition() tools.ToolDefinition {
	return at.definition
}

// GetCategory 获取工具类别
func (at *AgentTool) GetCategory() tools.ToolCategory {
	return at.definition.Category
}

// Execute 执行Agent工具
func (at *AgentTool) Execute(ctx context.Context, execution tools.ToolExecution) (*tools.ToolResult, error) {
	startTime := time.Now()

	// 验证参数
	if err := at.ValidateParameters(execution.Parameters); err != nil {
		return &tools.ToolResult{
			Success:   false,
			Error:     fmt.Sprintf("参数验证失败: %v", err),
			Duration:  time.Since(startTime),
			Timestamp: startTime,
		}, nil
	}

	// 提取任务参数
	task, ok := execution.Parameters["task"].(string)
	if !ok {
		return &tools.ToolResult{
			Success:   false,
			Error:     "参数'task'必须是字符串类型",
			Duration:  time.Since(startTime),
			Timestamp: startTime,
		}, nil
	}

	// 提取其他参数
	var parameters map[string]interface{}
	if params, ok := execution.Parameters["parameters"]; ok {
		if paramMap, ok := params.(map[string]interface{}); ok {
			parameters = paramMap
		}
	}

	if parameters == nil {
		parameters = make(map[string]interface{})
	}

	// 添加上下文信息
	if execution.Context != nil {
		parameters["context"] = execution.Context
	}

	utils.Info("通过工具调用Agent: %s, 任务: %s", at.agent.GetConfig().Name, task)

	// 执行Agent任务
	result, err := at.agent.Execute(ctx, task, parameters)
	duration := time.Since(startTime)

	// 构建工具结果
	toolResult := &tools.ToolResult{
		Success:   err == nil,
		Duration:  duration,
		Timestamp: startTime,
		Metadata: map[string]interface{}{
			"agent_id":   at.agent.GetConfig().ID,
			"agent_name": at.agent.GetConfig().Name,
			"agent_role": at.agent.GetConfig().Role,
			"task":       task,
		},
	}

	if err != nil {
		toolResult.Error = err.Error()
		toolResult.Output = map[string]interface{}{
			"success":  false,
			"error":    err.Error(),
			"agent_id": at.agent.GetConfig().ID,
		}
	} else if result != nil {
		// 序列化Agent结果
		resultJSON, jsonErr := json.Marshal(result)
		if jsonErr != nil {
			utils.Warn("序列化Agent结果失败: %v", jsonErr)
			toolResult.Output = map[string]interface{}{
				"success":  true,
				"agent_id": at.agent.GetConfig().ID,
				"result":   "执行成功，但结果序列化失败",
			}
		} else {
			var resultMap map[string]interface{}
			if jsonErr := json.Unmarshal(resultJSON, &resultMap); jsonErr != nil {
				utils.Warn("解析Agent结果失败: %v", jsonErr)
				toolResult.Output = map[string]interface{}{
					"success":  true,
					"agent_id": at.agent.GetConfig().ID,
					"result":   "执行成功，但结果解析失败",
				}
			} else {
				toolResult.Output = map[string]interface{}{
					"success":        true,
					"agent_id":       at.agent.GetConfig().ID,
					"agent_name":     at.agent.GetConfig().Name,
					"result":         resultMap,
					"execution_time": duration.Seconds(),
				}
			}
		}
	} else {
		toolResult.Output = map[string]interface{}{
			"success":  true,
			"agent_id": at.agent.GetConfig().ID,
			"message":  "任务执行完成，无返回结果",
		}
	}

	// 更新统计信息
	at.baseTool.IncrementExecutionCount(duration)

	utils.Info("Agent工具执行完成: %s, 耗时: %v, 成功: %v",
		at.agent.GetConfig().Name, duration, toolResult.Success)

	return toolResult, nil
}

// ValidateParameters 验证参数
func (at *AgentTool) ValidateParameters(parameters map[string]interface{}) error {
	if parameters == nil {
		return fmt.Errorf("参数不能为空")
	}

	task, ok := parameters["task"]
	if !ok {
		return fmt.Errorf("缺少必需的参数'task'")
	}

	if _, ok := task.(string); !ok {
		return fmt.Errorf("参数'task'必须是字符串类型")
	}

	// 验证parameters参数（如果存在）
	if params, ok := parameters["parameters"]; ok {
		if _, ok := params.(map[string]interface{}); !ok {
			return fmt.Errorf("参数'parameters'必须是对象类型")
		}
	}

	return nil
}

// Initialize 初始化工具
func (at *AgentTool) Initialize() error {
	// Agent已经在AgentManager中初始化，这里不需要额外操作
	return nil
}

// Shutdown 关闭工具
func (at *AgentTool) Shutdown() error {
	// Agent的关闭由AgentManager管理，这里不需要额外操作
	return nil
}

// IsEnabled 检查工具是否启用
func (at *AgentTool) IsEnabled() bool {
	return at.agent.IsEnabled()
}

// SetEnabled 设置工具启用状态
func (at *AgentTool) SetEnabled(enabled bool) {
	at.agent.SetEnabled(enabled)
	at.definition.Enabled = enabled
}

// GetExecutionCount 获取执行次数
func (at *AgentTool) GetExecutionCount() int {
	return at.baseTool.GetExecutionCount()
}

// GetAverageDuration 获取平均执行时间
func (at *AgentTool) GetAverageDuration() time.Duration {
	return at.baseTool.GetAverageDuration()
}

// GetLastExecutionTime 获取最后执行时间
func (at *AgentTool) GetLastExecutionTime() time.Time {
	return at.baseTool.GetLastExecutionTime()
}
