import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TagCloud from '../../components/TagCloud.vue'

describe('TagCloud', () => {
  const sampleTags = [
    { name: '美妆', count: 45, color: '#f97316' },
    { name: '科技', count: 32, color: '#3b82f6' },
    { name: '游戏', count: 28, color: '#10b981' },
    { name: '美食', count: 15, color: '#8b5cf6' },
    { name: '旅行', count: 8, color: '#ec4899' },
  ]

  const defaultProps = {
    tags: sampleTags,
    selectedTag: null,
    loading: false,
    showRecommend: true,
    isDark: false,
  }

  describe('cloud view (default)', () => {
    it('renders tags in cloud view by default', () => {
      const wrapper = mount(TagCloud, { props: defaultProps })
      expect(wrapper.find('.cloud-view').exists()).toBe(true)
    })

    it('renders all tags as cloud-tag elements', () => {
      const wrapper = mount(TagCloud, { props: defaultProps })
      const cloudTags = wrapper.findAll('.cloud-tag')
      expect(cloudTags.length).toBe(sampleTags.length)
    })

    it('displays tag name and count', () => {
      const wrapper = mount(TagCloud, { props: defaultProps })
      expect(wrapper.text()).toContain('美妆')
      expect(wrapper.text()).toContain('45')
    })

    it('highlights selected tag', () => {
      const wrapper = mount(TagCloud, {
        props: { ...defaultProps, selectedTag: '科技' },
      })
      const selected = wrapper.findAll('.cloud-tag.active')
      expect(selected.length).toBe(1)
      expect(selected[0].text()).toContain('科技')
    })

    it('shows recommend badge for recommended tags', () => {
      const tagsWithRec = [...sampleTags]
      tagsWithRec[0].recommended = true
      const wrapper = mount(TagCloud, {
        props: { ...defaultProps, tags: tagsWithRec },
      })
      expect(wrapper.find('.recommend-badge').exists()).toBe(true)
      expect(wrapper.find('.recommend-badge').text()).toBe('AI')
    })

    it('emits select event when clicking a tag', async () => {
      const wrapper = mount(TagCloud, { props: defaultProps })
      await wrapper.findAll('.cloud-tag')[0].trigger('click')
      expect(wrapper.emitted('select')).toBeTruthy()
      expect(wrapper.emitted('select')[0]).toEqual(['美妆'])
    })
  })

  describe('list view', () => {
    it('switches to list view on toggle click', async () => {
      const wrapper = mount(TagCloud, { props: defaultProps })
      await wrapper.find('.toggle-btn').trigger('click')
      expect(wrapper.find('.list-view').exists()).toBe(true)
      expect(wrapper.find('.cloud-view').exists()).toBe(false)
    })

    it('renders list items in list mode', async () => {
      const wrapper = mount(TagCloud, { props: defaultProps })
      await wrapper.find('.toggle-btn').trigger('click')
      const items = wrapper.findAll('.list-tag-item')
      expect(items.length).toBe(sampleTags.length)
    })
  })

  describe('empty state', () => {
    it('shows empty state when no tags and not loading', () => {
      const wrapper = mount(TagCloud, {
        props: { ...defaultProps, tags: [], loading: false },
      })
      expect(wrapper.find('.empty-state').exists()).toBe(true)
      expect(wrapper.text()).toContain('暂无标签')
    })

    it('hides empty state when loading', () => {
      const wrapper = mount(TagCloud, {
        props: { ...defaultProps, tags: [], loading: true },
      })
      expect(wrapper.find('.empty-state').exists()).toBe(false)
      expect(wrapper.find('.loading-spinner').exists()).toBe(true)
    })
  })

  describe('loading state', () => {
    it('shows spinner when loading is true', () => {
      const wrapper = mount(TagCloud, {
        props: { ...defaultProps, loading: true },
      })
      expect(wrapper.find('.loading-spinner').exists()).toBe(true)
      expect(wrapper.find('.spinner').exists()).toBe(true)
    })

    it('hides cloud view when loading', () => {
      const wrapper = mount(TagCloud, {
        props: { ...defaultProps, loading: true },
      })
      expect(wrapper.find('.cloud-view').exists()).toBe(false)
    })
  })

  describe('refresh button', () => {
    it('shows refresh button when showRecommend is true', () => {
      const wrapper = mount(TagCloud, {
        props: { ...defaultProps, showRecommend: true },
      })
      expect(wrapper.find('.refresh-btn').exists()).toBe(true)
    })

    it('hides refresh button when showRecommend is false', () => {
      const wrapper = mount(TagCloud, {
        props: { ...defaultProps, showRecommend: false },
      })
      expect(wrapper.find('.refresh-btn').exists()).toBe(false)
    })

    it('emits refresh event on click', async () => {
      const wrapper = mount(TagCloud, {
        props: { ...defaultProps, showRecommend: true },
      })
      await wrapper.find('.refresh-btn').trigger('click')
      expect(wrapper.emitted('refresh')).toBeTruthy()
    })
  })

  describe('dark mode', () => {
    it('applies dark-mode class when isDark is true', () => {
      const wrapper = mount(TagCloud, {
        props: { ...defaultProps, isDark: true },
      })
      expect(wrapper.find('.tag-cloud-panel').classes()).toContain('dark-mode')
    })
  })

  describe('tooltip on hover', () => {
    it('shows tooltip on mouseenter', async () => {
      const wrapper = mount(TagCloud, { props: defaultProps })
      const tag = wrapper.findAll('.cloud-tag')[0]
      await tag.trigger('mouseenter')
      expect(wrapper.find('.tag-tooltip').exists()).toBe(true)
    })

    it('hides tooltip on mouseleave', async () => {
      const wrapper = mount(TagCloud, { props: defaultProps })
      const tag = wrapper.findAll('.cloud-tag')[0]
      await tag.trigger('mouseenter')
      await tag.trigger('mouseleave')
      expect(wrapper.find('.tag-tooltip').exists()).toBe(false)
    })
  })
})
