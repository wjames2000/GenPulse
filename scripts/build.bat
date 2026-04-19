@echo off
REM GenPulse Windows 编译脚本
REM 支持基本的构建和开发命令

setlocal enabledelayedexpansion

REM 颜色定义
for /f %%a in ('echo prompt $E ^| cmd') do set "ESC=%%a"
set "RED=%ESC%[31m"
set "GREEN=%ESC%[32m"
set "YELLOW=%ESC%[33m"
set "BLUE=%ESC%[34m"
set "NC=%ESC%[0m"

REM 项目根目录
set "PROJECT_ROOT=%~dp0.."
set "FRONTEND_DIR=%PROJECT_ROOT%\frontend"
set "BUILD_DIR=%PROJECT_ROOT%\build"
set "BIN_DIR=%PROJECT_ROOT%\bin"
set "LOG_DIR=%PROJECT_ROOT%\logs"

REM 确保必要的目录存在
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"
if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"
if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"

REM 日志函数
:log_info
    echo %BLUE%[INFO]%NC% %*
    goto :eof

:log_success
    echo %GREEN%[SUCCESS]%NC% %*
    goto :eof

:log_warning
    echo %YELLOW%[WARNING]%NC% %*
    goto :eof

:log_error
    echo %RED%[ERROR]%NC% %*
    goto :eof

REM 检查命令是否存在
:check_command
    where %1 >nul 2>nul
    if errorlevel 1 (
        call :log_error "命令 '%1' 未找到，请先安装"
        exit /b 1
    )
    exit /b 0

REM 检查环境
:check_environment
    call :log_info "检查开发环境..."
    
    REM 检查 Go
    call :check_command go
    if not errorlevel 1 (
        for /f "tokens=3" %%v in ('go version') do set "GO_VERSION=%%v"
        call :log_info "Go 版本: !GO_VERSION!"
    )
    
    REM 检查 Wails
    call :check_command wails
    if not errorlevel 1 (
        for /f "delims=" %%v in ('wails version 2^>nul ^| findstr /v "If Wails"') do set "WAILS_VERSION=%%v"
        if "!WAILS_VERSION!"=="" set "WAILS_VERSION=未知"
        call :log_info "Wails 版本: !WAILS_VERSION!"
    )
    
    REM 检查 Node.js
    call :check_command node
    if not errorlevel 1 (
        for /f %%v in ('node --version') do set "NODE_VERSION=%%v"
        for /f %%v in ('npm --version') do set "NPM_VERSION=%%v"
        call :log_info "Node.js 版本: !NODE_VERSION!"
        call :log_info "npm 版本: !NPM_VERSION!"
    )
    
    call :log_success "环境检查完成"
    goto :eof

REM 清理函数
:clean
    call :log_info "清理构建文件..."
    
    REM 清理前端构建文件
    if exist "%FRONTEND_DIR%\dist" (
        rmdir /s /q "%FRONTEND_DIR%\dist"
        call :log_info "已清理前端 dist 目录"
    )
    
    if exist "%FRONTEND_DIR%\node_modules" (
        rmdir /s /q "%FRONTEND_DIR%\node_modules"
        call :log_info "已清理前端 node_modules"
    )
    
    REM 清理二进制文件
    if exist "%BIN_DIR%" (
        del /q "%BIN_DIR%\*" >nul 2>nul
        call :log_info "已清理 bin 目录"
    )
    
    REM 清理日志文件
    if exist "%LOG_DIR%" (
        del /q "%LOG_DIR%\*.log" >nul 2>nul
        call :log_info "已清理日志文件"
    )
    
    call :log_success "清理完成"
    goto :eof

REM 安装依赖
:install_deps
    call :log_info "安装依赖..."
    
    REM 安装 Go 依赖
    call :log_info "安装 Go 依赖..."
    cd /d "%PROJECT_ROOT%"
    go mod tidy
    go mod download
    
    REM 安装前端依赖
    call :log_info "安装前端依赖..."
    cd /d "%FRONTEND_DIR%"
    call npm install --no-audit --no-fund
    
    call :log_success "依赖安装完成"
    goto :eof

REM 开发模式
:dev
    call :log_info "启动开发模式..."
    call :check_environment
    
    REM 生成 Wails 绑定
    call :log_info "生成 Wails 绑定..."
    cd /d "%PROJECT_ROOT%"
    wails generate module
    
    REM 启动开发服务器
    call :log_info "启动开发服务器..."
    wails dev
    goto :eof

REM 构建函数
:build
    set "PLATFORM=%~1"
    set "ARCH=%~2"
    
    if "!PLATFORM!"=="" set "PLATFORM=windows"
    if "!ARCH!"=="" set "ARCH=amd64"
    
    call :log_info "构建 !PLATFORM!/!ARCH! 版本..."
    call :check_environment
    
    REM 安装依赖
    call :install_deps
    
    REM 构建命令
    cd /d "%PROJECT_ROOT%"
    
    set "BUILD_CMD=wails build -platform !PLATFORM! -arch !ARCH!"
    set "OUTPUT_NAME=GenPulse.exe"
    
    REM 执行构建
    call :log_info "执行构建命令: !BUILD_CMD!"
    !BUILD_CMD!
    
    if not errorlevel 1 (
        call :log_success "构建成功: !OUTPUT_NAME!"
        
        REM 复制到 bin 目录
        set "BUILD_OUTPUT=%PROJECT_ROOT%\build\!OUTPUT_NAME!"
        if exist "!BUILD_OUTPUT!" (
            if not exist "%BIN_DIR%\!PLATFORM!-!ARCH!" mkdir "%BIN_DIR%\!PLATFORM!-!ARCH!"
            copy "!BUILD_OUTPUT!" "%BIN_DIR%\!PLATFORM!-!ARCH!\" >nul
            call :log_info "已复制到: %BIN_DIR%\!PLATFORM!-!ARCH!\"
        )
    ) else (
        call :log_error "构建失败"
        exit /b 1
    )
    goto :eof

REM 显示帮助信息
:show_help
    echo %BLUE%GenPulse Windows 编译脚本%NC%
    echo.
    echo 用法: %~n0 [命令]
    echo.
    echo 命令:
    echo   check         检查开发环境
    echo   clean         清理构建文件
    echo   deps          安装依赖
    echo   dev           启动开发模式
    echo   build [架构]  构建 Windows 版本
    echo                 架构: amd64 ^(默认^), arm64
    echo   help          显示此帮助信息
    echo.
    echo 示例:
    echo   %~n0 check              # 检查环境
    echo   %~n0 dev                # 启动开发模式
    echo   %~n0 build amd64        # 构建 Windows 64位版本
    echo.
    goto :eof

REM 主函数
:main
    set "COMMAND=%~1"
    
    if "!COMMAND!"=="" (
        call :show_help
        goto :eof
    )
    
    if "!COMMAND!"=="check" (
        call :check_environment
    ) else if "!COMMAND!"=="clean" (
        call :clean
    ) else if "!COMMAND!"=="deps" (
        call :install_deps
    ) else if "!COMMAND!"=="dev" (
        call :dev
    ) else if "!COMMAND!"=="build" (
        call :build %2
    ) else if "!COMMAND!"=="help" (
        call :show_help
    ) else (
        call :log_error "未知命令: !COMMAND!"
        call :show_help
        exit /b 1
    )
    goto :eof

REM 执行主函数
call :main %*