import React from 'react';

const ProjectsView: React.FC = () => {
  return (
    <div className="p-6 md:p-8 lg:p-12">
      <div className="max-w-6xl mx-auto">
        <div className="mb-10">
          <h2 className="text-[2.25rem] font-bold tracking-[-0.02em] text-on-surface mb-2">
            项目管理
          </h2>
          <p className="text-on-surface-variant text-sm">
            创建、管理和监控您的AI项目开发流水线。
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {/* 项目卡片 */}
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <div key={i} className="bg-surface-container rounded-xl p-6 hover:bg-surface-container-high transition-colors cursor-pointer group">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="p-2 rounded-lg bg-surface-container-low text-primary">
                    <span className="material-symbols-outlined">folder</span>
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold text-on-surface">Project {i}</h3>
                    <p className="text-sm text-outline">AI Development Pipeline</p>
                  </div>
                </div>
                <span className="material-symbols-outlined text-outline group-hover:text-primary transition-colors">
                  arrow_forward
                </span>
              </div>
              
              <div className="space-y-3">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-outline">状态</span>
                  <span className="font-mono text-sys-success">运行中</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-outline">代理数</span>
                  <span className="font-mono text-primary">3</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-outline">最后更新</span>
                  <span className="font-mono text-outline">2小时前</span>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* 新建项目卡片 */}
        <div className="mt-8 bg-surface-container-low rounded-xl p-8 border-2 border-dashed border-outline-variant/30 hover:border-primary/50 transition-colors cursor-pointer text-center">
          <div className="w-16 h-16 rounded-full bg-surface-container mx-auto mb-4 flex items-center justify-center">
            <span className="material-symbols-outlined text-3xl text-primary">add</span>
          </div>
          <h3 className="text-xl font-semibold text-on-surface mb-2">新建项目</h3>
          <p className="text-outline mb-6">开始一个新的AI开发流水线项目</p>
          <button className="bg-primary-container text-on-primary-container hover:bg-inverse-primary rounded-full px-6 py-2 text-sm font-medium shadow-lg shadow-primary/20 transition-all">
            创建项目
          </button>
        </div>
      </div>
    </div>
  );
};

export default ProjectsView;