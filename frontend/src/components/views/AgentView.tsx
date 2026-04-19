import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/Card';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Select } from '../ui/Select';
import { TextArea } from '../ui/TextArea';
import { useAgentStore } from '../../stores/agentStore';
import { useNotificationStore } from '../../stores/notificationStore';
import './AgentView.css';

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

const AgentView: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'agents' | 'execution' | 'config'>('agents');
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
    
    // 每10秒刷新一次状态
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

      // 根据任务类型添加额外参数
      if (taskInput.toLowerCase().includes('创建项目') || taskInput.toLowerCase().includes('create project')) {
        parameters.project_path = projectPath;
        parameters.project_type = projectType;
        parameters.project_name = projectName;
      }

      await executeAgent(selectedAgent, taskInput, parameters);
      addNotification('任务已开始执行', 'success');
      setTaskInput('');
      
      // 等待一会儿后刷新执行列表
      setTimeout(() => {
        fetchExecutions();
      }, 1000);
    } catch (err) {
      addNotification(`执行失败: ${err instanceof Error ? err.message : '未知错误'}`, 'error');
    }
  };

  const handleCancelExecution = async (executionId: string) => {
    try {
      // 这里需要实现取消执行的API调用
      addNotification('取消执行功能待实现', 'info');
    } catch (err) {
      addNotification(`取消失败: ${err instanceof Error ? err.message : '未知错误'}`, 'error');
    }
  };

  const getAgentStatusColor = (state: string) => {
    switch (state) {
      case 'idle': return 'bg-green-100 text-green-800';
      case 'thinking': return 'bg-blue-100 text-blue-800';
      case 'executing': return 'bg-yellow-100 text-yellow-800';
      case 'waiting': return 'bg-purple-100 text-purple-800';
      case 'completed': return 'bg-green-100 text-green-800';
      case 'failed': return 'bg-red-100 text-red-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  const getExecutionStatusColor = (state: string) => {
    switch (state) {
      case 'completed': return 'bg-green-100 text-green-800';
      case 'failed': return 'bg-red-100 text-red-800';
      case 'thinking':
      case 'executing':
      case 'waiting': return 'bg-yellow-100 text-yellow-800';
      default: return 'bg-gray-100 text-gray-800';
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
    <div className="agent-view">
      <div className="flex mb-6">
        <button
          className={`px-4 py-2 mr-2 rounded-lg ${activeTab === 'agents' ? 'bg-blue-600 text-white' : 'bg-gray-200'}`}
          onClick={() => setActiveTab('agents')}
        >
          Agent列表
        </button>
        <button
          className={`px-4 py-2 mr-2 rounded-lg ${activeTab === 'execution' ? 'bg-blue-600 text-white' : 'bg-gray-200'}`}
          onClick={() => setActiveTab('execution')}
        >
          任务执行
        </button>
        <button
          className={`px-4 py-2 rounded-lg ${activeTab === 'config' ? 'bg-blue-600 text-white' : 'bg-gray-200'}`}
          onClick={() => setActiveTab('config')}
        >
          执行历史
        </button>
      </div>

      {activeTab === 'agents' && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {agents.map(agent => (
            <Card 
              key={agent.id} 
              className={`cursor-pointer ${selectedAgent === agent.id ? 'ring-2 ring-blue-500' : ''}`}
              onClick={() => setSelectedAgent(agent.id)}
            >
              <CardHeader>
                <CardTitle className="flex justify-between items-center">
                  <span>{agent.name}</span>
                  <span className={`px-2 py-1 text-xs rounded-full ${agent.enabled ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}`}>
                    {agent.enabled ? '启用' : '禁用'}
                  </span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-gray-600 mb-2">{agent.description}</p>
                <div className="mb-2">
                  <span className="text-xs font-semibold">角色: </span>
                  <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">
                    {agent.role}
                  </span>
                </div>
                <div className="mb-2">
                  <span className="text-xs font-semibold">模型: </span>
                  <span className="text-xs">{agent.model_config.name} ({agent.model_config.provider})</span>
                </div>
                <div className="mb-2">
                  <span className="text-xs font-semibold">能力: </span>
                  <div className="flex flex-wrap gap-1 mt-1">
                    {agent.capabilities.slice(0, 3).map((cap, index) => (
                      <span key={index} className="text-xs bg-gray-100 text-gray-800 px-2 py-1 rounded">
                        {cap}
                      </span>
                    ))}
                    {agent.capabilities.length > 3 && (
                      <span className="text-xs text-gray-500">+{agent.capabilities.length - 3}更多</span>
                    )}
                  </div>
                </div>
                <div>
                  <span className="text-xs font-semibold">工具: </span>
                  <span className="text-xs">{agent.tools.length}个工具</span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {activeTab === 'execution' && (
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>执行任务</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium mb-1">选择Agent</label>
                  <Select
                    value={selectedAgent}
                    onChange={(e) => setSelectedAgent(e.target.value)}
                    options={agents.map(agent => ({
                      value: agent.id,
                      label: `${agent.name} (${agent.role})`
                    }))}
                    placeholder="请选择Agent"
                  />
                </div>

                {selectedAgent && (
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                    <div>
                      <label className="block text-sm font-medium mb-1">项目路径</label>
                      <Input
                        type="text"
                        value={projectPath}
                        onChange={(e) => setProjectPath(e.target.value)}
                        placeholder="./project-path"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium mb-1">项目类型</label>
                      <Select
                        value={projectType}
                        onChange={(e) => setProjectType(e.target.value)}
                        options={[
                          { value: 'go', label: 'Go项目' },
                          { value: 'nodejs', label: 'Node.js项目' },
                          { value: 'react', label: 'React项目' },
                          { value: 'python', label: 'Python项目' },
                          { value: 'static', label: '静态网站' }
                        ]}
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium mb-1">项目名称</label>
                      <Input
                        type="text"
                        value={projectName}
                        onChange={(e) => setProjectName(e.target.value)}
                        placeholder="MyProject"
                      />
                    </div>
                  </div>
                )}

                <div>
                  <label className="block text-sm font-medium mb-1">任务描述</label>
                  <TextArea
                    value={taskInput}
                    onChange={(e) => setTaskInput(e.target.value)}
                    placeholder="例如：创建一个Go Web API项目，包含用户认证和数据库连接"
                    rows={4}
                  />
                  <p className="text-xs text-gray-500 mt-1">
                    示例任务：创建Go项目、生成React组件、修复代码错误、运行测试等
                  </p>
                </div>

                <div className="flex justify-end">
                  <Button
                    onClick={handleExecuteTask}
                    disabled={!selectedAgent || !taskInput.trim() || loading}
                    className="px-6"
                  >
                    {loading ? '执行中...' : '开始执行'}
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {selectedAgent && (
            <Card>
              <CardHeader>
                <CardTitle>Agent状态</CardTitle>
              </CardHeader>
              <CardContent>
                {(() => {
                  const agent = getSelectedAgent();
                  if (!agent) return null;
                  
                  const agentExecutions = executions.filter(exec => exec.agent_id === selectedAgent);
                  const latestExecution = agentExecutions[0];
                  
                  return (
                    <div className="space-y-4">
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        <div className="bg-gray-50 p-4 rounded-lg">
                          <div className="text-sm text-gray-500">状态</div>
                          <div className={`text-lg font-semibold ${getAgentStatusColor(latestExecution?.state || 'idle')} px-2 py-1 rounded inline-block`}>
                            {latestExecution?.state || 'idle'}
                          </div>
                        </div>
                        <div className="bg-gray-50 p-4 rounded-lg">
                          <div className="text-sm text-gray-500">执行次数</div>
                          <div className="text-lg font-semibold">{agentExecutions.length}</div>
                        </div>
                        <div className="bg-gray-50 p-4 rounded-lg">
                          <div className="text-sm text-gray-500">成功率</div>
                          <div className="text-lg font-semibold">
                            {agentExecutions.length > 0 
                              ? `${Math.round((agentExecutions.filter(e => e.state === 'completed').length / agentExecutions.length) * 100)}%`
                              : '0%'
                            }
                          </div>
                        </div>
                        <div className="bg-gray-50 p-4 rounded-lg">
                          <div className="text-sm text-gray-500">工具数量</div>
                          <div className="text-lg font-semibold">{agent.tools.length}</div>
                        </div>
                      </div>

                      {latestExecution && (
                        <div className="border-t pt-4">
                          <h4 className="font-medium mb-2">最新执行</h4>
                          <div className="bg-gray-50 p-4 rounded-lg">
                            <div className="flex justify-between items-start mb-2">
                              <div>
                                <div className="font-medium">{latestExecution.task}</div>
                                <div className="text-sm text-gray-500">
                                  开始时间: {formatDate(latestExecution.started_at)}
                                </div>
                              </div>
                              <div className={`px-3 py-1 rounded-full text-sm ${getExecutionStatusColor(latestExecution.state)}`}>
                                {latestExecution.state}
                              </div>
                            </div>
                            {latestExecution.error && (
                              <div className="mt-2 p-2 bg-red-50 text-red-700 rounded text-sm">
                                <strong>错误:</strong> {latestExecution.error}
                              </div>
                            )}
                            {latestExecution.result && (
                              <div className="mt-2">
                                <button
                                  className="text-sm text-blue-600 hover:text-blue-800"
                                  onClick={() => {
                                    // 显示详细结果
                                    console.log('Execution result:', latestExecution.result);
                                    addNotification('请在控制台查看详细结果', 'info');
                                  }}
                                >
                                  查看结果详情
                                </button>
                              </div>
                            )}
                          </div>
                        </div>
                      )}
                    </div>
                  );
                })()}
              </CardContent>
            </Card>
          )}
        </div>
      )}

      {activeTab === 'config' && (
        <Card>
          <CardHeader>
            <CardTitle>执行历史</CardTitle>
          </CardHeader>
          <CardContent>
            {executions.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                暂无执行历史
              </div>
            ) : (
              <div className="space-y-4">
                {executions.map(execution => (
                  <div key={execution.id} className="border rounded-lg p-4">
                    <div className="flex justify-between items-start mb-2">
                      <div>
                        <div className="font-medium">{execution.task}</div>
                        <div className="text-sm text-gray-500">
                          Agent: {agents.find(a => a.id === execution.agent_id)?.name || execution.agent_id}
                        </div>
                      </div>
                      <div className="flex items-center space-x-2">
                        <div className={`px-3 py-1 rounded-full text-sm ${getExecutionStatusColor(execution.state)}`}>
                          {execution.state}
                        </div>
                        {execution.state === 'thinking' || execution.state === 'executing' || execution.state === 'waiting' ? (
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => handleCancelExecution(execution.id)}
                          >
                            取消
                          </Button>
                        ) : null}
                      </div>
                    </div>
                    <div className="text-sm text-gray-600 mb-2">
                      <span>开始: {formatDate(execution.started_at)}</span>
                      {execution.completed_at && (
                        <span className="ml-4">完成: {formatDate(execution.completed_at)}</span>
                      )}
                    </div>
                    {execution.error && (
                      <div className="mt-2 p-2 bg-red-50 text-red-700 rounded text-sm">
                        <strong>错误:</strong> {execution.error}
                      </div>
                    )}
                    {execution.result && (
                      <div className="mt-2">
                        <details>
                          <summary className="cursor-pointer text-sm text-blue-600 hover:text-blue-800">
                            查看结果
                          </summary>
                          <pre className="mt-2 p-2 bg-gray-50 rounded text-xs overflow-auto max-h-60">
                            {JSON.stringify(execution.result, null, 2)}
                          </pre>
                        </details>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {error && (
        <div className="mt-4 p-4 bg-red-50 text-red-700 rounded-lg">
          {error}
        </div>
      )}
    </div>
  );
};

export default AgentView;