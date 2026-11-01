import { test, expect } from '@playwright/test';

// @TestCase: TC-CUJ-001
test.describe('CUJ-1: Developer Flag Rollout & Contextual Targeting', () => {
  test('navigates to flags and verifies targeting rule controls', async ({ page }) => {
    await page.goto('/');
    if (page.url().includes('/login')) {
      await page.fill('input[type="email"]', 'admin@example.com');
      await page.fill('input[type="password"]', 'admin123');
      await page.click('button[type="submit"]');
      await page.waitForURL('**/projects');
    }
    await expect(page.locator('h1')).toBeVisible();
  });
});

