# 系统设置页面文档

## 概述
系统设置页面为GenPulse应用程序提供了一个完整的配置界面，允许用户自定义应用程序的各种设置。

## 功能特性

### 1. 通用设置 (General)
- **应用名称**: 自定义应用程序显示名称
- **语言**: 支持多语言界面 (英语、中文、日语、韩语)
- **主题**: 支持亮色、暗色和自动(跟随系统)主题
- **开机自启**: 设置应用程序是否随系统启动
- **通知**: 启用或禁用系统通知

### 2. API设置
- **OpenAI API密钥**: 配置OpenAI API访问密钥
- **Anthropic API密钥**: 配置Claude API访问密钥
- **Google Gemini API密钥**: 配置Google Gemini API访问密钥
- **OpenRouter API密钥**: 配置OpenRouter API访问密钥
- **API基础URL**: 自定义API端点地址
- **超时设置**: 设置API请求超时时间(秒)

### 3. 系统设置
- **最大并发代理数**: 控制同时运行的代理数量(1-10)
- **内存限制**: 设置应用程序内存使用限制(MB)
- **日志级别**: 设置日志详细程度(调试、信息、警告、错误)
- **自动保存**: 启用或禁用项目自动保存
- **保存间隔**: 设置自动保存间隔时间(秒)

### 4. 路径设置
- **项目目录**: 设置项目文件存储路径
- **日志目录**: 设置日志文件存储路径
- **缓存目录**: 设置缓存文件存储路径
- **模型目录**: 设置AI模型文件存储路径
- **浏览按钮**: 提供图形化目录选择器(待实现)

### 5. 高级设置
- **实验性功能**: 启用开发中的新功能
- **遥测数据**: 启用匿名使用数据收集
- **开发者模式**: 显示开发者工具和高级选项
- **自定义CSS**: 提供界面自定义样式

## 技术实现

### 状态管理
- 使用Zustand进行状态管理
- 设置数据持久化到localStorage
- 实时检测设置变化

### 组件结构
```
SettingsView
├── 标签导航 (Tab Navigation)
├── 设置卡片 (Settings Card)
│   ├── 通用设置表单
│   ├── API设置表单
│   ├── 系统设置表单
│   ├── 路径设置表单
│   └── 高级设置表单
├── 操作按钮 (Save/Reset)
└── 设置信息面板
```

### 数据模型
```typescript
interface SettingsData {
  general: {
    appName: string;
    language: string;
    theme: 'light' | 'dark' | 'auto';
    autoStart: boolean;
    notifications: boolean;
  };
  api: {
    openaiApiKey: string;
    anthropicApiKey: string;
    geminiApiKey: string;
    openrouterApiKey: string;
    baseUrl: string;
    timeout: number;
  };
  system: {
    maxConcurrentAgents: number;
    memoryLimit: number;
    logLevel: 'debug' | 'info' | 'warn' | 'error';
    autoSave: boolean;
    saveInterval: number;
  };
  paths: {
    projectsPath: string;
    logsPath: string;
    cachePath: string;
    modelsPath: string;
  };
  advanced: {
    enableExperimental: boolean;
    enableTelemetry: boolean;
    developerMode: boolean;
    customCss: string;
  };
}
```

## 使用说明

### 访问设置页面
1. 启动GenPulse应用程序
2. 点击左侧边栏的"Settings"图标(⚙️)
3. 系统设置页面将显示在主内容区域

### 修改设置
1. 点击左侧标签切换到相应设置类别
2. 修改表单中的设置值
3. 点击"Save Changes"按钮保存设置
4. 如需恢复默认设置，点击"Reset to Default"按钮

### 设置生效
- 大多数设置立即生效
- 部分设置(如主题)可能需要刷新页面
- 系统相关设置可能需要重启应用程序

## 开发说明

### 添加新设置项
1. 在`SettingsData`接口中添加新的设置字段
2. 在`SettingsView`组件的初始状态中添加默认值
3. 在相应的设置类别表单中添加输入控件
4. 更新`handleSettingChange`函数处理新字段

### 样式定制
- 使用Tailwind CSS类名进行样式控制
- 支持暗色/亮色主题
- 可通过高级设置中的自定义CSS进一步定制

### 后端集成
当前版本使用localStorage进行设置存储。如需后端集成：
1. 创建设置API端点
2. 修改`loadSettings`和`saveSettings`函数
3. 添加错误处理和加载状态

## 注意事项
- API密钥等敏感信息仅存储在本地
- 设置更改前会提示用户确认
- 提供设置重置功能以防配置错误
- 支持设置导入/导出功能(待实现)