import { ref, onMounted, onUnmounted } from 'vue'

export function useOnlineStatus() {
  const isOnline = ref(navigator.onLine)
  const wasOffline = ref(false)
  const lastOfflineTime = ref(null)

  const handleOnline = () => {
    isOnline.value = true
    if (wasOffline.value) {
      lastOfflineTime.value = Date.now()
    }
  }

  const handleOffline = () => {
    isOnline.value = false
    wasOffline.value = true
  }

  onMounted(() => {
    window.addEventListener('online', handleOnline)
    window.addEventListener('offline', handleOffline)
  })

  onUnmounted(() => {
    window.removeEventListener('online', handleOnline)
    window.removeEventListener('offline', handleOffline)
  })

  return { isOnline, wasOffline, lastOfflineTime }
}
