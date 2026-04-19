import React from 'react';
import { useAppStore, selectCurrentView } from '../stores/appStore';

const SidebarNew: React.FC = () => {
  const currentView = useAppStore(selectCurrentView);
  const { setCurrentView, projectList, currentProject } = useAppStore();

  const menuItems = [
    { id: 'dashboard', label: 'Dashboard', icon: 'grid_view', description: '概览与监控' },
    { id: 'projects', label: 'Projects', icon: 'account_tree', description: '项目管理' },
    { id: 'agents', label: 'Agents', icon: 'psychology', description: '代理管理' },
    { id: 'skills', label: 'Skills', icon: 'memory', description: '技能库' },
    { id: 'memory', label: 'Memory', icon: 'database', description: '神经资产' },
    { id: 'kanban', label: 'Kanban', icon: 'view_kanban', description: '执行看板' },
    { id: 'terminal', label: 'Terminal', icon: 'terminal', description: '终端与对比' },
    { id: 'settings', label: 'Settings', icon: 'settings', description: '系统设置' },
  ];

  return (
    <nav className="h-screen w-64 fixed left-0 top-0 flex flex-col bg-surface border-r border-outline-variant/15 z-40 hidden md:flex">
      {/* 头部 */}
      <div className="px-6 py-6 flex items-center gap-3">
        <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-primary-container to-inverse-primary flex items-center justify-center shadow-[0_0_15px_rgba(91,95,255,0.3)]">
          <span className="material-symbols-outlined text-on-primary-container text-sm" data-weight="fill">
            token
          </span>
        </div>
        <div>
          <h1 className="text-primary font-bold text-lg leading-tight tracking-tight">Genpulse</h1>
          <p className="font-body text-xs tracking-wide text-outline">AI Development</p>
        </div>
      </div>

      {/* CTA按钮 */}
      <div className="px-4 mb-6">
        <button className="w-full bg-primary-container hover:bg-surface-bright text-on-primary-container py-2.5 px-4 rounded-lg flex items-center justify-center gap-2 transition-all duration-200">
          <span className="material-symbols-outlined text-sm">add</span>
          <span className="font-medium text-sm">New Pipeline</span>
        </button>
      </div>

      {/* 主导航 */}
      <div className="flex-1 px-2 space-y-1">
        {menuItems.map((item) => (
          <button
            key={item.id}
            className={`flex items-center gap-3 py-3 px-4 rounded-lg transition-all duration-200 cursor-pointer active:opacity-80 w-full text-left ${
              currentView === item.id
                ? 'text-primary bg-primary-container/10 border-r-2 border-primary'
                : 'text-outline hover:bg-surface-container-low hover:text-primary'
            }`}
            onClick={() => setCurrentView(item.id)}
          >
            <span className="material-symbols-outlined">{item.icon}</span>
            <div className="flex-1">
              <span className="font-label text-sm tracking-wide font-medium block">{item.label}</span>
              <span className="font-body text-xs text-outline/70 block mt-0.5">{item.description}</span>
            </div>
          </button>
        ))}
      </div>

      {/* 项目列表 */}
      <div className="px-2 py-4 border-t border-outline-variant/15">
        <h3 className="font-label text-sm font-medium text-outline px-4 mb-2">Projects</h3>
        {projectList.length === 0 ? (
          <p className="text-xs text-outline/50 px-4 py-2">No projects yet</p>
        ) : (
          <ul className="space-y-1">
            {projectList.map((project) => (
              <li key={project}>
                <button
                  className={`flex items-center gap-2 py-2 px-4 rounded-lg w-full text-left transition-colors ${
                    currentProject === project
                      ? 'bg-surface-container-low text-primary'
                      : 'text-outline hover:bg-surface-container-low hover:text-on-surface'
                  }`}
                  onClick={() => useAppStore.getState().setCurrentProject(project)}
                >
                  <span className="material-symbols-outlined text-sm">folder</span>
                  <span className="font-body text-sm truncate">{project}</span>
                </button>
              </li>
            ))}
          </ul>
        )}
        <button 
          className="w-full mt-2 text-xs font-medium text-primary hover:text-primary-container transition-colors flex items-center justify-center gap-1 py-2"
          onClick={() => {/* TODO: 实现新建项目 */}}
        >
          <span className="material-symbols-outlined text-sm">add</span>
          New Project
        </button>
      </div>

      {/* 底部导航 */}
      <div className="px-2 pb-6 pt-4 border-t border-outline-variant/15 space-y-1">
        <button className="flex items-center gap-3 py-2.5 px-4 rounded-lg text-outline hover:bg-surface-container-low hover:text-primary transition-all duration-200 w-full">
          <span className="material-symbols-outlined text-[20px]">description</span>
          <span className="font-body text-sm tracking-wide">Documentation</span>
        </button>
        <button className="flex items-center gap-3 py-2.5 px-4 rounded-lg text-outline hover:bg-surface-container-low hover:text-primary transition-all duration-200 w-full">
          <span className="material-symbols-outlined text-[20px]">contact_support</span>
          <span className="font-body text-sm tracking-wide">Support</span>
        </button>
      </div>

      {/* 系统状态 */}
      <div className="px-4 py-3 border-t border-outline-variant/15">
        <div className="text-xs space-y-1.5">
          <div className="flex justify-between">
            <span className="text-outline">CPU</span>
            <span className="font-mono text-sys-success">12%</span>
          </div>
          <div className="flex justify-between">
            <span className="text-outline">Memory</span>
            <span className="font-mono text-primary">45%</span>
          </div>
          <div className="flex justify-between">
            <span className="text-outline">Agents</span>
            <span className="font-mono text-secondary">3/5</span>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default SidebarNew;