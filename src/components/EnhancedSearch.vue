<template>
  <div class="enhanced-search" :class="{ 'is-focus': isFocus, 'has-value': modelValue }">
    <div class="search-input-wrapper">
      <el-icon class="search-icon"><Search /></el-icon>
      <input
        ref="inputRef"
        v-model="inputValue"
        type="text"
        :placeholder="placeholder"
        class="search-input"
        @focus="handleFocus"
        @blur="handleBlur"
        @input="handleInput"
        @keydown.enter="handleSearch"
        @keydown.esc="handleClear"
        @keydown.down.prevent="moveDown"
        @keydown.up.prevent="moveUp"
      />
      <div v-if="isLoading" class="loading-spinner">
        <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" class="spin">
          <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
        </svg>
      </div>
      <el-icon
        v-else-if="modelValue"
        class="clear-icon"
        @click="handleClear"
      >
        <CircleClose />
      </el-icon>
    </div>

    <div v-if="searchModeLabel && isFocus && inputValue" class="search-mode-hint">
      {{ searchModeLabel }}
    </div>

    <transition name="el-fade-in">
      <div v-if="isFocus && showPanel" class="search-suggestions">
        <div v-if="showHistory && searchHistoryList.length > 0 && !inputValue" class="suggestion-section">
          <div class="suggestion-header">
            <span class="section-title">搜索历史</span>
            <el-icon class="clear-btn" @click="clearHistory"><Delete /></el-icon>
          </div>
          <div class="suggestion-tags">
            <el-tag
              v-for="(item, index) in searchHistoryList"
              :key="'h-' + index"
              size="small"
              class="history-tag"
              @click="selectSuggestion(item)"
            >
              {{ item }}
              <el-icon class="tag-close" @click.stop="removeHistory(index)"><Close /></el-icon>
            </el-tag>
          </div>
        </div>

        <div v-if="showHotTags && hotTagsList.length > 0 && !inputValue" class="suggestion-section">
          <div class="suggestion-header">
            <span class="section-title">热门搜索</span>
            <el-icon class="refresh-btn" @click="refreshHotTags"><Refresh /></el-icon>
          </div>
          <div class="suggestion-tags">
            <el-tag
              v-for="(tag, index) in hotTagsList"
              :key="'hot-' + index"
              size="small"
              class="hot-tag"
              :type="index < 3 ? 'warning' : ''"
              @click="selectSuggestion(tag)"
            >
              <span v-if="index < 3" class="hot-index">{{ index + 1 }}</span>
              {{ tag }}
            </el-tag>
          </div>
        </div>

        <div v-if="isLoading" class="suggestion-section loading-section">
          <div class="loading-dots">
            <span></span><span></span><span></span>
          </div>
          <p class="loading-text">搜索中...</p>
        </div>

        <div v-else-if="suggestions.length > 0" class="suggestion-section">
          <div class="suggestion-header">
            <span class="section-title">搜索建议</span>
            <span class="suggestion-count">{{ suggestions.length }} 个结果</span>
          </div>
          <ul class="suggestion-list" ref="suggestionListRef">
            <li
              v-for="(item, index) in suggestions"
              :key="'s-' + index"
              class="suggestion-item"
              :class="{ 'is-active': activeIndex === index }"
              @mouseenter="activeIndex = index"
              @click="selectSuggestion(item.text || item)"
            >
              <el-icon><Search /></el-icon>
              <span class="item-text" v-html="highlightKeyword(item.text || item)"></span>
              <span v-if="item.category" class="item-category">{{ item.category }}</span>
            </li>
          </ul>
        </div>

        <div v-else-if="inputValue && !isLoading && hasSearched" class="suggestion-section">
          <div class="empty-suggestions">
            <el-icon :size="36"><Search /></el-icon>
            <p>未找到 "{{ inputValue }}" 相关结果</p>
            <p class="empty-hint">试试拼音搜索，如输入"zfb"匹配"支付宝"</p>
          </div>
        </div>

        <div v-if="!inputValue && !showHistory && !showHotTags" class="suggestion-section">
          <div class="empty-suggestions">
            <el-icon :size="36"><Search /></el-icon>
            <p>输入关键词开始搜索</p>
            <p class="empty-hint">支持中文、拼音首字母、完整拼音搜索</p>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { Search, CircleClose, Delete, Close, Refresh } from '@element-plus/icons-vue'
import { SearchHistory, HotSearchTags } from '@/utils/search'
import { getSearchMode, highlightMatch } from '@/utils/pinyin-search'
import { bloggerSuggestionsAPI } from '@/api'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: '搜索...'
  },
  showHistory: {
    type: Boolean,
    default: true
  },
  showHotTags: {
    type: Boolean,
    default: true
  },
  debounce: {
    type: Number,
    default: 200
  },
  maxSuggestions: {
    type: Number,
    default: 8
  }
})

const emit = defineEmits(['update:modelValue', 'search'])

const inputValue = ref(props.modelValue)
const isFocus = ref(false)
const activeIndex = ref(-1)
const suggestions = ref([])
const isLoading = ref(false)
const hasSearched = ref(false)
const inputRef = ref(null)
const suggestionListRef = ref(null)

const searchHistory = new SearchHistory()
const hotTags = new HotSearchTags()
const searchHistoryList = ref(searchHistory.get())
const hotTagsList = ref(hotTags.getList())

let debounceTimer = null
let abortController = null

const searchModeLabel = computed(() => {
  if (!inputValue.value) return ''
  const mode = getSearchMode(inputValue.value)
  return mode.label || ''
})

const showPanel = computed(() => {
  return true
})

const handleInput = () => {
  emit('update:modelValue', inputValue.value)
  activeIndex.value = -1

  if (debounceTimer) clearTimeout(debounceTimer)

  if (!inputValue.value || !inputValue.value.trim()) {
    suggestions.value = []
    hasSearched.value = false
    isLoading.value = false
    return
  }

  isLoading.value = true

  debounceTimer = setTimeout(() => {
    fetchSuggestions(inputValue.value.trim())
  }, props.debounce)
}

const fetchSuggestions = async (query) => {
  if (!query) {
    suggestions.value = []
    isLoading.value = false
    return
  }

  if (abortController) {
    abortController.abort()
  }

  try {
    abortController = new AbortController()
    const res = await bloggerSuggestionsAPI(query)

    if (res.code === 200 && res.data) {
      const results = []

      if (res.data.nicknames && res.data.nicknames.length > 0) {
        results.push(...res.data.nicknames.slice(0, props.maxSuggestions).map(n => ({
          text: n.label || n.value || n,
          category: '博主'
        })))
      }

      if (res.data.categories && res.data.categories.length > 0) {
        results.push(...res.data.categories.slice(0, 3).map(c => ({
          text: c.label || c.value || c,
          category: '分类'
        })))
      }

      if (res.data.platforms && res.data.platforms.length > 0) {
        results.push(...res.data.platforms.slice(0, 3).map(p => ({
          text: p.label || p.value || p,
          category: '平台'
        })))
      }

      suggestions.value = results.slice(0, props.maxSuggestions)
    } else {
      suggestions.value = []
    }
  } catch (e) {
    if (e.name !== 'AbortError') {
      suggestions.value = []
    }
  } finally {
    isLoading.value = false
    hasSearched.value = true
  }
}

const handleFocus = () => {
  isFocus.value = true
  searchHistoryList.value = searchHistory.get()
  hotTagsList.value = hotTags.getList()
}

const handleBlur = () => {
  setTimeout(() => {
    isFocus.value = false
    activeIndex.value = -1
  }, 200)
}

const handleSearch = () => {
  const query = inputValue.value.trim()
  if (!query) return

  searchHistory.add(query)
  searchHistoryList.value = searchHistory.get()

  hotTags.update(query)
  hotTagsList.value = hotTags.getList()

  emit('search', query)

  suggestions.value = []
  isFocus.value = false
}

const handleClear = () => {
  inputValue.value = ''
  emit('update:modelValue', '')
  suggestions.value = []
  activeIndex.value = -1
  hasSearched.value = false
  isLoading.value = false
  if (inputRef.value) inputRef.value.focus()
}

const selectSuggestion = (value) => {
  const text = typeof value === 'object' ? (value.text || value.label || value.value) : value
  inputValue.value = text
  emit('update:modelValue', text)
  handleSearch()
}

const removeHistory = (index) => {
  const item = searchHistoryList.value[index]
  searchHistory.remove(item)
  searchHistoryList.value = searchHistory.get()
}

const clearHistory = () => {
  searchHistory.clear()
  searchHistoryList.value = []
}

const refreshHotTags = () => {
  hotTagsList.value = hotTags.getList()
}

const moveDown = () => {
  if (suggestions.value.length === 0) return
  if (activeIndex.value < suggestions.value.length - 1) {
    activeIndex.value++
  } else {
    activeIndex.value = 0
  }
  nextTick(() => scrollToActive())
}

const moveUp = () => {
  if (suggestions.value.length === 0) return
  if (activeIndex.value > 0) {
    activeIndex.value--
  } else {
    activeIndex.value = suggestions.value.length - 1
  }
  nextTick(() => scrollToActive())
}

const scrollToActive = () => {
  const listEl = suggestionListRef.value
  if (!listEl) return
  const activeEl = listEl.children[activeIndex.value]
  if (!activeEl) return

  const container = listEl.closest('.search-suggestions') || listEl.parentElement
  if (!container) return

  const containerRect = container.getBoundingClientRect()
  const activeRect = activeEl.getBoundingClientRect()

  if (activeRect.top < containerRect.top) {
    container.scrollTop -= (containerRect.top - activeRect.top)
  } else if (activeRect.bottom > containerRect.bottom) {
    container.scrollTop += (activeRect.bottom - containerRect.bottom)
  }
}

const highlightKeyword = (text) => {
  if (!inputValue.value || !text) return escapeHtml(text || '')
  return highlightMatch(text, inputValue.value)
}

const escapeHtml = (str) => {
  if (!str) return ''
  return String(str).replace(/[&<>"'/]/g, c => ({
    '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;',
    "'": '&#39;', '/': '&#x2F;'
  }[c]))
}

watch(() => props.modelValue, (val) => {
  inputValue.value = val
})

onUnmounted(() => {
  if (debounceTimer) clearTimeout(debounceTimer)
  if (abortController) abortController.abort()
})
</script>

<style scoped>
.enhanced-search {
  position: relative;
  width: 100%;
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: var(--text-muted);
  font-size: 18px;
  pointer-events: none;
  transition: color 0.3s;
}

.search-input {
  width: 100%;
  height: 44px;
  padding: 0 40px;
  border: 2px solid var(--border-color);
  border-radius: 22px;
  background: var(--bg-card);
  color: var(--text-primary);
  font-size: 14px;
  transition: all 0.3s;
  outline: none;
}

.search-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

.search-input:focus ~ .search-icon,
.search-input:focus + .search-icon {
  color: var(--primary);
}

.clear-icon {
  position: absolute;
  right: 12px;
  color: var(--text-muted);
  font-size: 18px;
  cursor: pointer;
  transition: color 0.3s;
}

.clear-icon:hover {
  color: var(--danger);
}

.loading-spinner {
  position: absolute;
  right: 12px;
  color: var(--primary);
  display: flex;
  align-items: center;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.search-mode-hint {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  font-size: 11px;
  color: var(--primary);
  background: rgba(249, 115, 22, 0.08);
  padding: 2px 10px;
  border-radius: 4px;
  white-space: nowrap;
  z-index: 101;
}

.search-suggestions {
  position: absolute;
  top: calc(100% + 24px);
  left: 0;
  right: 0;
  background: var(--bg-card);
  border-radius: var(--radius-lg, 12px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  border: 1px solid var(--border-color);
  padding: 12px;
  z-index: 100;
  max-height: 420px;
  overflow-y: auto;
  animation: suggestDrop 0.18s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes suggestDrop {
  from { opacity: 0; transform: translateY(-6px); }
  to { opacity: 1; transform: translateY(0); }
}

.suggestion-section {
  margin-bottom: 12px;
}

.suggestion-section:last-child {
  margin-bottom: 0;
}

.suggestion-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  padding-bottom: 6px;
  border-bottom: 1px solid var(--border-light, #f0f0f0);
}

.section-title {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 600;
}

.suggestion-count {
  font-size: 11px;
  color: var(--text-muted);
}

.clear-btn, .refresh-btn {
  cursor: pointer;
  color: var(--text-muted);
  transition: color 0.3s;
}

.clear-btn:hover, .refresh-btn:hover {
  color: var(--primary);
}

.suggestion-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.history-tag, .hot-tag {
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 4px;
}

.history-tag:hover, .hot-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.tag-close {
  margin-left: 2px;
  font-size: 10px;
}

.tag-close:hover {
  color: var(--danger);
}

.hot-index {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 14px;
  height: 14px;
  background: var(--danger);
  color: white;
  font-size: 9px;
  font-weight: bold;
  border-radius: 50%;
  margin-right: 2px;
}

.suggestion-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.suggestion-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
}

.suggestion-item:hover,
.suggestion-item.is-active {
  background: var(--bg-hover, rgba(249, 115, 22, 0.06));
}

.suggestion-item .el-icon {
  color: var(--text-muted);
  font-size: 14px;
  flex-shrink: 0;
}

.item-text {
  flex: 1;
  color: var(--text-primary);
  font-size: 13px;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-text :deep(mark.highlight) {
  background: linear-gradient(120deg, #f97316 0%, #fb923c 100%);
  color: white;
  padding: 0 3px;
  border-radius: 3px;
  font-weight: 600;
}

.item-category {
  font-size: 10px;
  color: var(--text-muted);
  background: var(--bg-tertiary, #f5f5f5);
  padding: 1px 6px;
  border-radius: 4px;
  flex-shrink: 0;
}

.loading-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.loading-dots {
  display: flex;
  gap: 6px;
}

.loading-dots span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--primary);
  animation: dotPulse 1.4s ease-in-out infinite;
}

.loading-dots span:nth-child(2) { animation-delay: 0.2s; }
.loading-dots span:nth-child(3) { animation-delay: 0.4s; }

@keyframes dotPulse {
  0%, 80%, 100% { transform: scale(0.6); opacity: 0.4; }
  40% { transform: scale(1); opacity: 1; }
}

.loading-text {
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-muted);
}

.empty-suggestions {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px 16px;
  color: var(--text-muted);
}

.empty-suggestions .el-icon {
  margin-bottom: 8px;
  opacity: 0.4;
}

.empty-suggestions p {
  margin: 0;
  font-size: 13px;
}

.empty-hint {
  font-size: 11px !important;
  opacity: 0.6;
  margin-top: 4px !important;
}

@media (max-width: 768px) {
  .search-suggestions {
    position: fixed;
    top: auto;
    left: 16px;
    right: 16px;
    max-height: 60vh;
  }

  .search-input {
    height: 40px;
  }
}
</style>
