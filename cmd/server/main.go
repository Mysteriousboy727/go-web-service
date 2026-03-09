package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"go-industry-server/configs"
	"go-industry-server/internal/handler"
	"go-industry-server/internal/middleware"
	"go-industry-server/internal/repository"
	"go-industry-server/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg := configs.Load()

	// Connect to PostgreSQL
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connected")

	userRepo := repository.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepo, logger)
	userHandler := handler.NewUserHandler(userService, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("POST /api/v1/users", userHandler.Create)
	mux.HandleFunc("GET /api/v1/users", userHandler.List)
	mux.HandleFunc("GET /api/v1/users/{id}", userHandler.GetByID)
	mux.HandleFunc("PUT /api/v1/users/{id}", userHandler.Update)
	mux.HandleFunc("DELETE /api/v1/users/{id}", userHandler.Delete)

	//promethus part for metrics
	mux.Handle("GET /metrics", promhttp.Handler())

	stack := middleware.Chain(
		mux,
		middleware.Recovery(logger),
		middleware.Logger(logger),
		middleware.RequestID,
	)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: stack,
	}

	logger.Info("server starting", slog.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server failed", slog.Any("error", err))
		os.Exit(1)
	}
}
