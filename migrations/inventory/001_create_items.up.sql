CREATE TABLE IF NOT EXISTS items (
    id          BIGSERIAL PRIMARY KEY,
    sku         VARCHAR(100) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    quantity    BIGINT       NOT NULL DEFAULT 0,
    price       NUMERIC(12,2) NOT NULL,
    created_at  TIMESTAMP    NOT NULL,
    created_by  VARCHAR(255) NOT NULL,
    updated_at  TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP,
    deleted_by  VARCHAR(255)
);
