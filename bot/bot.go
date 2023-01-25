package bot

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/turnage/graw/reddit"

	"github.com/itsnoproblem/mall-fountain-cop-bot/log"
)

const (
	minWaitTime                     = 10
	maxWaitTime                     = 300
	alwaysReplyToTopLevelCommentsBy = "HallOfTheMountainCop"
	neverRespondToCommentsBy        = "MallFountainCop"
)

type CommentResponder interface {
	RespondToMessage(ctx context.Context, msg string) (string, error)
}

type rgb struct {
	bot             reddit.Bot
	responseCreator CommentResponder
	logger          log.Logger
	minWaitTime     int
	maxWaitTime     int
	subReddit       string
}

func NewBot(bot reddit.Bot, commenter CommentResponder, logger log.Logger, sub string) *rgb {
	return &rgb{
		bot:             bot,
		responseCreator: commenter,
		logger:          logger,
		minWaitTime:     minWaitTime,
		maxWaitTime:     maxWaitTime,
		subReddit:       sub,
	}
}

func (bot *rgb) Post(post *reddit.Post) error {
	waitTime := time.Duration(bot.waitRandom()) * time.Second

	bot.logger.Info("Post "+post.Permalink,
		"subreddit", post.Subreddit,
		"author", post.Author,
		"wait", waitTime,
	)

	<-time.After(waitTime)

	replyText, err := bot.responseCreator.RespondToMessage(context.Background(), post.Title+" - "+post.SelfText)
	if err != nil {
		return fmt.Errorf("rgb.Post: %s", err.Error())
	}

	replyText = strings.Replace(replyText, "\n", "  \n", -1)

	return bot.bot.Reply(post.Name, replyText)
}

func (bot *rgb) Comment(post *reddit.Comment) error {
	if !bot.shouldRespondToComment(post) {
		bot.logger.Warn("Comment: "+post.Permalink,
			"subreddit", post.Subreddit,
			"author", post.Author,
			"action", "ignore")
		return nil
	}

	waitTime := time.Duration(bot.waitRandom()) * time.Second

	bot.logger.Info("Comment: "+post.Permalink,
		"subreddit", post.Subreddit,
		"author", post.Author,
		"wait", waitTime,
	)

	<-time.After(waitTime)

	replyText, err := bot.responseCreator.RespondToMessage(context.Background(), post.Body)
	if err != nil {
		bot.logger.Error("failed to create response", "error", err.Error())
	}

	return bot.bot.Reply(post.Name, replyText)
}

func (bot *rgb) CommentReply(reply *reddit.Message) error {
	bot.logger.Info("CommentReply: "+reply.Name, "subreddit", reply.Subreddit)

	if reply.Subreddit != bot.subReddit {
		bot.logger.Warn("CommentReply",
			"subreddit", reply.Subreddit,
			"author", reply.Author,
			"action", "ignore")

		return nil
	}

	waitTime := time.Duration(bot.waitRandom()) * time.Second

	bot.logger.Info("CommentReply: "+reply.Subject,
		"subreddit", reply.Subreddit,
		"author", reply.Author,
		"wait", waitTime,
	)

	<-time.After(waitTime)

	replyBackText, err := bot.responseCreator.RespondToMessage(context.Background(), reply.BodyHTML)
	if err != nil {
		return fmt.Errorf("rgb.CommentReply: %s", err.Error())
	}

	return bot.bot.Reply(reply.Name, replyBackText)
}

func (bot rgb) Mention(mention *reddit.Message) error {
	waitTime := time.Duration(bot.waitRandom()) * time.Second

	bot.logger.Info("Mention: ("+mention.Name+") "+mention.Subject,
		"subreddit", mention.Subreddit,
		"author", mention.Author,
		"wait", waitTime)

	<-time.After(waitTime)

	replyBackText, err := bot.responseCreator.RespondToMessage(context.Background(), mention.BodyHTML)
	if err != nil {
		return fmt.Errorf("rgb.Mention: %s", err.Error())
	}

	return bot.bot.Reply(mention.Name, replyBackText)
}

func (bot *rgb) shouldRespondToComment(comment *reddit.Comment) bool {
	if comment.Author == neverRespondToCommentsBy {
		return false
	}

	return bot.shouldRespondToCommentText(comment.Body) ||
		bot.shouldRespondToCommentText(comment.Author) ||
		(comment.IsTopLevel() && comment.Author == alwaysReplyToTopLevelCommentsBy)
}

func (bot *rgb) shouldRespondToCommentText(text string) bool {
	triggers := []string{
		"cop",
		"police",
		"arrest",
		"justice",
		"crime",
		"drugs",
		"meth",
		"crack",
		"edibles",
		"weed",
		"911",
		"karen",
		"bigot",
		"racist",
		"war",
		"gun",
		"tony",
		"christmas",
		"xmas",
	}

	text = strings.ToLower(text)
	for _, trigger := range triggers {
		if strings.Contains(text, trigger) {
			return true
		}
	}

	return false
}

func (bot *rgb) waitRandom() int {
	rand.Seed(time.Now().Unix())
	waitTime := rand.Intn(bot.maxWaitTime-bot.minWaitTime+1) + bot.minWaitTime
	return waitTime
}
