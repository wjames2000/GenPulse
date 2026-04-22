package genkit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"GenPulse/internal/agents"
	genkitconfig "GenPulse/internal/genkit/config"
	"GenPulse/internal/genkit/flows"
	"GenPulse/internal/genkit/models"
	"GenPulse/internal/genkit/tools"
	mcpconfig "GenPulse/internal/mcp/config"
	"GenPulse/internal/mcp/discovery"
	"GenPulse/internal/mcp/host"
	"GenPulse/internal/memory"
	"GenPulse/internal/skills"
	"GenPulse/internal/utils"
)

// GenkitManager 管理Genkit运行时
type GenkitManager struct {
	ctx              context.Context
	config           *genkitconfig.AppConfig
	initialized      bool
	modelAdapter     *models.UnifiedModelAdapter
	toolRegistry     *tools.ToolRegistry
	flowEngine       *flows.FlowEngine
	agentManager     *agents.AgentManager
	skillManager     *skills.SkillManager
	memoryManager    *memory.SearchEngine
	mcpHost          *host.MCPHost
	mcpConfig        *mcpconfig.MCPConfigManager
	toolDiscovery    *discovery.ToolDiscovery
	StartupOptimizer *StartupOptimizer
	phase2Ready      bool
	phase3Ready      bool
	// genkit     interface{} // 暂时用interface{}，等确定具体类型后替换
}

// NewGenkitManager 创建新的Genkit管理器
func NewGenkitManager() *GenkitManager {
	return &GenkitManager{
		initialized:      false,
		StartupOptimizer: GetStartupOptimizer(),
	}
}

// Initialize 初始化Genkit运行时（带启动阶段优化）
func (gm *GenkitManager) Initialize(ctx context.Context) error {
	if gm.initialized {
		return nil
	}

	gm.ctx = ctx

	// 获取配置
	cfg := genkitconfig.GetConfig()
	gm.config = cfg

	utils.Info("初始化Genkit运行时（启动优化模式）...")

	// 预加载配置文件
	if gm.StartupOptimizer.preloadConfig.Enabled {
		gm.StartupOptimizer.PreloadConfigFiles()
	}

	// ========== Phase 1: 关键服务 ==========
	phase1Err := RunPhaseWithTimeout("Phase 1: 关键服务", Phase1Critical, gm.StartupOptimizer.GetPhase1Timeout(), gm.StartupOptimizer.metrics, func() error {
		utils.Info("Phase 1: 初始化关键服务...")

		// 配置模型提供商（核心依赖）
		if err := gm.configureModelProviders(); err != nil {
			return fmt.Errorf("failed to configure model providers: %w", err)
		}

		// 初始化工具注册表（核心依赖）
		if err := gm.initToolRegistry(); err != nil {
			return fmt.Errorf("failed to initialize tool registry: %w", err)
		}

		return nil
	})
	if phase1Err != nil {
		utils.Error("Phase 1 初始化失败: %v", phase1Err)
	}

	// ========== Phase 2: 重要服务（并行） ==========
	var wg sync.WaitGroup
	var phase2Mu sync.Mutex
	var phase2Errors []error

	wg.Add(2)
	go func() {
		defer wg.Done()
		start := time.Now()
		err := gm.initFlowEngine()
		elapsed := time.Since(start)
		gm.StartupOptimizer.metrics.RecordPhase("Flow引擎", Phase2Important, err == nil, err, elapsed)
		if err != nil {
			phase2Mu.Lock()
			phase2Errors = append(phase2Errors, err)
			phase2Mu.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		start := time.Now()
		err := gm.initAgentManager()
		elapsed := time.Since(start)
		gm.StartupOptimizer.metrics.RecordPhase("Agent管理器", Phase2Important, err == nil, err, elapsed)
		if err != nil {
			phase2Mu.Lock()
			phase2Errors = append(phase2Errors, err)
			phase2Mu.Unlock()
		}
	}()

	wg.Wait()
	gm.phase2Ready = true

	for _, err := range phase2Errors {
		utils.Warn("Phase 2 部分初始化失败: %v", err)
	}

	// ========== Phase 3: 后台服务（延迟初始化） ==========
	go func() {
		start := time.Now()
		var phase3Wg sync.WaitGroup

		// 技能管理器（延迟）
		if !gm.StartupOptimizer.LazyInitSkills() {
			phase3Wg.Add(1)
			go func() {
				defer phase3Wg.Done()
				s := time.Now()
				err := gm.initSkillManager()
				gm.StartupOptimizer.metrics.RecordPhase("技能管理器", Phase3Background, err == nil, err, time.Since(s))
				if err != nil {
					utils.Warn("技能管理器初始化失败（Phase 3）: %v", err)
				}
			}()
		}

		// 记忆管理器（延迟）
		phase3Wg.Add(1)
		go func() {
			defer phase3Wg.Done()
			s := time.Now()
			err := gm.initMemoryManager()
			gm.StartupOptimizer.metrics.RecordPhase("记忆管理器", Phase3Background, err == nil, err, time.Since(s))
			if err != nil {
				utils.Warn("记忆管理器初始化失败（Phase 3）: %v", err)
			}
		}()

		// MCP功能（延迟）
		if !gm.StartupOptimizer.LazyInitMCP() {
			phase3Wg.Add(1)
			go func() {
				defer phase3Wg.Done()
				s := time.Now()
				err := gm.initMCP()
				gm.StartupOptimizer.metrics.RecordPhase("MCP功能", Phase3Background, err == nil, err, time.Since(s))
				if err != nil {
					utils.Warn("MCP初始化失败（Phase 3）: %v", err)
				}
			}()
		}

		phase3Wg.Wait()
		gm.phase3Ready = true
		elapsed := time.Since(start)
		gm.StartupOptimizer.metrics.RecordPhase("Phase 3: 后台服务", Phase3Background, true, nil, elapsed)

		gm.initialized = true
		utils.Info("Genkit运行时初始化完成（总耗时: %v）", elapsed)

		// 报告启动指标
		gm.StartupOptimizer.ReportMetrics()
	}()

	// 立即标记为已初始化（前端可交互）
	gm.initialized = true
	return nil
}

func (gm *GenkitManager) IsPhase2Ready() bool {
	return gm.phase2Ready
}

func (gm *GenkitManager) IsPhase3Ready() bool {
	return gm.phase3Ready
}

// InitializeSync 同步初始化（阻塞直到全部完成）
func (gm *GenkitManager) InitializeSync(ctx context.Context) error {
	if err := gm.Initialize(ctx); err != nil {
		return err
	}

	for !gm.phase3Ready {
		time.Sleep(50 * time.Millisecond)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
	return nil
}

// RunPhase2 手动触发Phase 2初始化（如果尚未完成）
func (gm *GenkitManager) RunPhase2() error {
	if gm.phase2Ready {
		return nil
	}

	var wg sync.WaitGroup
	var lastErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		if gm.flowEngine == nil {
			start := time.Now()
			err := gm.initFlowEngine()
			gm.StartupOptimizer.metrics.RecordPhase("Flow引擎(按需)", Phase2Important, err == nil, err, time.Since(start))
			if err != nil {
				lastErr = err
			}
		}
	}()
	go func() {
		defer wg.Done()
		if gm.agentManager == nil {
			start := time.Now()
			err := gm.initAgentManager()
			gm.StartupOptimizer.metrics.RecordPhase("Agent管理器(按需)", Phase2Important, err == nil, err, time.Since(start))
			if err != nil {
				lastErr = err
			}
		}
	}()

	wg.Wait()
	gm.phase2Ready = true
	return lastErr
}

// configureModelProviders 配置模型提供商
func (gm *GenkitManager) configureModelProviders() error {
	utils.Info("配置模型提供商...")

	// 这里可以根据配置动态加载不同的模型提供商
	// 目前先实现基础框架

	// 检查是否有API密钥配置
	// 注意：这里需要从配置管理器获取，而不是直接访问config
	configMgr := genkitconfig.GetGlobalConfig()
	if configMgr == nil {
		utils.Warn("配置管理器未初始化，跳过模型提供商配置")
		return nil
	}

	googleAPIKey := configMgr.GetAPIKey("google")
	openaiAPIKey := configMgr.GetAPIKey("openai")
	anthropicAPIKey := configMgr.GetAPIKey("anthropic")

	if googleAPIKey == "" && openaiAPIKey == "" && anthropicAPIKey == "" {
		utils.Warn("未配置任何API密钥，模型功能将受限")
		utils.Info("请在配置文件中配置API密钥:")
		utils.Info("  - google: Gemini/Vertex AI")
		utils.Info("  - openai: GPT系列")
		utils.Info("  - anthropic: Claude系列")
	}

	// 创建模型适配器工厂
	factory := &models.DefaultModelAdapterFactory{}

	// 创建统一模型适配器
	modelAdapter := models.NewUnifiedModelAdapter(factory)
	if modelAdapter == nil {
		return fmt.Errorf("failed to create model adapter")
	}

	// 注册可用的模型
	if googleAPIKey != "" {
		utils.Info("已配置Google API密钥，注册Gemini模型")

		geminiConfig := models.ModelConfig{
			Type:     models.ModelTypeGemini,
			Name:     "gemini-1.5-pro",
			Provider: "google",
			APIKey:   googleAPIKey,
		}

		if err := modelAdapter.RegisterModel(geminiConfig); err != nil {
			utils.Warn("注册Gemini模型失败: %v", err)
		}
	}

	if openaiAPIKey != "" {
		utils.Info("已配置OpenAI API密钥，注册GPT模型")

		gptConfig := models.ModelConfig{
			Type:     models.ModelTypeGPT,
			Name:     "gpt-4",
			Provider: "openai",
			APIKey:   openaiAPIKey,
		}

		if err := modelAdapter.RegisterModel(gptConfig); err != nil {
			utils.Warn("注册GPT模型失败: %v", err)
		}
	}

	if anthropicAPIKey != "" {
		utils.Info("已配置Anthropic API密钥，注册Claude模型")

		claudeConfig := models.ModelConfig{
			Type:     models.ModelTypeClaude,
			Name:     "claude-3-opus",
			Provider: "anthropic",
			APIKey:   anthropicAPIKey,
		}

		if err := modelAdapter.RegisterModel(claudeConfig); err != nil {
			utils.Warn("注册Claude模型失败: %v", err)
		}
	}

	// 注册Ollama模型（本地运行，无需API密钥）
	ollamaConfig := models.ModelConfig{
		Type:     models.ModelTypeOllama,
		Name:     "llama3",
		Provider: "ollama",
	}

	if err := modelAdapter.RegisterModel(ollamaConfig); err != nil {
		utils.Warn("注册Ollama模型失败: %v", err)
	}

	gm.modelAdapter = modelAdapter
	utils.Info("模型提供商配置完成，注册模型数量: %d", len(modelAdapter.ListModels()))
	return nil
}

// initToolRegistry 初始化工具注册表
func (gm *GenkitManager) initToolRegistry() error {
	utils.Info("初始化工具注册表...")

	// 获取全局工具注册表
	toolRegistry := tools.GetGlobalToolRegistry()
	if toolRegistry == nil {
		return fmt.Errorf("failed to get global tool registry")
	}

	// 创建工作区路径
	workspacePath := "/tmp/genpulse-workspace"
	if gm.config != nil && gm.config.WorkspacePath != "" {
		workspacePath = gm.config.WorkspacePath
	}

	// 注册文件系统工具
	fsTool, err := tools.NewFSTool(workspacePath)
	if err != nil {
		utils.Warn("Failed to create file system tool: %v", err)
	} else {
		if err := toolRegistry.RegisterTool(fsTool); err != nil {
			utils.Warn("Failed to register file system tool: %v", err)
		} else {
			utils.Info("已注册文件系统工具")
		}
	}

	// 注册Git工具
	gitTool, err := tools.NewGitTool(workspacePath)
	if err != nil {
		utils.Warn("Failed to create Git tool: %v", err)
	} else {
		if err := toolRegistry.RegisterTool(gitTool); err != nil {
			utils.Warn("Failed to register Git tool: %v", err)
		} else {
			utils.Info("已注册Git工具")
		}
	}

	// 注册Shell工具
	shellTool, err := tools.NewShellTool(workspacePath)
	if err != nil {
		utils.Warn("Failed to create Shell tool: %v", err)
	} else {
		if err := toolRegistry.RegisterTool(shellTool); err != nil {
			utils.Warn("Failed to register Shell tool: %v", err)
		} else {
			utils.Info("已注册Shell工具")
		}
	}

	// 注册项目管理工具
	projectTool, err := tools.NewProjectTool(workspacePath)
	if err != nil {
		utils.Warn("Failed to create Project tool: %v", err)
	} else {
		if err := toolRegistry.RegisterTool(projectTool); err != nil {
			utils.Warn("Failed to register Project tool: %v", err)
		} else {
			utils.Info("已注册项目管理工具")
		}
	}

	// 初始化所有已注册的工具
	toolRegistry.InitializeAllTools()

	// 获取工具统计信息
	stats := toolRegistry.GetToolStatistics()
	utils.Info("工具注册表初始化完成，注册工具数量: %d", stats["total_tools"])

	return nil
}

// initFlowEngine 初始化Flow引擎
func (gm *GenkitManager) initFlowEngine() error {
	utils.Info("初始化Flow引擎...")

	// 创建Flow引擎实例
	flowEngine := flows.NewFlowEngine(gm.modelAdapter, gm.toolRegistry)
	if flowEngine == nil {
		return fmt.Errorf("failed to create flow engine")
	}

	// Flow引擎创建成功
	utils.Info("Flow引擎创建成功")

	gm.flowEngine = flowEngine
	utils.Info("Flow引擎初始化完成")
	return nil
}

// initAgentManager 初始化Agent管理器
func (gm *GenkitManager) initAgentManager() error {
	utils.Info("初始化Agent管理器...")

	// 检查依赖
	if gm.modelAdapter == nil {
		return fmt.Errorf("model adapter is required for agent manager")
	}

	if gm.toolRegistry == nil {
		return fmt.Errorf("tool registry is required for agent manager")
	}

	// 创建Agent管理器
	agentManager := agents.NewAgentManager(gm.modelAdapter, gm.toolRegistry, gm.flowEngine)
	if agentManager == nil {
		return fmt.Errorf("failed to create agent manager")
	}

	// 初始化Agent管理器
	if err := agentManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize agent manager: %w", err)
	}

	gm.agentManager = agentManager
	utils.Info("Agent管理器初始化完成")
	return nil
}

// GetGenkit 获取Genkit实例
func (gm *GenkitManager) GetGenkit() interface{} {
	if !gm.initialized {
		utils.Warn("Genkit管理器未初始化")
		return nil
	}
	return nil // TODO: 返回实际的genkit实例
}

// GetModelAdapter 获取模型适配器
func (gm *GenkitManager) GetModelAdapter() *models.UnifiedModelAdapter {
	return gm.modelAdapter
}

// GetToolRegistry 获取工具注册表
func (gm *GenkitManager) GetToolRegistry() *tools.ToolRegistry {
	return gm.toolRegistry
}

// GetFlowEngine 获取Flow引擎
func (gm *GenkitManager) GetFlowEngine() *flows.FlowEngine {
	return gm.flowEngine
}

// GetAgentManager 获取Agent管理器
func (gm *GenkitManager) GetAgentManager() *agents.AgentManager {
	return gm.agentManager
}

// GetSkillManager 获取技能管理器
func (gm *GenkitManager) GetSkillManager() *skills.SkillManager {
	return gm.skillManager
}

// GetMemoryManager 获取记忆管理器
func (gm *GenkitManager) GetMemoryManager() *memory.SearchEngine {
	return gm.memoryManager
}

// IsInitialized 检查是否已初始化
func (gm *GenkitManager) IsInitialized() bool {
	return gm.initialized
}

// Shutdown 关闭Genkit运行时
func (gm *GenkitManager) Shutdown() error {
	if !gm.initialized {
		return nil
	}

	utils.Info("关闭Genkit运行时...")

	// 关闭Agent管理器
	if gm.agentManager != nil {
		if err := gm.agentManager.Shutdown(); err != nil {
			utils.Warn("关闭Agent管理器失败: %v", err)
		}
	}

	// Flow引擎关闭
	if gm.flowEngine != nil {
		utils.Info("Flow引擎已关闭")
	}

	// 关闭工具注册表（如果有关闭方法）
	// 注意：toolRegistry目前没有Shutdown方法

	// 关闭工具发现服务
	if gm.toolDiscovery != nil {
		if err := gm.toolDiscovery.Stop(); err != nil {
			utils.Warn("关闭工具发现服务失败: %v", err)
		}
	}

	// 关闭MCP主机
	if gm.mcpHost != nil {
		if err := gm.mcpHost.Stop(); err != nil {
			utils.Warn("关闭MCP主机失败: %v", err)
		}
	}

	// 清理资源
	gm.modelAdapter = nil
	gm.toolRegistry = nil
	gm.flowEngine = nil
	gm.agentManager = nil
	gm.skillManager = nil
	gm.memoryManager = nil
	gm.mcpHost = nil
	gm.mcpConfig = nil
	gm.toolDiscovery = nil

	gm.initialized = false
	utils.Info("Genkit运行时已关闭")
	return nil
}

// Global Genkit manager instance
var globalGenkitManager *GenkitManager

// InitGlobalGenkit 初始化全局Genkit管理器
func InitGlobalGenkit(ctx context.Context) error {
	if globalGenkitManager == nil {
		globalGenkitManager = NewGenkitManager()
	}

	return globalGenkitManager.Initialize(ctx)
}

// GetGlobalGenkitManager 获取全局Genkit管理器
func GetGlobalGenkitManager() *GenkitManager {
	return globalGenkitManager
}

// GetGlobalGenkit 获取全局Genkit实例
func GetGlobalGenkit() interface{} {
	if globalGenkitManager == nil {
		return nil
	}
	return globalGenkitManager.GetGenkit()
}

// initSkillManager 初始化技能管理器
func (gm *GenkitManager) initSkillManager() error {
	utils.Info("初始化技能管理器...")

	// 获取技能目录路径
	skillsDir := "./data/skills"

	// 创建LLM客户端包装器
	var llmClient skills.LLMClient
	if gm.modelAdapter != nil {
		// 使用默认模型ID
		modelID := "gemini-1.5-flash" // 默认模型
		llmClient = skills.NewModelAdapterWrapper(gm.modelAdapter, modelID)
	} else {
		utils.Warn("模型适配器未初始化，使用模拟LLM客户端")
		// 使用模拟客户端
		llmClient = skills.NewLLMClientWrapper(nil)
	}

	// 创建技能管理器
	skillManager, err := skills.NewSkillManager(skillsDir, llmClient)
	if err != nil {
		return fmt.Errorf("failed to create skill manager: %w", err)
	}

	gm.skillManager = skillManager
	utils.Info("技能管理器初始化完成")
	return nil
}

// initMemoryManager 初始化记忆管理器
func (gm *GenkitManager) initMemoryManager() error {
	utils.Info("初始化记忆管理器...")

	// 获取记忆数据库路径
	dbPath := "./data/memory.db"

	// 创建记忆组件
	workingMemory := memory.NewWorkingMemoryManager(100, 24*time.Hour)
	episodicMemory, err := memory.NewEpisodicMemory(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create episodic memory: %w", err)
	}

	semanticMemory, err := memory.NewSemanticMemory(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create semantic memory: %w", err)
	}

	// 创建记忆检索引擎
	memoryManager := memory.NewSearchEngine(workingMemory, episodicMemory, semanticMemory)

	gm.memoryManager = memoryManager
	utils.Info("记忆管理器初始化完成")
	return nil
}

// initMCP 初始化MCP功能
func (gm *GenkitManager) initMCP() error {
	utils.Info("初始化MCP功能...")

	// 初始化MCP配置管理器
	mcpConfigPath := "./data/mcp_config.json"
	mcpConfigManager, err := mcpconfig.NewMCPConfigManager(mcpConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create MCP config manager: %w", err)
	}
	gm.mcpConfig = mcpConfigManager

	// 获取MCP配置
	mcpConfig := mcpConfigManager.GetConfig()

	// 创建MCP主机
	mcpHost := host.NewMCPHost(mcpConfig)
	gm.mcpHost = mcpHost

	// 启动MCP主机
	if mcpConfig.AutoStart {
		if err := mcpHost.Start(gm.ctx); err != nil {
			utils.Warn("启动MCP主机失败: %v", err)
		} else {
			utils.Info("MCP主机已启动")
		}
	}

	// 初始化工具发现服务
	if gm.toolRegistry != nil {
		toolDiscovery := discovery.NewToolDiscovery(mcpHost, gm.toolRegistry)
		gm.toolDiscovery = toolDiscovery

		// 设置全局工具发现实例
		discovery.SetGlobalToolDiscovery(toolDiscovery)

		// 启动工具发现服务
		if err := toolDiscovery.Start(gm.ctx); err != nil {
			utils.Warn("启动工具发现服务失败: %v", err)
		} else {
			utils.Info("工具发现服务已启动")
		}
	}

	utils.Info("MCP功能初始化完成")
	return nil
}

// GetMCPHost 获取MCP主机
func (gm *GenkitManager) GetMCPHost() *host.MCPHost {
	return gm.mcpHost
}

// GetMCPConfig 获取MCP配置管理器
func (gm *GenkitManager) GetMCPConfig() *mcpconfig.MCPConfigManager {
	return gm.mcpConfig
}

// GetToolDiscovery 获取工具发现服务
func (gm *GenkitManager) GetToolDiscovery() *discovery.ToolDiscovery {
	return gm.toolDiscovery
}
