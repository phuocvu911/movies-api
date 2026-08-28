package handlers

import (
	"movies-api/internal/customerrors"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"movies-api/internal/validation"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	service *service.MovieService
}

func NewMovieHandler(s *service.MovieService) *MovieHandler {
	return &MovieHandler{service: s}
}

func (h *MovieHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, size, err := pagination(r)
	if err != nil {
		respondError(w, err)
		return
	}

	filter := models.MovieFilter{}

	genre := r.URL.Query().Get("genre")
	if genre != "" {
		genreID, err := strconv.ParseInt(genre, 10, 64)
		if err != nil || genreID <= 0 {
			respondError(w, customerrors.Validationf("Invalid genre ID %v", genreID))
			return
		}
		movies, err := h.service.GetByGenreID(genreID)
		if err != nil {
			respondError(w, err)
			return
		}

		filter.GenreID = &genreID
	}

	year := r.URL.Query().Get("year")
	if year != "" {
		yearNo, err := strconv.Atoi(year)
		if err != nil || yearNo <= 0 {
			respondError(w, customerrors.Validationf("Invalid release year %v", yearNo))
			return
		}

		filter.Year = &yearNo
	}

	actor := r.URL.Query().Get("actor")
	if actor != "" {
		actorID, err := strconv.ParseInt(actor, 10, 64)
		if err != nil || actorID <= 0 {
			respondError(w, customerrors.Validationf("Invalid actor ID %v", actorID))
			return
		}

		filter.ActorID = &actorID
	}

	movies, err := h.service.GetAll(filter, page, size)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, movies)
}

func (h *MovieHandler) Search(w http.ResponseWriter, r *http.Request) {
	page, size, err := pagination(r)
	if err != nil {
		respondError(w, err)
		return
	}

	query := r.URL.Query().Get("title")

	if query == "" {
		respondError(w, customerrors.Validationf("Title query parameter is missing or empty"))
		return
	}

	movies, err := h.service.Search(query, page, size)
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
	movieID, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}

	movie, err := h.service.GetByID(movieID)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, movie)
}

func (h *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
	movieID, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}

	var moviePatchRequest models.MoviePatch
	if err := decodeJSON(r, &moviePatchRequest); err != nil {
		respondError(w, err)
		return
	}

	if err := validation.V.Struct(moviePatchRequest); err != nil {
		respondError(w, customerrors.Validationf("validation error: %v", err))
		return
	}

	movie, err := h.service.Update(movieID, moviePatchRequest)
	if err != nil {
		respondError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, movie)
}

func (h *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	movieID, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}

	force, err := forceParam(r)
	if err != nil {
		respondError(w, err)
		return
	}

	if err := h.service.Delete(movieID, force); err != nil {
		respondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MovieHandler) Actors(w http.ResponseWriter, r *http.Request) {
	page, size, err := pagination(r)
	if err != nil {
		respondError(w, err)
		return
	}
	movieID, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}

	actors, err := h.service.GetActorsByMovieID(movieID, page, size)
	if err != nil {
		respondError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, actors)
}
