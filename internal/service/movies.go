package service

import (
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

// MovieService provides business logic for movie-related operations
type MovieService struct {
	repo *repository.MovieRepository
}

// NewMovieService creates a new MovieService
func NewMovieService(repo *repository.MovieRepository) *MovieService {
	return &MovieService{repo: repo}
}

// GetAll returns a paginated list of movies and supports optional filters
// by genres, release year and actor.
func (s *MovieService) GetAll(filter models.MovieFilter, page, size int) (models.Page[models.Movie], error) {
	offset := page * size
	movies, total, err := s.repo.GetAll(filter, size, offset)
	if err != nil {
		return models.Page[models.Movie]{}, err
	}
	return newPage(movies, page, size, total), nil
}

// Search returns a paginated list of movies matching a given title.
func (s *MovieService) Search(title string, page, size int) (models.Page[models.Movie], error) {
	offset := page * size
	movies, total, err := s.repo.Search(title, size, offset)
	if err != nil {
		return models.Page[models.Movie]{}, err
	}
	return newPage(movies, page, size, total), nil
}

// Create creates a new movie with movie's details and its associated genres and actors.
func (s *MovieService) Create(input models.MovieRequest) (models.Movie, error) {
	return s.repo.Create(input)
}

// GetByID returns a movie by its ID
func (s *MovieService) GetByID(id int64) (models.Movie, error) {
	return s.repo.GetByID(id)
}

// Update updates the provided movie fields and relationships.
func (s *MovieService) Update(id int64, u models.MoviePatch) (models.MovieDetail, error) {
	if err := s.repo.Update(id, u); err != nil {
		return models.MovieDetail{}, err
	}
	return s.repo.GetDetailByID(id)
}

// Delete deletes a movie. The force query parameter can be used
// when the movie has associated genres and actors.
func (s *MovieService) Delete(id int64, force bool) error {
	return s.repo.Delete(id, force)
}

// GetActorsByMovieID returns a paginated list of actors associated with the movie
func (s *MovieService) GetActorsByMovieID(movieID int64, page, size int) (models.Page[models.Actor], error) {
	offset := page * size
	actors, total, err := s.repo.GetActorsByMovieID(movieID, size, offset)
	if err != nil {
		return models.Page[models.Actor]{}, err
	}
	return newPage(actors, page, size, total), nil
}
