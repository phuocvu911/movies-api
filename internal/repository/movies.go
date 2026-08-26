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
	if err := rows.Err(); err != nil {
		return nil, err
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
	if err := rows.Err(); err != nil {
		return nil, err
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

	res, err := tx.Exec("INSERT INTO movies (title, release_year, duration) VALUES (?, ?, ?)", input.Title, input.Year, input.Duration)
	if err != nil {
		return models.Movie{}, err
	}

	movieID, err := res.LastInsertId()
	if err != nil {
		return models.Movie{}, err
	}

	for _, genreID := range dedupe(input.GenreIDs) {
		if _, err := tx.Exec("INSERT INTO movie_genre (movie_id, genre_id) VALUES (?,?)", movieID, genreID); err != nil {
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				return models.Movie{}, customerrors.NotFoundf("Genre with id %d not found", genreID)
			}
			return models.Movie{}, err
		}
	}

	for _, actorID := range dedupe(input.ActorIDs) {
		if _, err := tx.Exec("INSERT INTO movie_actor (movie_id, actor_id) VALUES (?,?)", movieID, actorID); err != nil {
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				return models.Movie{}, customerrors.NotFoundf("Actor with id %d not found", actorID)
			}
			return models.Movie{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return models.Movie{}, err
	}
	return models.Movie{
		ID:       movieID,
		Title:    input.Title,
		Year:     input.Year,
		Duration: input.Duration,
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

	if len(sets) == 0 && u.GenreIDs == nil && u.ActorIDs == nil {
		return customerrors.Validationf("No fields to update")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists int
	err = tx.QueryRow("SELECT 1 FROM movies WHERE id = ?", id).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return customerrors.NotFoundf("Movie with id %d not found", id)
	}
	if err != nil {
		return err
	}

	if len(sets) > 0 {
		args = append(args, id)
		_, err := tx.Exec("UPDATE movies SET "+strings.Join(sets, ", ")+" WHERE id = ?", args...)
		if err != nil {
			return err
		}
	}

	if u.GenreIDs != nil {
		_, err := tx.Exec("DELETE FROM movie_genre WHERE movie_id = ?", id)
		if err != nil {
			return err
		}

		for _, genreID := range dedupe(*u.GenreIDs) {
			if _, err := tx.Exec("INSERT INTO movie_genre (movie_id, genre_id) VALUES (?,?)", id, genreID); err != nil {
				if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
					return customerrors.NotFoundf("Genre with id %d not found", genreID)
				}
				return err
			}
		}
	}

	if u.ActorIDs != nil {
		if _, err := tx.Exec("DELETE FROM movie_actor WHERE movie_id = ?", id); err != nil {
			return err
		}

		for _, actorID := range dedupe(*u.ActorIDs) {
			if _, err := tx.Exec("INSERT INTO movie_actor (movie_id, actor_id) VALUES (?,?)", id, actorID); err != nil {
				if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
					return customerrors.NotFoundf("Actor with id %d not found", actorID)
				}
				return err
			}
		}
	}
	return tx.Commit()
}

func (r *MovieRepository) Delete(id int64, force bool) error {
	var genreCount int
	err := r.db.QueryRow("SELECT COUNT(*) FROM movie_genre WHERE movie_id = ?", id).Scan(&genreCount)
	if err != nil {
		return err
	}

	var actorCount int
	err = r.db.QueryRow("SELECT COUNT(*) FROM movie_actor WHERE movie_id = ?", id).Scan(&actorCount)
	if err != nil {
		return err
	}

	if genreCount > 0 || actorCount > 0 {
		if !force {
			return customerrors.Conflictf("Movie with id %d has associated genres or actors", id)
		}
	}

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
		return nil, err
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

func (r *MovieRepository) GetActorsByMovieID(movieID int64) ([]models.Actor, error) {
	rows, err := r.db.Query("SELECT a.id, a.name, a.birth_date FROM actors a JOIN movie_actor ma ON a.id = ma.actor_id WHERE ma.movie_id = ?", movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []models.Actor
	for rows.Next() {
		var actor models.Actor
		if err := rows.Scan(&actor.ID, &actor.Name, &actor.BirthDate); err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(actors) == 0 {
		return nil, customerrors.NotFoundf("No actors found for movie ID %d", movieID)
	}
	return actors, nil
}

func (r *MovieRepository) GetDetailByID(id int64) (models.MovieDetail, error) {
	var movie models.MovieDetail

	err := r.db.QueryRow("SELECT id, title, release_year, duration FROM movies WHERE id = ?", id).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Duration)
	if errors.Is(err, sql.ErrNoRows) {
		return models.MovieDetail{}, customerrors.NotFoundf("Movie with id %d not found", id)
	}
	if err != nil {
		return models.MovieDetail{}, err
	}

	genreRows, err := r.db.Query("SELECT genre_id FROM movie_genre WHERE movie_id = ?", id)
	if err != nil {
		return models.MovieDetail{}, err
	}
	defer genreRows.Close()

	var genreIDs []int64
	for genreRows.Next() {
		var genreID int64
		if err := genreRows.Scan(&genreID); err != nil {
			return models.MovieDetail{}, err
		}
		genreIDs = append(genreIDs, genreID)
	}
	if err := genreRows.Err(); err != nil {
		return models.MovieDetail{}, err
	}

	actorRows, err := r.db.Query("SELECT actor_id FROM movie_actor WHERE movie_id = ?", id)
	if err != nil {
		return models.MovieDetail{}, err
	}
	defer actorRows.Close()

	var actorIDs []int64
	for actorRows.Next() {
		var actorID int64
		if err := actorRows.Scan(&actorID); err != nil {
			return models.MovieDetail{}, err
		}
		actorIDs = append(actorIDs, actorID)
	}
	if err := actorRows.Err(); err != nil {
		return models.MovieDetail{}, err
	}

	movie.GenreIDs = genreIDs
	movie.ActorIDs = actorIDs
	return movie, nil
}
