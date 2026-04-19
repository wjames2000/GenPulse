# GenPulse Makefile
# 跨平台构建工具

.PHONY: help check clean deps dev build build-all test lint all

# 项目变量
PROJECT_ROOT := $(shell pwd)
FRONTEND_DIR := $(PROJECT_ROOT)/frontend
BUILD_DIR := $(PROJECT_ROOT)/build
BIN_DIR := $(PROJECT_ROOT)/bin
LOG_DIR := $(PROJECT_ROOT)/logs

# 颜色定义
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m

# 帮助信息
help:
	@echo "$(BLUE)GenPulse 构建工具$(NC)"
	@echo ""
	@echo "可用命令:"
	@echo "  make check         检查开发环境"
	@echo "  make clean         清理构建文件"
	@echo "  make deps          安装依赖"
	@echo "  make dev           启动开发模式"
	@echo "  make build         构建当前平台版本"
	@echo "  make build-all     构建所有平台版本"
	@echo "  make test          运行测试"
	@echo "  make lint          运行代码检查"
	@echo "  make all           执行完整构建流程"
	@echo ""
	@echo "平台特定构建:"
	@echo "  make build-darwin  构建 macOS 版本"
	@echo "  make build-windows 构建 Windows 版本"
	@echo "  make build-linux   构建 Linux 版本"
	@echo ""
	@echo "示例:"
	@echo "  make check          # 检查环境"
	@echo "  make dev            # 启动开发模式"
	@echo "  make build-darwin   # 构建 macOS 版本"
	@echo "  make all            # 执行完整构建流程"

# 检查环境
check:
	@echo "$(BLUE)[INFO]$(NC) 检查开发环境..."
	@command -v go >/dev/null 2>&1 || { echo "$(RED)[ERROR]$(NC) Go 未安装，请先安装: https://golang.org/dl/"; exit 1; }
	@echo "$(BLUE)[INFO]$(NC) Go 版本: $(shell go version | awk '{print $$3}')"
	@command -v wails >/dev/null 2>&1 || { echo "$(RED)[ERROR]$(NC) Wails 未安装，请安装: go install github.com/wailsapp/wails/v2/cmd/wails@latest"; exit 1; }
	@echo "$(BLUE)[INFO]$(NC) Wails 版本: $(shell wails version 2>/dev/null | head -1 || echo "未知")"
	@command -v node >/dev/null 2>&1 || { echo "$(RED)[ERROR]$(NC) Node.js 未安装，请先安装: https://nodejs.org/"; exit 1; }
	@echo "$(BLUE)[INFO]$(NC) Node.js 版本: $(shell node --version)"
	@echo "$(BLUE)[INFO]$(NC) npm 版本: $(shell npm --version)"
	@echo "$(GREEN)[SUCCESS]$(NC) 环境检查完成"

# 清理
clean:
	@echo "$(BLUE)[INFO]$(NC) 清理构建文件..."
	@rm -rf $(FRONTEND_DIR)/dist
	@rm -rf $(FRONTEND_DIR)/node_modules
	@rm -f $(PROJECT_ROOT)/go.sum
	@rm -rf $(BIN_DIR)/*
	@rm -rf $(LOG_DIR)/*.log
	@find $(PROJECT_ROOT) -name ".DS_Store" -delete
	@find $(PROJECT_ROOT) -name "*.app" -type d -exec rm -rf {} + 2>/dev/null || true
	@echo "$(GREEN)[SUCCESS]$(NC) 清理完成"

# 安装依赖
deps: check
	@echo "$(BLUE)[INFO]$(NC) 安装依赖..."
	@cd $(PROJECT_ROOT) && go mod tidy && go mod download
	@cd $(FRONTEND_DIR) && npm install --no-audit --no-fund
	@echo "$(GREEN)[SUCCESS]$(NC) 依赖安装完成"

# 开发模式
dev: deps
	@echo "$(BLUE)[INFO]$(NC) 启动开发模式..."
	@cd $(PROJECT_ROOT) && wails dev

# 构建当前平台
build: deps
	@echo "$(BLUE)[INFO]$(NC) 构建当前平台版本..."
	@echo "$(BLUE)[INFO]$(NC) 构建前端..."
	@cd $(FRONTEND_DIR) && npm run build
	@cd $(PROJECT_ROOT) && wails build
	@mkdir -p $(BIN_DIR)/$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed 's/x86_64/amd64/')
	@cp -r $(BUILD_DIR)/bin/* $(BIN_DIR)/$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed 's/x86_64/amd64/')/ 2>/dev/null || true
	@echo "$(GREEN)[SUCCESS]$(NC) 构建完成，输出在: $(BIN_DIR)/$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed 's/x86_64/amd64/')/"

# 构建 macOS
build-darwin: deps
	@echo "$(BLUE)[INFO]$(NC) 构建 macOS 版本..."
	@cd $(PROJECT_ROOT) && wails build -platform darwin -arch universal
	@mkdir -p $(BIN_DIR)/darwin-universal
	@cp -r $(BUILD_DIR)/GenPulse.app $(BIN_DIR)/darwin-universal/ 2>/dev/null || true
	@echo "$(GREEN)[SUCCESS]$(NC) macOS 构建完成"

# 构建 Windows
build-windows: deps
	@echo "$(BLUE)[INFO]$(NC) 构建 Windows 版本..."
	@cd $(PROJECT_ROOT) && wails build -platform windows -arch amd64
	@mkdir -p $(BIN_DIR)/windows-amd64
	@cp $(BUILD_DIR)/GenPulse.exe $(BIN_DIR)/windows-amd64/ 2>/dev/null || true
	@echo "$(GREEN)[SUCCESS]$(NC) Windows 构建完成"

# 构建 Linux
build-linux: deps
	@echo "$(BLUE)[INFO]$(NC) 构建 Linux 版本..."
	@cd $(PROJECT_ROOT) && wails build -platform linux -arch amd64
	@mkdir -p $(BIN_DIR)/linux-amd64
	@cp $(BUILD_DIR)/GenPulse $(BIN_DIR)/linux-amd64/ 2>/dev/null || true
	@echo "$(GREEN)[SUCCESS]$(NC) Linux 构建完成"

# 构建所有平台
build-all: build-darwin build-windows build-linux
	@echo "$(GREEN)[SUCCESS]$(NC) 所有平台构建完成"
	@echo "$(BLUE)[INFO]$(NC) 构建结果保存在: $(BIN_DIR)/"

# 运行测试
test: deps
	@echo "$(BLUE)[INFO]$(NC) 运行测试..."
	@mkdir -p $(LOG_DIR)
	@cd $(PROJECT_ROOT) && go test ./... -v 2>&1 | tee $(LOG_DIR)/go_test.log
	@if [ -f "$(FRONTEND_DIR)/package.json" ] && grep -q '"test":' "$(FRONTEND_DIR)/package.json"; then \
		cd $(FRONTEND_DIR) && npm test 2>&1 | tee $(LOG_DIR)/frontend_test.log; \
	else \
		echo "$(YELLOW)[WARNING]$(NC) 前端 package.json 中没有找到 test 脚本"; \
	fi
	@echo "$(GREEN)[SUCCESS]$(NC) 测试完成"

# 代码检查
lint: deps
	@echo "$(BLUE)[INFO]$(NC) 运行代码检查..."
	@mkdir -p $(LOG_DIR)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd $(PROJECT_ROOT) && golangci-lint run 2>&1 | tee $(LOG_DIR)/go_lint.log; \
	else \
		echo "$(YELLOW)[WARNING]$(NC) golangci-lint 未安装，跳过 Go 代码检查"; \
		echo "$(BLUE)[INFO]$(NC) 安装命令: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi
	@if [ -f "$(FRONTEND_DIR)/package.json" ] && grep -q '"lint":' "$(FRONTEND_DIR)/package.json"; then \
		cd $(FRONTEND_DIR) && npm run lint 2>&1 | tee $(LOG_DIR)/frontend_lint.log; \
	else \
		echo "$(YELLOW)[WARNING]$(NC) 前端 package.json 中没有找到 lint 脚本"; \
	fi
	@echo "$(GREEN)[SUCCESS]$(NC) 代码检查完成"

# 完整构建流程
all: clean deps lint test build-all
	@echo "$(GREEN)[SUCCESS]$(NC) 完整构建流程完成"
	@echo "$(BLUE)[INFO]$(NC) 构建结果保存在: $(BIN_DIR)/"