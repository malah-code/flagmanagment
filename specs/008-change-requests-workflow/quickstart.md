# Quickstart Validation: Change Requests Workflow

## Prerequisites
- The backend must be running.
- You must have an Admin token and two test user tokens (one Editor, one Release Manager).

## Scenario: End-to-End Change Request Lifecycle

1. **Protect the Environment**
   ```bash
   curl -X PATCH http://localhost:8080/environments/env-1 \
     -H "Authorization: Bearer $ADMIN_TOKEN" \
     -d '{"is_protected": true}'
   ```

2. **Propose a Change (as Editor)**
   ```bash
   curl -X PUT http://localhost:8080/flags/flag-1/environments/env-1 \
     -H "Authorization: Bearer $EDITOR_TOKEN" \
     -d '{"state": "ENABLED"}'
   ```
   *Expected Outcome*: Returns `202 Accepted` instead of `200 OK`. The response body contains the `ChangeRequest` object with `status: "Pending"`. Note the `id` of the change request.

3. **Verify Self-Approval Prevention**
   ```bash
   curl -X POST http://localhost:8080/change-requests/REQ_ID/approve \
     -H "Authorization: Bearer $EDITOR_TOKEN"
   ```
   *Expected Outcome*: Returns `403 Forbidden` because the author cannot approve their own request (and may not have the Release Manager role).

4. **Approve the Change (as Release Manager)**
   ```bash
   curl -X POST http://localhost:8080/change-requests/REQ_ID/approve \
     -H "Authorization: Bearer $RELEASE_MANAGER_TOKEN"
   ```
   *Expected Outcome*: Returns `200 OK`. The response shows the Change Request is now `Applied`.

5. **Verify the Change is Live**
   ```bash
   curl http://localhost:8080/flags/flag-1/environments/env-1
   ```
   *Expected Outcome*: The flag state is now `ENABLED`.
