package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type jsonResponse struct {
	Data    any     `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	Error   string  `json:"error,omitempty"`
}

func WriteJSONResponse(rw http.ResponseWriter, statusCode int, body any, message *string, err error) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	resp := jsonResponse{}
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Data = body
		resp.Message = message
	}
	err = json.NewEncoder(rw).Encode(resp)
	if err != nil {
		log.Printf("Error encoding response as json: %v", err.Error())
	}
}

type CreateLoanResponseInstallment struct {
	Id      string `json:"id"`
	Week    int    `json:"week"`
	Amount  int64  `json:"amount"`
	Status  int    `json:"status"`
	DueDate string `json:"due_date"`
}

type CreateLoanResponse struct {
	Id           string                          `json:"id"`
	Installments []CreateLoanResponseInstallment `json:"installments"`
}

type GetOutstandingResponse struct {
	Amount int64 `json:"amount"`
}

type GetInstallmentsResponseInstallment struct {
	Id      string `json:"id"`
	Week    int    `json:"week"`
	Amount  int64  `json:"amount"`
	Status  int    `json:"status"`
	DueDate string `json:"due_date"`
}

type GetInstallmentsResponse struct {
	Installments []GetInstallmentsResponseInstallment `json:"installments"`
}

type IsDeliquentResponse struct {
	IsDeliquent bool `json:"is_deliquent"`
}

type MakePaymentResponse struct {
	InstallmentId string `json:"installment_id"`
}
