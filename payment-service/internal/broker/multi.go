package broker

import (
	"context"
	"fmt"
)

type MultiBroker struct {
	brokers []MessageBroker
}

type MessageBroker interface {
	PublishPaymentSuccess(ctx context.Context, orderID string) error
}

func NewMultiBroker(brokers ...MessageBroker) *MultiBroker {
	return &MultiBroker{brokers: brokers}
}

func (m *MultiBroker) PublishPaymentSuccess(ctx context.Context, orderID string) error {
	for _, b := range m.brokers {
		if err := b.PublishPaymentSuccess(ctx, orderID); err != nil {
			return fmt.Errorf("failed to publish to one of the brokers: %w", err)
		}
	}
	return nil
}
