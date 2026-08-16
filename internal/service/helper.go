package service

import "movies-api/internal/models"

// newPage builds a Page from a slice of results and pagination parameters.
func newPage[T any](results []T, page, size, total int) models.Page[T] {
	if results == nil {
		results = []T{}
	}
	totalPages := 0
	if size > 0 {
		totalPages = (total + size - 1) / size
	}
	return models.Page[T]{
		Results: results,
		Pagination: models.PageInfo{
			Page:          page,
			Size:          size,
			TotalElements: total,
			TotalPages:    totalPages,
		},
	}
}
