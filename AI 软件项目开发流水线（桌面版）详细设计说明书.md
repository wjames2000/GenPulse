# AI 软件项目开发流水线（桌面版）详细设计说明书

基于《项目需求说明书 V4.0（Go 语言版）》及《概要设计说明书 V1.0》


## 1. 引言

### 1.1 编写目的
本文档详细描述 AI 软件项目开发流水线（桌面版）各功能模块的内部实现细节，包括类结构、接口定义、函数流程、数据结构及关键算法。文档旨在指导开发人员进行编码实现，确保系统架构的一致性和代码质量。

### 1.2 适用范围
适用于项目开发工程师、测试工程师及后续维护人员。

### 1.3 参考文档
- 《项目需求说明书 V4.0（Go 语言版）》
- 《概要设计说明书 V1.0》
- 《数据库设计说明书 V1.0》
- Genkit Go SDK 官方文档
- Wails v2 官方文档

### 1.4 术语与缩略语
| 术语/缩略语 | 说明                                           |
| ----------- | ---------------------------------------------- |
| DI          | 依赖注入（Dependency Injection）               |
| DTO         | 数据传输对象（Data Transfer Object）           |
| ORM         | 对象关系映射（本系统不使用 ORM，采用原生 SQL） |
| FSM         | 有限状态机（Finite State Machine）             |


## 2. 系统总体设计

### 2.1 模块结构图

系统采用分层模块化设计，各模块职责清晰、边界明确：

```
┌─────────────────────────────────────────────────────────────────┐
│                        前端 React 应用                           │
│   ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐          │
│   │Dashboard │ │Config    │ │Monitor   │ │Skills    │          │
│   │Page      │ │Page      │ │Page      │ │Page      │          │
│   └──────────┘ └──────────┘ └──────────┘ └──────────┘          │
└─────────────────────────────┬───────────────────────────────────┘
                              │ Wails Bindings (IPC)
┌─────────────────────────────┴───────────────────────────────────┐
│                        服务层 (Service Layer)                    │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   Wails App Service                       │   │
│  │  - PipelineService: 流水线控制                            │   │
│  │  - ConfigService: 配置管理                                │   │
│  │  - MonitorService: 状态查询与事件推送                     │   │
│  │  - SkillService: 技能管理接口                             │   │
│  │  - MemoryService: 记忆查询接口                            │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────┬───────────────────────────────────┘
                              │
┌─────────────────────────────┴───────────────────────────────────┐
│                      引擎层 (Engine Layer)                       │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                PipelineOrchestrator                      │   │
│  │  - 主 Flow 定义与执行                                     │   │
│  │  - 阶段调度与错误恢复                                     │   │
│  └─────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Agent Registry                        │   │
│  │  - Agent 实例管理                                         │   │
│  │  - 配置驱动创建                                           │   │
│  └─────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                 自进化子系统 (Evolution)                  │   │
│  │  ┌──────────────────┐  ┌──────────────────┐             │   │
│  │  │  SkillManager    │  │  MemoryManager   │             │   │
│  │  └──────────────────┘  └──────────────────┘             │   │
│  └─────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    MCP Gateway                           │   │
│  │  - 客户端连接池                                          │   │
│  │  - 工具代理与命名空间                                    │   │
│  └─────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Tool Registry                         │   │
│  │  - 内置工具注册                                          │   │
│  │  - MCP 工具动态注册                                      │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────┬───────────────────────────────────┘
                              │
┌─────────────────────────────┴───────────────────────────────────┐
│                   基础设施层 (Infrastructure)                    │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐          │
│  │FileSystem│ │  go-git  │ │ os/exec  │ │  SQLite  │          │
│  │  Client  │ │  Client  │ │  Runner  │ │   Repo   │          │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 核心设计模式

| 模式           | 应用场景                         | 实现方式                                 |
| -------------- | -------------------------------- | ---------------------------------------- |
| **工厂模式**   | Agent 实例创建                   | `AgentFactory.Create(role string) Agent` |
| **策略模式**   | 不同模型的调用策略               | `ModelProvider` 接口 + 具体实现          |
| **观察者模式** | 执行状态推送                     | Wails Events 作为事件总线                |
| **装饰器模式** | 工具调用拦截（日志、权限）       | `ToolMiddleware` 包装 Genkit Tool        |
| **模板方法**   | 流水线阶段执行                   | 基类定义骨架，子类实现具体阶段           |
| **单例模式**   | 全局服务实例（数据库连接、配置） | `sync.Once` 懒加载                       |

### 2.3 依赖注入

使用构造函数注入，避免全局变量污染：

```go
type App struct {
    ctx         context.Context
    pipelineSvc *PipelineService
    configSvc   *ConfigService
    monitorSvc  *MonitorService
}

func NewApp() *App {
    db := database.NewDB()
    skillMgr := skills.NewManager(db)
    memoryMgr := memory.NewManager(db)
    agentReg := agents.NewRegistry(skillMgr, memoryMgr)
    
    return &App{
        pipelineSvc: services.NewPipelineService(agentReg, skillMgr, memoryMgr),
        configSvc:   services.NewConfigService(db),
        monitorSvc:  services.NewMonitorService(db),
    }
}
```


## 3. 模块详细设计

### 3.1 服务层（Wails 绑定）

#### 3.1.1 PipelineService

**职责**：提供流水线控制的对外接口，供前端调用。

**接口定义**：
```go
type PipelineService struct {
    orch   *orchestrator.PipelineOrchestrator
    db     *sql.DB
    eventBus *events.Bus
}

// RunPipeline 启动流水线执行
func (s *PipelineService) RunPipeline(projectName string, requirement string) (*ExecutionResult, error)

// GetStatus 获取执行状态
func (s *PipelineService) GetStatus(executionID string) (*ExecutionStatus, error)

// Cancel 取消执行
func (s *PipelineService) Cancel(executionID string) error

// Resume 从断点恢复
func (s *PipelineService) Resume(executionID string) error

// ListHistory 获取历史记录
func (s *PipelineService) ListHistory(limit int) ([]ExecutionSummary, error)
```

**执行流程**：
1. 参数校验：检查项目名称是否存在，需求非空。
2. 加载项目配置：通过 `ConfigService` 获取 `ProjectConfig`。
3. 创建执行上下文：包含工作目录、取消信号、事件通道。
4. 异步执行：启动 goroutine 调用 `PipelineOrchestrator.Execute()`，不阻塞前端调用。
5. 实时推送：通过事件总线将进度推送到前端。

#### 3.1.2 MonitorService

**职责**：提供实时状态查询和历史数据接口。

```go
type MonitorService struct {
    db *sql.DB
    bus *events.Bus
}

func (m *MonitorService) SubscribeEvents(ctx context.Context) (<-chan ExecutionEvent, error)

func (m *MonitorService) GetAgentActivities(executionID string) ([]AgentActivity, error)

func (m *MonitorService) GetToolCallLogs(executionID string) ([]ToolCallLog, error)

func (m *MonitorService) GetTraceData(executionID string) (*TraceData, error)
```

#### 3.1.3 前端通信数据模型

**DTO 定义**（位于 `pkg/dto` 包）：

```go
type ExecutionResult struct {
    ExecutionID string `json:"executionId"`
    Status      string `json:"status"`
    OutputDir   string `json:"outputDir,omitempty"`
    Error       string `json:"error,omitempty"`
}

type ExecutionStatus struct {
    ExecutionID  string          `json:"executionId"`
    ProjectName  string          `json:"projectName"`
    Status       string          `json:"status"`
    CurrentStage string          `json:"currentStage"`
    Progress     int             `json:"progress"` // 0-100
    Agents       []AgentSnapshot `json:"agents"`
}

type AgentSnapshot struct {
    Name      string `json:"name"`
    Status    string `json:"status"`
    StartTime string `json:"startTime,omitempty"`
    Thought   string `json:"thought,omitempty"`
}
```

### 3.2 引擎层核心模块

#### 3.2.1 PipelineOrchestrator

**职责**：定义并执行主流水线 Flow，协调各 Agent 执行。

**Genkit Flow 定义**（Go 语言风格）：

```go
func (o *PipelineOrchestrator) DefineFlow() *genkit.Flow {
    return genkit.DefineFlow("MainPipeline",
        // 输入类型
        func(ctx context.Context, input PipelineInput) (*PipelineOutput, error) {
            // Step 1: Orchestrator 分解任务
            plan, err := o.runOrchestrator(ctx, input)
            if err != nil {
                return nil, err
            }
            
            // Step 2: Product Manager
            prd, err := o.runPM(ctx, plan)
            if err != nil {
                return nil, err
            }
            
            // Step 3: Architect
            arch, err := o.runArchitect(ctx, prd)
            if err != nil {
                return nil, err
            }
            
            // Step 4: 并行执行 Frontend 和 Backend
            var wg sync.WaitGroup
            errCh := make(chan error, 2)
            
            wg.Add(2)
            go func() {
                defer wg.Done()
                if err := o.runFrontend(ctx, arch); err != nil {
                    errCh <- err
                }
            }()
            go func() {
                defer wg.Done()
                if err := o.runBackend(ctx, arch); err != nil {
                    errCh <- err
                }
            }()
            wg.Wait()
            close(errCh)
            
            // 检查并行错误
            for err := range errCh {
                if err != nil {
                    return nil, err
                }
            }
            
            // Step 5: QA Engineer
            if err := o.runQA(ctx); err != nil {
                return nil, err
            }
            
            // Step 6: Reviewer
            if err := o.runReviewer(ctx); err != nil {
                // 审查失败不阻断流程，仅记录
                log.Warn("review failed", "error", err)
            }
            
            // Step 7: DevOps
            if err := o.runDevOps(ctx); err != nil {
                return nil, err
            }
            
            return &PipelineOutput{Success: true}, nil
        },
    )
}
```

**错误处理与重试策略**：

```go
func (o *PipelineOrchestrator) runWithRetry(ctx context.Context, stage string, fn func() error) error {
    maxRetries := 3
    backoff := 2 * time.Second
    
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if i < maxRetries-1 {
            log.Warn("stage failed, retrying", "stage", stage, "attempt", i+1, "error", err)
            time.Sleep(backoff)
            backoff *= 2
        } else {
            return fmt.Errorf("stage %s failed after %d attempts: %w", stage, maxRetries, err)
        }
    }
    return nil
}
```

#### 3.2.2 Agent 实现模板

所有专业 Agent 实现统一的 `Agent` 接口：

```go
type Agent interface {
    Name() string
    Role() string
    Execute(ctx context.Context, task AgentTask) (*AgentResult, error)
}

type AgentTask struct {
    Description string
    Context     map[string]interface{}
    Files       []string // 相关文件路径
}

type AgentResult struct {
    Output      string
    Artifacts   []Artifact // 生成的文件、文档等
    TokenUsage  int
    ToolCalls   []ToolCallRecord
}
```

**基础实现结构体**（模板方法）：

```go
type BaseAgent struct {
    name        string
    role        string
    model       genkit.Model
    tools       []genkit.Tool
    skillMgr    *skills.Manager
    memoryMgr   *memory.Manager
}

func (a *BaseAgent) Execute(ctx context.Context, task AgentTask) (*AgentResult, error) {
    // 1. 记忆检索与注入
    memCtx := a.memoryMgr.RetrieveForTask(ctx, task.Description, a.role)
    
    // 2. 技能加载（渐进式披露）
    skills := a.skillMgr.LoadForTask(ctx, task.Description)
    
    // 3. 构建 Prompt
    prompt := a.buildPrompt(task, memCtx, skills)
    
    // 4. 调用 Genkit Generate
    response, err := genkit.Generate(ctx, a.model, prompt,
        genkit.WithTools(a.tools...),
        genkit.WithTemperature(0.2),
    )
    if err != nil {
        return nil, err
    }
    
    // 5. 后处理：更新记忆、触发技能生成
    a.postProcess(ctx, task, response)
    
    return &AgentResult{
        Output:     response.Text(),
        TokenUsage: response.Usage.TotalTokens,
        ToolCalls:  extractToolCalls(response),
    }, nil
}
```

**各专业 Agent 实现要点**：

| Agent               | 特殊处理                                                     |
| ------------------- | ------------------------------------------------------------ |
| **Orchestrator**    | 输出结构化任务分解（JSON），通过 `genkit.WithOutputSchema` 约束 |
| **Frontend Dev**    | 工具集中包含 `npm_install`, `npm_run_dev`，偏好使用 React/Vue 模板 |
| **Backend Dev**     | 工具集包含 `go_mod_init`, `sql_schema_gen`，遵循 Go 项目布局 |
| **QA Engineer**     | 生成测试文件后自动执行 `go test` / `npm test`，收集覆盖率    |
| **DevOps Engineer** | 执行 `git_init`, `git_commit`，可选推送到远程                |

#### 3.2.3 工具实现（Tool Registry）

**文件系统工具示例**：

```go
func RegisterFileSystemTools(reg *ToolRegistry, rootDir string) {
    reg.Register(genkit.DefineTool(
        "fs_write",
        "Write content to a file",
        func(ctx context.Context, input struct {
            Path    string `json:"path"`
            Content string `json:"content"`
        }) (string, error) {
            // 路径安全检查
            safePath, err := securePath(rootDir, input.Path)
            if err != nil {
                return "", err
            }
            
            // 确保目录存在
            dir := filepath.Dir(safePath)
            if err := os.MkdirAll(dir, 0755); err != nil {
                return "", fmt.Errorf("create dir: %w", err)
            }
            
            // 写入文件
            if err := os.WriteFile(safePath, []byte(input.Content), 0644); err != nil {
                return "", fmt.Errorf("write file: %w", err)
            }
            
            return fmt.Sprintf("Successfully wrote to %s", input.Path), nil
        },
    ))
}
```

**工具调用拦截器（Middleware）**：

```go
type ToolMiddleware func(next genkit.ToolFunc) genkit.ToolFunc

func LoggingMiddleware(logger *slog.Logger) ToolMiddleware {
    return func(next genkit.ToolFunc) genkit.ToolFunc {
        return func(ctx context.Context, args map[string]any) (any, error) {
            start := time.Now()
            logger.Info("tool call start", "tool", ctx.Value("tool_name"), "args", args)
            
            result, err := next(ctx, args)
            
            logger.Info("tool call end",
                "tool", ctx.Value("tool_name"),
                "duration", time.Since(start),
                "error", err,
            )
            return result, err
        }
    }
}
```

### 3.3 自进化子系统详细设计

#### 3.3.1 SkillManager

**类结构**：

```go
type SkillManager struct {
    db        *sql.DB
    skillDir  string
    validator *SkillValidator
    generator *SkillGenerator
    loader    *SkillLoader
}

// 渐进式披露加载器
type SkillLoader struct {
    cache map[string]*Skill // L1 完整技能缓存
    mu    sync.RWMutex
}

func (l *SkillLoader) LoadL0() ([]SkillMeta, error) {
    // 从 skills_index 表查询元数据
    rows, err := db.Query(`SELECT name, version, description, trigger FROM skills_index WHERE enabled=1`)
    // ...
}

func (l *SkillLoader) LoadL1(name string) (*Skill, error) {
    l.mu.RLock()
    if skill, ok := l.cache[name]; ok {
        l.mu.RUnlock()
        return skill, nil
    }
    l.mu.RUnlock()
    
    // 从文件加载完整内容
    var filePath string
    db.QueryRow(`SELECT file_path FROM skills_index WHERE name=?`, name).Scan(&filePath)
    
    content, err := os.ReadFile(filePath)
    // 解析 YAML frontmatter + Markdown
    skill, err := parseSkill(content)
    
    l.mu.Lock()
    l.cache[name] = skill
    l.mu.Unlock()
    
    return skill, nil
}
```

**技能自动生成算法**：

```go
func (g *SkillGenerator) GenerateFromTrace(trace *ExecutionTrace) (*Skill, error) {
    // 触发条件检查
    if len(trace.ToolCalls) < 5 || trace.Status != "completed" {
        return nil, ErrNotEligible
    }
    
    // 构建生成 Prompt
    prompt := fmt.Sprintf(`
你是一个经验丰富的开发者。以下是一次成功的任务执行轨迹，请从中提取可复用的工作流程，生成一个 Skill 文档。

任务描述：%s
工具调用序列：
%s
最终结果：成功

请按以下格式输出 Skill（YAML frontmatter + Markdown）：
---
name: 技能名称
version: 1
description: 简短描述
trigger: 触发条件关键词（逗号分隔）
tools_required: [工具列表]
---
# 技能名称
## 适用场景
...
## 执行步骤
1. ...
## 常见陷阱
- ...
`, trace.Description, formatToolCalls(trace.ToolCalls))
    
    // 调用轻量模型生成
    response, err := genkit.Generate(ctx, g.model, prompt)
    if err != nil {
        return nil, err
    }
    
    // 解析并验证
    skill, err := parseSkillFromResponse(response.Text())
    if err != nil {
        return nil, err
    }
    
    if err := g.validator.Validate(skill); err != nil {
        return nil, err
    }
    
    return skill, nil
}
```

#### 3.3.2 MemoryManager

**三层记忆实现**：

```go
type MemoryManager struct {
    db          *sql.DB
    workingMem  *WorkingMemory
    semanticDir string
}

// L1 工作记忆（会话级，基于 context）
type WorkingMemory struct {
    data sync.Map
}

func (w *WorkingMemory) Set(key string, value interface{}) {
    w.data.Store(key, value)
}

func (w *WorkingMemory) Get(key string) interface{} {
    v, _ := w.data.Load(key)
    return v
}

// L2 情节记忆（持久化 + 全文检索）
func (m *MemoryManager) StoreEpisodic(mem *EpisodicMemory) error {
    _, err := m.db.Exec(`
        INSERT INTO episodic_memories (id, task_type, description, steps, outcome, agent_role, timestamp)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `, mem.ID, mem.TaskType, mem.Description, mem.Steps, mem.Outcome, mem.AgentRole, time.Now())
    return err
}

func (m *MemoryManager) RetrieveEpisodic(query string, limit int) ([]*EpisodicMemory, error) {
    rows, err := m.db.Query(`
        SELECT id, task_type, description, steps, outcome, agent_role, timestamp
        FROM episodic_memories
        WHERE episodic_memories MATCH ?
        ORDER BY rank
        LIMIT ?
    `, query, limit)
    // 解析 rows...
}

// L3 语义记忆（文件读写）
func (m *MemoryManager) LoadSemantic() (*SemanticMemory, error) {
    data := &SemanticMemory{
        UserProfile:  loadUserProfile(filepath.Join(m.semanticDir, "USER.md")),
        ProjectFacts: loadProjectFacts(filepath.Join(m.semanticDir, "MEMORY.md")),
    }
    return data, nil
}

func (m *MemoryManager) UpdateSemantic(updateFn func(*SemanticMemory)) error {
    mem, _ := m.LoadSemantic()
    updateFn(mem)
    return m.saveSemantic(mem)
}
```

**记忆检索优先级**：

```go
func (m *MemoryManager) RetrieveForTask(ctx context.Context, task, role string) *MemoryContext {
    mc := &MemoryContext{}
    
    // L3 优先（项目约定、用户偏好）
    sem, _ := m.LoadSemantic()
    mc.UserPreferences = sem.UserProfile.Preferences
    mc.ProjectConventions = sem.ProjectFacts
    
    // L2 补充（相似历史经验）
    similar, _ := m.RetrieveEpisodic(task, 3)
    mc.HistoricalExamples = similar
    
    return mc
}
```

### 3.4 MCP Gateway

```go
type MCPGateway struct {
    clients map[string]*mcp.Client
    tools   map[string]genkit.Tool // namespace:toolName -> Tool
    mu      sync.RWMutex
}

func (g *MCPGateway) Connect(id string, config MCPConfig) error {
    var client *mcp.Client
    var err error
    
    switch config.Transport {
    case "stdio":
        client, err = mcp.NewStdioClient(config.Command, config.Args...)
    case "sse":
        client, err = mcp.NewSSEClient(config.URL)
    default:
        return fmt.Errorf("unsupported transport: %s", config.Transport)
    }
    if err != nil {
        return err
    }
    
    g.mu.Lock()
    g.clients[id] = client
    g.mu.Unlock()
    
    // 发现并注册工具
    tools, err := client.ListTools(context.Background())
    if err != nil {
        return err
    }
    
    for _, t := range tools {
        nsName := fmt.Sprintf("%s:%s", id, t.Name)
        genkitTool := convertMCPToolToGenkit(client, t)
        g.tools[nsName] = genkitTool
    }
    
    return nil
}
```


## 4. 核心算法与流程

### 4.1 流水线并行执行算法

使用 goroutine + channel 实现并行阶段，错误收集后统一处理：

```go
type parallelResult struct {
    Stage string
    Err   error
}

func runParallelStages(ctx context.Context, stages []StageFunc) error {
    results := make(chan parallelResult, len(stages))
    var wg sync.WaitGroup
    
    for _, stage := range stages {
        wg.Add(1)
        go func(s StageFunc) {
            defer wg.Done()
            err := s(ctx)
            results <- parallelResult{Stage: s.Name(), Err: err}
        }(stage)
    }
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    var firstErr error
    for res := range results {
        if res.Err != nil && firstErr == nil {
            firstErr = res.Err
        }
    }
    return firstErr
}
```

### 4.2 上下文传递机制

Agent 间通过文件系统和共享数据结构传递上下文：

1. **文件传递**：PRD 写入 `docs/PRD.md`，架构设计写入 `docs/ARCHITECTURE.md`。
2. **结构化数据**：通过 Context 传递 `ProjectContext` 结构体：

```go
type ProjectContext struct {
    ProjectName   string
    RootDir       string
    TechStack     TechStack
    PRDPath       string
    ArchPath      string
    GeneratedFiles []string
}
```

### 4.3 工具调用安全校验

**路径穿越防护**：

```go
func securePath(root, sub string) (string, error) {
    dest := filepath.Join(root, sub)
    absRoot, _ := filepath.Abs(root)
    absDest, _ := filepath.Abs(dest)
    
    if !strings.HasPrefix(absDest, absRoot) {
        return "", fmt.Errorf("path traversal detected: %s", sub)
    }
    return absDest, nil
}
```

**命令注入防护**：

```go
var dangerousPatterns = []*regexp.Regexp{
    regexp.MustCompile(`rm\s+-rf\s+/`),
    regexp.MustCompile(`>\s*/dev/`),
    regexp.MustCompile(`\|.*sh`),
}

func validateCommand(cmd string) error {
    for _, p := range dangerousPatterns {
        if p.MatchString(cmd) {
            return fmt.Errorf("dangerous command pattern detected")
        }
    }
    return nil
}
```


## 5. 状态机设计

### 5.1 流水线执行状态机

```
                  ┌─────────┐
                  │  IDLE   │
                  └────┬────┘
                       │ RunPipeline
                       ▼
                  ┌─────────┐
         ┌────────│ RUNNING │────────┐
         │        └────┬────┘        │
         │             │             │
         │       ┌─────┴─────┐       │
         │       ▼           ▼       │
         │  ┌─────────┐ ┌─────────┐  │
         │  │PAUSED   │ │CANCELLED│  │
         │  └────┬────┘ └─────────┘  │
         │       │ Resume            │
         │       ▼                   │
         │  ┌─────────┐              │
         └─►│COMPLETED│◄─────────────┘
            └─────────┘
                   │ Error
                   ▼
            ┌─────────┐
            │ FAILED  │
            └─────────┘
```

**状态迁移表**：

| 当前状态 | 事件        | 下一状态  | 动作                             |
| -------- | ----------- | --------- | -------------------------------- |
| IDLE     | RunPipeline | RUNNING   | 初始化执行上下文，启动 goroutine |
| RUNNING  | Complete    | COMPLETED | 更新数据库，触发后处理           |
| RUNNING  | Error       | FAILED    | 记录错误，清理资源               |
| RUNNING  | Cancel      | CANCELLED | 发送取消信号，等待退出           |
| RUNNING  | Pause       | PAUSED    | 保存断点，暂停调度               |
| PAUSED   | Resume      | RUNNING   | 恢复执行                         |
| FAILED   | Retry       | RUNNING   | 重置错误计数，重新执行           |

### 5.2 Agent 活动状态机

```
   ┌─────────┐
   │ PENDING │
   └────┬────┘
        │ Start
        ▼
   ┌─────────┐
   │ RUNNING │
   └────┬────┘
        │
   ┌────┴────┐
   ▼         ▼
┌─────────┐ ┌─────────┐
│COMPLETED│ │ FAILED  │
└─────────┘ └─────────┘
```


## 6. 异常处理设计

### 6.1 错误分类

```go
type ErrorCode int

const (
    ErrConfigInvalid ErrorCode = iota + 1
    ErrModelUnavailable
    ErrToolExecution
    ErrAgentTimeout
    ErrPathNotAllowed
    ErrCommandNotAllowed
)

type AppError struct {
    Code    ErrorCode
    Message string
    Cause   error
}

func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}
```

### 6.2 错误处理策略

| 错误类型     | 处理策略               | 用户提示                              |
| ------------ | ---------------------- | ------------------------------------- |
| 配置无效     | 终止执行，返回配置页面 | "项目配置无效：{详情}"                |
| 模型不可用   | 自动降级到备用模型     | 静默降级，日志记录                    |
| 工具执行失败 | 重试 3 次，指数退避    | "工具 {name} 执行失败，正在重试..."   |
| Agent 超时   | 取消当前阶段，标记失败 | "Agent {name} 执行超时，流水线已停止" |
| 路径越界     | 拒绝操作，记录审计日志 | "操作被拒绝：试图访问项目外路径"      |

### 6.3 优雅关闭

```go
func (o *PipelineOrchestrator) Execute(ctx context.Context, input PipelineInput) error {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    
    // 监听系统信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-sigCh
        log.Info("received shutdown signal, cancelling pipeline")
        cancel()
    }()
    
    // 执行流水线...
    return o.run(ctx, input)
}
```


## 7. 安全详细设计

### 7.1 API Key 管理

```go
type KeyringService struct {
    serviceName string
}

func (k *KeyringService) Set(key, value string) error {
    return keyring.Set(k.serviceName, key, value)
}

func (k *KeyringService) Get(key string) (string, error) {
    return keyring.Get(k.serviceName, key)
}

func (k *KeyringService) Delete(key string) error {
    return keyring.Delete(k.serviceName, key)
}
```

### 7.2 操作审计日志

所有敏感操作记录到结构化日志：

```go
type AuditLog struct {
    Timestamp   time.Time `json:"timestamp"`
    Operation   string    `json:"operation"`   // fs_write, git_commit, shell_exec
    User        string    `json:"user"`
    Project     string    `json:"project"`
    Details     string    `json:"details"`
    Success     bool      `json:"success"`
    ClientIP    string    `json:"client_ip,omitempty"`
}

func (l *AuditLogger) Log(entry AuditLog) {
    // 写入独立审计日志文件
    l.logger.Info("audit", slog.Any("entry", entry))
}
```


## 8. 接口详细设计

### 8.1 前端调用 Go 方法列表

| 方法签名                                                     | 描述         |
| ------------------------------------------------------------ | ------------ |
| `RunPipeline(projectName, requirement string) (ExecutionResult, error)` | 启动流水线   |
| `GetExecutionStatus(executionID string) (ExecutionStatus, error)` | 获取状态     |
| `CancelExecution(executionID string) error`                  | 取消执行     |
| `ListProjects() ([]ProjectSummary, error)`                   | 列出所有项目 |
| `SaveProject(config ProjectConfig) error`                    | 保存项目配置 |
| `GetSkills() ([]SkillMeta, error)`                           | 获取技能列表 |
| `DeleteSkill(name string) error`                             | 删除技能     |
| `SearchMemories(query string) ([]EpisodicMemory, error)`     | 搜索记忆     |
| `GetAppSettings() (AppSettings, error)`                      | 获取应用设置 |
| `SaveAppSettings(settings AppSettings) error`                | 保存应用设置 |

### 8.2 事件推送定义

| 事件名                   | 数据结构                                       |
| ------------------------ | ---------------------------------------------- |
| `pipeline:started`       | `{executionID, projectName, timestamp}`        |
| `pipeline:stage_changed` | `{executionID, stage, status}`                 |
| `pipeline:completed`     | `{executionID, outputDir, duration}`           |
| `pipeline:failed`        | `{executionID, error}`                         |
| `agent:thought`          | `{executionID, agentName, thought}`            |
| `agent:tool_call`        | `{executionID, agentName, tool, args, result}` |
| `skill:generated`        | `{skillName, version}`                         |


## 9. 数据结构与序列化

### 9.1 核心结构体定义

位于 `pkg/models` 包：

```go
// 项目配置
type ProjectConfig struct {
    Name        string           `json:"name"`
    Path        string           `json:"path"`
    Type        string           `json:"type"`
    TechStack   TechStack        `json:"techStack"`
    Agents      []AgentConfig    `json:"agents"`
    Models      ModelConfig      `json:"models"`
    MCPServers  []MCPServerConfig `json:"mcpServers"`
}

// 技能元数据
type SkillMeta struct {
    Name        string   `json:"name"`
    Version     int      `json:"version"`
    Description string   `json:"description"`
    Trigger     []string `json:"trigger"`
    Enabled     bool     `json:"enabled"`
}

// 情节记忆
type EpisodicMemory struct {
    ID          string    `json:"id"`
    TaskType    string    `json:"taskType"`
    Description string    `json:"description"`
    Steps       string    `json:"steps"`
    Outcome     string    `json:"outcome"`
    AgentRole   string    `json:"agentRole"`
    Timestamp   time.Time `json:"timestamp"`
}
```


## 10. 部署与配置

### 10.1 配置文件格式

`~/.ai-dev-pipeline/config.yaml`：

```yaml
app:
  theme: dark
  language: zh-CN
  log_level: info
  max_history_days: 90

models:
  default: gemini-2.5-flash
  providers:
    google:
      api_key_id: gemini_key
    openai:
      api_key_id: openai_key

mcp_servers:
  - name: filesystem
    command: npx
    args: ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed/dir"]
    enabled: true
```

### 10.2 环境变量

| 变量名                 | 说明                    | 默认值               |
| ---------------------- | ----------------------- | -------------------- |
| `AI_DEV_PIPELINE_HOME` | 应用数据目录            | `~/.ai-dev-pipeline` |
| `GENKIT_ENV`           | Genkit 环境（dev/prod） | `prod`               |
| `LOG_LEVEL`            | 日志级别                | `info`               |


## 11. 附录

### 11.1 关键代码片段索引

| 模块       | 文件位置                            | 说明                    |
| ---------- | ----------------------------------- | ----------------------- |
| Wails 入口 | `main.go`                           | 应用启动                |
| 服务层     | `internal/services/*.go`            | 前端调用接口实现        |
| 流水线编排 | `internal/orchestrator/pipeline.go` | 主 Flow 定义            |
| Agent 实现 | `internal/agents/*.go`              | 各专业 Agent            |
| 技能管理   | `internal/skills/manager.go`        | SkillManager 实现       |
| 记忆管理   | `internal/memory/manager.go`        | MemoryManager 实现      |
| 工具注册   | `internal/tools/registry.go`        | ToolRegistry 与内置工具 |

### 11.2 单元测试覆盖要求

- 所有公共接口至少一个正向测试和一个异常测试。
- 工具函数（如路径校验、命令验证）覆盖率 ≥ 90%。
- 使用 `go test -race` 检测数据竞争。

---

**文档版本**：V1.0
**创建日期**：2026-04-18
**最后更新**：2026-04-18