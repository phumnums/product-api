CREATE TABLE IF NOT EXISTS products (
    
    id VARCHAR(26) PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    sale_price NUMERIC(10, 2) NULL,
    price NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT products_price_positive CHECK (price > 0),
    CONSTRAINT products_sale_price_positive CHECK (
        sale_price IS NULL 
        OR sale_price >= 0
    )
);