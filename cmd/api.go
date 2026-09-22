package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sushanthach12/ecom-go/internal/config"
	"github.com/sushanthach12/ecom-go/internal/products"
)

type application struct {
	config config.ConfigVars
	logger *slog.Logger
	db     *sql.DB
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

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("All good"))
	})

	// Routes
	products.Register(router, app.db, app.logger)

	return router
}

func (app *application) run(handler http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server listening at port %s", app.config.Port)

	return srv.ListenAndServe()
}
