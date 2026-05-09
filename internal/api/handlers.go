package api

import (
	"encoding/json"
	"net/http"

	"github.com/louispy/miniloan/internal/services"
)

func (a API) CreateLoanHandler(rw http.ResponseWriter, r *http.Request) {
	req := CreateLoanRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}

	if err := req.Validate(); err != nil {
		WriteJSONResponse(rw, 400, nil, nil, err)
		return
	}

	input := services.CreateLoanInput{
		Principal: req.Principal,
	}

	output, err := a.loanService.Create(r.Context(), input)
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
