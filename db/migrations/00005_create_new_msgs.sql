-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Polymorphic target: either to a user or to a channel
    receiver_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    channel_id INTEGER REFERENCES channels(id) ON DELETE CASCADE,

    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    CHECK (
        (receiver_id IS NULL AND channel_id IS NOT NULL) OR
        (receiver_id IS NOT NULL AND channel_id IS NULL)
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
-- +goose StatementEnd
