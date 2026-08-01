package handlers

import (
	"movies-api/internal/customerrors"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"movies-api/internal/validation"
	"net/http"
)

// actorRequest is used to decode the request body for creating an actor.
type actorRequest struct {
	Name      *string `json:"name" validate:"required"`
	BirthDate *string `json:"birth_date" validate:"required,datetime=2006-01-02,pastdate"`
}

// ActorHandler exposes the actor endpoints.
type ActorHandler struct {
	service *service.ActorService
}

func NewActorHandler(s *service.ActorService) *ActorHandler {
	return &ActorHandler{service: s}
}

// Create handles POST /api/actors.
func (h *ActorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var actorRequest actorRequest

	if err := decodeJSON(r, &actorRequest); err != nil {
		respondError(w, err)
		return
	}

	if err := validation.V.Struct(actorRequest); err != nil {
		respondError(w, customerrors.Validationf("validation err: %v", err))
		return
	}

	input := service.ActorInput{
		Name:      *actorRequest.Name,
		BirthDate: *actorRequest.BirthDate,
	}

	actor, err := h.service.Create(input)
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

	var actors []models.Actor
	var err error

	if name != "" {
		// Filter by name (partially, case-insensitive)
		actors, err = h.service.GetByName(name)
	} else {
		// Return all actors
		actors, err = h.service.GetAll()
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
// Movies handles GET /api/actors/{id}/movies.
func (h *ActorHandler) Movies(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		respondError(w, err)
		return
	}

	result, err := h.service.Movies(id)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
