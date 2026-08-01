package handlers

import (
	"net/http"
)

// NewRouter is where we register endpoint to the mux, then I can wrap it around the middleware -> cleaner main.
func NewRouter( /*genres *GenreHandler, movies *MovieHandler,*/ actors *ActorHandler) http.Handler {
	mux := http.NewServeMux()

	// register here
	//actors endpoints
	mux.HandleFunc("POST /api/actors", actors.Create)
	mux.HandleFunc("GET /api/actors", actors.GetAll)
	mux.HandleFunc("GET /api/actors/{id}", actors.GetByID)
	mux.HandleFunc("PATCH /api/actors/{id}", actors.Update)
	mux.HandleFunc("DELETE /api/actors/{id}", actors.Delete)

	return loggingMiddleware(mux)
}
