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

func (r defaultLoansRepository) Create(ctx context.Context, loan models.Loan) (uuid.UUID, error) {
	return uuid.Nil, nil
}
