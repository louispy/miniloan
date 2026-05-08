package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/louispy/miniloan/internal/domain/models"
	"github.com/louispy/miniloan/internal/domain/repositories"
)

type loanService struct {
	loansRepo repositories.LoansRepository
}

type LoanServiceOpts struct {
	LoansRepo repositories.LoansRepository
}

func NewLoanService(opts LoanServiceOpts) LoanService {
	return &loanService{
		loansRepo: opts.LoansRepo,
	}
}

func (s loanService) Create(ctx context.Context, loan models.Loan) (uuid.UUID, error) {
	return uuid.Nil, nil
}
