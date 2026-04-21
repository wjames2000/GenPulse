import React, { useState, useEffect } from 'react';
import {
  Brain,
  Search,
  Filter,
  RefreshCw,
  Calendar,
  Clock,
  CheckCircle,
  XCircle,
  User,
  Database,
  Layers,
  BarChart3,
  TrendingUp,
  Hash,
  Tag,
  Zap,
  ChevronRight,
  Wrench
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';
import { api, EpisodicMemory, SemanticMemory } from '../services/api';

export default function MemoryViewer() {
  const [activeTab, setActiveTab] = useState<'episodic' | 'semantic' | 'stats'>('episodic');
  const [episodicMemories, setEpisodicMemories] = useState<EpisodicMemory[]>([]);
  const [semanticMemory, setSemanticMemory] = useState<SemanticMemory | null>(null);
  const [memoryStats, setMemoryStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedMemory, setSelectedMemory] = useState<EpisodicMemory | null>(null);

  useEffect(() => {
    loadData();
  }, [activeTab]);

  const loadData = async () => {
    setLoading(true);
    try {
      switch (activeTab) {
        case 'episodic':
          const memories = await api.getEpisodicMemories(searchQuery, 20);
          setEpisodicMemories(memories);
          if (memories.length > 0 && !selectedMemory) {
            setSelectedMemory(memories[0]);
          }
          break;
        case 'semantic':
          const semantic = await api.getSemanticMemory();
          setSemanticMemory(semantic);
          break;
        case 'stats':
          const stats = await api.getMemoryStats();
          setMemoryStats(stats);
          break;
      }
    } catch (error) {
      console.error('Failed to load memory data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = () => {
    if (activeTab === 'episodic') {
      loadData();
    }
  };

  const getAgentIcon = (agentName: string) => {
    if (agentName.includes('前端')) return '🎨';
    if (agentName.includes('后端')) return '⚙️';
    if (agentName.includes('质量保证') || agentName.includes('QA')) return '🧪';
    if (agentName.includes('架构')) return '🏗️';
    if (agentName.includes('运维')) return '🚀';
    return '🤖';
  };

  const formatDuration = (ms: number) => {
    if (ms < 1000) return `${ms}ms`;
    return `${(ms / 1000).toFixed(2)}s`;
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className="flex-1 overflow-hidden flex flex-col p-8 gap-8 h-full bg-background font-sans">
      {/* Header */}
      <header className="flex items-end justify-between shrink-0">
        <div>
          <h2 className="text-5xl font-black text-on-surface tracking-tight leading-tight mb-2">记忆系统</h2>
          <p className="text-outline text-base font-medium">查看和管理Agent的学习经验与用户画像。</p>
        </div>
        <div className="flex gap-3">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-outline-variant" size={16} />
            <input
              type="text"
              placeholder="搜索记忆..."
              className="pl-10 pr-4 py-2.5 rounded-xl bg-surface-container hover:bg-surface-container-highest text-on-surface-variant text-xs font-bold border border-outline-variant/10 shadow-sm transition-all focus:outline-none focus:ring-2 focus:ring-primary/20"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
            />
          </div>
          <button
            className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-surface-container hover:bg-surface-container-highest text-on-surface-variant text-xs font-bold border border-outline-variant/10 shadow-sm transition-all active:scale-95 uppercase tracking-widest"
            onClick={loadData}
            disabled={loading}
          >
            <RefreshCw size={16} className={loading ? 'animate-spin' : ''} />
            {loading ? '加载中...' : '刷新'}
          </button>
        </div>
      </header>

      {/* Tabs */}
      <div className="flex gap-2 border-b border-outline-variant/10 pb-2">
        <button
          className={cn(
            "flex items-center gap-2 px-5 py-2.5 rounded-xl text-xs font-bold tracking-widest transition-all",
            activeTab === 'episodic'
              ? "bg-primary text-white shadow-lg"
              : "bg-surface-container hover:bg-surface-container-highest text-on-surface-variant border border-outline-variant/10"
          )}
          onClick={() => setActiveTab('episodic')}
        >
          <Database size={14} />
          情节记忆
        </button>
        <button
          className={cn(
            "flex items-center gap-2 px-5 py-2.5 rounded-xl text-xs font-bold tracking-widest transition-all",
            activeTab === 'semantic'
              ? "bg-primary text-white shadow-lg"
              : "bg-surface-container hover:bg-surface-container-highest text-on-surface-variant border border-outline-variant/10"
          )}
          onClick={() => setActiveTab('semantic')}
        >
          <User size={14} />
          语义记忆
        </button>
        <button
          className={cn(
            "flex items-center gap-2 px-5 py-2.5 rounded-xl text-xs font-bold tracking-widest transition-all",
            activeTab === 'stats'
              ? "bg-primary text-white shadow-lg"
              : "bg-surface-container hover:bg-surface-container-highest text-on-surface-variant border border-outline-variant/10"
          )}
          onClick={() => setActiveTab('stats')}
        >
          <BarChart3 size={14} />
          统计信息
        </button>
      </div>

      {/* Content */}
      <div className="flex-1 flex gap-8 overflow-hidden min-h-0">
        {/* Left Pane: List/Overview */}
        <div className="w-[420px] shrink-0 flex flex-col gap-5 overflow-y-auto pr-3 pb-8 custom-scrollbar">
          {loading ? (
            <div className="flex items-center justify-center h-32">
              <RefreshCw className="animate-spin text-primary" size={24} />
            </div>
          ) : (
            <>
              {activeTab === 'episodic' && (
                <>
                  {episodicMemories.length === 0 ? (
                    <div className="text-center py-8 text-outline-variant">
                      <p>没有找到记忆记录</p>
                    </div>
                  ) : (
                    episodicMemories.map((memory) => (
                      <div
                        key={memory.id}
                        className={cn(
                          "rounded-3xl p-6 flex flex-col gap-4 relative cursor-pointer border transition-all group",
                          selectedMemory?.id === memory.id
                            ? "bg-surface-container-highest/60 border-primary/20 shadow-2xl hover:bg-surface-container-highest"
                            : "bg-surface-container-low/40 border-outline-variant/5 hover:border-outline-variant/20 hover:bg-surface-container-high"
                        )}
                        onClick={() => setSelectedMemory(memory)}
                      >
                        {/* Selection Indicator */}
                        {selectedMemory?.id === memory.id && (
                          <div className="absolute left-0 top-1/2 -translate-y-1/2 w-1.5 h-12 bg-primary rounded-r-full shadow-[0_0_15px_rgba(91,95,255,0.8)]" />
                        )}

                        <div className="flex justify-between items-start">
                          <div className="flex items-center gap-3">
                            <div className="text-2xl">{getAgentIcon(memory.agent_name)}</div>
                            <div>
                              <h4 className={cn(
                                "text-base font-bold",
                                selectedMemory?.id === memory.id ? "text-on-surface" : "text-on-surface-variant group-hover:text-on-surface"
                              )}>
                                {memory.agent_name}
                              </h4>
                              <p className="text-sm text-outline opacity-80">{memory.task_type}</p>
                            </div>
                          </div>
                          <div className="flex items-center gap-2">
                            {memory.success ? (
                              <CheckCircle size={14} className="text-green-500" />
                            ) : (
                              <XCircle size={14} className="text-red-500" />
                            )}
                            <span className="text-[10px] font-mono font-bold text-outline opacity-70">
                              {formatDuration(memory.duration_ms)}
                            </span>
                          </div>
                        </div>

                        <p className={cn(
                          "text-sm leading-relaxed line-clamp-2",
                          selectedMemory?.id === memory.id ? "text-on-surface-variant opacity-80" : "text-outline opacity-80 group-hover:opacity-100"
                        )}>
                          {memory.description}
                        </p>

                        <div className="flex items-center justify-between pt-3 border-t border-outline-variant/10">
                          <div className="flex items-center gap-2">
                            <Calendar size={12} className="text-outline-variant" />
                            <span className="text-[10px] text-outline-variant">
                              {formatDate(memory.created_at)}
                            </span>
                          </div>
                          <div className="flex items-center gap-1">
                            <Hash size={12} className="text-outline-variant" />
                            <span className="text-[10px] font-mono font-bold text-primary">
                              {memory.relevance_score.toFixed(2)}
                            </span>
                          </div>
                        </div>

                        {/* Keywords */}
                        {memory.keywords && memory.keywords.length > 0 && (
                          <div className="flex flex-wrap gap-1">
                            {memory.keywords.slice(0, 3).map((keyword, index) => (
                              <span
                                key={index}
                                className="px-2 py-1 text-[10px] font-bold uppercase tracking-widest rounded-full bg-surface-container-highest/50 text-outline-variant"
                              >
                                {keyword}
                              </span>
                            ))}
                            {memory.keywords.length > 3 && (
                              <span className="px-2 py-1 text-[10px] font-bold uppercase tracking-widest rounded-full bg-surface-container-highest/50 text-outline-variant">
                                +{memory.keywords.length - 3}
                              </span>
                            )}
                          </div>
                        )}
                      </div>
                    ))
                  )}
                </>
              )}

              {activeTab === 'semantic' && semanticMemory && (
                <div className="space-y-6">
                  {/* User Profile Card */}
                  <div className="bg-surface-container-highest/60 rounded-3xl p-6">
                    <div className="flex items-center gap-4 mb-6">
                      <div className="w-16 h-16 rounded-2xl bg-primary-container/20 flex items-center justify-center text-primary shadow-inner border border-primary/20">
                        <User size={28} />
                      </div>
                      <div>
                        <h3 className="text-xl font-black text-on-surface">{semanticMemory.username}</h3>
                        <p className="text-sm text-outline">用户画像</p>
                      </div>
                    </div>

                    <div className="space-y-4">
                      <div>
                        <h4 className="text-sm font-bold text-on-surface-variant mb-2 flex items-center gap-2">
                          <Brain size={14} />
                          技能领域
                        </h4>
                        <div className="flex flex-wrap gap-2">
                          {semanticMemory.skills.map((skill, index) => (
                            <span
                              key={index}
                              className="px-3 py-1.5 text-xs font-bold rounded-xl bg-surface-container-high text-on-surface-variant"
                            >
                              {skill}
                            </span>
                          ))}
                        </div>
                      </div>

                      <div>
                        <h4 className="text-sm font-bold text-on-surface-variant mb-2 flex items-center gap-2">
                          <TrendingUp size={14} />
                          知识领域
                        </h4>
                        <div className="space-y-2">
                          {Object.entries(semanticMemory.knowledge_areas || {}).map(([area, score]) => (
                            <div key={area} className="flex items-center justify-between">
                              <span className="text-sm text-outline">{area}</span>
                              <div className="flex items-center gap-2">
                                <div className="w-24 h-2 rounded-full bg-surface-container-high overflow-hidden">
                                  <div
                                    className="h-full bg-primary rounded-full"
                                    style={{ width: `${score}%` }}
                                  />
                                </div>
                                <span className="text-xs font-mono font-bold text-primary">{score}%</span>
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Goals Card */}
                  {semanticMemory.goals && semanticMemory.goals.length > 0 && (
                    <div className="bg-surface-container-low/40 rounded-3xl p-6">
                      <h4 className="text-sm font-bold text-on-surface-variant mb-3 flex items-center gap-2">
                        <Zap size={14} />
                        项目目标
                      </h4>
                      <ul className="space-y-2">
                        {semanticMemory.goals.map((goal, index) => (
                          <li key={index} className="flex items-start gap-2">
                            <ChevronRight size={12} className="text-primary mt-1 flex-shrink-0" />
                            <span className="text-sm text-outline">{goal}</span>
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                </div>
              )}

              {activeTab === 'stats' && memoryStats && (
                <div className="space-y-6">
                  {/* Stats Overview */}
                  <div className="bg-surface-container-highest/60 rounded-3xl p-6">
                    <h3 className="text-lg font-black text-on-surface mb-4">记忆系统概览</h3>
                    <div className="grid grid-cols-2 gap-4">
                      <StatCard
                        icon={Database}
                        label="总记忆数"
                        value={memoryStats.total_memories}
                        color="text-primary"
                      />
                      <StatCard
                        icon={Layers}
                        label="情节记忆"
                        value={memoryStats.episodic_memories}
                        color="text-secondary"
                      />
                      <StatCard
                        icon={User}
                        label="语义记忆"
                        value={memoryStats.semantic_memories}
                        color="text-tertiary"
                      />
                      <StatCard
                        icon={BarChart3}
                        label="缓存命中率"
                        value={`${((memoryStats.cache_hit_rate || 0) * 100).toFixed(1)}%`}
                        color="text-green-500"
                      />
                    </div>
                  </div>

                  {/* Performance Stats */}
                  <div className="bg-surface-container-low/40 rounded-3xl p-6">
                    <h4 className="text-sm font-bold text-on-surface-variant mb-3">性能指标</h4>
                    <div className="space-y-3">
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-outline">平均搜索时间</span>
                        <span className="text-sm font-mono font-bold text-primary">
                          {memoryStats.avg_search_time_ms?.toFixed(1)}ms
                        </span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-outline">总搜索次数</span>
                        <span className="text-sm font-mono font-bold text-primary">
                          {memoryStats.total_searches}
                        </span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-outline">最后更新</span>
                        <span className="text-sm font-mono font-bold text-outline-variant">
                          {new Date(memoryStats.last_updated).toLocaleDateString('zh-CN')}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </>
          )}
        </div>

        {/* Right Pane: Details */}
        <div className="flex-1 flex flex-col bg-surface-container-lowest/40 rounded-3xl border border-outline-variant/10 overflow-hidden shadow-2xl backdrop-blur-xl">
          <div className="bg-surface-container-highest/50 border-b border-outline-variant/10 px-6 py-4 shrink-0">
            <h3 className="text-lg font-black text-on-surface">
              {activeTab === 'episodic' && selectedMemory
                ? '记忆详情'
                : activeTab === 'semantic'
                ? '用户画像详情'
                : '统计详情'}
            </h3>
          </div>

          <div className="flex-1 p-8 overflow-y-auto custom-scrollbar">
            {activeTab === 'episodic' && selectedMemory ? (
              <div className="space-y-6">
                {/* Memory Header */}
                <div className="flex items-start justify-between">
                  <div>
                    <h2 className="text-2xl font-black text-on-surface mb-2">{selectedMemory.description}</h2>
                    <div className="flex items-center gap-4">
                      <div className="flex items-center gap-2">
                        <div className="text-2xl">{getAgentIcon(selectedMemory.agent_name)}</div>
                        <span className="text-sm font-bold text-on-surface-variant">{selectedMemory.agent_name}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <Tag size={14} className="text-outline-variant" />
                        <span className="text-sm text-outline-variant">{selectedMemory.task_type}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        {selectedMemory.success ? (
                          <>
                            <CheckCircle size={14} className="text-green-500" />
                            <span className="text-sm text-green-500 font-bold">成功</span>
                          </>
                        ) : (
                          <>
                            <XCircle size={14} className="text-red-500" />
                            <span className="text-sm text-red-500 font-bold">失败</span>
                          </>
                        )}
                      </div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-sm text-outline-variant mb-1">执行时间</div>
                    <div className="text-xl font-mono font-bold text-primary">
                      {formatDuration(selectedMemory.duration_ms)}
                    </div>
                  </div>
                </div>

                {/* Metadata Grid */}
                <div className="grid grid-cols-2 gap-6">
                  <div className="bg-surface-container-high/50 rounded-2xl p-5">
                    <h4 className="text-sm font-bold text-on-surface-variant mb-3 flex items-center gap-2">
                      <Calendar size={14} />
                      时间信息
                    </h4>
                    <div className="space-y-2">
                      <div className="flex justify-between">
                        <span className="text-sm text-outline">创建时间</span>
                        <span className="text-sm font-mono font-bold text-on-surface">
                          {new Date(selectedMemory.created_at).toLocaleString('zh-CN')}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-sm text-outline">任务ID</span>
                        <span className="text-sm font-mono font-bold text-outline-variant">
                          {selectedMemory.task_id}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-sm text-outline">相关度评分</span>
                        <span className="text-sm font-mono font-bold text-primary">
                          {selectedMemory.relevance_score.toFixed(3)}
                        </span>
                      </div>
                    </div>
                  </div>

                  <div className="bg-surface-container-high/50 rounded-2xl p-5">
                    <h4 className="text-sm font-bold text-on-surface-variant mb-3 flex items-center gap-2">
                      <Wrench size={14} />
                      工具使用
                    </h4>
                    <div className="space-y-2">
                      {Object.entries(selectedMemory.tool_usage || {}).map(([tool, count]) => (
                        <div key={tool} className="flex justify-between">
                          <span className="text-sm text-outline">{tool}</span>
                          <span className="text-sm font-mono font-bold text-primary">{count}次</span>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>

                {/* Keywords */}
                {selectedMemory.keywords && selectedMemory.keywords.length > 0 && (
                  <div className="bg-surface-container-high/50 rounded-2xl p-5">
                    <h4 className="text-sm font-bold text-on-surface-variant mb-3">关键词</h4>
                    <div className="flex flex-wrap gap-2">
                      {selectedMemory.keywords.map((keyword, index) => (
                        <span
                          key={index}
                          className="px-3 py-1.5 text-xs font-bold rounded-xl bg-surface-container-highest text-on-surface-variant"
                        >
                          {keyword}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Context Data */}
                {selectedMemory.context_data && Object.keys(selectedMemory.context_data).length > 0 && (
                  <div className="bg-surface-container-high/50 rounded-2xl p-5">
                    <h4 className="text-sm font-bold text-on-surface-variant mb-3">上下文数据</h4>
                    <pre className="text-sm text-outline font-mono overflow-x-auto">
                      {JSON.stringify(selectedMemory.context_data, null, 2)}
                    </pre>
                  </div>
                )}
              </div>
            ) : activeTab === 'semantic' && semanticMemory ? (
              <div className="space-y-6">
                {/* Preferences */}
                {semanticMemory.preferences && Object.keys(semanticMemory.preferences).length > 0 && (
                  <div className="bg-surface-container-high/50 rounded-2xl p-5">
                    <h4 className="text-sm font-bold text-on-surface-variant mb-3">偏好设置</h4>
                    <div className="grid grid-cols-2 gap-4">
                      {Object.entries(semanticMemory.preferences).map(([key, value]) => (
                        <div key={key} className="flex justify-between">
                          <span className="text-sm text-outline">{key}</span>
                          <span className="text-sm font-mono font-bold text-primary">{String(value)}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Interests */}
                {semanticMemory.interests && semanticMemory.interests.length > 0 && (
                  <div className="bg-surface-container-high/50 rounded-2xl p-5">
                    <h4 className="text-sm font-bold text-on-surface-variant mb-3">兴趣领域</h4>
                    <div className="flex flex-wrap gap-2">
                      {semanticMemory.interests.map((interest, index) => (
                        <span
                          key={index}
                          className="px-3 py-1.5 text-xs font-bold rounded-xl bg-surface-container-highest text-on-surface-variant"
                        >
                          {interest}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Working Style */}
                {semanticMemory.working_style && Object.keys(semanticMemory.working_style).length > 0 && (
                  <div className="bg-surface-container-high/50 rounded-2xl p-5">
                    <h4 className="text-sm font-bold text-on-surface-variant mb-3">工作风格</h4>
                    <div className="grid grid-cols-2 gap-4">
                      {Object.entries(semanticMemory.working_style).map(([key, value]) => (
                        <div key={key} className="flex justify-between">
                          <span className="text-sm text-outline">{key}</span>
                          <span className="text-sm font-mono font-bold text-primary">{String(value)}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Timeline */}
                <div className="bg-surface-container-high/50 rounded-2xl p-5">
                  <h4 className="text-sm font-bold text-on-surface-variant mb-3">时间线</h4>
                  <div className="space-y-3">
                    <div className="flex justify-between">
                      <span className="text-sm text-outline">用户ID</span>
                      <span className="text-sm font-mono font-bold text-outline-variant">
                        {semanticMemory.user_id}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-sm text-outline">创建时间</span>
                      <span className="text-sm font-mono font-bold text-on-surface">
                        {new Date(semanticMemory.created_at).toLocaleString('zh-CN')}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-sm text-outline">最后更新</span>
                      <span className="text-sm font-mono font-bold text-on-surface">
                        {new Date(semanticMemory.last_updated).toLocaleString('zh-CN')}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            ) : activeTab === 'stats' && memoryStats ? (
              <div className="space-y-6">
                {/* Detailed Stats */}
                <div className="grid grid-cols-2 gap-6">
                  <div className="bg-surface-container-high/50 rounded-2xl p-5">
                    <h4 className="text-sm font-bold text-on-surface-variant mb-3">记忆分布</h4>
                    <div className="space-y-3">
                      <div className="flex justify-between">
                        <span className="text-sm text-outline">情节记忆占比</span>
                        <span className="text-sm font-mono font-bold text-primary">
                          {((memoryStats.episodic_memories / memoryStats.total_memories) * 100).toFixed(1)}%
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-sm text-outline">语义记忆占比</span>
                        <span className="text-sm font-mono font-bold text-primary">
                          {((memoryStats.semantic_memories / memoryStats.total_memories) * 100).toFixed(1)}%
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-sm text-outline">工作记忆占比</span>
                        <span className="text-sm font-mono font-bold text-primary">
                          {((memoryStats.working_memories / memoryStats.total_memories) * 100).toFixed(1)}%
                        </span>
                      </div>
                    </div>
                  </div>

                  <div className="bg-surface-container-high/50 rounded-2xl p-5">
                    <h4 className="text-sm font-bold text-on-surface-variant mb-3">性能指标</h4>
                    <div className="space-y-3">
                      <div className="flex justify-between">
                        <span className="text-sm text-outline">搜索成功率</span>
                        <span className="text-sm font-mono font-bold text-green-500">
                          {(memoryStats.cache_hit_rate * 100).toFixed(1)}%
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-sm text-outline">平均响应时间</span>
                        <span className="text-sm font-mono font-bold text-primary">
                          {memoryStats.avg_search_time_ms?.toFixed(1)}ms
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-sm text-outline">日均搜索量</span>
                        <span className="text-sm font-mono font-bold text-primary">
                          {Math.round(memoryStats.total_searches / 30)}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Recommendations */}
                <div className="bg-surface-container-high/50 rounded-2xl p-5">
                  <h4 className="text-sm font-bold text-on-surface-variant mb-3">优化建议</h4>
                  <ul className="space-y-2">
                    <li className="flex items-start gap-2">
                      <ChevronRight size={12} className="text-primary mt-1 flex-shrink-0" />
                      <span className="text-sm text-outline">
                        缓存命中率较高，考虑增加缓存容量以进一步提升性能
                      </span>
                    </li>
                    <li className="flex items-start gap-2">
                      <ChevronRight size={12} className="text-primary mt-1 flex-shrink-0" />
                      <span className="text-sm text-outline">
                        情节记忆数量充足，建议定期清理低相关度的记忆
                      </span>
                    </li>
                    <li className="flex items-start gap-2">
                      <ChevronRight size={12} className="text-primary mt-1 flex-shrink-0" />
                      <span className="text-sm text-outline">
                        语义记忆相对较少，建议收集更多用户偏好数据
                      </span>
                    </li>
                  </ul>
                </div>
              </div>
            ) : (
              <div className="flex items-center justify-center h-full text-outline-variant">
                {activeTab === 'episodic' ? '请选择一个记忆查看详情' : '加载中...'}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

function StatCard({ icon: Icon, label, value, color }: any) {
  return (
    <div className="flex flex-col items-center p-3 rounded-xl bg-surface-container-high">
      <Icon size={20} className={cn("mb-2", color)} />
      <span className="text-[10px] uppercase font-black tracking-widest text-outline-variant text-center">
        {label}
      </span>
      <span className={cn("text-lg font-mono font-bold leading-none mt-1", color)}>
        {value}
      </span>
    </div>
  );
}