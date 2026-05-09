package constants

// hardcoded per task spec
// should there be more loan configurations,
// a new table (i.e. loan_plans) can be introduced to store loan configs.
const LOAN_WEEKS = 50
const LOAN_INTEREST_RATE = 0.1
const LOAN_PRINCIPAL = 5000000

const (
	_ = iota
	LOAN_STATUS_REQUESTED
	LOAN_STATUS_APPROVED
	LOAN_STATUS_DISBURSED
	LOAN_STATUS_FULLY_PAID
)

const (
	_ = iota
	INSTALLMENT_STATUS_UNPAID
	INSTALLMENT_STATUS_PAID
)

const DELIQUENT_LATE_COUNT = 2
