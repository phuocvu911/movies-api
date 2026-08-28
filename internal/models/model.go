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
	Year     int     `json:"release_year" validate:"required,gt=0,max=2027"`
	Duration int     `json:"duration" validate:"required,gt=0"`
	GenreIDs []int64 `json:"genre_ids" validate:"omitempty,dive,gt=0"`
	ActorIDs []int64 `json:"actor_ids" validate:"omitempty,dive,gt=0"`
}

// MoviePatch describes a partial update (PATCH); nil fields are left unchanged
type MoviePatch struct {
	ID       int64    `json:"id"`
	Title    *string  `json:"title"`
	Year     *int     `json:"release_year" validate:"omitempty,gt=0"`
	Duration *int     `json:"duration" validate:"omitempty,gt=0"`
	GenreIDs *[]int64 `json:"genre_ids" validate:"omitempty,dive,gt=0"`
	ActorIDs *[]int64 `json:"actor_ids" validate:"omitempty,dive,gt=0"`
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

// PageInfo carries pagination metadata for list responses
type PageInfo struct {
	Page          int `json:"page"`
	Size          int `json:"size"`
	TotalElements int `json:"total_elements`
	TotalPages    int `json:"total_pages`
}

// Page is a paginated list response envelope, generic over the type of results it contains
type Page[T any] struct {
	Results    []T      `json:"results"`
	Pagination PageInfo `json:"pagination"`
}

type MovieFilter struct {
	GenreID *int64
	Year    *int
	ActorID *int64
}
