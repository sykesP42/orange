import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import PWAStatusBar from '../../components/PWAStatusBar.vue'

describe('PWAStatusBar', () => {
  describe('offline banner', () => {
    it('shows offline banner when isOnline is false', () => {
      const wrapper = mount(PWAStatusBar, {
        props: { isOnline: false, showUpdateBanner: false, isInstallable: false },
      })
      expect(wrapper.find('.offline-banner').exists()).toBe(true)
      expect(wrapper.text()).toContain('离线模式')
    })

    it('hides offline banner when isOnline is true', () => {
      const wrapper = mount(PWAStatusBar, {
        props: { isOnline: true, showUpdateBanner: false, isInstallable: false },
      })
      expect(wrapper.find('.offline-banner').exists()).toBe(false)
    })
  })

  describe('update banner', () => {
    it('shows update banner when showUpdateBanner is true and online', () => {
      const wrapper = mount(PWAStatusBar, {
        props: { isOnline: true, showUpdateBanner: true, isInstallable: false },
      })
      expect(wrapper.find('.update-banner').exists()).toBe(true)
      expect(wrapper.text()).toContain('新版本')
    })

    it('update button emits update event', async () => {
      const wrapper = mount(PWAStatusBar, {
        props: { isOnline: true, showUpdateBanner: true, isInstallable: false },
      })
      await wrapper.find('.update-btn').trigger('click')
      expect(wrapper.emitted('update')).toBeTruthy()
    })

    it('dismiss button exists and is clickable', async () => {
      const wrapper = mount(PWAStatusBar, {
        props: { isOnline: true, showUpdateBanner: true, isInstallable: false },
      })
      const btn = wrapper.find('.update-dismiss')
      expect(btn.exists()).toBe(true)
      await btn.trigger('click')
    })
  })

  describe('install banner', () => {
    it('shows install banner when isInstallable is true and no other banners', () => {
      const wrapper = mount(PWAStatusBar, {
        props: { isOnline: true, showUpdateBanner: false, isInstallable: true },
      })
      expect(wrapper.find('.install-banner').exists()).toBe(true)
      expect(wrapper.text()).toContain('安装到桌面')
    })

    it('install button emits install event', async () => {
      const wrapper = mount(PWAStatusBar, {
        props: { isOnline: true, showUpdateBanner: false, isInstallable: true },
      })
      await wrapper.find('.install-btn').trigger('click')
      expect(wrapper.emitted('install')).toBeTruthy()
    })

    it('hides install banner when hideInstall is true', () => {
      const wrapper = mount(PWAStatusBar, {
        props: { isOnline: true, showUpdateBanner: false, isInstallable: true, hideInstall: true },
      })
      expect(wrapper.find('.install-banner').exists()).toBe(false)
    })
  })

  describe('priority order', () => {
    it('offline banner takes priority over update banner', () => {
      const wrapper = mount(PWAStatusBar, {
        props: { isOnline: false, showUpdateBanner: true, isInstallable: false },
      })
      expect(wrapper.find('.offline-banner').exists()).toBe(true)
      expect(wrapper.find('.update-banner').exists()).toBe(false)
    })

    it('update banner takes priority over install banner', () => {
      const wrapper = mount(PWAStatusBar, {
        props: { isOnline: true, showUpdateBanner: true, isInstallable: true },
      })
      expect(wrapper.find('.update-banner').exists()).toBe(true)
      expect(wrapper.find('.install-banner').exists()).toBe(false)
    })

    it('shows offline dot when offline but hideOfflineBar', () => {
      const wrapper = mount(PWAStatusBar, {
        props: { isOnline: false, showUpdateBanner: false, isInstallable: false, showOfflineBar: false },
      })
      expect(wrapper.find('.offline-dot').exists()).toBe(true)
    })
  })

  describe('root element', () => {
    it('applies is-offline class when offline', () => {
      const wrapper = mount(PWAStatusBar, {
        props: { isOnline: false },
      })
      expect(wrapper.find('.pwa-status-bar').classes()).toContain('is-offline')
    })
  })
})
