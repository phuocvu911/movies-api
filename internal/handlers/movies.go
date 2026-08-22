package handlers

import (
	"movies-api/internal/customerrors"
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

func (h *MovieHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("title")

	if query == "" {
		respondError(w, customerrors.Validationf("Title query parameter is missing or empty"))
		return
	}

	movies, err := h.service.Search(query)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, movies)
}

func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *MovieHandler) GetByID(w http.ResponseWriter, r *http.Request) {

}
