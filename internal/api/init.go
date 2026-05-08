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
func (r *API) GetRouter() *mux.Router {
	return r.router
}

func (r *API) Register() {
	r.router.HandleFunc("/hello", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "Hello World")
	})
}
