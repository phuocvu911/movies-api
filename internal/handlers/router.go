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

	return loggingMiddleware(mux)
}
