# 🚀 Orange 云端部署完整指南（Docker + 宝塔面板）

本文档提供从本地到云端的**一键部署流程**，支持宝塔面板管理。

---

## 📋 目录

- [部署架构](#部署架构)
- [方式一：Docker 镜像推送部署（推荐）](#方式一docker-镜像推送部署推荐)
- [方式二：Git 代码同步部署](#方式二git-代码同步部署)
- [宝塔面板集成方案](#宝塔面板集成方案)
- [一键部署脚本](#一键部署脚本)
- [常见问题排查](#常见问题排查)

---

## 🏗️ 部署架构

```
本地开发环境 (Windows/Mac/Linux)
    │
    ├── ① 构建镜像: docker build -t orange:latest .
    │
    ├── ② 推送镜像:
    │   ├── 方案A: 推送到 Docker Hub / 阿里云 ACR
    │   └── 方案B: 导出 tar 文件，SCP 上传到服务器
    │
    └── ③ 服务器部署 (Linux 云服务器)
        │
        ├── 拉取/加载镜像
        ├── docker compose up -d
        └── 配置宝塔面板（可选）
            │
            ├── Nginx 反向代理（域名 + SSL）
            ├── 监控告警
            └── 定时备份
```

---

## 方式一：Docker 镜像推送部署（推荐）

### ✅ 适用场景
- 多台服务器需要部署
- 需要 CI/CD 自动化
- 希望版本化管理镜像

### 步骤 1：本地构建并导出镜像

```bash
# 进入项目目录
cd orange

# 1. 确保已构建前端
npm run build

# 2. 构建 Docker 镜像
docker build -t orange:v1.0.0 .

# 3. 导出镜像为 tar 文件（方便传输）
docker save orange:v1.0.0 | gzip > orange-v1.0.0.tar.gz

# 4. 查看文件大小
ls -lh orange-v1.0.0.tar.gz
# 输出: ~50MB (压缩后)
```

### 步骤 2：上传镜像到服务器

#### **方法 A：使用 SCP 直接传输**

```bash
# 在本地执行（替换为你的服务器信息）
scp orange-v1.0.0.tar.gz root@你的服务器IP:/opt/orange/
```

#### **方法 B：使用宝塔面板上传**

1. 登录宝塔面板 → **文件**
2. 进入 `/opt/orange/` 目录
3. 点击 **上传** → 选择 `orange-v1.0.0.tar.gz`
4. 等待上传完成

### 步骤 3：服务器端部署

**SSH 登录服务器后执行：**

```bash
# 1. 创建项目目录
mkdir -p /opt/orange && cd /opt/orange

# 2. 加载 Docker 镜像
gunzip -c orange-v1.0.0.tar.gz | docker load

# 3. 验证镜像加载成功
docker images | grep orange

# 4. 创建 docker-compose.yml（如果不存在）
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  orange-frontend:
    image: orange:v1.0.0
    container_name: orange-frontend
    ports:
      - "80:80"
      - "443:443"
    restart: unless-stopped
    networks:
      - orange-network

networks:
  orange-network:
    driver: bridge
EOF

# 5. 启动容器
docker compose up -d

# 6. 查看运行状态
docker ps | grep orange

# 7. 测试访问
curl http://localhost
# 应该返回 HTML 内容
```

### 步骤 4：验证部署成功

```bash
# 检查容器状态
docker compose ps

# 查看日志
docker compose logs -f orange-foreground

# 测试 HTTP 访问
curl -I http://localhost

# 如果配置了域名
curl -I http://your-domain.com
```

---

## 方式二：Git 代码同步部署

### ✅ 适用场景
- 单台服务器部署
- 希望在服务器上直接构建
- 团队协作开发

### 步骤 1：服务器准备环境

```bash
# 1. 安装 Docker（如果未安装）
curl -fsSL https://get.docker.com | sh
systemctl start docker && systemctl enable docker

# 2. 安装 Docker Compose
curl -SL https://github.com/docker/compose/releases/latest/download/docker-linux-x86_64 \
  -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# 3. 安装 Git
yum install -y git  # CentOS/Ubuntu: apt install git

# 4. 克隆项目
cd /opt
git clone https://github.com/your-org/orange.git
cd orange
```

### 步骤 2：构建并启动

```bash
# 1. 构建前端（可选：如果 dist 已存在可跳过）
npm install && npm run build

# 2. 使用 Docker Compose 构建并启动
docker compose up -d --build

# 3. 查看状态
docker compose ps
```

---

## 🔧 宝塔面板集成方案

### 为什么用宝塔？

| 功能 | 手动配置 | 宝塔面板 |
|------|---------|---------|
| **Nginx 反向代理** | 手动写配置 | 图形界面操作 |
| **SSL 证书** | Let's Encrypt 命令行 | 一键申请 |
| **防火墙** | iptables/firewalld | 可视化规则 |
| **监控** | 自建脚本 | 内置监控 |
| **备份** | cron 脚本 | 一键定时备份 |

### 步骤 1：安装宝塔面板

```bash
# CentOS 系统
yum install -y wget && wget -O install.sh https://download.bt.cn/install/install_6.0.sh && sh install.sh ed8484bec

# Ubuntu/Debian 系统
wget -O install.sh https://download.bt.cn/install/install-ubuntu_6.0.sh && sudo bash install.sh ed8484bec

# 安装完成后会显示：
# ==================================================================
# 外网面板地址: http://YOUR_SERVER_IP:8888/xxxxxx
# 内网面板地址: http://内网IP:8888/xxxxxx
# username: xxxxxx
# password: xxxxxx
# ==================================================================
```

> ⚠️ **安全提示**: 安装后立即修改默认用户名和密码！

### 步骤 2：配置反向代理（关键！）

**场景**: Docker 容器运行在 80 端口，但你想通过 **域名 + 443 (HTTPS)** 访问。

#### **方法 A：宝塔图形化配置（推荐新手）**

1. 登录宝塔面板 → **网站** → **添加站点**
   ```
   域名: orange.example.com
   根目录: /www/wwwroot/orange.example.com
   PHP版本: 纯静态
   ```

2. 点击站点名称 → **反向代理** → **添加反向代理**
   ```
   代理名称: orange-docker
   目标URL: http://127.0.0.1:80
   发送域名: $host
   ```

3. **配置 SSL 证书**
   - 点击站点 → **SSL** → **Let's Encrypt**
   - 选择域名 → 申请证书
   - 开启 **强制 HTTPS**

4. **完成！** 现在可以通过 `https://orange.example.com` 访问

#### **方法 B：修改 Nginx 配置（高级用户）**

1. 宝塔面板 → **网站** → 找到站点 → **配置文件**

2. 替换为以下内容：

```nginx
server {
    listen 80;
    server_name orange.example.com;
    
    # 强制跳转 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name orange.example.com;

    # SSL 证书路径（宝塔自动生成）
    ssl_certificate    /www/server/panel/vhost/cert/orange.example.com/fullchain.pem;
    ssl_certificate_key    /www/server/panel/vhost/cert/orange.example.com/privkey.pem;
    
    # SSL 优化
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;

    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # 反向代理到 Docker 容器
    location / {
        proxy_pass http://127.0.0.1:80;  # Docker 容器端口
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持（如果需要）
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 日志
    access_log /www/wwwlogs/orange.example.com.log;
    error_log  /www/wwwlogs/orange.example.com.error.log;
}
```

3. 保存后重载 Nginx：
   ```bash
   nginx -t && nginx -s reload
   ```

### 步骤 3：宝塔安全配置

#### **1. 修改面板端口（防止扫描）**
- 宝塔面板 → **设置** → **面板端口** → 改为 `88888` 或其他非常见端口

#### **2. 配置防火墙**
- **安全** → **防火墙**
- 允许端口：
  - `80` (HTTP)
  - `443` (HTTPS)
  - `88888` (宝塔面板)
  - `22` (SSH)

#### **3. 绑定域名访问面板**
- **设置** → **面板设置** → **域名访问限制** → 填入 `bt.yourdomain.com`

### 步骤 4：定时备份（重要！）

**宝塔面板 → 计划任务 → 添加任务：**

```bash
# 任务类型: 备份网站
# 执行周期: 每天 02:00
# 备份到: 服务器磁盘 + 云存储（阿里云OSS/腾讯云COS）

# 或者自定义 Shell 脚本：
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
tar -czf /backup/orange_$DATE.tar.gz /opt/orange
# 删除7天前的备份
find /backup -name "orange_*.tar.gz" -mtime +7 -delete
```

---

## 📜 一键部署脚本

我已为你准备好**自动化部署脚本**，适用于生产环境：

### 脚本 1：服务器初始化 (`deploy/init-server.sh`)

```bash
#!/bin/bash
set -e

echo "🚀 Orange 服务器初始化..."

# 1. 更新系统
apt update && apt upgrade -y  # CentOS: yum update -y

# 2. 安装 Docker
if ! command -v docker &> /dev/null; then
    echo "📦 Installing Docker..."
    curl -fsSL https://get.docker.com | sh
    systemctl start docker && systemctl enable docker
fi

# 3. 安装 Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "📦 Installing Docker Compose..."
    curl -SL "https://github.com/docker/compose/releases/latest/download/docker-linux-x86_64" \
        -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
fi

# 4. 创建目录结构
mkdir -p /opt/orange/{data,backups,logs,ssl}

# 5. 设置权限
chown -R $USER:$USER /opt/orange

echo "✅ 服务器初始化完成！"
echo "下一步: 运行 deploy.sh 开始部署"
```

### 脚本 2：一键部署 (`deploy/deploy.sh`)

```bash
#!/bin/bash
set -e

VERSION=${1:-"latest"}
DOMAIN=${2:-""}

echo "🍊 Orange 一键部署 v${VERSION}"

# 加载镜像
if [ -f "orange-${VERSION}.tar.gz" ]; then
    echo "📦 Loading Docker image..."
    gunzip -c orange-${VERSION}.tar.gz | docker load
else
    echo "❌ Error: Image file not found!"
    exit 1
fi

# 启动容器
echo "🚀 Starting container..."
cat > docker-compose.yml << EOF
version: '3.8'
services:
  orange-frontend:
    image: orange:${VERSION}
    container_name: orange-frontend
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./ssl:/etc/nginx/ssl:ro
      - ./logs:/var/log/nginx
    restart: unless-stopped
    networks:
      - orange-network

networks:
  orange-network:
    driver: bridge
EOF

docker compose up -d

# 等待启动
sleep 5

# 验证
if docker ps | grep -q orange-frontend; then
    echo ""
    echo "✅✅✅ 部署成功！✅✅✅"
    echo ""
    echo "访问地址:"
    if [ -n "$DOMAIN" ]; then
        echo "  http://${DOMAIN}"
        echo "  https://${DOMAIN} (需配置SSL)"
    else
        echo "  http://$(curl -s ifconfig.me)"
        echo "  http://localhost"
    fi
    echo ""
    echo "常用命令:"
    echo "  查看日志: docker compose logs -f"
    echo "  停止服务: docker compose down"
    echo "  重启服务: docker compose restart"
else
    echo "❌ 部署失败，请检查日志:"
    docker compose logs
    exit 1
fi
```

### 使用方法：

```bash
# 1. 将脚本上传到服务器
scp deploy/*.sh root@服务器IP:/opt/orange/deploy/

# 2. SSH 登录服务器
ssh root@服务器IP

# 3. 初始化服务器
chmod +x deploy/init-server.sh && ./deploy/init-server.sh

# 4. 上传镜像后执行部署
./deploy/deploy.sh v1.0.0 your-domain.com
```

---

## ❓ 常见问题排查

### Q1: 端口被占用
```bash
# 检查 80 端口占用
lsof -i :80
netstat -tlnp | grep :80

# 解决：停止占用进程或修改 docker-compose.yml 的端口映射
ports:
  - "8080:80"  # 映射到 8080
```

### Q2: 权限问题
```bash
# Docker 权限错误
sudo usermod -aG docker $USER
# 重新登录或执行 newgrp docker
```

### Q3: 内存不足
```bash
# 查看 Docker 占用内存
docker stats

# 清理无用资源
docker system prune -a
```

### Q4: 宝塔无法连接 Docker
```bash
# 确保 Docker 服务运行中
systemctl status docker

# 重启 Docker
systemctl restart docker
```

### Q5: SSL 证书失败
```bash
# 检查 DNS 解析是否正确
nslookup your-domain.com

# 确保 80/443 端口开放（Let's Encrypt 验证）
# 宝塔面板 → 安全 → 放行端口
```

---

## 🎯 最佳实践总结

### ✅ 生产环境 Checklist

- [ ] 使用强密码 + SSH 密钥登录
- [ ] 修改宝塔默认端口
- [ ] 配置 Let's Encrypt SSL 证书
- [ ] 开启防火墙（仅开放必要端口）
- [ ] 配置定时备份（每日）
- [ ] 监控服务器资源（CPU/内存/磁盘）
- [ ] 设置日志轮转（防止单个日志过大）
- [ ] 定期更新系统和 Docker

### 🔄 更新部署流程

```bash
# 1. 本地构建新版本
npm run build
docker build -t orange:v1.1.0 .
docker save orange:v1.1.0 | gzip > orange-v1.1.0.tar.gz

# 2. 上传到服务器
scp orange-v1.1.0.tar.gz root@server:/opt/orange/

# 3. 服务器上更新
ssh root@server
cd /opt/orange
./deploy/deploy.sh v1.1.0 your-domain.com

# 4. 验证新版本
curl -I https://your-domain.com
```

---

## 💡 总结

### **你的思路完全正确！** ✅

**标准流程：**
```
本地构建测试 
  → 导出 Docker 镜像 
  → 上传到云服务器 
  → docker compose up -d 
  → 宝塔配置域名 + SSL 
  → 完成！✨
```

**优势：**
- ✅ **环境一致性**：本地和服务器完全相同
- ✅ **一键部署**：命令即可完成，无需手动配置依赖
- ✅ **易于回滚**：保留旧版本镜像，随时回退
- ✅ **宝塔友好**：可视化管理和监控
- ✅ **扩展简单**：后续加数据库/Redis 只需修改 compose 文件

---

**🎉 现在你已经掌握了从开发到生产的完整部署流程！**

需要我帮你：
1. 🐳 **修复本地 Docker 问题**并实际运行测试？
2. ☁️ **编写更详细的云端部署脚本**？
3. 🔒 **补充安全加固配置**？
4. 📖 **讲解其他技术细节**？

告诉我你的选择！🚀
