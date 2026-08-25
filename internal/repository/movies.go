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

func (r *MovieRepository) GetAll(filter models.MovieFilter) ([]models.Movie, error) {
	var command string
	var arg int64

	if filter.ActorID != nil {
		command = `
			SELECT id, title, release_year, duration 
			FROM movies 
			INNER JOIN movie_actor ON movies.id = movie_actor.movie_id 
			WHERE movie_actor.actor_id = ?`
		arg = *filter.ActorID
	} else if filter.GenreID != nil {
		command = `
			SELECT id, title, release_year, duration 
			FROM movies 
			INNER JOIN movie_genre ON movies.id = movie_genre.movie_id 
			WHERE movie_genre.genre_id = ?`
		arg = *filter.GenreID
	} else if filter.Year != nil {
		command = `
			SELECT id, title, release_year, duration 
			FROM movies 
			WHERE release_year = ?`
		arg = *filter.Year
	} else {
		command = `SELECT id, title, release_year, duration FROM movies`
	}

	rows, err := r.db.Query(command, arg)
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
	if err = rows.Err(); err != nil {
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
	rows, err := r.db.Query(`SELECT id, title, release_year, duration FROM movies WHERE LOWER(title) LIKE ?`, "%"+title+"%")
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
	if err = rows.Err(); err != nil {
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

	result, err := tx.Exec(`INSERT INTO movies (title, release_year, duration) VALUES (?, ?, ?)`, input.Title, input.Year, input.Duration)
	if err != nil {
		return models.Movie{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Movie{}, err
	}

	//add relationships
	for _, genreID := range dedupe(input.GenreIDs) {
		_, err := tx.Exec(`INSERT INTO movie_genre (movie_id, genre_id) VALUES (?, ?)`, id, genreID)
		if err != nil {
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				return models.Movie{}, customerrors.NotFoundf("Genre with ID %d not found", genreID)
			}
			return models.Movie{}, err
		}
	}

	for _, actorID := range dedupe(input.ActorIDs) {
		_, err := tx.Exec(`INSERT INTO movie_actor (movie_id, actor_id) VALUES (?, ?)`, id, actorID)
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

func (r *MovieRepository) GetByID(id int64) (models.Movie, error) {
	var movie models.Movie
	err := r.db.QueryRow(`SELECT id, title, release_year, duration FROM movies WHERE id = ?`, id).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Duration)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Movie{}, customerrors.NotFoundf("Movie with ID %d not found", id)
	}
	return movie, err
}

func (r *MovieRepository) Update(id int64, p models.MoviePatchRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//update movies table
	var sets []string
	var args []any

	if p.Title != nil {
		sets = append(sets, "title = ?")
		args = append(args, *p.Title)
	}

	if p.Year != nil {
		sets = append(sets, "release_year = ?")
		args = append(args, *p.Year)
	}

	if p.Duration != nil {
		sets = append(sets, "duration = ?")
		args = append(args, *p.Duration)
	}

	if len(sets) > 0 {
		args = append(args, id)
		result, err := tx.Exec(`UPDATE movies SET `+strings.Join(sets, ", ")+` WHERE id = ?`, args...)
		if err != nil {
			return err
		}
		if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
			return customerrors.NotFoundf("Movie with ID %d not found", id)
		}
	}

	//update relationships
	if p.GenreIDs != nil {
		//delete all existing genres
		_, err := tx.Exec(`DELETE FROM movie_genre WHERE movie_id = ?`, id)
		if err != nil {
			return err
		}
		//insert new genres
		for _, genreID := range dedupe(*p.GenreIDs) {
			_, err := tx.Exec(`INSERT INTO movie_genre (movie_id, genre_id) VALUES (?, ?)`, id, genreID)
			if err != nil {
				if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
					return customerrors.NotFoundf("Genre with ID %d not found", genreID)
				}
				return err
			}
		}
	}

	if p.ActorIDs != nil {
		//delete all existing actors
		_, err := tx.Exec(`DELETE FROM movie_actor WHERE movie_id = ?`, id)
		if err != nil {
			return err
		}
		//insert new actors
		for _, actorID := range dedupe(*p.ActorIDs) {
			_, err := tx.Exec(`INSERT INTO movie_actor (movie_id, actor_id) VALUES (?, ?)`, id, actorID)
			if err != nil {
				if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
					return customerrors.NotFoundf("Actor with ID %d not found", actorID)
				}
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *MovieRepository) GetByIDForPatch(id int64) (models.MoviePatchRequest, error) {
	var movie models.MoviePatchRequest
	err := r.db.QueryRow(`SELECT id, title, release_year, duration FROM movies WHERE id = ?`, id).Scan(&movie.Id, &movie.Title, &movie.Year, &movie.Duration)
	rows, err := r.db.Query(`SELECT genre_id FROM movie_genre WHERE movie_id = ?`, id)
	if err != nil {
		return models.MoviePatchRequest{}, err
	}
	defer rows.Close()

	var genreIDs []int64
	for rows.Next() {
		var genreID int64
		if err := rows.Scan(&genreID); err != nil {
			return models.MoviePatchRequest{}, err
		}
		genreIDs = append(genreIDs, genreID)
	}
	if err = rows.Err(); err != nil {
		return models.MoviePatchRequest{}, err
	}
	if len(genreIDs) > 0 {
		movie.GenreIDs = &genreIDs
	}

	rows, err = r.db.Query(`SELECT actor_id FROM movie_actor WHERE movie_id = ?`, id)
	if err != nil {
		return models.MoviePatchRequest{}, err
	}
	defer rows.Close()

	var actorIDs []int64
	for rows.Next() {
		var actorID int64
		if err := rows.Scan(&actorID); err != nil {
			return models.MoviePatchRequest{}, err
		}
		actorIDs = append(actorIDs, actorID)
	}
	if err = rows.Err(); err != nil {
		return models.MoviePatchRequest{}, err
	}
	if len(actorIDs) > 0 {
		movie.ActorIDs = &actorIDs
	}

	return movie, nil
}

func (r *MovieRepository) GenreCount(movieID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM movie_genre WHERE movie_id = ?`, movieID).Scan(&count)
	return count, err
}

func (r *MovieRepository) ActorCount(movieID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM movie_actor WHERE movie_id = ?`, movieID).Scan(&count)
	return count, err
}

func (r *MovieRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM movies WHERE id = ?`, id)
	return err
}

func (r *MovieRepository) Actors(movieID int64) ([]models.Actor, error) {
	rows, err := r.db.Query(`
		SELECT
			actors.id,
			actors.name,
			actors.birth_date
		FROM movie_actor
		JOIN actors ON movie_actor.actor_id = actors.id
		WHERE movie_actor.movie_id = ?
		ORDER BY actor_id`, movieID)
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
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(actors) == 0 {
		return nil, customerrors.NotFoundf("No actors found for movie with ID %d", movieID)
	}
	return actors, nil
}
