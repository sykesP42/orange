import { ref, readonly } from 'vue'

const isVisible = ref(false)
const title = ref('确认操作')
const message = ref('确定要执行此操作吗？')
const dialogType = ref('info')
const confirmText = ref('确定')
const cancelText = ref('取消')
const resolver = ref(null)

export const useConfirm = () => {
  const confirm = (options = {}) => {
    return new Promise((resolve) => {
      isVisible.value = true
      title.value = options.title || '确认操作'
      message.value = options.message || '确定要执行此操作吗？'
      dialogType.value = options.type || 'info'
      confirmText.value = options.confirmText || '确定'
      cancelText.value = options.cancelText || '取消'
      resolver.value = resolve
    })
  }

  const confirmDanger = (messageText, options = {}) => {
    return confirm({
      ...options,
      title: options.title || '确认删除',
      message: messageText,
      type: 'danger',
      confirmText: options.confirmText || '删除'
    })
  }

  const confirmWarning = (messageText, options = {}) => {
    return confirm({
      ...options,
      title: options.title || '操作警告',
      message: messageText,
      type: 'warning'
    })
  }

  const handleConfirm = () => {
    if (resolver.value) {
      resolver.value(true)
    }
    isVisible.value = false
    resolver.value = null
  }

  const handleCancel = () => {
    if (resolver.value) {
      resolver.value(false)
    }
    isVisible.value = false
    resolver.value = null
  }

  return {
    isVisible: readonly(isVisible),
    title: readonly(title),
    message: readonly(message),
    dialogType: readonly(dialogType),
    confirmText: readonly(confirmText),
    cancelText: readonly(cancelText),
    confirm,
    confirmDanger,
    confirmWarning,
    handleConfirm,
    handleCancel
  }
}

export default useConfirm
