# GenPulse 编译脚本

本项目提供了多种编译和构建工具，支持跨平台开发和部署。

## 可用工具

### 1. Shell 脚本 (Linux/macOS)
**文件**: `scripts/build.sh`

**功能**:
- 完整的构建流程管理
- 跨平台构建支持
- 环境检查
- 依赖管理
- 测试和代码检查

**使用方法**:
```bash
# 给予执行权限
chmod +x scripts/build.sh

# 检查环境
./scripts/build.sh check

# 启动开发模式
./scripts/build.sh dev

# 构建当前平台
./scripts/build.sh build

# 构建 macOS 版本
./scripts/build.sh build darwin universal

# 构建所有平台
./scripts/build.sh build-all

# 执行完整构建流程
./scripts/build.sh all

# 显示帮助
./scripts/build.sh help
```

### 2. Makefile (所有平台)
**文件**: `Makefile`

**功能**:
- 跨平台兼容
- 简单的命令接口
- 集成到开发工作流

**使用方法**:
```bash
# 检查环境
make check

# 清理构建文件
make clean

# 安装依赖
make deps

# 启动开发模式
make dev

# 构建当前平台
make build

# 构建 macOS 版本
make build-darwin

# 构建 Windows 版本
make build-windows

# 构建 Linux 版本
make build-linux

# 构建所有平台
make build-all

# 运行测试
make test

# 运行代码检查
make lint

# 完整构建流程
make all

# 显示帮助
make help
```

### 3. Windows 批处理脚本
**文件**: `scripts/build.bat`

**功能**:
- Windows 专用
- 基本的构建功能
- 环境检查

**使用方法**:
```cmd
# 检查环境
scripts\build.bat check

# 清理构建文件
scripts\build.bat clean

# 安装依赖
scripts\build.bat deps

# 启动开发模式
scripts\build.bat dev

# 构建 Windows 版本
scripts\build.bat build amd64

# 显示帮助
scripts\build.bat help
```

## 构建流程

### 开发模式
```bash
# 使用 Shell 脚本
./scripts/build.sh dev

# 使用 Makefile
make dev
```

开发模式会:
1. 检查环境
2. 安装依赖
3. 生成 Wails 绑定
4. 启动开发服务器

### 生产构建
```bash
# 构建当前平台
make build

# 或使用 Shell 脚本
./scripts/build.sh build
```

生产构建会:
1. 检查环境
2. 安装依赖
3. 编译前端代码
4. 编译 Go 代码
5. 打包应用程序

### 完整构建流程
```bash
# 使用 Shell 脚本
./scripts/build.sh all

# 使用 Makefile
make all
```

完整构建流程包括:
1. 清理旧文件
2. 安装依赖
3. 代码检查
4. 运行测试
5. 构建所有平台版本

## 平台支持

| 平台 | 架构 | 输出文件 | 构建命令 |
|------|------|----------|----------|
| macOS | universal | GenPulse.app | `make build-darwin` |
| macOS | amd64 | GenPulse.app | `./scripts/build.sh build darwin amd64` |
| macOS | arm64 | GenPulse.app | `./scripts/build.sh build darwin arm64` |
| Windows | amd64 | GenPulse.exe | `make build-windows` |
| Windows | arm64 | GenPulse.exe | `./scripts/build.sh build windows arm64` |
| Linux | amd64 | GenPulse | `make build-linux` |
| Linux | arm64 | GenPulse | `./scripts/build.sh build linux arm64` |

## 输出目录

构建结果保存在以下目录:

- `bin/` - 所有构建输出
  - `darwin-universal/` - macOS 通用二进制
  - `windows-amd64/` - Windows 64位
  - `linux-amd64/` - Linux 64位
- `build/` - Wails 临时构建文件
- `logs/` - 构建和测试日志

## 环境要求

### 必需工具
1. **Go** (1.20+)
   - 下载: https://golang.org/dl/
   - 验证: `go version`

2. **Wails CLI**
   - 安装: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
   - 验证: `wails version`

3. **Node.js** (16+)
   - 下载: https://nodejs.org/
   - 验证: `node --version` 和 `npm --version`

### 可选工具
1. **golangci-lint** (代码检查)
   - 安装: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

2. **Make** (Linux/macOS)
   - 通常已预装

## 故障排除

### 常见问题

1. **命令未找到**
   ```
   [ERROR] 命令 'wails' 未找到，请先安装
   ```
   解决方案: 确保 Go 二进制目录在 PATH 中
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

2. **前端依赖安装失败**
   ```
   npm ERR! 各种错误
   ```
   解决方案: 清理并重试
   ```bash
   ./scripts/build.sh clean
   ./scripts/build.sh deps
   ```

3. **Go 模块问题**
   ```
   go: 模块错误
   ```
   解决方案: 清理 Go 缓存
   ```bash
   go clean -modcache
   ./scripts/build.sh deps
   ```

### 日志文件
所有构建和测试日志都保存在 `logs/` 目录:
- `go_test.log` - Go 测试输出
- `frontend_test.log` - 前端测试输出
- `go_lint.log` - Go 代码检查输出
- `frontend_lint.log` - 前端代码检查输出

## 自动化集成

### CI/CD 示例 (GitHub Actions)
```yaml
name: Build

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Set up Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
    
    - name: Install Wails
      run: go install github.com/wailsapp/wails/v2/cmd/wails@latest
    
    - name: Build
      run: make all
```

### 开发工作流建议
1. 开发时使用 `make dev`
2. 提交前运行 `make lint` 和 `make test`
3. 发布时使用 `make all` 构建所有平台

## 更新日志

### v1.0.0
- 初始版本
- 支持 macOS、Windows、Linux 构建
- 提供 Shell 脚本、Makefile 和 Windows 批处理脚本
- 完整的开发和生产构建流程