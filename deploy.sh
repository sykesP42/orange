#!/bin/bash

# ============================================
# 🍊 Orange 系统一键部署脚本
# ============================================
# 功能：
# 1. 更新后端程序（修复Bing壁纸、删除清理等）
# 2. 更新前端文件（修复列表点击、Bing壁纸显示）
# 3. 更新Nginx配置（图片显示、CSP策略优化）
# 4. 重启所有服务
# ============================================

set -e  # 遇到错误立即退出

echo "╔══════════════════════════════════════╗"
echo "║   🍊 Orange 系统一键部署工具 v1.0    ║"
echo "╚══════════════════════════════════════╝"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目目录
PROJECT_DIR="/www/orange"
BACKUP_DIR="/www/orange/backups/$(date +%Y%m%d_%H%M%S)"

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
}

# 检查是否为root用户
if [ "$EUID" -ne 0 ]; then
    log_error "请使用 root 用户执行此脚本"
    log_info "使用方式: sudo ./deploy.sh"
    exit 1
fi

# 检查项目目录是否存在
if [ ! -d "$PROJECT_DIR" ]; then
    log_error "项目目录不存在: $PROJECT_DIR"
    exit 1
fi

cd $PROJECT_DIR

# ============================================
# 步骤1: 备份当前版本
# ============================================
log_info "步骤 1/6: 备份当前版本..."
mkdir -p "$BACKUP_DIR"

# 备份关键文件
if [ -f "orange-server-linux" ]; then
    cp orange-server-linux "$BACKUP_DIR/" 2>/dev/null || true
fi

if [ -d "dist" ]; then
    cp -r dist "$BACKUP_DIR/dist_backup" 2>/dev/null || true
fi

if [ -f "nginx.conf" ]; then
    cp nginx.conf "$BACKUP_DIR/" 2>/dev/null || true
fi

log_success "备份完成 → $BACKUP_DIR"

# ============================================
# 步骤2: 停止服务（可选）
# ============================================
log_info "步骤 2/6: 准备更新..."

# 如果使用Docker Compose
if [ -f "docker-compose.yml" ]; then
    log_info "检测到 Docker Compose 环境"
    USE_DOCKER=true
else
    USE_DOCKER=false
    log_warning "未检测到 docker-compose.yml，将使用手动部署模式"
fi

# ============================================
# 步骤3: 部署后端程序
# ============================================
log_info "步骤 3/6: 部署后端程序..."

if [ -f "orange-server-linux.new" ] || [ -f "orange-server-linux" ]; then
    if [ -f "orange-server-linux.new" ]; then
        mv orange-server-linux.new orange-server-linux
        log_success "后端程序已更新 (orange-server-linux)"
    else
        log_warning "未找到新的后端程序文件 (orange-server-linux.new)"
        log_info "如果需要更新后端，请先上传 orange-server-linux.new 文件"
    fi
else
    log_warning "后端程序文件不存在，跳过"
fi

# ============================================
# 步骤4: 部署前端文件
# ============================================
log_info "步骤 4/6: 部署前端文件..."

if [ -d "dist.new" ]; then
    # 备份旧的前端文件
    if [ -d "dist.old" ]; then
        rm -rf dist.old
    fi
    mv dist dist.old
    
    # 使用新的前端文件
    mv dist.new dist
    log_success "前端文件已更新 (dist/)"
elif [ -d "dist" ]; then
    log_info "前端文件已存在，检查是否需要更新..."
    
    # 检查Home组件的构建时间
    HOME_JS=$(find dist/assets -name "Home-*.js" -type f 2>/dev/null | head -1)
    if [ -n "$HOME_JS" ]; then
        MOD_TIME=$(stat -c %Y "$HOME_JS" 2>/dev/null || stat -f %m "$HOME_JS")
        CURRENT_TIME=$(date +%s)
        DIFF=$((CURRENT_TIME - MOD_TIME))
        
        if [ $DIFF -lt 300 ]; then  # 5分钟内修改过
            log_success "前端文件是最新状态 (< 5分钟前更新)"
        else
            log_warning "前端文件可能不是最新版本 (上次更新: ${DIFF}秒前)"
        fi
    fi
else
    log_error "前端文件目录不存在！"
    log_info "请确保 dist 目录存在"
fi

# ============================================
# 步骤5: 更新Nginx配置
# ============================================
log_info "步骤 5/6: 更新Nginx配置..."

if [ -f "nginx.conf.new" ]; then
    # 备份旧的nginx配置
    cp nginx.conf nginx.conf.backup 2>/dev/null || true
    
    # 使用新的配置
    mv nginx.conf.new nginx.conf
    log_success "Nginx配置已更新"
else
    if [ -f "nginx.conf" ]; then
        # 检查CSP配置是否包含 https:
        if grep -q "img-src 'self' data: blob:" nginx.conf && ! grep -q "img-src 'self' data: blob: https:" nginx.conf; then
            log_warning "Nginx CSP配置可能需要更新以支持外部图片"
            log_info "建议手动添加 https: 到 img-src 指令"
        else
            log_success "Nginx配置已是最新"
        fi
    fi
fi

# ============================================
# 步骤6: 重启服务
# ============================================
log_info "步骤 6/6: 重启服务..."

if [ "$USE_DOCKER" = true ]; then
    log_info "使用 Docker Compose 重启服务..."
    
    # 重新构建并启动（仅重建有变化的容器）
    docker-compose up -d --build --force-recreate orange-backend orange-frontend
    
    # 等待服务启动
    sleep 5
    
    # 检查容器状态
    if docker-compose ps | grep -q "Up"; then
        log_success "Docker 容器启动成功"
        
        # 显示运行状态
        echo ""
        docker-compose ps
    else
        log_error "Docker 容器启动失败！"
        log_info "查看日志: docker-compose logs --tail=50"
        exit 1
    fi
    
    # 重载Nginx配置
    docker exec orange-frontend nginx -t 2>/dev/null && \
    docker exec orange-frontend nginx -s reload && \
    log_success "Nginx 配置已重载"
    
else
    # 手动部署模式
    log_info "手动部署模式..."
    
    # 查找并重启Go进程
    PID=$(pgrep -f "orange-server" || echo "")
    if [ -n "$PID" ]; then
        kill $PID
        log_success "旧的后端进程已停止 (PID: $PID)"
        sleep 2
    fi
    
    # 启动新的后端进程
    if [ -f "orange-server-linux" ]; then
        nohup ./orange-server-linux > server.log 2>&1 &
        NEW_PID=$!
        log_success "新的后端进程已启动 (PID: $NEW_PID)"
    fi
    
    # 重启Nginx
    if command -v nginx &> /dev/null; then
        nginx -t && nginx -s reload
        log_success "Nginx 已重载"
    elif docker ps | grep -q "nginx"; then
        docker exec $(docker ps -qf "name=nginx") nginx -s reload
        log_success "Docker Nginx 已重载"
    fi
fi

# ============================================
# 验证部署结果
# ============================================
echo ""
echo "╔══════════════════════════════════════╗"
echo "║       🎉 部署完成！验证中...        ║"
echo "╚══════════════════════════════════════╝"
echo ""

sleep 3

# 测试后端API
log_info "测试后端健康检查..."
if curl -sf http://localhost:3003/health > /dev/null 2>&1; then
    log_success "后端API正常 ✓"
else
    log_warning "后端API可能还未就绪，稍后再试"
fi

# 测试前端
log_info "测试前端访问..."
if curl -sf http://localhost:8080 > /dev/null 2>&1 || curl -sf http://123.56.99.154:8080 > /dev/null 2>&1; then
    log_success "前端页面可访问 ✓"
else
    log_warning "前端可能还在启动中"
fi

# 测试Bing壁纸API
log_info "测试Bing壁纸API..."
BING_RESULT=$(curl -sf http://localhost:3003/api/bing-wallpaper 2>/dev/null || echo "{}")
if echo "$BING_RESULT" | grep -q '"code":200'; then
    BING_URL=$(echo "$BING_RESULT" | grep -o '"url":"[^"]*"' | head -1 | cut -d'"' -f4)
    if [ -n "$BING_URL" ]; then
        log_success "Bing壁纸API正常 ✓"
        log_info "壁纸URL: ${BING_URL:0:50}..."
    fi
else
    log_warning "Bing壁纸API返回异常（可能是网络问题）"
fi

echo ""
echo "=========================================="
echo "📋  后续操作清单"
echo "=========================================="
echo ""
echo "1️⃣  清除浏览器缓存："
echo "   按 Ctrl+Shift+R (Windows) 或 Cmd+Shift+R (Mac)"
echo "   或在浏览器控制台执行:"
echo "   localStorage.clear(); location.reload();"
echo ""
echo "2️⃣  测试功能："
echo "   ✅ 登录页背景图（Bing壁纸）"
echo "   ✅ 列表视图点击进入详情页"
echo "   ✅ 图片显示（头像、Logo等）"
echo ""
echo "3️⃣  如有问题，查看日志："
echo "   Docker:  docker-compose logs -f --tail=100"
echo "   后端:  tail -f server.log"
echo ""
echo "4️⃣  回滚方案（如有问题）："
echo "   cd $BACKUP_DIR"
echo "   将备份文件恢复即可"
echo ""
echo "=========================================="
echo "✅ 部署脚本执行完毕！"
echo "=========================================="
