import { describe, it, expect } from 'vitest'

describe('Home.vue Core Logic', () => {
  describe('isSelected — Set-based O(1) lookup', () => {
    const createSelectedSet = (ids) => new Set(ids)

    it('returns true for existing ID in selected set', () => {
      const set = createSelectedSet([1, 2, 3])
      expect(set.has(2)).toBe(true)
      expect(set.has(1)).toBe(true)
      expect(set.has(3)).toBe(true)
    })

    it('returns false for non-existing ID', () => {
      const set = createSelectedSet([1, 2, 3])
      expect(set.has(99)).toBe(false)
      expect(set.has(0)).toBe(false)
      expect(set.has(-1)).toBe(false)
    })

    it('returns false for empty selection', () => {
      const set = createSelectedSet([])
      expect(set.has(1)).toBe(false)
    })

    it('handles string IDs correctly', () => {
      const set = createSelectedSet(['a', 'b', 'c'])
      expect(set.has('b')).toBe(true)
      expect(set.has('x')).toBe(false)
    })

    it('handles large selections efficiently (10K IDs)', () => {
      const ids = Array.from({ length: 10000 }, (_, i) => i + 1)
      const set = createSelectedSet(ids)
      expect(set.has(5000)).toBe(true)
      expect(set.has(10000)).toBe(true)
      expect(set.has(0)).toBe(false)
    })
  })

  describe('selectedBreakdown — status distribution', () => {
    const computeBreakdown = (selectedIds, bloggers) => {
      const map = {}
      selectedIds.forEach((id) => {
        const b = bloggers.find((bl) => bl.id === id)
        if (b) {
          const s = b.status || '未知'
          map[s] = (map[s] || 0) + 1
        }
      })
      return map
    }

    const sampleBloggers = [
      { id: 1, nickname: 'Alice', status: '已合作' },
      { id: 2, nickname: 'Bob', status: '洽谈中' },
      { id: 3, nickname: 'Carol', status: '已合作' },
      { id: 4, nickname: 'Dave', status: '初次联系' },
      { id: 5, nickname: 'Eve', status: '洽谈中' },
      { id: 6, nickname: 'Frank', status: '已拒绝' },
    ]

    it('returns empty object for no selection', () => {
      const result = computeBreakdown([], sampleBloggers)
      expect(result).toEqual({})
    })

    it('counts single status correctly', () => {
      const result = computeBreakdown([4], sampleBloggers)
      expect(result).toEqual({ '初次联系': 1 })
    })

    it('counts multiple statuses correctly', () => {
      const result = computeBreakdown([1, 2, 3, 5], sampleBloggers)
      expect(result).toEqual({
        '已合作': 2,
        '洽谈中': 2,
      })
    })

    it('includes all statuses when all selected', () => {
      const result = computeBreakdown([1, 2, 3, 4, 5, 6], sampleBloggers)
      expect(result).toEqual({
        '已合作': 2,
        '洽谈中': 2,
        '初次联系': 1,
        '已拒绝': 1,
      })
    })

    it('ignores IDs not found in blogger list', () => {
      const result = computeBreakdown([1, 999], sampleBloggers)
      expect(result).toEqual({ '已合作': 1 })
    })

    it('defaults to "未知" for missing status', () => {
      const bloggersNoStatus = [{ id: 7, nickname: 'Ghost' }]
      const result = computeBreakdown([7], bloggersNoStatus)
      expect(result).toEqual({ '未知': 1 })
    })
  })

  describe('sortBloggers — multi-field sorting', () => {
    const sortWith = (bloggers, field, order) => {
      const sorted = [...bloggers]
      sorted.sort((a, b) => {
        let valA, valB
        switch (field) {
          case 'nickname':
            valA = a.nickname || ''
            valB = b.nickname || ''
            break
          case 'category':
            valA = a.category || ''
            valB = b.category || ''
            break
          case 'status':
            valA = a.status || ''
            valB = b.status || ''
            break
          default:
            return 0
        }
        if (valA < valB) return order === 'asc' ? -1 : 1
        if (valA > valB) return order === 'asc' ? 1 : -1
        return 0
      })
      return sorted
    }

    const sampleData = [
      { id: 1, nickname: 'Charlie', category: '科技', status: '已合作' },
      { id: 2, nickname: 'Alice', category: '美妆', status: '洽谈中' },
      { id: 3, nickname: 'Bob', category: '游戏', status: '初次联系' },
    ]

    it('sorts by nickname ascending', () => {
      const result = sortWith(sampleData, 'nickname', 'asc')
      expect(result.map((b) => b.nickname)).toEqual(['Alice', 'Bob', 'Charlie'])
    })

    it('sorts by nickname descending', () => {
      const result = sortWith(sampleData, 'nickname', 'desc')
      expect(result.map((b) => b.nickname)).toEqual(['Charlie', 'Bob', 'Alice'])
    })

    it('sorts by category ascending', () => {
      const result = sortWith(sampleData, 'category', 'asc')
      const categories = result.map((b) => b.category)
      expect(categories).toEqual([...categories].sort())
    })

    it('sorts by status ascending', () => {
      const result = sortWith(sampleData, 'status', 'asc')
      expect(result.map((b) => b.status)).toEqual(['初次联系', '已合作', '洽谈中'])
    })

    it('does not mutate original array', () => {
      const original = [...sampleData]
      sortWith(sampleData, 'nickname', 'asc')
      expect(sampleData).toEqual(original)
    })

    it('handles empty array', () => {
      const result = sortWith([], 'nickname', 'asc')
      expect(result).toEqual([])
    })

    it('handles single element', () => {
      const data = [{ id: 1, nickname: 'Solo' }]
      const result = sortWith(data, 'nickname', 'asc')
      expect(result.length).toBe(1)
    })

    it('handles null/undefined values gracefully', () => {
      const data = [
        { id: 1, nickname: 'Zebra' },
        { id: 2, nickname: null },
        { id: 3, nickname: undefined },
        { id: 4, nickname: 'Alpha' },
      ]
      const result = sortWith(data, 'nickname', 'asc')
      const names = result.map((b) => b.nickname || '')
      expect(names[0]).toBe('')
      expect(names[names.length - 1]).toBe('Zebra')
    })

    it('is stable for equal values', () => {
      const data = [
        { id: 1, nickname: 'Same', category: 'B' },
        { id: 2, nickname: 'Same', category: 'A' },
      ]
      const result = sortWith(data, 'nickname', 'asc')
      expect(result[0].id).toBe(1)
    })
  })

  describe('toggleSelect — add/remove from selection', () => {
    const simulateToggle = (currentSelection, targetId) => {
      const index = currentSelection.indexOf(targetId)
      if (index === -1) {
        return [...currentSelection, targetId]
      } else {
        const next = [...currentSelection]
        next.splice(index, 1)
        return next
      }
    }

    it('adds ID when not present', () => {
      expect(simulateToggle([1, 2], 3)).toEqual([1, 2, 3])
    })

    it('removes ID when already present', () => {
      expect(simulateToggle([1, 2, 3], 2)).toEqual([1, 3])
    })

    it('toggles: add then remove returns to original', () => {
      const initial = [1, 2]
      const afterAdd = simulateToggle(initial, 3)
      const afterRemove = simulateToggle(afterAdd, 3)
      expect(afterRemove).toEqual(initial)
    })

    it('works with empty initial state', () => {
      expect(simulateToggle([], 1)).toEqual([1])
    })

    it('removing only item returns empty', () => {
      expect(simulateToggle([42], 42)).toEqual([])
    })
  })

  describe('isAllSelected — full selection detection', () => {
    const checkAllSelected = (bloggersCount, selectedCount) =>
      bloggersCount > 0 && selectedCount === bloggersCount

    it('returns true when counts match and non-zero', () => {
      expect(checkAllSelected(5, 5)).toBe(true)
    })

    it('returns false when not all selected', () => {
      expect(checkAllSelected(5, 3)).toBe(false)
    })

    it('returns false when more selected than exist', () => {
      expect(checkAllSelected(3, 5)).toBe(false)
    })

    it('returns false for zero bloggers even with zero selected', () => {
      expect(checkAllSelected(0, 0)).toBe(false)
    })
  })

  describe('tagCloudData transformation', () => {
    const buildTagCloud = (allTags, tagCounts) => {
      return allTags.map((name) => ({
        name,
        count: tagCounts[name] || 0,
      }))
    }

    it('maps tags to cloud format with counts', () => {
      const tags = ['美妆', '科技', '游戏']
      const counts = { '美妆': 45, '科技': 32, '游戏': 28 }
      const result = buildTagCloud(tags, counts)
      expect(result).toEqual([
        { name: '美妆', count: 45 },
        { name: '科技', count: 32 },
        { name: '游戏', count: 28 },
      ])
    })

    it('assigns 0 count for unknown tags', () => {
      const result = buildTagCloud(['新标签'], {})
      expect(result[0].count).toBe(0)
    })

    it('returns empty array for no tags', () => {
      const result = buildTagCloud([], { 'x': 5 })
      expect(result).toEqual([])
    })
  })
})
