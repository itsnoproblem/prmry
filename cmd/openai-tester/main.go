package main

import (
	"bufio"
	"context"
	gosql "database/sql"
	"fmt"
	"github.com/itsnoproblem/mall-fountain-cop-bot/env"
	"github.com/joho/godotenv"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	gogpt "github.com/sashabaranov/go-gpt3"

	"github.com/itsnoproblem/mall-fountain-cop-bot/interacting"
	"github.com/itsnoproblem/mall-fountain-cop-bot/sql"
)

const (
	separator  = "-----------------------------------"
	prompt     = "> "
	rawKeyword = "RAW"
)

func banner() string {
	return `
Reddit GPT Bot
- - - - - - - ,
- - - Type something for the bot personality to respond to, or prefix with 
- - - ` + rawKeyword + ` to bypass the bot and send your input directly to the AI.
- - - - - - - - -
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
	defer func(db *gosql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("db.Close: %s", err)
		}
	}(db)

	interactionsRepo := sql.NewInteractionsRepo(db)
	moderationsRepo := sql.NewModerationsRepo(db)

	openAIKey := os.Getenv(env.VarOpenAIKey)
	gptClient := gogpt.NewClient(openAIKey)
	commenter := interacting.NewService(gptClient, &interactionsRepo, &moderationsRepo)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(banner())
	fmt.Print(prompt)
	for scanner.Scan() {
		response, err := commenter.RespondToMessage(context.Background(), scanner.Text())
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
