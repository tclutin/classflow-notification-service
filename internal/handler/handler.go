package handler

import (
	"net/http"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Init() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", h.CheckHealth)

	return router
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
