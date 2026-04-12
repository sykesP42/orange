import { onMounted, onUnmounted } from 'vue'

let focusTrapStack = 0
let trappedElement = null
let previouslyFocused = null

export function useA11y() {
  const trapFocus = (containerEl) => {
    if (!containerEl) return () => {}
    if (focusTrapStack === 0) {
      previouslyFocused = document.activeElement
      trappedElement = containerEl
      containerEl.setAttribute('data-focus-trapped', 'true')
    }
    focusTrapStack++
    const handleTab = (e) => {
      if (e.key !== 'Tab') return
      const focusable = containerEl.querySelectorAll(
        'a[href], button:not([disabled]), textarea, input:not([disabled]), select, [tabindex]:not([tabindex="-1"])'
      )
      if (focusable.length === 0) { e.preventDefault(); return }
      const first = focusable[0]
      const last = focusable[focusable.length - 1]
      if (e.shiftKey) {
        if (document.activeElement === first) { e.preventDefault(); last.focus() }
      } else {
        if (document.activeElement === last) { e.preventDefault(); first.focus() }
      }
    }
    document.addEventListener('keydown', handleTab)
    return () => {
      document.removeEventListener('keydown', handleTab)
      focusTrapStack--
      if (focusTrapStack <= 0) {
        focusTrapStack = 0
        if (trappedElement) trappedElement.removeAttribute('data-focus-trapped')
        trappedElement = null
        if (previouslyFocused && typeof previouslyFocused.focus === 'function') {
          previouslyFocused.focus()
        }
        previouslyFocused = null
      }
    }
  }

  const announceToScreenReader = (message, priority = 'polite') => {
    const announcer = document.getElementById('sr-announcer') || (() => {
      const el = document.createElement('div')
      el.id = 'sr-announcer'
      el.setAttribute('aria-live', priority)
      el.setAttribute('aria-atomic', 'true')
      el.style.cssText = 'position:absolute;width:1px;height:1px;padding:0;margin:-1px;overflow:hidden;clip:rect(0,0,0,0);white-space:nowrap;border:0;'
      document.body.appendChild(el)
      return el
    })()
    announcer.textContent = message
    setTimeout(() => { announcer.textContent = '' }, 1000)
  }

  const getFocusableElements = (container) => {
    if (!container) return []
    return Array.from(container.querySelectorAll(
      'a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])'
    ))
  }

  return { trapFocus, announceToScreenReader, getFocusableElements }
}
