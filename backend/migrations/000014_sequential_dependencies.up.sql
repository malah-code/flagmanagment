ALTER TABLE feature_flags
ADD COLUMN parent_flag_id UUID REFERENCES feature_flags(id) ON DELETE RESTRICT;

CREATE INDEX idx_feature_flags_parent_flag_id ON feature_flags(parent_flag_id);
