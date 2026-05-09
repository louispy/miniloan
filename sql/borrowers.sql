CREATE TABLE borrowers (
    id uuid NOT NULL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);