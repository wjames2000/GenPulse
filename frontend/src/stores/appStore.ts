import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { AppInfo, HealthStatus } from '../services/baseService';

// 应用状态接口
export interface AppState {
  // 应用信息
  appInfo: AppInfo | null;
  healthStatus: HealthStatus | null;
  isInitialized: boolean;
  
  // UI状态
  sidebarOpen: boolean;
  darkMode: boolean;
  currentView: string;
  
  // 项目状态
  currentProject: string | null;
  projectList: string[];
  
  // Agent状态
  activeAgents: string[];
  agentStatus: Record<string, 'idle' | 'running' | 'completed' | 'error'>;
  
  // 执行状态
  isExecuting: boolean;
  executionProgress: number;
  executionLogs: Array<{
    id: string;
    level: string;
    message: string;
    timestamp: string;
  }>;
  
  // Actions
  setAppInfo: (info: AppInfo) => void;
  setHealthStatus: (status: HealthStatus) => void;
  setInitialized: (initialized: boolean) => void;
  
  toggleSidebar: () => void;
  toggleDarkMode: () => void;
  setCurrentView: (view: string) => void;
  
  setCurrentProject: (project: string | null) => void;
  setProjectList: (projects: string[]) => void;
  addProject: (project: string) => void;
  removeProject: (project: string) => void;
  
  setAgentStatus: (agentId: string, status: 'idle' | 'running' | 'completed' | 'error') => void;
  setActiveAgents: (agents: string[]) => void;
  
  startExecution: () => void;
  stopExecution: () => void;
  setExecutionProgress: (progress: number) => void;
  addExecutionLog: (level: string, message: string) => void;
  clearExecutionLogs: () => void;
  
  // 重置状态
  reset: () => void;
}

// 创建应用状态store
export const useAppStore = create<AppState>()(
  persist(
    (set, get) => ({
      // 初始状态
      appInfo: null,
      healthStatus: null,
      isInitialized: false,
      
      sidebarOpen: true,
      darkMode: false,
      currentView: 'dashboard',
      
      currentProject: null,
      projectList: [],
      
      activeAgents: [],
      agentStatus: {},
      
      isExecuting: false,
      executionProgress: 0,
      executionLogs: [],
      
      // Actions实现
      setAppInfo: (info) => set({ appInfo: info }),
      setHealthStatus: (status) => set({ healthStatus: status }),
      setInitialized: (initialized) => set({ isInitialized: initialized }),
      
      toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
      toggleDarkMode: () => set((state) => ({ darkMode: !state.darkMode })),
      setCurrentView: (view) => set({ currentView: view }),
      
      setCurrentProject: (project) => set({ currentProject: project }),
      setProjectList: (projects) => set({ projectList: projects }),
      addProject: (project) => 
        set((state) => ({
          projectList: [...state.projectList, project].filter((p, i, arr) => arr.indexOf(p) === i)
        })),
      removeProject: (project) => 
        set((state) => ({
          projectList: state.projectList.filter(p => p !== project),
          currentProject: state.currentProject === project ? null : state.currentProject
        })),
      
      setAgentStatus: (agentId, status) => 
        set((state) => ({
          agentStatus: {
            ...state.agentStatus,
            [agentId]: status
          }
        })),
      setActiveAgents: (agents) => set({ activeAgents: agents }),
      
      startExecution: () => set({ 
        isExecuting: true,
        executionProgress: 0,
        executionLogs: []
      }),
      stopExecution: () => set({ 
        isExecuting: false,
        executionProgress: 0
      }),
      setExecutionProgress: (progress) => set({ executionProgress: Math.min(100, Math.max(0, progress)) }),
      addExecutionLog: (level, message) =>
        set((state) => ({
          executionLogs: [
            ...state.executionLogs,
            {
              id: Date.now().toString(),
              level,
              message,
              timestamp: new Date().toISOString()
            }
          ].slice(-100) // 只保留最近100条日志
        })),
      clearExecutionLogs: () => set({ executionLogs: [] }),
      
      // 重置状态
      reset: () => set({
        appInfo: null,
        healthStatus: null,
        isInitialized: false,
        sidebarOpen: true,
        currentView: 'dashboard',
        currentProject: null,
        projectList: [],
        activeAgents: [],
        agentStatus: {},
        isExecuting: false,
        executionProgress: 0,
        executionLogs: []
      })
    }),
    {
      name: 'genpulse-app-storage', // localStorage key
      partialize: (state) => ({
        // 只持久化部分状态
        sidebarOpen: state.sidebarOpen,
        darkMode: state.darkMode,
        currentView: state.currentView,
        currentProject: state.currentProject,
        projectList: state.projectList
      })
    }
  )
);

// 选择器函数（用于优化性能）
export const selectAppInfo = (state: AppState) => state.appInfo;
export const selectHealthStatus = (state: AppState) => state.healthStatus;
export const selectIsInitialized = (state: AppState) => state.isInitialized;
export const selectSidebarOpen = (state: AppState) => state.sidebarOpen;
export const selectDarkMode = (state: AppState) => state.darkMode;
export const selectCurrentView = (state: AppState) => state.currentView;
export const selectIsExecuting = (state: AppState) => state.isExecuting;
export const selectExecutionProgress = (state: AppState) => state.executionProgress;
export const selectExecutionLogs = (state: AppState) => state.executionLogs;
export const selectAgentStatus = (state: AppState) => state.agentStatus;