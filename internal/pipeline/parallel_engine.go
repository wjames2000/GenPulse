package pipeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"GenPulse/internal/agents"
	"GenPulse/internal/utils"
)

// ParallelTask 并行任务定义
type ParallelTask struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	AgentID     string                 `json:"agent_id"`
	Task        string                 `json:"task"`
	Parameters  map[string]interface{} `json:"parameters"`
	Priority    int                    `json:"priority"` // 优先级，数字越小优先级越高
	Timeout     time.Duration          `json:"timeout"`
	RetryCount  int                    `json:"retry_count"`
}

// ParallelTaskResult 并行任务结果
type ParallelTaskResult struct {
	TaskID     string                 `json:"task_id"`
	Success    bool                   `json:"success"`
	Output     interface{}            `json:"output"`
	Error      error                  `json:"error,omitempty"`
	Duration   time.Duration          `json:"duration"`
	RetryCount int                    `json:"retry_count"`
	AgentName  string                 `json:"agent_name"`
	AgentRole  string                 `json:"agent_role"`
	Artifacts  []agents.AgentArtifact `json:"artifacts,omitempty"`
}

// ParallelEngine 并行执行引擎
type ParallelEngine struct {
	agentManager *agents.AgentManager
	maxWorkers   int
	taskQueue    chan ParallelTask
	results      chan ParallelTaskResult
	workerPool   []*parallelWorker
	wg           sync.WaitGroup
	mu           sync.RWMutex
	running      bool
}

// parallelWorker 并行工作器
type parallelWorker struct {
	id       int
	engine   *ParallelEngine
	taskChan chan ParallelTask
	stopChan chan bool
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewParallelEngine 创建并行执行引擎
func NewParallelEngine(agentManager *agents.AgentManager, maxWorkers int) *ParallelEngine {
	if maxWorkers <= 0 {
		maxWorkers = 5 // 默认5个worker
	}

	return &ParallelEngine{
		agentManager: agentManager,
		maxWorkers:   maxWorkers,
		taskQueue:    make(chan ParallelTask, 100), // 缓冲队列
		results:      make(chan ParallelTaskResult, 100),
		workerPool:   make([]*parallelWorker, 0, maxWorkers),
		running:      false,
	}
}

// Start 启动并行引擎
func (pe *ParallelEngine) Start(ctx context.Context) error {
	pe.mu.Lock()
	defer pe.mu.Unlock()

	if pe.running {
		return fmt.Errorf("并行引擎已经在运行")
	}

	utils.Info("启动并行执行引擎，worker数量: %d", pe.maxWorkers)

	// 创建工作器
	for i := 0; i < pe.maxWorkers; i++ {
		workerCtx, cancel := context.WithCancel(ctx)
		worker := &parallelWorker{
			id:       i + 1,
			engine:   pe,
			taskChan: make(chan ParallelTask),
			stopChan: make(chan bool),
			ctx:      workerCtx,
			cancel:   cancel,
		}
		pe.workerPool = append(pe.workerPool, worker)
		pe.wg.Add(1)
		go worker.start()
	}

	// 启动任务分发器
	pe.wg.Add(1)
	go pe.taskDispatcher()

	pe.running = true
	utils.Info("并行执行引擎启动完成")
	return nil
}

// Stop 停止并行引擎
func (pe *ParallelEngine) Stop() {
	pe.mu.Lock()
	defer pe.mu.Unlock()

	if !pe.running {
		return
	}

	utils.Info("停止并行执行引擎")

	// 停止所有worker
	for _, worker := range pe.workerPool {
		worker.stop()
	}

	// 关闭通道
	close(pe.taskQueue)
	close(pe.results)

	// 等待所有goroutine结束
	pe.wg.Wait()

	pe.running = false
	utils.Info("并行执行引擎已停止")
}

// SubmitTask 提交单个任务
func (pe *ParallelEngine) SubmitTask(task ParallelTask) error {
	pe.mu.RLock()
	defer pe.mu.RUnlock()

	if !pe.running {
		return fmt.Errorf("并行引擎未运行")
	}

	// 设置默认值
	if task.Timeout == 0 {
		task.Timeout = 5 * time.Minute
	}
	if task.RetryCount == 0 {
		task.RetryCount = 3
	}

	pe.taskQueue <- task
	utils.Debug("提交并行任务: %s (%s)", task.Name, task.AgentID)
	return nil
}

// SubmitTasks 批量提交任务
func (pe *ParallelEngine) SubmitTasks(tasks []ParallelTask) error {
	for _, task := range tasks {
		if err := pe.SubmitTask(task); err != nil {
			return err
		}
	}
	return nil
}

// GetResults 获取结果通道
func (pe *ParallelEngine) GetResults() <-chan ParallelTaskResult {
	return pe.results
}

// ExecuteParallel 执行并行任务并等待所有结果
func (pe *ParallelEngine) ExecuteParallel(ctx context.Context, tasks []ParallelTask) ([]ParallelTaskResult, error) {
	// 确保引擎正在运行
	if !pe.running {
		if err := pe.Start(ctx); err != nil {
			return nil, err
		}
		defer pe.Stop()
	}

	// 提交所有任务
	if err := pe.SubmitTasks(tasks); err != nil {
		return nil, err
	}

	// 收集结果
	results := make([]ParallelTaskResult, 0, len(tasks))
	resultMap := make(map[string]ParallelTaskResult)
	resultsChan := pe.GetResults()

	// 等待所有任务完成
	for i := 0; i < len(tasks); i++ {
		select {
		case result, ok := <-resultsChan:
			if !ok {
				return results, fmt.Errorf("结果通道已关闭")
			}
			resultMap[result.TaskID] = result
			results = append(results, result)
		case <-ctx.Done():
			return results, ctx.Err()
		}
	}

	// 按任务顺序排序结果
	sortedResults := make([]ParallelTaskResult, 0, len(tasks))
	for _, task := range tasks {
		if result, ok := resultMap[task.ID]; ok {
			sortedResults = append(sortedResults, result)
		}
	}

	return sortedResults, nil
}

// taskDispatcher 任务分发器
func (pe *ParallelEngine) taskDispatcher() {
	defer pe.wg.Done()

	utils.Debug("任务分发器启动")

	// 简单的轮询调度
	workerIndex := 0
	for task := range pe.taskQueue {
		// 选择worker
		worker := pe.workerPool[workerIndex]
		worker.taskChan <- task

		// 更新worker索引
		workerIndex = (workerIndex + 1) % len(pe.workerPool)
	}

	// 关闭所有worker的任务通道
	for _, worker := range pe.workerPool {
		close(worker.taskChan)
	}

	utils.Debug("任务分发器停止")
}

// start 启动worker
func (pw *parallelWorker) start() {
	defer pw.engine.wg.Done()

	utils.Debug("Worker %d 启动", pw.id)

	for {
		select {
		case task, ok := <-pw.taskChan:
			if !ok {
				utils.Debug("Worker %d 任务通道已关闭", pw.id)
				return
			}
			pw.executeTask(task)
		case <-pw.stopChan:
			utils.Debug("Worker %d 收到停止信号", pw.id)
			return
		case <-pw.ctx.Done():
			utils.Debug("Worker %d 上下文取消", pw.id)
			return
		}
	}
}

// stop 停止worker
func (pw *parallelWorker) stop() {
	select {
	case pw.stopChan <- true:
		// 发送停止信号成功
	default:
		// 通道已满或已关闭
	}
	pw.cancel()
}

// executeTask 执行任务
func (pw *parallelWorker) executeTask(task ParallelTask) {
	utils.Debug("Worker %d 开始执行任务: %s", pw.id, task.Name)

	startTime := time.Now()
	var result ParallelTaskResult
	var lastErr error

	// 重试逻辑
	for retry := 0; retry <= task.RetryCount; retry++ {
		// 创建带超时的上下文
		taskCtx, cancel := context.WithTimeout(pw.ctx, task.Timeout)

		// 获取Agent
		agent, err := pw.engine.agentManager.GetAgent(task.AgentID)
		if err != nil {
			lastErr = fmt.Errorf("获取Agent失败: %w", err)
			cancel()
			if retry < task.RetryCount {
				utils.Warn("Worker %d 任务 %s 重试 %d/%d: %v",
					pw.id, task.Name, retry+1, task.RetryCount, lastErr)
				time.Sleep(time.Duration(retry+1) * time.Second) // 指数退避
				continue
			}
			break
		}

		// 执行任务
		agentResult, err := agent.Execute(taskCtx, task.Task, task.Parameters)
		cancel()

		if err != nil {
			lastErr = fmt.Errorf("任务执行失败: %w", err)
			if retry < task.RetryCount {
				utils.Warn("Worker %d 任务 %s 重试 %d/%d: %v",
					pw.id, task.Name, retry+1, task.RetryCount, lastErr)
				time.Sleep(time.Duration(retry+1) * time.Second) // 指数退避
				continue
			}
			break
		}

		// 任务执行成功
		result = ParallelTaskResult{
			TaskID:     task.ID,
			Success:    agentResult.Success,
			Output:     agentResult.Output,
			Error:      nil,
			Duration:   time.Since(startTime),
			RetryCount: retry,
			AgentName:  agent.GetConfig().Name,
			AgentRole:  string(agent.GetConfig().Role),
			Artifacts:  agentResult.Artifacts,
		}

		if !agentResult.Success {
			lastErr = fmt.Errorf("Agent执行失败: %v", agentResult.Output)
			if retry < task.RetryCount {
				utils.Warn("Worker %d 任务 %s 重试 %d/%d: %v",
					pw.id, task.Name, retry+1, task.RetryCount, lastErr)
				time.Sleep(time.Duration(retry+1) * time.Second)
				continue
			}
		}

		// 成功或最终失败
		break
	}

	// 如果所有重试都失败
	if result.TaskID == "" {
		result = ParallelTaskResult{
			TaskID:     task.ID,
			Success:    false,
			Output:     nil,
			Error:      lastErr,
			Duration:   time.Since(startTime),
			RetryCount: task.RetryCount,
			AgentName:  "未知",
			AgentRole:  "未知",
		}
	}

	// 发送结果
	select {
	case pw.engine.results <- result:
		// 结果发送成功
		if result.Success {
			utils.Info("Worker %d 任务 %s 执行成功，耗时: %v",
				pw.id, task.Name, result.Duration)
		} else {
			utils.Error("Worker %d 任务 %s 执行失败: %v",
				pw.id, task.Name, result.Error)
		}
	case <-pw.ctx.Done():
		// 上下文取消，丢弃结果
		utils.Warn("Worker %d 任务 %s 结果丢弃（上下文取消）", pw.id, task.Name)
	}
}

// CreateFrontendBackendTasks 创建前端和后端并行任务
func (pe *ParallelEngine) CreateFrontendBackendTasks(pipelineCtx *PipelineContext) []ParallelTask {
	projectName := pipelineCtx.Parameters["project_name"].(string)

	// 前端开发任务
	frontendTask := ParallelTask{
		ID:          "frontend_development_" + time.Now().Format("20060102150405"),
		Name:        "前端开发",
		Description: "开发项目前端界面和组件",
		AgentID:     "frontend_dev_001",
		Task:        "开发项目前端界面",
		Parameters: map[string]interface{}{
			"project_name":        projectName,
			"project_description": pipelineCtx.Parameters["project_description"],
			"prd_document":        pipelineCtx.GetArtifact("prd_document", ""),
			"architecture_design": pipelineCtx.GetArtifact("architecture_design", ""),
			"task_plan":           pipelineCtx.GetArtifact("task_plan", ""),
			"tech_stack":          pipelineCtx.Parameters["tech_stack"],
		},
		Priority:   1,
		Timeout:    10 * time.Minute,
		RetryCount: 2,
	}

	// 后端开发任务
	backendTask := ParallelTask{
		ID:          "backend_development_" + time.Now().Format("20060102150405"),
		Name:        "后端开发",
		Description: "开发项目后端API和业务逻辑",
		AgentID:     "backend_dev_001",
		Task:        "开发项目后端服务",
		Parameters: map[string]interface{}{
			"project_name":        projectName,
			"project_description": pipelineCtx.Parameters["project_description"],
			"prd_document":        pipelineCtx.GetArtifact("prd_document", ""),
			"architecture_design": pipelineCtx.GetArtifact("architecture_design", ""),
			"task_plan":           pipelineCtx.GetArtifact("task_plan", ""),
			"tech_stack":          pipelineCtx.Parameters["tech_stack"],
		},
		Priority:   1,
		Timeout:    10 * time.Minute,
		RetryCount: 2,
	}

	return []ParallelTask{frontendTask, backendTask}
}
