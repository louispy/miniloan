package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	Id            uuid.UUID `db:"id"`
	Amount        int64     `db:"amount"`
	InstallmentId uuid.UUID `db:"installment_id"`
	CreatedBy     uuid.UUID `db:"created_by"`
	UpdatedBy     uuid.UUID `db:"updated_by"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
