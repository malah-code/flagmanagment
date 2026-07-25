# Phase 0: Research

## Webhook Authentication Strategy
- **Decision**: Authenticate webhooks using the Environment's API Key passed as a Bearer token in the Authorization header.
- **Rationale**: Reuses the existing robust API key infrastructure for environment isolation. APM tools (Datadog, New Relic) support custom headers for webhooks out of the box.
- **Alternatives considered**: HMAC signatures. While more secure against tampering, it requires adding a new `webhook_secret` to the Environment model and complicates APM tool configuration.

## Alert Identification
- **Decision**: Trigger kill-switches based on a user-provided string `alert_identifier`. The webhook payload must contain a flat `alert_id` or `alert_name` field in the JSON body that matches exactly.
- **Rationale**: Simple to implement, covers 90% of use cases. Most APMs can customize webhook JSON bodies to map complex conditions down to a single identifying string.
- **Alternatives considered**: Full JSONPath expression evaluation on the incoming payload. Rejected as too complex for MVP.

## Event Processing & Performance
- **Decision**: Webhook endpoint will locate the environment, parse the `alert_identifier`, query `kill_switches` for that identifier, disable the associated flags via `FlagService.UpdateFlagState()`, log to `AuditService`, and invalidate the cache. All performed synchronously but fast.
- **Rationale**: A typical system won't receive thousands of kill-switch alerts per second. Processing it synchronously ensures immediate 202 acceptance and instant invalidation, keeping edge SDKs updated quickly.
- **Alternatives considered**: Asynchronous processing via a message queue (Kafka/RabbitMQ). Too heavy for current architecture.
