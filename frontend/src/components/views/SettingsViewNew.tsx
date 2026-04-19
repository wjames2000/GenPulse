import React, { useState, useEffect } from 'react';
import { useAppStore } from '../../stores/appStore';

// 定义设计系统中的颜色常量
const DESIGN_COLORS = {
  surface: '#131318',
  'surface-container': '#1f1f25',
  'surface-container-low': '#1b1b20',
  'surface-container-high': '#2a292f',
  'surface-container-highest': '#35343a',
  'surface-container-lowest': '#0e0e13',
  'surface-variant': '#35343a',
  'surface-bright': '#39383e',
  'on-surface': '#e4e1e9',
  'on-surface-variant': '#c6c4d8',
  'on-background': '#e4e1e9',
  primary: '#c0c1ff',
  'primary-container': '#5b5fff',
  'on-primary': '#0d00aa',
  'on-primary-container': '#fffcff',
  secondary: '#c0c1ff',
  'secondary-container': '#3f4287',
  'on-secondary': '#26286c',
  'on-secondary-container': '#afb2ff',
  outline: '#908fa1',
  'outline-variant': '#454555',
  error: '#ffb4ab',
  'error-container': '#93000a',
  'on-error': '#690005',
  'on-error-container': '#ffdad6',
  'sys-success': '#4ADE80',
};

interface ModelProvider {
  id: string;
  name: string;
  icon: string;
  endpoint: string;
  status: 'connected' | 'disconnected' | 'unconfigured';
  defaultModel: string;
  availableModels: string[];
}

interface AgentRole {
  id: string;
  name: string;
  description: string;
  icon: string;
  status: 'active' | 'inactive';
  assignedModel: string;
  systemPrompt: string;
  availableModels: string[];
}

const SettingsViewNew: React.FC = () => {
  const [theme, setTheme] = useState<'light' | 'dark' | 'auto'>('dark');
  const [language, setLanguage] = useState('zh');
  const [logLevel, setLogLevel] = useState('info');
  
  const [modelProviders, setModelProviders] = useState<ModelProvider[]>([
    {
      id: 'openai',
      name: 'OpenAI',
      icon: 'data_object',
      endpoint: 'api.openai.com/v1',
      status: 'connected',
      defaultModel: 'gpt-4-turbo-preview',
      availableModels: ['gpt-4-turbo-preview', 'gpt-4', 'gpt-3.5-turbo'],
    },
    {
      id: 'anthropic',
      name: 'Anthropic Claude',
      icon: 'psychology_alt',
      endpoint: 'api.anthropic.com',
      status: 'connected',
      defaultModel: 'claude-3-opus-20240229',
      availableModels: ['claude-3-opus-20240229', 'claude-3-sonnet-20240229'],
    },
    {
      id: 'gemini',
      name: 'Google Gemini',
      icon: 'auto_awesome',
      endpoint: 'generativelanguage.googleapis',
      status: 'unconfigured',
      defaultModel: 'gemini-1.5-pro',
      availableModels: ['gemini-1.5-pro'],
    },
  ]);

  const [agentRoles, setAgentRoles] = useState<AgentRole[]>([
    {
      id: 'orchestrator',
      name: 'Orchestrator',
      description: '任务路由与分解',
      icon: 'account_tree',
      status: 'active',
      assignedModel: 'gpt-4-turbo-preview',
      systemPrompt: '你是 Orchestrator。你的主要目标是将复杂的用户请求分解为特定子代理可执行的独立步骤...',
      availableModels: ['gpt-4-turbo-preview', 'gpt-4', 'gpt-3.5-turbo'],
    },
    {
      id: 'architect',
      name: 'Architect',
      description: '系统设计与架构',
      icon: 'architecture',
      status: 'active',
      assignedModel: 'claude-3-opus-20240229',
      systemPrompt: '作为首席架构师，评估提出的需求并设计稳健、可扩展的数据模型。优先考虑关注点分离和...',
      availableModels: ['claude-3-opus-20240229', 'claude-3-sonnet-20240229'],
    },
    {
      id: 'frontend',
      name: 'Frontend Dev',
      description: 'UI 生成与逻辑',
      icon: 'web',
      status: 'active',
      assignedModel: 'gpt-4',
      systemPrompt: '扮演一位利用现代 Web 标准的专家前端开发者。\n\n在生成 UI 组件时：\n- 严格遵守提供的设计系统规范。\n- 实施优先考虑移动端工作流的响应式布局。\n- 确保所有交互状态（hover、focus、active）都得到明确定义。\n- 输出整洁、语义化的 HTML5 结构。',
      availableModels: ['gpt-4-turbo-preview', 'gpt-4', 'gpt-3.5-turbo'],
    },
  ]);

  const [hasChanges, setHasChanges] = useState(false);

  const handleThemeChange = (newTheme: 'light' | 'dark' | 'auto') => {
    setTheme(newTheme);
    setHasChanges(true);
  };

  const handleLanguageChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setLanguage(e.target.value);
    setHasChanges(true);
  };

  const handleLogLevelChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setLogLevel(e.target.value);
    setHasChanges(true);
  };

  const handleModelChange = (providerId: string, model: string) => {
    setModelProviders(prev => prev.map(p => 
      p.id === providerId ? { ...p, defaultModel: model } : p
    ));
    setHasChanges(true);
  };

  const handleAgentModelChange = (roleId: string, model: string) => {
    setAgentRoles(prev => prev.map(r => 
      r.id === roleId ? { ...r, assignedModel: model } : r
    ));
    setHasChanges(true);
  };

  const handleSave = () => {
    console.log('Settings saved:', { theme, language, logLevel, modelProviders, agentRoles });
    setHasChanges(false);
  };

  const handleDiscard = () => {
    // 重置为默认值
    setTheme('dark');
    setLanguage('zh');
    setLogLevel('info');
    setHasChanges(false);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'connected':
      case 'active':
        return DESIGN_COLORS['sys-success'];
      case 'disconnected':
        return DESIGN_COLORS.error;
      case 'unconfigured':
        return DESIGN_COLORS.error;
      default:
        return DESIGN_COLORS.outline;
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'connected':
        return '已连接';
      case 'disconnected':
        return '未连接';
      case 'unconfigured':
        return '未配置';
      case 'active':
        return '活跃';
      case 'inactive':
        return '未激活';
      default:
        return '未知';
    }
  };

  return (
    <div className="h-screen overflow-hidden flex">
      {/* 主内容区域 */}
      <div className="flex-1 ml-0 md:ml-64 flex flex-col h-screen overflow-hidden bg-surface relative">
        {/* 顶部应用栏 */}
        <header className="flex justify-between items-center w-full px-6 py-3 bg-[#131318]/80 backdrop-blur-xl shadow-2xl shadow-black/40 docked full-width top-0 z-50 sticky">
          {/* 左侧：品牌/搜索 */}
          <div className="flex items-center gap-6 flex-1">
            <span className="text-xl font-bold tracking-tight text-[#C0C1FF] font-['Inter'] antialiased hidden lg:block">
              Genpulse AI
            </span>
            <div className="relative w-full max-w-xs md:w-64">
              <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#908FA1] text-sm">
                search
              </span>
              <input
                className="w-full bg-surface-container-lowest text-on-surface text-sm rounded-lg pl-9 pr-3 py-1.5 border-none focus:ring-0 focus:border-b-2 focus:border-b-primary transition-all placeholder:text-outline outline-none"
                placeholder="搜索设置..."
                type="text"
              />
            </div>
          </div>
          
          {/* 右侧：操作和用户 */}
          <div className="flex items-center gap-4">
            <div className="hidden md:flex items-center gap-1 text-[#C0C1FF]">
              <button className="p-2 rounded-full hover:bg-[#35343a]/50 transition-colors duration-200 active:scale-95 transition-transform flex items-center justify-center">
                <span className="material-symbols-outlined">notifications</span>
              </button>
              <button className="p-2 rounded-full hover:bg-[#35343a]/50 transition-colors duration-200 active:scale-95 transition-transform flex items-center justify-center">
                <span className="material-symbols-outlined">settings</span>
              </button>
              <button className="p-2 rounded-full hover:bg-[#35343a]/50 transition-colors duration-200 active:scale-95 transition-transform flex items-center justify-center">
                <span className="material-symbols-outlined">help</span>
              </button>
            </div>
            <div className="w-8 h-8 rounded-full bg-surface-container-high overflow-hidden border border-outline-variant/15 flex-shrink-0 cursor-pointer hover:ring-2 hover:ring-primary/50 transition-all">
              <div className="w-full h-full bg-gradient-to-br from-primary-container to-inverse-primary flex items-center justify-center">
                <span className="text-on-primary-container text-xs font-bold">U</span>
              </div>
            </div>
          </div>
        </header>

        {/* 可滚动内容区域 */}
        <main className="flex-1 overflow-y-auto p-6 md:p-8 lg:p-12 relative">
          {/* 页面标题 */}
          <div className="max-w-6xl mx-auto mb-10">
            <h2 className="text-[2.25rem] font-bold tracking-[-0.02em] text-on-surface mb-2">
              系统设置
            </h2>
            <p className="text-on-surface-variant text-sm">
              管理全局环境变量、模型路由及代理的基础行为。
            </p>
          </div>

          <div className="max-w-6xl mx-auto space-y-12 pb-20">
            {/* 第一部分：通用设置 */}
            <section>
              <h3 className="text-lg font-semibold text-secondary mb-4 flex items-center gap-2">
                <span className="material-symbols-outlined text-[20px]">tune</span>
                通用设置
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {/* 主题设置 */}
                <div className="bg-surface-container-low rounded-xl p-5 flex flex-col justify-between relative overflow-hidden group">
                  <div className="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none"></div>
                  <div>
                    <span className="material-symbols-outlined text-outline mb-3">palette</span>
                    <h4 className="text-sm font-medium text-on-surface mb-1">界面主题</h4>
                    <p className="text-xs text-outline mb-4">确定全局界面的外观。</p>
                  </div>
                  <div className="flex bg-surface-container-highest p-1 rounded-lg">
                    <button
                      className={`flex-1 py-1.5 text-xs font-medium rounded-md transition-colors ${
                        theme === 'light'
                          ? 'bg-surface-variant text-primary shadow-sm ring-1 ring-white/5'
                          : 'text-on-surface-variant hover:text-on-surface'
                      }`}
                      onClick={() => handleThemeChange('light')}
                    >
                      浅色
                    </button>
                    <button
                      className={`flex-1 py-1.5 text-xs font-medium rounded-md transition-colors ${
                        theme === 'dark'
                          ? 'bg-surface-variant text-primary shadow-sm ring-1 ring-white/5'
                          : 'text-on-surface-variant hover:text-on-surface'
                      }`}
                      onClick={() => handleThemeChange('dark')}
                    >
                      深色
                    </button>
                    <button
                      className={`flex-1 py-1.5 text-xs font-medium rounded-md transition-colors ${
                        theme === 'auto'
                          ? 'bg-surface-variant text-primary shadow-sm ring-1 ring-white/5'
                          : 'text-on-surface-variant hover:text-on-surface'
                      }`}
                      onClick={() => handleThemeChange('auto')}
                    >
                      跟随系统
                    </button>
                  </div>
                </div>

                {/* 语言设置 */}
                <div className="bg-surface-container-low rounded-xl p-5 flex flex-col justify-between relative overflow-hidden group">
                  <div className="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none"></div>
                  <div>
                    <span className="material-symbols-outlined text-outline mb-3">language</span>
                    <h4 className="text-sm font-medium text-on-surface mb-1">语言设置</h4>
                    <p className="text-xs text-outline mb-4">设置偏好的界面语言。</p>
                  </div>
                  <div className="relative">
                    <select
                      className="w-full appearance-none bg-surface-container-highest border-none text-sm text-on-surface rounded-lg py-2 pl-3 pr-8 focus:ring-0 focus:border-b-2 focus:border-b-primary outline-none cursor-pointer"
                      value={language}
                      onChange={handleLanguageChange}
                    >
                      <option value="zh">简体中文</option>
                      <option value="en">English (US)</option>
                    </select>
                    <span className="material-symbols-outlined absolute right-2 top-1/2 -translate-y-1/2 text-outline pointer-events-none text-[20px]">
                      arrow_drop_down
                    </span>
                  </div>
                </div>

                {/* 日志级别设置 */}
                <div className="bg-surface-container-low rounded-xl p-5 flex flex-col justify-between relative overflow-hidden group">
                  <div className="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none"></div>
                  <div>
                    <span className="material-symbols-outlined text-outline mb-3">terminal</span>
                    <h4 className="text-sm font-medium text-on-surface mb-1">系统日志级别</h4>
                    <p className="text-xs text-outline mb-4">执行追踪的详细程度。</p>
                  </div>
                  <div className="relative">
                    <select
                      className="w-full appearance-none bg-surface-container-highest border-none text-sm text-on-surface rounded-lg py-2 pl-3 pr-8 focus:ring-0 focus:border-b-2 focus:border-b-primary outline-none font-mono cursor-pointer"
                      value={logLevel}
                      onChange={handleLogLevelChange}
                    >
                      <option value="debug">DEBUG</option>
                      <option value="info">INFO</option>
                      <option value="warn">WARN</option>
                      <option value="error">ERROR</option>
                    </select>
                    <span className="material-symbols-outlined absolute right-2 top-1/2 -translate-y-1/2 text-outline pointer-events-none text-[20px]">
                      arrow_drop_down
                    </span>
                  </div>
                </div>
              </div>
            </section>

            {/* 第二部分：模型提供商配置 */}
            <section>
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold text-secondary flex items-center gap-2">
                  <span className="material-symbols-outlined text-[20px]">dns</span>
                  模型提供商配置
                </h3>
                <button className="text-xs font-medium text-primary hover:text-primary-container transition-colors flex items-center gap-1">
                  <span className="material-symbols-outlined text-[16px]">add</span>
                  添加自定义提供商
                </button>
              </div>
              <div className="bg-surface-container-low rounded-xl flex flex-col overflow-hidden">
                {modelProviders.map((provider, index) => (
                  <React.Fragment key={provider.id}>
                    {index > 0 && <div className="h-[1px] bg-outline-variant/15 mx-4" />}
                    <div className="flex items-center justify-between p-4 hover:bg-white/[0.02] transition-colors relative group">
                      <div className="flex items-center gap-4 w-1/3">
                        <div className={`w-10 h-10 rounded-lg bg-surface-container-highest flex items-center justify-center shadow-inner ${provider.status === 'unconfigured' ? 'opacity-60' : ''}`}>
                          <span className="material-symbols-outlined text-on-surface">
                            {provider.icon}
                          </span>
                        </div>
                        <div className={provider.status === 'unconfigured' ? 'opacity-60' : ''}>
                          <h4 className="text-sm font-semibold text-on-surface">{provider.name}</h4>
                          <p className="text-xs font-mono text-outline mt-0.5">{provider.endpoint}</p>
                        </div>
                      </div>
                      <div className="flex items-center gap-3 w-1/4">
                        <span
                          className="flex h-2 w-2 rounded-full shadow-[0_0_8px_rgba(74,222,128,0.4)]"
                          style={{ backgroundColor: getStatusColor(provider.status) }}
                        ></span>
                        <span className={`text-xs font-medium uppercase tracking-wider ${
                          provider.status === 'unconfigured' ? 'text-error' : 'text-on-surface'
                        }`}>
                          {getStatusText(provider.status)}
                        </span>
                      </div>
                      <div className="w-1/4">
                        <select
                          className={`w-full bg-surface-container-highest border-none text-xs text-on-surface font-mono rounded-md py-1.5 px-3 outline-none focus:ring-1 focus:ring-primary/50 appearance-none ${
                            provider.status === 'unconfigured' ? 'opacity-50 pointer-events-none' : ''
                          }`}
                          value={provider.defaultModel}
                          onChange={(e) => handleModelChange(provider.id, e.target.value)}
                          disabled={provider.status === 'unconfigured'}
                        >
                          {provider.availableModels.map(model => (
                            <option key={model} value={model}>{model}</option>
                          ))}
                        </select>
                      </div>
                      <div className="flex justify-end w-[15%]">
                        {provider.status === 'unconfigured' ? (
                          <button className="text-xs font-medium bg-surface-container-highest text-on-surface hover:bg-surface-bright px-3 py-1.5 rounded-md transition-all shadow-sm">
                            配置密钥
                          </button>
                        ) : (
                          <button className="text-xs font-medium text-primary hover:bg-surface-container-highest px-3 py-1.5 rounded-md transition-all">
                            测试连接
                          </button>
                        )}
                      </div>
                    </div>
                  </React.Fragment>
                ))}
              </div>
            </section>

            {/* 第三部分：Agent 角色配置 */}
            <section>
              <h3 className="text-lg font-semibold text-secondary mb-4 flex items-center gap-2">
                <span className="material-symbols-outlined text-[20px]">group_work</span>
                全局 Agent 角色模板
              </h3>
              <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
                {/* Orchestrator */}
                <div className="bg-surface-container-high rounded-xl p-6 relative overflow-hidden group shadow-[0_4px_24px_rgba(0,0,0,0.2)] border-t border-l border-white/5">
                  <div className="absolute -top-10 -right-10 w-32 h-32 bg-primary/5 rounded-full blur-2xl group-hover:bg-primary/10 transition-colors duration-500"></div>
                  <div className="flex justify-between items-start mb-6 relative z-10">
                    <div className="flex items-center gap-3">
                      <div className="p-2 rounded-lg bg-surface-container-lowest text-primary">
                        <span className="material-symbols-outlined">account_tree</span>
                      </div>
                      <div>
                        <h4 className="text-base font-semibold text-on-surface">Orchestrator</h4>
                        <p className="text-xs text-outline">任务路由与分解。</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2 bg-surface-container-lowest px-2 py-1 rounded-md">
                      <span className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse"></span>
                      <span className="text-[10px] font-mono text-primary uppercase">活跃</span>
                    </div>
                  </div>
                  <div className="space-y-4 relative z-10">
                    <div>
                      <label className="text-[10px] uppercase font-semibold text-outline tracking-wider block mb-1.5">
                        分配模型
                      </label>
                      <div className="bg-surface-container-lowest rounded-md px-3 py-2 flex items-center justify-between">
                        <span className="font-mono text-xs text-secondary">
                          {agentRoles.find(r => r.id === 'orchestrator')?.assignedModel}
                        </span>
                        <span className="material-symbols-outlined text-[16px] text-outline cursor-pointer hover:text-on-surface">
                          swap_horiz
                        </span>
                      </div>
                    </div>
                    <div>
                      <label className="text-[10px] uppercase font-semibold text-outline tracking-wider block mb-1.5 flex justify-between">
                        <span>系统提示词</span>
                        <span className="material-symbols-outlined text-[14px] cursor-pointer hover:text-primary transition-colors">
                          edit
                        </span>
                      </label>
                      <div className="bg-surface-container-lowest rounded-md p-3 h-20 overflow-hidden relative">
                        <p className="font-mono text-[11px] text-on-surface-variant leading-relaxed">
                          {agentRoles.find(r => r.id === 'orchestrator')?.systemPrompt}
                        </p>
                        <div className="absolute bottom-0 left-0 right-0 h-8 bg-gradient-to-t from-surface-container-lowest to-transparent pointer-events-none"></div>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Architect */}
                <div className="bg-surface-container-high rounded-xl p-6 relative overflow-hidden group shadow-[0_4px_24px_rgba(0,0,0,0.2)] border-t border-l border-white/5">
                  <div className="absolute -top-10 -right-10 w-32 h-32 bg-primary/5 rounded-full blur-2xl group-hover:bg-primary/10 transition-colors duration-500"></div>
                  <div className="flex justify-between items-start mb-6 relative z-10">
                    <div className="flex items-center gap-3">
                      <div className="p-2 rounded-lg bg-surface-container-lowest text-primary">
                        <span className="material-symbols-outlined">architecture</span>
                      </div>
                      <div>
                        <h4 className="text-base font-semibold text-on-surface">Architect</h4>
                        <p className="text-xs text-outline">系统设计与架构。</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2 bg-surface-container-lowest px-2 py-1 rounded-md">
                      <span className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse"></span>
                      <span className="text-[10px] font-mono text-primary uppercase">活跃</span>
                    </div>
                  </div>
                  <div className="space-y-4 relative z-10">
                    <div>
                      <label className="text-[10px] uppercase font-semibold text-outline tracking-wider block mb-1.5">
                        分配模型
                      </label>
                      <div className="bg-surface-container-lowest rounded-md px-3 py-2 flex items-center justify-between">
                        <span className="font-mono text-xs text-secondary">
                          {agentRoles.find(r => r.id === 'architect')?.assignedModel}
                        </span>
                        <span className="material-symbols-outlined text-[16px] text-outline cursor-pointer hover:text-on-surface">
                          swap_horiz
                        </span>
                      </div>
                    </div>
                    <div>
                      <label className="text-[10px] uppercase font-semibold text-outline tracking-wider block mb-1.5 flex justify-between">
                        <span>系统提示词</span>
                        <span className="material-symbols-outlined text-[14px] cursor-pointer hover:text-primary transition-colors">
                          edit
                        </span>
                      </label>
                      <div className="bg-surface-container-lowest rounded-md p-3 h-20 overflow-hidden relative">
                        <p className="font-mono text-[11px] text-on-surface-variant leading-relaxed">
                          {agentRoles.find(r => r.id === 'architect')?.systemPrompt}
                        </p>
                        <div className="absolute bottom-0 left-0 right-0 h-8 bg-gradient-to-t from-surface-container-lowest to-transparent pointer-events-none"></div>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Frontend Dev - 跨两列 */}
                <div className="bg-surface-container-high rounded-xl p-6 relative overflow-hidden group shadow-[0_4px_24px_rgba(0,0,0,0.2)] border-t border-l border-white/5 xl:col-span-2">
                  <div className="flex flex-col md:flex-row gap-6 relative z-10">
                    {/* 左侧：身份和模型 */}
                    <div className="w-full md:w-1/3 flex flex-col justify-between">
                      <div className="flex justify-between items-start mb-6">
                        <div className="flex items-center gap-3">
                          <div className="p-2 rounded-lg bg-surface-container-lowest text-primary">
                            <span className="material-symbols-outlined">web</span>
                          </div>
                          <div>
                            <h4 className="text-base font-semibold text-on-surface">Frontend Dev</h4>
                            <p className="text-xs text-outline">UI 生成与逻辑。</p>
                          </div>
                        </div>
                      </div>
                      <div className="mt-auto">
                        <label className="text-[10px] uppercase font-semibold text-outline tracking-wider block mb-1.5">
                          分配模型
                        </label>
                        <div className="bg-surface-container-lowest rounded-md px-3 py-2 flex items-center justify-between">
                          <span className="font-mono text-xs text-secondary">
                            {agentRoles.find(r => r.id === 'frontend')?.assignedModel}
                          </span>
                          <span className="material-symbols-outlined text-[16px] text-outline cursor-pointer hover:text-on-surface">
                            swap_horiz
                          </span>
                        </div>
                      </div>
                    </div>
                    {/* 右侧：提示词编辑器 */}
                    <div className="w-full md:w-2/3">
                      <label className="text-[10px] uppercase font-semibold text-outline tracking-wider block mb-1.5 flex justify-between">
                        <span>系统提示词</span>
                        <div className="flex gap-2">
                          <span className="material-symbols-outlined text-[14px] cursor-pointer text-outline hover:text-on-surface transition-colors">
                            history
                          </span>
                          <span className="material-symbols-outlined text-[14px] cursor-pointer text-outline hover:text-primary transition-colors">
                            open_in_full
                          </span>
                        </div>
                      </label>
                      <textarea
                        className="w-full bg-surface-container-lowest border-none rounded-md p-3 font-mono text-[11px] text-on-surface-variant leading-relaxed h-32 resize-none focus:ring-1 focus:ring-primary/30 outline-none"
                        spellCheck="false"
                        value={agentRoles.find(r => r.id === 'frontend')?.systemPrompt}
                        onChange={(e) => {
                          const newRoles = agentRoles.map(r => 
                            r.id === 'frontend' ? { ...r, systemPrompt: e.target.value } : r
                          );
                          setAgentRoles(newRoles);
                          setHasChanges(true);
                        }}
                      />
                    </div>
                  </div>
                </div>
              </div>
            </section>
          </div>

          {/* 浮动保存操作 */}
          {hasChanges && (
            <div className="fixed bottom-8 right-8 z-50">
              <div className="bg-surface-variant/80 backdrop-blur-xl p-1.5 rounded-full shadow-[0_8px_32px_rgba(0,0,0,0.5)] border border-outline-variant/20 flex items-center gap-2 translate-y-2 opacity-0 hover:translate-y-0 hover:opacity-100 transition-all duration-300" style={{ animation: 'slideUpFade 0.5s ease forwards 1s' }}>
                <button
                  className="bg-surface-container-lowest text-on-surface hover:bg-surface-bright rounded-full px-4 py-2 text-sm font-medium transition-colors"
                  onClick={handleDiscard}
                >
                  放弃更改
                </button>
                <button
                  className="bg-primary-container text-on-primary-container hover:bg-inverse-primary rounded-full px-6 py-2 text-sm font-medium shadow-lg shadow-primary/20 transition-all flex items-center gap-2"
                  onClick={handleSave}
                >
                  <span className="material-symbols-outlined text-[18px]">save</span>
                  保存配置
                </button>
              </div>
            </div>
          )}
        </main>
      </div>
    </div>
  );
};

export default SettingsViewNew;