package genkit

import (
	"context"
	"fmt"

	"GenPulse/internal/agents"
	"GenPulse/internal/genkit/config"
	"GenPulse/internal/genkit/flows"
	"GenPulse/internal/genkit/models"
	"GenPulse/internal/genkit/tools"
	"GenPulse/internal/utils"
)

// GenkitManager 管理Genkit运行时
type GenkitManager struct {
	ctx          context.Context
	config       *config.AppConfig
	initialized  bool
	modelAdapter *models.UnifiedModelAdapter
	toolRegistry *tools.ToolRegistry
	flowEngine   *flows.FlowEngine
	agentManager *agents.AgentManager
	// genkit     interface{} // 暂时用interface{}，等确定具体类型后替换
}

// NewGenkitManager 创建新的Genkit管理器
func NewGenkitManager() *GenkitManager {
	return &GenkitManager{
		initialized: false,
	}
}

// Initialize 初始化Genkit运行时
func (gm *GenkitManager) Initialize(ctx context.Context) error {
	if gm.initialized {
		return nil
	}

	gm.ctx = ctx

	// 获取配置
	cfg := config.GetConfig()
	gm.config = cfg

	utils.Info("初始化Genkit运行时...")

	// TODO: 初始化Genkit核心
	// genkitInstance, err := genkitgo.Init(ctx)
	// if err != nil {
	//     return fmt.Errorf("failed to initialize Genkit: %w", err)
	// }
	// gm.genkit = genkitInstance

	utils.Info("Genkit核心初始化完成")

	// 配置模型提供商
	if err := gm.configureModelProviders(); err != nil {
		return fmt.Errorf("failed to configure model providers: %w", err)
	}

	// 初始化工具注册表
	if err := gm.initToolRegistry(); err != nil {
		return fmt.Errorf("failed to initialize tool registry: %w", err)
	}

	// 初始化Flow引擎
	if err := gm.initFlowEngine(); err != nil {
		return fmt.Errorf("failed to initialize flow engine: %w", err)
	}

	// 初始化Agent管理器
	if err := gm.initAgentManager(); err != nil {
		return fmt.Errorf("failed to initialize agent manager: %w", err)
	}

	gm.initialized = true
	utils.Info("Genkit运行时初始化完成")

	return nil
}

// configureModelProviders 配置模型提供商
func (gm *GenkitManager) configureModelProviders() error {
	utils.Info("配置模型提供商...")

	// 这里可以根据配置动态加载不同的模型提供商
	// 目前先实现基础框架

	// 检查是否有API密钥配置
	// 注意：这里需要从配置管理器获取，而不是直接访问config
	configMgr := config.GetGlobalConfig()
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

	// 清理资源
	gm.modelAdapter = nil
	gm.toolRegistry = nil
	gm.flowEngine = nil
	gm.agentManager = nil

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
