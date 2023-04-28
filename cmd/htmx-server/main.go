package main

import (
	"encoding/hex"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	gogpt "github.com/sashabaranov/go-gpt3"

	"github.com/itsnoproblem/mall-fountain-cop-bot/env"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/auth"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/authorizing"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/htmx"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/interacting"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/profiling"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/sql"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/templates"
)

const (
	defaultListen = "9999"
)

func main() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env: %s", err)
	}

	listen := os.Getenv(env.VarListenAddress)
	if listen == "" {
		listen = defaultListen
	}

	db := initDb()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("db.Close: %s", err)
		}
	}(db)

	openAIKey := os.Getenv(env.VarOpenAIKey)
	gptClient := gogpt.NewClient(openAIKey)

	// Authorizer Middleware
	authSecret, err := hex.DecodeString("1ad6bbbff4d1c1e08608f814570d562c9b5ef2fc4c9e6b5ec4c9f3234b595bcb")
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(render.SetContentType(render.ContentTypeHTML))
	r.Use(auth.Middleware(authSecret))

	googleClient := google.New(
		os.Getenv(env.GoogleClientID),
		os.Getenv(env.GoogleClientSecret),
		os.Getenv(env.GoogleCallbackURL),
		"email",
	)

	githubClient := github.New(
		os.Getenv(env.GithubClientID),
		os.Getenv(env.GithubClientSecret),
		os.Getenv(env.GithubCallbackURL),
		"user:email",
	)

	goth.UseProviders(
		githubClient,
		googleClient,
	)

	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv(env.SessionSecret)))

	tpl, err := templates.Parse()
	if err != nil {
		log.Fatalf("Failed to parse templates: %s", err)
	}
	renderer := htmx.NewRenderer(tpl)

	usersRepo := sql.NewUsersRepo(db)
	authService := authorizing.NewService(usersRepo)
	authResource, err := authorizing.NewResource(renderer, authSecret, authService)
	if err != nil {
		log.Fatal(err.Error())
	}

	profileResource, err := profiling.NewResource(tpl)
	if err != nil {
		log.Fatal(err.Error())
	}

	ixnRepo := sql.NewInteractionsRepo(db)
	modRepo := sql.NewModerationsRepo(db)
	ixnService := interacting.NewService(gptClient, &ixnRepo, &modRepo)
	ixnResource := interacting.NewResource(renderer, ixnService)

	r.Mount("/", profileResource.Routes())
	r.Mount("/auth", authResource.Routes())
	r.Mount("/interactions", ixnResource.Routes())

	log.Println("Listening on " + listen)
	if err := http.ListenAndServe(listen, r); err != nil {
		panic(err)
	}
}

func initDb() *sqlx.DB {
	var (
		dbHost = os.Getenv(env.VarDBHost)
		dbUser = os.Getenv(env.VarDBUser)
		dbPass = os.Getenv(env.VarDBPass)
		dbName = os.Getenv(env.VarDBName)
	)

	db, err := sqlx.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName+"?parseTime=true")
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("DB Ping failed: %s", err)
	}
	log.Println("DB Ping: OK")

	return db
}
