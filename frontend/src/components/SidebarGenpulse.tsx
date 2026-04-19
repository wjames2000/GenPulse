import React from 'react';
import { useAppStore } from '../stores/appStore';

const SidebarGenpulse: React.FC = () => {
  const { currentView, setCurrentView } = useAppStore();

  const navItems = [
    { id: 'dashboard', label: 'Dashboard', icon: 'grid_view' },
    { id: 'projects', label: 'Pipeline', icon: 'account_tree' },
    { id: 'agents', label: 'Agents', icon: 'psychology' },
    { id: 'skills', label: 'Skills', icon: 'memory' },
    { id: 'memory', label: 'Memory', icon: 'psychology' },
    { id: 'kanban', label: 'Kanban', icon: 'dashboard' },
    { id: 'terminal', label: 'Terminal', icon: 'terminal' },
    { id: 'settings', label: 'Settings', icon: 'settings' },
  ];

  const footerItems = [
    { label: 'Documentation', icon: 'description' },
    { label: 'Support', icon: 'contact_support' },
  ];

  return (
    <nav className="h-screen w-64 flex flex-col bg-[#131318] border-r border-[#454555]/15">
      {/* Header */}
      <div className="px-6 py-6 flex items-center gap-3">
        <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-primary-container to-inverse-primary flex items-center justify-center shadow-[0_0_15px_rgba(91,95,255,0.3)]">
          <span className="material-symbols-outlined text-on-primary-container text-sm" data-weight="fill">token</span>
        </div>
        <div>
          <h1 className="text-[#C0C1FF] font-bold text-lg leading-tight tracking-tight">Genpulse</h1>
          <p className="font-['Inter'] text-xs tracking-wide text-[#908FA1]">AI Development</p>
        </div>
      </div>

      {/* CTA Button */}
      <div className="px-4 mb-6">
        <button className="w-full bg-primary-container hover:bg-surface-bright text-on-primary-container py-2.5 px-4 rounded-lg flex items-center justify-center gap-2 transition-all duration-200">
          <span className="material-symbols-outlined text-sm">add</span>
          <span className="font-medium text-sm">New Pipeline</span>
        </button>
      </div>

      {/* Main Navigation */}
      <div className="flex-1 px-2 space-y-1">
        {navItems.map((item) => (
          <button
            key={item.id}
            onClick={() => setCurrentView(item.id)}
            className={`w-full flex items-center gap-3 py-3 px-4 rounded-lg transition-all duration-200 cursor-pointer active:opacity-80 ${
              currentView === item.id
                ? 'text-[#C0C1FF] bg-[#5B5FFF]/10 border-r-2 border-[#5B5FFF]'
                : 'text-[#908FA1] hover:bg-[#1B1B20] hover:text-[#C0C1FF]'
            }`}
          >
            <span className="material-symbols-outlined" data-weight={currentView === item.id ? 'fill' : undefined}>
              {item.icon}
            </span>
            <span className="font-['Inter'] text-sm tracking-wide font-medium">{item.label}</span>
          </button>
        ))}
      </div>

      {/* Footer Navigation */}
      <div className="px-2 pb-6 pt-4 border-t border-[#454555]/15 space-y-1">
        {footerItems.map((item) => (
          <button
            key={item.label}
            className="w-full flex items-center gap-3 py-2.5 px-4 rounded-lg text-[#908FA1] hover:bg-[#1B1B20] hover:text-[#C0C1FF] transition-all duration-200"
          >
            <span className="material-symbols-outlined text-[20px]">{item.icon}</span>
            <span className="font-['Inter'] text-sm tracking-wide">{item.label}</span>
          </button>
        ))}
      </div>
    </nav>
  );
};

export default SidebarGenpulse;