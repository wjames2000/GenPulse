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

// ==================== MCP 配置接口 ====================

export interface MCPClientConfig {
  server_type: string;
  command: string;
  args: string[];
  namespace: string;
  timeout: number;
}

export interface MCPServerConfig {
  type: string;
  tool_filter?: string;
}

export interface MCPServer {
  id: string;
  name: string;
  type: 'client' | 'server';
  enabled: boolean;
  priority: number;
  client_config?: MCPClientConfig;
  server_config?: MCPServerConfig;
}

export interface MCPConfig {
  auto_start: boolean;
  tool_discovery_interval: number;
  max_concurrent_calls: number;
  servers: MCPServer[];
}

export interface MCPServerStatus {
  id: string;
  name: string;
  type: 'client' | 'server';
  enabled: boolean;
  connected: boolean;
  last_error?: string;
  tool_count: number;
  last_update: string;
}

export interface MCPTool {
  name: string;
  description: string;
  server_id: string;
  server_name: string;
  namespace: string;
  input_schema: {
    type: string;
    properties: Record<string, any>;
  };
}

export interface MCPConnectionTestResult {
  success: boolean;
  error?: string;
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

  // ==================== MCP 配置管理 API ====================

  async getMCPConfig(): Promise<MCPConfig> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      // 暂时返回模拟数据
      return this.getMockMCPConfig();
    } catch (error) {
      console.error('Failed to get MCP config:', error);
      return this.getMockMCPConfig();
    }
  }

  async updateMCPConfig(config: MCPConfig): Promise<boolean> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      console.log('Updating MCP config:', config);
      return true;
    } catch (error) {
      console.error('Failed to update MCP config:', error);
      return false;
    }
  }

  async addMCPServer(server: MCPServer): Promise<MCPServer> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      console.log('Adding MCP server:', server);
      return server;
    } catch (error) {
      console.error('Failed to add MCP server:', error);
      throw error;
    }
  }

  async removeMCPServer(serverId: string): Promise<boolean> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      console.log('Removing MCP server:', serverId);
      return true;
    } catch (error) {
      console.error('Failed to remove MCP server:', error);
      return false;
    }
  }

  async updateMCPServer(serverId: string, server: MCPServer): Promise<boolean> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      console.log('Updating MCP server:', serverId, server);
      return true;
    } catch (error) {
      console.error('Failed to update MCP server:', error);
      return false;
    }
  }

  async getMCPServerStatus(serverId: string): Promise<MCPServerStatus> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      return this.getMockMCPServerStatus(serverId);
    } catch (error) {
      console.error('Failed to get MCP server status:', error);
      return this.getMockMCPServerStatus(serverId);
    }
  }

  async getMCPTools(): Promise<MCPTool[]> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      return this.getMockMCPTools();
    } catch (error) {
      console.error('Failed to get MCP tools:', error);
      return this.getMockMCPTools();
    }
  }

  async testMCPServerConnection(serverId: string): Promise<MCPConnectionTestResult> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      console.log('Testing MCP server connection:', serverId);
      return { success: true };
    } catch (error) {
      console.error('Failed to test MCP server connection:', error);
      return { success: false, error: error instanceof Error ? error.message : 'Unknown error' };
    }
  }

  // ==================== MCP 模拟数据 ====================

  private getMockMCPConfig(): MCPConfig {
    return {
      auto_start: true,
      tool_discovery_interval: 60,
      max_concurrent_calls: 10,
      servers: [
        {
          id: 'local-fs-tools',
          name: '本地文件系统工具',
          type: 'server',
          enabled: true,
          priority: 100,
          server_config: {
            type: 'stdio',
            tool_filter: '',
          },
        },
        {
          id: 'local-git-tools',
          name: '本地Git工具',
          type: 'server',
          enabled: true,
          priority: 90,
          server_config: {
            type: 'stdio',
            tool_filter: '',
          },
        },
        {
          id: 'weather-api',
          name: '天气API',
          type: 'client',
          enabled: true,
          priority: 80,
          client_config: {
            server_type: 'stdio',
            command: 'npx',
            args: ['@modelcontextprotocol/server-weather'],
            namespace: 'weather',
            timeout: 30,
          },
        },
        {
          id: 'filesystem-browser',
          name: '文件系统浏览器',
          type: 'client',
          enabled: true,
          priority: 70,
          client_config: {
            server_type: 'stdio',
            command: 'npx',
            args: ['@modelcontextprotocol/server-filesystem'],
            namespace: 'fs',
            timeout: 30,
          },
        },
      ],
    };
  }

  private getMockMCPServerStatus(serverId: string): MCPServerStatus {
    const now = new Date();
    const servers = this.getMockMCPConfig().servers;
    const server = servers.find(s => s.id === serverId) || servers[0];

    return {
      id: server.id,
      name: server.name,
      type: server.type,
      enabled: server.enabled,
      connected: true,
      last_error: '',
      tool_count: server.type === 'client' ? 5 : 8,
      last_update: now.toISOString(),
    };
  }

  private getMockMCPTools(): MCPTool[] {
    return [
      {
        name: 'read_file',
        description: '读取文件内容',
        server_id: 'local-fs-tools',
        server_name: '本地文件系统工具',
        namespace: 'fs',
        input_schema: {
          type: 'object',
          properties: {
            path: {
              type: 'string',
              description: '文件路径',
            },
          },
        },
      },
      {
        name: 'write_file',
        description: '写入文件内容',
        server_id: 'local-fs-tools',
        server_name: '本地文件系统工具',
        namespace: 'fs',
        input_schema: {
          type: 'object',
          properties: {
            path: {
              type: 'string',
              description: '文件路径',
            },
            content: {
              type: 'string',
              description: '文件内容',
            },
          },
        },
      },
      {
        name: 'get_current_weather',
        description: '获取当前天气',
        server_id: 'weather-api',
        server_name: '天气API',
        namespace: 'weather',
        input_schema: {
          type: 'object',
          properties: {
            location: {
              type: 'string',
              description: '地理位置',
            },
            unit: {
              type: 'string',
              description: '温度单位 (celsius/fahrenheit)',
              enum: ['celsius', 'fahrenheit'],
            },
          },
        },
      },
      {
        name: 'list_files',
        description: '列出目录中的文件',
        server_id: 'filesystem-browser',
        server_name: '文件系统浏览器',
        namespace: 'fs',
        input_schema: {
          type: 'object',
          properties: {
            path: {
              type: 'string',
              description: '目录路径',
            },
            recursive: {
              type: 'boolean',
              description: '是否递归',
              default: false,
            },
          },
        },
      },
    ];
  }

  // 监控仪表盘相关方法
  async getAgentsStatus(): Promise<any[]> {
    try {
      // 暂时返回模拟数据，实际应该调用后端API
      return [
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
    } catch (error) {
      console.error('Failed to get agents status:', error);
      return [];
    }
  }

  async getTimelineEvents(): Promise<any[]> {
    try {
      // 暂时返回模拟数据
      return [
        { id: "1", agent: "Orchestrator", action: "init_project", time: "T-14m", width: "25%", offset: "2%", isComplete: true },
        { id: "2", agent: "Orchestrator", action: "analyze_reqs", time: "T-10m", width: "15%", offset: "28%", isComplete: false },
        { id: "3", agent: "Architect", action: "design_db_schema", time: "T-8m", width: "30%", offset: "28%", isComplete: true },
        { id: "4", agent: "Architect", action: "auth_flow", time: "T-4m", width: "20%", offset: "60%", isComplete: false },
      ];
    } catch (error) {
      console.error('Failed to get timeline events:', error);
      return [];
    }
  }

  async getThoughts(): Promise<any[]> {
    try {
      // 暂时返回模拟数据
      return [
        {
          id: "1",
          type: "internal",
          content: "为了确保新旧系统的兼容性，我们需要保留旧的 session_id 映射，同时引入基于 JWT 的无状态验证。",
          agent: "Architect",
          timestamp: new Date().toISOString()
        },
        {
          id: "2",
          type: "internal",
          isCode: true,
          filename: "schema.prisma",
          code: `model User {
  id String @id @default(uuid())
  // legacy_token String?
  password_hash String
  refresh_tokens RefreshToken[]
  createdAt DateTime @default(now())
}`,
          agent: "Architect",
          timestamp: new Date().toISOString()
        },
        {
          id: "3",
          type: "formulating",
          content: "考虑到前端是 React SPA，建议使用 HttpOnly cookie 存储 refresh token，内存存储 access token 以防止 XSS 攻击。正在构建最终架构文档...",
          agent: "Architect",
          timestamp: new Date().toISOString()
        }
      ];
    } catch (error) {
      console.error('Failed to get thoughts:', error);
      return [];
    }
  }

  async getToolLogs(): Promise<any[]> {
    try {
      // 暂时返回模拟数据
      return [
        {
          id: "1",
          toolName: "fs_read",
          toolType: "fs",
          agent: "Orchestrator",
          timestamp: new Date().toISOString(),
          duration: 120,
          status: "success",
          parameters: { path: "/project/config.json" },
          result: { content: "{}" }
        },
        {
          id: "2",
          toolName: "git_init",
          toolType: "git",
          agent: "Orchestrator",
          timestamp: new Date().toISOString(),
          duration: 450,
          status: "success",
          parameters: { path: "/project" },
          result: { initialized: true }
        }
      ];
    } catch (error) {
      console.error('Failed to get tool logs:', error);
      return [];
    }
  }

  async getTerminalOutput(): Promise<string[]> {
    try {
      // 暂时返回模拟数据
      return [
        "[10:42:01] INFO: Orchestrator: Spawned new task cluster 'auth_module'",
        "[10:42:03] INFO: MessageBus: Routed prompt to Architect ID-4A9B",
        "[10:42:05] DEBUG: Architect: Analyzing legacy codebase for auth dependencies...",
        "[10:42:12] SUCCESS: Architect: Legacy analysis complete. Found 3 deprecation warnings.",
        "[10:42:15] DEBUG: Architect: Generating initial schema proposal...",
        "[10:42:28] SYS: Memory allocation at 64%",
        "[10:42:35] DEBUG: Architect: Evaluating OAuth2 vs JWT implementation paths..."
      ];
    } catch (error) {
      console.error('Failed to get terminal output:', error);
      return [];
    }
  }

  async getFileDiffs(): Promise<any[]> {
    try {
      // 暂时返回模拟数据
      return [
        {
          id: "1",
          filePath: "src/models/user.go",
          changeType: "added",
          agent: "Backend Dev",
          timestamp: new Date().toISOString(),
          diff: `+package models
+
++type User struct {
++    ID        string    \`json:"id"\`
++    Email     string    \`json:"email"\`
++    CreatedAt time.Time \`json:"created_at"\`
++}`,
          linesAdded: 7,
          linesDeleted: 0,
          size: 256
        }
      ];
    } catch (error) {
      console.error('Failed to get file diffs:', error);
      return [];
    }
  }

  async getCostMetrics(): Promise<any[]> {
    try {
      // 暂时返回模拟数据
      return [
        {
          id: "1",
          costType: "llm",
          agent: "Architect",
          timestamp: new Date().toISOString(),
          amount: 0.25,
          tokenCount: 1250,
          description: "Architecture design generation"
        },
        {
          id: "2",
          costType: "api",
          agent: "Orchestrator",
          timestamp: new Date().toISOString(),
          amount: 0.05,
          description: "Git API calls"
        }
      ];
    } catch (error) {
      console.error('Failed to get cost metrics:', error);
      return [];
    }
  }

  async getEvolutionEvents(): Promise<any[]> {
    try {
      // 暂时返回模拟数据
      return [
        {
          id: "1",
          eventType: "skill_generated",
          agent: "Architect",
          timestamp: new Date().toISOString(),
          description: "Generated skill for JWT authentication implementation",
          efficiencyGain: 15,
          tokenSavings: 5000,
          timeSavings: 30
        }
      ];
    } catch (error) {
      console.error('Failed to get evolution events:', error);
      return [];
    }
  }

  async getMonitoringStats(): Promise<any> {
    try {
      // 暂时返回模拟数据
      return {
        activeAgents: 2,
        totalAgents: 4,
        successRate: 85,
        uptime: 99,
        totalExecutions: 1248,
        avgResponseTime: 2.4,
        tokenUsage: 12800,
        toolCalls: 42,
        filesChanged: 8,
        costToday: 12.50,
        skillsGenerated: 3
      };
    } catch (error) {
      console.error('Failed to get monitoring stats:', error);
      return {
        activeAgents: 0,
        totalAgents: 0,
        successRate: 0,
        uptime: 0,
        totalExecutions: 0,
        avgResponseTime: 0,
        tokenUsage: 0,
        toolCalls: 0,
        filesChanged: 0,
        costToday: 0,
        skillsGenerated: 0
      };
    }
  }

  async sendIntervention(message: string): Promise<void> {
    try {
      console.log('Sending intervention:', message);
      // 暂时只是记录，实际应该调用后端API
      await this.logMessage('info', `User intervention: ${message}`);
    } catch (error) {
      console.error('Failed to send intervention:', error);
    }
  }

  // ==================== 历史记录与回放 API ====================

  async getHistoryRecords(query: any): Promise<[any[], number]> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      // 暂时返回模拟数据
      return this.getMockHistoryRecords(query);
    } catch (error) {
      console.error('Failed to get history records:', error);
      return [[], 0];
    }
  }

  async getHistoryStatistics(): Promise<any> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      // 暂时返回模拟数据
      return this.getMockHistoryStatistics();
    } catch (error) {
      console.error('Failed to get history statistics:', error);
      return {};
    }
  }

  async deleteHistoryRecord(recordId: string): Promise<boolean> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      console.log('Deleting history record:', recordId);
      return true;
    } catch (error) {
      console.error('Failed to delete history record:', error);
      return false;
    }
  }

  async startReplay(recordId: string, speed: number = 1.0): Promise<any> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      console.log('Starting replay:', recordId, speed);
      return this.getMockReplayState(recordId);
    } catch (error) {
      console.error('Failed to start replay:', error);
      throw error;
    }
  }

  async controlReplay(replayId: string, action: string, params?: any): Promise<any> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      console.log('Controlling replay:', replayId, action, params);
      return this.getMockReplayState(replayId);
    } catch (error) {
      console.error('Failed to control replay:', error);
      throw error;
    }
  }

  async getReplayState(replayId: string): Promise<any> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      // 暂时返回模拟数据
      return this.getMockReplayState(replayId);
    } catch (error) {
      console.error('Failed to get replay state:', error);
      throw error;
    }
  }

  async getReplayData(replayId: string, fromIndex: number, limit: number): Promise<any[]> {
    try {
      // 注意：这里需要等待Wails生成新的绑定
      // 暂时返回模拟数据
      return this.getMockReplayData(fromIndex, limit);
    } catch (error) {
      console.error('Failed to get replay data:', error);
      return [];
    }
  }

  subscribeToReplayEvents(replayId: string, callback: (event: string, data: any) => void): () => void {
    // 模拟事件订阅
    console.log('Subscribing to replay events:', replayId);
    
    const interval = setInterval(() => {
      const state = this.getMockReplayState(replayId);
      if (state.progress < 100) {
        state.progress += 1;
        state.current_span_index = Math.floor(state.progress / 100 * state.total_spans);
        callback('replay:progress', { state });
      } else {
        state.status = 'completed';
        callback('replay:ended', { state });
        clearInterval(interval);
      }
    }, 1000);

    return () => {
      clearInterval(interval);
      console.log('Unsubscribed from replay events:', replayId);
    };
  }

  // ==================== 历史记录模拟数据 ====================

  private getMockHistoryRecords(query: any): [any[], number] {
    const now = new Date();
    const records = [
      {
        id: 'exec-001',
        trace_id: 'trace-001',
        pipeline_id: 'pipeline-001',
        name: '用户认证系统开发',
        description: '完整的JWT认证系统实现，包含前端React组件和后端Go API',
        status: 'completed',
        start_time: new Date(now.getTime() - 2 * 60 * 60 * 1000).toISOString(),
        end_time: new Date(now.getTime() - 1.5 * 60 * 60 * 1000).toISOString(),
        duration: 30 * 60 * 1000, // 30分钟
        agent_count: 4,
        tool_call_count: 42,
        token_usage: 12500,
        cost_estimate: 0.85,
        metadata: {
          primary_agent: 'Orchestrator',
          project_type: 'web',
          framework: 'react+go'
        },
        tags: ['production', 'authentication', 'security'],
        created_at: new Date(now.getTime() - 2 * 60 * 60 * 1000).toISOString(),
        updated_at: new Date(now.getTime() - 1.5 * 60 * 60 * 1000).toISOString()
      },
      {
        id: 'exec-002',
        trace_id: 'trace-002',
        pipeline_id: 'pipeline-002',
        name: '电商产品页面优化',
        description: '优化产品展示页面性能，添加图片懒加载和缓存策略',
        status: 'completed',
        start_time: new Date(now.getTime() - 5 * 60 * 60 * 1000).toISOString(),
        end_time: new Date(now.getTime() - 4 * 60 * 60 * 1000).toISOString(),
        duration: 60 * 60 * 1000, // 1小时
        agent_count: 3,
        tool_call_count: 28,
        token_usage: 8500,
        cost_estimate: 0.55,
        metadata: {
          primary_agent: 'Frontend Dev',
          project_type: 'ecommerce',
          framework: 'react'
        },
        tags: ['optimization', 'performance', 'frontend'],
        created_at: new Date(now.getTime() - 5 * 60 * 60 * 1000).toISOString(),
        updated_at: new Date(now.getTime() - 4 * 60 * 60 * 1000).toISOString()
      },
      {
        id: 'exec-003',
        trace_id: 'trace-003',
        pipeline_id: 'pipeline-003',
        name: '数据库迁移脚本生成',
        description: '从MySQL迁移到PostgreSQL的自动化脚本生成',
        status: 'failed',
        start_time: new Date(now.getTime() - 8 * 60 * 60 * 1000).toISOString(),
        end_time: new Date(now.getTime() - 7.5 * 60 * 60 * 1000).toISOString(),
        duration: 30 * 60 * 1000, // 30分钟
        agent_count: 2,
        tool_call_count: 15,
        token_usage: 5200,
        cost_estimate: 0.35,
        metadata: {
          primary_agent: 'Backend Dev',
          project_type: 'database',
          error: '数据类型转换失败'
        },
        tags: ['database', 'migration', 'failed'],
        created_at: new Date(now.getTime() - 8 * 60 * 60 * 1000).toISOString(),
        updated_at: new Date(now.getTime() - 7.5 * 60 * 60 * 1000).toISOString()
      },
      {
        id: 'exec-004',
        trace_id: 'trace-004',
        pipeline_id: 'pipeline-004',
        name: 'API文档自动生成',
        description: '基于现有Go代码生成OpenAPI 3.0规范文档',
        status: 'running',
        start_time: new Date(now.getTime() - 0.5 * 60 * 60 * 1000).toISOString(),
        end_time: undefined,
        duration: undefined,
        agent_count: 2,
        tool_call_count: 8,
        token_usage: 3200,
        cost_estimate: 0.22,
        metadata: {
          primary_agent: 'Backend Dev',
          project_type: 'documentation'
        },
        tags: ['documentation', 'api', 'automation'],
        created_at: new Date(now.getTime() - 0.5 * 60 * 60 * 1000).toISOString(),
        updated_at: new Date(now.getTime() - 0.5 * 60 * 60 * 1000).toISOString()
      },
      {
        id: 'exec-005',
        trace_id: 'trace-005',
        pipeline_id: 'pipeline-005',
        name: '移动端响应式布局',
        description: '为现有Web应用添加移动端适配',
        status: 'cancelled',
        start_time: new Date(now.getTime() - 12 * 60 * 60 * 1000).toISOString(),
        end_time: new Date(now.getTime() - 11.5 * 60 * 60 * 1000).toISOString(),
        duration: 30 * 60 * 1000, // 30分钟
        agent_count: 1,
        tool_call_count: 5,
        token_usage: 1800,
        cost_estimate: 0.12,
        metadata: {
          primary_agent: 'Frontend Dev',
          project_type: 'mobile',
          reason: '需求变更'
        },
        tags: ['mobile', 'responsive', 'cancelled'],
        created_at: new Date(now.getTime() - 12 * 60 * 60 * 1000).toISOString(),
        updated_at: new Date(now.getTime() - 11.5 * 60 * 60 * 1000).toISOString()
      }
    ];

    // 简单的查询过滤
    let filtered = records;
    
    if (query.searchText) {
      const search = query.searchText.toLowerCase();
      filtered = filtered.filter(r => 
        r.name.toLowerCase().includes(search) || 
        r.description?.toLowerCase().includes(search)
      );
    }
    
    if (query.statuses && query.statuses.length > 0) {
      filtered = filtered.filter(r => query.statuses.includes(r.status));
    }
    
    if (query.tags && query.tags.length > 0) {
      filtered = filtered.filter(r => 
        r.tags?.some(tag => query.tags.includes(tag))
      );
    }

    // 排序
    filtered.sort((a, b) => {
      const order = query.sortOrder === 'asc' ? 1 : -1;
      switch (query.sortBy) {
        case 'start_time':
          return (new Date(a.start_time).getTime() - new Date(b.start_time).getTime()) * order;
        case 'end_time':
          const aEnd = a.end_time ? new Date(a.end_time).getTime() : 0;
          const bEnd = b.end_time ? new Date(b.end_time).getTime() : 0;
          return (aEnd - bEnd) * order;
        case 'duration':
          return ((a.duration || 0) - (b.duration || 0)) * order;
        case 'agent_count':
          return (a.agent_count - b.agent_count) * order;
        case 'token_usage':
          return (a.token_usage - b.token_usage) * order;
        case 'cost_estimate':
          return ((a.cost_estimate || 0) - (b.cost_estimate || 0)) * order;
        default:
          return (new Date(a.start_time).getTime() - new Date(b.start_time).getTime()) * order;
      }
    });

    const total = filtered.length;
    const start = query.offset || 0;
    const limit = query.limit || 20;
    const paginated = filtered.slice(start, start + limit);

    return [paginated, total];
  }

  private getMockHistoryStatistics(): any {
    return {
      total_records: 5,
      completed_records: 2,
      failed_records: 1,
      success_rate: 40.0,
      total_duration: '2.5h',
      avg_duration: 30 * 60 * 1000, // 30分钟
      total_token_usage: 31200,
      avg_token_usage: 6240,
      total_tool_calls: 98,
      avg_tool_calls: 19.6,
      total_cost: 2.09,
      avg_cost: 0.418,
      agent_statistics: {
        'Orchestrator': {
          count: 1,
          total_duration: 30 * 60 * 1000,
          success_count: 1,
          avg_duration: 30 * 60 * 1000,
          success_rate: 100.0
        },
        'Frontend Dev': {
          count: 2,
          total_duration: 90 * 60 * 1000,
          success_count: 1,
          avg_duration: 45 * 60 * 1000,
          success_rate: 50.0
        },
        'Backend Dev': {
          count: 2,
          total_duration: 60 * 60 * 1000,
          success_count: 1,
          avg_duration: 30 * 60 * 1000,
          success_rate: 50.0
        }
      }
    };
  }

  private getMockReplayState(replayId: string): any {
    const now = new Date();
    return {
      record_id: replayId,
      trace_id: 'trace-001',
      status: 'playing',
      current_time: new Date(now.getTime() - 1.8 * 60 * 60 * 1000).toISOString(),
      start_time: new Date(now.getTime() - 2 * 60 * 60 * 1000).toISOString(),
      end_time: new Date(now.getTime() - 1.5 * 60 * 60 * 1000).toISOString(),
      playback_speed: 1.0,
      current_span_index: 15,
      total_spans: 25,
      progress: 60.0,
      metadata: {
        record_name: '用户认证系统开发',
        record_status: 'completed'
      }
    };
  }

  private getMockReplayData(fromIndex: number, limit: number): any[] {
    const spans = [];
    const now = new Date();
    
    for (let i = fromIndex; i < fromIndex + limit && i < 25; i++) {
      const startTime = new Date(now.getTime() - (2 - i * 0.02) * 60 * 60 * 1000);
      const endTime = new Date(startTime.getTime() + 60 * 1000); // 1分钟
      
      spans.push({
        id: `span-${i}`,
        trace_id: 'trace-001',
        parent_id: i > 0 ? `span-${Math.floor(i/2)}` : undefined,
        name: i % 5 === 0 ? 'Agent Execution' : 'Tool Invocation',
        kind: i % 5 === 0 ? 'INTERNAL' : 'CLIENT',
        start_time: startTime.toISOString(),
        end_time: endTime.toISOString(),
        duration: 60 * 1000,
        attributes: {
          'agent.name': i % 5 === 0 ? (i < 5 ? 'Orchestrator' : i < 10 ? 'Architect' : i < 15 ? 'Backend Dev' : 'Frontend Dev') : undefined,
          'tool.name': i % 5 !== 0 ? `tool-${i % 3}` : undefined,
          'task.id': `task-${i}`,
          'iteration': i
        },
        events: i % 3 === 0 ? [
          {
            name: 'task_started',
            timestamp: startTime.toISOString(),
            attributes: { 'task.type': 'code_generation' }
          },
          {
            name: 'task_completed',
            timestamp: endTime.toISOString(),
            attributes: { 'result': 'success' }
          }
        ] : [],
        status: i === 12 ? 'ERROR' : 'OK',
        status_message: i === 12 ? 'Failed to generate code: timeout' : undefined,
        resource: {
          'service.name': 'genpulse',
          'service.version': '1.0.0'
        },
        agent_name: i % 5 === 0 ? (i < 5 ? 'Orchestrator' : i < 10 ? 'Architect' : i < 15 ? 'Backend Dev' : 'Frontend Dev') : undefined,
        tool_name: i % 5 !== 0 ? `tool-${i % 3}` : undefined,
        pipeline_id: 'pipeline-001'
      });
    }
    
    return spans;
  }
}

export const api = new ApiService();