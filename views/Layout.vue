﻿<template>
  <div class="layout">
    <header class="header">
      <div class="header-left">
        <button class="menu-toggle" @click="showMenu = !showMenu">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="3" y1="12" x2="21" y2="12"/>
            <line x1="3" y1="6" x2="21" y2="6"/>
            <line x1="3" y1="18" x2="21" y2="18"/>
          </svg>
        </button>
        <div class="logo" @click="router.push('/')">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
              <polyline points="9 22 9 12 15 12 15 22"/>
            </svg>
          </div>
          <span class="logo-text">Orange</span>
        </div>
      </div>

      <nav class="nav">
        <router-link to="/" class="nav-item" :class="{ active: route.path === '/' }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
            <polyline points="9,22 9,12 15,12 15,22"/>
          </svg>
          <span>首页</span>
        </router-link>
        <router-link to="/add" class="nav-item" :class="{ active: route.path === '/add' }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="16"/>
            <line x1="8" y1="12" x2="16" y2="12"/>
          </svg>
          <span>录入</span>
        </router-link>
        <router-link to="/my" class="nav-item" :class="{ active: route.path === '/my' }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          <span>我的</span>
        </router-link>
        <router-link to="/team" class="nav-item" :class="{ active: route.path === '/team' }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
          <span>团队</span>
        </router-link>
        <router-link v-if="userStore.role === 'admin'" to="/admin" class="nav-item" :class="{ active: route.path === '/admin' }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/>
            <path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20z"/>
          </svg>
          <span>管理</span>
        </router-link>
        <router-link to="/statistics" class="nav-item" :class="{ active: route.path === '/statistics' }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="20" x2="18" y2="10"/>
            <line x1="12" y1="20" x2="12" y2="4"/>
            <line x1="6" y1="20" x2="6" y2="14"/>
          </svg>
          <span>统计</span>
        </router-link>
        <router-link to="/calendar" class="nav-item" :class="{ active: route.path === '/calendar' }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
            <line x1="16" y1="2" x2="16" y2="6"/>
            <line x1="8" y1="2" x2="8" y2="6"/>
            <line x1="3" y1="10" x2="21" y2="10"/>
          </svg>
          <span>日历</span>
        </router-link>
        <router-link to="/workflow" class="nav-item" :class="{ active: route.path === '/workflow' }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
          </svg>
          <span>工作流</span>
        </router-link>
      </nav>

      <div class="header-right">
        <button class="command-palette-trigger" @click="openCommandPalette" title="命令面板 (Ctrl+K)">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
          <span class="trigger-text">搜索</span>
          <kbd>⌘K</kbd>
        </button>
        <button class="theme-toggle" @click="themeStore.toggleTheme" :title="themeStore.isDark ? '切换到浅色模式' : '切换到深色模式'">
          <svg v-if="!themeStore.isDark" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="5"/>
            <line x1="12" y1="1" x2="12" y2="3"/>
            <line x1="12" y1="21" x2="12" y2="23"/>
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
            <line x1="1" y1="12" x2="3" y2="12"/>
            <line x1="21" y1="12" x2="23" y2="12"/>
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
          </svg>
        </button>
        <button v-if="reminderStore.hasUnreadReminders" class="reminder-btn reminder-btn-mobile" @click="handleReminderClose" title="临期博主提醒">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <span class="reminder-dot"></span>
        </button>
        <button v-if="reminderStore.hasUnreadReminders" class="mark-reminders-read-btn mark-reminders-read-btn-mobile" @click="handleMarkAllRemindersRead" title="全部已读">
          已读
        </button>
        <button class="notification-btn" @click="showNotifications = !showNotifications" :class="{ 'has-unread': unreadNotificationCount > 0 }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
            <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
          </svg>
          <span v-if="unreadNotificationCount > 0" class="notification-badge">{{ unreadNotificationCount > 99 ? '99+' : unreadNotificationCount }}</span>
        </button>
        <div class="user-info">
          <div class="user-avatar">
            {{ userStore.real_name?.charAt(0) || 'U' }}
          </div>
          <div class="user-detail">
            <span class="user-name">{{ userStore.real_name }}</span>
            <span class="user-role">{{ userStore.role === 'admin' ? '管理员' : '用户' }}</span>
          </div>
        </div>
        <button class="logout-btn" @click="handleLogout">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
            <polyline points="16,17 21,12 16,7"/>
            <line x1="21" y1="12" x2="9" y2="12"/>
          </svg>
        </button>
      </div>
    </header>

    <transition name="dropdown-fade">
    <div v-if="showNotifications" class="notification-dropdown" :class="{ 'dark-mode': themeStore.isDark }">
      <div class="notification-header">
        <h4>消息中心</h4>
        <div class="notification-tabs">
          <button :class="{ active: notificationTab === 'notifications' }" @click="notificationTab = 'notifications'">
            通知 <span v-if="unreadNotificationCount > 0" class="notification-badge">{{ unreadNotificationCount }}</span>
          </button>
          <button :class="{ active: notificationTab === 'announcements' }" @click="notificationTab = 'announcements'">
            公告
          </button>
          <button :class="{ active: notificationTab === 'messages' }" @click="notificationTab = 'messages'">
            私信 <span v-if="privateMessageUnreadCount > 0" class="notification-badge">{{ privateMessageUnreadCount }}</span>
          </button>
        </div>
      </div>
      <div v-if="notificationTab === 'notifications'" class="notification-search-bar">
        <input v-model="notificationSearch" type="text" placeholder="搜索通知..." class="notification-search-input" />
        <div class="notification-filter-buttons">
          <button :class="{ active: notificationFilter === 'all' }" @click="notificationFilter = 'all'">全部</button>
          <button :class="{ active: notificationFilter === 'unread' }" @click="notificationFilter = 'unread'">未读</button>
          <button :class="{ active: notificationFilter === 'important' }" @click="notificationFilter = 'important'">重要</button>
          <button :class="{ active: notificationFilter === 'team' }" @click="notificationFilter = 'team'">团队</button>
          <button :class="{ active: notificationFilter === 'blogger' }" @click="notificationFilter = 'blogger'">博主</button>
        </div>
        <div class="notification-batch-actions" v-if="selectedNotifications.length > 0">
          <button class="batch-btn" @click="batchMarkRead">批量已读</button>
          <button class="batch-btn delete" @click="batchDelete">批量删除</button>
          <button class="batch-btn cancel" @click="clearSelection">取消</button>
        </div>
      </div>
      <div class="notification-list">
        <div v-if="notificationTab === 'notifications'">
          <div v-if="filteredNotifications.length > 0">
            <div v-for="notif in filteredNotifications" :key="notif.id" class="notification-item" :class="{ unread: !notif.read, important: notif.priority === 'important', selected: selectedNotifications.includes(notif.id) }" @click="handleNotificationClick(notif)" @contextmenu.prevent="toggleNotificationSelection(notif.id)">
              <div class="notification-checkbox" @click.stop="toggleNotificationSelection(notif.id)">
                <input type="checkbox" :checked="selectedNotifications.includes(notif.id)" />
              </div>
              <div class="notification-icon" :class="notif.type">
                <span v-if="notif.type === 'mention'">@</span>
                <span v-else-if="notif.type === 'announcement'">📢</span>
                <span v-else-if="notif.type === 'system'">⚙️</span>
                <span v-else-if="notif.type === 'invalid_blogger'">⚠️</span>
                <span v-else-if="notif.type === 'countdown'">⏰</span>
                <span v-else-if="notif.type === 'team_change'">👥</span>
                <span v-else-if="notif.type === 'new_post'">📝</span>
                <span v-else-if="notif.type === 'post_reply'">💬</span>
                <span v-else-if="notif.type === 'post_like'">❤️</span>
                <span v-else-if="notif.type === 'team_message'">📨</span>
                <span v-else-if="notif.type === 'blogger_transfer'">🔄</span>
                <span v-else-if="notif.type === 'team_join'">🎉</span>
                <span v-else-if="notif.type === 'team_leave'">👋</span>
                <span v-else>📬</span>
              </div>
              <div class="notification-content">
                <div class="notification-title">{{ notif.title }}</div>
                <div class="notification-text">{{ notif.content }}</div>
                <div class="notification-time">{{ formatTime(notif.create_time) }}</div>
              </div>
              <div class="notification-actions" @click.stop>
                <button v-if="!notif.read" class="quick-action-btn" @click="markSingleRead(notif)" title="标记已读">
                  ✓
                </button>
                <button class="quick-action-btn delete" @click="deleteSingleNotification(notif)" title="删除">
                  🗑️
                </button>
              </div>
              <div v-if="!notif.read" class="unread-dot"></div>
              <div v-if="notif.priority === 'important'" class="important-indicator">重要</div>
            </div>
          </div>
          <div v-else class="notification-empty">
            <span>{{ notificationSearch ? '未找到匹配的通知' : '暂无通知' }}</span>
          </div>
        </div>
        <div v-else-if="notificationTab === 'announcements'">
          <div v-if="announcements.length > 0">
            <div v-for="announcement in announcements" :key="announcement.id" class="notification-item" @click="handleAnnouncementClick(announcement)">
              <div class="notification-icon announcement">
                📢
              </div>
              <div class="notification-content">
                <div class="notification-title">{{ announcement.title }}</div>
                <div class="notification-text">{{ announcement.content }}</div>
                <div class="notification-time">{{ formatTime(announcement.create_time) }}</div>
              </div>
            </div>
          </div>
          <div v-else class="notification-empty">
            <span>暂无公告</span>
          </div>
        </div>
        <div v-else-if="notificationTab === 'messages'">
          <div v-if="conversations.length > 0">
            <div v-for="conv in conversations" :key="conv.user_id" class="notification-item" @click="handleConversationClick(conv)">
              <div class="notification-icon">
                💬
              </div>
              <div class="notification-content">
                <div class="notification-title">{{ conv.username || conv.real_name || '用户' }}</div>
                <div class="notification-text">{{ conv.last_message || '暂无消息' }}</div>
                <div class="notification-time">{{ formatTime(conv.last_time) }}</div>
              </div>
              <div v-if="conv.unread_count > 0" class="unread-dot"></div>
            </div>
          </div>
          <div v-else class="notification-empty">
            <span>暂无私信</span>
          </div>
        </div>
      </div>
      <div v-if="notificationTab === 'notifications' && notifications.length > 0" class="notification-footer">
        <button class="mark-all-read" @click="markAllRead">全部已读</button>
        <span class="notification-count">{{ filteredNotifications.length }} 条</span>
      </div>
    </div>
    </transition>

    <transition name="modal-fade">
      <div v-if="showReminderPopup" class="reminder-modal-overlay" @click="showReminderPopup = false">
        <div class="reminder-modal" @click.stop>
          <div class="reminder-modal-header">
            <h3>⏰ 临期博主提醒</h3>
            <button class="reminder-modal-close" @click="handleReminderClose">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="reminder-modal-body">
            <p>您有 {{ reminderStore.expiringBloggers.length }} 位博主即将到期，且未填写联系方式，请及时处理！</p>
            <div class="reminder-list">
              <div v-for="blogger in reminderStore.expiringBloggers" :key="blogger.id" class="reminder-item" @click="router.push(`/blogger/${blogger.id}`); handleReminderClose()">
                <div class="reminder-item-avatar" :style="{ backgroundColor: blogger.color || '#3b82f6' }">
                  {{ blogger.nickname.charAt(0) }}
                </div>
                <div class="reminder-item-info">
                  <div class="reminder-item-name">{{ blogger.nickname }}</div>
                  <div class="reminder-item-days">剩余 {{ blogger.days_left }} 天</div>
                </div>
                <div class="reminder-item-arrow">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="9,6 15,12 9,18"/>
                  </svg>
                </div>
              </div>
            </div>
          </div>
          <div class="reminder-modal-footer">
            <button class="btn btn-secondary" @click="handleReminderClose">稍后处理</button>
            <button class="btn btn-primary" @click="handleMarkAllRemindersRead; handleReminderClose">全部已读</button>
          </div>
        </div>
      </div>
    </transition>

    <div class="mobile-menu" :class="{ active: showMenu }" @click="showMenu = false">
      <div class="mobile-menu-content" @click.stop>
        <div class="menu-section-label">核心功能</div>
        <router-link to="/" class="mobile-menu-item" :class="{ active: route.path === '/' }" @click="showMenu = false">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
            <polyline points="9,22 9,12 15,12 15,22"/>
          </svg>
          <span>首页</span>
        </router-link>
        <router-link to="/add" class="mobile-menu-item" :class="{ active: route.path === '/add' }" @click="showMenu = false">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="16"/>
            <line x1="8" y1="12" x2="16" y2="12"/>
          </svg>
          <span>录入博主</span>
        </router-link>
        <router-link to="/my" class="mobile-menu-item" :class="{ active: route.path === '/my' }" @click="showMenu = false">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          <span>我的</span>
        </router-link>
        <router-link to="/team" class="mobile-menu-item" :class="{ active: route.path === '/team' }" @click="showMenu = false">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
          <span>团队</span>
        </router-link>
        <div class="menu-section-label">更多工具</div>
        <router-link to="/statistics" class="mobile-menu-item" :class="{ active: route.path === '/statistics' }" @click="showMenu = false">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="20" x2="18" y2="10"/>
            <line x1="12" y1="20" x2="12" y2="4"/>
            <line x1="6" y1="20" x2="6" y2="14"/>
          </svg>
          <span>数据统计</span>
        </router-link>
        <router-link to="/calendar" class="mobile-menu-item" :class="{ active: route.path === '/calendar' }" @click="showMenu = false">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
            <line x1="16" y1="2" x2="16" y2="6"/>
            <line x1="8" y1="2" x2="8" y2="6"/>
            <line x1="3" y1="10" x2="21" y2="10"/>
          </svg>
          <span>日历日程</span>
        </router-link>
        <router-link to="/workflow" class="mobile-menu-item" :class="{ active: route.path === '/workflow' }" @click="showMenu = false">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
          </svg>
          <span>工作流</span>
        </router-link>
        <router-link to="/forums" class="mobile-menu-item" :class="{ active: route.path.startsWith('/forums') }" @click="showMenu = false">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
          <span>论坛社区</span>
        </router-link>
        <router-link to="/chat" class="mobile-menu-item" :class="{ active: route.path === '/chat' }" @click="showMenu = false">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
          <span>即时聊天</span>
        </router-link>
        <div class="menu-section-label" v-if="userStore.role === 'admin'">系统管理</div>
        <router-link v-if="userStore.role === 'admin'" to="/admin" class="mobile-menu-item" :class="{ active: route.path === '/admin' }" @click="showMenu = false">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/>
            <path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20z"/>
          </svg>
          <span>管理面板</span>
        </router-link>
      </div>
    </div>

    <main class="main">
      <router-view />
    </main>

    <nav class="mobile-nav" v-if="!isTeamPage">
      <router-link to="/" class="mobile-nav-item" :class="{ active: route.path === '/' }">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
          <polyline points="9,22 9,12 15,12 15,22"/>
        </svg>
        <span>首页</span>
      </router-link>
      <router-link to="/add" class="mobile-nav-item" :class="{ active: route.path === '/add' }">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="8" x2="12" y2="16"/>
          <line x1="8" y1="12" x2="16" y2="12"/>
        </svg>
        <span>录入</span>
      </router-link>
      <router-link to="/my" class="mobile-nav-item" :class="{ active: route.path === '/my' }">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
          <circle cx="12" cy="7" r="4"/>
        </svg>
        <span>我的</span>
      </router-link>
      <router-link to="/team" class="mobile-nav-item" :class="{ active: route.path === '/team' }">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
          <circle cx="9" cy="7" r="4"/>
          <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
          <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
        </svg>
        <span>团队</span>
      </router-link>
    </nav>
  </div>
</template>

<script setup>
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useThemeStore } from '../stores/theme'
import { useReminderStore } from '../stores/reminder'
import { useToastStore } from '../stores/toast'
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { getNotificationsAPI, markNotificationReadAPI, markAllNotificationsReadAPI, deleteNotificationAPI, batchDeleteNotificationsAPI, getExpiringBloggersWithoutContactAPI, getAnnouncementsAPI, getConversationsAPI, getUnreadCountAPI } from '../api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const themeStore = useThemeStore()
const reminderStore = useReminderStore()
const toastStore = useToastStore()

const showReminderPopup = ref(false)
const showMenu = ref(false)
const showNotifications = ref(false)
const notifications = ref([])
const announcements = ref([])
const conversations = ref([])
const notificationTimer = ref(null)
const notificationTab = ref('notifications')
const privateMessageUnreadCount = ref(0)
const notificationSearch = ref('')
const notificationFilter = ref('all')
const selectedNotifications = ref([])
const newNotificationSound = ref(null)

const unreadNotificationCount = computed(() => notifications.value.filter(n => !n.read).length)
const importantNotificationCount = computed(() => notifications.value.filter(n => n.priority === 'important' && !n.read).length)
const teamNotificationCount = computed(() => notifications.value.filter(n => ['new_post', 'post_reply', 'post_like', 'team_message', 'team_change', 'team_join', 'team_leave', 'mention'].includes(n.type) && !n.read).length)
const bloggerNotificationCount = computed(() => notifications.value.filter(n => ['invalid_blogger', 'countdown', 'blogger_transfer'].includes(n.type) && !n.read).length)

const filteredNotifications = computed(() => {
  let result = notifications.value
  if (notificationSearch.value) {
    const search = notificationSearch.value.toLowerCase()
    result = result.filter(n => 
      n.title?.toLowerCase().includes(search) || 
      n.content?.toLowerCase().includes(search)
    )
  }
  if (notificationFilter.value === 'unread') {
    result = result.filter(n => !n.read)
  } else if (notificationFilter.value === 'important') {
    result = result.filter(n => n.priority === 'important')
  } else if (notificationFilter.value === 'team') {
    result = result.filter(n => ['new_post', 'post_reply', 'post_like', 'team_message', 'team_change', 'team_join', 'team_leave', 'mention'].includes(n.type))
  } else if (notificationFilter.value === 'blogger') {
    result = result.filter(n => ['invalid_blogger', 'countdown', 'blogger_transfer'].includes(n.type))
  }
  return result
})


const loadAnnouncements = async () => {
  try {
    const res = await getAnnouncementsAPI()
    if (res && res.code === 200) {
      announcements.value = res.data || []
    }
  } catch (e) {
    announcements.value = []
  }
}

const loadConversations = async () => {
  try {
    const res = await getConversationsAPI()
    if (res.code === 200) {
      conversations.value = res.data || []
    }
    const unreadRes = await getUnreadCountAPI()
    if (unreadRes.code === 200) {
      privateMessageUnreadCount.value = unreadRes.data?.total || 0
    }
  } catch (e) {
    toastStore.error('加载对话失败', '请稍后重试')
  }
}

const isTeamPage = computed(() => {
  return route.name === 'TeamHome' || route.path.startsWith('/team/')
})

const loadExpiringBloggers = async () => {
  try {
    const res = await getExpiringBloggersWithoutContactAPI()
    if (res.code === 200) {
      reminderStore.setExpiringBloggers(res.data || [])
      
      if (reminderStore.hasExpiringBloggers && !reminderStore.hasShownLoginPopup) {
        showReminderPopup.value = true
      }
    }
  } catch (e) {
    toastStore.error('加载临期博主提醒失败', '请稍后重试')
  }
}

const loadNotifications = async () => {
  try {
    const res = await getNotificationsAPI()
    if (res.code === 200) {
      notifications.value = res.data || []
    }
  } catch (e) {
    toastStore.error('加载通知失败', '请稍后重试')
  }
}

const loadAllMessageCenterData = async () => {
  await Promise.all([
    loadNotifications(),
    loadAnnouncements(),
    loadConversations()
  ])
}

const markAllRead = async () => {
  try {
    await markAllNotificationsReadAPI()
    notifications.value.forEach(n => n.read = true)
  } catch (e) {
    console.error('标记已读失败', e)
  }
}

const handleNotificationClick = async (notif) => {
  if (!notif.read) {
    try {
      await markNotificationReadAPI(notif.id)
      notif.read = true
    } catch (e) {
      console.error('标记已读失败', e)
    }
  }
  showNotifications.value = false
  if (notif.type === 'invalid_blogger' && notif.blogger_id) {
    router.push(`/blogger/${notif.blogger_id}`)
  } else if (notif.type === 'countdown' && notif.blogger_id) {
    router.push(`/blogger/${notif.blogger_id}`)
  } else if (notif.type === 'new_post' && notif.team_id && notif.post_id) {
    router.push(`/team?teamId=${notif.team_id}&postId=${notif.post_id}`)
  } else if (notif.type === 'post_reply' && notif.team_id && notif.post_id) {
    router.push(`/team?teamId=${notif.team_id}&postId=${notif.post_id}`)
  } else if (notif.type === 'post_like' && notif.team_id && notif.post_id) {
    router.push(`/team?teamId=${notif.team_id}&postId=${notif.post_id}`)
  } else if (notif.type === 'mention' && notif.team_id && notif.post_id) {
    router.push(`/team?teamId=${notif.team_id}&postId=${notif.post_id}`)
  } else if (notif.type === 'team_message' && notif.team_id) {
    router.push(`/team?teamId=${notif.team_id}&tab=message`)
  } else if (notif.type === 'blogger_transfer' && notif.blogger_id) {
    router.push(`/blogger/${notif.blogger_id}`)
  } else if (notif.team_id) {
    router.push(`/team?teamId=${notif.team_id}`)
  }
}

const handleConversationClick = (conv) => {
  showNotifications.value = false
  router.push('/chat?userId=' + conv.user_id)
}

const handleAnnouncementClick = (announcement) => {
  showNotifications.value = false
}

const markSingleRead = async (notif) => {
  if (!notif.read) {
    try {
      await markNotificationReadAPI(notif.id)
      notif.read = true
    } catch (e) {
      console.error('标记已读失败', e)
    }
  }
}

const deleteSingleNotification = async (notif) => {
  try {
    await deleteNotificationAPI(notif.id)
    notifications.value = notifications.value.filter(n => n.id !== notif.id)
  } catch (e) {
    console.error('删除通知失败', e)
  }
}

const toggleNotificationSelection = (id) => {
  const index = selectedNotifications.value.indexOf(id)
  if (index === -1) {
    selectedNotifications.value.push(id)
  } else {
    selectedNotifications.value.splice(index, 1)
  }
}

const batchMarkRead = async () => {
  try {
    for (const id of selectedNotifications.value) {
      const notif = notifications.value.find(n => n.id === id)
      if (notif && !notif.read) {
        await markNotificationReadAPI(id)
        notif.read = true
      }
    }
    clearSelection()
  } catch (e) {
    console.error('批量标记已读失败', e)
  }
}

const batchDelete = async () => {
  try {
    await batchDeleteNotificationsAPI(selectedNotifications.value)
    notifications.value = notifications.value.filter(n => !selectedNotifications.value.includes(n.id))
    clearSelection()
  } catch (e) {
    console.error('批量删除失败', e)
  }
}

const clearSelection = () => {
  selectedNotifications.value = []
}

const formatTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  return `${Math.floor(diff / 86400000)}天前`
}

const handleClickOutside = (e) => {
  if (showNotifications.value && !e.target.closest('.notification-btn') && !e.target.closest('.notification-dropdown')) {
    showNotifications.value = false
  }
}

onMounted(() => {
  if (userStore.isLoggedIn) {
    reminderStore.loadFromStorage()
    loadAllMessageCenterData()
    loadExpiringBloggers()
    notificationTimer.value = setInterval(loadAllMessageCenterData, 30000)
  }
  document.addEventListener('click', handleClickOutside)
})

watch(showNotifications, (newVal) => {
  if (newVal && userStore.isLoggedIn) {
    loadAllMessageCenterData()
  }
})

onUnmounted(() => {
  if (notificationTimer.value) {
    clearInterval(notificationTimer.value)
  }
  document.removeEventListener('click', handleClickOutside)
})

const handleReminderClose = () => {
  showReminderPopup.value = false
  reminderStore.markAsRead()
}

const handleMarkAllRemindersRead = () => {
  reminderStore.markAllAsRead()
  loadExpiringBloggers()
}

const handleLogout = () => {
  userStore.logout()
  reminderStore.resetForNewLogin()
  router.push('/login')
}

const openCommandPalette = () => {
  window.dispatchEvent(new CustomEvent('open-command-palette'))
}
</script>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  height: 64px;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.menu-toggle {
  display: none;
  width: 40px;
  height: 40px;
  background: transparent;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  align-items: center;
  justify-content: center;
  color: var(--text-primary);
  transition: background 0.3s;
}

.menu-toggle:hover {
  background: var(--bg-card-hover);
}

.menu-toggle svg {
  width: 24px;
  height: 24px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  transition: opacity 0.3s;
}

.logo:hover {
  opacity: 0.8;
}

.logo-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.logo-icon svg {
  width: 20px;
  height: 20px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.nav {
  display: flex;
  gap: 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border-radius: 10px;
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
  text-decoration: none;
}

.nav-item svg {
  width: 18px;
  height: 18px;
}

.nav-item:hover {
  background: var(--bg-card-hover);
  color: var(--text-primary);
}

.nav-item.active {
  background: rgba(249, 115, 22, 0.1);
  color: var(--primary);
}

.nav-item-highlight {
  background: linear-gradient(135deg, #ff6b35, #f7931e);
  color: white;
}

.nav-item-highlight:hover {
  background: linear-gradient(135deg, #ff5722, #e64a19);
  color: white;
}

.nav-item-highlight.active {
  background: linear-gradient(135deg, #ff5722, #e64a19);
  color: white;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.command-palette-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.2s;
  font-size: 13px;
}

.command-palette-trigger:hover {
  border-color: var(--primary);
  color: var(--text-secondary);
  background: var(--bg-card-hover);
}

.command-palette-trigger svg {
  width: 16px;
  height: 16px;
}

.command-palette-trigger .trigger-text {
  font-size: 13px;
}

.command-palette-trigger kbd {
  padding: 1px 6px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 11px;
  font-family: inherit;
  color: var(--text-muted);
}

.theme-toggle {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-card-hover);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s ease;
}

.theme-toggle:hover {
  background: var(--bg-card);
  color: var(--text-primary);
  border-color: var(--primary);
}

.theme-toggle svg {
  width: 18px;
  height: 18px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 14px;
  box-shadow: 0 2px 8px rgba(249, 115, 22, 0.2);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.user-avatar:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.user-detail {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.user-role {
  font-size: 12px;
  color: var(--text-muted);
}

.logout-btn {
  width: 36px;
  height: 36px;
  background: rgba(239, 68, 68, 0.1);
  border: none;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ef4444;
  cursor: pointer;
  transition: all 0.3s ease;
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.15);
}

.logout-btn svg {
  width: 18px;
  height: 18px;
}

.notification-btn {
  width: 36px;
  height: 36px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  margin-right: 8px;
}

.notification-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.notification-btn.has-unread {
  color: #f59e0b;
}

.notification-btn svg {
  width: 18px;
  height: 18px;
}

.notification-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background: #ef4444;
  color: white;
  font-size: 10px;
  font-weight: bold;
  min-width: 16px;
  height: 16px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 4px;
}

.notification-dropdown {
  position: fixed;
  top: 64px;
  right: 80px;
  width: 360px;
  max-height: 480px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.notification-dropdown.dark-mode {
  background: var(--bg-card);
}

.notification-header {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;
}

.notification-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.notification-tabs {
  display: flex;
  gap: 4px;
  flex: 1;
  justify-content: center;
}

.notification-tabs button {
  padding: 4px 12px;
  border: none;
  background: var(--bg-secondary);
  border-radius: 16px;
  font-size: 12px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 4px;
}

.notification-tabs button.active {
  background: var(--primary);
  color: white;
}

.notification-tabs button span {
  font-size: 10px;
}

.important-badge {
  background: #ef4444;
  color: white;
  padding: 1px 5px;
  border-radius: 8px;
  font-size: 10px;
}

.mark-all-read {
  background: none;
  border: none;
  color: #3b82f6;
  font-size: 12px;
  cursor: pointer;
}

.notification-search-bar {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-search-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 13px;
  outline: none;
  transition: all 0.2s;
}

.notification-search-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

.notification-filter-buttons {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.notification-filter-buttons button {
  padding: 4px 10px;
  border: 1px solid var(--border-color);
  background: var(--bg-card);
  border-radius: 14px;
  font-size: 11px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.notification-filter-buttons button.active {
  background: var(--primary);
  border-color: var(--primary);
  color: white;
}

.notification-filter-buttons button:hover {
  border-color: var(--primary);
}

.notification-batch-actions {
  display: flex;
  gap: 8px;
  padding-top: 8px;
  border-top: 1px solid var(--border-color);
}

.batch-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.batch-btn:not(.delete):not(.cancel) {
  background: var(--primary);
  color: white;
}

.batch-btn.delete {
  background: #ef4444;
  color: white;
}

.batch-btn.cancel {
  background: var(--bg-secondary);
  color: var(--text-secondary);
}

.batch-btn:hover {
  opacity: 0.9;
}

.notification-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.notification-count {
  font-size: 12px;
  color: var(--text-muted);
}

.notification-checkbox {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.notification-checkbox input {
  cursor: pointer;
}

.notification-item.selected {
  background: rgba(249, 115, 22, 0.08);
}

.notification-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.notification-item:hover .notification-actions {
  opacity: 1;
}

.quick-action-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: var(--bg-secondary);
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.quick-action-btn:hover {
  background: var(--bg-hover);
}

.quick-action-btn.delete:hover {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.mark-all-read:hover {
  text-decoration: underline;
}

.notification-list {
  flex: 1;
  overflow-y: auto;
  max-height: 400px;
}

.notification-item {
  padding: 12px 16px;
  display: flex;
  align-items: flex-start;
  gap: 12px;
  cursor: pointer;
  transition: background 0.2s;
  position: relative;
}

.notification-item:hover {
  background: var(--bg-hover);
}

.notification-item.unread {
  background: rgba(59, 130, 246, 0.05);
}

.notification-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.notification-text {
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.notification-time {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 4px;
}

.unread-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #3b82f6;
  position: absolute;
  top: 50%;
  right: 16px;
  transform: translateY(-50%);
}

.notification-item.important {
  background: rgba(239, 68, 68, 0.05);
  border-left: 3px solid #ef4444;
}

.notification-item.important.unread {
  background: rgba(239, 68, 68, 0.08);
}

.notification-icon.mention {
  background: rgba(99, 102, 241, 0.15);
  color: var(--purple);
}

.notification-icon.announcement {
  background: rgba(249, 115, 22, 0.15);
  color: var(--primary);
}

.notification-icon.system {
  background: rgba(59, 130, 246, 0.15);
  color: var(--info);
}

.important-indicator {
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 10px;
  color: #ef4444;
  font-weight: 600;
  background: rgba(239, 68, 68, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
}

.notification-empty {
  padding: 40px 16px;
  text-align: center;
  color: var(--text-muted);
}

.dropdown-fade-enter-active,
.dropdown-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.dropdown-fade-enter-from,
.dropdown-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.main {
  flex: 1;
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}

.mobile-nav {
  display: none;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 68px;
  background-color: var(--bg-card) !important;
  background-image: none !important;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-top: 1px solid var(--border-color);
  justify-content: space-between;
  align-items: center;
  padding: 0 8px;
  z-index: 100;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.08);
}

[data-theme="dark"] .mobile-nav {
  background-color: var(--bg-card) !important;
  background-image: none !important;
}

.mobile-nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 6px 16px;
  min-width: 60px;
  color: var(--text-muted);
  text-decoration: none;
  transition: all 0.2s ease;
  border-radius: 8px;
  min-width: 56px;
}

.mobile-nav-item:active {
  background: var(--bg-secondary);
  transform: scale(0.95);
}

.mobile-nav-item svg {
  width: 24px;
  height: 24px;
}

.mobile-nav-item span {
  font-size: 10px;
  font-weight: 500;
  line-height: 1.2;
}

.mobile-nav-item.active {
  color: var(--primary);
  background: rgba(249, 115, 22, 0.08);
}

.mobile-menu {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 200;
  opacity: 0;
  visibility: hidden;
  transition: all 0.3s ease;
}

.mobile-menu.active {
  opacity: 1;
  visibility: visible;
}

.mobile-menu-content {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 280px;
  background: var(--bg-card);
  padding: 80px 20px 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  transform: translateX(100%);
  transition: transform 0.3s ease;
  overflow-y: auto;
}

.mobile-menu.active .mobile-menu-content {
  transform: translateX(0);
}

.menu-section-label {
  padding: 12px 20px 8px;
  font-size: 11px;
  font-weight: 700;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
  border-top: 1px solid var(--border-light);
  margin-top: 4px;
}

.menu-section-label:first-child {
  border-top: none;
  margin-top: 0;
}

.mobile-menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-radius: 12px;
  color: var(--text-secondary);
  font-size: 15px;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s ease;
}

.mobile-menu-item svg {
  width: 22px;
  height: 22px;
}

.mobile-menu-item:hover {
  background: var(--bg-card-hover);
}

.mobile-menu-item.active {
  background: rgba(249, 115, 22, 0.1);
  color: var(--primary);
}

@media (max-width: 768px) {
  .layout .header {
    padding: 0 12px;
  }

  .menu-toggle {
    display: flex;
    width: 36px;
    height: 36px;
  }

  .nav {
    display: none;
  }

  .mobile-nav {
    display: flex;
  }

  .mobile-menu {
    display: block;
  }

  .main {
    padding: 16px;
    padding-bottom: 92px;
  }

  .header-right .user-detail {
    display: none;
  }

  .command-palette-trigger .trigger-text {
    display: none;
  }

  .command-palette-trigger kbd {
    display: none;
  }

  /* 优化右侧按钮间距 */
  .header-right {
    gap: 8px;
  }

  .reminder-btn-mobile {
    display: none;
  }

  .mark-reminders-read-btn-mobile {
    display: none;
  }

  .notification-btn {
    width: 32px;
    height: 32px;
    margin-right: 4px;
  }

  .notification-btn svg {
    width: 16px;
    height: 16px;
  }

  .theme-toggle {
    width: 32px;
    height: 32px;
  }

  .theme-toggle svg {
    width: 16px;
    height: 16px;
  }

  .user-avatar {
    width: 32px;
    height: 32px;
    font-size: 13px;
  }

  .logout-btn {
    width: 32px;
    height: 32px;
  }

  .logout-btn svg {
    width: 16px;
    height: 16px;
  }

  .notification-dropdown {
    right: 12px;
    width: calc(100vw - 24px);
    max-width: 360px;
  }
}

@media (max-width: 480px) {
  .logo-text {
    display: none;
  }

  .layout .header {
    height: 56px;
    padding: 0 12px !important;
    min-width: 100%;
  }

  .header-left {
    gap: 6px;
    flex-shrink: 0;
  }

  .menu-toggle {
    width: 32px;
    height: 32px;
    flex-shrink: 0;
  }
  
  .logo-icon {
    width: 30px;
    height: 30px;
    flex-shrink: 0;
  }
  
  .header-right {
    gap: 6px;
    flex-shrink: 0;
  }

  .menu-toggle svg {
    width: 20px;
    height: 20px;
  }

  .logo-icon {
    width: 32px;
    height: 32px;
  }

  .logo-icon svg {
    width: 18px;
    height: 18px;
  }

  .user-avatar {
    width: 28px;
    height: 28px;
    font-size: 12px;
  }

  .notification-btn {
    width: 28px;
    height: 28px;
    margin-right: 4px;
  }

  .notification-btn svg {
    width: 14px;
    height: 14px;
  }

  .theme-toggle {
    width: 28px;
    height: 28px;
  }

  .theme-toggle svg {
    width: 14px;
    height: 14px;
  }

  .logout-btn {
    width: 28px;
    height: 28px;
  }

  .logout-btn svg {
    width: 14px;
    height: 14px;
  }

  .header-right {
    gap: 4px;
  }

  .notification-dropdown {
    top: 56px;
    right: 8px;
    width: calc(100vw - 16px);
  }

  .theme-toggle {
    width: 32px;
    height: 32px;
  }

  .notification-dropdown {
    position: fixed;
    top: 56px;
    right: 0;
    left: 0;
    width: 100%;
    max-height: calc(100vh - 56px);
    border-radius: 0;
    border-left: none;
    border-right: none;
  }
}

.reminder-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background: transparent;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  color: var(--text-primary);
  transition: background 0.3s;
  position: relative;
}

.reminder-btn:hover {
  background: var(--bg-card-hover);
}

.reminder-btn svg {
  width: 20px;
  height: 20px;
}

.reminder-dot {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 8px;
  height: 8px;
  background: #ef4444;
  border-radius: 50%;
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.mark-reminders-read-btn {
  padding: 6px 12px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: background 0.3s;
}

.mark-reminders-read-btn:hover {
  background: var(--primary-dark);
}

.reminder-modal-overlay {
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
}

.reminder-modal {
  background: var(--bg-card);
  border-radius: 16px;
  width: 90%;
  max-width: 480px;
  max-height: 80vh;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.reminder-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.reminder-modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.reminder-modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: transparent;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.3s;
}

.reminder-modal-close:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.reminder-modal-body {
  padding: 24px;
  max-height: 50vh;
  overflow-y: auto;
}

.reminder-modal-body p {
  margin: 0 0 20px 0;
  color: var(--text-secondary);
  line-height: 1.5;
}

.reminder-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.reminder-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--bg-hover);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s;
}

.reminder-item:hover {
  background: var(--primary);
  transform: translateY(-2px);
}

.reminder-item:hover .reminder-item-name,
.reminder-item:hover .reminder-item-days,
.reminder-item:hover .reminder-item-arrow {
  color: white;
}

.reminder-item-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 16px;
  flex-shrink: 0;
}

.reminder-item-info {
  flex: 1;
  min-width: 0;
}

.reminder-item-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.reminder-item-days {
  font-size: 12px;
  color: #ef4444;
  font-weight: 500;
}

.reminder-item:hover .reminder-item-days {
  color: white;
}

.reminder-item-arrow {
  color: var(--text-secondary);
  flex-shrink: 0;
}

.reminder-item-arrow svg {
  width: 16px;
  height: 16px;
}

.reminder-modal-footer {
  display: flex;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid var(--border-color);
}

.reminder-modal-footer .btn {
  flex: 1;
}

.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: all 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-from .reminder-modal,
.modal-fade-leave-to .reminder-modal {
  transform: scale(0.9);
}

.modal-fade-enter-active .reminder-modal,
.modal-fade-leave-active .reminder-modal {
  transition: transform 0.3s ease;
}
</style>