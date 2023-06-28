-- +migrate Up
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    phone_number VARCHAR(16) NOT NULL UNIQUE,
    email VARCHAR(64) NULL,
    password VARCHAR(256) NOT NULL,
    first_name VARCHAR(128) NULL,
    last_name VARCHAR(128) NULL,
    display_name VARCHAR(128) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT(FALSE),
    is_admin BOOLEAN NOT NULL DEFAULT(FALSE),
    is_superuser BOOLEAN NOT NULL DEFAULT(FALSE),
    created_at DATETIME NOT NULL
);
-- +migrate Down
DROP TABLE users;
