import * as App from '../../wailsjs/go/main/App';

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

export interface ProjectConfig {
  name: string;
  path: string;
  architecture: 'web' | 'cli' | 'micro';
  frontend: string;
  backend: string;
  database: string;
  mobile: string;
}

export interface Execution {
  id: string;
  agent_id: string;
  task: string;
  state: 'pending' | 'executing' | 'completed' | 'failed';
  started_at: string;
  completed_at?: string;
  result?: any;
  error?: string;
  parameters: Record<string, any>;
}

export interface LogEntry {
  timestamp: string;
  level: 'info' | 'debug' | 'success' | 'warn' | 'error' | 'sys';
  message: string;
}

class ApiService {
  async getAppInfo() {
    try {
      return await App.GetAppInfo();
    } catch (error) {
      console.error('Failed to get app info:', error);
      return null;
    }
  }

  async listAgents() {
    try {
      const agents = await App.ListAgents();
      return agents || [];
    } catch (error) {
      console.error('Failed to list agents:', error);
      return [];
    }
  }

  async getAgentStatus(agentId: string) {
    try {
      return await App.GetAgentStatus(agentId);
    } catch (error) {
      console.error(`Failed to get agent status for ${agentId}:`, error);
      return null;
    }
  }

  async getAllAgentsStatus() {
    try {
      return await App.GetAllAgentsStatus();
    } catch (error) {
      console.error('Failed to get all agents status:', error);
      return {};
    }
  }

  async executeAgent(agentId: string, task: string, parameters: Record<string, any> = {}) {
    try {
      return await App.ExecuteAgent(agentId, task, parameters);
    } catch (error) {
      console.error(`Failed to execute agent ${agentId}:`, error);
      throw error;
    }
  }

  async getAgentExecutions() {
    try {
      return await App.GetAgentExecutions();
    } catch (error) {
      console.error('Failed to get agent executions:', error);
      return [];
    }
  }

  async cancelAgentExecution(executionId: string) {
    try {
      await App.CancelAgentExecution(executionId);
      return true;
    } catch (error) {
      console.error(`Failed to cancel execution ${executionId}:`, error);
      return false;
    }
  }

  async healthCheck() {
    try {
      return await App.HealthCheck();
    } catch (error) {
      console.error('Health check failed:', error);
      return { status: 'error', message: 'Health check failed' };
    }
  }

  async logMessage(level: string, message: string) {
    try {
      await App.LogMessage(level, message);
    } catch (error) {
      console.error('Failed to log message:', error);
    }
  }

  async getLogs() {
    try {
      return await App.GetLogs();
    } catch (error) {
      console.error('Failed to get logs:', error);
      return [];
    }
  }

  async initializeProject(config: ProjectConfig) {
    try {
      await this.logMessage('info', `Initializing project: ${config.name}`);
      
      const result = await this.executeAgent('orchestrator', 'initialize_project', {
        project_name: config.name,
        project_path: config.path,
        architecture: config.architecture,
        tech_stack: {
          frontend: config.frontend,
          backend: config.backend,
          database: config.database,
          mobile: config.mobile
        }
      });

      await this.logMessage('success', `Project ${config.name} initialized successfully`);
      return result;
    } catch (error) {
      await this.logMessage('error', `Failed to initialize project: ${error}`);
      throw error;
    }
  }

  async startPipeline(projectId: string) {
    try {
      await this.logMessage('info', `Starting pipeline for project: ${projectId}`);
      
      const result = await this.executeAgent('orchestrator', 'start_pipeline', {
        project_id: projectId
      });

      await this.logMessage('success', `Pipeline started for project: ${projectId}`);
      return result;
    } catch (error) {
      await this.logMessage('error', `Failed to start pipeline: ${error}`);
      throw error;
    }
  }

  async getProjectStatus(projectId: string) {
    try {
      const status = await this.getAllAgentsStatus();
      return {
        project_id: projectId,
        agents_status: status,
        overall_status: this.calculateOverallStatus(status)
      };
    } catch (error) {
      console.error(`Failed to get project status for ${projectId}:`, error);
      return null;
    }
  }

  private calculateOverallStatus(agentsStatus: any): 'idle' | 'active' | 'completed' | 'error' {
    if (!agentsStatus || typeof agentsStatus !== 'object') {
      return 'idle';
    }

    const statuses = Object.values(agentsStatus);
    if (statuses.some((s: any) => s?.state === 'error')) return 'error';
    if (statuses.some((s: any) => s?.state === 'active')) return 'active';
    if (statuses.every((s: any) => s?.state === 'completed')) return 'completed';
    return 'idle';
  }
}

export const api = new ApiService();