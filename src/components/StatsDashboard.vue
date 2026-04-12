<template>
  <div class="stats-dashboard" :class="{ 'dark-mode': isDark }">
    <div class="stats-grid">
      <div
        v-for="(stat, index) in stats"
        :key="stat.key"
        class="stat-card"
        :class="[stat.key, { 'stat-loading': loading }]"
        :style="{ '--accent': stat.color, '--delay': index * 0.06 + 's' }"
      >
        <div class="stat-icon" :style="{ background: stat.bgGradient || `${stat.color}12` }">
          <component :is="'svg'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" v-html="stat.icon"></component>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ loading ? '...' : formatNumber(stat.value) }}</span>
          <span class="stat-label">{{ stat.label }}</span>
        </div>
        <div class="stat-trend" :class="stat.trendClass">
          <svg v-if="stat.trend === 'up'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/><polyline points="17 6 23 6 23 12"/></svg>
          <svg v-else-if="stat.trend === 'down'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 18 13.5 8.5 8.5 13.5 1 6"/><polyline points="17 18 23 18 23 12"/></svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="5" y1="12" x2="19" y2="12"/></svg>
          <span>{{ stat.trendText }}</span>
        </div>
      </div>
    </div>

    <div class="dashboard-row">
      <div class="quick-actions-panel">
        <h4 class="panel-title">⚡ 快捷操作</h4>
        <div class="action-grid">
          <button v-for="action in quickActions" :key="action.label" class="qa-btn" @click="$emit(action.emit)" :style="{ '--btn-accent': action.color }">
            <span class="qa-icon" v-html="action.icon"></span>
            <span class="qa-label">{{ action.label }}</span>
          </button>
        </div>
      </div>

      <div class="activity-feed-panel">
        <h4 class="panel-title">📈 近期动态</h4>
        <div class="feed-list" v-if="recentActivity.length > 0">
          <div v-for="(item, i) in recentActivity" :key="i" class="feed-item">
            <span class="feed-dot" :class="item.type"></span>
            <span class="feed-text">{{ item.text }}</span>
            <span class="feed-time">{{ item.time }}</span>
          </div>
        </div>
        <div v-else class="feed-empty">暂无动态</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  total: { type: Number, default: 0 },
  activeCount: { type: Number, default: 0 },
  pendingCount: { type: Number, default: 0 },
  invalidCount: { type: Number, default: 0 },
  todayNew: { type: Number, default: 0 },
  weekGrowth: { type: Number, default: 0 },
  cooperationRate: { type: Number, default: 0 },
  loading: { type: Boolean, default: false },
  isDark: { type: Boolean, default: false }
})

defineEmits(['add', 'invalid', 'export', 'import', 'kanban', 'workflow', 'calendar'])

const stats = computed(() => [
  {
    key: 'total',
    label: '博主总数',
    value: props.total,
    color: '#f97316',
    icon: '<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/>',
    trend: props.weekGrowth > 0 ? 'up' : props.weekGrowth < 0 ? 'down' : 'flat',
    trendText: props.weekGrowth > 0 ? `+${props.weekGrowth}%` : props.weekGrowth < 0 ? `${props.weekGrowth}%` : '持平'
  },
  {
    key: 'active',
    label: '合作中',
    value: props.activeCount,
    color: '#22c55e',
    icon: '<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/>',
    trend: props.cooperationRate > 50 ? 'up' : 'flat',
    trendText: `${props.cooperationRate}% 转化率`
  },
  {
    key: 'pending',
    label: '跟进中',
    value: props.pendingCount,
    color: '#3b82f6',
    icon: '<circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>',
    trend: props.pendingCount > 3 ? 'up' : 'flat',
    trendText: props.todayNew > 0 ? `今日+${props.todayNew}` : '待处理'
  },
  {
    key: 'invalid',
    label: '待处理',
    value: props.invalidCount,
    color: '#ef4444',
    icon: '<circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>',
    trend: props.invalidCount > 0 ? 'down' : 'flat',
    trendText: props.invalidCount > 0 ? '需处理' : '已清零'
  }
])

const quickActions = [
  { label: '录入博主', emit: 'add', color: '#f97316', icon: '<circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/>' },
  { label: '失效处理', emit: 'invalid', color: '#ef4444', icon: '<circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>' },
  { label: '看板视图', emit: 'kanban', color: '#8b5cf6', icon: '<rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/>' },
  { label: '工作流', emit: 'workflow', color: '#ec4899', icon: '<path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>' }
]

const recentActivity = computed(() => {
  const items = []
  if (props.todayNew > 0) items.push({ type: 'new', text: `今日新增 ${props.todayNew} 位博主`, time: '今天' })
  if (props.weekGrowth > 0) items.push({ type: 'growth', text: `本周增长 ${props.weekGrowth}%`, time: '本周' })
  if (props.invalidCount > 0) items.push({ type: 'warning', text: `${props.invalidCount} 位博主需要处理`, time: '待办' })
  items.push({ type: 'info', text: `当前共 ${props.total} 位博主在库`, time: '概览' })
  return items.slice(0, 4)
})

function formatNumber(n) {
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w'
  return n.toLocaleString()
}
</script>

<style scoped>
.stats-dashboard {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
}

.stat-card {
  position: relative;
  padding: 18px 20px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  display: flex;
  align-items: center;
  gap: 14px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  animation: card-in 0.4s ease-out both;
  animation-delay: var(--delay);
  overflow: hidden;
}

@keyframes card-in {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}

.stat-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 28px rgba(0, 0, 0, 0.08), 0 0 0 1px var(--border-color);
  border-color: transparent;
}

.stat-card.stat-loading { opacity: 0.6; }

.stat-icon {
  width: 46px; height: 46px;
  border-radius: 12px;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
  color: var(--accent, #f97316);
}
.stat-icon svg { width: 22px; height: 22px; }

.stat-content {
  display: flex; flex-direction: column; gap: 2px;
  min-width: 0;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
  letter-spacing: -0.5px;
}

.stat-label {
  font-size: 12.5px;
  font-weight: 500;
  color: var(--text-muted);
}

.stat-trend {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 3px;
  font-size: 11.5px;
  font-weight: 600;
  padding: 3px 7px;
  border-radius: 6px;
  white-space: nowrap;
  flex-shrink: 0;
}

.stat-trend svg { width: 13px; height: 13px; }

.stat-trend.up { color: #22c55e; background: rgba(34, 197, 94, 0.08); }
.stat-trend.down { color: #ef4444; background: rgba(239, 68, 68, 0.08); }
.stat-trend.flat { color: var(--text-muted); background: rgba(128, 128, 128, 0.06); }

.dashboard-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.quick-actions-panel,
.activity-feed-panel {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 14px;
  padding: 18px;
}

.panel-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 14px;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.qa-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 11px 14px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 11px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 13.5px;
  font-weight: 500;
  color: var(--text-secondary);
}

.qa-btn:hover {
  background: linear-gradient(135deg, rgba(var(--btn-accent-rgb, 249,115,22), 0.06), rgba(var(--btn-accent-rgb, 249,115,22), 0.02));
  border-color: var(--btn-accent, #f97316);
  color: var(--btn-accent, #f97316);
  transform: translateY(-1px);
}

.qa-icon { display: flex; }
.qa-icon :deep(svg) { width: 17px; height: 17px; }

.feed-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.feed-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 9px;
  transition: background 0.15s;
}

.feed-item:hover { background: var(--bg-secondary); }

.feed-dot {
  width: 8px; height: 8px;
  border-radius: 50%; flex-shrink: 0;
}
.feed-dot.new { background: #3b82f6; }
.feed-dot.growth { background: #22c55e; }
.feed-dot.warning { background: #f59e0b; }
.feed-dot.info { background: #8b5cf6; }

.feed-text {
  font-size: 13px;
  color: var(--text-secondary);
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.feed-time {
  font-size: 11px;
  color: var(--text-muted);
  flex-shrink: 0;
}

.feed-empty {
  text-align: center;
  padding: 20px 0;
  color: var(--text-muted);
  font-size: 13px;
}

@media (max-width: 1200px) {
  .stats-grid { grid-template-columns: repeat(2, 1fr); }
  .dashboard-row { grid-template-columns: 1fr; }
}

@media (max-width: 768px) {
  .stats-grid { grid-template-columns: repeat(2, 1fr); gap: 10px; }
  .stat-card { padding: 14px 14px; gap: 10px; }
  .stat-icon { width: 38px; height: 38px; }
  .stat-icon svg { width: 18px; height: 18px; }
  .stat-value { font-size: 20px; }
  .stat-trend { display: none; }
  .action-grid { grid-template-columns: repeat(2, 1fr); gap: 8px; }
  .qa-btn { padding: 10px 10px; font-size: 12.5px; }
  .qa-icon :deep(svg) { width: 15px; height: 15px; }
}

@media (max-width: 480px) {
  .stats-grid { grid-template-columns: 1fr 1fr; gap: 8px; }
  .stat-card { padding: 12px 10px; }
  .stat-value { font-size: 18px; }
  .stat-label { font-size: 11px; }
  .action-grid { grid-template-columns: 1fr 1fr; }
}
</style>
