<template>
  <div class="skeleton-detail">
    <div class="skeleton-detail-header">
      <div class="skeleton-avatar-lg"></div>
      <div class="skeleton-header-info">
        <div class="skeleton-title-lg"></div>
        <div class="skeleton-text-sm"></div>
        <div class="skeleton-text-xs"></div>
      </div>
    </div>

    <div class="skeleton-detail-body">
      <div class="skeleton-section">
        <div class="skeleton-section-title"></div>
        <div class="skeleton-stats-row">
          <div class="skeleton-stat-box" v-for="i in 3" :key="'stat-'+i"></div>
        </div>
      </div>

      <div class="skeleton-section">
        <div class="skeleton-section-title"></div>
        <div class="skeleton-text" v-for="i in 3" :key="'line-'+i"></div>
      </div>

      <div class="skeleton-actions">
        <div class="skeleton-btn primary"></div>
        <div class="skeleton-btn secondary"></div>
        <div class="skeleton-btn danger"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({})
</script>

<style scoped>
.skeleton-detail {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 24px;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-color);
}

.skeleton-detail-header {
  display: flex;
  gap: 20px;
  margin-bottom: 28px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--border-light);
}

.skeleton-avatar-lg {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  flex-shrink: 0;
  background: linear-gradient(90deg, var(--bg-tertiary) 25%, var(--bg-dark) 50%, var(--bg-tertiary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.skeleton-header-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
  justify-content: center;
}

.skeleton-title-lg {
  height: 26px;
  width: 55%;
  border-radius: 6px;
  background: linear-gradient(90deg, var(--bg-tertiary) 25%, var(--bg-dark) 50%, var(--bg-tertiary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.skeleton-text-sm {
  height: 16px;
  width: 35%;
  border-radius: 4px;
  background: linear-gradient(90deg, var(--bg-tertiary) 25%, var(--bg-dark) 50%, var(--bg-tertiary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  animation-delay: 0.15s;
}

.skeleton-text-xs {
  height: 14px;
  width: 25%;
  border-radius: 4px;
  background: linear-gradient(90deg, var(--bg-tertiary) 25%, var(--bg-dark) 50%, var(--bg-tertiary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  animation-delay: 0.25s;
}

.skeleton-detail-body {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.skeleton-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.skeleton-section-title {
  height: 18px;
  width: 120px;
  border-radius: 4px;
  background: linear-gradient(90deg, var(--bg-tertiary) 25%, var(--bg-dark) 50%, var(--bg-tertiary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.skeleton-stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.skeleton-stat-box {
  height: 72px;
  border-radius: 12px;
  background: linear-gradient(90deg, var(--bg-tertiary) 25%, var(--bg-dark) 50%, var(--bg-tertiary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.skeleton-text {
  height: 16px;
  border-radius: 4px;
  background: linear-gradient(90deg, var(--bg-tertiary) 25%, var(--bg-dark) 50%, var(--bg-tertiary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.skeleton-text:nth-child(2) { width: 88%; animation-delay: 0.1s; }
.skeleton-text:nth-child(3) { width: 65%; animation-delay: 0.2s; }

.skeleton-actions {
  display: flex;
  gap: 12px;
  padding-top: 8px;
}

.skeleton-btn {
  height: 40px;
  border-radius: 10px;
  background: linear-gradient(90deg, var(--bg-tertiary) 25%, var(--bg-dark) 50%, var(--bg-tertiary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.skeleton-btn.primary { flex: 2; }
.skeleton-btn.secondary { flex: 1; }
.skeleton-btn.danger { flex: 1; }

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
