import React from 'react';
import { 
  GitBranch, 
  Ruler, 
  Terminal, 
  Layout, 
  CheckCircle2,
  Brain,
  AlertCircle,
  Clock,
  TrendingUp,
  BarChart3,
  Activity,
  Zap,
  Cpu,
  Database,
  Server,
  HardDrive,
  MemoryStick,
  Network,
  Shield,
  Users,
  Settings,
  Eye,
  Filter,
  Download,
  Maximize2,
  Minimize2
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../../utils';
import { Agent } from '../../types';

interface AgentStatusDashboardProps {
  agents: Agent[];
  stats: {
    activeAgents: number;
    totalAgents: number;
    successRate: number;
    uptime: number;
    totalExecutions: number;
    avgResponseTime: number;
    tokenUsage: number;
    toolCalls: number;
    filesChanged: number;
    costToday: number;
    skillsGenerated: number;
  };
  onRefresh: () => void;
  isRefreshing: boolean;
}

export default function AgentStatusDashboard({ 
  agents, 
  stats, 
  onRefresh, 
  isRefreshing 
}: AgentStatusDashboardProps) {
  const getAgentIcon = (type: string) => {
    switch (type) {
      case 'orchestrator': return GitBranch;
      case 'architect': return Ruler;
      case 'backend': return Terminal;
      case 'frontend': return Layout;
      case 'qa': return CheckCircle2;
      case 'reviewer': return Shield;
      case 'devops': return Server;
      case 'product': return Users;
      default: return Brain;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'text-primary';
      case 'waiting': return 'text-yellow-500';
      case 'completed': return 'text-green-500';
      case 'error': return 'text-red-500';
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
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-3">
            <Users size={24} />
            Agent Status Dashboard
          </h2>
          <p className="text-sm text-white/60 mt-1">
            Real-time monitoring of all AI agents in the pipeline
          </p>
        </div>
        
        <div className="flex items-center gap-3">
          <div className="text-sm text-white/40">
            {stats.activeAgents} Active / {stats.totalAgents} Total
          </div>
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 bg-primary rounded-full"></div>
            <span className="text-xs text-white/40">Active</span>
            <div className="w-2 h-2 bg-yellow-500 rounded-full ml-3"></div>
            <span className="text-xs text-white/40">Waiting</span>
            <div className="w-2 h-2 bg-green-500 rounded-full ml-3"></div>
            <span className="text-xs text-white/40">Completed</span>
          </div>
        </div>
      </div>

      {/* Agent Cards Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {agents.map((agent) => {
          const Icon = getAgentIcon(agent.type);
          return (
            <div
              key={agent.id}
              className={cn(
                "border rounded-xl p-5 transition-all hover:scale-[1.02]",
                getStatusBgColor(agent.status),
                agent.status === 'idle' && "opacity-60"
              )}
            >
              <div className="flex justify-between items-start mb-4">
                <div className="flex items-center gap-3">
                  <div className={cn(
                    "p-2 rounded-lg",
                    agent.status === 'active' ? "bg-primary/20" :
                    agent.status === 'waiting' ? "bg-yellow-500/20" :
                     agent.status === 'completed' || agent.status === 'error' ? "bg-green-500/20" :
                    "bg-white/10"
                  )}>
                    <Icon size={20} className={getStatusColor(agent.status)} />
                  </div>
                  <div>
                    <div className="font-bold text-lg">{agent.name}</div>
                    <div className="text-xs text-white/60 uppercase tracking-wider">{agent.role}</div>
                  </div>
                </div>
                
                <div className={cn(
                  "px-2 py-1 text-xs font-bold uppercase rounded",
                  getStatusColor(agent.status),
                  getStatusBgColor(agent.status)
                )}>
                  {getStatusText(agent.status)}
                </div>
              </div>

              <div className="space-y-4">
                {/* Current Task */}
                <div>
                  <div className="text-xs text-white/40 mb-1">Current Task</div>
                  <div className="text-sm font-medium line-clamp-2 h-10">
                    {agent.currentTask || "No active task"}
                  </div>
                </div>

                {/* Progress Bar */}
                <div>
                  <div className="flex justify-between text-xs mb-1">
                    <span className="text-white/60">Progress</span>
                    <span className="font-bold">{agent.progress}%</span>
                  </div>
                  <div className="h-2 bg-white/10 rounded-full overflow-hidden">
                    <motion.div
                      initial={{ width: 0 }}
                      animate={{ width: `${agent.progress}%` }}
                      className={cn(
                        "h-full rounded-full transition-all",
                        agent.status === 'active' ? "bg-primary" :
                        agent.status === 'waiting' ? "bg-yellow-500" :
                         agent.status === 'completed' || agent.status === 'error' ? "bg-green-500" :
                        "bg-white/20"
                      )}
                    />
                  </div>
                </div>

                {/* Metrics */}
                <div className="grid grid-cols-2 gap-3 pt-3 border-t border-white/10">
                  <div>
                    <div className="text-xs text-white/40">Time Active</div>
                    <div className="text-sm font-bold flex items-center gap-1">
                      <Clock size={12} />
                      {agent.timeActive || "00:00"}
                    </div>
                  </div>
                  <div>
                    <div className="text-xs text-white/40">Tasks</div>
                    <div className="text-sm font-bold">12</div>
                  </div>
                </div>

                {/* Additional Info */}
                <div className="flex items-center justify-between text-xs text-white/40">
                  <div className="flex items-center gap-1">
                    <Cpu size={12} />
                    <span>CPU: 24%</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <MemoryStick size={12} />
                    <span>Mem: 128MB</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <Network size={12} />
                    <span>Net: 12KB/s</span>
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Detailed Statistics */}
      <div className="border border-white/10 rounded-xl p-6 bg-white/[0.02]">
        <h3 className="text-lg font-bold mb-6 flex items-center gap-3">
          <BarChart3 size={20} />
          Agent Performance Statistics
        </h3>
        
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Total Executions</div>
            <div className="text-2xl font-bold">{stats.totalExecutions.toLocaleString()}</div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Success Rate</div>
            <div className="text-2xl font-bold text-green-500">{stats.successRate}%</div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Avg Response Time</div>
            <div className="text-2xl font-bold">{stats.avgResponseTime.toFixed(1)}s</div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Token Usage</div>
            <div className="text-2xl font-bold">{(stats.tokenUsage / 1000).toFixed(1)}K</div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Tool Calls</div>
            <div className="text-2xl font-bold">{stats.toolCalls}</div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Uptime</div>
            <div className="text-2xl font-bold">{stats.uptime}%</div>
          </div>
        </div>

        {/* Performance Chart Placeholder */}
        <div className="mt-6 p-4 border border-white/10 rounded-lg">
          <div className="flex justify-between items-center mb-4">
            <div className="text-sm font-medium">Performance Over Time</div>
            <div className="flex items-center gap-2">
              <button className="text-xs px-3 py-1 bg-white/5 rounded hover:bg-white/10 transition-colors">
                1H
              </button>
              <button className="text-xs px-3 py-1 bg-white/5 rounded hover:bg-white/10 transition-colors">
                24H
              </button>
              <button className="text-xs px-3 py-1 bg-primary/20 text-primary rounded">
                7D
              </button>
            </div>
          </div>
          
          <div className="h-40 flex items-end gap-1">
            {Array.from({ length: 24 }).map((_, i) => (
              <div
                key={i}
                className="flex-1 bg-primary/40 rounded-t hover:bg-primary transition-colors"
                style={{ height: `${20 + Math.sin(i / 3) * 30 + Math.random() * 20}%` }}
                title={`Hour ${i}: ${Math.round(20 + Math.sin(i / 3) * 30 + Math.random() * 20)}%`}
              />
            ))}
          </div>
          
          <div className="flex justify-between text-xs text-white/40 mt-2">
            <span>00:00</span>
            <span>06:00</span>
            <span>12:00</span>
            <span>18:00</span>
            <span>24:00</span>
          </div>
        </div>
      </div>

      {/* Agent Distribution */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="border border-white/10 rounded-xl p-6">
          <h4 className="text-lg font-bold mb-4">Agent Distribution by Type</h4>
          <div className="space-y-4">
            {[
              { type: 'orchestrator', count: agents.filter(a => a.type === 'orchestrator').length, color: 'bg-primary' },
              { type: 'architect', count: agents.filter(a => a.type === 'architect').length, color: 'bg-blue-500' },
              { type: 'backend', count: agents.filter(a => a.type === 'backend').length, color: 'bg-green-500' },
              { type: 'frontend', count: agents.filter(a => a.type === 'frontend').length, color: 'bg-purple-500' },
              { type: 'qa', count: agents.filter(a => a.type === 'qa').length, color: 'bg-yellow-500' },
              { type: 'devops', count: agents.filter(a => a.type === 'devops').length, color: 'bg-red-500' },
            ].map((item) => (
              <div key={item.type} className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className={`w-3 h-3 ${item.color} rounded-full`} />
                  <span className="text-sm capitalize">{item.type}</span>
                </div>
                <div className="flex items-center gap-4">
                  <div className="w-32 h-2 bg-white/10 rounded-full overflow-hidden">
                    <div 
                      className={`h-full ${item.color} rounded-full`}
                      style={{ width: `${(item.count / agents.length) * 100}%` }}
                    />
                  </div>
                  <span className="text-sm font-bold w-8 text-right">{item.count}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="border border-white/10 rounded-xl p-6">
          <h4 className="text-lg font-bold mb-4">Agent Status Distribution</h4>
          <div className="flex items-center justify-center h-48">
            <div className="relative w-40 h-40">
              {/* Pie Chart */}
              <svg viewBox="0 0 100 100" className="w-full h-full">
                {[
                  { status: 'active', count: agents.filter(a => a.status === 'active').length, color: '#FBDF24' },
                  { status: 'waiting', count: agents.filter(a => a.status === 'waiting').length, color: '#F59E0B' },
                  { status: 'completed', count: agents.filter(a => a.status === 'completed').length, color: '#10B981' },
                  { status: 'idle', count: agents.filter(a => a.status === 'idle').length, color: '#6B7280' },
                  { status: 'error', count: agents.filter(a => a.status === 'error').length, color: '#EF4444' },
                ]
                .filter(item => item.count > 0)
                .reduce((acc, item, index, array) => {
                  const total = agents.length;
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
                      key: item.status,
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
                  <div className="text-2xl font-bold">{agents.length}</div>
                  <div className="text-xs text-white/60">Total Agents</div>
                </div>
              </div>
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-3 mt-4">
            {[
              { status: 'active', color: 'bg-primary', label: 'Active' },
              { status: 'waiting', color: 'bg-yellow-500', label: 'Waiting' },
              { status: 'completed', color: 'bg-green-500', label: 'Completed' },
              { status: 'idle', color: 'bg-gray-500', label: 'Idle' },
            ].map((item) => (
              <div key={item.status} className="flex items-center gap-2">
                <div className={`w-3 h-3 ${item.color} rounded-full`} />
                <span className="text-sm">{item.label}</span>
                <span className="text-sm font-bold ml-auto">
                  {agents.filter(a => a.status === item.status).length}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}