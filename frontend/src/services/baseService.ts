// Wails运行时导入
import { EventsOn, EventsOff } from '../../wailsjs/runtime';

// 基础服务接口
export interface AppInfo {
  name: string;
  version: string;
  status: string;
}

export interface HealthStatus {
  status: string;
  service: string;
  timestamp: number;
}

export interface LogMessage {
  level: string;
  message: string;
  time: string;
}

// 基础服务类
class BaseService {
  private eventListeners: Map<string, Function[]> = new Map();

  /**
   * 获取应用信息
   */
  async getAppInfo(): Promise<AppInfo> {
    try {
      // 注意：这里使用Wails的invoke方式调用Go后端
      // 实际调用方式需要根据Wails的绑定机制调整
      const result = await (window as any).go.main.App.GetAppInfo();
      return result as AppInfo;
    } catch (error) {
      console.error('Failed to get app info:', error);
      return {
        name: 'GenPulse',
        version: '1.0.0',
        status: 'error'
      };
    }
  }

  /**
   * 健康检查
   */
  async healthCheck(): Promise<HealthStatus> {
    try {
      const result = await (window as any).go.main.App.HealthCheck();
      return result as HealthStatus;
    } catch (error) {
      console.error('Health check failed:', error);
      return {
        status: 'unhealthy',
        service: 'base_service',
        timestamp: Date.now()
      };
    }
  }

  /**
   * 发送日志消息
   */
  async logMessage(level: string, message: string): Promise<void> {
    try {
      await (window as any).go.main.App.LogMessage(level, message);
    } catch (error) {
      console.error('Failed to log message:', error);
    }
  }

  /**
   * 监听事件
   */
  async listenToEvent(eventName: string, callback: (data: any) => void): Promise<void> {
    try {
      // Wails的事件监听
      const cleanup = EventsOn(eventName, (data: any) => {
        try {
          const parsedData = typeof data === 'string' 
            ? JSON.parse(data)
            : data;
          callback(parsedData);
        } catch (error) {
          console.error(`Failed to parse event data for ${eventName}:`, error);
        }
      });

      // 保存监听器以便后续清理
      if (!this.eventListeners.has(eventName)) {
        this.eventListeners.set(eventName, []);
      }
      this.eventListeners.get(eventName)!.push(() => {
        EventsOff(eventName);
        if (cleanup && typeof cleanup === 'function') {
          cleanup();
        }
      });
    } catch (error) {
      console.error(`Failed to listen to event ${eventName}:`, error);
    }
  }

  /**
   * 监听日志事件
   */
  async listenToLogs(callback: (log: LogMessage) => void): Promise<void> {
    await this.listenToEvent('log', callback);
  }

  /**
   * 清理所有事件监听器
   */
  cleanupEventListeners(): void {
    for (const [eventName, listeners] of this.eventListeners.entries()) {
      listeners.forEach(cleanup => {
        if (typeof cleanup === 'function') {
          cleanup();
        }
      });
      console.log(`Cleaned up ${listeners.length} listeners for event: ${eventName}`);
    }
    this.eventListeners.clear();
  }

  /**
   * 初始化服务
   */
  async initialize(): Promise<boolean> {
    try {
      // 执行健康检查
      const health = await this.healthCheck();
      if (health.status !== 'healthy') {
        console.warn('Service health check failed:', health);
        return false;
      }

      // 获取应用信息
      const appInfo = await this.getAppInfo();
      console.log('BaseService initialized:', appInfo);

      // 监听日志
      await this.listenToLogs((log) => {
        console.log(`[${log.level}] ${log.time}: ${log.message}`);
      });

      return true;
    } catch (error) {
      console.error('Failed to initialize BaseService:', error);
      return false;
    }
  }
}

// 导出单例实例
export const baseService = new BaseService();

// 导出类型和实例
export default baseService;