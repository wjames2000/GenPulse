# 监控数据采集系统

本文档描述了GenPulse项目的监控数据采集系统，实现了WBS第四阶段4.2节的所有功能。

## 功能概述

监控数据采集系统包括以下核心组件：

1. **OpenTelemetry集成** - 分布式追踪和指标收集
2. **指标收集器** - 自定义业务指标收集
3. **执行轨迹存储** - SQLite存储执行轨迹
4. **实时状态推送** - 事件总线与前端通信

## 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    监控数据采集系统                           │
├─────────────────────────────────────────────────────────────┤
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │ OpenTelemetry │  │ 指标收集器  │  │ 执行轨迹存储 │            │
│  │  集成配置    │  │ Metrics   │  │ Trace Store│            │
│  └────────────┘  └────────────┘  └────────────┘            │
│                                                            │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                  事件总线 (Event Bus)                  │  │
│  │  • 发布/订阅模式                                     │  │
│  │  • 实时状态推送                                      │  │
│  │  • Wails前端桥接                                     │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## 组件详细说明

### 1. OpenTelemetry集成 (`internal/monitoring/telemetry/`)

**功能**：
- 配置OpenTelemetry追踪导出器（控制台、OTLP）
- 提供Tracer包装器简化使用
- 支持函数级追踪装饰器

**关键文件**：
- `config.go` - OpenTelemetry配置和初始化
- `tracer.go` - Tracer包装器和追踪辅助函数

**使用示例**：
```go
tracer := telemetry.NewTracer("my-service")
ctx, spanCtx := tracer.Start(ctx, "my-operation")
defer spanCtx.End()

// 设置属性
spanCtx.SetAttributes(
    attribute.String("key", "value"),
    attribute.Int("count", 42),
)

// 记录错误
if err != nil {
    spanCtx.RecordError(err)
    spanCtx.SetStatus(codes.Error, err.Error())
}
```

### 2. 指标收集器 (`internal/monitoring/metrics/`)

**功能**：
- 定义和注册自定义业务指标
- 支持计数器、仪表盘、直方图三种指标类型
- 自动记录Agent执行、工具调用、Token消耗等指标

**预定义指标**：
- `agent.execution.count` - Agent执行次数
- `agent.execution.duration` - Agent执行耗时
- `tool.invocation.count` - 工具调用次数
- `token.consumption.total` - Token消耗总量
- `pipeline.execution.count` - 流水线执行次数
- `skill.generation.count` - 技能生成次数
- `skill.reuse.count` - 技能复用次数

**使用示例**：
```go
collector, _ := metrics.NewMetricsCollector("genpulse")

// 记录Agent执行
collector.RecordAgentExecution("orchestrator", 150*time.Millisecond, "success")

// 记录工具调用
collector.RecordToolInvocation("fs_write", 50*time.Millisecond, "success")

// 记录Token消耗
collector.RecordTokenConsumption("gpt-4", "completion", 1500)
```

### 3. 执行轨迹存储 (`internal/monitoring/storage/`)

**功能**：
- SQLite存储执行轨迹数据
- 支持复杂的查询和统计
- 自动清理旧数据（可配置保留天数）

**数据结构**：
```go
type TraceSpan struct {
    ID            string                 // 跨度ID
    TraceID       string                 // 追踪ID
    ParentID      string                 // 父跨度ID
    Name          string                 // 跨度名称
    Kind          string                 // 跨度类型
    StartTime     time.Time              // 开始时间
    EndTime       time.Time              // 结束时间
    Duration      time.Duration          // 持续时间
    Attributes    map[string]interface{} // 属性
    Events        []TraceEvent           // 事件列表
    Status        string                 // 状态
    StatusMessage string                 // 状态消息
    AgentName     string                 // Agent名称
    ToolName      string                 // 工具名称
    PipelineID    string                 // 流水线ID
}
```

**使用示例**：
```go
store, _ := storage.NewTraceStore("data/traces.db")

// 存储跨度
span := &storage.TraceSpan{
    ID:        "span_1",
    TraceID:   "trace_1",
    Name:      "agent.execution",
    StartTime: startTime,
    EndTime:   endTime,
    Duration:  duration,
    Status:    "OK",
    AgentName: "orchestrator",
}
store.StoreSpan(ctx, span)

// 查询跨度
spans, _ := store.GetRecentSpans(ctx, 10)
spans, _ := store.GetAgentSpans(ctx, "orchestrator", 5)
spans, _ := store.GetPipelineSpans(ctx, "pipeline_123")

// 获取统计
stats, _ := store.GetStatistics(ctx, startTime, endTime)
```

### 4. 事件总线 (`internal/monitoring/events/`)

**功能**：
- 发布/订阅模式的事件系统
- 支持多种事件类型
- Wails前端桥接
- 事件缓冲和历史记录

**事件类型**：
- `agent_started` - Agent开始执行
- `agent_completed` - Agent执行完成
- `tool_invoked` - 工具调用开始
- `tool_completed` - 工具调用完成
- `pipeline_started` - 流水线开始
- `pipeline_completed` - 流水线完成
- `metric_updated` - 指标更新
- `trace_recorded` - 轨迹记录
- `skill_generated` - 技能生成
- `error_occurred` - 错误发生
- `status_changed` - 状态变更

**使用示例**：
```go
eventBus := events.NewEventBus(events.DefaultConfig())

// 发布事件
eventBus.Publish(ctx, events.EventTypeAgentStarted, "orchestrator", map[string]interface{}{
    "agent_name": "orchestrator",
    "input": "Create web app",
})

// 订阅事件
eventBus.Subscribe(events.EventTypeAgentCompleted, func(ctx context.Context, event events.Event) error {
    fmt.Printf("Agent completed: %v\n", event.Data["agent_name"])
    return nil
})

// Wails前端桥接
wailsBridge := events.NewWailsEventBridge(eventBus)
wailsBridge.SetAppContext(appCtx)
wailsBridge.StartForwarding()
```

### 5. 监控服务 (`internal/monitoring/service.go`)

**功能**：
- 集成所有监控组件
- 提供统一接口
- 自动事件处理
- 生命周期管理

**使用示例**：
```go
config := monitoring.DefaultConfig()
config.TraceStorePath = "data/traces.db"

service, _ := monitoring.NewMonitoringService(config)
defer service.Shutdown(ctx)

// 获取各个组件
telemetry := service.GetTelemetry()
metrics := service.GetMetrics()
traceStore := service.GetTraceStore()
eventBus := service.GetEventBus()

// 发布事件
service.PublishEvent(ctx, events.EventTypeAgentStarted, "demo", map[string]interface{}{
    "agent_name": "test_agent",
})
```

## 配置说明

### 监控服务配置
```go
type Config struct {
    Enabled           bool              // 是否启用监控
    TelemetryConfig   telemetry.Config  // OpenTelemetry配置
    TraceStorePath    string            // 轨迹存储路径
    EventBufferSize   int               // 事件缓冲区大小
    RetentionDays     int               // 数据保留天数
}
```

### OpenTelemetry配置
```go
type Config struct {
    Enabled           bool
    ServiceName       string
    ServiceVersion    string
    Environment       string
    OTLPEndpoint      string    // OTLP端点地址
    OTLPInsecure      bool      // 是否使用不安全连接
    ConsoleExport     bool      // 是否输出到控制台
    BatchTimeout      time.Duration
    ExportTimeout     time.Duration
    MaxExportBatchSize int
}
```

## 示例代码

### 完整示例
参考 `examples/monitoring_demo.go` 和 `examples/monitoring_simple_demo.go`

### 测试示例
参考 `test/monitoring/monitoring_integration_test.go`

## 与现有系统集成

### 与Agent系统集成
监控系统可以无缝集成到现有的Agent系统中：

1. **在Agent执行前后添加追踪**：
```go
func (a *Agent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    tracer := monitoring.GetTracer("agent")
    ctx, spanCtx := tracer.Start(ctx, fmt.Sprintf("agent.%s", a.Name))
    defer spanCtx.End()

    startTime := time.Now()
    result, err := a.doExecute(ctx, input)
    duration := time.Since(startTime)

    // 记录指标
    monitoring.RecordAgentExecution(a.Name, duration, err)

    // 发布事件
    if err != nil {
        monitoring.PublishEvent(ctx, events.EventTypeAgentCompleted, 
            a.Name, map[string]interface{}{
                "agent_name": a.Name,
                "duration":   duration.Milliseconds(),
                "status":     "error",
                "error":      err.Error(),
            })
    } else {
        monitoring.PublishEvent(ctx, events.EventTypeAgentCompleted,
            a.Name, map[string]interface{}{
                "agent_name": a.Name,
                "duration":   duration.Milliseconds(),
                "status":     "success",
                "output":     result,
            })
    }

    return result, err
}
```

### 与流水线系统集成
在流水线执行过程中集成监控：

```go
func (p *Pipeline) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    // 发布流水线开始事件
    monitoring.PublishEvent(ctx, events.EventTypePipelineStarted,
        p.Name, map[string]interface{}{
            "pipeline_id":   p.ID,
            "pipeline_name": p.Name,
            "input":         input,
        })

    startTime := time.Now()
    result, err := p.executeStages(ctx, input)
    duration := time.Since(startTime)

    // 记录流水线指标
    status := "success"
    if err != nil {
        status = "error"
    }
    monitoring.RecordPipelineExecution(p.Name, status)

    // 发布流水线完成事件
    monitoring.PublishEvent(ctx, events.EventTypePipelineCompleted,
        p.Name, map[string]interface{}{
            "pipeline_id":   p.ID,
            "pipeline_name": p.Name,
            "output":        result,
            "duration":      duration.Milliseconds(),
            "status":        status,
        })

    return result, err
}
```

## 前端监控仪表盘

监控系统通过Wails事件桥接将数据推送到前端，前端可以：

1. **实时显示Agent状态**
2. **展示执行时间线**
3. **显示指标图表**
4. **查看执行日志**
5. **监控技能生成和复用**

前端组件可以通过订阅 `monitoring_event` 事件接收实时数据：

```javascript
// 前端JavaScript示例
window.runtime.EventsOn('monitoring_event', (eventData) => {
    const event = JSON.parse(eventData);
    switch (event.type) {
        case 'agent_started':
            updateAgentStatus(event.data.agent_name, 'running');
            break;
        case 'agent_completed':
            updateAgentStatus(event.data.agent_name, 'completed');
            addToTimeline(event);
            break;
        case 'metric_updated':
            updateMetricsChart(event.data);
            break;
    }
});
```

## 性能考虑

1. **异步处理**：所有监控操作都是异步的，不会阻塞主业务流程
2. **批量处理**：OpenTelemetry使用批量处理器减少IO操作
3. **内存管理**：事件缓冲区有大小限制，防止内存泄漏
4. **数据库优化**：SQLite使用索引优化查询性能
5. **资源清理**：定期清理旧数据，控制存储增长

## 扩展性

监控系统设计为可扩展的：

1. **添加新指标**：通过 `MetricsCollector.RegisterMetric()` 注册
2. **添加新事件类型**：在 `EventType` 常量中添加
3. **自定义导出器**：实现OpenTelemetry的SpanExporter接口
4. **替换存储后端**：实现TraceStore接口支持其他数据库

## 故障排除

### 常见问题

1. **OpenTelemetry初始化失败**
   - 检查依赖是否正确安装
   - 验证配置参数
   - 查看控制台错误输出

2. **SQLite数据库锁**
   - 确保只有一个进程访问数据库
   - 检查文件权限
   - 使用内存数据库测试

3. **事件丢失**
   - 增加事件缓冲区大小
   - 检查事件处理器是否阻塞
   - 验证Wails连接状态

4. **性能问题**
   - 减少控制台输出
   - 调整批量处理参数
   - 限制历史数据量

### 调试模式

启用调试模式查看详细日志：
```go
config := monitoring.DefaultConfig()
config.TelemetryConfig.ConsoleExport = true
config.TelemetryConfig.Enabled = true
```

## 总结

监控数据采集系统为GenPulse项目提供了完整的可观测性解决方案，实现了：

1. ✅ **OpenTelemetry集成** - 分布式追踪和指标
2. ✅ **指标收集器** - 业务指标监控
3. ✅ **执行轨迹存储** - 历史数据持久化
4. ✅ **实时状态推送** - 前端实时更新
5. ✅ **完整测试覆盖** - 确保系统可靠性

该系统为第四阶段的监控仪表盘开发提供了坚实的数据基础。