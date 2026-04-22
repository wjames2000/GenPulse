#!/bin/bash

# GenPulse 跨平台兼容性测试脚本
# 用于测试在不同平台上的兼容性

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# 显示帮助信息
show_help() {
    echo "GenPulse 跨平台兼容性测试脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help          显示此帮助信息"
    echo "  -p, --platform PLAT 测试指定平台 (darwin, linux, windows)"
    echo "  -a, --all           测试所有平台"
    echo "  -b, --build         构建跨平台二进制文件"
    echo "  -t, --test          运行跨平台测试"
    echo "  -v, --verbose       显示详细输出"
    echo ""
    echo "示例:"
    echo "  $0 --all            测试所有平台"
    echo "  $0 --platform linux 测试Linux平台"
    echo "  $0 --build --test   构建并测试"
}

# 检测当前平台
detect_platform() {
    case "$(uname -s)" in
        Darwin)
            echo "darwin"
            ;;
        Linux)
            echo "linux"
            ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*)
            echo "windows"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# 检测当前架构
detect_arch() {
    case "$(uname -m)" in
        x86_64)
            echo "amd64"
            ;;
        arm64|aarch64)
            echo "arm64"
            ;;
        i386|i686)
            echo "386"
            ;;
        armv7l)
            echo "armv7"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# 运行平台特定测试
run_platform_tests() {
    local platform=$1
    local arch=$2
    
    log_info "在 ${platform}/${arch} 上运行测试..."
    
    # 设置平台特定的环境变量
    export GOOS="$platform"
    export GOARCH="$arch"
    
    # 创建测试目录
    TEST_DIR="./cross_platform_test/${platform}_${arch}"
    mkdir -p "$TEST_DIR"
    
    # 运行基础功能测试
    log_info "运行基础功能测试..."
    
    # 测试文件系统操作
    if [ "$platform" = "windows" ]; then
        # Windows平台的特殊测试
        log_info "运行Windows平台特殊测试..."
        # 测试Windows路径处理
        go test ./internal/genkit/tools -run "TestFsTools.*" -v 2>&1 | tee "$TEST_DIR/fs_tests.log"
    else
        # Unix-like平台测试
        go test ./internal/genkit/tools -run "TestFsTools.*" -v 2>&1 | tee "$TEST_DIR/fs_tests.log"
    fi
    
    # 测试命令行执行
    log_info "运行命令行执行测试..."
    go test ./internal/genkit/tools -run "TestShellTools.*" -v 2>&1 | tee "$TEST_DIR/shell_tests.log"
    
    # 测试Git操作
    log_info "运行Git操作测试..."
    go test ./internal/genkit/tools -run "TestGitTools.*" -v 2>&1 | tee "$TEST_DIR/git_tests.log"
    
    # 检查测试结果
    if grep -q "FAIL" "$TEST_DIR"/*.log; then
        log_error "${platform}/${arch} 测试失败"
        return 1
    else
        log_success "${platform}/${arch} 测试通过"
        return 0
    fi
}

# 构建跨平台二进制文件
build_cross_platform() {
    local platform=$1
    local arch=$2
    
    log_info "为 ${platform}/${arch} 构建二进制文件..."
    
    # 设置构建环境
    export GOOS="$platform"
    export GOARCH="$arch"
    
    # 创建构建目录
    BUILD_DIR="./build/${platform}_${arch}"
    mkdir -p "$BUILD_DIR"
    
    # 构建主程序
    log_info "构建主程序..."
    if ! go build -o "$BUILD_DIR/genpulse" ./cmd/genpulse 2>&1 | tee "$BUILD_DIR/build.log"; then
        log_error "${platform}/${arch} 构建失败"
        return 1
    fi
    
    # 验证构建结果
    if [ -f "$BUILD_DIR/genpulse" ]; then
        log_success "${platform}/${arch} 构建成功"
        
        # 显示文件信息
        file "$BUILD_DIR/genpulse" | tee -a "$BUILD_DIR/build.log"
        
        # 如果是当前平台，测试运行
        CURRENT_PLATFORM=$(detect_platform)
        CURRENT_ARCH=$(detect_arch)
        
        if [ "$platform" = "$CURRENT_PLATFORM" ] && [ "$arch" = "$CURRENT_ARCH" ]; then
            log_info "测试运行构建的二进制文件..."
            if "$BUILD_DIR/genpulse" --version 2>&1 | tee -a "$BUILD_DIR/build.log"; then
                log_success "二进制文件运行测试通过"
            else
                log_warning "二进制文件运行测试失败（可能是预期行为）"
            fi
        fi
        
        return 0
    else
        log_error "${platform}/${arch} 构建失败：未生成二进制文件"
        return 1
    fi
}

# 清理函数
cleanup() {
    log_info "清理测试文件..."
    rm -rf ./cross_platform_test ./build
}

# 主函数
main() {
    # 默认参数
    PLATFORM=""
    TEST_ALL=false
    DO_BUILD=false
    DO_TEST=false
    VERBOSE=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -p|--platform)
                PLATFORM="$2"
                shift 2
                ;;
            -a|--all)
                TEST_ALL=true
                shift
                ;;
            -b|--build)
                DO_BUILD=true
                shift
                ;;
            -t|--test)
                DO_TEST=true
                shift
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            *)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 如果没有指定操作，默认运行测试
    if [ "$DO_BUILD" = false ] && [ "$DO_TEST" = false ]; then
        DO_TEST=true
    fi
    
    # 检测当前平台
    CURRENT_PLATFORM=$(detect_platform)
    CURRENT_ARCH=$(detect_arch)
    
    log_info "当前平台: ${CURRENT_PLATFORM}/${CURRENT_ARCH}"
    
    # 定义要测试的平台和架构组合
    declare -A PLATFORMS
    
    # 根据参数设置要测试的平台
    if [ -n "$PLATFORM" ]; then
        # 测试指定平台
        case "$PLATFORM" in
            darwin)
                PLATFORMS["darwin"]="amd64 arm64"
                ;;
            linux)
                PLATFORMS["linux"]="amd64 arm64 386 armv7"
                ;;
            windows)
                PLATFORMS["windows"]="amd64 386"
                ;;
            *)
                log_error "不支持的平台: $PLATFORM"
                exit 1
                ;;
        esac
    elif [ "$TEST_ALL" = true ]; then
        # 测试所有平台
        PLATFORMS["darwin"]="amd64 arm64"
        PLATFORMS["linux"]="amd64 arm64 386 armv7"
        PLATFORMS["windows"]="amd64 386"
    else
        # 只测试当前平台
        PLATFORMS["$CURRENT_PLATFORM"]="$CURRENT_ARCH"
    fi
    
    # 打印测试配置
    log_info "=== 跨平台兼容性测试配置 ==="
    log_info "构建: $DO_BUILD"
    log_info "测试: $DO_TEST"
    log_info "详细输出: $VERBOSE"
    log_info "测试平台:"
    for platform in "${!PLATFORMS[@]}"; do
        for arch in ${PLATFORMS[$platform]}; do
            log_info "  - ${platform}/${arch}"
        done
    done
    log_info "============================"
    
    # 记录开始时间
    START_TIME=$(date +%s)
    
    # 清理旧的测试文件
    cleanup
    
    # 执行构建和测试
    FAILED=0
    TOTAL=0
    
    for platform in "${!PLATFORMS[@]}"; do
        for arch in ${PLATFORMS[$platform]}; do
            TOTAL=$((TOTAL + 1))
            
            log_info "处理 ${platform}/${arch}..."
            
            # 构建
            if [ "$DO_BUILD" = true ]; then
                if ! build_cross_platform "$platform" "$arch"; then
                    FAILED=$((FAILED + 1))
                    continue
                fi
            fi
            
            # 测试
            if [ "$DO_TEST" = true ]; then
                if ! run_platform_tests "$platform" "$arch"; then
                    FAILED=$((FAILED + 1))
                fi
            fi
        done
    done
    
    # 计算并显示测试结果
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))
    
    SUCCESS=$((TOTAL - FAILED))
    
    log_info "=== 跨平台兼容性测试结果 ==="
    log_info "总计: $TOTAL"
    log_info "成功: $SUCCESS"
    log_info "失败: $FAILED"
    log_info "耗时: ${DURATION}秒"
    
    if [ $FAILED -eq 0 ]; then
        log_success "所有跨平台测试通过！"
    else
        log_error "有 $FAILED 个测试失败"
        
        # 显示失败详情
        if [ -d "./cross_platform_test" ]; then
            log_info "失败详情:"
            find ./cross_platform_test -name "*.log" -exec grep -l "FAIL" {} \; | while read -r file; do
                log_warning "失败文件: $file"
                tail -20 "$file"
            done
        fi
        
        exit 1
    fi
    
    # 生成测试报告
    if [ $TOTAL -gt 1 ]; then
        log_info "生成跨平台测试报告..."
        
        REPORT_FILE="./cross_platform_report.md"
        cat > "$REPORT_FILE" << EOF
# GenPulse 跨平台兼容性测试报告

## 测试概况
- 测试时间: $(date)
- 总测试平台: $TOTAL
- 成功平台: $SUCCESS
- 失败平台: $FAILED
- 测试耗时: ${DURATION}秒

## 测试详情

EOF
        
        # 添加每个平台的测试结果
        for platform in "${!PLATFORMS[@]}"; do
            for arch in ${PLATFORMS[$platform]}; do
                TEST_DIR="./cross_platform_test/${platform}_${arch}"
                if [ -d "$TEST_DIR" ]; then
                    echo "### ${platform}/${arch}" >> "$REPORT_FILE"
                    echo "" >> "$REPORT_FILE"
                    
                    # 检查测试结果
                    if find "$TEST_DIR" -name "*.log" -exec grep -l "FAIL" {} \; > /dev/null; then
                        echo "❌ 测试失败" >> "$REPORT_FILE"
                    else
                        echo "✅ 测试通过" >> "$REPORT_FILE"
                    fi
                    echo "" >> "$REPORT_FILE"
                    
                    # 添加测试日志摘要
                    for log_file in "$TEST_DIR"/*.log; do
                        if [ -f "$log_file" ]; then
                            echo "#### $(basename "$log_file")" >> "$REPORT_FILE"
                            echo '```' >> "$REPORT_FILE"
                            tail -10 "$log_file" >> "$REPORT_FILE"
                            echo '```' >> "$REPORT_FILE"
                            echo "" >> "$REPORT_FILE"
                        fi
                    done
                fi
            done
        done
        
        log_success "测试报告已生成: $REPORT_FILE"
    fi
    
    exit 0
}

# 运行主函数
main "$@"