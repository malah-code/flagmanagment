CREATE TABLE scheduled_changes (
    id             UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    project_id     UUID        NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    environment_id UUID        NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    target_type    VARCHAR(20) NOT NULL CHECK (target_type IN ('FLAG', 'CHANGE_REQUEST')),
    target_id      UUID        NOT NULL,
    action         VARCHAR(20) NOT NULL CHECK (action IN ('ENABLE', 'DISABLE', 'APPLY')),
    scheduled_for  TIMESTAMPTZ NOT NULL,
    status         VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'EXECUTED', 'CANCELLED')),
    created_by     UUID        NOT NULL REFERENCES users(id),
    executed_at    TIMESTAMPTZ,
    cancelled_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_scheduled_changes_pending_flag
  ON scheduled_changes (target_id)
  WHERE status = 'PENDING' AND target_type = 'FLAG';

CREATE INDEX idx_scheduled_changes_due
  ON scheduled_changes (scheduled_for, status)
  WHERE status = 'PENDING';

CREATE INDEX idx_scheduled_changes_env
  ON scheduled_changes (environment_id, status);
