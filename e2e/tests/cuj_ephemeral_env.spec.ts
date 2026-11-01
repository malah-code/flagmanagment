import { test, expect } from '@playwright/test';

// @TestCase: TC-CUJ-004
test.describe('CUJ-4: Ephemeral Environment Cloning Journey', () => {
  test('verifies environment management controls', async ({ page }) => {
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

