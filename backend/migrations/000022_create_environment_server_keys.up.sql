CREATE TABLE IF NOT EXISTS environment_server_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    key_hash VARCHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_env_server_keys_hash ON environment_server_keys(key_hash);
CREATE INDEX IF NOT EXISTS idx_env_server_keys_env_id ON environment_server_keys(environment_id);
