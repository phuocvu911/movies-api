package repository

import (
	"database/sql"
	"errors"
	"movies-api/internal/customerrors"
	"movies-api/internal/models"
	"strings"
)

// GenreRepository provides raw-SQL data access for genres.
type GenreRepository struct {
	db *sql.DB
}

func NewGenreRepository(db *sql.DB) *GenreRepository {
	return &GenreRepository{db: db}
}

// GetByName return genre with matching input (case-intensitive)
func (r *GenreRepository) GetByName(name string) (models.Genre, error) {
	row := r.db.QueryRow("SELECT id, name FROM genres WHERE LOWER(name) = LOWER(?)", name)
	var genre models.Genre
	if err := row.Scan(&genre.ID, &genre.Name); err != nil {
		return models.Genre{}, err
	}
	return genre, nil
}

// Create create a new genre.
func (r *GenreRepository) Create(name string) (models.Genre, error) {
	res, err := r.db.Exec("INSERT INTO genres (name) VALUES (?)", name)
	if err != nil {
		return models.Genre{}, err
	}

	GenreID, err := res.LastInsertId()
	if err != nil {
		return models.Genre{}, err
	}

	return models.Genre{
		ID:   GenreID,
		Name: name,
	}, nil
}

// GetAll returns all genres and returns a NotFoundError if no genres exist.
func (r *GenreRepository) GetAll() ([]models.Genre, error) {
	rows, err := r.db.Query("SELECT id, name FROM genres")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If no genres found, return a NotFoundError
	if len(genres) == 0 {
		return nil, customerrors.NotFoundf("No genre found")
	}
	return genres, nil
}

// GetByID returns a genre by its ID and returns a NotFoundError if the genre does not exist.
func (r *GenreRepository) GetByID(id int64) (models.Genre, error) {
	var genre models.Genre
	err := r.db.QueryRow("SELECT id, name FROM genres WHERE id = ?", id).Scan(&genre.ID, &genre.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Genre{}, customerrors.NotFoundf("Genre with ID %d not found", id)
	}
	return genre, nil
}

// GetMoviesByGenreID returns all movies associated with a specific genre ID.
func (r *GenreRepository) GetMoviesByGenreID(genreID int64) ([]models.Movie, error) {
	rows, err := r.db.Query(`
		SELECT m.id, m.title, m.release_year, m.duration
		FROM movies m
		JOIN movie_genre mg ON m.id = mg.movie_id
		WHERE mg.genre_id = ?`, genreID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Duration); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If no movies found for the genre, return a NotFoundError
	if len(movies) == 0 {
		return nil, customerrors.NotFoundf("No movies found for genre ID %d", genreID)
	}
	return movies, nil
}

// Update updates a genre's name by its ID.
func (r *GenreRepository) Update(id int64, input models.GenreRequest) (models.Genre, error) {
	result, err := r.db.Exec("UPDATE genres SET name = ? WHERE id = ?", input.Name, id)
	if err != nil {
		// If the error is due to a unique constraint violation, return a ConflictError
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return models.Genre{}, customerrors.Conflictf("Genre with the same name already exists")
		}
		return models.Genre{}, err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return models.Genre{}, customerrors.NotFoundf("Genre with ID %d not found", id)
	}

	return models.Genre{
		ID:   id,
		Name: input.Name,
	}, nil
}

// MovieCount returns the number of movies associated with a specific genre ID.
func (r *GenreRepository) MovieCount(genreID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM movie_genre WHERE genre_id = ?`, genreID).Scan(&count)
	return count, err
}

// Delete removes a genre by its ID.
func (r *GenreRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM genres WHERE id = ?`, id)
	return err
}
