package service

import (
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

type MovieService struct {
	repo *repository.MovieRepository
}

func NewMovieService(repo *repository.MovieRepository) *MovieService {
	return &MovieService{repo: repo}
}

func (s *MovieService) GetAll() ([]models.Movie, error) {
	return s.repo.GetAll()
}

// Search searches for movies by title.
func (s *MovieService) Search(title string) ([]models.Movie, error) {
	return s.repo.Search(title)
}

func (s *MovieService) Create(input models.MovieRequest) (models.Movie, error) {
	return s.repo.Create(input)
}

func (s *MovieService) GetByID(id int64) (models.Movie, error) {
	return s.repo.GetByID(id)
}
