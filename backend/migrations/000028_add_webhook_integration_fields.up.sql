ALTER TABLE webhook_integrations ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT 'Custom Webhook';
ALTER TABLE webhook_integrations ADD COLUMN integration_type VARCHAR(50) NOT NULL DEFAULT 'custom';
