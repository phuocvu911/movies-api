package repository

import (
	"database/sql"
	"errors"
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

func (r *MovieRepository) Update(id int64, u models.MoviePatch) error {
	var sets []string
	var args []any

	if u.Title != nil {
		sets = append(sets, "title = ?")
		args = append(args, *u.Title)
	}
	if u.Year != nil {
		sets = append(sets, "release_year = ?")
		args = append(args, *u.Year)
	}
	if u.Duration != nil {
		sets = append(sets, "duration = ?")
		args = append(args, *u.Duration)
	}

	if len(sets) == 0 {
		return customerrors.Validationf("No fields to update")
	}

	args = append(args, id)
	result, err := r.db.Exec("UPDATE movies SET "+strings.Join(sets, ", ")+" WHERE id = ?", args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return customerrors.NotFoundf("Movie with id %d not found", id)
	}

	return nil
}

func (r *MovieRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM movies WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return customerrors.NotFoundf("Movie with id %d not found", id)
	}

	return nil
}

func (r *MovieRepository) GetByGenreID(genreID int64) ([]models.Movie, error) {
	rows, err := r.db.Query("SELECT m.id, m.title, m.release_year, m.duration FROM movies m JOIN movie_genre mg ON m.id = mg.movie_id WHERE mg.genre_id = ?", genreID)
	if err != nil {
		return []models.Movie{}, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Duration); err != nil {
			return movies, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	// If no movies found, return a NotFoundError
	if len(movies) == 0 {
		return nil, customerrors.NotFoundf("No movie found for genre ID %d", genreID)
	}
	return movies, nil
}

func (r *MovieRepository) GetByYear(releaseYear int) ([]models.Movie, error) {
	rows, err := r.db.Query("SELECT id, title, release_year, duration FROM movies WHERE release_year = ?", releaseYear)
	if err != nil {
		return []models.Movie{}, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Duration); err != nil {
			return movies, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	// If no movies found, return a NotFoundError
	if len(movies) == 0 {
		return nil, customerrors.NotFoundf("No movie found for release year %d", releaseYear)
	}
	return movies, nil
}

func (r *MovieRepository) GetByActorID(actorID int64) ([]models.Movie, error) {
	rows, err := r.db.Query("SELECT m.id, m.title, m.release_year, m.duration FROM movies m JOIN movie_actor ma ON ma.movie_id = m.id WHERE ma.actor_id = ?", actorID)
	if err != nil {
		return []models.Movie{}, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Duration); err != nil {
			return movies, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	// If no movies found, return a NotFoundError
	if len(movies) == 0 {
		return nil, customerrors.NotFoundf("No movie found for actor ID %d", actorID)
	}
	return movies, nil
}

func (r *MovieRepository) GetActorsByMovieID(id int) []models.Actor
