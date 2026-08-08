DROP TABLE IF EXISTS stale_flag_policies;

ALTER TABLE environment_flag_states
    DROP COLUMN IF EXISTS lifecycle_state,
    DROP COLUMN IF EXISTS last_evaluated_at,
    DROP COLUMN IF EXISTS last_state_change_at;

DROP TYPE IF EXISTS flag_lifecycle_state;
