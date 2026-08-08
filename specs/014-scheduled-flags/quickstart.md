# Quickstart: Validating Scheduled Flag Changes

**Feature**: `014-scheduled-flags`
**Phase**: 1 — Design & Contracts

---

## Prerequisites

- Docker Compose stack running (`make dev` or `docker compose up`)
- A valid JWT token for a `RELEASE_MANAGER` or `ADMIN` user
- A project, environment, and at least one feature flag already created
- `curl` and `jq` installed locally

### Environment Variables

```bash
export BASE_URL="http://localhost:8080/api/v1"
export AUTH_TOKEN="<your JWT token>"
export ENV_ID="<uuid of target environment>"
export FLAG_ID="<uuid of target feature flag>"
```

---

## Scenario 1: Schedule a Flag to Turn ON (US-1, AC-1)

### Step 1 — Confirm the flag is currently OFF

```bash
curl -s -H "Authorization: Bearer $AUTH_TOKEN" \
  "$BASE_URL/environments/$ENV_ID/flags/$FLAG_ID/state" | jq '.enabled'
# Expected: false
```

### Step 2 — Create a schedule 2 minutes in the future

```bash
SCHEDULED_FOR=$(date -u -d '+2 minutes' +"%Y-%m-%dT%H:%M:%SZ")
curl -s -X POST \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"target_type\":\"FLAG\",\"target_id\":\"$FLAG_ID\",\"action\":\"ENABLE\",\"scheduled_for\":\"$SCHEDULED_FOR\"}" \
  "$BASE_URL/environments/$ENV_ID/scheduled-changes" | jq .
# Expected: JSON with status="PENDING", action="ENABLE"
```

### Step 3 — Confirm visible on the flag detail listing

```bash
curl -s -H "Authorization: Bearer $AUTH_TOKEN" \
  "$BASE_URL/environments/$ENV_ID/scheduled-changes?status=PENDING" | jq '.data[0]'
# Expected: The PENDING record for your flag
```

### Step 4 — Wait for execution and verify

After the scheduled time + up to 30 seconds (scheduler poll interval):

```bash
curl -s -H "Authorization: Bearer $AUTH_TOKEN" \
  "$BASE_URL/environments/$ENV_ID/flags/$FLAG_ID/state" | jq '.enabled'
# Expected: true
```

### Step 5 — Confirm audit log entry

```bash
curl -s -H "Authorization: Bearer $AUTH_TOKEN" \
  "$BASE_URL/environments/$ENV_ID/audit-logs" | jq '[.data[] | select(.action=="SCHEDULED_EXECUTION")][0]'
# Expected: entry with action="SCHEDULED_EXECUTION", target_type="FEATURE_FLAG"
```

---

## Scenario 2: Cancel a Scheduled Change (US-1, AC-3)

```bash
# 1. Create a schedule 10 minutes from now
SCHEDULED_FOR=$(date -u -d '+10 minutes' +"%Y-%m-%dT%H:%M:%SZ")
SC_ID=$(curl -s -X POST \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"target_type\":\"FLAG\",\"target_id\":\"$FLAG_ID\",\"action\":\"DISABLE\",\"scheduled_for\":\"$SCHEDULED_FOR\"}" \
  "$BASE_URL/environments/$ENV_ID/scheduled-changes" | jq -r '.id')

# 2. Cancel it
curl -s -X DELETE \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  "$BASE_URL/scheduled-changes/$SC_ID" | jq '.status'
# Expected: "CANCELLED"

# 3. Flag state should remain unchanged after the original scheduled time
```

---

## Scenario 3: Reject Conflicting Schedule (FR-006)

```bash
# With a PENDING schedule already existing for $FLAG_ID, try to create another
SCHEDULED_FOR=$(date -u -d '+5 minutes' +"%Y-%m-%dT%H:%M:%SZ")
curl -s -X POST \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"target_type\":\"FLAG\",\"target_id\":\"$FLAG_ID\",\"action\":\"DISABLE\",\"scheduled_for\":\"$SCHEDULED_FOR\"}" \
  "$BASE_URL/environments/$ENV_ID/scheduled-changes" | jq .
# Expected: HTTP 409, error message about existing pending schedule
```

---

## Scenario 4: Schedule an Approved Change Request (US-2, AC-1)

```bash
export CR_ID="<uuid of an APPROVED change request>"

SCHEDULED_FOR=$(date -u -d '+3 minutes' +"%Y-%m-%dT%H:%M:%SZ")
curl -s -X POST \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"target_type\":\"CHANGE_REQUEST\",\"target_id\":\"$CR_ID\",\"action\":\"APPLY\",\"scheduled_for\":\"$SCHEDULED_FOR\"}" \
  "$BASE_URL/environments/$ENV_ID/scheduled-changes" | jq .
# Expected: JSON with status="PENDING", action="APPLY"

# After scheduled_for + up to 30s, verify change request status:
curl -s -H "Authorization: Bearer $AUTH_TOKEN" \
  "$BASE_URL/change-requests/$CR_ID" | jq '.status'
# Expected: "APPLIED"
```

---

## Scenario 5: Modify a Scheduled Time (FR-005)

```bash
NEW_TIME=$(date -u -d '+20 minutes' +"%Y-%m-%dT%H:%M:%SZ")
curl -s -X PATCH \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"scheduled_for\":\"$NEW_TIME\"}" \
  "$BASE_URL/scheduled-changes/$SC_ID" | jq '.scheduled_for'
# Expected: Updated UTC timestamp
```

---

## Expected Outcomes Summary

| Scenario              | Verification Point            | Expected Result        |
|-----------------------|-------------------------------|------------------------|
| Flag scheduled ON     | Flag state after scheduled time | `enabled: true`      |
| Cancel schedule       | Flag state unchanged          | No change to flag      |
| Conflict rejection    | HTTP status code              | 409 Conflict           |
| CR scheduled APPLY    | CR status after scheduled time | `APPLIED`             |
| Audit log             | Action field in audit log     | `SCHEDULED_EXECUTION`  |
| Modify schedule       | `scheduled_for` in response   | New UTC timestamp      |

---

## References

- Data Model: [data-model.md](./data-model.md)
- API Contract: [contracts/scheduled-changes.openapi.yaml](./contracts/scheduled-changes.openapi.yaml)
- Feature Spec: [spec.md](./spec.md)
