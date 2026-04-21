import React, { useState, useEffect } from 'react';
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
  Smartphone,
  CheckCircle,
  AlertCircle,
  Loader2
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';
import { api, ProjectConfig } from '../services/api';

export default function InitializeView() {
  const [archType, setArchType] = useState('web');
  const [projectName, setProjectName] = useState('nexus-core-api');
  const [projectPath, setProjectPath] = useState('~/Lab/nexus-core');
  const [techStack, setTechStack] = useState({
    frontend: 'React / Next.js',
    backend: 'Go / Chi',
    database: 'PostgreSQL',
    mobile: 'React Native'
  });
  const [isInitializing, setIsInitializing] = useState(false);
  const [initializationStatus, setInitializationStatus] = useState<'idle' | 'success' | 'error'>('idle');
  const [statusMessage, setStatusMessage] = useState('');

  const handleInitialize = async () => {
    if (!projectName.trim() || !projectPath.trim()) {
      setStatusMessage('项目名称和路径不能为空');
      setInitializationStatus('error');
      return;
    }

    setIsInitializing(true);
    setInitializationStatus('idle');
    setStatusMessage('正在初始化项目...');

    try {
      const config: ProjectConfig = {
        name: projectName,
        path: projectPath,
        architecture: archType as 'web' | 'cli' | 'micro',
        frontend: techStack.frontend,
        backend: techStack.backend,
        database: techStack.database,
        mobile: techStack.mobile
      };

      const result = await api.initializeProject(config);
      
      setInitializationStatus('success');
      setStatusMessage(`项目 "${projectName}" 初始化成功！`);
      
      // 自动跳转到仪表板视图
      setTimeout(() => {
        window.location.hash = '#dashboard';
      }, 2000);
    } catch (error) {
      setInitializationStatus('error');
      setStatusMessage(`初始化失败: ${error instanceof Error ? error.message : '未知错误'}`);
    } finally {
      setIsInitializing(false);
    }
  };

  const handleDiscard = () => {
    setProjectName('nexus-core-api');
    setProjectPath('~/Lab/nexus-core');
    setArchType('web');
    setTechStack({
      frontend: 'React / Next.js',
      backend: 'Go / Chi',
      database: 'PostgreSQL',
      mobile: 'React Native'
    });
    setInitializationStatus('idle');
    setStatusMessage('');
  };

  const handleTechStackChange = (key: keyof typeof techStack, value: string) => {
    setTechStack(prev => ({
      ...prev,
      [key]: value
    }));
  };

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
            <button 
              onClick={handleDiscard}
              className="text-[10px] font-black uppercase tracking-widest text-white/40 hover:text-white transition-colors"
              disabled={isInitializing}
            >
              Discard Changes
            </button>
            <button 
              onClick={handleInitialize}
              disabled={isInitializing}
              className={cn(
                "bg-primary text-black px-12 py-5 font-black uppercase text-xs tracking-widest hover:scale-105 transition-all shadow-2xl flex items-center gap-3",
                isInitializing && "opacity-70 cursor-not-allowed"
              )}
            >
              {isInitializing ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  Initializing...
                </>
              ) : (
                'Initialize System'
              )}
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
                    value={projectName}
                    onChange={(e) => setProjectName(e.target.value)}
                    className="w-full bg-transparent border-b border-white/20 focus:border-primary text-white font-black text-3xl uppercase tracking-tighter py-3 outline-none transition-all placeholder:text-white/10"
                    placeholder="Enter project name"
                    disabled={isInitializing}
                  />
                </div>
                <div className="space-y-4">
                  <label className="text-[10px] font-black text-white/40 uppercase tracking-[0.2em]">Deployment Path</label>
                  <div className="relative group">
                    <input 
                      type="text" 
                      value={projectPath}
                      onChange={(e) => setProjectPath(e.target.value)}
                      className="w-full bg-transparent border-b border-white/20 focus:border-primary text-white font-mono text-lg py-3 outline-none transition-all"
                      placeholder="~/path/to/project"
                      disabled={isInitializing}
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
                  disabled={isInitializing}
                />
                <ArchitectureCard 
                  id="cli"
                  label="CLI tool"
                  description="System Runtime"
                  icon={Terminal}
                  isSelected={archType === 'cli'}
                  onClick={() => setArchType('cli')}
                  disabled={isInitializing}
                />
                <ArchitectureCard 
                  id="micro"
                  label="Microservice"
                  description="Isolated Node"
                  icon={Settings}
                  isSelected={archType === 'micro'}
                  onClick={() => setArchType('micro')}
                  disabled={isInitializing}
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
              <TechItem 
                icon={Layout} 
                label="Frontend" 
                value={techStack.frontend}
                active={techStack.frontend !== 'React / Next.js'}
                onChange={(value) => handleTechStackChange('frontend', value)}
                disabled={isInitializing}
              />
              <TechItem 
                icon={Cpu} 
                label="Backend" 
                value={techStack.backend}
                active={techStack.backend !== 'Go / Chi'}
                onChange={(value) => handleTechStackChange('backend', value)}
                disabled={isInitializing}
              />
              <TechItem 
                icon={Database} 
                label="Data Layer" 
                value={techStack.database}
                active={techStack.database !== 'PostgreSQL'}
                onChange={(value) => handleTechStackChange('database', value)}
                disabled={isInitializing}
              />
              <TechItem 
                icon={Smartphone} 
                label="Mobile" 
                value={techStack.mobile}
                active={techStack.mobile !== 'React Native'}
                onChange={(value) => handleTechStackChange('mobile', value)}
                disabled={isInitializing}
              />
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

        {/* Status Message */}
        {statusMessage && (
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className={cn(
              "p-6 border flex items-center gap-4",
              initializationStatus === 'success' 
                ? "border-primary/30 bg-primary/5 text-primary" 
                : initializationStatus === 'error'
                ? "border-red-500/30 bg-red-500/5 text-red-500"
                : "border-white/10 bg-white/5 text-white/80"
            )}
          >
            {initializationStatus === 'success' ? (
              <CheckCircle className="w-5 h-5" />
            ) : initializationStatus === 'error' ? (
              <AlertCircle className="w-5 h-5" />
            ) : (
              <Loader2 className="w-5 h-5 animate-spin" />
            )}
            <span className="text-sm font-medium">{statusMessage}</span>
          </motion.div>
        )}
      </div>
    </div>
  );
}

function ArchitectureCard({ label, description, icon: Icon, isSelected, onClick, disabled }: any) {
  return (
    <div 
      onClick={disabled ? undefined : onClick} 
      className={cn(
        "relative p-10 flex flex-col gap-12 transition-all duration-500 border border-white/5",
        isSelected 
          ? "bg-primary text-black border-transparent shadow-[0_40px_80px_rgba(251,223,36,0.15)]" 
          : "hover:bg-white/[0.03] group",
        disabled ? "opacity-40 cursor-not-allowed" : "cursor-pointer"
      )}
    >
      <div className="flex justify-between items-start">
        <Icon 
          size={32} 
          strokeWidth={isSelected ? 3 : 1} 
          className={cn(
            !isSelected && "text-white/20 group-hover:text-white transition-colors",
            disabled && "text-white/10"
          )} 
        />
        {isSelected && <span className="text-[10px] font-black uppercase tracking-widest bg-black text-primary px-3 py-1">Active</span>}
      </div>
      <div>
        <div className="text-2xl font-black uppercase tracking-tighter mb-1">{label}</div>
        <div className={cn("text-[9px] uppercase font-black tracking-widest", isSelected ? "text-black/60" : "text-white/20")}>{description}</div>
      </div>
    </div>
  );
}

function TechItem({ icon: Icon, label, value, active, onChange, disabled }: any) {
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState(value);

  const handleSave = () => {
    if (editValue.trim() && onChange) {
      onChange(editValue);
    }
    setIsEditing(false);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSave();
    } else if (e.key === 'Escape') {
      setEditValue(value);
      setIsEditing(false);
    }
  };

  return (
    <div className="flex flex-col gap-2 group">
      <div className="flex items-center justify-between">
        <div className="text-[9px] font-black uppercase tracking-[0.3em] text-white/20 group-hover:text-primary transition-colors">{label}</div>
        {!disabled && (
          <button 
            onClick={() => setIsEditing(true)}
            className="text-white/10 hover:text-white transition-colors"
          >
            <Edit size={12} />
          </button>
        )}
      </div>
      
      {isEditing ? (
        <div className="relative">
          <input
            type="text"
            value={editValue}
            onChange={(e) => setEditValue(e.target.value)}
            onKeyDown={handleKeyDown}
            onBlur={handleSave}
            className="w-full bg-transparent border-b border-primary text-white font-black text-lg uppercase tracking-tighter py-1 outline-none"
            autoFocus
          />
          <div className="absolute right-0 top-1/2 -translate-y-1/2 flex gap-2">
            <button 
              onClick={handleSave}
              className="text-[8px] font-black uppercase tracking-widest text-primary hover:text-white transition-colors"
            >
              Save
            </button>
            <button 
              onClick={() => {
                setEditValue(value);
                setIsEditing(false);
              }}
              className="text-[8px] font-black uppercase tracking-widest text-white/40 hover:text-white transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      ) : (
        <div className={cn(
          "text-lg font-black uppercase tracking-tighter transition-all",
          active ? "text-primary text-stroke !text-white" : "text-white/40"
        )}>
          {value}
        </div>
      )}
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
