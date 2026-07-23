package health

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var Version = "0.1.0"

type CheckResult struct {
	Status    string `json:"status"`
	LatencyMs int64  `json:"latency_ms,omitempty"`
	Error     string `json:"error,omitempty"`
}

type Response struct {
	Status        string                 `json:"status"`
	Version       string                 `json:"version"`
	UptimeSeconds int64                  `json:"uptime_seconds"`
	Checks        map[string]CheckResult `json:"checks"`
}

var startTime = time.Now()

func NewHandler(db *sql.DB, rdb *redis.Client, logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		overallStatus := "healthy"
		checks := make(map[string]CheckResult)

		// Check Postgres
		start := time.Now()
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		err := db.PingContext(ctx)
		cancel()
		latency := time.Since(start).Milliseconds()

		if err != nil {
			overallStatus = "unhealthy"
			checks["postgres"] = CheckResult{
				Status: "unhealthy",
				Error:  err.Error(),
			}
			logger.Error().Err(err).Msg("postgres health check failed")
		} else {
			checks["postgres"] = CheckResult{
				Status:    "healthy",
				LatencyMs: latency,
			}
		}

		// Check Redis
		start = time.Now()
		ctx2, cancel2 := context.WithTimeout(r.Context(), 2*time.Second)
		err = rdb.Ping(ctx2).Err()
		cancel2()
		latency = time.Since(start).Milliseconds()

		if err != nil {
			overallStatus = "unhealthy"
			checks["redis"] = CheckResult{
				Status: "unhealthy",
				Error:  err.Error(),
			}
			logger.Error().Err(err).Msg("redis health check failed")
		} else {
			checks["redis"] = CheckResult{
				Status:    "healthy",
				LatencyMs: latency,
			}
		}

		resp := Response{
			Status:        overallStatus,
			Version:       Version,
			UptimeSeconds: int64(time.Since(startTime).Seconds()),
			Checks:        checks,
		}

		w.Header().Set("Content-Type", "application/json")
		if overallStatus == "healthy" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error().Err(err).Msg("failed to encode health check response")
		}
	}
}
