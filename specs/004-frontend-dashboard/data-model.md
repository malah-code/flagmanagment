# Data Model: Frontend Dashboard

The frontend consumes the REST API defined in the backend. These are the primary TypeScript interfaces it needs:

## Entities

### `Project`
```typescript
interface Project {
  id: string;
  key: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
}
```

### `Environment`
```typescript
interface Environment {
  id: string;
  projectId: string;
  key: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

interface CreateEnvironmentResponse {
  environment: Environment;
  apiKey: string; // Only returned once upon creation
}
```

### `FeatureFlag`
```typescript
interface FeatureFlag {
  id: string;
  projectId: string;
  key: string;
  name: string;
  description: string;
  type: "boolean" | "multivariate";
  createdAt: string;
  updatedAt: string;
}
```

### `FlagState`
```typescript
interface TargetingRule {
  attribute: string;
  operator: "eq" | "in" | "contains" | "regex";
  value: any;
}

interface FlagState {
  environmentId: string;
  flagId: string;
  enabled: boolean;
  rules: TargetingRule[]; // simplified rule structure for UI MVP
  updatedAt: string;
}
```

### `PaginatedResponse<T>`
```typescript
interface PaginatedResponse<T> {
  items: T[];
  nextPageToken?: string;
  totalSize: number;
}
```

## Validation Rules (Frontend)
- **Keys**: All `key` fields (Project, Environment, FeatureFlag) must be valid slugs: lowercase alphanumeric separated by hyphens (e.g., `my-project-key`).
- **Names**: Required and must be between 3 and 100 characters.
- **Descriptions**: Optional, maximum 255 characters.
