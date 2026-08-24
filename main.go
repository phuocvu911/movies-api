package main

import (
	"fmt"
	"log"
	"movies-api/internal/database"
	"movies-api/internal/handlers"
	"movies-api/internal/repository"
	"movies-api/internal/service"
	"net/http"
)

func main() {
	//open db
	db, err := database.Open()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	//migrate
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	//seeding data
	if err := database.Seed(db); err != nil {
		log.Fatalf("database seeding failed: %v", err)
	}

	//repo init
	actorRepo := repository.NewActorRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	movieRepo := repository.NewMovieRepository(db)

	//service init
	actorService := service.NewActorService(actorRepo /*movieRepo*/)
	genreService := service.NewGenreService(genreRepo)
	movieService := service.NewMovieService(movieRepo)

	//handlers through newRouter()
	router := handlers.NewRouter(
		handlers.NewGenreHandler(genreService),
		handlers.NewMovieHandler(movieService),
		handlers.NewActorHandler(actorService),
	)

	// Start the cleanup goroutine for rate limiters
	go handlers.CleanupVisitors()

	//ListenandServe()
	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
