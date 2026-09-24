package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sushanthach12/ecom-go/internal/config"
	"github.com/sushanthach12/ecom-go/internal/database"
	"github.com/sushanthach12/ecom-go/internal/health"
	"github.com/sushanthach12/ecom-go/internal/inventory"
	"github.com/sushanthach12/ecom-go/internal/orders"
	"github.com/sushanthach12/ecom-go/internal/products"
)

type application struct {
	config config.ConfigVars
	logger *slog.Logger
	db     *database.DB
}

func (app *application) mount() http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.ClientIPFromRemoteAddr) // pick one ClientIPFrom* based on your infra, see below
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	// Health
	health.Register(router)

	// Routes
	productAdapters := products.Register(router, app.db, app.logger)
	inventoryAdapters := inventory.Register(router, app.db, app.logger, productAdapters.Adapter)
	orders.Register(router, app.db, app.logger, productAdapters.Adapter, inventoryAdapters.Adapter)

	return router
}

func (app *application) run(ctx context.Context, handler http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server listening at port %s", app.config.Port)

	// Channel to catch errors from ListenAndServe
	serverErr := make(chan error, 1)

	go func() {
		app.logger.Info("Server started", "addr", app.config.Port)
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server error: %w", err)
		}
	case <-ctx.Done():
		app.logger.Info("Shutdown signal received, shutting down server...")

		// Give in-flight requests time to finish
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		app.logger.Info("Server stopped")
	}

	return nil
}
