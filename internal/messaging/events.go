package messaging

import "time"

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// OrderStatusChangedEvent event saat order status berubah
type OrderStatusChangedEvent struct {
	OrderID      string      `json:"order_id"`
	UserID       string      `json:"user_id"`
	Email        string      `json:"email"`
	OldStatus    OrderStatus `json:"old_status"`
	NewStatus    OrderStatus `json:"new_status"`
	TotalAmount  float64     `json:"total_amount"`
	ProductCount int         `json:"product_count"`
	ChangedAt    time.Time   `json:"changed_at"`
}
