<template>
  <div class="public-users">
    <div class="page-header">
      <h1>团队成员</h1>
      <p>认识我们的团队成员</p>
    </div>

    <div class="tabs">
      <button
        class="tab-btn"
        :class="{ active: activeTab === 'approved' }"
        @click="activeTab = 'approved'"
      >
        已审核成员 ({{ approvedUsers.length }})
      </button>
      <button
        v-if="isAdmin"
        class="tab-btn"
        :class="{ active: activeTab === 'pending' }"
        @click="activeTab = 'pending'"
      >
        待审核申请 ({{ pendingUsers.length }})
      </button>
    </div>

    <div v-if="activeTab === 'approved'" class="users-grid">
      <div
        v-for="user in approvedUsers"
        :key="user.id"
        class="user-card"
        @click="viewUserDetail(user.username)"
      >
        <div class="user-avatar">
          <img v-if="user.avatar" :src="user.avatar" :alt="user.real_name" v-avatar />
          <div v-else class="avatar-placeholder">
            {{ user.real_name?.charAt(0) || user.username?.charAt(0) || '?' }}
          </div>
        </div>
        <div class="user-info">
          <h3 class="user-name">{{ user.real_name }}</h3>
          <span class="user-role" :class="user.role">
            {{ user.role === 'admin' ? '管理员' : '成员' }}
          </span>
          <p v-if="user.bio" class="user-bio">{{ user.bio }}</p>
          <p class="user-join">加入时间：{{ formatDate(user.create_time) }}</p>
        </div>
        <div v-if="isSuperAdmin && user.username !== 'admin'" class="user-actions">
          <button class="action-btn role" @click.stop="handleSetRole(user)">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
            设置权限
          </button>
          <button class="action-btn deactivate" @click.stop="handleDeactivate(user)">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
              <polyline points="16,17 21,12 16,7"/>
              <line x1="21" y1="12" x2="9" y2="12"/>
            </svg>
            注销
          </button>
          <button v-if="isSuperAdmin" class="action-btn delete" @click.stop="handleDeleteUser(user)">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3,6 5,6 21,6"/>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
            </svg>
            删除
          </button>
        </div>
      </div>
    </div>

    <div v-if="activeTab === 'pending'" class="users-grid">
      <div
        v-for="user in pendingUsers"
        :key="user.id"
        class="user-card pending"
      >
        <div class="user-avatar">
          <img v-if="user.avatar" :src="user.avatar" :alt="user.real_name" v-avatar />
          <div v-else class="avatar-placeholder">
            {{ user.real_name?.charAt(0) || user.username?.charAt(0) || '?' }}
          </div>
        </div>
        <div class="user-info">
          <h3 class="user-name">{{ user.real_name }}</h3>
          <span class="user-role">待审核</span>
          <p class="user-username">@{{ user.username }}</p>
          <p class="user-join">申请时间：{{ formatDate(user.create_time) }}</p>
        </div>
        <div class="user-actions">
          <button class="action-btn approve" @click="handleApprove(user.id)">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20,6 9,17 4,12"/>
            </svg>
            批准
          </button>
          <button class="action-btn reject" @click="handleReject(user.id)">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
            拒绝
          </button>
        </div>
      </div>
    </div>

    <div v-if="currentUsers.length === 0" class="empty-state">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
        <circle cx="9" cy="7" r="4"/>
        <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
        <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
      </svg>
      <p>{{ activeTab === 'approved' ? '暂无团队成员' : '暂无待审核用户' }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'
import { useConfirm } from '../utils/confirm'
import { getPublicUsersAPI, getPendingUsersAPI, approveUserAPI, rejectUserAPI, deactivateUserAPI, deleteUserAPI, setUserRoleAPI } from '../api'

const router = useRouter()
const userStore = useUserStore()
const { success, error: notifyError } = useNotification()
const { confirm, confirmDanger, confirmWarning } = useConfirm()

const activeTab = ref('approved')
const users = ref([])
const pendingUsers = ref([])

const isAdmin = computed(() => userStore.role === 'admin')

const isSuperAdmin = computed(() => userStore.username === 'admin' || userStore.role === 'admin')

const approvedUsers = computed(() => users.value.filter(u => u.status === 'active'))

const currentUsers = computed(() => activeTab.value === 'approved' ? approvedUsers.value : pendingUsers.value)

const loadUsers = async () => {
  try {
    const res = await getPublicUsersAPI()
    if (res.code === 200) {
      users.value = res.data || []
    }
  } catch (error) {
    notifyError('加载用户列表失败')
  }
}

const loadPendingUsers = async () => {
  if (!isAdmin.value) return
  
  try {
    const res = await getPendingUsersAPI()
    if (res.code === 200) {
      pendingUsers.value = res.data || []
    }
  } catch (error) {
    notifyError('加载待审核用户失败')
  }
}

const handleApprove = async (id) => {
  if (!await confirm('确定批准该用户的注册申请？')) return

  try {
    const res = await approveUserAPI(id)
    if (res.code === 200) {
      success('已批准该用户的注册申请')
      loadPendingUsers()
      loadUsers()
    } else {
      notifyError(res.message || '操作失败')
    }
  } catch (error) {
    notifyError('操作失败')
  }
}

const handleReject = async (id) => {
  if (!await confirmDanger('确定拒绝该用户的注册申请？')) return

  try {
    const res = await rejectUserAPI(id)
    if (res.code === 200) {
      success('已拒绝该用户的注册申请')
      loadPendingUsers()
    } else {
      notifyError(res.message || '操作失败')
    }
  } catch (error) {
    notifyError('操作失败')
  }
}

const viewUserDetail = (username) => {
  router.push(`/user/${username}`)
}

const handleSetRole = (user) => {
  const newRole = user.role === 'admin' ? 'user' : 'admin'
  const roleName = newRole === 'admin' ? '管理员' : '普通用户'

  if (!confirmWarning(`确定要将【${user.real_name}】的权限设置为【${roleName}】吗？`)) {
    return
  }

  setUserRoleAPI(user.id, newRole)
    .then(data => {
      if (data.code === 200) {
        success('权限设置成功')
        loadUsers()
      } else {
        notifyError(data.message || '操作失败')
      }
    })
    .catch(error => {
      notifyError('操作失败')
  })
}

const handleDeactivate = async (user) => {
  if (!await confirmDanger(`确定要注销【${user.real_name}】吗？此操作不可恢复！`)) {
    return
  }

  try {
    const res = await deactivateUserAPI(user.id)
    if (res.code === 200) {
      success('已注销该用户')
      loadUsers()
    } else {
      notifyError(res.message || '操作失败')
    }
  } catch (error) {
    notifyError('操作失败')
  }
}

const handleDeleteUser = async (user) => {
  if (!await confirmDanger(`确定要永久删除【${user.real_name}】吗？此操作不可恢复！`)) {
    return
  }

  try {
    const res = await deleteUserAPI(user.id)
    if (res.code === 200) {
      success('已永久删除该用户')
      loadUsers()
    } else {
      notifyError(res.message || '操作失败')
    }
  } catch (error) {
    notifyError('操作失败')
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

onMounted(() => {
  loadUsers()
  if (isAdmin.value) {
    loadPendingUsers()
  }
})
</script>

<style scoped>
.public-users {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.page-header p {
  color: #666;
  font-size: 16px;
}

.tabs {
  display: flex;
  gap: 12px;
  margin-bottom: 32px;
  justify-content: center;
}

.tab-btn {
  padding: 12px 24px;
  background: white;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  color: var(--text-tertiary);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.tab-btn.active {
  background: var(--primary);
  border-color: var(--primary);
  color: white;
}

.users-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
}

.user-card {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  cursor: pointer;
}

.user-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.user-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  margin: 0 auto 16px;
  overflow: hidden;
  background: linear-gradient(135deg, #ff6b35 0%, #f7931e 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 28px;
  font-weight: 600;
  color: #fff;
}

.user-name {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.user-role {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  margin-bottom: 12px;
}

.user-role.admin {
  background: #fff2f0;
  color: #ff4d4f;
}

.user-role.user {
  background: #f6ffed;
  color: #52c41a;
}

.user-username {
  font-size: 14px;
  color: var(--text-tertiary);
  margin-bottom: 8px;
}

.user-bio {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.user-join {
  font-size: 12px;
  color: #999;
}

.user-card.pending {
  cursor: default;
}

.user-card.pending:hover {
  transform: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.user-actions {
  display: flex;
  gap: 8px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 12px;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn svg {
  width: 16px;
  height: 16px;
}

.action-btn.approve {
  background: #52c41a;
  color: white;
}

.action-btn.approve:hover {
  background: #389e0d;
}

.action-btn.reject {
  background: #ff4d4f;
  color: white;
}

.action-btn.reject:hover {
  background: #cf1322;
}

.action-btn.role {
  background: var(--info);
  color: white;
}

.action-btn.role:hover {
  background: #2563eb;
}

.action-btn.deactivate {
  background: var(--danger);
  color: white;
}

.action-btn.deactivate:hover {
  background: var(--danger);
}

.action-btn.delete {
  background: #7c3aed;
  color: white;
}

.action-btn.delete:hover {
  background: #6d28d9;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.empty-state svg {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state p {
  font-size: 14px;
}
</style>
