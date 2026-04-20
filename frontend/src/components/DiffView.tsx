import React from 'react';
import { 
  Diff, 
  Terminal, 
  History, 
  Check, 
  X, 
  RotateCcw, 
  Maximize2, 
  MoreHorizontal,
  Cloud,
  Bell,
  Code2,
  Cpu,
  AlertTriangle,
  Zap,
  CheckCircle2,
  Trash2
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';

export default function DiffView() {
  return (
    <div className="flex-1 flex flex-col lg:flex-row p-6 gap-6 overflow-hidden bg-background h-full font-sans">
      {/* Diff Engine View */}
      <section className="flex-[1.5] min-w-0 flex flex-col bg-surface-container-low/40 rounded-3xl shadow-2xl overflow-hidden border border-outline-variant/10 backdrop-blur-3xl">
        {/* Pane Header */}
        <div className="h-16 bg-surface-container-highest/40 flex items-center justify-between px-6 shrink-0 border-b border-outline-variant/10">
          <div className="flex items-center gap-4">
            <div className="p-2 rounded-xl bg-primary/10 text-primary">
              <Diff size={20} />
            </div>
            <div>
              <h3 className="text-sm font-black text-on-surface tracking-widest uppercase">core_logic.py</h3>
              <div className="flex items-center gap-2 mt-0.5">
                <span className="px-2 py-0.5 rounded-lg text-[9px] font-black bg-surface-container-high text-outline uppercase tracking-widest border border-outline-variant/10">MODIFIED</span>
                <span className="text-[10px] font-mono text-outline-variant opacity-70">commit: #af239b</span>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-6 text-xs font-mono font-black">
            <div className="flex items-center gap-2 text-error px-3 py-1 rounded-lg bg-error/5 border border-error/10">
              <span className="w-1.5 h-1.5 rounded-full bg-error" /> -4
            </div>
            <div className="flex items-center gap-2 text-success px-3 py-1 rounded-lg bg-success/5 border border-success/10">
              <span className="w-1.5 h-1.5 rounded-full bg-success" /> +6
            </div>
          </div>
        </div>

        {/* Diff Content (Side-by-Side) */}
        <div className="flex-1 flex overflow-hidden font-mono text-[13px] leading-relaxed">
          {/* Old Code Pane */}
          <div className="flex-1 bg-surface-container-lowest/30 overflow-y-auto border-r border-outline-variant/10 flex flex-col custom-scrollbar">
            <div className="sticky top-0 bg-surface-container-highest/60 backdrop-blur-md px-6 py-2 text-[10px] text-outline font-black uppercase tracking-widest border-b border-outline-variant/10 flex justify-between items-center z-20">
              <span className="flex items-center gap-2 opacity-70"><History size={12} /> 原始版本 (Agent-Base)</span>
            </div>
            <div className="flex-1 p-0 select-text relative">
              <table className="w-full text-left border-collapse">
                <tbody>
                  <CodeLine number={42} content="def analyze_sentiment(text):" />
                  <CodeLine number={43} content='    """分析文本的基础情感"""' />
                  <CodeLine number={44} content="    if not text:" deleted />
                  <CodeLine number={45} content='        return "neutral"' deleted />
                  <CodeLine number={46} content=" " />
                  <CodeLine number={47} content="    result = model.predict(text)" />
                  <CodeLine number={48} content="    return result.label" />
                </tbody>
              </table>
            </div>
          </div>

          {/* New Code Pane */}
          <div className="flex-1 bg-surface-container-lowest/40 overflow-y-auto flex flex-col custom-scrollbar">
            <div className="sticky top-0 bg-primary/10 backdrop-blur-md px-6 py-2 text-[10px] text-primary font-black uppercase tracking-widest border-b border-primary/20 flex justify-between items-center z-20">
              <span className="flex items-center gap-2 italic"><Zap size={12} /> 认知迭代 (Gen-Pulse-v2)</span>
              <div className="w-2 h-2 bg-primary rounded-full animate-pulse shadow-[0_0_8px_rgba(192,193,255,0.8)]" />
            </div>
            <div className="flex-1 p-0 select-text relative">
              <table className="w-full text-left border-collapse">
                <tbody>
                  <CodeLine number={42} content="def analyze_sentiment(text, deep_scan=False):" />
                  <CodeLine number={43} content='    """分析文本情感，支持深度扫描模式"""' />
                  <CodeLine number={44} content="    if not text.strip():" added />
                  <CodeLine number={45} content='        raise ValueError("输入内容不可为空")' added />
                  <CodeLine number={46} content=" " added />
                  <CodeLine number={47} content='    params = {"mode": "deep"} if deep_scan else {}' added />
                  <CodeLine number={48} content="    result = model.predict(text, **params)" />
                  <CodeLine number={49} content="    return result.label" />
                </tbody>
              </table>
              <div className="absolute top-0 left-0 w-80 h-80 bg-primary/5 rounded-full blur-[100px] pointer-events-none" />
            </div>
          </div>
        </div>

        {/* Footer Actions */}
        <div className="h-20 bg-surface-container-highest/60 shrink-0 flex items-center justify-between px-8 border-t border-outline-variant/10 backdrop-blur-md">
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-xl bg-primary/10 text-primary">
              <Cpu size={20} />
            </div>
            <div>
              <span className="text-sm font-bold text-on-surface">AI 审计建议</span>
              <p className="text-[10px] text-outline font-medium tracking-tight mt-0.5">建议提高错误处理的严谨性，增加边界条件探测。</p>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <button className="px-6 py-2.5 rounded-xl text-sm font-bold text-outline-variant hover:text-white hover:bg-white/5 transition-all uppercase tracking-widest">
              拒绝
            </button>
            <button className="px-8 py-2.5 rounded-xl text-sm font-black bg-gradient-to-br from-primary-container to-inverse-primary text-on-primary-container hover:brightness-110 shadow-xl shadow-primary/20 transition-all flex items-center gap-3 uppercase tracking-widest group">
              <Check size={18} className="group-active:scale-125 transition-transform" />
              接受变更
            </button>
          </div>
        </div>
      </section>

      {/* Terminal View */}
      <section className="flex-1 flex flex-col bg-surface-container-low/60 rounded-3xl shadow-2xl overflow-hidden border border-outline-variant/10 backdrop-blur-2xl">
        {/* Term Header */}
        <div className="h-16 bg-surface-container-highest/40 flex items-center justify-between px-6 shrink-0 border-b border-outline-variant/10">
          <div className="flex items-center gap-4">
            <div className="p-2 rounded-xl bg-surface-container-highest text-outline">
              <Terminal size={18} />
            </div>
            <h3 className="text-sm font-black text-on-surface uppercase tracking-widest">系统控制台</h3>
          </div>
          <div className="flex items-center gap-2">
            <button className="text-outline hover:text-on-surface transition-all p-2 rounded-xl hover:bg-white/5">
              <Trash2 size={16} />
            </button>
            <button className="text-outline hover:text-on-surface transition-all p-2 rounded-xl hover:bg-white/5">
              <Maximize2 size={16} />
            </button>
          </div>
        </div>

        {/* Term Body */}
        <div className="flex-1 bg-surface-container-lowest/20 p-6 overflow-y-auto font-mono text-[11px] leading-loose relative custom-scrollbar">
          <div className="text-outline-variant mb-6 opacity-60 italic">
            Genpulse AI Terminal v2.4.0 [Type 'help' for commands]<br/>
            Connected to cognitive kernel #kernel-0x12a...
          </div>
          
          <div className="space-y-4">
            <div>
              <div className="flex gap-3 text-primary font-black mb-1">
                <span className="shrink-0 opacity-40">~ user@genpulse:</span>
                <span className="text-on-surface selection:bg-primary/30">agent diff --target core_logic.py</span>
              </div>
              <div className="pl-6 space-y-1 opacity-80">
                <div className="flex items-center gap-3 text-outline">
                  <RotateCcw size={12} className="animate-spin text-primary" />
                  正在分析上下文结构...
                </div>
                <div className="text-success flex items-center gap-3">
                  <CheckCircle2 size={12} />
                  ✓ 依赖树解析完成 (120ms)
                </div>
                <div className="text-outline flex items-center gap-3">
                  <ChevronRight size={12} />
                  发现潜在的边界漏洞 (line: 44)
                </div>
                <div className="text-outline flex items-center gap-3">
                  <ChevronRight size={12} />
                  生成重构方案 [deep_scan 支持]
                </div>
                <div className="text-success flex items-center gap-3 font-bold">
                  <CheckCircle2 size={12} />
                  ✓ 方案生成完毕。等待代码审查。
                </div>
              </div>
            </div>

            <div className="flex gap-3 items-center p-3 rounded-xl bg-error/5 border border-error/10">
              <AlertTriangle size={14} className="text-error" />
              <span className="text-error/80 font-bold uppercase tracking-tighter">WARN:</span>
              <span className="text-on-surface-variant font-medium">发现未捕获异常风险，建议增加 strict_mode 参数验证。</span>
            </div>

            {/* Active Prompt */}
            <div className="flex gap-3 items-center pt-2">
              <span className="text-primary-container font-black opacity-40">~ sys@agent:</span>
              <motion.div 
                animate={{ opacity: [0, 1] }}
                transition={{ repeat: Infinity, duration: 0.5 }}
                className="w-2 h-4 bg-primary rounded-sm shadow-[0_0_8px_rgba(192,193,255,0.5)]"
              />
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}

function CodeLine({ number, content, added, deleted }: any) {
  return (
    <tr className={cn(
      "group transition-all duration-200",
      added ? "bg-success/5 hover:bg-success/10" : deleted ? "bg-error/5 hover:bg-error/10" : "hover:bg-white/5"
    )}>
      <td className={cn(
        "w-12 text-right pr-6 select-none py-1 border-r border-outline-variant/5 text-[10px] font-black opacity-30 group-hover:opacity-100 transition-opacity",
        added ? "text-success border-success/20 opacity-100" : deleted ? "text-error border-error/20 opacity-100" : "text-outline"
      )}>
        {number}
      </td>
      <td className={cn(
        "py-1 whitespace-pre px-6 relative transition-all overflow-hidden",
        added ? "text-success font-bold" : deleted ? "text-error line-through opacity-60" : "text-on-surface-variant"
      )}>
        {added && <div className="absolute left-0 top-0 bottom-0 w-0.5 bg-success shadow-[0_0_10px_rgba(74,222,128,0.5)]" />}
        {deleted && <div className="absolute left-0 top-0 bottom-0 w-0.5 bg-error shadow-[0_0_10px_rgba(255,180,171,0.5)]" />}
        {content}
      </td>
    </tr>
  );
}

function ChevronRight({ size, className }: any) {
  return (
    <svg 
      width={size} 
      height={size} 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      strokeWidth="3" 
      strokeLinecap="round" 
      strokeLinejoin="round" 
      className={className}
    >
      <polyline points="9 18 15 12 9 6" />
    </svg>
  );
}
