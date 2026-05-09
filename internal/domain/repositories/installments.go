package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
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
	"is_final",
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
	for i, installment := range installments {
		placeholders := []string{}
		for j := range installmentCols {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i*len(installmentCols)+j))
		}
		query += fmt.Sprintf("\n(%s)", strings.Join(placeholders, ","))
		args = append(args, []any{
			installment.Id,
			installment.LoanId,
			installment.Amount,
			installment.DueDate,
			installment.Week,
			installment.Status,
			installment.IsFinal,
			installment.CreatedBy,
			installment.UpdatedBy,
			installment.CreatedAt,
			installment.UpdatedAt,
		}...)
	}
	tx := utils.SqlxTxFromCtx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = r.db.ExecContext(ctx, query, args...)
	}

	return err
}
