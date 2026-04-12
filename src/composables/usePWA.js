import { ref, onMounted, onUnmounted } from 'vue'

const isOnline = ref(navigator.onLine)
const showUpdateBanner = ref(false)
const isInstallable = ref(false)
const deferredPrompt = ref(null)
const swRegistration = ref(null)

let onlineHandler = null
let offlineHandler = null

export function usePWA() {
  const updateSW = async () => {
    if (!swRegistration.value || !swRegistration.value.waiting) return

    swRegistration.value.waiting.postMessage({ type: 'SKIP_WAITING' })

    window.addEventListener('controllerchange', () => {
      window.location.reload()
    }, { once: true })
  }

  const clearCache = async () => {
    if (!swRegistration.value) return false

    try {
      return new Promise((resolve) => {
        const handler = (event) => {
          if (event.data?.type === 'CACHES_CLEARED') {
            navigator.serviceWorker.removeEventListener('message', handler)
            resolve(true)
          }
        }
        navigator.serviceWorker.addEventListener('message', handler)
        swRegistration.value.active.postMessage({ type: 'CLEAR_ALL_CACHES' })
        setTimeout(() => resolve(false), 5000)
      })
    } catch (e) {
      return false
    }
  }

  const installPWA = async () => {
    if (!deferredPrompt.value) return false

    deferredPrompt.value.prompt()
    const result = await deferredPrompt.value.userChoice
    deferredPrompt.value = null

    if (result.outcome === 'accepted') {
      isInstallable.value = false
      return true
    }
    return false
  }

  onMounted(() => {
    onlineHandler = () => { isOnline.value = true }
    offlineHandler = () => { isOnline.value = false }

    window.addEventListener('online', onlineHandler)
    window.addEventListener('offline', offlineHandler)

    window.addEventListener('sw-update-available', () => {
      showUpdateBanner.value = true
    })

    window.addEventListener('beforeinstallprompt', (e) => {
      e.preventDefault()
      deferredPrompt.value = e
      isInstallable.value = true
    })

    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.ready.then((reg) => {
        swRegistration.value = reg
      }).catch(() => {})
    }
  })

  onUnmounted(() => {
    if (onlineHandler) window.removeEventListener('online', onlineHandler)
    if (offlineHandler) window.removeEventListener('offline', offlineHandler)
  })

  return {
    isOnline,
    showUpdateBanner,
    isInstallable,
    updateSW,
    clearCache,
    installPWA,
  }
}
