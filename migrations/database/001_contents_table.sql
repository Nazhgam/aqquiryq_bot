-- +goose Up
CREATE TABLE IF NOT EXISTS contents (
                                        id            BIGSERIAL PRIMARY KEY,
                                        title         TEXT NOT NULL,
                                        canva_url     TEXT NOT NULL,
                                        class         INTEGER NOT NULL CHECK (class > 0),
                                        quarter       INTEGER NOT NULL CHECK (quarter BETWEEN 1 AND 4),
                                        lesson_number INTEGER NOT NULL CHECK (class > 0),
                                        is_active     BOOLEAN NOT NULL DEFAULT TRUE,
                                        created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        updated_at    TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS contents;