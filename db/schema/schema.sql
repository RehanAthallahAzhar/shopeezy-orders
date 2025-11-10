CREATE TYPE order_status AS ENUM (
    'PENDING_PAYMENT',
    'PROCESSING',
    'SHIPPED',
    'COMPLETED',
    'CANCELED',
    'REFUNDED'
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL, 
    order_status order_status NOT NULL DEFAULT 'PENDING_PAYMENT',
    total_price DECIMAL(10, 2) NOT NULL,
    
    -- Snapshot of delivery data from User Service
    shipping_address TEXT NOT NULL,
    shipping_method TEXT NOT NULL,
    shipping_tracking_code TEXT,
    
    -- Snapshot of payment data from Payment Service
    payment_method TEXT NOT NULL,
    payment_gateway_id TEXT, -- Transaction ID from the payment gateway
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index untuk mempercepat pencarian pesanan berdasarkan pengguna
CREATE INDEX idx_orders_user_id ON orders(user_id);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE, -- Foreign key ke tabel orders
    product_id UUID NOT NULL, 
    seller_id UUID NOT NULL,
    
    price DECIMAL(10, 2) NOT NULL,
    
    -- Snapshot of product data from Product Service
    product_name TEXT NOT NULL,
    quantity INT NOT NULL,
    "description" TEXT,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index untuk mempercepat pencarian item berdasarkan order_id
CREATE INDEX idx_order_items_order_id ON order_items(order_id);