import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useReminderStore = defineStore('reminder', () => {
  const expiringBloggers = ref([])
  const hasShownLoginPopup = ref(false)
  const reminderRead = ref(false)

  const hasUnreadReminders = computed(() => {
    return expiringBloggers.value.length > 0 && !reminderRead.value
  })

  const hasExpiringBloggers = computed(() => {
    return expiringBloggers.value.length > 0
  })

  function setExpiringBloggers(bloggers) {
    expiringBloggers.value = bloggers
    if (bloggers.length === 0) {
      reminderRead.value = true
    }
  }

  function markAsRead() {
    reminderRead.value = true
    localStorage.setItem('reminderRead', 'true')
  }

  function markAllAsRead() {
    reminderRead.value = true
    hasShownLoginPopup.value = true
    localStorage.setItem('reminderRead', 'true')
    localStorage.setItem('hasShownLoginPopup', 'true')
  }

  function checkLoginPopupShown() {
    const saved = localStorage.getItem('hasShownLoginPopup')
    hasShownLoginPopup.value = saved === 'true'
  }

  function checkReminderRead() {
    const saved = localStorage.getItem('reminderRead')
    reminderRead.value = saved === 'true'
  }

  function resetForNewLogin() {
    hasShownLoginPopup.value = false
    localStorage.removeItem('hasShownLoginPopup')
  }

  function loadFromStorage() {
    checkLoginPopupShown()
    checkReminderRead()
  }

  return {
    expiringBloggers,
    hasShownLoginPopup,
    reminderRead,
    hasUnreadReminders,
    hasExpiringBloggers,
    setExpiringBloggers,
    markAsRead,
    markAllAsRead,
    checkLoginPopupShown,
    checkReminderRead,
    resetForNewLogin,
    loadFromStorage
  }
})
