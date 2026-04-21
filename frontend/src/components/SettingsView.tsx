import React from 'react';
import { 
  Palette, 
  Globe as Language, 
  Terminal, 
  Settings as SettingsIcon, 
  Search, 
  Bell, 
  HelpCircle, 
  Plus as Add, 
  Server as Dns, 
  Database as DataObject, 
  Brain as PsychologyAlt, 
  Sparkles as AutoAwesome, 
  Grid as GroupWork, 
  GitBranch as AccountTree, 
  Ruler as Architecture, 
  Layout as Web, 
  ArrowLeftRight as SwapHoriz, 
  Edit, 
  History, 
  Maximize as OpenInFull, 
  Save,
  Globe,
  Settings as Tune,
  Command,
  Layout,
  Cpu,
  ShieldCheck,
  BrainCircuit,
  Settings2,
  Trash2,
  CheckCircle2,
  Activity,
  Server,
  Layers,
  Network
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '../utils';

export default function SettingsView() {
  return (
    <main className="flex-1 overflow-y-auto p-6 md:p-8 lg:p-12 relative bg-background custom-scrollbar">
      {/* Page Header */}
      <div className="max-w-6xl mx-auto mb-12">
        <h2 className="text-4xl font-black tracking-tight text-on-surface mb-2">系统设置</h2>
        <p className="text-outline text-base">管理全局环境变量、模型路由及代理的基础行为。</p>
      </div>

      <div className="max-w-6xl mx-auto space-y-16 pb-24">
        {/* SECTION 1: General Settings */}
        <section>
          <h3 className="text-xl font-bold text-primary mb-6 flex items-center gap-2">
            <Tune size={20} />
            通用设置
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <SettingsCard icon={Palette} title="界面主题" description="确定全局界面的外观。">
              <div className="flex bg-surface-container-highest/50 p-1.5 rounded-xl border border-outline-variant/10">
                <button className="flex-1 py-2 text-xs font-bold rounded-lg text-outline hover:text-on-surface transition-all">浅色</button>
                <button className="flex-1 py-2 text-xs font-bold rounded-lg bg-surface-container shadow-lg text-primary border border-primary/20">深色</button>
                <button className="flex-1 py-2 text-xs font-bold rounded-lg text-outline hover:text-on-surface transition-all">跟随系统</button>
              </div>
            </SettingsCard>

            <SettingsCard icon={Globe} title="语言设置" description="设置偏好的界面语言。">
              <div className="relative group">
                <select className="w-full appearance-none bg-surface-container-highest/50 border border-outline-variant/10 text-sm text-on-surface font-semibold rounded-xl py-3 pl-4 pr-10 focus:ring-1 focus:ring-primary focus:border-primary outline-none cursor-pointer transition-all">
                  <option selected value="zh">简体中文</option>
                  <option value="en">English (US)</option>
                </select>
                <Tune className="absolute right-3 top-1/2 -translate-y-1/2 text-outline pointer-events-none" size={16} />
              </div>
            </SettingsCard>

            <SettingsCard icon={Terminal} title="系统日志级别" description="执行追踪的详细程度。">
              <div className="relative group">
                <select className="w-full appearance-none bg-surface-container-highest/50 border border-outline-variant/10 text-sm text-on-surface font-mono rounded-xl py-3 pl-4 pr-10 focus:ring-1 focus:ring-primary focus:border-primary outline-none cursor-pointer transition-all">
                  <option value="debug">DEBUG</option>
                  <option selected value="info">INFO</option>
                  <option value="warn">WARN</option>
                  <option value="error">ERROR</option>
                </select>
                <Tune className="absolute right-3 top-1/2 -translate-y-1/2 text-outline pointer-events-none" size={16} />
              </div>
            </SettingsCard>
          </div>
        </section>

        {/* SECTION 2: Model Providers */}
        <section>
          <div className="flex items-center justify-between mb-8">
            <h3 className="text-xl font-bold text-primary flex items-center gap-2">
              <Cpu size={20} />
              模型提供商配置
            </h3>
            <button className="text-xs font-bold text-primary hover:text-white transition-all flex items-center gap-1.5 uppercase tracking-widest">
              <Add size={14} /> 添加自定义提供商
            </button>
          </div>
          <div className="bg-surface-container-low/40 rounded-3xl border border-outline-variant/10 overflow-hidden backdrop-blur-xl">
            <ProviderRow name="OpenAI" endpoint="api.openai.com/v1" connected model="gpt-4-turbo-preview" icon={Settings2} />
            <div className="h-px bg-outline-variant/10 mx-6" />
            <ProviderRow name="Anthropic Claude" endpoint="api.anthropic.com" connected model="claude-3-opus-20240229" icon={BrainCircuit} />
            <div className="h-px bg-outline-variant/10 mx-6" />
            <ProviderRow name="Google Gemini" endpoint="generativelanguage.googleapis" model="gemini-1.5-pro" icon={AutoAwesome} />
          </div>
        </section>

        {/* SECTION 3: Agent Role Templates */}
        <section>
          <h3 className="text-xl font-bold text-primary mb-8 flex items-center gap-2">
            <ShieldCheck size={20} />
            全局 Agent 角色模板
          </h3>
          <div className="grid grid-cols-1 xl:grid-cols-2 gap-8">
            <AgentRoleCard 
              name="Orchestrator" 
              role="任务路由与分解" 
              model="gpt-4-turbo-preview" 
              active
              icon={AccountTree}
              prompt="你是 Orchestrator。你的主要目标是将复杂的用户请求分解为特定子代理可执行的独立步骤..." 
            />
            <AgentRoleCard 
              name="Architect" 
              role="系统设计与架构" 
              model="claude-3-opus-20240229" 
              active
              icon={Layout}
              prompt="作为首席架构师，评估提出的需求并设计稳健、可扩展的数据模型。优先考虑关注点分离和..." 
            />
          </div>
        </section>

        {/* SECTION 4: MCP Configuration */}
        <section>
          <h3 className="text-xl font-bold text-primary mb-8 flex items-center gap-2">
            <Server size={20} />
            MCP 协议配置
          </h3>
          <div className="bg-surface-container-low/40 rounded-3xl border border-outline-variant/10 overflow-hidden backdrop-blur-xl">
            <div className="p-6">
              <div className="flex items-center justify-between mb-6">
                <div>
                  <h4 className="text-lg font-bold text-on-surface">Model Context Protocol (MCP)</h4>
                  <p className="text-sm text-outline mt-1">连接外部工具和服务，扩展 Agent 能力</p>
                </div>
                <div className="flex items-center gap-2">
                  <span className="h-2 w-2 rounded-full bg-success animate-pulse" />
                  <span className="text-[10px] font-bold uppercase tracking-widest text-success">已集成</span>
                </div>
              </div>
              
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-surface-container-high rounded-2xl p-5">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="p-2 rounded-xl bg-surface-container-lowest text-primary">
                      <Server size={20} />
                    </div>
                    <div>
                      <h5 className="font-bold text-on-surface">服务器管理</h5>
                      <p className="text-xs text-outline">配置 MCP 服务器连接</p>
                    </div>
                  </div>
                  <p className="text-sm text-on-surface-variant mb-4">
                    管理本地和远程 MCP 服务器，支持 stdio 和 SSE 连接方式
                  </p>
                  <button 
                    onClick={() => window.location.hash = '#/mcp-config'}
                    className="w-full bg-primary-container text-on-primary-container hover:brightness-110 rounded-xl py-2.5 text-sm font-bold transition-all flex items-center justify-center gap-2"
                  >
                    <Settings2 size={16} />
                    管理服务器
                  </button>
                </div>

                <div className="bg-surface-container-high rounded-2xl p-5">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="p-2 rounded-xl bg-surface-container-lowest text-success">
                      <Layers size={20} />
                    </div>
                    <div>
                      <h5 className="font-bold text-on-surface">工具发现</h5>
                      <p className="text-xs text-outline">自动发现可用工具</p>
                    </div>
                  </div>
                  <p className="text-sm text-on-surface-variant mb-4">
                    动态发现 MCP 服务器提供的工具，自动注册到工具注册表
                  </p>
                  <div className="flex items-center justify-between text-xs text-outline">
                    <span>发现间隔</span>
                    <span className="font-mono font-bold">60 秒</span>
                  </div>
                </div>

                <div className="bg-surface-container-high rounded-2xl p-5">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="p-2 rounded-xl bg-surface-container-lowest text-warning">
                      <Network size={20} />
                    </div>
                    <div>
                      <h5 className="font-bold text-on-surface">连接状态</h5>
                      <p className="text-xs text-outline">监控服务器连接</p>
                    </div>
                  </div>
                  <p className="text-sm text-on-surface-variant mb-4">
                    实时监控 MCP 服务器连接状态，支持自动重连和错误处理
                  </p>
                  <div className="flex items-center justify-between text-xs text-outline">
                    <span>活动连接</span>
                    <span className="font-mono font-bold text-success">4/4</span>
                  </div>
                </div>
              </div>

              <div className="mt-6 pt-6 border-t border-outline-variant/10">
                <div className="flex items-center justify-between">
                  <div>
                    <h5 className="font-bold text-on-surface">MCP 功能特性</h5>
                    <p className="text-sm text-outline mt-1">已实现完整的 MCP 协议支持</p>
                  </div>
                  <div className="flex items-center gap-4">
                    <div className="text-center">
                      <div className="text-2xl font-black text-primary">4</div>
                      <div className="text-[10px] uppercase tracking-widest text-outline">服务器</div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-black text-success">18</div>
                      <div className="text-[10px] uppercase tracking-widest text-outline">工具</div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-black text-warning">100%</div>
                      <div className="text-[10px] uppercase tracking-widest text-outline">可用性</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>

      {/* Floating Save Action */}
      <div className="fixed bottom-10 right-10 z-50">
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-surface-variant/90 backdrop-blur-3xl p-2 rounded-2xl shadow-2xl border border-white/10 flex items-center gap-3"
        >
          <button className="bg-surface-container-lowest text-on-surface hover:bg-surface-bright rounded-xl px-5 py-2.5 text-sm font-bold transition-all border border-outline-variant/10">放弃更改</button>
          <button className="bg-primary-container text-on-primary-container hover:brightness-110 rounded-xl px-7 py-2.5 text-sm font-bold shadow-xl shadow-primary/20 transition-all flex items-center gap-2">
            <Save size={18} />
            保存配置
          </button>
        </motion.div>
      </div>
    </main>
  );
}

function SettingsCard({ icon: Icon, title, description, children }: any) {
  return (
    <div className="bg-surface-container-low rounded-3xl p-6 flex flex-col justify-between relative overflow-hidden group border border-outline-variant/5 hover:border-outline-variant/20 transition-all shadow-sm">
      <div className="absolute inset-0 bg-gradient-to-br from-white/[0.02] to-transparent opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none" />
      <div className="mb-8">
        <div className="w-10 h-10 rounded-xl bg-surface-container-highest/50 flex items-center justify-center text-outline mb-4 group-hover:text-primary transition-all">
          <Icon size={20} />
        </div>
        <h4 className="text-base font-bold text-on-surface mb-1">{title}</h4>
        <p className="text-xs text-outline leading-relaxed">{description}</p>
      </div>
      {children}
    </div>
  );
}

function ProviderRow({ name, endpoint, connected, model, icon: Icon }: any) {
  return (
    <div className="flex items-center justify-between p-6 hover:bg-white/[0.02] transition-colors group">
      <div className="flex items-center gap-5 w-1/3">
        <div className={cn(
          "w-12 h-12 rounded-2xl flex items-center justify-center shadow-inner border transition-all",
          connected ? "bg-surface-container border-primary/20 text-primary" : "bg-surface-container-highest/30 border-outline-variant/10 text-outline opacity-60"
        )}>
          <Icon size={24} />
        </div>
        <div>
          <h4 className={cn("text-base font-bold", !connected && "opacity-60")}>{name}</h4>
          <p className="text-[10px] font-mono text-outline uppercase tracking-widest mt-1 opacity-70">{endpoint}</p>
        </div>
      </div>
      
      <div className="flex items-center gap-3 w-1/4">
        <span className={cn(
          "flex h-2.5 w-2.5 rounded-full shadow-lg",
          connected ? "bg-success shadow-success/40" : "bg-error shadow-error/40"
        )} />
        <span className={cn(
          "text-[10px] font-bold uppercase tracking-widest",
          connected ? "text-success" : "text-error"
        )}>
          {connected ? "已连接" : "未配置"}
        </span>
      </div>

      <div className="w-1/4">
        <div className={cn(
          "bg-surface-container-highest px-4 py-2 rounded-xl text-xs font-mono font-bold text-on-surface-variant flex items-center justify-between border border-outline-variant/10",
          !connected && "opacity-40"
        )}>
          {model}
        </div>
      </div>

      <div className="flex justify-end w-[15%]">
        <button className={cn(
          "text-[10px] font-bold uppercase tracking-widest px-4 py-2 rounded-lg transition-all",
          connected 
            ? "text-primary hover:bg-primary/10" 
            : "bg-primary-container text-on-primary-container hover:brightness-110 shadow-lg shadow-primary/20"
        )}>
          {connected ? "测试连接" : "配置密钥"}
        </button>
      </div>
    </div>
  );
}

function AgentRoleCard({ name, role, model, active, icon: Icon, prompt }: any) {
  return (
    <div className="bg-surface-container-high rounded-3xl p-7 relative overflow-hidden group shadow-2xl border border-outline-variant/10">
      <div className="absolute -top-10 -right-10 w-40 h-40 bg-primary/5 rounded-full blur-3xl group-hover:bg-primary/10 transition-all duration-700 pointer-events-none" />
      
      <div className="flex justify-between items-start mb-8 relative z-10">
        <div className="flex items-center gap-4">
          <div className="p-3 rounded-2xl bg-surface-container-lowest text-primary shadow-inner">
            <Icon size={24} />
          </div>
          <div>
            <h4 className="text-lg font-black text-on-surface">{name}</h4>
            <p className="text-xs text-outline font-medium">{role}</p>
          </div>
        </div>
        {active && (
          <div className="flex items-center gap-2 bg-surface-container-lowest/80 px-3 py-1 rounded-full border border-primary/20">
            <span className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse" />
            <span className="text-[10px] font-mono font-black text-primary uppercase tracking-widest">活跃</span>
          </div>
        )}
      </div>

      <div className="space-y-6 relative z-10">
        <div>
          <label className="text-[10px] uppercase font-black text-outline tracking-widest block mb-2">分配模型</label>
          <div className="bg-surface-container-lowest/50 rounded-xl px-4 py-3 flex items-center justify-between border border-outline-variant/10">
            <span className="font-mono text-sm text-primary font-bold">{model}</span>
            <Tune size={16} className="text-outline cursor-pointer hover:text-on-surface transition-colors" />
          </div>
        </div>
        <div>
          <div className="flex justify-between items-center mb-2">
            <label className="text-[10px] uppercase font-black text-outline tracking-widest block">系统提示词</label>
            <Edit size={14} className="text-outline cursor-pointer hover:text-primary transition-colors" />
          </div>
          <div className="bg-surface-container-lowest/50 rounded-2xl p-4 h-28 overflow-hidden relative border border-outline-variant/5">
            <p className="font-mono text-xs text-on-surface-variant leading-loose opacity-80 italic">
              {prompt}
            </p>
            <div className="absolute bottom-0 left-0 right-0 h-10 bg-gradient-to-t from-surface-container-lowest to-transparent" />
          </div>
        </div>
      </div>
    </div>
  );
}
