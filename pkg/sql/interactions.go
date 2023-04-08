package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/interaction"
	"github.com/jmoiron/sqlx"
	"time"

	"github.com/google/uuid"
)

type interactionsRepo struct {
	db *sqlx.DB
}

func NewInteractionsRepo(db *sqlx.DB) interactionsRepo {
	return interactionsRepo{
		db: db,
	}
}

func (r *interactionsRepo) Add(ctx context.Context, in interaction.Interaction) (id string, err error) {
	id = uuid.New().String()

	reqJSON, err := json.Marshal(in.Request)
	if err != nil {
		return "", fmt.Errorf("sql.interactionsRepo: %s", err)
	}

	resJSON, err := json.Marshal(in.Response)
	if err != nil {
		return "", fmt.Errorf("sql.interactionsRepo: %s", err)
	}

	sql := `
		INSERT INTO interactions (
			id,
			request,
			response,
			err,
			created_at
		) VALUES (?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, sql,
		id,
		reqJSON,
		resJSON,
		in.Error,
		time.Now(),
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
	var in interaction.Interaction
	row := r.db.QueryRowContext(ctx, query, id)

	var req json.RawMessage
	var res json.RawMessage

	err := row.Scan(&in.ID, &req, &res, &in.Error, &in.CreatedAt)
	if err != nil { // && err == sql.ErrNoRows{
		return interaction.Interaction{}, fmt.Errorf("interactionsRepo.Interaction: %s", err)
	}

	if err = json.Unmarshal(req, &in.Request); err != nil {
		return interaction.Interaction{}, fmt.Errorf("interactionsRepo.Interaction unmarshal request: %s", err)
	}

	if err = json.Unmarshal(res, &in.Response); err != nil {
		return interaction.Interaction{}, fmt.Errorf("interactionsRepo.Interaction unmarshal response: %s", err)
	}

	return in, nil
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
