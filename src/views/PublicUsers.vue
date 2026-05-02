<template>
  <div class="public-users">
    <div class="page-header">
      <div class="header-content">
        <h1>用户中心</h1>
        <p>查看全部成员信息，快速联系协作</p>
      </div>
      <div class="header-stats">
        <div class="stat-chip">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
          <span>{{ approvedUsers.length }} 位成员</span>
        </div>
        <div v-if="pendingUsers.length > 0" class="stat-chip pending">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
          <span>{{ pendingUsers.length }} 待审核</span>
        </div>
      </div>
    </div>

    <div class="toolbar">
      <div class="search-box">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <input v-model="searchQuery" id="public-users-search" name="search" type="text" placeholder="搜索姓名、用户名..." />
      </div>
      <div class="filter-group">
        <select v-model="teamFilter" class="team-filter">
          <option value="">全部团队</option>
          <option v-for="team in teams" :key="team.id" :value="team.id">{{ team.name }}</option>
        </select>
        <select v-model="roleFilter" class="role-filter">
          <option value="">全部角色</option>
          <option value="admin">管理员</option>
          <option value="user">成员</option>
        </select>
      </div>
    </div>

    <div class="tabs">
      <button class="tab-btn" :class="{ active: activeTab === 'approved' }" @click="activeTab = 'approved'">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/></svg>
        成员 ({{ approvedUsers.length }})
      </button>
      <button v-if="isAdmin" class="tab-btn" :class="{ active: activeTab === 'pending' }" @click="activeTab = 'pending'">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
        待审核 ({{ pendingUsers.length }})
      </button>
    </div>

    <div v-if="activeTab === 'approved'" class="users-grid">
      <div v-for="user in filteredUsers" :key="user.id" class="user-card" @click="viewUserDetail(user.username)">
        <div class="card-top">
          <div class="user-avatar">
            <img v-if="user.avatar" :src="user.avatar" :alt="user.real_name" />
            <div v-else class="avatar-placeholder">{{ user.real_name?.charAt(0) || '?' }}</div>
          </div>
          <div class="user-badges">
            <span class="badge role" :class="user.role">{{ user.role === 'admin' ? '管理员' : '成员' }}</span>
            <span v-if="getTeamForUser(user)" class="badge team" :style="{ backgroundColor: getTeamForUser(user)?.color || '#6b7280' }">{{ getTeamForUser(user)?.name }}</span>
          </div>
        </div>
        <div class="card-body">
          <h3 class="user-name">{{ user.real_name }}</h3>
          <p class="user-username">@{{ user.username }}</p>
          <p v-if="user.bio" class="user-bio">{{ user.bio }}</p>
        </div>
        <div class="card-footer">
          <span class="join-date">{{ formatDate(user.create_time) }}加入</span>
          <div class="quick-actions" @click.stop>
            <button class="quick-btn" @click="goToChat(user)" title="发私信">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
            </button>
            <button class="quick-btn" @click="viewUserBloggers(user.username)" title="查看博主">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
            </button>
            <button class="quick-btn" @click="viewUserDetail(user.username)" title="查看详情">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9,6 15,12 9,18"/></svg>
            </button>
          </div>
        </div>
        <div v-if="isAdmin && user.username !== 'admin'" class="admin-actions">
          <button class="action-btn role" @click.stop="handleSetRole(user)">
            {{ user.role === 'admin' ? '取消管理员' : '设为管理员' }}
          </button>
          <button class="action-btn deactivate" @click.stop="handleDeactivate(user)">注销</button>
          <button v-if="isAdmin" class="action-btn delete" @click.stop="handleDeleteUser(user)">删除</button>
        </div>
      </div>
    </div>

    <div v-if="activeTab === 'pending'" class="users-grid">
      <div v-for="user in pendingUsers" :key="user.id" class="user-card pending">
        <div class="card-top">
          <div class="user-avatar">
            <img v-if="user.avatar" :src="user.avatar" :alt="user.real_name" />
            <div v-else class="avatar-placeholder">{{ user.real_name?.charAt(0) || '?' }}</div>
          </div>
          <span class="badge pending-badge">待审核</span>
        </div>
        <div class="card-body">
          <h3 class="user-name">{{ user.real_name }}</h3>
          <p class="user-username">@{{ user.username }}</p>
        </div>
        <div class="card-footer">
          <span class="join-date">{{ formatDate(user.create_time) }}申请</span>
          <div class="pending-actions">
            <button class="action-btn approve" id="approve-btn" @click="handleApprove(user.id)">批准</button>
            <button class="action-btn reject" id="reject-btn" @click="handleReject(user.id)">拒绝</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="currentUsers.length === 0 && !loading" class="empty-state">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/>
      </svg>
      <p>{{ searchQuery || teamFilter || roleFilter ? '没有匹配的用户' : (activeTab === 'approved' ? '暂无团队成员' : '暂无待审核用户') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'
import { useConfirm } from '../utils/confirm'
import { getPublicUsersAPI, getPendingUsersAPI, approveUserAPI, rejectUserAPI, deactivateUserAPI, deleteUserAPI, setUserRoleAPI, getTeamsAPI } from '../api'

const router = useRouter()
const userStore = useUserStore()
const { success, error: notifyError } = useNotification()
const { confirm, confirmDanger, confirmWarning } = useConfirm()

const activeTab = ref('approved')
const users = ref([])
const pendingUsers = ref([])
const teams = ref([])
const searchQuery = ref('')
const teamFilter = ref('')
const roleFilter = ref('')
const loading = ref(false)

const isAdmin = computed(() => userStore.role === 'admin')

const approvedUsers = computed(() => users.value.filter(u => u.status === 'active'))

const filteredUsers = computed(() => {
  let result = approvedUsers.value
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(u =>
      (u.real_name || '').toLowerCase().includes(q) ||
      (u.username || '').toLowerCase().includes(q)
    )
  }
  if (teamFilter.value) {
    result = result.filter(u => String(u.team_id) === String(teamFilter.value))
  }
  if (roleFilter.value) {
    result = result.filter(u => u.role === roleFilter.value)
  }
  return result
})

const currentUsers = computed(() => activeTab.value === 'approved' ? filteredUsers.value : pendingUsers.value)

const getTeamForUser = (user) => {
  if (!user.team_id) return null
  return teams.value.find(t => String(t.id) === String(user.team_id))
}

const loadUsers = async () => {
  loading.value = true
  try {
    const res = await getPublicUsersAPI()
    if (res.code === 200) {
      users.value = res.data || []
    }
  } catch (error) {
    notifyError('加载用户列表失败')
  }
  loading.value = false
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

const loadTeams = async () => {
  try {
    const res = await getTeamsAPI()
    if (res.code === 200) {
      teams.value = res.data || []
    }
  } catch (e) {}
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

const viewUserBloggers = (username) => {
  router.push(`/user/${username}`)
}

const goToChat = (user) => {
  router.push(`/chat?userId=${user.id}`)
}

const handleSetRole = async (user) => {
  const newRole = user.role === 'admin' ? 'user' : 'admin'
  const roleName = newRole === 'admin' ? '管理员' : '普通用户'
  if (!await confirmWarning(`确定要将【${user.real_name}】的权限设置为【${roleName}】吗？`)) return
  try {
    const res = await setUserRoleAPI(user.id, newRole)
    if (res.code === 200) {
      success('权限设置成功')
      loadUsers()
    } else {
      notifyError(res.message || '操作失败')
    }
  } catch (error) {
    notifyError('操作失败')
  }
}

const handleDeactivate = async (user) => {
  if (!await confirmDanger(`确定要注销【${user.real_name}】吗？此操作不可恢复！`)) return
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
  if (!await confirmDanger(`确定要永久删除【${user.real_name}】吗？此操作不可恢复！`)) return
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
  loadTeams()
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
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.page-header p {
  color: var(--text-secondary);
  font-size: 14px;
}

.header-stats {
  display: flex;
  gap: 12px;
}

.stat-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 20px;
  font-size: 13px;
  color: var(--text-secondary);
}

.stat-chip svg {
  width: 16px;
  height: 16px;
  color: var(--primary);
}

.stat-chip.pending svg {
  color: #f59e0b;
}

.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.search-box {
  flex: 1;
  min-width: 200px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  transition: border-color 0.2s;
}

.search-box:focus-within {
  border-color: var(--primary);
}

.search-box svg {
  width: 18px;
  height: 18px;
  color: var(--text-muted);
  flex-shrink: 0;
}

.search-box input {
  flex: 1;
  border: none;
  background: none;
  outline: none;
  font-size: 14px;
  color: var(--text-primary);
}

.filter-group {
  display: flex;
  gap: 8px;
}

.team-filter,
.role-filter {
  padding: 10px 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  font-size: 14px;
  color: var(--text-primary);
  outline: none;
  cursor: pointer;
}

.team-filter:focus,
.role-filter:focus {
  border-color: var(--primary);
}

.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn svg {
  width: 16px;
  height: 16px;
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
  gap: 20px;
}

.user-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 14px;
  padding: 20px;
  transition: all 0.3s ease;
  cursor: pointer;
  display: flex;
  flex-direction: column;
}

.user-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  border-color: rgba(249, 115, 22, 0.2);
}

.card-top {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  margin-bottom: 14px;
}

.user-avatar {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  overflow: hidden;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 22px;
  font-weight: 600;
  color: #fff;
}

.user-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding-top: 4px;
}

.badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

.badge.role.admin {
  background: rgba(249, 115, 22, 0.1);
  color: var(--primary);
}

.badge.role.user {
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.badge.team {
  color: white;
}

.badge.pending-badge {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

.card-body {
  flex: 1;
  margin-bottom: 12px;
}

.user-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 2px;
}

.user-username {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 6px;
}

.user-bio {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.join-date {
  font-size: 12px;
  color: var(--text-muted);
}

.quick-actions {
  display: flex;
  gap: 4px;
}

.quick-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-secondary);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.2s;
}

.quick-btn svg {
  width: 16px;
  height: 16px;
}

.quick-btn:hover {
  background: var(--primary);
  color: white;
}

.admin-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.action-btn {
  flex: 1 1 auto;
  min-width: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 7px 8px;
  border: none;
  border-radius: 8px;
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.action-btn.approve {
  background: #22c55e;
  color: white;
}

.action-btn.approve:hover {
  background: #16a34a;
}

.action-btn.reject {
  background: #ef4444;
  color: white;
}

.action-btn.reject:hover {
  background: #dc2626;
}

.action-btn.role {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.action-btn.role:hover {
  background: #3b82f6;
  color: white;
}

.action-btn.deactivate {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.action-btn.deactivate:hover {
  background: #ef4444;
  color: white;
}

.action-btn.delete {
  background: rgba(124, 58, 237, 0.1);
  color: #7c3aed;
}

.action-btn.delete:hover {
  background: #7c3aed;
  color: white;
}

.pending-actions {
  display: flex;
  gap: 6px;
}

.user-card.pending {
  cursor: default;
}

.user-card.pending:hover {
  transform: none;
  box-shadow: none;
  border-color: var(--border-color);
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-muted);
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

@media (max-width: 768px) {
  .public-users {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .toolbar {
    flex-direction: column;
  }

  .filter-group {
    width: 100%;
  }

  .team-filter,
  .role-filter {
    flex: 1;
  }

  .users-grid {
    grid-template-columns: 1fr;
  }
}
</style>
