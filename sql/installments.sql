CREATE TABLE installments (
    id uuid NOT NULL PRIMARY KEY,
    loan_id uuid NOT NULL,
    amount BIGINT NOT NULL,
    due_date TIMESTAMP NOT NULL,
    week INT NOT NULL,
    status INT NOT NULL,
    is_final BOOLEAN NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);