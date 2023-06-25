package main

import (
	"encoding/hex"
	"flag"
	"github.com/itsnoproblem/prmry/pkg/env"
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

	"github.com/itsnoproblem/prmry/pkg/auth"
	"github.com/itsnoproblem/prmry/pkg/authorizing"
	"github.com/itsnoproblem/prmry/pkg/components"
	"github.com/itsnoproblem/prmry/pkg/flowing"
	"github.com/itsnoproblem/prmry/pkg/htmx"
	"github.com/itsnoproblem/prmry/pkg/interacting"
	"github.com/itsnoproblem/prmry/pkg/profiling"
	"github.com/itsnoproblem/prmry/pkg/sql"
)

const (
	defaultListen = "9999"
)

func main() {
	flag.Parse()

	if fileExists(".env") {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Failed to load .env: %s", err)
		}
	}

	listen := os.Getenv(env.VarListenPort)
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
	r.Use(htmx.Middleware)
	r.Use(auth.Middleware(authSecret))

	githubCallback := os.Getenv(env.AppURL) + "/auth/github/callback"
	googleCallback := os.Getenv(env.AppURL) + "/auth/google/callback"

	googleClient := google.New(
		os.Getenv(env.GoogleClientID),
		os.Getenv(env.GoogleClientSecret),
		googleCallback,
		"email",
	)

	githubClient := github.New(
		os.Getenv(env.GithubClientID),
		os.Getenv(env.GithubClientSecret),
		githubCallback,
		"user:email",
	)

	goth.UseProviders(
		githubClient,
		googleClient,
	)

	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv(env.SessionSecret)))
	renderer := components.NewRenderer()

	usersRepo := sql.NewUsersRepo(db)
	ixnRepo := sql.NewInteractionsRepo(db)
	modRepo := sql.NewModerationsRepo(db)
	flowsRepo := sql.NewFlowsRepository(db)

	authService := authorizing.NewService(usersRepo)
	ixnService := interacting.NewService(gptClient, &ixnRepo, &modRepo, flowsRepo)
	flowService := flowing.NewService(flowsRepo)

	authResource, err := authorizing.NewResource(renderer, authSecret, authService)
	if err != nil {
		log.Fatal(err.Error())
	}

	homeResource := profiling.NewResource(renderer)
	ixnResource := interacting.NewResource(renderer, ixnService)
	flowResource := flowing.NewResource(renderer, flowService)

	r.Mount("/", homeResource.Routes())
	r.Mount("/auth", authResource.Routes())
	r.Mount("/interactions", ixnResource.Routes())
	r.Mount("/flows", flowResource.Routes())

	staticFS := http.FileServer(http.Dir("www/static"))
	wellknownFS := http.FileServer(http.Dir("www/.well-known"))
	r.Handle("/static/*", http.StripPrefix("/static/", staticFS))
	r.Handle("/.well-known/*", http.StripPrefix("/.well-known/", wellknownFS))

	log.Println("Listening on " + listen)
	if err := http.ListenAndServe(":"+listen, r); err != nil {
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
