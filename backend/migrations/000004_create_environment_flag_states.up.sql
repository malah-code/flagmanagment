CREATE TABLE IF NOT EXISTS environment_flag_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    feature_flag_id UUID NOT NULL REFERENCES feature_flags(id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    targeting_rules JSONB DEFAULT '{}'::jsonb,
    remote_config JSONB DEFAULT '{}'::jsonb,
    variations JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(environment_id, feature_flag_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_env_flag_states_lookup ON environment_flag_states(environment_id, feature_flag_id);
CREATE INDEX IF NOT EXISTS idx_env_flag_states_env_id ON environment_flag_states(environment_id);
