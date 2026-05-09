package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/louispy/miniloan/internal/services"
)

func (a API) CreateLoanHandler(rw http.ResponseWriter, r *http.Request) {
	output, err := a.loanService.Create(r.Context(), services.CreateLoanInput{})
	if err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}
	resp := CreateLoanResponse{
		Id: output.Id,
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
