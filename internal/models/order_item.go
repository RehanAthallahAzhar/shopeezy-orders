package models

import (
	"time"
)

type RedisOrderItem struct {
	Quantity        int       `json:"quantity"`
	CartDescription string    `json:"cart_description "`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type OrderItemReq struct {
	ID          string `json:"id" validate:"required"`
	Quantity    int32  `json:"quantity" validate:"required,min=1"`
	Description string `json:"description"`
}

type OrderItemRes struct {
	ID              string  `json:"id"`
	OrderID         string  `json:"order_id"`
	SellerID        string  `json:"seller_id"`
	SellerName      string  `json:"seller_name"`
	ProductID       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	Quantity        int     `json:"quantity"`
	ProductPrice    float64 `json:"product_price"`
	CartDescription string  `json:"cart_description,omitempty"`
}
