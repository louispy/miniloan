package repo_mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/louispy/miniloan/internal/domain/models"
)

type MockLoansRepo struct {
	CreateFn  func() (uuid.UUID, error)
	GetByIdFn func() (*models.Loan, error)
}

func (r MockLoansRepo) Create(ctx context.Context, loan models.Loan) (uuid.UUID, error) {
	return r.CreateFn()
}

func (r MockLoansRepo) GetById(ctx context.Context, loanId uuid.UUID) (*models.Loan, error) {
	return r.GetByIdFn()
}

type MockInstallmentRepo struct {
	CreateManyFn                func() error
	GetSumByLoanIdAndStatusFn   func() (int64, error)
	GetLateCountByLoanIdFn      func() (int, error)
	GetFirstByLoanIdAndStatusFn func() (*models.Installment, error)
	GetByIdFn                   func() (*models.Installment, error)
	GetByIdForUpdateFn          func() (*models.Installment, error)
	UpdatePaymentFn             func() error
}

func (r MockInstallmentRepo) CreateMany(ctx context.Context, installments []models.Installment) error {
	return r.CreateManyFn()
}

func (r MockInstallmentRepo) GetSumByLoanIdAndStatus(ctx context.Context, loanId uuid.UUID, status int) (int64, error) {
	return r.GetSumByLoanIdAndStatusFn()
}

func (r MockInstallmentRepo) GetLateCountByLoanId(ctx context.Context, loanId uuid.UUID, cutoffTime time.Time) (int, error) {
	return r.GetLateCountByLoanIdFn()
}

func (r MockInstallmentRepo) GetFirstByLoanIdAndStatus(ctx context.Context, loanId uuid.UUID, status int) (*models.Installment, error) {
	return r.GetFirstByLoanIdAndStatusFn()
}

func (r MockInstallmentRepo) GetById(ctx context.Context, installmentId uuid.UUID) (*models.Installment, error) {
	return r.GetByIdFn()
}
func (r MockInstallmentRepo) GetByIdForUpdate(ctx context.Context, installmentId uuid.UUID) (*models.Installment, error) {
	return r.GetByIdForUpdateFn()
}

func (r MockInstallmentRepo) UpdatePayment(ctx context.Context, installmentId uuid.UUID, timestamp time.Time) error {
	return r.UpdatePaymentFn()
}

type MockPaymentsRepo struct {
	CreateFn func() (uuid.UUID, error)
}

func (r MockPaymentsRepo) Create(ctx context.Context, payment models.Payment) (uuid.UUID, error) {
	return r.CreateFn()
}
