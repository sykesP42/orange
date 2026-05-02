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

    it('activeTab defaults to all', () => {
      expect(createWrapper().vm.activeTab).toBe('all')
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

  describe('search functionality', () => {
    it('commands list is populated', () => {
      const w = createWrapper()
      expect(w.vm.commands.length).toBeGreaterThan(0)
    })

    it('pages list is populated', () => {
      const w = createWrapper()
      expect(w.vm.pages.length).toBeGreaterThan(0)
    })

    it('filteredCommands returns all commands when no query', () => {
      const w = createWrapper()
      expect(w.vm.filteredCommands.length).toBe(w.vm.commands.length)
    })

    it('filteredPages returns all pages when no query', () => {
      const w = createWrapper()
      expect(w.vm.filteredPages.length).toBe(w.vm.pages.length)
    })
  })

  describe('highlightMatch utility', () => {
    it('returns original text for empty query', () => {
      expect(createWrapper().vm.highlightMatch('test', '')).toBe('test')
    })

    it('returns original text when text is empty', () => {
      expect(createWrapper().vm.highlightMatch('', 'test')).toBe('')
    })

    it('highlights English substring match', () => {
      const w = createWrapper()
      w.vm.query = 'hello'
      const result = w.vm.highlightMatch('hello world')
      expect(result).toContain('<mark')
    })
  })
})
