import React, { useState } from 'react';

const DashboardViewGenpulse: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'logs' | 'terminal' | 'files'>('logs');

  // 模拟 Agent 数据
  const agents = [
    {
      id: 1,
      name: 'Product Manager',
      role: '产品经理',
      status: 'completed' as const,
      task: '需求分析与用户故事创建',
      progress: 100,
      tokens: 2150,
      color: 'agent-product-manager'
    },
    {
      id: 2,
      name: 'Architect',
      role: '架构师',
      status: 'running' as const,
      task: '设计微服务架构',
      progress: 75,
      tokens: 3480,
      color: 'agent-architect'
    },
    {
      id: 3,
      name: 'Frontend Dev',
      role: '前端开发',
      status: 'running' as const,
      task: '生成 React 组件库',
      progress: 45,
      tokens: 1890,
      color: 'agent-frontend'
    },
    {
      id: 4,
      name: 'Backend Dev',
      role: '后端开发',
      status: 'waiting' as const,
      task: '等待 API 规范',
      progress: 20,
      tokens: 920,
      color: 'agent-backend'
    },
    {
      id: 5,
      name: 'QA Engineer',
      role: '测试工程师',
      status: 'idle' as const,
      task: '待命',
      progress: 0,
      tokens: 0,
      color: 'agent-qa'
    },
    {
      id: 6,
      name: 'DevOps Engineer',
      role: '运维工程师',
      status: 'running' as const,
      task: '配置 CI/CD 流水线',
      progress: 60,
      tokens: 2760,
      color: 'agent-devops'
    }
  ];

  // 模拟时间线数据
  const timelineData = [
    { agent: 'Product Manager', start: 0, duration: 120, status: 'completed' },
    { agent: 'Architect', start: 60, duration: 180, status: 'running' },
    { agent: 'Frontend Dev', start: 150, duration: 120, status: 'running' },
    { agent: 'Backend Dev', start: 240, duration: 90, status: 'waiting' },
    { agent: 'DevOps Engineer', start: 180, duration: 150, status: 'running' },
  ];

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'idle': return 'bg-pending text-on-surface-variant';
      case 'running': return 'bg-running/10 text-running border border-running/20';
      case 'waiting': return 'bg-warning/10 text-warning border border-warning/20';
      case 'completed': return 'bg-success/10 text-success border border-success/20';
      default: return 'bg-surface-container-high text-on-surface-variant';
    }
  };

  const getStatusIndicator = (status: string) => {
    switch (status) {
      case 'idle': return <div className="w-2 h-2 rounded-full bg-pending"></div>;
      case 'running': return <div className="w-2 h-2 rounded-full bg-running animate-pulse"></div>;
      case 'waiting': return <div className="w-2 h-2 rounded-full bg-warning animate-pulse"></div>;
      case 'completed': return <div className="w-2 h-2 rounded-full bg-success"></div>;
      default: return <div className="w-2 h-2 rounded-full bg-outline"></div>;
    }
  };

  return (
    <div className="h-full overflow-y-auto p-6">
      {/* Page Header */}
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-on-surface mb-2">AI 开发流水线仪表盘</h1>
        <p className="text-sm text-on-surface-variant">
          监控 AI Agent 团队的执行状态、查看实时日志与项目输出
        </p>
      </div>

      {/* Agent Status Dashboard (设计规范 5.1.1) */}
      <section className="mb-10">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-on-surface flex items-center gap-2">
            <span className="material-symbols-outlined text-primary">groups</span>
            Agent 状态仪表盘
          </h2>
          <div className="text-xs text-outline">
            最后更新: 刚刚
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {agents.map((agent) => (
            <div
              key={agent.id}
              className={`bg-surface-container rounded-xl p-4 border transition-all duration-200 hover:border-outline-variant/50 ${
                agent.status === 'running' ? 'animate-pulse-border border-running/30' : 'border-outline-variant/15'
              }`}
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-3">
                  <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${
                    agent.color === 'agent-product-manager' ? 'bg-agent-product-manager/10' :
                    agent.color === 'agent-architect' ? 'bg-agent-architect/10' :
                    agent.color === 'agent-frontend' ? 'bg-agent-frontend/10' :
                    agent.color === 'agent-backend' ? 'bg-agent-backend/10' :
                    agent.color === 'agent-qa' ? 'bg-agent-qa/10' :
                    agent.color === 'agent-devops' ? 'bg-agent-devops/10' :
                    'bg-surface-container-high'
                  }`}>
                    <span className="material-symbols-outlined text-lg text-on-surface-variant">
                      smart_toy
                    </span>
                  </div>
                  <div>
                    <h3 className="font-medium text-on-surface">{agent.name}</h3>
                    <p className="text-xs text-outline">{agent.role}</p>
                  </div>
                </div>
                <div className={`flex items-center gap-1.5 px-2 py-1 rounded-full text-xs ${getStatusColor(agent.status)}`}>
                  {getStatusIndicator(agent.status)}
                  <span className="capitalize">{agent.status}</span>
                </div>
              </div>

              <p className="text-sm text-on-surface-variant mb-3 line-clamp-2">{agent.task}</p>

              {agent.status === 'running' && (
                <div className="space-y-2">
                  <div className="flex items-center justify-between text-xs">
                    <span className="text-outline">进度</span>
                    <span className="font-mono text-on-surface">{agent.progress}%</span>
                  </div>
                  <div className="h-1.5 bg-surface-container-high rounded-full overflow-hidden">
                    <div
                      className="h-full bg-running rounded-full transition-all duration-500"
                      style={{ width: `${agent.progress}%` }}
                    ></div>
                  </div>
                </div>
              )}

              <div className="flex items-center justify-between mt-3 pt-3 border-t border-outline-variant/15">
                <div className="text-xs text-outline">
                  <span className="font-mono">{agent.tokens.toLocaleString()}</span> tokens
                </div>
                <button className="text-xs text-primary hover:text-primary-container transition-colors">
                  查看详情 →
                </button>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Execution Timeline (设计规范 5.1.2) */}
      <section className="mb-10">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-on-surface flex items-center gap-2">
            <span className="material-symbols-outlined text-primary">timeline</span>
            执行时间线
          </h2>
          <div className="flex items-center gap-2 text-xs">
            <button className="px-2 py-1 rounded bg-surface-container-high text-outline hover:text-on-surface transition-colors">
              1小时
            </button>
            <button className="px-2 py-1 rounded bg-primary/10 text-primary">
              全部
            </button>
            <button className="px-2 py-1 rounded bg-surface-container-high text-outline hover:text-on-surface transition-colors">
              24小时
            </button>
          </div>
        </div>

        <div className="bg-surface-container rounded-xl p-4">
          <div className="relative h-32">
            {/* Timeline grid */}
            <div className="absolute inset-0 flex">
              {[...Array(6)].map((_, i) => (
                <div key={i} className="flex-1 border-r border-outline-variant/15 relative">
                  <div className="absolute bottom-0 left-0 right-0 text-center text-xs text-outline pt-1">
                    {i * 60}min
                  </div>
                </div>
              ))}
            </div>

            {/* Agent timeline bars */}
            {timelineData.map((item, index) => (
              <div
                key={index}
                className="absolute h-6 rounded-md opacity-80 hover:opacity-100 transition-opacity"
                style={{
                  left: `${(item.start / 360) * 100}%`,
                  width: `${(item.duration / 360) * 100}%`,
                  top: `${index * 28 + 20}px`,
                  backgroundColor: `var(--color-agent-${item.agent.toLowerCase().replace(' ', '-')})`,
                }}
              >
                <div className="absolute inset-0 flex items-center px-2">
                  <span className="text-xs font-medium text-white truncate">{item.agent}</span>
                </div>
              </div>
            ))}
          </div>

          <div className="flex items-center justify-between mt-6 pt-4 border-t border-outline-variant/15">
            <div className="text-xs text-outline">
              时间范围: 360分钟 | 并行 Agent: {timelineData.filter(d => d.status === 'running').length}个
            </div>
            <button className="text-xs text-primary hover:text-primary-container transition-colors flex items-center gap-1">
              <span className="material-symbols-outlined text-sm">open_in_new</span>
              展开时间线
            </button>
          </div>
        </div>
      </section>

      {/* Tab切换面板 (设计规范 5.1.3) */}
      <section>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-on-surface flex items-center gap-2">
            <span className="material-symbols-outlined text-primary">monitoring</span>
            实时监控面板
          </h2>
        </div>

        <div className="bg-surface-container rounded-xl overflow-hidden">
          {/* Tab Navigation */}
          <div className="flex border-b border-outline-variant/15">
            <button
              onClick={() => setActiveTab('logs')}
              className={`px-6 py-3 text-sm font-medium transition-colors border-b-2 ${
                activeTab === 'logs'
                  ? 'text-primary border-primary'
                  : 'text-outline border-transparent hover:text-on-surface-variant'
              }`}
            >
              <span className="flex items-center gap-2">
                <span className="material-symbols-outlined text-base">list_alt</span>
                执行日志
              </span>
            </button>
            <button
              onClick={() => setActiveTab('terminal')}
              className={`px-6 py-3 text-sm font-medium transition-colors border-b-2 ${
                activeTab === 'terminal'
                  ? 'text-primary border-primary'
                  : 'text-outline border-transparent hover:text-on-surface-variant'
              }`}
            >
              <span className="flex items-center gap-2">
                <span className="material-symbols-outlined text-base">terminal</span>
                终端输出
              </span>
            </button>
            <button
              onClick={() => setActiveTab('files')}
              className={`px-6 py-3 text-sm font-medium transition-colors border-b-2 ${
                activeTab === 'files'
                  ? 'text-primary border-primary'
                  : 'text-outline border-transparent hover:text-on-surface-variant'
              }`}
            >
              <span className="flex items-center gap-2">
                <span className="material-symbols-outlined text-base">folder</span>
                输出文件
              </span>
            </button>
          </div>

          {/* Tab Content */}
          <div className="p-4 min-h-[300px]">
            {activeTab === 'logs' && (
              <div className="space-y-3">
                <div className="flex items-center gap-3 p-3 rounded-lg bg-surface-container-high">
                  <div className="w-2 h-2 rounded-full bg-info"></div>
                  <div className="flex-1">
                    <div className="text-sm font-mono text-on-surface">[INFO] 2025-04-19 14:30:25</div>
                    <div className="text-sm text-on-surface-variant">Product Manager 完成了需求分析，生成了 12 个用户故事</div>
                  </div>
                </div>
                <div className="flex items-center gap-3 p-3 rounded-lg bg-surface-container-high">
                  <div className="w-2 h-2 rounded-full bg-running"></div>
                  <div className="flex-1">
                    <div className="text-sm font-mono text-on-surface">[RUNNING] 2025-04-19 14:32:10</div>
                    <div className="text-sm text-on-surface-variant">Architect 正在设计微服务架构，评估技术选型...</div>
                  </div>
                </div>
                <div className="flex items-center gap-3 p-3 rounded-lg bg-surface-container-high">
                  <div className="w-2 h-2 rounded-full bg-success"></div>
                  <div className="flex-1">
                    <div className="text-sm font-mono text-on-surface">[SUCCESS] 2025-04-19 14:28:45</div>
                    <div className="text-sm text-on-surface-variant">Frontend Dev 成功创建了 5 个 React 组件</div>
                  </div>
                </div>
                <div className="text-center py-4">
                  <p className="text-sm text-outline">日志实时流式更新中...</p>
                </div>
              </div>
            )}

            {activeTab === 'terminal' && (
              <div className="bg-surface-container-high rounded-lg p-4 font-mono text-sm">
                <div className="text-success">$ npm run dev</div>
                <div className="text-on-surface-variant mt-2">
                  <div>Starting development server...</div>
                  <div>Compiled successfully!</div>
                  <div className="mt-4 text-outline">Server running at http://localhost:5173</div>
                </div>
                <div className="mt-4 flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-success animate-pulse"></div>
                  <span className="text-xs text-outline">终端活跃中</span>
                </div>
              </div>
            )}

            {activeTab === 'files' && (
              <div className="space-y-4">
                <div className="flex items-center justify-between p-3 rounded-lg bg-surface-container-high">
                  <div className="flex items-center gap-3">
                    <span className="material-symbols-outlined text-outline">description</span>
                    <div>
                      <div className="text-sm font-medium text-on-surface">src/components/Login.tsx</div>
                      <div className="text-xs text-outline">新增 | +120 行</div>
                    </div>
                  </div>
                  <button className="text-xs text-primary hover:text-primary-container transition-colors">
                    查看差异
                  </button>
                </div>
                <div className="flex items-center justify-between p-3 rounded-lg bg-surface-container-high">
                  <div className="flex items-center gap-3">
                    <span className="material-symbols-outlined text-outline">description</span>
                    <div>
                      <div className="text-sm font-medium text-on-surface">package.json</div>
                      <div className="text-xs text-outline">修改 | +3 依赖项</div>
                    </div>
                  </div>
                  <button className="text-xs text-primary hover:text-primary-container transition-colors">
                    查看差异
                  </button>
                </div>
                <div className="text-center py-4">
                  <p className="text-sm text-outline">共生成 8 个文件，修改 12 个文件</p>
                </div>
              </div>
            )}
          </div>
        </div>
      </section>
    </div>
  );
};

export default DashboardViewGenpulse;