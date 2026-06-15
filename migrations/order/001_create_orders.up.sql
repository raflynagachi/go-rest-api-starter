CREATE TABLE IF NOT EXISTS orders (
    id           BIGSERIAL PRIMARY KEY,
    user_email   VARCHAR(255)  NOT NULL,
    status       VARCHAR(50)   NOT NULL DEFAULT 'pending',
    total_amount NUMERIC(12,2) NOT NULL,
    created_at   TIMESTAMP     NOT NULL,
    created_by   VARCHAR(255)  NOT NULL,
    updated_at   TIMESTAMP,
    updated_by   VARCHAR(255),
    deleted_at   TIMESTAMP,
    deleted_by   VARCHAR(255)
);
