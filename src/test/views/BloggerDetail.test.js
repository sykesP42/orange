import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn() }),
  useRoute: () => ({ params: { id: '1' }, query: {}, path: '/blogger/1' }),
  createRouter: vi.fn(() => ({})),
  createWebHistory: vi.fn(),
}))

vi.mock('../../stores/user', () => ({
  useUserStore: () => ({ role: 'admin', username: 'admin', real_name: 'Admin' }),
}))
vi.mock('../../stores/notification', () => ({
  useNotification: () => ({ success: vi.fn(), error: vi.fn(), warning: vi.fn(), info: vi.fn() }),
}))
vi.mock('../../utils/confirm', () => ({
  useConfirm: () => ({ confirmDanger: vi.fn().mockResolvedValue(true) }),
}))
vi.mock('../../api', () => ({
  getBloggerDetail: vi.fn().mockResolvedValue({
    code: 200,
    data: {
      id: 1, nickname: 'TestBlogger', category: '美妆', status: '已合作',
      platform: '小红书', contact: '13800138000', wechat: 'wx_test',
      avatar: '/avatar.png', description: '测试博主描述', tags: [],
      product: null, custom_contact: null,
      create_time: '2024-01-15T10:00:00Z',
      owner_username: 'admin', owner_real_name: '管理员', owner_avatar: '',
      fans_count: 50,
    },
  }),
  followupAddAPI: vi.fn().mockResolvedValue({ code: 200 }),
  followupListAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  followupDeleteAPI: vi.fn().mockResolvedValue({ code: 200 }),
  followupUpdateAPI: vi.fn().mockResolvedValue({ code: 200 }),
  bloggerUpdateAPI: vi.fn().mockResolvedValue({ code: 200 }),
  bloggerDeleteAPI: vi.fn().mockResolvedValue({ code: 200 }),
  categoryListAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  requestBloggerTransfer: vi.fn().mockResolvedValue({ code: 200 }),
  cooperationListAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  cooperationAddAPI: vi.fn().mockResolvedValue({ code: 200 }),
  cooperationUpdateAPI: vi.fn().mockResolvedValue({ code: 200 }),
  cooperationDeleteAPI: vi.fn().mockResolvedValue({ code: 200 }),
  getBloggerStatusList: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  getTeamsAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  getBloggerLogsAPI: vi.fn().mockResolvedValue({ code: 200, data: [], total: 0 }),
  getPublicUsersAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  getTagsAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  getBloggerTagsAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  setBloggerTagsAPI: vi.fn().mockResolvedValue({ code: 200 }),
  getSimilarBloggersAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
}))

const BloggerDetail = (await import('../../views/BloggerDetail.vue')).default

describe('BloggerDetail.vue', () => {
  const createWrapper = () =>
    mount(BloggerDetail, {
      global: {
        stubs: {
          'router-link': { template: '<a><slot /></a>' },
          AvatarCropper: { template: '<div>AvatarCropper</div>' },
          TemplateSelector: { template: '<div>TemplateSelector</div>' },
          RouterView: { template: '<div><slot /></div>' },
        },
      },
    })

  describe('rendering', () => {
    it('renders detail page container', () => {
      expect(createWrapper().find('.detail-page').exists()).toBe(true)
    })

    it('component mounts without crashing', () => {
      const w = createWrapper()
      expect(w.exists()).toBe(true)
      expect(w.text().length).toBeGreaterThanOrEqual(0)
    })
  })

  describe('loading state', () => {
    it('shows skeleton or loading text initially', () => {
      const w = createWrapper()
      const hasSkeleton = w.find('.skeleton-loader').exists()
      const hasLoadingText = w.text().includes('加载')
      expect(hasSkeleton || hasLoadingText).toBe(true)
    })
  })

  describe('methods exist on instance', () => {
    it('handleDelete method exists', () => {
      expect(typeof createWrapper().vm.handleDelete).toBe('function')
    })

    it('addFollowup method exists', () => {
      expect(typeof createWrapper().vm.addFollowup).toBe('function')
    })
  })

  describe('component structure', () => {
    it('mounts without crashing', () => {
      expect(createWrapper().exists()).toBe(true)
    })
  })
})
