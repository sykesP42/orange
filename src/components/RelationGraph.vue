<template>
  <div class="relation-graph" :class="{ dark: isDark }">
    <div class="graph-controls">
      <div class="graph-legend">
        <button
          v-for="item in legendItems"
          :key="item.key"
          class="legend-item"
          :class="{ active: visibleCategories[item.key], dimmed: !visibleCategories[item.key] }"
          @click="toggleCategory(item.key)"
        >
          <span class="legend-dot" :style="{ background: item.color }"></span>
          <span class="legend-label">{{ item.label }}</span>
          <span class="legend-count">{{ item.count }}</span>
        </button>
      </div>
      <div class="graph-actions">
        <div class="graph-search">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
          <input v-model="searchText" placeholder="搜索节点..." />
          <button v-if="searchText" class="search-clear" @click="searchText = ''">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </div>
        <div class="graph-zoom">
          <button @click="zoomIn" title="放大">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          </button>
          <span class="zoom-level">{{ Math.round(currentZoom * 100) }}%</span>
          <button @click="zoomOut" title="缩小">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="5" y1="12" x2="19" y2="12"/></svg>
          </button>
          <button @click="resetZoom" title="重置" class="reset-btn">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/></svg>
          </button>
        </div>
      </div>
    </div>

    <div class="graph-canvas-wrapper">
      <div id="relation-graph-chart" class="graph-chart"></div>
      <div v-if="bloggers.length === 0" class="graph-empty">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><path d="M8 12h8"/></svg>
        <p>暂无关系数据</p>
        <span>录入博主后将自动生成关系图谱</span>
      </div>
      <div v-if="hoveredNode" class="node-tooltip" :style="tooltipStyle">
        <div class="tooltip-header">
          <span class="tooltip-dot" :style="{ background: getNodeColor(hoveredNode.category) }"></span>
          <span class="tooltip-name">{{ hoveredNode.name }}</span>
        </div>
        <div class="tooltip-meta">
          <span class="tooltip-type">{{ getCategoryLabel(hoveredNode.category) }}</span>
          <span v-if="hoveredNode.connections" class="tooltip-connections">{{ hoveredNode.connections }} 个关联</span>
        </div>
      </div>
      <div class="graph-stats">
        <span>{{ filteredNodes.length }} 节点</span>
        <span>{{ filteredLinks.length }} 关系</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import { useThemeStore } from '../stores/theme'

const props = defineProps({
  bloggers: { type: Array, default: () => [] }
})

const themeStore = useThemeStore()
const isDark = computed(() => themeStore.isDark)

let chart = null
let currentZoom = 1
const searchText = ref('')
const hoveredNode = ref(null)
const tooltipStyle = ref({})

const visibleCategories = ref({
  blogger: true,
  category: true,
  tag: true,
  platform: true
})

const categoryIndexMap = { blogger: 0, category: 1, tag: 2, platform: 3 }

const COLORS = {
  blogger: { light: '#f97316', dark: '#fb923c', glow: 'rgba(249,115,22,0.3)' },
  category: { light: '#3b82f6', dark: '#60a5fa', glow: 'rgba(59,130,246,0.3)' },
  tag: { light: '#8b5cf6', dark: '#a78bfa', glow: 'rgba(139,92,246,0.3)' },
  platform: { light: '#10b981', dark: '#34d399', glow: 'rgba(16,185,129,0.3)' }
}

const LINK_COLORS = {
  category: { light: 'rgba(59,130,246,0.2)', dark: 'rgba(96,165,250,0.25)' },
  tag: { light: 'rgba(139,92,246,0.15)', dark: 'rgba(167,139,250,0.2)' },
  platform: { light: 'rgba(16,185,129,0.2)', dark: 'rgba(52,211,153,0.25)' }
}

const getNodeColor = (catIdx) => {
  const keys = ['blogger', 'category', 'tag', 'platform']
  const key = keys[catIdx] || 'blogger'
  return isDark.value ? COLORS[key].dark : COLORS[key].light
}

const getCategoryLabel = (catIdx) => {
  return ['博主', '分类', '标签', '平台'][catIdx] || '未知'
}

const buildGraphData = () => {
  const nodes = []
  const links = []
  const categories = [
    { name: '博主', itemStyle: { color: isDark.value ? COLORS.blogger.dark : COLORS.blogger.light } },
    { name: '分类', itemStyle: { color: isDark.value ? COLORS.category.dark : COLORS.category.light } },
    { name: '标签', itemStyle: { color: isDark.value ? COLORS.tag.dark : COLORS.tag.light } },
    { name: '平台', itemStyle: { color: isDark.value ? COLORS.platform.dark : COLORS.platform.light } }
  ]

  const categorySet = new Set()
  const tagSet = new Set()
  const platformSet = new Set()
  const connectionCount = {}

  const addConnection = (id) => {
    connectionCount[id] = (connectionCount[id] || 0) + 1
  }

  props.bloggers.forEach(b => {
    const bloggerId = `blogger_${b.id}`
    addConnection(bloggerId)
    if (b.category) addConnection(`cat_${b.category}`)
    if (b.platform) addConnection(`plat_${b.platform}`)
    if (b.tags && Array.isArray(b.tags)) {
      b.tags.forEach(tag => addConnection(`tag_${tag}`))
    }
  })

  props.bloggers.forEach(b => {
    const bloggerId = `blogger_${b.id}`
    const connCount = connectionCount[bloggerId] || 1
    const size = Math.min(20 + connCount * 4, 50)

    nodes.push({
      id: bloggerId,
      name: b.nickname || b.real_name || '未知',
      category: 0,
      symbolSize: size,
      value: connCount,
      itemStyle: {
        color: isDark.value ? COLORS.blogger.dark : COLORS.blogger.light,
        shadowBlur: size > 30 ? 12 : 6,
        shadowColor: COLORS.blogger.glow
      },
      label: {
        show: true,
        fontSize: Math.min(10 + connCount, 14),
        color: isDark.value ? '#f9fafb' : '#1f2937',
        fontWeight: connCount > 3 ? 'bold' : 'normal'
      },
      connections: connCount
    })

    if (b.category && !categorySet.has(b.category)) {
      categorySet.add(b.category)
      const catId = `cat_${b.category}`
      const catConn = connectionCount[catId] || 1
      nodes.push({
        id: catId,
        name: b.category,
        category: 1,
        symbolSize: Math.min(16 + catConn * 3, 38),
        value: catConn,
        itemStyle: {
          color: isDark.value ? COLORS.category.dark : COLORS.category.light,
          shadowBlur: catConn > 2 ? 10 : 4,
          shadowColor: COLORS.category.glow
        },
        label: {
          show: true,
          fontSize: Math.min(9 + catConn, 13),
          color: isDark.value ? '#f9fafb' : '#1f2937'
        },
        connections: catConn
      })
    }
    if (b.category) {
      links.push({
        source: bloggerId,
        target: `cat_${b.category}`,
        lineStyle: {
          color: isDark.value ? 'rgba(96,165,250,0.25)' : 'rgba(59,130,246,0.2)',
          width: 1.2,
          curveness: 0.15
        }
      })
    }

    if (b.tags && Array.isArray(b.tags)) {
      b.tags.forEach(tag => {
        if (!tagSet.has(tag)) {
          tagSet.add(tag)
          const tagId = `tag_${tag}`
          const tagConn = connectionCount[tagId] || 1
          nodes.push({
            id: tagId,
            name: tag,
            category: 2,
            symbolSize: Math.min(12 + tagConn * 2, 30),
            value: tagConn,
            itemStyle: {
              color: isDark.value ? COLORS.tag.dark : COLORS.tag.light,
              shadowBlur: tagConn > 2 ? 8 : 3,
              shadowColor: COLORS.tag.glow
            },
            label: {
              show: true,
              fontSize: Math.min(8 + tagConn, 12),
              color: isDark.value ? '#f9fafb' : '#1f2937'
            },
            connections: tagConn
          })
        }
        links.push({
          source: bloggerId,
          target: `tag_${tag}`,
          lineStyle: {
            color: isDark.value ? 'rgba(167,139,250,0.2)' : 'rgba(139,92,246,0.15)',
            width: 0.8,
            curveness: 0.1
          }
        })
      })
    }

    if (b.platform && !platformSet.has(b.platform)) {
      platformSet.add(b.platform)
      const platId = `plat_${b.platform}`
      const platConn = connectionCount[platId] || 1
      nodes.push({
        id: platId,
        name: b.platform,
        category: 3,
        symbolSize: Math.min(14 + platConn * 3, 36),
        value: platConn,
        itemStyle: {
          color: isDark.value ? COLORS.platform.dark : COLORS.platform.light,
          shadowBlur: platConn > 2 ? 10 : 4,
          shadowColor: COLORS.platform.glow
        },
        label: {
          show: true,
          fontSize: Math.min(9 + platConn, 13),
          color: isDark.value ? '#f9fafb' : '#1f2937'
        },
        connections: platConn
      })
    }
    if (b.platform) {
      links.push({
        source: bloggerId,
        target: `plat_${b.platform}`,
        lineStyle: {
          color: isDark.value ? 'rgba(52,211,153,0.25)' : 'rgba(16,185,129,0.2)',
          width: 1,
          curveness: 0.12
        }
      })
    }
  })

  return { nodes, links, categories }
}

const legendItems = computed(() => [
  { key: 'blogger', label: '博主', color: isDark.value ? COLORS.blogger.dark : COLORS.blogger.light, count: props.bloggers.length },
  { key: 'category', label: '分类', color: isDark.value ? COLORS.category.dark : COLORS.category.light, count: new Set(props.bloggers.map(b => b.category).filter(Boolean)).size },
  { key: 'tag', label: '标签', color: isDark.value ? COLORS.tag.dark : COLORS.tag.light, count: new Set(props.bloggers.flatMap(b => b.tags || [])).size },
  { key: 'platform', label: '平台', color: isDark.value ? COLORS.platform.dark : COLORS.platform.light, count: new Set(props.bloggers.map(b => b.platform).filter(Boolean)).size }
])

const filteredNodes = computed(() => {
  const { nodes } = buildGraphData()
  return nodes.filter(n => {
    const keys = ['blogger', 'category', 'tag', 'platform']
    return visibleCategories.value[keys[n.category]]
  })
})

const filteredLinks = computed(() => {
  const { nodes, links } = buildGraphData()
  const nodeIds = new Set(filteredNodes.value.map(n => n.id))
  return links.filter(l => nodeIds.has(l.source) && nodeIds.has(l.target))
})

const toggleCategory = (key) => {
  visibleCategories.value[key] = !visibleCategories.value[key]
  renderChart()
}

const renderChart = () => {
  const el = document.getElementById('relation-graph-chart')
  if (!el || props.bloggers.length === 0) return

  if (!chart) {
    chart = echarts.init(el)
    chart.on('mouseover', (params) => {
      if (params.dataType === 'node') {
        hoveredNode.value = params.data
        const rect = el.getBoundingClientRect()
        tooltipStyle.value = {
          left: `${params.event.offsetX + 16}px`,
          top: `${params.event.offsetY - 8}px`
        }
      }
    })
    chart.on('mouseout', (params) => {
      if (params.dataType === 'node') {
        hoveredNode.value = null
      }
    })
  }

  const { nodes: allNodes, links: allLinks, categories } = buildGraphData()

  const activeKeys = Object.entries(visibleCategories.value)
    .filter(([, v]) => v)
    .map(([k]) => categoryIndexMap[k])

  const nodes = allNodes.filter(n => activeKeys.includes(n.category))

  const nodeIds = new Set(nodes.map(n => n.id))
  const links = allLinks.filter(l => nodeIds.has(l.source) && nodeIds.has(l.target))

  const searchLower = searchText.value.toLowerCase()
  const matchedIds = new Set()
  if (searchLower) {
    nodes.forEach(n => {
      if (n.name.toLowerCase().includes(searchLower)) {
        matchedIds.add(n.id)
      }
    })
    allLinks.forEach(l => {
      if (matchedIds.has(l.source)) matchedIds.add(l.target)
      if (matchedIds.has(l.target)) matchedIds.add(l.source)
    })
  }

  const finalNodes = nodes.map(n => {
    const isMatch = searchLower && matchedIds.has(n.id)
    const isDimmed = searchLower && !matchedIds.has(n.id)
    return {
      ...n,
      symbolSize: isMatch ? n.symbolSize * 1.3 : n.symbolSize,
      itemStyle: {
        ...n.itemStyle,
        opacity: isDimmed ? 0.15 : 1,
        shadowBlur: isMatch ? 20 : n.itemStyle.shadowBlur
      },
      label: {
        ...n.label,
        show: isMatch || !searchLower,
        fontSize: isMatch ? n.label.fontSize + 2 : n.label.fontSize,
        opacity: isDimmed ? 0.15 : 1
      }
    }
  })

  const finalLinks = links.map(l => {
    const isMatch = searchLower && (matchedIds.has(l.source) || matchedIds.has(l.target))
    return {
      ...l,
      lineStyle: {
        ...l.lineStyle,
        opacity: searchLower ? (isMatch ? 0.6 : 0.03) : l.lineStyle.width > 1 ? 0.4 : 0.2
      }
    }
  })

  chart.setOption({
    tooltip: {
      show: false
    },
    series: [{
      type: 'graph',
      layout: 'force',
      animation: true,
      animationDuration: 800,
      animationEasingUpdate: 'quinticInOut',
      draggable: true,
      data: finalNodes,
      links: finalLinks,
      categories: categories,
      roam: true,
      zoom: currentZoom,
      force: {
        repulsion: 280,
        gravity: 0.08,
        edgeLength: [50, 180],
        layoutAnimation: true,
        friction: 0.6
      },
      label: {
        position: 'right',
        formatter: '{b}',
        distance: 5
      },
      lineStyle: {
        curveness: 0.15
      },
      emphasis: {
        focus: 'adjacency',
        itemStyle: {
          shadowBlur: 20,
          borderWidth: 2,
          borderColor: isDark.value ? '#f9fafb' : '#1f2937'
        },
        lineStyle: {
          width: 3,
          opacity: 0.8
        },
        label: {
          show: true,
          fontSize: 14,
          fontWeight: 'bold'
        }
      },
      blur: {
        itemStyle: { opacity: 0.08 },
        lineStyle: { opacity: 0.03 },
        label: { opacity: 0.08 }
      }
    }]
  }, true)
}

const zoomIn = () => {
  currentZoom = Math.min(currentZoom * 1.2, 3)
  if (chart) chart.setOption({ series: [{ zoom: currentZoom }] })
}

const zoomOut = () => {
  currentZoom = Math.max(currentZoom / 1.2, 0.3)
  if (chart) chart.setOption({ series: [{ zoom: currentZoom }] })
}

const resetZoom = () => {
  currentZoom = 1
  if (chart) chart.setOption({ series: [{ zoom: 1, center: undefined }] })
}

const handleResize = () => {
  if (chart) chart.resize()
}

watch(() => props.bloggers, () => {
  nextTick(renderChart)
}, { deep: true })

watch(searchText, () => {
  if (chart) renderChart()
})

watch(() => themeStore.isDark, () => {
  nextTick(renderChart)
})

onMounted(() => {
  nextTick(renderChart)
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (chart) {
    chart.dispose()
    chart = null
  }
})
</script>

<style scoped>
.relation-graph {
  position: relative;
  width: 100%;
  height: 600px;
  overflow: hidden;
}

.graph-controls {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
  gap: 12px;
  flex-wrap: wrap;
}

.graph-legend {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-radius: 20px;
  border: 1px solid var(--border-color, #e5e7eb);
  background: var(--bg-card, #fff);
  cursor: pointer;
  font-size: 13px;
  color: var(--text-secondary, #4b5563);
  transition: all 0.2s;
  -webkit-tap-highlight-color: transparent;
  user-select: none;
}

.legend-item:hover {
  border-color: var(--primary, #f97316);
  background: var(--bg-hover, #f3f4f6);
}

.legend-item.dimmed {
  opacity: 0.4;
}

.legend-item.active {
  border-color: var(--primary, #f97316);
  background: rgba(249, 115, 22, 0.06);
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.legend-count {
  font-size: 11px;
  color: var(--text-muted, #9ca3af);
  margin-left: 2px;
}

.graph-actions {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.graph-search {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-radius: 20px;
  border: 1px solid var(--border-color, #e5e7eb);
  background: var(--bg-card, #fff);
  transition: border-color 0.2s;
}

.graph-search:focus-within {
  border-color: var(--primary, #f97316);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

.graph-search svg {
  width: 14px;
  height: 14px;
  color: var(--text-muted, #9ca3af);
  flex-shrink: 0;
}

.graph-search input {
  border: none;
  outline: none;
  background: transparent;
  font-size: 13px;
  color: var(--text-primary, #1f2937);
  width: 120px;
}

.graph-search input::placeholder {
  color: var(--text-muted, #9ca3af);
}

.search-clear {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border: none;
  background: var(--bg-hover, #f3f4f6);
  border-radius: 50%;
  cursor: pointer;
  color: var(--text-muted, #9ca3af);
  padding: 0;
  flex-shrink: 0;
}

.search-clear svg {
  width: 10px;
  height: 10px;
}

.graph-zoom {
  display: flex;
  gap: 2px;
  align-items: center;
  border-radius: 20px;
  border: 1px solid var(--border-color, #e5e7eb);
  background: var(--bg-card, #fff);
  padding: 2px;
}

.graph-zoom button {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary, #4b5563);
  transition: all 0.15s;
  -webkit-tap-highlight-color: transparent;
}

.graph-zoom button:hover {
  background: var(--bg-hover, #f3f4f6);
  color: var(--primary, #f97316);
}

.graph-zoom button svg {
  width: 16px;
  height: 16px;
}

.zoom-level {
  font-size: 11px;
  color: var(--text-muted, #9ca3af);
  min-width: 36px;
  text-align: center;
  font-variant-numeric: tabular-nums;
}

.reset-btn {
  margin-left: 2px;
  border-left: 1px solid var(--border-color, #e5e7eb) !important;
  padding-left: 2px;
}

.graph-canvas-wrapper {
  position: relative;
  width: 100%;
  height: calc(100% - 52px);
  border-radius: 16px;
  border: 1px solid var(--border-color, #e5e7eb);
  background: var(--bg-card, #fff);
  overflow: hidden;
}

.dark .graph-canvas-wrapper {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
}

.graph-chart {
  width: 100%;
  height: 100%;
}

.graph-empty {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  color: var(--text-muted, #9ca3af);
}

.graph-empty svg {
  width: 48px;
  height: 48px;
  margin-bottom: 12px;
  opacity: 0.4;
}

.graph-empty p {
  font-size: 15px;
  margin: 0 0 4px;
  color: var(--text-secondary, #4b5563);
}

.graph-empty span {
  font-size: 12px;
}

.node-tooltip {
  position: absolute;
  pointer-events: none;
  padding: 10px 14px;
  border-radius: 10px;
  background: var(--bg-card, #fff);
  border: 1px solid var(--border-color, #e5e7eb);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  z-index: 10;
  min-width: 120px;
  transition: opacity 0.15s;
}

.dark .node-tooltip {
  background: #1e293b;
  border-color: #334155;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

.tooltip-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.tooltip-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.tooltip-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary, #1f2937);
}

.dark .tooltip-name {
  color: #f9fafb;
}

.tooltip-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--text-muted, #9ca3af);
}

.tooltip-type {
  padding: 1px 8px;
  border-radius: 10px;
  background: var(--bg-hover, #f3f4f6);
  font-size: 11px;
}

.dark .tooltip-type {
  background: #334155;
}

.graph-stats {
  position: absolute;
  bottom: 12px;
  left: 12px;
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: var(--text-muted, #9ca3af);
  padding: 4px 10px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(8px);
}

.dark .graph-stats {
  background: rgba(15, 23, 42, 0.7);
  color: var(--text-muted);
}

@media (max-width: 768px) {
  .relation-graph {
    height: 450px;
  }
  .graph-controls {
    flex-direction: column;
    gap: 8px;
  }
  .graph-legend {
    gap: 4px;
  }
  .legend-item {
    padding: 4px 8px;
    font-size: 12px;
  }
  .graph-actions {
    width: 100%;
    justify-content: space-between;
  }
  .graph-search input {
    width: 80px;
  }
  .graph-zoom button {
    width: 36px;
    height: 36px;
  }
  .node-tooltip {
    display: none;
  }
}

@media (max-width: 480px) {
  .relation-graph {
    height: 350px;
  }
  .legend-item {
    padding: 3px 6px;
    font-size: 11px;
    gap: 4px;
  }
  .legend-count {
    display: none;
  }
  .graph-search input {
    width: 60px;
  }
}
</style>
