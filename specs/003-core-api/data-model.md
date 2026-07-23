# API Data Models (DTOs)

The API layer introduces Data Transfer Objects (DTOs) that separate the internal database models (defined in Feature 002) from the external API contract.

## Management API DTOs

### Project
```json
{
  "id": "uuid",
  "name": "string",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

### Environment
*Note: `apiKey` is only returned in the response of the POST create endpoint.*
```json
{
  "id": "uuid",
  "projectId": "uuid",
  "name": "string",
  "apiKey": "string (create only)",
  "isProtected": "boolean",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

### Feature Flag
```json
{
  "id": "uuid",
  "projectId": "uuid",
  "key": "string",
  "type": "string (boolean, string, number, json)",
  "parentFlagId": "uuid (optional)",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

### Flag State (Environment specific)
```json
{
  "environmentId": "uuid",
  "featureFlagId": "uuid",
  "state": "boolean",
  "targetingRules": [
    // JSONB structure defining segments/rules
  ],
  "remoteConfig": {
    // JSONB payload for non-boolean flags
  },
  "updatedAt": "timestamp"
}
```

## SDK API Payload
*This is an optimized view aggregating flags and their states for a single environment.*
```json
{
  "environmentId": "uuid",
  "flags": {
    "my-flag-key": {
      "state": true,
      "type": "boolean",
      "rules": [],
      "value": null
    }
  }
}
```
