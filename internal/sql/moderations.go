package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/itsnoproblem/prmry/internal/moderation"
)

type moderationsRepo struct {
	db *sqlx.DB
}

func NewModerationsRepo(db *sqlx.DB) moderationsRepo {
	return moderationsRepo{
		db: db,
	}
}

func (r *moderationsRepo) Insert(ctx context.Context, mod moderation.Moderation) error {

	resJSON, err := json.Marshal(mod.Results)
	if err != nil {
		return fmt.Errorf("sql.moderationsRepo: %s", err)
	}

	sql := `
		INSERT INTO moderations (
			id,
			interaction_id,
			model,
			results,
			created_at
		) VALUES (?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, sql,
		mod.ID,
		mod.InteractionID,
		mod.Model,
		resJSON,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("sql.moderationsRepo: %s", err)
	}

	return nil
}

func (r *moderationsRepo) Remove(ctx context.Context, id string) error {
	return fmt.Errorf("moderationsRepo.Delete is not implemented")
}

func (r *moderationsRepo) All(ctx context.Context) ([]moderation.Moderation, error) {
	query := `
		SELECT 
		    id,
			interaction_id,
			model,
			results,
			created_at
		FROM moderations
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("sql.interactionsRepo: %s", err)
	}

	moderations := make([]moderation.Moderation, 0)

	for rows.Next() {
		var mod moderation.Moderation
		if err := rows.Scan(&mod.ID, &mod.InteractionID, &mod.Model, &mod.Results, &mod.CreatedAt); err != nil {
			return nil, fmt.Errorf("sql.moderationsRepo: %s", err)
		}
		moderations = append(moderations, mod)
	}

	return moderations, nil
}
