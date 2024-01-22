package sql

import (
	"context"
	"github.com/itsnoproblem/prmry/internal/funnel"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type funnelRow struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Name      string    `db:"name"`
	Path      string    `db:"path"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func newFunnelRow(funnel funnel.Funnel) funnelRow {
	return funnelRow{
		ID:        funnel.ID,
		UserID:    funnel.UserID,
		Name:      funnel.Name,
		Path:      funnel.Path,
		CreatedAt: funnel.CreatedAt,
		UpdatedAt: funnel.UpdatedAt,
	}
}

type funnelSummary struct {
	UserID    string    `db:"user_id"`
	Name      string    `db:"name"`
	Path      string    `db:"path"`
	Flows     int       `db:"flows"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s funnelSummary) toFunnelSummary() (funnel.Summary, error) {
	return funnel.Summary{
		Name:      s.Name,
		Path:      s.Path,
		FlowCount: s.Flows,
		UpdatedAt: time.Time{},
	}, nil
}

type funnelsRepo struct {
	db *sqlx.DB
}

func NewFunnelsRepository(db *sqlx.DB) *funnelsRepo {
	return &funnelsRepo{
		db: db,
	}
}

func (r funnelsRepo) GetFunnelSummariesForUser(ctx context.Context, userID string) ([]funnel.Summary, error) {
	sql := `
		SELECT 
		    f.id,
		    f.user_id,
		    f.name,
		    f.path,
		    f.updated_at,
		    count(*) AS flows
		FROM funnels AS f
		LEFT JOIN funnel_flows AS fl ON fl.funnel_id = f.id
		WHERE f.user_id = ?
		ORDER BY f.updated_at DESC
	`
	result, err := r.db.QueryxContext(ctx, sql, userID)
	defer result.Close()
	if err != nil {
		return nil, errors.Wrap(err, "sql.flows.GetFunnelSummariesForUser")
	}

	rows := make([]funnelSummary, 0)
	for result.Next() {
		var row funnelSummary
		if err = result.StructScan(&row); err != nil {
			return nil, errors.Wrap(err, "sql.flows.GetFunnelSummariesForUser")
		}
		rows = append(rows, row)
	}

	funnels := make([]funnel.Summary, len(rows))
	for i, row := range rows {
		if funnels[i], err = row.toFunnelSummary(); err != nil {
			return nil, errors.Wrap(err, "sql.flows.GetFunnelSummariesForUser")
		}
	}

	return funnels, nil
}

func (r funnelsRepo) InsertFunnel(ctx context.Context, f funnel.Funnel) error {
	newRow := newFunnelRow(f)

	sql := `
		INSERT INTO funnels (
			id,
			user_id,
			name,
			path,
			created_at,
			updated_at
		) VALUES (
			:id,
			:user_id,
			:name,
			:path,
			:created_at,
			:updated_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, sql, newRow)
	if err != nil {
		return errors.Wrap(err, "sql.flows.InsertFunnel")
	}

	return nil
}

func (r funnelsRepo) AddFlowsToFunnel(ctx context.Context, funnelID string, flowIDs ...string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "sql.flows.AddFlowsToFunnel: begin transaction")
	}

	sql := `
		INSERT INTO funnel_flows (
			funnel_id,
			flow_id
		) VALUES (
			:funnel_id,
			:flow_id
		)
	`

	for _, flowID := range flowIDs {
		_, err := tx.NamedExecContext(ctx, sql, map[string]interface{}{
			"funnel_id": funnelID,
			"flow_id":   flowID,
		})
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return errors.Wrap(rollbackErr, "sql.flows.AddFlowsToFunnel: rollback due to error")
			}
			return errors.Wrap(err, "sql.flows.AddFlowsToFunnel")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "sql.flows.AddFlowsToFunnel: commit transaction")
	}

	return nil
}

func (r funnelsRepo) RemoveFlowsFromFunnel(ctx context.Context, funnelID string, flowIDs ...string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "sql.flows.RemoveFlowsFromFunnel: begin transaction")
	}

	sql := `
		DELETE FROM funnel_flows
		WHERE funnel_id = :funnel_id AND flow_id = :flow_id
	`

	for _, flowID := range flowIDs {
		_, err := tx.NamedExecContext(ctx, sql, map[string]interface{}{
			"funnel_id": funnelID,
			"flow_id":   flowID,
		})
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return errors.Wrap(rollbackErr, "sql.flows.RemoveFlowsFromFunnel: rollback due to error")
			}
			return errors.Wrap(err, "sql.flows.RemoveFlowsFromFunnel")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "sql.flows.RemoveFlowsFromFunnel: commit transaction")
	}

	return nil
}
