package services

import "github.com/google/uuid"

type CreateLoanInput struct {
	BorrowerId uuid.UUID
}
