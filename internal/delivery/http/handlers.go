package httpdelivery

import (
	"encoding/json"
	"github.com/buguzei/medods-task/internal/models"
	"github.com/buguzei/medods-task/internal/usecase"
	"github.com/buguzei/medods-task/pkg/config"
	"net/http"
)

type Handler struct {
	uc  usecase.AuthUC
	cfg config.Config
}

func NewHandler(cfg *config.Config, uc usecase.AuthUC) Handler {
	return Handler{
		uc:  uc,
		cfg: *cfg,
	}
}

func (h Handler) NewPair(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenPair, err := h.uc.NewPair(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bResp, err := json.Marshal(tokenPair)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Authorization")

	resp, err := h.uc.Refresh(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
