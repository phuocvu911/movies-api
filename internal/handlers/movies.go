package handlers

import (
	"movies-api/internal/service"
	"net/http"
)

type MovieHandler struct {
	service *service.MovieService
}

func NewMovieHandler(s *service.MovieService) *MovieHandler {
	return &MovieHandler{service: s}
}

func (h *MovieHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	movies, err := h.service.GetAll()
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, movies)
}
