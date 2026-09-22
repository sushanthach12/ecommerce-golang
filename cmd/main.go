package main

import (
	"log/slog"
	"os"

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

	if err := api.run(api.mount()); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
