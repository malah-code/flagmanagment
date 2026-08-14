---
name: test-admin-ui
description: Systematically tests the local Admin UI page-by-page using the browser subagent, capturing screenshots and checking for errors.
---
# Admin UI Testing Skill

When the user invokes this skill, you must utilize the Browser Subagent to conduct an end-to-end visual and functional test of the admin dashboard.

## Execution Steps
1. **Navigate:** Open `http://localhost:3000` (or the URL provided by the user).
2. **Authenticate:** If a login screen is present, use the credentials provided by the user to log in. Wait for the dashboard to fully load.
3. **Map Navigation:** Analyze the DOM to identify the primary sidebar or top navigation menu. Extract all top-level page links.
4. **Systematic Click-Through:** For every link discovered in the navigation menu:
   - Click the link to navigate.
   - Wait for React/Next.js client-side rendering to finish (wait for data tables, charts, or lists to populate).
   - Capture a screenshot of the fully rendered page.
   - Check the browser console logs for any JavaScript errors, hydration mismatches, or API network failures.
5. **Report:** Return to the main chat and generate a comprehensive Markdown report containing the screenshots, page statuses, and any console errors found.

## Constraints
- **Destructive Actions:** NEVER click buttons labeled "Delete", "Remove", "Drop", or "Submit" inside form modules during the test. Read-only navigation only.
- **Patience:** Always wait at least 3 seconds after a route change before taking a screenshot to ensure API data has loaded.
