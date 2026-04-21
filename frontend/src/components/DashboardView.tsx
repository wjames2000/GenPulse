import React, { useState, useEffect, useRef } from 'react';
import { 
  GitBranch, 
  Ruler, 
  Terminal, 
  Layout, 
  Clock, 
  Search, 
  MoreVertical, 
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
  AlertTriangle
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';
import { Agent } from '../types';
import { AGENTS, TIMELINE, LOGS, THOUGHTS } from '../constants';
import { api } from '../services/api';
import LogViewer from './LogViewer';

export default function DashboardView() {
  const [activeTab, setActiveTab] = useState<'logs' | 'terminal' | 'network'>('logs');
  const [agents, setAgents] = useState<Agent[]>(AGENTS);
  const [logs, setLogs] = useState(LOGS);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [isAutoRefresh, setIsAutoRefresh] = useState(true);
  const [pipelineStatus, setPipelineStatus] = useState<'idle' | 'active' | 'completed' | 'error'>('idle');
  const [activeAgentsCount, setActiveAgentsCount] = useState(2);
  const [successRate, setSuccessRate] = useState(85);
  const [uptime, setUptime] = useState(99);
  const refreshIntervalRef = useRef<NodeJS.Timeout>();

  const fetchAgentsStatus = async () => {
    try {
      const status = await api.getAllAgentsStatus();
      
      // 更新agents状态
      const updatedAgents = agents.map(agent => {
        const agentStatus = status[agent.id];
        if (agentStatus) {
          return {
            ...agent,
            status: agentStatus.state === 'active' ? 'active' : 
                    agentStatus.state === 'waiting' ? 'waiting' : 'idle',
            progress: agentStatus.progress || agent.progress,
            currentTask: agentStatus.current_task || agent.currentTask,
            timeActive: agentStatus.time_active || agent.timeActive
          };
        }
        return agent;
      });
      
      setAgents(updatedAgents);
      
      // 计算统计数据
      const activeCount = updatedAgents.filter(a => a.status === 'active').length;
      const completedCount = updatedAgents.filter(a => a.progress === 100).length;
      
      setActiveAgentsCount(activeCount);
      setSuccessRate(completedCount > 0 ? Math.round((completedCount / updatedAgents.length) * 100) : 0);
      
      // 更新管道状态
      if (activeCount > 0) {
        setPipelineStatus('active');
      } else if (completedCount === updatedAgents.length) {
        setPipelineStatus('completed');
      } else {
        setPipelineStatus('idle');
      }
      
      // 获取最新日志
      const newLogs = await api.getLogs();
      if (Array.isArray(newLogs) && newLogs.length > 0) {
        const formattedLogs = newLogs.map((log: any) => ({
          timestamp: new Date(log.timestamp).toLocaleTimeString('en-US', { hour12: false }),
          level: log.level as 'info' | 'debug' | 'success' | 'warn' | 'error' | 'sys',
          message: log.message
        }));
        setLogs(prev => [...formattedLogs, ...prev].slice(0, 50)); // 保持最多50条日志
      }
    } catch (error) {
      console.error('Failed to fetch agents status:', error);
      setPipelineStatus('error');
    }
  };

  const handleRefresh = async () => {
    setIsRefreshing(true);
    await fetchAgentsStatus();
    setIsRefreshing(false);
  };

  const handleStartPipeline = async () => {
    try {
      await api.startPipeline('current-project');
      await fetchAgentsStatus();
    } catch (error) {
      console.error('Failed to start pipeline:', error);
    }
  };

  const handleIntervene = async (message: string) => {
    try {
      await api.logMessage('info', `User intervention: ${message}`);
      // 这里可以添加更多的干预逻辑
      const newLog = {
        timestamp: new Date().toLocaleTimeString('en-US', { hour12: false }),
        level: 'info' as const,
        message: `User intervention: ${message}`
      };
      setLogs(prev => [newLog, ...prev].slice(0, 50));
    } catch (error) {
      console.error('Failed to send intervention:', error);
    }
  };

  useEffect(() => {
    // 初始加载
    fetchAgentsStatus();

    // 设置自动刷新
    if (isAutoRefresh) {
      refreshIntervalRef.current = setInterval(fetchAgentsStatus, 5000); // 每5秒刷新一次
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
      refreshIntervalRef.current = setInterval(fetchAgentsStatus, 5000);
    }
  }, [isAutoRefresh]);

  return (
    <div className="flex flex-1 h-full overflow-hidden p-12 gap-12 bg-[#0A0A0A]">
      {/* Left Column: Dashboard & Timeline */}
      <div className="flex-1 flex flex-col gap-12 overflow-y-auto pr-4 custom-scrollbar">
        {/* Page Header */}
        <section className="relative">
          <div className="flex justify-between items-start mb-4">
            <span className="text-[10px] uppercase tracking-[0.4em] text-white/40">
              Pipeline Execution / Live Monitoring / 2024
            </span>
            <div className="flex items-center gap-4">
              <button
                onClick={() => setIsAutoRefresh(!isAutoRefresh)}
                className={cn(
                  "text-[10px] font-black uppercase tracking-widest flex items-center gap-2 transition-colors",
                  isAutoRefresh ? "text-primary" : "text-white/40 hover:text-white"
                )}
              >
                {isAutoRefresh ? <Pause size={12} /> : <Play size={12} />}
                {isAutoRefresh ? "Auto Refresh ON" : "Auto Refresh OFF"}
              </button>
              <button
                onClick={handleRefresh}
                disabled={isRefreshing}
                className="text-[10px] font-black uppercase tracking-widest text-white/40 hover:text-white transition-colors flex items-center gap-2"
              >
                <RefreshCw size={12} className={cn(isRefreshing && "animate-spin")} />
                Refresh
              </button>
            </div>
          </div>
          
          <h1 className="text-[120px] leading-[0.8] font-black tracking-tighter uppercase mb-6">
            Hyper<br/>Pulse
          </h1>
          <div className="absolute top-0 right-0 w-1/3 text-xs leading-relaxed text-white/40 pt-4 uppercase tracking-widest font-bold">
            Real-time heuristic analysis of autonomic agent pathways. Monitoring the nexus of digital intelligence.
          </div>
          
          <div className="flex items-center gap-6 mt-8">
            <div className="flex flex-col gap-1">
              <span className="text-[40px] font-black leading-none">{activeAgentsCount}</span>
              <span className="text-[10px] uppercase font-black tracking-widest text-primary">Active Agents</span>
            </div>
            <div className="w-[1px] h-12 bg-white/10" />
            <div className="flex flex-col gap-1">
              <span className="text-[40px] font-black leading-none text-primary">{successRate}%</span>
              <span className="text-[10px] uppercase font-black tracking-widest text-white/40">Success Rate</span>
            </div>
            <div className="w-[1px] h-12 bg-white/10" />
            <div className="flex flex-col gap-1 text-primary">
              <span className="text-[40px] font-black leading-none">{uptime}%</span>
              <span className="text-[10px] uppercase font-black tracking-widest text-white/40">Uptime</span>
            </div>
            
            <div className="ml-auto">
              <button
                onClick={handleStartPipeline}
                className={cn(
                  "px-8 py-3 font-black uppercase text-[10px] tracking-widest transition-all shadow-lg rounded-none flex items-center gap-3",
                  pipelineStatus === 'active' 
                    ? "bg-primary/20 text-primary border border-primary/30" 
                    : "bg-primary text-black hover:scale-105"
                )}
              >
                {pipelineStatus === 'active' ? (
                  <>
                    <AlertTriangle size={14} />
                    Pipeline Running
                  </>
                ) : (
                  <>
                    <Play size={14} />
                    Start Pipeline
                  </>
                )}
              </button>
            </div>
          </div>
        </section>

        {/* Agent Status Cards */}
        <div className="border-y border-white/10 py-12">
          <div className="flex items-center justify-between mb-8">
            <h2 className="text-xl font-black uppercase tracking-tighter flex items-center gap-4">
              <span className="w-8 h-[1px] bg-primary"></span>
              Agent Status Dashboard
            </h2>
            <div className="flex items-center gap-4">
              <div className="text-[10px] uppercase tracking-[0.3em] font-black text-white/40">
                {agents.filter(a => a.status === 'active').length} Active / {agents.length} Total
              </div>
              <div className="flex items-center gap-2">
                <div className="w-2 h-2 bg-primary rounded-full"></div>
                <span className="text-[8px] uppercase tracking-widest text-white/40">Active</span>
                <div className="w-2 h-2 bg-yellow-500 rounded-full ml-4"></div>
                <span className="text-[8px] uppercase tracking-widest text-white/40">Waiting</span>
                <div className="w-2 h-2 bg-green-500 rounded-full ml-4"></div>
                <span className="text-[8px] uppercase tracking-widest text-white/40">Completed</span>
              </div>
            </div>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-1">
            {agents.map((agent) => (
              <AgentCard key={agent.id} agent={agent} />
            ))}
          </div>

          {/* Agent Statistics */}
          <div className="grid grid-cols-4 gap-1 mt-12 pt-12 border-t border-white/10">
            <div className="p-6 border-r border-white/10 last:border-r-0">
              <div className="text-[8px] text-white/20 uppercase font-black tracking-widest mb-2">Total Executions</div>
              <div className="text-2xl font-black text-primary">1,248</div>
            </div>
            <div className="p-6 border-r border-white/10">
              <div className="text-[8px] text-white/20 uppercase font-black tracking-widest mb-2">Success Rate</div>
              <div className="text-2xl font-black text-green-500">94.2%</div>
            </div>
            <div className="p-6 border-r border-white/10">
              <div className="text-[8px] text-white/20 uppercase font-black tracking-widest mb-2">Avg Response Time</div>
              <div className="text-2xl font-black text-primary">2.4s</div>
            </div>
            <div className="p-6">
              <div className="text-[8px] text-white/20 uppercase font-black tracking-widest mb-2">Token Usage</div>
              <div className="text-2xl font-black text-primary">12.8K</div>
            </div>
          </div>
        </div>

        {/* Execution Timeline */}
        <div className="border border-white/10 p-10 mt-4 bg-white/[0.02]">
          <div className="flex items-center justify-between mb-10">
            <h2 className="text-xl font-black uppercase tracking-tighter flex items-center gap-4">
              <span className="w-8 h-[1px] bg-primary"></span>
              Execution Timeline
            </h2>
            <div className="text-[10px] uppercase tracking-[0.3em] font-black text-white/40">
              Synchronizing Nodes...
            </div>
          </div>
          
          <div className="space-y-8">
            {['Orchestrator', 'Architect', 'Backend Dev'].map((agentName) => (
              <div key={agentName} className="relative h-12 flex items-center group">
                <div className="w-40 text-[10px] font-black uppercase tracking-widest text-white/40 group-hover:text-white transition-colors">{agentName}</div>
                <div className="flex-1 bg-white/5 h-2 rounded-none relative">
                  {TIMELINE.filter(e => e.agent === agentName).map(event => (
                    <div 
                      key={event.id}
                      className={cn(
                        "absolute top-1/2 -translate-y-1/2 h-4 flex items-center px-4 text-[9px] font-black uppercase tracking-widest border transition-all truncate",
                        event.isComplete 
                          ? "bg-primary text-black border-transparent" 
                          : "bg-transparent text-primary border-primary animate-pulse"
                      )}
                      style={{ left: event.offset, width: event.width }}
                    >
                      {event.action}
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Lower Panel: Tabs */}
        <div className="border border-white/10 flex flex-col flex-1 min-h-[400px]">
          <div className="flex border-b border-white/10">
            {(['logs', 'terminal', 'network', 'agents'] as const).map(tab => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={cn(
                  "px-10 py-5 text-[10px] font-black uppercase tracking-[0.3em] transition-all relative border-r border-white/10",
                  activeTab === tab ? "text-black bg-primary" : "text-white/40 hover:text-white hover:bg-white/5"
                )}
              >
                {tab === 'logs' && "System Logs"}
                {tab === 'terminal' && "Terminal Out"}
                {tab === 'network' && "Network Nodes"}
                {tab === 'agents' && "Agent Details"}
              </button>
            ))}
          </div>
          
          <div className="flex-1 overflow-hidden">
            {activeTab === 'logs' && (
              <div className="h-full">
                <LogViewer 
                  compact={true}
                  autoRefresh={isAutoRefresh}
                  showFilters={false}
                  className="h-full border-0 rounded-none"
                />
              </div>
            )}
            
            {activeTab === 'terminal' && (
              <div className="p-8 flex-1 font-mono text-[11px] space-y-3 custom-scrollbar uppercase tracking-tight overflow-y-auto max-h-[300px]">
                <div className="flex gap-6 items-baseline opacity-60 hover:opacity-100 transition-opacity">
                  <span className="text-white/20 w-24 shrink-0 font-black tracking-widest">[14:32:45]</span>
                  <span className="font-black w-20 shrink-0 text-primary">INFO</span>
                  <span className="text-white/80 flex-1">Starting pipeline execution...</span>
                </div>
                <div className="flex gap-6 items-baseline opacity-60 hover:opacity-100 transition-opacity">
                  <span className="text-white/20 w-24 shrink-0 font-black tracking-widest">[14:32:46]</span>
                  <span className="font-black w-20 shrink-0 text-primary">INFO</span>
                  <span className="text-white/80 flex-1">Orchestrator agent initialized project structure</span>
                </div>
                <div className="flex gap-6 items-baseline opacity-60 hover:opacity-100 transition-opacity">
                  <span className="text-white/20 w-24 shrink-0 font-black tracking-widest">[14:32:47]</span>
                  <span className="font-black w-20 shrink-0 text-green-500">SUCCESS</span>
                  <span className="text-white/80 flex-1">Architect agent generated technical design</span>
                </div>
                <div className="flex gap-6 items-baseline opacity-60 hover:opacity-100 transition-opacity">
                  <span className="text-white/20 w-24 shrink-0 font-black tracking-widest">[14:32:48]</span>
                  <span className="font-black w-20 shrink-0 text-primary">INFO</span>
                  <span className="text-white/80 flex-1">Frontend agent creating React components...</span>
                </div>
                <div className="flex gap-6 items-baseline opacity-60 hover:opacity-100 transition-opacity">
                  <span className="text-white/20 w-24 shrink-0 font-black tracking-widest">[14:32:49]</span>
                  <span className="font-black w-20 shrink-0 text-primary">INFO</span>
                  <span className="text-white/80 flex-1">Backend agent generating API endpoints...</span>
                </div>
                <div className="flex gap-6 items-baseline opacity-60 hover:opacity-100 transition-opacity">
                  <span className="text-white/20 w-24 shrink-0 font-black tracking-widest">[14:32:50]</span>
                  <span className="font-black w-20 shrink-0 text-yellow-500">WARN</span>
                  <span className="text-white/80 flex-1">QA agent detected potential performance issue</span>
                </div>
                <div className="flex gap-6 items-baseline opacity-60 hover:opacity-100 transition-opacity">
                  <span className="text-white/20 w-24 shrink-0 font-black tracking-widest">[14:32:51]</span>
                  <span className="font-black w-20 shrink-0 text-green-500">SUCCESS</span>
                  <span className="text-white/80 flex-1">DevOps agent completed deployment configuration</span>
                </div>
                <div className="flex gap-6 items-baseline opacity-60 hover:opacity-100 transition-opacity">
                  <span className="text-white/20 w-24 shrink-0 font-black tracking-widest">[14:32:52]</span>
                  <span className="font-black w-20 shrink-0 text-primary">INFO</span>
                  <span className="text-white/80 flex-1">Reviewer agent completed code review</span>
                </div>
                <div className="flex gap-6 items-baseline opacity-60 hover:opacity-100 transition-opacity">
                  <span className="text-white/20 w-24 shrink-0 font-black tracking-widest">[14:32:53]</span>
                  <span className="font-black w-20 shrink-0 text-green-500">SUCCESS</span>
                  <span className="text-white/80 flex-1">Pipeline execution completed successfully!</span>
                </div>
              </div>
            )}
            
            {activeTab === 'network' && (
              <div className="p-8">
                <div className="text-center text-white/40 py-8">
                  <div className="text-lg mb-2">Network Nodes View</div>
                  <div className="text-sm">Visualization of agent communication network</div>
                </div>
              </div>
            )}
            
            {activeTab === 'agents' && (
              <div className="p-8">
                <div className="grid grid-cols-2 gap-4">
                  {agents.map(agent => (
                    <div key={agent.id} className="border border-white/10 p-4 hover:bg-white/5 transition-colors">
                      <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center gap-3">
                          <div className={cn(
                            "w-2 h-2 rounded-full",
                            agent.status === 'active' ? "bg-primary" : 
                            agent.status === 'waiting' ? "bg-yellow-500" :
                            agent.status === 'completed' ? "bg-green-500" :
                            "bg-white/20"
                          )} />
                          <span className="text-sm font-black uppercase tracking-wider">{agent.name}</span>
                        </div>
                        <span className="text-[10px] text-white/40 uppercase tracking-widest">{agent.role}</span>
                      </div>
                      <div className="text-xs text-white/60 mb-2">{agent.currentTask || "Idle"}</div>
                      <div className="flex items-center justify-between text-[10px] text-white/40">
                        <span>Progress: {agent.progress}%</span>
                        <span>Time: {agent.timeActive}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Right Column: Thought Bubble (Glass Panel) */}
      <aside className="w-[400px] shrink-0 glass-panel border-l border-white/10 flex flex-col relative">
        <div className="p-10 border-b border-white/10">
          <div className="text-[10px] uppercase font-black tracking-[0.5em] text-white/40 mb-2">Cognitive Engine</div>
          <h3 className="text-3xl font-black uppercase tracking-tighter text-primary">Architect Node</h3>
        </div>

        <div className="flex-1 overflow-y-auto p-10 space-y-12 custom-scrollbar">
          {THOUGHTS.map((thought, i) => (
            <motion.div 
              key={i}
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: i * 0.1 }}
            >
              {thought.isCode ? (
                <div className="border border-white/10 overflow-hidden">
                  <div className="bg-white/5 px-6 py-3 text-[9px] font-black uppercase tracking-widest text-white/40 flex items-center justify-between">
                    <span className="flex items-center gap-3">
                      <Terminal size={12} />
                      {thought.filename}
                    </span>
                    <Copy size={12} className="cursor-pointer hover:text-primary transition-colors" />
                  </div>
                  <div className="p-6 font-mono text-[10px] leading-relaxed overflow-x-auto text-white/60 bg-black/40">
                    <pre><code>{thought.code}</code></pre>
                  </div>
                </div>
              ) : (
                <div className="relative pl-8 border-l-2 border-white/10 py-2">
                  <div className="text-[9px] text-primary mb-4 flex items-center gap-3 uppercase font-black tracking-widest">
                    {thought.type === 'internal' ? <Bolt size={12} /> : <Hourglass size={12} />}
                    {thought.type === 'internal' ? "Heuristic thought" : "Neural formulation"}
                  </div>
                  <p className="text-sm font-medium text-white/80 leading-[1.6]">
                    {thought.content}
                  </p>
                </div>
              )}
            </motion.div>
          ))}
        </div>

        {/* Input Area */}
        <div className="p-10 bg-black/40 border-t border-white/10">
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
            <div className="relative">
              <input 
                type="text" 
                placeholder="INTERVENE IN ARCHITECTURE..." 
                className="w-full bg-white/5 border border-white/10 py-5 px-6 text-[10px] font-black tracking-widest text-white focus:bg-white/10 outline-none transition-all placeholder:text-white/20 uppercase"
              />
              <button 
                type="submit"
                className="absolute right-6 top-1/2 -translate-y-1/2 text-primary hover:scale-125 transition-all"
              >
                <Send size={18} />
              </button>
            </div>
          </form>
        </div>

        <div className="absolute left-0 top-1/2 -translate-y-1/2 flex items-center gap-4 origin-left -rotate-90 -ml-16 pointer-events-none">
          <div className="w-12 h-[1px] bg-white/10"></div>
          <span className="text-[8px] uppercase tracking-[0.5em] text-white/20">Agent Consciousness</span>
        </div>
      </aside>
    </div>
  );
}

interface AgentCardProps {
  agent: Agent;
  key?: string;
}

function AgentCard({ agent }: AgentCardProps) {
  const Icon = agent.type === 'orchestrator' ? GitBranch : 
               agent.type === 'architect' ? Ruler : 
               agent.type === 'backend' ? Terminal : 
               agent.type === 'frontend' ? Layout :
               agent.type === 'qa' ? CheckCircle2 : Brain;
  
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'text-primary';
      case 'waiting': return 'text-yellow-500';
      default: return 'text-white/40';
    }
  };

  const getStatusBgColor = (status: string) => {
    switch (status) {
      case 'active': return 'bg-primary/20 border-primary/30';
      case 'waiting': return 'bg-yellow-500/20 border-yellow-500/30';
      case 'completed': return 'bg-green-500/20 border-green-500/30';
      case 'error': return 'bg-red-500/20 border-red-500/30';
      default: return 'bg-white/5 border-white/10';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'active': return 'ACTIVE';
      case 'waiting': return 'WAITING';
      case 'completed': return 'COMPLETED';
      case 'error': return 'ERROR';
      default: return 'IDLE';
    }
  };

  return (
    <div className={cn(
      "p-8 transition-all duration-500 border-r border-white/10 last:border-r-0 hover:bg-white/[0.03] group relative",
      agent.status === 'idle' && "opacity-40 grayscale"
    )}>
      {/* Status indicator */}
      <div className="absolute top-4 right-4">
        <div className={cn(
          "w-3 h-3 rounded-none rotate-45 border transition-all",
          agent.status === 'active' ? "bg-primary border-primary shadow-[0_0_15px_#FBDF24]" : 
          agent.status === 'waiting' ? "bg-yellow-500 border-yellow-500" :
          "border-white/20"
        )} />
      </div>

      <div className="flex justify-between items-start mb-10">
        <div className="flex flex-col gap-1">
          <div className="flex items-center gap-3">
            <Icon size={16} className={cn(
              "transition-colors",
              agent.status === 'active' ? "text-primary" : 
              agent.status === 'waiting' ? "text-yellow-500" :
              "text-white/40"
            )} />
            <span className="text-[9px] text-white/40 uppercase font-black tracking-widest">{agent.role}</span>
          </div>
          <span className="text-lg font-black uppercase tracking-tighter mt-2">{agent.name}</span>
        </div>
      </div>

      <div className="space-y-8">
        {/* Status badge */}
        <div className={cn(
          "inline-flex items-center px-3 py-1 text-[8px] font-black uppercase tracking-widest border transition-all",
          getStatusBgColor(agent.status),
          getStatusColor(agent.status)
        )}>
          {getStatusText(agent.status)}
        </div>

        <div className="flex flex-col gap-2">
          <span className="text-[8px] text-white/20 uppercase font-black tracking-widest">Active Task</span>
          <div className="text-[11px] font-black uppercase tracking-wide text-white/60 leading-tight line-clamp-2 h-8">
            {agent.currentTask || "No active task"}
          </div>
        </div>
        
        <div className="grid grid-cols-2 gap-4">
          <div className="flex flex-col">
            <span className="text-[8px] text-white/20 uppercase font-black tracking-widest">Elapsed</span>
            <span className="text-[10px] font-mono font-black text-primary mt-1">{agent.timeActive || "00:00"}</span>
          </div>
          <div className="flex flex-col">
            <span className="text-[8px] text-white/20 uppercase font-black tracking-widest">Progress</span>
            <span className="text-[10px] font-mono font-black text-primary mt-1">{agent.progress}%</span>
          </div>
        </div>
        
        <div className="w-full h-[2px] bg-white/5 relative">
          <motion.div 
            initial={{ width: 0 }}
            animate={{ width: `${agent.progress}%` }}
            className={cn(
              "h-full transition-all",
              agent.status === 'active' ? "bg-primary" : 
              agent.status === 'waiting' ? "bg-yellow-500" :
              "bg-white/20"
            )}
          />
        </div>

        {/* Additional metrics */}
        <div className="grid grid-cols-2 gap-3 pt-4 border-t border-white/10">
          <div className="flex flex-col">
            <span className="text-[7px] text-white/20 uppercase font-black tracking-widest">Tasks</span>
            <span className="text-[9px] font-mono font-black text-white/60 mt-1">12</span>
          </div>
          <div className="flex flex-col">
            <span className="text-[7px] text-white/20 uppercase font-black tracking-widest">Success</span>
            <span className="text-[9px] font-mono font-black text-green-500 mt-1">92%</span>
          </div>
        </div>
      </div>
    </div>
  );
}
