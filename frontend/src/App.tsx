/**
 * @license
 * SPDX-License-Identifier: Apache-2.0
 */

import React, { useState } from 'react';
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
  CheckCircle2
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

export default function App() {
  const [currentView, setCurrentView] = useState<View>('dashboard');
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);

  const navItems = [
    { id: 'neural' as View, icon: Beaker, label: '认知实验室' },
    { id: 'dashboard' as View, icon: GitBranch, label: '项目流水线' },
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
          <div className="text-xl font-black tracking-tighter text-primary">VOID™</div>
          {isSidebarOpen && (
            <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="flex flex-col">
              <div className="text-[10px] uppercase font-bold tracking-[0.3em] text-white/40 leading-none">Cognitive</div>
              <div className="text-[10px] uppercase font-bold tracking-[0.3em] text-white/40 leading-none mt-1">Lab / V2.4</div>
            </motion.div>
          )}
        </div>

        <div className="px-6 mb-12">
          <button 
            onClick={() => setCurrentView('pipeline')}
            className={cn(
              "w-full bg-primary-container text-on-primary-container py-4 px-4 font-black uppercase text-xs tracking-widest hover:scale-105 transition-all shadow-xl rounded-none",
              !isSidebarOpen && "px-0 justify-center"
            )}
          >
            {isSidebarOpen ? "New Pipeline" : <Plus size={18} />}
          </button>
        </div>

        <div className="flex-1 flex flex-col space-y-2 px-4 uppercase text-[10px] font-bold tracking-[0.2em] text-white/60">
          {navItems.map((item) => (
            <button
              key={item.id}
              onClick={() => setCurrentView(item.id)}
              className={cn(
                "py-4 px-4 flex items-center space-x-4 transition-all hover:text-white relative",
                currentView === item.id ? "text-primary bg-white/5" : ""
              )}
            >
              <item.icon size={16} className="shrink-0" />
              {isSidebarOpen && <span className="whitespace-nowrap">{item.label}</span>}
              {currentView === item.id && (
                <motion.div 
                  layoutId="activeNav"
                  className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-primary rounded-r-full"
                />
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
            
            <button className="bg-primary text-black px-8 py-3 font-black uppercase text-[10px] tracking-widest hover:scale-105 transition-all shadow-lg rounded-none">
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
