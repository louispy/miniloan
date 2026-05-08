package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type APIRouter struct {
	router *mux.Router
}

func NewAPIRouter() *APIRouter {
	r := mux.Router{}
	return &APIRouter{router: &r}
}
func (r *APIRouter) GetRouter() *mux.Router {
	return r.router
}

func (r *APIRouter) Register() {
	r.router.HandleFunc("/hello", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "Hello World")
	})
}
