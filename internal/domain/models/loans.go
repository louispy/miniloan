package models

import (
	"time"

	"github.com/google/uuid"
)

type Loan struct {
	Id             uuid.UUID `db:"id"`
	Principal      int64     `db:"principal"`
	InterestRate   float64   `db:"interest_rate"`
	PeriodDays     int64     `db:"period_days"` // weekly: 7
	NoInstallments int64     `db:"no_installments"`
	Status         int64     `db:"status"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
