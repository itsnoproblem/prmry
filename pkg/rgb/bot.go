package rgb

import (
	"context"
	"fmt"
	"github.com/itsnoproblem/prmry/pkg/log"
	"math/rand"
	"strings"
	"time"

	"github.com/turnage/graw/reddit"
)

const (
	minWaitTime                     = 10
	maxWaitTime                     = 300
	alwaysReplyToTopLevelCommentsBy = "HallOfTheMountainCop"
	neverRespondToCommentsBy        = "MallFountainCop"
)

type ResponseGenerator interface {
	GenerateResponse(ctx context.Context, msg string) (string, error)
}

type bot struct {
	logger      log.Logger
	responder   ResponseGenerator
	bot         reddit.Bot
	minWaitTime int
	maxWaitTime int
	subReddit   string
}

func NewBot(rb reddit.Bot, commenter ResponseGenerator, logger log.Logger, sub string) *bot {
	return &bot{
		bot:         rb,
		responder:   commenter,
		logger:      logger,
		minWaitTime: minWaitTime,
		maxWaitTime: maxWaitTime,
		subReddit:   sub,
	}
}

func (b *bot) Post(post *reddit.Post) error {
	waitTime := time.Duration(b.waitRandom()) * time.Second

	b.logger.Info("Post "+post.Permalink,
		"subreddit", post.Subreddit,
		"author", post.Author,
		"wait", waitTime,
	)

	<-time.After(waitTime)

	replyText, err := b.responder.GenerateResponse(context.Background(), post.Title+" - "+post.SelfText)
	if err != nil {
		return fmt.Errorf("bot.Post: %s", err.Error())
	}

	replyText = strings.Replace(replyText, "\n", "  \n", -1)

	return b.bot.Reply(post.Name, replyText)
}

func (b *bot) Comment(post *reddit.Comment) error {
	if !b.shouldRespondToComment(post) {
		b.logger.Warn("Comment: "+post.Permalink,
			"subreddit", post.Subreddit,
			"author", post.Author,
			"action", "ignore")
		return nil
	}

	waitTime := time.Duration(b.waitRandom()) * time.Second

	b.logger.Info("Comment: "+post.Permalink,
		"subreddit", post.Subreddit,
		"author", post.Author,
		"wait", waitTime,
	)

	<-time.After(waitTime)

	replyText, err := b.responder.GenerateResponse(context.Background(), post.Body)
	if err != nil {
		b.logger.Error("failed to create response", "error", err.Error())
	}

	b.logger.Info("Comment: "+post.Permalink,
		"action", "reply",
		"post", post.Name,
	)
	return b.bot.Reply(post.Name, replyText)
}

func (b *bot) CommentReply(reply *reddit.Message) error {
	b.logger.Info("CommentReply: "+reply.Name, "subreddit", reply.Subreddit)

	if reply.Subreddit != b.subReddit {
		b.logger.Warn("CommentReply",
			"subreddit", reply.Subreddit,
			"author", reply.Author,
			"action", "ignore")

		return nil
	}

	waitTime := time.Duration(b.waitRandom()) * time.Second

	b.logger.Info("CommentReply: "+reply.Subject,
		"subreddit", reply.Subreddit,
		"author", reply.Author,
		"wait", waitTime,
	)

	<-time.After(waitTime)

	replyBackText, err := b.responder.GenerateResponse(context.Background(), reply.BodyHTML)
	if err != nil {
		return fmt.Errorf("bot.CommentReply: %s", err.Error())
	}

	return b.bot.Reply(reply.Name, replyBackText)
}

func (b bot) Mention(mention *reddit.Message) error {
	waitTime := time.Duration(b.waitRandom()) * time.Second

	b.logger.Info("Mention: ("+mention.Name+") "+mention.Subject,
		"subreddit", mention.Subreddit,
		"author", mention.Author,
		"wait", waitTime)

	<-time.After(waitTime)

	replyBackText, err := b.responder.GenerateResponse(context.Background(), mention.BodyHTML)
	if err != nil {
		return fmt.Errorf("bot.Mention: %s", err.Error())
	}

	return b.bot.Reply(mention.Name, replyBackText)
}

func (b *bot) shouldRespondToComment(comment *reddit.Comment) bool {
	if comment.Author == neverRespondToCommentsBy {
		return false
	}

	return b.shouldRespondToCommentText(comment.Body) ||
		b.shouldRespondToCommentText(comment.Author) ||
		(comment.IsTopLevel() && comment.Author == alwaysReplyToTopLevelCommentsBy)
}

func (b *bot) shouldRespondToCommentText(text string) bool {
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

func (b *bot) waitRandom() int {
	rand.Seed(time.Now().Unix())
	waitTime := rand.Intn(b.maxWaitTime-b.minWaitTime+1) + b.minWaitTime
	return waitTime
}
