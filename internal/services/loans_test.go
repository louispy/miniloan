package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/louispy/miniloan/internal/constants"
	"github.com/louispy/miniloan/internal/custerr"
	database_mocks "github.com/louispy/miniloan/internal/database/mocks"
	"github.com/louispy/miniloan/internal/domain/models"
	repo_mocks "github.com/louispy/miniloan/internal/domain/repositories/mocks"
)

func TestLoanService_Create(t *testing.T) {
	err := errors.New("oops")

	okBegin := func() (context.Context, error) { return context.Background(), nil }
	okLoanCreate := func() (uuid.UUID, error) { return uuid.New(), nil }
	okCreateMany := func() error { return nil }
	okCommit := func() error { return nil }
	okRollback := func() error { return nil }

	tests := []struct {
		name             string
		loansRepo        repo_mocks.MockLoansRepo
		installmentsRepo repo_mocks.MockInstallmentRepo
		txManager        database_mocks.MockTxManager
		wantErr          error
	}{
		{
			name:             "ok",
			loansRepo:        repo_mocks.MockLoansRepo{CreateFn: okLoanCreate},
			installmentsRepo: repo_mocks.MockInstallmentRepo{CreateManyFn: okCreateMany},
			txManager: database_mocks.MockTxManager{
				BeginFn:  okBegin,
				CommitFn: okCommit,
			},
		},
		{
			name: "begin fails",
			txManager: database_mocks.MockTxManager{
				BeginFn: func() (context.Context, error) { return nil, err },
			},
			wantErr: err,
		},
		{
			name: "loans repo create fails",
			loansRepo: repo_mocks.MockLoansRepo{
				CreateFn: func() (uuid.UUID, error) { return uuid.Nil, err },
			},
			txManager: database_mocks.MockTxManager{
				BeginFn:    okBegin,
				RollbackFn: okRollback,
			},
			wantErr: err,
		},
		{
			name:      "installments create many fails",
			loansRepo: repo_mocks.MockLoansRepo{CreateFn: okLoanCreate},
			installmentsRepo: repo_mocks.MockInstallmentRepo{
				CreateManyFn: func() error { return err },
			},
			txManager: database_mocks.MockTxManager{
				BeginFn:    okBegin,
				RollbackFn: okRollback,
			},
			wantErr: err,
		},
		{
			name:             "commit fails",
			loansRepo:        repo_mocks.MockLoansRepo{CreateFn: okLoanCreate},
			installmentsRepo: repo_mocks.MockInstallmentRepo{CreateManyFn: okCreateMany},
			txManager: database_mocks.MockTxManager{
				BeginFn:    okBegin,
				CommitFn:   func() error { return err },
				RollbackFn: okRollback,
			},
			wantErr: err,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewLoanService(LoanServiceOpts{
				LoansRepo:        tt.loansRepo,
				InstallmentsRepo: tt.installmentsRepo,
				TxManager:        tt.txManager,
				BusinessTZ:       time.UTC,
			})

			out, err := svc.Create(context.Background(), CreateLoanInput{})

			if err != tt.wantErr {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if out != nil {
					t.Fatalf("expected nil output on error, got %+v", out)
				}
				return
			}
			if out == nil {
				t.Fatal("expected output, got nil")
			}
			if out.Id == "" {
				t.Fatal("expected non-empty loan id")
			}
			if _, err := uuid.Parse(out.Id); err != nil {
				t.Fatalf("output id is not a valid uuid: %v", err)
			}
			if len(out.Installments) != constants.LOAN_WEEKS {
				t.Fatalf("installments len = %d, want %d", len(out.Installments), constants.LOAN_WEEKS)
			}
			wantAmount := int64(5500000) / int64(constants.LOAN_WEEKS)
			for i, ins := range out.Installments {
				if ins.Week != i+1 {
					t.Errorf("installment[%d].Week = %d, want %d", i, ins.Week, i+1)
				}
				if ins.Amount != wantAmount {
					t.Errorf("installment[%d].Amount = %d, want %d", i, ins.Amount, wantAmount)
				}
				if ins.Status != constants.INSTALLMENT_STATUS_UNPAID {
					t.Errorf("installment[%d].Status = %d, want UNPAID", i, ins.Status)
				}
				if ins.LoanId != out.Id {
					t.Errorf("installment[%d].LoanId = %s, want %s", i, ins.LoanId, out.Id)
				}
			}
		})
	}
}

func TestLoanService_GetOutstanding(t *testing.T) {
	err := errors.New("oops")

	okLoanGet := func() (*models.Loan, error) { return &models.Loan{}, nil }

	tests := []struct {
		name             string
		loansRepo        repo_mocks.MockLoansRepo
		installmentsRepo repo_mocks.MockInstallmentRepo
		wantErr          error
		wantAmount       int64
	}{
		{
			name: "ok",
			loansRepo: repo_mocks.MockLoansRepo{
				GetByIdFn: okLoanGet,
			},
			installmentsRepo: repo_mocks.MockInstallmentRepo{
				GetSumByLoanIdAndStatusFn: func() (int64, error) { return 550000, nil },
			},
			wantAmount: 550000,
		},
		{
			name: "ok zero outstanding",
			loansRepo: repo_mocks.MockLoansRepo{
				GetByIdFn: okLoanGet,
			},
			installmentsRepo: repo_mocks.MockInstallmentRepo{
				GetSumByLoanIdAndStatusFn: func() (int64, error) { return 0, nil },
			},
			wantAmount: 0,
		},
		{
			name: "loan not found",
			loansRepo: repo_mocks.MockLoansRepo{
				GetByIdFn: func() (*models.Loan, error) { return nil, custerr.ErrDataNotFound },
			},
			wantErr: custerr.ErrDataNotFound,
		},
		{
			name: "loan repo fails",
			loansRepo: repo_mocks.MockLoansRepo{
				GetByIdFn: func() (*models.Loan, error) { return nil, err },
			},
			wantErr: err,
		},
		{
			name: "installments sum fails",
			loansRepo: repo_mocks.MockLoansRepo{
				GetByIdFn: okLoanGet,
			},
			installmentsRepo: repo_mocks.MockInstallmentRepo{
				GetSumByLoanIdAndStatusFn: func() (int64, error) { return 0, err },
			},
			wantErr: err,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewLoanService(LoanServiceOpts{
				LoansRepo:        tt.loansRepo,
				InstallmentsRepo: tt.installmentsRepo,
				BusinessTZ:       time.UTC,
			})

			out, err := svc.GetOutstanding(context.Background(), GetOutstandingInput{LoanId: uuid.New()})

			if err != tt.wantErr {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if out != nil {
					t.Fatalf("expected nil output on error, got %+v", out)
				}
				return
			}
			if out == nil {
				t.Fatal("expected output, got nil")
			}
			if out.Amount != tt.wantAmount {
				t.Errorf("Amount = %d, want %d", out.Amount, tt.wantAmount)
			}
		})
	}
}
