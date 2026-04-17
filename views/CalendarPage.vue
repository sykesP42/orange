<template>
  <div class="calendar-page">
    <div class="page-header">
      <button class="back-btn" @click="router.push('/')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M19 12H5M12 19l-7-7 7-7"/>
        </svg>
        返回主页
      </button>
      <h1>📅 博主提醒日历</h1>
    </div>
    <Calendar />
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import Calendar from '../components/Calendar.vue'

const router = useRouter()
</script>

<style scoped>
.calendar-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 24px;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.back-btn:hover {
  background: var(--bg-card-hover);
  border-color: var(--primary);
  color: var(--primary);
}

.back-btn svg {
  width: 18px;
  height: 18px;
}

h1 {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}
</style>
