<template>
  <div class="log-page">
    <div class="page-header">
      <div class="header-content">
        <h1>📜 操作记录</h1>
        <p>所有操作均已记录，公开透明</p>
      </div>
    </div>

    <div class="filters">
      <select v-model="filters.action" @change="loadLogs" class="filter-select">
        <option value="">全部操作</option>
        <option value="注册申请">注册申请</option>
        <option value="找回密码">找回密码</option>
        <option value="审核通过">审核通过</option>
        <option value="审核拒绝">审核拒绝</option>
        <option value="新增博主">新增博主</option>
        <option value="编辑博主">编辑博主</option>
        <option value="删除博主">删除博主</option>
        <option value="登录">登录</option>
      </select>

      <select v-model="filters.username" @change="loadLogs" class="filter-select">
        <option value="">全部用户</option>
        <option v-for="user in users" :key="user" :value="user">
          {{ user }}
        </option>
      </select>
    </div>

    <div class="log-list">
      <div
        v-for="(log, index) in logs"
        :key="log.id"
        class="log-item"
        :style="{ animationDelay: index * 0.03 + 's' }"
        @touchstart="touchStart(index, $event)"
        @touchmove="touchMove(index, $event)"
        @touchend="touchEnd(index)"
        :class="{ 'swiped': swipedLogIndex === index }"
      >
        <div class="log-icon" :class="getIconClass(log.action)">
          <span>{{ getIcon(log.action) }}</span>
        </div>
        <div class="log-content">
          <div class="log-header">
            <span class="log-action">{{ log.action }}</span>
            <span class="log-target" v-if="log.target">: {{ log.target }}</span>
          </div>
          <div class="log-detail" v-if="log.detail">{{ log.detail }}</div>
          <div class="log-meta">
            <span class="log-operator">{{ log.operator }}</span>
            <span class="log-time">{{ formatDateTime(log.create_time) }}</span>
          </div>
        </div>
        <div class="log-actions" v-if="userRole === 'admin'">
          <button class="delete-btn" @click.stop="confirmDelete(log)">
            🗑️
          </button>
        </div>
      </div>

      <div v-if="loading" class="loading">
        <span>加载中...</span>
      </div>

      <div v-if="!loading && logs.length === 0" class="empty">
        <span>暂无记录</span>
      </div>
    </div>

    <div class="pagination" v-if="totalPages > 1">
      <button class="page-btn" :disabled="page <= 1" @click="handlePageChange(page - 1)">
        上一页
      </button>
      <span class="page-info">{{ page }} / {{ totalPages }}</span>
      <button class="page-btn" :disabled="page >= totalPages" @click="handlePageChange(page + 1)">
        下一页
      </button>
    </div>

    <!-- 删除确认对话框 -->
    <div v-if="showDeleteConfirm" class="confirm-dialog">
      <div class="confirm-content">
        <h3>确认删除</h3>
        <p>确定要删除这条操作记录吗？</p>
        <div class="confirm-buttons">
          <button class="cancel-btn" @click="showDeleteConfirm = false">取消</button>
          <button class="confirm-btn" @click="deleteLog">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { logger } from '../utils/logger'
import { useNotification } from '../stores/notification'
import { useUserStore } from '../stores/user'
import { getOperationLog, getPublicUsersAPI, deleteLogAPI } from '../api'

const { error: notifyError } = useNotification()
const userStore = useUserStore()

const logs = ref([])
const users = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const showDeleteConfirm = ref(false)
const logToDelete = ref(null)
const userRole = ref('user')

// 触摸手势相关
const swipedLogIndex = ref(-1)
const touchStartX = ref(0)
const touchStartY = ref(0)

const filters = reactive({
  action: '',
  username: ''
})

const totalPages = computed(() => Math.ceil(total.value / pageSize))

const loadLogs = async () => {
  loading.value = true
  try {
    const res = await getOperationLog({
      page: page.value,
      pageSize,
      action: filters.action,
      username: filters.username
    })
    if (res.code === 200) {
      logs.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载日志失败')
  }
  loading.value = false
}

const loadUsers = async () => {
  try {
    const res = await getPublicUsersAPI()
    if (res.code === 200) {
      users.value = (res.data || []).map(u => u.real_name || u.username)
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载用户失败')
  }
}

// 加载用户角色
const loadUserRole = () => {
  userRole.value = userStore.role || localStorage.getItem('role') || 'user'
}

const handlePageChange = (newPage) => {
  page.value = newPage
  loadLogs()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const getIcon = (action) => {
  const icons = {
    '注册申请': '📝',
    '找回密码': '🔑',
    '审核通过': '✅',
    '审核拒绝': '❌',
    '新增博主': '➕',
    '编辑博主': '✏️',
    '删除博主': '🗑️',
    '登录': '🔓'
  }
  return icons[action] || '📋'
}

const getIconClass = (action) => {
  if (action.includes('通过') || action.includes('新增')) return 'success'
  if (action.includes('拒绝') || action.includes('删除')) return 'danger'
  if (action.includes('登录')) return 'info'
  return 'default'
}

const formatDateTime = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 触摸手势处理
const touchStart = (index, event) => {
  touchStartX.value = event.touches[0].clientX
  touchStartY.value = event.touches[0].clientY
  swipedLogIndex.value = -1
}

const touchMove = (index, event) => {
  const touchX = event.touches[0].clientX
  const touchY = event.touches[0].clientY
  const diffX = touchStartX.value - touchX
  const diffY = Math.abs(touchStartY.value - touchY)
  
  // 只有水平滑动且滑动距离足够大时才触发
  if (diffX > 50 && diffY < 50) {
    swipedLogIndex.value = index
  }
}

const touchEnd = (index) => {
  // 可以在这里添加自动回弹逻辑
}

// 删除日志
const confirmDelete = (log) => {
  logToDelete.value = log
  showDeleteConfirm.value = true
}

const deleteLog = async () => {
  if (!logToDelete.value) return
  
  try {
    const data = await deleteLogAPI(logToDelete.value.id)
    if (data.code === 200) {
      loadLogs()
      showDeleteConfirm.value = false
      logToDelete.value = null
    }
  } catch (error) {
    logger.error('', error)
    notifyError('删除日志失败')
  }
}

onMounted(() => {
  loadUserRole()
  loadLogs()
  loadUsers()
})
</script>

<style scoped>
.log-page {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.page-header {
  margin-bottom: 24px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.page-header p {
  font-size: 14px;
  color: var(--text-muted);
}

.filters {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.filter-select {
  height: 40px;
  padding: 0 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  color: var(--text-primary);
  cursor: pointer;
}

.log-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.log-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  animation: slideUp 0.3s ease forwards;
  opacity: 0;
  transition: all 0.2s;
  position: relative;
  overflow: hidden;
}

.log-item:hover {
  border-color: var(--primary);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.log-item.swiped {
  transform: translateX(-80px);
  transition: transform 0.3s ease;
}

.log-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.log-icon.success {
  background: rgba(34, 197, 94, 0.1);
}

.log-icon.danger {
  background: rgba(239, 68, 68, 0.1);
}

.log-icon.info {
  background: rgba(59, 130, 246, 0.1);
}

.log-icon.default {
  background: rgba(148, 163, 184, 0.1);
}

.log-content {
  flex: 1;
  min-width: 0;
}

.log-header {
  margin-bottom: 4px;
}

.log-action {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.log-target {
  font-size: 15px;
  color: var(--text-secondary);
}

.log-detail {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 6px;
}

.log-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: var(--text-muted);
}

.log-actions {
  position: absolute;
  right: -80px;
  top: 0;
  height: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
  background: var(--bg-card);
  border-left: 1px solid var(--border-color);
}

.delete-btn {
  width: 44px;
  height: 44px;
  border: none;
  border-radius: 8px;
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger);
  font-size: 20px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.delete-btn:hover {
  background: rgba(239, 68, 68, 0.2);
  transform: scale(1.05);
}

.loading, .empty {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
  font-size: 14px;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 24px;
}

.page-btn {
  padding: 8px 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  background: var(--bg-card-hover);
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: var(--text-muted);
}

/* 确认对话框 */
.confirm-dialog {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

.confirm-content {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  padding: 24px;
  max-width: 400px;
  width: 90%;
  animation: slideUp 0.3s ease;
}

.confirm-content h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
  text-align: center;
}

.confirm-content p {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 24px;
  text-align: center;
}

.confirm-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.cancel-btn, .confirm-btn {
  padding: 8px 24px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.cancel-btn {
  background: var(--bg-card);
  color: var(--text-primary);
}

.cancel-btn:hover {
  background: var(--bg-card-hover);
}

.confirm-btn {
  background: var(--danger);
  color: var(--bg-card);
  border-color: var(--danger);
}

.confirm-btn:hover {
  background: var(--danger);
}

@media (max-width: 640px) {
  .filters {
    flex-direction: column;
  }

  .log-item {
    flex-direction: column;
    gap: 12px;
  }

  .log-icon {
    width: 36px;
    height: 36px;
    font-size: 16px;
  }

  .log-actions {
    right: -70px;
  }

  .log-item.swiped {
    transform: translateX(-70px);
  }

  .delete-btn {
    width: 36px;
    height: 36px;
    font-size: 16px;
  }
}
</style>
