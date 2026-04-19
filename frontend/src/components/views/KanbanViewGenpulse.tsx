import React from 'react';

const KanbanViewGenpulse: React.FC = () => {
  return (
    <div className="flex-1 flex flex-col h-full overflow-hidden">
      {/* 页面标题栏 - 简化版本，因为TopAppBar已经提供主要导航 */}
      <div className="w-full flex justify-between items-center px-6 py-4 bg-surface-container-low z-20 border-b border-outline-variant/10">
        <div>
          <h2 className="text-2xl font-bold tracking-tight text-on-surface font-headline">执行看板</h2>
          <p className="text-sm text-outline font-body mt-1">Pipeline Alpha-7 运行中</p>
        </div>
        <div className="flex items-center space-x-6">
          <div className="flex items-center space-x-2 text-sm text-outline">
            <span className="material-symbols-outlined text-[1.1rem]">schedule</span>
            <span className="font-code text-xs">最后同步: 2分钟前</span>
          </div>
          <div className="flex items-center space-x-4">
            <button className="p-2 rounded-lg text-outline hover:bg-surface-variant hover:text-primary transition-all duration-200">
              <span className="material-symbols-outlined">search</span>
            </button>
            <button className="p-2 rounded-lg text-outline hover:bg-surface-variant hover:text-primary transition-all duration-200">
              <span className="material-symbols-outlined">filter_list</span>
            </button>
            <button className="bg-gradient-to-r from-primary-container to-inverse-primary px-4 py-2 rounded-lg text-on-primary-container text-sm font-medium hover:brightness-110 transition-all active:scale-95 shadow-lg shadow-primary-container/20">
              部署流水线
            </button>
          </div>
        </div>
      </div>

      {/* 看板画布 */}
      <main className="flex-1 p-6 overflow-x-auto flex space-x-6 items-start pb-12">
        {/* 协调者列 */}
        <div className="w-[320px] shrink-0 flex flex-col max-h-full">
          <div className="flex items-center justify-between mb-4 px-2">
            <div className="flex items-center space-x-2">
              <span className="material-symbols-outlined text-primary text-[1.2rem]">psychology</span>
              <h3 className="text-sm font-medium uppercase tracking-widest text-on-surface font-label">协调者</h3>
              <span className="bg-surface-container-highest text-outline text-xs px-2 py-0.5 rounded-full font-code">2</span>
            </div>
            <button className="text-outline hover:text-primary transition-colors">
              <span className="material-symbols-outlined text-[1.2rem]">add</span>
            </button>
          </div>
          <div className="flex-1 overflow-y-auto space-y-4 pr-2 pb-4">
            {/* 卡片 1 */}
            <div className="bg-surface-container-high rounded-lg p-5 ghost-border relative group hover:bg-surface-bright transition-colors duration-200 cursor-pointer">
              <div className="flex justify-between items-start mb-3">
                <span className="text-[0.65rem] uppercase tracking-wider font-semibold text-primary bg-primary/10 px-2 py-1 rounded">进行中</span>
                <img 
                  alt="Agent Icon" 
                  className="w-6 h-6 rounded-full object-cover opacity-80 group-hover:opacity-100 transition-opacity"
                  src="https://lh3.googleusercontent.com/aida-public/AB6AXuC9DEocclmBdS-GR6YVlqSu_mwzNIEbFjCsOQSvs-60Q_1ofdXAr1ECQBLkyr_op1he15PffgtYBSQFhzl98unNPq2JnBerTNH75x2tdVrADo6tyWWtgPY7PfQZle-ZQ-f8em1pBEU333jDHM77i8wnAvMy-gF236jFIU35t5tAX3VOnY6SCHqtuzpirL1phT33_tVYheWjIP-nWhhQ3G-Irar8zvM2nBNp6WH_EB7uPY4WoXQ8uY_lQH2-4Am5jh4uk3XYvrdqFtw9"
                />
              </div>
              <h4 className="text-base font-semibold text-on-surface mb-2 font-headline">分析需求</h4>
              <p className="text-xs text-outline font-body leading-relaxed mb-4 line-clamp-2">解析用户意图，将初始提示分解为离散的微服务架构。</p>
              <div className="flex items-center justify-between text-xs text-outline">
                <div className="flex items-center space-x-1 font-code">
                  <span className="material-symbols-outlined text-[1rem]">terminal</span>
                  <span>CMD-01</span>
                </div>
                <div className="flex items-center space-x-1">
                  <div className="pulse-dot"></div>
                  <span className="ml-2 text-primary font-medium">活跃</span>
                </div>
              </div>
            </div>

            {/* 卡片 2 */}
            <div className="bg-surface-container-high rounded-lg p-5 ghost-border relative group hover:bg-surface-bright transition-colors duration-200 cursor-pointer">
              <div className="flex justify-between items-start mb-3">
                <span className="text-[0.65rem] uppercase tracking-wider font-semibold text-outline bg-surface-container-highest px-2 py-1 rounded">待处理</span>
                <img 
                  alt="Agent Icon" 
                  className="w-6 h-6 rounded-full object-cover opacity-60"
                  src="https://lh3.googleusercontent.com/aida-public/AB6AXuB2FxDsOAUm5lD9hp_wd1wAu2XjKo0omT2_erq61ROpIwSXnjC6WHChTfDniDSk97CxXC7N5EIDXAoky3___p_98rPWzYaXV6hnoh17LfoKqlm2nYclXtttoTkShEgsZ14yfY9a1eYv8Ofi7HnJRPGpd6yWWXVsvhLy08TS0CyoHPY3QOfC5VsuRpVMbLSCPk_0NM5sNuFvvIJWDukV5f6jrsNa0OLt5n1kESmdHWJ8I7ACjeL9gUj3NPG358HcSXzn4KKSAN7VxYV6"
                />
              </div>
              <h4 className="text-base font-semibold text-on-surface mb-2 font-headline">分配资源</h4>
              <p className="text-xs text-outline font-body leading-relaxed mb-4 line-clamp-2">确定计算需求，并将虚拟代理分配给开发容器。</p>
              <div className="flex items-center justify-between text-xs text-outline">
                <div className="flex items-center space-x-1 font-code">
                  <span className="material-symbols-outlined text-[1rem]">terminal</span>
                  <span>CMD-02</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* 产品经理列 */}
        <div className="w-[320px] shrink-0 flex flex-col max-h-full">
          <div className="flex items-center justify-between mb-4 px-2">
            <div className="flex items-center space-x-2">
              <span className="material-symbols-outlined text-primary-fixed-dim text-[1.2rem]">person</span>
              <h3 className="text-sm font-medium uppercase tracking-widest text-on-surface font-label">产品经理</h3>
              <span className="bg-surface-container-highest text-outline text-xs px-2 py-0.5 rounded-full font-code">1</span>
            </div>
            <button className="text-outline hover:text-primary transition-colors">
              <span className="material-symbols-outlined text-[1.2rem]">add</span>
            </button>
          </div>
          <div className="flex-1 overflow-y-auto space-y-4 pr-2 pb-4">
            {/* 卡片 3 */}
            <div className="bg-surface-container-high rounded-lg p-5 ghost-border relative group hover:bg-surface-bright transition-colors duration-200 cursor-pointer">
              <div className="flex justify-between items-start mb-3">
                <span className="text-[0.65rem] uppercase tracking-wider font-semibold text-primary bg-primary/10 px-2 py-1 rounded">进行中</span>
                <img 
                  alt="Agent Icon" 
                  className="w-6 h-6 rounded-full object-cover opacity-80 group-hover:opacity-100 transition-opacity"
                  src="https://lh3.googleusercontent.com/aida-public/AB6AXuC9DEocclmBdS-GR6YVlqSu_mwzNIEbFjCsOQSvs-60Q_1ofdXAr1ECQBLkyr_op1he15PffgtYBSQFhzl98unNPq2JnBerTNH75x2tdVrADo6tyWWtgPY7PfQZle-ZQ-f8em1pBEU333jDHM77i8wnAvMy-gF236jFIU35t5tAX3VOnY6SCHqtuzpirL1phT33_tVYheWjIP-nWhhQ3G-Irar8zvM2nBNp6WH_EB7uPY4WoXQ8uY_lQH2-4Am5jh4uk3XYvrdqFtw9"
                />
              </div>
              <h4 className="text-base font-semibold text-on-surface mb-2 font-headline">创建用户故事</h4>
              <p className="text-xs text-outline font-body leading-relaxed mb-4 line-clamp-2">将需求转化为可执行的用户故事，并确定优先级。</p>
              <div className="flex items-center justify-between text-xs text-outline">
                <div className="flex items-center space-x-1 font-code">
                  <span className="material-symbols-outlined text-[1rem]">story</span>
                  <span>US-01</span>
                </div>
                <div className="flex items-center space-x-1">
                  <div className="w-2 h-2 rounded-full bg-primary relative">
                    <div className="absolute -inset-1 rounded-full bg-primary opacity-40 animate-ping"></div>
                  </div>
                  <span className="ml-2 text-primary font-medium">活跃</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* 架构师列 */}
        <div className="w-[320px] shrink-0 flex flex-col max-h-full">
          <div className="flex items-center justify-between mb-4 px-2">
            <div className="flex items-center space-x-2">
              <span className="material-symbols-outlined text-tertiary text-[1.2rem]">architecture</span>
              <h3 className="text-sm font-medium uppercase tracking-widest text-on-surface font-label">架构师</h3>
              <span className="bg-surface-container-highest text-outline text-xs px-2 py-0.5 rounded-full font-code">1</span>
            </div>
            <button className="text-outline hover:text-primary transition-colors">
              <span className="material-symbols-outlined text-[1.2rem]">add</span>
            </button>
          </div>
          <div className="flex-1 overflow-y-auto space-y-4 pr-2 pb-4">
            {/* 卡片 4 */}
            <div className="bg-surface-container-high rounded-lg p-5 ghost-border relative group hover:bg-surface-bright transition-colors duration-200 cursor-pointer">
              <div className="flex justify-between items-start mb-3">
                <span className="text-[0.65rem] uppercase tracking-wider font-semibold text-tertiary bg-tertiary/10 px-2 py-1 rounded">进行中</span>
                <img 
                  alt="Agent Icon" 
                  className="w-6 h-6 rounded-full object-cover opacity-80 group-hover:opacity-100 transition-opacity"
                  src="https://lh3.googleusercontent.com/aida-public/AB6AXuBBzd0c5qdfyDrPDv65gjZAblfipGGZYWVi0puMuYLV86wZHkQELAXy90W_TrHoMhAkaDVxcvQQ0Q5jyquTmwY-4GAAYGWQzfyQ14a9ZXtOr9GsyDlDr3XX6gytIMzXQvctblOnvBmGLvVWiHYj-zxQMGL1PFwKy_pjv0gKTZtjL1LeAuvqJGJd67fkuz6cgMtI8RDt8mn3ECXDWdaiJMjKizXaOUvZuvZnkTm4DqjwxRXPrHM6800N-K0l1d-3veuZ_7G7bx_v6ViN"
                />
              </div>
              <h4 className="text-base font-semibold text-on-surface mb-2 font-headline">设计数据库架构</h4>
              <p className="text-xs text-outline font-body leading-relaxed mb-4 line-clamp-2">为用户配置文件和认证令牌设计关系表。</p>
              <div className="flex items-center justify-between text-xs text-outline">
                <div className="flex items-center space-x-1 font-code">
                  <span className="material-symbols-outlined text-[1rem]">database</span>
                  <span>DB-01</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* 开发者列 */}
        <div className="w-[320px] shrink-0 flex flex-col max-h-full">
          <div className="flex items-center justify-between mb-4 px-2">
            <div className="flex items-center space-x-2">
              <span className="material-symbols-outlined text-secondary text-[1.2rem]">code</span>
              <h3 className="text-sm font-medium uppercase tracking-widest text-on-surface font-label">开发者</h3>
              <span className="bg-surface-container-highest text-outline text-xs px-2 py-0.5 rounded-full font-code">3</span>
            </div>
            <button className="text-outline hover:text-primary transition-colors">
              <span className="material-symbols-outlined text-[1.2rem]">add</span>
            </button>
          </div>
          <div className="flex-1 overflow-y-auto space-y-4 pr-2 pb-4">
            {/* 卡片 5 - 高亮卡片 */}
            <div className="liquid-glass rounded-lg p-5 ghost-border relative group hover:bg-surface-bright/80 transition-colors duration-200 cursor-pointer ambient-shadow border-l-2 border-primary">
              <div className="flex justify-between items-start mb-3">
                <span className="text-[0.65rem] uppercase tracking-wider font-semibold text-primary bg-primary/10 px-2 py-1 rounded">进行中</span>
                <div className="flex -space-x-2">
                  <img 
                    alt="Agent Icon" 
                    className="w-6 h-6 rounded-full object-cover border border-surface-variant"
                    src="https://lh3.googleusercontent.com/aida-public/AB6AXuBZQgleY3VkBjdiMGQ9sdcbJRSnZCT2XxDAKiRoLTZCqjXSQoEf8DqiZcPt9CZtjfwky-13Yv4gMYNLC_jZqD-0XfGAdwRFTOaf3TQKZNGO2bj1AFt1nW_mQ1T9QnHSFo6HmL_XWEmhTE4heLi2GeFYL1gzU72zph9n9iuXSaGEzstpKOS2rAR66hisMseknLBF1IM_VBQYfW2av36cS8awcDt9GDbgnrUiOmPh0t7YTi0ZGSnFRspFWGEb4utePuSfGmh7hZvBkucj"
                  />
                  <img 
                    alt="Agent Icon" 
                    className="w-6 h-6 rounded-full object-cover border border-surface-variant"
                    src="https://lh3.googleusercontent.com/aida-public/AB6AXuCx6zC6Fli2k3lNRTwmpGO9gqddMcuZH0LxdW_O7SdibSa6z7o9MuXq_iQm6d74chtIQdSnVU7FMHGpe84El7s2HK6WDf_gLbubdgdlRjKPNJdRK9dib8HOQ7ST8cn6jOWlvdSRT20F64TCysfYrcPlMnpA8SVZijmB-MORHWLf6czlk9_u21L64qPyMlEM6Ub_cnOtmqvJ9d9XJMyaCgS0_ik3h2dxQOvlb-9UgrZjo5hmxnyh12BCf7sAW4NvmbY_4rF26jodBC4C"
                  />
                </div>
              </div>
              <h4 className="text-base font-semibold text-on-surface mb-2 font-headline">API 路由生成</h4>
              <p className="text-xs text-on-surface-variant font-body leading-relaxed mb-4 line-clamp-2">基于架构师设计的架构，为认证和用户数据检索搭建 Express.js 端点。</p>
              <div className="bg-surface-container-lowest rounded p-2 mb-4">
                <code className="text-[0.65rem] font-code text-primary-fixed block">POST /api/v1/auth/login</code>
                <code className="text-[0.65rem] font-code text-primary-fixed block">GET /api/v1/users/:id</code>
              </div>
              <div className="flex items-center justify-between text-xs text-outline">
                <div className="flex items-center space-x-1 font-code">
                  <span className="material-symbols-outlined text-[1rem]">api</span>
                  <span>BE-04</span>
                </div>
                <div className="flex items-center space-x-1">
                  <div className="w-2 h-2 rounded-full bg-primary relative">
                    <div className="absolute -inset-1 rounded-full bg-primary opacity-40 animate-ping"></div>
                  </div>
                  <span className="ml-2 text-primary font-medium">活跃</span>
                </div>
              </div>
            </div>

            {/* 卡片 6 */}
            <div className="bg-surface-container-high rounded-lg p-5 ghost-border relative group hover:bg-surface-bright transition-colors duration-200 cursor-pointer">
              <div className="flex justify-between items-start mb-3">
                <span className="text-[0.65rem] uppercase tracking-wider font-semibold text-success bg-success/10 px-2 py-1 rounded">已完成</span>
                <img 
                  alt="Agent Icon" 
                  className="w-6 h-6 rounded-full object-cover opacity-60"
                  src="https://lh3.googleusercontent.com/aida-public/AB6AXuB3usKmDTsEUd7aGwNOge2a7RASrLd93JhOiC_A-fU3LrwHAHV4dzXvAuZo3_1r34UREGmuc-h8E5NDepquSVuNZEzbImTtRGWQWNsCIgTmYbmfIGzjZCptl7eDHsdm5ibgvrX2HObVzhmBdFQoq8LQOBk1tYzVnHLJHtWBAq-RamIX_CyiWR9fjfmNsY5bSEgXDVSAYzVWa87n5ku2_sm7wZ9JISQFSGA4bY6y0lSU3oYRf5v1Yt7ay3mcw1h8hawTfKq7A4SBc65H"
                />
              </div>
              <h4 className="text-base font-semibold text-on-surface mb-2 font-headline line-through opacity-70">初始化仓库</h4>
              <p className="text-xs text-outline font-body leading-relaxed mb-4 line-clamp-2 opacity-70">使用 Tailwind CSS 和严格的 TypeScript 配置搭建 Next.js 环境。</p>
              <div className="flex items-center justify-between text-xs text-outline opacity-70">
                <div className="flex items-center space-x-1 font-code">
                  <span className="material-symbols-outlined text-[1rem]">folder</span>
                  <span>FE-01</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* QA列 */}
        <div className="w-[320px] shrink-0 flex flex-col max-h-full">
          <div className="flex items-center justify-between mb-4 px-2">
            <div className="flex items-center space-x-2">
              <span className="material-symbols-outlined text-error text-[1.2rem]">fact_check</span>
              <h3 className="text-sm font-medium uppercase tracking-widest text-on-surface font-label">质量保证</h3>
              <span className="bg-surface-container-highest text-outline text-xs px-2 py-0.5 rounded-full font-code">1</span>
            </div>
            <button className="text-outline hover:text-primary transition-colors">
              <span className="material-symbols-outlined text-[1.2rem]">add</span>
            </button>
          </div>
          <div className="flex-1 overflow-y-auto space-y-4 pr-2 pb-4">
            {/* 卡片 7 */}
            <div className="bg-surface-container-high rounded-lg p-5 ghost-border relative group hover:bg-surface-bright transition-colors duration-200 cursor-pointer">
              <div className="flex justify-between items-start mb-3">
                <span className="text-[0.65rem] uppercase tracking-wider font-semibold text-outline bg-surface-container-highest px-2 py-1 rounded">待处理</span>
                <img 
                  alt="Agent Icon" 
                  className="w-6 h-6 rounded-full object-cover opacity-60"
                  src="https://lh3.googleusercontent.com/aida-public/AB6AXuCXtzWeijKxea2pHe1vp5ItwIsX9jUhV1xqCDo6KjyvcrwcPVS-YQpnZNmjCvlBIN4rPq1XaeJetkFBIUITPLPy_dtTq9vG-6_LqNNal_xUvFS2Bv4tXVcWbU1EaoLvp6k84VTP0891YIe70PfKCp_DYTdVqPImF0YeRX3w4CHN55k9C9GKtL38p6MuC4QfO0x1cHHFE8CNhudTibRt4KaINSyT-QSxDw0qOXH0HWGxNxRld7AJTQmRMWl4wmfIfGiEHjMkYi51yqQt"
                />
              </div>
              <h4 className="text-base font-semibold text-on-surface mb-2 font-headline">单元测试数据库连接</h4>
              <p className="text-xs text-outline font-body leading-relaxed mb-4 line-clamp-2">验证 ORM 模型与架构师设计的架构是否一致，并测试插入操作的边界情况。</p>
              <div className="flex items-center justify-between text-xs text-outline">
                <div className="flex items-center space-x-1 font-code">
                  <span className="material-symbols-outlined text-[1rem]">bug_report</span>
                  <span>QA-01</span>
                </div>
                <span className="text-[0.65rem] px-2 py-0.5 bg-error-container text-error rounded font-medium">受阻</span>
              </div>
            </div>
          </div>
        </div>

        {/* 运维列 */}
        <div className="w-[320px] shrink-0 flex flex-col max-h-full">
          <div className="flex items-center justify-between mb-4 px-2">
            <div className="flex items-center space-x-2">
              <span className="material-symbols-outlined text-secondary-fixed-dim text-[1.2rem]">cloud</span>
              <h3 className="text-sm font-medium uppercase tracking-widest text-on-surface font-label">运维</h3>
              <span className="bg-surface-container-highest text-outline text-xs px-2 py-0.5 rounded-full font-code">2</span>
            </div>
            <button className="text-outline hover:text-primary transition-colors">
              <span className="material-symbols-outlined text-[1.2rem]">add</span>
            </button>
          </div>
          <div className="flex-1 overflow-y-auto space-y-4 pr-2 pb-4">
            {/* 卡片 8 */}
            <div className="bg-surface-container-high rounded-lg p-5 ghost-border relative group hover:bg-surface-bright transition-colors duration-200 cursor-pointer">
              <div className="flex justify-between items-start mb-3">
                <span className="text-[0.65rem] uppercase tracking-wider font-semibold text-warning bg-warning/10 px-2 py-1 rounded">配置中</span>
                <img 
                  alt="Agent Icon" 
                  className="w-6 h-6 rounded-full object-cover opacity-80 group-hover:opacity-100 transition-opacity"
                  src="https://lh3.googleusercontent.com/aida-public/AB6AXuC9DEocclmBdS-GR6YVlqSu_mwzNIEbFjCsOQSvs-60Q_1ofdXAr1ECQBLkyr_op1he15PffgtYBSQFhzl98unNPq2JnBerTNH75x2tdVrADo6tyWWtgPY7PfQZle-ZQ-f8em1pBEU333jDHM77i8wnAvMy-gF236jFIU35t5tAX3VOnY6SCHqtuzpirL1phT33_tVYheWjIP-nWhhQ3G-Irar8zvM2nBNp6WH_EB7uPY4WoXQ8uY_lQH2-4Am5jh4uk3XYvrdqFtw9"
                />
              </div>
              <h4 className="text-base font-semibold text-on-surface mb-2 font-headline">配置 CI/CD 流水线</h4>
              <p className="text-xs text-outline font-body leading-relaxed mb-4 line-clamp-2">为自动构建、测试和部署设置 Jenkins/GitHub Actions 工作流。</p>
              <div className="flex items-center justify-between text-xs text-outline">
                <div className="flex items-center space-x-1 font-code">
                  <span className="material-symbols-outlined text-[1rem]">settings</span>
                  <span>OPS-01</span>
                </div>
                <div className="flex items-center space-x-1">
                  <div className="w-2 h-2 rounded-full bg-primary relative">
                    <div className="absolute -inset-1 rounded-full bg-primary opacity-40 animate-ping"></div>
                  </div>
                  <span className="ml-2 text-primary font-medium">活跃</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>

      {/* 自定义样式 - 已通过Tailwind类实现 */}
    </div>
  );
};

export default KanbanViewGenpulse;