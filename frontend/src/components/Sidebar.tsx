import React from 'react';
import { useAppStore, selectCurrentView } from '../stores/appStore';

const Sidebar: React.FC = () => {
  const currentView = useAppStore(selectCurrentView);
  const { setCurrentView, projectList, currentProject } = useAppStore();

  const menuItems = [
    { id: 'dashboard', label: 'Dashboard', icon: '📊' },
    { id: 'projects', label: 'Projects', icon: '📁' },
    { id: 'agents', label: 'Agents', icon: '🤖' },
    { id: 'skills', label: 'Skills', icon: '🧠' },
    { id: 'memory', label: 'Memory', icon: '💾' },
    { id: 'monitoring', label: 'Monitoring', icon: '📈' },
    { id: 'settings', label: 'Settings', icon: '⚙️' },
  ];

  return (
    <aside className="sidebar">
      <nav className="sidebar-nav">
        <ul className="sidebar-menu">
          {menuItems.map((item) => (
            <li key={item.id}>
              <button
                className={`sidebar-menu-item ${currentView === item.id ? 'active' : ''}`}
                onClick={() => setCurrentView(item.id)}
              >
                <span className="sidebar-menu-icon">{item.icon}</span>
                <span className="sidebar-menu-label">{item.label}</span>
              </button>
            </li>
          ))}
        </ul>
      </nav>

      <div className="sidebar-projects">
        <h3 className="sidebar-section-title">Projects</h3>
        {projectList.length === 0 ? (
          <p className="sidebar-empty">No projects yet</p>
        ) : (
          <ul className="project-list">
            {projectList.map((project) => (
              <li key={project}>
                <button
                  className={`project-item ${currentProject === project ? 'active' : ''}`}
                  onClick={() => useAppStore.getState().setCurrentProject(project)}
                >
                  📄 {project}
                </button>
              </li>
            ))}
          </ul>
        )}
        <button 
          className="sidebar-button"
          onClick={() => {/* TODO: 实现新建项目 */}}
        >
          + New Project
        </button>
      </div>

      <div className="sidebar-footer">
        <div className="system-status">
          <div className="status-item">
            <span className="status-label">CPU:</span>
            <span className="status-value">12%</span>
          </div>
          <div className="status-item">
            <span className="status-label">Memory:</span>
            <span className="status-value">45%</span>
          </div>
        </div>
      </div>
    </aside>
  );
};

export default Sidebar;