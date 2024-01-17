package interacting

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	openai "github.com/sashabaranov/go-openai"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/flow"
	"github.com/itsnoproblem/prmry/internal/interaction"
	"github.com/itsnoproblem/prmry/internal/moderation"
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
	Insert(ctx context.Context, mod moderation.Moderation) error
}

type FlowService interface {
	GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error)
	GetFlow(ctx context.Context, flowID string) (flow.Flow, error)
}

func NewService(c *openai.Client, r InteractionRepo, m ModerationRepo, f FlowService) *service {
	log.Printf("PRMRY - max tokens: [%d] - char per token: [%d]\n",
		GPTMaxTokens, GPTCharactersPerToken)

	return &service{
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
	flows       FlowService
	input       string
}

func (s service) GenerateResponse(ctx context.Context, msg, flowID string, params map[string]string) (string, error) {
	ix, err := s.NewInteraction(ctx, msg, flowID, params)
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

func (s service) Moderation(ctx context.Context, interactionID string) (moderation.Moderation, error) {
	return moderation.Moderation{}, fmt.Errorf("Not implemented")
}

func (s service) ModerationByID(ctx context.Context, moderationID string) (moderation.Moderation, error) {
	return moderation.Moderation{}, fmt.Errorf("Not implemented")
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

func (s service) NewInteraction(ctx context.Context, msg, flowID string, params map[string]string) (interaction.Interaction, error) {
	if auth.UserFromContext(ctx) == nil {
		return interaction.Interaction{}, errors.New("Unauthorized")
	}

	ixn, executes, err := s.executeFlowAndInteract(ctx, msg, flowID, params)
	if err != nil {
		return interaction.Interaction{}, errors.Wrap(err, "interacting.NewInteraction")
	}

	if !executes {
		return interaction.Interaction{}, fmt.Errorf("interacting.NewInteraction: flow does not execute")
	}

	return ixn, nil
}

func (s service) ExecuteFlow(ctx context.Context, inputText, flowID string, params map[string]string) (exec flow.Execution, err error) {
	flw, err := s.flows.GetFlow(ctx, flowID)
	if err != nil {
		return flow.Execution{}, errors.Wrap(err, "flowing.ExecuteFlow")
	}

	s.input = inputText
	prompt, err := s.getPromptFromFlow(ctx, inputText, flw, params)
	if err != nil {
		return flow.Execution{}, errors.Wrap(err, "interacting.ExecuteFlow")
	}

	if prompt == "" {
		return flow.Execution{}, nil
	}

	return flow.Execution{
		Model:       flw.Model,
		Temperature: flw.Temperature,
		Executes:    true,
		Prompt:      prompt,
	}, nil
}

// ---- private -----

func (s service) executeFlowAndInteract(ctx context.Context, inputText, flowID string, params map[string]string) (ixn interaction.Interaction, executes bool, err error) {
	if inputText == "" {
		return interaction.Interaction{}, false, fmt.Errorf("executeFlowAndInteract: input text cannot be empty")
	}
	s.input = inputText
	prompt := inputText

	if flowID != "" {
		exec, err := s.ExecuteFlow(ctx, inputText, flowID, params)
		if err != nil {
			return interaction.Interaction{}, false, errors.Wrap(err, "executeFlowAndInteract")
		}

		if !exec.Executes {
			return interaction.Interaction{}, false, nil
		}

		prompt = exec.Prompt
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
		return interaction.Interaction{}, true, gptErr
	}

	var completion string
	if len(resp.Choices) > 0 {
		completion = resp.Choices[0].Message.Content
	}

	ixn = interaction.Interaction{
		ID:               uuid.New().String(),
		Request:          req,
		Response:         resp,
		CreatedAt:        time.Now(),
		UserID:           auth.UserFromContext(ctx).ID,
		FlowID:           flowID,
		Type:             resp.Object,
		Model:            resp.Model,
		Prompt:           prompt,
		Completion:       completion,
		TokensPrompt:     resp.Usage.PromptTokens,
		TokensCompletion: resp.Usage.CompletionTokens,
	}

	err = s.history.Insert(ctx, ixn)
	if err != nil {
		return interaction.Interaction{}, true, errors.Wrap(err, "executeFlowAndInteract")
	}

	go s.moderate(ctx, ixn.ID, inputText)

	ixn, err = s.history.GetInteraction(ctx, ixn.ID)
	if err != nil {
		ixn.FlowName = "Error fetching flow"
	}

	return ixn, true, nil
}

func (s service) getPromptFromFlow(ctx context.Context, inputText string, flw flow.Flow, params map[string]string) (string, error) {
	for _, cond := range flw.Triggers {
		// comparator is the value we're comparing against
		var comparator string

		switch cond.Field.Source {
		case flow.FieldSourceInputArg:
			var ok bool
			comparator, ok = params[cond.Field.Value]
			if !ok {
				return "", fmt.Errorf("getPromptFromFlow: missing input arg '%s'", cond.Field.Value)
			}
		case flow.FieldSourceModeration:
			return "", fmt.Errorf("getPromptFromFlow: moderation not implemented")
		case flow.FieldSourceInput:
			comparator = inputText
		case flow.FieldSourceFlow:
			if flw.ID == cond.Field.Value {
				return "", fmt.Errorf("getPromptFromFlow: flow '%s' cannot reference itself", flw.Name)
			}
			ixn, executes, err := s.executeFlowAndInteract(ctx, inputText, cond.Field.Value, params)
			if err != nil {
				return "", errors.Wrap(err, "getPromptFromFlow")
			}

			if executes {
				comparator = ixn.ResponseText()
			}
		}

		conditionMatches, err := cond.Matches(comparator)
		if err != nil {
			return "", fmt.Errorf("getPromptFromFlow: %s", err)
		}

		if conditionMatches {
			prompt, err := s.compilePrompt(ctx, flw, params)
			if err != nil {
				return "", errors.Wrap(err, "getPromptFromFlow")
			}

			return prompt, nil
		}
	}

	return "", nil
}

func (s service) compilePrompt(ctx context.Context, flw flow.Flow, params map[string]string) (string, error) {
	args := make([]interface{}, 0)

	for _, arg := range flw.PromptArgs {
		switch arg.Source {
		case flow.FieldSourceInput:
			args = append(args, s.input)

		case flow.FieldSourceInputArg:
			args = append(args, params[arg.Value])

		case flow.FieldSourceFlow:
			ixn, executes, err := s.executeFlowAndInteract(ctx, s.input, arg.Value, params)
			if err == nil {
				return "", errors.Wrap(err, "compilePrompt")
			}

			var responseText string
			if executes {
				responseText = ixn.ResponseText()
			}

			args = append(args, responseText)
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
		log.Printf("ERROR - Service.moderate: %s", err)
	}

	s.moderations.Insert(ctx, moderation.Moderation{
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
