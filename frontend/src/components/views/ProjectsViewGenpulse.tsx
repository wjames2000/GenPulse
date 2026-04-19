import React, { useState } from 'react';

type Project = {
  id: string;
  name: string;
  description: string;
  type: 'web' | 'cli' | 'api' | 'library';
  status: 'running' | 'completed' | 'paused' | 'error';
  agents: number;
  lastUpdated: string;
  techStack: string[];
  progress: number;
};

type ProjectFormData = {
  name: string;
  description: string;
  projectType: string;
  projectPath: string;
  techStack: {
    frontend: string;
    backend: string;
    database: string;
    ui: string;
  };
  agents: {
    productManager: boolean;
    architect: boolean;
    frontendDev: boolean;
    backendDev: boolean;
    qaEngineer: boolean;
    devopsEngineer: boolean;
    reviewer: boolean;
  };
};

const ProjectsViewGenpulse: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'list' | 'new'>('list');
  const [selectedProject, setSelectedProject] = useState<string | null>(null);
  const [formData, setFormData] = useState<ProjectFormData>({
    name: '',
    description: '',
    projectType: 'web',
    projectPath: './projects/',
    techStack: {
      frontend: 'react',
      backend: 'go',
      database: 'sqlite',
      ui: 'shadcn/ui'
    },
    agents: {
      productManager: true,
      architect: true,
      frontendDev: true,
      backendDev: true,
      qaEngineer: false,
      devopsEngineer: false,
      reviewer: true
    }
  });

  const projects: Project[] = [
    {
      id: '1',
      name: 'E-Commerce Platform',
      description: 'Full-stack e-commerce solution with AI-powered recommendations',
      type: 'web',
      status: 'running',
      agents: 4,
      lastUpdated: '2小时前',
      techStack: ['React', 'Node.js', 'PostgreSQL', 'Tailwind'],
      progress: 65
    },
    {
      id: '2',
      name: 'API Gateway',
      description: 'Microservices API gateway with rate limiting and authentication',
      type: 'api',
      status: 'completed',
      agents: 3,
      lastUpdated: '1天前',
      techStack: ['Go', 'Redis', 'JWT'],
      progress: 100
    },
    {
      id: '3',
      name: 'CLI Tool',
      description: 'Developer productivity CLI for project scaffolding',
      type: 'cli',
      status: 'running',
      agents: 2,
      lastUpdated: '3小时前',
      techStack: ['Go', 'Cobra', 'Viper'],
      progress: 30
    },
    {
      id: '4',
      name: 'UI Component Library',
      description: 'Design system and reusable React components',
      type: 'library',
      status: 'paused',
      agents: 1,
      lastUpdated: '1周前',
      techStack: ['React', 'TypeScript', 'Storybook'],
      progress: 45
    },
    {
      id: '5',
      name: 'Data Analytics Dashboard',
      description: 'Real-time analytics dashboard with ML insights',
      type: 'web',
      status: 'error',
      agents: 5,
      lastUpdated: '刚刚',
      techStack: ['Next.js', 'Python', 'MySQL', 'Chart.js'],
      progress: 80
    },
    {
      id: '6',
      name: 'Authentication Service',
      description: 'Centralized auth service with OAuth2 and SSO',
      type: 'api',
      status: 'running',
      agents: 3,
      lastUpdated: '5小时前',
      techStack: ['Go', 'OAuth2', 'PostgreSQL'],
      progress: 90
    }
  ];

  const getStatusColor = (status: Project['status']) => {
    switch (status) {
      case 'running': return 'bg-running/10 text-running border border-running/20';
      case 'completed': return 'bg-success/10 text-success border border-success/20';
      case 'paused': return 'bg-warning/10 text-warning border border-warning/20';
      case 'error': return 'bg-error/10 text-error border border-error/20';
    }
  };

  const getStatusIcon = (status: Project['status']) => {
    switch (status) {
      case 'running': return 'play_circle';
      case 'completed': return 'check_circle';
      case 'paused': return 'pause_circle';
      case 'error': return 'error';
    }
  };

  const getTypeColor = (type: Project['type']) => {
    switch (type) {
      case 'web': return 'bg-blue-500/10 text-blue-400';
      case 'cli': return 'bg-green-500/10 text-green-400';
      case 'api': return 'bg-purple-500/10 text-purple-400';
      case 'library': return 'bg-orange-500/10 text-orange-400';
    }
  };

  const handleInputChange = (field: keyof ProjectFormData, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleTechStackChange = (field: keyof ProjectFormData['techStack'], value: string) => {
    setFormData(prev => ({
      ...prev,
      techStack: { ...prev.techStack, [field]: value }
    }));
  };

  const handleAgentToggle = (agent: keyof ProjectFormData['agents']) => {
    setFormData(prev => ({
      ...prev,
      agents: { ...prev.agents, [agent]: !prev.agents[agent] }
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log('Form submitted:', formData);
    setActiveTab('list');
  };

  if (activeTab === 'new') {
    return (
      <div className="h-full overflow-y-auto p-6">
        {/* Page Header */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-on-surface mb-2">初始化流水线</h1>
              <p className="text-sm text-on-surface-variant">
                Configure basic parameters and autonomous agents for the new initiative.
              </p>
            </div>
            <div className="flex gap-3">
              <button
                onClick={() => setActiveTab('list')}
                className="px-4 py-2 rounded-lg text-sm font-medium text-on-surface-variant hover:bg-surface-bright transition-colors"
              >
                取消
              </button>
              <button
                onClick={handleSubmit}
                className="px-5 py-2 rounded-lg text-sm font-semibold bg-gradient-to-br from-primary-container to-inverse-primary text-on-primary-container hover:brightness-110 transition-all flex items-center gap-2"
              >
                <span className="material-symbols-outlined text-sm">precision_manufacturing</span>
                初始化系统
              </button>
            </div>
          </div>
        </div>

        {/* Form Canvas */}
        <div className="bg-surface-container-low rounded-xl p-8 space-y-12 max-w-5xl mx-auto">
          {/* Section: Basic Info (设计规范 5.3.1) */}
          <section className="space-y-6">
            <div className="flex items-center gap-3 border-b border-outline-variant/15 pb-2">
              <span className="material-symbols-outlined text-primary">data_object</span>
              <h2 className="text-lg font-semibold text-on-surface">基本信息配置</h2>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              <div className="space-y-2">
                <label className="text-[0.6875rem] font-medium uppercase tracking-wider text-on-surface-variant block">
                  项目名称
                </label>
                <input
                  className="w-full bg-surface-container-lowest border-none rounded-lg px-4 py-2.5 text-sm text-on-surface font-mono focus:ring-0 focus:border-b-2 focus:border-primary outline-none transition-all"
                  type="text"
                  placeholder="my-project"
                  value={formData.name}
                  onChange={(e) => handleInputChange('name', e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <label className="text-[0.6875rem] font-medium uppercase tracking-wider text-on-surface-variant block">
                  项目类型
                </label>
                <select
                  className="w-full bg-surface-container-lowest border-none rounded-lg px-4 py-2.5 text-sm text-on-surface font-mono focus:ring-0 focus:border-b-2 focus:border-primary outline-none transition-all appearance-none"
                  value={formData.projectType}
                  onChange={(e) => handleInputChange('projectType', e.target.value)}
                >
                  <option value="web">Web 应用</option>
                  <option value="cli">CLI 工具</option>
                  <option value="api">API 服务</option>
                  <option value="library">库</option>
                </select>
              </div>
              <div className="space-y-2">
                <label className="text-[0.6875rem] font-medium uppercase tracking-wider text-on-surface-variant block">
                  项目路径
                </label>
                <div className="flex bg-surface-container-lowest rounded-lg overflow-hidden focus-within:border-b-2 focus-within:border-primary transition-all">
                  <span className="px-4 py-2.5 text-outline text-sm font-mono border-r border-outline-variant/20 bg-surface-container-highest/30">
                    ~/projects/
                  </span>
                  <input
                    className="flex-1 bg-transparent border-none px-4 py-2.5 text-sm text-on-surface font-mono outline-none"
                    type="text"
                    placeholder="project-name"
                    value={formData.projectPath}
                    onChange={(e) => handleInputChange('projectPath', e.target.value)}
                  />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-[0.6875rem] font-medium uppercase tracking-wider text-on-surface-variant block">
                  需求描述
                </label>
                <textarea
                  className="w-full bg-surface-container-lowest border-none rounded-lg px-4 py-2.5 text-sm text-on-surface focus:ring-0 focus:border-b-2 focus:border-primary outline-none transition-all min-h-[120px] resize-y"
                  placeholder="Describe your project requirements..."
                  value={formData.description}
                  onChange={(e) => handleInputChange('description', e.target.value)}
                />
                <p className="text-xs text-outline mt-1">支持 Markdown 格式</p>
              </div>
            </div>
          </section>

          {/* Section: Tech Stack (设计规范 5.3.2) */}
          <section className="space-y-6">
            <div className="flex items-center gap-3 border-b border-outline-variant/15 pb-2">
              <span className="material-symbols-outlined text-primary">deployed_code</span>
              <h2 className="text-lg font-semibold text-on-surface">技术栈配置</h2>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium text-on-surface-variant">前端框架</label>
                <select
                  className="w-full bg-surface-container-lowest border-none rounded-lg px-3 py-2 text-sm text-on-surface focus:ring-0 focus:border-b-2 focus:border-primary outline-none transition-all"
                  value={formData.techStack.frontend}
                  onChange={(e) => handleTechStackChange('frontend', e.target.value)}
                >
                  <option value="react">React</option>
                  <option value="vue">Vue.js</option>
                  <option value="nextjs">Next.js</option>
                  <option value="none">无</option>
                </select>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-on-surface-variant">后端框架</label>
                <select
                  className="w-full bg-surface-container-lowest border-none rounded-lg px-3 py-2 text-sm text-on-surface focus:ring-0 focus:border-b-2 focus:border-primary outline-none transition-all"
                  value={formData.techStack.backend}
                  onChange={(e) => handleTechStackChange('backend', e.target.value)}
                >
                  <option value="go">Go</option>
                  <option value="express">Express.js</option>
                  <option value="fastify">Fastify</option>
                  <option value="none">无</option>
                </select>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-on-surface-variant">数据库</label>
                <select
                  className="w-full bg-surface-container-lowest border-none rounded-lg px-3 py-2 text-sm text-on-surface focus:ring-0 focus:border-b-2 focus:border-primary outline-none transition-all"
                  value={formData.techStack.database}
                  onChange={(e) => handleTechStackChange('database', e.target.value)}
                >
                  <option value="sqlite">SQLite</option>
                  <option value="postgresql">PostgreSQL</option>
                  <option value="mysql">MySQL</option>
                  <option value="none">无</option>
                </select>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-on-surface-variant">UI 组件库</label>
                <select
                  className="w-full bg-surface-container-lowest border-none rounded-lg px-3 py-2 text-sm text-on-surface focus:ring-0 focus:border-b-2 focus:border-primary outline-none transition-all"
                  value={formData.techStack.ui}
                  onChange={(e) => handleTechStackChange('ui', e.target.value)}
                >
                  <option value="shadcn/ui">shadcn/ui</option>
                  <option value="tailwind">Tailwind</option>
                  <option value="none">无</option>
                </select>
              </div>
            </div>
          </section>

          {/* Section: Agent Team (设计规范 5.3.3) */}
          <section className="space-y-6">
            <div className="flex items-center gap-3 border-b border-outline-variant/15 pb-2">
              <span className="material-symbols-outlined text-primary">groups</span>
              <h2 className="text-lg font-semibold text-on-surface">Agent 团队配置</h2>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              {[
                { key: 'productManager' as const, label: '产品经理', color: 'agent-product-manager' },
                { key: 'architect' as const, label: '架构师', color: 'agent-architect' },
                { key: 'frontendDev' as const, label: '前端开发', color: 'agent-frontend' },
                { key: 'backendDev' as const, label: '后端开发', color: 'agent-backend' },
                { key: 'qaEngineer' as const, label: '测试工程师', color: 'agent-qa' },
                { key: 'devopsEngineer' as const, label: '运维工程师', color: 'agent-devops' },
                { key: 'reviewer' as const, label: '审查员', color: 'agent-reviewer' },
              ].map((agent) => (
                <div
                  key={agent.key}
                  className={`bg-surface-container rounded-lg p-4 cursor-pointer transition-all ${
                    formData.agents[agent.key] ? 'ring-2 ring-primary/30' : ''
                  }`}
                  onClick={() => handleAgentToggle(agent.key)}
                >
                  <div className="flex items-center justify-between mb-3">
                    <div className={`w-8 h-8 rounded-lg flex items-center justify-center ${
                      formData.agents[agent.key] ? `bg-${agent.color}/10` : 'bg-surface-container-high'
                    }`}>
                      <span className={`material-symbols-outlined text-lg ${
                        formData.agents[agent.key] ? `text-${agent.color}` : 'text-outline'
                      }`}>
                        smart_toy
                      </span>
                    </div>
                    <div className={`w-10 h-5 rounded-full transition-all ${
                      formData.agents[agent.key] ? 'bg-primary' : 'bg-surface-container-high'
                    }`}>
                      <div className={`w-5 h-5 rounded-full bg-white transition-all ${
                        formData.agents[agent.key] ? 'translate-x-5' : ''
                      }`}></div>
                    </div>
                  </div>
                  <h3 className="text-sm font-medium text-on-surface">{agent.label}</h3>
                  <p className="text-xs text-outline mt-1">
                    {formData.agents[agent.key] ? '已启用' : '未启用'}
                  </p>
                </div>
              ))}
            </div>
          </section>

          {/* Configuration Template (设计规范 5.3.6) */}
          <section className="space-y-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span className="material-symbols-outlined text-primary">content_copy</span>
                <h2 className="text-lg font-semibold text-on-surface">配置模板</h2>
              </div>
              <div className="flex gap-3">
                <button className="px-3 py-1.5 text-sm font-medium text-primary hover:bg-primary/10 rounded-lg transition-colors">
                  保存为模板
                </button>
                <button className="px-3 py-1.5 text-sm font-medium text-primary hover:bg-primary/10 rounded-lg transition-colors">
                  从模板加载
                </button>
              </div>
            </div>
          </section>
        </div>
      </div>
    );
  }

  // Project List View
  return (
    <div className="h-full overflow-y-auto p-6">
      {/* Page Header */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-on-surface mb-2">项目流水线</h1>
            <p className="text-sm text-on-surface-variant">
              管理您的 AI 项目开发流水线，监控执行状态
            </p>
          </div>
          <button
            onClick={() => setActiveTab('new')}
            className="bg-gradient-to-br from-primary-container to-inverse-primary text-on-primary-container hover:brightness-110 transition-all duration-200 py-2 px-4 rounded-lg flex items-center justify-center gap-2 font-medium text-sm"
          >
            <span className="material-symbols-outlined text-sm">add</span>
            新建流水线
          </button>
        </div>
      </div>

      {/* Project Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {projects.map((project) => (
          <div
            key={project.id}
            className={`bg-surface-container rounded-xl p-5 hover:bg-surface-container-high transition-all duration-200 cursor-pointer ${
              selectedProject === project.id ? 'ring-2 ring-primary/30' : ''
            }`}
            onClick={() => setSelectedProject(project.id)}
          >
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-surface-container-low text-primary">
                  <span className="material-symbols-outlined">folder</span>
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-on-surface">{project.name}</h3>
                  <p className="text-xs text-outline">{project.type.toUpperCase()}</p>
                </div>
              </div>
              <div className={`flex items-center gap-1.5 px-2 py-1 rounded-full text-xs ${getStatusColor(project.status)}`}>
                <span className="material-symbols-outlined text-xs">{getStatusIcon(project.status)}</span>
                <span className="capitalize">{project.status}</span>
              </div>
            </div>

            <p className="text-sm text-on-surface-variant mb-4 line-clamp-2">{project.description}</p>

            <div className="space-y-3">
              <div className="flex items-center justify-between text-sm">
                <span className="text-outline">代理数</span>
                <span className="font-mono text-primary">{project.agents}</span>
              </div>

              {project.status === 'running' && (
                <div className="space-y-1">
                  <div className="flex items-center justify-between text-xs">
                    <span className="text-outline">进度</span>
                    <span className="font-mono text-on-surface">{project.progress}%</span>
                  </div>
                  <div className="h-1.5 bg-surface-container-high rounded-full overflow-hidden">
                    <div
                      className="h-full bg-running rounded-full transition-all duration-500"
                      style={{ width: `${project.progress}%` }}
                    ></div>
                  </div>
                </div>
              )}

              <div className="flex flex-wrap gap-1 mt-2">
                {project.techStack.slice(0, 3).map((tech, idx) => (
                  <span key={idx} className="text-xs px-2 py-1 rounded bg-surface-container-high text-outline">
                    {tech}
                  </span>
                ))}
                {project.techStack.length > 3 && (
                  <span className="text-xs px-2 py-1 rounded bg-surface-container-high text-outline">
                    +{project.techStack.length - 3}
                  </span>
                )}
              </div>

              <div className="flex items-center justify-between pt-3 border-t border-outline-variant/15 text-xs text-outline">
                <span>最后更新: {project.lastUpdated}</span>
                <button className="text-primary hover:text-primary-container transition-colors">
                  查看详情 →
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Empty State */}
      {projects.length === 0 && (
        <div className="text-center py-16">
          <div className="w-20 h-20 rounded-full bg-surface-container mx-auto mb-6 flex items-center justify-center">
            <span className="material-symbols-outlined text-3xl text-outline">folder_open</span>
          </div>
          <h3 className="text-lg font-medium text-on-surface mb-2">暂无项目</h3>
          <p className="text-sm text-on-surface-variant mb-6">创建您的第一个 AI 开发流水线项目</p>
          <button
            onClick={() => setActiveTab('new')}
            className="bg-primary-container text-on-primary-container hover:bg-inverse-primary rounded-full px-6 py-2.5 text-sm font-medium shadow-lg shadow-primary/20 transition-all"
          >
            创建新项目
          </button>
        </div>
      )}
    </div>
  );
};

export default ProjectsViewGenpulse;