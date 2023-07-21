-- +migrate Up
CREATE TABLE tokens (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    token VARCHAR(256) NOT NULL,
    is_refresh_token BOOLEAN NOT NULL DEFAULT(FALSE),
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +migrate Down
DROP TABLE tokens;
