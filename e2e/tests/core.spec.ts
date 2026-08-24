import { test, expect } from '@playwright/test';

test.describe.serial('FlagManagment Core User Journeys', () => {
  let projectName = `Feature Project ${Date.now()}`;
  let envName = `Staging Env ${Date.now()}`;
  let flagKey = `checkout-v2-${Date.now()}`;

  // Use a single browser context to maintain auth state across all tests in this describe block
  test.beforeAll(async ({ browser }) => {
    // optional setup
  });

  test('Authentication & Login Flow', async ({ page }) => {
    await page.goto('/');
    if (page.url().includes('/login')) {
      await page.fill('input[type="email"]', 'admin@example.com');
      await page.fill('input[type="password"]', 'admin123');
      await page.click('button[type="submit"]');
    }
    await page.waitForURL('**/projects');
    await expect(page.locator('h1')).toContainText('Projects');
    // Save state so next tests can use it
    await page.context().storageState({ path: 'auth.json' });
  });

  // Only the tests AFTER login use the saved auth state. Declaring
  // test.use() at the outer describe level would apply it to the login
  // test itself, which fails because auth.json doesn't exist yet.
  test.describe('authenticated', () => {
    test.use({ storageState: 'auth.json' });

    async function openProject(page: import('@playwright/test').Page) {
      await page.goto('/projects');
      const targetCard = page.locator('div', { has: page.locator(`span:has-text("${projectName}")`) })
                             .filter({ hasText: 'Manage' }).last();
      await targetCard.locator('a:has-text("Manage")').click();
      await page.waitForURL('**/projects/*');
    }

    test('Projects Creation & Listing', async ({ page }) => {
      await page.goto('/projects');
      await page.click('button:has-text("New Project")');
      await page.waitForSelector('form');

      await page.fill('input[placeholder*="Name"], input[name="name"], form input:first-of-type', projectName);
      await page.fill('textarea, input[placeholder*="description"]', 'Project created for feature test');
      await page.click('button[type="submit"]:has-text("Create")');

      await page.waitForSelector('div.fixed.inset-0', { state: 'detached' });
      await expect(page.locator('body')).toContainText(projectName);
    });

    test('Project Workspace Navigation', async ({ page }) => {
      await page.goto('/projects');
      const targetCard = page.locator('div', { has: page.locator(`span:has-text("${projectName}")`) })
                             .filter({ hasText: 'Manage' }).last();

      await targetCard.locator('a:has-text("Manage")').click();
      await page.waitForURL('**/projects/*');
      await expect(page.locator('body')).toContainText(projectName);
    });

    test('Environments Management', async ({ page }) => {
      await openProject(page);

      // Environment management lives under "Environment Settings" in the sidebar
      await page.click('a:has-text("Environment Settings")');
      await page.waitForURL('**/projects/*/settings');

      await page.click('button:has-text("Add Environment")');
      await page.waitForSelector('form');
      await page.fill('form input[type="text"]', envName);
      await page.click('button[type="submit"]:has-text("Create Environment")');

      // After creation the dialog shows the generated SDK key; dismiss it
      await page.click('button:has-text("I have securely copied this key")');
      await page.waitForSelector('div.fixed.inset-0', { state: 'detached' });
      await expect(page.locator('body')).toContainText(envName);
    });

    test('Feature Flags Creation', async ({ page }) => {
      await openProject(page);

      // Flag definitions live under "All Flag Definitions" in the sidebar
      await page.click('a:has-text("All Flag Definitions")');
      await page.waitForURL('**/projects/*/flags');

      await page.click('button:has-text("New Feature Flag")');
      await page.waitForSelector('form');
      await page.fill('input[placeholder*="checkout"]', flagKey);
      await page.fill('textarea', 'Boolean flag created for testing');
      await page.click('button[type="submit"]:has-text("Create Flag")');
      await page.waitForSelector('div.fixed.inset-0', { state: 'detached' });
      await expect(page.locator('body')).toContainText(flagKey);
    });
  });
});
