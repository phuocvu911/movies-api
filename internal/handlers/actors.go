package handlers

import (
	"movies-api/internal/customerrors"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"movies-api/internal/validation"
	"net/http"
)

// ActorHandler exposes the actor endpoints.
type ActorHandler struct {
	service *service.ActorService
}

// NewActorHandler creates a new ActorHandler.
func NewActorHandler(s *service.ActorService) *ActorHandler {
	return &ActorHandler{service: s}
}

// Create handles POST /api/actors.
func (h *ActorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var actorRequest models.ActorRequest

	if err := decodeJSON(r, &actorRequest); err != nil {
		respondError(w, err)
		return
	}

	if err := validation.V.Struct(actorRequest); err != nil {
		respondError(w, customerrors.Validationf("validation err: %v", err))
		return
	}

	actor, err := h.service.Create(actorRequest)
	if err != nil {
		respondError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, actor)
}

// GetAll handles GET /api/actors.
func (h *ActorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Check for optional name query parameter
	name := r.URL.Query().Get("name")

	//take the pagination parameters from the query string
	page, size, err := pagination(r)
	if err != nil {
		respondError(w, err)
		return
	}

	var actors models.Page[models.Actor]

	if name != "" {
		// Filter by name (partially, case-insensitive)
		actors, err = h.service.GetByName(name, page, size)
	} else {
		// Return all actors
		actors, err = h.service.GetAll(page, size)
	}

	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, actors)
}

// GetByID handles GET /api/actors/{id}.
func (h *ActorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}
	actor, err := h.service.GetByID(id)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, actor)
}

// Update handles PATCH /api/actors/{id}.
func (h *ActorHandler) Update(w http.ResponseWriter, r *http.Request) {
	//get actor id from the path
	id, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}
	var actorPatchRequest models.ActorPatch
	if err := decodeJSON(r, &actorPatchRequest); err != nil {
		respondError(w, err)
		return
	}

	//validate the request body
	if err := validation.V.Struct(actorPatchRequest); err != nil {
		respondError(w, customerrors.Validationf("validation err: %v", err))
		return
	}
	//call the service to update the actor
	actor, err := h.service.Update(id, actorPatchRequest)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, actor)
}

// Delete handles DELETE /api/actors/{id}.
func (h *ActorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	//get actor id from the path
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

	if err := h.service.Delete(id, force); err != nil {
		respondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Movies handles GET /api/actors/{id}/movies.
func (h *ActorHandler) Movies(w http.ResponseWriter, r *http.Request) {
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
