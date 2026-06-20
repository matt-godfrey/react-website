package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matt-godfrey/react-website/internal/auth"
	appmiddleware "github.com/matt-godfrey/react-website/internal/middleware"
	"github.com/matt-godfrey/react-website/internal/quotes"
	"github.com/matt-godfrey/react-website/internal/sessions"
	"github.com/matt-godfrey/react-website/internal/users"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(chimiddleware.RequestID) // important for rate limiting
	r.Use(chimiddleware.RealIP)    // important for rate limiting and analytics
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// TODO: add proper allowed origins
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173",
			"https://mattgodfrey.xyz"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
		// AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(chimiddleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("All good"))
	})

	slidingWindowRateLimiter := appmiddleware.NewSlidingWindowRateLimiter(5, 1*time.Minute)

	userRepo := users.NewRepository(app.db)
	sessionRepo := sessions.NewRepository(app.db)
	quoteRepo := quotes.NewRepository(app.db, app.mongo)
	quoteService := quotes.NewService(quoteRepo)
	quoteHandler := quotes.NewHandler(quoteService)

	authService := auth.NewService(userRepo, sessionRepo)
	authHandler := auth.NewHandler(authService)

	r.Post("/register", authHandler.Register)

	r.With(appmiddleware.RateLimiterMiddleware(slidingWindowRateLimiter)).Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)
	r.Get("/auth/me", authHandler.GetCurrentUser)

	r.Get("/quotes/random", quoteHandler.GetRandomQuote)
	r.Get("/quotes/author", quoteHandler.GetAllQuotesByAuthor)
	return r

}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("listening on %s", app.config.addr)
	return srv.ListenAndServe()

}

type application struct {
	config config
	db     *pgxpool.Pool
	mongo  *mongo.Client
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn      string
	mongoURI string
}
