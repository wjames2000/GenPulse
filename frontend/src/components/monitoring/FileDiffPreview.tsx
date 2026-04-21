import React, { useState, useMemo } from 'react';
import { 
  FileText, 
  GitBranch, 
  Folder, 
  Code, 
  Copy, 
  Download, 
  Eye, 
  EyeOff,
  ChevronUp,
  ChevronDown,
  Filter,
  Search,
  RefreshCw,
  CheckCircle2,
  XCircle,
  Plus,
  Minus,
  Edit,
  Trash2,
  Save,
  FolderOpen,
  FilePlus,
  FileMinus,
  FileEdit,
  GitCommit,
  Clock,
  User,
  Hash,
  Maximize2,
  Minimize2,
  Settings,
  MoreVertical,
  ExternalLink,
  FolderTree,
  FileCode,
  FileJson,
  FileType,
  FileImage,
  FileArchive,
  FileVideo,
  FileAudio,
  FileSpreadsheet,
  File,
  FolderPlus,
  FolderMinus
} from 'lucide-react';
import { cn } from '../../utils';
import { FileDiff } from '../../types';

interface FileDiffPreviewProps {
  diffs: FileDiff[];
}

export default function FileDiffPreview({ diffs }: FileDiffPreviewProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const [fileTypeFilter, setFileTypeFilter] = useState<string>('all');
  const [changeTypeFilter, setChangeTypeFilter] = useState<'all' | 'added' | 'modified' | 'deleted'>('all');
  const [expandedFiles, setExpandedFiles] = useState<Set<string>>(new Set());
  const [viewMode, setViewMode] = useState<'unified' | 'split'>('unified');
  const [showLineNumbers, setShowLineNumbers] = useState(true);
  const [showWhitespace, setShowWhitespace] = useState(false);
  const [selectedFile, setSelectedFile] = useState<string | null>(null);

  // 获取所有文件类型
  const fileTypes = useMemo(() => {
    const types = new Set(diffs.map(diff => {
      const ext = diff.filePath.split('.').pop()?.toLowerCase();
      return ext || 'unknown';
    }));
    return Array.from(types);
  }, [diffs]);

  // 过滤文件差异
  const filteredDiffs = useMemo(() => {
    return diffs.filter(diff => {
      // 搜索过滤
      if (searchQuery && !diff.filePath.toLowerCase().includes(searchQuery.toLowerCase())) {
        return false;
      }
      
      // 文件类型过滤
      if (fileTypeFilter !== 'all') {
        const ext = diff.filePath.split('.').pop()?.toLowerCase();
        if (ext !== fileTypeFilter) {
          return false;
        }
      }
      
      // 变更类型过滤
      if (changeTypeFilter !== 'all' && diff.changeType !== changeTypeFilter) {
        return false;
      }
      
      return true;
    });
  }, [diffs, searchQuery, fileTypeFilter, changeTypeFilter]);

  const selectedDiff = selectedFile ? diffs.find(d => d.id === selectedFile) : null;

  const getFileIcon = (filePath: string) => {
    const ext = filePath.split('.').pop()?.toLowerCase();
    
    switch (ext) {
      case 'go': return Code;
      case 'ts': case 'tsx': case 'js': case 'jsx': return FileCode;
      case 'json': return FileJson;
      case 'md': case 'txt': return FileText;
      case 'yml': case 'yaml': return Settings;
      case 'png': case 'jpg': case 'jpeg': case 'gif': case 'svg': return FileImage;
      case 'zip': case 'tar': case 'gz': return FileArchive;
      case 'mp4': case 'avi': case 'mov': return FileVideo;
      case 'mp3': case 'wav': return FileAudio;
      case 'csv': case 'xlsx': return FileSpreadsheet;
      default: return File;
    }
  };

  const getChangeTypeColor = (changeType: string) => {
    switch (changeType) {
      case 'added': return 'text-green-500';
      case 'modified': return 'text-yellow-500';
      case 'deleted': return 'text-red-500';
      default: return 'text-white/60';
    }
  };

  const getChangeTypeBgColor = (changeType: string) => {
    switch (changeType) {
      case 'added': return 'bg-green-500/10 border-green-500/20';
      case 'modified': return 'bg-yellow-500/10 border-yellow-500/20';
      case 'deleted': return 'bg-red-500/10 border-red-500/20';
      default: return 'bg-white/5 border-white/10';
    }
  };

  const getChangeTypeIcon = (changeType: string) => {
    switch (changeType) {
      case 'added': return Plus;
      case 'modified': return Edit;
      case 'deleted': return Trash2;
      default: return File;
    }
  };

  const getChangeTypeLabel = (changeType: string) => {
    switch (changeType) {
      case 'added': return 'Added';
      case 'modified': return 'Modified';
      case 'deleted': return 'Deleted';
      default: return 'Unknown';
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  const toggleExpandFile = (fileId: string) => {
    const newExpanded = new Set(expandedFiles);
    if (newExpanded.has(fileId)) {
      newExpanded.delete(fileId);
    } else {
      newExpanded.add(fileId);
    }
    setExpandedFiles(newExpanded);
  };

  const handleDownloadDiff = (diff: FileDiff) => {
    const content = diff.diff || '';
    const blob = new Blob([content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${diff.filePath.split('/').pop()}.diff`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const handleDownloadAll = () => {
    const content = filteredDiffs.map(diff => 
      `=== ${diff.filePath} (${diff.changeType}) ===\n${diff.diff || ''}\n\n`
    ).join('\n');
    
    const blob = new Blob([content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `file-diffs-${new Date().toISOString().split('T')[0]}.txt`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const renderDiffLines = (diff: string) => {
    const lines = diff.split('\n');
    
    return lines.map((line, index) => {
      let bgColor = '';
      let textColor = '';
      let prefix = '';
      
      if (line.startsWith('+')) {
        bgColor = 'bg-green-500/10';
        textColor = 'text-green-500';
        prefix = '+';
      } else if (line.startsWith('-')) {
        bgColor = 'bg-red-500/10';
        textColor = 'text-red-500';
        prefix = '-';
      } else if (line.startsWith('@@')) {
        bgColor = 'bg-blue-500/10';
        textColor = 'text-blue-500';
        prefix = '@@';
      } else {
        bgColor = 'bg-transparent';
        textColor = 'text-white/60';
        prefix = ' ';
      }
      
      return (
        <div key={index} className={cn("flex font-mono text-sm", bgColor)}>
          {showLineNumbers && (
            <div className="w-12 text-right pr-3 text-white/40 select-none border-r border-white/10">
              {index + 1}
            </div>
          )}
          <div className="w-8 text-center select-none border-r border-white/10">
            <span className={textColor}>{prefix}</span>
          </div>
          <div className={cn("flex-1 px-3 py-1", textColor)}>
            {showWhitespace ? line : line.replace(/\s+$/g, '')}
          </div>
        </div>
      );
    });
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-3">
            <FileText size={24} />
            File Changes Diff Preview
          </h2>
          <p className="text-sm text-white/60 mt-1">
            Visual diff of all files created, modified, or deleted by agents
          </p>
        </div>
        
        <div className="flex items-center gap-3">
          <div className="text-sm text-white/40">
            {filteredDiffs.length} files • {diffs.reduce((sum, diff) => sum + (diff.linesAdded || 0), 0)}+ {diffs.reduce((sum, diff) => sum + (diff.linesDeleted || 0), 0)}-
          </div>
        </div>
      </div>

      {/* Controls */}
      <div className="bg-white/5 border border-white/10 rounded-xl p-4">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="relative">
              <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-white/40" />
              <input
                type="text"
                placeholder="Search files..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 pr-4 py-2 bg-white/5 border border-white/10 rounded-lg text-sm focus:outline-none focus:border-primary transition-colors"
              />
            </div>
            
            <select
              value={fileTypeFilter}
              onChange={(e) => setFileTypeFilter(e.target.value)}
              className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
            >
              <option value="all">All File Types</option>
              {fileTypes.map(type => (
                <option key={type} value={type}>.{type}</option>
              ))}
            </select>
            
            <select
              value={changeTypeFilter}
              onChange={(e) => setChangeTypeFilter(e.target.value as any)}
              className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary transition-colors"
            >
              <option value="all">All Changes</option>
              <option value="added">Added</option>
              <option value="modified">Modified</option>
              <option value="deleted">Deleted</option>
            </select>
          </div>
          
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2">
              <span className="text-sm text-white/60">View:</span>
              <button
                onClick={() => setViewMode('unified')}
                className={cn(
                  "px-3 py-1 text-xs rounded transition-colors",
                  viewMode === 'unified' 
                    ? "bg-primary/20 text-primary" 
                    : "bg-white/5 text-white/60 hover:bg-white/10"
                )}
              >
                Unified
              </button>
              <button
                onClick={() => setViewMode('split')}
                className={cn(
                  "px-3 py-1 text-xs rounded transition-colors",
                  viewMode === 'split' 
                    ? "bg-primary/20 text-primary" 
                    : "bg-white/5 text-white/60 hover:bg-white/10"
                )}
              >
                Split
              </button>
            </div>
            
            <button
              onClick={() => setShowLineNumbers(!showLineNumbers)}
              className={cn(
                "p-2 rounded-lg transition-colors",
                showLineNumbers
                  ? "bg-primary/20 text-primary" 
                  : "bg-white/5 text-white/60 hover:bg-white/10"
              )}
              title="Toggle line numbers"
            >
              <Hash size={20} />
            </button>
            
            <button
              onClick={() => setShowWhitespace(!showWhitespace)}
              className={cn(
                "p-2 rounded-lg transition-colors",
                showWhitespace
                  ? "bg-primary/20 text-primary" 
                  : "bg-white/5 text-white/60 hover:bg-white/10"
              )}
              title="Toggle whitespace"
            >
              {showWhitespace ? <Eye size={20} /> : <EyeOff size={20} />}
            </button>
            
            <button
              onClick={handleDownloadAll}
              className="p-2 rounded-lg bg-white/5 text-white/60 hover:bg-white/10 transition-colors"
              title="Download all diffs"
            >
              <Download size={20} />
            </button>
          </div>
        </div>
        
        {/* Statistics */}
        <div className="grid grid-cols-4 gap-4 mt-4">
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Total Files</div>
            <div className="text-lg font-bold">{diffs.length}</div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Lines Added</div>
            <div className="text-lg font-bold text-green-500">
              +{diffs.reduce((sum, diff) => sum + (diff.linesAdded || 0), 0)}
            </div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Lines Deleted</div>
            <div className="text-lg font-bold text-red-500">
              -{diffs.reduce((sum, diff) => sum + (diff.linesDeleted || 0), 0)}
            </div>
          </div>
          
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="text-xs text-white/40 mb-1">Total Size</div>
            <div className="text-lg font-bold">
              {formatFileSize(diffs.reduce((sum, diff) => sum + (diff.size || 0), 0))}
            </div>
          </div>
        </div>
      </div>

      {/* File List */}
      <div className="border border-white/10 rounded-xl overflow-hidden bg-white/[0.02]">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-white/10">
                <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                  File
                </th>
                <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                  Change Type
                </th>
                <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                  Size
                </th>
                <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                  Lines
                </th>
                <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                  Agent
                </th>
                <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                  Time
                </th>
                <th className="text-left p-4 text-xs font-bold uppercase tracking-wider text-white/40">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody>
              {filteredDiffs.map((diff) => {
                const FileIcon = getFileIcon(diff.filePath);
                const ChangeIcon = getChangeTypeIcon(diff.changeType);
                const isExpanded = expandedFiles.has(diff.id);
                const isSelected = selectedFile === diff.id;
                
                return (
                  <React.Fragment key={diff.id}>
                    <tr 
                      className={cn(
                        "border-b border-white/5 hover:bg-white/5 transition-colors cursor-pointer",
                        isSelected && "bg-primary/10"
                      )}
                      onClick={() => setSelectedFile(diff.id === selectedFile ? null : diff.id)}
                    >
                      <td className="p-4">
                        <div className="flex items-center gap-3">
                          <div className={cn(
                            "p-2 rounded-lg",
                            getChangeTypeBgColor(diff.changeType)
                          )}>
                            <FileIcon size={16} className={getChangeTypeColor(diff.changeType)} />
                          </div>
                          <div>
                            <div className="font-bold truncate max-w-xs">
                              {diff.filePath.split('/').pop()}
                            </div>
                            <div className="text-xs text-white/40 truncate max-w-xs">
                              {diff.filePath}
                            </div>
                          </div>
                        </div>
                      </td>
                      <td className="p-4">
                        <div className={cn(
                          "inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-bold",
                          getChangeTypeBgColor(diff.changeType),
                          getChangeTypeColor(diff.changeType)
                        )}>
                          <ChangeIcon size={12} />
                          {getChangeTypeLabel(diff.changeType)}
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="text-sm">{formatFileSize(diff.size || 0)}</div>
                      </td>
                      <td className="p-4">
                        <div className="text-sm">
                          <span className="text-green-500">+{diff.linesAdded || 0}</span>{' '}
                          <span className="text-red-500">-{diff.linesDeleted || 0}</span>
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="text-sm">{diff.agent || 'Unknown'}</div>
                      </td>
                      <td className="p-4">
                        <div className="text-sm text-white/60">
                          {diff.timestamp ? new Date(diff.timestamp).toLocaleTimeString() : 'Unknown'}
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="flex items-center gap-2">
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              toggleExpandFile(diff.id);
                            }}
                            className="p-1 rounded hover:bg-white/10 transition-colors"
                            title={isExpanded ? "Collapse diff" : "Expand diff"}
                          >
                            {isExpanded ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
                          </button>
                          
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleDownloadDiff(diff);
                            }}
                            className="p-1 rounded hover:bg-white/10 transition-colors"
                            title="Download diff"
                          >
                            <Download size={16} />
                          </button>
                          
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              navigator.clipboard.writeText(diff.diff || '');
                            }}
                            className="p-1 rounded hover:bg-white/10 transition-colors"
                            title="Copy diff"
                          >
                            <Copy size={16} />
                          </button>
                        </div>
                      </td>
                    </tr>
                    
                    {/* Expanded diff view */}
                    {isExpanded && (
                      <tr>
                        <td colSpan={7} className="p-0">
                          <div className="border-t border-white/10 bg-black/40">
                            <div className="p-4">
                              <div className="flex justify-between items-center mb-4">
                                <div className="text-sm font-bold">File Diff</div>
                                <div className="text-xs text-white/40">
                                  {diff.linesAdded || 0} additions, {diff.linesDeleted || 0} deletions
                                </div>
                              </div>
                              
                              <div className="border border-white/10 rounded-lg overflow-hidden">
                                <div className="bg-white/5 border-b border-white/10 p-3 flex justify-between items-center">
                                  <div className="font-mono text-sm">
                                    {diff.filePath}
                                  </div>
                                  <div className="flex items-center gap-2">
                                    <span className={cn(
                                      "text-xs px-2 py-1 rounded",
                                      getChangeTypeBgColor(diff.changeType),
                                      getChangeTypeColor(diff.changeType)
                                    )}>
                                      {getChangeTypeLabel(diff.changeType)}
                                    </span>
                                  </div>
                                </div>
                                
                                <div className="max-h-96 overflow-y-auto">
                                  {diff.diff ? (
                                    renderDiffLines(diff.diff)
                                  ) : (
                                    <div className="p-8 text-center text-white/40">
                                      No diff content available
                                    </div>
                                  )}
                                </div>
                              </div>
                            </div>
                          </div>
                        </td>
                      </tr>
                    )}
                  </React.Fragment>
                );
              })}
              
              {filteredDiffs.length === 0 && (
                <tr>
                  <td colSpan={7} className="p-8 text-center">
                    <FileText size={48} className="mx-auto text-white/20 mb-4" />
                    <div className="text-lg font-bold text-white/40">No file changes found</div>
                    <p className="text-white/60 mt-2">
                      {searchQuery 
                        ? `No files match "${searchQuery}"`
                        : "Wait for agents to modify files or check your filters"}
                    </p>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Selected File Details */}
      {selectedDiff && (
        <div className="border border-white/10 rounded-xl p-6 bg-white/[0.02]">
          <div className="flex justify-between items-start mb-6">
            <div>
              <h3 className="text-lg font-bold">File Change Details</h3>
              <p className="text-sm text-white/60 mt-1">
                {selectedDiff.filePath}
              </p>
            </div>
            
            <button
              onClick={() => setSelectedFile(null)}
              className="p-2 rounded-lg hover:bg-white/10 transition-colors"
            >
              <ChevronUp size={20} className="rotate-90" />
            </button>
          </div>
          
          <div className="grid grid-cols-2 gap-6">
            <div className="space-y-4">
              <div>
                <div className="text-sm font-bold mb-2">File Information</div>
                <div className="bg-white/5 border border-white/10 rounded-lg p-4 space-y-3">
                  <div className="flex justify-between">
                    <span className="text-white/60">Path:</span>
                    <span className="font-mono text-sm">{selectedDiff.filePath}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-white/60">Change Type:</span>
                    <span className={cn(
                      "font-bold",
                      getChangeTypeColor(selectedDiff.changeType)
                    )}>
                      {getChangeTypeLabel(selectedDiff.changeType)}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-white/60">Size:</span>
                    <span className="font-bold">{formatFileSize(selectedDiff.size || 0)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-white/60">Lines:</span>
                    <span className="font-bold">
                      <span className="text-green-500">+{selectedDiff.linesAdded || 0}</span>{' '}
                      <span className="text-red-500">-{selectedDiff.linesDeleted || 0}</span>
                    </span>
                  </div>
                </div>
              </div>
              
              <div>
                <div className="text-sm font-bold mb-2">Metadata</div>
                <div className="bg-white/5 border border-white/10 rounded-lg p-4 space-y-3">
                  <div className="flex justify-between">
                    <span className="text-white/60">Agent:</span>
                    <span className="font-bold">{selectedDiff.agent || 'Unknown'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-white/60">Timestamp:</span>
                    <span className="font-bold">
                      {selectedDiff.timestamp ? new Date(selectedDiff.timestamp).toLocaleString() : 'Unknown'}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-white/60">Commit Hash:</span>
                    <span className="font-mono text-sm">{selectedDiff.commitHash || 'N/A'}</span>
                  </div>
                </div>
              </div>
            </div>
            
            <div>
              <div className="text-sm font-bold mb-2">Diff Preview</div>
              <div className="border border-white/10 rounded-lg overflow-hidden">
                <div className="bg-white/5 border-b border-white/10 p-3">
                  <div className="font-mono text-sm truncate">
                    {selectedDiff.filePath.split('/').pop()}
                  </div>
                </div>
                <div className="max-h-80 overflow-y-auto">
                  {selectedDiff.diff ? (
                    renderDiffLines(selectedDiff.diff)
                  ) : (
                    <div className="p-8 text-center text-white/40">
                      No diff content available
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>
          
          {selectedDiff.metadata && (
            <div className="mt-6">
              <div className="text-sm font-bold mb-2">Additional Metadata</div>
              <div className="bg-white/5 border border-white/10 rounded-lg p-4">
                <pre className="text-sm text-white/80 overflow-x-auto">
                  {JSON.stringify(selectedDiff.metadata, null, 2)}
                </pre>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Change Type Distribution */}
      <div className="border border-white/10 rounded-xl p-6">
        <h4 className="text-lg font-bold mb-4">Change Type Distribution</h4>
        <div className="space-y-4">
          {[
            { type: 'added', label: 'Added', color: 'bg-green-500', icon: Plus },
            { type: 'modified', label: 'Modified', color: 'bg-yellow-500', icon: Edit },
            { type: 'deleted', label: 'Deleted', color: 'bg-red-500', icon: Trash2 },
          ].map((item) => {
            const count = diffs.filter(d => d.changeType === item.type).length;
            const percentage = (count / diffs.length) * 100;
            const Icon = item.icon;
            
            return (
              <div key={item.type} className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className={`p-2 ${item.color}/20 rounded-lg`}>
                    <Icon size={16} className={item.color.replace('bg-', 'text-')} />
                  </div>
                  <span className="text-sm">{item.label}</span>
                </div>
                <div className="flex items-center gap-4">
                  <div className="w-48 h-2 bg-white/10 rounded-full overflow-hidden">
                    <div 
                      className={`h-full ${item.color} rounded-full`}
                      style={{ width: `${percentage}%` }}
                    />
                  </div>
                  <div className="text-right w-16">
                    <span className="text-sm font-bold">{count}</span>
                    <span className="text-xs text-white/40 ml-1">({percentage.toFixed(1)}%)</span>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* File Tree Preview */}
      <div className="border border-white/10 rounded-xl p-6">
        <h4 className="text-lg font-bold mb-4">File Tree Structure</h4>
        <div className="bg-white/5 border border-white/10 rounded-lg p-4 max-h-60 overflow-y-auto">
          <div className="space-y-1">
            {diffs.slice(0, 20).map((diff) => {
              const parts = diff.filePath.split('/');
              const fileName = parts.pop();
              const folderPath = parts.join('/');
              
              return (
                <div key={diff.id} className="flex items-center gap-2 py-1">
                  <div className="text-white/40" style={{ paddingLeft: `${parts.length * 1}rem` }}>
                    {parts.length > 0 && (
                      <>
                        <FolderOpen size={12} className="inline mr-1" />
                        {parts[parts.length - 1]}/
                      </>
                    )}
                  </div>
                  <div className={cn(
                    "flex items-center gap-1",
                    getChangeTypeColor(diff.changeType)
                  )}>
                    {React.createElement(getChangeTypeIcon(diff.changeType), { size: 12 })}
                    <span className="font-mono text-sm">{fileName}</span>
                  </div>
                </div>
              );
            })}
            
            {diffs.length > 20 && (
              <div className="text-center py-2 text-white/40 text-sm">
                ... and {diffs.length - 20} more files
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}