package services

type CreateLoanOutputInstallment struct {
	Week    int
	Amount  float64
	Status  int
	DueDate string
	LoanId  string
	IsFinal string
}

type CreateLoanOutput struct {
	Id string
}
