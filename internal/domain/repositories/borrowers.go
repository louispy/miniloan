package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/louispy/miniloan/internal/domain/models"
)

type defaultBorrowersRepository struct {
	db *sqlx.DB
}

type BorrowerRepoOpts struct {
	DB *sqlx.DB
}

func NewBorrowersRepository(opts LoanRepoOpts) BorrowersRepository {
	return &defaultBorrowersRepository{
		db: opts.DB,
	}
}

const insertBorrowerQuery = `
	INSERT INTO
		borrowers (
			id,
			full_name,
			created_by,
			updated_by,
			created_at,
			updated_at
		)
	VALUES
		($1, $2, $3, $4, $5, $6)
`

func (r defaultBorrowersRepository) Create(ctx context.Context, borrower models.Borrower) (uuid.UUID, error) {
	if borrower.Id == uuid.Nil {
		borrower.Id = uuid.New()
	}
	args := []any{
		borrower.Id,
		borrower.FullName,
		borrower.CreatedBy,
		borrower.UpdatedBy,
		borrower.CreatedAt,
		borrower.UpdatedAt,
	}
	_, err := r.db.ExecContext(ctx, insertBorrowerQuery, args...)

	return borrower.Id, err
}
