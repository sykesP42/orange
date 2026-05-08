@echo off
chcp 65001 >nul 2>&1
title 🍊 Orange 系统一键部署工具

echo.
echo ╔══════════════════════════════════════╗
echo ║   🍊 Orange 系统部署助手 (Windows)    ║
echo ╚══════════════════════════════════════╝
echo.

:: 检查SCP命令
where scp >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未找到 scp 命令，请确保已安装 OpenSSH 客户端
    echo        或使用 Git Bash / WSL 执行此脚本
    pause
    exit /b 1
)

:: 设置服务器信息（请根据实际情况修改）
set SERVER_IP=120.76.48.218
set SERVER_USER=root
set REMOTE_DIR=/www/orange
set LOCAL_DIR=%~dp0..

echo [信息] 项目目录: %LOCAL_DIR%
echo [信息] 目标服务器: %SERVER_IP%
echo.

:: ============================================
:: 步骤1: 上传后端程序
:: ============================================
echo ══════════════════════════════════════
echo   步骤 1/3: 上传后端程序...
echo ══════════════════════════════════════
echo.

if exist "%LOCAL_DIR%\server-go\orange-server-linux" (
    echo [上传] orange-server-linux ...
    scp "%LOCAL_DIR%\server-go\orange-server-linux" %SERVER_USER%@%SERVER_IP%:%REMOTE_DIR%/orange-server-linux.new
    
    if %errorlevel% equ 0 (
        echo [✓] 后端程序上传成功
    ) else (
        echo [✗] 后端程序上传失败！
        pause
        exit /b 1
    )
) else (
    echo [!] 未找到 orange-server-linux 文件
    echo     请先执行: cd server-go ^&^& go build -o orange-server-linux ./cmd
    pause
    exit /b 1
)

echo.

:: ============================================
:: 步骤2: 上传前端文件
:: ============================================
echo ══════════════════════════════════════
echo   步骤 2/3: 上传前端文件...
echo ══════════════════════════════════════
echo.

if exist "%LOCAL_DIR%\dist" (
    echo [上传] 前端 dist 目录（这可能需要几分钟）...
    
    :: 在服务器上创建临时目录
    ssh %SERVER_USER%@%SERVER_IP% "mkdir -p %REMOTE_DIR%/dist.new"
    
    :: 使用 rsync 或 scp 上传整个目录（scp较慢但兼容性好）
    scp -r "%LOCAL_DIR%\dist\*" %SERVER_USER%@%SERVER_IP%:%REMOTE_DIR%/dist.new/
    
    if %errorlevel% equ 0 (
        echo [✓] 前端文件上传成功
    ) else (
        echo [✗] 前端文件上传失败！
        pause
        exit /b 1
    )
) else (
    echo [!] 未找到 dist 目录
    echo     请先执行: npm run build
    pause
    exit /b 1
)

echo.

:: ============================================
:: 步骤3: 上传Nginx配置和部署脚本
:: ============================================
echo ══════════════════════════════════════
echo   步骤 3/3: 上传配置文件...
echo ══════════════════════════════════════
echo.

if exist "%LOCAL_DIR%\nginx.conf" (
    echo [上传] nginx.conf ...
    scp "%LOCAL_DIR%\nginx.conf" %SERVER_USER%@%SERVER_IP%:%REMOTE_DIR%/nginx.conf.new
    echo [✓] Nginx配置上传成功
)

if exist "%LOCAL_DIR%\deploy.sh" (
    echo [上传] deploy.sh 部署脚本 ...
    scp "%LOCAL_DIR%\deploy.sh" %SERVER_USER%@%SERVER_IP%:%REMOTE_DIR%/deploy.sh
    echo [✓] 部署脚本上传成功
)

echo.

:: ============================================
:: 远程执行部署
:: ============================================
echo ══════════════════════════════════════
echo   正在远程服务器上执行部署...
echo ══════════════════════════════════════
echo.

ssh %SERVER_USER%@%SERVER_IP% "cd %REMOTE_DIR% && chmod +x deploy.sh && ./deploy.sh"

echo.
echo ╔══════════════════════════════════════╗
echo ║       🎉 部署完成！                  ║
echo ╚══════════════════════════════════════╝
echo.
echo 请在浏览器中执行以下操作：
echo.
echo   1. 打开登录页: http://123.56.99.154:8080/login
echo   2. 按 Ctrl+Shift+R 强制刷新（清除缓存）
echo   3. 测试以下功能：
echo      ✓ 登录页背景图（Bing壁纸）
echo      ✓ 列表视图点击进入详情页
echo      ✓ 图片显示正常
echo.
pause
