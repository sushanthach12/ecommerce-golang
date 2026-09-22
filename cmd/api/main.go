package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sushanthach12/ecom-go/internal/config"
	"github.com/sushanthach12/ecom-go/internal/database"
)

func main() {
	cfg := config.MustLoad()

	loggerHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource: true,
		Level: slog.LevelDebug,
	})
	logger := slog.New(loggerHandler)

	// Initialize the database connection
	db, err := database.Connect(cfg.DatabaseUrl)
	if err != nil {
		logger.Error("Failed to connect to database:", "error", err)
		os.Exit(1)
	}

	api := application{
		config: cfg,
		logger: logger,
		db:     db,
	}

	// Listen for interrupt (Ctrl+C) and termination signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := api.run(ctx, api.mount()); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
