import React, { useState } from 'react';
import { 
  FolderOpen, 
  Terminal, 
  Layout, 
  Network as Hub, 
  Cpu as Memory, 
  Database, 
  Edit, 
  Ruler as Schema, 
  Code2 as DataObject, 
  Layout as WebAsset, 
  Users as Groups, 
  Play, 
  Plus,
  Rocket,
  Globe,
  Settings,
  FolderCog as FolderManaged,
  Ruler as Architecture,
  Layers,
  Sparkles as AutoAwesome,
  Folder as FolderIcon,
  Cpu,
  Smartphone
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';

export default function InitializeView() {
  const [archType, setArchType] = useState('web');

  return (
    <div className="flex-1 overflow-y-auto p-12 bg-[#0A0A0A] custom-scrollbar min-h-screen">
      <div className="max-w-6xl mx-auto space-y-16 pb-20">
        <header className="flex flex-col md:flex-row md:items-end justify-between gap-12 border-b border-white/10 pb-12">
          <div className="flex-1">
            <span className="text-[10px] uppercase tracking-[0.5em] text-white/40 block mb-6">Pipeline Initialization / Phase 01</span>
            <h1 className="text-[140px] leading-[0.8] font-black tracking-tighter uppercase">
              Start<br/>Nexus
            </h1>
          </div>
          <div className="flex items-center gap-6 pb-4">
            <button className="text-[10px] font-black uppercase tracking-widest text-white/40 hover:text-white transition-colors">
              Discard Changes
            </button>
            <button className="bg-primary text-black px-12 py-5 font-black uppercase text-xs tracking-widest hover:scale-105 transition-all shadow-2xl">
              Initialize System
            </button>
          </div>
        </header>

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-1 border-b border-white/10">
          <div className="lg:col-span-8 border-r border-white/10 pr-12 py-12">
            <section className="mb-20">
              <div className="text-[10px] uppercase font-black tracking-[0.4em] text-primary mb-10 flex items-center gap-4">
                <span className="w-8 h-[1px] bg-primary"></span>
                01 / Basic Configuration
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
                <div className="space-y-4">
                  <label className="text-[10px] font-black text-white/40 uppercase tracking-[0.2em]">Project Name</label>
                  <input 
                    type="text" 
                    defaultValue="nexus-core-api"
                    className="w-full bg-transparent border-b border-white/20 focus:border-primary text-white font-black text-3xl uppercase tracking-tighter py-3 outline-none transition-all placeholder:text-white/10"
                  />
                </div>
                <div className="space-y-4">
                  <label className="text-[10px] font-black text-white/40 uppercase tracking-[0.2em]">Deployment Path</label>
                  <div className="relative group">
                    <input 
                      type="text" 
                      defaultValue="~/Lab/nexus-core"
                      className="w-full bg-transparent border-b border-white/20 focus:border-primary text-white font-mono text-lg py-3 outline-none transition-all"
                    />
                  </div>
                </div>
              </div>
            </section>

            <section>
              <div className="text-[10px] uppercase font-black tracking-[0.4em] text-primary mb-10 flex items-center gap-4">
                <span className="w-8 h-[1px] bg-primary"></span>
                02 / Architecture
              </div>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-1">
                <ArchitectureCard 
                  id="web"
                  label="Web App"
                  description="Fullstack SPA"
                  icon={Layout}
                  isSelected={archType === 'web'}
                  onClick={() => setArchType('web')}
                />
                <ArchitectureCard 
                  id="cli"
                  label="CLI tool"
                  description="System Runtime"
                  icon={Terminal}
                  isSelected={archType === 'cli'}
                  onClick={() => setArchType('cli')}
                />
                <ArchitectureCard 
                  id="micro"
                  label="Microservice"
                  description="Isolated Node"
                  icon={Settings}
                  isSelected={archType === 'micro'}
                  onClick={() => setArchType('micro')}
                />
              </div>
            </section>
          </div>

          <div className="lg:col-span-4 p-12">
            <div className="text-[10px] uppercase font-black tracking-[0.4em] text-primary mb-10 flex items-center gap-4">
              <span className="w-8 h-[1px] bg-primary"></span>
              03 / Tech Matrix
            </div>
            <div className="space-y-12">
              <TechItem icon={Layout} label="Frontend" value="React / Next.js" />
              <TechItem icon={Cpu} label="Backend" value="Go / Chi" active />
              <TechItem icon={Database} label="Data Layer" value="PostgreSQL" />
              <TechItem icon={Smartphone} label="Mobile" value="React Native" />
            </div>
          </div>
        </div>

        <div className="pt-20">
          <div className="text-[10px] uppercase font-black tracking-[0.4em] text-primary mb-12 flex items-center gap-4">
            <span className="w-8 h-[1px] bg-primary"></span>
            04 / Autonomous Team
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-1">
            <AgentConfigCard 
              index="01"
              title="Architect"
              subtitle="System Lead"
              description="Infrastructure design and domain modeling."
              logo="A"
            />
            <AgentConfigCard 
              index="02"
              title="Backend"
              subtitle="Core Dev"
              description="Business logic and persistence layers."
              logo="B"
            />
            <AgentConfigCard 
              index="03"
              title="Frontend"
              subtitle="UI Expert"
              description="Interface engineering and state management."
              logo="F"
              disabled
            />
          </div>
        </div>
      </div>
    </div>
  );
}

function ArchitectureCard({ label, description, icon: Icon, isSelected, onClick }: any) {
  return (
    <div 
      onClick={onClick} 
      className={cn(
        "relative cursor-pointer p-10 flex flex-col gap-12 transition-all duration-500 border border-white/5",
        isSelected 
          ? "bg-primary text-black border-transparent shadow-[0_40px_80px_rgba(251,223,36,0.15)]" 
          : "hover:bg-white/[0.03] group"
      )}
    >
      <div className="flex justify-between items-start">
        <Icon size={32} strokeWidth={isSelected ? 3 : 1} className={cn(!isSelected && "text-white/20 group-hover:text-white transition-colors")} />
        {isSelected && <span className="text-[10px] font-black uppercase tracking-widest bg-black text-primary px-3 py-1">Active</span>}
      </div>
      <div>
        <div className="text-2xl font-black uppercase tracking-tighter mb-1">{label}</div>
        <div className={cn("text-[9px] uppercase font-black tracking-widest", isSelected ? "text-black/60" : "text-white/20")}>{description}</div>
      </div>
    </div>
  );
}

function TechItem({ icon: Icon, label, value, active }: any) {
  return (
    <div className="flex flex-col gap-2 group">
      <div className="flex items-center justify-between">
        <div className="text-[9px] font-black uppercase tracking-[0.3em] text-white/20 group-hover:text-primary transition-colors">{label}</div>
        <Edit size={12} className="text-white/10 cursor-pointer hover:text-white transition-colors" />
      </div>
      <div className={cn(
        "text-lg font-black uppercase tracking-tighter transition-all",
        active ? "text-primary text-stroke !text-white" : "text-white/40"
      )}>
        {value}
      </div>
    </div>
  );
}

function AgentConfigCard({ index, title, subtitle, description, logo, disabled }: any) {
  return (
    <div className={cn(
      "p-12 border border-white/5 flex flex-col gap-10 transition-all duration-500",
      disabled ? "opacity-20 grayscale" : "hover:bg-white/[0.03]"
    )}>
      <div className="flex justify-between items-start">
        <span className="text-[40px] font-black leading-none text-white/10">{index}</span>
        <div className={cn("w-14 h-14 font-black text-2xl flex items-center justify-center border", disabled ? "border-white/10" : "border-primary text-primary")}>
          {logo}
        </div>
      </div>
      
      <div>
        <h4 className="text-2xl font-black uppercase tracking-tighter mb-1">{title}</h4>
        <p className="text-[9px] font-black uppercase tracking-widest text-white/40">{subtitle}</p>
        <p className="text-xs text-white/30 mt-6 leading-relaxed max-w-[200px]">
          {description}
        </p>
      </div>

      <div className="mt-auto pt-10 border-t border-white/10 flex items-center justify-between">
        <span className="text-[9px] font-black uppercase tracking-widest text-white/20">Model Engine</span>
        <span className="text-[10px] font-black uppercase tracking-widest text-primary">Gemini 2.5 Pro</span>
      </div>
    </div>
  );
}
