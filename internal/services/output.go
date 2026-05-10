package services

type CreateLoanOutputInstallment struct {
	Id      string
	Week    int
	Amount  int64
	Status  int
	DueDate string
	LoanId  string
	IsFinal bool
}

type CreateLoanOutput struct {
	Id           string
	Installments []CreateLoanOutputInstallment
}

type GetOutstandingOutput struct {
	Amount int64
}

type GetInstallmentsOutputInstallment struct {
	Id      string
	Week    int
	Amount  int64
	Status  int
	DueDate string
	LoanId  string
}

type GetInstallmentsOutput struct {
	Installments []GetInstallmentsOutputInstallment
}

type IsDeliquentOutput struct {
	IsDeliquent bool
}

type MakePaymentOutput struct {
	InstallmentId string
}
