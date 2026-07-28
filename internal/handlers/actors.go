package handlers

import (
	"movies-api/internal/service"
	"movies-api/internal/validation"
	"net/http"
)

// actorRequest is the body for both POST and PATCH; pointers distinguish
// absent fields from zero values so PATCH can be partial.
type actorRequest struct {
	Name      *string  `json:"name" validate:"required"`
	BirthDate *string  `json:"birth_date" validate:"required,datetime=2006-01-02,pastdate"`
	MovieIDs  *[]int64 `json:"movie_ids"` // optional, can be empty or null
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.V.Struct(actorRequest); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity) // 422
		return
	}

	input:= service.ActorInput{
		Name:      *actorRequest.Name,
		BirthDate: *actorRequest.BirthDate,
	}

	if actorRequest.MovieIDs != nil {
		input.MovieIDs = *actorRequest.MovieIDs
	}

	actor, err:= h.service.Create(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, actor)
}
