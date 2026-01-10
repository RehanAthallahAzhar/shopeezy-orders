package models

import (
	"time"
)

type OrderDetailsResponse struct {
	Order OrderResponse  `json:"order"`
	Items []OrderItemRes `json:"items"`
}

type OrderResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Status          string    `json:"status"`
	TotalPrice      float64   `json:"total_price"`
	ShippingAddress string    `json:"shipping_address"`
	CreatedAt       time.Time `json:"created_at"`
}

type OrderReq struct {
	TotalPrice           float64 `json:"total_price" validate:"required"`
	ShippingAddress      string  `json:"shipping_address" validate:"required"`
	ShippingMethod       string  `json:"shipping_method" validate:"required"`
	PaymentMethod        string  `json:"payment_method" validate:"required"`
	ShippingTrackingCode string  `json:"shipping_tracking_code" validate:"required"`
	PaymentGatewayID     string  `json:"payment_gateway_id" validate:"required"`
}

type OrderDetailReq struct {
	Order OrderReq       `json:"order"`
	Items []OrderItemReq `json:"items"`
}

type UpdateOrderStatusReq struct {
	Status        string `json:"status" validate:"required"`
	CurrentStatus string `json:"current_status" validate:"required"`
}
