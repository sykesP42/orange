<template>
  <div class="team-page">
    <div class="page-header">
      <div class="header-left">
        <h1>团队中心</h1>
        <p>浏览团队、加入小组、结识伙伴</p>
      </div>
      <div class="header-right" v-if="isAtTop">
        <div class="action-buttons">
          <button @click="goToChat" class="action-btn chat-btn">
            <div class="btn-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 0 1-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
              </svg>
            </div>
            <div class="btn-text">
              <span class="btn-title">消息中心</span>
              <span class="btn-subtitle">即时通讯</span>
            </div>
          </button>
          <button @click="goToUsers" class="action-btn public-forum-btn">
            <div class="btn-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              </svg>
            </div>
            <div class="btn-text">
              <span class="btn-title">用户列表</span>
              <span class="btn-subtitle">查看成员</span>
            </div>
          </button>
          <button v-if="hasTeam" class="action-btn my-team-btn" @click="goToMyTeam" title="进入我的团队主页">
            <div class="btn-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
                <polyline points="9 22 9 12 15 12 15 22"/>
              </svg>
            </div>
            <div class="btn-text">
              <span class="btn-title">我的团队</span>
              <span class="btn-subtitle">团队主页</span>
            </div>
          </button>
        </div>
      </div>
    </div>

    <div class="fixed-action-buttons" :class="{ 'show': showButtons, 'expanded': expanded }" @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave">
      <button @click="goToChat" class="fixed-action-btn chat-btn" title="前往消息中心">
        <div class="btn-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 0 1-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
          </svg>
        </div>
        <div class="btn-text">
          <span class="btn-title">消息中心</span>
          <span class="btn-subtitle">即时通讯</span>
        </div>
      </button>
      <button @click="goToUsers" class="fixed-action-btn public-forum-btn" title="查看用户列表">
        <div class="btn-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
        </div>
        <div class="btn-text">
          <span class="btn-title">用户列表</span>
          <span class="btn-subtitle">查看成员</span>
        </div>
      </button>
      <button v-if="hasTeam" class="fixed-action-btn my-team-btn" @click="goToMyTeam" title="进入我的团队主页">
        <div class="btn-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
            <polyline points="9 22 9 12 15 12 15 22"/>
          </svg>
        </div>
        <div class="btn-text">
          <span class="btn-title">我的团队</span>
          <span class="btn-subtitle">团队主页</span>
        </div>
      </button>
    </div>

    <button class="floating-trigger" @click="toggleExpanded" :class="{ 'active': expanded }" title="快捷菜单">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="1"/>
        <circle cx="12" cy="5" r="1"/>
        <circle cx="12" cy="19" r="1"/>
      </svg>
    </button>

    <div class="view-toggle" v-if="teams.length > 0">
      <button class="toggle-btn" :class="{ active: viewMode === 'card' }" @click="viewMode = 'card'">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="7" height="7" rx="1"/>
          <rect x="14" y="3" width="7" height="7" rx="1"/>
          <rect x="3" y="14" width="7" height="7" rx="1"/>
          <rect x="14" y="14" width="7" height="7" rx="1"/>
        </svg>
        卡片
      </button>
      <button class="toggle-btn" :class="{ active: viewMode === 'list' }" @click="viewMode = 'list'">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="3" y1="6" x2="21" y2="6"/>
          <line x1="3" y1="12" x2="21" y2="12"/>
          <line x1="3" y1="18" x2="21" y2="18"/>
        </svg>
        列表
      </button>
    </div>

    <div class="team-actions">
      <template v-if="!hasTeam && !hasPendingRequest">
        <button class="btn btn-primary" @click="openCreateTeam">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          {{ userStore.role === 'admin' ? '新建团队' : '申请创建团队' }}
        </button>
      </template>

      <template v-else-if="hasPendingRequest">
        <div class="pending-request-card">
          <div class="pending-icon">⏳</div>
          <div class="pending-info">
            <span class="pending-title">团队创建申请审核中</span>
            <span class="pending-desc">{{ pendingRequestInfo?.name || '您的申请' }}</span>
            <span class="pending-time">提交时间: {{ formatTime(pendingRequestInfo?.create_time) }}</span>
          </div>
          <button class="btn btn-sm btn-secondary" @click="checkPendingStatus">刷新状态</button>
        </div>
      </template>

      <template v-else-if="hasTeam">
        <div class="has-team-info">
          <span class="info-text">您已加入团队</span>
          <button class="btn btn-secondary" @click="goToMyTeam">查看我的团队</button>
        </div>
      </template>
    </div>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else-if="viewMode === 'card'" class="team-grid">
      <div v-for="team in teams" :key="team.id" class="team-card" @click="openTeamDetail(team)">
        <div class="card-inner">
          <div class="card-bg"></div>
          <div class="card-avatar-bg" :style="team.bg_image ? { backgroundImage: `url(${team.bg_image})` } : {}">
            <img v-if="team.bg_image" :src="team.bg_image" :alt="team.name" />
          </div>
          <div class="team-card-content">
            <div class="team-logo" :style="{ backgroundColor: team.color }">
              <img v-if="team.logo" :src="team.logo" :alt="team.name" />
              <span v-else>{{ team.name.charAt(0) }}</span>
            </div>
            <h3 class="team-name">{{ team.name }}</h3>
            <p class="team-leader">组长: {{ team.leader_name || '未知' }}</p>
            <p class="team-desc">{{ team.description || '暂无描述' }}</p>
            <div class="team-stats">
              <span class="stat">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                  <circle cx="9" cy="7" r="4"/>
                </svg>
                {{ getTeamMembers(team.id).length }} 名成员
              </span>
              <span class="stat">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                  <line x1="16" y1="2" x2="16" y2="6"/>
                  <line x1="8" y1="2" x2="8" y2="6"/>
                  <line x1="3" y1="10" x2="21" y2="10"/>
                </svg>
                {{ formatDate(team.create_time) }}
              </span>
            </div>
            <div class="team-members-preview">
              <div v-for="user in getTeamMembers(team.id).slice(0, 5)" :key="user.id" class="member-avatar" :style="{ backgroundColor: team.color }">
                {{ user.real_name?.charAt(0) || 'U' }}
              </div>
              <div v-if="getTeamMembers(team.id).length > 5" class="member-avatar more">
                +{{ getTeamMembers(team.id).length - 5 }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="team-list">
      <div class="list-header">
        <div class="list-cell logo-cell">Logo</div>
        <div class="list-cell name-cell">团队名称</div>
        <div class="list-cell leader-cell">组长</div>
        <div class="list-cell desc-cell">描述</div>
        <div class="list-cell members-cell">成员数</div>
        <div class="list-cell action-cell">操作</div>
      </div>
      <div v-for="team in teams" :key="team.id" class="list-row" @click="openTeamDetail(team)">
        <div class="list-cell logo-cell">
          <div class="team-logo-sm" :style="{ backgroundColor: team.color }">
            <img v-if="team.logo" :src="team.logo" :alt="team.name" />
            <span v-else>{{ team.name.charAt(0) }}</span>
          </div>
        </div>
        <div class="list-cell name-cell">{{ team.name }}</div>
        <div class="list-cell leader-cell">{{ team.leader_name || '未知' }}</div>
        <div class="list-cell desc-cell">{{ team.description || '-' }}</div>
        <div class="list-cell members-cell">{{ getTeamMembers(team.id).length }}</div>
        <div class="list-cell action-cell">
          <button v-if="!userStore.team_id || userStore.team_id !== team.id" class="btn btn-sm btn-primary" @click.stop="joinTeam(team)">
            加入
          </button>
          <span v-else class="joined-tag">已加入</span>
        </div>
      </div>
    </div>

    <div v-if="teams.length === 0 && !loading" class="empty-state">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
        <circle cx="9" cy="7" r="4"/>
        <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
        <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
      </svg>
      <p>暂无团队</p>
      <button v-if="userStore.role === 'admin'" class="btn btn-primary" @click="openCreateTeam">创建第一个团队</button>
    </div>

    <div v-if="detailPanelVisible" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="closeDetailPanel">
      <div class="modal team-detail-modal" style="background:var(--bg-card);">
        <div class="modal-header">
          <h3>团队详情</h3>
          <button class="close-btn" @click="closeDetailPanel">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body" v-if="currentTeam">
          <div class="team-detail-header" :style="currentTeam.bg_image ? { backgroundImage: `url(${currentTeam.bg_image})` } : { background: currentTeam.color }">
            <div class="team-logo-xl" :style="{ backgroundColor: currentTeam.color }">
              <img v-if="currentTeam.logo" :src="currentTeam.logo" :alt="currentTeam.name" />
              <span v-else>{{ currentTeam.name.charAt(0) }}</span>
            </div>
            <h2 class="team-title">{{ currentTeam.name }}</h2>
            <p class="team-description">{{ currentTeam.description || '暂无描述' }}</p>
          </div>

          <div class="team-info-stats">
            <div class="stat-item">
              <span class="stat-value">{{ getTeamMembers(currentTeam.id).length }}</span>
              <span class="stat-label">成员</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ currentTeam.leader_name || '未知' }}</span>
              <span class="stat-label">组长</span>
            </div>
          </div>

          <div v-if="currentTeam.announcement" class="announcement-bar">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 17H2a3 3 0 0 0 3-3V9a7 7 0 0 1 14 0v5a3 3 0 0 0 3 3zm-8.27 4a2 2 0 0 1-3.46 0"/>
            </svg>
            <span>{{ currentTeam.announcement }}</span>
          </div>

          <div class="action-buttons">
            <button v-if="!userStore.team_id || String(userStore.team_id) !== String(currentTeam.id)" class="btn btn-primary btn-lg" @click="joinTeam(currentTeam)">
              加入团队
            </button>
            <template v-else>
              <button class="btn btn-secondary" @click="goToTeamHome(currentTeam)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
                  <polyline points="9 22 9 12 15 12 15 22"/>
                </svg>
                进入小组主页
              </button>
              <button v-if="isTeamLeaderOrAdmin(currentTeam)" class="btn btn-secondary" @click="editTeam(currentTeam)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
                编辑团队
              </button>
              <button v-if="isTeamLeaderOrAdmin(currentTeam)" class="btn btn-secondary" @click="showAnnouncementDialog">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 17H2a3 3 0 0 0 3-3V9a7 7 0 0 1 14 0v5a3 3 0 0 0 3 3zm-8.27 4a2 2 0 0 1-3.46 0"/>
                </svg>
                发布公告
              </button>
              <button v-if="!isTeamLeaderOrAdmin(currentTeam)" class="btn btn-danger" @click="leaveTeam">
                退出团队
              </button>
            </template>
          </div>

          <div v-if="userStore.team_id && String(userStore.team_id) !== String(currentTeam.id)" class="action-buttons" style="margin-top: 12px;">
            <button class="btn btn-secondary btn-block" @click="goToTeamHome(currentTeam)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
                <polyline points="9 22 9 12 15 12 15 22"/>
              </svg>
              访问该小组主页
            </button>
            <button class="btn btn-secondary btn-block" @click="goToUsers" style="margin-top: 8px;">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              </svg>
              查看用户列表
            </button>
          </div>

          <div class="members-section">
            <h3>团队成员</h3>
            <div class="members-list">
              <div v-for="user in getTeamMembers(currentTeam.id)" :key="user.id" class="member-item" :class="{ 'is-leader': user.id === currentTeam.leader_id }">
                <div class="member-avatar-lg" :style="{ backgroundColor: currentTeam.color }">
                  <img v-if="user.avatar" :src="user.avatar" :alt="user.real_name" v-avatar />
                  <span v-else>{{ user.real_name?.charAt(0) || 'U' }}</span>
                </div>
                <div class="member-info">
                  <span class="member-name">{{ user.real_name }} <span v-if="user.id === currentTeam.leader_id" class="leader-badge">组长</span></span>
                  <span class="member-username">@{{ user.username }}</span>
                </div>
              </div>
              <div v-if="getTeamMembers(currentTeam.id).length === 0" class="no-members">
                暂无成员
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="formPanelVisible" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="closeFormPanel">
      <div class="modal team-form-modal" style="background:var(--bg-card);">
        <div class="modal-header">
          <h3>{{ isEditing ? '编辑团队' : (userStore.role === 'admin' ? '新建团队' : '申请创建团队') }}</h3>
          <button class="close-btn" @click="closeFormPanel">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div v-if="formError" class="error-message">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="15" y1="9" x2="9" y2="15"/>
              <line x1="9" y1="9" x2="15" y2="15"/>
            </svg>
            {{ formError }}
          </div>
          <div v-if="!isEditing && userStore.role !== 'admin'" class="info-message">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="16" x2="12" y2="12"/>
              <line x1="12" y1="8" x2="12.01" y2="8"/>
            </svg>
            <span>普通用户创建团队需要管理员审核，审核通过后即可使用</span>
          </div>
          <div class="form-row logo-row">
            <div class="form-field">
              <label>团队Logo</label>
              <div class="image-upload-preview">
                <div class="preview-box logo-preview" :style="{ backgroundColor: formData.color }">
                  <img v-if="formData.logo" :src="formData.logo" alt="logo" />
                  <span v-else>{{ formData.name?.charAt(0) || 'L' }}</span>
                </div>
                <div class="upload-actions">
                  <button type="button" class="btn btn-sm btn-secondary" @click="triggerLogoUpload">上传Logo</button>
                  <input type="file" ref="logoInput" accept="image/*" style="display: none" @change="handleLogoUpload" />
                  <button v-if="formData.logo" type="button" class="btn btn-sm btn-danger" @click="formData.logo = ''">移除</button>
                </div>
              </div>
            </div>
            <div class="form-field">
              <label>团队背景</label>
              <div class="image-upload-preview">
                <div class="preview-box bg-preview" :style="formData.bg_image ? { backgroundImage: `url(${formData.bg_image})` } : { background: formData.color }">
                  <div v-if="formData.bg_image" class="blur-hint">背景图片</div>
                  <span v-else class="bg-placeholder">背景图片</span>
                </div>
                <div class="upload-actions">
                  <button type="button" class="btn btn-sm btn-secondary" @click="triggerBgUpload">上传背景</button>
                  <input type="file" ref="bgInput" accept="image/*" style="display: none" @change="handleBgUpload" />
                  <button v-if="formData.bg_image" type="button" class="btn btn-sm btn-danger" @click="formData.bg_image = ''">移除</button>
                </div>
              </div>
            </div>
          </div>
          <div class="form-field">
            <label>团队名称</label>
            <input v-model="formData.name" type="text" placeholder="请输入团队名称" />
          </div>
          <div class="form-field">
            <label>团队描述</label>
            <textarea v-model="formData.description" rows="3" placeholder="请输入团队描述（可选）"></textarea>
          </div>
          <div class="form-field">
            <label>团队颜色</label>
            <div class="color-input-wrapper">
              <input v-model="formData.color" type="color" />
              <input v-model="formData.color" type="text" placeholder="#6366f1" />
            </div>
          </div>
          <div v-if="isEditing && isTeamLeaderOrAdmin(currentTeam)" class="form-field">
            <label>组长</label>
            <select v-model="formData.leader_id">
              <option value="">请选择组长</option>
              <option v-for="member in getTeamMembers(currentTeam.id)" :key="member.id" :value="member.id">
                {{ member.real_name }} (@{{ member.username }})
              </option>
            </select>
          </div>
          <div v-if="isEditing && isTeamLeaderOrAdmin(currentTeam)" class="form-field">
            <label>公告</label>
            <textarea v-model="formData.announcement" rows="2" placeholder="输入团队公告（可选）"></textarea>
          </div>
          <div v-if="isEditing && userStore.role === 'admin'" class="form-field">
            <label>可见性设置</label>
            <div class="toggle-field">
              <label class="toggle-label">
                <input type="checkbox" v-model="formData.is_public" />
                <span class="toggle-text">{{ formData.is_public ? '公开' : '私密' }}</span>
              </label>
              <p class="field-hint">私密团队不会在公共论坛页面显示</p>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeFormPanel">取消</button>
          <button class="btn btn-primary" @click="saveTeam" :disabled="saving">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="announcementDialogVisible" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="announcementDialogVisible = false">
      <div class="modal announcement-modal" style="background:var(--bg-card);">
        <div class="modal-header">
          <h3>发布公告</h3>
          <button class="close-btn" @click="announcementDialogVisible = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-field">
            <label>公告内容</label>
            <textarea v-model="announcementContent" rows="4" placeholder="输入团队公告内容..."></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="announcementDialogVisible = false">取消</button>
          <button class="btn btn-primary" @click="handleAnnouncement">发布</button>
        </div>
      </div>
    </div>
  </div>

  <!-- 底部空白区域 -->
  <div class="page-bottom-spacer"></div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'
import { useConfirm } from '../utils/confirm'
import { getTeamsAPI, createTeamAPI, updateTeamAPI, deleteTeamAPI, setUserTeamAPI, updateMyTeamAPI, getPublicUsersAPI, uploadTeamLogoAPI, uploadTeamBgAPI } from '../api'

const router = useRouter()
const route = useRoute()
const { success, error: notifyError, warning } = useNotification()
const { confirm, confirmDanger } = useConfirm()

const userStore = useUserStore()
const hasTeam = computed(() => {
  const teamId = userStore.team_id || localStorage.getItem('team_id')
  return !!teamId && teamId !== 'null' && teamId !== 'undefined'
})
const hasPendingRequest = ref(false)
const pendingRequestInfo = ref(null)

const scrollY = ref(0)
const showButtons = ref(false)
const canScroll = ref(false)
const expanded = ref(false)
const isAtTop = ref(true)
let hoverTimeout = null

let scrollTimeout = null

const handleScroll = () => {
  if (scrollTimeout) clearTimeout(scrollTimeout)
  
  scrollTimeout = setTimeout(() => {
    scrollY.value = window.scrollY
    isAtTop.value = window.scrollY <= 100
    
    if (window.scrollY > 100) {
      showButtons.value = true
    } else {
      showButtons.value = false
      expanded.value = false
    }
  }, 10)
}

const checkCanScroll = () => {
  setTimeout(() => {
    canScroll.value = document.documentElement.scrollHeight > document.documentElement.clientHeight
    // 如果页面不能滚动，默认显示按钮
    if (!canScroll.value) {
      showButtons.value = true
    }
  }, 100)
}

const handleMouseEnter = () => {
  if (hoverTimeout) clearTimeout(hoverTimeout)
  showButtons.value = true
  // 添加 Header 高亮效果
  triggerHeaderEffect()
}

const handleMouseLeave = () => {
  // 延迟隐藏，给用户一点反应时间
  hoverTimeout = setTimeout(() => {
    if (!canScroll.value || scrollY.value > 100) {
      showButtons.value = true
    } else {
      showButtons.value = false
    }
    // 移除 Header 高亮效果
    removeHeaderGlow()
  }, 300)
}

const triggerHeaderEffect = () => {
  const header = document.querySelector('.page-header')
  if (header) {
    header.classList.add('header-glow')
    header.classList.add('header-shine')
    setTimeout(() => {
      header.classList.remove('header-shine')
    }, 600)
  }
}

const removeHeaderGlow = () => {
  const header = document.querySelector('.page-header')
  if (header) {
    header.classList.remove('header-glow')
  }
}

const toggleExpanded = () => {
  expanded.value = !expanded.value
  showButtons.value = expanded.value
}

const viewMode = ref('card')
const teams = ref([])
const users = ref([])
const loading = ref(false)
const detailPanelVisible = ref(false)
const formPanelVisible = ref(false)
const currentTeam = ref(null)
const isEditing = ref(false)
const saving = ref(false)
const formError = ref('')

const formData = ref({
  name: '',
  description: '',
  color: '#6366f1',
  logo: '',
  bg_image: '',
  is_public: true
})

const logoInput = ref(null)
const bgInput = ref(null)

const loadTeams = async () => {
  loading.value = true
  try {
    const res = await getTeamsAPI()
    if (res.code === 200) {
      teams.value = res.data || []
    }
  } catch (error) {
    notifyError('加载团队失败')
  } finally {
    loading.value = false
  }
}

const loadUsers = async () => {
  try {
    const res = await getPublicUsersAPI()
    if (res.code === 200) {
      users.value = res.data || []
    }
  } catch (error) {
    console.error('加载用户失败', error)
  }
}

const getTeamMembers = (teamId) => {
  return users.value.filter(u => u.team_id === teamId && u.status === 'active')
}

const openTeamDetail = (team) => {
  currentTeam.value = team
  detailPanelVisible.value = true
}

const closeDetailPanel = () => {
  detailPanelVisible.value = false
}

const goToTeamHome = (team) => {
  router.push({ name: 'TeamHome', params: { teamId: team.id } })
}

const goToMyTeam = () => {
  const teamId = userStore.team_id || localStorage.getItem('team_id')
  if (teamId) {
    router.push({ name: 'TeamHome', params: { teamId } })
  }
}

const formatTime = (timeStr) => {
  if (!timeStr) return '-'
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const loadPendingRequest = async () => {
  try {
    const res = await fetch('/api/team/my-request', {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    const data = await res.json()
    if (data.code === 200 && data.data) {
      const req = data.data
      if (req.status === 'pending') {
        hasPendingRequest.value = true
        pendingRequestInfo.value = req
      } else if (req.status === 'rejected') {
        hasPendingRequest.value = false
        pendingRequestInfo.value = req
        warning(`您的团队创建申请被拒绝: ${req.reject_reason || '无原因'}`)
      } else {
        hasPendingRequest.value = false
        pendingRequestInfo.value = null
      }
    } else {
      hasPendingRequest.value = false
    }
  } catch (e) {
    console.error('加载待审核申请失败:', e)
  }
}

const checkPendingStatus = async () => {
  await loadPendingRequest()
}

const goToUsers = () => {
  router.push('/users')
}

const goToChat = () => {
  router.push('/chat')
}

const openCreateTeam = () => {
  const teamId = userStore.team_id || localStorage.getItem('team_id')
  if (teamId && teamId !== 'null' && teamId !== 'undefined') {
    warning('您已加入其他团队，如需创建新团队请先退出当前团队')
    return
  }
  if (hasPendingRequest.value) {
    warning('您已有待审核的团队创建申请，请等待管理员审核')
    return
  }

  isEditing.value = false
  formData.value = {
    name: '',
    description: '',
    color: '#6366f1',
    logo: '',
    bg_image: '',
    is_public: true
  }
  formPanelVisible.value = true
}

const editTeam = (team) => {
  isEditing.value = true
  currentTeam.value = team
  formData.value = {
    id: team.id,
    name: team.name,
    description: team.description || '',
    color: team.color || '#6366f1',
    logo: team.logo || '',
    bg_image: team.bg_image || '',
    leader_id: team.leader_id || '',
    announcement: team.announcement || '',
    is_public: team.is_public !== false
  }
  detailPanelVisible.value = false
  formPanelVisible.value = true
}

const closeFormPanel = () => {
  formPanelVisible.value = false
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
  
  try {
    const base64Data = await new Promise((resolve) => {
      const reader = new FileReader()
      reader.onload = (e) => resolve(e.target.result)
      reader.readAsDataURL(file)
    })
    
    const res = await uploadTeamLogoAPI({ image: base64Data })
    if (res.code === 200) {
      formData.value.logo = res.data.url
      success('Logo上传成功')
    } else {
      notifyError(res.message || '上传失败')
    }
  } catch (error) {
    console.error('Logo upload error:', error)
    notifyError('上传失败')
  } finally {
    e.target.value = ''
  }
}

const handleBgUpload = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  
  try {
    const base64Data = await new Promise((resolve) => {
      const reader = new FileReader()
      reader.onload = (e) => resolve(e.target.result)
      reader.readAsDataURL(file)
    })
    
    const res = await uploadTeamBgAPI({ image: base64Data })
    if (res.code === 200) {
      formData.value.bg_image = res.data.url
      success('背景上传成功')
    } else {
      notifyError(res.message || '上传失败')
    }
  } catch (error) {
    console.error('Background upload error:', error)
    notifyError('上传失败')
  } finally {
    e.target.value = ''
  }
}

const saveTeam = async () => {
  formError.value = ''
  if (!formData.value.name?.trim()) {
    formError.value = '请输入团队名称'
    return
  }

  saving.value = true
  try {
    let res
    if (isEditing.value) {
      res = await updateTeamAPI(formData.value.id, formData.value)
    } else {
      res = await createTeamAPI(formData.value)
    }

    if (res.code === 200) {
      if (isEditing.value) {
        success('团队更新成功')
      } else if (userStore.role === 'admin') {
        success('团队创建成功')
      } else {
        success('团队创建申请已提交，请等待管理员审核')
      }
      closeFormPanel()
      loadTeams()
    } else {
      formError.value = res.message || '操作失败'
    }
  } catch (error) {
    formError.value = '操作失败，请重试'
  } finally {
    saving.value = false
  }
}

const joinTeam = async (team) => {
  try {
    const res = await updateMyTeamAPI(team.id)
    if (res.code === 200) {
      success(`已成功加入【${team.name}】`)
      userStore.updateTeam(team.id, team.name, team.color)
      loadUsers()
      closeDetailPanel()
    } else {
      notifyError(res.message || '加入失败')
    }
  } catch (error) {
    notifyError('加入失败')
  }
}

const leaveTeam = async () => {
  if (!await confirmDanger('确定退出当前团队？')) return

  try {
    const res = await updateMyTeamAPI(null)
    if (res.code === 200) {
      success('已退出团队')
      userStore.updateTeam(null, null, null)
      loadUsers()
      closeDetailPanel()
    } else {
      notifyError(res.message || '退出失败')
    }
  } catch (error) {
    notifyError('退出失败')
  }
}

const announcementDialogVisible = ref(false)
const announcementContent = ref('')

const isTeamLeaderOrAdmin = (team) => {
  if (userStore.role === 'admin') return true
  if (team && String(team.leader_id) === String(userStore.id)) return true
  return false
}

const showAnnouncementDialog = () => {
  announcementContent.value = currentTeam.value?.announcement || ''
  announcementDialogVisible.value = true
}

const handleAnnouncement = async () => {
  if (!currentTeam.value) return
  if (announcementContent.value && !await confirm('确定发布此公告？')) return
  try {
    const res = await updateTeamAPI(currentTeam.value.id, {
      announcement: announcementContent.value
    })
    if (res.code === 200) {
      success(announcementContent.value ? '公告发布成功' : '公告已清除')
      currentTeam.value.announcement = announcementContent.value
      announcementDialogVisible.value = false
    } else {
      notifyError(res.message || '发布失败')
    }
  } catch (error) {
    notifyError('发布失败')
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

onMounted(async () => {
  await loadTeams()
  await loadUsers()
  await loadPendingRequest()
  window.addEventListener('scroll', handleScroll)
  checkCanScroll()
  const referer = document.referrer
  const isFromTeamHome = referer.includes('/team/') || route.query.from === 'teamHome'
  const storedTeamId = localStorage.getItem('team_id')
  if (storedTeamId && !isFromTeamHome) {
    const myTeam = teams.value.find(t => String(t.id) === String(storedTeamId))
    if (myTeam) {
      router.push({ name: 'TeamHome', params: { teamId: myTeam.id } })
    }
  }
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
  if (scrollTimeout) clearTimeout(scrollTimeout)
})
</script>

<style scoped>
.team-page {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-bottom-spacer {
  height: 400px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 24px;
  background: linear-gradient(135deg, rgba(255,255,255,0.9), rgba(255,255,255,0.5));
  backdrop-filter: blur(20px);
  border-radius: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255,255,255,0.3);
  position: relative;
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.page-header::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.08) 0%, transparent 70%);
  animation: rotateBg 20s linear infinite;
  pointer-events: none;
  transition: all 0.4s ease;
}

.page-header::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.4), transparent);
  transition: left 0.6s ease;
  pointer-events: none;
}

.page-header.header-glow::before {
  background: radial-gradient(circle, rgba(99, 102, 241, 0.15) 0%, rgba(255, 107, 53, 0.1) 50%, transparent 70%);
  animation: rotateBg 10s linear infinite;
}

.page-header.header-glow {
  box-shadow: 0 12px 48px rgba(99, 102, 241, 0.2), 0 0 80px rgba(255, 107, 53, 0.1);
  border-color: rgba(99, 102, 241, 0.3);
}

.page-header.header-shine::after {
  left: 100%;
}

@keyframes rotateBg {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.header-left h1 {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
  transition: all 0.3s ease;
}

.header-left p {
  color: var(--text-secondary);
  font-size: 14px;
  transition: all 0.3s ease;
}

.header-right {
  display: flex;
  align-items: center;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-buttons:has(.public-forum-btn:hover) .my-team-btn,
.action-buttons:has(.my-team-btn:hover) .public-forum-btn {
  transform: scale(0.95);
  opacity: 0.7;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 200px;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  color: #333;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.action-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.btn-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  flex-shrink: 0;
  color: white;
}

.btn-icon svg {
  width: 20px;
  height: 20px;
}

.btn-text {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  text-align: left;
}

.btn-title {
  font-size: 14px;
  font-weight: 600;
  line-height: 1.2;
}

.btn-subtitle {
  font-size: 12px;
  opacity: 0.7;
  margin-top: 2px;
}

.action-btn.chat-btn {
  background: rgba(102, 126, 234, 0.05);
  border: 1px solid rgba(102, 126, 234, 0.2);
}

.action-btn.chat-btn .btn-icon {
  background: #667eea;
}

.action-btn.public-forum-btn {
  background: rgba(99, 102, 241, 0.05);
  border: 1px solid rgba(99, 102, 241, 0.2);
}

.action-btn.public-forum-btn .btn-icon {
  background: #6366f1;
}

.action-btn.my-team-btn {
  background: rgba(59, 130, 246, 0.05);
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.action-btn.my-team-btn .btn-icon {
  background: #3b82f6;
}

.action-btn.my-team-btn {
  background: rgba(255, 107, 53, 0.05);
  border: 1px solid rgba(255, 107, 53, 0.2);
}

.action-btn.my-team-btn .btn-icon {
  background: #ff6b35;
}

.fixed-action-buttons {
  position: fixed;
  right: -180px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 12px;
  z-index: 100;
  transition: right 0.3s ease;
}

.fixed-action-buttons.show {
  right: 24px;
}

.fixed-action-buttons.expanded {
  right: 24px;
}

.floating-trigger {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: linear-gradient(135deg, #6366f1, #ff6b35);
  border: none;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  transition: all 0.3s ease;
}

.floating-trigger svg {
  width: 24px;
  height: 24px;
  color: white;
}

.floating-trigger:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
}

.floating-trigger.active {
  background: linear-gradient(135deg, #ff6b35, #6366f1);
  transform: rotate(90deg);
}

.fixed-action-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  width: 180px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.fixed-action-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.2), transparent);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.fixed-action-btn:hover::before {
  opacity: 1;
}

.fixed-action-btn:hover {
  transform: translateX(-4px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
}

.fixed-action-btn .btn-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255,255,255,0.2);
  border-radius: 10px;
  flex-shrink: 0;
}

.fixed-action-btn .btn-icon svg {
  width: 24px;
  height: 24px;
}

.fixed-action-btn .btn-text {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  text-align: left;
  opacity: 1;
  transition: opacity 0.2s ease;
}

.fixed-action-btn .btn-title {
  font-size: 14px;
  font-weight: 700;
  line-height: 1.2;
}

.fixed-action-btn .btn-subtitle {
  font-size: 11px;
  opacity: 0.8;
  margin-top: 2px;
}

.fixed-action-btn.public-forum-btn {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
}

.fixed-action-btn.my-team-btn {
  background: linear-gradient(135deg, #ff6b35, #f7931e);
  color: white;
}

.right-hover-trigger {
  position: fixed;
  right: 0;
  top: 0;
  width: 30px;
  height: 100vh;
  z-index: 99;
  cursor: pointer;
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
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.toggle-btn svg {
  width: 18px;
  height: 18px;
}

.toggle-btn.active {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.team-actions {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.team-actions .has-team-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.team-actions .info-text {
  color: var(--text-secondary);
  font-size: 14px;
}

.pending-request-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  background: linear-gradient(135deg, #fef3c7, #fde68a);
  border-radius: 12px;
  border: 1px solid #f59e0b;
}

.pending-icon {
  font-size: 32px;
  animation: pulse 2s ease-in-out infinite;
}

.pending-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.pending-title {
  font-weight: 600;
  color: #92400e;
  font-size: 15px;
}

.pending-desc {
  color: #b45309;
  font-size: 13px;
}

.pending-time {
  color: #d97706;
  font-size: 12px;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

.loading {
  text-align: center;
  padding: 80px 20px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border-color);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.team-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.team-card {
  position: relative;
  background: transparent;
  border: none;
  border-radius: 0;
  overflow: visible;
  cursor: pointer;
  transition: z-index 0.4s ease;
  min-height: 280px;
}

.team-card:hover {
  z-index: 5;
}

.team-card .card-inner {
  position: relative;
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.95) 0%, rgba(255, 255, 255, 0.92) 100%);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 20px;
  overflow: hidden;
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  min-height: 280px;
}

.dark .team-card .card-inner {
  background: linear-gradient(145deg, rgba(30, 41, 59, 0.96) 0%, rgba(15, 23, 42, 0.98) 100%);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.team-card .card-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 60%;
  height: 100%;
  background: radial-gradient(circle at 80% 20%, rgba(59, 130, 246, 0.18) 0%, transparent 55%);
  opacity: 0.7;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .team-card .card-bg {
  background: radial-gradient(circle at 80% 20%, rgba(59, 130, 246, 0.25) 0%, transparent 55%);
}

.team-card:hover .card-bg {
  opacity: 0.35;
  width: 45%;
}

.team-card .card-avatar-bg {
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

.team-card .card-avatar-bg img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.team-card:hover .card-avatar-bg {
  opacity: 1;
  transform: translateX(0);
}

.team-card-content {
  position: relative;
  padding: 24px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.team-logo {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 600;
  overflow: hidden;
  margin-bottom: 16px;
  color: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.team-logo:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

.team-logo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.team-logo > span {
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.team-name {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 8px;
}

.team-desc {
  font-size: 14px;
  margin-bottom: 16px;
  flex: 1;
  color: var(--text-secondary, #64748b);
}

.team-leader {
  font-size: 13px;
  margin-bottom: 4px;
  color: var(--text-secondary, #64748b);
}

.team-stats {
  display: flex;
  gap: 16px;
  font-size: 13px;
  margin-bottom: 16px;
  color: var(--text-secondary, #64748b);
}

.team-stats .stat {
  display: flex;
  align-items: center;
  gap: 4px;
}

.team-stats .stat svg {
  width: 14px;
  height: 14px;
}

.team-members-preview {
  display: flex;
  margin-left: -8px;
}

.member-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 500;
  border: 2px solid white;
  margin-left: -8px;
  overflow: hidden;
}

.member-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.member-avatar.more {
  background: rgba(255, 255, 255, 0.3);
  font-size: 10px;
}

.team-list {
  background: var(--bg-card);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--border-color);
}

.list-header {
  display: flex;
  padding: 12px 20px;
  background: var(--bg-dark);
  font-weight: 500;
  font-size: 13px;
  color: var(--text-secondary);
}

.list-row {
  display: flex;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  transition: background 0.2s;
}

.list-row:last-child {
  border-bottom: none;
}

.list-row:hover {
  background: var(--bg-hover);
}

.list-cell {
  display: flex;
  align-items: center;
}

.logo-cell { width: 60px; }
.name-cell { width: 120px; }
.leader-cell { width: 100px; color: var(--text-secondary); }
.desc-cell { flex: 1; color: var(--text-secondary); }
.members-cell { width: 80px; justify-content: center; }
.action-cell { width: 100px; justify-content: center; }

.team-logo-sm {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  color: white;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.team-logo-sm:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.team-logo-sm img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.team-logo-sm > span {
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.joined-tag {
  color: var(--success);
  font-size: 13px;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: var(--text-muted);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.empty-state svg {
  width: 80px;
  height: 80px;
  opacity: 0.5;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  border-radius: var(--radius-md);
  border: none;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn svg {
  width: 18px;
  height: 18px;
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  color: white;
  box-shadow: var(--shadow-sm);
}

.btn-primary:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.btn-secondary {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
  border: 1px solid #6366f1;
  padding: 10px 20px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.btn-secondary:hover {
  background: linear-gradient(135deg, #8b5cf6, #a78bfa);
  border-color: #8b5cf6;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(99, 102, 241, 0.4);
}

.btn-primary {
  background: linear-gradient(135deg, #ff6b35, #f7931e);
  color: white;
  border: 1px solid #ff6b35;
  padding: 10px 24px;
  font-weight: 600;
  letter-spacing: 0.5px;
  box-shadow: 0 4px 12px rgba(255, 107, 53, 0.3);
}

.btn-primary:hover {
  background: linear-gradient(135deg, #f7931e, #ff8c5a);
  border-color: #f7931e;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(255, 107, 53, 0.4);
}

.btn-danger {
  background: linear-gradient(135deg, var(--danger), #dc2626);
  color: white;
}

.btn-danger:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.btn-lg {
  padding: 12px 32px;
  font-size: 16px;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 12px;
}

.slide-panel {
  position: fixed;
  inset: 0;
  z-index: 1000;
  visibility: hidden;
}

.slide-panel.open {
  visibility: visible;
}

.panel-backdrop {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
}

.panel-content {
  position: absolute;
  top: 0;
  right: 0;
  width: 480px;
  max-width: 90vw;
  height: 100%;
  background-size: cover;
  background-position: center;
  transition: transform 0.3s ease;
  transform: translateX(100%);
}

.slide-panel.open .panel-content {
  transform: translateX(0);
}

.panel-blur-bg {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  filter: blur(40px);
  transform: scale(1.2);
}

.panel-inner {
  position: relative;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(30px);
  color: white;
  display: flex;
  flex-direction: column;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.close-btn {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 50%;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: white;
  transition: background 0.2s;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.close-btn svg {
  width: 20px;
  height: 20px;
}

.edit-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.15);
  border: none;
  border-radius: 8px;
  color: white;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.2s;
}

.edit-btn:hover {
  background: rgba(255, 255, 255, 0.25);
}

.edit-btn svg {
  width: 16px;
  height: 16px;
}

.panel-body {
  flex: 1;
  overflow-y: auto;
  padding: 32px 24px;
}

.team-logo-lg {
  width: 80px;
  height: 80px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: 600;
  margin: 0 auto 20px;
  overflow: hidden;
}

.team-logo-lg img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.team-title {
  font-size: 24px;
  font-weight: 600;
  text-align: center;
  margin-bottom: 8px;
}

.team-description {
  text-align: center;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
  margin-bottom: 24px;
}

.team-info-stats {
  display: flex;
  justify-content: center;
  gap: 40px;
  margin-bottom: 24px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-bottom: 32px;
}

.members-section h3 {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.members-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.member-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
}

.member-avatar-lg {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 500;
  overflow: hidden;
  flex-shrink: 0;
}

.member-avatar-lg img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.member-info {
  flex: 1;
}

.member-name {
  display: block;
  font-weight: 500;
  margin-bottom: 2px;
}

.member-username {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}

.no-members {
  text-align: center;
  padding: 24px;
  color: rgba(255, 255, 255, 0.5);
}

.modal-overlay {
  position: fixed;
  inset: 0;
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
  width: 500px;
  max-width: 90vw;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2), 0 0 0 1px var(--border-color);
}

html.dark .modal {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(255, 255, 255, 0.05);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
}

.modal-header .close-btn {
  background: var(--bg-dark);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.modal-header .close-btn:hover {
  background: var(--bg-hover);
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  color: #ef4444;
  font-size: 14px;
  margin-bottom: 16px;
}

.error-message svg {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.info-message {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 8px;
  color: #3b82f6;
  font-size: 14px;
  margin-bottom: 16px;
}

.info-message svg {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.announcement-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: rgba(249, 115, 22, 0.1);
  border: 1px solid rgba(249, 115, 22, 0.3);
  border-radius: 8px;
  color: var(--primary);
  font-size: 14px;
  margin-bottom: 16px;
}

.announcement-bar svg {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.leader-badge {
  display: inline-block;
  padding: 2px 6px;
  background: var(--primary);
  color: white;
  font-size: 10px;
  border-radius: 4px;
  margin-left: 6px;
}

.member-item.is-leader {
  background: rgba(249, 115, 22, 0.05);
  border-radius: 8px;
  padding: 8px 12px;
}

.announcement-modal {
  width: 400px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
}

.form-row {
  display: flex;
  gap: 16px;
}

.logo-row .form-field {
  flex: 1;
}

.form-field {
  margin-bottom: 20px;
}

.form-field label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 8px;
  color: var(--text-primary);
}

.form-field input,
.form-field textarea,
.form-field select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  background: var(--input-bg);
  color: var(--text-primary);
}

.form-field textarea {
  resize: vertical;
}

.image-upload-preview {
  display: flex;
  gap: 16px;
  align-items: flex-start;
}

.preview-box {
  width: 80px;
  height: 80px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 600;
  color: white;
  background-size: cover;
  background-position: center;
  overflow: hidden;
  position: relative;
}

.preview-box img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.bg-preview {
  width: 120px;
  height: 80px;
}

.blur-hint {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.4);
  font-size: 12px;
}

.bg-placeholder {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
}

.upload-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.color-input-wrapper {
  display: flex;
  gap: 8px;
}

.color-input-wrapper input[type="color"] {
  width: 48px;
  height: 40px;
  padding: 4px;
  border-radius: 8px;
  cursor: pointer;
}

.color-input-wrapper input[type="text"] {
  flex: 1;
}

.toggle-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.toggle-label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.toggle-label input[type="checkbox"] {
  width: 20px;
  height: 20px;
  cursor: pointer;
  accent-color: var(--primary);
}

.toggle-text {
  font-size: 14px;
  color: var(--text-primary);
  font-weight: 500;
}

.field-hint {
  font-size: 12px;
  color: var(--text-muted);
  margin: 0;
}

.team-detail-modal {
  width: 560px;
  max-width: 95vw;
  max-height: 85vh;
}

.team-detail-header {
  background-size: cover;
  background-position: center;
  padding: 40px 24px 24px;
  margin: -24px -24px 24px -24px;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: 200px;
  justify-content: flex-end;
}

.team-detail-header::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(to bottom, rgba(0,0,0,0.2) 0%, rgba(0,0,0,0.75) 100%);
}

.team-detail-header > * {
  position: relative;
  z-index: 1;
}

.team-logo-xl {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: 600;
  margin-bottom: 16px;
  overflow: hidden;
  border: 3px solid white;
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.team-logo-xl:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 16px rgba(0,0,0,0.4);
}

.team-logo-xl img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.team-logo-xl:hover img {
  transform: scale(1.08);
}

.team-logo-xl > span {
  color: white;
  text-shadow: 0 2px 4px rgba(0,0,0,0.3);
}

.team-detail-header .team-title {
  color: white;
  text-shadow: 0 2px 4px rgba(0,0,0,0.4);
  font-size: 22px;
  margin-bottom: 4px;
}

.team-detail-header .team-description {
  color: rgba(255,255,255,0.95);
  text-align: center;
  font-size: 14px;
  line-height: 1.5;
}

.team-detail-modal .modal-body {
  padding: 20px 24px 24px;
}

.team-detail-modal .team-info-stats {
  display: flex;
  justify-content: center;
  gap: 32px;
  padding: 16px 20px;
  background: var(--bg-primary);
  border-radius: 12px;
  margin-bottom: 16px;
  border: 1px solid var(--border-color);
}

.team-detail-modal .stat-item {
  text-align: center;
}

.team-detail-modal .stat-value {
  display: block;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 2px;
}

.team-detail-modal .stat-label {
  font-size: 12px;
  color: var(--text-secondary);
}

.team-detail-modal .announcement-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: rgba(249, 115, 22, 0.08);
  border: 1px solid rgba(249, 115, 22, 0.2);
  border-radius: 10px;
  color: var(--primary);
  font-size: 14px;
  margin-bottom: 16px;
}

.team-detail-modal .announcement-bar svg {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.team-detail-modal .action-buttons {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.team-detail-modal .action-buttons .btn {
  padding: 10px 20px;
  font-size: 14px;
}

.team-detail-modal .action-buttons .btn svg {
  width: 16px;
  height: 16px;
}

.team-detail-modal .members-section {
  margin-top: 8px;
}

.team-detail-modal .members-section h3 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border-color);
}

.team-detail-modal .members-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.team-detail-modal .member-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  background: var(--bg-primary);
  border-radius: 10px;
  border: 1px solid var(--border-color);
  transition: background 0.2s;
}

.team-detail-modal .member-item:hover {
  background: var(--bg-hover);
}

.team-detail-modal .member-item.is-leader {
  border-color: var(--primary);
  background: rgba(249, 115, 22, 0.05);
}

.team-detail-modal .member-avatar-lg {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 500;
  overflow: hidden;
  flex-shrink: 0;
  border: 2px solid var(--border-color);
}

.team-detail-modal .member-avatar-lg img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.team-detail-modal .member-info {
  flex: 1;
  min-width: 0;
}

.team-detail-modal .member-name {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 500;
  font-size: 14px;
  color: var(--text-primary);
}

.team-detail-modal .member-username {
  font-size: 12px;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.team-detail-modal .leader-badge {
  display: inline-block;
  padding: 2px 6px;
  background: var(--primary);
  color: white;
  font-size: 10px;
  border-radius: 4px;
  font-weight: 500;
}

.team-detail-modal .no-members {
  text-align: center;
  padding: 24px;
  color: var(--text-muted);
  font-size: 14px;
}

@media (max-width: 768px) {
  .team-page {
    padding: 12px;
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .header-left h1 {
    font-size: 24px;
  }

  .header-left p {
    font-size: 13px;
  }

  .header-right {
    justify-content: center;
  }

  .action-buttons {
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: center;
    gap: 8px;
  }

  .action-btn {
    min-width: auto;
    padding: 10px 16px;
  }

  .btn-icon {
    width: 36px;
    height: 36px;
  }

  .btn-icon svg {
    width: 20px;
    height: 20px;
  }

  .btn-subtitle {
    display: none;
  }

  .btn-title {
    font-size: 14px;
  }

  .fixed-action-buttons {
    right: -60px;
    gap: 8px;
  }

  .fixed-action-buttons.show {
    right: 12px;
  }

  .fixed-action-buttons.expanded {
    right: 12px;
  }

  .floating-trigger {
    right: 12px;
    bottom: 12px;
    width: 48px;
    height: 48px;
  }

  .floating-trigger svg {
    width: 20px;
    height: 20px;
  }

  .fixed-action-btn {
    width: 48px;
    height: 48px;
    padding: 0;
    justify-content: center;
  }

  .fixed-action-btn .btn-icon {
    width: 36px;
    height: 36px;
    margin: 0;
  }

  .fixed-action-btn .btn-icon svg {
    width: 20px;
    height: 20px;
  }

  .fixed-action-btn .btn-text {
    display: none;
  }

  .team-grid {
    grid-template-columns: 1fr;
  }

  .form-row {
    flex-direction: column;
  }

  .panel-content {
    width: 100%;
    max-width: 100%;
  }

  .form-field input,
  .form-field textarea,
  .form-field select {
    padding: 12px;
    font-size: 16px;
  }

  .modal {
    width: 100%;
    max-width: 100%;
    border-radius: 16px 16px 0 0;
    max-height: 85vh;
  }

  .modal-header {
    padding: 16px;
  }

  .modal-body {
    padding: 16px;
  }

  .modal-footer {
    padding: 12px 16px;
  }
}

@media (max-width: 480px) {
  .team-page {
    padding: 8px;
  }

  .header-left h1 {
    font-size: 20px;
  }

  .view-toggle {
    width: 100%;
    justify-content: center;
  }

  .team-actions {
    width: 100%;
  }

  .team-actions .btn {
    width: 100%;
    justify-content: center;
  }

  .team-card .card-inner {
    padding: 16px;
  }

  .team-card .card-avatar {
    width: 60px;
    height: 60px;
  }

  .card-info h3 {
    font-size: 16px;
  }

  .card-info p {
    font-size: 13px;
  }

  .member-count {
    font-size: 12px;
  }

  .slide-panel .panel-content {
    padding: 16px;
  }

  .form-section {
    padding: 16px;
  }

  .action-buttons {
    flex-direction: column;
  }

  .action-buttons .btn {
    width: 100%;
    justify-content: center;
  }
}
</style>