CREATE TABLE borrowers (
    id uuid NOT NULL PRIMARY KEY,
    principal DECIMAL(18,2) NOT NULL,
    interest_rate DECIMAL(18,2) NOT NULL,
    weeks INT NOT NULL,
    status INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);