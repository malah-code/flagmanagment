DO $$ BEGIN
    CREATE TYPE flag_lifecycle_state AS ENUM ('ACTIVE', 'STALE', 'DEPRECATED', 'ARCHIVED');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

ALTER TABLE environment_flag_states 
    ADD COLUMN IF NOT EXISTS lifecycle_state flag_lifecycle_state NOT NULL DEFAULT 'ACTIVE',
    ADD COLUMN IF NOT EXISTS last_evaluated_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS last_state_change_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_env_flag_states_lifecycle ON environment_flag_states (environment_id, lifecycle_state);
CREATE INDEX IF NOT EXISTS idx_env_flag_states_staleness ON environment_flag_states (environment_id, last_evaluated_at, last_state_change_at);

CREATE TABLE IF NOT EXISTS stale_flag_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    environment_id UUID REFERENCES environments(id) ON DELETE CASCADE,
    stale_after_days INT NOT NULL DEFAULT 30,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_stale_policy_proj_env ON stale_flag_policies (project_id, COALESCE(environment_id, '00000000-0000-0000-0000-000000000000'));
