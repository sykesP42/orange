<template>
  <div class="team-home-page" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>正在加载小组数据...</p>
    </div>

    <div v-else-if="!team" class="not-found-state">
      <div class="not-found-icon">?</div>
      <h2>小组不存在</h2>
      <p>该小组可能已被删除或您没有访问权限</p>
      <button class="btn btn-primary" @click="goBack">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="19" y1="12" x2="5" y2="12"/>
          <polyline points="12 19 5 12 12 5"/>
        </svg>
        返回团队中心
      </button>
    </div>

    <div v-else-if="isPrivate && !isTeamMember" class="privacy-error-state">
      <div class="privacy-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
          <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
        </svg>
      </div>
      <h2>页面隐私设置</h2>
      <p>该小组的可见性被设置为私密</p>
      <p class="privacy-hint">请联系小组管理员申请加入</p>
      <button class="btn btn-primary" @click="goBack">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="19" y1="12" x2="5" y2="12"/>
          <polyline points="12 19 5 12 12 5"/>
        </svg>
        返回团队中心
      </button>
    </div>

    <div v-else>
      <div class="back-btn-fixed" @click="handleBackBtnClick" title="返回团队中心">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="19" y1="12" x2="5" y2="12"/>
          <polyline points="12 19 5 12 12 5"/>
        </svg>
      </div>

      <div class="team-banner" :style="bannerStyle">
        <div class="banner-particles"></div>
        <div class="banner-overlay"></div>
        <div class="banner-content">
          <div class="banner-left">
            <div class="team-avatar" :style="{ backgroundColor: team.color }">
              <img v-if="team.logo" :src="team.logo" :alt="team.name" />
              <span v-else>{{ team.name.charAt(0) }}</span>
              <div v-if="isLeaderOrAdmin" class="admin-badge">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                </svg>
              </div>
            </div>
            <div class="team-intro">
              <h1 class="team-title">{{ team.name }}</h1>
              <p class="team-subtitle">{{ team.description || '暂无小组描述' }}</p>
              <div class="team-tags">
                <span class="tag">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                    <circle cx="9" cy="7" r="4"/>
                  </svg>
                  {{ members.length }} 名成员
                </span>
                <span class="tag">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                    <line x1="16" y1="2" x2="16" y2="6"/>
                    <line x1="8" y1="2" x2="8" y2="6"/>
                    <line x1="3" y1="10" x2="21" y2="10"/>
                  </svg>
                  {{ formatDate(team.create_time) }} 创建
                </span>
              </div>
            </div>
          </div>
          <div class="banner-right">
            <div v-if="team.announcement" class="announcement-marquee">
              <div class="marquee-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 17H2a3 3 0 0 0 3-3V9a7 7 0 0 1 14 0v5a3 3 0 0 0 3 3zm-8.27 4a2 2 0 0 1-3.46 0"/>
                </svg>
                公告
              </div>
              <div class="marquee-content">
                <span>{{ team.announcement }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="team-container">
        <nav class="sidebar" :class="{ collapsed: sidebarCollapsed }">
          <div class="sidebar-toggle" @click="sidebarCollapsed = !sidebarCollapsed">
            <svg v-if="sidebarCollapsed" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="15 18 9 12 15 6"/>
            </svg>
          </div>

          <div class="nav-section">
            <div class="nav-section-title" v-if="!sidebarCollapsed">导航</div>
            <a class="nav-item" :class="{ active: activeTab === 'home' }" @click="activeTab = 'home'">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
                <polyline points="9 22 9 12 15 12 15 22"/>
              </svg>
              <span v-if="!sidebarCollapsed">小组首页</span>
            </a>
            <a class="nav-item" :class="{ active: activeTab === 'members' }" @click="activeTab = 'members'">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              </svg>
              <span v-if="!sidebarCollapsed">成员列表</span>
              <span class="badge" v-if="!sidebarCollapsed">{{ members.length }}</span>
            </a>
            <a class="nav-item" :class="{ active: activeTab === 'bloggers' }" @click="activeTab = 'bloggers'">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
              <span v-if="!sidebarCollapsed">对接博主</span>
              <span class="badge" v-if="!sidebarCollapsed">{{ teamStats.total || 0 }}</span>
            </a>
            <a class="nav-item" :class="{ active: activeTab === 'forum' }" @click="activeTab = 'forum'">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              </svg>
              <span v-if="!sidebarCollapsed">讨论区</span>
              <span class="badge" v-if="!sidebarCollapsed">{{ posts.length }}</span>
            </a>
            <a class="nav-item" :class="{ active: activeTab === 'message' }" @click="switchToMessage">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 17H2a3 3 0 0 0 3-3V9a7 7 0 0 1 14 0v5a3 3 0 0 0 3 3zm-8.27 4a2 2 0 0 1-3.46 0"/>
              </svg>
              <span v-if="!sidebarCollapsed">消息中心</span>
              <span class="badge badge-danger" v-if="!sidebarCollapsed && unreadCount > 0">{{ unreadCount }}</span>
            </a>
            <a class="nav-item" @click="showAllUsersDialog">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              </svg>
              <span v-if="!sidebarCollapsed">全部用户</span>
            </a>
          </div>

          <div class="nav-section" v-if="isTeamAdmin">
            <div class="nav-section-title" v-if="!sidebarCollapsed">管理</div>
            <a class="nav-item" :class="{ active: activeTab === 'edit' }" @click="showEditDialog">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
              </svg>
              <span v-if="!sidebarCollapsed">编辑资料</span>
            </a>
            <a class="nav-item" :class="{ active: activeTab === 'announcement' }" @click="showAnnouncementDialog">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
1                <path d="M22 17H2a3 3 0 0 0 3-3V9a7 7 0 0 1 14 0v5a3 3 0 0 0 3 3zm-8.27 4a2 2 0 0 1-3.46 0"/>
              </svg>
              <span v-if="!sidebarCollapsed">发布公告</span>
            </a>
            <a class="nav-item" :class="{ active: activeTab === 'settings' }" @click="activeTab = 'settings'">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="3"/>
                <path d="M12 1v6m0 6v6m-5.5-13.5l4.24 4.24m4.52 4.52L20 21M1 12h6m6 0h6M6.34 17.66l4.24-4.24m4.52-4.52L20 3"/>
              </svg>
              <span v-if="!sidebarCollapsed">管理设置</span>
            </a>
          </div>

          <div class="nav-section nav-bottom">
            <div class="nav-section-title" v-if="!sidebarCollapsed">操作</div>
            <a v-if="!isMember && !isLeaderOrAdmin" class="nav-item action" @click="joinTeam">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="8.5" cy="7" r="4"/>
                <line x1="20" y1="8" x2="20" y2="14"/>
                <line x1="23" y1="11" x2="17" y2="11"/>
              </svg>
              <span v-if="!sidebarCollapsed">加入小组</span>
            </a>
            <a v-if="isMember && !isLeaderOrAdmin" class="nav-item danger" @click="leaveTeam">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                <polyline points="16 17 21 12 16 7"/>
                <line x1="21" y1="12" x2="9" y2="12"/>
              </svg>
              <span v-if="!sidebarCollapsed">退出小组</span>
            </a>
            <a class="nav-item" @click="handleBackBtnClick">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="19" y1="12" x2="5" y2="12"/>
                <polyline points="12 19 5 12 12 5"/>
              </svg>
              <span v-if="!sidebarCollapsed">返回中心</span>
            </a>
          </div>
        </nav>

        <main class="main-content">
          <div v-if="activeTab === 'home'" class="tab-content home-tab">
            <div class="content-card reveal" style="--delay: 0">
              <div class="card-header">
                <h3>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                    <circle cx="9" cy="7" r="4"/>
                    <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                    <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                  </svg>
                  小组成员
                </h3>
                <span class="count-badge">{{ members.length }} 人</span>
              </div>
              <div class="members-showcase">
                <div v-for="(member, index) in members.slice(0, 8)" :key="member.id"
                     class="member-avatar-lg reveal" :style="{ '--delay': index * 0.05 }"
                     :class="{ 'is-leader': member.id === team.leader_id }">
                  <img v-if="member.avatar" :src="member.avatar" :alt="member.real_name" v-avatar />
                  <span v-else>{{ member.real_name?.charAt(0) || 'U' }}</span>
                  <div class="member-tooltip">
                    <strong>{{ member.real_name }}</strong>
                    <span v-if="member.id === team.leader_id" class="leader-tag">组长</span>
                    <span>@{{ member.username }}</span>
                  </div>
                </div>
                <div v-if="members.length > 8" class="more-members">
                  +{{ members.length - 8 }}
                </div>
              </div>
            </div>

            <div class="content-card reveal" style="--delay: 0.1">
              <div class="card-header">
                <h3>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="18" y1="20" x2="18" y2="10"/>
                    <line x1="12" y1="20" x2="12" y2="4"/>
                    <line x1="6" y1="20" x2="6" y2="14"/>
                  </svg>
                  小组数据
                </h3>
              </div>
              <div class="stats-grid">
                <div class="stat-box">
                  <div class="stat-icon blue">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                      <circle cx="12" cy="7" r="4"/>
                    </svg>
                  </div>
                  <div class="stat-info">
                    <span class="stat-number">{{ teamStats.total || 0 }}</span>
                    <span class="stat-text">对接博主</span>
                  </div>
                </div>
                <div class="stat-box">
                  <div class="stat-icon green">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/>
                      <polyline points="17 6 23 6 23 12"/>
                    </svg>
                  </div>
                  <div class="stat-info">
                    <span class="stat-number">{{ teamStats.today || 0 }}</span>
                    <span class="stat-text">今日新增</span>
                  </div>
                </div>
                <div class="stat-box">
                  <div class="stat-icon orange">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                      <line x1="16" y1="2" x2="16" y2="6"/>
                      <line x1="8" y1="2" x2="8" y2="6"/>
                      <line x1="3" y1="10" x2="21" y2="10"/>
                    </svg>
                  </div>
                  <div class="stat-info">
                    <span class="stat-number">{{ teamStats.month || 0 }}</span>
                    <span class="stat-text">本月新增</span>
                  </div>
                </div>
                <div class="stat-box">
                  <div class="stat-icon purple">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                      <circle cx="9" cy="7" r="4"/>
                    </svg>
                  </div>
                  <div class="stat-info">
                    <span class="stat-number">{{ teamStats.memberCount || 0 }}</span>
                    <span class="stat-text">参与成员</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="content-card reveal" style="--delay: 0.2" v-if="teamBloggerList.length > 0">
              <div class="card-header">
                <h3>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                  </svg>
                  对接博主
                </h3>
                <span class="count-badge">{{ teamBloggerList.length }}</span>
              </div>
              <div class="blogger-grid">
                <div v-for="(blogger, index) in teamBloggerList.slice(0, 12)" :key="blogger.id"
                     class="blogger-mini-card reveal" :style="{ '--delay': 0.3 + index * 0.03 }"
                     @click="goToBloggerDetail(blogger.id)">
                  <div class="mini-avatar" :style="{ backgroundColor: getCategoryColor(blogger.category) }">
                    <img v-if="blogger.avatar" :src="blogger.avatar" :alt="blogger.nickname" v-avatar />
                    <span v-else>{{ blogger.nickname?.charAt(0) || '?' }}</span>
                  </div>
                  <div class="mini-info">
                    <span class="mini-name">{{ blogger.nickname }}</span>
                    <span class="mini-category" :style="{ color: getCategoryColor(blogger.category) }">
                      {{ blogger.category || '其他' }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="activeTab === 'members'" class="tab-content members-tab">
            <div class="content-card">
              <div class="card-header">
                <h3>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                    <circle cx="9" cy="7" r="4"/>
                    <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                    <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                  </svg>
                  全部成员
                </h3>
                <span class="count-badge">{{ members.length }} 人</span>
                <div class="batch-actions" v-if="isTeamAdmin">
                  <button class="btn btn-outline btn-sm" @click.stop="startMemberSelection" v-if="!isSelectingMembers && selectedMembers.length === 0">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                      <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                      <circle cx="8.5" cy="7" r="4"/>
                      <line x1="20" y1="8" x2="20" y2="14"/>
                      <line x1="23" y1="11" x2="17" y2="11"/>
                    </svg>
                    批量管理
                  </button>
                  <template v-if="selectedMembers.length > 0">
                    <span class="selected-count">已选择 {{ selectedMembers.length }} 人</span>
                    <button class="btn btn-secondary btn-sm" @click="openBatchMessageDialog">发送消息</button>
                    <button class="btn btn-danger btn-sm" @click="batchRemoveMembers">移出组员</button>
                    <button class="btn btn-text btn-sm" @click="clearSelection">清除</button>
                  </template>
                </div>
              </div>
              <div class="members-list" @click="handleMembersAreaClick">
                <div v-for="member in members" :key="member.id" class="member-list-item" :class="{ selected: selectedMembers.includes(member.id), 'selectable': isSelectingMembers }" @click.stop="toggleMemberSelection(member)">
                  <div class="member-checkbox" v-if="isSelectingMembers && member.id !== team.leader_id">
                    <input type="checkbox" :checked="selectedMembers.includes(member.id)" @click.stop />
                  </div>
                  <div class="member-avatar-md" :style="{ backgroundColor: team.color }">
                    <img v-if="member.avatar" :src="member.avatar" :alt="member.real_name" v-avatar />
                    <span v-else>{{ member.real_name?.charAt(0) || 'U' }}</span>
                    <span v-if="member.id === team.leader_id" class="leader-indicator"></span>
                  </div>
                  <div class="member-details">
                    <div class="member-name-row">
                      <span class="member-name">{{ member.real_name }}</span>
                      <span v-if="getMemberRole(member.id) === '组长'" class="role-badge leader">
                        <svg viewBox="0 0 24 24" fill="currentColor">
                          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                        </svg>
                        组长
                      </span>
                      <span v-else-if="getMemberRole(member.id) === '管理员'" class="role-badge admin">
                        <svg viewBox="0 0 24 24" fill="currentColor">
                          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                        </svg>
                        管理员
                      </span>
                      <span v-else class="role-badge member">组员</span>
                    </div>
                    <span class="member-username">@{{ member.username }}</span>
                  </div>
                  <div class="member-actions" v-if="isTeamAdmin && member.id !== team.leader_id">
                    <button class="action-btn-sm" :class="{ active: memberAdminMap[member.id] }" @click="toggleAdmin(member, $event)" :title="memberAdminMap[member.id] ? '取消管理员' : '设为管理员'">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                      </svg>
                    </button>
                  </div>
                  <div class="member-join-date">
                    {{ formatDate(member.approved_time || member.create_time) }}
                  </div>
                </div>
                <div v-if="members.length === 0" class="empty-tip">
                  暂无成员
                </div>
              </div>
            </div>
          </div>

          <div v-if="activeTab === 'message'" class="tab-content message-tab">
            <div class="content-card">
              <div class="card-header">
                <h3>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M22 17H2a3 3 0 0 0 3-3V9a7 7 0 0 1 14 0v5a3 3 0 0 0 3 3zm-8.27 4a2 2 0 0 1-3.46 0"/>
                  </svg>
                  消息中心
                </h3>
                <span class="count-badge" v-if="unreadCount > 0">{{ unreadCount }} 未读</span>
                <button v-if="isMember || isTeamAdmin" class="btn btn-primary btn-sm" @click="showMessageDialog">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                  </svg>
                  发消息
                </button>
              </div>
              <div class="messages-list">
                <div v-for="message in messages" :key="message.id" class="message-item" :class="{ unread: !message.is_read }" @click="markAsRead(message)">
                  <div class="message-avatar">
                    <img v-if="message.sender_avatar" :src="message.sender_avatar" :alt="message.sender_name" v-avatar />
                    <span v-else>{{ message.sender_name?.charAt(0) || '?' }}</span>
                  </div>
                  <div class="message-content">
                    <div class="message-header">
                      <span class="message-title">{{ message.title }}</span>
                      <span class="message-time">{{ formatRelativeTime(message.create_time) }}</span>
                    </div>
                    <p class="message-text">{{ message.content }}</p>
                    <span class="message-sender">来自 {{ message.sender_name }}</span>
                  </div>
                  <div class="message-actions" @click.stop>
                    <button class="action-btn-sm danger" @click="deleteMessage(message)" title="删除">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                      </svg>
                    </button>
                  </div>
                </div>
                <div v-if="messages.length === 0" class="empty-tip">
                  暂无消息
                </div>
              </div>
            </div>
          </div>

          <div v-if="activeTab === 'settings'" class="tab-content settings-tab">
            <div class="content-card">
              <div class="card-header">
                <h3>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="3"/>
                    <path d="M12 1v6m0 6v6m-5.5-13.5l4.24 4.24m4.52 4.52L20 21M1 12h6m6 0h6M6.34 17.66l4.24-4.24m4.52-4.52L20 3"/>
                  </svg>
                  管理设置
                </h3>
              </div>
              <div class="settings-content">
                <div class="settings-section">
                  <h4>成员管理</h4>
                  <div class="settings-info">
                    <p>小组成员总数：<strong>{{ members.length }}</strong> 人</p>
                    <p>组长：<strong>{{ leader?.real_name || '未知' }}</strong></p>
                    <p>管理员：<strong>{{ team?.admin_ids?.length || 0 }}</strong> 人</p>
                  </div>
                </div>
                <div class="settings-section" v-if="isTeamAdmin">
                  <h4>快速操作</h4>
                  <div class="settings-actions">
                    <button class="btn btn-secondary" @click="showEditDialog">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                      </svg>
                      编辑小组信息
                    </button>
                    <button class="btn btn-secondary" @click="showAnnouncementDialog">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M22 17H2a3 3 0 0 0 3-3V9a7 7 0 0 1 14 0v5a3 3 0 0 0 3 3zm-8.27 4a2 2 0 0 1-3.46 0"/>
                      </svg>
                      发布公告
                    </button>
                    <button class="btn btn-secondary" @click="openLogsDialog">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                        <polyline points="14 2 14 8 20 8"/>
                        <line x1="16" y1="13" x2="8" y2="13"/>
                        <line x1="16" y1="17" x2="8" y2="17"/>
                      </svg>
                      操作日志
                    </button>
                  </div>
                </div>
                <div class="settings-section">
                  <h4>小组统计</h4>
                  <div class="settings-info">
                    <p>帖子总数：<strong>{{ posts.length }}</strong></p>
                    <p>讨论区话题：<strong>{{ posts.length }}</strong></p>
                    <p>小组成员：<strong>{{ members.length }}</strong> 人</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="activeTab === 'bloggers'" class="tab-content bloggers-tab">
            <div class="content-card">
              <div class="card-header">
                <h3>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                  </svg>
                  全部对接博主
                </h3>
                <span class="count-badge">{{ teamBloggerList.length }}</span>
              </div>
              <div class="blogger-full-grid">
                <div v-for="blogger in teamBloggerList" :key="blogger.id" class="blogger-card" @click="goToBloggerDetail(blogger.id)">
                  <div class="blogger-avatar" :style="{ backgroundColor: getCategoryColor(blogger.category) }">
                    <img v-if="blogger.avatar" :src="blogger.avatar" :alt="blogger.nickname" v-avatar />
                    <span v-else>{{ blogger.nickname?.charAt(0) || '?' }}</span>
                  </div>
                  <div class="blogger-info">
                    <span class="blogger-name">{{ blogger.nickname }}</span>
                    <span class="blogger-category" :style="{ color: getCategoryColor(blogger.category) }">
                      {{ blogger.category || '其他' }}
                    </span>
                  </div>
                </div>
              </div>
              <div v-if="teamBloggerList.length === 0" class="empty-tip">
                暂无对接博主
              </div>
            </div>
          </div>

          <div v-if="activeTab === 'forum'" class="tab-content forum-tab">
            <div class="content-card">
              <div class="card-header forum-header">
                <div class="forum-tabs">
                  <button class="forum-tab-btn" :class="{ active: !collectedTab }" @click="switchToForumTab(false)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
                    </svg>
                    讨论区
                  </button>
                  <button class="forum-tab-btn" :class="{ active: collectedTab }" @click="switchToForumTab(true)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/>
                    </svg>
                    我的收藏
                  </button>
                </div>
                <div class="category-filter" v-if="!collectedTab">
                  <button class="category-btn" :class="{ active: selectedCategory === '全部' }" @click="filterByCategory('全部')">全部</button>
                  <button class="category-btn" :class="{ active: selectedCategory === '综合' }" @click="filterByCategory('综合')">综合</button>
                  <button class="category-btn" :class="{ active: selectedCategory === '讨论' }" @click="filterByCategory('讨论')">讨论</button>
                  <button class="category-btn" :class="{ active: selectedCategory === '分享' }" @click="filterByCategory('分享')">分享</button>
                  <button class="category-btn" :class="{ active: selectedCategory === '求助' }" @click="filterByCategory('求助')">求助</button>
                  <button class="category-btn" :class="{ active: selectedCategory === '公告' }" @click="filterByCategory('公告')">公告</button>
                </div>
                <div class="forum-actions">
                  <div class="search-box" v-if="!collectedTab">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <circle cx="11" cy="11" r="8"/>
                      <line x1="21" y1="21" x2="16.65" y2="16.65"/>
                    </svg>
                    <input type="text" v-model="searchKeyword" placeholder="搜索话题..." @keyup.enter="searchPosts" />
                  </div>
                  <button v-if="isMember || isLeaderOrAdmin" class="btn btn-primary btn-sm" @click="showPostDialog">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <line x1="12" y1="5" x2="12" y2="19"/>
                      <line x1="5" y1="12" x2="19" y2="12"/>
                    </svg>
                    发布话题
                  </button>
                </div>
              </div>
              <div class="posts-list">
                <div v-for="(post, index) in posts" :key="post?.id || index"
                     class="post-item reveal" :style="{ '--delay': index * 0.05 }"
                     @click="post && viewPostDetail(post)">
                  <div class="post-badges" v-if="post && (post.is_pinned || post.is_featured)">
                    <span v-if="post.is_pinned" class="post-badge pinned">
                      <svg viewBox="0 0 24 24" fill="currentColor" width="12" height="12">
                        <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
                      </svg>
                      置顶
                    </span>
                    <span v-if="post.is_featured" class="post-badge featured">
                      <svg viewBox="0 0 24 24" fill="currentColor" width="12" height="12">
                        <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                      </svg>
                      精华
                    </span>
                  </div>
                  <div class="post-header" v-if="post">
                    <div class="post-author-avatar">
                      <img v-if="post.author_avatar" :src="post.author_avatar" :alt="post.author_name || '作者'" v-avatar />
                      <span v-else>{{ post.author_name?.charAt?.(0) || post.author_username?.charAt?.(0) || '?' }}</span>
                    </div>
                    <div class="post-author-info">
                      <span class="post-author-name">{{ post.author_name || post.author_username || '匿名用户' }}</span>
                      <span class="post-category" :style="{ backgroundColor: getCategoryColor(post.category) }">{{ post.category || '综合' }}</span>
                    </div>
                    <span class="post-time">{{ post.create_time ? formatRelativeTime(post.create_time) : '' }}</span>
                  </div>
                  <h4 class="post-title" v-if="post">{{ post.title || '无标题' }}</h4>
                  <p class="post-excerpt" v-if="post">{{ (post.content || '').substring?.(0, 120) }}{{ (post.content?.length || 0) > 120 ? '...' : '' }}</p>
                  <div class="post-footer" v-if="post">
                    <button class="post-action" :class="{ active: postLikedMap[post.id] }" @click.stop="likePost(post, $event)">
                      <svg viewBox="0 0 24 24" :fill="postLikedMap[post.id] ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                        <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
                      </svg>
                      {{ post.like_count ?? 0 }}
                    </button>
                    <span class="post-stat">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
                      </svg>
                      {{ post.comment_count ?? 0 }}
                    </span>
                    <span class="post-stat">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                        <circle cx="12" cy="12" r="3"/>
                      </svg>
                      {{ post.view_count ?? 0 }}
                    </span>
                    <button class="post-action" :class="{ active: postCollectedMap[post.id] }" @click.stop="collectPost(post, $event)">
                      <svg viewBox="0 0 24 24" :fill="postCollectedMap[post.id] ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                        <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/>
                      </svg>
                    </button>
                    <div class="post-admin-actions" v-if="isLeaderOrAdmin" @click.stop>
                      <button class="admin-btn" @click="pinPost(post, $event)" :class="{ active: post.is_pinned }" title="置顶">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
                        </svg>
                      </button>
                      <button class="admin-btn" @click="featurePost(post, $event)" :class="{ active: post.is_featured }" title="加精">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
                <div v-if="posts.length === 0" class="empty-tip">
                  <p>{{ collectedTab ? '暂无收藏' : '暂无话题' }}</p>
                  <p v-if="!collectedTab && (isMember || isLeaderOrAdmin)">点击上方按钮发布第一个话题</p>
                </div>
              </div>
            </div>
          </div>
        </main>

        <aside class="right-sidebar">
          <div class="sidebar-card leader-card reveal" style="--delay: 0">
            <div class="card-header">
              <h3>组长</h3>
            </div>
            <div class="leader-profile" v-if="leader">
              <div class="leader-avatar" :style="{ backgroundColor: team.color }">
                <img v-if="leader.avatar" :src="leader.avatar" :alt="leader.real_name" v-avatar />
                <span v-else>{{ leader.real_name?.charAt(0) || 'L' }}</span>
              </div>
              <div class="leader-info">
                <span class="leader-name">{{ leader.real_name }}</span>
                <span class="leader-username">@{{ leader.username }}</span>
              </div>
            </div>
          </div>

          <div class="sidebar-card reveal" style="--delay: 0.05" v-if="team.announcement">
            <div class="card-header">
              <h3>团队公告</h3>
            </div>
            <div class="announcement-sidebar">
              <div class="announcement-content-preview">
                <p>{{ team.announcement }}</p>
              </div>
              <div class="announcement-meta">
                <span class="announcement-from">{{ team.leader_name || '管理员' }}</span>
                <span class="announcement-time">{{ formatAnnouncementTime(team.update_time) }}</span>
              </div>
            </div>
          </div>

          <div class="sidebar-card reveal" style="--delay: 0.1" v-if="isLeaderOrAdmin">
            <div class="card-header">
              <h3>快捷操作</h3>
            </div>
            <div class="quick-actions">
              <button class="action-btn" @click="showEditDialog">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
                编辑资料
              </button>
              <button class="action-btn" @click="showAnnouncementDialog">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 17H2a3 3 0 0 0 3-3V9a7 7 0 0 1 14 0v5a3 3 0 0 0 3 3zm-8.27 4a2 2 0 0 1-3.46 0"/>
                </svg>
                发布公告
              </button>
            </div>
          </div>

          <div class="sidebar-card reveal" style="--delay: 0.2">
            <div class="card-header">
              <h3>热门成员</h3>
            </div>
            <div class="top-members">
              <div v-for="(member, index) in members.slice(0, 5)" :key="member.id" class="top-member-item">
                <span class="rank" :class="{ gold: index === 0, silver: index === 1, bronze: index === 2 }">
                  {{ index + 1 }}
                </span>
                <div class="member-avatar-sm" :style="{ backgroundColor: team.color }">
                  <span>{{ member.real_name?.charAt(0) || 'U' }}</span>
                </div>
                <span class="member-name-sm">{{ member.real_name }}</span>
                <span v-if="member.id === team.leader_id" class="is-leader-tag">组长</span>
              </div>
            </div>
          </div>
        </aside>
      </div>
    </div>

    <div v-if="editDialogVisible" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="closeEditDialog">
      <div class="modal team-edit-modal" style="background:var(--bg-card);">
        <div class="modal-header">
          <h3>编辑小组信息</h3>
          <button class="close-btn" @click="closeEditDialog">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-field">
            <label>小组名称</label>
            <input v-model="editForm.name" type="text" placeholder="请输入小组名称" />
          </div>
          <div class="form-field">
            <label>小组描述</label>
            <textarea v-model="editForm.description" placeholder="请输入小组描述" rows="3"></textarea>
          </div>
          <div class="form-field">
            <label>小组图标</label>
            <div class="image-upload-row">
              <div class="preview-avatar" :style="{ backgroundColor: editForm.color }">
                <img v-if="editForm.logo" :src="editForm.logo" alt="logo" />
                <span v-else>{{ editForm.name?.charAt(0) || 'L' }}</span>
              </div>
              <button class="btn btn-secondary" @click="triggerLogoUpload">上传图标</button>
              <input ref="logoInput" type="file" accept="image/*" hidden @change="handleLogoUpload" />
            </div>
          </div>
          <div class="form-field">
            <label>背景图片</label>
            <div class="image-upload-row">
              <div class="preview-bg" :style="editForm.bg_image ? { backgroundImage: `url(${editForm.bg_image})` } : {}">
                <span v-if="!editForm.bg_image">背景图</span>
              </div>
              <button class="btn btn-secondary" @click="triggerBgUpload">上传背景</button>
              <input ref="bgInput" type="file" accept="image/*" hidden @change="handleBgUpload" />
            </div>
          </div>
          <div class="form-field">
            <label>主题颜色</label>
            <div class="color-picker-row">
              <input type="color" v-model="editForm.color" />
              <input type="text" v-model="editForm.color" placeholder="#6366f1" />
            </div>
          </div>
          <div class="form-field">
            <label>组长</label>
            <select v-model="editForm.leader_id">
              <option value="">请选择组长</option>
              <option v-for="member in members" :key="member.id" :value="member.id">
                {{ member.real_name }} (@{{ member.username }})
              </option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeEditDialog">取消</button>
          <button class="btn btn-primary" @click="saveTeamEdit" :disabled="saving">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="announcementDialogVisible" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="closeAnnouncementDialog">
      <div class="modal announcement-modal" style="background:var(--bg-card);">
        <div class="modal-header">
          <h3>发布公告</h3>
          <button class="close-btn" @click="closeAnnouncementDialog">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-field">
            <label>公告内容</label>
            <textarea v-model="announcementText" placeholder="请输入公告内容" rows="4"></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="clearAnnouncement">清空</button>
          <button class="btn btn-primary" @click="saveAnnouncement" :disabled="saving">
            {{ saving ? '发布中...' : '发布' }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="announcementPopupVisible && team.announcement" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="dismissAnnouncement">
      <div class="modal announcement-popup-modal" style="background:var(--bg-card);">
        <div class="modal-header">
          <div class="announcement-popup-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 17H2a3 3 0 0 0 3-3V9a7 7 0 0 1 14 0v5a3 3 0 0 0 3 3zm-8.27 4a2 2 0 0 1-3.46 0"/>
            </svg>
          </div>
          <h3>团队公告</h3>
          <button class="close-btn" @click="dismissAnnouncement">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="announcement-popup-content">
            <div class="announcement-meta">
              <span class="announcement-from">{{ team.leader_name || '管理员' }}</span>
              <span class="announcement-time">{{ formatAnnouncementTime(team.update_time) }}</span>
            </div>
            <div class="announcement-text">{{ team.announcement }}</div>
          </div>
        </div>
        <div class="modal-footer">
          <label class="dont-show-again">
            <input type="checkbox" v-model="dontShowAgain" />
            <span>不再弹出此公告</span>
          </label>
          <button class="btn btn-primary" @click="dismissAnnouncement">我知道了</button>
        </div>
      </div>
    </div>

    <div v-if="messageDialogVisible" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="closeMessageDialog">
      <div class="modal message-modal" style="background:var(--bg-card);">
        <div class="modal-header">
          <h3>发送消息</h3>
          <button class="close-btn" @click="closeMessageDialog">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-field">
            <label>消息标题</label>
            <input type="text" v-model="messageForm.title" placeholder="请输入消息标题" maxlength="50" />
          </div>
          <div class="form-field">
            <label>消息内容</label>
            <textarea v-model="messageForm.content" placeholder="请输入消息内容" rows="4"></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeMessageDialog">取消</button>
          <button class="btn btn-primary" @click="sendMessage">发送</button>
        </div>
      </div>
    </div>

    <div v-if="batchMessageDialogVisible" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="batchMessageDialogVisible = false">
      <div class="modal message-modal" style="background:var(--bg-card);">
        <div class="modal-header">
          <h3>批量发送消息 ({{ selectedMembers.length }}人)</h3>
          <button class="close-btn" @click="batchMessageDialogVisible = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-field">
            <label>消息标题</label>
            <input type="text" v-model="messageForm.title" placeholder="请输入消息标题" maxlength="50" />
          </div>
          <div class="form-field">
            <label>消息内容</label>
            <textarea v-model="messageForm.content" placeholder="请输入消息内容" rows="4"></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="batchMessageDialogVisible = false">取消</button>
          <button class="btn btn-primary" @click="sendBatchMessage">发送</button>
        </div>
      </div>
    </div>

    <div v-if="logsDialogVisible" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="logsDialogVisible = false">
      <div class="modal logs-modal" style="background:var(--bg-card);">
        <div class="modal-header">
          <h3>操作日志</h3>
          <button class="close-btn" @click="logsDialogVisible = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body logs-body">
          <div v-if="logsLoading" class="loading">加载中...</div>
          <div v-else-if="teamLogs.length === 0" class="empty-tip">暂无日志记录</div>
          <div v-else class="log-list">
            <div v-for="log in teamLogs" :key="log.id" class="log-item">
              <div class="log-icon">
                <span>📋</span>
              </div>
              <div class="log-content">
                <div class="log-action">{{ log.action }}</div>
                <div class="log-detail" v-if="log.detail">{{ log.detail }}</div>
                <div class="log-meta">
                  <span>{{ log.operator }}</span>
                  <span>{{ formatDateTime(log.create_time) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer" v-if="logsTotal > 50">
          <button class="btn btn-secondary btn-sm" :disabled="logsPage <= 1" @click="loadTeamLogs(logsPage - 1)">上一页</button>
          <span class="page-info">{{ logsPage }} / {{ Math.ceil(logsTotal / 50) }}</span>
          <button class="btn btn-secondary btn-sm" :disabled="logsPage >= Math.ceil(logsTotal / 50)" @click="loadTeamLogs(logsPage + 1)">下一页</button>
        </div>
      </div>
    </div>

    <div v-if="postDialogVisible" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="closePostDialog">
      <div class="modal post-modal" style="background:var(--bg-card);">
        <div class="modal-header">
          <h3>发布话题</h3>
          <button class="close-btn" @click="closePostDialog">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-field">
            <label>标题</label>
            <input type="text" v-model="newPostForm.title" placeholder="请输入话题标题" maxlength="50" />
          </div>
          <div class="form-field">
            <label>分类</label>
            <select v-model="newPostForm.category">
              <option value="综合">综合</option>
              <option value="讨论">讨论</option>
              <option value="分享">分享</option>
              <option value="求助">求助</option>
              <option value="公告">公告</option>
            </select>
          </div>
          <div class="form-field">
            <label>内容</label>
            <textarea v-model="newPostForm.content" placeholder="请输入话题内容" rows="6"></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closePostDialog">取消</button>
          <button class="btn btn-primary" @click="submitPost" :disabled="isPosting">
            {{ isPosting ? '发布中...' : '发布' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 帖子详情弹窗 -->
    <Teleport to="body">
      <div v-if="postDetailVisible" class="post-detail-modal-wrapper" @click.self="closePostDetail">
        <div class="post-detail-modal-container" v-if="postDetail">
          <div class="post-detail-modal-header">
            <h3>{{ postDetail.title }}</h3>
            <div class="post-detail-modal-actions">
              <button class="modal-action-btn" @click="sharePost" title="分享">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="18" cy="5" r="3"/>
                  <circle cx="6" cy="12" r="3"/>
                  <circle cx="18" cy="19" r="3"/>
                  <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/>
                  <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/>
                </svg>
              </button>
              <button class="modal-close-btn" @click="closePostDetail">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"/>
                  <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
              </button>
            </div>
          </div>
          <div class="post-detail-modal-body">
            <div class="post-meta">
              <div class="post-author">
                <div class="author-avatar">
                  <img v-if="postDetail.author_avatar" :src="postDetail.author_avatar" :alt="postDetail.author_name" v-avatar />
                  <span v-else>{{ postDetail.author_name?.charAt(0) || '?' }}</span>
                </div>
                <div class="author-info">
                  <span class="author-name">{{ postDetail.author_name }}</span>
                  <span class="post-category" :style="{ backgroundColor: getCategoryColor(postDetail.category) }">{{ postDetail.category }}</span>
                </div>
              </div>
              <span class="post-time">{{ formatDate(postDetail.create_time) }}</span>
            </div>
            <div class="post-content">
              {{ postDetail.content }}
            </div>
            <div class="post-stats">
              <span class="stat-item">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
                </svg>
                {{ postDetail.comment_count || 0 }} 评论
              </span>
              <span class="stat-item">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
                {{ postDetail.view_count || 0 }} 浏览
              </span>
            </div>
            <div class="comments-section">
              <h4 class="comments-title">评论</h4>
              <div class="comments-list">
                <div v-for="comment in postDetail.comments" :key="comment.id" class="comment-item">
                  <div class="comment-avatar">
                    <img v-if="comment.author_avatar" :src="comment.author_avatar" :alt="comment.author_name" v-avatar />
                    <span v-else>{{ comment.author_name?.charAt(0) || '?' }}</span>
                  </div>
                  <div class="comment-body">
                    <div class="comment-header">
                      <span class="comment-author">{{ comment.author_name }}</span>
                      <span class="comment-time">{{ formatRelativeTime(comment.create_time) }}</span>
                    </div>
                    <p class="comment-text">{{ comment.content }}</p>
                  </div>
                </div>
                <div v-if="!postDetail.comments || postDetail.comments.length === 0" class="no-comments">
                  暂无评论
                </div>
              </div>
              <div class="comment-input-area" v-if="isMember || isLeaderOrAdmin">
                <div class="mention-input-wrapper">
                  <input type="text" ref="commentInputRef" v-model="postComment" placeholder="写下你的评论...(@提及)" @input="handleCommentInput" @keydown="handleCommentKeydown" @keyup.enter="submitComment" />
                  <div class="mention-dropdown" v-if="mentionSearchVisible" :style="{ bottom: '100%', top: 'auto' }">
                    <div v-for="user in mentionSuggestions" :key="user.id" class="mention-item" @click="selectMention(user)">
                      <div class="mention-avatar">
                        <img v-if="user.avatar" :src="user.avatar" :alt="user.real_name" v-avatar />
                        <span v-else>{{ user.real_name?.charAt(0) }}</span>
                      </div>
                      <div class="mention-info">
                        <span class="mention-name">{{ user.real_name }}</span>
                        <span class="mention-username">@{{ user.username }}</span>
                      </div>
                    </div>
                    <div v-if="mentionSuggestions.length === 0" class="mention-empty">无匹配用户</div>
                  </div>
                </div>
                <button class="submit-comment-btn btn btn-primary btn-sm" @click="submitComment">发送</button>
              </div>
            </div>
          </div>
          <div class="post-detail-modal-footer" v-if="postDetail.author_id === userStore.id || isLeaderOrAdmin">
            <button class="btn btn-danger" @click="deletePost(postDetail)">删除帖子</button>
          </div>
        </div>
      </div>
    </Teleport>

    <nav class="mobile-bottom-nav" v-if="isMobile">
      <a class="mobile-nav-item" :class="{ active: activeTab === 'home' }" @click="activeTab = 'home'">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
          <polyline points="9 22 9 12 15 12 15 22"/>
        </svg>
        <span>首页</span>
      </a>
      <a class="mobile-nav-item" :class="{ active: activeTab === 'forum' }" @click="activeTab = 'forum'">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
        </svg>
        <span>讨论区</span>
      </a>
      <a class="mobile-nav-item" :class="{ active: activeTab === 'members' }" @click="activeTab = 'members'">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
          <circle cx="9" cy="7" r="4"/>
          <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
          <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
        </svg>
        <span>成员</span>
      </a>
      <a class="mobile-nav-item" :class="{ active: activeTab === 'message' }" @click="switchToMessage">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M22 17H2a3 3 0 0 0 3-3V9a7 7 0 0 1 14 0v5a3 3 0 0 0 3 3zm-8.27 4a2 2 0 0 1-3.46 0"/>
        </svg>
        <span class="mobile-badge" v-if="unreadCount > 0">{{ unreadCount > 99 ? '99+' : unreadCount }}</span>
        <span>消息</span>
      </a>
      <a class="mobile-nav-item" :class="{ active: activeTab === 'settings' }" @click="activeTab = 'settings'" v-if="isTeamAdmin">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/>
          <path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20z"/>
        </svg>
        <span>管理</span>
      </a>
    </nav>

    <Teleport to="body">
      <div v-if="allUsersDialogVisible" class="modal-overlay" @click.self="closeAllUsersDialog">
        <div class="all-users-modal">
          <div class="modal-header">
            <h3>全部用户</h3>
            <button class="close-btn" @click="closeAllUsersDialog">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-search">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
            <input type="text" v-model="userSearchKeyword" placeholder="搜索用户..." />
          </div>
          <div class="users-list">
            <div v-for="user in filteredAllUsers" :key="user.id" class="user-item">
              <div class="user-avatar" :style="{ backgroundColor: getUserColor(user) }">
                <img v-if="user.avatar" :src="user.avatar" :alt="user.real_name" v-avatar />
                <span v-else>{{ user.real_name?.charAt(0) || 'U' }}</span>
              </div>
              <div class="user-info">
                <span class="user-name">{{ user.real_name }}</span>
                <span class="user-team" v-if="user.team_name">{{ user.team_name }}</span>
                <span class="user-role" :class="user.role">{{ user.role === 'admin' ? '管理员' : '成员' }}</span>
              </div>
              <div class="user-actions">
                <button class="btn btn-sm btn-secondary" @click="viewUserProfile(user)">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                  </svg>
                  资料
                </button>
                <button class="btn btn-sm btn-primary" @click="sendPrivateMessage(user)" v-if="user.id !== userStore.userId">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
                  </svg>
                  私信
                </button>
              </div>
            </div>
            <div v-if="filteredAllUsers.length === 0" class="empty-users">
              未找到匹配的用户
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="privateMessageDialogVisible" class="modal-overlay" @click.self="closePrivateMessageDialog">
        <div class="private-message-modal">
          <div class="modal-header">
            <h3>发送私信给 {{ selectedUser?.real_name }}</h3>
            <button class="close-btn" @click="closePrivateMessageDialog">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <textarea v-model="privateMessageContent" placeholder="输入私信内容..." rows="5"></textarea>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="closePrivateMessageDialog">取消</button>
            <button class="btn btn-primary" @click="submitPrivateMessage" :disabled="!privateMessageContent.trim()">发送</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { logger } from '../utils/logger'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'
import { useConfirm } from '../utils/confirm'
import { getTeamsAPI, updateTeamAPI, setUserTeamAPI, getPublicUsersAPI, updateMyTeamAPI, getTeamBloggerStatAPI, getTeamBloggerChartsAPI, uploadTeamLogoAPI, uploadTeamBgAPI, getTeamPostsAPI, createTeamPostAPI, getTeamPostDetailAPI, deleteTeamPostAPI, createTeamPostCommentAPI, likeTeamPostAPI, getTeamPostLikeStatusAPI, pinTeamPostAPI, featureTeamPostAPI, searchTeamPostsAPI, collectTeamPostAPI, getTeamPostCollectStatusAPI, getTeamPostsCollectedAPI, setTeamAdminAPI, getTeamMessagesAPI, sendTeamMessageAPI, markMessageReadAPI, deleteTeamMessageAPI, getTeamOperationLogsAPI, removeTeamMemberAPI, sendPrivateMessageAPI } from '../api'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const { success, error: notifyError, warning } = useNotification()
const { confirmDanger } = useConfirm()

const loading = ref(true)
const team = ref(null)
const members = ref([])
const teamStats = ref({})
const teamBloggerList = ref([])
const sidebarCollapsed = ref(false)
const activeTab = ref('home')
const isMobile = ref(window.innerWidth <= 768)

const posts = ref([])
const postDialogVisible = ref(false)
const postDetailVisible = ref(false)
const postDetail = ref(null)
const newPostForm = ref({ title: '', content: '', category: '综合' })
const postComment = ref('')
const mentionSuggestions = ref([])
const mentionSearchVisible = ref(false)
const mentionSearchQuery = ref('')
const isPosting = ref(false)
const postLikedMap = ref({})
const postCollectedMap = ref({})
const searchKeyword = ref('')
const collectedTab = ref(false)
const selectedCategory = ref('全部')
const postsPage = ref(1)
const postsTotal = ref(0)
const messages = ref([])
const messageDialogVisible = ref(false)
const batchMessageDialogVisible = ref(false)
const messageForm = ref({ title: '', content: '', type: 'normal' })
const unreadCount = ref(0)
const memberAdminMap = ref({})
const selectedMembers = ref([])
const isSelectingMembers = ref(false)
const teamLogs = ref([])

const allUsersDialogVisible = ref(false)
const allUsersList = ref([])
const userSearchKeyword = ref('')
const privateMessageDialogVisible = ref(false)
const selectedUser = ref(null)
const privateMessageContent = ref('')

const filteredAllUsers = computed(() => {
  if (!userSearchKeyword.value) return allUsersList.value
  const keyword = userSearchKeyword.value.toLowerCase()
  return allUsersList.value.filter(u => 
    u.real_name?.toLowerCase().includes(keyword) || 
    u.username?.toLowerCase().includes(keyword)
  )
})
const logsDialogVisible = ref(false)
const logsLoading = ref(false)
const logsPage = ref(1)
const logsTotal = ref(0)

const editDialogVisible = ref(false)
const announcementDialogVisible = ref(false)
const announcementPopupVisible = ref(false)
const editForm = ref({})
const announcementText = ref('')
const saving = ref(false)
const logoInput = ref(null)
const bgInput = ref(null)
const dismissedAnnouncements = ref(new Set())
const dontShowAgain = ref(false)

const teamId = computed(() => route.params.teamId)

const isMember = computed(() => {
  return userStore.team_id && String(userStore.team_id) === String(teamId.value)
})

const isLeaderOrAdmin = computed(() => {
  if (!team.value) return false
  if (userStore.role === 'admin') return true
  if (userStore.id === team.value.leader_id) return true
  if (team.value.admin_ids?.includes(userStore.id)) return true
  return false
})

const isTeamAdmin = computed(() => {
  if (!team.value) return false
  if (userStore.role === 'admin') return true
  if (userStore.id === team.value.leader_id) return true
  if (team.value.admin_ids?.includes(userStore.id)) return true
  return false
})

const isPrivate = computed(() => {
  return team.value?.is_public === false
})

const goToBloggerDetail = (bloggerId) => {
  router.push(`/blogger/${bloggerId}`)
}

const isTeamMember = computed(() => {
  if (!team.value) return false
  return userStore.team_id === team.value.id
})

const leader = computed(() => {
  return members.value.find(m => m.id === team.value?.leader_id)
})

const bannerStyle = computed(() => {
  if (team.value?.bg_image) {
    return { backgroundImage: `url(${team.value.bg_image})` }
  }
  return { background: team.value?.color || 'var(--purple)' }
})

const categoryColors = {
  '健身': 'var(--success)',
  '美妆': 'var(--purple)',
  '美食': 'var(--warning)',
  '旅游': 'var(--info)',
  '数码': 'var(--info)',
  '穿搭': 'var(--purple)',
  '母婴': 'var(--primary)',
  '家居': 'var(--success)',
  '职场': 'var(--purple)',
  '情感': 'var(--danger)',
  '教育': 'var(--success)',
  '游戏': 'var(--purple)',
  '萌宠': 'var(--primary)',
  '其他': '#64748b'
}

const getCategoryColor = (category) => {
  return categoryColors[category] || categoryColors['其他']
}

const getUserColor = (user) => {
  const colors = ['var(--primary)', 'var(--purple)', 'var(--success)', 'var(--warning)', 'var(--danger)', 'var(--purple)', 'var(--purple)', 'var(--success)']
  const index = user.id % colors.length
  return colors[index]
}

const showAllUsersDialog = async () => {
  allUsersDialogVisible.value = true
  try {
    const data = await getPublicUsersAPI()
    if (data.code === 200) {
      allUsersList.value = data.data.filter(u => u.status === 'active') || []
    }
  } catch (error) {
    notifyError('加载用户列表失败')
  }
}

const closeAllUsersDialog = () => {
  allUsersDialogVisible.value = false
  userSearchKeyword.value = ''
}

const sendPrivateMessage = (user) => {
  selectedUser.value = user
  privateMessageDialogVisible.value = true
  privateMessageContent.value = ''
}

const closePrivateMessageDialog = () => {
  privateMessageDialogVisible.value = false
  selectedUser.value = null
  privateMessageContent.value = ''
}

const submitPrivateMessage = async () => {
  if (!privateMessageContent.value.trim() || !selectedUser.value) return
  
  try {
    const data = await sendPrivateMessageAPI({
      to_user_id: selectedUser.value.id,
      content: privateMessageContent.value
    })
    if (data.code === 200) {
      success('私信发送成功')
      closePrivateMessageDialog()
    } else {
      notifyError(data.message || '发送失败')
    }
  } catch (error) {
    notifyError('发送失败')
  }
}

const viewUserProfile = (user) => {
  closeAllUsersDialog()
  router.push(`/user/${user.username}`)
}

const formatDate = (date) => {
  if (!date) return '-'
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const loadTeamData = async () => {
  loading.value = true
  try {
    const [teamsRes, usersRes, postsRes, messagesRes] = await Promise.all([
      getTeamsAPI(),
      getPublicUsersAPI(),
      getTeamPostsAPI(teamId.value),
      getTeamMessagesAPI(teamId.value)
    ])

    if (teamsRes.code === 200) {
      team.value = teamsRes.data.find(t => String(t.id) === String(teamId.value))
    }

    if (usersRes.code === 200) {
      const users = usersRes.data || []
      members.value = users.filter(u => String(u.team_id) === String(teamId.value) && u.status === 'active')
    }

    if (postsRes.code === 200) {
      posts.value = postsRes.data.list || postsRes.data || []
    }

    if (messagesRes.code === 200) {
      messages.value = messagesRes.data || []
      unreadCount.value = messages.value.filter(m => !m.is_read).length
    }

    if (team.value) {
      team.value.admin_ids = team.value.admin_ids || []
      members.value.forEach(m => {
        memberAdminMap.value[m.id] = team.value.admin_ids?.includes(m.id) || false
      })
      loadTeamBloggerStats()
      nextTick(() => {
        checkAnnouncementPopup()
      })
    } else {
      // 团队不存在
      notifyError('团队不存在或已被删除')
    }
  } catch (err) {
    logger.error('', err)
    notifyError('加载小组数据失败')
    const errorMsg = err.response?.status === 401 
      ? '请先登录' 
      : '加载数据失败，请检查网络连接'
    notifyError(errorMsg)
  } finally {
    loading.value = false
  }
}

const loadTeamBloggerStats = async () => {
  try {
    const [statRes, chartRes] = await Promise.all([
      getTeamBloggerStatAPI(teamId.value),
      getTeamBloggerChartsAPI(teamId.value)
    ])

    if (statRes.code === 200) {
      teamStats.value = statRes.data.stats || {}
    }

    if (chartRes.code === 200) {
      teamBloggerList.value = chartRes.data.bloggerList || chartRes.data.byMember?.flatMap(m => m.bloggers || []) || []
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载博主统计数据失败')
  }
}

const loadPosts = async () => {
  try {
    const category = selectedCategory.value === '全部' ? '' : selectedCategory.value
    const res = await getTeamPostsAPI(teamId.value, postsPage.value, 20, category)
    if (res.code === 200) {
      const rawPosts = res.data.list || res.data || []
      posts.value = rawPosts.filter(p => p && p.id != null)
      postsTotal.value = res.data.total || res.total || 0
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载帖子失败')
  }
}

const filterByCategory = (category) => {
  selectedCategory.value = category
  postsPage.value = 1
  loadPosts()
}

const showPostDialog = () => {
  if (!isMember.value && !isLeaderOrAdmin.value) {
    warning('只有团队成员才能发帖')
    return
  }
  newPostForm.value = { title: '', content: '', category: '综合' }
  postDialogVisible.value = true
}

const closePostDialog = () => {
  postDialogVisible.value = false
}

const submitPost = async () => {
  if (!newPostForm.value.title) {
    warning('请输入帖子标题')
    return
  }
  if (!newPostForm.value.content) {
    warning('请输入帖子内容')
    return
  }
  isPosting.value = true
  try {
    const res = await createTeamPostAPI(teamId.value, newPostForm.value)
    if (res.code === 200) {
      success('发帖成功')
      posts.value.unshift(res.data)
      closePostDialog()
    } else {
      notifyError(res.message || '发帖失败')
    }
  } catch (error) {
    notifyError('发帖失败')
  } finally {
    isPosting.value = false
  }
}

const viewPostDetail = async (post) => {
  try {
    const res = await getTeamPostDetailAPI(teamId.value, post.id)
    if (res.code === 200) {
      postDetail.value = res.data
      postDetailVisible.value = true
      if (isMember.value || isLeaderOrAdmin.value) {
        const [likeRes, collectRes] = await Promise.all([
          getTeamPostLikeStatusAPI(teamId.value, post.id),
          getTeamPostCollectStatusAPI(teamId.value, post.id)
        ])
        postLikedMap.value[post.id] = likeRes.data?.liked || false
        postCollectedMap.value[post.id] = collectRes.data?.collected || false
      }
    } else {
      notifyError(res.message || '加载帖子详情失败')
    }
  } catch (error) {
    notifyError('加载帖子详情失败')
  }
}

const closePostDetail = () => {
  postDetailVisible.value = false
  postDetail.value = null
  postComment.value = ''
}

const sharePost = async () => {
  if (!postDetail.value) return
  const shareUrl = `${window.location.origin}/team/${teamId.value}/post/${postDetail.value.id}`
  const shareText = `${postDetail.value.title} - 来自 ${team.value?.name || '小组'}论坛`
  try {
    if (navigator.clipboard) {
      await navigator.clipboard.writeText(shareUrl)
      success('链接已复制到剪贴板')
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = shareUrl
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
      success('链接已复制到剪贴板')
    }
  } catch (error) {
    notifyError('复制失败')
  }
}

const handleCommentInput = () => {
  const text = postComment.value
  const lastAtIndex = text.lastIndexOf('@')
  if (lastAtIndex !== -1) {
    const searchText = text.slice(lastAtIndex + 1)
    if (!searchText.includes(' ')) {
      mentionSearchQuery.value = searchText
      mentionSuggestions.value = members.value.filter(m =>
        m.username.toLowerCase().includes(searchText.toLowerCase()) ||
        m.real_name.toLowerCase().includes(searchText.toLowerCase())
      ).slice(0, 5)
      mentionSearchVisible.value = true
      return
    }
  }
  mentionSearchVisible.value = false
}

const handleCommentKeydown = (e) => {
  if (mentionSearchVisible.value) {
    if (e.key === 'Escape') {
      mentionSearchVisible.value = false
      e.preventDefault()
    }
  }
}

const selectMention = (user) => {
  const text = postComment.value
  const lastAtIndex = text.lastIndexOf('@')
  postComment.value = text.slice(0, lastAtIndex) + '@' + user.username + ' '
  mentionSearchVisible.value = false
}

const submitComment = async () => {
  if (mentionSearchVisible.value) {
    mentionSearchVisible.value = false
    return
  }
  if (!postComment.value) {
    warning('请输入评论内容')
    return
  }
  try {
    const res = await createTeamPostCommentAPI(teamId.value, postDetail.value.id, { content: postComment.value })
    if (res.code === 200) {
      success('评论成功')
      if (!postDetail.value.comments) {
        postDetail.value.comments = []
      }
      postDetail.value.comments.push(res.data)
      postDetail.value.comment_count = (postDetail.value.comment_count || 0) + 1
      postComment.value = ''
      const postIndex = posts.value.findIndex(p => p.id === postDetail.value.id)
      if (postIndex !== -1) {
        posts.value[postIndex].comment_count = postDetail.value.comment_count
      }
    } else {
      notifyError(res.message || '评论失败')
    }
  } catch (error) {
    notifyError('评论失败')
  }
}

const deletePost = async (post) => {
  if (!await confirmDanger('确定要删除这篇帖子吗？')) return
  try {
    const res = await deleteTeamPostAPI(teamId.value, post.id)
    if (res.code === 200) {
      success('删除成功')
      posts.value = posts.value.filter(p => p.id !== post.id)
      closePostDetail()
    } else {
      notifyError(res.message || '删除失败')
    }
  } catch (error) {
    notifyError('删除失败')
  }
}

const likePost = async (post, event) => {
  event?.stopPropagation()
  if (!isMember.value && !isLeaderOrAdmin.value) {
    warning('只有团队成员才能点赞')
    return
  }
  try {
    const res = await likeTeamPostAPI(teamId.value, post.id)
    if (res.code === 200) {
      post.like_count = res.like_count
      postLikedMap.value[post.id] = res.liked
      if (res.liked) {
        success('点赞成功')
      }
    } else {
      notifyError(res.message || '操作失败')
    }
  } catch (error) {
    notifyError('操作失败')
  }
}

const collectPost = async (post, event) => {
  event?.stopPropagation()
  if (!isMember.value && !isLeaderOrAdmin.value) {
    warning('只有团队成员才能收藏')
    return
  }
  try {
    const res = await collectTeamPostAPI(teamId.value, post.id)
    if (res.code === 200) {
      postCollectedMap.value[post.id] = res.collected
      if (res.collected) {
        success('收藏成功')
      } else {
        success('已取消收藏')
      }
    } else {
      notifyError(res.message || '操作失败')
    }
  } catch (error) {
    notifyError('操作失败')
  }
}

const pinPost = async (post, event) => {
  event?.stopPropagation()
  try {
    const res = await pinTeamPostAPI(teamId.value, post.id)
    if (res.code === 200) {
      post.is_pinned = res.is_pinned
      success(res.message)
    } else {
      notifyError(res.message || '操作失败')
    }
  } catch (error) {
    notifyError('操作失败')
  }
}

const featurePost = async (post, event) => {
  event?.stopPropagation()
  try {
    const res = await featureTeamPostAPI(teamId.value, post.id)
    if (res.code === 200) {
      post.is_featured = res.is_featured
      success(res.message)
    } else {
      notifyError(res.message || '操作失败')
    }
  } catch (error) {
    notifyError('操作失败')
  }
}

const searchPosts = async () => {
  if (!searchKeyword.value.trim()) {
    loadPosts()
    return
  }
  try {
    const res = await searchTeamPostsAPI(teamId.value, searchKeyword.value)
    if (res.code === 200) {
      posts.value = res.data || []
    }
  } catch (error) {
    notifyError('搜索失败')
  }
}

const loadCollectedPosts = async () => {
  try {
    const res = await getTeamPostsCollectedAPI(teamId.value)
    if (res.code === 200) {
      posts.value = res.data || []
    }
  } catch (error) {
    notifyError('加载失败')
  }
}

const switchToForumTab = (isCollected) => {
  collectedTab.value = isCollected
  if (isCollected) {
    loadCollectedPosts()
  } else {
    loadPosts()
  }
}

const switchToMessage = () => {
  activeTab.value = 'message'
  loadMessages()
}

const loadTeamLogs = async (page = 1) => {
  logsLoading.value = true
  try {
    const res = await getTeamOperationLogsAPI(teamId.value, page, 50)
    if (res.code === 200) {
      teamLogs.value = res.data.logs
      logsTotal.value = res.data.total
      logsPage.value = res.data.page
    }
  } catch (error) {
    logger.error('', error)
    notifyError('加载团队日志失败')
  } finally {
    logsLoading.value = false
  }
}

const openLogsDialog = () => {
  logsDialogVisible.value = true
  loadTeamLogs()
}

const getMemberRole = (memberId) => {
  if (team.value?.leader_id === memberId) return '组长'
  if (team.value?.admin_ids?.includes(memberId)) return '管理员'
  return '组员'
}

const toggleAdmin = async (member, event) => {
  event?.stopPropagation()
  if (!isTeamAdmin.value) {
    warning('只有组长或管理员才能设置管理员')
    return
  }
  if (team.value?.leader_id === member.id) {
    warning('不能对组长进行此操作')
    return
  }
  try {
    const res = await setTeamAdminAPI(teamId.value, member.id)
    if (res.code === 200) {
      if (!team.value.admin_ids) team.value.admin_ids = []
      if (res.is_admin) {
        team.value.admin_ids.push(member.id)
        memberAdminMap.value[member.id] = true
        success(`${member.real_name} 已设为管理员`)
      } else {
        team.value.admin_ids = team.value.admin_ids.filter(id => id !== member.id)
        memberAdminMap.value[member.id] = false
        success(`已取消 ${member.real_name} 的管理员权限`)
      }
    } else {
      notifyError(res.message || '操作失败')
    }
  } catch (error) {
    notifyError('操作失败')
  }
}

const toggleMemberSelection = (member) => {
  if (!isTeamAdmin.value || member.id === team.value?.leader_id) return
  const index = selectedMembers.value.indexOf(member.id)
  if (index > -1) {
    selectedMembers.value.splice(index, 1)
  } else {
    selectedMembers.value.push(member.id)
  }
}

const handleMembersAreaClick = () => {
  if (selectedMembers.value.length === 0 && !isSelectingMembers.value) {
    isSelectingMembers.value = true
  }
}

const clearSelection = () => {
  selectedMembers.value = []
  isSelectingMembers.value = false
}

const startMemberSelection = () => {
  isSelectingMembers.value = true
}

const openBatchMessageDialog = () => {
  if (selectedMembers.value.length === 0) return
  batchMessageDialogVisible.value = true
}

const batchRemoveMembers = async () => {
  if (selectedMembers.value.length === 0) return
  if (!await confirmDanger(`确定要移出选中的 ${selectedMembers.value.length} 位组员吗？`)) return
  try {
    for (const userId of selectedMembers.value) {
      await removeTeamMemberAPI(teamId.value, userId)
    }
    members.value = members.value.filter(m => !selectedMembers.value.includes(m.id))
    success(`已成功移出 ${selectedMembers.value.length} 位组员`)
    clearSelection()
  } catch (error) {
    notifyError('批量移出失败')
  }
}

const batchSetAdmin = async () => {
  if (selectedMembers.value.length === 0) return
  try {
    const res = await batchSetAdminAPI(teamId.value, selectedMembers.value, 'add')
    if (res.code === 200) {
      res.results?.forEach(r => {
        if (r.status === 'added' && team.value && !team.value.admin_ids.includes(r.user_id)) {
          team.value.admin_ids.push(r.user_id)
        }
        memberAdminMap.value[r.user_id] = true
      })
      success(`已成功设置 ${res.results?.length || 0} 位管理员`)
      clearSelection()
    } else {
      notifyError(res.message || '批量设置失败')
    }
  } catch (error) {
    notifyError('批量设置失败')
  }
}

const batchRemoveAdmin = async () => {
  if (selectedMembers.value.length === 0) return
  try {
    const res = await batchSetAdminAPI(teamId.value, selectedMembers.value, 'remove')
    if (res.code === 200) {
      res.results?.forEach(r => {
        if (r.status === 'removed' && team.value) {
          team.value.admin_ids = team.value.admin_ids.filter(id => id !== r.user_id)
        }
        memberAdminMap.value[r.user_id] = false
      })
      success(`已成功取消 ${res.results?.length || 0} 位管理员`)
      clearSelection()
    } else {
      notifyError(res.message || '批量取消失败')
    }
  } catch (error) {
    notifyError('批量取消失败')
  }
}

const loadMessages = async () => {
  try {
    const res = await getTeamMessagesAPI(teamId.value)
    if (res.code === 200) {
      messages.value = res.data || []
      unreadCount.value = messages.value.filter(m => !m.is_read).length
    }
  } catch (error) {
    notifyError('加载消息失败')
  }
}

const showMessageDialog = () => {
  messageForm.value = { title: '', content: '', type: 'normal' }
  messageDialogVisible.value = true
}

const closeMessageDialog = () => {
  messageDialogVisible.value = false
}

const sendMessage = async () => {
  if (!messageForm.value.title) {
    warning('请输入消息标题')
    return
  }
  if (!messageForm.value.content) {
    warning('请输入消息内容')
    return
  }
  try {
    const res = await sendTeamMessageAPI(teamId.value, messageForm.value)
    if (res.code === 200) {
      success('发送成功')
      messages.value.unshift(res.data)
      unreadCount.value++
      closeMessageDialog()
    } else {
      notifyError(res.message || '发送失败')
    }
  } catch (error) {
    notifyError('发送失败')
  }
}

const sendBatchMessage = async () => {
  if (!messageForm.value.title) {
    warning('请输入消息标题')
    return
  }
  if (!messageForm.value.content) {
    warning('请输入消息内容')
    return
  }
  try {
    for (const userId of selectedMembers.value) {
      await sendTeamMessageAPI(teamId.value, {
        ...messageForm.value,
        receiver_id: userId
      })
    }
    success(`已成功发送消息给 ${selectedMembers.value.length} 位成员`)
    batchMessageDialogVisible.value = false
    messageForm.value = { title: '', content: '', type: 'normal' }
    clearSelection()
  } catch (error) {
    notifyError('批量发送失败')
  }
}

const markAsRead = async (message) => {
  if (message.is_read) return
  try {
    const res = await markMessageReadAPI(teamId.value, message.id)
    if (res.code === 200) {
      message.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    }
  } catch (error) {
    console.error('标记已读失败')
    notifyError('标记已读失败')
  }
}

const deleteMessage = async (message) => {
  if (!await confirmDanger('确定要删除此消息吗？')) return
  try {
    const res = await deleteTeamMessageAPI(teamId.value, message.id)
    if (res.code === 200) {
      success('删除成功')
      messages.value = messages.value.filter(m => m.id !== message.id)
      if (!message.is_read) {
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    } else {
      notifyError(res.message || '删除失败')
    }
  } catch (error) {
    notifyError('删除失败')
  }
}

const formatRelativeTime = (dateStr) => {
  if (!dateStr) return ''
  const now = new Date()
  const date = new Date(dateStr)
  const diff = now - date
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 30) return `${days}天前`
  return formatDate(dateStr)
}

const goBack = () => {
  window.location.href = '/team?from=teamHome'
}

const handleBackBtnClick = () => {
  window.location.href = '/team?from=teamHome'
}

const showEditDialog = () => {
  editForm.value = {
    name: team.value.name,
    description: team.value.description || '',
    color: team.value.color || 'var(--purple)',
    logo: team.value.logo || '',
    bg_image: team.value.bg_image || '',
    leader_id: team.value.leader_id || ''
  }
  editDialogVisible.value = true
}

const closeEditDialog = () => {
  editDialogVisible.value = false
}

const triggerLogoUpload = () => {
  logoInput.value?.click()
}

const triggerBgUpload = () => {
  bgInput.value?.click()
}

const handleLogoUpload = async (e) => {
  const file = e.target.files[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = async (event) => {
    try {
      const res = await uploadTeamLogoAPI({ image: event.target.result })
      if (res.code === 200) {
        editForm.value.logo = res.data.url
      }
    } catch (error) {
      notifyError('上传失败')
    }
  }
  reader.readAsDataURL(file)
}

const handleBgUpload = async (e) => {
  const file = e.target.files[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = async (event) => {
    try {
      const res = await uploadTeamBgAPI({ image: event.target.result })
      if (res.code === 200) {
        editForm.value.bg_image = res.data.url
      }
    } catch (error) {
      notifyError('上传失败')
    }
  }
  reader.readAsDataURL(file)
}

const saveTeamEdit = async () => {
  if (!editForm.value.name) {
    warning('请输入小组名称')
    return
  }

  saving.value = true
  try {
    const res = await updateTeamAPI(teamId.value, editForm.value)
    if (res.code === 200) {
      success('保存成功')
      team.value = { ...team.value, ...editForm.value }
      closeEditDialog()
    } else {
      notifyError(res.message || '保存失败')
    }
  } catch (error) {
    notifyError('保存失败')
  } finally {
    saving.value = false
  }
}

const showAnnouncementDialog = () => {
  announcementText.value = team.value.announcement || ''
  announcementDialogVisible.value = true
}

const closeAnnouncementDialog = () => {
  announcementDialogVisible.value = false
}

const saveAnnouncement = async () => {
  saving.value = true
  try {
    const res = await updateTeamAPI(teamId.value, { announcement: announcementText.value })
    if (res.code === 200) {
      success('公告发布成功')
      team.value.announcement = announcementText.value
      team.value.update_time = new Date().toISOString()
      closeAnnouncementDialog()
    } else {
      notifyError(res.message || '发布失败')
    }
  } catch (error) {
    notifyError('发布失败')
  } finally {
    saving.value = false
  }
}

const clearAnnouncement = () => {
  announcementText.value = ''
}

const formatAnnouncementTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)}天前`
  return date.toLocaleDateString('zh-CN')
}

const dismissAnnouncement = () => {
  if (dontShowAgain.value && team.value?.announcement) {
    const key = `${teamId.value}_${team.value.update_time}`
    dismissedAnnouncements.value.add(key)
    localStorage.setItem('dismissedAnnouncements', JSON.stringify([...dismissedAnnouncements.value]))
  }
  announcementPopupVisible.value = false
  dontShowAgain.value = false
}

const checkAnnouncementPopup = () => {
  if (!team.value?.announcement) return
  const key = `${teamId.value}_${team.value.update_time}`
  const saved = localStorage.getItem('dismissedAnnouncements')
  if (saved) {
    const dismissed = JSON.parse(saved)
    if (dismissed.includes(key)) return
  }
  if (!dismissedAnnouncements.value.has(key)) {
    announcementPopupVisible.value = true
  }
}

const joinTeam = async () => {
  try {
    const res = await setUserTeamAPI(teamId.value)
    if (res.code === 200) {
      success(`已成功加入【${team.value.name}】`)
      userStore.updateTeam(team.value.id, team.value.name, team.value.color)
      loadTeamData()
    } else {
      notifyError(res.message || '加入失败')
    }
  } catch (error) {
    notifyError('加入失败')
  }
}

const leaveTeam = async () => {
  if (!await confirmDanger('确定要退出该小组吗？')) return

  try {
    const res = await updateMyTeamAPI(null)
    if (res.code === 200) {
      success('已退出小组')
      userStore.updateTeam(null, null, null)
      loadTeamData()
    } else {
      notifyError(res.message || '退出失败')
    }
  } catch (error) {
    notifyError('退出失败')
  }
}

onMounted(() => {
  loadTeamData()
  window.addEventListener('resize', handleResize)
  const saved = localStorage.getItem('dismissedAnnouncements')
  if (saved) {
    dismissedAnnouncements.value = new Set(JSON.parse(saved))
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

const handleResize = () => {
  isMobile.value = window.innerWidth <= 768
}

watch(() => route.params.teamId, () => {
  loadTeamData()
})
</script>

<style scoped>
.team-home-page {
  min-height: 100vh;
  background: var(--bg-primary);
  --primary-color: var(--primary);
  --secondary-color: var(--primary-dark);
  --bg-page: var(--bg-primary);
  --bg-card: var(--bg-card);
  --bg-secondary: var(--bg-secondary);
  --bg-hover: var(--bg-card-hover);
  --text-primary: var(--text-primary);
  --text-secondary: var(--text-secondary);
  --text-muted: var(--text-muted);
  --border-color: var(--border-color);
}

.team-home-page.dark {
  --bg-page: var(--bg-primary);
  --bg-card: var(--bg-card);
  --bg-secondary: var(--bg-secondary);
  --bg-hover: var(--bg-card-hover);
  --text-primary: var(--text-primary);
  --text-secondary: var(--text-secondary);
  --text-muted: var(--text-muted);
  --border-color: var(--border-color);
}

.back-btn-fixed {
  position: fixed;
  top: 80px;
  left: 16px;
  z-index: 1000;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  border-radius: 50%;
  color: var(--color-on-brand);
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(255, 107, 53, 0.4);
}

.back-btn-fixed:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(255, 107, 53, 0.5);
}

.back-btn-fixed svg {
  width: 20px;
  height: 20px;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  gap: 20px;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.loading-text {
  color: var(--text-secondary);
  font-size: 14px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.not-found-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  gap: 16px;
  text-align: center;
  padding: 40px;
}

.not-found-icon {
  width: 80px;
  height: 80px;
  background: var(--bg-card);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40px;
  color: var(--text-muted);
}

.privacy-error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  gap: 16px;
  text-align: center;
  padding: 40px;
}

.privacy-icon {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #ef4444, #f97316);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.privacy-icon svg {
  width: 40px;
  height: 40px;
  color: var(--color-on-brand);
}

.privacy-error-state h2 {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.privacy-error-state p {
  color: var(--text-secondary);
  margin: 0;
}

.privacy-hint {
  font-size: 14px;
  color: var(--text-muted);
}

.team-banner {
  position: relative;
  padding: 60px 40px 40px;
  background-size: cover;
  background-position: center;
  overflow: hidden;
}

.banner-particles {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(255,107,53,0.9) 0%, rgba(247,147,30,0.85) 100%);
  animation: pulse-bg 8s ease-in-out infinite;
}

@keyframes pulse-bg {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.85; }
}

.banner-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(0,0,0,0.3) 0%, rgba(0,0,0,0.6) 100%);
}

.banner-content {
  position: relative;
  z-index: 1;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  max-width: 1400px;
  margin: 0 auto;
  gap: 40px;
}

.banner-left {
  display: flex;
  gap: 24px;
  align-items: center;
}

.team-avatar {
  width: 100px;
  height: 100px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40px;
  color: var(--color-on-brand);
  overflow: hidden;
  position: relative;
  box-shadow: 0 8px 32px rgba(0,0,0,0.3);
}

.team-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.admin-badge {
  position: absolute;
  top: -5px;
  right: -5px;
  width: 28px;
  height: 28px;
  background: var(--warning);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-on-brand);
  box-shadow: 0 2px 8px rgba(0,0,0,0.3);
}

.admin-badge svg {
  width: 16px;
  height: 16px;
}

.team-intro {
  color: var(--color-on-brand);
}

.team-title {
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 8px;
  text-shadow: 0 2px 4px rgba(0,0,0,0.3);
}

.team-subtitle {
  font-size: 14px;
  opacity: 0.9;
  margin: 0 0 16px;
}

.team-tags {
  display: flex;
  gap: 16px;
}

.team-tags .tag {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  background: rgba(255,255,255,0.2);
  padding: 6px 12px;
  border-radius: 20px;
  backdrop-filter: blur(10px);
}

.team-tags .tag svg {
  width: 14px;
  height: 14px;
}

.announcement-marquee {
  background: rgba(255,255,255,0.15);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  padding: 16px 20px;
  max-width: 400px;
  display: flex;
  gap: 12px;
  align-items: center;
}

.marquee-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  white-space: nowrap;
  color: var(--warning);
}

.marquee-label svg {
  width: 18px;
  height: 18px;
}

.marquee-content {
  overflow: hidden;
  flex: 1;
}

.marquee-content span {
  display: inline-block;
  animation: marquee 15s linear infinite;
  white-space: nowrap;
}

@keyframes marquee {
  0% { transform: translateX(100%); }
  100% { transform: translateX(-100%); }
}

.team-container {
  display: flex;
  max-width: 1400px;
  margin: 0 auto;
  padding: 24px;
  gap: 24px;
  position: relative;
}

.sidebar {
  width: 220px;
  flex-shrink: 0;
  transition: all 0.3s ease;
  position: relative;
}

.sidebar.collapsed {
  width: 70px;
}

.sidebar-toggle {
  position: absolute;
  right: -16px;
  top: 20px;
  width: 32px;
  height: 32px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  z-index: 10;
}

.sidebar-toggle:hover {
  background: var(--primary-color);
  color: var(--color-on-brand);
  border-color: var(--primary-color);
}

.sidebar-toggle svg {
  width: 16px;
  height: 16px;
}

.nav-section {
  background: var(--bg-card);
  border-radius: 16px;
  padding: 16px 12px;
  margin-bottom: 16px;
  border: 1px solid var(--border-color);
}

.nav-section-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 0 12px;
  margin-bottom: 12px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 10px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.nav-item:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.nav-item.active {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: var(--color-on-brand);
}

.nav-item svg {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.nav-item .badge {
  margin-left: auto;
  background: var(--primary);
  color: var(--color-on-brand);
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 600;
}

.nav-item:hover .badge {
  background: var(--primary-dark);
  color: var(--color-on-brand);
}

.badge-danger {
  background: var(--danger) !important;
}

.nav-item.action {
  color: var(--primary-color);
}

.nav-item.danger {
  color: var(--danger);
}

.nav-item.danger:hover {
  background: rgba(239,68,68,0.1);
}

.nav-bottom {
  margin-top: auto;
}

.main-content {
  flex: 1;
  min-width: 0;
  transition: all 0.3s ease;
}

.content-card {
  background: var(--bg-card);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 20px;
  border: 1px solid var(--border-color);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header h3 {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.card-header h3 svg {
  width: 20px;
  height: 20px;
  color: var(--primary-color);
}

.count-badge {
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-on-brand);
  box-shadow: 0 2px 8px rgba(249, 115, 22, 0.2);
}

.members-showcase {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: center;
}

.member-avatar-lg {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  overflow: hidden;
  position: relative;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 3px solid transparent;
  background: var(--bg-brand-gradient);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.15);
}

.member-avatar-lg:hover {
  transform: translateY(-4px) scale(1.08);
  border-color: var(--primary-color);
  box-shadow: 0 8px 20px rgba(249, 115, 22, 0.3);
}

.member-avatar-lg.is-leader {
  border-color: var(--warning);
  box-shadow: 0 4px 12px rgba(251, 191, 36, 0.3);
}

.member-avatar-lg.is-leader:hover {
  box-shadow: 0 8px 20px rgba(251, 191, 36, 0.5);
}

.member-avatar-lg img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.member-avatar-lg:hover img {
  transform: scale(1.05);
}

.member-avatar-lg > span {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 600;
  color: var(--color-on-brand);
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.member-tooltip {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%) translateY(-8px);
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 8px 12px;
  white-space: nowrap;
  opacity: 0;
  visibility: hidden;
  transition: all 0.2s ease;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  z-index: 100;
  display: flex;
  flex-direction: column;
  gap: 4px;
  text-align: center;
}

.member-tooltip strong {
  font-size: 13px;
  color: var(--text-primary);
}

.member-tooltip .leader-tag {
  font-size: 10px;
  color: var(--warning);
}

.member-tooltip span:not(.leader-tag) {
  font-size: 11px;
  color: var(--text-muted);
}

.member-avatar-lg:hover .member-tooltip {
  opacity: 1;
  visibility: visible;
}

.more-members {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-muted);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-box {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  gap: 16px;
  align-items: center;
  transition: all 0.3s ease;
}

.stat-box:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon svg {
  width: 24px;
  height: 24px;
  color: var(--color-on-brand);
}

.stat-icon.blue { background: linear-gradient(135deg, #3b82f6, #60a5fa); }
.stat-icon.green { background: linear-gradient(135deg, #10b981, #34d399); }
.stat-icon.orange { background: linear-gradient(135deg, #f59e0b, #fbbf24); }
.stat-icon.purple { background: linear-gradient(135deg, #8b5cf6, #a78bfa); }

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-text {
  font-size: 12px;
  color: var(--text-muted);
}

.blogger-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 12px;
}

.blogger-mini-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 12px 8px;
  background: var(--bg-secondary);
  border-radius: 12px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.blogger-mini-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
}

.mini-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-brand-gradient);
}

.mini-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.mini-avatar > span {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-on-brand);
}

.mini-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.mini-name {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-primary);
  text-align: center;
  max-width: 70px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.mini-category {
  font-size: 10px;
}

.members-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.member-list-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  background: var(--bg-secondary);
  border-radius: 12px;
  transition: all 0.2s ease;
}

.member-list-item:hover {
  transform: translateX(4px);
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.member-avatar-md {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background: var(--bg-brand-gradient);
}

.member-avatar-md img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.member-avatar-md > span {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-on-brand);
}

.leader-indicator {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 14px;
  height: 14px;
  background: var(--warning);
  border-radius: 50%;
  border: 2px solid var(--bg-secondary);
}

.member-list-item.selected {
  background: rgba(249, 115, 22, 0.1);
  border: 1px solid var(--primary-color);
}

.member-checkbox {
  display: flex;
  align-items: center;
  justify-content: center;
}

.member-checkbox input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
  accent-color: var(--primary-color);
}

.batch-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}

.selected-count {
  font-size: 13px;
  color: var(--primary-color);
  font-weight: 500;
}

.member-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.member-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.member-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.leader-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  color: var(--warning);
  background: rgba(251,191,36,0.1);
  padding: 2px 8px;
  border-radius: 10px;
}

.leader-badge svg {
  width: 10px;
  height: 10px;
}

.member-username {
  font-size: 12px;
  color: var(--text-muted);
}

.member-join-date {
  font-size: 12px;
  color: var(--text-muted);
}

.blogger-full-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.blogger-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 20px;
  background: var(--bg-secondary);
  border-radius: 12px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.blogger-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
}

.blogger-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-brand-gradient);
}

.blogger-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.blogger-avatar > span {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-on-brand);
}

.blogger-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.blogger-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.blogger-category {
  font-size: 12px;
}

.empty-tip {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
}

.right-sidebar {
  width: 280px;
  flex-shrink: 0;
}

.sidebar-card {
  background: var(--bg-card);
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 16px;
  border: 1px solid var(--border-color);
}

.sidebar-card .card-header {
  margin-bottom: 16px;
}

.sidebar-card .card-header h3 {
  font-size: 14px;
}

.leader-profile {
  display: flex;
  align-items: center;
  gap: 16px;
}

.leader-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-brand-gradient);
}

.leader-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.leader-avatar > span {
  font-size: 22px;
  font-weight: 600;
  color: var(--color-on-brand);
}

.leader-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.leader-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.leader-username {
  font-size: 12px;
  color: var(--text-muted);
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: var(--bg-hover);
  border: none;
  border-radius: 10px;
  font-size: 13px;
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn:hover {
  background: var(--bg-brand-gradient);
  color: var(--color-on-brand);
}

.action-btn svg {
  width: 18px;
  height: 18px;
}

.top-members {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.top-member-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.rank {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  background: var(--bg-secondary);
  color: var(--text-muted);
}

.rank.gold {
  background: linear-gradient(135deg, #fbbf24, #f59e0b);
  color: var(--color-on-brand);
}

.rank.silver {
  background: linear-gradient(135deg, #9ca3af, #6b7280);
  color: var(--color-on-brand);
}

.rank.bronze {
  background: var(--bg-warning-gradient-alt);
  color: var(--color-on-brand);
}

.announcement-sidebar {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.announcement-content-preview {
  background: var(--bg-secondary);
  border-radius: 10px;
  padding: 14px;
}

.announcement-content-preview p {
  margin: 0;
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-primary);
  word-break: break-word;
}

.announcement-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  color: var(--text-muted);
}

.announcement-from {
  font-weight: 500;
}

.announcement-popup-modal {
  width: 480px;
  max-width: 90vw;
}

.announcement-popup-icon {
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, #f59e0b, #d97706);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.announcement-popup-icon svg {
  width: 22px;
  height: 22px;
  color: var(--color-on-brand);
}

.announcement-popup-modal .modal-header {
  display: flex;
  align-items: center;
  gap: 14px;
}

.announcement-popup-modal .modal-header h3 {
  margin: 0;
  font-size: 18px;
}

.announcement-popup-content {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 20px;
}

.announcement-popup-content .announcement-text {
  font-size: 15px;
  line-height: 1.7;
  color: var(--text-primary);
  margin-top: 12px;
  word-break: break-word;
}

.dont-show-again {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 13px;
  color: var(--text-secondary);
}

.dont-show-again input {
  width: 16px;
  height: 16px;
  cursor: pointer;
  accent-color: var(--primary);
}

.modal-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.member-avatar-sm {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-brand-gradient);
}

.member-avatar-sm > span {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-on-brand);
}

.member-name-sm {
  flex: 1;
  font-size: 13px;
  color: var(--text-primary);
}

.is-leader-tag {
  font-size: 10px;
  color: var(--warning);
  background: rgba(251,191,36,0.1);
  padding: 2px 6px;
  border-radius: 8px;
}

.modal-overlay {
  position: fixed !important;
  inset: 0 !important;
  background: rgba(0,0,0,0.6) !important;
  backdrop-filter: blur(4px) !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  z-index: 2000 !important;
  animation: fadeIn 0.2s ease;
}

html.dark .modal-overlay {
  background: rgba(0, 0, 0, 0.8) !important;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal {
  background: var(--bg-card) !important;
  border-radius: 20px !important;
  width: 90% !important;
  max-width: 500px !important;
  max-height: 90vh !important;
  overflow: auto !important;
  animation: slideUp 0.3s ease !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2) !important;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 2px solid var(--border-color);
  background: var(--bg-card-hover);
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 700;
  margin: 0;
  color: var(--text-primary);
}

.close-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: var(--danger);
  border-color: var(--danger);
  color: var(--color-on-brand);
}

.close-btn svg {
  width: 18px;
  height: 18px;
}

.post-detail-actions {
  display: flex;
  gap: 8px;
}

.action-icon-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-secondary);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--text-secondary);
}

.action-icon-btn:hover {
  background: var(--primary-color);
  color: var(--color-on-brand);
}

.action-icon-btn svg {
  width: 18px;
  height: 18px;
}

.modal-body {
  padding: 24px;
  background: var(--bg-card);
}

.logs-modal .modal-body {
  max-height: 400px;
  overflow-y: auto;
}

.logs-body {
  padding: 16px;
}

.log-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.log-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: var(--bg-card-hover);
  border: 1px solid var(--border-color);
  border-radius: 10px;
}

.log-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 16px;
}

.log-content {
  flex: 1;
}

.log-action {
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.log-detail {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.log-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-muted);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 24px;
  border-top: 2px solid var(--border-color);
  background: var(--bg-card-hover);
}

.form-field {
  margin-bottom: 20px;
}

.form-field label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.form-field input[type="text"],
.form-field textarea,
.form-field select {
  width: 100%;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 2px solid var(--border-color);
  border-radius: 10px;
  font-size: 14px;
  color: var(--text-primary);
  transition: all 0.2s ease;
}

.form-field input:focus,
.form-field textarea:focus,
.form-field select:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(255, 107, 53, 0.1);
  background: var(--bg-card);
}

.form-field textarea {
  resize: vertical;
  min-height: 100px;
}

.image-upload-row {
  display: flex;
  align-items: center;
  gap: 16px;
}

.preview-avatar {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: var(--bg-brand-gradient);
}

.preview-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-avatar > span {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-on-brand);
}

.preview-bg {
  width: 120px;
  height: 64px;
  border-radius: 12px;
  background: var(--bg-secondary);
  background-size: cover;
  background-position: center;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  font-size: 12px;
  border: 1px dashed var(--border-color);
}

.color-picker-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.color-picker-row input[type="color"] {
  width: 48px;
  height: 48px;
  border: none;
  border-radius: 10px;
  cursor: pointer;
}

.color-picker-row input[type="text"] {
  flex: 1;
}

.btn {
  padding: 10px 20px;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.btn-primary {
  background: var(--bg-brand-gradient);
  color: var(--color-on-brand);
  box-shadow: var(--shadow-sm);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-secondary {
  background: var(--bg-hover);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: var(--bg-tertiary);
}

.btn-danger {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: var(--color-on-brand);
}

.btn-danger:hover {
  box-shadow: var(--shadow-md);
}

.btn-block {
  width: 100%;
  justify-content: center;
}

.btn svg {
  width: 16px;
  height: 16px;
}

.btn-sm {
  padding: 6px 14px;
  font-size: 12px;
}

.reveal {
  animation: revealUp 0.5s ease forwards;
  animation-delay: calc(var(--delay, 0) * 1s);
  opacity: 0;
}

@keyframes revealUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 1200px) {
  .right-sidebar {
    display: none;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .blogger-grid {
    grid-template-columns: repeat(4, 1fr);
  }

  .blogger-full-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .team-container {
    padding-bottom: 92px;
  }

  .sidebar {
    display: none;
  }

  .sidebar.collapsed {
    display: none;
  }

  .banner-content {
    flex-direction: column;
  }

  .announcement-marquee {
    max-width: 100%;
  }

  .team-banner {
    min-height: 180px;
  }

  .banner-content {
    gap: 12px;
  }

  .team-avatar {
    width: 56px;
    height: 56px;
  }

  .blogger-grid {
    grid-template-columns: repeat(3, 1fr);
  }

  .blogger-full-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  /* 确保内容不会被底部导航遮挡 */
  .main-content {
    padding-bottom: 92px;
  }
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.post-item {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid transparent;
}

.post-item:hover {
  transform: translateY(-2px);
  border-color: var(--primary-color);
  box-shadow: 0 8px 24px rgba(255, 107, 53, 0.15);
}

.post-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.post-author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-brand-gradient);
}

.post-author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.post-author-avatar > span {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-on-brand);
}

.post-author-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.post-author-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.post-category {
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 10px;
  color: var(--color-on-brand);
  width: fit-content;
}

.post-time {
  margin-left: auto;
  font-size: 12px;
  color: var(--text-muted);
}

.post-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px;
  line-height: 1.4;
}

.post-excerpt {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0 0 12px;
}

.post-footer {
  display: flex;
  gap: 20px;
}

.post-stat {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-muted);
}

.post-stat svg {
  width: 14px;
  height: 14px;
}

/* 帖子详情弹窗 - 使用 Teleport 到 body */
.post-detail-modal-wrapper {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  animation: modalFadeIn 0.3s ease;
}

html.dark .post-detail-modal-wrapper {
  background: rgba(0, 0, 0, 0.8);
}

@keyframes modalFadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.post-detail-modal-container {
  background: var(--bg-card);
  border-radius: 16px;
  width: 100%;
  max-width: 700px;
  max-height: calc(100vh - 40px);
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  overflow: hidden;
  animation: modalSlideUp 0.3s ease;
}

@keyframes modalSlideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.post-detail-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.post-detail-modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.post-detail-modal-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.modal-action-btn, .modal-close-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: none;
  background: var(--bg-hover);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.modal-action-btn:hover, .modal-close-btn:hover {
  background: var(--bg-brand-gradient);
  color: var(--color-on-brand);
}

.modal-action-btn:focus, .modal-close-btn:focus {
  outline: 2px solid #ff6b35;
  outline-offset: 2px;
}

.modal-close-btn:hover {
  background: var(--danger);
}

.modal-action-btn svg, .modal-close-btn svg {
  width: 18px;
  height: 18px;
}

.post-detail-modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.post-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 12px;
}

.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.author-avatar span {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-brand-gradient);
  color: var(--color-on-brand);
  font-weight: 600;
  font-size: 16px;
}

.author-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.author-name {
  font-weight: 600;
  font-size: 14px;
}

.post-category {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 8px;
  font-size: 12px;
  color: var(--color-on-brand);
  font-weight: 500;
}

.post-time {
  font-size: 13px;
  color: var(--text-muted);
}

.post-content {
  margin-bottom: 20px;
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-word;
}

.post-stats {
  display: flex;
  gap: 20px;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--border-color);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: var(--text-muted);
}

.stat-item svg {
  width: 16px;
  height: 16px;
}

.comments-section {
  margin-top: 20px;
}

.comments-title {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 20px;
}

.comment-item {
  display: flex;
  gap: 12px;
}

.comment-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.comment-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.comment-avatar span {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-brand-gradient);
  color: var(--color-on-brand);
  font-weight: 600;
  font-size: 14px;
}

.comment-body {
  flex: 1;
  min-width: 0;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.comment-author {
  font-weight: 600;
  font-size: 14px;
}

.comment-time {
  font-size: 12px;
  color: var(--text-muted);
}

.comment-text {
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}

.no-comments {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-muted);
  font-size: 14px;
}

.comment-input-area {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.mention-input-wrapper {
  flex: 1;
  position: relative;
}

.mention-input-wrapper input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
  outline: none;
  transition: border-color 0.2s;
}

.mention-input-wrapper input:focus {
  border-color: var(--primary-color);
}

.mention-dropdown {
  position: absolute;
  bottom: 100%;
  left: 0;
  right: 0;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  box-shadow: 0 -4px 12px rgba(0, 0, 0, 0.15);
  max-height: 200px;
  overflow-y: auto;
  margin-bottom: 8px;
  z-index: 10000;
}

.mention-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.mention-item:hover {
  background: var(--bg-hover);
}

.mention-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.mention-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.mention-avatar span {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-brand-gradient);
  color: var(--color-on-brand);
  font-weight: 600;
  font-size: 12px;
}

.mention-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.mention-name {
  font-size: 13px;
  font-weight: 500;
}

.mention-username {
  font-size: 12px;
  color: var(--text-muted);
}

.mention-empty {
  padding: 8px 12px;
  color: var(--text-muted);
  font-size: 13px;
}

.submit-comment-btn {
  flex-shrink: 0;
}

.post-detail-modal-footer {
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
}

.post-detail-modal-footer .btn {
  padding: 8px 16px;
}

.forum-header {
  flex-wrap: wrap;
  gap: 12px;
  justify-content: space-between;
}

.forum-tabs {
  display: flex;
  gap: 8px;
}

.forum-tab-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.forum-tab-btn svg {
  width: 16px;
  height: 16px;
}

.forum-tab-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.forum-tab-btn.active {
  background: var(--primary-color);
  border-color: var(--primary-color);
  color: var(--color-on-brand);
}

.category-filter {
  display: flex;
  gap: 8px;
  padding: 12px 0;
  flex-wrap: wrap;
}

.category-btn {
  padding: 6px 14px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 20px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.category-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.category-btn.active {
  background: var(--bg-brand-gradient);
  border-color: var(--primary);
  color: var(--color-on-brand);
}

.forum-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  transition: all 0.2s ease;
}

.search-box:focus-within {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(255, 107, 53, 0.1);
}

.search-box svg {
  width: 16px;
  height: 16px;
  color: var(--text-muted);
  flex-shrink: 0;
}

.search-box input {
  border: none;
  background: transparent;
  font-size: 14px;
  color: var(--text-primary);
  outline: none;
  width: 160px;
}

.search-box input::placeholder {
  color: var(--text-muted);
}

.post-badges {
  display: flex;
  gap: 8px;
  margin-bottom: 10px;
}

.post-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.post-badge.pinned {
  background: var(--bg-brand-gradient);
  color: var(--color-on-brand);
}

.post-badge.featured {
  background: linear-gradient(135deg, #ffd700, #ffaa00);
  color: var(--text-primary);
}

.post-footer {
  display: flex;
  align-items: center;
  gap: 12px;
}

.post-action {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.post-action svg {
  width: 14px;
  height: 14px;
}

.post-action:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.post-action.active {
  color: var(--primary);
  border-color: var(--primary);
  background: rgba(255, 107, 53, 0.1);
}

.post-admin-actions {
  display: flex;
  gap: 4px;
  margin-left: auto;
}

.admin-btn {
  padding: 4px 8px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.2s ease;
}

.admin-btn svg {
  width: 14px;
  height: 14px;
}

.admin-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.admin-btn.active {
  background: var(--primary-color);
  border-color: var(--primary-color);
  color: var(--color-on-brand);
}

.badge-danger {
  background: var(--danger);
  color: var(--color-on-brand);
}

.role-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.role-badge svg {
  width: 12px;
  height: 12px;
}

.role-badge.leader {
  background: var(--bg-brand-gradient);
  color: var(--color-on-brand);
}

.role-badge.admin {
  background: var(--bg-purple-gradient-btn);
  color: var(--color-on-brand);
}

.role-badge.member {
  background: var(--bg-secondary);
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.member-actions {
  display: flex;
  gap: 8px;
  margin-left: auto;
}

.action-btn-sm {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  padding: 0;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn-sm svg {
  width: 14px;
  height: 14px;
}

.action-btn-sm:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.action-btn-sm.active {
  background: linear-gradient(135deg, var(--purple), #a855f7);
  border-color: var(--purple);
  color: var(--color-on-brand);
}

.action-btn-sm.danger:hover {
  background: var(--danger);
  border-color: var(--danger);
  color: var(--color-on-brand);
}

.message-tab .messages-list,
.settings-tab .settings-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message-item {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.message-item:hover {
  background: var(--bg-hover);
}

.message-item.unread {
  border-color: var(--primary-color);
  background: rgba(255, 107, 53, 0.05);
}

.message-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--bg-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  color: var(--text-secondary);
  flex-shrink: 0;
  overflow: hidden;
}

.message-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.message-content {
  flex: 1;
  min-width: 0;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.message-title {
  font-weight: 600;
  color: var(--text-primary);
}

.message-time {
  font-size: 12px;
  color: var(--text-muted);
}

.message-text {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 6px;
  line-height: 1.5;
}

.message-sender {
  font-size: 12px;
  color: var(--text-muted);
}

.message-actions {
  display: flex;
  gap: 4px;
}

.settings-section {
  padding: 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
}

.settings-section h4 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.settings-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.settings-info p {
  font-size: 14px;
  color: var(--text-secondary);
}

.settings-info strong {
  color: var(--text-primary);
}

.settings-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.settings-actions .btn {
  display: flex;
  align-items: center;
  gap: 8px;
}

.settings-actions .btn svg {
  width: 16px;
  height: 16px;
}

.mobile-bottom-nav {
  display: none;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 68px;
  background-color: var(--bg-card) !important;
  background-image: none !important;
  border-top: 1px solid var(--border-color);
  padding: 0 8px;
  padding-bottom: env(safe-area-inset-bottom);
  z-index: 1000;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 -2px 10px rgba(0,0,0,0.08);
}

.team-home-page.dark .mobile-bottom-nav {
  background-color: var(--bg-card) !important;
  background-image: none !important;
}

@media (max-width: 768px) {
  .mobile-bottom-nav {
    display: flex;
  }
}

.mobile-nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 6px 10px;
  color: var(--text-muted);
  text-decoration: none;
  font-size: 10px;
  position: relative;
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

.mobile-nav-item.active {
  color: var(--primary-color);
  background: rgba(249, 115, 22, 0.08);
}

.mobile-nav-item.active svg {
  transform: scale(1.05);
}

.mobile-badge {
  position: absolute;
  top: -2px;
  right: 2px;
  background: var(--danger);
  color: var(--color-on-brand);
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 10px;
  min-width: 16px;
  text-align: center;
}

.all-users-modal,
.private-message-modal {
  background: var(--bg-card);
  border-radius: 16px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
}

.all-users-modal .modal-header,
.private-message-modal .modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
}

.all-users-modal .modal-header h3,
.private-message-modal .modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.modal-search {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  position: relative;
}

.modal-search svg {
  position: absolute;
  left: 28px;
  top: 50%;
  transform: translateY(-50%);
  width: 18px;
  height: 18px;
  color: var(--text-muted);
}

.modal-search input {
  width: 100%;
  padding: 10px 12px 10px 36px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
}

.users-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  transition: background 0.2s;
}

.user-item:hover {
  background: var(--bg-secondary);
}

.user-item .user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  overflow: hidden;
}

.user-item .user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-item .user-avatar span {
  color: var(--color-on-brand);
  font-size: 16px;
  font-weight: 600;
}

.user-item .user-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.user-item .user-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.user-item .user-team {
  font-size: 12px;
  color: var(--text-muted);
}

.user-item .user-role {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 4px;
  width: fit-content;
}

.user-item .user-role.admin {
  background: var(--warning-bg);
  color: var(--warning);
}

.user-item .user-role.member {
  background: var(--info-bg);
  color: var(--info);
}

.user-item .user-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.empty-users {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
}

.private-message-modal .modal-body {
  padding: 20px;
}

.private-message-modal .modal-body textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
  resize: vertical;
  min-height: 120px;
}

.private-message-modal .modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid var(--border-color);
  background: var(--bg-secondary);
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
    gap: 8px;
  }

  .blogger-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .blogger-full-grid {
    grid-template-columns: 1fr;
  }

  .team-home {
    padding: 10px;
  }

  .section-card {
    padding: 14px;
    border-radius: 12px;
  }

  .team-title {
    font-size: 20px !important;
  }

  .team-subtitle {
    font-size: 13px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .team-banner {
    min-height: 140px !important;
  }

  .back-btn-fixed {
    width: 32px;
    height: 32px;
    top: 8px;
    left: 8px;
  }

  .back-btn-fixed svg {
    width: 16px;
    height: 16px;
  }

  .mobile-bottom-nav {
    padding: 4px 8px;
  }

  .post-item {
    padding: 14px;
    border-radius: 12px;
  }
}
</style>
