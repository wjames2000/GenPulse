# AI 软件项目开发流水线（桌面版）

> 🚀 **从 0 到 1，全自动写项目** — 基于 Go + Wails + Genkit 的多 Agent 协作开发平台

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Wails](https://img.shields.io/badge/Wails-v2-red?style=flat&logo=wails)](https://wails.io/)
[![Genkit](https://img.shields.io/badge/Genkit-Go_SDK-blue?style=flat)](https://firebase.google.com/docs/genkit-go)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)

[English](#README.md) | 中文

---

## 📖 简介

**GenPulse** 是一款桌面端 AI 开发工具，能够**全自动**地将自然语言需求转化为完整、可运行的软件项目。它并不是简单的代码补全助手，而是一个拥有**多角色 Agent 团队**、**自进化能力**和**可观测性**的“虚拟开发团队”。

你只需要输入一句话需求，系统就会自动调度 **产品经理、架构师、前后端开发、测试、运维** 等 AI Agent，协同完成从需求分析、架构设计到代码实现、测试部署的全部工作。更关键的是，它会**记住经验、沉淀技能**，越用越聪明。

---

## ✨ 核心特性

### 🤖 多 Agent 协作流水线
- **9 大内置角色**：Orchestrator（总指挥）、产品经理、架构师、前端/后端/全栈开发、测试、运维、代码审查。
- **层级化调度**：基于 Genkit 的 `Agent-as-Tool` 架构，Orchestrator 负责任务分解与分发，专业 Agent 各司其职。
- **并行执行**：前端与后端 Agent 可并行开发，大幅缩短项目生成时间。

### 🧠 自进化子系统（Hermes 风格）
- **Skills 闭环**：Agent 成功执行复杂任务后，自动提取关键步骤生成可复用的 Skill 文档（符合 agentskills.io 标准），下次遇到类似任务直接复用，效率提升 **40%+**。
- **三层记忆架构**：
  - **L1 工作记忆**：会话级上下文。
  - **L2 情节记忆**：跨会话持久化存储历史经验（SQLite + 全文检索）。
  - **L3 语义记忆**：长期记忆项目约定与用户编码偏好。

### 🔌 MCP 协议原生集成
- **无限扩展能力**：通过 Model Context Protocol 接入外部工具（GitHub、Figma、Puppeteer、数据库等）。
- **双向支持**：既可作为 MCP Client 消费外部工具，也可作为 MCP Server 暴露自身能力。

### 📊 全方位可观测性
- **实时监控仪表盘**：可视化展示每个 Agent 的状态、思维链、工具调用、Token 消耗和费用。
- **执行时间线**：甘特图展示并行/串行阶段。
- **文件变更预览**：内嵌 Diff 视图，查看 Agent 修改了哪些代码。

### 🖥️ 原生桌面体验（Go + Wails）
- **轻量高速**：安装包体积 < 30MB，冷启动 < 2 秒。
- **跨平台**：支持 Windows、macOS、Linux。
- **安全隔离**：文件操作与命令执行严格限制在项目目录内。

---

## 🎬 效果演示

> *占位：请替换为实际截图或 GIF*

| 配置界面                            | 执行监控                              | 技能管理                            |
| ----------------------------------- | ------------------------------------- | ----------------------------------- |
| ![config](./docs/images/config.png) | ![monitor](./docs/images/monitor.png) | ![skills](./docs/images/skills.png) |

---

## 🏗️ 技术架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Wails 桌面应用                           │
│   ┌──────────┐   ┌──────────┐   ┌──────────┐              │
│   │ React    │◄──│ Wails    │──►│ Go 后端  │              │
│   │ 前端界面  │   │ IPC 绑定  │   │ 核心引擎 │              │
│   └──────────┘   └──────────┘   └──────────┘              │
├─────────────────────────────────────────────────────────────┤
│                   Genkit Go SDK (AI 编排)                   │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐      │
│  │  Flow    │ │  Agent   │ │  Tool    │ │   MCP    │      │
│  │  引擎    │ │  管理器  │ │  注册表  │ │  Gateway │      │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘      │
├─────────────────────────────────────────────────────────────┤
│                    自进化子系统                             │
│  ┌────────────────────────┐ ┌──────────────────────────┐   │
│  │   Skills 闭环引擎      │ │   三层记忆 (SQLite)      │   │
│  └────────────────────────┘ └──────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│  文件系统 (os) │ go-git │ os/exec │ SQLite │ MCP Client    │
└─────────────────────────────────────────────────────────────┘
```

**核心依赖**：
- [Genkit Go SDK](https://firebase.google.com/docs/genkit-go) - AI 编排框架
- [Wails v2](https://wails.io/) - Go + Web 技术桌面框架
- [go-git](https://github.com/go-git/go-git) - 纯 Go 实现的 Git 操作库
- [SQLite](https://sqlite.org/) + FTS5 - 嵌入式数据库与全文检索
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/) - 可观测性追踪

---

## 🚀 快速开始

### 前置要求
- **Go** 1.21 或更高版本
- **Node.js** 18+ (用于前端构建)
- **Wails CLI**：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **API Key**：至少一个模型提供商（Gemini / OpenAI / Anthropic）

### 安装与运行

```bash
# 克隆仓库
git clone https://github.com/wjames2000/GenPulse.git
cd GenPulse

# 安装前端依赖
cd frontend && npm install && cd ..

# 开发模式运行
wails dev

# 生产构建
wails build
```

### 配置模型 API Key

首次启动后，在设置界面配置你的模型 API Key（密钥将安全存储在系统钥匙串中）。

---

## 📁 项目结构

```
GenPulse/
├── main.go                 # Wails 应用入口
├── app.go                  # 导出给前端的方法
├── frontend/               # React + TypeScript 前端
│   ├── src/
│   │   ├── components/     # UI 组件
│   │   ├── pages/          # 页面
│   │   └── services/       # Go 函数调用封装
│   └── package.json
├── internal/
│   ├── orchestrator/       # 流水线编排核心
│   ├── agents/             # 各 Agent 角色实现
│   ├── tools/              # 内置工具（fs, git, shell...）
│   ├── skills/             # 技能管理（生成/加载/验证）
│   ├── memory/             # 三层记忆系统
│   ├── mcp/                # MCP 网关
│   └── database/           # SQLite 数据层
├── skills/                 # 技能库存放目录
├── docs/                   # 详细设计文档
└── build/                  # 构建产物
```

---

## 📚 详细文档

- [项目需求说明书 V4.0](./docs/需求说明书_V4.0.md)
- [概要设计说明书](./docs/概要设计说明书.md)
- [数据库设计说明书](./docs/数据库设计说明书.md)
- [详细设计说明书](./docs/详细设计说明书.md)
- [WBS 任务分解](./docs/WBS.md)

---

## 🤝 贡献指南

我们非常欢迎社区贡献！无论是报告 Bug、提交功能建议，还是直接贡献代码，请先阅读 [CONTRIBUTING.md](./CONTRIBUTING.md)。

### 开发环境设置

```bash
# 安装开发依赖
make dev-setup

# 运行单元测试
make test

# 代码检查
make lint
```

---

## 📄 许可证

本项目采用 [MIT License](./LICENSE) 开源协议。

---

## 🙏 致谢

本项目深受以下优秀工作的启发：
- [Google Genkit](https://firebase.google.com/docs/genkit) - 生产级 AI 应用框架
- [Nous Research Hermes](https://github.com/NousResearch/Hermes) - 自进化 Agent 架构
- [Model Context Protocol](https://modelcontextprotocol.io/) - AI 工具互联标准
- [Wails](https://wails.io/) - 轻量级 Go 桌面框架

---

## ⭐ Star 历史

如果这个项目对你有帮助，请给我们一个 Star ⭐ 支持一下！

[![Star History Chart](https://api.star-history.com/svg?repos=wjames2000/GenPulse&type=Date)](https://star-history.com/#your-org/ai-dev-pipeline&Date)