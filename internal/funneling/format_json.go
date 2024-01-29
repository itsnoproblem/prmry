package funneling

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itsnoproblem/prmry/internal/flow"
)

type executeFunnelResponse struct {
	FlowName    string  `json:"flowName"`
	FlowID      string  `json:"flowID"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	Executes    bool    `json:"executes"`
	Prompt      string  `json:"prompt"`
}

func formatFunnelResponse(ctx context.Context, response interface{}) (json.RawMessage, error) {
	res, ok := response.(flow.Execution)
	if !ok {
		return nil, fmt.Errorf("formatFunnelResponse: failed to cast response to flow.Execution")
	}

	formatted := executeFunnelResponse{
		FlowName:    res.FlowName,
		FlowID:      res.FlowID,
		Model:       res.Model,
		Temperature: res.Temperature,
		Executes:    res.Executes,
		Prompt:      res.Prompt,
	}

	encoded, err := json.Marshal(formatted)
	if err != nil {
		return nil, fmt.Errorf("formatFunnelResponse: failed to marshal response")
	}

	return encoded, nil
}
