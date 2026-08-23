package handlers

import (
	"net/http"
)

// NewRouter is where we register endpoint to the mux, then I can wrap it around the middleware -> cleaner main.
func NewRouter(genres *GenreHandler, movies *MovieHandler, actors *ActorHandler) http.Handler {
	mux := http.NewServeMux()

	// register here
	//actors endpoints
	mux.HandleFunc("POST /api/actors", actors.Create)
	mux.HandleFunc("GET /api/actors", actors.GetAll)
	mux.HandleFunc("GET /api/actors/{id}", actors.GetByID)
	mux.HandleFunc("PATCH /api/actors/{id}", actors.Update)
	mux.HandleFunc("DELETE /api/actors/{id}", actors.Delete)
	mux.HandleFunc("GET /api/actors/{id}/movies", actors.Movies)

	//genres enpoints
	mux.HandleFunc("POST /api/genres", genres.Create)
	mux.HandleFunc("GET /api/genres", genres.GetAll)
	mux.HandleFunc("GET /api/genres/{id}", genres.GetByID)
	mux.HandleFunc("PATCH /api/genres/{id}", genres.Update)
	mux.HandleFunc("DELETE /api/genres/{id}", genres.Delete)
	mux.HandleFunc("GET /api/genres/{id}/movies", genres.Movies)

	//movies endpoints
	mux.HandleFunc("GET /api/movies", movies.GetAll)
	mux.HandleFunc("GET /api/movies/search", movies.Search)
	mux.HandleFunc("POST /api/movies", movies.Create)
	mux.HandleFunc("GET /api/movies/{id}", movies.GetByID)
	mux.HandleFunc("PATCH /api/movies/{id}", movies.Update)
	// mux.HandleFunc("DELETE /api/movies/{id}", movies.Delete)

	// wrap the mux with middleware
	return loggingMiddleware(rateLimit(requireJSON(mux)))
}
