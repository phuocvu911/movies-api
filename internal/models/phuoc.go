// models that is widely used accross the codebase got registered here
package models

// ActorRequest is used to decode the request body for creating an actor(POST).
type ActorRequest struct {
	Name      string `json:"name" validate:"required"`
	BirthDate string `json:"birth_date" validate:"required,datetime=2006-01-02,pastdate"`
}

// Actor represents an actor. BirthDate uses ISO 8601 (YYYY-MM-DD).
type Actor struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}

// ActorPatch describes a partial update(PATCH); nil fields are left unchanged.
type ActorPatch struct {
	Id        int64    `json:"id"`
	Name      *string  `json:"name"`
	BirthDate *string  `json:"birth_date" validate:"omitempty,datetime=2006-01-02,pastdate"`
	MovieIDs  *[]int64 `json:"movie_ids" validate:"omitempty,dive,gt=0"`
}

// GenreRequest is used to decode the request body for creating a genre.
type GenreRequest struct {
	Name string `json:"name" validate:"required"`
}

// Genre represents a genre.
type Genre struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
