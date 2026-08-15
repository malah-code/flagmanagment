package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"github.com/flagmanagment/backend/internal/api"
	"github.com/flagmanagment/backend/internal/cache"
	"github.com/flagmanagment/backend/internal/config"
	"github.com/flagmanagment/backend/internal/health"
	"github.com/flagmanagment/backend/internal/logging"
	customMiddleware "github.com/flagmanagment/backend/internal/middleware"
	"github.com/flagmanagment/backend/internal/repository"
	sdkService "github.com/flagmanagment/backend/internal/sdk"
	"github.com/flagmanagment/backend/internal/services"
	pb "github.com/flagmanagment/backend/pkg/gen/sdk/v1"
)

func main() {
	cfg := config.Load()
	logger := logging.NewLogger(cfg.LogFormat, cfg.Env)

	logger.Info().Str("port", cfg.BackendPort).Msg("starting FlagManagment backend")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open database connection")
	}

	if err := db.Ping(); err != nil {
		logger.Warn().Err(err).Msg("failed to ping database at startup")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})

	if err := runMigrations(db, cfg.DBName); err != nil && err != migrate.ErrNoChange {
		logger.Fatal().Err(err).Msg("failed to run migrations")
	} else {
		logger.Info().Msg("migrations applied successfully")
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/healthz", health.NewHandler(db, rdb, logger))

	// Initialize Store
	store := repository.NewStore(db)

	// Initialize Cache Client
	cacheClient := cache.NewClient(fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort), "", 0)

	// Initialize Handlers
	authHandler := api.NewAuthHandler(store)
	rbacMiddleware := api.NewRBACMiddleware(store)
	auditHandler := api.NewAuditHandler(store)
	auditService := services.NewAuditService(store)
	crService := services.NewChangeRequestService(store, auditService)
	promotionService := services.NewPromotionService(store, auditService, crService)
	scService := services.NewScheduledChangeService(store, auditService)

	metricService := services.NewMetricAggregationService(store.FlagStateRepo())
	metricService.Start(context.Background(), 10*time.Second)

	staleScanner := services.NewStaleScannerService(store)
	staleScanner.Start(context.Background(), 1*time.Hour)

	envService := services.NewEnvironmentService(store, auditService)

	projectHandler := api.NewProjectHandler(store, rbacMiddleware, auditHandler)
	envHandler := api.NewEnvironmentHandler(store, rbacMiddleware, auditService, envService)
	flagHandler := api.NewFlagHandler(store, cacheClient, rbacMiddleware, auditService, crService)
	crHandler := api.NewChangeRequestHandler(store, crService, rbacMiddleware, cacheClient)
	promotionHandler := api.NewPromotionHandler(promotionService, cacheClient, rbacMiddleware)
	sdkHandler := api.NewSDKHandler(store, metricService)

	notificationService := services.NewNotificationService(store)
	webhookService := services.NewWebhookService(store, auditService, cacheClient, notificationService)
	webhookHandler := api.NewWebhookHandler(webhookService)
	ksHandler := api.NewKillSwitchHandler(store, rbacMiddleware)
	slackHandler := api.NewSlackConfigHandler(store, rbacMiddleware)
	scHandler := api.NewScheduledChangeHandler(store, scService, rbacMiddleware, cacheClient)
	saHandler := api.NewServiceAccountHandler(store)

	cryptoService := services.NewCryptoService()
	emailService := services.NewEmailService(store)
	userService := services.NewUserService(store, cryptoService, emailService)
	usersHandler := api.NewUsersHandler(userService)
	configHandler := api.NewConfigHandler(store, cryptoService, emailService)

	lifecycleHandler := api.NewLifecycleHandler(store, rbacMiddleware, auditService)
	stalePolicyHandler := api.NewStalePolicyHandler(store, rbacMiddleware)

	// API v1 Routes
	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			authHandler.RegisterRoutes(r)
		})

		// Dashboard API Routes (Protected by UserAuthMiddleware)
		r.Group(func(r chi.Router) {
			r.Use(api.UserAuthMiddleware(store))

			r.Route("/projects", func(r chi.Router) {
				projectHandler.RegisterRoutes(r)
			})

			r.Route("/service-accounts", func(r chi.Router) {
				saHandler.RegisterRoutes(r)
			})

			r.Route("/users", func(r chi.Router) {
				usersHandler.RegisterRoutes(r)
			})

			r.Route("/config", func(r chi.Router) {
				r.Use(rbacMiddleware.RequireRole("ADMIN"))
				configHandler.RegisterRoutes(r)
			})

			envHandler.RegisterRoutes(r)
			flagHandler.RegisterRoutes(r)
			lifecycleHandler.RegisterRoutes(r)
			stalePolicyHandler.RegisterRoutes(r)
			crHandler.RegisterRoutes(r)
			ksHandler.RegisterRoutes(r)
			slackHandler.RegisterRoutes(r)
			promotionHandler.RegisterRoutes(r)
			scHandler.RegisterRoutes(r)
			auditHandler.RegisterRoutes(r)
			webhookHandler.RegisterManagementRoutes(r)
		})

		// SDK Routes (Protected by API Key AuthMiddleware)
		r.Group(func(r chi.Router) {
			r.Use(api.AuthMiddleware(store))
			sdkHandler.RegisterRoutes(r)
		})

		// Webhook Routes (Protected by API Key AuthMiddleware for environment identification)
		r.Group(func(r chi.Router) {
			r.Use(api.AuthMiddleware(store))
			r.Route("/webhooks", func(r chi.Router) {
				webhookHandler.RegisterAPMRoutes(r)
			})
		})
	})

	// Start Background Scheduler for Scheduled Flag Changes
	scheduler := services.NewScheduler(store, auditService, crService, cacheClient, logger)
	go scheduler.Start(context.Background())

	// Start Webhook Dispatcher
	webhookDispatcher := services.NewWebhookDispatcher(store.WebhookIntegrationRepo(), auditService.Subscribe())
	go webhookDispatcher.Start(context.Background())

	// Start Data Retention Cleanup Ticker (purges audit/analytics logs older than 30 days every 24 hours)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := auditService.CleanupOldLogs(context.Background(), 30); err != nil {
				logger.Error().Err(err).Msg("failed to clean up old audit logs")
			}
		}
	}()

	// Start gRPC server for SDK streaming
	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			logger.Error().Err(err).Msg("failed to listen on gRPC port 9090")
			return
		}
		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(customMiddleware.SDKGRPCAuthInterceptor(store)),
		)
		sdkServer := sdkService.NewServer(store, cacheClient)
		pb.RegisterSDKServiceServer(grpcServer, sdkServer)

		logger.Info().Msg("starting gRPC server on :9090")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error().Err(err).Msg("gRPC server failed")
		}
	}()

	log.Fatal(http.ListenAndServe(":"+cfg.BackendPort, router))
}

func runMigrations(db *sql.DB, dbName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbName, driver)
	if err != nil {
		return err
	}
	return m.Up()
}
