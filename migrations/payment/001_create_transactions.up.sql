CREATE TABLE IF NOT EXISTS payment_transactions (
    id          BIGSERIAL PRIMARY KEY,
    order_id    BIGINT        NOT NULL,
    amount      NUMERIC(12,2) NOT NULL,
    status      VARCHAR(50)   NOT NULL DEFAULT 'pending',
    method      VARCHAR(100)  NOT NULL,
    created_at  TIMESTAMP     NOT NULL,
    created_by  VARCHAR(255)  NOT NULL,
    updated_at  TIMESTAMP,
    updated_by  VARCHAR(255)
);
