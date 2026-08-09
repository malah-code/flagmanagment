# API Contract: Remote Configuration Payload UI

No external or internal API contract modifications are required for this feature.

## Rationale
The existing `/api/v1/projects/{id}/flags` POST and PUT endpoints accept a standard `CreateFlagRequest` schema:
```json
{
  "key": "string",
  "name": "string",
  "type": "JSON",
  "variations": [
    {
      "id": "var_a",
      "name": "Default Payload",
      "value": { "color": "blue", "limit": 10 }
    }
  ]
}
```

Since the `value` property mapped to `interface{}` in Go is serialized out-of-the-box via standard JSON reflection rules, the API correctly persists and retrieves this payload without modification. The SDK evaluation endpoints similarly forward this `value` identically as it is stored in `JSONB`.

## Next Steps
The frontend needs to be capable of parsing stringified JSON inputs from a code editor component, serializing them into valid Javascript objects (`JSON.parse(val)`), and submitting them within the `value` field of the existing Variation shape.
