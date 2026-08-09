package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/flagmanagment/backend/internal/auth"
	"github.com/flagmanagment/backend/internal/repository"
)

const UserIDKey contextKey = "user_id"
const ServiceAccountIDKey contextKey = "service_account_id"

func UserAuthMiddleware(store repository.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			if strings.HasPrefix(tokenStr, "fm_sa_") {
				hash := sha256.Sum256([]byte(tokenStr))
				keyHash := hex.EncodeToString(hash[:])

				key, err := store.ServiceAccountRepo().GetKeyByHash(r.Context(), keyHash)
				if err != nil {
					http.Error(w, "Unauthorized: invalid service account key", http.StatusUnauthorized)
					return
				}

				if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
					http.Error(w, "Unauthorized: service account key expired", http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), ServiceAccountIDKey, key.ServiceAccountID)
				// For compatibility with handlers expecting a user ID string
				ctx = context.WithValue(ctx, UserIDKey, key.ServiceAccountID.String())
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			claims, err := auth.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
