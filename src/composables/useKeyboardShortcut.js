import { onMounted, onUnmounted } from 'vue'

export function useKeyboardShortcut(keyMap, options = {}) {
  const { preventDefault = true, target = null } = options

  const pressedKeys = new Set()

  const parseShortcut = (shortcut) => {
    return shortcut
      .toLowerCase()
      .split('+')
      .map((k) => k.trim())
  }

  const matchShortcut = (event, shortcut) => {
    const keys = parseShortcut(shortcut)
    const eventKey = event.key.toLowerCase()

    const requiredModifiers = {
      ctrl: keys.includes('ctrl') || keys.includes('cmd'),
      alt: keys.includes('alt'),
      shift: keys.includes('shift'),
      meta: keys.includes('meta') || keys.includes('cmd')
    }

    const mainKey = keys.find(
      (k) => !['ctrl', 'alt', 'shift', 'meta', 'cmd'].includes(k)
    )

    const modifierMatch =
      requiredModifiers.ctrl === (event.ctrlKey || event.metaKey) &&
      requiredModifiers.alt === event.altKey &&
      requiredModifiers.shift === event.shiftKey &&
      requiredModifiers.meta === event.metaKey

    if (!modifierMatch) return false

    if (!mainKey) return true

    const keyAliases = {
      escape: 'escape',
      esc: 'escape',
      enter: 'enter',
      return: 'enter',
      space: ' ',
      up: 'arrowup',
      down: 'arrowdown',
      left: 'arrowleft',
      right: 'arrowright',
      '/': '/',
      '?': '?'
    }

    const normalizedMain = keyAliases[mainKey] || mainKey
    const normalizedEvent = keyAliases[eventKey] || eventKey

    return normalizedMain === normalizedEvent
  }

  const handleKeyDown = (event) => {
    if (event.target.tagName === 'INPUT' || event.target.tagName === 'TEXTAREA' || event.target.tagName === 'SELECT' || event.target.isContentEditable) {
      const inputShortcuts = Object.entries(keyMap).filter(([key]) => {
        const keys = parseShortcut(key)
        return keys.includes('/')
      })
      if (inputShortcuts.length === 0) return
      let matched = false
      for (const [shortcut, handler] of inputShortcuts) {
        if (matchShortcut(event, shortcut)) {
          if (preventDefault) event.preventDefault()
          handler(event)
          matched = true
          break
        }
      }
      if (!matched) return
      return
    }

    for (const [shortcut, handler] of Object.entries(keyMap)) {
      if (matchShortcut(event, shortcut)) {
        if (preventDefault) event.preventDefault()
        handler(event)
        break
      }
    }
  }

  onMounted(() => {
    const el = target?.value || target || document
    el.addEventListener('keydown', handleKeyDown)
  })

  onUnmounted(() => {
    const el = target?.value || target || document
    el.removeEventListener('keydown', handleKeyDown)
  })
}
