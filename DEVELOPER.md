# Orange 开发者文档

> 博主归属锁定管理系统 — 完整开发指南

---

## 目录

- [项目概览](#项目概览)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [前端开发](#前端开发)
- [后端开发](#后端开发)
- [数据库](#数据库)
- [API 接口](#api-接口)
- [安全机制](#安全机制)
- [部署](#部署)
- [开发规范](#开发规范)
- [常见问题](#常见问题)

---

## 项目概览

Orange 是一个面向团队协作场景的博主资源管理平台，用于追踪博主信息、跟进记录、合作管理、团队协作等。

**核心功能模块**：博主管理 / 跟进记录 / 合作管理 / 标签系统 / 团队协作 / 公共论坛 / 私信 / 日历提醒 / 数据分析 / 工作流自动化 / 数据快照 / 导入导出

---

## 技术栈

| 层 | 技术 | 版本 | 说明 |
|---|------|------|------|
| 前端框架 | Vue 3 | ^3.4 | Composition API + `<script setup>` |
| 状态管理 | Pinia | ^2.1 | 5 个 Store |
| UI 组件库 | Element Plus | ^2.13 | 全局注册所有图标 |
| 图表 | ECharts | ^6.0 | 数据可视化 |
| HTTP 客户端 | Axios | ^1.6 | 自动重试 + 401 跳转 |
| 路由 | Vue Router | ^4.2 | 懒加载 + 过渡动画 |
| 构建工具 | Vite | ^5.0 | 手动分包 |
| 后端框架 | Gin | v1.9.1 | Go 1.21 |
| 数据库 | SQLite | - | modernc.org/sqlite（纯 Go，无需 CGO） |
| 认证 | JWT | v5.2.0 | HMAC-SHA256 签名 |
| 密码加密 | bcrypt | - | cost=12 |
| 实时通信 | WebSocket | - | gorilla/websocket |
| PWA | Service Worker | - | 缓存优先策略 |

---

## 快速开始

### 环境要求

- Node.js >= 18
- Go >= 1.21
- 无需 CGO（SQLite 使用纯 Go 实现）

### 1. 克隆项目

```bash
git clone <repo-url>
cd orange
```

### 2. 启动前端

```bash
npm install
npm run dev
```

前端开发服务器运行在 `http://localhost:3000`，API 请求自动代理到 `http://localhost:3003`。

### 3. 启动后端

```bash
cd server-go
go mod download
go run cmd/main.go
```

后端服务运行在 `http://localhost:3003`。首次运行自动创建数据库和默认数据。

### 4. 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | 管理员 |
| user1 | admin123 | 普通用户 |
| user2 | admin123 | 普通用户 |
| user3 | admin123 | 普通用户 |

> ⚠️ 生产环境必须修改默认密码和 JWT 密钥。

### 5. 构建生产版本

```bash
# 前端
npm run build          # 输出到 dist/

# 后端
cd server-go
go build -o server cmd/main.go
```

---

## 项目结构

```
orange/
├── public/                    # 静态资源（PWA 文件）
│   ├── manifest.json          # PWA 清单
│   ├── sw.js                  # Service Worker
│   └── pwa-icons/             # PWA 图标
├── src/                       # 前端源码
│   ├── api/                   # API 层
│   │   └── index.js           # 所有 API 调用（Axios 封装）
│   ├── components/            # 公共组件（20 个）
│   ├── composables/           # 组合式函数（5 个）
│   ├── router/                # 路由配置
│   │   └── index.js           # 路由表 + 守卫
│   ├── stores/                # Pinia Store（5 个）
│   ├── styles/                # 全局样式
│   │   ├── header.css         # Header 毛玻璃效果
│   │   └── transitions.css    # 过渡动画库
│   ├── utils/                 # 工具函数
│   ├── views/                 # 页面视图（19 个）
│   ├── App.vue                # 根组件
│   ├── main.js                # 入口（全局快捷键/PWA注册）
│   └── style.css              # 全局样式 + CSS 变量
├── server-go/                 # Go 后端
│   ├── cmd/
│   │   └── main.go            # 程序入口 + 路由注册
│   ├── internal/
│   │   ├── config/            # 配置加载
│   │   │   └── config.go      # 环境变量 + .env 文件
│   │   ├── database/          # 数据库
│   │   │   └── database.go    # 建表/索引/默认数据
│   │   ├── handlers/          # HTTP 处理器（21 个文件）
│   │   ├── middleware/        # 中间件
│   │   │   ├── auth.go        # JWT 认证
│   │   │   ├── cors.go        # 跨域
│   │   │   ├── csrf.go        # CSRF 防护
│   │   │   └── rate_limit.go  # 限流
│   │   └── models/            # 数据模型
│   │       └── all.go         # 所有结构体定义
│   ├── pkg/
│   │   └── auth/              # 认证工具
│   │       ├── jwt.go         # JWT 生成/验证
│   │       └── password.go    # bcrypt 哈希
│   └── go.mod
├── index.html                 # HTML 入口
├── vite.config.js             # Vite 配置
├── package.json
└── nginx-config.conf          # Nginx 生产配置
```

---

## 前端开发

### 路由

所有路由组件使用懒加载，定义在 `src/router/index.js`。

**添加新页面**：

```javascript
// router/index.js
{
  path: '/my-page',
  name: 'MyPage',
  component: () => import('@/views/MyPage.vue'),
  meta: {
    requiresAuth: true,        // 需要登录
    requiresAdmin: false,      // 需要管理员
    transition: 'slide-up'     // 过渡动画
  }
}
```

**路由守卫**：
- 未登录 → 重定向 `/login`
- `requiresAdmin` → 检查 `role === 'admin'`
- 刷新时从 `localStorage` 恢复用户状态

### 状态管理（Pinia Stores）

| Store | 文件 | 用途 |
|-------|------|------|
| `useUserStore` | `stores/user.js` | 用户信息/登录状态/团队信息 |
| `useThemeStore` | `stores/theme.js` | 深色/浅色模式切换 |
| `useToastStore` | `stores/toast.js` | 全局轻提示消息 |
| `useNotification` | `stores/notification.js` | 通知管理 |
| `useReminderStore` | `stores/reminder.js` | 博主过期提醒 |

**使用示例**：

```javascript
import { useUserStore } from '@/stores/user'
const userStore = useUserStore()

userStore.setUser({ token, username, realName, role, id })
userStore.logout()
```

### 组合式函数（Composables）

| Composable | 文件 | 用途 |
|------------|------|------|
| `useClipboard` | `composables/useClipboard.js` | 剪贴板复制（2s 反馈） |
| `useDateFormat` | `composables/useDateFormat.js` | 日期格式化（相对时间等） |
| `useKeyboardShortcut` | `composables/useKeyboardShortcut.js` | 键盘快捷键绑定 |
| `useOnlineStatus` | `composables/useOnlineStatus.js` | 在线/离线状态检测 |
| `useWebSocket` | `composables/useWebSocket.js` | WebSocket 连接（自动重连+心跳） |

**使用示例**：

```javascript
import { useKeyboardShortcut } from '@/composables'

useKeyboardShortcut({
  'ctrl+alt+s': () => save()
})
```

### API 层

所有 API 调用集中在 `src/api/index.js`，基于 Axios 封装。

**特性**：
- 自动附加 `Authorization: Bearer {token}` 请求头
- 401 自动跳转登录页
- 自动重试（最多 3 次，4xx 不重试）

**添加新 API**：

```javascript
// api/index.js
export const myNewAPI = (data) => api.post('/my-endpoint', data)
```

### CSS 变量与主题

主题切换通过 `useThemeStore` 控制，在 `<html>` 元素上切换 `.dark` 类名。

**核心变量**（定义在 `src/style.css`）：

```css
:root {
  --primary: #f97316;
  --bg-primary: #f5f7fa;
  --bg-card: #ffffff;
  --text-primary: #1f2937;
  --text-muted: #6b7280;
  --border-color: #e5e7eb;
  --header-height: 64px;
  --mobile-nav-height: 68px;
  --sidebar-width: 240px;
}

.dark {
  --bg-primary: #111827;
  --bg-card: #1f2937;
  --text-primary: #f9fafb;
  --text-muted: #9ca3af;
  --border-color: #374151;
}
```

**规范**：所有颜色必须使用 CSS 变量，禁止硬编码。

### 全局快捷键

定义在 `src/main.js`：

| 快捷键 | 功能 |
|--------|------|
| `Ctrl+Alt+K` | 打开命令面板 |
| `Ctrl+Alt+N` | 快捷笔记 |
| `Ctrl+Alt+D` | 切换深色模式 |

### 公共组件

| 组件 | 用途 |
|------|------|
| `BaseCard` | 通用卡片（支持类型/悬停/覆盖层） |
| `EmptyState` | 空状态展示 |
| `SkeletonCard` / `SkeletonList` | 骨架屏加载 |
| `EnhancedSearch` | 增强搜索（历史/热门/建议） |
| `CommandPalette` | 命令面板（Ctrl+Alt+K） |
| `QuickNote` | 快捷笔记 |
| `Toast` | 轻提示 |
| `ConfirmDialog` | 确认对话框 |
| `LazyImage` | 懒加载图片 |
| `AvatarCropper` | 头像裁剪 |
| `PullRefresh` | 下拉刷新 |
| `PwaInstall` | PWA 安装提示 |
| `OnlineStatus` | 在线状态指示 |
| `WebSocketStatus` | WebSocket 连接状态 |
| `MessageNotify` | 消息通知 |
| `ImportPanel` | 导入面板 |
| `TemplateSelector` | 模板选择器 |
| `ShortcutHelp` | 快捷键帮助 |
| `AnalyticsDashboard` | 数据分析仪表盘 |
| `Calendar` | 日历组件 |

---

## 后端开发

### 环境变量

通过 `server-go/internal/config/config.go` 加载，支持 `.env` 文件，环境变量优先。

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `SERVER_PORT` | `3003` | 服务端口 |
| `DATABASE_PATH` | `./data.db` | SQLite 数据库路径 |
| `JWT_SECRET` | `your-secret-key-change-in-production` | JWT 签名密钥（**必须修改**） |
| `UPLOAD_PATH` | `./uploads` | 文件上传目录 |
| `MAX_UPLOAD_SIZE` | `10485760` | 最大上传大小（10MB） |
| `CORS_ORIGIN` | `*` | CORS 允许源 |

**创建 .env 文件**：

```env
SERVER_PORT=3003
DATABASE_PATH=./data.db
JWT_SECRET=your-production-secret-key-here
UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE=10485760
CORS_ORIGIN=https://your-domain.com
```

### 中间件链

请求经过以下中间件处理（按顺序）：

```
请求 → gin.Recovery → CORS → RateLimit → Logger → 路由匹配
                                                      ↓
                                              AuthMiddleware（认证路由）
                                                      ↓
                                              CSRFMiddleware（认证路由）
                                                      ↓
                                              AdminMiddleware（管理员路由）
                                                      ↓
                                              Handler 处理
```

| 中间件 | 文件 | 作用域 | 说明 |
|--------|------|--------|------|
| `CORSMiddleware` | `middleware/cors.go` | 全局 | 跨域处理，通配符时 AllowCredentials=false |
| `RateLimitMiddleware` | `middleware/rate_limit.go` | 全局 | 令牌桶限流，100 token，10/s 补充 |
| `AuthMiddleware` | `middleware/auth.go` | 认证路由 | JWT 验证 + HMAC 算法校验 |
| `CSRFMiddleware` | `middleware/csrf.go` | 认证路由 | Double Submit Cookie 模式 |
| `AdminMiddleware` | `middleware/auth.go` | 管理员路由 | 检查 role=admin |

### 添加新 API 端点

**1. 创建 Handler 方法**（在 `internal/handlers/` 对应文件中）：

```go
func (h *MyHandler) MyEndpoint(c *gin.Context) {
    // 获取认证信息
    userIDVal, _ := c.Get("userID")
    userID, ok := userIDVal.(int)
    if !ok {
        c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
        return
    }
    usernameVal, _ := c.Get("username")
    username, _ := usernameVal.(string)

    // 业务逻辑...
    c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}
```

**2. 注册路由**（在 `cmd/main.go` 中）：

```go
// 公开路由
api.GET("/public/endpoint", myHandler.MyEndpoint)

// 认证路由
authorized.POST("/my-endpoint", myHandler.MyEndpoint)

// 管理员路由
admin.POST("/admin-endpoint", myHandler.MyEndpoint)
```

### Context 中可用的认证信息

在 `AuthMiddleware` 之后，以下值可通过 `c.Get()` 获取：

| Key | 类型 | 说明 |
|-----|------|------|
| `userID` | `int` | 用户 ID |
| `username` | `string` | 用户名 |
| `realName` | `string` | 真实姓名 |
| `role` | `string` | 角色（admin/user） |

### 数据库操作规范

```go
// ✅ 正确：参数化查询
h.DB.Exec("UPDATE blogger SET status = ? WHERE id = ?", status, id)

// ❌ 错误：字符串拼接（SQL 注入风险）
h.DB.Exec("UPDATE blogger SET status = '" + status + "' WHERE id = " + strconv.Itoa(id))

// ✅ 正确：IN 查询使用占位符
placeholders := strings.Repeat("?,", len(ids))
placeholders = placeholders[:len(placeholders)-1]
args := make([]interface{}, len(ids))
for i, id := range ids { args[i] = id }
h.DB.Exec("DELETE FROM table WHERE id IN ("+placeholders+")", args...)

// ✅ 正确：事务操作
tx, err := h.DB.Begin()
if err != nil { return }
defer tx.Rollback()
// ... 操作
tx.Commit()
```

### 文件上传规范

所有上传必须：
1. 检查扩展名白名单
2. 校验文件头魔数（magic bytes）
3. 使用 UUID 生成文件名（不保留原始文件名）
4. 限制文件大小

```go
// 文件头魔数校验（已定义在 upload.go）
var magicSignatures = map[string][]byte{
    ".jpg":  {0xFF, 0xD8, 0xFF},
    ".png":  {0x89, 0x50, 0x4E, 0x47},
    ".gif":  {0x47, 0x49, 0x46},
    ".webp": {0x52, 0x49, 0x46, 0x46},
    ".pdf":  {0x25, 0x50, 0x44, 0x46},
}
```

### 路径安全规范

所有从用户输入构造的文件路径必须：
1. 禁止 `..` 和路径分隔符 `/` `\`
2. 验证绝对路径在预期目录内
3. 强制文件名格式

```go
func sanitizeSnapshotFilename(filename string) (string, error) {
    if strings.Contains(filename, "..") || strings.ContainsAny(filename, "/\\") {
        return "", fmt.Errorf("非法文件名")
    }
    if !strings.HasPrefix(filename, "snapshot_") || !strings.HasSuffix(filename, ".db") {
        return "", fmt.Errorf("非法快照文件名")
    }
    return filename, nil
}
```

---

## 数据库

### 连接配置

```go
// SQLite WAL 模式，连接池
MaxOpenConns: 10    // SQLite 文件锁，不宜过大
MaxIdleConns: 5
BusyTimeout:  5s    // 等待锁释放
CacheSize:    64MB   // 页面缓存
SyncMode:     NORMAL // 平衡性能与安全
```

### 完整表结构

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `users` | 用户 | id, username, password, real_name, role, status, team_id |
| `blogger` | 博主 | id, nickname, category, platform, status, user_name, is_deleted, is_invalid |
| `blogger_tags` | 标签定义 | id, name, color, is_system |
| `blogger_tag_relations` | 博主-标签关联 | blogger_id, tag_id (UNIQUE) |
| `blogger_status_history` | 状态变更历史 | blogger_id, old_status, new_status, operator |
| `blogger_evaluations` | 博主评价 | blogger_id, user_id, rating, comment (UNIQUE) |
| `blogger_transfer_requests` | 转让请求 | blogger_id, from_user_id, to_user_id, status |
| `categories` | 分类 | id, name (UNIQUE), color, icon |
| `platforms` | 平台 | id, name (UNIQUE), icon |
| `products` | 产品 | id, name (UNIQUE) |
| `followup` | 跟进记录 | id, blogger_id, content, type, operator |
| `cooperation` | 合作记录 | id, blogger_id, title, date, product, amount |
| `followup_templates` | 跟进模板 | id, user_id, name, content, is_default |
| `teams` | 团队 | id, name, leader_id, invite_code (UNIQUE), member_count |
| `team_posts` | 团队帖子 | id, team_id, author_id, title, content, like_count |
| `team_post_comments` | 帖子评论 | id, post_id, author_id, content |
| `team_post_likes` | 帖子点赞 | post_id, user_id (UNIQUE) |
| `team_post_collects` | 帖子收藏 | post_id, user_id (UNIQUE) |
| `team_messages` | 团队消息 | id, team_id, sender_id, title, content, type |
| `public_posts` | 公共论坛帖子 | id, author_id, title, content, like_count |
| `public_post_comments` | 论坛评论 | id, post_id, author_id, content |
| `public_post_likes` | 论坛点赞 | post_id, user_id (UNIQUE) |
| `public_post_collects` | 论坛收藏 | post_id, user_id (UNIQUE) |
| `notifications` | 通知 | id, user_id, type, title, content, is_read |
| `private_messages` | 私信 | id, from_user_id, to_user_id, content, is_read |
| `operation_log` | 操作日志 | id, action, target, operator, detail |
| `system_settings` | 系统设置 | key (UNIQUE), value |
| `announcements` | 公告 | id, title, content |
| `snapshots` | 数据快照 | id, name, filename, size |
| `saved_filters` | 保存的筛选器 | id, user_id, name, filters |
| `workflow_rules` | 工作流规则 | id, name, trigger, conditions, actions, enabled |

### 重要字段说明

- `blogger.user_name`：存储的是**用户名字符串**（如 "admin"），不是 userID 整数。查询时必须用字符串匹配。
- `blogger.is_deleted`：软删除标记，0=正常，1=已删除。
- `blogger.is_invalid`：失效标记，0=正常，1=失效。
- `users.status`：用户状态，`active`/`pending`/`disabled`。

### 默认数据

首次运行时自动创建：
- 1 个管理员账户 + 3 个测试用户
- 20 个默认分类（美妆、穿搭、美食等）
- 10 个默认平台（小红书、抖音、微博等）
- 12 个默认产品
- 8 个系统标签
- 5 个跟进模板

---

## API 接口

### 路由分组

```
/api                          ← 公开路由（登录/注册/分类/平台/产品）
/api/                         ← 认证路由（需 JWT + CSRF）
/api/admin/                   ← 管理员路由（需 admin 角色）
```

### 认证

**登录**：
```
POST /api/login
Body: { "username": "admin", "password": "admin123" }
Response: { "code": 200, "data": { "token": "eyJ...", "username": "admin", "role": "admin" } }
```

**注册**：
```
POST /api/register
Body: { "username": "newuser", "password": "123456", "real_name": "张三" }
Response: { "code": 200, "message": "注册成功，等待管理员审核" }
```

> 注册后 status=pending，不发放 token，需管理员审批后才能登录。

### CSRF 防护

认证路由启用 CSRF Double Submit Cookie 模式：

1. GET 请求时服务器设置 `XSRF-TOKEN` Cookie
2. 写请求（POST/PUT/DELETE）需在请求头中携带 `X-XSRF-TOKEN`
3. 服务器比对 Cookie 和 Header 中的 token

**前端适配**（如使用 axios）：

```javascript
// 在 api/index.js 中添加
api.interceptors.request.use(config => {
  const csrfToken = document.cookie
    .split('; ')
    .find(row => row.startsWith('XSRF-TOKEN='))
    ?.split('=')[1]
  if (csrfToken) {
    config.headers['X-XSRF-TOKEN'] = decodeURIComponent(csrfToken)
  }
  return config
})
```

### 响应格式

所有 API 统一返回格式：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {}
}
```

| code | 含义 |
|------|------|
| 200 | 成功 |
| 400 | 参数错误 |
| 401 | 未授权 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 429 | 请求过于频繁 |
| 500 | 服务器错误 |

> 注意：即使业务错误，HTTP 状态码仍为 200，通过 body 中的 code 区分。

### 核心 API 列表

#### 博主管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/blogger/list` | 博主列表（分页/筛选） |
| GET | `/blogger/my` | 我的博主 |
| GET | `/blogger/:id` | 博主详情 |
| POST | `/blogger/add` | 添加博主 |
| PUT | `/blogger/:id` | 更新博主 |
| DELETE | `/blogger/:id` | 删除博主（软删除） |
| POST | `/blogger/batch/status` | 批量更新状态 |
| POST | `/blogger/batch/tags` | 批量更新标签 |
| POST | `/blogger/batch/delete` | 批量删除 |

#### 跟进/合作

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/followup/list` | 跟进列表 |
| POST | `/followup/add` | 添加跟进 |
| PUT | `/followup/:id` | 更新跟进 |
| DELETE | `/followup/:id` | 删除跟进 |
| GET | `/cooperation/list` | 合作列表 |
| POST | `/cooperation/add` | 添加合作 |
| PUT | `/cooperation/:id` | 更新合作 |
| DELETE | `/cooperation/:id` | 删除合作 |

#### 团队

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/teams` | 团队列表 |
| POST | `/teams` | 创建团队 |
| PUT | `/teams/:id` | 更新团队 |
| DELETE | `/teams/:id` | 删除团队 |
| GET | `/team/members` | 团队成员 |
| GET | `/team/:id/posts` | 团队帖子 |
| POST | `/team/:id/messages` | 发送团队消息 |

#### 论坛

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/public/posts` | 公开帖子列表 |
| POST | `/forum/posts` | 创建帖子 |
| POST | `/forum/posts/:id/like` | 点赞 |
| POST | `/forum/posts/:id/collect` | 收藏 |

#### 数据分析

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/analytics/overview` | 总览数据 |
| GET | `/analytics/blogger/:id` | 单个博主分析 |
| GET | `/blogger/stat` | 博主统计 |
| GET | `/blogger/charts` | 图表数据 |

#### 管理后台

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/registration-mode` | 注册模式 |
| POST | `/admin/registration-mode` | 设置注册模式 |
| GET | `/admin/log` | 操作日志 |
| POST | `/admin/purge` | 清理旧数据 |
| GET | `/admin/export` | 导出数据 |

---

## 安全机制

### 已实施的安全措施

| 措施 | 说明 |
|------|------|
| JWT HMAC 校验 | 防止算法混淆攻击，仅允许 HMAC 签名方法 |
| JWT claims 安全断言 | 防止类型断言 panic |
| bcrypt cost=12 | 安全的密码哈希强度 |
| CSRF Double Submit Cookie | 防止跨站请求伪造 |
| 限流中间件 | 令牌桶算法，100 token，10/s 补充 |
| SQL 参数化查询 | 所有查询使用 `?` 占位符 |
| 文件头魔数校验 | 防止文件类型伪造 |
| 路径穿越防护 | 禁止 `..`、验证绝对路径 |
| 密码字段排除 | GetProfile/GetUserDetail 不返回密码 |
| 待审核用户无 token | 注册后需管理员审批才能登录 |
| WebSocket 认证 | WebSocket 路由在认证组内 |

### 安全注意事项

1. **JWT_SECRET**：生产环境必须设置强随机密钥
2. **默认密码**：所有测试用户密码为 `admin123`，上线前必须修改
3. **CORS_ORIGIN**：生产环境应设为具体域名，不要使用 `*`
4. **HTTPS**：生产环境必须启用 HTTPS
5. **CSRF 前端适配**：前端需在请求中携带 `X-XSRF-TOKEN` 头

---

## 部署

### Nginx 配置

参考 `nginx-config.conf`：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    client_max_body_size 10m;

    # 前端静态文件
    location / {
        root /www/wwwroot/blogger;
        try_files $uri $uri/ /index.html;
    }

    # API 反向代理
    location /api {
        proxy_pass http://127.0.0.1:3003;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 上传文件代理
    location /uploads {
        proxy_pass http://127.0.0.1:3003;
    }
}
```

### 进程管理

推荐使用 Supervisor 或 systemd 管理 Go 后端进程：

```ini
# /etc/supervisor/conf.d/orange.conf
[program:orange-server]
command=/www/wwwroot/orange-server/server
directory=/www/wwwroot/orange-server
autostart=true
autorestart=true
stderr_logfile=/var/log/orange.err.log
stdout_logfile=/var/log/orange.out.log
environment=JWT_SECRET="your-production-secret",CORS_ORIGIN="https://your-domain.com"
```

### 部署步骤

1. 前端构建：`npm run build`，将 `dist/` 内容上传到服务器
2. 后端构建：`cd server-go && go build -o server cmd/main.go`
3. 上传二进制文件到服务器
4. 配置 Nginx 反向代理
5. 配置 Supervisor/systemd 进程管理
6. 设置环境变量（JWT_SECRET、CORS_ORIGIN 等）
7. 启动服务

---

## 开发规范

### 前端规范

1. **组件风格**：使用 `<script setup>` + Composition API
2. **样式**：使用 CSS 变量，禁止硬编码颜色；支持 `.dark` 深色模式
3. **API 调用**：集中在 `src/api/index.js`，使用封装的 axios 实例
4. **状态管理**：全局状态用 Pinia Store，组件内状态用 `ref`/`reactive`
5. **加载状态**：使用 `SkeletonCard`/`SkeletonList` 骨架屏
6. **空状态**：使用 `EmptyState` 组件
7. **响应式**：所有页面必须适配移动端（断点：1200px/1024px/768px/480px）

### 后端规范

1. **代码风格**：`go fmt` + `go vet`
2. **错误处理**：统一返回 `gin.H{"code": xxx, "message": "..."}`，不使用 panic
3. **SQL 查询**：必须参数化，禁止字符串拼接
4. **文件操作**：必须校验路径安全
5. **权限检查**：每个 Handler 开头验证用户身份和权限
6. **日志记录**：关键操作使用 `database.AddLog()` 记录

### Git 规范

- `feat:` 新功能
- `fix:` 修复 Bug
- `refactor:` 重构
- `docs:` 文档
- `style:` 格式调整
- `chore:` 构建/工具变更

---

## 常见问题

### Q: 前端启动后 API 请求 404？

确保 Go 后端已启动在 3003 端口。Vite 开发服务器会将 `/api` 请求代理到 `http://localhost:3003`。

### Q: 数据库 "database is locked" 错误？

SQLite 文件级锁，MaxOpenConns 不宜过大（当前设为 10）。如果并发写入频繁，考虑减少连接数或增加 BusyTimeout。

### Q: 暗色模式不生效？

确保使用 `.dark` 类名选择器（不是 `[data-theme="dark"]` 或 `html.dark`）。主题切换由 `useThemeStore` 控制，在 `<html>` 上添加/移除 `.dark` 类。

### Q: JWT 认证失败？

1. 检查 `JWT_SECRET` 环境变量是否与签发时一致
2. 确认 token 未过期
3. 检查请求头 `Authorization: Bearer <token>` 格式

### Q: CSRF 请求被拒绝？

前端需要在写请求中携带 `X-XSRF-TOKEN` 请求头，值从 `XSRF-TOKEN` Cookie 中读取。参见 [CSRF 防护](#csrf-防护) 章节。

### Q: 如何重置数据库？

删除 `data.db` 文件，重启后端即可自动重建。

### Q: 如何添加新的数据库表？

在 `server-go/internal/database/database.go` 的 `createTables()` 函数中添加建表 SQL，重启后端自动执行。

### Q: blogger.user_name 是整数还是字符串？

**字符串**。`blogger.user_name` 存储的是用户名字符串（如 "admin"），不是用户 ID 整数。查询时必须使用字符串匹配。

---

**最后更新**: 2026-04-11
