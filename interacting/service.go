package interacting

import (
	"context"
	"fmt"
	"github.com/itsnoproblem/mall-fountain-cop-bot/interaction"
	"log"
	"math/rand"
	"time"

	gogpt "github.com/sashabaranov/go-gpt3"
)

const (
	RawPromptTrigger      = "RAW"
	GPTModel              = gogpt.GPT3TextDavinci003
	GPTMaxTokens          = 4000
	GPTCharactersPerToken = 4
	maxInputChars         = 3000 * GPTCharactersPerToken
)

type Responder interface {
	RespondToMessage(ctx context.Context, msg string) (string, error)
}

type InteractionRepo interface {
	Add(ctx context.Context, in interaction.Interaction) (id string, err error)
	Remove(ctx context.Context, id string) error
	Summaries(ctx context.Context) ([]interaction.Summary, error)
	Interaction(ctx context.Context, id string) (interaction.Interaction, error)
}

type ModerationRepo interface {
	Add(ctx context.Context, mod interaction.Moderation) error
	Remove(ctx context.Context, id string) error
	All(ctx context.Context) ([]interaction.Moderation, error)
}

func NewService(c *gogpt.Client, r InteractionRepo, m ModerationRepo) service {
	log.Printf("RGB - model: [%s] - max tokens: [%d] - char per token: [%d]\n",
		GPTModel, GPTMaxTokens, GPTCharactersPerToken)

	return service{
		gptClient:   c,
		history:     r,
		moderations: m,
	}
}

type service struct {
	gptClient   *gogpt.Client
	history     InteractionRepo
	moderations ModerationRepo
}

func (s service) RespondToMessage(ctx context.Context, msg string) (string, error) {
	prompt := s.prompt(msg)
	promptTokens := s.tokenCount(prompt)
	maxTokens := GPTMaxTokens - promptTokens

	req := gogpt.CompletionRequest{
		Model:     GPTModel,
		MaxTokens: maxTokens,
		Prompt:    prompt,
	}

	resp, gptErr := s.gptClient.CreateCompletion(ctx, req)
	err := ""
	if gptErr != nil {
		err = gptErr.Error()
	}

	interactionID, histErr := s.history.Add(ctx, interaction.Interaction{
		Request:   req,
		Response:  resp,
		Error:     err,
		CreatedAt: time.Now(),
	})
	if histErr != nil {
		log.Printf("ERROR - Failed to save interaction history: %s", histErr)
	}

	if gptErr != nil {
		return "", gptErr
	}

	go s.moderate(ctx, interactionID, msg)

	return resp.Choices[0].Text, nil
}

func (s service) Interactions(ctx context.Context) ([]interaction.Summary, error) {
	history, err := s.history.Summaries(ctx)
	if err != nil {
		return nil, fmt.Errorf("interacting.Interactions: %s", err)
	}

	return history, nil
}

func (s service) Interaction(ctx context.Context, interactionID string) (interaction.Interaction, error) {
	in, err := s.history.Interaction(ctx, interactionID)
	if err != nil {
		return interaction.Interaction{}, fmt.Errorf("interacting.Interaction: %s", err)
	}

	return in, nil
}

func (s service) Moderation(ctx context.Context, interactionID string) (interaction.Moderation, error) {
	return interaction.Moderation{}, fmt.Errorf("Not implemented")
}

func (s service) ModerationByID(ctx context.Context, moderationID string) (interaction.Moderation, error) {
	return interaction.Moderation{}, fmt.Errorf("Not implemented")
}

func (s service) moderate(ctx context.Context, interactionID, msg string) {
	modReq := gogpt.ModerationRequest{
		Input: msg,
		Model: nil,
	}

	modRes, err := s.gptClient.Moderations(ctx, modReq)
	if err != nil {
		log.Printf("ERROR - interactingService.moderate: %s", err)
	}

	s.moderations.Add(ctx, interaction.Moderation{
		ID:            modRes.ID,
		InteractionID: interactionID,
		Model:         modRes.Model,
		Results:       modRes.Results,
		CreatedAt:     time.Now(),
	})

}

// private

func (s service) prompt(msg string) string {
	var prompt string
	if len(msg) > maxInputChars {
		msg = msg[:maxInputChars]
	}

	if len(msg) >= len(RawPromptTrigger) && msg[:len(RawPromptTrigger)] == RawPromptTrigger {
		prompt = msg[len(RawPromptTrigger):]
	} else {
		prompt = s.generatePrompt(msg)
	}

	return prompt
}

func (s service) generatePrompt(msg string) string {
	forms := []string{
		"in the form of a short poem",
		"as a sarcastic cop who speaks in riddles",
		"in the form of a limmerick",
		"as a dramatic, macho cop, who comes to surprising conclusions",
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(forms), func(i, j int) { forms[i], forms[j] = forms[j], forms[i] })

	prompt :=
		`respond to the text below ` + forms[0] + `.

Text: """
` + msg + `
"""
	`

	return prompt
}

func (s service) tokenCount(text string) int {
	return len(text) / GPTCharactersPerToken
}
