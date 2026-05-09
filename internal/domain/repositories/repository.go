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
	GetFirstByLoanIdAndStatus(ctx context.Context, loanId uuid.UUID, status int) (*models.Installment, error)
	GetById(ctx context.Context, installmentId uuid.UUID) (*models.Installment, error)
	GetByIdForUpdate(ctx context.Context, installmentId uuid.UUID) (*models.Installment, error)
	UpdatePayment(ctx context.Context, installmentId uuid.UUID, timestamp time.Time) error
}

type PaymentsRepository interface {
	Create(ctx context.Context, payment models.Payment) (uuid.UUID, error)
}
