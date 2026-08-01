// models that is widely used accross the codebase got registered here
package models

// Actor represents an actor. BirthDate uses ISO 8601 (YYYY-MM-DD).
type Actor struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}

// ActorPatch describes a partial update; nil fields are left unchanged.
type ActorPatch struct {
	Id        int64    `json:"id"`
	Name      *string  `json:"name"`
	BirthDate *string  `json:"birth_date" validate:"omitempty,datetime=2006-01-02,pastdate"`
	MovieIDs  *[]int64 `json:"movie_ids" validate:"omitempty,dive,gt=0"`
}
