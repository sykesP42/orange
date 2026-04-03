# Orange 项目开发文档

> 本文档记录了项目中所有自研组件、工具函数和最佳实践

---

## 📦 目录

- [通用组件](#通用组件)
- [工具函数](#工具函数)
- [样式系统](#样式系统)
- [路由配置](#路由配置)
- [最佳实践](#最佳实践)

---

## 🧩 通用组件

### 1. BaseCard - 基础卡片组件

**文件位置**: `src/components/BaseCard.vue`

**功能描述**: 通用的卡片容器组件，支持多种类型、悬停效果和覆盖层

**Props**:
```javascript
{
  type: String,          // 卡片类型：'default' | 'primary' | 'success' | 'warning' | 'danger' | 'info'
  hoverable: Boolean,    // 是否可悬停
  clickable: Boolean,    // 是否可点击
  shadow: String         // 阴影显示：'always' | 'hover' | 'never'
}
```

**Slots**:
- `header` - 卡片头部
- `default` - 卡片主体内容
- `footer` - 卡片底部
- `overlay` - 悬停时显示的覆盖层

**Events**:
- `click` - 点击事件（当 clickable=true 时触发）

**使用示例**:
```vue
<BaseCard 
  type="primary" 
  hoverable 
  @click="handleClick"
>
  <template #header>
    <h3>卡片标题</h3>
  </template>
  
  <p>这是卡片内容</p>
  
  <template #footer>
    <span>底部信息</span>
  </template>
  
  <template #overlay>
    <el-button type="primary">操作</el-button>
  </template>
</BaseCard>
```

---

### 2. EmptyState - 空状态组件

**文件位置**: `src/components/EmptyState.vue`

**功能描述**: 统一的空状态展示组件，带有浮动动画效果

**Props**:
```javascript
{
  icon: String,          // 图标名称，默认 'DocumentDelete'
  title: String,         // 标题，默认 '暂无数据'
  description: String,   // 描述文字
  height: Number,        // 组件高度，默认 300
  iconSize: Number       // 图标大小，默认 64
}
```

**Slots**:
- `icon` - 自定义图标
- `action` - 操作按钮区域

**使用示例**:
```vue
<EmptyState 
  icon="FolderOpened"
  title="暂无文件"
  description="还没有上传任何文件，点击下面按钮开始上传"
>
  <template #action>
    <el-button type="primary">上传文件</el-button>
  </template>
</EmptyState>
```

---

### 3. SkeletonCard - 卡片骨架屏

**文件位置**: `src/components/SkeletonCard.vue`

**功能描述**: 卡片加载时的占位组件，带有流光动画效果

**Props**:
```javascript
{
  height: Number,   // 卡片高度，默认 200
  lines: Number     // 文本行数，默认 3
}
```

**使用示例**:
```vue
<!-- 加载时显示骨架屏 -->
<SkeletonCard :height="200" :lines="3" v-if="loading" />

<!-- 加载完成显示实际内容 -->
<BaseCard v-else>
  <p>实际内容</p>
</BaseCard>
```

---

### 4. SkeletonList - 列表骨架屏

**文件位置**: `src/components/SkeletonList.vue`

**功能描述**: 列表加载时的占位组件，模拟列表项结构

**Props**:
```javascript
{
  count: Number   // 列表项数量，默认 5
}
```

**使用示例**:
```vue
<SkeletonList :count="10" v-if="loading" />
```

---

### 5. EnhancedSearch - 增强搜索框

**文件位置**: `src/components/EnhancedSearch.vue`

**功能描述**: 带搜索历史、热门推荐和搜索建议的智能搜索组件

**Props**:
```javascript
{
  modelValue: String,    // 搜索关键词（v-model）
  placeholder: String,   // 占位文字
  showHistory: Boolean,  // 显示搜索历史，默认 true
  showHotTags: Boolean,  // 显示热门标签，默认 true
  debounce: Number,      // 防抖延迟（ms），默认 300
  maxSuggestions: Number // 最大建议数，默认 5
}
```

**Events**:
- `update:modelValue` - 值更新事件
- `search` - 搜索事件（返回搜索关键词）

**使用示例**:
```vue
<EnhancedSearch 
  v-model="searchQuery"
  placeholder="搜索博主、团队..."
  :showHistory="true"
  :showHotTags="true"
  :debounce="300"
  @search="handleSearch"
/>
```

**配合搜索工具使用**:
```javascript
import { SearchHistory, HotSearchTags } from '@/utils/search'

const searchHistory = new SearchHistory()
const hotTags = new HotSearchTags()

// 搜索历史会自动保存，热门标签会根据频率排序
```

---

## 🔧 工具函数

### 1. 搜索工具库

**文件位置**: `src/utils/search.js`

#### useDebounce - 防抖函数

```javascript
/**
 * @param {Function} fn - 需要防抖的函数
 * @param {number} delay - 延迟时间（毫秒）
 * @returns {Function} - 防抖后的函数
 */
export function useDebounce(fn, delay = 300)
```

**使用示例**:
```javascript
const handleInput = useDebounce((value) => {
  // 搜索逻辑
}, 300)
```

---

#### SearchHistory - 搜索历史管理

```javascript
class SearchHistory {
  constructor(key = 'search_history', max = 10)
  
  get()           // 获取搜索历史
  add(query)      // 添加搜索记录
  remove(query)   // 删除搜索记录
  clear()         // 清空搜索历史
}
```

**使用示例**:
```javascript
import { SearchHistory } from '@/utils/search'

const searchHistory = new SearchHistory()

// 添加记录
searchHistory.add('关键词')

// 获取记录
const history = searchHistory.get()

// 删除记录
searchHistory.remove('关键词')

// 清空
searchHistory.clear()
```

---

#### HotSearchTags - 热门搜索标签

```javascript
class HotSearchTags {
  constructor(key = 'hot_search_tags', max = 8)
  
  get()           // 获取热门标签
  update(query)   // 更新热门标签
  getList()       // 获取标签文本列表
  clear()         // 清空热门标签
}
```

**使用示例**:
```javascript
import { HotSearchTags } from '@/utils/search'

const hotTags = new HotSearchTags()

// 更新（搜索时调用）
hotTags.update('关键词')

// 获取列表
const tags = hotTags.getList()
```

---

#### SearchSuggestions - 搜索建议缓存

```javascript
class SearchSuggestions {
  constructor()
  
  get(query)      // 获取缓存的建议
  set(query, data)// 设置建议缓存
  clear()         // 清空缓存
}
```

**缓存策略**: 5 分钟自动过期

---

### 2. 手势工具库

**文件位置**: `src/utils/gestures.js`

#### initSwipeGesture - 滑动手势

```javascript
/**
 * @param {HTMLElement} element - 目标元素
 * @param {Function} onSwipeLeft - 左滑回调
 * @param {Function} onSwipeRight - 右滑回调
 * @returns {Function} - 清理函数
 */
export function initSwipeGesture(element, onSwipeLeft, onSwipeRight)
```

**使用示例**:
```javascript
import { initSwipeGesture } from '@/utils/gestures'

onMounted(() => {
  const cleanup = initSwipeGesture(
    document.querySelector('.card'),
    () => console.log('左滑删除'),
    () => console.log('右滑编辑')
  )
  
  onUnmounted(() => cleanup())
})
```

---

#### initLongPress - 长按手势

```javascript
/**
 * @param {HTMLElement} element - 目标元素
 * @param {Function} onLongPress - 长按回调
 * @param {number} duration - 长按时长（ms），默认 800
 * @returns {Function} - 清理函数
 */
export function initLongPress(element, onLongPress, duration = 800)
```

**特性**: 带震动反馈（需要设备支持）

---

#### initDoubleTap - 双击手势

```javascript
/**
 * @param {HTMLElement} element - 目标元素
 * @param {Function} onDoubleTap - 双击回调
 * @param {number} maxDelay - 双击间隔（ms），默认 300
 * @returns {Function} - 清理函数
 */
export function initDoubleTap(element, onDoubleTap, maxDelay = 300)
```

---

#### initPinchZoom - 捏合缩放

```javascript
/**
 * @param {HTMLElement} element - 目标元素
 * @param {Function} onPinch - 缩放回调（接收 scale 参数）
 * @returns {Function} - 清理函数
 */
export function initPinchZoom(element, onPinch)
```

---

### 3. Header 效果工具

**文件位置**: `src/utils/headerEffect.js`

#### triggerHeaderEffect - 触发 Header 光效

```javascript
/**
 * @param {string} headerSelector - Header 选择器，默认 '.page-header'
 */
export function triggerHeaderEffect(headerSelector = '.page-header')
```

**使用示例**:
```javascript
import { triggerHeaderEffect } from '@/utils/headerEffect'

// 鼠标悬停时触发
const handleMouseEnter = () => {
  triggerHeaderEffect()
}
```

---

#### removeHeaderGlow - 移除 Header 高亮

```javascript
/**
 * @param {string} headerSelector - Header 选择器，默认 '.page-header'
 */
export function removeHeaderGlow(headerSelector = '.page-header')
```

---

## 🎨 样式系统

### 1. 全局 Header 样式

**文件位置**: `src/styles/header.css`

**特性**:
- 现代毛玻璃设计
- 双圆心静态光晕（左紫右橙）
- 扫过光泽效果
- hover 时光晕增强
- 深色模式支持

**类名**:
- `.page-header` - Header 容器
- `.header-left` / `.header-content` - 左侧内容区
- `.header-right` / `.header-stats` - 右侧统计区
- `.header-glow` - 高亮状态（JS 添加）
- `.header-shine` - 扫光动画（JS 添加）

**使用示例**:
```vue
<header class="page-header">
  <div class="header-left">
    <h1>页面标题</h1>
    <p class="page-desc">页面描述</p>
  </div>
  <div class="header-right">
    <span class="stat-item">统计 1</span>
    <span class="stat-item">统计 2</span>
  </div>
</header>
```

---

### 2. 动画库

**文件位置**: `src/styles/transitions.css`

**页面过渡动画**:
- `fade` - 淡入淡出
- `zoom` - 缩放效果
- `slide` - 水平滑动
- `slide-up` - 上浮效果（默认）
- `push` - 推进效果
- `fold` - 折叠效果
- `rotate` - 旋转效果
- `bounce` - 弹性效果
- `quick-slide` - 快速滑动

**组件动画**:
- `component-fade` - 组件淡入淡出
- `list` - 列表过渡
- `list-item` - 列表项过渡
- `modal` - 模态框过渡

**微动画**:
- `pulse` - 脉冲动画
- `shake` - 抖动动画
- `bounce-in` - 弹入动画
- `flip-in` - 翻入动画
- `expand` - 展开动画
- `ripple` - 按钮涟漪效果

**使用示例**:
```vue
<!-- 在路由中设置动画 -->
{
  path: '/home',
  component: () => import('@/views/Home.vue'),
  meta: { transition: 'fade' }
}

<!-- 在组件中使用 -->
<transition name="bounce">
  <div v-if="show">内容</div>
</transition>
```

---

### 3. 基础样式

**文件位置**: `src/style.css`

**CSS 变量**:
```css
:root {
  --bg-primary: #f5f7fa;
  --bg-secondary: #ffffff;
  --text-primary: #1f2937;
  --text-secondary: #4b5563;
  --primary: #f97316;
  --success: #22c55e;
  --warning: #f59e0b;
  --danger: #ef4444;
  /* ... 更多变量 */
}
```

**移动端优化**:
- 触摸反馈（scale 0.95）
- 最小触摸区域 44px
- `-webkit-overflow-scrolling: touch` 平滑滚动
- 防止意外选中文本

---

## 🛣️ 路由配置

### 路由懒加载

所有路由组件都已配置懒加载：

```javascript
{
  path: '/home',
  name: 'Home',
  component: () => import('@/views/Home.vue')
}
```

### 滚动行为

```javascript
scrollBehavior(to, from, savedPosition) {
  if (savedPosition) {
    return savedPosition  // 返回时恢复位置
  } else {
    return { top: 0, behavior: 'smooth' }  // 平滑滚动到顶部
  }
}
```

### 页面过渡动画配置

```javascript
// 在 Layout.vue 中
<router-view v-slot="{ Component, route }">
  <transition 
    :name="route.meta.transition || 'slide-up'" 
    mode="out-in"
    appear
  >
    <component :is="Component" :key="route.fullPath" />
  </transition>
</router-view>
```

**注意**: 动画只作用于内容区域，Header 和导航栏保持静止

---

## 📚 最佳实践

### 1. 组件使用规范

#### 骨架屏使用
```vue
<template>
  <div>
    <SkeletonCard v-if="loading" />
    <BaseCard v-else>
      <p>{{ data.content }}</p>
    </BaseCard>
  </div>
</template>

<script setup>
const loading = ref(true)
const data = ref(null)

onMounted(async () => {
  loading.value = true
  data.value = await fetchData()
  loading.value = false
})
</script>
```

#### 空状态使用
```vue
<template>
  <div>
    <EmptyState 
      v-if="!items.length"
      title="暂无数据"
      description="点击按钮添加新数据"
    >
      <template #action>
        <el-button type="primary">添加数据</el-button>
      </template>
    </EmptyState>
    
    <div v-else>
      <!-- 数据列表 -->
    </div>
  </div>
</template>
```

---

### 2. 搜索功能实现

```vue
<template>
  <EnhancedSearch 
    v-model="query"
    @search="handleSearch"
  />
</template>

<script setup>
import { ref } from 'vue'
import { SearchHistory, HotSearchTags } from '@/utils/search'

const query = ref('')
const searchHistory = new SearchHistory()
const hotTags = new HotSearchTags()

const handleSearch = (keyword) => {
  // 1. 添加到搜索历史
  searchHistory.add(keyword)
  
  // 2. 更新热门标签
  hotTags.update(keyword)
  
  // 3. 执行搜索
  performSearch(keyword)
}
</script>
```

---

### 3. 手势操作实现

```vue
<template>
  <div class="card" ref="cardRef">
    <p>卡片内容</p>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { initSwipeGesture, initLongPress } from '@/utils/gestures'

const cardRef = ref(null)

onMounted(() => {
  // 滑动删除
  const swipeCleanup = initSwipeGesture(
    cardRef.value,
    () => handleDelete(),  // 左滑
    () => handleEdit()     // 右滑
  )
  
  // 长按更多操作
  const longPressCleanup = initLongPress(
    cardRef.value,
    () => handleMore(),
    800
  )
  
  // 清理
  onUnmounted(() => {
    swipeCleanup()
    longPressCleanup()
  })
})
</script>
```

---

### 4. Header 效果使用

```vue
<template>
  <header class="page-header" ref="headerRef">
    <h1>页面标题</h1>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { triggerHeaderEffect, removeHeaderGlow } from '@/utils/headerEffect'

const headerRef = ref(null)

const handleMouseEnter = () => {
  triggerHeaderEffect('.page-header')
}

const handleMouseLeave = () => {
  removeHeaderGlow('.page-header')
}
</script>

<style scoped>
.page-header {
  transition: all 0.4s ease;
}
</style>
```

---

### 5. 移动端优化

#### 触摸友好设计
```vue
<template>
  <!-- 按钮最小触摸区域 -->
  <el-button class="touch-friendly">按钮</el-button>
</template>

<style scoped>
.touch-friendly {
  min-height: 44px;
  min-width: 44px;
}

/* 卡片触摸反馈 */
.card:active {
  transform: scale(0.98);
}
</style>
```

#### 手势增强
```javascript
// 在列表项中添加手势
import { initSwipeGesture } from '@/utils/gestures'

// 左滑删除，右滑编辑
initSwipeGesture(
  document.querySelector('.list-item'),
  () => console.log('删除'),
  () => console.log('编辑')
)
```

---

### 6. 性能优化

#### 防抖优化
```javascript
import { useDebounce } from '@/utils/search'

// 搜索输入防抖
const handleInput = useDebounce((value) => {
  performSearch(value)
}, 300)

// 滚动事件防抖
const handleScroll = useDebounce(() => {
  loadMore()
}, 100)
```

#### 缓存策略
```javascript
import { SearchSuggestions } from '@/utils/search'

const suggestions = new SearchSuggestions()

// 设置缓存
suggestions.set('keyword', ['建议 1', '建议 2'])

// 获取缓存（5 分钟内有效）
const cached = suggestions.get('keyword')
```

---

## 🎯 快速开始

### 1. 创建新页面

```vue
<template>
  <div class="my-page">
    <!-- Header -->
    <header class="page-header">
      <div class="header-left">
        <h1>页面标题</h1>
        <p class="page-desc">页面描述</p>
      </div>
    </header>
    
    <!-- 内容区 -->
    <main>
      <SkeletonCard v-if="loading" />
      <EmptyState v-else-if="!data.length" />
      <BaseCard v-else v-for="item in data" :key="item.id">
        {{ item.content }}
      </BaseCard>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { triggerHeaderEffect } from '@/utils/headerEffect'

const loading = ref(true)
const data = ref([])

onMounted(async () => {
  loading.value = true
  data.value = await fetchData()
  loading.value = false
})
</script>
```

### 2. 添加路由

```javascript
// router/index.js
{
  path: '/my-page',
  name: 'MyPage',
  component: () => import('@/views/MyPage.vue'),
  meta: { 
    transition: 'slide-up',  // 可选动画
    requiresAuth: true       // 需要登录
  }
}
```

### 3. 使用搜索

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

---

## 📝 更新日志

### v1.0.0
- ✅ 新增 BaseCard 通用卡片组件
- ✅ 新增 EmptyState 空状态组件
- ✅ 新增 SkeletonCard/SkeletonList 骨架屏组件
- ✅ 新增 EnhancedSearch 增强搜索组件
- ✅ 新增搜索工具库（search.js）
- ✅ 新增手势工具库（gestures.js）
- ✅ 新增 Header 效果工具（headerEffect.js）
- ✅ 新增全局 Header 样式（header.css）
- ✅ 新增动画库（transitions.css）
- ✅ 优化移动端体验（触摸反馈、手势支持）
- ✅ 优化路由配置（懒加载、滚动行为）
- ✅ 优化页面过渡动画（只作用于内容区）

---

## 🤝 贡献指南

### 添加新组件

1. 在 `src/components/` 创建 `.vue` 文件
2. 遵循现有组件的代码风格
3. 添加完整的 Props、Events、Slots 文档
4. 提供使用示例

### 添加工具函数

1. 在 `src/utils/` 创建 `.js` 文件
2. 使用 JSDoc 注释
3. 提供清晰的 API 文档
4. 添加使用示例

### 样式规范

1. 使用 CSS 变量
2. 遵循 BEM 命名规范
3. 支持深色模式
4. 移动端优先

---

## 📞 联系方式

如有问题或建议，请联系开发团队。

---

**最后更新**: 2026-03-24
