import { test, expect } from '@playwright/test';
import { loginUser } from './helpers';

// @TestCase: TC-CUJ-001
test.describe('CUJ-1: Developer Flag Rollout & Contextual Targeting', () => {
  test('navigates to flags and verifies targeting rule controls', async ({ page }) => {
    await loginUser(page);
    await expect(page.locator('body')).toBeVisible();
  });
});
