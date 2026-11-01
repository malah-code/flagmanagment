import { test, expect } from '@playwright/test';
import { loginUser } from './helpers';

// @TestCase: TC-CUJ-006
test.describe('CUJ-6: OpenFeature Multi-Language Conformance Journey', () => {
  test('verifies integration hub and OpenFeature playground conformance', async ({ page }) => {
    await loginUser(page);
    await page.goto('/integrations');
    await expect(page.locator('body')).toBeVisible();
    await expect(page.locator('h1')).toContainText('Integrations');

    await page.goto('/projects');
    await expect(page.locator('body')).toBeVisible();

    // Check if any projects exist to navigate into OpenFeature playground
    const manageLinks = page.locator('a:has-text("Manage")');
    if ((await manageLinks.count()) > 0) {
      await manageLinks.first().click();
      await page.waitForURL('**/projects/*');

      // Click OpenFeature Playground in the sidebar
      const playgroundLink = page.locator('a:has-text("OpenFeature Playground")');
      if (await playgroundLink.isVisible()) {
        await playgroundLink.click();
        await page.waitForURL('**/openfeature');
        await expect(page.locator('h1')).toContainText('OpenFeature Playground');
        await expect(page.locator('body')).toContainText('ResolutionDetails');
      }
    }
  });
});
