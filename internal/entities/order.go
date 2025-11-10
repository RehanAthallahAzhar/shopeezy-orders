package entities

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	UserID               uuid.UUID `json:"user_id" db:"user_id"`
	OrderStatus          string    `json:"order_status" db:"order_status"`
	TotalPrice           float64   `json:"total_amount" db:"total_price"`
	ShippingAddress      string    `json:"shipping_address" db:"shipping_address"`
	ShippingMethod       string    `json:"shipping_method" db:"shipping_method"`
	ShippingTrackingCode string    `json:"shipping_tracking_code" db:"shipping_tracking_code"`
	PaymentMethod        string    `json:"payment_method" db:"payment_method"`
	PaymentGatewayID     string    `json:"payment_gateway_id" db:"payment_gateway_id"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

type OrderDetails struct {
	Order Order
	Items []OrderItem
}
