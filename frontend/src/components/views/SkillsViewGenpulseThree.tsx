import React, { useMemo, useState } from 'react';

type Skill = {
  id: string; name: string; description: string; category: string; version: string;
  trigger: string; usageCount: number; successRate: number; enabled: boolean; lastUsed: string; tools: string[];
};

const SkillsViewGenpulseThree: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const categories = ['all','开发','质量','测试','文档','设计','优化','安全','运维'];

  const skills: Skill[] = useMemo(() => [
    { id:'1', name:'代码生成', description:'根据需求描述生成高质量的代码', category:'开发', version:'1.2.0', trigger:'用户描述代码需求时', usageCount:245, successRate:92, enabled:true, lastUsed:'2小时前', tools:['code_generation','context_analysis'] },
    { id:'2', name:'代码审查', description:'分析代码质量，识别潜在问题', category:'质量', version:'1.1.3', trigger:'代码提交时', usageCount:189, successRate:88, enabled:true, lastUsed:'1天前', tools:['code_analysis','security_check'] },
    { id:'3', name:'测试生成', description:'自动生成单元测试与集成测试', category:'测试', version:'1.0.8', trigger:'功能完成时', usageCount:156, successRate:85, enabled:true, lastUsed:'3小时前', tools:['test_generation','coverage_analysis'] },
  ],[]);

  const filteredSkills = skills.filter(s => {
    const m1 = s.name.toLowerCase().includes(searchQuery.toLowerCase()) || s.description.toLowerCase().includes(searchQuery.toLowerCase());
    const m2 = selectedCategory==='all' || s.category===selectedCategory;
    return m1 && m2;
  });

  return (
    <div className="h-full overflow-y-auto p-6">
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-on-surface mb-2">技能库三栏视图</h2>
        <p className="text-on-surface-variant text-sm">浏览、筛选并管理技能库。</p>
      </div>

      <div className="flex gap-4 mb-4">
        <div className="flex-1 rounded-lg bg-surface-container-low p-3">
          <div className="flex items-center gap-2">
            <span className="material-symbols-outlined text-outline">search</span>
            <input className="bg-transparent outline-none placeholder-on-surface-variant" placeholder="搜索技能..." value={searchQuery} onChange={e=>setSearchQuery(e.target.value)} />
          </div>
        </div>
        <div className="text-sm text-on-surface-variant flex items-center gap-2">
          分类:
          {categories.map(c => (
            <button key={c} onClick={()=>setSelectedCategory(c)} className={`px-2 py-1 rounded ${selectedCategory===c?'bg-primary-container text-on-primary-container':'bg-surface-container-high text-on-surface-variant'}`}>{c==='all'?'All':c}</button>
          ))}
        </div>
      </div>

      {/* Main grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredSkills.map(s => (
          <div key={s.id} className="bg-surface-container rounded-xl p-5 border border-outline-variant/15 hover:bg-surface-container-high transition-colors">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-2">
                <span className="material-symbols-outlined text-primary">memory</span>
                <div>
                  <div className="font-medium text-on-surface">{s.name}</div>
                  <div className="text-xs text-outline">{s.category} • v{s.version}</div>
                </div>
              </div>
              <span className={`text-xs px-2 py-1 rounded ${s.enabled?'bg-primary/10 text-primary':'bg-surface-container-high text-on-surface-variant'}`}>{s.enabled?'Enabled':'Disabled'}</span>
            </div>
            <p className="text-sm text-on-surface-variant mb-2">{s.description}</p>
            <div className="text-xs text-outline">触发: {s.trigger}</div>
            <div className="flex items-center justify-between mt-3 text-xs text-outline">
              <span>使用: {s.usageCount.toLocaleString()}</span>
              <span>成功率: {s.successRate}%</span>
            </div>
          </div>
        ))}
      </div>

      {/* Right Panel */}
      <aside className="hidden xl:block absolute right-0 top-0 h-full w-[360px] p-4">
        <div className="bg-surface-container-low h-full rounded-xl p-4 border border-outline-variant/15">
          <div className="font-semibold text-on-surface mb-2">Top Skills</div>
          {skills.slice(0,3).map((sk, idx) => (
            <div key={sk.id} className="flex items-center justify-between mb-2 text-sm">
              <span>{idx+1}. {sk.name}</span>
              <span className="text-xs text-outline">{sk.successRate}%</span>
            </div>
          ))}
        </div>
      </aside>
    </div>
  );
};

export default SkillsViewGenpulseThree;
