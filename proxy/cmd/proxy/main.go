package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flagmanagment/proxy/internal/broadcaster"
	"github.com/flagmanagment/proxy/internal/config"
	"github.com/flagmanagment/proxy/internal/health"
	"github.com/flagmanagment/proxy/internal/server"
	"github.com/flagmanagment/proxy/internal/store"
	"github.com/flagmanagment/proxy/internal/upstream"
	pb "github.com/flagmanagment/proxy/pkg/gen/sdk/v1"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	cfg := config.Load()

	// Zerolog setup
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if cfg.LogFormat == "text" || (cfg.LogFormat == "auto" && os.Getenv("TERM") != "") {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	log.Info().
		Str("grpc_port", cfg.GRPCPort).
		Str("health_port", cfg.HealthPort).
		Str("backend_addr", cfg.BackendAddr).
		Msg("Starting FlagManagment Edge Proxy")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rStore := store.NewRulesetStore()
	bCaster := broadcaster.NewBroadcaster()
	uState := upstream.NewUpstreamState()
	uClient := upstream.NewUpstreamClient(cfg, rStore, bCaster, uState)

	// Start Upstream Client Goroutine
	go uClient.Run(ctx)

	// Start Downstream gRPC Server
	grpcLis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to listen on gRPC port %s", cfg.GRPCPort)
	}

	var grpcOpts []grpc.ServerOption
	if cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertFile, cfg.TLSKeyFile)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to load TLS certificates for downstream gRPC server")
		}
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(grpcOpts...)
	proxyServer := server.NewProxyServer(cfg, rStore, bCaster)
	pb.RegisterSDKServiceServer(grpcServer, proxyServer)

	go func() {
		log.Info().Str("port", cfg.GRPCPort).Msg("Downstream gRPC server listening")
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Error().Err(err).Msg("Downstream gRPC server stopped")
		}
	}()

	// Start Health HTTP Server
	healthRouter := chi.NewRouter()
	healthRouter.Handle("/healthz", health.NewHealthHandler(cfg, rStore, bCaster, uState))

	healthServer := &http.Server{
		Addr:    ":" + cfg.HealthPort,
		Handler: healthRouter,
	}

	go func() {
		log.Info().Str("port", cfg.HealthPort).Msg("Health HTTP server listening")
		if err := healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Health HTTP server stopped")
		}
	}()

	// Signal handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info().Msg("Shutting down Edge Proxy gracefully...")

	cancel() // Stop upstream client
	grpcServer.GracefulStop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = healthServer.Shutdown(shutdownCtx)

	log.Info().Msg("Edge Proxy shutdown complete")
}
