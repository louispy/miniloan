package services

import (
	"time"

	"github.com/google/uuid"
)

type CreateLoanInput struct{}

type GetOutstandingInput struct {
	LoanId uuid.UUID
}

type GetInstallmentsInput struct {
	LoanId uuid.UUID
	Status int
}

type IsDeliquentInput struct {
	Timestamp time.Time
	LoanId    uuid.UUID
}

// Per task specification,
// MakePayment should only take amount (+loan_id) as param
// However amount alone does not guarantee idempotency should there be a retry by client
// While being out of scope, adding InstallmentId as a required param is a reasonable assumption
// InstallmentId can be used as unique id for payment, ensuring payment to a specific installment
type MakePaymentInput struct {
	Amount        int64
	LoanId        uuid.UUID
	InstallmentId uuid.UUID
}
