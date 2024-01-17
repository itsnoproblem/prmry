package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/flow"
)

type repository struct {
	db *sqlx.DB
}

func NewFlowsRepository(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error) {
	sql := `
		SELECT 
		    id,
		    user_id,
		    name,
		    model,
		    temperature,
		    rules,
		    require_all,
		    prompt,
		    prompt_args,
		    inputs,
		    created_at,
		    updated_at
		FROM flows
		WHERE user_id = ?
		ORDER BY updated_at DESC
	`
	result, err := r.db.QueryxContext(ctx, sql, userID)
	defer result.Close()
	if err != nil {
		return nil, errors.Wrap(err, "sql.flows.GetFlowsForUser")
	}

	rows := make([]flowsRow, 0)
	for result.Next() {
		var row flowsRow
		if err = result.StructScan(&row); err != nil {
			return nil, errors.Wrap(err, "sql.flows.GetFlowsForUser")
		}
		rows = append(rows, row)
	}

	flows := make([]flow.Flow, len(rows))
	for i, row := range rows {
		if flows[i], err = row.toFlow(); err != nil {
			return nil, errors.Wrap(err, "sql.flows.GetFlowsForUser")
		}
	}

	return flows, nil
}

func (r *repository) GetFlow(ctx context.Context, flowID string) (flow.Flow, error) {
	sql := `
		SELECT 
		    id,
		    user_id,
		    name,
		    model,
		    temperature,
		    rules,
		    require_all,
		    prompt,
		    prompt_args,
		    inputs,
		    created_at,
		    updated_at
		FROM flows
		WHERE id = ?
	`
	result := r.db.QueryRowxContext(ctx, sql, flowID)
	if result == nil {
		return flow.Flow{}, fmt.Errorf("sql.flow: no flow found for %s", flowID)
	}

	var row flowsRow
	if err := result.StructScan(&row); err != nil {
		return flow.Flow{}, errors.Wrap(err, "sql.flows.GetFlow")
	}

	flw, err := row.toFlow()
	if err != nil {
		return flow.Flow{}, errors.Wrap(err, "sql.flows.GetFlow")
	}

	return flw, nil
}

func (r *repository) InsertFlow(ctx context.Context, flw flow.Flow) error {
	var requireAll int
	if flw.RequireAll {
		requireAll = 1
	}

	rules, err := json.Marshal(flw.Triggers)
	if err != nil {
		return errors.Wrap(err, "sql.flows.InsertFlow")
	}

	promptArgs, err := json.Marshal(flw.PromptArgs)
	if err != nil {
		return errors.Wrap(err, "sql.flows.InsertFlow")
	}

	inputs, err := json.Marshal(flw.InputParams)
	if err != nil {
		return errors.Wrap(err, "sql.flows.InsertFlow")
	}

	sql := `
		INSERT INTO flows (
			id,
			user_id,
		    name,
		    model,
		    temperature,
		    rules,
		    require_all,
		    prompt,
		    prompt_args,
		    inputs,
		    created_at,
		    updated_at
	  	) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)	
	`

	args := []interface{}{
		flw.ID,
		flw.UserID,
		flw.Name,
		flw.Model,
		flw.Temperature,
		rules,
		requireAll,
		flw.Prompt,
		promptArgs,
		inputs,
		flw.CreatedAt,
		flw.UpdatedAt,
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "sql.flows.InsertFlow")
	}

	return nil
}

func (r *repository) UpdateFlow(ctx context.Context, flw flow.Flow) error {
	var requireAll int
	if flw.RequireAll {
		requireAll = 1
	}

	rules, err := json.Marshal(flw.Triggers)
	if err != nil {
		return errors.Wrap(err, "sql.flows.InsertFlow")
	}

	promptArgs, err := json.Marshal(flw.PromptArgs)
	if err != nil {
		return errors.Wrap(err, "sql.flows.InsertFlow")
	}

	inputs, err := json.Marshal(flw.InputParams)
	if err != nil {
		return errors.Wrap(err, "sql.flows.InsertFlow")
	}

	sql := `
		UPDATE flows SET 
		    name = ?,
		    model = ?,
		    temperature = ?,
		    rules = ?,
		    require_all = ?,
		    prompt = ?,
		    prompt_args = ?,
		    inputs = ?,
		    updated_at = ?
	  	WHERE id = ?	
	`

	args := []interface{}{
		flw.Name,
		flw.Model,
		flw.Temperature,
		rules,
		requireAll,
		flw.Prompt,
		promptArgs,
		inputs,
		flw.UpdatedAt,
		flw.ID,
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "sql.flows.UpdateFlow")
	}

	return nil
}

func (r repository) DeleteFlow(ctx context.Context, flowID string) error {
	sql := `
		DELETE FROM flows WHERE id = ?
	`
	if _, err := r.db.ExecContext(ctx, sql, flowID); err != nil {
		return errors.Wrap(err, "sql.DeleteFLow")
	}

	return nil
}

// private

type flowsRow struct {
	ID          string           `db:"id"`
	UserID      string           `db:"user_id"`
	Name        string           `db:"name"`
	Model       string           `db:"model"`
	Temperature float64          `db:"temperature"`
	Rules       *json.RawMessage `db:"rules"`
	RequireAll  bool             `db:"require_all"`
	Prompt      string           `db:"prompt"`
	PromptArgs  *json.RawMessage `db:"prompt_args"`
	Inputs      *json.RawMessage `db:"inputs"`
	CreatedAt   time.Time        `db:"created_at"`
	UpdatedAt   time.Time        `db:"updated_at"`
}

func (r flowsRow) toFlow() (flow.Flow, error) {
	rules := make([]flow.Trigger, 0)
	if r.Rules != nil {
		if err := json.Unmarshal(*r.Rules, &rules); err != nil {
			return flow.Flow{}, errors.Wrap(err, "toFlow: rules")
		}
	}

	promptArgs := make([]flow.Field, 0)
	if r.PromptArgs != nil {
		if err := json.Unmarshal(*r.PromptArgs, &promptArgs); err != nil {
			return flow.Flow{}, errors.Wrap(err, "toFlow: promptArgs")
		}
	}

	inputs := make([]flow.InputParam, 0)
	if r.Inputs != nil {
		if err := json.Unmarshal(*r.Inputs, &inputs); err != nil {
			return flow.Flow{}, errors.Wrap(err, "toFlow: inputs")
		}
	}

	model := flow.DefaultModel
	if r.Model != "" {
		model = r.Model
	}

	temp := flow.DefaultTemperature
	if r.Temperature != 0 {
		temp = r.Temperature
	}

	return flow.Flow{
		ID:          r.ID,
		UserID:      r.UserID,
		Name:        r.Name,
		Model:       model,
		Temperature: temp,
		Triggers:    rules,
		RequireAll:  r.RequireAll,
		Prompt:      r.Prompt,
		PromptArgs:  promptArgs,
		InputParams: inputs,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}, nil
}
