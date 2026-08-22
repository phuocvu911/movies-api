package handlers

import (
	"movies-api/internal/customerrors"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"movies-api/internal/validation"
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
	var movieRequest models.MovieRequest

	if err := decodeJSON(r, &movieRequest); err != nil {
		respondError(w, err)
		return
	}

	if err := validation.V.Struct(movieRequest); err != nil {
		respondError(w, customerrors.Validationf("validation err: %v", err))
		return
	}

	movie, err := h.service.Create(movieRequest)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, movie)
}
