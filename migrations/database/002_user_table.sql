-- +goose Up
CREATE TABLE IF NOT EXISTS users (
                       id           BIGINT PRIMARY KEY,
                       username     TEXT,
                       is_admin     BOOLEAN NOT NULL DEFAULT FALSE,
                       created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS users;