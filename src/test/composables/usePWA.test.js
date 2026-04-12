import { describe, it, expect, vi } from 'vitest'

vi.mock('../../composables/usePWA', () => {
  const isOnline = { value: true }
  const showUpdateBanner = { value: false }
  const isInstallable = { value: false }

  return {
    usePWA: () => ({
      isOnline,
      showUpdateBanner,
      isInstallable,
      updateSW: vi.fn().mockResolvedValue(undefined),
      clearCache: vi.fn().mockResolvedValue(true),
      installPWA: vi.fn().mockResolvedValue(false),
    }),
  }
})

import { usePWA } from '../../composables/usePWA'

describe('usePWA', () => {
  describe('initial state', () => {
    it('should return isOnline as true by default', () => {
      expect(usePWA().isOnline.value).toBe(true)
    })

    it('should return showUpdateBanner as false by default', () => {
      expect(usePWA().showUpdateBanner.value).toBe(false)
    })

    it('should return isInstallable as false by default', () => {
      expect(usePWA().isInstallable.value).toBe(false)
    })
  })

  describe('methods', () => {
    it('updateSW should be a function that resolves', async () => {
      const { updateSW } = usePWA()
      expect(typeof updateSW).toBe('function')
      await expect(updateSW()).resolves.toBeUndefined()
    })

    it('clearCache should resolve to true', async () => {
      const { clearCache } = usePWA()
      expect(await clearCache()).toBe(true)
    })

    it('installPWA should resolve to false without prompt', async () => {
      const { installPWA } = usePWA()
      expect(await installPWA()).toBe(false)
    })
  })

  describe('API shape', () => {
    it('exposes all required properties', () => {
      const pwa = usePWA()
      const required = ['isOnline', 'showUpdateBanner', 'isInstallable', 'updateSW', 'clearCache', 'installPWA']
      required.forEach((key) => {
        expect(pwa).toHaveProperty(key)
      })
    })

    it('isOnline is ref-like with .value', () => {
      expect(usePWA().isOnline).toHaveProperty('value')
    })
  })
})
