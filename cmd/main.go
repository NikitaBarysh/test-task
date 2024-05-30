package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"test-task/internal/app"
	"test-task/internal/config"
	"test-task/internal/handler"
	"test-task/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Warn("Main: ", err)
		panic(err)
	}

	router := chi.NewRouter()
	newService := service.NewService(cfg)
	newHandler := handler.NewHandler(newService, cfg)
	newHandler.Register(router)

	srv := app.NewServer(cfg.ServerHost, cfg.ServerPort, router)

	if err := srv.Run(); err != nil {
		log.Warn("Main server: Err to start server: ", err)
		panic(err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Warn("Main server: Err to shutdown server: ", err)
		panic(err)
	}

}
