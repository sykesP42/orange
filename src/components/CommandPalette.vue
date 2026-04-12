<template>
  <Teleport to="body">
    <Transition name="command-fade">
      <div v-if="isOpen" class="cmd-overlay" @click="close" @keydown.esc="close">
        <div class="cmd-palette" @click.stop>
          <div class="cmd-search">
            <svg class="cmd-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
            <input
              ref="searchInput"
              v-model="query"
              type="text"
              class="cmd-input"
              placeholder="搜索或输入命令..."
              @keydown.down.prevent="moveDown"
              @keydown.up.prevent="moveUp"
              @keydown.enter.prevent="executeSelected"
              @keydown.tab.prevent="toggleFilter"
            />
            <div class="cmd-filter" v-if="activeFilter">
              <span>{{ activeFilterLabel }}</span>
              <button @click="activeFilter = null">✕</button>
            </div>
          </div>

          <div class="cmd-body" ref="commandBody">
            <div v-if="isLoadingBloggers" class="cmd-loading">
              <div class="cmd-spinner"></div>
            </div>

            <template v-else>
              <div v-if="!query && recentQueries.length > 0" class="cmd-section">
                <div class="section-label">最近搜索</div>
                <div class="cmd-row-list">
                  <div
                    v-for="(q, i) in recentQueries"
                    :key="i"
                    class="cmd-row query-row"
                    :data-idx="i"
                    :class="{ selected: selectedIndex === i }"
                    @click="query = q; selectedIndex = 0"
                    @mouseenter="selectedIndex = i"
                  >
                    <span class="row-icon">🕐</span>
                    <span class="row-text">{{ q }}</span>
                  </div>
                </div>
              </div>

              <template v-if="query">
                <div class="cmd-section" v-if="filteredActions.length > 0">
                  <div class="section-label">操作</div>
                  <div class="cmd-row-list">
                    <div
                      v-for="(action, i) in filteredActions"
                      :key="action.id"
                      class="cmd-row"
                      :data-idx="getIndex('action', i)"
                      :class="{ selected: selectedIndex === getIndex('action', i) }"
                      @click="executeAction(action)"
                      @mouseenter="selectedIndex = getIndex('action', i)"
                    >
                      <span class="row-icon" :style="{ background: action.bg }">{{ action.icon }}</span>
                      <div class="row-content">
                        <span class="row-title" v-html="highlight(action.name)"></span>
                        <span class="row-desc" v-html="highlight(action.desc)"></span>
                      </div>
                      <kbd class="row-shortcut" v-if="action.key">{{ action.key }}</kbd>
                    </div>
                  </div>
                </div>

                <div class="cmd-section" v-if="filteredBloggers.length > 0">
                  <div class="section-label">博主 <span class="count">({{ filteredBloggers.length }})</span></div>
                  <div class="cmd-row-list">
                    <div
                      v-for="(blogger, i) in filteredBloggers"
                      :key="blogger.id"
                      class="cmd-row"
                      :data-idx="getIndex('blogger', i)"
                      :class="{ selected: selectedIndex === getIndex('blogger', i) }"
                      @click="goToBlogger(blogger.id)"
                      @mouseenter="selectedIndex = getIndex('blogger', i)"
                    >
                      <span class="row-avatar" :style="{ background: getColor(blogger.category) }">
                        <img v-if="blogger.avatar" :src="blogger.avatar" />
                        <span v-else>{{ blogger.nickname?.charAt(0) }}</span>
                      </span>
                      <div class="row-content">
                        <span class="row-title" v-html="highlight(blogger.nickname)"></span>
                        <div class="row-meta">
                          <span class="meta-tag" v-if="blogger.category">{{ blogger.category }}</span>
                          <span class="meta-tag" :class="getStatusClass(blogger.status)" v-if="blogger.status">{{ blogger.status }}</span>
                        </div>
                      </div>
                      <span class="row-arrow">→</span>
                    </div>
                  </div>
                </div>

                <div class="cmd-section" v-if="filteredPages.length > 0">
                  <div class="section-label">页面</div>
                  <div class="cmd-row-list">
                    <div
                      v-for="(page, i) in filteredPages"
                      :key="page.path"
                      class="cmd-row"
                      :data-idx="getIndex('page', i)"
                      :class="{ selected: selectedIndex === getIndex('page', i) }"
                      @click="goToPage(page.path)"
                      @mouseenter="selectedIndex = getIndex('page', i)"
                    >
                      <span class="row-icon" :style="{ background: page.bg }">{{ page.icon }}</span>
                      <div class="row-content">
                        <span class="row-title" v-html="highlight(page.name)"></span>
                        <span class="row-path">{{ page.path }}</span>
                      </div>
                      <span class="row-arrow">→</span>
                    </div>
                  </div>
                </div>

                <div class="cmd-empty" v-if="totalResults === 0">
                  <span class="empty-icon">🔍</span>
                  <span>未找到 "{{ query }}" 相关结果</span>
                </div>
              </template>

              <template v-else>
                <div class="cmd-section" v-if="quickActions.length > 0">
                  <div class="section-label">快速操作</div>
                  <div class="cmd-row-list">
                    <div
                      v-for="(action, i) in quickActions"
                      :key="action.id"
                      class="cmd-row"
                      :data-idx="i"
                      :class="{ selected: selectedIndex === i }"
                      @click="executeAction(action)"
                      @mouseenter="selectedIndex = i"
                    >
                      <span class="row-icon" :style="{ background: action.bg }">{{ action.icon }}</span>
                      <div class="row-content">
                        <span class="row-title">{{ action.name }}</span>
                        <span class="row-desc">{{ action.desc }}</span>
                      </div>
                      <kbd class="row-shortcut" v-if="action.key">{{ action.key }}</kbd>
                    </div>
                  </div>
                </div>

                <div class="cmd-section" v-if="navPages.length > 0">
                  <div class="section-label">导航</div>
                  <div class="cmd-row-list">
                    <div
                      v-for="(page, i) in navPages"
                      :key="page.path"
                      class="cmd-row"
                      :data-idx="getIndex('nav', i)"
                      :class="{ selected: selectedIndex === getIndex('nav', i) }"
                      @click="goToPage(page.path)"
                      @mouseenter="selectedIndex = getIndex('nav', i)"
                    >
                      <span class="row-icon" :style="{ background: page.bg }">{{ page.icon }}</span>
                      <div class="row-content">
                        <span class="row-title">{{ page.name }}</span>
                      </div>
                      <kbd class="row-shortcut" v-if="page.key">{{ page.key }}</kbd>
                    </div>
                  </div>
                </div>
              </template>
            </template>
          </div>

          <div class="cmd-footer">
            <div class="footer-keys">
              <span><kbd>↑↓</kbd> 导航</span>
              <span><kbd>↵</kbd> 确认</span>
              <span><kbd>Tab</kbd> 切换类型</span>
              <span><kbd>Esc</kbd> 关闭</span>
            </div>
            <div class="footer-info" v-if="selectedItem">
              <span>{{ selectedItem.name || selectedItem.nickname || selectedItem.title }}</span>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { bloggerListAPI } from '../api'

const router = useRouter()

const props = defineProps({
  visible: { type: Boolean, default: false }
})

const emit = defineEmits(['update:visible', 'close'])

const isOpen = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

const query = ref('')
const selectedIndex = ref(0)
const searchInput = ref(null)
const commandBody = ref(null)
const bloggers = ref([])
const isLoadingBloggers = ref(false)
const activeFilter = ref(null)

const filters = [
  { id: 'action', label: '操作' },
  { id: 'blogger', label: '博主' },
  { id: 'page', label: '页面' }
]

const activeFilterLabel = computed(() => filters.find(f => f.id === activeFilter.value)?.label || '')

const toggleFilter = () => {
  const current = filters.findIndex(f => f.id === activeFilter.value)
  const next = (current + 1) % (filters.length + 1)
  activeFilter.value = next === filters.length ? null : filters[next].id
  selectedIndex.value = 0
}

const close = () => {
  isOpen.value = false
  emit('close')
}

const recentQueries = ref(JSON.parse(localStorage.getItem('orange_recent_queries') || '[]'))

const saveQuery = (q) => {
  if (!q || q.trim().length < 2) return
  const list = [q.trim(), ...recentQueries.value.filter(s => s !== q.trim())].slice(0, 5)
  recentQueries.value = list
  localStorage.setItem('orange_recent_queries', JSON.stringify(list))
}

const highlight = (text) => {
  if (!query.value || !text) return text || ''
  const idx = text.toLowerCase().indexOf(query.value.toLowerCase())
  if (idx === -1) return text
  return text.slice(0, idx) + '<mark>' + text.slice(idx, idx + query.value.length) + '</mark>' + text.slice(idx + query.value.length)
}

const actions = ref([
  { id: 'add', name: '录入博主', desc: '添加新博主信息', icon: '➕', bg: 'rgba(249,115,22,0.15)', key: '>add', action: () => router.push('/add') },
  { id: 'home', name: '返回首页', desc: '跳转到首页', icon: '🏠', bg: 'rgba(59,130,246,0.15)', key: '>home', action: () => router.push('/') },
  { id: 'my', name: '个人中心', desc: '管理个人信息', icon: '👤', bg: 'rgba(34,197,94,0.15)', key: '>my', action: () => router.push('/my') },
  { id: 'team', name: '团队中心', desc: '管理团队协作', icon: '👥', bg: 'rgba(99,102,241,0.15)', key: '>team', action: () => router.push('/team') },
  { id: 'stats', name: '数据统计', desc: '查看数据报表', icon: '📊', bg: 'rgba(249,115,22,0.15)', key: '>stats', action: () => router.push('/statistics') },
  { id: 'calendar', name: '提醒日历', desc: '查看待办事项', icon: '📅', bg: 'rgba(236,72,153,0.15)', key: '>cal', action: () => router.push('/calendar') },
  { id: 'admin', name: '系统管理', desc: '管理后台设置', icon: '⚙️', bg: 'rgba(107,114,128,0.15)', key: '>admin', action: () => router.push('/admin') },
  { id: 'theme', name: '切换主题', desc: '深色/浅色模式', icon: '🌓', bg: 'rgba(156,163,175,0.15)', key: '>theme', action: () => toggleTheme() },
  { id: 'export', name: '导出数据', desc: '导出博主数据', icon: '📤', bg: 'rgba(34,197,94,0.15)', key: '>export', action: () => router.push('/admin?tab=export') },
  { id: 'import', name: '导入数据', desc: '导入博主数据', icon: '📥', bg: 'rgba(59,130,246,0.15)', key: '>import', action: () => router.push('/admin?tab=import') }
])

const quickActions = computed(() => actions.value.slice(0, 6))

const navPages = ref([
  { name: '首页', path: '/', icon: '🏠', bg: 'rgba(59,130,246,0.15)', key: '/home' },
  { name: '录入博主', path: '/add', icon: '➕', bg: 'rgba(249,115,22,0.15)', key: '/add' },
  { name: '个人中心', path: '/my', icon: '👤', bg: 'rgba(34,197,94,0.15)', key: '/my' },
  { name: '团队中心', path: '/team', icon: '👥', bg: 'rgba(99,102,241,0.15)', key: '/team' },
  { name: '数据统计', path: '/statistics', icon: '📊', bg: 'rgba(249,115,22,0.15)', key: '/stats' },
  { name: '提醒日历', path: '/calendar', icon: '📅', bg: 'rgba(236,72,153,0.15)', key: '/cal' },
  { name: '系统管理', path: '/admin', icon: '⚙️', bg: 'rgba(107,114,128,0.15)', key: '/admin' }
])

const filteredActions = computed(() => {
  if (!query.value) return []
  const q = query.value.toLowerCase()
  return actions.value.filter(a =>
    (!activeFilter.value || activeFilter.value === 'action') &&
    (a.name.toLowerCase().includes(q) || a.desc.toLowerCase().includes(q))
  )
})

const filteredBloggers = computed(() => {
  if (!query.value) return []
  const q = query.value.toLowerCase().trim()
  if (!q) return []
  return bloggers.value.filter(b => {
    if (activeFilter.value && activeFilter.value !== 'blogger') return false
    const fields = [b.nickname, b.category, b.platform, b.status, b.product].filter(Boolean).map(f => f.toLowerCase())
    return fields.some(f => f.includes(q)) ||
      (b.tags || []).some(t => t.toLowerCase().includes(q))
  }).slice(0, 8)
})

const filteredPages = computed(() => {
  if (!query.value) return []
  const q = query.value.toLowerCase()
  return navPages.value.filter(p =>
    (!activeFilter.value || activeFilter.value === 'page') &&
    (p.name.toLowerCase().includes(q) || p.path.toLowerCase().includes(q))
  )
})

const totalResults = computed(() =>
  filteredActions.value.length + filteredBloggers.value.length + filteredPages.value.length
)

const globalItems = computed(() => {
  if (query.value) {
    const items = []
    if (filteredActions.value.length) {
      items.push(...filteredActions.value.map((a, i) => ({ type: 'action', data: a, index: i })))
    }
    if (filteredBloggers.value.length) {
      items.push(...filteredBloggers.value.map((b, i) => ({ type: 'blogger', data: b, index: i })))
    }
    if (filteredPages.value.length) {
      items.push(...filteredPages.value.map((p, i) => ({ type: 'page', data: p, index: i })))
    }
    return items
  } else {
    const items = []
    items.push(...quickActions.value.map((a, i) => ({ type: 'quick', data: a, index: i })))
    const navOffset = items.length
    items.push(...navPages.value.map((p, i) => ({ type: 'nav', data: p, index: i })))
    return items
  }
})

const selectedItem = computed(() => globalItems.value[selectedIndex.value]?.data || null)

const getIndex = (type, localIndex) => {
  return globalItems.value.findIndex(item => item.type === type && item.index === localIndex)
}

const scrollToSelected = () => {
  nextTick(() => {
    const container = commandBody.value
    if (!container) return

    if (selectedIndex.value === 0) {
      container.scrollTop = 0
      return
    }

    const item = globalItems.value[selectedIndex.value]
    if (!item) return
    const el = container.querySelector(`[data-idx="${selectedIndex.value}"]`)
    if (!el) return

    const elTop = el.offsetTop
    const elBottom = elTop + el.offsetHeight
    const viewTop = container.scrollTop
    const viewBottom = viewTop + container.clientHeight

    if (elTop < viewTop) {
      container.scrollTop = Math.max(0, elTop - 12)
    } else if (elBottom > viewBottom) {
      container.scrollTop = elBottom - container.clientHeight + 12
    }
  })
}

const moveDown = () => {
  const len = globalItems.value.length
  if (!len) return
  selectedIndex.value = Math.min(selectedIndex.value + 1, len - 1)
  scrollToSelected()
}

const moveUp = () => {
  if (!globalItems.value.length) return
  selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
  scrollToSelected()
}

const executeSelected = () => {
  const item = selectedItem.value
  if (!item) return
  const type = globalItems.value[selectedIndex.value]?.type
  if (type === 'action' || type === 'quick') executeAction(item)
  else if (type === 'blogger') goToBlogger(item.id)
  else if (type === 'page' || type === 'nav') goToPage(item.path)
}

const executeAction = (action) => {
  saveQuery(query.value)
  action.action?.()
  close()
}

const goToBlogger = (id) => {
  saveQuery(query.value)
  router.push(`/blogger/${id}`)
  close()
}

const goToPage = (path) => {
  saveQuery(query.value)
  router.push(path)
  close()
}

const toggleTheme = () => {
  document.documentElement.classList.toggle('dark')
  localStorage.setItem('theme', document.documentElement.classList.contains('dark') ? 'dark' : 'light')
  close()
}

const getColor = (str) => {
  const colors = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899']
  if (!str) return colors[0]
  let hash = 0
  for (let i = 0; i < str.length; i++) hash = str.charCodeAt(i) + ((hash << 5) - hash)
  return colors[Math.abs(hash) % colors.length]
}

const getStatusClass = (status) => {
  const map = { '初次联系': 's-init', '洽谈中': 's-talk', '已合作': 's-ok', '已拒绝': 's-no', '暂停': 's-pause' }
  return map[status] || ''
}

const loadBloggers = async () => {
  if (isLoadingBloggers.value) return
  isLoadingBloggers.value = true
  try {
    const res = await bloggerListAPI({ page: 1, pageSize: 200 })
    if (res.code === 200) bloggers.value = res.data?.list || res.data || []
  } catch (e) {
    bloggers.value = []
  } finally {
    isLoadingBloggers.value = false
  }
}

const open = () => {
  query.value = ''
  selectedIndex.value = 0
  activeFilter.value = null
  loadBloggers()
  nextTick(() => searchInput.value?.focus())
}

watch(query, () => { selectedIndex.value = 0 })
watch(isOpen, val => { if (val) open() })

onMounted(() => window.addEventListener('open-command-palette', () => { isOpen.value = true; open() }))
onUnmounted(() => window.removeEventListener('open-command-palette', () => {}))
</script>

<style scoped>
.cmd-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 6vh;
  z-index: 9999;
}

.cmd-palette {
  width: 94%;
  max-width: 720px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.35);
  overflow: hidden;
  max-height: 72vh;
  display: flex;
  flex-direction: column;
}

.cmd-search {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 16px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
}

.cmd-icon {
  width: 20px;
  height: 20px;
  color: var(--text-muted);
  flex-shrink: 0;
}

.cmd-input {
  flex: 1;
  height: 56px;
  background: transparent;
  border: none;
  font-size: 16px;
  color: var(--text-primary);
  outline: none;
}

.cmd-input::placeholder { color: var(--text-muted); }

.cmd-filter {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: var(--primary);
  color: white;
  border-radius: 6px;
  font-size: 12px;
}

.cmd-filter button {
  background: none;
  border: none;
  color: inherit;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.cmd-body {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
  min-height: 280px;
  max-height: 520px;
}

.cmd-section {
  margin-bottom: 4px;
}

.section-label {
  padding: 6px 16px;
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.section-label .count {
  font-weight: 400;
  margin-left: 4px;
}

.cmd-row-list {
  display: flex;
  flex-direction: column;
}

.cmd-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  cursor: pointer;
  transition: background 0.1s;
  border-left: 3px solid transparent;
}

.cmd-row:hover { background: var(--bg-hover); }

.cmd-row.selected {
  background: linear-gradient(90deg, rgba(249,115,22,0.1) 0%, transparent 100%);
  border-left-color: var(--primary);
}

.row-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
}

.row-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  color: white;
  flex-shrink: 0;
  overflow: hidden;
}

.row-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.row-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.row-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.row-title :deep(mark) {
  background: rgba(249,115,22,0.25);
  color: var(--primary);
  border-radius: 2px;
  padding: 0 1px;
}

.row-desc, .row-path {
  font-size: 12px;
  color: var(--text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.row-meta {
  display: flex;
  gap: 6px;
}

.meta-tag {
  font-size: 11px;
  padding: 1px 6px;
  background: var(--bg-secondary);
  border-radius: 4px;
  color: var(--text-secondary);
}

.meta-tag.s-init { background: rgba(59,130,246,0.15); color: var(--info); }
.meta-tag.s-talk { background: rgba(249,115,22,0.15); color: var(--primary); }
.meta-tag.s-ok { background: rgba(34,197,94,0.15); color: var(--success); }
.meta-tag.s-no { background: rgba(239,68,68,0.15); color: var(--danger); }
.meta-tag.s-pause { background: rgba(107,114,128,0.15); color: var(--text-tertiary); }

.row-badge {
  font-size: 11px;
  padding: 2px 8px;
  background: var(--primary);
  color: white;
  border-radius: 10px;
  flex-shrink: 0;
}

.row-shortcut {
  padding: 3px 8px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 11px;
  color: var(--text-muted);
  font-family: inherit;
  flex-shrink: 0;
}

.row-arrow {
  color: var(--text-muted);
  font-size: 14px;
  flex-shrink: 0;
}

.cmd-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 48px 24px;
  color: var(--text-muted);
}

.empty-icon { font-size: 36px; }

.cmd-loading {
  display: flex;
  justify-content: center;
  padding: 40px;
}

.cmd-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--border-color);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.cmd-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
}

.footer-keys {
  display: flex;
  gap: 16px;
}

.footer-keys span {
  font-size: 11px;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  gap: 4px;
}

.footer-keys kbd {
  padding: 2px 5px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 3px;
  font-size: 10px;
  font-family: inherit;
}

.footer-info {
  font-size: 12px;
  color: var(--text-secondary);
  max-width: 180px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.query-row .row-icon { font-size: 14px; }

.command-fade-enter-active, .command-fade-leave-active { transition: all 0.15s ease; }
.command-fade-enter-from, .command-fade-leave-to { opacity: 0; transform: scale(0.97) translateY(-8px); }
</style>