# Phase 0: Research

## Decisions

- **Decision:** AES-256-GCM encryption for SMTP passwords.
  - **Rationale:** Standard, secure authenticated encryption supported natively by Go's `crypto/cipher` and `crypto/aes`. It ensures that SMTP credentials are not stored in plaintext in PostgreSQL.
  - **Alternatives considered:** Hashing (not viable for SMTP authentication since the plaintext is needed to connect), external secret managers like HashiCorp Vault (overkill for simple self-hosted deployments).

- **Decision:** Cryptographically secure random tokens (SHA-256 hashed in DB) for invitations.
  - **Rationale:** Need a secure, time-bound way to identify the invited user when they click the registration link. Storing the hash prevents compromise if the DB is leaked.
  - **Alternatives considered:** Simple database IDs (insecure, predictable), JWTs (unnecessary complexity for a one-time use token).
