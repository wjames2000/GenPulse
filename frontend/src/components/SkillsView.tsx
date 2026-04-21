import React, { useState, useEffect } from 'react';
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
  Plus,
  RefreshCw,
  CheckCircle,
  XCircle,
  AlertCircle
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';
import { api, Skill, SkillDetails } from '../services/api';

export default function SkillsView() {
  const [skills, setSkills] = useState<Skill[]>([]);
  const [selectedSkill, setSelectedSkill] = useState<Skill | null>(null);
  const [skillDetails, setSkillDetails] = useState<SkillDetails | null>(null);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    loadSkills();
  }, []);

  useEffect(() => {
    if (selectedSkill) {
      loadSkillDetails(selectedSkill.id);
    }
  }, [selectedSkill]);

  const loadSkills = async () => {
    setLoading(true);
    try {
      const data = await api.getSkills();
      setSkills(data);
      if (data.length > 0 && !selectedSkill) {
        setSelectedSkill(data[0]);
      }
    } catch (error) {
      console.error('Failed to load skills:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadSkillDetails = async (skillId: string) => {
    try {
      const details = await api.getSkillDetails(skillId);
      setSkillDetails(details);
    } catch (error) {
      console.error('Failed to load skill details:', error);
    }
  };

  const handleEnableSkill = async (skillId: string) => {
    try {
      await api.enableSkill(skillId);
      await loadSkills(); // 重新加载技能列表
    } catch (error) {
      console.error('Failed to enable skill:', error);
    }
  };

  const handleDisableSkill = async (skillId: string) => {
    try {
      await api.disableSkill(skillId);
      await loadSkills(); // 重新加载技能列表
    } catch (error) {
      console.error('Failed to disable skill:', error);
    }
  };

  const handleValidateSkill = async (skillId: string) => {
    try {
      await api.validateSkill(skillId);
      await loadSkills(); // 重新加载技能列表
    } catch (error) {
      console.error('Failed to validate skill:', error);
    }
  };

  const filteredSkills = skills.filter(skill =>
    skill.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    skill.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
    skill.tags.some(tag => tag.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  const getSkillIcon = (category: string) => {
    switch (category) {
      case 'frontend':
        return Code2;
      case 'backend':
        return Database;
      case 'devops':
        return GitBranch;
      default:
        return BoxIcon;
    }
  };

  const getComplexityColor = (complexity: string) => {
    switch (complexity) {
      case 'simple':
        return 'text-green-500';
      case 'medium':
        return 'text-yellow-500';
      case 'complex':
        return 'text-red-500';
      default:
        return 'text-outline';
    }
  };

  return (
    <div className="flex-1 overflow-hidden flex flex-col p-8 gap-8 h-full bg-background font-sans">
      {/* Workspace Header */}
      <header className="flex items-end justify-between shrink-0">
        <div>
          <h2 className="text-5xl font-black text-on-surface tracking-tight leading-tight mb-2">技能库</h2>
          <p className="text-outline text-base font-medium">管理和配置代理的认知能力模块。</p>
        </div>
        <div className="flex gap-3">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-outline-variant" size={16} />
            <input
              type="text"
              placeholder="搜索技能..."
              className="pl-10 pr-4 py-2.5 rounded-xl bg-surface-container hover:bg-surface-container-highest text-on-surface-variant text-xs font-bold border border-outline-variant/10 shadow-sm transition-all focus:outline-none focus:ring-2 focus:ring-primary/20"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
          <button 
            className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-surface-container hover:bg-surface-container-highest text-on-surface-variant text-xs font-bold border border-outline-variant/10 shadow-sm transition-all active:scale-95 uppercase tracking-widest"
            onClick={loadSkills}
            disabled={loading}
          >
            <RefreshCw size={16} className={loading ? 'animate-spin' : ''} />
            {loading ? '加载中...' : '刷新'}
          </button>
        </div>
      </header>

      {/* Bento Grid Layout */}
      <div className="flex-1 flex gap-8 overflow-hidden min-h-0">
        {/* Left Pane: Skill List */}
        <div className="w-[420px] shrink-0 flex flex-col gap-5 overflow-y-auto pr-3 pb-8 custom-scrollbar">
          {loading ? (
            <div className="flex items-center justify-center h-32">
              <RefreshCw className="animate-spin text-primary" size={24} />
            </div>
          ) : filteredSkills.length === 0 ? (
            <div className="text-center py-8 text-outline-variant">
              <p>没有找到技能</p>
            </div>
          ) : (
            <>
              {filteredSkills.map((skill, index) => {
                const Icon = getSkillIcon(skill.category);
                const isSelected = selectedSkill?.id === skill.id;
                
                return (
                  <div
                    key={skill.id}
                    className={`rounded-3xl p-6 flex flex-col gap-6 relative cursor-pointer border transition-all group ${
                      isSelected
                        ? 'bg-surface-container-highest/60 border-primary/20 shadow-2xl hover:bg-surface-container-highest'
                        : 'bg-surface-container-low/40 border-outline-variant/5 hover:border-outline-variant/20 hover:bg-surface-container-high'
                    }`}
                    onClick={() => setSelectedSkill(skill)}
                  >
                    {/* Selection Indicator */}
                    {isSelected && (
                      <div className="absolute left-0 top-1/2 -translate-y-1/2 w-1.5 h-12 bg-primary rounded-r-full shadow-[0_0_15px_rgba(91,95,255,0.8)]" />
                    )}
                    
                    <div className="flex justify-between items-start">
                      <div className="flex items-center gap-4">
                        <div className={`w-12 h-12 rounded-2xl flex items-center justify-center shadow-inner border ${
                          isSelected
                            ? 'bg-primary-container/20 text-primary border-primary/20'
                            : 'bg-surface-container-highest/50 text-outline border-outline-variant/10 group-hover:text-primary group-hover:bg-primary/10 group-hover:border-primary/20'
                        }`}>
                          <Icon size={24} />
                        </div>
                        <div>
                          <h3 className="text-lg font-black flex items-center gap-3">
                            <span className={isSelected ? 'text-on-surface' : 'text-on-surface-variant group-hover:text-on-surface'}>
                              {skill.name}
                            </span>
                            {skill.enabled && (
                              <span className="relative flex h-2.5 w-2.5">
                                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-secondary opacity-50" />
                                <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-secondary" />
                              </span>
                            )}
                          </h3>
                          <div className="flex items-center gap-2 mt-1">
                            <span className="text-[10px] text-outline font-mono font-bold uppercase tracking-widest opacity-70">
                              {skill.version} • {skill.category}
                            </span>
                            {skill.validated && (
                              <CheckCircle size={10} className="text-green-500" />
                            )}
                          </div>
                        </div>
                      </div>
                      <div className="flex items-center gap-1">
                        {skill.enabled ? (
                          <button
                            className="text-outline hover:text-red-500 transition-colors p-1.5 rounded-xl hover:bg-white/5"
                            onClick={(e) => {
                              e.stopPropagation();
                              handleDisableSkill(skill.id);
                            }}
                            title="禁用技能"
                          >
                            <XCircle size={16} />
                          </button>
                        ) : (
                          <button
                            className="text-outline hover:text-green-500 transition-colors p-1.5 rounded-xl hover:bg-white/5"
                            onClick={(e) => {
                              e.stopPropagation();
                              handleEnableSkill(skill.id);
                            }}
                            title="启用技能"
                          >
                            <CheckCircle size={16} />
                          </button>
                        )}
                        {!skill.validated && (
                          <button
                            className="text-outline hover:text-yellow-500 transition-colors p-1.5 rounded-xl hover:bg-white/5"
                            onClick={(e) => {
                              e.stopPropagation();
                              handleValidateSkill(skill.id);
                            }}
                            title="验证技能"
                          >
                            <AlertCircle size={16} />
                          </button>
                        )}
                        <button className="text-outline hover:text-primary transition-colors p-1.5 rounded-xl hover:bg-white/5">
                          <MoreVertical size={20} />
                        </button>
                      </div>
                    </div>
                    
                    <p className={`text-sm leading-relaxed line-clamp-2 ${
                      isSelected ? 'text-on-surface-variant opacity-80' : 'text-outline opacity-80 group-hover:opacity-100'
                    }`}>
                      {skill.description}
                    </p>
                    
                    <div className="flex items-center gap-6 pt-5 border-t border-outline-variant/10">
                      <SkillStat 
                        label="成功率" 
                        value={`${(skill.success_rate * 100).toFixed(1)}%`} 
                        color="text-secondary" 
                      />
                      <div className="w-px h-8 bg-outline-variant/10" />
                      <SkillStat 
                        label="调用次数" 
                        value={skill.usage_count.toLocaleString()} 
                      />
                      <div className="w-px h-8 bg-outline-variant/10" />
                      <div className="flex flex-col gap-1">
                        <span className="text-[10px] uppercase font-black tracking-widest text-outline-variant leading-none">复杂度</span>
                        <span className={cn("text-base font-mono font-bold leading-none", getComplexityColor(skill.complexity))}>
                          {skill.complexity}
                        </span>
                      </div>
                    </div>
                    
                    {/* Tags */}
                    {skill.tags.length > 0 && (
                      <div className="flex flex-wrap gap-1">
                        {skill.tags.map((tag, tagIndex) => (
                          <span
                            key={tagIndex}
                            className="px-2 py-1 text-[10px] font-bold uppercase tracking-widest rounded-full bg-surface-container-highest/50 text-outline-variant"
                          >
                            {tag}
                          </span>
                        ))}
                      </div>
                    )}
                  </div>
                );
              })}
            </>
          )}
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
              details.json
            </div>
            <div className="flex-1" />
            <div className="flex items-center gap-2 px-2">
              <button 
                className="text-outline hover:text-primary transition-all p-2 rounded-xl hover:bg-white/5" 
                title="Copy to clipboard"
                onClick={() => {
                  if (skillDetails) {
                    navigator.clipboard.writeText(JSON.stringify(skillDetails, null, 2));
                  }
                }}
              >
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
            {skillDetails ? (
              <pre className="m-0 relative z-10 select-text"><code className="text-on-surface-variant/90 font-medium">
{`name: ${skillDetails.name}
id: ${skillDetails.id}
version: "${skillDetails.version}"
category: ${skillDetails.category}
complexity: ${skillDetails.complexity}
enabled: ${skillDetails.enabled}
validated: ${skillDetails.validated}

description: |
  ${skillDetails.description}

usage_stats:
  usage_count: ${skillDetails.usage_count}
  success_rate: ${(skillDetails.success_rate * 100).toFixed(1)}%

tags:`}
{skillDetails.tags.map(tag => `\n  - ${tag}`).join('')}
{skillDetails.agent_types && skillDetails.agent_types.length > 0 ? `
\nagent_types:` + skillDetails.agent_types.map(type => `\n  - ${type}`).join('') : ''}
{skillDetails.steps && skillDetails.steps.length > 0 ? `
\nsteps:` + skillDetails.steps.map((step: any, i: number) => `
  - step_${i + 1}:
      action: ${step.action}
      tool: ${step.tool}`).join('') : ''}
{skillDetails.examples && skillDetails.examples.length > 0 ? `
\nexamples:` + skillDetails.examples.map((example: string, i: number) => `
  - example_${i + 1}: "${example}"`).join('') : ''}
{skillDetails.tips && skillDetails.tips.length > 0 ? `
\ntips:` + skillDetails.tips.map((tip: string, i: number) => `
  - tip_${i + 1}: "${tip}"`).join('') : ''}
{skillDetails.warnings && skillDetails.warnings.length > 0 ? `
\nwarnings:` + skillDetails.warnings.map((warning: string, i: number) => `
  - warning_${i + 1}: "${warning}"`).join('') : ''}
{skillDetails.prerequisites && skillDetails.prerequisites.length > 0 ? `
\nprerequisites:` + skillDetails.prerequisites.map((prereq: string) => `\n  - ${prereq}`).join('') : ''}
{skillDetails.related_tools && skillDetails.related_tools.length > 0 ? `
\nrelated_tools:` + skillDetails.related_tools.map((tool: string) => `\n  - ${tool}`).join('') : ''}
{`
\ntoken_estimate: ${skillDetails.token_estimate}
avg_execution_time: ${skillDetails.avg_execution_time}
source_task_id: ${skillDetails.source_task_id}`}
              </code></pre>
            ) : (
              <div className="flex items-center justify-center h-full text-outline-variant">
                {selectedSkill ? '加载技能详情中...' : '请选择一个技能查看详情'}
              </div>
            )}
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
