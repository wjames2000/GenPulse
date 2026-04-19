import React, { useState } from 'react';

interface RightPanelProps {
  onClose: () => void;
  width: number;
  onWidthChange: (width: number) => void;
}

const RightPanel: React.FC<RightPanelProps> = ({ onClose, width, onWidthChange }) => {
  const [activeTab, setActiveTab] = useState<'agents' | 'thought' | 'files'>('agents');
  const [isResizing, setIsResizing] = useState(false);

  const handleMouseDown = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsResizing(true);
    
    const startX = e.clientX;
    const startWidth = width;
    
    const handleMouseMove = (moveEvent: MouseEvent) => {
      const deltaX = startX - moveEvent.clientX;
      const newWidth = startWidth + deltaX;
      
      // 限制宽度在 280px 到 480px 之间 (来自设计规范 4.3.3)
      if (newWidth >= 280 && newWidth <= 480) {
        onWidthChange(newWidth);
      }
    };
    
    const handleMouseUp = () => {
      setIsResizing(false);
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
    
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
  };

  return (
    <div className="h-full flex flex-col bg-surface-container-low">
      {/* Panel Header */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-outline-variant/15">
        <div className="flex items-center gap-2">
          <h3 className="text-sm font-semibold text-on-surface">辅助面板</h3>
          <span className="text-xs px-2 py-0.5 rounded bg-surface-container-high text-outline">
            Beta
          </span>
        </div>
        <button
          onClick={onClose}
          className="p-1 rounded hover:bg-surface-container-high text-outline hover:text-on-surface transition-colors"
        >
          <span className="material-symbols-outlined text-lg">close</span>
        </button>
      </div>

      {/* Resize Handle */}
      <div
        className="w-1.5 h-full absolute left-0 top-0 cursor-col-resize hover:bg-primary/50 active:bg-primary"
        onMouseDown={handleMouseDown}
      />

      {/* Tab Navigation */}
      <div className="flex border-b border-outline-variant/15">
        <button
          onClick={() => setActiveTab('agents')}
          className={`flex-1 py-2.5 text-xs font-medium transition-colors ${
            activeTab === 'agents'
              ? 'text-primary border-b-2 border-primary'
              : 'text-outline hover:text-on-surface-variant'
          }`}
        >
          <span className="flex items-center justify-center gap-1">
            <span className="material-symbols-outlined text-sm">smart_toy</span>
            Agent 状态
          </span>
        </button>
        <button
          onClick={() => setActiveTab('thought')}
          className={`flex-1 py-2.5 text-xs font-medium transition-colors ${
            activeTab === 'thought'
              ? 'text-primary border-b-2 border-primary'
              : 'text-outline hover:text-on-surface-variant'
          }`}
        >
          <span className="flex items-center justify-center gap-1">
            <span className="material-symbols-outlined text-sm">chat</span>
            思维链
          </span>
        </button>
        <button
          onClick={() => setActiveTab('files')}
          className={`flex-1 py-2.5 text-xs font-medium transition-colors ${
            activeTab === 'files'
              ? 'text-primary border-b-2 border-primary'
              : 'text-outline hover:text-on-surface-variant'
          }`}
        >
          <span className="flex items-center justify-center gap-1">
            <span className="material-symbols-outlined text-sm">difference</span>
            文件变更
          </span>
        </button>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-4">
        {activeTab === 'agents' && (
          <div className="space-y-3">
            <div className="bg-surface-container rounded-lg p-3">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-running animate-pulse"></div>
                  <span className="text-xs font-medium text-on-surface">Frontend Dev</span>
                </div>
                <span className="text-xs text-outline">React</span>
              </div>
              <p className="text-xs text-on-surface-variant line-clamp-2">
                正在生成 Login 组件...
              </p>
              <div className="flex items-center justify-between mt-3">
                <div className="text-xs text-outline">65%</div>
                <div className="text-xs font-mono text-outline">1,248 tokens</div>
              </div>
              <div className="h-1.5 bg-surface-container-high rounded-full overflow-hidden mt-1">
                <div className="h-full bg-running w-2/3"></div>
              </div>
            </div>

            <div className="bg-surface-container rounded-lg p-3">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-pending"></div>
                  <span className="text-xs font-medium text-on-surface">Backend Dev</span>
                </div>
                <span className="text-xs text-outline">Go</span>
              </div>
              <p className="text-xs text-on-surface-variant line-clamp-2">
                等待数据库连接测试...
              </p>
              <div className="flex items-center justify-between mt-3">
                <div className="text-xs text-outline">等待中</div>
                <div className="text-xs font-mono text-outline">842 tokens</div>
              </div>
            </div>

            <div className="bg-surface-container rounded-lg p-3">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-success"></div>
                  <span className="text-xs font-medium text-on-surface">Architect</span>
                </div>
                <span className="text-xs text-outline">完成</span>
              </div>
              <p className="text-xs text-on-surface-variant line-clamp-2">
                项目架构设计已完成
              </p>
              <div className="flex items-center justify-between mt-3">
                <div className="text-xs text-outline">100%</div>
                <div className="text-xs font-mono text-outline">3,521 tokens</div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'thought' && (
          <div className="space-y-4">
            <div className="text-center py-8">
              <span className="material-symbols-outlined text-3xl text-outline/50 mb-3">
                psychology
              </span>
              <p className="text-sm text-outline">选择 Agent 以查看思维链</p>
            </div>
          </div>
        )}

        {activeTab === 'files' && (
          <div className="space-y-4">
            <div className="text-center py-8">
              <span className="material-symbols-outlined text-3xl text-outline/50 mb-3">
                difference
              </span>
              <p className="text-sm text-outline">暂无文件变更</p>
              <p className="text-xs text-outline/70 mt-1">文件变更将在此处显示</p>
            </div>
          </div>
        )}
      </div>

      {/* Panel Footer */}
      <div className="px-4 py-3 border-t border-outline-variant/15">
        <div className="flex items-center justify-between text-xs text-outline">
          <span>系统运行中</span>
          <span className="flex items-center gap-1">
            <span className="w-1.5 h-1.5 rounded-full bg-success animate-pulse"></span>
            正常
          </span>
        </div>
      </div>
    </div>
  );
};

export default RightPanel;