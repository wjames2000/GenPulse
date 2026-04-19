#!/bin/bash

# 测试构建是否成功的简单脚本

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="$PROJECT_ROOT/build/bin"
BIN_DIR="$PROJECT_ROOT/bin"

echo "🔍 检查构建结果..."

# 检查构建文件是否存在
if [ -f "$BUILD_DIR/GenPulse" ]; then
    echo "✅ 构建成功: $BUILD_DIR/GenPulse"
    echo "   文件大小: $(ls -lh "$BUILD_DIR/GenPulse" | awk '{print $5}')"
    echo "   构建时间: $(stat -f "%Sm" "$BUILD_DIR/GenPulse")"
    
    # 检查文件类型
    echo "   文件类型: $(file "$BUILD_DIR/GenPulse" | cut -d: -f2-)"
else
    echo "❌ 构建失败: 未找到 GenPulse 可执行文件"
    exit 1
fi

# 检查是否复制到 bin 目录
if [ -f "$BIN_DIR/darwin-amd64/GenPulse" ]; then
    echo "✅ 已复制到 bin 目录: $BIN_DIR/darwin-amd64/GenPulse"
else
    echo "⚠️  未复制到 bin 目录"
fi

# 检查前端构建
if [ -d "$PROJECT_ROOT/frontend/dist" ]; then
    DIST_SIZE=$(du -sh "$PROJECT_ROOT/frontend/dist" | awk '{print $1}')
    echo "✅ 前端构建成功: $PROJECT_ROOT/frontend/dist ($DIST_SIZE)"
    
    # 检查关键文件
    if [ -f "$PROJECT_ROOT/frontend/dist/index.html" ]; then
        echo "   ✅ index.html 存在"
    fi
    if [ -f "$PROJECT_ROOT/frontend/dist/assets/index.*.js" ]; then
        JS_FILE=$(ls "$PROJECT_ROOT/frontend/dist/assets/index.*.js" 2>/dev/null | head -1)
        echo "   ✅ JavaScript 文件: $(basename "$JS_FILE")"
    fi
else
    echo "❌ 前端构建失败: 未找到 dist 目录"
fi

# 检查 Go 模块
echo "📦 Go 模块状态:"
cd "$PROJECT_ROOT"
go version
go mod verify 2>&1 | grep -v "go:" || true

# 检查前端依赖
echo "📦 前端依赖状态:"
cd "$PROJECT_ROOT/frontend"
npm --version
npm list --depth=0 2>&1 | grep -E "(react|vite|typescript|zustand)" || true

echo ""
echo "🎉 构建测试完成！"
echo "   应用已成功构建，可以运行开发模式进行测试:"
echo "   ./scripts/build.sh dev"
echo "   或"
echo "   make dev"