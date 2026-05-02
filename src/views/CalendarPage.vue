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
  padding: 16px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}

.back-btn:hover {
  background: #f9fafb;
  border-color: var(--primary);
  color: var(--primary);
}

.back-btn svg {
  width: 16px;
  height: 16px;
}

h1 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

@media (max-width: 480px) {
  .calendar-page {
    padding: 12px;
  }

  .page-header {
    gap: 8px;
    margin-bottom: 16px;
  }

  .back-btn {
    padding: 6px 10px;
    font-size: 12px;
  }

  h1 {
    font-size: 15px;
  }
}
</style>
