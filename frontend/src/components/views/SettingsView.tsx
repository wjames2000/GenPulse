import React, { useState, useEffect } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '../ui/Card';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Select } from '../ui/Select';
import { TextArea } from '../ui/TextArea';
import { useAppStore } from '../../stores/appStore';

interface ModelProvider {
  id: string;
  name: string;
  icon: string;
  endpoint: string;
  status: 'connected' | 'disconnected' | 'unconfigured';
  defaultModel: string;
  availableModels: string[];
  apiKey: string;
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

interface SettingsData {
  general: {
    theme: 'light' | 'dark' | 'auto';
    language: string;
    logLevel: 'debug' | 'info' | 'warn' | 'error';
  };
  modelProviders: ModelProvider[];
  agentRoles: AgentRole[];
  advanced: {
    enableExperimental: boolean;
    enableTelemetry: boolean;
    developerMode: boolean;
  };
}

const SettingsView: React.FC = () => {
  const { darkMode, toggleDarkMode } = useAppStore();
  const [settings, setSettings] = useState<SettingsData>({
    general: {
      theme: 'dark',
      language: 'zh',
      logLevel: 'info',
    },
    modelProviders: [
      {
        id: 'openai',
        name: 'OpenAI',
        icon: 'data_object',
        endpoint: 'api.openai.com/v1',
        status: 'connected',
        defaultModel: 'gpt-4-turbo-preview',
        availableModels: ['gpt-4-turbo-preview', 'gpt-4', 'gpt-3.5-turbo'],
        apiKey: '',
      },
      {
        id: 'anthropic',
        name: 'Anthropic Claude',
        icon: 'psychology_alt',
        endpoint: 'api.anthropic.com',
        status: 'connected',
        defaultModel: 'claude-3-opus-20240229',
        availableModels: ['claude-3-opus-20240229', 'claude-3-sonnet-20240229'],
        apiKey: '',
      },
      {
        id: 'gemini',
        name: 'Google Gemini',
        icon: 'auto_awesome',
        endpoint: 'generativelanguage.googleapis.com',
        status: 'unconfigured',
        defaultModel: 'gemini-1.5-pro',
        availableModels: ['gemini-1.5-pro'],
        apiKey: '',
      },
    ],
    agentRoles: [
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
    ],
    advanced: {
      enableExperimental: false,
      enableTelemetry: true,
      developerMode: false,
    },
  });

  const [isSaving, setIsSaving] = useState(false);
  const [hasChanges, setHasChanges] = useState(false);
  const [activeSection, setActiveSection] = useState<'general' | 'models' | 'agents' | 'advanced'>('general');

  // 模拟加载设置
  useEffect(() => {
    const loadSettings = async () => {
      const savedSettings = localStorage.getItem('genpulse-settings-v2');
      if (savedSettings) {
        try {
          setSettings(JSON.parse(savedSettings));
        } catch (error) {
          console.error('Failed to load settings:', error);
        }
      }
    };
    loadSettings();
  }, []);

  // 检测设置变化
  useEffect(() => {
    const savedSettings = localStorage.getItem('genpulse-settings-v2');
    const currentSettings = JSON.stringify(settings);
    setHasChanges(savedSettings !== currentSettings);
  }, [settings]);

  const handleGeneralChange = (key: keyof SettingsData['general'], value: any) => {
    setSettings(prev => ({
      ...prev,
      general: {
        ...prev.general,
        [key]: value,
      },
    }));
  };

  const handleModelProviderChange = (providerId: string, key: keyof ModelProvider, value: any) => {
    setSettings(prev => ({
      ...prev,
      modelProviders: prev.modelProviders.map(provider =>
        provider.id === providerId ? { ...provider, [key]: value } : provider
      ),
    }));
  };

  const handleAgentRoleChange = (roleId: string, key: keyof AgentRole, value: any) => {
    setSettings(prev => ({
      ...prev,
      agentRoles: prev.agentRoles.map(role =>
        role.id === roleId ? { ...role, [key]: value } : role
      ),
    }));
  };

  const handleAdvancedChange = (key: keyof SettingsData['advanced'], value: any) => {
    setSettings(prev => ({
      ...prev,
      advanced: {
        ...prev.advanced,
        [key]: value,
      },
    }));
  };

  const handleSave = async () => {
    setIsSaving(true);
    try {
      localStorage.setItem('genpulse-settings-v2', JSON.stringify(settings));
      
      // 应用主题设置
      if (settings.general.theme === 'dark' && !darkMode) {
        toggleDarkMode();
      } else if (settings.general.theme === 'light' && darkMode) {
        toggleDarkMode();
      }
      
      await new Promise(resolve => setTimeout(resolve, 500));
      
      setHasChanges(false);
      // 显示保存成功提示
      console.log('Settings saved successfully!');
    } catch (error) {
      console.error('Failed to save settings:', error);
    } finally {
      setIsSaving(false);
    }
  };

  const handleReset = () => {
    if (confirm('确定要重置所有设置为默认值吗？')) {
      setSettings({
        general: {
          theme: 'dark',
          language: 'zh',
          logLevel: 'info',
        },
        modelProviders: [
          {
            id: 'openai',
            name: 'OpenAI',
            icon: 'data_object',
            endpoint: 'api.openai.com/v1',
            status: 'connected',
            defaultModel: 'gpt-4-turbo-preview',
            availableModels: ['gpt-4-turbo-preview', 'gpt-4', 'gpt-3.5-turbo'],
            apiKey: '',
          },
          {
            id: 'anthropic',
            name: 'Anthropic Claude',
            icon: 'psychology_alt',
            endpoint: 'api.anthropic.com',
            status: 'connected',
            defaultModel: 'claude-3-opus-20240229',
            availableModels: ['claude-3-opus-20240229', 'claude-3-sonnet-20240229'],
            apiKey: '',
          },
          {
            id: 'gemini',
            name: 'Google Gemini',
            icon: 'auto_awesome',
            endpoint: 'generativelanguage.googleapis.com',
            status: 'unconfigured',
            defaultModel: 'gemini-1.5-pro',
            availableModels: ['gemini-1.5-pro'],
            apiKey: '',
          },
        ],
        agentRoles: [
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
        ],
        advanced: {
          enableExperimental: false,
          enableTelemetry: true,
          developerMode: false,
        },
      });
    }
  };

  const testConnection = (providerId: string) => {
    console.log(`Testing connection for ${providerId}`);
    // TODO: 实现连接测试
  };

  const configureApiKey = (providerId: string) => {
    const apiKey = prompt(`请输入 ${providerId} 的 API 密钥:`);
    if (apiKey) {
      handleModelProviderChange(providerId, 'apiKey', apiKey);
      handleModelProviderChange(providerId, 'status', 'connected');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'connected':
        return 'bg-green-500';
      case 'disconnected':
        return 'bg-yellow-500';
      case 'unconfigured':
        return 'bg-red-500';
      default:
        return 'bg-gray-500';
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
      default:
        return '未知';
    }
  };

  return (
    <div className="settings-view">
      <div className="settings-header mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">系统设置</h1>
        <p className="text-gray-600 dark:text-gray-400 text-sm">
          管理全局环境变量、模型路由及代理的基础行为
        </p>
      </div>

      {/* 导航标签 */}
      <div className="settings-tabs mb-6 border-b border-gray-200 dark:border-gray-700">
        <nav className="flex space-x-8">
          <button
            className={`pb-3 px-1 text-sm font-medium border-b-2 transition-colors ${
              activeSection === 'general'
                ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
            }`}
            onClick={() => setActiveSection('general')}
          >
            <span className="flex items-center gap-2">
              <span className="material-symbols-outlined text-base">tune</span>
              通用设置
            </span>
          </button>
          <button
            className={`pb-3 px-1 text-sm font-medium border-b-2 transition-colors ${
              activeSection === 'models'
                ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
            }`}
            onClick={() => setActiveSection('models')}
          >
            <span className="flex items-center gap-2">
              <span className="material-symbols-outlined text-base">dns</span>
              模型提供商
            </span>
          </button>
          <button
            className={`pb-3 px-1 text-sm font-medium border-b-2 transition-colors ${
              activeSection === 'agents'
                ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
            }`}
            onClick={() => setActiveSection('agents')}
          >
            <span className="flex items-center gap-2">
              <span className="material-symbols-outlined text-base">group_work</span>
              Agent 角色
            </span>
          </button>
          <button
            className={`pb-3 px-1 text-sm font-medium border-b-2 transition-colors ${
              activeSection === 'advanced'
                ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
            }`}
            onClick={() => setActiveSection('advanced')}
          >
            <span className="flex items-center gap-2">
              <span className="material-symbols-outlined text-base">code</span>
              高级设置
            </span>
          </button>
        </nav>
      </div>

      {/* 通用设置 */}
      {activeSection === 'general' && (
        <div className="space-y-8">
          <section>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
              <span className="material-symbols-outlined text-xl">tune</span>
              通用设置
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {/* 主题设置 */}
              <Card className="bg-gray-50 dark:bg-gray-800 rounded-xl p-5">
                <div className="mb-4">
                  <span className="material-symbols-outlined text-gray-500 dark:text-gray-400 mb-3 block">palette</span>
                  <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-1">界面主题</h4>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mb-4">确定全局界面的外观</p>
                </div>
                <div className="flex bg-gray-100 dark:bg-gray-700 p-1 rounded-lg">
                  <button
                    className={`flex-1 py-1.5 text-xs font-medium rounded-md transition-colors ${
                      settings.general.theme === 'light'
                        ? 'bg-white dark:bg-gray-600 text-gray-900 dark:text-white shadow-sm'
                        : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
                    }`}
                    onClick={() => handleGeneralChange('theme', 'light')}
                  >
                    浅色
                  </button>
                  <button
                    className={`flex-1 py-1.5 text-xs font-medium rounded-md transition-colors ${
                      settings.general.theme === 'dark'
                        ? 'bg-gray-800 dark:bg-gray-600 text-white shadow-sm'
                        : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
                    }`}
                    onClick={() => handleGeneralChange('theme', 'dark')}
                  >
                    深色
                  </button>
                  <button
                    className={`flex-1 py-1.5 text-xs font-medium rounded-md transition-colors ${
                      settings.general.theme === 'auto'
                        ? 'bg-gray-100 dark:bg-gray-600 text-gray-900 dark:text-white shadow-sm'
                        : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
                    }`}
                    onClick={() => handleGeneralChange('theme', 'auto')}
                  >
                    跟随系统
                  </button>
                </div>
              </Card>

              {/* 语言设置 */}
              <Card className="bg-gray-50 dark:bg-gray-800 rounded-xl p-5">
                <div className="mb-4">
                  <span className="material-symbols-outlined text-gray-500 dark:text-gray-400 mb-3 block">language</span>
                  <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-1">语言设置</h4>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mb-4">设置偏好的界面语言</p>
                </div>
                <div className="relative">
                  <Select
                    value={settings.general.language}
                    onChange={(e) => handleGeneralChange('language', e.target.value)}
                    className="w-full bg-gray-100 dark:bg-gray-700 border-none text-sm text-gray-900 dark:text-white rounded-lg py-2 pl-3 pr-8"
                    options={[
                      { value: 'zh', label: '简体中文' },
                      { value: 'en', label: 'English (US)' },
                    ]}
                  />
                </div>
              </Card>

              {/* 日志级别设置 */}
              <Card className="bg-gray-50 dark:bg-gray-800 rounded-xl p-5">
                <div className="mb-4">
                  <span className="material-symbols-outlined text-gray-500 dark:text-gray-400 mb-3 block">terminal</span>
                  <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-1">系统日志级别</h4>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mb-4">执行追踪的详细程度</p>
                </div>
                <div className="relative">
                  <Select
                    value={settings.general.logLevel}
                    onChange={(e) => handleGeneralChange('logLevel', e.target.value as any)}
                    className="w-full bg-gray-100 dark:bg-gray-700 border-none text-sm text-gray-900 dark:text-white rounded-lg py-2 pl-3 pr-8 font-mono"
                    options={[
                      { value: 'debug', label: 'DEBUG' },
                      { value: 'info', label: 'INFO' },
                      { value: 'warn', label: 'WARN' },
                      { value: 'error', label: 'ERROR' },
                    ]}
                  />
                </div>
              </Card>
            </div>
          </section>
        </div>
      )}

      {/* 模型提供商设置 */}
      {activeSection === 'models' && (
        <div className="space-y-8">
          <section>
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
                <span className="material-symbols-outlined text-xl">dns</span>
                模型提供商配置
              </h3>
              <button className="text-xs font-medium text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 transition-colors flex items-center gap-1">
                <span className="material-symbols-outlined text-base">add</span>
                添加自定义提供商
              </button>
            </div>
            <Card className="bg-gray-50 dark:bg-gray-800 rounded-xl">
              {settings.modelProviders.map((provider, index) => (
                <React.Fragment key={provider.id}>
                  {index > 0 && <div className="h-px bg-gray-200 dark:bg-gray-700 mx-4" />}
                  <div className="flex items-center justify-between p-4 hover:bg-gray-100 dark:hover:bg-gray-700/50 transition-colors">
                    <div className="flex items-center gap-4 w-1/3">
                      <div className="w-10 h-10 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                        <span className="material-symbols-outlined text-gray-700 dark:text-gray-300">
                          {provider.icon}
                        </span>
                      </div>
                      <div>
                        <h4 className="text-sm font-semibold text-gray-900 dark:text-white">{provider.name}</h4>
                        <p className="text-xs font-mono text-gray-500 dark:text-gray-400 mt-0.5">{provider.endpoint}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-3 w-1/4">
                      <span className={`flex h-2 w-2 rounded-full ${getStatusColor(provider.status)}`} />
                      <span className="text-xs font-medium text-gray-900 dark:text-white uppercase tracking-wider">
                        {getStatusText(provider.status)}
                      </span>
                    </div>
                    <div className="w-1/4">
                      <Select
                        value={provider.defaultModel}
                        onChange={(e) => handleModelProviderChange(provider.id, 'defaultModel', e.target.value)}
                        className="w-full bg-gray-100 dark:bg-gray-700 border-none text-xs text-gray-900 dark:text-white font-mono rounded-md py-1.5 px-3"
                        disabled={provider.status === 'unconfigured'}
                        options={provider.availableModels.map(model => ({
                          value: model,
                          label: model,
                        }))}
                      />
                    </div>
                    <div className="flex justify-end w-[15%]">
                      {provider.status === 'unconfigured' ? (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => configureApiKey(provider.id)}
                          className="text-xs"
                        >
                          配置密钥
                        </Button>
                      ) : (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => testConnection(provider.id)}
                          className="text-xs"
                        >
                          测试连接
                        </Button>
                      )}
                    </div>
                  </div>
                </React.Fragment>
              ))}
            </Card>
          </section>
        </div>
      )}

      {/* Agent 角色设置 */}
      {activeSection === 'agents' && (
        <div className="space-y-8">
          <section>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
              <span className="material-symbols-outlined text-xl">group_work</span>
              全局 Agent 角色模板
            </h3>
            <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
              {settings.agentRoles.map((role) => (
                <Card key={role.id} className="bg-gray-50 dark:bg-gray-800 rounded-xl p-6">
                  <div className="flex justify-between items-start mb-6">
                    <div className="flex items-center gap-3">
                      <div className="p-2 rounded-lg bg-gray-100 dark:bg-gray-700 text-blue-600 dark:text-blue-400">
                        <span className="material-symbols-outlined">{role.icon}</span>
                      </div>
                      <div>
                        <h4 className="text-base font-semibold text-gray-900 dark:text-white">{role.name}</h4>
                        <p className="text-xs text-gray-500 dark:text-gray-400">{role.description}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2 bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded-md">
                      <span className="h-1.5 w-1.5 rounded-full bg-green-500 animate-pulse" />
                      <span className="text-[10px] font-mono text-green-600 dark:text-green-400 uppercase">
                        {role.status === 'active' ? '活跃' : '未激活'}
                      </span>
                    </div>
                  </div>
                  <div className="space-y-4">
                    <div>
                      <label className="text-[10px] uppercase font-semibold text-gray-500 dark:text-gray-400 tracking-wider block mb-1.5">
                        分配模型
                      </label>
                      <div className="bg-gray-100 dark:bg-gray-700 rounded-md px-3 py-2 flex items-center justify-between">
                        <span className="font-mono text-xs text-blue-600 dark:text-blue-400">
                          {role.assignedModel}
                        </span>
                        <span className="material-symbols-outlined text-base text-gray-500 dark:text-gray-400 cursor-pointer hover:text-gray-700 dark:hover:text-gray-300">
                          swap_horiz
                        </span>
                      </div>
                    </div>
                    <div>
                      <label className="text-[10px] uppercase font-semibold text-gray-500 dark:text-gray-400 tracking-wider block mb-1.5 flex justify-between">
                        <span>系统提示词</span>
                        <span className="material-symbols-outlined text-sm cursor-pointer hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
                          edit
                        </span>
                      </label>
                      <div className="bg-gray-100 dark:bg-gray-700 rounded-md p-3 h-20 overflow-hidden relative">
                        <p className="font-mono text-[11px] text-gray-600 dark:text-gray-400 leading-relaxed">
                          {role.systemPrompt}
                        </p>
                        <div className="absolute bottom-0 left-0 right-0 h-8 bg-gradient-to-t from-gray-100 dark:from-gray-700 to-transparent pointer-events-none" />
                      </div>
                    </div>
                  </div>
                </Card>
              ))}
            </div>
          </section>
        </div>
      )}

      {/* 高级设置 */}
      {activeSection === 'advanced' && (
        <div className="space-y-8">
          <section>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
              <span className="material-symbols-outlined text-xl">code</span>
              高级设置
            </h3>
            <Card className="bg-gray-50 dark:bg-gray-800 rounded-xl p-6">
              <div className="space-y-6">
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-1">实验性功能</h4>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                      启用开发中的新功能，可能不稳定
                    </p>
                  </div>
                  <label className="relative inline-flex items-center cursor-pointer">
                    <input
                      type="checkbox"
                      className="sr-only peer"
                      checked={settings.advanced.enableExperimental}
                      onChange={(e) => handleAdvancedChange('enableExperimental', e.target.checked)}
                    />
                    <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600" />
                  </label>
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-1">匿名遥测数据</h4>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                      发送匿名使用数据以帮助改进产品
                    </p>
                  </div>
                  <label className="relative inline-flex items-center cursor-pointer">
                    <input
                      type="checkbox"
                      className="sr-only peer"
                      checked={settings.advanced.enableTelemetry}
                      onChange={(e) => handleAdvancedChange('enableTelemetry', e.target.checked)}
                    />
                    <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600" />
                  </label>
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-1">开发者模式</h4>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                      显示开发者工具和高级调试选项
                    </p>
                  </div>
                  <label className="relative inline-flex items-center cursor-pointer">
                    <input
                      type="checkbox"
                      className="sr-only peer"
                      checked={settings.advanced.developerMode}
                      onChange={(e) => handleAdvancedChange('developerMode', e.target.checked)}
                    />
                    <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600" />
                  </label>
                </div>
              </div>
            </Card>
          </section>
        </div>
      )}

      {/* 浮动保存按钮 */}
      {hasChanges && (
        <div className="fixed bottom-8 right-8 z-50 animate-slideUpFade">
          <div className="bg-white dark:bg-gray-800/80 backdrop-blur-xl p-1.5 rounded-full shadow-lg border border-gray-200 dark:border-gray-700 flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handleReset}
              disabled={isSaving}
              className="rounded-full"
            >
              放弃更改
            </Button>
            <Button
              variant="primary"
              size="sm"
              onClick={handleSave}
              disabled={isSaving}
              className="rounded-full flex items-center gap-2"
            >
              <span className="material-symbols-outlined text-base">save</span>
              {isSaving ? '保存中...' : '保存配置'}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};

export default SettingsView;