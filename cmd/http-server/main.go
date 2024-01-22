package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/itsnoproblem/prmry/internal/funneling"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	openai "github.com/sashabaranov/go-openai"

	"github.com/itsnoproblem/prmry/internal/accounting"
	"github.com/itsnoproblem/prmry/internal/api"
	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/authenticating"
	"github.com/itsnoproblem/prmry/internal/envvars"
	"github.com/itsnoproblem/prmry/internal/flowing"
	"github.com/itsnoproblem/prmry/internal/htmx"
	"github.com/itsnoproblem/prmry/internal/interacting"
	"github.com/itsnoproblem/prmry/internal/prmrying"
	"github.com/itsnoproblem/prmry/internal/sql"
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
	AuthSecret         auth.Byte32
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

	if fileExists(".env") {
		log.Println("Loading .env file...")
		if err := godotenv.Load(); err != nil {
			log.Fatalf("loading .env: %s", err)
		}
	}

	authSecret, err := hex.DecodeString(os.Getenv(envvars.AuthSecret))
	if err != nil {
		log.Fatalf("ERROR: %s (must be 32 byte hex string): %s", envvars.AuthSecret, err)
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
		GithubClientSecret: os.Getenv(envvars.GithubClientSecret),
		GoogleClientID:     os.Getenv(envvars.GoogleClientID),
		GoogleClientSecret: os.Getenv(envvars.GoogleClientSecret),
		SessionSecret:      os.Getenv(envvars.SessionSecret),
		AuthSecret:         authSecret,
	}

	//if debug, err := json.Marshal(cfg); err == nil {
	//	log.Printf("%s", string(debug))
	//}

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

	renderer := htmx.NewRenderer()
	apiRenderer := api.NewRenderer()
	gptClient := openai.NewClient(appConfig.OpenAPIKey)

	usersRepo := sql.NewUsersRepo(db)
	interactionsRepo := sql.NewInteractionsRepo(db)
	moderationsRepo := sql.NewModerationsRepo(db)
	flowsRepo := sql.NewFlowsRepository(db)
	funnelsRepo := sql.NewFunnelsRepository(db)

	accountingService := accounting.NewService(&usersRepo)
	authService := authenticating.NewService(&usersRepo)
	interactingService := interacting.NewService(gptClient, &interactionsRepo, &moderationsRepo, flowsRepo)
	flowingService := flowing.NewService(flowsRepo)
	funnelingService := funneling.NewService(funnelsRepo)

	// TODO: replace this shortcut OAuth implementation that uses goth / gothic
	authResource, err := authenticating.NewResource(renderer, appConfig.AuthSecret, authService)
	if err != nil {
		log.Fatal(err.Error())
	}

	// UI Transports
	interactingTransport := interacting.HTMXRouteHandler(interactingService, flowingService, renderer)
	flowingTransport := flowing.RouteHandler(flowingService, renderer)
	staticTransport := prmrying.RouteHandler(renderer)
	accountingTransport := accounting.RouteHandler(accountingService, renderer)
	funnelingTransport := funneling.RouteHandler(funnelingService, renderer)
	authenticatingTransport := authResource.RouteHandler()

	// API Transports
	interactingAPITransport := interacting.JSONRouteHandler(interactingService, apiRenderer)
	flowingAPITransport := flowing.JSONRouteHandler(flowingService, apiRenderer)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(render.SetContentType(render.ContentTypeHTML))
	r.Use(middleware.Compress(5))
	r.Use(htmx.Middleware)
	fixHostAndProto := appConfig.Env != "DEV"
	r.Use(auth.Middleware(appConfig.AuthSecret, fixHostAndProto))
	r.Use(api.Middleware(&usersRepo, apiRenderer))

	goth.UseProviders(oauthProviders(appConfig)...)
	gothic.Store = sessions.NewCookieStore(appConfig.AuthSecret)
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return chi.URLParam(req, "provider"), nil
	}

	// UI Routes
	r.Group(staticTransport)
	r.Group(authenticatingTransport)
	r.Group(accountingTransport)
	r.Group(interactingTransport)
	r.Group(flowingTransport)
	r.Group(funnelingTransport)

	// API Routes
	r.Group(interactingAPITransport)
	r.Group(flowingAPITransport)

	staticFS := http.FileServer(http.Dir("www/static"))
	wellknownFS := http.FileServer(http.Dir("www/.well-known"))
	r.Handle("/static/*", http.StripPrefix("/static/", staticFS))
	r.Handle("/.well-known/*", http.StripPrefix("/.well-known/", wellknownFS))

	port := ":" + strings.TrimPrefix(appConfig.ListenPort, ":")
	log.Println("Listening on " + port)
	if err := http.ListenAndServe(port, r); err != nil {
		panic(err)
	}
}

func oauthProviders(appConfig AppConfig) []goth.Provider {
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

	return []goth.Provider{
		githubClient,
		googleClient,
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
