package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/louispy/miniloan/internal/constants"
	"github.com/louispy/miniloan/internal/custerr"
	"github.com/louispy/miniloan/internal/domain/models"
	"github.com/louispy/miniloan/internal/utils"
)

type defaultInstallmentsRepository struct {
	db *sqlx.DB
}

type InstallmentRepoOpts struct {
	DB *sqlx.DB
}

func NewInstallmentsRepository(opts InstallmentRepoOpts) InstallmentsRepository {
	return &defaultInstallmentsRepository{
		db: opts.DB,
	}
}

var installmentCols = []string{
	"id",
	"loan_id",
	"amount",
	"due_date",
	"week",
	"status",
	"paid_at",
	"created_by",
	"updated_by",
	"created_at",
	"updated_at",
}

const insertInstallmentQuery = `
	INSERT INTO
		installments (%s)
	VALUES
`

func (r defaultInstallmentsRepository) CreateMany(ctx context.Context, installments []models.Installment) (err error) {
	args := []any{}
	query := fmt.Sprintf(insertInstallmentQuery, strings.Join(installmentCols, ","))
	values := []string{}
	for i, installment := range installments {
		placeholders := []string{}
		for j := range installmentCols {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i*len(installmentCols)+j+1))
		}
		values = append(values, fmt.Sprintf("\n(%s)", strings.Join(placeholders, ",")))
		args = append(args, []any{
			installment.Id,
			installment.LoanId,
			installment.Amount,
			installment.DueDate,
			installment.Week,
			installment.Status,
			installment.PaidAt,
			installment.CreatedBy,
			installment.UpdatedBy,
			installment.CreatedAt,
			installment.UpdatedAt,
		}...)
	}
	query += strings.Join(values, ",")
	tx := utils.SqlxTxFromCtx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = r.db.ExecContext(ctx, query, args...)
	}

	return err
}

const getSumByLoanIdAndStatusQuery = `
	SELECT
		COALESCE(SUM(amount), 0)
	FROM
		installments
	WHERE
		loan_id = $1
		AND status = $2
`

func (r defaultInstallmentsRepository) GetSumByLoanIdAndStatus(ctx context.Context, loanId uuid.UUID, status int) (int64, error) {
	var (
		total int64
		err   error
	)

	args := []any{
		loanId,
		status,
	}
	tx := utils.SqlxTxFromCtx(ctx)
	if tx != nil {
		err = tx.GetContext(ctx, &total, getSumByLoanIdAndStatusQuery, args...)
	} else {
		err = r.db.GetContext(ctx, &total, getSumByLoanIdAndStatusQuery, args...)
	}

	if err != nil {
		return 0, err
	}

	return total, nil
}

const getLateInstallmentsCountByLoanIdQuery = `
	SELECT
		COUNT(*)
	FROM
		installments
	WHERE
		loan_id = $1
		AND status = $2
		AND due_date < $3
`

func (r defaultInstallmentsRepository) GetLateCountByLoanId(ctx context.Context, loanId uuid.UUID, cutoffTime time.Time) (int, error) {
	var (
		count int
		err   error
	)

	args := []any{
		loanId,
		constants.INSTALLMENT_STATUS_UNPAID,
		cutoffTime,
	}
	tx := utils.SqlxTxFromCtx(ctx)
	if tx != nil {
		err = tx.GetContext(ctx, &count, getLateInstallmentsCountByLoanIdQuery, args...)
	} else {
		err = r.db.GetContext(ctx, &count, getLateInstallmentsCountByLoanIdQuery, args...)
	}

	if err != nil {
		return 0, err
	}

	return count, nil
}

const getInstallmentQuery = `
	SELECT
		%s
	FROM
		installments
	WHERE
		%s
`

func (r defaultInstallmentsRepository) GetFirstByLoanIdAndStatus(ctx context.Context, loanId uuid.UUID, status int) (*models.Installment, error) {
	var err error
	query := fmt.Sprintf(
		getInstallmentQuery,
		strings.Join(installmentCols, ","),
		"loan_id = $1 AND status = $2 ORDER BY week ASC LIMIT 1",
	)

	installment := models.Installment{}
	args := []any{
		loanId,
		status,
	}

	tx := utils.SqlxTxFromCtx(ctx)
	if tx != nil {
		err = tx.GetContext(ctx, &installment, query, args...)
	} else {
		err = r.db.GetContext(ctx, &installment, query, args...)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custerr.ErrDataNotFound
		}
		return nil, err
	}

	return &installment, nil

}

func (r defaultInstallmentsRepository) GetById(ctx context.Context, installmentId uuid.UUID) (*models.Installment, error) {
	var err error
	query := fmt.Sprintf(
		getInstallmentQuery,
		strings.Join(installmentCols, ","),
		"id = $1 LIMIT 1",
	)

	installment := models.Installment{}
	args := []any{
		installmentId,
	}

	tx := utils.SqlxTxFromCtx(ctx)
	if tx != nil {
		err = tx.GetContext(ctx, &installment, query, args...)
	} else {
		err = r.db.GetContext(ctx, &installment, query, args...)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custerr.ErrDataNotFound
		}
		return nil, err
	}

	return &installment, nil
}

const updateInstallmentPaymentQuery = `
	UPDATE
		installments
	SET
		status = $1,
		paid_at = $2,
		updated_at = $2
	WHERE
		id = $3
`

func (r defaultInstallmentsRepository) UpdatePayment(ctx context.Context, installmentId uuid.UUID, timestamp time.Time) error {
	var err error
	args := []any{
		constants.INSTALLMENT_STATUS_PAID,
		timestamp,
		installmentId,
	}
	query := updateInstallmentPaymentQuery
	tx := utils.SqlxTxFromCtx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = r.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (r defaultInstallmentsRepository) GetByIdForUpdate(ctx context.Context, installmentId uuid.UUID) (*models.Installment, error) {
	var err error
	query := fmt.Sprintf(
		getInstallmentQuery,
		strings.Join(installmentCols, ","),
		"id = $1 LIMIT 1 FOR UPDATE",
	)

	installment := models.Installment{}
	args := []any{
		installmentId,
	}

	tx := utils.SqlxTxFromCtx(ctx)
	if tx == nil {
		return nil, custerr.ErrNoTransaction
	}
	err = tx.GetContext(ctx, &installment, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custerr.ErrDataNotFound
		}
		return nil, err
	}

	return &installment, nil
}
