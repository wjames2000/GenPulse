package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"GenPulse/internal/agents"
	"GenPulse/internal/genkit/models"
	"gopkg.in/yaml.v3"
)

// ConfigLoader 配置加载器
type ConfigLoader struct {
	configDir string
}

// NewConfigLoader 创建配置加载器
func NewConfigLoader(configDir string) *ConfigLoader {
	return &ConfigLoader{
		configDir: configDir,
	}
}

// LoadAgentConfig 加载Agent配置
func (cl *ConfigLoader) LoadAgentConfig(filePath string) (*agents.AgentConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 根据文件扩展名决定解析方式
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		return cl.parseJSONConfig(data)
	case ".yaml", ".yml":
		return cl.parseYAMLConfig(data)
	default:
		return nil, fmt.Errorf("不支持的配置文件格式: %s", ext)
	}
}

// LoadAgentConfigsFromDir 从目录加载所有Agent配置
func (cl *ConfigLoader) LoadAgentConfigsFromDir(dirPath string) ([]*agents.AgentConfig, error) {
	var configs []*agents.AgentConfig

	// 如果目录不存在，创建它
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return nil, fmt.Errorf("创建配置目录失败: %w", err)
		}
		return configs, nil
	}

	// 读取目录
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置目录失败: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		ext := strings.ToLower(filepath.Ext(filePath))
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			continue
		}

		config, err := cl.LoadAgentConfig(filePath)
		if err != nil {
			// 记录错误但继续加载其他文件
			fmt.Printf("警告: 加载配置文件 %s 失败: %v\n", filePath, err)
			continue
		}

		configs = append(configs, config)
	}

	return configs, nil
}

// parseJSONConfig 解析JSON配置
func (cl *ConfigLoader) parseJSONConfig(data []byte) (*agents.AgentConfig, error) {
	var rawConfig RawAgentConfig
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("解析JSON配置失败: %w", err)
	}

	return cl.convertRawConfig(&rawConfig)
}

// parseYAMLConfig 解析YAML配置
func (cl *ConfigLoader) parseYAMLConfig(data []byte) (*agents.AgentConfig, error) {
	var rawConfig RawAgentConfig
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("解析YAML配置失败: %w", err)
	}

	return cl.convertRawConfig(&rawConfig)
}

// convertRawConfig 转换原始配置为AgentConfig
func (cl *ConfigLoader) convertRawConfig(raw *RawAgentConfig) (*agents.AgentConfig, error) {
	// 验证必需字段
	if raw.ID == "" {
		return nil, fmt.Errorf("Agent配置缺少ID字段")
	}
	if raw.Name == "" {
		return nil, fmt.Errorf("Agent配置缺少Name字段")
	}
	if raw.Role == "" {
		return nil, fmt.Errorf("Agent配置缺少Role字段")
	}

	// 转换角色
	role, err := parseAgentRole(raw.Role)
	if err != nil {
		return nil, fmt.Errorf("无效的Agent角色: %w", err)
	}

	// 转换能力
	capabilities, err := parseCapabilities(raw.Capabilities)
	if err != nil {
		return nil, fmt.Errorf("无效的能力配置: %w", err)
	}

	// 转换模型配置
	modelConfig, err := parseModelConfig(raw.ModelConfig)
	if err != nil {
		return nil, fmt.Errorf("无效的模型配置: %w", err)
	}

	// 转换超时
	timeout, err := parseDuration(raw.Timeout)
	if err != nil {
		return nil, fmt.Errorf("无效的超时配置: %w", err)
	}

	// 设置默认值
	if raw.MaxRetries == 0 {
		raw.MaxRetries = 3
	}
	if timeout == 0 {
		timeout = 5 * time.Minute
	}
	if raw.PromptTemplates == nil {
		raw.PromptTemplates = make(map[string]string)
	}

	config := &agents.AgentConfig{
		ID:               raw.ID,
		Name:             raw.Name,
		Role:             role,
		Description:      raw.Description,
		ModelConfig:      modelConfig,
		Capabilities:     capabilities,
		Tools:            raw.Tools,
		PromptTemplates:  raw.PromptTemplates,
		MaxRetries:       raw.MaxRetries,
		Timeout:          timeout,
		Enabled:          raw.Enabled,
	}

	return config, nil
}

// SaveAgentConfig 保存Agent配置到文件
func (cl *ConfigLoader) SaveAgentConfig(config *agents.AgentConfig, filePath string) error {
	rawConfig := convertToRawConfig(config)

	ext := strings.ToLower(filepath.Ext(filePath))
	var data []byte
	var err error

	switch ext {
	case ".json":
		data, err = json.MarshalIndent(rawConfig, "", "  ")
	case ".yaml", ".yml":
		data, err = yaml.Marshal(rawConfig)
	default:
		return fmt.Errorf("不支持的配置文件格式: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// ExportAgentConfig 导出Agent配置为JSON或YAML
func (cl *ConfigLoader) ExportAgentConfig(config *agents.AgentConfig, format string, writer io.Writer) error {
	rawConfig := convertToRawConfig(config)

	var data []byte
	var err error

	switch strings.ToLower(format) {
	case "json":
		data, err = json.MarshalIndent(rawConfig, "", "  ")
	case "yaml", "yml":
		data, err = yaml.Marshal(rawConfig)
	default:
		return fmt.Errorf("不支持的导出格式: %s", format)
	}

	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	_, err = writer.Write(data)
	return err
}

// CreateDefaultConfigs 创建默认配置
func (cl *ConfigLoader) CreateDefaultConfigs() map[string]*agents.AgentConfig {
	defaultConfigs := make(map[string]*agents.AgentConfig)

	// 全栈开发Agent
	fullstackConfig := &agents.AgentConfig{
		ID:          "fullstack-developer",
		Name:        "全栈开发工程师",
		Role:        agents.RoleFullStackDev,
		Description: "全栈开发工程师，能够处理前后端开发任务",
		ModelConfig: models.ModelConfig{
			Name:     "gpt-4",
			Provider: "openai",
			Type:     "chat",
		},
		Capabilities: []agents.AgentCapability{
			agents.CapabilityCodeGeneration,
			agents.CapabilityFileOperation,
			agents.CapabilityGitOperation,
			agents.CapabilityShellExecution,
			agents.CapabilityProjectSetup,
		},
		Tools: []string{
			"fs_tool",
			"git_tool",
			"shell_tool",
			"project_tool",
		},
		PromptTemplates: map[string]string{
			"task_analysis": `你是一个全栈开发工程师。请分析以下任务需求：

任务：{{task}}

请按照以下步骤进行分析：
1. 理解需求：明确用户想要实现什么功能
2. 技术选型：根据需求选择合适的技术栈（前端、后端、数据库等）
3. 项目结构：设计合理的项目目录结构
4. 实现步骤：列出具体的实现步骤
5. 文件清单：需要创建哪些文件，每个文件的内容概要

请开始你的分析：`,
		},
		MaxRetries: 3,
		Timeout:    5 * time.Minute,
		Enabled:    true,
	}
	defaultConfigs[fullstackConfig.ID] = fullstackConfig

	// 前端开发Agent
	frontendConfig := &agents.AgentConfig{
		ID:          "frontend-developer",
		Name:        "前端开发工程师",
		Role:        agents.RoleFrontendDev,
		Description: "前端开发工程师，专注于React/Vue等前端技术",
		ModelConfig: models.ModelConfig{
			Name:     "gpt-4",
			Provider: "openai",
			Type:     "chat",
		},
		Capabilities: []agents.AgentCapability{
			agents.CapabilityCodeGeneration,
			agents.CapabilityFileOperation,
		},
		Tools: []string{
			"fs_tool",
		},
		PromptTemplates: map[string]string{
			"frontend_code": `你是一个前端开发工程师。请根据以下需求生成前端代码：

项目类型：{{project_type}}
技术栈：{{tech_stack}}
需求描述：{{requirement}}

请生成高质量的前端代码，包括：
1. 组件结构设计
2. 样式实现
3. 状态管理
4. 交互逻辑

请开始生成代码：`,
		},
		MaxRetries: 3,
		Timeout:    3 * time.Minute,
		Enabled:    true,
	}
	defaultConfigs[frontendConfig.ID] = frontendConfig

	// 后端开发Agent
	backendConfig := &agents.AgentConfig{
		ID:          "backend-developer",
		Name:        "后端开发工程师",
		Role:        agents.RoleBackendDev,
		Description: "后端开发工程师，专注于Go/Node.js等后端技术",
		ModelConfig: models.ModelConfig{
			Name:     "gpt-4",
			Provider: "openai",
			Type:     "chat",
		},
		Capabilities: []agents.AgentCapability{
			agents.CapabilityCodeGeneration,
			agents.CapabilityFileOperation,
		},
		Tools: []string{
			"fs_tool",
		},
		PromptTemplates: map[string]string{
			"backend_code": `你是一个后端开发工程师。请根据以下需求生成后端代码：

项目类型：{{project_type}}
技术栈：{{tech_stack}}
需求描述：{{requirement}}

请生成高质量的后端代码，包括：
1. API设计
2. 数据库模型
3. 业务逻辑
4. 错误处理

请开始生成代码：`,
		},
		MaxRetries: 3,
		Timeout:    3 * time.Minute,
		Enabled:    true,
	}
	defaultConfigs[backendConfig.ID] = backendConfig

	return defaultConfigs
}

// SaveDefaultConfigs 保存默认配置到文件
func (cl *ConfigLoader) SaveDefaultConfigs() error {
	defaultConfigs := cl.CreateDefaultConfigs()

	// 确保配置目录存在
	if err := os.MkdirAll(cl.configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 保存每个默认配置
	for _, config := range defaultConfigs {
		filePath := filepath.Join(cl.configDir, fmt.Sprintf("%s.yaml", config.ID))
		if err := cl.SaveAgentConfig(config, filePath); err != nil {
			return fmt.Errorf("保存默认配置 %s 失败: %w", config.ID, err)
		}
	}

	return nil
}

// GetConfigDir 获取配置目录
func (cl *ConfigLoader) GetConfigDir() string {
	return cl.configDir
}

// SetConfigDir 设置配置目录
func (cl *ConfigLoader) SetConfigDir(configDir string) {
	cl.configDir = configDir
}