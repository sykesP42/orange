<template>
  <div class="invalid-bloggers">
    <div class="page-header">
      <div class="header-content">
        <h1>失效博主库</h1>
        <p>这些博主因15天未补充联系方式已失效，您可以查看并绑定</p>
      </div>
    </div>

    <div class="announcement-bar" v-if="announcements.length > 0">
      <div class="announcement-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
          <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
        </svg>
      </div>
      <div class="announcement-content">
        <div class="announcement-title">{{ announcements[0].title }}</div>
        <div class="announcement-text">{{ announcements[0].content }}</div>
      </div>
    </div>

    <div class="filters">
      <div class="search-box">
        <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/>
          <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <input
          v-model="filters.keyword"
          type="text"
          placeholder="搜索博主昵称..."
          @input="handleSearch"
        />
      </div>

      <select v-model="filters.category" @change="loadBloggers" class="filter-select">
        <option value="">全部分类</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.name">
          {{ cat.name }}
        </option>
      </select>

      <select v-model="filters.platform" @change="loadBloggers" class="filter-select">
        <option value="">全部平台</option>
        <option v-for="p in platforms" :key="p.id" :value="p.name">
          {{ p.name }}
        </option>
      </select>
    </div>

    <div class="blogger-grid" v-if="bloggers.length > 0">
      <div v-for="blogger in bloggers" :key="blogger.id" class="blogger-card">
        <div class="card-inner">
          <div class="card-bg invalid-bg"></div>
          
          <div v-if="blogger.avatar" class="card-avatar-bg">
            <img :src="blogger.avatar" :alt="blogger.nickname" v-avatar />
          </div>
          
          <div class="card-content">
            <div class="card-top">
              <div class="left-area">
                <div class="status-badge invalid-badge">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="15" y1="9" x2="9" y2="15"/>
                    <line x1="9" y1="9" x2="15" y2="15"/>
                  </svg>
                  <span class="status-text">失效</span>
                </div>
              </div>
            </div>
            
            <div class="nickname">{{ blogger.nickname }}</div>
            <div class="info">
              <span class="tag" :style="{ backgroundColor: getCategoryColor(blogger.category) }">
                {{ blogger.category }}
              </span>
              <span class="platform">{{ blogger.platform }}</span>
            </div>
            <div class="product" v-if="blogger.product">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
                <polyline points="3.27 6.96 12 12.01 20.73 6.96"/>
                <line x1="12" y1="22.08" x2="12" y2="12"/>
              </svg>
              {{ blogger.product }}
            </div>
            <div class="description" v-if="blogger.description">{{ blogger.description }}</div>
            
            <div class="card-bottom">
              <div class="user-info">
                <div class="avatar">
                  {{ blogger.original_owner_name?.charAt(0) || '?' }}
                </div>
                <div class="user-details">
                  <div class="user-name">原归属：{{ blogger.original_owner_name || blogger.original_owner }}</div>
                  <div class="expire-time">失效时间：{{ formatDate(blogger.expire_time) }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="evaluation-section" v-if="blogger.evaluation">
          <div class="evaluation-header">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/>
            </svg>
            <span>失效评价</span>
          </div>
          <div class="evaluation-reason">原因：{{ blogger.evaluation.reason }}</div>
          <div class="evaluation-content" v-if="blogger.evaluation.content">
            {{ blogger.evaluation.content }}
          </div>
          <div class="evaluation-author">
            — {{ blogger.evaluation.real_name }} · {{ formatDate(blogger.evaluation.create_time) }}
          </div>
        </div>

        <div class="action-buttons">
          <button class="bind-btn" @click="bindBlogger(blogger)">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="8.5" cy="7" r="4"/>
              <line x1="20" y1="8" x2="20" y2="14"/>
              <line x1="23" y1="11" x2="17" y2="11"/>
            </svg>
            绑定此博主
          </button>
        </div>
      </div>
    </div>

    <div class="empty-state" v-else-if="!loading">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"/>
        <line x1="15" y1="9" x2="9" y2="15"/>
        <line x1="9" y1="9" x2="15" y2="15"/>
      </svg>
      <h3>暂无失效博主</h3>
      <p>所有博主都在有效期内</p>
    </div>

    <div class="pagination" v-if="total > pageSize">
      <button class="page-btn" :disabled="page === 1" @click="page--; loadBloggers()">
        上一页
      </button>
      <span class="page-info">第 {{ page }} 页 / 共 {{ totalPages }} 页</span>
      <button class="page-btn" :disabled="page === totalPages" @click="page++; loadBloggers()">
        下一页
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { invalidBloggerListAPI, bindInvalidBloggerAPI, getAnnouncementsAPI, categoryListAPI, platformListAPI } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const bloggers = ref([])
const categories = ref([])
const platforms = ref([])
const announcements = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const filters = ref({
  keyword: '',
  category: '',
  platform: ''
})

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const loadBloggers = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value,
      ...filters.value
    }
    const res = await invalidBloggerListAPI(params)
    if (res.code === 200) {
      bloggers.value = res.data.list
      total.value = res.data.total
    }
  } catch (error) {
    console.error('加载失效博主失败:', error)
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  try {
    const res = await categoryListAPI()
    if (res.code === 200) {
      categories.value = res.data
    }
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

const loadPlatforms = async () => {
  try {
    const res = await platformListAPI()
    if (res.code === 200) {
      platforms.value = res.data
    }
  } catch (error) {
    console.error('加载平台失败:', error)
  }
}

const loadAnnouncements = async () => {
  try {
    const res = await getAnnouncementsAPI()
    if (res.code === 200) {
      announcements.value = res.data || []
    }
  } catch (error) {
    console.error('加载公告失败:', error)
  }
}

const getCategoryColor = (categoryName) => {
  const cat = categories.value.find(c => c.name === categoryName)
  return cat?.color || '#6b7280'
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

const handleSearch = () => {
  page.value = 1
  loadBloggers()
}

const bindBlogger = async (blogger) => {
  try {
    await ElMessageBox.confirm(
      `确定要绑定博主【${blogger.nickname}】吗？绑定后该博主将归您所有。`,
      '确认绑定',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const res = await bindInvalidBloggerAPI(blogger.id)
    if (res.code === 200) {
      ElMessage.success('绑定成功！')
      loadBloggers()
    } else {
      ElMessage.error(res.message || '绑定失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('绑定博主失败:', error)
      ElMessage.error('绑定失败')
    }
  }
}

onMounted(() => {
  loadBloggers()
  loadCategories()
  loadPlatforms()
  loadAnnouncements()
})
</script>

<style scoped>
.invalid-bloggers {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-content h1 {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.header-content p {
  color: var(--text-tertiary);
  font-size: 14px;
}

.announcement-bar {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  border-radius: 12px;
  padding: 16px 20px;
  margin-bottom: 24px;
  display: flex;
  align-items: flex-start;
  gap: 12px;
  border: 1px solid var(--border-warning);
}

.announcement-icon {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  color: var(--warning);
  margin-top: 2px;
}

.announcement-content {
  flex: 1;
}

.announcement-title {
  font-weight: 600;
  color: var(--warning);
  font-size: 14px;
  margin-bottom: 4px;
}

.announcement-text {
  color: var(--warning);
  font-size: 13px;
}

.filters {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.search-box {
  flex: 1;
  min-width: 200px;
  position: relative;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  width: 20px;
  height: 20px;
  color: #9ca3af;
}

.search-box input {
  width: 100%;
  padding: 10px 12px 10px 40px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  transition: all 0.2s;
}

.search-box input:focus {
  border-color: var(--info);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.filter-select {
  padding: 10px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  cursor: pointer;
  background: white;
  transition: all 0.2s;
}

.filter-select:focus {
  border-color: var(--info);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.blogger-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.blogger-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.3s;
  cursor: pointer;
}

.blogger-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
}

.card-inner {
  position: relative;
  padding: 20px;
}

.card-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 80px;
  border-radius: 16px 16px 0 0;
  opacity: 0.1;
}

.invalid-bg {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
}

.card-avatar-bg {
  position: absolute;
  top: 40px;
  right: 20px;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  border: 4px solid white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.card-avatar-bg img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.view-detail-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 12px;
  font-weight: 600;
  opacity: 0;
  transition: opacity 0.3s;
  white-space: nowrap;
}

.blogger-card:hover .view-detail-text {
  opacity: 1;
}

.card-content {
  position: relative;
  z-index: 1;
}

.card-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.invalid-badge {
  background: #fef2f2;
  color: #dc2626;
}

.nickname {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
  padding-right: 90px;
}

.info {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.tag {
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  color: white;
}

.platform {
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  background: var(--bg-hover);
  color: var(--text-secondary);
}

.product {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--text-tertiary);
  font-size: 13px;
  margin-bottom: 8px;
}

.product svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.description {
  color: var(--text-tertiary);
  font-size: 13px;
  line-height: 1.5;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-info .avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  font-weight: 600;
}

.user-details {
  font-size: 12px;
}

.user-name {
  color: var(--text-secondary);
  font-weight: 500;
}

.expire-time {
  color: #9ca3af;
  font-size: 11px;
}

.evaluation-section {
  padding: 16px 20px;
  background: var(--bg-card-hover);
  border-top: 1px solid #e5e7eb;
}

.evaluation-header {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-tertiary);
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 8px;
}

.evaluation-header svg {
  width: 16px;
  height: 16px;
}

.evaluation-reason {
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 8px;
}

.evaluation-content {
  color: var(--text-tertiary);
  font-size: 13px;
  line-height: 1.6;
  margin-bottom: 8px;
}

.evaluation-author {
  color: #9ca3af;
  font-size: 12px;
}

.action-buttons {
  padding: 12px 20px 20px;
}

.bind-btn {
  width: 100%;
  padding: 12px;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.bind-btn:hover {
  transform: scale(1.02);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
}

.bind-btn svg {
  width: 18px;
  height: 18px;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
}

.empty-state svg {
  width: 80px;
  height: 80px;
  color: #d1d5db;
  margin-bottom: 20px;
}

.empty-state h3 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.empty-state p {
  color: var(--text-tertiary);
  font-size: 14px;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 32px;
}

.page-btn {
  padding: 8px 16px;
  background: white;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  border-color: var(--info);
  color: var(--info);
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  color: var(--text-tertiary);
  font-size: 14px;
}

@media (max-width: 768px) {
  .invalid-bloggers {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .header-content h1 {
    font-size: 22px;
  }
  
  .filters {
    flex-direction: column;
  }
  
  .search-box {
    min-width: 100%;
  }
  
  .blogger-grid {
    grid-template-columns: 1fr;
  }
  
  .nickname {
    padding-right: 0;
  }
  
  .card-avatar-bg {
    display: none;
  }
}
</style>
