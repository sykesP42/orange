<template>
  <div class="tag-cloud-panel" :class="{ 'dark-mode': isDark }">
    <div class="panel-header">
      <h3 class="panel-title">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"/><line x1="7" y1="7" x2="7.01" y2="7"/></svg>
        标签云
      </h3>
      <div class="panel-actions">
        <button v-if="showRecommend" class="action-btn refresh-btn" @click="$emit('refresh')" title="刷新推荐">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>
        </button>
        <button class="action-btn toggle-btn" @click="toggleViewMode" :title="viewMode === 'cloud' ? '切换列表视图' : '切换云视图'">
          <svg v-if="viewMode === 'cloud'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M16 12l-4-4-4 4M12 16V8"/></svg>
        </button>
      </div>
    </div>

    <div v-if="tags.length === 0 && !loading" class="empty-state">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"/><line x1="7" y1="7" x2="7.01" y2="7"/></svg>
      <span>暂无标签</span>
    </div>

    <div v-if="loading" class="loading-spinner">
      <div class="spinner"></div>
    </div>

    <!-- Cloud View -->
    <div v-if="viewMode === 'cloud' && tags.length > 0 && !loading" class="cloud-view">
      <div
        v-for="(tag, index) in sortedTags"
        :key="tag.name"
        class="cloud-tag"
        :class="{ active: selectedTag === tag.name, recommended: tag.recommended }"
        :style="getTagStyle(tag, index)"
        @click="$emit('select', tag.name)"
        @mouseenter="hoverTag = tag.name"
        @mouseleave="hoverTag = null"
      >
        {{ tag.name }}
        <span class="tag-count">{{ tag.count }}</span>
        <span v-if="tag.recommended" class="recommend-badge">AI</span>
        <div v-if="hoverTag === tag.name" class="tag-tooltip">
          <strong>{{ tag.name }}</strong>
          <span>{{ tag.count }} 个博主</span>
          <span v-if="tag.color">颜色: {{ tag.color }}</span>
        </div>
      </div>
    </div>

    <!-- List View -->
    <div v-if="viewMode === 'list' && tags.length > 0 && !loading" class="list-view">
      <div
        v-for="tag in sortedTags"
        :key="tag.name"
        class="list-tag-item"
        :class="{ active: selectedTag === tag.name, recommended: tag.recommended }"
        @click="$emit('select', tag.name)"
      >
        <div class="list-tag-left">
          <span v-if="tag.color" class="color-dot" :style="{ backgroundColor: tag.color }"></span>
          <span v-else class="color-dot default"></span>
          <span class="list-tag-name">{{ tag.name }}</span>
          <span v-if="tag.recommended" class="recommend-badge small">推荐</span>
        </div>
        <div class="list-tag-right">
          <div class="tag-bar-wrapper">
            <div class="tag-bar" :style="{ width: getBarWidth(tag.count), backgroundColor: tag.color || 'var(--primary)' }"></div>
          </div>
          <span class="list-tag-count">{{ tag.count }}</span>
        </div>
      </div>
    </div>

    <div v-if="tags.length > 0" class="panel-footer">
      <span class="stats-text">{{ tags.length }} 个标签 · 共 {{ totalCount }} 次使用</span>
      <button v-if="selectedTag" class="clear-btn" @click="$emit('select', '')">清除筛选</button>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  tags: { type: Array, default: () => [] },
  selectedTag: { type: String, default: '' },
  loading: { type: Boolean, default: false },
  showRecommend: { type: Boolean, default: false },
  isDark: { type: Boolean, default: false }
})

defineEmits(['select', 'refresh'])

const viewMode = ref('cloud')
const hoverTag = ref(null)

const sortedTags = computed(() => {
  return [...props.tags].sort((a, b) => (b.count || 0) - (a.count || 0))
})

const totalCount = computed(() => {
  return props.tags.reduce((sum, t) => sum + (t.count || 0), 0)
})

const maxCount = computed(() => {
  return Math.max(...props.tags.map(t => t.count || 0), 1)
})

function getTagStyle(tag, index) {
  const count = tag.count || 1
  const ratio = count / maxCount.value
  const minSize = 11, maxSize = 26
  const fontSize = minSize + ratio * (maxSize - minSize)

  const hues = [22, 25, 28, 32, 35, 160, 200, 260]
  const hue = tag.color ? 0 : hues[index % hues.length]
  const saturation = tag.color ? 70 : 65 + ratio * 20
  const lightness = 45 + (1 - ratio) * 15

  let color
  if (tag.color) {
    color = tag.color
  } else {
    color = `hsl(${hue}, ${saturation}%, ${lightness}%)`
  }

  const opacity = 0.5 + ratio * 0.5

  return {
    fontSize: `${fontSize}px`,
    color,
    opacity,
    fontWeight: ratio > 0.6 ? '600' : ratio > 0.3 ? '500' : '400',
    transform: `scale(${0.95 + ratio * 0.1})`
  }
}

function getBarWidth(count) {
  const max = maxCount.value
  if (!max) return '0%'
  return `${Math.max((count / max) * 100, 8)}%`
}

function toggleViewMode() {
  viewMode.value = viewMode.value === 'cloud' ? 'list' : 'cloud'
}
</script>

<style scoped>
.tag-cloud-panel {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 14px;
  padding: 18px;
  transition: all 0.3s;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.panel-title svg {
  width: 18px;
  height: 18px;
  color: var(--primary);
}

.panel-actions {
  display: flex;
  gap: 4px;
}

.action-btn {
  width: 30px;
  height: 30px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  transition: all 0.2s;
}

.action-btn:hover {
  background: var(--bg-card-hover);
  color: var(--text-primary);
}

.action-btn svg {
  width: 16px;
  height: 16px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 36px 0;
  color: var(--text-muted);
  font-size: 13px;
}

.empty-state svg {
  width: 40px;
  height: 40px;
  opacity: 0.35;
}

.loading-spinner {
  display: flex;
  justify-content: center;
  padding: 24px 0;
}

.spinner {
  width: 28px;
  height: 28px;
  border: 2.5px solid var(--border-color);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.cloud-view {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  justify-content: center;
  padding: 8px 0;
}

.cloud-tag {
  position: relative;
  padding: 5px 12px;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  white-space: nowrap;
  line-height: 1.3;
  user-select: none;
  -webkit-tap-highlight-color: transparent;
}

.cloud-tag:hover {
  transform: scale(1.12) !important;
  z-index: 2;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

.cloud-tag.active {
  background: var(--primary);
  color: white !important;
  box-shadow: 0 4px 16px rgba(249, 115, 22, 0.3);
  opacity: 1 !important;
}

.cloud-tag.recommended::after {
  content: '';
  position: absolute;
  top: -2px;
  right: -2px;
  width: 8px;
  height: 8px;
  background: #22c55e;
  border-radius: 50%;
  border: 2px solid var(--bg-card);
}

.tag-count {
  font-size: 0.68em;
  opacity: 0.55;
  margin-left: 3px;
  font-weight: 400;
}

.cloud-tag.active .tag-count {
  opacity: 0.75;
  color: inherit;
}

.recommend-badge {
  position: absolute;
  top: -6px;
  right: -6px;
  background: linear-gradient(135deg, #22c55e, #16a34a);
  color: white;
  font-size: 9px;
  font-weight: 700;
  padding: 1px 4px;
  border-radius: 6px;
  line-height: 1.3;
}

.tag-tooltip {
  position: absolute;
  bottom: calc(100% + 8px);
  left: 50%;
  transform: translateX(-50%);
  background: var(--bg-elevated);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 8px 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  z-index: 10;
  display: flex;
  flex-direction: column;
  gap: 3px;
  font-size: 12px;
  color: var(--text-secondary);
  white-space: nowrap;
  pointer-events: none;
  animation: tooltip-in 0.15s ease-out;
}

@keyframes tooltip-in {
  from { opacity: 0; transform: translateX(-50%) translateY(4px); }
  to { opacity: 1; transform: translateX(-50%) translateY(0); }
}

.list-view {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 320px;
  overflow-y: auto;
}

.list-tag-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 9px 12px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
}

.list-tag-item:hover {
  background: var(--bg-card-hover);
}

.list-tag-item.active {
  background: rgba(249, 115, 22, 0.08);
  border: 1px solid rgba(249, 115, 22, 0.2);
}

.list-tag-left {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.color-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.color-dot.default {
  background: var(--primary);
}

.list-tag-name {
  font-size: 13.5px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.recommend-badge.small {
  font-size: 9px;
  padding: 1px 5px;
  background: linear-gradient(135deg, #22c55e, #16a34a);
  color: white;
  border-radius: 4px;
  font-weight: 600;
}

.list-tag-right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.tag-bar-wrapper {
  width: 80px;
  height: 6px;
  background: var(--border-color);
  border-radius: 3px;
  overflow: hidden;
}

.tag-bar {
  height: 100%;
  border-radius: 3px;
  transition: width 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.list-tag-count {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  min-width: 24px;
  text-align: right;
}

.panel-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.stats-text {
  font-size: 11.5px;
  color: var(--text-muted);
}

.clear-btn {
  font-size: 12px;
  color: var(--primary);
  background: none;
  border: none;
  cursor: pointer;
  font-weight: 500;
  padding: 3px 8px;
  border-radius: 5px;
  transition: all 0.2s;
}

.clear-btn:hover {
  background: rgba(249, 115, 22, 0.08);
}

@media (max-width: 768px) {
  .tag-cloud-panel {
    padding: 14px;
  }

  .cloud-view {
    gap: 6px;
  }

  .cloud-tag {
    padding: 4px 10px;
  }

  .list-view {
    max-height: 240px;
  }
}
</style>
