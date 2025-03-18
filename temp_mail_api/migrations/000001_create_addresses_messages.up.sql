CREATE TABLE addresses (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    from_address TEXT NOT NULL,
    to_address_id INT NOT NULL,
    subject TEXT,
    body TEXT,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (to_address_id) REFERENCES addresses(id) ON DELETE CASCADE
);
