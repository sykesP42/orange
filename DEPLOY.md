# Orange 部署指南

> **最后更新**: 2026-04-15  
> **部署方式**: Docker Compose（推荐）  
> **目标服务器**: Ubuntu + 宝塔面板

---

## 一、部署架构概览

```
┌─────────────────────────────────────────────────────┐
│                    用户浏览器                        │
│                  http://your-domain                 │
└──────────────────────┬──────────────────────────────┘
                       │ :80 / :443
┌──────────────────────▼──────────────────────────────┐
│              宝塔 Nginx                              │
│  ┌──────────────────────────────────────────┐       │
│  │  SSL证书管理 / 域名路由 / 安全规则         │       │
│  │  反向代理 → http://127.0.0.1:8080        │       │
│  └──────────────────────────────────────────┘       │
└──────────────────────┬──────────────────────────────┘
                       │ :8080
┌──────────────────────▼──────────────────────────────┐
│           Docker 前端容器 (orange-frontend)           │
│              nginx:1.27-alpine                       │
│  ┌──────────────────────────────────────────┐       │
│  │  Vue SPA 静态文件 (dist/)                 │       │
│  │  Nginx 反代配置:                          │       │
│  │    /api/*   → orange-backend:3003          │       │
│  │    /ws      → orange-backend:3003 (WS)     │       │
│  │    /uploads/* → orange-backend:3003        │       │
│  └──────────────────────────────────────────┘       │
└──────────────────────┬──────────────────────────────┘
                       │ Docker 内部网络
┌──────────────────────▼──────────────────────────────┐
│           Docker 后端容器 (orange-backend)            │
│              alpine:3.21                             │
│  ┌──────────────────────────────────────────┐       │
│  │  Go Gin HTTP Server (:3003)               │       │
│  │  SQLite 数据库 (WAL模式)                  │       │
│  │  JWT 认证 / CORS / 日志中间件             │       │
│  └──────────────────────────────────────────┘       │
│                                                     │
│  数据卷:                                            │
│    ./data → /app/data        (数据库文件)            │
│    ./uploads → /app/uploads  (上传文件)              │
└─────────────────────────────────────────────────────┘
```

---

## 二、服务器要求

### 最低配置

| 资源 | 最低值 | 推荐值 |
|------|--------|--------|
| CPU | 1 核 | 2 核+ |
| 内存 | 1 GB | 2 GB+ |
| 磁盘 | 10 GB | 20 GB+ |
| 带宽 | 1 Mbps | 5 Mbps+ |
| 系统 | Ubuntu 20.04+ | Ubuntu 22.04 LTS |
| Docker | 24.0+ | 最新版 |

### 用户规模参考

| 日活用户 | 所需配置 | 带宽建议 |
|---------|---------|---------|
| < 50 人 | 1C1G 2M | 够用 |
| 50-200 人 | 2C2G 5M | 推荐 |
| 200-500 人 | 2C4G 10M | 必须 |
| 500+ 人 | 4C8G 20M+ | 升级方案 |

---

## 三、快速部署（5分钟）

### 前置条件

```bash
# 1. 安装 Docker（如果没有）
curl -fsSL https://get.docker.com | bash
systemctl start docker && systemctl enable docker

# 2. 安装 Docker Compose（如果没有）
apt install -y docker-compose

# 3. 配置国内镜像加速（中国大陆服务器必须）
mkdir -p /etc/docker
cat > /etc/docker/daemon.json << 'EOF'
{
  "registry-mirrors": [
    "https://docker.1ms.run",
    "https://docker.xuanyuan.me"
  ]
}
EOF
systemctl daemon-reload && systemctl restart docker
```

### 部署步骤

#### 方式 A：使用预编译部署包（推荐，避免 OOM）

**本地执行：**
```bash
cd orange
npm run build                    # 本地编译前端产物
tar czf orange-deploy.tar.gz \    # 打包（排除垃圾文件）
  --exclude=node_modules \
  --exclude=.git \
  --exclude=dist \
  --exclude="*.md" \
  --exclude=server \
  --exclude=apps \
  .
scp orange-deploy.tar.gz root@你的服务器IP:/www/
```

**服务器执行：**
```bash
mkdir -p /www/orange/data /www/orange/uploads
cd /www/orange && rm -rf * .env 2>/dev/null
tar xzf ../orange-deploy.tar.gz --strip-components=1
echo "JWT_SECRET=OrangeProd$(openssl rand -hex 16)" > .env
docker-compose up -d --build
```

#### 方式 B：直接 Git 拉取构建（需要 2G+ 内存）

```bash
cd /www
git clone <your-repo-url> orange
cd orange
echo "JWT_SECRET=$(openssl rand -hex 32)" > .env
docker-compose up -d --build
```

### 验证部署

```bash
# 检查容器状态
docker-compose ps

# 测试后端 API
curl http://localhost:3003/health
# 预期输出: {"status":"ok"}

# 测试前端
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/
# 预期输出: 200

# 查看日志
docker-compose logs -f orange-backend
docker-compose logs -f orange-frontend
```

---

## 四、宝塔面板集成

### 方案一：纯反代模式（推荐）

如果已有宝塔站点，将 Nginx 改为纯反代：

```nginx
# /www/server/panel/vhost/nginx/your-domain.conf
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 300s;
        proxy_send_timeout 300s;
    }

    access_log /www/wwwlogs/your-domain.log;
    error_log /www/wwwlogs/your-domain.error.log;
}
```

### 方案二：Docker 直接占用 80 端口

不通过宝塔 Nginx，Docker 前端容器直接监听 80：

修改 `docker-compose.yml`：
```yaml
orange-frontend:
  ports:
    - "80:80"    # 直接占用 80 端口
```

### SSL 证书（HTTPS）

使用宝塔申请 Let's Encrypt 免费证书：
1. 宝塔面板 → 网站 → 设置 → SSL
2. 申请免费证书
3. 开启强制 HTTPS

或使用 [deploy/nginx.conf](./deploy/nginx.conf) 中的 SSL 配置。

---

## 五、环境变量说明

### 必须修改

| 变量 | 默认值 | 说明 | 安全级别 |
|------|--------|------|---------|
| `JWT_SECRET` | `change-me-in-production` | JWT 签名密钥 | 🔴 **必须修改** |

### 可选调整

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `SERVER_PORT` | `3003` | 后端监听端口 |
| `DATABASE_PATH` | `/app/data/data.db` | SQLite 数据库路径 |
| `UPLOAD_PATH` | `/app/uploads` | 文件上传目录 |
| `MAX_UPLOAD_SIZE` | `10485760` (10MB) | 最大上传文件大小 |
| `CORS_ORIGIN` | `http://localhost` | 允许的跨域来源 |

### 设置方式

```bash
# 方式 1: .env 文件（推荐）
cat > .env << 'EOF'
JWT_SECRET=your-strong-secret-here
CORS_ORIGIN=https://your-domain.com
EOF

# 方式 2: docker-compose.yml environment 字段
# 方式 3: 运行时传入
docker run -e JWT_SECRET=xxx ...
```

---

## 六、数据持久化

### 数据卷映射

| 宿主机路径 | 容器路径 | 内容 | 备份频率 |
|-----------|---------|------|---------|
| `./data/` | `/app/data/` | SQLite 数据库文件 | 每日 |
| `./uploads/` | `/app/uploads/` | 上传的图片/文件 | 每周 |

### 备份命令

```bash
# 数据库备份
cp data/data.db data/data.db.backup.$(date +%Y%m%d_%H%M%S)

# 上传文件备份
tar czf uploads-backup-$(date +%Y%m%d).tar.gz uploads/

# 一键全备份
BACKUP_DIR="backup-$(date +%Y%m%d)"
mkdir -p $BACKUP_DIR
cp data/*.db* $BACKUP_DIR/
cp -r uploads/* $BACKUP_DIR/ 2>/dev/null
tar czf $BACKUP_DIR.tar.gz $BACKUP_DIR/
rm -rf $BACKUP_DIR
```

### 恢复命令

```bash
# 停止服务
docker-compose down

# 恢复数据
cp backup-YYYYMMDD/*.db* data/
cp -r backup-YYYYMMDD/* uploads/

# 重启
docker-compose up -d
```

---

## 七、更新部署

### 流程

```bash
# 1. 本地重新编译前端
npm run build

# 2. 重新打包
tar czf orange-deploy.tar.gz --exclude=node_modules --exclude=.git ...

# 3. 上传到服务器
scp orange-deploy.tar.gz root@server:/www/

# 4. 服务器上更新（不丢失数据）
cd /www/orange
tar xzf ../orange-deploy.tar.gz --strip-components=1
docker-compose up -d --build

# 5. 清理旧镜像（可选）
docker image prune -f
```

### 回滚

```bash
# 如果新版本有问题，回滚到上一版本
cd /www/orange
git log --oneline -5    # 查看提交历史
git checkout <commit-hash>
docker-compose up -d --build
```

---

## 八、运维常用命令

```bash
# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f                    # 全部日志
docker-compose logs -f orange-backend     # 仅后端
docker-compose logs -f orange-frontend    # 仅前端

# 重启服务
docker-compose restart                   # 重启所有
docker-compose restart orange-backend     # 仅重启后端

# 停止服务
docker-compose down                      # 停止并删除容器
docker-compose down -v                   # 同时删除数据卷（⚠️ 会丢数据）

# 进入容器调试
docker exec -it orange-backend sh         # 进入后端容器
docker exec -it orange-frontend sh        # 进入前端容器

# 查看资源使用
docker stats                              # 实时资源监控

# 清理空间
docker system prune -f                    # 清理无用镜像和缓存
df -h                                    # 查看磁盘使用
free -h                                   # 查看内存使用
```

---

## 九、故障排查

### 常见问题

| 问题 | 原因 | 解决方案 |
|------|------|---------|
| 容器启动失败 | 端口被占用 | `lsof -i :8080` 找到并释放端口 |
| 前端反复重启 | DNS 解析失败 | 检查 nginx.conf 中 resolver 配置 |
| API 返回 502 | 后端未启动 | `docker-compose logs orange-backend` |
| 数据库锁定 | WAL 文件冲突 | 确保只有一个后端进程在运行 |
| 内存不足 OOM | 编译阶段内存溢出 | 使用预编译 dist 方式部署 |
| 镜像拉取失败 | docker.io 不通 | 配置国内镜像加速源 |
| Go 依赖下载超时 | proxy.golang.org 不通 | Dockerfile 中设置 GOPROXY=https://goproxy.cn |

### 性能优化

详见 [docs/STRESS_TEST_REPORT.md](../docs/STRESS_TEST_REPORT.md)

关键优化项：
1. **升级带宽** — 最显著的性能提升方式
2. **开启 Gzip** — 已在 nginx.conf 中配置
3. **CDN 加速静态资源** — 减少服务器带宽压力
4. **Redis 缓存热点 API** — 减少数据库查询

---

## 十、多项目共存

一台服务器可运行多个项目：

```yaml
# 项目A: Orange (当前项目)
orange-frontend:  ports: ["8080:80"]
orange-backend:   ports: ["3003:3003"]

# 项目B: 其他应用
project-b:         ports: ["8081:80"]

# 宝塔 Nginx 根据域名分发
# a.example.com → 127.0.0.1:8080
# b.example.com → 127.0.0.1:8081
```

---

## 十一、安全加固建议

1. **修改默认密码** — 首次登录立即更改 admin 密码
2. **强 JWT_SECRET** — 至少 32 位随机字符串
3. **防火墙** — 只开放 80/443 端口
4. **定期备份** — 自动化每日备份脚本
5. **SSL 证书** — 生产环境务必启用 HTTPS
6. **定期更新** — `docker pull` 更新基础镜像安全补丁
