package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderCreatedEvent struct {
	OrderID     uuid.UUID `json:"order_id"`
	UserID      uuid.UUID `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"order_status"`
	CreatedAt   time.Time `json:"created_at"`
}
