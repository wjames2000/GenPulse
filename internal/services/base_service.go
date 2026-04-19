package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// BaseService 提供前端调用的基础服务封装
type BaseService struct {
	ctx context.Context
}

// NewBaseService 创建新的基础服务实例
func NewBaseService() *BaseService {
	return &BaseService{}
}

// Startup 在应用启动时调用，保存上下文
func (b *BaseService) Startup(ctx context.Context) {
	b.ctx = ctx
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
	if b.ctx == nil {
		return fmt.Errorf("context not initialized")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	runtime.EventsEmit(b.ctx, eventName, string(jsonData))
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
	log.Printf("[%s] %s", level, message)

	// 同时发送日志事件到前端
	b.SendEventToFrontend("log", map[string]interface{}{
		"level":   level,
		"message": message,
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	})
}
