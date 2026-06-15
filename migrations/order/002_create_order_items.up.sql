CREATE TABLE IF NOT EXISTS order_items (
    id       BIGSERIAL PRIMARY KEY,
    order_id BIGINT        NOT NULL REFERENCES orders(id),
    sku      VARCHAR(100)  NOT NULL,
    qty      BIGINT        NOT NULL,
    price    NUMERIC(12,2) NOT NULL
);
