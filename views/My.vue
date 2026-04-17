<template>
  <div class="my-page">
    <div class="page-header">
      <h1>我的</h1>
      <p>管理我的博主和个人信息</p>
    </div>

    <!-- 个人资料卡片（可折叠） -->
    <div class="profile-section">
      <div class="section-header" @click="showProfile = !showProfile">
        <h2>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ 'rotated': showProfile }">
            <polyline points="6,9 12,15 18,9"/>
          </svg>
          个人资料
        </h2>
      </div>

      <div v-show="showProfile" class="profile-content">
        <div v-if="profileLoading" class="section-loading">
          <div class="loading-spinner"></div>
          <span>加载个人资料中...</span>
        </div>
        <template v-else>
        <div class="settings-card">
          <div class="card-header">
            <h3>头像设置</h3>
          </div>
          <div class="card-body">
            <div class="avatar-section">
              <div class="avatar-preview">
                <img v-if="profile.avatar" :src="profile.avatar" alt="头像" v-avatar />
                <div v-else class="avatar-placeholder">
                  {{ profile.real_name?.charAt(0) || profile.username?.charAt(0) || '?' }}
                </div>
              </div>
              <div class="avatar-actions">
                <input
                  type="file"
                  ref="fileInput"
                  accept="image/*"
                  style="display: none"
                  @change="handleFileChange"
                />
                <button class="btn btn-primary" @click="triggerFileInput">
                  上传头像
                </button>
                <p class="hint">支持 JPG、PNG 格式，文件大小不超过 10MB</p>
              </div>
            </div>
          </div>
        </div>

        <AvatarCropper
          v-model:visible="cropDialogVisible"
          :image-file="selectedImageFile"
          :uploadAPI="uploadAvatarAPI"
          @success="handleAvatarUploadSuccess"
        />

        <div class="settings-card">
          <div class="card-header">
            <h3>基本信息</h3>
          </div>
          <div class="card-body">
            <form @submit.prevent="updateProfile">
              <div class="form-row">
                <div class="form-group">
                  <label>用户名</label>
                  <input type="text" v-model="profile.username" disabled />
                  <span class="hint">用户名不可修改</span>
                </div>
                <div class="form-group">
                  <label>真实姓名</label>
                  <input type="text" v-model="form.real_name" placeholder="请输入真实姓名" />
                </div>
              </div>
              <div class="form-row">
                <div class="form-group">
                  <label>邮箱</label>
                  <input type="email" v-model="form.email" placeholder="请输入邮箱" />
                </div>
                <div class="form-group">
                  <label>手机号</label>
                  <input type="tel" v-model="form.phone" placeholder="请输入手机号" />
                </div>
              </div>
              <div class="form-group">
                <label>个人简介</label>
                <textarea v-model="form.bio" rows="3" placeholder="请输入个人简介"></textarea>
              </div>
              <div class="form-group" v-if="profile.approved_by">
                <label>注册批准人</label>
                <input type="text" :value="profile.approved_by" disabled />
                <span class="hint">批准时间：{{ formatDateTime(profile.approved_time) }}</span>
              </div>
              <div class="form-actions">
                <button type="submit" class="btn btn-primary" :disabled="saving">
                  {{ saving ? '保存中...' : '保存修改' }}
                </button>
              </div>
            </form>
          </div>
        </div>

        <div class="settings-card">
          <div class="card-header">
            <h3>我的团队</h3>
          </div>
          <div class="card-body">
            <div v-if="currentTeam" class="team-info-card">
              <div class="team-badge" :style="{ backgroundColor: currentTeam.color }">
                {{ currentTeam.name.charAt(0) }}
              </div>
              <div class="team-details">
                <div class="team-name">{{ currentTeam.name }}</div>
                <div class="team-desc">{{ currentTeam.description || '暂无描述' }}</div>
              </div>
              <button class="btn btn-danger" :disabled="leavingTeam" @click="handleLeaveTeam">{{ leavingTeam ? '退出中...' : '退出团队' }}</button>
            </div>
            <div v-else class="no-team">
              <p>您还没有加入任何团队</p>
              <button class="btn btn-primary" @click="joinTeamDialogVisible = true">加入团队</button>
            </div>
          </div>
        </div>

        <div class="settings-card">
          <div class="card-header">
            <h3>修改密码</h3>
          </div>
          <div class="card-body">
            <form @submit.prevent="updatePassword">
              <div class="form-row">
                <div class="form-group">
                  <label>当前密码</label>
                  <input type="password" v-model="passwordForm.oldPassword" placeholder="请输入当前密码" />
                </div>
                <div class="form-group">
                  <label>新密码</label>
                  <input type="password" v-model="passwordForm.newPassword" placeholder="请输入新密码" />
                </div>
              </div>
              <div class="form-group">
                <label>确认新密码</label>
                <input type="password" v-model="passwordForm.confirmPassword" placeholder="请再次输入新密码" />
              </div>
              <div class="form-actions">
                <button type="submit" class="btn btn-primary" :disabled="updatingPassword">
                  {{ updatingPassword ? '修改中...' : '修改密码' }}
                </button>
              </div>
            </form>
          </div>
        </div>
        </template>
      </div>
    </div>

    <!-- 我的博主 -->
      <div class="blogger-section">
      <div class="section-header">
        <h2>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          我的博主
        </h2>
        <button class="refresh-btn" @click="loadBloggers" :disabled="bloggersLoading">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ 'spinning': bloggersLoading }">
            <path d="M23 4v6h-6"/>
            <path d="M1 20v-6h6"/>
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10"/>
            <path d="M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
          </svg>
        </button>
      </div>

      <div class="stats-cards">
        <div class="stat-card">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
              <circle cx="12" cy="7" r="4"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ total }}</div>
            <div class="stat-label">博主总数</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ categoryCount }}</div>
            <div class="stat-label">分类数</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ noContactCount }}</div>
            <div class="stat-label">待补充联系方式</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
              <line x1="16" y1="2" x2="16" y2="6"/>
              <line x1="8" y1="2" x2="8" y2="6"/>
              <line x1="3" y1="10" x2="21" y2="10"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ thisMonthCount }}</div>
            <div class="stat-label">本月新增</div>
          </div>
        </div>
      </div>

      <!-- 搜索和筛选 -->
      <div class="search-filter-section">
        <div class="search-box">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="搜索博主名称、平台账号、标签..." 
            class="search-input"
          />
        </div>
        
        <div class="filter-box">
          <select v-model="filterPlatform" class="filter-select">
            <option value="">所有平台</option>
            <option v-for="platform in platforms" :key="platform" :value="platform">{{ platform }}</option>
          </select>
          
          <select v-model="filterCategory" class="filter-select">
            <option value="">所有分类</option>
            <option v-for="category in categories" :key="category" :value="category">{{ category }}</option>
          </select>
          
          <button v-if="searchQuery || filterPlatform || filterCategory" @click="clearFilters" class="clear-btn">
            清除筛选
          </button>
        </div>
      </div>

      <div class="blogger-list">
        <div v-if="bloggersLoading" class="loading">
          <div class="spinner"></div>
          <p>加载中...</p>
        </div>

        <div v-else-if="filteredBloggers.length === 0" class="empty">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          <p v-if="searchQuery || filterPlatform || filterCategory">没有找到匹配的博主</p>
          <p v-else>还没有对接任何博主</p>
          <router-link to="/add" class="btn btn-primary">立即录入</router-link>
        </div>

        <div v-else class="blogger-cards">
          <div v-for="blogger in filteredBloggers" :key="blogger.id" class="blogger-card" @click="viewDetail(blogger.id)">
            <div class="card-header">
              <div class="avatar" :class="{ 'has-image': blogger.avatar }">
                <img v-if="blogger.avatar" :src="blogger.avatar" alt="博主头像" v-avatar />
                <template v-else>
                  {{ blogger.nickname?.charAt(0) || '?' }}
                </template>
              </div>
              <div class="blogger-info">
                <h3>{{ blogger.nickname }}</h3>
                <div class="meta">
                  <span class="platform">{{ blogger.platform }}</span>
                  <span class="account">@{{ blogger.platform_account }}</span>
                </div>
              </div>
              <div class="card-actions">
                <span class="category-tag" :style="{ backgroundColor: getCategoryColor(blogger.category) }">
                  <span class="category-icon">{{ getCategoryIcon(blogger.category) }}</span>
                  {{ blogger.category }}
                </span>
              </div>
            </div>
            <div class="card-body">
              <p v-if="blogger.description" class="description">{{ blogger.description }}</p>
              <div v-if="blogger.tags && blogger.tags.length > 0" class="tags">
                <span v-for="tag in blogger.tags" :key="tag" class="tag">{{ tag }}</span>
              </div>
              <div class="product-info">
                <span class="product-label">对接产品：</span>
                <span class="product-value">{{ blogger.product || '未指定' }}</span>
              </div>
              <div class="contact-info">
                <span class="contact-label">联系方式：</span>
                <span class="contact-value" :class="{ 'empty': !blogger.contact }">
                  {{ blogger.contact || '待补充' }}
                </span>
              </div>
            </div>
            <div class="card-footer">
              <span class="date">{{ formatDate(blogger.create_time) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="joinTeamDialogVisible" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="joinTeamDialogVisible = false">
      <div class="modal" style="background:var(--bg-card);">
        <div class="modal-header">
          <h3>加入团队</h3>
          <button class="close-btn" @click="joinTeamDialogVisible = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>选择团队</label>
            <select v-model="selectedTeamId" class="team-select">
              <option value="">请选择团队</option>
              <option v-for="team in teams" :key="team.id" :value="team.id">{{ team.name }}</option>
            </select>
          </div>
          <div v-if="selectedTeamId" class="team-preview">
            <div class="preview-badge" :style="{ backgroundColor: teams.find(t => t.id == selectedTeamId)?.color }">
              {{ teams.find(t => t.id == selectedTeamId)?.name.charAt(0) }}
            </div>
            <div class="preview-info">
              <div class="preview-name">{{ teams.find(t => t.id == selectedTeamId)?.name }}</div>
              <div class="preview-desc">{{ teams.find(t => t.id == selectedTeamId)?.description || '暂无描述' }}</div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="joinTeamDialogVisible = false">取消</button>
          <button class="btn btn-primary" :disabled="joiningTeam" @click="handleJoinTeam">{{ joiningTeam ? '加入中...' : '确认加入' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { logger } from '../utils/logger'
import { useRouter } from 'vue-router'
import { useNotification } from '../stores/notification'
import { useUserStore } from '../stores/user'
import { useConfirm } from '../utils/confirm'
import { getUserProfileAPI, updateUserProfileAPI, updatePasswordAPI, uploadAvatarAPI, bloggerMyAPI, categoryListAPI, getTeamsAPI, updateMyTeamAPI } from '../api'
import AvatarCropper from '../components/AvatarCropper.vue'

const router = useRouter()
const userStore = useUserStore()
const { success, error: notifyError, warning } = useNotification()
const { confirmDanger } = useConfirm()

// 个人资料相关
const showProfile = ref(true)
const fileInput = ref(null)
const cropDialogVisible = ref(false)
const selectedImageFile = ref(null)
const profile = ref({})
const form = ref({
  real_name: '',
  email: '',
  phone: '',
  bio: ''
})
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const saving = ref(false)
const updatingPassword = ref(false)
const profileLoading = ref(false)
const joiningTeam = ref(false)
const leavingTeam = ref(false)
const teams = ref([])
const currentTeam = ref(null)
const joinTeamDialogVisible = ref(false)
const selectedTeamId = ref('')

const loadProfile = async () => {
  profileLoading.value = true
  try {
    const res = await getUserProfileAPI()
    if (res.code === 200) {
      profile.value = res.data
      form.value = {
        real_name: res.data.real_name || '',
        email: res.data.email || '',
        phone: res.data.phone || '',
        bio: res.data.bio || ''
      }
      loadTeams()
    }
  } catch (error) {
    notifyError('加载用户信息失败')
  } finally {
    profileLoading.value = false
  }
}

const loadTeams = async () => {
  try {
    const res = await getTeamsAPI()
    if (res.code === 200) {
      teams.value = res.data || []
      currentTeam.value = teams.value.find(t => t.id === profile.value.team_id) || null
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载团队失败')
  }
}

const handleJoinTeam = async () => {
  if (!selectedTeamId.value) {
    warning('请选择团队')
    return
  }
  joiningTeam.value = true
  const teamId = Number(selectedTeamId.value)
  const team = teams.value.find(t => t.id === teamId)
  try {
    const res = await updateMyTeamAPI(teamId)
    if (res.code === 200) {
      success(`已成功加入【${team?.name || '团队'}】`)
      joinTeamDialogVisible.value = false
      selectedTeamId.value = ''
      if (team) {
        userStore.updateTeam(team.id, team.name, team.color)
      }
      loadProfile()
    } else {
      notifyError(res.message || '加入团队失败')
    }
  } catch (error) {
    notifyError('加入团队失败')
  } finally {
    joiningTeam.value = false
  }
}

const handleLeaveTeam = async () => {
  if (!await confirmDanger('确定退出当前团队？退出后您将不再属于该小组的所有博主数据。')) return
  leavingTeam.value = true
  try {
    const res = await updateMyTeamAPI(null)
    if (res.code === 200) {
      success('已退出团队')
      userStore.updateTeam(null, null, null)
      loadProfile()
    } else {
      notifyError(res.message || '退出团队失败')
    }
  } catch (error) {
    notifyError('退出团队失败')
  } finally {
    leavingTeam.value = false
  }
}

const updateProfile = async () => {
  if (!form.value.real_name?.trim()) {
    warning('请输入真实姓名')
    return
  }
  if (form.value.real_name.trim().length > 20) {
    warning('姓名不能超过20个字符')
    return
  }
  if (form.value.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.value.email)) {
    warning('请输入有效的邮箱地址')
    return
  }
  if (form.value.phone && !/^1[3-9]\d{9}$/.test(form.value.phone)) {
    warning('请输入有效的手机号码')
    return
  }
  if (form.value.bio && form.value.bio.length > 200) {
    warning('个人简介不能超过200个字符')
    return
  }
  form.value.real_name = form.value.real_name.trim()
  if (form.value.email) form.value.email = form.value.email.trim()
  if (form.value.phone) form.value.phone = form.value.phone.trim()
  saving.value = true
  try {
    const res = await updateUserProfileAPI(form.value)
    if (res.code === 200) {
      success('保存成功')
      loadProfile()
    } else {
      notifyError(res.message || '保存失败')
    }
  } catch (error) {
    notifyError('保存失败')
  } finally {
    saving.value = false
  }
}

const updatePassword = async () => {
  if (!passwordForm.value.oldPassword || !passwordForm.value.newPassword) {
    warning('请填写完整密码信息')
    return
  }
  if (passwordForm.value.newPassword.length < 6) {
    warning('新密码长度不能少于6位')
    return
  }
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    warning('两次输入的新密码不一致')
    return
  }
  if (passwordForm.value.oldPassword === passwordForm.value.newPassword) {
    warning('新密码不能与旧密码相同')
    return
  }

  updatingPassword.value = true
  try {
    const res = await updatePasswordAPI({
      oldPassword: passwordForm.value.oldPassword,
      newPassword: passwordForm.value.newPassword
    })
    if (res.code === 200) {
      success('密码修改成功')
      passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
    } else {
      notifyError(res.message || '密码修改失败')
    }
  } catch (error) {
    notifyError('密码修改失败')
  } finally {
    updatingPassword.value = false
  }
}

const handleFileChange = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  if (file.size > 10 * 1024 * 1024) {
    warning('文件大小不能超过 10MB')
    return
  }

  selectedImageFile.value = file
  cropDialogVisible.value = true

  event.target.value = ''
}

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleAvatarUploadSuccess = async (url) => {
  try {
    const res = await updateUserProfileAPI({ avatar: url })
    if (res.code === 200) {
      success('头像上传成功')
      loadProfile()
    } else {
      notifyError(res.message || '更新失败')
    }
  } catch (error) {
    notifyError('更新失败')
  }
}

const formatDateTime = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 我的博主相关
const bloggers = ref([])
const bloggersLoading = ref(false)
const isRefreshing = ref(false)
const allCategories = ref([])

// 搜索和筛选
const searchQuery = ref('')
const filterPlatform = ref('')
const filterCategory = ref('')

const total = computed(() => bloggers.value.length)

const thisMonthCount = computed(() => {
  const now = new Date()
  return bloggers.value.filter(b => {
    const date = new Date(b.create_time)
    return date.getFullYear() === now.getFullYear() && date.getMonth() === now.getMonth()
  }).length
})

const categoryCount = computed(() => {
  const categories = new Set(bloggers.value.map(b => b.category))
  return categories.size
})

const noContactCount = computed(() => {
  return bloggers.value.filter(b => !b.contact).length
})

// 获取所有平台
const platforms = computed(() => {
  const platformSet = new Set(bloggers.value.map(b => b.platform).filter(Boolean))
  return Array.from(platformSet)
})

// 获取所有分类
const categories = computed(() => {
  const categorySet = new Set(bloggers.value.map(b => b.category).filter(Boolean))
  return Array.from(categorySet)
})

// 筛选后的博主列表
const filteredBloggers = computed(() => {
  let result = bloggers.value

  // 搜索
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(b => {
      const matchName = b.nickname?.toLowerCase().includes(query)
      const matchAccount = b.platform_account?.toLowerCase().includes(query)
      const matchTags = b.tags?.some(tag => tag.toLowerCase().includes(query))
      return matchName || matchAccount || matchTags
    })
  }

  // 平台筛选
  if (filterPlatform.value) {
    result = result.filter(b => b.platform === filterPlatform.value)
  }

  // 分类筛选
  if (filterCategory.value) {
    result = result.filter(b => b.category === filterCategory.value)
  }

  return result
})

const clearFilters = () => {
  searchQuery.value = ''
  filterPlatform.value = ''
  filterCategory.value = ''
}

const loadBloggers = async () => {
  bloggersLoading.value = true
  isRefreshing.value = true
  try {
    const res = await bloggerMyAPI()
    if (res.code === 200) {
      // API 返回的是 { list, total, page, pageSize } 对象
      bloggers.value = res.data.list || []
      success('刷新成功')
    }
  } catch (error) {
    notifyError('加载博主列表失败')
  } finally {
    bloggersLoading.value = false
    isRefreshing.value = false
  }
}

const loadCategories = async () => {
  try {
    const res = await categoryListAPI()
    if (res.code === 200) {
      allCategories.value = res.data || []
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载分类失败')
  }
}

const getCategoryColor = (categoryName) => {
  const category = allCategories.value.find(c => c.name === categoryName)
  if (category && category.color) {
    return category.color
  }
  const colorMap = {
    '美妆': 'var(--danger)',
    '美食': 'var(--success)',
    '穿搭': 'var(--info)',
    '旅游': 'var(--success)',
    '科技': 'var(--warning)',
    '其他': '#DFE6E9'
  }
  return colorMap[categoryName] || '#74B9FF'
}

const getCategoryIcon = (categoryName) => {
  const category = allCategories.value.find(c => c.name === categoryName)
  if (category && category.icon) {
    return category.icon
  }
  const iconMap = {
    '美妆': '💄',
    '美食': '🍜',
    '穿搭': '👗',
    '旅游': '✈️',
    '科技': '📱',
    '其他': '📋'
  }
  return iconMap[categoryName] || '📌'
}

const viewDetail = (id) => {
  router.push(`/blogger/${id}`)
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

onMounted(() => {
  loadProfile()
  loadBloggers()
  loadCategories()
})
</script>

<style scoped>
.my-page {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.page-header p {
  color: var(--text-secondary);
  font-size: 14px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  margin-bottom: 16px;
}

.section-header h2 {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.section-header h2 svg {
  width: 24px;
  height: 24px;
  transition: transform 0.3s ease;
}

.section-header h2 svg.rotated {
  transform: rotate(180deg);
}

.refresh-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  transition: all 0.2s;
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.refresh-btn svg {
  width: 20px;
  height: 20px;
}

.refresh-btn svg.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.collapse-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
}

.collapse-btn svg {
  width: 20px;
  height: 20px;
  transition: transform 0.3s ease;
}

.profile-section {
  margin-bottom: 32px;
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.settings-card {
  background: var(--bg-card);
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.card-header {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.card-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.card-body {
  padding: 20px;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.avatar-preview {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  overflow: hidden;
  background: linear-gradient(135deg, #ff6b35 0%, #f7931e 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 36px;
  font-weight: 600;
  color: var(--color-on-brand);
}

.avatar-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.btn {
  padding: 10px 20px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-primary {
  background: linear-gradient(135deg, #ff6b35 0%, #f7931e 100%);
  color: var(--color-on-brand);
}

.btn-primary:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.hint {
  font-size: 12px;
  color: var(--text-muted);
  margin: 0;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.2s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--primary);
}

.form-group input:disabled {
  background: var(--bg-primary);
  cursor: not-allowed;
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}

/* 博主列表样式 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #ff6b35 0%, #f7931e 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-on-brand);
}

.stat-icon svg {
  width: 24px;
  height: 24px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
}

.stat-label {
  font-size: 13px;
  color: var(--text-secondary);
  margin-top: 4px;
}

/* 搜索和筛选 */
.search-filter-section {
  margin-bottom: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
}

.search-box svg {
  position: absolute;
  left: 16px;
  width: 20px;
  height: 20px;
  color: var(--text-muted);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 12px 16px 12px 48px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  font-size: 14px;
  color: var(--text-primary);
  transition: all 0.2s;
}

.search-input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

.filter-box {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.filter-select {
  padding: 10px 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.2s;
}

.filter-select:focus {
  outline: none;
  border-color: var(--primary);
}

.clear-btn {
  padding: 10px 20px;
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.clear-btn:hover {
  background: var(--border-color);
}

.blogger-list {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.loading,
.empty {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-muted);
}

.loading .spinner,
.empty svg {
  width: 48px;
  height: 48px;
  margin-bottom: 16px;
}

.spinner {
  border: 3px solid var(--border-light);
  border-top: 3px solid #ff6b35;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty p {
  margin-bottom: 16px;
}

.empty .btn {
  display: inline-block;
  text-decoration: none;
}

.blogger-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.blogger-card {
  background: var(--bg-card);
  border-radius: 12px;
  border: 1px solid var(--border-color);
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
}

.blogger-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.card-header {
  display: flex;
  gap: 16px;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #ff6b35 0%, #f7931e 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.avatar.has-image {
  background: none;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-on-brand);
}

.blogger-info h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.meta {
  display: flex;
  gap: 12px;
  font-size: 13px;
  color: var(--text-secondary);
}

.card-actions {
  margin-left: auto;
}

.category-tag {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  color: var(--color-on-brand);
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.15);
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.category-tag:hover {
  transform: scale(1.05);
  box-shadow: 0 5px 14px rgba(0, 0, 0, 0.2);
}

.category-icon {
  font-size: 14px;
  line-height: 1;
}

.card-body {
  padding: 16px;
}

.description {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.tag {
  padding: 4px 10px;
  background: var(--bg-hover);
  color: var(--text-secondary);
  border-radius: 6px;
  font-size: 12px;
}

.product-info,
.contact-info {
  display: flex;
  gap: 8px;
  font-size: 13px;
  margin-bottom: 8px;
}

.product-label,
.contact-label {
  color: var(--text-secondary);
  font-weight: 500;
}

.product-value,
.contact-value {
  color: var(--text-primary);
}

.contact-value.empty {
  color: var(--danger);
}

.card-footer {
  padding: 12px 16px;
  background: var(--bg-card-hover);
  font-size: 12px;
  color: var(--text-muted);
}

@media (max-width: 768px) {
  .my-page {
    padding: 16px;
  }

  .page-header h1 {
    font-size: 24px;
  }

  .page-header p {
    font-size: 13px;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .blogger-cards {
    grid-template-columns: 1fr;
  }

  .search-filter-section {
    gap: 12px;
  }

  .filter-box {
    gap: 8px;
  }

  .filter-select {
    flex: 1;
    min-width: 120px;
  }

  .clear-btn {
    width: 100%;
  }

  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  /* 优化移动端输入体验 */
  .edit-input,
  .edit-select,
  .edit-textarea,
  input[type="text"],
  input[type="email"],
  input[type="tel"],
  textarea,
  select {
    font-size: 16px;
    -webkit-tap-highlight-color: transparent;
  }

  /* 优化按钮触摸反馈 */
  .btn,
  button {
    -webkit-tap-highlight-color: transparent;
    user-select: none;
    -webkit-user-select: none;
  }

  /* 优化卡片点击区域 */
  .blogger-card {
    -webkit-tap-highlight-color: transparent;
  }

  /* 优化搜索框 */
  .search-input {
    font-size: 16px;
    padding: 14px 16px 14px 48px;
  }

  /* 优化统计卡片 */
  .stat-card {
    padding: 12px;
  }

  .stat-value {
    font-size: 24px;
  }

  .stat-label {
    font-size: 12px;
  }
}

/* 添加平滑滚动 */
html {
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
}

/* 优化移动端整体体验 */
@media (max-width: 480px) {
  .my-page {
    padding: 12px;
  }

  .profile-section,
  .blogger-section {
    margin-bottom: 20px;
  }

  .section-header h2 {
    font-size: 18px;
  }

  .stats-cards {
    gap: 10px;
  }

  .blogger-card {
    padding: 0;
  }

  .card-header {
    padding: 12px;
  }

  .card-body {
    padding: 12px;
  }

  .blogger-cards {
    grid-template-columns: 1fr;
  }

  .filter-select {
    width: 100%;
  }

  .stats-cards {
    grid-template-columns: 1fr;
  }

  .stat-card {
    padding: 12px;
  }

  .stat-icon {
    width: 36px;
    height: 36px;
  }

  .stat-value {
    font-size: 20px;
  }
}

/* 暗色模式样式 */
.dark .my-page {
  color: var(--text-primary);
}

.dark .page-header h1 {
  color: var(--text-primary);
}

.dark .page-header p {
  color: var(--text-secondary);
}

.dark .settings-card {
  background: var(--bg-card);
  border-color: var(--border-color);
}

.dark .settings-card .card-header h3 {
  color: var(--text-primary);
}

.dark .settings-card label {
  color: var(--text-secondary);
}

.dark .section-header h2 {
  color: var(--text-primary);
}

.dark .stat-card {
  background: var(--bg-card);
  border-color: var(--border-color);
}

.dark .stat-value {
  color: var(--text-primary);
}

.dark .stat-label {
  color: var(--text-secondary);
}

.dark .search-filter-section {
  background: var(--bg-card);
  border-color: var(--border-color);
}

.dark .search-input {
  background: var(--bg-dark);
  color: var(--text-primary);
  border-color: var(--border-color);
}

.dark .filter-select {
  background: var(--bg-dark);
  color: var(--text-primary);
  border-color: var(--border-color);
}

.dark .clear-btn {
  background: var(--bg-dark);
  color: var(--text-secondary);
  border-color: var(--border-color);
}

.dark .blogger-card {
  background: var(--bg-card);
  border-color: var(--border-color);
}

.dark .blogger-card:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.dark .blogger-info h3 {
  color: var(--text-primary);
}

.dark .platform {
  color: var(--text-secondary);
}

.dark .account {
  color: var(--text-muted);
}

.dark .description {
  color: var(--text-secondary);
}

.dark .tag {
  background: var(--bg-dark);
  color: var(--text-secondary);
}

.dark .product-label,
.dark .contact-label {
  color: var(--text-muted);
}

.dark .product-value {
  color: var(--text-primary);
}

.dark .contact-value {
  color: var(--text-primary);
}

.dark .contact-value.empty {
  color: var(--warning);
}

.dark .card-footer {
  background: var(--bg-dark);
  border-color: var(--border-color);
}

.dark .date {
  color: var(--text-muted);
}

.dark .section-header svg {
  color: var(--text-secondary);
}

.dark .hint {
  color: var(--text-muted);
}

.dark .form-group input:disabled {
  background: var(--bg-dark);
  color: var(--text-primary);
  border-color: var(--border-color);
}

.dark .form-group input[type="text"],
.dark .form-group input[type="email"],
.dark .form-group input[type="tel"],
.dark .form-group textarea {
  background: var(--bg-dark);
  color: var(--text-primary);
  border-color: var(--border-color);
}

.dark .form-group input:focus,
.dark .form-group textarea:focus {
  border-color: var(--primary);
}

.team-info-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: rgba(0, 0, 0, 0.02);
  border-radius: 8px;
}

.team-badge {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-on-brand);
  font-size: 20px;
  font-weight: 600;
}

.team-details {
  flex: 1;
}

.team-details .team-name {
  font-weight: 600;
  font-size: 16px;
  margin-bottom: 4px;
}

.team-details .team-desc {
  color: var(--text-secondary);
  font-size: 14px;
}

.no-team {
  text-align: center;
  padding: 24px;
  color: var(--text-secondary);
}

.no-team p {
  margin-bottom: 16px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

html.dark .modal-overlay {
  background: rgba(0, 0, 0, 0.8);
}

.modal {
  background: var(--bg-card);
  border-radius: 16px;
  width: 400px;
  max-width: 90vw;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3), 0 0 0 1px var(--border-color);
}

html.dark .modal {
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.05);
}

.dark .modal {
  background: var(--bg-card);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  color: var(--text-secondary);
}

.modal-body {
  padding: 20px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid var(--border-color);
}

.team-select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 14px;
  background: var(--input-bg, #fff);
  color: var(--text-primary);
}

.team-preview {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 16px;
  padding: 12px;
  background: rgba(0, 0, 0, 0.03);
  border-radius: 8px;
}

.preview-badge {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-on-brand);
  font-weight: 600;
}

.preview-info .preview-name {
  font-weight: 500;
  font-size: 14px;
}

.preview-info .preview-desc {
  font-size: 12px;
  color: var(--text-secondary);
}

.btn-secondary {
  background: var(--bg-dark);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-danger {
  background: var(--danger);
  color: var(--color-on-brand);
}

.dark .team-info-card {
  background: rgba(255, 255, 255, 0.05);
}

.dark .team-preview {
  background: rgba(255, 255, 255, 0.05);
}

.dark .modal {
  background: var(--bg-card);
}

.dark .team-select {
  background: var(--bg-dark);
  color: var(--text-primary);
}

.section-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--text-tertiary);
  font-size: 14px;
  gap: 12px;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--border-color);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
