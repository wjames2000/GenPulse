#!/bin/bash

# GenPulse 集成测试运行脚本
# 用于运行端到端流水线测试和集成测试

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
    echo "GenPulse 集成测试运行脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help          显示此帮助信息"
    echo "  -a, --all           运行所有测试（单元测试+集成测试）"
    echo "  -u, --unit          只运行单元测试"
    echo "  -i, --integration   只运行集成测试"
    echo "  -e, --e2e           只运行端到端测试"
    echo "  -v, --verbose       显示详细输出"
    echo "  -c, --coverage      生成测试覆盖率报告"
    echo "  -p, --parallel      并行运行测试"
    echo "  -t, --timeout N     设置测试超时时间（秒）"
    echo ""
    echo "示例:"
    echo "  $0 --all            运行所有测试"
    echo "  $0 --integration    运行集成测试"
    echo "  $0 --e2e --verbose  运行端到端测试并显示详细输出"
}

# 默认参数
RUN_UNIT=false
RUN_INTEGRATION=false
RUN_E2E=false
VERBOSE=false
COVERAGE=false
PARALLEL=false
TIMEOUT=600

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -a|--all)
            RUN_UNIT=true
            RUN_INTEGRATION=true
            RUN_E2E=true
            shift
            ;;
        -u|--unit)
            RUN_UNIT=true
            shift
            ;;
        -i|--integration)
            RUN_INTEGRATION=true
            shift
            ;;
        -e|--e2e)
            RUN_E2E=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -p|--parallel)
            PARALLEL=true
            shift
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        *)
            log_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 如果没有指定任何测试类型，默认运行所有测试
if [ "$RUN_UNIT" = false ] && [ "$RUN_INTEGRATION" = false ] && [ "$RUN_E2E" = false ]; then
    RUN_UNIT=true
    RUN_INTEGRATION=true
    RUN_E2E=true
fi

# 构建测试参数
TEST_ARGS=""
if [ "$VERBOSE" = true ]; then
    TEST_ARGS="$TEST_ARGS -v"
fi
if [ "$COVERAGE" = true ]; then
    TEST_ARGS="$TEST_ARGS -cover"
fi
if [ "$PARALLEL" = true ]; then
    TEST_ARGS="$TEST_ARGS -parallel=4"
fi

# 设置环境变量
export GENPULSE_TEST_MODE=true
export GENPULSE_TEST_TIMEOUT=$TIMEOUT

# 打印测试配置
log_info "=== GenPulse 集成测试配置 ==="
log_info "运行单元测试: $RUN_UNIT"
log_info "运行集成测试: $RUN_INTEGRATION"
log_info "运行端到端测试: $RUN_E2E"
log_info "详细输出: $VERBOSE"
log_info "覆盖率报告: $COVERAGE"
log_info "并行测试: $PARALLEL"
log_info "超时时间: ${TIMEOUT}秒"
log_info "测试参数: $TEST_ARGS"
log_info "=============================="

# 记录开始时间
START_TIME=$(date +%s)

# 运行单元测试
if [ "$RUN_UNIT" = true ]; then
    log_info "开始运行单元测试..."
    
    # 工具模块单元测试
    log_info "运行工具模块单元测试..."
    if ! go test ./internal/genkit/tools/... $TEST_ARGS -timeout=${TIMEOUT}s; then
        log_error "工具模块单元测试失败"
        exit 1
    fi
    
    # 记忆模块单元测试
    log_info "运行记忆模块单元测试..."
    if ! go test ./internal/memory/... $TEST_ARGS -timeout=${TIMEOUT}s; then
        log_error "记忆模块单元测试失败"
        exit 1
    fi
    
    # 技能模块单元测试
    log_info "运行技能模块单元测试..."
    if ! go test ./internal/skills/... $TEST_ARGS -timeout=${TIMEOUT}s; then
        log_error "技能模块单元测试失败"
        exit 1
    fi
    
    # Agent模块单元测试
    log_info "运行Agent模块单元测试..."
    if ! go test ./internal/agents/... $TEST_ARGS -timeout=${TIMEOUT}s; then
        log_error "Agent模块单元测试失败"
        exit 1
    fi
    
    log_success "单元测试全部通过"
fi

# 运行集成测试
if [ "$RUN_INTEGRATION" = true ]; then
    log_info "开始运行集成测试..."
    
    # MCP集成测试
    log_info "运行MCP集成测试..."
    if ! go test ./internal/mcp/... $TEST_ARGS -timeout=${TIMEOUT}s; then
        log_error "MCP集成测试失败"
        exit 1
    fi
    
    # 监控集成测试
    log_info "运行监控集成测试..."
    if ! go test ./internal/monitoring/... $TEST_ARGS -timeout=${TIMEOUT}s; then
        log_error "监控集成测试失败"
        exit 1
    fi
    
    # Genkit集成测试
    log_info "运行Genkit集成测试..."
    if ! go test ./internal/genkit/... $TEST_ARGS -timeout=${TIMEOUT}s; then
        log_error "Genkit集成测试失败"
        exit 1
    fi
    
    log_success "集成测试全部通过"
fi

# 运行端到端测试
if [ "$RUN_E2E" = true ]; then
    log_info "开始运行端到端测试..."
    
    # 创建测试输出目录
    TEST_OUTPUT_DIR="./test_output"
    mkdir -p "$TEST_OUTPUT_DIR"
    
    # 设置测试输出环境变量
    export GENPULSE_TEST_OUTPUT_DIR="$TEST_OUTPUT_DIR"
    
    # 运行端到端流水线测试
    log_info "运行端到端流水线测试..."
    if ! go test ./test/... $TEST_ARGS -timeout=${TIMEOUT}s; then
        log_error "端到端流水线测试失败"
        
        # 显示测试输出目录内容
        if [ -d "$TEST_OUTPUT_DIR" ]; then
            log_info "测试输出目录内容:"
            ls -la "$TEST_OUTPUT_DIR"
        fi
        
        exit 1
    fi
    
    log_success "端到端测试全部通过"
    
    # 清理测试输出目录
    if [ -d "$TEST_OUTPUT_DIR" ]; then
        log_info "清理测试输出目录..."
        rm -rf "$TEST_OUTPUT_DIR"
    fi
fi

# 计算并显示测试耗时
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

log_success "=== 所有测试通过 ==="
log_info "总耗时: ${DURATION}秒"

# 生成测试覆盖率报告
if [ "$COVERAGE" = true ]; then
    log_info "生成测试覆盖率报告..."
    
    # 创建覆盖率报告目录
    COVERAGE_DIR="./coverage"
    mkdir -p "$COVERAGE_DIR"
    
    # 生成覆盖率数据
    go test ./... -coverprofile="$COVERAGE_DIR/coverage.out" -covermode=atomic
    
    # 生成HTML报告
    go tool cover -html="$COVERAGE_DIR/coverage.out" -o "$COVERAGE_DIR/coverage.html"
    
    # 生成文本报告
    go tool cover -func="$COVERAGE_DIR/coverage.out" > "$COVERAGE_DIR/coverage.txt"
    
    # 显示覆盖率摘要
    log_info "测试覆盖率报告:"
    tail -n 5 "$COVERAGE_DIR/coverage.txt"
    
    log_success "覆盖率报告已生成到: $COVERAGE_DIR/"
fi

exit 0