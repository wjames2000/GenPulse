import React from 'react';
import { useAppStore } from '../../stores/appStore';

const DashboardView: React.FC = () => {
  const { appInfo, healthStatus, isInitialized } = useAppStore();

  return (
    <div className="p-6 md:p-8 lg:p-12">
      <div className="max-w-7xl mx-auto">
        {/* 欢迎区域 */}
        <div className="mb-10">
          <h2 className="text-[2.25rem] font-bold tracking-[-0.02em] text-on-surface mb-2">
            欢迎使用 Genpulse AI
          </h2>
          <p className="text-on-surface-variant text-sm">
            智能AI软件项目开发流水线 - 让开发更高效、更智能
          </p>
        </div>

        {/* 状态卡片 */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div className="bg-surface-container rounded-xl p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 rounded-lg bg-primary/10 text-primary">
                <span className="material-symbols-outlined">rocket_launch</span>
              </div>
              <span className="text-xs font-mono text-sys-success bg-sys-success/10 px-2 py-1 rounded">
                运行中
              </span>
            </div>
            <h3 className="text-2xl font-bold text-on-surface mb-1">3</h3>
            <p className="text-sm text-outline">活跃代理</p>
          </div>

          <div className="bg-surface-container rounded-xl p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 rounded-lg bg-secondary/10 text-secondary">
                <span className="material-symbols-outlined">folder</span>
              </div>
              <span className="text-xs font-mono text-primary bg-primary/10 px-2 py-1 rounded">
                5个项目
              </span>
            </div>
            <h3 className="text-2xl font-bold text-on-surface mb-1">12</h3>
            <p className="text-sm text-outline">进行中任务</p>
          </div>

          <div className="bg-surface-container rounded-xl p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 rounded-lg bg-tertiary/10 text-tertiary">
                <span className="material-symbols-outlined">memory</span>
              </div>
              <span className="text-xs font-mono text-tertiary bg-tertiary/10 px-2 py-1 rounded">
                1.2k记录
              </span>
            </div>
            <h3 className="text-2xl font-bold text-on-surface mb-1">87%</h3>
            <p className="text-sm text-outline">记忆命中率</p>
          </div>

          <div className="bg-surface-container rounded-xl p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 rounded-lg bg-sys-success/10 text-sys-success">
                <span className="material-symbols-outlined">speed</span>
              </div>
              <span className="text-xs font-mono text-sys-success bg-sys-success/10 px-2 py-1 rounded">
                良好
              </span>
            </div>
            <h3 className="text-2xl font-bold text-on-surface mb-1">245ms</h3>
            <p className="text-sm text-outline">平均响应时间</p>
          </div>
        </div>

        {/* 主要区域 */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* 左侧：快速操作 */}
          <div className="lg:col-span-2">
            <div className="bg-surface-container rounded-xl p-6 mb-6">
              <h3 className="text-lg font-semibold text-on-surface mb-4">快速操作</h3>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                <button className="bg-surface-container-low hover:bg-surface-container-high rounded-lg p-4 text-center transition-colors group">
                  <span className="material-symbols-outlined text-primary text-2xl mb-2">add</span>
                  <p className="text-sm font-medium text-on-surface">新建项目</p>
                  <p className="text-xs text-outline">开始AI开发流水线</p>
                </button>
                <button className="bg-surface-container-low hover:bg-surface-container-high rounded-lg p-4 text-center transition-colors group">
                  <span className="material-symbols-outlined text-primary text-2xl mb-2">play_arrow</span>
                  <p className="text-sm font-medium text-on-surface">运行代理</p>
                  <p className="text-xs text-outline">启动AI工作流</p>
                </button>
                <button className="bg-surface-container-low hover:bg-surface-container-high rounded-lg p-4 text-center transition-colors group">
                  <span className="material-symbols-outlined text-primary text-2xl mb-2">analytics</span>
                  <p className="text-sm font-medium text-on-surface">查看报告</p>
                  <p className="text-xs text-outline">性能分析与洞察</p>
                </button>
                <button className="bg-surface-container-low hover:bg-surface-container-high rounded-lg p-4 text-center transition-colors group">
                  <span className="material-symbols-outlined text-primary text-2xl mb-2">settings</span>
                  <p className="text-sm font-medium text-on-surface">系统设置</p>
                  <p className="text-xs text-outline">配置与环境变量</p>
                </button>
                <button className="bg-surface-container-low hover:bg-surface-container-high rounded-lg p-4 text-center transition-colors group">
                  <span className="material-symbols-outlined text-primary text-2xl mb-2">help</span>
                  <p className="text-sm font-medium text-on-surface">帮助文档</p>
                  <p className="text-xs text-outline">使用指南与API</p>
                </button>
                <button className="bg-surface-container-low hover:bg-surface-container-high rounded-lg p-4 text-center transition-colors group">
                  <span className="material-symbols-outlined text-primary text-2xl mb-2">download</span>
                  <p className="text-sm font-medium text-on-surface">导出数据</p>
                  <p className="text-xs text-outline">项目与日志导出</p>
                </button>
              </div>
            </div>

            {/* 最近活动 */}
            <div className="bg-surface-container rounded-xl p-6">
              <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-semibold text-on-surface">最近活动</h3>
                <button className="text-sm font-medium text-primary hover:text-primary-container transition-colors">
                  查看全部
                </button>
              </div>
              
              <div className="space-y-4">
                {[
                  { action: '项目创建', project: '电商平台后端', time: '10分钟前', user: 'Architect' },
                  { action: '代码生成', project: '用户认证模块', time: '25分钟前', user: 'Frontend Dev' },
                  { action: '测试通过', project: '支付接口测试', time: '1小时前', user: 'QA Agent' },
                  { action: '部署完成', project: 'API网关v2', time: '2小时前', user: 'DevOps' },
                  { action: '性能优化', project: '数据库查询', time: '3小时前', user: 'DBA' },
                ].map((activity, index) => (
                  <div key={index} className="flex items-center gap-4 p-3 hover:bg-surface-container-low rounded-lg transition-colors">
                    <div className="p-2 rounded-md bg-surface-container-lowest text-primary">
                      <span className="material-symbols-outlined">bolt</span>
                    </div>
                    <div className="flex-1">
                      <p className="text-sm font-medium text-on-surface">
                        <span className="font-semibold">{activity.user}</span> {activity.action}
                      </p>
                      <p className="text-xs text-outline">{activity.project}</p>
                    </div>
                    <span className="text-xs text-outline">{activity.time}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* 右侧：系统信息 */}
          <div>
            <div className="bg-surface-container rounded-xl p-6 mb-6">
              <h3 className="text-lg font-semibold text-on-surface mb-4">系统信息</h3>
              <div className="space-y-4">
                <div>
                  <p className="text-sm text-outline mb-1">应用版本</p>
                  <p className="text-sm font-mono text-on-surface">{appInfo?.version || '1.0.0'}</p>
                </div>
                <div>
                  <p className="text-sm text-outline mb-1">运行状态</p>
                  <p className={`text-sm font-mono ${healthStatus?.status === 'healthy' ? 'text-sys-success' : 'text-error'}`}>
                    {healthStatus?.status === 'healthy' ? '健康' : '异常'}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-outline mb-1">初始化状态</p>
                  <p className={`text-sm font-mono ${isInitialized ? 'text-sys-success' : 'text-error'}`}>
                    {isInitialized ? '已初始化' : '未初始化'}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-outline mb-1">运行时间</p>
                  <p className="text-sm font-mono text-on-surface">2天14小时</p>
                </div>
              </div>
            </div>

            {/* 资源使用 */}
            <div className="bg-surface-container rounded-xl p-6">
              <h3 className="text-lg font-semibold text-on-surface mb-4">资源使用</h3>
              <div className="space-y-4">
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span className="text-outline">CPU使用率</span>
                    <span className="font-mono text-primary">32%</span>
                  </div>
                  <div className="h-2 bg-surface-container-low rounded-full overflow-hidden">
                    <div className="h-full bg-primary rounded-full" style={{ width: '32%' }}></div>
                  </div>
                </div>
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span className="text-outline">内存使用</span>
                    <span className="font-mono text-secondary">45%</span>
                  </div>
                  <div className="h-2 bg-surface-container-low rounded-full overflow-hidden">
                    <div className="h-full bg-secondary rounded-full" style={{ width: '45%' }}></div>
                  </div>
                </div>
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span className="text-outline">磁盘空间</span>
                    <span className="font-mono text-tertiary">78%</span>
                  </div>
                  <div className="h-2 bg-surface-container-low rounded-full overflow-hidden">
                    <div className="h-full bg-tertiary rounded-full" style={{ width: '78%' }}></div>
                  </div>
                </div>
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span className="text-outline">网络带宽</span>
                    <span className="font-mono text-sys-success">12%</span>
                  </div>
                  <div className="h-2 bg-surface-container-low rounded-full overflow-hidden">
                    <div className="h-full bg-sys-success rounded-full" style={{ width: '12%' }}></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* 底部：模型状态 */}
        <div className="mt-8 bg-surface-container rounded-xl p-6">
          <h3 className="text-lg font-semibold text-on-surface mb-4">模型提供商状态</h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="flex items-center gap-4 p-4 bg-surface-container-low rounded-lg">
              <div className="p-2 rounded-md bg-surface-container-lowest text-primary">
                <span className="material-symbols-outlined">data_object</span>
              </div>
              <div>
                <p className="text-sm font-medium text-on-surface">OpenAI</p>
                <p className="text-xs text-outline">gpt-4-turbo</p>
              </div>
              <div className="ml-auto">
                <span className="flex h-2 w-2 rounded-full bg-sys-success shadow-[0_0_8px_rgba(74,222,128,0.4)]"></span>
              </div>
            </div>
            
            <div className="flex items-center gap-4 p-4 bg-surface-container-low rounded-lg">
              <div className="p-2 rounded-md bg-surface-container-lowest text-primary">
                <span className="material-symbols-outlined">psychology_alt</span>
              </div>
              <div>
                <p className="text-sm font-medium text-on-surface">Anthropic</p>
                <p className="text-xs text-outline">claude-3-opus</p>
              </div>
              <div className="ml-auto">
                <span className="flex h-2 w-2 rounded-full bg-sys-success shadow-[0_0_8px_rgba(74,222,128,0.4)]"></span>
              </div>
            </div>
            
            <div className="flex items-center gap-4 p-4 bg-surface-container-low rounded-lg opacity-60">
              <div className="p-2 rounded-md bg-surface-container-lowest text-on-surface">
                <span className="material-symbols-outlined">auto_awesome</span>
              </div>
              <div>
                <p className="text-sm font-medium text-on-surface">Google Gemini</p>
                <p className="text-xs text-outline">未配置</p>
              </div>
              <div className="ml-auto">
                <span className="flex h-2 w-2 rounded-full bg-error shadow-[0_0_8px_rgba(255,180,171,0.4)]"></span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardView;