CREATE TABLE IF NOT EXISTS feature_flags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('BOOLEAN', 'MULTIVARIATE', 'JSON')),
    parent_flag_id UUID REFERENCES feature_flags(id) ON DELETE SET NULL,
    last_evaluated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(project_id, key),
    CHECK (parent_flag_id != id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_feature_flags_project_key ON feature_flags(project_id, key);
