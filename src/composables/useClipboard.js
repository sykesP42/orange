import { ref } from 'vue'

export function useClipboard() {
  const copied = ref(false)
  const copyTimer = ref(null)

  const copy = async (text) => {
    try {
      await navigator.clipboard.writeText(text)
      copied.value = true
      clearTimeout(copyTimer.value)
      copyTimer.value = setTimeout(() => {
        copied.value = false
      }, 2000)
      return true
    } catch {
      const textarea = document.createElement('textarea')
      textarea.value = text
      textarea.style.cssText = 'position:fixed;opacity:0;left:-9999px'
      document.body.appendChild(textarea)
      textarea.select()
      try {
        document.execCommand('copy')
        copied.value = true
        clearTimeout(copyTimer.value)
        copyTimer.value = setTimeout(() => {
          copied.value = false
        }, 2000)
        return true
      } catch {
        return false
      } finally {
        document.body.removeChild(textarea)
      }
    }
  }

  return { copied, copy }
}
