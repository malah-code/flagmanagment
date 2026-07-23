package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
)

type contextKey string

const EnvKey contextKey = "environment"

// AuthMiddleware extracts the API key, hashes it, and looks up the Environment.
func AuthMiddleware(store repository.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				RespondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
				return
			}

			apiKey := strings.TrimPrefix(authHeader, "Bearer ")
			if apiKey == "" {
				RespondWithError(w, http.StatusUnauthorized, "Empty API key")
				return
			}

			// Hash the API key using SHA-256
			hash := sha256.Sum256([]byte(apiKey))
			hashHex := hex.EncodeToString(hash[:])

			// Lookup environment
			env, err := store.EnvironmentRepo().GetByAPIKeyHash(r.Context(), hashHex)
			if err != nil {
				RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
				return
			}

			// Add to context
			ctx := context.WithValue(r.Context(), EnvKey, env)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetEnvironmentFromContext extracts the environment from the request context.
func GetEnvironmentFromContext(ctx context.Context) *models.Environment {
	if env, ok := ctx.Value(EnvKey).(*models.Environment); ok {
		return env
	}
	return nil
}

// Pagination represents standard API pagination fields
type Pagination struct {
	PageSize  int
	PageToken string
}

// GetPagination extracts pagination parameters from the URL
func GetPagination(r *http.Request) Pagination {
	pageSize := 50
	if sizeStr := r.URL.Query().Get("pageSize"); sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil && size > 0 {
			pageSize = size
			if pageSize > 1000 {
				pageSize = 1000
			}
		}
	}

	pageToken := r.URL.Query().Get("pageToken")

	return Pagination{
		PageSize:  pageSize,
		PageToken: pageToken,
	}
}

// TokenToOffset converts a pageToken string back to an integer offset.
// For MVP, we'll just encode the integer offset directly in base10 strings, or just assume it's the offset string.
func TokenToOffset(token string) int {
	if token == "" {
		return 0
	}
	offset, err := strconv.Atoi(token)
	if err != nil {
		return 0
	}
	return offset
}

// OffsetToToken converts a numeric offset to an opaque token.
func OffsetToToken(offset int) string {
	return strconv.Itoa(offset)
}
