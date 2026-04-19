import React, { useState } from 'react';

const TerminalView: React.FC = () => {
  const [command, setCommand] = useState('');
  const [activeTab, setActiveTab] = useState('terminal');

  const terminalOutput = [
    { type: 'info', text: '🚀 Genpulse AI Terminal v1.0.0' },
    { type: 'info', text: '📦 初始化系统环境...' },
    { type: 'success', text: '✅ 环境初始化完成' },
    { type: 'info', text: '🤖 启动AI代理服务...' },
    { type: 'success', text: '✅ Agent服务已启动 (3个代理活跃)' },
    { type: 'info', text: '🔗 连接模型提供商...' },
    { type: 'success', text: '✅ OpenAI: 已连接 (gpt-4-turbo)' },
    { type: 'success', text: '✅ Anthropic: 已连接 (claude-3-opus)' },
    { type: 'warning', text: '⚠️  Google Gemini: 未配置' },
    { type: 'info', text: '📊 系统状态: 运行正常' },
    { type: 'prompt', text: 'genpulse@ai-dev ~ $' },
  ];

  const diffContent = `diff --git a/src/components/Button.tsx b/src/components/Button.tsx
index a1b2c3d..e4f5g6h 100644
--- a/src/components/Button.tsx
+++ b/src/components/Button.tsx
@@ -1,5 +1,5 @@
 import React from 'react';
-import './Button.css';
+import styles from './Button.module.css';

 interface ButtonProps {
   children: React.ReactNode;
@@ -10,7 +10,7 @@ interface ButtonProps {
 export const Button: React.FC<ButtonProps> = ({ children, variant = 'primary', onClick }) => {
   return (
     <button
-      className={\`btn btn-\${variant}\`}
+      className={\`\${styles.button} \${styles[variant]}\`}
       onClick={onClick}
     >
       {children}
@@ -18,4 +18,4 @@ export const Button: React.FC<ButtonProps> = ({ children, variant = 'primary', onClick }) => {
   );
 };

-export default Button;
\\ No newline at end of file
+export default Button;`;

  const handleCommandSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (command.trim()) {
      console.log('执行命令:', command);
      setCommand('');
    }
  };

  const getOutputClass = (type: string) => {
    switch (type) {
      case 'info': return 'text-blue-400';
      case 'success': return 'text-green-400';
      case 'warning': return 'text-yellow-400';
      case 'error': return 'text-red-400';
      case 'prompt': return 'text-primary font-mono';
      default: return 'text-on-surface';
    }
  };

  return (
    <div className="p-6 md:p-8 lg:p-12">
      <div className="max-w-7xl mx-auto">
        <div className="mb-10">
          <h2 className="text-[2.25rem] font-bold tracking-[-0.02em] text-on-surface mb-2">
            终端与对比
          </h2>
          <p className="text-on-surface-variant text-sm">
            执行命令、查看日志和代码差异对比。
          </p>
        </div>

        {/* 标签页 */}
        <div className="flex border-b border-outline-variant/15 mb-6">
          <button
            className={`pb-3 px-4 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'terminal'
                ? 'border-primary text-primary'
                : 'border-transparent text-outline hover:text-on-surface'
            }`}
            onClick={() => setActiveTab('terminal')}
          >
            <span className="flex items-center gap-2">
              <span className="material-symbols-outlined">terminal</span>
              终端
            </span>
          </button>
          <button
            className={`pb-3 px-4 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'diff'
                ? 'border-primary text-primary'
                : 'border-transparent text-outline hover:text-on-surface'
            }`}
            onClick={() => setActiveTab('diff')}
          >
            <span className="flex items-center gap-2">
              <span className="material-symbols-outlined">difference</span>
              代码对比
            </span>
          </button>
          <button
            className={`pb-3 px-4 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'logs'
                ? 'border-primary text-primary'
                : 'border-transparent text-outline hover:text-on-surface'
            }`}
            onClick={() => setActiveTab('logs')}
          >
            <span className="flex items-center gap-2">
              <span className="material-symbols-outlined">list_alt</span>
              日志
            </span>
          </button>
        </div>

        {/* 终端内容 */}
        {activeTab === 'terminal' && (
          <div className="bg-surface-container rounded-xl overflow-hidden">
            {/* 终端头部 */}
            <div className="bg-surface-container-high px-4 py-3 flex items-center justify-between border-b border-outline-variant/15">
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded-full bg-error"></div>
                <div className="w-3 h-3 rounded-full bg-tertiary"></div>
                <div className="w-3 h-3 rounded-full bg-sys-success"></div>
                <span className="text-sm font-mono text-outline ml-2">genpulse-terminal</span>
              </div>
              <div className="flex items-center gap-3">
                <button className="text-outline hover:text-primary transition-colors">
                  <span className="material-symbols-outlined text-sm">content_copy</span>
                </button>
                <button className="text-outline hover:text-primary transition-colors">
                  <span className="material-symbols-outlined text-sm">delete</span>
                </button>
              </div>
            </div>

            {/* 终端输出 */}
            <div className="p-4 font-mono text-sm h-96 overflow-y-auto bg-surface-container-lowest">
              {terminalOutput.map((line, index) => (
                <div key={index} className={`mb-1 ${getOutputClass(line.type)}`}>
                  {line.text}
                </div>
              ))}
              <form onSubmit={handleCommandSubmit} className="flex items-center mt-2">
                <span className="text-primary font-mono mr-2">genpulse@ai-dev ~ $</span>
                <input
                  className="flex-1 bg-transparent text-on-surface outline-none font-mono"
                  type="text"
                  value={command}
                  onChange={(e) => setCommand(e.target.value)}
                  autoFocus
                  placeholder="输入命令..."
                />
              </form>
            </div>

            {/* 终端底部 */}
            <div className="bg-surface-container-high px-4 py-2 border-t border-outline-variant/15">
              <div className="flex items-center justify-between text-xs text-outline">
                <div className="flex items-center gap-4">
                  <span>UTF-8</span>
                  <span>•</span>
                  <span>bash</span>
                </div>
                <div className="flex items-center gap-4">
                  <span>Ln 11, Col 45</span>
                  <span>•</span>
                  <span>100%</span>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* 代码对比内容 */}
        {activeTab === 'diff' && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="bg-surface-container rounded-xl overflow-hidden">
              <div className="bg-surface-container-high px-4 py-3 border-b border-outline-variant/15">
                <h3 className="text-sm font-semibold text-on-surface">原文件</h3>
                <p className="text-xs text-outline">src/components/Button.tsx</p>
              </div>
              <div className="p-4 font-mono text-sm h-96 overflow-y-auto bg-surface-container-lowest">
                <pre className="text-on-surface-variant whitespace-pre-wrap">
                  {diffContent.split('\n').map((line, index) => {
                    if (line.startsWith('-') && !line.startsWith('---')) {
                      return <div key={index} className="text-error bg-error/10">{line}</div>;
                    }
                    return <div key={index}>{line}</div>;
                  })}
                </pre>
              </div>
            </div>

            <div className="bg-surface-container rounded-xl overflow-hidden">
              <div className="bg-surface-container-high px-4 py-3 border-b border-outline-variant/15">
                <h3 className="text-sm font-semibold text-on-surface">修改后</h3>
                <p className="text-xs text-outline">src/components/Button.tsx</p>
              </div>
              <div className="p-4 font-mono text-sm h-96 overflow-y-auto bg-surface-container-lowest">
                <pre className="text-on-surface-variant whitespace-pre-wrap">
                  {diffContent.split('\n').map((line, index) => {
                    if (line.startsWith('+') && !line.startsWith('+++')) {
                      return <div key={index} className="text-sys-success bg-sys-success/10">{line}</div>;
                    }
                    return <div key={index}>{line}</div>;
                  })}
                </pre>
              </div>
            </div>
          </div>
        )}

        {/* 日志内容 */}
        {activeTab === 'logs' && (
          <div className="bg-surface-container rounded-xl overflow-hidden">
            <div className="bg-surface-container-high px-4 py-3 border-b border-outline-variant/15 flex items-center justify-between">
              <div>
                <h3 className="text-sm font-semibold text-on-surface">系统日志</h3>
                <p className="text-xs text-outline">实时应用日志流</p>
              </div>
              <div className="flex items-center gap-3">
                <button className="text-xs font-medium text-primary hover:text-primary-container transition-colors">
                  <span className="material-symbols-outlined text-sm">play_arrow</span>
                  开始
                </button>
                <button className="text-xs font-medium text-outline hover:text-on-surface transition-colors">
                  <span className="material-symbols-outlined text-sm">pause</span>
                  暂停
                </button>
                <button className="text-xs font-medium text-outline hover:text-on-surface transition-colors">
                  <span className="material-symbols-outlined text-sm">download</span>
                  导出
                </button>
              </div>
            </div>
            
            <div className="p-4 font-mono text-sm h-96 overflow-y-auto bg-surface-container-lowest">
              {[
                { time: '10:30:45', level: 'INFO', message: '应用程序启动成功' },
                { time: '10:31:12', level: 'INFO', message: '数据库连接已建立' },
                { time: '10:32:05', level: 'WARN', message: '内存使用率较高: 78%' },
                { time: '10:33:21', level: 'INFO', message: '用户登录: admin@example.com' },
                { time: '10:34:47', level: 'ERROR', message: 'API请求失败: 连接超时' },
                { time: '10:35:12', level: 'INFO', message: '重试机制已触发' },
                { time: '10:36:33', level: 'INFO', message: 'API请求成功' },
                { time: '10:37:55', level: 'INFO', message: '缓存已更新' },
                { time: '10:38:42', level: 'DEBUG', message: '处理用户请求: GET /api/projects' },
                { time: '10:39:18', level: 'INFO', message: '响应时间: 245ms' },
              ].map((log, index) => (
                <div key={index} className="flex items-start gap-4 py-2 border-b border-outline-variant/10 last:border-0">
                  <span className="text-xs text-outline flex-shrink-0">{log.time}</span>
                  <span className={`text-xs px-2 py-0.5 rounded flex-shrink-0 ${
                    log.level === 'ERROR' ? 'bg-error/20 text-error' :
                    log.level === 'WARN' ? 'bg-tertiary/20 text-tertiary' :
                    log.level === 'INFO' ? 'bg-primary/20 text-primary' :
                    'bg-outline/20 text-outline'
                  }`}>
                    {log.level}
                  </span>
                  <span className="text-sm text-on-surface-variant flex-1">{log.message}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* 快速命令 */}
        <div className="mt-8 grid grid-cols-2 md:grid-cols-4 gap-4">
          <button className="bg-surface-container hover:bg-surface-container-high rounded-lg p-4 text-center transition-colors group">
            <span className="material-symbols-outlined text-primary text-2xl mb-2">refresh</span>
            <p className="text-sm font-medium text-on-surface">重启服务</p>
            <p className="text-xs text-outline">systemctl restart genpulse</p>
          </button>
          <button className="bg-surface-container hover:bg-surface-container-high rounded-lg p-4 text-center transition-colors group">
            <span className="material-symbols-outlined text-primary text-2xl mb-2">monitoring</span>
            <p className="text-sm font-medium text-on-surface">查看状态</p>
            <p className="text-xs text-outline">genpulse status</p>
          </button>
          <button className="bg-surface-container hover:bg-surface-container-high rounded-lg p-4 text-center transition-colors group">
            <span className="material-symbols-outlined text-primary text-2xl mb-2">backup</span>
            <p className="text-sm font-medium text-on-surface">备份数据</p>
            <p className="text-xs text-outline">genpulse backup</p>
          </button>
          <button className="bg-surface-container hover:bg-surface-container-high rounded-lg p-4 text-center transition-colors group">
            <span className="material-symbols-outlined text-primary text-2xl mb-2">update</span>
            <p className="text-sm font-medium text-on-surface">更新系统</p>
            <p className="text-xs text-outline">genpulse update</p>
          </button>
        </div>
      </div>
    </div>
  );
};

export default TerminalView;