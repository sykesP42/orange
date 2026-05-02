<template>
  <Teleport to="body">
    <Transition name="command-fade">
      <div v-if="isOpen" class="command-overlay" tabindex="-1" ref="overlayRef" @click="close" @keydown.down.prevent="handleKeyDown" @keydown.up.prevent="handleKeyUp" @keydown.enter.prevent="executeSelected" @keydown.tab.prevent="switchTab" @keydown.esc="close" @after-enter="onOverlayEntered">
        <div class="command-palette" @click.stop>
          <div class="command-header">
            <div class="command-search-wrapper">
              <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <line x1="21" y1="21" x2="16.65" y2="16.65"/>
              </svg>
              <input
                ref="searchInput"
                v-model="query"
                type="text"
                class="command-input"
                :placeholder="activeTab === 'all' ? '搜索命令、博主、页面...' : '输入关键词...'"
                @keydown.down.prevent
                @keydown.up.prevent
                @keydown.enter.prevent
                @keydown.tab.prevent
                @keydown.esc="close"
              />
              <kbd v-if="query" class="clear-btn" @click="query = ''">✕</kbd>
              <kbd v-else class="esc-hint">ESC</kbd>
            </div>
          </div>

          <div class="command-tabs">
            <button
              v-for="tab in tabs"
              :key="tab.id"
              class="tab-btn"
              :class="{ active: activeTab === tab.id }"
              @click="activeTab = tab.id"
            >
              <span v-html="tab.icon"></span>
              {{ tab.name }}
              <span v-if="tab.id !== 'all'" class="tab-count">{{ getCount(tab.id) }}</span>
            </button>
          </div>

          <div ref="commandBody" class="command-body">
            <div v-if="activeTab === 'all'" class="all-results">
              <div v-if="filteredCommands.length > 0" class="result-group">
                <div class="group-label">⚡ 命令</div>
                <div
                  v-for="(cmd, index) in filteredCommands.slice(0, 5)"
                  :key="'cmd-' + cmd.id"
                  class="result-item command-item"
                  :class="{ selected: selectedIndex === getGlobalIndex('cmd', index) }"
                  @click="executeCommand(cmd)"
                  @mouseenter="selectedIndex = getGlobalIndex('cmd', index); $event.currentTarget.scrollIntoView({ block: 'center', behavior: 'smooth' })"
                >
                  <div class="item-icon" :style="{ background: cmd.bgColor }">
                    <span v-html="cmd.icon"></span>
                  </div>
                  <div class="item-info">
                    <span class="item-name" v-html="highlightMatch(cmd.name)"></span>
                    <span class="item-desc" v-html="highlightMatch(cmd.description)"></span>
                  </div>
                  <div class="item-shortcut" v-if="cmd.shortcut">
                    <kbd v-for="key in cmd.shortcut" :key="key">{{ key }}</kbd>
                  </div>
                </div>
              </div>

              <div v-if="filteredBloggers.length > 0" class="result-group">
                <div class="group-label">👤 博主</div>
                <div
                  v-for="(blogger, index) in filteredBloggers.slice(0, 5)"
                  :key="'blog-' + blogger.id"
                  class="result-item blogger-item"
                  :class="{ selected: selectedIndex === getGlobalIndex('blog', index) }"
                  @click="goToBlogger(blogger.id)"
                  @mouseenter="selectedIndex = getGlobalIndex('blog', index); $event.currentTarget.scrollIntoView({ block: 'center', behavior: 'smooth' })"
                >
                  <div class="item-avatar" :style="{ backgroundColor: getCategoryColor(blogger.category) }">
                    <img v-if="blogger.avatar" :src="blogger.avatar" :alt="blogger.nickname" />
                    <span v-else>{{ blogger.nickname?.charAt(0) || '?' }}</span>
                  </div>
                  <div class="item-info">
                    <span class="item-name" v-html="highlightMatch(blogger.nickname)"></span>
                    <span class="item-desc">{{ blogger.category }} · {{ blogger.platform || '未知平台' }}</span>
                  </div>
                  <span class="item-arrow">→</span>
                </div>
              </div>

              <div v-if="filteredPages.length > 0" class="result-group">
                <div class="group-label">📄 页面</div>
                <div
                  v-for="(page, index) in filteredPages.slice(0, 5)"
                  :key="'page-' + page.path"
                  class="result-item page-item"
                  :class="{ selected: selectedIndex === getGlobalIndex('page', index) }"
                  @click="goToPage(page.path)"
                  @mouseenter="selectedIndex = getGlobalIndex('page', index); $event.currentTarget.scrollIntoView({ block: 'center', behavior: 'smooth' })"
                >
                  <div class="item-icon" :style="{ background: page.bgColor }">
                    <span v-html="page.icon"></span>
                  </div>
                  <div class="item-info">
                    <span class="item-name" v-html="highlightMatch(page.name)"></span>
                    <span class="item-path">{{ page.path }}</span>
                  </div>
                  <span class="item-arrow">→</span>
                </div>
              </div>

              <div v-if="totalResults === 0 && query" class="empty-state">
                <span>未找到 "{{ query }}" 相关结果</span>
                <p class="empty-hint">试试其他关键词或使用拼音搜索</p>
              </div>
              <div v-if="!query" class="recent-hint">
                <p>💡 输入关键词开始搜索，支持：</p>
                <div class="feature-list">
                  <span>✨ 模糊匹配</span>
                  <span>🔤 拼音搜索</span>
                  <span>🎯 首字母缩写</span>
                  <span>⌨️ 路径跳转 (如 add, my)</span>
                </div>
              </div>
            </div>

            <div v-else-if="activeTab === 'commands'" class="command-list">
              <div
                v-for="(cmd, index) in filteredCommands"
                :key="cmd.id"
                class="command-item"
                :class="{ selected: selectedIndex === index }"
                @click="executeCommand(cmd)"
                @mouseenter="selectedIndex = index; $event.currentTarget.scrollIntoView({ block: 'center', behavior: 'smooth' })"
              >
                <div class="command-icon" :style="{ background: cmd.bgColor }">
                  <span v-html="cmd.icon"></span>
                </div>
                <div class="command-info">
                  <span class="command-name" v-html="highlightMatch(cmd.name)"></span>
                  <span class="command-desc" v-html="highlightMatch(cmd.description)"></span>
                </div>
                <div class="command-shortcut" v-if="cmd.shortcut">
                  <kbd v-for="key in cmd.shortcut" :key="key">{{ key }}</kbd>
                </div>
              </div>
              <div v-if="filteredCommands.length === 0" class="empty-state">
                <span>未找到匹配的命令</span>
              </div>
            </div>

            <div v-else-if="activeTab === 'bloggers'" class="blogger-list">
              <div
                v-for="(blogger, index) in filteredBloggers"
                :key="blogger.id"
                class="blogger-item"
                :class="{ selected: selectedIndex === index }"
                @click="goToBlogger(blogger.id)"
                @mouseenter="selectedIndex = index; $event.currentTarget.scrollIntoView({ block: 'center', behavior: 'smooth' })"
              >
                <div class="blogger-avatar" :style="{ backgroundColor: getCategoryColor(blogger.category) }">
                  <img v-if="blogger.avatar" :src="blogger.avatar" :alt="blogger.nickname" />
                  <span v-else>{{ blogger.nickname?.charAt(0) || '?' }}</span>
                </div>
                <div class="blogger-info">
                  <span class="blogger-name" v-html="highlightMatch(blogger.nickname)"></span>
                  <span class="blogger-meta">{{ blogger.category }} · {{ blogger.platform || '未知平台' }} · {{ blogger.status || '初次联系' }}</span>
                </div>
                <span class="blogger-arrow">→</span>
              </div>
              <div v-if="filteredBloggers.length === 0" class="empty-state">
                <span>未找到匹配的博主</span>
              </div>
            </div>

            <div v-else-if="activeTab === 'pages'" class="page-list">
              <div
                v-for="(page, index) in filteredPages"
                :key="page.path"
                class="page-item"
                :class="{ selected: selectedIndex === index }"
                @click="goToPage(page.path)"
                @mouseenter="selectedIndex = index; $event.currentTarget.scrollIntoView({ block: 'center', behavior: 'smooth' })"
              >
                <div class="page-icon" :style="{ background: page.bgColor }">
                  <span v-html="page.icon"></span>
                </div>
                <div class="page-info">
                  <span class="page-name" v-html="highlightMatch(page.name)"></span>
                  <span class="page-path">{{ page.path }}</span>
                </div>
                <span class="page-arrow">→</span>
              </div>
              <div v-if="filteredPages.length === 0" class="empty-state">
                <span>未找到匹配的页面</span>
              </div>
            </div>
          </div>

          <div class="command-footer">
            <div class="footer-hint">
              <kbd>↑↓</kbd> 导航
              <kbd>↵</kbd> 选择
              <kbd>⇥</kbd> 切换标签
              <kbd>esc</kbd> 关闭
            </div>
            <div class="footer-stats" v-if="activeTab === 'all' && totalResults > 0">
              共 {{ totalResults }} 个结果
            </div>
            <div class="footer-stats" v-else-if="activeTab === 'bloggers'">
              共 {{ filteredBloggers.length }} 位博主
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { bloggerListAPI } from '../api'
import { searchList, highlightText as engineHighlight, getMatchInfo, MATCH_TYPE } from '../utils/search-engine'
import { highlightMatch as pinyinHighlight } from '../utils/pinyin-search'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'close'])

const router = useRouter()
const isOpen = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})
const query = ref('')
const activeTab = ref('all')
const selectedIndex = ref(0)
const searchInput = ref(null)
const commandBody = ref(null)
const overlayRef = ref(null)
const bloggers = ref([])

const close = () => {
  isOpen.value = false
  emit('close')
}

const tabs = [
  { id: 'all', name: '全部', icon: '🔍' },
  { id: 'commands', name: '命令', icon: '⚡' },
  { id: 'bloggers', name: '博主', icon: '👤' },
  { id: 'pages', name: '页面', icon: '📄' }
]

const commands = ref([
  { id: 'add', name: '录入博主', description: '添加新的博主信息', icon: '➕', bgColor: 'rgba(249, 115, 22, 0.1)', shortcut: ['add'], keywords: ['add', '录入', '新增', '创建'], action: () => router.push('/add') },
  { id: 'home', name: '返回首页', description: '查看博主列表', icon: '🏠', bgColor: 'rgba(59, 130, 246, 0.1)', shortcut: ['/'], keywords: ['home', '首页', '列表', '/'], action: () => router.push('/') },
  { id: 'my', name: '我的博主', description: '查看我负责的博主', icon: '👤', bgColor: 'rgba(34, 197, 94, 0.1)', shortcut: ['/my'], keywords: ['my', '我的', '个人', '/my'], action: () => router.push('/my') },
  { id: 'team', name: '团队协作', description: '团队博主管理', icon: '👥', bgColor: 'rgba(99, 102, 241, 0.1)', shortcut: ['/team'], keywords: ['team', '团队', '协作', '/team'], action: () => router.push('/team') },
  { id: 'calendar', name: '日历视图', description: '查看日程安排', icon: '📅', bgColor: 'rgba(236, 72, 153, 0.1)', shortcut: ['/calendar'], keywords: ['calendar', '日历', '日程', '安排', '/calendar'], action: () => router.push('/calendar') },
  { id: 'stats', name: '统计分析', description: '查看数据统计', icon: '📊', bgColor: 'rgba(249, 115, 22, 0.1)', shortcut: ['/statistics'], keywords: ['stats', 'statistics', '统计', '分析', '数据', '/statistics'], action: () => router.push('/statistics') },
  { id: 'admin', name: '管理后台', description: '系统管理设置', icon: '⚙️', bgColor: 'rgba(107, 114, 128, 0.1)', shortcut: ['/admin'], keywords: ['admin', '管理', '后台', '设置', '/admin'], action: () => router.push('/admin') },
  { id: 'dark', name: '切换主题', description: '深色/浅色模式切换', icon: '🌓', bgColor: 'rgba(156, 163, 175, 0.1)', shortcut: ['⌘D'], keywords: ['dark', 'theme', '主题', '切换', '深色', '浅色'], action: () => toggleTheme() },
  { id: 'export', name: '导出数据', description: '导出博主数据', icon: '📤', bgColor: 'rgba(34, 197, 94, 0.1)', keywords: ['export', '导出', '数据', 'excel'], action: () => router.push('/admin?tab=export') },
  { id: 'import', name: '导入数据', description: '批量导入博主', icon: '📥', bgColor: 'rgba(59, 130, 246, 0.1)', keywords: ['import', '导入', '批量'], action: () => router.push('/admin?tab=import') },
  { id: 'invalid', name: '失效博主', description: '查看失效博主库', icon: '⚠️', bgColor: 'rgba(239, 68, 68, 0.1)', keywords: ['invalid', '失效', '失效博主', '/invalid-bloggers'], action: () => router.push('/invalid-bloggers') },
  { id: 'chat', name: '私信聊天', description: '消息中心', icon: '💬', bgColor: 'rgba(59, 130, 246, 0.1)', shortcut: ['/chat'], keywords: ['chat', '私信', '聊天', '消息', '/chat'], action: () => router.push('/chat') }
])

const pages = ref([
  { name: '首页', path: '/', icon: '🏠', bgColor: 'rgba(59, 130, 246, 0.1)', keywords: ['home', '首页', '列表'] },
  { name: '录入博主', path: '/add', icon: '➕', bgColor: 'rgba(249, 115, 22, 0.1)', keywords: ['add', '录入', '新增'] },
  { name: '我的博主', path: '/my', icon: '👤', bgColor: 'rgba(34, 197, 94, 0.1)', keywords: ['my', '我的', '个人'] },
  { name: '团队协作', path: '/team', icon: '👥', bgColor: 'rgba(99, 102, 241, 0.1)', keywords: ['team', '团队', '协作'] },
  { name: '日历视图', path: '/calendar', icon: '📅', bgColor: 'rgba(236, 72, 153, 0.1)', keywords: ['calendar', '日历', '日程'] },
  { name: '统计分析', path: '/statistics', icon: '📊', bgColor: 'rgba(249, 115, 22, 0.1)', keywords: ['statistics', '统计', '分析', '数据'] },
  { name: '失效博主', path: '/invalid-bloggers', icon: '⚠️', bgColor: 'rgba(239, 68, 68, 0.1)', keywords: ['invalid', '失效'] },
  { name: '管理后台', path: '/admin', icon: '⚙️', bgColor: 'rgba(107, 114, 128, 0.1)', keywords: ['admin', '管理', '后台'] },
  { name: '私信聊天', path: '/chat', icon: '💬', bgColor: 'rgba(59, 130, 246, 0.1)', keywords: ['chat', '私信', '聊天', '消息'] }
])

const commandSearchResults = computed(() => {
  if (!query.value) return commands.value.map(cmd => ({ item: cmd, score: 0 }))
  return searchList(commands.value, query.value, {
    fields: ['name', 'description'],
    fieldWeights: { name: 1.2, description: 0.8 },
    keywordsField: 'keywords',
    limit: 20
  })
})

const bloggerSearchResults = computed(() => {
  if (!query.value) return bloggers.value.slice(0, 20).map(b => ({ item: b, score: 0 }))
  return searchList(bloggers.value, query.value, {
    fields: ['nickname', 'category', 'platform'],
    fieldWeights: { nickname: 1.3, category: 0.9, platform: 0.8 },
    keywordsField: 'keywords',
    limit: 20
  })
})

const pageSearchResults = computed(() => {
  if (!query.value) return pages.value.map(p => ({ item: p, score: 0 }))
  return searchList(pages.value, query.value, {
    fields: ['name', 'path'],
    fieldWeights: { name: 1.2, path: 0.7 },
    keywordsField: 'keywords',
    limit: 20
  })
})

const filteredCommands = computed(() => commandSearchResults.value.map(r => r.item))
const filteredBloggers = computed(() => bloggerSearchResults.value.map(r => r.item))
const filteredPages = computed(() => pageSearchResults.value.map(r => r.item))

const totalResults = computed(() => {
  return filteredCommands.value.length + filteredBloggers.value.length + filteredPages.value.length
})

const getGlobalIndex = (type, index) => {
  if (type === 'cmd') return index
  if (type === 'blog') return filteredCommands.value.slice(0, 5).length + index
  if (type === 'page') return filteredCommands.value.slice(0, 5).length + filteredBloggers.value.slice(0, 5).length + index
  return index
}

const getTotalItems = () => {
  if (activeTab.value === 'all') {
    return Math.min(filteredCommands.value.length, 5) +
           Math.min(filteredBloggers.value.length, 5) +
           Math.min(filteredPages.value.length, 5)
  }
  if (activeTab.value === 'commands') return filteredCommands.value.length
  if (activeTab.value === 'bloggers') return filteredBloggers.value.length
  if (activeTab.value === 'pages') return filteredPages.value.length
  return 0
}

const getCount = (tabId) => {
  if (tabId === 'commands') return filteredCommands.value.length
  if (tabId === 'bloggers') return filteredBloggers.value.length
  if (tabId === 'pages') return filteredPages.value.length
  return 0
}

const moveDown = () => {
  const total = getTotalItems()
  if (total === 0) return
  if (selectedIndex.value >= total - 1) {
    selectedIndex.value = 0
    commandBody.value?.scrollTo({ top: 0, behavior: 'instant' })
  } else {
    selectedIndex.value++
    nextTick(() => scrollCurrentIntoView())
  }
}

const moveUp = () => {
  const total = getTotalItems()
  if (total === 0) return
  if (selectedIndex.value <= 0) {
    selectedIndex.value = total - 1
    commandBody.value?.scrollTo({ top: commandBody.value.scrollHeight, behavior: 'instant' })
  } else {
    selectedIndex.value--
    nextTick(() => scrollCurrentIntoView())
  }
}

const handleKeyDown = (e) => {
  e.preventDefault()
  e.stopPropagation()
  moveDown()
}

const handleKeyUp = (e) => {
  e.preventDefault()
  e.stopPropagation()
  moveUp()
}

const scrollCurrentIntoView = () => {
  const container = commandBody.value
  if (!container) return
  const selectedEl = container.querySelector('.result-item.selected, .command-item.selected, .blogger-item.selected, .page-item.selected')
  if (!selectedEl) return
  
  const containerRect = container.getBoundingClientRect()
  const elementRect = selectedEl.getBoundingClientRect()
  const offsetTop = elementRect.top - containerRect.top + container.scrollTop
  const containerHeight = container.clientHeight
  const elementHeight = elementRect.height
  const scrollToPos = offsetTop - (containerHeight / 2) + (elementHeight / 2)
  
  container.scrollTo({
    top: Math.max(0, Math.min(scrollToPos, container.scrollHeight - containerHeight)),
    behavior: 'instant'
  })
}

const switchTab = () => {
  const currentIndex = tabs.findIndex(t => t.id === activeTab.value)
  const nextIndex = (currentIndex + 1) % tabs.length
  activeTab.value = tabs[nextIndex].id
  selectedIndex.value = 0
}

const executeSelected = () => {
  if (activeTab.value === 'all') {
    const cmdCount = Math.min(filteredCommands.value.length, 5)
    const blogCount = Math.min(filteredBloggers.value.length, 5)

    if (selectedIndex.value < cmdCount) {
      executeCommand(filteredCommands.value[selectedIndex.value])
    } else if (selectedIndex.value < cmdCount + blogCount) {
      goToBlogger(filteredBloggers.value[selectedIndex.value - cmdCount].id)
    } else {
      const pageIdx = selectedIndex.value - cmdCount - blogCount
      goToPage(filteredPages.value[pageIdx].path)
    }
    return
  }

  if (activeTab.value === 'commands') {
    if (filteredCommands.value[selectedIndex.value]) {
      executeCommand(filteredCommands.value[selectedIndex.value])
    }
  } else if (activeTab.value === 'bloggers') {
    if (filteredBloggers.value[selectedIndex.value]) {
      goToBlogger(filteredBloggers.value[selectedIndex.value].id)
    }
  } else if (activeTab.value === 'pages') {
    if (filteredPages.value[selectedIndex.value]) {
      goToPage(filteredPages.value[selectedIndex.value].path)
    }
  }
}

const executeCommand = (cmd) => {
  cmd.action()
  close()
}

const goToBlogger = (id) => {
  router.push(`/blogger/${id}`)
  close()
}

const goToPage = (path) => {
  router.push(path.startsWith('/') ? path : '/' + path)
  close()
}

const toggleTheme = () => {
  document.documentElement.classList.toggle('dark')
  localStorage.setItem('theme', document.documentElement.classList.contains('dark') ? 'dark' : 'light')
  close()
}

const highlightMatch = (text) => {
  if (!query.value || !text) return text
  return pinyinHighlight(text, query.value)
}

const getCategoryColor = (category) => {
  const colors = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#06b6d4', '#84cc16']
  let hash = 0
  if (category) {
    for (let i = 0; i < category.length; i++) {
      hash = category.charCodeAt(i) + ((hash << 5) - hash)
    }
  }
  return colors[Math.abs(hash) % colors.length]
}

const loadBloggers = async () => {
  try {
    const res = await bloggerListAPI({ page: 1, pageSize: 100 })
    console.log('博主API返回:', res)
    if (res.code === 200) {
      const list = res.data?.list || res.data || []
      console.log('博主列表:', list.length, '条')
      bloggers.value = list.map(b => ({
        ...b,
        keywords: [b.nickname, b.category, b.platform, b.status].filter(Boolean)
      }))
    } else {
      console.warn('博主API返回非200:', res.code, res.message)
    }
  } catch (e) {
    console.error('加载博主列表失败', e)
  }
}

const open = () => {
  query.value = ''
  selectedIndex.value = 0
  activeTab.value = 'all'
  loadBloggers()
}

const onOverlayEntered = () => {
  focusInput()
  setTimeout(() => {
    overlayRef.value?.focus()
  }, 100)
}

const focusInput = () => {
  const tryFocus = (retries = 10) => {
    if (retries <= 0) return
    const el = searchInput.value
    if (el && document.contains(el)) {
      el.focus()
      el.select()
    } else {
      requestAnimationFrame(() => tryFocus(retries - 1))
    }
  }
  requestAnimationFrame(() => tryFocus())
}

watch(isOpen, (val) => {
  if (val) {
    open()
    focusInput()
  }
})

watch(query, () => {
  selectedIndex.value = 0
})

watch(activeTab, () => {
  selectedIndex.value = 0
})

const handleOpenCommandPalette = () => {
  isOpen.value = true
  open()
}

onMounted(() => {
  window.addEventListener('open-command-palette', handleOpenCommandPalette)
})

onUnmounted(() => {
  window.removeEventListener('open-command-palette', handleOpenCommandPalette)
})
</script>

<style scoped>
.command-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 12vh;
  z-index: 9999;
}

.command-palette {
  width: 90%;
  max-width: 700px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  box-shadow: 0 24px 48px rgba(0, 0, 0, 0.24);
  overflow: hidden;
  max-height: 75vh;
  display: flex;
  flex-direction: column;
}

.command-header {
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.command-search-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
  background: var(--bg-hover);
  border: 2px solid transparent;
  border-radius: 12px;
  padding: 0 16px;
  transition: all 0.2s;
}

.command-search-wrapper:focus-within {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.15);
  background: var(--bg-card);
}

.search-icon {
  width: 20px;
  height: 20px;
  color: var(--primary);
  flex-shrink: 0;
}

.command-input {
  flex: 1;
  height: 48px;
  background: transparent;
  border: none;
  font-size: 16px;
  color: var(--text-primary);
  outline: none;
  font-family: inherit;
}

.command-input::placeholder {
  color: var(--text-muted);
}

.clear-btn,
.esc-hint {
  padding: 4px 10px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.clear-btn:hover,
.esc-hint:hover {
  background: var(--bg-hover);
  border-color: var(--primary);
  color: var(--primary);
}

.command-tabs {
  display: flex;
  gap: 4px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: transparent;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.tab-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.tab-btn.active {
  background: var(--primary);
  color: white;
  font-weight: 500;
}

.tab-count {
  font-size: 10px;
  background: rgba(255, 255, 255, 0.2);
  padding: 1px 6px;
  border-radius: 10px;
  margin-left: 4px;
}

.tab-btn.active .tab-count {
  background: rgba(255, 255, 255, 0.3);
}

.command-body {
  flex: 1;
  overflow-y: auto;
  min-height: 200px;
  max-height: 450px;
}

.all-results {
  padding: 4px 0;
}

.result-group {
  margin-bottom: 4px;
}

.group-label {
  padding: 8px 16px 4px;
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.result-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  cursor: pointer;
  transition: all 0.15s;
  margin: 2px 8px;
  border-radius: 10px;
}

.result-item:hover {
  background: var(--bg-hover);
}

.result-item.selected {
  background: rgba(249, 115, 22, 0.12);
}

.item-icon,
.item-avatar,
.command-icon,
.page-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
}

.item-avatar {
  border-radius: 50%;
}

.item-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item-avatar span {
  font-size: 16px;
  font-weight: 600;
  color: white;
}

.item-info,
.command-info,
.blogger-info,
.page-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.item-name,
.command-name,
.blogger-name,
.page-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  line-height: 1.3;
}

.item-desc,
.command-desc,
.blogger-meta,
.page-path {
  font-size: 12px;
  color: var(--text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-shortcut,
.command-shortcut {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.item-shortcut kbd,
.command-shortcut kbd {
  padding: 3px 7px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 5px;
  font-size: 10px;
  color: var(--text-muted);
  font-family: inherit;
}

.item-arrow,
.blogger-arrow,
.page-arrow {
  color: var(--text-muted);
  font-size: 16px;
  flex-shrink: 0;
}

.command-list,
.blogger-list,
.page-list {
  padding: 8px;
}

.command-item,
.blogger-item,
.page-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.15s;
}

.command-item:hover,
.blogger-item:hover,
.page-item:hover {
  background: var(--bg-hover);
}

.command-item.selected,
.blogger-item.selected,
.page-item.selected {
  background: rgba(249, 115, 22, 0.12);
}

.blogger-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  flex-shrink: 0;
}

.blogger-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.blogger-avatar span {
  font-size: 16px;
  font-weight: 600;
  color: white;
}

.empty-state {
  padding: 48px 24px;
  text-align: center;
  color: var(--text-muted);
  font-size: 14px;
}

.empty-hint {
  margin-top: 8px;
  font-size: 12px;
  opacity: 0.7;
}

.recent-hint {
  padding: 32px 24px;
  text-align: center;
}

.recent-hint p {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 16px;
}

.feature-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: center;
}

.feature-list span {
  padding: 6px 12px;
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 12px;
  color: var(--text-secondary);
}

.highlight {
  background: rgba(249, 115, 22, 0.25);
  color: var(--primary);
  padding: 0 2px;
  border-radius: 2px;
  font-weight: 500;
}

.command-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--bg-secondary);
}

.footer-hint {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: var(--text-muted);
}

.footer-hint kbd {
  padding: 2px 6px;
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 10px;
  margin-right: 2px;
}

.footer-stats {
  font-size: 11px;
  color: var(--text-muted);
}

.command-fade-enter-active,
.command-fade-leave-active {
  transition: opacity 0.2s ease;
}

.command-fade-enter-active .command-palette,
.command-fade-leave-active .command-palette {
  transition: transform 0.2s ease, opacity 0.2s ease;
}

.command-fade-enter-from,
.command-fade-leave-to {
  opacity: 0;
}

.command-fade-enter-from .command-palette,
.command-fade-leave-to .command-palette {
  transform: scale(0.95) translateY(-10px);
  opacity: 0;
}

@media (max-width: 768px) {
  .command-overlay {
    padding-top: 8vh;
  }

  .command-palette {
    width: 96%;
    max-height: 85vh;
  }

  .command-body {
    max-height: 350px;
  }

  .footer-hint {
    flex-wrap: wrap;
    gap: 8px;
  }

  .command-footer {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
