package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/louispy/miniloan/internal/domain/models"
)

type LoansRepository interface {
	Create(ctx context.Context, loan models.Loan) (uuid.UUID, error)
}

type BorrowersRepository interface {
	Create(ctx context.Context, borrower models.Borrower) (uuid.UUID, error)
}

type InstallmentsRepository interface {
	CreateMany(ctx context.Context, installments []models.Installment) error
}
