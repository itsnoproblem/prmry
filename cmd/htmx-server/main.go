package main

import (
	gosql "database/sql"
	"encoding/hex"
	"flag"
	"github.com/gorilla/sessions"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/auth"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/authorizing"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/interacting"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/profiling"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/sql"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/templates"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	gogpt "github.com/sashabaranov/go-gpt3"

	"github.com/itsnoproblem/mall-fountain-cop-bot/env"
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
	defer func(db *gosql.DB) {
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
	)

	githubClient := github.New(
		os.Getenv(env.GithubClientID),
		os.Getenv(env.GithubClientSecret),
		os.Getenv(env.GithubCallbackURL),
	)

	goth.UseProviders(
		githubClient,
		googleClient,
	)

	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv(env.SessionSecret)))

	usersRepo := sql.NewUsersRepo(db)
	authService := authorizing.NewService(usersRepo)
	authResource, err := authorizing.NewResource(authSecret, authService)
	if err != nil {
		log.Fatal(err.Error())
	}

	tpl, err := templates.Parse()
	if err != nil {
		log.Fatalf("Failed to parse templates: %s", err)
	}

	profileResource, err := profiling.NewResource(tpl)
	if err != nil {
		log.Fatal(err.Error())
	}

	ixnRepo := sql.NewInteractionsRepo(db)
	modRepo := sql.NewModerationsRepo(db)
	ixnService := interacting.NewService(gptClient, &ixnRepo, &modRepo)
	ixnResource := interacting.NewResource(tpl, ixnService)

	r.Mount("/", profileResource.Routes())
	r.Mount("/auth", authResource.Routes())
	r.Mount("/interactions", ixnResource.Routes())

	log.Println("Listening on " + listen)
	if err := http.ListenAndServe(listen, r); err != nil {
		panic(err)
	}
}

func initDb() *gosql.DB {
	var (
		dbHost = os.Getenv(env.VarDBHost)
		dbUser = os.Getenv(env.VarDBUser)
		dbPass = os.Getenv(env.VarDBPass)
		dbName = os.Getenv(env.VarDBName)
	)

	db, err := gosql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName+"?parseTime=true")
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("DB Ping failed: %s", err)
	}
	log.Println("DB Ping: OK")

	return db
}
