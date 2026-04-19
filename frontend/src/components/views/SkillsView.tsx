import React from 'react';

const SkillsView: React.FC = () => {
  const skills = [
    { name: '代码生成', icon: 'code', category: '开发', status: 'active' },
    { name: '代码审查', icon: 'rate_review', category: '质量', status: 'active' },
    { name: '测试生成', icon: 'bug_report', category: '测试', status: 'active' },
    { name: '文档生成', icon: 'description', category: '文档', status: 'inactive' },
    { name: '架构设计', icon: 'architecture', category: '设计', status: 'active' },
    { name: '性能优化', icon: 'speed', category: '优化', status: 'inactive' },
    { name: '安全审计', icon: 'security', category: '安全', status: 'active' },
    { name: '部署自动化', icon: 'deployed_code', category: '运维', status: 'inactive' },
  ];

  return (
    <div className="p-6 md:p-8 lg:p-12">
      <div className="max-w-6xl mx-auto">
        <div className="flex items-center justify-between mb-10">
          <div>
            <h2 className="text-[2.25rem] font-bold tracking-[-0.02em] text-on-surface mb-2">
              技能库
            </h2>
            <p className="text-on-surface-variant text-sm">
              管理和配置AI代理的可用技能集合。
            </p>
          </div>
          <button className="bg-primary-container text-on-primary-container hover:bg-inverse-primary rounded-full px-6 py-2 text-sm font-medium shadow-lg shadow-primary/20 transition-all flex items-center gap-2">
            <span className="material-symbols-outlined">add</span>
            添加技能
          </button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {skills.map((skill, index) => (
            <div key={index} className="bg-surface-container rounded-xl p-6 hover:bg-surface-container-high transition-colors cursor-pointer group">
              <div className="flex items-center justify-between mb-4">
                <div className="p-3 rounded-lg bg-surface-container-low text-primary">
                  <span className="material-symbols-outlined">{skill.icon}</span>
                </div>
                <div className={`px-2 py-1 rounded-md text-xs font-mono uppercase ${
                  skill.status === 'active' 
                    ? 'bg-sys-success/10 text-sys-success' 
                    : 'bg-error/10 text-error'
                }`}>
                  {skill.status === 'active' ? '活跃' : '未激活'}
                </div>
              </div>
              
              <h3 className="text-lg font-semibold text-on-surface mb-1">{skill.name}</h3>
              <p className="text-sm text-outline mb-4">{skill.category}</p>
              
              <div className="flex items-center justify-between">
                <span className="text-xs text-outline">点击配置</span>
                <span className="material-symbols-outlined text-outline group-hover:text-primary transition-colors">
                  settings
                </span>
              </div>
            </div>
          ))}
        </div>

        {/* 技能统计 */}
        <div className="mt-12 grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-surface-container rounded-xl p-6">
            <h3 className="text-lg font-semibold text-on-surface mb-4">技能统计</h3>
            <div className="space-y-4">
              <div>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-outline">活跃技能</span>
                  <span className="font-mono text-sys-success">5/8</span>
                </div>
                <div className="h-2 bg-surface-container-low rounded-full overflow-hidden">
                  <div className="h-full bg-sys-success rounded-full" style={{ width: '62.5%' }}></div>
                </div>
              </div>
              <div>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-outline">使用频率</span>
                  <span className="font-mono text-primary">高</span>
                </div>
                <div className="h-2 bg-surface-container-low rounded-full overflow-hidden">
                  <div className="h-full bg-primary rounded-full" style={{ width: '80%' }}></div>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-surface-container rounded-xl p-6 md:col-span-2">
            <h3 className="text-lg font-semibold text-on-surface mb-4">最近活动</h3>
            <div className="space-y-3">
              {[
                { skill: '代码生成', time: '2分钟前', agent: 'Frontend Dev' },
                { skill: '代码审查', time: '15分钟前', agent: 'Architect' },
                { skill: '测试生成', time: '1小时前', agent: 'QA Agent' },
              ].map((activity, index) => (
                <div key={index} className="flex items-center justify-between py-2 border-b border-outline-variant/15 last:border-0">
                  <div className="flex items-center gap-3">
                    <div className="p-1.5 rounded-md bg-surface-container-low">
                      <span className="material-symbols-outlined text-sm text-primary">bolt</span>
                    </div>
                    <div>
                      <p className="text-sm font-medium text-on-surface">{activity.skill}</p>
                      <p className="text-xs text-outline">由 {activity.agent} 执行</p>
                    </div>
                  </div>
                  <span className="text-xs text-outline">{activity.time}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SkillsView;