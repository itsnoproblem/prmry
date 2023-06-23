package main

import (
	"github.com/itsnoproblem/prmry/pkg/env"
	golog "log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-gpt3"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"

	"github.com/itsnoproblem/prmry/pkg/interacting"
	"github.com/itsnoproblem/prmry/pkg/log"
	"github.com/itsnoproblem/prmry/pkg/rgb"
	"github.com/itsnoproblem/prmry/pkg/sql"
)

const (
	botAgentConfig  = "rgb.agent"
	redditRateLimit = 0
	subReddit       = "hotspotbattlebots"
)

func main() {
	golog.SetOutput(os.Stdout)
	golog.SetFlags(golog.LstdFlags)

	err := godotenv.Load()
	if err != nil {
		golog.Fatalf("Failed loading .env: %s", err)
	}

	db := initDb()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			golog.Fatalf("db.Close: %s", err)
		}
	}(db)

	interactionsRepo := sql.NewInteractionsRepo(db)
	moderationsRepo := sql.NewModerationsRepo(db)

	openAIKey := os.Getenv(env.VarOpenAIKey)
	gptClient := gogpt.NewClient(openAIKey)
	commenter := interacting.NewService(gptClient, &interactionsRepo, &moderationsRepo)

	// logger
	var mainLog log.Logger
	mainLog = log.NewLogger()

	// reddit bot / event handler
	subredditsToListenOn := []string{subReddit}
	if bot, err := reddit.NewBotFromAgentFile(botAgentConfig, redditRateLimit); err != nil {
		mainLog.Error("failed to create bot handle: " + err.Error())
	} else {
		cfg := graw.Config{
			Subreddits:        subredditsToListenOn,
			SubredditComments: subredditsToListenOn,
			CommentReplies:    true,
			Logger:            golog.Default(),
		}

		handler := rgb.NewBot(bot, commenter, mainLog, subReddit)

		mainLog.Info("Starting patrol",
			"subreddit", subReddit,
			"agent", botAgentConfig)

		if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
			mainLog.Error("failed to start")
		} else {
			mainLog.Warn("graw run failed", "wait", wait())
		}
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
		golog.Fatal(err.Error())
	}

	if err = db.Ping(); err != nil {
		golog.Fatalf("DB Ping failed: %s", err)
	}
	golog.Println("DB Ping: OK")

	return db
}
