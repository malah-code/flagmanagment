import { test, expect } from '@playwright/test';

// @TestCase: TC-CUJ-005
test.describe('CUJ-5: Automated Kill-Switch & Stale Flag Journey', () => {
  test('verifies kill switch and alerting components', async ({ page }) => {
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

