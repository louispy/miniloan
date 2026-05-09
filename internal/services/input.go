package services

import "github.com/google/uuid"

type CreateLoanInput struct {
	Principal  int64
	BorrowerId uuid.UUID
}
