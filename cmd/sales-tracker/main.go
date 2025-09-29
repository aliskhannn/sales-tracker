package main

import (
	"context"
	"errors"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"

	"github.com/aliskhannn/sales-tracker/internal/api/handler/analytics"
	"github.com/aliskhannn/sales-tracker/internal/api/handler/category"
	"github.com/aliskhannn/sales-tracker/internal/api/handler/item"
	"github.com/aliskhannn/sales-tracker/internal/api/router"
	"github.com/aliskhannn/sales-tracker/internal/api/server"
	"github.com/aliskhannn/sales-tracker/internal/config"
	repoanalytics "github.com/aliskhannn/sales-tracker/internal/repository/analytics"
	repocategory "github.com/aliskhannn/sales-tracker/internal/repository/category"
	repoitem "github.com/aliskhannn/sales-tracker/internal/repository/item"
	srvcanalytics "github.com/aliskhannn/sales-tracker/internal/service/analytics"
	srvccategory "github.com/aliskhannn/sales-tracker/internal/service/category"
	srvcitem "github.com/aliskhannn/sales-tracker/internal/service/item"
)

func main() {
	// Initialize logger, configuration and validator.
	zlog.Init()
	cfg := config.MustLoad()
	val := validator.New()

	// Connect to PostgreSQL master and slave databases.
	opts := &dbpg.Options{
		MaxOpenConns:    cfg.Database.MaxOpenConnections,
		MaxIdleConns:    cfg.Database.MaxIdleConnections,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}

	slaveDNSs := make([]string, 0, len(cfg.Database.Slaves))

	for _, s := range cfg.Database.Slaves {
		slaveDNSs = append(slaveDNSs, s.DSN())
	}

	db, err := dbpg.New(cfg.Database.Master.DSN(), slaveDNSs, opts)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	// Initialize category repository, service, and handler for category endpoints.
	categoryRepo := repocategory.NewRepository(db)
	categoryService := srvccategory.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService, val)

	// Initialize item repository, service, and handler for item endpoints.
	itemRepo := repoitem.NewRepository(db)
	itemService := srvcitem.NewService(itemRepo)
	itemHandler := item.NewHandler(itemService, val)

	// Initialize analytics repository, service, and handler for analytics endpoints.
	analyticsRepo := repoanalytics.NewRepository(db)
	analyticsService := srvcanalytics.NewService(analyticsRepo)
	analyticsHandler := analytics.NewHandler(analyticsService, cfg)

	// Initialize API router and HTTP server.
	r := router.New(categoryHandler, itemHandler, analyticsHandler)
	s := server.New(cfg, r)

	// Start HTTP server in a separate goroutine.
	go func() {
		if err := s.ListenAndServe(); err != nil {
			zlog.Logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// Setup context to handle SIGINT and SIGTERM for graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Wait for shutdown signal.
	<-ctx.Done()
	zlog.Logger.Print("shutdown signal received")

	// Gracefully shutdown server with timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	zlog.Logger.Print("gracefully shutting down server...\n")
	if err := s.Shutdown(shutdownCtx); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to shutdown server")
	}
	if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
		zlog.Logger.Info().Msg("timeout exceeded, forcing shutdown")
	}

	zlog.Logger.Print("closing master and slave databases...\n")

	// Close master database connection.
	if err := db.Master.Close(); err != nil {
		zlog.Logger.Printf("failed to close master DB: %v", err)
	}

	// Close slave database connections.
	for i, s := range db.Slaves {
		if err := s.Close(); err != nil {
			zlog.Logger.Printf("failed to close slave DB %d: %v", i, err)
		}
	}
}
