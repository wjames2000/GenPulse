import React from 'react';

const DashboardViewGenpulseThree: React.FC = () => {
  // 简化版数据，方便三栏布局演示
  const agents = [
    { id: 1, name: 'Product Manager', role: '产品经理', status: 'completed', progress: 100, task: '需求分析与用户故事创建' },
    { id: 2, name: 'Architect', role: '架构师', status: 'running', progress: 75, task: '设计微服务架构' },
    { id: 3, name: 'Frontend Dev', role: '前端开发', status: 'running', progress: 45, task: '生成 React 组件库' },
  ];
  return (
    <div className="h-full flex w-full">
      {/* Left Sidebar */}
      <aside className="w-64 bg-surface-container-low border-r border-outline-variant/15 p-4 overflow-y-auto hidden md:block">
        <div className="font-semibold text-on-surface mb-3">导航</div>
        <nav className="space-y-2 text-sm text-on-surface-variant">
          <div>仪表盘总览</div>
          <div>Agent 状态</div>
          <div>执行时间线</div>
          <div>设置</div>
        </nav>
      </aside>

      {/* Middle Content */}
      <main className="flex-1 p-6 overflow-auto">
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-on-surface mb-2">Agent 状态仪表盘</h2>
          <p className="text-on-surface-variant text-sm">统一视图，展示核心代理的健康状态和进度</p>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {agents.map(a => (
            <div key={a.id} className="bg-surface-container rounded-xl p-4 border border-outline-variant/15">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-2">
                  <span className="material-symbols-outlined text-primary">person</span>
                  <span className="text-on-surface font-medium">{a.name}</span>
                </div>
                <span className="text-xs text-on-surface-variant">{a.role}</span>
              </div>
              <div className="text-xs text-on-surface-variant mb-2">任务: {a.task}</div>
              <div className="h-2 rounded-full bg-surface-container-high mb-2">
                <div className={
                  a.status === 'running' ? 'h-full bg-running' : 'h-full bg-success'
                } style={{ width: a.progress + '%' }}></div>
              </div>
              <div className="flex items-center justify-between text-xs text-on-surface-variant">
                <span>状态：{a.status}</span>
                <span>进度 {a.progress}%</span>
              </div>
            </div>
          ))}
        </div>
      </main>

      {/* Right Panel */}
      <aside className="w-[360px] bg-surface-container-low border-l border-outline-variant/15 p-4 overflow-y-auto hidden xl:block">
        <div className="font-semibold text-on-surface mb-3">执行摘要</div>
        <div className="space-y-3 text-sm text-on-surface-variant">
          <div>总代理数: {agents.length}</div>
          <div>正在执行: 2</div>
          <div>最近更新时间: 刚刚</div>
        </div>
        <hr className="my-4 border-outline-variant/20" />
        <div className="font-semibold text-on-surface mb-2">时间线</div>
        <div className="text-xs text-on-surface-variant">Just a lightweight timeline mock</div>
      </aside>
    </div>
  );
};

export default DashboardViewGenpulseThree;
