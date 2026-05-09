package api

import (
	"net/http"

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
