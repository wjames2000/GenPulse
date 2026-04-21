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

export interface Skill {
  id: string;
  name: string;
  description: string;
  category: string;
  version: string;
  enabled: boolean;
  validated: boolean;
  complexity: string;
  usage_count: number;
  success_rate: number;
  tags: string[];
  created_at?: string;
  updated_at?: string;
  agent_types?: string[];
  type?: string;
}

export interface SkillDetails extends Skill {
  steps: any[];
  examples: string[];
  tips: string[];
  warnings: string[];
  prerequisites: string[];
  related_tools: string[];
  token_estimate: number;
  avg_execution_time: string;
  source_task_id: string;
}

export interface EpisodicMemory {
  id: string;
  task_id: string;
  task_type: string;
  description: string;
  agent_id: string;
  agent_name: string;
  success: boolean;
  duration_ms: number;
  created_at: string;
  keywords: string[];
  tool_usage: Record<string, number>;
  context_data: Record<string, any>;
  relevance_score: number;
}

export interface SemanticMemory {
  user_id: string;
  username: string;
  preferences: Record<string, any>;
  skills: string[];
  interests: string[];
  goals: string[];
  working_style: Record<string, any>;
  communication: Record<string, any>;
  knowledge_areas: Record<string, number>;
  created_at: string;
  last_updated: string;
}

export interface EvolutionBenefits {
  skill_stats: {
    total_skills: number;
    enabled_skills: number;
    total_usage: number;
    average_success_rate: number;
  };
  memory_stats?: {
    total_memories: number;
    episodic_memories: number;
    semantic_memories: number;
    cache_hit_rate: number;
  };
  benefit_metrics: {
    total_token_savings: number;
    total_time_savings_seconds: number;
    automation_rate: number;
    estimated_cost_savings: number;
    productivity_gain: number;
  };
  trends: {
    skill_growth: number[];
    usage_growth: number[];
    success_rate_trend: number[];
  };
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

  // ==================== 技能管理 API ====================

  async getSkills(): Promise<Skill[]> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      // 暂时返回模拟数据
      return this.getMockSkills();
    } catch (error) {
      console.error('Failed to get skills:', error);
      return this.getMockSkills();
    }
  }

  async getSkillDetails(skillId: string): Promise<SkillDetails> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      // 暂时返回模拟数据
      return this.getMockSkillDetails(skillId);
    } catch (error) {
      console.error(`Failed to get skill details for ${skillId}:`, error);
      return this.getMockSkillDetails(skillId);
    }
  }

  async enableSkill(skillId: string): Promise<boolean> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      console.log(`Enabling skill: ${skillId}`);
      return true;
    } catch (error) {
      console.error(`Failed to enable skill ${skillId}:`, error);
      return false;
    }
  }

  async disableSkill(skillId: string): Promise<boolean> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      console.log(`Disabling skill: ${skillId}`);
      return true;
    } catch (error) {
      console.error(`Failed to disable skill ${skillId}:`, error);
      return false;
    }
  }

  async validateSkill(skillId: string): Promise<any> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      console.log(`Validating skill: ${skillId}`);
      return { overall_pass: true, timestamp: new Date().toISOString() };
    } catch (error) {
      console.error(`Failed to validate skill ${skillId}:`, error);
      return { overall_pass: false, timestamp: new Date().toISOString() };
    }
  }

  async getSkillStats(): Promise<any> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      return this.getMockSkillStats();
    } catch (error) {
      console.error('Failed to get skill stats:', error);
      return this.getMockSkillStats();
    }
  }

  // ==================== 记忆管理 API ====================

  async getEpisodicMemories(query: string = '', limit: number = 10): Promise<EpisodicMemory[]> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      return this.getMockEpisodicMemories();
    } catch (error) {
      console.error('Failed to get episodic memories:', error);
      return this.getMockEpisodicMemories();
    }
  }

  async getSemanticMemory(): Promise<SemanticMemory> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      return this.getMockSemanticMemory();
    } catch (error) {
      console.error('Failed to get semantic memory:', error);
      return this.getMockSemanticMemory();
    }
  }

  async getMemoryStats(): Promise<any> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      return this.getMockMemoryStats();
    } catch (error) {
      console.error('Failed to get memory stats:', error);
      return this.getMockMemoryStats();
    }
  }

  // ==================== 进化收益 API ====================

  async getEvolutionBenefits(): Promise<EvolutionBenefits> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      return this.getMockEvolutionBenefits();
    } catch (error) {
      console.error('Failed to get evolution benefits:', error);
      return this.getMockEvolutionBenefits();
    }
  }

  // ==================== 模拟数据 ====================

  private getMockSkills(): Skill[] {
    const now = new Date();
    return [
      {
        id: 'react-expert',
        name: 'React 专家',
        description: '精通现代 React 架构，能够生成高性能组件、处理复杂状态管理并优化渲染流水线。',
        category: 'frontend',
        version: 'v1.4.2',
        enabled: true,
        validated: true,
        complexity: 'complex',
        usage_count: 12450,
        success_rate: 0.984,
        tags: ['react', 'frontend', 'ui'],
        type: 'cognitive-skill',
        created_at: new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000).toISOString(),
        updated_at: new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000).toISOString(),
        agent_types: ['frontend', 'fullstack'],
      },
      {
        id: 'go-backend',
        name: 'Go 后端',
        description: '专注于高并发微服务架构，提供稳健的 API 设计、数据库分片策略与协程管理方案。',
        category: 'backend',
        version: 'v2.1.0',
        enabled: true,
        validated: true,
        complexity: 'complex',
        usage_count: 8204,
        success_rate: 0.952,
        tags: ['go', 'backend', 'api'],
        type: 'logic-processing',
        created_at: new Date(now.getTime() - 45 * 24 * 60 * 60 * 1000).toISOString(),
        updated_at: new Date(now.getTime() - 14 * 24 * 60 * 60 * 1000).toISOString(),
        agent_types: ['backend', 'devops'],
      },
      {
        id: 'git-pipeline',
        name: 'Git 流水线',
        description: '自动化 CI/CD 流程构建，代码合并冲突智能解决，部署策略优化。',
        category: 'devops',
        version: 'v1.0.8',
        enabled: true,
        validated: true,
        complexity: 'medium',
        usage_count: 45112,
        success_rate: 0.991,
        tags: ['git', 'ci-cd', 'automation'],
        type: 'ops-automation',
        created_at: new Date(now.getTime() - 60 * 24 * 60 * 60 * 1000).toISOString(),
        updated_at: new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000).toISOString(),
        agent_types: ['devops', 'reviewer'],
      },
    ];
  }

  private getMockSkillDetails(skillId: string): SkillDetails {
    const baseSkill = this.getMockSkills().find(s => s.id === skillId) || this.getMockSkills()[0];
    return {
      ...baseSkill,
      steps: [
        {
          id: 'analyze-requirements',
          order: 1,
          action: '分析需求并确定组件结构',
          tool: 'llm',
          parameters: [],
        },
        {
          id: 'create-component',
          order: 2,
          action: '创建React组件文件',
          tool: 'fs_write',
          parameters: [],
        },
      ],
      examples: [
        '生成一个带有搜索功能的用户列表组件',
        '创建支持拖拽排序的看板组件',
      ],
      tips: [
        '使用React.memo优化性能',
        '优先使用函数组件和hooks',
      ],
      warnings: [
        '避免在渲染函数中创建新对象',
        '注意useEffect的依赖数组',
      ],
      prerequisites: [],
      related_tools: ['fs_write', 'llm'],
      token_estimate: 1500,
      avg_execution_time: '2.5s',
      source_task_id: 'task-12345',
    };
  }

  private getMockSkillStats() {
    return {
      total_skills: 3,
      enabled_skills: 3,
      validated_skills: 3,
      by_category: {
        frontend: 1,
        backend: 1,
        devops: 1,
      },
      by_complexity: {
        complex: 2,
        medium: 1,
      },
      total_usage: 65866,
      average_success_rate: 0.975,
    };
  }

  private getMockEpisodicMemories(): EpisodicMemory[] {
    const now = new Date();
    return [
      {
        id: 'mem-001',
        task_id: 'task-001',
        task_type: 'code_generation',
        description: '生成React用户管理界面',
        agent_id: 'frontend-agent',
        agent_name: '前端开发Agent',
        success: true,
        duration_ms: 2450,
        created_at: new Date(now.getTime() - 2 * 60 * 60 * 1000).toISOString(),
        keywords: ['react', 'ui', 'user-management'],
        tool_usage: { 'fs_write': 3, 'llm': 2 },
        context_data: { framework: 'react', components: 5 },
        relevance_score: 0.95,
      },
      {
        id: 'mem-002',
        task_id: 'task-002',
        task_type: 'api_development',
        description: '创建用户认证API',
        agent_id: 'backend-agent',
        agent_name: '后端开发Agent',
        success: true,
        duration_ms: 3200,
        created_at: new Date(now.getTime() - 5 * 60 * 60 * 1000).toISOString(),
        keywords: ['go', 'api', 'authentication'],
        tool_usage: { 'fs_write': 4, 'llm': 3 },
        context_data: { language: 'go', endpoints: 3 },
        relevance_score: 0.88,
      },
      {
        id: 'mem-003',
        task_id: 'task-003',
        task_type: 'testing',
        description: '执行单元测试套件',
        agent_id: 'qa-agent',
        agent_name: '质量保证Agent',
        success: true,
        duration_ms: 1800,
        created_at: new Date(now.getTime() - 8 * 60 * 60 * 1000).toISOString(),
        keywords: ['testing', 'go-test', 'coverage'],
        tool_usage: { 'shell_exec': 5, 'fs_read': 10 },
        context_data: { tests: 25, coverage: 0.85 },
        relevance_score: 0.76,
      },
    ];
  }

  private getMockSemanticMemory(): SemanticMemory {
    const now = new Date();
    return {
      user_id: 'user-001',
      username: '开发者',
      preferences: {
        language: 'go',
        framework: 'react',
        database: 'postgresql',
        testing: '单元测试优先',
      },
      skills: [
        'Go 后端开发',
        'React 前端开发',
        '数据库设计',
        '微服务架构',
      ],
      interests: [
        'AI 辅助编程',
        '性能优化',
        '系统架构',
        '开发者工具',
      ],
      goals: [
        '构建高效的AI开发流水线',
        '实现代码自动生成与优化',
        '提升开发效率50%以上',
      ],
      working_style: {
        indentation: 'tabs',
        line_length: 100,
        naming: 'camelCase',
        comments: '必要的文档注释',
      },
      communication: {},
      knowledge_areas: {
        'go': 85,
        'react': 80,
        'postgresql': 75,
        'docker': 70,
      },
      created_at: new Date(now.getTime() - 90 * 24 * 60 * 60 * 1000).toISOString(),
      last_updated: new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000).toISOString(),
    };
  }

  private getMockMemoryStats() {
    return {
      total_memories: 156,
      episodic_memories: 150,
      semantic_memories: 1,
      working_memories: 5,
      total_searches: 245,
      avg_search_time_ms: 45.2,
      cache_hit_rate: 0.78,
      last_updated: new Date().toISOString(),
    };
  }

  private getMockEvolutionBenefits(): EvolutionBenefits {
    return {
      skill_stats: {
        total_skills: 3,
        enabled_skills: 3,
        total_usage: 65866,
        average_success_rate: 0.975,
      },
      memory_stats: {
        total_memories: 156,
        episodic_memories: 150,
        semantic_memories: 1,
        cache_hit_rate: 0.78,
      },
      benefit_metrics: {
        total_token_savings: 6586600,
        total_time_savings_seconds: 329330,
        automation_rate: 1.0,
        estimated_cost_savings: 13.1732,
        productivity_gain: 100,
      },
      trends: {
        skill_growth: [5, 8, 12, 15, 18, 22, 25, 28, 32, 35],
        usage_growth: [10, 25, 45, 70, 100, 135, 175, 220, 270, 325],
        success_rate_trend: [0.85, 0.87, 0.89, 0.91, 0.92, 0.93, 0.94, 0.94, 0.95, 0.95],
      },
    };
  }
}

export const api = new ApiService();