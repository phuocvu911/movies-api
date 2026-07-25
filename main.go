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

	//repo init

	//service init

	//handlers through newRouter()

	//ListenandServe()
}
