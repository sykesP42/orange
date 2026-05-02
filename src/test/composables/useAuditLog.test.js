import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { useAuditLog } from '../../composables/useAuditLog'

describe('useAuditLog', () => {
  let auditLog

  beforeEach(() => {
    localStorage.clear()
    auditLog = useAuditLog()
    auditLog.clearLogs()
  })

  it('should log an action with details', () => {
    const entry = auditLog.log('test_action', { key: 'value' })

    expect(entry).toMatchObject({
      action: 'test_action',
      details: { key: 'value' }
    })
    expect(entry.id).toBeDefined()
    expect(entry.timestamp).toBeDefined()
    expect(entry.user).toBe('unknown')
  })

  it('should log CRUD operations', () => {
    auditLog.logCreate('blogger', '测试博主', '1')
    auditLog.logUpdate('blogger', '测试博主', '1', { status: '已合作' })
    auditLog.logDelete('blogger', '测试博主', '1')

    const logs = auditLog.getLogs()
    expect(logs).toHaveLength(3)
    expect(logs[0].action).toBe('delete')
    expect(logs[0].details.entityType).toBe('blogger')
    expect(logs[2].action).toBe('create')
  })

  it('should log batch operations', () => {
    auditLog.logBatch('status', 'blogger', 5, { newStatus: '洽谈中' })
    auditLog.logBatch('tag', 'blogger', 3, { tag: '公' })

    const logs = auditLog.getLogs({ action: 'batch_status' })
    expect(logs).toHaveLength(1)
    expect(logs[0].details.count).toBe(5)
    expect(logs[0].details.newStatus).toBe('洽谈中')
  })

  it('should log export operations', () => {
    auditLog.logExport('xlsx', 42)

    const logs = auditLog.getLogs({ action: 'export' })
    expect(logs).toHaveLength(1)
    expect(logs[0].details.format).toBe('xlsx')
    expect(logs[0].details.count).toBe(42)
  })

  it('should filter logs by action prefix', () => {
    auditLog.logBatch('status', 'blogger', 1)
    auditLog.logBatch('tag', 'blogger', 1)
    auditLog.logBatch('delete', 'blogger', 1)
    auditLog.logExport('xlsx', 10)

    const batchLogs = auditLog.getLogs({ action: 'batch' })
    expect(batchLogs).toHaveLength(3)
  })

  it('should respect limit parameter', () => {
    for (let i = 0; i < 10; i++) {
      auditLog.log(`action_${i}`)
    }

    const limited = auditLog.getLogs({ limit: 3 })
    expect(limited).toHaveLength(3)
  })

  it('should return stats correctly', () => {
    auditLog.logCreate('blogger', 'A', '1')
    auditLog.logCreate('blogger', 'B', '2')
    auditLog.logUpdate('blogger', 'A', '1', {})
    auditLog.logDelete('blogger', 'C', '3')

    const stats = auditLog.getStats()
    expect(stats.total).toBe(4)
    expect(stats.byAction.create).toBe(2)
    expect(stats.byAction.update).toBe(1)
    expect(stats.byAction.delete).toBe(1)
  })

  it('should clear all logs', () => {
    auditLog.log('a')
    auditLog.log('b')
    auditLog.log('c')

    expect(auditLog.getStats().total).toBe(3)
    auditLog.clearLogs()
    expect(auditLog.getStats().total).toBe(0)
  })

  it('should persist to localStorage', () => {
    auditLog.log('persist_test')

    const stored = JSON.parse(localStorage.getItem('orange_audit_log'))
    expect(stored).toHaveLength(1)
    expect(stored[0].action).toBe('persist_test')
  })

  it('should cap at MAX_LOGS entries', () => {
    for (let i = 0; i < 600; i++) {
      auditLog.log(`overflow_${i}`)
    }

    const stats = auditLog.getStats()
    expect(stats.total).toBeLessThanOrEqual(500)
  })
})
