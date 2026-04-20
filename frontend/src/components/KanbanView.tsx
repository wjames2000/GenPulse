import React from 'react';
import { 
  GitBranch, 
  Ruler, 
  Code2, 
  CheckSquare as FactCheck, 
  MoreHorizontal, 
  CheckCircle2, 
  History, 
  Terminal, 
  Settings,
  CheckSquare as FactCheckIcon,
  ShieldCheck,
  Layout,
  Brain,
  MessageSquare
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';

export default function KanbanView() {
  return (
    <div className="flex-1 overflow-x-auto p-12 bg-[#0A0A0A] h-full custom-scrollbar">
      <div className="flex gap-1 min-w-max h-full pb-8">
        <KanbanColumn 
          index="01"
          title="Orchestrator" 
          count={2} 
          icon={Brain} 
          color="text-primary"
        >
          <KanbanCard 
            tag="Analyzing" 
            tagColor="text-primary border-primary"
            title="Logic Parsing"
            content="Deconstructing user intent into discrete functional blocks for architecture mapping..."
            ai
            tasks={12}
          />
          <KanbanCard 
            tag="Pending" 
            tagColor="text-white/20 border-white/10"
            title="Service Mesh"
            content="Defining isolated node communication protocols for the core repository."
          />
        </KanbanColumn>

        <KanbanColumn 
          index="02"
          title="Architect" 
          count={1} 
          icon={Ruler} 
          color="text-white"
        >
          <div className="bg-white/5 border border-white/10 p-10 flex flex-col gap-8 group relative overflow-hidden">
            <div className="flex justify-between items-start">
               <span className="text-[10px] font-black uppercase tracking-widest text-primary animate-pulse">Processing...</span>
               <Terminal size={14} className="text-white/20" />
            </div>
            
            <div>
              <h3 className="text-2xl font-black uppercase tracking-tighter text-white mb-6">Database Schema</h3>
              <div className="bg-black border border-white/10 p-6">
                <code className="text-xs font-mono text-primary leading-relaxed uppercase">
                  CREATE TABLE users (<br/>
                  &nbsp;&nbsp;id UUID PRIMARY KEY,<br/>
                  &nbsp;&nbsp;email VARCHAR(255) UNIQUE<br/>
                  );
                </code>
              </div>
            </div>
            <div className="absolute top-0 right-0 w-32 h-32 bg-primary/5 rounded-full blur-3xl -z-10" />
          </div>
        </KanbanColumn>

        <KanbanColumn 
          index="03"
          title="Developer" 
          count={3} 
          icon={Code2} 
          color="text-white/40"
        >
          <KanbanCard 
            tag="Completed" 
            tagColor="text-white/10 border-white/5"
            title="Express Routes"
            content="Standardized API surface area defined and implemented."
            completed
          />
          <KanbanCard 
            tag="Completed" 
            tagColor="text-white/10 border-white/5"
            title="JWT Middleware"
            content="Strict token validation layer integration finalized."
            completed
          />
        </KanbanColumn>

        <KanbanColumn 
          index="04"
          title="QA Lead" 
          count={1} 
          icon={FactCheckIcon} 
          color="text-error"
        >
          <div className="bg-white/5 border border-error/20 p-10 flex flex-col gap-8 group">
            <div className="flex justify-between items-start">
              <span className="text-[10px] font-black uppercase tracking-widest text-error">Failure Detected</span>
              <AlertCircle size={18} className="text-error" />
            </div>
            
            <div>
              <h3 className="text-2xl font-black uppercase tracking-tighter text-white mb-6">Integration: Auth</h3>
              <div className="bg-black border border-error/10 font-mono text-[11px] overflow-hidden">
                <div className="bg-error/10 text-error px-6 py-3 border-b border-error/5 flex items-center gap-3">
                  <span className="font-black">-</span>
                  <span className="line-through">expect(res.status).toBe(200)</span>
                </div>
                <div className="bg-primary/5 text-primary px-6 py-3 flex items-center gap-3">
                  <span className="font-black">+</span>
                  <span className="font-black">expect(res.status).toBe(401)</span>
                </div>
              </div>
            </div>
            
            <div className="flex gap-1">
              <button className="flex-1 py-4 text-[9px] font-black uppercase tracking-widest border border-white/10 text-white/40 hover:text-white transition-all">Logs</button>
              <button className="flex-1 py-4 text-[9px] font-black uppercase tracking-widest bg-primary text-black hover:scale-105 transition-all">Re-run</button>
            </div>
          </div>
        </KanbanColumn>
      </div>
    </div>
  );
}

function KanbanColumn({ index, title, count, icon: Icon, color, children }: any) {
  return (
    <div className="w-[450px] flex flex-col gap-10 border-r border-white/5 px-10">
      <div className="flex items-end justify-between border-b border-white/10 pb-10">
        <div className="flex flex-col gap-4">
          <span className="text-[32px] font-black leading-none text-white/10">{index}</span>
          <div className="flex items-center gap-4">
            <Icon size={18} className={color} />
            <h2 className="text-xs font-black text-white uppercase tracking-[0.4em]">{title}</h2>
          </div>
        </div>
        <span className="text-stroke text-4xl font-black text-white/5">
          {count}
        </span>
      </div>
      <div className="flex flex-col gap-1">
        {children}
      </div>
    </div>
  );
}

function KanbanCard({ tag, tagColor, title, content, ai, tasks, completed }: any) {
  return (
    <div className={cn(
      "p-10 flex flex-col gap-8 transition-all duration-500 group relative border border-transparent hover:bg-white/[0.03]",
      completed && "opacity-20 grayscale",
      !completed && "border-b border-white/5"
    )}>
      <div className="flex justify-between items-start">
        <span className={cn("text-[9px] font-black uppercase tracking-widest px-3 py-1 border", tagColor)}>
          {tag}
        </span>
        {completed ? (
          <CheckCircle2 size={16} className="text-primary" />
        ) : (
          <div className="flex gap-1">
             <div className="w-1 h-1 bg-white/20 rounded-full" />
             <div className="w-1 h-1 bg-white/20 rounded-full" />
             <div className="w-1 h-1 bg-white/20 rounded-full" />
          </div>
        )}
      </div>

      <div>
        <h3 className={cn("text-2xl font-black uppercase tracking-tighter mb-4 transition-all", completed ? "text-white/20 line-through" : "text-white group-hover:text-primary")}>
          {title}
        </h3>
        <p className={cn("text-xs leading-relaxed uppercase font-bold tracking-wide", completed ? "text-white/10 italic" : "text-white/40")}>
          {content}
        </p>
      </div>

      {(ai || tasks) && (
        <div className="flex items-center justify-between pt-8 border-t border-white/5">
          <div className="flex items-center gap-4">
            {ai && (
              <div className="flex items-center gap-2">
                <Brain size={14} className="text-primary" />
                <span className="text-[9px] font-black uppercase tracking-widest text-primary">Neural Node</span>
              </div>
            )}
          </div>
          {tasks && (
            <div className="flex items-center gap-3 text-white/20">
              <MessageSquare size={12} />
              <span className="text-[10px] font-black">{tasks}</span>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

function AlertCircle({ size, className }: any) {
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
      <circle cx="12" cy="12" r="10" />
      <line x1="12" y1="8" x2="12" y2="12" />
      <line x1="12" x2="12.01" y1="16" y2="16" />
    </svg>
  );
}
