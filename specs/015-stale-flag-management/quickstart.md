# Quickstart Validation Guide: Stale Flag Detection

**Feature**: `015-stale-flag-management`

## Prerequisites
- A running backend server connected to PostgreSQL and Redis.
- A project and environment created.
- At least one active feature flag.

## Scenario 1: Simulating Stale Flag Detection

1. **Backdate a Flag's State Change**:
   Directly update the database to simulate a flag that has been unchanged for 45 days.
   ```bash
   psql -U flagmanagement -d flagdb -c "UPDATE environment_flag_states SET last_state_change_at = NOW() - INTERVAL '45 days', boolean_state = true WHERE feature_flag_id = '<YOUR_FLAG_ID>';"
   ```

2. **Trigger the Stale Scanner**:
   Trigger the background job manually via the admin health/trigger endpoint.
   ```bash
   curl -X POST http://localhost:8080/api/v1/system/jobs/scan-stale-flags
   ```

3. **Verify Lifecycle State**:
   Fetch the flag from the API and verify the state is `STALE`.
   ```bash
   curl -X GET -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/v1/environments/<ENV_ID>/flags
   ```
   **Expected**: The flag object includes `"lifecycle_state": "STALE"`.

## Scenario 2: Archiving a Stale Flag

1. **Archive the Flag via API**:
   ```bash
   curl -X POST -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/v1/environments/<ENV_ID>/flags/<FLAG_ID>/lifecycle/archive
   ```

2. **Verify SDK Payload Exclusion**:
   Fetch the SDK streaming payload / initial ruleset.
   ```bash
   curl -X GET -H "Authorization: Bearer <SDK_CLIENT_TOKEN>" http://localhost:8080/sdk/v1/flags
   ```
   **Expected**: The archived flag is *missing* from the payload entirely, ensuring client payloads remain lightweight and free of technical debt.
