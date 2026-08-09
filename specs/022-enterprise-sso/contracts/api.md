# API Contracts: Enterprise SSO & Authentication

## 1. SSO Endpoints

### `GET /api/v1/auth/sso/login?provider={oidc|saml}`
Initiates the SSO login flow.
- **Response**: `302 Redirect` to the Identity Provider's authorization URL.

### `GET /api/v1/auth/sso/callback/oidc`
The callback URL for OIDC authentication.
- **Query Params**: `code`, `state`
- **Response**: `302 Redirect` to the frontend dashboard with a secure session cookie (or JWT in URL fragment if using a stateless SPA approach, though cookies are preferred for security).

### `POST /api/v1/auth/sso/callback/saml`
The Assertion Consumer Service (ACS) endpoint for SAML.
- **Form Data**: `SAMLResponse`, `RelayState`
- **Response**: `302 Redirect` to the frontend dashboard with a secure session cookie.

## 2. Service Account Endpoints

### `POST /api/v1/service-accounts`
Creates a new Service Account.
- **Request Body**:
  ```json
  {
    "name": "Terraform Provisioner",
    "description": "Used by CI pipeline"
  }
  ```
- **Response**: `201 Created`
  ```json
  {
    "id": "uuid",
    "name": "Terraform Provisioner",
    "created_at": "timestamp"
  }
  ```

### `POST /api/v1/service-accounts/{id}/keys`
Generates a new API key for the Service Account.
- **Request Body**:
  ```json
  {
    "name": "CI Token",
    "expires_in_days": 30
  }
  ```
- **Response**: `201 Created`
  ```json
  {
    "id": "uuid",
    "name": "CI Token",
    "token": "fm_sa_XXXXXXXXXXXXXXXXXXXXXXXXX" // ONLY SHOWN ONCE
  }
  ```

### `GET /api/v1/service-accounts/{id}/keys`
Lists keys for a Service Account. (Does not return the plaintext token).
- **Response**: `200 OK`
  ```json
  [
    {
      "id": "uuid",
      "name": "CI Token",
      "expires_at": "timestamp",
      "last_used_at": "timestamp"
    }
  ]
  ```
