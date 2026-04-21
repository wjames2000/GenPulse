import React, { useState, useEffect } from 'react';
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
  MessageSquare,
  AlertCircle as AlertCircleIcon,
  Play,
  RefreshCw,
  Loader2
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';
import { api } from '../services/api';

interface Task {
  id: string;
  title: string;
  content: string;
  tag: string;
  tagColor: string;
  ai?: boolean;
  tasks?: number;
  completed?: boolean;
  agent: 'orchestrator' | 'architect' | 'developer' | 'qa';
}

export default function KanbanView() {
  const [tasks, setTasks] = useState<Task[]>([
    {
      id: '1',
      title: 'Logic Parsing',
      content: 'Deconstructing user intent into discrete functional blocks for architecture mapping...',
      tag: 'Analyzing',
      tagColor: 'text-primary border-primary',
      ai: true,
      tasks: 12,
      agent: 'orchestrator'
    },
    {
      id: '2',
      title: 'Service Mesh',
      content: 'Defining isolated node communication protocols for the core repository.',
      tag: 'Pending',
      tagColor: 'text-white/20 border-white/10',
      agent: 'orchestrator'
    },
    {
      id: '3',
      title: 'Database Schema',
      content: 'Designing database schema for user authentication and data persistence.',
      tag: 'Processing',
      tagColor: 'text-primary border-primary',
      agent: 'architect'
    },
    {
      id: '4',
      title: 'Express Routes',
      content: 'Standardized API surface area defined and implemented.',
      tag: 'Completed',
      tagColor: 'text-white/10 border-white/5',
      completed: true,
      agent: 'developer'
    },
    {
      id: '5',
      title: 'JWT Middleware',
      content: 'Strict token validation layer integration finalized.',
      tag: 'Completed',
      tagColor: 'text-white/10 border-white/5',
      completed: true,
      agent: 'developer'
    },
    {
      id: '6',
      title: 'Integration: Auth',
      content: 'Authentication integration testing and validation.',
      tag: 'Failure Detected',
      tagColor: 'text-error border-error',
      agent: 'qa'
    }
  ]);

  const [isRefreshing, setIsRefreshing] = useState(false);
  const [activeTask, setActiveTask] = useState<string | null>(null);

  const handleRefresh = async () => {
    setIsRefreshing(true);
    try {
      // 模拟从API获取更新
      const agents = await api.getAllAgentsStatus();
      
      // 根据agent状态更新任务
      const updatedTasks = tasks.map(task => {
        const agentStatus = agents[task.agent];
        if (agentStatus) {
          const isActive = agentStatus.state === 'active';
          const isCompleted = agentStatus.progress === 100;
          
          let newTag = task.tag;
          let newTagColor = task.tagColor;
          
          if (isCompleted && task.tag !== 'Completed') {
            newTag = 'Completed';
            newTagColor = 'text-white/10 border-white/5';
          } else if (isActive && task.tag !== 'Processing') {
            newTag = 'Processing';
            newTagColor = 'text-primary border-primary';
          }
          
          return {
            ...task,
            tag: newTag,
            tagColor: newTagColor,
            completed: isCompleted
          };
        }
        return task;
      });
      
      setTasks(updatedTasks);
    } catch (error) {
      console.error('Failed to refresh tasks:', error);
    } finally {
      setIsRefreshing(false);
    }
  };

  const handleRunTask = async (taskId: string) => {
    setActiveTask(taskId);
    try {
      const task = tasks.find(t => t.id === taskId);
      if (task) {
        await api.executeAgent(task.agent, task.title, {
          description: task.content
        });
        
        // 更新任务状态
        setTasks(prev => prev.map(t => 
          t.id === taskId 
            ? { ...t, tag: 'Processing', tagColor: 'text-primary border-primary', completed: false }
            : t
        ));
      }
    } catch (error) {
      console.error('Failed to run task:', error);
    } finally {
      setActiveTask(null);
    }
  };

  const handleViewLogs = (taskId: string) => {
    // 这里可以跳转到日志视图或显示模态框
    console.log('View logs for task:', taskId);
  };

  const getTasksByAgent = (agent: Task['agent']) => {
    return tasks.filter(task => task.agent === agent);
  };

  const getAgentCount = (agent: Task['agent']) => {
    return getTasksByAgent(agent).length;
  };

  return (
    <div className="flex-1 overflow-x-auto p-12 bg-[#0A0A0A] h-full custom-scrollbar">
      <div className="flex justify-between items-center mb-12">
        <div>
          <span className="text-[10px] uppercase tracking-[0.5em] text-white/40 block mb-2">Execution Kanban / Real-time</span>
          <h1 className="text-[80px] leading-[0.8] font-black tracking-tighter uppercase">
            Task<br/>Flow
          </h1>
        </div>
        <button
          onClick={handleRefresh}
          disabled={isRefreshing}
          className="text-[10px] font-black uppercase tracking-widest text-white/40 hover:text-white transition-colors flex items-center gap-2"
        >
          <RefreshCw size={12} className={cn(isRefreshing && "animate-spin")} />
          {isRefreshing ? 'Refreshing...' : 'Refresh'}
        </button>
      </div>

      <div className="flex gap-1 min-w-max h-full pb-8">
        <KanbanColumn 
          index="01"
          title="Orchestrator" 
          count={getAgentCount('orchestrator')} 
          icon={Brain} 
          color="text-primary"
        >
          {getTasksByAgent('orchestrator').map(task => (
            <KanbanCard 
              key={task.id}
              task={task}
              onRun={() => handleRunTask(task.id)}
              onViewLogs={() => handleViewLogs(task.id)}
              isActive={activeTask === task.id}
            />
          ))}
        </KanbanColumn>

        <KanbanColumn 
          index="02"
          title="Architect" 
          count={getAgentCount('architect')} 
          icon={Ruler} 
          color="text-white"
        >
          {getTasksByAgent('architect').map(task => (
            <KanbanCard 
              key={task.id}
              task={task}
              onRun={() => handleRunTask(task.id)}
              onViewLogs={() => handleViewLogs(task.id)}
              isActive={activeTask === task.id}
            />
          ))}
        </KanbanColumn>

        <KanbanColumn 
          index="03"
          title="Developer" 
          count={getAgentCount('developer')} 
          icon={Code2} 
          color="text-white/40"
        >
          {getTasksByAgent('developer').map(task => (
            <KanbanCard 
              key={task.id}
              task={task}
              onRun={() => handleRunTask(task.id)}
              onViewLogs={() => handleViewLogs(task.id)}
              isActive={activeTask === task.id}
            />
          ))}
        </KanbanColumn>

        <KanbanColumn 
          index="04"
          title="QA Lead" 
          count={getAgentCount('qa')} 
          icon={FactCheckIcon} 
          color="text-error"
        >
          {getTasksByAgent('qa').map(task => (
            <KanbanCard 
              key={task.id}
              task={task}
              onRun={() => handleRunTask(task.id)}
              onViewLogs={() => handleViewLogs(task.id)}
              isActive={activeTask === task.id}
            />
          ))}
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

interface KanbanCardProps {
  task: Task;
  onRun: () => void;
  onViewLogs: () => void;
  isActive: boolean;
  key?: string; // React key prop
}

function KanbanCard({ task, onRun, onViewLogs, isActive }: KanbanCardProps) {
  const isError = task.tag.includes('Failure') || task.tag.includes('Error');
  const isProcessing = task.tag === 'Processing';
  const isCompleted = task.completed;

  return (
    <div className={cn(
      "p-10 flex flex-col gap-8 transition-all duration-500 group relative border",
      isCompleted 
        ? "border-white/5 opacity-20 grayscale" 
        : isError
        ? "border-error/20 hover:border-error/40"
        : "border-white/10 hover:border-white/20 hover:bg-white/[0.03]",
      isActive && "border-primary/50 bg-primary/5"
    )}>
      <div className="flex justify-between items-start">
        <span className={cn("text-[9px] font-black uppercase tracking-widest px-3 py-1 border", task.tagColor)}>
          {task.tag}
        </span>
        {isCompleted ? (
          <CheckCircle2 size={16} className="text-primary" />
        ) : isProcessing ? (
          <Loader2 size={16} className="text-primary animate-spin" />
        ) : (
          <div className="flex gap-1">
            <div className="w-1 h-1 bg-white/20 rounded-full" />
            <div className="w-1 h-1 bg-white/20 rounded-full" />
            <div className="w-1 h-1 bg-white/20 rounded-full" />
          </div>
        )}
      </div>

      <div>
        <h3 className={cn(
          "text-2xl font-black uppercase tracking-tighter mb-4 transition-all",
          isCompleted 
            ? "text-white/20 line-through" 
            : isError
            ? "text-error"
            : "text-white group-hover:text-primary"
        )}>
          {task.title}
        </h3>
        <p className={cn(
          "text-xs leading-relaxed uppercase font-bold tracking-wide",
          isCompleted 
            ? "text-white/10 italic" 
            : isError
            ? "text-error/60"
            : "text-white/40"
        )}>
          {task.content}
        </p>
      </div>

      {(task.ai || task.tasks) && (
        <div className="flex items-center justify-between pt-8 border-t border-white/5">
          <div className="flex items-center gap-4">
            {task.ai && (
              <div className="flex items-center gap-2">
                <Brain size={14} className="text-primary" />
                <span className="text-[9px] font-black uppercase tracking-widest text-primary">Neural Node</span>
              </div>
            )}
          </div>
          {task.tasks && (
            <div className="flex items-center gap-3 text-white/20">
              <MessageSquare size={12} />
              <span className="text-[10px] font-black">{task.tasks}</span>
            </div>
          )}
        </div>
      )}

      {!isCompleted && (
        <div className="flex gap-1 mt-4">
          <button 
            onClick={onViewLogs}
            className="flex-1 py-4 text-[9px] font-black uppercase tracking-widest border border-white/10 text-white/40 hover:text-white transition-all"
          >
            Logs
          </button>
          <button 
            onClick={onRun}
            disabled={isProcessing || isActive}
            className={cn(
              "flex-1 py-4 text-[9px] font-black uppercase tracking-widest transition-all flex items-center justify-center gap-2",
              isError
                ? "bg-error text-black hover:bg-error/90"
                : "bg-primary text-black hover:scale-105"
            )}
          >
            {isProcessing || isActive ? (
              <>
                <Loader2 size={12} className="animate-spin" />
                Running...
              </>
            ) : (
              <>
                <Play size={12} />
                {isError ? 'Re-run' : 'Run'}
              </>
            )}
          </button>
        </div>
      )}
    </div>
  );
}


