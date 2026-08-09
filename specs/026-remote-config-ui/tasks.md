# Implementation Tasks: Remote Configuration Payload UI

## T001: Install Dependencies
- [ ] Run `npm install @monaco-editor/react` in the frontend directory to provide the JSON code editor.

## T002: Update Flag Type Selector
- [ ] In `frontend/src/components/flags/CreateFlagDialog.tsx`, update the conditional rendering that reveals the Variations section from `type === 'MULTIVARIATE'` to `type === 'MULTIVARIATE' || type === 'JSON'`.

## T003: Initialize JSON Variations
- [ ] In `frontend/src/components/flags/CreateFlagDialog.tsx`, when `type` changes to `'JSON'`, automatically reset the `variations` array to valid JSON string templates (e.g., `"{}"` or `"{\n  \n}"`) so the editor has a valid base.

## T004: Integrate Monaco Editor
- [ ] In `frontend/src/components/flags/CreateFlagDialog.tsx`, within the Variations loop, check if `type === 'JSON'`.
- [ ] If true, render the Monaco `<Editor>` component for the variation `value` field. Configure it with `language="json"`, a reasonable height (e.g., `120px`), and a minimap disabled for compactness.
- [ ] Fall back to the standard text `<input>` if `type === 'MULTIVARIATE'`.

## T005: Add JSON Validation on Submit
- [ ] In the `handleSubmit` function of `CreateFlagDialog.tsx`, add a validation pass for `type === 'JSON'`.
- [ ] Iterate through all variations and attempt `JSON.parse(v.value)`.
- [ ] If parsing fails, set the form `error` state (e.g., `"Invalid JSON payload in [Variation Name]"`) and abort submission.

## T006: Map Parsed JSON to API Request
- [ ] In the `createMutation.mutateAsync` payload mapping within `handleSubmit`, conditionally map the variations array.
- [ ] For `type === 'JSON'`, map the `value` field to `JSON.parse(v.value)` so the API receives a proper JSON object instead of a stringified payload.

## T007: Update Types & Build
- [ ] Ensure `frontend/src/types/index.ts` is fully compatible with these changes.
- [ ] Run `npm run build` in the frontend directory to verify TypeScript compilation and lack of structural errors.
