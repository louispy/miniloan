package services

import (
	"context"
)

type LoanService interface {
	Create(ctx context.Context, loan CreateLoanInput) (*CreateLoanOutput, error)
}
