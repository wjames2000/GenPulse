package main

import (
	"context"
	"fmt"

	"GenPulse/internal/genkit"
	"GenPulse/internal/services"
)

// App struct
type App struct {
	ctx           context.Context
	baseService   *services.BaseService
	genkitManager *genkit.GenkitManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		baseService:   services.NewBaseService(),
		genkitManager: genkit.NewGenkitManager(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 启动基础服务
	a.baseService.Startup(ctx)

	// 启动Genkit运行时
	if err := a.genkitManager.Initialize(ctx); err != nil {
		fmt.Printf("警告: Genkit运行时初始化失败: %v\n", err)
		// 不阻止应用启动，记录错误继续
		a.baseService.LogMessage("error", fmt.Sprintf("Genkit初始化失败: %v", err))
	} else {
		a.baseService.LogMessage("info", "Genkit运行时初始化完成")
	}
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
