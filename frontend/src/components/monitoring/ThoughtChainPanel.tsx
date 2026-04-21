import React, { useState, useRef, useEffect } from 'react';
import { 
  Brain, 
  MessageSquare, 
  Code, 
  Copy, 
  Play, 
  Pause, 
  SkipBack, 
  SkipForward, 
  Volume2, 
  VolumeX,
  Maximize2,
  Minimize2,
  Download,
  Filter,
  Search,
  ChevronRight,
  ChevronLeft,
  ChevronUp,
  ChevronDown,
  Zap,
  Cpu,
  MemoryStick,
  Network,
  Terminal,
  FileText,
  GitBranch,
  Ruler,
  Layout,
  CheckCircle2,
  Shield,
  Server,
  Users,
  AlertCircle,
  Clock,
  TrendingUp,
  BarChart3,
  Activity,
  Settings,
  Eye,
  EyeOff
} from 'lucide-react';
import { motion, AnimatePresence } from 'motion/react';
import { cn } from '../../utils';
import { Thought } from '../../types';

interface ThoughtChainPanelProps {
  thoughts: Thought[];
}

export default function ThoughtChainPanel({ thoughts }: ThoughtChainPanelProps) {
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentThoughtIndex, setCurrentThoughtIndex] = useState(0);
  const [playbackSpeed, setPlaybackSpeed] = useState(1);
  const [isMuted, setIsMuted] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [filterType, setFilterType] = useState<'all' | 'internal' | 'formulating' | 'code'>('all');
  const [expandedThoughts, setExpandedThoughts] = useState<Set<string>>(new Set());
  const [selectedAgent, setSelectedAgent] = useState<string>('all');
  const thoughtsContainerRef = useRef<HTMLDivElement>(null);
  const playIntervalRef = useRef<NodeJS.Timeout>();

  // 获取所有唯一的代理
  const agents = Array.from(new Set(thoughts.map(t => t.agent || 'Unknown')));

  // 过滤想法
  const filteredThoughts = thoughts.filter(thought => {
    // 搜索过滤
    if (searchQuery && !thought.content.toLowerCase().includes(searchQuery.toLowerCase())) {
      return false;
    }
    
    // 类型过滤
    if (filterType !== 'all') {
      if (filterType === 'code' && !thought.isCode) return false;
      if (filterType !== 'code' && thought.type !== filterType) return false;
    }
    
    // 代理过滤
    if (selectedAgent !== 'all' && thought.agent !== selectedAgent) {
      return false;
    }
    
    return true;
  });

  const currentThought = filteredThoughts[currentThoughtIndex];

  const getThoughtIcon = (thought: Thought) => {
    if (thought.isCode) return Code;
    switch (thought.type) {
      case 'internal': return Brain;
      case 'formulating': return MessageSquare;
      default: return MessageSquare;
    }
  };

  const getThoughtColor = (thought: Thought) => {
    if (thought.isCode) return 'text-blue-500';
    switch (thought.type) {
      case 'internal': return 'text-primary';
      case 'formulating': return 'text-purple-500';
      default: return 'text-white/60';
    }
  };

  const getThoughtBgColor = (thought: Thought) => {
    if (thought.isCode) return 'bg-blue-500/10 border-blue-500/20';
    switch (thought.type) {
      case 'internal': return 'bg-primary/10 border-primary/20';
      case 'formulating': return 'bg-purple-500/10 border-purple-500/20';
      default: return 'bg-white/5 border-white/10';
    }
  };

  const getThoughtTypeLabel = (thought: Thought) => {
    if (thought.isCode) return 'Code Generation';
    switch (thought.type) {
      case 'internal': return 'Internal Reasoning';
      case 'formulating': return 'Solution Formulation';
      default: return 'Thought';
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
      playIntervalRef.current = setInterval(() => {
        setCurrentThoughtIndex(prev => {
          if (prev >= filteredThoughts.length - 1) {
            if (playIntervalRef.current) {
              clearInterval(playIntervalRef.current);
            }
            setIsPlaying(false);
            return 0;
          }
          return prev + 1;
        });
      }, 3000 / playbackSpeed);
    }
  };

  const handleSkipBack = () => {
    setCurrentThoughtIndex(prev => Math.max(0, prev - 1));
  };

  const handleSkipForward = () => {
    setCurrentThoughtIndex(prev => Math.min(filteredThoughts.length - 1, prev + 1));
  };

  const handleCopyCode = (code: string) => {
    navigator.clipboard.writeText(code);
    // 可以添加复制成功的提示
  };

  const handleDownloadThoughts = () => {
    const data = {
      thoughts: filteredThoughts,
      timestamp: new Date().toISOString(),
      filter: filterType,
      agent: selectedAgent
    };
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `thoughts-${new Date().toISOString().split('T')[0]}.json`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const toggleExpandThought = (thoughtId: string) => {
    const newExpanded = new Set(expandedThoughts);
    if (newExpanded.has(thoughtId)) {
      newExpanded.delete(thoughtId);
    } else {
      newExpanded.add(thoughtId);
    }
    setExpandedThoughts(newExpanded);
  };

  useEffect(() => {
    return () => {
      if (playIntervalRef.current) {
        clearInterval(playIntervalRef.current);
      }
    };
  }, []);

  useEffect(() => {
    // 当过滤条件变化时重置播放
    if (playIntervalRef.current) {
      clearInterval(playIntervalRef.current);
    }
    setIsPlaying(false);
    setCurrentThoughtIndex(0);
  }, [searchQuery, filterType, selectedAgent]);

  return (
    <div className={cn(
      "space-y-6",
      isFullscreen && "fixed inset-0 z-50 bg-[#0A0A0A] p-6 overflow-auto"
    )}>
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-3">
            <Brain size={24} />
            Thought Chain Panel
          </h2>
          <p className="text-sm text-white/60 mt-1">
            Real-time streaming of AI agent reasoning and decision-making processes
          </p>
        </div>
        
        <div className="flex items-center gap-3">
          <div className="text-sm text-white/40">
            {filteredThoughts.length} thoughts • {agents.length} agents
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
              onClick={handleSkipBack}
              disabled={currentThoughtIndex === 0}
              className={cn(
                "p-2 rounded-lg transition-colors",
                currentThoughtIndex === 0
                  ? "bg-white/5 text-white/20"
                  : "bg-white/5 text-white/60 hover:bg-white/10"
              )}
            >
              <SkipBack size={20} />
            </button>
            
            <div className="text-sm font-medium min-w-[80px] text-center">
              {currentThoughtIndex + 1} / {filteredThoughts.length}
            </div>
            
            <button
              onClick={handleSkipForward}
              disabled={currentThoughtIndex === filteredThoughts.length - 1}
              className={cn(
                "p-2 rounded-lg transition-colors",
                currentThoughtIndex === filteredThoughts.length - 1
                  ? "bg-white/5 text-white/20"
                  : "bg-white/5 text-white/60 hover:bg-white/10"
              )}
            >
              <SkipForward size={20} />
            </button>
            
            <div className="h-6 w-px bg-white/10" />
            
            <button
              onClick={() => setIsMuted(!isMuted)}
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
            >
              {isMuted ? <VolumeX size={20} /> : <Volume2 size={20} />}
            </button>
            
            <select
              value={playbackSpeed}
              onChange={(e) => setPlaybackSpeed(parseFloat(e.target.value))}
              className="bg-transparent text-sm border-none outline-none"
            >
              <option value="0.5">0.5x</option>
              <option value="1">1x</option>
              <option value="1.5">1.5x</option>
              <option value="2">2x</option>
              <option value="3">3x</option>
            </select>
          </div>
          
          <div className="flex items-center gap-3">
            <div className="relative">
              <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-white/40" />
              <input
                type="text"
                placeholder="Search thoughts..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 pr-4 py-2 bg-white/5 border border-white/10 rounded-lg text-sm focus:outline-none focus:border-primary transition-colors"
              />
            </div>
            
            <select
              value={filterType}
              onChange={(e) => setFilterType(e.target.value as any)}
              className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
            >
              <option value="all">All Types</option>
              <option value="internal">Internal Reasoning</option>
              <option value="formulating">Solution Formulation</option>
              <option value="code">Code Generation</option>
            </select>
            
            <select
              value={selectedAgent}
              onChange={(e) => setSelectedAgent(e.target.value)}
              className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
            >
              <option value="all">All Agents</option>
              {agents.map(agent => (
                <option key={agent} value={agent}>{agent}</option>
              ))}
            </select>
            
            <button
              onClick={handleDownloadThoughts}
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
            >
              <Download size={20} />
            </button>
          </div>
        </div>
      </div>

      {/* Current Thought Highlight */}
      {currentThought && (
        <div className={cn(
          "border rounded-xl p-6 transition-all",
          getThoughtBgColor(currentThought),
          "ring-2 ring-primary/30"
        )}>
          <div className="flex justify-between items-start mb-4">
            <div className="flex items-center gap-3">
              <div className={cn(
                "p-2 rounded-lg",
                getThoughtBgColor(currentThought).replace('border-', 'bg-').replace('/20', '/20')
              )}>
                {React.createElement(getThoughtIcon(currentThought), {
                  size: 24,
                  className: getThoughtColor(currentThought)
                })}
              </div>
              <div>
                <div className="text-lg font-bold">{getThoughtTypeLabel(currentThought)}</div>
                <div className="text-sm text-white/60">
                  {currentThought.agent || 'Unknown Agent'} • {currentThought.timestamp || 'Just now'}
                </div>
              </div>
            </div>
            
            <div className="flex items-center gap-2">
              <button
                onClick={() => toggleExpandThought(currentThought.id || `thought-${currentThoughtIndex}`)}
                className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
              >
                {expandedThoughts.has(currentThought.id || `thought-${currentThoughtIndex}`) ? 
                  <ChevronUp size={16} /> : <ChevronDown size={16} />}
              </button>
              
              {currentThought.isCode && (
                <button
                  onClick={() => handleCopyCode(currentThought.code || '')}
                  className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
                >
                  <Copy size={16} />
                </button>
              )}
            </div>
          </div>
          
          <AnimatePresence>
            {expandedThoughts.has(currentThought.id || `thought-${currentThoughtIndex}`) && (
              <motion.div
                initial={{ opacity: 0, height: 0 }}
                animate={{ opacity: 1, height: 'auto' }}
                exit={{ opacity: 0, height: 0 }}
                className="overflow-hidden"
              >
                {currentThought.isCode ? (
                  <div className="mt-4">
                    <div className="flex items-center justify-between mb-2">
                      <div className="text-sm font-bold">{currentThought.filename}</div>
                      <div className="text-xs text-white/40">
                        {currentThought.code?.split('\n').length || 0} lines
                      </div>
                    </div>
                    <pre className="bg-black/40 p-4 rounded-lg overflow-x-auto text-sm">
                      <code className="text-white/80">{currentThought.code}</code>
                    </pre>
                  </div>
                ) : (
                  <p className="text-white/80 leading-relaxed mt-2">
                    {currentThought.content}
                  </p>
                )}
              </motion.div>
            )}
          </AnimatePresence>
          
          {!expandedThoughts.has(currentThought.id || `thought-${currentThoughtIndex}`) && (
            <p className="text-white/80 leading-relaxed line-clamp-3">
              {currentThought.isCode ? (
                <div className="font-mono text-sm opacity-60">
                  {currentThought.code?.split('\n').slice(0, 3).join('\n')}
                  {currentThought.code && currentThought.code.split('\n').length > 3 && '...'}
                </div>
              ) : (
                currentThought.content
              )}
            </p>
          )}
          
          {currentThought.metadata && (
            <div className="mt-4 pt-4 border-t border-white/10">
              <div className="text-xs text-white/40 mb-2">Metadata</div>
              <div className="flex flex-wrap gap-2">
                {Object.entries(currentThought.metadata).map(([key, value]) => (
                  <div key={key} className="px-2 py-1 bg-white/5 rounded text-xs">
                    <span className="text-white/60">{key}:</span>{' '}
                    <span className="text-white/80">{String(value)}</span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* All Thoughts List */}
      <div 
        ref={thoughtsContainerRef}
        className="border border-white/10 rounded-xl p-6 bg-white/[0.02] max-h-[600px] overflow-y-auto"
      >
        <div className="space-y-4">
          <AnimatePresence>
            {filteredThoughts.map((thought, index) => {
              const thoughtId = thought.id || `thought-${index}`;
              const isCurrent = index === currentThoughtIndex;
              const Icon = getThoughtIcon(thought);
              
              return (
                <motion.div
                  key={thoughtId}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                  className={cn(
                    "border rounded-lg p-4 transition-all cursor-pointer hover:scale-[1.01]",
                    getThoughtBgColor(thought),
                    isCurrent && "ring-2 ring-primary"
                  )}
                  onClick={() => setCurrentThoughtIndex(index)}
                >
                  <div className="flex justify-between items-start">
                    <div className="flex items-start gap-3">
                      <div className={cn(
                        "p-2 rounded-lg mt-1",
                        getThoughtBgColor(thought).replace('border-', 'bg-').replace('/20', '/20')
                      )}>
                        <Icon size={16} className={getThoughtColor(thought)} />
                      </div>
                      
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-sm font-bold">{getThoughtTypeLabel(thought)}</span>
                          <span className="text-xs px-2 py-0.5 bg-white/10 rounded">
                            {thought.agent || 'Unknown'}
                          </span>
                          <span className="text-xs text-white/40">
                            {thought.timestamp || `${index * 2}m ago`}
                          </span>
                        </div>
                        
                        {thought.isCode ? (
                          <div className="font-mono text-sm opacity-80">
                            {thought.filename}
                          </div>
                        ) : (
                          <p className="text-white/80 line-clamp-2">
                            {thought.content}
                          </p>
                        )}
                      </div>
                    </div>
                    
                    <div className="flex items-center gap-2">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          toggleExpandThought(thoughtId);
                        }}
                        className="p-1 rounded hover:bg-white/10 transition-colors"
                      >
                        {expandedThoughts.has(thoughtId) ? 
                          <ChevronUp size={14} /> : <ChevronDown size={14} />}
                      </button>
                      
                      {thought.isCode && (
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleCopyCode(thought.code || '');
                          }}
                          className="p-1 rounded hover:bg-white/10 transition-colors"
                        >
                          <Copy size={14} />
                        </button>
                      )}
                    </div>
                  </div>
                  
                  <AnimatePresence>
                    {expandedThoughts.has(thoughtId) && (
                      <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="overflow-hidden mt-3 pt-3 border-t border-white/10"
                      >
                        {thought.isCode ? (
                          <div>
                            <div className="flex items-center justify-between mb-2">
                              <div className="text-sm font-bold">{thought.filename}</div>
                              <div className="text-xs text-white/40">
                                {thought.code?.split('\n').length || 0} lines
                              </div>
                            </div>
                            <pre className="bg-black/40 p-3 rounded overflow-x-auto text-xs">
                              <code className="text-white/80">{thought.code}</code>
                            </pre>
                          </div>
                        ) : (
                          <p className="text-white/80 text-sm leading-relaxed">
                            {thought.content}
                          </p>
                        )}
                        
                        {thought.metadata && (
                          <div className="mt-3 pt-3 border-t border-white/10">
                            <div className="text-xs text-white/40 mb-2">Metadata</div>
                            <div className="flex flex-wrap gap-2">
                              {Object.entries(thought.metadata).map(([key, value]) => (
                                <div key={key} className="px-2 py-1 bg-white/5 rounded text-xs">
                                  <span className="text-white/60">{key}:</span>{' '}
                                  <span className="text-white/80">{String(value)}</span>
                                </div>
                              ))}
                            </div>
                          </div>
                        )}
                      </motion.div>
                    )}
                  </AnimatePresence>
                </motion.div>
              );
            })}
          </AnimatePresence>
          
          {filteredThoughts.length === 0 && (
            <div className="text-center py-12">
              <Brain size={48} className="mx-auto text-white/20 mb-4" />
              <div className="text-lg font-bold text-white/40">No thoughts found</div>
              <p className="text-white/60 mt-2">
                {searchQuery 
                  ? `No thoughts match "${searchQuery}"`
                  : "Try changing your filters or wait for agents to generate thoughts"}
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Statistics */}
      <div className="grid grid-cols-4 gap-4">
        <div className="bg-white/5 border border-white/10 rounded-lg p-4">
          <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Total Thoughts</div>
          <div className="text-2xl font-bold">{thoughts.length}</div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-lg p-4">
          <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Internal Reasoning</div>
          <div className="text-2xl font-bold text-primary">
            {thoughts.filter(t => t.type === 'internal').length}
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-lg p-4">
          <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Code Generation</div>
          <div className="text-2xl font-bold text-blue-500">
            {thoughts.filter(t => t.isCode).length}
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-lg p-4">
          <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Avg Length</div>
          <div className="text-2xl font-bold">
            {Math.round(thoughts.reduce((sum, t) => sum + (t.content?.length || 0), 0) / Math.max(1, thoughts.length))} chars
          </div>
        </div>
      </div>
    </div>
  );
}