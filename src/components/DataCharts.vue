<template>
  <div class="data-charts" v-if="bloggers.length > 0">
    <div class="charts-header" @click="expanded = !expanded">
      <h3>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
          <path d="M18 20V10M12 20V4M6 20v-6"/>
        </svg>
        数据看板
      </h3>
      <span class="charts-subtitle">基于当前筛选结果的 {{ bloggers.length }} 条数据</span>
      <button class="collapse-toggle" :class="{ expanded }" :aria-label="expanded ? '折叠' : '展开'">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <polyline points="6,9 12,15 18,9"/>
        </svg>
      </button>
    </div>

    <transition name="charts-collapse">
    <div v-show="expanded" class="charts-grid">
      <div class="chart-card status-chart-card">
        <div class="chart-title">状态分布</div>
        <div ref="statusChartRef" class="chart-container"></div>
      </div>

      <div class="chart-card category-chart-card">
        <div class="chart-title">分类分布</div>
        <div ref="categoryChartRef" class="chart-container"></div>
      </div>

      <div class="chart-card platform-chart-card">
        <div class="chart-title">平台分布</div>
        <div ref="platformChartRef" class="chart-container"></div>
      </div>

      <div class="chart-card trend-chart-card full-width">
        <div class="chart-title">近期趋势</div>
        <div ref="trendChartRef" class="chart-container"></div>
      </div>
    </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'

const props = defineProps({
  bloggers: { type: Array, default: () => [] },
  isDark: { type: Boolean, default: false }
})

const expanded = ref(false)

const statusChartRef = ref(null)
const categoryChartRef = ref(null)
const platformChartRef = ref(null)
const trendChartRef = ref(null)

let statusChart = null
let categoryChart = null
let platformChart = null
let trendChart = null

const statusColors = {
  '初次联系': '#6366f1',
  '洽谈中': '#f59e0b',
  '已合作': '#22c55e',
  '已拒绝': '#ef4444',
  '暂停合作': '#94a3b8'
}

const categoryColors = [
  '#f97316', '#8b5cf6', '#06b6d4', '#ec4899',
  '#eab308', '#14b8a6', '#f43f5e', '#84cc16',
  '#a855f7', '#0ea5e9', '#ef4444', '#22c55e'
]

function aggregateStatus(data) {
  const map = {}
  data.forEach(b => {
    const s = b.status || '初次联系'
    map[s] = (map[s] || 0) + 1
  })
  return Object.entries(map).map(([name, value]) => ({ name, value }))
}

function aggregateCategory(data) {
  const map = {}
  data.forEach(b => {
    if (b.category) {
      map[b.category] = (map[b.category] || 0) + 1
    }
  })
  return Object.entries(map)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 8)
    .map(([name, value], i) => ({ name, value, itemStyle: { color: categoryColors[i % categoryColors.length] } }))
}

function aggregatePlatform(data) {
  const map = {}
  data.forEach(b => {
    if (b.platform) {
      map[b.platform] = (map[b.platform] || 0) + 1
    }
  })
  return Object.entries(map)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 6)
    .map(([name, value]) => ({ name, value }))
}

function aggregateTrend(data) {
  const map = {}
  data.forEach(b => {
    if (b.create_time) {
      const d = new Date(b.create_time)
      const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
      map[key] = (map[key] || 0) + 1
    }
  })
  return Object.entries(map)
    .sort((a, b) => a[0].localeCompare(b[0]))
    .slice(-6)
    .map(([name, value]) => ({ name, value }))
}

function getThemeColors() {
  return {
    text: props.isDark ? '#c4c9d4' : '#374151',
    textSecondary: props.isDark ? '#8b95a5' : '#9ca3af',
    border: props.isDark ? '#2d3748' : '#e5e7eb',
    bg: props.isDark ? '#1a1f2e' : '#ffffff',
    gridLine: props.isDark ? '#2a3441' : '#f0f0f0'
  }
}

function initStatusChart() {
  if (!statusChartRef.value) return
  try { statusChart?.dispose() } catch {}
  statusChart = echarts.init(statusChartRef.value)
  updateStatusChart()
}

function initCategoryChart() {
  if (!categoryChartRef.value) return
  try { categoryChart?.dispose() } catch {}
  categoryChart = echarts.init(categoryChartRef.value)
  updateCategoryChart()
}

function initPlatformChart() {
  if (!platformChartRef.value) return
  try { platformChart?.dispose() } catch {}
  platformChart = echarts.init(platformChartRef.value)
  updatePlatformChart()
}

function initTrendChart() {
  if (!trendChartRef.value) return
  try { trendChart?.dispose() } catch {}
  trendChart = echarts.init(trendChartRef.value)
  updateTrendChart()
}

function updateStatusChart() {
  if (!statusChart) return
  try {
    const c = getThemeColors()
    const data = aggregateStatus(props.bloggers)
    statusChart.setOption({
      tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
      legend: { bottom: 0, textStyle: { color: c.textSecondary, fontSize: 11 }, itemWidth: 12, itemHeight: 12 },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['50%', '45%'],
        avoidLabelOverlap: true,
        itemStyle: { borderRadius: 6, borderColor: c.bg, borderWidth: 2 },
        label: { show: false },
        emphasis: { label: { show: true, fontSize: 13, fontWeight: 'bold' } },
        data: data.map(d => ({ ...d, itemStyle: { color: statusColors[d.name] || '#94a3b8' } }))
      }],
      color: Object.values(statusColors)
    }, true)
  } catch (e) { console.warn('statusChart error:', e) }
}

function updateCategoryChart() {
  if (!categoryChart) return
  try {
    const c = getThemeColors()
    const data = aggregateCategory(props.bloggers)
    categoryChart.setOption({
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      grid: { left: 3, right: 10, top: 10, bottom: 24, containLabel: true },
      xAxis: { type: 'value', show: false },
      yAxis: {
        type: 'category',
        data: data.map(d => d.name).reverse(),
        axisLine: { lineStyle: { color: c.border } },
        axisTick: { show: false },
        axisLabel: { color: c.text, fontSize: 11 }
      },
      series: [{
        type: 'bar',
        data: data.map(d => d.value).reverse(),
        barWidth: 16,
        itemStyle: { borderRadius: [0, 4, 4, 0] },
        label: { show: true, position: 'right', color: c.textSecondary, fontSize: 11 }
      }]
    }, true)
  } catch (e) { console.warn('categoryChart error:', e) }
}

function updatePlatformChart() {
  if (!platformChart) return
  try {
    const c = getThemeColors()
    const data = aggregatePlatform(props.bloggers)
    platformChart.setOption({
      tooltip: { trigger: 'item' },
      grid: { left: 10, right: 10, top: 10, bottom: 24, containLabel: true },
      series: [{
        type: 'pie',
        radius: ['25%', '60%'],
        center: ['50%', '48%'],
        roseType: 'area',
        itemStyle: { borderRadius: 5, borderColor: c.bg, borderWidth: 2 },
        label: { color: c.text, fontSize: 11 },
        data: data.map((d, i) => ({
          ...d,
          itemStyle: { color: categoryColors[i % categoryColors.length] }
        }))
      }]
    }, true)
  } catch (e) { console.warn('platformChart error:', e) }
}

function updateTrendChart() {
  if (!trendChart) return
  try {
    const c = getThemeColors()
    const data = aggregateTrend(props.bloggers)
    trendChart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: 40, right: 16, top: 12, bottom: 28, containLabel: true },
      xAxis: {
        type: 'category',
        data: data.map(d => d.name),
        boundaryGap: false,
        axisLine: { lineStyle: { color: c.border } },
        axisTick: { show: false },
        axisLabel: { color: c.textSecondary, fontSize: 11 }
      },
      yAxis: {
        type: 'value',
        splitLine: { lineStyle: { color: c.gridLine, type: 'dashed' } },
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: { color: c.textSecondary, fontSize: 11 }
      },
      series: [{
        type: 'line',
        data: data.map(d => d.value),
        smooth: true,
        symbol: 'circle',
        symbolSize: 7,
        lineStyle: { width: 2.5, color: '#f97316' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(249,115,22,0.25)' },
              { offset: 1, color: 'rgba(249,115,22,0.02)' }
            ]
          }
        },
        itemStyle: { color: '#f97316', borderColor: '#fff', borderWidth: 2 }
      }]
    }, true)
  } catch (e) { console.warn('trendChart error:', e) }
}

function updateAllCharts() {
  nextTick(() => {
    updateStatusChart()
    updateCategoryChart()
    updatePlatformChart()
    updateTrendChart()
  })
}

function handleResize() {
  statusChart?.resize()
  categoryChart?.resize()
  platformChart?.resize()
  trendChart?.resize()
}

watch(() => props.bloggers, () => { if (expanded.value) updateAllCharts() }, { deep: true })
watch(() => props.isDark, () => { if (expanded.value) updateAllCharts() })
watch(expanded, (val) => {
  if (val) {
    nextTick(() => setTimeout(initAllCharts, 80))
  } else {
    statusChart?.dispose(); statusChart = null
    categoryChart?.dispose(); categoryChart = null
    platformChart?.dispose(); platformChart = null
    trendChart?.dispose(); trendChart = null
  }
})

function initAllCharts() {
  initStatusChart()
  initCategoryChart()
  initPlatformChart()
  initTrendChart()
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  statusChart?.dispose()
  categoryChart?.dispose()
  platformChart?.dispose()
  trendChart?.dispose()
})
</script>

<style scoped>
.data-charts {
  margin-bottom: 28px;
  animation: fadeInUp 0.4s ease-out;
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(16px); }
  to { opacity: 1; transform: translateY(0); }
}

.charts-header {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 16px;
  cursor: pointer;
  user-select: none;
}

.charts-header:hover .collapse-toggle {
  opacity: 1;
}

.charts-header h3 {
  font-size: 17px;
  font-weight: 600;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
}

.charts-header h3 svg {
  width: 18px;
  height: 18px;
  color: var(--primary);
}

.charts-subtitle {
  font-size: 13px;
  color: var(--text-muted);
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.chart-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 14px;
  padding: 18px;
  transition: all 0.3s ease;
}

.chart-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
  border-color: rgba(249, 115, 22, 0.15);
}

.full-width {
  grid-column: 1 / -1;
}

.chart-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 10px;
}

.chart-container {
  height: 220px;
  width: 100%;
}

.dark .chart-card {
  background: var(--bg-secondary);
  border-color: rgba(255, 255, 255, 0.06);
}

.collapse-toggle {
  width: 28px;
  height: 28px;
  border: none;
  background: var(--bg-secondary);
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  margin-left: auto;
  flex-shrink: 0;
  opacity: 0.5;
  transition: all 0.3s ease;
}

.collapse-toggle svg {
  width: 14px;
  height: 14px;
  transition: transform 0.3s ease;
}

.collapse-toggle.expanded svg {
  transform: rotate(180deg);
}

.collapse-toggle:hover {
  opacity: 1;
  background: var(--bg-hover);
}

.charts-collapse-enter-active {
  transition: all 0.35s ease;
  overflow: hidden;
}
.charts-collapse-leave-active {
  transition: all 0.25s ease;
  overflow: hidden;
}
.charts-collapse-enter-from,
.charts-collapse-leave-to {
  opacity: 0;
  max-height: 0;
  margin-top: 0;
}
.charts-collapse-enter-to,
.charts-collapse-leave-from {
  max-height: 800px;
}

@media (max-width: 768px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }

  .full-width {
    grid-column: auto;
  }

  .chart-container {
    height: 200px;
  }

  .charts-header {
    flex-direction: column;
    gap: 4px;
  }

  .charts-header .collapse-toggle {
    position: absolute;
    right: 12px;
    top: 10px;
  }
}
</style>
