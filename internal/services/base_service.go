package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// LogEntry 日志条目
type LogEntry struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	AgentID   string                 `json:"agent_id,omitempty"`
	TaskID    string                 `json:"task_id,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Duration  int64                  `json:"duration,omitempty"` // 毫秒
	Tags      []string               `json:"tags,omitempty"`
}

// BaseService 提供前端调用的基础服务封装
type BaseService struct {
	Ctx        context.Context
	logs       []LogEntry
	logsMutex  sync.RWMutex
	maxLogSize int
}

// NewBaseService 创建新的基础服务实例
func NewBaseService() *BaseService {
	return &BaseService{
		logs:       make([]LogEntry, 0),
		maxLogSize: 1000, // 最多保存1000条日志
	}
}

// Startup 在应用启动时调用，保存上下文
func (b *BaseService) Startup(ctx context.Context) {
	b.Ctx = ctx
	log.Println("BaseService started")
}

// CallMethod 通用方法调用封装
func (b *BaseService) CallMethod(method string, params interface{}) (interface{}, error) {
	log.Printf("Calling method: %s with params: %v", method, params)

	// 这里可以添加通用的前置处理逻辑
	// 比如日志记录、权限验证、参数验证等

	return nil, fmt.Errorf("method %s not implemented", method)
}

// SendEventToFrontend 发送事件到前端
func (b *BaseService) SendEventToFrontend(eventName string, data interface{}) error {
	if b.Ctx == nil {
		return fmt.Errorf("context not initialized")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	runtime.EventsEmit(b.Ctx, eventName, string(jsonData))
	log.Printf("Event sent to frontend: %s", eventName)
	return nil
}

// GetAppInfo 获取应用信息
func (b *BaseService) GetAppInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":    "GenPulse",
		"version": "1.0.0",
		"status":  "running",
	}
}

// HealthCheck 健康检查
func (b *BaseService) HealthCheck() map[string]interface{} {
	return map[string]interface{}{
		"status":    "healthy",
		"service":   "base_service",
		"timestamp": time.Now().Unix(),
	}
}

// LogMessage 记录日志消息
func (b *BaseService) LogMessage(level string, message string) {
	b.LogMessageWithDetails(level, message, nil, "", "", nil, 0)
}

// LogMessageWithDetails 记录带详细信息的日志消息
func (b *BaseService) LogMessageWithDetails(
	level string,
	message string,
	details map[string]interface{},
	agentID string,
	taskID string,
	tags []string,
	duration int64,
) {
	// 创建日志条目
	logEntry := LogEntry{
		ID:        fmt.Sprintf("log-%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		AgentID:   agentID,
		TaskID:    taskID,
		Details:   details,
		Duration:  duration,
		Tags:      tags,
	}

	// 保存到内存
	b.logsMutex.Lock()
	b.logs = append(b.logs, logEntry)

	// 限制日志数量
	if len(b.logs) > b.maxLogSize {
		b.logs = b.logs[len(b.logs)-b.maxLogSize:]
	}
	b.logsMutex.Unlock()

	// 输出到控制台
	log.Printf("[%s] %s", level, message)

	// 发送日志事件到前端（仅在Wails上下文中）
	if b.Ctx != nil {
		_ = b.SendEventToFrontend("log", map[string]interface{}{
			"id":        logEntry.ID,
			"timestamp": logEntry.Timestamp.Format(time.RFC3339),
			"level":     level,
			"message":   message,
			"agent_id":  agentID,
			"task_id":   taskID,
			"details":   details,
			"duration":  duration,
			"tags":      tags,
		})
	}
}

// GetLogs 获取日志
func (b *BaseService) GetLogs() []map[string]interface{} {
	b.logsMutex.RLock()
	defer b.logsMutex.RUnlock()

	logs := make([]map[string]interface{}, len(b.logs))
	for i, entry := range b.logs {
		logs[i] = map[string]interface{}{
			"id":        entry.ID,
			"timestamp": entry.Timestamp.Format(time.RFC3339),
			"level":     entry.Level,
			"message":   entry.Message,
			"agent_id":  entry.AgentID,
			"task_id":   entry.TaskID,
			"details":   entry.Details,
			"duration":  entry.Duration,
			"tags":      entry.Tags,
		}
	}

	return logs
}

// GetLogsByLevel 按级别获取日志
func (b *BaseService) GetLogsByLevel(level string) []map[string]interface{} {
	b.logsMutex.RLock()
	defer b.logsMutex.RUnlock()

	var filtered []map[string]interface{}
	for _, entry := range b.logs {
		if entry.Level == level {
			filtered = append(filtered, map[string]interface{}{
				"id":        entry.ID,
				"timestamp": entry.Timestamp.Format(time.RFC3339),
				"level":     entry.Level,
				"message":   entry.Message,
				"agent_id":  entry.AgentID,
				"task_id":   entry.TaskID,
				"details":   entry.Details,
				"duration":  entry.Duration,
				"tags":      entry.Tags,
			})
		}
	}

	return filtered
}

// GetLogsByAgent 按Agent获取日志
func (b *BaseService) GetLogsByAgent(agentID string) []map[string]interface{} {
	b.logsMutex.RLock()
	defer b.logsMutex.RUnlock()

	var filtered []map[string]interface{}
	for _, entry := range b.logs {
		if entry.AgentID == agentID {
			filtered = append(filtered, map[string]interface{}{
				"id":        entry.ID,
				"timestamp": entry.Timestamp.Format(time.RFC3339),
				"level":     entry.Level,
				"message":   entry.Message,
				"agent_id":  entry.AgentID,
				"task_id":   entry.TaskID,
				"details":   entry.Details,
				"duration":  entry.Duration,
				"tags":      entry.Tags,
			})
		}
	}

	return filtered
}

// GetLogsByTimeRange 按时间范围获取日志
func (b *BaseService) GetLogsByTimeRange(startTime, endTime time.Time) []map[string]interface{} {
	b.logsMutex.RLock()
	defer b.logsMutex.RUnlock()

	var filtered []map[string]interface{}
	for _, entry := range b.logs {
		if !entry.Timestamp.Before(startTime) && !entry.Timestamp.After(endTime) {
			filtered = append(filtered, map[string]interface{}{
				"id":        entry.ID,
				"timestamp": entry.Timestamp.Format(time.RFC3339),
				"level":     entry.Level,
				"message":   entry.Message,
				"agent_id":  entry.AgentID,
				"task_id":   entry.TaskID,
				"details":   entry.Details,
				"duration":  entry.Duration,
				"tags":      entry.Tags,
			})
		}
	}

	return filtered
}

// ClearLogs 清空日志
func (b *BaseService) ClearLogs() {
	b.logsMutex.Lock()
	b.logs = make([]LogEntry, 0)
	b.logsMutex.Unlock()
}

// GetLogStatistics 获取日志统计信息
func (b *BaseService) GetLogStatistics() map[string]interface{} {
	b.logsMutex.RLock()
	defer b.logsMutex.RUnlock()

	stats := map[string]interface{}{
		"total_logs": len(b.logs),
		"levels":     map[string]int{},
		"agents":     map[string]int{},
		"time_range": map[string]interface{}{},
	}

	if len(b.logs) > 0 {
		// 统计各级别日志数量
		for _, entry := range b.logs {
			levelStats := stats["levels"].(map[string]int)
			levelStats[entry.Level] = levelStats[entry.Level] + 1

			if entry.AgentID != "" {
				agentStats := stats["agents"].(map[string]int)
				agentStats[entry.AgentID] = agentStats[entry.AgentID] + 1
			}
		}

		// 时间范围
		timeRange := stats["time_range"].(map[string]interface{})
		timeRange["first"] = b.logs[0].Timestamp.Format(time.RFC3339)
		timeRange["last"] = b.logs[len(b.logs)-1].Timestamp.Format(time.RFC3339)
		timeRange["duration"] = b.logs[len(b.logs)-1].Timestamp.Sub(b.logs[0].Timestamp).String()
	}

	return stats
}
