# Data Model: Change Requests Workflow

## 1. Environment (Update)
- `id` (UUID, Primary Key)
- `project_id` (UUID, Foreign Key)
- `name` (String)
- `is_protected` (Boolean, default: false) - **[NEW FIELD]**

## 2. ChangeRequest (New Entity)
- `id` (UUID, Primary Key)
- `environment_id` (UUID, Foreign Key to Environment)
- `flag_id` (UUID, Foreign Key to Flag)
- `author_id` (UUID, Foreign Key to User)
- `status` (Enum: `Pending`, `Approved`, `Applied`, `Rejected`)
- `proposed_state` (JSONB) - The proposed flag configuration (state, targeting rules)
- `current_state` (JSONB) - The flag configuration at the time the request was created
- `created_at` (Timestamp)
- `updated_at` (Timestamp)

## 3. Approval (New Entity or Embed in ChangeRequest)
Since a Change Request can have one terminal approval/rejection decision, this can be embedded in the `ChangeRequest` table to reduce complexity, or kept separate. We will embed it:
- `reviewer_id` (UUID, Foreign Key to User, nullable)
- `review_comment` (String, nullable)
- `reviewed_at` (Timestamp, nullable)

## State Transitions
1. `Pending` -> `Approved` (Upon Approval by Release Manager) -> `Applied` (Synchronous auto-transition immediately after applying to live flag).
2. `Pending` -> `Rejected` (Upon Rejection by Release Manager).
