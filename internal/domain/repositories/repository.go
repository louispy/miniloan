package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/louispy/miniloan/internal/domain/models"
)

type LoansRepository interface {
	Create(ctx context.Context, loan models.Loan) (uuid.UUID, error)
	GetById(ctx context.Context, loanId uuid.UUID) (*models.Loan, error)
}

type BorrowersRepository interface {
	Create(ctx context.Context, borrower models.Borrower) (uuid.UUID, error)
}

type InstallmentsRepository interface {
	CreateMany(ctx context.Context, installments []models.Installment) error
	GetSumByLoanIdAndStatus(ctx context.Context, loanId uuid.UUID, status int) (int64, error)
	GetLateCountByLoanId(ctx context.Context, loanId uuid.UUID, cutoffTime time.Time) (int, error)
}
