package sql

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/funnel"
)

type funnelsRepo struct {
	db *sqlx.DB
}

func NewFunnelsRepository(db *sqlx.DB) *funnelsRepo {
	return &funnelsRepo{
		db: db,
	}
}

type funnelRow struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Name      string    `db:"name"`
	Path      string    `db:"path"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r funnelRow) toFunnel() funnel.Funnel {
	return funnel.Funnel{
		ID:        r.ID,
		UserID:    r.UserID,
		Name:      r.Name,
		Path:      r.Path,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
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
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Name      string    `db:"name"`
	Path      string    `db:"path"`
	FlowCount int       `db:"flows"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s funnelSummary) toFunnelSummary() (funnel.Summary, error) {
	return funnel.Summary{
		ID:        s.ID,
		Name:      s.Name,
		Path:      s.Path,
		FlowCount: s.FlowCount,
		UpdatedAt: time.Time{},
	}, nil
}

type funnelFlowRow struct {
	FunnelID string `db:"funnel_id"`
	FlowID   string `db:"flow_id"`
	Name     string `db:"name"`
}

func (r *funnelsRepo) GetFunnelSummariesForUser(ctx context.Context, userID string) ([]funnel.Summary, error) {
	sql := `
		SELECT 
		    f.id,
		    f.user_id,
		    f.name,
		    f.path,
		    f.updated_at,
		    (SELECT count(*) FROM funnel_flows AS ff WHERE ff.funnel_id = f.id) AS flows
		FROM funnels AS f
		WHERE f.user_id = ?
		ORDER BY f.updated_at DESC
	`
	result, err := r.db.QueryxContext(ctx, sql, userID)
	if err != nil {
		return nil, errors.Wrap(err, "sql.funnels.GetFunnelSummariesForUser")
	}
	defer result.Close()

	rows := make([]funnelSummary, 0)
	for result.Next() {
		var row funnelSummary
		if err = result.StructScan(&row); err != nil {
			return nil, errors.Wrap(err, "sql.funnels.GetFunnelSummariesForUser")
		}
		rows = append(rows, row)
	}

	funnels := make([]funnel.Summary, len(rows))
	for i, row := range rows {
		if funnels[i], err = row.toFunnelSummary(); err != nil {
			return nil, errors.Wrap(err, "sql.funnels.GetFunnelSummariesForUser")
		}
	}

	return funnels, nil
}

func (r *funnelsRepo) GetFunnel(ctx context.Context, id string) (funnel.Funnel, error) {
	sql := `
		SELECT 
		    f.id,
		    f.user_id,
		    f.name,
		    f.path,
		    f.created_at,
		    f.updated_at
		FROM funnels AS f
		WHERE f.id = ?
		LIMIT 1
	`
	result, err := r.db.QueryxContext(ctx, sql, id)
	defer result.Close()
	if err != nil {
		return funnel.Funnel{}, errors.Wrap(err, "sql.funnels.GetFunnel")
	}

	var row funnelRow
	for result.Next() {

		if err = result.StructScan(&row); err != nil {
			return funnel.Funnel{}, errors.Wrap(err, "sql.funnels.GetFunnel")
		}

		return row.toFunnel(), nil
	}

	return funnel.Funnel{}, errors.New("sql.funnels.GetFunnel: no funnel found")
}

func (r *funnelsRepo) GetFunnelByPath(ctx context.Context, path string) (funnel.Funnel, error) {
	sql := `
		SELECT 
		    f.id,
		    f.user_id,
		    f.name,
		    f.path,
		    f.created_at,
		    f.updated_at
		FROM funnels AS f
		WHERE f.path = ?
		LIMIT 1
	`
	result, err := r.db.QueryxContext(ctx, sql, path)
	defer result.Close()
	if err != nil {
		return funnel.Funnel{}, errors.Wrap(err, "sql.funnels.GetFunnelByPath")
	}

	var row funnelRow
	for result.Next() {

		if err = result.StructScan(&row); err != nil {
			return funnel.Funnel{}, errors.Wrap(err, "sql.funnels.GetFunnelByPath")
		}

		return row.toFunnel(), nil
	}

	return funnel.Funnel{}, errors.New("sql.funnels.GetFunnelByPath: no funnel found")
}

func (r *funnelsRepo) InsertFunnel(ctx context.Context, f funnel.Funnel) error {
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
		return errors.Wrap(err, "sql.funnels.InsertFunnel")
	}

	return nil
}

func (r *funnelsRepo) UpdateFunnel(ctx context.Context, f funnel.Funnel) error {
	row := newFunnelRow(f)

	sql := `
		UPDATE funnels 
		SET name = :name,
			path = :path,
			updated_at = :updated_at
		WHERE user_id = :user_id AND id = :id
	`
	_, err := r.db.NamedExecContext(ctx, sql, row)
	if err != nil {
		return errors.Wrap(err, "sql.funnels.InsertFunnel")
	}

	return nil
}

func (r *funnelsRepo) ExistsFunnelFlow(ctx context.Context, funnelID, flowID string) (bool, error) {
	sql := `
		SELECT EXISTS (
			SELECT 1 FROM funnel_flows WHERE funnel_id = ? AND flow_id = ?
		)
	`
	var exists bool
	if err := r.db.GetContext(ctx, &exists, sql, funnelID, flowID); err != nil {
		return false, errors.Wrap(err, "sql.funnels.ExistsFunnelFlow")
	}

	return exists, nil
}

func (r *funnelsRepo) AddFlowsToFunnel(ctx context.Context, funnelID string, flowIDs ...string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "sql.funnels.AddFlowsToFunnel: begin transaction")
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
		row := funnelFlowRow{
			FunnelID: funnelID,
			FlowID:   flowID,
		}
		if _, err := tx.NamedExecContext(ctx, sql, row); err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return errors.Wrap(rollbackErr, "sql.funnels.AddFlowsToFunnel: rollback due to error")
			}

			return errors.Wrap(err, "sql.funnels.AddFlowsToFunnel")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "sql.funnels.AddFlowsToFunnel: commit transaction")
	}

	return nil
}

func (r *funnelsRepo) RemoveFlowsFromFunnel(ctx context.Context, funnelID string, flowIDs ...string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "sql.funnels.RemoveFlowsFromFunnel: begin transaction")
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
				return errors.Wrap(rollbackErr, "sql.funnels.RemoveFlowsFromFunnel: rollback due to error")
			}
			return errors.Wrap(err, "sql.funnels.RemoveFlowsFromFunnel")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "sql.funnels.RemoveFlowsFromFunnel: commit transaction")
	}

	return nil
}

func (r *funnelsRepo) DeleteFunnel(ctx context.Context, funnelID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "sql.funnels.transactionDeleteFunnelAndFunnelFlows: begin transaction")
	}

	sql := `
		DELETE FROM funnels WHERE id = ?
	`
	if _, err := tx.ExecContext(ctx, sql, funnelID); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.Wrap(rollbackErr, "sql.funnels.transactionDeleteFunnelAndFunnelFlows: rollback due to error")
		}
		return errors.Wrap(err, "sql.funnels.transactionDeleteFunnelAndFunnelFlows")
	}

	sql = `
		DELETE FROM funnel_flows WHERE funnel_id = ?
	`
	if _, err := tx.ExecContext(ctx, sql, funnelID); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.Wrap(rollbackErr, "sql.funnels.transactionDeleteFunnelAndFunnelFlows: rollback due to error")
		}
		return errors.Wrap(err, "sql.funnels.transactionDeleteFunnelAndFunnelFlows")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "sql.funnels.transactionDeleteFunnelAndFunnelFlows: commit transaction")
	}

	return nil
}
