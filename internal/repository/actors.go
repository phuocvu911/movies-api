package repository

import (
	"database/sql"
	"errors"
	"movies-api/internal/customerrors"
	"movies-api/internal/models"
	"strings"
)

// ActorUpdate describes a partial update; nil fields are left unchanged.
type ActorUpdate struct {
	Name      *string
	BirthDate *string
	MovieIDs  *[]int64
}

// ActorRepository provides raw-SQL data access for actors.
type ActorRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{db: db}
}

// Create creates a new actor and associates it with the given movie IDs.
func (r *ActorRepository) Create(name, birthDate string) (models.Actor, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.Actor{}, err
	}
	defer tx.Rollback()

	res, err := tx.Exec("INSERT INTO actors (name, birth_date) VALUES (?, ?)", name, birthDate)
	if err != nil {
		return models.Actor{}, err
	}

	actorID, err := res.LastInsertId()
	if err != nil {
		return models.Actor{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Actor{}, err
	}

	return models.Actor{
		ID:        int64(actorID),
		Name:      name,
		BirthDate: birthDate,
	}, nil
}

// GetAll returns all actors.
func (r *ActorRepository) GetAll() ([]models.Actor, error) {
	rows, err := r.db.Query("SELECT id, name, birth_date FROM actors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	actors := []models.Actor{}
	for rows.Next() {
		var actor models.Actor
		if err := rows.Scan(&actor.ID, &actor.Name, &actor.BirthDate); err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}
	return actors, nil
}

// GetByName returns actors matching the given name (partially, case-insensitive)
func (r *ActorRepository) GetByName(name string) ([]models.Actor, error) {
	rows, err := r.db.Query("SELECT id, name, birth_date FROM actors WHERE LOWER(name) LIKE ?", "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	actors := []models.Actor{}
	for rows.Next() {
		var actor models.Actor
		if err := rows.Scan(&actor.ID, &actor.Name, &actor.BirthDate); err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}

	// If no actors found, return a NotFoundError
	if len(actors) == 0 {
		return nil, customerrors.NotFoundf("no actors found with name %s", name)
	}
	return actors, nil
}

// GetByID returns a single actor or a NotFoundError.
func (r *ActorRepository) GetByID(id int64) (models.Actor, error) {
	var actor models.Actor
	err := r.db.QueryRow(`SELECT id, name, birth_date FROM actors WHERE id = ?`, id).Scan(&actor.ID, &actor.Name, &actor.BirthDate)

	//if we dont have actor for that id
	if errors.Is(err, sql.ErrNoRows) {
		return models.Actor{}, customerrors.NotFoundf("actor with id %d not found", id)
	}

	return actor, err
}

// GetByIDForPatch returns a single actor or a NotFoundError, including associated movie IDs.
func (r *ActorRepository) GetByIDForPatch(id int64) (models.ActorPatch, error) {
	var actor models.ActorPatch
	err := r.db.QueryRow(`SELECT id, name, birth_date FROM actors WHERE id = ?`, id).Scan(&actor.Id, &actor.Name, &actor.BirthDate)

	//if we dont have actor for that id
	if errors.Is(err, sql.ErrNoRows) {
		return models.ActorPatch{}, customerrors.NotFoundf("actor with id %d not found", id)
	}

	// get the movie ids associated with the actor
	rows, err := r.db.Query(`SELECT movie_id FROM movie_actor WHERE actor_id = ?`, id)
	if err != nil {
		return models.ActorPatch{}, err
	}
	defer rows.Close()

	var movieIDs []int64
	for rows.Next() {
		var movieID int64
		if err := rows.Scan(&movieID); err != nil {
			return models.ActorPatch{}, err
		}
		movieIDs = append(movieIDs, movieID)
	}
	if len(movieIDs) > 0 {
		actor.MovieIDs = &movieIDs
	}

	return actor, err
}

// GetByNameAndBirthDate returns an actor by name and birth date (to check for duplicate actor that client try to write to db)
func (r *ActorRepository) GetByNameAndBirthDate(name, birthDate string) (models.Actor, error) {
	row := r.db.QueryRow("SELECT id, name, birth_date FROM actors WHERE name = ? AND birth_date = ?", name, birthDate)
	var actor models.Actor
	if err := row.Scan(&actor.ID, &actor.Name, &actor.BirthDate); err != nil {
		return models.Actor{}, err
	}
	return actor, nil
}

// Update updates an actor and its associated movie IDs.
func (r *ActorRepository) Update(id int64, u models.ActorPatch) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var sets []string
	var args []any // string and int mixed, so use any instead of string
	if u.Name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *u.Name)
	}
	if u.BirthDate != nil {
		sets = append(sets, "birth_date = ?")
		args = append(args, *u.BirthDate)
	}
	if len(sets) > 0 {
		args = append(args, id)
		if _, err := tx.Exec(`UPDATE actors SET `+strings.Join(sets, ", ")+` WHERE id = ?`, args...); err != nil {
			return err
		}
	}

	if u.MovieIDs != nil {
		// Clear existing relastionships
		if _, err := tx.Exec(`DELETE FROM movie_actor WHERE actor_id = ?`, id); err != nil {
			return err
		}
		// Insert new relationships
		for _, m_id := range dedupe(*u.MovieIDs) {
			if _, err := tx.Exec(`INSERT INTO movie_actor (movie_id, actor_id) VALUES (?, ?)`, m_id, id); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// Delete removes an actor and its associated movie relationships, with force true already specified.
func (r *ActorRepository) Delete(id int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete relationships in movie_actor table
	if _, err := tx.Exec(`DELETE FROM movie_actor WHERE actor_id = ?`, id); err != nil {
		return err
	}

	// Delete the actor
	if _, err := tx.Exec(`DELETE FROM actors WHERE id = ?`, id); err != nil {
		return err
	}

	return tx.Commit()
}

// move this to movie repository later
// GetMovieByID returns a single movie by ID or a NotFoundError.
func (r *ActorRepository) GetMovieByID(id int64) (int, error) {
	var movie_id int
	err := r.db.QueryRow(`SELECT id FROM movies WHERE id = ?`, id).Scan(&movie_id)

	//if we dont have movie for that id
	if errors.Is(err, sql.ErrNoRows) {
		return 0, customerrors.NotFoundf("movie with id %d not found", id)
	}
	return movie_id, err
}

// MovieCount returns how many movies are linked to the actor.
func (r *ActorRepository) MovieCount(id int64) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM movie_actor WHERE actor_id = ?`, id).Scan(&count)
	return count, err
}

// GetMoviesByActorID returns all movies associated with an actor.
func (r *ActorRepository) GetMoviesByActorID(actorID int64) ([]models.Movie, error) {
	rows, err := r.db.Query(`
		SELECT m.id, m.title, m.release_year, m.duration
		FROM movies m
		JOIN movie_actor ma ON m.id = ma.movie_id
		WHERE ma.actor_id = ?`, actorID)
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
	if len(movies) == 0 {
		return nil, customerrors.NotFoundf("no movies found for actor with id %d", actorID)
	}

	return movies, nil
}
