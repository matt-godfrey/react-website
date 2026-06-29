package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/matt-godfrey/react-website/internal/database"
	"github.com/matt-godfrey/react-website/internal/queue"
)

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
	MONGO_URI := os.Getenv("MONGO_URI")
	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn:      DATABASE_URL,
			mongoURI: MONGO_URI,
		},
	}
	db, err := database.Connect(ctx, cfg.db.dsn)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Connected to postgres database")

	mongoClient, err := database.NewMongoClient(cfg.db.mongoURI)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Connected to mongo database")

	rabbitMQUser := os.Getenv("RABBITMQ_USER")
	rabbitMQPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	rabbitMQVHost := os.Getenv("RABBITMQ_VHOST")

	rabbitConn, err := queue.ConnectRabbitMQ(rabbitMQUser, rabbitMQPassword, rabbitMQHost, rabbitMQVHost)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Connected to rabbitmq")

	defer rabbitConn.Close()

	api := application{
		config:     cfg,
		db:         db,
		mongo:      mongoClient,
		rabbitConn: rabbitConn,
	}
	defer db.Close()
	defer mongoClient.Disconnect(ctx)
	defer api.rabbitConn.Close()

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
