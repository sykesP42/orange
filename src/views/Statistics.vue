<template>
  <div class="statistics">
    <div class="page-header">
      <h1>数据统计</h1>
      <p>实时数据一览</p>
    </div>

    <div class="view-toggle">
      <button class="toggle-btn" :class="{ active: viewMode === 'blogger' }" @click="viewMode = 'blogger'">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
          <circle cx="9" cy="7" r="4"/>
          <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
          <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
        </svg>
        博主统计
      </button>
      <button class="toggle-btn" :class="{ active: viewMode === 'team' }" @click="switchToTeamView">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
          <circle cx="9" cy="7" r="4"/>
          <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
          <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
        </svg>
        团队统计
      </button>
      <button class="toggle-btn" :class="{ active: viewMode === 'team-blogger' }" @click="switchToTeamBloggerView">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"/>
          <line x1="7" y1="7" x2="7.01" y2="7"/>
        </svg>
        团队对接
      </button>
    </div>

    <div v-if="viewMode === 'blogger'">
    <div class="stats-cards">
      <div class="stat-card" style="--delay: 0">
        <div class="stat-icon total">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.total }}</span>
          <span class="stat-label">博主总数</span>
        </div>
        <div class="stat-trend up">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="23,6 13.5,15.5 8.5,10.5 1,18"/>
            <polyline points="17,6 23,6 23,12"/>
          </svg>
        </div>
      </div>

      <div class="stat-card" style="--delay: 1">
        <div class="stat-icon users">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.userCount }}</span>
          <span class="stat-label">录入人数</span>
        </div>
      </div>

      <div class="stat-card" style="--delay: 2">
        <div class="stat-icon categories">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.categoryCount }}</span>
          <span class="stat-label">分类数量</span>
        </div>
      </div>

      <div class="stat-card" style="--delay: 3">
        <div class="stat-icon products">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.productCount }}</span>
          <span class="stat-label">产品数量</span>
        </div>
      </div>
    </div>

    <div class="charts-grid">
      <div class="chart-card">
        <div class="chart-header">
          <h3>平台分布</h3>
        </div>
        <div class="chart-body">
          <div ref="platformChart" class="chart"></div>
        </div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>分类统计</h3>
        </div>
        <div class="chart-body">
          <div ref="categoryChart" class="chart"></div>
        </div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>月度录入趋势</h3>
        </div>
        <div class="chart-body">
          <div ref="monthlyChart" class="chart"></div>
        </div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>团队成员贡献排行</h3>
        </div>
        <div class="chart-body">
          <div ref="userChart" class="chart"></div>
        </div>
      </div>
    </div>

    <div class="detail-section">
      <div class="section-header">
        <h3>博主详情</h3>
      </div>
      <div class="blogger-details">
        <div
          v-for="(blogger, index) in recentBloggers"
          :key="blogger.id"
          class="blogger-item"
          :style="{ '--item-delay': index * 0.05 }"
        >
          <div class="blogger-avatar">
            {{ blogger.nickname?.charAt(0) || '?' }}
          </div>
          <div class="blogger-info">
            <div class="blogger-main">
              <span class="nickname">{{ blogger.nickname }}</span>
              <span class="category-badge">{{ blogger.category }}</span>
            </div>
            <div class="blogger-sub">
              <span class="product">{{ blogger.product }}</span>
              <span class="divider">•</span>
              <span class="owner">{{ blogger.real_name }}</span>
              <span class="divider">•</span>
              <span class="time">{{ formatDate(blogger.create_time) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    </div>

    <div v-if="viewMode === 'team'">
    <div class="stats-cards">
      <div class="stat-card" style="--delay: 0">
        <div class="stat-icon total">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ teamStats.totalTeams }}</span>
          <span class="stat-label">团队总数</span>
        </div>
      </div>

      <div class="stat-card" style="--delay: 1">
        <div class="stat-icon users">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ teamStats.totalMembers }}</span>
          <span class="stat-label">小组成员总数</span>
        </div>
      </div>

      <div class="stat-card" style="--delay: 2">
        <div class="stat-icon categories">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ teamStats.avgMembersPerTeam }}</span>
          <span class="stat-label">平均每组人数</span>
        </div>
      </div>

      <div class="stat-card" style="--delay: 3">
        <div class="stat-icon products">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ teamStats.largestTeam?.members || 0 }}</span>
          <span class="stat-label">最大团队人数</span>
        </div>
      </div>
    </div>

    <div class="charts-grid">
      <div class="chart-card">
        <div class="chart-header">
          <h3>团队成员分布</h3>
        </div>
        <div class="chart-body">
          <div ref="teamDistChart" class="chart"></div>
        </div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>团队对比</h3>
        </div>
        <div class="chart-body">
          <div ref="teamCompareChart" class="chart"></div>
        </div>
      </div>
    </div>

    <div class="detail-section">
      <div class="section-header">
        <h3>团队详情</h3>
      </div>
      <div class="blogger-details">
        <div
          v-for="(team, index) in teamList"
          :key="team.id"
          class="blogger-item"
          :style="{ '--item-delay': index * 0.05 }"
        >
          <div class="blogger-avatar" :style="{ backgroundColor: team.color }">
            {{ team.name?.charAt(0) || '?' }}
          </div>
          <div class="blogger-info">
            <div class="blogger-main">
              <span class="nickname">{{ team.name }}</span>
              <span class="category-badge">{{ team.leader_name }}</span>
            </div>
            <div class="blogger-sub">
              <span class="product">{{ team.memberCount }} 名成员</span>
              <span class="divider">•</span>
              <span class="owner">{{ team.description || '暂无描述' }}</span>
              <span class="divider">•</span>
              <span class="time">创建于 {{ formatDate(team.create_time) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    </div>

    <div v-if="viewMode === 'team-blogger'">
    <div class="team-blogger-header">
      <div class="team-info-badge" :style="{ backgroundColor: teamBloggerStats.teamInfo?.color || '#6366f1' }">
        <span class="team-name">{{ teamBloggerStats.teamInfo?.name || '我的团队' }}</span>
        <span class="team-leader">组长: {{ teamBloggerStats.teamInfo?.leader_name || '未知' }}</span>
      </div>
    </div>

    <div class="stats-cards">
      <div class="stat-card" style="--delay: 0">
        <div class="stat-icon total">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ teamBloggerStats.total || 0 }}</span>
          <span class="stat-label">对接博主总数</span>
        </div>
        <div class="stat-trend up">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="23,6 13.5,15.5 8.5,10.5 1,18"/>
            <polyline points="17,6 23,6 23,12"/>
          </svg>
        </div>
      </div>

      <div class="stat-card" style="--delay: 1">
        <div class="stat-icon users">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ teamBloggerStats.memberCount || 0 }}</span>
          <span class="stat-label">团队成员数</span>
        </div>
      </div>

      <div class="stat-card" style="--delay: 2">
        <div class="stat-icon categories">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ teamBloggerStats.today || 0 }}</span>
          <span class="stat-label">今日新增</span>
        </div>
      </div>

      <div class="stat-card" style="--delay: 3">
        <div class="stat-icon products">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ teamBloggerStats.month || 0 }}</span>
          <span class="stat-label">本月新增</span>
        </div>
      </div>
    </div>

    <div class="charts-grid">
      <div class="chart-card">
        <div class="chart-header">
          <h3>平台分布</h3>
        </div>
        <div class="chart-body">
          <div ref="teamBloggerPlatformChart" class="chart"></div>
        </div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>分类统计</h3>
        </div>
        <div class="chart-body">
          <div ref="teamBloggerCategoryChart" class="chart"></div>
        </div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>对接状态</h3>
        </div>
        <div class="chart-body">
          <div ref="teamBloggerStatusChart" class="chart"></div>
        </div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>成员贡献排行</h3>
        </div>
        <div class="chart-body">
          <div ref="teamBloggerMemberChart" class="chart"></div>
        </div>
      </div>
    </div>

    <div class="detail-section">
      <div class="section-header">
        <h3>组内博主详情</h3>
      </div>
      <div class="blogger-details">
        <div
          v-for="(blogger, index) in teamBloggerList"
          :key="blogger.id"
          class="blogger-item"
          :style="{ '--item-delay': index * 0.05 }"
        >
          <div class="blogger-avatar">
            {{ blogger.nickname?.charAt(0) || '?' }}
          </div>
          <div class="blogger-info">
            <div class="blogger-main">
              <span class="nickname">{{ blogger.nickname }}</span>
              <span class="category-badge">{{ blogger.category }}</span>
            </div>
            <div class="blogger-sub">
              <span class="product">{{ blogger.product }}</span>
              <span class="divider">•</span>
              <span class="owner">{{ blogger.real_name }}</span>
              <span class="divider">•</span>
              <span class="time">{{ formatDate(blogger.create_time) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { useNotification } from '../stores/notification'
import { bloggerStatAPI, getBloggerCharts, getTeamsAPI, getTeamBloggerStatAPI, getTeamBloggerChartsAPI } from '../api'
import { useUserStore } from '../stores/user'

const { success, error: notifyError, warning, info } = useNotification()

const userStore = useUserStore()

const viewMode = ref('blogger')
const stats = ref({
  total: 0,
  userCount: 0,
  categoryCount: 0,
  productCount: 0
})

const recentBloggers = ref([])

const teamStats = ref({
  totalTeams: 0,
  totalMembers: 0,
  avgMembersPerTeam: 0,
  largestTeam: null
})
const teamList = ref([])

const teamBloggerStats = ref({
  total: 0,
  today: 0,
  month: 0,
  memberCount: 0,
  teamInfo: null
})
const teamBloggerList = ref([])

const platformChart = ref(null)
const categoryChart = ref(null)
const monthlyChart = ref(null)
const userChart = ref(null)
const teamDistChart = ref(null)
const teamCompareChart = ref(null)
const teamBloggerPlatformChart = ref(null)
const teamBloggerCategoryChart = ref(null)
const teamBloggerStatusChart = ref(null)
const teamBloggerMemberChart = ref(null)

let platformChartInstance = null
let categoryChartInstance = null
let monthlyChartInstance = null
let userChartInstance = null
let teamDistChartInstance = null
let teamCompareChartInstance = null
let teamBloggerPlatformChartInstance = null
let teamBloggerCategoryChartInstance = null
let teamBloggerStatusChartInstance = null
let teamBloggerMemberChartInstance = null

const loadStats = async () => {
  try {
    const res = await bloggerStatAPI()
    if (res.code === 200) {
      const data = res.data
      const statsData = data.stats || {}
      stats.value = {
        total: statsData.total || 0,
        userCount: statsData.users || 0,
        categoryCount: data.byCategory?.length || 0,
        productCount: data.byPlatform?.length || 0
      }
      recentBloggers.value = data.recent || []
    }
  } catch (error) {
    console.error('加载统计数据失败', error)
  }
}

const loadCharts = async () => {
  try {
    const res = await getBloggerCharts()
    if (res.code === 200) {
      const data = res.data

      initPlatformChart(data.byPlatform)
      initCategoryChart(data.byCategory)
      initMonthlyChart(data.monthly)
      initUserChart(data.byUser)
    }
  } catch (error) {
    console.error('加载图表数据失败', error)
  }
}

const loadTeamStats = async () => {
  try {
    const [teamsRes, usersRes] = await Promise.all([getTeamsAPI(), fetch('/api/users/list', {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    }).then(r => r.json())])

    if (teamsRes.code === 200) {
      const teams = teamsRes.data || []
      const users = usersRes.data || []
      const activeUsers = users.filter(u => u.status === 'active')

      const teamsWithMembers = teams.map(team => ({
        ...team,
        memberCount: activeUsers.filter(u => u.team_id === team.id).length
      }))

      const totalMembers = teamsWithMembers.reduce((sum, t) => sum + t.memberCount, 0)

      teamStats.value = {
        totalTeams: teams.length,
        totalMembers: totalMembers,
        avgMembersPerTeam: teams.length > 0 ? Math.round(totalMembers / teams.length * 10) / 10 : 0,
        largestTeam: teamsWithMembers.reduce((max, t) => t.memberCount > (max?.memberCount || 0) ? t : max, null)
      }

      teamList.value = teamsWithMembers

      initTeamDistChart(teamsWithMembers)
      initTeamCompareChart(teamsWithMembers)
    }
  } catch (error) {
    console.error('加载团队统计数据失败', error)
  }
}

const switchToTeamView = () => {
  viewMode.value = 'team'
  if (teamStats.value.totalTeams === 0) {
    loadTeamStats()
  }
  setTimeout(() => {
    teamDistChartInstance?.resize()
    teamCompareChartInstance?.resize()
  }, 100)
}

const switchToTeamBloggerView = async () => {
  const storedTeamId = localStorage.getItem('team_id')
  if (!storedTeamId) {
    warning('请先加入一个团队')
    return
  }
  viewMode.value = 'team-blogger'
  if (!teamBloggerStats.value.teamInfo) {
    await loadTeamBloggerStats()
  }
  if (teamBloggerStats.value.total === 0) {
    info('当前团队暂无对接博主数据')
  } else if (teamBloggerStats.value.memberCount < 2) {
    warning('团队人数较少，统计数据可能不够准确')
  }
  setTimeout(() => {
    teamBloggerPlatformChartInstance?.resize()
    teamBloggerCategoryChartInstance?.resize()
    teamBloggerStatusChartInstance?.resize()
    teamBloggerMemberChartInstance?.resize()
  }, 100)
}

const loadTeamBloggerStats = async () => {
  try {
    const teamId = localStorage.getItem('team_id')
    if (!teamId) return

    const [statRes, chartRes] = await Promise.all([
      getTeamBloggerStatAPI(teamId),
      getTeamBloggerChartsAPI(teamId)
    ])

    if (statRes.code === 200) {
      const data = statRes.data
      teamBloggerStats.value = {
        total: data.stats?.total || 0,
        today: data.stats?.today || 0,
        month: data.stats?.month || 0,
        memberCount: data.stats?.memberCount || 0,
        teamInfo: data.teamInfo || null
      }

      const teamMembers = data.byMember || []
      let memberUsernames = []
      const usersRes = await fetch('/api/users/list', {
        headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
      }).then(r => r.json())

      if (usersRes.code === 200) {
        const teamUsers = usersRes.data?.filter(u => u.team_id === teamId) || []
        memberUsernames = teamUsers.map(u => u.username)
      }

      const bloggersRes = await fetch('/api/blogger/list?pageSize=100', {
        headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
      }).then(r => r.json())

      if (bloggersRes.code === 200) {
        teamBloggerList.value = bloggersRes.data?.list?.filter(b => memberUsernames.includes(b.user_name)) || []
      }
    }

    if (chartRes.code === 200) {
      const chartData = chartRes.data
      initTeamBloggerPlatformChart(chartData.byPlatform || [])
      initTeamBloggerCategoryChart(chartData.byCategory || [])
      initTeamBloggerStatusChart(chartData.byStatus || [])
      initTeamBloggerMemberChart(chartData.byMember || [])
    }
  } catch (error) {
    console.error('加载团队对接统计数据失败', error)
  }
}

const initTeamBloggerPlatformChart = (data) => {
  if (!teamBloggerPlatformChart.value) return
  teamBloggerPlatformChartInstance = echarts.init(teamBloggerPlatformChart.value)

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center'
    },
    series: [{
      name: '平台分布',
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: false,
        position: 'center'
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 16,
          fontWeight: 'bold'
        }
      },
      labelLine: {
        show: false
      },
      data: data
    }]
  }

  teamBloggerPlatformChartInstance.setOption(option)
}

const initTeamBloggerCategoryChart = (data) => {
  if (!teamBloggerCategoryChart.value) return
  teamBloggerCategoryChartInstance = echarts.init(teamBloggerCategoryChart.value)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(d => d.name),
      axisLabel: {
        interval: 0,
        rotate: 30
      }
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      name: '分类统计',
      type: 'bar',
      data: data.map(d => d.value),
      itemStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [
            { offset: 0, color: '#8b5cf6' },
            { offset: 1, color: '#a78bfa' }
          ]
        }
      },
      barWidth: '60%'
    }]
  }

  teamBloggerCategoryChartInstance.setOption(option)
}

const initTeamBloggerStatusChart = (data) => {
  if (!teamBloggerStatusChart.value) return
  teamBloggerStatusChartInstance = echarts.init(teamBloggerStatusChart.value)

  const statusColors = {
    '初次联系': '#94a3b8',
    '洽谈中': '#f59e0b',
    '已合作': '#10b981',
    '已拒绝': '#ef4444',
    '暂停合作': '#6b7280'
  }

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center'
    },
    series: [{
      name: '对接状态',
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: false,
        position: 'center'
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 16,
          fontWeight: 'bold'
        }
      },
      labelLine: {
        show: false
      },
      data: data.map(d => ({
        name: d.name,
        value: d.value,
        itemStyle: { color: statusColors[d.name] || '#6366f1' }
      }))
    }]
  }

  teamBloggerStatusChartInstance.setOption(option)
}

const initTeamBloggerMemberChart = (data) => {
  if (!teamBloggerMemberChart.value) return
  teamBloggerMemberChartInstance = echarts.init(teamBloggerMemberChart.value)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(d => d.name),
      axisLabel: {
        interval: 0,
        rotate: 30
      }
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      name: '贡献数量',
      type: 'bar',
      data: data.map(d => d.value),
      itemStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [
            { offset: 0, color: '#f97316' },
            { offset: 1, color: '#fbbf24' }
          ]
        }
      },
      barWidth: '60%'
    }]
  }

  teamBloggerMemberChartInstance.setOption(option)
}

const initTeamDistChart = (data) => {
  if (!teamDistChart.value) return
  teamDistChartInstance = echarts.init(teamDistChart.value)

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center'
    },
    series: [{
      name: '团队成员分布',
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: false,
        position: 'center'
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 20,
          fontWeight: 'bold'
        }
      },
      labelLine: {
        show: false
      },
      data: data.map(t => ({ name: t.name, value: t.memberCount, itemStyle: { color: t.color } }))
    }]
  }

  teamDistChartInstance.setOption(option)
}

const initTeamCompareChart = (data) => {
  if (!teamCompareChart.value) return
  teamCompareChartInstance = echarts.init(teamCompareChart.value)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(t => t.name),
      axisLabel: {
        interval: 0,
        rotate: 30
      }
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      name: '成员数量',
      type: 'bar',
      data: data.map(t => t.memberCount),
      itemStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [
            { offset: 0, color: '#f97316' },
            { offset: 1, color: '#fbbf24' }
          ]
        }
      },
      barWidth: '60%'
    }]
  }

  teamCompareChartInstance.setOption(option)
}

const initPlatformChart = (data) => {
  if (!platformChart.value) return
  
  platformChartInstance = echarts.init(platformChart.value)
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center'
    },
    series: [{
      name: '平台分布',
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: false,
        position: 'center'
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 20,
          fontWeight: 'bold'
        }
      },
      labelLine: {
        show: false
      },
      data: data
    }]
  }
  
  platformChartInstance.setOption(option)
}

const initCategoryChart = (data) => {
  if (!categoryChart.value) return
  
  categoryChartInstance = echarts.init(categoryChart.value)
  
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.name),
      axisLabel: {
        interval: 0,
        rotate: 30
      }
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      name: '博主数量',
      type: 'bar',
      data: data.map(item => item.value),
      itemStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [
            { offset: 0, color: '#f97316' },
            { offset: 1, color: '#fbbf24' }
          ]
        }
      },
      barWidth: '60%'
    }]
  }
  
  categoryChartInstance.setOption(option)
}

const initMonthlyChart = (data) => {
  if (!monthlyChart.value) return
  
  monthlyChartInstance = echarts.init(monthlyChart.value)
  
  const option = {
    tooltip: {
      trigger: 'axis'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.name),
      boundaryGap: false
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      name: '录入数量',
      type: 'line',
      smooth: true,
      data: data.map(item => item.value),
      areaStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(249, 115, 22, 0.3)' },
            { offset: 1, color: 'rgba(249, 115, 22, 0.05)' }
          ]
        }
      },
      lineStyle: {
        color: '#f97316',
        width: 3
      },
      itemStyle: {
        color: '#f97316'
      }
    }]
  }
  
  monthlyChartInstance.setOption(option)
}

const initUserChart = (data) => {
  if (!userChart.value) return
  
  userChartInstance = echarts.init(userChart.value)
  
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'value'
    },
    yAxis: {
      type: 'category',
      data: data.map(item => item.name)
    },
    series: [{
      name: '博主数量',
      type: 'bar',
      data: data.map(item => item.value),
      itemStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 1,
          y2: 0,
          colorStops: [
            { offset: 0, color: '#3b82f6' },
            { offset: 1, color: '#60a5fa' }
          ]
        }
      },
      barWidth: '50%'
    }]
  }
  
  userChartInstance.setOption(option)
}

const handleResize = () => {
  platformChartInstance?.resize()
  categoryChartInstance?.resize()
  monthlyChartInstance?.resize()
  userChartInstance?.resize()
  teamDistChartInstance?.resize()
  teamCompareChartInstance?.resize()
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

onMounted(() => {
  loadStats()
  loadCharts()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  platformChartInstance?.dispose()
  categoryChartInstance?.dispose()
  monthlyChartInstance?.dispose()
  userChartInstance?.dispose()
  teamDistChartInstance?.dispose()
  teamCompareChartInstance?.dispose()
  teamBloggerPlatformChartInstance?.dispose()
  teamBloggerCategoryChartInstance?.dispose()
  teamBloggerStatusChartInstance?.dispose()
  teamBloggerMemberChartInstance?.dispose()
})
</script>

<style scoped>
.statistics {
  animation: fadeIn 0.4s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes countUp {
  from { opacity: 0; transform: scale(0.8); }
  to { opacity: 1; transform: scale(1); }
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
  font-size: 14px;
  color: var(--text-muted);
}

.view-toggle {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
}

.toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border: 1px solid rgba(0, 0, 0, 0.12);
  border-radius: 8px;
  background: white;
  color: #666;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.toggle-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.toggle-btn svg {
  width: 18px;
  height: 18px;
}

.toggle-btn:hover:not(:disabled) {
  border-color: #3b82f6;
  color: #3b82f6;
}

.toggle-btn.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.dark .toggle-btn {
  background: rgba(30, 41, 59, 0.96);
  border-color: rgba(255, 255, 255, 0.12);
  color: #94a3b8;
}

.dark .toggle-btn:hover:not(:disabled) {
  border-color: #3b82f6;
  color: #3b82f6;
}

.dark .toggle-btn.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.team-blogger-header {
  margin-bottom: 24px;
}

.team-info-badge {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  border-radius: 12px;
  color: white;
}

.team-info-badge .team-name {
  font-size: 18px;
  font-weight: 600;
}

.team-info-badge .team-leader {
  font-size: 13px;
  opacity: 0.85;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 20px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  position: relative;
  overflow: hidden;
  animation: slideUp 0.5s ease forwards;
  animation-delay: calc(var(--delay) * 0.1s);
  opacity: 0;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--primary), var(--primary-light));
  opacity: 0;
  transition: opacity 0.3s ease;
}

.stat-card:hover::before {
  opacity: 1;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon svg {
  width: 28px;
  height: 28px;
}

.stat-icon.total {
  background: linear-gradient(135deg, rgba(249, 115, 22, 0.2), rgba(249, 115, 22, 0.1));
  color: var(--primary);
}

.stat-icon.users {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.2), rgba(59, 130, 246, 0.1));
  color: #60a5fa;
}

.stat-icon.categories {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.2), rgba(34, 197, 94, 0.1));
  color: #4ade80;
}

.stat-icon.products {
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.2), rgba(168, 85, 247, 0.1));
  color: #c084fc;
}

.stat-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
  animation: countUp 0.6s ease forwards;
  animation-delay: calc(var(--delay) * 0.1s + 0.2s);
  opacity: 0;
}

.stat-label {
  font-size: 13px;
  color: var(--text-muted);
}

.stat-trend {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-trend svg {
  width: 16px;
  height: 16px;
}

.stat-trend.up {
  background: rgba(34, 197, 94, 0.1);
  color: #4ade80;
}

.charts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 24px;
}

.chart-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 20px;
  padding: 24px;
  animation: slideUp 0.5s ease forwards;
  animation-delay: 0.4s;
  opacity: 0;
}

.chart-header {
  margin-bottom: 20px;
}

.chart-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.chart {
  width: 100%;
  height: 300px;
}

.chart-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.chart-item {
  animation: slideUp 0.4s ease forwards;
  animation-delay: calc(var(--item-delay) + 0.5s);
  opacity: 0;
}

.chart-item-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.chart-item-name {
  font-size: 14px;
  color: var(--text-secondary);
}

.chart-item-value {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.chart-bar {
  height: 8px;
  background: var(--bg-dark);
  border-radius: 4px;
  overflow: hidden;
}

.chart-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--primary), var(--primary-light));
  border-radius: 4px;
  transition: width 1s ease;
}

.product-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.product-tag {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: var(--bg-dark);
  border-radius: 12px;
  animation: slideUp 0.4s ease forwards;
  animation-delay: calc(var(--item-delay) + 0.5s);
  opacity: 0;
  transition: all 0.3s ease;
}

.product-tag:hover {
  transform: translateY(-2px);
  background: var(--bg-card-hover);
}

.tag-name {
  font-size: 14px;
  color: var(--text-secondary);
}

.tag-count {
  padding: 2px 8px;
  background: var(--primary);
  color: white;
  font-size: 12px;
  font-weight: 600;
  border-radius: 6px;
}

.detail-section {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 20px;
  padding: 24px;
  animation: slideUp 0.5s ease forwards;
  animation-delay: 0.6s;
  opacity: 0;
}

.section-header {
  margin-bottom: 20px;
}

.section-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.blogger-details {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.blogger-item {
  position: relative;
  background: transparent;
  border: none;
  border-radius: 0;
  overflow: visible;
  cursor: pointer;
  transition: z-index 0.4s ease;
  padding: 0;
  animation: slideUp 0.4s ease forwards;
  animation-delay: calc(var(--item-delay) * 0.1s);
  opacity: 0;
}

.blogger-item:hover {
  z-index: 5;
}

.blogger-item .card-inner {
  position: relative;
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.95) 0%, rgba(255, 255, 255, 0.92) 100%);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 20px;
  overflow: hidden;
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  padding: 16px;
}

:global(.dark) .blogger-item .card-inner {
  background: linear-gradient(145deg, rgba(30, 41, 59, 0.96) 0%, rgba(15, 23, 42, 0.98) 100%);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

:global(.dark) .blogger-item:hover .card-inner {
  background: linear-gradient(145deg, rgba(30, 41, 59, 0.98) 0%, rgba(15, 23, 42, 0.99) 100%);
}

.blogger-item .card-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 60%;
  height: 100%;
  background: radial-gradient(circle at 80% 20%, rgba(59, 130, 246, 0.18) 0%, transparent 55%);
  opacity: 0.7;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.blogger-item:hover .card-bg {
  opacity: 0.35;
  width: 45%;
}

.blogger-avatar {
  position: relative;
  width: 52px;
  height: 52px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
  font-weight: 600;
  flex-shrink: 0;
  z-index: 2;
}

.blogger-info {
  position: relative;
  flex: 1;
  min-width: 0;
  z-index: 2;
}

.blogger-main {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.nickname {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.category-badge {
  padding: 2px 8px;
  background: rgba(59, 130, 246, 0.15);
  color: #93c5fd;
  font-size: 11px;
  border-radius: 4px;
}

.blogger-sub {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-muted);
}

.divider {
  color: var(--border-color);
}

@media (max-width: 1024px) {
  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .charts-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .stats-cards {
    grid-template-columns: 1fr;
  }

  .stat-card {
    padding: 20px;
  }

  .stat-value {
    font-size: 28px;
  }
}
</style>