-- name: CreateOrder :one
INSERT INTO orders (
    id,
    user_id,
    total_price,
    shipping_address,
    shipping_method,
    payment_method
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetOrderByID :one
SELECT 
    id, 
    user_id, 
    total_price,
    shipping_address, 
    shipping_method, 
    shipping_tracking_code, 
    payment_method, 
    payment_gateway_id, 
    order_status, 
    created_at,
    updated_at
FROM orders 
WHERE id = $1 AND order_status != 'CANCELED' 
LIMIT 1;

-- name: GetOrdersByUserID :many
SELECT 
    id, 
    user_id, 
    total_price,
    shipping_address, 
    shipping_method, 
    shipping_tracking_code, 
    payment_method, 
    payment_gateway_id, 
    order_status, 
    created_at,
    updated_at
FROM orders
WHERE user_id = $1 AND order_status != 'CANCELED'
ORDER BY created_at DESC;

-- name: GetItemsForRestock :many
SELECT product_id, quantity FROM order_items
WHERE order_id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders
SET order_status = $2, updated_at = NOW()
WHERE id = $1 AND order_status != 'CANCELED'
RETURNING *;

-- name: CancelOrder :one
UPDATE orders
SET
    order_status = 'CANCELED',
    updated_at = NOW()
WHERE
    id = $1
    AND user_id = $2
    AND order_status = 'PENDING_PAYMENT' 
RETURNING *;