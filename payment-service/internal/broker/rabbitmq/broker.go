package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewBroker(url string) (*Broker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		"payment_events", // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &Broker{conn: conn, ch: ch}, nil
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

	return b.ch.PublishWithContext(ctx,
		"payment_events",   // exchange
		"payment.success", // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (b *Broker) Close() {
	b.ch.Close()
	b.conn.Close()
}
