package repository

import (
	"database/sql"
	"errors"
	"movies-api/internal/customerrors"
	"movies-api/internal/models"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) GetAll() ([]models.Movie, error) {
	rows, err := r.db.Query("SELECT id, title, release_year, duration FROM movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := []models.Movie{}
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Duration); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	// If no movies found, return a NotFoundError
	if len(movies) == 0 {
		return nil, customerrors.NotFoundf("No movie found")
	}
	return movies, nil
}

// Search searches for movies by title (partial match, case-insensitive).
func (r *MovieRepository) Search(title string) ([]models.Movie, error) {
	rows, err := r.db.Query("SELECT id, title, release_year, duration FROM movies WHERE LOWER(title) LIKE ?", "%"+title+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := []models.Movie{}
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Duration); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	// If no movies found, return a NotFoundError
	if len(movies) == 0 {
		return nil, customerrors.NotFoundf("No movie found for title containing '%s'", title)
	}
	return movies, nil
}

func (r *MovieRepository) Create(title string, releaseYear, duration int) (models.Movie, error) {
	res, err := r.db.Exec("INSERT INTO movies (title, release_year, duration) VALUES (?, ?, ?)", title, releaseYear, duration)
	if err != nil {
		return models.Movie{}, err
	}

	movieID, err := res.LastInsertId()
	if err != nil {
		return models.Movie{}, err
	}

	return models.Movie{
		ID:       movieID,
		Title:    title,
		Year:     releaseYear,
		Duration: duration,
	}, nil
}

func (r *MovieRepository) GetByID(id int64) (models.Movie, error) {
	var movie models.Movie
	err := r.db.QueryRow("SELECT id, title, release_year, duration FROM movies WHERE id = ?", id).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Duration)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Movie{}, customerrors.NotFoundf("Movie with id %d not found", id)
	}
	return movie, err
}
