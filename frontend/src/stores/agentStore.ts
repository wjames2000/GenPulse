import { create } from 'zustand';

interface AgentConfig {
  id: string;
  name: string;
  role: string;
  description: string;
  model_config: {
    type: string;
    name: string;
    provider: string;
  };
  capabilities: string[];
  tools: string[];
  enabled: boolean;
}

interface AgentExecution {
  id: string;
  agent_id: string;
  task: string;
  state: string;
  started_at: string;
  completed_at?: string;
  error?: string;
  result?: any;
}

interface AgentStatus {
  id: string;
  name: string;
  state: string;
  execution_count: number;
  success_rate: number;
  avg_duration: string;
  current_execution?: {
    id: string;
    task: string;
    state: string;
    started_at: string;
  };
}

interface AgentStore {
  // 状态
  agents: AgentConfig[];
  executions: AgentExecution[];
  agentStatus: Record<string, AgentStatus>;
  loading: boolean;
  error: string | null;
  
  // 动作
  fetchAgents: () => Promise<void>;
  fetchAgentStatus: (agentId: string) => Promise<void>;
  fetchAllAgentsStatus: () => Promise<void>;
  executeAgent: (agentId: string, task: string, parameters?: Record<string, any>) => Promise<void>;
  fetchExecutions: () => Promise<void>;
  cancelExecution: (executionId: string) => Promise<void>;
  clearError: () => void;
}

export const useAgentStore = create<AgentStore>((set, get) => ({
  // 初始状态
  agents: [],
  executions: [],
  agentStatus: {},
  loading: false,
  error: null,

  // 获取Agent列表
  fetchAgents: async () => {
    set({ loading: true, error: null });
    try {
      // 调用后端API获取Agent列表
      const agents = await (window as any).go.main.App.ListAgents();
      set({ agents, loading: false });
    } catch (err) {
      set({ 
        error: err instanceof Error ? err.message : '获取Agent列表失败',
        loading: false 
      });
    }
  },

  // 获取单个Agent状态
  fetchAgentStatus: async (agentId: string) => {
    try {
      // 调用后端API获取Agent状态
      const status = await (window as any).go.main.App.GetAgentStatus(agentId);
      set((state) => ({
        agentStatus: {
          ...state.agentStatus,
          [agentId]: status,
        },
      }));
    } catch (err) {
      console.error('获取Agent状态失败:', err);
    }
  },

  // 获取所有Agent状态
  fetchAllAgentsStatus: async () => {
    try {
      // 调用后端API获取所有Agent状态
      const allStatus = await (window as any).go.main.App.GetAllAgentsStatus();
      set({ agentStatus: allStatus });
    } catch (err) {
      console.error('获取所有Agent状态失败:', err);
    }
  },

  // 执行Agent任务
  executeAgent: async (agentId: string, task: string, parameters: Record<string, any> = {}) => {
    set({ loading: true, error: null });
    try {
      // 调用后端API执行Agent任务
      const result = await (window as any).go.main.App.ExecuteAgent(agentId, task, parameters);
      
      // 更新执行列表
      set((state) => ({
        executions: [result, ...state.executions],
        loading: false,
      }));
      
      // 刷新Agent状态
      get().fetchAgentStatus(agentId);
    } catch (err) {
      set({ 
        error: err instanceof Error ? err.message : '执行Agent任务失败',
        loading: false 
      });
    }
  },

  // 获取执行历史
  fetchExecutions: async () => {
    try {
      // 调用后端API获取执行历史
      const executions = await (window as any).go.main.App.GetAgentExecutions();
      set({ executions });
    } catch (err) {
      console.error('获取执行历史失败:', err);
    }
  },

  // 取消执行
  cancelExecution: async (executionId: string) => {
    try {
      // 调用后端API取消执行
      await (window as any).go.main.App.CancelAgentExecution(executionId);
      
      // 刷新执行列表
      get().fetchExecutions();
    } catch (err) {
      set({ 
        error: err instanceof Error ? err.message : '取消执行失败'
      });
    }
  },

  // 清除错误
  clearError: () => {
    set({ error: null });
  },
}));