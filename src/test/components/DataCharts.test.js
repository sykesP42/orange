import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import DataCharts from '../../components/DataCharts.vue'

const mockBloggers = [
  { id: 1, nickname: '博主A', status: '初次联系', category: '美妆', platform: '小红书', create_time: '2026-01-15T10:00:00Z', tags: ['公'] },
  { id: 2, nickname: '博主B', status: '洽谈中', category: '美食', platform: '抖音', create_time: '2026-02-20T10:00:00Z', tags: ['公', 'KOL'] },
  { id: 3, nickname: '博主C', status: '已合作', category: '美妆', platform: '小红书', create_time: '2026-03-10T10:00:00Z', tags: [] },
  { id: 4, nickname: '博主D', status: '初次联系', category: '数码', platform: 'B站', create_time: '2026-04-01T10:00:00Z', tags: ['KOL'] },
  { id: 5, nickname: '博主E', status: '洽谈中', category: '美食', platform: '抖音', create_time: '2026-04-10T10:00:00Z', tags: ['公'] },
]

const mockChartInstance = {
  setOption: vi.fn(),
  resize: vi.fn(),
  dispose: vi.fn(),
}

vi.mock('echarts', () => {
  const echartsModule = {
    init: vi.fn(() => mockChartInstance),
    graphic: {
      LinearGradient: class {
        constructor(...args) { Object.assign(this, {}) }
      },
    },
  }
  return echartsModule
})

describe('DataCharts', () => {
  it('should not render when bloggers is empty', () => {
    const wrapper = mount(DataCharts, {
      props: { bloggers: [], isDark: false }
    })
    expect(wrapper.find('.data-charts').exists()).toBe(false)
  })

  it('should render when bloggers has data', () => {
    const wrapper = mount(DataCharts, {
      props: { bloggers: mockBloggers, isDark: false }
    })
    expect(wrapper.find('.data-charts').exists()).toBe(true)
    expect(wrapper.find('.charts-header').exists()).toBe(true)
  })

  it('should display correct blogger count in subtitle', () => {
    const wrapper = mount(DataCharts, {
      props: { bloggers: mockBloggers, isDark: false }
    })
    expect(wrapper.text()).toContain('5 条数据')
  })

  it('should render all 4 chart cards', () => {
    const wrapper = mount(DataCharts, {
      props: { bloggers: mockBloggers, isDark: false }
    })

    expect(wrapper.findAll('.chart-card')).toHaveLength(4)
    expect(wrapper.text()).toContain('状态分布')
    expect(wrapper.text()).toContain('分类分布')
    expect(wrapper.text()).toContain('平台分布')
    expect(wrapper.text()).toContain('近期趋势')
  })

  it('should have chart containers for echarts instances', () => {
    const wrapper = mount(DataCharts, {
      props: { bloggers: mockBloggers, isDark: false }
    })

    expect(wrapper.findAll('.chart-container')).toHaveLength(4)
  })

  it('should show full-width class on trend chart card', () => {
    const wrapper = mount(DataCharts, {
      props: { bloggers: mockBloggers, isDark: false }
    })

    const trendCard = wrapper.find('.trend-chart-card')
    expect(trendCard.classes()).toContain('full-width')
  })

  it('should respond to bloggers data change', async () => {
    const wrapper = mount(DataCharts, {
      props: { bloggers: mockBloggers.slice(0, 2), isDark: false }
    })

    expect(wrapper.text()).toContain('2 条数据')

    await wrapper.setProps({ bloggers: mockBloggers })
    expect(wrapper.text()).toContain('5 条数据')
  })

  it('should have correct ARIA labels on header', () => {
    const wrapper = mount(DataCharts, {
      props: { bloggers: mockBloggers, isDark: false }
    })

    const svg = wrapper.find('.charts-header svg')
    expect(svg.attributes('aria-hidden')).toBe('true')
  })
})
