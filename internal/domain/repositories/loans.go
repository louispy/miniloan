package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/louispy/miniloan/internal/domain/models"
)

type defaultLoansRepository struct {
}

func NewLoansRepository() LoansRepository {
	return &defaultLoansRepository{}
}

func (r defaultLoansRepository) Create(ctx context.Context, loan models.Loan) (uuid.UUID, error) {
	return uuid.Nil, nil
}
