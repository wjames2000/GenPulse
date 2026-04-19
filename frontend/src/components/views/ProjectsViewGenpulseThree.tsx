import React from 'react';

type Project = {
  id: string;
  name: string;
  description: string;
  type: string;
  status: string;
  agents: number;
  lastUpdated: string;
  progress: number;
};

const ProjectsViewGenpulseThree: React.FC = () => {
  const projects: Project[] = [
    { id: '1', name: 'E-Commerce Platform', description: 'Full-stack e-commerce with AI', type: 'Web', status: 'running', agents: 4, lastUpdated: '2小时前', progress: 65 },
    { id: '2', name: 'API Gateway', description: 'Gateway with auth', type: 'API', status: 'completed', agents: 3, lastUpdated: '1天前', progress: 100 },
    { id: '3', name: 'CLI Tool', description: 'CLI for scaffolding', type: 'CLI', status: 'running', agents: 2, lastUpdated: '3小时前', progress: 30 },
  ];

  return (
    <div className="h-full flex w-full">
      {/* Left Sidebar */}
      <aside className="w-64 bg-surface-container-low border-r border-outline-variant/15 p-4 overflow-y-auto hidden md:block">
        <div className="font-semibold text-on-surface mb-2">过滤与操作</div>
        <div className="space-y-2 text-sm text-on-surface-variant">
          <div>全部</div>
          <div>进行中</div>
          <div>已完成</div>
        </div>
        <hr className="my-4 border-outline-variant/20" />
        <button className="w-full rounded-lg bg-primary-container text-on-primary-container py-2">新建流水线</button>
      </aside>

      {/* Middle Content */}
      <main className="flex-1 p-6 overflow-auto">
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-on-surface mb-2">项目流水线</h2>
          <p className="text-on-surface-variant text-sm">管理您的 AI 项目开发流水线，监控执行状态</p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {projects.map(p => (
            <div key={p.id} className="bg-surface-container rounded-xl p-4 border border-outline-variant/15 hover:bg-surface-container-high transition-colors">
              <div className="flex items-start justify-between mb-2">
                <div className="flex items-center gap-2">
                  <span className="material-symbols-outlined text-primary">folder</span>
                  <div>
                    <div className="font-medium text-on-surface">{p.name}</div>
                    <div className="text-xs text-outline">{p.type} • {p.status}</div>
                  </div>
                </div>
                <div className={`text-xs px-2 py-1 rounded-full ${p.status==='running'?'bg-running/10 text-running border border-running/20':'bg-success/10 text-success border border-success/20'}`}>
                  {p.status}
                </div>
              </div>
              <p className="text-sm text-on-surface-variant mb-2 line-clamp-2">{p.description}</p>
              <div className="h-1.5 bg-surface-container-high rounded-full">
                <div className="h-full bg-running rounded-full" style={{ width: p.progress + '%' }}></div>
              </div>
              <div className="flex items-center justify-between mt-3 text-xs text-outline">
                <span>进度 {p.progress}%</span>
                <span>最后更新: {p.lastUpdated}</span>
              </div>
            </div>
          ))}
        </div>
      </main>

      {/* Right Panel */}
      <aside className="w-[360px] bg-surface-container-low border-l border-outline-variant/15 p-4 overflow-y-auto hidden xl:block">
        <div className="font-semibold text-on-surface mb-3">快速摘要</div>
        <div className="text-sm text-on-surface-variant">最近活动、统计和快速操作将显示在这里。</div>
      </aside>
    </div>
  );
};

export default ProjectsViewGenpulseThree;
