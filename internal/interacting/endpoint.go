package interacting

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/flow"
	"github.com/itsnoproblem/prmry/internal/htmx"
	"github.com/itsnoproblem/prmry/internal/interaction"
	"github.com/itsnoproblem/prmry/internal/moderation"
)

type interactingService interface {
	Interactions(ctx context.Context) ([]interaction.Summary, error)
	Interaction(ctx context.Context, interactionID string) (interaction.Interaction, error)
	Moderation(ctx context.Context, interactionID string) (moderation.Moderation, error)
	ModerationByID(ctx context.Context, moderationID string) (moderation.Moderation, error)
	NewInteraction(ctx context.Context, msg, flowID string, params map[string]string) (interaction.Interaction, error)
}

type flowService interface {
	GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error)
}

type chatPromptRequest struct {
	SelectedFlow string
	FlowParams   map[string]string
}

type chatPromptResponse struct {
	Flows          []flow.Flow
	SelectedFlow   string
	RequiredParams map[string]bool
}

func makeChatPromptEndpoint(svc flowService) htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, err := getAuthorizedUser(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "makeChatPromptEndpoint")
		}

		req, ok := request.(chatPromptRequest)
		if !ok {
			return nil, fmt.Errorf("makeChatPromptEndpoint: failed to parse request")
		}

		flows, err := svc.GetFlowsForUser(ctx, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "makeChatPromptEndpoint")
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

type createInteractionRequest struct {
	FlowID string
	Input  string
	Params map[string]string
}

func (req createInteractionRequest) validate() error {
	if req.Input == "" {
		return fmt.Errorf("input was empty")
	}
	return nil
}

func makeCreateInteractionEndpoint(svc interactingService) htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(createInteractionRequest)
		if !ok {
			return nil, fmt.Errorf("makeCreateInteractionEndpoint: failed to parse request")
		}

		if err = req.validate(); err != nil {
			return nil, errors.Wrap(err, "makeCreateInteractionEndpoint")
		}

		ixn, err := svc.NewInteraction(ctx, req.Input, req.FlowID, req.Params)
		if err != nil {
			return nil, errors.Wrap(err, "makeCreateInteractionEndpoint")
		}

		return ixn, nil
	}
}

func makeListInteractionsEndpoint(svc interactingService) htmx.HandlerFunc {
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

func makeGetInteractionEndpoint(svc interactingService) htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(getInteractionRequest)
		if !ok {
			return nil, fmt.Errorf("makeGetInteractionEndpoint: failed to parse request")
		}

		interaction, err := svc.Interaction(ctx, req.ID)
		if err != nil {
			return nil, fmt.Errorf("interacting.makeListInteractionsEndpoint: %s", err)
		}

		return interaction, nil
	}
}

func getAuthorizedUser(ctx context.Context) (auth.User, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return auth.User{}, fmt.Errorf("user is missing")
	}

	return *user, nil
}
