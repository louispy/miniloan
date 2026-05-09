package models

import (
	"time"

	"github.com/google/uuid"
)

type Borrower struct {
	Id        uuid.UUID `db:"id"`
	FullName  string    `db:"full_name"`
	CreatedBy uuid.UUID `db:"created_by"`
	UpdatedBy uuid.UUID `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
