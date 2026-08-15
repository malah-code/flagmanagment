# Tasks: CI/CD Pipeline Enhancements & Release Automation

**Feature**: [spec.md](spec.md) | **Plan**: [plan.md](plan.md)

## Phase 1: Quality Gates & Vulnerability Scanning (US1 & US2)

- [x] T001 Configure push triggers for `main` branch in `.github/workflows/ci.yml` per FR-001
- [x] T002 Integrate Aqua Security Trivy container vulnerability scanner in `.github/workflows/ci.yml` per FR-002
- [x] T003 Ensure backend test coverage threshold (>=80%) and frontend Vitest suite execute on all pushes per FR-001

## Phase 2: Release Automation & Multi-Arch Publishing (US3)

- [x] T004 Update `.github/workflows/publish.yml` with `docker/metadata-action` and semver tag triggers per FR-003, FR-004
- [x] T005 Create `.github/workflows/release.yml` for automated GitHub Releases on `v*.*.*` tags per FR-005
- [x] T006 Commit and push workflow definitions to repository `main` branch
