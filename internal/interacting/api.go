package interacting

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/interaction"
)

type createInteractionAPIRequest struct {
	FlowID     string            `json:"flowID"`
	FlowParams map[string]string `json:"flowParams"`
	Message    string            `json:"message"`
}

func decodeCreateInteractionAPIRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createInteractionAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "decodeCreateInteractionRequest")
	}

	return Input{
		FlowID:       req.FlowID,
		InputMessage: req.Message,
		Params:       req.FlowParams,
	}, nil
}

type interactionAPIResponse struct {
	ID               string    `json:"id"`
	FlowID           string    `json:"flowID"`
	FlowName         string    `json:"flowName"`
	Type             string    `json:"type"`
	Model            string    `json:"model"`
	Prompt           string    `json:"prompt"`
	Completion       string    `json:"completion"`
	TokensPrompt     int       `json:"tokensPrompt"`
	TokensCompletion int       `json:"tokensCompletion"`
	CreatedAt        time.Time `json:"createdAt"`
	UserID           string    `json:"userID"`
}

func formatInteractionAPIResponse(ctx context.Context, response interface{}) (json.RawMessage, error) {
	ixn, ok := response.(interaction.Interaction)
	if !ok {
		return nil, fmt.Errorf("formatInteractionAPIResponse: failed to parse response")
	}

	formatted := map[string]interactionAPIResponse{
		"data": {
			ID:               ixn.ID,
			FlowID:           ixn.FlowID,
			FlowName:         ixn.FlowName,
			Type:             ixn.Type,
			Model:            ixn.Model,
			Prompt:           ixn.Prompt,
			Completion:       ixn.Completion,
			TokensPrompt:     ixn.TokensPrompt,
			TokensCompletion: ixn.TokensCompletion,
			CreatedAt:        ixn.CreatedAt,
			UserID:           ixn.UserID,
		},
	}

	encoded, err := json.Marshal(formatted)
	if err != nil {
		return nil, errors.Wrap(err, "formatInteractionAPIResponse")
	}

	return encoded, nil
}

type summaryAPIResponse struct {
	ID             string    `json:"id"`
	FlowID         string    `json:"flowID"`
	FlowName       string    `json:"flowName"`
	Type           string    `json:"type"`
	Model          string    `json:"model"`
	Prompt         string    `json:"prompt"`
	TokensUsed     int       `json:"tokensUsed"`
	ResponseLength int       `json:"responseLength"`
	CreatedAt      time.Time `json:"createdAt"`
	UserID         string    `json:"userID"`
}

func formatInteractionSummariesAPIResponse(ctx context.Context, response interface{}) (json.RawMessage, error) {
	interactions, ok := response.([]interaction.Summary)
	if !ok {
		return nil, fmt.Errorf("formatInteractionSummariesAPIResponse: failed to parse response")
	}

	formatted := make([]summaryAPIResponse, len(interactions))
	for i, ixn := range interactions {
		formatted[i] = summaryAPIResponse{
			ID:             ixn.ID,
			FlowID:         ixn.FlowID,
			FlowName:       ixn.FlowName,
			Type:           ixn.Type,
			Model:          ixn.Model,
			Prompt:         ixn.PromptFragment(250),
			TokensUsed:     ixn.TokensUsed,
			ResponseLength: ixn.ResponseLength,
			CreatedAt:      ixn.CreatedAt,
			UserID:         ixn.UserID,
		}
	}

	res := map[string][]summaryAPIResponse{
		"data": formatted,
	}

	encoded, err := json.Marshal(res)
	if err != nil {
		return nil, errors.Wrap(err, "formatInteractionSummariesAPIResponse")
	}

	return encoded, nil
}
