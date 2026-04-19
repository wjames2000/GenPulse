#!/bin/bash

# 快速构建测试脚本

set -e

echo "🚀 快速构建测试..."
echo "=================="

# 检查基本环境
echo "1. 检查环境..."
command -v go >/dev/null 2>&1 && echo "   ✅ Go: $(go version)" || echo "   ❌ Go 未安装"
command -v node >/dev/null 2>&1 && echo "   ✅ Node.js: $(node --version)" || echo "   ❌ Node.js 未安装"

# 检查 Wails
if command -v wails >/dev/null 2>&1; then
    echo "   ✅ Wails: $(wails version 2>/dev/null | head -1)"
else
    # 尝试从 GOBIN 查找
    GOBIN="$(go env GOPATH)/bin"
    if [ -f "$GOBIN/wails" ]; then
        export PATH="$GOBIN:$PATH"
        echo "   ✅ Wails: $(wails version 2>/dev/null | head -1) (从 GOBIN 找到)"
    else
        echo "   ❌ Wails 未安装"
    fi
fi

echo ""
echo "2. 清理旧构建..."
rm -rf frontend/dist build/bin 2>/dev/null || true
echo "   ✅ 清理完成"

echo ""
echo "3. 构建前端..."
cd frontend
npm run build 2>&1 | tail -5
echo "   ✅ 前端构建完成"

echo ""
echo "4. 构建应用..."
cd ..
echo "   执行: wails build -s"  # -s 跳过前端构建（因为我们已经构建了）
wails build -s 2>&1 | grep -E "(Building|Generating|Compiling|Done|ERROR|Error)" || true

echo ""
echo "5. 检查结果..."
if [ -f "build/bin/GenPulse" ]; then
    echo "   ✅ 构建成功!"
    echo "   文件: build/bin/GenPulse"
    echo "   大小: $(ls -lh build/bin/GenPulse | awk '{print $5}')"
    echo "   类型: $(file build/bin/GenPulse | cut -d: -f2-)"
    
    # 复制到 bin 目录
    mkdir -p bin/darwin-amd64
    cp build/bin/GenPulse bin/darwin-amd64/ 2>/dev/null || true
    echo "   已复制到: bin/darwin-amd64/GenPulse"
else
    echo "   ❌ 构建失败"
    exit 1
fi

echo ""
echo "🎉 快速构建测试完成!"
echo "   应用已成功构建，可以运行开发模式进行完整测试:"
echo "   ./scripts/build.sh dev"