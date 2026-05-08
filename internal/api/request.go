package api

import "errors"

type CreateLoanRequest struct {
	Principal    int64   `json:"principal"`
	InterestRate float64 `json:"interestRate"`
}

func (r CreateLoanRequest) Validate() error {
	if r.Principal <= 0 {
		return errors.New("principal should be > 0")
	}
	if r.InterestRate <= 0 {
		return errors.New("interestRate should be > 0")
	}
	return nil
}
