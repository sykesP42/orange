<template>
  <div class="statistics">
    <div class="page-header">
      <div class="header-left">
        <h1>数据统计</h1>
        <p>实时数据一览</p>
      </div>
      <div class="header-right">
        <div class="time-range-selector">
          <button
            v-for="range in timeRanges"
            :key="range.value"
            class="range-btn"
            :class="{ active: selectedTimeRange === range.value }"
            @click="selectedTimeRange = range.value"
          >
            {{ range.label }}
          </button>
        </div>
        <button class="layout-btn" @click="resetLayout" title="重置布局">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/>
            <path d="M3 3v5h5"/>
          </svg>
        </button>
      </div>
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
      <button class="toggle-btn" :class="{ active: viewMode === 'graph' }" @click="switchToGraphView">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="6" cy="6" r="3"/>
          <circle cx="18" cy="6" r="3"/>
          <circle cx="12" cy="18" r="3"/>
          <line x1="8.5" y1="7.5" x2="15.5" y2="7.5"/>
          <line x1="7.5" y1="8.5" x2="10.5" y2="15.5"/>
          <line x1="16.5" y1="8.5" x2="13.5" y2="15.5"/>
        </svg>
        关系图谱
      </button>
    </div>

    <div v-if="loading" class="page-loading">
      <div class="loading-spinner"></div>
      <span>加载统计数据中...</span>
    </div>

    <div v-if="!loading && viewMode === 'blogger' && stats.total === 0" class="empty-state">
      <div class="empty-icon">📊</div>
      <h3>暂无统计数据</h3>
      <p>开始录入博主后，这里将展示详细的数据分析</p>
    </div>

    <div v-show="viewMode === 'blogger' && !loading && stats.total > 0" class="dashboard-container">
    <div class="stats-cards" :class="{ 'edit-mode': layoutEditMode }" @dragover.prevent @drop="handleStatsDrop">
      <div
        v-for="(card, index) in statsCardsOrder"
        :key="card.id"
        class="stat-card"
        :class="{ 'dragging': dragFromStats === card.id, 'drag-over': dragOverStats === card.id }"
        :style="{ '--delay': index }"
        :draggable="layoutEditMode"
        @dragstart="handleStatsDragStart($event, card.id, 'stats')"
        @dragend="handleDragEnd"
        @dragenter="handleDragEnter(card.id, 'stats')"
        @dragleave="handleDragLeave('stats')"
      >
        <div v-if="layoutEditMode" class="drag-handle">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="9" cy="5" r="1"/><circle cx="9" cy="12" r="1"/><circle cx="9" cy="19" r="1"/>
            <circle cx="15" cy="5" r="1"/><circle cx="15" cy="12" r="1"/><circle cx="15" cy="19" r="1"/>
          </svg>
        </div>
        <template v-if="card.id === 'total'">
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
        </template>
        <template v-else-if="card.id === 'userCount'">
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
        </template>
        <template v-else-if="card.id === 'categoryCount'">
        <div class="stat-icon categories">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.categoryCount }}</span>
          <span class="stat-label">分类数量</span>
        </div>
        </template>
        <template v-else-if="card.id === 'productCount'">
        <div class="stat-icon products">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.productCount }}</span>
          <span class="stat-label">产品数量</span>
        </div>
        </template>
      </div>
    </div>

    <div class="charts-grid" :class="{ 'edit-mode': layoutEditMode }" @dragover.prevent @drop="handleChartsDrop">
      <div
        v-for="(chart, index) in chartsOrder"
        :key="chart.id"
        class="chart-card"
        :class="{ 'dragging': dragFromCharts === chart.id, 'drag-over': dragOverCharts === chart.id }"
        :draggable="layoutEditMode"
        @dragstart="handleChartsDragStart($event, chart.id, 'charts')"
        @dragend="handleDragEnd"
        @dragenter="handleDragEnter(chart.id, 'charts')"
        @dragleave="handleDragLeave('charts')"
      >
        <div v-if="layoutEditMode" class="drag-handle">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="9" cy="5" r="1"/><circle cx="9" cy="12" r="1"/><circle cx="9" cy="19" r="1"/>
            <circle cx="15" cy="5" r="1"/><circle cx="15" cy="12" r="1"/><circle cx="15" cy="19" r="1"/>
          </svg>
        </div>
        <template v-if="chart.id === 'platform'">
        <div class="chart-header">
          <h3>平台分布</h3>
        </div>
        <div class="chart-body">
          <div id="platform-chart" class="chart"></div>
        </div>
        </template>
        <template v-else-if="chart.id === 'category'">
        <div class="chart-header">
          <h3>分类统计</h3>
        </div>
        <div class="chart-body">
          <div id="category-chart" class="chart"></div>
        </div>
        </template>
        <template v-else-if="chart.id === 'monthly'">
        <div class="chart-header">
          <h3>月度录入趋势</h3>
        </div>
        <div class="chart-body">
          <div id="monthly-chart" class="chart"></div>
        </div>
        </template>
        <template v-else-if="chart.id === 'user'">
        <div class="chart-header">
          <h3>团队成员贡献排行</h3>
        </div>
        <div class="chart-body">
          <div id="user-chart" class="chart"></div>
        </div>
        </template>
      </div>
    </div>

    <div class="detail-section" :class="{ 'edit-mode': layoutEditMode }" @dragover.prevent @drop="handleDetailsDrop">
      <div
        class="section-header-wrapper"
        :class="{ 'dragging': dragFromDetails === 'header', 'drag-over': dragOverDetails === 'header' }"
        :draggable="layoutEditMode"
        @dragstart="handleDetailsDragStart($event, 'header')"
        @dragend="handleDragEnd"
        @dragenter="handleDragEnter('header', 'details')"
        @dragleave="handleDragLeave('details')"
      >
        <div v-if="layoutEditMode" class="drag-handle">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="9" cy="5" r="1"/><circle cx="9" cy="12" r="1"/><circle cx="9" cy="19" r="1"/>
            <circle cx="15" cy="5" r="1"/><circle cx="15" cy="12" r="1"/><circle cx="15" cy="19" r="1"/>
          </svg>
        </div>
      <div class="section-header">
        <h3>博主详情</h3>
      </div>
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

    <div v-show="viewMode === 'team'">
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
          <div id="team-dist-chart" class="chart"></div>
        </div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>团队对比</h3>
        </div>
        <div class="chart-body">
          <div id="team-compare-chart" class="chart"></div>
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

    <div v-show="viewMode === 'team-blogger'">
    <div class="team-blogger-header">
      <div class="team-info-badge" :style="{ backgroundColor: teamBloggerStats.teamInfo?.color || 'var(--purple)' }">
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
          <div id="team-blogger-platform-chart" class="chart"></div>
        </div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>分类统计</h3>
        </div>
        <div class="chart-body">
          <div id="team-blogger-category-chart" class="chart"></div>
        </div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>对接状态</h3>
        </div>
        <div class="chart-body">
          <div id="team-blogger-status-chart" class="chart"></div>
        </div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>成员贡献排行</h3>
        </div>
        <div class="chart-body">
          <div id="team-blogger-member-chart" class="chart"></div>
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

    <div v-show="viewMode === 'graph'" class="graph-container">
      <RelationGraph :bloggers="graphBloggers" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { logger } from '../utils/logger'
import * as echarts from 'echarts'
import RelationGraph from '../components/RelationGraph.vue'
import { useNotification } from '../stores/notification'
import { bloggerStatAPI, getBloggerChartsAPI, getTeamsAPI, getTeamBloggerStatAPI, getTeamBloggerChartsAPI, getUsersListAPI, bloggerListAPI } from '../api'
import { useUserStore } from '../stores/user'
import { useThemeStore } from '../stores/theme'

const { success, error: notifyError, warning, info } = useNotification()

const userStore = useUserStore()
const themeStore = useThemeStore()

const getChartTheme = () => {
  const isDark = themeStore.isDark
  return {
    textColor: isDark ? '#d1d5db' : '#1f2937',
    axisLineColor: isDark ? '#374151' : '#e5e7eb',
    splitLineColor: isDark ? '#1f2937' : '#f3f4f6',
    borderColor: isDark ? '#1f2937' : '#fff',
    tooltipBg: isDark ? '#1f2937' : '#fff',
    tooltipBorder: isDark ? '#374151' : '#e5e7eb',
    tooltipText: isDark ? '#f9fafb' : '#1f2937'
  }
}

const viewMode = ref('blogger')
const loading = ref(false)
const selectedTimeRange = ref('week')
const timeRanges = [
  { label: '本周', value: 'week' },
  { label: '本月', value: 'month' },
  { label: '本季', value: 'quarter' },
  { label: '本年', value: 'year' }
]
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

const layoutEditMode = ref(false)
const dragFromStats = ref(null)
const dragOverStats = ref(null)
const dragFromCharts = ref(null)
const dragOverCharts = ref(null)
const dragFromDetails = ref(null)
const dragOverDetails = ref(null)

const defaultStatsOrder = [
  { id: 'total', label: '博主总数' },
  { id: 'userCount', label: '录入人数' },
  { id: 'categoryCount', label: '分类数量' },
  { id: 'productCount', label: '产品数量' }
]

const defaultChartsOrder = [
  { id: 'platform', label: '平台分布' },
  { id: 'category', label: '分类统计' },
  { id: 'monthly', label: '月度录入趋势' },
  { id: 'user', label: '团队成员贡献排行' }
]

const statsCardsOrder = ref([...defaultStatsOrder])
const chartsOrder = ref([...defaultChartsOrder])

const handleStatsDragStart = (e, id, type) => {
  dragFromStats.value = id
  e.dataTransfer.effectAllowed = 'move'
}

const handleChartsDragStart = (e, id, type) => {
  dragFromCharts.value = id
  e.dataTransfer.effectAllowed = 'move'
}

const handleDetailsDragStart = (e, id) => {
  dragFromDetails.value = id
  e.dataTransfer.effectAllowed = 'move'
}

const handleDragEnter = (id, type) => {
  if (type === 'stats') dragOverStats.value = id
  else if (type === 'charts') dragOverCharts.value = id
  else if (type === 'details') dragOverDetails.value = id
}

const handleDragLeave = (type) => {
  if (type === 'stats') dragOverStats.value = null
  else if (type === 'charts') dragOverCharts.value = null
  else if (type === 'details') dragOverDetails.value = null
}

const handleDragEnd = () => {
  dragFromStats.value = null
  dragOverStats.value = null
  dragFromCharts.value = null
  dragOverCharts.value = null
  dragFromDetails.value = null
  dragOverDetails.value = null
}

const handleStatsDrop = (e) => {
  e.preventDefault()
  if (!dragFromStats.value || !dragOverStats.value) return
  const fromIndex = statsCardsOrder.value.findIndex(c => c.id === dragFromStats.value)
  const toIndex = statsCardsOrder.value.findIndex(c => c.id === dragOverStats.value)
  if (fromIndex !== -1 && toIndex !== -1) {
    const item = statsCardsOrder.value.splice(fromIndex, 1)[0]
    statsCardsOrder.value.splice(toIndex, 0, item)
    saveLayoutPreference()
  }
  handleDragEnd()
}

const handleChartsDrop = (e) => {
  e.preventDefault()
  if (!dragFromCharts.value || !dragOverCharts.value) return
  const fromIndex = chartsOrder.value.findIndex(c => c.id === dragFromCharts.value)
  const toIndex = chartsOrder.value.findIndex(c => c.id === dragOverCharts.value)
  if (fromIndex !== -1 && toIndex !== -1) {
    const item = chartsOrder.value.splice(fromIndex, 1)[0]
    chartsOrder.value.splice(toIndex, 0, item)
    saveLayoutPreference()
  }
  handleDragEnd()
}

const handleDetailsDrop = (e) => {
  e.preventDefault()
  handleDragEnd()
}

const resetLayout = () => {
  statsCardsOrder.value = [...defaultStatsOrder]
  chartsOrder.value = [...defaultChartsOrder]
  localStorage.removeItem('dashboard_layout')
  success('布局已重置')
}

const saveLayoutPreference = () => {
  const layout = {
    statsCards: statsCardsOrder.value.map(c => c.id),
    charts: chartsOrder.value.map(c => c.id)
  }
  localStorage.setItem('dashboard_layout', JSON.stringify(layout))
}

const loadLayoutPreference = () => {
  try {
    const saved = localStorage.getItem('dashboard_layout')
    if (saved) {
      const layout = JSON.parse(saved)
      if (layout.statsCards) {
        statsCardsOrder.value = layout.statsCards
          .map(id => defaultStatsOrder.find(c => c.id === id) || { id, label: id })
          .filter(c => c.id);
      }
      if (layout.charts) {
        chartsOrder.value = layout.charts
          .map(id => defaultChartsOrder.find(c => c.id === id) || { id, label: id })
          .filter(c => c.id);
      }
    }
  } catch (e) {
    logger.error('', e);
    notifyError('加载布局失败')
  }
};

const loadStats = async () => {
  loading.value = true
  try {
    const res = await bloggerStatAPI()
    if (Number(res.code) === 200) {
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
    logger.error('', error)
    notifyError('加载统计数据失败')
  } finally {
    loading.value = false
  }
}

const loadCharts = async () => {
  try {
    const res = await getBloggerChartsAPI()
    if (Number(res.code) === 200) {
      const data = res.data
      await nextTick()
      initPlatformChart(data?.byPlatform)
      initCategoryChart(data?.byCategory)
      initMonthlyChart(data?.monthly)
      initUserChart(data?.byUser)
    } else {
      console.error('加载图表数据失败: code=', res.code, 'message=', res.message)
      notifyError('加载图表数据失败')
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载图表数据失败')
  }
}

const loadTeamStats = async () => {
  try {
    const [teamsRes, usersRes] = await Promise.all([getTeamsAPI(), getUsersListAPI()])

    if (Number(teamsRes.code) === 200) {
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

      await nextTick()
      initTeamDistChart(teamsWithMembers)
      initTeamCompareChart(teamsWithMembers)
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载团队统计数据失败')
  }
}

const switchToTeamView = async () => {
  viewMode.value = 'team'
  await nextTick()
  if (teamStats.value.totalTeams === 0) {
    await loadTeamStats()
  }
  nextTick(() => {
    teamDistChartInstance?.resize()
    teamCompareChartInstance?.resize()
  })
}

const graphBloggers = ref([])

const switchToGraphView = async () => {
  viewMode.value = 'graph'
  await nextTick()
  if (graphBloggers.value.length === 0) {
    try {
      const res = await bloggerListAPI({ page: 1, pageSize: 200 })
      if (Number(res.code) === 200) {
        graphBloggers.value = res.data?.list || res.data?.bloggers || (Array.isArray(res.data) ? res.data : []) || []
      }
    } catch (e) {
      graphBloggers.value = []
    }
  }
}

const switchToTeamBloggerView = async () => {
  const storedTeamId = localStorage.getItem('team_id')
  if (!storedTeamId) {
    warning('请先加入一个团队')
    return
  }
  viewMode.value = 'team-blogger'
  await nextTick()
  if (!teamBloggerStats.value.teamInfo) {
    await loadTeamBloggerStats()
  }
  if (teamBloggerStats.value.total === 0) {
    info('当前团队暂无对接博主数据')
  } else if (teamBloggerStats.value.memberCount < 2) {
    warning('团队人数较少，统计数据可能不够准确')
  }
  nextTick(() => {
    teamBloggerPlatformChartInstance?.resize()
    teamBloggerCategoryChartInstance?.resize()
    teamBloggerStatusChartInstance?.resize()
    teamBloggerMemberChartInstance?.resize()
  })
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
      const usersRes = await getUsersListAPI()

      if (Number(usersRes.code) === 200) {
        const teamUsers = usersRes.data?.filter(u => u.team_id === teamId) || []
        memberUsernames = teamUsers.map(u => u.username)
      }

      const bloggersRes = await bloggerListAPI({ pageSize: 100 })

      if (Number(bloggersRes.code) === 200) {
        teamBloggerList.value = bloggersRes.data?.list?.filter(b => memberUsernames.includes(b.user_name)) || []
      }
    }

    if (Number(chartRes.code) === 200) {
      const chartData = chartRes.data
      await nextTick()
      initTeamBloggerPlatformChart(chartData.byPlatform || [])
      initTeamBloggerCategoryChart(chartData.byCategory || [])
      initTeamBloggerStatusChart(chartData.byStatus || [])
      initTeamBloggerMemberChart(chartData.byMember || [])
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载团队对接统计数据失败')
  }
}

const initTeamBloggerPlatformChart = (data) => {
  const el = document.getElementById('team-blogger-platform-chart')
  if (!el) return
  teamBloggerPlatformChartInstance = echarts.init(el)

  const option = {
    tooltip: {
      trigger: 'item',
          backgroundColor: getChartTheme().tooltipBg,
          borderColor: getChartTheme().tooltipBorder,
          textStyle: { color: getChartTheme().tooltipText },
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      
        textStyle: { color: getChartTheme().textColor },
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
        borderColor: getChartTheme().borderColor,
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
  const el = document.getElementById('team-blogger-category-chart')
  if (!el) return
  teamBloggerCategoryChartInstance = echarts.init(el)

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
      axisLabel: { color: getChartTheme().textColor, interval: 0,
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
            { offset: 1, color: '#7c3aed' }
          ]
        }
      },
      barWidth: '60%'
    }]
  }

  teamBloggerCategoryChartInstance.setOption(option)
}

const initTeamBloggerStatusChart = (data) => {
  const el = document.getElementById('team-blogger-status-chart')
  if (!el) return
  teamBloggerStatusChartInstance = echarts.init(el)

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
          backgroundColor: getChartTheme().tooltipBg,
          borderColor: getChartTheme().tooltipBorder,
          textStyle: { color: getChartTheme().tooltipText },
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      
        textStyle: { color: getChartTheme().textColor },
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
        borderColor: getChartTheme().borderColor,
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
        itemStyle: { color: statusColors[d.name] || '#8b5cf6' }
      }))
    }]
  }

  teamBloggerStatusChartInstance.setOption(option)
}

const initTeamBloggerMemberChart = (data) => {
  const el = document.getElementById('team-blogger-member-chart')
  if (!el) return
  teamBloggerMemberChartInstance = echarts.init(el)

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
      axisLabel: { color: getChartTheme().textColor, interval: 0,
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
            { offset: 1, color: '#f59e0b' }
          ]
        }
      },
      barWidth: '60%'
    }]
  }

  teamBloggerMemberChartInstance.setOption(option)
}

const initTeamDistChart = (data) => {
  const el = document.getElementById('team-dist-chart')
  if (!el) return
  teamDistChartInstance = echarts.init(el)

  const option = {
    tooltip: {
      trigger: 'item',
          backgroundColor: getChartTheme().tooltipBg,
          borderColor: getChartTheme().tooltipBorder,
          textStyle: { color: getChartTheme().tooltipText },
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      
        textStyle: { color: getChartTheme().textColor },
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
        borderColor: getChartTheme().borderColor,
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
  const el = document.getElementById('team-compare-chart')
  if (!el) return
  teamCompareChartInstance = echarts.init(el)

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
      axisLabel: { color: getChartTheme().textColor, interval: 0,
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
            { offset: 1, color: '#f59e0b' }
          ]
        }
      },
      barWidth: '60%'
    }]
  }

  teamCompareChartInstance.setOption(option)
}

const initPlatformChart = (data) => {
  const el = document.getElementById('platform-chart')
  if (!el) return

  if (platformChartInstance) {
    platformChartInstance.dispose()
  }

  try {
    platformChartInstance = echarts.init(el)
  } catch (e) {
    logger.error('', e)
    notifyError('初始化平台图表失败')
    return
  }

  const option = {
    tooltip: {
      trigger: 'item',
          backgroundColor: getChartTheme().tooltipBg,
          borderColor: getChartTheme().tooltipBorder,
          textStyle: { color: getChartTheme().tooltipText },
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      
        textStyle: { color: getChartTheme().textColor },
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
        borderColor: getChartTheme().borderColor,
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
      data: data || []
    }]
  }

  platformChartInstance.setOption(option)
}

const initCategoryChart = (data) => {
  const el = document.getElementById('category-chart')
  if (!el) return

  if (categoryChartInstance) {
    categoryChartInstance.dispose()
  }

  try {
    categoryChartInstance = echarts.init(el)
  } catch (e) {
    logger.error('', e)
    notifyError('初始化分类图表失败')
    return
  }

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
      axisLabel: { color: getChartTheme().textColor, interval: 0,
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
            { offset: 1, color: '#f59e0b' }
          ]
        }
      },
      barWidth: '60%'
    }]
  }
  
  categoryChartInstance.setOption(option)
}

const initMonthlyChart = (data) => {
  const el = document.getElementById('monthly-chart')
  if (!el) return

  if (monthlyChartInstance) {
    monthlyChartInstance.dispose()
  }

  try {
    monthlyChartInstance = echarts.init(el)
  } catch (e) {
    logger.error('', e)
    notifyError('初始化月度图表失败')
    return
  }

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
  const el = document.getElementById('user-chart')
  if (!el) return

  if (userChartInstance) {
    userChartInstance.dispose()
  }

  try {
    userChartInstance = echarts.init(el)
  } catch (e) {
    logger.error('', e)
    notifyError('初始化用户图表失败')
    return
  }

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
            { offset: 1, color: '#2563eb' }
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

onMounted(async () => {
  loadLayoutPreference()
  await loadStats()
  await nextTick()
  loadCharts()
  window.addEventListener('resize', handleResize)
})


watch(() => themeStore.isDark, () => {
  nextTick(() => {
    loadCharts()
  })
})

watch(() => viewMode.value, (mode) => {
  nextTick(() => {
    if (mode === 'blogger') {
      platformChartInstance?.resize()
      categoryChartInstance?.resize()
      monthlyChartInstance?.resize()
      userChartInstance?.resize()
    } else if (mode === 'team') {
      teamDistChartInstance?.resize()
      teamCompareChartInstance?.resize()
    } else if (mode === 'team-blogger') {
      teamBloggerPlatformChartInstance?.resize()
      teamBloggerCategoryChartInstance?.resize()
      teamBloggerStatusChartInstance?.resize()
      teamBloggerMemberChartInstance?.resize()
    }
  })
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.header-left {
  flex: 1;
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

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.time-range-selector {
  display: flex;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 4px;
  gap: 2px;
}

.range-btn {
  padding: 6px 14px;
  background: transparent;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.range-btn:hover {
  color: var(--text-primary);
  background: var(--bg-hover);
}

.range-btn.active {
  background: var(--bg-purple-gradient);
  color: var(--color-on-brand);
  font-weight: 500;
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
  background: var(--bg-card);
  color: var(--text-secondary);
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
  border-color: var(--info);
  color: var(--info);
}

.toggle-btn.active {
  background: var(--info);
  border-color: var(--info);
  color: var(--color-on-brand);
}

.dark .toggle-btn {
  background: rgba(30, 41, 59, 0.96);
  border-color: rgba(255, 255, 255, 0.12);
  color: var(--text-muted);
}

.dark .toggle-btn:hover:not(:disabled) {
  border-color: var(--info);
  color: var(--info);
}

.dark .toggle-btn.active {
  background: var(--info);
  border-color: var(--info);
  color: var(--color-on-brand);
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
  color: var(--color-on-brand);
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
  color: var(--info);
}

.stat-icon.categories {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.2), rgba(34, 197, 94, 0.1));
  color: var(--success);
}

.stat-icon.products {
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.2), rgba(168, 85, 247, 0.1));
  color: var(--purple);
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
  color: var(--success);
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
  color: var(--color-on-brand);
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
  color: var(--color-on-brand);
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
  color: var(--info);
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

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .header-right {
    width: 100%;
    justify-content: space-between;
  }

  .time-range-selector {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    scrollbar-width: none;
  }

  .time-range-selector::-webkit-scrollbar {
    display: none;
  }

  .range-btn {
    padding: 6px 10px;
    font-size: 12px;
    white-space: nowrap;
  }

  .view-toggle {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    scrollbar-width: none;
    gap: 6px;
    padding-bottom: 4px;
  }

  .view-toggle::-webkit-scrollbar {
    display: none;
  }

  .toggle-btn {
    padding: 6px 12px;
    font-size: 13px;
    white-space: nowrap;
    flex-shrink: 0;
  }

  .toggle-btn svg {
    width: 16px;
    height: 16px;
  }

  .stat-card {
    padding: 16px;
    border-radius: 14px;
  }

  .stat-icon {
    width: 44px;
    height: 44px;
    border-radius: 12px;
  }

  .stat-icon svg {
    width: 22px;
    height: 22px;
  }

  .stat-value {
    font-size: 24px;
  }

  .chart-card {
    padding: 16px;
    border-radius: 14px;
  }

  .chart {
    height: 250px;
  }

  .detail-section {
    padding: 16px;
    border-radius: 14px;
  }

  .blogger-details {
    grid-template-columns: 1fr;
  }

  .blogger-avatar {
    width: 40px;
    height: 40px;
    font-size: 16px;
    border-radius: 10px;
  }

  .blogger-sub {
    flex-wrap: wrap;
    gap: 4px;
  }

  .team-info-badge {
    padding: 10px 16px;
    flex-wrap: wrap;
    gap: 8px;
  }

  .team-info-badge .team-name {
    font-size: 16px;
  }
}

@media (max-width: 640px) {
  .stats-cards {
    grid-template-columns: 1fr;
  }

  .stat-card {
    padding: 14px;
  }

  .stat-value {
    font-size: 22px;
  }

  .stat-icon {
    width: 40px;
    height: 40px;
  }
}

@media (max-width: 480px) {
  .page-header h1 {
    font-size: 22px;
  }

  .stat-card {
    padding: 12px;
    gap: 10px;
  }

  .stat-value {
    font-size: 20px;
  }

  .chart-card {
    padding: 12px;
  }

  .chart {
    height: 220px;
  }

  .detail-section {
    padding: 12px;
  }

  .toggle-btn {
    padding: 5px 10px;
    font-size: 12px;
  }

  .range-btn {
    padding: 5px 8px;
    font-size: 11px;
  }
}

.layout-btn {
  width: 36px;
  height: 36px;
  padding: 0;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.layout-btn:hover {
  background: var(--bg-hover);
  border-color: var(--primary);
}

.layout-btn svg {
  width: 18px;
  height: 18px;
  color: var(--text-secondary);
}

.stats-cards.edit-mode,
.charts-grid.edit-mode,
.detail-section.edit-mode {
  cursor: default;
}

.stats-cards .stat-card,
.charts-grid .chart-card {
  transition: transform 0.2s, box-shadow 0.2s;
}

.stats-cards.edit-mode .stat-card,
.charts-grid.edit-mode .chart-card {
  cursor: grab;
  position: relative;
}

.stats-cards.edit-mode .stat-card:active,
.charts-grid.edit-mode .chart-card:active {
  cursor: grabbing;
}

.stat-card.dragging,
.chart-card.dragging {
  opacity: 0.5;
  transform: scale(1.02);
}

.stat-card.drag-over,
.chart-card.drag-over {
  border: 2px dashed var(--primary);
  background: rgba(102, 126, 234, 0.05);
}

.drag-handle {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 4px;
  cursor: grab;
  opacity: 0;
  transition: opacity 0.2s;
}

.stats-cards.edit-mode .stat-card:hover .drag-handle,
.charts-grid.edit-mode .chart-card:hover .drag-handle,
.detail-section.edit-mode .section-header-wrapper:hover .drag-handle {
  opacity: 1;
}

.drag-handle svg {
  width: 16px;
  height: 16px;
  color: var(--primary);
}

.section-header-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  margin: -8px;
  border-radius: 8px;
  transition: background 0.2s;
}

.detail-section.edit-mode .section-header-wrapper:hover {
  background: var(--bg-hover);
}

.detail-section.edit-mode .section-header-wrapper.dragging {
  opacity: 0.5;
}

.detail-section.edit-mode .section-header-wrapper.drag-over {
  border: 2px dashed var(--primary);
  background: rgba(102, 126, 234, 0.05);
}

.graph-container {
  padding: 0;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

.page-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: var(--text-tertiary);
  font-size: 16px;
  gap: 16px;
}

.loading-spinner {
  width: 36px;
  height: 36px;
  border: 3px solid var(--border-color);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: var(--text-tertiary);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.6;
}

.empty-state h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-secondary);
  margin: 0 0 8px 0;
}

.empty-state p {
  font-size: 14px;
  color: var(--text-tertiary);
  margin: 0;
}
</style>