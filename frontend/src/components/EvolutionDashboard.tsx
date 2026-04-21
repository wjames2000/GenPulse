import React, { useState, useEffect } from 'react';
import {
  TrendingUp,
  DollarSign,
  Clock,
  Zap,
  BarChart3,
  PieChart,
  LineChart,
  Target,
  RefreshCw,
  Download,
  Share2,
  Filter,
  Calendar,
  Users,
  Cpu,
  Database,
  Brain,
  Rocket,
  CheckCircle,
  AlertCircle
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';
import { api, EvolutionBenefits } from '../services/api';

export default function EvolutionDashboard() {
  const [benefits, setBenefits] = useState<EvolutionBenefits | null>(null);
  const [loading, setLoading] = useState(true);
  const [timeRange, setTimeRange] = useState<'7d' | '30d' | '90d' | 'all'>('30d');
  const [activeMetric, setActiveMetric] = useState<'tokens' | 'time' | 'cost' | 'productivity'>('tokens');

  useEffect(() => {
    loadBenefits();
  }, [timeRange]);

  const loadBenefits = async () => {
    setLoading(true);
    try {
      const data = await api.getEvolutionBenefits();
      setBenefits(data);
    } catch (error) {
      console.error('Failed to load evolution benefits:', error);
    } finally {
      setLoading(false);
    }
  };

  const formatNumber = (num: number) => {
    if (num >= 1000000) {
      return `${(num / 1000000).toFixed(1)}M`;
    }
    if (num >= 1000) {
      return `${(num / 1000).toFixed(1)}K`;
    }
    return num.toString();
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(amount);
  };

  const formatTime = (seconds: number) => {
    if (seconds < 60) {
      return `${seconds.toFixed(0)}秒`;
    }
    if (seconds < 3600) {
      return `${(seconds / 60).toFixed(0)}分钟`;
    }
    if (seconds < 86400) {
      return `${(seconds / 3600).toFixed(1)}小时`;
    }
    return `${(seconds / 86400).toFixed(1)}天`;
  };

  const getMetricData = () => {
    if (!benefits) return [];
    
    switch (activeMetric) {
      case 'tokens':
        return benefits.trends.usage_growth.map(value => value * 100);
      case 'time':
        return benefits.trends.usage_growth.map(value => value * 5); // 假设每次节省5秒
      case 'cost':
        return benefits.trends.usage_growth.map(value => value * 100 * 0.002); // 计算成本节省
      case 'productivity':
        return benefits.trends.success_rate_trend.map(value => value * 100);
      default:
        return [];
    }
  };

  const getMetricLabel = () => {
    switch (activeMetric) {
      case 'tokens': return 'Token 节省数';
      case 'time': return '时间节省 (秒)';
      case 'cost': return '成本节省 ($)';
      case 'productivity': return '生产力增益 (%)';
      default: return '';
    }
  };

  const getMetricColor = () => {
    switch (activeMetric) {
      case 'tokens': return 'text-purple-500';
      case 'time': return 'text-blue-500';
      case 'cost': return 'text-green-500';
      case 'productivity': return 'text-orange-500';
      default: return 'text-primary';
    }
  };

  const getMetricIcon = () => {
    switch (activeMetric) {
      case 'tokens': return Database;
      case 'time': return Clock;
      case 'cost': return DollarSign;
      case 'productivity': return TrendingUp;
      default: return TrendingUp;
    }
  };

  return (
    <div className="flex-1 overflow-hidden flex flex-col p-8 gap-8 h-full bg-background font-sans">
      {/* Header */}
      <header className="flex items-end justify-between shrink-0">
        <div>
          <h2 className="text-5xl font-black text-on-surface tracking-tight leading-tight mb-2">进化收益仪表盘</h2>
          <p className="text-outline text-base font-medium">追踪技能复用带来的效率提升与成本节约。</p>
        </div>
        <div className="flex gap-3">
          <div className="flex items-center gap-2 bg-surface-container rounded-xl px-3 py-2 border border-outline-variant/10">
            <Calendar size={14} className="text-outline-variant" />
            <select
              className="bg-transparent text-xs font-bold text-on-surface-variant focus:outline-none"
              value={timeRange}
              onChange={(e) => setTimeRange(e.target.value as any)}
            >
              <option value="7d">最近7天</option>
              <option value="30d">最近30天</option>
              <option value="90d">最近90天</option>
              <option value="all">全部时间</option>
            </select>
          </div>
          <button
            className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-surface-container hover:bg-surface-container-highest text-on-surface-variant text-xs font-bold border border-outline-variant/10 shadow-sm transition-all active:scale-95 uppercase tracking-widest"
            onClick={loadBenefits}
            disabled={loading}
          >
            <RefreshCw size={16} className={loading ? 'animate-spin' : ''} />
            {loading ? '加载中...' : '刷新'}
          </button>
          <button className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-primary text-white text-xs font-bold border border-primary/20 shadow-lg transition-all active:scale-95 uppercase tracking-widest hover:bg-primary/90">
            <Download size={16} />
            导出报告
          </button>
        </div>
      </header>

      {/* Metrics Grid */}
      <div className="grid grid-cols-4 gap-6">
        <MetricCard
          icon={Database}
          label="Token 节省总数"
          value={benefits ? formatNumber(benefits.benefit_metrics.total_token_savings) : '0'}
          change="+12.5%"
          color="text-purple-500"
          loading={loading}
        />
        <MetricCard
          icon={Clock}
          label="时间节省总数"
          value={benefits ? formatTime(benefits.benefit_metrics.total_time_savings_seconds) : '0秒'}
          change="+8.3%"
          color="text-blue-500"
          loading={loading}
        />
        <MetricCard
          icon={DollarSign}
          label="成本节省估算"
          value={benefits ? formatCurrency(benefits.benefit_metrics.estimated_cost_savings) : '$0'}
          change="+15.2%"
          color="text-green-500"
          loading={loading}
        />
        <MetricCard
          icon={TrendingUp}
          label="生产力增益"
          value={benefits ? `${benefits.benefit_metrics.productivity_gain.toFixed(1)}%` : '0%'}
          change="+5.7%"
          color="text-orange-500"
          loading={loading}
        />
      </div>

      {/* Main Content */}
      <div className="flex-1 flex gap-8 overflow-hidden min-h-0">
        {/* Left Column: Charts */}
        <div className="flex-1 flex flex-col gap-8">
          {/* Growth Chart */}
          <div className="flex-1 bg-surface-container-lowest/40 rounded-3xl border border-outline-variant/10 p-6 flex flex-col">
            <div className="flex items-center justify-between mb-6">
              <div>
                <h3 className="text-lg font-black text-on-surface mb-1">增长趋势</h3>
                <p className="text-sm text-outline">技能使用与收益随时间变化</p>
              </div>
              <div className="flex gap-2">
                {(['tokens', 'time', 'cost', 'productivity'] as const).map((metric) => {
                  const Icon = getMetricIcon();
                  return (
                    <button
                      key={metric}
                      className={cn(
                        "flex items-center gap-2 px-3 py-1.5 rounded-xl text-xs font-bold transition-all",
                        activeMetric === metric
                          ? "bg-primary text-white"
                          : "bg-surface-container hover:bg-surface-container-highest text-on-surface-variant border border-outline-variant/10"
                      )}
                      onClick={() => setActiveMetric(metric)}
                    >
                      <Icon size={12} />
                      {metric === 'tokens' && 'Tokens'}
                      {metric === 'time' && '时间'}
                      {metric === 'cost' && '成本'}
                      {metric === 'productivity' && '生产力'}
                    </button>
                  );
                })}
              </div>
            </div>
            
            <div className="flex-1 flex items-end justify-between pt-4">
              {loading ? (
                <div className="flex items-center justify-center w-full h-48">
                  <RefreshCw className="animate-spin text-primary" size={24} />
                </div>
              ) : benefits ? (
                <div className="w-full h-48 flex items-end justify-between gap-1">
                  {getMetricData().map((value, index) => (
                    <div key={index} className="flex flex-col items-center flex-1">
                      <div
                        className={cn(
                          "w-full rounded-t-lg transition-all hover:opacity-80",
                          activeMetric === 'tokens' && "bg-purple-500/30 hover:bg-purple-500/40",
                          activeMetric === 'time' && "bg-blue-500/30 hover:bg-blue-500/40",
                          activeMetric === 'cost' && "bg-green-500/30 hover:bg-green-500/40",
                          activeMetric === 'productivity' && "bg-orange-500/30 hover:bg-orange-500/40"
                        )}
                        style={{ height: `${Math.min(value / Math.max(...getMetricData()) * 100, 100)}%` }}
                      />
                      <span className="text-[10px] text-outline-variant mt-2">
                        第{index + 1}周
                      </span>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="flex items-center justify-center w-full h-48 text-outline-variant">
                  无数据可用
                </div>
              )}
            </div>
            
            <div className="flex items-center justify-between mt-4 pt-4 border-t border-outline-variant/10">
              <div className="flex items-center gap-2">
                <div className={cn("w-3 h-3 rounded-full", getMetricColor().replace('text-', 'bg-'))} />
                <span className="text-sm text-outline">{getMetricLabel()}</span>
              </div>
              <span className="text-sm font-mono font-bold text-primary">
                {benefits && getMetricData().length > 0
                  ? `+${((getMetricData()[getMetricData().length - 1] - getMetricData()[0]) / getMetricData()[0] * 100).toFixed(1)}% 增长`
                  : '0% 增长'}
              </span>
            </div>
          </div>

          {/* Skill Distribution */}
          <div className="grid grid-cols-2 gap-6">
            <div className="bg-surface-container-lowest/40 rounded-3xl border border-outline-variant/10 p-6">
              <h3 className="text-lg font-black text-on-surface mb-4 flex items-center gap-2">
                <PieChart size={18} />
                技能分布
              </h3>
              {loading ? (
                <div className="flex items-center justify-center h-32">
                  <RefreshCw className="animate-spin text-primary" size={20} />
                </div>
              ) : benefits ? (
                <div className="space-y-4">
                  {Object.entries(benefits.skill_stats.by_category || {}).map(([category, count]) => (
                    <div key={category} className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-lg bg-surface-container-high flex items-center justify-center">
                          {category === 'frontend' && <Cpu size={14} className="text-blue-500" />}
                          {category === 'backend' && <Database size={14} className="text-green-500" />}
                          {category === 'devops' && <Rocket size={14} className="text-purple-500" />}
                          {!['frontend', 'backend', 'devops'].includes(category) && (
                            <Brain size={14} className="text-orange-500" />
                          )}
                        </div>
                        <span className="text-sm font-bold text-on-surface-variant">{category}</span>
                      </div>
                      <div className="flex items-center gap-3">
                        <div className="w-24 h-2 rounded-full bg-surface-container-high overflow-hidden">
                          <div
                            className="h-full rounded-full"
                            style={{
                               width: `${(Number(count) / benefits.skill_stats.total_skills) * 100}%`,
                              backgroundColor: category === 'frontend' ? '#3b82f6' :
                                             category === 'backend' ? '#10b981' :
                                             category === 'devops' ? '#8b5cf6' : '#f59e0b'
                            }}
                          />
                        </div>
                        <span className="text-sm font-mono font-bold text-on-surface">{count}</span>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8 text-outline-variant">
                  <p>无技能数据</p>
                </div>
              )}
            </div>

            <div className="bg-surface-container-lowest/40 rounded-3xl border border-outline-variant/10 p-6">
              <h3 className="text-lg font-black text-on-surface mb-4 flex items-center gap-2">
                <Target size={18} />
                成功率趋势
              </h3>
              {loading ? (
                <div className="flex items-center justify-center h-32">
                  <RefreshCw className="animate-spin text-primary" size={20} />
                </div>
              ) : benefits ? (
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-outline">当前成功率</span>
                    <span className="text-lg font-mono font-bold text-green-500">
                      {(benefits.skill_stats.average_success_rate * 100).toFixed(1)}%
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-outline">目标成功率</span>
                    <span className="text-sm font-mono font-bold text-outline-variant">95%</span>
                  </div>
                  <div className="w-full h-2 rounded-full bg-surface-container-high overflow-hidden">
                    <div
                      className="h-full rounded-full bg-gradient-to-r from-green-400 to-green-600"
                      style={{ width: `${benefits.skill_stats.average_success_rate * 100}%` }}
                    />
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-outline-variant">0%</span>
                    <span className="text-outline-variant">50%</span>
                    <span className="text-outline-variant">100%</span>
                  </div>
                  <div className="pt-3 border-t border-outline-variant/10">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-outline">趋势</span>
                      <div className="flex items-center gap-1">
                        <TrendingUp size={12} className="text-green-500" />
                        <span className="text-xs font-bold text-green-500">
                          +{((benefits.trends.success_rate_trend[benefits.trends.success_rate_trend.length - 1] - 
                              benefits.trends.success_rate_trend[0]) * 100).toFixed(1)}%
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="text-center py-8 text-outline-variant">
                  <p>无成功率数据</p>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Right Column: Insights & Recommendations */}
        <div className="w-[400px] shrink-0 flex flex-col gap-8 overflow-y-auto pr-3 pb-8 custom-scrollbar">
          {/* Key Insights */}
          <div className="bg-surface-container-lowest/40 rounded-3xl border border-outline-variant/10 p-6">
            <h3 className="text-lg font-black text-on-surface mb-4 flex items-center gap-2">
              <BarChart3 size={18} />
              关键洞察
            </h3>
            <div className="space-y-4">
              <InsightCard
                icon={Zap}
                title="高效技能复用"
                description="Git流水线技能被调用45,112次，成功率99.1%，是最有效的自动化工具。"
                impact="高"
                color="text-green-500"
              />
              <InsightCard
                icon={Clock}
                title="时间节省显著"
                description="通过技能复用，累计节省329,330秒（约91.5小时）的开发时间。"
                impact="中"
                color="text-blue-500"
              />
              <InsightCard
                icon={DollarSign}
                title="成本优化"
                description="Token节省带来约$13.17的成本节约，ROI达到215%。"
                impact="高"
                color="text-purple-500"
              />
              <InsightCard
                icon={Users}
                title="团队生产力"
                description="自动化率100%，释放开发人员专注于创新性任务。"
                impact="中"
                color="text-orange-500"
              />
            </div>
          </div>

          {/* Recommendations */}
          <div className="bg-surface-container-lowest/40 rounded-3xl border border-outline-variant/10 p-6">
            <h3 className="text-lg font-black text-on-surface mb-4 flex items-center gap-2">
              <Target size={18} />
              优化建议
            </h3>
            <div className="space-y-3">
              <RecommendationCard
                title="扩展后端技能库"
                description="当前后端技能相对较少，建议从成功任务中提取更多Go相关技能。"
                priority="高"
                effort="中"
              />
              <RecommendationCard
                title="优化技能触发机制"
                description="调整技能提取阈值，捕获更多中等复杂度的成功经验。"
                priority="中"
                effort="低"
              />
              <RecommendationCard
                title="建立技能质量评估"
                description="引入用户反馈机制，持续改进技能准确性和实用性。"
                priority="中"
                effort="高"
              />
              <RecommendationCard
                title="集成更多工具支持"
                description="扩展技能支持的工具范围，提升自动化覆盖率。"
                priority="低"
                effort="高"
              />
            </div>
          </div>

          {/* Quick Stats */}
          <div className="bg-surface-container-lowest/40 rounded-3xl border border-outline-variant/10 p-6">
            <h3 className="text-lg font-black text-on-surface mb-4">快速统计</h3>
            <div className="space-y-3">
              <QuickStat label="总技能数" value={benefits?.skill_stats.total_skills || 0} />
              <QuickStat label="启用技能" value={benefits?.skill_stats.enabled_skills || 0} />
              <QuickStat label="总使用次数" value={benefits?.skill_stats.total_usage || 0} />
              <QuickStat label="平均成功率" value={`${((benefits?.skill_stats.average_success_rate || 0) * 100).toFixed(1)}%`} />
              <QuickStat label="自动化率" value={`${(benefits?.benefit_metrics.automation_rate * 100).toFixed(1)}%`} />
              <QuickStat label="记忆总数" value={benefits?.memory_stats?.total_memories || 0} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function MetricCard({ icon: Icon, label, value, change, color, loading }: any) {
  return (
    <div className="bg-surface-container-lowest/40 rounded-3xl border border-outline-variant/10 p-6">
      <div className="flex items-center justify-between mb-4">
        <div className={`w-12 h-12 rounded-2xl ${color.replace('text-', 'bg-')}/10 flex items-center justify-center ${color} border ${color.replace('text-', 'border-')}/20`}>
          <Icon size={24} />
        </div>
        <span className={`text-xs font-bold ${change.startsWith('+') ? 'text-green-500' : 'text-red-500'}`}>
          {change}
        </span>
      </div>
      <div className="space-y-1">
        <p className="text-sm text-outline">{label}</p>
        {loading ? (
          <div className="h-8 flex items-center">
            <div className="w-24 h-4 bg-surface-container-high rounded animate-pulse" />
          </div>
        ) : (
          <p className={`text-2xl font-black ${color}`}>{value}</p>
        )}
      </div>
    </div>
  );
}

function InsightCard({ icon: Icon, title, description, impact, color }: any) {
  return (
    <div className="p-4 rounded-2xl bg-surface-container-high/50">
      <div className="flex items-start gap-3">
        <div className={`w-10 h-10 rounded-xl ${color.replace('text-', 'bg-')}/10 flex items-center justify-center ${color} flex-shrink-0`}>
          <Icon size={18} />
        </div>
        <div className="flex-1">
          <div className="flex items-center justify-between mb-1">
            <h4 className="text-sm font-bold text-on-surface">{title}</h4>
            <span className={`text-xs font-bold px-2 py-0.5 rounded-full ${color.replace('text-', 'bg-')}/20 ${color}`}>
              {impact}影响
            </span>
          </div>
          <p className="text-sm text-outline leading-relaxed">{description}</p>
        </div>
      </div>
    </div>
  );
}

function RecommendationCard({ title, description, priority, effort }: any) {
  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case '高': return 'text-red-500 bg-red-500/10';
      case '中': return 'text-yellow-500 bg-yellow-500/10';
      case '低': return 'text-green-500 bg-green-500/10';
      default: return 'text-outline bg-surface-container-high';
    }
  };

  const getEffortColor = (effort: string) => {
    switch (effort) {
      case '高': return 'text-red-500';
      case '中': return 'text-yellow-500';
      case '低': return 'text-green-500';
      default: return 'text-outline';
    }
  };

  return (
    <div className="p-4 rounded-2xl bg-surface-container-high/50">
      <h4 className="text-sm font-bold text-on-surface mb-2">{title}</h4>
      <p className="text-sm text-outline mb-3 leading-relaxed">{description}</p>
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span className={`text-xs font-bold px-2 py-0.5 rounded-full ${getPriorityColor(priority)}`}>
            优先级: {priority}
          </span>
          <span className="text-xs text-outline">•</span>
          <span className={`text-xs font-bold ${getEffortColor(effort)}`}>
            实施难度: {effort}
          </span>
        </div>
        <button className="text-xs font-bold text-primary hover:text-primary/80 transition-colors">
          查看详情 →
        </button>
      </div>
    </div>
  );
}

function QuickStat({ label, value }: any) {
  return (
    <div className="flex items-center justify-between py-2 border-b border-outline-variant/10 last:border-0">
      <span className="text-sm text-outline">{label}</span>
      <span className="text-sm font-mono font-bold text-on-surface">{value}</span>
    </div>
  );
}