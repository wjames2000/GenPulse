/**
 * @license
 * SPDX-License-Identifier: Apache-2.0
 */

import React, { useState, useEffect } from 'react';
import { 
  Beaker, 
  GitBranch, 
  History, 
  Brain, 
  Settings, 
  FileText, 
  HelpCircle, 
  Plus, 
  Bell, 
  Puzzle, 
  Search, 
  Play, 
  Globe, 
  Rocket,
  Menu,
  MoreVertical,
  Layout,
  Terminal as TerminalIcon,
  Diff as DiffIcon,
  CheckCircle2,
  Activity,
  Cpu,
  Database,
  Code2
} from 'lucide-react';
import { motion, AnimatePresence } from 'motion/react';
import { cn } from './utils';
import { View } from './types';
import DashboardView from './components/DashboardView';
import InitializeView from './components/InitializeView';
import SettingsView from './components/SettingsView';
import NeuralAssetsView from './components/NeuralAssetsView';
import KanbanView from './components/KanbanView';
import SkillsView from './components/SkillsView';
import DiffView from './components/DiffView';
import { api } from './services/api';

export default function App() {
  const [currentView, setCurrentView] = useState<View>('dashboard');
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);
  const [systemStatus, setSystemStatus] = useState<'healthy' | 'degraded' | 'error'>('healthy');
  const [activeAgents, setActiveAgents] = useState(0);

  useEffect(() => {
    const checkSystemStatus = async () => {
      try {
        const health = await api.healthCheck();
        if (health.status === 'healthy') {
          setSystemStatus('healthy');
        } else {
          setSystemStatus('degraded');
        }

        const agentsStatus = await api.getAllAgentsStatus();
        const activeCount = Object.values(agentsStatus).filter((agent: any) => 
          agent?.state === 'active'
        ).length;
        setActiveAgents(activeCount);
      } catch (error) {
        setSystemStatus('error');
        console.error('Failed to check system status:', error);
      }
    };

    checkSystemStatus();
    const interval = setInterval(checkSystemStatus, 30000); // Check every 30 seconds

    return () => clearInterval(interval);
  }, []);

  const navItems = [
    { id: 'dashboard' as View, icon: GitBranch, label: '项目流水线', badge: activeAgents > 0 ? activeAgents.toString() : undefined },
    { id: 'pipeline' as View, icon: Plus, label: '新建项目' },
    { id: 'kanban' as View, icon: Layout, label: '任务看板' },
    { id: 'neural' as View, icon: Brain, label: '认知资产' },
    { id: 'skills' as View, icon: Code2, label: '技能库' },
    { id: 'history' as View, icon: History, label: '开发历史' },
    { id: 'settings' as View, icon: Settings, label: '系统设置' },
  ];

  return (
    <div className="flex bg-background text-on-surface min-h-screen font-sans">
      {/* Sidebar */}
      <nav className={cn(
        "fixed left-0 top-0 h-full bg-background border-r border-white/10 transition-all duration-300 z-40 flex flex-col pt-10 pb-8",
        isSidebarOpen ? "w-64" : "w-20"
      )}>
        <div className="px-8 mb-12 flex items-center space-x-3">
          <div className="text-xl font-black tracking-tighter text-primary">GenPulse™</div>
          {isSidebarOpen && (
            <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="flex flex-col">
              <div className="text-[10px] uppercase font-bold tracking-[0.3em] text-white/40 leading-none">Cognitive</div>
              <div className="text-[10px] uppercase font-bold tracking-[0.3em] text-white/40 leading-none mt-1">Lab / V2.4</div>
            </motion.div>
          )}
        </div>

        <div className="px-6 mb-12">
          <div className="flex items-center gap-3 mb-4">
            <div className={cn(
              "w-2 h-2 rounded-full",
              systemStatus === 'healthy' ? "bg-primary animate-pulse" :
              systemStatus === 'degraded' ? "bg-yellow-500" :
              "bg-red-500"
            )} />
            <span className="text-[9px] font-black uppercase tracking-widest text-white/40">
              {systemStatus === 'healthy' ? 'System Online' :
               systemStatus === 'degraded' ? 'System Degraded' :
               'System Error'}
            </span>
          </div>
          <button 
            onClick={() => setCurrentView('pipeline')}
            className={cn(
              "w-full bg-primary text-black py-4 px-4 font-black uppercase text-xs tracking-widest hover:scale-105 transition-all shadow-xl rounded-none flex items-center justify-center gap-3",
              !isSidebarOpen && "px-0"
            )}
          >
            {isSidebarOpen ? (
              <>
                <Plus size={16} />
                New Pipeline
              </>
            ) : (
              <Plus size={18} />
            )}
          </button>
        </div>

        <div className="flex-1 flex flex-col space-y-1 px-4 uppercase text-[10px] font-bold tracking-[0.2em] text-white/60">
          {navItems.map((item) => (
            <button
              key={item.id}
              onClick={() => setCurrentView(item.id)}
              className={cn(
                "py-4 px-4 flex items-center justify-between transition-all hover:text-white hover:bg-white/5 relative group",
                currentView === item.id ? "text-primary bg-white/5" : ""
              )}
            >
              <div className="flex items-center space-x-4">
                <item.icon size={16} className="shrink-0" />
                {isSidebarOpen && <span className="whitespace-nowrap">{item.label}</span>}
              </div>
              
              {isSidebarOpen && item.badge && (
                <span className="text-[8px] font-black bg-primary text-black px-2 py-1 rounded-full min-w-[20px] text-center">
                  {item.badge}
                </span>
              )}
              
              {currentView === item.id && (
                <motion.div 
                  layoutId="activeNav"
                  className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-primary rounded-r-full"
                />
              )}
              
              {!isSidebarOpen && currentView === item.id && (
                <div className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-primary rounded-r-full" />
              )}
            </button>
          ))}
        </div>

        <div className="mt-auto px-8 space-y-4 pt-10 border-t border-white/10 uppercase text-[9px] font-black tracking-[0.3em] text-white/30">
          <a href="#" className="flex items-center space-x-3 hover:text-white transition-colors">
            <FileText size={14} />
            {isSidebarOpen && <span>Documentation</span>}
          </a>
          <a href="#" className="flex items-center space-x-3 hover:text-white transition-colors">
            <HelpCircle size={14} />
            {isSidebarOpen && <span>Support</span>}
          </a>
        </div>
      </nav>

      {/* Main Content Area */}
      <div className={cn(
        "flex-1 flex flex-col transition-all duration-300 min-h-screen",
        isSidebarOpen ? "ml-64" : "ml-20"
      )}>
        {/* Header */}
        <header className="h-24 flex items-center justify-between px-12 sticky top-0 z-50 bg-background border-b border-white/10">
          <div className="flex items-center space-x-8">
            <button 
              onClick={() => setIsSidebarOpen(!isSidebarOpen)}
              className="p-2 text-white/40 hover:text-white transition-all"
            >
              <Menu size={20} />
            </button>
            <div className="text-sm font-black uppercase tracking-[0.4em] text-white/80">
              Control Panel
            </div>
            
            <div className="relative hidden lg:flex items-center text-white/40 group">
              <Search size={14} className="absolute left-0 group-focus-within:text-primary transition-colors" />
              <input 
                type="text" 
                placeholder="SEARCH COGNITIVE ASSETS" 
                className="bg-transparent border-none pl-6 pr-4 py-2 text-[10px] font-black tracking-widest text-white focus:ring-0 outline-none w-64 opacity-60 focus:opacity-100 transition-all placeholder:text-white/20"
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    const target = e.target as HTMLInputElement;
                    console.log('Search for:', target.value);
                    // 这里可以添加搜索功能
                  }
                }}
              />
            </div>
          </div>

          <div className="flex items-center space-x-10">
            <div className="flex space-x-6 text-white/40">
              <button className="hover:text-white transition-all relative">
                <Bell size={18} />
                <span className="absolute -top-1 -right-1 w-1.5 h-1.5 bg-primary rounded-full" />
              </button>
              <button className="hover:text-white transition-all">
                <Puzzle size={18} />
              </button>
            </div>
            
            <div className="h-6 w-[1px] bg-white/10" />
            
            <button 
              onClick={async () => {
                try {
                  await api.logMessage('info', 'Manual agent deployment triggered');
                  // 这里可以添加部署逻辑
                  console.log('Deploying agent...');
                } catch (error) {
                  console.error('Failed to deploy agent:', error);
                }
              }}
              className="bg-primary text-black px-8 py-3 font-black uppercase text-[10px] tracking-widest hover:scale-105 transition-all shadow-lg rounded-none"
            >
              Deploy Agent
            </button>

            <div className="w-10 h-10 border border-white/20 overflow-hidden cursor-pointer hover:border-primary transition-all p-1">
              <img 
                src="https://picsum.photos/seed/user/100/100" 
                alt="Profile" 
                referrerPolicy="no-referrer"
                className="w-full h-full object-cover grayscale active:grayscale-0 transition-all"
              />
            </div>
          </div>
        </header>

        {/* View Transition Canvas */}
        <main className="flex-1 overflow-hidden relative">
          <AnimatePresence mode="wait">
            <motion.div
              key={currentView}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.2, ease: "easeOut" }}
              className="h-full"
            >
              {currentView === 'dashboard' && <DashboardView />}
              {currentView === 'pipeline' && <InitializeView />}
              {currentView === 'settings' && <SettingsView />}
              {currentView === 'neural' && <NeuralAssetsView />}
              {currentView === 'kanban' && <KanbanView />}
              {currentView === 'history' && <DiffView />}
              {currentView === 'skills' && <SkillsView />}
            </motion.div>
          </AnimatePresence>
        </main>
      </div>
    </div>
  );
}
