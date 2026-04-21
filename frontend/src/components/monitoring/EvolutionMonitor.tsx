import React, { useState, useMemo } from 'react';
import { 
  TrendingUp, 
  Brain, 
  Zap, 
  Target, 
  Award, 
  Trophy, 
  Crown, 
  Gem, 
  Diamond,
  Star,
  Heart,
  Flag,
  Bookmark,
  Tag,
  Hash,
  Users,
  Cpu,
  Database,
  Server,
  FileText,
  Code,
  MessageSquare,
  GitBranch,
  Layout,
  CheckCircle2,
  Shield,
  Terminal,
  Settings,
  Filter,
  Search,
  Download,
  RefreshCw,
  ChevronUp,
  ChevronDown,
  Eye,
  EyeOff,
  Clock,
  Calendar,
  BarChart3,
  PieChart,
  LineChart,
  Activity,
  AlertCircle,
  CheckCircle,
  XCircle,
  Plus,
  Minus,
  Edit,
  Trash2,
  Save,
  Copy,
  ExternalLink,
  Maximize2,
  Minimize2,
  MoreVertical,
  Play,
  Pause,
  SkipBack,
  SkipForward,
  Volume2,
  VolumeX
} from 'lucide-react';
import { cn } from '../../utils';
import { EvolutionEvent } from '../../types';

interface EvolutionMonitorProps {
  events: EvolutionEvent[];
}

export default function EvolutionMonitor({ events }: EvolutionMonitorProps) {
  const [timeRange, setTimeRange] = useState<'day' | 'week' | 'month' | 'year'>('week');
  const [eventTypeFilter, setEventTypeFilter] = useState<'all' | 'skill_generated' | 'memory_updated' | 'agent_improved' | 'error_fixed'>('all');
  const [agentFilter, setAgentFilter] = useState<string>('all');
  const [showDetails, setShowDetails] = useState(false);
  const [selectedEvent, setSelectedEvent] = useState<string | null>(null);

  // 获取所有代理
  const agents = useMemo(() => {
    const agentSet = new Set(events.map(e => e.agent).filter(Boolean));
    return Array.from(agentSet);
  }, [events]);

  // 获取所有事件类型
  const eventTypes = useMemo(() => {
    const typeSet = new Set(events.map(e => e.eventType));
    return Array.from(typeSet);
  }, [events]);

  // 过滤事件
  const filteredEvents = useMemo(() => {
    let filtered = events.filter(event => {
      // 事件类型过滤
      if (eventTypeFilter !== 'all' && event.eventType !== eventTypeFilter) {
        return false;
      }
      
      // 代理过滤
      if (agentFilter !== 'all' && event.agent !== agentFilter) {
        return false;
      }
      
      return true;
    });
    
    // 按时间排序（最新的在前）
    filtered.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime());
    
    return filtered;
  }, [events, eventTypeFilter, agentFilter]);

  // 计算统计数据
  const stats = useMemo(() => {
    const totalEvents = filteredEvents.length;
    const skillEvents = filteredEvents.filter(e => e.eventType === 'skill_generated').length;
    const memoryEvents = filteredEvents.filter(e => e.eventType === 'memory_updated').length;
    const improvementEvents = filteredEvents.filter(e => e.eventType === 'agent_improved').length;
    
    // 计算效率提升
    const efficiencyGain = filteredEvents.reduce((sum, event) => sum + (event.efficiencyGain || 0), 0);
    const avgEfficiencyGain = totalEvents > 0 ? efficiencyGain / totalEvents : 0;
    
    // 计算token节省
    const tokenSavings = filteredEvents.reduce((sum, event) => sum + (event.tokenSavings || 0), 0);
    
    // 计算时间节省（分钟）
    const timeSavings = filteredEvents.reduce((sum, event) => sum + (event.timeSavings || 0), 0);
    
    return {
      totalEvents,
      skillEvents,
      memoryEvents,
      improvementEvents,
      efficiencyGain: avgEfficiencyGain,
      tokenSavings,
      timeSavings
    };
  }, [filteredEvents]);

  const getEventTypeIcon = (eventType: string) => {
    switch (eventType) {
      case 'skill_generated': return Brain;
      case 'memory_updated': return Database;
      case 'agent_improved': return TrendingUp;
      case 'error_fixed': return CheckCircle2;
      default: return Activity;
    }
  };

  const getEventTypeColor = (eventType: string) => {
    switch (eventType) {
      case 'skill_generated': return 'text-purple-500';
      case 'memory_updated': return 'text-blue-500';
      case 'agent_improved': return 'text-green-500';
      case 'error_fixed': return 'text-yellow-500';
      default: return 'text-white/60';
    }
  };

  const getEventTypeBgColor = (eventType: string) => {
    switch (eventType) {
      case 'skill_generated': return 'bg-purple-500/10';
      case 'memory_updated': return 'bg-blue-500/10';
      case 'agent_improved': return 'bg-green-500/10';
      case 'error_fixed': return 'bg-yellow-500/10';
      default: return 'bg-white/5';
    }
  };

  const getEventTypeLabel = (eventType: string) => {
    switch (eventType) {
      case 'skill_generated': return 'Skill Generated';
      case 'memory_updated': return 'Memory Updated';
      case 'agent_improved': return 'Agent Improved';
      case 'error_fixed': return 'Error Fixed';
      default: return 'Evolution Event';
    }
  };

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString('en-US', { 
      hour12: false,
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const formatDate = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric'
    });
  };

  const formatTokenCount = (tokens: number) => {
    if (tokens < 1000) return `${tokens}`;
    if (tokens < 1000000) return `${(tokens / 1000).toFixed(1)}K`;
    return `${(tokens / 1000000).toFixed(2)}M`;
  };

  const formatTime = (minutes: number) => {
    if (minutes < 60) return `${minutes}m`;
    if (minutes < 1440) return `${(minutes / 60).toFixed(1)}h`;
    return `${(minutes / 1440).toFixed(1)}d`;
  };

  const handleDownload = () => {
    const data = {
      events: filteredEvents,
      stats,
      timestamp: new Date().toISOString(),
      filters: { timeRange, eventTypeFilter, agentFilter }
    };
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `evolution-events-${new Date().toISOString().split('T')[0]}.json`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const selectedEventData = selectedEvent ? events.find(e => e.id === selectedEvent) : null;

  // 模拟时间序列数据
  const timeSeriesData = useMemo(() => {
    const days = timeRange === 'day' ? 24 : 
                 timeRange === 'week' ? 7 : 
                 timeRange === 'month' ? 30 : 365;
    
    return Array.from({ length: days }, (_, i) => {
      const date = new Date();
      date.setDate(date.getDate() - (days - 1 - i));
      
      // 模拟进化事件数据
      const skillCount = Math.floor(Math.random() * 5);
      const memoryCount = Math.floor(Math.random() * 10);
      const improvementCount = Math.floor(Math.random() * 3);
      const totalEvents = skillCount + memoryCount + improvementCount;
      
      return {
        date: date.toISOString().split('T')[0],
        skillCount,
        memoryCount,
        improvementCount,
        totalEvents,
        efficiency: 70 + Math.random() * 20 // 模拟效率提升
      };
    });
  }, [timeRange]);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-3">
            <TrendingUp size={24} />
            Self-Evolution Monitor
          </h2>
          <p className="text-sm text-white/60 mt-1">
            Tracking AI agent learning, skill generation, and memory updates
          </p>
        </div>
        
        <div className="flex items-center gap-3">
          <div className="text-sm text-white/40">
            {stats.totalEvents} events • {stats.skillEvents} skills • {formatTokenCount(stats.tokenSavings)} tokens saved
          </div>
        </div>
      </div>

      {/* Controls */}
      <div className="bg-white/5 border border-white/10 rounded-xl p-4">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <select
              value={timeRange}
              onChange={(e) => setTimeRange(e.target.value as any)}
              className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
            >
              <option value="day">Last 24 Hours</option>
              <option value="week">Last 7 Days</option>
              <option value="month">Last 30 Days</option>
              <option value="year">Last Year</option>
            </select>
            
            <select
              value={eventTypeFilter}
              onChange={(e) => setEventTypeFilter(e.target.value as any)}
              className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
            >
              <option value="all">All Event Types</option>
              {eventTypes.map(type => (
                <option key={type} value={type}>{getEventTypeLabel(type)}</option>
              ))}
            </select>
            
            <select
              value={agentFilter}
              onChange={(e) => setAgentFilter(e.target.value)}
              className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
            >
              <option value="all">All Agents</option>
              {agents.map(agent => (
                <option key={agent} value={agent}>{agent}</option>
              ))}
            </select>
          </div>
          
          <div className="flex items-center gap-3">
            <button
              onClick={() => setShowDetails(!showDetails)}
              className={cn(
                "px-3 py-2 text-sm rounded-lg transition-colors flex items-center gap-2",
                showDetails
                  ? "bg-primary/20 text-primary" 
                  : "bg-white/5 text-white/60 hover:bg-white/10"
              )}
            >
              {showDetails ? <EyeOff size={16} /> : <Eye size={16} />}
              {showDetails ? 'Hide Details' : 'Show Details'}
            </button>
            
            <button
              onClick={handleDownload}
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
              title="Download evolution data"
            >
              <Download size={20} />
            </button>
          </div>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white/5 border border-white/10 rounded-xl p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="text-xs text-white/40 uppercase tracking-wider">Total Evolution Events</div>
            <Activity size={20} className="text-primary" />
          </div>
          <div className="text-3xl font-bold">{stats.totalEvents}</div>
          <div className="text-sm text-white/60 mt-2">
            Across {agents.length} agents
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-xl p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="text-xs text-white/40 uppercase tracking-wider">Efficiency Gain</div>
            <Zap size={20} className="text-green-500" />
          </div>
          <div className="text-3xl font-bold text-green-500">
            {stats.efficiencyGain.toFixed(1)}%
          </div>
          <div className="text-sm text-white/60 mt-2">
            Average improvement per event
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-xl p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="text-xs text-white/40 uppercase tracking-wider">Token Savings</div>
            <Brain size={20} className="text-purple-500" />
          </div>
          <div className="text-3xl font-bold">
            {formatTokenCount(stats.tokenSavings)}
          </div>
          <div className="text-sm text-white/60 mt-2">
            Estimated cost reduction
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-xl p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="text-xs text-white/40 uppercase tracking-wider">Time Saved</div>
            <Clock size={20} className="text-blue-500" />
          </div>
          <div className="text-3xl font-bold">{formatTime(stats.timeSavings)}</div>
          <div className="text-sm text-white/60 mt-2">
            Cumulative time optimization
          </div>
        </div>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Events Over Time */}
        <div className="border border-white/10 rounded-xl p-6">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-lg font-bold">Evolution Events Over Time</h3>
            <div className="flex items-center gap-2">
              <button className="text-xs px-3 py-1 bg-white/5 rounded hover:bg-white/10 transition-colors">
                Day
              </button>
              <button className="text-xs px-3 py-1 bg-primary/20 text-primary rounded">
                Week
              </button>
              <button className="text-xs px-3 py-1 bg-white/5 rounded hover:bg-white/10 transition-colors">
                Month
              </button>
            </div>
          </div>
          
          <div className="h-64">
            <div className="flex items-end h-48 gap-1">
              {timeSeriesData.map((day, i) => (
                <div key={i} className="flex-1 flex flex-col items-center">
                  <div className="flex items-end justify-center w-full gap-1" style={{ height: '100%' }}>
                    <div 
                      className="flex-1 bg-purple-500/60 rounded-t hover:bg-purple-500 transition-colors"
                      style={{ height: `${(day.skillCount / 5) * 100}%` }}
                      title={`Skills: ${day.skillCount}`}
                    />
                    <div 
                      className="flex-1 bg-blue-500/60 rounded-t hover:bg-blue-500 transition-colors"
                      style={{ height: `${(day.memoryCount / 10) * 100}%` }}
                      title={`Memory: ${day.memoryCount}`}
                    />
                    <div 
                      className="flex-1 bg-green-500/60 rounded-t hover:bg-green-500 transition-colors"
                      style={{ height: `${(day.improvementCount / 3) * 100}%` }}
                      title={`Improvements: ${day.improvementCount}`}
                    />
                  </div>
                  <div className="text-xs text-white/40 mt-2">
                    {timeRange === 'day' 
                      ? `${i}:00` 
                      : day.date.split('-').slice(1).join('/')}
                  </div>
                </div>
              ))}
            </div>
            
            <div className="flex justify-between text-xs text-white/40 mt-4">
              <span>Start</span>
              <span>End</span>
            </div>
          </div>
          
          <div className="flex items-center justify-center gap-6 mt-6">
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 bg-purple-500 rounded-full" />
              <span className="text-sm">Skills</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 bg-blue-500 rounded-full" />
              <span className="text-sm">Memory</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 bg-green-500 rounded-full" />
              <span className="text-sm">Improvements</span>
            </div>
          </div>
        </div>

        {/* Event Type Distribution */}
        <div className="border border-white/10 rounded-xl p-6">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-lg font-bold">Event Type Distribution</h3>
            <PieChart size={20} className="text-white/40" />
          </div>
          
          <div className="flex items-center justify-center h-48">
            <div className="relative w-40 h-40">
              {/* Pie Chart */}
              <svg viewBox="0 0 100 100" className="w-full h-full">
                {[
                  { type: 'skill_generated', count: stats.skillEvents, color: '#A855F7' },
                  { type: 'memory_updated', count: stats.memoryEvents, color: '#3B82F6' },
                  { type: 'agent_improved', count: stats.improvementEvents, color: '#10B981' },
                  { type: 'error_fixed', count: stats.totalEvents - stats.skillEvents - stats.memoryEvents - stats.improvementEvents, color: '#F59E0B' },
                ]
                .filter(item => item.count > 0)
                .reduce((acc, item, index, array) => {
                  const total = stats.totalEvents;
                  const percentage = (item.count / total) * 100;
                  const startAngle = acc.currentAngle;
                  const endAngle = startAngle + (percentage * 3.6);
                  
                  const startX = 50 + 40 * Math.cos((startAngle - 90) * Math.PI / 180);
                  const startY = 50 + 40 * Math.sin((startAngle - 90) * Math.PI / 180);
                  const endX = 50 + 40 * Math.cos((endAngle - 90) * Math.PI / 180);
                  const endY = 50 + 40 * Math.sin((endAngle - 90) * Math.PI / 180);
                  
                  const largeArcFlag = percentage > 50 ? 1 : 0;
                  
                  const path = `M 50 50 L ${startX} ${startY} A 40 40 0 ${largeArcFlag} 1 ${endX} ${endY} Z`;
                  
                  acc.elements.push(
                    React.createElement('path', {
                      key: item.type,
                      d: path,
                      fill: item.color,
                      stroke: "#0A0A0A",
                      strokeWidth: "1"
                    })
                  );
                  
                  acc.currentAngle = endAngle;
                  return acc;
                }, { currentAngle: 0, elements: [] as JSX.Element[] }).elements}
                
                <circle cx="50" cy="50" r="15" fill="#0A0A0A" />
              </svg>
              
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="text-center">
                  <div className="text-2xl font-bold">{stats.totalEvents}</div>
                  <div className="text-xs text-white/60">Total Events</div>
                </div>
              </div>
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-3 mt-6">
            {[
              { type: 'skill_generated', label: 'Skills', color: 'bg-purple-500' },
              { type: 'memory_updated', label: 'Memory', color: 'bg-blue-500' },
              { type: 'agent_improved', label: 'Improvements', color: 'bg-green-500' },
              { type: 'error_fixed', label: 'Error Fixes', color: 'bg-yellow-500' },
            ].map((item) => (
              <div key={item.type} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className={`w-3 h-3 ${item.color} rounded-full`} />
                  <span className="text-sm">{item.label}</span>
                </div>
                <div className="text-right">
                  <div className="text-sm font-bold">
                    {filteredEvents.filter(e => e.eventType === item.type).length}
                  </div>
                  <div className="text-xs text-white/40">
                    {stats.totalEvents > 0 
                      ? `${((filteredEvents.filter(e => e.eventType === item.type).length / stats.totalEvents) * 100).toFixed(1)}%`
                      : '0%'}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Recent Events */}
      <div className="border border-white/10 rounded-xl p-6">
        <h3 className="text-lg font-bold mb-6">Recent Evolution Events</h3>
        
        <div className="space-y-4">
          {filteredEvents.slice(0, 5).map((event) => {
            const Icon = getEventTypeIcon(event.eventType);
            const isSelected = selectedEvent === event.id;
            
            return (
              <div
                key={event.id}
                className={cn(
                  "border rounded-lg p-4 transition-all cursor-pointer hover:scale-[1.01]",
                  getEventTypeBgColor(event.eventType),
                  isSelected && "ring-2 ring-primary"
                )}
                onClick={() => setSelectedEvent(event.id === selectedEvent ? null : event.id)}
              >
                <div className="flex justify-between items-start">
                  <div className="flex items-start gap-3">
                    <div className={cn(
                      "p-2 rounded-lg mt-1",
                      getEventTypeBgColor(event.eventType).replace('border-', 'bg-').replace('/20', '/20')
                    )}>
                      <Icon size={16} className={getEventTypeColor(event.eventType)} />
                    </div>
                    
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="font-bold">{getEventTypeLabel(event.eventType)}</span>
                        <span className="text-xs px-2 py-0.5 bg-white/10 rounded">
                          {event.agent || 'Unknown'}
                        </span>
                        <span className="text-xs text-white/40">
                          {formatTimestamp(event.timestamp)}
                        </span>
                      </div>
                      
                      <p className="text-white/80">
                        {event.description || 'No description available'}
                      </p>
                      
                      {(event.efficiencyGain || event.tokenSavings || event.timeSavings) && (
                        <div className="flex items-center gap-4 mt-3">
                          {event.efficiencyGain && (
                            <div className="flex items-center gap-1 text-xs">
                              <Zap size={12} className="text-green-500" />
                              <span className="text-green-500">+{event.efficiencyGain}% efficiency</span>
                            </div>
                          )}
                          {event.tokenSavings && (
                            <div className="flex items-center gap-1 text-xs">
                              <Brain size={12} className="text-purple-500" />
                              <span className="text-purple-500">Saved {formatTokenCount(event.tokenSavings)} tokens</span>
                            </div>
                          )}
                          {event.timeSavings && (
                            <div className="flex items-center gap-1 text-xs">
                              <Clock size={12} className="text-blue-500" />
                              <span className="text-blue-500">Saved {formatTime(event.timeSavings)}</span>
                            </div>
                          )}
                        </div>
                      )}
                    </div>
                  </div>
                  
                  <div className="flex items-center gap-2">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        navigator.clipboard.writeText(JSON.stringify(event, null, 2));
                      }}
                      className="p-1 rounded hover:bg-white/10 transition-colors"
                      title="Copy event data"
                    >
                      <Copy size={14} />
                    </button>
                  </div>
                </div>
              </div>
            );
          })}
          
          {filteredEvents.length === 0 && (
            <div className="text-center py-12">
              <TrendingUp size={48} className="mx-auto text-white/20 mb-4" />
              <div className="text-lg font-bold text-white/40">No evolution events found</div>
              <p className="text-white/60 mt-2">
                Try changing your filters or wait for agents to generate evolution events
              </p>
            </div>
          )}
          
          {filteredEvents.length > 5 && (
            <div className="text-center pt-4">
              <button className="text-sm text-white/60 hover:text-white transition-colors">
                Show all {filteredEvents.length} events
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Selected Event Details */}
      {selectedEventData && (
        <div className="border border-white/10 rounded-xl p-6 bg-white/[0.02]">
          <div className="flex justify-between items-start mb-6">
            <div>
              <h3 className="text-lg font-bold">Evolution Event Details</h3>
              <p className="text-sm text-white/60 mt-1">
                {getEventTypeLabel(selectedEventData.eventType)} • {selectedEventData.agent || 'Unknown Agent'}
              </p>
            </div>
            
            <button
              onClick={() => setSelectedEvent(null)}
              className="p-2 rounded-lg hover:bg-white/10 transition-colors"
            >
              <ChevronUp size={20} className="rotate-90" />
            </button>
          </div>
          
          <div className="grid grid-cols-2 gap-6">
            <div className="space-y-4">
              <div>
                <div className="text-sm font-bold mb-2">Event Information</div>
                <div className="bg-white/5 border border-white/10 rounded-lg p-4 space-y-3">
                  <div className="flex justify-between">
                    <span className="text-white/60">Type:</span>
                    <span className={cn(
                      "font-bold",
                      getEventTypeColor(selectedEventData.eventType)
                    )}>
                      {getEventTypeLabel(selectedEventData.eventType)}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-white/60">Agent:</span>
                    <span className="font-bold">{selectedEventData.agent || 'Unknown'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-white/60">Timestamp:</span>
                    <span className="font-bold">
                      {new Date(selectedEventData.timestamp).toLocaleString()}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-white/60">Severity:</span>
                    <span className={cn(
                      "font-bold",
                      selectedEventData.severity === 'high' ? 'text-red-500' :
                      selectedEventData.severity === 'medium' ? 'text-yellow-500' :
                      'text-green-500'
                    )}>
                      {selectedEventData.severity?.toUpperCase() || 'LOW'}
                    </span>
                  </div>
                </div>
              </div>
              
              <div>
                <div className="text-sm font-bold mb-2">Impact Metrics</div>
                <div className="bg-white/5 border border-white/10 rounded-lg p-4 space-y-3">
                  {selectedEventData.efficiencyGain && (
                    <div className="flex justify-between">
                      <span className="text-white/60">Efficiency Gain:</span>
                      <span className="font-bold text-green-500">+{selectedEventData.efficiencyGain}%</span>
                    </div>
                  )}
                  {selectedEventData.tokenSavings && (
                    <div className="flex justify-between">
                      <span className="text-white/60">Token Savings:</span>
                      <span className="font-bold text-purple-500">
                        {formatTokenCount(selectedEventData.tokenSavings)}
                      </span>
                    </div>
                  )}
                  {selectedEventData.timeSavings && (
                    <div className="flex justify-between">
                      <span className="text-white/60">Time Savings:</span>
                      <span className="font-bold text-blue-500">
                        {formatTime(selectedEventData.timeSavings)}
                      </span>
                    </div>
                  )}
                  {selectedEventData.successRate && (
                    <div className="flex justify-between">
                      <span className="text-white/60">Success Rate:</span>
                      <span className="font-bold">{selectedEventData.successRate}%</span>
                    </div>
                  )}
                </div>
              </div>
            </div>
            
            <div>
              <div className="text-sm font-bold mb-2">Description</div>
              <div className="bg-white/5 border border-white/10 rounded-lg p-4">
                <p className="text-white/80 leading-relaxed">
                  {selectedEventData.description || 'No description available'}
                </p>
              </div>
              
              {selectedEventData.metadata && (
                <div className="mt-4">
                  <div className="text-sm font-bold mb-2">Metadata</div>
                  <div className="bg-white/5 border border-white/10 rounded-lg p-4 max-h-60 overflow-y-auto">
                    <pre className="text-sm text-white/80">
                      {JSON.stringify(selectedEventData.metadata, null, 2)}
                    </pre>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Evolution Insights */}
      <div className="border border-white/10 rounded-xl p-6">
        <h3 className="text-lg font-bold mb-6 flex items-center gap-3">
          <Brain size={24} />
          Evolution Insights & Recommendations
        </h3>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="flex items-center gap-3 mb-3">
              <Target size={20} className="text-green-500" />
              <div className="text-sm font-bold">Learning Patterns</div>
            </div>
            <ul className="text-sm text-white/80 space-y-2">
              <li>• Agents learn fastest during complex tasks</li>
              <li>• Error recovery generates valuable skills</li>
              <li>• Memory updates peak after successful completions</li>
              <li>• Cross-agent knowledge sharing accelerates learning</li>
            </ul>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="flex items-center gap-3 mb-3">
              <Zap size={20} className="text-yellow-500" />
              <div className="text-sm font-bold">Optimization Opportunities</div>
            </div>
            <ul className="text-sm text-white/80 space-y-2">
              <li>• Increase task complexity for faster learning</li>
              <li>• Implement cross-training between agents</li>
              <li>• Prioritize error-prone areas for skill generation</li>
              <li>• Schedule regular knowledge consolidation</li>
            </ul>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="flex items-center gap-3 mb-3">
              <TrendingUp size={20} className="text-blue-500" />
              <div className="text-sm font-bold">Growth Projection</div>
            </div>
            <div className="space-y-3">
              <div>
                <div className="flex justify-between text-sm">
                  <span>Skill Library Size</span>
                  <span className="font-bold">+{Math.round(stats.skillEvents * 1.5)} in 30 days</span>
                </div>
                <div className="h-2 bg-white/10 rounded-full mt-1 overflow-hidden">
                  <div className="h-full bg-blue-500 rounded-full" style={{ width: '65%' }} />
                </div>
              </div>
              <div>
                <div className="flex justify-between text-sm">
                  <span>Efficiency Gain</span>
                  <span className="font-bold text-green-500">+{Math.round(stats.efficiencyGain * 2)}% in 30 days</span>
                </div>
                <div className="h-2 bg-white/10 rounded-full mt-1 overflow-hidden">
                  <div className="h-full bg-green-500 rounded-full" style={{ width: '45%' }} />
                </div>
              </div>
            </div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="flex items-center gap-3 mb-3">
              <Award size={20} className="text-purple-500" />
              <div className="text-sm font-bold">Achievements</div>
            </div>
            <div className="space-y-3">
              <div className="flex items-center gap-2">
                <Trophy size={16} className="text-yellow-500" />
                <span className="text-sm">Fast Learner: Generated 10+ skills this week</span>
              </div>
              <div className="flex items-center gap-2">
                <Crown size={16} className="text-purple-500" />
                <span className="text-sm">Efficiency Master: 25% average improvement</span>
              </div>
              <div className="flex items-center gap-2">
                <Gem size={16} className="text-blue-500" />
                <span className="text-sm">Memory Champion: Updated memory 50+ times</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}