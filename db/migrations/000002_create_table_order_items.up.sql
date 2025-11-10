CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE, -- Foreign key ke tabel orders
    product_id UUID NOT NULL, -- Hanya ID, BUKAN foreign key ke database lain
    seller_id UUID NOT NULL,

    product_name TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    
    quantity INT NOT NULL,
    description TEXT,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index untuk mempercepat pencarian item berdasarkan order_id
CREATE INDEX idx_order_items_order_id ON order_items(order_id);