package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	gogpt "github.com/sashabaranov/go-gpt3"

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
	Error          string    `db:"err"`
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
		Error:          row.Error,
		CreatedAt:      row.CreatedAt,
		UserID:         row.UserID,
	}
}

type interactionRow struct {
	ID        string          `db:"id"`
	Request   json.RawMessage `db:"request"`
	Response  json.RawMessage `db:"response"`
	Error     string          `db:"err"`
	CreatedAt time.Time       `db:"created_at"`
	UserID    string          `db:"user_id"`
	FlowID    string          `db:"flow_id"`
	FlowName  string          `db:"flow_name"`
}

func (row interactionRow) toInteraction() (interaction.Interaction, error) {
	var (
		req gogpt.CompletionRequest
		res gogpt.CompletionResponse
	)

	if err := json.Unmarshal(row.Request, &req); err != nil {
		return interaction.Interaction{}, errors.Wrap(err, "toInteraction")
	}

	if err := json.Unmarshal(row.Response, &res); err != nil {
		return interaction.Interaction{}, errors.Wrap(err, "toInteraction")
	}

	return interaction.Interaction{
		ID:        row.ID,
		Request:   req,
		Response:  res,
		Error:     row.Error,
		CreatedAt: row.CreatedAt,
		UserID:    row.UserID,
		FlowID:    row.FlowID,
		FlowName:  row.FlowName,
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

func (r *interactionsRepo) Add(ctx context.Context, ixn interaction.Interaction) (id string, err error) {
	id = uuid.New().String()

	reqJSON, err := json.Marshal(ixn.Request)
	if err != nil {
		return "", fmt.Errorf("sql.interactionsRepo: %s", err)
	}

	resJSON, err := json.Marshal(ixn.Response)
	if err != nil {
		return "", fmt.Errorf("sql.interactionsRepo: %s", err)
	}

	sql := `
		INSERT INTO interactions (
			id,
			request,
			response,
			err,
			created_at,
		    user_id, 
		    flow_id
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, sql,
		id,
		reqJSON,
		resJSON,
		ixn.Error,
		time.Now(),
		ixn.UserID,
		ixn.FlowID,
	)
	if err != nil {
		return "", fmt.Errorf("sql.interactionsRepo: %s", err)
	}

	return id, nil
}

func (r *interactionsRepo) Remove(ctx context.Context, id string) error {
	return nil
}

func (r *interactionsRepo) Interaction(ctx context.Context, id string) (interaction.Interaction, error) {
	query := `
		SELECT 
			i.id, 
			i.user_id,
			i.flow_id,
			IF(f.name IS NULL, '', f.name) AS flow_name,
			i.request, 
			i.response, 
			i.err, 
			i.created_at
		FROM interactions AS i
		LEFT JOIN flows AS f ON f.id = i.flow_id
		WHERE i.id = ?
		LIMIT 1
	`

	var ixnRow interactionRow
	row := r.db.QueryRowxContext(ctx, query, id)
	err := row.StructScan(&ixnRow)

	if err != nil { // && err == sql.ErrNoRows{
		return interaction.Interaction{}, errors.Wrap(err, "interactionsRepo.interaction")
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
		  IFNULL(i.response->>'$.object', '') AS type,
		  IFNULL(i.request->>'$.model', '') AS model,
		  IFNULL(i.request->>'$.prompt', '') AS prompt,
		  IFNULL(i.response->>'$.usage.total_tokens', 0) AS tokens_used,
		  IFNULL(LENGTH(i.response->>'$.choices[0].text'), 0) AS response_length,
		  i.err,
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

func (r *interactionsRepo) SummariesForUser(ctx context.Context, userID string) ([]interaction.Summary, error) {
	query := `
		SELECT 
		  i.id,
		  i.user_id,
		  IFNULL(i.response->>'$.object', '') AS type,
		  IFNULL(i.request->>'$.model', '') AS model,
		  IFNULL(i.request->>'$.prompt', '') AS prompt,
		  IFNULL(i.response->>'$.usage.total_tokens', 0) AS tokens_used,
		  IFNULL(LENGTH(i.response->>'$.choices[0].text'), 0) AS response_length,
		  i.err,
		  i.flow_id,
		  IF(f.name IS NULL, '', f.name) AS flow_name,
		  i.created_at
		FROM interactions AS i
		LEFT JOIN flows AS f ON f.id = i.flow_id 
		WHERE i.user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, errors.Wrap(err, "sql.interactionsRepo.SummariesForUser")
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
