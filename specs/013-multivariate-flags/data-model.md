# Phase 1: Data Model - Multivariate Flags

## Feature Flag Enhancements
- `feature_flags` table already supports `Type` (`BOOLEAN` vs `MULTIVARIATE`).
- We need to define variations. These can be stored inside `feature_flags` as a JSONB column `variations`, or in a separate table. Given the JSONB heavy approach of FlagManagment, storing it in `feature_flags.variations` (JSONB) is optimal.

### `feature_flags` (Updated)
- Add `variations` (JSONB): An array of objects:
  ```json
  [
    {
      "id": "var_1",
      "name": "Red Button",
      "value": "red" // can be string, number, or json
    },
    {
      "id": "var_2",
      "name": "Blue Button",
      "value": "blue"
    }
  ]
  ```

## Environment Flag State Enhancements
- `environment_flag_states` already exists. We need to store the `percentage_rollout` allocations.
- Add `rollout_rules` (JSONB) to `environment_flag_states`, or update the existing `targeting_rules` payload to allow rollout splits.

### `environment_flag_states` (Updated)
- `default_variation`: (string) ID of the variation to serve by default (or when Identity is missing).
- `rollout_rules`: (JSONB) An array defining percentages for variations.
  ```json
  [
    {
      "variation_id": "var_1",
      "percentage": 3333 // representing 33.33% (out of 10000 basis points)
    },
    {
      "variation_id": "var_2",
      "percentage": 6667 // representing 66.67%
    }
  ]
  ```
