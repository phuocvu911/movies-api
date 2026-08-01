package service

import (
	"movies-api/internal/customerrors"
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

// ActorInput is the validated(clean) input for creating an actor.
type ActorInput struct {
	Name      string
	BirthDate string
}

// ActorPatch is the input for partially updating an actor; nil fields are left unchanged.
type ActorPatch struct {
	Name      *string
	BirthDate *string
	MovieIDs  *[]int64
}

// ActorService implements business logic and validation for actors.
type ActorService struct {
	repo *repository.ActorRepository
	//movies *repository.MovieRepository
}

func NewActorService(actors *repository.ActorRepository /*, movies *repository.MovieRepository*/) *ActorService {
	return &ActorService{repo: actors /*movies: movies*/}
}

// Create creates a new actor and associates it with the given movie IDs.
func (s *ActorService) Create(input ActorInput) (models.Actor, error) {
	//check if actor name and birthdate already exist in the actors table
	existingActor, err := s.repo.GetByNameAndBirthDate(input.Name, input.BirthDate)
	if err == nil && existingActor.ID != 0 {
		return models.Actor{}, customerrors.Conflictf("actor with the same name and birth date already exists")
	}
	return s.repo.Create(input.Name, input.BirthDate)
}

// GetAll returns all actors.
func (s *ActorService) GetAll() ([]models.Actor, error) {
	return s.repo.GetAll()
}

// GetByName returns actors matching the given name (case-insensitive)
func (s *ActorService) GetByName(name string) ([]models.Actor, error) {
	return s.repo.GetByName(name)
}

// GetByID returns an actor by ID.
func (s *ActorService) GetByID(id int64) (models.Actor, error) {
	return s.repo.GetByID(id)
}

// Update validates and applies a partial update, returning the updated actor.
func (s *ActorService) Update(id int64, p models.ActorPatch) (models.ActorPatch, error) {
	//check if the actor exists
	if _, err := s.repo.GetByID(id); err != nil {
		return models.ActorPatch{}, err
	}
	//check if the movie ids exist in the movies table
	if p.MovieIDs != nil {
		if err := s.checkMovieIDs(*p.MovieIDs); err != nil {
			return models.ActorPatch{}, err
		}
	}

	if err := s.repo.Update(id, p); err != nil {
		return models.ActorPatch{}, err
	}
	return s.repo.GetByIDForPatch(id)
}

// checkMovieIDs checks if all movie IDs in request body exist in the movies table.
func (s *ActorService) checkMovieIDs(movieIDs []int64) error {
	for _, movieID := range movieIDs {
		if _, err := s.repo.GetMovieByID(movieID); err != nil {
			return err
		}
	}
	return nil
}

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
		return customerrors.Validationf("Unable to delete actor '%s' as they are associated with %d movies", actor.Name, count)
	}
	return s.repo.Delete(id)
}

// Movies returns all movies associated with an actor.
func (s *ActorService) Movies(actorID int64) ([]models.Movie, error) {
	//check if the actor exists
	if _, err := s.repo.GetByID(actorID); err != nil {
		return nil, err
	}
	return s.repo.GetMoviesByActorID(actorID)
}
