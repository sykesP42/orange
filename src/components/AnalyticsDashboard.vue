<template>
  <div class="analytics-dashboard">
    <div class="dashboard-header">
      <h2>📊 数据分析</h2>
      <div class="time-range">
        <button
          v-for="range in timeRanges"
          :key="range.value"
          :class="{ active: selectedRange === range.value }"
          @click="selectedRange = range.value"
        >
          {{ range.label }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <template v-else-if="data">
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">👥</div>
          <div class="stat-info">
            <span class="stat-value">{{ data.total_bloggers }}</span>
            <span class="stat-label">总博主数</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">✅</div>
          <div class="stat-info">
            <span class="stat-value">{{ data.active_bloggers }}</span>
            <span class="stat-label">活跃博主</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">🤝</div>
          <div class="stat-info">
            <span class="stat-value">{{ data.cooperation_count }}</span>
            <span class="stat-label">合作次数</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">👁️</div>
          <div class="stat-info">
            <span class="stat-value">{{ formatNumber(data.total_exposure) }}</span>
            <span class="stat-label">总曝光量</span>
          </div>
        </div>
      </div>

      <div class="charts-row">
        <div class="chart-card">
          <h3>📈 月度趋势</h3>
          <div class="chart-content">
            <div class="bar-chart">
              <div
                v-for="(month, index) in data.monthly_trend"
                :key="index"
                class="bar-group"
              >
                <div class="bar-container">
                  <div
                    class="bar added"
                    :style="{ height: getBarHeight(month.added) + '%' }"
                    :title="`新增: ${month.added}`"
                  ></div>
                  <div
                    class="bar cooperated"
                    :style="{ height: getBarHeight(month.cooperated) + '%' }"
                    :title="`合作: ${month.cooperated}`"
                  ></div>
                </div>
                <span class="bar-label">{{ month.month.slice(5) }}</span>
              </div>
            </div>
            <div class="chart-legend">
              <span class="legend-item"><span class="dot added"></span> 新增博主</span>
              <span class="legend-item"><span class="dot cooperated"></span> 合作次数</span>
            </div>
          </div>
        </div>

        <div class="chart-card">
          <h3>🏷️ 状态分布</h3>
          <div class="chart-content">
            <div class="pie-chart">
              <div
                v-for="(status, index) in data.status_distribution"
                :key="index"
                class="pie-item"
              >
                <div class="pie-label">
                  <span class="pie-color" :style="{ backgroundColor: getStatusColor(index) }"></span>
                  <span class="pie-name">{{ status.status }}</span>
                </div>
                <div class="pie-bar">
                  <div
                    class="pie-fill"
                    :style="{
                      width: status.percent + '%',
                      backgroundColor: getStatusColor(index)
                    }"
                  ></div>
                </div>
                <span class="pie-percent">{{ status.percent.toFixed(1) }}%</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="charts-row">
        <div class="chart-card">
          <h3>📂 分类排名</h3>
          <div class="chart-content">
            <div class="rank-list">
              <div
                v-for="(cat, index) in data.top_categories"
                :key="index"
                class="rank-item"
              >
                <span class="rank-number" :class="getRankClass(index)">{{ index + 1 }}</span>
                <span class="rank-name">{{ cat.category }}</span>
                <span class="rank-value">{{ cat.count }}</span>
                <div class="rank-bar">
                  <div
                    class="rank-fill"
                    :style="{
                      width: (cat.count / data.top_categories[0].count * 100) + '%',
                      backgroundColor: getCategoryColor(index)
                    }"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="chart-card">
          <h3>📱 平台分布</h3>
          <div class="chart-content">
            <div class="rank-list">
              <div
                v-for="(plt, index) in data.top_platforms"
                :key="index"
                class="rank-item"
              >
                <span class="rank-number" :class="getRankClass(index)">{{ index + 1 }}</span>
                <span class="rank-name">{{ plt.platform }}</span>
                <span class="rank-value">{{ plt.count }}</span>
                <div class="rank-bar">
                  <div
                    class="rank-fill"
                    :style="{
                      width: (plt.count / data.top_platforms[0].count * 100) + '%',
                      backgroundColor: getPlatformColor(index)
                    }"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAnalyticsOverviewAPI } from '../api'

const loading = ref(false)
const data = ref(null)
const selectedRange = ref('6m')
const timeRanges = [
  { label: '1个月', value: '1m' },
  { label: '3个月', value: '3m' },
  { label: '6个月', value: '6m' },
  { label: '1年', value: '1y' },
]

const loadAnalytics = async () => {
  loading.value = true
  try {
    const res = await getAnalyticsOverviewAPI()
    if (res.code === 200) {
      data.value = res.data
    }
  } catch (error) {
    console.error('加载统计数据失败', error)
  } finally {
    loading.value = false
  }
}

const formatNumber = (num) => {
  if (num >= 100000000) {
    return (num / 100000000).toFixed(1) + '亿'
  }
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toLocaleString()
}

const getBarHeight = (value) => {
  if (!data.value || !data.value.monthly_trend) return 0
  const max = Math.max(...data.value.monthly_trend.map(m => Math.max(m.added, m.cooperated)))
  return max > 0 ? (value / max * 100) : 0
}

const getStatusColor = (index) => {
  const colors = ['#10b981', '#3b82f6', '#f59e0b', '#ef4444', '#8b5cf6', '#6b7280']
  return colors[index % colors.length]
}

const getCategoryColor = (index) => {
  const colors = ['#f97316', '#3b82f6', '#10b981', '#8b5cf6', '#ec4899']
  return colors[index % colors.length]
}

const getPlatformColor = (index) => {
  const colors = ['#ef4444', '#f97316', '#eab308', '#22c55e', '#3b82f6']
  return colors[index % colors.length]
}

const getRankClass = (index) => {
  if (index === 0) return 'gold'
  if (index === 1) return 'silver'
  if (index === 2) return 'bronze'
  return ''
}

onMounted(() => {
  loadAnalytics()
})
</script>

<style scoped>
.analytics-dashboard {
  padding: 24px;
  background: #f9fafb;
  min-height: 100vh;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.dashboard-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.time-range {
  display: flex;
  gap: 8px;
  background: white;
  padding: 4px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.time-range button {
  padding: 8px 16px;
  background: transparent;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
}

.time-range button.active {
  background: var(--primary);
  color: white;
}

.loading {
  text-align: center;
  padding: 60px;
  color: #9ca3af;
  font-size: 16px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  font-size: 40px;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
}

.charts-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.chart-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.chart-card h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 20px 0;
}

.bar-chart {
  display: flex;
  justify-content: space-around;
  align-items: flex-end;
  height: 200px;
  padding: 20px 0;
}

.bar-group {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.bar-container {
  display: flex;
  gap: 4px;
  align-items: flex-end;
  height: 160px;
}

.bar {
  width: 24px;
  border-radius: 4px 4px 0 0;
  transition: all 0.3s;
}

.bar.added {
  background: linear-gradient(180deg, var(--primary), var(--primary-light));
}

.bar.cooperated {
  background: linear-gradient(180deg, #10b981, #34d399);
}

.bar:hover {
  transform: scaleY(1.05);
  transform-origin: bottom;
}

.bar-label {
  font-size: 12px;
  color: #6b7280;
}

.chart-legend {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-top: 16px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #6b7280;
}

.legend-item .dot {
  width: 12px;
  height: 12px;
  border-radius: 3px;
}

.legend-item .dot.added {
  background: var(--primary);
}

.legend-item .dot.cooperated {
  background: #10b981;
}

.pie-chart {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.pie-item {
  display: grid;
  grid-template-columns: 100px 1fr 60px;
  align-items: center;
  gap: 12px;
}

.pie-label {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pie-color {
  width: 12px;
  height: 12px;
  border-radius: 3px;
}

.pie-name {
  font-size: 13px;
  color: #374151;
}

.pie-bar {
  height: 8px;
  background: #f3f4f6;
  border-radius: 4px;
  overflow: hidden;
}

.pie-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.5s ease;
}

.pie-percent {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  text-align: right;
}

.rank-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.rank-item {
  display: grid;
  grid-template-columns: 30px 1fr 60px 100px;
  align-items: center;
  gap: 12px;
}

.rank-number {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  background: #f3f4f6;
  color: #6b7280;
}

.rank-number.gold {
  background: linear-gradient(135deg, #fbbf24, #f59e0b);
  color: white;
}

.rank-number.silver {
  background: linear-gradient(135deg, #9ca3af, #6b7280);
  color: white;
}

.rank-number.bronze {
  background: linear-gradient(135deg, #d97706, #92400e);
  color: white;
}

.rank-name {
  font-size: 14px;
  color: #374151;
}

.rank-value {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  text-align: right;
}

.rank-bar {
  height: 8px;
  background: #f3f4f6;
  border-radius: 4px;
  overflow: hidden;
}

.rank-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.5s ease;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .charts-row {
    grid-template-columns: 1fr;
  }
}
</style>
