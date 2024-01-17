package interacting

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/flow"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
	"github.com/itsnoproblem/prmry/internal/interaction"
	"github.com/itsnoproblem/prmry/internal/moderation"
)

type Service interface {
	Interactions(ctx context.Context) ([]interaction.Summary, error)
	Interaction(ctx context.Context, interactionID string) (interaction.Interaction, error)
	Moderation(ctx context.Context, interactionID string) (moderation.Moderation, error)
	ModerationByID(ctx context.Context, moderationID string) (moderation.Moderation, error)
	NewInteraction(ctx context.Context, msg, flowID string, params map[string]string) (interaction.Interaction, error)
	ExecuteFlow(ctx context.Context, inputText, flowID string, params map[string]string) (exec flow.Execution, err error)
}

type flowService interface {
	GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error)
}

type chatPrompt struct {
	SelectedFlow string
	InputParams  components.SortedMap
}

type chatPromptResponse struct {
	Flows          []flow.Flow
	SelectedFlow   string
	RequiredParams map[string]bool
}

func makeChatPromptEndpoint(svc flowService) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, err := getAuthorizedUser(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "interacting.makeChatPromptEndpoint")
		}

		req, ok := request.(chatPrompt)
		if !ok {
			return nil, fmt.Errorf("interacting.makeChatPromptEndpoint: failed to parse request")
		}

		flows, err := svc.GetFlowsForUser(ctx, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "interacting.makeChatPromptEndpoint")
		}

		res := chatPromptResponse{
			Flows:          flows,
			SelectedFlow:   req.SelectedFlow,
			RequiredParams: make(map[string]bool),
		}

		if req.SelectedFlow != "" {
			for _, flw := range flows {
				if flw.ID == req.SelectedFlow {
					for _, param := range flw.InputParams {
						res.RequiredParams[param.Key] = param.IsRequired
					}

					break
				}
			}
		}

		return res, nil
	}
}

func (req createInteractionRequest) validate() error {
	if req.Input == "" {
		return fmt.Errorf("validate: input was empty")
	}
	return nil
}

func makeCreateInteractionEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(Input)
		if !ok {
			return nil, fmt.Errorf("interacting.makeCreateInteractionEndpoint: failed to parse request")
		}

		ixn, err := svc.NewInteraction(ctx, req.InputMessage, req.FlowID, req.Params)
		if err != nil {
			return nil, errors.Wrap(err, "interacting.makeCreateInteractionEndpoint")
		}

		return ixn, nil
	}
}

func makeListInteractionsEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ixns, err := svc.Interactions(ctx)
		if err != nil {
			return nil, fmt.Errorf("interacting.makeListInteractionsEndpoint: %s", err)
		}

		return ixns, nil
	}
}

type getInteractionRequest struct {
	ID string
}

func makeGetInteractionEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(getInteractionRequest)
		if !ok {
			return nil, fmt.Errorf("interacting.makeGetInteractionEndpoint: failed to parse request")
		}

		ixn, err := svc.Interaction(ctx, req.ID)
		if err != nil {
			return nil, fmt.Errorf("interacting.makeListInteractionsEndpoint: %s", err)
		}

		return ixn, nil
	}
}

type executeFlowResponse struct {
	Model       string
	Temperature float64
	Prompt      string
	Executes    bool
}

func makeExecuteFlowEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(Input)
		if !ok {
			return nil, fmt.Errorf("interacting.makeExecuteFlowEndpoint: failed to parse request")
		}

		exec, err := svc.ExecuteFlow(ctx, req.InputMessage, req.FlowID, req.Params)
		if err != nil {
			return nil, fmt.Errorf("interacting.makeExecuteFlowEndpoint: %s", err)
		}

		return executeFlowResponse{
			Model:    exec.Model,
			Prompt:   exec.Prompt,
			Executes: exec.Executes,
		}, nil
	}
}

func getAuthorizedUser(ctx context.Context) (auth.User, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return auth.User{}, fmt.Errorf("user is missing")
	}

	return *user, nil
}
