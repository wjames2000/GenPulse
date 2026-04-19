import React, { useState } from 'react';
import { useAppStore } from '../stores/appStore';
import SidebarGenpulse from './SidebarGenpulse';
import TopAppBar from './TopAppBar';
import RightPanel from './RightPanel';

interface LayoutProps {
  children: React.ReactNode;
}

const LayoutGenpulse: React.FC<LayoutProps> = ({ children }) => {
  const { darkMode } = useAppStore();
  const [rightPanelOpen, setRightPanelOpen] = useState(true);
  const [rightPanelWidth, setRightPanelWidth] = useState(360); // 默认宽度 360px

  return (
    <div className={`h-screen overflow-hidden ${darkMode ? 'dark' : ''}`}>
      <div className="h-screen overflow-hidden flex bg-surface text-on-surface antialiased font-body">
        {/* 左侧边栏 - 固定 260px (来自设计规范 4.3.2) */}
        <div className="w-64 flex-shrink-0 h-screen">
          <SidebarGenpulse />
        </div>

        {/* 主内容区域 - 自适应 */}
        <div className="flex-1 flex flex-col h-screen overflow-hidden">
          {/* 顶部应用栏 - 固定 56px (来自设计规范 4.3.1) */}
          <div className="h-14 flex-shrink-0">
            <TopAppBar />
          </div>

          {/* 主要内容区域 + 右侧面板 */}
          <div className="flex-1 flex overflow-hidden">
            {/* 主内容区 */}
            <main className={`flex-1 overflow-y-auto ${rightPanelOpen ? '' : ''}`}>
              <div className="h-full">
                {children}
              </div>
            </main>

            {/* 右侧辅助面板 - 可收起 (来自设计规范 4.3.3) */}
            {rightPanelOpen && (
              <div 
                className="flex-shrink-0 h-full overflow-hidden border-l border-outline-variant/15"
                style={{ width: `${rightPanelWidth}px` }}
              >
                <RightPanel 
                  onClose={() => setRightPanelOpen(false)}
                  width={rightPanelWidth}
                  onWidthChange={setRightPanelWidth}
                />
              </div>
            )}
          </div>

          {/* 右侧面板切换按钮 */}
          {!rightPanelOpen && (
            <button
              onClick={() => setRightPanelOpen(true)}
              className="absolute right-0 top-1/2 transform -translate-y-1/2 bg-surface-container-high text-outline hover:text-primary hover:bg-surface-container-highest p-2 rounded-l-lg border border-outline-variant/30 border-r-0 transition-all"
            >
              <span className="material-symbols-outlined text-lg">
                chevron_left
              </span>
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

export default LayoutGenpulse;