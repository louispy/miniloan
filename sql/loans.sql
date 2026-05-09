CREATE TABLE loans (
    id uuid NOT NULL PRIMARY KEY,
    principal BIGINT NOT NULL,
    interest_rate DECIMAL(18,2) NOT NULL,
    weeks INT NOT NULL,
    status INT NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);