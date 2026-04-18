# AI 软件项目开发流水线（桌面版）概要设计说明书

基于《项目需求说明书 V4.0（Go 语言版）》


## 1. 引言

### 1.1 编写目的
本文档旨在对 AI 软件项目开发流水线（桌面版）进行概要设计，明确系统的总体架构、模块划分、关键接口、数据结构和核心流程，为后续详细设计和编码实现提供指导。

### 1.2 适用范围
本文档适用于项目开发团队、测试团队和技术管理人员，作为技术评审、开发实施和系统集成的依据。

### 1.3 参考文档
- 《项目需求说明书 V4.0（Go 语言版）》
- Genkit Go SDK 官方文档（https://firebase.google.com/docs/genkit-go）
- Wails v2 官方文档（https://wails.io）
- go-git 文档（https://github.com/go-git/go-git）

### 1.4 设计原则
- **模块化**：各功能模块高内聚、低耦合，便于扩展和维护。
- **可观测性**：全链路追踪与监控，执行过程透明化。
- **安全第一**：文件操作、命令执行严格权限控制。
- **自进化**：Agent 可自动沉淀经验为 Skill，实现持续能力提升。
- **跨平台**：基于 Wails + Go 实现 Windows/macOS/Linux 一致体验。


## 2. 总体设计

### 2.1 系统总体架构

系统采用分层架构，自上而下分为**表现层**、**服务层**、**引擎层**和**基础设施层**。

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           表现层 (Presentation)                         │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │
│  │ 项目管理 UI │ │ Agent监控台 │ │ 技能管理 UI │ │ 配置管理 UI │       │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘       │
│                     React + TypeScript + Wails Runtime                   │
├─────────────────────────────────────────────────────────────────────────┤
│                           服务层 (Service)                               │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     Wails App Service (Go)                       │   │
│  │  - PipelineService (流水线控制)   - ConfigService (配置管理)      │   │
│  │  - MonitorService (状态查询)      - EventService (事件推送)       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────────────────┤
│                           引擎层 (Engine)                                │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                      Genkit Flow Orchestrator                    │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐           │   │
│  │  │   Agent  │ │   Flow   │ │   Tool   │ │   MCP    │           │   │
│  │  │  Registry│ │  Engine  │ │ Registry │ │  Gateway │           │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘           │   │
│  ├─────────────────────────────────────────────────────────────────┤   │
│  │                      自进化子系统 (Evolution)                     │   │
│  │  ┌─────────────────────┐  ┌─────────────────────────────────┐  │   │
│  │  │   Skills 闭环引擎   │  │        三层记忆架构              │  │   │
│  │  │  (提取/存储/加载)    │  │  (工作记忆/情节记忆/语义记忆)    │  │   │
│  │  └─────────────────────┘  └─────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────────────────┤
│                         基础设施层 (Infrastructure)                      │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐     │
│  │ 文件系统 │ │  go-git  │ │ os/exec  │ │  SQLite  │ │  Keyring │     │
│  │  (os)    │ │          │ │          │ │          │ │          │     │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘     │
├─────────────────────────────────────────────────────────────────────────┤
│                           外部依赖 (External)                            │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐     │
│  │  Gemini  │ │  OpenAI  │ │ Anthropic│ │  Ollama  │ │MCP Server│     │
│  │   API    │ │   API    │ │   API    │ │  (本地)  │ │ (外部)   │     │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘     │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.2 核心设计决策

| 决策点      | 选型                | 理由                                                        |
| ----------- | ------------------- | ----------------------------------------------------------- |
| 桌面框架    | Wails               | 轻量（10-20MB）、启动快（<80ms）、Go 原生支持               |
| AI 编排框架 | Genkit Go SDK       | 类型安全的 Flow 定义、原生 OpenTelemetry 集成、MCP 插件支持 |
| Git 操作库  | go-git              | 纯 Go 实现，无需系统 Git 依赖，跨平台一致                   |
| 记忆存储    | SQLite + FTS5       | 嵌入式、全文检索支持、无需独立数据库服务                    |
| 前端技术栈  | React + TypeScript  | 生态成熟、组件丰富、类型安全                                |
| 进程通信    | Wails 绑定 + Events | 原生 Go ↔ JS 互操作，支持事件推送                           |
| 模型接入    | Genkit 插件体系     | 统一 API，支持多模型无缝切换                                |

### 2.3 技术架构图（组件视图）

```
┌─────────────────────────────────────────────────────────────────┐
│                       Wails Application                         │
│  ┌─────────────────────┐      ┌─────────────────────────────┐  │
│  │   前端 (React)       │◄────┤    Go 后端 (main.go)        │  │
│  │  - Dashboard         │ IPC  │  - App 结构体               │  │
│  │  - Config Panel      │      │  - 导出方法                 │  │
│  │  - Monitor View      │      │  - 事件发射                 │  │
│  └─────────────────────┘      └──────────┬──────────────────┘  │
│                                          │                      │
│                           ┌──────────────┴──────────────┐      │
│                           │      服务层 (Service)        │      │
│                           │  ┌───────────────────────┐  │      │
│                           │  │   PipelineOrchestrator│  │      │
│                           │  └───────────────────────┘  │      │
│                           └──────────────┬──────────────┘      │
│                                          │                      │
│       ┌──────────────────────────────────┼──────────────────┐  │
│       │                    Genkit 引擎层                     │  │
│       │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │  │
│       │  │AgentManager │  │ FlowEngine  │  │ ToolRegistry│  │  │
│       │  └─────────────┘  └─────────────┘  └─────────────┘  │  │
│       │  ┌─────────────────────────────────────────────────┐│  │
│       │  │              自进化子系统                        ││  │
│       │  │  ┌───────────────┐    ┌────────────────────┐   ││  │
│       │  │  │ SkillManager  │    │  MemoryManager     │   ││  │
│       │  │  └───────────────┘    └────────────────────┘   ││  │
│       │  └─────────────────────────────────────────────────┘│  │
│       └──────────────────────────────────────────────────────┘  │
│                                          │                      │
│                           ┌──────────────┴──────────────┐      │
│                           │       基础设施工具           │      │
│                           │  ┌───────┐ ┌───────┐ ┌────┐ │      │
│                           │  │FileSys│ │ go-git│ │Exec│ │      │
│                           │  └───────┘ └───────┘ └────┘ │      │
│                           └──────────────────────────────┘      │
└─────────────────────────────────────────────────────────────────┘
```


## 3. 模块设计

### 3.1 模块划分

系统划分为以下核心模块：

| 模块名称           | 职责                                 | 主要接口                                      |
| ------------------ | ------------------------------------ | --------------------------------------------- |
| **AppService**     | Wails 服务层，提供前端调用的 Go 方法 | `RunPipeline`, `GetStatus`, `SubscribeEvents` |
| **AgentManager**   | Agent 注册、生命周期管理、配置加载   | `Register`, `Get`, `List`                     |
| **FlowEngine**     | 流水线执行引擎，支持顺序/并行        | `Execute`, `Resume`, `Cancel`                 |
| **ToolRegistry**   | 工具注册、调用拦截、权限检查         | `Register`, `Call`                            |
| **SkillManager**   | 技能库管理、自动生成、渐进式加载     | `Load`, `Generate`, `Search`                  |
| **MemoryManager**  | 三层记忆读写、检索、更新             | `Store`, `Retrieve`, `Update`                 |
| **MCPGateway**     | MCP 客户端/服务端管理                | `Connect`, `ListTools`, `CallTool`            |
| **MonitorService** | 执行追踪、指标采集、状态推送         | `Trace`, `Metrics`, `Emit`                    |
| **ConfigService**  | 项目配置、模型配置持久化             | `Load`, `Save`, `Validate`                    |

### 3.2 模块间调用关系

```
┌──────────────┐
│ AppService   │ (Wails 导出)
└──────┬───────┘
       │ 调用
       ▼
┌──────────────┐      ┌──────────────┐
│ FlowEngine   │─────►│ AgentManager │
└──────┬───────┘      └──────┬───────┘
       │                     │
       │ 调用                │ 调用
       ▼                     ▼
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│ ToolRegistry │◄─────│ SkillManager │◄─────│ MemoryManager│
└──────────────┘      └──────────────┘      └──────────────┘
       │                     │                     │
       │                     │                     │
       └─────────┬───────────┴──────────┬─────────┘
                 │                      │
                 ▼                      ▼
         ┌──────────────┐       ┌──────────────┐
         │ MCPGateway   │       │MonitorService│
         └──────────────┘       └──────────────┘
```

### 3.3 关键模块详细设计

#### 3.3.1 FlowEngine（流水线引擎）

**职责**：执行 Genkit Flow，管理流水线生命周期。

**核心结构**：
```go
type FlowEngine struct {
    flows       map[string]*genkit.Flow
    executor    *Executor
    cancelFuncs map[string]context.CancelFunc
}

type PipelineRequest struct {
    ProjectConfig ProjectConfig
    UserRequirement string
    WorkDir      string
}

type PipelineResult struct {
    Success      bool
    OutputDir    string
    ExecutionID  string
    Error        error
}
```

**主要方法**：
- `Execute(ctx context.Context, req PipelineRequest) (*PipelineResult, error)`
- `Resume(executionID string) error`
- `Cancel(executionID string) error`
- `GetStatus(executionID string) (*ExecutionStatus, error)`

**执行流程**（详见 4.1 节）。

#### 3.3.2 AgentManager

**职责**：管理所有 Agent 实例，支持动态注册与配置加载。

**Agent 接口定义**：
```go
type Agent interface {
    Name() string
    Role() string
    SystemPrompt() string
    Model() genkit.Model
    Tools() []genkit.Tool
    Execute(ctx context.Context, input string) (string, error)
}
```

**内置 Agent 实现**（以 FrontendDev 为例）：
```go
type FrontendDevAgent struct {
    name   string
    model  genkit.Model
    tools  []genkit.Tool
    skills *SkillManager
    memory *MemoryManager
}

func (a *FrontendDevAgent) Execute(ctx context.Context, task string) (string, error) {
    // 1. 记忆检索注入
    memories := a.memory.Retrieve(ctx, task)
    // 2. 技能加载
    skills := a.skills.LoadForTask(task)
    // 3. 构建提示词
    prompt := a.buildPrompt(task, memories, skills)
    // 4. 调用 Genkit Generate
    return genkit.Generate(ctx, a.model, prompt)
}
```

#### 3.3.3 SkillManager（技能管理）

**职责**：管理技能库，实现自动技能生成和渐进式加载。

**Skill 数据结构**：
```go
type Skill struct {
    Name        string    `yaml:"name"`
    Version     int       `yaml:"version"`
    Description string    `yaml:"description"`
    Trigger     string    `yaml:"trigger"`
    Tools       []string  `yaml:"tools_required"`
    Content     string    // Markdown 正文
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**核心方法**：
- `LoadL0() []SkillMeta`：仅加载元数据（名称、触发词）
- `LoadL1(name string) (*Skill, error)`：加载完整 Skill
- `GenerateFromTrace(trace *ExecutionTrace) (*Skill, error)`：从执行轨迹生成 Skill
- `Search(query string) []*Skill`：全文检索 Skill

**技能生成触发条件**：
- 工具调用次数 ≥ 5
- 执行成功且无人工干预
- 任务类型可泛化（非一次性特定任务）

#### 3.3.4 MemoryManager（三层记忆）

**职责**：实现工作记忆、情节记忆、语义记忆的分层存储与检索。

**数据结构**：
```go
// L1: 工作记忆（内存）
type WorkingMemory struct {
    SessionID   string
    Context     map[string]interface{}
    History     []Message
}

// L2: 情节记忆（SQLite）
type EpisodicMemory struct {
    ID          string
    TaskType    string
    Description string
    Steps       []Step
    Outcome     string
    Timestamp   time.Time
    Embedding   []float32  // 可选，用于语义检索
}

// L3: 语义记忆（文件）
type SemanticMemory struct {
    UserProfile   UserProfile   // USER.md
    ProjectFacts  []Fact        // MEMORY.md
}
```

**检索流程**：
1. 优先检索 L3 语义记忆（项目约定、用户偏好）
2. 检索 L2 情节记忆（相似任务历史）
3. 构建上下文注入 Agent Prompt

#### 3.3.5 MCPGateway

**职责**：管理 MCP 客户端连接，实现工具发现与调用。

**接口设计**：
```go
type MCPGateway struct {
    clients map[string]*mcp.Client
    tools   map[string]genkit.Tool  // 命名空间化工具
}

// 连接 MCP Server
func (g *MCPGateway) Connect(id, transport, address string) error

// 获取所有 MCP 工具（已转换为 Genkit Tool）
func (g *MCPGateway) GetTools() []genkit.Tool

// 将本地工具暴露为 MCP Server
func (g *MCPGateway) StartServer(port int) error
```

#### 3.3.6 MonitorService（监控服务）

**职责**：采集执行追踪数据，推送到前端。

**事件类型**：
```go
type ExecutionEvent struct {
    Type        EventType  // AgentStart, ToolCall, Thought, Error, SkillGenerated
    ExecutionID string
    AgentName   string
    Timestamp   time.Time
    Data        map[string]interface{}
}
```

**推送方式**：通过 Wails `Runtime.EventsEmit` 向前端推送事件，前端通过 `EventsOn` 订阅。


## 4. 核心流程设计

### 4.1 主流水线执行流程

```
用户输入需求
      │
      ▼
┌─────────────────────────────────────────────────────────────┐
│                    FlowEngine.Execute()                      │
│  1. 创建工作目录，初始化项目结构                               │
│  2. 启动 OpenTelemetry Span                                   │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Stage 1: Orchestrator                     │
│  - 分析需求，分解任务                                          │
│  - 生成执行计划（顺序/并行步骤）                                │
│  - 输出任务分解文档                                            │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Stage 2: Product Manager                  │
│  - 详细需求分析                                                │
│  - 生成 PRD.md                                                │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Stage 3: Architect                        │
│  - 技术架构设计                                                │
│  - 生成 ARCHITECTURE.md                                       │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                Stage 4: 并行开发 (goroutine)                 │
│  ┌───────────────────┐        ┌───────────────────┐         │
│  │   Frontend Dev    │        │    Backend Dev    │         │
│  │  - 生成前端代码    │        │  - 生成后端 API   │         │
│  │  - 组件/页面/样式  │        │  - 数据库 Schema  │         │
│  └───────────────────┘        └───────────────────┘         │
│           │                            │                     │
│           └──────────┬─────────────────┘                     │
│                      ▼                                       │
│              结果合并（通过 Channel）                          │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Stage 5: QA Engineer                      │
│  - 生成单元测试/集成测试                                       │
│  - 执行测试（可选）                                            │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Stage 6: Reviewer                         │
│  - 代码审查                                                    │
│  - 输出审查报告                                                │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Stage 7: DevOps Engineer                  │
│  - 安装依赖                                                    │
│  - 构建项目                                                    │
│  - 启动验证（可选）                                            │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    后处理阶段                                 │
│  - 触发 Skill 生成（如满足条件）                               │
│  - 更新情节记忆                                                │
│  - 完成 Trace，导出遥测数据                                    │
│  - 返回执行结果                                                │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 自进化流程（Skills 闭环）

```
┌─────────────────────────────────────────────────────────────────┐
│                        执行任务                                  │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│              记录执行轨迹 (ExecutionTrace)                        │
│  - 工具调用序列                                                  │
│  - 输入输出                                                      │
│  - 错误与修复                                                    │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
                   ┌─────────────────────┐
                   │  满足触发条件？      │
                   │  - 工具调用 ≥5      │
                   │  - 执行成功         │
                   │  - 任务可泛化       │
                   └─────────┬───────────┘
                             │ 是
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                SkillGenerator.GenerateFromTrace()                │
│  1. 调用 LLM 提取关键步骤                                        │
│  2. 识别常见陷阱和边界条件                                        │
│  3. 生成 YAML frontmatter + Markdown 正文                        │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                SkillValidator.Validate()                         │
│  - 格式校验                                                      │
│  - 安全扫描（禁止危险命令）                                       │
│  - 去重检查                                                      │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                SkillManager.Save()                               │
│  - 存储到 ~/.hermes/skills/                                      │
│  - 更新索引                                                      │
│  - 触发 SkillGenerated 事件                                      │
└─────────────────────────────────────────────────────────────────┘
```

### 4.3 记忆检索流程

```
Agent 接收任务
      │
      ▼
┌─────────────────────────────────────────────────────────────┐
│                    MemoryManager.Retrieve()                   │
│  输入：任务描述、Agent 角色                                    │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│              Step 1: L3 语义记忆检索                          │
│  - 读取 USER.md（用户偏好）                                   │
│  - 读取项目 MEMORY.md（项目约定）                              │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│              Step 2: L2 情节记忆检索                          │
│  - SQLite FTS5 全文检索                                       │
│  - 按任务类型、关键词匹配                                      │
│  - 返回 Top-K 相似历史任务                                     │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│              Step 3: 结果整合                                 │
│  - 构建记忆上下文文本                                          │
│  - 注入 Agent 系统提示词                                       │
└─────────────────────────────────────────────────────────────┘
```


## 5. 接口设计

### 5.1 前端 ↔ Go 后端接口（Wails 绑定）

| 方法签名                                                     | 描述            | 请求参数           | 返回值         |
| ------------------------------------------------------------ | --------------- | ------------------ | -------------- |
| `RunPipeline(config ProjectConfig, requirement string) (ExecutionResult, error)` | 启动流水线      | 项目配置、需求文本 | 执行结果       |
| `GetExecutionStatus(executionID string) (ExecutionStatus, error)` | 获取执行状态    | 执行 ID            | 状态对象       |
| `CancelExecution(executionID string) error`                  | 取消执行        | 执行 ID            | 错误信息       |
| `ListAgents() ([]AgentInfo, error)`                          | 获取 Agent 列表 | -                  | Agent 信息数组 |
| `GetSkills() ([]SkillMeta, error)`                           | 获取技能列表    | -                  | 技能元数据数组 |
| `GetSkillDetail(name string) (Skill, error)`                 | 获取技能详情    | 技能名             | 完整 Skill     |
| `DeleteSkill(name string) error`                             | 删除技能        | 技能名             | 错误信息       |
| `GetMemories(query string) ([]EpisodicMemory, error)`        | 检索记忆        | 查询关键词         | 记忆数组       |
| `SaveConfig(config ProjectConfig) error`                     | 保存项目配置    | 配置对象           | 错误信息       |
| `LoadConfig(projectName string) (ProjectConfig, error)`      | 加载项目配置    | 项目名             | 配置对象       |

### 5.2 事件推送接口（Wails Events）

前端通过 `EventsOn` 订阅以下事件：

| 事件名              | 数据格式                                       | 触发时机           |
| ------------------- | ---------------------------------------------- | ------------------ |
| `agent:start`       | `{executionID, agentName, task}`               | Agent 开始执行     |
| `agent:thought`     | `{executionID, agentName, thought}`            | Agent 输出推理过程 |
| `agent:tool_call`   | `{executionID, agentName, tool, args, result}` | 工具调用完成       |
| `agent:complete`    | `{executionID, agentName, output}`             | Agent 执行完成     |
| `agent:error`       | `{executionID, agentName, error}`              | Agent 执行出错     |
| `flow:stage_change` | `{executionID, stage, status}`                 | 流水线阶段变更     |
| `skill:generated`   | `{skillName, triggerTrace}`                    | 新 Skill 生成      |
| `memory:updated`    | `{memoryType, key}`                            | 记忆更新           |

### 5.3 内部模块接口

#### 5.3.1 ToolRegistry 接口

```go
type ToolRegistry interface {
    Register(tool genkit.Tool) error
    RegisterFromMCP(namespace string, mcpTool mcp.Tool) error
    Get(name string) (genkit.Tool, error)
    List() []genkit.Tool
    Call(ctx context.Context, name string, args map[string]any) (any, error)
}
```

#### 5.3.2 SkillManager 接口

```go
type SkillManager interface {
    LoadL0() []SkillMeta
    LoadL1(name string) (*Skill, error)
    Save(skill *Skill) error
    Delete(name string) error
    Search(query string) []*Skill
    GenerateFromTrace(trace *ExecutionTrace) (*Skill, error)
}
```

#### 5.3.3 MemoryManager 接口

```go
type MemoryManager interface {
    // L1 工作记忆
    SetWorking(key string, value interface{})
    GetWorking(key string) interface{}
    ClearWorking()
    
    // L2 情节记忆
    StoreEpisodic(memory *EpisodicMemory) error
    RetrieveEpisodic(query string, limit int) ([]*EpisodicMemory, error)
    
    // L3 语义记忆
    LoadSemantic() (*SemanticMemory, error)
    UpdateSemantic(updateFunc func(*SemanticMemory)) error
    
    // 组合检索
    RetrieveForTask(task string, agentRole string) (*MemoryContext, error)
}
```


## 6. 数据设计

### 6.1 数据库设计（SQLite）

#### 6.1.1 执行轨迹表（execution_traces）
```sql
CREATE TABLE execution_traces (
    id TEXT PRIMARY KEY,
    project_name TEXT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    status TEXT NOT NULL,  -- running, completed, failed, cancelled
    total_tokens INTEGER DEFAULT 0,
    total_cost REAL DEFAULT 0.0,
    trace_data JSON,       -- OpenTelemetry Span 数据
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### 6.1.2 情节记忆表（episodic_memories）
```sql
CREATE VIRTUAL TABLE episodic_memories USING fts5(
    id UNINDEXED,
    task_type,
    description,
    steps,
    outcome,
    agent_role,
    timestamp,
    embedding
);
```

#### 6.1.3 技能索引表（skills_index）
```sql
CREATE TABLE skills_index (
    name TEXT PRIMARY KEY,
    version INTEGER NOT NULL,
    description TEXT,
    trigger TEXT,
    tools_required TEXT,  -- JSON 数组
    file_path TEXT NOT NULL,
    enabled BOOLEAN DEFAULT 1,
    usage_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### 6.1.4 项目配置表（project_configs）
```sql
CREATE TABLE project_configs (
    name TEXT PRIMARY KEY,
    path TEXT NOT NULL,
    type TEXT,
    tech_stack TEXT,      -- JSON
    agent_configs TEXT,   -- JSON
    model_configs TEXT,   -- JSON
    mcp_servers TEXT,     -- JSON
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 6.2 文件存储结构

```
~/.ai-dev-pipeline/
├── config.yaml                 # 全局配置
├── data/
│   └── pipeline.db             # SQLite 数据库
├── skills/                     # 技能库
│   ├── react-component.md
│   ├── api-endpoint.md
│   └── ...
├── memory/
│   ├── USER.md                 # 用户画像
│   └── projects/               # 项目级记忆
│       └── {project_hash}/
│           └── MEMORY.md
├── traces/                     # 执行轨迹导出（可选）
│   └── {execution_id}.json
└── logs/
    └── app.log
```

### 6.3 配置数据结构

#### 6.3.1 项目配置（ProjectConfig）
```go
type ProjectConfig struct {
    Name        string            `json:"name"`
    Path        string            `json:"path"`
    Type        string            `json:"type"`        // web, cli, api
    TechStack   TechStack         `json:"techStack"`
    Agents      []AgentConfig     `json:"agents"`
    Models      ModelConfig       `json:"models"`
    MCPServers  []MCPServerConfig `json:"mcpServers"`
}

type AgentConfig struct {
    Role    string `json:"role"`     // orchestrator, frontend, backend...
    Enabled bool   `json:"enabled"`
    Model   string `json:"model"`    // gemini-2.5-pro, gemini-2.5-flash...
}

type ModelConfig struct {
    Provider string            `json:"provider"` // google, openai, anthropic
    APIKey   string            `json:"-"`        // 从 keyring 读取
    Options  map[string]any    `json:"options"`  // temperature, maxTokens...
}
```


## 7. 安全设计

### 7.1 API Key 存储
- 使用操作系统原生密钥链：
  - macOS: Keychain (`github.com/keybase/go-keychain`)
  - Windows: Credential Manager (`github.com/danieljoos/wincred`)
  - Linux: Secret Service (`github.com/zalando/go-keyring`)
- 在配置文件中仅存储 Key ID 引用，不存储明文。

### 7.2 工作区隔离
- 所有文件操作通过 `SecureFileSystem` 包装器执行。
- 操作前校验路径前缀是否为项目根目录：
```go
func (fs *SecureFileSystem) validatePath(path string) (string, error) {
    absPath, err := filepath.Abs(path)
    if err != nil {
        return "", err
    }
    if !strings.HasPrefix(absPath, fs.rootDir) {
        return "", ErrPathNotAllowed
    }
    return absPath, nil
}
```

### 7.3 命令执行沙箱
- 维护可执行命令白名单：
```go
var allowedCommands = map[string]bool{
    "npm":   true,
    "yarn":  true,
    "pnpm":  true,
    "go":    true,
    "git":   true,
    "node":  true,
}
```
- 危险参数模式检测（如 `rm -rf /`、`format`）。
- 所有命令执行记录审计日志。

### 7.4 Skill 安全扫描
- 自动生成的 Skill 需通过安全检查：
  - 不包含 `rm -rf` 等危险命令
  - 不包含文件删除/移动敏感路径的指令
  - 不包含网络请求到未知域名


## 8. 部署与运维设计

### 8.1 构建流程
```bash
# 开发模式
wails dev

# 生产构建
wails build -platform windows/amd64 -o ai-dev-pipeline.exe
wails build -platform darwin/amd64 -o ai-dev-pipeline.app
wails build -platform darwin/arm64 -o ai-dev-pipeline.app
wails build -platform linux/amd64 -o ai-dev-pipeline
```

### 8.2 自动更新
- 集成 Wails 自动更新方案：
  - 在 GitHub Releases 托管更新文件
  - 客户端定期检查 `latest.yml` 版本信息
  - 下载增量更新包并替换

### 8.3 日志管理
- 使用 `log/slog` 结构化日志，输出到文件和控制台。
- 日志级别：DEBUG、INFO、WARN、ERROR。
- 日志轮转：按天切割，保留最近 30 天。


## 9. 附录

### 9.1 术语表

| 术语          | 解释                                                |
| ------------- | --------------------------------------------------- |
| Wails         | 基于 Go 和 Web 技术的桌面应用框架                   |
| Genkit        | Google 开源 AI 应用开发框架                         |
| Flow          | Genkit 中类型安全的 AI 工作流定义                   |
| Agent         | 具备自主规划、工具调用能力的 AI 实体                |
| Skill         | 可复用的任务执行知识包，包含步骤、陷阱和最佳实践    |
| MCP           | Model Context Protocol，AI 与外部工具的标准连接协议 |
| FTS5          | SQLite 全文检索扩展                                 |
| OpenTelemetry | 可观测性数据采集标准                                |

### 9.2 参考架构图

（见第 2 节系统总体架构图和技术架构图）

---

**文档版本**：V1.0
**创建日期**：2026-04-18
**最后更新**：2026-04-18