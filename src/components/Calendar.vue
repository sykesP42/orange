<template>
  <div class="calendar-view">
    <div class="calendar-header">
      <div class="header-left">
        <button class="nav-btn" @click="previousPeriod">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15 18 9 12 15 6"/>
          </svg>
        </button>
        <div class="title-section">
          <h2>{{ currentTitle }}</h2>
          <button v-if="!isViewingToday" class="today-btn" @click="goToToday" title="跳转到今天">
            ← 今天
          </button>
          <span v-else class="today-badge">今天</span>
        </div>
        <button class="nav-btn" @click="nextPeriod">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
      </div>
      <div class="header-right">
        <div class="view-tabs">
          <button :class="{ active: viewMode === 'month' }" @click="switchView('month')">月</button>
          <button :class="{ active: viewMode === 'week' }" @click="switchView('week')">周</button>
          <button :class="{ active: viewMode === 'day' }" @click="switchView('day')">日</button>
        </div>
      </div>
    </div>

    <div class="calendar-body">
      <!-- 加载状态 -->
      <div v-if="isLoading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <!-- 月视图 -->
      <div class="month-view" v-show="viewMode === 'month' && !isLoading">
        <div class="weekday-header">
          <div v-for="day in weekdays" :key="day" class="weekday">
            <span class="weekday-name">{{ day }}</span>
          </div>
        </div>
        <div class="days-grid">
          <div
            v-for="(day, index) in calendarDays"
            :key="index"
            class="day-cell"
            :class="{
              'other-month': !day.isCurrentMonth,
              'today': day.isToday,
              'selected': day.isSelected,
              'has-events': day.events.length > 0
            }"
            @click="selectDay(day)"
          >
            <div class="day-content">
              <span class="day-number">{{ day.date.getDate() }}</span>
              <div class="event-list-mini" v-if="day.events.length > 0">
                <div
                  v-for="(event, idx) in day.events.slice(0, 3)"
                  :key="idx"
                  class="event-item-mini"
                  :style="{ borderLeftColor: event.color || '#f97316' }"
                  @click.stop="selectDay(day)"
                  :title="event.title"
                >
                  <span class="event-type-icon-mini">{{ getEventIcon(event.type) }}</span>
                  <span class="event-title-mini">{{ event.title.length > 8 ? event.title.substring(0, 8) + '...' : event.title }}</span>
                </div>
                <div v-if="day.events.length > 3" class="more-events-mini" @click.stop="selectDay(day)">
                  +{{ day.events.length - 3 }} 更多
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 周视图 -->
      <div class="week-view" v-show="viewMode === 'week' && !isLoading">
        <div class="week-header">
          <div class="time-column-header"></div>
          <div v-for="day in weekDays" :key="day.date.toISOString()" 
               class="week-day-header" 
               :class="{ today: day.isToday, selected: day.isSelected }"
               @click="selectDay(day)">
            <span class="weekday-name">{{ day.weekday }}</span>
            <span class="weekday-date">{{ day.date.getDate() }}</span>
            <div class="event-count-badge" v-if="day.events.length > 0">{{ day.events.length }}</div>
          </div>
        </div>
        <div class="week-body">
          <div class="time-column">
            <div v-for="hour in timeSlots" :key="hour" class="time-slot">{{ hour }}:00</div>
          </div>
          <div v-for="day in weekDays" :key="'col-' + day.date.toISOString()" 
               class="week-day-column"
               :class="{ today: day.isToday }">
            <div v-for="hour in timeSlots" :key="hour" 
                 class="hour-slot"
                 @click="selectDateTime(day.date, hour)"></div>
            <div
              v-for="(event, idx) in getDayEventsForWeek(day)"
              :key="'evt-' + idx"
              class="week-event"
              :style="getEventStyle(event)"
              @click.stop="showEventDetail(event)"
            >
              {{ event.title }}
            </div>
          </div>
        </div>
      </div>

      <!-- 日视图 -->
      <div class="day-view" v-show="viewMode === 'day' && !isLoading">
        <div class="day-view-header">
          <button class="back-to-month-btn" @click="switchView('month')">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="15 18 9 12 15 6"/>
            </svg>
            返回月视图
          </button>
          <h3 class="day-title">{{ dayViewTitle }}</h3>
          <div class="day-stats">
            <span class="stat-item" v-if="dayEvents.length > 0">
              📅 {{ dayEvents.length }} 个提醒
            </span>
          </div>
        </div>

        <div class="day-events-container">
          <!-- 空状态 -->
          <div v-if="dayEvents.length === 0" class="empty-state">
            <div class="empty-icon">📅</div>
            <h4>今日暂无提醒事项</h4>
            <p>这一天很清闲，或者还没有添加任何跟进计划</p>
            <button class="add-event-btn" @click="showAddEventModal = true">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              添加提醒
            </button>
          </div>

          <!-- 事件列表 -->
          <div v-else class="events-list">
            <div class="time-column-mini">
              <div v-for="(event, idx) in dayEvents" :key="idx" class="time-label">
                {{ formatEventTime(event.start) }}
              </div>
            </div>
            <div class="events-column">
              <div
                v-for="(event, idx) in dayEvents"
                :key="idx"
                class="event-card"
                :style="{ borderLeftColor: event.color || '#f97316' }"
                @click="showEventDetail(event)"
              >
                <div class="event-card-header">
                  <span class="event-type-badge" :style="{ backgroundColor: event.color || '#f97316' }">
                    {{ getEventIcon(event.type) }}
                  </span>
                  <span class="event-status" :class="getEventStatus(event)">{{ getStatusText(event.status, event.start) }}</span>
                  <button class="delete-event-btn" @click.stop="deleteEvent(event)" title="删除提醒">🗑️</button>
                </div>
                <h4 class="event-title">{{ event.title }}</h4>
                <p class="event-desc" v-if="event.remark">{{ event.remark }}</p>
                <div class="event-footer">
                  <span class="event-time-full" v-if="event.start">{{ formatFullTime(event.start) }}</span>
                  <button class="view-blogger-btn" v-if="event.blog_id" @click.stop="goToBlogger(event.blog_id)">
                    查看博主 →
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="isLoading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>
    </div>

    <!-- 事件详情弹窗 -->
    <Transition name="modal">
      <div class="event-detail-overlay" v-if="showEvent" @click.self="showEvent = false">
        <div class="event-detail-card">
          <div class="event-detail-header" :style="{ background: `linear-gradient(135deg, ${selectedEvent?.color || '#f97316'}, ${selectedEvent?.color || '#f97316'}dd)` }">
            <div class="detail-header-content">
              <span class="event-type-icon">{{ getEventIcon(selectedEvent?.type) }}</span>
              <div>
                <h3>{{ selectedEvent?.title }}</h3>
                <span class="event-type-text">{{ getEventTypeName(selectedEvent?.type) }}</span>
              </div>
            </div>
            <button class="close-btn" @click="showEvent = false">×</button>
          </div>
          <div class="event-detail-body">
            <div class="detail-section">
              <h4>📅 时间信息</h4>
              <div class="detail-row" v-if="selectedEvent?.start">
                <span class="detail-label">开始时间</span>
                <span class="detail-value">{{ formatFullTime(selectedEvent.start) }}</span>
              </div>
              <div class="detail-row" v-if="selectedEvent?.end">
                <span class="detail-label">结束时间</span>
                <span class="detail-value">{{ formatFullTime(selectedEvent.end) }}</span>
              </div>
            </div>
            <div class="detail-section">
              <h4>🏷️ 状态与备注</h4>
              <div class="detail-row" v-if="selectedEvent?.start">
                <span class="detail-label">当前状态</span>
                <span class="status-badge" :class="getEventStatus(selectedEvent)">{{ getStatusText(selectedEvent?.status, selectedEvent?.start) }}</span>
              </div>
              <div class="detail-row" v-if="selectedEvent?.remark">
                <span class="detail-label">备注说明</span>
                <span class="detail-value remark">{{ selectedEvent.remark }}</span>
              </div>
            </div>
          </div>
          <div class="event-detail-footer">
            <button class="btn-danger" @click="deleteSelectedEvent" v-if="selectedEvent?.id">🗑️ 删除提醒</button>
            <button class="btn-secondary" @click="showEvent = false">关闭</button>
            <button v-if="selectedEvent?.blog_id" class="btn-primary" @click="goToBlogger(selectedEvent.blog_id)">
              查看博主详情
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- 添加提醒弹窗 -->
    <Transition name="modal">
      <div class="event-detail-overlay" v-if="showAddEventModal" @click.self="closeAddEventModal">
        <div class="event-detail-card add-event-modal">
          <div class="event-detail-header" style="background: linear-gradient(135deg, #f97316 0%, #fb923c 100%);">
            <div class="detail-header-content">
              <span class="event-type-icon">📝</span>
              <h3>添加提醒事项</h3>
            </div>
            <button class="close-btn" @click="closeAddEventModal">×</button>
          </div>
          <div class="event-detail-body">
            <div class="form-group">
              <label>提醒标题 *</label>
              <input type="text" v-model="newEvent.title" placeholder="例如：联系小红书博主张三" maxlength="50">
            </div>
            <div class="form-row">
              <div class="form-group half">
                <label>提醒类型</label>
                <select v-model="newEvent.type">
                  <option value="followup">跟进提醒</option>
                  <option value="cooperation_start">合作开始</option>
                  <option value="cooperation_end">合作结束</option>
                  <option value="other">其他事项</option>
                </select>
              </div>
              <div class="form-group half">
                <label>提醒时间</label>
                <input type="datetime-local" v-model="newEvent.start">
              </div>
            </div>
            <div class="form-group">
              <label>备注说明</label>
              <textarea v-model="newEvent.remark" placeholder="可选，添加详细备注..." rows="3"></textarea>
            </div>
          </div>
          <div class="event-detail-footer">
            <button class="btn-secondary" @click="closeAddEventModal">取消</button>
            <button class="btn-primary" @click="saveNewEvent" :disabled="isSaving || !newEvent.title.trim()">
              {{ isSaving ? '保存中...' : '保存提醒' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
  
  <ConfirmDialog />
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { getCalendarEventsAPI, createCalendarEventAPI, deleteCalendarEventAPI } from '../api'
import ConfirmDialog from './ConfirmDialog.vue'
import { useConfirm } from '../utils/confirm'

const router = useRouter()
const { confirmDanger } = useConfirm()

const viewMode = ref('month')
const currentDate = ref(new Date())
const currentYear = ref(new Date().getFullYear())
const currentMonth = ref(new Date().getMonth())
const events = ref([])
const showEvent = ref(false)
const showAddEventModal = ref(false)
const selectedEvent = ref(null)
const isLoading = ref(false)
const selectedDateStr = ref('')
const isSaving = ref(false)

const newEvent = ref({
  title: '',
  type: 'followup',
  remark: '',
  start: '',
  color: '#f97316'
})

const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
const timeSlots = Array.from({ length: 12 }, (_, i) => i + 8)

const isViewingToday = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  
  if (viewMode.value === 'day') {
    const viewing = new Date(currentDate.value)
    viewing.setHours(0, 0, 0, 0)
    return viewing.getTime() === today.getTime()
  } else if (viewMode.value === 'month') {
    return currentYear.value === today.getFullYear() && 
           currentMonth.value === today.getMonth()
  } else {
    return false
  }
})

const currentTitle = computed(() => {
  if (viewMode.value === 'month') {
    return `${currentYear.value}年${currentMonth.value + 1}月`
  } else if (viewMode.value === 'week') {
    const weekStart = getWeekStart(currentDate.value)
    const weekEnd = new Date(weekStart)
    weekEnd.setDate(weekEnd.getDate() + 6)
    return `${weekStart.getMonth() + 1}月${weekStart.getDate()}日 - ${weekEnd.getMonth() + 1}月${weekEnd.getDate()}日`
  } else {
    return currentDate.value.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
  }
})

const dayViewTitle = computed(() => {
  return currentDate.value.toLocaleDateString('zh-CN', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

const calendarDays = computed(() => {
  const year = currentYear.value
  const month = currentMonth.value
  const firstDay = new Date(year, month, 1)
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
      isSelected: dateStr === selectedDateStr.value,
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

    days.push({
      date,
      weekday: weekdays[date.getDay()],
      isToday: date.getTime() === today.getTime(),
      isSelected: formatDate(date) === selectedDateStr.value,
      events: events.value.filter(e => e.start && e.start.startsWith(formatDate(date)))
    })
  }
  return days
})

const dayEvents = computed(() => {
  const dateStr = formatDate(currentDate.value)
  return events.value.filter(e => e.start && e.start.startsWith(dateStr))
})

const loadEvents = async () => {
  isLoading.value = true
  try {
    const start = getMonthStart()
    const end = getMonthEnd()
    const res = await getCalendarEventsAPI(start, end)
    if (res.code === 200) {
      events.value = res.data || []
    }
  } catch (error) {
    console.error('加载日历事件失败', error)
  } finally {
    isLoading.value = false
  }
}

const switchView = (mode) => {
  viewMode.value = mode
  if (mode === 'month') {
    loadEvents()
  }
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
  const today = new Date()
  currentDate.value = today
  currentYear.value = today.getFullYear()
  currentMonth.value = today.getMonth()
  selectedDateStr.value = formatDate(today)
}

const selectDay = (day) => {
  currentDate.value = day.date
  selectedDateStr.value = formatDate(day.date)
  viewMode.value = 'day'
}

const selectDateTime = (date, hour) => {
  currentDate.value = date
  selectedDateStr.value = formatDate(date)
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

const openAddEventModal = async () => {
  const dateStr = selectedDateStr.value || formatDate(currentDate.value)
  const now = new Date()
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  
  newEvent.value = {
    title: '',
    type: 'followup',
    remark: '',
    start: `${dateStr}T${hours}:${minutes}`,
    color: '#f97316'
  }
  
  showAddEventModal.value = true
  
  await nextTick()
}

const closeAddEventModal = () => {
  showAddEventModal.value = false
}

const saveNewEvent = async () => {
  if (!newEvent.value.title.trim()) {
    return
  }
  
  isSaving.value = true
  try {
    await createCalendarEventAPI({
      title: newEvent.value.title,
      type: newEvent.value.type,
      start: newEvent.value.start,
      end: newEvent.value.start,
      remark: newEvent.value.remark,
      blog_id: null
    })
    
    showAddEventModal.value = false
    
    await loadEvents()
  } catch (error) {
    console.error('保存提醒失败', error)
  } finally {
    isSaving.value = false
  }
}

const deleteEvent = async (event) => {
  const confirmed = await confirmDanger(`确定要删除提醒「${event.title}」吗？此操作不可恢复！`)
  if (!confirmed) {
    return
  }
  
  try {
    await deleteCalendarEventAPI(event.id)
    events.value = events.value.filter(e => e.id !== event.id)
    if (showEvent.value && selectedEvent.value?.id === event.id) {
      showEvent.value = false
    }
  } catch (error) {
    console.error('删除提醒失败', error)
  }
}

const deleteSelectedEvent = async () => {
  if (!selectedEvent.value?.id) return
  
  const confirmed = await confirmDanger(`确定要删除提醒「${selectedEvent.value.title}」吗？此操作不可恢复！`)
  if (!confirmed) {
    return
  }
  
  try {
    await deleteCalendarEventAPI(selectedEvent.value.id)
    events.value = events.value.filter(e => e.id !== selectedEvent.value.id)
    showEvent.value = false
  } catch (error) {
    console.error('删除提醒失败', error)
  }
}

const getDayEventsForWeek = (day) => {
  return day.events || []
}

const getEventStyle = (event) => {
  const startHour = parseInt(event.start?.split(' ')[1]?.split(':')[0] || '9')
  const top = Math.max(0, (startHour - 8) * 48)
  return {
    backgroundColor: event.color || '#f97316',
    top: `${top}px`,
    height: '40px'
  }
}

const formatEventTime = (dateStr) => {
  if (!dateStr) return ''
  const parts = dateStr.split(' ')
  return parts[1] ? parts[1].substring(0, 5) : ''
}

const formatFullTime = (dateStr) => {
  if (!dateStr) return '-'
  try {
    const [datePart, timePart] = dateStr.split(' ')
    if (datePart && timePart) {
      return `${datePart} ${timePart.substring(0, 5)}`
    }
    return dateStr
  } catch {
    return dateStr
  }
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
  return names[type] || '日程事件'
}

const getStatusText = (status, startTime) => {
  if (startTime) {
    const eventTime = new Date(startTime.replace(' ', 'T'))
    const now = new Date()
    if (eventTime < now) {
      return '已过期'
    }
  }
  const texts = {
    'pending': '待处理',
    'completed': '已完成',
    'cancelled': '已取消'
  }
  return texts[status] || status || '进行中'
}

const getEventStatus = (event) => {
  if (event.start) {
    const eventTime = new Date(event.start.replace(' ', 'T'))
    const now = new Date()
    if (eventTime < now) {
      return 'expired'
    }
  }
  return event.status || 'pending'
}

const getMonthStart = () => {
  return formatDate(new Date(currentYear.value, currentMonth.value, 1))
}

const getMonthEnd = () => {
  return formatDate(new Date(currentYear.value, currentMonth.value + 1, 0))
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

watch([currentYear, currentMonth], () => {
  if (viewMode.value === 'month') {
    loadEvents()
  }
}, { immediate: true })

onMounted(() => {
  selectedDateStr.value = formatDate(new Date())
})
</script>

<style scoped>
.calendar-view {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Header */
.calendar-header {
  padding: 20px 28px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #e5e7eb;
  background: linear-gradient(135deg, #f97316 0%, #fb923c 100%);
  color: white;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left h2 {
  font-size: 20px;
  font-weight: 700;
  margin: 0;
  min-width: 140px;
  letter-spacing: 0.5px;
}

.nav-btn {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.25);
  border: none;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(10px);
}

.nav-btn:hover {
  background: rgba(255, 255, 255, 0.4);
  transform: scale(1.05);
}

.nav-btn:active {
  transform: scale(0.95);
}

.nav-btn svg {
  width: 20px;
  height: 20px;
}

.today-btn {
  padding: 8px 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.25);
  border: 1px solid rgba(255, 255, 255, 0.4);
  color: white;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.25s;
  backdrop-filter: blur(10px);
}

.today-btn:hover {
  background: white;
  color: #f97316;
  border-color: white;
}

.today-badge {
  padding: 6px 14px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.35);
  color: white;
  font-size: 13px;
  font-weight: 700;
  border: 2px solid rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(10px);
  animation: pulse-badge 2s ease-in-out infinite;
}

@keyframes pulse-badge {
  0%, 100% { 
    opacity: 1; 
    transform: scale(1);
  }
  50% { 
    opacity: 0.85; 
    transform: scale(1.02);
  }
}

.view-tabs {
  display: flex;
  gap: 4px;
  background: rgba(255, 255, 255, 0.2);
  padding: 4px;
  border-radius: 10px;
  backdrop-filter: blur(10px);
}

.view-tabs button {
  padding: 8px 18px;
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.85);
  font-size: 14px;
  font-weight: 600;
  border-radius: 7px;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.view-tabs button:hover {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.view-tabs button.active {
  background: white;
  color: #f97316;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

/* Calendar Body */
.calendar-body {
  min-height: 500px;
  position: relative;
}

/* Month View */
.month-view {
  padding: 20px;
}

.weekday-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 8px;
  margin-bottom: 12px;
  padding: 0 4px;
}

.weekday {
  text-align: center;
  padding: 10px 0;
}

.weekday-name {
  font-size: 13px;
  font-weight: 700;
  color: #6b7280;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.days-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 8px;
}

.day-cell {
  min-height: 110px;
  background: #fafbfc;
  border-radius: 12px;
  padding: 8px;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  border: 2px solid transparent;
  position: relative;
  overflow: hidden;
}

.day-cell:hover {
  background: #f3f4f6;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border-color: #f9731633;
}

.day-cell.other-month {
  opacity: 0.35;
  background: #f9fafb;
}

.day-cell.other-month:hover {
  opacity: 0.55;
}

.day-cell.today {
  background: linear-gradient(135deg, #fff7ed 0%, #ffedd5 100%);
  border-color: #f97316;
  box-shadow: 0 0 0 2px #f9731633;
}

.day-cell.selected:not(.today) {
  background: linear-gradient(135deg, #eff6ff 0%, #dbeafe 100%);
  border-color: #3b82f6;
}

.day-cell.has-events::after {
  content: '';
  position: absolute;
  bottom: 6px;
  left: 50%;
  transform: translateX(-50%);
  width: 6px;
  height: 6px;
  background: #f97316;
  border-radius: 50%;
}

.day-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.day-number {
  font-size: 14px;
  font-weight: 700;
  color: #374151;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: all 0.2s;
}

.day-cell.today .day-number {
  background: #f97316;
  color: white;
  box-shadow: 0 2px 6px rgba(249, 115, 22, 0.35);
}

.event-indicators {
  display: flex;
  gap: 3px;
  margin-top: 6px;
  flex-wrap: wrap;
  align-items: center;
}

.event-indicator {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.more-indicator {
  font-size: 10px;
  color: #6b7280;
  font-weight: 600;
  margin-left: 2px;
}

/* Week View */
.week-view {
  padding: 16px;
}

.week-header {
  display: grid;
  grid-template-columns: 60px repeat(7, 1fr);
  gap: 6px;
  margin-bottom: 8px;
}

.time-column-header {
  padding: 8px;
}

.week-day-header {
  text-align: center;
  padding: 10px 8px;
  background: #f9fafb;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
  border: 2px solid transparent;
}

.week-day-header:hover {
  background: #eff6ff;
}

.week-day-header.today {
  background: linear-gradient(135deg, #fff7ed 0%, #ffedd5 100%);
  border-color: #f97316;
}

.week-day-header.selected {
  border-color: #3b82f6;
  background: #eff6ff;
}

.weekday-name {
  display: block;
  font-size: 12px;
  color: #6b7280;
  font-weight: 600;
  margin-bottom: 4px;
}

.weekday-date {
  display: block;
  font-size: 20px;
  font-weight: 700;
  color: #374151;
}

.week-day-header.today .weekday-date {
  color: #f97316;
}

.event-count-badge {
  position: absolute;
  top: 4px;
  right: 4px;
  background: #f97316;
  color: white;
  font-size: 10px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
}

.week-body {
  display: grid;
  grid-template-columns: 60px repeat(7, 1fr);
  gap: 6px;
  max-height: 520px;
  overflow-y: auto;
  border-left: 1px solid #e5e7eb;
}

.time-column {
  padding-right: 4px;
}

.time-slot {
  height: 48px;
  font-size: 11px;
  color: #9ca3af;
  text-align: right;
  padding-right: 8px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.week-day-column {
  position: relative;
  background: #fafbfc;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.hour-slot {
  height: 48px;
  border-bottom: 1px solid #f3f4f6;
  cursor: pointer;
  transition: background 0.15s;
}

.hour-slot:hover {
  background: #eff6ff;
}

.week-event {
  position: absolute;
  left: 3px;
  right: 3px;
  padding: 6px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  color: white;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
  z-index: 2;
  transition: all 0.2s;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
}

.week-event:hover {
  transform: scale(1.02);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  z-index: 3;
}

/* Day View */
.day-view {
  padding: 24px;
  min-height: 480px;
}

.day-view-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 2px solid #f3f4f6;
}

.back-to-month-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
}

.back-to-month-btn:hover {
  background: #e5e7eb;
  color: #374151;
  border-color: #d1d5db;
}

.back-to-month-btn svg {
  width: 16px;
  height: 16px;
}

.day-title {
  font-size: 20px;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
  flex: 1;
}

.day-stats {
  display: flex;
  gap: 12px;
}

.stat-item {
  font-size: 13px;
  color: #6b7280;
  padding: 6px 12px;
  background: #f9fafb;
  border-radius: 8px;
  font-weight: 500;
}

.day-events-container {
  min-height: 320px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  animation: fadeInUp 0.5s ease;
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  filter: grayscale(30%);
}

.empty-state h4 {
  font-size: 20px;
  font-weight: 700;
  color: #374151;
  margin: 0 0 8px 0;
}

.empty-state p {
  font-size: 14px;
  color: #9ca3af;
  margin: 0 0 24px 0;
}

.add-event-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #f97316 0%, #fb923c 100%);
  border: none;
  border-radius: 10px;
  color: white;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.25s;
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.add-event-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(249, 115, 22, 0.4);
}

.add-event-btn svg {
  width: 18px;
  height: 18px;
}

.events-list {
  display: grid;
  grid-template-columns: 80px 1fr;
  gap: 12px;
  animation: fadeInUp 0.4s ease;
}

.time-column-mini {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.time-label {
  font-size: 13px;
  font-weight: 600;
  color: #6b7280;
  padding-top: 16px;
  text-align: right;
  padding-right: 12px;
}

.events-column {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.event-card {
  padding: 18px 20px;
  background: white;
  border-radius: 14px;
  border-left: 4px solid;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border: 1px solid #f3f4f6;
}

.event-card:hover {
  transform: translateX(6px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
  border-color: #e5e7eb;
}

.event-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.event-type-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  font-size: 16px;
}

.event-status {
  font-size: 11px;
  font-weight: 700;
  padding: 4px 10px;
  border-radius: 6px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.event-status.pending {
  background: #fef3c7;
  color: #d97706;
}

.event-status.completed {
  background: #d1fae5;
  color: #059669;
}

.event-status.cancelled {
	background: #fee2e2;
	color: #dc2626;
}

.event-status.expired {
	background: #f3f4f6;
	color: #9ca3af;
}

.event-title {
  font-size: 16px;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 6px 0;
}

.event-desc {
  font-size: 13px;
  color: #6b7280;
  margin: 0 0 12px 0;
  line-height: 1.5;
}

.event-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 12px;
  border-top: 1px solid #f3f4f6;
}

.event-time-full {
  font-size: 12px;
  color: #9ca3af;
  font-weight: 500;
}

.view-blogger-btn {
  font-size: 12px;
  font-weight: 600;
  color: #f97316;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s;
}

.view-blogger-btn:hover {
  background: #fff7ed;
}

/* Loading State */
.loading-state {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: white;
  z-index: 10;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #f3f4f6;
  border-top-color: #f97316;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-state p {
  margin-top: 16px;
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

/* Event Detail Modal */
.event-detail-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(8px);
  padding: 20px;
}

.event-detail-card {
  background: white;
  border-radius: 20px;
  width: 100%;
  max-width: 480px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.3);
}

.event-detail-header {
  padding: 28px 28px 24px;
  color: white;
  position: relative;
}

.detail-header-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.detail-header-content h3 {
  font-size: 20px;
  font-weight: 700;
  margin: 0 0 4px 0;
}

.event-type-icon {
  font-size: 36px;
}

.event-type-text {
  font-size: 14px;
  opacity: 0.9;
  font-weight: 500;
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.25);
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  line-height: 1;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.4);
  transform: rotate(90deg);
}

.event-detail-body {
  padding: 28px;
}

.detail-section {
  margin-bottom: 24px;
}

.detail-section:last-child {
  margin-bottom: 0;
}

.detail-section h4 {
  font-size: 14px;
  font-weight: 700;
  color: #6b7280;
  margin: 0 0 14px 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f3f4f6;
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-label {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

.detail-value {
  font-size: 14px;
  color: #1f2937;
  font-weight: 600;
  text-align: right;
  max-width: 60%;
}

.detail-value.remark {
  text-align: left;
  max-width: 100%;
  line-height: 1.6;
  color: #4b5563;
  font-weight: 400;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
}

.status-badge.pending {
  background: #fef3c7;
  color: #d97706;
}

.status-badge.completed {
  background: #d1fae5;
  color: #059669;
}

.status-badge.cancelled {
	background: #fee2e2;
	color: #dc2626;
}

.status-badge.expired {
	background: #fee2e2;
	color: #9ca3af;
}

.event-detail-footer {
  padding: 20px 28px;
  border-top: 1px solid #e5e7eb;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  background: #f9fafb;
}

.btn-secondary {
  padding: 10px 20px;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  color: #4b5563;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: #f3f4f6;
  border-color: #9ca3af;
}

.btn-primary {
  padding: 10px 24px;
  background: linear-gradient(135deg, #f97316 0%, #fb923c 100%);
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  color: white;
  cursor: pointer;
  transition: all 0.25s;
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.btn-primary:hover {
	transform: translateY(-2px);
	box-shadow: 0 6px 20px rgba(249, 115, 22, 0.4);
}

.btn-danger {
	padding: 10px 20px;
	background: linear-gradient(135deg, #ef4444 0%, #f87171 100%);
	border: none;
	border-radius: 10px;
	font-size: 14px;
	font-weight: 600;
	color: white;
	cursor: pointer;
	transition: all 0.25s;
	box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
	margin-right: auto;
}

.btn-danger:hover {
	transform: translateY(-2px);
	box-shadow: 0 6px 20px rgba(239, 68, 68, 0.4);
}

.delete-event-btn {
	background: none;
	border: none;
	font-size: 16px;
	cursor: pointer;
	padding: 4px;
	border-radius: 6px;
	transition: all 0.2s;
	margin-left: auto;
}

.delete-event-btn:hover {
	background: #fee2e2;
	transform: scale(1.1);
}

/* Modal Transitions */
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .event-detail-card,
.modal-leave-active .event-detail-card {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.modal-enter-from .event-detail-card,
.modal-leave-to .event-detail-card {
  transform: scale(0.9) translateY(20px);
  opacity: 0;
}

/* Responsive */
@media (max-width: 1024px) {
  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .charts-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .calendar-header {
    padding: 16px 20px;
    flex-wrap: wrap;
    gap: 12px;
  }

  .header-left {
    flex-wrap: wrap;
    justify-content: center;
  }

  .header-left h2 {
    font-size: 17px;
    min-width: auto;
  }

  .day-cell {
    min-height: 64px;
    padding: 6px;
  }

  .day-number {
    font-size: 13px;
    width: 24px;
    height: 24px;
  }

  .event-indicators {
    margin-top: 4px;
  }

  .event-indicator {
    width: 5px;
    height: 5px;
  }

  .week-body {
    max-height: 400px;
  }

  .day-view {
    padding: 16px;
  }

  .day-view-header {
    flex-wrap: wrap;
    gap: 10px;
  }

  .day-title {
    font-size: 17px;
  }

  .events-list {
    grid-template-columns: 60px 1fr;
    gap: 8px;
  }

  .event-card {
    padding: 14px 16px;
  }

  .event-title {
    font-size: 14px;
  }
}

@media (max-width: 480px) {
  .calendar-header {
    padding: 12px 16px;
    flex-direction: column;
    gap: 10px;
  }

  .header-left {
    width: 100%;
    justify-content: space-between;
  }

  .header-right {
    width: 100%;
    justify-content: center;
  }

  .title-section {
    flex: 1;
    justify-content: center;
  }

  .month-view {
    padding: 12px;
  }

  .days-grid {
    gap: 4px;
  }

  .day-cell {
    min-height: 52px;
    padding: 4px;
    border-radius: 8px;
  }

  .day-number {
    font-size: 12px;
    width: 22px;
    height: 22px;
    border-radius: 6px;
  }

  .weekday-header {
    gap: 4px;
    padding: 0 2px;
  }

  .weekday-name {
    font-size: 11px;
  }

  .event-dot {
    width: 4px;
    height: 4px;
  }

  .more-events {
    font-size: 9px;
  }

  .week-header {
    grid-template-columns: 44px repeat(7, 1fr);
  }

  .week-body {
    grid-template-columns: 44px repeat(7, 1fr);
    max-height: 340px;
  }

  .time-slot {
    height: 40px;
    font-size: 10px;
  }

  .hour-slot {
    height: 40px;
  }

  .week-event {
    font-size: 10px;
    padding: 4px 6px;
  }

  .day-view {
    padding: 12px;
  }

  .back-to-month-btn {
    padding: 6px 10px;
    font-size: 12px;
  }

  .day-title {
    font-size: 15px;
  }

  .stat-item {
    font-size: 11px;
    padding: 4px 8px;
  }

  .empty-icon {
    font-size: 48px;
  }

  .empty-state h4 {
    font-size: 17px;
  }

  .empty-state p {
    font-size: 13px;
  }

  .add-event-btn {
    padding: 10px 18px;
    font-size: 13px;
  }

  .events-list {
    grid-template-columns: 50px 1fr;
  }

  .time-label {
    font-size: 11px;
    padding-top: 12px;
  }

  .event-card {
    padding: 12px 14px;
  }

  .event-type-badge {
    width: 28px;
    height: 28px;
    font-size: 14px;
  }

  .event-title {
    font-size: 13px;
  }

  .event-desc {
    font-size: 12px;
  }

  .event-time-full {
    font-size: 11px;
  }

  .view-blogger-btn {
    font-size: 11px;
  }

  .event-detail-card {
    max-width: 100%;
    border-radius: 16px;
  }

  .event-detail-header {
    padding: 20px;
  }

  .detail-header-content h3 {
    font-size: 17px;
  }

  .event-detail-body {
    padding: 20px;
  }

  .event-detail-footer {
    padding: 16px 20px;
    flex-direction: column;
  }

  .btn-primary,
  .btn-secondary {
    width: 100%;
    text-align: center;
  }
}

/* 事件缩略显示 */
.event-list-mini {
  margin-top: 4px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.event-item-mini {
  display: flex;
  align-items: center;
  gap: 3px;
  padding: 2px 5px;
  font-size: 10px;
  background: white;
  border-radius: 4px;
  border-left: 3px solid #f97316;
  cursor: pointer;
  transition: all 0.15s;
  overflow: hidden;
  white-space: nowrap;
}

.event-item-mini:hover {
  background: #fff7ed;
  transform: scale(1.02);
}

.event-type-icon-mini {
  flex-shrink: 0;
  font-size: 9px;
}

.event-title-mini {
  overflow: hidden;
  text-overflow: ellipsis;
  color: #374151;
  font-weight: 500;
}

.more-events-mini {
  font-size: 9px;
  color: #f97316;
  font-weight: 600;
  text-align: center;
  padding: 2px;
  cursor: pointer;
  background: #fff7ed;
  border-radius: 4px;
}

.more-events-mini:hover {
  background: #fed7aa;
}

/* 添加提醒弹窗表单 */
.add-event-modal {
  max-width: 480px;
}

.form-group {
  margin-bottom: 18px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 6px;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.2s;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #f97316;
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

.form-row {
  display: flex;
  gap: 12px;
}

.form-group.half {
  flex: 1;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
