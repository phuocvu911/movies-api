package repository

import (
	"database/sql"
	"movies-api/internal/customerrors"
	"movies-api/internal/models"
	"strings"
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

func (r *MovieRepository) Create(input models.MovieRequest) (models.Movie, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.Movie{}, err
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO movies (title, release_year, duration) VALUES (?, ?, ?)", input.Title, input.Year, input.Duration)
	if err != nil {
		return models.Movie{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Movie{}, err
	}

	//add relationships
	for _, genreID := range dedupe(input.GenreIDs) {
		_, err := tx.Exec("INSERT INTO movie_genre (movie_id, genre_id) VALUES (?, ?)", id, genreID)
		if err != nil {
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				return models.Movie{}, customerrors.NotFoundf("Genre with ID %d not found", genreID)
			}
			return models.Movie{}, err
		}
	}

	for _, actorID := range dedupe(input.ActorIDs) {
		_, err := tx.Exec("INSERT INTO movie_actor (movie_id, actor_id) VALUES (?, ?)", id, actorID)
		if err != nil {
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				return models.Movie{}, customerrors.NotFoundf("Actor with ID %d not found", actorID)
			}
			return models.Movie{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return models.Movie{}, err
	}

	return models.Movie{
		ID:       id,
		Title:    input.Title,
		Year:     input.Year,
		Duration: input.Duration,
	}, nil
}
