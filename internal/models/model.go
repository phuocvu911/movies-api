package models

type Movie struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"release_year"`
	Duration int    `json:"duration"`
}

type MovieRequest struct {
	Title    string `json:"title" validate:"required"`
	Year     int    `json:"release_year"`
	Duration int    `json:"duration"`
}
