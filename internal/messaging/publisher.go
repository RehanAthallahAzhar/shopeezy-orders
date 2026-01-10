package messaging

import (
	"context"
	"fmt"

	messaging "github.com/RehanAthallahAzhar/tokohobby-messaging-go"
	"github.com/sirupsen/logrus"
)

type EventPublisher struct {
	publisher *messaging.Publisher
	log       *logrus.Logger
}

func NewEventPublisher(rmq *messaging.RabbitMQ, log *logrus.Logger) *EventPublisher {
	return &EventPublisher{
		publisher: messaging.NewPublisher(rmq),
		log:       log,
	}
}

// PublishOrderStatusChanged publish event order status changed
func (p *EventPublisher) PublishOrderStatusChanged(ctx context.Context, event OrderStatusChangedEvent) error {
	// Dynamic routing key based on status
	routingKey := fmt.Sprintf("order.status.%s", event.NewStatus)

	opts := messaging.PublishOptions{
		Exchange:   "order.events",
		RoutingKey: routingKey, // order.status.paid, order.status.shipped, etc
		Mandatory:  false,
		Immediate:  false,
	}

	err := p.publisher.Publish(ctx, opts, event)
	if err != nil {
		p.log.Errorf("Failed to publish order status changed event: %v", err)
		return err
	}

	p.log.Debugf("Published order.status.%s event for order: %s", event.NewStatus, event.OrderID)
	return nil
}
