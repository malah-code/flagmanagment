package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const EnvKey contextKey = "environment"

// ValidateSDKToken hashes the input token and looks up the Environment.
func ValidateSDKToken(ctx context.Context, store repository.Store, token string) (*models.Environment, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}

	hash := sha256.Sum256([]byte(token))
	hashHex := hex.EncodeToString(hash[:])

	env, err := store.EnvironmentRepo().GetByAPIKeyHash(ctx, hashHex)
	if err != nil {
		return nil, errors.New("invalid SDK token")
	}
	return env, nil
}

// SDKHTTPAuthMiddleware checks the SDK token for REST HTTP requests.
func SDKHTTPAuthMiddleware(store repository.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"Missing or invalid Authorization header"}`, http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			env, err := ValidateSDKToken(r.Context(), store, token)
			if err != nil {
				http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), EnvKey, env)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SDKGRPCAuthInterceptor provides unary and stream auth interceptors for gRPC.
func SDKGRPCAuthInterceptor(store repository.Store) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		token := strings.TrimPrefix(tokens[0], "Bearer ")
		env, err := ValidateSDKToken(ctx, store, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization token")
		}

		newCtx := context.WithValue(ctx, EnvKey, env)
		return handler(newCtx, req)
	}
}

// GetEnvironmentFromContext retrieves the Environment struct from the context.
func GetEnvironmentFromContext(ctx context.Context) *models.Environment {
	if env, ok := ctx.Value(EnvKey).(*models.Environment); ok {
		return env
	}
	return nil
}
