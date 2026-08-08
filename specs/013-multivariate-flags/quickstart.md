# Phase 1: Quickstart & Validation - Multivariate Flags

## Pre-requisites
1. The server must be running locally.
2. An environment and project must exist with API keys available.

## Validation Scenario: Create and evaluate a Multivariate Flag

1. **Create the flag**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/projects/{project_id}/flags \
     -H "Content-Type: application/json" \
     -d '{
       "key": "hero-banner-ab-test",
       "name": "Hero Banner A/B Test",
       "type": "MULTIVARIATE",
       "variations": [
         {"id": "var_a", "name": "Variant A", "value": "green"},
         {"id": "var_b", "name": "Variant B", "value": "red"}
       ]
     }'
   ```

2. **Configure Rollout (50/50)**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/environments/{env_id}/flags/hero-banner-ab-test/state \
     -H "Content-Type: application/json" \
     -d '{
       "isEnabled": true,
       "default_variation": "var_a",
       "rollout_rules": [
         {"variation_id": "var_a", "percentage": 5000},
         {"variation_id": "var_b", "percentage": 5000}
       ]
     }'
   ```

3. **Evaluate Flag via SDK API**:
   ```bash
   # User 1 evaluates (expected output consistently Variant A or B depending on hash)
   curl -X POST http://localhost:8080/api/v1/sdk/evaluate \
     -H "Authorization: Bearer {env_key}" \
     -H "Content-Type: application/json" \
     -d '{
       "flagKey": "hero-banner-ab-test",
       "context": {
         "userId": "user_123"
       }
     }'
   ```

## Expected Outcome
The exact same `userId` always receives the same result. A distribution of multiple different user IDs will split roughly evenly between `var_a` and `var_b`.
