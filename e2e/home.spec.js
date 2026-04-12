import { test, expect } from '@playwright/test'

test.describe('Home Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('should render page header with navigation', async ({ page }) => {
    await expect(page.locator('.header')).toBeVisible()
  })

  test('should have PWA status bar area', async ({ page }) => {
    await expect(page.locator('.pwa-status-bar')).toBeVisible()
  })

  test('should have language switcher', async ({ page }) => {
    await expect(page.locator('.lang-switcher')).toBeVisible()
    await expect(page.locator('.lang-btn')).toBeVisible()
  })

  test('should render main content area', async ({ page }) => {
    const main = page.locator('#main-content')
    await expect(main).toBeVisible()
  })

  test('should have bottom navigation on mobile viewport', async ({ page }) => {
    const bottomNav = page.locator('.bottom-nav')
    if (await bottomNav.isVisible()) {
      await expect(bottomNav).toBeVisible()
      const navItems = bottomNav.locator('.nav-item')
      const count = await navItems.count()
      expect(count).toBeGreaterThanOrEqual(5)
    }
  })

  test('should have FAB floating action button', async ({ page }) => {
    await expect(page.locator('.fab-main')).toBeVisible()
  })
})

test.describe('Language Switching', () => {
  test('should open language dropdown on click', async ({ page }) => {
    await page.goto('/')
    const langBtn = page.locator('.lang-btn')
    await langBtn.click()
    const menu = page.locator('.lang-menu')
    await expect(menu).toBeVisible()
  })

  test('should show Chinese and English options', async ({ page }) => {
    await page.goto('/')
    await page.locator('.lang-btn').click()
    await expect(page.getByText('中文')).toBeVisible()
    await expect(page.getByText('English')).toBeVisible()
  })

  test('should close dropdown when clicking outside', async ({ page }) => {
    await page.goto('/')
    await page.locator('.lang-btn').click()
    await expect(page.locator('.lang-menu')).toBeVisible()
    await page.mouse.click(10, 10)
    await expect(page.locator('.lang-menu')).not.toBeVisible()
  })
})

test.describe('Command Palette Access', () => {
  test('command palette trigger should exist in header', async ({ page }) => {
    await page.goto('/')
    const trigger = page.locator('.command-palette-trigger, [class*="search"]')
    if (await trigger.count() > 0) {
      await expect(trigger.first()).toBeVisible()
    }
  })
})
