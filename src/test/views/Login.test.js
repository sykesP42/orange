import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'

vi.mock('../../api', () => ({
  loginAPI: vi.fn().mockResolvedValue({ code: 200, data: { token: 'test-token' } }),
  registerAPI: vi.fn().mockResolvedValue({ code: 200 }),
  forgotPasswordAPI: vi.fn().mockResolvedValue({ code: 200 }),
}))

vi.mock('../../stores/notification', () => ({
  useNotification: () => ({ success: vi.fn(), error: vi.fn(), warning: vi.fn(), info: vi.fn() }),
}))

vi.mock('../../stores/user', () => ({
  useUserStore: () => ({ setToken: vi.fn(), fetchUserInfo: vi.fn() }),
}))

vi.mock('vue-router', () => ({ useRouter: () => ({ push: vi.fn() }) }))

const Login = (await import('../../views/Login.vue')).default

describe('Login.vue', () => {
  const createWrapper = () =>
    mount(Login, {
      global: { stubs: { 'router-link': { template: '<a><slot /></a>' } } },
    })

  describe('rendering', () => {
    it('renders login page container', () => {
      expect(createWrapper().find('.login-page').exists()).toBe(true)
    })

    it('renders canvas element for background animation', () => {
      expect(createWrapper().find('.geo-canvas').exists()).toBe(true)
    })

    it('renders login card container', () => {
      expect(createWrapper().find('.login-card').exists()).toBe(true)
    })
  })

  describe('version display', () => {
    it('shows version number containing "1.0"', () => {
      expect(createWrapper().text()).toContain('1.0')
    })
  })

  describe('login form fields', () => {
    it('has text input for username', () => {
      const w = createWrapper()
      expect(w.find('input[type="text"]').exists()).toBe(true)
    })

    it('has password input field', () => {
      expect(createWrapper().find('input[type="password"]').exists()).toBe(true)
    })

    it('username input binds to form.username', async () => {
      const w = createWrapper()
      await w.find('input[type="text"]').setValue('testuser')
      expect(w.vm.form.username).toBe('testuser')
    })

    it('password input binds to form.password', async () => {
      const w = createWrapper()
      await w.find('input[type="password"]').setValue('pass123')
      expect(w.vm.form.password).toBe('pass123')
    })
  })

  describe('submit button', () => {
    it('has submit button with 登 text', () => {
      const btn = createWrapper().find('.login-btn')
      expect(btn.exists()).toBe(true)
      expect(btn.text()).toContain('登')
    })
  })

  describe('form switching defaults', () => {
    it('shows login form by default (not register, not forgot)', () => {
      const vm = createWrapper().vm
      expect(vm.showRegister).toBe(false)
      expect(vm.showForgot).toBe(false)
    })
  })

  describe('keyboard shortcuts hint area', () => {
    it('contains Ctrl+K hint text', () => {
      expect(createWrapper().text()).toContain('Ctrl+K')
    })
  })

  describe('password toggle', () => {
    it('has password visibility toggle element', () => {
      const w = createWrapper()
      const hasToggle = w.find('.toggle-password').exists() ||
                        w.find('[class*="eye"]').exists() ||
                        w.find('[class*="toggle"]').exists()
      expect(hasToggle).toBe(true)
    })
  })

  describe('form validation guard', () => {
    it('handleLogin exists as method', () => {
      expect(typeof createWrapper().vm.handleLogin).toBe('function')
    })
  })
})
