import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SkeletonLoader from '../../components/SkeletonLoader.vue'

describe('SkeletonLoader', () => {
  describe('rendering variants', () => {
    it('renders blogger-card variant with correct count of cards', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'blogger-card', count: 6 },
      })
      const cards = wrapper.findAll('.ske-card')
      expect(cards.length).toBe(6)
    })

    it('renders blogger-list variant with header and rows', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'blogger-list', count: 8 },
      })
      expect(wrapper.find('.ske-list-header').exists()).toBe(true)
      const rows = wrapper.findAll('.ske-list-row')
      expect(rows.length).toBe(8)
    })

    it('renders stats variant with 4 stat cards', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'stats' },
      })
      const statCards = wrapper.findAll('.ske-stat-card')
      expect(statCards.length).toBe(4)
    })

    it('renders table variant with correct columns and rows', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'table', columns: 5, rows: 10 },
      })
      expect(wrapper.find('.ske-thead').exists()).toBe(true)
      const ths = wrapper.findAll('.ske-th')
      expect(ths.length).toBe(5)
      const trs = wrapper.findAll('.ske-tr')
      expect(trs.length).toBe(10)
    })

    it('renders detail variant with header and body sections', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'detail' },
      })
      expect(wrapper.find('.ske-detail-header').exists()).toBe(true)
      expect(wrapper.find('.ske-detail-body').exists()).toBe(true)
      expect(wrapper.find('.ske-detail-avatar-lg').exists()).toBe(true)
    })

    it('renders text (default) variant with count blocks', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'text', count: 4 },
      })
      const blocks = wrapper.findAll('.ske-block')
      expect(blocks.length).toBe(4)
    })
  })

  describe('dark mode', () => {
    it('applies dark-mode class when isDark is true', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'blogger-card', count: 3, isDark: true },
      })
      expect(wrapper.classes()).toContain('dark-mode')
    })

    it('does not apply dark-mode class by default', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'blogger-card', count: 3 },
      })
      expect(wrapper.classes()).not.toContain('dark-mode')
    })
  })

  describe('variant class', () => {
    it('applies variant as CSS class', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'stats' },
      })
      expect(wrapper.classes()).toContain('stats')
    })

    it('applies blogger-card class for card variant', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'blogger-card' },
      })
      expect(wrapper.classes()).toContain('blogger-card')
    })
  })

  describe('shimmer animation elements', () => {
    it('has shimmer-animated skeleton lines in card variant', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'blogger-card', count: 1 },
      })
      const lines = wrapper.findAll('.ske-line')
      expect(lines.length).toBeGreaterThan(0)
    })

    it('has avatar placeholder in card variant', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'blogger-card', count: 1 },
      })
      expect(wrapper.find('.ske-avatar').exists()).toBe(true)
    })

    it('has checkbox in card variant', () => {
      const wrapper = mount(SkeletonLoader, {
        props: { variant: 'blogger-card', count: 1 },
      })
      expect(wrapper.find('.ske-checkbox').exists()).toBe(true)
    })
  })
})
