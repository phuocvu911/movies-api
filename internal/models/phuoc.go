// models that is widely used accross the codebase got registered here
package models

// Actor represents an actor. BirthDate uses ISO 8601 (YYYY-MM-DD).
type Actor struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}
