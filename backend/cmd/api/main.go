package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/matt-godfrey/react-website/internal/database"
)

func register(w http.ResponseWriter, r *http.Request) {
	// io.WriteString(w, "Hello from a HandleFunc #1!\n");
	fmt.Println("Hello, stdout")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "registered",
	})
}

func main() {
	ctx := context.Background()

	// Global logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file", err.Error(), 3)
	}

	// DATABASE_URL=postgres://matt:yourpassword@localhost:5432/react_website?sslmode=disable
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	DATABASE_URL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: DATABASE_URL,
		},
	}
	db, err := database.Connect(ctx, cfg.db.dsn)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Connected to database")

	api := application{
		config: cfg,
		db:     db,
	}
	defer db.Close()

	// handlers
	//

	// http.HandleFunc("/register", register)

	if err := api.run(api.mount()); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// log.Println("server running on :8080")
	// http.ListenAndServe(":8080", nil)
}
