package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/louispy/miniloan/internal/domain/models"
)

type loanService struct {
}

func NewLoanService() LoanService {
	return &loanService{}
}

func (s loanService) Create(ctx context.Context, loan models.Loan) (uuid.UUID, error) {
	return uuid.Nil, nil
}
