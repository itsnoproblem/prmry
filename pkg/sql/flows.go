package sql

import (
	"context"
	"encoding/json"
	"github.com/itsnoproblem/prmry/pkg/flow"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type repository struct {
	db *sqlx.DB
}

func NewFlowsRepository(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}

type flowsRow struct {
	ID         string          `db:"id"`
	UserID     string          `db:"user_id"`
	Name       string          `db:"name"`
	Rules      json.RawMessage `db:"rules"`
	RequireAll bool            `db:"require_all"`
	Prompt     string          `db:"prompt"`
	PromptArgs json.RawMessage `db:"prompt_args"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
}

func (r *repository) InsertFlow(ctx context.Context, flw flow.Flow) error {
	var requireAll int
	if flw.RequireAll {
		requireAll = 1
	}

	rules, err := json.Marshal(flw.Rules)
	if err != nil {
		return errors.Wrap(err, "sql.flows.InsertFlow")
	}

	promptArgs, err := json.Marshal(flw.PromptArgs)
	if err != nil {
		return errors.Wrap(err, "sql.flows.InsertFlow")
	}

	sql := `
		INSERT INTO flows (
			id,
			user_id,
		    name,
		    rules,
		    require_all,
		    prompt,
		    prompt_args,
		    created_at,
		    updated_at
	  	) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)	
	`

	args := []interface{}{
		flw.ID,
		flw.UserID,
		flw.Name,
		rules,
		requireAll,
		flw.Prompt,
		promptArgs,
		flw.CreatedAt,
		flw.UpdatedAt,
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "sql.flows.InsertFlow")
	}

	return nil
}
