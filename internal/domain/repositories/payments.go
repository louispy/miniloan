package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/louispy/miniloan/internal/domain/models"
	"github.com/louispy/miniloan/internal/utils"
)

type defaultPaymentsRepository struct {
	db *sqlx.DB
}

type PaymentRepoOpts struct {
	DB *sqlx.DB
}

func NewPaymentsRepository(opts PaymentRepoOpts) PaymentsRepository {
	return &defaultPaymentsRepository{
		db: opts.DB,
	}
}

const insertPaymentQuery = `
	INSERT INTO
		payments (
			id,
			amount,
			installment_id,
			created_by,
			updated_by,
			created_at,
			updated_at
		)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)
`

func (r defaultPaymentsRepository) Create(ctx context.Context, payment models.Payment) (uuid.UUID, error) {
	var err error

	if payment.Id == uuid.Nil {
		payment.Id = uuid.New()
	}
	args := []any{
		payment.Id,
		payment.Amount,
		payment.InstallmentId,
		payment.CreatedBy,
		payment.UpdatedBy,
		payment.CreatedAt,
		payment.UpdatedAt,
	}

	tx := utils.SqlxTxFromCtx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, insertPaymentQuery, args...)
	} else {
		_, err = r.db.ExecContext(ctx, insertPaymentQuery, args...)
	}

	return payment.Id, err
}
