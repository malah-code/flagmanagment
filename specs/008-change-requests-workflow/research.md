# Research: Change Requests Workflow

## 1. Database Representation of Change Requests

**Decision**: Create a `ChangeRequest` table in PostgreSQL.
**Rationale**: Needs to store environment ID, flag ID, author ID, proposed state, current state, status, and approval details. Utilizing GORM for this. The states (`Pending`, `Approved`, `Applied`, `Rejected`) will be string enums.
**Alternatives considered**: Storing change requests as a JSON array on the environment or flag. Rejected because it would make querying and auditing approvals much more complex and less relational.

## 2. Intercepting Flag Mutations

**Decision**: Modify the existing flag update API/Service logic. When an update is requested for a flag, the service will check if the target Environment `is_protected`. If yes, it creates a `ChangeRequest` and returns a specific status (e.g., 202 Accepted) instead of 200 OK. If no, it updates immediately.
**Rationale**: Keeps the logic centralized in the service layer.
**Alternatives considered**: Using a separate API endpoint for creating change requests. Rejected because it shifts the burden to the client to know if an environment is protected.

## 3. Visual JSON Diffing

**Decision**: On the frontend, use a library like `react-diff-viewer` (or build a simple custom diff component if dependencies need to be minimized) to show the `current_state` vs `proposed_state` side-by-side.
**Rationale**: The backend will provide the raw JSON for both states. The frontend is best suited to visually render the differences for the Release Manager.
**Alternatives considered**: Backend generates an HTML diff. Rejected because it breaks the API-first JSON contract paradigm.

## 4. Role Enforcement (Release Manager)

**Decision**: Enhance the existing `RBACMiddleware` or add a specific check in the Change Request approval service to ensure the user has the `Release Manager` role (or `Admin`). It must also verify that `current_user.id != change_request.author_id`.
**Rationale**: Aligns with the Constitution's "Governance by Default" and the specific requirement to prevent self-approval.
