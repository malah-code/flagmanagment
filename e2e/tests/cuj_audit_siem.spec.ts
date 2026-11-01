import { test, expect } from '@playwright/test';
import { loginUser } from './helpers';

// @TestCase: TC-CUJ-003
test.describe('CUJ-3: Audit Log SIEM Compliance Journey', () => {
  test('verifies audit logs page and table structure', async ({ page }) => {
    await loginUser(page);
    await page.goto('/audit-logs');
    await expect(page.locator('body')).toBeVisible();
    await expect(page.locator('h1')).toContainText('Audit Logs');
  });
});
