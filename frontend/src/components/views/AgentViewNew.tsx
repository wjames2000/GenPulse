import React, { useState, useEffect } from 'react';
import { useAgentStore } from '../../stores/agentStore';
import { useNotificationStore } from '../../stores/notificationStore';
import { colors } from '../../design-system/colors';

interface AgentConfig {
  id: string;
  name: string;
  role: string;
  description: string;
  model_config: {
    type: string;
    name: string;
    provider: string;
  };
  capabilities: string[];
  tools: string[];
  enabled: boolean;
}

interface AgentExecution {
  id: string;
  agent_id: string;
  task: string;
  state: string;
  started_at: string;
  completed_at?: string;
  error?: string;
  result?: any;
}

const AgentViewNew: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'agents' | 'execution' | 'history'>('agents');
  const [selectedAgent, setSelectedAgent] = useState<string>('');
  const [taskInput, setTaskInput] = useState<string>('');
  const [projectPath, setProjectPath] = useState<string>('./test-project');
  const [projectType, setProjectType] = useState<string>('go');
  const [projectName, setProjectName] = useState<string>('TestProject');
  
  const { 
    agents, 
    executions, 
    loading, 
    error,
    fetchAgents,
    fetchAgentStatus,
    executeAgent,
    fetchExecutions
  } = useAgentStore();
  
  const { addNotification } = useNotificationStore();

  useEffect(() => {
    fetchAgents();
    fetchExecutions();
    
    const interval = setInterval(() => {
      if (selectedAgent) {
        fetchAgentStatus(selectedAgent);
      }
      fetchExecutions();
    }, 10000);
    
    return () => clearInterval(interval);
  }, [selectedAgent]);

  const handleExecuteTask = async () => {
    if (!selectedAgent || !taskInput.trim()) {
      addNotification('请选择Agent并输入任务', 'warning');
      return;
    }

    try {
      const parameters: Record<string, any> = {
        project_path: projectPath,
        project_type: projectType,
        project_name: projectName,
      };

      if (taskInput.toLowerCase().includes('创建项目') || taskInput.toLowerCase().includes('create project')) {
        parameters.project_path = projectPath;
        parameters.project_type = projectType;
        parameters.project_name = projectName;
      }

      await executeAgent(selectedAgent, taskInput, parameters);
      addNotification('任务已开始执行', 'success');
      setTaskInput('');
      
      setTimeout(() => {
        fetchExecutions();
      }, 1000);
    } catch (err) {
      addNotification(`执行失败: ${err instanceof Error ? err.message : '未知错误'}`, 'error');
    }
  };

  const handleCancelExecution = async (executionId: string) => {
    try {
      addNotification('取消执行功能待实现', 'info');
    } catch (err) {
      addNotification(`取消失败: ${err instanceof Error ? err.message : '未知错误'}`, 'error');
    }
  };

  const getAgentStatusColor = (state: string) => {
    switch (state) {
      case 'idle': return 'bg-green-500/10 text-green-400';
      case 'thinking': return 'bg-blue-500/10 text-blue-400';
      case 'executing': return 'bg-yellow-500/10 text-yellow-400';
      case 'waiting': return 'bg-purple-500/10 text-purple-400';
      case 'completed': return 'bg-green-500/10 text-green-400';
      case 'failed': return 'bg-red-500/10 text-red-400';
      default: return 'bg-surface-container text-on-surface-variant';
    }
  };

  const getExecutionStatusColor = (state: string) => {
    switch (state) {
      case 'completed': return 'bg-green-500/10 text-green-400';
      case 'failed': return 'bg-red-500/10 text-red-400';
      case 'thinking':
      case 'executing':
      case 'waiting': return 'bg-yellow-500/10 text-yellow-400';
      default: return 'bg-surface-container text-on-surface-variant';
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN');
  };

  const getSelectedAgent = () => {
    return agents.find(agent => agent.id === selectedAgent);
  };

  return (
    <div className="p-6 space-y-6">
      {/* 页面标题 */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-on-surface">代理管理</h1>
          <p className="text-sm text-on-surface-variant mt-1">管理您的AI代理，执行任务并监控状态</p>
        </div>
      </div>

      {/* 标签页导航 */}
      <div className="flex space-x-2 border-b border-outline-variant pb-2">
        <button
          className={`px-4 py-2 rounded-t-lg transition-all duration-200 ${
            activeTab === 'agents' 
              ? 'bg-surface-container-high text-primary border-b-2 border-primary' 
              : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-container'
          }`}
          onClick={() => setActiveTab('agents')}
        >
          <span className="flex items-center">
            <span className="material-symbols-rounded text-sm mr-2">person</span>
            Agent列表
          </span>
        </button>
        <button
          className={`px-4 py-2 rounded-t-lg transition-all duration-200 ${
            activeTab === 'execution' 
              ? 'bg-surface-container-high text-primary border-b-2 border-primary' 
              : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-container'
          }`}
          onClick={() => setActiveTab('execution')}
        >
          <span className="flex items-center">
            <span className="material-symbols-rounded text-sm mr-2">play_arrow</span>
            任务执行
          </span>
        </button>
        <button
          className={`px-4 py-2 rounded-t-lg transition-all duration-200 ${
            activeTab === 'history' 
              ? 'bg-surface-container-high text-primary border-b-2 border-primary' 
              : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-container'
          }`}
          onClick={() => setActiveTab('history')}
        >
          <span className="flex items-center">
            <span className="material-symbols-rounded text-sm mr-2">history</span>
            执行历史
          </span>
        </button>
      </div>

      {/* Agent列表 */}
      {activeTab === 'agents' && (
        <div className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {agents.map(agent => (
              <div 
                key={agent.id} 
                className={`bg-surface-container rounded-xl p-4 border transition-all duration-200 cursor-pointer hover:border-primary/30 ${
                  selectedAgent === agent.id 
                    ? 'border-primary ring-1 ring-primary/20' 
                    : 'border-outline-variant'
                }`}
                onClick={() => setSelectedAgent(agent.id)}
              >
                <div className="flex justify-between items-start mb-3">
                  <div>
                    <h3 className="font-medium text-on-surface">{agent.name}</h3>
                    <div className={`inline-flex items-center px-2 py-1 rounded-full text-xs mt-1 ${
                      agent.enabled ? 'bg-green-500/10 text-green-400' : 'bg-surface-container-high text-on-surface-variant'
                    }`}>
                      <span className="w-1.5 h-1.5 rounded-full mr-1.5 bg-current"></span>
                      {agent.enabled ? '启用' : '禁用'}
                    </div>
                  </div>
                  <div className="text-xs px-2 py-1 rounded bg-surface-container-high text-on-surface-variant">
                    {agent.role}
                  </div>
                </div>
                
                <p className="text-sm text-on-surface-variant mb-3 line-clamp-2">{agent.description}</p>
                
                <div className="space-y-2">
                  <div className="flex items-center text-sm">
                    <span className="material-symbols-rounded text-base mr-2 text-primary">smart_toy</span>
                    <span className="text-on-surface-variant">{agent.model_config.name}</span>
                    <span className="text-xs text-on-surface-variant/70 ml-2">({agent.model_config.provider})</span>
                  </div>
                  
                  <div className="flex items-center text-sm">
                    <span className="material-symbols-rounded text-base mr-2 text-secondary">build</span>
                    <span className="text-on-surface-variant">{agent.tools.length}个工具</span>
                  </div>
                  
                  <div className="flex flex-wrap gap-1 mt-2">
                    {agent.capabilities.slice(0, 3).map((cap, index) => (
                      <span key={index} className="text-xs px-2 py-1 rounded bg-surface-container-high text-on-surface-variant">
                        {cap}
                      </span>
                    ))}
                    {agent.capabilities.length > 3 && (
                      <span className="text-xs px-2 py-1 rounded bg-surface-container-high text-on-surface-variant">
                        +{agent.capabilities.length - 3}更多
                      </span>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
          
          {agents.length === 0 && (
            <div className="text-center py-12">
              <div className="material-symbols-rounded text-4xl text-on-surface-variant/50 mb-4">smart_toy</div>
              <h3 className="text-lg font-medium text-on-surface mb-2">暂无Agent</h3>
              <p className="text-sm text-on-surface-variant">请先创建或启用Agent</p>
            </div>
          )}
        </div>
      )}

      {/* 任务执行 */}
      {activeTab === 'execution' && (
        <div className="space-y-6">
          <div className="bg-surface-container rounded-xl p-6 border border-outline-variant">
            <h2 className="text-lg font-medium text-on-surface mb-4">执行任务</h2>
            
            <div className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-on-surface mb-2">选择Agent</label>
                <div className="relative">
                  <select
                    value={selectedAgent}
                    onChange={(e) => setSelectedAgent(e.target.value)}
                    className="w-full bg-surface-container-high border border-outline-variant rounded-lg px-4 py-3 text-on-surface focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary transition-all appearance-none"
                  >
                    <option value="">请选择Agent</option>
                    {agents.map(agent => (
                      <option key={agent.id} value={agent.id}>
                        {agent.name} ({agent.role})
                      </option>
                    ))}
                  </select>
                  <div className="absolute right-3 top-3 pointer-events-none">
                    <span className="material-symbols-rounded text-on-surface-variant">expand_more</span>
                  </div>
                </div>
              </div>

              {selectedAgent && (
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-on-surface mb-2">项目路径</label>
                    <input
                      type="text"
                      value={projectPath}
                      onChange={(e) => setProjectPath(e.target.value)}
                      className="w-full bg-surface-container-high border border-outline-variant rounded-lg px-4 py-3 text-on-surface focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary transition-all"
                      placeholder="./project-path"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-on-surface mb-2">项目类型</label>
                    <div className="relative">
                      <select
                        value={projectType}
                        onChange={(e) => setProjectType(e.target.value)}
                        className="w-full bg-surface-container-high border border-outline-variant rounded-lg px-4 py-3 text-on-surface focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary transition-all appearance-none"
                      >
                        <option value="go">Go项目</option>
                        <option value="nodejs">Node.js项目</option>
                        <option value="react">React项目</option>
                        <option value="python">Python项目</option>
                        <option value="static">静态网站</option>
                      </select>
                      <div className="absolute right-3 top-3 pointer-events-none">
                        <span className="material-symbols-rounded text-on-surface-variant">expand_more</span>
                      </div>
                    </div>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-on-surface mb-2">项目名称</label>
                    <input
                      type="text"
                      value={projectName}
                      onChange={(e) => setProjectName(e.target.value)}
                      className="w-full bg-surface-container-high border border-outline-variant rounded-lg px-4 py-3 text-on-surface focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary transition-all"
                      placeholder="MyProject"
                    />
                  </div>
                </div>
              )}

              <div>
                <label className="block text-sm font-medium text-on-surface mb-2">任务描述</label>
                <textarea
                  value={taskInput}
                  onChange={(e) => setTaskInput(e.target.value)}
                  className="w-full bg-surface-container-high border border-outline-variant rounded-lg px-4 py-3 text-on-surface focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary transition-all min-h-[120px] resize-y"
                  placeholder="例如：创建一个Go Web API项目，包含用户认证和数据库连接"
                  rows={4}
                />
                <p className="text-xs text-on-surface-variant mt-2">
                  示例任务：创建Go项目、生成React组件、修复代码错误、运行测试等
                </p>
              </div>

              <div className="flex justify-end">
                <button
                  onClick={handleExecuteTask}
                  disabled={!selectedAgent || !taskInput.trim() || loading}
                  className={`px-6 py-3 rounded-lg font-medium transition-all duration-200 ${
                    !selectedAgent || !taskInput.trim() || loading
                      ? 'bg-surface-container-high text-on-surface-variant cursor-not-allowed'
                      : 'bg-primary text-on-primary hover:bg-primary/90 active:scale-95'
                  }`}
                >
                  {loading ? (
                    <span className="flex items-center">
                      <span className="material-symbols-rounded animate-spin mr-2">refresh</span>
                      执行中...
                    </span>
                  ) : (
                    <span className="flex items-center">
                      <span className="material-symbols-rounded mr-2">play_arrow</span>
                      开始执行
                    </span>
                  )}
                </button>
              </div>
            </div>
          </div>

          {selectedAgent && (
            <div className="bg-surface-container rounded-xl p-6 border border-outline-variant">
              <h2 className="text-lg font-medium text-on-surface mb-4">Agent状态</h2>
              
              {(() => {
                const agent = getSelectedAgent();
                if (!agent) return null;
                
                const agentExecutions = executions.filter(exec => exec.agent_id === selectedAgent);
                const latestExecution = agentExecutions[0];
                
                return (
                  <div className="space-y-6">
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                      <div className="bg-surface-container-high rounded-lg p-4">
                        <div className="text-sm text-on-surface-variant mb-1">状态</div>
                        <div className={`text-lg font-semibold px-3 py-1.5 rounded-full inline-block ${getAgentStatusColor(latestExecution?.state || 'idle')}`}>
                          {latestExecution?.state || 'idle'}
                        </div>
                      </div>
                      <div className="bg-surface-container-high rounded-lg p-4">
                        <div className="text-sm text-on-surface-variant mb-1">执行次数</div>
                        <div className="text-2xl font-semibold text-on-surface">{agentExecutions.length}</div>
                      </div>
                      <div className="bg-surface-container-high rounded-lg p-4">
                        <div className="text-sm text-on-surface-variant mb-1">成功率</div>
                        <div className="text-2xl font-semibold text-on-surface">
                          {agentExecutions.length > 0 
                            ? `${Math.round((agentExecutions.filter(e => e.state === 'completed').length / agentExecutions.length) * 100)}%`
                            : '0%'
                          }
                        </div>
                      </div>
                      <div className="bg-surface-container-high rounded-lg p-4">
                        <div className="text-sm text-on-surface-variant mb-1">工具数量</div>
                        <div className="text-2xl font-semibold text-on-surface">{agent.tools.length}</div>
                      </div>
                    </div>

                    {latestExecution && (
                      <div className="border-t border-outline-variant pt-6">
                        <h3 className="font-medium text-on-surface mb-3">最新执行</h3>
                        <div className="bg-surface-container-high rounded-lg p-4">
                          <div className="flex justify-between items-start mb-3">
                            <div>
                              <div className="font-medium text-on-surface">{latestExecution.task}</div>
                              <div className="text-sm text-on-surface-variant mt-1">
                                开始时间: {formatDate(latestExecution.started_at)}
                              </div>
                            </div>
                            <div className={`px-3 py-1.5 rounded-full text-sm ${getExecutionStatusColor(latestExecution.state)}`}>
                              {latestExecution.state}
                            </div>
                          </div>
                          {latestExecution.error && (
                            <div className="mt-3 p-3 bg-red-500/10 text-red-400 rounded-lg text-sm">
                              <div className="font-medium mb-1">错误信息</div>
                              {latestExecution.error}
                            </div>
                          )}
                          {latestExecution.result && (
                            <div className="mt-3">
                              <details className="group">
                                <summary className="cursor-pointer text-sm text-primary hover:text-primary/80 flex items-center">
                                  <span className="material-symbols-rounded text-base mr-1 group-open:rotate-90 transition-transform">chevron_right</span>
                                  查看结果详情
                                </summary>
                                <div className="mt-3 p-3 bg-surface rounded-lg">
                                  <pre className="text-xs text-on-surface-variant overflow-auto max-h-60">
                                    {JSON.stringify(latestExecution.result, null, 2)}
                                  </pre>
                                </div>
                              </details>
                            </div>
                          )}
                        </div>
                      </div>
                    )}
                  </div>
                );
              })()}
            </div>
          )}
        </div>
      )}

      {/* 执行历史 */}
      {activeTab === 'history' && (
        <div className="bg-surface-container rounded-xl p-6 border border-outline-variant">
          <h2 className="text-lg font-medium text-on-surface mb-4">执行历史</h2>
          
          {executions.length === 0 ? (
            <div className="text-center py-12">
              <div className="material-symbols-rounded text-4xl text-on-surface-variant/50 mb-4">history</div>
              <h3 className="text-lg font-medium text-on-surface mb-2">暂无执行历史</h3>
              <p className="text-sm text-on-surface-variant">执行任务后，历史记录将显示在这里</p>
            </div>
          ) : (
            <div className="space-y-4">
              {executions.map(execution => (
                <div key={execution.id} className="bg-surface-container-high rounded-lg p-4 border border-outline-variant">
                  <div className="flex justify-between items-start mb-3">
                    <div>
                      <div className="font-medium text-on-surface">{execution.task}</div>
                      <div className="text-sm text-on-surface-variant mt-1">
                        Agent: {agents.find(a => a.id === execution.agent_id)?.name || execution.agent_id}
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <div className={`px-3 py-1.5 rounded-full text-sm ${getExecutionStatusColor(execution.state)}`}>
                        {execution.state}
                      </div>
                      {(execution.state === 'thinking' || execution.state === 'executing' || execution.state === 'waiting') && (
                        <button
                          onClick={() => handleCancelExecution(execution.id)}
                          className="px-3 py-1.5 text-sm border border-outline-variant text-on-surface-variant rounded-lg hover:bg-surface-container hover:text-on-surface transition-colors"
                        >
                          取消
                        </button>
                      )}
                    </div>
                  </div>
                  
                  <div className="text-sm text-on-surface-variant mb-3">
                    <span>开始: {formatDate(execution.started_at)}</span>
                    {execution.completed_at && (
                      <span className="ml-4">完成: {formatDate(execution.completed_at)}</span>
                    )}
                  </div>
                  
                  {execution.error && (
                    <div className="mt-2 p-3 bg-red-500/10 text-red-400 rounded-lg text-sm">
                      <div className="font-medium mb-1">错误信息</div>
                      {execution.error}
                    </div>
                  )}
                  
                  {execution.result && (
                    <div className="mt-2">
                      <details className="group">
                        <summary className="cursor-pointer text-sm text-primary hover:text-primary/80 flex items-center">
                          <span className="material-symbols-rounded text-base mr-1 group-open:rotate-90 transition-transform">chevron_right</span>
                          查看执行结果
                        </summary>
                        <div className="mt-3 p-3 bg-surface rounded-lg">
                          <pre className="text-xs text-on-surface-variant overflow-auto max-h-60">
                            {JSON.stringify(execution.result, null, 2)}
                          </pre>
                        </div>
                      </details>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* 错误提示 */}
      {error && (
        <div className="bg-red-500/10 border border-red-500/20 text-red-400 rounded-lg p-4">
          <div className="flex items-center">
            <span className="material-symbols-rounded mr-2">error</span>
            <span>{error}</span>
          </div>
        </div>
      )}
    </div>
  );
};

export default AgentViewNew;