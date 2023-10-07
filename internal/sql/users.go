package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
)

type UserRow struct {
	ID        string `db:"id"`
	Email     string `db:"email"`
	Name      string `db:"name"`
	Nickname  string `db:"nickname"`
	AvatarURL string `db:"avatar_url"`
}

func (r UserRow) ToUser() auth.User {
	return auth.User{
		ID:        r.ID,
		Email:     r.Email,
		Name:      r.Name,
		Nickname:  r.Nickname,
		AvatarURL: r.AvatarURL,
	}
}

type usersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) usersRepo {
	return usersRepo{
		db: db,
	}
}

func (r usersRepo) InsertUser(ctx context.Context, usr auth.User) error {
	//id = uuid.NewCookie().String()
	query := `
		INSERT INTO 
		    users (
			   id, 
			   email, 
			   name, 
			   nickname, 
			   avatar_url
		   ) VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, usr.ID, usr.Email, usr.Name, usr.Nickname, usr.AvatarURL)
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

func (r usersRepo) FindUserViaOAuth(ctx context.Context, provider, providerUserID string) (usr auth.User, exists bool, err error) {
	query := `
		SELECT 
			u.id, 
			u.email, 
			u.name,
			u.nickname,
			u.avatar_url
		FROM oauth_users AS oa
		INNER JOIN users AS u ON u.id = oa.user_id
		WHERE oa.provider = ?
		AND oa.provider_user_id = ?
		LIMIT 1
	`

	rows, err := r.db.QueryxContext(ctx, query, provider, providerUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return auth.User{}, false, nil
		}
		return auth.User{}, false, errors.Wrap(err, "usersRepo.FindUserViaOAuth")
	}

	for rows.Next() {
		var userRow UserRow
		if err = rows.StructScan(&userRow); err != nil {
			return auth.User{}, false, errors.Wrap(err, "usersRepo.FindUserViaOAuth")
		}
	}

	return usr, true, nil
}

func (r usersRepo) FindUserByEmail(ctx context.Context, email string) (usr auth.User, exists bool, err error) {
	query := `
		SELECT 
			id, 
			email, 
			name,
			nickname,
			avatar_url
		FROM users
		WHERE email = ?
	`

	var userRow UserRow
	if err := r.db.QueryRowxContext(ctx, query, email).StructScan(&userRow); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return auth.User{}, false, nil
		}
		return auth.User{}, false, errors.Wrapf(err, "usersRepo.FindUserByEmail")
	}

	return userRow.ToUser(), true, nil
}

func (r usersRepo) FindUserByAPIKey(ctx context.Context, key string) (usr auth.User, exists bool, err error) {
	if key == "" {
		return auth.User{}, false, nil
	}

	query := `
		SELECT 
			id, 
			email, 
			name,
			nickname,
			avatar_url
		FROM users
		WHERE api_key = ?
	`

	var userRow UserRow
	if err := r.db.QueryRowxContext(ctx, query, key).StructScan(&userRow); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return auth.User{}, false, nil
		}
		return auth.User{}, false, errors.Wrapf(err, "usersRepo.FindUserByAPIKey")
	}

	return userRow.ToUser(), true, nil
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

func (r usersRepo) SaveUserFromOAuth(ctx context.Context, usr auth.User, oauthProvider, providerUserID string) error {
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
