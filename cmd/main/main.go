package main

import (
	"context"
	"log"
	"medods-jwt/internal/config"
	"medods-jwt/internal/repository"
	"medods-jwt/internal/service"
	"medods-jwt/internal/trasnport/rest"
	"medods-jwt/pkg/db/migrations"
	"medods-jwt/pkg/db/postgres"
	"medods-jwt/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg := config.New()
	if cfg == nil {
		log.Fatal("Config is nil")
	}

	mainLogger := logger.New(cfg.Debug)
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)

	mainLogger.Info("Starting application")

	db, err := postgres.New(&ctx, cfg.PostgresConfig)
	if err != nil {
		mainLogger.Fatal("Database connection error", zap.Error(err))
	}
	mainLogger.Debug("Database connection success")

	migrationsVersion, err := migrations.Up(db.DB)
	if err != nil {
		mainLogger.Fatal("Failed to apply migrations", zap.Error(err))
	}
	if migrationsVersion == 0 {
		mainLogger.Debug("No new migrations")
	}

	authRepo := repository.NewAuthRepository(&ctx, db)
	emailRepo := repository.NewEmailRepository(&ctx)

	authService := service.NewAuthService(&ctx, cfg.SecretKey, authRepo, emailRepo)

	server := rest.New(&ctx, cfg.ServerConfig, authService, cfg.Debug)

	graceChannel := make(chan os.Signal, 1)
	signal.Notify(graceChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Run(); err != nil {
			mainLogger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	<-graceChannel
	db.Close()
	mainLogger.Debug("Database connection closed")
	mainLogger.Info("Graceful shutdown!")
}
