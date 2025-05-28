-- +goose Up
-- +goose StatementBegin
ALTER TABLE messages DROP COLUMN receiver_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE messages ADD COLUMN receiver_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd
