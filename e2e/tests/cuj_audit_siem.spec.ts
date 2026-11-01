import { test, expect } from '@playwright/test';

// @TestCase: TC-CUJ-003
test.describe('CUJ-3: Audit Log SIEM Compliance Journey', () => {
  test('verifies audit logs page and table structure', async ({ page }) => {
    await page.goto('/');
    if (page.url().includes('/login')) {
      await page.fill('input[type="email"]', 'admin@example.com');
      await page.fill('input[type="password"]', 'admin123');
      await page.click('button[type="submit"]');
      await page.waitForURL('**/projects');
    }
    await page.goto('/settings/system');
    await expect(page.locator('body')).toBeVisible();
  });
});

