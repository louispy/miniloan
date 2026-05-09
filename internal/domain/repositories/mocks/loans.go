package repo_mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/louispy/miniloan/internal/domain/models"
)

type MockLoansRepo struct {
	CreateFn func() (uuid.UUID, error)
}

func (r MockLoansRepo) Create(ctx context.Context, loan models.Loan) (uuid.UUID, error) {
	return r.CreateFn()
}

type MockBorrowersRepo struct {
	CreateFn func() (uuid.UUID, error)
}

func (r MockBorrowersRepo) Create(ctx context.Context, borrower models.Borrower) (uuid.UUID, error) {
	return r.CreateFn()
}

type MockInstallmentRepo struct {
	CreateManyFn func() error
}

func (r MockInstallmentRepo) CreateMany(ctx context.Context, installments []models.Installment) error {
	return r.CreateManyFn()
}

// type LoansRepository interface {
// 	Create(ctx context.Context, loan models.Loan) (uuid.UUID, error)
// }

// type BorrowersRepository interface {
// 	Create(ctx context.Context, borrower models.Borrower) (uuid.UUID, error)
// }

// type InstallmentsRepository interface {
// 	CreateMany(ctx context.Context, installments []models.Installment) error
// }
