# Contract: Go Repository Interfaces

**Feature**: 002-data-model-state
**Date**: 2026-07-20

This document defines the Go interfaces for database access. These interfaces serve as contracts between the service layer and the repository implementations.

---

## Core Repository Interface

```go
// Store provides access to all repository operations.
// It is the single entry point for database access from the service layer.
type Store interface {
    ProjectRepo() ProjectRepository
    EnvironmentRepo() EnvironmentRepository
    FlagRepo() FlagRepository
    FlagStateRepo() FlagStateRepository
    AuditRepo() AuditRepository
    ChangeRequestRepo() ChangeRequestRepository
    RoleRepo() RoleRepository

    // Transaction support
    WithTx(ctx context.Context, fn func(Store) error) error

    // Migration support
    MigrateUp() error
    MigrateDown() error
}
```

---

## Entity Repositories

### ProjectRepository

```go
type ProjectRepository interface {
    Create(ctx context.Context, project *Project) error
    GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
    GetByKey(ctx context.Context, key string) (*Project, error)
    List(ctx context.Context, limit, offset int) ([]*Project, int, error)
    Update(ctx context.Context, project *Project) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### EnvironmentRepository

```go
type EnvironmentRepository interface {
    Create(ctx context.Context, env *Environment) error
    GetByID(ctx context.Context, id uuid.UUID) (*Environment, error)
    GetByAPIKeyHash(ctx context.Context, apiKeyHash string) (*Environment, error)
    ListByProject(ctx context.Context, projectID uuid.UUID) ([]*Environment, error)
    Update(ctx context.Context, env *Environment) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### FlagRepository

```go
type FlagRepository interface {
    Create(ctx context.Context, flag *FeatureFlag) error
    GetByID(ctx context.Context, id uuid.UUID) (*FeatureFlag, error)
    GetByKey(ctx context.Context, projectID uuid.UUID, key string) (*FeatureFlag, error)
    ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*FeatureFlag, int, error)
    Update(ctx context.Context, flag *FeatureFlag) error
    Delete(ctx context.Context, id uuid.UUID) error
    UpdateLastEvaluatedAt(ctx context.Context, ids []uuid.UUID) error
}
```

### FlagStateRepository

```go
type FlagStateRepository interface {
    Create(ctx context.Context, state *EnvironmentFlagState) error
    GetByID(ctx context.Context, id uuid.UUID) (*EnvironmentFlagState, error)
    GetByEnvAndFlag(ctx context.Context, envID, flagID uuid.UUID) (*EnvironmentFlagState, error)
    ListByEnvironment(ctx context.Context, envID uuid.UUID) ([]*EnvironmentFlagState, error)
    Update(ctx context.Context, state *EnvironmentFlagState) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### AuditRepository

```go
type AuditRepository interface {
    Create(ctx context.Context, log *AuditLog) error
    ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*AuditLog, int, error)
    ListByEnvironment(ctx context.Context, envID uuid.UUID, limit, offset int) ([]*AuditLog, int, error)
}
```

### ChangeRequestRepository

```go
type ChangeRequestRepository interface {
    Create(ctx context.Context, cr *ChangeRequest) error
    GetByID(ctx context.Context, id uuid.UUID) (*ChangeRequest, error)
    ListByEnvironment(ctx context.Context, envID uuid.UUID, status string, limit, offset int) ([]*ChangeRequest, int, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status string, appliedBy *uuid.UUID) error
    AddApproval(ctx context.Context, approval *ChangeRequestApproval) error
    ListApprovals(ctx context.Context, crID uuid.UUID) ([]*ChangeRequestApproval, error)
}
```

### RoleRepository

```go
type RoleRepository interface {
    Create(ctx context.Context, role *Role) error
    GetByID(ctx context.Context, id uuid.UUID) (*Role, error)
    GetByName(ctx context.Context, name string) (*Role, error)
    List(ctx context.Context) ([]*Role, error)
    AssignUserRole(ctx context.Context, ur *UserRole) error
    RemoveUserRole(ctx context.Context, id uuid.UUID) error
    GetUserRoles(ctx context.Context, userID uuid.UUID, projectID *uuid.UUID) ([]*UserRole, error)
}
```

---

## Error Contract

```go
var (
    ErrNotFound      = errors.New("resource not found")
    ErrAlreadyExists = errors.New("resource already exists")
    ErrConflict      = errors.New("resource conflict")
)
```

All repository methods return `ErrNotFound` when a lookup fails, `ErrAlreadyExists` for unique constraint violations, and wrapped `error` for unexpected database failures.
