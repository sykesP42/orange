import { ref, computed, watch } from 'vue'

const STORAGE_KEY = 'orange_notification_prefs'
const DEFAULT_PREFS = {
  emailEnabled: false,
  emailAddress: '',
  pushEnabled: true,
  soundEnabled: true,
  quietHoursStart: '22:00',
  quietHoursEnd: '08:00',
  quietHoursEnabled: false,
  types: {
    mention: { email: true, push: true, sound: true },
    announcement: { email: true, push: true, sound: true },
    system: { email: true, push: true, sound: false },
    invalid_blogger: { email: true, push: true, sound: true },
    countdown: { email: true, push: true, sound: true },
    team_change: { email: true, push: true, sound: true },
    new_post: { email: false, push: true, sound: true },
    post_reply: { email: false, push: true, sound: true },
    post_like: { email: false, push: false, sound: false },
    team_message: { email: false, push: true, sound: true },
    blogger_transfer: { email: true, push: true, sound: true },
    team_join: { email: true, push: true, sound: true },
    team_leave: { email: true, push: true, sound: true },
    status_change: { email: true, push: true, sound: true },
    deadline_reminder: { email: true, push: true, sound: true },
    data_backup: { email: true, push: false, sound: false }
  }
}

let prefs = ref({ ...DEFAULT_PREFS })
let initialized = false

export function useNotificationPrefs() {
  if (!initialized) {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      try {
        const parsed = JSON.parse(stored)
        prefs.value = deepMerge(DEFAULT_PREFS, parsed)
      } catch {
        prefs.value = { ...DEFAULT_PREFS }
      }
    }
    initialized = true
  }

  function _reset() {
    initialized = false
    prefs.value = { ...DEFAULT_PREFS }
  }

  function save() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(prefs.value))
  }

  function updatePref(key, value) {
    prefs.value[key] = value
    save()
  }

  function updateTypePref(typeName, channel, value) {
    if (prefs.value.types[typeName]) {
      prefs.value.types[typeName][channel] = value
      save()
    }
  }

  function shouldNotify(typeName, channel) {
    const typeConfig = prefs.value.types[typeName]
    if (!typeConfig) return channel === 'push'
    return typeConfig[channel] !== false
  }

  function isQuietHours() {
    if (!prefs.value.quietHoursEnabled) return false
    const now = new Date()
    const currentMinutes = now.getHours() * 60 + now.getMinutes()
    const [startH, startM] = prefs.value.quietHoursStart.split(':').map(Number)
    const [endH, endM] = prefs.value.quietHoursEnd.split(':').map(Number)
    const startMinutes = startH * 60 + startM
    const endMinutes = endH * 60 + endM
    if (startMinutes > endMinutes) {
      return currentMinutes >= startMinutes || currentMinutes < endMinutes
    }
    return currentMinutes >= startMinutes && currentMinutes < endMinutes
  }

  function resetToDefaults() {
    prefs.value = { ...DEFAULT_PREFS, types: { ...DEFAULT_PREFS.types, ...Object.fromEntries(Object.keys(DEFAULT_PREFS.types).map(k => [k, { ...DEFAULT_PREFS.types[k] }])) } }
    save()
  }

  function getSummary() {
    const enabledTypes = Object.entries(prefs.value.types).filter(([, cfg]) => cfg.push || cfg.email).length
    return {
      totalTypes: Object.keys(prefs.value.types).length,
      enabledTypes,
      emailOn: prefs.value.emailEnabled,
      pushOn: prefs.value.pushEnabled,
      soundOn: prefs.value.soundEnabled,
      quietOn: prefs.value.quietHoursEnabled
    }
  }

  watch(() => ({ ...prefs.value }), () => save(), { deep: true })

  return {
    prefs,
    save,
    updatePref,
    updateTypePref,
    shouldNotify,
    isQuietHours,
    resetToDefaults,
    getSummary,
    _reset
  }
}

function deepMerge(target, source) {
  const result = { ...target }
  for (const key of Object.keys(source)) {
    if (source[key] && typeof source[key] === 'object' && !Array.isArray(source[key]) && target[key] && typeof target[key] === 'object') {
      result[key] = deepMerge(target[key], source[key])
    } else {
      result[key] = source[key]
    }
  }
  return result
}
