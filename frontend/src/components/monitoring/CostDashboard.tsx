import React, { useState, useMemo } from 'react';
import { 
  DollarSign, 
  TrendingUp, 
  TrendingDown, 
  PieChart, 
  BarChart3, 
  LineChart,
  Calendar,
  Filter,
  Download,
  RefreshCw,
  ChevronUp,
  ChevronDown,
  Eye,
  EyeOff,
  Settings,
  AlertCircle,
  CheckCircle2,
  Clock,
  Users,
  Cpu,
  Database,
  Server,
  Brain,
  MessageSquare,
  Code,
  FileText,
  Terminal,
  GitBranch,
  Layout,
  Shield,
  Zap,
  Target,
  Award,
  Trophy,
  Crown,
  Gem,
  Diamond,
  Coins,
  Banknote,
  CreditCard,
  Wallet,
  Receipt,
  ShoppingCart,
  Package,
  Box
} from 'lucide-react';
import { cn } from '../../utils';
import { CostMetric } from '../../types';

interface CostDashboardProps {
  metrics: CostMetric[];
}

export default function CostDashboard({ metrics }: CostDashboardProps) {
  const [timeRange, setTimeRange] = useState<'today' | 'week' | 'month' | 'year'>('week');
  const [costTypeFilter, setCostTypeFilter] = useState<'all' | 'llm' | 'api' | 'storage' | 'compute'>('all');
  const [agentFilter, setAgentFilter] = useState<string>('all');
  const [showDetails, setShowDetails] = useState(false);
  const [currency, setCurrency] = useState<'USD' | 'EUR' | 'GBP' | 'CNY'>('USD');

  // 获取所有代理
  const agents = useMemo(() => {
    const agentSet = new Set(metrics.map(m => m.agent).filter(Boolean));
    return Array.from(agentSet);
  }, [metrics]);

  // 获取所有成本类型
  const costTypes = useMemo(() => {
    const typeSet = new Set(metrics.map(m => m.costType));
    return Array.from(typeSet);
  }, [metrics]);

  // 过滤指标
  const filteredMetrics = useMemo(() => {
    let filtered = metrics.filter(metric => {
      // 时间范围过滤（这里简化处理，实际应该根据时间戳过滤）
      // 成本类型过滤
      if (costTypeFilter !== 'all' && metric.costType !== costTypeFilter) {
        return false;
      }
      
      // 代理过滤
      if (agentFilter !== 'all' && metric.agent !== agentFilter) {
        return false;
      }
      
      return true;
    });
    
    return filtered;
  }, [metrics, costTypeFilter, agentFilter]);

  // 计算统计数据
  const stats = useMemo(() => {
    const totalCost = filteredMetrics.reduce((sum, metric) => sum + metric.amount, 0);
    const avgCost = filteredMetrics.length > 0 ? totalCost / filteredMetrics.length : 0;
    const maxCost = Math.max(...filteredMetrics.map(m => m.amount), 0);
    const minCost = Math.min(...filteredMetrics.map(m => m.amount), Infinity) || 0;
    
    // 按成本类型分组
    const byType: Record<string, { total: number; count: number }> = {};
    filteredMetrics.forEach(metric => {
      if (!byType[metric.costType]) {
        byType[metric.costType] = { total: 0, count: 0 };
      }
      byType[metric.costType].total += metric.amount;
      byType[metric.costType].count += 1;
    });
    
    // 按代理分组
    const byAgent: Record<string, { total: number; count: number }> = {};
    filteredMetrics.forEach(metric => {
      const agent = metric.agent || 'Unknown';
      if (!byAgent[agent]) {
        byAgent[agent] = { total: 0, count: 0 };
      }
      byAgent[agent].total += metric.amount;
      byAgent[agent].count += 1;
    });
    
    return {
      totalCost,
      avgCost,
      maxCost,
      minCost,
      count: filteredMetrics.length,
      byType,
      byAgent
    };
  }, [filteredMetrics]);

  const getCostTypeIcon = (costType: string) => {
    switch (costType) {
      case 'llm': return Brain;
      case 'api': return Server;
      case 'storage': return Database;
      case 'compute': return Cpu;
      default: return DollarSign;
    }
  };

  const getCostTypeColor = (costType: string) => {
    switch (costType) {
      case 'llm': return 'text-purple-500';
      case 'api': return 'text-blue-500';
      case 'storage': return 'text-green-500';
      case 'compute': return 'text-orange-500';
      default: return 'text-white/60';
    }
  };

  const getCostTypeBgColor = (costType: string) => {
    switch (costType) {
      case 'llm': return 'bg-purple-500/10';
      case 'api': return 'bg-blue-500/10';
      case 'storage': return 'bg-green-500/10';
      case 'compute': return 'bg-orange-500/10';
      default: return 'bg-white/5';
    }
  };

  const formatCurrency = (amount: number) => {
    const formatter = new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
      minimumFractionDigits: 2,
      maximumFractionDigits: 4
    });
    return formatter.format(amount);
  };

  const formatTokenCount = (tokens: number) => {
    if (tokens < 1000) return `${tokens}`;
    if (tokens < 1000000) return `${(tokens / 1000).toFixed(1)}K`;
    return `${(tokens / 1000000).toFixed(2)}M`;
  };

  const handleDownload = () => {
    const data = {
      metrics: filteredMetrics,
      stats,
      timestamp: new Date().toISOString(),
      filters: { timeRange, costTypeFilter, agentFilter },
      currency
    };
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `cost-metrics-${new Date().toISOString().split('T')[0]}.json`;
    a.click();
    URL.revokeObjectURL(url);
  };

  // 模拟时间序列数据
  const timeSeriesData = useMemo(() => {
    const days = timeRange === 'today' ? 24 : 
                 timeRange === 'week' ? 7 : 
                 timeRange === 'month' ? 30 : 365;
    
    return Array.from({ length: days }, (_, i) => {
      const date = new Date();
      date.setDate(date.getDate() - (days - 1 - i));
      
      // 模拟成本数据
      const baseCost = 10 + Math.random() * 50;
      const llmCost = baseCost * 0.6;
      const apiCost = baseCost * 0.2;
      const storageCost = baseCost * 0.1;
      const computeCost = baseCost * 0.1;
      
      return {
        date: date.toISOString().split('T')[0],
        total: baseCost,
        llm: llmCost,
        api: apiCost,
        storage: storageCost,
        compute: computeCost
      };
    });
  }, [timeRange]);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-3">
            <DollarSign size={24} />
            Cost Dashboard
          </h2>
          <p className="text-sm text-white/60 mt-1">
            Token consumption tracking and cost estimation across all agents
          </p>
        </div>
        
        <div className="flex items-center gap-3">
          <div className="text-sm text-white/40">
            {formatCurrency(stats.totalCost)} total • {stats.count} transactions
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
              <option value="today">Today</option>
              <option value="week">Last 7 Days</option>
              <option value="month">Last 30 Days</option>
              <option value="year">Last Year</option>
            </select>
            
            <select
              value={costTypeFilter}
              onChange={(e) => setCostTypeFilter(e.target.value as any)}
              className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
            >
              <option value="all">All Cost Types</option>
              {costTypes.map(type => (
                <option key={type} value={type}>{type.toUpperCase()}</option>
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
            
            <select
              value={currency}
              onChange={(e) => setCurrency(e.target.value as any)}
              className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
            >
              <option value="USD">USD ($)</option>
              <option value="EUR">EUR (€)</option>
              <option value="GBP">GBP (£)</option>
              <option value="CNY">CNY (¥)</option>
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
              title="Download cost data"
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
            <div className="text-xs text-white/40 uppercase tracking-wider">Total Cost</div>
            <DollarSign size={20} className="text-primary" />
          </div>
          <div className="text-3xl font-bold">{formatCurrency(stats.totalCost)}</div>
          <div className="flex items-center gap-2 mt-2">
            <TrendingUp size={16} className="text-green-500" />
            <span className="text-sm text-green-500">+12.5% from last period</span>
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-xl p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="text-xs text-white/40 uppercase tracking-wider">Avg Cost per Task</div>
            <BarChart3 size={20} className="text-blue-500" />
          </div>
          <div className="text-3xl font-bold">{formatCurrency(stats.avgCost)}</div>
          <div className="flex items-center gap-2 mt-2">
            <TrendingDown size={16} className="text-red-500" />
            <span className="text-sm text-red-500">-3.2% from last period</span>
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-xl p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="text-xs text-white/40 uppercase tracking-wider">Total Tokens</div>
            <Brain size={20} className="text-purple-500" />
          </div>
          <div className="text-3xl font-bold">
            {formatTokenCount(filteredMetrics.reduce((sum, m) => sum + (m.tokenCount || 0), 0))}
          </div>
          <div className="text-sm text-white/60 mt-2">
            {filteredMetrics.length} transactions
          </div>
        </div>
        
        <div className="bg-white/5 border border-white/10 rounded-xl p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="text-xs text-white/40 uppercase tracking-wider">Cost Efficiency</div>
            <Target size={20} className="text-green-500" />
          </div>
          <div className="text-3xl font-bold">84%</div>
          <div className="text-sm text-white/60 mt-2">
            Higher than industry average
          </div>
        </div>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Cost Over Time */}
        <div className="border border-white/10 rounded-xl p-6">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-lg font-bold">Cost Over Time</h3>
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
                  <div className="flex items-end justify-center w-full" style={{ height: '100%' }}>
                    <div 
                      className="w-full bg-primary/40 rounded-t hover:bg-primary transition-colors"
                      style={{ height: `${(day.total / 100) * 100}%` }}
                      title={`${day.date}: ${formatCurrency(day.total)}`}
                    />
                  </div>
                  <div className="text-xs text-white/40 mt-2">
                    {timeRange === 'today' 
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
          
          <div className="grid grid-cols-4 gap-4 mt-6">
            <div className="text-center">
              <div className="text-xs text-white/40">LLM</div>
              <div className="text-lg font-bold text-purple-500">
                {formatCurrency(timeSeriesData.reduce((sum, day) => sum + day.llm, 0))}
              </div>
            </div>
            <div className="text-center">
              <div className="text-xs text-white/40">API</div>
              <div className="text-lg font-bold text-blue-500">
                {formatCurrency(timeSeriesData.reduce((sum, day) => sum + day.api, 0))}
              </div>
            </div>
            <div className="text-center">
              <div className="text-xs text-white/40">Storage</div>
              <div className="text-lg font-bold text-green-500">
                {formatCurrency(timeSeriesData.reduce((sum, day) => sum + day.storage, 0))}
              </div>
            </div>
            <div className="text-center">
              <div className="text-xs text-white/40">Compute</div>
              <div className="text-lg font-bold text-orange-500">
                {formatCurrency(timeSeriesData.reduce((sum, day) => sum + day.compute, 0))}
              </div>
            </div>
          </div>
        </div>

        {/* Cost by Type */}
        <div className="border border-white/10 rounded-xl p-6">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-lg font-bold">Cost by Type</h3>
            <PieChart size={20} className="text-white/40" />
          </div>
          
          <div className="flex items-center justify-center h-48">
            <div className="relative w-40 h-40">
              {/* Pie Chart */}
              <svg viewBox="0 0 100 100" className="w-full h-full">
                {Object.entries(stats.byType)
                  .filter(([_, data]) => data.total > 0)
                  .reduce((acc, [type, data], index, array) => {
                    const total = stats.totalCost;
                    const percentage = (data.total / total) * 100;
                    const startAngle = acc.currentAngle;
                    const endAngle = startAngle + (percentage * 3.6);
                    
                    const startX = 50 + 40 * Math.cos((startAngle - 90) * Math.PI / 180);
                    const startY = 50 + 40 * Math.sin((startAngle - 90) * Math.PI / 180);
                    const endX = 50 + 40 * Math.cos((endAngle - 90) * Math.PI / 180);
                    const endY = 50 + 40 * Math.sin((endAngle - 90) * Math.PI / 180);
                    
                    const largeArcFlag = percentage > 50 ? 1 : 0;
                    
                    const path = `M 50 50 L ${startX} ${startY} A 40 40 0 ${largeArcFlag} 1 ${endX} ${endY} Z`;
                    
                    let color = '';
                    switch (type) {
                      case 'llm': color = '#A855F7'; break;
                      case 'api': color = '#3B82F6'; break;
                      case 'storage': color = '#10B981'; break;
                      case 'compute': color = '#F59E0B'; break;
                      default: color = '#6B7280';
                    }
                    
                    acc.elements.push(
                      React.createElement('path', {
                        key: type,
                        d: path,
                        fill: color,
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
                  <div className="text-2xl font-bold">{formatCurrency(stats.totalCost)}</div>
                  <div className="text-xs text-white/60">Total</div>
                </div>
              </div>
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-3 mt-6">
            {Object.entries(stats.byType).map(([type, data]) => {
              const Icon = getCostTypeIcon(type);
              const percentage = (data.total / stats.totalCost) * 100;
              
              return (
                <div key={type} className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <div className={cn(
                      "p-2 rounded-lg",
                      getCostTypeBgColor(type)
                    )}>
                      <Icon size={16} className={getCostTypeColor(type)} />
                    </div>
                    <span className="text-sm capitalize">{type}</span>
                  </div>
                  <div className="text-right">
                    <div className="text-sm font-bold">{formatCurrency(data.total)}</div>
                    <div className="text-xs text-white/40">{percentage.toFixed(1)}%</div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>

      {/* Cost by Agent */}
      <div className="border border-white/10 rounded-xl p-6">
        <h3 className="text-lg font-bold mb-6">Cost by Agent</h3>
        
        <div className="space-y-4">
          {Object.entries(stats.byAgent)
            .sort((a, b) => b[1].total - a[1].total)
            .map(([agent, data]) => {
              const percentage = (data.total / stats.totalCost) * 100;
              
              return (
                <div key={agent} className="space-y-2">
                  <div className="flex justify-between items-center">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-white/5 rounded-lg">
                        <Users size={16} className="text-white/60" />
                      </div>
                      <span className="font-bold">{agent}</span>
                    </div>
                    <div className="text-right">
                      <div className="text-lg font-bold">{formatCurrency(data.total)}</div>
                      <div className="text-sm text-white/40">{data.count} transactions</div>
                    </div>
                  </div>
                  
                  <div className="h-2 bg-white/10 rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-primary rounded-full"
                      style={{ width: `${percentage}%` }}
                    />
                  </div>
                  
                  <div className="flex justify-between text-xs text-white/40">
                    <span>{percentage.toFixed(1)}% of total cost</span>
                    <span>Avg: {formatCurrency(data.total / data.count)} per transaction</span>
                  </div>
                </div>
              );
            })}
        </div>
      </div>

      {/* Detailed Metrics Table */}
      {showDetails && (
        <div className="border border-white/10 rounded-xl overflow-hidden bg-white/[0.02]">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-white/10">
                  <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                    Time
                  </th>
                  <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                    Agent
                  </th>
                  <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                    Type
                  </th>
                  <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                    Amount
                  </th>
                  <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                    Tokens
                  </th>
                  <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                    Description
                  </th>
                </tr>
              </thead>
              <tbody>
                {filteredMetrics.slice(0, 10).map((metric) => {
                  const Icon = getCostTypeIcon(metric.costType);
                  
                  return (
                    <tr key={metric.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                      <td className="p-4">
                        <div className="text-sm">
                          {metric.timestamp ? new Date(metric.timestamp).toLocaleTimeString() : 'Unknown'}
                        </div>
                        <div className="text-xs text-white/40">
                          {metric.timestamp ? new Date(metric.timestamp).toLocaleDateString() : ''}
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="text-sm font-bold">{metric.agent || 'Unknown'}</div>
                      </td>
                      <td className="p-4">
                        <div className={cn(
                          "inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-bold",
                          getCostTypeBgColor(metric.costType),
                          getCostTypeColor(metric.costType)
                        )}>
                          <Icon size={12} />
                          {metric.costType.toUpperCase()}
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="text-lg font-bold">{formatCurrency(metric.amount)}</div>
                      </td>
                      <td className="p-4">
                        <div className="text-sm">
                          {metric.tokenCount ? formatTokenCount(metric.tokenCount) : 'N/A'}
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="text-sm text-white/80 max-w-xs truncate">
                          {metric.description || 'No description'}
                        </div>
                      </td>
                    </tr>
                  );
                })}
                
                {filteredMetrics.length === 0 && (
                  <tr>
                    <td colSpan={6} className="p-8 text-center">
                      <DollarSign size={48} className="mx-auto text-white/20 mb-4" />
                      <div className="text-lg font-bold text-white/40">No cost metrics found</div>
                      <p className="text-white/60 mt-2">
                        Try changing your filters or wait for cost data to be collected
                      </p>
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
          
          {filteredMetrics.length > 10 && (
            <div className="p-4 border-t border-white/10 text-center">
              <button className="text-sm text-white/60 hover:text-white transition-colors">
                Show all {filteredMetrics.length} transactions
              </button>
            </div>
          )}
        </div>
      )}

      {/* Cost Optimization Tips */}
      <div className="border border-white/10 rounded-xl p-6">
        <h3 className="text-lg font-bold mb-6 flex items-center gap-3">
          <Target size={24} />
          Cost Optimization Tips
        </h3>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="flex items-center gap-3 mb-3">
              <Brain size={20} className="text-purple-500" />
              <div className="text-sm font-bold">LLM Cost Reduction</div>
            </div>
            <ul className="text-sm text-white/80 space-y-2">
              <li>• Use smaller models for simple tasks</li>
              <li>• Implement response caching</li>
              <li>• Set token limits per request</li>
              <li>• Use streaming for long responses</li>
            </ul>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="flex items-center gap-3 mb-3">
              <Server size={20} className="text-blue-500" />
              <div className="text-sm font-bold">API Cost Management</div>
            </div>
            <ul className="text-sm text-white/80 space-y-2">
              <li>• Implement request batching</li>
              <li>• Use webhooks instead of polling</li>
              <li>• Cache API responses</li>
              <li>• Monitor rate limits</li>
            </ul>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="flex items-center gap-3 mb-3">
              <Database size={20} className="text-green-500" />
              <div className="text-sm font-bold">Storage Optimization</div>
            </div>
            <ul className="text-sm text-white/80 space-y-2">
              <li>• Implement data compression</li>
              <li>• Use tiered storage</li>
              <li>• Clean up old data regularly</li>
              <li>• Monitor storage growth</li>
            </ul>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-4">
            <div className="flex items-center gap-3 mb-3">
              <Cpu size={20} className="text-orange-500" />
              <div className="text-sm font-bold">Compute Efficiency</div>
            </div>
            <ul className="text-sm text-white/80 space-y-2">
              <li>• Use serverless functions</li>
              <li>• Implement auto-scaling</li>
              <li>• Optimize container sizes</li>
              <li>• Monitor CPU utilization</li>
            </ul>
          </div>
        </div>
        
        <div className="mt-6 p-4 bg-primary/10 border border-primary/20 rounded-lg">
          <div className="flex items-center gap-3">
            <AlertCircle size={20} className="text-primary" />
            <div className="text-sm font-bold">Estimated Monthly Savings</div>
          </div>
          <div className="text-2xl font-bold text-primary mt-2">
            {formatCurrency(stats.totalCost * 0.15)}
          </div>
          <p className="text-sm text-white/60 mt-1">
            Potential 15% cost reduction by implementing optimization strategies
          </p>
        </div>
      </div>
    </div>
  );
}