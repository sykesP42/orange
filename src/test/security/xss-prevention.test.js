import { describe, it, expect, beforeEach } from 'vitest'
import { highlightMatch } from '../../utils/pinyin-search'

describe('Security: highlightMatch XSS Prevention', () => {
  it('should escape HTML in text parameter', () => {
    const result = highlightMatch('<script>alert("xss")</script>', 'script')
    expect(result).not.toContain('<script')
    expect(result).toContain('&lt;')
    expect(result).toContain('&gt;')
    expect(result).toContain('<mark')
  })

  it('should escape HTML in keyword parameter when not matching', () => {
    const result = highlightMatch('hello world', '<img onerror=alert(1)>')
    expect(result).toBe('hello world')
    expect(result).not.toContain('<img')
    expect(result).not.toContain('onerror')
  })

  it('should handle combined XSS attack vectors', () => {
    const maliciousText = '<div onclick="steal()"><script>evil()</script></div>'
    const maliciousKeyword = '<svg onload="hack()">'
    const result = highlightMatch(maliciousText, maliciousKeyword)
    expect(result).not.toContain('onclick="')
    expect(result).not.toContain('onload="')
    expect(result).not.toContain('<script>')
    expect(result).not.toContain('<svg')
    expect(result).toContain('&lt;')
    expect(result).toContain('&gt;')
  })

  it('should safely render mark tags only', () => {
    const result = highlightMatch('hello world', 'world')
    expect(result).toBe('hello <mark class="highlight">world</mark>')
    const markCount = (result.match(/<mark/g) || []).length
    expect(markCount).toBe(1)
  })

  it('should handle null/undefined gracefully', () => {
    expect(highlightMatch(null, 'test')).toBe('')
    expect(highlightMatch(undefined, 'test')).toBe('')
    expect(highlightMatch('test', null)).toBe('test')
    expect(highlightMatch('', '')).toBe('')
  })

  it('should handle empty string inputs', () => {
    expect(highlightMatch('', 'test')).toBe('')
    expect(highlightMatch('test', '')).toBe('test')
  })

  it('should handle regex special characters in keyword', () => {
    const result = highlightMatch('price: $100 (USD)', '$100')
    expect(result).toContain('<mark class="highlight">$100</mark>')
  })

  it('should handle very long strings without performance issues', () => {
    const longText = 'a'.repeat(10000) + 'target' + 'b'.repeat(10000)
    const start = Date.now()
    const result = highlightMatch(longText, 'target')
    expect(Date.now() - start).toBeLessThan(100)
    expect(result).toContain('<mark class="highlight">target</mark>')
  })
})
