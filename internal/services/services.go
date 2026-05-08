package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/louispy/miniloan/internal/domain/models"
)

type LoanService interface {
	Create(ctx context.Context, loan models.Loan) (uuid.UUID, error)
}
