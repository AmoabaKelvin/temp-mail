-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS addresses;
-- +goose StatementEnd
