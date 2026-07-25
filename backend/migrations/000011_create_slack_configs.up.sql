CREATE TABLE IF NOT EXISTS slack_webhook_configs (
    id UUID PRIMARY KEY,
    environment_id UUID NOT NULL UNIQUE REFERENCES environments(id) ON DELETE CASCADE,
    webhook_url TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_slack_configs_env ON slack_webhook_configs(environment_id);
