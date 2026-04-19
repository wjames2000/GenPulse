import React from 'react';
import LayoutGenpulse from '../LayoutGenpulse';

const SettingsViewDesignThree: React.FC = () => {
  return (
    <LayoutGenpulse>
      <div className="p-6">
        <div className="mb-4">
          <h2 className="text-2xl font-bold text-on-surface">系统设置</h2>
          <p className="text-on-surface-variant text-sm">通过三栏模板统一管理全局设置、主题、语言等</p>
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
          <section className="bg-surface-container-low rounded-xl p-4 border border-outline-variant/15">
            <h3 className="text-sm font-semibold text-on-surface mb-2">外观</h3>
            <div>主题：暗色模式</div>
          </section>
          <section className="bg-surface-container-low rounded-xl p-4 border border-outline-variant/15">
            <h3 className="text-sm font-semibold text-on-surface mb-2">语言</h3>
            <div>语言：中文</div>
          </section>
          <section className="bg-surface-container-low rounded-xl p-4 border border-outline-variant/15">
            <h3 className="text-sm font-semibold text-on-surface mb-2">日志等级</h3>
            <div>Info</div>
          </section>
        </div>
      </div>
    </LayoutGenpulse>
  );
};

export default SettingsViewDesignThree;
