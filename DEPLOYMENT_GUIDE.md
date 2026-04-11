# Orange 博主管理系统 - 阿里云服务器完整部署指南

> 服务器公网IP: `123.56.99.154`
> 后端技术栈: Go (Gin) + SQLite
> 前端技术栈: Vue 3 + Vite + Element Plus
> 反向代理: Nginx

---

## 一、项目架构总览

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   浏览器      │────▶│   Nginx      │────▶│  Go Backend  │
│  (Vue3 SPA)  │◀────│  :80         │◀────│  :3003       │
└──────────────┘     └──────────────┘     └──────┬───────┘
                           │                     │
                     静态文件服务              SQLite DB
                   /www/wwwroot/blogger        ./data.db
```

### 服务器目录结构

```
/www/wwwroot/
├── blogger/                    # 前端静态文件（Vite构建产物）
│   ├── index.html
│   ├── assets/
│   │   ├── index-*.js
│   │   └── index-*.css
│   └── ...
│
└── orange-server/              # Go后端项目
    ├── server                  # Go编译后的二进制文件
    ├── .env                    # 环境变量配置
    ├── data.db                 # SQLite数据库
    ├── uploads/                # 上传文件目录
    │   ├── avatars/
    │   ├── team-logos/
    │   ├── team-bgimages/
    │   └── snapshots/
    └── logs/                   # 日志目录
```

---

## 二、服务器环境准备

### 2.1 阿里云ECS实例信息

| 项目 | 值 |
|------|-----|
| 公网IP | 123.56.99.154 |
| 配置 | 2核2G，40G ESSD |
| 系统 | Alibaba Cloud Linux / Ubuntu 22.04 |
| 面板 | 宝塔Linux面板 |

### 2.2 SSH登录服务器

```bash
ssh root@123.56.99.154
```

### 2.3 安装基础环境

```bash
# 更新系统
apt update && apt upgrade -y

# 安装必要工具
apt install -y curl wget git nginx supervisor

# 安装Go环境（用于编译后端）
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
source /etc/profile
go version

# 安装Node.js（用于构建前端）
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt install -y nodejs
node -v && npm -v
```

---

## 三、Go后端部署

### 3.1 上传后端代码

```bash
# 创建项目目录
mkdir -p /www/wwwroot/orange-server

# 方式1: 从本地上传（在本地执行）
scp -r server-go/* root@123.56.99.154:/www/wwwroot/orange-server/

# 方式2: 通过Git拉取（在服务器执行）
cd /www/wwwroot/orange-server
git clone <your-repo-url> .
```

### 3.2 编译Go后端

```bash
cd /www/wwwroot/orange-server

# 下载依赖
go mod download

# 编译为Linux可执行文件
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

# 赋予执行权限
chmod +x server
```

### 3.3 配置环境变量

```bash
cat > /www/wwwroot/orange-server/.env << 'EOF'
SERVER_PORT=3003
DATABASE_PATH=./data.db
JWT_SECRET=your-production-secret-key-change-this
UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE=10485760
EOF
```

**重要**: 生产环境必须修改 `JWT_SECRET` 为强随机字符串：
```bash
JWT_SECRET=$(openssl rand -hex 32)
sed -i "s/your-production-secret-key-change-this/$JWT_SECRET/" /www/wwwroot/orange-server/.env
```

### 3.4 创建必要目录

```bash
cd /www/wwwroot/orange-server
mkdir -p uploads/avatars uploads/team-logos uploads/team-bgimages uploads/snapshots logs
```

### 3.5 使用Supervisor管理Go进程

Supervisor 比 PM2 更适合管理 Go 二进制文件。

```bash
# 安装Supervisor
apt install -y supervisor

# 创建配置文件
cat > /etc/supervisor/conf.d/orange-server.conf << 'EOF'
[program:orange-server]
directory=/www/wwwroot/orange-server
command=/www/wwwroot/orange-server/server
autostart=true
autorestart=true
startsecs=5
startretries=3
user=root
redirect_stderr=true
stdout_logfile=/www/wwwroot/orange-server/logs/server.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=10
environment=SERVER_PORT="3003",DATABASE_PATH="./data.db",JWT_SECRET="%(ENV_JWT_SECRET)s",UPLOAD_PATH="./uploads"
EOF

# 加载配置并启动
supervisorctl reread
supervisorctl update
supervisorctl start orange-server

# 查看状态
supervisorctl status orange-server
```

### 3.6 验证后端运行

```bash
# 检查进程
ps aux | grep server

# 检查端口监听
ss -tlnp | grep 3003

# 测试健康检查
curl http://127.0.0.1:3003/health

# 测试登录API
curl -X POST http://127.0.0.1:3003/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}'
```

---

## 四、前端部署

### 4.1 本地构建前端

```bash
# 在本地开发机执行
cd /path/to/orange
npm install
npm run build
```

构建产物在 `dist/` 目录下。

### 4.2 上传前端文件

```bash
# 创建目标目录
mkdir -p /www/wwwroot/blogger

# 上传构建产物
scp -r dist/* root@123.56.99.154:/www/wwwroot/blogger/
```

### 4.3 验证前端文件

```bash
ls -la /www/wwwroot/blogger/
# 应看到 index.html 和 assets/ 目录
```

---

## 五、Nginx配置

### 5.1 创建Nginx配置

```bash
cat > /etc/nginx/sites-available/orange << 'NGINXEOF'
server {
    listen 80;
    server_name 123.56.99.154;

    # 前端静态文件
    location / {
        root /www/wwwroot/blogger;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    # 后端API代理
    location /api {
        proxy_pass http://127.0.0.1:3003;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # 超时配置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;

        # 缓冲配置
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
        proxy_busy_buffers_size 8k;
    }

    # 上传文件代理
    location /uploads {
        proxy_pass http://127.0.0.1:3003;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 上传文件大小限制
        client_max_body_size 10m;
    }

    # 健康检查端点（代理到Go后端）
    location /health {
        proxy_pass http://127.0.0.1:3003/health;
        access_log off;
    }

    # Gzip压缩
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_min_length 1024;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/json
        application/javascript
        application/x-javascript
        application/xml
        application/xml+rss
        application/vnd.ms-fontobject
        application/x-font-ttf
        font/opentype
        image/svg+xml;

    # 静态资源缓存策略
    location ~* \.(js|css)$ {
        root /www/wwwroot/blogger;
        expires 1y;
        add_header Cache-Control "public, immutable";
        access_log off;
    }

    location ~* \.(jpg|jpeg|png|gif|ico|svg|webp)$ {
        root /www/wwwroot/blogger;
        expires 6M;
        add_header Cache-Control "public";
        access_log off;
    }

    location ~* \.(woff|woff2|ttf|eot)$ {
        root /www/wwwroot/blogger;
        expires 1y;
        add_header Cache-Control "public, immutable";
        add_header Access-Control-Allow-Origin "*";
        access_log off;
    }

    # 禁止访问隐藏文件
    location ~ /\. {
        deny all;
        access_log off;
        log_not_found off;
    }
}
NGINXEOF
```

### 5.2 启用站点配置

```bash
# 创建符号链接
ln -sf /etc/nginx/sites-available/orange /etc/nginx/sites-enabled/

# 删除默认站点（可选）
rm -f /etc/nginx/sites-enabled/default

# 测试配置
nginx -t

# 重载Nginx
systemctl reload nginx
```

### 5.3 宝塔面板方式配置

如果使用宝塔面板：

1. 点击「网站」→「添加站点」
2. 域名填写 `123.56.99.154`
3. 根目录选择 `/www/wwwroot/blogger`
4. PHP版本选择「纯静态」
5. 点击「设置」→「配置文件」，替换为上面的Nginx配置
6. 保存并重启Nginx

---

## 六、阿里云安全组配置

### 6.1 必须开放的端口

| 端口 | 协议 | 用途 | 授权对象 |
|------|------|------|----------|
| 22 | TCP | SSH | 你的IP或0.0.0.0/0 |
| 80 | TCP | HTTP | 0.0.0.0/0 |
| 443 | TCP | HTTPS | 0.0.0.0/0 |
| 8888 | TCP | 宝塔面板 | 你的IP |

### 6.2 配置步骤

1. 登录阿里云控制台
2. 找到ECS实例 → 安全组
3. 点击「配置规则」→「添加安全组规则」
4. 按上表添加入方向规则

**注意**: 3003端口不需要对外开放，所有请求通过Nginx代理。

---

## 七、HTTPS配置（推荐）

### 7.1 使用Let's Encrypt免费证书

```bash
# 安装certbot
apt install -y certbot python3-certbot-nginx

# 申请证书（需要有域名）
certbot --nginx -d your-domain.com

# 自动续期
certbot renew --dry-run
```

### 7.2 使用阿里云免费证书

1. 在阿里云控制台申请免费SSL证书
2. 下载Nginx格式证书
3. 上传到服务器：

```bash
mkdir -p /etc/nginx/ssl
# 上传证书文件
scp your-domain.pem root@123.56.99.154:/etc/nginx/ssl/
scp your-domain.key root@123.56.99.154:/etc/nginx/ssl/
```

4. 修改Nginx配置添加HTTPS：

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/nginx/ssl/your-domain.pem;
    ssl_certificate_key /etc/nginx/ssl/your-domain.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # ... 其余配置与HTTP版本相同
}

server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$host$request_uri;
}
```

---

## 八、一键部署脚本

### 8.1 完整部署脚本

```bash
#!/bin/bash
set -e

SERVER_IP="123.56.99.154"
PROJECT_DIR="/www/wwwroot/orange-server"
FRONTEND_DIR="/www/wwwroot/blogger"

echo "=========================================="
echo "  Orange 博主管理系统 - 一键部署脚本"
echo "=========================================="

# 1. 创建目录
echo "[1/7] 创建项目目录..."
mkdir -p $PROJECT_DIR/uploads/{avatars,team-logos,team-bgimages,snapshots}
mkdir -p $PROJECT_DIR/logs
mkdir -p $FRONTEND_DIR

# 2. 编译Go后端
echo "[2/7] 编译Go后端..."
cd $PROJECT_DIR
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go
chmod +x server

# 3. 配置环境变量
echo "[3/7] 配置环境变量..."
if [ ! -f .env ]; then
    JWT_SECRET=$(openssl rand -hex 32)
    cat > .env << EOF
SERVER_PORT=3003
DATABASE_PATH=./data.db
JWT_SECRET=$JWT_SECRET
UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE=10485760
EOF
    echo "已生成新的JWT密钥"
fi

# 4. 配置Supervisor
echo "[4/7] 配置Supervisor..."
cat > /etc/supervisor/conf.d/orange-server.conf << 'EOF'
[program:orange-server]
directory=/www/wwwroot/orange-server
command=/www/wwwroot/orange-server/server
autostart=true
autorestart=true
startsecs=5
startretries=3
user=root
redirect_stderr=true
stdout_logfile=/www/wwwroot/orange-server/logs/server.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=10
EOF

supervisorctl reread
supervisorctl update

# 5. 构建前端
echo "[5/7] 构建前端..."
cd /path/to/orange/frontend
npm install
npm run build
cp -r dist/* $FRONTEND_DIR/

# 6. 配置Nginx
echo "[6/7] 配置Nginx..."
cat > /etc/nginx/sites-available/orange << 'NGINXEOF'
server {
    listen 80;
    server_name 123.56.99.154;

    location / {
        root /www/wwwroot/blogger;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://127.0.0.1:3003;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    location /uploads {
        proxy_pass http://127.0.0.1:3003;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        client_max_body_size 10m;
    }

    location /health {
        proxy_pass http://127.0.0.1:3003/health;
        access_log off;
    }

    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/json application/javascript application/xml+rss;

    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf|eot)$ {
        root /www/wwwroot/blogger;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
NGINXEOF

ln -sf /etc/nginx/sites-available/orange /etc/nginx/sites-enabled/
nginx -t && systemctl reload nginx

# 7. 验证部署
echo "[7/7] 验证部署..."
sleep 3
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:3003/health)
if [ "$HTTP_CODE" = "200" ]; then
    echo "后端健康检查: 通过 (HTTP $HTTP_CODE)"
else
    echo "后端健康检查: 失败 (HTTP $HTTP_CODE)"
fi

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1/)
if [ "$HTTP_CODE" = "200" ]; then
    echo "前端页面检查: 通过 (HTTP $HTTP_CODE)"
else
    echo "前端页面检查: 失败 (HTTP $HTTP_CODE)"
fi

echo ""
echo "=========================================="
echo "  部署完成！"
echo "=========================================="
echo "访问地址: http://$SERVER_IP"
echo "后端API: http://$SERVER_IP/api"
echo "健康检查: http://$SERVER_IP/health"
echo ""
echo "管理命令:"
echo "  查看后端状态: supervisorctl status orange-server"
echo "  查看后端日志: tail -f $PROJECT_DIR/logs/server.log"
echo "  重启后端: supervisorctl restart orange-server"
echo "  重启Nginx: systemctl reload nginx"
```

---

## 九、更新部署

### 9.1 更新Go后端

```bash
cd /www/wwwroot/orange-server

# 拉取最新代码
git pull

# 重新编译
CGO_ENABLED=1 go build -o server ./cmd/main.go

# 重启服务
supervisorctl restart orange-server

# 验证
supervisorctl status orange-server
curl http://127.0.0.1:3003/health
```

### 9.2 更新前端

```bash
# 在本地构建
cd /path/to/orange
npm run build

# 上传到服务器
scp -r dist/* root@123.56.99.154:/www/wwwroot/blogger/

# 或在服务器上构建
cd /www/wwwroot/orange-source
git pull
npm install
npm run build
cp -r dist/* /www/wwwroot/blogger/
```

### 9.3 零停机更新策略

```bash
# 1. 编译新版本
CGO_ENABLED=1 go build -o server-new ./cmd/main.go

# 2. 替换二进制文件
mv server server-old
mv server-new server

# 3. 重启服务（Supervisor会自动重启）
supervisorctl restart orange-server

# 4. 验证新版本
sleep 2
curl http://127.0.0.1:3003/health

# 5. 如果有问题，回滚
# mv server-old server
# supervisorctl restart orange-server
```

---

## 十、部署检查清单

- [ ] Go后端编译成功
- [ ] .env配置正确（JWT_SECRET已修改）
- [ ] uploads目录权限正确
- [ ] Supervisor配置完成
- [ ] 后端进程运行中
- [ ] 后端端口3003监听正常
- [ ] 健康检查接口返回200
- [ ] 前端文件已上传
- [ ] Nginx配置完成
- [ ] Nginx配置测试通过
- [ ] Nginx已重载
- [ ] 安全组80端口已开放
- [ ] 前端页面可访问
- [ ] 登录功能正常
- [ ] API请求正常
- [ ] 文件上传功能正常
- [ ] 备份脚本已配置
