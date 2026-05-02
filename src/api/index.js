import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 15000
})

const MAX_RETRIES = 3
const RETRY_DELAY = 1000

const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms))

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  response => {
    if (response.config.responseType === 'blob') {
      return response.data
    }
    const data = response.data
    if (data && data.code === 401) {
      localStorage.clear()
      window.location.href = '/login'
      return Promise.reject(new Error(data.message || 'Unauthorized'))
    }
    return data
  },
  async error => {
    const config = error.config

    if (error.response?.status === 401) {
      localStorage.clear()
      window.location.href = '/login'
      return Promise.reject(error)
    }

    if (!config || !config._retry) {
      config._retry = 0
    }

    if (config._retry < MAX_RETRIES && !error.response?.status?.toString().startsWith('4')) {
      config._retry += 1
      console.warn(`API 请求重试 (${config._retry}/${MAX_RETRIES}): ${config.url}`)
      await delay(RETRY_DELAY * config._retry)
      return api(config)
    }

    console.error('API 请求失败:', error)
    return Promise.reject(error)
  }
)

export const loginAPI = (data) => api.post('/login', data)

export const registerAPI = (data) => api.post('/register', data)

export const forgotPasswordAPI = () => Promise.reject(new Error('功能暂未实现'))

export const bloggerListAPI = (params) => api.get('/blogger/list', { params })

export const bloggerSuggestionsAPI = (keyword) => api.get('/blogger/suggestions', { params: { keyword } })

export const bloggerMyAPI = (params) => api.get('/blogger/my', { params })

export const getMyBloggersAPI = (params) => api.get('/blogger/my', { params })

export const bloggerAddAPI = (data) => api.post('/blogger/add', data)

export const bloggerDeleteAPI = (id) => api.post('/admin/blogger/delete', { id })

export const bloggerStatAPI = () => api.get('/blogger/stat')

export const getBloggerCharts = () => api.get('/blogger/charts')

export const getBloggerChartsAPI = () => api.get('/blogger/charts')

export const followupAddAPI = (data) => api.post('/followup/add', data)

export const followupListAPI = (params) => api.get('/followup/list', { params })

export const followupDeleteAPI = (id) => api.delete(`/followup/${id}`)

export const bloggerUpdateAPI = (id, data) => api.put(`/blogger/${id}`, data)

export const followupUpdateAPI = (id, data) => api.put(`/followup/${id}`, data)

export const cooperationListAPI = (params) => api.get('/cooperation/list', { params })
export const cooperationAddAPI = (data) => api.post('/cooperation/add', data)
export const cooperationUpdateAPI = (id, data) => api.put(`/cooperation/${id}`, data)
export const cooperationDeleteAPI = (id) => api.delete(`/cooperation/${id}`)

export const getCategoryList = () => api.get('/category/list')

export const categoryListAPI = () => api.get('/category/list')

export const categoryAddAPI = (data) => api.post('/categories', data)

export const categoryUpdateAPI = (data) => api.put(`/categories/${data.id}`, data)

export const categoryDeleteAPI = (id) => api.delete(`/categories/${id}`)

export const platformListAPI = () => api.get('/platform/list')

export const productListAPI = () => api.get('/product/list')

export const productAddAPI = (data) => api.post('/products', data)

export const productUpdateAPI = (id, data) => api.put(`/products/${id}`, data)

export const productDeleteAPI = (id) => api.delete(`/products/${id}`)

export const userListAPI = () => api.get('/users/public')

export const deactivateUserAPI = (id) => api.put(`/users/${id}`, { status: 'inactive' })

export const deleteUserAPI = (id) => api.delete(`/users/${id}`)

export const resetPasswordAPI = (id) => api.post('/admin/users/reset-password', { id })

export const exportAPI = (type, fields) => api.get('/admin/export', { params: { type, fields }, responseType: 'blob' })

export const getOperationLog = (params) => api.get('/admin/log', { params })

export const editBlogger = (data) => api.post('/admin/blogger/edit', data)

export const getBloggerStatusList = () => api.get('/blogger/status/list')

export const getBloggerDetail = (id) => api.get(`/blogger/${id}`)

export const getBloggerHotData = (nickname) => api.get(`/blogger/hot/${nickname}`)

export const getBloggerNews = (nickname) => api.get(`/blogger/news/${nickname}`)

export const getBloggerStatusHistory = (id) => api.get(`/blogger/${id}/history`)

export const requestBloggerTransfer = (data) => api.post(`/blogger/${data.id}/transfer`, data)

export const confirmBloggerTransfer = (data) => api.post(`/blogger/transfers/${data.id}/confirm`, data)

export const getBloggerTransferRequests = () => api.get('/blogger/transfers')

export const deleteBlogger = (id) => api.post('/admin/blogger/delete', { id })

export const exportData = (type) => api.get('/admin/export', { params: { type }, responseType: 'blob' })

// 用户设置相关API
export const getUserProfileAPI = () => api.get('/user/profile')

export const updateUserProfileAPI = (data) => api.put('/user/profile', data)

export const updatePasswordAPI = (data) => api.post('/user/change-password', data)

export const uploadAvatarAPI = (data) => api.post('/user/avatar', data)

export const getPublicUsersAPI = () => api.get('/users/public')

export const getPendingUsersAPI = () => api.get('/users/pending')

export const batchSetAdminAPI = async (teamId, userIds) => {
  return Promise.all(userIds.map(userId => setTeamAdminAPI(teamId, userId)))
}

export const approveUserAPI = (id) => api.post('/users/approve', { id })

export const rejectUserAPI = (id) => api.put(`/users/${id}`, { status: 'rejected' })

export const setAdminAPI = (id) => api.put(`/users/${id}`, { role: 'admin' })
export const removeAdminAPI = (id) => api.put(`/users/${id}`, { role: 'user' })

export const getTeamsAPI = () => api.get('/teams')
export const createTeamAPI = (data) => api.post('/teams', data)
export const updateTeamAPI = (id, data) => api.put(`/teams/${id}`, data)
export const deleteTeamAPI = (id, force = false) => api.delete(`/teams/${id}${force ? '?force=true' : ''}`)
export const setUserTeamAPI = (id, team_id) => api.put(`/users/${id}`, { team_id })
export const updateMyTeamAPI = (team_id) => api.put(`/my/team`, { team_id })
export const generateInviteCodeAPI = () => api.post('/user/generate-invite-code')
export const removeTeamMemberAPI = (teamId, userId) => api.delete(`/team/${teamId}/members/${userId}`)

export const getTeamBloggerStatAPI = (team_id) => api.get(`/team/blogger/stat?team_id=${team_id}`)
export const getTeamBloggerChartsAPI = (team_id) => api.get(`/team/blogger/charts?team_id=${team_id}`)

export const uploadTeamLogoAPI = (data) => api.post('/upload/team-logo', data)
export const uploadTeamBgAPI = (data) => api.post('/upload/team-bg', data)

export const getTeamPostsAPI = (teamId, page = 1, pageSize = 20, category = '') => {
  const params = { page, pageSize }
  if (category) params.category = category
  return api.get(`/team/${teamId}/posts`, { params })
}
export const createTeamPostAPI = (teamId, data) => api.post(`/team/${teamId}/posts`, data)
export const getTeamPostDetailAPI = (teamId, postId) => api.get(`/team/posts/${postId}`)
export const deleteTeamPostAPI = (teamId, postId) => api.delete(`/team/posts/${postId}`)
export const createTeamPostCommentAPI = (teamId, postId, data) => api.post(`/team/posts/${postId}/comments`, data)
export const likeTeamPostAPI = (teamId, postId) => api.post(`/team/posts/${postId}/like`)
export const getTeamPostLikeStatusAPI = (postId) => api.get(`/team/posts/${postId}/like-status`)
export const pinTeamPostAPI = (teamId, postId) => api.post(`/team/posts/${postId}/pin`)
export const featureTeamPostAPI = (teamId, postId) => api.post(`/team/posts/${postId}/feature`)
export const collectTeamPostAPI = (teamId, postId) => api.post(`/team/posts/${postId}/collect`)
export const getTeamPostCollectStatusAPI = (postId) => api.get(`/team/posts/collect-status/${postId}`)
export const getTeamPostsCollectedAPI = () => api.get('/team/posts/collected')
export const searchTeamPostsAPI = (keyword) => api.get('/team/posts/search', { params: { keyword } })
export const setTeamAdminAPI = (teamId, userId) => api.put(`/team/${teamId}/admins/${userId}`)
export const getTeamMessagesAPI = (teamId) => api.get(`/team/${teamId}/messages`)
export const sendTeamMessageAPI = (teamId, data) => api.post(`/team/${teamId}/messages`, data)
export const markMessageReadAPI = (teamId, messageId) => api.put(`/team/${teamId}/messages/${messageId}/read`)
export const deleteTeamMessageAPI = (teamId, messageId) => api.delete(`/team/${teamId}/messages/${messageId}`)
export const getTeamOperationLogsAPI = (teamId, page = 1, pageSize = 50) => api.get(`/team/${teamId}/logs`, { params: { page, pageSize } })

export const uploadBloggerAvatarAPI = (data) => api.post('/upload/blogger-avatar', data)

export const getUserDetailAPI = (username) => api.get(`/user/detail/${username}`)

export const getUserBloggersAPI = (username) => api.get(`/user/${username}/bloggers`)

export const getUserFollowupStatsAPI = (username) => api.get(`/user/${username}/followup-stats`)

export const getNotificationsAPI = () => api.get('/notifications')

export const markNotificationReadAPI = (id) => api.post(`/notifications/${id}/read`)

export const markAllNotificationsReadAPI = () => api.post('/notifications/mark-all-read')

export const deleteNotificationAPI = (id) => api.delete(`/notifications/${id}`)

export const batchDeleteNotificationsAPI = (ids) => api.post('/notifications/batch-delete', { ids })

export const sendEmailNotificationAPI = (data) => api.post('/notifications/email-send', data)

export const updateNotificationPrefsAPI = (prefs) => api.put('/notifications/preferences', prefs)

export const getNotificationPrefsAPI = () => api.get('/notifications/preferences')

export const createDataBackupAPI = (label) => api.post('/system/backup', { label })

export const listBackupsAPI = () => api.get('/system/backups')

export const restoreBackupAPI = (backupId) => api.post(`/system/backup/${backupId}/restore`)

export const deleteBackupAPI = (backupId) => api.delete(`/system/backup/${backupId}`)

export const exportAllDataAPI = () => api.get('/system/export')

export const setInvalidBloggerAPI = (id, data) => api.post(`/blogger/${id}/invalid`, data)

export const invalidBloggerListAPI = (params) => api.get('/blogger/invalid/list', { params })

export const bindInvalidBloggerAPI = () => Promise.reject(new Error('功能暂未实现'))

export const getInvalidBloggerCountAPI = () => api.get('/blogger/invalid/count')

export const submitBloggerEvaluationAPI = (data) => api.post('/blogger/evaluation', data)

export const getBloggerEvaluationAPI = (blogger_id) => api.get(`/blogger/evaluation/${blogger_id}`)

export const getExpiringBloggersWithoutContactAPI = () => api.get('/blogger/expiring-without-contact')

export const getInviteInfoAPI = () => api.get('/user/invite-info')
export const getBloggerAddStatsAPI = () => api.get('/user/blogger-add-stats')
export const getInviteCodeRecordsAPI = () => api.get('/user/invite-code-records')

export const getAnnouncementsAPI = () => api.get('/announcements')

export const createAnnouncementAPI = (data) => api.post('/announcements', data)

export const getConversationsAPI = () => api.get('/messages/conversations')
export const getMessagesAPI = (userId, page = 1, pageSize = 50) => api.get('/messages', { params: { user_id: userId, page, pageSize } })
export const sendMessageAPI = (data) => api.post('/messages', data)
export const markMessagesReadAPI = (userId) => api.post('/messages/read', null, { params: { user_id: userId } })
export const getUnreadCountAPI = () => api.get('/messages/unread')

export const getUsersListAPI = () => api.get('/users/list')
export const deleteLogAPI = (id) => api.delete(`/admin/log/${id}`)
export const clearLogsAPI = () => api.delete('/admin/log/clear')
export const clearOldLogsAPI = (count) => api.delete(`/admin/log/clear/${count}`)
export const batchApproveUsersAPI = (user_ids) => api.post('/admin/users/batch-approve', { user_ids })
export const batchRejectUsersAPI = (user_ids) => api.post('/admin/users/batch-reject', { user_ids })
export const getRegistrationModeAPI = () => api.get('/admin/registration-mode')
export const setRegistrationModeAPI = (mode) => api.post('/admin/registration-mode', { mode })
export const deactivateUserAdminAPI = (id) => api.post('/admin/users/deactivate', { id })
export const deleteUserAdminAPI = (id) => api.post('/admin/users/delete', { id })
export const setUserRoleAPI = (id, role) => api.post('/admin/users/role', { id, role })
export const getBloggerLogsAPI = (nickname) => api.get('/blogger/logs', { params: { nickname } })
export const sendPrivateMessageAPI = (data) => api.post('/messages/private', data)
export const createBackupAPI = (data) => api.post('/backup', data)
export const getSnapshotsAPI = () => api.get('/snapshots')
export const getSnapshotSettingsAPI = () => api.get('/snapshots/settings')
export const setSnapshotSettingsAPI = (data) => api.post('/snapshots/settings', data)
export const createSnapshotAPI = (data) => api.post('/snapshots', data)
export const restoreSnapshotAPI = (filename) => api.post(`/snapshots/file/${filename}/restore`)
export const downloadSnapshotAPI = (filename) => api.get(`/snapshots/file/${filename}`, { responseType: 'blob' })
export const deleteSnapshotAPI = (filename) => api.delete(`/snapshots/${filename}`)
export const clearDataAPI = () => api.post('/clear')
export const purgeDataAPI = (data) => api.post('/admin/purge', data)

export const batchUpdateStatusAPI = (ids, status) => api.post('/blogger/batch/status', { ids, status })
export const batchUpdateTagsAPI = (ids, tag_ids, action) => api.post('/blogger/batch/tags', { ids, tag_ids, action })
export const batchDeleteAPI = (ids) => api.post('/blogger/batch/delete', { ids })

export const getSavedFiltersAPI = () => api.get('/saved-filters')
export const createSavedFilterAPI = (data) => api.post('/saved-filters', data)
export const updateSavedFilterAPI = (id, data) => api.put(`/saved-filters/${id}`, data)
export const deleteSavedFilterAPI = (id) => api.delete(`/saved-filters/${id}`)
export const setDefaultFilterAPI = (id) => api.post(`/saved-filters/${id}/default`)

export const getRecommendTagsAPI = (params) => api.get('/blogger/recommend-tags', { params })
export const getSimilarBloggersAPI = (params) => api.get('/blogger/similar', { params })

export const getWorkflowRulesAPI = () => api.get('/workflow/rules')
export const createWorkflowRuleAPI = (data) => api.post('/workflow/rules', data)
export const updateWorkflowRuleAPI = (id, data) => api.put(`/workflow/rules/${id}`, data)
export const deleteWorkflowRuleAPI = (id) => api.delete(`/workflow/rules/${id}`)
export const toggleWorkflowRuleAPI = (id) => api.post(`/workflow/rules/${id}/toggle`)
export const triggerWorkflowAPI = (data) => api.post('/workflow/trigger', data)

export const getTagsAPI = () => api.get('/tags')
export const createTagAPI = (data) => api.post('/tags', data)
export const updateTagAPI = (id, data) => api.put(`/tags/${id}`, data)
export const deleteTagAPI = (id) => api.delete(`/tags/${id}`)
export const getBloggerTagsAPI = (bloggerId) => api.get(`/blogger/${bloggerId}/tags`)
export const setBloggerTagsAPI = (bloggerId, tagIds) => api.post(`/blogger/${bloggerId}/tags`, { tag_ids: tagIds })
export const getBloggersByTagAPI = (tagId, page = 1, pageSize = 20) => api.get('/bloggers/by-tag', { params: { tag_id: tagId, page, pageSize } })

export const createCalendarEventAPI = (data) => api.post('/calendar/events', data)

export const getCalendarEventsAPI = (start, end) => api.get('/calendar/events', { params: { start, end } })

export const deleteCalendarEventAPI = (id) => api.delete(`/calendar/events/${id}`)

export const getTemplatesAPI = () => api.get('/templates')
export const createTemplateAPI = (data) => api.post('/templates', data)
export const updateTemplateAPI = (id, data) => api.put(`/templates/${id}`, data)
export const deleteTemplateAPI = (id) => api.delete(`/templates/${id}`)
export const useTemplateAPI = (id) => api.post(`/templates/${id}/use`)
export const setDefaultTemplateAPI = (id) => api.post(`/templates/${id}/set-default`)

export const getAnalyticsOverviewAPI = () => api.get('/analytics/overview')
export const getBloggerAnalyticsAPI = (bloggerId) => api.get(`/analytics/blogger/${bloggerId}`)

export const getTeamMembersAPI = () => api.get('/team/members')
export const requestBloggerTransferAPI = (data) => api.post('/blogger/transfer/request', data)
export const getTransferRequestsAPI = (status) => api.get('/blogger/transfer/requests', { params: { status } })
export const handleTransferRequestAPI = (data) => api.post('/blogger/transfer/handle', data)
export const getTeamTasksAPI = () => api.get('/team/tasks')
export const createTeamTaskAPI = (data) => api.post('/team/tasks', data)
export const updateTeamTaskStatusAPI = (data) => api.put('/team/tasks/status', data)

export const previewImportAPI = (formData) => api.post('/import/preview', formData, {
  headers: { 'Content-Type': 'multipart/form-data' }
})
export const doImportAPI = (data) => api.post('/import/do', data)
export const downloadTemplateAPI = (format) => api.get('/import/template', { params: { format }, responseType: 'blob' })
export const getImportHistoryAPI = () => api.get('/import/history')
export const exportFailedRowsAPI = (errors) => api.post('/import/export-errors', errors, { responseType: 'blob' })

export default api