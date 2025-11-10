package entities

type OrderItem struct {
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
