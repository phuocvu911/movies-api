package handlers

import (
	"log"
	"net/http"
	"time"
)

// NewRouter is where we register endpoint to the mux, then I can wrap it around the middleware -> cleaner main.
func NewRouter(genres *GenreHandler, movies *MovieHandler, actors *ActorHandler) http.Handler {
	mux := http.NewServeMux()

	// register here

	return loggingMiddleware(mux)
}

// loggingMiddleware logs each request's method, path and duration.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}
