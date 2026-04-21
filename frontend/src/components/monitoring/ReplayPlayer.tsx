import React, { useState, useEffect, useRef } from 'react';
import { 
  Play, 
  Pause, 
  StopCircle, 
  SkipBack, 
  SkipForward, 
  FastForward, 
  Rewind,
  Maximize2,
  Minimize2,
  Volume2,
  VolumeX,
  Settings,
  Clock,
  Zap,
  Terminal,
  User,
  FileText,
  AlertCircle,
  ChevronRight,
  ChevronLeft,
  X,
  Expand,
  Minimize,
  RotateCcw,
  RotateCw
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../../utils';
import { api } from '../../services/api';

interface ReplayPlayerProps {
  recordId: string;
  onClose?: () => void;
}

interface ReplayState {
  record_id: string;
  trace_id: string;
  status: 'initializing' | 'playing' | 'paused' | 'completed' | 'error';
  current_time: string;
  start_time: string;
  end_time: string;
  playback_speed: number;
  current_span_index: number;
  total_spans: number;
  progress: number;
  metadata?: Record<string, any>;
}

interface SpanData {
  id: string;
  trace_id: string;
  parent_id?: string;
  name: string;
  kind: string;
  start_time: string;
  end_time: string;
  duration: number;
  attributes: Record<string, any>;
  events?: any[];
  status: string;
  status_message?: string;
  resource?: Record<string, any>;
  agent_name?: string;
  tool_name?: string;
  pipeline_id?: string;
}

export default function ReplayPlayer({ recordId, onClose }: ReplayPlayerProps) {
  const [replayState, setReplayState] = useState<ReplayState | null>(null);
  const [spans, setSpans] = useState<SpanData[]>([]);
  const [currentSpanIndex, setCurrentSpanIndex] = useState(0);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [isMuted, setIsMuted] = useState(false);
  const [showSettings, setShowSettings] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  const playerRef = useRef<HTMLDivElement>(null);
  const timelineRef = useRef<HTMLDivElement>(null);
  const eventSubscriptionRef = useRef<(() => void) | null>(null);

  const loadReplayData = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const [state, initialSpans] = await Promise.all([
        api.getReplayState(recordId),
        api.getReplayData(recordId, 0, 100)
      ]);
      
      setReplayState(state);
      setSpans(initialSpans as SpanData[]);
      
      if (state.status === 'playing' || state.status === 'paused') {
        setCurrentSpanIndex(state.current_span_index || 0);
      }
      
      subscribeToReplayEvents();
    } catch (err: any) {
      setError(err.message || 'Failed to load replay data');
      console.error('Failed to load replay data:', err);
    } finally {
      setLoading(false);
    }
  };

  const subscribeToReplayEvents = () => {
    if (eventSubscriptionRef.current) {
      eventSubscriptionRef.current();
    }
    
    const unsubscribe = api.subscribeToReplayEvents(recordId, (event, data) => {
      if (event === 'replay:progress') {
        setReplayState(data.state);
        setCurrentSpanIndex(data.state.current_span_index || 0);
        
        if (data.state.current_span_index > spans.length - 10) {
          loadMoreSpans();
        }
      } else if (event === 'replay:ended') {
        setReplayState(data.state);
        if (eventSubscriptionRef.current) {
          eventSubscriptionRef.current();
          eventSubscriptionRef.current = null;
        }
      }
    });
    
    eventSubscriptionRef.current = unsubscribe;
  };

  const loadMoreSpans = async () => {
    try {
      const moreSpans = await api.getReplayData(recordId, spans.length, 50);
      setSpans(prev => [...prev, ...(moreSpans as SpanData[])]);
    } catch (err) {
      console.error('Failed to load more spans:', err);
    }
  };

  const handleControl = async (action: string, params?: any) => {
    try {
      const newState = await api.controlReplay(recordId, action, params);
      setReplayState(newState);
      
      if (action === 'stop') {
        if (eventSubscriptionRef.current) {
          eventSubscriptionRef.current();
          eventSubscriptionRef.current = null;
        }
      }
    } catch (err) {
      console.error('Failed to control replay:', err);
    }
  };

  const handleSeek = (progress: number) => {
    handleControl('seek', { progress });
  };

  const handleSpeedChange = (speed: number) => {
    handleControl('speed', { speed });
  };

  const handleFullscreen = () => {
    if (!playerRef.current) return;
    
    if (!isFullscreen) {
      if (playerRef.current.requestFullscreen) {
        playerRef.current.requestFullscreen();
      }
    } else {
      if (document.exitFullscreen) {
        document.exitFullscreen();
      }
    }
    
    setIsFullscreen(!isFullscreen);
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
  };

  const formatDuration = (ms: number) => {
    const seconds = Math.floor(ms / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);

    if (hours > 0) {
      return `${hours}:${minutes % 60}:${seconds % 60}`;
    } else if (minutes > 0) {
      return `${minutes}:${seconds % 60}`;
    } else {
      return `0:${seconds}`;
    }
  };

  const getSpanColor = (span: SpanData) => {
    if (span.status === 'OK') return 'bg-green-500';
    if (span.status === 'ERROR') return 'bg-red-500';
    if (span.agent_name) return 'bg-blue-500';
    if (span.tool_name) return 'bg-purple-500';
    return 'bg-gray-500';
  };

  const getCurrentSpan = () => {
    if (currentSpanIndex >= 0 && currentSpanIndex < spans.length) {
      return spans[currentSpanIndex];
    }
    return null;
  };

  const renderSpanDetails = (span: SpanData) => {
    return (
      <div className="space-y-4">
        <div>
          <h4 className="text-lg font-bold mb-2">{span.name}</h4>
          <div className="flex items-center gap-4 text-sm text-white/60">
            {span.agent_name && (
              <div className="flex items-center gap-1">
                <User size={14} />
                <span>{span.agent_name}</span>
              </div>
            )}
            {span.tool_name && (
              <div className="flex items-center gap-1">
                <Terminal size={14} />
                <span>{span.tool_name}</span>
              </div>
            )}
            <div className="flex items-center gap-1">
              <Clock size={14} />
              <span>{formatDuration(span.duration)}</span>
            </div>
            <div className={cn(
              "px-2 py-1 rounded text-xs",
              span.status === 'OK' ? "bg-green-500/20 text-green-400" :
              span.status === 'ERROR' ? "bg-red-500/20 text-red-400" :
              "bg-gray-500/20 text-gray-400"
            )}>
              {span.status}
            </div>
          </div>
        </div>

        {span.status_message && (
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-sm font-medium mb-1">Status Message</div>
            <div className="text-sm text-white/80">{span.status_message}</div>
          </div>
        )}

        {span.attributes && Object.keys(span.attributes).length > 0 && (
          <div>
            <div className="text-sm font-medium mb-2">Attributes</div>
            <div className="bg-white/5 border border-white/10 rounded-lg overflow-hidden">
              <table className="w-full">
                <tbody>
                  {Object.entries(span.attributes).map(([key, value]) => (
                    <tr key={key} className="border-b border-white/10 last:border-b-0">
                      <td className="p-3 text-sm font-medium text-white/60">{key}</td>
                      <td className="p-3 text-sm text-white/80">
                        {typeof value === 'object' ? JSON.stringify(value, null, 2) : String(value)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {span.events && span.events.length > 0 && (
          <div>
            <div className="text-sm font-medium mb-2">Events</div>
            <div className="space-y-2">
              {span.events.map((event: any, index: number) => (
                <div key={index} className="bg-white/5 border border-white/10 rounded-lg p-3">
                  <div className="flex justify-between items-start mb-1">
                    <div className="text-sm font-medium">{event.name}</div>
                    <div className="text-xs text-white/40">
                      {formatTime(event.timestamp)}
                    </div>
                  </div>
                  {event.attributes && Object.keys(event.attributes).length > 0 && (
                    <div className="text-xs text-white/60 mt-2">
                      {JSON.stringify(event.attributes, null, 2)}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    );
  };

  useEffect(() => {
    loadReplayData();

    return () => {
      if (eventSubscriptionRef.current) {
        eventSubscriptionRef.current();
      }
    };
  }, [recordId]);

  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };

    document.addEventListener('fullscreenchange', handleFullscreenChange);
    return () => {
      document.removeEventListener('fullscreenchange', handleFullscreenChange);
    };
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4" />
          <div className="text-white/60">Loading replay data...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-500/10 border border-red-500/30 rounded-lg p-6">
        <div className="flex items-center gap-3 mb-3">
          <AlertCircle className="text-red-400" size={24} />
          <div className="text-lg font-bold">Failed to load replay</div>
        </div>
        <div className="text-white/80 mb-4">{error}</div>
        <button
          onClick={loadReplayData}
          className="px-4 py-2 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
        >
          Retry
        </button>
      </div>
    );
  }

  if (!replayState) {
    return null;
  }

  const currentSpan = getCurrentSpan();
  const speedOptions = [0.25, 0.5, 1, 2, 4, 8];

  return (
    <div 
      ref={playerRef}
      className={cn(
        "bg-black/40 border border-white/10 rounded-lg overflow-hidden",
        isFullscreen && "fixed inset-0 z-50 bg-black"
      )}
    >
      {/* Header */}
      <div className="border-b border-white/10 p-4 flex justify-between items-center">
        <div>
          <h3 className="text-lg font-bold">
            Replay: {replayState.metadata?.record_name || 'Execution'}
          </h3>
          <div className="text-sm text-white/60 flex items-center gap-4">
            <span>Status: {replayState.status}</span>
            <span>Speed: {replayState.playback_speed}x</span>
            <span>Progress: {replayState.progress.toFixed(1)}%</span>
          </div>
        </div>
        
        <div className="flex items-center gap-2">
          <button
            onClick={() => setIsMuted(!isMuted)}
            className="p-2 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
            title={isMuted ? "Unmute" : "Mute"}
          >
            {isMuted ? <VolumeX size={16} /> : <Volume2 size={16} />}
          </button>
          
          <button
            onClick={() => setShowSettings(!showSettings)}
            className="p-2 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
            title="Settings"
          >
            <Settings size={16} />
          </button>
          
          <button
            onClick={handleFullscreen}
            className="p-2 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
            title={isFullscreen ? "Exit Fullscreen" : "Fullscreen"}
          >
            {isFullscreen ? <Minimize size={16} /> : <Expand size={16} />}
          </button>
          
          {onClose && (
            <button
              onClick={onClose}
              className="p-2 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
              title="Close"
            >
              <X size={16} />
            </button>
          )}
        </div>
      </div>

      {/* Main Content */}
      <div className="flex h-[calc(100vh-200px)] min-h-[600px]">
        {/* Timeline */}
        <div className="w-1/3 border-r border-white/10 p-4 overflow-y-auto">
          <div className="mb-4">
            <h4 className="text-sm font-bold mb-2">Timeline</h4>
            <div ref={timelineRef} className="space-y-1">
              {spans.map((span, index) => (
                <button
                  key={span.id}
                  onClick={() => {
                    const progress = (index / Math.max(1, spans.length - 1)) * 100;
                    handleSeek(progress);
                  }}
                  className={cn(
                    "w-full text-left p-3 rounded-lg border transition-colors",
                    index === currentSpanIndex
                      ? "bg-primary/20 border-primary/30"
                      : "bg-white/5 border-white/10 hover:bg-white/10"
                  )}
                >
                  <div className="flex items-center gap-2 mb-1">
                    <div className={cn("w-2 h-2 rounded-full", getSpanColor(span))} />
                    <div className="text-sm font-medium truncate">{span.name}</div>
                  </div>
                  <div className="text-xs text-white/60 flex justify-between">
                    <span>{span.agent_name || span.tool_name || 'System'}</span>
                    <span>{formatDuration(span.duration)}</span>
                  </div>
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Player and Details */}
        <div className="flex-1 flex flex-col">
          {/* Player Controls */}
          <div className="border-b border-white/10 p-4">
            <div className="flex items-center justify-center gap-4 mb-4">
              <button
                onClick={() => handleControl('seek', { progress: 0 })}
                className="p-3 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
                title="Start"
              >
                <SkipBack size={20} />
              </button>
              
              <button
                onClick={() => handleSpeedChange(replayState.playback_speed / 2)}
                className="p-3 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
                title="Slower"
              >
                <Rewind size={20} />
              </button>
              
              {replayState.status === 'playing' ? (
                <button
                  onClick={() => handleControl('pause')}
                  className="p-4 bg-primary/20 border border-primary/30 rounded-lg hover:bg-primary/30 transition-colors"
                  title="Pause"
                >
                  <Pause size={24} className="text-primary" />
                </button>
              ) : (
                <button
                  onClick={() => handleControl('resume')}
                  className="p-4 bg-primary border border-primary rounded-lg hover:bg-primary/90 transition-colors"
                  title="Play"
                >
                  <Play size={24} className="text-black" />
                </button>
              )}
              
              <button
                onClick={() => handleSpeedChange(replayState.playback_speed * 2)}
                className="p-3 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
                title="Faster"
              >
                <FastForward size={20} />
              </button>
              
              <button
                onClick={() => handleControl('stop')}
                className="p-3 bg-white/5 border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
                title="Stop"
              >
                <StopCircle size={20} />
              </button>
            </div>

            {/* Progress Bar */}
            <div className="space-y-2">
              <div className="flex justify-between text-sm text-white/60">
                <span>{formatTime(replayState.current_time)}</span>
                <span>{formatTime(replayState.end_time)}</span>
              </div>
              <div className="relative">
                <div 
                  className="absolute top-0 left-0 h-1 bg-primary transition-all"
                  style={{ width: `${replayState.progress}%` }}
                />
                <div className="w-full h-1 bg-white/10 rounded-full" />
                <input
                  type="range"
                  min="0"
                  max="100"
                  step="0.1"
                  value={replayState.progress}
                  onChange={(e) => handleSeek(parseFloat(e.target.value))}
                  className="absolute top-0 left-0 w-full h-1 opacity-0 cursor-pointer"
                />
              </div>
              <div className="flex justify-between text-xs text-white/40">
                <span>Span {currentSpanIndex + 1} of {spans.length}</span>
                <span>{replayState.progress.toFixed(1)}%</span>
              </div>
            </div>

            {/* Speed Controls */}
            {showSettings && (
              <div className="mt-4 p-4 bg-white/5 border border-white/10 rounded-lg">
                <div className="text-sm font-medium mb-2">Playback Speed</div>
                <div className="flex flex-wrap gap-2">
                  {speedOptions.map((speed) => (
                    <button
                      key={speed}
                      onClick={() => handleSpeedChange(speed)}
                      className={cn(
                        "px-3 py-1 rounded-lg border transition-colors",
                        replayState.playback_speed === speed
                          ? "bg-primary/20 text-primary border-primary/30"
                          : "bg-white/5 text-white/60 border-white/10 hover:bg-white/10"
                      )}
                    >
                      {speed}x
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Current Span Details */}
          <div className="flex-1 overflow-y-auto p-4">
            {currentSpan ? (
              renderSpanDetails(currentSpan)
            ) : (
              <div className="flex items-center justify-center h-full">
                <div className="text-center text-white/60">
                  <Clock size={48} className="mx-auto mb-4" />
                  <div className="text-lg">No active span</div>
                  <div className="text-sm">Select a span from the timeline or start playback</div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Footer Stats */}
      <div className="border-t border-white/10 p-4">
        <div className="grid grid-cols-4 gap-4">
          <div className="text-center">
            <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Current Time</div>
            <div className="text-lg font-bold">{formatTime(replayState.current_time)}</div>
          </div>
          <div className="text-center">
            <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Elapsed</div>
            <div className="text-lg font-bold">
              {formatDuration(new Date(replayState.current_time).getTime() - new Date(replayState.start_time).getTime())}
            </div>
          </div>
          <div className="text-center">
            <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Remaining</div>
            <div className="text-lg font-bold">
              {formatDuration(new Date(replayState.end_time).getTime() - new Date(replayState.current_time).getTime())}
            </div>
          </div>
          <div className="text-center">
            <div className="text-xs text-white/40 uppercase tracking-wider mb-1">Speed</div>
            <div className="text-lg font-bold">{replayState.playback_speed}x</div>
          </div>
        </div>
      </div>
    </div>
  );
}