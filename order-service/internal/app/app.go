package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	productClient "github.com/QosmuratSamat/order-service/internal/client/product"
	"github.com/QosmuratSamat/order-service/internal/config"
	domain "github.com/QosmuratSamat/order-service/internal/domain/order"
	httpHandler "github.com/QosmuratSamat/order-service/internal/handler/http"
	"github.com/QosmuratSamat/order-service/internal/lib/metrics"
	orderRepo "github.com/QosmuratSamat/order-service/internal/repository/order/postgres"
	orderService "github.com/QosmuratSamat/order-service/internal/service/order"
	orderUseCase "github.com/QosmuratSamat/order-service/internal/usecase/order"
)

type App struct {
	cfg    *config.Config
	router *chi.Mux
	dbPool *pgxpool.Pool
	nc     *nats.Conn
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

	prodClient := productClient.New(cfg.ProductServiceURL)
	ordRepo := orderRepo.New(dbPool)
	ordService := orderService.NewService(prodClient)
	ordUC := orderUseCase.NewUseCase(ordRepo, ordService)

	r := chi.NewRouter()

	// Metrics middleware
	r.Use(metrics.Middleware)

	// Metrics endpoint for Prometheus
	r.Handle("/metrics", promhttp.Handler())

	orderModule := httpHandler.NewOrderModule(ordUC, cfg.JWTSecret)
	orderModule.RegisterRoutes(r)

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to nats: %w", err)
	}

	// Subscribe to payment success events to update order status
	_, err = nc.Subscribe("payment.success", func(m *nats.Msg) {
		var data struct {
			OrderID string `json:"order_id"`
			Status  string `json:"status"`
		}
		if err := json.Unmarshal(m.Data, &data); err != nil {
			log.Printf("failed to unmarshal payment success event: %v", err)
			return
		}

		if data.Status == "success" {
			log.Printf("Updating order %s status to paid", data.OrderID)
			if err := ordRepo.UpdateOrderStatus(ctx, data.OrderID, domain.StatusPaid); err != nil {
				log.Printf("failed to update order status for order %s: %v", data.OrderID, err)
			}
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to payment success: %w", err)
	}

	return &App{
		cfg:    cfg,
		router: r,
		dbPool: dbPool,
		nc:     nc,
	}, nil
}

func (a *App) Run() {
	defer a.dbPool.Close()
	defer a.nc.Close()

	srv := &http.Server{
		Addr:         a.cfg.HTTPServer.Address,
		Handler:      a.router,
		ReadTimeout:  a.cfg.HTTPServer.Timeout,
		WriteTimeout: a.cfg.HTTPServer.Timeout,
		IdleTimeout:  a.cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		log.Printf("Starting order-service on %s", a.cfg.HTTPServer.Address)
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
