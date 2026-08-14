ALTER TABLE feature_flags ADD COLUMN tags TEXT[] DEFAULT '{}'::TEXT[];
