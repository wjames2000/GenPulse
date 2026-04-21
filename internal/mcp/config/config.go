package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"GenPulse/internal/mcp/client"
	"GenPulse/internal/mcp/host"
	"GenPulse/internal/utils"
)

// MCPConfigManager MCP配置管理器
type MCPConfigManager struct {
	configPath string
	config     host.MCPHostConfig
	mu         sync.RWMutex
}

// NewMCPConfigManager 创建新的MCP配置管理器
func NewMCPConfigManager(configPath string) (*MCPConfigManager, error) {
	manager := &MCPConfigManager{
		configPath: configPath,
		config: host.MCPHostConfig{
			AutoStart:             true,
			ToolDiscoveryInterval: 60,
			MaxConcurrentCalls:    10,
			Servers:               []host.MCPHostServerConfig{},
		},
	}

	// 确保配置目录存在
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// 加载配置
	if err := manager.loadConfig(); err != nil {
		utils.Warn("加载MCP配置失败: %v，使用默认配置", err)
		// 保存默认配置
		if err := manager.saveConfig(); err != nil {
			utils.Warn("保存默认MCP配置失败: %v", err)
		}
	}

	return manager, nil
}

// loadConfig 加载配置
func (m *MCPConfigManager) loadConfig() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查配置文件是否存在
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", m.configPath)
	}

	// 读取配置文件
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析JSON
	var config host.MCPHostConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config JSON: %w", err)
	}

	m.config = config
	utils.Info("MCP配置已加载: %s", m.configPath)
	return nil
}

// saveConfig 保存配置
func (m *MCPConfigManager) saveConfig() error {
	m.mu.RLock()
	configToSave := m.config
	m.mu.RUnlock()

	return m.saveConfigWithConfig(configToSave)
}

// saveConfigWithConfig 使用指定的配置保存
func (m *MCPConfigManager) saveConfigWithConfig(config host.MCPHostConfig) error {
	// 序列化为JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	utils.Debug("MCP配置已保存: %s", m.configPath)
	return nil
}

// GetConfig 获取配置
func (m *MCPConfigManager) GetConfig() host.MCPHostConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// UpdateConfig 更新配置
func (m *MCPConfigManager) UpdateConfig(config host.MCPHostConfig) error {
	m.mu.Lock()
	m.config = config
	configToSave := m.config
	m.mu.Unlock()

	// 保存到文件
	if err := m.saveConfigWithConfig(configToSave); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	utils.Info("MCP配置已更新")
	return nil
}

// AddServer 添加服务器配置
func (m *MCPConfigManager) AddServer(serverConfig host.MCPHostServerConfig) error {
	m.mu.Lock()

	// 检查ID是否已存在
	for _, server := range m.config.Servers {
		if server.ID == serverConfig.ID {
			m.mu.Unlock()
			return fmt.Errorf("server with ID %s already exists", serverConfig.ID)
		}
	}

	// 设置默认值
	if serverConfig.Priority == 0 {
		serverConfig.Priority = 50
	}

	if !serverConfig.Enabled {
		serverConfig.Enabled = true
	}

	// 添加服务器
	m.config.Servers = append(m.config.Servers, serverConfig)

	// 保存配置前释放锁
	configToSave := m.config
	m.mu.Unlock()

	// 保存配置
	if err := m.saveConfigWithConfig(configToSave); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	utils.Info("已添加MCP服务器配置: %s (%s)", serverConfig.Name, serverConfig.ID)
	return nil
}

// RemoveServer 移除服务器配置
func (m *MCPConfigManager) RemoveServer(serverID string) error {
	m.mu.Lock()

	// 查找服务器
	found := false
	newServers := make([]host.MCPHostServerConfig, 0, len(m.config.Servers))

	for _, server := range m.config.Servers {
		if server.ID == serverID {
			found = true
			continue
		}
		newServers = append(newServers, server)
	}

	if !found {
		m.mu.Unlock()
		return fmt.Errorf("server not found: %s", serverID)
	}

	m.config.Servers = newServers

	// 保存配置前释放锁
	configToSave := m.config
	m.mu.Unlock()

	// 保存配置
	if err := m.saveConfigWithConfig(configToSave); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	utils.Info("已移除MCP服务器配置: %s", serverID)
	return nil
}

// UpdateServer 更新服务器配置
func (m *MCPConfigManager) UpdateServer(serverID string, serverConfig host.MCPHostServerConfig) error {
	m.mu.Lock()

	// 查找服务器
	found := false
	for i, server := range m.config.Servers {
		if server.ID == serverID {
			// 确保ID不变
			serverConfig.ID = serverID
			m.config.Servers[i] = serverConfig
			found = true
			break
		}
	}

	if !found {
		m.mu.Unlock()
		return fmt.Errorf("server not found: %s", serverID)
	}

	// 保存配置前释放锁
	configToSave := m.config
	m.mu.Unlock()

	// 保存配置
	if err := m.saveConfigWithConfig(configToSave); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	utils.Info("已更新MCP服务器配置: %s", serverID)
	return nil
}

// GetServer 获取服务器配置
func (m *MCPConfigManager) GetServer(serverID string) (*host.MCPHostServerConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, server := range m.config.Servers {
		if server.ID == serverID {
			return &server, nil
		}
	}

	return nil, fmt.Errorf("server not found: %s", serverID)
}

// ListServers 列出所有服务器配置
func (m *MCPConfigManager) ListServers() []host.MCPHostServerConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 返回副本
	servers := make([]host.MCPHostServerConfig, len(m.config.Servers))
	copy(servers, m.config.Servers)
	return servers
}

// GetEnabledServers 获取启用的服务器配置
func (m *MCPConfigManager) GetEnabledServers() []host.MCPHostServerConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	enabledServers := make([]host.MCPHostServerConfig, 0)
	for _, server := range m.config.Servers {
		if server.Enabled {
			enabledServers = append(enabledServers, server)
		}
	}

	return enabledServers
}

// GetConfigPath 获取配置文件路径
func (m *MCPConfigManager) GetConfigPath() string {
	return m.configPath
}

// GetDefaultConfig 获取默认配置示例
func GetDefaultConfig() host.MCPHostConfig {
	return host.MCPHostConfig{
		AutoStart:             true,
		ToolDiscoveryInterval: 60,
		MaxConcurrentCalls:    10,
		Servers: []host.MCPHostServerConfig{
			{
				ID:   "local-fs-tools",
				Name: "本地文件系统工具",
				Type: "server",
				ServerConfig: host.MCPServerConfig{
					Type: "stdio",
				},
				Enabled:  true,
				Priority: 100,
			},
			{
				ID:   "local-git-tools",
				Name: "本地Git工具",
				Type: "server",
				ServerConfig: host.MCPServerConfig{
					Type: "stdio",
				},
				Enabled:  true,
				Priority: 90,
			},
			{
				ID:   "weather-api",
				Name: "天气API",
				Type: "client",
				ClientConfig: client.MCPClientConfig{
					ServerType: "stdio",
					Command:    "npx",
					Args:       []string{"@modelcontextprotocol/server-weather"},
					Namespace:  "weather",
					Timeout:    30,
				},
				Enabled:  true,
				Priority: 80,
			},
			{
				ID:   "filesystem-browser",
				Name: "文件系统浏览器",
				Type: "client",
				ClientConfig: client.MCPClientConfig{
					ServerType: "stdio",
					Command:    "npx",
					Args:       []string{"@modelcontextprotocol/server-filesystem"},
					Namespace:  "fs",
					Timeout:    30,
				},
				Enabled:  true,
				Priority: 70,
			},
		},
	}
}

// CreateDefaultConfig 创建默认配置文件
func CreateDefaultConfig(configPath string) error {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// 获取默认配置
	defaultConfig := GetDefaultConfig()

	// 序列化为JSON
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	utils.Info("已创建默认MCP配置文件: %s", configPath)
	return nil
}
