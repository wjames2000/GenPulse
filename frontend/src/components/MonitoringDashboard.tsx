import React, { useState, useEffect, useRef } from 'react';
import { 
  GitBranch, 
  Ruler, 
  Terminal, 
  Layout, 
  Clock, 
  Search, 
  Bolt, 
  Hourglass, 
  Copy, 
  Send, 
  Brain,
  CheckCircle2,
  AlertCircle,
  RefreshCw,
  Play,
  Pause,
  AlertTriangle,
  FileText,
  DollarSign,
  TrendingUp,
  Code,
  Cpu,
  Database,
  Network,
  BarChart3,
  Activity,
  Shield,
  Zap,
  Users,
  Settings,
  Eye,
  EyeOff,
  Maximize2,
  Minimize2,
  Filter,
  Download,
  Trash2,
  Edit,
  Plus,
  Minus,
  ChevronRight,
  ChevronLeft,
  ChevronUp,
  ChevronDown,
  Star,
  Heart,
  Target,
  Server,
  HardDrive,
  MemoryStick,
  MessageSquare,
  Mail,
  User,
  CreditCard,
  Wallet,
  Coins,
  ShoppingCart,
  Package,
  Box,
  Tag,
  Hash,
  List,
  Bell
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';
import { Agent, LogEntry, TimelineEvent, Thought, ToolInvocation, CostMetric, EvolutionEvent, FileDiff } from '../types';
import { api } from '../services/api';
import AgentStatusDashboard from './monitoring/AgentStatusDashboard';
import ExecutionTimeline from './monitoring/ExecutionTimeline';
import ThoughtChainPanel from './monitoring/ThoughtChainPanel';
import ToolLogsTable from './monitoring/ToolLogsTable';
import TerminalOutputPanel from './monitoring/TerminalOutputPanel';
import FileDiffPreview from './monitoring/FileDiffPreview';
import CostDashboard from './monitoring/CostDashboard';
import EvolutionMonitor from './monitoring/EvolutionMonitor';

export default function MonitoringDashboard() {
  const [activeView, setActiveView] = useState<'overview' | 'agents' | 'timeline' | 'thoughts' | 'tools' | 'terminal' | 'files' | 'cost' | 'evolution'>('overview');
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [isAutoRefresh, setIsAutoRefresh] = useState(true);
  const [pipelineStatus, setPipelineStatus] = useState<'idle' | 'active' | 'completed' | 'error'>('idle');
  const [agents, setAgents] = useState<Agent[]>([]);
  const [timelineEvents, setTimelineEvents] = useState<TimelineEvent[]>([]);
  const [thoughts, setThoughts] = useState<Thought[]>([]);
  const [toolLogs, setToolLogs] = useState<ToolInvocation[]>([]);
  const [terminalOutput, setTerminalOutput] = useState<string[]>([]);
  const [fileDiffs, setFileDiffs] = useState<FileDiff[]>([]);
  const [costMetrics, setCostMetrics] = useState<CostMetric[]>([]);
  const [evolutionEvents, setEvolutionEvents] = useState<EvolutionEvent[]>([]);
  const [stats, setStats] = useState({
    activeAgents: 0,
    totalAgents: 0,
    successRate: 0,
    uptime: 0,
    totalExecutions: 0,
    avgResponseTime: 0,
    tokenUsage: 0,
    toolCalls: 0,
    filesChanged: 0,
    costToday: 0,
    skillsGenerated: 0
  });
  
  const refreshIntervalRef = useRef<NodeJS.Timeout>();

  const fetchMonitoringData = async () => {
    try {
      setIsRefreshing(true);
      
      // 并行获取所有监控数据
      const [
        agentsData,
        timelineData,
        thoughtsData,
        toolLogsData,
        terminalData,
        fileDiffsData,
        costData,
        evolutionData,
        statsData
      ] = await Promise.all([
        api.getAgentsStatus(),
        api.getTimelineEvents(),
        api.getThoughts(),
        api.getToolLogs(),
        api.getTerminalOutput(),
        api.getFileDiffs(),
        api.getCostMetrics(),
        api.getEvolutionEvents(),
        api.getMonitoringStats()
      ]);
      
      setAgents(agentsData);
      setTimelineEvents(timelineData);
      setThoughts(thoughtsData);
      setToolLogs(toolLogsData);
      setTerminalOutput(terminalData);
      setFileDiffs(fileDiffsData);
      setCostMetrics(costData);
      setEvolutionEvents(evolutionData);
      setStats(statsData);
      
      // 更新管道状态
      const activeCount = agentsData.filter(a => a.status === 'active').length;
      const completedCount = agentsData.filter(a => a.status === 'completed').length;
      
      if (activeCount > 0) {
        setPipelineStatus('active');
      } else if (completedCount === agentsData.length) {
        setPipelineStatus('completed');
      } else {
        setPipelineStatus('idle');
      }
      
    } catch (error) {
      console.error('Failed to fetch monitoring data:', error);
      setPipelineStatus('error');
    } finally {
      setIsRefreshing(false);
    }
  };

  const handleRefresh = async () => {
    await fetchMonitoringData();
  };

  const handleStartPipeline = async () => {
    try {
      await api.startPipeline('current-project');
      await fetchMonitoringData();
    } catch (error) {
      console.error('Failed to start pipeline:', error);
    }
  };

  const handleIntervene = async (message: string) => {
    try {
      await api.sendIntervention(message);
      await fetchMonitoringData();
    } catch (error) {
      console.error('Failed to send intervention:', error);
    }
  };

  useEffect(() => {
    // 初始加载
    fetchMonitoringData();

    // 设置自动刷新
    if (isAutoRefresh) {
      refreshIntervalRef.current = setInterval(fetchMonitoringData, 3000); // 每3秒刷新一次
    }

    return () => {
      if (refreshIntervalRef.current) {
        clearInterval(refreshIntervalRef.current);
      }
    };
  }, [isAutoRefresh]);

  useEffect(() => {
    // 更新自动刷新定时器
    if (refreshIntervalRef.current) {
      clearInterval(refreshIntervalRef.current);
    }
    
    if (isAutoRefresh) {
      refreshIntervalRef.current = setInterval(fetchMonitoringData, 3000);
    }
  }, [isAutoRefresh]);

  const renderActiveView = () => {
    switch (activeView) {
      case 'overview':
        return (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <div className="col-span-2">
              <AgentStatusDashboard 
                agents={agents}
                stats={stats}
                onRefresh={fetchMonitoringData}
                isRefreshing={isRefreshing}
              />
            </div>
            <div>
              <ExecutionTimeline events={timelineEvents} />
            </div>
            <div>
              <ThoughtChainPanel thoughts={thoughts} />
            </div>
            <div className="col-span-2">
              <ToolLogsTable logs={toolLogs} />
            </div>
          </div>
        );
      case 'agents':
        return <AgentStatusDashboard agents={agents} stats={stats} onRefresh={fetchMonitoringData} isRefreshing={isRefreshing} />;
      case 'timeline':
        return <ExecutionTimeline events={timelineEvents} />;
      case 'thoughts':
        return <ThoughtChainPanel thoughts={thoughts} />;
      case 'tools':
        return <ToolLogsTable logs={toolLogs} />;
      case 'terminal':
        return <TerminalOutputPanel output={terminalOutput} />;
      case 'files':
        return <FileDiffPreview diffs={fileDiffs} />;
      case 'cost':
        return <CostDashboard metrics={costMetrics} />;
      case 'evolution':
        return <EvolutionMonitor events={evolutionEvents} />;
      default:
        return null;
    }
  };

  return (
    <div className="flex flex-1 h-full overflow-hidden bg-[#0A0A0A]">
      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <div className="border-b border-white/10 p-6">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-black uppercase tracking-tighter mb-2">
                Monitoring Dashboard
              </h1>
              <p className="text-sm text-white/60">
                Real-time monitoring of AI agents, execution traces, and system metrics
              </p>
            </div>
            
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                <div className={cn(
                  "w-3 h-3 rounded-full",
                  pipelineStatus === 'active' ? "bg-primary animate-pulse" :
                  pipelineStatus === 'completed' ? "bg-green-500" :
                  pipelineStatus === 'error' ? "bg-red-500" :
                  "bg-white/20"
                )} />
                <span className="text-sm font-medium">
                  {pipelineStatus === 'active' ? 'Pipeline Active' :
                   pipelineStatus === 'completed' ? 'Pipeline Completed' :
                   pipelineStatus === 'error' ? 'Pipeline Error' :
                   'Pipeline Idle'}
                </span>
              </div>
              
              <button
                onClick={() => setIsAutoRefresh(!isAutoRefresh)}
                className={cn(
                  "px-4 py-2 text-xs font-medium rounded-lg border transition-colors flex items-center gap-2",
                  isAutoRefresh 
                    ? "bg-primary/20 text-primary border-primary/30" 
                    : "bg-white/5 text-white/60 border-white/10 hover:bg-white/10"
                )}
              >
                {isAutoRefresh ? <Pause size={14} /> : <Play size={14} />}
                {isAutoRefresh ? 'Auto Refresh ON' : 'Auto Refresh OFF'}
              </button>
              
              <button
                onClick={handleRefresh}
                disabled={isRefreshing}
                className="px-4 py-2 text-xs font-medium bg-white/5 text-white/60 border border-white/10 rounded-lg hover:bg-white/10 transition-colors flex items-center gap-2"
              >
                <RefreshCw size={14} className={cn(isRefreshing && "animate-spin")} />
                Refresh
              </button>
              
              <button
                onClick={handleStartPipeline}
                className={cn(
                  "px-6 py-2 text-sm font-medium rounded-lg transition-all flex items-center gap-2",
                  pipelineStatus === 'active'
                    ? "bg-primary/20 text-primary border border-primary/30"
                    : "bg-primary text-black hover:scale-105"
                )}
              >
                {pipelineStatus === 'active' ? (
                  <>
                    <AlertTriangle size={16} />
                    Pipeline Running
                  </>
                ) : (
                  <>
                    <Play size={16} />
                    Start Pipeline
                  </>
                )}
              </button>
            </div>
          </div>
          
          {/* Quick Stats */}
          <div className="grid grid-cols-6 gap-4 mt-6">
            <div className="bg-white/5 border border-white/10 rounded-lg p-4">
              <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Active Agents</div>
              <div className="text-2xl font-bold">{stats.activeAgents}/{stats.totalAgents}</div>
            </div>
            <div className="bg-white/5 border border-white/10 rounded-lg p-4">
              <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Success Rate</div>
              <div className="text-2xl font-bold text-green-500">{stats.successRate}%</div>
            </div>
            <div className="bg-white/5 border border-white/10 rounded-lg p-4">
              <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Tool Calls</div>
              <div className="text-2xl font-bold">{stats.toolCalls}</div>
            </div>
            <div className="bg-white/5 border border-white/10 rounded-lg p-4">
              <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Token Usage</div>
              <div className="text-2xl font-bold">{stats.tokenUsage.toLocaleString()}</div>
            </div>
            <div className="bg-white/5 border border-white/10 rounded-lg p-4">
              <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Files Changed</div>
              <div className="text-2xl font-bold">{stats.filesChanged}</div>
            </div>
            <div className="bg-white/5 border border-white/10 rounded-lg p-4">
              <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Cost Today</div>
              <div className="text-2xl font-bold">${stats.costToday.toFixed(2)}</div>
            </div>
          </div>
        </div>
        
        {/* Navigation Tabs */}
        <div className="border-b border-white/10">
          <div className="flex overflow-x-auto">
            {[
              { id: 'overview', label: 'Overview', icon: <Activity size={16} /> },
              { id: 'agents', label: 'Agents', icon: <Users size={16} /> },
              { id: 'timeline', label: 'Timeline', icon: <Clock size={16} /> },
              { id: 'thoughts', label: 'Thoughts', icon: <Brain size={16} /> },
              { id: 'tools', label: 'Tools', icon: <Settings size={16} /> },
              { id: 'terminal', label: 'Terminal', icon: <Terminal size={16} /> },
              { id: 'files', label: 'Files', icon: <FileText size={16} /> },
              { id: 'cost', label: 'Cost', icon: <DollarSign size={16} /> },
              { id: 'evolution', label: 'Evolution', icon: <TrendingUp size={16} /> }
            ].map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveView(tab.id as any)}
                className={cn(
                  "px-6 py-4 text-sm font-medium border-b-2 transition-colors flex items-center gap-2 whitespace-nowrap",
                  activeView === tab.id
                    ? "border-primary text-primary"
                    : "border-transparent text-white/60 hover:text-white"
                )}
              >
                {tab.icon}
                {tab.label}
              </button>
            ))}
          </div>
        </div>
        
        {/* Content Area */}
        <div className="flex-1 overflow-y-auto p-6">
          {renderActiveView()}
        </div>
      </div>
      
      {/* Sidebar - Intervention Panel */}
      <div className="w-80 border-l border-white/10 bg-black/40">
        <div className="p-6 border-b border-white/10">
          <h3 className="text-lg font-bold mb-2">Agent Intervention</h3>
          <p className="text-sm text-white/60">
            Send instructions or feedback to active agents
          </p>
        </div>
        
        <div className="p-6">
          <form 
            onSubmit={async (e) => {
              e.preventDefault();
              const form = e.target as HTMLFormElement;
              const input = form.querySelector('input') as HTMLInputElement;
              if (input.value.trim()) {
                await handleIntervene(input.value);
                input.value = '';
              }
            }}
          >
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">Message</label>
                <input 
                  type="text" 
                  placeholder="Enter intervention message..."
                  className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm focus:outline-none focus:border-primary transition-colors"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-2">Target Agent</label>
                <select className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm focus:outline-none focus:border-primary transition-colors">
                  <option value="">All Active Agents</option>
                  {agents.filter(a => a.status === 'active').map(agent => (
                    <option key={agent.id} value={agent.id}>{agent.name} ({agent.role})</option>
                  ))}
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-2">Priority</label>
                <select className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm focus:outline-none focus:border-primary transition-colors">
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                  <option value="critical">Critical</option>
                </select>
              </div>
              
              <button
                type="submit"
                className="w-full bg-primary text-black font-medium py-3 rounded-lg hover:bg-primary/90 transition-colors flex items-center justify-center gap-2"
              >
                <Send size={16} />
                Send Intervention
              </button>
            </div>
          </form>
        </div>
        
        <div className="p-6 border-t border-white/10">
          <h4 className="text-sm font-bold mb-4">Recent Interventions</h4>
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-white/5 border border-white/10 rounded-lg p-3">
                <div className="flex justify-between items-start mb-1">
                  <span className="text-xs font-medium">User Intervention</span>
                  <span className="text-xs text-white/40">2m ago</span>
                </div>
                <p className="text-sm text-white/80">Please optimize the database queries for better performance.</p>
                <div className="flex items-center gap-2 mt-2">
                  <span className="text-xs px-2 py-1 bg-primary/20 text-primary rounded">To: Architect</span>
                  <span className="text-xs px-2 py-1 bg-yellow-500/20 text-yellow-500 rounded">Priority: High</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}