# 🛠️ Orange 开发维护指南

本文档为开发团队提供完整的开发规范、工作流程、故障排查和最佳实践。

---

## 📋 目录

- [开发环境搭建](#开发环境搭建)
- [代码规范](#代码规范)
- [Git 工作流](#git-工作流)
- [日常开发流程](#日常开发流程)
- [测试指南](#测试指南)
- [性能优化](#性能优化)
- [故障排查](#故障排查)
- [版本发布流程](#版本发布流程)
- [监控与告警](#监控与告警)

---

## 💻 开发环境搭建

### 必需软件

| 软件 | 版本 | 用途 |
|------|------|------|
| Node.js | >= 18.0 LTS | 前端运行时 |
| npm / pnpm | >= 9.0 | 包管理器 |
| Git | >= 2.30 | 版本控制 |
| VS Code | 最新版 | 推荐编辑器 |
| Go | 1.21+ | 后端开发（可选）|
| Docker | 24.0+ | 容器化（可选）|

### VS Code 推荐插件

```json
{
  "recommendations": [
    "vue.volar",
    "dbaeumer.vscode-eslint",
    "esbenp.prettier-vscode",
    "bradlc.vscode-tailwindcss",
    "formulahendry.auto-rename-tag",
    "christian-kohler.path-intellisense",
    "ms-azuretools.vscode-docker"
  ]
}
```

### 环境初始化

```bash
# 1. 克隆仓库
git clone https://github.com/your-org/orange.git
cd orange

# 2. 安装前端依赖
npm install

# 3. 安装移动端依赖（可选）
cd apps/mobile && npm install && cd ../..

# 4. 安装桌面端依赖（可选）
cd apps/desktop && npm install && cd ../..

# 5. 复制环境变量模板
cp .env.example .env

# 6. 启动后端服务（需要 Go 环境）
cd server-go
go mod download
go run cmd/main.go
# 后端运行在 http://localhost:3003

# 7. 启动前端开发服务器
cd ../..
npm run dev
# 前端运行在 http://localhost:3000
```

### 环境变量配置

创建 `.env` 文件：

```env
# ===================
# 应用配置
# ===================
VITE_APP_NAME=Orange Blogger Manager
VITE_APP_VERSION=1.0.0

# ===================
# API 配置
# ===================
VITE_API_BASE_URL=http://localhost:3003
VITE_API_TIMEOUT=30000

# ===================
# WebSocket 配置
# ===================
VITE_WS_URL=ws://localhost:3003
VITE_WS_RECONNECT_INTERVAL=5000

# ===================
# 功能开关
# ===================
VITE_ENABLE_PWA=true
VITE_ENABLE_ANALYTICS=false
VITE_ENABLE_DEBUG_MODE=true

# ===================
# 第三方服务（生产环境）
# ===================
# VITE_SENTRY_DSN=https://xxx@sentry.io/xxx
# VITE_GA_TRACKING_ID=G-XXXXXXXXXX
```

---

## 📝 代码规范

### 命名约定

```javascript
// ✅ 正确示例
const userName = 'John'           // camelCase (变量/函数)
const MAX_RETRY_COUNT = 3          // UPPER_SNAKE_CASE (常量)
class UserService {}               // PascalCase (类/组件)
interface BloggersResponse {}      // PascalCase (接口/类型)

// ❌ 错误示例
const username = 'John'           // 不一致
const maxretrycount = 3            // 难读
class user_service {}              // 非 PascalCase
```

### 组件编写规范

```vue
<template>
  <!-- 使用语义化 HTML -->
  <article class="blogger-card" :class="{ 'is-selected': selected }">
    <header class="card-header">
      <!-- ... -->
    </header>
    
    <section class="card-body">
      <!-- ... -->
    </section>
  </article>
</template>

<script setup>
// 1. 导入顺序: Vue API → 第三方库 → 本地组件 → 工具函数
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { formatDateTime } from '@/utils/format'

// 2. Props 定义（带类型和默认值）
const props = defineProps({
  blogger: {
    type: Object,
    required: true
  },
  selected: {
    type: Boolean,
    default: false
  }
})

// 3. Emits 定义
const emit = defineEmits(['select', 'delete'])

// 4. 组合式函数
const { t } = useI18n()
const router = useRouter()

// 5. 响应式状态
const isLoading = ref(false)
const showDetail = ref(false)

// 6. 计算属性
const displayName = computed(() => {
  return props.blogger.nickname || t('common.unknown')
})

// 7. 方法
async function handleSelect() {
  emit('select', props.blogger.id)
}

function handleDelete() {
  emit('delete', props.blogger.id)
}

// 8. 生命周期钩子
onMounted(() => {
  // 初始化逻辑
})
</script>

<style scoped>
/* BEM 命名规范 */
.blogger-card {
  /* ... */
  
  &.is-selected {
    border-color: var(--primary);
  }
  
  &__header {
    display: flex;
    align-items: center;
  }
  
  &__body {
    padding: 16px;
  }
}
</style>
```

### API 调用规范

```javascript
// api/blogger.js
import request from './index'

export const bloggerAPI = {
  /**
   * 获取博主列表
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.pageSize - 每页数量
   * @returns {Promise<{code: number, data: Array, pagination: Object}>}
   */
  getList(params = {}) {
    return request.get('/bloggers', { params })
  },
  
  /**
   * 创建博主
   * @param {Object} data - 博主数据
   */
  create(data) {
    return request.post('/bloggers', data)
  },
  
  /**
   * 更新博主
   * @param {number} id - 博主ID
   * @param {Object} data - 更新数据
   */
  update(id, data) {
    return request.put(`/bloggers/${id}`, data)
  },
  
  /**
   * 删除博主
   * @param {number} id - 博主ID
   */
  delete(id) {
    return request.delete(`/bloggers/${id}`)
  }
}
```

---

## 🔄 Git 工作流

### 分支策略

```
main          ← 生产环境（受保护）
  │
  develop      ← 开发集成分支
  │
  ├── feature/xxx     ← 新功能分支
  ├── bugfix/xxx      ← Bug 修复分支
  └── hotfix/xxx      ← 紧急修复（直接合并到 main）
```

### Commit Message 规范

采用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type 类型：**

| Type | 描述 | 示例 |
|------|------|------|
| feat | 新功能 | `feat(blogger): add batch export feature` |
| fix | Bug 修复 | `fix(login): resolve token expiration issue` |
| docs | 文档更新 | `docs(readme): update installation guide` |
| style | 代码格式 | `style(home): fix indentation` |
| refactor | 重构 | `refactor(api): simplify error handling` |
| perf | 性能优化 | `perf(virtual-list): improve scroll performance` |
| test | 测试相关 | `test(auth): add login unit tests` |
| chore | 构建/工具 | `chore(deps): update dependencies` |

**示例：**

```bash
feat(blogger): add filter by status dropdown

- Add status filter component to Home page
- Implement backend API endpoint for filtering
- Update i18n translations for new filter options

Closes #123
```

### Pull Request 流程

1. **创建分支**
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/add-blogger-filter
   ```

2. **开发并提交**
   ```bash
   git add .
   git commit -m "feat(blogger): add status filter"
   git push origin feature/add-blogger-filter
   ```

3. **创建 PR**
   - 标题：`[Feature] Add blogger status filter`
   - 模板：
     ```
     ## 变更描述
     简要说明本次变更的内容...
     
     ## 测试清单
     - [ ] 单元测试通过
     - [ ] E2E 测试通过
     - [ ] 手动测试完成
     
     ## 截图（UI 变更）
     （附上截图）
     
     ## 关联 Issue
     Closes #123
     ```

4. **Code Review**
   - 至少 1 人 Review 通过
   - CI 检查全部通过
   - 解决所有 Review 意见

5. **合并**
   - Squash and merge 到 develop
   - 定期从 develop 合并到 main

---

## 🚀 日常开发流程

### 新功能开发 Checklist

```markdown
## 开发前
- [ ] 需求明确，有对应 Issue 或文档
- [ ] 已创建功能分支
- [ ] 了解现有代码结构

## 开发中
- [ ] 编写组件/页面代码
- [ ] 实现业务逻辑
- [ ] 添加国际化文本
- [ ] 编写单元测试
- [ ] 更新 API 接口（如需要）

## 开发后
- [ ] 本地测试通过
- [ ] Lint 检查通过
- [ ] 构建成功
- [ ] 提交代码并创建 PR
- [ ] Code Review 通过
- [ ] 合并到 develop
```

### Bug 修复流程

```markdown
1. 复现 Bug
   - 收集错误信息、截图、复现步骤
   
2. 定位问题
   - 查看控制台日志
   - 检查网络请求
   - 分析代码逻辑
   
3. 编写修复
   - 创建 bugfix 分支
   - 编写修复代码
   - 添加回归测试
   
4. 验证修复
   - 本地验证
   - 提交 PR
   - 部署到测试环境验证
   
5. 关闭 Issue
   - 合并代码
   - 关联并关闭 Issue
```

---

## 🧪 测试指南

### 单元测试

```javascript
// test/components/BloggerCard.test.js
import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import BloggerCard from '@/components/BloggerCard.vue'

describe('BloggerCard.vue', () => {
  const mockBlogger = {
    id: 1,
    nickname: 'Test Blogger',
    category: 'Tech',
    status: 'cooperated'
  }

  it('renders blogger info correctly', () => {
    const wrapper = mount(BloggerCard, {
      props: { blogger: mockBlogger }
    })
    
    expect(wrapper.text()).toContain('Test Blogger')
    expect(wrapper.text()).toContain('Tech')
  })
  
  it('emits select event when clicked', async () => {
    const wrapper = mount(BloggerCard, {
      props: { blogger: mockBlogger }
    })
    
    await wrapper.trigger('click')
    
    expect(wrapper.emitted('select')).toBeTruthy()
    expect(wrapper.emitted('select')[0]).toEqual([1])
  })
  
  it('applies selected class when prop is true', () => {
    const wrapper = mount(BloggerCard, {
      props: { 
        blogger: mockBlogger,
        selected: true 
      }
    })
    
    expect(wrapper.find('.blogger-card').classes()).toContain('is-selected')
  })
})
```

### 运行测试命令

```bash
# 运行所有单元测试
npm run test:unit

# 运行特定文件测试
npx vitest run src/components/BloggerCard.test.js

# 监听模式（文件变化自动重跑）
npx vitest --watch

# 生成覆盖率报告
npm run test:coverage

# 运行 E2E 测试
npm run test:e2e

# 运行特定 E2E 测试
npx playwright test home.spec.js
```

---

## ⚡ 性能优化

### 前端性能检查清单

```markdown
## 构建产物优化
- [ ] 代码分割合理（按路由/组件）
- [ ] Vendor 包分离（vue/element-plus/echarts）
- [ ] 图片资源压缩（WebP 格式）
- [ ] CSS/JS 压缩（Gzip/Brotli）

## 运行时性能
- [ ] 大列表使用虚拟滚动
- [ ] 图片懒加载
- [ ] 路由懒加载
- [ ] API 请求缓存
- [ ] 防抖/节流处理

## 加载体验
- [ ] 骨架屏占位
- [ ] 渐进式加载
- [ ] 预加载关键资源
- [ ] Service Worker 缓存策略
```

### 性能分析工具

```bash
# 1. Lighthouse 审计
npx lighthouse http://localhost:3000 --view

# 2. Webpack Bundle Analyzer
npm run build -- --mode analyze
# 打开 dist/report.html 查看 bundle 构成

# 3. Chrome DevTools Performance
# 打开 DevTools → Performance → Record → 分析

# 4. Vue DevTools
# 安装 Vue DevTools 浏览器插件
```

---

## 🔧 故障排查

### 常见问题及解决方案

#### 问题 1: 页面白屏

**可能原因：**
- JavaScript 执行错误
- 资源加载失败（404）
- 路由配置错误

**排查步骤：**
```bash
# 1. 检查控制台错误
# 打开浏览器 DevTools → Console

# 2. 检查网络请求
# Network 面板 → 查看失败的请求

# 3. 检查构建产物
npm run build
# 检查 dist/index.html 是否正确引用资源

# 4. 检查路由配置
# 确保路由组件正确导入
```

#### 问题 2: API 请求失败

**可能原因：**
- 后端服务未启动
- CORS 配置错误
- Token 过期或无效

**排查步骤：**
```bash
# 1. 检查后端是否运行
curl http://localhost:3003/health

# 2. 检查 CORS 配置
# 后端日志中查看 OPTIONS 请求是否返回 200

# 3. 检查 Token
# Application → Local Storage → 查看 token 是否存在且未过期

# 4. 直接调用 API
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:3003/api/v1/bloggers
```

#### 问题 3: 样式错乱

**可能原因：**
- CSS 未正确加载
- 样式优先级冲突
- 响应式断点问题

**排查步骤：**
```bash
# 1. 检查 CSS 文件是否加载
# Network 面板 → Filters: CSS → 查看是否有 404

# 2. 检查样式作用域
# 确保 <style scoped> 正确使用

# 3. 检查响应式布局
# DevTools → Toggle Device Toolbar → 测试不同屏幕尺寸
```

#### 问题 4: 移动端打包失败

**常见错误：**

```bash
# Gradle 同步失败
cd apps/mobile/android
./gradlew clean
./gradlew --refresh-dependencies

# iOS Pod 安装失败
cd apps/mobile/ios
pod deintegrate
pod install

# Capacitor 同步失败
npx cap sync android --verbose
# 查看详细错误信息
```

---

## 🚢 版本发布流程

### 版本号规则

采用 [Semantic Versioning](https://semver.org/)：

```
MAJOR.MINOR.PATCH

例: 1.2.3
- MAJOR (1): 不兼容的 API 变更
- MINOR (2): 向下兼容的功能新增
- PATCH (3): 向下兼容的问题修正
```

### 发布前 Checklist

```markdown
## 代码质量
- [ ] 所有测试通过
- [ ] Lint 检查无错误
- [ ] 无 console.log/debugger 残留
- [ ] 依赖已更新到最新稳定版

## 文档更新
- [ ] README.md 更新（如有功能变更）
- [ ] CHANGELOG.md 添加版本记录
- [ ] PACKAGING.md 更新（如有打包流程变更）
- [ ] API.md 更新（如有接口变更）

## 多端构建
- [ ] Web 构建成功
- [ ] Android APK/AAB 构建成功
- [ ] iOS IPA 构建成功（如适用）
- [ ] Windows/Mac/Linux 桌面包构建成功

## 测试验证
- [ ] Chrome/Firefox/Safari 浏览器测试
- [ ] Android 真机/模拟器测试
- [ ] iOS 真机/模拟器测试（如适用）
- [ ] 桌面应用测试（Windows/Mac）

## 安全检查
- [ ] 无敏感信息泄露（密钥/密码）
- [ ] 依赖无已知漏洞（npm audit）
- [ ] HTTPS 配置正确
```

### 发布脚本

```bash
#!/bin/bash
# scripts/release.sh

VERSION=$1

if [ -z "$VERSION" ]; then
  echo "Usage: $0 <version>"
  echo "Example: $0 1.2.3"
  exit 1
fi

echo "🚀 Releasing version $VERSION..."

# 1. 更新版本号
npm version $VERSION --no-git-tag-version

# 2. 构建 Web
echo "📦 Building Web..."
npm run build

# 3. 构建 Android
echo "📱 Building Android..."
cd apps/mobile && npm run build:android || true
cd ../..

# 4. 构建桌面端
echo "🖥️ Building Desktop..."
cd apps/desktop && npm run build || true
cd ../..

# 5. 创建 Git Tag
git add .
git commit -m "chore(release): v$VERSION"
git tag -a v$VERSION -m "Release version $VERSION"

# 6. 推送到远程
git push origin main
git push origin v$VERSION

echo "✅ Release $VERSION completed!"
```

使用方法：
```bash
chmod +x scripts/release.sh
./scripts/release.sh 1.2.3
```

---

## 📊 监控与告警

### 前端监控指标

```javascript
// utils/monitor.js
class FrontendMonitor {
  constructor() {
    this.initPerformanceObserver()
    this.initErrorCapture()
  }

  initPerformanceObserver() {
    // Core Web Vitals
    new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        this.reportMetric(entry.name, entry.startTime)
      }
    }).observe({ entryTypes: ['paint'] })

    // Long Tasks
    new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        if (entry.duration > 50) {
          this.reportSlowTask(entry)
        }
      }
    }).observe({ entryTypes: ['longtask'] })
  }

  initErrorCapture() {
    window.addEventListener('error', (event) => {
      this.reportError({
        message: event.message,
        filename: event.filename,
        lineno: event.lineno,
        colno: event.colno
      })
    })

    window.addEventListener('unhandledrejection', (event) => {
      this.reportError({
        message: event.reason?.message || 'Unhandled Promise Rejection',
        stack: event.reason?.stack
      })
    })
  }

  reportMetric(name, value) {
    // 发送到监控系统
    fetch('/api/v1/metrics', {
      method: 'POST',
      body: JSON.stringify({ name, value, timestamp: Date.now() })
    })
  }

  reportError(error) {
    // 发送到 Sentry 或自建系统
    console.error('Frontend Error:', error)
  }
}
```

### 日志规范

```javascript
// utils/logger.js
const LOG_LEVELS = {
  ERROR: 0,
  WARN: 1,
  INFO: 2,
  DEBUG: 3
}

class Logger {
  constructor(level = LOG_LEVELS.INFO) {
    this.level = level
  }

  error(...args) {
    if (this.level >= LOG_LEVELS.ERROR) {
      console.error('[ERROR]', ...args)
      this.sendToServer('error', args)
    }
  }

  warn(...args) {
    if (this.level >= LOG_LEVELS.WARN) {
      console.warn('[WARN]', ...args)
    }
  }

  info(...args) {
    if (this.level >= LOG_LEVELS.INFO) {
      console.info('[INFO]', ...args)
    }
  }

  debug(...args) {
    if (this.level >= LOG_LEVELS.DEBUG) {
      console.debug('[DEBUG]', ...args)
    }
  }

  sendToServer(level, args) {
    // 生产环境发送到日志服务器
    if (import.meta.env.PROD) {
      navigator.sendBeacon('/api/v1/logs', JSON.stringify({
        level,
        message: args.join(' '),
        timestamp: Date.now(),
        url: window.location.href,
        userAgent: navigator.userAgent
      }))
    }
  }
}

export const logger = new Logger(
  import.meta.env.DEV ? LOG_LEVELS.DEBUG : LOG_LEVELS.WARN
)
```

---

## 📞 联系与支持

### 内部资源

- **技术文档**: `/docs/` 目录
- **API 文档**: [API.md](./API.md)
- **架构设计**: [ARCHITECTURE.md](./ARCHITECTURE.md)
- **部署指南**: [DEPLOYMENT.md](./DEPLOYMENT.md)
- **打包流程**: [PACKAGING.md](./PACKAGING.md)

### 外部支持

- **GitHub Issues**: https://github.com/your-org/orange/issues
- **Slack/Discord**: #orange-dev 频道
- **邮件**: dev-team@orange.com

---

## 🎯 最佳实践总结

### ✅ Do（推荐做法）

- ✅ 遵循代码规范和命名约定
- ✅ 编写清晰的 Commit Message
- ✅ 为新功能编写测试
- ✅ 及时更新文档
- ✅ 定期重构和优化代码
- ✅ 关注性能和安全
- ✅ Code Review 要认真对待

### ❌ Don't（避免做法）

- ❌ 不要提交敏感信息（密钥/密码）
- ❌ 不要在主分支直接修改
- ❌ 不要忽略测试和文档
- ❌ 不要写过于复杂的函数（保持单一职责）
- ❌ 不要忽略 Lint 和类型错误
- ❌ 不要阻塞其他开发者（频繁提交大改动）

---

**📘 保持学习，持续改进！**

Happy Coding! 🚀
