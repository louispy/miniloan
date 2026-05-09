package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/louispy/miniloan/internal/constants"
	"github.com/louispy/miniloan/internal/services"
)

func (a API) CreateLoanHandler(rw http.ResponseWriter, r *http.Request) {
	output, err := a.loanService.Create(r.Context(), services.CreateLoanInput{})
	if err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}
	installments := make([]CreateLoanResponseInstallment, 0, len(output.Installments))
	for _, ins := range output.Installments {
		installments = append(installments, CreateLoanResponseInstallment{
			Week:    ins.Week,
			Amount:  ins.Amount,
			Status:  ins.Status,
			DueDate: ins.DueDate,
		})
	}
	resp := CreateLoanResponse{
		Id:           output.Id,
		Installments: installments,
	}
	message := "Successfully created a new loan request"

	WriteJSONResponse(rw, 200, resp, &message, nil)
}

func (a API) GetOutstanding(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuId, err := uuid.Parse(id)
	if err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}
	output, err := a.loanService.GetOutstanding(r.Context(), services.GetOutstandingInput{
		LoanId: uuId,
	})
	if err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}
	resp := GetOutstandingResponse{
		Amount: output.Amount,
	}
	message := "Successfully get outstanding amount"

	WriteJSONResponse(rw, 200, resp, &message, nil)
}

func (a API) GetInstallments(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuId, err := uuid.Parse(id)
	if err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}

	status := 0
	switch r.URL.Query().Get("status") {
	case "paid":
		status = constants.INSTALLMENT_STATUS_PAID
	case "unpaid":
		status = constants.INSTALLMENT_STATUS_UNPAID
	case "":
	default:
		WriteJSONResponse(rw, 400, nil, nil, errors.New("invalid status filter"))
		return
	}

	output, err := a.loanService.GetInstallments(r.Context(), services.GetInstallmentsInput{
		LoanId: uuId,
		Status: status,
	})
	if err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}

	installments := make([]GetInstallmentsResponseInstallment, 0, len(output.Installments))
	for _, ins := range output.Installments {
		installments = append(installments, GetInstallmentsResponseInstallment{
			Week:    ins.Week,
			Amount:  ins.Amount,
			Status:  ins.Status,
			DueDate: ins.DueDate,
		})
	}
	resp := GetInstallmentsResponse{
		Installments: installments,
	}
	message := "Successfully get installments"

	WriteJSONResponse(rw, 200, resp, &message, nil)
}

func (a API) IsDeliquent(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuId, err := uuid.Parse(id)
	if err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}
	output, err := a.loanService.IsDeliquent(r.Context(), services.IsDeliquentInput{
		LoanId:    uuId,
		Timestamp: time.Now(),
	})
	if err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}
	resp := IsDeliquentResponse{
		IsDeliquent: output.IsDeliquent,
	}
	message := "Successfully get IsDeliquent value"

	WriteJSONResponse(rw, 200, resp, &message, nil)
}

func (a API) MakePayment(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuId, err := uuid.Parse(id)
	if err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}

	r.Body = http.MaxBytesReader(rw, r.Body, 1<<20)
	req := MakePaymentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}

	if err := req.Validate(); err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}

	installmentId, err := uuid.Parse(req.InstallmentId)
	if err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}

	output, err := a.loanService.MakePayment(r.Context(), services.MakePaymentInput{
		LoanId:        uuId,
		Amount:        req.Amount,
		InstallmentId: installmentId,
	})
	if err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}
	resp := MakePaymentResponse{
		InstallmentId: output.InstallmentId,
	}
	message := "Successfully MakePayment"

	WriteJSONResponse(rw, 200, resp, &message, nil)
}
