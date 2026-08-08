DROP INDEX IF EXISTS idx_feature_flags_parent_flag_id;

ALTER TABLE feature_flags
DROP COLUMN IF EXISTS parent_flag_id;
