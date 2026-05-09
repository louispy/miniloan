package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/louispy/miniloan/internal/services"
)

type API struct {
	loanService services.LoanService
	router      *mux.Router
}

type Opts struct {
	LoanService services.LoanService
}

func NewAPI(o Opts) *API {
	r := mux.Router{}
	return &API{
		loanService: o.LoanService,
		router:      &r,
	}
}
func (a *API) GetRouter() *mux.Router {
	return a.router
}

func (a *API) Register() {
	a.router.HandleFunc("/hello", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "Hello World")
	})
	a.router.Methods(http.MethodPost).Path("/loans").HandlerFunc(a.CreateLoanHandler)
	a.router.Methods(http.MethodGet).Path("/loans/{id}/outstanding").HandlerFunc(a.GetOutstanding)
	a.router.Methods(http.MethodGet).Path("/loans/{id}/deliquent").HandlerFunc(a.IsDeliquent)
	a.router.Methods(http.MethodPost).Path("/loans/{id}/payment").HandlerFunc(a.MakePayment)
}
