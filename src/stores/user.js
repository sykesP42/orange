import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref('')
  const username = ref('')
  const realName = ref('')
  const role = ref('user')
  const id = ref(null)
  const team_id = ref(null)
  const team_name = ref(null)
  const team_color = ref(null)
  const user_code = ref('')
  const invited_by = ref('')

  const isLoggedIn = ref(false)
  const isAdmin = ref(false)
  const _roleVerified = ref(false)

  function setUser(data) {
    token.value = data.token
    username.value = data.username
    realName.value = data.real_name
    role.value = data.role
    id.value = data.id
    team_id.value = data.team_id || null
    team_name.value = data.team_name || null
    team_color.value = data.team_color || null
    user_code.value = data.user_code || ''
    invited_by.value = data.invited_by || ''
    isLoggedIn.value = true
    isAdmin.value = data.role === 'admin'
    _roleVerified.value = true

    localStorage.setItem('token', data.token)
    localStorage.setItem('username', data.username)
    localStorage.setItem('realName', data.real_name)
    localStorage.setItem('role', data.role)
    localStorage.setItem('id', data.id)
    localStorage.setItem('team_id', data.team_id || '')
    localStorage.setItem('team_name', data.team_name || '')
    localStorage.setItem('team_color', data.team_color || '')
    localStorage.setItem('user_code', data.user_code || '')
    localStorage.setItem('invited_by', data.invited_by || '')
  }

  function logout() {
    token.value = ''
    username.value = ''
    realName.value = ''
    role.value = 'user'
    id.value = null
    team_id.value = null
    team_name.value = null
    team_color.value = null
    user_code.value = ''
    invited_by.value = ''
    isLoggedIn.value = false
    isAdmin.value = false
    _roleVerified.value = false

    localStorage.removeItem('token')
    localStorage.removeItem('username')
    localStorage.removeItem('realName')
    localStorage.removeItem('role')
    localStorage.removeItem('id')
    localStorage.removeItem('team_id')
    localStorage.removeItem('team_name')
    localStorage.removeItem('team_color')
    localStorage.removeItem('user_code')
    localStorage.removeItem('invited_by')
  }

  function updateTeam(teamId, teamName, teamColor) {
    team_id.value = teamId
    team_name.value = teamName
    team_color.value = teamColor
    localStorage.setItem('team_id', teamId || '')
    localStorage.setItem('team_name', teamName || '')
    localStorage.setItem('team_color', teamColor || '')
  }

  return {
    token,
    username,
    realName,
    role,
    id,
    team_id,
    team_name,
    team_color,
    user_code,
    invited_by,
    isLoggedIn,
    isAdmin,
    _roleVerified,
    setUser,
    logout,
    updateTeam
  }
})
