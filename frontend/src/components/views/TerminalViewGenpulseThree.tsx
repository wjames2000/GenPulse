import React from 'react';
import LayoutGenpulse from '../LayoutGenpulse';

const TerminalViewGenpulseThree: React.FC = () => {
  return (
    <LayoutGenpulse>
      <div className="p-6">
        <div className="mb-4">
          <h2 className="text-2xl font-bold text-on-surface">终端输出</h2>
          <p className="text-on-surface-variant text-sm">简易终端输出面板，与三栏模板对齐</p>
        </div>
        <div className="bg-surface-container-low rounded-xl p-6 border border-outline-variant/15 font-mono text-sm h-80 overflow-auto">$ ls -la
          index.html
          main.ts
          vite.config.ts
        </div>
      </div>
    </LayoutGenpulse>
  );
};

export default TerminalViewGenpulseThree;
