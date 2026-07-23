# Quickstart Validation Guide

Once the API is implemented and running (e.g., `go run cmd/server/main.go`), you can validate it end-to-end using `curl`.

1. **Create a Project**
   ```bash
   PROJECT_ID=$(curl -s -X POST http://localhost:8080/api/v1/projects \
     -H "Content-Type: application/json" \
     -d '{"name": "Quickstart Project"}' | jq -r '.id')
   echo "Project ID: $PROJECT_ID"
   ```

2. **Create an Environment**
   ```bash
   ENV_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/projects/$PROJECT_ID/environments \
     -H "Content-Type: application/json" \
     -d '{"name": "Production"}')
   ENV_ID=$(echo $ENV_RESPONSE | jq -r '.id')
   API_KEY=$(echo $ENV_RESPONSE | jq -r '.apiKey')
   echo "Environment ID: $ENV_ID"
   echo "API Key: $API_KEY"
   ```

3. **Create a Flag**
   ```bash
   FLAG_ID=$(curl -s -X POST http://localhost:8080/api/v1/projects/$PROJECT_ID/flags \
     -H "Content-Type: application/json" \
     -d '{"key": "new-dashboard", "type": "boolean"}' | jq -r '.id')
   echo "Flag ID: $FLAG_ID"
   ```

4. **Update Flag State**
   ```bash
   curl -s -X PUT http://localhost:8080/api/v1/environments/$ENV_ID/flags/$FLAG_ID/state \
     -H "Content-Type: application/json" \
     -d '{"state": true, "targetingRules": [], "remoteConfig": {}}'
   ```

5. **Fetch SDK Payload**
   ```bash
   curl -s -i -X GET http://localhost:8080/api/v1/evaluate/flags \
     -H "Authorization: Bearer $API_KEY"
   ```
   *Note: This will return a 200 OK with the ETag header. A subsequent request with `If-None-Match: <etag>` should return 304 Not Modified.*
