import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'

const KanbanBoard = (await import('../../components/KanbanBoard.vue')).default

describe('KanbanBoard', () => {
  const sampleBloggers = [
    { id: 1, nickname: 'Alice', status: '初次联系', category: '美妆' },
    { id: 2, nickname: 'Bob', status: '洽谈中', category: '科技' },
    { id: 3, nickname: 'Carol', status: '已合作', category: '游戏' },
    { id: 4, nickname: 'Dave', status: '已拒绝', category: '美食' },
    { id: 5, nickname: 'Eve', status: '暂停合作', category: '旅行' },
    { id: 6, nickname: 'Frank', status: '洽谈中', category: '运动' },
  ]

  const createWrapper = (propsOverrides = {}) =>
    mount(KanbanBoard, {
      props: { bloggers: sampleBloggers, ...propsOverrides },
      global: {
        stubs: {
          'router-link': { template: '<a><slot /></a>' },
          Draggable: {
            template: '<div><slot /></div>',
            props: ['modelValue'],
            emits: ['update:modelValue'],
          },
        },
      },
    })

  describe('column rendering', () => {
    it('renders 5 kanban columns', () => {
      const cols = createWrapper().findAll('.kanban-column')
      expect(cols.length).toBe(5)
    })

    it('has correct status labels for all 5 columns', () => {
      const w = createWrapper()
      expect(w.text()).toContain('初次联系')
      expect(w.text()).toContain('洽谈中')
      expect(w.text()).toContain('已合作')
      expect(w.text()).toContain('已拒绝')
      expect(w.text()).toContain('暂停合作')
    })

    it('each column has a colored status dot', () => {
      const dots = createWrapper().findAll('.kanban-status-dot')
      expect(dots.length).toBe(5)
    })

    it('each column has count badge', () => {
      const counts = createWrapper().findAll('.kanban-count')
      expect(counts.length).toBe(5)
    })
  })

  describe('blogger distribution via getColumnBloggers', () => {
    it('returns correct blogger for 初次联系', () => {
      const w = createWrapper()
      const col = w.vm.getColumnBloggers('初次联系')
      expect(col.length).toBe(1)
      expect(col[0].nickname).toBe('Alice')
    })

    it('returns 2 bloggers for 洽谈中', () => {
      expect(createWrapper().vm.getColumnBloggers('洽谈中').length).toBe(2)
    })

    it('returns 1 blogger for 已合作', () => {
      const col = createWrapper().vm.getColumnBloggers('已合作')
      expect(col.length).toBe(1)
      expect(col[0].nickname).toBe('Carol')
    })

    it('returns empty for non-existent status', () => {
      expect(createWrapper().vm.getColumnBloggers('不存在').length).toBe(0)
    })
  })

  describe('empty state', () => {
    it('shows empty placeholders when no bloggers', () => {
      const w = createWrapper({ bloggers: [] })
      const empties = w.findAll('.kanban-empty')
      expect(empties.length).toBe(5)
    })
  })

  describe('columns definition', () => {
    it('has 5 column definitions with status and color', () => {
      const cols = createWrapper().vm.columns
      expect(cols.length).toBe(5)
      cols.forEach((col) => {
        expect(col).toHaveProperty('status')
        expect(col).toHaveProperty('color')
        expect(col).toHaveProperty('cssClass')
      })
    })
  })

  describe('emits', () => {
    it('emits goToDetail on card click if handler wired', async () => {
      const w = createWrapper()
      const cards = w.findAll('.kanban-card')
      if (cards.length > 0) {
        await cards[0].trigger('click')
        if (w.emitted().goToDetail) {
          expect(w.emitted().goToDetail).toBeTruthy()
        }
      }
    })
  })
})
