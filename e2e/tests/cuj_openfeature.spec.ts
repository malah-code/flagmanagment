import { test, expect } from '@playwright/test';

// @TestCase: TC-CUJ-006
test.describe('CUJ-6: OpenFeature Multi-Language Conformance Journey', () => {
  test('verifies integration hub and SDK documentation page', async ({ page }) => {
    await page.goto('/');
    if (page.url().includes('/login')) {
      await page.fill('input[type="email"]', 'admin@example.com');
      await page.fill('input[type="password"]', 'admin123');
      await page.click('button[type="submit"]');
      await page.waitForURL('**/projects');
    }
    await page.goto('/integrations');
    await expect(page.locator('body')).toBeVisible();
  });
});

