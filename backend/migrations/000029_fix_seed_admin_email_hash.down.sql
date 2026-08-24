-- Revert: restore the plaintext fallback backfill from migration 000027.
-- Only touches the seed admin row; safe to run on a fresh DB.
UPDATE users
SET email_hash = email
WHERE id = '00000000-0000-0000-0000-000000000099'
  AND email = 'admin@example.com';
