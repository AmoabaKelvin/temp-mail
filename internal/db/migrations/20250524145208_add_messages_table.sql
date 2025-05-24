-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    from_address varchar(255) NOT NULL,
    to_address_id INT NOT NULL REFERENCES addresses(id) ON DELETE CASCADE,
    subject varchar(255),
    headers JSONB,
    body text,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_messages_to_address_id ON messages (to_address_id);
CREATE INDEX IF NOT EXISTS idx_messages_received_at ON messages (received_at);
CREATE INDEX IF NOT EXISTS idx_messages_read_at ON messages (read_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
DROP INDEX IF EXISTS idx_messages_to_address_id;
DROP INDEX IF EXISTS idx_messages_received_at;
DROP INDEX IF EXISTS idx_messages_read_at;
-- +goose StatementEnd
