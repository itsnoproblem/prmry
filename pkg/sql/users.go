package sql

import (
	"context"
	"database/sql"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/user"
	"time"

	"github.com/pkg/errors"
)

type usersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) usersRepo {
	return usersRepo{
		db: db,
	}
}

func (r usersRepo) InsertUser(ctx context.Context, usr user.User) error {
	//id = uuid.New().String()
	query := `
		INSERT INTO users (id, name, email) VALUES (?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, usr.ID, usr.Name, usr.Email)
	if err != nil {
		return errors.Wrapf(err, "usersRepo.InsertUser")
	}

	return nil
}

func (r usersRepo) DeleteUser(ctx context.Context, id string) error {
	userQuery := "DELETE FROM users WHERE id = ?"
	oauthQuery := "DELETE FROM oauth_users WHERE user_id = ?"

	if _, err := r.db.ExecContext(ctx, userQuery, id); err != nil {
		return errors.Wrap(err, "usersRepo.DeleteUser")
	}

	if _, err := r.db.ExecContext(ctx, oauthQuery, id); err != nil {
		return errors.Wrap(err, "usersRepo.DeleteUser")
	}

	return nil
}

func (r usersRepo) FindUserViaOAuth(ctx context.Context, provider, providerUserID string) (usr user.User, exists bool, err error) {
	query := `
		SELECT 
			u.id, 
			u.email, 
			u.name
		FROM oauth_users AS oa
		INNER JOIN users AS u ON u.id = oa.user_id
		WHERE oa.provider = ?
		AND oa.provider_user_id = ?
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, provider, providerUserID)
	err = row.Scan(&usr.ID, &usr.Email, &usr.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, false, nil
		}
		return user.User{}, false, errors.Wrapf(err, "usersRepo.FindUserViaOAuth")
	}

	return usr, true, nil
}

func (r usersRepo) FindUserByEmail(ctx context.Context, email string) (usr user.User, exists bool, err error) {
	query := `
		SELECT 
			id, 
			email, 
			name
		FROM users
		WHERE email = ?
	`

	row := r.db.QueryRowContext(ctx, query, email)
	err = row.Scan(&usr.ID, &usr.Email, &usr.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, false, nil
		}
		return user.User{}, false, errors.Wrapf(err, "usersRepo.FindUserByEmail")
	}

	return usr, true, nil
}

func (r usersRepo) ExistsUserViaOAuth(ctx context.Context, provider, providerUserID string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM oauth_users 
		WHERE provider = ? and provider_user_id = ?
	`
	var numRows int
	err := r.db.QueryRowContext(ctx, query, provider, providerUserID).Scan(&numRows)
	if err != nil {
		return false, errors.Wrap(err, "usersRepo.ExistsViaOAuth")
	}
	return numRows > 0, nil
}

func (r usersRepo) ExistsUserByID(ctx context.Context, userID string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM users 
		WHERE id = ?
	`
	var numRows int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&numRows)
	if err != nil {
		return false, errors.Wrap(err, "usersRepo.ExistsUserByID")
	}
	return numRows > 0, nil
}

func (r usersRepo) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM users 
		WHERE email = ?
	`
	var numRows int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&numRows)
	if err != nil {
		return false, errors.Wrap(err, "usersRepo.ExistsUserByEmail")
	}
	return numRows > 0, nil
}

func (r usersRepo) SaveUserFromOAuth(ctx context.Context, usr user.User, oauthProvider, providerUserID string) error {
	query := `
		INSERT INTO oauth_users (
			user_id, 
			provider_user_id, 
		    provider, 
		    last_authenticated_at
		) 
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE last_authenticated_at = ?
	`
	args := []interface{}{
		usr.ID,
		providerUserID,
		oauthProvider,
		time.Now(),
		time.Now(),
	}
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrapf(err, "usersRepo.SaveUserFromOAuth")
	}

	return nil
}
