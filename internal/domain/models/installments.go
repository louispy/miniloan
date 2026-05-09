package models

import (
	"time"

	"github.com/google/uuid"
)

type Installment struct {
	Id      uuid.UUID `db:"id"`
	LoanId  uuid.UUID `db:"loan_id"`
	Amount  int64     `db:"amount"`
	DueDate time.Time `db:"due_date"`
	Week    int       `db:"week"`
	Status  int       `db:"status"`
	// IsFinal is true when it is the final week installment record
	IsFinal   bool      `db:"is_final"`
	CreatedBy uuid.UUID `db:"created_by"`
	UpdatedBy uuid.UUID `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
