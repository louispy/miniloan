package services

import (
	"context"
	"log"
	"time"

	"github.com/louispy/miniloan/internal/constants"
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

func (s loanService) Create(ctx context.Context, loanInput CreateLoanInput) (*CreateLoanOutput, error) {
	now := time.Now()
	loan := models.Loan{
		Principal:    loanInput.Principal,
		InterestRate: loanInput.InterestRate,
		Weeks:        constants.LOAN_WEEKS,
		Status:       constants.LOAN_STATUS_APPROVED,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	loanId, err := s.loansRepo.Create(ctx, loan)
	if err != nil {
		log.Printf("Error creating loan: %v\n", err.Error())
		return nil, err
	}

	return &CreateLoanOutput{
		Id: loanId.String(),
	}, nil
}
