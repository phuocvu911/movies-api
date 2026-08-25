package service

import (
	"movies-api/internal/customerrors"
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
	movies, total, err := s.repo.GetAll(filter, size, page*size)
	if err != nil {
		return models.Page[models.Movie]{}, err
	}
	return newPage(movies, page, size, total), nil
}

// Search searches for movies by title.
func (s *MovieService) Search(title string, page, size int) (models.Page[models.Movie], error) {
	movies, total, err := s.repo.Search(title, size, page*size)
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

func (s *MovieService) Update(id int64, u models.MoviePatchRequest) (models.MoviePatchRequest, error) {
	err := s.repo.Update(id, u)
	if err != nil {
		return models.MoviePatchRequest{}, err
	}
	return s.repo.GetByIDForPatch(id)
}

func (s *MovieService) Delete(id int64, force bool) error {
	movie, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	genreCount, err := s.repo.GenreCount(id)
	if err != nil {
		return err
	}
	actorCount, err := s.repo.ActorCount(id)
	if err != nil {
		return err
	}

	if !force && (genreCount > 0 || actorCount > 0) {
		return customerrors.Validationf("Unable to delete the movie '%s' as it is associated with %d genres and %d actors", movie.Title, genreCount, actorCount)
	}

	return s.repo.Delete(id)
}

func (s *MovieService) Actors(id int64) ([]models.Actor, error) {
	return s.repo.Actors(id)
}
