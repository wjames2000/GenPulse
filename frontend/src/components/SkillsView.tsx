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
  AlertCircle,
  Trash2,
  RotateCcw,
  Save,
  FileText,
  Undo2,
  Clock,
  User,
  Tag,
  Download,
  Upload,
  Globe,
  Link,
  Loader2
} from 'lucide-react';
import { motion, AnimatePresence } from 'motion/react';
import { cn } from '../utils';
import { api, Skill, SkillDetails, SkillVersion } from '../services/api';

type ModalMode = 'create' | 'edit' | 'versions' | 'online_import' | 'official_library' | null;

interface OnlineSkillEntry {
  id: string;
  name: string;
  description: string;
  category: string;
  tags: string[];
  complexity: string;
  author: string;
  version: string;
}

export default function SkillsView() {
  const [skills, setSkills] = useState<Skill[]>([]);
  const [selectedSkill, setSelectedSkill] = useState<Skill | null>(null);
  const [skillDetails, setSkillDetails] = useState<SkillDetails | null>(null);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [modalMode, setModalMode] = useState<ModalMode>(null);
  const [versions, setVersions] = useState<SkillVersion[]>([]);
  const [versionsLoading, setVersionsLoading] = useState(false);
  const [formData, setFormData] = useState<Record<string, any>>({});
  const [saving, setSaving] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState<string | null>(null);

  const [importURL, setImportURL] = useState('');
  const [importing, setImporting] = useState(false);
  const [importError, setImportError] = useState<string | null>(null);
  const [importSuccess, setImportSuccess] = useState<string | null>(null);
  const [onlineSources, setOnlineSources] = useState<any[]>([]);
  const [selectedSource, setSelectedSource] = useState('custom_url');
  const [sourceSkillID, setSourceSkillID] = useState('');

  const [officialSkills, setOfficialSkills] = useState<OnlineSkillEntry[]>([]);
  const [officialLoading, setOfficialLoading] = useState(false);
  const [officialSearchQuery, setOfficialSearchQuery] = useState('');
  const [installingSkill, setInstallingSkill] = useState<string | null>(null);

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
      } else if (data.length === 0) {
        setSelectedSkill(null);
        setSkillDetails(null);
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

  const loadVersions = async (skillId: string) => {
    setVersionsLoading(true);
    try {
      const data = await api.getSkillVersions(skillId);
      setVersions(data);
    } catch (error) {
      console.error('Failed to load versions:', error);
      setVersions([]);
    } finally {
      setVersionsLoading(false);
    }
  };

  const handleEnableSkill = async (skillId: string) => {
    try {
      await api.enableSkill(skillId);
      await loadSkills();
    } catch (error) {
      console.error('Failed to enable skill:', error);
    }
  };

  const handleDisableSkill = async (skillId: string) => {
    try {
      await api.disableSkill(skillId);
      await loadSkills();
    } catch (error) {
      console.error('Failed to disable skill:', error);
    }
  };

  const handleValidateSkill = async (skillId: string) => {
    try {
      await api.validateSkill(skillId);
      await loadSkills();
    } catch (error) {
      console.error('Failed to validate skill:', error);
    }
  };

  const handleDeleteSkill = async (skillId: string) => {
    try {
      await api.deleteSkill(skillId);
      setConfirmDelete(null);
      if (selectedSkill?.id === skillId) {
        setSelectedSkill(null);
        setSkillDetails(null);
      }
      await loadSkills();
    } catch (error) {
      console.error('Failed to delete skill:', error);
    }
  };

  const openCreateModal = () => {
    setFormData({
      name: '',
      description: '',
      author: 'user',
      category: 'general',
      complexity: 'medium',
      tags: [],
    });
    setModalMode('create');
  };

  const openOfficialLibrary = async () => {
    setOfficialSkills([]);
    setOfficialSearchQuery('');
    setInstallingSkill(null);
    setModalMode('official_library');
    setOfficialLoading(true);
    try {
      const skills = await api.searchOnlineSkills('', 'official_repo');
      setOfficialSkills(skills || []);
    } catch {
      setOfficialSkills([]);
    } finally {
      setOfficialLoading(false);
    }
  };

  const handleInstallOfficialSkill = async (skillId: string) => {
    setInstallingSkill(skillId);
    try {
      const result = await api.importSkillFromOnline('official_repo', skillId);
      await loadSkills();
      alert(`技能 "${result.name}" 安装成功！`);
    } catch (error: any) {
      alert(error?.message || '安装失败');
    } finally {
      setInstallingSkill(null);
    }
  };

  const openOnlineImportModal = async () => {
    setImportURL('');
    setImportError(null);
    setImportSuccess(null);
    setSourceSkillID('');
    setSelectedSource('custom_url');
    try {
      const sources = await api.listOnlineSources();
      setOnlineSources(sources);
    } catch {
      setOnlineSources([]);
    }
    setModalMode('online_import');
  };

  const handleImportFromURL = async () => {
    if (!importURL.trim()) return;
    setImporting(true);
    setImportError(null);
    setImportSuccess(null);
    try {
      const result = await api.importSkillFromURL(importURL.trim());
      setImportSuccess(`技能 "${result.name}" 导入成功！`);
      setImportURL('');
      await loadSkills();
    } catch (error: any) {
      setImportError(error?.message || '导入失败，请检查URL是否正确');
    } finally {
      setImporting(false);
    }
  };

  const handleImportFromOnline = async () => {
    if (!sourceSkillID.trim() || !selectedSource) return;
    setImporting(true);
    setImportError(null);
    setImportSuccess(null);
    try {
      const result = await api.importSkillFromOnline(selectedSource, sourceSkillID.trim());
      setImportSuccess(`技能 "${result.name}" 导入成功！`);
      setSourceSkillID('');
      await loadSkills();
    } catch (error: any) {
      setImportError(error?.message || '导入失败');
    } finally {
      setImporting(false);
    }
  };

  const openEditModal = () => {
    if (!skillDetails) return;
    setFormData({
      name: skillDetails.name,
      description: skillDetails.description,
      category: skillDetails.category,
      complexity: skillDetails.complexity,
      tags: skillDetails.tags,
      steps: skillDetails.steps || [],
      examples: skillDetails.examples || [],
      tips: skillDetails.tips || [],
      warnings: skillDetails.warnings || [],
      agent_types: skillDetails.agent_types || [],
      token_estimate: skillDetails.token_estimate,
      change_log: '',
    });
    setModalMode('edit');
  };

  const openVersionsModal = async () => {
    if (!selectedSkill) return;
    setModalMode('versions');
    await loadVersions(selectedSkill.id);
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      if (modalMode === 'create') {
        await api.createSkill(formData);
      } else if (modalMode === 'edit' && selectedSkill) {
        await api.updateSkill(selectedSkill.id, formData);
      }
      setModalMode(null);
      await loadSkills();
      if (selectedSkill) {
        await loadSkillDetails(selectedSkill.id);
      }
    } catch (error) {
      console.error('Failed to save skill:', error);
    } finally {
      setSaving(false);
    }
  };

  const handleRollback = async (versionId: string) => {
    if (!selectedSkill) return;
    try {
      await api.rollbackSkill(selectedSkill.id, versionId);
      setModalMode(null);
      await loadSkills();
      await loadSkillDetails(selectedSkill.id);
    } catch (error) {
      console.error('Failed to rollback skill:', error);
    }
  };

  const filteredSkills = skills.filter(skill =>
    skill.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    skill.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
    skill.tags.some(tag => tag.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  const getSkillIcon = (category: string) => {
    switch (category) {
      case 'frontend': return Code2;
      case 'backend': return Database;
      case 'devops': return GitBranch;
      default: return BoxIcon;
    }
  };

  const getComplexityColor = (complexity: string) => {
    switch (complexity) {
      case 'simple': return 'text-green-500';
      case 'medium': return 'text-yellow-500';
      case 'complex': return 'text-red-500';
      default: return 'text-outline';
    }
  };

  return (
    <div className="flex-1 overflow-hidden flex flex-col p-8 gap-8 h-full bg-background font-sans">
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
            className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-surface-container hover:bg-surface-container-highest text-on-surface-variant text-xs font-bold border border-outline-variant/10 shadow-sm transition-all active:scale-95"
            onClick={openOfficialLibrary}
            title="官方技能库"
          >
            <Database size={16} />
            官方技能库
          </button>
          <button 
            className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-surface-container hover:bg-surface-container-highest text-on-surface-variant text-xs font-bold border border-outline-variant/10 shadow-sm transition-all active:scale-95"
            onClick={openOnlineImportModal}
            title="在线导入技能"
          >
            <Globe size={16} />
            在线导入
          </button>
          <button 
            className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-primary text-black font-bold text-xs uppercase tracking-widest shadow-sm transition-all active:scale-95 hover:brightness-110"
            onClick={openCreateModal}
          >
            <Plus size={16} />
            新建技能
          </button>
          <button 
            className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-surface-container hover:bg-surface-container-highest text-on-surface-variant text-xs font-bold border border-outline-variant/10 shadow-sm transition-all active:scale-95"
            onClick={loadSkills}
            disabled={loading}
          >
            <RefreshCw size={16} className={loading ? 'animate-spin' : ''} />
          </button>
        </div>
      </header>

      <div className="flex-1 flex gap-8 overflow-hidden min-h-0">
        <div className="w-[420px] shrink-0 flex flex-col gap-5 overflow-y-auto pr-3 pb-8 custom-scrollbar">
          {loading ? (
            <div className="flex items-center justify-center h-32">
              <RefreshCw className="animate-spin text-primary" size={24} />
            </div>
          ) : filteredSkills.length === 0 ? (
            <div className="text-center py-8 text-outline-variant">
              <p>没有找到技能</p>
              {skills.length === 0 && (
                <button 
                  className="mt-4 text-primary text-sm font-bold hover:underline"
                  onClick={openCreateModal}
                >
                  创建第一个技能
                </button>
              )}
            </div>
          ) : (
            filteredSkills.map((skill) => {
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
                          onClick={(e) => { e.stopPropagation(); handleDisableSkill(skill.id); }}
                          title="禁用技能"
                        >
                          <XCircle size={16} />
                        </button>
                      ) : (
                        <button
                          className="text-outline hover:text-green-500 transition-colors p-1.5 rounded-xl hover:bg-white/5"
                          onClick={(e) => { e.stopPropagation(); handleEnableSkill(skill.id); }}
                          title="启用技能"
                        >
                          <CheckCircle size={16} />
                        </button>
                      )}
                      {!skill.validated && (
                        <button
                          className="text-outline hover:text-yellow-500 transition-colors p-1.5 rounded-xl hover:bg-white/5"
                          onClick={(e) => { e.stopPropagation(); handleValidateSkill(skill.id); }}
                          title="验证技能"
                        >
                          <AlertCircle size={16} />
                        </button>
                      )}
                      <button 
                        className="text-outline hover:text-red-500 transition-colors p-1.5 rounded-xl hover:bg-white/5"
                        onClick={(e) => { e.stopPropagation(); setConfirmDelete(skill.id); }}
                        title="删除技能"
                      >
                        <Trash2 size={16} />
                      </button>
                    </div>
                  </div>
                  
                  <p className={`text-sm leading-relaxed line-clamp-2 ${
                    isSelected ? 'text-on-surface-variant opacity-80' : 'text-outline opacity-80 group-hover:opacity-100'
                  }`}>
                    {skill.description}
                  </p>
                  
                  <div className="flex items-center gap-6 pt-5 border-t border-outline-variant/10">
                    <SkillStat label="成功率" value={`${(skill.success_rate * 100).toFixed(1)}%`} color="text-secondary" />
                    <div className="w-px h-8 bg-outline-variant/10" />
                    <SkillStat label="调用次数" value={skill.usage_count.toLocaleString()} />
                    <div className="w-px h-8 bg-outline-variant/10" />
                    <div className="flex flex-col gap-1">
                      <span className="text-[10px] uppercase font-black tracking-widest text-outline-variant leading-none">复杂度</span>
                      <span className={cn("text-base font-mono font-bold leading-none", getComplexityColor(skill.complexity))}>
                        {skill.complexity}
                      </span>
                    </div>
                  </div>
                  
                  {skill.tags.length > 0 && (
                    <div className="flex flex-wrap gap-1">
                      {skill.tags.map((tag, tagIndex) => (
                        <span key={tagIndex} className="px-2 py-1 text-[10px] font-bold uppercase tracking-widest rounded-full bg-surface-container-highest/50 text-outline-variant">
                          {tag}
                        </span>
                      ))}
                    </div>
                  )}

                  {confirmDelete === skill.id && (
                    <div className="flex items-center gap-2 pt-2 border-t border-red-500/20">
                      <span className="text-xs text-red-500 font-bold">确认删除?</span>
                      <button
                        className="px-3 py-1 text-xs font-bold bg-red-500 text-white rounded-lg hover:bg-red-600"
                        onClick={(e) => { e.stopPropagation(); handleDeleteSkill(skill.id); }}
                      >
                        确认
                      </button>
                      <button
                        className="px-3 py-1 text-xs font-bold bg-surface-container-highest text-on-surface rounded-lg hover:bg-surface-container-high"
                        onClick={(e) => { e.stopPropagation(); setConfirmDelete(null); }}
                      >
                        取消
                      </button>
                    </div>
                  )}
                </div>
              );
            })
          )}
        </div>

        <div className="flex-1 flex flex-col bg-surface-container-lowest/40 rounded-3xl border border-outline-variant/10 overflow-hidden shadow-2xl backdrop-blur-xl">
          <div className="bg-surface-container-highest/50 border-b border-outline-variant/10 flex items-center px-4 py-2 gap-2 shrink-0">
            <div className="flex items-center gap-2.5 px-5 py-2.5 bg-surface-container-low rounded-xl text-primary text-xs font-bold tracking-widest border-t border-l border-r border-outline-variant/10 shadow-lg">
              <FileText size={14} className="opacity-70" />
              definition.yaml
            </div>
            <div className="flex-1" />
            {selectedSkill && (
              <>
                <button 
                  className="text-outline hover:text-primary transition-all p-2 rounded-xl hover:bg-white/5" 
                  title="版本历史"
                  onClick={openVersionsModal}
                >
                  <History size={16} />
                </button>
                <button 
                  className="text-outline hover:text-primary transition-all p-2 rounded-xl hover:bg-white/5" 
                  title="编辑技能"
                  onClick={openEditModal}
                >
                  <Edit size={16} />
                </button>
                <button 
                  className="text-outline hover:text-primary transition-all p-2 rounded-xl hover:bg-white/5" 
                  title="复制到剪贴板"
                  onClick={() => {
                    if (skillDetails) {
                      navigator.clipboard.writeText(JSON.stringify(skillDetails, null, 2));
                    }
                  }}
                >
                  <Copy size={16} />
                </button>
              </>
            )}
          </div>

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
                {skills.length === 0 ? '技能库为空，点击"新建技能"开始' : '请选择一个技能查看详情'}
              </div>
            )}
          </div>
        </div>
      </div>

      <AnimatePresence>
        {(modalMode === 'create' || modalMode === 'edit') && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center"
            onClick={() => setModalMode(null)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              className="bg-surface-container rounded-3xl border border-outline-variant/10 shadow-2xl w-[600px] max-h-[80vh] overflow-y-auto"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="p-6 border-b border-outline-variant/10 flex items-center justify-between">
                <h3 className="text-xl font-black">
                  {modalMode === 'create' ? '新建技能' : '编辑技能'}
                </h3>
                <button 
                  className="text-outline hover:text-on-surface p-2 rounded-xl hover:bg-white/5"
                  onClick={() => setModalMode(null)}
                >
                  <XCircle size={20} />
                </button>
              </div>

              <div className="p-6 space-y-5">
                <FormField label="名称" required>
                  <input
                    type="text"
                    className="w-full px-4 py-3 rounded-xl bg-surface-container-highest text-on-surface border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20 font-medium"
                    placeholder="技能名称"
                    value={formData.name || ''}
                    onChange={(e) => setFormData({...formData, name: e.target.value})}
                  />
                </FormField>

                <FormField label="描述" required>
                  <textarea
                    className="w-full px-4 py-3 rounded-xl bg-surface-container-highest text-on-surface border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20 font-medium resize-none"
                    placeholder="技能描述"
                    rows={3}
                    value={formData.description || ''}
                    onChange={(e) => setFormData({...formData, description: e.target.value})}
                  />
                </FormField>

                <div className="grid grid-cols-2 gap-4">
                  <FormField label="分类">
                    <select
                      className="w-full px-4 py-3 rounded-xl bg-surface-container-highest text-on-surface border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20 font-medium"
                      value={formData.category || 'general'}
                      onChange={(e) => setFormData({...formData, category: e.target.value})}
                    >
                      <option value="general">通用</option>
                      <option value="frontend">前端</option>
                      <option value="backend">后端</option>
                      <option value="devops">DevOps</option>
                      <option value="testing">测试</option>
                      <option value="database">数据库</option>
                      <option value="security">安全</option>
                    </select>
                  </FormField>

                  <FormField label="复杂度">
                    <select
                      className="w-full px-4 py-3 rounded-xl bg-surface-container-highest text-on-surface border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20 font-medium"
                      value={formData.complexity || 'medium'}
                      onChange={(e) => setFormData({...formData, complexity: e.target.value})}
                    >
                      <option value="simple">简单</option>
                      <option value="medium">中等</option>
                      <option value="complex">复杂</option>
                    </select>
                  </FormField>
                </div>

                <FormField label="标签 (逗号分隔)">
                  <input
                    type="text"
                    className="w-full px-4 py-3 rounded-xl bg-surface-container-highest text-on-surface border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20 font-medium"
                    placeholder="react, frontend, ui"
                    value={(formData.tags || []).join(', ')}
                    onChange={(e) => setFormData({...formData, tags: e.target.value.split(',').map((t: string) => t.trim()).filter(Boolean)})}
                  />
                </FormField>

                <FormField label="适用Agent类型 (逗号分隔)">
                  <input
                    type="text"
                    className="w-full px-4 py-3 rounded-xl bg-surface-container-highest text-on-surface border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20 font-medium"
                    placeholder="frontend, fullstack, backend"
                    value={(formData.agent_types || []).join(', ')}
                    onChange={(e) => setFormData({...formData, agent_types: e.target.value.split(',').map((t: string) => t.trim()).filter(Boolean)})}
                  />
                </FormField>

                <FormField label="Token估算">
                  <input
                    type="number"
                    className="w-full px-4 py-3 rounded-xl bg-surface-container-highest text-on-surface border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20 font-medium"
                    placeholder="1000"
                    value={formData.token_estimate || 1000}
                    onChange={(e) => setFormData({...formData, token_estimate: parseInt(e.target.value) || 0})}
                  />
                </FormField>

                {modalMode === 'edit' && (
                  <FormField label="变更说明">
                    <textarea
                      className="w-full px-4 py-3 rounded-xl bg-surface-container-highest text-on-surface border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20 font-medium resize-none"
                      placeholder="描述此次修改的内容"
                      rows={2}
                      value={formData.change_log || ''}
                      onChange={(e) => setFormData({...formData, change_log: e.target.value})}
                    />
                  </FormField>
                )}
              </div>

              <div className="p-6 border-t border-outline-variant/10 flex justify-end gap-3">
                <button
                  className="px-6 py-3 rounded-xl bg-surface-container-highest text-on-surface font-bold text-sm hover:bg-surface-container-high transition-all"
                  onClick={() => setModalMode(null)}
                >
                  取消
                </button>
                <button
                  className="px-6 py-3 rounded-xl bg-primary text-black font-bold text-sm flex items-center gap-2 hover:brightness-110 transition-all"
                  onClick={handleSave}
                  disabled={saving || !formData.name}
                >
                  <Save size={16} />
                  {saving ? '保存中...' : '保存'}
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}

        {modalMode === 'online_import' && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center"
            onClick={() => setModalMode(null)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              className="bg-surface-container rounded-3xl border border-outline-variant/10 shadow-2xl w-[540px] max-h-[80vh] overflow-y-auto"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="p-6 border-b border-outline-variant/10 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <Globe size={20} className="text-primary" />
                  <h3 className="text-xl font-black">在线导入技能</h3>
                </div>
                <button 
                  className="text-outline hover:text-on-surface p-2 rounded-xl hover:bg-white/5"
                  onClick={() => setModalMode(null)}
                >
                  <XCircle size={20} />
                </button>
              </div>

              <div className="p-6 space-y-6">
                <div>
                  <FormField label="来源类型">
                    <select
                      className="w-full px-4 py-3 rounded-xl bg-surface-container-highest text-on-surface border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20 font-medium"
                      value={selectedSource}
                      onChange={(e) => setSelectedSource(e.target.value)}
                    >
                      <option value="custom_url">自定义URL</option>
                      <option value="github">GitHub仓库</option>
                      <option value="marketplace">技能市场</option>
                      {onlineSources.map((source: any) => (
                        <option key={source.id} value={source.id}>
                          {source.name}
                        </option>
                      ))}
                    </select>
                  </FormField>
                </div>

                {selectedSource === 'custom_url' && (
                  <div>
                    <FormField label="技能文件URL" required>
                      <input
                        type="url"
                        className="w-full px-4 py-3 rounded-xl bg-surface-container-highest text-on-surface border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20 font-medium"
                        placeholder="https://example.com/skill.yaml"
                        value={importURL}
                        onChange={(e) => setImportURL(e.target.value)}
                      />
                    </FormField>
                    <p className="mt-2 text-xs text-outline flex items-center gap-1">
                      <Link size={12} />
                      支持YAML/JSON格式的技能定义文件URL，或GitHub raw URL
                    </p>
                  </div>
                )}

                {selectedSource === 'github' && (
                  <div>
                    <FormField label="GitHub URL" required>
                      <input
                        type="url"
                        className="w-full px-4 py-3 rounded-xl bg-surface-container-highest text-on-surface border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20 font-medium"
                        placeholder="https://github.com/user/repo 或 https://github.com/user/repo/blob/main/skill.yaml"
                        value={importURL}
                        onChange={(e) => setImportURL(e.target.value)}
                      />
                    </FormField>
                    <p className="mt-2 text-xs text-outline flex items-center gap-1">
                      <Link size={12} />
                      支持GitHub仓库或raw文件URL
                    </p>
                  </div>
                )}

                {selectedSource !== 'custom_url' && selectedSource !== 'github' && (
                  <div>
                    <FormField label="技能ID" required>
                      <input
                        type="text"
                        className="w-full px-4 py-3 rounded-xl bg-surface-container-highest text-on-surface border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20 font-medium"
                        placeholder="输入技能ID或名称"
                        value={sourceSkillID}
                        onChange={(e) => setSourceSkillID(e.target.value)}
                      />
                    </FormField>
                  </div>
                )}

                {importError && (
                  <div className="p-4 rounded-2xl bg-red-500/10 border border-red-500/20 flex items-start gap-3">
                    <XCircle size={18} className="text-red-500 shrink-0 mt-0.5" />
                    <p className="text-sm text-red-500 font-medium">{importError}</p>
                  </div>
                )}

                {importSuccess && (
                  <div className="p-4 rounded-2xl bg-green-500/10 border border-green-500/20 flex items-start gap-3">
                    <CheckCircle size={18} className="text-green-500 shrink-0 mt-0.5" />
                    <p className="text-sm text-green-500 font-medium">{importSuccess}</p>
                  </div>
                )}
              </div>

              <div className="p-6 border-t border-outline-variant/10 flex justify-end gap-3">
                <button
                  className="px-6 py-3 rounded-xl bg-surface-container-highest text-on-surface font-bold text-sm hover:bg-surface-container-high transition-all"
                  onClick={() => setModalMode(null)}
                >
                  取消
                </button>
                <button
                  className="px-6 py-3 rounded-xl bg-primary text-black font-bold text-sm flex items-center gap-2 hover:brightness-110 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                  onClick={selectedSource === 'custom_url' || selectedSource === 'github' ? handleImportFromURL : handleImportFromOnline}
                  disabled={importing || (selectedSource === 'custom_url' || selectedSource === 'github' ? !importURL.trim() : !sourceSkillID.trim())}
                >
                  {importing ? (
                    <Loader2 size={16} className="animate-spin" />
                  ) : (
                    <Download size={16} />
                  )}
                  {importing ? '导入中...' : '导入技能'}
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}

        {modalMode === 'versions' && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center"
            onClick={() => setModalMode(null)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              className="bg-surface-container rounded-3xl border border-outline-variant/10 shadow-2xl w-[500px] max-h-[70vh] overflow-y-auto"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="p-6 border-b border-outline-variant/10 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <History size={20} className="text-primary" />
                  <h3 className="text-xl font-black">版本历史</h3>
                  {selectedSkill && (
                    <span className="text-sm text-outline font-mono">{selectedSkill.name} v{selectedSkill.version}</span>
                  )}
                </div>
                <button 
                  className="text-outline hover:text-on-surface p-2 rounded-xl hover:bg-white/5"
                  onClick={() => setModalMode(null)}
                >
                  <XCircle size={20} />
                </button>
              </div>

              <div className="p-6">
                {versionsLoading ? (
                  <div className="flex items-center justify-center py-8">
                    <RefreshCw className="animate-spin text-primary" size={24} />
                  </div>
                ) : versions.length === 0 ? (
                  <div className="text-center py-8 text-outline-variant">
                    <Clock size={32} className="mx-auto mb-3 opacity-50" />
                    <p className="font-bold">暂无版本历史</p>
                    <p className="text-sm mt-1">编辑技能时将自动保存版本快照</p>
                  </div>
                ) : (
                  <div className="space-y-3">
                    {versions.map((version, index) => (
                      <div key={version.id} className="p-4 rounded-2xl bg-surface-container-highest/50 border border-outline-variant/10">
                        <div className="flex items-center justify-between mb-2">
                          <div className="flex items-center gap-2">
                            <Tag size={14} className="text-primary" />
                            <span className="font-bold font-mono text-sm">{version.version}</span>
                            {index === versions.length - 1 && (
                              <span className="text-[10px] px-2 py-0.5 rounded-full bg-primary/20 text-primary font-bold">当前</span>
                            )}
                          </div>
                          <button
                            className="px-3 py-1.5 rounded-lg bg-surface-container text-xs font-bold flex items-center gap-1 hover:bg-primary hover:text-black transition-all"
                            onClick={() => handleRollback(version.id)}
                          >
                            <Undo2 size={12} />
                            回滚
                          </button>
                        </div>
                        <div className="flex items-center gap-3 text-xs text-outline">
                          <span className="flex items-center gap-1">
                            <Clock size={12} />
                            {new Date(version.created_at).toLocaleString('zh-CN')}
                          </span>
                          <span className="flex items-center gap-1">
                            <User size={12} />
                            {version.created_by}
                          </span>
                        </div>
                        {version.change_log && (
                          <p className="mt-2 text-sm text-on-surface-variant/80">{version.change_log}</p>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </motion.div>
          </motion.div>
        )}

        {modalMode === 'official_library' && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center"
            onClick={() => setModalMode(null)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              className="bg-surface-container rounded-3xl border border-outline-variant/10 shadow-2xl w-[680px] max-h-[85vh] overflow-hidden flex flex-col"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="p-6 border-b border-outline-variant/10 flex items-center justify-between shrink-0">
                <div className="flex items-center gap-3">
                  <Database size={20} className="text-primary" />
                  <h3 className="text-xl font-black">官方技能库</h3>
                  <span className="text-xs px-2.5 py-1 rounded-full bg-primary/15 text-primary font-bold tracking-wider">genpulse-skills</span>
                </div>
                <button 
                  className="text-outline hover:text-on-surface p-2 rounded-xl hover:bg-white/5"
                  onClick={() => setModalMode(null)}
                >
                  <XCircle size={20} />
                </button>
              </div>

              <div className="p-4 border-b border-outline-variant/10 shrink-0">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-outline-variant" size={16} />
                  <input
                    type="text"
                    placeholder="搜索官方技能..."
                    className="w-full pl-10 pr-4 py-2.5 rounded-xl bg-surface-container-highest text-on-surface-variant text-xs font-bold border border-outline-variant/10 focus:outline-none focus:ring-2 focus:ring-primary/20"
                    value={officialSearchQuery}
                    onChange={(e) => setOfficialSearchQuery(e.target.value)}
                  />
                </div>
              </div>

              <div className="flex-1 overflow-y-auto p-4 custom-scrollbar">
                {officialLoading ? (
                  <div className="flex items-center justify-center py-16">
                    <RefreshCw className="animate-spin text-primary" size={28} />
                  </div>
                ) : officialSkills.length === 0 ? (
                  <div className="text-center py-16 text-outline-variant">
                    <Database size={40} className="mx-auto mb-4 opacity-40" />
                    <p className="font-bold text-lg">暂无可用技能</p>
                    <p className="text-sm mt-2">官方技能库正在建设中，请稍后再试</p>
                  </div>
                ) : (
                  <div className="space-y-3">
                    {officialSkills
                      .filter(skill =>
                        !officialSearchQuery ||
                        skill.name.toLowerCase().includes(officialSearchQuery.toLowerCase()) ||
                        skill.description.toLowerCase().includes(officialSearchQuery.toLowerCase()) ||
                        skill.tags.some(tag => tag.toLowerCase().includes(officialSearchQuery.toLowerCase()))
                      )
                      .map((skill) => (
                        <div key={skill.id} className="p-5 rounded-2xl bg-surface-container-highest/40 border border-outline-variant/10 hover:border-primary/20 transition-all group">
                          <div className="flex items-start justify-between gap-4">
                            <div className="flex-1 min-w-0">
                              <div className="flex items-center gap-3 mb-1">
                                <h4 className="font-black text-on-surface truncate">{skill.name}</h4>
                                <span className="text-[10px] font-mono font-bold text-outline-variant px-2 py-0.5 rounded-full bg-surface-container-lowest/50 shrink-0">
                                  v{skill.version}
                                </span>
                              </div>
                              <p className="text-sm text-outline line-clamp-2 leading-relaxed">{skill.description}</p>
                              <div className="flex items-center gap-3 mt-3">
                                <span className="text-[10px] uppercase font-black tracking-widest text-outline-variant px-2 py-0.5 rounded-full bg-surface-container-lowest/50">
                                  {skill.category}
                                </span>
                                <span className="text-[10px] font-mono font-bold text-outline-variant">
                                  {skill.complexity}
                                </span>
                                <span className="text-[10px] text-outline-variant">
                                  by {skill.author}
                                </span>
                              </div>
                              {skill.tags.length > 0 && (
                                <div className="flex flex-wrap gap-1 mt-3">
                                  {skill.tags.map((tag, i) => (
                                    <span key={i} className="px-2 py-0.5 text-[9px] font-bold uppercase tracking-widest rounded-full bg-primary/8 text-primary/70">
                                      {tag}
                                    </span>
                                  ))}
                                </div>
                              )}
                            </div>
                            <button
                              className="shrink-0 px-4 py-2.5 rounded-xl bg-primary text-black font-bold text-xs flex items-center gap-1.5 hover:brightness-110 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                              onClick={() => handleInstallOfficialSkill(skill.id)}
                              disabled={installingSkill === skill.id}
                            >
                              {installingSkill === skill.id ? (
                                <Loader2 size={14} className="animate-spin" />
                              ) : (
                                <Download size={14} />
                              )}
                              {installingSkill === skill.id ? '安装中...' : '安装'}
                            </button>
                          </div>
                        </div>
                      ))}
                  </div>
                )}
              </div>

              <div className="p-4 border-t border-outline-variant/10 flex items-center justify-between shrink-0">
                <p className="text-xs text-outline-variant">
                  共 {officialSkills.filter(skill =>
                    !officialSearchQuery ||
                    skill.name.toLowerCase().includes(officialSearchQuery.toLowerCase()) ||
                    skill.description.toLowerCase().includes(officialSearchQuery.toLowerCase()) ||
                    skill.tags.some(tag => tag.toLowerCase().includes(officialSearchQuery.toLowerCase()))
                  ).length } 个技能
                </p>
                <button
                  className="px-5 py-2 rounded-xl bg-surface-container-highest text-on-surface font-bold text-xs hover:bg-surface-container-high transition-all"
                  onClick={() => setModalMode(null)}
                >
                  关闭
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
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

function FormField({ label, children, required }: { label: string; children: React.ReactNode; required?: boolean }) {
  return (
    <div>
      <label className="block text-xs font-bold uppercase tracking-widest text-outline mb-2">
        {label}{required && <span className="text-red-500 ml-1">*</span>}
      </label>
      {children}
    </div>
  );
}
