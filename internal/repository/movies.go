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

func (r *MovieRepository) GetAll(filter models.MovieFilter, limit, offset int) ([]models.Movie, int, error) {
	conditions := []string{}
	args := []any{}
	joins := ""

	if filter.Year != nil {
		conditions = append(conditions, "m.release_year = ?")
		args = append(args, *filter.Year)
	}

	if filter.GenreID != nil {
		joins += " JOIN movie_genre mg ON mg.movie_id = m.id"
		conditions = append(conditions, "mg.genre_id = ?")
		args = append(args, *filter.GenreID)
	}

	if filter.ActorID != nil {
		joins += " JOIN movie_actor ma ON ma.movie_id = m.id"
		conditions = append(conditions, "ma.actor_id = ?")
		args = append(args, *filter.ActorID)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT (*) FROM movies m" + joins + whereClause
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, customerrors.NotFoundf("No movie found")
	}

	if total < offset {
		return nil, 0, customerrors.NotFoundf("Page out of range")
	}

	query := "SELECT m.id, m.title, m.release_year, m.duration FROM movies m" + joins + whereClause + " ORDER BY m.id LIMIT ? OFFSET ?"
	queryArgs := append(args, limit, offset)
	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	movies := []models.Movie{}
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Duration); err != nil {
			return nil, 0, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

// Search searches for movies by title (partial match, case-insensitive).
func (r *MovieRepository) Search(title string, limit, offset int) ([]models.Movie, int, error) {
	var total int
	err := r.db.QueryRow("SELECT COUNT (*) FROM movies WHERE LOWER(title) LIKE ?", "%"+title+"%").Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	// If no movies found, return a NotFoundError
	if total == 0 {
		return nil, 0, customerrors.NotFoundf("No movie found")
	}

	if total < offset {
		return nil, 0, customerrors.NotFoundf("Page out of range")
	}
	rows, err := r.db.Query("SELECT id, title, release_year, duration FROM movies WHERE LOWER(title) LIKE ? LIMIT ? OFFSET ?", "%"+title+"%", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	movies := []models.Movie{}
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Duration); err != nil {
			return nil, 0, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return movies, total, nil
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
			return customerrors.Validationf("Movie with id %d has associated genres or actors", id)
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

func (r *MovieRepository) GetActorsByMovieID(movieID int64, limit, offset int) ([]models.Actor, int, error) {
	var total int
	err := r.db.QueryRow("SELECT COUNT (*) FROM movie_actor WHERE movie_id = ?", movieID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, customerrors.NotFoundf("No actor found")
	}

	if total < offset {
		return nil, 0, customerrors.NotFoundf("Page out of range")
	}

	rows, err := r.db.Query("SELECT a.id, a.name, a.birth_date FROM actors a JOIN movie_actor ma ON a.id = ma.actor_id WHERE ma.movie_id = ? LIMIT ? OFFSET ?", movieID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var actors []models.Actor
	for rows.Next() {
		var actor models.Actor
		if err := rows.Scan(&actor.ID, &actor.Name, &actor.BirthDate); err != nil {
			return nil, 0, err
		}
		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return actors, total, nil
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
