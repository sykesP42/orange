import { describe, it, expect, beforeEach } from 'vitest'
import { useDataBackup } from '../../composables/useDataBackup'

describe('useDataBackup', () => {
  let backup

  beforeEach(() => {
    Object.keys(localStorage).forEach(key => {
      if (key.startsWith('orange_backup_')) {
        localStorage.removeItem(key)
      }
    })
    backup = useDataBackup()
    backup._reset()
  })

  it('should initialize with empty backup list', () => {
    expect(backup.backupList.value).toHaveLength(0)
    expect(backup.isBackingUp.value).toBe(false)
    expect(backup.isRestoring.value).toBe(false)
  })

  it('should create a backup and store it', async () => {
    const data = { bloggers: [{ id: 1, name: 'test' }], timestamp: new Date().toISOString() }
    const result = await backup.createBackup(data)
    expect(result.id).toMatch(/^bk_/)
    expect(result.label).toBe('手动备份')
    expect(result.recordCount).toBe(1)
    expect(backup.backupList.value).toHaveLength(1)
    expect(backup.lastBackupTime.value).toBeTruthy()
  })

  it('should store backup data in localStorage', async () => {
    const data = { bloggers: [], test: true }
    await backup.createBackup(data)
    const keys = Object.keys(localStorage).filter(k => k.startsWith('orange_backup_'))
    expect(keys.length).toBeGreaterThanOrEqual(2)
  })

  it('should enforce max backups limit (10)', async () => {
    for (let i = 0; i < 12; i++) {
      await backup.createBackup({ index: i })
    }
    expect(backup.backupList.value.length).toBeLessThanOrEqual(10)
  })

  it('should delete a specific backup', async () => {
    const data = { test: 1 }
    const result = await backup.createBackup(data)
    expect(backup.backupList.value).toHaveLength(1)
    backup.deleteBackup(result.id)
    expect(backup.backupList.value).toHaveLength(0)
    expect(localStorage.getItem(`orange_backup_${result.id}`)).toBeNull()
  })

  it('should clear all backups', async () => {
    for (let i = 0; i < 3; i++) {
      await backup.createBackup({ i })
    }
    backup.clearAllBackups()
    expect(backup.backupList.value).toHaveLength(0)
    expect(backup.lastBackupTime.value).toBeNull()
  })

  it('should restore from local storage', async () => {
    const originalData = { bloggers: [{ id: 42 }], source: 'test' }
    const bk = await backup.createBackup(originalData)
    const restored = await backup.restoreFromLocal(bk.id)
    expect(restored.bloggers).toHaveLength(1)
    expect(restored.source).toBe('test')
  })

  it('should import from a JSON file', async () => {
    const fileContent = JSON.stringify({ imported: true, count: 99 })
    const file = new File([fileContent], 'test.json', { type: 'application/json' })
    const data = await backup.importFromFile(file)
    expect(data.imported).toBe(true)
    expect(data.count).toBe(99)
  })

  it('should reject invalid import files', async () => {
    const file = new File(['not json'], 'bad.txt', { type: 'text/plain' })
    await expect(backup.importFromFile(file)).rejects.toThrow()
  })

  it('should format file sizes correctly', () => {
    expect(backup.formatSize(500)).toBe('500 B')
    expect(backup.formatSize(2048)).toContain('KB')
    expect(backup.formatSize(3 * 1024 * 1024)).toContain('MB')
  })

  it('should compute stats correctly', async () => {
    expect(backup.stats.value.totalBackups).toBe(0)
    await backup.createBackup({ a: 1 })
    expect(backup.stats.value.totalBackups).toBe(1)
    expect(backup.stats.value.maxBackups).toBe(10)
  })

  it('should handle timeAgo formatting', () => {
    const now = new Date().toISOString()
    expect(backup.timeAgo(now)).toBe('刚刚')
    const past = new Date(Date.now() - 5 * 60000).toISOString()
    expect(backup.timeAgo(past)).toContain('分钟')
    const older = new Date(Date.now() - 3 * 3600000).toISOString()
    expect(backup.timeAgo(older)).toContain('小时')
  })

  it('should reject invalid backup data', async () => {
    await expect(backup.createBackup(null)).rejects.toThrow()
    await expect(backup.createBackup('string')).rejects.toThrow()
  })

  it('should fail to restore non-existent backup', async () => {
    await expect(backup.restoreFromLocal('non_existent_id')).rejects.toThrow()
  })

  it('should support custom labels on createBackup', async () => {
    const result = await backup.createBackup({}, '自动备份')
    expect(result.label).toBe('自动备份')
  })
})
