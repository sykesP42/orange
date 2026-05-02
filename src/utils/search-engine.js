import { getFullPinyin, getInitialPinyin, getPinyin, getPinyinInitial } from './pinyin-dict'

const MATCH_TYPE = {
  EXACT: 'exact',
  CONTAINS: 'contains',
  PINYIN_FULL: 'pinyin_full',
  PINYIN_INITIAL: 'pinyin_initial',
  SUBSEQUENCE: 'subsequence',
  NONE: 'none'
}

const MATCH_SCORES = {
  [MATCH_TYPE.EXACT]: 100,
  [MATCH_TYPE.CONTAINS]: 80,
  [MATCH_TYPE.PINYIN_FULL]: 60,
  [MATCH_TYPE.PINYIN_INITIAL]: 50,
  [MATCH_TYPE.SUBSEQUENCE]: 30,
  [MATCH_TYPE.NONE]: 0
}

class SearchCache {
  constructor(maxSize = 200) {
    this.cache = new Map()
    this.maxSize = maxSize
  }

  get(key) {
    const entry = this.cache.get(key)
    if (!entry) return null
    if (Date.now() - entry.timestamp > 5 * 60 * 1000) {
      this.cache.delete(key)
      return null
    }
    return entry.data
  }

  set(key, data) {
    if (this.cache.size >= this.maxSize) {
      const oldestKey = this.cache.keys().next().value
      this.cache.delete(oldestKey)
    }
    this.cache.set(key, { data, timestamp: Date.now() })
  }

  clear() {
    this.cache.clear()
  }
}

const globalCache = new SearchCache()

function detectSearchMode(keyword) {
  if (!keyword || !keyword.trim()) return 'normal'
  const hasChinese = /[\u4e00-\u9fa5]/.test(keyword)
  const hasLatin = /[a-zA-Z]/.test(keyword)
  if (hasChinese && hasLatin) return 'mixed'
  if (hasLatin) return 'pinyin'
  return 'chinese'
}

function matchItem(text, query) {
  if (!text || !query) return { type: MATCH_TYPE.NONE, score: 0, positions: [] }

  const lowerText = text.toLowerCase()
  const lowerQuery = query.toLowerCase()

  if (lowerText === lowerQuery) {
    return { type: MATCH_TYPE.EXACT, score: MATCH_SCORES[MATCH_TYPE.EXACT], positions: Array.from({ length: text.length }, (_, i) => i) }
  }

  const containsIdx = lowerText.indexOf(lowerQuery)
  if (containsIdx !== -1) {
    const positions = Array.from({ length: lowerQuery.length }, (_, i) => containsIdx + i)
    let score = MATCH_SCORES[MATCH_TYPE.CONTAINS]
    if (containsIdx === 0) score += 10
    return { type: MATCH_TYPE.CONTAINS, score, positions }
  }

  const fullPinyin = getFullPinyin(text)
  if (fullPinyin) {
    const pinyinIdx = fullPinyin.indexOf(lowerQuery)
    if (pinyinIdx !== -1) {
      const positions = mapPinyinPositionToOriginal(text, pinyinIdx, lowerQuery.length)
      let score = MATCH_SCORES[MATCH_TYPE.PINYIN_FULL]
      if (pinyinIdx === 0) score += 5
      return { type: MATCH_TYPE.PINYIN_FULL, score, positions }
    }
  }

  const initialPinyin = getInitialPinyin(text)
  if (initialPinyin) {
    const initialIdx = initialPinyin.indexOf(lowerQuery)
    if (initialIdx !== -1) {
      const positions = mapInitialPositionToOriginal(text, initialIdx, lowerQuery.length)
      let score = MATCH_SCORES[MATCH_TYPE.PINYIN_INITIAL]
      if (initialIdx === 0) score += 5
      if (lowerQuery.length <= 3) score += 5
      return { type: MATCH_TYPE.PINYIN_INITIAL, score, positions }
    }
  }

  if (lowerQuery.length >= 2) {
    const subseqResult = subsequenceMatch(lowerText, lowerQuery)
    if (subseqResult.matched) {
      return { type: MATCH_TYPE.SUBSEQUENCE, score: MATCH_SCORES[MATCH_TYPE.SUBSEQUENCE], positions: subseqResult.positions }
    }

    if (fullPinyin) {
      const pinyinSubseq = subsequenceMatch(fullPinyin, lowerQuery)
      if (pinyinSubseq.matched) {
        const positions = mapSubseqPositionsToOriginal(text, pinyinSubseq.positions, fullPinyin)
        return { type: MATCH_TYPE.SUBSEQUENCE, score: MATCH_SCORES[MATCH_TYPE.SUBSEQUENCE] - 5, positions }
      }
    }
  }

  return { type: MATCH_TYPE.NONE, score: 0, positions: [] }
}

function mapPinyinPositionToOriginal(text, pinyinStart, pinyinLength) {
  const positions = []
  let pinyinOffset = 0
  for (let i = 0; i < text.length; i++) {
    const charPinyin = getPinyin(text[i])
    const charPinyinLen = charPinyin ? charPinyin.length : 1
    const charStart = pinyinOffset
    const charEnd = pinyinOffset + charPinyinLen
    if (charStart < pinyinStart + pinyinLength && charEnd > pinyinStart) {
      positions.push(i)
    }
    pinyinOffset += charPinyinLen
  }
  return positions
}

function mapInitialPositionToOriginal(text, initialStart, initialLength) {
  const positions = []
  let initialOffset = 0
  for (let i = 0; i < text.length; i++) {
    const initial = getPinyinInitial(text[i])
    if (initial) {
      if (initialOffset >= initialStart && initialOffset < initialStart + initialLength) {
        positions.push(i)
      }
      initialOffset++
    } else if (/[a-zA-Z0-9]/.test(text[i])) {
      if (initialOffset >= initialStart && initialOffset < initialStart + initialLength) {
        positions.push(i)
      }
      initialOffset++
    }
  }
  return positions
}

function mapSubseqPositionsToOriginal(text, pinyinPositions, fullPinyin) {
  const positions = new Set()
  let pinyinOffset = 0
  for (let i = 0; i < text.length; i++) {
    const charPinyin = getPinyin(text[i])
    const charPinyinLen = charPinyin ? charPinyin.length : 1
    for (let p = pinyinOffset; p < pinyinOffset + charPinyinLen; p++) {
      if (pinyinPositions.includes(p)) {
        positions.add(i)
        break
      }
    }
    pinyinOffset += charPinyinLen
  }
  return Array.from(positions)
}

function subsequenceMatch(text, query) {
  let textIdx = 0
  let queryIdx = 0
  const positions = []

  while (textIdx < text.length && queryIdx < query.length) {
    if (text[textIdx] === query[queryIdx]) {
      positions.push(textIdx)
      queryIdx++
    }
    textIdx++
  }

  return {
    matched: queryIdx === query.length,
    positions
  }
}

function searchList(items, query, options = {}) {
  const {
    fields = ['name', 'nickname', 'title'],
    fieldWeights = {},
    keywordsField = 'keywords',
    limit = 50,
    minScore = 1,
    useCache = true
  } = options

  if (!query || !query.trim()) {
    return items.slice(0, limit).map(item => ({ item, score: 0, matchType: MATCH_TYPE.NONE, matchPositions: {} }))
  }

  const cacheKey = `${query}::${items.length}::${fields.join(',')}::${JSON.stringify(fieldWeights)}`
  if (useCache) {
    const cached = globalCache.get(cacheKey)
    if (cached) return cached.slice(0, limit)
  }

  const trimmedQuery = query.trim()
  const results = []

  for (const item of items) {
    let bestScore = 0
    let bestMatchType = MATCH_TYPE.NONE
    const matchPositions = {}

    for (const field of fields) {
      const value = item[field]
      if (!value) continue

      const weight = fieldWeights[field] || 1
      const match = matchItem(String(value), trimmedQuery)

      const weightedScore = match.score * weight
      if (weightedScore > bestScore) {
        bestScore = weightedScore
        bestMatchType = match.type
      }

      if (match.positions.length > 0) {
        matchPositions[field] = match.positions
      }
    }

    if (keywordsField && item[keywordsField] && Array.isArray(item[keywordsField])) {
      for (const kw of item[keywordsField]) {
        const match = matchItem(String(kw), trimmedQuery)
        if (match.score > bestScore) {
          bestScore = match.score * 0.8
          bestMatchType = match.type
        }
      }
    }

    if (bestScore >= minScore) {
      results.push({
        item,
        score: bestScore,
        matchType: bestMatchType,
        matchPositions
      })
    }
  }

  results.sort((a, b) => {
    if (b.score !== a.score) return b.score - a.score
    const aName = a.item.name || a.item.nickname || a.item.title || ''
    const bName = b.item.name || b.item.nickname || b.item.title || ''
    return aName.localeCompare(bName, 'zh-CN')
  })

  const limited = results.slice(0, limit)

  if (useCache) {
    globalCache.set(cacheKey, limited)
  }

  return limited
}

function highlightText(text, positions, highlightTag = 'mark') {
  if (!text || !positions || positions.length === 0) return escapeHtml(text || '')

  const safeText = escapeHtml(text)
  const sortedPositions = [...new Set(positions)].sort((a, b) => a - b)

  if (sortedPositions.length === 0) return safeText

  const ranges = []
  let rangeStart = sortedPositions[0]
  let rangeEnd = sortedPositions[0]

  for (let i = 1; i < sortedPositions.length; i++) {
    if (sortedPositions[i] === rangeEnd + 1) {
      rangeEnd = sortedPositions[i]
    } else {
      ranges.push([rangeStart, rangeEnd])
      rangeStart = sortedPositions[i]
      rangeEnd = sortedPositions[i]
    }
  }
  ranges.push([rangeStart, rangeEnd])

  let result = ''
  let lastIdx = 0
  for (const [start, end] of ranges) {
    result += safeText.substring(lastIdx, start)
    result += `<${highlightTag} class="highlight">${safeText.substring(start, end + 1)}</${highlightTag}>`
    lastIdx = end + 1
  }
  result += safeText.substring(lastIdx)

  return result
}

function escapeHtml(str) {
  if (!str) return ''
  return String(str).replace(/[&<>"'/]/g, c => ({
    '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;',
    "'": '&#39;', '/': '&#x2F;'
  }[c]))
}

function fuzzyMatch(text, query) {
  if (!query) return true
  const result = matchItem(text, query)
  return result.type !== MATCH_TYPE.NONE
}

function getMatchInfo(text, query) {
  return matchItem(text, query)
}

export {
  searchList,
  highlightText,
  fuzzyMatch,
  getMatchInfo,
  detectSearchMode,
  MATCH_TYPE,
  MATCH_SCORES,
  SearchCache,
  globalCache
}

export default {
  searchList,
  highlightText,
  fuzzyMatch,
  getMatchInfo,
  detectSearchMode,
  MATCH_TYPE,
  MATCH_SCORES,
  SearchCache
}
