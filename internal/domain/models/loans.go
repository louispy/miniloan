package models

import (
	"time"

	"github.com/google/uuid"
)

type Loan struct {
	Id           uuid.UUID `db:"id"`
	BorrowerId   uuid.UUID `db:"borrower_id"`
	Principal    int64     `db:"principal"`
	InterestRate float64   `db:"interest_rate"`
	Weeks        int64     `db:"weeks"`
	Status       int64     `db:"status"`
	CreatedBy    uuid.UUID `db:"created_by"`
	UpdatedBy    uuid.UUID `db:"created_by"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
