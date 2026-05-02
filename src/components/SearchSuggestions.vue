<template>
  <Teleport to="body">
    <div v-if="visible && hasSuggestions" class="search-suggestions" :class="{ 'is-pinyin': isPinyin }" :style="suggestionStyle">
    <div class="suggestions-header">
      <span class="suggestions-title">搜索建议</span>
      <span class="suggestions-count">{{ totalCount }} 个结果</span>
    </div>

    <div v-if="nicknames.length > 0" class="suggestion-group">
      <div class="group-label">
        <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
        博主
      </div>
      <div
        v-for="item in nicknames"
        :key="'n-' + item.value"
        class="suggestion-item"
        @click="$emit('select', item.value)"
        @mouseenter="hoveredIndex = item._idx"
        @mouseleave="hoveredIndex = -1"
        :class="{ 'is-hovered': hoveredIndex === item._idx }"
      >
        <span class="item-label" v-html="highlightText(item.label, query)"></span>
        <span class="item-pinyin" v-if="item.pinyin && isPinyin">{{ item.pinyin }}</span>
        <span class="item-count">{{ item.count }} 位博主</span>
      </div>
    </div>

    <div v-if="categories.length > 0" class="suggestion-group">
      <div class="group-label">
        <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
        分类
      </div>
      <div
        v-for="item in categories"
        :key="'c-' + item.value"
        class="suggestion-item"
        @click="$emit('select', item.value)"
        @mouseenter="hoveredIndex = item._idx"
        @mouseleave="hoveredIndex = -1"
        :class="{ 'is-hovered': hoveredIndex === item._idx }"
      >
        <span class="item-tag category-tag"></span>
        <span class="item-label" v-html="highlightText(item.label, query)"></span>
        <span class="item-pinyin" v-if="item.pinyin && isPinyin">{{ item.pinyin }}</span>
        <span class="item-count">{{ item.count }}</span>
      </div>
    </div>

    <div v-if="platforms.length > 0" class="suggestion-group">
      <div class="group-label">
        <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="3" width="20" height="14" rx="2"/><path d="M8 21h8m-4-4v4"/></svg>
        平台
      </div>
      <div
        v-for="item in platforms"
        :key="'p-' + item.value"
        class="suggestion-item"
        @click="$emit('select', item.value)"
        @mouseenter="hoveredIndex = item._idx"
        @mouseleave="hoveredIndex = -1"
        :class="{ 'is-hovered': hoveredIndex === item._idx }"
      >
        <span class="item-tag platform-tag"></span>
        <span class="item-label" v-html="highlightText(item.label, query)"></span>
        <span class="item-count">{{ item.count }}</span>
      </div>
    </div>

    <div v-if="!hasSuggestions && query" class="suggestions-empty">
      <span>无匹配结果</span>
    </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { bloggerSuggestionsAPI } from '../api'
import { highlightMatch, isPinyinInput } from '../utils/pinyin-search'

const props = defineProps({
  keyword: { type: String, default: '' },
  visible: { type: Boolean, default: false }
})

const emit = defineEmits(['select'])

const nicknames = ref([])
const categories = ref([])
const platforms = ref([])
const hoveredIndex = ref(-1)
let fetchTimer = null
const anchorRect = ref({ top: 0, left: 0, width: 0 })

const query = computed(() => props.keyword?.trim() || '')
const isPinyin = computed(() => query.value ? isPinyinInput(query.value) : false)

const totalCount = computed(() => {
  return nicknames.value.length + categories.value.length + platforms.value.length
})

const hasSuggestions = computed(() => totalCount.value > 0)

const suggestionStyle = computed(() => ({
  position: 'fixed',
  top: `${anchorRect.value.bottom + 4}px`,
  left: `${anchorRect.value.left}px`,
  width: `${anchorRect.value.width}px`,
}))

function updateAnchorPosition() {
  const searchInput = document.querySelector('.search-box input')
  if (searchInput) {
    const rect = searchInput.getBoundingClientRect()
    anchorRect.value = {
      top: rect.top,
      left: rect.left,
      bottom: rect.bottom,
      width: rect.width,
    }
  }
}

onMounted(() => {
  updateAnchorPosition()
  window.addEventListener('resize', updateAnchorPosition)
  window.addEventListener('scroll', updateAnchorPosition, true)
})

function highlightText(text, kw) {
  if (!text || !kw) return text
  return highlightMatch(text, kw)
}

async function fetchSuggestions(kw) {
  if (!kw || kw.length < 1) {
    nicknames.value = []
    categories.value = []
    platforms.value = []
    return
  }

  try {
    const res = await bloggerSuggestionsAPI(kw)
    if (res.code === 200 && res.data) {
      let idx = 0
      const nicks = (res.data.nicknames || []).map(item => ({ ...item, _idx: idx++ }))
      const cats = (res.data.categories || []).map(item => ({ ...item, _idx: idx++ }))
      const plats = (res.data.platforms || []).map(item => ({ ...item, _idx: idx++ }))
      nicknames.value = nicks
      categories.value = cats
      platforms.value = plats
    }
  } catch (e) {
    nicknames.value = []
    categories.value = []
    platforms.value = []
  }
}

watch(() => props.keyword, (val) => {
  clearTimeout(fetchTimer)
  hoveredIndex.value = -1
  if (!val || val.trim().length < 1) {
    nicknames.value = []
    categories.value = []
    platforms.value = []
    return
  }
  fetchTimer = setTimeout(() => fetchSuggestions(val.trim()), 150)
})

watch(() => props.visible, (val) => {
  if (val) {
    nextTick(updateAnchorPosition)
  } else {
    nicknames.value = []
    categories.value = []
    platforms.value = []
  }
})
</script>

<style scoped>
.search-suggestions {
  background: var(--bg-card, #fff);
  border: 1px solid var(--border-color, #e4e7ed);
  border-radius: 12px;
  box-shadow:
    0 8px 30px rgba(0, 0, 0, 0.12),
    0 2px 8px rgba(0, 0, 0, 0.06),
    0 0 1px rgba(0, 0, 0, 0.04);
  z-index: 9999;
  max-height: 360px;
  overflow-y: auto;
  animation: suggestDrop 0.18s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes suggestDrop {
  from { opacity: 0; transform: translateY(-6px); }
  to { opacity: 1; transform: translateY(0); }
}

.suggestions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 14px;
  border-bottom: 1px solid var(--border-light, #f0f0f0);
  font-size: 11px;
  color: var(--text-secondary, #909399);
}

.suggestions-title { font-weight: 600; }
.suggestions-count { color: var(--text-placeholder, #c0c4cc); }

.suggestion-group {
  padding: 4px 0;
}

.suggestion-group + .suggestion-group {
  border-top: 1px solid var(--border-light, #f0f0f0);
}

.group-label {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px 4px;
  font-size: 11px;
  font-weight: 600;
  color: var(--text-secondary, #909399);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.group-label svg { opacity: 0.6; }

.suggestion-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 14px;
  cursor: pointer;
  transition: all 0.12s ease;
  border-radius: 6px;
  margin: 1px 6px;
  font-size: 13px;
}

.suggestion-item:hover,
.suggestion-item.is-hovered {
  background: var(--bg-hover, #fff8f0);
}

.item-tag {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.category-tag { background: linear-gradient(135deg, #ff9500, #ff5e00); }
.platform-tag { background: linear-gradient(135deg, #007aff, #5856d6); }

.item-label {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--text-primary, #303133);
}

.item-label :deep(mark) {
  background: linear-gradient(120deg, #ff9500 0%, #ff6b00 100%);
  color: var(--color-on-brand);
  padding: 0 3px;
  border-radius: 3px;
  font-weight: 600;
}

.item-pinyin {
  font-size: 11px;
  color: var(--text-placeholder, #c0c4cc);
  font-family: 'SF Mono', Monaco, Consolas, monospace;
}

.item-count {
  font-size: 11px;
  color: var(--text-placeholder, #c0c4cc);
  background: var(--bg-page, #f5f7fa);
  padding: 1px 7px;
  border-radius: 8px;
  flex-shrink: 0;
}

.suggestions-empty {
  padding: 20px;
  text-align: center;
  color: var(--text-placeholder, #c0c4cc);
  font-size: 13px;
}

.search-suggestions.is-pinyin .suggestions-header::after {
  content: '🔤 拼音模式';
  display: inline-block;
  margin-left: 8px;
  font-size: 10px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: var(--color-on-brand);
  padding: 1px 6px;
  border-radius: 4px;
}
</style>
