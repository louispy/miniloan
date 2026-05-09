package services

import (
	"context"
)

type LoanService interface {
	Create(ctx context.Context, loan CreateLoanInput) (*CreateLoanOutput, error)
	GetOutstanding(ctx context.Context, inp GetOutstandingInput) (*GetOutstandingOutput, error)
	GetInstallments(ctx context.Context, inp GetInstallmentsInput) (*GetInstallmentsOutput, error)
	IsDeliquent(ctx context.Context, inp IsDeliquentInput) (*IsDeliquentOutput, error)
	MakePayment(ctx context.Context, inp MakePaymentInput) (*MakePaymentOutput, error)
}
