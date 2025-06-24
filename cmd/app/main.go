package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adal4ik/people-enrichment-service/internal/config"
	"github.com/adal4ik/people-enrichment-service/internal/handler"
	"github.com/adal4ik/people-enrichment-service/internal/logger"
	"github.com/adal4ik/people-enrichment-service/internal/repository"
	"github.com/adal4ik/people-enrichment-service/internal/service"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	db, err := repository.ConnectDB(ctx, cfg)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	logger.Info("successfully connected to database")
	repositories := repository.New(db, logger)
	services := service.New(repositories, cfg, logger)
	handlers := handler.New(services, logger)
	mux := handler.Router(*handlers)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		log.Println("Server is running on port: http://localhost" + httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("ListenAndServe error", zap.Error(err))
		}
	}()
	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	logger.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("HTTP server shutdown error", zap.Error(err))
	}

	logger.Info("Server stopped gracefully")
}
