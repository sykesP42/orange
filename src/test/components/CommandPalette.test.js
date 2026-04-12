import { describe, it, expect, vi, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'

const CommandPalette = (await import('../../components/CommandPalette.vue')).default

describe('CommandPalette', () => {
  const createWrapper = (propsOverrides = {}) =>
    mount(CommandPalette, {
      props: { visible: true, ...propsOverrides },
      global: {
        stubs: { 'router-link': { template: '<a><slot /></a>' } },
        provide: { router: {} },
      },
      attachTo: document.body,
    })

  afterEach(() => { document.body.innerHTML = '' })

  describe('rendering', () => {
    it('mounts without error when visible=true', () => {
      const w = createWrapper()
      expect(w.exists()).toBe(true)
    })

    it('isOpen computed returns true when visible prop is true', () => {
      expect(createWrapper().vm.isOpen).toBe(true)
    })
  })

  describe('search input state', () => {
    it('query ref exists as empty string by default', () => {
      expect(createWrapper().vm.query).toBe('')
    })

    it('activeTab defaults to commands', () => {
      expect(createWrapper().vm.activeTab).toBe('commands')
    })

    it('selectedIndex is a number (0)', () => {
      expect(typeof createWrapper().vm.selectedIndex).toBe('number')
    })
  })

  describe('close behavior', () => {
    it('close method exists', () => {
      expect(typeof createWrapper().vm.close).toBe('function')
    })
  })

  describe('search history integration', () => {
    it('searchHistory is an array', () => {
      const w = createWrapper()
      expect(Array.isArray(w.vm.searchHistory)).toBe(true)
    })

    it('saveToHistory exists', () => {
      expect(typeof createWrapper().vm.saveToHistory).toBe('function')
    })

    it('clearHistory exists', () => {
      expect(typeof createWrapper().vm.clearHistory).toBe('function')
    })
  })

  describe('highlightMatch utility', () => {
    it('wraps matched text in <mark class="hl">', () => {
      const w = createWrapper()
      const result = w.vm.highlightMatch('录入博主', '录')
      expect(result).toContain('<mark')
      expect(result).toContain('hl')
    })

    it('returns original text on no match', () => {
      expect(createWrapper().vm.highlightMatch('hello', 'xyz')).toBe('hello')
    })

    it('returns original text for empty query', () => {
      expect(createWrapper().vm.highlightMatch('test', '')).toBe('test')
    })
  })
})
