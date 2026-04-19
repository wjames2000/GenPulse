import React from 'react';
import { useAppStore } from '../stores/appStore';

const Header: React.FC = () => {
  const { toggleSidebar, toggleDarkMode, darkMode, appInfo } = useAppStore();

  return (
    <header className="header">
      <div className="header-left">
        <button 
          className="header-button" 
          onClick={toggleSidebar}
          title="Toggle Sidebar"
        >
          ☰
        </button>
        <div className="header-title">
          <h1>GenPulse</h1>
          {appInfo && (
            <span className="header-subtitle">
              v{appInfo.version} • <span className={`status-${appInfo.status}`}>{appInfo.status}</span>
            </span>
          )}
        </div>
      </div>
      
      <div className="header-right">
        <button 
          className="header-button" 
          onClick={toggleDarkMode}
          title={darkMode ? 'Switch to Light Mode' : 'Switch to Dark Mode'}
        >
          {darkMode ? '☀️' : '🌙'}
        </button>
        
        <div className="header-actions">
          {/* 这里可以添加更多操作按钮 */}
        </div>
      </div>
    </header>
  );
};

export default Header;