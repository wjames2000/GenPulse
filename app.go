package main

import (
	"context"
	"fmt"

	"GenPulse/internal/services"
)

// App struct
type App struct {
	ctx         context.Context
	baseService *services.BaseService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		baseService: services.NewBaseService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.baseService.Startup(ctx)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetAppInfo 获取应用信息
func (a *App) GetAppInfo() map[string]interface{} {
	return a.baseService.GetAppInfo()
}

// HealthCheck 健康检查
func (a *App) HealthCheck() map[string]interface{} {
	return a.baseService.HealthCheck()
}

// LogMessage 记录日志消息
func (a *App) LogMessage(level string, message string) {
	a.baseService.LogMessage(level, message)
}
