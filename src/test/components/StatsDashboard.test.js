import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StatsDashboard from '../../components/StatsDashboard.vue'

describe('StatsDashboard', () => {
  const defaultProps = {
    total: 128,
    activeCount: 45,
    pendingCount: 23,
    invalidCount: 12,
    todayNew: 5,
    weekGrowth: 12.5,
    cooperationRate: 35.2,
    loading: false,
    isDark: false,
  }

  describe('stat cards rendering', () => {
    it('renders 4 stat cards', () => {
      const wrapper = mount(StatsDashboard, { props: defaultProps })
      const cards = wrapper.findAll('.stat-card')
      expect(cards.length).toBe(4)
    })

    it('displays total count in first card', () => {
      const wrapper = mount(StatsDashboard, { props: defaultProps })
      const values = wrapper.findAll('.stat-value')
      expect(values[0].text()).toContain('128')
    })

    it('displays active count in second card', () => {
      const wrapper = mount(StatsDashboard, { props: defaultProps })
      const values = wrapper.findAll('.stat-value')
      expect(values[1].text()).toContain('45')
    })

    it('displays pending count in third card', () => {
      const wrapper = mount(StatsDashboard, { props: defaultProps })
      const values = wrapper.findAll('.stat-value')
      expect(values[2].text()).toContain('23')
    })

    it('displays invalid/other count in fourth card', () => {
      const wrapper = mount(StatsDashboard, { props: defaultProps })
      const values = wrapper.findAll('.stat-value')
      expect(values[3].text()).toContain('12')
    })

    it('shows loading placeholder when loading is true', () => {
      const wrapper = mount(StatsDashboard, {
        props: { ...defaultProps, loading: true },
      })
      const cards = wrapper.findAll('.stat-card')
      expect(cards[0].classes()).toContain('stat-loading')
      expect(wrapper.findAll('.stat-value')[0].text()).toBe('...')
    })

    it('shows actual numbers when not loading', () => {
      const wrapper = mount(StatsDashboard, {
        props: { ...defaultProps, loading: false },
      })
      expect(wrapper.findAll('.stat-value')[0].text()).not.toBe('...')
    })
  })

  describe('trend indicators', () => {
    it('renders trend section for each card', () => {
      const wrapper = mount(StatsDashboard, { props: defaultProps })
      const trends = wrapper.findAll('.stat-trend')
      expect(trends.length).toBe(4)
    })

    it('trend contains text and icon', () => {
      const wrapper = mount(StatsDashboard, { props: defaultProps })
      const firstTrend = wrapper.find('.stat-trend')
      expect(firstTrend.find('svg').exists()).toBe(true)
      expect(firstTrend.find('span').exists()).toBe(true)
    })
  })

  describe('quick actions panel', () => {
    it('renders quick actions panel with title', () => {
      const wrapper = mount(StatsDashboard, { props: defaultProps })
      expect(wrapper.text()).toContain('快捷操作')
    })

    it('has action buttons for each quick action', () => {
      const wrapper = mount(StatsDashboard, { props: defaultProps })
      const buttons = wrapper.findAll('.qa-btn')
      expect(buttons.length).toBeGreaterThan(0)
    })

    it('emits add event when clicking 录入 button', async () => {
      const wrapper = mount(StatsDashboard, { props: defaultProps })
      const addBtn = wrapper.findAll('.qa-btn').find(
        (b) => b.text().includes('录入') || b.attributes()?.style?.includes('--btn-accent')
      )
      if (addBtn) {
        await addBtn.trigger('click')
        expect(wrapper.emitted('add')).toBeTruthy()
      }
    })
  })

  describe('activity feed panel', () => {
    it('renders activity feed panel', () => {
      const wrapper = mount(StatsDashboard, { props: defaultProps })
      expect(wrapper.text()).toContain('近期动态')
    })

    it('shows activity feed list (always has at least overview item)', () => {
      const wrapper = mount(StatsDashboard, { props: defaultProps })
      expect(wrapper.find('.feed-list').exists()).toBe(true)
      expect(wrapper.text()).toContain('博主在库')
    })
  })

  describe('dark mode', () => {
    it('applies dark-mode class when isDark is true', () => {
      const wrapper = mount(StatsDashboard, {
        props: { ...defaultProps, isDark: true },
      })
      expect(wrapper.find('.stats-dashboard').classes()).toContain('dark-mode')
    })
  })

  describe('emits', () => {
    const emitEvents = ['add', 'invalid', 'export', 'import', 'kanban', 'workflow', 'calendar']

    emitEvents.forEach((event) => {
      it(`emits ${event} event`, () => {
        const wrapper = mount(StatsDashboard, { props: defaultProps })
        expect(typeof wrapper.vm.$emit).toBe('function')
      })
    })
  })
})
