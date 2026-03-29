import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

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
    return response.data
  },
  error => {
    if (error.response?.status === 401) {
      localStorage.clear()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const loginAPI = (data) => api.post('/login', data)

export const registerAPI = (data) => api.post('/register', data)

export const forgotPasswordAPI = (data) => api.post('/forgot-password', data)

export const bloggerListAPI = (params) => api.get('/blogger/list', { params })

export const bloggerMyAPI = (params) => api.get('/blogger/my', { params })

export const getMyBloggersAPI = (params) => api.get('/blogger/my', { params })

export const bloggerAddAPI = (data) => api.post('/blogger/add', data)

export const bloggerDeleteAPI = (id) => api.post('/admin/blogger/delete', { id })

export const bloggerStatAPI = () => api.get('/blogger/stat')

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

export const categoryAddAPI = (data) => api.post('/category/add', data)

export const categoryUpdateAPI = (data) => api.post('/category/update', data)

export const categoryDeleteAPI = (id) => api.post('/category/delete', { id })

export const platformListAPI = () => api.get('/platform/list')

export const productListAPI = () => api.get('/product/list')

export const productAddAPI = (data) => api.post('/product/add', data)

export const productUpdateAPI = (data) => api.put('/product/update', data)

export const productDeleteAPI = (id) => api.delete(`/product/${id}`)

export const userListAPI = () => api.get('/users/public')

export const deactivateUserAPI = (id) => api.post('/admin/users/deactivate', { id })

export const deleteUserAPI = (id) => api.post('/admin/users/delete', { id })

export const exportAPI = (type) => api.get('/admin/export', { params: { type }, responseType: 'blob' })

export const getOperationLog = (params) => api.get('/admin/log', { params })

export const editBlogger = (data) => api.post('/admin/blogger/edit', data)

export const getBloggerStatusList = () => api.get('/blogger/status/list')

export const getBloggerDetail = (id) => api.get(`/blogger/${id}`)

export const getBloggerHotData = (nickname) => api.get(`/blogger/${nickname}/hot`)

export const getBloggerNews = (nickname) => api.get(`/blogger/${nickname}/news`)

export const getBloggerStatusHistory = (id) => api.get(`/blogger/${id}/status/history`)

export const getBloggerCharts = () => api.get('/blogger/charts')

export const requestBloggerTransfer = (data) => api.post('/blogger/transfer/request', data)

export const confirmBloggerTransfer = (data) => api.post('/blogger/transfer/confirm', data)

export const adminConfirmBloggerTransfer = (data) => api.post('/blogger/transfer/admin-confirm', data)

export const getBloggerTransferRequests = () => api.get('/blogger/transfer/requests')

export const getPendingBloggerTransfers = () => api.get('/admin/blogger/transfer/pending')

export const deleteBlogger = (id) => api.post('/admin/blogger/delete', { id })

export const exportData = (type) => api.get('/admin/export', { params: { type }, responseType: 'blob' })

// 用户设置相关API
export const getUserProfileAPI = () => api.get('/user/profile')

export const updateUserProfileAPI = (data) => api.put('/user/profile', data)

export const updatePasswordAPI = (data) => api.put('/user/password', data)

export const uploadAvatarAPI = (data) => api.post('/upload/avatar', data)

export const getPublicUsersAPI = () => api.get('/users/public')

export const getPendingUsersAPI = () => api.get('/users/pending')

export const batchSetAdminAPI = (teamId, userIds, action) => api.post(`/team/${teamId}/admins/batch`, { user_ids: userIds, action })

export const getPublicForumsAPI = () => api.get('/public/forums')

export const getPublicForumDetailAPI = (forumId) => api.get(`/public/forums/${forumId}`)

export const approveUserAPI = (id) => api.post('/users/approve', { id })

export const rejectUserAPI = (id) => api.post('/users/reject', { id })

export const setAdminAPI = (id) => api.post('/admin/users/setadmin', { id })
export const removeAdminAPI = (id) => api.post('/admin/users/removeadmin', { id })

export const getTeamsAPI = () => api.get('/teams')
export const createTeamAPI = (data) => api.post('/teams', data)
export const updateTeamAPI = (id, data) => api.put(`/teams/${id}`, data)
export const deleteTeamAPI = (id) => api.delete(`/teams/${id}`)
export const setUserTeamAPI = (id, team_id) => api.put(`/users/${id}/team`, { team_id })
export const updateMyTeamAPI = (team_id) => api.put(`/my/team`, { team_id })
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
export const getTeamPostDetailAPI = (teamId, postId) => api.get(`/team/${teamId}/posts/${postId}`)
export const deleteTeamPostAPI = (teamId, postId) => api.delete(`/team/${teamId}/posts/${postId}`)
export const createTeamPostCommentAPI = (teamId, postId, data) => api.post(`/team/${teamId}/posts/${postId}/comments`, data)
export const likeTeamPostAPI = (teamId, postId) => api.post(`/team/${teamId}/posts/${postId}/like`)
export const getTeamPostLikeStatusAPI = (teamId, postId) => api.get(`/team/${teamId}/posts/${postId}/like-status`)
export const pinTeamPostAPI = (teamId, postId) => api.put(`/team/${teamId}/posts/${postId}/pin`)
export const featureTeamPostAPI = (teamId, postId) => api.put(`/team/${teamId}/posts/${postId}/feature`)
export const searchTeamPostsAPI = (teamId, keyword) => api.get(`/team/${teamId}/posts/search`, { params: { keyword } })
export const collectTeamPostAPI = (teamId, postId) => api.post(`/team/${teamId}/posts/${postId}/collect`)
export const getTeamPostCollectStatusAPI = (teamId, postId) => api.get(`/team/${teamId}/posts/collect-status/${postId}`)
export const getTeamPostsCollectedAPI = (teamId) => api.get(`/team/${teamId}/posts/collected`)
export const setTeamAdminAPI = (teamId, userId) => api.put(`/team/${teamId}/admins/${userId}`)
export const getTeamMessagesAPI = (teamId) => api.get(`/team/${teamId}/messages`)
export const sendTeamMessageAPI = (teamId, data) => api.post(`/team/${teamId}/messages`, data)
export const markMessageReadAPI = (teamId, messageId) => api.put(`/team/${teamId}/messages/${messageId}/read`)
export const deleteTeamMessageAPI = (teamId, messageId) => api.delete(`/team/${teamId}/messages/${messageId}`)
export const getTeamOperationLogsAPI = (teamId, page = 1, pageSize = 50) => api.get(`/team/${teamId}/logs?page=${page}&pageSize=${pageSize}`)

export const uploadBloggerAvatarAPI = (data) => api.post('/upload/blogger-avatar', data)

export const getUserDetailAPI = (username) => api.get(`/user/detail/${username}`)

export const getUserBloggersAPI = (username) => api.get(`/user/${username}/bloggers`)

export const getUserFollowupStatsAPI = (username) => api.get(`/user/${username}/followup-stats`)

export const getNotificationsAPI = () => api.get('/notifications')

export const markNotificationReadAPI = (id) => api.put(`/notifications/${id}/read`)

export const markAllNotificationsReadAPI = () => api.put('/notifications/read-all')

export const invalidBloggerListAPI = (params) => api.get('/blogger/invalid/list', { params })

export const bindInvalidBloggerAPI = (id) => api.post('/blogger/invalid/bind', { id })

export const getInvalidBloggerCountAPI = () => api.get('/blogger/invalid/count')

export const submitBloggerEvaluationAPI = (data) => api.post('/blogger/evaluation', data)

export const getBloggerEvaluationAPI = (blogger_id) => api.get(`/blogger/evaluation/${blogger_id}`)

export const getExpiringBloggersWithoutContactAPI = () => api.get('/blogger/expiring-without-contact')

export const getAnnouncementsAPI = () => api.get('/announcements')

export const createAnnouncementAPI = (data) => api.post('/announcements', data)

export default api