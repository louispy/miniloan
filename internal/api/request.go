package api

import (
	"errors"
)

type MakePaymentRequest struct {
	Amount        int64  `json:"amount"`
	InstallmentId string `json:"installment_id"`
}

func (r MakePaymentRequest) Validate() error {
	if r.Amount <= 0 {
		return errors.New("amount should be > 0")
	}

	return nil
}
