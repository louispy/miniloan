package services

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/louispy/miniloan/internal/constants"
	"github.com/louispy/miniloan/internal/database"
	"github.com/louispy/miniloan/internal/domain/models"
	"github.com/louispy/miniloan/internal/domain/repositories"
)

type loanService struct {
	loansRepo        repositories.LoansRepository
	installmentsRepo repositories.InstallmentsRepository
	txManager        database.TxManager
	weeks            int
	interestRate     float64
}

type LoanServiceOpts struct {
	LoansRepo        repositories.LoansRepository
	InstallmentsRepo repositories.InstallmentsRepository
	TxManager        database.TxManager
}

func NewLoanService(opts LoanServiceOpts) LoanService {
	return &loanService{
		loansRepo:        opts.LoansRepo,
		installmentsRepo: opts.InstallmentsRepo,
		txManager:        opts.TxManager,
		weeks:            constants.LOAN_WEEKS,
		interestRate:     constants.LOAN_INTEREST_RATE,
	}
}

func (s loanService) Create(ctx context.Context, loanInput CreateLoanInput) (*CreateLoanOutput, error) {
	now := time.Now()
	interestRate := s.interestRate
	weeks := s.weeks
	// for scope simplicity, loans are approved and disbursed immediately
	loanId := uuid.New()
	loan := models.Loan{
		Id:           loanId,
		Principal:    loanInput.Principal,
		InterestRate: interestRate,
		Weeks:        int64(weeks),
		Status:       constants.LOAN_STATUS_APPROVED,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	principal := float64(loanInput.Principal)
	total := principal + (principal * interestRate)
	repaymentAmountPerWeek := total / float64(weeks)

	// ASSUME rounding to nearest whole integer
	repaymentAmountPerWeekRounded := math.Round(repaymentAmountPerWeek)

	// Assume that the final week repayment amount should fill the rounding gap
	remainderAmount := total - (repaymentAmountPerWeekRounded * float64(weeks-1))

	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	installments := []models.Installment{}
	for i := range weeks {
		amount := repaymentAmountPerWeekRounded
		if i == weeks-1 {
			amount = remainderAmount
		}
		dueDate := currentDate.AddDate(0, 0, (int(i)+1)*7)
		installments = append(installments, models.Installment{
			Id:        uuid.New(),
			LoanId:    loanId,
			Amount:    amount,
			DueDate:   dueDate,
			Week:      int(i + 1), // assuming week is 1-indexed
			Status:    constants.INSTALLMENT_STATUS_UNPAID,
			IsFinal:   false,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	txCtx, err := s.txManager.Begin(ctx)
	if err != nil {
		log.Printf("error starting transaction: %s", err.Error())
		return nil, err
	}
	defer func(cause error) {
		if cause == nil {
			return
		}
		if err = s.txManager.Rollback(txCtx); err != nil {
			log.Printf("error rolling back transaction: %s", err.Error())
		}
	}(err)

	_, err = s.loansRepo.Create(txCtx, loan)
	if err != nil {
		log.Printf("Error creating loan: %v\n", err.Error())
		return nil, err
	}

	err = s.installmentsRepo.CreateMany(txCtx, installments)
	if err != nil {
		log.Printf("Error creating installments: %v\n", err.Error())
		return nil, err
	}

	if err := s.txManager.Commit(txCtx); err != nil {
		log.Printf("Error committing transaction: %v\n", err.Error())
		return nil, err
	}

	return &CreateLoanOutput{
		Id: loanId.String(),
	}, nil

}
