import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'

const FloatingActionBar = (await import('../../components/FloatingActionBar.vue')).default

describe('FloatingActionBar', () => {
  const createWrapper = (props = {}) => {
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/', component: { template: '<div/>' } }, { path: '/add', component: { template: '<div/>' } }],
    })
    return mount(FloatingActionBar, {
      props: { isDark: false, ...props },
      global: { plugins: [router] },
      attachTo: document.body,
    })
  }

  afterEach(() => { document.body.innerHTML = '' })

  describe('initial state', () => {
    it('renders FAB main button', () => {
      const w = createWrapper()
      expect(w.find('.fab-main').exists()).toBe(true)
    })

    it('has SVG icon inside main button', () => {
      const w = createWrapper()
      expect(w.find('.fab-main svg').exists()).toBe(true)
    })

    it('action items exist but are hidden via display:none (v-show)', () => {
      const w = createWrapper()
      const actions = w.findAll('.fab-action')
      expect(actions.length).toBe(5)
      expect(w.find('.fab-container').classes()).not.toContain('open')
    })

    it('backdrop does not exist initially (v-if)', () => {
      const w = createWrapper()
      expect(w.find('.fab-backdrop').exists()).toBe(false)
    })
  })

  describe('toggle behavior', () => {
    it('opens menu on click — adds open class and shows backdrop', async () => {
      const w = createWrapper()
      await w.find('.fab-main').trigger('click')
      expect(w.find('.fab-container').classes()).toContain('open')
      expect(w.find('.fab-backdrop').exists()).toBe(true)
    })

    it('closes menu on second click — removes open class', async () => {
      const w = createWrapper()
      await w.find('.fab-main').trigger('click')
      await w.find('.fab-main').trigger('click')
      expect(w.find('.fab-container').classes()).not.toContain('open')
    })
  })

  describe('action items', () => {
    it('has 5 action items defined', () => {
      const w = createWrapper()
      expect(w.findAll('.fab-action').length).toBe(5)
    })

    it('shows correct labels after opening', async () => {
      const w = createWrapper()
      await w.find('.fab-main').trigger('click')
      const labels = w.findAll('.fab-action-label').map((l) => l.text())
      expect(labels).toEqual(expect.arrayContaining(['录入博主', '命令面板', '日历', '工作流', '看板']))
    })
  })

  describe('backdrop close', () => {
    it('removes open class after clicking backdrop (Transition keeps element during animation)', async () => {
      const w = createWrapper()
      await w.find('.fab-main').trigger('click')
      expect(w.find('.fab-container').classes()).toContain('open')
      await w.find('.fab-backdrop').trigger('click')
      expect(w.find('.fab-container').classes()).not.toContain('open')
    })
  })

  describe('dark mode', () => {
    it('applies dark-mode class', () => {
      const r = createRouter({ history: createMemoryHistory(), routes: [] })
      const w = mount(FloatingActionBar, { props: { isDark: true }, global: { plugins: [r] } })
      expect(w.find('.fab-container').classes()).toContain('dark-mode')
    })
  })

  describe('accessibility', () => {
    it('main button has aria-label', () => {
      const w = createWrapper()
      expect(w.find('.fab-main').attributes('aria-label')).toBe('快捷操作菜单')
    })
  })
})
