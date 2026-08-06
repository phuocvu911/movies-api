package service

import (
	"movies-api/internal/customerrors"
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

// GenreService implements business logic and validation for genres.
type GenreService struct {
	repo *repository.GenreRepository
}

func NewGenreService(genres *repository.GenreRepository) *GenreService {
	return &GenreService{repo: genres}
}

// Create new genre
func (s *GenreService) Create(input models.GenreRequest) (models.Genre, error) {
	//checking if the genre name is existed in db
	existingGenre, err := s.repo.GetByName(input.Name)
	if err == nil && existingGenre.ID != 0 {
		return models.Genre{}, customerrors.Conflictf("Genre with the same name already exists")
	}
	return s.repo.Create(input.Name)
}

// GetAll returns all genres.
func (s *GenreService) GetAll() ([]models.Genre, error) {
	return s.repo.GetAll()
}

// GetByID returns a genre by its ID.
func (s *GenreService) GetByID(id int64) (models.Genre, error) {
	return s.repo.GetByID(id)
}

// Movies returns all movies associated with a specific genre ID.
func (s *GenreService) Movies(id int64) ([]models.Movie, error) {
	// Check if the genre ID exists
	if _, err := s.repo.GetByID(id); err != nil {
		return nil, err
	}

	return s.repo.GetMoviesByGenreID(id)
}

// Update handles updating a genre by ID
func (s *GenreService) Update(id int64, input models.GenreRequest) (models.Genre, error) {
	return s.repo.Update(id, input)
}

// Delete handles deleting a genre by ID
func (s *GenreService) Delete(id int64, force bool) error {
	// Check if the genre ID exists
	genre, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	count, err := s.repo.MovieCount(id)
	if err != nil {
		return err
	}
	if count > 0 && !force {
		return customerrors.Validationf("Unable to delete genre '%s' as it is associated with %d movies", genre.Name, count)
	}

	// Delete the genre
	return s.repo.Delete(id)
}
