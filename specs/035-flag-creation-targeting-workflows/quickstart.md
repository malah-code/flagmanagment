# Quickstart & Verification: Feature Flag Creation & Targeting Workflows

## Validation Steps via Puppeteer

### Scenario 1: Multi-Type Flag Creation (US1)
1. **Navigate**: Open `/projects/:projectId/flags`
2. **Action**: Click `+ New Feature Flag`.
3. **Fill Form**:
   - Key: `dark-mode-v2`
   - Name: `Dark Mode V2`
   - Type: `BOOLEAN`
4. **Submit**: Click `Create Flag`.
5. **Verify**: Flag row appears in the table with key `dark-mode-v2` and `Active` status.

---

### Scenario 2: Contextual Targeting Configuration (US2)
1. **Navigate**: Switch to `Development` environment view.
2. **Action**: Click `Targeting` button on `dark-mode-v2`.
3. **Configure Rule**:
   - Attribute: `country`
   - Operator: `EQUALS`
   - Value: `US`
   - Variation: `ON`
4. **Submit**: Save targeting rule.
5. **Verify**: State persists and success confirmation toast displays.

---

### Scenario 3: Emergency Kill Switch (US3)
1. **Action**: Click `Kill Switch` on an active flag in `Development`.
2. **Reason**: Enter `Production anomaly mitigation`.
3. **Confirm**: Click `Confirm Kill`.
4. **Verify**: Flag status immediately transitions to `INACTIVE (OFF)`.
