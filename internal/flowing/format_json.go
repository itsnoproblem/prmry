package flowing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/flow"
)

type listFlowsAPIResponse struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Trigger        string          `json:"trigger"`
	Prompt         string          `json:"prompt"`
	RequiredParams map[string]bool `json:"requiredParams"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

func formatListFlowsAPIResponse(ctx context.Context, response interface{}) (json.RawMessage, error) {
	res, ok := response.([]flow.Flow)
	if !ok {
		return nil, fmt.Errorf("formatListFlowsAPIResponse: failed to parse response")
	}

	flows := make([]listFlowsAPIResponse, 0)
	for _, flw := range res {
		params := make(map[string]bool)
		for _, param := range flw.InputParams {
			params[param.Key] = param.IsRequired
		}

		flows = append(flows, listFlowsAPIResponse{
			ID:             flw.ID,
			Name:           flw.Name,
			Trigger:        flw.TriggerDescription(),
			Prompt:         flw.Prompt,
			RequiredParams: params,
			UpdatedAt:      flw.UpdatedAt,
		})
	}

	formatted := map[string][]listFlowsAPIResponse{
		"data": flows,
	}

	encoded, err := json.Marshal(formatted)
	if err != nil {
		return nil, errors.Wrap(err, "formatInteractionAPIResponse")
	}

	return encoded, nil
}
