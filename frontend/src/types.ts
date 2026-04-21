export type View = 'dashboard' | 'pipeline' | 'history' | 'neural' | 'settings' | 'kanban' | 'skills' | 'diff';

export interface Agent {
  id: string;
  name: string;
  role: string;
  status: 'active' | 'waiting' | 'idle';
  currentTask: string;
  progress: number;
  timeActive: string;
  type: 'orchestrator' | 'architect' | 'backend' | 'frontend' | 'qa';
}

export interface TimelineEvent {
  id: string;
  agent: string;
  action: string;
  time: string;
  width: string;
  offset: string;
  isComplete: boolean;
}

export interface LogEntry {
  timestamp: string;
  level: 'info' | 'debug' | 'success' | 'warn' | 'error' | 'sys';
  message: string;
}

export interface Thought {
  type: 'internal' | 'formulating';
  content?: string;
  isCode?: boolean;
  code?: string;
  filename?: string;
}
