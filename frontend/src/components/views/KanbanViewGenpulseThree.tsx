import React from 'react';
import LayoutGenpulse from '../LayoutGenpulse';

const KanbanViewGenpulseThree: React.FC = () => {
  return (
    <LayoutGenpulse>
      <div className="p-6">
        <div className="mb-4">
          <h2 className="text-2xl font-bold text-on-surface">执行看板</h2>
          <p className="text-on-surface-variant text-sm">三列看板示例：待办、进行中、已完成</p>
        </div>
        <div className="flex gap-4">
          {['To Do','In Progress','Done'].map((col, idx)=> (
            <div key={idx} className="flex-1 bg-surface-container-low rounded-xl p-4 border border-outline-variant/15">
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs uppercase text-primary font-semibold">{col}</span>
              </div>
              <div className="space-y-3">
                {[0,1].map(i => (
                  <div key={i} className="p-2 rounded-lg bg-surface-container-high hover:bg-surface-container-highest">
                    <div className="flex items-center justify-between text-xs">
                      <span>任务 #{i+1}</span>
                      <span className="text-xs text-on-surface-variant">2m</span>
                    </div>
                    <div className="text-sm text-on-surface-variant mt-1">描述简要</div>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </LayoutGenpulse>
  );
};

export default KanbanViewGenpulseThree;
