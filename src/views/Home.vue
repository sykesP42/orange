<template>
  <div class="home">
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

    <div class="invalid-blogger-tip" v-if="invalidBloggerCount > 0" @click="router.push('/invalid-bloggers')">
      <div class="tip-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="8" x2="12" y2="12"/>
          <line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
      </div>
      <div class="tip-content">
        <div class="tip-title">失效博主更新</div>
        <div class="tip-text">当前有 {{ invalidBloggerCount }} 位失效博主可查看绑定</div>
      </div>
      <div class="tip-arrow">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="9 18 15 12 9 6"/>
        </svg>
      </div>
    </div>

    <div class="page-header">
      <div class="header-content">
        <h1>博主列表</h1>
        <p>共 <span class="highlight">{{ total }}</span> 位博主</p>
      </div>
      <div class="header-buttons">
        <button class="invalid-btn" @click="router.push('/invalid-bloggers')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="15" y1="9" x2="9" y2="15"/>
            <line x1="9" y1="9" x2="15" y2="15"/>
          </svg>
          失效博主库
          <span v-if="invalidBloggerCount > 0" class="badge">{{ invalidBloggerCount }}</span>
        </button>
        <button class="add-btn" @click="router.push('/add')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="16"/>
            <line x1="8" y1="12" x2="16" y2="12"/>
          </svg>
          录入博主
        </button>
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

      <select v-model="filters.team_id" @change="loadBloggers" class="filter-select">
        <option value="">全部小组</option>
        <option v-for="team in teams" :key="team.id" :value="team.id">
          {{ team.name }}
        </option>
      </select>

      <select v-model="filters.status" @change="loadBloggers" class="filter-select">
        <option value="">全部状态</option>
        <option v-for="status in statusList" :key="status.id" :value="status.id">
          {{ status.name }}
        </option>
      </select>

      <select v-model="filters.tag" @change="loadBloggers" class="filter-select">
        <option value="">全部标签</option>
        <option v-for="tag in allTags" :key="tag" :value="tag">
          {{ tag }}
        </option>
      </select>
    </div>

    <div class="wordcloud-section" v-if="allTags.length > 0">
      <h3>标签词云</h3>
      <div class="wordcloud">
        <span
          v-for="tag in allTags"
          :key="tag"
          class="cloud-tag"
          :class="{ active: filters.tag === tag }"
          :style="{ fontSize: getTagSize(tag) }"
          @click="toggleTag(tag)"
        >
          {{ tag }}
        </span>
      </div>
    </div>

    <div class="view-toggle">
      <button 
        class="toggle-btn" 
        :class="{ active: viewMode === 'card' }"
        @click="viewMode = 'card'"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="7" height="7" rx="1"/>
          <rect x="14" y="3" width="7" height="7" rx="1"/>
          <rect x="3" y="14" width="7" height="7" rx="1"/>
          <rect x="14" y="14" width="7" height="7" rx="1"/>
        </svg>
        卡片
      </button>
      <button 
        class="toggle-btn" 
        :class="{ active: viewMode === 'list' }"
        @click="viewMode = 'list'"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="3" y1="6" x2="21" y2="6"/>
          <line x1="3" y1="12" x2="21" y2="12"/>
          <line x1="3" y1="18" x2="21" y2="18"/>
        </svg>
        列表
      </button>
    </div>

    <div class="blogger-grid" v-if="bloggers.length > 0 && viewMode === 'card'">
      <div v-for="blogger in bloggers" :key="blogger.id" class="blogger-card" @click="goToDetail(blogger.id)">
        <div class="card-inner">
          <div class="card-bg" :class="getStatusClass(blogger.status)"></div>
          
          <div v-if="blogger.avatar" class="card-avatar-bg">
            <img :src="blogger.avatar" :alt="blogger.nickname" v-avatar />
          </div>
          
          <div class="view-detail-text" :class="getStatusClass(blogger.status)">
            点击查看详情
          </div>
          
          <div class="card-content">
            <div class="card-top">
              <div class="left-area">
                <div class="status-badge" :class="getStatusClass(blogger.status)">
                  <svg v-if="blogger.status === '初次联系'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <path d="M12 6v6l4 2"/>
                  </svg>
                  <svg v-else-if="blogger.status === '洽谈中'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
                  </svg>
                  <svg v-else-if="blogger.status === '已合作'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                    <polyline points="22 4 12 14.01 9 11.01"/>
                  </svg>
                  <svg v-else-if="blogger.status === '已拒绝'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="15" y1="9" x2="9" y2="15"/>
                    <line x1="9" y1="9" x2="15" y2="15"/>
                  </svg>
                  <svg v-else-if="blogger.status === '暂停合作'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="10" y1="9" x2="10" y2="15"/>
                    <line x1="14" y1="9" x2="14" y2="15"/>
                  </svg>
                  <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <path d="M12 6v6l4 2"/>
                  </svg>
                  <span class="status-text">{{ blogger.status || '初次联系' }}</span>
                </div>
              </div>
              
              <div class="right-area">
                <div class="blogger-name">
                  <h3>{{ blogger.nickname }}</h3>
                  <span class="category-small" :style="{ color: getCategoryColor(blogger.category) }">{{ blogger.category }}</span>
                </div>
              </div>
              
              <div class="small-avatar" :class="{ 'has-image': blogger.avatar }">
                <img v-if="blogger.avatar" :src="blogger.avatar" :alt="blogger.nickname" v-avatar />
                <template v-else>{{ blogger.nickname?.charAt(0) || '?' }}</template>
              </div>
            </div>
            
            <div class="divider-line"></div>
            
            <div class="contact-info-bottom">
              <span class="icon">👤</span>
              <span class="text">{{ blogger.real_name }}</span>
            </div>
            
            <div class="card-info">
              <div class="info-row">
                <div class="info-item">
                  <span class="icon">📦</span>
                  <span class="text">{{ blogger.product }}</span>
                </div>
                <div class="info-item">
                  <span class="icon">👤</span>
                  <span class="text">{{ blogger.real_name }}</span>
                </div>
              </div>
              
              <div class="info-row status-info-row">
                <div class="info-item">
                  <span class="icon">📅</span>
                  <span class="text">{{ formatDate(blogger.create_time) }}</span>
                </div>
                <div class="info-item status-pill" :class="getStatusClass(blogger.status)">
                  <span class="text">{{ blogger.status || '初次联系' }}</span>
                </div>
              </div>
              
              <div class="info-row">
                <div class="info-item">
                  <span class="icon">📞</span>
                  <span class="text">{{ blogger.contact || '待添加' }}</span>
                  <button v-if="blogger.contact" class="copy-btn" @click.stop="copyToClipboard(blogger.contact)">复制</button>
                </div>
                <div class="info-item">
                  <span class="icon">🌐</span>
                  <span class="text">{{ blogger.platform || '待添加' }}</span>
                </div>
              </div>
              
              <div class="info-row">
                <div class="info-item">
                  <span class="icon">💬</span>
                  <span class="text">{{ blogger.wechat || '待添加' }}</span>
                  <button v-if="blogger.wechat" class="copy-btn" @click.stop="copyToClipboard(blogger.wechat)">复制</button>
                </div>
              </div>
              
              <div v-if="blogger.tags && blogger.tags.length" class="tags">
                <span v-for="(tag, i) in blogger.tags.slice(0, 3)" :key="i" class="tag">{{ tag }}</span>
                <span v-if="blogger.tags.length > 3" class="tag more">+{{ blogger.tags.length - 3 }}</span>
              </div>
            </div>
            
            <div class="category-bottom">
              <span class="category-icon">{{ getCategoryIcon(blogger.category) }}</span>
              <span class="category-text" :style="{ backgroundColor: getCategoryColor(blogger.category), color: getStatusColor(blogger.status) }">{{ blogger.category }}</span>
            </div>
          </div>
        </div>
        
        <div v-if="userStore.role === 'admin'" class="delete-btn" @click.stop="handleDelete(blogger)">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3,6 5,6 21,6"/>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
          </svg>
        </div>
      </div>
    </div>

    <div class="blogger-list" v-if="bloggers.length > 0 && viewMode === 'list'">
      <div class="list-header">
        <div class="list-cell avatar-cell sortable" @click="toggleSort('avatar')">
          头像
          <span v-if="sortField === 'avatar'" class="sort-icon">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
        </div>
        <div class="list-cell nickname-cell sortable" @click="toggleSort('nickname')">
          昵称
          <span v-if="sortField === 'nickname'" class="sort-icon">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
        </div>
        <div class="list-cell category-cell sortable" @click="toggleSort('category')">
          分类
          <span v-if="sortField === 'category'" class="sort-icon">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
        </div>
        <div class="list-cell status-cell sortable" @click="toggleSort('status')">
          状态
          <span v-if="sortField === 'status'" class="sort-icon">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
        </div>
        <div class="list-cell platform-cell sortable" @click="toggleSort('platform')">
          平台
          <span v-if="sortField === 'platform'" class="sort-icon">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
        </div>
        <div class="list-cell contact-cell sortable" @click="toggleSort('contact')">
          联系方式
          <span v-if="sortField === 'contact'" class="sort-icon">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
        </div>
        <div class="list-cell tags-cell sortable" @click="toggleSort('tags')">
          标签
          <span v-if="sortField === 'tags'" class="sort-icon">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
        </div>
        <div class="list-cell username-cell sortable" @click="toggleSort('team')">
          归属小组
          <span v-if="sortField === 'team'" class="sort-icon">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
        </div>
      </div>
      <div 
        v-for="blogger in bloggers" 
        :key="blogger.id" 
        class="list-row"
        @click="goToDetail(blogger.id)"
      >
        <div class="list-cell avatar-cell">
          <img v-if="blogger.avatar" :src="blogger.avatar" class="avatar-img" :alt="blogger.nickname" v-avatar />
          <div v-else class="avatar-mini" :style="{ background: getCategoryColor(blogger.category) }">
            {{ blogger.nickname?.charAt(0) || '?' }}
          </div>
        </div>
        <div class="list-cell nickname-cell">
          <span class="nickname">{{ blogger.nickname }}</span>
          <span v-if="blogger.fans_count" class="fans">粉丝{{ blogger.fans_count }}万</span>
        </div>
        <div class="list-cell category-cell">
          <span class="category-tag" :style="{ backgroundColor: getCategoryColor(blogger.category), color: getStatusColor(blogger.status) }">
            {{ blogger.category }}
          </span>
        </div>
        <div class="list-cell status-cell">
          <span class="status-dot" :class="getStatusClass(blogger.status)"></span>
          <span class="status-text">{{ blogger.status || '初次联系' }}</span>
        </div>
        <div class="list-cell platform-cell">{{ blogger.platform || '-' }}</div>
        <div class="list-cell contact-cell">
          <span v-if="blogger.contact" class="contact-item">
            <span class="contact-text">{{ blogger.contact }}</span>
            <button class="copy-icon-btn" @click.stop="copyToClipboard(blogger.contact)">📋</button>
          </span>
          <span v-else-if="blogger.wechat" class="contact-item">
            <span class="contact-text">微信: {{ blogger.wechat }}</span>
            <button class="copy-icon-btn" @click.stop="copyToClipboard(blogger.wechat)">📋</button>
          </span>
          <span v-else>-</span>
        </div>
        <div class="list-cell tags-cell">
          <span 
            v-if="blogger.tags && blogger.tags.length" 
            class="tag-mini" 
            v-for="(tag, i) in blogger.tags.slice(0, 2)" 
            :key="i"
            :style="{ backgroundColor: getCategoryColor(blogger.category) + '20', color: getCategoryColor(blogger.category) }"
          >{{ tag }}</span>
          <span v-else>-</span>
        </div>
        <div class="list-cell username-cell">{{ blogger.team_name || '无小组' }}</div>
      </div>
    </div>

    <div v-if="bloggers.length === 0" class="empty-state">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
        <circle cx="12" cy="7" r="4"/>
      </svg>
      <h3>暂无博主数据</h3>
      <p>点击"录入博主"添加第一条记录</p>
    </div>

    <div class="pagination" v-if="total > pageSize">
      <button
        class="page-btn"
        :disabled="page <= 1"
        @click="handlePageChange(page - 1)"
      >
        上一页
      </button>
      <span class="page-info">{{ page }} / {{ totalPages }}</span>
      <button
        class="page-btn"
        :disabled="page >= totalPages"
        @click="handlePageChange(page + 1)"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'
import { useConfirm } from '../utils/confirm'
import { bloggerListAPI, categoryListAPI, userListAPI, bloggerDeleteAPI, getBloggerStatusList, getTeamsAPI, getAnnouncementsAPI, getInvalidBloggerCountAPI } from '../api'

const router = useRouter()
const userStore = useUserStore()
const notification = useNotification()
const { confirmDanger } = useConfirm()

const bloggers = ref([])
const categories = ref([])
const users = ref([])
const teams = ref([])
const statusList = ref([])
const announcements = ref([])
const invalidBloggerCount = ref(0)
const total = ref(0)
const page = ref(1)
const viewMode = ref('card')
const sortField = ref('')
const sortOrder = ref('asc')
const pageSize = 12

const filters = reactive({
  keyword: '',
  category: '',
  team_id: '',
  status: '',
  tag: ''
})

const totalPages = computed(() => Math.ceil(total.value / pageSize))
const allTags = ref([])
const tagCounts = ref({})

const loadStatusList = async () => {
  try {
    const res = await getBloggerStatusList()
    if (res.code === 200) {
      statusList.value = res.data || []
    }
  } catch (error) {
    console.error('加载状态列表失败', error)
  }
}

const loadBloggers = async () => {
  try {
    const params = {
      page: page.value,
      pageSize,
      keyword: filters.keyword,
      category: filters.category,
      team_id: filters.team_id,
      status: filters.status
    }
    const res = await bloggerListAPI(params)
    if (res.code === 200) {
      let list = res.data.list || []
      if (filters.tag) {
        list = list.filter(b => b.tags && b.tags.includes(filters.tag))
      }
      list = list.map(b => ({
        ...b,
        team_name: getTeamNameById(b.team_id)
      }))
      bloggers.value = list
      total.value = res.data.total || 0

      const tags = {}
      res.data.list.forEach(b => {
        if (b.tags) {
          b.tags.forEach(t => {
            tags[t] = (tags[t] || 0) + 1
          })
        }
      })
      allTags.value = Object.keys(tags)
      tagCounts.value = tags
    }
  } catch (error) {
    console.error('加载博主列表失败', error)
  }
}

const loadCategories = async () => {
  try {
    const res = await categoryListAPI()
    if (res.code === 200) {
      categories.value = res.data || []
    }
  } catch (error) {
    console.error('加载分类失败', error)
  }
}

const loadAnnouncements = async () => {
  try {
    const res = await getAnnouncementsAPI()
    if (res.code === 200) {
      announcements.value = res.data || []
    }
  } catch (error) {
    console.error('加载公告失败', error)
  }
}

const loadInvalidBloggerCount = async () => {
  try {
    const res = await getInvalidBloggerCountAPI()
    if (res.code === 200) {
      invalidBloggerCount.value = res.data.count || 0
    }
  } catch (error) {
    console.error('加载失效博主数量失败', error)
  }
}

const loadUsers = async () => {
  try {
    const res = await userListAPI()
    if (res.code === 200) {
      users.value = res.data || []
    }
  } catch (error) {
    console.error('加载用户失败', error)
  }
}

const loadTeams = async () => {
  try {
    const res = await getTeamsAPI()
    if (res.code === 200) {
      teams.value = res.data || []
    }
  } catch (error) {
    console.error('加载团队失败', error)
  }
}

let searchTimer = null
const handleSearch = () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    page.value = 1
    loadBloggers()
  }, 300)
}

const handlePageChange = (newPage) => {
  page.value = newPage
  loadBloggers()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handleDelete = async (blogger) => {
  if (!await confirmDanger(`确定要删除博主"${blogger.nickname}"吗？`)) return

  try {
    const res = await bloggerDeleteAPI(blogger.id)
    if (res.code === 200) {
      notification.success('删除成功')
      loadBloggers()
    } else {
      notification.error(res.message || '删除失败')
    }
  } catch (error) {
    notification.error('删除失败')
  }
}

const goToDetail = (id) => {
  router.push(`/blogger/${id}`)
}

const toggleSort = (field) => {
  if (sortField.value === field) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortOrder.value = 'asc'
  }
  sortBloggers()
}

const sortBloggers = () => {
  if (!sortField.value) return
  
  bloggers.value.sort((a, b) => {
    let valA, valB
    
    switch (sortField.value) {
      case 'nickname':
        valA = a.nickname || ''
        valB = b.nickname || ''
        break
      case 'category':
        valA = a.category || ''
        valB = b.category || ''
        break
      case 'status':
        valA = a.status || ''
        valB = b.status || ''
        break
      case 'platform':
        valA = a.platform || ''
        valB = b.platform || ''
        break
      case 'contact':
        valA = a.contact || a.wechat || ''
        valB = b.contact || b.wechat || ''
        break
      case 'tags':
        valA = (a.tags && a.tags.length) ? a.tags[0] : ''
        valB = (b.tags && b.tags.length) ? b.tags[0] : ''
        break
      case 'team':
        valA = a.team_name || ''
        valB = b.team_name || ''
        break
      default:
        return 0
    }
    
    if (valA < valB) return sortOrder.value === 'asc' ? -1 : 1
    if (valA > valB) return sortOrder.value === 'asc' ? 1 : -1
    return 0
  })
}

const toggleTag = (tag) => {
  if (filters.tag === tag) {
    filters.tag = ''
  } else {
    filters.tag = tag
  }
  loadBloggers()
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    notification.success('复制成功！')
  } catch (error) {
    const textarea = document.createElement('textarea')
    textarea.value = text
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    notification.success('复制成功！')
  }
}

const getTagSize = (tag) => {
  const count = tagCounts.value[tag] || 1
  const minSize = 14
  const maxSize = 32
  const scale = Math.min(count / 5, 1)
  return (minSize + scale * (maxSize - minSize)) + 'px'
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const getStatusClass = (status) => {
  const statusMap = {
    '初次联系': 'status-first',
    '洽谈中': 'status-negotiating',
    '已合作': 'status-cooperated',
    '已拒绝': 'status-rejected',
    '暂停合作': 'status-paused'
  }
  return statusMap[status] || 'status-first'
}

const getStatusColor = (status) => {
  const colorMap = {
    '初次联系': '#3b82f6',
    '洽谈中': '#f97316',
    '已合作': '#22c55e',
    '已拒绝': '#ef4444',
    '暂停合作': '#6b7280'
  }
  return colorMap[status] || '#3b82f6'
}

const getCategoryColor = (categoryName) => {
  const category = categories.value.find(cat => cat.name === categoryName)
  return category ? category.color : '#6b7280'
}

const getCategoryIcon = (categoryName) => {
  const category = categories.value.find(cat => cat.name === categoryName)
  return category ? category.icon : '🏷️'
}

const getRealName = (username) => {
  if (!username) return '-'
  const user = users.value.find(u => u.username === username)
  return user?.real_name || user?.username || username
}

const getTeamNameById = (teamId) => {
  if (!teamId) return ''
  const team = teams.value.find(t => t.id === teamId)
  return team?.name || ''
}

const getTextColorForBackground = (bgColor) => {
  if (!bgColor) return 'white'
  
  let r, g, b
  
  if (bgColor.startsWith('#')) {
    let hex = bgColor.slice(1)
    if (hex.length === 3) {
      hex = hex[0] + hex[0] + hex[1] + hex[1] + hex[2] + hex[2]
    }
    r = parseInt(hex.slice(0, 2), 16)
    g = parseInt(hex.slice(2, 4), 16)
    b = parseInt(hex.slice(4, 6), 16)
  } else if (bgColor.startsWith('rgb')) {
    const match = bgColor.match(/\d+/g)
    if (match && match.length >= 3) {
      r = parseInt(match[0])
      g = parseInt(match[1])
      b = parseInt(match[2])
    }
  }
  
  if (r === undefined || g === undefined || b === undefined) return 'white'
  
  const brightness = (r * 299 + g * 587 + b * 114) / 1000
  return brightness > 128 ? '#1f2937' : 'white'
}

const getStatusGradientClass = (status) => {
  const gradientMap = {
    '初次联系': 'gradient-first',
    '洽谈中': 'gradient-negotiating',
    '已合作': 'gradient-cooperated',
    '已拒绝': 'gradient-rejected',
    '暂停合作': 'gradient-paused'
  }
  return gradientMap[status] || 'gradient-first'
}

onMounted(() => {
  loadBloggers()
  loadCategories()
  loadUsers()
  loadTeams()
  loadStatusList()
  loadAnnouncements()
  loadInvalidBloggerCount()
})
</script>

<style scoped>
.home {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.announcement-bar {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  border-radius: 12px;
  padding: 16px 20px;
  margin-bottom: 20px;
  display: flex;
  align-items: flex-start;
  gap: 12px;
  border: 1px solid #fcd34d;
}

.announcement-icon {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  color: #d97706;
  margin-top: 2px;
}

.announcement-content {
  flex: 1;
}

.announcement-title {
  font-weight: 600;
  color: #92400e;
  font-size: 14px;
  margin-bottom: 4px;
}

.announcement-text {
  color: #78350f;
  font-size: 13px;
}

.invalid-blogger-tip {
  background: linear-gradient(135deg, #eff6ff 0%, #dbeafe 100%);
  border-radius: 12px;
  padding: 16px 20px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  border: 1px solid #bfdbfe;
  cursor: pointer;
  transition: all 0.3s ease;
}

.invalid-blogger-tip:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.15);
}

.tip-icon {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  color: #2563eb;
}

.tip-content {
  flex: 1;
}

.tip-title {
  font-weight: 600;
  color: #1e40af;
  font-size: 14px;
  margin-bottom: 4px;
}

.tip-text {
  color: #3b82f6;
  font-size: 13px;
}

.tip-arrow {
  flex-shrink: 0;
  width: 20px;
  height: 20px;
  color: #2563eb;
}

.header-buttons {
  display: flex;
  gap: 12px;
  align-items: center;
}

.invalid-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  border: none;
  border-radius: 12px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
  position: relative;
}

.invalid-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(99, 102, 241, 0.4);
}

.invalid-btn svg {
  width: 18px;
  height: 18px;
}

.invalid-btn .badge {
  position: absolute;
  top: -8px;
  right: -8px;
  background: #ef4444;
  color: white;
  font-size: 12px;
  font-weight: 600;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid white;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-content h1 {
  font-size: 28px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.header-content p {
  font-size: 14px;
  color: var(--text-muted);
}

.highlight {
  color: var(--primary);
  font-weight: 600;
}

.add-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  border: none;
  border-radius: 12px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(249, 115, 22, 0.4);
}

.add-btn svg {
  width: 18px;
  height: 18px;
}

.filters {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.search-box {
  position: relative;
  flex: 1;
  min-width: 200px;
}

.search-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  width: 18px;
  height: 18px;
  color: var(--text-muted);
}

.search-box input {
  width: 100%;
  height: 44px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 0 16px 0 44px;
  font-size: 14px;
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.search-box input::placeholder {
  color: var(--text-muted);
}

.search-box input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

.filter-select {
  height: 44px;
  padding: 0 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  font-size: 14px;
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.3s ease;
}

.filter-select:focus {
  outline: none;
  border-color: var(--primary);
}

.wordcloud-section {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 24px;
}

.wordcloud-section h3 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}

.wordcloud {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
  padding: 10px;
}

.cloud-tag {
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.3s ease;
  padding: 6px 12px;
  border-radius: 20px;
  font-weight: 500;
}

.cloud-tag:hover {
  color: var(--primary);
  background: rgba(249, 115, 22, 0.1);
}

.cloud-tag.active {
  color: white;
  background: var(--primary);
}

.blogger-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 24px;
}

.view-toggle {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
}

.toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-card);
  color: var(--text-muted);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.toggle-btn svg {
  width: 18px;
  height: 18px;
}

.toggle-btn:hover {
  border-color: var(--info);
  color: var(--info);
}

.toggle-btn.active {
  background: var(--info);
  border-color: var(--info);
  color: white;
}

.dark .toggle-btn {
  background: var(--bg-card);
  border-color: var(--border-color);
  color: var(--text-muted);
}

.dark .toggle-btn:hover {
  border-color: var(--info);
  color: var(--info);
}

.dark .toggle-btn.active {
  background: var(--info);
  border-color: var(--info);
  color: white;
}

.blogger-list {
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.blogger-list .list-header {
  display: flex;
  background: var(--bg-secondary);
  border-bottom: 2px solid var(--border-color);
  font-weight: 600;
  font-size: 13px;
  color: var(--text-secondary);
}

.dark .blogger-list .list-header {
  background: var(--bg-tertiary);
  border-color: var(--border-color);
  color: var(--text-secondary);
}

.blogger-list .list-header .sortable {
  cursor: pointer;
  user-select: none;
  transition: background 0.2s;
}

.blogger-list .list-header .sortable:hover {
  background: rgba(0, 0, 0, 0.04);
}

.dark .blogger-list .list-header .sortable:hover {
  background: rgba(255, 255, 255, 0.08);
}

.sort-icon {
  margin-left: 4px;
  font-size: 12px;
}

.blogger-list .list-row {
  display: flex;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  cursor: pointer;
  transition: background 0.15s ease;
}

.blogger-list .list-row:last-child {
  border-bottom: none;
}

.blogger-list .list-row:hover {
  background: var(--bg-hover);
}

.dark .blogger-list .list-row:hover {
  background: var(--bg-hover);
}

.list-cell {
  padding: 12px 8px;
  display: flex;
  align-items: center;
  font-size: 13px;
}

.avatar-cell { width: 50px; justify-content: center; }
.nickname-cell { width: 150px; flex-direction: column; align-items: flex-start; gap: 2px; }
.category-cell { width: 90px; }
.status-cell { width: 100px; }
.platform-cell { width: 80px; }
.contact-cell { width: 140px; flex: 1; }
.tags-cell { width: 120px; }
.username-cell { width: 80px; }

.avatar-mini {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  color: white;
}

.avatar-img {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  object-fit: cover;
}

.list-row .nickname {
  font-weight: 600;
  color: var(--text-primary);
}

.dark .list-row .nickname {
  color: var(--text-primary);
}

.list-row .fans {
  font-size: 11px;
  color: var(--text-muted);
}

.category-tag {
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  border: 1px solid;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 6px;
}

.status-dot.status-first { background: #3b82f6; }
.status-dot.status-negotiating { background: #f97316; }
.status-dot.status-cooperated { background: #22c55e; }
.status-dot.status-rejected { background: #ef4444; }
.status-dot.status-paused { background: #6b7280; }

.status-text {
  font-size: 12px;
  color: #64748b;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.contact-text {
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.copy-icon-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 12px;
  opacity: 0.6;
  transition: opacity 0.2s;
  padding: 2px;
}

.copy-icon-btn:hover {
  opacity: 1;
}

.tag-mini {
  display: inline-block;
  padding: 2px 8px;
  margin-right: 4px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  border: 1px solid;
}

.blogger-card {
  position: relative;
  background: transparent;
  border: none;
  border-radius: 0;
  overflow: visible;
  cursor: pointer;
  transition: z-index 0.4s ease;
}

.blogger-card:hover {
  z-index: 5;
}

.card-inner {
  position: relative;
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.95) 0%, rgba(255, 255, 255, 0.92) 100%);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 20px;
  overflow: hidden;
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
}

.dark .card-inner {
  background: linear-gradient(145deg, rgba(30, 41, 59, 0.96) 0%, rgba(15, 23, 42, 0.98) 100%);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.card-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 60%;
  height: 100%;
  background: radial-gradient(circle at 80% 20%, rgba(59, 130, 246, 0.18) 0%, transparent 55%);
  opacity: 0.7;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .card-bg {
  background: radial-gradient(circle at 80% 20%, rgba(59, 130, 246, 0.25) 0%, transparent 55%);
}

.card-bg.status-first {
  background: radial-gradient(circle at 80% 20%, rgba(59, 130, 246, 0.22) 0%, transparent 55%);
}

.dark .card-bg.status-first {
  background: radial-gradient(circle at 80% 20%, rgba(59, 130, 246, 0.3) 0%, transparent 55%);
}

.card-bg.status-negotiating {
  background: radial-gradient(circle at 80% 20%, rgba(249, 115, 22, 0.22) 0%, transparent 55%);
}

.dark .card-bg.status-negotiating {
  background: radial-gradient(circle at 80% 20%, rgba(249, 115, 22, 0.3) 0%, transparent 55%);
}

.card-bg.status-cooperated {
  background: radial-gradient(circle at 80% 20%, rgba(34, 197, 94, 0.22) 0%, transparent 55%);
}

.dark .card-bg.status-cooperated {
  background: radial-gradient(circle at 80% 20%, rgba(34, 197, 94, 0.3) 0%, transparent 55%);
}

.card-bg.status-rejected {
  background: radial-gradient(circle at 80% 20%, rgba(239, 68, 68, 0.22) 0%, transparent 55%);
}

.dark .card-bg.status-rejected {
  background: radial-gradient(circle at 80% 20%, rgba(239, 68, 68, 0.3) 0%, transparent 55%);
}

.card-bg.status-paused {
  background: radial-gradient(circle at 80% 20%, rgba(156, 163, 175, 0.22) 0%, transparent 55%);
}

.dark .card-bg.status-paused {
  background: radial-gradient(circle at 80% 20%, rgba(156, 163, 175, 0.3) 0%, transparent 55%);
}

.blogger-card:hover .card-bg {
  opacity: 0.35;
  width: 45%;
}

.card-avatar-bg {
  position: absolute;
  top: 0;
  right: 0;
  width: 45%;
  height: 100%;
  overflow: hidden;
  opacity: 0;
  transform: translateX(100%);
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.card-avatar-bg img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: 0.35;
  filter: blur(3px);
  transform: scale(1.2);
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .card-avatar-bg img {
  opacity: 0.5;
}

.blogger-card:hover .card-avatar-bg {
  opacity: 1;
  transform: translateX(0);
}

.blogger-card:hover .card-avatar-bg img {
  opacity: 0.45;
  filter: blur(0);
  transform: scale(1.15);
}

.dark .blogger-card:hover .card-avatar-bg img {
  opacity: 0.6;
}

.view-detail-text {
  position: absolute;
  top: 18px;
  right: -100px;
  z-index: 20;
  font-size: 12px;
  font-weight: 600;
  color: #3b82f6;
  padding: 4px 10px;
  border-radius: 6px;
  opacity: 0;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.view-detail-text.status-first {
  color: #3b82f6;
}

.view-detail-text.status-negotiating {
  color: #f97316;
}

.view-detail-text.status-cooperated {
  color: #22c55e;
}

.view-detail-text.status-rejected {
  color: #ef4444;
}

.view-detail-text.status-paused {
  color: #6b7280;
}

.dark .view-detail-text.status-paused {
  color: #9ca3af;
}

.blogger-card:hover .view-detail-text {
  right: 18px;
  opacity: 1;
}

.card-info {
  display: flex;
  flex-direction: column;
  gap: 10px;
  opacity: 0;
  max-height: 0;
  overflow: hidden;
  transform: translateY(8px);
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.blogger-card:hover .card-info {
  opacity: 1;
  max-height: 500px;
  transform: translateY(0);
}

.card-content {
  position: relative;
  z-index: 2;
  padding: 22px;
  padding-right: 22px;
  padding-bottom: 60px;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  pointer-events: none;
}

.card-content > * {
  pointer-events: auto;
}

.blogger-card:hover .card-content {
  transform: translateX(-12px);
  padding-right: 55%;
}

.card-top {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 20px;
  position: relative;
}

.left-area {
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex-shrink: 0;
}

.status-badge {
  width: 50px;
  height: 50px;
  border-radius: 16px;
  background: rgba(0, 0, 0, 0.04);
  border: 1px solid rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: visible;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .status-badge {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.status-badge::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.1;
  transform: translate(-50%, -50%);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .status-badge::before {
  opacity: 0.12;
}

.blogger-card:hover .status-badge::before {
  width: 100%;
  height: 100%;
}

.status-badge svg {
  width: 24px;
  height: 24px;
  position: relative;
  z-index: 1;
  flex-shrink: 0;
}

.status-text {
  position: absolute;
  left: 100%;
  margin-left: 12px;
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
  opacity: 0;
  transform: translateX(-8px);
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  color: inherit;
}

.blogger-card:hover .status-text {
  opacity: 1;
  transform: translateX(0);
}

.right-area {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.blogger-card:hover .right-area {
  transform: translateY(40px);
}

.divider-line {
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(0, 0, 0, 0.12), transparent);
  margin: 8px 0 16px 0;
  opacity: 0;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .divider-line {
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.15), transparent);
}

.blogger-card:hover .divider-line {
  opacity: 1;
}

.category-bottom {
  position: absolute;
  bottom: 6px;
  left: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
  opacity: 0;
  transform: translateY(8px) scale(0.92);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 10;
}

.blogger-card:hover .category-bottom {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.category-icon {
  font-size: 18px;
  animation: bounce 2s ease-in-out infinite;
}

@keyframes bounce {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-3px);
  }
}

.category-text {
  font-size: 16px;
  font-weight: 700;
  padding: 5px 12px;
  border-radius: 8px;
  letter-spacing: 0.5px;
  position: relative;
  outline: 1px solid rgba(0, 0, 0, 0.3);
  outline-offset: 1px;
}

.copy-btn {
  margin-left: 8px;
  padding: 2px 8px;
  font-size: 11px;
  font-weight: 600;
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.copy-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.35);
}

.small-avatar {
  position: absolute;
  top: 0;
  right: 0;
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, var(--primary), #7c3aed);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 22px;
  font-weight: 700;
  overflow: hidden;
  border: 2px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
  opacity: 1;
  transform: scale(1);
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .small-avatar {
  border: 2px solid rgba(255, 255, 255, 0.18);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.35);
}

.blogger-card:hover .small-avatar {
  opacity: 0;
  transform: scale(0.4);
}

.small-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.contact-info-bottom {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #4b5563;
  margin-bottom: 16px;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .contact-info-bottom {
  color: rgba(255, 255, 255, 0.8);
}

.blogger-card:hover .contact-info-bottom {
  opacity: 0;
  max-height: 0;
  margin-bottom: 0;
  overflow: hidden;
}

.contact-info-bottom .icon {
  font-size: 18px;
}

.blogger-name {
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.blogger-name h3 {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  line-height: 1.2;
}

.category-small {
  font-size: 14px;
  font-weight: 600;
  opacity: 0.9;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .blogger-name h3 {
  color: white;
}

.blogger-card:hover .blogger-name h3 {
  font-size: 22px;
}

.blogger-card:hover .category-small {
  opacity: 0;
  transform: translateY(-10px);
}

.status-badge.status-first {
  color: #3b82f6;
}

.status-badge.status-negotiating {
  color: #f97316;
}

.status-badge.status-cooperated {
  color: #22c55e;
}

.status-badge.status-rejected {
  color: #ef4444;
}

.status-badge.status-paused {
  color: #6b7280;
}

.dark .status-badge.status-paused {
  color: #9ca3af;
}

.category {
  font-size: 11px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.6px;
}

.dark .category {
  color: rgba(255, 255, 255, 0.6);
}

.card-info {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.info-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.blogger-card:hover .info-row {
  grid-template-columns: 1fr;
}

.blogger-card:hover .status-info-row {
  display: none;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(0, 0, 0, 0.02);
  padding: 9px 12px;
  border-radius: 12px;
  border: 1px solid rgba(0, 0, 0, 0.04);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .info-item {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.blogger-card:hover .info-item {
  backdrop-filter: blur(14px);
  background: rgba(0, 0, 0, 0.03);
}

.dark .blogger-card:hover .info-item {
  background: rgba(255, 255, 255, 0.06);
}

.info-item .icon {
  font-size: 16px;
}

.info-item .text {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  opacity: 0.9;
  transition: all 0.3s ease;
}

.dark .info-item .text {
  color: rgba(255, 255, 255, 0.88);
}

.blogger-card:hover .info-item .text {
  opacity: 1;
  color: #1f2937;
}

.dark .blogger-card:hover .info-item .text {
  color: rgba(255, 255, 255, 0.95);
}

.status-pill {
  background: rgba(59, 130, 246, 0.12);
  border: 1px solid rgba(59, 130, 246, 0.25);
  justify-content: center;
}

.dark .status-pill {
  background: rgba(59, 130, 246, 0.18);
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.status-pill .text {
  font-weight: 700;
  font-size: 12px;
}

.status-pill.status-first {
  background: rgba(59, 130, 246, 0.12);
  border-color: rgba(59, 130, 246, 0.25);
}

.dark .status-pill.status-first {
  background: rgba(59, 130, 246, 0.18);
  border-color: rgba(59, 130, 246, 0.3);
}

.status-pill.status-first .text {
  color: #3b82f6;
}

.status-pill.status-negotiating {
  background: rgba(249, 115, 22, 0.12);
  border-color: rgba(249, 115, 22, 0.25);
}

.dark .status-pill.status-negotiating {
  background: rgba(249, 115, 22, 0.18);
  border-color: rgba(249, 115, 22, 0.3);
}

.status-pill.status-negotiating .text {
  color: #f97316;
}

.status-pill.status-cooperated {
  background: rgba(34, 197, 94, 0.12);
  border-color: rgba(34, 197, 94, 0.25);
}

.dark .status-pill.status-cooperated {
  background: rgba(34, 197, 94, 0.18);
  border-color: rgba(34, 197, 94, 0.3);
}

.status-pill.status-cooperated .text {
  color: #22c55e;
}

.status-pill.status-rejected {
  background: rgba(239, 68, 68, 0.12);
  border-color: rgba(239, 68, 68, 0.25);
}

.dark .status-pill.status-rejected {
  background: rgba(239, 68, 68, 0.18);
  border-color: rgba(239, 68, 68, 0.3);
}

.status-pill.status-rejected .text {
  color: #ef4444;
}

.status-pill.status-paused {
  background: rgba(107, 114, 128, 0.12);
  border-color: rgba(107, 114, 128, 0.25);
}

.dark .status-pill.status-paused {
  background: rgba(156, 163, 175, 0.18);
  border-color: rgba(156, 163, 175, 0.3);
}

.status-pill.status-paused .text {
  color: #6b7280;
}

.dark .status-pill.status-paused .text {
  color: #d1d5db;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 4px;
  margin-bottom: 8px;
  opacity: 0.85;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.blogger-card:hover .tags {
  opacity: 1;
}

.tag {
  padding: 5px 12px;
  background: linear-gradient(135deg, var(--primary) 0%, #7c3aed 100%);
  border-radius: 20px;
  font-size: 11px;
  color: white;
  font-weight: 600;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.15);
  position: relative;
  overflow: hidden;
  opacity: 0.92;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .tag {
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.35);
}

.blogger-card:hover .tag {
  opacity: 1;
}

.tag::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.25), transparent);
  transition: left 0.7s ease;
}

.blogger-card:hover .tag::before {
  left: 100%;
}

.tag.more {
  background: rgba(0, 0, 0, 0.06);
  border: 1px dashed rgba(0, 0, 0, 0.15);
  box-shadow: none;
  color: #6b7280;
}

.dark .tag.more {
  background: rgba(255, 255, 255, 0.12);
  border: 1px dashed rgba(255, 255, 255, 0.22);
  color: rgba(255, 255, 255, 0.7);
}

.delete-btn {
  position: absolute;
  bottom: 12px;
  right: 12px;
  width: 38px;
  height: 38px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.25);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #dc2626;
  cursor: pointer;
  z-index: 10;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  opacity: 0.75;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .delete-btn {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #f87171;
}

.blogger-card:hover .delete-btn {
  opacity: 1;
}

.delete-btn:hover {
  background: rgba(239, 68, 68, 0.18);
  transform: scale(1.1);
  box-shadow: 0 4px 14px rgba(239, 68, 68, 0.2);
  backdrop-filter: blur(14px);
}

.dark .delete-btn:hover {
  background: rgba(239, 68, 68, 0.25);
  box-shadow: 0 4px 14px rgba(239, 68, 68, 0.35);
}

.delete-btn svg {
  width: 18px;
  height: 18px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.empty-state svg {
  width: 80px;
  height: 80px;
  color: var(--text-muted);
  margin-bottom: 24px;
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 18px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.empty-state p {
  font-size: 14px;
  color: var(--text-muted);
}

.empty-state .btn {
  margin-top: 20px;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 32px;
}

.page-btn {
  padding: 10px 20px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  color: var(--text-secondary);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.page-btn:hover:not(:disabled) {
  background: var(--bg-card-hover);
  color: var(--text-primary);
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: var(--text-muted);
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .add-btn {
    width: 100%;
    justify-content: center;
  }

  .filters {
    flex-direction: column;
  }

  .search-box {
    min-width: 100%;
  }

  .filter-select {
    width: 100%;
  }

  .blogger-grid {
    grid-template-columns: 1fr;
  }

  .card {
    padding: 20px;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .card-actions {
    width: 100%;
    justify-content: flex-end;
  }
}

@media (max-width: 480px) {
  .page-header h1 {
    font-size: 20px;
  }

  .card-header h3 {
    font-size: 16px;
  }

  .info-row {
    font-size: 13px;
  }
}
</style>