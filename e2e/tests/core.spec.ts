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

  test.use({ storageState: 'auth.json' });

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
    await page.goto('/projects');
    const targetCard = page.locator('div', { has: page.locator(`span:has-text("${projectName}")`) })
                           .filter({ hasText: 'Manage' }).last();
    await targetCard.locator('a:has-text("Manage")').click();
    await page.waitForURL('**/projects/*');

    await page.click('button:has-text("Environments"), a:has-text("Environments")');
    const addEnvBtn = page.locator('button:has-text("Add Environment"), button:has-text("New Environment")');
    
    if (await addEnvBtn.isVisible()) {
      await addEnvBtn.click();
      await page.fill('form input:first-of-type', envName);
      await page.click('button[type="submit"]:has-text("Create"), button[type="submit"]:has-text("Add")');
      await page.waitForSelector('div.fixed.inset-0', { state: 'detached' });
      await expect(page.locator('body')).toContainText(envName);
    }
  });

  test('Feature Flags Creation', async ({ page }) => {
    await page.goto('/projects');
    const targetCard = page.locator('div', { has: page.locator(`span:has-text("${projectName}")`) })
                           .filter({ hasText: 'Manage' }).last();
    await targetCard.locator('a:has-text("Manage")').click();
    await page.waitForURL('**/projects/*');

    await page.click('button:has-text("Feature Flags"), a:has-text("Feature Flags")');
    const createFlagBtn = page.locator('button:has-text("New Flag"), button:has-text("Create Flag")');
    
    if (await createFlagBtn.isVisible()) {
      await createFlagBtn.click();
      await page.waitForSelector('form');
      await page.fill('input[placeholder*="key"], input[name="key"]', flagKey);
      await page.fill('textarea, input[placeholder*="description"]', 'Boolean flag created for testing');
      await page.click('button[type="submit"]:has-text("Create")');
      await page.waitForSelector('div.fixed.inset-0', { state: 'detached' });
      await expect(page.locator('body')).toContainText(flagKey);
    }
  });
});
