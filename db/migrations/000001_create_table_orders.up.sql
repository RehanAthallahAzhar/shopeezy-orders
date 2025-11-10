CREATE TYPE order_status AS ENUM (
    'PENDING_PAYMENT',
    'PROCESSING',
    'SHIPPED',
    'COMPLETED',
    'CANCELED',
    'REFUNDED'
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL, -- Hanya ID, BUKAN foreign key ke database lain
    order_status order_status NOT NULL DEFAULT 'PENDING_PAYMENT',
    total_price DECIMAL(12, 2) NOT NULL,
    
    -- Snapshot data pengiriman dari User Service
    shipping_address TEXT NOT NULL,
    shipping_method TEXT NOT NULL,
    shipping_tracking_code TEXT,
    
    -- Snapshot data pembayaran dari Payment Service
    payment_method TEXT NOT NULL,
    payment_gateway_id TEXT, -- ID transaksi dari payment gateway
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index untuk mempercepat pencarian pesanan berdasarkan pengguna
CREATE INDEX idx_orders_user_id ON orders(user_id);