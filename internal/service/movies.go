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

func (s *MovieService) GetAll(filter models.MovieFilter, page, size int) (models.Page[models.Movie], error) {
	offset := page * size
	movies, total, err := s.repo.GetAll(filter, size, offset)
	if err != nil {
		return models.Page[models.Movie]{}, err
	}
	return newPage(movies, page, size, total), nil
}

// Search searches for movies by title.
func (s *MovieService) Search(title string, page, size int) (models.Page[models.Movie], error) {
	offset := page * size
	movies, total, err := s.repo.Search(title, size, offset)
	if err != nil {
		return models.Page[models.Movie]{}, err
	}
	return newPage(movies, page, size, total), nil
}

func (s *MovieService) Create(input models.MovieRequest) (models.Movie, error) {
	return s.repo.Create(input)
}

func (s *MovieService) GetByID(id int64) (models.Movie, error) {
	return s.repo.GetByID(id)
}

func (s *MovieService) Update(id int64, u models.MoviePatch) (models.MovieDetail, error) {
	if err := s.repo.Update(id, u); err != nil {
		return models.MovieDetail{}, err
	}
	return s.repo.GetDetailByID(id)
}

func (s *MovieService) Delete(id int64, force bool) error {
	return s.repo.Delete(id, force)
}

func (s *MovieService) GetActorsByMovieID(movieID int64, page, size int) (models.Page[models.Actor], error) {
	offset := page * size
	actors, total, err := s.repo.GetActorsByMovieID(movieID, size, offset)
	if err != nil {
		return models.Page[models.Actor]{}, err
	}
	return newPage(actors, page, size, total), nil
}
