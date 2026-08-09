INSERT INTO roles (id, name, description, permissions)
VALUES 
  ('00000000-0000-0000-0000-000000000001', 'ADMIN', 'Full Administrator', '["*"]'::jsonb),
  ('00000000-0000-0000-0000-000000000002', 'EDITOR', 'Can edit flags and environments', '["read", "write"]'::jsonb),
  ('00000000-0000-0000-0000-000000000003', 'VIEWER', 'Read-only access', '["read"]'::jsonb)
ON CONFLICT (name) DO NOTHING;

INSERT INTO users (id, email, password_hash, auth_provider)
VALUES 
  ('00000000-0000-0000-0000-000000000099', 'admin@example.com', '$2b$12$MQm0g8thOB2IGGkkVjoxOuU.pGqNGJgmTzrw0KmxXe7glyznodoq.', 'local')
ON CONFLICT (email) DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
VALUES 
  ('00000000-0000-0000-0000-000000000099', '00000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;
