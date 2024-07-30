CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(320) NOT NULL UNIQUE,
    created_at timestamp NOT NULL DEFAULT NOW(),
    created_by VARCHAR(320) NOT NULL,
    updated_at timestamp,
    updated_by VARCHAR(320),
    deleted_at timestamp,
    deleted_by VARCHAR(320)
);