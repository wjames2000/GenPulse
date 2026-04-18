# AI Software Development Pipeline (Desktop Edition)

> 🚀 **From 0 to 1, Fully Automated Project Generation** — A Multi-Agent Collaborative Development Platform Built with Go + Wails + Genkit

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Wails](https://img.shields.io/badge/Wails-v2-red?style=flat&logo=wails)](https://wails.io/)
[![Genkit](https://img.shields.io/badge/Genkit-Go_SDK-blue?style=flat)](https://firebase.google.com/docs/genkit-go)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)

English | [中文](https://github.com/wjames2000/GenPulse/blob/main/README_cn.md)

---

## 📖 Introduction

**GenPulse** is a desktop AI development tool capable of **fully automatically** transforming natural language requirements into complete, runnable software projects. It is not merely a code completion assistant, but a "virtual development team" equipped with **multi-role Agents**, **self-evolution capabilities**, and **full observability**.

Simply input a one-sentence requirement, and the system will automatically orchestrate AI Agents acting as **Product Manager, Architect, Frontend/Backend Developers, QA, DevOps**, etc., to collaboratively complete the entire workflow—from requirements analysis and architecture design to code implementation, testing, and deployment. More importantly, it **remembers experiences and accumulates skills**, becoming smarter with every use.

---

## ✨ Core Features

### 🤖 Multi-Agent Collaborative Pipeline
- **9 Built-in Roles**: Orchestrator, Product Manager, Architect, Frontend/Backend/Full-stack Developer, QA Engineer, DevOps Engineer, Reviewer.
- **Hierarchical Scheduling**: Based on Genkit's `Agent-as-Tool` architecture, the Orchestrator decomposes and distributes tasks while specialized Agents execute their responsibilities.
- **Parallel Execution**: Frontend and Backend Agents can develop simultaneously, drastically reducing project generation time.

### 🧠 Self-Evolution Subsystem (Hermes Style)
- **Skills Closed Loop**: After successfully executing a complex task, Agents automatically extract key steps to generate reusable Skill documents (compliant with agentskills.io standard). Next time a similar task arises, efficiency improves by **over 40%** through reuse.
- **Three-Tier Memory Architecture**:
  - **L1 Working Memory**: Session-level context.
  - **L2 Episodic Memory**: Persistent cross-session storage of historical experiences (SQLite + Full-Text Search).
  - **L3 Semantic Memory**: Long-term retention of project conventions and user coding preferences.

### 🔌 Native MCP Protocol Integration
- **Unlimited Extensibility**: Connect external tools (GitHub, Figma, Puppeteer, databases, etc.) via Model Context Protocol.
- **Bidirectional Support**: Acts as both an MCP Client consuming external tools and an MCP Server exposing local capabilities.

### 📊 Comprehensive Observability
- **Real-Time Monitoring Dashboard**: Visualize each Agent's status, thought process, tool calls, token consumption, and cost.
- **Execution Timeline**: Gantt chart displaying parallel/serial stages.
- **File Change Preview**: Built-in diff view to inspect code modifications made by Agents.

### 🖥️ Native Desktop Experience (Go + Wails)
- **Lightweight & Fast**: Installer size < 30MB, cold start < 2 seconds.
- **Cross-Platform**: Supports Windows, macOS, and Linux.
- **Secure Isolation**: File operations and command executions are strictly confined within the project directory.

---

## 🎬 Demo

> *Placeholder: Replace with actual screenshots or GIFs*

| Configuration                       | Execution Monitor                     | Skill Management                    |
| ----------------------------------- | ------------------------------------- | ----------------------------------- |
| ![config](./docs/images/config.png) | ![monitor](./docs/images/monitor.png) | ![skills](./docs/images/skills.png) |

---

## 🏗️ Technical Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Wails Desktop Application                 │
│   ┌──────────┐   ┌──────────┐   ┌──────────┐              │
│   │ React    │◄──│ Wails    │──►│ Go Backend│             │
│   │ Frontend │   │ IPC Bind │   │ Core Engine│            │
│   └──────────┘   └──────────┘   └──────────┘              │
├─────────────────────────────────────────────────────────────┤
│                   Genkit Go SDK (AI Orchestration)          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐      │
│  │  Flow    │ │  Agent   │ │  Tool    │ │   MCP    │      │
│  │  Engine  │ │  Manager │ │ Registry │ │  Gateway │      │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘      │
├─────────────────────────────────────────────────────────────┤
│                   Self-Evolution Subsystem                  │
│  ┌────────────────────────┐ ┌──────────────────────────┐   │
│  │   Skills Closed-Loop   │ │  Three-Tier Memory       │   │
│  │        Engine          │ │      (SQLite)            │   │
│  └────────────────────────┘ └──────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│  FileSystem (os) │ go-git │ os/exec │ SQLite │ MCP Client  │
└─────────────────────────────────────────────────────────────┘
```

**Core Dependencies**:
- [Genkit Go SDK](https://firebase.google.com/docs/genkit-go) - AI orchestration framework
- [Wails v2](https://wails.io/) - Go + Web tech desktop framework
- [go-git](https://github.com/go-git/go-git) - Pure Go Git implementation
- [SQLite](https://sqlite.org/) + FTS5 - Embedded database with full-text search
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/) - Observability tracing

---

## 🚀 Quick Start

### Prerequisites
- **Go** 1.21 or higher
- **Node.js** 18+ (for frontend build)
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **API Key**: At least one model provider (Gemini / OpenAI / Anthropic)

### Installation & Running

```bash
# Clone the repository
git clone https://github.com/wjames2000/GenPulse.git
cd GenPulse

# Install frontend dependencies
cd frontend && npm install && cd ..

# Run in development mode
wails dev

# Build for production
wails build
```

### Configuring Model API Keys

After the first launch, configure your model API keys in the settings interface (keys will be securely stored in the system keychain).

---

## 📁 Project Structure

```
GenPulse/
├── main.go                 # Wails application entry point
├── app.go                  # Methods exposed to the frontend
├── frontend/               # React + TypeScript frontend
│   ├── src/
│   │   ├── components/     # UI components
│   │   ├── pages/          # Pages
│   │   └── services/       # Go function call wrappers
│   └── package.json
├── internal/
│   ├── orchestrator/       # Pipeline orchestration core
│   ├── agents/             # Implementation of Agent roles
│   ├── tools/              # Built-in tools (fs, git, shell...)
│   ├── skills/             # Skill management (generate/load/validate)
│   ├── memory/             # Three-tier memory system
│   ├── mcp/                # MCP gateway
│   └── database/           # SQLite data layer
├── skills/                 # Skill library storage directory
├── docs/                   # Detailed design documentation
└── build/                  # Build artifacts
```

---

## 📚 Detailed Documentation

- [Project Requirements Specification V4.0](./docs/requirements_V4.0.md)
- [High-Level Design Document](./docs/high_level_design.md)
- [Database Design Document](./docs/database_design.md)
- [Detailed Design Document](./docs/detailed_design.md)
- [WBS Task Breakdown](./docs/WBS.md)

---

## 🤝 Contributing

Community contributions are highly welcome! Whether reporting bugs, suggesting features, or contributing code directly, please read [CONTRIBUTING.md](./CONTRIBUTING.md) first.

### Development Environment Setup

```bash
# Install development dependencies
make dev-setup

# Run unit tests
make test

# Lint code
make lint
```

---

## 📄 License

This project is open-sourced under the [MIT License](./LICENSE).

---

## 🙏 Acknowledgements

This project is heavily inspired by the following excellent work:
- [Google Genkit](https://firebase.google.com/docs/genkit) - Production-grade AI framework
- [Nous Research Hermes](https://github.com/NousResearch/Hermes) - Self-evolving Agent architecture
- [Model Context Protocol](https://modelcontextprotocol.io/) - Standard for AI tool interconnection
- [Wails](https://wails.io/) - Lightweight Go desktop framework

---

## ⭐ Star History

If this project is helpful to you, please give us a Star ⭐ to show your support!

[![Star History Chart](https://api.star-history.com/svg?repos=wjames2000/GenPulse&type=Date)](https://star-history.com/#your-org/ai-dev-pipeline&Date)