import React, { useState, useRef, useEffect } from 'react';
import { 
  Clock, 
  Play, 
  Pause, 
  SkipBack, 
  SkipForward, 
  ZoomIn, 
  ZoomOut, 
  Maximize2, 
  Minimize2,
  Filter,
  Download,
  Calendar,
  ChevronRight,
  ChevronLeft,
  GitBranch,
  Ruler,
  Terminal,
  Layout,
  CheckCircle2,
  Brain,
  AlertCircle,
  Server,
  Users,
  Shield
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../../utils';
import { TimelineEvent } from '../../types';

interface ExecutionTimelineProps {
  events: TimelineEvent[];
}

export default function ExecutionTimeline({ events }: ExecutionTimelineProps) {
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [zoomLevel, setZoomLevel] = useState(1);
  const [selectedEvent, setSelectedEvent] = useState<string | null>(null);
  const [timeRange, setTimeRange] = useState<'1h' | '6h' | '12h' | '24h'>('6h');
  const [viewMode, setViewMode] = useState<'compact' | 'detailed' | 'gantt'>('gantt');
  const timelineRef = useRef<HTMLDivElement>(null);
  const playIntervalRef = useRef<NodeJS.Timeout>();

  const getAgentIcon = (agent: string) => {
    switch (agent.toLowerCase()) {
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

  const getEventColor = (event: TimelineEvent) => {
    if (event.isComplete) return 'bg-green-500';
    if (event.status === 'error') return 'bg-red-500';
    if (event.status === 'warning') return 'bg-yellow-500';
    if (event.status === 'running') return 'bg-primary animate-pulse';
    return 'bg-blue-500';
  };

  const getEventBorderColor = (event: TimelineEvent) => {
    if (event.isComplete) return 'border-green-500/30';
    if (event.status === 'error') return 'border-red-500/30';
    if (event.status === 'warning') return 'border-yellow-500/30';
    if (event.status === 'running') return 'border-primary/30';
    return 'border-blue-500/30';
  };

  const getEventTextColor = (event: TimelineEvent) => {
    if (event.isComplete) return 'text-green-500';
    if (event.status === 'error') return 'text-red-500';
    if (event.status === 'warning') return 'text-yellow-500';
    if (event.status === 'running') return 'text-primary';
    return 'text-blue-500';
  };

  const formatTime = (time: string) => {
    // 简单的时间格式化
    return time.replace('T-', '').replace('m', ' min');
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
        setCurrentTime(prev => {
          const maxTime = events.reduce((max, event) => {
            const timeNum = parseInt(event.time.replace('T-', '').replace('m', ''));
            return Math.max(max, timeNum);
          }, 0);
          return prev >= maxTime ? 0 : prev + 1;
        });
      }, 1000);
    }
  };

  const handleSkipBack = () => {
    setCurrentTime(prev => Math.max(0, prev - 5));
  };

  const handleSkipForward = () => {
    const maxTime = events.reduce((max, event) => {
      const timeNum = parseInt(event.time.replace('T-', '').replace('m', ''));
      return Math.max(max, timeNum);
    }, 0);
    setCurrentTime(prev => Math.min(maxTime, prev + 5));
  };

  const handleZoomIn = () => {
    setZoomLevel(prev => Math.min(3, prev + 0.25));
  };

  const handleZoomOut = () => {
    setZoomLevel(prev => Math.max(0.5, prev - 0.25));
  };

  const handleDownload = () => {
    // 导出时间线数据
    const data = {
      events,
      timestamp: new Date().toISOString(),
      timeRange,
      viewMode
    };
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `timeline-${new Date().toISOString().split('T')[0]}.json`;
    a.click();
    URL.revokeObjectURL(url);
  };

  useEffect(() => {
    return () => {
      if (playIntervalRef.current) {
        clearInterval(playIntervalRef.current);
      }
    };
  }, []);

  // 按代理分组事件
  const eventsByAgent = events.reduce((acc, event) => {
    if (!acc[event.agent]) {
      acc[event.agent] = [];
    }
    acc[event.agent].push(event);
    return acc;
  }, {} as Record<string, TimelineEvent[]>);

  const agents = Object.keys(eventsByAgent);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-3">
            <Clock size={24} />
            Execution Timeline
          </h2>
          <p className="text-sm text-white/60 mt-1">
            Gantt chart visualization of agent execution over time
          </p>
        </div>
        
        <div className="flex items-center gap-3">
          <div className="text-sm text-white/40">
            {events.length} events • {agents.length} agents
          </div>
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
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
            >
              <SkipBack size={20} />
            </button>
            
            <button
              onClick={handleSkipForward}
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
            >
              <SkipForward size={20} />
            </button>
            
            <div className="h-6 w-px bg-white/10" />
            
            <button
              onClick={handleZoomOut}
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
            >
              <ZoomOut size={20} />
            </button>
            
            <div className="text-sm font-medium min-w-[60px] text-center">
              {zoomLevel.toFixed(2)}x
            </div>
            
            <button
              onClick={handleZoomIn}
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
            >
              <ZoomIn size={20} />
            </button>
            
            <div className="h-6 w-px bg-white/10" />
            
            <button
              onClick={handleDownload}
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
            >
              <Download size={20} />
            </button>
          </div>
          
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2">
              <Filter size={16} className="text-white/40" />
              <select 
                value={timeRange}
                onChange={(e) => setTimeRange(e.target.value as any)}
                className="bg-transparent text-sm border-none outline-none"
              >
                <option value="1h">Last 1 hour</option>
                <option value="6h">Last 6 hours</option>
                <option value="12h">Last 12 hours</option>
                <option value="24h">Last 24 hours</option>
              </select>
            </div>
            
            <div className="flex items-center gap-2">
              <Calendar size={16} className="text-white/40" />
              <select 
                value={viewMode}
                onChange={(e) => setViewMode(e.target.value as any)}
                className="bg-transparent text-sm border-none outline-none"
              >
                <option value="compact">Compact</option>
                <option value="detailed">Detailed</option>
                <option value="gantt">Gantt</option>
              </select>
            </div>
          </div>
        </div>
        
        {/* Timeline Slider */}
        <div className="mt-4">
          <div className="flex justify-between text-xs text-white/40 mb-2">
            <span>Start</span>
            <span>Current: {currentTime}m</span>
            <span>End</span>
          </div>
          <input
            type="range"
            min="0"
            max={events.reduce((max, event) => {
              const timeNum = parseInt(event.time.replace('T-', '').replace('m', ''));
              return Math.max(max, timeNum);
            }, 0)}
            value={currentTime}
            onChange={(e) => setCurrentTime(parseInt(e.target.value))}
            className="w-full h-2 bg-white/10 rounded-lg appearance-none [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:h-4 [&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-primary"
          />
        </div>
      </div>

      {/* Timeline Visualization */}
      <div 
        ref={timelineRef}
        className="border border-white/10 rounded-xl p-6 bg-white/[0.02] overflow-x-auto"
      >
        <div className="min-w-[800px]">
          {/* Time Scale */}
          <div className="flex mb-6">
            <div className="w-40" /> {/* Agent names column */}
            <div className="flex-1 relative h-8">
              {Array.from({ length: 13 }).map((_, i) => {
                const time = i * 5; // 5-minute intervals
                return (
                  <div
                    key={i}
                    className="absolute top-0 text-xs text-white/40"
                    style={{ left: `${(time / 60) * 100}%` }}
                  >
                    <div className="h-4 border-l border-white/20" />
                    <div className="mt-1">{time}m</div>
                  </div>
                );
              })}
              
              {/* Current time indicator */}
              <div 
                className="absolute top-0 w-px h-full bg-primary z-10"
                style={{ left: `${(currentTime / 60) * 100}%` }}
              >
                <div className="absolute -top-2 -left-2 w-4 h-4 bg-primary rounded-full" />
              </div>
            </div>
          </div>

          {/* Agent Rows */}
          <div className="space-y-8">
            {agents.map((agent) => {
              const Icon = getAgentIcon(agent);
              const agentEvents = eventsByAgent[agent];
              
              return (
                <div key={agent} className="flex items-center group">
                  <div className="w-40 flex items-center gap-3">
                    <div className="p-2 bg-white/5 rounded-lg">
                      <Icon size={16} className="text-white/60" />
                    </div>
                    <div>
                      <div className="font-bold">{agent}</div>
                      <div className="text-xs text-white/40">
                        {agentEvents.length} events
                      </div>
                    </div>
                  </div>
                  
                  <div className="flex-1 relative h-12">
                    {/* Background grid */}
                    <div className="absolute inset-0 flex">
                      {Array.from({ length: 13 }).map((_, i) => (
                        <div
                          key={i}
                          className="flex-1 border-l border-white/5 last:border-r"
                        />
                      ))}
                    </div>
                    
                    {/* Events */}
                    {agentEvents.map((event) => {
                      const startPercent = (parseInt(event.time.replace('T-', '').replace('m', '')) / 60) * 100;
                      const widthPercent = (parseInt(event.width) / 100) * 100;
                      
                      return (
                        <motion.div
                          key={event.id}
                          initial={{ opacity: 0, scale: 0.9 }}
                          animate={{ opacity: 1, scale: 1 }}
                          className={cn(
                            "absolute top-1/2 -translate-y-1/2 h-8 rounded-lg border px-3 flex items-center cursor-pointer transition-all hover:scale-105 hover:z-10",
                            getEventColor(event),
                            getEventBorderColor(event),
                            selectedEvent === event.id && "ring-2 ring-white"
                          )}
                          style={{
                            left: `${startPercent}%`,
                            width: `${widthPercent}%`,
                          }}
                          onClick={() => setSelectedEvent(event.id === selectedEvent ? null : event.id)}
                          title={`${event.action} - ${formatTime(event.time)}`}
                        >
                          <div className="flex items-center gap-2 w-full">
                            <div className={cn(
                              "w-2 h-2 rounded-full",
                              getEventTextColor(event).replace('text-', 'bg-')
                            )} />
                            <div className="text-xs font-bold truncate">
                              {event.action.replace(/_/g, ' ')}
                            </div>
                            <div className="text-xs opacity-60 ml-auto">
                              {formatTime(event.time)}
                            </div>
                          </div>
                        </motion.div>
                      );
                    })}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>

      {/* Event Details Panel */}
      {selectedEvent && (
        <div className="border border-white/10 rounded-xl p-6 bg-white/[0.02]">
          <div className="flex justify-between items-start mb-4">
            <h3 className="text-lg font-bold">Event Details</h3>
            <button
              onClick={() => setSelectedEvent(null)}
              className="p-1 hover:bg-white/10 rounded transition-colors"
            >
              <ChevronRight size={20} className="rotate-90" />
            </button>
          </div>
          
          {(() => {
            const event = events.find(e => e.id === selectedEvent);
            if (!event) return null;
            
            const Icon = getAgentIcon(event.agent);
            
            return (
              <div className="space-y-4">
                <div className="flex items-center gap-4">
                  <div className={cn(
                    "p-3 rounded-lg",
                    getEventColor(event).replace('bg-', 'bg-').replace('text-', 'bg-') + '/20'
                  )}>
                    <Icon size={24} className={getEventTextColor(event)} />
                  </div>
                  
                  <div>
                    <div className="text-2xl font-bold capitalize">
                      {event.action.replace(/_/g, ' ')}
                    </div>
                    <div className="text-sm text-white/60">
                      {event.agent} • {formatTime(event.time)}
                    </div>
                  </div>
                  
                  <div className={cn(
                    "ml-auto px-4 py-2 rounded-lg text-sm font-bold",
                    getEventColor(event),
                    getEventBorderColor(event)
                  )}>
                    {event.isComplete ? 'COMPLETED' : 
                     event.status === 'running' ? 'RUNNING' :
                     event.status === 'error' ? 'ERROR' : 'PENDING'}
                  </div>
                </div>
                
                <div className="grid grid-cols-3 gap-4">
                  <div className="bg-white/5 border border-white/10 rounded-lg p-4">
                    <div className="text-xs text-white/40 mb-1">Duration</div>
                    <div className="text-lg font-bold">{event.width}</div>
                  </div>
                  
                  <div className="bg-white/5 border border-white/10 rounded-lg p-4">
                    <div className="text-xs text-white/40 mb-1">Start Time</div>
                    <div className="text-lg font-bold">{formatTime(event.time)}</div>
                  </div>
                  
                  <div className="bg-white/5 border border-white/10 rounded-lg p-4">
                    <div className="text-xs text-white/40 mb-1">Status</div>
                    <div className={cn(
                      "text-lg font-bold",
                      getEventTextColor(event)
                    )}>
                      {event.isComplete ? 'Complete' : 
                       event.status === 'running' ? 'Running' :
                       event.status === 'error' ? 'Error' : 'Pending'}
                    </div>
                  </div>
                </div>
                
                <div className="bg-white/5 border border-white/10 rounded-lg p-4">
                  <div className="text-sm font-bold mb-2">Description</div>
                  <p className="text-white/80">
                    {event.description || `The ${event.agent} is performing ${event.action.replace(/_/g, ' ')}. This is part of the pipeline execution.`}
                  </p>
                </div>
                
                {event.metadata && (
                  <div className="bg-white/5 border border-white/10 rounded-lg p-4">
                    <div className="text-sm font-bold mb-2">Metadata</div>
                    <pre className="text-xs text-white/60 overflow-x-auto">
                      {JSON.stringify(event.metadata, null, 2)}
                    </pre>
                  </div>
                )}
              </div>
            );
          })()}
        </div>
      )}

      {/* Statistics */}
      <div className="grid grid-cols-4 gap-4">
        <div className="bg-white/5 border border-white/10 rounded-lg p-4">
          <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Total Events</div>
          <div className="text-2xl font-bold">{events.length}</div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-lg p-4">
          <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Completed</div>
          <div className="text-2xl font-bold text-green-500">
            {events.filter(e => e.isComplete).length}
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-lg p-4">
          <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Running</div>
          <div className="text-2xl font-bold text-primary">
            {events.filter(e => e.status === 'running').length}
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-lg p-4">
          <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Errors</div>
          <div className="text-2xl font-bold text-red-500">
            {events.filter(e => e.status === 'error').length}
          </div>
        </div>
      </div>
    </div>
  );
}