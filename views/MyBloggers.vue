<template>
  <div class="my-bloggers-page">
    <div class="page-header">
      <h1>我的博主</h1>
      <p>查看和管理您对接的所有博主</p>
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
    </div>

    <div class="blogger-list">
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="bloggers.length === 0" class="empty">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="8" x2="12" y2="12"/>
          <line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
        <p>还没有对接任何博主</p>
        <router-link to="/add" class="btn btn-primary">立即录入</router-link>
      </div>

      <div v-else class="blogger-cards">
        <div v-for="blogger in bloggers" :key="blogger.id" class="blogger-card" @click="viewDetail(blogger.id)">
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
            <span class="status">{{ blogger.contact ? '已完善' : '待完善' }}</span>
          </div>
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
import { bloggerMyAPI, categoryListAPI } from '../api'

const { error: notifyError } = useNotification()

const router = useRouter()
const loading = ref(false)
const bloggers = ref([])
const total = ref(0)
const allCategories = ref([])

const categoryCount = computed(() => {
  const categories = new Set(bloggers.value.map(b => b.category))
  return categories.size
})

const noContactCount = computed(() => {
  return bloggers.value.filter(b => !b.contact || b.contact.trim() === '').length
})

const loadBloggers = async () => {
  loading.value = true
  try {
    const res = await bloggerMyAPI({ page: 1, pageSize: 100 })
    if (res.code === 200) {
      bloggers.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载博主列表失败')
  } finally {
    loading.value = false
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

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

onMounted(() => {
  loadBloggers()
  loadCategories()
})
</script>

<style scoped>
.my-bloggers-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  text-align: center;
  margin-bottom: 30px;
}

.page-header h1 {
  font-size: 28px;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.page-header p {
  color: var(--text-secondary);
  font-size: 14px;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-icon {
  width: 48px;
  height: 48px;
  background: var(--bg-purple-gradient);
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
  font-size: 28px;
  font-weight: bold;
  color: var(--text-primary);
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: var(--text-muted);
}

.loading {
  text-align: center;
  padding: 60px 20px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border-light);
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading p {
  color: var(--text-muted);
}

.empty {
  text-align: center;
  padding: 80px 20px;
}

.empty svg {
  width: 80px;
  height: 80px;
  color: #ddd;
  margin-bottom: 20px;
}

.empty p {
  color: var(--text-muted);
  margin-bottom: 20px;
  font-size: 16px;
}

.blogger-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.blogger-card {
  background: var(--bg-card);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  cursor: pointer;
  transition: all 0.2s;
}

.blogger-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
}

.card-header {
  display: flex;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #f0f0f0;
  gap: 12px;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: var(--bg-purple-gradient);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-on-brand);
  font-size: 20px;
  font-weight: bold;
  flex-shrink: 0;
  overflow: hidden;
}

.avatar.has-image {
  background: none;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.blogger-info {
  flex: 1;
  min-width: 0;
}

.blogger-info h3 {
  font-size: 16px;
  color: var(--text-primary);
  margin: 0 0 4px 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--text-muted);
}

.platform {
  color: var(--info);
  font-weight: 500;
}

.account {
  color: var(--text-muted);
}

.card-actions {
  flex-shrink: 0;
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
  padding: 20px;
}

.description {
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.6;
  margin: 0 0 12px 0;
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
  display: inline-block;
  padding: 4px 10px;
  background: var(--bg-primary);
  color: var(--text-muted);
  border-radius: 4px;
  font-size: 12px;
}

.product-info,
.contact-info {
  display: flex;
  align-items: center;
  font-size: 13px;
  margin-bottom: 8px;
}

.product-label,
.contact-label {
  color: var(--text-muted);
  min-width: 80px;
}

.product-value {
  color: var(--text-primary);
  font-weight: 500;
}

.contact-value {
  color: var(--text-primary);
}

.contact-value.empty {
  color: #ff9800;
  font-style: italic;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: var(--bg-tertiary);
  border-top: 1px solid #f0f0f0;
}

.date {
  font-size: 12px;
  color: var(--text-muted);
}

.status {
  font-size: 12px;
  font-weight: 500;
  padding: 2px 8px;
  border-radius: 4px;
}

.status {
  background: var(--success-bg);
  color: #4caf50;
}

.btn {
  display: inline-block;
  padding: 10px 24px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  text-decoration: none;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}

.btn-primary {
  background: var(--bg-purple-gradient);
  color: var(--color-on-brand);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

/* 暗色模式样式 */
.dark .my-bloggers-page {
  color: var(--text-primary);
}

.dark .page-header h1 {
  color: var(--text-primary);
}

.dark .page-header p {
  color: var(--text-secondary);
}

.dark .stat-card {
  background: var(--bg-card);
  border-color: var(--border-color);
}

.dark .stat-card .stat-value {
  color: var(--text-primary);
}

.dark .stat-card .stat-label {
  color: var(--text-secondary);
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

.dark .product-value,
.dark .contact-value {
  color: var(--text-primary);
}

.dark .card-footer {
  background: var(--bg-dark);
  border-color: var(--border-color);
}

.dark .date {
  color: var(--text-muted);
}

.dark .status {
  background: rgba(76, 175, 80, 0.2);
  color: var(--success);
}

@media (max-width: 768px) {
  .my-bloggers-page {
    padding: 16px;
  }

  .page-header {
    margin-bottom: 20px;
  }

  .page-header h1 {
    font-size: 22px;
  }

  .stats-cards {
    grid-template-columns: 1fr;
    gap: 12px;
    margin-bottom: 20px;
  }

  .stat-card {
    padding: 16px;
    gap: 12px;
  }

  .stat-icon {
    width: 40px;
    height: 40px;
  }

  .stat-icon svg {
    width: 20px;
    height: 20px;
  }

  .stat-value {
    font-size: 22px;
  }

  .blogger-cards {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .card-header {
    padding: 14px;
    flex-wrap: wrap;
    gap: 8px;
  }

  .avatar {
    width: 40px;
    height: 40px;
    font-size: 16px;
  }

  .category-tag {
    padding: 4px 10px;
    font-size: 11px;
  }

  .card-body {
    padding: 14px;
  }

  .product-label,
  .contact-label {
    min-width: auto;
  }

  .card-footer {
    padding: 10px 14px;
  }
}

@media (max-width: 480px) {
  .my-bloggers-page {
    padding: 12px;
  }

  .page-header h1 {
    font-size: 20px;
  }

  .stat-card {
    padding: 12px;
  }

  .stat-value {
    font-size: 20px;
  }

  .card-header {
    padding: 12px;
  }

  .card-body {
    padding: 12px;
  }

  .card-footer {
    padding: 8px 12px;
  }
}
</style>
