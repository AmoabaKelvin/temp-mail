-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS attachments (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    filename varchar(255) NOT NULL,
    content_type varchar(255) NOT NULL,
    file_location TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_attachments_message_id ON attachments (message_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS attachments;
DROP INDEX IF EXISTS idx_attachments_message_id;
-- +goose StatementEnd
