# Quickstart Validation: Competitor UX Parity

## Prerequisites
- Frontend development server running (`npm run dev` in `frontend/`)
- Backend server running
- Seed data with multiple environments and flags with tags

## Validation Steps

1. **Sidebar Navigation**:
   - Open the dashboard. Verify the left sidebar is visible and displays environments.
   - Click different environments in the sidebar and ensure the main content area updates to that environment's flags.

2. **Flag Toggling**:
   - In the flag list for an environment, locate the boolean toggle switch on a row.
   - Click it and verify it immediately changes state, showing a loading spinner momentarily, and a success toast appears.

3. **Flag Creation**:
   - Click "Create Flag".
   - Verify the "Initial Value" field is present and optional.
   - Create a flag with an initial value and verify it defaults to that value.

4. **Tag Filtering**:
   - Locate the tag filter input above the flag list.
   - Select a tag and verify the list is filtered to only show flags with that tag.
