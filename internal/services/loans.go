package services

import (
	"context"
	"errors"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/louispy/miniloan/internal/constants"
	"github.com/louispy/miniloan/internal/custerr"
	"github.com/louispy/miniloan/internal/database"
	"github.com/louispy/miniloan/internal/domain/models"
	"github.com/louispy/miniloan/internal/domain/repositories"
)

type loanService struct {
	loansRepo        repositories.LoansRepository
	installmentsRepo repositories.InstallmentsRepository
	paymentsRepo     repositories.PaymentsRepository
	txManager        database.TxManager
	weeks            int
	interestRate     float64
	principal        int64
	businessTZ       *time.Location
}

type LoanServiceOpts struct {
	LoansRepo        repositories.LoansRepository
	InstallmentsRepo repositories.InstallmentsRepository
	PaymentsRepo     repositories.PaymentsRepository
	TxManager        database.TxManager
	BusinessTZ       *time.Location
}

func NewLoanService(opts LoanServiceOpts) LoanService {
	return &loanService{
		loansRepo:        opts.LoansRepo,
		installmentsRepo: opts.InstallmentsRepo,
		paymentsRepo:     opts.PaymentsRepo,
		txManager:        opts.TxManager,
		weeks:            constants.LOAN_WEEKS,
		interestRate:     constants.LOAN_INTEREST_RATE,
		principal:        constants.LOAN_PRINCIPAL,
		businessTZ:       opts.BusinessTZ,
	}
}

func (s loanService) Create(ctx context.Context, loanInput CreateLoanInput) (*CreateLoanOutput, error) {
	now := time.Now().UTC()
	interestRate := s.interestRate
	weeks := s.weeks
	principal := s.principal
	// for scope simplicity, loans are approved and disbursed immediately
	loanId := uuid.New()
	loan := models.Loan{
		Id:           loanId,
		Principal:    principal,
		InterestRate: interestRate,
		Weeks:        int64(weeks),
		Status:       constants.LOAN_STATUS_APPROVED,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	total := principal + int64(math.Round(float64(principal)*interestRate))
	repaymentAmountPerWeek := total / int64(weeks)

	nowBiz := now.In(s.businessTZ)
	today := time.Date(nowBiz.Year(), nowBiz.Month(), nowBiz.Day(), 0, 0, 0, 0, s.businessTZ)

	installments := []models.Installment{}
	for i := range weeks {
		dueDate := today.AddDate(0, 0, (int(i)+1)*7)
		installments = append(installments, models.Installment{
			Id:        uuid.New(),
			LoanId:    loanId,
			Amount:    repaymentAmountPerWeek,
			DueDate:   dueDate,
			Week:      int(i + 1), // assuming week is 1-indexed
			Status:    constants.INSTALLMENT_STATUS_UNPAID,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	txCtx, err := s.txManager.Begin(ctx)
	if err != nil {
		log.Printf("error starting transaction: %s", err.Error())
		return nil, err
	}
	defer func() {
		if err != nil {
			s.txManager.Rollback(txCtx)
		}
	}()

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

	err = s.txManager.Commit(txCtx)
	if err != nil {
		log.Printf("Error committing transaction: %v\n", err.Error())
		return nil, err
	}

	outputInstallments := []CreateLoanOutputInstallment{}
	for _, installment := range installments {
		outputInstallments = append(outputInstallments, CreateLoanOutputInstallment{
			Week:    installment.Week,
			Amount:  installment.Amount,
			Status:  installment.Status,
			DueDate: installment.DueDate.String(),
			LoanId:  installment.LoanId.String(),
		})
	}
	return &CreateLoanOutput{
		Id: loanId.String(),
	}, nil

}

func (s loanService) GetOutstanding(ctx context.Context, inp GetOutstandingInput) (*GetOutstandingOutput, error) {
	_, err := s.loansRepo.GetById(ctx, inp.LoanId)
	if err != nil {
		if err != custerr.ErrDataNotFound {
			log.Printf("Error retrieving Loan: %v\n", err.Error())
		}
		return nil, err
	}
	amount, err := s.installmentsRepo.GetSumByLoanIdAndStatus(ctx, inp.LoanId, constants.INSTALLMENT_STATUS_UNPAID)
	if err != nil {
		log.Printf("Error retrieving outstanding sum: %v\n", err.Error())
		return nil, err
	}
	return &GetOutstandingOutput{Amount: amount}, nil
}

func (s loanService) IsDeliquent(ctx context.Context, inp IsDeliquentInput) (*IsDeliquentOutput, error) {
	_, err := s.loansRepo.GetById(ctx, inp.LoanId)
	if err != nil {
		if err != custerr.ErrDataNotFound {
			log.Printf("Error retrieving Loan: %v\n", err.Error())
		}
		return nil, err
	}

	nowBiz := inp.Timestamp.In(s.businessTZ)
	todayBiz := time.Date(nowBiz.Year(), nowBiz.Month(), nowBiz.Day(), 0, 0, 0, 0, s.businessTZ)

	lateCount, err := s.installmentsRepo.GetLateCountByLoanId(ctx, inp.LoanId, todayBiz)
	if err != nil {
		log.Printf("Error retrieving outstanding sum: %v\n", err.Error())
		return nil, err
	}

	return &IsDeliquentOutput{
		IsDeliquent: lateCount >= constants.DELIQUENT_LATE_COUNT,
	}, nil
}

func (s loanService) MakePayment(ctx context.Context, inp MakePaymentInput) (*MakePaymentOutput, error) {
	_, err := s.loansRepo.GetById(ctx, inp.LoanId)
	if err != nil {
		if err != custerr.ErrDataNotFound {
			log.Printf("Error retrieving Loan: %v\n", err.Error())
		}
		return nil, err
	}

	txCtx, err := s.txManager.Begin(ctx)
	if err != nil {
		log.Printf("error starting transaction: %s", err.Error())
		return nil, err
	}
	defer func() {
		if err != nil {
			s.txManager.Rollback(txCtx)
		}
	}()

	installment, err := s.installmentsRepo.GetByIdForUpdate(txCtx, inp.InstallmentId)
	if err != nil {
		if err != custerr.ErrDataNotFound {
			log.Printf("Error retrieving installment: %v\n", err.Error())
		}
		return nil, err
	}

	if installment.LoanId != inp.LoanId {
		return nil, errors.New("invalid installment")
	}

	if installment.Status != constants.INSTALLMENT_STATUS_UNPAID {
		return nil, errors.New("invalid installment status")
	}

	if inp.Amount != installment.Amount {
		return nil, errors.New("invalid payment amount")
	}

	now := time.Now().UTC()
	_, err = s.paymentsRepo.Create(txCtx, models.Payment{
		Id:            uuid.New(),
		Amount:        inp.Amount,
		InstallmentId: installment.Id,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		log.Printf("Error creating payment: %v\n", err.Error())
		return nil, err
	}
	err = s.installmentsRepo.UpdatePayment(txCtx, installment.Id, now)
	if err != nil {
		log.Printf("Error updating installment: %v\n", err.Error())
		return nil, err
	}

	err = s.txManager.Commit(txCtx)
	if err != nil {
		log.Printf("Error committing transaction: %v\n", err.Error())
		return nil, err
	}

	return &MakePaymentOutput{InstallmentId: installment.Id.String()}, nil
}
