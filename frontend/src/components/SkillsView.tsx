import React from 'react';
import { 
  Box, 
  Search, 
  Filter, 
  MoreVertical, 
  History, 
  Copy, 
  Edit, 
  Code2, 
  Database, 
  GitBranch, 
  Terminal,
  Cpu,
  Fingerprint,
  Layers,
  Puzzle,
  Box as BoxIcon,
  ChevronRight,
  Plus
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';

export default function SkillsView() {
  return (
    <div className="flex-1 overflow-hidden flex flex-col p-8 gap-8 h-full bg-background font-sans">
      {/* Workspace Header */}
      <header className="flex items-end justify-between shrink-0">
        <div>
          <h2 className="text-5xl font-black text-on-surface tracking-tight leading-tight mb-2">技能库</h2>
          <p className="text-outline text-base font-medium">管理和配置代理的认知能力模块。</p>
        </div>
        <div className="flex gap-3">
          <button className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-surface-container hover:bg-surface-container-highest text-on-surface-variant text-xs font-bold border border-outline-variant/10 shadow-sm transition-all active:scale-95 uppercase tracking-widest">
            <Filter size={16} />
            过滤
          </button>
          <button className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-surface-container hover:bg-surface-container-highest text-on-surface-variant text-xs font-bold border border-outline-variant/10 shadow-sm transition-all active:scale-95 uppercase tracking-widest">
            <Layers size={16} />
            排序
          </button>
        </div>
      </header>

      {/* Bento Grid Layout */}
      <div className="flex-1 flex gap-8 overflow-hidden min-h-0">
        {/* Left Pane: Skill List */}
        <div className="w-[420px] shrink-0 flex flex-col gap-5 overflow-y-auto pr-3 pb-8 custom-scrollbar">
          {/* Active Skill Card */}
          <div className="bg-surface-container-highest/60 rounded-3xl p-6 flex flex-col gap-6 relative cursor-pointer border border-primary/20 shadow-2xl group transition-all hover:bg-surface-container-highest">
            {/* Selection Indicator */}
            <div className="absolute left-0 top-1/2 -translate-y-1/2 w-1.5 h-12 bg-primary rounded-r-full shadow-[0_0_15px_rgba(91,95,255,0.8)]" />
            
            <div className="flex justify-between items-start">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-2xl bg-primary-container/20 flex items-center justify-center text-primary shadow-inner border border-primary/20">
                  <BoxIcon size={24} />
                </div>
                <div>
                  <h3 className="text-lg font-black text-on-surface flex items-center gap-3">
                    React 专家
                    <span className="relative flex h-2.5 w-2.5">
                      <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-secondary opacity-50" />
                      <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-secondary" />
                    </span>
                  </h3>
                  <span className="text-[10px] text-outline font-mono font-bold uppercase tracking-widest opacity-70">v1.4.2 • 认知模块</span>
                </div>
              </div>
              <button className="text-outline hover:text-primary transition-colors p-1.5 rounded-xl hover:bg-white/5">
                <MoreVertical size={20} />
              </button>
            </div>
            
            <p className="text-sm text-on-surface-variant leading-relaxed opacity-80">
              精通现代 React 架构，能够生成高性能组件、处理复杂状态管理并优化渲染流水线。
            </p>
            
            <div className="flex items-center gap-6 pt-5 border-t border-outline-variant/10">
              <SkillStat label="成功率" value="98.4%" color="text-secondary" />
              <div className="w-px h-8 bg-outline-variant/10" />
              <SkillStat label="调用次数" value="12,450" />
              <div className="w-px h-8 bg-outline-variant/10" />
              <SkillStat label="延迟" value="240ms" />
            </div>
          </div>

          {/* Inactive Skill Cards */}
          <div className="space-y-4">
            <InactiveSkillCard 
              icon={Terminal} 
              title="Go 后端" 
              version="v2.1.0" 
              type="逻辑处理" 
              description="专注于高并发微服务架构，提供稳健的 API 设计、数据库分片策略与协程管理方案。"
              stat="95.2%"
              calls="8,204"
            />
            <InactiveSkillCard 
              icon={GitBranch} 
              title="Git 流水线" 
              version="v1.0.8" 
              type="运维自动化" 
              description="自动化 CI/CD 流程构建，代码合并冲突智能解决，部署策略优化。"
              stat="99.1%"
              calls="45,112"
            />
          </div>
        </div>

        {/* Right Pane: Preview Editor */}
        <div className="flex-1 flex flex-col bg-surface-container-lowest/40 rounded-3xl border border-outline-variant/10 overflow-hidden shadow-2xl backdrop-blur-xl">
          {/* Editor Header / Tabs */}
          <div className="bg-surface-container-highest/50 border-b border-outline-variant/10 flex items-center px-4 py-2 gap-2 shrink-0">
            <div className="flex items-center gap-2.5 px-5 py-2.5 bg-surface-container-low rounded-xl text-primary text-xs font-bold tracking-widest border-t border-l border-r border-outline-variant/10 shadow-lg">
              <FileText size={14} className="opacity-70" />
              definition.yaml
            </div>
            <div className="flex items-center gap-2.5 px-5 py-2.5 text-outline-variant hover:text-on-surface hover:bg-white/5 rounded-xl text-xs font-bold tracking-widest transition-all cursor-pointer">
              <History size={14} className="opacity-50" />
              README.md
            </div>
            <div className="flex-1" />
            <div className="flex items-center gap-2 px-2">
              <button className="text-outline hover:text-primary transition-all p-2 rounded-xl hover:bg-white/5" title="Copy to clipboard">
                <Copy size={16} />
              </button>
              <button className="text-outline hover:text-primary transition-all p-2 rounded-xl hover:bg-white/5" title="Edit definition">
                <Edit size={16} />
              </button>
            </div>
          </div>

          {/* Editor Body (Code View) */}
          <div className="flex-1 p-8 overflow-y-auto bg-surface-container-lowest/30 font-mono text-sm leading-relaxed custom-scrollbar relative">
            <div className="absolute top-0 right-0 w-80 h-80 bg-primary/5 rounded-full blur-[100px] pointer-events-none" />
            <pre className="m-0 relative z-10 select-text"><code className="text-on-surface-variant/90 font-medium">
<span className="text-primary font-bold">name:</span> <span className="text-tertiary">react-expert</span><br/>
<span className="text-primary font-bold">type:</span> <span className="text-tertiary">cognitive-skill</span><br/>
<span className="text-primary font-bold">version:</span> <span className="text-tertiary">"1.4.2"</span><br/>
<span className="text-primary font-bold">description:</span> <span className="text-tertiary">"高级 React 组件生成与优化专家"</span><br/>
<br/>
<span className="text-primary font-bold">parameters:</span><br/>
&nbsp;&nbsp;<span className="text-secondary font-bold">framework_version:</span> <span className="text-tertiary">"&gt;=18.0.0"</span><br/>
&nbsp;&nbsp;<span className="text-secondary font-bold">strict_mode:</span> <span className="text-primary">true</span><br/>
&nbsp;&nbsp;<span className="text-secondary font-bold">styling_preference:</span><br/>
&nbsp;&nbsp;&nbsp;&nbsp;<span className="text-outline">-</span> <span className="text-tertiary">"tailwind"</span><br/>
&nbsp;&nbsp;&nbsp;&nbsp;<span className="text-outline">-</span> <span className="text-tertiary">"css-modules"</span><br/>
<br/>
<span className="text-primary font-bold">capabilities:</span><br/>
&nbsp;&nbsp;<span className="text-outline">-</span> <span className="text-secondary font-bold">id:</span> <span className="text-tertiary">ui-generation</span><br/>
&nbsp;&nbsp;&nbsp;&nbsp;<span className="text-secondary font-bold">confidence_score:</span> <span className="text-secondary">0.98</span><br/>
&nbsp;&nbsp;&nbsp;&nbsp;<span className="text-secondary font-bold">dependencies:</span> <span className="text-outline">[]</span><br/>
<br/>
&nbsp;&nbsp;<span className="text-outline">-</span> <span className="text-secondary font-bold">id:</span> <span className="text-tertiary">state-management</span><br/>
&nbsp;&nbsp;&nbsp;&nbsp;<span className="text-secondary font-bold">supported_libs:</span><br/>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span className="text-outline">-</span> <span className="text-tertiary">"zustand"</span><br/>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span className="text-outline">-</span> <span className="text-tertiary">"redux-toolkit"</span><br/>
<br/>
<span className="text-primary font-bold">execution_context:</span><br/>
&nbsp;&nbsp;<span className="text-secondary font-bold">max_tokens:</span> <span className="text-secondary">4096</span><br/>
&nbsp;&nbsp;<span className="text-secondary font-bold">temperature:</span> <span className="text-secondary">0.2</span><br/>
&nbsp;&nbsp;<span className="text-secondary font-bold">system_prompt:</span> <span className="text-tertiary">|</span><br/>
&nbsp;&nbsp;&nbsp;&nbsp;<span className="text-tertiary italic opacity-70">You are an elite React engineer. Focus on performance,<br/>
&nbsp;&nbsp;&nbsp;&nbsp;accessibility, and modern hooks patterns. Do not use<br/>
&nbsp;&nbsp;&nbsp;&nbsp;class components unless explicitly requested.</span><br/>
<br/>
<span className="text-primary font-bold">metrics_target:</span><br/>
&nbsp;&nbsp;<span className="text-secondary font-bold">min_success_rate:</span> <span className="text-secondary">0.95</span><br/>
&nbsp;&nbsp;<span className="text-secondary font-bold">max_latency_ms:</span> <span className="text-secondary">500</span>
            </code></pre>
          </div>
        </div>
      </div>
    </div>
  );
}

function SkillStat({ label, value, color }: any) {
  return (
    <div className="flex flex-col gap-1">
      <span className="text-[10px] uppercase font-black tracking-widest text-outline-variant leading-none">{label}</span>
      <span className={cn("text-base font-mono font-bold leading-none", color || "text-on-surface")}>{value}</span>
    </div>
  );
}

function InactiveSkillCard({ icon: Icon, title, version, type, description, stat, calls }: any) {
  return (
    <div className="bg-surface-container-low/40 rounded-3xl p-6 flex flex-col gap-5 border border-outline-variant/5 hover:border-outline-variant/20 hover:bg-surface-container-high transition-all group">
      <div className="flex justify-between items-start">
        <div className="flex items-center gap-4">
          <div className="w-12 h-12 rounded-2xl bg-surface-container-highest/50 flex items-center justify-center text-outline group-hover:text-primary group-hover:bg-primary/10 transition-all border border-outline-variant/10 group-hover:border-primary/20">
            <Icon size={24} />
          </div>
          <div>
            <h4 className="text-base font-bold text-on-surface-variant group-hover:text-on-surface transition-all">{title}</h4>
            <span className="text-[10px] text-outline font-mono font-bold uppercase tracking-widest opacity-60">{version} • {type}</span>
          </div>
        </div>
        <button className="text-outline-variant hover:text-primary transition-all p-1.5 opacity-0 group-hover:opacity-100">
          <ChevronRight size={18} />
        </button>
      </div>
      <p className="text-sm text-outline leading-relaxed line-clamp-2 opacity-80 group-hover:opacity-100 transition-opacity">
        {description}
      </p>
      <div className="flex items-center gap-6 pt-5 border-t border-outline-variant/10">
        <div className="flex flex-col gap-1">
          <span className="text-[10px] font-mono font-bold text-secondary">{stat}</span>
        </div>
        <div className="w-px h-4 bg-outline-variant/10" />
        <div className="flex flex-col gap-1">
          <span className="text-[10px] font-mono font-bold text-outline opacity-70">{calls} calls</span>
        </div>
      </div>
    </div>
  );
}

function FileText({ size, className }: any) {
  return (
    <svg 
      width={size} 
      height={size} 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      strokeWidth="2.5" 
      strokeLinecap="round" 
      strokeLinejoin="round" 
      className={className}
    >
      <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z" />
      <polyline points="14 2 14 8 20 8" />
      <line x1="16" x2="8" y1="13" y2="13" />
      <line x1="16" x2="8" y1="17" y2="17" />
      <line x1="10" x2="8" y1="9" y2="9" />
    </svg>
  );
}
