#!/bin/bash

# GenPulse 编译脚本
# 支持开发、构建、清理和跨平台打包

set -e  # 遇到错误时退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="$PROJECT_ROOT/frontend"
BUILD_DIR="$PROJECT_ROOT/build"
BIN_DIR="$PROJECT_ROOT/bin"
LOG_DIR="$PROJECT_ROOT/logs"

# Go 环境变量
export GOBIN="$(go env GOPATH)/bin"
export PATH="$GOBIN:$PATH"

# 确保必要的目录存在
mkdir -p "$BUILD_DIR" "$BIN_DIR" "$LOG_DIR"

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令是否存在
check_command() {
    if ! command -v "$1" &> /dev/null; then
        log_error "命令 '$1' 未找到，请先安装"
        return 1
    fi
    return 0
}

# 检查环境
check_environment() {
    log_info "检查开发环境..."
    
    # 检查 Go
    if check_command "go"; then
        GO_VERSION=$(go version | awk '{print $3}')
        log_info "Go 版本: $GO_VERSION"
    else
        log_error "请安装 Go (https://golang.org/dl/)"
        exit 1
    fi
    
    # 检查 Wails
    # 首先尝试直接查找
    if command -v wails &> /dev/null; then
        WAILS_VERSION=$(wails version 2>/dev/null | head -1 || echo "未知")
        log_info "Wails 版本: $WAILS_VERSION"
    # 然后检查 Go 二进制目录
    elif [ -f "$GOBIN/wails" ]; then
        export PATH="$GOBIN:$PATH"
        WAILS_VERSION=$(wails version 2>/dev/null | head -1 || echo "未知")
        log_info "Wails 版本: $WAILS_VERSION (在 GOBIN 中找到)"
    else
        log_error "Wails 未找到，请安装: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
        log_info "安装后请确保 $(go env GOPATH)/bin 在 PATH 环境变量中"
        exit 1
    fi
    
    # 检查 Node.js
    if check_command "node"; then
        NODE_VERSION=$(node --version)
        NPM_VERSION=$(npm --version)
        log_info "Node.js 版本: $NODE_VERSION"
        log_info "npm 版本: $NPM_VERSION"
    else
        log_error "请安装 Node.js (https://nodejs.org/)"
        exit 1
    fi
    
    log_success "环境检查完成"
}

# 清理函数
clean() {
    log_info "清理构建文件..."
    
    # 清理前端构建文件
    if [ -d "$FRONTEND_DIR/dist" ]; then
        rm -rf "$FRONTEND_DIR/dist"
        log_info "已清理前端 dist 目录"
    fi
    
    if [ -d "$FRONTEND_DIR/node_modules" ]; then
        rm -rf "$FRONTEND_DIR/node_modules"
        log_info "已清理前端 node_modules"
    fi
    
    # 清理 Go 构建文件
    if [ -f "$PROJECT_ROOT/go.sum" ]; then
        rm -f "$PROJECT_ROOT/go.sum"
        log_info "已清理 go.sum"
    fi
    
    # 清理二进制文件
    if [ -d "$BIN_DIR" ]; then
        rm -rf "$BIN_DIR"/*
        log_info "已清理 bin 目录"
    fi
    
    # 清理日志文件
    if [ -d "$LOG_DIR" ]; then
        rm -rf "$LOG_DIR"/*.log
        log_info "已清理日志文件"
    fi
    
    # 清理 macOS 特定文件
    find "$PROJECT_ROOT" -name ".DS_Store" -delete
    find "$PROJECT_ROOT" -name "*.app" -type d -exec rm -rf {} + 2>/dev/null || true
    
    log_success "清理完成"
}

# 安装依赖
install_deps() {
    log_info "安装依赖..."
    
    # 安装 Go 依赖
    log_info "安装 Go 依赖..."
    cd "$PROJECT_ROOT"
    go mod tidy
    go mod download
    
    # 安装前端依赖
    log_info "安装前端依赖..."
    cd "$FRONTEND_DIR"
    npm install --no-audit --no-fund
    
    log_success "依赖安装完成"
}

# 开发模式
dev() {
    log_info "启动开发模式..."
    check_environment
    
    # 安装依赖
    install_deps
    
    # 启动开发服务器
    log_info "启动开发服务器..."
    cd "$PROJECT_ROOT"
    wails dev
}

# 构建函数
build() {
    local PLATFORM=$1
    local ARCH=$2
    
    log_info "构建 $PLATFORM/$ARCH 版本..."
    check_environment
    
    # 安装依赖
    install_deps
    
    # 确保前端已构建
    log_info "构建前端..."
    cd "$FRONTEND_DIR"
    npm run build
    
    # 构建命令
    cd "$PROJECT_ROOT"
    
    local BUILD_CMD="wails build"
    local OUTPUT_NAME="GenPulse"
    
    # 根据平台设置输出文件名
    case "$PLATFORM" in
        "darwin")
            OUTPUT_NAME="$OUTPUT_NAME.app"
            ;;
        "windows")
            OUTPUT_NAME="$OUTPUT_NAME.exe"
            ;;
        "linux")
            OUTPUT_NAME="$OUTPUT_NAME"
            ;;
    esac
    
    # 设置平台参数
    if [ -n "$ARCH" ] && [ "$ARCH" != "x86_64" ] && [ "$ARCH" != "amd64" ]; then
        # 对于非默认架构，需要指定平台/架构组合
        BUILD_CMD="$BUILD_CMD -platform $PLATFORM/$ARCH"
    elif [ "$PLATFORM" != "$(uname -s | tr '[:upper:]' '[:lower:]')" ]; then
        # 跨平台构建
        BUILD_CMD="$BUILD_CMD -platform $PLATFORM"
    fi
    
    # 执行构建
    log_info "执行构建命令: $BUILD_CMD"
    if eval "$BUILD_CMD"; then
        log_success "构建成功"
        
        # 复制到 bin 目录
        local BUILD_OUTPUT_DIR="$PROJECT_ROOT/build/bin"
        if [ -d "$BUILD_OUTPUT_DIR" ]; then
            mkdir -p "$BIN_DIR/$PLATFORM-$ARCH"
            
            # 复制所有构建文件
            cp -r "$BUILD_OUTPUT_DIR"/* "$BIN_DIR/$PLATFORM-$ARCH/" 2>/dev/null || true
            
            # 检查是否复制成功
            if [ "$(ls -A "$BIN_DIR/$PLATFORM-$ARCH/" 2>/dev/null | wc -l)" -gt 0 ]; then
                log_info "已复制到: $BIN_DIR/$PLATFORM-$ARCH/"
                log_info "构建文件: $(ls "$BIN_DIR/$PLATFORM-$ARCH/")"
            else
                log_warning "构建文件未复制到 bin 目录"
            fi
        else
            log_warning "未找到构建输出目录: $BUILD_OUTPUT_DIR"
        fi
    else
        log_error "构建失败"
        exit 1
    fi
}

# 构建所有平台
build_all() {
    log_info "构建所有平台版本..."
    
    # macOS
    build darwin universal
    build darwin amd64
    build darwin arm64
    
    # Windows
    build windows amd64
    
    # Linux
    build linux amd64
    build linux arm64
    
    log_success "所有平台构建完成"
    log_info "构建结果保存在: $BIN_DIR/"
}

# 运行测试
test() {
    log_info "运行测试..."
    
    # Go 测试
    log_info "运行 Go 测试..."
    cd "$PROJECT_ROOT"
    go test ./... -v 2>&1 | tee "$LOG_DIR/go_test.log"
    
    # 前端测试（如果有）
    if [ -f "$FRONTEND_DIR/package.json" ]; then
        log_info "运行前端测试..."
        cd "$FRONTEND_DIR"
        if grep -q "\"test\":" package.json; then
            npm test 2>&1 | tee "$LOG_DIR/frontend_test.log"
        else
            log_warning "前端 package.json 中没有找到 test 脚本"
        fi
    fi
    
    log_success "测试完成"
}

# 代码检查
lint() {
    log_info "运行代码检查..."
    
    # Go 代码检查
    if check_command "golangci-lint"; then
        log_info "运行 golangci-lint..."
        cd "$PROJECT_ROOT"
        golangci-lint run 2>&1 | tee "$LOG_DIR/go_lint.log"
    else
        log_warning "golangci-lint 未安装，跳过 Go 代码检查"
        log_info "安装命令: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    fi
    
    # 前端代码检查
    if [ -f "$FRONTEND_DIR/package.json" ]; then
        cd "$FRONTEND_DIR"
        if grep -q "\"lint\":" package.json; then
            log_info "运行前端代码检查..."
            npm run lint 2>&1 | tee "$LOG_DIR/frontend_lint.log"
        else
            log_warning "前端 package.json 中没有找到 lint 脚本"
        fi
    fi
    
    log_success "代码检查完成"
}

# 显示帮助信息
show_help() {
    echo -e "${BLUE}GenPulse 编译脚本${NC}"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  check         检查开发环境"
    echo "  clean         清理构建文件"
    echo "  deps          安装依赖"
    echo "  dev           启动开发模式"
    echo "  build [平台] [架构]  构建指定平台版本"
    echo "                  平台: darwin, windows, linux (默认: 当前平台)"
    echo "                  架构: amd64, arm64, universal (默认: 当前架构)"
    echo "  build-all     构建所有平台版本"
    echo "  test          运行测试"
    echo "  lint          运行代码检查"
    echo "  all           执行完整构建流程 (clean -> deps -> lint -> test -> build-all)"
    echo "  help          显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 check              # 检查环境"
    echo "  $0 dev                # 启动开发模式"
    echo "  $0 build darwin       # 构建 macOS 版本"
    echo "  $0 build windows amd64 # 构建 Windows 64位版本"
    echo "  $0 all                # 执行完整构建流程"
    echo ""
}

# 完整构建流程
all() {
    log_info "开始完整构建流程..."
    
    # 检查环境
    check_environment
    
    # 清理
    clean
    
    # 安装依赖
    install_deps
    
    # 代码检查
    lint
    
    # 运行测试
    test
    
    # 构建所有平台
    build_all
    
    log_success "完整构建流程完成"
    log_info "构建结果保存在: $BIN_DIR/"
}

# 主函数
main() {
    local COMMAND=$1
    local PLATFORM=$2
    local ARCH=$3
    
    # 转换架构名称
    if [ "$ARCH" = "x86_64" ]; then
        ARCH="amd64"
    fi
    
    case "$COMMAND" in
        "check")
            check_environment
            ;;
        "clean")
            clean
            ;;
        "deps")
            install_deps
            ;;
        "dev")
            dev
            ;;
        "build")
            build "${PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}" "${ARCH:-$(uname -m)}"
            ;;
        "build-all")
            build_all
            ;;
        "test")
            test
            ;;
        "lint")
            lint
            ;;
        "all")
            all
            ;;
        "help"|"-h"|"--help"|"")
            show_help
            ;;
        *)
            log_error "未知命令: $COMMAND"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"