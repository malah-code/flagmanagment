# API Contracts: Feature Flag Creation & Targeting Workflows

## 1. Feature Flag Management API

### `POST /api/v1/projects/{projectId}/flags`
Creates a new feature flag and initializes default state across all project environments.

**Request Body**:
```json
{
  "key": "enable-smart-checkout",
  "name": "Enable Smart Checkout",
  "description": "Dynamic checkout funnel with Apple Pay and Stripe Elements",
  "type": "BOOLEAN",
  "tags": ["checkout", "v2"],
  "variations": [
    { "id": "on", "name": "Enabled", "value": true },
    { "id": "off", "name": "Disabled", "value": false }
  ]
}
```

**Response (201 Created)**:
```json
{
  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "project_id": "80840e67-78eb-46d9-8335-ece3c11e79c6",
  "key": "enable-smart-checkout",
  "name": "Enable Smart Checkout",
  "description": "Dynamic checkout funnel with Apple Pay and Stripe Elements",
  "type": "BOOLEAN",
  "tags": ["checkout", "v2"],
  "created_at": "2026-08-15T02:00:00Z"
}
```

---

## 2. Environment Targeting & State API

### `PUT /api/v1/projects/{projectId}/environments/{envId}/flags/{flagId}/state`
Toggles or updates the flag configuration in an environment.

**Request Body**:
```json
{
  "is_enabled": true,
  "default_variation": "on",
  "off_variation": "off",
  "rules": [
    {
      "id": "rule_beta_users",
      "name": "Beta User Segment",
      "conditions": [
        {
          "attribute": "beta_tester",
          "operator": "EQUALS",
          "values": ["true"]
        }
      ],
      "variation": "on"
    }
  ],
  "rollout_weights": {
    "on": 100,
    "off": 0
  }
}
```

**Response (200 OK)**:
```json
{
  "id": "b1a7d65e-2f34-4b51-9e73-04859a721c5f",
  "environment_id": "1ccc6a07-50d5-468b-b578-cdf9644f933c",
  "flag_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "is_enabled": true,
  "rules": [...],
  "rollout_weights": { "on": 100, "off": 0 },
  "updated_at": "2026-08-15T02:05:00Z"
}
```

---

## 3. Emergency Kill Switch API

### `POST /api/v1/projects/{projectId}/environments/{envId}/flags/{flagId}/kill`
Instantly triggers emergency deactivation for a flag in the selected environment.

**Request Body**:
```json
{
  "reason": "Payment service degraded downstream"
}
```

**Response (200 OK)**:
```json
{
  "status": "killed",
  "is_enabled": false,
  "killed_at": "2026-08-15T02:10:00Z"
}
```
