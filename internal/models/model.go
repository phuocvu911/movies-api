package models

type Movie struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"release_year"`
	Duration int    `json:"duration"`
}

// MovieRequest is used to decode the request body for creating a movie (POST)
type MovieRequest struct {
	Title    string  `json:"title" validate:"required"`
	Year     int     `json:"release_year" validate:"required,gt=0"`
	Duration int     `json:"duration" validate:"required,gt=0"`
	GenreIDs []int64 `json:"genre_ids" validate:"required,dive,gt=0"`
	ActorIDs []int64 `json:"actor_ids" validate:"required,dive,gt=0"`
}

// MoviePatch describes a partial update (PATCH); nil fields are left unchanged
type MoviePatch struct {
	ID       int64
	Title    *string  `json:"title"`
	Year     *int     `json:"release_year" validate:"omitempty,gt=0"`
	Duration *int     `json:"duration" validate:"omitempty,gt=0"`
	GenreIDs *[]int64 `json:"genre_id" validate:"omitempty,dive,gt=0"`
	ActorIDs *[]int64 `json:"actor_id" validate:"omitempty,dive,gt=0"`
}

// has all relationships
type MovieDetail struct {
	ID       int64   `json:"id"`
	Title    string  `json:"title"`
	Year     int     `json:"release_year"`
	Duration int     `json:"duration"`
	GenreIDs []int64 `json:"genre_ids"`
	ActorIDs []int64 `json:"actor_ids"`
}
