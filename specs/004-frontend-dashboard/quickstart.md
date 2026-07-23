# Quickstart Validation Guide: Frontend Dashboard

This guide walks you through verifying the frontend dashboard locally.

## Prerequisites
- Node.js installed
- The Go backend must be running on `http://localhost:8080`.

## Setup
1. Navigate to the `frontend` directory.
2. Install dependencies: `npm install`
3. If not already present, the frontend connects to `/api/v1` which should be proxied to the backend via Vite's proxy settings.
4. Run the development server: `npm run dev`

## Validation Scenarios

### Scenario 1: Project Creation
1. Open the dashboard at `http://localhost:5173`.
2. Click "New Project".
3. Enter `test-project-1` as the key and `Test Project 1` as the name.
4. Submit the form.
5. **Expected**: The new project appears in the projects list and you are navigated to its details.

### Scenario 2: Environment Creation
1. Navigate to `Test Project 1`.
2. Click "New Environment".
3. Enter `production` as the key and `Production` as the name.
4. Submit.
5. **Expected**: A modal or alert displays the securely generated API key for `Production`. Copy it. It should not be visible again. The environment appears in the environments list.

### Scenario 3: Feature Flag Management
1. Inside `Test Project 1`, navigate to the "Feature Flags" tab.
2. Click "New Flag".
3. Enter `new-checkout-flow` as the key and select type `boolean`.
4. Submit.
5. **Expected**: The flag appears in the flags list.

### Scenario 4: Flag State Configuration
1. Select the `Production` environment.
2. Go to its "Flag States" section.
3. Locate `new-checkout-flow`.
4. Toggle it from `Disabled` to `Enabled`.
5. **Expected**: The toggle visibly changes to the ON position and the change persists when refreshing the page.
