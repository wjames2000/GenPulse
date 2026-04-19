import React, { useState } from 'react';

type Skill = {
  id: string;
  name: string;
  description: string;
  category: string;
  version: string;
  trigger: string;
  usageCount: number;
  successRate: number;
  enabled: boolean;
  lastUsed: string;
  tools: string[];
};

type SkillDetail = {
  id: string;
  name: string;
  description: string;
  content: string;
  version: string;
  author: string;
  createdAt: string;
  updatedAt: string;
  usageStats: {
    totalUses: number;
    successRate: number;
    avgTokensSaved: number;
  };
  tools: string[];
  triggers: string[];
};

const SkillsViewGenpulse: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'list' | 'stats'>('list');
  const [selectedSkill, setSelectedSkill] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('all');

  const skills: Skill[] = [
    {
      id: '1',
      name: '代码生成',
      description: '根据需求描述生成高质量的代码',
      category: '开发',
      version: '1.2.0',
      trigger: '当用户描述代码需求时',
      usageCount: 245,
      successRate: 92,
      enabled: true,
      lastUsed: '2小时前',
      tools: ['code_generation', 'context_analysis']
    },
    {
      id: '2',
      name: '代码审查',
      description: '分析代码质量，识别潜在问题',
      category: '质量',
      version: '1.1.3',
      trigger: '当代码提交时',
      usageCount: 189,
      successRate: 88,
      enabled: true,
      lastUsed: '1天前',
      tools: ['code_analysis', 'security_check']
    },
    {
      id: '3',
      name: '测试生成',
      description: '自动生成单元测试和集成测试',
      category: '测试',
      version: '1.0.8',
      trigger: '当功能开发完成时',
      usageCount: 156,
      successRate: 85,
      enabled: true,
      lastUsed: '3小时前',
      tools: ['test_generation', 'coverage_analysis']
    },
    {
      id: '4',
      name: '文档生成',
      description: '为代码库生成详细的API文档',
      category: '文档',
      version: '1.3.2',
      trigger: '当API变更时',
      usageCount: 98,
      successRate: 95,
      enabled: true,
      lastUsed: '1周前',
      tools: ['doc_generation', 'api_parsing']
    },
    {
      id: '5',
      name: '架构设计',
      description: '设计系统架构和技术选型',
      category: '设计',
      version: '2.0.1',
      trigger: '当开始新项目时',
      usageCount: 76,
      successRate: 90,
      enabled: true,
      lastUsed: '2天前',
      tools: ['architecture_planning', 'tech_evaluation']
    },
    {
      id: '6',
      name: '性能优化',
      description: '识别性能瓶颈并提供优化建议',
      category: '优化',
      version: '1.4.0',
      trigger: '当应用性能下降时',
      usageCount: 64,
      successRate: 82,
      enabled: false,
      lastUsed: '3天前',
      tools: ['performance_analysis', 'optimization_suggestions']
    },
    {
      id: '7',
      name: '安全审计',
      description: '检测安全漏洞和潜在风险',
      category: '安全',
      version: '1.5.2',
      trigger: '当代码包含安全敏感操作时',
      usageCount: 112,
      successRate: 96,
      enabled: true,
      lastUsed: '5小时前',
      tools: ['security_scan', 'vulnerability_detection']
    },
    {
      id: '8',
      name: '部署自动化',
      description: '自动化部署流程和CI/CD配置',
      category: '运维',
      version: '1.2.4',
      trigger: '当项目需要部署时',
      usageCount: 87,
      successRate: 94,
      enabled: false,
      lastUsed: '1月前',
      tools: ['deployment_automation', 'ci_cd_config']
    }
  ];

  const categories = ['all', '开发', '质量', '测试', '文档', '设计', '优化', '安全', '运维'];

  const filteredSkills = skills.filter(skill => {
    const matchesSearch = skill.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         skill.description.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesCategory = selectedCategory === 'all' || skill.category === selectedCategory;
    return matchesSearch && matchesCategory;
  });

  const handleSkillToggle = (skillId: string, enabled: boolean) => {
    console.log(`Skill ${skillId} ${enabled ? 'enabled' : 'disabled'}`);
    // 在实际应用中，这里会调用API更新技能状态
  };

  const getSuccessRateColor = (rate: number) => {
    if (rate >= 90) return 'text-success';
    if (rate >= 80) return 'text-warning';
    return 'text-error';
  };

  return (
    <div className="h-full overflow-y-auto p-6">
      {/* Page Header */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-on-surface mb-2">技能库管理</h1>
            <p className="text-sm text-on-surface-variant">
              管理 Skills 技能库，查看使用统计，编辑或创建新 Skill
            </p>
          </div>
          <button className="bg-gradient-to-br from-primary-container to-inverse-primary text-on-primary-container hover:brightness-110 transition-all duration-200 py-2 px-4 rounded-lg flex items-center justify-center gap-2 font-medium text-sm">
            <span className="material-symbols-outlined text-sm">add</span>
            创建新技能
          </button>
        </div>

        {/* Tab Navigation */}
        <div className="flex items-center border-b border-outline-variant/15 mb-6">
          <button
            onClick={() => setActiveTab('list')}
            className={`px-6 py-3 text-sm font-medium transition-colors border-b-2 ${
              activeTab === 'list'
                ? 'text-primary border-primary'
                : 'text-outline border-transparent hover:text-on-surface-variant'
            }`}
          >
            <span className="flex items-center gap-2">
              <span className="material-symbols-outlined text-base">apps</span>
              技能库列表
            </span>
          </button>
          <button
            onClick={() => setActiveTab('stats')}
            className={`px-6 py-3 text-sm font-medium transition-colors border-b-2 ${
              activeTab === 'stats'
                ? 'text-primary border-primary'
                : 'text-outline border-transparent hover:text-on-surface-variant'
            }`}
          >
            <span className="flex items-center gap-2">
              <span className="material-symbols-outlined text-base">monitoring</span>
              使用统计
            </span>
          </button>
        </div>
      </div>

      {activeTab === 'list' ? (
        <>
          {/* Filter and Search Bar */}
          <div className="mb-8">
            <div className="flex flex-col md:flex-row gap-4 items-center justify-between">
              <div className="flex-1 w-full">
                <div className="relative">
                  <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-outline text-sm">
                    search
                  </span>
                  <input
                    className="w-full bg-surface-container-lowest border-none rounded-lg pl-10 pr-4 py-2.5 text-sm text-on-surface placeholder:text-outline focus:ring-0 focus:border-b-2 focus:border-primary outline-none transition-all"
                    placeholder="搜索技能名称或描述..."
                    type="text"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div className="text-sm text-on-surface-variant">分类:</div>
                <div className="flex flex-wrap gap-2">
                  {categories.map((category) => (
                    <button
                      key={category}
                      onClick={() => setSelectedCategory(category)}
                      className={`px-3 py-1.5 text-xs rounded-lg transition-colors ${
                        selectedCategory === category
                          ? 'bg-primary-container text-on-primary-container'
                          : 'bg-surface-container-high text-on-surface-variant hover:bg-surface-container-highest'
                      }`}
                    >
                      {category === 'all' ? '全部' : category}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </div>

          {/* Skills Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredSkills.map((skill) => (
              <div
                key={skill.id}
                className={`bg-surface-container rounded-xl p-5 transition-all duration-200 ${
                  selectedSkill === skill.id ? 'ring-2 ring-primary/30' : ''
                }`}
                onClick={() => setSelectedSkill(skill.id)}
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <div className={`p-2 rounded-lg ${
                      skill.enabled ? 'bg-primary/10 text-primary' : 'bg-surface-container-high text-outline'
                    }`}>
                      <span className="material-symbols-outlined">memory</span>
                    </div>
                    <div>
                      <h3 className="text-lg font-semibold text-on-surface">{skill.name}</h3>
                      <div className="flex items-center gap-2">
                        <span className="text-xs px-2 py-0.5 rounded bg-surface-container-high text-outline">
                          {skill.category}
                        </span>
                        <span className="text-xs text-outline">v{skill.version}</span>
                      </div>
                    </div>
                  </div>
                  <div className="relative">
                    <div className={`w-10 h-5 rounded-full transition-all ${
                      skill.enabled ? 'bg-primary' : 'bg-surface-container-high'
                    }`}>
                      <div className={`w-5 h-5 rounded-full bg-white transition-all ${
                        skill.enabled ? 'translate-x-5' : ''
                      }`}
                      onClick={(e) => {
                        e.stopPropagation();
                        handleSkillToggle(skill.id, !skill.enabled);
                      }}
                      ></div>
                    </div>
                  </div>
                </div>

                <p className="text-sm text-on-surface-variant mb-4 line-clamp-2">{skill.description}</p>

                <div className="space-y-3">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-outline">触发条件</span>
                    <span className="text-on-surface-variant text-right max-w-[60%] truncate">{skill.trigger}</span>
                  </div>

                  <div className="flex items-center justify-between text-sm">
                    <span className="text-outline">使用次数</span>
                    <span className="font-mono text-on-surface">{skill.usageCount.toLocaleString()}</span>
                  </div>

                  <div className="flex items-center justify-between text-sm">
                    <span className="text-outline">成功率</span>
                    <span className={`font-mono ${getSuccessRateColor(skill.successRate)}`}>
                      {skill.successRate}%
                    </span>
                  </div>

                  <div className="flex flex-wrap gap-1 mt-2">
                    {skill.tools.slice(0, 2).map((tool, idx) => (
                      <span key={idx} className="text-xs px-2 py-1 rounded bg-surface-container-high text-outline">
                        {tool.replace('_', ' ')}
                      </span>
                    ))}
                    {skill.tools.length > 2 && (
                      <span className="text-xs px-2 py-1 rounded bg-surface-container-high text-outline">
                        +{skill.tools.length - 2}
                      </span>
                    )}
                  </div>

                  <div className="flex items-center justify-between pt-3 border-t border-outline-variant/15 text-xs text-outline">
                    <span>最后使用: {skill.lastUsed}</span>
                    <button className="text-primary hover:text-primary-container transition-colors">
                      编辑 →
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* Empty State */}
          {filteredSkills.length === 0 && (
            <div className="text-center py-16">
              <div className="w-20 h-20 rounded-full bg-surface-container mx-auto mb-6 flex items-center justify-center">
                <span className="material-symbols-outlined text-3xl text-outline">search_off</span>
              </div>
              <h3 className="text-lg font-medium text-on-surface mb-2">未找到匹配的技能</h3>
              <p className="text-sm text-on-surface-variant">尝试其他搜索关键词或分类</p>
            </div>
          )}
        </>
      ) : (
        /* Statistics View (设计规范 5.4.3) */
        <div className="space-y-8">
          {/* Usage Trends */}
          <div className="bg-surface-container rounded-xl p-6">
            <h2 className="text-lg font-semibold text-on-surface mb-4">技能使用趋势</h2>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              <div>
                <div className="h-64 bg-surface-container-low rounded-lg flex items-center justify-center">
                  <div className="text-center">
                    <span className="material-symbols-outlined text-4xl text-outline/50 mb-3">
                      trending_up
                    </span>
                    <p className="text-sm text-outline">使用趋势图表</p>
                  </div>
                </div>
              </div>
              <div className="space-y-6">
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm text-on-surface-variant">总使用次数</span>
                    <span className="text-2xl font-bold text-on-surface">1,027</span>
                  </div>
                  <div className="h-2 bg-surface-container-low rounded-full overflow-hidden">
                    <div className="h-full bg-primary rounded-full w-3/4"></div>
                  </div>
                </div>
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm text-on-surface-variant">平均成功率</span>
                    <span className="text-2xl font-bold text-success">89%</span>
                  </div>
                  <div className="h-2 bg-surface-container-low rounded-full overflow-hidden">
                    <div className="h-full bg-success rounded-full w-[89%]"></div>
                  </div>
                </div>
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm text-on-surface-variant">平均节省 Token</span>
                    <span className="text-2xl font-bold text-on-surface">1,248</span>
                  </div>
                  <div className="h-2 bg-surface-container-low rounded-full overflow-hidden">
                    <div className="h-full bg-running rounded-full w-2/3"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Top Skills */}
          <div className="bg-surface-container rounded-xl p-6">
            <h2 className="text-lg font-semibold text-on-surface mb-4">最常用技能</h2>
            <div className="space-y-4">
              {skills
                .sort((a, b) => b.usageCount - a.usageCount)
                .slice(0, 5)
                .map((skill, index) => (
                  <div key={skill.id} className="flex items-center justify-between p-3 rounded-lg bg-surface-container-high">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center">
                        <span className="text-xs font-bold text-primary">{index + 1}</span>
                      </div>
                      <div>
                        <div className="text-sm font-medium text-on-surface">{skill.name}</div>
                        <div className="text-xs text-outline">{skill.category}</div>
                      </div>
                    </div>
                    <div className="flex items-center gap-6">
                      <div className="text-right">
                        <div className="text-sm font-mono text-on-surface">{skill.usageCount.toLocaleString()}</div>
                        <div className="text-xs text-outline">使用次数</div>
                      </div>
                      <div className="text-right">
                        <div className={`text-sm font-mono ${getSuccessRateColor(skill.successRate)}`}>
                          {skill.successRate}%
                        </div>
                        <div className="text-xs text-outline">成功率</div>
                      </div>
                    </div>
                  </div>
                ))}
            </div>
          </div>

          {/* Recent Activities */}
          <div className="bg-surface-container rounded-xl p-6">
            <h2 className="text-lg font-semibold text-on-surface mb-4">最近执行记录</h2>
            <div className="space-y-3">
              {[
                { skill: '代码生成', agent: 'Frontend Dev', result: '成功', time: '2分钟前', tokens: 1248 },
                { skill: '代码审查', agent: 'Architect', result: '成功', time: '15分钟前', tokens: 876 },
                { skill: '测试生成', agent: 'QA Engineer', result: '成功', time: '1小时前', tokens: 1542 },
                { skill: '安全审计', agent: 'Security Agent', result: '警告', time: '3小时前', tokens: 2103 },
                { skill: '文档生成', agent: 'Documentation Agent', result: '成功', time: '5小时前', tokens: 987 },
              ].map((activity, index) => (
                <div key={index} className="flex items-center justify-between py-3 border-b border-outline-variant/15 last:border-0">
                  <div className="flex items-center gap-3">
                    <div className={`p-1.5 rounded-md ${
                      activity.result === '成功' ? 'bg-success/10' : 'bg-warning/10'
                    }`}>
                      <span className={`material-symbols-outlined text-sm ${
                        activity.result === '成功' ? 'text-success' : 'text-warning'
                      }`}>
                        {activity.result === '成功' ? 'check_circle' : 'warning'}
                      </span>
                    </div>
                    <div>
                      <p className="text-sm font-medium text-on-surface">{activity.skill}</p>
                      <p className="text-xs text-outline">由 {activity.agent} 执行</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-xs font-mono text-outline">{activity.tokens.toLocaleString()} tokens</div>
                    <div className="text-xs text-outline">{activity.time}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default SkillsViewGenpulse;