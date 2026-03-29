<template>
  <div class="skeleton-card" :style="{ height: height + 'px' }">
    <div class="skeleton-header">
      <div class="skeleton-avatar"></div>
      <div class="skeleton-title"></div>
    </div>
    <div class="skeleton-content">
      <div class="skeleton-text" v-for="i in lines" :key="i"></div>
    </div>
    <div class="skeleton-footer">
      <div class="skeleton-stat"></div>
      <div class="skeleton-stat"></div>
      <div class="skeleton-stat"></div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  height: {
    type: Number,
    default: 200
  },
  lines: {
    type: Number,
    default: 3
  }
})
</script>

<style scoped>
.skeleton-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 20px;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-color);
  overflow: hidden;
  position: relative;
}

.skeleton-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.skeleton-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(90deg, var(--bg-tertiary) 25%, var(--bg-dark) 50%, var(--bg-tertiary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.skeleton-title {
  flex: 1;
  height: 20px;
  border-radius: 4px;
  background: linear-gradient(90deg, var(--bg-tertiary) 25%, var(--bg-dark) 50%, var(--bg-tertiary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.skeleton-content {
  margin-bottom: 16px;
}

.skeleton-text {
  height: 16px;
  margin-bottom: 12px;
  border-radius: 4px;
  background: linear-gradient(90deg, var(--bg-tertiary) 25%, var(--bg-dark) 50%, var(--bg-tertiary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.skeleton-text:last-child {
  margin-bottom: 0;
}

.skeleton-footer {
  display: flex;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid var(--border-light);
}

.skeleton-stat {
  flex: 1;
  height: 24px;
  border-radius: 4px;
  background: linear-gradient(90deg, var(--bg-tertiary) 25%, var(--bg-dark) 50%, var(--bg-tertiary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

@keyframes skeleton-loading {
  0% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0 50%;
  }
}
</style>
