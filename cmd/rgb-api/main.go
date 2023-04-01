package main

import (
	gosql "database/sql"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/itsnoproblem/mall-fountain-cop-bot/env"
	interacting2 "github.com/itsnoproblem/mall-fountain-cop-bot/pkg/interacting"
	sql2 "github.com/itsnoproblem/mall-fountain-cop-bot/pkg/sql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

const (
	defaultListen = ":3333"
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

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	interactionsRepo := sql2.NewInteractionsRepo(db)
	moderationsRepo := sql2.NewModerationsRepo(db)

	interactingService := interacting2.NewService(nil, &interactionsRepo, &moderationsRepo)
	interactionRenderer, err := interacting2.NewRenderer()
	if err != nil {
		log.Fatalf("Failed to load interaction renderer: %s", err)
	}

	r.Group(interacting2.RouteHandler(interactingService, interactionRenderer))

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
