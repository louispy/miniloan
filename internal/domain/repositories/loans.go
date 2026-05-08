package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/louispy/miniloan/internal/domain/models"
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
			principal,
			interest_rate,
			weeks,
			status,
			created_at,
			updated_at
		)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)
`

func (r defaultLoansRepository) Create(ctx context.Context, loan models.Loan) (uuid.UUID, error) {
	if loan.Id == uuid.Nil {
		loan.Id = uuid.New()
	}
	args := []any{
		loan.Id,
		loan.Principal,
		loan.InterestRate,
		loan.Weeks,
		loan.Status,
		loan.CreatedAt,
		loan.UpdatedAt,
	}
	_, err := r.db.ExecContext(ctx, insertLoanQuery, args...)

	return loan.Id, err
}
