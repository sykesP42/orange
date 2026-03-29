import { ref } from 'vue'

const notifications = ref([])
let notificationId = 0
const MAX_NOTIFICATIONS = 50

export const useNotification = () => {
  const addNotification = ({ message, type = 'info', duration = 3000 }) => {
    if (notifications.value.length >= MAX_NOTIFICATIONS) {
      notifications.value.shift()
    }

    const id = ++notificationId
    notifications.value.push({ id, message, type, timestamp: Date.now() })

    if (duration > 0) {
      setTimeout(() => {
        removeNotification(id)
      }, duration)
    }

    return id
  }

  const removeNotification = (id) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
  }

  const clearAll = () => {
    notifications.value = []
  }

  const success = (message, duration) => addNotification({ message, type: 'success', duration })
  const error = (message, duration) => addNotification({ message, type: 'error', duration })
  const warning = (message, duration) => addNotification({ message, type: 'warning', duration })
  const info = (message, duration) => addNotification({ message, type: 'info', duration })

  return {
    notifications,
    addNotification,
    removeNotification,
    clearAll,
    success,
    error,
    warning,
    info
  }
}
