import React from 'react';

const MemoryViewGenpulse: React.FC = () => {
  return (
    <main className="flex-1 p-6 grid grid-cols-12 gap-6 h-full overflow-hidden">
      {/* 左侧：情节记忆检索 */}
      <section className="col-span-3 flex flex-col gap-4 overflow-hidden">
        <div className="flex items-center justify-between px-1">
          <h2 className="text-headline-sm font-headline font-semibold tracking-tight text-on-surface">情节记忆检索</h2>
          <span className="material-symbols-outlined text-on-surface-variant text-[20px]">filter_list</span>
        </div>
        <div className="flex flex-col gap-3 flex-1 overflow-y-auto pr-2 custom-scrollbar">
          {/* 搜索和过滤器 */}
          <div className="bg-surface-container-lowest rounded-xl p-1 mb-2">
            <div className="flex items-center px-3 py-2 gap-2 text-on-surface-variant">
              <span className="material-symbols-outlined text-[18px]">search</span>
              <input 
                className="bg-transparent border-none outline-none text-sm w-full placeholder-on-surface-variant/50 focus:ring-0" 
                placeholder="搜索记忆..." 
                type="text"
              />
            </div>
          </div>
          <div className="flex gap-2 mb-2 flex-wrap">
            <span className="bg-surface-container-high text-on-surface text-[10px] font-label uppercase px-2.5 py-1 rounded-full border border-outline-variant/15 flex items-center gap-1 cursor-pointer hover:bg-surface-variant transition-colors">
              <span className="material-symbols-outlined text-[12px]">smart_toy</span> Agent 角色
            </span>
            <span className="bg-surface-container-high text-on-surface text-[10px] font-label uppercase px-2.5 py-1 rounded-full border border-outline-variant/15 flex items-center gap-1 cursor-pointer hover:bg-surface-variant transition-colors">
              <span className="material-symbols-outlined text-[12px]">schedule</span> 时间范围
            </span>
          </div>
          
          {/* 记忆卡片 */}
          <div className="bg-surface-container-low rounded-xl p-4 flex flex-col gap-3 cursor-pointer hover:bg-surface-bright transition-all duration-200 group">
            <div className="flex justify-between items-start">
              <div className="flex items-center gap-2">
                <div className="w-2 h-2 rounded-full bg-primary-container"></div>
                <span className="text-xs font-mono text-on-surface-variant">ID: 0x9F2A</span>
              </div>
              <span className="text-[10px] text-on-surface-variant/70">2h ago</span>
            </div>
            <h3 className="text-sm font-medium text-on-surface group-hover:text-primary transition-colors">鉴权模块重构</h3>
            <p className="text-xs text-on-surface-variant line-clamp-2 leading-relaxed">
              Updated the OAuth2 flow to support strict JWT validation. Replaced symmetric keys with asymmetric RSA payload signing.
            </p>
            <div className="flex gap-2 mt-1">
              <span className="text-[10px] text-primary-fixed-dim bg-primary-container/10 px-1.5 py-0.5 rounded">#security</span>
              <span className="text-[10px] text-on-surface-variant bg-surface-container-highest px-1.5 py-0.5 rounded">#core</span>
            </div>
          </div>

          <div className="bg-surface-container-low rounded-xl p-4 flex flex-col gap-3 cursor-pointer hover:bg-surface-bright transition-all duration-200 group">
            <div className="flex justify-between items-start">
              <div className="flex items-center gap-2">
                <div className="w-2 h-2 rounded-full bg-tertiary"></div>
                <span className="text-xs font-mono text-on-surface-variant">ID: 0x4C1B</span>
              </div>
              <span className="text-[10px] text-on-surface-variant/70">5h ago</span>
            </div>
            <h3 className="text-sm font-medium text-on-surface group-hover:text-primary transition-colors">API 端点生成</h3>
            <p className="text-xs text-on-surface-variant line-clamp-2 leading-relaxed">
              Generated RESTful endpoints for user profile management. Implemented rate limiting and schema validation middleware.
            </p>
            <div className="flex gap-2 mt-1">
              <span className="text-[10px] text-tertiary-fixed-dim bg-tertiary-container/20 px-1.5 py-0.5 rounded">#backend</span>
            </div>
          </div>

          {/* 更多记忆卡片 */}
          {[1, 2, 3].map((i) => (
            <div key={i} className="bg-surface-container-low rounded-xl p-4 flex flex-col gap-3 cursor-pointer hover:bg-surface-bright transition-all duration-200 group">
              <div className="flex justify-between items-start">
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-secondary"></div>
                  <span className="text-xs font-mono text-on-surface-variant">ID: 0x{Math.random().toString(16).substring(2, 6).toUpperCase()}</span>
                </div>
                <span className="text-[10px] text-on-surface-variant/70">{i * 3}h ago</span>
              </div>
              <h3 className="text-sm font-medium text-on-surface group-hover:text-primary transition-colors">记忆项目 #{i}</h3>
              <p className="text-xs text-on-surface-variant line-clamp-2 leading-relaxed">
                这是一个示例记忆描述，包含重要的项目决策和技术细节。
              </p>
              <div className="flex gap-2 mt-1">
                <span className="text-[10px] text-secondary-fixed-dim bg-secondary-container/10 px-1.5 py-0.5 rounded">#project</span>
                <span className="text-[10px] text-on-surface-variant bg-surface-container-highest px-1.5 py-0.5 rounded">#decision</span>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* 中间：记忆编辑器 */}
      <section className="col-span-6 flex flex-col overflow-hidden bg-surface-container-low rounded-2xl relative shadow-[0px_24px_48px_rgba(0,0,0,0.2)]">
        {/* 标签页头部 */}
        <div className="flex items-center bg-surface-container-highest/50 px-2 pt-2 gap-1 rounded-t-2xl">
          <button className="bg-surface-container-low text-primary px-4 py-2.5 rounded-t-lg flex items-center gap-2 text-xs font-mono font-medium relative overflow-hidden group">
            <span className="material-symbols-outlined text-[16px]">person</span>
            USER.md
            <span className="bg-surface-container-high text-[10px] px-1.5 rounded text-on-surface-variant ml-2">全局偏好</span>
            <div className="absolute bottom-0 left-0 w-full h-[2px] bg-primary"></div>
          </button>
          <button className="text-on-surface-variant hover:text-on-surface px-4 py-2.5 rounded-t-lg flex items-center gap-2 text-xs font-mono transition-colors">
            <span className="material-symbols-outlined text-[16px]">folder_data</span>
            MEMORY.md
            <span className="bg-surface-container-high/50 text-[10px] px-1.5 rounded text-on-surface-variant/70 ml-2">项目上下文</span>
          </button>
          <div className="ml-auto pr-2 flex items-center gap-2">
            <button className="text-on-surface-variant hover:text-primary transition-colors p-1 rounded-md hover:bg-surface-variant">
              <span className="material-symbols-outlined text-[18px]">save</span>
            </button>
            <button className="text-on-surface-variant hover:text-primary transition-colors p-1 rounded-md hover:bg-surface-variant">
              <span className="material-symbols-outlined text-[18px]">more_vert</span>
            </button>
          </div>
        </div>
        
        {/* 编辑器主体 */}
        <div className="flex-1 overflow-y-auto p-6 font-mono text-sm text-on-surface leading-relaxed flex flex-col gap-6 custom-scrollbar relative">
          {/* 装饰性模糊元素，增强玻璃质感 */}
          <div className="absolute top-0 right-0 w-64 h-64 bg-primary/5 rounded-full blur-[80px] pointer-events-none"></div>
          
          {/* 编辑器内容部分 */}
          <div>
            <h4 className="text-primary-fixed-dim font-medium mb-3 flex items-center gap-2">
              <span className="material-symbols-outlined text-[18px]">code_blocks</span>
              编码风格指南
            </h4>
            <div className="bg-surface-container-lowest/80 rounded-xl p-4 border border-outline-variant/10 text-[13px]">
              <div className="flex gap-4">
                <div className="text-outline-variant select-none text-right flex flex-col gap-1 w-6">
                  <span>1</span><span>2</span><span>3</span><span>4</span><span>5</span>
                </div>
                <div className="flex flex-col gap-1 text-on-surface/90">
                  <span className="text-secondary-fixed-dim"># 严格遵循 TypeScript 强类型</span>
                  <span><span className="text-tertiary">interface</span> <span className="text-primary-container">CognitivePayload</span> {`{`}</span>
                  <span className="pl-4">intent: <span className="text-success">string</span>;</span>
                  <span className="pl-4">confidence: <span className="text-success">number</span>;</span>
                  <span>{`}`}</span>
                </div>
              </div>
            </div>
          </div>
          
          <div>
            <h4 className="text-primary-fixed-dim font-medium mb-3 flex items-center gap-2">
              <span className="material-symbols-outlined text-[18px]">gavel</span>
              技术决策约定
            </h4>
            <ul className="list-none space-y-3">
              <li className="flex items-start gap-3">
                <div className="mt-1 w-1.5 h-1.5 rounded-full bg-outline-variant"></div>
                <div>
                  <span className="font-medium text-on-surface block mb-0.5">状态管理</span>
                  <span className="text-on-surface-variant text-xs">优先使用 Zustand。禁止在小型模块中引入 Redux，以保持 Cognitive Loop 的轻量化执行。</span>
                </div>
              </li>
              <li className="flex items-start gap-3">
                <div className="mt-1 w-1.5 h-1.5 rounded-full bg-outline-variant"></div>
                <div>
                  <span className="font-medium text-on-surface block mb-0.5">API 响应结构</span>
                  <span className="text-on-surface-variant text-xs">所有响应必须包裹在标准 <code className="bg-surface-container px-1 py-0.5 rounded text-[11px] text-tertiary-fixed-dim">DataWrapper&lt;T&gt;</code> 中，包含 meta 诊断信息。</span>
                </div>
              </li>
              <li className="flex items-start gap-3">
                <div className="mt-1 w-1.5 h-1.5 rounded-full bg-outline-variant"></div>
                <div>
                  <span className="font-medium text-on-surface block mb-0.5">错误处理</span>
                  <span className="text-on-surface-variant text-xs">使用结构化错误对象，包含错误码、用户友好消息和技术详情。</span>
                </div>
              </li>
            </ul>
          </div>
          
          <div>
            <h4 className="text-primary-fixed-dim font-medium mb-3 flex items-center gap-2">
              <span className="material-symbols-outlined text-[18px]">security</span>
              安全规范
            </h4>
            <div className="bg-surface-container-lowest rounded-xl p-4 border border-outline-variant/10">
              <p className="text-sm text-on-surface-variant mb-2">所有鉴权令牌必须使用 RSA 非对称加密签名，公钥存储在配置中心。</p>
              <p className="text-sm text-on-surface-variant">敏感环境变量必须通过密钥管理服务访问，禁止硬编码。</p>
            </div>
          </div>
        </div>
      </section>

      {/* 右侧：交互时间线 */}
      <section className="col-span-3 flex flex-col gap-4 overflow-hidden">
        <div className="flex items-center justify-between px-1">
          <h2 className="text-headline-sm font-headline font-semibold tracking-tight text-on-surface">交互时间线</h2>
          <div className="flex items-center gap-2">
            <span className="relative flex h-2 w-2">
              <span className="absolute inline-flex h-full w-full rounded-full bg-success opacity-40"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-success"></span>
            </span>
            <span className="text-[10px] font-label uppercase text-success tracking-wider">Live</span>
          </div>
        </div>
        <div className="flex-1 overflow-y-auto pr-2 custom-scrollbar relative">
          {/* 时间线 */}
          <div className="absolute left-[15px] top-4 bottom-0 w-px bg-outline-variant/30 border-l border-dashed border-outline-variant/50"></div>
          <div className="flex flex-col gap-6 pt-2">
            {/* 时间线项目：AI 活动 */}
            <div className="relative flex gap-4">
              <div className="relative z-10 w-8 h-8 rounded-full bg-surface-container-highest border border-outline-variant/30 flex items-center justify-center shrink-0 mt-1 shadow-[0_0_15px_rgba(91,95,255,0.15)]">
                <span className="material-symbols-outlined text-[16px] text-primary">psychology</span>
              </div>
              <div className="flex-1 bg-surface-container-highest/60 backdrop-blur-md rounded-xl p-3 border border-outline-variant/10 rounded-tl-none">
                <div className="flex justify-between items-center mb-1">
                  <span className="text-[11px] font-semibold text-primary">分析意图 (AI)</span>
                  <span className="text-[10px] text-on-surface-variant font-mono">Just now</span>
                </div>
                <p className="text-xs text-on-surface-variant leading-relaxed">
                  正在解析用户输入的鉴权需求，交叉对比 <code className="text-[10px] bg-surface-container px-1 rounded text-secondary">USER.md</code> 中的安全偏好...
                </p>
              </div>
            </div>
            
            {/* 时间线项目：差异对比 */}
            <div className="relative flex gap-4 opacity-80 hover:opacity-100 transition-opacity">
              <div className="relative z-10 w-8 h-8 rounded-full bg-surface-container-low flex items-center justify-center shrink-0 mt-1">
                <span className="material-symbols-outlined text-[14px] text-tertiary">difference</span>
              </div>
              <div className="flex-1 bg-surface-container-low rounded-xl p-3 border border-outline-variant/5">
                <div className="flex justify-between items-center mb-2">
                  <span className="text-[11px] font-semibold text-on-surface">提议代码变更</span>
                  <span className="text-[10px] text-on-surface-variant font-mono">2m ago</span>
                </div>
                <div className="rounded-lg overflow-hidden text-[10px] font-mono border border-outline-variant/10">
                  <div className="bg-error-container/20 text-on-surface px-2 py-1 flex gap-2">
                    <span className="text-error/70">-</span> <span className="text-on-surface-variant line-through">validateToken(req.body)</span>
                  </div>
                  <div className="bg-primary-container/15 text-on-surface px-2 py-1 flex gap-2">
                    <span className="text-primary/70">+</span> <span className="text-on-primary-container">await jwt.verify(header, RS256)</span>
                  </div>
                </div>
              </div>
            </div>
            
            {/* 时间线项目：用户操作 */}
            <div className="relative flex gap-4 opacity-60">
              <div className="relative z-10 w-8 h-8 rounded-full bg-surface-container-low flex items-center justify-center shrink-0 mt-1">
                <span className="material-symbols-outlined text-[14px] text-on-surface-variant">person</span>
              </div>
              <div className="flex-1 bg-transparent p-1 pt-2">
                <div className="flex justify-between items-center">
                  <span className="text-[11px] text-on-surface-variant">操作员批准了 <strong>auth.ts</strong> 的修改</span>
                  <span className="text-[10px] text-on-surface-variant/50 font-mono">5m ago</span>
                </div>
              </div>
            </div>
            
            {/* 更多时间线项目 */}
            {[1, 2, 3].map((i) => (
              <div key={i} className="relative flex gap-4 opacity-70">
                <div className="relative z-10 w-8 h-8 rounded-full bg-surface-container-low flex items-center justify-center shrink-0 mt-1">
                  <span className="material-symbols-outlined text-[14px] text-secondary">code</span>
                </div>
                <div className="flex-1 bg-surface-container-low rounded-xl p-3 border border-outline-variant/5">
                  <div className="flex justify-between items-center mb-1">
                    <span className="text-[11px] font-semibold text-on-surface">代码生成</span>
                    <span className="text-[10px] text-on-surface-variant font-mono">{i * 10}m ago</span>
                  </div>
                  <p className="text-xs text-on-surface-variant">生成了新的 API 端点用于用户管理。</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* 自定义滚动条样式 - 浏览器默认 */}
    </main>
  );
};

export default MemoryViewGenpulse;