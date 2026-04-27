# Orange Blogger Manager

<p align="center">
  <img src="public/pwa-icons/icon.svg" alt="Orange Logo" width="120" height="120">
</p>

<h3 align="center">小而美的博主归属锁定管理工具</h3>

<p align="center">
  <strong>Web · Mobile (iOS/Android) · Desktop (Windows/Mac/Linux)</strong>
</p>

<p align="center">
  <a href="#功能特性">功能特性</a> •
  <a href="#命令面板">命令面板</a> •
  <a href="#技术栈">技术栈</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#部署">部署</a> •
  <a href="#项目结构">项目结构</a>
</p>

---

## 功能特性

### 核心功能

- **博主管理** - 完整的 CRUD 操作，支持批量导入导出，标签分类体系，AI 智能标签推荐
- **智能搜索** - 全文检索、标签筛选、高级过滤、保存筛选条件、词云可视化
- **数据分析** - 可视化统计、平台分布饼图、分类柱状图、月度趋势折线图、团队排行榜
- **日历视图** - 任务安排、跟进提醒、日程管理、合作到期预警
- **团队协作** - 多角色权限、任务分配、团队动态、帖子系统、实时消息
- **工作流自动化** - 自定义规则、状态自动流转、触发器机制
- **数据安全** - 快照备份恢复（自动/手动）、操作审计日志、数据导入导出 Excel
- **失效博主库** - 自动识别失效博主，统一管理归档

### 命令面板 (CommandPalette ⭐)

按 `Ctrl+K` 或点击搜索按钮打开，支持以下搜索方式：

| 方式 | 示例 | 说明 |
|------|------|------|
| **路径代码** | `add` `/my` `/team` | 直接跳转到对应页面 |
| **拼音首字母** | `lr` → 录入, `tj` → 统计, `gl` → 管理 | 中文拼音缩写匹配 |
| **模糊搜索** | `博` → 匹配所有含"博"的结果 | 子串/连续字符匹配 |
| **关键词搜索** | `数据导出` | 名称+描述+关键词多字段 |

**快捷键：**
- `↑↓` 导航选择，`Enter` 执行，`Tab` 切换标签页，`Esc` 关闭
- 支持全局搜索（同时显示命令+博主+页面三类结果）
- 选中项自动居中滚动

### 用户体验

- **主题切换** - 深色/浅色模式自适应，一键切换
- **国际化** - 中英双语无缝切换
- **PWA 支持** - 离线可用、桌面快捷方式安装、Service Worker 缓存
- **响应式设计** - 手机/平板/桌面全适配，卡片/列表视图切换
- **在线状态** - 实时显示用户在线/离线状态
- **私信聊天** - 站内消息通知与对话系统

---

## 技术栈

| 层 | 技术 | 说明 |
|---|------|------|
| **前端框架** | Vue 3 + Composition API + `<script setup>` | 单文件组件 |
| **构建工具** | Vite 5 | 极速热更新，Gzip 压缩 |
| **UI 组件库** | Element Plus | 企业级 UI 组件 |
| **状态管理** | Pinia | 轻量级状态管理 |
| **图表库** | ECharts 5 | 数据可视化大屏 |
| **路由** | Vue Router 4 | 路由守卫 + 权限控制 |
| **后端语言** | Go 1.25+ (Gin) | 高性能 HTTP 框架 |
| **数据库** | SQLite (WAL模式) | 轻量零配置 |
| **缓存** | Redis (可选) | 高频数据缓存 |
| **认证** | JWT + bcrypt | 安全认证 |
| **容器化** | Docker + Docker Compose | 一键部署 |
| **反向代理** | Nginx | SSL / API 反代 / WebSocket |

---

## 快速开始

### 本地开发

```bash
# 克隆项目
git clone https://github.com/your-org/orange.git
cd orange

# 安装前端依赖
npm install

# 启动前端开发服务器（端口 3000）
npm run dev

# 后端服务（另一个终端）
cd server-go
go run ./cmd/main.go
```

访问 http://localhost:3000

### 环境变量

前端 `.env`:
```env
VITE_API_BASE_URL=/api
VITE_WS_URL=ws://localhost:3003/ws
```

后端 `server-go/.env`:
```env
SERVER_PORT=3003
DATABASE_PATH=./data.db
JWT_SECRET=your-secret-key
UPLOAD_PATH=./uploads
CORS_ORIGIN=http://localhost
REDIS_URL=redis:6379
```

---

## 部署

### 一键 Docker 部署（推荐）

```bash
# 1. 构建生产版本
npm run build

# 2. 上传到服务器
scp -r dist root@your-server:/opt/orange/dist_new/

# 3. SSH 登录并部署
ssh root@your-server
docker cp /opt/orange/dist_new/. orange-frontend:/tmp/dist/
docker exec orange-frontend sh -c '
  rm -rf /usr/share/nginx/html/assets/*
  cp -r /tmp/dist/assets/* /usr/share/nginx/html/assets/
  cp /tmp/dist/index.html /usr/share/nginx/html/index.html
'
docker exec orange-frontend nginx -s reload
```

### Docker Compose 编排

```bash
# 启动全部服务
docker-compose up -d --build

# 查看日志
docker-compose logs -f

# 重启服务
docker-compose restart
```

### 架构说明

```
用户浏览器 :80 / :443 (HTTPS)
    ↓
Nginx 反向代理
    ├── 静态资源 → orange-frontend (Nginx)
    ├── /api/*     → orange-backend (:3003)
    ├── /ws/*      → orange-backend (:3003) WebSocket
    └── /uploads/* → orange-backend (:3003)
        ↓
Go Gin + SQLite + Redis (可选)
```

### 文件说明

| 文件 | 用途 |
|------|------|
| [Dockerfile](./Dockerfile) | 前端镜像（Nginx + 预编译 dist） |
| [server-go/Dockerfile](./server-go/Dockerfile) | 后端镜像（Go 编译） |
| [nginx.conf](./nginx.conf) | 前端 Nginx 配置（API 反代 + WebSocket） |
| [docker-compose.yml](./docker-compose.yml) | 容器编排（3个服务） |
| [docker-compose.local.yml](./docker-compose.local.yml) | 本地开发环境 |
| [.dockerignore](./.dockerignore) | 构建排除规则 |

---

## 项目结构

```
orange/
├── src/                          # Vue 3 前端源码
│   ├── api/                      # API 接口层 (axios 封装)
│   │   └── index.js             # 所有 API 接口定义
│   ├── components/               # 公共组件
│   │   ├── CommandPalette.vue   # ⭐ 命令面板 (模糊搜索/拼音/路径代码)
│   │   ├── ConfirmDialog.vue    # 确认对话框
│   │   ├── Toast.vue            # 全局提示
│   │   ├── AvatarCropper.vue   # 头像裁剪上传
│   │   ├── PwaInstall.vue       # PWA 安装提示
│   │   ├── MessageNotify.vue    # 消息通知
│   │   ├── OnlineStatus.vue     # 在线状态指示
│   │   ├── ShortcutHelp.vue     # 快捷键帮助
│   │   └── QuickNote.vue        # 快速笔记
│   ├── views/                    # 页面视图
│   │   ├── Home.vue             # 博主列表首页
│   │   ├── Add.vue              # 录入博主表单
│   │   ├── My.vue               # 个人中心/我的博主
│   │   ├── Team.vue             # 团队协作
│   │   ├── TeamHome.vue         # 团队主页/帖子详情
│   │   ├── CalendarPage.vue     # 日历视图
│   │   ├── Statistics.vue       # 统计分析 (ECharts)
│   │   ├── Admin.vue            # 管理后台 (8个Tab)
│   │   ├── InvalidBloggers.vue  # 失效博主库
│   │   ├── Chat.vue             # 私信聊天
│   │   ├── Login.vue            # 登录注册
│   │   ├── PublicUsers.vue      # 公开用户列表
│   │   ├── UserDetail.vue       # 用户详情
│   │   ├── BloggerDetail.vue    # 博主详情
│   │   ├── UserSettings.vue     # 用户设置
│   │   ├── AnalyticsDashboard.vue # 分析仪表板
│   │   ├── MyBloggers.vue       # 我的博主列表
│   │   └── Workflow.vue         # 工作流管理
│   ├── stores/                   # Pinia 状态管理
│   │   ├── user.js              # 用户认证
│   │   └── app.js               # 应用全局状态
│   ├── router/                   # Vue Router 配置
│   │   └── index.js             # 路由定义 + 守卫
│   ├── utils/                    # 工具函数
│   ├── composables/              # 组合式函数
│   └── assets/                   # 样式资源
├── server-go/                    # Go 后端
│   ├── cmd/main.go              # 入口文件
│   ├── internal/
│   │   ├── handlers/            # API 处理器
│   │   │   ├── auth.go          # 认证 (登录/注册/JWT)
│   │   │   ├── blogger.go       # 博主 CRUD
│   │   │   ├── category.go      # 分类管理
│   │   │   ├── team.go          # 团队管理
│   │   │   ├── user.go          # 用户管理
│   │   │   ├── statistics.go    # 统计数据
│   │   │   ├── snapshot.go      # 快照备份
│   │   │   ├── import.go        # 数据导入
│   │   │   ├── export.go        # 数据导出
│   │   │   ├── notification.go  # 通知
│   │   │   ├── message.go       # 消息/聊天
│   │   │   ├── tag.go           # 标签管理
│   │   │   ├── workflow.go      # 工作流
│   │   │   └── misc.go          # 杂项 (Excel/上传)
│   │   ├── middleware/           # 中间件
│   │   │   ├── auth.go          # JWT 认证
│   │   │   ├── cors.go          # CORS 跨域
│   │   │   ├── ratelimit.go     # 限流
│   │   │   └── logger.go        # 日志
│   │   ├── database/            # 数据库初始化
│   │   └── config/              # 配置加载
│   └── pkg/
│       └── auth/                # JWT + 密码工具
├── public/                       # 静态资源
│   └── pwa-icons/              # PWA 图标 (192x192)
├── scripts/                      # 部署脚本
│   └── deploy.sh                # 一键部署脚本
├── .github/workflows/            # CI/CD
│   └── deploy.yml               # 自动部署
├── Dockerfile                    # 前端镜像
├── docker-compose.yml            # 生产环境编排
├── docker-compose.local.yml      # 本地开发环境
├── nginx.conf                   # Nginx 配置
├── package.json                 # 前端依赖
├── vite.config.js               # Vite 配置
└── README.md                    # 本文件
```

---

## 页面路由

| 路径 | 页面 | 说明 |
|------|------|------|
| `/` | 首页 | 博主列表、搜索筛选 |
| `/add` | 录入博主 | 新增博主表单 |
| `/my` | 个人中心 | 我的资料/博主/团队 |
| `/team` | 团队协作 | 团队列表/管理 |
| `/team/:id` | 团队主页 | 动态/帖子/成员 |
| `/calendar` | 日历视图 | 日程安排 |
| `/statistics` | 统计分析 | ECharts 图表 |
| `/admin` | 管理后台 | 分类/产品/用户/导入/导出/日志 |
| `/invalid-bloggers` | 失效博主 | 失效博主库 |
| `/chat` | 私信聊天 | 消息中心 |
| `/blogger/:id` | 博主详情 | 详细信息/操作记录 |
| `/user/:id` | 用户详情 | 用户资料 |
| `/login` | 登录注册 | 认证页面 |

---

## 开发命令

```bash
# 前端
npm run dev          # 启动开发服务器 :3000
npm run build        # 构建生产版本
npm run preview      # 预览生产版本
npm run lint         # 代码检查

# 后端
cd server-go
go run ./cmd/main.go # 启动后端 :3003
go build -o orange  # 编译二进制
./orange             # 运行

# Docker
docker-compose up -d --build   # 启动全部服务
docker-compose logs -f         # 查看日志
docker-compose down            # 停止服务
```

---

## 安全特性

- ✅ JWT Token 认证 + 过期时间控制
- ✅ bcrypt 密码哈希加密存储
- ✅ CSRF Token 防护
- ✅ XSS 攻击防护（输入过滤 + CSP）
- ✅ SQL 注入防护（参数化查询）
- ✅ Rate Limiting 接口速率限制
- ✅ CORS 白名单精确配置
- ✅ 敏感接口权限校验（Admin only）
- ✅ 操作审计日志完整记录
- ✅ 文件上传类型/大小校验
- ✅ Nginx 安全头配置

---

## 已知问题 & TODO

- [ ] 后端 Redis 缓存优化（当前可选）
- [ ] 单元测试覆盖率提升
- [ ] CI/CD 自动化测试流水线
- [ ] 多租户支持
- [ ] 移动端 App (React Native / Flutter)

---

## 许可证

MIT License

---

<p align="center">
  Made with ❤️ by Orange Team
</p>
