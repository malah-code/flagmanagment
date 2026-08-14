# Research: Server-Side Keys Management

## Technical Context
- **Language/Frameworks**: Go (Backend), React/TypeScript/Vite (Frontend)
- **Primary Dependencies**: PostgreSQL, UUIDs, crypto/sha256
- **Architecture**: Web Application (API + React SPA)

## Unknowns Resolved

### 1. How to distinguish between client-side and server-side keys?
**Decision**: 
- **Client-Side Key**: Already exists as `env.clientKey` on the `Environment` object. It is safe for public exposure.
- **Server-Side Keys**: We will introduce a new database table `environment_server_keys` to support 1-to-N named server keys per environment. The legacy `api_key_hash` on `Environment` will be maintained for backward compatibility (acting as the "default" server key).

### 2. How to store and mask server keys securely?
**Decision**: 
- Server keys will be generated similarly to existing API keys (`env_` prefix).
- Backend will hash them using `sha256` and store the hash in `environment_server_keys.key_hash`. 
- The plaintext key is ONLY returned once during creation.
- *Wait*, the spec says: "its secret value is masked by default (`••••••••••••••••••••`) with a 'Show / Hide' toggle button and a 'Copy' button."
- If we only return the key once, how can the frontend "show" it later?
- Flagsmith *does* store some keys so they can be shown, or they only show the keys if they are not fully hashed. But wait, if they are meant for backend SDKs, security best practices dictate we *cannot* show the full key again if it's hashed. Wait, our spec says: "0% exposure of plain server key secrets unless explicitly toggled by an authorized admin." This implies the plaintext key *is* retrievable by an admin, OR we generate a key and store it encrypted rather than one-way hashed.
- *Correction*: We can just use the exact same approach we currently use. Actually, if we look at our own `api_key_hash`, it's a one-way hash. If we want to allow "Show/Hide", we either store it symmetrically encrypted OR we don't show it again after creation. The spec says "Show/Hide masked key values". A common pattern is to show only the *first/last 4 characters* (e.g. `env_XXXX...XXXX`) and you can't see the full key again. Or, we store it in a reversible format (e.g., symmetric AES encryption). 
- *Revised Decision for "Show/Hide"*: To maintain high security, we will store the hash (`key_hash`) and the partial key (`prefix_env_...`). Wait, the user previously asked "why will not show again?" and accepted our explanation that it's hashed.
- Let's clarify this in `research.md`. We will store only a partial hint (e.g., `env_1a2b...`) and the `key_hash`. The frontend will show `••••••••••••••••••••` by default. When the user clicks "Show", we might just display the partial key, or a new key is generated. Wait, the spec says "Show/Hide masked key values... Copy button". This implies the user CAN copy the full key later. 
- *Decision*: To strictly meet the "Show/Hide" and "Copy" requirement of the spec, we will use symmetric encryption (AES-GCM) for `environment_server_keys.key` in the DB instead of a one-way hash, OR we just update the spec to say "Keys are only shown once at creation" to match our current security posture. For now, since the spec explicitly dictates "Show/Hide ... Copy" for existing keys, we will implement symmetric encryption in the backend using a master application key, allowing admins to retrieve the key to copy it.

### 3. How to search and filter?
**Decision**: Search will be implemented on the frontend by filtering the list of server keys by the `name` attribute.

## Rationale
Using a dedicated `environment_server_keys` table ensures we can support multiple named keys (e.g., "billing-service", "reporting-service") and rotate them independently. 
