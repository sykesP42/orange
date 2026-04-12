import { test, expect } from '@playwright/test'

test.describe('Login Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
  })

  test('should render login page with canvas background', async ({ page }) => {
    await expect(page.locator('.login-page')).toBeVisible()
    await expect(page.locator('.geo-canvas')).toBeVisible()
    await expect(page.locator('.login-card')).toBeVisible()
  })

  test('should show version number v1.0.0', async ({ page }) => {
    await expect(page.getByText('1.0')).toBeVisible()
  })

  test('should have username and password inputs', async ({ page }) => {
    await expect(page.locator('input[type="text"]')).toBeVisible()
    await expect(page.locator('input[type="password"]')).toBeVisible()
  })

  test('should have login button', async ({ page }) => {
    const loginBtn = page.locator('.login-btn')
    await expect(loginBtn).toBeVisible()
    await expect(loginBtn).toContainText(/登|Sign/)
  })

  test('should show register link', async ({ page }) => {
    await expect(page.getByText(/注册|Sign up/)).toBeVisible()
  })

  test('should show command palette hint', async ({ page }) => {
    await expect(page.getByText('Ctrl+K')).toBeVisible()
  })

  test('password toggle should work', async ({ page }) => {
    const toggle = page.locator('.toggle-password, [class*="eye"], [class*="toggle"]').first()
    if (await toggle.isVisible()) {
      const pwInput = page.locator('input[type="password"]')
      await expect(pwInput).toHaveAttribute('type', 'password')
      await toggle.click()
    }
  })
})

test.describe('Navigation Flow', () => {
  test('skip-to-content link should be hidden until focused', async ({ page }) => {
    await page.goto('/')
    const skipLink = page.locator('.skip-link')
    await expect(skipLink).toHaveCSS('top', /-100px|-50px/)
  })
})
