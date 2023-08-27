package interacting

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"

	"github.com/pkg/errors"
	openai "github.com/sashabaranov/go-openai"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/flow"
	"github.com/itsnoproblem/prmry/internal/interaction"
)

const (
	GPTModel              = openai.GPT3Dot5Turbo
	GPTMaxTokens          = 4000
	GPTCharactersPerToken = 4
)

type Responder interface {
	GenerateResponse(ctx context.Context, msg string) (string, error)
}

type InteractionRepo interface {
	Insert(ctx context.Context, in interaction.Interaction) (err error)
	Delete(ctx context.Context, id string) error
	GetInteractionSummaries(ctx context.Context, userID string) ([]interaction.Summary, error)
	GetInteraction(ctx context.Context, id string) (interaction.Interaction, error)
}

type ModerationRepo interface {
	Insert(ctx context.Context, mod interaction.Moderation) error
}

type FlowRepo interface {
	GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error)
	GetFlow(ctx context.Context, flowID string) (flow.Flow, error)
}

func NewService(c *openai.Client, r InteractionRepo, m ModerationRepo, f FlowRepo) service {
	log.Printf("PRMRY - model: [%s] - max tokens: [%d] - char per token: [%d]\n",
		GPTModel, GPTMaxTokens, GPTCharactersPerToken)

	return service{
		gptClient:   c,
		history:     r,
		moderations: m,
		flows:       f,
	}
}

type service struct {
	gptClient   *openai.Client
	history     InteractionRepo
	moderations ModerationRepo
	flows       FlowRepo
	input       string
}

func (s service) GenerateResponse(ctx context.Context, msg, flowID string) (string, error) {
	ix, err := s.NewInteraction(ctx, msg, flowID)
	if err != nil {
		return "", errors.Wrap(err, "inetracting.GenerateResponse")
	}

	if len(ix.Response.Choices) == 0 {
		return "", errors.Wrap(err, "interacting.GenerateResponse: no choices")
	}

	return ix.Response.Choices[0].Message.Content, nil
}

func (s service) Interactions(ctx context.Context) ([]interaction.Summary, error) {
	usr := auth.UserFromContext(ctx)
	history, err := s.history.GetInteractionSummaries(ctx, usr.ID)
	if err != nil {
		return nil, fmt.Errorf("interacting.Interactions: %s", err)
	}

	return history, nil
}

func (s service) Interaction(ctx context.Context, interactionID string) (interaction.Interaction, error) {
	in, err := s.history.GetInteraction(ctx, interactionID)
	if err != nil {
		return interaction.Interaction{}, fmt.Errorf("interacting.GetInteraction: %s", err)
	}

	return in, nil
}

func (s service) Moderation(ctx context.Context, interactionID string) (interaction.Moderation, error) {
	return interaction.Moderation{}, fmt.Errorf("Not implemented")
}

func (s service) ModerationByID(ctx context.Context, moderationID string) (interaction.Moderation, error) {
	return interaction.Moderation{}, fmt.Errorf("Not implemented")
}

func (s service) GetFlows(ctx context.Context) ([]flow.Flow, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("interacting.GetFlows: missing user")
	}

	flows, err := s.flows.GetFlowsForUser(ctx, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "interacting.GetFlows")
	}

	return flows, nil
}

func (s service) NewInteraction(ctx context.Context, msg, flowID string) (interaction.Interaction, error) {
	if auth.UserFromContext(ctx) == nil {
		return interaction.Interaction{}, errors.New("Unauthorized")
	}

	ixn, err := s.executeFlow(ctx, msg, flowID)
	if err != nil {
		return interaction.Interaction{}, errors.Wrap(err, "interacting.NewInteraction")
	}

	return ixn, nil
}

// ---- private -----

func (s service) executeFlow(ctx context.Context, inputText, flowID string) (interaction.Interaction, error) {
	if inputText == "" {
		return interaction.Interaction{}, fmt.Errorf("executeFlow: input text cannot be empty")
	}
	s.input = inputText
	prompt := inputText

	if flowID != "" {
		flw, err := s.flows.GetFlow(ctx, flowID)
		if err != nil {
			return interaction.Interaction{}, errors.Wrap(err, "executeFlow")
		}

		prompt, err = s.getPromptFromFlow(ctx, inputText, flw)
		if err != nil {
			return interaction.Interaction{}, errors.Wrap(err, "executeFlow")
		}

		if prompt == "" {
			return interaction.Interaction{}, fmt.Errorf("Flow '%s' does not execute.", flw.Name)
		}
	}

	promptTokens := s.tokenCount(prompt)
	maxTokens := GPTMaxTokens - promptTokens

	req := openai.ChatCompletionRequest{
		Model:     GPTModel,
		MaxTokens: maxTokens,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	resp, gptErr := s.gptClient.CreateChatCompletion(ctx, req)
	if gptErr != nil {
		return interaction.Interaction{}, gptErr
	}

	newInteraction := interaction.Interaction{
		ID:        uuid.New().String(),
		Request:   req,
		Response:  resp,
		CreatedAt: time.Now(),
		UserID:    auth.UserFromContext(ctx).ID,
		FlowID:    flowID,
	}

	err := s.history.Insert(ctx, newInteraction)
	if err != nil {
		return interaction.Interaction{}, errors.Wrap(err, "executeFlow")
	}

	go s.moderate(ctx, newInteraction.ID, inputText)

	ixn, err := s.history.GetInteraction(ctx, newInteraction.ID)

	return ixn, nil
}

func (s service) getPromptFromFlow(ctx context.Context, inputText string, flw flow.Flow) (string, error) {
	for _, cond := range flw.Rules {
		comparator := inputText
		if cond.Field.Source == flow.FieldSourceFlow {
			ixn, err := s.executeFlow(ctx, inputText, cond.Field.Value)
			if err != nil {
				return "", errors.Wrap(err, "getPromptFromFlow")
			}

			comparator = ixn.ResponseText()
		}

		conditionMatches, err := cond.Matches(comparator)
		if err != nil {
			return "", fmt.Errorf("getPromptFromFlow: %s", err)
		}

		if conditionMatches {
			prompt, err := s.compilePrompt(ctx, flw)
			if err != nil {
				return "", errors.Wrap(err, "getPromptFromFlow")
			}

			return prompt, nil
		}
	}

	return "", nil
}

func (s service) compilePrompt(ctx context.Context, flw flow.Flow) (string, error) {
	args := make([]interface{}, 0)
	for _, arg := range flw.PromptArgs {
		switch arg.Source {
		case flow.FieldSourceInput:
			args = append(args, s.input)
		case flow.FieldSourceFlow:
			ixn, err := s.executeFlow(ctx, s.input, arg.Value)
			if err == nil {
				return "", errors.Wrap(err, "compilePrompt")
			}

			args = append(args, ixn.ResponseText())
		}
	}

	return fmt.Sprintf(flw.Prompt, args...), nil
}

func (s service) moderate(ctx context.Context, interactionID, msg string) {
	modReq := openai.ModerationRequest{
		Input: msg,
		Model: openai.ModerationTextStable,
	}

	modRes, err := s.gptClient.Moderations(ctx, modReq)
	if err != nil {
		log.Printf("ERROR - interactingService.moderate: %s", err)
	}

	s.moderations.Insert(ctx, interaction.Moderation{
		ID:            modRes.ID,
		InteractionID: interactionID,
		Model:         modRes.Model,
		Results:       modRes.Results,
		CreatedAt:     time.Now(),
	})

}

func (s service) tokenCount(text string) int {
	return len(text) / GPTCharactersPerToken
}
