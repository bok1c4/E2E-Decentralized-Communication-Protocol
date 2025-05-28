-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- optionally recreate old version here if needed
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

