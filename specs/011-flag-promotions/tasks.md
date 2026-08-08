# Tasks: One-Click Flag Environment Promotions

**Input**: Design documents from `/specs/011-flag-promotions/`

## Phase 1: Backend API & Service

- [x] T001 Create `PromotionService` in `backend/internal/services/promotion_service.go` containing logic to copy state or generate Change Request
- [x] T002 Add unit tests for `PromotionService` in `backend/internal/services/promotion_service_test.go`
- [x] T003 Create `PromotionHandler` REST API in `backend/internal/api/promotion.go` mapping to `POST /api/v1/projects/{projectId}/flags/{flagId}/promote`
- [x] T004 Register `PromotionHandler` routes in `backend/cmd/server/main.go`

## Phase 2: Frontend Component & Integration

- [x] T005 Add `promoteFlag` to `frontend/src/services/flags.ts`
- [x] T006 Create `PromoteFlagModal.tsx` in `frontend/src/components/flagStates/PromoteFlagModal.tsx`
- [x] T007 Integrate `PromoteFlagModal` into flag states view / action menu
