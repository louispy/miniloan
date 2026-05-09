package services

import (
	"time"

	"github.com/google/uuid"
)

type CreateLoanInput struct {
	BorrowerId uuid.UUID
}

type GetOutstandingInput struct {
	Timestamp time.Time
	LoanId    uuid.UUID
}

type IsDeliquentInput struct {
	Timestamp time.Time
	LoanId    uuid.UUID
}

type MakePaymentInput struct {
	Timestamp time.Time
	Amount    int64
}
