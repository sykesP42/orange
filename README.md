# Orange - 博主归属锁定管理系统

> 🍊 一个现代化的博主管理与协作平台

[![Vue 3](https://img.shields.io/badge/Vue-3.4.0-42b883?logo=vue.js)](https://vuejs.org/)
[![Vite](https://img.shields.io/badge/Vite-5.0.0-646cff?logo=vite)](https://vitejs.dev/)
[![Element Plus](https://img.shields.io/badge/Element_Plus-2.13.5-409EFF?logo=element)](https://element-plus.org/)

---

## 📖 项目简介

Orange 是一个功能强大的博主管理与协作系统，提供博主信息录入、团队管理、公共论坛、数据统计等功能。采用现代化的技术栈，拥有优秀的用户体验和精美的界面设计。

### ✨ 主要特性

- 🎨 **现代化 UI 设计** - 毛玻璃效果、流畅动画、深色模式支持
- 📱 **移动端优先** - 完美适配各种屏幕尺寸，触摸手势支持
- 🔍 **智能搜索** - 搜索历史、热门推荐、实时建议
- 🎭 **页面动画** - 多种过渡效果，可自定义配置
- 🧩 **组件化开发** - 丰富的通用组件库，提高开发效率
- 🔐 **权限管理** - 基于角色的访问控制（用户/管理员）
- 📊 **数据统计** - 可视化数据展示，支持导出 Excel
- 💬 **公共论坛** - 团队内部交流平台，支持话题讨论
- 🔔 **实时通知** - 系统通知、重要提醒

---

## 🚀 快速开始

### 环境要求

- Node.js >= 16.0.0
- npm >= 8.0.0

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd orange
```

2. **安装依赖**
```bash
npm install
```

3. **启动开发服务器**
```bash
npm run dev
```

访问 http://localhost:5173 查看应用

4. **构建生产版本**
```bash
npm run build
```

5. **预览生产构建**
```bash
npm run preview
```

6. **启动后端服务器**（可选）
```bash
npm run server
```

---

## 📁 项目结构

```
orange/
├── src/
│   ├── api/              # API 接口
│   │   └── index.js
│   ├── components/       # 通用组件
│   │   ├── AvatarCropper.vue    # 头像裁剪
│   │   ├── BaseCard.vue         # 基础卡片 ⭐
│   │   ├── ConfirmDialog.vue    # 确认对话框
│   │   ├── EmptyState.vue       # 空状态 ⭐
│   │   ├── EnhancedSearch.vue   # 增强搜索 ⭐
│   │   ├── LazyImage.vue        # 懒加载图片
│   │   ├── MessageNotify.vue    # 消息通知
│   │   ├── PullRefresh.vue      # 下拉刷新
│   │   ├── SkeletonCard.vue     # 卡片骨架屏 ⭐
│   │   └── SkeletonList.vue     # 列表骨架屏 ⭐
│   ├── router/           # 路由配置
│   │   └── index.js
│   ├── stores/           # Pinia 状态管理
│   │   ├── notification.js
│   │   ├── theme.js
│   │   └── user.js
│   ├── styles/           # 全局样式
│   │   ├── header.css           # Header 样式 ⭐
│   │   └── transitions.css      # 动画库 ⭐
│   ├── utils/            # 工具函数
│   │   ├── confirm.js           # 确认对话框
│   │   ├── gestures.js          # 手势工具 ⭐
│   │   ├── headerEffect.js      # Header 效果 ⭐
│   │   └── search.js            # 搜索工具 ⭐
│   ├── views/            # 页面组件
│   │   ├── Add.vue
│   │   ├── Admin.vue
│   │   ├── BloggerDetail.vue
│   │   ├── Forums.vue
│   │   ├── Home.vue
│   │   ├── Layout.vue
│   │   ├── Log.vue
│   │   ├── Login.vue
│   │   ├── My.vue
│   │   ├── MyBloggers.vue
│   │   ├── PublicUsers.vue
│   │   ├── Statistics.vue
│   │   ├── Team.vue
│   │   ├── TeamHome.vue
│   │   ├── UserDetail.vue
│   │   └── UserSettings.vue
│   ├── App.vue           # 根组件
│   ├── main.js           # 入口文件
│   └── style.css         # 基础样式
├── server/               # 后端服务器
├── DEVELOPMENT.md        # 开发文档 📚
├── package.json
└── vite.config.js
```

> ⭐ 标记的为核心自研组件和工具

---

## 🛠️ 技术栈

### 前端框架
- **Vue 3.4** - 渐进式 JavaScript 框架
- **Vite 5.0** - 下一代前端构建工具
- **Vue Router 4.2** - Vue.js 官方路由
- **Pinia 2.1** - Vue 3 专用状态管理库

### UI 组件
- **Element Plus 2.13** - 基于 Vue 3 的组件库
- **@element-plus/icons-vue** - Element Plus 图标库

### 样式
- **CSS3** - 自定义样式和动画
- **CSS Variables** - CSS 变量主题系统
- **Backdrop Filter** - 毛玻璃效果

### 工具库
- **Axios 1.6** - HTTP 客户端
- **ECharts 6.0** - 数据可视化图表
- **XLSX 0.18** - Excel 文件处理

### 后端（可选）
- **Express 4.18** - Node.js Web 框架
- **JWT** - 身份验证
- **Bcrypt** - 密码加密
- **CORS** - 跨域支持

---

## 📱 功能模块

### 1. 首页（Home）
- 博主列表展示
- 搜索和筛选
- 用户管理（管理员）
- 卡片式布局

### 2. 录入博主（Add）
- 批量导入（Excel）
- 手动录入
- 数据验证

### 3. 我的（My）
- 个人中心
- 个人资料编辑
- 头像裁剪上传
- 密码修改

### 4. 统计（Statistics）
- 数据可视化
- 博主状态分布
- 合作数据统计
- 图表导出

### 5. 团队（Team）
- 团队列表
- 创建团队
- 加入团队
- 团队管理

### 6. 团队主页（TeamHome）
- 团队成员展示
- 团队公告
- 帖子发布
- 话题讨论

### 7. 公共论坛（Forums）
- 话题列表
- 发布话题
- 回复评论
- 点赞收藏

### 8. 博主详情（BloggerDetail）
- 详细信息展示
- 跟进记录
- 合作历史
- 操作日志
- 在线编辑

### 9. 管理后台（Admin）
- 用户管理
- 博主管理
- 团队管理
- 系统设置

---

## 🎨 核心功能详解

### 智能搜索系统

```vue
<template>
  <EnhancedSearch 
    v-model="query"
    :showHistory="true"
    :showHotTags="true"
    @search="handleSearch"
  />
</template>
```

**特性**：
- 🔍 实时搜索（300ms 防抖）
- 📜 搜索历史记录
- 🔥 热门搜索标签
- 💡 搜索建议
- ⌨️ 键盘快捷键支持

### 手势操作

```javascript
import { initSwipeGesture, initLongPress } from '@/utils/gestures'

// 滑动删除/编辑
initSwipeGesture(element, onSwipeLeft, onSwipeRight)

// 长按更多操作
initLongPress(element, onLongPress, 800)
```

**支持手势**：
- 👆 滑动（左/右）
- 🖱️ 长按
- 👆👆 双击
- 🤏 捏合缩放

### 页面动画

```javascript
// 路由配置
{
  path: '/page',
  component: () => import('@/views/Page.vue'),
  meta: { transition: 'slide-up' }
}
```

**动画效果**：
- 📤 slide-up - 上浮（默认）
- 💨 fade - 淡入淡出
- 🔍 zoom - 缩放
- ➡️ slide - 滑动
- 🚀 push - 推进
- 📄 fold - 折叠
- 🔄 rotate - 旋转
- ⚡ bounce - 弹性

### Header 效果

```javascript
import { triggerHeaderEffect } from '@/utils/headerEffect'

triggerHeaderEffect('.page-header')
```

**视觉效果**：
- 💎 毛玻璃背景
- ✨ 双圆心光晕（左紫右橙）
- 🌟 扫过光泽
- 🌓 深色模式支持

---

## 📚 开发文档

详细的开发文档请查看 [DEVELOPMENT.md](./DEVELOPMENT.md)

### 内容包括：
- 🧩 通用组件 API 文档
- 🔧 工具函数使用指南
- 🎨 样式系统说明
- 🛣️ 路由配置详解
- 📖 最佳实践示例

---

## 🎯 最佳实践

### 1. 使用骨架屏

```vue
<SkeletonCard v-if="loading" />
<BaseCard v-else>
  <p>{{ data.content }}</p>
</BaseCard>
```

### 2. 空状态处理

```vue
<EmptyState 
  title="暂无数据"
  description="点击按钮添加新数据"
>
  <template #action>
    <el-button type="primary">添加数据</el-button>
  </template>
</EmptyState>
```

### 3. 卡片组件

```vue
<BaseCard 
  type="primary" 
  hoverable
  @click="handleClick"
>
  <template #header>标题</template>
  <p>内容</p>
  <template #footer>底部</template>
</BaseCard>
```

### 4. 移动端优化

```css
/* 最小触摸区域 */
button {
  min-height: 44px;
  min-width: 44px;
}

/* 触摸反馈 */
.card:active {
  transform: scale(0.98);
}
```

---

## 🔐 权限说明

### 角色分类
- **普通用户** - 基础功能使用
- **管理员** - 系统管理权限

### 权限控制
- 路由守卫验证
- 按钮级权限控制
- API 接口鉴权
- Token 自动续期

---

## 🌐 浏览器支持

- ✅ Chrome >= 90
- ✅ Firefox >= 88
- ✅ Safari >= 14
- ✅ Edge >= 90
- ✅ 移动端浏览器

---

## 📦 构建优化

### 代码分割
- 路由懒加载
- 组件异步加载
- 第三方库分离

### 性能优化
- Gzip 压缩
- Tree Shaking
- 图片懒加载
- 缓存策略

### 样式优化
- CSS 变量主题
- 按需加载
- 自动前缀

---

## 🤝 贡献指南

### 开发流程

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 代码规范

- 使用 ESLint 检查代码
- 遵循 Vue 3 组合式 API 风格
- 组件命名使用 PascalCase
- 文件命名与组件名一致

### 提交信息

- `feat:` 新功能
- `fix:` 修复 bug
- `docs:` 文档更新
- `style:` 代码格式调整
- `refactor:` 代码重构
- `perf:` 性能优化
- `test:` 测试相关
- `chore:` 构建/工具链相关

---

## 📝 更新日志

### v1.0.0 (2026-03-24)

**新增功能**
- ✅ BaseCard 通用卡片组件
- ✅ EmptyState 空状态组件
- ✅ SkeletonCard/SkeletonList 骨架屏组件
- ✅ EnhancedSearch 增强搜索组件
- ✅ 搜索工具库（历史/热门/建议）
- ✅ 手势工具库（滑动/长按/双击/缩放）
- ✅ Header 效果工具
- ✅ 全局 Header 样式（毛玻璃设计）
- ✅ 动画库（10+ 种过渡效果）

**优化改进**
- ✅ 移动端体验优化（触摸反馈、手势支持）
- ✅ 路由配置优化（懒加载、滚动行为）
- ✅ 页面过渡动画优化（只作用于内容区）
- ✅ 深色模式支持
- ✅ 性能优化（防抖、缓存）

**文档更新**
- ✅ 开发文档（DEVELOPMENT.md）
- ✅ 项目说明文档（README.md）

---

## 📄 开源协议

MIT License

---

## 👨‍💻 开发团队

Orange Development Team

---

## 📞 联系方式

如有问题或建议，请通过以下方式联系我们：

- 📧 Email: team@orange.com
- 💬 Issues: [GitHub Issues](https://github.com/your-repo/orange/issues)

---

## 🙏 致谢

感谢以下开源项目：

- [Vue.js](https://vuejs.org/)
- [Element Plus](https://element-plus.org/)
- [Vite](https://vitejs.dev/)
- [ECharts](https://echarts.apache.org/)
- [Pinia](https://pinia.vuejs.org/)

---

**Made with ❤️ by Orange Team**
