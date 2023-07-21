-- +migrate Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(16) NOT NULL UNIQUE,
    email VARCHAR(64),
    password VARCHAR(256) NOT NULL,
    first_name VARCHAR(128),
    last_name VARCHAR(128),
    display_name VARCHAR(128) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    is_superuser BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL
);
-- +migrate Down
DROP TABLE users;
