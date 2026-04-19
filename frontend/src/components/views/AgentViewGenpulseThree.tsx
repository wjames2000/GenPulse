import React from 'react';
import LayoutGenpulse from '../LayoutGenpulse';

const AgentViewGenpulseThree: React.FC = () => {
  const agents = [
    { id:1, name:'Product Manager', role:'产品经理', status:'running', progress:60 },
    { id:2, name:'Architect', role:'架构师', status:'completed', progress:100 },
  ];
  return (
    <LayoutGenpulse>
      <div className="p-6">
        <div className="mb-4">
          <h2 className="text-2xl font-bold text-on-surface">代理管理</h2>
          <p className="text-on-surface-variant text-sm">统一三栏模板下的代理视图</p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {agents.map(a => (
            <div key={a.id} className="bg-surface-container-high rounded-xl p-4 border border-outline-variant/15">
              <div className="flex justify-between items-center mb-2">
                <div className="flex items-center gap-2">
                  <span className="material-symbols-outlined text-primary">person</span>
                  <span className="text-on-surface font-medium">{a.name}</span>
                </div>
                <span className="text-xs text-on-surface-variant">{a.role}</span>
              </div>
              <div className="h-1.5 bg-surface-container-high rounded-full mb-2">
                <div className={`h-full rounded-full ${a.status==='running'?'bg-running':'bg-success'}`} style={{ width: `${a.progress}%` }} />
              </div>
              <div className="text-xs text-on-surface-variant">状态: {a.status}</div>
            </div>
          ))}
        </div>
      </div>
    </LayoutGenpulse>
  );
};

export default AgentViewGenpulseThree;
