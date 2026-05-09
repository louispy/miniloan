package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/louispy/miniloan/internal/custerr"
	"github.com/louispy/miniloan/internal/domain/models"
	"github.com/louispy/miniloan/internal/utils"
)

type defaultLoansRepository struct {
	db *sqlx.DB
}

type LoanRepoOpts struct {
	DB *sqlx.DB
}

func NewLoansRepository(opts LoanRepoOpts) LoansRepository {
	return &defaultLoansRepository{
		db: opts.DB,
	}
}

const insertLoanQuery = `
	INSERT INTO
		loans (
			id,
			borrower_id,
			principal,
			interest_rate,
			weeks,
			status,
			created_by,
			updated_by,
			created_at,
			updated_at
		)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

func (r defaultLoansRepository) Create(ctx context.Context, loan models.Loan) (id uuid.UUID, err error) {
	if loan.Id == uuid.Nil {
		loan.Id = uuid.New()
	}
	args := []any{
		loan.Id,
		loan.BorrowerId,
		loan.Principal,
		loan.InterestRate,
		loan.Weeks,
		loan.Status,
		loan.CreatedBy,
		loan.UpdatedBy,
		loan.CreatedAt,
		loan.UpdatedAt,
	}

	tx := utils.SqlxTxFromCtx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, insertLoanQuery, args...)
	} else {
		_, err = r.db.ExecContext(ctx, insertLoanQuery, args...)
	}

	return loan.Id, err
}

const getLoanQuery = `
	SELECT
		id,
		borrower_id,
		principal,
		interest_rate,
		weeks,
		status,
		created_by,
		updated_by,
		created_at,
		updated_at
	FROM
		loans
	WHERE
		%s
`

func (r defaultLoansRepository) GetById(ctx context.Context, loanId uuid.UUID) (*models.Loan, error) {
	var err error
	query := fmt.Sprintf(getLoanQuery, "id = $1")
	loan := models.Loan{}
	tx := utils.SqlxTxFromCtx(ctx)
	args := []any{
		loanId,
	}
	if tx != nil {
		err = tx.GetContext(ctx, &loan, query, args...)
	} else {
		err = r.db.GetContext(ctx, &loan, query, args...)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custerr.ErrDataNotFound
		}
		return nil, err
	}
	return &loan, nil
}
