<template>
  <div class="skeleton-loader" :class="{ 'dark-mode': isDark, [variant]: true }">
    <template v-if="variant === 'blogger-card'">
      <div v-for="n in count" :key="n" class="ske-card" :style="{ '--i': n }">
        <div class="ske-checkbox"></div>
        <div class="ske-inner">
          <div class="ske-avatar"></div>
          <div class="ske-content">
            <div class="ske-line ske-title"></div>
            <div class="ske-line ske-subtitle"></div>
            <div class="ske-tags">
              <div class="ske-tag"></div>
              <div class="ske-tag short"></div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <template v-else-if="variant === 'blogger-list'">
      <div class="ske-list-header">
        <div class="ske-line ske-h wide"></div>
      </div>
      <div v-for="n in count" :key="n" class="ske-list-row" :style="{ '--i': n }">
        <div class="ske-cell narrow"><div class="ske-circle"></div></div>
        <div class="ske-cell"><div class="ske-line ske-text"></div></div>
        <div class="ske-cell narrow"><div class="ske-badge"></div></div>
        <div class="ske-cell narrow"><div class="ske-line ske-text sm"></div></div>
        <div class="ske-cell narrow"><div class="ske-line ske-text sm"></div></div>
        <div class="ske-cell med"><div class="ske-tags-inline"><div class="ske-tag-xs"></div><div class="ske-tag-xs"></div></div></div>
        <div class="ske-cell narrow"><div class="ske-actions"></div></div>
      </div>
    </template>

    <template v-else-if="variant === 'stats'">
      <div class="ske-stats-grid">
        <div v-for="n in 4" :key="n" class="ske-stat-card" :style="{ '--i': n }">
          <div class="ske-stat-icon"></div>
          <div class="ske-stat-body">
            <div class="ske-line ske-value"></div>
            <div class="ske-line ske-label"></div>
          </div>
        </div>
      </div>
    </template>

    <template v-else-if="variant === 'table'">
      <div class="ske-table">
        <div class="ske-thead">
          <div v-for="n in columns" :key="n" class="ske-th"><div class="ske-line ske-th-text"></div></div>
        </div>
        <div v-for="r in rows" :key="r" class="ske-tr" :style="{ '--i': r }">
          <div v-for="c in columns" :key="c" class="ske-td">
            <div class="ske-line" :class="{ 'ske-sm': c % 3 === 0, 'ske-md': c % 2 === 0 }"></div>
          </div>
        </div>
      </div>
    </template>

    <template v-else-if="variant === 'detail'">
      <div class="ske-detail">
        <div class="ske-detail-header">
          <div class="ske-detail-avatar-lg"></div>
          <div class="ske-detail-info">
            <div class="ske-line ske-d-name"></div>
            <div class="ske-line ske-d-meta"></div>
            <div class="ske-badges-row">
              <div class="ske-badge-pill"></div>
              <div class="ske-badge-pill short"></div>
            </div>
          </div>
        </div>
        <div class="ske-detail-body">
          <div class="ske-section">
            <div class="ske-line ske-section-title"></div>
            <div class="ske-field"><div class="ske-label-sm"></div><div class="ske-line ske-field-val"></div></div>
            <div class="ske-field"><div class="ske-label-sm"></div><div class="ske-line ske-field-val mid"></div></div>
            <div class="ske-field"><div class="ske-label-sm"></div><div class="ske-line ske-field-val long"></div></div>
          </div>
          <div class="ske-section">
            <div class="ske-line ske-section-title"></div>
            <div class="ske-field"><div class="ske-label-sm"></div><div class="ske-line ske-field-val"></div></div>
            <div class="ske-field"><div class="ske-label-sm"></div><div class="ske-line ske-field-val mid"></div></div>
          </div>
        </div>
      </div>
    </template>

    <template v-else>
      <div v-for="n in count" :key="n" class="ske-block" :style="{ '--i': n }">
        <div class="ske-line" :style="{ width: widths[n % widths.length] }"></div>
        <div v-if="n < count" class="ske-line sm" :style="{ width: widths[(n + 1) % widths.length] }"></div>
      </div>
    </template>
  </div>
</template>

<script setup>
const props = defineProps({
  variant: {
    type: String,
    default: 'text',
    validator: (v) => ['text', 'blogger-card', 'blogger-list', 'stats', 'table', 'detail'].includes(v),
  },
  count: { type: Number, default: 6 },
  columns: { type: Number, default: 5 },
  rows: { type: Number, default: 8 },
  isDark: { type: Boolean, default: false },
})

const widths = ['85%', '95%', '70%', '60%', '90%']
</script>

<style scoped>
.skeleton-loader { --ske-bg: #f0f0f0; --ske-shine: linear-gradient(90deg, transparent 0%, rgba(255,255,255,0.4) 50%, transparent 100%); }
.dark-mode { --ske-bg: #2a2a35; --ske-shine: linear-gradient(90deg, transparent 0%, rgba(255,255,255,0.06) 50%, transparent 100%); }

@keyframes ske-shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}

@keyframes ske-fade-in {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}

.ske-line,
.ske-circle,
.ske-avatar,
.ske-tag,
.ske-badge,
.ske-checkbox,
.ske-stat-icon,
.ske-detail-avatar-lg,
.ske-badge-pill,
.ske-tag-xs,
.ske-actions,
.ske-label-sm {
  background: var(--ske-bg);
  border-radius: 6px;
  animation: ske-shimmer 1.8s ease-in-out infinite;
  background-size: 200% 100%;
}

.ske-block,
.ske-card,
.ske-list-row,
.ske-stat-card,
.ske-tr,
.ske-section {
  animation: ske-fade-in 0.4s ease both;
  animation-delay: calc((var(--i) - 1) * 60ms);
}

.ske-line { height: 14px; margin-bottom: 8px; }
.ske-line.sm { height: 10px; width: 60% !important; }

/* Blogger Card Skeleton */
.ske-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-radius: 12px;
  background: var(--card-bg, #fff);
  border: 1px solid var(--border-color, #eee);
  min-height: 96px;
}
.dark-mode .ske-card { --card-bg: #1e1e2d; --border-color: #2a2a38; }
.ske-checkbox { width: 18px; height: 18px; border-radius: 4px; flex-shrink: 0; }
.ske-inner { display: flex; align-items: center; gap: 14px; flex: 1; }
.ske-avatar { width: 48px; height: 48px; border-radius: 12px; flex-shrink: 0; }
.ske-content { flex: 1; display: flex; flex-direction: column; gap: 7px; }
.ske-title { width: 45%; height: 16px; }
.ske-subtitle { width: 30%; height: 12px; }
.ske-tags { display: flex; gap: 6px; }
.ske-tag { width: 52px; height: 20px; border-radius: 10px; }
.ske-tag.short { width: 36px; }

/* Blogger List Skeleton */
.ske-list-header { padding: 10px 16px; }
.ske-list-row {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  border-bottom: 1px solid var(--border-color, #f0f0f0);
}
.dark-mode .ske-list-row { --border-color: #2a2a38; }
.ske-cell { flex: 1; display: flex; align-items: center; }
.ske-cell.narrow { flex: 0 0 56px; justify-content: center; }
.ske-cell.med { flex: 1.5; }
.ske-circle { width: 32px; height: 32px; border-radius: 50%; }
.ske-text { width: 65%; height: 13px; }
.ske-text.sm { width: 45%; height: 11px; }
.ske-badge { width: 64px; height: 22px; border-radius: 11px; }
.ske-tags-inline { display: flex; gap: 4px; }
.ske-tag-xs { width: 40px; height: 16px; border-radius: 8px; }
.ske-actions { width: 60px; height: 24px; border-radius: 6px; }
.ske-h.wide { width: 40%; height: 15px; }

/* Stats Skeleton */
.ske-stats-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; }
.ske-stat-card {
  display: flex; align-items: center; gap: 12px;
  padding: 16px; border-radius: 12px;
  background: var(--card-bg, #fff);
  border: 1px solid var(--border-color, #eee);
}
.dark-mode .ske-stat-card { --card-bg: #1e1e2d; --border-color: #2a2a38; }
.ske-stat-icon { width: 42px; height: 42px; border-radius: 10px; flex-shrink: 0; }
.ske-stat-body { flex: 1; display: flex; flex-direction: column; gap: 6px; }
.ske-value { width: 55%; height: 22px; }
.ske-label { width: 35%; height: 11px; }

/* Table Skeleton */
.ske-table { width: 100%; border-radius: 10px; overflow: hidden; border: 1px solid var(--border-color, #eee); }
.dark-mode .ske-table { --border-color: #2a2a38; }
.ske-thead { display: flex; background: var(--ske-bg); padding: 10px 12px; }
.ske-th { flex: 1; }
.ske-th-text { width: 60%; height: 12px; }
.ske-tr { display: flex; padding: 10px 12px; border-top: 1px solid var(--border-color, #f0f0f0); }
.dark-mode .ske-tr { --border-color: #2a2a38; }
.ske-td { flex: 1; padding: 2px 4px; }
.ske-sm { width: 40% !important; }
.ske-md { width: 65% !important; }

/* Detail Page Skeleton */
.ske-detail { max-width: 800px; }
.ske-detail-header { display: flex; gap: 20px; padding: 24px; }
.ske-detail-avatar-lg { width: 88px; height: 88px; border-radius: 20px; flex-shrink: 0; }
.ske-detail-info { flex: 1; display: flex; flex-direction: column; gap: 10px; }
.ske-d-name { width: 40%; height: 22px; }
.ske-d-meta { width: 55%; height: 13px; }
.ske-badges-row { display: flex; gap: 8px; }
.ske-badge-pill { width: 72px; height: 24px; border-radius: 12px; }
.ske-badge-pill.short { width: 52px; }
.ske-detail-body { padding: 0 24px 24px; }
.ske-section { margin-bottom: 20px; }
.ske-section-title { width: 25%; height: 16px; margin-bottom: 14px; }
.ske-field { display: flex; align-items: center; gap: 12px; margin-bottom: 10px; }
.ske-label-sm { width: 70px; height: 11px; flex-shrink: 0; }
.ske-field-val { flex: 1; height: 13px; }
.ske-field-val.mid { width: 55%; }
.ske-field-val.long { width: 80%; }

@media (max-width: 768px) {
  .ske-stats-grid { grid-template-columns: repeat(2, 1fr); }
  .ske-table { overflow-x: auto; }
  .ske-detail-header { flex-direction: column; align-items: center; text-align: center; }
  .ske-detail-info { align-items: center; }
}

@media (max-width: 480px) {
  .ske-stats-grid { grid-template-columns: 1fr 1fr; gap: 8px; }
  .ske-stat-card { padding: 12px; }
  .ske-stat-icon { width: 34px; height: 34px; }
  .ske-card { padding: 12px; min-height: 80px; }
  .ske-avatar { width: 40px; height: 40px; border-radius: 10px; }
}
</style>
