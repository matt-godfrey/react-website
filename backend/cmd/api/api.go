package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	httprateredis "github.com/go-chi/httprate-redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matt-godfrey/react-website/internal/auth"
	"github.com/matt-godfrey/react-website/internal/mailer"
	appmiddleware "github.com/matt-godfrey/react-website/internal/middleware"
	"github.com/matt-godfrey/react-website/internal/queue"
	"github.com/matt-godfrey/react-website/internal/quotes"
	"github.com/matt-godfrey/react-website/internal/sessions"
	"github.com/matt-godfrey/react-website/internal/users"
	amqp "github.com/rabbitmq/amqp091-go"
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

	// slidingWindowRateLimiter := appmiddleware.NewSlidingWindowRateLimiter(5, 1*time.Minute)

	host := os.Getenv("VALKEY_HOST")
	port, err := strconv.Atoi(os.Getenv("VALKEY_PORT"))
	if err != nil {
		port = 6379
	}

	dscfg := httprateredis.Config{Host: host, Port: uint16(port)}
	keyFuncs := httprate.WithKeyFuncs(httprate.KeyByIP)
	chiRateLimiter := appmiddleware.NewChiRateLimiter(3, 10*time.Second, keyFuncs, httprateredis.WithRedisLimitCounter(&dscfg))

	userRepo := users.NewRepository(app.db)
	sessionRepo := sessions.NewRepository(app.db)
	quoteRepo := quotes.NewRepository(app.db, app.mongo)
	quoteService := quotes.NewService(quoteRepo)
	quoteHandler := quotes.NewHandler(quoteService)

	from := os.Getenv("RESEND_FROM")
	mailer := mailer.NewResendMailer(os.Getenv("RESEND_API_KEY"), from)

	rabbitClient, err := queue.NewRabbitClient(app.rabbitConn)
	if err != nil {
		log.Fatal(err)
	}

	err = rabbitClient.CreateQueue("email", true, false)
	err = rabbitClient.CreateExchange("email", "direct", true, false)
	err = rabbitClient.CreateBinding("email", "email", "email")
	if err != nil {
		log.Fatal(err)
	}

	// defer rabbitClient.Close()

	authService := auth.NewService(userRepo, sessionRepo, mailer, rabbitClient)
	authHandler := auth.NewHandler(authService)

	r.Post("/register", authHandler.Register)

	// r.With(appmiddleware.RateLimiterMiddleware(slidingWindowRateLimiter)).Post("/login", authHandler.Login)
	r.With(appmiddleware.ChiRateLimiterMiddleware(chiRateLimiter)).Post("/login", authHandler.Login)
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
	config     config
	db         *pgxpool.Pool
	mongo      *mongo.Client
	rabbitConn *amqp.Connection
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn      string
	mongoURI string
}
