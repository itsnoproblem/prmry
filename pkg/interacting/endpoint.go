package interacting

import (
	"context"
	"fmt"

	"github.com/itsnoproblem/prmry/pkg/api"
	"github.com/itsnoproblem/prmry/pkg/interaction"
)

type interactingService interface {
	Interactions(ctx context.Context) ([]interaction.Summary, error)
	Interaction(ctx context.Context, interactionID string) (interaction.Interaction, error)
	Moderation(ctx context.Context, interactionID string) (interaction.Moderation, error)
	ModerationByID(ctx context.Context, moderationID string) (interaction.Moderation, error)
	NewInteraction(ctx context.Context, msg string) (interaction.Interaction, error)
}

func makeListInteractionsEndpoint(svc interactingService) api.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		history, err := svc.Interactions(ctx)
		if err != nil {
			return nil, fmt.Errorf("interacting.makeListInteractionsEndpoint: %s", err)
		}

		return history, nil
	}
}

type getInteractionRequest struct {
	ID string
}

func makeGetInteractionEndpoint(svc interactingService) api.HandlerFunc {
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

func makeGetInteractionHTMLEndpoint(svc interactingService) api.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(getInteractionRequest)
		if !ok {
			return nil, fmt.Errorf("makeGetInteractionHTMLEndpoint: failed to parse request")
		}

		inter, err := svc.Interaction(ctx, req.ID)
		if err != nil {
			return nil, fmt.Errorf("makeGetInteractionHTMLEndpoint: %s", err)
		}

		return inter, nil
	}
}
