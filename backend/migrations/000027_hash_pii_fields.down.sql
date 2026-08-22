ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_hash_key;
ALTER TABLE invitations DROP CONSTRAINT IF EXISTS invitations_email_hash_key;

DROP INDEX IF EXISTS idx_users_email_hash;
CREATE INDEX idx_users_email ON users(email);

DROP INDEX IF EXISTS idx_invitations_email_hash;
CREATE INDEX idx_invitations_email ON invitations(email);

ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
ALTER TABLE invitations ADD CONSTRAINT invitations_email_key UNIQUE (email);

ALTER TABLE users DROP COLUMN IF EXISTS email_hash;
ALTER TABLE invitations DROP COLUMN IF EXISTS email_hash;
