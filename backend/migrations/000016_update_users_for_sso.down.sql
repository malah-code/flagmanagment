DROP INDEX IF EXISTS idx_users_sso_identity;
ALTER TABLE users DROP COLUMN external_id;
ALTER TABLE users DROP COLUMN auth_provider;
-- Cannot safely re-add NOT NULL to password_hash if SSO users exist
