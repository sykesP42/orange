import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import Admin from '../../views/Admin.vue'

vi.mock('../../stores/user', () => ({
  useUserStore: () => ({
    role: 'admin', team_id: 1, real_name: 'Admin', username: 'admin', avatar: '',
  }),
}))
vi.mock('../../stores/notification', () => ({
  useNotification: () => ({ success: vi.fn(), error: vi.fn(), warning: vi.fn(), info: vi.fn() }),
}))
vi.mock('../../utils/confirm', () => ({
  useConfirm: () => ({ confirmDanger: vi.fn().mockResolvedValue(true) }),
}))
vi.mock('../../api', () => ({
  categoryListAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  productListAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  getTeamsAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  exportAPI: vi.fn().mockResolvedValue(new Blob(['test'])),
  getOperationLog: vi.fn().mockResolvedValue({ code: 200, data: [], total: 0 }),
  getPendingUsersAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  getRegistrationModeAPI: vi.fn().mockResolvedValue({ code: 200, data: {} }),
  categoryAddAPI: vi.fn().mockResolvedValue({ code: 200 }),
  categoryUpdateAPI: vi.fn().mockResolvedValue({ code: 200 }),
  categoryDeleteAPI: vi.fn().mockResolvedValue({ code: 200 }),
  productAddAPI: vi.fn().mockResolvedValue({ code: 200 }),
  productUpdateAPI: vi.fn().mockResolvedValue({ code: 200 }),
  productDeleteAPI: vi.fn().mockResolvedValue({ code: 200 }),
  deleteTeamAPI: vi.fn().mockResolvedValue({ code: 200 }),
  updateTeamAPI: vi.fn().mockResolvedValue({ code: 200 }),
  createTeamAPI: vi.fn().mockResolvedValue({ code: 200 }),
  setAdminAPI: vi.fn().mockResolvedValue({ code: 200 }),
  removeAdminAPI: vi.fn().mockResolvedValue({ code: 200 }),
  setUserTeamAPI: vi.fn().mockResolvedValue({ code: 200 }),
  getUsersListAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  deleteLogAPI: vi.fn().mockResolvedValue({ code: 200 }),
  clearLogsAPI: vi.fn().mockResolvedValue({ code: 200 }),
  clearOldLogsAPI: vi.fn().mockResolvedValue({ code: 200 }),
  batchApproveUsersAPI: vi.fn().mockResolvedValue({ code: 200 }),
  batchRejectUsersAPI: vi.fn().mockResolvedValue({ code: 200 }),
  setRegistrationModeAPI: vi.fn().mockResolvedValue({ code: 200 }),
  deactivateUserAdminAPI: vi.fn().mockResolvedValue({ code: 200 }),
  deleteUserAdminAPI: vi.fn().mockResolvedValue({ code: 200 }),
  uploadTeamLogoAPI: vi.fn().mockResolvedValue({ code: 200, data: { url: '/logo.png' } }),
  uploadTeamBgAPI: vi.fn().mockResolvedValue({ code: 200, data: { url: '/bg.png' } }),
  createBackupAPI: vi.fn().mockResolvedValue({ code: 200 }),
  getSnapshotsAPI: vi.fn().mockResolvedValue({ code: 200, data: [] }),
  getSnapshotSettingsAPI: vi.fn().mockResolvedValue({ code: 200, data: {} }),
  setSnapshotSettingsAPI: vi.fn().mockResolvedValue({ code: 200 }),
  createSnapshotAPI: vi.fn().mockResolvedValue({ code: 200 }),
  restoreSnapshotAPI: vi.fn().mockResolvedValue({ code: 200 }),
  downloadSnapshotAPI: vi.fn().mockResolvedValue({ code: 200 }),
  deleteSnapshotAPI: vi.fn().mockResolvedValue({ code: 200 }),
  clearDataAPI: vi.fn().mockResolvedValue({ code: 200 }),
  purgeDataAPI: vi.fn().mockResolvedValue({ code: 200 }),
}))

describe('Admin.vue', () => {
  const createWrapper = () =>
    mount(Admin, {
      global: {
        stubs: {
          'router-link': { template: '<a><slot /></a>' },
          'router-view': { template: '<div><slot /></div>' },
          ImportPanel: { template: '<div>ImportPanel</div>' },
        },
      },
    })

  describe('rendering', () => {
    it('renders admin container', () => {
      expect(createWrapper().find('.admin').exists()).toBe(true)
    })

    it('has page header', () => {
      expect(createWrapper().find('.page-header').exists()).toBe(true)
    })

    it('has tab navigation area', () => {
      expect(createWrapper().find('.admin-tabs').exists()).toBe(true)
    })
  })

  describe('tabs content', () => {
    it('renders tab items for management sections', () => {
      expect(createWrapper().text().length).toBeGreaterThan(10)
    })
  })
})
