package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/matt-godfrey/react-website/internal/database"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// DATABASE_URL=postgres://matt:yourpassword@localhost:5432/react_website?sslmode=disable
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	DATABASE_URL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := database.Connect(ctx, DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create router, handlers, services here

	log.Println("server running on :8080")
	http.ListenAndServe(":8080", nil)
}
