<template>
  <div class="enhanced-search" :class="{ 'is-focus': isFocus, 'has-value': modelValue }">
    <div class="search-input-wrapper">
      <el-icon class="search-icon"><Search /></el-icon>
      <input
        v-model="inputValue"
        type="text"
        :placeholder="placeholder"
        class="search-input"
        @focus="handleFocus"
        @blur="handleBlur"
        @input="handleInput"
        @keydown.enter="handleSearch"
        @keydown.esc="handleClear"
      />
      <el-icon 
        v-if="modelValue" 
        class="clear-icon" 
        @click="handleClear"
      >
        <CircleClose />
      </el-icon>
    </div>
    
    <!-- 搜索历史和热门推荐 -->
    <transition name="el-fade-in">
      <div v-if="isFocus && showSuggestions" class="search-suggestions">
        <!-- 搜索历史 -->
        <div v-if="showHistory && searchHistory.length > 0" class="suggestion-section">
          <div class="suggestion-header">
            <span class="section-title">搜索历史</span>
            <el-icon class="clear-btn" @click="clearHistory"><Delete /></el-icon>
          </div>
          <div class="suggestion-tags">
            <el-tag
              v-for="(item, index) in searchHistory"
              :key="index"
              size="small"
              class="history-tag"
              @click="selectSuggestion(item)"
            >
              {{ item }}
              <el-icon class="tag-close" @click.stop="removeHistory(index)"><Close /></el-icon>
            </el-tag>
          </div>
        </div>
        
        <!-- 热门搜索 -->
        <div v-if="showHotTags && hotTags.length > 0" class="suggestion-section">
          <div class="suggestion-header">
            <span class="section-title">热门搜索</span>
            <el-icon class="refresh-btn" @click="refreshHotTags"><Refresh /></el-icon>
          </div>
          <div class="suggestion-tags">
            <el-tag
              v-for="(tag, index) in hotTags"
              :key="index"
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
        
        <!-- 搜索建议 -->
        <div v-if="showSuggestions && suggestions.length > 0" class="suggestion-section">
          <div class="suggestion-header">
            <span class="section-title">搜索建议</span>
          </div>
          <ul class="suggestion-list">
            <li
              v-for="(item, index) in suggestions"
              :key="index"
              class="suggestion-item"
              :class="{ 'is-active': activeIndex === index }"
              @mouseenter="activeIndex = index"
              @click="selectSuggestion(item)"
            >
              <el-icon><Search /></el-icon>
              <span class="item-text" v-html="highlightKeyword(item)"></span>
            </li>
          </ul>
        </div>
        
        <!-- 空状态 -->
        <div v-if="isFocus && !showHistory && !showHotTags && !suggestions.length" class="empty-suggestions">
          <el-icon :size="48"><Search /></el-icon>
          <p>输入关键词开始搜索</p>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { Search, CircleClose, Delete, Close, Refresh } from '@element-plus/icons-vue'
import { useDebounce, SearchHistory, HotSearchTags, SearchSuggestions } from '@/utils/search'

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
    default: 300
  },
  maxSuggestions: {
    type: Number,
    default: 5
  }
})

const emit = defineEmits(['update:modelValue', 'search'])

const inputValue = ref(props.modelValue)
const isFocus = ref(false)
const activeIndex = ref(-1)
const suggestions = ref([])

// 搜索工具实例
const searchHistory = new SearchHistory()
const hotTags = new HotSearchTags()
const searchSuggestions = new SearchSuggestions()

// 计算属性
const showSuggestions = computed(() => {
  return inputValue.value && inputValue.value.trim().length > 0
})

const searchHistoryList = ref(searchHistory.get())
const hotTagsList = ref(hotTags.getList())

// 监听输入
const handleInput = useDebounce((event) => {
  const value = event.target.value
  emit('update:modelValue', value)
  
  // 获取搜索建议
  fetchSuggestions(value)
}, props.debounce)

// 获取搜索建议（可以替换为 API 调用）
const fetchSuggestions = async (query) => {
  if (!query || query.trim().length === 0) {
    suggestions.value = []
    return
  }
  
  // 尝试从缓存获取
  const cached = searchSuggestions.get(query)
  if (cached) {
    suggestions.value = cached
    return
  }
  
  // TODO: 调用 API 获取搜索建议
  // const res = await searchSuggestAPI(query)
  // suggestions.value = res.data
  
  // 临时模拟数据
  searchSuggestions.set(query, [])
  suggestions.value = []
}

// 聚焦
const handleFocus = () => {
  isFocus.value = true
  searchHistoryList.value = searchHistory.get()
  hotTagsList.value = hotTags.getList()
}

// 失焦
const handleBlur = () => {
  setTimeout(() => {
    isFocus.value = false
    activeIndex.value = -1
  }, 200)
}

// 搜索
const handleSearch = () => {
  const query = inputValue.value.trim()
  if (!query) return
  
  // 添加到搜索历史
  searchHistory.add(query)
  searchHistoryList.value = searchHistory.get()
  
  // 更新热门搜索
  hotTags.update(query)
  hotTagsList.value = hotTags.getList()
  
  // 触发搜索事件
  emit('search', query)
  
  // 清空建议
  suggestions.value = []
}

// 清空
const handleClear = () => {
  inputValue.value = ''
  emit('update:modelValue', '')
  suggestions.value = []
  activeIndex.value = -1
}

// 选择建议
const selectSuggestion = (value) => {
  inputValue.value = value
  emit('update:modelValue', value)
  handleSearch()
}

// 删除历史记录
const removeHistory = (index) => {
  const item = searchHistoryList.value[index]
  searchHistory.remove(item)
  searchHistoryList.value = searchHistory.get()
}

// 清空历史
const clearHistory = () => {
  searchHistory.clear()
  searchHistoryList.value = []
}

// 刷新热门标签
const refreshHotTags = () => {
  hotTagsList.value = hotTags.getList()
}

// 高亮关键词
const highlightKeyword = (text) => {
  if (!inputValue.value) return text
  const regex = new RegExp(`(${inputValue.value})`, 'gi')
  return text.replace(regex, '<span class="highlight">$1</span>')
}

// 监听外部值变化
watch(() => props.modelValue, (val) => {
  inputValue.value = val
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

/* 搜索建议面板 */
.search-suggestions {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  right: 0;
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  border: 1px solid var(--border-color);
  padding: 16px;
  z-index: 100;
  max-height: 400px;
  overflow-y: auto;
}

.suggestion-section {
  margin-bottom: 16px;
}

.suggestion-section:last-child {
  margin-bottom: 0;
}

.suggestion-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-light);
}

.section-title {
  font-size: 13px;
  color: var(--text-secondary);
  font-weight: 500;
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
  gap: 8px;
}

.history-tag, .hot-tag {
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 4px;
}

.history-tag:hover, .hot-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.tag-close {
  margin-left: 4px;
  font-size: 12px;
}

.tag-close:hover {
  color: var(--danger);
}

.hot-index {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  background: var(--danger);
  color: white;
  font-size: 10px;
  font-weight: bold;
  border-radius: 50%;
  margin-right: 4px;
}

/* 搜索建议列表 */
.suggestion-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.suggestion-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.suggestion-item:hover,
.suggestion-item.is-active {
  background: var(--bg-tertiary);
}

.suggestion-item .el-icon {
  color: var(--text-muted);
}

.item-text {
  flex: 1;
  color: var(--text-primary);
}

.item-text :deep(.highlight) {
  color: var(--primary);
  font-weight: 600;
}

/* 空状态 */
.empty-suggestions {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: var(--text-muted);
}

.empty-suggestions .el-icon {
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-suggestions p {
  margin: 0;
  font-size: 14px;
}

/* 移动端优化 */
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
