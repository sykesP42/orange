# 🏗️ Orange 技术架构设计文档

本文档详细描述 Orange Blogger Manager 的系统架构、技术选型、设计决策和实现细节。

---

## 📋 目录

- [系统概览](#系统概览)
- [架构设计原则](#架构设计原则)
- [整体架构图](#整体架构图)
- [前端架构](#前端架构)
  - [技术栈选型](#技术栈选型)
  - [目录结构](#前端目录结构)
  - [状态管理](#状态管理)
  - [路由设计](#路由设计)
  - [组件体系](#组件体系)
- [后端架构](#后端架构)
  - [Go 服务架构](#go-服务架构)
  - [数据库设计](#数据库设计)
  - [API 设计规范](#api-设计规范)
- [多端架构](#多端架构)
  - [Capacitor 移动端](#capacitor-移动端)
  - [Tauri 桌面端](#tauri-桌面端)
  - [跨平台抽象层](#跨平台抽象层)
- [DevOps 架构](#devops-架构)
  - [容器化方案](#容器化方案)
  - [CI/CD 流水线](#cicd-流水线)
  - [监控告警](#监控告警)
- [性能优化策略](#性能优化策略)
- [安全架构](#安全架构)
- [扩展性设计](#扩展性设计)

---

## 🎯 系统概览

### 项目定位

Orange 是一个**企业级博主管理系统**，采用 **Monorepo + 多端一体化** 架构，支持 Web、移动端（iOS/Android）、桌面端（Windows/Mac/Linux）三大平台。

### 核心特性

- 🌐 **多端统一**: 一套代码库，覆盖所有主流平台
- ⚡ **高性能**: 虚拟列表、懒加载、代码分割
- 🔒 **安全可靠**: JWT 认证、RBAC 权限、数据加密
- 📊 **可观测**: 完整的日志、监控、告警体系
- 🚀 **可扩展**: 微服务就绪、插件化架构

---

## ✨ 架构设计原则

### 1. 单一职责 (SRP)
每个模块、组件、服务只负责一个明确的功能。

### 2. 开放封闭原则 (OCP)
对扩展开放，对修改关闭。通过插件机制和配置实现扩展。

### 3. 依赖倒置 (DIP)
高层模块不依赖低层模块，都依赖于抽象。

### 4. 关注点分离 (SoC)
前端展示、业务逻辑、数据访问严格分层。

### 5. 约定优于配置
提供合理的默认值，减少配置负担。

---

## 🏛️ 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        用户层 (Users)                        │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐                │
│   │   Web    │  │  Mobile  │  │ Desktop  │                │
│   │ Browser  │  │ iOS/And. │  │ Win/Mac/ │                │
│   └────┬─────┘  └────┬─────┘  └────┬─────┘                │
│        │              │             │                      │
├────────┼──────────────┼─────────────┼──────────────────────┤
│        ▼              ▼             ▼                      │
│   ┌─────────────────────────────────────────┐               │
│   │         前端应用层 (Vue 3 + Vite)       │               │
│   │  ┌─────┐ ┌─────┐ ┌─────┐ ┌──────────┐ │               │
│   │  │Views│ │Comp.│ │Store│ │ Services │ │               │
│   │  └─────┘ └─────┘ └─────┘ └──────────┘ │               │
│   └────────────────┬──────────────────────┘               │
│                    │                                       │
├────────────────────┼───────────────────────────────────────┤
│                    ▼                                       │
│   ┌─────────────────────────────────────────┐               │
│   │          API 网关层 (Nginx)              │               │
│   │  • 反向代理  • SSL 终结  • 负载均衡      │               │
│   └────────────────┬────────────────────────┘               │
│                    │                                       │
├────────────────────┼───────────────────────────────────────┤
│                    ▼                                       │
│   ┌─────────────────────────────────────────┐               │
│   │        后端服务层 (Go + Gin)             │               │
│   │  ┌───────┐ ┌───────┐ ┌───────┐         │               │
│   │  │ Auth  │ │Blogger│ │ Team  │ ...     │               │
│   │  └───────┘ └───────┘ └───────┘         │               │
│   └────────────────┬────────────────────────┘               │
│                    │                                       │
├────────────────────┼───────────────────────────────────────┤
│        ┌───────────┼───────────┐                           │
│        ▼           ▼           ▼                           │
│   ┌────────┐ ┌────────┐ ┌────────┐                        │
│   │PostgreSQL│ │ Redis  │ │MinIO/S3│                        │
│   │(主数据库)│ │ (缓存) │ │(文件)  │                        │
│   └────────┘ └────────┘ └────────┘                        │
└─────────────────────────────────────────────────────────────┘
```

---

## 💻 前端架构

### 技术栈选型

| 技术 | 版本 | 选型理由 |
|------|------|---------|
| Vue | 3.x | Composition API、响应式系统优秀 |
| Vite | 5.x | 极速 HMR、原生 ESM、构建速度快 |
| Pinia | 2.x | TypeScript 友好、轻量级状态管理 |
| Vue Router | 4.x | 官方路由、支持懒加载 |
| Element Plus | 2.x | 企业级 UI 组件库、定制性强 |
| ECharts | 5.x | 数据可视化能力强 |
| Vitest | 1.x | 快速单元测试、Vite 生态 |
| Playwright | 1.x | 跨浏览器 E2E 测试 |

### 前端目录结构

```
src/
├── api/                  # API 接口层
│   ├── index.js         # API 统一封装
│   ├── blogger.js       # 博主相关接口
│   ├── auth.js          # 认证接口
│   └── ...
│
├── components/           # 公共组件 (29个)
│   ├── layout/          # 布局组件
│   │   ├── Layout.vue
│   │   └── Header.vue
│   ├── business/        # 业务组件
│   │   ├── BloggerCard.vue
│   │   ├── StatsDashboard.vue
│   │   └── ...
│   └── common/          # 通用组件
│       ├── Button.vue
│       ├── Modal.vue
│       └── ...
│
├── composables/         # 组合式函数 (Hooks)
│   ├── useAuth.js       # 认证逻辑
│   ├── useTheme.js      # 主题切换
│   ├── useWebSocket.js  # WebSocket 连接
│   └── ...
│
├── locales/             # 国际化资源
│   ├── zh.json          # 中文
│   └── en.json          # 英文
│
├── router/              # 路由配置
│   ├── index.js         # 路由定义
│   └── guards.js        # 路由守卫
│
├── services/            # 平台服务
│   ├── mobile.js        # Capacitor 移动端服务
│   └── desktop.js       # Tauri 桌面端服务
│
├── stores/              # Pinia 状态管理
│   ├── user.js          # 用户状态
│   ├── theme.js         # 主题状态
│   └── ...
│
├── styles/              # 全局样式
│   ├── variables.css    # CSS 变量
│   ├── reset.css        # 样式重置
│   └── animations.css   # 动画库
│
├── utils/               # 工具函数
│   ├── platform.js      # 平台检测
│   ├── format.js        # 格式化工具
│   └── validate.js      # 验证函数
│
├── views/               # 页面组件 (20+)
│   ├── Home.vue         # 首页
│   ├── Admin.vue        # 管理后台
│   └── ...
│
├── App.vue              # 根组件
├── main.js              # 入口文件
└── style.css            # 全局样式
```

### 状态管理 (Pinia)

采用 **模块化 Store** 设计：

```javascript
// stores/user.js
export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: null,
    permissions: []
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.userInfo?.role === 'admin'
  },
  
  actions: {
    async login(credentials) { /* ... */ },
    async logout() { /* ... */ },
    async fetchUserInfo() { /* ... */ }
  },
  
  // 持久化配置
  persist: {
    key: 'orange-user',
    storage: localStorage,
    paths: ['token', 'userInfo']
  }
})
```

### 路由设计

```javascript
// router/index.js
const routes = [
  {
    path: '/',
    component: Layout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/Home.vue'),  // 懒加载
        meta: { title: '首页' }
      },
      // ... 其他路由
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { guest: true }
  }
]
```

**路由守卫：**
```javascript
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  
  // 需要认证的路由
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    return next('/login')
  }
  
  // 访客页面（已登录用户跳转）
  if (to.meta.guest && userStore.isLoggedIn) {
    return next('/')
  }
  
  // 权限检查
  if (to.meta.permission && !hasPermission(to.meta.permission)) {
    return next('/403')
  }
  
  next()
})
```

### 组件体系

#### 三层组件架构

1. **基础组件 (Base Components)**
   - Button, Input, Modal, Toast...
   - 高度可复用、无业务逻辑

2. **业务组件 (Business Components)**
   - BloggerCard, StatsDashboard, Calendar...
   - 包含特定业务逻辑

3. **页面组件 (View Components)**
   - Home, Admin, Settings...
   - 组合业务组件、处理路由和数据

#### 组件通信模式

```
父子组件: Props / Emits
跨层级: Provide / Inject
全局状态: Pinia Store
事件总线: mitt (简单场景)
```

---

## 🖥️ 后端架构

### Go 服务架构

```
server-go/
├── cmd/
│   └── main.go              # 应用入口
│
├── internal/
│   ├── config/              # 配置管理
│   │   └── config.go
│   │
│   ├── database/            # 数据库连接
│   │   └── database.go
│   │
│   ├── handlers/            # HTTP 处理器
│   │   ├── auth.go          # 认证处理器
│   │   ├── blogger.go       # 博主 CRUD
│   │   ├── team.go          # 团队管理
│   │   ├── admin.go         # 管理后台
│   │   └── ...
│   │
│   ├── middleware/          # 中间件
│   │   ├── auth.go          # JWT 认证
│   │   ├── cors.go          # CORS 处理
│   │   ├── csrf.go          # CSRF 防护
│   │   └── rate_limit.go    # 速率限制
│   │
│   ├── models/              # 数据模型
│   │   └── all.go
│   │
│   └── pkg/                 # 内部包
│       └── auth/
│           ├── jwt.go       # JWT 工具
│           └── password.go  # 密码加密
│
├── migrations/              # 数据库迁移
│   └── 001_add_features.sql
│
├── go.mod
└── go.sum
```

### 数据库设计

#### 核心表结构

```sql
-- 博主表
CREATE TABLE bloggers (
    id SERIAL PRIMARY KEY,
    nickname VARCHAR(100) NOT NULL,
    avatar TEXT,
    category VARCHAR(50),
    platform VARCHAR(50),
    status VARCHAR(20) DEFAULT 'initial_contact',
    contact_info JSONB,
    tags TEXT[],
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 用户表
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'operator',
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 团队表
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    leader_id INTEGER REFERENCES users(id),
    members INTEGER[],
    created_at TIMESTAMP DEFAULT NOW()
);
```

### API 设计规范

遵循 RESTful 规范：

```
GET    /api/v1/bloggers          # 获取博主列表
POST   /api/v1/bloggers          # 创建博主
GET    /api/v1/bloggers/:id      # 获取博主详情
PUT    /api/v1/bloggers/:id      # 更新博主
DELETE /api/v1/bloggers/:id      # 删除博主

POST   /api/v1/auth/login        # 登录
POST   /api/v1/auth/logout       # 登出
GET    /api/v1/auth/me           # 当前用户信息

GET    /api/v1/stats/dashboard   # 仪表盘数据
GET    /api/v1/stats/analytics   # 分析数据
```

**统一响应格式：**

```json
{
  "code": 200,
  "message": "success",
  "data": { ... },
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "total": 100,
    "totalPages": 5
  }
}
```

---

## 📱 多端架构

### Capacitor 移动端

**架构特点：**

```
Web 应用 (src/)
    ↓ npm run build
dist/ (构建产物)
    ↓ npx cap sync
apps/mobile/android/  ← Android 原生壳
apps/mobile/ios/      ← iOS 原生壳
```

**核心能力集成：**

```javascript
// services/mobile.js
class MobileService {
  // 相机拍照
  async takePhoto() { /* Camera plugin */ }
  
  // 地理位置
  async getCurrentPosition() { /* Geolocation plugin */ }
  
  // 推送通知
  async sendNotification() { /* LocalNotifications */ }
  
  // 分享功能
  async shareContent() { /* Share plugin */ }
  
  // 本地存储
  async savePreference() { /* Preferences plugin */ }
}
```

**原生桥接示例：**

```typescript
// Android: android/app/src/main/java/com/orange/OrangePlugin.java
@CapacitorPlugin(name = "OrangePlugin")
public class OrangePlugin extends Plugin {
    @PluginMethod
    public void exportToExcel(PluginCall call) {
        String data = call.getString("data");
        String path = call.getString("path");
        
        // 使用 Apache POI 库生成 Excel
        boolean success = ExcelExporter.export(data, path);
        
        call.resolve(new JSObject().put("success", success));
    }
}
```

### Tauri 桌面端

**架构特点：**

```
Web 应用 (src/)
    ↓ npm run build
dist/
    ↓ Tauri 打包
apps/desktop/src-tauri/
    ├── src/lib.rs    ← Rust 后端
    ├── Cargo.toml    ← Rust 依赖
    └── icons/        ← 应用图标
```

**Rust 后端能力：**

```rust
// src-tauri/src/lib.rs
#[tauri::command]
async fn export_to_excel(data: Vec<BloggerRecord>, path: String) -> Result<(), String> {
    let mut workbook = Xlsx::new();
    // 使用 calamine 库生成 Excel
    workbook.save(&path).map_err(|e| e.to_string())?;
    Ok(())
}

#[tauri::command]
fn get_system_info() -> SystemInfo {
    SystemInfo {
        os: std::env::consts::OS.to_string(),
        arch: std::env::consts::ARCH.to_string(),
    }
}
```

**系统集成：**

- ✅ 文件系统访问
- ✅ 系统通知
- ✅ 系统托盘图标
- ✅ 全局快捷键
- ✅ 自动更新

### 跨平台抽象层

统一的多端适配接口：

```javascript
// utils/platform.js
export const Platform = {
  isMobile: () => window.Capacitor || window.innerWidth < 768,
  isDesktop: () => window.__TAURI__ || window.innerWidth >= 1024,
  isWeb: () => !Platform.isMobile() && !Platform.isDesktop(),
  
  getType: () => {
    if (window.__TAURI__) return 'desktop'
    if (window.Capacitor) return 'mobile'
    return 'web'
  }
}

// 使用示例
import { mobileService } from '@/services/mobile'
import { desktopService } from '@/services/desktop'

function exportData(data) {
  switch (Platform.getType()) {
    case 'mobile':
      return mobileService.shareContent({ text: JSON.stringify(data) })
    case 'desktop':
      return desktopService.exportToExcel(data)
    case 'web':
      downloadFile(JSON.stringify(data), 'data.json')
  }
}
```

---

## 🚢 DevOps 架构

### 容器化方案

**Docker 镜像分层：**

```dockerfile
# 阶段1: 构建前端
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# 阶段2: 生产镜像
FROM nginx:1.25-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

**Docker Compose 编排：**

```yaml
services:
  frontend:
    build: .
    ports:
      - "80:80"
    depends_on:
      - backend
      
  backend:
    build: ./server-go
    ports:
      - "3003:3003"
    environment:
      - DB_HOST=db
    depends_on:
      db:
        condition: service_healthy
        
  db:
    image: postgres:15
    volumes:
      - postgres_data:/var/lib/postgresql/data
```

### CI/CD 流水线

**GitHub Actions 工作流：**

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: npm ci
      - run: npm run lint
      - run: npm run test:unit
      - run: npm run build

  build-docker:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: docker/build-push-action@v5
        with:
          push: true
          tags: orange:latest

  deploy:
    needs: build-docker
    runs-on: ubuntu-latest
    steps:
      - uses: appleboy/ssh-action@v1
        with:
          script: docker compose up -d
```

### 监控告警

**前端监控：**

```javascript
// utils/monitor.js
export class PerformanceMonitor {
  trackPageLoad() {
    const [paint] = performance.getEntriesByType('paint')
    this.report('FCP', paint.startTime)
  }
  
  trackError(error, context) {
    this.sendToSentry(error, context)
  }
  
  trackAPIPerformance(url, duration) {
    if (duration > 2000) {
      console.warn(`Slow API: ${url} took ${duration}ms`)
    }
  }
}
```

**后端监控：**

```go
// 中间件: 记录请求耗时
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        log.Printf("%s %s %d %v", 
            c.Request.Method, 
            c.Request.URL.Path, 
            c.Writer.Status(),
            duration,
        )
        
        // 上报到监控系统
        metrics.RecordAPILatency(c.Request.URL.Path, duration)
    }
}
```

---

## ⚡ 性能优化策略

### 前端性能

1. **代码分割 (Code Splitting)**
   ```javascript
   // 路由懒加载
   const Admin = () => import('@/views/Admin.vue')
   
   // Vite manualChunks
   manualChunks: {
     'vendor-vue': ['vue', 'vue-router', 'pinia'],
     'vendor-ui': ['element-plus'],
   }
   ```

2. **虚拟列表 (Virtual Scrolling)**
   ```vue
   <VirtualList :items="bloggers" :item-size="120" />
   ```

3. **图片优化**
   ```html
   <!-- 懒加载 -->
   <img v-lazy="avatar" />
   
   <!-- WebP 格式 -->
   <picture>
     <source srcset="avatar.webp" type="image/webp">
     <img src="avatar.jpg" />
   </picture>
   ```

4. **缓存策略**
   ```javascript
   // API 缓存
   class ApiCache {
     async get(key, fetcher) {
       if (this.cache.has(key)) return this.cache.get(key)
       const data = await fetcher()
       this.cache.set(key, data, { ttl: 5 * 60 * 1000 })
       return data
     }
   }
   ```

### 后端性能

1. **数据库索引**
   ```sql
   CREATE INDEX idx_bloggers_status ON bloggers(status);
   CREATE INDEX idx_bloggers_category ON bloggers(category);
   ```

2. **Redis 缓存**
   ```go
   func GetBloggers() ([]Blogger, error) {
       cached, err := redis.Get("bloggers:list").Result()
       if err == nil {
           return json.Unmarshal([]byte(cached))
       }
       
       // 查询数据库
       bloggers := db.FindBloggers()
       
       // 写入缓存
       redis.Set("bloggers:list", json.Marshal(bloggers), 5*time.Minute)
       
       return bloggers
   }
   ```

3. **连接池**
   ```go
   db.SetMaxOpenConns(100)
   db.SetMaxIdleConns(10)
   db.SetConnMaxLifetime(time.Hour)
   ```

---

## 🔒 安全架构

### 认证授权

**JWT Token 流程：**

```
客户端 → POST /auth/login (用户名+密码)
       → 服务器验证密码
       → 生成 JWT Token (包含 userId, role, permissions)
       → 返回 Token
       → 客户端存储 Token (localStorage/cookie)
       
后续请求 → Header: Authorization: Bearer <token>
         → 服务器验证 Token
         → 解析用户信息
         → 执行业务逻辑
```

**权限控制 (RBAC)：**

```javascript
// stores/permission.js
const roles = {
  admin: ['*'],  // 全部权限
  operator: ['bloggers:read', 'bloggers:write'],
  viewer: ['bloggers:read']
}

function hasPermission(permission) {
  const userRole = userStore.role
  return roles[userRole]?.includes(permission) || 
         roles[userRole] === '*'
}
```

### 安全防护措施

| 威胁 | 防护措施 | 实现 |
|------|---------|------|
| XSS | 输出转义、CSP | Vue 自动转义 + Nginx CSP |
| CSRF | Token 验证 | Double Submit Cookie |
| SQL 注入 | 参数化查询 | GORM ORM |
| 点击劫持 | X-Frame-Options | Nginx 配置 |
| 中间人攻击 | HTTPS | SSL/TLS 证书 |
| 暴力破解 | Rate Limiting | Redis 计数器 |

---

## 🔌 扩展性设计

### 插件化架构

```javascript
// plugins/system.js
export class PluginManager {
  plugins = new Map()
  
  register(plugin) {
    this.plugins.set(plugin.name, plugin)
    plugin.install(this.app)
  }
  
  getPlugin(name) {
    return this.plugins.get(name)
  }
}

// 示例插件
const analyticsPlugin = {
  name: 'analytics',
  install(app) {
    app.provide('analytics', new AnalyticsService())
  }
}

pluginManager.register(analyticsPlugin)
```

### 微服务准备

当前单体架构可平滑迁移到微服务：

```
当前: Monolith (Go 单体)
  ↓
阶段1: 模块化 (按领域拆分 Handler)
  ↓
阶段2: 服务拆分 (独立部署)
  ├── auth-service
  ├── blogger-service
  ├── team-service
  └── stats-service
  ↓
阶段3: 引入消息队列 (RabbitMQ/Kafka)
  ↓
阶段4: Service Mesh (Istio)
```

---

## 📈 技术债务管理

### 定期重构计划

- **每月**: 代码审查、依赖更新
- **每季度**: 架构评审、性能优化
- **每年**: 技术栈升级评估

### 代码质量门禁

```yaml
# CI/CD 质量检查
quality_gates:
  coverage_threshold: 80%
  complexity_max: 15
  duplicate_lines_max: 100
  security_scan: required
  performance_test: required
```

---

## 🎯 未来规划

### 短期目标 (Q1-Q2 2024)
- ✅ 多端打包流程完善
- ✅ 生产环境部署
- ✅ 监控告警体系

### 中期目标 (Q3-Q4 2024)
- 🔲 微信小程序版本
- 🔲 AI 智能推荐
- 🔲 数据导出增强

### 长期目标 (2025+)
- 🔲 SaaS 多租户
- 🔲 国际化扩展
- 🔲 开放 API 平台

---

## 📚 相关文档

- [PACKAGING.md](./PACKAGING.md) - 多端打包指南
- [MAINTENANCE.md](./MAINTENANCE.md) - 开发维护手册
- [DEPLOYMENT.md](./DEPLOYMENT.md) - 部署运维指南
- [API.md](./API.md) - API 接口文档

---

**🏗️ 架构设计持续演进中...**

如有问题或建议，欢迎提交 Issue 或 PR！
