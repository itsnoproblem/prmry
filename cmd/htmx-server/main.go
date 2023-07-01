package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/itsnoproblem/prmry/internal/envvars"
	"github.com/itsnoproblem/prmry/internal/staticrendering"
	gogpt "github.com/sashabaranov/go-gpt3"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/authorizing"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/flowing"
	"github.com/itsnoproblem/prmry/internal/htmx"
	"github.com/itsnoproblem/prmry/internal/interacting"
	"github.com/itsnoproblem/prmry/internal/profiling"
	"github.com/itsnoproblem/prmry/internal/sql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

const (
	defaultListen = "9999"
)

type AppConfig struct {
	Env                string
	AppURL             string
	ListenPort         string
	OpenAPIKey         string
	DBHost             string
	DBUser             string
	DBPass             string
	DBName             string
	GithubClientID     string
	GithubClientSecret string
	GoogleClientID     string
	GoogleClientSecret string
	SessionSecret      string
}

func (cfg AppConfig) validate() error {
	requiredParams := map[string]string{
		cfg.AppURL:             envvars.AppURL,
		cfg.ListenPort:         envvars.ListenPort,
		cfg.OpenAPIKey:         envvars.OpenAIKey,
		cfg.DBHost:             envvars.DbHost,
		cfg.DBUser:             envvars.DbUser,
		cfg.DBName:             envvars.DbName,
		cfg.ListenPort:         envvars.ListenPort,
		cfg.GithubClientID:     envvars.GithubClientID,
		cfg.GithubClientSecret: envvars.GithubClientSecret,
		cfg.GoogleClientID:     envvars.GoogleClientID,
		cfg.GoogleClientSecret: envvars.GoogleClientSecret,
		cfg.SessionSecret:      envvars.SessionSecret,
	}

	for val, name := range requiredParams {
		if val == "" {
			return fmt.Errorf("env var not found: %s", name)
		}
	}

	return nil
}

func mustLoadAppConfig() AppConfig {
	if !fileExists(".env") {
		log.Fatalf("missing .env file!")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("loading .env: %s", err)
	}

	cfg := AppConfig{
		Env:                os.Getenv(envvars.Env),
		AppURL:             os.Getenv(envvars.AppURL),
		ListenPort:         os.Getenv(envvars.ListenPort),
		OpenAPIKey:         os.Getenv(envvars.OpenAIKey),
		DBHost:             os.Getenv(envvars.DbHost),
		DBUser:             os.Getenv(envvars.DbUser),
		DBPass:             os.Getenv(envvars.DbPass),
		DBName:             os.Getenv(envvars.DbName),
		GithubClientID:     os.Getenv(envvars.GithubClientID),
		GithubClientSecret: os.Getenv(envvars.GoogleClientSecret),
		GoogleClientID:     os.Getenv(envvars.GoogleClientID),
		GoogleClientSecret: os.Getenv(envvars.GoogleClientSecret),
		SessionSecret:      os.Getenv(envvars.SessionSecret),
	}

	if err := cfg.validate(); err != nil {
		log.Fatalf("app config: %s", err.Error())
	}

	return cfg
}

func main() {
	flag.Parse()
	appConfig := mustLoadAppConfig()

	db := mustInitDB(appConfig.DBHost, appConfig.DBUser, appConfig.DBPass, appConfig.DBName)
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("db.Close: %s", err)
		}
	}(db)

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

	fixHostAndProto := appConfig.Env != "DEV"
	r.Use(auth.Middleware(authSecret, fixHostAndProto))

	googleClient := google.New(
		appConfig.GoogleClientID,
		appConfig.GoogleClientSecret,
		appConfig.AppURL+"/auth/google/callback",
		"email",
	)

	githubClient := github.New(
		appConfig.GithubClientID,
		appConfig.GithubClientSecret,
		appConfig.AppURL+"/auth/github/callback",
		"user:email",
	)

	gptClient := gogpt.NewClient(appConfig.OpenAPIKey)

	goth.UseProviders(
		githubClient,
		googleClient,
	)

	gothic.Store = sessions.NewCookieStore([]byte(appConfig.SessionSecret))
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

	staticResource := staticrendering.NewResource(renderer)

	r.Mount("/", homeResource.Routes())
	r.Mount("/auth", authResource.Routes())
	r.Mount("/interactions", ixnResource.Routes())
	r.Mount("/flows", flowResource.Routes())

	r.Get("/terms", staticResource.Terms)
	r.Get("/privacy", staticResource.Privacy)

	staticFS := http.FileServer(http.Dir("www/static"))
	wellknownFS := http.FileServer(http.Dir("www/.well-known"))
	r.Handle("/static/*", http.StripPrefix("/static/", staticFS))
	r.Handle("/.well-known/*", http.StripPrefix("/.well-known/", wellknownFS))

	log.Println("Listening on " + appConfig.ListenPort)
	if err := http.ListenAndServe(":"+appConfig.ListenPort, r); err != nil {
		panic(err)
	}
}

func mustInitDB(dbHost, dbUser, dbPass, dbName string) *sqlx.DB {
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
