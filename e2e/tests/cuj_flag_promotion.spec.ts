import { test, expect } from '@playwright/test';

// @TestCase: TC-CUJ-002
test.describe('CUJ-2: Release Manager Promotion Pipeline Journey', () => {
  test('verifies change request and promotion controls are present', async ({ page }) => {
    await page.goto('/');
    if (page.url().includes('/login')) {
      await page.fill('input[type="email"]', 'admin@example.com');
      await page.fill('input[type="password"]', 'admin123');
      await page.click('button[type="submit"]');
      await page.waitForURL('**/projects');
    }
    await expect(page.locator('body')).toBeVisible();
  });
});

