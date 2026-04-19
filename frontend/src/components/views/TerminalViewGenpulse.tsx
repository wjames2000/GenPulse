import React from 'react';

const TerminalViewGenpulse: React.FC = () => {
  return (
    <div className="flex-1 flex flex-col h-full overflow-hidden">
      {/* 上下文标题和操作按钮 */}
      <div className="flex items-end justify-between shrink-0 p-6 pb-4">
        <div className="flex flex-col gap-2">
          <div className="flex items-center gap-2 text-xs font-medium text-outline uppercase tracking-wider font-label">
            <span className="material-symbols-outlined text-[16px]">folder_open</span>
            <span>src / core / agents /</span>
            <span className="text-on-surface">reasoning_engine.ts</span>
          </div>
          <h1 className="text-[2.25rem] font-bold leading-none tracking-[-0.02em] font-headline text-on-surface flex items-center gap-4">
            代码对比解析
            <span className="px-2 py-1 rounded bg-surface-container-high border border-outline-variant/15 text-[0.6875rem] font-medium tracking-widest text-primary align-middle">待审核</span>
          </h1>
        </div>
        <div className="flex items-center gap-3">
          <button className="px-5 py-2 rounded-lg bg-surface-container-high text-on-surface text-sm font-medium hover:bg-surface-bright transition-all flex items-center gap-2 shadow-[0_4px_12px_rgba(0,0,0,0.2)]">
            <span className="material-symbols-outlined text-[18px] text-error">close</span>
            拒绝变更
          </button>
          <button className="px-5 py-2 rounded-lg bg-primary-container text-on-primary-container text-sm font-medium hover:brightness-110 transition-all flex items-center gap-2 shadow-[0_4px_12px_rgba(0,0,0,0.2)]">
            <span className="material-symbols-outlined text-[18px]">check</span>
            接受并应用
          </button>
        </div>
      </div>

      {/* 并排对比容器 */}
      <div className="flex-1 flex gap-4 overflow-hidden p-6 pt-0">
        <div className="flex-1 flex gap-4 overflow-hidden rounded-xl border border-outline-variant/15 p-1 bg-surface-container-highest shadow-[0_8px_32px_rgba(0,0,0,0.2)]">
          {/* 左侧面板：原始代码 */}
          <div className="flex-1 flex flex-col bg-surface overflow-hidden rounded-lg">
            <div className="px-4 py-2 bg-surface-container border-b border-outline-variant/15 flex justify-between items-center text-xs text-outline font-medium">
              <span>原始: HEAD</span>
              <span className="font-code opacity-50">v1.2.8</span>
            </div>
            <div className="flex-1 overflow-auto font-code text-[0.875rem] leading-relaxed select-text pb-4">
              {/* 代码行 */}
              <div className="flex hover:bg-surface-container-low/50 transition-colors">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-outline-variant select-none border-r border-outline-variant/15">42</div>
                <div className="pl-4 py-0.5 text-on-surface-variant whitespace-pre">  async function analyzeContext(data: any) {'{'}</div>
              </div>
              <div className="flex hover:bg-surface-container-low/50 transition-colors">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-outline-variant select-none border-r border-outline-variant/15">43</div>
                <div className="pl-4 py-0.5 text-on-surface-variant whitespace-pre">    const prompt = buildPrompt(data);</div>
              </div>
              {/* 删除块 */}
              <div className="flex bg-error-container/20 border-l-2 border-error">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-error/70 select-none border-r border-error/20">44</div>
                <div className="pl-4 py-0.5 text-on-surface whitespace-pre">    const response = await llm.generate(prompt);</div>
              </div>
              <div className="flex bg-error-container/20 border-l-2 border-error">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-error/70 select-none border-r border-error/20">45</div>
                <div className="pl-4 py-0.5 text-on-surface whitespace-pre">    return parseResult(response);</div>
              </div>
              <div className="flex hover:bg-surface-container-low/50 transition-colors">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-outline-variant select-none border-r border-outline-variant/15">46</div>
                <div className="pl-4 py-0.5 text-on-surface-variant whitespace-pre">  {'}'}</div>
              </div>
              
              {/* 更多示例代码行 */}
              {[47, 48, 49, 50].map((lineNum) => (
                <div key={lineNum} className="flex hover:bg-surface-container-low/50 transition-colors">
                  <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-outline-variant select-none border-r border-outline-variant/15">{lineNum}</div>
                  <div className="pl-4 py-0.5 text-on-surface-variant whitespace-pre">    // 示例代码行 {lineNum}</div>
                </div>
              ))}
            </div>
          </div>

          {/* 右侧面板：修改后的代码 */}
          <div className="flex-1 flex flex-col bg-surface overflow-hidden rounded-lg">
            <div className="px-4 py-2 bg-surface-container border-b border-outline-variant/15 flex justify-between items-center text-xs text-primary font-medium">
              <span className="flex items-center gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-primary animate-pulse"></div>
                Agent 修改
              </span>
              <span className="font-code opacity-50">本地编辑</span>
            </div>
            <div className="flex-1 overflow-auto font-code text-[0.875rem] leading-relaxed select-text pb-4">
              <div className="flex hover:bg-surface-container-low/50 transition-colors">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-outline-variant select-none border-r border-outline-variant/15">42</div>
                <div className="pl-4 py-0.5 text-on-surface-variant whitespace-pre">  async function analyzeContext(data: any) {'{'}</div>
              </div>
              <div className="flex hover:bg-surface-container-low/50 transition-colors">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-outline-variant select-none border-r border-outline-variant/15">43</div>
                <div className="pl-4 py-0.5 text-on-surface-variant whitespace-pre">    const prompt = buildPrompt(data);</div>
              </div>
              {/* 添加块 */}
              <div className="flex bg-primary-container/15 border-l-2 border-primary">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-primary/70 select-none border-r border-primary/20">44</div>
                <div className="pl-4 py-0.5 text-on-surface whitespace-pre">    const contextOpts = {`{ temperature: 0.2, top_p: 0.9 }`};</div>
              </div>
               <div className="flex bg-primary-container/15 border-l-2 border-primary">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-primary/70 select-none border-r border-primary/20">45</div>
                <div className="pl-4 py-0.5 text-on-surface whitespace-pre">    try {'{'}</div>
              </div>
              <div className="flex bg-primary-container/15 border-l-2 border-primary">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-primary/70 select-none border-r border-primary/20">46</div>
                <div className="pl-4 py-0.5 text-on-surface whitespace-pre">      const response = await llm.generate(prompt, contextOpts);</div>
              </div>
              <div className="flex bg-primary-container/15 border-l-2 border-primary">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-primary/70 select-none border-r border-primary/20">47</div>
                <div className="pl-4 py-0.5 text-on-surface whitespace-pre">      return validateAndParse(response);</div>
              </div>
               <div className="flex bg-primary-container/15 border-l-2 border-primary">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-primary/70 select-none border-r border-primary/20">48</div>
                <div className="pl-4 py-0.5 text-on-surface whitespace-pre">    {'}'} catch (e) {'{'}</div>
              </div>
              <div className="flex bg-primary-container/15 border-l-2 border-primary">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-primary/70 select-none border-r border-primary/20">49</div>
                <div className="pl-4 py-0.5 text-on-surface whitespace-pre">      logger.error('Generation failed', e);</div>
              </div>
              <div className="flex bg-primary-container/15 border-l-2 border-primary">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-primary/70 select-none border-r border-primary/20">50</div>
                <div className="pl-4 py-0.5 text-on-surface whitespace-pre">      throw new CriticalAgentFault();</div>
              </div>
               <div className="flex bg-primary-container/15 border-l-2 border-primary">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-primary/70 select-none border-r border-primary/20">51</div>
                <div className="pl-4 py-0.5 text-on-surface whitespace-pre">    {'}'}</div>
              </div>
               <div className="flex hover:bg-surface-container-low/50 transition-colors">
                <div className="w-12 shrink-0 text-right pr-4 py-0.5 text-outline-variant select-none border-r border-outline-variant/15">52</div>
                <div className="pl-4 py-0.5 text-on-surface-variant whitespace-pre">  {'}'}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 终端面板 */}
      <div className="h-64 shrink-0 bg-surface-container-lowest border-t border-outline-variant/20 flex flex-col font-code relative z-10 shadow-[0_-8px_24px_rgba(0,0,0,0.3)]">
        {/* 终端标题 */}
        <div className="flex items-center justify-between px-4 py-2 bg-[#0a0a0d] border-b border-outline-variant/10">
          <div className="flex items-center gap-4">
            <div className="flex gap-1.5">
              <div className="w-2.5 h-2.5 rounded-full bg-surface-variant"></div>
              <div className="w-2.5 h-2.5 rounded-full bg-surface-variant"></div>
              <div className="w-2.5 h-2.5 rounded-full bg-surface-variant"></div>
            </div>
            <span className="text-[0.6875rem] text-outline tracking-wider uppercase font-body font-medium">终端进程 : bash</span>
          </div>
          <div className="flex gap-2">
            <button className="text-outline hover:text-white transition-colors">
              <span className="material-symbols-outlined text-[16px]">clear_all</span>
            </button>
            <button className="text-outline hover:text-white transition-colors">
              <span className="material-symbols-outlined text-[16px]">expand_less</span>
            </button>
          </div>
        </div>
        {/* 终端输出 */}
        <div className="flex-1 overflow-auto p-4 text-[0.8125rem] text-on-surface-variant leading-relaxed">
          <div className="mb-1">
            <span className="text-success">genpulse@system</span>:<span className="text-primary">~/project-alpha</span>$ npm run build
          </div>
          <div className="text-outline mb-4">
            &gt; alpha-core@2.4.0 build<br/>&gt; tsc &amp;&amp; vite build
          </div>
          <div className="mb-1 text-primary-fixed-dim">[*] 初始化 Agent Reasoner...</div>
          <div className="mb-1 text-outline">[-] 解析目标文件的 AST...</div>
          <div className="mb-1 text-outline">
            [-] 通过 <span className="text-secondary">GPT-4-Turbo</span> 生成建议的修改...
          </div>
          <div className="mb-3 text-success">[+] 为 reasoning_engine.ts 成功生成差异</div>
          <div className="mb-1">
            <span className="text-success">genpulse@system</span>:<span className="text-primary">~/project-alpha</span>$ git status
          </div>
          <div className="text-outline">
            On branch feature/agent-reasoning-update<br/>Changes not staged for commit:<br/>  (use "git add &lt;file&gt;..." to update what will be committed)<br/>  (use "git restore &lt;file&gt;..." to discard changes in working directory)
          </div>
          <div className="text-error-container pl-4 mb-3">modified:   src/core/agents/reasoning_engine.ts</div>
          <div className="flex items-center">
            <span className="text-success">genpulse@system</span>:<span className="text-primary">~/project-alpha</span>$ 
            <span className="w-2 h-4 bg-primary ml-1 animate-[pulse_1s_step-end_infinite]"></span>
          </div>
          
          {/* 更多终端输出 */}
          <div className="mt-4 pt-4 border-t border-outline-variant/10">
            <div className="mb-1">
              <span className="text-success">genpulse@system</span>:<span className="text-primary">~/project-alpha</span>$ npm test
            </div>
            <div className="text-outline mb-2">
              &gt; alpha-core@2.4.0 test<br/>&gt; jest --coverage
            </div>
            <div className="text-outline mb-1">PASS src/core/agents/reasoning_engine.test.ts</div>
            <div className="text-success mb-1">✓ analyzeContext function (42ms)</div>
            <div className="text-outline">Test Suites: 1 passed, 1 total</div>
          </div>
        </div>
      </div>

      {/* 自定义样式 - 已通过Tailwind类实现 */}
    </div>
  );
};

export default TerminalViewGenpulse;