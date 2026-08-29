package models

// Movie represents a movie
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
	Year     *int     `json:"release_year" validate:"omitempty,gt=0,max=2027"`
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
	TotalElements int `json:"total_elements"`
	TotalPages    int `json:"total_pages"`
}

// Page is a paginated list response envelope, generic over the type of results it contains
type Page[T any] struct {
	Results    []T      `json:"results"`
	Pagination PageInfo `json:"pagination"`
}

// MovieFilter represents optional filters for querying movies
type MovieFilter struct {
	Year    *int
	GenreID *int64
	ActorID *int64
}

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

