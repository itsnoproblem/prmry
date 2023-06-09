package interacting

import (
	"fmt"
	"time"

	"github.com/itsnoproblem/prmry/pkg/interaction"
)

type interactionSummaryResponse struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Model          string    `json:"model"`
	Prompt         string    `json:"prompt"`
	TokensUsed     int       `json:"tokens_used"`
	ResponseLength int       `json:"response_length"`
	CreatedAt      time.Time `json:"created_at"`
}

func formatInteractionSummaries(summaries interface{}) (interface{}, error) {
	sum, ok := summaries.([]interaction.Summary)
	if !ok {
		return nil, fmt.Errorf("formatInteractionSummaries: failed to parse response")
	}

	res := make([]interactionSummaryResponse, len(sum))
	for i, s := range sum {
		res[i] = interactionSummaryResponse{
			ID:             s.ID,
			Type:           s.Type,
			Model:          s.Model,
			Prompt:         s.Prompt,
			TokensUsed:     s.TokensUsed,
			ResponseLength: s.ResponseLength,
			CreatedAt:      s.CreatedAt,
		}
	}
	return res, nil
}

type interactionResponse struct {
	ID        string      `json:"id"`
	Request   interface{} `json:"request"`
	Response  interface{} `json:"response"`
	Error     string      `json:"error"`
	CreatedAt time.Time   `json:"created_at"`
}

func formatGetInteractionResponse(in interface{}) (interface{}, error) {
	res, ok := in.(interaction.Interaction)
	if !ok {
		return nil, fmt.Errorf("formatGetInteractionResponse: failed to parse response")
	}

	return interactionResponse{
		ID:        res.ID,
		Request:   res.Request,
		Response:  res.Response,
		Error:     res.Error,
		CreatedAt: res.CreatedAt,
	}, nil

}
