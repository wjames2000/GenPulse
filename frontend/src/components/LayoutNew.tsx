import React from 'react';
import { useAppStore } from '../stores/appStore';
import SidebarNew from './SidebarNew';
import HeaderNew from './HeaderNew';

interface LayoutProps {
  children: React.ReactNode;
}

const LayoutNew: React.FC<LayoutProps> = ({ children }) => {
  const { darkMode, sidebarOpen } = useAppStore();

  return (
    <div className={`h-screen overflow-hidden ${darkMode ? 'dark' : ''}`}>
      <div className="h-screen overflow-hidden flex bg-surface text-on-surface antialiased">
        {/* 侧边栏 */}
        <SidebarNew />

        {/* 主内容区域 */}
        <div className={`flex-1 ${sidebarOpen ? 'ml-0 md:ml-64' : 'ml-0'} flex flex-col h-screen overflow-hidden relative`}>
          {/* 顶部应用栏 */}
          <HeaderNew />

          {/* 主要内容 */}
          <main className="flex-1 overflow-y-auto">
            {children}
          </main>

          {/* 移动端侧边栏遮罩 */}
          {sidebarOpen && (
            <div 
              className="md:hidden fixed inset-0 bg-black/50 z-30"
              onClick={() => useAppStore.getState().toggleSidebar()}
            />
          )}
        </div>
      </div>
    </div>
  );
};

export default LayoutNew;