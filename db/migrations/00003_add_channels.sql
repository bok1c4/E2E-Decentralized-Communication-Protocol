-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    is_direct BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE channels;
-- +goose StatementEnd
