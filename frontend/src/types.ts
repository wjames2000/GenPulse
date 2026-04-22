export type View = 'dashboard' | 'pipeline' | 'history' | 'neural' | 'settings' | 'kanban' | 'skills' | 'diff' | 'mcp-config' | 'monitoring';

export interface Agent {
  id: string;
  name: string;
  role: string;
  status: 'active' | 'waiting' | 'idle' | 'completed' | 'error';
  currentTask: string;
  progress: number;
  timeActive: string;
  type: 'orchestrator' | 'architect' | 'backend' | 'frontend' | 'qa' | 'devops' | 'reviewer' | 'product';
}

export interface TimelineEvent {
  id: string;
  agent: string;
  action: string;
  time: string;
  width: string;
  offset: string;
  isComplete: boolean;
  status?: 'running' | 'error' | 'warning';
  description?: string;
  metadata?: Record<string, any>;
}

export interface LogEntry {
  id: string;
  timestamp: string;
  level: 'info' | 'debug' | 'success' | 'warn' | 'error' | 'sys';
  message: string;
  agentId?: string;
  taskId?: string;
  details?: any;
  duration?: number;
  tags?: string[];
}

export interface Thought {
  id?: string;
  type: 'internal' | 'formulating';
  content?: string;
  isCode?: boolean;
  code?: string;
  filename?: string;
  agent?: string;
  timestamp?: string;
  metadata?: Record<string, any>;
}

export interface ToolInvocation {
  id: string;
  toolName: string;
  toolType?: string;
  agent?: string;
  timestamp: string;
  duration?: number;
  status: 'success' | 'error' | 'pending' | 'running';
  parameters?: Record<string, any>;
  result?: any;
  error?: string;
  metadata?: Record<string, any>;
}

export interface CostMetric {
  id: string;
  costType: 'llm' | 'api' | 'storage' | 'compute';
  agent?: string;
  timestamp: string;
  amount: number;
  tokenCount?: number;
  description?: string;
  metadata?: Record<string, any>;
}

export interface EvolutionEvent {
  id: string;
  eventType: 'skill_generated' | 'memory_updated' | 'agent_improved' | 'error_fixed';
  agent?: string;
  timestamp: string;
  description?: string;
  efficiencyGain?: number;
  tokenSavings?: number;
  timeSavings?: number;
  successRate?: number;
  severity?: 'low' | 'medium' | 'high';
  metadata?: Record<string, any>;
}

export interface FileDiff {
  id: string;
  filePath: string;
  changeType: 'added' | 'modified' | 'deleted';
  agent?: string;
  timestamp: string;
  diff?: string;
  linesAdded?: number;
  linesDeleted?: number;
  size?: number;
  commitHash?: string;
  metadata?: Record<string, any>;
}
