package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/QosmuratSamat0/notification-service/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
)

type App struct {
	cfg     *config.Config
	rmq     *amqp.Connection
	ch      *amqp.Channel
	clients map[chan string]bool
	mu      sync.Mutex
}

func New(cfg *config.Config) (*App, error) {
	conn, err := amqp.Dial(cfg.RabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &App{
		cfg:     cfg,
		rmq:     conn,
		ch:      ch,
		clients: make(map[chan string]bool),
	}, nil
}

func (a *App) Run() {
	log.Printf("Notification service starting...")
	log.Printf("Listening to RabbitMQ at %s", a.cfg.RabbitmqURL)
	log.Printf("Serving SSE at :%s/events", a.cfg.HttpPort)

	err := a.ch.ExchangeDeclare(
		"payment_events", // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	q, err := a.ch.QueueDeclare(
		"notification_service_queue", // name
		false,                        // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	err = a.ch.QueueBind(
		q.Name,            // queue name
		"payment.success", // routing key
		"payment_events",  // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	msgs, err := a.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	go func() {
		for d := range msgs {
			msg := string(d.Body)
			log.Printf("Received payment success event: %s", msg)
			a.broadcast(msg)
		}
	}()

	// HTTP server for SSE
	mux := http.NewServeMux()
	mux.HandleFunc("/events", a.handleSSE)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr:    ":" + a.cfg.HttpPort,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Println("Shutting down notification-service")
	a.ch.Close()
	a.rmq.Close()
}

func (a *App) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan string)

	a.mu.Lock()
	a.clients[messageChan] = true
	a.mu.Unlock()

	defer func() {
		a.mu.Lock()
		delete(a.clients, messageChan)
		a.mu.Unlock()
		close(messageChan)
	}()

	notify := r.Context().Done()

	for {
		select {
		case msg := <-messageChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.(http.Flusher).Flush()
		case <-notify:
			return
		}
	}
}

func (a *App) broadcast(msg string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for clientChan := range a.clients {
		select {
		case clientChan <- msg:
		default:
			log.Println("Skipping slow client")
		}
	}
}
