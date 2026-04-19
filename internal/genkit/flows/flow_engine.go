package flows

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"GenPulse/internal/genkit/models"
	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/utils"
)

// FlowStatus Flow状态
type FlowStatus string

const (
	FlowStatusPending   FlowStatus = "pending"
	FlowStatusRunning   FlowStatus = "running"
	FlowStatusPaused    FlowStatus = "paused"
	FlowStatusCompleted FlowStatus = "completed"
	FlowStatusFailed    FlowStatus = "failed"
	FlowStatusCancelled FlowStatus = "cancelled"
)

// FlowType Flow类型
type FlowType string

const (
	FlowTypeSequential  FlowType = "sequential"
	FlowTypeParallel    FlowType = "parallel"
	FlowTypeConditional FlowType = "conditional"
	FlowTypeLoop        FlowType = "loop"
)

// NodeType 节点类型
type NodeType string

const (
	NodeTypeAction    NodeType = "action"
	NodeTypeCondition NodeType = "condition"
	NodeTypeParallel  NodeType = "parallel"
	NodeTypeLoop      NodeType = "loop"
	NodeTypeModel     NodeType = "model"
	NodeTypeTool      NodeType = "tool"
	NodeTypeCustom    NodeType = "custom"
)

// FlowDefinition Flow定义
type FlowDefinition struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        FlowType               `json:"type"`
	Version     string                 `json:"version"`
	Nodes       []NodeDefinition       `json:"nodes"`
	Edges       []EdgeDefinition       `json:"edges"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Returns     map[string]interface{} `json:"returns,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// NodeDefinition 节点定义
type NodeDefinition struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        NodeType               `json:"type"`
	Description string                 `json:"description,omitempty"`
	Config      map[string]interface{} `json:"config"`
	Position    Position               `json:"position,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// EdgeDefinition 边定义
type EdgeDefinition struct {
	ID          string                 `json:"id"`
	Source      string                 `json:"source"`
	Target      string                 `json:"target"`
	Condition   string                 `json:"condition,omitempty"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Position 位置信息
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// FlowExecution Flow执行
type FlowExecution struct {
	ID          string                 `json:"id"`
	FlowID      string                 `json:"flow_id"`
	Status      FlowStatus             `json:"status"`
	Parameters  map[string]interface{} `json:"parameters"`
	Results     map[string]interface{} `json:"results,omitempty"`
	Error       string                 `json:"error,omitempty"`
	StartTime   time.Time              `json:"start_time,omitempty"`
	EndTime     time.Time              `json:"end_time,omitempty"`
	Duration    time.Duration          `json:"duration,omitempty"`
	NodeResults map[string]NodeResult  `json:"node_results,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NodeResult 节点执行结果
type NodeResult struct {
	NodeID    string                 `json:"node_id"`
	Status    FlowStatus             `json:"status"`
	Output    interface{}            `json:"output,omitempty"`
	Error     string                 `json:"error,omitempty"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time,omitempty"`
	Duration  time.Duration          `json:"duration,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// FlowEngine Flow引擎
type FlowEngine struct {
	flows        map[string]*FlowDefinition
	executions   map[string]*FlowExecution
	modelAdapter *models.UnifiedModelAdapter
	toolRegistry *tools.ToolRegistry
	mutex        sync.RWMutex
}

// NewFlowEngine 创建Flow引擎
func NewFlowEngine(modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry) *FlowEngine {
	return &FlowEngine{
		flows:        make(map[string]*FlowDefinition),
		executions:   make(map[string]*FlowExecution),
		modelAdapter: modelAdapter,
		toolRegistry: toolRegistry,
	}
}

// RegisterFlow 注册Flow
func (fe *FlowEngine) RegisterFlow(definition FlowDefinition) error {
	fe.mutex.Lock()
	defer fe.mutex.Unlock()

	// 设置时间戳
	now := time.Now()
	if definition.CreatedAt.IsZero() {
		definition.CreatedAt = now
	}
	definition.UpdatedAt = now

	// 检查是否已注册
	if _, exists := fe.flows[definition.ID]; exists {
		return fmt.Errorf("flow already registered: %s", definition.ID)
	}

	// 验证Flow定义
	if err := fe.validateFlowDefinition(&definition); err != nil {
		return fmt.Errorf("invalid flow definition: %w", err)
	}

	// 注册Flow
	fe.flows[definition.ID] = &definition

	utils.Info("注册Flow: %s (%s)", definition.Name, definition.Type)
	return nil
}

// validateFlowDefinition 验证Flow定义
func (fe *FlowEngine) validateFlowDefinition(definition *FlowDefinition) error {
	// 检查必需字段
	if definition.ID == "" {
		return fmt.Errorf("flow ID is required")
	}
	if definition.Name == "" {
		return fmt.Errorf("flow name is required")
	}
	if definition.Type == "" {
		return fmt.Errorf("flow type is required")
	}

	// 检查节点
	if len(definition.Nodes) == 0 {
		return fmt.Errorf("flow must have at least one node")
	}

	// 检查节点ID唯一性
	nodeIDs := make(map[string]bool)
	for _, node := range definition.Nodes {
		if node.ID == "" {
			return fmt.Errorf("node ID is required")
		}
		if nodeIDs[node.ID] {
			return fmt.Errorf("duplicate node ID: %s", node.ID)
		}
		nodeIDs[node.ID] = true
	}

	// 检查边
	for _, edge := range definition.Edges {
		if edge.Source == "" || edge.Target == "" {
			return fmt.Errorf("edge source and target are required")
		}
		if !nodeIDs[edge.Source] {
			return fmt.Errorf("edge source node not found: %s", edge.Source)
		}
		if !nodeIDs[edge.Target] {
			return fmt.Errorf("edge target node not found: %s", edge.Target)
		}
	}

	return nil
}

// ExecuteFlow 执行Flow
func (fe *FlowEngine) ExecuteFlow(ctx context.Context, flowID string, parameters map[string]interface{}) (*FlowExecution, error) {
	// 获取Flow定义
	definition, err := fe.GetFlow(flowID)
	if err != nil {
		return nil, err
	}

	// 创建执行记录
	execution := &FlowExecution{
		ID:          generateExecutionID(flowID),
		FlowID:      flowID,
		Status:      FlowStatusPending,
		Parameters:  parameters,
		Results:     make(map[string]interface{}),
		NodeResults: make(map[string]NodeResult),
		Context:     make(map[string]interface{}),
		Metadata:    make(map[string]interface{}),
		StartTime:   time.Now(),
	}

	// 保存执行记录
	fe.mutex.Lock()
	fe.executions[execution.ID] = execution
	fe.mutex.Unlock()

	// 异步执行
	go fe.executeFlowAsync(ctx, definition, execution)

	return execution, nil
}

// executeFlowAsync 异步执行Flow
func (fe *FlowEngine) executeFlowAsync(ctx context.Context, definition *FlowDefinition, execution *FlowExecution) {
	// 更新状态
	execution.Status = FlowStatusRunning
	fe.updateExecution(execution)

	utils.Info("开始执行Flow: %s (%s)", definition.Name, execution.ID)

	var err error
	var results map[string]interface{}

	// 根据Flow类型执行
	switch definition.Type {
	case FlowTypeSequential:
		results, err = fe.executeSequentialFlow(ctx, definition, execution)
	case FlowTypeParallel:
		results, err = fe.executeParallelFlow(ctx, definition, execution)
	case FlowTypeConditional:
		results, err = fe.executeConditionalFlow(ctx, definition, execution)
	case FlowTypeLoop:
		results, err = fe.executeLoopFlow(ctx, definition, execution)
	default:
		err = fmt.Errorf("unsupported flow type: %s", definition.Type)
	}

	// 更新执行结果
	execution.EndTime = time.Now()
	execution.Duration = execution.EndTime.Sub(execution.StartTime)

	if err != nil {
		execution.Status = FlowStatusFailed
		execution.Error = err.Error()
		utils.Error("Flow执行失败 %s: %v", execution.ID, err)
	} else {
		execution.Status = FlowStatusCompleted
		execution.Results = results
		utils.Info("Flow执行完成: %s (耗时: %v)", execution.ID, execution.Duration)
	}

	fe.updateExecution(execution)
}

// executeSequentialFlow 执行顺序Flow
func (fe *FlowEngine) executeSequentialFlow(ctx context.Context, definition *FlowDefinition, execution *FlowExecution) (map[string]interface{}, error) {
	// 构建执行顺序（基于边的拓扑排序）
	executionOrder, err := fe.getExecutionOrder(definition)
	if err != nil {
		return nil, fmt.Errorf("failed to determine execution order: %w", err)
	}

	// 顺序执行节点
	context := make(map[string]interface{})
	for _, nodeID := range executionOrder {
		node := fe.getNodeByID(definition, nodeID)
		if node == nil {
			return nil, fmt.Errorf("node not found: %s", nodeID)
		}

		// 执行节点
		result, err := fe.executeNode(ctx, node, context, execution)
		if err != nil {
			return nil, fmt.Errorf("node execution failed %s: %w", nodeID, err)
		}

		// 更新上下文
		context[nodeID] = result.Output

		// 记录节点结果
		execution.NodeResults[nodeID] = *result
		fe.updateExecution(execution)

		// 检查是否应该继续执行
		if result.Status == FlowStatusFailed {
			return nil, fmt.Errorf("node %s failed: %s", nodeID, result.Error)
		}
	}

	return context, nil
}

// executeParallelFlow 执行并行Flow
func (fe *FlowEngine) executeParallelFlow(ctx context.Context, definition *FlowDefinition, execution *FlowExecution) (map[string]interface{}, error) {
	// 找出可以并行执行的节点组
	nodeGroups, err := fe.getParallelNodeGroups(definition)
	if err != nil {
		return nil, fmt.Errorf("failed to determine parallel groups: %w", err)
	}

	context := make(map[string]interface{})

	// 按组顺序执行
	for _, group := range nodeGroups {
		// 并行执行组内节点
		results := make(map[string]*NodeResult)
		errors := make(chan error, len(group))

		var wg sync.WaitGroup
		for _, nodeID := range group {
			wg.Add(1)

			go func(nid string) {
				defer wg.Done()

				node := fe.getNodeByID(definition, nid)
				if node == nil {
					errors <- fmt.Errorf("node not found: %s", nid)
					return
				}

				result, err := fe.executeNode(ctx, node, context, execution)
				if err != nil {
					errors <- fmt.Errorf("node %s failed: %w", nid, err)
					return
				}

				results[nid] = result
			}(nodeID)
		}

		wg.Wait()
		close(errors)

		// 检查错误
		for err := range errors {
			return nil, err
		}

		// 更新上下文和记录结果
		for nodeID, result := range results {
			context[nodeID] = result.Output
			execution.NodeResults[nodeID] = *result
		}

		fe.updateExecution(execution)
	}

	return context, nil
}

// executeConditionalFlow 执行条件Flow
func (fe *FlowEngine) executeConditionalFlow(ctx context.Context, definition *FlowDefinition, execution *FlowExecution) (map[string]interface{}, error) {
	// 简化实现：执行第一个条件节点，根据结果选择分支
	// TODO: 实现完整的条件逻辑

	context := make(map[string]interface{})

	// 找到条件节点
	var conditionNode *NodeDefinition
	for _, node := range definition.Nodes {
		if node.Type == NodeTypeCondition {
			conditionNode = &node
			break
		}
	}

	if conditionNode == nil {
		return nil, fmt.Errorf("no condition node found in conditional flow")
	}

	// 执行条件节点
	conditionResult, err := fe.executeNode(ctx, conditionNode, context, execution)
	if err != nil {
		return nil, fmt.Errorf("condition node failed: %w", err)
	}

	// 记录条件节点结果
	execution.NodeResults[conditionNode.ID] = *conditionResult
	fe.updateExecution(execution)

	// 根据条件结果选择分支
	// 这里简化处理，实际应该根据边的条件判断
	context[conditionNode.ID] = conditionResult.Output

	// 执行选中的分支
	// TODO: 实现分支选择逻辑

	return context, nil
}

// executeLoopFlow 执行循环Flow
func (fe *FlowEngine) executeLoopFlow(ctx context.Context, definition *FlowDefinition, execution *FlowExecution) (map[string]interface{}, error) {
	// 简化实现：执行固定次数的循环
	// TODO: 实现完整的循环逻辑

	context := make(map[string]interface{})

	// 获取循环配置
	maxIterations := 3 // 默认3次
	if config, ok := definition.Parameters["max_iterations"].(float64); ok {
		maxIterations = int(config)
	}

	// 找到循环体节点
	var loopBodyNodes []*NodeDefinition
	for i := range definition.Nodes {
		node := &definition.Nodes[i]
		if node.Type == NodeTypeLoop || (node.Config != nil && node.Config["is_loop_body"] == true) {
			loopBodyNodes = append(loopBodyNodes, node)
		}
	}

	if len(loopBodyNodes) == 0 {
		return nil, fmt.Errorf("no loop body nodes found")
	}

	// 执行循环
	for iteration := 0; iteration < maxIterations; iteration++ {
		utils.Debug("循环迭代 %d/%d", iteration+1, maxIterations)

		// 设置迭代上下文
		context["iteration"] = iteration
		context["max_iterations"] = maxIterations

		// 执行循环体节点
		for _, node := range loopBodyNodes {
			result, err := fe.executeNode(ctx, node, context, execution)
			if err != nil {
				return nil, fmt.Errorf("loop iteration %d failed: %w", iteration, err)
			}

			// 记录节点结果（带迭代信息）
			result.Metadata["iteration"] = iteration
			execution.NodeResults[fmt.Sprintf("%s_iter%d", node.ID, iteration)] = *result
		}

		fe.updateExecution(execution)

		// 检查是否应该提前退出循环
		// TODO: 实现退出条件检查
	}

	return context, nil
}

// executeNode 执行单个节点
func (fe *FlowEngine) executeNode(ctx context.Context, node *NodeDefinition, context map[string]interface{}, execution *FlowExecution) (*NodeResult, error) {
	startTime := time.Now()

	result := &NodeResult{
		NodeID:    node.ID,
		Status:    FlowStatusRunning,
		StartTime: startTime,
		Metadata:  make(map[string]interface{}),
	}

	utils.Debug("执行节点: %s (%s)", node.Name, node.Type)

	var output interface{}
	var err error

	// 根据节点类型执行
	switch node.Type {
	case NodeTypeModel:
		output, err = fe.executeModelNode(ctx, node, context)
	case NodeTypeTool:
		output, err = fe.executeToolNode(ctx, node, context)
	case NodeTypeAction:
		output, err = fe.executeActionNode(ctx, node, context)
	case NodeTypeCondition:
		output, err = fe.executeConditionNode(ctx, node, context)
	default:
		err = fmt.Errorf("unsupported node type: %s", node.Type)
	}

	// 更新结果
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(startTime)

	if err != nil {
		result.Status = FlowStatusFailed
		result.Error = err.Error()
		utils.Error("节点执行失败 %s: %v", node.ID, err)
	} else {
		result.Status = FlowStatusCompleted
		result.Output = output
		utils.Debug("节点执行完成: %s (耗时: %v)", node.ID, result.Duration)
	}

	result.Metadata["node_type"] = string(node.Type)
	result.Metadata["duration"] = result.Duration.Seconds()

	return result, nil
}

// executeModelNode 执行模型节点
func (fe *FlowEngine) executeModelNode(ctx context.Context, node *NodeDefinition, context map[string]interface{}) (interface{}, error) {
	if fe.modelAdapter == nil {
		return nil, fmt.Errorf("model adapter not available")
	}

	// 获取模型配置
	modelID, _ := node.Config["model_id"].(string)
	prompt, _ := node.Config["prompt"].(string)

	// 处理模板变量
	resolvedPrompt := fe.resolveTemplate(prompt, context)

	// 构建模型请求
	req := models.ModelRequest{
		Prompt:      resolvedPrompt,
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	// 调用模型
	response, err := fe.modelAdapter.Generate(ctx, modelID, req)
	if err != nil {
		return nil, fmt.Errorf("model generation failed: %w", err)
	}

	return map[string]interface{}{
		"content": response.Content,
		"usage":   response.Usage,
		"model":   modelID,
	}, nil
}

// executeToolNode 执行工具节点
func (fe *FlowEngine) executeToolNode(ctx context.Context, node *NodeDefinition, context map[string]interface{}) (interface{}, error) {
	if fe.toolRegistry == nil {
		return nil, fmt.Errorf("tool registry not available")
	}

	// 获取工具配置
	toolID, _ := node.Config["tool_id"].(string)
	parameters, _ := node.Config["parameters"].(map[string]interface{})

	// 解析参数中的模板变量
	resolvedParams := fe.resolveTemplatesInMap(parameters, context)

	// 构建工具执行请求
	execution := tools.ToolExecution{
		ToolID:     toolID,
		Parameters: resolvedParams,
		Context:    context,
	}

	// 执行工具
	result, err := fe.toolRegistry.ExecuteTool(ctx, execution)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	return result, nil
}

// executeActionNode 执行动作节点
func (fe *FlowEngine) executeActionNode(ctx context.Context, node *NodeDefinition, context map[string]interface{}) (interface{}, error) {
	// 动作节点可以执行各种自定义操作
	// 这里实现一个简单的示例动作

	actionType, _ := node.Config["action_type"].(string)

	switch actionType {
	case "log":
		message, _ := node.Config["message"].(string)
		resolvedMessage := fe.resolveTemplate(message, context)
		utils.Info("动作节点日志: %s", resolvedMessage)
		return map[string]interface{}{
			"action":  "log",
			"message": resolvedMessage,
			"logged":  true,
		}, nil

	case "transform":
		// 数据转换动作
		input, _ := node.Config["input"].(string)
		operation, _ := node.Config["operation"].(string)

		resolvedInput := fe.resolveTemplate(input, context)

		var output string
		switch operation {
		case "uppercase":
			output = strings.ToUpper(resolvedInput)
		case "lowercase":
			output = strings.ToLower(resolvedInput)
		case "trim":
			output = strings.TrimSpace(resolvedInput)
		default:
			output = resolvedInput
		}

		return map[string]interface{}{
			"action":    "transform",
			"input":     resolvedInput,
			"operation": operation,
			"output":    output,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported action type: %s", actionType)
	}
}

// executeConditionNode 执行条件节点
func (fe *FlowEngine) executeConditionNode(ctx context.Context, node *NodeDefinition, context map[string]interface{}) (interface{}, error) {
	// 条件节点评估条件并返回布尔值
	condition, _ := node.Config["condition"].(string)

	// 解析条件中的变量
	resolvedCondition := fe.resolveTemplate(condition, context)

	// 简化条件评估
	// 实际实现应该使用表达式求值引擎
	isTrue := false

	// 简单字符串匹配作为示例
	if strings.Contains(strings.ToLower(resolvedCondition), "true") {
		isTrue = true
	} else if strings.Contains(strings.ToLower(resolvedCondition), "yes") {
		isTrue = true
	} else if strings.Contains(resolvedCondition, ">") {
		// 简单比较
		parts := strings.Split(resolvedCondition, ">")
		if len(parts) == 2 {
			// 这里应该解析和比较数值
			isTrue = true // 简化处理
		}
	}

	return map[string]interface{}{
		"condition": resolvedCondition,
		"result":    isTrue,
	}, nil
}

// resolveTemplate 解析模板变量
func (fe *FlowEngine) resolveTemplate(template string, context map[string]interface{}) string {
	if template == "" {
		return ""
	}

	// 简单模板变量替换：{{variable}}
	result := template

	for key, value := range context {
		placeholder := fmt.Sprintf("{{%s}}", key)
		if strValue, ok := value.(string); ok {
			result = strings.ReplaceAll(result, placeholder, strValue)
		} else {
			// 尝试转换为字符串
			jsonValue, err := json.Marshal(value)
			if err == nil {
				result = strings.ReplaceAll(result, placeholder, string(jsonValue))
			}
		}
	}

	return result
}

// resolveTemplatesInMap 解析map中的模板变量
func (fe *FlowEngine) resolveTemplatesInMap(data map[string]interface{}, context map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		switch v := value.(type) {
		case string:
			result[key] = fe.resolveTemplate(v, context)
		case map[string]interface{}:
			result[key] = fe.resolveTemplatesInMap(v, context)
		case []interface{}:
			result[key] = fe.resolveTemplatesInSlice(v, context)
		default:
			result[key] = value
		}
	}

	return result
}

// resolveTemplatesInSlice 解析slice中的模板变量
func (fe *FlowEngine) resolveTemplatesInSlice(data []interface{}, context map[string]interface{}) []interface{} {
	result := make([]interface{}, len(data))

	for i, value := range data {
		switch v := value.(type) {
		case string:
			result[i] = fe.resolveTemplate(v, context)
		case map[string]interface{}:
			result[i] = fe.resolveTemplatesInMap(v, context)
		case []interface{}:
			result[i] = fe.resolveTemplatesInSlice(v, context)
		default:
			result[i] = value
		}
	}

	return result
}

// getExecutionOrder 获取执行顺序（拓扑排序）
func (fe *FlowEngine) getExecutionOrder(definition *FlowDefinition) ([]string, error) {
	// 简化实现：返回节点定义顺序
	// TODO: 实现基于边的拓扑排序

	var order []string
	for _, node := range definition.Nodes {
		order = append(order, node.ID)
	}

	return order, nil
}

// getParallelNodeGroups 获取并行节点组
func (fe *FlowEngine) getParallelNodeGroups(definition *FlowDefinition) ([][]string, error) {
	// 简化实现：将所有节点放在一个组中并行执行
	// TODO: 实现基于依赖关系的分组

	var group []string
	for _, node := range definition.Nodes {
		group = append(group, node.ID)
	}

	return [][]string{group}, nil
}

// getNodeByID 根据ID获取节点
func (fe *FlowEngine) getNodeByID(definition *FlowDefinition, nodeID string) *NodeDefinition {
	for i := range definition.Nodes {
		if definition.Nodes[i].ID == nodeID {
			return &definition.Nodes[i]
		}
	}
	return nil
}

// updateExecution 更新执行记录
func (fe *FlowEngine) updateExecution(execution *FlowExecution) {
	fe.mutex.Lock()
	defer fe.mutex.Unlock()

	fe.executions[execution.ID] = execution
}

// GetFlow 获取Flow定义
func (fe *FlowEngine) GetFlow(flowID string) (*FlowDefinition, error) {
	fe.mutex.RLock()
	defer fe.mutex.RUnlock()

	definition, exists := fe.flows[flowID]
	if !exists {
		return nil, fmt.Errorf("flow not found: %s", flowID)
	}

	return definition, nil
}

// GetExecution 获取执行记录
func (fe *FlowEngine) GetExecution(executionID string) (*FlowExecution, error) {
	fe.mutex.RLock()
	defer fe.mutex.RUnlock()

	execution, exists := fe.executions[executionID]
	if !exists {
		return nil, fmt.Errorf("execution not found: %s", executionID)
	}

	return execution, nil
}

// ListFlows 列出所有Flow
func (fe *FlowEngine) ListFlows() []FlowDefinition {
	fe.mutex.RLock()
	defer fe.mutex.RUnlock()

	var flows []FlowDefinition
	for _, flow := range fe.flows {
		flows = append(flows, *flow)
	}

	return flows
}

// ListExecutions 列出所有执行记录
func (fe *FlowEngine) ListExecutions() []FlowExecution {
	fe.mutex.RLock()
	defer fe.mutex.RUnlock()

	var executions []FlowExecution
	for _, execution := range fe.executions {
		executions = append(executions, *execution)
	}

	return executions
}

// CancelExecution 取消执行
func (fe *FlowEngine) CancelExecution(executionID string) error {
	fe.mutex.Lock()
	defer fe.mutex.Unlock()

	execution, exists := fe.executions[executionID]
	if !exists {
		return fmt.Errorf("execution not found: %s", executionID)
	}

	if execution.Status == FlowStatusCompleted || execution.Status == FlowStatusFailed || execution.Status == FlowStatusCancelled {
		return fmt.Errorf("cannot cancel execution in status: %s", execution.Status)
	}

	execution.Status = FlowStatusCancelled
	execution.EndTime = time.Now()
	execution.Duration = execution.EndTime.Sub(execution.StartTime)

	utils.Info("取消Flow执行: %s", executionID)
	return nil
}

// GetStatistics 获取统计信息
func (fe *FlowEngine) GetStatistics() map[string]interface{} {
	fe.mutex.RLock()
	defer fe.mutex.RUnlock()

	stats := make(map[string]interface{})
	stats["total_flows"] = len(fe.flows)
	stats["total_executions"] = len(fe.executions)

	// 按状态统计执行
	statusCounts := make(map[string]int)
	for _, execution := range fe.executions {
		statusCounts[string(execution.Status)]++
	}
	stats["executions_by_status"] = statusCounts

	return stats
}

// generateExecutionID 生成执行ID
func generateExecutionID(flowID string) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d", flowID, timestamp)
}

// Global flow engine instance
var globalFlowEngine *FlowEngine

// InitGlobalFlowEngine 初始化全局Flow引擎
func InitGlobalFlowEngine(modelAdapter *models.UnifiedModelAdapter, toolRegistry *tools.ToolRegistry) error {
	if globalFlowEngine == nil {
		globalFlowEngine = NewFlowEngine(modelAdapter, toolRegistry)
		utils.Info("初始化全局Flow引擎")
	}
	return nil
}

// GetGlobalFlowEngine 获取全局Flow引擎
func GetGlobalFlowEngine() *FlowEngine {
	return globalFlowEngine
}
