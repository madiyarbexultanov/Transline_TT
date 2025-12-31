CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idn TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
