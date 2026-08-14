ALTER TABLE feature_flags ADD COLUMN variations JSONB DEFAULT '[]'::jsonb;
