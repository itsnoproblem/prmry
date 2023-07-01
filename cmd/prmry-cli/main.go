package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/itsnoproblem/prmry/internal/envvars"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"

	"github.com/itsnoproblem/prmry/internal/interacting"
	"github.com/itsnoproblem/prmry/internal/sql"
)

const (
	separator  = "-----------------------------------"
	prompt     = "> "
	rawKeyword = "RAW"
)

func banner() string {
	return `
  _____   ______  ______,  
 |  ,  | |   ___||  .  / 
 |     \ |   |  ||  ,   \  
 |__|\__\|______||______/   
 = reddit gpt bot (v0.0.1)     
 
 Condition something for the rgb personality to respond to, or prefix with 
 "` + rawKeyword + `" to bypass the rgb and send your input directly to the AI.
`
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed loading .env: %s", err)
	}

	db := initDb()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("db.Close: %s", err)
		}
	}(db)

	interactionsRepo := sql.NewInteractionsRepo(db)
	moderationsRepo := sql.NewModerationsRepo(db)
	flowsRepo := sql.NewFlowsRepository(db)

	openAIKey := os.Getenv(envvars.OpenAIKey)
	gptClient := gogpt.NewClient(openAIKey)
	commenter := interacting.NewService(gptClient, &interactionsRepo, &moderationsRepo, flowsRepo)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(banner())
	fmt.Print(prompt)

	for scanner.Scan() {
		flowID := ""
		response, err := commenter.GenerateResponse(context.Background(), scanner.Text(), flowID)
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
		}

		fmt.Println(response)
		fmt.Println(separator)
		fmt.Print(prompt)
	}

	if scanner.Err() != nil {
		log.Printf("ERROR: %s\n", scanner.Err())
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
