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

	"github.com/itsnoproblem/prmry/pkg/interaction"
)

type InteractionRow struct {
	ID        string          `db:"id"`
	Request   json.RawMessage `db:"request"`
	Response  json.RawMessage `db:"response"`
	Error     string          `db:"err"`
	CreatedAt time.Time       `db:"created_at"`
	UserID    string          `db:"user_id"`
}

func (row InteractionRow) toInteraction() (interaction.Interaction, error) {
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
		    user_id
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, sql,
		id,
		reqJSON,
		resJSON,
		ixn.Error,
		time.Now(),
		ixn.UserID,
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
			id, 
			request, 
			response, 
			err, 
			created_at
		FROM interactions
		WHERE id = ?
		LIMIT 1
	`
	var (
		ixn interaction.Interaction
		req json.RawMessage
		res json.RawMessage
	)

	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&ixn.ID, &req, &res, &ixn.Error, &ixn.CreatedAt)
	if err != nil { // && err == sql.ErrNoRows{
		return interaction.Interaction{}, fmt.Errorf("interactionsRepo.Interaction: %s", err)
	}

	if err = json.Unmarshal(req, &ixn.Request); err != nil {
		return interaction.Interaction{}, fmt.Errorf("interactionsRepo.Interaction unmarshal request: %s", err)
	}

	if err = json.Unmarshal(res, &ixn.Response); err != nil {
		return interaction.Interaction{}, fmt.Errorf("interactionsRepo.Interaction unmarshal response: %s", err)
	}

	return ixn, nil
}

func (r *interactionsRepo) Summaries(ctx context.Context) ([]interaction.Summary, error) {
	query := `
		SELECT 
		  id, 
		  IFNULL(response->>'$.object', '') AS type,
		  IFNULL(request->>'$.model', '') AS model,
		  IFNULL(request->>'$.prompt', '') AS prompt,
		  IFNULL(response->>'$.usage.total_tokens', 0) AS tokens_used,
		  IFNULL(LENGTH(response->>'$.choices[0].text'), 0) AS response_length,
		  err,
		  created_at
		FROM interactions
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("sql.interactionsRepo: %s", err)
	}
	defer rows.Close()

	interactions := make([]interaction.Summary, 0)
	for rows.Next() {
		var in interaction.Summary

		if err = rows.Scan(
			&in.ID,
			&in.Type,
			&in.Model,
			&in.Prompt,
			&in.TokensUsed,
			&in.ResponseLength,
			&in.Error,
			&in.CreatedAt); err != nil {
			return nil, fmt.Errorf("sql.interactionsRepo: %s", err)
		}

		interactions = append(interactions, in)
	}

	return interactions, nil
}

func (r *interactionsRepo) SummariesForUser(ctx context.Context, userID string) ([]interaction.Summary, error) {
	query := `
		SELECT 
		  id, 
		  IFNULL(response->>'$.object', '') AS type,
		  IFNULL(request->>'$.model', '') AS model,
		  IFNULL(request->>'$.prompt', '') AS prompt,
		  IFNULL(response->>'$.usage.total_tokens', 0) AS tokens_used,
		  IFNULL(LENGTH(response->>'$.choices[0].text'), 0) AS response_length,
		  err,
		  created_at
		FROM interactions
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("sql.interactionsRepo.SummariesForUser(%s): %s", userID, err)
	}
	defer rows.Close()

	interactions := make([]interaction.Summary, 0)
	for rows.Next() {
		var in interaction.Summary

		if err = rows.Scan(
			&in.ID,
			&in.Type,
			&in.Model,
			&in.Prompt,
			&in.TokensUsed,
			&in.ResponseLength,
			&in.Error,
			&in.CreatedAt); err != nil {
			return nil, fmt.Errorf("sql.interactionsRepo: %s", err)
		}

		interactions = append(interactions, in)
	}

	return interactions, nil
}
