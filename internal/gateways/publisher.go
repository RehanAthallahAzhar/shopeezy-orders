package gateway

import (
	"context"

	"github.com/RehanAthallahAzhar/shopeezy-orders/internal/models"
)

type EventPublisher interface {
	PublishOrderCreated(ctx context.Context, event models.OrderCreatedEvent) error
	// Mungkin ada fungsi lain di sini nanti, seperti PublishOrderCanceled
}
