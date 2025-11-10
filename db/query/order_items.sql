-- name: CreateOrderItem :one
INSERT INTO order_items (
    id,
    order_id,
    product_id,
    seller_id,
    product_name,
    price,
    quantity,
    "description"
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetOrderItemsByOrderIDs :many
SELECT id, order_id, product_id, seller_id, product_name,price, quantity, "description" FROM order_items
WHERE order_id = ANY(sqlc.arg(order_ids)::uuid[]);

-- name: GetOrderItemsByOrderID :many
SELECT id, order_id, product_id, seller_id, product_name,price, quantity, "description" FROM order_items
WHERE order_id = $1;

-- name: UpdateOrderWithTrackingCode :one
UPDATE orders
SET shipping_tracking_code = $2, order_status = 'SHIPPED', updated_at = NOW()
WHERE id = $1
RETURNING *;