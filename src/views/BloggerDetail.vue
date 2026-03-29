<template>
  <div class="detail-page" v-if="blogger">
    <div class="page-header">
      <div class="header-left">
        <button class="back-btn" @click="router.push('/')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
            <polyline points="9 22 9 12 15 12 15 22"/>
          </svg>
          返回主页
        </button>
      </div>
      <div class="header-right">
        <button class="edit-btn" @click="startEdit" v-if="canEdit && !showEditDialog">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
          </svg>
          编辑信息
        </button>
        <button class="transfer-btn" @click="showTransferDialog = true" v-if="canEdit && !showEditDialog">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 1l4 4-4 4"/>
            <path d="M3 11V9a4 4 0 0 1 4-4h14"/>
            <path d="M7 23l-4-4 4-4"/>
            <path d="M21 13v2a4 4 0 0 1-4 4H3"/>
          </svg>
          转移归属
        </button>
        <button class="delete-btn" @click="confirmDelete" v-if="canDelete">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3,6 5,6 21,6"/>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
          </svg>
          删除博主
        </button>
      </div>
    </div>

    <div class="detail-content">
      <div class="main-info">
        <div class="profile-card">
          <div class="avatar" :class="{ 'has-image': blogger.avatar }">
            <img v-if="blogger.avatar" :src="blogger.avatar" alt="博主头像" v-avatar />
            <template v-else>
              {{ blogger.nickname?.charAt(0) || '?' }}
            </template>
          </div>
          <h1>{{ blogger.nickname }}</h1>
          <div class="tags-display" v-if="blogger.tags && blogger.tags.length">
            <span v-for="(tag, index) in blogger.tags" :key="index" class="tag">
              {{ tag }}
            </span>
          </div>
          <div class="meta-info">
            <span class="meta-item category-badge" :style="{ backgroundColor: getCategoryColor(blogger.category), color: getTextColorForBackground(getCategoryColor(blogger.category)) }">
              <span class="category-icon">{{ getCategoryIcon(blogger.category) }}</span>
              <span class="category-text">{{ blogger.category }}</span>
            </span>
          </div>
        </div>

        <div class="info-card">
          <h3>基本信息</h3>
          <div class="info-grid">
            <div class="info-item">
              <label>对接产品</label>
              <span>{{ blogger.product || '-' }}</span>
            </div>
            <div class="info-item">
              <label>联系方式</label>
              <div class="contact-with-copy">
                <span>{{ blogger.contact || '-' }}</span>
                <button v-if="blogger.contact" class="copy-btn" @click="copyToClipboard(blogger.contact)" title="复制联系方式">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                  </svg>
                </button>
              </div>
            </div>
            <div class="info-item">
              <label>微信号</label>
              <div class="contact-with-copy">
                <span>{{ blogger.wechat || '-' }}</span>
                <button v-if="blogger.wechat" class="copy-btn" @click="copyToClipboard(blogger.wechat)" title="复制微信号">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                  </svg>
                </button>
              </div>
            </div>
            <div class="info-item" v-if="blogger.custom_contact">
              <label>自定义联系方式</label>
              <div class="contact-with-copy">
                <span>{{ blogger.custom_contact }}</span>
                <button v-if="blogger.custom_contact" class="copy-btn" @click="copyToClipboard(blogger.custom_contact)" title="复制联系方式">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                  </svg>
                </button>
              </div>
            </div>
            <div class="info-item">
              <label>活跃平台</label>
              <span>{{ blogger.platform || '-' }}</span>
            </div>
            <div class="info-item">
              <label>跟进状态</label>
              <span class="status-tag" :class="getStatusClass(blogger.status)">{{ blogger.status || '初次联系' }}</span>
            </div>
            <div class="info-item">
              <label>录入时间</label>
              <span>{{ formatDateTime(blogger.create_time) }}</span>
            </div>
          </div>
          <div class="info-item full-width">
            <label>简介</label>
            <span class="description">{{ blogger.description || '-' }}</span>
          </div>
        </div>

        <!-- 对接人信息卡片 -->
        <div class="info-card contact-person-card" v-if="blogger.owner_username || blogger.owner_real_name">
          <div v-if="blogger.owner_username" class="contact-person-link" @click="goToUserProfile(blogger.owner_username)">
            <h3>👤 对接人</h3>
            <div class="contact-person-content">
              <div class="contact-person-avatar">
                <img v-if="blogger.owner_avatar" :src="blogger.owner_avatar" alt="头像" v-avatar />
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
              </div>
              <div class="contact-person-info">
                <div class="contact-person-name">{{ blogger.owner_real_name || blogger.owner_username }}</div>
                <div class="contact-person-username">{{ blogger.owner_username }}</div>
                <div class="contact-person-team">小组: {{ bloggerTeamName }}</div>
              </div>
              <button class="jump-profile-btn" @click.stop="goToUserProfile(blogger.owner_username)" title="查看个人资料">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                  <polyline points="15 3 21 3 21 9"/>
                  <line x1="10" y1="14" x2="21" y2="3"/>
                </svg>
              </button>
            </div>
          </div>
          <div v-else>
            <h3>👤 对接人</h3>
            <div class="contact-person-content">
              <div class="contact-person-avatar">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
              </div>
              <div class="contact-person-info">
                <div class="contact-person-name">{{ blogger.owner_real_name }}</div>
                <div class="contact-person-team">小组: {{ bloggerTeamName }}</div>
              </div>
            </div>
          </div>
        </div>

        <div class="info-card" v-if="blogger.tags && blogger.tags.length">
          <h3>标签 Tags</h3>
          <div class="tags-cloud">
            <span v-for="(tag, index) in blogger.tags" :key="index" class="cloud-tag" @click="searchByTag(tag)">
              {{ tag }}
            </span>
          </div>
        </div>

        <div class="info-card">
          <div class="card-header-with-action">
            <h3>📝 跟进记录</h3>
            <button class="refresh-btn" @click="loadFollowup(blogger.id)" :disabled="followupLoading">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ 'spinning': followupLoading }">
                <path d="M23 4v6h-6"/>
                <path d="M1 20v-6h6"/>
                <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10"/>
                <path d="M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
              </svg>
            </button>
          </div>
          <div class="followup-form" v-if="canEdit">
            <textarea
              v-model="followupContent"
              placeholder="添加跟进记录..."
              rows="3"
            ></textarea>
            <div class="followup-actions">
              <select v-model="followupType" class="type-select">
                <option value="跟进">跟进</option>
                <option value="电话">电话</option>
                <option value="微信">微信</option>
                <option value="邮件">邮件</option>
                <option value="会面">会面</option>
              </select>
              <button class="add-btn" @click="addFollowup" :disabled="!followupContent.trim()">
                添加记录
              </button>
            </div>
          </div>

          <div class="followup-list" v-if="followupList.length">
            <div v-for="item in followupList" :key="item.id" class="followup-item">
              <div class="followup-header">
                <span class="followup-type" :class="item.type">{{ item.type }}</span>
                <span class="followup-time">{{ formatDateTime(item.create_time) }}</span>
              </div>
              <div v-if="editingFollowupId === item.id" class="followup-edit">
                <select v-model="editFollowupType" class="type-select">
                  <option value="跟进">跟进</option>
                  <option value="电话">电话</option>
                  <option value="微信">微信</option>
                  <option value="邮件">邮件</option>
                  <option value="会面">会面</option>
                </select>
                <textarea v-model="editFollowupContent" class="edit-textarea" rows="2"></textarea>
                <div class="edit-actions">
                  <button class="save-btn" @click="saveFollowupEdit(item.id)">保存</button>
                  <button class="cancel-btn" @click="cancelFollowupEdit">取消</button>
                </div>
              </div>
              <div v-else class="followup-content-text">{{ item.content }}</div>
              <div class="followup-footer">
                <span class="followup-operator">{{ item.operator }}</span>
                <div class="followup-ops" v-if="canEditFollowup(item)">
                  <button class="op-btn" @click="startFollowupEdit(item)">编辑</button>
                  <button class="delete-btn" @click="deleteFollowup(item.id)">删除</button>
                </div>
              </div>
            </div>
          </div>
          <div class="empty-state" v-else>
            <span>暂无跟进记录</span>
          </div>
        </div>

        <div class="info-card">
          <h3>📋 操作日志</h3>
          <div class="log-list" v-if="operationLogs.length">
            <div v-for="log in operationLogs" :key="log.id" class="log-item">
              <span class="log-action">{{ log.action }}</span>
              <span class="log-detail">{{ log.detail }}</span>
              <span class="log-operator">{{ log.operator }}</span>
              <span class="log-time">{{ formatDateTime(log.create_time) }}</span>
            </div>
          </div>
          <div class="empty-state" v-else>
            <span>暂无操作日志</span>
          </div>
        </div>

        <div class="info-card">
          <div class="card-header-with-action">
            <h3>🤝 合作记录</h3>
            <button class="add-btn" @click="openAddCooperation">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              添加记录
            </button>
          </div>
          <div class="timeline" v-if="cooperationHistory.length">
            <div v-for="(item, index) in cooperationHistory" :key="item.id" class="timeline-item">
              <div class="timeline-dot"></div>
              <div class="timeline-content">
                <div class="timeline-header">
                  <span class="timeline-title">{{ item.title }}</span>
                  <span class="timeline-date">{{ formatDate(item.date) }}</span>
                </div>
                <div class="timeline-body">
                  <div class="timeline-info">
                    <span class="info-label">状态：</span>
                    <span class="info-value" :class="getCooperationStatusClass(item.status)">{{ item.status }}</span>
                  </div>
                  <div class="timeline-info" v-if="item.product">
                    <span class="info-label">产品：</span>
                    <span class="info-value">{{ item.product }}</span>
                  </div>
                  <div class="timeline-info" v-if="item.amount">
                    <span class="info-label">金额：</span>
                    <span class="info-value">¥{{ item.amount }}</span>
                  </div>
                  <div class="timeline-note" v-if="item.note">
                    <span class="note-label">备注：</span>
                    <span class="note-text">{{ item.note }}</span>
                  </div>
                </div>
                <div class="timeline-actions" v-if="canEdit">
                  <button class="op-btn" @click="openEditCooperation(item)">编辑</button>
                  <button class="delete-btn" @click="deleteCooperation(item.id)">删除</button>
                </div>
              </div>
            </div>
          </div>
          <div class="empty-state" v-else>
            <span>暂无合作记录</span>
          </div>
        </div>
      </div>

      <div class="side-info">
        <div class="info-card countdown-card" v-if="showCountdown">
          <h3>⏰ 联系方式倒计时</h3>
          <div class="countdown-display">
            <div class="countdown-value" :class="countdownClass">
              <span class="countdown-number">{{ countdownDays }}</span>
              <span class="countdown-unit">天</span>
            </div>
            <div class="countdown-text" :class="countdownClass">
              {{ countdownText }}
            </div>
          </div>
          <div class="countdown-tip">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <path d="M12 6v6l4 2"/>
            </svg>
            请及时补充联系方式
          </div>
        </div>

        <div class="info-card">
          <h3>🔍 快速搜索</h3>
          
          <div class="search-section">
            <div class="search-select-group">
              <label class="search-label">搜索引擎</label>
              <select v-model="selectedSearchEngine" class="form-select search-select">
                <option v-for="engine in searchEngines" :key="engine.id" :value="engine.id">
                  {{ engine.name }}
                </option>
              </select>
            </div>
            
            <div class="search-prompt-group">
              <label class="search-label">搜索关键词</label>
              <div class="prompt-buttons">
                <button 
                  v-for="prompt in searchPrompts" 
                  :key="prompt.id"
                  class="prompt-btn"
                  :class="{ active: selectedPrompt === prompt.id }"
                  @click="selectedPrompt = prompt.id"
                >
                  {{ prompt.name }}
                </button>
              </div>
            </div>
            
            <div class="search-preview">
              <label class="search-label">搜索内容</label>
              <input 
                type="text" 
                v-model="customSearchQuery"
                class="search-preview-input"
                placeholder="输入或编辑搜索内容..."
              />
            </div>
            
            <button class="search-btn" @click="openSearch">
              开始搜索
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showTransferDialog" class="dialog-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click="showTransferDialog = false">
      <div class="dialog" style="background:var(--bg-card);" @click.stop>
        <div class="dialog-header">
          <h3>转移博主归属</h3>
          <button class="close-btn" @click="showTransferDialog = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label>选择接收人</label>
            <select v-model="transferForm.toUserId" class="form-select">
              <option value="">请选择团队成员</option>
              <option v-for="user in teamMembers" :key="user.id" :value="user.id">
                {{ user.real_name }} ({{ user.username }})
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>转移原因</label>
            <textarea v-model="transferForm.reason" class="form-textarea" rows="3" placeholder="请填写转移原因..."></textarea>
          </div>
          <div class="dialog-actions">
            <button class="btn-secondary" @click="showTransferDialog = false">取消</button>
            <button class="btn-primary" @click="handleTransfer" :disabled="!transferForm.toUserId">确认转移</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 合作记录编辑弹窗 -->
    <div class="dialog-overlay" v-if="showCooperationDialog" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="showCooperationDialog = false">
      <div class="dialog cooperation-dialog" style="background:var(--bg-card);">
        <div class="dialog-header">
          <h3>{{ editingCooperationId ? '编辑合作记录' : '添加合作记录' }}</h3>
          <button class="close-btn" @click="showCooperationDialog = false">×</button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label>合作标题 *</label>
            <input type="text" v-model="cooperationForm.title" placeholder="请输入合作标题" class="form-input" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>日期 *</label>
              <input type="date" v-model="cooperationForm.date" class="form-input" />
            </div>
            <div class="form-group">
              <label>状态 *</label>
              <select v-model="cooperationForm.status" class="form-input">
                <option v-for="status in cooperationStatuses" :key="status" :value="status">{{ status }}</option>
              </select>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>产品</label>
              <input type="text" v-model="cooperationForm.product" placeholder="请输入合作产品" class="form-input" />
            </div>
            <div class="form-group">
              <label>金额 (元)</label>
              <input type="number" v-model="cooperationForm.amount" placeholder="请输入金额" class="form-input" min="0" step="0.01" />
            </div>
          </div>
          <div class="form-group">
            <label>备注</label>
            <textarea v-model="cooperationForm.note" placeholder="请输入备注信息" rows="3" class="form-input"></textarea>
          </div>
        </div>
        <div class="dialog-actions">
          <button class="btn-secondary" @click="showCooperationDialog = false">取消</button>
          <button class="btn-primary" @click="saveCooperation">
            {{ editingCooperationId ? '保存修改' : '添加记录' }}
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- 编辑博主信息弹窗 -->
  <div class="dialog-overlay" v-if="showEditDialog" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click="closeEditDialog">
    <div class="edit-dialog" style="background:var(--bg-card);" @click.stop>
      <div class="dialog-header">
        <h2>编辑博主信息</h2>
        <button class="close-btn" @click="closeEditDialog">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      <div class="dialog-body">
        <div class="edit-form">
          <!-- 头像上传 -->
          <div class="form-section">
            <label class="form-label">博主头像</label>
            <div class="avatar-upload">
              <div class="avatar-preview" :class="{ 'has-image': editForm.avatar || blogger.avatar }" @click="triggerAvatarUpload()">
                <img v-if="editForm.avatar || blogger.avatar" :src="editForm.avatar || blogger.avatar" alt="头像" v-avatar />
                <template v-else>
                  {{ (editForm.nickname || blogger.nickname)?.charAt(0) || '?' }}
                </template>
                <div class="avatar-overlay">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/>
                    <circle cx="12" cy="13" r="4"/>
                  </svg>
                  <span>更换头像</span>
                </div>
              </div>
              <input
                type="file"
                ref="avatarInput"
                accept="image/*"
                style="display: none"
                @change="handleAvatarChange"
              />
              <p class="upload-tip">点击头像上传或更换</p>
            </div>
          </div>

          <!-- 基本信息 -->
          <div class="form-section">
            <label class="form-label">昵称 *</label>
            <input v-model="editForm.nickname" class="form-input" placeholder="请输入博主昵称" />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">活跃平台 *</label>
              <input v-model="editForm.platform" class="form-input" placeholder="如：抖音、小红书" />
            </div>
            <div class="form-group">
              <label class="form-label">分类 *</label>
              <select v-model="editForm.category" class="form-select">
                <option value="">请选择分类</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
              </select>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">对接产品</label>
              <input v-model="editForm.product" class="form-input" placeholder="请输入对接产品" />
            </div>
            <div class="form-group">
              <label class="form-label">联系方式</label>
              <input v-model="editForm.contact" class="form-input" placeholder="请输入联系方式" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">微信号</label>
              <input v-model="editForm.wechat" class="form-input" placeholder="请输入微信号" />
            </div>
            <div class="form-group">
              <label class="form-label">自定义联系方式</label>
              <input v-model="editForm.custom_contact" class="form-input" placeholder="请输入其他联系方式" />
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">跟进状态</label>
              <select v-model="editForm.status" class="form-select">
                <option v-for="s in bloggerStatusList" :key="s.id" :value="s.id">{{ s.name }}</option>
              </select>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">状态备注</label>
              <input v-model="editForm.status_remark" class="form-input" placeholder="请输入状态变更备注" />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">对接人</label>
            <input v-model="editForm.contact_person" class="form-input" placeholder="请输入对接人姓名" />
          </div>

          <div class="form-group">
            <label class="form-label">博主简介</label>
            <textarea v-model="editForm.description" class="form-textarea" rows="3" placeholder="请输入博主简介"></textarea>
          </div>

          <!-- 标签编辑 -->
          <div class="form-group">
            <label class="form-label">标签</label>
            <div class="tags-edit-container">
              <div class="tags-display" v-if="editForm.tags && editForm.tags.length">
                <span v-for="(tag, index) in editForm.tags" :key="index" class="tag">
                  {{ tag }}
                  <button class="remove-tag" @click="removeTag(index)">×</button>
                </span>
              </div>
              <div class="tag-input-wrapper">
                <input 
                  v-model="tagInput" 
                  placeholder="输入标签后回车添加" 
                  @keydown.enter.prevent="addTag" 
                  class="form-input"
                />
                <button class="add-tag-btn" @click="addTag">添加</button>
              </div>
            </div>
          </div>
        </div>

        <div class="dialog-actions">
          <button class="btn-secondary" @click="closeEditDialog">取消</button>
          <button class="btn-primary" @click="saveEdit" :disabled="!editForm.nickname || !editForm.platform">
            保存修改
          </button>
        </div>
      </div>
    </div>
  </div>

  <div class="detail-page" v-else>
    <div class="loading">加载中...</div>
  </div>

  <AvatarCropper
    v-model:visible="cropDialogVisible"
    :image-file="selectedImageFile"
    @success="handleAvatarUploadSuccess"
  />
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getBloggerDetail, followupAddAPI, followupListAPI, followupDeleteAPI, followupUpdateAPI, bloggerUpdateAPI, bloggerDeleteAPI, categoryListAPI, requestBloggerTransfer, cooperationListAPI, cooperationAddAPI, cooperationUpdateAPI, cooperationDeleteAPI, getBloggerStatusList, getTeamsAPI } from '../api'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'
import { useConfirm } from '../utils/confirm'
import AvatarCropper from '../components/AvatarCropper.vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const notification = useNotification()
const { confirmDanger } = useConfirm()

const blogger = ref(null)
const followupList = ref([])
const followupLoading = ref(false)
const operationLogs = ref([])
const cooperationHistory = ref([])
const showTransferDialog = ref(false)
const transferForm = ref({
  toUserId: '',
  reason: ''
})
const teamMembers = ref([])
const followupContent = ref('')
const followupType = ref('跟进')
const showEditDialog = ref(false)
const editForm = ref({})
const tagInput = ref('')
const categories = ref([])
const teams = ref([])
const editingFollowupId = ref(null)
const editFollowupContent = ref('')
const editFollowupType = ref('跟进')
const avatarInput = ref(null)
const showCooperationDialog = ref(false)
const editingCooperationId = ref(null)
const cooperationForm = ref({
  title: '',
  date: new Date().toISOString().split('T')[0],
  status: '进行中',
  product: '',
  amount: '',
  note: ''
})
const cooperationStatuses = ['洽谈中', '进行中', '已完成', '已终止']
const bloggerStatusList = ref([])

let countdownTimer = null
const countdownDays = ref(15)
const countdownText = ref('剩余时间')

const showCountdown = computed(() => {
  if (!blogger.value) return false
  if (blogger.value.contact || blogger.value.wechat || blogger.value.custom_contact) return false
  return true
})

const countdownClass = computed(() => {
  if (countdownDays.value <= 3) return 'urgent'
  if (countdownDays.value <= 7) return 'warning'
  return 'normal'
})

const calculateCountdown = () => {
  if (!blogger.value?.create_time) return
  
  const createTime = new Date(blogger.value.create_time)
  const now = new Date()
  const diffTime = createTime.getTime() + 15 * 24 * 60 * 60 * 1000 - now.getTime()
  
  if (diffTime <= 0) {
    countdownDays.value = 0
    countdownText.value = '已过期'
    return
  }
  
  countdownDays.value = Math.ceil(diffTime / (24 * 60 * 60 * 1000))
  
  if (countdownDays.value === 1) {
    countdownText.value = '最后1天'
  } else if (countdownDays.value <= 3) {
    countdownText.value = '即将过期'
  } else if (countdownDays.value <= 7) {
    countdownText.value = '请注意时间'
  } else {
    countdownText.value = '剩余时间'
  }
}

watch([() => blogger.value?.contact, () => blogger.value?.wechat, () => blogger.value?.custom_contact], () => {
  if (blogger.value?.contact || blogger.value?.wechat || blogger.value?.custom_contact) {
    if (countdownTimer) {
      clearInterval(countdownTimer)
      countdownTimer = null
    }
  }
})

watch(() => blogger.value, () => {
  if (blogger.value) {
    calculateCountdown()
  }
}, { immediate: true })

const cropDialogVisible = ref(false)
const selectedImageFile = ref(null)

const handleAvatarChange = (event) => {
  const file = event.target.files[0]
  if (!file) return

  if (file.size > 10 * 1024 * 1024) {
    notification.warning('文件大小不能超过10MB')
    return
  }

  selectedImageFile.value = file
  cropDialogVisible.value = true

  event.target.value = ''
}

const handleAvatarUploadSuccess = (url) => {
  editForm.value.avatar = url
}

const searchEngines = [
  { id: 'baidu', name: '百度', url: 'https://www.baidu.com/s?wd=' },
  { id: 'google', name: 'Google', url: 'https://www.google.com/search?q=' },
  { id: 'bing', name: 'Bing', url: 'https://www.bing.com/search?q=' },
  { id: 'sogou', name: '搜狗', url: 'https://www.sogou.com/web?query=' },
  { id: '360', name: '360搜索', url: 'https://www.so.com/s?q=' }
]

const searchPrompts = [
  { id: 'basic', name: '基本信息', template: '{nickname}' },
  { id: 'news', name: '最新资讯', template: '{nickname} 最新动态' },
  { id: 'works', name: '作品展示', template: '{nickname} 作品合集' },
  { id: 'collab', name: '合作案例', template: '{nickname} 合作品牌' },
  { id: 'price', name: '报价参考', template: '{nickname} 报价 费用' },
  { id: 'fans', name: '粉丝分析', template: '{nickname} 粉丝画像' },
  { id: 'platform', name: '平台数据', template: '{nickname} 粉丝数 播放量' },
  { id: 'interview', name: '采访报道', template: '{nickname} 专访 采访' },
  { id: 'controversy', name: '争议舆情', template: '{nickname} 黑料 负面新闻' },
  { id: 'similar', name: '同类博主', template: '{nickname} 同类博主 竞品博主' },
  { id: 'agency', name: '签约机构', template: '{nickname} MCN 签约公司' },
  { id: 'history', name: '发展历程', template: '{nickname} 成长史 成名经历' }
]

const selectedSearchEngine = ref('baidu')
const selectedPrompt = ref('basic')
const customSearchQuery = ref('')

const updateSearchQuery = () => {
  if (!blogger.value) return
  const prompt = searchPrompts.find(p => p.id === selectedPrompt.value)
  if (prompt) {
    customSearchQuery.value = prompt.template.replace('{nickname}', blogger.value.nickname)
  }
}

watch(selectedPrompt, () => {
  updateSearchQuery()
})

watch(blogger, () => {
  updateSearchQuery()
}, { immediate: true })

const openSearch = () => {
  const query = customSearchQuery.value
  
  if (!query) {
    return
  }
  
  const engine = searchEngines.find(e => e.id === selectedSearchEngine.value)
  
  if (engine) {
    const searchUrl = engine.url + encodeURIComponent(query)
    window.open(searchUrl, '_blank')
  }
}

const getStatusClass = (status) => {
  const classMap = {
    '初次联系': 'status-init',
    '洽谈中': 'status-talking',
    '已合作': 'status-cooperated',
    '已拒绝': 'status-rejected',
    '暂停合作': 'status-paused'
  }
  return classMap[status] || 'status-init'
}

const canEdit = computed(() => {
  return userStore.role === 'admin' || blogger.value?.user_name === userStore.username
})

const canDelete = computed(() => {
  return userStore.role === 'admin' || blogger.value?.user_name === userStore.username
})

const bloggerTeamName = computed(() => {
  if (!blogger.value?.team_id) return '无小组'
  const team = teams.value.find(t => t.id === blogger.value.team_id)
  return team?.name || '无小组'
})

const canEditFollowup = (item) => {
  return userStore.role === 'admin' || item.operator === userStore.real_name
}

const loadBlogger = async () => {
  try {
    const res = await getBloggerDetail(route.params.id)
    if (res.code === 200) {
      blogger.value = res.data
      if (blogger.value) {
        loadFollowup(blogger.value.id)
        loadOperationLogs(blogger.value.nickname)
        loadCooperationHistory()
      }
    } else {
      console.error('加载博主信息失败', res.message)
    }
  } catch (error) {
    console.error('加载博主信息失败', error)
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

const getCategoryColor = (categoryName) => {
  const category = categories.value.find(c => c.name === categoryName)
  if (category && category.color) {
    return category.color
  }
  const colorMap = {
    '美妆': '#FF6B6B',
    '美食': '#4ECDC4',
    '穿搭': '#45B7D1',
    '旅游': '#96CEB4',
    '科技': '#FFEAA7',
    '其他': '#DFE6E9'
  }
  return colorMap[categoryName] || '#74B9FF'
}

const getCategoryIcon = (categoryName) => {
  const category = categories.value.find(c => c.name === categoryName)
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

const loadFollowup = async (bloggerId) => {
  followupLoading.value = true
  try {
    const res = await followupListAPI({ blogger_id: bloggerId })
    if (res.code === 200) {
      followupList.value = res.data || []
    }
  } catch (error) {
    console.error('加载跟进记录失败', error)
  } finally {
    followupLoading.value = false
  }
}

const loadOperationLogs = async (nickname) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/blogger/logs?nickname=${encodeURIComponent(nickname)}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    if (res.ok) {
      const data = await res.json()
      if (data.code === 200) {
        operationLogs.value = data.data || []
      }
    }
  } catch (error) {
    console.error('操作日志加载失败', error)
  }
}

const loadCooperationHistory = async () => {
  try {
    const res = await cooperationListAPI({ blogger_id: blogger.value?.id })
    if (res.code === 200) {
      cooperationHistory.value = res.data || []
    }
  } catch (error) {
    console.error('加载合作记录失败', error)
  }
}

const openAddCooperation = () => {
  editingCooperationId.value = null
  cooperationForm.value = {
    title: '',
    date: new Date().toISOString().split('T')[0],
    status: '进行中',
    product: '',
    amount: '',
    note: ''
  }
  showCooperationDialog.value = true
}

const openEditCooperation = (item) => {
  editingCooperationId.value = item.id
  cooperationForm.value = {
    title: item.title,
    date: item.date,
    status: item.status,
    product: item.product || '',
    amount: item.amount ? String(item.amount) : '',
    note: item.note || ''
  }
  showCooperationDialog.value = true
}

const saveCooperation = async () => {
  if (!cooperationForm.value.title || !cooperationForm.value.date || !cooperationForm.value.status) {
    notification.warning('请填写完整信息')
    return
  }

  try {
    const data = {
      blogger_id: blogger.value.id,
      ...cooperationForm.value,
      amount: cooperationForm.value.amount ? Number(cooperationForm.value.amount) : 0
    }

    if (editingCooperationId.value) {
      await cooperationUpdateAPI(editingCooperationId.value, data)
      notification.success('更新成功')
    } else {
      await cooperationAddAPI(data)
      notification.success('添加成功')
    }

    showCooperationDialog.value = false
    loadCooperationHistory()
  } catch (error) {
    console.error('保存合作记录失败', error)
    notification.error('保存失败')
  }
}

const deleteCooperation = async (id) => {
  if (!await confirmDanger('确定要删除这条合作记录吗？')) return

  try {
    await cooperationDeleteAPI(id)
    notification.success('删除成功')
    loadCooperationHistory()
  } catch (error) {
    console.error('删除合作记录失败', error)
    notification.error('删除失败')
  }
}

const loadTeamMembers = async () => {
  try {
    const res = await fetch('/api/users/public')
    if (res.ok) {
      const data = await res.json()
      if (data.code === 200) {
        teamMembers.value = data.data.filter(u => u.status === 'active' && u.id !== userStore.id)
      }
    }
  } catch (error) {
    console.error('加载团队成员失败', error)
  }
}

const confirmDelete = () => {
  if (!confirmDanger('确定要删除该博主吗？此操作不可恢复。')) return
  handleDelete()
}

const handleDelete = async () => {
  try {
    const res = await bloggerDeleteAPI(blogger.value.id)
    if (res.code === 200) {
      notification.success('删除成功')
      router.push('/')
    } else {
      notification.error(res.message || '删除失败')
    }
  } catch (error) {
    notification.error('删除失败')
  }
}

const handleTransfer = async () => {
  if (!transferForm.value.toUserId) {
    notification.warning('请选择接收人')
    return
  }

  try {
    const res = await requestBloggerTransfer({
      bloggerId: blogger.value.id,
      toUserId: transferForm.value.toUserId,
      reason: transferForm.value.reason
    })
    if (res.code === 200) {
      notification.success('转移请求已发送，等待对方确认')
      showTransferDialog.value = false
      transferForm.value = { toUserId: '', reason: '' }
    } else {
      notification.error(res.message || '转移失败')
    }
  } catch (error) {
    notification.error('转移失败')
  }
}

const getCooperationStatusClass = (status) => {
  const statusMap = {
    '已完成': 'status-completed',
    '进行中': 'status-progress',
    '已取消': 'status-cancelled',
    '待确认': 'status-pending'
  }
  return statusMap[status] || ''
}

const startEdit = () => {
  editForm.value = {
    id: blogger.value.id,
    nickname: blogger.value.nickname,
    category: blogger.value.category,
    product: blogger.value.product,
    contact: blogger.value.contact,
    wechat: blogger.value.wechat || '',
    custom_contact: blogger.value.custom_contact || '',
    contact_person: blogger.value.contact_person || '',
    platform: blogger.value.platform,
    description: blogger.value.description,
    tags: [...(blogger.value.tags || [])],
    avatar: blogger.value.avatar || '',
    status: blogger.value.status || '初次联系',
    status_remark: ''
  }
  showEditDialog.value = true
}

const closeEditDialog = () => {
  showEditDialog.value = false
}

const triggerAvatarUpload = () => {
  avatarInput.value?.click()
}

const addTag = () => {
  const tag = tagInput.value.trim()
  if (tag && !editForm.value.tags.includes(tag)) {
    editForm.value.tags.push(tag)
    tagInput.value = ''
  }
}

const saveEdit = async () => {
  try {
    const res = await bloggerUpdateAPI(blogger.value.id, editForm.value)
    
    if (res.code === 200) {
      await loadBlogger()
      showEditDialog.value = false
      loadOperationLogs(blogger.value.nickname)
      notification.success('保存成功')
    } else {
      notification.error(res.message || '保存失败')
    }
  } catch (error) {
    console.error('保存失败:', error)
    notification.error('保存失败')
  }
}

const addFollowup = async () => {
  if (!followupContent.value.trim()) return

  try {
    const res = await followupAddAPI({
      blogger_id: blogger.value.id,
      content: followupContent.value,
      type: followupType.value
    })
    if (res.code === 200) {
      followupList.value.unshift(res.data)
      followupContent.value = ''
      loadOperationLogs(blogger.value.nickname)
      notification.success('添加成功')
    } else {
      notification.error(res.message || '添加失败')
    }
  } catch (error) {
    notification.error('添加失败')
  }
}

const startFollowupEdit = (item) => {
  editingFollowupId.value = item.id
  editFollowupContent.value = item.content
  editFollowupType.value = item.type
}

const cancelFollowupEdit = () => {
  editingFollowupId.value = null
  editFollowupContent.value = ''
  editFollowupType.value = '跟进'
}

const saveFollowupEdit = async (id) => {
  try {
    const res = await followupUpdateAPI(id, { 
      content: editFollowupContent.value,
      type: editFollowupType.value
    })
    if (res.code === 200) {
      const index = followupList.value.findIndex(f => f.id === id)
      if (index !== -1) {
        followupList.value[index].content = editFollowupContent.value
        followupList.value[index].type = editFollowupType.value
      }
      editingFollowupId.value = null
      editFollowupContent.value = ''
      editFollowupType.value = '跟进'
      loadOperationLogs(blogger.value.nickname)
      notification.success('保存成功')
    } else {
      notification.error(res.message || '保存失败')
    }
  } catch (error) {
    notification.error('保存失败')
  }
}

const deleteFollowup = async (id) => {
  if (!await confirmDanger('确定要删除这条跟进记录吗？')) return

  try {
    const res = await followupDeleteAPI(id)
    if (res.code === 200) {
      followupList.value = followupList.value.filter(f => f.id !== id)
      loadOperationLogs(blogger.value.nickname)
      notification.success('删除成功')
    } else {
      notification.error(res.message || '删除失败')
    }
  } catch (error) {
    notification.error('删除失败')
  }
}

const searchByTag = (tag) => {
  router.push({ path: '/', query: { tag } })
}

const goToUserProfile = (username) => {
  router.push(`/user/${username}`)
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    notification.success('已复制到剪贴板')
  } catch (error) {
    const textarea = document.createElement('textarea')
    textarea.value = text
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    notification.success('已复制到剪贴板')
  }
}

const formatDateTime = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

onMounted(() => {
  loadBlogger()
  loadCategories()
  loadTeams()
  loadTeamMembers()
  loadStatusList()
  
  countdownTimer = setInterval(() => {
    calculateCountdown()
  }, 60000)
})

onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
})

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

const loadStatusList = async () => {
  try {
    const res = await getBloggerStatusList()
    if (res.code === 200) {
      bloggerStatusList.value = res.data
    }
  } catch (error) {
    console.error('加载状态列表失败', error)
  }
}
</script>

<style scoped>
.detail-page {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.countdown-card {
  background: linear-gradient(135deg, #fff7ed 0%, #ffedd5 100%) !important;
  border: 2px solid #fdba74 !important;
  margin-bottom: 20px;
}

.countdown-card h3 {
  color: #c2410c !important;
  margin-bottom: 16px;
}

.countdown-display {
  text-align: center;
  margin-bottom: 16px;
}

.countdown-value {
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 4px;
  margin-bottom: 8px;
}

.countdown-number {
  font-size: 48px;
  font-weight: 800;
  line-height: 1;
}

.countdown-unit {
  font-size: 20px;
  font-weight: 600;
}

.countdown-text {
  font-size: 14px;
  font-weight: 500;
}

.countdown-value.normal .countdown-number,
.countdown-value.normal .countdown-unit,
.countdown-text.normal {
  color: #ea580c;
}

.countdown-value.warning .countdown-number,
.countdown-value.warning .countdown-unit,
.countdown-text.warning {
  color: #dc2626;
}

.countdown-value.urgent .countdown-number,
.countdown-value.urgent .countdown-unit,
.countdown-text.urgent {
  color: #991b1b;
  animation: pulse 1s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}

.countdown-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.7);
  border-radius: 8px;
  color: #c2410c;
  font-size: 13px;
}

.countdown-tip svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.page-header {
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left,
.header-right {
  display: flex;
  gap: 12px;
  align-items: center;
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

.back-btn svg {
  width: 18px;
  height: 18px;
}

.edit-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: var(--primary);
  border: none;
  border-radius: 10px;
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.edit-btn:hover {
  background: var(--primary-dark);
}

.edit-btn svg {
  width: 18px;
  height: 18px;
}

.transfer-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: #3b82f6;
  border: none;
  border-radius: 10px;
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.transfer-btn:hover {
  background: #2563eb;
}

.transfer-btn svg {
  width: 18px;
  height: 18px;
}

.delete-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: #ef4444;
  border: none;
  border-radius: 10px;
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.delete-btn:hover {
  background: #dc2626;
}

.delete-btn svg {
  width: 18px;
  height: 18px;
}

.detail-content {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 24px;
}

.main-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.profile-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 20px;
  padding: 32px;
  text-align: center;
}

.avatar {
  width: 100px;
  height: 100px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 36px;
  color: white;
  font-weight: 600;
  margin: 0 auto 16px;
  overflow: hidden;
  position: relative;
  cursor: default;
}

.avatar.has-image {
  cursor: pointer;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.avatar:hover .avatar-overlay {
  opacity: 1;
}

.avatar-overlay svg {
  width: 24px;
  height: 24px;
  color: white;
  margin-bottom: 4px;
}

.avatar-overlay span {
  font-size: 12px;
  color: white;
}

.profile-card h1 {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.profile-link {
  cursor: pointer;
  transition: all 0.3s ease;
}

.profile-link:hover {
  transform: translateY(-2px);
}

.profile-link:hover .avatar {
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.profile-link:hover h1 {
  color: var(--primary);
}

.owner-link {
  color: var(--primary);
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
}

.owner-link:hover {
  color: var(--primary-dark);
  text-decoration: underline;
}

.edit-input, .edit-select, .edit-textarea {
  padding: 8px 12px;
  background: var(--bg-dark);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 14px;
  width: 100%;
  box-sizing: border-box;
}

.edit-input:focus, .edit-select:focus, .edit-textarea:focus {
  outline: none;
  border-color: var(--primary);
}

.profile-card h1 {
  margin-bottom: 8px;
}

.tags-display {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 8px;
  margin-bottom: 16px;
}

.tag {
  padding: 6px 14px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border-radius: 20px;
  font-size: 13px;
  color: white;
  font-weight: 500;
}

.tags-edit {
  display: flex;
  gap: 8px;
  justify-content: center;
  margin-bottom: 12px;
}

.add-tag-btn {
  padding: 8px 16px;
  background: var(--primary);
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 13px;
  cursor: pointer;
}

.meta-info {
  display: flex;
  justify-content: center;
  gap: 24px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: var(--text-secondary);
}

.meta-item svg {
  width: 16px;
  height: 16px;
}

.category-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  border-radius: 8px;
  font-weight: 700;
  font-size: 14px;
  letter-spacing: 0.5px;
  transition: transform 0.3s ease;
  border: 2px solid rgba(255, 255, 255, 0.7);
}

.category-badge:hover {
  transform: scale(1.05);
}

.category-icon {
  font-size: 20px;
  line-height: 1;
}

.category-text {
  font-size: 14px;
  line-height: 1;
}

.edit-select {
  padding: 4px 8px;
}

.info-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  padding: 24px;
}

.info-card h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-color);
}

.search-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.search-select-group,
.search-prompt-group,
.search-preview {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.search-label {
  font-size: 13px;
  color: var(--text-secondary);
  font-weight: 500;
}

.search-select {
  width: 100%;
}

.prompt-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.prompt-btn {
  padding: 6px 12px;
  border: 1px solid var(--border-color);
  background: var(--bg-card);
  border-radius: 8px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.prompt-btn:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.prompt-btn.active {
  background: var(--primary);
  border-color: var(--primary);
  color: white;
}

.search-preview-content {
  padding: 10px 12px;
  background: var(--bg-page);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 13px;
  color: var(--text-primary);
  word-break: break-all;
}

.search-preview-input {
  width: 100%;
  padding: 10px 12px;
  background: var(--bg-page);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 13px;
  color: var(--text-primary);
  transition: all 0.2s;
}

.search-preview-input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.search-btn {
  width: 100%;
  padding: 12px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border: none;
  border-radius: 10px;
  color: white;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.search-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.card-header-with-action {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-color);
}

.card-header-with-action h3 {
  margin: 0;
  padding: 0;
  border: none;
}

.refresh-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  transition: all 0.2s;
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.refresh-btn:hover:not(:disabled) {
  color: var(--primary);
}

.refresh-btn svg {
  width: 18px;
  height: 18px;
}

.refresh-btn svg.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.contact-with-copy {
  display: flex;
  align-items: center;
  gap: 8px;
}

.copy-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  transition: all 0.2s;
  border-radius: 4px;
}

.copy-btn:hover {
  background: var(--primary-light);
  color: var(--primary);
}

.copy-btn svg {
  width: 16px;
  height: 16px;
}

/* 对接人信息卡片样式 */
.contact-person-card {
  margin-top: 20px;
}

.contact-person-link {
  cursor: pointer;
  transition: all 0.3s ease;
}

.contact-person-link:hover {
  transform: translateX(4px);
}

.contact-person-link:hover .contact-person-avatar {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.contact-person-link:hover .jump-profile-btn {
  background: linear-gradient(135deg, var(--primary-dark), var(--primary));
  transform: translateX(4px);
}

.jump-profile-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(249, 115, 22, 0.2);
}

.jump-profile-btn svg {
  width: 20px;
  height: 20px;
  color: #fff;
}

.jump-profile-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 16px rgba(249, 115, 22, 0.4);
}

.contact-person-content {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: var(--bg-dark);
  border-radius: 12px;
}

.contact-person-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
  transition: all 0.3s ease;
  overflow: hidden;
}

.contact-person-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.contact-person-avatar svg {
  width: 32px;
  height: 32px;
  transition: all 0.3s ease;
}

.contact-person-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.contact-person-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.contact-person-username {
  font-size: 13px;
  color: var(--text-secondary);
}

.contact-person-team {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-item label {
  font-size: 12px;
  color: var(--text-muted);
}

.info-item span {
  font-size: 14px;
  color: var(--text-primary);
}

.info-item.full-width {
  grid-column: 1 / -1;
  margin-top: 8px;
}

.description {
  line-height: 1.6;
  color: var(--text-secondary) !important;
}

.edit-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 16px;
}

.edit-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.edit-item label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

.edit-actions {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.save-btn {
  padding: 8px 20px;
  background: var(--primary);
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 13px;
  cursor: pointer;
}

.cancel-btn {
  padding: 8px 20px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-secondary);
  font-size: 13px;
  cursor: pointer;
}

.tags-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.cloud-tag {
  padding: 8px 16px;
  background: var(--bg-card-hover);
  border: 1px solid var(--border-color);
  border-radius: 20px;
  font-size: 14px;
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.2s;
}

.cloud-tag:hover {
  background: var(--primary);
  border-color: var(--primary);
  color: white;
}

.followup-form {
  margin-bottom: 20px;
}

.followup-form textarea {
  width: 100%;
  padding: 12px;
  background: var(--bg-dark);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  color: var(--text-primary);
  font-size: 14px;
  resize: vertical;
  font-family: inherit;
  box-sizing: border-box;
}

.followup-form textarea:focus {
  outline: none;
  border-color: var(--primary);
}

.followup-actions {
  display: flex;
  gap: 12px;
  margin-top: 12px;
  align-items: center;
}

.type-select {
  padding: 8px 12px;
  background: var(--bg-dark);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 13px;
}

.followup-form .add-btn {
  padding: 8px 20px;
  background: var(--primary);
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.followup-form .add-btn:hover:not(:disabled) {
  background: var(--primary-dark);
}

.followup-form .add-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.followup-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.followup-item {
  padding: 14px;
  background: var(--bg-card-hover);
  border-radius: 10px;
  border: 1px solid var(--border-color);
}

.followup-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.followup-type {
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

.followup-type.跟进 {
  background: rgba(76, 175, 80, 0.2);
  color: #4CAF50;
}

.followup-type.电话 {
  background: rgba(33, 150, 243, 0.2);
  color: #2196F3;
}

.followup-type.微信 {
  background: rgba(76, 175, 80, 0.2);
  color: #4CAF50;
}

.followup-type.邮件 {
  background: rgba(156, 39, 176, 0.2);
  color: #9C27B0;
}

.followup-type.会面 {
  background: rgba(255, 152, 0, 0.2);
  color: #FF9800;
}

.followup-time {
  font-size: 12px;
  color: var(--text-muted);
}

.followup-edit {
  margin-bottom: 8px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.followup-edit .type-select {
  width: fit-content;
  padding: 8px 12px;
  background: var(--bg-dark);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 13px;
}

.followup-edit .edit-textarea {
  width: 100%;
  padding: 12px;
  background: var(--bg-dark);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  color: var(--text-primary);
  font-size: 14px;
  resize: vertical;
  font-family: inherit;
  box-sizing: border-box;
}

.followup-edit .edit-textarea:focus {
  outline: none;
  border-color: var(--primary);
}

.followup-edit .edit-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.followup-edit .save-btn,
.followup-edit .cancel-btn {
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.followup-edit .save-btn {
  background: var(--primary);
  color: white;
}

.followup-edit .save-btn:hover {
  background: var(--primary-dark);
}

.followup-edit .cancel-btn {
  background: var(--bg-dark);
  color: var(--text-secondary);
}

.followup-edit .cancel-btn:hover {
  background: var(--border-color);
}

.followup-content-text {
  font-size: 14px;
  color: var(--text-primary);
  line-height: 1.6;
  margin-bottom: 8px;
}

.followup-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.followup-operator {
  font-size: 12px;
  color: var(--text-muted);
}

.followup-ops {
  display: flex;
  gap: 8px;
}

.op-btn {
  padding: 4px 10px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-secondary);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.op-btn:hover {
  background: var(--primary);
  border-color: var(--primary);
  color: white;
}

.delete-btn {
  padding: 4px 10px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-muted);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.delete-btn:hover {
  background: #f44336;
  border-color: #f44336;
  color: white;
}

.log-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-height: 300px;
  overflow-y: auto;
}

.log-item {
  display: grid;
  grid-template-columns: 60px 1fr 60px 140px;
  gap: 10px;
  padding: 10px;
  background: var(--bg-card-hover);
  border-radius: 8px;
  font-size: 12px;
  align-items: center;
}

.log-action {
  font-weight: 500;
  color: var(--primary);
}

.log-detail {
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

html.dark .dialog-overlay {
  background: rgba(0, 0, 0, 0.8);
}

.dialog {
  background: #ffffff;
  border-radius: 16px;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: slideUp 0.3s ease;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 2px solid #e5e7eb;
  background: #f9fafb;
}

.dialog-header h3 {
  font-size: 18px;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
}

.dialog-header .close-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
}

.dialog-header .close-btn:hover {
  background: #ef4444;
  border-color: #ef4444;
  color: white;
}

.dialog-header .close-btn svg {
  width: 20px;
  height: 20px;
}

.dialog-body {
  padding: 24px;
  background: #ffffff;
}

.dialog-body .form-group {
  margin-bottom: 20px;
}

.dialog-body .form-group:last-child {
  margin-bottom: 0;
}

.dialog-body .form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

.dialog-body .form-select,
.dialog-body .form-textarea {
  width: 100%;
  padding: 12px;
  background: #f9fafb;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  color: #1f2937;
  font-size: 14px;
  font-family: inherit;
  box-sizing: border-box;
}

.dialog-body .form-select:focus,
.dialog-body .form-textarea:focus {
  outline: none;
  border-color: #ff6b35;
  box-shadow: 0 0 0 3px rgba(255, 107, 53, 0.1);
  background: #ffffff;
}

.dialog-body .form-textarea {
  resize: vertical;
}

.dialog-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 2px solid #e5e7eb;
}

.dialog-actions .btn-secondary {
  padding: 10px 20px;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  color: #374151;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.dialog-actions .btn-secondary:hover {
  background: #e5e7eb;
  border-color: #d1d5db;
}

.dialog-actions .btn-primary {
  padding: 10px 20px;
  background: linear-gradient(135deg, #ff6b35, #f7931e);
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.dialog-actions .btn-primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #f7931e, #ff8c5a);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(255, 107, 53, 0.3);
}

.dialog-actions .btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.timeline {
  position: relative;
  padding-left: 30px;
}

.timeline-item {
  position: relative;
  padding-bottom: 24px;
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-dot {
  position: absolute;
  left: -30px;
  top: 0;
  width: 12px;
  height: 12px;
  background: var(--primary);
  border-radius: 50%;
  border: 3px solid var(--bg-card);
  box-shadow: 0 0 0 2px var(--primary);
}

.timeline-item:not(:last-child)::before {
  content: '';
  position: absolute;
  left: -25px;
  top: 12px;
  bottom: 0;
  width: 2px;
  background: var(--border-color);
}

.timeline-content {
  background: var(--bg-card-hover);
  border-radius: 12px;
  padding: 16px;
  transition: all 0.3s ease;
}

.timeline-content:hover {
  background: var(--bg-card);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-color);
}

.timeline-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.timeline-date {
  font-size: 13px;
  color: var(--text-muted);
}

.timeline-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.timeline-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.timeline-info .info-label {
  color: var(--text-muted);
}

.timeline-info .info-value {
  color: var(--text-primary);
  font-weight: 500;
}

.timeline-info .info-value.status-completed {
  color: var(--success);
}

.timeline-info .info-value.status-progress {
  color: var(--primary);
}

.timeline-info .info-value.status-cancelled {
  color: var(--danger);
}

.timeline-info .info-value.status-pending {
  color: var(--warning);
}

.timeline-note {
  padding: 12px;
  background: var(--bg-dark);
  border-radius: 8px;
  font-size: 14px;
}

.timeline-note .note-label {
  color: var(--text-muted);
  font-weight: 500;
}

.timeline-note .note-text {
  color: var(--text-secondary);
  line-height: 1.6;
}

.timeline-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.timeline-actions .op-btn,
.timeline-actions .delete-btn {
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.timeline-actions .op-btn {
  background: var(--bg-dark);
  color: var(--text-primary);
}

.timeline-actions .op-btn:hover {
  background: var(--primary-light);
  color: var(--primary);
}

.timeline-actions .delete-btn {
  background: #fee2e2;
  color: #dc2626;
}

.timeline-actions .delete-btn:hover {
  background: #fecaca;
}

.card-header-with-action .add-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
  z-index: 10;
}

.card-header-with-action .add-btn:hover {
  background: var(--primary-dark);
}

.card-header-with-action .add-btn svg {
  width: 16px;
  height: 16px;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.cooperation-dialog {
  max-width: 550px;
}

@media (max-width: 768px) {
  .page-header {
    flex-wrap: wrap;
    gap: 8px;
  }

  .page-header button {
    padding: 8px 12px;
    font-size: 12px;
  }

  .page-header button span {
    display: none;
  }

  .page-header button svg {
    width: 18px;
    height: 18px;
  }

  .detail-content {
    grid-template-columns: 1fr;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .edit-grid {
    grid-template-columns: 1fr;
  }

  .log-item {
    grid-template-columns: 1fr;
    gap: 8px;
  }

  .log-item > span {
    margin-bottom: 4px;
  }

  .followup-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .followup-actions select {
    width: 100%;
  }

  .timeline {
    padding-left: 20px;
  }

  .timeline-dot {
    left: -20px;
  }

  .timeline-item:not(:last-child)::before {
    left: -15px;
  }

  .profile-card {
    padding: 20px;
  }

  .info-card {
    padding: 16px;
  }

  .avatar {
    width: 80px;
    height: 80px;
    font-size: 28px;
  }

  .profile-card h1 {
    font-size: 20px;
  }

  .tags-edit {
    flex-direction: column;
  }

  .tags-edit .add-tag-btn {
    width: 100%;
  }

  .meta-info {
    flex-direction: column;
    align-items: center;
    gap: 12px;
  }

  /* 移动端对接人卡片优化 */
  .contact-person-content {
    flex-direction: column;
    text-align: center;
  }

  .contact-person-avatar {
    width: 80px;
    height: 80px;
  }

  .contact-person-info {
    align-items: center;
  }

  .contact-person-name {
    justify-content: center;
    flex-wrap: wrap;
  }

  .contact-person-meta {
    justify-content: center;
  }

  .jump-profile-btn {
    align-self: center;
    margin-top: 8px;
  }

  .view-profile-btn {
    align-self: center;
    width: 100%;
    justify-content: center;
  }

  /* 移动端复制按钮优化 */
  .contact-with-copy {
    justify-content: space-between;
  }

  .copy-btn {
    padding: 8px;
  }

  /* 移动端博主信息卡片优化 */
  .info-card.contact-person-card {
    padding: 16px;
  }

  .contact-person-link h3,
  .contact-person-card h3 {
    font-size: 14px;
    margin-bottom: 12px;
  }

  .dialog {
    width: 95%;
    max-width: 400px;
  }

  .dialog-header {
    padding: 16px 20px;
  }

  .dialog-body {
    padding: 20px;
  }
}

@media (max-width: 480px) {
  .page-container {
    padding: 8px;
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }

  .page-header button {
    justify-content: center;
  }

  .profile-card {
    padding: 16px;
  }

  .info-card {
    padding: 12px;
  }

  .info-card h3 {
    font-size: 14px;
    margin-bottom: 16px;
  }

  .followup-form textarea {
    padding: 10px;
    font-size: 13px;
  }

  .followup-item {
    padding: 12px;
  }

  .timeline-content {
    padding: 12px;
  }

  .timeline-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
    padding-bottom: 8px;
  }

  .timeline-title {
    font-size: 14px;
  }

  .timeline-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
}

/* 编辑弹窗样式 */
.dialog-overlay {
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
  animation: fadeIn 0.2s ease;
}

.edit-dialog {
  background: #fff;
  border-radius: 16px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
}

.dialog-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #1a1a1a;
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  transition: all 0.2s;
  border-radius: 4px;
}

.close-btn:hover {
  background: #f3f4f6;
  color: #1a1a1a;
}

.close-btn svg {
  width: 20px;
  height: 20px;
}

.dialog-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  color: #1a1a1a;
  transition: all 0.2s;
  box-sizing: border-box;
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.avatar-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.avatar-preview {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: linear-gradient(135deg, #f3f4f6, #e5e7eb);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: 600;
  color: #666;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.avatar-preview.has-image {
  background: none;
}

.avatar-preview.has-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-preview:hover {
  transform: scale(1.05);
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.avatar-preview:hover .avatar-overlay {
  opacity: 1;
}

.avatar-overlay svg {
  width: 24px;
  height: 24px;
  margin-bottom: 4px;
}

.avatar-overlay span {
  font-size: 12px;
}

.upload-tip {
  font-size: 12px;
  color: #666;
  margin: 0;
}

.tags-edit-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tag-input-wrapper {
  display: flex;
  gap: 8px;
}

.tag-input-wrapper .form-input {
  flex: 1;
}

.add-tag-btn {
  padding: 10px 20px;
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.add-tag-btn:hover {
  background: var(--primary-dark);
}

.tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  color: #fff;
  border-radius: 16px;
  font-size: 13px;
  font-weight: 500;
}

.remove-tag {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  padding: 0;
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  line-height: 1;
  opacity: 0.8;
  transition: opacity 0.2s;
}

.remove-tag:hover {
  opacity: 1;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #e5e7eb;
  background: #f9fafb;
}

.btn-secondary {
  padding: 10px 24px;
  background: #f3f4f6;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: #e5e7eb;
  border-color: #9ca3af;
}

.btn-primary {
  padding: 10px 24px;
  background: linear-gradient(135deg, #ff6b35, #f7931e);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 107, 53, 0.3);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 移动端弹窗优化 */
@media (max-width: 768px) {
  .edit-dialog {
    width: 95%;
    max-height: 95vh;
  }

  .dialog-header {
    padding: 16px 20px;
  }

  .dialog-body {
    padding: 20px;
  }

  .dialog-actions {
    padding: 16px 20px;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .avatar-preview {
    width: 80px;
    height: 80px;
    font-size: 28px;
  }

  .tag-input-wrapper {
    flex-direction: column;
  }

  .add-tag-btn {
    width: 100%;
  }
}

.status-tag {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 500;
}

.status-init {
  background: #e3f2fd;
  color: #1976d2;
}

.status-talking {
  background: #fff3e0;
  color: #f57c00;
}

.status-cooperated {
  background: #e8f5e9;
  color: #388e3c;
}

.status-rejected {
  background: #ffebee;
  color: #d32f2f;
}

.status-paused {
  background: #f3e5f5;
  color: #7b1fa2;
}

.crop-dialog {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.crop-dialog-content {
  background: var(--bg-card);
  border-radius: 16px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.crop-dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.crop-dialog-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.crop-dialog-header .close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: var(--bg-page);
  border-radius: 8px;
  font-size: 24px;
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.crop-dialog-header .close-btn:hover {
  background: var(--border-color);
  color: var(--text-primary);
}

.crop-dialog-body {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.crop-container {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-page);
  border-radius: 12px;
  padding: 20px;
}

.crop-canvas {
  border-radius: 8px;
  cursor: grab;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.crop-canvas:active {
  cursor: grabbing;
}

.crop-controls {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.control-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.control-group label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

.control-group input[type="range"] {
  width: 100%;
  height: 6px;
  border-radius: 3px;
  background: var(--border-color);
  outline: none;
  -webkit-appearance: none;
}

.control-group input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: var(--primary);
  cursor: pointer;
  transition: all 0.2s;
}

.control-group input[type="range"]::-webkit-slider-thumb:hover {
  transform: scale(1.1);
  box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.2);
}

.crop-dialog-actions {
  display: flex;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid var(--border-color);
}

.crop-dialog-actions .btn-secondary,
.crop-dialog-actions .btn-primary {
  flex: 1;
}

/* 暗色模式样式 */
.dark .status-init {
  background: rgba(25, 118, 210, 0.2);
  color: #64b5f6;
}

.dark .status-talking {
  background: rgba(245, 124, 0, 0.2);
  color: #ffb74d;
}

.dark .status-cooperated {
  background: rgba(56, 142, 60, 0.2);
  color: #81c784;
}

.dark .status-rejected {
  background: rgba(211, 47, 47, 0.2);
  color: #e57373;
}

.dark .status-paused {
  background: rgba(123, 31, 162, 0.2);
  color: #ba68c8;
}

.dark .search-preview-input {
  background: var(--bg-dark);
  color: var(--text-primary);
}

.dark .search-select {
  background: var(--bg-dark);
  color: var(--text-primary);
}

.dark .prompt-btn {
  background: var(--bg-dark);
  color: var(--text-primary);
  border-color: var(--border-color);
}

.dark .prompt-btn.active {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}
</style>