package quotes

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	db    *pgxpool.Pool
	mongo *mongo.Client
}

type mongoQuote struct {
	Q string `bson:"q"`
	A string `bson:"a"`
	C any    `bson:"c"`
	H string `bson:"h"`
}

func NewRepository(db *pgxpool.Pool, mongo *mongo.Client) *Repository {
	return &Repository{
		db:    db,
		mongo: mongo,
	}
}

func (r *Repository) FindAllQuotes(ctx context.Context) ([]*Quote, error) {

	quotesCollection := r.mongo.Database("quotes").Collection("quotes")

	// empty filter to get all docs
	filter := bson.M{}

	cursor, err := quotesCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var quotes []*Quote
	for cursor.Next(ctx) {
		var q Quote
		var mq mongoQuote
		if err := cursor.Decode(&mq); err != nil {
			return nil, err
		}
		q.Text = mq.Q
		q.Author = mq.A
		q.CharCount = mq.C
		q.Html = mq.H
		quotes = append(quotes, &q)
	}
	return quotes, nil
}

func (r *Repository) FindAllQuotesByAuthor(ctx context.Context, author string) ([]*Quote, error) {
	quotesCollection := r.mongo.Database("quotes").Collection("quotes")

	filter := bson.M{"a": author}

	cursor, err := quotesCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var quotes []*Quote
	for cursor.Next(ctx) {
		var q Quote
		var mq mongoQuote
		if err := cursor.Decode(&mq); err != nil {
			return nil, err
		}
		q.Text = mq.Q
		q.Author = mq.A
		q.CharCount = mq.C
		q.Html = mq.H
		quotes = append(quotes, &q)
	}
	return quotes, nil
}
