package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/QosmuratSamat0/product-service/internal/config"
	httpHandler "github.com/QosmuratSamat0/product-service/internal/handler/http"
	"github.com/QosmuratSamat0/product-service/internal/lib/metrics"
	productRepo "github.com/QosmuratSamat0/product-service/internal/repository/product/postgres"
	productUseCase "github.com/QosmuratSamat0/product-service/internal/usecase/product"
)

type App struct {
	cfg    *config.Config
	router *chi.Mux
	dbPool *pgxpool.Pool
}

func New(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	// Repositories
	prodRepo := productRepo.New(dbPool)

	// UseCases
	prodUC := productUseCase.NewUseCase(prodRepo)

	// HTTP
	r := chi.NewRouter()

	// Metrics middleware
	r.Use(metrics.Middleware)

	// Metrics endpoint for Prometheus
	r.Handle("/metrics", promhttp.Handler())

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	prodModule := httpHandler.NewProductModule(prodUC, cfg.JWTSecret)
	prodModule.RegisterRoutes(r)

	return &App{
		cfg:    cfg,
		router: r,
		dbPool: dbPool,
	}, nil
}

func (a *App) Run() {
	defer a.dbPool.Close()

	srv := &http.Server{
		Addr:    a.cfg.HTTPServer.Address,
		Handler: a.router,
	}

	go func() {
		log.Printf("Starting product-service on %s", a.cfg.HTTPServer.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
