CREATE TABLE IF NOT EXISTS shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route TEXT NOT NULL,
    price NUMERIC NOT NULL,
    status TEXT NOT NULL DEFAULT 'CREATED',
    customer_id UUID NOT NULL REFERENCES customers(id),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
