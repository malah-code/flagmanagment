# Tasks: Slack Notification Webhooks for Flag Updates

**Input**: Design documents from `/specs/010-slack-notifications/`

## Phase 1: Setup & Data Model

- [x] T001 Create database migration `backend/migrations/000011_create_slack_configs.up.sql`
- [x] T002 Create `SlackWebhookConfig` model in `backend/internal/models/slack_config.go`
- [x] T003 Create `SlackConfigRepository` in `backend/internal/repository/slack_config_repo.go` and register in `store.go`

---

## Phase 2: Backend API & Service

- [x] T004 Create `NotificationService` in `backend/internal/services/notification_service.go` to send Slack Block Kit messages asynchronously
- [x] T005 Create `SlackConfigHandler` REST API in `backend/internal/api/slack.go`
- [x] T006 Integrate `NotificationService` triggers into flag state toggles, change request approvals, and kill-switches
- [x] T007 Add unit tests for `NotificationService` in `backend/internal/services/notification_service_test.go`

---

## Phase 3: Frontend Component & Verification

- [x] T008 Create `slackApi.ts` in `frontend/src/services/slackApi.ts`
- [x] T009 Create `SlackConfigForm.tsx` in `frontend/src/components/SlackConfigForm.tsx`
- [x] T010 Integrate `SlackConfigForm` into environment settings view
