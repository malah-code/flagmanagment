# Research: Frontend Dashboard

## State Management and Data Fetching
- **Decision**: `@tanstack/react-query`
- **Rationale**: The frontend needs to interact with multiple REST API endpoints (Projects, Environments, Flags, Flag States). It must handle loading states (FR-007), caching, and optimistic updates (e.g., toggling a flag state). React Query provides built-in hooks for querying and mutations that will dramatically reduce boilerplate.
- **Alternatives considered**: Standard `fetch` with `useEffect` (too much boilerplate, error-prone for race conditions), Redux (overkill for simple CRUD dashboard).

## UI Component Library
- **Decision**: Shadcn/UI with Tailwind CSS
- **Rationale**: The Constitution mandates Shadcn/UI or MUI. Shadcn/UI integrates perfectly with Vite/React, provides highly customizable components, and enforces a clean, modern aesthetic out of the box.
- **Alternatives considered**: MUI (heavier, harder to customize cleanly with Tailwind), custom CSS (violates Constitution's mandate for standardized library).

## Routing
- **Decision**: `react-router-dom`
- **Rationale**: We need standard client-side routing to navigate between `/projects`, `/projects/:id/environments`, and `/projects/:id/flags`.
- **Alternatives considered**: TanStack Router (newer, type-safe, but React Router is the industry standard and sufficient for this MVP).

## API Client
- **Decision**: Native `fetch` API wrapped in service functions
- **Rationale**: Keeps the bundle size small and works perfectly with React Query.
- **Alternatives considered**: Axios (adds unnecessary dependency when native fetch is well-supported).
