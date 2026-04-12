<template>
  <div class="relation-graph">
    <div class="graph-controls">
      <div class="graph-legend">
        <span class="legend-item"><span class="legend-dot" style="background: var(--primary)"></span>博主</span>
        <span class="legend-item"><span class="legend-dot" style="background: var(--info)"></span>分类</span>
        <span class="legend-item"><span class="legend-dot" style="background:#8b5cf6"></span>标签</span>
        <span class="legend-item"><span class="legend-dot" style="background: var(--success)"></span>平台</span>
      </div>
      <div class="graph-zoom">
        <button @click="zoomIn" title="放大">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </button>
        <button @click="zoomOut" title="缩小">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </button>
        <button @click="resetZoom" title="重置">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/></svg>
        </button>
      </div>
    </div>
    <div ref="chartRef" class="graph-chart"></div>
    <div v-if="bloggers.length === 0" class="graph-empty">暂无数据</div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts'

const props = defineProps({
  bloggers: { type: Array, default: () => [] }
})

const chartRef = ref(null)
let chart = null
let currentZoom = 1

const buildGraphData = () => {
  const nodes = []
  const links = []
  const categories = [
    { name: '博主' },
    { name: '分类' },
    { name: '标签' },
    { name: '平台' }
  ]

  const categorySet = new Set()
  const tagSet = new Set()
  const platformSet = new Set()

  props.bloggers.forEach(b => {
    nodes.push({
      id: `blogger_${b.id}`,
      name: b.nickname || b.real_name || '未知',
      category: 0,
      symbolSize: 30,
      itemStyle: { color: '#f97316' },
      label: { show: true, fontSize: 11 }
    })

    if (b.category && !categorySet.has(b.category)) {
      categorySet.add(b.category)
      nodes.push({
        id: `cat_${b.category}`,
        name: b.category,
        category: 1,
        symbolSize: 22,
        itemStyle: { color: '#3b82f6' },
        label: { show: true, fontSize: 10 }
      })
    }
    if (b.category) {
      links.push({ source: `blogger_${b.id}`, target: `cat_${b.category}`, lineStyle: { opacity: 0.3 } })
    }

    if (b.tags && Array.isArray(b.tags)) {
      b.tags.forEach(tag => {
        if (!tagSet.has(tag)) {
          tagSet.add(tag)
          nodes.push({
            id: `tag_${tag}`,
            name: tag,
            category: 2,
            symbolSize: 16,
            itemStyle: { color: '#8b5cf6' },
            label: { show: true, fontSize: 9 }
          })
        }
        links.push({ source: `blogger_${b.id}`, target: `tag_${tag}`, lineStyle: { opacity: 0.2 } })
      })
    }

    if (b.platform && !platformSet.has(b.platform)) {
      platformSet.add(b.platform)
      nodes.push({
        id: `plat_${b.platform}`,
        name: b.platform,
        category: 3,
        symbolSize: 20,
        itemStyle: { color: '#22c55e' },
        label: { show: true, fontSize: 10 }
      })
    }
    if (b.platform) {
      links.push({ source: `blogger_${b.id}`, target: `plat_${b.platform}`, lineStyle: { opacity: 0.25 } })
    }
  })

  return { nodes, links, categories }
}

const renderChart = () => {
  if (!chartRef.value || props.bloggers.length === 0) return

  if (!chart) {
    chart = echarts.init(chartRef.value)
  }

  const { nodes, links, categories } = buildGraphData()

  chart.setOption({
    tooltip: {
      formatter: (params) => {
        if (params.dataType === 'node') {
          return params.name
        }
        return `${params.data.source} → ${params.data.target}`
      }
    },
    legend: { show: false },
    series: [{
      type: 'graph',
      layout: 'force',
      animation: true,
      draggable: true,
      data: nodes,
      links: links,
      categories: categories,
      roam: true,
      zoom: currentZoom,
      force: {
        repulsion: 200,
        gravity: 0.1,
        edgeLength: [60, 150],
        layoutAnimation: true
      },
      label: {
        position: 'right',
        formatter: '{b}'
      },
      lineStyle: {
        color: 'source',
        curveness: 0.1
      },
      emphasis: {
        focus: 'adjacency',
        lineStyle: { width: 3 }
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
  height: 500px;
  overflow: hidden;
}

.graph-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.graph-legend {
  display: flex;
  gap: 16px;
  font-size: 13px;
  color: var(--text-secondary, #4b5563);
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.graph-zoom {
  display: flex;
  gap: 4px;
}

.graph-zoom button {
  width: 36px;
  height: 36px;
  border: 1px solid var(--border-color, #e5e7eb);
  background: var(--bg-card, #fff);
  border-radius: 8px;
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
  border-color: var(--primary, #f97316);
}

.graph-zoom button svg {
  width: 16px;
  height: 16px;
}

.graph-chart {
  width: 100%;
  height: calc(100% - 44px);
  border-radius: 12px;
  border: 1px solid var(--border-color, #e5e7eb);
  background: var(--bg-card, #fff);
}

.graph-empty {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: var(--text-muted, #9ca3af);
  font-size: 14px;
}

@media (max-width: 768px) {
  .relation-graph {
    height: 350px;
  }
  .graph-legend {
    flex-wrap: wrap;
    gap: 8px;
    font-size: 12px;
  }
  .graph-controls {
    flex-wrap: wrap;
    gap: 8px;
  }
  .graph-zoom button {
    width: 40px;
    height: 40px;
  }
  .graph-zoom button svg {
    width: 20px;
    height: 20px;
  }
}

@media (max-width: 480px) {
  .relation-graph {
    height: 280px;
  }
  .graph-legend {
    font-size: 11px;
    gap: 6px;
  }
  .legend-dot {
    width: 8px;
    height: 8px;
  }
}
</style>
