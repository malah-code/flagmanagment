# Implementation Plan: Remote Configuration Payload UI

This plan outlines the steps to introduce a fully-fledged JSON editor into the frontend for authoring and managing Remote Config feature flags, completing the final missing MVP capability.

## 1. Dependency Integration
- **Package**: `@monaco-editor/react`
- **Action**: Install via `npm install @monaco-editor/react` in the `/frontend` directory.
- **Why Monaco**: It provides native JSON syntax highlighting, structural bracket matching, auto-formatting, and live linting out of the box.

## 2. Frontend Updates

### `frontend/src/components/flags/CreateFlagDialog.tsx`
- **Render Variations for JSON**: Update the conditional rendering `type === 'MULTIVARIATE'` to `type === 'MULTIVARIATE' || type === 'JSON'`.
- **Default Variation Initialization**: When switching to `type === 'JSON'`, initialize the variations state to valid default JSON (e.g., `"{}"` as a string, because the Monaco editor works with strings).
- **Monaco Editor Component**: 
  - If `type === 'JSON'`, render the `<Editor />` component for the value field instead of a standard `<input>`.
  - Set language to `json` and height to `120px` (or resizable).
- **Validation Before Submit**: 
  - On form submit (`handleSubmit`), if `type === 'JSON'`, iterate through the `variations` array.
  - Attempt `JSON.parse(v.value)`. If any fail, intercept the submission and set the `error` state (e.g., `"Invalid JSON in Variation A"`).
  - Map the variations to pass the *parsed* JSON object to `createMutation` rather than the string literal (e.g., `value: JSON.parse(v.value)`), so the API receives an object.

### `frontend/src/types/index.ts`
- Ensure no strict typing blocks `value: unknown` (already correctly typed).
- Ensure `FlagType` includes `JSON` (already correctly typed).

### `frontend/src/components/ChangeRequestDiff.tsx`
- No code changes required! `react-diff-viewer-continued` already serializes `proposedChanges` using `JSON.stringify(..., null, 2)` internally. JSON payloads created by the new editor will natively render structurally in this component when requested via Change Request flows.

## Constitution Compliance Matrix

| Principle | Adherence Check | Notes |
| :--- | :--- | :--- |
| **I. API-First Contract Design** | ✅ PASS | No backend API changes needed. Existing contract supports arbitrary JSON payload structures via `JSONB`. |
| **II. Environment Isolation** | ✅ PASS | JSON variations are scoped exactly like standard string/multivariate flags. |
| **III. Governance by Default** | ✅ PASS | Complex JSON configurations will natively display in existing Change Request and Audit Diff viewers. |
| **IV. Local Evaluation Performance** | ✅ PASS | No impact to evaluation speed. SDKs serialize the JSON locally. |
| **V. Test-First Quality Gates** | ✅ PASS | JSON parse testing inherently provided by the compiler/browser engine. |
| **VI. OpenFeature Interoperability**| ✅ PASS | OpenFeature `Object` flag evaluation natively maps to this capability. |
| **VII. PII Protection & Compliance**| ✅ PASS | Follows standard data rules. |
| **VIII. Cloud-Native Portability** | ✅ PASS | Frontend dependency addition. No backend/deployment changes. |

## Open Questions for User
- The Monaco editor component requires downloading some web workers to the client. Are you okay with using `@monaco-editor/react` (which handles this gracefully via CDN by default), or do you prefer a lighter-weight editor like `textarea` with manual validation for the MVP?
