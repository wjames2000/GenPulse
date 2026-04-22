import { Agent, LogEntry, TimelineEvent, Thought } from "./types";

export const AGENTS: Agent[] = [
  {
    id: "orchestrator-1",
    name: "编排者",
    role: "Orchestrator",
    status: "active",
    currentTask: "分析项目需求并分配微服务架构任务至各开发节点",
    progress: 80,
    timeActive: "14m 23s",
    type: "orchestrator"
  },
  {
    id: "architect-1",
    name: "架构师",
    role: "Architect",
    status: "active",
    currentTask: "设计用户认证模块的数据流图与数据库 schema",
    progress: 45,
    timeActive: "08m 45s",
    type: "architect"
  },
  {
    id: "backend-1",
    name: "后端开发",
    role: "Backend Dev",
    status: "waiting",
    currentTask: "等待架构师输出 schema 定义...",
    progress: 0,
    timeActive: "--:--",
    type: "backend"
  },
  {
    id: "frontend-1",
    name: "前端开发",
    role: "Frontend Dev",
    status: "idle",
    currentTask: "等待 API 接口定义...",
    progress: 0,
    timeActive: "--:--",
    type: "frontend"
  }
];

export const TIMELINE: TimelineEvent[] = [
  { id: "1", agent: "Orchestrator", action: "init_project", time: "T-14m", width: "25%", offset: "2%", isComplete: true },
  { id: "2", agent: "Orchestrator", action: "analyze_reqs", time: "T-10m", width: "15%", offset: "28%", isComplete: false },
  { id: "3", agent: "Architect", action: "design_db_schema", time: "T-8m", width: "30%", offset: "28%", isComplete: true },
  { id: "4", agent: "Architect", action: "auth_flow", time: "T-4m", width: "20%", offset: "60%", isComplete: false },
];

export const LOGS: LogEntry[] = [
  { id: "log-1", timestamp: "10:42:01", level: "info", message: "Orchestrator: Spawned new task cluster 'auth_module'" },
  { id: "log-2", timestamp: "10:42:03", level: "info", message: "MessageBus: Routed prompt to Architect ID-4A9B" },
  { id: "log-3", timestamp: "10:42:05", level: "debug", message: "Architect: Analyzing legacy codebase for auth dependencies..." },
  { id: "log-4", timestamp: "10:42:12", level: "success", message: "Architect: Legacy analysis complete. Found 3 deprecation warnings." },
  { id: "log-5", timestamp: "10:42:15", level: "debug", message: "Architect: Generating initial schema proposal..." },
  { id: "log-6", timestamp: "10:42:28", level: "sys", message: "Memory allocation at 64%" },
  { id: "log-7", timestamp: "10:42:35", level: "debug", message: "Architect: Evaluating OAuth2 vs JWT implementation paths..." },
];

export const THOUGHTS: Thought[] = [
  {
    type: "internal",
    content: "为了确保新旧系统的兼容性，我们需要保留旧的 session_id 映射，同时引入基于 JWT 的无状态验证。"
  },
  {
    type: "internal",
    isCode: true,
    filename: "schema.prisma",
    code: `model User {
  id String @id @default(uuid())
  // legacy_token String?
  password_hash String
  refresh_tokens RefreshToken[]
  createdAt DateTime @default(now())
}`
  },
  {
    type: "formulating",
    content: "考虑到前端是 React SPA，建议使用 HttpOnly cookie 存储 refresh token，内存存储 access token 以防止 XSS 攻击。正在构建最终架构文档..."
  }
];
