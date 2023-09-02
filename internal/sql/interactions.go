package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	openai "github.com/sashabaranov/go-openai"

	"github.com/itsnoproblem/prmry/internal/interaction"
)

type summaryRow struct {
	ID             string    `db:"id"`
	UserID         string    `db:"user_id"`
	Type           string    `db:"type"`
	Model          string    `db:"model"`
	Prompt         string    `db:"prompt"`
	TokensUsed     int       `db:"tokens_used"`
	ResponseLength int       `db:"response_length"`
	FlowID         string    `db:"flow_id"`
	FlowName       string    `db:"flow_name"`
	CreatedAt      time.Time `db:"created_at"`
}

func (row summaryRow) toSummary() interaction.Summary {
	return interaction.Summary{
		ID:             row.ID,
		FlowID:         row.FlowID,
		FlowName:       row.FlowName,
		Type:           row.Type,
		Model:          row.Model,
		Prompt:         row.Prompt,
		TokensUsed:     row.TokensUsed,
		ResponseLength: row.ResponseLength,
		CreatedAt:      row.CreatedAt,
		UserID:         row.UserID,
	}
}

type interactionRow struct {
	ID               string          `db:"id"`
	Type             string          `db:"type"`
	Model            string          `db:"model"`
	Prompt           string          `db:"prompt"`
	Completion       string          `db:"completion"`
	CompletionTokens int             `db:"completion_tokens"`
	PromptTokens     int             `db:"prompt_tokens"`
	Request          json.RawMessage `db:"request"`
	Response         json.RawMessage `db:"response"`
	CreatedAt        time.Time       `db:"created_at"`
	UserID           string          `db:"user_id"`
	FlowID           string          `db:"flow_id"`
	FlowName         string          `db:"flow_name"`
}

func (row *interactionRow) toInteraction() (interaction.Interaction, error) {
	var (
		req openai.ChatCompletionRequest
		res openai.ChatCompletionResponse
	)

	if err := json.Unmarshal(row.Request, &req); err != nil {
		return interaction.Interaction{}, errors.Wrap(err, "toInteraction")
	}

	if err := json.Unmarshal(row.Response, &res); err != nil {
		return interaction.Interaction{}, errors.Wrap(err, "toInteraction")
	}

	prompt := row.Prompt
	if prompt == "" && len(req.Messages) > 0 {
		prompt = req.Messages[0].Content
	}

	completion := row.Completion
	if completion == "" && len(res.Choices) > 0 {
		completion = res.Choices[0].Message.Content
	}

	return interaction.Interaction{
		ID:               row.ID,
		Type:             row.Type,
		Model:            row.Model,
		Prompt:           prompt,
		Completion:       completion,
		TokensCompletion: row.CompletionTokens,
		TokensPrompt:     row.PromptTokens,
		Request:          req,
		Response:         res,
		CreatedAt:        row.CreatedAt,
		UserID:           row.UserID,
		FlowID:           row.FlowID,
		FlowName:         row.FlowName,
	}, nil
}

type interactionsRepo struct {
	db *sqlx.DB
}

func NewInteractionsRepo(db *sqlx.DB) interactionsRepo {
	return interactionsRepo{
		db: db,
	}
}

func (r *interactionsRepo) Insert(ctx context.Context, ixn interaction.Interaction) (err error) {
	reqJSON, err := json.Marshal(ixn.Request)
	if err != nil {
		return fmt.Errorf("sql.interactionsRepo: %s", err)
	}

	resJSON, err := json.Marshal(ixn.Response)
	if err != nil {
		return fmt.Errorf("sql.interactionsRepo: %s", err)
	}

	sql := `
		INSERT INTO interactions (
			id,
			request,
			response,
			created_at,
		    user_id, 
		    flow_id,
		    type,
		    model,
		    prompt,
		    completion,
		    prompt_tokens,
		    completion_tokens
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, sql,
		ixn.ID,
		reqJSON,
		resJSON,
		time.Now(),
		ixn.UserID,
		ixn.FlowID,
		ixn.Type,
		ixn.Model,
		ixn.Prompt,
		ixn.Completion,
		ixn.TokensPrompt,
		ixn.TokensCompletion,
	)
	if err != nil {
		return fmt.Errorf("sql.interactionsRepo: %s", err)
	}

	return nil
}

func (r *interactionsRepo) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("sql.interactionsRepo.Delete: not implemented")
}

func (r *interactionsRepo) GetInteraction(ctx context.Context, id string) (interaction.Interaction, error) {
	query := `
		SELECT 
			i.id,
			i.request, 
			i.response, 
			i.created_at,
			i.user_id,
			i.flow_id,
			i.type,
			i.model,
			i.prompt,
			i.completion,
			i.prompt_tokens,
			i.completion_tokens,
			IF(f.name IS NULL, '', f.name) AS flow_name
		FROM interactions AS i
		LEFT JOIN flows AS f ON f.id = i.flow_id
		WHERE i.id = ?
		LIMIT 1
	`

	var ixnRow interactionRow
	row := r.db.QueryRowxContext(ctx, query, id)
	if err := row.StructScan(&ixnRow); err != nil {
		return interaction.Interaction{}, errors.Wrap(err, "interactionsRepo.interaction")
	}
	if row == nil {
		return interaction.Interaction{}, fmt.Errorf("interactionsRepo.interaction: not found for id %s", id)
	}

	ixn, err := ixnRow.toInteraction()
	if err != nil {
		return interaction.Interaction{}, errors.Wrap(err, "interactionsRepo.interaction")
	}

	return ixn, nil
}

func (r *interactionsRepo) Summaries(ctx context.Context) ([]interaction.Summary, error) {
	query := `
		SELECT 
		  i.id,
		  i.user_id,
		  i.type,
		  i.model,
		  i.prompt,
		  (i.completion_tokens + i.prompt_tokens) AS tokens_used,
		  LENGTH(i.response) AS response_length,
		  i.flow_id,
		  if(f.name IS NULL, '', f.name) AS flow_name,
		  i.created_at
		FROM interactions AS i
		LEFT JOIN flows AS f ON f.id = i.flow_id 
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "sql.interactionsRepo.Summaries")
	}
	defer rows.Close()

	interactions := make([]interaction.Summary, 0)
	for rows.Next() {
		var result summaryRow
		if err = rows.StructScan(&result); err != nil {
			return nil, fmt.Errorf("sql.interactionsRepo: %s", err)
		}

		interactions = append(interactions, result.toSummary())
	}

	return interactions, nil
}

func (r *interactionsRepo) GetInteractionSummaries(ctx context.Context, userID string) ([]interaction.Summary, error) {
	query := `
		SELECT 
		  i.id,
		  i.user_id,
		  i.type,
		  i.model,
		  i.prompt,
		  i.flow_id,
		  i.created_at,
		  (i.prompt_tokens + i.completion_tokens) AS tokens_used,
		  IFNULL(LENGTH(i.response->>'$.choices[0].text'), 0) AS response_length,
		  IF(f.name IS NULL, '', f.name) AS flow_name
		FROM interactions AS i
		LEFT JOIN flows AS f ON f.id = i.flow_id 
		WHERE i.user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, errors.Wrap(err, "sql.interactionsRepo.GetInteractionSummaries")
	}
	defer rows.Close()

	interactions := make([]interaction.Summary, 0)
	for rows.Next() {
		var result summaryRow
		if err = rows.StructScan(&result); err != nil {
			return nil, fmt.Errorf("sql.interactionsRepo: %s", err)
		}

		interactions = append(interactions, result.toSummary())
	}

	return interactions, nil
}
