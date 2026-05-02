import { describe, it, expect, beforeEach } from 'vitest'
import { useNotificationPrefs } from '../../composables/useNotificationPrefs'
import { useDataBackup } from '../../composables/useDataBackup'
import { useAuditLog } from '../../composables/useAuditLog'

describe('Boundary: useNotificationPrefs edge cases', () => {
  let np

  beforeEach(() => {
    localStorage.removeItem('orange_notification_prefs')
    np = useNotificationPrefs()
    np._reset()
  })

  it('should handle corrupted localStorage data', () => {
    localStorage.setItem('orange_notification_prefs', 'not-json{{{')
    np._reset()
    np = useNotificationPrefs()
    expect(np.prefs.value.pushEnabled).toBe(true)
  })

  it('should handle partial prefs data', () => {
    localStorage.setItem('orange_notification_prefs', JSON.stringify({ emailEnabled: true }))
    np._reset()
    np = useNotificationPrefs()
    expect(np.prefs.value.emailEnabled).toBe(true)
    expect(np.prefs.value.pushEnabled).toBe(true)
    expect(np.prefs.value.types).toBeDefined()
  })

  it('should handle unknown notification type in shouldNotify', () => {
    expect(np.shouldNotify('nonexistent_type_xyz', 'push')).toBe(true)
    expect(np.shouldNotify('nonexistent_type_xyz', 'email')).toBe(false)
  })

  it('should handle midnight crossover in quiet hours', () => {
    np.updatePref('quietHoursEnabled', true)
    np.updatePref('quietHoursStart', '23:00')
    np.updatePref('quietHoursEnd', '06:00')
    const isQuiet = np.isQuietHours()
    expect(typeof isQuiet).toBe('boolean')
  })

  it('should survive rapid pref updates without error', () => {
    for (let i = 0; i < 100; i++) {
      np.updatePref('pushEnabled', i % 2 === 0)
      np.updateTypePref('mention', 'push', i % 3 === 0)
      np.updateTypePref('mention', 'email', i % 4 === 0)
    }
    expect(typeof np.prefs.value.pushEnabled).toBe('boolean')
  })

  it('should handle all channels disabled on a type', () => {
    np.updateTypePref('post_like', 'push', false)
    np.updateTypePref('post_like', 'email', false)
    np.updateTypePref('post_like', 'sound', false)
    expect(np.shouldNotify('post_like', 'push')).toBe(false)
    expect(np.shouldNotify('post_like', 'email')).toBe(false)
    expect(np.shouldNotify('post_like', 'sound')).toBe(false)
  })
})

describe('Boundary: useDataBackup edge cases', () => {
  let backup

  beforeEach(() => {
    Object.keys(localStorage).forEach(key => {
      if (key.startsWith('orange_backup_')) localStorage.removeItem(key)
    })
    backup = useDataBackup()
    backup._reset()
  })

  it('should handle empty object backup', async () => {
    const result = await backup.createBackup({})
    expect(result.id).toMatch(/^bk_/)
    expect(result.recordCount).toBe(0)
  })

  it('should handle very large backup data', async () => {
    const largeData = { bloggers: Array(1000).fill({ id: 1, name: 'test' }) }
    const result = await backup.createBackup(largeData)
    expect(result.recordCount).toBe(1000)
  })

  it('should handle special characters in backup label', async () => {
    const result = await backup.createBackup({}, '<script>alert(1)</script>"&\'')
    expect(result.label).toBeTruthy()
  })

  it('should handle corrupted localStorage meta', () => {
    localStorage.setItem('orange_backup_meta', '{corrupted')
    backup._reset()
    backup = useDataBackup()
    expect(Array.isArray(backup.backupList.value)).toBe(true)
  })

  it('should handle restore of non-JSON backup storage', async () => {
    localStorage.setItem('orange_backup_test_corrupt', 'not-json-data')
    await expect(backup.restoreFromLocal('test_corrupt')).rejects.toThrow()
  })

  it('should enforce max backups exactly at limit', async () => {
    for (let i = 0; i < 10; i++) {
      await backup.createBackup({ index: i })
    }
    expect(backup.backupList.value.length).toBe(10)
    const tenthId = backup.backupList.value[9].id
    await backup.createBackup({ index: 10 })
    expect(backup.backupList.value.length).toBe(10)
    const newIds = backup.backupList.value.map(b => b.id)
    expect(newIds).not.toContain(tenthId)
  })

  it('should handle delete of non-existent backup gracefully', () => {
    expect(() => backup.deleteBackup('non_existent_123')).not.toThrow()
  })
})

describe('Boundary: useAuditLog edge cases', () => {
  beforeEach(() => {
    localStorage.removeItem('orange_audit_log')
    const auditLog = useAuditLog()
    auditLog.clearLogs()
  })

  it('should handle very long action details', () => {
    const auditLog = useAuditLog()
    const longDetails = 'x'.repeat(10000)
    auditLog.log('test_action', longDetails)
    const logs = auditLog.getLogs()
    expect(logs[0].details).toBe(longDetails)
  })

  it('should handle special characters in log details', () => {
    const auditLog = useAuditLog()
    auditLog.log('test', '<script>alert(1)</script>&"\'')
    const logs = auditLog.getLogs()
    expect(logs.length).toBe(1)
  })

  it('should handle getStats with no logs', () => {
    const auditLog = useAuditLog()
    const stats = auditLog.getStats()
    expect(stats.totalLogs || 0).toBe(0)
    expect(Object.keys(stats.byAction || {})).toHaveLength(0)
  })

  it('should respect MAX_LOGS limit', () => {
    const auditLog = useAuditLog()
    for (let i = 0; i < 600; i++) {
      auditLog.log(`action_${i}`, `detail_${i}`)
    }
    const logs = auditLog.getLogs()
    expect(logs.length).toBeLessThanOrEqual(500)
  })

  it('should handle filter with non-matching criteria', () => {
    const auditLog = useAuditLog()
    auditLog.log('create', 'test')
    const filtered = auditLog.getLogs({ action: 'nonexistent' })
    expect(filtered).toHaveLength(0)
  })
})
