import React, { useState, useEffect, useCallback } from 'react';
import { 
  Clock, 
  Play, 
  Pause, 
  StopCircle, 
  Trash2, 
  Filter, 
  Search, 
  ChevronRight, 
  ChevronLeft,
  Calendar,
  User,
  Cpu,
  Database,
  Terminal,
  FileText,
  AlertCircle,
  CheckCircle2,
  XCircle,
  MoreVertical,
  Download,
  Eye,
  RefreshCw,
  BarChart3,
  Tag,
  Hash,
  Zap,
  DollarSign,
  Timer
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../../utils';
import { api } from '../../services/api';

interface ExecutionRecord {
  id: string;
  trace_id: string;
  pipeline_id: string;
  name: string;
  description?: string;
  status: 'running' | 'completed' | 'failed' | 'cancelled';
  start_time: string;
  end_time?: string;
  duration?: number;
  agent_count: number;
  tool_call_count: number;
  token_usage: number;
  cost_estimate?: number;
  metadata?: Record<string, any>;
  tags?: string[];
  created_at: string;
  updated_at: string;
}

interface HistoryQuery {
  searchText?: string;
  statuses?: string[];
  tags?: string[];
  startTime?: string;
  endTime?: string;
  sortBy?: string;
  sortOrder?: string;
  limit?: number;
  offset?: number;
}

export default function HistoryList() {
  const [records, setRecords] = useState<ExecutionRecord[]>([]);
  const [totalRecords, setTotalRecords] = useState(0);
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedStatuses, setSelectedStatuses] = useState<string[]>([]);
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [sortBy, setSortBy] = useState('start_time');
  const [sortOrder, setSortOrder] = useState('desc');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20);
  const [stats, setStats] = useState<any>(null);
  const [activeReplay, setActiveReplay] = useState<string | null>(null);
  const [replayStates, setReplayStates] = useState<Record<string, any>>({});

  const statusOptions = [
    { value: 'running', label: 'Running', color: 'bg-blue-500', icon: <Clock size={12} /> },
    { value: 'completed', label: 'Completed', color: 'bg-green-500', icon: <CheckCircle2 size={12} /> },
    { value: 'failed', label: 'Failed', color: 'bg-red-500', icon: <XCircle size={12} /> },
    { value: 'cancelled', label: 'Cancelled', color: 'bg-gray-500', icon: <StopCircle size={12} /> }
  ];

  const tagOptions = [
    'production', 'test', 'demo', 'experiment', 'optimization', 'debug', 'feature', 'bugfix'
  ];

  const fetchRecords = useCallback(async () => {
    setLoading(true);
    try {
      const query: HistoryQuery = {
        searchText: searchQuery || undefined,
        statuses: selectedStatuses.length > 0 ? selectedStatuses : undefined,
        tags: selectedTags.length > 0 ? selectedTags : undefined,
        sortBy,
        sortOrder,
        limit: pageSize,
        offset: (currentPage - 1) * pageSize
      };

      const [recordsData, total] = await api.getHistoryRecords(query);
      setRecords(recordsData);
      setTotalRecords(total);
    } catch (error) {
      console.error('Failed to fetch history records:', error);
    } finally {
      setLoading(false);
    }
  }, [searchQuery, selectedStatuses, selectedTags, sortBy, sortOrder, currentPage, pageSize]);

  const fetchStatistics = useCallback(async () => {
    try {
      const statsData = await api.getHistoryStatistics();
      setStats(statsData);
    } catch (error) {
      console.error('Failed to fetch history statistics:', error);
    }
  }, []);

  const handleStartReplay = async (recordId: string) => {
    try {
      const state = await api.startReplay(recordId, 1.0);
      setActiveReplay(recordId);
      setReplayStates(prev => ({ ...prev, [recordId]: state }));
      
      // 订阅回放事件
      api.subscribeToReplayEvents(recordId, (event, data) => {
        if (event === 'replay:progress' || event === 'replay:ended') {
          setReplayStates(prev => ({ ...prev, [recordId]: data.state }));
        }
        if (event === 'replay:ended') {
          setActiveReplay(null);
        }
      });
    } catch (error) {
      console.error('Failed to start replay:', error);
    }
  };

  const handleControlReplay = async (replayId: string, action: string, params?: any) => {
    try {
      const state = await api.controlReplay(replayId, action, params);
      setReplayStates(prev => ({ ...prev, [replayId]: state }));
    } catch (error) {
      console.error('Failed to control replay:', error);
    }
  };

  const handleDeleteRecord = async (recordId: string) => {
    if (!confirm('Are you sure you want to delete this execution record?')) {
      return;
    }

    try {
      await api.deleteHistoryRecord(recordId);
      await fetchRecords();
    } catch (error) {
      console.error('Failed to delete record:', error);
    }
  };

  const handleExportRecord = async (recordId: string) => {
    try {
      const record = records.find(r => r.id === recordId);
      if (!record) return;

      const dataStr = JSON.stringify(record, null, 2);
      const dataBlob = new Blob([dataStr], { type: 'application/json' });
      const url = URL.createObjectURL(dataBlob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `execution-${recordId}.json`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Failed to export record:', error);
    }
  };

  const formatDuration = (ms?: number) => {
    if (!ms) return 'N/A';
    const seconds = Math.floor(ms / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);

    if (hours > 0) {
      return `${hours}h ${minutes % 60}m`;
    } else if (minutes > 0) {
      return `${minutes}m ${seconds % 60}s`;
    } else {
      return `${seconds}s`;
    }
  };

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'running': return 'bg-blue-500/20 text-blue-400 border-blue-500/30';
      case 'completed': return 'bg-green-500/20 text-green-400 border-green-500/30';
      case 'failed': return 'bg-red-500/20 text-red-400 border-red-500/30';
      case 'cancelled': return 'bg-gray-500/20 text-gray-400 border-gray-500/30';
      default: return 'bg-gray-500/20 text-gray-400 border-gray-500/30';
    }
  };

  useEffect(() => {
    fetchRecords();
    fetchStatistics();
  }, [fetchRecords, fetchStatistics]);

  const totalPages = Math.ceil(totalRecords / pageSize);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold">Execution History</h2>
          <p className="text-sm text-white/60">
            View and replay past AI agent executions
          </p>
        </div>
        
        <div className="flex items-center gap-4">
          <button
            onClick={fetchRecords}
            disabled={loading}
            className="px-4 py-2 text-sm font-medium bg-white/5 text-white/60 border border-white/10 rounded-lg hover:bg-white/10 transition-colors flex items-center gap-2"
          >
            <RefreshCw size={16} className={cn(loading && "animate-spin")} />
            Refresh
          </button>
          
          <button
            onClick={() => {
              setSearchQuery('');
              setSelectedStatuses([]);
              setSelectedTags([]);
              setCurrentPage(1);
            }}
            className="px-4 py-2 text-sm font-medium bg-white/5 text-white/60 border border-white/10 rounded-lg hover:bg-white/10 transition-colors flex items-center gap-2"
          >
            <Filter size={16} />
            Clear Filters
          </button>
        </div>
      </div>

      {/* Statistics */}
      {stats && (
        <div className="grid grid-cols-4 gap-4">
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="flex items-center gap-3 mb-2">
              <div className="p-2 bg-blue-500/20 rounded-lg">
                <BarChart3 size={20} className="text-blue-400" />
              </div>
              <div>
                <div className="text-xs text-white/40 uppercase tracking-wider">Total Executions</div>
                <div className="text-2xl font-bold">{stats.total_records || 0}</div>
              </div>
            </div>
            <div className="text-xs text-white/60">Success Rate: {stats.success_rate?.toFixed(1) || 0}%</div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="flex items-center gap-3 mb-2">
              <div className="p-2 bg-green-500/20 rounded-lg">
                <Timer size={20} className="text-green-400" />
              </div>
              <div>
                <div className="text-xs text-white/40 uppercase tracking-wider">Avg Duration</div>
                <div className="text-2xl font-bold">
                  {stats.avg_duration ? formatDuration(stats.avg_duration) : 'N/A'}
                </div>
              </div>
            </div>
            <div className="text-xs text-white/60">Total: {stats.total_duration || '0s'}</div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="flex items-center gap-3 mb-2">
              <div className="p-2 bg-purple-500/20 rounded-lg">
                <Zap size={20} className="text-purple-400" />
              </div>
              <div>
                <div className="text-xs text-white/40 uppercase tracking-wider">Total Tokens</div>
                <div className="text-2xl font-bold">
                  {(stats.total_token_usage || 0).toLocaleString()}
                </div>
              </div>
            </div>
            <div className="text-xs text-white/60">Avg: {(stats.avg_token_usage || 0).toLocaleString()}</div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="flex items-center gap-3 mb-2">
              <div className="p-2 bg-yellow-500/20 rounded-lg">
                <DollarSign size={20} className="text-yellow-400" />
              </div>
              <div>
                <div className="text-xs text-white/40 uppercase tracking-wider">Total Cost</div>
                <div className="text-2xl font-bold">
                  ${(stats.total_cost || 0).toFixed(2)}
                </div>
              </div>
            </div>
            <div className="text-xs text-white/60">Avg: ${(stats.avg_cost || 0).toFixed(2)}</div>
          </div>
        </div>
      )}

      {/* Filters */}
      <div className="bg-white/5 border border-white/10 rounded-lg p-4">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {/* Search */}
          <div>
            <label className="block text-sm font-medium mb-2">Search</label>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-white/40" size={16} />
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search by name or description..."
                className="w-full bg-white/5 border border-white/10 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
              />
            </div>
          </div>

          {/* Status Filter */}
          <div>
            <label className="block text-sm font-medium mb-2">Status</label>
            <div className="flex flex-wrap gap-2">
              {statusOptions.map((status) => (
                <button
                  key={status.value}
                  onClick={() => {
                    setSelectedStatuses(prev =>
                      prev.includes(status.value)
                        ? prev.filter(s => s !== status.value)
                        : [...prev, status.value]
                    );
                    setCurrentPage(1);
                  }}
                  className={cn(
                    "px-3 py-1 text-xs rounded-full border transition-colors flex items-center gap-1",
                    selectedStatuses.includes(status.value)
                      ? "bg-primary/20 text-primary border-primary/30"
                      : "bg-white/5 text-white/60 border-white/10 hover:bg-white/10"
                  )}
                >
                  <div className={cn("w-2 h-2 rounded-full", status.color)} />
                  {status.label}
                </button>
              ))}
            </div>
          </div>

          {/* Tags Filter */}
          <div>
            <label className="block text-sm font-medium mb-2">Tags</label>
            <div className="flex flex-wrap gap-2">
              {tagOptions.map((tag) => (
                <button
                  key={tag}
                  onClick={() => {
                    setSelectedTags(prev =>
                      prev.includes(tag)
                        ? prev.filter(t => t !== tag)
                        : [...prev, tag]
                    );
                    setCurrentPage(1);
                  }}
                  className={cn(
                    "px-3 py-1 text-xs rounded-full border transition-colors flex items-center gap-1",
                    selectedTags.includes(tag)
                      ? "bg-primary/20 text-primary border-primary/30"
                      : "bg-white/5 text-white/60 border-white/10 hover:bg-white/10"
                  )}
                >
                  <Tag size={12} />
                  {tag}
                </button>
              ))}
            </div>
          </div>

          {/* Sort */}
          <div>
            <label className="block text-sm font-medium mb-2">Sort By</label>
            <div className="flex gap-2">
              <select
                value={sortBy}
                onChange={(e) => {
                  setSortBy(e.target.value);
                  setCurrentPage(1);
                }}
                className="flex-1 bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
              >
                <option value="start_time">Start Time</option>
                <option value="end_time">End Time</option>
                <option value="duration">Duration</option>
                <option value="agent_count">Agent Count</option>
                <option value="token_usage">Token Usage</option>
                <option value="cost_estimate">Cost</option>
              </select>
              <button
                onClick={() => {
                  setSortOrder(sortOrder === 'desc' ? 'asc' : 'desc');
                  setCurrentPage(1);
                }}
                className="px-3 py-2 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
              >
                {sortOrder === 'desc' ? '↓' : '↑'}
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Records Table */}
      <div className="bg-white/5 border border-white/10 rounded-lg overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-white/10">
                <th className="text-left p-4 text-sm font-medium text-white/60">Name</th>
                <th className="text-left p-4 text-sm font-medium text-white/60">Status</th>
                <th className="text-left p-4 text-sm font-medium text-white/60">Time</th>
                <th className="text-left p-4 text-sm font-medium text-white/60">Agents</th>
                <th className="text-left p-4 text-sm font-medium text-white/60">Tools</th>
                <th className="text-left p-4 text-sm font-medium text-white/60">Tokens</th>
                <th className="text-left p-4 text-sm font-medium text-white/60">Actions</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={7} className="p-8 text-center text-white/60">
                    <div className="flex justify-center">
                      <RefreshCw size={24} className="animate-spin" />
                    </div>
                  </td>
                </tr>
              ) : records.length === 0 ? (
                <tr>
                  <td colSpan={7} className="p-8 text-center text-white/60">
                    No execution records found
                  </td>
                </tr>
              ) : (
                records.map((record) => {
                  const replayState = replayStates[record.id];
                  const isReplaying = activeReplay === record.id;
                  
                  return (
                    <tr key={record.id} className="border-b border-white/10 hover:bg-white/5">
                      <td className="p-4">
                        <div>
                          <div className="font-medium">{record.name}</div>
                          {record.description && (
                            <div className="text-sm text-white/60 mt-1">{record.description}</div>
                          )}
                          <div className="flex flex-wrap gap-1 mt-2">
                            {record.tags?.map((tag) => (
                              <span
                                key={tag}
                                className="px-2 py-1 text-xs bg-white/5 text-white/60 rounded"
                              >
                                {tag}
                              </span>
                            ))}
                          </div>
                        </div>
                      </td>
                      <td className="p-4">
                        <span className={cn(
                          "px-3 py-1 text-xs rounded-full border",
                          getStatusColor(record.status)
                        )}>
                          {record.status.charAt(0).toUpperCase() + record.status.slice(1)}
                        </span>
                        {isReplaying && replayState && (
                          <div className="mt-2">
                            <div className="text-xs text-white/60 mb-1">
                              Replay: {replayState.progress?.toFixed(1)}%
                            </div>
                            <div className="w-full bg-white/10 rounded-full h-1">
                              <div 
                                className="bg-primary h-1 rounded-full transition-all"
                                style={{ width: `${replayState.progress || 0}%` }}
                              />
                            </div>
                          </div>
                        )}
                      </td>
                      <td className="p-4">
                        <div className="text-sm">{formatDateTime(record.start_time)}</div>
                        <div className="text-xs text-white/60 mt-1">
                          Duration: {formatDuration(record.duration)}
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="flex items-center gap-2">
                          <User size={14} className="text-white/40" />
                          <span>{record.agent_count}</span>
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="flex items-center gap-2">
                          <Terminal size={14} className="text-white/40" />
                          <span>{record.tool_call_count}</span>
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="flex items-center gap-2">
                          <Zap size={14} className="text-white/40" />
                          <span>{record.token_usage.toLocaleString()}</span>
                        </div>
                        {record.cost_estimate && (
                          <div className="text-xs text-white/60 mt-1">
                            ${record.cost_estimate.toFixed(2)}
                          </div>
                        )}
                      </td>
                      <td className="p-4">
                        <div className="flex items-center gap-2">
                          {isReplaying ? (
                            <>
                              <button
                                onClick={() => handleControlReplay(record.id, 'pause')}
                                className="p-2 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
                                title="Pause"
                              >
                                <Pause size={16} />
                              </button>
                              <button
                                onClick={() => handleControlReplay(record.id, 'stop')}
                                className="p-2 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
                                title="Stop"
                              >
                                <StopCircle size={16} />
                              </button>
                            </>
                          ) : (
                            <button
                              onClick={() => handleStartReplay(record.id)}
                              disabled={record.status !== 'completed'}
                              className={cn(
                                "p-2 border rounded-lg transition-colors flex items-center gap-2",
                                record.status === 'completed'
                                  ? "bg-primary/20 text-primary border-primary/30 hover:bg-primary/30"
                                  : "bg-white/5 text-white/40 border-white/10 cursor-not-allowed"
                              )}
                              title={record.status === 'completed' ? "Replay" : "Can only replay completed executions"}
                            >
                              <Play size={16} />
                            </button>
                          )}
                          
                          <button
                            onClick={() => handleExportRecord(record.id)}
                            className="p-2 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
                            title="Export"
                          >
                            <Download size={16} />
                          </button>
                          
                          <button
                            onClick={() => handleDeleteRecord(record.id)}
                            className="p-2 bg-white/5 border border-white/10 rounded-lg hover:bg-red-500/20 hover:text-red-400 hover:border-red-500/30 transition-colors"
                            title="Delete"
                          >
                            <Trash2 size={16} />
                          </button>
                        </div>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="border-t border-white/10 p-4 flex justify-between items-center">
            <div className="text-sm text-white/60">
              Showing {(currentPage - 1) * pageSize + 1} to {Math.min(currentPage * pageSize, totalRecords)} of {totalRecords} records
            </div>
            <div className="flex items-center gap-2">
              <button
                onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                disabled={currentPage === 1}
                className="p-2 bg-white/5 border border-white/10 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white/10 transition-colors"
              >
                <ChevronLeft size={16} />
              </button>
              
              <div className="flex items-center gap-1">
                {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                  let pageNum;
                  if (totalPages <= 5) {
                    pageNum = i + 1;
                  } else if (currentPage <= 3) {
                    pageNum = i + 1;
                  } else if (currentPage >= totalPages - 2) {
                    pageNum = totalPages - 4 + i;
                  } else {
                    pageNum = currentPage - 2 + i;
                  }
                  
                  return (
                    <button
                      key={pageNum}
                      onClick={() => setCurrentPage(pageNum)}
                      className={cn(
                        "w-8 h-8 rounded-lg text-sm font-medium transition-colors",
                        currentPage === pageNum
                          ? "bg-primary text-black"
                          : "bg-white/5 text-white/60 hover:bg-white/10"
                      )}
                    >
                      {pageNum}
                    </button>
                  );
                })}
              </div>
              
              <button
                onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
                disabled={currentPage === totalPages}
                className="p-2 bg-white/5 border border-white/10 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white/10 transition-colors"
              >
                <ChevronRight size={16} />
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}