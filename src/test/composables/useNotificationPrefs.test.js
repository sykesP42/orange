import { describe, it, expect, beforeEach } from 'vitest'
import { useNotificationPrefs } from '../../composables/useNotificationPrefs'

describe('useNotificationPrefs', () => {
  let np

  beforeEach(() => {
    localStorage.removeItem('orange_notification_prefs')
    np = useNotificationPrefs()
    np._reset()
    np = useNotificationPrefs()
  })

  it('should initialize with default preferences', () => {
    expect(np.prefs.value.emailEnabled).toBe(false)
    expect(np.prefs.value.pushEnabled).toBe(true)
    expect(np.prefs.value.soundEnabled).toBe(true)
    expect(np.prefs.value.quietHoursEnabled).toBe(false)
    expect(np.prefs.value.quietHoursStart).toBe('22:00')
    expect(np.prefs.value.quietHoursEnd).toBe('08:00')
    expect(Object.keys(np.prefs.value.types)).toContain('mention')
    expect(Object.keys(np.prefs.value.types)).toContain('status_change')
    expect(Object.keys(np.prefs.value.types)).toContain('data_backup')
  })

  it('should persist preferences to localStorage', () => {
    np.updatePref('emailEnabled', true)
    np.updatePref('soundEnabled', false)
    const stored = JSON.parse(localStorage.getItem('orange_notification_prefs'))
    expect(stored.emailEnabled).toBe(true)
    expect(stored.soundEnabled).toBe(false)
  })

  it('should load saved preferences from localStorage', () => {
    localStorage.setItem('orange_notification_prefs', JSON.stringify({
      emailEnabled: true,
      emailAddress: 'test@example.com',
      pushEnabled: false
    }))
    np._reset()
    np = useNotificationPrefs()
    expect(np.prefs.value.emailEnabled).toBe(true)
    expect(np.prefs.value.emailAddress).toBe('test@example.com')
    expect(np.prefs.value.pushEnabled).toBe(false)
  })

  it('should update type-specific preferences', () => {
    np.updateTypePref('mention', 'email', false)
    expect(np.prefs.value.types.mention.email).toBe(false)
    np.updateTypePref('deadline_reminder', 'sound', true)
    expect(np.prefs.value.types.deadline_reminder.sound).toBe(true)
  })

  it('should return correct shouldNotify result', () => {
    expect(np.shouldNotify('mention', 'push')).toBe(true)
    expect(np.shouldNotify('post_like', 'push')).toBe(false)
    expect(np.shouldNotify('unknown_type', 'push')).toBe(true)
  })

  it('should detect quiet hours correctly', () => {
    expect(np.isQuietHours()).toBe(false)
    np.updatePref('quietHoursEnabled', true)
    np.updatePref('quietHoursStart', '00:00')
    np.updatePref('quietHoursEnd', '23:59')
    expect(np.isQuietHours()).toBe(true)
  })

  it('should reset to defaults', () => {
    np.updatePref('emailEnabled', true)
    np.updatePref('pushEnabled', false)
    np.resetToDefaults()
    expect(np.prefs.value.emailEnabled).toBe(false)
    expect(np.prefs.value.pushEnabled).toBe(true)
  })

  it('should provide summary stats', () => {
    const summary = np.getSummary()
    expect(summary.totalTypes).toBe(16)
    expect(typeof summary.enabledTypes).toBe('number')
    expect(summary.emailOn).toBe(false)
    expect(summary.pushOn).toBe(true)
    expect(summary.soundOn).toBe(true)
  })

  it('should handle all 16 notification types with defaults', () => {
    const expectedTypes = [
      'mention', 'announcement', 'system', 'invalid_blogger', 'countdown',
      'team_change', 'new_post', 'post_reply', 'post_like', 'team_message',
      'blogger_transfer', 'team_join', 'team_leave',
      'status_change', 'deadline_reminder', 'data_backup'
    ]
    for (const t of expectedTypes) {
      expect(np.prefs.value.types[t]).toBeDefined()
      expect(typeof np.prefs.value.types[t].push).toBe('boolean')
      expect(typeof np.prefs.value.types[t].email).toBe('boolean')
      expect(typeof np.prefs.value.types[t].sound).toBe('boolean')
    }
  })
})
