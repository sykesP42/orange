<template>
  <div class="user-detail-page">
    <div class="page-header">
      <button class="back-btn" @click="router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="19" y1="12" x2="5" y2="12"/>
          <polyline points="12,19 5,12 12,5"/>
        </svg>
        返回
      </button>
    </div>

    <div class="loading" v-if="loading">加载中...</div>

    <div class="user-content" v-else-if="user">
      <div class="user-card">
        <div class="avatar" :class="{ 'has-image': user.avatar }">
          <img v-if="user.avatar" :src="user.avatar" alt="用户头像" v-avatar />
          <template v-else>
            {{ user.real_name?.charAt(0) || user.username?.charAt(0) || '?' }}
          </template>
        </div>

        <h1>{{ user.real_name || user.username }}</h1>
        <p class="username">@{{ user.username }}</p>

        <div class="role-tag" :class="user.role">
          {{ user.role === 'admin' ? '管理员' : '普通用户' }}
        </div>

        <div class="info-section">
          <h3>基本信息</h3>
          <div class="info-item">
            <span class="label">用户名</span>
            <span class="value">{{ user.username }}</span>
          </div>
          <div class="info-item" v-if="user.real_name">
            <span class="label">真实姓名</span>
            <span class="value">{{ user.real_name }}</span>
          </div>
          <div class="info-item" v-if="user.email">
            <span class="label">邮箱</span>
            <span class="value">{{ user.email }}</span>
          </div>
          <div class="info-item" v-if="user.phone">
            <span class="label">电话</span>
            <span class="value">{{ user.phone }}</span>
          </div>
          <div class="info-item" v-if="user.bio">
            <span class="label">简介</span>
            <span class="value">{{ user.bio }}</span>
          </div>
          <div class="info-item">
            <span class="label">注册时间</span>
            <span class="value">{{ formatDateTime(user.create_time) }}</span>
          </div>
          <div class="info-item" v-if="user.approved_by">
            <span class="label">批准人</span>
            <span class="value">{{ user.approved_by }}</span>
          </div>
          <div class="info-item" v-if="user.approved_time">
            <span class="label">批准时间</span>
            <span class="value">{{ formatDateTime(user.approved_time) }}</span>
          </div>
        </div>

        <!-- 跟进统计 -->
        <div class="stats-section" v-if="followupStats">
          <h3>跟进统计</h3>
          <div class="stats-grid">
            <div class="stat-card">
              <div class="stat-value">{{ followupStats.total }}</div>
              <div class="stat-label">总跟进数</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ followupStats.thisMonth }}</div>
              <div class="stat-label">本月跟进</div>
            </div>
            <div class="stat-card" v-for="(count, type) in followupStats.byType" :key="type">
              <div class="stat-value">{{ count }}</div>
              <div class="stat-label">{{ type }}</div>
            </div>
          </div>
        </div>

        <!-- 博主列表 -->
        <div class="bloggers-section" v-if="bloggers.length">
          <h3>TA 的博主（{{ bloggers.length }}）</h3>
          <div class="blogger-grid">
            <div v-for="blogger in bloggers" :key="blogger.id" class="blogger-card" @click="viewBlogger(blogger.id)">
              <div class="blogger-avatar" :class="{ 'has-image': blogger.avatar }">
                <img v-if="blogger.avatar" :src="blogger.avatar" alt="博主头像" v-avatar />
                <template v-else>
                  {{ blogger.nickname?.charAt(0) || '?' }}
                </template>
              </div>
              <div class="blogger-card-info">
                <h4>{{ blogger.nickname }}</h4>
                <div class="blogger-meta">
                  <span class="platform-tag">{{ blogger.platform }}</span>
                  <span class="category-tag">{{ blogger.category }}</span>
                </div>
                <div class="blogger-account">@{{ blogger.platform_account }}</div>
              </div>
            </div>
          </div>
        </div>

        <div class="action-section" v-if="user.username !== userStore.username">
          <h3>联系方式</h3>
          <div class="contact-buttons">
            <button v-if="user.phone" class="contact-btn" @click="copyContact(user.phone, '电话号码')">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/>
              </svg>
              <span>复制电话</span>
            </button>
            <button v-if="user.email" class="contact-btn" @click="copyContact(user.email, '邮箱')">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                <polyline points="22,6 12,13 2,6"/>
              </svg>
              <span>复制邮箱</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="empty" v-else>
      <p>用户不存在</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { logger } from '../utils/logger'
import { useRouter, useRoute } from 'vue-router'
import { getUserDetailAPI, getUserBloggersAPI, getUserFollowupStatsAPI } from '../api'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const { success, error: notifyError } = useNotification()

const user = ref(null)
const bloggers = ref([])
const followupStats = ref(null)
const loading = ref(true)

const loadUser = async () => {
  loading.value = true
  try {
    const res = await getUserDetailAPI(route.params.username)
    if (res.code === 200) {
      user.value = res.data
    } else {
      notifyError(res.message || '加载失败')
    }
  } catch (error) {
    notifyError('加载用户信息失败')
  } finally {
    loading.value = false
  }
}

const loadBloggers = async () => {
  try {
    const res = await getUserBloggersAPI(route.params.username)
    if (res.code === 200) {
      bloggers.value = res.data.list || []
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载博主列表失败')
  }
}

const loadFollowupStats = async () => {
  try {
    const res = await getUserFollowupStatsAPI(route.params.username)
    if (res.code === 200) {
      followupStats.value = res.data
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载跟进统计失败')
  }
}

const formatDateTime = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const copyContact = async (value, type) => {
  try {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(value)
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = value
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    }
    success(`已复制${type}到剪贴板`)
  } catch {
    error('复制失败，请手动复制')
  }
}

const viewBlogger = (id) => {
  router.push(`/blogger/${id}`)
}

onMounted(() => {
  loadUser()
  loadBloggers()
  loadFollowupStats()
})
</script>

<style scoped>
.user-detail-page {
  min-height: 100vh;
  background: var(--bg-page);
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.back-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  color: var(--text-secondary);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.back-btn:hover {
  background: var(--bg-card-hover);
  color: var(--text-primary);
}

.loading, .empty {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  color: var(--text-secondary);
  font-size: 16px;
}

.user-content {
  max-width: 800px;
  margin: 0 auto;
}

.user-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 20px;
  padding: 40px;
  text-align: center;
}

.avatar {
  width: 120px;
  height: 120px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  color: var(--bg-card);
  font-weight: 600;
  margin: 0 auto 24px;
  overflow: hidden;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-card h1 {
  font-size: 28px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.username {
  font-size: 16px;
  color: var(--text-secondary);
  margin: 0 0 16px 0;
}

.role-tag {
  display: inline-block;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 32px;
}

.role-tag.admin {
  background: rgba(249, 115, 22, 0.1);
  color: var(--primary);
}

.role-tag.user {
  background: rgba(100, 116, 139, 0.1);
  color: var(--text-secondary);
}

.info-section {
  text-align: left;
  margin-bottom: 32px;
}

.info-section h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 2px solid var(--border-color);
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}

.info-item:last-child {
  border-bottom: none;
}

.info-item .label {
  font-size: 14px;
  color: var(--text-secondary);
  font-weight: 500;
}

.info-item .value {
  font-size: 14px;
  color: var(--text-primary);
}

.bio-text {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.8;
  padding: 12px 0;
}

/* 跟进统计 */
.stats-section {
  margin-bottom: 32px;
}

.stats-section h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 2px solid var(--border-color);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 16px;
}

.stat-card {
  background: var(--bg-card-hover);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 16px;
  text-align: center;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--primary);
  margin-bottom: 8px;
}

.stat-label {
  font-size: 13px;
  color: var(--text-secondary);
}

/* 博主列表 */
.bloggers-section {
  margin-bottom: 32px;
}

.bloggers-section h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 2px solid var(--border-color);
}

.blogger-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.blogger-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: var(--bg-card-hover);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.blogger-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: var(--primary);
}

.blogger-avatar {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: var(--bg-card);
  font-weight: 600;
  flex-shrink: 0;
  overflow: hidden;
}

.blogger-avatar.has-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.blogger-card-info {
  flex: 1;
  min-width: 0;
}

.blogger-card-info h4 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.blogger-meta {
  display: flex;
  gap: 8px;
  margin-bottom: 6px;
}

.platform-tag,
.category-tag {
  font-size: 11px;
  padding: 3px 8px;
  border-radius: 6px;
  background: rgba(249, 115, 22, 0.1);
  color: var(--primary);
  font-weight: 500;
}

.category-tag {
  background: rgba(100, 116, 139, 0.1);
  color: var(--text-secondary);
}

.blogger-account {
  font-size: 12px;
  color: var(--text-secondary);
}

/* 联系方式 */
.action-section {
  padding-top: 24px;
  border-top: 2px solid var(--border-color);
}

.action-section h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 16px 0;
}

.contact-buttons {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.contact-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: var(--bg-card-hover);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.contact-btn:hover {
  background: var(--primary);
  color: var(--bg-card);
  border-color: var(--primary);
  transform: translateY(-2px);
}

.contact-btn svg {
  width: 18px;
  height: 18px;
}

@media (max-width: 768px) {
  .user-detail-page {
    padding: 16px;
  }

  .user-card {
    padding: 24px;
  }

  .avatar {
    width: 100px;
    height: 100px;
    font-size: 40px;
  }

  .user-card h1 {
    font-size: 24px;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .blogger-grid {
    grid-template-columns: 1fr;
  }

  .blogger-card {
    flex-direction: column;
    text-align: center;
  }

  .blogger-meta {
    justify-content: center;
  }

  .contact-buttons {
    flex-direction: column;
  }

  .contact-btn {
    width: 100%;
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .user-detail-page {
    padding: 12px;
  }

  .user-card {
    padding: 16px;
  }

  .avatar {
    width: 72px;
    height: 72px;
    font-size: 28px;
  }

  .user-card h1 {
    font-size: 20px;
  }

  .stats-grid {
    grid-template-columns: 1fr 1fr;
    gap: 8px;
  }

  .stat-value {
    font-size: 22px;
  }

  .blogger-grid {
    grid-template-columns: 1fr;
  }
}
</style>
