package repository

import (
	"database/sql"
	"movies-api/internal/models"
)

// ActorRepository provides raw-SQL data access for actors.
type ActorRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{db: db}
}

func (r *ActorRepository) Create(name, birthDate string, movieIDs []int64) (models.Actor, error) {
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

	for _, movieID := range movieIDs {
		_, err := tx.Exec("INSERT INTO movie_actor (actor_id, movie_id) VALUES (?, ?)", actorID, movieID)
		if err != nil {
			return models.Actor{}, err
		}
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
