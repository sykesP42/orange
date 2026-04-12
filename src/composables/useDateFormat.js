export function useDateFormat() {
  const formatRelative = (time) => {
    if (!time) return ''
    const date = new Date(time)
    const now = new Date()
    const diff = now - date

    if (diff < 60000) return '刚刚'
    if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
    if (diff < 172800000) return '昨天'
    if (diff < 604800000) return `${Math.floor(diff / 86400000)}天前`
    if (diff < 2592000000) return `${Math.floor(diff / 604800000)}周前`
    return formatDate(time)
  }

  const formatDate = (time) => {
    if (!time) return '-'
    return new Date(time).toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    })
  }

  const formatDateTime = (time) => {
    if (!time) return '-'
    return new Date(time).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const formatTime = (time) => {
    if (!time) return ''
    return new Date(time).toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const formatDuration = (ms) => {
    if (ms < 1000) return `${ms}ms`
    if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
    if (ms < 3600000) return `${Math.floor(ms / 60000)}m ${Math.floor((ms % 60000) / 1000)}s`
    return `${Math.floor(ms / 3600000)}h ${Math.floor((ms % 3600000) / 60000)}m`
  }

  const getDaysAgo = (time) => {
    if (!time) return -1
    const diff = Date.now() - new Date(time).getTime()
    return Math.floor(diff / 86400000)
  }

  return {
    formatRelative,
    formatDate,
    formatDateTime,
    formatTime,
    formatDuration,
    getDaysAgo
  }
}
