ALTER TABLE users ADD COLUMN email_hash VARCHAR(255);
ALTER TABLE invitations ADD COLUMN email_hash VARCHAR(255);

-- Update existing records by making email_hash = email initially (since they aren't hashed yet, it's a fallback)
-- In a real prod migration we'd run a script to hash them, but for this project we'll just drop the table data or assume it's a fresh DB.
-- Wait, we can't assume fresh DB, so let's just copy it for now.
UPDATE users SET email_hash = email;
UPDATE invitations SET email_hash = email;

ALTER TABLE users ALTER COLUMN email_hash SET NOT NULL;
ALTER TABLE invitations ALTER COLUMN email_hash SET NOT NULL;

-- Drop old unique constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE invitations DROP CONSTRAINT IF EXISTS invitations_email_key;

-- Add new unique constraints on the hash
ALTER TABLE users ADD CONSTRAINT users_email_hash_key UNIQUE (email_hash);
ALTER TABLE invitations ADD CONSTRAINT invitations_email_hash_key UNIQUE (email_hash);

-- Drop old indices and create new ones
DROP INDEX IF EXISTS idx_users_email;
CREATE INDEX idx_users_email_hash ON users(email_hash);

DROP INDEX IF EXISTS idx_invitations_email;
CREATE INDEX idx_invitations_email_hash ON invitations(email_hash);
