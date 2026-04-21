import React, { useState, useEffect } from 'react';
import { 
  Server, 
  Plus as Add, 
  Edit as EditIcon,
  Trash2, 
  RefreshCw,
  TestTube,
  Settings as SettingsIcon,
  ChevronDown,
  ChevronUp,
  Save,
  Database,
  CheckCircle2,
  XCircle,
  Layers,
  Power,
  PowerOff,
  Globe,
  Cpu,
  Cloud,
  Terminal,
  Wrench,
  Zap,
  X
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';
import { api, MCPConfig, MCPServer, MCPTool, MCPServerStatus } from '../services/api';

// 辅助函数：获取服务器图标
const getServerIcon = (type: string) => {
  switch (type) {
    case 'http':
      return <Globe className="w-5 h-5 text-primary" />;
    case 'stdio':
      return <Terminal className="w-5 h-5 text-primary" />;
    case 'sse':
      return <Cloud className="w-5 h-5 text-primary" />;
    default:
      return <Server className="w-5 h-5 text-primary" />;
  }
};

// 辅助函数：获取状态图标
const getStatusIcon = (connected: boolean) => {
  return connected ? (
    <CheckCircle2 className="w-3 h-3 text-success" />
  ) : (
    <XCircle className="w-3 h-3 text-error" />
  );
};

// 辅助函数：获取状态颜色类名
const getStatusColor = (connected: boolean) => {
  return connected ? 'text-success' : 'text-error';
};

export default function MCPConfigView() {
  const [config, setConfig] = useState<MCPConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [tools, setTools] = useState<MCPTool[]>([]);
  const [expandedServer, setExpandedServer] = useState<string | null>(null);
  const [serverStatuses, setServerStatuses] = useState<Record<string, MCPServerStatus>>({});
  const [editingServer, setEditingServer] = useState<MCPServer | null>(null);
  const [showEditModal, setShowEditModal] = useState(false);

  // 加载配置
  useEffect(() => {
    loadConfig();
    loadTools();
  }, []);

  // 加载服务器状态
  useEffect(() => {
    if (config?.servers) {
      loadServerStatuses();
    }
  }, [config]);

  const loadConfig = async () => {
    try {
      setLoading(true);
      const configData = await api.getMCPConfig();
      setConfig(configData);
    } catch (error) {
      console.error('Failed to load MCP config:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadTools = async () => {
    try {
      const toolsData = await api.getMCPTools();
      setTools(toolsData);
    } catch (error) {
      console.error('Failed to load MCP tools:', error);
    }
  };

  const loadServerStatuses = async () => {
    if (!config?.servers) return;

    const statuses: Record<string, MCPServerStatus> = {};
    for (const server of config.servers) {
      try {
        const status = await api.getMCPServerStatus(server.id);
        statuses[server.id] = status;
      } catch (error) {
        console.error(`Failed to load status for server ${server.id}:`, error);
      }
    }
    setServerStatuses(statuses);
  };

  const handleSaveConfig = async () => {
    if (!config) return;

    try {
      const success = await api.updateMCPConfig(config);
      if (success) {
        alert('配置已保存');
        loadConfig();
      } else {
        alert('保存失败');
      }
    } catch (error) {
      console.error('Failed to save config:', error);
      alert('保存失败');
    }
  };

  const handleRemoveServer = async (serverId: string) => {
    if (!confirm('确定要删除这个服务器吗？')) return;

    try {
      const success = await api.removeMCPServer(serverId);
      if (success) {
        loadConfig();
      } else {
        alert('删除失败');
      }
    } catch (error) {
      console.error('Failed to remove server:', error);
      alert('删除失败');
    }
  };

  const handleUpdateServer = async (serverId: string, updates: Partial<MCPServer>) => {
    if (!config) return;

    const updatedServers = config.servers.map(server => 
      server.id === serverId ? { ...server, ...updates } : server
    );

    setConfig({ ...config, servers: updatedServers });
  };

  const handleEditServer = (server: MCPServer) => {
    setEditingServer(server);
    setShowEditModal(true);
  };

  const handleSaveEditedServer = async () => {
    if (!editingServer || !config) return;

    try {
      const success = await api.updateMCPServer(editingServer.id, editingServer);
      if (success) {
        // 更新本地配置
        const updatedServers = config.servers.map(server => 
          server.id === editingServer.id ? editingServer : server
        );
        setConfig({ ...config, servers: updatedServers });
        setShowEditModal(false);
        setEditingServer(null);
        alert('服务器配置已更新');
      } else {
        alert('更新失败');
      }
    } catch (error) {
      console.error('Failed to update server:', error);
      alert('更新失败');
    }
  };

  const handleAddNewServer = () => {
    // 创建新的服务器配置
    const newServer: MCPServer = {
      id: `server-${Date.now()}`,
      name: '新服务器',
      type: 'client',
      enabled: true,
      priority: 50,
      client_config: {
        server_type: 'stdio',
        command: '',
        args: [],
        namespace: '',
        timeout: 30
      }
    };
    setEditingServer(newServer);
    setShowEditModal(true);
  };

  const handleSaveNewServer = async () => {
    if (!editingServer || !config) return;

    try {
      const savedServer = await api.addMCPServer(editingServer);
      if (savedServer) {
        // 添加新服务器到配置
        const updatedServers = [...config.servers, editingServer];
        setConfig({ ...config, servers: updatedServers });
        setShowEditModal(false);
        setEditingServer(null);
        alert('服务器已添加');
        loadConfig(); // 重新加载配置以确保同步
      } else {
        alert('添加失败');
      }
    } catch (error) {
      console.error('Failed to add server:', error);
      alert('添加失败');
    }
  };

  const handleTestConnection = async (serverId: string) => {
    try {
      const result = await api.testMCPServerConnection(serverId);
      if (result.success) {
        alert('连接测试成功');
      } else {
        alert(`连接测试失败: ${result.error}`);
      }
    } catch (error) {
      console.error('Failed to test connection:', error);
      alert('测试失败');
    }
  };

  const handleRefreshServer = async (serverId: string) => {
    try {
      const status = await api.getMCPServerStatus(serverId);
      setServerStatuses(prev => ({ ...prev, [serverId]: status }));
      alert('服务器状态已刷新');
    } catch (error) {
      console.error('Failed to refresh server:', error);
      alert('刷新失败');
    }
  };

  const getServerIcon = (serverType: string) => {
    switch (serverType) {
      case 'client':
        return <Server size={20} className="text-primary" />;
      case 'server':
        return <Database size={20} className="text-success" />;
      default:
        return <Server size={20} className="text-outline" />;
    }
  };

  const getStatusColor = (connected: boolean) => {
    return connected ? 'text-success' : 'text-error';
  };

  const getStatusIcon = (connected: boolean) => {
    return connected ? (
      <CheckCircle2 size={16} className="text-success" />
    ) : (
      <XCircle size={16} className="text-error" />
    );
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-outline">加载中...</div>
      </div>
    );
  }

  if (!config) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-error">无法加载配置</div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* 编辑模态框 */}
      {showEditModal && editingServer && (
        <EditServerModal
          server={editingServer}
          onClose={() => {
            setShowEditModal(false);
            setEditingServer(null);
          }}
          onSave={editingServer.id.startsWith('server-') ? handleSaveNewServer : handleSaveEditedServer}
          onChange={(updates) => setEditingServer({ ...editingServer, ...updates })}
          isNew={editingServer.id.startsWith('server-')}
        />
      )}

      {/* 全局配置 */}
      <section className="bg-surface-container-low rounded-3xl p-6 border border-outline-variant/10">
        <h3 className="text-xl font-bold text-primary mb-6 flex items-center gap-2">
          <SettingsIcon size={20} />
          全局配置
        </h3>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <label className="text-sm font-bold text-on-surface mb-2 block">自动启动</label>
            <div className="flex items-center gap-3">
              <button
                onClick={() => setConfig({ ...config, auto_start: true })}
                className={cn(
                  "px-4 py-2 rounded-xl text-sm font-bold transition-all",
                  config.auto_start
                    ? "bg-primary-container text-on-primary-container"
                    : "bg-surface-container-highest text-outline hover:text-on-surface"
                )}
              >
                启用
              </button>
              <button
                onClick={() => setConfig({ ...config, auto_start: false })}
                className={cn(
                  "px-4 py-2 rounded-xl text-sm font-bold transition-all",
                  !config.auto_start
                    ? "bg-error-container text-on-error-container"
                    : "bg-surface-container-highest text-outline hover:text-on-surface"
                )}
              >
                禁用
              </button>
            </div>
          </div>

          <div>
            <label className="text-sm font-bold text-on-surface mb-2 block">工具发现间隔（秒）</label>
            <input
              type="number"
              value={config.tool_discovery_interval}
              onChange={(e) => setConfig({ ...config, tool_discovery_interval: parseInt(e.target.value) || 60 })}
              className="w-full bg-surface-container-highest/50 border border-outline-variant/10 rounded-xl px-4 py-3 text-on-surface font-mono"
              min="10"
              max="3600"
            />
          </div>

          <div>
            <label className="text-sm font-bold text-on-surface mb-2 block">最大并发调用数</label>
            <input
              type="number"
              value={config.max_concurrent_calls}
              onChange={(e) => setConfig({ ...config, max_concurrent_calls: parseInt(e.target.value) || 10 })}
              className="w-full bg-surface-container-highest/50 border border-outline-variant/10 rounded-xl px-4 py-3 text-on-surface font-mono"
              min="1"
              max="100"
            />
          </div>
        </div>
      </section>

      {/* 服务器列表 */}
      <section className="bg-surface-container-low rounded-3xl p-6 border border-outline-variant/10">
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-xl font-bold text-primary flex items-center gap-2">
            <Server size={20} />
            MCP 服务器配置 ({config.servers.length})
          </h3>
          <button
            onClick={handleAddNewServer}
            className="bg-primary-container text-on-primary-container hover:brightness-110 rounded-xl px-5 py-2.5 text-sm font-bold shadow-xl shadow-primary/20 transition-all flex items-center gap-2"
          >
            <Add size={18} />
            添加服务器
          </button>
        </div>

        <div className="space-y-4">
          {config.servers.length === 0 ? (
            <div className="text-center py-12 text-outline">
              <Server size={48} className="mx-auto mb-4 opacity-30" />
              <p>暂无服务器配置</p>
            </div>
          ) : (
            config.servers.map((server) => (
              <div key={server.id}>
                <ServerCard
                  server={server}
                  status={serverStatuses[server.id]}
                  expanded={expandedServer === server.id}
                  onToggleExpand={() => setExpandedServer(expandedServer === server.id ? null : server.id)}
                  onUpdate={(updates) => handleUpdateServer(server.id, updates)}
                  onRemove={() => handleRemoveServer(server.id)}
                  onTestConnection={() => handleTestConnection(server.id)}
                  onRefresh={() => handleRefreshServer(server.id)}
                  onEdit={() => handleEditServer(server)}
                  tools={tools.filter(tool => tool.server_id === server.id)}
                />
              </div>
            ))
          )}
        </div>
      </section>

      {/* 工具列表 */}
      <section className="bg-surface-container-low rounded-3xl p-6 border border-outline-variant/10">
        <h3 className="text-xl font-bold text-primary mb-6 flex items-center gap-2">
          <Layers size={20} />
          可用工具 ({tools.length})
        </h3>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {tools.map((tool) => (
            <div key={`${tool.server_id}-${tool.name}`}>
              <ToolCard tool={tool} />
            </div>
          ))}
        </div>
      </section>

      {/* 保存按钮 */}
      <div className="fixed bottom-10 right-10 z-50">
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-surface-variant/90 backdrop-blur-3xl p-2 rounded-2xl shadow-2xl border border-white/10 flex items-center gap-3"
        >
          <button 
            onClick={() => loadConfig()}
            className="bg-surface-container-lowest text-on-surface hover:bg-surface-bright rounded-xl px-5 py-2.5 text-sm font-bold transition-all border border-outline-variant/10 flex items-center gap-2"
          >
            <RefreshCw size={16} />
            刷新
          </button>
          <button 
            onClick={handleSaveConfig}
            className="bg-primary-container text-on-primary-container hover:brightness-110 rounded-xl px-7 py-2.5 text-sm font-bold shadow-xl shadow-primary/20 transition-all flex items-center gap-2"
          >
            <Save size={18} />
            保存配置
          </button>
        </motion.div>
      </div>
    </div>
  );
}

interface ServerCardProps {
  server: MCPServer;
  status?: MCPServerStatus;
  expanded: boolean;
  onToggleExpand: () => void;
  onUpdate: (updates: Partial<MCPServer>) => void;
  onRemove: () => void;
  onTestConnection: () => void;
  onRefresh: () => void;
  onEdit: () => void;
  tools: MCPTool[];
}

function ServerCard({
  server,
  status,
  expanded,
  onToggleExpand,
  onUpdate,
  onRemove,
  onTestConnection,
  onRefresh,
  onEdit,
  tools
}: ServerCardProps) {
  return (
    <div className="bg-surface-container-high rounded-2xl border border-outline-variant/10 overflow-hidden">
      {/* 服务器头部 */}
      <div className="p-4 flex items-center justify-between hover:bg-white/[0.02] transition-colors">
        <div className="flex items-center gap-4 flex-1">
          <div className="p-2 rounded-xl bg-surface-container-lowest">
            {getServerIcon(server.type)}
          </div>
          
          <div className="flex-1">
            <div className="flex items-center gap-3">
              <h4 className="text-base font-bold text-on-surface">{server.name}</h4>
              <span className="text-[10px] font-mono font-bold text-outline uppercase tracking-widest bg-surface-container-highest/50 px-2 py-1 rounded">
                {server.type}
              </span>
              {status && (
                <div className="flex items-center gap-1">
                  {getStatusIcon(status.connected)}
                  <span className={cn("text-[10px] font-bold uppercase tracking-widest", getStatusColor(status.connected))}>
                    {status.connected ? '已连接' : '未连接'}
                  </span>
                </div>
              )}
            </div>
            
            <div className="flex items-center gap-4 mt-2">
              <div className="flex items-center gap-2">
                <span className="text-[10px] text-outline uppercase tracking-widest">ID:</span>
                <span className="text-xs font-mono text-on-surface-variant">{server.id}</span>
              </div>
              
              <div className="flex items-center gap-2">
                <span className="text-[10px] text-outline uppercase tracking-widest">优先级:</span>
                <span className="text-xs font-mono text-on-surface-variant">{server.priority}</span>
              </div>
              
              <div className="flex items-center gap-2">
                <span className="text-[10px] text-outline uppercase tracking-widest">工具:</span>
                <span className="text-xs font-mono text-on-surface-variant">{tools.length}</span>
              </div>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={() => onUpdate({ enabled: !server.enabled })}
            className={cn(
              "text-[10px] font-bold uppercase tracking-widest px-3 py-1.5 rounded-lg",
              server.enabled
                ? "bg-success-container text-on-success-container hover:brightness-110"
                : "bg-error-container text-on-error-container hover:brightness-110"
            )}
          >
            {server.enabled ? '启用' : '禁用'}
          </button>
          
          <button
            onClick={onTestConnection}
            className="text-[10px] font-bold uppercase tracking-widest bg-primary-container text-on-primary-container hover:brightness-110 px-3 py-1.5 rounded-lg flex items-center gap-1"
          >
            <TestTube size={12} />
            测试
          </button>
          
          <button
            onClick={onRefresh}
            className="text-[10px] font-bold uppercase tracking-widest bg-surface-container-highest text-on-surface hover:bg-surface-bright px-3 py-1.5 rounded-lg flex items-center gap-1"
          >
            <RefreshCw size={12} />
            刷新
          </button>
          
          <button
            onClick={onEdit}
            className="text-[10px] font-bold uppercase tracking-widest bg-surface-container-highest text-on-surface hover:bg-surface-bright px-3 py-1.5 rounded-lg flex items-center gap-1"
          >
            <EditIcon size={12} />
            编辑
          </button>
          
          <button
            onClick={onToggleExpand}
            className="text-[10px] font-bold uppercase tracking-widest bg-surface-container-highest text-on-surface hover:bg-surface-bright px-3 py-1.5 rounded-lg"
          >
            {expanded ? <ChevronUp size={12} /> : <ChevronDown size={12} />}
          </button>
          
          <button
            onClick={onRemove}
            className="text-[10px] font-bold uppercase tracking-widest bg-error-container text-on-error-container hover:brightness-110 px-3 py-1.5 rounded-lg flex items-center gap-1"
          >
            <Trash2 size={12} />
          </button>
        </div>
      </div>

      {/* 展开的详细信息 */}
      {expanded && (
        <div className="border-t border-outline-variant/10 p-4 bg-surface-container-lowest/50">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* 服务器配置详情 */}
            <div>
              <h5 className="text-sm font-bold text-on-surface mb-3">服务器配置</h5>
              <div className="space-y-3">
                {server.type === 'client' && server.client_config && (
                  <>
                    <div>
                      <label className="text-[10px] uppercase font-bold text-outline tracking-widest block mb-1">命令</label>
                      <div className="bg-surface-container-highest/50 rounded-lg px-3 py-2 font-mono text-sm">
                        {server.client_config.command} {server.client_config.args.join(' ')}
                      </div>
                    </div>
                    <div>
                      <label className="text-[10px] uppercase font-bold text-outline tracking-widest block mb-1">命名空间</label>
                      <div className="bg-surface-container-highest/50 rounded-lg px-3 py-2 font-mono text-sm">
                        {server.client_config.namespace}
                      </div>
                    </div>
                    <div>
                      <label className="text-[10px] uppercase font-bold text-outline tracking-widest block mb-1">超时时间</label>
                      <div className="bg-surface-container-highest/50 rounded-lg px-3 py-2 font-mono text-sm">
                        {server.client_config.timeout} 秒
                      </div>
                    </div>
                  </>
                )}
                
                {server.type === 'server' && server.server_config && (
                  <>
                    <div>
                      <label className="text-[10px] uppercase font-bold text-outline tracking-widest block mb-1">服务器类型</label>
                      <div className="bg-surface-container-highest/50 rounded-lg px-3 py-2 font-mono text-sm">
                        {server.server_config.type}
                      </div>
                    </div>
                    {server.server_config.tool_filter && (
                      <div>
                        <label className="text-[10px] uppercase font-bold text-outline tracking-widest block mb-1">工具过滤器</label>
                        <div className="bg-surface-container-highest/50 rounded-lg px-3 py-2 font-mono text-sm">
                          {server.server_config.tool_filter}
                        </div>
                      </div>
                    )}
                  </>
                )}
              </div>
            </div>

            {/* 工具列表 */}
            <div>
              <h5 className="text-sm font-bold text-on-surface mb-3">可用工具 ({tools.length})</h5>
              {tools.length === 0 ? (
                <div className="text-outline text-sm">暂无工具</div>
              ) : (
                <div className="space-y-2 max-h-48 overflow-y-auto">
                  {tools.map((tool) => (
                    <div key={tool.name} className="bg-surface-container-highest/30 rounded-lg p-3">
                      <div className="flex justify-between items-start mb-1">
                        <span className="font-mono text-sm font-bold text-primary">{tool.name}</span>
                        <span className="text-[10px] text-outline uppercase tracking-widest">{tool.namespace}</span>
                      </div>
                      <p className="text-xs text-on-surface-variant">{tool.description}</p>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* 状态信息 */}
          {status && (
            <div className="mt-4 pt-4 border-t border-outline-variant/10">
              <h5 className="text-sm font-bold text-on-surface mb-3">状态信息</h5>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div>
                  <label className="text-[10px] uppercase font-bold text-outline tracking-widest block mb-1">连接状态</label>
                  <div className={cn("flex items-center gap-2", getStatusColor(status.connected))}>
                    {getStatusIcon(status.connected)}
                    <span className="text-sm font-bold">{status.connected ? '已连接' : '未连接'}</span>
                  </div>
                </div>
                <div>
                  <label className="text-[10px] uppercase font-bold text-outline tracking-widest block mb-1">工具数量</label>
                  <div className="text-sm font-bold text-on-surface">{status.tool_count}</div>
                </div>
                <div>
                  <label className="text-[10px] uppercase font-bold text-outline tracking-widest block mb-1">最后更新</label>
                  <div className="text-sm font-mono text-on-surface-variant">
                    {new Date(status.last_update).toLocaleString()}
                  </div>
                </div>
                {status.last_error && (
                  <div>
                    <label className="text-[10px] uppercase font-bold text-outline tracking-widest block mb-1">最后错误</label>
                    <div className="text-sm font-mono text-error truncate" title={status.last_error}>
                      {status.last_error}
                    </div>
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

interface ToolCardProps {
  tool: MCPTool;
}

function ToolCard({ tool }: ToolCardProps) {
  const [expanded, setExpanded] = useState(false);

  return (
    <div className="bg-surface-container-high rounded-2xl p-4 border border-outline-variant/10 hover:border-outline-variant/30 transition-all">
      <div className="flex justify-between items-start mb-3">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <span className="font-mono text-sm font-bold text-primary">{tool.name}</span>
            <span className="text-[10px] text-outline uppercase tracking-widest bg-surface-container-highest/50 px-2 py-0.5 rounded">
              {tool.namespace}
            </span>
          </div>
          <p className="text-xs text-on-surface-variant line-clamp-2">{tool.description}</p>
        </div>
        <button
          onClick={() => setExpanded(!expanded)}
          className="text-outline hover:text-on-surface transition-colors"
        >
          {expanded ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
        </button>
      </div>

      <div className="flex items-center justify-between text-xs text-outline mb-3">
        <span className="truncate" title={tool.server_name}>
          {tool.server_name}
        </span>
        <span className="font-mono">{tool.server_id}</span>
      </div>

      {expanded && tool.input_schema.properties && (
        <div className="mt-3 pt-3 border-t border-outline-variant/10">
          <h6 className="text-[10px] uppercase font-bold text-outline tracking-widest mb-2">输入参数</h6>
          <div className="space-y-2">
            {Object.entries(tool.input_schema.properties).map(([key, schema]) => (
              <div key={key} className="flex items-center justify-between text-xs">
                <span className="font-mono text-primary">{key}</span>
                <span className="text-outline">{schema.type}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

interface EditServerModalProps {
  server: MCPServer;
  onClose: () => void;
  onSave: () => void;
  onChange: (updates: Partial<MCPServer>) => void;
  isNew: boolean;
}

function EditServerModal({ server, onClose, onSave, onChange, isNew }: EditServerModalProps) {
  const handleInputChange = (field: string, value: any) => {
    onChange({ [field]: value });
  };

  const handleClientConfigChange = (field: string, value: any) => {
    const updatedConfig = { ...server.client_config, [field]: value };
    onChange({ client_config: updatedConfig });
  };

  const handleArgsChange = (index: number, value: string) => {
    const args = [...(server.client_config?.args || [])];
    args[index] = value;
    handleClientConfigChange('args', args);
  };

  const addArg = () => {
    const args = [...(server.client_config?.args || []), ''];
    handleClientConfigChange('args', args);
  };

  const removeArg = (index: number) => {
    const args = [...(server.client_config?.args || [])];
    args.splice(index, 1);
    handleClientConfigChange('args', args);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div className="bg-surface-container-low rounded-3xl p-6 border border-outline-variant/10 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-xl font-bold text-primary flex items-center gap-2">
            <EditIcon size={20} />
            {isNew ? '添加新服务器' : '编辑服务器配置'}
          </h3>
          <button
            onClick={onClose}
            className="text-outline hover:text-on-surface transition-colors"
          >
            <X size={24} />
          </button>
        </div>

        <div className="space-y-6">
          {/* 基本信息 */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-bold text-on-surface mb-2 block">服务器名称</label>
              <input
                type="text"
                value={server.name}
                onChange={(e) => handleInputChange('name', e.target.value)}
                className="w-full bg-surface-container-highest/50 border border-outline-variant/10 rounded-xl px-4 py-3 text-on-surface"
                placeholder="输入服务器名称"
              />
            </div>

            <div>
              <label className="text-sm font-bold text-on-surface mb-2 block">服务器ID</label>
              <input
                type="text"
                value={server.id}
                onChange={(e) => handleInputChange('id', e.target.value)}
                className="w-full bg-surface-container-highest/50 border border-outline-variant/10 rounded-xl px-4 py-3 text-on-surface font-mono"
                placeholder="输入唯一ID"
                disabled={!isNew}
              />
              {!isNew && (
                <p className="text-xs text-outline mt-1">服务器ID创建后不可修改</p>
              )}
            </div>

            <div>
              <label className="text-sm font-bold text-on-surface mb-2 block">服务器类型</label>
              <div className="flex gap-2">
                <button
                  onClick={() => handleInputChange('type', 'client')}
                  className={cn(
                    "flex-1 px-4 py-3 rounded-xl text-sm font-bold transition-all",
                    server.type === 'client'
                      ? "bg-primary-container text-on-primary-container"
                      : "bg-surface-container-highest text-outline hover:text-on-surface"
                  )}
                >
                  客户端
                </button>
                <button
                  onClick={() => handleInputChange('type', 'server')}
                  className={cn(
                    "flex-1 px-4 py-3 rounded-xl text-sm font-bold transition-all",
                    server.type === 'server'
                      ? "bg-primary-container text-on-primary-container"
                      : "bg-surface-container-highest text-outline hover:text-on-surface"
                  )}
                >
                  服务器
                </button>
              </div>
            </div>

            <div>
              <label className="text-sm font-bold text-on-surface mb-2 block">优先级 (0-100)</label>
              <input
                type="number"
                value={server.priority}
                onChange={(e) => handleInputChange('priority', parseInt(e.target.value) || 50)}
                className="w-full bg-surface-container-highest/50 border border-outline-variant/10 rounded-xl px-4 py-3 text-on-surface"
                min="0"
                max="100"
              />
            </div>
          </div>

          {/* 启用状态 */}
          <div>
            <label className="text-sm font-bold text-on-surface mb-2 block">启用状态</label>
            <div className="flex gap-2">
              <button
                onClick={() => handleInputChange('enabled', true)}
                className={cn(
                  "flex-1 px-4 py-3 rounded-xl text-sm font-bold transition-all",
                  server.enabled
                    ? "bg-success-container text-on-success-container"
                    : "bg-surface-container-highest text-outline hover:text-on-surface"
                )}
              >
                启用
              </button>
              <button
                onClick={() => handleInputChange('enabled', false)}
                className={cn(
                  "flex-1 px-4 py-3 rounded-xl text-sm font-bold transition-all",
                  !server.enabled
                    ? "bg-error-container text-on-error-container"
                    : "bg-surface-container-highest text-outline hover:text-on-surface"
                )}
              >
                禁用
              </button>
            </div>
          </div>

          {/* 客户端配置 */}
          {server.type === 'client' && (
            <div className="border-t border-outline-variant/10 pt-6">
              <h4 className="text-lg font-bold text-primary mb-4">客户端配置</h4>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-bold text-on-surface mb-2 block">服务器类型</label>
                  <input
                    type="text"
                    value={server.client_config?.server_type || ''}
                    onChange={(e) => handleClientConfigChange('server_type', e.target.value)}
                    className="w-full bg-surface-container-highest/50 border border-outline-variant/10 rounded-xl px-4 py-3 text-on-surface font-mono"
                    placeholder="例如: stdio"
                  />
                </div>

                <div>
                  <label className="text-sm font-bold text-on-surface mb-2 block">命令</label>
                  <input
                    type="text"
                    value={server.client_config?.command || ''}
                    onChange={(e) => handleClientConfigChange('command', e.target.value)}
                    className="w-full bg-surface-container-highest/50 border border-outline-variant/10 rounded-xl px-4 py-3 text-on-surface font-mono"
                    placeholder="例如: npx, python, node"
                  />
                </div>

                <div className="md:col-span-2">
                  <label className="text-sm font-bold text-on-surface mb-2 block">命名空间</label>
                  <input
                    type="text"
                    value={server.client_config?.namespace || ''}
                    onChange={(e) => handleClientConfigChange('namespace', e.target.value)}
                    className="w-full bg-surface-container-highest/50 border border-outline-variant/10 rounded-xl px-4 py-3 text-on-surface font-mono"
                    placeholder="例如: weather, fs"
                  />
                </div>

                <div>
                  <label className="text-sm font-bold text-on-surface mb-2 block">超时时间 (秒)</label>
                  <input
                    type="number"
                    value={server.client_config?.timeout || 30}
                    onChange={(e) => handleClientConfigChange('timeout', parseInt(e.target.value) || 30)}
                    className="w-full bg-surface-container-highest/50 border border-outline-variant/10 rounded-xl px-4 py-3 text-on-surface"
                    min="1"
                    max="300"
                  />
                </div>
              </div>

              {/* 参数列表 */}
              <div className="mt-4">
                <div className="flex items-center justify-between mb-2">
                  <label className="text-sm font-bold text-on-surface">参数</label>
                  <button
                    onClick={addArg}
                    className="text-xs font-bold uppercase tracking-widest bg-surface-container-highest text-on-surface hover:bg-surface-bright px-3 py-1.5 rounded-lg"
                  >
                    添加参数
                  </button>
                </div>
                
                <div className="space-y-2">
                  {(server.client_config?.args || []).map((arg, index) => (
                    <div key={index} className="flex gap-2">
                      <input
                        type="text"
                        value={arg}
                        onChange={(e) => handleArgsChange(index, e.target.value)}
                        className="flex-1 bg-surface-container-highest/50 border border-outline-variant/10 rounded-xl px-4 py-2 text-on-surface font-mono"
                        placeholder={`参数 ${index + 1}`}
                      />
                      <button
                        onClick={() => removeArg(index)}
                        className="bg-error-container text-on-error-container hover:brightness-110 rounded-xl px-4 py-2"
                      >
                        <X size={16} />
                      </button>
                    </div>
                  ))}
                  
                  {(server.client_config?.args || []).length === 0 && (
                    <div className="text-outline text-sm text-center py-4">
                      暂无参数，点击"添加参数"按钮添加
                    </div>
                  )}
                </div>
              </div>
            </div>
          )}

          {/* 服务器配置 */}
          {server.type === 'server' && (
            <div className="border-t border-outline-variant/10 pt-6">
              <h4 className="text-lg font-bold text-primary mb-4">服务器配置</h4>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-bold text-on-surface mb-2 block">服务器类型</label>
                  <input
                    type="text"
                    value={server.server_config?.type || ''}
                    onChange={(e) => onChange({ 
                      server_config: { ...server.server_config, type: e.target.value }
                    })}
                    className="w-full bg-surface-container-highest/50 border border-outline-variant/10 rounded-xl px-4 py-3 text-on-surface font-mono"
                    placeholder="例如: stdio"
                  />
                </div>

                <div className="md:col-span-2">
                  <label className="text-sm font-bold text-on-surface mb-2 block">工具过滤器 (正则表达式)</label>
                  <input
                    type="text"
                    value={server.server_config?.tool_filter || ''}
                    onChange={(e) => onChange({ 
                      server_config: { ...server.server_config, tool_filter: e.target.value }
                    })}
                    className="w-full bg-surface-container-highest/50 border border-outline-variant/10 rounded-xl px-4 py-3 text-on-surface font-mono"
                    placeholder="例如: ^fs\\. 匹配以'fs.'开头的工具"
                  />
                </div>
              </div>
            </div>
          )}

          {/* 操作按钮 */}
          <div className="flex justify-end gap-4 pt-6 border-t border-outline-variant/10">
            <button
              onClick={onClose}
              className="bg-surface-container-highest text-on-surface hover:bg-surface-bright rounded-xl px-6 py-3 text-sm font-bold transition-all"
            >
              取消
            </button>
            <button
              onClick={onSave}
              className="bg-primary-container text-on-primary-container hover:brightness-110 rounded-xl px-8 py-3 text-sm font-bold shadow-xl shadow-primary/20 transition-all flex items-center gap-2"
            >
              <Save size={18} />
              {isNew ? '添加服务器' : '保存更改'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}