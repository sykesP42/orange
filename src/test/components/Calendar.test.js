import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'

vi.mock('../../api', () => ({
  getCalendarEventsAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  createCalendarEventAPI: vi.fn().mockResolvedValue({ code: 200 }),
  updateCalendarEventAPI: vi.fn().mockResolvedValue({ code: 200 }),
  deleteCalendarEventAPI: vi.fn().mockResolvedValue({ code: 200 }),
}))

const Calendar = (await import('../../components/Calendar.vue')).default

describe('Calendar', () => {
  const createWrapper = () =>
    mount(Calendar, {
      global: { stubs: { 'router-link': { template: '<a><slot /></a>' } } },
    })

  describe('rendering', () => {
    it('renders calendar container', () => {
      const w = createWrapper()
      expect(w.find('.calendar-container').exists() || w.find('.calendar').exists() || w.find('[class*="calendar"]').exists()).toBe(true)
    })

    it('has year display in header', () => {
      const w = createWrapper()
      expect(w.text()).toContain(String(new Date().getFullYear()))
    })

    it('has navigation controls', () => {
      const w = createWrapper()
      expect(w.findAll('button').length).toBeGreaterThanOrEqual(2)
    })
  })

  describe('state', () => {
    it('viewMode defaults to month', () => {
      expect(createWrapper().vm.viewMode).toBe('month')
    })

    it('currentDate is a Date object', () => {
      expect(createWrapper().vm.currentDate).toBeInstanceOf(Date)
    })

    it('events is an array', () => {
      expect(Array.isArray(createWrapper().vm.events)).toBe(true)
    })
  })

  describe('upcoming events', () => {
    it('upcomingEvents is a computed property', () => {
      const w = createWrapper()
      expect(typeof w.vm.upcomingEvents).not.toBe('undefined')
    })

    it('upcomingEvents returns array', () => {
      expect(Array.isArray(createWrapper().vm.upcomingEvents)).toBe(true)
    })
  })

  describe('quick add event', () => {
    it('showQuickAdd ref exists', () => {
      expect(typeof createWrapper().vm.showQuickAdd).not.toBe('undefined')
    })

    it('newEvent reactive object has required fields', () => {
      const ev = createWrapper().vm.newEvent
      expect(ev).toHaveProperty('title')
      expect(ev).toHaveProperty('date')
      expect(ev).toHaveProperty('type')
      expect(ev).toHaveProperty('remark')
    })

    it('handleQuickAdd method exists', () => {
      expect(typeof createWrapper().vm.handleQuickAdd).toBe('function')
    })
  })

  describe('event methods', () => {
    it('getEventIcon method exists', () => {
      expect(typeof createWrapper().vm.getEventIcon).toBe('function')
    })

    it('getEventTypeName method exists', () => {
      expect(typeof createWrapper().vm.getEventTypeName).toBe('function')
    })
  })
})
