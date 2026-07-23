# Phase 0: Research & Technical Decisions

## API Framework & Routing
- **Decision**: Use `go-chi/chi` for REST routing.
- **Rationale**: It is already imported into the backend for the health endpoint. It is standard, fast, and 100% compatible with `net/http`.
- **Alternatives considered**: `gin` (too heavy, non-standard signature), standard `http.ServeMux` (Go 1.22 has better routing, but `chi` provides excellent middleware like RequestID, Logger, and Recoverer out of the box).

## Validation
- **Decision**: Use `github.com/go-playground/validator/v10` for JSON payload validation.
- **Rationale**: It provides rich, struct-tag based validation rules (e.g., `validate:"required,uuid"`), saving boilerplate validation code in every handler.
- **Alternatives considered**: Manual validation (error-prone and verbose).

## Pagination
- **Decision**: Cursor-based pagination (`pageToken`, `pageSize`) via encoding the last seen entity ID.
- **Rationale**: Required by the Clarification session to follow Google API standards and scale better than offset.

## API Key Authentication
- **Decision**: Middleware checking the `Authorization: Bearer <key>` header, hashing it with SHA-256, and looking it up in the `environments` table via the `idx_environments_api_key_hash` index.
- **Rationale**: Highly secure and extremely fast. SDKs simply pass the raw key; the server hashes it once per request.

## Caching / ETag for SDK
- **Decision**: The SDK evaluation endpoint will hash the concatenated `updated_at` timestamps of all flags in the environment to generate an ETag. 
- **Rationale**: Allows HTTP 304 Not Modified responses when no flags have been updated, saving massive bandwidth on polling SDKs.
