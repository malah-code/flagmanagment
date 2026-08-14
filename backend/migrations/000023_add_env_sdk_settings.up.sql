ALTER TABLE environments ADD COLUMN IF NOT EXISTS sdk_settings JSONB NOT NULL DEFAULT '{}'::jsonb;
