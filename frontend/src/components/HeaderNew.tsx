import React, { useState } from 'react';
import { useAppStore } from '../stores/appStore';

const HeaderNew: React.FC = () => {
  const { sidebarOpen, toggleSidebar } = useAppStore();
  const [searchQuery, setSearchQuery] = useState('');

  return (
    <header className="flex justify-between items-center w-full px-6 py-3 bg-surface/80 backdrop-blur-xl shadow-2xl shadow-black/40 sticky top-0 z-50">
      {/* 左侧：品牌/搜索 */}
      <div className="flex items-center gap-6 flex-1">
        {/* 移动端菜单按钮 */}
        <button
          className="md:hidden flex items-center gap-3"
          onClick={toggleSidebar}
        >
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-primary-container to-inverse-primary flex items-center justify-center">
            <span className="material-symbols-outlined text-on-primary-container text-sm" data-weight="fill">
              token
            </span>
          </div>
        </button>

        {/* 品牌标题 */}
        <span className="text-xl font-bold tracking-tight text-primary font-headline antialiased hidden lg:block">
          Genpulse AI
        </span>

        {/* 搜索框 */}
        <div className="relative w-full max-w-xs md:w-64">
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-outline text-sm">
            search
          </span>
          <input
            className="w-full bg-surface-container-lowest text-on-surface text-sm rounded-lg pl-9 pr-3 py-1.5 border-none focus:ring-0 focus:border-b-2 focus:border-b-primary transition-all placeholder:text-outline outline-none"
            placeholder="搜索..."
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
      </div>

      {/* 右侧：操作和用户 */}
      <div className="flex items-center gap-4">
        {/* 操作按钮 */}
        <div className="hidden md:flex items-center gap-1 text-primary">
          <button 
            className="p-2 rounded-full hover:bg-surface-container-high/50 transition-colors duration-200 active:scale-95 flex items-center justify-center"
            title="Notifications"
          >
            <span className="material-symbols-outlined">notifications</span>
          </button>
          <button 
            className="p-2 rounded-full hover:bg-surface-container-high/50 transition-colors duration-200 active:scale-95 flex items-center justify-center"
            title="Settings"
          >
            <span className="material-symbols-outlined">settings</span>
          </button>
          <button 
            className="p-2 rounded-full hover:bg-surface-container-high/50 transition-colors duration-200 active:scale-95 flex items-center justify-center"
            title="Help"
          >
            <span className="material-symbols-outlined">help</span>
          </button>
        </div>

        {/* 用户头像 */}
        <div className="w-8 h-8 rounded-full bg-surface-container-high overflow-hidden border border-outline-variant/15 flex-shrink-0 cursor-pointer hover:ring-2 hover:ring-primary/50 transition-all">
          <div className="w-full h-full bg-gradient-to-br from-primary-container to-inverse-primary flex items-center justify-center">
            <span className="text-on-primary-container text-xs font-bold">U</span>
          </div>
        </div>
      </div>
    </header>
  );
};

export default HeaderNew;