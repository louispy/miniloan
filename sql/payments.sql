CREATE TABLE IF NOT EXISTS payments (
    id uuid NOT NULL PRIMARY KEY,
    amount bigint NOT NULL,
    installment_id uuid NOT NULL UNIQUE,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);