import React from 'react';

const MemoryView: React.FC = () => {
  return (
    <div className="p-6 md:p-8 lg:p-12">
      <div className="max-w-6xl mx-auto">
        <div className="mb-10">
          <h2 className="text-[2.25rem] font-bold tracking-[-0.02em] text-on-surface mb-2">
            神经资产
          </h2>
          <p className="text-on-surface-variant text-sm">
            管理和探索AI代理的长期记忆与知识库。
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* 左侧：记忆概览 */}
          <div className="lg:col-span-2">
            <div className="bg-surface-container rounded-xl p-6 mb-6">
              <h3 className="text-lg font-semibold text-on-surface mb-4">记忆概览</h3>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="bg-surface-container-low rounded-lg p-4">
                  <p className="text-sm text-outline mb-1">总记忆数</p>
                  <p className="text-2xl font-bold text-primary">1,247</p>
                </div>
                <div className="bg-surface-container-low rounded-lg p-4">
                  <p className="text-sm text-outline mb-1">今日新增</p>
                  <p className="text-2xl font-bold text-sys-success">42</p>
                </div>
                <div className="bg-surface-container-low rounded-lg p-4">
                  <p className="text-sm text-outline mb-1">记忆大小</p>
                  <p className="text-2xl font-bold text-secondary">3.2GB</p>
                </div>
                <div className="bg-surface-container-low rounded-lg p-4">
                  <p className="text-sm text-outline mb-1">检索次数</p>
                  <p className="text-2xl font-bold text-tertiary">1.2k</p>
                </div>
              </div>
            </div>

            {/* 记忆列表 */}
            <div className="bg-surface-container rounded-xl p-6">
              <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-semibold text-on-surface">最近记忆</h3>
                <button className="text-sm font-medium text-primary hover:text-primary-container transition-colors">
                  查看全部
                </button>
              </div>
              
              <div className="space-y-4">
                {[1, 2, 3, 4, 5].map((i) => (
                  <div key={i} className="flex items-start gap-4 p-4 hover:bg-surface-container-low rounded-lg transition-colors cursor-pointer group">
                    <div className="p-2 rounded-md bg-surface-container-lowest text-primary flex-shrink-0">
                      <span className="material-symbols-outlined">memory</span>
                    </div>
                    <div className="flex-1">
                      <h4 className="text-sm font-semibold text-on-surface mb-1">
                        项目架构设计模式 #{i}
                      </h4>
                      <p className="text-xs text-outline mb-2">
                        关于微服务架构的最佳实践和设计模式，包含容器化部署策略...
                      </p>
                      <div className="flex items-center gap-4 text-xs text-outline">
                        <span>由 Architect 创建</span>
                        <span>•</span>
                        <span>2小时前</span>
                        <span>•</span>
                        <span className="font-mono text-primary">相关性: 92%</span>
                      </div>
                    </div>
                    <span className="material-symbols-outlined text-outline group-hover:text-primary transition-colors">
                      arrow_forward
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* 右侧：记忆分类 */}
          <div>
            <div className="bg-surface-container rounded-xl p-6 sticky top-6">
              <h3 className="text-lg font-semibold text-on-surface mb-4">记忆分类</h3>
              <div className="space-y-3">
                {[
                  { name: '架构设计', count: 124, color: 'primary' },
                  { name: '代码模式', count: 89, color: 'secondary' },
                  { name: '测试策略', count: 67, color: 'tertiary' },
                  { name: '部署配置', count: 42, color: 'sys-success' },
                  { name: '性能优化', count: 38, color: 'error' },
                  { name: '安全实践', count: 31, color: 'outline' },
                ].map((category, index) => (
                  <div key={index} className="flex items-center justify-between p-3 hover:bg-surface-container-low rounded-lg transition-colors cursor-pointer group">
                    <div className="flex items-center gap-3">
                      <div className={`w-2 h-2 rounded-full bg-${category.color}`}></div>
                      <span className="text-sm text-on-surface">{category.name}</span>
                    </div>
                    <span className="text-xs font-mono text-outline">{category.count}</span>
                  </div>
                ))}
              </div>

              <div className="mt-8 pt-6 border-t border-outline-variant/15">
                <h4 className="text-sm font-semibold text-on-surface mb-3">记忆搜索</h4>
                <div className="relative">
                  <input
                    className="w-full bg-surface-container-lowest text-on-surface text-sm rounded-lg pl-9 pr-3 py-2 border-none focus:ring-0 focus:border-b-2 focus:border-b-primary transition-all placeholder:text-outline outline-none"
                    placeholder="搜索记忆..."
                    type="text"
                  />
                  <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-outline text-sm">
                    search
                  </span>
                </div>
                <button className="w-full mt-4 bg-primary-container text-on-primary-container hover:bg-inverse-primary rounded-lg py-2 text-sm font-medium shadow-lg shadow-primary/20 transition-all">
                  新建记忆
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default MemoryView;