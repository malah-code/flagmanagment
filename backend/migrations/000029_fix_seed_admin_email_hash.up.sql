-- Fix the seeded admin user's email_hash.
--
-- Migration 000027 backfilled email_hash with the plaintext email
-- ("email_hash = email" as a temporary fallback), but the repository
-- looks up users by the SHA-256 hash of the email (HashToken in
-- internal/crypto). The seeded admin@example.com could therefore never
-- be found, and login always returned 401 "Invalid email or password".
--
-- This sets the seed row's email_hash to sha256("admin@example.com")
-- so GetByEmail matches. Other rows are unaffected (a fresh DB has no
-- other users; real deployments migrate PII data out-of-band).
UPDATE users
SET email_hash = '258d8dc916db8cea2cafb6c3cd0cb0246efe061421dbd83ec3a350428cabda4f'
WHERE id = '00000000-0000-0000-0000-000000000099'
  AND email = 'admin@example.com';
