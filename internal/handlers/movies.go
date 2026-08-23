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

func (h *MovieHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}

	movie, err := h.service.GetByID(id)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, movie)
}

func (h *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
	//get movie id from the path
	id, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}
	var moviePatchRequest models.MoviePatchRequest
	if err := decodeJSON(r, &moviePatchRequest); err != nil {
		respondError(w, err)
		return
	}

	//validate the request body
	if err := validation.V.Struct(moviePatchRequest); err != nil {
		respondError(w, customerrors.Validationf("validation err: %v", err))
		return
	}
	//call the service to update the movie
	movie, err := h.service.Update(id, moviePatchRequest)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, movie)
}

func (h *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusNoContent)
}
