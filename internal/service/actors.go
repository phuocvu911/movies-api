package service

import (
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

// ActorInput is the validated(clean) input for creating an actor.
type ActorInput struct {
	Name      string
	BirthDate string
	MovieIDs  []int64
}

// ActorService implements business logic and validation for actors.
type ActorService struct {
	repo *repository.ActorRepository
	//movies *repository.MovieRepository
}

func NewActorService(actors *repository.ActorRepository /*, movies *repository.MovieRepository*/) *ActorService {
	return &ActorService{repo: actors /*movies: movies*/}
}

func (s *ActorService) Create(input ActorInput) (models.Actor, error) {
	//check if movie ids exist in the movies table when its movie repo ready
	return s.repo.Create(input.Name, input.BirthDate, input.MovieIDs)
}
