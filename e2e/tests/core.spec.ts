import { test, expect, Page } from '@playwright/test';

test.describe.serial('FlagManagment Core User Journeys', () => {
  let projectName = `Feature Project ${Date.now()}`;
  let envName = `Staging Env ${Date.now()}`;
  let flagKey = `checkout-v2-${Date.now()}`;

  async function ensureLoggedIn(page: Page) {
    await page.goto('/');
    if (page.url().includes('/login')) {
      await page.fill('input[type="email"]', 'admin@example.com');
      await page.fill('input[type="password"]', 'admin123');
      await page.click('button[type="submit"]');
    }
    await page.waitForURL('**/projects');
    await expect(page.locator('h1')).toContainText('Projects');
  }

  async function openProject(page: Page) {
    await page.goto('/projects');
    const targetCard = page.locator('div', { has: page.locator(`span:has-text("${projectName}")`) })
                           .filter({ hasText: 'Manage' }).last();
    await targetCard.locator('a:has-text("Manage")').click();
    await page.waitForURL('**/projects/*');
  }

  test('Authentication & Login Flow', async ({ page }) => {
    await ensureLoggedIn(page);
  });

  test('Projects Creation & Listing', async ({ page }) => {
    await ensureLoggedIn(page);
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
    await ensureLoggedIn(page);
    await openProject(page);
    await expect(page.locator('body')).toContainText(projectName);
  });

  test('Environments Management', async ({ page }) => {
    await ensureLoggedIn(page);
    await openProject(page);

    await page.click('a:has-text("Environment Settings")');
    await page.waitForURL('**/projects/*/settings');

    await page.click('button:has-text("Add Environment")');
    await page.waitForSelector('form');
    await page.fill('form input[type="text"]', envName);
    await page.click('button[type="submit"]:has-text("Create Environment")');

    await page.click('button:has-text("I have securely copied this key")');
    await page.waitForSelector('div.fixed.inset-0', { state: 'detached' });
    await expect(page.locator('body')).toContainText(envName);
  });

  test('Feature Flags Creation', async ({ page }) => {
    await ensureLoggedIn(page);
    await openProject(page);

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
