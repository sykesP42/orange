<template>
  <div class="calendar-view">
    <div class="calendar-header">
      <div class="header-left">
        <button class="nav-btn" @click="previousPeriod">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15 18 9 12 15 6"/>
          </svg>
        </button>
        <h2>{{ currentTitle }}</h2>
        <button class="nav-btn" @click="nextPeriod">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        <button class="today-btn" @click="goToToday">今天</button>
      </div>
      <div class="header-right">
        <div class="view-tabs">
          <button :class="{ active: viewMode === 'month' }" @click="viewMode = 'month'">月</button>
          <button :class="{ active: viewMode === 'week' }" @click="viewMode = 'week'">周</button>
          <button :class="{ active: viewMode === 'day' }" @click="viewMode = 'day'">日</button>
        </div>
        <button class="add-event-btn" @click="showQuickAdd = true" title="快速添加事件">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </button>
      </div>
    </div>

    <div class="upcoming-bar" v-if="upcomingEvents.length > 0">
      <div class="upcoming-label">🔔 即将到期</div>
      <div class="upcoming-list">
        <button v-for="(evt, i) in upcomingEvents.slice(0, 5)" :key="i" class="upcoming-chip" :class="evt.urgency" @click="showEventDetail(evt)">
          {{ evt.title }}
        </button>
      </div>
    </div>

    <div class="calendar-grid" v-if="viewMode === 'month'">
      <div class="weekday-header">
        <div v-for="day in weekdays" :key="day" class="weekday">{{ day }}</div>
      </div>
      <div class="days-grid">
        <div
          v-for="(day, index) in calendarDays"
          :key="index"
          class="day-cell"
          :class="{
            'other-month': !day.isCurrentMonth,
            'today': day.isToday,
            'has-events': day.events.length > 0
          }"
          @click="selectDay(day)"
        >
          <span class="day-number">{{ day.date.getDate() }}</span>
          <div class="day-events">
            <div
              v-for="(event, idx) in day.events.slice(0, 3)"
              :key="idx"
              class="event-dot"
              :style="{ backgroundColor: event.color }"
              :title="event.title"
            >
              {{ event.title.substring(0, 8) }}
            </div>
            <div v-if="day.events.length > 3" class="more-events">
              +{{ day.events.length - 3 }} 更多
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="week-view" v-if="viewMode === 'week'">
      <div class="week-header">
        <div class="time-gutter"></div>
        <div v-for="day in weekDays" :key="day.date.toISOString()" class="week-day-header" :class="{ today: day.isToday }">
          <span class="weekday-name">{{ day.weekday }}</span>
          <span class="weekday-date">{{ day.date.getDate() }}</span>
        </div>
      </div>
      <div class="week-body">
        <div class="time-gutter">
          <div v-for="hour in 24" :key="hour" class="time-slot">{{ hour - 1 }}:00</div>
        </div>
        <div v-for="day in weekDays" :key="day.date.toISOString()" class="week-day-column">
          <div v-for="hour in 24" :key="hour" class="hour-slot" @click="selectDateTime(day.date, hour - 1)"></div>
          <div
            v-for="(event, idx) in day.events"
            :key="idx"
            class="week-event"
            :style="{ backgroundColor: event.color, top: '0px', height: '30px' }"
            @click="showEventDetail(event)"
          >
            {{ event.title }}
          </div>
        </div>
      </div>
    </div>

    <div class="day-view" v-if="viewMode === 'day'">
      <div class="day-header">
        <span class="day-title">{{ currentDate.toLocaleDateString('zh-CN', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}</span>
      </div>
      <div class="day-events-list">
        <div v-if="dayEvents.length === 0" class="no-events">
          <span>📅 今日暂无提醒</span>
        </div>
        <div
          v-for="(event, idx) in dayEvents"
          :key="idx"
          class="day-event-card"
          :style="{ borderLeftColor: event.color }"
          @click="showEventDetail(event)"
        >
          <div class="event-time">{{ event.start }}</div>
          <div class="event-title">{{ event.title }}</div>
          <div class="event-info" v-if="event.remark">{{ event.remark }}</div>
        </div>
      </div>
    </div>

    <div class="event-detail-overlay" v-if="showEvent" @click.self="showEvent = false">
      <div class="event-detail-card">
        <div class="event-detail-header" :style="{ backgroundColor: selectedEvent?.color }">
          <span class="event-type-icon">{{ getEventIcon(selectedEvent?.type) }}</span>
          <span class="event-type-text">{{ getEventTypeName(selectedEvent?.type) }}</span>
        </div>
        <div class="event-detail-body">
          <h3>{{ selectedEvent?.title }}</h3>
          <div class="detail-row" v-if="selectedEvent?.start">
            <span class="detail-label">📅 时间</span>
            <span class="detail-value">{{ selectedEvent.start }}</span>
          </div>
          <div class="detail-row" v-if="selectedEvent?.status">
            <span class="detail-label">🏷️ 状态</span>
            <span class="detail-value">{{ selectedEvent.status }}</span>
          </div>
          <div class="detail-row" v-if="selectedEvent?.remark">
            <span class="detail-label">📝 备注</span>
            <span class="detail-value">{{ selectedEvent.remark }}</span>
          </div>
        </div>
        <div class="event-detail-footer">
          <button class="btn-secondary" @click="showEvent = false">关闭</button>
          <button v-if="selectedEvent?.blog_id" class="btn-primary" @click="goToBlogger(selectedEvent.blog_id)">
            查看博主
          </button>
        </div>
      </div>
    </div>

    <div class="quick-add-overlay" v-if="showQuickAdd" @click.self="showQuickAdd = false">
      <div class="quick-add-card">
        <div class="quick-add-header">
          <h3>📌 快速添加事件</h3>
          <button class="close-btn" @click="showQuickAdd = false">✕</button>
        </div>
        <div class="quick-add-body">
          <div class="form-row">
            <label>标题</label>
            <input v-model="newEvent.title" type="text" placeholder="例如: 跟进张三" />
          </div>
          <div class="form-row">
            <label>日期</label>
            <input v-model="newEvent.date" type="date" />
          </div>
          <div class="form-row">
            <label>类型</label>
            <select v-model="newEvent.type">
              <option value="followup">跟进提醒</option>
              <option value="cooperation_start">合作开始</option>
              <option value="cooperation_end">合作结束</option>
              <option value="custom">自定义</option>
            </select>
          </div>
          <div class="form-row">
            <label>备注</label>
            <textarea v-model="newEvent.remark" placeholder="可选备注..." rows="2"></textarea>
          </div>
        </div>
        <div class="quick-add-footer">
          <button class="btn-secondary" @click="showQuickAdd = false">取消</button>
          <button class="btn-primary" @click="handleQuickAdd">添加</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getCalendarEventsAPI } from '../api'

const router = useRouter()

const viewMode = ref('month')
const currentDate = ref(new Date())
const currentYear = ref(new Date().getFullYear())
const currentMonth = ref(new Date().getMonth())
const events = ref([])
const showEvent = ref(false)
const selectedEvent = ref(null)
const showQuickAdd = ref(false)
const newEvent = reactive({ title: '', date: '', type: 'followup', remark: '' })

const upcomingEvents = computed(() => {
  const now = new Date()
  return events.value
    .filter(e => {
      const d = new Date(e.date)
      const diffDays = Math.ceil((d - now) / (1000 * 60 * 60 * 24))
      return diffDays >= 0 && diffDays <= 14
    })
    .sort((a, b) => new Date(a.date) - new Date(b.date))
    .map(e => ({
      ...e,
      urgency: (() => {
        const d = new Date(e.date)
        const diff = Math.ceil((d - now) / (1000 * 60 * 60 * 24))
        if (diff <= 3) return 'urgent'
        if (diff <= 7) return 'warning'
        return 'normal'
      })()
    }))
})

const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']

const currentTitle = computed(() => {
  if (viewMode.value === 'month') {
    return `${currentYear.value}年 ${currentMonth.value + 1}月`
  } else if (viewMode.value === 'week') {
    const weekStart = getWeekStart(currentDate.value)
    const weekEnd = new Date(weekStart)
    weekEnd.setDate(weekEnd.getDate() + 6)
    return `${weekStart.getMonth() + 1}月${weekStart.getDate()}日 - ${weekEnd.getMonth() + 1}月${weekEnd.getDate()}日`
  } else {
    return currentDate.value.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
  }
})

const calendarDays = computed(() => {
  const year = currentYear.value
  const month = currentMonth.value
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const startDate = new Date(firstDay)
  startDate.setDate(startDate.getDate() - firstDay.getDay())

  const days = []
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  for (let i = 0; i < 42; i++) {
    const date = new Date(startDate)
    date.setDate(startDate.getDate() + i)
    const dateStr = formatDate(date)

    days.push({
      date,
      isCurrentMonth: date.getMonth() === month,
      isToday: date.getTime() === today.getTime(),
      events: events.value.filter(e => e.start && e.start.startsWith(dateStr))
    })
  }
  return days
})

const weekDays = computed(() => {
  const weekStart = getWeekStart(currentDate.value)
  const days = []
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  for (let i = 0; i < 7; i++) {
    const date = new Date(weekStart)
    date.setDate(weekStart.getDate() + i)
    const dateStr = formatDate(date)

    days.push({
      date,
      weekday: weekdays[date.getDay()],
      isToday: date.getTime() === today.getTime(),
      events: events.value.filter(e => e.start && e.start.startsWith(dateStr))
    })
  }
  return days
})

const dayEvents = computed(() => {
  const dateStr = formatDate(currentDate.value)
  return events.value.filter(e => e.start && e.start.startsWith(dateStr))
})

const loadEvents = async () => {
  try {
    const start = getMonthStart()
    const end = getMonthEnd()
    const res = await getCalendarEventsAPI(start, end)
    if (res.code === 200) {
      events.value = res.data || []
    }
  } catch (error) {
    console.error('加载日历事件失败', error)
  }
}

const getMonthStart = () => {
  const date = new Date(currentYear.value, currentMonth.value, 1)
  return formatDate(date)
}

const getMonthEnd = () => {
  const date = new Date(currentYear.value, currentMonth.value + 1, 0)
  return formatDate(date)
}

const getWeekStart = (date) => {
  const d = new Date(date)
  const day = d.getDay()
  d.setDate(d.getDate() - day)
  return d
}

const formatDate = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const previousPeriod = () => {
  if (viewMode.value === 'month') {
    if (currentMonth.value === 0) {
      currentMonth.value = 11
      currentYear.value--
    } else {
      currentMonth.value--
    }
  } else if (viewMode.value === 'week') {
    currentDate.value = new Date(currentDate.value.setDate(currentDate.value.getDate() - 7))
  } else {
    currentDate.value = new Date(currentDate.value.setDate(currentDate.value.getDate() - 1))
  }
}

const nextPeriod = () => {
  if (viewMode.value === 'month') {
    if (currentMonth.value === 11) {
      currentMonth.value = 0
      currentYear.value++
    } else {
      currentMonth.value++
    }
  } else if (viewMode.value === 'week') {
    currentDate.value = new Date(currentDate.value.setDate(currentDate.value.getDate() + 7))
  } else {
    currentDate.value = new Date(currentDate.value.setDate(currentDate.value.getDate() + 1))
  }
}

const goToToday = () => {
  currentDate.value = new Date()
  currentYear.value = new Date().getFullYear()
  currentMonth.value = new Date().getMonth()
}

const handleQuickAdd = async () => {
  if (!newEvent.title || !newEvent.date) { notification.warning('请填写标题和日期'); return }
  try {
    await createCalendarEventAPI({
      title: newEvent.title,
      date: newEvent.date,
      type: newEvent.type,
      remark: newEvent.remark
    })
    notification.success('事件已创建')
    showQuickAdd.value = false
    Object.assign(newEvent, { title: '', date: '', type: 'followup', remark: '' })
    loadEvents()
  } catch (e) {
    notification.error('创建失败')
  }
}

const selectDay = (day) => {
  currentDate.value = day.date
  viewMode.value = 'day'
}

const selectDateTime = (date, hour) => {
  currentDate.value = date
  viewMode.value = 'day'
}

const showEventDetail = (event) => {
  selectedEvent.value = event
  showEvent.value = true
}

const goToBlogger = (blogId) => {
  showEvent.value = false
  router.push(`/blogger/${blogId}`)
}

const getEventIcon = (type) => {
  const icons = {
    'followup': '📅',
    'followup_record': '📋',
    'cooperation_start': '🤝',
    'cooperation_end': '📝'
  }
  return icons[type] || '📌'
}

const getEventTypeName = (type) => {
  const names = {
    'followup': '跟进提醒',
    'followup_record': '跟进记录',
    'cooperation_start': '合作开始',
    'cooperation_end': '合作结束'
  }
  return names[type] || '事件'
}

watch([currentYear, currentMonth, viewMode], () => {
  loadEvents()
}, { immediate: true })
</script>

<style scoped>
.calendar-view {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  overflow: hidden;
  border: 1px solid var(--border-color);
  transition: all 0.3s;
}

.calendar-header {
  padding: 20px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  color: white;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left h2 {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  min-width: 160px;
}

.nav-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.nav-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.nav-btn svg {
  width: 18px;
  height: 18px;
}

.today-btn {
  padding: 8px 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.today-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.view-tabs {
  display: flex;
  gap: 4px;
  background: rgba(255, 255, 255, 0.2);
  padding: 4px;
  border-radius: 8px;
}

.view-tabs button {
  padding: 6px 16px;
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  font-weight: 500;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.view-tabs button.active {
  background: white;
  color: var(--primary);
}

.calendar-grid {
  padding: 16px;
}

.weekday-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
  margin-bottom: 8px;
}

.weekday {
  text-align: center;
  font-size: 13px;
  font-weight: 600;
  color: #6b7280;
  padding: 8px;
}

.days-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
}

.day-cell {
  min-height: 80px;
  background: var(--bg-hover);
  border-radius: var(--radius-md);
  padding: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.day-cell:hover {
  background: var(--bg-card-hover);
  transform: translateY(-2px);
}

.day-cell.other-month {
  opacity: 0.4;
}

.day-cell.today {
  background: var(--primary-50);
  border: 2px solid var(--primary);
}

.day-number {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.day-cell.today .day-number {
  color: var(--primary);
}

.day-events {
  margin-top: 4px;
}

.event-dot {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 4px;
  color: white;
  margin-bottom: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.more-events {
  font-size: 10px;
  color: var(--text-muted);
  padding: 2px;
}

.week-view {
  padding: 16px;
}

.week-header {
  display: grid;
  grid-template-columns: 60px repeat(7, 1fr);
  gap: 4px;
  margin-bottom: 8px;
}

.time-gutter {
  width: 60px;
}

.week-day-header {
  text-align: center;
  padding: 8px;
  background: var(--bg-hover);
  border-radius: var(--radius-md);
}

.week-day-header.today {
  background: var(--primary-50);
  color: var(--primary);
}

.weekday-name {
  display: block;
  font-size: 12px;
  color: var(--text-muted);
}

.weekday-date {
  display: block;
  font-size: 18px;
  font-weight: 600;
}

.week-day-header.today .weekday-date {
  color: var(--primary);
}

.week-body {
  display: grid;
  grid-template-columns: 60px repeat(7, 1fr);
  gap: 4px;
  max-height: 600px;
  overflow-y: auto;
}

.time-slot {
  height: 48px;
  font-size: 11px;
  color: var(--text-muted);
  text-align: right;
  padding-right: 8px;
}

.week-day-column {
  position: relative;
  background: var(--bg-hover);
  border-radius: var(--radius-md);
}

.hour-slot {
  height: 48px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
}

.hour-slot:hover {
  background: var(--bg-card-hover);
}

.week-event {
  position: absolute;
  left: 2px;
  right: 2px;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 11px;
  color: white;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
  z-index: 1;
}

.day-view {
  padding: 16px;
}

.day-header {
  margin-bottom: 16px;
  padding: 16px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border-radius: 12px;
  color: white;
}

.day-title {
  font-size: 18px;
  font-weight: 600;
}

.day-events-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.no-events {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
  font-size: 16px;
}

.day-event-card {
  padding: 16px;
  background: var(--bg-hover);
  border-radius: var(--radius-lg);
  border-left: 4px solid;
  cursor: pointer;
  transition: all 0.2s;
}

.day-event-card:hover {
  transform: translateX(4px);
  box-shadow: var(--shadow-md);
}

.event-time {
  font-size: 12px;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.event-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.event-info {
  font-size: 13px;
  color: var(--text-muted);
}

.event-detail-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(4px);
}

.event-detail-card {
  background: var(--bg-card);
  border-radius: var(--radius-xl);
  width: 90%;
  max-width: 400px;
  overflow: hidden;
  box-shadow: var(--shadow-xl);
  animation: slideUp 0.3s ease;
}

.event-detail-header {
  padding: 20px 24px;
  color: white;
  display: flex;
  align-items: center;
  gap: 12px;
}

.event-type-icon {
  font-size: 24px;
}

.event-type-text {
  font-size: 16px;
  font-weight: 600;
}

.event-detail-body {
  padding: 24px;
}

.event-detail-body h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 16px 0;
  color: var(--text-primary);
}

.detail-row {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.detail-label {
  font-size: 14px;
  color: var(--text-muted);
  min-width: 60px;
}

.detail-value {
  font-size: 14px;
  color: var(--text-primary);
  flex: 1;
}

.event-detail-footer {
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.btn-secondary {
  padding: 10px 20px;
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: var(--bg-card-hover);
}

.btn-primary {
  padding: 10px 24px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.add-event-btn {
  width: 34px;
  height: 34px;
  border: none;
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  color: white;
  border-radius: 10px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.25);
}

.add-event-btn:hover { transform: scale(1.08); }
.add-event-btn svg { width: 18px; height: 18px; }

.upcoming-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(245, 158, 11, 0.05);
  border: 1px solid rgba(245, 158, 11, 0.12);
  border-radius: 12px;
  margin-bottom: 12px;
  overflow-x: auto;
}

.upcoming-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  white-space: nowrap;
}

.upcoming-list {
  display: flex;
  gap: 6px;
  flex-wrap: nowrap;
}

.upcoming-chip {
  padding: 4px 10px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
  background: var(--bg-secondary);
  color: var(--text-secondary);
}

.upcoming-chip:hover { opacity: 0.85; transform: translateY(-1px); }

.upcoming-chip.urgent {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.2);
}
.upcoming-chip.warning {
  background: rgba(245, 158, 11, 0.1);
  color: #d97706;
  border: 1px solid rgba(245, 158, 11, 0.15);
}
.upcoming-chip.normal {
  background: rgba(34, 197, 94, 0.08);
  color: #16a34a;
  border: 1px solid rgba(34, 197, 94, 0.12);
}

.quick-add-overlay {
  position: fixed; inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fade-in 0.2s ease-out;
}

@keyframes fade-in { from { opacity: 0; } to { opacity: 1; } }

.quick-add-card {
  width: 90%; max-width: 420px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 18px;
  padding: 24px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: slide-up 0.25s ease-out;
}

@keyframes slide-up { from { opacity: 0; transform: translateY(20px) scale(0.96); } to { opacity: 1; transform: translateY(0) scale(1); } }

.quick-add-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
}

.quick-add-header h3 { margin: 0; font-size: 17px; color: var(--text-primary); }

.close-btn {
  width: 30px; height: 30px;
  border: none; background: var(--bg-secondary);
  border-radius: 8px; cursor: pointer;
  color: var(--text-muted); font-size: 14px;
  transition: all 0.2s;
}
.close-btn:hover { background: rgba(239, 68, 68, 0.08); color: #ef4444; }

.quick-add-body { display: flex; flex-direction: column; gap: 12px; margin-bottom: 18px; }

.form-row label {
  display: block; font-size: 12px; font-weight: 600;
  color: var(--text-muted); margin-bottom: 4px;
}

.form-row input,
.form-row select,
.form-row textarea {
  width: 100%;
  padding: 9px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  font-size: 14px;
  color: var(--text-primary);
  box-sizing: border-box;
  transition: all 0.2s;
}
.form-row input:focus,
.form-row select:focus,
.form-row textarea:focus {
  outline: none; border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.08);
}

.quick-add-footer {
  display: flex; justify-content: flex-end; gap: 8px;
}

.btn-primary {
  padding: 9px 20px;
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  border: none; border-radius: 10px;
  color: white; font-weight: 600; cursor: pointer;
  font-size: 14px; transition: all 0.2s;
}
.btn-primary:hover { transform: translateY(-1px); box-shadow: 0 4px 12px rgba(249, 115, 22, 0.25); }

.btn-secondary {
  padding: 9px 20px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px; color: var(--text-secondary);
  font-weight: 500; cursor: pointer; font-size: 14px;
  transition: all 0.2s;
}
.btn-secondary:hover { background: var(--bg-card-hover); }

@media (max-width: 768px) {
  .quick-add-card { padding: 18px; }
  .upcoming-bar { padding: 10px 12px; gap: 8px; }
  .upcoming-chip { font-size: 11px; padding: 3px 8px; }
}
</style>
