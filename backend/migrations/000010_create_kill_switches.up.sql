CREATE TABLE IF NOT EXISTS kill_switches (
    id UUID PRIMARY KEY,
    flag_id UUID NOT NULL REFERENCES feature_flags(id) ON DELETE CASCADE,
    environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    alert_identifier VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL DEFAULT 'DISABLE',
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT uk_kill_switch_env_alert UNIQUE (environment_id, alert_identifier, flag_id)
);
