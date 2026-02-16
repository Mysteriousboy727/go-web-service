package main

import (
	"log/slog"
	"net/http"
	"os"

	"go-industry-server/configs"
	"go-industry-server/internal/handler"
	"go-industry-server/internal/middleware"
	"go-industry-server/internal/repository"
	"go-industry-server/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg := configs.Load()

	userRepo := repository.NewInMemoryUserRepository()
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
