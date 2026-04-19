package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"GenPulse/internal/utils"
)

// AppConfig 应用配置
type AppConfig struct {
	// 应用设置
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	LogLevel   string `json:"log_level"`
	LogToFile  bool   `json:"log_to_file"`

	// 工作区设置
	WorkspacePath string `json:"workspace_path"`
	MaxFileSize   int64  `json:"max_file_size"` // 字节

	// 模型设置
	DefaultModel  string  `json:"default_model"`
	FallbackModel string  `json:"fallback_model"`
	MaxTokens     int     `json:"max_tokens"`
	Temperature   float64 `json:"temperature"`

	// API密钥
	APIKeys map[string]string `json:"api_keys"`

	// 代理设置
	HTTPProxy  string `json:"http_proxy"`
	HTTPSProxy string `json:"https_proxy"`

	// 高级设置
	EnableAutoSave   bool `json:"enable_auto_save"`
	AutoSaveInterval int  `json:"auto_save_interval"` // 秒
	MaxHistoryItems  int  `json:"max_history_items"`
}

// DefaultConfig 默认配置
func DefaultConfig() *AppConfig {
	return &AppConfig{
		AppName:          "GenPulse",
		AppVersion:       "1.0.0",
		LogLevel:         "info",
		LogToFile:        true,
		WorkspacePath:    "workspace",
		MaxFileSize:      10 * 1024 * 1024, // 10MB
		DefaultModel:     "gemini-1.5-pro",
		FallbackModel:    "gpt-4-turbo",
		MaxTokens:        4096,
		Temperature:      0.7,
		APIKeys:          make(map[string]string),
		EnableAutoSave:   true,
		AutoSaveInterval: 300, // 5分钟
		MaxHistoryItems:  100,
	}
}

// ConfigManager 配置管理器
type ConfigManager struct {
	config     *AppConfig
	configPath string
}

// NewConfigManager 创建新的配置管理器
func NewConfigManager(configPath string) *ConfigManager {
	return &ConfigManager{
		config:     DefaultConfig(),
		configPath: configPath,
	}
}

// Load 加载配置
func (cm *ConfigManager) Load() error {
	// 确保配置目录存在
	configDir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// 配置文件不存在，使用默认配置并保存
		utils.Info("Config file not found, creating with default values: %s", cm.configPath)
		return cm.Save()
	}

	// 读取配置文件
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析JSON
	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	cm.config = &config
	utils.Info("Config loaded from: %s", cm.configPath)
	return nil
}

// Save 保存配置
func (cm *ConfigManager) Save() error {
	// 序列化为JSON
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	utils.Info("Config saved to: %s", cm.configPath)
	return nil
}

// GetConfig 获取配置
func (cm *ConfigManager) GetConfig() *AppConfig {
	return cm.config
}

// UpdateConfig 更新配置
func (cm *ConfigManager) UpdateConfig(updates map[string]interface{}) error {
	// 这里可以添加配置验证逻辑
	// 暂时简单更新

	// 保存更新后的配置
	return cm.Save()
}

// GetAPIKey 获取API密钥
func (cm *ConfigManager) GetAPIKey(provider string) string {
	if key, exists := cm.config.APIKeys[provider]; exists {
		return key
	}
	return ""
}

// SetAPIKey 设置API密钥
func (cm *ConfigManager) SetAPIKey(provider, key string) error {
	if cm.config.APIKeys == nil {
		cm.config.APIKeys = make(map[string]string)
	}
	cm.config.APIKeys[provider] = key
	return cm.Save()
}

// GetWorkspacePath 获取工作区路径
func (cm *ConfigManager) GetWorkspacePath() string {
	return cm.config.WorkspacePath
}

// SetWorkspacePath 设置工作区路径
func (cm *ConfigManager) SetWorkspacePath(path string) error {
	cm.config.WorkspacePath = path
	return cm.Save()
}

// Global config instance
var globalConfigManager *ConfigManager

// InitGlobalConfig 初始化全局配置
func InitGlobalConfig() error {
	// 确定配置文件路径
	configDir := "config"
	configPath := filepath.Join(configDir, "app_config.json")

	// 创建配置管理器
	cm := NewConfigManager(configPath)

	// 加载配置
	if err := cm.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	globalConfigManager = cm
	return nil
}

// GetGlobalConfig 获取全局配置管理器
func GetGlobalConfig() *ConfigManager {
	if globalConfigManager == nil {
		// 尝试初始化
		if err := InitGlobalConfig(); err != nil {
			utils.Error("Failed to initialize config: %v", err)
			// 返回一个临时配置管理器
			globalConfigManager = NewConfigManager("config/app_config.json")
		}
	}
	return globalConfigManager
}

// GetConfig 获取全局配置
func GetConfig() *AppConfig {
	return GetGlobalConfig().GetConfig()
}
