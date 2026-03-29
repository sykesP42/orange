<template>
  <div class="pull-refresh" @touchstart="onTouchStart" @touchmove="onTouchMove" @touchend="onTouchEnd">
    <!-- 下拉刷新提示 -->
    <div class="refresh-header" :style="{ height: headerHeight }">
      <div class="refresh-content" :style="{ transform: pullDistance > 0 ? `translateY(${pullDistance}px)` : 'translateY(0)' }">
        <div v-if="refreshing" class="refreshing-spinner">
          <svg class="spinner" viewBox="0 0 50 50">
            <circle cx="25" cy="25" r="20" fill="none" stroke-width="5"></circle>
          </svg>
        </div>
        <div v-else-if="pullDistance > pullDownRefreshThreshold" class="refresh-text">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="18,15 12,9 6,15"/>
          </svg>
          <span>释放刷新</span>
        </div>
        <div v-else-if="pullDistance > 0" class="refresh-text">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="6,9 12,15 18,9"/>
          </svg>
          <span>下拉刷新</span>
        </div>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="refresh-body">
      <slot></slot>
    </div>

    <!-- 上拉加载提示 -->
    <div v-if="loading && hasMore" class="loading-footer">
      <div class="loading-spinner">
        <svg class="spinner" viewBox="0 0 50 50">
          <circle cx="25" cy="25" r="20" fill="none" stroke-width="5"></circle>
        </svg>
      </div>
      <span>加载中...</span>
    </div>
    <div v-else-if="!hasMore" class="no-more">
      <span>没有更多了</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  refreshing: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  hasMore: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['refresh', 'loadMore'])

const pullDistance = ref(0)
const startY = ref(0)
const isScrolledToTop = ref(false)
const headerHeight = ref('0px')
const pullDownRefreshThreshold = 80

const containerRef = ref(null)

const onTouchStart = (event) => {
  startY.value = event.touches[0].clientY
  const container = event.currentTarget.querySelector('.refresh-body')
  isScrolledToTop.value = container && container.scrollTop <= 0
}

const onTouchMove = (event) => {
  if (!isScrolledToTop.value || props.refreshing) return

  const currentY = event.touches[0].clientY
  const diff = currentY - startY.value

  if (diff > 0) {
    event.preventDefault()
    // 添加阻力效果
    pullDistance.value = diff * 0.5
    headerHeight.value = `${pullDistance.value}px`
  }
}

const onTouchEnd = () => {
  if (!isScrolledToTop.value || props.refreshing) return

  if (pullDistance.value > pullDownRefreshThreshold) {
    emit('refresh')
  }
  pullDistance.value = 0
  headerHeight.value = '0px'
}

// 监听 refreshing 状态，刷新完成后重置
watch(() => props.refreshing, (newVal) => {
  if (!newVal) {
    pullDistance.value = 0
    headerHeight.value = '0px'
  }
})

const checkScroll = (event) => {
  const container = event.currentTarget
  isScrolledToTop.value = container.scrollTop <= 0

  // 检测是否滚动到底部
  const { scrollTop, scrollHeight, clientHeight } = container
  if (scrollHeight - scrollTop - clientHeight < 50 && !props.loading && props.hasMore) {
    emit('loadMore')
  }
}
</script>

<style scoped>
.pull-refresh {
  height: 100%;
  overflow: hidden;
}

.refresh-header {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: height 0.3s ease;
  overflow: hidden;
}

.refresh-content {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.3s ease;
}

.refresh-text {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #999;
  font-size: 14px;
}

.refresh-text svg {
  width: 20px;
  height: 20px;
  animation: bounce 0.5s ease infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}

.refreshing-spinner svg {
  width: 24px;
  height: 24px;
  animation: rotate 1s linear infinite;
}

.refreshing-spinner circle {
  stroke: var(--primary);
  stroke-dasharray: 90, 150;
  stroke-dashoffset: 0;
  stroke-linecap: round;
}

@keyframes rotate {
  to { transform: rotate(360deg); }
}

.refresh-body {
  height: calc(100% - var(--header-height, 0px));
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.loading-footer,
.no-more {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 20px;
  color: #999;
  font-size: 14px;
}

.loading-spinner svg {
  width: 20px;
  height: 20px;
  animation: rotate 1s linear infinite;
}

.loading-spinner circle {
  stroke: var(--primary);
  stroke-dasharray: 90, 150;
  stroke-dashoffset: 0;
  stroke-linecap: round;
}
</style>
