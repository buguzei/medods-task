package server

import (
	delivery "github.com/buguzei/medods-task/internal/delivery/http"
	"github.com/gorilla/mux"
)

func InitRoutes(r *mux.Router, h delivery.Handler) *mux.Router {
	r.HandleFunc("/new-pair", h.NewPair).Methods("POST")
	r.HandleFunc("/refresh", h.Refresh).Methods("POST")

	return r
}
