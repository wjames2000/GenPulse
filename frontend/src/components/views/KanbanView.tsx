import React from 'react';

const KanbanView: React.FC = () => {
  const columns = [
    {
      id: 'todo',
      title: '待处理',
      color: 'outline',
      tasks: [
        { id: 1, title: '设计用户认证系统', priority: 'high', assignee: 'Architect' },
        { id: 2, title: '编写API文档', priority: 'medium', assignee: 'Documenter' },
        { id: 3, title: '优化数据库查询', priority: 'low', assignee: 'DBA' },
      ]
    },
    {
      id: 'in-progress',
      title: '进行中',
      color: 'primary',
      tasks: [
        { id: 4, title: '开发前端仪表盘', priority: 'high', assignee: 'Frontend Dev' },
        { id: 5, title: '实现WebSocket通信', priority: 'high', assignee: 'Backend Dev' },
      ]
    },
    {
      id: 'review',
      title: '审查中',
      color: 'tertiary',
      tasks: [
        { id: 6, title: '代码审查：用户模块', priority: 'medium', assignee: 'Code Reviewer' },
      ]
    },
    {
      id: 'done',
      title: '已完成',
      color: 'sys-success',
      tasks: [
        { id: 7, title: '项目初始化配置', priority: 'low', assignee: 'DevOps' },
        { id: 8, title: 'CI/CD流水线设置', priority: 'medium', assignee: 'DevOps' },
      ]
    }
  ];

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high': return 'bg-error/20 text-error';
      case 'medium': return 'bg-tertiary/20 text-tertiary';
      case 'low': return 'bg-outline/20 text-outline';
      default: return 'bg-outline/20 text-outline';
    }
  };

  const getPriorityText = (priority: string) => {
    switch (priority) {
      case 'high': return '高';
      case 'medium': return '中';
      case 'low': return '低';
      default: return '低';
    }
  };

  return (
    <div className="p-6 md:p-8 lg:p-12">
      <div className="max-w-7xl mx-auto">
        <div className="flex items-center justify-between mb-10">
          <div>
            <h2 className="text-[2.25rem] font-bold tracking-[-0.02em] text-on-surface mb-2">
              执行看板
            </h2>
            <p className="text-on-surface-variant text-sm">
              可视化跟踪和管理AI代理的任务执行流程。
            </p>
          </div>
          <div className="flex items-center gap-3">
            <button className="bg-surface-container text-on-surface hover:bg-surface-container-high rounded-lg px-4 py-2 text-sm font-medium transition-colors flex items-center gap-2">
              <span className="material-symbols-outlined">filter_list</span>
              筛选
            </button>
            <button className="bg-primary-container text-on-primary-container hover:bg-inverse-primary rounded-lg px-6 py-2 text-sm font-medium shadow-lg shadow-primary/20 transition-all flex items-center gap-2">
              <span className="material-symbols-outlined">add</span>
              新建任务
            </button>
          </div>
        </div>

        {/* 看板列 */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {columns.map((column) => (
            <div key={column.id} className="flex flex-col">
              {/* 列标题 */}
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-2">
                  <div className={`w-3 h-3 rounded-full bg-${column.color}`}></div>
                  <h3 className="text-lg font-semibold text-on-surface">{column.title}</h3>
                  <span className="text-sm font-mono text-outline bg-surface-container-low px-2 py-0.5 rounded">
                    {column.tasks.length}
                  </span>
                </div>
                <button className="text-outline hover:text-primary transition-colors">
                  <span className="material-symbols-outlined">more_vert</span>
                </button>
              </div>

              {/* 任务卡片 */}
              <div className="flex-1 space-y-4">
                {column.tasks.map((task) => (
                  <div key={task.id} className="bg-surface-container rounded-lg p-4 hover:bg-surface-container-high transition-colors cursor-move group">
                    <div className="flex items-start justify-between mb-3">
                      <h4 className="text-sm font-semibold text-on-surface flex-1">{task.title}</h4>
                      <button className="text-outline hover:text-primary transition-colors opacity-0 group-hover:opacity-100">
                        <span className="material-symbols-outlined text-sm">drag_indicator</span>
                      </button>
                    </div>
                    
                    <div className="flex items-center justify-between">
                      <div className={`text-xs px-2 py-1 rounded ${getPriorityColor(task.priority)}`}>
                        {getPriorityText(task.priority)}优先级
                      </div>
                      <div className="flex items-center gap-2">
                        <div className="w-6 h-6 rounded-full bg-surface-container-low flex items-center justify-center">
                          <span className="text-xs font-bold text-primary">{task.assignee.charAt(0)}</span>
                        </div>
                        <span className="text-xs text-outline">{task.assignee}</span>
                      </div>
                    </div>

                    {/* 任务操作 */}
                    <div className="mt-4 pt-3 border-t border-outline-variant/15 flex items-center justify-between">
                      <button className="text-xs text-outline hover:text-primary transition-colors flex items-center gap-1">
                        <span className="material-symbols-outlined text-sm">comment</span>
                        <span>3</span>
                      </button>
                      <div className="flex items-center gap-2">
                        <button className="text-xs text-outline hover:text-primary transition-colors">
                          <span className="material-symbols-outlined text-sm">edit</span>
                        </button>
                        <button className="text-xs text-outline hover:text-primary transition-colors">
                          <span className="material-symbols-outlined text-sm">delete</span>
                        </button>
                      </div>
                    </div>
                  </div>
                ))}

                {/* 添加任务按钮 */}
                <button className="w-full border-2 border-dashed border-outline-variant/30 hover:border-primary/50 rounded-lg p-4 text-center transition-colors group">
                  <span className="material-symbols-outlined text-outline group-hover:text-primary transition-colors">
                    add
                  </span>
                  <span className="text-sm text-outline group-hover:text-primary transition-colors ml-2">
                    添加任务
                  </span>
                </button>
              </div>
            </div>
          ))}
        </div>

        {/* 看板统计 */}
        <div className="mt-12 grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-surface-container rounded-xl p-6">
            <h3 className="text-lg font-semibold text-on-surface mb-4">任务统计</h3>
            <div className="space-y-4">
              <div>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-outline">总任务数</span>
                  <span className="font-mono text-primary">8</span>
                </div>
                <div className="h-2 bg-surface-container-low rounded-full"></div>
              </div>
              <div>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-outline">完成率</span>
                  <span className="font-mono text-sys-success">25%</span>
                </div>
                <div className="h-2 bg-surface-container-low rounded-full overflow-hidden">
                  <div className="h-full bg-sys-success rounded-full" style={{ width: '25%' }}></div>
                </div>
              </div>
              <div>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-outline">平均周期</span>
                  <span className="font-mono text-tertiary">2.5天</span>
                </div>
                <div className="h-2 bg-surface-container-low rounded-full overflow-hidden">
                  <div className="h-full bg-tertiary rounded-full" style={{ width: '50%' }}></div>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-surface-container rounded-xl p-6 md:col-span-2">
            <h3 className="text-lg font-semibold text-on-surface mb-4">活动时间线</h3>
            <div className="space-y-3">
              {[
                { action: '任务创建', task: '设计用户认证系统', time: '10:30 AM', user: 'Architect' },
                { action: '状态更新', task: '开发前端仪表盘', time: '11:15 AM', user: 'Frontend Dev' },
                { action: '代码提交', task: '项目初始化配置', time: '2:45 PM', user: 'DevOps' },
                { action: '审查通过', task: 'CI/CD流水线设置', time: '4:20 PM', user: 'Code Reviewer' },
              ].map((activity, index) => (
                <div key={index} className="flex items-center gap-3 py-2">
                  <div className="w-2 h-2 rounded-full bg-primary"></div>
                  <div className="flex-1">
                    <p className="text-sm text-on-surface">
                      <span className="font-medium">{activity.user}</span> {activity.action}: {activity.task}
                    </p>
                    <p className="text-xs text-outline">{activity.time}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default KanbanView;