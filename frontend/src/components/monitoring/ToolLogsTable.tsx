import React, { useState, useMemo } from 'react';
import { 
  Settings, 
  Filter, 
  Search, 
  Download, 
  RefreshCw, 
  ChevronUp, 
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  MoreVertical,
  Copy,
  Eye,
  EyeOff,
  Clock,
  Zap,
  Cpu,
  Database,
  Terminal,
  FileText,
  GitBranch,
  Layout,
  CheckCircle2,
  AlertCircle,
  TrendingUp,
  BarChart3,
  Activity,
  Server,
  Users,
  Shield,
  Brain,
  MessageSquare,
  Code,
  Play,
  Pause,
  SkipBack,
  SkipForward,
  Volume2,
  VolumeX,
  Maximize2,
  Minimize2
} from 'lucide-react';
import { cn } from '../../utils';
import { ToolInvocation } from '../../types';

interface ToolLogsTableProps {
  logs: ToolInvocation[];
}

export default function ToolLogsTable({ logs }: ToolLogsTableProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<'all' | 'success' | 'error' | 'pending'>('all');
  const [toolTypeFilter, setToolTypeFilter] = useState<string>('all');
  const [sortColumn, setSortColumn] = useState<'timestamp' | 'toolName' | 'duration' | 'status'>('timestamp');
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [selectedLog, setSelectedLog] = useState<string | null>(null);
  const [showParams, setShowParams] = useState(false);
  const [showResults, setShowResults] = useState(false);

  // 获取所有工具类型
  const toolTypes = useMemo(() => {
    const types = new Set(logs.map(log => log.toolType || 'unknown'));
    return Array.from(types);
  }, [logs]);

  // 过滤和排序日志
  const filteredLogs = useMemo(() => {
    let filtered = logs.filter(log => {
      // 搜索过滤
      if (searchQuery && !log.toolName.toLowerCase().includes(searchQuery.toLowerCase())) {
        return false;
      }
      
      // 状态过滤
      if (statusFilter !== 'all' && log.status !== statusFilter) {
        return false;
      }
      
      // 工具类型过滤
      if (toolTypeFilter !== 'all' && log.toolType !== toolTypeFilter) {
        return false;
      }
      
      return true;
    });
    
    // 排序
    filtered.sort((a, b) => {
      let aValue: any, bValue: any;
      
      switch (sortColumn) {
        case 'timestamp':
          aValue = new Date(a.timestamp).getTime();
          bValue = new Date(b.timestamp).getTime();
          break;
        case 'toolName':
          aValue = a.toolName.toLowerCase();
          bValue = b.toolName.toLowerCase();
          break;
        case 'duration':
          aValue = a.duration || 0;
          bValue = b.duration || 0;
          break;
        case 'status':
          aValue = a.status;
          bValue = b.status;
          break;
        default:
          return 0;
      }
      
      if (sortDirection === 'asc') {
        return aValue > bValue ? 1 : -1;
      } else {
        return aValue < bValue ? 1 : -1;
      }
    });
    
    return filtered;
  }, [logs, searchQuery, statusFilter, toolTypeFilter, sortColumn, sortDirection]);

  // 分页
  const totalPages = Math.ceil(filteredLogs.length / pageSize);
  const paginatedLogs = filteredLogs.slice((page - 1) * pageSize, page * pageSize);

  const getToolIcon = (toolType: string) => {
    switch (toolType) {
      case 'fs': return FileText;
      case 'git': return GitBranch;
      case 'shell': return Terminal;
      case 'project': return Layout;
      case 'skill': return Brain;
      case 'memory': return Database;
      case 'api': return Server;
      case 'auth': return Shield;
      default: return Settings;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'success': return 'text-green-500';
      case 'error': return 'text-red-500';
      case 'pending': return 'text-yellow-500';
      case 'running': return 'text-primary';
      default: return 'text-white/60';
    }
  };

  const getStatusBgColor = (status: string) => {
    switch (status) {
      case 'success': return 'bg-green-500/20';
      case 'error': return 'bg-red-500/20';
      case 'pending': return 'bg-yellow-500/20';
      case 'running': return 'bg-primary/20';
      default: return 'bg-white/5';
    }
  };

  const getDurationColor = (duration: number) => {
    if (duration < 100) return 'text-green-500';
    if (duration < 500) return 'text-yellow-500';
    if (duration < 1000) return 'text-orange-500';
    return 'text-red-500';
  };

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString('en-US', { 
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  };

  const formatDuration = (duration: number) => {
    if (duration < 1000) return `${duration}ms`;
    return `${(duration / 1000).toFixed(2)}s`;
  };

  const handleSort = (column: typeof sortColumn) => {
    if (sortColumn === column) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
    } else {
      setSortColumn(column);
      setSortDirection('desc');
    }
  };

  const handleDownload = () => {
    const data = {
      logs: filteredLogs,
      timestamp: new Date().toISOString(),
      filters: { searchQuery, statusFilter, toolTypeFilter },
      sort: { column: sortColumn, direction: sortDirection }
    };
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `tool-logs-${new Date().toISOString().split('T')[0]}.json`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const selectedLogData = selectedLog ? logs.find(log => log.id === selectedLog) : null;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-3">
            <Settings size={24} />
            Tool Invocation Logs
          </h2>
          <p className="text-sm text-white/60 mt-1">
            Detailed records of all tool calls with parameters and results
          </p>
        </div>
        
        <div className="flex items-center gap-3">
          <div className="text-sm text-white/40">
            {filteredLogs.length} logs • {toolTypes.length} tool types
          </div>
        </div>
      </div>

      {/* Filters and Controls */}
      <div className="bg-white/5 border border-white/10 rounded-xl p-4">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="relative">
              <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-white/40" />
              <input
                type="text"
                placeholder="Search tool names..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 pr-4 py-2 bg-white/5 border border-white/10 rounded-lg text-sm focus:outline-none focus:border-primary transition-colors"
              />
            </div>
            
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value as any)}
              className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
            >
              <option value="all">All Status</option>
              <option value="success">Success</option>
              <option value="error">Error</option>
              <option value="pending">Pending</option>
              <option value="running">Running</option>
            </select>
            
            <select
              value={toolTypeFilter}
              onChange={(e) => setToolTypeFilter(e.target.value)}
              className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
            >
              <option value="all">All Tool Types</option>
              {toolTypes.map(type => (
                <option key={type} value={type}>{type}</option>
              ))}
            </select>
          </div>
          
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2">
              <span className="text-sm text-white/60">Show:</span>
              <button
                onClick={() => setShowParams(!showParams)}
                className={cn(
                  "px-3 py-1 text-xs rounded transition-colors",
                  showParams ? "bg-primary/20 text-primary" : "bg-white/5 text-white/60 hover:bg-white/10"
                )}
              >
                Params
              </button>
              <button
                onClick={() => setShowResults(!showResults)}
                className={cn(
                  "px-3 py-1 text-xs rounded transition-colors",
                  showResults ? "bg-primary/20 text-primary" : "bg-white/5 text-white/60 hover:bg-white/10"
                )}
              >
                Results
              </button>
            </div>
            
            <button
              onClick={handleDownload}
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
            >
              <Download size={20} />
            </button>
          </div>
        </div>
        
        {/* Statistics */}
        <div className="grid grid-cols-4 gap-4 mt-4">
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Total Calls</div>
            <div className="text-lg font-bold">{logs.length}</div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Success Rate</div>
            <div className="text-lg font-bold text-green-500">
              {logs.length > 0 
                ? `${Math.round((logs.filter(l => l.status === 'success').length / logs.length) * 100)}%`
                : '0%'}
            </div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Avg Duration</div>
            <div className="text-lg font-bold">
              {logs.length > 0 
                ? formatDuration(logs.reduce((sum, log) => sum + (log.duration || 0), 0) / logs.length)
                : '0ms'}
            </div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Errors</div>
            <div className="text-lg font-bold text-red-500">
              {logs.filter(l => l.status === 'error').length}
            </div>
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="border border-white/10 rounded-xl overflow-hidden bg-white/[0.02]">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-white/10">
                <th className="text-left p-4">
                  <button
                    onClick={() => handleSort('timestamp')}
                    className="flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-white/40 hover:text-white transition-colors"
                  >
                    Time
                    {sortColumn === 'timestamp' && (
                      sortDirection === 'asc' ? <ChevronUp size={12} /> : <ChevronDown size={12} />
                    )}
                  </button>
                </th>
                <th className="text-left p-4">
                  <button
                    onClick={() => handleSort('toolName')}
                    className="flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-white/40 hover:text-white transition-colors"
                  >
                    Tool
                    {sortColumn === 'toolName' && (
                      sortDirection === 'asc' ? <ChevronUp size={12} /> : <ChevronDown size={12} />
                    )}
                  </button>
                </th>
                <th className="text-left p-4">
                  <button
                    onClick={() => handleSort('duration')}
                    className="flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-white/40 hover:text-white transition-colors"
                  >
                    Duration
                    {sortColumn === 'duration' && (
                      sortDirection === 'asc' ? <ChevronUp size={12} /> : <ChevronDown size={12} />
                    )}
                  </button>
                </th>
                <th className="text-left p-4">
                  <button
                    onClick={() => handleSort('status')}
                    className="flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-white/40 hover:text-white transition-colors"
                  >
                    Status
                    {sortColumn === 'status' && (
                      sortDirection === 'asc' ? <ChevronUp size={12} /> : <ChevronDown size={12} />
                    )}
                  </button>
                </th>
                {showParams && (
                  <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                    Parameters
                  </th>
                )}
                {showResults && (
                  <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                    Result
                  </th>
                )}
                <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody>
              {paginatedLogs.map((log) => {
                const Icon = getToolIcon(log.toolType || 'unknown');
                
                return (
                  <tr 
                    key={log.id}
                    className={cn(
                      "border-b border-white/5 hover:bg-white/5 transition-colors",
                      selectedLog === log.id && "bg-primary/10"
                    )}
                    onClick={() => setSelectedLog(log.id === selectedLog ? null : log.id)}
                  >
                    <td className="p-4">
                      <div className="text-sm font-mono">{formatTimestamp(log.timestamp)}</div>
                      <div className="text-xs text-white/40">{log.agent || 'Unknown'}</div>
                    </td>
                    <td className="p-4">
                      <div className="flex items-center gap-3">
                        <div className={cn(
                          "p-2 rounded-lg",
                          getStatusBgColor(log.status)
                        )}>
                          <Icon size={16} className={getStatusColor(log.status)} />
                        </div>
                        <div>
                          <div className="font-bold">{log.toolName}</div>
                          <div className="text-xs text-white/40">{log.toolType || 'unknown'}</div>
                        </div>
                      </div>
                    </td>
                    <td className="p-4">
                      <div className={cn(
                        "text-sm font-bold",
                        getDurationColor(log.duration || 0)
                      )}>
                        {formatDuration(log.duration || 0)}
                      </div>
                    </td>
                    <td className="p-4">
                      <div className={cn(
                        "inline-flex items-center px-3 py-1 rounded-full text-xs font-bold",
                        getStatusBgColor(log.status),
                        getStatusColor(log.status)
                      )}>
                        {log.status.toUpperCase()}
                      </div>
                    </td>
                    {showParams && (
                      <td className="p-4">
                        <div className="text-sm max-w-xs truncate">
                          {log.parameters ? JSON.stringify(log.parameters) : 'None'}
                        </div>
                      </td>
                    )}
                    {showResults && (
                      <td className="p-4">
                        <div className="text-sm max-w-xs truncate">
                          {log.result ? JSON.stringify(log.result) : 'None'}
                        </div>
                      </td>
                    )}
                    <td className="p-4">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          // 复制日志数据
                          navigator.clipboard.writeText(JSON.stringify(log, null, 2));
                        }}
                        className="p-1 rounded hover:bg-white/10 transition-colors"
                        title="Copy log data"
                      >
                        <Copy size={14} />
                      </button>
                    </td>
                  </tr>
                );
              })}
              
              {paginatedLogs.length === 0 && (
                <tr>
                  <td colSpan={7} className="p-8 text-center">
                    <Settings size={48} className="mx-auto text-white/20 mb-4" />
                    <div className="text-lg font-bold text-white/40">No tool logs found</div>
                    <p className="text-white/60 mt-2">
                      {searchQuery 
                        ? `No logs match "${searchQuery}"`
                        : "Try changing your filters or wait for tools to be invoked"}
                    </p>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
        
        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between p-4 border-t border-white/10">
            <div className="text-sm text-white/60">
              Showing {(page - 1) * pageSize + 1} to {Math.min(page * pageSize, filteredLogs.length)} of {filteredLogs.length} logs
            </div>
            
            <div className="flex items-center gap-2">
              <button
                onClick={() => setPage(1)}
                disabled={page === 1}
                className={cn(
                  "p-2 rounded-lg transition-colors",
                  page === 1
                    ? "bg-white/5 text-white/20"
                    : "bg-white/5 text-white/60 hover:bg-white/10"
                )}
              >
                <ChevronLeft size={16} className="-ml-1" />
                <ChevronLeft size={16} className="-ml-3" />
              </button>
              
              <button
                onClick={() => setPage(Math.max(1, page - 1))}
                disabled={page === 1}
                className={cn(
                  "p-2 rounded-lg transition-colors",
                  page === 1
                    ? "bg-white/5 text-white/20"
                    : "bg-white/5 text-white/60 hover:bg-white/10"
                )}
              >
                <ChevronLeft size={16} />
              </button>
              
              <div className="flex items-center gap-1">
                {Array.from({ length: Math.min(5, totalPages) }).map((_, i) => {
                  let pageNum: number;
                  if (totalPages <= 5) {
                    pageNum = i + 1;
                  } else if (page <= 3) {
                    pageNum = i + 1;
                  } else if (page >= totalPages - 2) {
                    pageNum = totalPages - 4 + i;
                  } else {
                    pageNum = page - 2 + i;
                  }
                  
                  return (
                    <button
                      key={pageNum}
                      onClick={() => setPage(pageNum)}
                      className={cn(
                        "w-8 h-8 rounded-lg text-sm transition-colors",
                        page === pageNum
                          ? "bg-primary text-black font-bold"
                          : "bg-white/5 text-white/60 hover:bg-white/10"
                      )}
                    >
                      {pageNum}
                    </button>
                  );
                })}
              </div>
              
              <button
                onClick={() => setPage(Math.min(totalPages, page + 1))}
                disabled={page === totalPages}
                className={cn(
                  "p-2 rounded-lg transition-colors",
                  page === totalPages
                    ? "bg-white/5 text-white/20"
                    : "bg-white/5 text-white/60 hover:bg-white/10"
                )}
              >
                <ChevronRight size={16} />
              </button>
              
              <button
                onClick={() => setPage(totalPages)}
                disabled={page === totalPages}
                className={cn(
                  "p-2 rounded-lg transition-colors",
                  page === totalPages
                    ? "bg-white/5 text-white/20"
                    : "bg-white/5 text-white/60 hover:bg-white/10"
                )}
              >
                <ChevronRight size={16} className="-ml-1" />
                <ChevronRight size={16} className="-ml-3" />
              </button>
            </div>
            
            <div className="flex items-center gap-2">
              <span className="text-sm text-white/60">Show:</span>
              <select
                value={pageSize}
                onChange={(e) => {
                  setPageSize(parseInt(e.target.value));
                  setPage(1);
                }}
                className="bg-white/5 border border-white/10 rounded-lg px-3 py-1 text-sm focus:outline-none focus:border-primary transition-colors"
              >
                <option value="10">10</option>
                <option value="20">20</option>
                <option value="50">50</option>
                <option value="100">100</option>
              </select>
            </div>
          </div>
        )}
      </div>

      {/* Log Details Panel */}
      {selectedLogData && (
        <div className="border border-white/10 rounded-xl p-6 bg-white/[0.02]">
          <div className="flex justify-between items-start mb-6">
            <div>
              <h3 className="text-lg font-bold">Tool Invocation Details</h3>
              <p className="text-sm text-white/60 mt-1">
                {selectedLogData.toolName} • {formatTimestamp(selectedLogData.timestamp)}
              </p>
            </div>
            
            <button
              onClick={() => setSelectedLog(null)}
              className="p-2 rounded-lg hover:bg-white/10 transition-colors"
            >
              <ChevronRight size={20} className="rotate-90" />
            </button>
          </div>
          
          <div className="grid grid-cols-2 gap-6">
            <div className="space-y-4">
              <div>
                <div className="text-sm font-bold mb-2">Basic Information</div>
                <div className="bg-white/5 border border-white/10 rounded-lg p-4 space-y-3">
                  <div className="flex justify-between">
                    <span className="text-white/60">Tool Name:</span>
                    <span className="font-bold">{selectedLogData.toolName}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-white/60">Tool Type:</span>
                    <span className="font-bold">{selectedLogData.toolType || 'unknown'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-white/60">Agent:</span>
                    <span className="font-bold">{selectedLogData.agent || 'Unknown'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-white/60">Timestamp:</span>
                    <span className="font-mono">{formatTimestamp(selectedLogData.timestamp)}</span>
                  </div>
                </div>
              </div>
              
              <div>
                <div className="text-sm font-bold mb-2">Performance</div>
                <div className="bg-white/5 border border-white/10 rounded-lg p-4 space-y-3">
                  <div className="flex justify-between">
                    <span className="text-white/60">Duration:</span>
                    <span className={cn(
                      "font-bold",
                      getDurationColor(selectedLogData.duration || 0)
                    )}>
                      {formatDuration(selectedLogData.duration || 0)}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-white/60">Status:</span>
                    <span className={cn(
                      "font-bold",
                      getStatusColor(selectedLogData.status)
                    )}>
                      {selectedLogData.status.toUpperCase()}
                    </span>
                  </div>
                  {selectedLogData.error && (
                    <div className="flex justify-between">
                      <span className="text-white/60">Error:</span>
                      <span className="font-bold text-red-500">{selectedLogData.error}</span>
                    </div>
                  )}
                </div>
              </div>
            </div>
            
            <div className="space-y-4">
              <div>
                <div className="text-sm font-bold mb-2">Parameters</div>
                <div className="bg-black/40 border border-white/10 rounded-lg p-4 max-h-60 overflow-y-auto">
                  <pre className="text-sm text-white/80">
                    {selectedLogData.parameters 
                      ? JSON.stringify(selectedLogData.parameters, null, 2)
                      : 'No parameters'}
                  </pre>
                </div>
              </div>
              
              <div>
                <div className="text-sm font-bold mb-2">Result</div>
                <div className="bg-black/40 border border-white/10 rounded-lg p-4 max-h-60 overflow-y-auto">
                  <pre className="text-sm text-white/80">
                    {selectedLogData.result 
                      ? JSON.stringify(selectedLogData.result, null, 2)
                      : 'No result'}
                  </pre>
                </div>
              </div>
            </div>
          </div>
          
          {selectedLogData.metadata && (
            <div className="mt-6">
              <div className="text-sm font-bold mb-2">Metadata</div>
              <div className="bg-white/5 border border-white/10 rounded-lg p-4">
                <pre className="text-sm text-white/80 overflow-x-auto">
                  {JSON.stringify(selectedLogData.metadata, null, 2)}
                </pre>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Tool Type Distribution */}
      <div className="border border-white/10 rounded-xl p-6">
        <h4 className="text-lg font-bold mb-4">Tool Type Distribution</h4>
        <div className="space-y-4">
          {toolTypes.map((type) => {
            const count = logs.filter(l => l.toolType === type).length;
            const percentage = (count / logs.length) * 100;
            const Icon = getToolIcon(type);
            
            return (
              <div key={type} className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-white/5 rounded-lg">
                    <Icon size={16} className="text-white/60" />
                  </div>
                  <span className="text-sm capitalize">{type}</span>
                </div>
                <div className="flex items-center gap-4">
                  <div className="w-48 h-2 bg-white/10 rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-primary rounded-full"
                      style={{ width: `${percentage}%` }}
                    />
                  </div>
                  <div className="text-right w-16">
                    <span className="text-sm font-bold">{count}</span>
                    <span className="text-xs text-white/40 ml-1">({percentage.toFixed(1)}%)</span>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}