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
	"github.com/nats-io/nats.go"
)

type App struct {
	cfg     *config.Config
	nc      *nats.Conn
	clients map[chan string]bool
	mu      sync.Mutex
}

func New(cfg *config.Config) (*App, error) {
	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:     cfg,
		nc:      nc,
		clients: make(map[chan string]bool),
	}, nil
}

func (a *App) Run() {
	log.Printf("Notification service starting...")
	log.Printf("Listening to NATS at %s", a.cfg.NatsURL)
	log.Printf("Serving SSE at :%s/events", a.cfg.HttpPort)

	// Subscribe to payment success events
	_, err := a.nc.Subscribe("payment.success", func(m *nats.Msg) {
		msg := string(m.Data)
		log.Printf("Received payment success event: %s", msg)
		a.broadcast(msg)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// HTTP server for SSE
	mux := http.NewServeMux()
	mux.HandleFunc("/events", a.handleSSE)

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
	a.nc.Close()
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
