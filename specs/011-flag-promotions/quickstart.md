# Quickstart: Flag Promotions Validation

1. Configure two environments: Dev (unprotected) and Prod (protected).
2. Set Flag A to ENABLED in Dev.
3. Call `POST /api/v1/projects/{projectId}/flags/{flagId}/promote` with `source_env_id=Dev` and `target_env_id=Prod`.
4. Verify response is a created Change Request for Prod.
