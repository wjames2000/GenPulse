#!/bin/bash

# Genkit运行时集成演示脚本

set -e

echo "🚀 Genkit运行时集成演示"
echo "========================"

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# 检查Go环境
check_go() {
    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go"
        exit 1
    fi
    log_info "Go版本: $(go version)"
}

# 清理和构建
build_project() {
    log_info "构建项目..."
    
    # 清理
    rm -rf frontend/dist build/bin 2>/dev/null || true
    
    # 安装依赖
    cd frontend && npm run build
    cd ..
    
    # 构建应用
    go build -o build/bin/GenPulse
    
    if [ -f "build/bin/GenPulse" ]; then
        log_success "构建成功: build/bin/GenPulse"
    else
        log_error "构建失败"
        exit 1
    fi
}

# 运行集成测试
run_integration_test() {
    log_info "运行集成测试..."
    
    # 运行Go测试
    if go test ./internal/genkit -v 2>&1 | grep -q "PASS"; then
        log_success "集成测试通过"
    else
        log_warning "集成测试有警告或失败"
    fi
}

# 显示系统架构
show_architecture() {
    echo ""
    echo "📁 项目结构概览:"
    echo "----------------"
    
    echo "internal/genkit/"
    find internal/genkit -type f -name "*.go" | sort | while read file; do
        filename=$(basename "$file")
        dirname=$(dirname "$file" | sed "s|internal/genkit/||")
        if [ "$dirname" = "internal/genkit" ]; then
            echo "  ├── $filename"
        else
            echo "  ├── ${dirname#*/}/$filename"
        fi
    done | head -20
    
    echo ""
    echo "🧩 核心组件:"
    echo "------------"
    echo "1. Genkit管理器 (init.go) - 统一初始化入口"
    echo "2. 模型适配器层 (models/) - 统一AI模型接口"
    echo "   ├── gemini_adapter.go - Google Gemini"
    echo "   ├── gpt_adapter.go    - OpenAI GPT"
    echo "   ├── claude_adapter.go - Anthropic Claude"
    echo "   ├── ollama_adapter.go - 本地Ollama"
    echo "   └── custom_adapter.go - 自定义模型"
    echo "3. 工具注册表 (tools/) - 可扩展工具系统"
    echo "   ├── tool_registry.go - 工具管理核心"
    echo "   └── fs_tools.go     - 文件系统工具"
    echo "4. Flow引擎 (flows/) - 工作流编排"
    echo "   └── flow_engine.go - Flow执行引擎"
    echo "5. 配置管理 (config/) - 统一配置"
    echo "   └── app_config.go  - 应用配置"
}

# 显示使用示例
show_examples() {
    echo ""
    echo "💡 使用示例:"
    echo "-----------"
    
    cat << 'EOF'
// 1. 初始化Genkit系统
manager := genkit.NewGenkitManager()
err := manager.Initialize(ctx)

// 2. 使用模型适配器
factory := &models.DefaultModelAdapterFactory{}
adapter := models.NewUnifiedModelAdapter(factory)

// 注册模型
config := models.ModelConfig{
    Type:        models.ModelTypeGemini,
    Name:        "gemini-1.5-pro",
    Provider:    "google",
    APIKey:      "your-api-key",
}
adapter.RegisterModel(config)

// 调用模型
response, err := adapter.Generate(ctx, "gemini-1.5-pro", models.ModelRequest{
    Prompt: "Hello, world!",
})

// 3. 使用工具系统
registry := tools.NewToolRegistry()
fsTool, _ := tools.NewFSTool("/workspace")
registry.RegisterTool(fsTool)

// 执行工具
result, err := registry.ExecuteTool(ctx, tools.ToolExecution{
    ToolID: "fs_tool",
    Parameters: map[string]interface{}{
        "operation": "write",
        "path":      "test.txt",
        "content":   "Hello from tool!",
    },
})

// 4. 使用Flow引擎
flowEngine := flows.NewFlowEngine(adapter, registry)

// 定义和执行Flow
flow := flows.FlowDefinition{
    ID:   "demo-flow",
    Name: "演示Flow",
    Type: flows.FlowTypeSequential,
    Nodes: []flows.NodeDefinition{
        {
            ID:   "analyze",
            Name: "分析",
            Type: flows.NodeTypeModel,
            Config: map[string]interface{}{
                "model_id": "gemini-1.5-pro",
                "prompt":   "分析需求",
            },
        },
    },
}

flowEngine.RegisterFlow(flow)
execution, err := flowEngine.ExecuteFlow(ctx, "demo-flow", nil)
EOF
}

# 显示下一步计划
show_next_steps() {
    echo ""
    echo "🎯 下一步开发计划:"
    echo "-----------------"
    echo "1. 完善模型适配器 - 集成实际API调用"
    echo "2. 扩展工具系统 - 添加Git、Shell、项目管理工具"
    echo "3. 实现Agent框架 - 基于Genkit的多Agent系统"
    echo "4. 添加技能系统 - 经验沉淀和复用"
    echo "5. 实现记忆系统 - 三层记忆架构"
    echo "6. 集成MCP协议 - 模型上下文协议"
    echo "7. 添加监控系统 - OpenTelemetry集成"
    echo ""
    echo "📚 相关文档:"
    echo "-----------"
    echo "• WBS任务分解: docs/WBS 任务分解.md"
    echo "• 项目结构规划: docs/项目结构规划.md"
    echo "• Genkit文档: https://firebase.google.com/docs/genkit"
}

# 主函数
main() {
    echo ""
    log_info "开始Genkit运行时集成演示"
    
    # 检查环境
    check_go
    
    # 显示架构
    show_architecture
    
    # 构建项目
    build_project
    
    # 运行测试
    run_integration_test
    
    # 显示示例
    show_examples
    
    # 显示下一步
    show_next_steps
    
    echo ""
    log_success "演示完成！"
    echo ""
    echo "🚀 下一步: 运行应用测试Genkit功能"
    echo "   ./build/bin/GenPulse"
    echo ""
    echo "🔧 开发模式:"
    echo "   make dev"
    echo ""
    echo "📊 查看完整测试:"
    echo "   go test ./internal/genkit -v"
}

# 运行主函数
main