package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"movies-api/internal/customerrors"
	"net/http"
	"strconv"
)

// writeJSON serializes v to the response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("write response error: %v", err)
		}
	}
}

// decodeJSON decode json body into struct variable
func decodeJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return customerrors.Validationf("invalid request body: %v", err)
	}
	return nil
}

// respondError is the centralized error handler: it maps the custom error
// types to HTTP status codes and a JSON error body.
func respondError(w http.ResponseWriter, err error) {
	var validationErr *customerrors.ValidationError
	var conflictErr *customerrors.ConflictError

	switch {
	case errors.As(err, &validationErr):
		writeJSON(w, http.StatusBadRequest, validationErr.Message)
	case errors.As(err, &conflictErr):
		writeJSON(w, http.StatusConflict, conflictErr.Message)
	case errors.Is(err, customerrors.ErrNotFound):
		writeJSON(w, http.StatusNotFound, err.Error())
	default:
		//log.Printf("internal error: %v", err)
		writeJSON(w, http.StatusInternalServerError, "internal server error")
	}
}

// pathID parses the {id} path segment.
func pathID(r *http.Request) (int64, error) {
	raw := r.PathValue("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, customerrors.Validationf("invalid id '%s': must be a positive integer", raw)
	}
	return id, nil
}

// forceParam parses the optional force query parameter for delete handlers.
func forceParam(r *http.Request) (bool, error) {
	raw := r.URL.Query().Get("force")
	if raw == "" {
		return false, nil
	}
	force, err := strconv.ParseBool(raw)
	if err != nil {
		return false, customerrors.Validationf("invalid force '%s': must be true or false", raw)
	}
	return force, nil
}

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

// pagination parses the custom page and size query parameters
func pagination(r *http.Request) (int, int, error) {
	page := 0
	size := defaultPageSize
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			return 0, 0, customerrors.Validationf("Invalid page '%s': must be a non-negative integer", pageStr)
		}
	}

	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		var err error
		size, err = strconv.Atoi(sizeStr)
		if err != nil || size < 1 || size > maxPageSize {
			return 0, 0, customerrors.Validationf("Invalid size '%s': must be between 1 and %d", sizeStr, maxPageSize)
		}
	}

	return page, size, nil
}
