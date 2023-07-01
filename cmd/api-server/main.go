package main

import (
	"flag"
	"github.com/itsnoproblem/prmry/internal/envvars"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"github.com/itsnoproblem/prmry/internal/interacting"
	"github.com/itsnoproblem/prmry/internal/sql"
)

const (
	defaultListen = ":3332"
)

func main() {

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env: %s", err)
	}
	listen := os.Getenv(envvars.ListenPort)
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

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	interactionsRepo := sql.NewInteractionsRepo(db)
	moderationsRepo := sql.NewModerationsRepo(db)
	flowsRepo := sql.NewFlowsRepository(db)
	interactingService := interacting.NewService(nil, &interactionsRepo, &moderationsRepo, flowsRepo)

	r.Group(interacting.RouteHandler(interactingService))

	log.Println("RGB API Listening on " + listen)
	if err := http.ListenAndServe(listen, r); err != nil {
		panic(err)
	}
}

func initDb() *sqlx.DB {
	var (
		dbHost = os.Getenv(envvars.DbHost)
		dbUser = os.Getenv(envvars.DbUser)
		dbPass = os.Getenv(envvars.DbPass)
		dbName = os.Getenv(envvars.DbName)
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
