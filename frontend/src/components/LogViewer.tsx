import React, { useState, useEffect, useRef, useMemo, useCallback } from 'react';
import { 
  Search, 
  Filter, 
  Download, 
  Trash2, 
  RefreshCw,
  ChevronDown,
  ChevronUp,
  Clock,
  AlertCircle,
  CheckCircle2,
  Info,
  XCircle,
  Terminal,
  Copy,
  ExternalLink,
  Maximize2,
  Minimize2,
  Play,
  Pause,
  ZoomIn,
  ZoomOut,
  Calendar,
  User,
  ArrowDownToLine,
  Zap
} from 'lucide-react';
import { motion, AnimatePresence } from 'motion/react';
import { cn } from '../utils';
import { LogEntry, Agent } from '../types';
import { api } from '../services/api';
import { useVirtualScroll } from '../hooks/useVirtualScroll';

const LOG_ITEM_HEIGHT = 72;
const LOG_ITEM_HEIGHT_COMPACT = 56;
const LOG_ITEM_HEIGHT_EXPANDED = 200;
const PERF_MODE_THRESHOLD = 1000;
const ANIMATION_THRESHOLD = 500;

interface LogViewerProps {
  className?: string;
  initialLogs?: LogEntry[];
  autoRefresh?: boolean;
  showFilters?: boolean;
  compact?: boolean;
  onLogClick?: (log: LogEntry) => void;
}

type LogLevel = 'all' | 'info' | 'debug' | 'success' | 'warn' | 'error' | 'sys';
type TimeRange = 'all' | '5m' | '30m' | '1h' | 'today' | 'custom';

export default function LogViewer({ 
  className, 
  initialLogs = [], 
  autoRefresh = true,
  showFilters = true,
  compact = false,
  onLogClick 
}: LogViewerProps) {
  const [logs, setLogs] = useState<LogEntry[]>(initialLogs);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedLevel, setSelectedLevel] = useState<LogLevel>('all');
  const [selectedTimeRange, setSelectedTimeRange] = useState<TimeRange>('all');
  const [selectedAgent, setSelectedAgent] = useState<string>('all');
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [isAutoRefresh, setIsAutoRefresh] = useState(autoRefresh);
  const [isFiltersOpen, setIsFiltersOpen] = useState(true);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [agents, setAgents] = useState<Agent[]>([]);
  const [expandedLogs, setExpandedLogs] = useState<Set<string>>(new Set());
  const [perfMode, setPerfMode] = useState(false);
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const refreshIntervalRef = useRef<NodeJS.Timeout>();
  const prevLogsLengthRef = useRef(0);

  const [logLevels] = useState<{ level: LogLevel; label: string; icon: React.ReactNode; count: number }[]>([
    { level: 'all', label: '全部', icon: <Filter className="w-4 h-4" />, count: 0 },
    { level: 'error', label: '错误', icon: <XCircle className="w-4 h-4" />, count: 0 },
    { level: 'warn', label: '警告', icon: <AlertCircle className="w-4 h-4" />, count: 0 },
    { level: 'success', label: '成功', icon: <CheckCircle2 className="w-4 h-4" />, count: 0 },
    { level: 'info', label: '信息', icon: <Info className="w-4 h-4" />, count: 0 },
    { level: 'debug', label: '调试', icon: <Terminal className="w-4 h-4" />, count: 0 },
    { level: 'sys', label: '系统', icon: <Terminal className="w-4 h-4" />, count: 0 },
  ]);
  const [timeRanges] = useState<{ range: TimeRange; label: string; icon: React.ReactNode }[]>([
    { range: 'all', label: '全部时间', icon: <Clock className="w-4 h-4" /> },
    { range: '5m', label: '最近5分钟', icon: <Clock className="w-4 h-4" /> },
    { range: '30m', label: '最近30分钟', icon: <Clock className="w-4 h-4" /> },
    { range: '1h', label: '最近1小时', icon: <Clock className="w-4 h-4" /> },
    { range: 'today', label: '今天', icon: <Calendar className="w-4 h-4" /> },
    { range: 'custom', label: '自定义', icon: <Calendar className="w-4 h-4" /> },
  ]);

  const fetchLogs = async () => {
    try {
      setIsRefreshing(true);
      const newLogs = await api.getLogs();
      if (Array.isArray(newLogs)) {
        const formattedLogs = newLogs.map((log: any) => ({
          id: log.id || `${log.timestamp}-${log.level}-${Math.random()}`,
          timestamp: log.timestamp,
          level: log.level as LogLevel,
          message: log.message,
          agentId: log.agent_id,
          taskId: log.task_id,
          details: log.details,
          duration: log.duration,
          tags: log.tags || [],
        }));
        setLogs(formattedLogs);
      }
    } catch (error) {
      console.error('获取日志失败:', error);
    } finally {
      setIsRefreshing(false);
    }
  };

  const fetchAgents = async () => {
    try {
      const agentsList = await api.listAgents();
      if (Array.isArray(agentsList)) {
        const formattedAgents = agentsList.map((agent: any) => ({
          id: agent.id,
          name: agent.name,
          role: agent.role,
          status: 'idle' as const,
          currentTask: '',
          progress: 0,
          timeActive: '',
          type: agent.role.toLowerCase().includes('frontend') ? 'frontend' : 
                agent.role.toLowerCase().includes('backend') ? 'backend' :
                agent.role.toLowerCase().includes('architect') ? 'architect' :
                agent.role.toLowerCase().includes('qa') ? 'qa' : 'orchestrator',
        }));
        setAgents(formattedAgents);
      }
    } catch (error) {
      console.error('获取Agent列表失败:', error);
    }
  };

  useEffect(() => {
    fetchLogs();
    fetchAgents();
  }, []);

  useEffect(() => {
    if (isAutoRefresh) {
      refreshIntervalRef.current = setInterval(fetchLogs, 5000);
    } else if (refreshIntervalRef.current) {
      clearInterval(refreshIntervalRef.current);
    }
    return () => {
      if (refreshIntervalRef.current) {
        clearInterval(refreshIntervalRef.current);
      }
    };
  }, [isAutoRefresh]);

  const filteredLogs = useMemo(() => {
    let filtered = [...logs];
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(log => 
        log.message.toLowerCase().includes(query) ||
        log.agentId?.toLowerCase().includes(query) ||
        log.tags?.some(tag => tag.toLowerCase().includes(query))
      );
    }
    if (selectedLevel !== 'all') {
      filtered = filtered.filter(log => log.level === selectedLevel);
    }
    if (selectedTimeRange !== 'all') {
      const now = new Date();
      let cutoffTime = new Date();
      switch (selectedTimeRange) {
        case '5m':
          cutoffTime.setMinutes(now.getMinutes() - 5);
          break;
        case '30m':
          cutoffTime.setMinutes(now.getMinutes() - 30);
          break;
        case '1h':
          cutoffTime.setHours(now.getHours() - 1);
          break;
        case 'today':
          cutoffTime.setHours(0, 0, 0, 0);
          break;
      }
      filtered = filtered.filter(log => {
        const logTime = new Date(log.timestamp);
        return logTime >= cutoffTime;
      });
    }
    if (selectedAgent !== 'all') {
      filtered = filtered.filter(log => log.agentId === selectedAgent);
    }
    return filtered;
  }, [logs, searchQuery, selectedLevel, selectedTimeRange, selectedAgent]);

  const totalLogCount = logs.length;
  const shouldDisableAnimations = totalLogCount > ANIMATION_THRESHOLD || perfMode;
  const shouldEnablePerfMode = totalLogCount > PERF_MODE_THRESHOLD;

  useEffect(() => {
    if (shouldEnablePerfMode) {
      setPerfMode(true);
    }
  }, [shouldEnablePerfMode]);

  const getItemHeight = useCallback(
    (index: number): number => {
      const log = filteredLogs[index];
      if (!log) return compact ? LOG_ITEM_HEIGHT_COMPACT : LOG_ITEM_HEIGHT;
      const baseHeight = compact ? LOG_ITEM_HEIGHT_COMPACT : LOG_ITEM_HEIGHT;
      if (expandedLogs.has(log.id) && log.details) {
        return baseHeight + LOG_ITEM_HEIGHT_EXPANDED;
      }
      return baseHeight;
    },
    [filteredLogs, expandedLogs, compact]
  );

  const itemHeight = compact ? LOG_ITEM_HEIGHT_COMPACT : LOG_ITEM_HEIGHT;

  const virtual = useVirtualScroll({
    totalItems: filteredLogs.length,
    itemHeight,
    containerRef: scrollContainerRef,
    overscan: perfMode ? 3 : 10,
    getItemHeight: perfMode ? undefined : getItemHeight,
  });

  useEffect(() => {
    if (virtual.isScrolledToBottom && filteredLogs.length > prevLogsLengthRef.current) {
      virtual.scrollToBottom();
    }
    prevLogsLengthRef.current = filteredLogs.length;
  }, [filteredLogs.length, virtual.isScrolledToBottom, virtual.scrollToBottom]);

  const toggleLogExpansion = (logId: string) => {
    setExpandedLogs(prev => {
      const newSet = new Set(prev);
      if (newSet.has(logId)) {
        newSet.delete(logId);
      } else {
        newSet.add(logId);
      }
      return newSet;
    });
  };

  const copyLogToClipboard = (log: LogEntry) => {
    const logText = `[${log.timestamp}] [${log.level.toUpperCase()}] ${log.message}`;
    navigator.clipboard.writeText(logText);
  };

  const clearLogs = () => {
    setLogs([]);
  };

  const exportLogs = () => {
    const logData = filteredLogs.map(log => ({
      timestamp: log.timestamp,
      level: log.level,
      message: log.message,
      agentId: log.agentId,
      taskId: log.taskId,
      details: log.details,
      duration: log.duration,
      tags: log.tags,
    }));
    const blob = new Blob([JSON.stringify(logData, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `genpulse-logs-${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const getLogLevelStyle = (level: LogLevel) => {
    switch (level) {
      case 'error':
        return 'bg-red-500/10 text-red-400 border-red-500/20';
      case 'warn':
        return 'bg-yellow-500/10 text-yellow-400 border-yellow-500/20';
      case 'success':
        return 'bg-green-500/10 text-green-400 border-green-500/20';
      case 'info':
        return 'bg-blue-500/10 text-blue-400 border-blue-500/20';
      case 'debug':
        return 'bg-purple-500/10 text-purple-400 border-purple-500/20';
      case 'sys':
        return 'bg-gray-500/10 text-gray-400 border-gray-500/20';
      default:
        return 'bg-gray-500/10 text-gray-400 border-gray-500/20';
    }
  };

  const getLogLevelIcon = (level: LogLevel) => {
    switch (level) {
      case 'error':
        return <XCircle className="w-4 h-4" />;
      case 'warn':
        return <AlertCircle className="w-4 h-4" />;
      case 'success':
        return <CheckCircle2 className="w-4 h-4" />;
      case 'info':
        return <Info className="w-4 h-4" />;
      case 'debug':
      case 'sys':
        return <Terminal className="w-4 h-4" />;
      default:
        return <Info className="w-4 h-4" />;
    }
  };

  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString('zh-CN', { 
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  };

  const formatDate = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    });
  };

  const renderLogItem = useCallback(
    (log: LogEntry, index: number) => {
      const isExpanded = expandedLogs.has(log.id);
      const isToday = new Date(log.timestamp).toDateString() === new Date().toDateString();

      return (
        <div
          key={log.id}
          className={cn(
            "rounded-lg border p-3 cursor-pointer transition-all hover:border-white/20 virtual-log-item",
            getLogLevelStyle(log.level),
            isExpanded && "border-white/30"
          )}
          onClick={() => {
            toggleLogExpansion(log.id);
            onLogClick?.(log);
          }}
        >
          <div className="flex items-start justify-between">
            <div className="flex items-start space-x-3 flex-1">
              <div className="mt-0.5">
                {getLogLevelIcon(log.level)}
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center space-x-2 mb-1">
                  <span className="text-xs font-mono text-white/60">
                    {isToday ? formatTime(log.timestamp) : formatDate(log.timestamp)}
                  </span>
                  <span className={cn(
                    "text-xs px-1.5 py-0.5 rounded uppercase font-semibold",
                    getLogLevelStyle(log.level).split(' ')[0]
                  )}>
                    {log.level}
                  </span>
                  {log.agentId && (
                    <span className="text-xs px-1.5 py-0.5 rounded bg-white/10 text-white/60">
                      {agents.find(a => a.id === log.agentId)?.name || log.agentId}
                    </span>
                  )}
                  {log.duration && (
                    <span className="text-xs px-1.5 py-0.5 rounded bg-white/10 text-white/60">
                      {log.duration}ms
                    </span>
                  )}
                </div>
                <p className="text-sm text-white leading-relaxed break-words">
                  {log.message}
                </p>
                {log.tags && log.tags.length > 0 && (
                  <div className="flex flex-wrap gap-1 mt-2">
                    {log.tags.map((tag, idx) => (
                      <span
                        key={idx}
                        className="text-xs px-1.5 py-0.5 rounded bg-white/5 text-white/40"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            </div>
            <div className="flex items-center space-x-1 ml-2">
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  copyLogToClipboard(log);
                }}
                className="p-1 rounded hover:bg-white/10 text-white/40 hover:text-white transition-colors"
                title="复制日志"
              >
                <Copy className="w-3 h-3" />
              </button>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  toggleLogExpansion(log.id);
                }}
                className="p-1 rounded hover:bg-white/10 text-white/40 hover:text-white transition-colors"
                title={isExpanded ? "收起详情" : "展开详情"}
              >
                {isExpanded ? <ChevronUp className="w-3 h-3" /> : <ChevronDown className="w-3 h-3" />}
              </button>
            </div>
          </div>

          {isExpanded && log.details && (
            <div className="mt-3 pt-3 border-t border-white/10">
              <div className="bg-black/20 rounded p-3">
                <pre className="text-xs text-white/80 font-mono whitespace-pre-wrap break-words">
                  {typeof log.details === 'string' 
                    ? log.details 
                    : JSON.stringify(log.details, null, 2)}
                </pre>
              </div>
            </div>
          )}
        </div>
      );
    },
    [expandedLogs, agents, onLogClick]
  );

  return (
    <div className={cn(
      "flex flex-col bg-background border border-white/10 rounded-lg overflow-hidden",
      isFullscreen ? "fixed inset-4 z-50" : "h-full",
      className
    )}>
      <div className="flex items-center justify-between p-4 border-b border-white/10 bg-background/50 backdrop-blur-sm">
        <div className="flex items-center space-x-4">
          <h2 className="text-lg font-semibold text-on-surface">执行日志</h2>
          <div className="flex items-center space-x-2">
            <button
              onClick={fetchLogs}
              disabled={isRefreshing}
              className={cn(
                "p-2 rounded-lg transition-colors",
                isRefreshing 
                  ? "bg-primary/20 text-primary cursor-not-allowed" 
                  : "hover:bg-white/5 text-white/60 hover:text-white"
              )}
            >
              <RefreshCw className={cn("w-4 h-4", isRefreshing && "animate-spin")} />
            </button>
            <button
              onClick={() => setIsAutoRefresh(!isAutoRefresh)}
              className={cn(
                "p-2 rounded-lg transition-colors",
                isAutoRefresh
                  ? "bg-green-500/20 text-green-400" 
                  : "hover:bg-white/5 text-white/60 hover:text-white"
              )}
            >
              {isAutoRefresh ? <Pause className="w-4 h-4" /> : <Play className="w-4 h-4" />}
            </button>
            <button
              onClick={exportLogs}
              className="p-2 rounded-lg hover:bg-white/5 text-white/60 hover:text-white transition-colors"
            >
              <Download className="w-4 h-4" />
            </button>
            <button
              onClick={clearLogs}
              className="p-2 rounded-lg hover:bg-red-500/10 text-red-400 hover:text-red-300 transition-colors"
            >
              <Trash2 className="w-4 h-4" />
            </button>
            <button
              onClick={() => setIsFullscreen(!isFullscreen)}
              className="p-2 rounded-lg hover:bg-white/5 text-white/60 hover:text-white transition-colors"
            >
              {isFullscreen ? <Minimize2 className="w-4 h-4" /> : <Maximize2 className="w-4 h-4" />}
            </button>
          </div>
        </div>

        <div className="flex items-center space-x-2">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-white/40" />
            <input
              type="text"
              placeholder="搜索日志..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10 pr-4 py-2 bg-white/5 border border-white/10 rounded-lg text-sm text-white placeholder-white/40 focus:outline-none focus:border-primary/50 focus:ring-1 focus:ring-primary/20 w-64"
            />
          </div>
          {showFilters && (
            <button
              onClick={() => setIsFiltersOpen(!isFiltersOpen)}
              className="p-2 rounded-lg hover:bg-white/5 text-white/60 hover:text-white transition-colors"
            >
              {isFiltersOpen ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
            </button>
          )}
        </div>
      </div>

      <AnimatePresence>
        {showFilters && isFiltersOpen && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            className="overflow-hidden border-b border-white/10"
          >
            <div className="p-4 space-y-4">
              <div>
                <h3 className="text-sm font-medium text-white/60 mb-2">日志级别</h3>
                <div className="flex flex-wrap gap-2">
                  {logLevels.map(levelInfo => (
                    <button
                      key={levelInfo.level}
                      onClick={() => setSelectedLevel(levelInfo.level)}
                      className={cn(
                        "flex items-center space-x-2 px-3 py-1.5 rounded-lg border transition-colors",
                        selectedLevel === levelInfo.level
                          ? getLogLevelStyle(levelInfo.level)
                          : "border-white/10 hover:border-white/20 text-white/60 hover:text-white"
                      )}
                    >
                      {levelInfo.icon}
                      <span className="text-sm">{levelInfo.label}</span>
                      {levelInfo.count > 0 && (
                        <span className={cn(
                          "text-xs px-1.5 py-0.5 rounded",
                          selectedLevel === levelInfo.level
                            ? "bg-white/20"
                            : "bg-white/10"
                        )}>
                          {levelInfo.count}
                        </span>
                      )}
                    </button>
                  ))}
                </div>
              </div>

              <div>
                <h3 className="text-sm font-medium text-white/60 mb-2">时间范围</h3>
                <div className="flex flex-wrap gap-2">
                  {timeRanges.map(rangeInfo => (
                    <button
                      key={rangeInfo.range}
                      onClick={() => setSelectedTimeRange(rangeInfo.range)}
                      className={cn(
                        "flex items-center space-x-2 px-3 py-1.5 rounded-lg border transition-colors",
                        selectedTimeRange === rangeInfo.range
                          ? "bg-primary/10 text-primary border-primary/20"
                          : "border-white/10 hover:border-white/20 text-white/60 hover:text-white"
                      )}
                    >
                      {rangeInfo.icon}
                      <span className="text-sm">{rangeInfo.label}</span>
                    </button>
                  ))}
                </div>
              </div>

              {agents.length > 0 && (
                <div>
                  <h3 className="text-sm font-medium text-white/60 mb-2">Agent</h3>
                  <div className="flex flex-wrap gap-2">
                    <button
                      onClick={() => setSelectedAgent('all')}
                      className={cn(
                        "flex items-center space-x-2 px-3 py-1.5 rounded-lg border transition-colors",
                        selectedAgent === 'all'
                          ? "bg-primary/10 text-primary border-primary/20"
                          : "border-white/10 hover:border-white/20 text-white/60 hover:text-white"
                      )}
                    >
                      <User className="w-4 h-4" />
                      <span className="text-sm">全部Agent</span>
                    </button>
                    {agents.map(agent => (
                      <button
                        key={agent.id}
                        onClick={() => setSelectedAgent(agent.id)}
                        className={cn(
                          "flex items-center space-x-2 px-3 py-1.5 rounded-lg border transition-colors",
                          selectedAgent === agent.id
                            ? "bg-primary/10 text-primary border-primary/20"
                            : "border-white/10 hover:border-white/20 text-white/60 hover:text-white"
                        )}
                      >
                        <User className="w-4 h-4" />
                        <span className="text-sm">{agent.name}</span>
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      <div className="px-4 py-2 border-b border-white/10 bg-background/30">
        <div className="flex items-center justify-between text-sm">
          <div className="flex items-center space-x-4">
            <span className="text-white/60">
              共 {filteredLogs.length} 条日志
              {searchQuery && ` (搜索: "${searchQuery}")`}
            </span>
            {selectedLevel !== 'all' && (
              <span className="text-white/40">
                级别: {logLevels.find(l => l.level === selectedLevel)?.label}
              </span>
            )}
            {selectedTimeRange !== 'all' && (
              <span className="text-white/40">
                时间: {timeRanges.find(t => t.range === selectedTimeRange)?.label}
              </span>
            )}
          </div>
          <div className="flex items-center space-x-3">
            {totalLogCount > ANIMATION_THRESHOLD && (
              <button
                onClick={() => setPerfMode(!perfMode)}
                className={cn(
                  "flex items-center space-x-1 px-2 py-1 rounded text-xs transition-colors",
                  perfMode
                    ? "bg-yellow-500/20 text-yellow-400"
                    : "bg-white/5 text-white/40 hover:text-white/60"
                )}
                title={perfMode ? "关闭性能模式" : "开启性能模式 (禁用动画，减少DOM)"}
              >
                <Zap className="w-3 h-3" />
                <span>性能模式</span>
              </button>
            )}
            <button
              onClick={() => {
                setSearchQuery('');
                setSelectedLevel('all');
                setSelectedTimeRange('all');
                setSelectedAgent('all');
              }}
              className="text-sm text-white/60 hover:text-white transition-colors"
            >
              清除所有过滤器
            </button>
          </div>
        </div>
      </div>

      <div
        ref={scrollContainerRef}
        onScroll={virtual.onScroll}
        className={cn(
          "flex-1 virtual-scroll-container",
          compact ? "p-2" : "p-4"
        )}
      >
        {filteredLogs.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-white/40">
            <Terminal className="w-12 h-12 mb-4" />
            <p className="text-lg">暂无日志</p>
            <p className="text-sm mt-2">开始执行任务后，日志将显示在这里</p>
          </div>
        ) : (
          <div
            className="relative"
            style={{ height: virtual.totalHeight }}
          >
            <div
              className="space-y-2"
              style={{ transform: `translateY(${virtual.offsetY}px)` }}
            >
              {virtual.visibleItems.map(({ index }) => renderLogItem(filteredLogs[index], index))}
            </div>
          </div>
        )}

        {!virtual.isScrolledToBottom && filteredLogs.length > 0 && (
          <button
            onClick={virtual.scrollToBottom}
            className="scroll-to-bottom-btn fixed bottom-20 right-8 z-40 flex items-center space-x-2 px-4 py-2 rounded-full bg-primary/20 text-primary border border-primary/30 backdrop-blur-md hover:bg-primary/30 transition-all shadow-lg"
          >
            <ArrowDownToLine className="w-4 h-4" />
            <span className="text-sm font-medium">滚动到底部</span>
          </button>
        )}
      </div>

      <div className="px-4 py-2 border-t border-white/10 bg-background/50 flex items-center justify-between text-xs text-white/40">
        <div className="flex items-center space-x-4">
          <span>最后更新: {logs.length > 0 ? formatTime(logs[0].timestamp) : '--:--:--'}</span>
          <span className={cn(
            "flex items-center space-x-1",
            isAutoRefresh ? "text-green-400" : "text-yellow-400"
          )}>
            {isAutoRefresh ? (
              <>
                <RefreshCw className="w-3 h-3 animate-spin" />
                <span>自动刷新中</span>
              </>
            ) : (
              <>
                <Pause className="w-3 h-3" />
                <span>自动刷新已暂停</span>
              </>
            )}
          </span>
          {perfMode && (
            <span className="flex items-center space-x-1 text-yellow-400">
              <Zap className="w-3 h-3" />
              <span>性能模式</span>
            </span>
          )}
        </div>
        <div className="flex items-center space-x-2">
          <span>日志级别: {selectedLevel}</span>
          {filteredLogs.length > 0 && (
            <>
              <span>|</span>
              <span>
                显示 {virtual.startIndex + 1}-{Math.min(virtual.endIndex + 1, filteredLogs.length)} / {filteredLogs.length}
              </span>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
