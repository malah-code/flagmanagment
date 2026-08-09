# Research: Enterprise SSO & Authentication

## OIDC Client Implementation
- **Decision**: Use `coreos/go-oidc/v3`.
- **Rationale**: It is the industry standard for Go OIDC relying party implementation, maintained by the CoreOS/RedHat team. It handles discovery and signature verification automatically.
- **Alternatives considered**: Writing raw OAuth2 via `golang.org/x/oauth2` (we will use this for the base OAuth2 flow, wrapped by `go-oidc` for ID token verification).

## SAML 2.0 Implementation
- **Decision**: Use `github.com/crewjam/saml`.
- **Rationale**: Robust, widely used SAML Service Provider library for Go.
- **Alternatives considered**: Forcing customers to use SAML-to-OIDC proxies (like Dex), but native SAML support is frequently requested by enterprise buyers.

## Service Account Token Format
- **Decision**: Opaque cryptographically random tokens stored hashed in the DB, OR long-lived JWTs signed with a rotating asymmetric key. Given time constraints, opaque tokens hashed via SHA-256 in the database are the most secure and easiest to revoke.
- **Rationale**: Avoids the complexities of JWT revocation (blocklists). Checking a hash against the DB is fast enough for API requests.
- **Alternatives considered**: Signed JWTs.
