# Data Model: Enterprise SSO

## `users` (Updates)
- `auth_provider` (VARCHAR): Enum for the authentication method. E.g. "LOCAL", "OIDC", "SAML".
- `external_id` (VARCHAR): The IdP subject ID mapping (e.g. `sub` claim for OIDC, `NameID` for SAML). Nullable for local users.

## `service_accounts` (Existing)
- `id` (UUID): Primary key.
- `name` (VARCHAR): Display name.

## `service_account_keys` (Updates)
- Need to ensure robust hashing of tokens in `key_hash` or issue JWTs signed with RS256 without persisting the raw token.
- `token_prefix` (VARCHAR, Optional): Useful to show "fk_***abcd" in the UI.
