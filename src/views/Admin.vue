<template>
  <div class="admin">
    <div class="page-header">
      <h1>管理面板</h1>
      <p>系统管理与数据维护</p>
    </div>

    <div class="admin-tabs">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        class="tab-btn"
        :class="{ active: activeTab === tab.key }"
        @click="activeTab = tab.key"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" v-html="tab.icon"></svg>
        {{ tab.name }}
      </button>
    </div>

    <div class="tab-content">
      <div v-if="activeTab === 'category'" class="tab-panel">
        <div class="panel-header">
          <h3>分类管理</h3>
          <button class="add-btn" @click="showAddModal('category')">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            添加分类
          </button>
        </div>
        <div class="data-table">
          <div class="table-header">
            <span>ID</span>
            <span>图标</span>
            <span>分类名称</span>
            <span>颜色</span>
            <span>创建时间</span>
            <span>操作</span>
          </div>
          <div v-for="item in categories" :key="item.id" class="table-row">
            <span>{{ item.id }}</span>
            <span class="category-icon-preview">{{ item.icon || '🏷️' }}</span>
            <span>{{ item.name }}</span>
            <span>
              <span class="color-preview" :style="{ backgroundColor: item.color || '#6b7280' }"></span>
              {{ item.color || '#6b7280' }}
            </span>
            <span>{{ formatDate(item.create_time) }}</span>
            <span>
              <button class="action-btn" @click="editItem('category', item)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
              </button>
              <button class="action-btn danger" @click="deleteItem('category', item)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3,6 5,6 21,6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
              </button>
            </span>
          </div>
        </div>
      </div>

      <div v-if="activeTab === 'product'" class="tab-panel">
        <div class="panel-header">
          <h3>产品管理</h3>
          <button class="add-btn" @click="showAddModal('product')">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            添加产品
          </button>
        </div>
        <div class="data-table">
          <div class="table-header">
            <span>ID</span>
            <span>产品名称</span>
            <span>创建时间</span>
            <span>操作</span>
          </div>
          <div v-for="item in products" :key="item.id" class="table-row">
            <span>{{ item.id }}</span>
            <span>{{ item.name }}</span>
            <span>{{ formatDate(item.create_time) }}</span>
            <span>
              <button class="action-btn" @click="editItem('product', item)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
              </button>
              <button class="action-btn danger" @click="deleteItem('product', item)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3,6 5,6 21,6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
              </button>
            </span>
          </div>
        </div>
      </div>

      <div v-if="activeTab === 'team'" class="tab-panel">
        <div class="panel-header">
          <h3>团队管理</h3>
        </div>
        <div v-if="!userStore.team_id && userStore.role !== 'admin'" class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
          </svg>
          <p>您还未加入任何团队</p>
          <p class="empty-hint">请先在"团队"面板加入或创建团队</p>
        </div>
        <div v-if="loading" class="skeleton-wrap">
          <SkeletonLoader variant="table" :columns="4" :rows="5" :is-dark="document.documentElement.classList.contains('dark')" />
        </div>
        <div v-else>
          <div class="team-table" v-if="teams.length > 0">
            <div class="table-header">
              <div class="col-logo">Logo</div>
              <div class="col-name">团队名称</div>
              <div class="col-leader">组长</div>
              <div class="col-members">成员数</div>
              <div class="col-date">创建时间</div>
              <div class="col-actions">操作</div>
            </div>
            <div v-for="team in teams" :key="team.id" class="table-row">
              <div class="col-logo">
                <div class="team-logo-small" :style="{ backgroundColor: team.color }">
                  <img v-if="team.logo" :src="team.logo" :alt="team.name" />
                  <span v-else>{{ team.name.charAt(0) }}</span>
                </div>
              </div>
              <div class="col-name" @click="openTeamDetail(team)">
                <div class="team-name-cell clickable">{{ team.name }}</div>
                <div class="team-desc-cell">{{ team.description || '暂无描述' }}</div>
              </div>
              <div class="col-leader">{{ team.leader_name || '未知' }}</div>
              <div class="col-members">{{ getTeamMembers(team.id).length }} 人</div>
              <div class="col-date">{{ formatDate(team.create_time) }}</div>
              <div class="col-actions" @click.stop>
                <button class="action-btn small" @click="goToTeamHome(team)" title="进入">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
                    <polyline points="9 22 9 12 15 12 15 22"/>
                  </svg>
                </button>
                <button class="action-btn small" @click="editTeam(team)" title="编辑">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                  </svg>
                </button>
                <button class="action-btn small danger" @click="deleteTeam(team)" title="删除">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3,6 5,6 21,6"/>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>
          <div v-if="teams.length === 0" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
            </svg>
            <p>暂无团队</p>
          </div>
        </div>
      </div>

      <div v-if="teamDetailVisible" class="slide-panel" :class="{ open: teamDetailVisible }">
        <div class="panel-backdrop" @click="closeTeamDetail"></div>
        <div class="panel-content">
          <div class="panel-header-bar">
            <h3>团队详情</h3>
            <button class="close-btn" @click="closeTeamDetail">×</button>
          </div>
          <div v-if="selectedTeam" class="team-detail-content">
            <div class="team-detail-bg" :style="selectedTeam.bg_image ? { backgroundImage: `url(${selectedTeam.bg_image})` } : {}">
              <div class="team-bg-overlay"></div>
            </div>
            <div class="team-detail-info">
              <div class="team-detail-logo" :style="{ backgroundColor: selectedTeam.color }">
                <img v-if="selectedTeam.logo" :src="selectedTeam.logo" :alt="selectedTeam.name" />
                <span v-else>{{ selectedTeam.name.charAt(0) }}</span>
              </div>
              <h2 class="team-detail-name">{{ selectedTeam.name }}</h2>
              <p class="team-detail-desc">{{ selectedTeam.description || '暂无描述' }}</p>
              <div class="team-detail-meta">
                <span>创建于 {{ formatDate(selectedTeam.create_time) }}</span>
                <span v-if="selectedTeam.announcement">公告: {{ selectedTeam.announcement }}</span>
              </div>
            </div>
            <div class="team-detail-section">
              <h4>成员列表 ({{ getTeamMembers(selectedTeam.id).length }})</h4>
              <div class="team-detail-members">
                <div v-for="user in getTeamMembers(selectedTeam.id)" :key="user.id" class="detail-member-chip">
                  <span class="detail-member-avatar" :style="{ backgroundColor: selectedTeam.color }">
                    {{ user.real_name?.charAt(0) || 'U' }}
                  </span>
                  <span class="detail-member-name">{{ user.real_name }}</span>
                  <span class="detail-member-username">@{{ user.username }}</span>
                </div>
                <div v-if="getTeamMembers(selectedTeam.id).length === 0" class="no-members">暂无成员</div>
              </div>
            </div>
            <div class="team-detail-actions">
              <button v-if="isTeamMember(selectedTeam)" class="primary-btn" @click="goToTeamHome(selectedTeam)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
                  <polyline points="9 22 9 12 15 12 15 22"/>
                </svg>
                进入团队主页
              </button>
              <button v-if="userStore.username === 'admin'" class="danger-btn" @click="handlePermanentDeleteTeam(selectedTeam)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3,6 5,6 21,6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
                永久删除团队
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="formPanelVisible" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="closeFormPanel">
        <div class="modal team-form-modal" style="background:var(--bg-card);">
          <div class="modal-header">
            <h3>{{ isEdit ? '编辑团队' : '新建团队' }}</h3>
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
            <div v-if="isEdit && currentTeam && isTeamLeaderOrAdmin(currentTeam)" class="form-field">
              <label>组长</label>
              <select v-model="formData.leader_id">
                <option value="">请选择组长</option>
                <option v-for="member in getTeamMembers(currentTeam.id)" :key="member.id" :value="member.id">
                  {{ member.real_name }} (@{{ member.username }})
                </option>
              </select>
            </div>
            <div v-if="isEdit && currentTeam && isTeamLeaderOrAdmin(currentTeam)" class="form-field">
              <label>公告</label>
              <textarea v-model="formData.announcement" rows="2" placeholder="输入团队公告（可选）"></textarea>
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

      <div v-if="activeTab === 'users'" class="tab-panel">
        <div class="panel-header">
          <h3>用户管理</h3>
          <div class="panel-actions">
            <div class="registration-control">
              <button 
                class="reg-control-btn" 
                :class="{ active: registrationMode === 'open' }" 
                @click="setRegistrationMode('open')"
                title="限时开放注册，新用户无需审核直接通过"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <polyline points="12,6 12,12 16,14"/>
                </svg>
                开放注册
              </button>
              <button 
                class="reg-control-btn" 
                :class="{ active: registrationMode === 'closed' }" 
                @click="setRegistrationMode('closed')"
                title="限时禁止注册，新用户无法注册"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/>
                </svg>
                禁止注册
              </button>
              <button 
                class="reg-control-btn" 
                :class="{ active: registrationMode === 'normal' }" 
                @click="setRegistrationMode('normal')"
                title="正常审核模式，新用户需要管理员审核"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M9 12l2 2 4-4"/>
                  <path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                正常审核
              </button>
            </div>
            <div class="view-toggle">
              <button class="toggle-btn" :class="{ active: !compactView }" @click="compactView = false" title="卡片视图">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="3" width="7" height="7" rx="1"/>
                  <rect x="14" y="3" width="7" height="7" rx="1"/>
                  <rect x="3" y="14" width="7" height="7" rx="1"/>
                  <rect x="14" y="14" width="7" height="7" rx="1"/>
                </svg>
              </button>
              <button class="toggle-btn" :class="{ active: compactView }" @click="compactView = true" title="紧凑列表">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="3" y1="6" x2="21" y2="6"/>
                  <line x1="3" y1="12" x2="21" y2="12"/>
                  <line x1="3" y1="18" x2="21" y2="18"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
        
        <div v-if="selectedUsers.length > 0" class="batch-actions">
          <span class="selected-count">已选择 {{ selectedUsers.length }} 个用户</span>
          <button class="batch-btn approve" @click="handleBatchApprove">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20,6 9,17 4,12"/>
            </svg>
            批量批准
          </button>
          <button class="batch-btn reject" @click="handleBatchReject">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
            批量拒绝
          </button>
          <button class="batch-btn clear" @click="selectedUsers = []">
            取消选择
          </button>
        </div>
        
        <div v-if="loading" class="skeleton-wrap">
          <SkeletonLoader variant="table" :columns="5" :rows="6" :is-dark="document.documentElement.classList.contains('dark')" />
        </div>
        <div v-else :class="compactView ? 'user-table' : 'user-cards'">
          <template v-if="compactView">
            <div class="compact-table-header">
              <div class="col-check">
                <input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll" />
              </div>
              <div class="col-avatar">头像</div>
              <div class="col-name">姓名</div>
              <div class="col-username">用户名</div>
              <div class="col-role">角色</div>
              <div class="col-status">状态</div>
              <div class="col-team">团队</div>
              <div class="col-time">注册时间</div>
              <div class="col-actions">操作</div>
            </div>
            <div v-for="user in allUsers" :key="user.id" class="compact-table-row" :class="{ selected: selectedUsers.includes(user.id) }">
              <div class="col-check">
                <input type="checkbox" :checked="selectedUsers.includes(user.id)" @change="toggleSelectUser(user.id)" />
              </div>
              <div class="col-avatar">
                <div class="user-avatar-small" :class="{ admin: user.role === 'admin' }">
                  {{ user.real_name?.charAt(0) || 'U' }}
                </div>
              </div>
              <div class="col-name">{{ user.real_name }}</div>
              <div class="col-username">@{{ user.username }}</div>
              <div class="col-role">
                <span class="role-tag" :class="{ 'role-admin': user.role === 'admin' }">
                  {{ user.role === 'admin' ? '管理员' : '用户' }}
                </span>
              </div>
              <div class="col-status">
                <span class="status-tag" :class="'status-' + user.status">
                  {{ user.status === 'pending' ? '待审核' : user.status === 'active' ? '正常' : '已注销' }}
                </span>
              </div>
              <div class="col-team">
                <span v-if="user.team_name" class="team-tag-small" :style="{ backgroundColor: user.team_color + '20', color: user.team_color }">
                  {{ user.team_name }}
                </span>
                <span v-else class="team-tag-small no-team">未分配</span>
              </div>
              <div class="col-time">{{ formatDate(user.create_time) }}</div>
              <div class="col-actions">
                <template v-if="user.status === 'pending'">
                  <button class="action-btn-small approve" @click="handleApprove(user.id)" title="批准">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="20,6 9,17 4,12"/>
                    </svg>
                  </button>
                  <button class="action-btn-small reject" @click="handleReject(user.id)" title="拒绝">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <line x1="18" y1="6" x2="6" y2="18"/>
                      <line x1="6" y1="6" x2="18" y2="18"/>
                    </svg>
                  </button>
                </template>
                <template v-else-if="userStore.username === 'admin' && user.username !== 'admin'">
                  <button v-if="user.role !== 'admin'" class="action-btn-small admin" @click="handleSetAdmin(user.id)" title="设为管理员">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                    </svg>
                  </button>
                  <button v-else class="action-btn-small remove-admin" @click="handleRemoveAdmin(user.id)" title="取消管理员">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                    </svg>
                  </button>
                  <button v-if="user.status === 'active'" class="action-btn-small deactivate" @click="handleDeactivate(user.id)" title="注销">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <circle cx="12" cy="12" r="10"/>
                      <line x1="15" y1="9" x2="9" y2="15"/>
                      <line x1="9" y1="9" x2="15" y2="15"/>
                    </svg>
                  </button>
                  <button class="action-btn-small danger" @click="handleDeleteUser(user.id)" title="删除">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="3,6 5,6 21,6"/>
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                    </svg>
                  </button>
                </template>
              </div>
            </div>
          </template>
          <template v-else>
            <div v-for="user in allUsers" :key="user.id" class="user-card" :class="{ selected: selectedUsers.includes(user.id) }">
              <div class="user-select">
                <input type="checkbox" :checked="selectedUsers.includes(user.id)" @change="toggleSelectUser(user.id)" />
              </div>
              <div class="user-avatar" :class="{ admin: user.role === 'admin' }">
                {{ user.real_name?.charAt(0) || 'U' }}
              </div>
              <div class="user-info">
                <div class="user-name">{{ user.real_name }}</div>
                <div class="user-username">@{{ user.username }}</div>
                <div class="user-role" :class="{ 'role-admin': user.role === 'admin' }">
                  {{ user.role === 'admin' ? '管理员' : '用户' }}
                  <span v-if="user.username === 'admin'" class="super-admin-tag">超级</span>
                </div>
                <div class="user-status" :class="'status-' + user.status">{{ user.status === 'pending' ? '待审核' : user.status === 'active' ? '正常' : '已注销' }}</div>
                <div class="user-team">
                  <span v-if="user.team_name" class="team-tag" :style="{ backgroundColor: user.team_color + '20', color: user.team_color }">
                    {{ user.team_name }}
                  </span>
                  <span v-else class="team-tag no-team">未分配</span>
                </div>
                <div class="user-time">注册: {{ formatDate(user.create_time) }}</div>
              </div>
              <div class="user-actions">
                <template v-if="user.status === 'pending'">
                  <button class="action-btn approve" @click="handleApprove(user.id)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="20,6 9,17 4,12"/>
                    </svg>
                    批准
                  </button>
                  <button class="action-btn reject" @click="handleReject(user.id)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <line x1="18" y1="6" x2="6" y2="18"/>
                      <line x1="6" y1="6" x2="18" y2="18"/>
                    </svg>
                    拒绝
                  </button>
                </template>
                <template v-else-if="userStore.username === 'admin' && user.username !== 'admin'">
                  <button v-if="user.role !== 'admin'" class="action-btn admin-btn" @click="handleSetAdmin(user.id)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                    </svg>
                    设为管理员
                  </button>
                  <button v-else class="action-btn remove-admin-btn" @click="handleRemoveAdmin(user.id)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                    </svg>
                    取消管理员
                  </button>
                  <button v-if="user.status === 'active'" class="action-btn reject" @click="handleDeactivate(user.id)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <circle cx="12" cy="12" r="10"/>
                      <line x1="15" y1="9" x2="9" y2="15"/>
                      <line x1="9" y1="9" x2="15" y2="15"/>
                    </svg>
                    注销
                  </button>
                  <button v-if="userStore.username === 'admin' && user.username !== 'admin'" class="action-btn danger" @click="handleDeleteUser(user.id)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="3,6 5,6 21,6"/>
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                    </svg>
                    删除
                  </button>
                </template>
              </div>
            </div>
          </template>
          <div v-if="allUsers.length === 0" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
            </svg>
            <p>暂无用户</p>
          </div>
        </div>
      </div>

      <div v-if="activeTab === 'backup'" class="tab-panel">
        <div class="panel-header">
          <h3>数据库管理</h3>
        </div>
        <div class="backup-actions">
          <div class="backup-card" @click="handleCreateSnapshot">
            <div class="backup-icon green">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 5v14M5 12h14"/>
              </svg>
            </div>
            <div class="backup-info">
              <span class="backup-name">创建快照</span>
              <span class="backup-desc">立即生成当前数据库快照</span>
            </div>
          </div>
          <div class="backup-card" @click="handleBackup">
            <div class="backup-icon green">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/>
                <path d="M21 3v5h-5"/>
              </svg>
            </div>
            <div class="backup-info">
              <span class="backup-name">本地备份</span>
              <span class="backup-desc">导出数据库到本地文件</span>
            </div>
          </div>
          <div class="backup-card import" @click="triggerImport">
            <div class="backup-icon blue">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17,8 12,3 7,8"/>
                <line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
            </div>
            <div class="backup-info">
              <span class="backup-name">本地还原</span>
              <span class="backup-desc">从本地备份文件还原</span>
            </div>
          </div>
          <input type="file" ref="importFileInput" accept=".json" style="display: none" @change="handleImport">
          <div class="backup-card" @click="handlePurge">
            <div class="backup-icon orange">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3,6 5,6 21,6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                <line x1="10" y1="11" x2="10" y2="17"/>
                <line x1="14" y1="11" x2="14" y2="17"/>
              </svg>
            </div>
            <div class="backup-info">
              <span class="backup-name">清理垃圾数据</span>
              <span class="backup-desc">永久删除30天前的已删除记录</span>
            </div>
          </div>
        </div>

        <div class="snapshot-section">
          <div class="section-header">
            <h4>快照列表</h4>
            <span class="snapshot-count">共 {{ snapshots.length }} 个 / 最多 {{ maxSnapshots }} 个</span>
          </div>

          <div class="auto-backup-settings" v-if="userStore.username === 'admin'">
            <div class="setting-row">
              <span class="setting-label">自动快照间隔：</span>
              <select v-model="backupInterval" @change="handleUpdateBackupInterval" class="interval-select">
                <option :value="3600000">每 1 小时</option>
                <option :value="21600000">每 6 小时</option>
                <option :value="86400000">每 1 天</option>
                <option :value="259200000">每 3 天</option>
                <option :value="604800000">每 7 天</option>
                <option :value="2592000000">每 30 天</option>
              </select>
            </div>
          </div>

          <div class="snapshot-list" v-if="snapshots.length > 0">
            <div class="snapshot-item" v-for="snap in snapshots" :key="snap.filename">
              <div class="snapshot-info">
                <span class="snapshot-name">{{ snap.name }}</span>
                <span class="snapshot-meta">
                  {{ formatFileSize(snap.size) }} · {{ formatDate(snap.created) }}
                </span>
              </div>
              <div class="snapshot-actions">
                <button class="snapshot-btn restore" @click="handleRestoreSnapshot(snap)" title="还原此快照">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/>
                    <path d="M3 3v5h5"/>
                  </svg>
                </button>
                <button class="snapshot-btn download" @click="handleDownloadSnapshot(snap)" title="下载快照">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                    <polyline points="7,10 12,15 17,10"/>
                    <line x1="12" y1="15" x2="12" y2="3"/>
                  </svg>
                </button>
                <button class="snapshot-btn delete" @click="handleDeleteSnapshot(snap)" title="删除快照">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3,6 5,6 21,6"/>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>
          <div class="snapshot-empty" v-else>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14,2 14,8 20,8"/>
            </svg>
            <span>暂无快照</span>
          </div>
        </div>

        <div class="backup-warning">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
            <line x1="12" y1="9" x2="12" y2="13"/>
            <line x1="12" y1="17" x2="12.01" y2="17"/>
          </svg>
          <span>注意：还原操作会覆盖当前所有数据，且不可恢复！自动快照每{{ formatInterval(backupInterval) }}生成一次。</span>
        </div>

        <div class="danger-zone" v-if="userStore.username === 'admin'">
          <h4>危险操作区</h4>
          <button class="clear-btn" @click="handleClearData">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3,6 5,6 21,6"/>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              <line x1="10" y1="11" x2="10" y2="17"/>
              <line x1="14" y1="11" x2="14" y2="17"/>
            </svg>
            清空数据库
          </button>
        </div>
      </div>

      <div v-if="activeTab === 'export'" class="tab-panel">
        <div class="panel-header">
          <h3>数据导出</h3>
          <p class="panel-desc">自定义导出字段和格式，快速获取博主数据</p>
        </div>

        <div class="export-section">
          <div class="section-header">
            <h4>📋 字段选择</h4>
            <div class="quick-actions">
              <button class="quick-btn" @click="selectAllFields">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="9,11 12,14 22,4"/>
                  <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
                </svg>
                全选
              </button>
              <button class="quick-btn" @click="selectBasicFields">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                  <polyline points="14,2 14,8 20,8"/>
                </svg>
                基础字段
              </button>
              <button class="quick-btn danger" @click="clearFields">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"/>
                  <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
                清除
              </button>
            </div>
          </div>

          <div class="field-grid">
            <label
              v-for="field in exportFields"
              :key="field.key"
              class="field-card"
              :class="{ selected: field.selected }"
            >
              <input type="checkbox" v-model="field.selected" class="field-checkbox" />
              <div class="field-checkmark">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                  <polyline points="20,6 9,17 4,12"/>
                </svg>
              </div>
              <span class="field-label">{{ field.label }}</span>
              <span class="field-key">{{ field.key }}</span>
            </label>
          </div>
        </div>

        <div class="export-section">
          <h4>📥 导出格式</h4>
          <div class="format-grid">
            <div
              class="format-card"
              :class="{ active: exportFormat === 'csv' }"
              @click="exportFormat = 'csv'"
            >
              <div class="format-icon csv">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                  <polyline points="14,2 14,8 20,8"/>
                  <line x1="8" y1="13" x2="16" y2="13"/>
                  <line x1="8" y1="17" x2="16" y2="17"/>
                  <line x1="8" y1="9" x2="10" y2="9"/>
                </svg>
              </div>
              <div class="format-info">
                <span class="format-name">CSV</span>
                <span class="format-desc">Excel/表格通用</span>
              </div>
              <div class="format-badge">推荐</div>
            </div>

            <div
              class="format-card"
              :class="{ active: exportFormat === 'json' }"
              @click="exportFormat = 'json'"
            >
              <div class="format-icon json">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                  <polyline points="14,2 14,8 20,8"/>
                  <path d="M8 13h2"/>
                  <path d="M8 17h2"/>
                  <path d="M14 13h2"/>
                  <path d="M14 17h2"/>
                </svg>
              </div>
              <div class="format-info">
                <span class="format-name">JSON</span>
                <span class="format-desc">程序开发使用</span>
              </div>
            </div>

            <div
              class="format-card"
              :class="{ active: exportFormat === 'xlsx' }"
              @click="exportFormat = 'xlsx'"
            >
              <div class="format-icon xlsx">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                  <polyline points="14,2 14,8 20,8"/>
                  <path d="M8 11v5"/>
                  <path d="M11 11h3v5"/>
                  <path d="M11 11l2 2"/>
                  <path d="M15 14l2-2"/>
                </svg>
              </div>
              <div class="format-info">
                <span class="format-name">Excel</span>
                <span class="format-desc">.xlsx 格式</span>
              </div>
            </div>
          </div>
        </div>

        <div class="export-section">
          <h4>📊 导出预览</h4>
          <div class="export-preview-card">
            <div class="preview-stats">
              <div class="stat-item">
                <span class="stat-value">{{ selectedFieldCount }}</span>
                <span class="stat-label">已选字段</span>
              </div>
              <div class="stat-divider"></div>
              <div class="stat-item">
                <span class="stat-value">{{ bloggerCount }}</span>
                <span class="stat-label">数据条数</span>
              </div>
              <div class="stat-divider"></div>
              <div class="stat-item">
                <span class="stat-value">{{ exportFormat.toUpperCase() }}</span>
                <span class="stat-label">导出格式</span>
              </div>
            </div>
            <div class="preview-fields">
              <span class="preview-label">将导出：</span>
              <div class="field-tags">
                <span v-for="field in selectedFields" :key="field.key" class="field-tag">
                  {{ field.label }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="export-actions">
          <button class="export-btn primary" @click="handleExport(exportFormat)" :disabled="selectedFieldCount === 0">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="7,10 12,15 17,10"/>
              <line x1="12" y1="15" x2="12" y2="3"/>
            </svg>
            确认导出
          </button>
          <button class="export-btn secondary" @click="resetExport">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/>
              <path d="M3 3v5h5"/>
            </svg>
            重置
          </button>
        </div>

        <div v-if="exporting" class="export-progress">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: exportProgress + '%' }"></div>
          </div>
          <span class="progress-text">{{ exportProgress }}% {{ exportProgress < 30 ? '准备中...' : exportProgress < 70 ? '获取数据...' : exportProgress < 90 ? '生成文件...' : '完成！' }}</span>
        </div>

        <div v-if="exportHistory.length > 0 && !exporting" class="export-history">
          <h4 class="history-title">📋 导出历史</h4>
          <div class="history-list">
            <div v-for="(record, index) in exportHistory" :key="record.time" class="history-item">
              <span class="history-format">{{ record.format }}</span>
              <span class="history-fields">{{ record.fields }} 字段</span>
              <span class="history-size">{{ record.size }}</span>
              <span class="history-time">{{ formatExportTime(record.time) }}</span>
            </div>
          </div>
        </div>
      </div>

      <div v-if="activeTab === 'import'" class="tab-panel">
        <div class="panel-header">
          <h3>批量导入博主</h3>
        </div>
        <ImportPanel @complete="handleImportComplete" />
      </div>

      <div v-if="activeTab === 'log'" class="tab-panel">
        <div class="panel-header">
          <h3>操作日志</h3>
          <div class="log-actions">
            <button v-if="userStore.role === 'admin'" class="btn-secondary-sm" @click="handleDeleteOldLogs(10)">删除最早10条</button>
            <button v-if="userStore.role === 'admin'" class="btn-secondary-sm" @click="handleDeleteOldLogs(50)">删除最早50条</button>
            <button v-if="userStore.role === 'admin'" class="btn-danger-sm" @click="handleClearLogs">清空日志</button>
          </div>
        </div>
        <div class="log-list">
          <div v-for="(log, index) in logs" :key="log.id" class="log-item" :style="{ '--item-delay': index * 0.05 }">
            <div class="log-icon" :class="getLogIconClass(log.action)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" v-html="getLogIcon(log.action)"></svg>
            </div>
            <div class="log-content">
              <div class="log-main">
                <span class="log-action">{{ log.action }}</span>
                <span class="log-target">{{ log.target }}</span>
              </div>
              <div class="log-detail">{{ log.detail }}</div>
              <div class="log-meta">
                <span class="log-operator">{{ log.operator }}</span>
                <span class="log-time">{{ formatDateTime(log.create_time) }}</span>
              </div>
            </div>
            <button v-if="userStore.role === 'admin'" class="log-delete-btn" @click="handleDeleteLog(log.id)" title="删除">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3,6 5,6 21,6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="showModal = false">
      <div class="modal" style="background:var(--bg-card);">
        <div class="modal-header">
          <h3>{{ modalType === 'category' ? '分类' : modalType === 'team' ? '团队' : '产品' }} {{ isEdit ? '编辑' : '添加' }}</h3>
          <button class="close-btn" @click="showModal = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div v-if="modalType === 'category'" class="category-form">
            <div class="form-field">
              <label>分类名称</label>
              <input v-model="modalValue" type="text" placeholder="请输入分类名称" />
            </div>
            <div class="form-field">
              <label>分类图标</label>
              <input v-model="modalIcon" type="text" placeholder="请输入图标（如：🏷️）" />
            </div>
            <div class="form-field">
              <label>分类颜色</label>
              <div class="color-input-wrapper">
                <input v-model="modalColor" type="color" />
                <input v-model="modalColor" type="text" placeholder="请输入颜色（如：#ec4899）" />
              </div>
            </div>
          </div>
          <input v-else v-model="modalValue" type="text" :placeholder="'请输入' + (modalType === 'category' ? '分类' : '产品') + '名称'" />
        </div>
        <div class="modal-footer">
          <button class="btn-secondary" @click="showModal = false">取消</button>
          <button class="btn-primary" @click="handleSave">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useNotification } from '../stores/notification'
import { useUserStore } from '../stores/user'
import { useConfirm } from '../utils/confirm'
import { categoryListAPI, productListAPI, categoryAddAPI, categoryUpdateAPI, categoryDeleteAPI, productAddAPI, productUpdateAPI, productDeleteAPI, exportAPI, getOperationLog, getPendingUsersAPI, approveUserAPI, rejectUserAPI, getTeamsAPI, deleteTeamAPI, updateTeamAPI, createTeamAPI, setAdminAPI, removeAdminAPI, setUserTeamAPI, getUsersListAPI, deleteLogAPI, clearLogsAPI, clearOldLogsAPI, batchApproveUsersAPI, batchRejectUsersAPI, getRegistrationModeAPI, setRegistrationModeAPI, deactivateUserAdminAPI, deleteUserAdminAPI, uploadTeamLogoAPI, uploadTeamBgAPI, createBackupAPI, getSnapshotsAPI, getSnapshotSettingsAPI, setSnapshotSettingsAPI, createSnapshotAPI, restoreSnapshotAPI, downloadSnapshotAPI, deleteSnapshotAPI, clearDataAPI, purgeDataAPI } from '../api'
import ImportPanel from '../components/ImportPanel.vue'
import SkeletonLoader from '../components/SkeletonLoader.vue'

const notification = useNotification()
const userStore = useUserStore()
const router = useRouter()
const { confirm, confirmDanger, confirmWarning } = useConfirm()

const tabs = [
  { key: 'category', name: '分类管理', icon: '<path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>' },
  { key: 'product', name: '产品管理', icon: '<path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>' },
  { key: 'team', name: '团队管理', icon: '<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/>' },
  { key: 'users', name: '用户审核', icon: '<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/>' },
  { key: 'import', name: '批量导入', icon: '<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17,8 12,3 7,8"/><line x1="12" y1="3" x2="12" y2="15"/>' },
  { key: 'backup', name: '数据库管理', icon: '<path d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/>' },
  { key: 'export', name: '数据导出', icon: '<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7,10 12,15 17,10"/><line x1="12" y1="15" x2="12" y2="3"/>' },
  { key: 'log', name: '操作日志', icon: '<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14,2 14,8 20,8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/>' }
]

const activeTab = ref('category')
const categories = ref([])
const products = ref([])
const logs = ref([])
const pendingUsers = ref([])
const allUsers = ref([])
const teams = ref([])
const teamDetailVisible = ref(false)
const selectedTeam = ref(null)
const loading = ref(false)
const snapshots = ref([])
const maxSnapshots = 10
const backupInterval = ref(604800000)

const selectedUsers = ref([])
const compactView = ref(false)
const registrationMode = ref('normal')

const bloggerCount = ref(0)
const exportFields = ref([
  { key: 'id', label: 'ID', selected: true },
  { key: 'nickname', label: '昵称', selected: true },
  { key: 'category', label: '分类', selected: true },
  { key: 'platform', label: '平台', selected: true },
  { key: 'contact', label: '联系方式', selected: false },
  { key: 'wechat', label: '微信', selected: false },
  { key: 'product', label: '对接产品', selected: false },
  { key: 'status', label: '状态', selected: true },
  { key: 'blogger_tags', label: '博主标签', selected: true },
  { key: 'followup_count', label: '跟进次数', selected: true },
  { key: 'cooperation_count', label: '合作次数', selected: true },
  { key: 'user_name', label: '用户名', selected: false },
  { key: 'real_name', label: '真实姓名', selected: false },
  { key: 'description', label: '简介', selected: false },
  { key: 'platform_account', label: '平台账号', selected: false },
  { key: 'custom_contact', label: '自定义联系方式', selected: false },
  { key: 'tags', label: '标签', selected: false },
  { key: 'create_time', label: '创建时间', selected: true },
  { key: 'update_time', label: '更新时间', selected: false },
])

const selectedFieldCount = computed(() => exportFields.value.filter(f => f.selected).length)
const selectedFieldNames = computed(() => exportFields.value.filter(f => f.selected).map(f => f.label).join('、'))
const exporting = ref(false)
const exportProgress = ref(0)
const exportHistory = ref(JSON.parse(localStorage.getItem('orange_export_history') || '[]'))

const saveExportHistory = (record) => {
  const history = [record, ...exportHistory.value].slice(0, 10)
  exportHistory.value = history
  localStorage.setItem('orange_export_history', JSON.stringify(history))
}

const formatExportTime = (timestamp) => {
  const date = new Date(timestamp)
  const now = Date.now()
  const diff = now - timestamp
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + ' 分钟前'
  if (diff < 86400000) return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' }) + ' ' + date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

const selectAllFields = () => {
  exportFields.value.forEach(f => f.selected = true)
}

const exportFormat = ref('csv')

const resetExport = () => {
  exportFields.value.forEach(f => f.selected = false)
  exportFormat.value = 'csv'
}

const selectedFields = computed(() => exportFields.value.filter(f => f.selected))

const selectBasicFields = () => {
  exportFields.value.forEach(f => {
    f.selected = ['id', 'nickname', 'category', 'platform', 'status', 'create_time'].includes(f.key)
  })
}

const clearFields = () => {
  exportFields.value.forEach(f => f.selected = false)
}

const isAllSelected = computed(() => {
  if (allUsers.value.length === 0) return false
  return allUsers.value.every(u => selectedUsers.value.includes(u.id))
})

const pendingSelectedCount = computed(() => {
  return allUsers.value.filter(u => selectedUsers.value.includes(u.id) && u.status === 'pending').length
})

const showModal = ref(false)
const modalType = ref('')
const modalValue = ref('')
const modalIcon = ref('')
const modalColor = ref('')
const modalDescription = ref('')
const isEdit = ref(false)
const editingId = ref(null)
const formPanelVisible = ref(false)
const formData = ref({
  name: '',
  description: '',
  color: '#6366f1',
  logo: '',
  bg_image: '',
  leader_id: '',
  announcement: ''
})
const formError = ref('')
const saving = ref(false)
const currentTeam = ref(null)
const logoInput = ref(null)
const bgInput = ref(null)

const loadCategories = async () => {
  try {
    const res = await categoryListAPI()
    if (res.code === 200) categories.value = res.data || []
  } catch (error) {
    console.error('加载分类失败', error)
  }
}

const loadProducts = async () => {
  try {
    const res = await productListAPI()
    if (res.code === 200) products.value = res.data || []
  } catch (error) {
    console.error('加载产品失败', error)
  }
}

const loadLogs = async () => {
  try {
    const res = await getOperationLog({ page: 1, pageSize: 50 })
    if (res.code === 200) logs.value = res.data?.list || []
  } catch (error) {
    console.error('加载日志失败', error)
  }
}

const handleDeleteLog = async (id) => {
  if (!await confirmDanger('确定删除该日志？')) return
  try {
    const data = await deleteLogAPI(id)
    if (data.code === 200) {
      notification.success('删除成功')
      loadLogs()
    } else {
      notification.error(data.message || '删除失败')
    }
  } catch (error) {
    notification.error('删除失败')
  }
}

const handleClearLogs = async () => {
  if (!await confirmDanger('警告：确定清空所有操作日志？此操作不可恢复！')) return
  try {
    const data = await clearLogsAPI()
    if (data.code === 200) {
      notification.success('清空成功')
      loadLogs()
    } else {
      notification.error(data.message || '清空失败')
    }
  } catch (error) {
    notification.error('清空失败')
  }
}

const handleDeleteOldLogs = async (count) => {
  if (!await confirmDanger(`确定删除最早的 ${count} 条日志？`)) return
  try {
    const data = await clearOldLogsAPI(count)
    if (data.code === 200) {
      notification.success(data.message)
      loadLogs()
    } else {
      notification.error(data.message || '删除失败')
    }
  } catch (error) {
    notification.error('删除失败')
  }
}

const loadPendingUsers = async () => {
  loading.value = true
  try {
    const res = await getPendingUsersAPI()
    if (res.code === 200) pendingUsers.value = res.data || []
  } catch (error) {
    console.error('加载待审核用户失败', error)
  } finally {
    loading.value = false
  }
}

const loadAllUsers = async () => {
  loading.value = true
  try {
    const res = await getUsersListAPI()
    if (res.code === 200) allUsers.value = res.data || []
  } catch (error) {
    console.error('加载用户列表失败', error)
  } finally {
    loading.value = false
  }
}

const loadTeams = async () => {
  loading.value = true
  try {
    const res = await getTeamsAPI()
    if (res.code === 200) teams.value = res.data || []
  } catch (error) {
    console.error('加载团队列表失败', error)
  } finally {
    loading.value = false
  }
}

const getTeamMembers = (teamId) => {
  return allUsers.value.filter(u => u.team_id === teamId)
}

const isTeamMember = (team) => {
  if (!team || !userStore.team_id) return false
  return String(team.id) === String(userStore.team_id)
}

const getAvailableMembers = (teamId) => {
  return allUsers.value.filter(u => !u.team_id || u.team_id !== teamId)
}

const handleAddMember = async (teamId, userId) => {
  if (!userId) return
  await handleSetUserTeam(Number(userId), teamId)
  loadTeams()
}

const handleRemoveMember = async (userId, teamId) => {
  await handleSetUserTeam(userId, null)
  loadTeams()
}

const handleApprove = async (id) => {
  if (!await confirm('确定批准该用户的注册申请？')) return

  try {
    const res = await approveUserAPI(id)
    if (res.code === 200) {
      notification.success('已批准该用户的注册申请')
      loadPendingUsers()
      loadAllUsers()
      loadLogs()
    } else {
      notification.error(res.message || '操作失败')
    }
  } catch (error) {
    notification.error('操作失败，请重试')
  }
}

const handleReject = async (id) => {
  if (!await confirmDanger('确定拒绝该用户的注册申请？')) return

  try {
    const res = await rejectUserAPI(id)
    if (res.code === 200) {
      notification.success('已拒绝该用户的注册申请')
      loadPendingUsers()
      loadAllUsers()
      loadLogs()
    } else {
      notification.error(res.message || '操作失败')
    }
  } catch (error) {
    notification.error('操作失败，请重试')
  }
}

const toggleSelectUser = (id) => {
  const index = selectedUsers.value.indexOf(id)
  if (index > -1) {
    selectedUsers.value.splice(index, 1)
  } else {
    selectedUsers.value.push(id)
  }
}

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedUsers.value = []
  } else {
    selectedUsers.value = allUsers.value.map(u => u.id)
  }
}

const handleBatchApprove = async () => {
  const pendingIds = allUsers.value
    .filter(u => selectedUsers.value.includes(u.id) && u.status === 'pending')
    .map(u => u.id)
  
  if (pendingIds.length === 0) {
    notification.warning('请选择待审核的用户')
    return
  }
  
  if (!await confirm(`确定批量批准 ${pendingIds.length} 个用户的注册申请？`)) return
  
  try {
    const data = await batchApproveUsersAPI(pendingIds)
    if (data.code === 200) {
      notification.success(`已批量批准 ${pendingIds.length} 个用户`)
      selectedUsers.value = []
      loadAllUsers()
      loadPendingUsers()
      loadLogs()
    } else {
      notification.error(data.message || '操作失败')
    }
  } catch (error) {
    notification.error('操作失败，请重试')
  }
}

const handleBatchReject = async () => {
  const pendingIds = allUsers.value
    .filter(u => selectedUsers.value.includes(u.id) && u.status === 'pending')
    .map(u => u.id)
  
  if (pendingIds.length === 0) {
    notification.warning('请选择待审核的用户')
    return
  }
  
  if (!await confirmDanger(`确定批量拒绝 ${pendingIds.length} 个用户的注册申请？`)) return
  
  try {
    const data = await batchRejectUsersAPI(pendingIds)
    if (data.code === 200) {
      notification.success(`已批量拒绝 ${pendingIds.length} 个用户`)
      selectedUsers.value = []
      loadAllUsers()
      loadPendingUsers()
      loadLogs()
    } else {
      notification.error(data.message || '操作失败')
    }
  } catch (error) {
    notification.error('操作失败，请重试')
  }
}

const setRegistrationMode = async (mode) => {
  const modeNames = {
    'normal': '正常审核',
    'open': '开放注册',
    'closed': '禁止注册'
  }
  
  if (!await confirm(`确定切换到"${modeNames[mode]}"模式？`)) return
  
  try {
    const data = await setRegistrationModeAPI(mode)
    if (data.code === 200) {
      registrationMode.value = mode
      notification.success(`已切换到"${modeNames[mode]}"模式`)
      loadLogs()
    } else {
      notification.error(data.message || '操作失败')
    }
  } catch (error) {
    notification.error('操作失败，请重试')
  }
}

const loadRegistrationMode = async () => {
  try {
    const data = await getRegistrationModeAPI()
    if (data.code === 200) {
      registrationMode.value = data.data?.mode || 'normal'
    }
  } catch (error) {
    console.error('加载注册模式失败', error)
  }
}

const handleSetAdmin = async (id) => {
  if (!await confirmWarning('确定将该用户设置为管理员？')) return
  try {
    const res = await setAdminAPI(id)
    if (res.code === 200) {
      notification.success(res.message)
      loadAllUsers()
      loadLogs()
    } else {
      notification.error(res.message || '操作失败')
    }
  } catch (error) {
    notification.error('操作失败')
  }
}

const handleRemoveAdmin = async (id) => {
  if (!await confirmWarning('确定取消该用户的管理员权限？')) return
  try {
    const res = await removeAdminAPI(id)
    if (res.code === 200) {
      notification.success(res.message)
      loadAllUsers()
      loadLogs()
    } else {
      notification.error(res.message || '操作失败')
    }
  } catch (error) {
    notification.error('操作失败')
  }
}

const handleDeactivate = async (id) => {
  if (!await confirmDanger('确定注销该用户？')) return
  try {
    const data = await deactivateUserAdminAPI(id)
    if (data.code === 200) {
      notification.success('已注销该用户')
      loadAllUsers()
      loadLogs()
    } else {
      notification.error(data.message || '操作失败')
    }
  } catch (error) {
    notification.error('操作失败')
  }
}

const handleDeleteUser = async (id) => {
  if (!await confirmDanger('确定永久删除该用户？此操作不可恢复！')) return
  try {
    const data = await deleteUserAdminAPI(id)
    if (data.code === 200) {
      notification.success('已永久删除该用户')
      loadAllUsers()
      loadLogs()
    } else {
      notification.error(data.message || '操作失败')
    }
  } catch (error) {
    notification.error('操作失败')
  }
}

const editTeam = (team) => {
  formData.value = {
    name: team.name,
    description: team.description || '',
    color: team.color || '#6366f1',
    logo: team.logo || '',
    bg_image: team.bg_image || '',
    leader_id: team.leader_id || '',
    announcement: team.announcement || ''
  }
  currentTeam.value = team
  isEdit.value = true
  editingId.value = team.id
  formError.value = ''
  formPanelVisible.value = true
}

const closeFormPanel = () => {
  formPanelVisible.value = false
  isEdit.value = false
  editingId.value = null
  currentTeam.value = null
  formError.value = ''
  formData.value = {
    name: '',
    description: '',
    color: '#6366f1',
    logo: '',
    bg_image: '',
    leader_id: '',
    announcement: ''
  }
}

const isTeamLeaderOrAdmin = (team) => {
  if (!team) return false
  if (userStore.role === 'admin') return true
  if (team.leader_id && team.leader_id === userStore.id) return true
  return false
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
  reader.onload = async () => {
    const base64Data = reader.result
    try {
      const data = await uploadTeamLogoAPI({ image: base64Data })
      if (data.code === 200) {
        formData.value.logo = data.data.url
        notification.success('Logo上传成功')
      } else {
        notification.error(data.message || '上传失败')
      }
    } catch (error) {
      notification.error('上传失败')
    }
  }
  reader.readAsDataURL(file)
  e.target.value = ''
}

const handleBgUpload = async (e) => {
  const file = e.target.files[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = async () => {
    const base64Data = reader.result
    try {
      const data = await uploadTeamBgAPI({ image: base64Data })
      if (data.code === 200) {
        formData.value.bg_image = data.data.url
        notification.success('背景上传成功')
      } else {
        notification.error(data.message || '上传失败')
      }
    } catch (error) {
      notification.error('上传失败')
    }
  }
  reader.readAsDataURL(file)
  e.target.value = ''
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
    if (isEdit.value) {
      res = await updateTeamAPI(editingId.value, formData.value)
    } else {
      res = await createTeamAPI(formData.value)
    }

    if (res.code === 200) {
      notification.success(isEdit.value ? '团队更新成功' : '团队创建成功')
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

const handleSetUserTeam = async (userId, teamId) => {
  const team_id = teamId === '' ? null : Number(teamId)
  try {
    const res = await setUserTeamAPI(userId, team_id)
    if (res.code === 200) {
      notification.success('团队分配成功')
      loadAllUsers()
      loadLogs()
    } else {
      notification.error(res.message || '分配失败')
    }
  } catch (error) {
    notification.error('分配失败')
  }
}

const deleteTeam = async (team) => {
  if (!await confirmDanger(`确定删除团队【${team.name}】？删除后成员将脱离该团队。`)) return
  try {
    const res = await deleteTeamAPI(team.id)
    if (res.code === 200) {
      notification.success('团队已删除')
      loadTeams()
      loadAllUsers()
      loadLogs()
    } else {
      notification.error(res.message || '删除失败')
    }
  } catch (error) {
    notification.error('删除失败')
  }
}

const openTeamDetail = (team) => {
  selectedTeam.value = team
  teamDetailVisible.value = true
}

const goToTeamHome = (team) => {
  router.push({ name: 'TeamHome', params: { teamId: team.id } })
}

const closeTeamDetail = () => {
  teamDetailVisible.value = false
  selectedTeam.value = null
}

const handlePermanentDeleteTeam = async (team) => {
  if (!await confirmDanger(`确定永久删除团队【${team.name}】？此操作不可恢复！`)) return
  try {
    const res = await deleteTeamAPI(team.id)
    if (res.code === 200) {
      notification.success('团队已永久删除')
      closeTeamDetail()
      loadTeams()
      loadAllUsers()
      loadLogs()
    } else {
      notification.error(res.message || '删除失败')
    }
  } catch (error) {
    notification.error('删除失败')
  }
}

const showAddModal = (type) => {
  modalType.value = type
  modalValue.value = ''
  modalIcon.value = '🏷️'
  modalColor.value = type === 'team' ? '#6366f1' : '#6b7280'
  modalDescription.value = ''
  isEdit.value = false
  editingId.value = null
  showModal.value = true
}

const editItem = (type, item) => {
  modalType.value = type
  modalValue.value = item.name
  modalIcon.value = item.icon || '🏷️'
  modalColor.value = item.color || '#6b7280'
  isEdit.value = true
  editingId.value = item.id
  showModal.value = true
}

const handleSave = async () => {
  if (!modalValue.value.trim()) {
    notification.warning('请输入名称')
    return
  }

  try {
    let res
    if (modalType.value === 'category') {
      if (isEdit.value) {
        res = await categoryUpdateAPI({
          id: editingId.value,
          name: modalValue.value,
          color: modalColor.value,
          icon: modalIcon.value
        })
      } else {
        res = await categoryAddAPI({
          name: modalValue.value,
          color: modalColor.value,
          icon: modalIcon.value
        })
      }
      loadCategories()
    } else if (modalType.value === 'team') {
      if (isEdit.value) {
        res = await updateTeamAPI(editingId.value, {
          name: modalValue.value,
          description: modalDescription.value,
          color: modalColor.value
        })
        loadTeams()
      } else {
        res = await createTeamAPI({
          name: modalValue.value,
          description: modalDescription.value,
          color: modalColor.value
        })
        loadTeams()
      }
    } else {
      if (isEdit.value) {
        res = await productUpdateAPI({ id: editingId.value, name: modalValue.value })
      } else {
        res = await productAddAPI({ name: modalValue.value })
      }
      loadProducts()
    }

    if (res.code === 200) {
      showModal.value = false
      notification.success(isEdit.value ? '修改成功' : '添加成功')
    } else {
      notification.error(res.message || '操作失败')
    }
  } catch (error) {
    notification.error('操作失败')
  }
}

const deleteItem = async (type, item) => {
  if (!await confirmDanger(`确定要删除"${item.name}"吗？`)) return

  try {
    let res
    if (type === 'category') {
      res = await categoryDeleteAPI(item.id)
      loadCategories()
    } else {
      res = await productDeleteAPI(item.id)
      loadProducts()
    }

    if (res.code === 200) {
      notification.success('删除成功')
    } else {
      notification.error(res.message || '删除失败')
    }
  } catch (error) {
    notification.error('删除失败')
  }
}

const handleExport = async (type) => {
  exporting.value = true
  exportProgress.value = 10
  try {
    const selectedFields = exportFields.value.filter(f => f.selected).map(f => f.key).join(',')
    exportProgress.value = 30

    const res = await exportAPI(type, selectedFields)
    exportProgress.value = 70

    const blob = new Blob([res], {
      type: type === 'json' ? 'application/json' : type === 'xlsx' ? 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' : 'text/csv'
    })
    const fileSize = (blob.size / 1024).toFixed(1)
    exportProgress.value = 90

    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    const timestamp = new Date().toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }).replace(/\//g, '-')
    a.download = `orange_export_${timestamp}.${type}`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)

    exportProgress.value = 100

    saveExportHistory({
      time: Date.now(),
      format: type.toUpperCase(),
      fields: selectedFieldCount.value,
      size: `${fileSize}KB`
    })

    setTimeout(() => { exporting.value = false; exportProgress.value = 0 }, 600)
    notification.success('导出成功！', `已导出 ${selectedFieldCount.value} 个字段 (${fileSize}KB)`)
  } catch (error) {
    exporting.value = false
    exportProgress.value = 0
    console.error('导出失败:', error)
    notification.error('导出失败，请重试')
  }
}

const importFileInput = ref(null)

const triggerImport = () => {
  importFileInput.value.click()
}

const handleBackup = async () => {
  try {
    const data = await createBackupAPI({ name: 'backup_' + new Date().toISOString().slice(0, 10) })
    if (data.code === 200) {
      notification.success('备份成功！')
      loadSnapshots()
    } else {
      notification.error(data.message || '备份失败')
    }
  } catch (error) {
    console.error('备份失败:', error)
    notification.error('备份失败，请重试')
  }
}

const loadSnapshots = async () => {
  try {
    const data = await getSnapshotsAPI()
    if (data.code === 200) {
      snapshots.value = data.data || []
    }
    if (userStore.username === 'admin') {
      const settingsData = await getSnapshotSettingsAPI()
      if (settingsData.code === 200) {
        const intervalHours = settingsData.data.backup_interval || 24
        backupInterval.value = intervalHours * 3600000
      }
    }
  } catch (error) {
    console.error('加载快照失败:', error)
  }
}

const handleUpdateBackupInterval = async () => {
  try {
    const intervalHours = Math.round(backupInterval.value / 3600000)
    const data = await setSnapshotSettingsAPI({
      auto_backup: true,
      backup_interval: intervalHours,
      keep_backups: 7
    })
    if (data.code === 200) {
      notification.success('自动快照间隔已更新')
    } else {
      notification.error(data.message || '更新失败')
    }
  } catch (error) {
    console.error('更新失败:', error)
    notification.error('更新失败')
  }
}

const handleCreateSnapshot = async () => {
  try {
    const data = await createSnapshotAPI({ name: 'manual' })
    if (data.code === 200) {
      notification.success('快照创建成功！')
      loadSnapshots()
    } else {
      notification.error(data.message || '创建失败')
    }
  } catch (error) {
    console.error('创建快照失败:', error)
    notification.error('创建快照失败')
  }
}

const handleRestoreSnapshot = async (snap) => {
  if (!await confirmDanger(`确定要还原快照 "${snap.name}" 吗？当前所有数据将被覆盖！`)) return
  try {
    const data = await restoreSnapshotAPI(snap.filename)
    if (data.code === 200) {
      notification.success('还原成功！页面将刷新...')
      setTimeout(() => window.location.reload(), 1500)
    } else {
      notification.error(data.message || '还原失败')
    }
  } catch (error) {
    console.error('还原失败:', error)
    notification.error('还原失败')
  }
}

const handleDownloadSnapshot = async (snap) => {
  try {
    const blob = await downloadSnapshotAPI(snap.filename)
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = snap.filename
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)
    notification.success('下载成功！')
  } catch (error) {
    console.error('下载失败:', error)
    notification.error('下载失败')
  }
}

const handleDeleteSnapshot = async (snap) => {
  if (!await confirmDanger(`确定要删除快照 "${snap.name}" 吗？`)) return
  try {
    const data = await deleteSnapshotAPI(snap.filename)
    if (data.code === 200) {
      notification.success('删除成功！')
      loadSnapshots()
    } else {
      notification.error(data.message || '删除失败')
    }
  } catch (error) {
    console.error('删除失败:', error)
    notification.error('删除失败')
  }
}

const handleClearData = async () => {
  if (!await confirmDanger('⚠️ 警告：此操作将删除所有博主、跟进记录、合作记录和操作日志！\n\n此操作不可恢复！\n\n确定要继续吗？')) return
  if (!await confirmDanger('再次确认：清空后所有数据将被永久删除！\n\n点击"确定"开始清空...')) return

  try {
    const data = await clearDataAPI()
    if (data.code === 200) {
      notification.success('数据库已清空！')
      setTimeout(() => window.location.reload(), 1500)
    } else {
      notification.error(data.message || '清空失败')
    }
  } catch (error) {
    console.error('清空失败:', error)
    notification.error('清空失败')
  }
}

const handleImport = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  if (!await confirmDanger('警告：还原操作将覆盖当前所有数据，且不可恢复！确定要继续吗？')) {
    event.target.value = ''
    return
  }

  try {
    const text = await file.text()
    const data = JSON.parse(text)

    const result = await createBackupAPI(data)
    if (result.code === 200) {
      notification.success('还原成功！页面将刷新...')
      setTimeout(() => {
        window.location.reload()
      }, 1500)
    } else {
      notification.error(result.message || '还原失败')
    }
  } catch (error) {
    console.error('还原失败:', error)
    notification.error('还原失败，请检查文件格式是否正确')
  } finally {
    event.target.value = ''
  }
}

const handlePurge = async () => {
  if (!await confirmDanger('警告：将永久删除30天前的所有已删除记录，此操作不可恢复！\n\n是否继续？')) return
  try {
    const data = await purgeDataAPI({ confirm: true })
    if (data.code === 200) {
      notification.success(data.message)
    } else {
      notification.error(data.message || '清理失败')
    }
  } catch (error) {
    notification.error('清理失败')
  }
}

const handleImportComplete = (result) => {
  notification.success(`导入完成：成功 ${result.successCount} 条，失败 ${result.failCount} 条`)
}

const getLogIcon = (action) => {
  const icons = {
    '登录': '<path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10,17 15,12 10,7"/><line x1="15" y1="12" x2="3" y2="12"/>',
    '新增': '<circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/>',
    '编辑': '<path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>',
    '删除': '<polyline points="3,6 5,6 21,6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>'
  }
  return icons[action] || icons['新增']
}

const getLogIconClass = (action) => {
  const classes = {
    '登录': 'login',
    '新增': 'add',
    '编辑': 'edit',
    '删除': 'delete'
  }
  return classes[action] || 'add'
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const formatInterval = (ms) => {
  if (!ms) return '未知'
  const hours = ms / 3600000
  if (hours < 24) return `${hours}小时`
  const days = hours / 24
  return `${days}天`
}

const formatDateTime = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  loadCategories()
  loadProducts()
  loadLogs()
  loadPendingUsers()
  loadRegistrationMode()
})

watch(activeTab, (newTab) => {
  if (newTab === 'backup') {
    loadSnapshots()
  }
  if (newTab === 'users') {
    loadAllUsers()
    loadRegistrationMode()
  }
  if (newTab === 'team') {
    loadTeams()
    loadAllUsers()
  }
})

const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}
</script>

<style scoped>
.admin {
  animation: fadeIn 0.4s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
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

.admin-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  overflow-x: auto;
  padding-bottom: 8px;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.tab-btn svg {
  width: 18px;
  height: 18px;
}

.tab-btn:hover {
  background: var(--bg-card-hover);
  color: var(--text-primary);
}

.tab-btn.active {
  background: var(--primary);
  border-color: var(--primary);
  color: white;
}

.tab-panel {
  animation: slideUp 0.4s ease;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.panel-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.add-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  border: none;
  border-radius: 10px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.add-btn svg {
  width: 16px;
  height: 16px;
}

.add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.data-table {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  overflow: hidden;
}

.table-header {
  display: grid;
  grid-template-columns: 80px 1fr 150px 120px;
  gap: 16px;
  padding: 16px 20px;
  background: var(--bg-dark);
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
}

.table-row {
  display: grid;
  grid-template-columns: 80px 1fr 150px 120px;
  gap: 16px;
  padding: 16px 20px;
  align-items: center;
  border-bottom: 1px solid var(--border-color);
  transition: background 0.3s ease;
}

.table-row:last-child {
  border-bottom: none;
}

.table-row:hover {
  background: var(--bg-dark);
}

.table-row span {
  font-size: 14px;
  color: var(--text-secondary);
}

.action-btn {
  width: 32px;
  height: 32px;
  background: rgba(59, 130, 246, 0.1);
  border: none;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #60a5fa;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-right: 8px;
}

.action-btn:last-child {
  margin-right: 0;
}

.action-btn:hover {
  background: rgba(59, 130, 246, 0.2);
}

.action-btn.danger {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.action-btn.danger:hover {
  background: rgba(239, 68, 68, 0.2);
}

.action-btn svg {
  width: 16px;
  height: 16px;
}

.export-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.backup-section {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.backup-actions {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.backup-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 28px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.backup-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
}

.backup-card.import:hover {
  box-shadow: 0 12px 24px rgba(59, 130, 246, 0.3);
}

.backup-icon {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #10b981, #059669);
  border-radius: 14px;
}

.backup-card.import .backup-icon {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
}

.backup-card .backup-icon.orange {
  background: linear-gradient(135deg, #f97316, #ea580c);
}

.backup-icon svg {
  width: 28px;
  height: 28px;
  color: white;
}

.backup-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.backup-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.backup-desc {
  font-size: 14px;
  color: var(--text-secondary);
}

.backup-warning {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 12px;
  color: #ef4444;
  font-size: 14px;
}

.backup-warning svg {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.danger-zone {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 2px dashed rgba(239, 68, 68, 0.3);
}

.danger-zone h4 {
  font-size: 14px;
  font-weight: 600;
  color: #ef4444;
  margin-bottom: 16px;
}

.clear-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: #ef4444;
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.clear-btn:hover {
  background: #dc2626;
  transform: translateY(-2px);
}

.clear-btn svg {
  width: 18px;
  height: 18px;
}

.snapshot-section {
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h4 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.snapshot-count {
  font-size: 13px;
  color: var(--text-secondary);
}

.snapshot-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.auto-backup-settings {
  margin-bottom: 16px;
  padding: 12px 16px;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 8px;
}

.setting-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.setting-label {
  font-size: 14px;
  color: var(--text-primary);
}

.interval-select {
  padding: 6px 12px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-card);
  color: var(--text-primary);
  font-size: 14px;
  cursor: pointer;
}

.interval-select:focus {
  outline: none;
  border-color: var(--primary);
}

.snapshot-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  transition: all 0.2s;
}

.snapshot-item:hover {
  border-color: var(--primary);
}

.snapshot-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.snapshot-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.snapshot-meta {
  font-size: 12px;
  color: var(--text-secondary);
}

.snapshot-actions {
  display: flex;
  gap: 8px;
}

.snapshot-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.snapshot-btn svg {
  width: 18px;
  height: 18px;
}

.snapshot-btn.restore {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.snapshot-btn.restore:hover {
  background: #10b981;
  color: white;
}

.snapshot-btn.download {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.snapshot-btn.download:hover {
  background: #3b82f6;
  color: white;
}

.snapshot-btn.delete {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.snapshot-btn.delete:hover {
  background: #ef4444;
  color: white;
}

.snapshot-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px;
  background: var(--bg-card);
  border: 1px dashed var(--border-color);
  border-radius: 12px;
  color: var(--text-secondary);
  gap: 12px;
}

.snapshot-empty svg {
  width: 48px;
  height: 48px;
  opacity: 0.5;
}

.export-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 28px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.export-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.3);
  border-color: var(--primary);
}

.export-icon {
  width: 60px;
  height: 60px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.export-icon svg {
  width: 30px;
  height: 30px;
}

.export-icon.xlsx {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.2), rgba(34, 197, 94, 0.1));
  color: #4ade80;
}

.export-icon.csv {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.2), rgba(59, 130, 246, 0.1));
  color: #60a5fa;
}

.export-icon.json {
  background: linear-gradient(135deg, rgba(249, 115, 22, 0.2), rgba(249, 115, 22, 0.1));
  color: #f97316;
}

.export-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.export-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.export-desc {
  font-size: 13px;
  color: var(--text-muted);
}

.panel-desc {
  font-size: 14px;
  color: var(--text-muted);
  margin-top: 4px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h4 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
}

.quick-actions {
  display: flex;
  gap: 8px;
}

.quick-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.quick-btn svg {
  width: 16px;
  height: 16px;
}

.quick-btn:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.quick-btn.danger:hover {
  border-color: var(--danger);
  color: var(--danger);
}

.field-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
  margin-top: 16px;
}

.field-card {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px 12px;
  background: var(--bg-card);
  border: 2px solid var(--border-color);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.field-card:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.1);
}

.field-card.selected {
  border-color: var(--primary);
  background: linear-gradient(135deg, rgba(249, 115, 22, 0.05) 0%, rgba(234, 88, 12, 0.05) 100%);
}

.field-checkbox {
  display: none;
}

.field-checkmark {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--bg-dark);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.field-card.selected .field-checkmark {
  background: var(--primary);
}

.field-checkmark svg {
  width: 12px;
  height: 12px;
  color: white;
  opacity: 0;
  transition: opacity 0.2s;
}

.field-card.selected .field-checkmark svg {
  opacity: 1;
}

.field-label {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  text-align: center;
}

.field-key {
  font-size: 11px;
  color: var(--text-muted);
  font-family: var(--font-mono);
}

.format-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-top: 16px;
}

.format-card {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 24px 16px;
  background: var(--bg-card);
  border: 2px solid var(--border-color);
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.2s;
}

.format-card:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(249, 115, 22, 0.15);
}

.format-card.active {
  border-color: var(--primary);
  background: linear-gradient(135deg, rgba(249, 115, 22, 0.08) 0%, rgba(234, 88, 12, 0.08) 100%);
  box-shadow: 0 4px 16px rgba(249, 115, 22, 0.2);
}

.format-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.format-icon svg {
  width: 28px;
  height: 28px;
}

.format-icon.csv {
  background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
  color: white;
}

.format-icon.json {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
}

.format-icon.xlsx {
  background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
  color: white;
}

.format-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.format-name {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
}

.format-desc {
  font-size: 12px;
  color: var(--text-muted);
}

.format-badge {
  position: absolute;
  top: -8px;
  right: 16px;
  padding: 4px 10px;
  background: linear-gradient(135deg, #f97316 0%, #ea580c 100%);
  color: white;
  font-size: 11px;
  font-weight: 600;
  border-radius: 10px;
}

.export-preview-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  padding: 20px;
  margin-top: 16px;
}

.preview-stats {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 16px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--primary);
}

.stat-label {
  font-size: 12px;
  color: var(--text-muted);
}

.stat-divider {
  width: 1px;
  height: 40px;
  background: var(--border-color);
}

.preview-fields {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.preview-label {
  font-size: 13px;
  color: var(--text-muted);
}

.field-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.field-tag {
  padding: 6px 12px;
  background: rgba(249, 115, 22, 0.1);
  color: var(--primary);
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
}

.export-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.export-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px 24px;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.export-btn svg {
  width: 20px;
  height: 20px;
}

.export-btn.primary {
  background: linear-gradient(135deg, #f97316 0%, #ea580c 100%);
  border: none;
  color: white;
}

.export-btn.primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(249, 115, 22, 0.35);
}

.export-btn.primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.export-btn.secondary {
  background: var(--bg-card);
  border: 2px solid var(--border-color);
  color: var(--text-secondary);
}

.export-btn.secondary:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.export-progress {
  margin-top: 20px;
  padding: 16px;
  background: rgba(249, 115, 22, 0.04);
  border: 1px solid rgba(249, 115, 22, 0.12);
  border-radius: 12px;
}

.progress-bar {
  height: 6px;
  background: var(--border-color);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 10px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #f97316, #fb923c);
  border-radius: 3px;
  transition: width 0.3s ease-out;
}

.progress-text {
  font-size: 13px;
  color: var(--text-muted);
  font-weight: 500;
}

.export-history {
  margin-top: 20px;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 12px;
  border: 1px solid var(--border-color);
}

.history-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background: var(--bg-card);
  border-radius: 8px;
  font-size: 13px;
}

.history-format {
  padding: 2px 8px;
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
  border-radius: 4px;
  font-weight: 600;
  font-size: 11px;
}

.history-fields { color: var(--text-secondary); }
.history-size { color: var(--text-muted); margin-left: auto; }
.history-time { color: var(--text-muted); font-size: 12px; }

.log-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.log-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 14px;
  animation: slideUp 0.4s ease forwards;
  animation-delay: calc(var(--item-delay) * 1s);
  opacity: 0;
  position: relative;
  transition: all 0.3s ease;
}

.log-item:hover {
  transform: translateX(4px);
  border-color: rgba(249, 115, 22, 0.3);
}

.log-item:hover .log-delete-btn {
  opacity: 1;
}

.log-delete-btn {
  opacity: 0;
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  transition: all 0.2s;
}

.log-delete-btn:hover {
  color: var(--danger);
  background: rgba(239, 68, 68, 0.1);
}

.log-content {
  flex: 1;
  position: relative;
}

.btn-danger-sm {
  padding: 6px 12px;
  font-size: 12px;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-danger-sm:hover {
  background: #dc2626;
}

.btn-secondary-sm {
  padding: 6px 12px;
  font-size: 12px;
  background: #6b7280;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary-sm:hover {
  background: #4b5563;
}

.log-actions {
  display: flex;
  gap: 8px;
}

.log-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.log-icon svg {
  width: 18px;
  height: 18px;
}

.log-icon.login {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
}

.log-icon.add {
  background: rgba(34, 197, 94, 0.15);
  color: #4ade80;
}

.log-icon.edit {
  background: rgba(249, 115, 22, 0.15);
  color: var(--primary);
}

.log-icon.delete {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.log-content {
  flex: 1;
  min-width: 0;
}

.log-main {
  display: flex;
  gap: 8px;
  margin-bottom: 4px;
}

.log-action {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.log-target {
  font-size: 14px;
  color: var(--primary);
}

.log-detail {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.log-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--text-muted);
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.modal {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 20px;
  width: 100%;
  max-width: 520px;
  max-height: 90vh;
  overflow-y: auto;
  animation: slideUp 0.3s ease;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.team-form-modal {
  max-width: 560px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 2px solid #e5e7eb;
  background: #f9fafb;
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
}

.close-btn {
  width: 32px;
  height: 32px;
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: var(--danger);
  border-color: var(--danger);
  color: white;
}

.close-btn svg {
  width: 16px;
  height: 16px;
}

.modal-body {
  padding: 24px;
  background: var(--bg-card);
}

.modal-body input {
  width: 100%;
  height: 48px;
  background: #f9fafb;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  padding: 0 16px;
  font-size: 15px;
  color: #1f2937;
  transition: all 0.3s ease;
}

.modal-body input:focus {
  outline: none;
  border-color: #ff6b35;
  box-shadow: 0 0 0 3px rgba(255, 107, 53, 0.1);
  background: #ffffff;
}

.modal-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 20px 24px;
  border-top: 2px solid #e5e7eb;
  background: #f9fafb;
}

.btn-secondary,
.btn-primary {
  padding: 12px 24px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-secondary {
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
}

.btn-secondary:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary), #f7931e);
  border: none;
  color: white;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 107, 53, 0.3);
}

.loading {
  text-align: center;
  padding: 60px 20px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid var(--primary);
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

.skeleton-wrap {
  padding: 16px;
}

.team-table {
  background: var(--bg-card);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--border-color);
}

.table-header {
  display: grid;
  grid-template-columns: 60px 2fr 1fr 80px 120px 140px;
  gap: 12px;
  padding: 14px 20px;
  background: var(--bg-hover);
  border-bottom: 1px solid var(--border-color);
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
}

.table-row {
  display: grid;
  grid-template-columns: 60px 2fr 1fr 80px 120px 140px;
  gap: 12px;
  padding: 14px 20px;
  align-items: center;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  transition: background 0.2s;
}

.table-row:last-child {
  border-bottom: none;
}

.table-row:hover {
  background: var(--bg-hover);
}

.team-logo-small {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.team-logo-small img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.team-logo-small span {
  color: white;
  font-weight: 600;
  font-size: 14px;
}

.team-name-cell {
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 2px;
}

.team-desc-cell {
  font-size: 12px;
  color: var(--text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}

.col-leader, .col-members, .col-date {
  font-size: 13px;
  color: var(--text-secondary);
}

.col-members {
  font-weight: 500;
  color: var(--text-primary);
}

.col-actions {
  display: flex;
  gap: 6px;
  justify-content: flex-start;
}

.action-btn.small {
  padding: 6px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  min-width: 28px;
}

.action-btn.small svg {
  width: 14px;
  height: 14px;
}

/* 响应式优化 */
@media (max-width: 1200px) {
  .table-header,
  .table-row {
    grid-template-columns: 60px 2fr 1fr 70px 100px 120px;
  }
}

@media (max-width: 992px) {
  .table-header,
  .table-row {
    grid-template-columns: 50px 1.5fr 80px 60px 90px 110px;
  }
  
  .action-btn.small {
    padding: 4px;
    min-width: 24px;
  }
  
  .action-btn.small svg {
    width: 12px;
    height: 12px;
  }
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
  width: 64px;
  height: 64px;
  opacity: 0.5;
}

.empty-state p {
  margin: 0;
  font-size: 16px;
}

.empty-hint {
  font-size: 14px !important;
  color: var(--text-muted);
  opacity: 0.7;
}

.panel-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-left: auto;
}

.registration-control {
  display: flex;
  gap: 8px;
}

.reg-control-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s;
}

.reg-control-btn:hover {
  background: var(--bg-card-hover);
}

.reg-control-btn.active {
  background: var(--primary);
  border-color: var(--primary);
  color: white;
}

.reg-control-btn svg {
  width: 16px;
  height: 16px;
}

.view-toggle {
  display: flex;
  gap: 4px;
  background: var(--bg-hover);
  padding: 4px;
  border-radius: 8px;
}

.toggle-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s;
}

.toggle-btn:hover {
  background: var(--bg-card-hover);
}

.toggle-btn.active {
  background: var(--primary);
  color: white;
}

.toggle-btn svg {
  width: 18px;
  height: 18px;
}

.batch-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--primary);
  border-radius: 10px;
  margin-bottom: 16px;
}

.selected-count {
  color: white;
  font-size: 14px;
  font-weight: 500;
}

.batch-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 6px;
  color: white;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s;
}

.batch-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.batch-btn.approve {
  background: rgba(34, 197, 94, 0.8);
}

.batch-btn.approve:hover {
  background: rgba(34, 197, 94, 1);
}

.batch-btn.reject {
  background: rgba(239, 68, 68, 0.8);
}

.batch-btn.reject:hover {
  background: rgba(239, 68, 68, 1);
}

.batch-btn.clear {
  background: rgba(255, 255, 255, 0.1);
}

.batch-btn svg {
  width: 16px;
  height: 16px;
}

.user-table {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.compact-table-header,
.compact-table-row {
  display: grid;
  grid-template-columns: 40px 50px 100px 100px 80px 80px 100px 120px 140px;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: var(--bg-card);
}

.compact-table-header {
  background: var(--bg-hover);
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  border-radius: 8px 8px 0 0;
}

.compact-table-row {
  border: 1px solid var(--border-color);
  border-top: none;
  transition: all 0.3s;
}

.compact-table-row:last-child {
  border-radius: 0 0 8px 8px;
}

.compact-table-row:hover {
  background: var(--bg-card-hover);
}

.compact-table-row.selected {
  background: rgba(249, 115, 22, 0.1);
  border-color: var(--primary);
}

.col-check {
  display: flex;
  justify-content: center;
}

.col-check input {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.user-avatar-small {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  font-weight: bold;
}

.role-tag,
.status-tag {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 500;
}

.role-tag {
  background: rgba(59, 130, 246, 0.1);
  color: #60a5fa;
}

.role-tag.role-admin {
  background: rgba(249, 115, 22, 0.1);
  color: var(--primary);
}

.status-tag.status-pending {
  background: rgba(234, 179, 8, 0.1);
  color: #eab308;
}

.status-tag.status-active {
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.status-tag.status-deactivated {
  background: rgba(107, 114, 128, 0.1);
  color: #6b7280;
}

.team-tag-small {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 8px;
  font-size: 11px;
}

.team-tag-small.no-team {
  background: rgba(107, 114, 128, 0.1);
  color: #6b7280;
}

.action-btn-small {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: var(--bg-hover);
  border: none;
  border-radius: 6px;
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.3s;
}

.action-btn-small:hover {
  background: var(--bg-card-hover);
  color: var(--text-primary);
}

.action-btn-small.approve:hover {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.action-btn-small.reject:hover,
.action-btn-small.deactivate:hover,
.action-btn-small.danger:hover {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.action-btn-small.admin:hover {
  background: rgba(249, 115, 22, 0.2);
  color: var(--primary);
}

.action-btn-small.remove-admin:hover {
  background: rgba(107, 114, 128, 0.2);
  color: #6b7280;
}

.action-btn-small svg {
  width: 14px;
  height: 14px;
}

.user-select {
  display: flex;
  align-items: center;
}

.user-select input {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.user-card.selected {
  background: rgba(249, 115, 22, 0.1);
  border-color: var(--primary);
}

.user-cards {
  display: grid;
  gap: 16px;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 14px;
  transition: all 0.3s ease;
}

.user-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border-color: rgba(249, 115, 22, 0.3);
}

.user-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
  font-weight: bold;
  flex-shrink: 0;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.user-username {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.user-role {
  display: inline-block;
  padding: 2px 10px;
  background: rgba(59, 130, 246, 0.1);
  color: #60a5fa;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  margin-bottom: 4px;
}

.user-time {
  font-size: 12px;
  color: var(--text-muted);
}

.user-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.user-actions .action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.user-actions .action-btn svg {
  width: 16px;
  height: 16px;
}

.user-actions .action-btn.approve {
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.user-actions .action-btn.approve:hover {
  background: rgba(34, 197, 94, 0.2);
  transform: translateY(-2px);
}

.user-actions .action-btn.reject {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.user-actions .action-btn.reject:hover {
  background: rgba(239, 68, 68, 0.2);
  transform: translateY(-2px);
}

.user-actions .action-btn.admin-btn {
  background: rgba(168, 85, 247, 0.1);
  color: #a855f7;
}

.user-actions .action-btn.admin-btn:hover {
  background: rgba(168, 85, 247, 0.2);
  transform: translateY(-2px);
}

.user-actions .action-btn.remove-admin-btn {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

.user-actions .action-btn.remove-admin-btn:hover {
  background: rgba(245, 158, 11, 0.2);
  transform: translateY(-2px);
}

.user-avatar.admin {
  background: linear-gradient(135deg, #a855f7, #7c3aed);
  color: white;
}

.role-admin {
  color: #a855f7;
}

.super-admin-tag {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: white;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  margin-left: 4px;
}

.user-status {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  display: inline-block;
}

.user-status.status-pending {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

.user-status.status-active {
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.user-status.status-inactive {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.user-team {
  display: flex;
  align-items: center;
  gap: 6px;
}

.team-select {
  padding: 2px 8px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: white;
  font-size: 12px;
  cursor: pointer;
}

.team-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.team-tag.no-team {
  background: rgba(107, 114, 128, 0.1);
  color: #6b7280;
}

.team-list {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.team-section {
  background: var(--card-bg);
  border-radius: 12px;
  overflow: hidden;
}

.team-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-left: 4px solid;
  background: var(--card-bg);
}

.team-color {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  flex-shrink: 0;
}

.team-info {
  flex: 1;
  min-width: 0;
}

.team-name {
  font-weight: 600;
  font-size: 16px;
  margin-bottom: 4px;
}

.team-desc {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.team-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--text-secondary);
}

.team-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.member-select {
  padding: 6px 10px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--input-bg);
  font-size: 13px;
  cursor: pointer;
  max-width: 180px;
}

.team-members {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.02);
  border-top: 1px solid var(--border-color);
}

.member-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px 4px 4px;
  background: var(--card-bg);
  border: 1px solid;
  border-radius: 20px;
  font-size: 13px;
}

.member-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  font-weight: 600;
}

.member-name {
  font-weight: 500;
}

.member-username {
  color: var(--text-secondary);
  font-size: 12px;
}

.remove-member-btn {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: none;
  background: rgba(0, 0, 0, 0.1);
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 14px;
  line-height: 1;
  margin-left: 4px;
}

.remove-member-btn:hover {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.no-members {
  color: var(--text-secondary);
  font-size: 13px;
  padding: 8px 0;
}

.slide-panel {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.slide-panel.open {
  pointer-events: all;
  opacity: 1;
}

.panel-backdrop {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  animation: fadeIn 0.3s ease;
}

.panel-content {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 400px;
  max-width: 90vw;
  background: var(--bg-card);
  border-left: 1px solid var(--border-color);
  box-shadow: -4px 0 24px rgba(0, 0, 0, 0.15);
  overflow-y: auto;
  transform: translateX(100%);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-panel.open .panel-content {
  transform: translateX(0);
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.panel-header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-card);
  position: sticky;
  top: 0;
  z-index: 10;
}

.panel-header-bar h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.close-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 1px solid var(--border-color);
  background: var(--bg-elevated);
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  flex-shrink: 0;
}

.close-btn:hover {
  background: #ef4444;
  border-color: #ef4444;
  color: white;
}

.close-btn svg {
  width: 16px;
  height: 16px;
}

.team-detail-content {
  position: relative;
}

.team-detail-bg {
  height: 160px;
  background: linear-gradient(135deg, var(--theme-primary), var(--theme-secondary));
  background-size: cover;
  background-position: center;
  position: relative;
}

.team-bg-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to bottom, transparent, rgba(0, 0, 0, 0.6));
}

.team-detail-info {
  padding: 20px;
  text-align: center;
  margin-top: -50px;
  position: relative;
}

.team-detail-logo {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  margin: 0 auto 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: bold;
  color: white;
  border: 4px solid var(--bg-card);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.team-detail-logo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.team-detail-name {
  margin: 0 0 8px;
  font-size: 22px;
  font-weight: 600;
}

.team-detail-desc {
  margin: 0 0 12px;
  color: var(--text-secondary);
  font-size: 14px;
}

.team-detail-meta {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 13px;
  color: var(--text-muted);
}

.team-detail-section {
  padding: 16px 20px;
  border-top: 1px solid var(--border-color);
}

.team-detail-section h4 {
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
}

.team-detail-members {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-member-chip {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  background: var(--bg-elevated);
  border-radius: 8px;
}

.detail-member-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  font-weight: 600;
}

.detail-member-name {
  font-weight: 500;
  font-size: 14px;
}

.detail-member-username {
  color: var(--text-secondary);
  font-size: 12px;
  margin-left: auto;
}

.team-detail-actions {
  padding: 20px;
  border-top: 1px solid var(--border-color);
  display: flex;
  gap: 12px;
}

.clickable {
  cursor: pointer;
  color: var(--text-primary);
  transition: color 0.2s;
}

.clickable:hover {
  color: var(--theme-primary);
}

.primary-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 20px;
  background: linear-gradient(135deg, var(--theme-primary), var(--theme-secondary));
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.primary-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}

.primary-btn svg {
  width: 18px;
  height: 18px;
}

.danger-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 20px;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.danger-btn:hover {
  background: #dc2626;
  transform: translateY(-2px);
}

.danger-btn svg {
  width: 18px;
  height: 18px;
}

.team-section {
  cursor: pointer;
}

.category-icon-preview {
  font-size: 20px;
}

.color-preview {
  display: inline-block;
  width: 24px;
  height: 24px;
  border-radius: 6px;
  margin-right: 8px;
  vertical-align: middle;
  border: 1px solid var(--border-color);
}

.category-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-field label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
}

.color-input-wrapper {
  display: flex;
  gap: 12px;
  align-items: center;
}

.color-input-wrapper input[type="color"] {
  width: 60px;
  padding: 4px;
  cursor: pointer;
}

@media (max-width: 768px) {
  .export-cards {
    grid-template-columns: 1fr;
  }

  .table-header,
  .table-row {
    grid-template-columns: 60px 1fr 100px;
  }

  .table-header span:nth-child(2),
  .table-header span:nth-child(4),
  .table-row span:nth-child(2),
  .table-row span:nth-child(4) {
    display: none;
  }
}

.form-row.logo-row {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.form-row.logo-row .form-field {
  flex: 1;
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

.modal-body textarea {
  width: 100%;
  min-height: 80px;
  background: var(--bg-dark);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 12px 16px;
  font-size: 15px;
  color: var(--text-primary);
  resize: vertical;
  transition: all 0.3s ease;
}

.modal-body textarea:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.15);
}

.modal-body select {
  width: 100%;
  height: 48px;
  background: var(--bg-dark);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 0 16px;
  font-size: 15px;
  color: var(--text-primary);
  transition: all 0.3s ease;
  cursor: pointer;
}

.modal-body select:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.15);
}

.error-message {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 10px;
  color: #ef4444;
  font-size: 14px;
  margin-bottom: 16px;
}

.error-message svg {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}
</style>