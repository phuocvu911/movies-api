package validation

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// V is the single shared validator instance for the whole server.
var V = newValidator()

// newValidator creates a new validator instance.
func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	v.RegisterValidation("pastdate", validatePastDate)
	return v
}

// validatePastDate validates that the field is a past date.
func validatePastDate(fl validator.FieldLevel) bool {
	t, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	return t.Before(time.Now())
}
