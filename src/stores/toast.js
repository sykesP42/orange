import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useToastStore = defineStore('toast', () => {
  const toasts = ref([])
  let toastId = 0

  const addToast = (options) => {
    const id = ++toastId
    const toast = {
      id,
      type: options.type || 'info',
      title: options.title,
      message: options.message,
      duration: options.duration || 3000
    }
    toasts.value.push(toast)
    
    if (toast.duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, toast.duration)
    }
    
    return id
  }

  const removeToast = (id) => {
    const index = toasts.value.findIndex(t => t.id === id)
    if (index > -1) {
      toasts.value.splice(index, 1)
    }
  }

  const success = (title, message, duration) => addToast({ type: 'success', title, message, duration })
  const error = (title, message, duration) => addToast({ type: 'error', title, message, duration })
  const warning = (title, message, duration) => addToast({ type: 'warning', title, message, duration })
  const info = (title, message, duration) => addToast({ type: 'info', title, message, duration })

  return {
    toasts,
    addToast,
    removeToast,
    success,
    error,
    warning,
    info
  }
})
