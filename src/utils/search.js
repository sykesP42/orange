/**
 * 搜索防抖工具函数
 * @param {Function} fn - 需要防抖的函数
 * @param {number} delay - 延迟时间（毫秒）
 * @returns {Function} - 防抖后的函数
 */
export function useDebounce(fn, delay = 300) {
  let timer = null
  
  return function(...args) {
    if (timer) {
      clearTimeout(timer)
    }
    
    timer = setTimeout(() => {
      fn.apply(this, args)
      timer = null
    }, delay)
  }
}

/**
 * 搜索历史管理
 */
export class SearchHistory {
  constructor(key = 'search_history', max = 10) {
    this.key = key
    this.max = max
  }
  
  /**
   * 获取搜索历史
   */
  get() {
    try {
      const history = localStorage.getItem(this.key)
      return history ? JSON.parse(history) : []
    } catch (error) {
      console.error('读取搜索历史失败', error)
      return []
    }
  }
  
  /**
   * 添加搜索记录
   * @param {string} query - 搜索关键词
   */
  add(query) {
    if (!query || !query.trim()) return
    
    const trimmedQuery = query.trim()
    const history = this.get()
    
    // 移除已存在的记录
    const index = history.indexOf(trimmedQuery)
    if (index > -1) {
      history.splice(index, 1)
    }
    
    // 添加到开头
    history.unshift(trimmedQuery)
    
    // 限制数量
    if (history.length > this.max) {
      history.pop()
    }
    
    try {
      localStorage.setItem(this.key, JSON.stringify(history))
    } catch (error) {
      console.error('保存搜索历史失败', error)
    }
  }
  
  /**
   * 删除搜索记录
   * @param {string} query - 搜索关键词
   */
  remove(query) {
    const history = this.get()
    const index = history.indexOf(query)
    
    if (index > -1) {
      history.splice(index, 1)
      try {
        localStorage.setItem(this.key, JSON.stringify(history))
      } catch (error) {
        console.error('删除搜索记录失败', error)
      }
    }
  }
  
  /**
   * 清空搜索历史
   */
  clear() {
    try {
      localStorage.removeItem(this.key)
    } catch (error) {
      console.error('清空搜索历史失败', error)
    }
  }
}

/**
 * 热门搜索标签管理
 */
export class HotSearchTags {
  constructor(key = 'hot_search_tags', max = 8) {
    this.key = key
    this.max = max
  }
  
  /**
   * 获取热门标签
   */
  get() {
    try {
      const tags = localStorage.getItem(this.key)
      return tags ? JSON.parse(tags) : []
    } catch (error) {
      console.error('读取热门标签失败', error)
      return []
    }
  }
  
  /**
   * 更新热门标签（基于搜索频率）
   * @param {string} query - 搜索关键词
   */
  update(query) {
    if (!query || !query.trim()) return
    
    const trimmedQuery = query.trim()
    let tags = this.get()
    
    const existingTag = tags.find(tag => tag.text === trimmedQuery)
    
    if (existingTag) {
      // 增加搜索次数
      existingTag.count++
      existingTag.lastSearch = Date.now()
    } else {
      // 添加新标签
      tags.push({
        text: trimmedQuery,
        count: 1,
        lastSearch: Date.now()
      })
    }
    
    // 按搜索次数和时间排序
    tags.sort((a, b) => {
      if (b.count !== a.count) {
        return b.count - a.count
      }
      return b.lastSearch - a.lastSearch
    })
    
    // 限制数量
    if (tags.length > this.max) {
      tags = tags.slice(0, this.max)
    }
    
    try {
      localStorage.setItem(this.key, JSON.stringify(tags))
    } catch (error) {
      console.error('保存热门标签失败', error)
    }
  }
  
  /**
   * 获取热门标签列表（只返回文本）
   */
  getList() {
    return this.get().map(tag => tag.text)
  }
  
  /**
   * 清空热门标签
   */
  clear() {
    try {
      localStorage.removeItem(this.key)
    } catch (error) {
      console.error('清空热门标签失败', error)
    }
  }
}

/**
 * 搜索建议（本地缓存）
 */
export class SearchSuggestions {
  constructor() {
    this.cache = new Map()
    this.cacheExpire = 5 * 60 * 1000 // 5 分钟缓存
  }
  
  /**
   * 获取缓存的搜索建议
   * @param {string} query - 搜索关键词
   * @returns {Array|null}
   */
  get(query) {
    const cached = this.cache.get(query)
    
    if (!cached) return null
    
    if (Date.now() - cached.timestamp > this.cacheExpire) {
      this.cache.delete(query)
      return null
    }
    
    return cached.data
  }
  
  /**
   * 设置搜索建议缓存
   * @param {string} query - 搜索关键词
   * @param {Array} data - 建议数据
   */
  set(query, data) {
    this.cache.set(query, {
      data,
      timestamp: Date.now()
    })
  }
  
  /**
   * 清空缓存
   */
  clear() {
    this.cache.clear()
  }
}
