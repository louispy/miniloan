CREATE TABLE installments (
    id uuid NOT NULL PRIMARY KEY,
    loan_id uuid NOT NULL,
    amount BIGINT NOT NULL,
    due_date TIMESTAMP NOT NULL,
    week INT NOT NULL,
    status INT NOT NULL,
    paid_at TIMESTAMP NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_installments_loan_status_duedate
    ON installments (loan_id, status, due_date);