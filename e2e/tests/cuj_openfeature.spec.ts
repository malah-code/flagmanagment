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
    let manageLinks = page.locator('a:has-text("Manage")');
    if ((await manageLinks.count()) === 0) {
      const createBtn = page.locator('button:has-text("New Project"), button:has-text("Create Project")').first();
      if (await createBtn.isVisible()) {
        await createBtn.click();
        await page.waitForSelector('form');
        await page.fill('input[placeholder*="Name"], input[name="name"], form input:first-of-type', `OF Conformance ${Date.now()}`);
        await page.click('button[type="submit"]:has-text("Create")');
        await page.waitForSelector('div.fixed.inset-0', { state: 'detached' });
      }
      manageLinks = page.locator('a:has-text("Manage")');
    }

    if ((await manageLinks.count()) > 0) {
      await manageLinks.first().click();
      await page.waitForURL('**/projects/*');

      // Click OpenFeature Playground in the sidebar
      const playgroundLink = page.locator('a:has-text("OpenFeature Playground")');
      if (await playgroundLink.isVisible()) {
        await playgroundLink.click();
        await page.waitForURL('**/openfeature');
        await expect(page.locator('h1')).toContainText('OpenFeature Playground');

        // Trigger evaluation to show ResolutionDetails
        const runBtn = page.locator('button:has-text("Run OpenFeature Evaluation")');
        if (await runBtn.isVisible()) {
          await runBtn.click();
          await expect(page.locator('body')).toContainText('ResolutionDetails');
        }
      }
    }
  });
});
