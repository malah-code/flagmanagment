import { Page, expect } from '@playwright/test';

export async function loginUser(page: Page) {
  if (page.url().includes('/projects')) {
    const heading = page.locator('h1:has-text("Projects")');
    if (await heading.isVisible()) {
      return;
    }
  }

  await page.goto('/login');
  await page.waitForSelector('input[type="email"]', { timeout: 15000 });
  await page.fill('input[type="email"]', 'admin@example.com');
  await page.fill('input[type="password"]', 'admin123');
  await page.click('button[type="submit"]');
  await page.waitForURL('**/projects', { timeout: 15000 });
  await expect(page.locator('h1')).toContainText('Projects');
}

