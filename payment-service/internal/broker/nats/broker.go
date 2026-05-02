package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type Broker struct {
	nc *nats.Conn
}

func NewBroker(url string) (*Broker, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}
	return &Broker{nc: nc}, nil
}

func (b *Broker) PublishPaymentSuccess(ctx context.Context, orderID string) error {
	payload := map[string]string{
		"order_id": orderID,
		"status":   "success",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return b.nc.Publish("payment.success", body)
}

func (b *Broker) Close() {
	b.nc.Close()
}
