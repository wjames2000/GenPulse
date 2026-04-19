import React from 'react';
import LayoutGenpulse from '../LayoutGenpulse';

const MemoryViewGenpulseThree: React.FC = () => {
  return (
    <LayoutGenpulse>
      <div className="p-6">
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-on-surface mb-2">神经资产</h2>
          <p className="text-on-surface-variant text-sm">统一的三栏布局示例：情节记忆、最近记忆、分类等</p>
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <section className="bg-surface-container-low rounded-xl p-6 border border-outline-variant/10">
            <h3 className="text-lg font-semibold text-on-surface mb-4">记忆概览</h3>
            <div className="grid grid-cols-2 gap-4">
              <div className="p-4 bg-surface-container-high rounded-lg">总记忆数<br/><span className="text-xl">1,247</span></div>
              <div className="p-4 bg-surface-container-high rounded-lg">最近新增<br/><span className="text-xl">42</span></div>
            </div>
          </section>
          <section className="bg-surface-container-low rounded-xl p-6 border border-outline-variant/10">
            <h3 className="text-lg font-semibold text-on-surface mb-4">最近记忆</h3>
            {Array.from({length:4}).map((_,i)=> (
              <div key={i} className="flex items-center justify-between mb-3 p-2 rounded-lg hover:bg-surface-container-high/20">
                <div className="flex items-center gap-2">
                  <span className="material-symbols-outlined text-primary">memory</span>
                  <span className="text-sm">记忆项 #{i+1}</span>
                </div>
                <span className="text-xs text-on-surface-variant">2h前</span>
              </div>
            ))}
          </section>
        </div>
      </div>
    </LayoutGenpulse>
  );
};

export default MemoryViewGenpulseThree;
