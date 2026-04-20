import React from 'react';
import { 
  Search, 
  Filter, 
  Brain, 
  History, 
  Settings, 
  FileText, 
  HelpCircle, 
  Beaker,
  LucideIcon,
  MessageSquare,
  Clock,
  MoreVertical,
  Save,
  User,
  FolderLock,
  Code2,
  Gavel,
  CheckCircle2,
  Diff,
  Terminal,
  Zap,
  Box,
  ChevronRight,
  ShieldCheck,
  Sparkles as AutoAwesome,
  Maximize2
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';

export default function NeuralAssetsView() {
  return (
    <div className="flex flex-1 h-full overflow-hidden bg-[#0A0A0A]">
      {/* Sidebar: Memories List */}
      <aside className="w-96 border-r border-white/10 flex flex-col pt-12">
        <div className="px-10 mb-10 text-primary">
          <span className="text-[10px] uppercase font-black tracking-[0.4em] text-white/40 block mb-4">Memory Retrieval / Phase 02</span>
          <h2 className="text-4xl font-black uppercase tracking-tighter">Episodic<br/>Repository</h2>
        </div>
        
        <div className="px-10 mb-8 flex gap-2">
          <div className="flex-1 relative group">
            <Search className="absolute left-0 top-1/2 -translate-y-1/2 text-white/20 group-focus-within:text-primary transition-colors" size={14} />
            <input 
              type="text" 
              placeholder="SEARCH MEMORIES..."
              className="w-full bg-transparent border-b border-white/20 py-3 pl-6 text-[10px] font-black uppercase tracking-widest text-white outline-none focus:border-primary transition-all placeholder:text-white/10"
            />
          </div>
          <button className="border border-white/10 p-3 hover:bg-white/5 transition-all">
            <Filter size={14} className="text-white/40" />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto px-6 space-y-1 custom-scrollbar">
          {[...Array(6)].map((_, i) => (
            <MemoryCard 
              key={i} 
              active={i === 0}
              timestamp="2024.04.12 14:32"
              title={i === 0 ? "Refactor Logic" : "API Integration"}
              tags={["Core", "Logic"]}
            />
          ))}
        </div>
      </aside>

      {/* Main Content: Memory Editor */}
      <main className="flex-1 flex flex-col overflow-hidden border-r border-white/10">
        <div className="h-20 flex bg-white/[0.02] border-b border-white/10 shrink-0">
          <button className="px-12 h-full text-[10px] font-black uppercase tracking-[0.3em] text-black bg-primary border-r border-white/10">
            Memory.md
          </button>
          <button className="px-12 h-full text-[10px] font-black uppercase tracking-[0.3em] text-white/40 hover:text-white transition-all border-r border-white/10">
            User.md
          </button>
          <div className="ml-auto flex items-center px-10 text-[10px] font-black uppercase tracking-widest text-white/20 italic">
            Retrieval Mode: Heuristic
          </div>
        </div>

        <div className="flex-1 overflow-y-auto p-12 space-y-12 custom-scrollbar bg-[#050505]">
          <section className="border border-white/10 p-10 relative group">
            <div className="absolute top-0 right-0 w-32 h-32 bg-primary/5 rounded-full blur-3xl -z-10" />
            <div className="text-[10px] font-black uppercase tracking-[0.4em] text-primary mb-6 flex items-center gap-4">
              <span className="w-8 h-[1px] bg-primary"></span>
              Strategic Alignment
            </div>
            <p className="text-lg font-black uppercase tracking-tight text-white/80 leading-relaxed max-w-3xl">
              Maintain rigorous adherence to functional purity in the core module. All side-effects must be encapsulated within the persistence boundary to ensure deterministic agent execution.
            </p>
          </section>

          <section className="space-y-6">
             <div className="text-[10px] font-black uppercase tracking-[0.4em] text-white/40 mb-10 flex items-center gap-4">
              <span className="w-8 h-[1px] bg-white/10"></span>
              Logic Blueprint / V4.2
            </div>
            <div className="bg-black border border-white/10">
              <div className="bg-white/5 px-6 py-3 flex items-center justify-between border-b border-white/10">
                <div className="flex items-center gap-4">
                  <Code2 size={12} className="text-primary" />
                  <span className="text-[9px] font-black uppercase tracking-widest text-white/40">Core_Module_V2.py</span>
                </div>
                <div className="flex gap-4">
                  <button className="text-white/20 hover:text-primary transition-colors"><Maximize2 size={12} /></button>
                  <button className="text-white/20 hover:text-primary transition-colors"><Save size={12} /></button>
                </div>
              </div>
              <div className="p-8 font-mono text-[11px] text-white/60 leading-loose">
                <pre><code>{`class NeuralOptimizer:
    def __init__(self, weights: dict):
        self.weights = weights
        self.momentum = 0.9

    def process_node(self, node_id: str):
        # Heuristic retrieval bypass
        if node_id in self.weights:
            return self.weights[node_id] * self.momentum
        return 1.0`}</code></pre>
              </div>
            </div>
          </section>

          <section>
            <div className="text-[10px] font-black uppercase tracking-[0.4em] text-white/40 mb-10 flex items-center gap-4">
              <span className="w-8 h-[1px] bg-white/10"></span>
              Operator Constraints
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-1 border-t border-white/10 pt-10">
              <div className="flex flex-col gap-4">
                <span className="text-4xl font-black text-primary">01</span>
                <span className="text-[10px] uppercase font-black tracking-widest">Latency Limit</span>
                <p className="text-xs text-white/40 leading-relaxed uppercase font-bold tracking-wide">All retrieval requests must terminate within 200ms.</p>
              </div>
              <div className="flex flex-col gap-4">
                <span className="text-4xl font-black text-primary">02</span>
                <span className="text-[10px] uppercase font-black tracking-widest">Consistency Gate</span>
                <p className="text-xs text-white/40 leading-relaxed uppercase font-bold tracking-wide">Block updates that fail heuristic validation protocols.</p>
              </div>
            </div>
          </section>
        </div>
      </main>

      {/* Right Content: Interaction Timeline */}
      <aside className="w-[450px] shrink-0 flex flex-col bg-[#0A0A0A] pt-12">
        <div className="px-12 mb-10">
          <span className="text-[10px] uppercase font-black tracking-[0.4em] text-white/40 block mb-4">Live Interaction / S08</span>
          <h2 className="text-4xl font-black uppercase tracking-tighter">Interaction<br/>Pathways</h2>
        </div>

        <div className="flex-1 overflow-y-auto px-12 space-y-12 custom-scrollbar border-t border-white/10 pt-10">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="relative group">
              <div className="absolute left-[-49px] top-4 w-4 h-4 bg-primary rotate-45 border-4 border-[#0A0A0A] z-10" />
              <div className="mb-6 flex justify-between items-end">
                <div className="text-[9px] font-black uppercase tracking-widest text-primary">Analysis Node / 14:55</div>
                <div className="text-stroke text-2xl font-black text-white/10">0{i+1}</div>
              </div>
              <div className="bg-white/5 p-8 border border-white/10 group-hover:bg-white/[0.08] transition-all">
                <p className="text-xs font-bold uppercase tracking-wide text-white/60 mb-6 leading-relaxed">
                  Proposed structural changes to the internal weighting system observed. Recommendation: Accept changes to optimize retrieval speed.
                </p>
                <div className="flex gap-4">
                  <button className="flex-1 py-3 text-[9px] font-black uppercase tracking-widest bg-primary text-black hover:scale-105 transition-all">
                    Approve
                  </button>
                  <button className="flex-1 py-3 text-[9px] font-black uppercase tracking-widest border border-white/20 text-white/40 hover:text-white transition-all">
                    Reject
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </aside>
    </div>
  );
}

function FilterChip({ icon: Icon, label }: { icon: LucideIcon, label: string }) {
  return (
    <span className="bg-surface-container-high/50 text-on-surface text-[10px] font-bold uppercase px-3 py-1.5 rounded-full border border-outline-variant/10 flex items-center gap-2 cursor-pointer hover:bg-surface-variant hover:text-primary transition-all active:scale-95 shadow-sm">
      <Icon size={12} />
      {label}
    </span>
  );
}

function MemoryCard({ id, title, time, description, tags, color, active }: any) {
  return (
    <div className={cn(
      "rounded-2xl p-5 flex flex-col gap-3 cursor-pointer transition-all duration-300 group border relative overflow-hidden",
      active 
        ? "bg-surface-container-low border-primary shadow-xl scale-[1.02]" 
        : "bg-surface-container-low/40 border-outline-variant/5 hover:border-outline-variant/20 hover:scale-[1.01]"
    )}>
      {active && <div className="absolute top-0 left-0 w-1 h-full bg-primary" />}
      <div className="flex justify-between items-start">
        <div className="flex items-center gap-2">
          <div className={cn("w-2 h-2 rounded-full shadow-lg", color)} />
          <span className="text-[10px] font-mono font-bold text-outline uppercase tracking-widest transition-colors group-hover:text-primary">ID: {id}</span>
        </div>
        <span className="text-[10px] text-outline font-bold uppercase tracking-tighter opacity-60">{time}</span>
      </div>
      <h3 className={cn("text-sm font-bold transition-colors", active ? "text-primary" : "text-on-surface group-hover:text-primary")}>
        {title}
      </h3>
      <p className="text-xs text-on-surface-variant line-clamp-2 leading-relaxed opacity-80">
        {description}
      </p>
      <div className="flex gap-2 mt-1">
        {tags.map((tag: string) => (
          <span key={tag} className="text-[9px] font-bold py-1 px-2 rounded-lg bg-surface-container-highest text-outline uppercase tracking-widest border border-outline-variant/5">
            {tag}
          </span>
        ))}
      </div>
    </div>
  );
}

function DecisionCard({ title, content }: any) {
  return (
    <div className="flex items-start gap-4 p-4 rounded-2xl bg-surface-container-low/80 border border-outline-variant/10 hover:border-primary/20 transition-all group">
      <div className="mt-1.5 w-2 h-2 rounded-full bg-primary-container shrink-0 shadow-[0_0_8px_rgba(91,95,255,0.4)] group-hover:bg-primary transition-all" />
      <div>
        <span className="font-bold text-on-surface text-sm block mb-1 group-hover:text-primary transition-colors">{title}</span>
        <span className="text-on-surface-variant text-xs leading-relaxed opacity-80">{content}</span>
      </div>
    </div>
  );
}

function TimelineItem({ type, title, time, content, isDiff, active }: any) {
  const Icon = type === 'ai' ? Brain : type === 'diff' ? Diff : User;
  
  return (
    <div className="relative flex gap-5 group">
      <div className={cn(
        "relative z-10 w-9 h-9 rounded-full flex items-center justify-center shrink-0 mt-1 transition-all duration-300",
        active 
          ? "bg-surface-container-highest border border-primary/40 text-primary shadow-xl shadow-primary/10" 
          : "bg-surface-container-highest/60 border border-outline-variant/20 text-outline group-hover:border-primary/20"
      )}>
        <Icon size={18} className={cn(active && "animate-pulse")} />
        {active && <span className="absolute -inset-1 rounded-full border border-primary/20 animate-ping" />}
      </div>
      
      <div className={cn(
        "flex-1 rounded-2xl p-4 border transition-all duration-300 group-hover:translate-x-1 shadow-sm h-min",
        active 
          ? "bg-surface-container-highest/60 border-primary/20 backdrop-blur-xl" 
          : "bg-surface-container-low border-outline-variant/5 group-hover:border-outline-variant/20"
      )}>
        <div className="flex justify-between items-center mb-2">
          <span className={cn("text-xs font-bold uppercase tracking-widest", active ? "text-primary" : "text-on-surface")}>{title}</span>
          <span className="text-[10px] text-outline font-mono font-bold">{time}</span>
        </div>
        
        {isDiff ? (
          <div className="rounded-xl overflow-hidden text-[10px] font-mono border border-outline-variant/10 shadow-inner">
            <div className="bg-error/10 text-error px-3 py-1.5 flex gap-3 border-b border-error/5">
              <span className="opacity-50 font-black">-</span>
              <span className="line-through opacity-80">validateToken(req.body)</span>
            </div>
            <div className="bg-success/10 text-success px-3 py-1.5 flex gap-3">
              <span className="opacity-50 font-black">+</span>
              <span className="font-bold">await jwt.verify(header, RS256)</span>
            </div>
          </div>
        ) : (
          <p className="text-xs text-on-surface-variant leading-relaxed italic opacity-80">
            {content}
          </p>
        )}
      </div>
    </div>
  );
}
