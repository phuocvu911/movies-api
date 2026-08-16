package handlers

import (
	"movies-api/internal/customerrors"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"movies-api/internal/validation"
	"net/http"
)

// GenreHandler exposes genres endpoint.
type GenreHandler struct {
	service *service.GenreService
}

// NewGenreHandler creates a new GenreHandler.
func NewGenreHandler(s *service.GenreService) *GenreHandler {
	return &GenreHandler{service: s}
}

// Create handle POST /api/genres.
func (h *GenreHandler) Create(w http.ResponseWriter, r *http.Request) {
	var genreRequest models.GenreRequest

	if err := decodeJSON(r, &genreRequest); err != nil {
		respondError(w, err)
		return
	}

	if err := validation.V.Struct(genreRequest); err != nil {
		respondError(w, customerrors.Validationf("validation err: %v", err))
		return
	}

	genre, err := h.service.Create(genreRequest)
	if err != nil {
		respondError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, genre)
}

// GetAll handles GET /api/genres.
func (h *GenreHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, size, err := pagination(r)
	if err != nil {
		respondError(w, err)
		return
	}

	genres, err := h.service.GetAll(page, size)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, genres)
}

// GetByID handles GET /api/genres/{id}
func (h *GenreHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}
	genre, err := h.service.GetByID(id)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, genre)
}

// Movies handles GET /api/genres/{id}/movies
func (h *GenreHandler) Movies(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}

	page, size, err := pagination(r)
	if err != nil {
		respondError(w, err)
		return
	}

	movies, err := h.service.Movies(id, page, size)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, movies)
}

// Update handles PATCH /api/genres/{id}
func (h *GenreHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}

	var genreRequest models.GenreRequest

	if err := decodeJSON(r, &genreRequest); err != nil {
		respondError(w, err)
		return
	}

	genre, err := h.service.Update(id, genreRequest)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, genre)
}

// Delete handles DELETE /api/genres/{id}
func (h *GenreHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}

	//get the force query parameter
	force, err := forceParam(r)
	if err != nil {
		respondError(w, err)
		return
	}
	err = h.service.Delete(id, force)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
