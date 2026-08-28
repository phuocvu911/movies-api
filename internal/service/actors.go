package service

import (
	"movies-api/internal/customerrors"
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

// ActorService implements business logic and validation for actors.
type ActorService struct {
	repo *repository.ActorRepository
	//movies *repository.MovieRepository
}

func NewActorService(actors *repository.ActorRepository /*, movies *repository.MovieRepository*/) *ActorService {
	return &ActorService{repo: actors /*movies: movies*/}
}

// Create creates a new actor.
func (s *ActorService) Create(input models.ActorRequest) (models.Actor, error) {
	//check if actor name and birthdate already exist in the actors table
	existingActor, err := s.repo.GetByNameAndBirthDate(input.Name, input.BirthDate)
	if err == nil && existingActor.ID != 0 {
		return models.Actor{}, customerrors.Conflictf("Actor with the same name and birth date already exists")
	}
	return s.repo.Create(input.Name, input.BirthDate)
}

// GetAll returns all actors.
func (s *ActorService) GetAll(page, size int) (models.Page[models.Actor], error) {
	offset := page * size
	actors, total, err := s.repo.GetAll(size, offset)
	if err != nil {
		return models.Page[models.Actor]{}, err
	}
	return newPage(actors, page, size, total), nil
}

// GetByName returns actors matching the given name (case-insensitive)
func (s *ActorService) GetByName(name string, page, size int) (models.Page[models.Actor], error) {
	offset := page * size
	actors, total, err := s.repo.GetByName(name, size, offset)
	if err != nil {
		return models.Page[models.Actor]{}, err
	}
	return newPage(actors, page, size, total), nil
}

// GetByID returns an actor by ID.
func (s *ActorService) GetByID(id int64) (models.Actor, error) {
	return s.repo.GetByID(id)
}

// Update validates and applies a partial update, returning the updated actor.
func (s *ActorService) Update(id int64, p models.ActorPatch) (models.ActorPatch, error) {
	if err := s.repo.Update(id, p); err != nil {
		return models.ActorPatch{}, err
	}
	//return the updated actor with the latest data from the database
	return s.repo.GetByIDForPatch(id)
}

// checkMovieIDs checks if all movie IDs in request body exist in the movies table.
// func (s *ActorService) checkMovieIDs1(movieIDs []int64) error {
// 	for _, movieID := range movieIDs {
// 		if _, err := s.repo.GetMovieByID(movieID); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// Delete removes an actor by ID. Normally it should fail if the actor is associated with any movies, unless force is true, in which case it removes the associations first.
func (s *ActorService) Delete(id int64, force bool) error {
	actor, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	count, err := s.repo.MovieCount(id)
	if err != nil {
		return err
	}
	if count > 0 && !force {
		return customerrors.Validationf("Unable to delete actor '%s' as he/she is associated with %d movies", actor.Name, count)
	}
	return s.repo.Delete(id)
}

// Movies returns all movies associated with an actor.
func (s *ActorService) Movies(actorID int64, page, size int) (models.Page[models.Movie], error) {
	//check if the actor exists
	if _, err := s.repo.GetByID(actorID); err != nil {
		return models.Page[models.Movie]{}, err
	}
	offset := page * size
	movies, total, err := s.repo.GetMoviesByActorID(actorID, size, offset)
	if err != nil {
		return models.Page[models.Movie]{}, err
	}
	return newPage(movies, page, size, total), nil
}
