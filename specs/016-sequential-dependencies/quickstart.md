# Quickstart: Sequential Dependencies

This guide provides steps to manually validate that sequential dependencies between feature flags function correctly end-to-end.

## Prerequisites
- A running instance of the FlagManagment backend and frontend.
- A configured project and environment with at least one active API token.
- cURL installed.

## Validation Scenarios

### Scenario 1: API Prevents Circular Dependencies

1. **Create Flag A**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/projects/{projectId}/flags \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <Admin_Token>" \
     -d '{"key": "flag-a", "type": "boolean", "name": "Flag A"}'
   ```
   *(Note the returned ID for Flag A)*

2. **Create Flag B depending on Flag A**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/projects/{projectId}/flags \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <Admin_Token>" \
     -d '{"key": "flag-b", "type": "boolean", "name": "Flag B", "parent_flag_id": "<ID_OF_FLAG_A>"}'
   ```
   *(Note the returned ID for Flag B)*

3. **Attempt to update Flag A to depend on Flag B**:
   ```bash
   curl -X PUT http://localhost:8080/api/v1/projects/{projectId}/flags/<ID_OF_FLAG_A> \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <Admin_Token>" \
     -d '{"key": "flag-a", "type": "boolean", "name": "Flag A", "parent_flag_id": "<ID_OF_FLAG_B>"}'
   ```
   **Expected Outcome:** The API must return a `400 Bad Request` with an error message explicitly indicating that a circular dependency was detected.

---

### Scenario 2: SDK Short-Circuit Evaluation

1. **Configure Environment State**:
   - Using the Dashboard UI, turn **Flag A** to `OFF`.
   - Turn **Flag B** to `ON`.
   - Ensure Flag B has Flag A configured as its parent.

2. **Evaluate Flag B via SDK (or Evaluation Endpoint)**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/sdk/evaluate \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <Environment_Token>" \
     -d '{"identity": "user-123"}'
   ```
   **Expected Outcome:**
   - The evaluation response for `flag-a` should be `false`.
   - The evaluation response for `flag-b` must also be `false` (the fallback state), even though its environment state is explicitly set to `ON`. Because Flag A is `OFF`, Flag B is never fully evaluated and safely defaults to `false`.
