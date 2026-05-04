package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/QosmuratSamat0/payment-service/internal/broker/rabbitmq"
	"github.com/QosmuratSamat0/payment-service/internal/client/order"
	"github.com/QosmuratSamat0/payment-service/internal/client/user"
	"github.com/QosmuratSamat0/payment-service/internal/config"
	httpHandler "github.com/QosmuratSamat0/payment-service/internal/handler/http"
	"github.com/QosmuratSamat0/payment-service/internal/provider"
	"github.com/QosmuratSamat0/payment-service/internal/repository/payment/postgres"
	"github.com/QosmuratSamat0/payment-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	cfg            *config.Config
	log            *slog.Logger
	router         *chi.Mux
	paymentService *service.PaymentService
	dbPool         *pgxpool.Pool
	closers        []interface{ Close() }
}

func New(cfg *config.Config) (*App, error) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	repo := postgres.New(dbPool)
	orderClient := order.NewClient(cfg.OrderServiceURL)
	userClient := user.NewClient(cfg.UserServiceURL)

	providers := map[string]provider.PaymentProvider{
		"mock": provider.NewMockProvider(),
	}

	// Инициализируем брокер
	rmqBroker, err := rabbitmq.NewBroker(cfg.RabbitmqURL)
	if err != nil {
		log.Error("failed to connect to rabbitmq", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	paymentService := service.NewPaymentService(repo, rmqBroker, providers, orderClient, userClient)
	var closers []interface{ Close() }
	closers = append(closers, rmqBroker)
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	r.Handle("/metrics", promhttp.Handler())

	paymentModule := httpHandler.NewPaymentModule(paymentService)
	paymentModule.RegisterRoutes(r)

	return &App{
		cfg:            cfg,
		log:            log,
		router:         r,
		paymentService: paymentService,
		dbPool:         dbPool,
		closers:        closers,
	}, nil
}

func (a *App) Run() {
	defer a.dbPool.Close()
	defer func() {
		for _, closer := range a.closers {
			closer.Close()
		}
	}()

	srv := &http.Server{
		Addr:         a.cfg.HTTPServer.Address,
		Handler:      a.router,
		ReadTimeout:  a.cfg.HTTPServer.Timeout,
		WriteTimeout: a.cfg.HTTPServer.Timeout,
		IdleTimeout:  a.cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		log.Printf("Starting payment-service on %s", a.cfg.HTTPServer.Address)
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
