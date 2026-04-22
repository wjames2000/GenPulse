/**
 * @license
 * SPDX-License-Identifier: Apache-2.0
 */

import React, { useState, useEffect, useCallback, useRef } from 'react';
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime';
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
  Code2,
  Command,
} from 'lucide-react';
import { motion, AnimatePresence } from 'motion/react';
import { cn } from './utils';
import { View } from './types';
import { useKeyboardShortcuts } from './hooks/useKeyboardShortcuts';
import { PageTransition } from './components/common/PageTransition';
import { StatusDot } from './components/common/StatusDot';
import DashboardView from './components/DashboardView';
import InitializeView from './components/InitializeView';
import SettingsView from './components/SettingsView';
import NeuralAssetsView from './components/NeuralAssetsView';
import KanbanView from './components/KanbanView';
import SkillsView from './components/SkillsView';
import DiffView from './components/DiffView';
import MCPConfigView from './components/MCPConfigView';
import MonitoringDashboard from './components/MonitoringDashboard';
import StartupScreen from './components/StartupScreen';
import { api } from './services/api';
import * as AppBindings from '../wailsjs/go/main/App';

export default function App() {
  const [currentView, setCurrentView] = useState<View>('dashboard');
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);
  const [systemStatus, setSystemStatus] = useState<'healthy' | 'degraded' | 'error'>('healthy');
  const [activeAgents, setActiveAgents] = useState(0);
  const [searchFocused, setSearchFocused] = useState(false);
  const [viewLoading, setViewLoading] = useState(false);
  const [startupComplete, setStartupComplete] = useState(false);
  const searchInputRef = useRef<HTMLInputElement>(null);
  const prevViewRef = useRef<View>(currentView);

  useEffect(() => {
    const onPhase1 = () => setStartupComplete(true);
    EventsOn('startup:phase1_complete', onPhase1);

    const interval = setInterval(async () => {
      try {
        const phase = await AppBindings.GetStartupPhase();
        if (phase.phase >= 1) {
          setStartupComplete(true);
          clearInterval(interval);
        }
      } catch {
      }
    }, 200);

    return () => {
      EventsOff('startup:phase1_complete');
      clearInterval(interval);
    };
  }, []);

  const checkSystemStatus = useCallback(async () => {
    try {
      const health = await api.healthCheck();
      setSystemStatus(health.status === 'healthy' ? 'healthy' : 'degraded');

      const agentsStatus = await api.getAllAgentsStatus();
      const activeCount = Object.values(agentsStatus).filter((agent: any) =>
        agent?.state === 'active'
      ).length;
      setActiveAgents(activeCount);
    } catch (error) {
      setSystemStatus('error');
    }
  }, []);

  useEffect(() => {
    const onPhase2 = () => {
      checkSystemStatus();
    };
    EventsOn('startup:phase2_complete', onPhase2);
    return () => { EventsOff('startup:phase2_complete'); };
  }, [checkSystemStatus]);

  useEffect(() => {
    const interval = setInterval(checkSystemStatus, 30000);
    return () => clearInterval(interval);
  }, [checkSystemStatus]);

  const handleViewChange = useCallback((view: View) => {
    if (view === currentView) return;
    prevViewRef.current = currentView;
    setViewLoading(true);
    setCurrentView(view);
    const timer = setTimeout(() => setViewLoading(false), 200);
    return () => clearTimeout(timer);
  }, [currentView]);

  useKeyboardShortcuts({
    search: () => {
      searchInputRef.current?.focus();
      setSearchFocused(true);
    },
    toggleSidebar: () => setIsSidebarOpen(prev => !prev),
    escape: () => {
      if (searchFocused) {
        searchInputRef.current?.blur();
        setSearchFocused(false);
      }
    },
    nav1: () => handleViewChange('dashboard'),
    nav2: () => handleViewChange('pipeline'),
    nav3: () => handleViewChange('kanban'),
    nav4: () => handleViewChange('neural'),
    nav5: () => handleViewChange('skills'),
    nav6: () => handleViewChange('history'),
    nav7: () => handleViewChange('monitoring'),
    nav8: () => handleViewChange('settings'),
    nav9: () => handleViewChange('mcp-config'),
  });

  const navItems = [
    { id: 'dashboard' as View, icon: GitBranch, label: '项目流水线', badge: activeAgents > 0 ? activeAgents.toString() : undefined },
    { id: 'pipeline' as View, icon: Plus, label: '新建项目' },
    { id: 'kanban' as View, icon: Layout, label: '任务看板' },
    { id: 'neural' as View, icon: Brain, label: '认知资产' },
    { id: 'skills' as View, icon: Code2, label: '技能库' },
    { id: 'history' as View, icon: History, label: '开发历史' },
    { id: 'monitoring' as View, icon: Activity, label: '监控仪表盘' },
    { id: 'settings' as View, icon: Settings, label: '系统设置' },
    { id: 'mcp-config' as View, icon: TerminalIcon, label: 'MCP 配置' },
  ];

  const renderView = () => {
    switch (currentView) {
      case 'dashboard': return <DashboardView />;
      case 'pipeline': return <InitializeView />;
      case 'settings': return <SettingsView />;
      case 'neural': return <NeuralAssetsView />;
      case 'kanban': return <KanbanView />;
      case 'history': return <DiffView />;
      case 'skills': return <SkillsView />;
      case 'monitoring': return <MonitoringDashboard />;
      case 'mcp-config': return <MCPConfigView />;
      default: return <DashboardView />;
    }
  };

  return (
    <>
      {!startupComplete && <StartupScreen />}
      <div className="flex bg-background text-on-surface min-h-screen font-sans">
      {/* Sidebar */}
      <motion.nav
        layout
        transition={{ duration: 0.3, ease: [0.4, 0, 0.2, 1] }}
        className={cn(
          "fixed left-0 top-0 h-full bg-background border-r border-white/10 z-40 flex flex-col pt-10 pb-8 overflow-hidden",
          isSidebarOpen ? "w-64" : "w-20"
        )}
      >
        <div className="px-8 mb-12 flex items-center space-x-3">
          <div className="text-xl font-black tracking-tighter text-primary shrink-0">GenPulse™</div>
          <AnimatePresence>
            {isSidebarOpen && (
              <motion.div
                initial={{ opacity: 0, width: 0 }}
                animate={{ opacity: 1, width: 'auto' }}
                exit={{ opacity: 0, width: 0 }}
                transition={{ duration: 0.2 }}
                className="flex flex-col overflow-hidden"
              >
                <div className="text-[10px] uppercase font-bold tracking-[0.3em] text-white/40 leading-none whitespace-nowrap">Cognitive</div>
                <div className="text-[10px] uppercase font-bold tracking-[0.3em] text-white/40 leading-none mt-1 whitespace-nowrap">Lab / V2.4</div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>

        <div className="px-6 mb-12">
          <div className="flex items-center gap-3 mb-4">
            <StatusDot
              status={systemStatus === 'healthy' ? 'healthy' : systemStatus === 'degraded' ? 'warning' : 'error'}
              animate={systemStatus === 'healthy'}
            />
            <span className="text-[9px] font-black uppercase tracking-widest text-white/40 whitespace-nowrap">
              {systemStatus === 'healthy' ? 'System Online' :
               systemStatus === 'degraded' ? 'System Degraded' :
               'System Error'}
            </span>
          </div>
          <motion.button
            whileHover={{ scale: 1.03 }}
            whileTap={{ scale: 0.97 }}
            onClick={() => handleViewChange('pipeline')}
            className={cn(
              "w-full bg-primary text-black py-4 px-4 font-black uppercase text-xs tracking-widest transition-all shadow-xl rounded-none flex items-center justify-center gap-3",
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
          </motion.button>
        </div>

        <div className="flex-1 flex flex-col space-y-1 px-4 uppercase text-[10px] font-bold tracking-[0.2em] text-white/60">
          {navItems.map((item, index) => (
            <motion.button
              key={item.id}
              whileHover={{ scale: 1.02, x: 4 }}
              whileTap={{ scale: 0.98 }}
              onClick={() => handleViewChange(item.id)}
              className={cn(
                "py-4 px-4 flex items-center justify-between transition-colors relative group rounded-sm",
                currentView === item.id
                  ? "text-primary bg-white/[0.07]"
                  : "hover:text-white hover:bg-white/[0.03]"
              )}
            >
              <div className="flex items-center space-x-4">
                <item.icon size={16} className="shrink-0" />
                {isSidebarOpen && <span className="whitespace-nowrap">{item.label}</span>}
              </div>
              
              <AnimatePresence>
                {isSidebarOpen && item.badge && (
                  <motion.span
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    exit={{ scale: 0 }}
                    className="text-[8px] font-black bg-primary text-black px-2 py-1 rounded-full min-w-[20px] text-center"
                  >
                    {item.badge}
                  </motion.span>
                )}
              </AnimatePresence>
              
              {currentView === item.id && (
                <motion.div
                  layoutId="activeNav"
                  className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-primary rounded-r-full"
                  transition={{ type: 'spring', stiffness: 300, damping: 30 }}
                />
              )}
            </motion.button>
          ))}
        </div>

        <div className="mt-auto px-8 space-y-4 pt-10 border-t border-white/10 uppercase text-[9px] font-black tracking-[0.3em] text-white/30">
          <a href="#" className="flex items-center space-x-3 hover:text-white transition-colors group">
            <FileText size={14} className="group-hover:scale-110 transition-transform" />
            {isSidebarOpen && <span>Documentation</span>}
          </a>
          <a href="#" className="flex items-center space-x-3 hover:text-white transition-colors group">
            <HelpCircle size={14} className="group-hover:scale-110 transition-transform" />
            {isSidebarOpen && <span>Support</span>}
          </a>
        </div>
      </motion.nav>

      {/* Main Content Area */}
      <motion.div
        layout
        transition={{ duration: 0.3, ease: [0.4, 0, 0.2, 1] }}
        className={cn(
          "flex-1 flex flex-col min-h-screen",
          isSidebarOpen ? "ml-64" : "ml-20"
        )}
      >
        {/* Header */}
        <header className="h-24 flex items-center justify-between px-12 sticky top-0 z-50 bg-background/80 backdrop-blur-xl border-b border-white/10">
          <div className="flex items-center space-x-8">
            <motion.button
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.9 }}
              onClick={() => setIsSidebarOpen(!isSidebarOpen)}
              className="p-2 text-white/40 hover:text-white transition-colors"
              title={isSidebarOpen ? 'Collapse sidebar (⌘B)' : 'Expand sidebar (⌘B)'}
            >
              <Menu size={20} />
            </motion.button>
            <div className="text-sm font-black uppercase tracking-[0.4em] text-white/80">
              Control Panel
            </div>
            
            <div className={cn(
              "relative hidden lg:flex items-center group transition-all duration-200",
              searchFocused && "scale-[1.02]"
            )}>
              <Search size={14} className={cn(
                "absolute left-0 transition-colors z-10",
                searchFocused ? "text-primary" : "text-white/40 group-hover:text-white/60"
              )} />
              <input
                ref={searchInputRef}
                type="text"
                placeholder="SEARCH COGNITIVE ASSETS"
                className="bg-transparent border-b border-transparent focus:border-primary/30 pl-6 pr-10 py-2 text-[10px] font-black tracking-widest text-white focus:ring-0 outline-none w-64 opacity-60 focus:opacity-100 transition-all placeholder:text-white/20"
                onFocus={() => setSearchFocused(true)}
                onBlur={() => setSearchFocused(false)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    const target = e.target as HTMLInputElement;
                    console.log('Search for:', target.value);
                  }
                  if (e.key === 'Escape') {
                    (e.target as HTMLInputElement).blur();
                  }
                }}
              />
              <kbd className="absolute right-0 text-[8px] font-mono text-white/20 border border-white/10 px-1.5 py-0.5 rounded pointer-events-none">
                ⌘K
              </kbd>
            </div>
          </div>

          <div className="flex items-center space-x-10">
            <div className="flex space-x-6 text-white/40">
              <motion.button
                whileHover={{ scale: 1.1 }}
                whileTap={{ scale: 0.9 }}
                className="hover:text-white transition-colors relative"
              >
                <Bell size={18} />
                <span className="absolute -top-1 -right-1 w-1.5 h-1.5 bg-primary rounded-full animate-pulse" />
              </motion.button>
              <motion.button
                whileHover={{ scale: 1.1 }}
                whileTap={{ scale: 0.9 }}
                className="hover:text-white transition-colors"
              >
                <Puzzle size={18} />
              </motion.button>
            </div>
            
            <div className="h-6 w-[1px] bg-white/10" />
            
            <motion.button
              whileHover={{ scale: 1.03 }}
              whileTap={{ scale: 0.97 }}
              onClick={async () => {
                try {
                  await api.logMessage('info', 'Manual agent deployment triggered');
                  console.log('Deploying agent...');
                } catch (error) {
                  console.error('Failed to deploy agent:', error);
                }
              }}
              className="bg-primary text-black px-8 py-3 font-black uppercase text-[10px] tracking-widest transition-all shadow-lg rounded-none"
            >
              Deploy Agent
            </motion.button>

            <motion.div
              whileHover={{ scale: 1.05 }}
              className="w-10 h-10 border border-white/20 overflow-hidden cursor-pointer hover:border-primary transition-all p-1"
            >
              <img
                src="https://picsum.photos/seed/user/100/100"
                alt="Profile"
                referrerPolicy="no-referrer"
                className="w-full h-full object-cover grayscale hover:grayscale-0 transition-all duration-500"
              />
            </motion.div>
          </div>
        </header>

        {/* View Transition Canvas */}
        <main className="flex-1 overflow-hidden relative bg-gradient-subtle">
          <PageTransition uniqueKey={currentView} direction="up" duration={0.25}>
            {viewLoading ? (
              <div className="p-12 space-y-6 animate-pulse">
                <div className="h-6 bg-white/5 w-1/4" />
                <div className="h-3 bg-white/5 w-1/3" />
                <div className="grid grid-cols-4 gap-4 mt-12">
                  {Array.from({ length: 4 }).map((_, i) => (
                    <div key={i} className="h-40 bg-white/[0.03] border border-white/5" />
                  ))}
                </div>
                <div className="h-64 bg-white/[0.02] border border-white/5 mt-8" />
              </div>
            ) : (
              renderView()
            )}
          </PageTransition>
        </main>
      </motion.div>
    </div>
    </>
  );
}
