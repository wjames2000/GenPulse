import React, { useState, useRef, useEffect } from 'react';
import { 
  Terminal, 
  Play, 
  Pause, 
  SkipBack, 
  SkipForward, 
  RefreshCw,
  Copy,
  Download,
  Filter,
  Search,
  ChevronUp,
  ChevronDown,
  Maximize2,
  Minimize2,
  Trash2,
  Settings,
  Eye,
  EyeOff,
  Volume2,
  VolumeX,
  Zap,
  Cpu,
  MemoryStick,
  Network,
  HardDrive,
  Server,
  Database,
  GitBranch,
  Layout,
  CheckCircle2,
  AlertCircle,
  Brain,
  MessageSquare,
  Code,
  FileText,
  Users,
  Shield,
  Clock
} from 'lucide-react';
import { cn } from '../../utils';

interface TerminalOutputPanelProps {
  output: string[];
}

export default function TerminalOutputPanel({ output }: TerminalOutputPanelProps) {
  const [isPlaying, setIsPlaying] = useState(false);
  const [playbackSpeed, setPlaybackSpeed] = useState(1);
  const [isFollowing, setIsFollowing] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [filterLevel, setFilterLevel] = useState<'all' | 'info' | 'success' | 'warn' | 'error'>('all');
  const [selectedLine, setSelectedLine] = useState<number | null>(null);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [showTimestamps, setShowTimestamps] = useState(true);
  const [showLineNumbers, setShowLineNumbers] = useState(true);
  const [wrapLines, setWrapLines] = useState(false);
  const terminalRef = useRef<HTMLDivElement>(null);
  const playIntervalRef = useRef<NodeJS.Timeout>();

  // 解析终端输出行
  const parsedOutput = output.map((line, index) => {
    // 尝试解析常见的终端输出格式
    const timestampMatch = line.match(/\[(\d{2}:\d{2}:\d{2})\]/);
    const levelMatch = line.match(/(INFO|SUCCESS|WARN|ERROR|DEBUG)/);
    const agentMatch = line.match(/(Orchestrator|Architect|Backend|Frontend|QA|DevOps|Reviewer):/);
    
    return {
      id: index,
      raw: line,
      timestamp: timestampMatch ? timestampMatch[1] : null,
      level: levelMatch ? levelMatch[1].toLowerCase() as 'info' | 'success' | 'warn' | 'error' | 'debug' : 'info',
      agent: agentMatch ? agentMatch[1] : null,
      message: line.replace(/\[.*?\]/g, '').replace(/(INFO|SUCCESS|WARN|ERROR|DEBUG):?/, '').trim()
    };
  });

  // 过滤输出
  const filteredOutput = parsedOutput.filter(line => {
    // 搜索过滤
    if (searchQuery && !line.raw.toLowerCase().includes(searchQuery.toLowerCase())) {
      return false;
    }
    
    // 级别过滤
    if (filterLevel !== 'all' && line.level !== filterLevel) {
      return false;
    }
    
    return true;
  });

  const getLevelColor = (level: string) => {
    switch (level) {
      case 'success': return 'text-green-500';
      case 'error': return 'text-red-500';
      case 'warn': return 'text-yellow-500';
      case 'debug': return 'text-blue-500';
      default: return 'text-white/80';
    }
  };

  const getLevelBgColor = (level: string) => {
    switch (level) {
      case 'success': return 'bg-green-500/10';
      case 'error': return 'bg-red-500/10';
      case 'warn': return 'bg-yellow-500/10';
      case 'debug': return 'bg-blue-500/10';
      default: return 'bg-white/5';
    }
  };

  const getAgentIcon = (agent: string | null) => {
    if (!agent) return Terminal;
    switch (agent.toLowerCase()) {
      case 'orchestrator': return GitBranch;
      case 'architect': return Brain;
      case 'backend': return Server;
      case 'frontend': return Layout;
      case 'qa': return CheckCircle2;
      case 'devops': return Database;
      case 'reviewer': return Shield;
      default: return Users;
    }
  };

  const handlePlayPause = () => {
    if (isPlaying) {
      if (playIntervalRef.current) {
        clearInterval(playIntervalRef.current);
      }
      setIsPlaying(false);
    } else {
      setIsPlaying(true);
      // 模拟播放终端输出
      playIntervalRef.current = setInterval(() => {
        // 这里可以添加播放逻辑
      }, 1000 / playbackSpeed);
    }
  };

  const handleCopyAll = () => {
    const text = filteredOutput.map(line => line.raw).join('\n');
    navigator.clipboard.writeText(text);
  };

  const handleDownload = () => {
    const text = filteredOutput.map(line => line.raw).join('\n');
    const blob = new Blob([text], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `terminal-output-${new Date().toISOString().split('T')[0]}.log`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const handleClear = () => {
    // 这里可以添加清除逻辑
    console.log('Clear terminal output');
  };

  // 自动滚动到底部
  useEffect(() => {
    if (isFollowing && terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  }, [filteredOutput, isFollowing]);

  // 清理定时器
  useEffect(() => {
    return () => {
      if (playIntervalRef.current) {
        clearInterval(playIntervalRef.current);
      }
    };
  }, []);

  return (
    <div className={cn(
      "space-y-6",
      isFullscreen && "fixed inset-0 z-50 bg-[#0A0A0A] p-6 overflow-auto"
    )}>
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-3">
            <Terminal size={24} />
            Terminal Output
          </h2>
          <p className="text-sm text-white/60 mt-1">
            Real-time command execution output and system logs
          </p>
        </div>
        
        <div className="flex items-center gap-3">
          <div className="text-sm text-white/40">
            {filteredOutput.length} lines • {new Set(parsedOutput.map(l => l.agent)).size} agents
          </div>
          <button
            onClick={() => setIsFullscreen(!isFullscreen)}
            className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
          >
            {isFullscreen ? <Minimize2 size={20} /> : <Maximize2 size={20} />}
          </button>
        </div>
      </div>

      {/* Controls */}
      <div className="bg-white/5 border border-white/10 rounded-xl p-4">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <button
              onClick={handlePlayPause}
              className={cn(
                "p-2 rounded-lg transition-colors",
                isPlaying 
                  ? "bg-red-500/20 text-red-500 hover:bg-red-500/30" 
                  : "bg-primary/20 text-primary hover:bg-primary/30"
              )}
            >
              {isPlaying ? <Pause size={20} /> : <Play size={20} />}
            </button>
            
            <button
              onClick={() => setIsFollowing(!isFollowing)}
              className={cn(
                "p-2 rounded-lg transition-colors",
                isFollowing
                  ? "bg-primary/20 text-primary" 
                  : "bg-white/5 text-white/60 hover:bg-white/10"
              )}
              title={isFollowing ? "Auto-scroll enabled" : "Auto-scroll disabled"}
            >
              {isFollowing ? <Eye size={20} /> : <EyeOff size={20} />}
            </button>
            
            <div className="h-6 w-px bg-white/10" />
            
            <select
              value={playbackSpeed}
              onChange={(e) => setPlaybackSpeed(parseFloat(e.target.value))}
              className="bg-transparent text-sm border-none outline-none"
            >
              <option value="0.5">0.5x</option>
              <option value="1">1x</option>
              <option value="2">2x</option>
              <option value="5">5x</option>
            </select>
          </div>
          
          <div className="flex items-center gap-3">
            <div className="relative">
              <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-white/40" />
              <input
                type="text"
                placeholder="Search terminal output..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 pr-4 py-2 bg-white/5 border border-white/10 rounded-lg text-sm focus:outline-none focus:border-primary transition-colors"
              />
            </div>
            
            <select
              value={filterLevel}
              onChange={(e) => setFilterLevel(e.target.value as any)}
              className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
            >
              <option value="all">All Levels</option>
              <option value="info">Info</option>
              <option value="success">Success</option>
              <option value="warn">Warning</option>
              <option value="error">Error</option>
              <option value="debug">Debug</option>
            </select>
          </div>
          
          <div className="flex items-center gap-3">
            <button
              onClick={handleCopyAll}
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
              title="Copy all output"
            >
              <Copy size={20} />
            </button>
            
            <button
              onClick={handleDownload}
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
              title="Download output"
            >
              <Download size={20} />
            </button>
            
            <button
              onClick={handleClear}
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
              title="Clear terminal"
            >
              <Trash2 size={20} />
            </button>
            
            <button
              onClick={() => setShowTimestamps(!showTimestamps)}
              className={cn(
                "p-2 rounded-lg transition-colors",
                showTimestamps
                  ? "bg-primary/20 text-primary" 
                  : "bg-white/5 text-white/60 hover:bg-white/10"
              )}
              title="Toggle timestamps"
            >
              <Clock size={20} />
            </button>
            
            <button
              onClick={() => setShowLineNumbers(!showLineNumbers)}
              className={cn(
                "p-2 rounded-lg transition-colors",
                showLineNumbers
                  ? "bg-primary/20 text-primary" 
                  : "bg-white/5 text-white/60 hover:bg-white/10"
              )}
              title="Toggle line numbers"
            >
              <Hash size={20} />
            </button>
          </div>
        </div>
        
        {/* Statistics */}
        <div className="grid grid-cols-5 gap-4 mt-4">
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Total Lines</div>
            <div className="text-lg font-bold">{output.length}</div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Errors</div>
            <div className="text-lg font-bold text-red-500">
              {parsedOutput.filter(l => l.level === 'error').length}
            </div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Warnings</div>
            <div className="text-lg font-bold text-yellow-500">
              {parsedOutput.filter(l => l.level === 'warn').length}
            </div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Success</div>
            <div className="text-lg font-bold text-green-500">
              {parsedOutput.filter(l => l.level === 'success').length}
            </div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Agents</div>
            <div className="text-lg font-bold">
              {new Set(parsedOutput.map(l => l.agent)).size}
            </div>
          </div>
        </div>
      </div>

      {/* Terminal Output */}
      <div 
        ref={terminalRef}
        className={cn(
          "border border-white/10 rounded-xl bg-black font-mono text-sm overflow-y-auto",
          isFullscreen ? "h-[calc(100vh-300px)]" : "h-[500px]"
        )}
      >
        <div className="p-4 space-y-1">
          {filteredOutput.map((line) => {
            const AgentIcon = getAgentIcon(line.agent);
            const isSelected = selectedLine === line.id;
            
            return (
              <div
                key={line.id}
                className={cn(
                  "px-3 py-2 rounded-lg transition-colors cursor-pointer group",
                  getLevelBgColor(line.level),
                  isSelected && "ring-2 ring-primary"
                )}
                onClick={() => setSelectedLine(line.id === selectedLine ? null : line.id)}
              >
                <div className="flex items-start gap-3">
                  {/* Line number */}
                  {showLineNumbers && (
                    <div className="text-white/40 text-xs font-mono w-8 shrink-0">
                      {line.id + 1}
                    </div>
                  )}
                  
                  {/* Timestamp */}
                  {showTimestamps && line.timestamp && (
                    <div className="text-white/40 text-xs font-mono w-20 shrink-0">
                      [{line.timestamp}]
                    </div>
                  )}
                  
                  {/* Level badge */}
                  <div className={cn(
                    "px-2 py-0.5 rounded text-xs font-bold uppercase shrink-0",
                    getLevelColor(line.level),
                    getLevelBgColor(line.level)
                  )}>
                    {line.level}
                  </div>
                  
                  {/* Agent icon */}
                  {line.agent && (
                    <div className="flex items-center gap-1 text-white/60 shrink-0">
                      <AgentIcon size={12} />
                      <span className="text-xs">{line.agent}</span>
                    </div>
                  )}
                  
                  {/* Message */}
                  <div className={cn(
                    "flex-1 font-mono",
                    wrapLines ? "whitespace-pre-wrap" : "whitespace-nowrap overflow-x-auto",
                    getLevelColor(line.level)
                  )}>
                    {line.message}
                  </div>
                  
                  {/* Actions */}
                  <div className="opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        navigator.clipboard.writeText(line.raw);
                      }}
                      className="p-1 rounded hover:bg-white/10 transition-colors"
                      title="Copy line"
                    >
                      <Copy size={12} />
                    </button>
                  </div>
                </div>
                
                {/* Expanded view for selected line */}
                {isSelected && (
                  <div className="mt-3 pt-3 border-t border-white/10">
                    <div className="text-xs text-white/40 mb-2">Raw Output</div>
                    <pre className="text-xs bg-black/40 p-3 rounded overflow-x-auto">
                      {line.raw}
                    </pre>
                  </div>
                )}
              </div>
            );
          })}
          
          {filteredOutput.length === 0 && (
            <div className="text-center py-12">
              <Terminal size={48} className="mx-auto text-white/20 mb-4" />
              <div className="text-lg font-bold text-white/40">No terminal output</div>
              <p className="text-white/60 mt-2">
                {searchQuery 
                  ? `No output matches "${searchQuery}"`
                  : "Wait for commands to be executed or check your filters"}
              </p>
            </div>
          )}
        </div>
        
        {/* Terminal prompt */}
        <div className="sticky bottom-0 border-t border-white/10 bg-black/80 backdrop-blur-sm p-4">
          <div className="flex items-center gap-2">
            <div className="text-green-500 font-bold">$</div>
            <div className="flex-1">
              <input
                type="text"
                placeholder="Type a command to execute..."
                className="w-full bg-transparent border-none outline-none text-white/80 font-mono"
                readOnly
              />
            </div>
            <div className="text-white/40 text-xs">
              Press Enter to execute
            </div>
          </div>
        </div>
      </div>

      {/* Command History */}
      <div className="border border-white/10 rounded-xl p-6">
        <h4 className="text-lg font-bold mb-4">Recent Commands</h4>
        <div className="space-y-3">
          {[
            { command: 'npm run build', status: 'success', time: '2 minutes ago' },
            { command: 'go test ./...', status: 'success', time: '5 minutes ago' },
            { command: 'docker compose up', status: 'running', time: '10 minutes ago' },
            { command: 'git push origin main', status: 'error', time: '15 minutes ago' },
            { command: 'npm install', status: 'success', time: '20 minutes ago' },
          ].map((cmd, i) => (
            <div key={i} className="flex items-center justify-between p-3 bg-white/5 rounded-lg hover:bg-white/10 transition-colors">
              <div className="flex items-center gap-3">
                <div className={cn(
                  "w-2 h-2 rounded-full",
                  cmd.status === 'success' ? 'bg-green-500' :
                  cmd.status === 'error' ? 'bg-red-500' :
                  'bg-primary animate-pulse'
                )} />
                <code className="font-mono text-sm">{cmd.command}</code>
              </div>
              <div className="flex items-center gap-4">
                <span className="text-xs text-white/40">{cmd.time}</span>
                <button className="text-xs px-3 py-1 bg-white/5 rounded hover:bg-white/10 transition-colors">
                  Run Again
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* System Metrics */}
      <div className="grid grid-cols-4 gap-4">
        <div className="bg-white/5 border border-white/10 rounded-lg p-4">
          <div className="flex items-center gap-3 mb-3">
            <Cpu size={20} className="text-primary" />
            <div className="text-xs text-white/40 uppercase tracking-wider">CPU Usage</div>
          </div>
          <div className="text-2xl font-bold">24%</div>
          <div className="h-2 bg-white/10 rounded-full mt-2 overflow-hidden">
            <div className="h-full bg-primary rounded-full" style={{ width: '24%' }} />
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-lg p-4">
          <div className="flex items-center gap-3 mb-3">
            <MemoryStick size={20} className="text-blue-500" />
            <div className="text-xs text-white/40 uppercase tracking-wider">Memory</div>
          </div>
          <div className="text-2xl font-bold">1.2GB</div>
          <div className="h-2 bg-white/10 rounded-full mt-2 overflow-hidden">
            <div className="h-full bg-blue-500 rounded-full" style={{ width: '64%' }} />
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-lg p-4">
          <div className="flex items-center gap-3 mb-3">
            <HardDrive size={20} className="text-green-500" />
            <div className="text-xs text-white/40 uppercase tracking-wider">Disk I/O</div>
          </div>
          <div className="text-2xl font-bold">12MB/s</div>
          <div className="h-2 bg-white/10 rounded-full mt-2 overflow-hidden">
            <div className="h-full bg-green-500 rounded-full" style={{ width: '40%' }} />
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-lg p-4">
          <div className="flex items-center gap-3 mb-3">
            <Network size={20} className="text-purple-500" />
            <div className="text-xs text-white/40 uppercase tracking-wider">Network</div>
          </div>
          <div className="text-2xl font-bold">45KB/s</div>
          <div className="h-2 bg-white/10 rounded-full mt-2 overflow-hidden">
            <div className="h-full bg-purple-500 rounded-full" style={{ width: '18%' }} />
          </div>
        </div>
      </div>
    </div>
  );
}

// 添加缺失的Hash图标组件
function Hash(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      {...props}
    >
      <line x1="4" y1="9" x2="20" y2="9" />
      <line x1="4" y1="15" x2="20" y2="15" />
      <line x1="10" y1="3" x2="8" y2="21" />
      <line x1="16" y1="3" x2="14" y2="21" />
    </svg>
  );
}