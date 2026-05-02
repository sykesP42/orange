import { ref } from 'vue'

const STORAGE_KEY = 'orange_audit_log'
const MAX_LOGS = 500

const logs = ref([])
let initialized = false

function ensureInit() {
  if (initialized) return
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    logs.value = stored ? JSON.parse(stored) : []
  } catch { logs.value = [] }
  initialized = true
}

function persist() {
  try {
    const toSave = logs.value.slice(-MAX_LOGS)
    localStorage.setItem(STORAGE_KEY, JSON.stringify(toSave))
  } catch {}
}

export function useAuditLog() {
  ensureInit()

  function log(action, details = {}) {
    const entry = {
      id: Date.now().toString(36) + Math.random().toString(36).slice(2, 7),
      action,
      details,
      timestamp: new Date().toISOString(),
      user: localStorage.getItem('username') || 'unknown',
      ua: navigator.userAgent?.slice(0, 120) || ''
    }
    logs.value.push(entry)
    if (logs.value.length > MAX_LOGS) {
      logs.value = logs.value.slice(-MAX_LOGS)
    }
    persist()
    return entry
  }

  function logCreate(entityType, entityName, entityId) {
    return log('create', { entityType, entityName, entityId })
  }

  function logUpdate(entityType, entityName, entityId, changes) {
    return log('update', { entityType, entityName, entityId, changes })
  }

  function logDelete(entityType, entityName, entityId) {
    return log('delete', { entityType, entityName, entityId })
  }

  function logBatch(action, entityType, count, extra = {}) {
    return log(`batch_${action}`, { entityType, count, ...extra })
  }

  function logExport(format, count) {
    return log('export', { format, count })
  }

  function logLogin(success, method = 'password') {
    return log(success ? 'login' : 'login_failed', { method })
  }

  function getLogs(options = {}) {
    const { action, limit = 100, since, until } = options
    let result = [...logs.value]
    if (action) result = result.filter(l => l.action === action || l.action?.startsWith(action))
    if (since) result = result.filter(l => new Date(l.timestamp) >= new Date(since))
    if (until) result = result.filter(l => new Date(l.timestamp) <= new Date(until))
    return result.slice(-limit).reverse()
  }

  function getStats() {
    const now = new Date()
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
    const weekAgo = new Date(today.getTime() - 7 * 86400000)

    return {
      total: logs.value.length,
      today: logs.value.filter(l => new Date(l.timestamp) >= today).length,
      thisWeek: logs.value.filter(l => new Date(l.timestamp) >= weekAgo).length,
      byAction: logs.value.reduce((acc, l) => {
        acc[l.action] = (acc[l.action] || 0) + 1
        return acc
      }, {})
    }
  }

  function clearLogs() {
    logs.value = []
    persist()
  }

  return {
    logs,
    log,
    logCreate,
    logUpdate,
    logDelete,
    logBatch,
    logExport,
    logLogin,
    getLogs,
    getStats,
    clearLogs
  }
}
