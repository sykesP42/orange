import { describe, it, expect, beforeEach } from 'vitest'
import { searchList, highlightText, fuzzyMatch, getMatchInfo, detectSearchMode, MATCH_TYPE, SearchCache } from '../../utils/search-engine'

describe('search-engine', () => {
  describe('detectSearchMode', () => {
    it('returns normal for empty input', () => {
      expect(detectSearchMode('')).toBe('normal')
      expect(detectSearchMode(null)).toBe('normal')
    })

    it('returns chinese for Chinese characters', () => {
      expect(detectSearchMode('博主')).toBe('chinese')
    })

    it('returns pinyin for Latin characters', () => {
      expect(detectSearchMode('bo')).toBe('pinyin')
      expect(detectSearchMode('zfb')).toBe('pinyin')
    })

    it('returns mixed for Chinese and Latin mix', () => {
      expect(detectSearchMode('博bo')).toBe('mixed')
    })
  })

  describe('fuzzyMatch', () => {
    it('returns true for empty query', () => {
      expect(fuzzyMatch('hello', '')).toBe(true)
    })

    it('returns true for exact substring match', () => {
      expect(fuzzyMatch('录入博主', '录入')).toBe(true)
      expect(fuzzyMatch('hello world', 'world')).toBe(true)
    })

    it('returns true for pinyin full match', () => {
      expect(fuzzyMatch('博主', 'bo')).toBe(true)
      expect(fuzzyMatch('录入', 'lu')).toBe(true)
    })

    it('returns true for pinyin initial match', () => {
      expect(fuzzyMatch('支付宝', 'zfb')).toBe(true)
      expect(fuzzyMatch('录入博主', 'lrbz')).toBe(true)
    })

    it('returns false for no match', () => {
      expect(fuzzyMatch('hello', 'xyz')).toBe(false)
    })

    it('is case insensitive', () => {
      expect(fuzzyMatch('Hello World', 'hello')).toBe(true)
      expect(fuzzyMatch('Hello World', 'HELLO')).toBe(true)
    })
  })

  describe('getMatchInfo', () => {
    it('returns EXACT for exact match', () => {
      const result = getMatchInfo('hello', 'hello')
      expect(result.type).toBe(MATCH_TYPE.EXACT)
      expect(result.score).toBeGreaterThan(0)
    })

    it('returns CONTAINS for substring match', () => {
      const result = getMatchInfo('hello world', 'world')
      expect(result.type).toBe(MATCH_TYPE.CONTAINS)
    })

    it('returns PINYIN_FULL for pinyin full match', () => {
      const result = getMatchInfo('博主', 'bo')
      expect(result.type).toBe(MATCH_TYPE.PINYIN_FULL)
    })

    it('returns PINYIN_INITIAL for pinyin initial match', () => {
      const result = getMatchInfo('支付宝', 'zfb')
      expect(result.type).toBe(MATCH_TYPE.PINYIN_INITIAL)
    })

    it('returns NONE for no match', () => {
      const result = getMatchInfo('hello', 'xyz')
      expect(result.type).toBe(MATCH_TYPE.NONE)
      expect(result.score).toBe(0)
    })

    it('returns positions for CONTAINS match', () => {
      const result = getMatchInfo('hello world', 'world')
      expect(result.positions).toContain(6)
      expect(result.positions.length).toBe(5)
    })

    it('returns positions for PINYIN_INITIAL match', () => {
      const result = getMatchInfo('支付宝', 'zfb')
      expect(result.positions.length).toBe(3)
    })
  })

  describe('searchList', () => {
    const items = [
      { id: 1, name: '支付宝', category: '支付', platform: '阿里' },
      { id: 2, name: '微信支付', category: '支付', platform: '腾讯' },
      { id: 3, name: '抖音', category: '短视频', platform: '字节跳动' },
      { id: 4, name: '淘宝', category: '电商', platform: '阿里' },
      { id: 5, name: '美团', category: '本地生活', platform: '美团' },
    ]

    it('returns all items when query is empty', () => {
      const results = searchList(items, '', { fields: ['name'] })
      expect(results.length).toBe(5)
    })

    it('filters by Chinese substring', () => {
      const results = searchList(items, '支付', { fields: ['name'] })
      expect(results.length).toBe(2)
      expect(results.map(r => r.item.name)).toContain('支付宝')
      expect(results.map(r => r.item.name)).toContain('微信支付')
    })

    it('filters by pinyin full', () => {
      const results = searchList(items, 'zhifu', { fields: ['name'] })
      expect(results.length).toBeGreaterThanOrEqual(1)
    })

    it('filters by pinyin initial', () => {
      const results = searchList(items, 'zfb', { fields: ['name'] })
      expect(results.length).toBeGreaterThanOrEqual(1)
      expect(results[0].item.name).toBe('支付宝')
    })

    it('sorts by relevance score', () => {
      const results = searchList(items, '支付', { fields: ['name'] })
      expect(results[0].score).toBeGreaterThanOrEqual(results[1].score)
    })

    it('respects field weights', () => {
      const results1 = searchList(items, '支付', { fields: ['name', 'category'], fieldWeights: { name: 2, category: 0.5 } })
      expect(results1.length).toBeGreaterThan(0)
    })

    it('respects limit', () => {
      const results = searchList(items, '', { fields: ['name'], limit: 2 })
      expect(results.length).toBe(2)
    })

    it('supports keywords field', () => {
      const itemsWithKeywords = [
        { name: '测试', keywords: ['test', 'ceshi'] },
        { name: '其他', keywords: ['other'] },
      ]
      const results = searchList(itemsWithKeywords, 'test', { fields: ['name'], keywordsField: 'keywords' })
      expect(results.length).toBe(1)
      expect(results[0].item.name).toBe('测试')
    })
  })

  describe('highlightText', () => {
    it('highlights specified positions', () => {
      const result = highlightText('hello world', [0, 1, 2, 3, 4])
      expect(result).toContain('<mark')
      expect(result).toContain('hello')
    })

    it('escapes HTML in text', () => {
      const result = highlightText('<script>alert(1)</script>', [])
      expect(result).not.toContain('<script>')
      expect(result).toContain('&lt;script')
    })

    it('escapes HTML and highlights safely', () => {
      const result = highlightText('<script>alert(1)</script>', [0])
      expect(result).not.toContain('<script>')
      expect(result).toContain('<mark')
    })

    it('returns escaped text when no positions', () => {
      const result = highlightText('hello', [])
      expect(result).toBe('hello')
    })

    it('handles adjacent positions as continuous highlight', () => {
      const result = highlightText('hello', [0, 1, 2])
      expect(result).toContain('<mark class="highlight">hel</mark>')
    })
  })

  describe('SearchCache', () => {
    let cache

    beforeEach(() => {
      cache = new SearchCache(3)
    })

    it('stores and retrieves values', () => {
      cache.set('key1', 'value1')
      expect(cache.get('key1')).toBe('value1')
    })

    it('returns null for missing keys', () => {
      expect(cache.get('missing')).toBeNull()
    })

    it('evicts oldest entry when max size reached', () => {
      cache.set('key1', 'value1')
      cache.set('key2', 'value2')
      cache.set('key3', 'value3')
      cache.set('key4', 'value4')
      expect(cache.get('key1')).toBeNull()
      expect(cache.get('key4')).toBe('value4')
    })

    it('clears all entries', () => {
      cache.set('key1', 'value1')
      cache.clear()
      expect(cache.get('key1')).toBeNull()
    })
  })

  describe('performance', () => {
    it('searches 1000 items within 300ms', () => {
      const items = Array.from({ length: 1000 }, (_, i) => ({
        id: i,
        name: `博主${i}号`,
        category: ['科技', '财经', '娱乐', '体育', '教育'][i % 5],
        platform: ['微信', '微博', '抖音', 'B站', '小红书'][i % 5],
      }))

      const start = performance.now()
      const results = searchList(items, 'bozhu', { fields: ['name', 'category', 'platform'] })
      const elapsed = performance.now() - start

      expect(elapsed).toBeLessThan(300)
      expect(results.length).toBeGreaterThan(0)
    })

    it('searches with pinyin initials within 300ms', () => {
      const items = Array.from({ length: 1000 }, (_, i) => ({
        id: i,
        name: `测试博主${i}`,
        category: '科技',
      }))

      const start = performance.now()
      const results = searchList(items, 'csbz', { fields: ['name'] })
      const elapsed = performance.now() - start

      expect(elapsed).toBeLessThan(300)
    })
  })
})
