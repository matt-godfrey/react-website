package quotes

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	db    *pgxpool.Pool
	mongo *mongo.Client
}

func NewRepository(db *pgxpool.Pool, mongo *mongo.Client) *Repository {
	return &Repository{
		db:    db,
		mongo: mongo,
	}
}
