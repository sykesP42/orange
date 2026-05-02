import { getFullPinyin, getInitialPinyin, getPinyin, getPinyinInitial } from './pinyin-dict'
import { getMatchInfo, highlightText as engineHighlight, detectSearchMode, MATCH_TYPE } from './search-engine'

export function isPinyinInput(text) {
  if (!text) return false
  for (const char of text) {
    if (/[\u4e00-\u9fa5]/.test(char)) {
      return false
    }
  }
  return true
}

export function getSearchMode(keyword) {
  if (!keyword || keyword.trim() === '') {
    return { mode: 'normal', label: '' }
  }

  const mode = detectSearchMode(keyword)

  switch (mode) {
    case 'pinyin':
      return {
        mode: 'pinyin',
        label: '🔤 拼音搜索',
        description: '正在使用拼音首字母或全拼搜索'
      }
    case 'mixed':
      return {
        mode: 'mixed',
        label: '🔤 中英混合搜索',
        description: '支持中文和拼音同时匹配'
      }
    case 'chinese':
      return {
        mode: 'chinese',
        label: '🔍 中文搜索',
        description: '正在搜索中文内容'
      }
    default:
      return { mode: 'normal', label: '' }
  }
}

export function highlightMatch(text, keyword) {
  if (!text || !keyword) return escapeHtml(text || '')

  const matchInfo = getMatchInfo(text, keyword)

  if (matchInfo.type === MATCH_TYPE.NONE) {
    return escapeHtml(text)
  }

  return engineHighlight(text, matchInfo.positions, 'mark')
}

export function getSearchHint(keyword) {
  if (!keyword) {
    return '支持中文、拼音首字母、完整拼音搜索'
  }

  const mode = getSearchMode(keyword)

  switch (mode.mode) {
    case 'pinyin':
      if (keyword.length <= 3) {
        return `输入 "${keyword}" 将匹配拼音首字母包含 "${keyword}" 的内容`
      } else {
        return `正在搜索拼音包含 "${keyword}" 的内容`
      }
    case 'mixed':
      return '同时匹配中文内容和拼音'
    default:
      return ''
  }
}

export function useEnhancedDebounce(fn, delay = 300, onImmediate = null) {
  let timer = null

  const debounced = function (...args) {
    if (timer) clearTimeout(timer)

    if (onImmediate && !timer) {
      onImmediate.apply(this, args)
    }

    timer = setTimeout(() => {
      fn.apply(this, args)
      timer = null
    }, delay)
  }

  debounced.cancel = () => {
    if (timer) {
      clearTimeout(timer)
      timer = null
    }
  }

  debounced.flush = () => {
    if (timer) {
      clearTimeout(timer)
      fn()
      timer = null
    }
  }

  return debounced
}

export function getPinyinForText(text) {
  return getFullPinyin(text)
}

export function getInitialForText(text) {
  return getInitialPinyin(text)
}

function escapeHtml(str) {
  if (!str) return ''
  return String(str).replace(/[&<>"'/]/g, c => ({
    '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;',
    "'": '&#39;', '/': '&#x2F;'
  }[c]))
}
