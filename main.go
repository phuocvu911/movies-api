package main

import (
	"log"
	"movies-api/internal/database"
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

	//seeding data()
	if err := database.Seed(db); err != nil {
		log.Fatalf("database seeding failed: %v", err)
	}

	//repo init

	//service init

	//handlers through newRouter()

	//ListenandServe()
}
