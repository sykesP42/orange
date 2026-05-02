import { ref, computed } from 'vue'

const STORAGE_KEY = 'orange_backup_meta'
const BACKUP_PREFIX = 'orange_backup_'
const MAX_BACKUPS = 10

let backupList = ref([])
let isBackingUp = ref(false)
let isRestoring = ref(false)
let lastBackupTime = ref(null)

export function useDataBackup() {
  function loadMeta() {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      try {
        backupList.value = JSON.parse(stored)
      } catch {
        backupList.value = []
      }
    }
    if (backupList.value.length > 0) {
      lastBackupTime.value = backupList.value[0].timestamp
    }
  }

  function saveMeta() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(backupList.value))
  }

  async function createBackup(data, label = '手动备份') {
    if (!data || typeof data !== 'object') throw new Error('备份数据无效')
    isBackingUp.value = true
    try {
      const backupId = `bk_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
      const timestamp = new Date().toISOString()
      const backup = {
        id: backupId,
        label,
        timestamp,
        version: '1.0.0',
        dataSize: new Blob([JSON.stringify(data)]).size,
        recordCount: Array.isArray(data.bloggers) ? data.bloggers.length : Object.keys(data).length
      }
      localStorage.setItem(`${BACKUP_PREFIX}${backupId}`, JSON.stringify({ meta: backup, data }))
      backupList.value.unshift(backup)
      if (backupList.value.length > MAX_BACKUPS) {
        const removed = backupList.value.pop()
        localStorage.removeItem(`${BACKUP_PREFIX}${removed.id}`)
      }
      saveMeta()
      lastBackupTime.value = timestamp
      return backup
    } finally {
      isBackingUp.value = false
    }
  }

  async function downloadBackup(data, filename = null) {
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename || `orange_backup_${new Date().toISOString().slice(0, 19).replace(/[T:]/g, '-')}.json`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  async function restoreFromLocal(backupId) {
    isRestoring.value = true
    try {
      const stored = localStorage.getItem(`${BACKUP_PREFIX}${backupId}`)
      if (!stored) throw new Error('备份不存在')
      const parsed = JSON.parse(stored)
      return parsed.data
    } finally {
      isRestoring.value = false
    }
  }

  async function importFromFile(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = (e) => {
        try {
          const data = JSON.parse(e.target.result)
          if (!data || typeof data !== 'object') throw new Error('文件格式错误')
          resolve(data)
        } catch (err) {
          reject(new Error('解析备份文件失败: ' + err.message))
        }
      }
      reader.onerror = () => reject(new Error('读取文件失败'))
      reader.readAsText(file)
    })
  }

  function deleteBackup(backupId) {
    localStorage.removeItem(`${BACKUP_PREFIX}${backupId}`)
    backupList.value = backupList.value.filter(b => b.id !== backupId)
    saveMeta()
  }

  function clearAllBackups() {
    for (const b of backupList.value) {
      localStorage.removeItem(`${BACKUP_PREFIX}${b.id}`)
    }
    backupList.value = []
    saveMeta()
    lastBackupTime.value = null
  }

  function formatSize(bytes) {
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  }

  function timeAgo(isoString) {
    if (!isoString) return ''
    const date = new Date(isoString)
    const now = new Date()
    const diff = now - date
    if (diff < 60000) return '刚刚'
    if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
    return `${Math.floor(diff / 86400000)}天前`
  }

  const stats = computed(() => ({
    totalBackups: backupList.value.length,
    totalSize: backupList.value.reduce((sum, b) => sum + (b.dataSize || 0), 0),
    lastBackup: lastBackupTime.value ? timeAgo(lastBackupTime.value) : '无',
    maxBackups: MAX_BACKUPS
  }))

  if (backupList.value.length === 0) {
    loadMeta()
  }

  function _reset() {
    backupList.value = []
    isBackingUp.value = false
    isRestoring.value = false
    lastBackupTime.value = null
  }

  return {
    backupList,
    isBackingUp,
    isRestoring,
    lastBackupTime,
    stats,
    createBackup,
    downloadBackup,
    restoreFromLocal,
    importFromFile,
    deleteBackup,
    clearAllBackups,
    formatSize,
    timeAgo,
    loadMeta,
    _reset
  }
}
